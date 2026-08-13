# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-052
- Date: 2026-08-13
- Source discussion: FEATURE-0015 requirements-review recovery
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Close FEATURE-0015 validation-proof coverage and requirements-consistency gaps

## Summary

This is a bounded correction to make existing FEATURE-0015 behavior provable without changing product semantics. The approved feature authority already requires CloudPlatform and CloudProvider name uniqueness and registry-defined field validation for the canonical resource model. The active registry has no exact local conformance scenario for those behaviors, while the generated requirements artifact incorrectly claims unrelated cases prove them. This handoff adds two exact local proof cases, corrects their traceability, and makes the completeness checker fail closed on that mismatch. It also records the exact requirements-artifact consistency fixes that must be made only after the authorities validate.

## Classification

Correction — missing conformance mappings for already-decided behavior and inconsistent supplemental requirements prose.

## Existing approved baseline

- `REQ-F15-01` requires a unique CloudPlatform `metadata.name`.
- `REQ-F15-02` requires a unique CloudProvider `metadata.name` and a non-empty, unique, uppercase, assigned ISO-3166-1-alpha-2 `spec.operatingMarkets[]`.
- `REQ-F15-04` requires the topology schemas and their registered field constraints.
- `VS0-SCHEMA-008..014` already define the exact create-time required, optional, reference, cardinality, format, enumeration, and mutability constraints.
- `VS0-CF-F15-16` already defines deterministic scheduler expiry, and `VS0-CF-F15-21` already defines authenticated LIST/GET/read authorization outcomes.
- ADH-2026-045 through ADH-2026-051 remain controlling and are not superseded.

The current `VS0-CF-F15-01`, `02`, `12`, `13`, `14`, and `23` cases do not test duplicate CloudPlatform/CloudProvider names or general collection-create schema-constraint failure. They must not be represented as exact proof of those absent scenarios.

## Correction

### 1. Add exact validation conformance cases

Add exactly these two FEATURE-0015-local conformance cases. Do not alter the semantics of any existing case.

```text
id: VS0-CF-F15-32
owner: FEATURE-0015
inputs: an authorized CloudPlatform or CloudProvider collection create whose metadata.name duplicates an existing resource of the same kind in that kind's server-derived Platform-root scope
expectedState: unchanged
expectedError: ALREADY_EXISTS
expectedSideEffects: no duplicate resource persisted
gate: feature
```

```text
id: VS0-CF-F15-33
owner: FEATURE-0015
inputs: an authorized FEATURE-0015 collection create that, after closed-contract classification, violates a registry-declared semantic schema constraint in VS0-SCHEMA-008..014 other than a duplicate name, an unassigned ISO-3166-1 alpha-2 code covered by VS0-CF-F15-29, or a malformed/prohibited header/body condition covered by VS0-CF-F15-30
expectedState: unchanged
expectedError: VALIDATION_FAILED
expectedSideEffects: no resource persisted, no idempotency completion, no AuditEvent unless an already-approved audited-denial rule independently applies
gate: feature
```

`VS0-CF-F15-33` covers existing registered semantic constraints only: required fields, length, cardinality, enum, registered reference type/UID pinning, and the existing `operatingMarkets[]` uniqueness/uppercase/assigned-code constraints. It creates no field, route, Problem code, violation code, lifecycle transition, authorization rule, or persistence behavior.

### 2. Exact traceability correction

Update the F0015-local proof map and traceability matrix atomically:

- `REQ-F15-01` maps to `VS0-CF-F15-32` for CloudPlatform name uniqueness and to `VS0-CF-F15-33` for its schema constraints.
- `REQ-F15-02` maps to `VS0-CF-F15-32` for CloudProvider name uniqueness and to `VS0-CF-F15-33` for `operatingMarkets[]` and its other schema constraints; `VS0-CF-F15-29` remains the exact assigned-code case.
- `REQ-F15-04` maps to `VS0-CF-F15-33` for registry-declared topology field/reference constraints; `VS0-CF-F15-23` remains limited to its registered safe-denial/authorized-mismatch scenario.
- `REQ-F15-03` and `REQ-F15-06` map to `VS0-CF-F15-16` for scheduler expiry.
- The architecture authority explicitly records `VS0-CF-F15-21` as the local observable proof for authenticated LIST without read grant, filtered LIST, empty authorized LIST, and inaccessible GET/reference resolution. It must not be falsely assigned to an unrelated REQ or AC identifier.

### 3. Requirements-artifact consistency corrections after authority validation

Do not edit `.kiro/specs/**` during this handoff. Once this authority update and all validations pass, a separate bounded requirements revision must:

- preserve all three mandatory canonical ledgers verbatim;
- use only semantically matching proof cases and add the exact supplemental mappings above;
- state that scheduler expiry is an internal api-server transition invoked by the deterministic scheduler, not a public participation action; the public surface has exactly eight item-action routes;
- add an unnumbered normative read/list authorization boundary traced to `VS0-CF-F15-21`, without creating or renumbering a REQ or AC;
- state that prohibited concepts are never active behavior, while verbatim ledgers, traceability, and explicit negative-boundary sections may name them;
- state only the approved observable persistence boundary: no durable or external persistence dependency. It must not select an in-memory registry implementation;
- correct sole api-server status-writer coverage to `VS0-SCHEMA-008..014`.

## Rationale

The correction removes a false proof claim without changing what FEATURE-0015 does. It gives design and tasks exact, testable cases for all claimed create-time uniqueness and schema validation behavior, while retaining the existing scope, lifecycle, authorization, audit, and API contracts.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse existing FEATURE-0012 `ALREADY_EXISTS` / 409 and `VALIDATION_FAILED` / 422 Problem Details outcomes.
  - Reuse the existing VS-000 registry, F0015 proof map, readiness checker, contract checker, and semantic checker.
- Sovrunn-owned responsibility summary:
  - Register exact F0015-local conformance scenarios and fail closed when supplemental proof mappings overclaim their registered inputs.
- Non-goals summary:
  - No new API, resource field, scope, action grant, lifecycle state, authorization policy, persistence dependency, audit category, Problem code, violation code, Go implementation, or generated Kiro artifact.

## Phase impact

- Current phase allowed: Yes.
- Current phase boundary impact: FEATURE-0015 contract-closure correction only.

## Conflict check

- Conflicts with accepted DEC/RFC: No.
- Resolution required: Founder approval of this correction; no DEC, RFC, baseline, roadmap, or feature-sequence update.

## Required action

- Update the F0015 registry/specification/feature/architecture authorities, traceability, closure matrix, steering, control manifest, and deterministic checkers.
- Preserve ADH-2026-048 through ADH-2026-051 and add ADH-2026-052 to `feature.handoffs` without changing another manifest field.
- Do not update requirements, design, tasks, Go source, DEC, baseline, roadmap, feature sequence, or state files.

## Strict Kiro read/write boundary

Kiro may read only:

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
- `scripts/kiro-semantic-check.py`.

Kiro may write only the same files, plus this handoff for an application-status annotation. It must not search, glob, inspect, or follow references outside this list. If another file is required, stop and report its exact path and reason without opening it.

## Impacted features

- FEATURE-0015: exact local validation proof and accurate requirements traceability.
- FEATURE-0016 through FEATURE-0026: no behavior change.

## Acceptance criteria for Kiro update

- [ ] `VS0-CF-F15-32` and `VS0-CF-F15-33` match this handoff exactly and existing cases are unchanged.
- [ ] Registry, contract specification, F0015 feature/architecture authorities, traceability matrix, closure matrix, and steering list F15-01..33 consistently.
- [ ] The proof map contains the exact corrections in decision 2 and does not use a case outside its registered scenario.
- [ ] Readiness and contract checkers fail closed if name uniqueness, schema constraint, scheduler-expiry, or read/list proof mappings are missing or semantically overbroad.
- [ ] `FEATURE-0015.control.json` retains ADH-048/049/050/051 and adds ADH-052 with no other field change.
- [ ] No generated Kiro artifact, Go source, DEC, baseline, roadmap, feature sequence, or state file changes.
- [ ] `make feature-0015-architecture-readiness`, `make feature-contract-check FEATURE=FEATURE-0015`, `make feature-0015-formal-check`, `make vs000-contract-check`, `make phase2r-drift-check`, `git diff --check`, `mkdocs build --strict`, and `make structurizr-check` pass.

## Explicit instructions to Kiro

- Apply no decision beyond this handoff.
- Do not search outside the strict allowlist.
- Do not modify generated Kiro artifacts or Go source.
- Do not reset or bypass the requirements revision-attempt limit.
- Stop, do not infer, if this handoff cannot be represented with only the allowed files.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-13
- Notes: Approved bounded correction for already-decided validation proof coverage and requirements consistency.
