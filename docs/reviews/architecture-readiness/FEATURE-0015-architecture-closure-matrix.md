# FEATURE-0015 Architecture Closure Matrix

## Status

- Feature: FEATURE-0015 — Canonical Cloud Model and Alpha Migration Foundation
- Baseline: `ARCH-2026.08-PHASE2R-CANONICAL`
- Prepared: 2026-08-10
- Resolved: 2026-08-10
- Status: **Resolved — all gaps closed by ADH-2026-043 (approved by Sanjeev Kumar)**
- Purpose: consolidate every architecture-readiness finding before another requirements-stage attempt.

This is a closure artifact, not an Architecture Decision Handoff and not approval. It freezes the discovered gaps so they are resolved as one package rather than rediscovered in requirements review.

## Review boundary and sources

Sources were evaluated in source-of-truth order: current architecture baseline; canonical data model and contract catalog; accepted DEC-0037, DEC-0054, and DEC-0058; ADH-2026-041; the VS-000 contract registry/specification/traceability; then the FEATURE-0015 authority and generated requirements review.

Requirements generation and requirements review remain paused until the readiness gate passes. The existing generated requirements document is evidence only and must not be approved for design until every row below is closed and architecture-owner approval is recorded.

## Closure matrix

| ID | Area | Conflicting or missing authority | Classification | Required architecture-owner outcome | Repository closure targets | Status |
|----|------|----------------------------------|----------------|-------------------------------------|----------------------------|--------|
| ARC-F15-01 | Migration-record correction linkage | FEATURE-0015 and `VS0-STATE-010` require immutable corrections to create linked records. `VS0-SCHEMA-061` has only `record.predecessorRef`, which is reserved by `VS0-STATE-011` for the strict 1..9 milestone chain. No correction/supersession field or rule exists. | Correction | Choose one representable correction model: add an explicit immutable correction/supersession link that cannot substitute for the milestone predecessor chain, or prohibit correction records and require a new plan/run. Define the record-counting and ordering effect. | `VS-000-contract-registry.yaml` schemas 061 and states 010/011; F0015 feature and architecture authorities; VS-000 specification; traceability; conformance/checker coverage. | **RESOLVED** — ADH-2026-043 decision 1: corrections require a superseding plan and new migration run; prior records never altered or re-ordered |
| ARC-F15-02 | Migration sequence starts at Draft | The canonical model, canonical contract catalog, and ADH-2026-041 define `Draft → InventoryValidated → … → Completed`. `VS0-STATE-011` and F0015 define a 1..9 record sequence beginning at `InventoryValidated`; no authority defines whether Draft is a mutable pre-persistence plan state, a record, or intentionally excluded. | Clarification | Define the representation of `Draft` and make every source agree. Recommended direction: Draft is pre-persistence plan preparation; the signed immutable plan and append-only evidence sequence begin at `InventoryValidated`. | Canonical model/catalog migration text; ADH-2026-041 if necessary; registry state 011; F0015 feature/architecture; traceability. | **RESOLVED** — ADH-2026-043 decision 2: Draft is pre-persistence preparation; signed immutable plan precedes persisted 1..9 sequence |
| ARC-F15-03 | Alpha migration inventory and feature boundary | ADH-2026-041 requires coordinated classification/migration of provider, topology, catalog, enrollment, CLI/fixture, decision and audit references. F0015 dry-run evidence is only `Provider/topology` fixtures, while its non-goals exclude later catalog/enrollment semantics. `VS0-SCHEMA-060` nevertheless requires all alpha records covered. | Clarification | Publish an exact inventory table that states, for every alpha record family, whether F0015: (a) classifies it in the signed plan, (b) transforms it now, (c) retains it read-only, or (d) defers its implementation to its owning future feature. Define F0015's deterministic dry-run fixture set and its completion criterion without permitting dual authority. | ADH-2026-041 supplement or approved handoff; F0015 feature/architecture; schema 060 validation; `VS0-CF-MIG01`; traceability; F0016–0026 dependency notes where affected. | **RESOLVED** — ADH-2026-043 decision 3: F0015 classifies all families, transforms provider/topology now, defers catalog/enrollment/governance/placement/plugin to later owners; final cutover is FEATURE-0026 |
| ARC-F15-04 | Cross-provider reference and safe-denial semantics | Canonical contract rules require cross-provider references to deny by default without existence disclosure. Generated requirements use `VALIDATION_FAILED` for a cross-provider topology reference but use `RESOURCE_NOT_FOUND` + `VS0_AUTHORIZATION_SAFE_DENIAL` elsewhere. `VS0-CF-X03` registers 404 but no expected violation. | Correction | Separate authorized structural invalidity from inaccessible cross-provider resolution. Define the exact outcome for each; if safe denial includes `VS0_AUTHORIZATION_SAFE_DENIAL`, register it on `VS0-CF-X03`. Apply the decision consistently to topology and ExecutionTarget references. | Registry schema/reference validation and `VS0-CF-X03`; F0015 feature/architecture; VS-000 specification/traceability; semantic checker expectations. | **RESOLVED** — ADH-2026-043 decision 4: inaccessible cross-provider returns RESOURCE_NOT_FOUND/404 + VS0_AUTHORIZATION_SAFE_DENIAL; authorized invalid same-provider returns VALIDATION_FAILED/422; VS0-CF-X03 updated |
| ARC-F15-05 | Ownership/scope reference integrity | The canonical model says a CloudPlatform is owned by exactly one Organization and a participation joins exactly one CloudPlatform. Registry schema 008 requires `ownerOrganizationRef`, and schema 010 requires `cloudPlatformRef`, but neither defines UID equality against the resource's authoritative `metadata.scopeRef`. Contradictory graphs are therefore representable. | Clarification | Define the UID-pinned equality invariants: CloudPlatform `spec.ownerOrganizationRef.uid` equals its Organization scope UID; CloudProviderParticipation `spec.cloudPlatformRef.uid` equals its CloudPlatform scope UID. Specify the existing or newly approved failure mapping for a mismatch. | Registry schemas 008/010 validation; F0015 feature/architecture; local F0015 conformance/error mapping; traceability. | **RESOLVED** — ADH-2026-043 decision 5: UID-pinned invariants defined; mismatch returns VALIDATION_FAILED/422 + VS0_SCOPE_REFERENCE_MISMATCH; VS0-CF-F15-11 registered |
| ARC-F15-06 | Audit, correlation, redaction, and no-secret obligations | FEATURE-0013's AuditEvent contract is consumed, and architecture standards require material lifecycle/security decisions to have audit evidence and stable correlation. F0015 has no bounded event/correlation contract for participation lifecycle changes, plan publication, migration milestones/cutover, or security denials. | Clarification | State the F0015 events that produce/reuse AuditEvent, required correlation fields, safe projections/redaction, and the explicit no-secret rule. Reuse FEATURE-0013; do not create a second audit envelope or import FEATURE-0026's integration-only `VS0-CF-T01`. | F0015 feature/architecture; registry/schema 007 reference guidance if needed; traceability; requirements after approval. | **RESOLVED** — ADH-2026-043 decision 6: F0015 reuses FEATURE-0013 AuditEvent with intentional correlation, safe projection/redaction, no secrets; does not import FEATURE-0026 VS0-CF-T01 |
| ARC-F15-07 | F0015 architecture traceability has downstream-owned evidence | F0015 architecture §10 assigns `VS0-CF-HP01` to most local resources and `VS0-CF-F09` to participation. Those are owned by FEATURE-0026 and FEATURE-0023 respectively, whereas F0015 now owns `VS0-CF-F15-01..10`. This directly conflicts with the feature-local ownership boundary. | Correction | Replace downstream conformance references with the applicable F0015-local conformance IDs. Keep downstream IDs only as explicitly non-owning, future integration references where needed. | F0015 architecture §10; F0015 feature authority if needed; VS-000 traceability; architecture-readiness checker. | **RESOLVED** — ADH-2026-043 decision 7: §10 updated with F0015-local conformance IDs only; downstream HP01/F09 noted as non-owning |

## Findings classified as requirements/reviewer follow-up, not architecture decisions

| ID | Finding | Disposition after architecture closure |
|----|---------|----------------------------------------|
| RQ-F15-01 | Requirements §5.4 names the in-memory registry. | Remove or rephrase as the approved Phase 2R non-goal when requirements are revised; storage mechanics belong to design. No architecture decision is needed. |
| RQ-F15-02 | N-03 says excluded names may not appear while the exclusion section necessarily names them. | Permit explicit exclusion-only/reference-only mention; retain the active-authority anti-drift prohibition. No architecture decision is needed. |
| RQ-F15-03 | Reviewer asked to remove canonical REQ/AC ledgers and introduced-term mapping. | **Rejected as reviewer noise.** The requirements-stage semantic contract requires exact canonical REQ/AC ledgers, and the stage contract requires a terms section. Keep them; inherited FEATURE-0011/0012/0013 schemas remain references only. |
| RQ-F15-04 | `CURRENT_PHASE_CONTEXT.md` says F0015 is architecture-ready and requirements are not generated. | Update status only after this closure package is approved and applied; it must not be used to bypass an open readiness finding. |

## Required approval package

The architecture owner should approve or amend one handoff that contains explicit resolutions for ARC-F15-01 through ARC-F15-07. The handoff must:

1. classify ARC-F15-01 and ARC-F15-04/07 as corrections, and state whether ARC-F15-02/03/05/06 are clarifications or require a new decision;
2. list every impacted schema, state, conformance, failure mapping, feature authority, traceability, and checker change;
3. state compatibility and Phase 2R impact; and
4. name the exact post-approval architecture readiness checks that must pass before requirements regeneration.

## Exit criteria

- [x] Each ARC-F15 row is `RESOLVED`, with the approving ADH/DEC/ACR reference.
- [x] Architecture-owner approval is recorded against the complete impacted-document package.
- [x] Registry, feature authority, architecture boundary, contract specification, and traceability agree.
- [x] A feature-generic architecture-readiness check validates the seven categories represented by ARC-F15-01 through ARC-F15-07 before Kiro requirements generation can run.
- [ ] The architecture approval digest is refreshed only after the preceding checks pass.
- [ ] Requirements are regenerated once from that approved package, then reviewed for fidelity rather than architecture discovery.
