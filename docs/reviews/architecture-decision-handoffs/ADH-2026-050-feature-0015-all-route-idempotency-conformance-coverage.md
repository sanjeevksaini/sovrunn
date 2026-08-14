# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-050
- Date: 2026-08-12
- Source discussion: FEATURE-0015 requirements review proof-coverage correction
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Correct FEATURE-0015 all-route idempotency conformance coverage

## Summary

This is a documentation/registry conformance-coverage correction, not a new product or architecture decision. The approved FEATURE-0015 authority already states that every collection-create route requires `Idempotency-Key`, and REQ-F15-18 already applies replay and changed-digest conflict semantics to every FEATURE-0015 `POST` create/action route. The F15-18/F15-19 conformance inputs presently name participation create/actions too narrowly. They must prove the existing rule across the seven collection creates and eight participation actions.

## Classification

Correction.

## Existing approved baseline

FEATURE-0015 `REQ-F15-18` states: `POST` create/actions require `Idempotency-Key`; the same principal, route, key, and digest return the original result, while a changed digest returns `CONFLICT` / `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`.

The active FEATURE-0015 authority further states that all seven collection-create routes require `Idempotency-Key`; participation create/action details remain subject to ADH-2026-046 and ADH-2026-048. ADH-2026-048 decision 4 requires complete per-route input-to-conformance coverage for all 22 routes and fail-closed treatment of a missing mapping.

`VS0-CF-F15-18` currently limits its input wording to a participation create or action, and `VS0-CF-F15-19` does not explicitly state the all-route scope. That leaves the six non-participation collection-create routes without exact local replay/mismatch evidence despite the approved behavior.

## Correction

Broaden the existing FEATURE-0015-local conformance cases; do not create new conformance IDs:

| Existing case | Corrected exact coverage |
|---|---|
| `VS0-CF-F15-18` | A same principal, route, `Idempotency-Key`, and canonical request digest replayed against any of the seven FEATURE-0015 collection creates or any of the eight existing participation create/action routes returns the original result before version checking, with no duplicate resource, transition, or AuditEvent. |
| `VS0-CF-F15-19` | Reusing the same principal, route, and `Idempotency-Key` with a changed canonical request digest against any of those same fifteen client-mutating routes retains the original result and returns `CONFLICT` with `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`. |

The seven collection creates are `CloudPlatform`, `CloudProvider`, `CloudProviderParticipation`, `HostingLocation`, `Datacenter`, `FaultDomain`, and `InfrastructureStack`. The eight participation actions are `accept`, `reject`, `withdraw`, `suspend`, `resume`, `request-release`, `accept-release`, and `decline-release`.

No rule is added for GET, LIST, PATCH, scheduler expiry, or a future-feature route. Existing authorization, safe-reference, structural-classification, replay ordering, version-check, audit, and non-publication rules remain unchanged.

## Rationale

The correction aligns the registry evidence with the approved, already-observable all-create idempotency contract and with ADH-2026-048's required route-by-route proof. It avoids narrowing REQ-F15-18 or silently treating participation-only evidence as evidence for unrelated resource creates.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Reuse
- Summary of mature candidates / applicable standards:
  - Reuse FEATURE-0012 Problem Details `CONFLICT` and the existing Slice 0 `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH` violation.
  - Reuse the established F0015 idempotency semantics and existing F15-18/F15-19 conformance IDs.
- Sovrunn-owned responsibility summary:
  - Make the existing local proof scope exact for every approved F0015 client-mutating POST route.
- Non-goals summary:
  - No new resource, field, route, method, header rule, lifecycle state, writer, grant, Problem code, violation code, persistence model, feature ownership, Go implementation, or generated Kiro artifact.

## Phase impact

- Current phase allowed: Yes.
- Current phase boundary impact: F0015 evidence correction only; no Phase 2R scope expansion.

## Conflict check

- Conflicts with accepted DEC/RFC: No.
- Resolution required: None. This correction implements already-approved F0015 semantics.

## Required action

- Update the F0015 feature and architecture authorities, VS-000 registry/specification, traceability, closure matrix, steering contract, and exact validation checks only as necessary to prove the corrected F15-18/F15-19 scope.
- Update `.automation/features/FEATURE-0015.control.json` by retaining ADH-2026-048 and ADH-2026-049 and adding ADH-2026-050 to `feature.handoffs`; do not modify any other manifest field.
- Regenerate or revise requirements only after the architecture validations pass; do not modify requirements during this architecture handoff.

## Strict Kiro read/write boundary

Kiro must read only these files:

- this handoff;
- `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`;
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`;
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`;
- `docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md`;
- `scripts/feature-0015-architecture-readiness-check.py`;
- `scripts/vs000-contract-check.py`;
- `.kiro/steering/slice0-contract.md`;
- `.automation/features/FEATURE-0015.control.json`.

Kiro may write only the same files, plus this handoff for a status annotation. It must not search, glob, inspect, or follow references elsewhere in the repository. If another file appears necessary, it must stop and report its exact path and reason without opening it.

Kiro must not modify Go source, `.kiro/specs/**/requirements.md`, `design.md`, `tasks.md`, any DEC, baseline, roadmap, feature sequence, or state file.

## Impacted features

- FEATURE-0015: exact all-route idempotency replay and mismatch evidence.
- FEATURE-0016 through FEATURE-0026: no behavior change.

## Acceptance criteria for Kiro update

- [ ] `VS0-CF-F15-18` and `VS0-CF-F15-19` each explicitly cover all seven collection creates and all eight participation create/actions.
- [ ] The cases preserve their existing IDs, expected states, Problem code, and violation semantics.
- [ ] All F0015 resource/mapping authorities identify F15-18/F15-19 as applicable where required by the all-route contract.
- [ ] F0015 readiness and registry checks fail closed if any of the fifteen client-mutating routes lacks replay/mismatch conformance coverage.
- [ ] ADH-2026-048 and ADH-2026-049 remain in `feature.handoffs`, and ADH-2026-050 is added with no other manifest change.
- [ ] No new behavior, conformance ID, or out-of-scope route is introduced.
- [ ] Only files in the strict allowlist changed.
- [ ] `make feature-0015-architecture-readiness`, `make vs000-contract-check`, `make phase2r-drift-check`, `git diff --check`, `mkdocs build --strict`, and `make structurizr-check` pass.

## Explicit instructions to Kiro

- Apply no decision beyond this correction handoff.
- Do not search the repository outside the strict allowlist.
- Do not modify generated Kiro artifacts or Go source.
- Stop, do not infer, if the correction cannot be represented using only the allowed files.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-12
- Notes: Approved correction to already-approved all-route idempotency behavior; no new architecture decision.
