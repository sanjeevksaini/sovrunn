# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-048
- Date: 2026-08-12
- Source discussion: FEATURE-0015 design-review input-contract clarification
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Clarify FEATURE-0015 ISO validation and malformed input outcomes

## Summary

This bounded clarification closes four externally observable input-contract outcomes that are already within the approved FEATURE-0015 resource and API surface. It introduces no resource, field, route, lifecycle state, writer, grant, dependency, top-level Problem code, feature ownership, or Phase 2R capability. It also requires a deterministic per-route design-readiness simulation so an unresolved observable input outcome cannot reach Kiro requirements, design, tasks, or review again.

## Classification

Clarification of existing FEATURE-0012/F0015 validation and Problem Details contracts.

## Existing approved baseline

ADH-2026-045 through ADH-2026-047 define the canonical-bootstrap F0015 resource surface, closed collection-create requests, participation actions, server-derived scope, and FEATURE-0012 Problem Details reuse. FEATURE-0012 owns the closed top-level Problem-code set. ADH-2026-046 already defines missing, malformed, or stale `If-Match` on an existing participation action as `STALE_RESOURCE_VERSION` (412).

The authorities require ISO-3166 values and require/reject specific headers and bodies, but did not state assigned-code membership or exact outcomes for several malformed/prohibited inputs.

## Decision

### 1. ISO-3166 validation

For FEATURE-0015, an ISO-3166-1 alpha-2 value means an **assigned** code from a fixed, repository-owned, version-pinned static dataset. The dataset contains the assigned ISO-3166-1 alpha-2 code list as of `2026-08-12`; it is compiled into or otherwise bundled with the repository, requires no network request, and introduces no third-party runtime dependency. A syntactically valid but unassigned value such as `ZZ` is rejected with the existing `VALIDATION_FAILED` (422) Problem Details contract.

`HostingLocation.spec.administrativeAreaCode` remains syntax-plus-country-prefix validation only: it must match the existing ISO-3166-2 shape when present and have the same two-letter prefix as `spec.countryCode`. FEATURE-0015 does not introduce an ISO-3166-2 membership dataset.

### 2. Header and item-action body outcomes

Use only existing FEATURE-0012 top-level Problem Details codes:

| Input | Exact outcome |
|---|---|
| Required `Idempotency-Key` is missing, empty, malformed, or exceeds the shared length limit | `MALFORMED_REQUEST` / 400 |
| `If-Match` is supplied on participation collection create | `MALFORMED_REQUEST` / 400 |
| `If-Match` is missing, malformed, or non-current on an existing participation action | `STALE_RESOURCE_VERSION` / 412 (preserved ADH-2026-046 rule) |
| A participation item action has any non-empty JSON body | `MALFORMED_REQUEST` / 400 |

No new top-level Problem code or violation code is introduced. A malformed/prohibited input produces no mutation, no idempotency record, and no AuditEvent unless an already-approved audited denial rule independently applies; FEATURE-0015 does not broaden the ADH-046/047 audit-denial list through this clarification.

### 3. Required deterministic design mechanics

The following are implementation/design mechanics, not new product decisions. The F0015 design must define them exactly and must not seek new architecture authority for them:

- idempotency reservation states `InFlight`, `Completed`, and `Aborted`; waiters wake on both terminal states; validation failures, stale versions, audit failures, and recovered panics abort and remove/release an in-flight reservation so a later valid request can retry;
- idempotency records are written only for completed replayable outcomes consistent with the approved audit policy; an audited denial never receives an idempotency record unless the existing F0015 requirements expressly make that denial replayable;
- in-memory audit/mutation coordination must state the achievable guarantee precisely: a resource/idempotency change is not published until its required audit append succeeds; if append fails, no resource/idempotency change is published. The design must not claim a durable cross-store transaction or hide the fact that an already-appended audit record may survive an in-process failure before resource publication;
- Go 1.22 `http.ServeMux` registration uses explicit method/path patterns and `Request.PathValue`. F0015 has seven collection registrations, seven item registrations, and eight explicit participation action registrations (22 total); it uses no wildcard, reflection, or alpha-handler auto-registration.

### 4. Final per-route design-readiness simulation

Before a requirements-stage or design-stage artifact can be accepted, the F0015 readiness checker must prove, for every owned collection, item, and action route:

```text
request body shape + required/prohibited headers
→ authentication/current authorization
→ root/reference/safe-denial ordering
→ structural and semantic validation
→ exact Problem Details outcome
→ state/status result
→ audit effect
→ idempotency/race outcome
→ F0015-local conformance test
```

The checker must fail closed on an unspecified observable input outcome, missing per-route mapping, unassigned ISO code accepted by the defined validator, non-empty participation-action body accepted, or a route registration count/pattern inconsistent with the explicit 22-route design.

### 5. Canonical ledgers remain mandatory

Requirements must retain exactly one verbatim Canonical requirement ledger, Canonical acceptance ledger, and Exact conformance semantics ledger. The ledgers must not be replaced by citations. Compact traceability may supplement them only.

## Rationale

The clarification prevents design from selecting client-visible validation and error behavior. It preserves the approved F0015 model and reuses existing FEATURE-0012 and FEATURE-0013 contracts without a new dependency or generic architecture revision.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse FEATURE-0012 Problem Details, header validation, metadata, and status semantics.
  - Reuse FEATURE-0013 AuditEvent ownership and existing F0015 audit policy.
  - Use Go standard-library HTTP routing and a repository-owned static validation dataset.
- Sovrunn-owned responsibility summary:
  - Define F0015 exact input classification and per-route conformance composition.
- Non-goals summary:
  - No new API resource, route, state, grant, writer, error code, dependency, migration behavior, ExecutionTarget behavior, FEATURE-0021 behavior, Go implementation, or generated Kiro artifact in this architecture update.

## Phase impact

- Current phase allowed: Yes.
- Current phase boundary impact: clarification of approved F0015 behavior only.

## Conflict check

- Conflicts with accepted DEC/RFC: No.
- Resolution required: bounded active-authority and guardrail clarification under this approved handoff. No DEC, RFC, baseline, or feature-sequence update is required.

## Required action

- Update F0015 feature and architecture authority
- Update VS-000 registry/specification/conformance and traceability only as necessary for the exact decisions above
- Update F0015 closure matrix and readiness/contract checks
- Regenerate requirements only after architecture validations pass

## Strict Kiro read/write boundary

Kiro must read only these files:

- this handoff
- `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`
- `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md`
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`
- `docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md`
- `scripts/feature-0015-architecture-readiness-check.py`
- `scripts/vs000-contract-check.py`
- `.kiro/steering/slice0-contract.md`

Kiro may write only the same files, plus this handoff if a status annotation is required. It must not search, glob, inspect, or follow references elsewhere in the repository. If an additional file appears necessary, it must stop and report that exact path and reason without opening it.

Kiro must not modify Go source, `.kiro/specs/**/requirements.md`, `design.md`, `tasks.md`, any DEC, baseline, roadmap, feature sequence, control manifest, or state file.

## Impacted features

- FEATURE-0015: closed input-contract clarification and readiness proof.
- FEATURE-0016 through FEATURE-0026: no behavior change.

## Acceptance criteria for Kiro update

- [ ] The four exact input outcomes and ISO assigned-code semantics appear consistently in all allowed active F0015 authorities.
- [ ] No new Problem/violation code, resource, field, route, dependency, or ownership is introduced.
- [ ] The F0015-local conformance and REQ/AC mapping prove each clarified outcome.
- [ ] The readiness checker implements the complete per-route simulation and fails closed on the listed conditions.
- [ ] Exact canonical ledgers remain required by semantic validation.
- [ ] Only files within the strict Kiro allowlist changed.
- [ ] `make feature-0015-architecture-readiness`, `make vs000-contract-check`, `make phase2r-drift-check`, `git diff --check`, `mkdocs build --strict`, and `make structurizr-check` pass.

## Explicit instructions to Kiro

- Apply no decision beyond this handoff.
- Do not search the repository outside the strict allowlist.
- Do not modify any generated Kiro artifact or Go source.
- Stop, do not infer, if any decision cannot be represented using only the allowed files.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-12
- Notes: Bounded final input-contract clarification. No architecture exploration or unrelated repository discovery is authorized.
