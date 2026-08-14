# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-049
- Date: 2026-08-12
- Source discussion: FEATURE-0015 requirements review
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Define FEATURE-0015 audit-append failure Problem Details outcome

## Summary

This bounded clarification assigns the existing `INTERNAL_ERROR` / HTTP 500 Problem Details outcome to a failed required FEATURE-0013 AuditEvent append in a FEATURE-0015 mutation. The mutation remains unpublished and no idempotency result is completed. No durable or external audit dependency is introduced, so `DEPENDENCY_UNAVAILABLE` is not applicable.

## Classification

Clarification.

## Existing approved baseline

FEATURE-0012 owns the closed Problem Details code set, including `INTERNAL_ERROR` / 500 and `DEPENDENCY_UNAVAILABLE` / 503. FEATURE-0013 owns `AuditEvent`. ADH-2026-048 decision 3 requires that a FEATURE-0015 resource or idempotency change is not published until its required audit append succeeds, but it did not select the resulting client-visible 5xx Problem Details outcome. `VS0-CF-F15-24` currently uses the non-exact phrase “established standard 5xx problem.”

## Decision

For every FEATURE-0015 mutation that requires a FEATURE-0013 AuditEvent, a failure to append that required AuditEvent returns the inherited FEATURE-0012 `INTERNAL_ERROR` Problem Details response with HTTP status 500.

The resource mutation and its idempotency completion are not published. An AuditEvent that was successfully appended before a later in-process failure may remain, as already constrained by ADH-2026-048 decision 3; this clarification does not claim a durable cross-store transaction.

`DEPENDENCY_UNAVAILABLE` / 503 is not used for this outcome because FEATURE-0015 introduces no durable or external persistence dependency.

## Rationale

The failure is an internal failure of the locally required audit-append operation, not a declared unavailable external dependency. Selecting the inherited 500 mapping removes the last client-visible ambiguity without adding a code, route, resource, dependency, lifecycle, or implementation mechanism.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Reuse
- Summary of mature candidates / applicable standards:
  - Reuse FEATURE-0012 Problem Details `INTERNAL_ERROR` / 500 and FEATURE-0013 `AuditEvent` ownership.
- Sovrunn-owned responsibility summary:
  - Apply the inherited mapping to the already-approved F0015 atomic audit obligation.
- Non-goals summary:
  - No new Problem code, violation code, persistence dependency, transaction mechanism, route, resource, lifecycle state, writer, grant, Go implementation, or generated Kiro artifact.

## Phase impact

- Current phase allowed: Yes.
- Current phase boundary impact: clarification of existing FEATURE-0015 observable error behavior only.

## Conflict check

- Conflicts with accepted DEC/RFC: No.
- Resolution required: None. This handoff updates the active F0015 contract authorities and checks only.

## Required action

- Update F0015 feature and architecture authorities.
- Update the VS-000 registry, specification, and traceability for the exact `INTERNAL_ERROR` / 500 mapping.
- Update the F0015 closure matrix and architecture-readiness/contract checks to fail closed if F15-24 does not have that exact mapping.
- Regenerate requirements only after the architecture validations pass.

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

In the control manifest, retain ADH-2026-048 and add this ADH-2026-049 handoff to `feature.handoffs`; do not change any other manifest field.

Kiro must not modify Go source, `.kiro/specs/**/requirements.md`, `design.md`, `tasks.md`, any DEC, baseline, roadmap, feature sequence, control manifest, or state file.

## Impacted features

- FEATURE-0015: exact audit-append failure response.
- FEATURE-0016 through FEATURE-0026: no behavior change.

## Acceptance criteria for Kiro update

- [ ] `VS0-CF-F15-24` has `expectedError: INTERNAL_ERROR` and its exact 500 mapping is consistently stated in the allowed F0015 authorities.
- [ ] `DEPENDENCY_UNAVAILABLE` is not assigned to F0015 audit-append failure.
- [ ] The mutation/idempotency non-publication rule remains unchanged.
- [ ] The readiness and registry checks fail closed on a missing or different F15-24 mapping.
- [ ] The control manifest retains ADH-2026-048 and includes ADH-2026-049, with no other manifest changes.
- [ ] No new code, violation, resource, route, dependency, lifecycle, ownership, or design mechanism is introduced.
- [ ] Only files within the strict Kiro allowlist changed.
- [ ] `make feature-0015-architecture-readiness`, `make vs000-contract-check`, `make phase2r-drift-check`, `git diff --check`, `mkdocs build --strict`, and `make structurizr-check` pass.

## Explicit instructions to Kiro

- Apply no decision beyond this handoff.
- Do not search the repository outside the strict allowlist.
- Do not modify generated Kiro artifacts or Go source.
- Stop, do not infer, if this mapping cannot be represented using only the allowed files.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-12
- Notes: Bounded audit-failure clarification. `INTERNAL_ERROR` / 500 is approved; no architecture exploration is authorized.
