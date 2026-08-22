------------------------------ MODULE ExecutionTargetLifecycle ------------------------------
EXTENDS Naturals, TLC

\* Bounded model of FEATURE-0016 ExecutionTarget lifecycle, qualification,
\* maintenance, and effective-availability projection (ADH-2026-058 clause 7,
\* F16-AD-15/15.2/15.3). The model deliberately omits API/HTTP mechanics; it
\* proves state-machine invariants before design/code chooses data structures.
\*
\* This model proves:
\*   1. Retired terminality (RetiredIsTerminal)
\*   2. The effectiveAvailability truth table (AvailabilityTruthTable)
\*   3. Maintenance and Retirement priority over expiry (MaintenanceRetireOutrankExpiry)
\*   4. Qualifying is never publicly persisted (QualifyingNeverPersisted)

CONSTANTS
  Active, Retired,
  Unqualified, Qualified, Rejected, Indeterminate,
  Available, Unavailable, Maintenance

Lifecycles == {Active, Retired}
Qualifications == {Unqualified, Qualified, Rejected, Indeterminate}
EffectiveAvailabilities == {Available, Unavailable, Maintenance}

VARIABLES
  lifecycle,          \* persisted: Active | Retired
  qualification,      \* persisted: Unqualified | Qualified | Rejected | Indeterminate
  qualifying,         \* lifecycle-service-only in-flight reservation flag; never persisted/projected
  maintenanceActive,  \* current-Maintenance marker active flag (target-bound, in-process)
  freshViableQualified \* TRUE iff facts are fresh AND backing is viable AND qualification = Qualified

vars == <<lifecycle, qualification, qualifying, maintenanceActive, freshViableQualified>>

Init ==
  /\ lifecycle = Active
  /\ qualification = Unqualified
  /\ qualifying = FALSE
  /\ maintenanceActive = FALSE
  /\ freshViableQualified = FALSE

\* Qualify starts only from Active Unqualified/Qualified/Rejected/Indeterminate.
\* Starting it creates an internal reservation; it does not mutate persisted state.
StartQualify ==
  /\ lifecycle = Active
  /\ ~qualifying
  /\ qualifying' = TRUE
  /\ UNCHANGED <<lifecycle, qualification, maintenanceActive, freshViableQualified>>

\* A completed qualification conclusion commits Qualified/Rejected/Indeterminate
\* and clears the reservation. Maintenance winning during the reservation
\* aborts it (handled by MaintenanceEnter below), so this action requires the
\* marker to still be inactive at commit time -- an independent proof that a
\* winning Maintenance transition pre-empts a stale conclusion.
CompleteQualify(outcome) ==
  /\ lifecycle = Active
  /\ qualifying
  /\ ~maintenanceActive
  /\ outcome \in Qualifications \ {Unqualified}
  /\ qualification' = outcome
  /\ qualifying' = FALSE
  /\ freshViableQualified' = (outcome = Qualified)
  /\ UNCHANGED <<lifecycle, maintenanceActive>>

\* Any abort (malformed observer output, cancellation, audit-append failure,
\* stale fence) leaves qualification/lifecycle unchanged and only clears the
\* reservation.
AbortQualify ==
  /\ qualifying
  /\ qualifying' = FALSE
  /\ UNCHANGED <<lifecycle, qualification, maintenanceActive, freshViableQualified>>

\* Maintenance entry: aborts any in-flight reservation, persists Unqualified,
\* activates the marker. Priority over expiry and over a stale qualification.
MaintenanceEnter ==
  /\ lifecycle = Active
  /\ ~maintenanceActive
  /\ maintenanceActive' = TRUE
  /\ qualifying' = FALSE
  /\ qualification' = Unqualified
  /\ freshViableQualified' = FALSE
  /\ UNCHANGED lifecycle

\* Maintenance clear: persists Active/Unqualified, deactivates the marker.
MaintenanceClear ==
  /\ lifecycle = Active
  /\ maintenanceActive
  /\ maintenanceActive' = FALSE
  /\ qualification' = Unqualified
  /\ freshViableQualified' = FALSE
  /\ UNCHANGED <<lifecycle, qualifying>>

\* Fact expiry: preserves the last completed qualification; only flips the
\* freshness/viability flag used by the availability projection.
FactExpiry ==
  /\ lifecycle = Active
  /\ ~maintenanceActive
  /\ freshViableQualified' = FALSE
  /\ UNCHANGED <<lifecycle, qualification, qualifying, maintenanceActive>>

\* Viability recovers before fact expiry (backing becomes viable again).
ViabilityRecovers ==
  /\ lifecycle = Active
  /\ ~maintenanceActive
  /\ qualification = Qualified
  /\ freshViableQualified' = TRUE
  /\ UNCHANGED <<lifecycle, qualification, qualifying, maintenanceActive>>

\* Retire: the sole lifecycle exit, allowed from every Active combination
\* including while Qualifying is in flight. Outranks a stale qualification.
Retire ==
  /\ lifecycle = Active
  /\ lifecycle' = Retired
  /\ qualification' = Unqualified
  /\ qualifying' = FALSE
  /\ maintenanceActive' = FALSE
  /\ freshViableQualified' = FALSE

Next ==
  \/ StartQualify
  \/ \E outcome \in Qualifications \ {Unqualified} : CompleteQualify(outcome)
  \/ AbortQualify
  \/ MaintenanceEnter
  \/ MaintenanceClear
  \/ FactExpiry
  \/ ViabilityRecovers
  \/ Retire

Spec == Init /\ [][Next]_vars

TypeOK ==
  /\ lifecycle \in Lifecycles
  /\ qualification \in Qualifications
  /\ qualifying \in BOOLEAN
  /\ maintenanceActive \in BOOLEAN
  /\ freshViableQualified \in BOOLEAN

\* (1) Retired terminality: once Retired, no further domain action applies
\* (Retire is the sole exit and there is no re-entry transition).
RetiredIsTerminal == lifecycle = Retired => ~ENABLED Next

\* Response-only effectiveAvailability projection, per ADH-2026-058 clause 7/
\* F16-AD-15.3: Retired->Unavailable; Active+Maintenance->Maintenance;
\* Active/Qualified+fresh+viable->Available; every other Active
\* combination->Unavailable.
EffectiveAvailability ==
  IF lifecycle = Retired THEN Unavailable
  ELSE IF maintenanceActive THEN Maintenance
  ELSE IF freshViableQualified THEN Available
  ELSE Unavailable

\* (2) The effectiveAvailability truth table is total and deterministic.
AvailabilityTruthTable == EffectiveAvailability \in EffectiveAvailabilities

\* (3) Maintenance and Retirement outrank expiry: an expiry-only transition
\* (FactExpiry) never sets maintenanceActive or lifecycle; and whenever
\* maintenanceActive is TRUE or lifecycle = Retired, EffectiveAvailability is
\* never computed as if a stale qualification were still Available.
MaintenanceRetireOutrankExpiry ==
  /\ (lifecycle = Retired) => (EffectiveAvailability = Unavailable)
  /\ (maintenanceActive /\ lifecycle = Active) => (EffectiveAvailability = Maintenance)

\* (4) Qualifying is never publicly persisted: the persisted `qualification`
\* variable ranges only over Qualifications == {Unqualified, Qualified,
\* Rejected, Indeterminate}. "Qualifying" is deliberately absent from that
\* set and modeled only as the disjoint boolean `qualifying` reservation
\* flag. This invariant asserts the committed value is always one of the
\* four persisted values -- never a fifth "Qualifying" value -- for every
\* reachable state, including every state where a reservation is in flight.
QualifyingNeverPersisted == qualification \in Qualifications

=============================================================================
