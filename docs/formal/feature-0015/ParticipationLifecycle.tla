------------------------------ MODULE ParticipationLifecycle ------------------------------
EXTENDS Naturals, Sequences, TLC

\* Bounded model of FEATURE-0015 CloudProviderParticipation lifecycle and
\* independent suspension holds.  The model deliberately omits API mechanics;
\* it proves state-machine invariants before design/code chooses data structures.

CONSTANTS Pending, Active, Suspended, Rejected, Withdrawn, Expired, Terminating, Terminated

Phases == {Pending, Active, Suspended, Rejected, Withdrawn, Expired, Terminating, Terminated}
Terminal == {Rejected, Withdrawn, Expired, Terminated}

VARIABLES phase, platformHold, providerHold
vars == <<phase, platformHold, providerHold>>

Init == /\ phase = Pending
        /\ platformHold = FALSE
        /\ providerHold = FALSE

Accept == /\ phase = Pending
          /\ phase' = Active
          /\ UNCHANGED <<platformHold, providerHold>>

Reject == /\ phase = Pending
          /\ phase' = Rejected
          /\ UNCHANGED <<platformHold, providerHold>>

Withdraw == /\ phase = Pending
            /\ phase' = Withdrawn
            /\ UNCHANGED <<platformHold, providerHold>>

Expire == /\ phase = Pending
          /\ phase' = Expired
          /\ UNCHANGED <<platformHold, providerHold>>

RequestRelease == /\ phase \in {Active, Suspended}
                  /\ phase' = Terminating
                  /\ UNCHANGED <<platformHold, providerHold>>

AcceptRelease == /\ phase = Terminating
                 /\ phase' = Terminated
                 /\ UNCHANGED <<platformHold, providerHold>>

DeclineRelease == /\ phase = Terminating
                  /\ phase' = IF platformHold \/ providerHold THEN Suspended ELSE Active
                  /\ UNCHANGED <<platformHold, providerHold>>

SetPlatformHold == /\ phase \in {Active, Suspended}
                   /\ platformHold' = TRUE
                   /\ phase' = Suspended
                   /\ UNCHANGED providerHold

ClearPlatformHold == /\ phase = Suspended /\ platformHold
                     /\ platformHold' = FALSE
                     /\ phase' = IF providerHold THEN Suspended ELSE Active
                     /\ UNCHANGED providerHold

SetProviderHold == /\ phase \in {Active, Suspended}
                   /\ providerHold' = TRUE
                   /\ phase' = Suspended
                   /\ UNCHANGED platformHold

ClearProviderHold == /\ phase = Suspended /\ providerHold
                     /\ providerHold' = FALSE
                     /\ phase' = IF platformHold THEN Suspended ELSE Active
                     /\ UNCHANGED platformHold

Next == Accept \/ Reject \/ Withdraw \/ Expire \/ RequestRelease \/ AcceptRelease \/ DeclineRelease
        \/ SetPlatformHold \/ ClearPlatformHold \/ SetProviderHold \/ ClearProviderHold

TypeOK == /\ phase \in Phases
          /\ platformHold \in BOOLEAN
          /\ providerHold \in BOOLEAN

HoldDerivedState ==
  /\ phase = Active => /\ ~platformHold /\ ~providerHold
  /\ phase = Suspended => platformHold \/ providerHold

TerminalIsFinal == phase \in Terminal => ~ENABLED Next

Spec == Init /\ [][Next]_vars

=============================================================================
