# Phase 2R Architecture Completeness Audit

## Status

- Baseline: `ARCH-2026.08-PHASE2R-CANONICAL`
- Audit date: 2026-08-11
- Scope: Phase 2R canonical contracts and the admission path from architecture to requirements, design, and tasks.
- Classification: architecture-readiness review; **not** an Architecture Decision Handoff and **not** an approval to change architecture.
- **Superseded update (2026-08-11, ADH-2026-045):** ARC-2R-001 and ARC-2R-002 (migration evidence cardinality and migration continuation ownership) are moot — ADH-2026-045 (canonical bootstrap replaces alpha runtime migration; DEC-0059 supersedes DEC-0058) removes `CanonicalMigrationPlan`, `CanonicalMigrationRecord`, and the migration controller entirely, so there is no migration cardinality or continuation model to resolve. ARC-2R-003 (ExecutionTarget required-vs-forbidden status) is resolved by full attribution of `ExecutionTarget` to FEATURE-0016; FEATURE-0015 introduces no ExecutionTarget contract at all. This audit's findings below are retained as historical record; they are not deleted. Kiro requirements generation for FEATURE-0015 is unblocked pending a passing run of the revised `make feature-0015-architecture-readiness` admission gate.

## Purpose and boundary

This audit responds to requirements-stage findings that exposed unspecified architecture after requirements generation had already begun. It distinguishes three states that must not be conflated:

| State | Meaning |
|---|---|
| Named canonical concept | The data model or catalog establishes that a concept exists and gives its broad purpose. |
| Feature-architected concept | A feature authority binds the concept to an exact scope, ownership, lifecycle, cardinality, validation, error, audit, and conformance contract. |
| Requirements-ready concept | The feature-architected contract has passed the admission gate below and can be rendered into requirements without requiring Kiro to choose product semantics. |

The audit does not attempt to pre-design every future feature. The canonical catalog explicitly reserves exact payload shapes for approved feature architecture. Instead, it identifies cross-feature semantic gaps that must be decided before a feature can be requirements-ready.

## Authorities examined

In source-of-truth order:

1. `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
2. `docs/architecture/canonical/sovrunn-finalized-data-model.md`
3. `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`
4. DEC-0037 through DEC-0058 and their controlling handoffs, with ADH-2026-041, ADH-2026-043, and ADH-2026-044 examined for FEATURE-0015
5. `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`, specification, and traceability matrix
6. Phase 2R rebaseline, feature sequence, FEATURE-0015 authority, and its architecture boundary

## Phase 2R admission inventory

| Feature | Registry-owned contract surface | Feature authority present | Requirements-ready now | Audit disposition |
|---|---|---:|---:|---|
| FEATURE-0015 | VS0-SCHEMA-008..015, 060..061 | Yes | **No** | Blocked by ARC-2R-001 through ARC-2R-004. |
| FEATURE-0016 | VS0-SCHEMA-016..017 | No | No | Not yet eligible for requirements architecture audit. |
| FEATURE-0017 | VS0-SCHEMA-018..019 | No | No | Not yet eligible for requirements architecture audit. |
| FEATURE-0018 | VS0-SCHEMA-020..024 | No | No | Not yet eligible for requirements architecture audit. |
| FEATURE-0019 | VS0-SCHEMA-025..028 | No | No | Not yet eligible for requirements architecture audit. |
| FEATURE-0020 | VS0-SCHEMA-029..030 | No | No | Not yet eligible for requirements architecture audit. |
| FEATURE-0021 | VS0-SCHEMA-031..035 | No | No | Not yet eligible for requirements architecture audit. |
| FEATURE-0022 | VS0-SCHEMA-036..042, 056 | No | No | Not yet eligible for requirements architecture audit. |
| FEATURE-0023 | VS0-SCHEMA-043..045 | No | No | Not yet eligible for requirements architecture audit. |
| FEATURE-0024 | VS0-SCHEMA-046, 048 | No | No | Not yet eligible for requirements architecture audit. |
| FEATURE-0025 | VS0-SCHEMA-049 | No | No | Not yet eligible for requirements architecture audit. |
| FEATURE-0026 | Cross-feature integration/conformance owner | No | No | Not yet eligible for requirements architecture audit. |

`No` in this table does not mean that a future feature is architecturally rejected. It means that a feature may not enter requirements until it has a feature authority and passes the admission gate. The current baseline provides semantic direction, not a substitute for feature-level contracts.

## Confirmed architecture blockers

| ID | Finding | Evidence | Why requirements/design/tasks cannot decide it | Required owner outcome |
|---|---|---|---|---|
| ARC-2R-001 | **[SUPERSEDED by ADH-2026-045/DEC-0059 — moot]** Migration evidence has incompatible cardinality. `CanonicalMigrationRecord` requires one singular source/target/classification/transform mapping, while ADH-2026-044 requires exactly one nine-milestone chain per migration run. The provider-topology run has multiple mappings. | VS0-SCHEMA-061; VS0-STATE-011; ADH-2026-044; FEATURE-0015 §7.3 | Any implementation must invent aggregation, placeholder identities, repeated chains, or a second evidence type. Each changes observable product semantics. | Superseded: ADH-2026-045 removes `CanonicalMigrationPlan`/`CanonicalMigrationRecord` entirely; there is no migration evidence cardinality to resolve. |
| ARC-2R-002 | **[SUPERSEDED by ADH-2026-045/DEC-0059 — moot]** Migration continuation ownership is not representable. The signed plan claims all alpha families are classified, FEATURE-0015 transforms only provider/topology, and later feature owners are said to append their deferred runs. No authority defines which feature owns each run, how a later feature is authorized to append under a plan published by FEATURE-0015, or how a superseding plan affects deferred runs. | ADH-2026-041; ADH-2026-043 decision 3; ADH-2026-044; FEATURE-0015 §7.5; VS0-WRITER-020/021 | A later feature would have to choose authority delegation, plan version binding, and failure/retry semantics. Global completion cannot be proven deterministically. | Superseded: ADH-2026-045 removes the migration-controller and migration-run model entirely; there is no continuation ownership to define. |
| ARC-2R-003 | **[RESOLVED by ADH-2026-045/DEC-0059]** ExecutionTarget is both required to have status and forbidden to persist status in FEATURE-0015. The registry marks several `status.*` fields required in VS0-SCHEMA-015, while the F0015 authority and boundary prohibit F0015 from introducing, persisting, defaulting, validating, or writing them; FEATURE-0016 owns introduction and activation. | VS0-SCHEMA-015; FEATURE-0015 REQ-F15-05; FEATURE-0015 §2.3; FEATURE-0015 architecture §§3–4 | Requirements cannot decide whether an F0015-created ExecutionTarget validates as incomplete, carries absent fields, uses a separate identity schema, or waits to exist until FEATURE-0016. | Resolved: ADH-2026-045 attributes ExecutionTarget in its entirety to FEATURE-0016 (identity, schema, routes, status, writer, conformance, target lifecycle); FEATURE-0015 introduces no ExecutionTarget contract and ends at InfrastructureStack. |
| ARC-2R-004 | **The pre-requirements gate validates text agreement, not semantic satisfiability.** It passed ADH-2026-044 even though the record schema and state machine could not both be true for a multi-mapping run. | FEATURE-0015 readiness output; VS0-SCHEMA-061; VS0-STATE-011; requirements review blocker | Kiro is being used to discover contradictions that a deterministic architecture gate must find before prompt creation. | Add a feature-generic contract-completeness check implementing the admission criteria below, including cardinality-versus-state-machine and field-ownership-versus-requiredness checks. It must run before prompt/context rendering. |
| ARC-2R-005 | **The phase status reports FEATURE-0015 as architecture-ready despite an open architecture blocker.** | `docs/context/CURRENT_PHASE_CONTEXT.md` FEATURE-0015 status row | Agents can incorrectly treat a stale status as permission to regenerate requirements. | After the architecture package is resolved, update phase context atomically with the actual gate result. Until then, the status must not assert readiness. |

## Architecture guardrails to add before requirements

The current catalog’s section 19 is a valuable completion test, but it needs a machine-checkable feature admission record. For every resource or cross-resource protocol touched by a feature, the record must answer all of the following before a requirements prompt is rendered.

| Admission dimension | Required proof |
|---|---|
| Semantic boundary | One authoritative purpose and explicit non-goals. |
| Identity and scope | Resource identity, scope kinds, UID-pinned reference rules, and no-existence-disclosure rules. |
| Cardinality | One-per-what, collection bounds, aggregation rules, and relationship uniqueness. |
| Lifecycle | Initial/final states, transitions, terminal meaning, correction/retry/supersession semantics. |
| Field staging | Required/optional fields for this feature; introduction and activation owner; no field may be both required now and prohibited now. |
| Authority | One writer per mutable field/condition, delegations, forbidden writers, and concurrency behavior. |
| Security and audit | Classification, redaction/projections, stable errors/violations, AuditEvent/correlation, and no-secret rules. |
| Cross-feature continuation | Explicit future owner, activation feature, handoff input/output, and deferred-state behavior. |
| Conformance | Positive, negative, stale, authorization, idempotency, and no-side-effect scenarios that are satisfiable from the declared schema and state model. |

### Mandatory fail-closed checks

The pre-requirements checker must reject a feature when any of these is true:

1. a schema has singular mapping fields while its state machine asserts an aggregate chain without an aggregation rule;
2. one feature is prohibited from persisting a field that the active schema requires on creation;
3. a state says a feature or run is complete without defining its exact identity and cardinality;
4. a deferred run or field has no activation feature and authorized writer;
5. a conformance case requires a state, error, or side effect that no schema/state/authority combination can produce;
6. feature status says architecture-ready while any audit row is open; or
7. the feature authority, registry, traceability, and architecture boundary do not agree on the same exact IDs and owners.

## Required closure sequence

1. Freeze FEATURE-0015 requirements, design, and tasks generation. The present generated requirements artifact is non-authoritative evidence only.
2. Have the architecture owner resolve ARC-2R-001 through ARC-2R-003 as one controlled architecture package. This must be a single coherent migration model, not another local patch.
3. Apply the approved package atomically through the architecture handoff process: canonical model/catalog where necessary, registry, state machines, conformance, traceability, FEATURE-0015 authority/boundary, and phase context.
4. Implement ARC-2R-004’s generic admission check and make it a prerequisite of `render-prompt.py` / `ff-kiro-stage-auto` for requirements, design, and tasks.
5. Re-run this audit’s FEATURE-0015 row. It may move to requirements-ready only when every admission dimension is proven.
6. For FEATURE-0016 through FEATURE-0026, create and audit each feature authority before its first requirements run. Do not pre-fill detailed design choices now; use the same gate to prevent their discovery from being deferred to Kiro.

## Exit criteria for reopening FEATURE-0015 requirements

- [ ] ARC-2R-001, ARC-2R-002, and ARC-2R-003 are resolved by explicit architecture-owner decisions.
- [ ] Every affected authority agrees, including schema, state, conformance, feature scope, traceability, and phase context.
- [ ] The generic pre-requirements admission checker implements the mandatory fail-closed checks.
- [ ] The checker passes against FEATURE-0015 using only approved authorities.
- [ ] A human confirms that FEATURE-0015 is architecture-ready after reviewing this audit and the closure package.

Only then may Kiro be used to generate FEATURE-0015 requirements again. A successful requirements validation is necessary afterward, but it is not evidence that architecture was complete beforehand.
