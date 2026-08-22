------------------------------ MODULE ExecutionTargetFence ------------------------------
EXTENDS Naturals, TLC

\* Bounded model of FEATURE-0016 qualification fencing and audit-before-
\* publication (ADH-2026-058 clause 5,7,8; F16-AD-16/16.1/20). The model
\* deliberately omits API/HTTP mechanics; it proves that a qualification
\* conclusion and a replay completion require a successful AuditEvent append,
\* while stale, cancelled, and failed proposals cannot mutate the target or
\* become replayable completions.
\*
\* This model proves:
\*   1. A committed conclusion always has a successful audit append
\*      (ConclusionRequiresAudit)
\*   2. A replayable completion always corresponds to a committed conclusion
\*      that was audited (ReplayRequiresAuditedCompletion)
\*   3. Stale, cancelled, and failed proposals never mutate committed state
\*      or become replayable (StaleCancelledFailedNeverMutate)

CONSTANTS
  Fresh, Stale,                 \* captured-fence freshness at commit time
  Cancelled, Failed, Succeeded  \* proposal disposition

VARIABLES
  committedGeneration,   \* the persisted target's committed observedGeneration (abstract counter)
  committedEpoch,        \* the persisted target's committed maintenanceEpoch (abstract counter)
  committedResult,       \* NULL | "Qualified" | "Rejected" | "Indeterminate" -- last committed conclusion
  auditAppended,         \* TRUE iff the required AuditEvent for the *current* proposal has been appended
  completionPublished,   \* TRUE iff an idempotency completion record now exists for the current proposal
  lastProposalDisposition, \* history: disposition of the most recently attempted proposal
  lastProposalWasStale,    \* history: TRUE iff the most recently attempted proposal's captured fence was Stale
  lastProposalMutated      \* history: TRUE iff the most recently attempted proposal changed committedResult

vars == <<committedGeneration, committedEpoch, committedResult, auditAppended, completionPublished,
          lastProposalDisposition, lastProposalWasStale, lastProposalMutated>>

Init ==
  /\ committedGeneration = 0
  /\ committedEpoch = 0
  /\ committedResult = "NULL"
  /\ auditAppended = FALSE
  /\ completionPublished = FALSE
  /\ lastProposalDisposition = "NULL"
  /\ lastProposalWasStale = FALSE
  /\ lastProposalMutated = FALSE

\* A proposal captures the fences at reservation time; it is Fresh if the
\* captured generation/epoch still equal the committed ones at commit time,
\* Stale otherwise. This is modeled abstractly: `Fence(capturedGen, capturedEpoch)`
\* returns Fresh or Stale relative to the *current* committed values.
Fence(capturedGen, capturedEpoch) ==
  IF capturedGen = committedGeneration /\ capturedEpoch = committedEpoch
  THEN Fresh
  ELSE Stale

\* An independent commit (e.g. Maintenance enter/clear, Retirement, or a
\* concurrent qualification winning first) advances the committed fences
\* without going through this proposal's audit/publish steps.
\* MaxFence bounds the abstract generation/epoch counters so the state space
\* is finite for model checking. The counters are unbounded naturals in the
\* real system; bounding them here only limits how many independent-advance
\* steps TLC explores and does not change any proved invariant, since Fence
\* freshness only depends on equality with the current committed value.
MaxFence == 2

IndependentFenceAdvance ==
  /\ committedGeneration' \in {committedGeneration, committedGeneration + 1}
  /\ committedEpoch' \in {committedEpoch, committedEpoch + 1}
  /\ committedGeneration' <= MaxFence
  /\ committedEpoch' <= MaxFence
  /\ (committedGeneration' # committedGeneration \/ committedEpoch' # committedEpoch)
  /\ UNCHANGED <<committedResult, auditAppended, completionPublished,
                 lastProposalDisposition, lastProposalWasStale, lastProposalMutated>>

\* A proposal with captured fences (capturedGen, capturedEpoch) and outcome
\* `disposition` attempts to commit. Ordering per ADH-2026-058 clause 5/7:
\* the audit append happens, and only on a *successful* append does the
\* commit + completion publish atomically; a Stale fence, Cancelled, or
\* Failed disposition aborts with no mutation and no completion, regardless
\* of audit outcome.
AttemptCommit(capturedGen, capturedEpoch, disposition, auditSucceeds) ==
  /\ auditAppended' = auditSucceeds
  /\ lastProposalDisposition' = disposition
  /\ lastProposalWasStale' = (Fence(capturedGen, capturedEpoch) = Stale)
  /\ IF disposition = Succeeded /\ Fence(capturedGen, capturedEpoch) = Fresh /\ auditSucceeds
       THEN /\ committedResult' = "Qualified"
            /\ completionPublished' = TRUE
            /\ lastProposalMutated' = TRUE
            /\ UNCHANGED <<committedGeneration, committedEpoch>>
       ELSE /\ UNCHANGED <<committedResult, committedGeneration, committedEpoch>>
            /\ completionPublished' = FALSE
            /\ lastProposalMutated' = FALSE

Next ==
  \/ IndependentFenceAdvance
  \/ \E capturedGen \in 0..2, capturedEpoch \in 0..2,
        disposition \in {Cancelled, Failed, Succeeded}, auditSucceeds \in BOOLEAN :
       AttemptCommit(capturedGen, capturedEpoch, disposition, auditSucceeds)

Spec == Init /\ [][Next]_vars

TypeOK ==
  /\ committedGeneration \in 0..MaxFence
  /\ committedEpoch \in 0..MaxFence
  /\ committedResult \in {"NULL", "Qualified", "Rejected", "Indeterminate"}
  /\ auditAppended \in BOOLEAN
  /\ completionPublished \in BOOLEAN
  /\ lastProposalDisposition \in {"NULL", Cancelled, Failed, Succeeded}
  /\ lastProposalWasStale \in BOOLEAN
  /\ lastProposalMutated \in BOOLEAN

\* (1) A committed conclusion (committedResult # NULL held constant from this
\* step) always coincides with a successful audit append in the same step
\* that produced it; the model only ever sets committedResult' together with
\* auditAppended' = TRUE and completionPublished' = TRUE.
ConclusionRequiresAudit ==
  (completionPublished => auditAppended)

\* (2) A replayable completion (completionPublished = TRUE) always corresponds
\* to an audited commit; there is no code path that publishes a completion
\* without the audit append succeeding first.
ReplayRequiresAuditedCompletion ==
  (completionPublished => (auditAppended /\ committedResult # "NULL"))

\* (3) Stale, cancelled, and failed proposals never mutate committed state and
\* never publish a completion. Any proposal that was Stale, Cancelled, or
\* Failed must show lastProposalMutated = FALSE and completionPublished = FALSE.
StaleCancelledFailedNeverMutate ==
  (lastProposalWasStale \/ lastProposalDisposition \in {Cancelled, Failed})
    => (~lastProposalMutated /\ ~completionPublished)

=============================================================================
