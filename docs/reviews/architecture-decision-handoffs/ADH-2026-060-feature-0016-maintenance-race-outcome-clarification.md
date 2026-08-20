# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-060
- Date: 2026-08-20
- Source discussion: Codex architecture session
- Related feature: FEATURE-0016
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Make FEATURE-0016 Maintenance-race outcomes mutually exclusive

## Summary

This correction resolves an internal FEATURE-0016 authority conflict without
adding a resource, route, field, state, top-level Problem code, dependency, or
future-phase behavior. A qualifying request whose in-flight work loses to a
committed active Maintenance entry returns `409 CONFLICT` with
`VS0_TARGET_MAINTENANCE`. A stale maintenance-epoch fence returns `412
STALE_RESOURCE_VERSION` with `VS0_TARGET_EPOCH_STALE` only when no active
current-Maintenance marker exists at commit. This preserves the established
Maintenance-wins behavior while making the stale-fence case distinct.

## Classification

Correction.

## Existing approved baseline

- ADH-2026-058 establishes a target-bound, lifecycle-service-owned,
  in-process current-Maintenance marker and states that Maintenance entry
  winning during in-flight qualification returns `409 CONFLICT` with
  `VS0_TARGET_MAINTENANCE`.
- `VS0-STATE-004.additionalInvalid` gives the same `409` outcome for
  Maintenance entry winning during in-flight qualification.
- `VS0-CF-F16-89` correctly proves that winning-Maintenance case.
- `VS0-CF-F16-75` incorrectly labels that same situation as an epoch-stale
  `412` outcome, creating incompatible observable semantics.

## Decision

At qualification commit, the `ExecutionTargetLifecycleService` evaluates the
current-Maintenance marker before comparing the captured maintenance epoch.

1. If the target has an active current-Maintenance marker, the qualifying
   caller receives `409 CONFLICT` with `VS0_TARGET_MAINTENANCE`. This applies
   when Maintenance entry won after the qualification reservation captured its
   fences. The qualification reservation aborts without a qualification
   conclusion, qualification AuditEvent, or idempotency completion; the
   independently winning Maintenance transition remains audited.
2. Otherwise, if the captured maintenance epoch differs from the current
   maintenance epoch, the caller receives `412 STALE_RESOURCE_VERSION` with
   `VS0_TARGET_EPOCH_STALE`. The qualification reservation aborts without
   target mutation, qualification AuditEvent, or idempotency completion.
3. A changed referenced InfrastructureStack generation remains the existing
   `412 STALE_RESOURCE_VERSION` with `VS0_TARGET_EPOCH_STALE` case.

The predicates are mutually exclusive. An active current-Maintenance marker
always selects the `409` result before the maintenance-epoch stale comparison.

## Rationale

Maintenance is a current lifecycle condition and must be visible as such to
the qualifying caller when it wins the race. A stale epoch without an active
marker is a distinct optimistic-concurrency failure. This ordering gives one
deterministic outcome per request and leaves no semantic choice for
requirements, design, tasks, or implementation.

## Reuse-before-build assessment

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse FEATURE-0012 `CONFLICT` and `STALE_RESOURCE_VERSION` Problem
    envelopes unchanged.
  - Reuse FEATURE-0013 AuditEvent append-before-publication semantics
    unchanged.
- Sovrunn-owned responsibility summary:
  - Clarify F0016-local selection between existing top-level Problem outcomes
    and F0016-local violation codes for one lifecycle race.
- Non-goals summary:
  - No real adapter, credential, external call, durable persistence,
    controller, plugin execution, placement, provisioning, or ExecutionTarget
    subtype is introduced.

## Phase impact

- Current phase allowed? Yes.
- Current phase boundary impact:
  - This is a deterministic, in-process contract correction within Phase 2R.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting decisions, if any:
  - None. The correction reconciles an accidental conflict between local
    F0016 conformance cases and the already-approved F0016 state rule.
- Resolution required:
  - Approved ADH application only; no DEC, RFC, or baseline update.

## Required action

- Update architecture doc
- Update Kiro requirements.md
- Update traceability matrix
- Update feature gate/checks

The exact F16-75 conformance-ledger row in requirements must be mechanically
regenerated from the corrected registry. Design, tasks, Go source, and formal
models are not modified by this handoff. After it is applied, the Design stage
must revise its Maintenance-race row to consume this authority.

## Strict Kiro read/write boundary

Kiro may read only this handoff and the following files while applying it.
It may write only files marked **write**. It must not search outside this
list.

| Path | Access | Purpose |
|---|---|---|
| `docs/reviews/architecture-decision-handoffs/ADH-2026-060-feature-0016-maintenance-race-outcome-clarification.md` | read | Controlling approved decision |
| `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md` | write | Correct the qualification-commit race rule and controlling-handoff list |
| `docs/reviews/architecture-readiness/FEATURE-0016-architecture.md` | write | Correct the canonical decision register/proof disposition |
| `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` | write | Make F16-75 the inactive-marker epoch-stale case; retain F16-89 as Maintenance-wins 409 |
| `docs/architecture/vertical-slices/VS-000-contract-specification.md` | write | State the mutually exclusive predicate/order |
| `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md` | write | Correct F16-75 description and ADH provenance |
| `.kiro/specs/adapter-boundary-and-executiontarget-qualification/requirements.md` | write | Mechanically copy the corrected F16-75 registry fields into the exact conformance ledger and update only affected supplemental traceability |
| `.kiro/steering/slice0-contract.md` | write | Record the exact race precedence rule |
| `docs/reviews/architecture-readiness/FEATURE-0016-architecture-closure-matrix.md` | write | Add resolved closure row |
| `.automation/features/FEATURE-0016.control.json` | write | Add ADH-2026-060 only to `feature.handoffs` |
| `scripts/feature-0016-architecture-readiness-check.py` | write | Fail closed on overlapping or reversed Maintenance/epoch predicates |
| `scripts/feature-contract-check.py` | write | Verify exact F16-75/F16-89 differentiation |
| `scripts/vs000-contract-check.py` | write | Verify exact registry semantics |

## Exact application changes

1. Replace the F16-75 input with: `Captured maintenance epoch differs at
   qualification commit and no active current-Maintenance marker exists`.
   Preserve its existing `STALE_RESOURCE_VERSION`,
   `VS0_TARGET_EPOCH_STALE`, 412, and non-publication semantics.
2. Retain F16-89 exactly as the Maintenance-entry-wins, active-marker,
   `CONFLICT`/`VS0_TARGET_MAINTENANCE`/409 case.
3. State the decision's ordered predicate in every active F0016 architecture
   authority: active marker first; inactive marker plus changed captured epoch
   second; changed InfrastructureStack generation remains the independent
   epoch-stale case.
4. Add `ADH-2026-060` to controlling handoffs and traceability where F16-75
   and F16-89 are named.
5. Add deterministic fail-closed checks proving that F16-75 does not say
   Maintenance entry wins and F16-89 does, and that their error/violation pairs
   remain respectively `412/VS0_TARGET_EPOCH_STALE` and
   `409/VS0_TARGET_MAINTENANCE`.
6. Regenerate only the F16-75 requirements exact-ledger row from the corrected
   registry. Preserve all canonical REQ and AC ledger rows verbatim and retain
   exactly one final `STAGE_STATUS: COMPLETE` receipt.

## Acceptance criteria for Kiro update

- [ ] F16-75 and F16-89 have mutually exclusive inputs and exact registered outcomes.
- [ ] The active-marker-first ordering appears consistently in all listed authorities.
- [ ] No top-level Problem code, field, route, state, dependency, or phase scope changes.
- [ ] Requirements F16-75 exact-ledger evidence matches the corrected registry;
  canonical REQ/AC ledgers remain untouched.
- [ ] Design, tasks, Go source, and formal models remain untouched.
- [ ] Feature readiness, feature-contract, VS-000 contract, Phase 2R drift,
  diff, documentation, and Structurizr checks pass.

## Explicit instructions to Kiro

- Do not introduce a new decision beyond this correction.
- Do not modify any file outside the strict allowlist.
- Do not update `CURRENT_ARCHITECTURE_BASELINE.md`.
- Modify only the allowlisted requirements artifact as the mechanical exact-ledger
  consequence of this correction; do not modify design, tasks, Go source, or
  formal models.
- If an authority cannot express the mutually exclusive predicate without a
  new observable contract, stop with `ARCHITECTURE_DECISION_REQUIRED`.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-20
- Notes: Founder approved the bounded F16-75/F16-89 Maintenance-race
  correction and its strict architecture-only application boundary.
