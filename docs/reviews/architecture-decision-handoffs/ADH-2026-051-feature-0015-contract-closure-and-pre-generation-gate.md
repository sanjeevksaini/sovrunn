# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-051
- Date: 2026-08-12
- Source discussion: FEATURE-0015 contract-closure recovery
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Close FEATURE-0015 audit, routing, and authentication contract gaps

## Summary

This one bounded correction closes the three remaining architecture-level contradictions discovered by the FEATURE-0015 contract checker: inconsistent audit side effects for authenticated denials, impossible route-registration arithmetic, and absent FEATURE-0015-local proof for unauthenticated requests. It does not add product capability. It makes already-required observable behavior exact and introduces a deterministic pre-generation gate so Kiro cannot generate requirements, design, or tasks while a route contract is incomplete.

## Classification

Correction of internally inconsistent active authorities and missing local proof.

## Existing approved baseline

FEATURE-0015 reuses FEATURE-0012 Problem Details and FEATURE-0013 AuditEvent. REQ-F15-15 and `VS0-CF-F15-24` require AuditEvent evidence for resource mutations, participation lifecycle/expiry, authenticated authorization denials, and safe denials. ADH-2026-048 decision 4 requires a complete route pipeline with F0015-local conformance for every owned route and fail-closed rejection of an unspecified observable outcome or missing mapping.

The current registry conflicts with this baseline because `VS0-CF-F15-05`, `VS0-CF-F15-06`, and `VS0-CF-X03` use `expectedSideEffects: none` despite being covered denial categories. The current phrase “22 explicit Go 1.22 method/path registrations” also conflicts with the approved seven collection paths, seven item paths, eight action paths, and the no-PATCH participation boundary. Finally, inherited `VS0-CF-F01` cannot serve as F0015-local proof for `AUTH_REQUIRED` / 401.

## Correction

### 1. Exact durable audit boundary

Each of the following produces exactly one redacted FEATURE-0013 AuditEvent after authentication/current authorization has been evaluated, with request correlation and no secret or inaccessible-resource disclosure:

| F0015 outcome category | Required local evidence |
|---|---|
| successful resource collection create or PATCH | `VS0-CF-F15-24`, `VS0-CF-F15-26` |
| successful participation create/action or deterministic scheduler expiry | `VS0-CF-F15-16`, `VS0-CF-F15-24` |
| authenticated client attempt to write api-server-owned status or identity/metadata | `VS0-CF-F15-05`, `VS0-CF-F15-06` |
| CloudPlatform-root creation denial | `VS0-CF-F15-12` |
| authenticated collection LIST without read grant | `VS0-CF-F15-21` |
| authenticated header/body forged-grant attempt | `VS0-CF-F15-22` |
| authenticated inaccessible cross-provider/safe-denial reference | `VS0-CF-F15-23`, `VS0-CF-X03` |

The following do **not** produce a durable AuditEvent: missing or invalid authentication; malformed/prohibited body/header/field input; unsupported method/media type; stale `If-Match`; same-key replay; changed-digest idempotency conflict; pair-uniqueness conflict; or an invalid participation source state. They may produce redacted logs/metrics under the existing observability baseline, but logs/metrics are not AuditEvents.

If an AuditEvent required by the first table cannot be appended, the request returns inherited `INTERNAL_ERROR` / 500, retains safe non-disclosure, and publishes no resource mutation, lifecycle/status transition, or idempotency completion. This extends the existing ADH-2026-049 non-publication rule to audited denials without changing its code or persistence boundary. An already-appended AuditEvent may remain after a later in-process publication failure, as constrained by ADH-2026-048 decision 3.

### 2. Exact route model and registration arithmetic

FEATURE-0015 owns exactly 22 logical endpoint paths:

- seven collection paths: one per owned kind;
- seven item paths: one per owned kind;
- eight CloudProviderParticipation action paths.

It registers exactly 35 explicit Go 1.22 `http.ServeMux` method/path patterns:

| Path category | Patterns | Count |
|---|---|---:|
| seven collections | `GET` LIST and `POST` create each | 14 |
| six PATCHable resource items | `GET` and `PATCH` each | 12 |
| CloudProviderParticipation item | `GET` only; no PATCH | 1 |
| eight participation actions | `POST` each | 8 |
| **Total** |  | **35** |

No path-only handler registration, wildcard registration, reflection, or internal HTTP-method dispatch is allowed. Unsupported methods return the established method-not-allowed behavior and do not create a new endpoint path.

### 3. F0015-local unauthenticated-request conformance

Add one and only one new F0015-local conformance case:

```text
id: VS0-CF-F15-31
owner: FEATURE-0015
inputs: missing or invalid authentication on any F0015 owned method/path pattern
expectedState: unchanged
expectedError: AUTH_REQUIRED
expectedSideEffects: no mutation, no lifecycle/status write, no idempotency record, no AuditEvent
gate: security
```

`VS0-CF-F15-31` is the exact local proof for the 401 branch of the required all-route pipeline. It does not replace or modify FEATURE-0018 `VS0-CF-F01`, and it creates no new top-level Problem or violation code.

### 4. Mandatory deterministic pre-generation gate

The existing generic `feature-contract-check.py` and two F0015 formal models are enforcement tooling, not architecture authorities. Before Kiro may generate FEATURE-0015 requirements, design, or tasks, all must pass:

```bash
make feature-0015-architecture-readiness
make feature-contract-check FEATURE=FEATURE-0015
make feature-0015-formal-check
```

The feature-contract checker must fail closed unless every F0015 method/path operation has exact request, authentication, authorization, safe-denial, validation, response, state, audit, idempotency/race, and F0015-local conformance coverage. It must verify the 22-logical-path/35-method-pattern arithmetic and audit-side-effect agreement. The formal check must model-check the participation lifecycle/independent holds and audit-before-publication/idempotency invariants with TLC. Missing TLC is a failure, not a skipped check.

## Rationale

The correction prevents requirements, design, and tasks from being used as an architecture discovery mechanism. It resolves all currently known contract contradictions in one package and converts the closure conditions into deterministic gates before another Kiro stage consumes credits.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse FEATURE-0012 `AUTH_REQUIRED` / 401, `INTERNAL_ERROR` / 500, and existing method-not-allowed behavior.
  - Reuse FEATURE-0013 AuditEvent and ADH-2026-048/049 atomic-publication constraints.
  - Reuse the existing F0015 architecture-readiness checker and VS-000 registry; add deterministic coverage enforcement and bounded TLA+ models rather than another prose-only review.
- Sovrunn-owned responsibility summary:
  - Align F0015 local conformance with existing observable behavior and enforce complete pre-generation proof.
- Non-goals summary:
  - No resource, field, route path, scope, action grant, lifecycle state, persistence dependency, external call, new Problem/violation code, future-feature behavior, Go feature implementation, or generated Kiro artifact.

## Phase impact

- Current phase allowed: Yes.
- Current phase boundary impact: correction and architecture-validation tooling for F0015 only.

## Conflict check

- Conflicts with accepted DEC/RFC: No.
- Resolution required: Founder approval of this bounded correction; no DEC, RFC, baseline, roadmap, or feature-sequence update.

## Required action

- Update only F0015 feature/architecture authorities, registry/specification, traceability, closure matrix, steering, F0015 readiness/registry checkers, feature contract/formal checker tooling, and the F0015 control manifest.
- Preserve ADH-2026-048, ADH-2026-049, and ADH-2026-050 in `feature.handoffs`; add ADH-2026-051 without changing another manifest field.
- Do not modify requirements, design, or tasks during this architecture correction.
- Regenerate requirements only after every required pre-generation command passes.

## Strict Kiro read/write boundary

Kiro must read only these files:

- this handoff;
- `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`;
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`;
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`;
- `docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md`;
- `.kiro/steering/slice0-contract.md`;
- `.automation/features/FEATURE-0015.control.json`;
- `scripts/feature-0015-architecture-readiness-check.py`;
- `scripts/vs000-contract-check.py`;
- `scripts/feature-contract-check.py`;
- `scripts/run-feature-0015-formal-checks.sh`;
- `scripts/generic-kiro-boundary-check.py`;
- `docs/formal/feature-0015/ParticipationLifecycle.tla`;
- `docs/formal/feature-0015/ParticipationLifecycle.cfg`;
- `docs/formal/feature-0015/IdempotencyAuditPublication.tla`;
- `docs/formal/feature-0015/IdempotencyAuditPublication.cfg`;
- `Makefile`.

Kiro may write only the same files, plus this handoff for a status annotation. It must not search, glob, inspect, or follow references outside this list. If any additional file seems necessary, it must stop and report the exact path and reason without opening it.

Kiro must not modify Go source, `.kiro/specs/**/requirements.md`, `design.md`, `tasks.md`, any DEC, baseline, roadmap, feature sequence, or state file.

## Impacted features

- FEATURE-0015: closed local route contract and deterministic pre-generation checks.
- FEATURE-0016 through FEATURE-0026: no behavior change.

## Acceptance criteria for Kiro update

- [ ] The listed audited denial cases have exact redacted AuditEvent side effects; listed non-audited cases expressly remain non-AuditEvent outcomes.
- [ ] Every required-AuditEvent append failure maps to `INTERNAL_ERROR` / 500 with no observable mutable publication and safe non-disclosure.
- [ ] Authorities consistently state 22 logical paths and 35 explicit Go 1.22 method/path registrations, with no path-only/internal dispatch.
- [ ] `VS0-CF-F15-31` is registered exactly as specified, and all F0015 range/mapping/checker references are updated atomically.
- [ ] The pre-generation gate verifies method/path coverage, audit-side-effect agreement, and local proof for every observable branch.
- [ ] The two F0015 TLA+ models represent the stated lifecycle and publication invariants; their runner fails closed without TLC.
- [ ] The control manifest preserves ADH-048/049/050 and adds ADH-051 with no other field change.
- [ ] No generated Kiro artifact or Go implementation source changed.
- [ ] `make feature-0015-architecture-readiness`, `make feature-contract-check FEATURE=FEATURE-0015`, `make feature-0015-formal-check`, `make vs000-contract-check`, `make phase2r-drift-check`, `git diff --check`, `mkdocs build --strict`, and `make structurizr-check` pass.

## Explicit instructions to Kiro

- Apply no decision beyond this handoff.
- Do not search outside the strict allowlist.
- Do not install software or download dependencies; if TLC is unavailable, stop and report it as an environment prerequisite after applying only the representable authority changes.
- Do not modify generated Kiro artifacts or Go source.
- Stop, do not infer, if a decision cannot be represented using only the allowed files.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-12
- Notes: Approved exact audit boundary, 35-registration route model, F15-31 local authentication proof, and fail-closed pre-generation gate.
