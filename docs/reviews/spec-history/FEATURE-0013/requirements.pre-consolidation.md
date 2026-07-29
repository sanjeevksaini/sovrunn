---
feature_id: FEATURE-0013
feature_name: Decision Record and AuditEvent Standard
phase: "Phase 2: Reuse-First PaaS Fabric Foundation"
stage: Requirements
status: draft
review_status: PENDING_HUMAN_REVIEW
kiro_slug: decision-object-and-auditevent-standard
controlling_handoff: ADH-2026-014
controlling_handoffs:
  - ADH-2026-014
  - ADH-2026-015
  - ADH-2026-016
controlling_architecture: docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md
reuse_assessment_format_version: 1.0.0
depends_on:
  - FEATURE-0012
---

# FEATURE-0013 — Decision Record and AuditEvent Standard

> **Status: SUPERSEDED PRE-CONSOLIDATION DRAFT — DO NOT USE FOR GENERATION,
> REVIEW, DESIGN, TASKS, OR IMPLEMENTATION.** This file preserves the former
> multi-handoff requirements history only. Proposed ADH-2026-017 and the
> consolidated architecture replace its normative effect after approval.

Stage: Requirements

> **ADH-2026-016 contract boundary clarification notice (Approved 2026-07-27).**
> ADH-2026-014, ADH-2026-015, and ADH-2026-016 are the **joint controlling
> handoffs** for this feature. Their precise relationships are:
> ADH-2026-015 supersedes **only** the earlier six-scope `AuditEvent` statement
> in ADH-2026-014; all other ADH-2026-014 decisions remain controlling.
> ADH-2026-016 is a **clarification** that introduces no new architecture and
> clarifies previously unresolved contract boundaries of both ADH-2026-014 and
> ADH-2026-015 — `DecisionRecord` scope authority via `metadata.scopeRef` as the
> sole scope source; the closed ordered provider-neutral sensitivity vocabulary
> `PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED` with mandatory profile-bound
> rules; the structural-only security-validation boundary; the structural-trust
> versus deferred-cryptographic boundary (ADR-F13-002); the exact conformance
> coverage `F13-CF-01` through `F13-CF-28`; the SUPERSEDED disposition of
> six-scope `AuditEvent` evidence; and the completion of the
> `DecisionObject`-to-`DecisionRecord` migration in active normative documents.
> The seven-value `ScopeKind` vocabulary approved by ADH-2026-015 is preserved.
> ADH-2026-016 authorizes no runtime capability and no product, provider, or
> algorithm selection. Section 18 restates these clarified boundaries as
> requirements; sections 4, 12, and 17 are updated accordingly. This document
> returns to **architecture-remediation status**; it keeps
> `review_status: PENDING_HUMAN_REVIEW` and gate status
> `PENDING_INDEPENDENT_REVIEW`, and it does **not** self-approve or advance any
> stage. Design, tasks, and implementation remain unauthorized until this
> requirements document receives fresh independent human approval.

> **ADH-2026-015 scope-vocabulary supersession notice (Approved 2026-07-27).**
> For the scope-vocabulary correction governed by this notice, ADH-2026-014 and
> ADH-2026-015 are the controlling pair. ADH-2026-016 is also joint controlling
> for FEATURE-0013 overall and clarifies the additional contract boundaries
> listed in the preceding ADH-2026-016 notice. ADH-2026-015 supersedes **only**
> the earlier six-scope `AuditEvent` statement. Under ADH-2026-015, `ScopeKind` is the single canonical shared
> vocabulary with exactly seven values — `Platform`, `Organization`,
> `OrganizationUnit`, `Tenant`, `Project`, `Provider`, and `ServiceInstance` —
> `ServiceInstance` is additive, `Provider` is retained, and `AuditEvent`
> permits all seven values using `metadata.scopeRef` as its sole scope
> identity. **Section 17 is authoritative** for scope vocabulary. Where earlier
> text in this document (for example sections 1.5, 8.2, and AC-24) refers to
> "six" scopes or omits `Provider`, section 17 controls.
>
> The earlier `APPROVED_FOR_DESIGN: READY_FOR_REVIEWER` self-assessment is
> **superseded for continued design work** by ADH-2026-015 and by this
> revision. Section 12 now records the gate status as
> `PENDING_INDEPENDENT_REVIEW`. This requirements document is returned to
> `PENDING_HUMAN_REVIEW`. Requirements must receive **fresh independent
> approval** before design, tasks, or implementation resume. This document does
> not self-approve.

## 1. Introduction

FEATURE-0013 establishes the common, immutable, provider-neutral contract for governed decisions and audit-event linkage across Sovrunn.

### 1.1 Problem

Sovrunn must represent decisions ranging from a trivial allow/deny to a complex, hierarchical, multi-evaluator composite without changing the meaning of core objects. The representation must be deterministic, explainable, auditable, provider-neutral, horizontally scalable, and deployable in connected, disconnected, and air-gapped environments.

### 1.2 Scope

This feature owns:

- the common `DecisionRecord` envelope and immutable lifecycle;
- the versioned `DecisionProfile` extension mechanism;
- atomic `EvaluationResult` normalization contract;
- deterministic composition and bounded dependency-graph semantics;
- decision, evaluation, audit, operation, trace, and request correlation;
- `AuditEvent` linkage and local atomic acceptance contract;
- sovereign security, projection, portability, and evidence requirements;
- conformance schemas, validation, fixtures, and compatibility rules.

### 1.3 Phase 2 scope constraint

FEATURE-0013 is contract-only within Phase 2. It defines schemas, vocabularies, validation, conformance fixtures, and documentation. It does not implement a production decision service, audit service, persistence layer, orchestration runtime, or any provider-specific adapter. It also does not define `EvaluatorAdapter`, `PolicyEngineAdapter`, `IdentityProviderAdapter`, `SecretProviderAdapter`, `ObservabilityAdapter`, or any shared adapter interface; those interfaces are owned by FEATURE-0016 or the later owning feature. FEATURE-0013 owns only the normalized `DecisionRecord`, `DecisionProfile`, `EvaluationResult`, `AuditEvent`, validation, and adapter-facing data requirements. Adapter names appearing in this document are illustrative future consumers/producers only and generate no interface or implementation artifact here (DEC-0036 boundary intent preserved).

### 1.4 Terminology correction

Per ADH-2026-014 approval, `DecisionRecord` replaces the ambiguous common name `DecisionObject` used in the Phase 2 spine. This is a controlled terminology correction, not a new concept.

### 1.5 AuditEvent scope correction

Per ADH-2026-015 approval, `AuditEvent` scope is expanded from the FEATURE-0012 Organization-only alpha profile to all seven canonical `ScopeKind` values: `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, and `ServiceInstance`. `ServiceInstance` is additive to the FEATURE-0012 shared vocabulary and `Provider` is retained. `AuditEvent` uses `metadata.scopeRef` as its sole canonical scope identity. Section 17 is authoritative.

> **Superseded historical statement (ADH-2026-014).** ADH-2026-014 originally
> expanded `AuditEvent` to a six-scope list (`Platform`, `Organization`,
> `OrganizationUnit`, `Tenant`, `Project`, `ServiceInstance`) that omitted
> `Provider`. That six-scope statement is **superseded by ADH-2026-015**, which
> is authoritative and establishes the single canonical seven-value `ScopeKind`
> vocabulary above (retaining `Provider` and adding `ServiceInstance`), with
> `AuditEvent` permitting all seven values via `metadata.scopeRef`. Section 17
> controls.

### 1.6 Controlling references

- ADH-2026-014 (Approved, 2026-07-27), ADH-2026-015 (Approved, 2026-07-27), and ADH-2026-016 (Approved, 2026-07-27) — joint controlling handoffs. ADH-2026-015 supersedes only the earlier six-scope `AuditEvent` statement; all other ADH-2026-014 decisions remain controlling. ADH-2026-016 is a clarification that clarifies previously unresolved contract boundaries of ADH-2026-014 and ADH-2026-015 without introducing new architecture and without authorizing any runtime capability or product/provider/algorithm selection.
- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md` (sections 6.1, 12.1, 12.3, 12.4, 17, and 30 for the ADH-2026-016 clarified boundaries)
- `docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md`
- DEC-0026 — Reuse Before Build
- DEC-0036 — Adapter Boundaries Before External Integration
- RFC-0023 — Decision and Audit Standard
- FEATURE-0012 API/resource grammar foundation

## 2. Glossary

| Term | Definition |
|---|---|
| DecisionRecord | Immutable governed conclusion with explicit authority, rationale, inputs, composition, typed result, and audit linkage. Replaces the Phase 2 spine term `DecisionObject`. |
| DecisionProfile | Versioned contract defining one decision family's form, authority, inputs, typed result, rationale, obligations, projections, and limits. |
| EvaluationResult | Atomic, evaluator-specific normalized finding. Evidence for a decision, not the final decision itself. |
| DecisionRationale | Structured reasons, alternatives, obligations, warnings, and corrective suggestions embedded in a decision. |
| DecisionRequest | Bounded request to decide, with context references and constraints. |
| DecisionContext | Sanitized, policy-governed projection for AI and human explanation. Owned by FEATURE-0025. |
| AuditEvent | Immutable accountability record of actor/action/subject/outcome, scoped through `metadata.scopeRef` at one of the seven canonical `ScopeKind` values: `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, or `ServiceInstance` (ADH-2026-015; section 17 authoritative). `metadata.scopeRef` is the sole logical and serialized scope authority; the canonical `Platform` form is an absent/nil `metadata.scopeRef` (FEATURE-0012 `NormalizeScope`/`CanonicalScopeIdentity`). |
| Decision Form | Primary semantic shape of a decision: adjudication, selection, ranking, classification, resolution, allocation, plan, assessment, or recommendation. |
| Decision Authority | Whether the decision is authoritative, advisory, recommendation, or simulation. |
| Composition Strategy | Versioned, registered algorithm for combining bounded evaluation results into a decision. |

## 3. User Stories

### 3.1 Platform architect

As a platform architect, I need a common decision envelope so that all Sovrunn decision-producing domains (authorization, governance, placement, lifecycle, economics, AI) use the same immutable structure without forking it per domain.

### 3.2 Domain feature developer

As a developer of a later decision-producing feature (e.g., FEATURE-0023 PlacementDecision), I need a versioned `DecisionProfile` extension mechanism so that I can introduce typed results and domain semantics without modifying the common envelope.

### 3.3 Auditor

As an auditor, I need every meaningful decision to atomically link to an `AuditEvent` so that I can trace who decided what, when, under which authority, and with what outcome.

### 3.4 Security reviewer

As a security reviewer, I need advisory, recommendation, simulation, and AI-generated outputs to be structurally distinguishable from authoritative decisions so that they cannot be consumed as enforcement authority.

### 3.5 AI consumer

As an AI-assisted operations consumer, I need structured rationale, reason codes, obligations, and evidence references in decisions so that I can explain decisions without accessing unrestricted canonical data.

### 3.6 Sovereign deployment operator

As a sovereign deployment operator, I need decision and audit contracts that work in disconnected and air-gapped environments with local trust roots, offline profile registries, and no implicit remote control-plane dependencies.

### 3.7 Compliance officer

As a compliance officer, I need immutable decisions with correction, supersession, and revocation through linked records so that I can demonstrate audit integrity and evidence chain.

### 3.8 Downstream feature adopter

As a developer of a downstream feature, I need a lightweight adoption contract with one short section and one gate check so that adoption is consistent without bureaucratic overhead.

## 4. Acceptance Criteria

### 4.1 DecisionRecord envelope (CONTRACT_NOW)

- AC-01: A `DecisionRecord` schema defines: profile reference (name + version), form, authority, purpose, request reference, actor reference, subject references, scope identity expressed solely through `metadata.scopeRef` (see AC-01.1), effective policy context reference, input snapshot reference, evaluation result references, composition metadata, typed result, rationale, obligations, provenance, validity, correlation (operation + audit event references).
- AC-01.1 (CONTRACT_NOW — ADH-2026-016): `DecisionRecord` uses FEATURE-0012 `metadata.scopeRef` as its sole logical and serialized scope authority. No top-level `DecisionRecord` `scopeRef`, parallel scope enum, duplicated scope attribute, or any second scope source exists; a top-level or parallel scope source is prohibited and fails validation. The canonical `Platform` form is an absent/nil `metadata.scopeRef` where the applicable contract permits `Platform` (inherited from FEATURE-0012 `NormalizeScope`/`CanonicalScopeIdentity`); every non-Platform scope carries a non-nil `metadata.scopeRef` with a canonical `ScopeKind` and UID. Duplicate, conflicting, alternate, aliased, or parallel scope sources fail validation. A canonical nil `metadata.scopeRef` resolving to `Platform` is not a contradictory or duplicate scope. The seven-value `ScopeKind` vocabulary (ADH-2026-015; SV-AC-01) is preserved. See section 18.
- AC-02: The envelope is stable. New decision families extend through `DecisionProfile` and typed result schemas without adding fields to or forking the common envelope.
- AC-03: Every successfully persisted `DecisionRecord` is `FINAL` and immutable. Pending, running, and waiting states belong exclusively to `Operation`.

### 4.2 Decision forms and authority (CONTRACT_NOW)

- AC-04: The closed primary-form vocabulary is defined: `ADJUDICATION`, `SELECTION`, `RANKING`, `CLASSIFICATION`, `RESOLUTION`, `ALLOCATION`, `PLAN`, `ASSESSMENT`, `RECOMMENDATION`.
- AC-05: Authority is orthogonal to form: `AUTHORITATIVE`, `ADVISORY`, `RECOMMENDATION`, `SIMULATION`.
- AC-06: Only profiles with an adjudication facet use the closed adjudication vocabulary: `ALLOWED`, `DENIED`, `REQUIRES_APPROVAL`.
- AC-07: AI output is never `AUTHORITATIVE`. Advisory, recommendation, and simulation records are structurally distinguishable and cannot be consumed as enforcement authority.

### 4.3 DecisionProfile extension mechanism (CONTRACT_NOW)

- AC-08: A `DecisionProfile` schema declares: stable name, family, semantic version, owner, lifecycle status, primary form, allowed secondary facets, allowed authority levels, input/result/rationale/obligation schemas, accepted evaluation types, composition strategies, scope/actor/subject/policy/evidence/audit semantics, determinism/failure/timeout behavior, validity/expiry/correction/supersession/revocation rules, sensitivity/residency/retention/projection profiles, idempotency identity, bounded limits, and compatibility policy.
- AC-09: A new decision family can be registered by providing a conforming `DecisionProfile` and typed result schema without requiring a change to the common `DecisionRecord` envelope.
- AC-10: Profile registration requires conformance review. Accepting arbitrary profile URIs at runtime is prohibited.

### 4.4 EvaluationResult contract (CONTRACT_NOW)

- AC-11: An `EvaluationResult` schema defines: registered evaluator type, bounded input snapshot or digest, result, timing envelope, and explicit success/non-match/indeterminate/timeout/error semantics.
- AC-12: Evaluator-native payloads remain behind adapters. Normalized fields enter the common record.
- AC-13: Nondeterministic or AI evaluator output is captured as a governed, authorized projection conforming to the following rules:
  - AC-13.1: Captured fields are limited to: evaluator identity/version, governed input snapshot or digest (never raw prompt or unrestricted content), authorized output projection, execution configuration reference, evaluation time, trust boundary, safety/policy filter declarations, sensitivity classification, and integrity digest.
  - AC-13.2: "Authorized output projection" means a schema-valid, size-bounded, sensitivity-classified subset of evaluator output that the applicable `DecisionProfile` explicitly declares as retainable for the decision's purpose, scope, and consumer audience. The projection schema defines maximum field count, maximum total bytes, allowed value types, and prohibited content categories.
  - AC-13.3: Content classified above the applicable profile's declared sensitivity ceiling (per the closed ordered vocabulary `PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED`; see CB-AC-04 through CB-AC-07), content declared in a prohibited content category, or unrestricted-tenant-owned content must not be copied into the canonical `EvaluationResult` record. This is enforced by deterministic structural and typed validation of the declared sensitivity classification and content-category declarations, not by semantic content inspection. Such content is represented by: (a) a controlled reference (stable URI or identifier resolvable within the evaluator's trust boundary), and (b) an externally supplied integrity-digest carrier — opaque structural metadata associated with the content and carried without disclosing it. FEATURE-0013 computes no digest over content and makes no cryptographic-validity claim about this carrier; whether the carrier cryptographically proves that the content existed and was considered is DEFERRED to ADR-F13-002 (see the section 18.5.1 pre-ADR-F13-002 digest and trust conformance boundary).
  - AC-13.4: Raw prompts, model system instructions, full model responses exceeding the projection bound, secrets, credentials, provider tokens, and unauthorized tenant content must never appear in the canonical record. FEATURE-0013 enforces this structurally through typed content-category declarations, prohibited-field placement, declared sensitivity, and bounded schema rules — not through lexical or semantic content scanning. Producers remain responsible for redaction before submission; redaction is mandatory before persistence.
  - AC-13.5: Cross-scope safety validation: before an `EvaluationResult` is accepted into a `DecisionRecord`, the projection is validated deterministically against the target scope's declared sensitivity ceiling using declared scope metadata (`metadata.scopeRef`), typed references, and content-category declarations — not semantic inspection of arbitrary tenant content. An evaluation produced in Tenant A's scope must not carry typed references or content-category declarations that target or disclose Tenant B's existence, data, or resource identity relative to the projection's own `metadata.scopeRef`.
  - AC-13.6: Replay must not rerun the evaluator. The captured projection, digest, and timing envelope are sufficient for audit, explanation, and deterministic composition without re-execution.
  - AC-13.7: Conformance fixtures must include positive and negative cases proving the following. Every rejection in these cases is a deterministic structural and typed conformance check against typed content-category declarations, structural field placement, declared sensitivity, scope metadata, and bounded schema rules. No lexical secret-pattern or credential-pattern detection, and no semantic inspection of arbitrary content, is implemented or required. Fixtures must not contain or implement secret-pattern or credential-pattern detection (see SEC-15, CB-AC-08, CB-AC-09, and section 18.4).
    - (a) A conforming bounded projection within size/field/type limits is accepted;
    - (b) A projection exceeding the size bound is rejected;
    - (c) A projection is rejected on a deterministic typed/structural violation — for example, a projection declaring a prohibited content category such as `SECRET` or `CREDENTIAL`, or placing content in a field prohibited by the applicable profile/schema. Rejection is based on the typed content-category declaration and structural field placement, not on lexical secret-pattern or credential-pattern detection, which is neither implemented nor required;
    - (d) A projection is rejected when it declares a prohibited content category or uses a prohibited structural field representing `RAW_PROMPT` or `MODEL_SYSTEM_INSTRUCTION`. Rejection is based on the typed category declaration and structural field placement, not on semantic inspection of arbitrary text;
    - (e) A projection is rejected on a deterministic cross-scope violation, based on declared scope metadata, typed references, content-category declarations, and target scope relative to the projection's own `metadata.scopeRef`. Rejection is not based on semantic inspection of arbitrary tenant content;
    - (f) A protected-content reference accompanied by a present, structurally well-formed integrity-digest carrier is accepted in place of the raw content. Acceptance is a structural carrier-presence check, not a cryptographic-validity determination (see section 18.5.1);
    - (g) A protected-content reference is rejected when its integrity-digest carrier is missing, or when an externally supplied opaque observed digest state does not match an externally supplied opaque expected digest state. This mismatch check is a deterministic opaque expected-versus-observed comparison, not content-to-digest computation or cryptographic verification (see section 18.5.1);
    - (h) An evaluation result missing sensitivity classification is rejected.

**AC-13 data-minimization conformance relationship:** AC-13 is the specific application of the general data-minimization principles stated in AC-42, AC-43, AC-44, SEC-07, SEC-08, and SEC-09 to the `EvaluationResult` domain. Any conflict between AC-13 sub-criteria and those general principles must be resolved in favor of the more restrictive (data-minimizing) interpretation. AC-13 must not be read as permitting retention of content that AC-43 or SEC-07 prohibits.

### 4.5 Composition and dependency graph (CONTRACT_NOW)

- AC-14: A composite decision combines bounded evaluation results using a versioned, registered composition strategy. The strategy declares: stable identifier/version, accepted input types, ordering/precedence, short-circuit behavior, missing/timeout/conflict/error behavior, fail posture, deterministic output mapping, resource budgets, maximum fan-out, and explanation rules.
- AC-15: A bounded, directed acyclic decision-requirements graph is defined with: stable node/edge identifiers, typed nodes (input, evaluation, composition, decision), single root decision, deterministic topological ordering, declared parallelizable groups, and explicit maximum bounds (nodes, depth, width, fan-out, payload, wall-clock, remote calls, authorities, bytes, concurrent evaluations, per-evaluator timeout, retries, cross-jurisdiction calls, cost units).
- AC-16: The graph describes dependency and composition only. It cannot execute provisioning, human tasks, compensations, arbitrary code, or unbounded loops.
- AC-17: Hierarchy is resolved once into `EffectivePolicyContext` before composition. The decision layer must not independently traverse Organization/OU/Tenant/Project/Provider trees.

### 4.6 Immutability, correction, and effective state (CONTRACT_NOW)

- AC-18: Correction, supersession, and revocation create new final records with typed relationships (`CORRECTS`, `SUPERSEDES`, `REVOKES`) referencing predecessors. Predecessors are never mutated.
- AC-19: Relationship chains are bounded, cycle-free, and deterministically resolved. A read model may calculate `effectiveState: EFFECTIVE | SUPERSEDED | REVOKED` as projection metadata.
- AC-20: Retention or lawful erasure may remove protected payloads only through an approved cryptographic erasure/redaction process preserving a verifiable tombstone and chain integrity.

### 4.7 AuditEvent linkage (CONTRACT_NOW)

- AC-21: Every meaningful decision links to at least one `AuditEvent`. A single audit event may reference a decision, operation, actor, subject, resource, request, trace, and parent event. A "meaningful decision" is one that produces a durable `DecisionRecord`; ephemeral governed authorization evaluations that satisfy AC-40's aggregation profile are not individually persisted as `DecisionRecord` instances and therefore do not individually satisfy AC-21. Their audit evidence is governed by the aggregation profile's audit obligation, not by per-event `DecisionRecord` linkage. The distinction is: AC-21 applies to every persisted `DecisionRecord`; AC-40 defines which authorization events require durable `DecisionRecord` persistence versus governed aggregation.
- AC-22: A final decision must not be reported as durably recorded until its immutable `DecisionRecord` and its required audit obligation are atomically accepted by the same local persistence boundary.
- AC-23: The atomic boundary preallocates the immutable `AuditEvent` identity and stores that identity in both the `DecisionRecord` reference and audit obligation. Asynchronous materialization uses the same identity without mutating the `DecisionRecord`.
- AC-24: `AuditEvent` supports all seven canonical `ScopeKind` values through `metadata.scopeRef`: `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, and `ServiceInstance` (ADH-2026-015). `metadata.scopeRef` is the sole logical and serialized scope authority; no parallel `AuditEvent` scope enum, presence flag, `AuditScope` type, discriminator, alias, or independently mutable scope field is permitted. Per FEATURE-0012 `NormalizeScope` (`internal/apimeta/scope.go:84-92`) and `CanonicalScopeIdentity` (`internal/apimeta/scope.go:107-112`), the canonical form of `Platform` scope is an absent/nil `metadata.scopeRef`: an absent `scopeRef` resolves deterministically to `Platform` where the contract permits `Platform`, and returns the stable required-scope error where it does not (absence is not automatically valid). Every non-Platform scope requires a non-nil `metadata.scopeRef` with canonical kind and UID semantics inherited from FEATURE-0012. SV-AC-01 through SV-AC-09 (section 17) expand on and are authoritative for this criterion.
- AC-25: Audit outcomes are separate from decision outcomes: `Succeeded`, `Denied`, `Failed`. Successfully producing a `DENIED` adjudication is audit outcome `Succeeded`.

### 4.8 Execution model invariants (INVARIANT_FOR_LATER)

- AC-26: Simple and bounded decisions complete synchronously when all required inputs are available within the declared latency budget, without requiring a queue or durable workflow engine.
- AC-27: When evaluation requires human approval, unavailable evidence, or budget-exceeding work, an `Operation` owns lifecycle, retries, deadlines, and cancellation. The `DecisionRecord` links to the `Operation` but never becomes a mutable workflow-state object.
- AC-28: Orchestration is required for authoritative evaluation ordering, composition, retry, deadlines, approval, and finalization. Choreography is allowed after durable finalization for non-authoritative reactions only.

### 4.9 Statelessness and scalability invariants (INVARIANT_FOR_LATER)

- AC-29: The decision kernel is a pure, deterministic function of versioned inputs, registered strategy, and bounded configuration. Workers are interchangeable and hold no authoritative local session state.
- AC-30: Effective-once behavior uses caller-supplied idempotency key plus separate semantic decision digest. Compare-and-create or equivalent concurrency control prevents duplicate authoritative decisions.
- AC-31: Data-plane locality: all profile, schema, strategy, policy-context, trust, and evaluator material needed for a bounded synchronous decision must be locally available, version-pinned, and digest-verifiable.

### 4.10 Sovereign deployment requirements (INVARIANT_FOR_LATER)

- AC-32: Sovereignty is treated as seven enforceable dimensions: data, operational, technical, legal, cryptographic, software supply chain, and AI/model.
- AC-33: The same canonical contracts must operate in: connected public/sovereign cloud, private/Government Community Cloud, on-premises, intermittently connected edge, and fully disconnected/air-gapped deployments.
- AC-34: No canonical object contains a cloud provider, infrastructure platform, container orchestrator, AI service, product identifier, or provider-native resource type as a required semantic.
- AC-35: Disconnected synchronization uses signed origin-aware manifests, collision-resistant identifiers, replay protection, duplicate detection, per-origin ordering, and explicit trust-root version. Unknown/revoked trust roots are quarantined. Immutable records are never merged using last-writer-wins.

### 4.11 Provider neutrality (CONTRACT_NOW)

- AC-36: Provider-specific topology, capacity, jurisdiction, trust, and capability information enters only through adapters, normalized capability profiles, and typed references.
- AC-37: No provider, runtime, database, queue, workflow engine, message broker, identity provider, policy engine, AI service, or cryptographic service dependency exists in the core contract packages.

### 4.12 Authorization as a profile (CONTRACT_NOW)

- AC-38: `AuthorizationDecision` is a profile of `DecisionRecord` requiring an `ADJUDICATION` facet. It captures principal, delegation, action, resource, scope, assurance, policy context, evaluations, rationale, validity, and enforceable obligations.
- AC-39: Authentication, identity resolution, policy authoring, policy evaluation, and enforcement remain outside FEATURE-0013. An enforcement point must enforce all mandatory obligations or deny the action.
- AC-40: Denials, privileged access, cross-scope/cross-tenant access, break-glass access, data disclosure/export, policy/privilege changes, human approval, and regulated actions always require a durable `DecisionRecord` and audit obligation. High-volume low-risk checks may use a governed aggregation profile that satisfies the following boundary:
  - An aggregated low-risk authorization event is an ephemeral enforcement evaluation, not a persisted `DecisionRecord`. It is not individually linked to an `AuditEvent` per AC-21.
  - The aggregation profile must declare: which check categories qualify as low-risk, the maximum aggregation window, the minimum audit-evidence content preserved per window, the retention and reconciliation obligations, and the escalation criteria that force a check to become a durable `DecisionRecord`.
  - Aggregated audit evidence must remain queryable and attributable within the profile's declared retention. It must not suppress evidence required by any applicable security, compliance, or regulatory profile.
  - The profile-defined criteria that distinguish an aggregated low-risk authorization event from a meaningful durable decision are: (a) the action is classified as low-risk by the applicable `DecisionProfile`; (b) the outcome is `ALLOWED` with no mandatory obligations beyond standard enforcement; (c) no denial, privilege escalation, cross-scope access, break-glass, data movement, policy change, or regulated-action flag is present; (d) no profile-specific escalation trigger is active.
  - Aggregation must not reduce the total audit evidence below what the applicable profile declares as minimum. Aggregation reduces `DecisionRecord` volume but must not suppress audit evidence.

### 4.13 Security and trust (mixed classification — see per-criterion labels and section 11.2/11.3)

- AC-41 (CONTRACT_NOW — contract data only): The FEATURE-0013 contract defines the schema fields, validation rules, projection restrictions, and conformance data required to *represent* authenticated and authorized principals and hops within a `DecisionRecord`/`AuthorizationDecision`: actor references, delegation/impersonation fields, scope/purpose/action/resource binding fields, assurance-level fields, and their validation and projection restrictions. FEATURE-0013 owns only these fields, their validation, and their conformance fixtures.
- AC-41.1 (INVARIANT_FOR_LATER): The actual runtime enforcement — authenticating and authorizing every human, workload, evaluator, and service hop, and binding authorization to scope, purpose, action, and resource at a zero-trust boundary — is runtime behavior owned by the identity/enforcement features outside FEATURE-0013. It generates no FEATURE-0013 runtime task; FEATURE-0013 provides only the data contract that those enforcers populate and read.
- AC-42: Evaluator output, reason text, external evidence, and AI content are treated as untrusted data. Schema, size, depth, character, URI, and reference-kind limits are enforced.
- AC-43: Canonical governance records must not duplicate personal data, secrets, credentials, provider tokens, raw policy documents, model prompts, or unrestricted tenant content. Records prefer stable references, digests, classifications, and purpose-bound projections.
- AC-44: Different authorized projections serve customer, operator, auditor, regulator, provider, and AI consumers. Redaction is field- and reason-aware.
- AC-45 (CONTRACT_NOW — algorithm-agile carrier and structural trust metadata only; ADH-2026-016): The contract supports canonical serialization and content digests by carrying the canonicalization-method identifier, digest-algorithm identifier, covered-field descriptor, signature-algorithm identifier, and structural trust-state fields as versioned data fields; it remains algorithm-agile and hardcodes no specific algorithm. These algorithm-agile carrier fields and structural trust metadata are `CONTRACT_NOW`. Structural trust metadata is separated from cryptographic verification: structural fixtures may use opaque deterministic assertions to validate state and failure semantics but must not claim cryptographic validity or select an algorithm (see the section 18.5.1 pre-ADR-F13-002 digest and trust conformance boundary). When a profile requires verified trust, an unknown, absent, expired, revoked, mismatched, or unverified trust state fails closed. Deployment profiles may additionally require signatures, trusted timestamps, append-only/WORM retention, or external notarization; those assurance services are DEFERRED per section 11.2. RFC 8785 remains illustrative only and is not selected.
- AC-45.1 (DEFERRED; escalation marker: ARCHITECTURE_DECISION_REQUIRED — ADR-F13-002; preserved and reaffirmed by ADH-2026-016): The selection of the specific canonicalization algorithm/profile (for example RFC 8785 JCS or an alternative), the exact set of fields covered by the content digest, the signature algorithm, and cryptographic product/service selection are not supplied by any approved architecture source. ADH-2026-014 approves only algorithm-agility (AD item 16) and explicitly prohibits product selection; ADH-2026-016 confirms these remain DEFERRED and blocked by ADR-F13-002. Design, tasks, fixtures, and code must not select a canonicalization algorithm, a covered-field set, a signature algorithm, or a cryptographic product/service until ADR-F13-002 is resolved through the section 11.8 escalation process. Actual cryptographic verification, signing, HSM, notary, trusted-time, key, revocation-distribution, and WORM services are not FEATURE-0013 artifacts. Documentation and traceability may carry a clearly labelled placeholder such as `PLACEHOLDER_PENDING_ADR_F13_002`; executable conformance fixtures must not select or encode the unresolved algorithm, covered-field set, or signature algorithm. ADR-F13-002 remains a separate `ARCHITECTURE_DECISION_REQUIRED` escalation marker in the section 11.9 register.

### 4.14 AI governance invariants (mixed classification — AC-46 through AC-48 and AC-49.1 INVARIANT_FOR_LATER; AC-49.2 CONTRACT_NOW)

AC-46 through AC-49 constrain the behavior of a later AI consumer and AI runtime. The `DecisionContext` schema, the AI consumer, prompt/model integration, the model client, the AI runtime, and any production projection service are owned by FEATURE-0025 or the later owning feature. FEATURE-0013 does **not** define or implement them. FEATURE-0013 owns only the `DecisionProfile` projection constraints, authorized-projection semantics, references, labels, and non-executing conformance data that `DecisionRecord` and `AuditEvent` require; those genuinely FEATURE-0013-owned contract elements remain CONTRACT_NOW and are carried by AC-08 (DecisionProfile schema) and the projection semantics in SEC-09 and SEC-13. AC-46 through AC-48, and the runtime-guarantee portion AC-49.1, generate no FEATURE-0013 runtime or schema implementation task beyond documentation and non-executing conformance placeholders. The single, narrowly scoped exception is AC-49.2 (CONTRACT_NOW): a pure/in-memory conformance fixture may execute to prove that FEATURE-0013-owned `DecisionRecord` validation and deterministic composition have no AI dependency and accept the contractually defined empty/unavailable/absent AI projection state. That fixture must not implement `DecisionContext`, a model client, prompt processing, a projection service, or any AI runtime behavior.

- AC-46: (INVARIANT_FOR_LATER) A later AI consumer receives only an authorized `DecisionContext` projection. It may explain, summarize, identify anomalies, and suggest remediation through normal approval. The authorized-projection constraint fields that a `DecisionProfile` declares are CONTRACT_NOW (AC-08); the `DecisionContext` schema and AI consumer are owned by FEATURE-0025.
- AC-47: (INVARIANT_FOR_LATER) A later AI consumer must not mutate records, invent evidence, suppress audit, widen disclosure, bypass approval, or execute a suggestion. This constrains later AI-runtime behavior, not a FEATURE-0013 artifact.
- AC-48: (INVARIANT_FOR_LATER) AI-generated recommendations are labelled with model/provider/version, prompt-template version, data boundary, confidence, and human-review status. The label/field schema a `DecisionProfile` or `DecisionRecord` may carry is CONTRACT_NOW (AC-08); producing the labelled recommendation is later AI-runtime behavior owned by FEATURE-0025.
- AC-49: (SPLIT CLASSIFICATION)
  - AC-49.1 (INVARIANT_FOR_LATER): The runtime guarantee that a later production deployment's core decisions remain available and correct when the AI subsystem is disabled or disconnected is a runtime invariant that the later owning implementation (FEATURE-0025 or later) must satisfy. FEATURE-0013 does not implement AI runtime, `DecisionContext`, model clients, or prompt processing to prove this runtime guarantee.
  - AC-49.2 (CONTRACT_NOW): The contract-level AI-independence property that FEATURE-0013 owns is proven now with a narrowly scoped, pure/in-memory conformance fixture (EC-09) that demonstrates, without any AI runtime: (a) `DecisionRecord` schema validation succeeds with no AI subsystem present; (b) deterministic composition over captured, bounded evaluation results has no dependency on any AI evaluator, model client, or projection service; and (c) the contractually defined AI projection state resolves to the empty/unavailable/absent projection value defined by the contract rather than to an error. This fixture must not implement or require `DecisionContext`, a model client, prompt processing, a projection service, or any AI runtime behavior; it validates only that the FEATURE-0013-owned schemas and deterministic composition accept the absent/unavailable projection state.

### 4.15 Conformance fixtures (CONTRACT_NOW)

- AC-50 (CONTRACT_NOW — ADH-2026-016): Executable positive and negative fixtures are defined for the conformance scenarios of canonical architecture section 17, which map exactly and one-to-one to the stable IDs `F13-CF-01` through `F13-CF-28` in their existing order. The mapping must not be merged, renumbered, or omitted. The canonical scenario-to-ID mapping is:

  | Conformance ID | Scenario (architecture section 17) |
  |---|---|
  | F13-CF-01 | simplest synchronous atomic authorization decision |
  | F13-CF-02 | simplest deterministic denial with actionable rationale and obligations |
  | F13-CF-03 | pure selection with a typed result and no adjudication facet |
  | F13-CF-04 | ranking, classification, resolution, allocation, plan, assessment, recommendation, advisory, and simulation profiles |
  | F13-CF-05 | human-approval asynchronous decision linked to an `Operation` |
  | F13-CF-06 | maximally complex but bounded hierarchical composite decision |
  | F13-CF-07 | parallel evaluations with deterministic aggregation |
  | F13-CF-08 | timeout, missing evidence, evaluator failure, and policy conflict |
  | F13-CF-09 | idempotent replay and concurrent duplicate requests |
  | F13-CF-10 | partial decision/audit persistence failure and reconciliation |
  | F13-CF-11 | supersession, correction, revocation, cycle, and chain-limit behavior |
  | F13-CF-12 | unauthorized profile, authority, projection, obligation bypass, and cross-scope reference rejection |
  | F13-CF-13 | offline/air-gapped profile registry, trust, and signed-bundle operation |
  | F13-CF-14 | export/import across two conforming provider implementations |
  | F13-CF-15 | AI-disabled operation and AI projection redaction |
  | F13-CF-16 | graph size, depth, width, fan-out, timeout, and payload budget rejection |
  | F13-CF-17 | registration of a new decision family without a common-envelope change |
  | F13-CF-18 | immutable supersession and revocation with calculated effective state |
  | F13-CF-19 | local atomic decision/audit-obligation acceptance while downstream audit delivery is unavailable |
  | F13-CF-20 | idempotency retry identity versus semantic decision identity across policy, profile, strategy, authority, evaluator-version, and validity changes |
  | F13-CF-21 | deterministic replay from captured nondeterministic or AI evaluation output without rerunning the evaluator |
  | F13-CF-22 | synchronous operation with profile registry and policy resolver unavailable but valid local digest-pinned bundles present |
  | F13-CF-23 | rejection of unknown, expired, revoked, or digest-mismatched profile bundles |
  | F13-CF-24 | enforcement denial for unknown or unsupported mandatory obligations |
  | F13-CF-25 | local authorization artifact cache hit, expiry, revocation, scope mismatch, and stale-policy failure |
  | F13-CF-26 | cancellation and total remote-call, byte, concurrency, retry, jurisdiction, time, and cost budgets for composite graphs |
  | F13-CF-27 | local digest plus batched signing/notarization during remote cryptographic-service outage |
  | F13-CF-28 | disconnected export/import replay, duplicate, ordering, unknown trust root, revoked origin, and explicit conflict reconciliation |

  In addition to `F13-CF-01` through `F13-CF-28`, the following case families use separately named IDs and do **not** alter the canonical 28-scenario count:
  - **AC-13.7 evaluator-output governed-projection cases** (`F13-CF-EVAL-a` through `F13-CF-EVAL-h`, corresponding to AC-13.7 (a)–(h));
  - **Scope cases** (`SV-AC-02` through `SV-AC-09` scope-vocabulary fixtures, including the seven positive `AuditEvent` scope fixtures, the canonical nil-`Platform` fixture, the negative absence-when-`Platform`-disallowed fixture, the second-scope-source negative fixture, and the out-of-vocabulary negative fixture);
  - **Security cases** (structural conformance fixtures for AC-42, AC-43, SEC-05, SEC-06, SEC-07, and the fail-closed high-risk-profile fixtures per SEC-12, bounded to structural/typed conformance per section 18.4);
  - **Compatibility cases** (AC-54 / BC-AC-04: seven-value `AuditEvent` compatibility fixtures and the six-value FEATURE-0012 regression fixtures).

  Coverage is counted by scenario ID, not by fixture-file count. One fixture may cover multiple scenario IDs, but **every** scenario ID (both `F13-CF-01` through `F13-CF-28` and each separately named case above) requires its own explicit coverage-matrix entry. Contract-only fixtures remain pure/in-memory and must not implement prohibited runtime systems.

**AC-50 implementation scope constraint:** The "authorization cache" scenario fixture tests the contract semantics of a cached authorization decision (validity, expiry, staleness detection, refresh obligation, scope binding) using pure-function or in-memory conformance fixtures only. It must not generate implementation tasks for a runtime authorization cache, cache-invalidation service, distributed cache, TTL worker, persistence layer, or any background process. The fixture validates that a conforming cache implementation would satisfy the contract, without itself being or requiring that cache.

**AC-50 evaluator-output governed-projection fixture constraint:** The evaluator-output scenarios required by AC-13.7 test the contract semantics of governed capture (projection bounds, redaction, sensitivity classification, cross-scope safety, reference/digest handling) using pure-function or in-memory conformance fixtures only. They must not generate implementation tasks for an evaluator runtime, AI model client, prompt engine, redaction service, or any background process. The fixtures validate that a conforming evaluator adapter and record-acceptance gate would satisfy the AC-13 data-minimization contract, without themselves being or requiring those services.

**AC-50 AI-disabled fixture constraint:** The "AI-disabled" scenario fixture tests the contract semantics of AI-independent operation (per AC-49.2 and EC-09) using pure-function or in-memory conformance fixtures only. It proves that `DecisionRecord` validation and deterministic composition succeed with no AI subsystem present and that the contractually defined empty/unavailable/absent AI projection state is accepted rather than raising an error. It must not generate implementation tasks for `DecisionContext`, an AI model client, a prompt engine, a projection service, or any AI runtime or background process. The fixture validates that a conforming implementation would satisfy the AI-disabled contract, without itself being or requiring that AI runtime.

### 4.16 Compatibility and evolution (CONTRACT_NOW)

- AC-51: Every canonical schema has explicit version and compatibility rules.
- AC-52: Decision form, authority, adjudication, typed-result, reason-code, obligation, strategy, and digest semantics cannot change silently. Readers reject unsupported major versions.
- AC-53: Registries for decision profile, form, authority, evaluator type, strategy, reason code, obligation, event type, evidence type, and projection are versioned and governable offline.
- AC-54: Per ADH-2026-015, `AuditEvent` expansion requires compatibility fixtures for all seven canonical `ScopeKind` values (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance`), separate regression fixtures for the six pre-existing FEATURE-0012 `ScopeKind` values (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`), and a documented alpha compatibility decision. The `Platform` positive fixture uses the canonical absent/nil `metadata.scopeRef` form; a negative fixture must assert that an absent `metadata.scopeRef` is rejected with the stable required-scope error under a contract that does not permit `Platform`.

### 4.17 Downstream adoption contract (CONTRACT_NOW)

- AC-55: One normative downstream adoption contract is defined in this feature.
- AC-56: Each applicable downstream feature includes one short `FEATURE-0013 Adoption` section declaring applicability, roles, profiles, reused contracts, inherited risks, exceptions, and conformance evidence.
- AC-57: One lightweight feature-gate check validates adoption section presence, completeness, and vocabulary compliance.
- AC-58: No separate manifests, registry, index, or approval workflow is introduced for adoption. Structured manifests remain trigger-based per section 28.8 of the canonical architecture.

### 4.18 Observability (INVARIANT_FOR_LATER)

- AC-59: Metrics and traces expose latency, queueing, evaluation fan-out, cache result, composition strategy, failure class, retry, saturation, and audit-reconciliation health without exposing decision evidence or personal data as labels.
- AC-60: Logs are structured, redacted, locally routable, retention-controlled, and correlated by opaque identifiers.
- AC-61: Health distinguishes evaluator, evidence, store, audit, identity, key, time, and synchronization dependencies.

## 5. Non-Goals

FEATURE-0013 must not implement:

- NG-01: A production decision, authorization, policy, placement, audit, evidence, identity, profile-registry, synchronization, cryptographic, or AI service.
- NG-02: A database schema, migration, repository implementation, transaction/outbox, event store, WORM store, cache, distributed lock, or lease.
- NG-03: A queue, broker, workflow/durable-execution runtime, scheduler, controller, worker pool, background reconciler, retry daemon, or dead-letter processor.
- NG-04: A network registry client, remote policy resolver, external evidence fetcher, HSM/notary client, provider adapter, model client, vector store, or agent.
- NG-05: A general-purpose rules engine, policy language, workflow DSL, DAG execution engine, expression interpreter, dynamic plugin system, or code generator.
- NG-06: Production authorization tokens, cache invalidation, revocation distribution, multi-site consensus, disconnected synchronization, or conflict-resolution service.
- NG-07: Provider-, runtime-, database-, broker-, identity-, cryptographic-, or model-specific dependencies in core contract packages.
- NG-08: Selection of any specific workflow engine, message broker, database, identity provider, policy engine, AI model, or cryptographic service.
- NG-09: India government security-classification vocabulary, exact retention schedules, or jurisdiction-specific legal values (these belong in versioned profiles).
- NG-10: Real PostgreSQL provisioning, real plugin execution, full OPA/Cedar integration, full Keycloak/Vault/Temporal integration.
- NG-11: Production persistence, billing, autonomous AI operations, failover/DR execution, or production compliance engine.

## 6. Edge Cases

- EC-01: A decision profile with no adjudication facet must not require `ALLOWED`/`DENIED`/`REQUIRES_APPROVAL` outcomes. Pure selection, ranking, classification, resolution, allocation, plan, assessment, and recommendation carry answers in the typed result only.
- EC-02: Runtime failure, timeout, evaluator error, and insufficient evidence are evaluation/production states, not business adjudications. A registered fail-safe strategy may map them to denial or required approval and must record that mapping.
- EC-03: Unknown, revoked, expired, or digest-mismatched profiles are rejected or quarantined—they never silently proceed.
- EC-04: Concurrent duplicate requests with the same idempotency key must produce the same outcome. Requests with the same purpose but different policy context, profile version, or authority produce different semantic identities and are not duplicates.
- EC-05: Partial decision/audit persistence failure: the synchronous response must not fabricate success. Recovery, bounded delivery, poison-record quarantine, and reconciliation must be defined.
- EC-06: A supersession/correction/revocation chain that exceeds the bounded depth limit must be rejected with a clear error.
- EC-07: A composite graph exceeding maximum nodes, depth, width, fan-out, payload, or wall-clock budget must be rejected before execution begins.
- EC-08: An enforcement point receiving an unknown or unsupported mandatory obligation must deny the action, even if the adjudication was `ALLOWED`.
- EC-09: AI-disabled deployments must function correctly for all decision and audit behavior; AI projections return the contractually defined empty/unavailable/absent value rather than an error. FEATURE-0013 proves only the contract portion of this edge case (per AC-49.2, CONTRACT_NOW) with a pure/in-memory conformance fixture that demonstrates `DecisionRecord` validation and deterministic composition have no AI dependency and accept the empty/unavailable/absent projection state. The fixture must not implement or require `DecisionContext`, a model client, prompt processing, a projection service, or any AI runtime. The broader runtime availability guarantee for a later production deployment is AC-49.1 (INVARIANT_FOR_LATER) and generates no FEATURE-0013 runtime task.
- EC-10: Cross-scope references that would disclose existence of resources in another tenant or governance domain must be safe-denied without existence disclosure.

## 7. Security and Privacy Requirements

- SEC-01 (CONTRACT_NOW — contract data only): FEATURE-0013 defines the schema fields, validation rules, and conformance data required to represent authenticated and authorized actors and hops (actor references, delegation, assurance, and scope/purpose/action/resource binding fields) and their projection restrictions. FEATURE-0013 owns only these fields, their validation, and their conformance data.
- SEC-01.1 (INVARIANT_FOR_LATER): Runtime authentication and authorization of every human, workload, evaluator, and service hop at a zero-trust boundary is enforcement behavior owned by the identity/enforcement features outside FEATURE-0013 and generates no FEATURE-0013 runtime task.
- SEC-02: Bind authorization to scope, purpose, action, and resource—not network location alone.
- SEC-03: Record delegation and impersonation explicitly in actor references.
- SEC-04: Reject unsigned or untrusted evaluator/strategy definitions where signing is required by the deployment profile.
- SEC-05: Treat evaluator output, reason text, external evidence, and AI content as untrusted data requiring validation.
- SEC-06: Enforce schema, size, depth, character, URI, and reference-kind limits on all inputs to prevent injection and resource exhaustion.
- SEC-07: The canonical governance record must not indiscriminately duplicate personal data, secrets, credentials, provider tokens, raw policy documents, model prompts, or unrestricted tenant content.
- SEC-08: Records prefer stable references, digests, classifications, and purpose-bound projections.
- SEC-09: Confidentiality: different authorized projections serve different consumers. Redaction is field- and reason-aware and must not change the final outcome or falsely imply omitted evidence did not exist.
- SEC-10: Integrity: the contract carries canonical-serialization and content-digest metadata as algorithm-agile carrier fields (see AC-45, CB-AC-10, and the section 18.5.1 boundary). FEATURE-0013 does not execute canonical serialization and computes no digest over content; canonicalization-algorithm, digest-covered-field, and signature-algorithm selection are DEFERRED to ADR-F13-002. Deployment profiles may require additional assurance (signatures, timestamps, WORM, Merkle proofs, notarization); those assurance services are DEFERRED per section 11.2.
- SEC-11: Tiered assurance is by profile. A profile may require that an externally supplied digest carrier accompany a record rather than mandating per-record remote HSM or notarization on the hot path. Actual local digest computation over content is not performed by FEATURE-0013 and is DEFERRED to ADR-F13-002 (see section 18.5.1); FEATURE-0013 defines only the algorithm-agile carrier and structural trust-state contract. Per-record remote HSM or notarization is not a universal hot-path requirement.
- SEC-12: For security, authorization, sovereignty, compliance, privileged access, data movement, and break-glass profiles, fail-closed or `REQUIRES_APPROVAL` is mandatory. Fail-open is prohibited unless a bounded approved exception identifies scope, owner, compensating controls, expiry, audit treatment, and reassessment trigger.
- SEC-13 (CONTRACT_NOW — contract data only): FEATURE-0013 defines the data-minimization and authorized-projection restrictions that apply to any AI-related content a `DecisionProfile`, `EvaluationResult`, or `DecisionRecord` may carry — prohibited-content categories, redaction requirements before persistence, sensitivity classification, and tenant-boundary constraints on the retained projection. FEATURE-0013 owns only these schema/validation/projection restrictions and their conformance fixtures.
- SEC-13.1 (INVARIANT_FOR_LATER): The actual runtime handling of AI prompts, logs, traces, and model calls — including redaction at emission time and tenant-boundary enforcement in a running AI subsystem — is AI-runtime behavior owned by FEATURE-0025 or the later owning feature. It generates no FEATURE-0013 runtime task; FEATURE-0013 supplies only the contract those runtimes must satisfy.
- SEC-14: Cross-scope references must not disclose or influence another tenant or governance domain.
- SEC-15 (CONTRACT_NOW — structural boundary only; ADH-2026-016): FEATURE-0013 security validation is bounded to structural schema validation, typed classification, declared bounds, projection restrictions, and deterministic conformance. FEATURE-0013 does not claim comprehensive semantic detection of secrets, credentials, personal data, malware, or policy violations in arbitrary content. Producers must redact prohibited content before submission. Comprehensive semantic content scanning and runtime content inspection belong to later approved security features. No secret-pattern, credential-pattern, PII, malware, DLP, or content-scanning engine may be introduced in FEATURE-0013. Where AC-13, AC-42, AC-43, SEC-05, SEC-06, and SEC-07 describe rejecting prohibited content (for example the AC-13.7 fixtures), the rejection is structural and typed conformance against declared bounds, prohibited content-category rules, and sensitivity ceilings — not a semantic content-scanning capability. See section 18.4.

## 8. Compatibility with Completed Phase 1 Features

### 8.1 FEATURE-0005 Operation Resource

- `Operation` continues to own asynchronous lifecycle (progress, retries, deadlines, cancellation, status). `DecisionRecord` links to `Operation` via `correlation.operationRef` but never becomes a mutable workflow-state object.
- No change to the Phase 1 Operation resource shape is required. Decision-linked operations use the existing correlation and reference patterns.

### 8.2 FEATURE-0001 through FEATURE-0004 Organization hierarchy

- `DecisionRecord` scope references use the same Organization/OrganizationUnit/Tenant/Project hierarchy defined in Phase 1.
- `AuditEvent` scope expansion (ADH-2026-015) permits all seven canonical `ScopeKind` values via `metadata.scopeRef`. `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, and `Project` map to the Phase 1 hierarchy; `Provider` is retained from the FEATURE-0012 shared vocabulary; `ServiceInstance` is additive and extends naturally from FEATURE-0008.
- No Phase 1 resource shapes are modified.

### 8.3 FEATURE-0006 ServiceClass/ServicePlan and FEATURE-0008 ServiceInstance

- `DecisionRecord` may reference a `ServiceClass`, `ServicePlan`, or `ServiceInstance` as a typed subject or reference. These are typed references, not embedded duplicates. `ServicePlan` may appear only as a typed subject/reference; it is never a scope. `ServicePlan` is not a `ScopeKind` and must not be used as a `DecisionRecord` (or `AuditEvent`) scope.
- `DecisionRecord` scope uses exactly the seven canonical `ScopeKind` values approved by ADH-2026-015 (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance`), consistent with SV-AC-01 (section 17). `ServiceInstance` may be used as scope where the applicable contract permits it. No value outside the canonical seven — including `ServicePlan` — is a valid scope kind.
- ServiceInstance-scoped `AuditEvent` records use the same resource identity model from Phase 1.

### 8.4 FEATURE-0012 API/resource grammar

- FEATURE-0013 extends the FEATURE-0012 resource grammar without creating a parallel metadata, reference, error, or validation grammar.
- The seven-value `ScopeKind`/`AuditEvent` expansion (ADH-2026-015) is an approved expansion of the FEATURE-0012 alpha profile with explicit compatibility fixtures required (AC-54).
- `DecisionRecord` uses the FEATURE-0012 `apiVersion`, `kind`, `metadata`, `spec`, `status` conventions where applicable to resource-shaped contracts.

## 9. Reuse Assessment

### 9.1 Feature-level reuse summary

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| Common DecisionRecord envelope and vocabulary | Build | Sovrunn-owned differentiation; no existing standard provides the exact governed-conclusion, authority, form, typed-result, immutable, projection, sovereign-ready envelope. | Approved | ADH-2026-014, DEC-0026 |
| DecisionProfile extension mechanism | Build | Sovrunn-owned typed extensibility for decision families; must integrate with Sovrunn governance and profile semantics. | Approved | ADH-2026-014 |
| EvaluationResult normalization | Extend | Extends FEATURE-0012 grammar with evaluation-specific fields; evaluator adapters wrap external native types. | Approved | ADH-2026-014, DEC-0036 |
| Composition and bounded dependency graph concepts | Reuse | Reuses bounded decision-requirements concepts from DMN; does not select a DMN engine. | Approved | ADH-2026-014 |
| AuditEvent linkage contract | Extend | Extends the FEATURE-0012 alpha AuditEvent profile to all seven canonical `ScopeKind` values (ADH-2026-015) with atomic acceptance semantics. | Approved | ADH-2026-014, ADH-2026-015 |
| Canonical serialization and digest (algorithm-agile carrier concept) | Reuse | Reuses canonical-JSON and content-digest concepts to define an algorithm-agile carrier contract. Approval applies only to reuse of the algorithm-agile carrier concept (the canonicalization-method identifier, digest-algorithm identifier, covered-field descriptor, and signature-algorithm identifier fields per AC-45) and the structural trust metadata, which remain `CONTRACT_NOW`; it does not approve any specific algorithm, covered-field set, signature algorithm, or cryptographic product/service. The full ADR-F13-002 blocking scope — the specific canonicalization algorithm/profile, the exact digest-covered field set, the signature algorithm, and the cryptographic product/service selection — is not approved by any source and remains escalated as ADR-F13-002 (ARCHITECTURE_DECISION_REQUIRED) in the section 11.9 open-decision register; design, tasks, fixtures, and code must not select any of those four values until resolved. Actual cryptographic verification, signing, and production cryptographic integrations/services remain separately DEFERRED. | Approved | ADH-2026-014 (algorithm-agility carrier and structural trust metadata only; canonicalization algorithm/profile, digest-covered field set, signature algorithm, and cryptographic product/service selection escalated as ADR-F13-002) |
| Telemetry correlation | Reuse | Reuses OpenTelemetry correlation concepts for trace/span; never audit authority. | Approved | ADH-2026-014 |
| Event transport mapping | Reuse | Optional CloudEvents transport mapping; never the canonical decision or audit record. | Approved | ADH-2026-014 |
| Supply-chain evidence references | Reuse | Reuses DSSE/in-toto and SPDX/CycloneDX concepts for signed evidence references; does not duplicate their schemas. | Approved | ADH-2026-014 |

### 9.2 Capability assessment: Common DecisionRecord envelope and vocabulary

#### Identity

- Feature: FEATURE-0013
- Capability: Common DecisionRecord envelope, forms, authority, adjudication, typed result, rationale, obligations, immutability, correction, provenance, validity, and correlation
- Assessment owner: Architecture owner

#### Classification

- Disposition: Build
- Decision status: Approved

#### Analysis

- Assessment scope: The stable governed-decision envelope that all decision-producing Sovrunn domains share
- Candidate category: Decision/governance standards and frameworks
- Mature candidates / applicable standards: OMG DMN decision model (bounded-requirements graph concepts reused, engine not selected), CloudEvents (transport only), OpenTelemetry (correlation only), Kubernetes conditions/status (grammar reused via FEATURE-0012)
- Relevant candidate strengths: DMN provides proven graph and decision-table concepts; CloudEvents provides event-envelope semantics; OpenTelemetry provides trace correlation
- Material candidate constraints: No existing standard provides the exact combination of: immutable governed conclusion, authority orthogonal to form, typed profile extension, mandatory audit linkage, sovereignty-ready projections, bounded composition, offline-capable registries, and Sovrunn-specific governance semantics
- Rationale: The envelope is Sovrunn's primary differentiator for governed, explainable, sovereign-ready decisions. External standards are reused where they fit but none provides the complete governed envelope.
- Selected foundation or approach: Sovrunn-owned common envelope; reuse DMN concepts for bounded composition; reuse FEATURE-0012 API grammar for metadata/reference/validation patterns

#### Boundary

- Sovrunn-owned responsibility: Envelope schema, forms, authority, adjudication, typed-result extension, rationale, obligations, immutability, correction/supersession/revocation, provenance, validity, correlation, and conformance
- Reused or external responsibility: DMN composition concepts, OpenTelemetry correlation, canonical-serialization/content-digest concepts (canonicalization algorithm/profile, digest-covered field set, signature algorithm, and cryptographic product/service selection all pending ADR-F13-002), FEATURE-0012 grammar
- Data crossing the boundary: Evaluation results enter through normalized adapter contracts; external evidence enters as digest references
- Control crossing the boundary: Profile, strategy, and evaluator definitions enter through governed registration; runtime evaluator calls cross through adapter ports
- Adapter required: Yes — the adapter boundary applies, but the adapter interfaces are not owned by FEATURE-0013.
- Adapter rationale: External evaluators, policy engines, identity providers, secret providers, evidence stores, and observability/cryptographic services are accessed through adapters per DEC-0036. FEATURE-0013 owns only the normalized data that crosses that boundary — `DecisionRecord`, `DecisionProfile`, `EvaluationResult`, `AuditEvent`, their validation, and adapter-facing data requirements — not the adapter interfaces.
- Adapter or contract identifier: `EvaluatorAdapter`, `PolicyEngineAdapter`, `IdentityProviderAdapter`, `SecretProviderAdapter`, and `ObservabilityAdapter` are named as illustrative future consumers/producers only. FEATURE-0013 does not define these interfaces or any shared adapter interface; they are owned by FEATURE-0016 or the later owning feature, and no interface or implementation artifact for them is generated in FEATURE-0013. DEC-0036 boundary intent is preserved.
- Vendor-native types allowed: No

#### Suitability

- Sovereignty and deployment fit: Envelope carries sovereignty dimensions, residency, jurisdiction, and deployment-profile references. Works in connected, disconnected, and air-gapped environments with local trust and offline registries.
- Security and trust: Zero-trust, least-privilege projections, algorithm-agile integrity, tiered assurance, fail-closed for high-risk profiles
- Operational and supportability: Operators diagnose without vendor remote access; structured logs, health, and reconciliation
- Licensing and supply-chain: No external runtime dependency; concepts reused from open standards only
- Portability and provider-neutrality impact: No provider, runtime, or database dependency in core. Portable across public/private/sovereign/disconnected environments.

#### Phase and scope

- Allowed in current phase: Yes
- Current-phase work: Schema definitions, vocabulary registries, validation rules, conformance fixtures, compatibility checks, documentation, and traceability
- Deferred work: Production persistence, network registry, signing service, runtime evaluation, provider adapters, AI model integration, multi-site synchronization
- Explicit non-goals: See section 5 NG-01 through NG-11
- Exit or migration boundary: If a future standard emerges that provides the same governed-envelope semantics, Sovrunn may adopt it through ADH process. Current envelope remains stable until explicitly superseded.
- Phase 2 non-goal acknowledgement: No production runtime implementation in Phase 2.

#### Build justification

- Why Reuse is insufficient: No existing standard provides the exact combination of immutable governed conclusion, authority orthogonal to form, typed profile extension without envelope changes, mandatory audit-obligation acceptance, sovereignty-ready projections, offline profile registries, and Sovrunn governance semantics.
- Why Wrap is insufficient: There is no existing decision-record engine or service to wrap; the envelope is a data-model/contract concern, not a runtime system.
- Why Extend is insufficient: No single existing standard contains the core semantics to extend. DMN provides graph concepts (reused) but not the governed envelope, authority model, or audit linkage.
- Protected Sovrunn differentiation and long-term ownership: The governed decision envelope with explainability, sovereignty, and immutability is the primary Sovrunn differentiator for regulated and sovereign environments.

#### Risk mitigation

- Applicable architecture risks: F13-R01, F13-R02, F13-R03, F13-R05, F13-R22, F13-R27, F13-R29
- Preventive controls: Governed profile registration (F13-R01), versioned result schemas with semantic-override prohibition (F13-R02), explicit authority with projection marking (F13-R03), separate retry and semantic identity (F13-R05), contract-only Phase 2 scope (F13-R22), stable envelope with architecture change control (F13-R27), coordinated traceability update (F13-R29)
- Detection controls: Duplicate-profile fixtures, schema conformance tests, authority-confusion negative tests, policy/profile mutation tests, changed-file/dependency gates, baseline semantic-diff gate, repository-wide identity checks
- Corrective path: Reclassify overlapping profiles, enforce schema boundaries, reject confused authority, reject mismatched identity, remove runtime dependencies, restore semantic stability, apply migration aliases
- Residual risk: Medium (F13-R01, F13-R02, F13-R03, F13-R05 at Medium target; F13-R22, F13-R27, F13-R29 at Low target)
- Replacement risk: Low — envelope is Sovrunn-owned with no external runtime dependency
- Reassessment triggers: New decision profile family, stable API promotion, production runtime selection, common-envelope change proposal

#### Traceability

- Related DEC / RFC / ADH references: DEC-0026, DEC-0036, RFC-0023, ADH-2026-014
- Linked acceptance criteria: AC-01 through AC-03, AC-08 through AC-10, AC-18 through AC-20, AC-36, AC-37, AC-50 through AC-54
- Validation and review evidence: Conformance fixtures (28 scenarios), schema validation tests, negative compatibility tests, FEATURE-0012 baseline compatibility checks

#### Human-approval evidence

- Approving person: Sanjeev Kumar
- Approval date: 2026-07-27
- Approved ADH: ADH-2026-014
- Scope: Approved disposition and responsibility boundary for the common DecisionRecord envelope as Build within FEATURE-0013.

## 10. Architecture Drift Checks

| Drift check | Status | Evidence |
|---|---|---|
| No provider-specific hardcoding in core | PASS | AC-34, AC-36, AC-37; no provider identifiers in schemas or vocabularies |
| No Kubernetes-only assumptions in core | PASS | AC-34; contracts are HTTP/API-native per FEATURE-0012; no CRD dependency |
| No PostgreSQL lifecycle logic in core placement engine | PASS | Not applicable; FEATURE-0013 does not contain placement or PostgreSQL logic |
| No custom policy engine embedded in handlers | PASS | AC-17; hierarchy resolved through EffectivePolicyContext; policy engines accessed through adapter per DEC-0028 |
| No raw secret storage | PASS | AC-43, SEC-07; records use SecretRef/digest references, never raw credentials |
| No customer-facing IaaS leakage | PASS | AC-34, AC-36; provider details enter only through adapters and typed references |
| Explainable decision object | PASS | AC-01, AC-08, AC-46; structured rationale, reason codes, obligations, and projections are mandatory |
| Defined audit behavior | PASS | AC-21 through AC-25; atomic decision/audit acceptance, seven canonical `ScopeKind` scopes via `metadata.scopeRef` (ADH-2026-015), separated outcomes |
| Preserved adapter boundaries | PASS | AC-12, AC-37; evaluator-native payloads behind adapters; no external engine types in core; DEC-0036 referenced |

## 11. Scope Classification Summary

### 11.1 Classification rules

- `CONTRACT_NOW` — schema, vocabulary, static registry/bundle format, validation, compatibility rule, fixture, conformance check, or documentation implemented by FEATURE-0013.
- `INVARIANT_FOR_LATER` — normative constraint that a later implementation must satisfy. Generates no production runtime implementation tasks; only documentation, traceability, and non-executing conformance placeholders are permitted.
- `DEFERRED` — product, runtime, topology, numerical SLO, or operational choice that is not implemented and requires a later approved feature. Generates no implementation tasks except documentation and traceability.

### 11.2 Acceptance criteria classification

| Classification | Acceptance criteria |
|---|---|
| CONTRACT_NOW | AC-01 through AC-25, AC-36 through AC-40, AC-41 (contract data only), AC-42 through AC-44, AC-45 (algorithm-agile contract only), AC-49.2, AC-50 through AC-58 |
| INVARIANT_FOR_LATER | AC-26 through AC-35, AC-41.1, AC-46 through AC-48, AC-49.1, AC-59 through AC-61 |
| DEFERRED | AC-45.1 canonicalization algorithm/profile, digest-covered field set, signature algorithm, and cryptographic product/service selection (with separate `ARCHITECTURE_DECISION_REQUIRED` escalation marker ADR-F13-002); production persistence pattern; workflow/orchestration runtime; evidence store; identity provider; key manager; jurisdiction-specific values; provider adapters; AI runtime; signature services and production cryptographic integrations (signing, trusted timestamping, notarization, WORM enforcement); cross-site topology |

### 11.3 Security requirements classification

| Requirement | Classification | Rationale |
|---|---|---|
| SEC-01 | CONTRACT_NOW (contract data only) | Schema fields, validation, and conformance data representing authenticated/authorized actors and hops and their projection restrictions |
| SEC-01.1 | INVARIANT_FOR_LATER | Runtime authentication/authorization enforcement of every hop at a zero-trust boundary is owned by the identity/enforcement features; generates no FEATURE-0013 runtime task |
| SEC-02 | CONTRACT_NOW | Authorization binding semantics defined in DecisionProfile and typed references |
| SEC-03 | CONTRACT_NOW | Delegation/impersonation representation defined in actor reference schema |
| SEC-04 | CONTRACT_NOW | Signing/trust requirements defined in profile contract and bundle format |
| SEC-05 | CONTRACT_NOW | Untrusted-data validation rules defined in schema and conformance fixtures |
| SEC-06 | CONTRACT_NOW | Schema/size/depth/character/URI/reference-kind limits defined in validation contract |
| SEC-07 | CONTRACT_NOW | Data-minimization contract defined in record schemas and projection rules |
| SEC-08 | CONTRACT_NOW | Reference/digest/classification preference defined in record schema contract |
| SEC-09 | CONTRACT_NOW | Authorized-projection contract defined in projection profile semantics |
| SEC-10 | CONTRACT_NOW (carrier only) | Algorithm-agile canonical-serialization and content-digest **carrier** contract; canonicalization/digest-covered-field/signature-algorithm selection and digest computation over content are DEFERRED to ADR-F13-002 (section 18.5.1) |
| SEC-11 | CONTRACT_NOW (carrier only) | Tiered-assurance **carrier** contract defined in profile semantics and structural conformance fixtures; actual local digest computation over content is DEFERRED to ADR-F13-002 (section 18.5.1) |
| SEC-12 | CONTRACT_NOW | Fail-closed/fail-open semantics and bounded-exception contract in profile rules |
| SEC-13 | CONTRACT_NOW (contract data only) | Data-minimization/authorized-projection restrictions on AI-related retained content (prohibited categories, redaction-before-persistence, sensitivity classification, tenant-boundary constraints) defined in schema/validation/projection contract |
| SEC-13.1 | INVARIANT_FOR_LATER | Runtime handling of AI prompts, logs, traces, and model calls (emission-time redaction, running tenant-boundary enforcement) is owned by FEATURE-0025 or later; generates no FEATURE-0013 runtime task |
| SEC-14 | CONTRACT_NOW | Cross-scope reference safety-denial semantics defined in schema and negative fixtures |

### 11.4 Edge-case requirements classification

| Requirement | Classification | Rationale |
|---|---|---|
| EC-01 | CONTRACT_NOW | Non-adjudication profile schema and typed-result semantics; tested by conformance fixtures |
| EC-02 | CONTRACT_NOW | Failure/timeout/error semantic distinction defined in evaluation contract and composition strategy |
| EC-03 | CONTRACT_NOW | Profile rejection/quarantine semantics defined in bundle-validation contract |
| EC-04 | CONTRACT_NOW | Idempotency semantics and semantic-identity distinction defined in decision identity contract |
| EC-05 | CONTRACT_NOW | Partial-failure response contract and reconciliation obligation defined in acceptance semantics |
| EC-06 | CONTRACT_NOW | Chain-depth-limit rejection defined in relationship-chain validation contract |
| EC-07 | CONTRACT_NOW | Graph-budget rejection defined in composition-graph validation contract |
| EC-08 | CONTRACT_NOW | Unknown-obligation enforcement-denial defined in obligation conformance contract |
| EC-09 | CONTRACT_NOW (narrow) | AI-disabled contract behavior proven by a pure/in-memory conformance fixture per AC-49.2: `DecisionRecord` validation and deterministic composition have no AI dependency and accept the empty/unavailable/absent projection state. No `DecisionContext`, model client, prompt processing, projection service, or AI runtime is implemented. The runtime availability guarantee is AC-49.1 (INVARIANT_FOR_LATER). |
| EC-10 | CONTRACT_NOW | Cross-scope safe-denial semantics defined in reference-safety contract |

### 11.5 Compatibility requirements classification

| Requirement | Classification | Rationale |
|---|---|---|
| AC-51 (schema versioning) | CONTRACT_NOW | Version and compatibility rules defined in every canonical schema |
| AC-52 (no silent semantic change) | CONTRACT_NOW | Enforced by version-rejection semantics and conformance fixtures |
| AC-53 (versioned registries) | CONTRACT_NOW | Static registry format and offline governance contract defined |
| AC-54 (seven-value `AuditEvent` compatibility, ADH-2026-015) | CONTRACT_NOW | Seven-value compatibility fixtures, six-value FEATURE-0012 regression fixtures, and alpha-correction decision documented |
| Section 8 (Phase 1 compatibility) | CONTRACT_NOW | Compatibility evidence and non-modification constraints documented and tested |

### 11.6 Non-goal requirements classification

| Requirement | Classification | Rationale |
|---|---|---|
| NG-01 through NG-11 | DEFERRED | Each non-goal explicitly defers production runtime, service, persistence, and integration work to later approved features |

### 11.7 Classification enforcement rule

Items classified as `INVARIANT_FOR_LATER` and `DEFERRED` must not generate implementation tasks except:
- documentation and traceability records;
- non-executing conformance placeholders that assert the invariant without runtime behavior;
- architecture-decision-required markers for unresolved choices.

Any task that produces executable production behavior must trace to a `CONTRACT_NOW` requirement.

### 11.8 Normative escalation rule: ARCHITECTURE_DECISION_REQUIRED

Unresolved semantic values, numeric limits, retention periods, SLOs, product choices, component owners, legal values, jurisdiction-specific constants, and operational defaults must not be invented in design, tasks, fixtures, or code. They must be recorded as `ARCHITECTURE_DECISION_REQUIRED` with the following properties:

- **Identifier**: A stable reference label (e.g., `ADR-F13-001`).
- **Subject**: The specific value, limit, or choice that requires a decision.
- **Context**: Why this decision is needed and which acceptance criteria or design questions reference it.
- **Escalation owner**: Architecture Owner (default) or explicitly named delegate.
- **Resolution path**: The mechanism by which the decision will be resolved (ADH, DEC, RFC, or Architecture Owner directive).
- **Blocking scope**: Which downstream artifacts (design sections, tasks, fixtures, code) are blocked until resolution.

This rule applies transitively: if a design question or acceptance criterion references a value that no approved architecture source supplies, the design must escalate rather than invent. Conformance fixtures may use clearly labelled placeholder values (e.g., `PLACEHOLDER_PENDING_ADR_F13_001`) that fail validation until the architecture decision supplies the approved value.

Violations of this rule are treated as architecture drift and must be caught by the feature gate.

### 11.9 Open ARCHITECTURE_DECISION_REQUIRED register

The following open architecture decisions are escalated per section 11.8. Design, tasks, fixtures, and code must not resolve, invent, or select values for these entries until an approved decision supplies them.

| Identifier | Subject | Context (referencing items) | Escalation owner | Resolution path | Blocking scope |
|---|---|---|---|---|---|
| ADR-F13-001 | System-wide/baseline maximum ceiling for bounded decision-requirements graphs (an upper limit that no profile may exceed), if such a baseline ceiling is required | DQ-03; AC-15 | Architecture Owner | ADH / DEC / RFC or Architecture Owner directive | Any design/task/fixture/code that would assert a baseline numeric graph ceiling |
| ADR-F13-002 | Selection of the canonicalization algorithm/profile (for example RFC 8785 JCS or an alternative), the exact digest-covered field set, the signature algorithm, and the cryptographic product/service selection | AC-45; AC-45.1; DQ-04; CB-AC-10; CB-AC-11; reuse assessment section 9.1 canonicalization/trust reuse row; section 11.2 DEFERRED row | Architecture Owner | ADH / DEC / RFC or Architecture Owner directive | Design, tasks, fixtures, and code must not select the canonicalization algorithm/profile, the digest-covered field set, the signature algorithm, or the cryptographic product/service. The algorithm-agile carrier fields and structural trust metadata (AC-45) remain `CONTRACT_NOW` and may be defined; actual cryptographic verification and services remain DEFERRED until resolved |

Conformance fixtures affected by these entries may use clearly labelled placeholders (for example `PLACEHOLDER_PENDING_ADR_F13_001`, `PLACEHOLDER_PENDING_ADR_F13_002`) that fail validation until the approved decision supplies the value.

## 12. ADH-2026-014 Baseline Correction Prerequisites

### 12.1 Purpose

ADH-2026-014 approves two controlled baseline corrections. These corrections require coordinated repository updates. This section separates three prerequisite tiers: documentary prerequisites (which gate design authorization, APPROVED_FOR_DESIGN); tasks-planning prerequisites (which the approved tasks document must fully specify before APPROVED_FOR_CURSOR); and implementation and fixture prerequisites, whose executable artifacts and passing evidence gate implementation completion, final feature acceptance, and merge — not APPROVED_FOR_CURSOR.

### 12.2 Terminology reconciliation requirement

ADH-2026-014 mandates replacing `DecisionObject` with `DecisionRecord` across the repository baseline.

- PRE-01: A repository-wide audit of all occurrences of `DecisionObject` in normative documents (Phase 2 spine, current baseline, decision summary, RFC, feature index, traceability matrix, and relevant gates) must be completed.
- PRE-02: Each occurrence must be corrected to `DecisionRecord` or receive an explicit compatibility alias or migration note where backward compatibility requires it.
- PRE-03: The terminology reconciliation must be documented with a traceability record listing every changed file, the prior term, the corrected term, and the change rationale.

### 12.2.1 Canonical terminology-reconciliation traceability record

- **Canonical path:** `docs/traceability/ADH-2026-014-terminology-reconciliation.md`
- **Verification method:** Confirm the file exists at the canonical path, contains a "Corrected normative files" table listing every normative `DecisionObject` occurrence audited, and records for each entry: the file path, the prior content, the correction applied, and the rationale. Confirm a "Files intentionally NOT corrected" section documents each excluded file with disposition rationale. Confirm a "Verification" section states PRE-01, PRE-02, and PRE-03 as complete.
- **Deterministic check:** The design-authorization gate (section 12.8) verifies this artifact exists and is structurally complete before APPROVED_FOR_DESIGN can be asserted.

### 12.3 AuditEvent scope correction requirement (seven-value, ADH-2026-015)

ADH-2026-015 supersedes the earlier ADH-2026-014 six-scope statement and mandates expanding the FEATURE-0012 alpha `AuditEvent` profile from Organization-only to all seven canonical `ScopeKind` values: `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, and `ServiceInstance`, expressed through `metadata.scopeRef` as the sole scope identity. (The superseded ADH-2026-014 list omitted `Provider`.)

- PRE-05: The seven-value `AuditEvent` scope correction (ADH-2026-015) must be documented with explicit compatibility treatment: which FEATURE-0012 alpha AuditEvent consumers are affected, what migration path applies, and what backward-compatible behavior is preserved.
- PRE-06: Compatibility fixtures for all seven canonical `ScopeKind` values, plus separate regression fixtures for the six pre-existing FEATURE-0012 `ScopeKind` values, must be fully specified in the approved tasks document before APPROVED_FOR_CURSOR (covers AC-54). The executable fixtures themselves are produced by the authorized implementation after APPROVED_FOR_CURSOR; their existence and passing status are implementation-completion, final-acceptance, and merge prerequisites, not APPROVED_FOR_CURSOR prerequisites.
- PRE-07: The alpha compatibility decision must be recorded in a decision record or traceability entry with explicit approval.

### 12.4 Required repository updates

The following repository updates are mandated by ADH-2026-014:

| Update target | Required action | Evidence | Gates |
|---|---|---|---|
| `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md` | Replace `DecisionObject` with `DecisionRecord`; document the seven-value `AuditEvent` scope correction (ADH-2026-015) | File diff showing corrections | APPROVED_FOR_DESIGN |
| `docs/context/CURRENT_ARCHITECTURE_BASELINE.md` | Update baseline to reflect DecisionRecord terminology and the seven-value `AuditEvent` scope correction (ADH-2026-015) | File diff and controlled-update marker | APPROVED_FOR_DESIGN |
| `docs/context/CURRENT_DECISION_SUMMARY.md` | Add ADH-2026-014 decision entries | Entry presence verification | APPROVED_FOR_DESIGN |
| `docs/rfc/RFC-0023-decision-and-audit-standard.md` | Update to reference DecisionRecord terminology | File diff | APPROVED_FOR_DESIGN |
| `docs/features/FEATURE_INDEX.md` | Update FEATURE-0013 entry with correct terminology | Entry presence verification | APPROVED_FOR_DESIGN |
| `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md` | Update FEATURE-0013 traceability row and DecisionRecord references | Entry presence verification | APPROVED_FOR_DESIGN |
| `docs/phase2/PHASE2_ACCEPTANCE_GATES.md` | Add FEATURE-0013 gate criteria | Entry presence verification | APPROVED_FOR_DESIGN |

### 12.5 FEATURE-0012 compatibility evidence

- PRE-08: FEATURE-0012 compatibility evidence must demonstrate that the seven-value `AuditEvent` scope correction (ADH-2026-015) does not break FEATURE-0012 conformance fixtures or invalidate its feature-gate evidence.
- PRE-09: If FEATURE-0012 conformance fixtures require update for seven-value `AuditEvent` compatibility, those updates must be fully planned in the approved tasks document before APPROVED_FOR_CURSOR, and must be completed and pass before FEATURE-0013 implementation completion, final feature acceptance, and merge.
- PRE-10: The FEATURE-0012 human semantic review evidence (`docs/reviews/feature-gates/FEATURE-0012-human-semantic-review-evidence.md`) must confirm acceptance of the AuditEvent scope correction or identify required remediation.

PRE-08 through PRE-10 are implementation-completion, final-acceptance, and merge prerequisites. The approved tasks document must fully plan the work required to satisfy PRE-08 through PRE-10 before APPROVED_FOR_CURSOR, but the executable regression evidence and confirmation they require are produced by the authorized implementation after APPROVED_FOR_CURSOR and are not APPROVED_FOR_CURSOR prerequisites.

### 12.6 Prerequisite verification evidence

This section records author-collected evidence for each APPROVED_FOR_DESIGN prerequisite as of 2026-07-27. This evidence record is informational only. It is author-collected evidence and, where noted, an author self-check; it does not constitute or grant stage authorization and has not been independently verified. Fresh independent review must validate DGA-01 through DGA-05 before design authorization is established.

**DGA-01 author self-check status (2026-07-27, author-collected evidence after remediation):**

The traceability record at `docs/traceability/ADH-2026-014-terminology-reconciliation.md` exists and was inspected by the author after remediation. This is author-collected evidence, not independent verification. The following structural conditions were observed by the author:

- ✅ File exists at the canonical path.
- ✅ "Corrected normative files" table is present with 11 entries (ADH-2026-014 active-normative corrections).
- ✅ Each entry includes all four required columns: File, Prior content, Correction applied, Rationale.
- ✅ The `FEATURE_TRACEABILITY_MATRIX.md` row now includes its Rationale column: "Traceability matrix must reflect current feature status and approved corrections".
- ✅ The ADH-2026-016 update table records 2 additional active-normative corrections (`docs/architecture/api-resource-standard.md` and `docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md`) plus the 2 frozen FEATURE-0012 Kiro artifacts under an explicit artifact-level compatibility note.
- ✅ "Files intentionally NOT corrected" section is present with 7 rows; following the ADH-2026-016 update, 2 rows are reclassified as corrected active-normative (superseded), 2 rows receive frozen-artifact treatment (the `.kiro/specs` requirements.md and design.md), and 3 rows remain genuinely excluded historical/self-referencing records, each with a disposition rationale.
- ✅ "Verification" section states PRE-01, PRE-02, and PRE-03 as Complete.

**DGA-01 author self-check result: PASS (author-collected evidence; pending fresh independent verification).**

The traceability record is structurally complete on the author's inspection. The previously identified defect (missing Rationale column for the `FEATURE_TRACEABILITY_MATRIX.md` row) has been corrected; this correction is recorded here as author-collected evidence and still requires fresh independent verification.

| Prerequisite | Collected status | Evidence |
|---|---|---|
| PRE-01 (repository-wide DecisionObject audit) | SATISFIED | `docs/traceability/ADH-2026-014-terminology-reconciliation.md` "Corrected normative files" table lists 11 normative files audited with all four columns complete; the ADH-2026-016 update table records 2 further active-normative corrections; the "Files intentionally NOT corrected" section (7 rows) documents each disposition (2 superseded/now corrected, 2 frozen-artifact, 3 historical/self-referencing) with rationale; "Verification" section states PRE-01 Complete. DGA-01 author self-check: PASS (author-collected evidence; pending independent verification). |
| PRE-02 (each occurrence corrected or aliased) | SATISFIED | Same traceability record confirms all active-normative uses are replaced with `DecisionRecord` — 11 files under ADH-2026-014 plus `docs/architecture/api-resource-standard.md` and `docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md` under ADH-2026-016 (13 active-normative corrections). Only the two frozen FEATURE-0012 Kiro artifacts (`.kiro/specs/api-resource-naming-status-and-validation-standard/requirements.md` and `design.md`) retain the historical `DecisionObject` name, under the explicit artifact-level compatibility note authorized by ADH-2026-016 (no second schema, type, alias, or contract). DGA-01 author self-check: PASS (author-collected evidence; pending independent verification). |
| PRE-03 (traceability record exists) | SATISFIED | `docs/traceability/ADH-2026-014-terminology-reconciliation.md` exists at the canonical path, lists every changed file, prior term, corrected term, and rationale. All table rows structurally complete. DGA-01 author self-check: PASS (author-collected evidence; pending independent verification). |
| PRE-05 (seven-value AuditEvent scope correction documented) | EVIDENCE_RECORDED | ADH-2026-015 and section 17 document the canonical seven-value ScopeKind and AuditEvent contract, including the canonical absent/nil `Platform` form. The controlling Phase 2 spine, acceptance gates, current baseline, decision summary, RFC, feature index, and traceability matrix identify ADH-2026-014, ADH-2026-015, and ADH-2026-016 as the joint controlling handoffs (ADH-2026-015 authoritative for the seven-value ScopeKind contract; ADH-2026-016 authoritative for the clarified contract boundaries) and record the seven-value contract at the documentary level. Schema, Go-binding, validator, and fixture synchronization remains an obligation pending implementation evidence and is not claimed complete. Historical ADH-2026-014 six-scope evidence remains preserved as superseded compatibility history. |
| PRE-07 (alpha compatibility decision recorded with approval) | EVIDENCE_RECORDED | ADH-2026-015 approves the seven-value `ScopeKind`/`AuditEvent` correction (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance` via `metadata.scopeRef`), superseding the earlier ADH-2026-014 six-scope statement. Human approval: Sanjeev Kumar, 2026-07-27. The superseded ADH-2026-014 approval remains recorded in `docs/reviews/architecture-decision-handoffs/ADH-2026-014-feature-0013-decision-record-and-auditevent-standard.md` and `docs/traceability/ADH-2026-014-terminology-reconciliation.md` as compatibility history. |
| Section 12.4: PHASE2_ARCHITECTURE_SPINE.md | EVIDENCE_RECORDED | Uses DecisionRecord terminology and the canonical seven-value ScopeKind; AuditEvent permits all seven through metadata.scopeRef as sole scope identity. Identifies ADH-2026-014, ADH-2026-015, and ADH-2026-016 as the joint controlling handoffs: ADH-2026-015 supersedes only the earlier ADH-2026-014 six-scope statement (explicitly identified as superseded), and ADH-2026-016 is the clarification of the DecisionRecord scope-authority, sensitivity, security-validation, structural-trust, conformance-ID, evidence, and terminology boundaries. |
| Section 12.4: CURRENT_ARCHITECTURE_BASELINE.md | EVIDENCE_RECORDED | Identifies ADH-2026-016 as the latest controlled clarification/update. Records ADH-2026-014, ADH-2026-015, and ADH-2026-016 as joint controlling handoffs with their precise relationships: ADH-2026-015 supersedes only the earlier six-scope statement and remains authoritative for the seven-value ScopeKind contract; ADH-2026-016 is authoritative for the clarified DecisionRecord scope authority via metadata.scopeRef, the closed ordered sensitivity vocabulary, the structural-only security-validation boundary, the structural-trust versus deferred-cryptographic boundary (ADR-F13-002), the exact F13-CF-01 through F13-CF-28 conformance-ID coverage, the SUPERSEDED disposition of six-scope evidence, and the terminology boundaries. Records the seven canonical ScopeKind values, all-seven AuditEvent applicability, and metadata.scopeRef as sole scope identity. |
| Section 12.4: CURRENT_DECISION_SUMMARY.md | EVIDENCE_RECORDED | Identifies ADH-2026-014, ADH-2026-015, and ADH-2026-016 as joint controlling handoffs with their precise relationships. Preserves the superseded ADH-2026-014 six-scope statement as history; records ADH-2026-015 as authoritative for the seven-value ScopeKind correction; and records ADH-2026-016 as the authoritative clarification of the DecisionRecord scope-authority, sensitivity, security-validation, structural-trust, conformance-ID, evidence, and terminology boundaries. |
| Section 12.4: RFC-0023-decision-and-audit-standard.md | EVIDENCE_RECORDED | Uses DecisionRecord terminology and records ADH-2026-014, ADH-2026-015, and ADH-2026-016 as joint controlling handoffs with their precise relationships: ADH-2026-015 supersedes only the earlier six-scope statement and is authoritative for the seven-value scope contract; ADH-2026-016 is the authoritative clarification of the DecisionRecord, sensitivity, security-validation, trust, conformance-ID, evidence, and terminology boundaries. |
| Section 12.4: FEATURE_INDEX.md | EVIDENCE_RECORDED | FEATURE-0013 uses the canonical Decision Record title and identifies ADH-2026-014, ADH-2026-015, and ADH-2026-016 as joint controlling handoffs (ADH-2026-015 authoritative for the seven-value ScopeKind contract; ADH-2026-016 authoritative for the clarified contract boundaries). |
| Section 12.4: FEATURE_TRACEABILITY_MATRIX.md | EVIDENCE_RECORDED | FEATURE-0013 records the seven-value ScopeKind, all-seven AuditEvent applicability, metadata.scopeRef sole authority, ADH-2026-015 compatibility evidence, and pending fresh requirements review. Identifies ADH-2026-014, ADH-2026-015, and ADH-2026-016 as joint controlling handoffs with their precise relationships (ADH-2026-015 authoritative for the seven-value ScopeKind contract; ADH-2026-016 authoritative for the clarified DecisionRecord scope-authority, sensitivity, security-validation, structural-trust, conformance-ID, evidence, and terminology boundaries). |
| Section 12.4: PHASE2_ACCEPTANCE_GATES.md | EVIDENCE_RECORDED | FEATURE-0013 gates identify ADH-2026-014, ADH-2026-015, and ADH-2026-016 as joint controlling handoffs with their precise relationships (ADH-2026-015 superseding only the earlier six-scope statement and authoritative for the seven-value ScopeKind contract; ADH-2026-016 authoritative for the clarified DecisionRecord scope-authority, sensitivity, security-validation, structural-trust, conformance-ID, evidence, and terminology boundaries) and require the canonical seven-value ScopeKind, all-seven AuditEvent fixtures (canonical absent/nil `Platform`, non-Platform via `metadata.scopeRef`), a negative absence-when-`Platform`-disallowed fixture, six-value FEATURE-0012 regression coverage, and `metadata.scopeRef` as sole scope authority. Schema/binding/validator/fixture synchronization remains pending implementation evidence. |
| Terminology-reconciliation traceability record exists (section 12.2.1) | SATISFIED | `docs/traceability/ADH-2026-014-terminology-reconciliation.md` exists at canonical path. Structure is complete: all 11 "Corrected normative files" rows have four columns; the ADH-2026-016 update table records 2 further active-normative corrections and the 2 frozen FEATURE-0012 Kiro artifacts; the "Files intentionally NOT corrected" section (7 rows: 2 superseded/now corrected, 2 frozen-artifact, 3 historical/self-referencing) is documented with rationale; PRE-01/02/03 marked Complete. DGA-01 author self-check: PASS (author-collected evidence; pending independent verification). |
| FEATURE-0012 dependency readiness (section 13.4) | EVIDENCE_RECORDED | DEP-01 through DEP-04 all satisfied per section 13.7. |
| DGA-05 (ADH-2026-016 synchronization) | EVIDENCE_RECORDED (author self-check) | Author-collected evidence only, pending fresh independent verification. Author self-check observed: DecisionRecord uses `metadata.scopeRef` as sole scope authority (AC-01.1, CB-AC-01 through CB-AC-03); the closed ordered sensitivity vocabulary `PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED` with mandatory profile bounds (CB-AC-04 through CB-AC-07); structural-only security validation (SEC-15, CB-AC-08, CB-AC-09); structural-trust versus deferred-cryptographic separation under ADR-F13-002 (AC-45, AC-45.1, CB-AC-10, CB-AC-11, section 11.9), now reconciled by the single pre-ADR-F13-002 digest and trust conformance boundary (section 18.5.1, CB-AC-16, CB-AC-17) covering AC-13.3, AC-13.7(f)–(g), SEC-10, SEC-11, AC-45/45.1, and F13-CF-13/23/27/28 — with AC-13.3(b), SEC-10, and SEC-11 revised so no FEATURE-0013 digest fixture claims to "prove" content existed and no local digest computation over content is claimed now, and pre-ADR-F13-002 executable checks limited to carrier presence, identifier syntax, covered-field-descriptor presence, opaque expected-versus-observed state, and fail-closed trust-state transitions; exact `F13-CF-01` through `F13-CF-28` mapping (AC-50, CB-AC-12, CB-AC-13); six-scope evidence marked SUPERSEDED and ADH-2026-015 evidence authoritative for the seven-value contract (sections 1.5, 17); the authoritative compatibility evidence (`docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md`) corrected so `metadata.scopeRef` is the sole but conditionally-present scope source (no longer described as universally required) and consistent with the canonical nil-`Platform` representation (section 10 "Required fields" row and correction note; SV-AC-11); active `DecisionObject`-to-`DecisionRecord` migration recorded with active-normative corrections (including `docs/architecture/api-resource-standard.md` and `docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md` corrected under ADH-2026-016) and frozen-artifact compatibility notes limited to the two FEATURE-0012 Kiro specs (section 12.2.1, `docs/traceability/ADH-2026-014-terminology-reconciliation.md`); and active baseline/index/RFC/traceability synchronization identifying ADH-2026-016 as the latest controlled clarification/update (section 12.4 rows above). This is an author self-check only; a fresh independent reviewer must independently verify each condition. |

**Evidence collection status:** Author self-check complete; fresh independent review pending. The results recorded in this section are author-collected evidence and author self-checks; none of them are independent verifications, independent PASS verdicts, or stage authorizations. On the author's inspection PRE-01 through PRE-03 are satisfied, and the DGA-01 terminology-reconciliation structural defect has been remediated. The DGA-03 evidence condition has been corrected to cite the actual location where ADH-2026-015 records the supersession: the `Predecessor` entry in the body `## Metadata` section of `docs/reviews/architecture-decision-handoffs/ADH-2026-015-feature-0013-scope-vocabulary-clarification.md`, together with that handoff's Summary, Classification, Required action, and Human approval sections. ADH-2026-015 carries no YAML/front-matter `predecessor` field, so no requirements statement or deterministic check may assert one. DGA-05 records author-collected evidence and an author self-check for the ADH-2026-016 synchronization condition. On the corrected conditions, the author self-check records DGA-01 through DGA-05 as PASS (author-collected evidence only). This revision additionally applied two reviewer-required corrections and re-ran the DGA-01 through DGA-05 author self-check against them: (1) the authoritative compatibility evidence (`docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md`) was corrected so `metadata.scopeRef` is recorded as the sole but conditionally-present scope source rather than "universally required", removing the contradiction with the canonical nil-`Platform` representation (affecting the DGA-05 compatibility-consistency condition and reinforcing DGA-03); and (2) a single normative pre-ADR-F13-002 digest and trust conformance boundary was added (section 18.5.1, CB-AC-16/CB-AC-17) and AC-13.3(b), AC-13.7(f)–(g), SEC-10, and SEC-11 were revised so no digest-related scenario claims cryptographic "proof" or present-tense local digest computation over content, preserving the deferred cryptographic boundary (affecting the DGA-05 structural-trust/deferred-cryptographic condition). This is not a stage authorization: under ADH-2026-015 the previous `APPROVED_FOR_DESIGN` is superseded and, following the ADH-2026-016 clarification, fresh independent review must validate DGA-01 through DGA-05 (including the corrected DGA-03 condition and the DGA-05 condition, now including the compatibility-consistency and digest/trust-boundary corrections) and return the APPROVED_FOR_DESIGN token per section 12.8 before design may proceed.

### 12.7 Authorization gate structure

FEATURE-0013 authorization follows a strict, non-circular sequence. Each stage token is returned only after the prior stage is approved:

1. Requirements approval returns APPROVED_FOR_DESIGN. DGA-01 through DGA-05 (section 12.8) gate this transition.
2. Design approval returns APPROVED_FOR_TASKS.
3. Tasks approval returns APPROVED_FOR_CURSOR. At this gate the approved `tasks.md` must explicitly and completely plan PRE-06, PRE-08 through PRE-10, BC-AC-04, BC-AC-05, the seven-scope fixtures, the six-value FEATURE-0012 regression coverage, schema/binding/validator synchronization, and compatibility evidence. The resulting executable artifacts and passing evidence are not required to exist yet.
4. Authorized implementation then produces those executable artifacts and passing evidence.
5. Implementation completion, final feature acceptance, and merge are blocked until PRE-06, PRE-08 through PRE-10, BC-AC-04, BC-AC-05, all fixtures, compatibility checks, regression results, and synchronization evidence exist and pass.

The staged prerequisites are:

**APPROVED_FOR_DESIGN requires:**

1. PRE-01 through PRE-03 are satisfied (documentary terminology reconciliation);
2. PRE-05 and PRE-07 are satisfied (documentary seven-value `AuditEvent` scope correction and compatibility decision, ADH-2026-015);
3. All repository updates in section 12.4 are committed and verified;
4. The terminology reconciliation traceability record exists at `docs/traceability/ADH-2026-014-terminology-reconciliation.md`;
5. FEATURE-0012 dependency readiness condition (section 13.4) is confirmed;
6. Independent review validates DGA-01 through DGA-05 (section 12.8) and returns the APPROVED_FOR_DESIGN token.

**APPROVED_FOR_TASKS requires:**

7. APPROVED_FOR_DESIGN has been granted;
8. Independent design review approves the design and returns the APPROVED_FOR_TASKS token.

**APPROVED_FOR_CURSOR requires:**

9. APPROVED_FOR_TASKS has been granted;
10. An approved tasks document (`tasks.md`) explicitly and completely plans the work required to satisfy PRE-06, PRE-08 through PRE-10, BC-AC-04, and BC-AC-05, including: seven-scope `AuditEvent` compatibility fixtures; six-value FEATURE-0012 regression coverage; FEATURE-0012 compatibility and human semantic review evidence; schema, Go-binding, and validator synchronization for the seven-value `ScopeKind`/`AuditEvent` contract; and the corresponding compatibility evidence. The resulting executable artifacts and passing evidence are NOT required to exist at APPROVED_FOR_CURSOR; only their complete plan in the approved tasks document is required.

**Implementation completion, final feature acceptance, and merge additionally require:**

11. PRE-06 is satisfied (executable compatibility fixtures for all seven canonical `ScopeKind` values, plus regression fixtures for the six pre-existing FEATURE-0012 values, exist and pass);
12. PRE-08 through PRE-10 are satisfied (FEATURE-0012 regression evidence and human semantic review confirmation);
13. BC-AC-04 and BC-AC-05 are satisfied (executable fixture and regression evidence exist and pass);
14. Schema, Go-binding, validator, and fixture synchronization evidence for the seven-value contract exists and passes.

Design and task generation may proceed once the APPROVED_FOR_DESIGN prerequisites are met and the independent review gate validates DGA-01 through DGA-05 and returns the approval token. The executable fixtures, regression evidence, and synchronization artifacts required by PRE-06, PRE-08 through PRE-10, BC-AC-04, and BC-AC-05 are produced by the authorized implementation after APPROVED_FOR_CURSOR; they are not prerequisites for APPROVED_FOR_CURSOR. They gate implementation completion, final feature acceptance, and merge.

Failure of any APPROVED_FOR_DESIGN precondition pauses design authorization. Failure of the tasks-planning coverage precondition pauses APPROVED_FOR_CURSOR. Failure of any implementation-completion, final-acceptance, or merge precondition (PRE-06, PRE-08 through PRE-10, BC-AC-04, BC-AC-05, or the synchronization evidence) blocks implementation completion, final feature acceptance, and merge, but does not retroactively invalidate APPROVED_FOR_CURSOR.

### 12.8 Deterministic design-authorization check

Before FEATURE-0013 may progress from requirements to design, the following deterministic checks must pass. These checks are verifiable by any agent or reviewer without subjective judgement.

**Important: This requirements document defines the checks but does not grant its own stage authorization. An independent review gate must validate DGA-01 through DGA-05 and return the APPROVED_FOR_DESIGN token before design generation may proceed.**

**Check DGA-01: Terminology-reconciliation traceability record exists and is structurally complete.**

- Verify: `docs/traceability/ADH-2026-014-terminology-reconciliation.md` exists.
- Verify: The file contains a "Corrected normative files" table with at least one entry.
- Verify: Each entry includes columns: File, Prior content, Correction applied, Rationale.
- Verify: The file contains a "Files intentionally NOT corrected" section.
- Verify: The file contains a "Verification" section stating PRE-01, PRE-02, and PRE-03 as complete.
- Failure action: APPROVED_FOR_DESIGN is blocked until the record is created or repaired.

**Check DGA-02: Required baseline documents reflect DecisionRecord terminology.**

- Verify: `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md` contains `DecisionRecord` in normative text and does not contain `DecisionObject` as a normative term (historical/migration notes excluded).
- Verify: `docs/context/CURRENT_ARCHITECTURE_BASELINE.md` contains `DecisionRecord` and references ADH-2026-014.
- Verify: `docs/context/CURRENT_DECISION_SUMMARY.md` contains at least one entry referencing ADH-2026-014.
- Verify: `docs/rfc/RFC-0023-decision-and-audit-standard.md` contains `DecisionRecord`.
- Verify: `docs/features/FEATURE_INDEX.md` FEATURE-0013 entry contains "Decision Record".
- Verify: `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md` FEATURE-0013 row references ADH-2026-014.
- Verify: `docs/phase2/PHASE2_ACCEPTANCE_GATES.md` contains FEATURE-0013 gate criteria and identifies ADH-2026-014, ADH-2026-015, and ADH-2026-016 as joint controlling handoffs with their precise relationships: ADH-2026-015 supersedes only the earlier six-scope `AuditEvent` statement and remains authoritative for the seven-value scope contract; ADH-2026-016 clarifies the additional contract boundaries without granting stage authorization.
- Failure action: APPROVED_FOR_DESIGN is blocked until the missing update is applied.

**Check DGA-03: Seven-value `AuditEvent` scope correction (ADH-2026-015) is recorded with human approval.**

- Verify: The ADH-2026-015 architecture-decision handoff exists and records the seven-value `ScopeKind`/`AuditEvent` decision that supersedes the earlier ADH-2026-014 six-scope statement.
- Verify: The controlling architecture (`docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md` section 29) and `docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md` record the seven canonical `ScopeKind` values (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance`) and `metadata.scopeRef` as the sole scope identity.
- Verify: The "Human approval" section records approval status as "Approved", approver identity, and date.
- Verify: The supersession of the earlier six-scope `AuditEvent` statement is authoritatively recorded by the approved ADH-2026-015 handoff itself. `docs/reviews/architecture-decision-handoffs/ADH-2026-015-feature-0013-scope-vocabulary-clarification.md` records — in the `Predecessor` entry of its body `## Metadata` section, and in its `## Summary`, `## Classification`, `## Required action`, and `## Human approval` sections — that it supersedes the ADH-2026-014 six-scope `AuditEvent` statement while ADH-2026-014 otherwise remains controlling. ADH-2026-015 carries no YAML/front-matter `predecessor` field; the supersession is recorded in the body `## Metadata` `Predecessor` entry and the sections named above, and the check must not assert any front-matter `predecessor` field. This approved superseding handoff is the controlled-documentation supersession artifact; the corresponding supersession notices in sections 1.5 and 17 of this requirements document restate it. The check does not require an annotation to be added to the superseded ADH-2026-014 file, because the controlled-documentation practice records supersession in the approved superseding artifact. (This condition was corrected across revisions: an earlier condition incorrectly required an annotation inside the ADH-2026-014 file, and a later revision incorrectly cited a nonexistent `predecessor` front-matter field on ADH-2026-015; the accurate evidence is the body `## Metadata` `Predecessor` entry plus the Summary, Classification, Required action, and Human approval sections of ADH-2026-015.)
- Failure action: APPROVED_FOR_DESIGN is blocked until the seven-value compatibility decision is approved.

**Check DGA-04: FEATURE-0012 dependency readiness confirmed.**

- Verify: Section 13.7 of this requirements document records DEP-01 through DEP-04 as SATISFIED with evidence.
- Verify: `docs/reviews/feature-gates/FEATURE-0012-approval-review.md` contains "Final feature-review status: Approved".
- Failure action: APPROVED_FOR_DESIGN is blocked until FEATURE-0012 readiness is confirmed.

**Check DGA-05: ADH-2026-016 contract-boundary clarification is synchronized and recorded.**

This check requires and independently verifies each of the following. It records only an author self-check and author-collected evidence (section 12.6); only a fresh independent reviewer may return APPROVED_FOR_DESIGN.

- Verify: `DecisionRecord` uses FEATURE-0012 `metadata.scopeRef` as its sole logical and serialized scope authority, with no top-level or parallel scope source (AC-01.1, CB-AC-01, CB-AC-03).
- Verify: The sensitivity vocabulary is closed and ordered `PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED`, and any profile permitting captured classified content declares the mandatory profile bounds — sensitivity ceiling, allowed/prohibited content-category rules, maximum field count, maximum byte size, maximum nesting depth, and allowed value types (CB-AC-04 through CB-AC-07).
- Verify: FEATURE-0013 security validation is bounded to structural-only validation, typed classification, declared bounds, projection restrictions, and deterministic conformance, and introduces no secret-pattern, credential-pattern, PII, malware, DLP, or content-scanning engine (SEC-15, CB-AC-08, CB-AC-09).
- Verify: Structural-trust metadata and algorithm-agile carrier fields remain `CONTRACT_NOW` and are separated from the deferred cryptographic selections blocked by ADR-F13-002 — canonicalization algorithm/profile, digest-covered field set, signature algorithm, and cryptographic product/service selection (AC-45, AC-45.1, CB-AC-10, CB-AC-11, section 11.9 register).
- Verify: The single pre-ADR-F13-002 digest and trust conformance boundary (section 18.5.1, CB-AC-16 and CB-AC-17) is present and reconciles AC-13.3, AC-13.7(f)–(g), SEC-10, SEC-11, AC-45/AC-45.1, CB-AC-10/CB-AC-11, and scenarios F13-CF-13/23/27/28. Confirm that pre-ADR-F13-002 executable checks are limited to carrier presence, identifier syntax, covered-field-descriptor presence, opaque expected-versus-observed state, and fail-closed trust-state transitions; and that canonical-serialization execution, content-to-digest computation, signature generation/verification, algorithm/covered-field/cryptographic-product selection, and cryptographic-validity/"proof" claims are prohibited until ADR-F13-002 is resolved. Confirm AC-13.3(b), SEC-10, and SEC-11 no longer assert that a FEATURE-0013 digest fixture "proves" content existed or that local digest computation occurs now.
- Verify: The authoritative compatibility evidence `docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md` no longer describes `metadata.scopeRef` as universally required and is consistent with the canonical nil-`Platform` representation (section 10 "Required fields" row and correction note; SV-AC-03, SV-AC-11, AC-24, CB-AC-02).
- Verify: The architecture section 17 scenarios map exactly and one-to-one to `F13-CF-01` through `F13-CF-28` in their existing order, with no merge, renumber, or omission, and separately named case IDs do not alter the canonical 28-scenario count (AC-50, CB-AC-12, CB-AC-13).
- Verify: The earlier six-scope `AuditEvent` evidence is marked SUPERSEDED and the ADH-2026-015 seven-value `ScopeKind` evidence is authoritative (sections 1.5, 17, SV-AC-01).
- Verify: The active `DecisionObject`-to-`DecisionRecord` migration is recorded with every active-normative document corrected — including `docs/architecture/api-resource-standard.md` and the FEATURE-0012 feature document (`docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md`) corrected under ADH-2026-016 — and frozen-artifact compatibility notes limited to the two FEATURE-0012 Kiro specs (section 12.2.1; `docs/traceability/ADH-2026-014-terminology-reconciliation.md`).
- Verify: The active baseline, feature index, RFC, and traceability documents are synchronized with ADH-2026-016 and identify it as the latest controlled clarification/update alongside ADH-2026-014 and ADH-2026-015 as joint controlling handoffs (section 12.4 and section 12.6 evidence rows).
- Failure action: APPROVED_FOR_DESIGN is blocked until the ADH-2026-016 synchronization conditions are satisfied and independently verified.

**Gate result:** If DGA-01, DGA-02, DGA-03, DGA-04, and DGA-05 all pass under independent review, the reviewer may grant APPROVED_FOR_DESIGN. If any check fails, design generation must not proceed and the failing check must be remediated first.

**Current gate status (2026-07-27, author self-check pending independent review):**

| Check | Evidence collected | Evidence summary |
|---|---|---|
| DGA-01 | PASS (author self-check) | `docs/traceability/ADH-2026-014-terminology-reconciliation.md` exists with 11 corrected active-normative files (all four columns complete including Rationale) plus 2 further active-normative corrections in the ADH-2026-016 update table; the 7-row "Files intentionally NOT corrected" section documents 2 superseded/now-corrected rows, 2 frozen-artifact rows, and 3 historical/self-referencing rows with rationale; PRE-01/02/03 marked Complete. |
| DGA-02 | PASS (author self-check) | All seven section 12.4 documents verified per section 12.6 evidence table. `DecisionRecord` is present in active normative text; `DecisionObject` remains only in explicitly treated historical/migration notes. `PHASE2_ACCEPTANCE_GATES.md` identifies ADH-2026-014, ADH-2026-015, and ADH-2026-016 as joint controlling handoffs with their precise supersession and clarification relationships. |
| DGA-03 | PASS (author self-check, condition corrected) | The supersession of the ADH-2026-014 six-scope `AuditEvent` statement is authoritatively recorded by the approved superseding handoff `ADH-2026-015` itself — in the `Predecessor` entry of its body `## Metadata` section, and in its Summary, Classification, Required action, and Human approval sections — which is the controlled-documentation supersession artifact; sections 1.5 and 17 of this document restate it. ADH-2026-015 has no YAML/front-matter `predecessor` field; the supersession is recorded in the body `## Metadata` `Predecessor` entry and the sections named above. Human approval: Approved, Sanjeev Kumar, 2026-07-27. This is an author self-check only; a fresh independent reviewer must confirm this corrected condition. |
| DGA-04 | PASS (author self-check) | Section 13.7 confirms DEP-01 through DEP-04 SATISFIED; FEATURE-0012-approval-review.md confirms "Final feature-review status: Approved". |
| DGA-05 | PASS (author self-check) | ADH-2026-016 synchronization observed by the author: `metadata.scopeRef` sole authority (AC-01.1, CB-AC-01/03); closed ordered sensitivity vocabulary `PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED` with mandatory profile bounds (CB-AC-04 through CB-AC-07); structural-only security validation (SEC-15, CB-AC-08/09); structural-trust versus deferred-cryptographic separation under ADR-F13-002 covering canonicalization algorithm/profile, digest-covered field set, signature algorithm, and cryptographic product/service (AC-45, AC-45.1, CB-AC-10/11, section 11.9), now reconciled by the single pre-ADR-F13-002 digest and trust conformance boundary (section 18.5.1, CB-AC-16/17) with AC-13.3(b), AC-13.7(f)–(g), SEC-10, and SEC-11 revised to remove "proof"/"local digest computation now" wording and to limit pre-ADR-F13-002 checks to carrier presence, identifier syntax, covered-field-descriptor presence, opaque expected-versus-observed state, and fail-closed trust-state transitions; exact `F13-CF-01` through `F13-CF-28` mapping (AC-50, CB-AC-12/13); six-scope evidence marked SUPERSEDED and ADH-2026-015 evidence authoritative (sections 1.5, 17); the authoritative compatibility evidence corrected so `metadata.scopeRef` is the sole but conditionally-present scope source (no longer "universally required"), consistent with the canonical nil-`Platform` form (`docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md` section 10, SV-AC-11); active `DecisionObject` migration with every active-normative document corrected (including `api-resource-standard.md` and the FEATURE-0012 feature document under ADH-2026-016) and frozen-artifact compatibility notes limited to the two FEATURE-0012 Kiro specs (section 12.2.1); active baseline/index/RFC/traceability synchronized with ADH-2026-016 as latest controlled clarification/update (section 12.4, 12.6). Author self-check only; a fresh independent reviewer must confirm each condition. |

**Gate status: PENDING_INDEPENDENT_REVIEW. This document does not self-approve.**

Under ADH-2026-015 the previous `APPROVED_FOR_DESIGN` is superseded for continued design work, and this requirements document must receive fresh independent approval before design resumes. ADH-2026-016 (Approved 2026-07-27) subsequently returned this requirements document to **architecture-remediation status** to apply the clarified contract boundaries in section 18 and to synchronize the active baseline documents (section 12.4, 12.6); it did not grant, restore, or advance any stage authorization and did not resolve the pending independent review. The gate status therefore remains `PENDING_INDEPENDENT_REVIEW` and the document remains `PENDING_HUMAN_REVIEW`. The author self-check of DGA-01 through DGA-05 recorded above is evidence only; it is not a stage authorization and must not be read as one. An independent reviewer must validate DGA-01 through DGA-05 (including the corrected DGA-03 condition and the DGA-05 ADH-2026-016 synchronization condition) and return the APPROVED_FOR_DESIGN token before design generation may proceed. No design, tasks, or implementation may proceed on the basis of this self-check or on the basis of the ADH-2026-016 remediation.

### 12.9 Acceptance criteria for baseline corrections

- BC-AC-01: Every normative document listed in section 12.4 contains `DecisionRecord` (not `DecisionObject`) in all normative references to the common decision envelope. (CONTRACT_NOW — gates APPROVED_FOR_DESIGN)
- BC-AC-02: The Phase 2 spine reflects the seven canonical `ScopeKind` values for `AuditEvent` (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance`) per ADH-2026-015. (CONTRACT_NOW — gates APPROVED_FOR_DESIGN)
- BC-AC-03: A terminology reconciliation traceability record exists listing all corrected files and rationale. (CONTRACT_NOW — gates APPROVED_FOR_DESIGN)
- BC-AC-04: `AuditEvent` compatibility fixtures exist for all seven canonical `ScopeKind` values (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance`), and separate regression fixtures cover the six pre-existing FEATURE-0012 values (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`). (CONTRACT_NOW — the approved tasks document must fully plan this before APPROVED_FOR_CURSOR; the executable fixtures are produced by the authorized implementation and gate implementation completion, final feature acceptance, and merge)
- BC-AC-05: FEATURE-0012 conformance evidence demonstrates no regression from the AuditEvent scope correction. (CONTRACT_NOW — the approved tasks document must fully plan this before APPROVED_FOR_CURSOR; the passing regression evidence is produced by the authorized implementation and gates implementation completion, final feature acceptance, and merge)
- BC-AC-06: The alpha compatibility decision for AuditEvent scope expansion is recorded with human approval reference. (CONTRACT_NOW — gates APPROVED_FOR_DESIGN)

## 13. FEATURE-0012 Dependency Readiness Condition

### 13.1 Required dependency

FEATURE-0013 depends on FEATURE-0012 (API, Resource Naming, Status, and Validation Standard) for the common resource grammar, metadata conventions, reference semantics, validation patterns, error contracts, and compatibility rules.

### 13.2 Precise dependency specification

| Dependency attribute | Required value | Current evidence |
|---|---|---|
| Feature ID | FEATURE-0012 | Confirmed |
| Minimum stage | Implemented and Merged | Confirmed: approved 2026-07-24, merged via commit a1b74fb / PR #14 |
| Final human approval | Final feature-review approval by architecture owner | Confirmed: approved 2026-07-24 by Sanjeev Kumar (`docs/reviews/feature-gates/FEATURE-0012-approval-review.md`) |
| Required contract version | `sovrunn.io/v1alpha1` resource grammar as implemented in FEATURE-0012 tasks | Confirmed: checkpoint 18 passed, conformance fixtures passing |
| Required artifacts | Common metadata, reference, scope, condition, validation, error, and compatibility contracts from FEATURE-0012 implementation | Confirmed: 231 implementation paths completed |
| Controlling handoff | ADH-2026-012 (Approved) with ADH-2026-013 clarification | Confirmed |
| Required evidence | FEATURE-0012 feature-gate pass evidence (`docs/reviews/feature-gates/FEATURE-0012-approval-review.md`) | Confirmed: exists with `Final feature-review status: Approved` |

### 13.3 Grammar elements consumed

FEATURE-0013 consumes the following FEATURE-0012 grammar elements:

- `apiVersion`, `kind`, `metadata` resource conventions;
- typed scope and reference semantics;
- status and condition grammar;
- validation and stable error contracts;
- compatibility and version rules;
- AuditEvent alpha profile (extended to all seven canonical `ScopeKind` values per ADH-2026-015).

### 13.4 Dependency readiness gate

- DEP-01: FEATURE-0012 must have reached at minimum `automation_flow_status: final_checkpoint_passed_pending_approval` before FEATURE-0013 design authorization. **STATUS: SATISFIED** — checkpoint 18 passed 2026-07-23; final human approval 2026-07-24; merged via commit a1b74fb / PR #14.
- DEP-02: The FEATURE-0012 `sovrunn.io/v1alpha1` grammar contracts referenced by FEATURE-0013 must be implemented and passing their conformance fixtures. **STATUS: SATISFIED** — feature-gate passed; conformance fixtures passing per checkpoint 18 evidence.
- DEP-03: If FEATURE-0012 is rejected or requires rework that changes `sovrunn.io/v1alpha1` grammar semantics consumed by FEATURE-0013, an architecture escalation is mandatory (see section 13.5). **STATUS: NOT TRIGGERED** — FEATURE-0012 approved without rework.
- DEP-04: FEATURE-0013 may reference FEATURE-0012 grammar for requirements and conformance specification regardless of FEATURE-0012 human approval status, but design and tasks must not assume unimplemented or unresolved FEATURE-0012 shapes. **STATUS: SATISFIED** — all referenced shapes are implemented.

### 13.5 Architecture escalation conditions

Architecture escalation to the Sovrunn Architecture Owner is mandatory when:

- FEATURE-0012 human review rejects or requires material rework that changes grammar contracts consumed by FEATURE-0013;
- FEATURE-0012 `sovrunn.io/v1alpha1` grammar version is superseded or deprecated before FEATURE-0013 design approval;
- the seven-value `AuditEvent` scope correction (ADH-2026-015; PRE-05 through PRE-07) is rejected during FEATURE-0012 compatibility review;
- FEATURE-0012 introduces a breaking change to metadata, reference, scope, validation, or error semantics after FEATURE-0013 requirements approval;
- any FEATURE-0012 contract referenced by FEATURE-0013 becomes `ARCHITECTURE_DECISION_REQUIRED` or is marked for redesign.

Escalation pauses FEATURE-0013 design/task generation until the Architecture Owner records a resolution decision.

### 13.6 Dependency classification

This section is classified as `CONTRACT_NOW`. The dependency gate (DEP-01 through DEP-04) and escalation conditions are normative requirements that govern FEATURE-0013 authorization, not deferred work.

### 13.7 FEATURE-0012 readiness confirmation

FEATURE-0012 completed checkpoint 18, received final human approval (2026-07-24), and is implemented and merged via commit a1b74fb / PR #14. The dependency readiness condition for FEATURE-0013 design authorization is fully met.

Evidence:
- `.automation/state/FEATURE-0012.json`: `automation_flow_status: final_checkpoint_passed_pending_approval`, `status: PENDING_HUMAN_REVIEW` (historical pre-merge checkpoint state; final approval and PR #14 merge evidence below are authoritative)
- `docs/reviews/feature-gates/FEATURE-0012-approval-review.md`: `Final feature-review status: Approved`, reviewer Sanjeev Kumar, decision date 2026-07-24
- `docs/reviews/feature-gates/FEATURE-0012-human-semantic-review-evidence.md`: evidence staged and reviewed
- ADH-2026-013 accepted as architecture clarification within FEATURE-0012 scope
- Merge evidence: PR #14 merged as commit `a1b74fb` into `phase2-reuse-first-paas-fabric-foundation`

## 14. Audit-Recording Boundary Specification

### 14.1 Purpose

This section provides an unambiguous, testable boundary between durable `DecisionRecord` instances governed by AC-21 and ephemeral governed-aggregation authorization events governed by AC-40. This boundary is required for validation and conformance fixtures.

### 14.2 Definitions

- **Durable DecisionRecord**: An immutable, persisted `DecisionRecord` instance that satisfies AC-21 by linking to at least one `AuditEvent`. It is the authoritative Sovrunn decision artifact.
- **Ephemeral governed authorization evaluation**: A high-volume, low-risk authorization check that satisfies its audit obligation through the governed aggregation profile (AC-40) without producing an individual persisted `DecisionRecord`.

### 14.3 Distinguishing criteria

An authorization event becomes an ephemeral governed evaluation (not a durable DecisionRecord) when ALL of the following hold:

1. The applicable `DecisionProfile` explicitly classifies the action category as eligible for governed aggregation;
2. The authorization outcome is `ALLOWED`;
3. No mandatory obligations beyond standard enforcement are produced;
4. None of the following escalation flags are present: denial, privilege escalation, cross-scope access, cross-tenant access, break-glass, data disclosure, data export, policy change, privilege change, human approval requirement, or regulated-action classification;
5. No profile-specific escalation trigger is active for the current scope, actor, or resource;
6. The aggregation window and minimum-evidence parameters are within the profile's declared bounds.

If ANY criterion fails, the event must produce a durable `DecisionRecord` with AC-21 audit linkage.

### 14.4 Audit obligation for aggregated events

Aggregated events are not exempt from audit. Their audit obligation is satisfied differently:

- The aggregation profile declares a minimum audit-evidence contract per aggregation window (e.g., count, scope distribution, actor distribution, time range, sample detail);
- Aggregated audit evidence must be queryable, attributable, and retention-governed;
- Aggregated evidence must not suppress any evidence that the applicable profile classifies as mandatory;
- A reconciliation mechanism must detect and escalate anomalies (volume spikes, denied-but-aggregated errors, missing windows).

**Implementation scope constraint:** Within FEATURE-0013, "queryable" and "reconciliation" are contract semantics only. They define the schemas, validation rules, and pure/in-memory conformance-fixture coverage that a conforming aggregation profile must satisfy. They must not generate implementation tasks for a cache, query service, persistence layer, reconciliation worker, background runtime, or any other production service. Design and tasks for this section are limited to: schema definitions, validation contracts, and deterministic pure-function or in-memory conformance fixtures that verify profile declarations against the contract.

### 14.5 Conformance fixture requirements

Conformance fixtures must test:

- A low-risk check that satisfies all six criteria produces no individual `DecisionRecord` and is covered by aggregation audit evidence;
- A low-risk check where criterion 4 fails (e.g., cross-scope flag) must produce a durable `DecisionRecord` with AC-21 audit linkage;
- An aggregation window that exceeds profile bounds must escalate;
- Aggregation must not suppress per-profile mandatory audit evidence;
- AI-disabled deployments preserve the same boundary behavior.

### 14.6 Classification

This section is classified as `CONTRACT_NOW`. The boundary semantics, distinguishing criteria, and conformance fixture requirements are normative for FEATURE-0013 contract definition.

## 15. Design Questions to Resolve in design.md

- DQ-01: Exact `DecisionRecord` Go struct layout and field types. Should `typedResult` use `json.RawMessage` or a registered interface?
- DQ-02: DecisionProfile schema representation — Go struct with embedded validation rules, or a standalone schema document format?
- DQ-03: Bounded graph representation — adjacency list, edge list, or embedded struct? Design must define the representation and validation of profile-declared finite bounds (nodes, depth, width, fan-out, payload) without selecting numeric defaults. Each profile declares its own finite limits as a required field. If a baseline or system-wide default maximum is needed (e.g., an upper ceiling that no profile may exceed), that value must be escalated as `ARCHITECTURE_DECISION_REQUIRED` rather than invented by design, tasks, fixtures, or code.
- DQ-04: DEFERRED; RESOLVED AS ESCALATION (ADR-F13-002, ARCHITECTURE_DECISION_REQUIRED). The choice of canonicalization algorithm/profile (for example RFC 8785 JCS or an alternative), the exact digest-covered field set, the signature algorithm, and the cryptographic product/service selection are not supplied by any approved architecture source (ADH-2026-014 approves only algorithm-agility and prohibits product selection). Per section 11.8, AC-45.1, and CB-AC-11, design, tasks, fixtures, and code must not select any of those four values; the design must escalate ADR-F13-002 and treat each as a labelled placeholder (for example `PLACEHOLDER_PENDING_ADR_F13_002`) until an approved decision (ADH/DEC/RFC or Architecture Owner directive) supplies it. Design may define only the algorithm-agile carrier fields (canonicalization-method identifier, digest-algorithm identifier, covered-field descriptor, signature-algorithm identifier) and structural trust metadata (CONTRACT_NOW per AC-45/CB-AC-10) without choosing their values; actual cryptographic verification and services remain DEFERRED.
- DQ-05: RESOLVED (ADH-2026-015, refined by the ADH-2026-015 review-epoch correction). `AuditEvent` has no separate scope field or scope enum. `metadata.scopeRef` is the sole logical and serialized scope authority for all seven `ScopeKind` values (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance`). The canonical `Platform` form is an absent/nil `metadata.scopeRef` per FEATURE-0012 `NormalizeScope` (`internal/apimeta/scope.go:84-92`, test `TestNormalizeScope`) and `CanonicalScopeIdentity` (`internal/apimeta/scope.go:107-112`, test `TestCanonicalScopeIdentity`); an absent `scopeRef` resolves deterministically to `Platform` where the contract permits it and returns the stable required-scope error otherwise. Design must not introduce a parallel scope field, presence flag, enum, discriminator, or alias, and must not invent behavior incompatible with the inherited FEATURE-0012 normalization/canonical-identity functions; it must validate scope solely through `metadata.scopeRef`. Remaining design work is limited to the Go struct layout that references `metadata.scopeRef`, not scope-source discrimination.
- DQ-06: Composition strategy registration format — static Go registry map or a declarative file format?
- DQ-07: How to represent the atomic decision/audit-obligation acceptance in pure-function conformance fixtures without an actual persistence layer.
- DQ-08: Conformance fixture format — Go table-driven tests, YAML test cases, or both?
- DQ-09: Relationship between FEATURE-0013 schemas and FEATURE-0012 common resource types — embedding, composition, or parallel packages?
- DQ-10: Profile bundle signed-format for offline verification — structure and trust-root identity representation.
- DQ-11: Error codes for profile rejection, graph limit violation, obligation failure, scope mismatch, and chain-depth limit. Alignment with FEATURE-0012 stable error code patterns.
- DQ-12: Package structure — separate packages for decision, evaluation, audit, composition, profile, or a single `internal/decision/` package?
- DQ-13: EvaluationResult governed-projection schema — how should the `authorizedOutputProjection` field be structured to enforce AC-13.2 bounds (max field count, max total bytes, allowed value types, prohibited content categories)? Should bounds be declared per-profile, per-evaluator-type, or both?
- DQ-14: Cross-scope safety validation for evaluator output (AC-13.5) — should the validation be a separate validation pass, a schema-embedded constraint, or a composition-time check? How is scope identity propagated into the evaluation result acceptance gate?

## 16. Matrix E Risk Summary

Architecture risk governance per ADH-2026-014 and canonical architecture section 18.

- Total risks: 31 (F13-R01 through F13-R31)
- Architecture-stage treatment: APPROVE_TREATMENT (2026-07-27, Sanjeev Kumar)
- Explicitly open High target-residual: F13-R04, F13-R08, F13-R13
- Final residual-risk acceptance: PENDING_HUMAN_REVIEW (all risks)

Risk IDs, controls, verification plans, ownership, and reassessment triggers are preserved from the canonical architecture. Kiro must not renumber, merge, delete, downgrade, transfer, or accept risks. Final residual acceptance requires implementation evidence and human semantic review.

## 17. ADH-2026-015 Scope Vocabulary Clarification and Reconciliation

### 17.1 Authority and status

For the scope-vocabulary decisions in this section, ADH-2026-014 and ADH-2026-015 are the controlling pair, and ADH-2026-015 supersedes only the earlier six-scope `AuditEvent` statement. ADH-2026-016 is joint controlling for FEATURE-0013 overall; it preserves ADH-2026-015 as authoritative for the seven-value scope vocabulary while clarifying the additional contract boundaries recorded in section 18. This section remains authoritative for scope vocabulary and supersedes the historical six-scope statement retained in section 1.5. AC-24 directly expresses the ADH-2026-015 seven-value contract. All unaffected ADH-2026-014 requirements, scope classifications, and non-goals are preserved.

Controlling architecture: `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md` section 29.
Compatibility evidence: `docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md`.

This section adds requirements only. Continued design, tasks, and implementation remain unauthorized until this requirements document receives fresh independent human approval. This document does not self-approve.

### 17.2 Scope vocabulary acceptance criteria (CONTRACT_NOW)

- SV-AC-01: `ScopeKind` is the single canonical shared scope vocabulary. Its exact values are exactly seven: `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, and `ServiceInstance`. No eighth value and no fewer than these seven are permitted. (CONTRACT_NOW)
- SV-AC-02: `AuditEvent` accepts all seven canonical `ScopeKind` values. A conforming `AuditEvent` fixture must exist for each of the seven values. The positive `Platform` fixture uses the canonical absent/nil `metadata.scopeRef` form (per FEATURE-0012 `NormalizeScope`/`CanonicalScopeIdentity`); the six positive non-Platform fixtures each carry a non-nil `metadata.scopeRef` with a canonical kind and UID. (CONTRACT_NOW)
- SV-AC-03: `AuditEvent` uses `metadata.scopeRef` as its sole logical and serialized scope authority. Validation must confirm scope is expressed only through `metadata.scopeRef`. An absent/nil `metadata.scopeRef` is the canonical `Platform` form and resolves deterministically to `Platform` where the contract permits `Platform`; where the contract does not permit `Platform`, an absent `metadata.scopeRef` returns the stable required-scope error (absence is not automatically valid). (CONTRACT_NOW)
- SV-AC-04: No parallel `AuditEvent` scope enum, presence flag, `AuditScope` type, local scope enum, discriminator, alias, or independently mutable scope field may exist in any schema, Go binding, or validator. Conformance must include a negative check asserting the absence of a second scope source. The canonical absent/nil `metadata.scopeRef` `Platform` form is not a second scope source and must not be flagged by this check. (CONTRACT_NOW)
- SV-AC-05: Backward compatibility — every one of the six existing FEATURE-0012 `ScopeKind` values (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`) remains valid and unchanged. Regression fixtures must cover each of the six values. (CONTRACT_NOW)
- SV-AC-06: Existing Organization-scoped `AuditEvent` records remain valid. A regression fixture must exercise an Organization-scoped `AuditEvent` and pass. (CONTRACT_NOW)
- SV-AC-07: `ServiceInstance` positive fixtures — at least one positive fixture must exercise `ServiceInstance` as a valid `ScopeKind` value for `AuditEvent`. (CONTRACT_NOW)
- SV-AC-08: Contradictory or duplicate scope rejection — an `AuditEvent` that declares scope through any source other than `metadata.scopeRef`, or through more than one scope source, or with conflicting scope declarations, must be rejected. The contradictory/duplicate condition is the presence of an unauthorized second scope property or source; a canonical absent/nil `metadata.scopeRef` resolving to `Platform` is explicitly not contradictory and must not be rejected on this basis. A negative fixture must cover an unauthorized second scope source, and a positive fixture must confirm the canonical nil-`Platform` form is accepted. (CONTRACT_NOW)
- SV-AC-09: A scope value outside the canonical seven must be rejected by validators. A negative fixture must cover an out-of-vocabulary scope kind. (CONTRACT_NOW)

### 17.3 Dependency-reconciliation acceptance criteria (CONTRACT_NOW)

- SV-AC-10: Before design approval, a value-by-value FEATURE-0012 versus FEATURE-0013 `ScopeKind` compatibility table must exist as evidence, classifying every inherited value as added, removed, renamed, mapped, restricted, expanded, or semantically reinterpreted. The evidence is recorded in `docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md`. (CONTRACT_NOW — gates design approval)
- SV-AC-11: The Dependency Contract Reconciliation must address every dimension defined in architecture section 29.5: exact dependency version/files, enums, required fields, cardinality, ownership, serialization, validation, security semantics, Go bindings, fixtures, migration classification, and approval reference. For the `metadata.scopeRef` scope-identity dimension, the reconciliation must record it as the **sole** scope source that is **conditionally present**, not universally required: absent/nil for the canonical `Platform` form where the applicable contract permits `Platform`; non-nil and required for every non-Platform scope; and rejected with the stable required-scope error when absent under a contract that does not permit `Platform`. The reconciliation must not describe `metadata.scopeRef` as universally or unconditionally required. An incomplete reconciliation blocks design approval. (CONTRACT_NOW — gates design approval)
- SV-AC-12: The compatibility evidence must cite the exact FEATURE-0012 schema, Go-binding, validator, and fixture files inspected, and must record the six existing values, the additive `ServiceInstance` value, serialization compatibility, migration classification, affected consumers, and regression obligations. (CONTRACT_NOW — gates design approval)
- SV-AC-13: Any inherited dependency delta not explicitly approved by ADH-2026-014 or ADH-2026-015 must be recorded as `ARCHITECTURE_DECISION_REQUIRED` and escalated. Kiro and Cursor must not resolve unapproved deltas or choose mappings, aliases, removals, additional scopes, restrictions, or applicability exclusions. (CONTRACT_NOW)
- SV-AC-14: The compatibility evidence must not claim that schemas, Go bindings, validators, fixtures, or tests have already been updated or passed. Its status must remain pending implementation evidence, and design/tasks/implementation remain gated on fresh independent requirements review. (CONTRACT_NOW)

### 17.4 Migration classification

The `ScopeKind` change is an additive, backward-compatible extension: `ServiceInstance` is added; `Provider` is retained; no value is removed, renamed, mapped, restricted, or semantically reinterpreted. The `AuditEvent` permitted-scope set is expanded from Organization-only to all seven canonical values.

### 17.5 Classification and preservation

This section is classified `CONTRACT_NOW`. It preserves all existing scope classifications (section 11) and non-goals (section 5) unchanged. It does not authorize any runtime capability, persistence, or product selection. Synchronization of schemas, bindings, validators, and fixtures is an obligation for a later, separately authorized implementation stage.

## 18. ADH-2026-016 Contract Boundary Clarification

### 18.1 Authority and status

ADH-2026-014, ADH-2026-015, and ADH-2026-016 are the joint controlling handoffs for this feature. Their precise relationships are:

- ADH-2026-015 supersedes **only** the earlier six-scope `AuditEvent` statement in ADH-2026-014; all other ADH-2026-014 decisions remain controlling.
- ADH-2026-016 is a **clarification** that introduces no new architecture and clarifies previously unresolved contract boundaries of both ADH-2026-014 and ADH-2026-015. The seven-value `ScopeKind` vocabulary approved by ADH-2026-015 (section 17, SV-AC-01) is preserved unchanged.

ADH-2026-016 authorizes no runtime capability, persistence, or product/provider/algorithm selection. It returned this requirements document to architecture-remediation status. This document keeps `review_status: PENDING_HUMAN_REVIEW` and gate status `PENDING_INDEPENDENT_REVIEW` (section 12.8); it does **not** self-approve or advance any stage. Design, tasks, and implementation remain unauthorized until fresh independent human approval is granted.

Controlling architecture: `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md` section 30 (with sections 6.1, 12.1, 12.3, 12.4, and 17).

This section adds and restates requirements only. It preserves all prior valid corrections, the scope classifications (section 11), non-goals (section 5), Matrix E identifiers (section 16), the seven-value scope contract and its compatibility gates (section 17), and the non-circular authorization gate structure (section 12.7, 12.8) unchanged.

### 18.2 DecisionRecord scope authority (CONTRACT_NOW)

- CB-AC-01: `DecisionRecord` uses FEATURE-0012 `metadata.scopeRef` as its sole logical and serialized scope authority; no top-level `DecisionRecord` `scopeRef` or second scope source exists, and any such source is prohibited and fails validation (see AC-01.1). (CONTRACT_NOW)
- CB-AC-02: `Platform` uses an absent/nil `metadata.scopeRef` where `Platform` is permitted; every non-Platform scope uses a non-nil `metadata.scopeRef` with canonical `ScopeKind` and UID semantics. A canonical nil `metadata.scopeRef` resolving to `Platform` is not a contradictory or duplicate scope. (CONTRACT_NOW)
- CB-AC-03: Duplicate, conflicting, alternate, aliased, or parallel scope sources fail validation. The seven-value `ScopeKind` vocabulary (SV-AC-01) is preserved. (CONTRACT_NOW)

### 18.3 Provider-neutral sensitivity vocabulary (CONTRACT_NOW)

- CB-AC-04: The sensitivity classification core vocabulary is closed and ordered: `PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED`. No fifth core value and no fewer than these four are permitted. (CONTRACT_NOW)
- CB-AC-05: Jurisdiction-specific classification labels and caveats are versioned profile mappings onto the core vocabulary, not additional core enum values. Jurisdiction-specific mappings are preserved as profile data. (CONTRACT_NOW)
- CB-AC-06: A profile permitting captured evaluator output (or other captured classified content) must declare all of: a sensitivity ceiling, allowed/prohibited content-category rules, a maximum field count, a maximum byte size, a maximum nesting depth, and the allowed value types. This is the sensitivity-vocabulary application of the AC-13.2 projection bounds. (CONTRACT_NOW)
- CB-AC-07: Missing or unknown mandatory profile values invalidate the profile; a captured value classified above the declared ceiling fails validation. (CONTRACT_NOW)

### 18.4 Security validation boundary (CONTRACT_NOW)

- CB-AC-08: FEATURE-0013 security validation is bounded to structural schema validation, typed classification, declared bounds, projection restrictions, and deterministic conformance (see SEC-15). (CONTRACT_NOW)
- CB-AC-09: FEATURE-0013 does not claim comprehensive semantic detection of secrets, credentials, personal data, malware, or policy violations in arbitrary content. Comprehensive semantic content scanning and runtime content inspection belong to later approved security features. No secret-pattern, credential-pattern, PII, malware, DLP, or content-scanning engine may be introduced in FEATURE-0013. Producers must redact prohibited content before submission. (CONTRACT_NOW / boundary)

### 18.5 Structural trust versus deferred cryptographic boundary (mixed)

- CB-AC-10 (CONTRACT_NOW): Algorithm-agile carrier fields and structural trust metadata (canonicalization-method identifier, digest-algorithm identifier, covered-field descriptor, signature-algorithm identifier, and structural trust-state fields) are `CONTRACT_NOW` carrier data (see AC-45). Structural fixtures may use opaque deterministic assertions and must not claim cryptographic validity or select an algorithm; RFC 8785 remains illustrative only. When a profile requires verified trust, an unknown, absent, expired, revoked, mismatched, or unverified trust state fails closed.
- CB-AC-11 (DEFERRED — ADR-F13-002): The canonicalization algorithm/profile, digest-covered fields, signature algorithm, and cryptographic product/service selection remain DEFERRED and blocked by ADR-F13-002 (see AC-45.1 and the section 11.9 register). Actual cryptographic verification, signing, HSM, notary, trusted-time, key, revocation-distribution, and WORM services are not FEATURE-0013 artifacts.

### 18.5.1 Pre-ADR-F13-002 digest and trust conformance boundary

This subsection is the single normative pre-ADR-F13-002 digest and trust conformance boundary. It reconciles AC-13.3, AC-13.7(f)–(g), SEC-10, SEC-11, AC-45, AC-45.1, CB-AC-10, and CB-AC-11 with the digest/trust-related conformance scenarios F13-CF-13, F13-CF-23, F13-CF-27, and F13-CF-28. It authorizes no canonicalization, digest, signature, or cryptographic algorithm selection; it is a clarification of the already approved algorithm-agile carrier and structural-trust boundary.

- CB-AC-16 (CONTRACT_NOW — allowed structural opaque checks): Until ADR-F13-002 is approved (section 11.9 register), executable conformance fixtures and code covering the criteria above may validate only the following structural, opaque-state checks:
  - (a) **carrier presence** — that the required algorithm-agile carrier fields (canonicalization-method identifier, digest-algorithm identifier, covered-field descriptor, signature-algorithm identifier, and structural trust-state fields) are present or absent exactly as the schema requires;
  - (b) **identifier syntax** — that those identifier fields conform to declared structural, type, and format rules, without selecting, endorsing, or requiring any specific algorithm value;
  - (c) **covered-field-descriptor presence** — that a covered-field descriptor is present and structurally well-formed, without selecting or computing the set of fields any digest actually covers;
  - (d) **opaque expected-versus-observed state** — deterministic comparison of an externally supplied opaque expected digest/trust-state token against an externally supplied opaque observed token as opaque values only, with no computation of either token from content;
  - (e) **fail-closed trust-state transitions** — that an unknown, absent, expired, revoked, mismatched, or unverified trust state fails closed when a profile requires verified trust.
- CB-AC-17 (DEFERRED — ADR-F13-002 — prohibited until resolved): Until ADR-F13-002 is approved, fixtures, code, design, and tasks must not: execute canonical serialization; compute a digest over content (content-to-digest computation); generate or verify a signature; select a canonicalization algorithm/profile, a digest-covered field set, a signature algorithm, or a cryptographic product/service; or claim cryptographic proof or cryptographic validity of any digest, signature, or trust state. A digest carrier is treated solely as externally supplied opaque structural metadata; no fixture may assert that such a carrier cryptographically proves that content existed, was considered, or is unaltered. The opaque expected-versus-observed and fail-closed checks in CB-AC-16 are structural state checks only and constitute no cryptographic-validity claim. RFC 8785 and any specific digest or signature algorithm remain illustrative only and are not selected.

### 18.6 Exact conformance coverage (CONTRACT_NOW)

- CB-AC-12: The architecture section 17 scenarios map one-to-one to stable IDs `F13-CF-01` through `F13-CF-28` in their existing order, without merging, renumbering, or omission (see AC-50). (CONTRACT_NOW)
- CB-AC-13: Additional AC-13.7, scope, security, and compatibility cases use separately named IDs and do not alter the canonical 28-scenario count. Coverage is counted by scenario ID, not by fixture-file count. One fixture may cover multiple scenario IDs, but every scenario ID requires an explicit coverage-matrix entry. Contract-only fixtures remain pure/in-memory. (CONTRACT_NOW)

### 18.7 Preservation

- CB-AC-14: AC-45.1 is preserved as DEFERRED with ADR-F13-002 as a separate `ARCHITECTURE_DECISION_REQUIRED` escalation marker (section 11.9). (Preserved)
- CB-AC-15: All prior valid corrections, scope classifications (section 11), non-goals (section 5), Matrix E identifiers (section 16), the seven-value scope contract and compatibility gates (section 17, AC-54, BC-AC-04, BC-AC-05), and the non-circular authorization gate structure (sections 12.7 and 12.8) are preserved unchanged. No production decision/audit/security/cryptographic/profile-registry/synchronization service, provider-specific component, parallel scope field, second decision type, adapter interface, queue, worker, persistence layer, or product selection is authorized. (Preserved)

### 18.8 Classification

This section is classified `CONTRACT_NOW` except for CB-AC-11 and CB-AC-17 (both DEFERRED under ADR-F13-002). It authorizes no runtime capability, persistence, or product selection. Synchronization of schemas, bindings, validators, and fixtures remains an obligation for a later, separately authorized implementation stage, and continued design/tasks/implementation remain gated on fresh independent requirements review.

---

```text
Model Execution Report:
- Tool: kiro
- Stage or task: requirements
- Recommended priority list: Architecture-heavy (claude-opus-4.8), Fallback (claude-sonnet-4.5)
- Selected model: claude-sonnet-4.5
- Effort/reasoning setting: not directly visible
- Fallback used: yes
- Fallback reason: unavailable
```
