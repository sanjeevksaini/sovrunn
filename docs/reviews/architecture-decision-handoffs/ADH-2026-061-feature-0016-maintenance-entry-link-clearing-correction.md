# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-061
- Date: 2026-08-20
- Source discussion: Codex architecture session
- Related feature: FEATURE-0016
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Correct FEATURE-0016 Maintenance-entry link clearing and qualification-abort scope

## Summary

This correction removes one contradictory sentence from the FEATURE-0016
feature authority. The active VS-000 registry already requires every
successful Maintenance entry to clear the target's current
`NormalizedTargetFactSet` and `TargetQualificationResult` links and persist
`Active/Unqualified`. The feature authority incorrectly preserves the last
completed conclusion when no qualification is in flight.

No resource, route, field, state, top-level Problem code, dependency, or
future-phase behavior is added. This handoff makes the feature authority match
the already-approved registry and conformance cases.

## Classification

Correction.

## Existing approved baseline

- `VS0-STATE-004` requires Maintenance entry from every Active qualification
  state to persist `Active/Unqualified`, clear current FactSet/Result links,
  and advance the maintenance epoch and resource version.
- `VS0-CF-F16-42` proves that a successful current-fenced Maintenance entry
  clears both links and leaves an active current-Maintenance marker.
- `VS0-CF-F16-43` proves that Maintenance clear remains
  `Active/Unqualified` with both links cleared.
- `VS0-CF-F16-89` proves that Maintenance winning during an in-flight
  qualification returns `409 CONFLICT` with `VS0_TARGET_MAINTENANCE`, without
  a qualification completion or qualification AuditEvent.
- ADH-2026-060 establishes that the active current-Maintenance marker selects
  the F16-89 result before maintenance-epoch stale comparison.
- The F0016 feature authority currently says that Maintenance entry with no
  in-flight qualification preserves the last completed conclusion. That is
  incompatible with the state contract and F16-42.

## Decision

The `ExecutionTargetLifecycleService` applies the following rule to **every
successful Maintenance entry**:

1. It establishes the current-Maintenance marker, advances the maintenance
   epoch and resource version/ETag, clears both current
   `NormalizedTargetFactSet` and `TargetQualificationResult` links, and
   persists lifecycle `Active` with qualification `Unqualified`.
2. If a qualification is in flight, it additionally aborts **only that
   qualification** and its linked idempotency reservation, and wakes only the
   waiters for that reservation. The qualifying caller receives the existing
   F16-89 `409 CONFLICT` plus `VS0_TARGET_MAINTENANCE` outcome. No
   qualification result, qualification completion, or qualification AuditEvent
   is published; the independently successful Maintenance entry is audited.
3. If no qualification is in flight, no idempotency reservation is cancelled.
   The same link-clearing and `Active/Unqualified` transition still occurs and
   the successful Maintenance entry is audited.

Maintenance clear remains `Active/Unqualified` with both current links
cleared, as already specified by F16-43.

## Rationale

Maintenance makes prior qualification evidence non-current. Clearing both
links consistently prevents a stale qualified, rejected, or indeterminate
conclusion from being retained behind a Maintenance projection. Narrowing the
abort scope prevents one target's Maintenance operation from cancelling an
unrelated target's idempotency work.

## Example

An `Active/Qualified` target has links to FactSet `F7` and Result `R7`.

- A successful Maintenance trigger commits `Active/Unqualified`, clears both
  links, advances its epoch and ETag, records the active marker, and appends
  the Maintenance AuditEvent. It does **not** retain `R7`.
- If `/actions/qualify` had reserved work on that same target, its linked
  reservation aborts and its waiters wake. The qualifying request gets
  `409 CONFLICT` with `VS0_TARGET_MAINTENANCE`; it does not create a new
  conclusion or qualification AuditEvent.
- A simultaneous qualification reservation for a different target is not
  affected.

## Reuse-before-build assessment

- Disposition: Clarify existing reuse.
- FEATURE-0012 top-level Problem envelopes are reused unchanged.
- FEATURE-0013 AuditEvent append-before-publication semantics are reused
  unchanged.
- No real adapter, credential, external call, durable persistence, controller,
  plugin execution, placement, provisioning, or ExecutionTarget subtype is
  introduced.

## Phase impact

- Current phase allowed? Yes.
- Current phase boundary impact: none; this is a deterministic in-process
  correction to an existing Phase 2R lifecycle contract.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting authority: the F0016 feature authority's retained-conclusion
  sentence conflicts with `VS0-STATE-004`, F16-42, and F16-43.
- Resolution required: founder approval followed by this correction's strict
  application. No DEC, RFC, or architecture-baseline update is required.

## Required action

- Update the F0016 feature authority, architecture readiness record, contract
  prose, traceability, steering, closure evidence, and control manifest.
- Add fail-closed checks for unconditional Maintenance-entry link clearing and
  target-specific qualification/idempotency abort scope.
- Do not modify requirements, design, tasks, Go source, or formal models as
  part of this handoff. A later tasks-only revision consumes the corrected
  authority.

## Strict Kiro read/write boundary

Kiro may read only this handoff and the following files while applying it. It
may write only files marked **write**. It must not search outside this list.

| Path | Access | Purpose |
|---|---|---|
| `docs/reviews/architecture-decision-handoffs/ADH-2026-061-feature-0016-maintenance-entry-link-clearing-correction.md` | read | Controlling approved decision |
| `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` | read | Existing canonical state and F16-42/43/89 evidence |
| `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md` | write | Replace the retained-conclusion conflict and add controlling handoff provenance |
| `docs/reviews/architecture-readiness/FEATURE-0016-architecture.md` | write | Align F16 decision register and lifecycle explanation |
| `docs/architecture/vertical-slices/VS-000-contract-specification.md` | write | State the unconditional link-clearing and target-specific abort rule |
| `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md` | write | Add ADH-061 provenance to F16-42, F16-43, and F16-89 |
| `.kiro/steering/slice0-contract.md` | write | Record the exact Maintenance-entry invariant |
| `docs/reviews/architecture-readiness/FEATURE-0016-architecture-closure-matrix.md` | write | Add the resolved architecture-contradiction closure row |
| `.automation/features/FEATURE-0016.control.json` | write | Add ADH-2026-061 only to `feature.handoffs` |
| `scripts/feature-0016-architecture-readiness-check.py` | write | Fail closed on retained-conclusion wording or an over-broad abort scope |
| `scripts/feature-contract-check.py` | write | Verify F0016 contract evidence remains aligned |
| `scripts/vs000-contract-check.py` | write | Verify registry evidence used by the correction remains intact |

## Exact application changes

1. Replace the feature-authority statement that preserves the last completed
   conclusion when no qualification is in flight. State instead that every
   successful Maintenance entry clears both current FactSet/Result links and
   persists `Active/Unqualified`.
2. State that an in-flight qualification adds a target-specific abort: only
   the qualifying reservation and its linked idempotency reservation abort;
   only those waiters wake.
3. Preserve F16-42, F16-43, and F16-89 fields and outcomes exactly. The
   registry is already semantically correct and requires no product-contract
   alteration.
4. Add `ADH-2026-061` to F0016 controlling-handoff and traceability evidence.
5. Add deterministic checks that reject (a) retained completed links or a
   retained qualification conclusion on Maintenance entry, and (b) any claim
   that all target reservations are aborted.

## Acceptance criteria for Kiro update

- [ ] Every successful Maintenance entry clears both current FactSet/Result
  links and persists `Active/Unqualified`.
- [ ] F16-42, F16-43, and F16-89 retain their current exact observable
  conformance semantics.
- [ ] An in-flight qualification abort affects only its linked reservation and
  waiters; no unrelated reservation is cancelled.
- [ ] The winning Maintenance transition remains independently audited, with
  no qualification completion or qualification AuditEvent for F16-89.
- [ ] Requirements, design, tasks, Go source, and formal models are untouched.
- [ ] Feature readiness, feature-contract, VS-000 contract, Phase 2R drift,
  diff, documentation, and Structurizr checks pass.

## Explicit instructions to Kiro

- This handoff is founder-approved and may be applied only through the strict
  Kiro architecture-update workflow.
- Do not introduce a new resource, route, field, state, Problem code,
  dependency, or future-phase behavior.
- Do not modify any file outside the strict allowlist.
- Do not edit the registry's F16-42, F16-43, or F16-89 product semantics;
  they are the source being reconciled.
- Do not modify requirements, design, tasks, Go source, or formal models.
- If the listed authority cannot express the rule without a new observable
  contract, stop with `ARCHITECTURE_DECISION_REQUIRED`.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-20
- Notes: Founder approved the bounded Maintenance-entry link-clearing
  correction. It becomes active only after Kiro applies it through the strict
  architecture-update workflow.
