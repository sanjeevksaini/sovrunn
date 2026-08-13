# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-057
- Date: 2026-08-13
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Close FEATURE-0015 mutable-spec writer-enforcement and acceptance-proof coverage

## Classification

Clarification — the approved writer boundary exists, but it lacks one exact FEATURE-0015-local runtime proof case and the supplemental acceptance mapping overclaims unrelated evidence.

## Existing approved baseline

- REQ-F15-11 and AC-F15-08 require writer enforcement: only the CloudPlatform administrator may write CloudPlatform spec, and only the CloudProvider administrator may write CloudProvider/topology spec.
- The approved mutable PATCH surface remains unchanged: CloudPlatform `spec.description`; CloudProvider `spec.displayName` and `spec.operatingMarkets`; topology `spec.description`; CloudProviderParticipation has no PATCH surface.
- ADH-2026-056 already defines F15-38 PATCH conditional concurrency, F15-39 exact read authorization, and F15-40 replay-response behavior.
- Existing F15-05, F15-06, F15-14, and X03 prove status/system-owned/field-boundary/cross-provider behavior, not a wrong-administrator mutable-spec write denial.

## Decision

1. Add exactly one FEATURE-0015-local conformance case, `VS0-CF-F15-41`, for an authenticated principal that attempts an otherwise syntactically valid mutable-spec PATCH outside its exact writer boundary:
   - a non-CloudPlatform administrator PATCHing CloudPlatform `spec.description`;
   - a non-CloudProvider administrator PATCHing CloudProvider `spec.displayName` or `spec.operatingMarkets`;
   - a non-CloudProvider administrator PATCHing topology `spec.description`.
   The outcome is `AUTHORIZATION_DENIED`; no resource mutation, status/lifecycle update, AuditEvent, or idempotency record is published. The case must state that this authorization decision occurs before semantic PATCH processing and publication. It does not alter safe-denial rules for inaccessible cross-provider references.
2. Map `VS0-CF-F15-41` as the exact writer-enforcement proof for REQ-F15-11 and AC-F15-08. Retain existing F15-05/F15-06/F15-14/X03 mappings only for the distinct behavior each actually proves.
3. Correct the supplemental AC-F15-01 scenario and all directly related acceptance mappings without changing the verbatim canonical REQ, AC, or conformance ledgers:
   - collection creation remains F15-01/F15-02;
   - GET/LIST authorization and scoped visibility uses F15-39;
   - permitted PATCH conditional concurrency uses F15-38;
   - completed replay response semantics uses F15-40.
   A scenario must not claim evidence from a case that does not cover the described observable behavior.

## Non-goals

- No new resource, route, grant action, lifecycle state, mutable field, Problem code, violation code, audit category, or storage dependency.
- No change to PATCH `If-Match`, read-action, replay-retention, idempotency, or safe-denial semantics established by ADH-2026-056.
- No change to the canonical REQ or AC ledger wording.

## Required action

- Update the registry, contract specification, traceability, FEATURE-0015 feature and architecture authorities, closure matrix, steering, control manifest, and deterministic checks atomically.
- Update only the affected supplemental requirements and design mappings plus the exact conformance ledger after the registry change.
- Do not update tasks or Go code.

## Strict read/write boundary

- this handoff;
- `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`;
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`;
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`;
- `docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md`;
- `.kiro/steering/slice0-contract.md`;
- `.automation/features/FEATURE-0015.control.json`;
- `scripts/feature-0015-architecture-readiness-check.py`;
- `scripts/feature-contract-check.py`;
- `scripts/vs000-contract-check.py`;
- `.kiro/specs/canonical-cloud-model-and-alpha-migration/requirements.md`;
- `.kiro/specs/canonical-cloud-model-and-alpha-migration/design.md`.

## Acceptance criteria for Kiro update

- [ ] `VS0-CF-F15-41` exists exactly once, is owned by FEATURE-0015, and proves all three wrong-administrator mutable PATCH boundaries with `AUTHORIZATION_DENIED` and no publication side effects.
- [ ] REQ-F15-11 and AC-F15-08 map to F15-41 without repurposing existing proof cases.
- [ ] The AC-F15-01 supplemental scenario maps creation, reads, PATCH concurrency, and replay behavior only to their semantically matching F15-01/02, F15-38, F15-39, and F15-40 cases.
- [ ] The three canonical ledgers preserve their approved rows, except for the exact F15-41 registry-ledger addition required by this handoff.
- [ ] Readiness, feature-contract, VS-000 contract, formal, Phase 2R drift, semantic, documentation, and diff checks pass.

## Explicit instructions to Kiro

- Do not introduce decisions beyond this handoff.
- Do not modify `CURRENT_ARCHITECTURE_BASELINE.md`, a DEC, tasks, or Go code.
- Do not broaden an existing case to stand in for writer enforcement.
- Do not rewrite requirements or design; make only the narrow authority and mapping changes described above.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-13
- Notes: Approved as a bounded FEATURE-0015 conformance and acceptance-proof closure.
