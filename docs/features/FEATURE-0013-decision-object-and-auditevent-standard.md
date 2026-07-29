---
doc_type: feature
id: FEATURE-0013
title: Decision Record and AuditEvent Standard
status: implemented_pending_final_merge
phase: 2
reuse_assessment_format_version: 1.0.0
depends_on:
  - FEATURE-0011
  - FEATURE-0012
ai_load_priority: feature
ai_summary: Approved FEATURE-0013 contract-only implementation for DecisionRecord, DecisionProfile, EvaluationResult, AuditEvent decision linkage, validation, schemas, and conformance.
controlling_handoff: ADH-2026-017
canonical_architecture: docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md
kiro_slug: decision-object-and-auditevent-standard
---

# FEATURE-0013 — Decision Record and AuditEvent Standard

## Purpose

FEATURE-0013 defines the Phase 2 contract-only decision and accountability
foundation for Sovrunn. It establishes the immutable `DecisionRecord` envelope,
versioned `DecisionProfile` definitions, normalized `EvaluationResult` evidence,
additive `AuditEvent` decision linkage, structural validation, schemas,
conformance fixtures, and compatibility evidence.

FEATURE-0013 does not implement a production decision service, persistence,
provider adapter, policy engine, workflow engine, AI runtime, cryptographic
execution, approval workflow, erasure/key-destruction runtime, or later deferred
capability.

Canonical architecture: `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`

Controlling handoff: ADH-2026-017 (Approved)

Canonical reuse standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

## Acceptance Criteria

1. Preserve FEATURE-0012 resource grammar, metadata, references, boundaries,
   strict decoding, Problem Details, schema baseline, and conformance mechanisms.
2. Define `DecisionRecord`, `DecisionProfile`, and `EvaluationResult` contracts
   without adding a shared `DecisionRequest`, pending-decision envelope, runtime
   service, or persistence model.
3. Preserve the six FEATURE-0012 `ScopeKind` values; model `ServiceInstance` as
   a typed subject/reference under its governing Project scope, never as a
   `ScopeKind`.
4. Preserve FEATURE-0012 `AuditEvent` base ownership while adding only the
   approved FEATURE-0013 decision-linkage extension and compatibility fixtures.
5. Enforce closed `DECISION_*` violation codes, append-only immutable records,
   structural-only trust and `SecurityExceptionRef` checks, and the
   pre-ADR-F13-002 no-cryptographic-execution boundary.
6. Provide executable conformance coverage for every approved FEATURE-0013
   scenario identifier and maintain Matrix E residual-risk evidence.
7. Keep INVARIANT_FOR_LATER and DEFERRED decisions as documentation/governance
   evidence only, with no source implementation.
8. Require final Codex/human review and `Final feature-review status: Approved`
   before merge approval.

## Feature-level reuse summary

Feature identity: FEATURE-0013 — Decision Record and AuditEvent Standard.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 API/resource grammar | Extend | Reuse every shared profile, metadata, reference, boundary, validation, Problem Details, schema, and compatibility primitive; add only decision/audit domain semantics and approved per-kind constraints. | Approved | ADH-2026-017; ADH-2026-012; RFC-0022 |
| DecisionRecord, DecisionProfile, EvaluationResult, and AuditEvent domain contract | Build | No mature external model supplies Sovrunn's complete sovereign, provider-neutral authority, projection, scope, audit-obligation, and downstream-adoption semantics. The common contract is Sovrunn differentiation; external concepts remain inputs rather than native core types. | Approved | ADH-2026-017; DEC-0026; RFC-0023 |
| Bounded decision-requirements and composition concepts | Reuse | Reuse applicable OMG DMN decision-requirements concepts without selecting, embedding, or recreating a DMN engine or workflow language. | Approved | ADH-2026-017; RFC-0023 |
| Audit transport and operational correlation | Reuse | CloudEvents may be an optional transport mapping and OpenTelemetry may carry operational correlation; neither becomes the canonical record or audit authority. | Approved | ADH-2026-017; RFC-0023 |
| Supply-chain evidence references | Reuse | Reuse DSSE/in-toto and SPDX/CycloneDX concepts by typed reference without duplicating their schemas or implementing signing. | Approved | ADH-2026-017; ADR-F13-002 |
| Canonicalization, digest, signing, and cryptographic verification | Reuse | Preserve algorithm-agile carriers and evaluate mature standards later; no algorithm, covered-field set, product, or service is selected in FEATURE-0013. | Deferred | ADH-2026-017; ADR-F13-002 |
| Policy/evaluator and workflow runtimes | Wrap | CEL, OPA, Cedar, durable-execution, and workflow products are later candidates behind owning-feature adapters or ports; FEATURE-0013 defines no wrapper or runtime implementation. | Deferred | ADH-2026-017; DEC-0036 |


## Capability assessment: FEATURE-0013 contract-only decision/audit foundation

### Identity

| Field | Value |
|---|---|
| Feature identity | FEATURE-0013 |
| Capability or decision-unit identity | FEATURE-0013 contract-only decision/audit foundation |
| Assessment owner | Sovrunn Architecture Owner |

### Classification

| Field | Value |
|---|---|
| Disposition | Build |
| Decision status | Approved |

### Analysis

| Field | Value |
|---|---|
| Assessment scope | Define and implement the approved Phase 2 contract-only DecisionRecord, DecisionProfile, EvaluationResult, additive AuditEvent decision-linkage, validation, schemas, conformance, compatibility, and evidence foundation. |
| Candidate category | Decision/audit semantic contracts, decision-requirements concepts, audit/event correlation, supply-chain evidence references, schema validation, and conformance standards. |
| Mature candidates / applicable standards | FEATURE-0012 API/resource grammar; RFC 9457 Problem Details; RFC 6901 JSON Pointer; JSON Schema 2020-12; selected DMN decision-requirements concepts; CloudEvents and OpenTelemetry concepts for optional mapping/correlation; DSSE/in-toto and SPDX/CycloneDX concepts by typed reference only. |
| Relevant candidate strengths | Mature standards provide portable syntax, validation, error paths, schema tooling, audit/correlation vocabulary, and evidence-reference patterns without forcing Sovrunn to adopt a product runtime. |
| Material candidate constraints | No external candidate supplies Sovrunn's full provider-neutral sovereignty, governance scope, profile, projection, immutable evidence, decision-linkage, and compatibility semantics. Runtime engines and cryptographic products remain deferred. |
| Rationale | Build the Sovrunn-owned decision/audit contract while extending FEATURE-0012 and reusing mature standards only where they do not change ownership or introduce runtime/product coupling. |
| Selected foundation or approach | Sovrunn-owned contract-only implementation using FEATURE-0012 resource grammar, immutable records, versioned profiles, embedded evaluation evidence, additive AuditEvent linkage, closed DECISION_* codes, structural trust/exception carriers, JSON Schema, Go bindings, and executable conformance fixtures. |
| Why Reuse is insufficient | Reuse alone cannot supply Sovrunn-specific provider-neutral governance authority, immutable decision semantics, profile ownership, AuditEvent decision linkage, Matrix E traceability, and closed conformance behavior. |
| Why Wrap is insufficient | Wrapping a decision/policy/workflow product would introduce runtime/product ownership that FEATURE-0013 explicitly defers and would not produce the shared canonical contract required by later features. |
| Why Extend is insufficient | Extending FEATURE-0012 is necessary but not sufficient because FEATURE-0012 deliberately excludes DecisionRecord/AuditEvent payload semantics; FEATURE-0013 must build the domain contract on top of inherited grammar. |
| Protected Sovrunn differentiation and long-term ownership | Sovrunn retains long-term ownership of the decision/audit contract, profile semantics, scope/audit linkage, closed DECISION_* registry, conformance IDs, and compatibility boundary. |

### Boundary

| Field | Value |
|---|---|
| Sovrunn-owned responsibility | Sovrunn owns the FEATURE-0013 DecisionRecord, DecisionProfile, EvaluationResult, additive AuditEvent decision-linkage, structural validation, schemas, conformance identifiers, compatibility fixtures, and Matrix E evidence for the approved Phase 2 contract-only scope. |
| Reused or extended responsibility | FEATURE-0012 owns common metadata, TypeMeta/ObjectMeta, ScopeKind, typed references, resource profiles, strict decoding, Problem Details, schema annotations, TypeBindings, baseline manifest/approvals, and generic Operation/AuditEvent base grammar. Mature external concepts are reused only by reference or optional mapping. |
| Responsibility/control boundary | FEATURE-0013 does not own provider adapters, policy/evaluator runtimes, workflow engines, persistence, AI/model runtime, cryptographic execution, approval workflow, erasure/key-destruction runtime, or stable API promotion. Later owning features must add those through separate architecture decisions and gates. |
| Data crossing the boundary | Versioned JSON/YAML schema documents, Go value types, typed references, metadata.scopeRef, normalized evaluation evidence, structural trust/exception carriers, decision-linkage references, conformance fixtures, and Matrix E evidence. Raw secrets, credentials, provider-native objects, prompts, unrestricted evaluator payloads, and runtime state do not cross into common contracts. |
| Control crossing the boundary | Local pure validation, schema conformance, strict decoding, and fixture execution occur in FEATURE-0013. Runtime authorization, persistence, provider calls, workflow orchestration, cryptographic signing/verification, and AI execution do not cross this boundary. |
| Adapter required | No |
| Adapter rationale | FEATURE-0013 defines common contracts and pure validation only. Adapters and runtime ports are later-feature responsibilities under DEC-0036 and are not implemented here. |
| Adapter or contract identifier | none |
| Vendor-native types allowed | No |

### Suitability

| Field | Value |
|---|---|
| Sovereignty and deployment fit | Local schemas, Go value types, pure validators, and offline conformance fixtures support disconnected, sovereign, and air-gapped interpretation without external services. |
| Security and trust | Scope authority is limited to metadata.scopeRef; ServiceInstance is a typed subject only; validation fails closed; public codes are closed; structural trust and SecurityExceptionRef carriers perform no cryptographic or workflow execution. |
| Operational and supportability | Stable DECISION_* codes, JSON Pointers, schema baseline evidence, conformance matrix rows, deterministic validation passes, and Matrix E evidence make failures reviewable and reproducible. |
| Licensing and supply-chain | Implementation relies on existing repository Go modules and open standards. No new provider SDK, policy engine, crypto product, workflow engine, database, or network dependency is introduced by FEATURE-0013. |
| Portability and provider-neutrality impact | Contracts are provider-neutral, resource-profile-aligned, and boundary-classified. Provider-native, runtime-native, and adapter-native details remain outside common schemas. |

### Phase and scope

| Field | Value |
|---|---|
| Allowed in current phase | Yes |
| Current-phase work | Contract-only Go value types, pure validators, schema additions, baseline manifest/approval updates, conformance fixtures, coverage matrix, Matrix E worksheet, and final review evidence. |
| Deferred work | Production decision runtime, persistence, policy/evaluator execution, provider adapters, AI/model runtime, workflow runtime, cryptographic execution, approval workflow, erasure/key-destruction runtime, stable API promotion, and later integration behavior. |
| Explicit non-goals | No shared DecisionRequest; no pending-decision envelope; no runtime service; no persistence; no provider, policy, workflow, AI, or crypto execution; no adapter interface; no erasure/key destruction; no provider-native core fields; no second scope authority; no ServiceInstance ScopeKind. |
| Exit or migration boundary | Material semantic changes require a new architecture decision, updated schemas, baseline compatibility evidence, migration path, conformance fixtures, and final review. Stable promotion remains separately gated. |
| Phase 2 non-goal acknowledgement | Phase 2 remains a model, standard, decision, audit, adapter-boundary, and simulation foundation. FEATURE-0013 does not authorize real runtime integrations or later-phase execution. |

### Risk mitigation

#### Risk-control matrix

| Risk | Preventive control | Detection control | Corrective path |
|---|---|---|---|
| Scope or ownership drift | metadata.scopeRef sole authority; six ScopeKind values; ServiceInstance subject-only rule | Scope validator, negative fixtures, architecture-boundary checker | Reject drift, correct schema/validator/docs, rerun conformance and review |
| Runtime or product leakage | Contract-only package/import guardrails; no provider/policy/workflow/crypto/AI runtime dependencies | Import tests, guardrails, feature boundary checks, Docker verification | Remove leaked dependency/behavior and restore deferred ownership |
| AuditEvent compatibility break | Preserve FEATURE-0012 base AuditEvent and add only approved decision-linkage extension | Schema baseline checks, compatibility fixtures, apiconform TypeBinding checks | Restore base compatibility or record approved schema-diff/migration evidence |
| Public error drift | Closed 32 DECISION_* violation-code registry | Code registry tests and conformance negative fixtures | Remove invented code or approve architecture change before use |
| Evidence or risk under-review | Matrix E worksheet and final review gate | Feature gate and final approval-review artifact | Hold merge until human review accepts or assigns corrective path |

| Field | Value |
|---|---|
| Applicable architecture risks | F13-R01 through F13-R31 as staged in `docs/reviews/feature-gates/FEATURE-0013-matrix-e-worksheet.md`. |
| Residual risk | Residual risk is accepted only for FEATURE-0013's Phase 2 contract-only scope and remains bound to documented owners, corrective paths, reassessment triggers, and later-feature gates. |
| Replacement risk | Medium |
| Reassessment triggers | Material architecture change; stable API promotion; first production decision runtime; first provider/policy/evaluator/AI/workflow/crypto integration; scope incident; AuditEvent compatibility break; public code change; Matrix E risk reassessment. |

### Traceability

| Field | Value |
|---|---|
| Related DEC / RFC / ADH references | DEC-0026; DEC-0036; RFC-0022; RFC-0023; ADH-2026-012; ADH-2026-017; ADR-F13-002 |
| Linked acceptance criteria | FEATURE-0013 requirements AC-13.1 through AC-13.18; tasks T-001 through T-033; architecture AD-001 through AD-045; F13-CF/F13-SCOPE/F13-EVAL/F13-SEC/F13-TRUST/F13-COMPAT conformance IDs. |
| Validation and review evidence | `docs/reviews/feature-gates/FEATURE-0013-approval-review.md`; `docs/reviews/feature-gates/FEATURE-0013-matrix-e-worksheet.md`; `.automation/logs/FEATURE-0013/implementation-summary.json`; FEATURE-0013 boundary checks; `./scripts/verify.sh`; Docker race test; schema baseline manifest. |

### Human-approval evidence

| Field | Value |
|---|---|
| Approving person or role | Sovrunn Architecture Owner |
| Approval date | 2026-07-29 |
| Approved ADH or assessment-review reference | ADH-2026-017; `docs/reviews/feature-gates/FEATURE-0013-approval-review.md` |
| Structured approval-evidence record | docs/reviews/reuse-assessments/FEATURE-0013-approval-evidence.md |
| Approval applies to | FEATURE-0013 Phase 2 contract-only implementation from automation baseline `be835b17d9ed9425e0797fd69a9cc9a9635071c9`, plus final approval/cleanup changes, including schemas, Go bindings, validators, fixtures, compatibility evidence, Matrix E staging, and docs. |

## Boundary decision

Decision: Approved

FEATURE-0013 owns only the normalized `DecisionRecord`, `DecisionProfile`,
`EvaluationResult`, and additive `AuditEvent` decision-linkage contracts for the
approved Phase 2 scope. It references inherited FEATURE-0012 `Operation`,
`Problem`, metadata, reference, scope, schema, and validation contracts without
redefining their ownership.

`DecisionRequest` remains calling-domain-owned. FEATURE-0013 defines no shared,
canonical, or common request schema/type/resource and no `PendingDecision` or
other non-final response envelope.

## Implementation evidence

Decision: Approved

Implementation completed through T-033 at automation baseline
`be835b17d9ed9425e0797fd69a9cc9a9635071c9`. Final review evidence is recorded
in `docs/reviews/feature-gates/FEATURE-0013-approval-review.md`.

Machine evidence reviewed:

- FEATURE-0013 architecture-boundary readiness check — PASS
- FEATURE-0013 requirements/design/tasks post-generation boundary checks — PASS
- `./scripts/verify.sh` — PASS
- Docker `go test -race ./...` — PASS
- Schema baseline manifest check — PASS, no missing files or digest mismatches
- Closed `DECISION_*` registry check — PASS, exactly 32 public violation codes

## Phase and scope

| Field | Value |
|---|---|
| Allowed in current phase | Yes |
| Current-phase work | Contract-only Go value types, pure validation, JSON schemas, schema baseline updates, conformance fixtures, compatibility evidence, Matrix E staging, and documentation. |
| Deferred work | Production decision service, policy/evaluator runtime, provider adapters, AI/model runtime, workflow/runtime orchestration, persistence, cryptographic execution, approval workflow, erasure/key-destruction runtime, stable API promotion, and later-domain integrations. |
| Phase 2 non-goal acknowledgement | Phase 2 remains a model, standard, decision, audit, adapter-boundary, and simulation foundation. FEATURE-0013 does not authorize real runtime integrations or later-phase execution. |

## Human-approval evidence

| Field | Value |
|---|---|
| Approving person or role | Sovrunn Architecture Owner |
| Approval date | 2026-07-29 |
| Approved ADH or assessment-review reference | ADH-2026-017; `docs/reviews/feature-gates/FEATURE-0013-approval-review.md` |
| Approval applies to | FEATURE-0013 Phase 2 contract-only implementation from automation baseline `be835b17d9ed9425e0797fd69a9cc9a9635071c9`, plus final approval/cleanup changes, including schemas, Go bindings, validators, fixtures, compatibility evidence, Matrix E staging, and docs. |

## Merge authorization

Final feature-review status: Approved

FEATURE-0013 is approved for pull-request review and merge consideration after
the final approval/cleanup change set is committed and the feature gate remains
passing. Any later material change requires rerunning FEATURE-0013 boundary
checks, repository verification, and final review.
