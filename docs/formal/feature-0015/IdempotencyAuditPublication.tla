--------------------------- MODULE IdempotencyAuditPublication ---------------------------
EXTENDS Naturals, Sequences, TLC

\* Bounded model of ADH-2026-048/049 ordering: a required AuditEvent must be
\* accepted before a mutation/idempotency completion is published.  It permits
\* an audit record to remain after a later in-process publication failure.

VARIABLES reservation, auditForCurrent, auditHistory, mutationPublished, idempotencyCompleted,
          requestTarget, reservationTarget, completedTarget
vars == <<reservation, auditForCurrent, auditHistory, mutationPublished, idempotencyCompleted,
          requestTarget, reservationTarget, completedTarget>>

Init == /\ reservation = "Absent"
        /\ auditForCurrent = FALSE
        /\ auditHistory = FALSE
        /\ mutationPublished = FALSE
        /\ idempotencyCompleted = FALSE
        /\ requestTarget = "A"
        /\ reservationTarget = "None"
        /\ completedTarget = "None"

Reserve == /\ reservation = "Absent"
           /\ reservation' = "InFlight"
           /\ auditForCurrent' = FALSE
           /\ reservationTarget' = requestTarget
           /\ UNCHANGED <<auditHistory, mutationPublished, idempotencyCompleted, requestTarget, completedTarget>>

AppendAudit == /\ reservation = "InFlight"
               /\ auditForCurrent' = TRUE
               /\ auditHistory' = TRUE
               /\ UNCHANGED <<reservation, mutationPublished, idempotencyCompleted, requestTarget, reservationTarget, completedTarget>>

Publish == /\ reservation = "InFlight" /\ auditForCurrent
           /\ mutationPublished' = TRUE
           /\ idempotencyCompleted' = TRUE
           /\ reservation' = "Completed"
           /\ completedTarget' = reservationTarget
           /\ UNCHANGED <<auditForCurrent, auditHistory, requestTarget, reservationTarget>>

AbortBeforeAudit == /\ reservation = "InFlight" /\ ~auditForCurrent
                    /\ reservation' = "Aborted"
                    /\ UNCHANGED <<auditForCurrent, auditHistory, mutationPublished, idempotencyCompleted, requestTarget, reservationTarget, completedTarget>>

AbortAfterAudit == /\ reservation = "InFlight" /\ auditForCurrent
                   /\ reservation' = "Aborted"
                   /\ UNCHANGED <<auditForCurrent, auditHistory, mutationPublished, idempotencyCompleted, requestTarget, reservationTarget, completedTarget>>

ReleaseAborted == /\ reservation = "Aborted"
                  /\ reservation' = "Absent"
                  /\ UNCHANGED <<auditForCurrent, auditHistory, mutationPublished, idempotencyCompleted, requestTarget, reservationTarget, completedTarget>>

Replay == /\ reservation = "Completed" /\ requestTarget = completedTarget
          /\ UNCHANGED vars

StartOtherTarget == /\ reservation = "Completed" /\ requestTarget = "A"
                    /\ reservation' = "Absent"
                    /\ requestTarget' = "B"
                    /\ reservationTarget' = "None"
                    /\ completedTarget' = "None"
                    /\ auditForCurrent' = FALSE
                    /\ mutationPublished' = FALSE
                    /\ idempotencyCompleted' = FALSE
                    /\ UNCHANGED <<auditHistory>>

Next == Reserve \/ AppendAudit \/ Publish \/ AbortBeforeAudit \/ AbortAfterAudit \/ ReleaseAborted \/ Replay \/ StartOtherTarget

TypeOK == /\ reservation \in {"Absent", "InFlight", "Completed", "Aborted"}
          /\ auditForCurrent \in BOOLEAN
          /\ auditHistory \in BOOLEAN
          /\ mutationPublished \in BOOLEAN
          /\ idempotencyCompleted \in BOOLEAN
          /\ requestTarget \in {"A", "B"}
          /\ reservationTarget \in {"None", "A", "B"}
          /\ completedTarget \in {"None", "A", "B"}

PublishRequiresAudit == mutationPublished => auditForCurrent
CompletionRequiresPublish == idempotencyCompleted => mutationPublished
AbortedNeverPublishes == reservation = "Aborted" => /\ ~mutationPublished /\ ~idempotencyCompleted
CompletedReplayIsTargetBound == reservation = "Completed" => completedTarget = requestTarget

Spec == Init /\ [][Next]_vars

=============================================================================
