--------------------------- MODULE IdempotencyAuditPublication ---------------------------
EXTENDS Naturals, Sequences, TLC

\* Bounded model of ADH-2026-048/049 ordering: a required AuditEvent must be
\* accepted before a mutation/idempotency completion is published.  It permits
\* an audit record to remain after a later in-process publication failure.

VARIABLES reservation, auditForCurrent, auditHistory, mutationPublished, idempotencyCompleted
vars == <<reservation, auditForCurrent, auditHistory, mutationPublished, idempotencyCompleted>>

Init == /\ reservation = "Absent"
        /\ auditForCurrent = FALSE
        /\ auditHistory = FALSE
        /\ mutationPublished = FALSE
        /\ idempotencyCompleted = FALSE

Reserve == /\ reservation = "Absent"
           /\ reservation' = "InFlight"
           /\ auditForCurrent' = FALSE
           /\ UNCHANGED <<auditHistory, mutationPublished, idempotencyCompleted>>

AppendAudit == /\ reservation = "InFlight"
               /\ auditForCurrent' = TRUE
               /\ auditHistory' = TRUE
               /\ UNCHANGED <<reservation, mutationPublished, idempotencyCompleted>>

Publish == /\ reservation = "InFlight" /\ auditForCurrent
           /\ mutationPublished' = TRUE
           /\ idempotencyCompleted' = TRUE
           /\ reservation' = "Completed"
           /\ UNCHANGED <<auditForCurrent, auditHistory>>

AbortBeforeAudit == /\ reservation = "InFlight" /\ ~auditForCurrent
                    /\ reservation' = "Aborted"
                    /\ UNCHANGED <<auditForCurrent, auditHistory, mutationPublished, idempotencyCompleted>>

AbortAfterAudit == /\ reservation = "InFlight" /\ auditForCurrent
                   /\ reservation' = "Aborted"
                   /\ UNCHANGED <<auditForCurrent, auditHistory, mutationPublished, idempotencyCompleted>>

ReleaseAborted == /\ reservation = "Aborted"
                  /\ reservation' = "Absent"
                  /\ UNCHANGED <<auditForCurrent, auditHistory, mutationPublished, idempotencyCompleted>>

Replay == /\ reservation = "Completed"
          /\ UNCHANGED vars

Next == Reserve \/ AppendAudit \/ Publish \/ AbortBeforeAudit \/ AbortAfterAudit \/ ReleaseAborted \/ Replay

TypeOK == /\ reservation \in {"Absent", "InFlight", "Completed", "Aborted"}
          /\ auditForCurrent \in BOOLEAN
          /\ auditHistory \in BOOLEAN
          /\ mutationPublished \in BOOLEAN
          /\ idempotencyCompleted \in BOOLEAN

PublishRequiresAudit == mutationPublished => auditForCurrent
CompletionRequiresPublish == idempotencyCompleted => mutationPublished
AbortedNeverPublishes == reservation = "Aborted" => /\ ~mutationPublished /\ ~idempotencyCompleted

Spec == Init /\ [][Next]_vars

=============================================================================
