# FEATURE-0016 Architecture Closure Matrix

## Status

- Feature: FEATURE-0016 — Adapter Boundary and ExecutionTarget Qualification
- Baseline: `ARCH-2026.08-PHASE2R-CANONICAL`
- Prepared: 2026-08-14 (`FEATURE-0016-architecture-starting-dossier.md`)
- Resolved: 2026-08-20 by ADH-2026-058
- Status: **Resolved — all ARC-F16-01..10 gaps closed by ADH-2026-058 (approved by Sanjeev Kumar, 2026-08-20)**
- Purpose: consolidate every architecture-readiness finding from the FEATURE-0016 starting dossier before requirements-stage generation. Retained as historical review evidence; not deleted.

This is a closure artifact, not an Architecture Decision Handoff and not
approval by itself. It freezes the discovered gaps so they were resolved as
one bounded package (ADH-2026-058) rather than rediscovered in requirements
review.

## Review boundary and sources

Sources were evaluated in source-of-truth order: current architecture
baseline; canonical data model and contract catalog; accepted DEC-0036,
DEC-0042, DEC-0057; the FEATURE-0015 delegation boundary (ADH-2026-045); the
VS-000 contract registry/specification/traceability; then the
`FEATURE-0016-architecture-starting-dossier.md` and
`FEATURE-0016-architecture.md` proposed architecture drafts.

## Closure matrix

| ID | Area | Conflicting or missing authority | Classification | Required architecture-owner outcome | Repository closure targets | Status |
|----|------|----------------------------------|----------------|-------------------------------------|----------------------------|--------|
| ARC-F16-01 | API surface and route table | ExecutionTarget's registered schema had required spec/status fields, but no exact F0016 method/path table, create surface, request body, scope derivation, initial status, or local route conformance ledger. | Correction | Specify the complete F0016 public API surface, each route's request/headers/result/problem/status/audit effect, and its Go 1.22 feasible registration form. | `VS-000-contract-registry.yaml` VS0-SCHEMA-015, conformance VS0-CF-F16-01..122; F0016 architecture/feature authorities. | **RESOLVED** — ADH-2026-058 clauses 3,4: exactly five routes, one pre-ServeMux transport guard, complete request/response contract. |
| ARC-F16-02 | Backing-reference optionality | `spec.infrastructureStackRef` was required in VS0-SCHEMA-015, while the canonical data model allowed an ExecutionTarget without InfrastructureStack backing for external API/edge/federated realization. | Correction | Reconcile required/optional backing-environment semantics for the approved `synthetic-iaas` profile without foreclosing future target categories. | VS0-SCHEMA-015 required fields; F0016 architecture §4.2. | **RESOLVED** — ADH-2026-058 clause 1/2: `synthetic-iaas` requires both `spec.cloudProviderParticipationRef` and `spec.infrastructureStackRef`; other target categories remain deferred/future-only per the starting dossier §9. |
| ARC-F16-03 | Adapter contract | DEC-0036 required an adapter boundary, but the adapter contract had no exact operation set, input/output envelope, timeout/cancellation/failure classification, fact freshness rule, fake adapter behavior, or provenance/redaction rule. | Correction | Define the minimal F0016 qualification-observer contract and deterministic fake behavior without choosing a vendor or performing external calls. | VS0-SCHEMA-016/017; F0016 architecture §4.8. | **RESOLVED** — ADH-2026-058 clause 6: sole observer `sovrunn.synthetic-iaas-observer/v1`, exact fact/provenance contract, deterministic missing-fixture/timeout behavior. |
| ARC-F16-04 | Status/writer ownership | Schema 015 named a target-lifecycle-controller and writer 006 named a fake-adapter-and-qualification-controller for facts/results, but exact status-field ownership and handoff were not enumerated. | Correction | Assign every F0016 status/record field to exactly one writer, define initial values and publication/audit behavior, and prohibit direct client status writes. | VS0-SCHEMA-015..017, VS0-WRITER-006; F0016 architecture §4.9. | **RESOLVED** — ADH-2026-058 clause 7: `ExecutionTargetLifecycleService` is the sole committer of target state and internal fact/result records. |
| ARC-F16-05 | Maintenance/state vocabulary | VS0-STATE-004 had a compact six-state form (including `Draining`), while canonical maintenance authority described a richer vocabulary; `InfrastructureMaintenanceNotice` had no visible schema/owner. | Correction | Reconcile the state vocabulary; remove `Draining` and public `Qualifying`; define the maintenance-trigger input contract without a separate notice resource. | VS0-STATE-004; F0016 architecture §4.9/4.10. | **RESOLVED** — ADH-2026-058 clauses 7,8: `Active/Retired x Unqualified/Qualified/Rejected/Indeterminate`; `Draining` removed; maintenance enters/clears only via a deterministic fixture trigger, not an HTTP route or notice resource. |
| ARC-F16-06 | Local conformance completeness | Only `VS0-CF-F10` proved changed-epoch stale-work rejection. No F0016-local cases covered registration, fact freshness/expiry, qualified/rejected/indeterminate results, authorized versus inaccessible references, writer denial, audit failure, or zero external effect. | Correction | Add a complete F0016-owned conformance matrix with exact inputs, output state, Problem code/HTTP/violation, side effects, and gate for every observable requirement. | `conformance` VS0-CF-F16-01..122; traceability matrix. | **RESOLVED** — ADH-2026-058 §2: 122-row proof/outcome annex added, one row (or named finite equivalence family) per observable behavior. |
| ARC-F16-07 | Credential/authority boundary | The target-scoped SecretRef boundary was stated in the phase plan, but current schemas did not identify the allowed reference carrier, writers/readers, rotation behavior, or denial outcome. | Correction | Define a reference-only credential/authority boundary, or explicitly defer the field while preserving fake-adapter operation. | VS0-SCHEMA-015 forbidden-fields list; F0016 architecture §3.1. | **RESOLVED** — ADH-2026-058 clauses 1,9: no SecretRef extension is active; `adapterAuthorityRef` and any SecretRef extension field are rejected `UNKNOWN_FIELD`/400 (VS0-CF-F16-57,84). |
| ARC-F16-08 | Cross-feature safe-denial reuse | F0016 depends on F0015 but shares cross-CloudProvider target references; existing X03 is owned by F0015, while F0016 needed exact safe-denial and authorized-mismatch semantics for target/participation/stack resolution. | Clarification | State that F0016 defines its own local safe-resolution order and no-existence-disclosure/audit effects rather than reusing X03 as local proof. | VS0-CF-F16-09..13,17,19,47,66,102,104,111..112; F0016 architecture §4.2/4.4. | **RESOLVED** — ADH-2026-058 clause 2/§2: F0016-local safe-denial and authorized-mismatch proof cases defined; `VS0-CF-X03` remains a cross-feature/shared reference only, not F0016-local proof. |
| ARC-F16-09 | Audit/idempotency/replay matrix | Qualification, lifecycle publication, audit append failure, retry/replay behavior, and maintenance-epoch concurrency had no feature-local atomicity/replay matrix. | Correction | State exactly which F0016 outcomes are replayable, abort-only, audited, or non-audited. | VS0-CF-F16-01,08,20..23,31..35,39..43,46..47,49..50,76..78,81..82,86..89,99..104,112,121..122; F0016 architecture §4.7,4.9,4.10. | **RESOLVED** — ADH-2026-058 clauses 5,8,9: complete idempotency/audit/replay/shutdown matrix defined. |
| ARC-F16-10 | Stale context pointers | The baseline/phase-context "next feature" pointers still named F0015 after it merged, and the approved sequence had already selected F0016. | Correction | Correct the stale status pointers as an administrative follow-up; do not treat stale text as feature scope authority. | `docs/context/CURRENT_PHASE_CONTEXT.md`, `docs/context/CHATGPT_ARCHITECTURE_SESSION_PROMPT.md` (out of ADH-2026-058 write boundary; tracked as a separate administrative follow-up, not part of this closure package). | **NOTED — OUT OF SCOPE FOR THIS APPLICATION RUN.** ADH-2026-058's strict write boundary does not include `docs/context/CURRENT_PHASE_CONTEXT.md` or `docs/context/CHATGPT_ARCHITECTURE_SESSION_PROMPT.md`. This row remains open for a future administrative-only update; it does not block FEATURE-0016 requirements generation because it is a stale pointer, not a semantic gap. |

## Exit criteria

- [x] Each ARC-F16 row is `RESOLVED` (or explicitly noted out-of-scope for this application run), with the approving ADH reference.
- [x] Architecture-owner approval is recorded against the complete impacted-document package (ADH-2026-058, approved by Sanjeev Kumar, 2026-08-20).
- [x] Registry, feature authority, architecture boundary, contract specification, and traceability agree.
- [x] `make feature-0016-architecture-readiness`, `make feature-contract-check FEATURE=FEATURE-0016`, `make feature-0016-formal-check`, and `make vs000-contract-check` are wired and required before FEATURE-0016 requirements generation.
- [ ] Requirements are generated once from this approved package, then reviewed for fidelity rather than architecture discovery.
