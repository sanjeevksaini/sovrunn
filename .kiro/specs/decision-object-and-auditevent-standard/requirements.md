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
controlling_architecture: docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md
reuse_assessment_format_version: 1.0.0
depends_on:
  - FEATURE-0012
---

# FEATURE-0013 — Decision Record and AuditEvent Standard

Stage: Requirements

> **ADH-2026-015 supersession notice (Approved 2026-07-27).**
> ADH-2026-014 and ADH-2026-015 are the joint controlling handoffs for this
> feature. ADH-2026-015 supersedes **only** the earlier six-scope `AuditEvent`
> statement. Under ADH-2026-015, `ScopeKind` is the single canonical shared
> vocabulary with exactly seven values — `Platform`, `Organization`,
> `OrganizationUnit`, `Tenant`, `Project`, `Provider`, and `ServiceInstance` —
> `ServiceInstance` is additive, `Provider` is retained, and `AuditEvent`
> permits all seven values using `metadata.scopeRef` as its sole scope
> identity. **Section 17 is authoritative** for scope vocabulary. Where earlier
> text in this document (for example sections 1.5, 8.2, and AC-24) refers to
> "six" scopes or omits `Provider`, section 17 controls.
>
> The previous `APPROVED_FOR_DESIGN: READY_FOR_REVIEWER` assertion recorded in
> section 12 is **historically preserved as evidence but is superseded for
> continued design work** by ADH-2026-015. This requirements document is
> returned to `PENDING_HUMAN_REVIEW`. Requirements must receive **fresh
> independent approval** before design, tasks, or implementation resume. This
> document does not self-approve.

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

- ADH-2026-014 (Approved, 2026-07-27) and ADH-2026-015 (Approved, 2026-07-27) — joint controlling handoffs. ADH-2026-015 supersedes only the earlier six-scope `AuditEvent` statement; all other ADH-2026-014 decisions remain controlling.
- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
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

- AC-01: A `DecisionRecord` schema defines: profile reference (name + version), form, authority, purpose, request reference, actor reference, subject references, scope reference, effective policy context reference, input snapshot reference, evaluation result references, composition metadata, typed result, rationale, obligations, provenance, validity, correlation (operation + audit event references).
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
  - AC-13.3: Content classified as protected, restricted, or unrestricted-tenant-owned must not be copied into the canonical `EvaluationResult` record. Such content is represented by: (a) a controlled reference (stable URI or identifier resolvable within the evaluator's trust boundary), and (b) an integrity digest proving the content existed and was considered, without disclosing it.
  - AC-13.4: Raw prompts, model system instructions, full model responses exceeding the projection bound, secrets, credentials, provider tokens, and unauthorized tenant content must never appear in the canonical record. Redaction is mandatory before persistence.
  - AC-13.5: Cross-scope safety validation: before an `EvaluationResult` is accepted into a `DecisionRecord`, the projection must be validated against the target scope's sensitivity ceiling. An evaluation produced in Tenant A's scope must not project content that discloses Tenant B's existence, data, or resource identity.
  - AC-13.6: Replay must not rerun the evaluator. The captured projection, digest, and timing envelope are sufficient for audit, explanation, and deterministic composition without re-execution.
  - AC-13.7: Conformance fixtures must include positive and negative cases proving:
    - (a) A conforming bounded projection within size/field/type limits is accepted;
    - (b) A projection exceeding the size bound is rejected;
    - (c) A projection containing a secret or credential pattern is rejected;
    - (d) A projection containing a raw prompt or model system instruction is rejected;
    - (e) A projection containing unrestricted tenant content from another scope is rejected;
    - (f) A protected-content reference with valid integrity digest is accepted in place of the raw content;
    - (g) A protected-content reference with mismatched or missing digest is rejected;
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

### 4.13 Security and trust (CONTRACT_NOW)

- AC-41: Authenticate and authorize every human, workload, evaluator, and service hop. Bind authorization to scope, purpose, action, and resource.
- AC-42: Evaluator output, reason text, external evidence, and AI content are treated as untrusted data. Schema, size, depth, character, URI, and reference-kind limits are enforced.
- AC-43: Canonical governance records must not duplicate personal data, secrets, credentials, provider tokens, raw policy documents, model prompts, or unrestricted tenant content. Records prefer stable references, digests, classifications, and purpose-bound projections.
- AC-44: Different authorized projections serve customer, operator, auditor, regulator, provider, and AI consumers. Redaction is field- and reason-aware.
- AC-45: Canonical serialization and content digests are supported. Deployment profiles may require signatures, trusted timestamps, append-only/WORM retention, or external notarization. The contract remains algorithm-agile.

### 4.14 AI governance invariants (INVARIANT_FOR_LATER)

AC-46 through AC-49 constrain the behavior of a later AI consumer and AI runtime. The `DecisionContext` schema, the AI consumer, prompt/model integration, the model client, the AI runtime, and any production projection service are owned by FEATURE-0025 or the later owning feature. FEATURE-0013 does **not** define or implement them. FEATURE-0013 owns only the `DecisionProfile` projection constraints, authorized-projection semantics, references, labels, and non-executing conformance data that `DecisionRecord` and `AuditEvent` require; those genuinely FEATURE-0013-owned contract elements remain CONTRACT_NOW and are carried by AC-08 (DecisionProfile schema) and the projection semantics in SEC-09 and SEC-13. AC-46 through AC-49 themselves generate no FEATURE-0013 runtime or schema implementation task beyond documentation and non-executing conformance placeholders.

- AC-46: (INVARIANT_FOR_LATER) A later AI consumer receives only an authorized `DecisionContext` projection. It may explain, summarize, identify anomalies, and suggest remediation through normal approval. The authorized-projection constraint fields that a `DecisionProfile` declares are CONTRACT_NOW (AC-08); the `DecisionContext` schema and AI consumer are owned by FEATURE-0025.
- AC-47: (INVARIANT_FOR_LATER) A later AI consumer must not mutate records, invent evidence, suppress audit, widen disclosure, bypass approval, or execute a suggestion. This constrains later AI-runtime behavior, not a FEATURE-0013 artifact.
- AC-48: (INVARIANT_FOR_LATER) AI-generated recommendations are labelled with model/provider/version, prompt-template version, data boundary, confidence, and human-review status. The label/field schema a `DecisionProfile` or `DecisionRecord` may carry is CONTRACT_NOW (AC-08); producing the labelled recommendation is later AI-runtime behavior owned by FEATURE-0025.
- AC-49: (INVARIANT_FOR_LATER) Core decisions remain available and correct when AI is disabled or disconnected. This is a runtime invariant a later implementation must satisfy; FEATURE-0013 proves it only through non-executing conformance placeholders (EC-09).

### 4.15 Conformance fixtures (CONTRACT_NOW)

- AC-50: Executable positive and negative fixtures are defined for at least the 28 scenarios listed in the canonical architecture section 17 (simplest atomic authorization, deterministic denial, pure selection, all forms, human-approval async, maximally bounded composite, parallel evaluations, timeout/failure/conflict, idempotent replay, partial persistence failure, supersession/correction/revocation, unauthorized rejection, offline/air-gapped, export/import, AI-disabled, graph limit rejection, new-family registration, local atomic acceptance, idempotency vs. semantic identity, deterministic replay, local-bundle operation, profile rejection, obligation denial, authorization cache, composite budgets, local-digest with batch signing, disconnected import) plus the evaluator-output governed-projection scenarios required by AC-13.7.

**AC-50 implementation scope constraint:** The "authorization cache" scenario fixture tests the contract semantics of a cached authorization decision (validity, expiry, staleness detection, refresh obligation, scope binding) using pure-function or in-memory conformance fixtures only. It must not generate implementation tasks for a runtime authorization cache, cache-invalidation service, distributed cache, TTL worker, persistence layer, or any background process. The fixture validates that a conforming cache implementation would satisfy the contract, without itself being or requiring that cache.

**AC-50 evaluator-output governed-projection fixture constraint:** The evaluator-output scenarios required by AC-13.7 test the contract semantics of governed capture (projection bounds, redaction, sensitivity classification, cross-scope safety, reference/digest handling) using pure-function or in-memory conformance fixtures only. They must not generate implementation tasks for an evaluator runtime, AI model client, prompt engine, redaction service, or any background process. The fixtures validate that a conforming evaluator adapter and record-acceptance gate would satisfy the AC-13 data-minimization contract, without themselves being or requiring those services.

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
- EC-09: AI-disabled deployments must function correctly for all decision and audit behavior. AI projections return empty/unavailable rather than errors.
- EC-10: Cross-scope references that would disclose existence of resources in another tenant or governance domain must be safe-denied without existence disclosure.

## 7. Security and Privacy Requirements

- SEC-01: Authenticate and authorize every human, workload, evaluator, and service hop at a zero-trust boundary.
- SEC-02: Bind authorization to scope, purpose, action, and resource—not network location alone.
- SEC-03: Record delegation and impersonation explicitly in actor references.
- SEC-04: Reject unsigned or untrusted evaluator/strategy definitions where signing is required by the deployment profile.
- SEC-05: Treat evaluator output, reason text, external evidence, and AI content as untrusted data requiring validation.
- SEC-06: Enforce schema, size, depth, character, URI, and reference-kind limits on all inputs to prevent injection and resource exhaustion.
- SEC-07: The canonical governance record must not indiscriminately duplicate personal data, secrets, credentials, provider tokens, raw policy documents, model prompts, or unrestricted tenant content.
- SEC-08: Records prefer stable references, digests, classifications, and purpose-bound projections.
- SEC-09: Confidentiality: different authorized projections serve different consumers. Redaction is field- and reason-aware and must not change the final outcome or falsely imply omitted evidence did not exist.
- SEC-10: Integrity: canonical serialization and content digests. Algorithm-agile. Deployment profiles may require additional assurance (signatures, timestamps, WORM, Merkle proofs, notarization).
- SEC-11: Local digest computation is permitted on every record. Per-record remote HSM or notarization is not a universal hot-path requirement. Tiered assurance by profile.
- SEC-12: For security, authorization, sovereignty, compliance, privileged access, data movement, and break-glass profiles, fail-closed or `REQUIRES_APPROVAL` is mandatory. Fail-open is prohibited unless a bounded approved exception identifies scope, owner, compensating controls, expiry, audit treatment, and reassessment trigger.
- SEC-13: AI prompts, logs, traces, and model calls must not expose secrets or unauthorized tenant data. Redaction and tenant boundary enforcement are mandatory.
- SEC-14: Cross-scope references must not disclose or influence another tenant or governance domain.

## 8. Compatibility with Completed Phase 1 Features

### 8.1 FEATURE-0005 Operation Resource

- `Operation` continues to own asynchronous lifecycle (progress, retries, deadlines, cancellation, status). `DecisionRecord` links to `Operation` via `correlation.operationRef` but never becomes a mutable workflow-state object.
- No change to the Phase 1 Operation resource shape is required. Decision-linked operations use the existing correlation and reference patterns.

### 8.2 FEATURE-0001 through FEATURE-0004 Organization hierarchy

- `DecisionRecord` scope references use the same Organization/OrganizationUnit/Tenant/Project hierarchy defined in Phase 1.
- `AuditEvent` scope expansion (ADH-2026-015) permits all seven canonical `ScopeKind` values via `metadata.scopeRef`. `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, and `Project` map to the Phase 1 hierarchy; `Provider` is retained from the FEATURE-0012 shared vocabulary; `ServiceInstance` is additive and extends naturally from FEATURE-0008.
- No Phase 1 resource shapes are modified.

### 8.3 FEATURE-0006 ServiceClass/ServicePlan and FEATURE-0008 ServiceInstance

- `DecisionRecord` may reference a ServiceInstance or ServicePlan as a subject or scope. These are typed references, not embedded duplicates.
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
| Canonical serialization and digest | Reuse | Reuses RFC 8785 JCS concepts for canonicalization profile; exact covered fields decided in design. | Approved | ADH-2026-014 |
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
- Reused or external responsibility: DMN composition concepts, OpenTelemetry correlation, RFC 8785 canonicalization, FEATURE-0012 grammar
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
| CONTRACT_NOW | AC-01 through AC-25, AC-36 through AC-45, AC-50 through AC-58 |
| INVARIANT_FOR_LATER | AC-26 through AC-35, AC-46 through AC-49, AC-59 through AC-61 |
| DEFERRED | Production persistence pattern, workflow/orchestration runtime, evidence store, identity provider, key manager, jurisdiction-specific values, provider adapters, AI runtime, canonicalization/signature profiles, cross-site topology |

### 11.3 Security requirements classification

| Requirement | Classification | Rationale |
|---|---|---|
| SEC-01 | CONTRACT_NOW | Zero-trust authentication/authorization contract semantics defined in schemas and conformance fixtures |
| SEC-02 | CONTRACT_NOW | Authorization binding semantics defined in DecisionProfile and typed references |
| SEC-03 | CONTRACT_NOW | Delegation/impersonation representation defined in actor reference schema |
| SEC-04 | CONTRACT_NOW | Signing/trust requirements defined in profile contract and bundle format |
| SEC-05 | CONTRACT_NOW | Untrusted-data validation rules defined in schema and conformance fixtures |
| SEC-06 | CONTRACT_NOW | Schema/size/depth/character/URI/reference-kind limits defined in validation contract |
| SEC-07 | CONTRACT_NOW | Data-minimization contract defined in record schemas and projection rules |
| SEC-08 | CONTRACT_NOW | Reference/digest/classification preference defined in record schema contract |
| SEC-09 | CONTRACT_NOW | Authorized-projection contract defined in projection profile semantics |
| SEC-10 | CONTRACT_NOW | Canonical serialization and content-digest contract; algorithm-agility semantics |
| SEC-11 | CONTRACT_NOW | Tiered-assurance contract defined in profile semantics and conformance fixtures |
| SEC-12 | CONTRACT_NOW | Fail-closed/fail-open semantics and bounded-exception contract in profile rules |
| SEC-13 | CONTRACT_NOW | AI prompt/trace/model redaction and boundary rules defined in projection contract |
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
| EC-09 | CONTRACT_NOW | AI-disabled behavior contract defined in projection semantics and conformance fixtures |
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

## 12. ADH-2026-014 Baseline Correction Prerequisites

### 12.1 Purpose

ADH-2026-014 approves two controlled baseline corrections. These corrections require coordinated repository updates. This section separates documentary prerequisites (which gate design authorization) from implementation and fixture prerequisites (which gate APPROVED_FOR_CURSOR and final acceptance).

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
- PRE-06: Compatibility fixtures for all seven canonical `ScopeKind` values, plus separate regression fixtures for the six pre-existing FEATURE-0012 `ScopeKind` values, must be specified as acceptance evidence (covers AC-54).
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
- PRE-09: If FEATURE-0012 conformance fixtures require update for seven-value `AuditEvent` compatibility, those updates must be completed and pass before FEATURE-0013 implementation (APPROVED_FOR_CURSOR).
- PRE-10: The FEATURE-0012 human semantic review evidence (`docs/reviews/feature-gates/FEATURE-0012-human-semantic-review-evidence.md`) must confirm acceptance of the AuditEvent scope correction or identify required remediation.

### 12.6 Prerequisite verification evidence

This section records evidence collected for each APPROVED_FOR_DESIGN prerequisite as of 2026-07-27. This evidence record is informational only. It does not constitute or grant stage authorization. Independent review must validate DGA-01 through DGA-04 before design authorization is established.

**DGA-01 independent verification status (2026-07-27, re-verified after remediation):**

The traceability record at `docs/traceability/ADH-2026-014-terminology-reconciliation.md` exists and was independently inspected after remediation. The following structural conditions were verified:

- ✅ File exists at the canonical path.
- ✅ "Corrected normative files" table is present with 11 entries.
- ✅ Each entry includes all four required columns: File, Prior content, Correction applied, Rationale.
- ✅ The `FEATURE_TRACEABILITY_MATRIX.md` row now includes its Rationale column: "Traceability matrix must reflect current feature status and approved corrections".
- ✅ "Files intentionally NOT corrected" section is present with 7 entries and each has a disposition rationale.
- ✅ "Verification" section states PRE-01, PRE-02, and PRE-03 as Complete.

**DGA-01 verdict: PASS.**

The traceability record is structurally complete. The previously identified defect (missing Rationale column for the `FEATURE_TRACEABILITY_MATRIX.md` row) has been corrected and independently verified.

| Prerequisite | Collected status | Evidence |
|---|---|---|
| PRE-01 (repository-wide DecisionObject audit) | SATISFIED | `docs/traceability/ADH-2026-014-terminology-reconciliation.md` "Corrected normative files" table lists 11 normative files audited with all four columns complete; "Files intentionally NOT corrected" lists 7 excluded with rationale; "Verification" section states PRE-01 Complete. DGA-01 re-verified: PASS. |
| PRE-02 (each occurrence corrected or aliased) | SATISFIED | Same traceability record confirms all normative uses replaced; historical/non-goal references documented as intentionally retained with explicit rationale per file. DGA-01 re-verified: PASS. |
| PRE-03 (traceability record exists) | SATISFIED | `docs/traceability/ADH-2026-014-terminology-reconciliation.md` exists at the canonical path, lists every changed file, prior term, corrected term, and rationale. All table rows structurally complete. DGA-01 re-verified: PASS. |
| PRE-05 (seven-value AuditEvent scope correction documented) | EVIDENCE_RECORDED | ADH-2026-015 and section 17 document the canonical seven-value ScopeKind and AuditEvent contract, including the canonical absent/nil `Platform` form. The controlling Phase 2 spine, acceptance gates, current baseline, decision summary, RFC, feature index, and traceability matrix identify ADH-2026-014 and ADH-2026-015 as the joint controlling handoffs and record the seven-value contract at the documentary level. Schema, Go-binding, validator, and fixture synchronization remains an obligation pending implementation evidence and is not claimed complete. Historical ADH-2026-014 six-scope evidence remains preserved as superseded compatibility history. |
| PRE-07 (alpha compatibility decision recorded with approval) | EVIDENCE_RECORDED | ADH-2026-015 approves the seven-value `ScopeKind`/`AuditEvent` correction (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance` via `metadata.scopeRef`), superseding the earlier ADH-2026-014 six-scope statement. Human approval: Sanjeev Kumar, 2026-07-27. The superseded ADH-2026-014 approval remains recorded in `docs/reviews/architecture-decision-handoffs/ADH-2026-014-feature-0013-decision-record-and-auditevent-standard.md` and `docs/traceability/ADH-2026-014-terminology-reconciliation.md` as compatibility history. |
| Section 12.4: PHASE2_ARCHITECTURE_SPINE.md | EVIDENCE_RECORDED | Uses DecisionRecord terminology and the canonical seven-value ScopeKind; AuditEvent permits all seven through metadata.scopeRef as sole scope identity. The earlier ADH-2026-014 six-scope statement is explicitly identified as superseded. |
| Section 12.4: CURRENT_ARCHITECTURE_BASELINE.md | EVIDENCE_RECORDED | Records ADH-2026-015 as the latest controlled update, the seven canonical ScopeKind values, all-seven AuditEvent applicability, and metadata.scopeRef as sole scope identity. |
| Section 12.4: CURRENT_DECISION_SUMMARY.md | EVIDENCE_RECORDED | Preserves the superseded ADH-2026-014 six-scope statement as history and records the authoritative ADH-2026-015 seven-value correction. |
| Section 12.4: RFC-0023-decision-and-audit-standard.md | EVIDENCE_RECORDED | Uses DecisionRecord terminology and records ADH-2026-014 and ADH-2026-015 as joint controlling handoffs with the seven-value scope contract. |
| Section 12.4: FEATURE_INDEX.md | EVIDENCE_RECORDED | FEATURE-0013 uses the canonical Decision Record title and identifies ADH-2026-014 and ADH-2026-015 as joint controlling handoffs. |
| Section 12.4: FEATURE_TRACEABILITY_MATRIX.md | EVIDENCE_RECORDED | FEATURE-0013 records the seven-value ScopeKind, all-seven AuditEvent applicability, metadata.scopeRef sole authority, ADH-2026-015 compatibility evidence, and pending fresh requirements review. |
| Section 12.4: PHASE2_ACCEPTANCE_GATES.md | EVIDENCE_RECORDED | FEATURE-0013 gates identify ADH-2026-014 and ADH-2026-015 as joint controlling handoffs (ADH-2026-015 superseding only the earlier six-scope statement) and require the canonical seven-value ScopeKind, all-seven AuditEvent fixtures (canonical absent/nil `Platform`, non-Platform via `metadata.scopeRef`), a negative absence-when-`Platform`-disallowed fixture, six-value FEATURE-0012 regression coverage, and `metadata.scopeRef` as sole scope authority. Schema/binding/validator/fixture synchronization remains pending implementation evidence. |
| Terminology-reconciliation traceability record exists (section 12.2.1) | SATISFIED | `docs/traceability/ADH-2026-014-terminology-reconciliation.md` exists at canonical path. Structure is complete: all 11 rows have four columns; 7 excluded files documented with rationale; PRE-01/02/03 marked Complete. DGA-01 re-verified: PASS. |
| FEATURE-0012 dependency readiness (section 13.4) | EVIDENCE_RECORDED | DEP-01 through DEP-04 all satisfied per section 13.7. |

**Evidence collection status:** Complete. DGA-01 structural defect has been remediated and independently re-verified. PRE-01 through PRE-03 are satisfied. DGA-01 through DGA-04 all pass. Stage authorization awaits independent reviewer validation and return of the APPROVED_FOR_DESIGN token per section 12.8.

### 12.7 Authorization gate structure

FEATURE-0013 authorization is staged:

**APPROVED_FOR_DESIGN requires:**

1. PRE-01 through PRE-03 are satisfied (documentary terminology reconciliation);
2. PRE-05 and PRE-07 are satisfied (documentary seven-value `AuditEvent` scope correction and compatibility decision, ADH-2026-015);
3. All repository updates in section 12.4 are committed and verified;
4. The terminology reconciliation traceability record exists at `docs/traceability/ADH-2026-014-terminology-reconciliation.md`;
5. FEATURE-0012 dependency readiness condition (section 13.4) is confirmed;
6. Independent review validates DGA-01 through DGA-04 (section 12.8) and returns the APPROVED_FOR_DESIGN token.

**APPROVED_FOR_CURSOR and final acceptance additionally require:**

7. PRE-06 is satisfied (executable compatibility fixtures for all seven canonical `ScopeKind` values, plus regression fixtures for the six pre-existing FEATURE-0012 values, exist and pass);
8. PRE-08 through PRE-10 are satisfied (FEATURE-0012 regression evidence);
9. BC-AC-04 and BC-AC-05 are satisfied (executable fixture and regression evidence).

Design and task generation may proceed once the APPROVED_FOR_DESIGN prerequisites are met and the independent review gate validates DGA-01 through DGA-04 and returns the approval token. Implementation artifacts (fixtures, executable conformance, regression evidence) are produced during implementation, not before design authorization.

Failure of any APPROVED_FOR_DESIGN precondition pauses design authorization. Failure of any APPROVED_FOR_CURSOR precondition pauses implementation merge.

### 12.8 Deterministic design-authorization check

Before FEATURE-0013 may progress from requirements to design, the following deterministic checks must pass. These checks are verifiable by any agent or reviewer without subjective judgement.

**Important: This requirements document defines the checks but does not grant its own stage authorization. An independent review gate must validate DGA-01 through DGA-04 and return the APPROVED_FOR_DESIGN token before design generation may proceed.**

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
- Verify: `docs/phase2/PHASE2_ACCEPTANCE_GATES.md` contains FEATURE-0013 gate criteria and identifies ADH-2026-014 and ADH-2026-015 as the joint controlling handoffs, stating that ADH-2026-015 supersedes only the earlier six-scope `AuditEvent` statement.
- Failure action: APPROVED_FOR_DESIGN is blocked until the missing update is applied.

**Check DGA-03: Seven-value `AuditEvent` scope correction (ADH-2026-015) is recorded with human approval.**

- Verify: The ADH-2026-015 architecture-decision handoff exists and records the seven-value `ScopeKind`/`AuditEvent` decision that supersedes the earlier ADH-2026-014 six-scope statement.
- Verify: The controlling architecture (`docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md` section 29) and `docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md` record the seven canonical `ScopeKind` values (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance`) and `metadata.scopeRef` as the sole scope identity.
- Verify: The "Human approval" section records approval status as "Approved", approver identity, and date.
- Verify: The superseded ADH-2026-014 six-scope statement in `docs/reviews/architecture-decision-handoffs/ADH-2026-014-feature-0013-decision-record-and-auditevent-standard.md` is identified as superseded by ADH-2026-015 (compatibility history).
- Failure action: APPROVED_FOR_DESIGN is blocked until the seven-value compatibility decision is approved.

**Check DGA-04: FEATURE-0012 dependency readiness confirmed.**

- Verify: Section 13.7 of this requirements document records DEP-01 through DEP-04 as SATISFIED with evidence.
- Verify: `docs/reviews/feature-gates/FEATURE-0012-approval-review.md` contains "Final feature-review status: Approved".
- Failure action: APPROVED_FOR_DESIGN is blocked until FEATURE-0012 readiness is confirmed.

**Gate result:** If DGA-01, DGA-02, DGA-03, and DGA-04 all pass under independent review, the reviewer may grant APPROVED_FOR_DESIGN. If any check fails, design generation must not proceed and the failing check must be remediated first.

**Current gate status (2026-07-27, re-verified after remediation):**

| Check | Evidence collected | Evidence summary |
|---|---|---|
| DGA-01 | PASS | `docs/traceability/ADH-2026-014-terminology-reconciliation.md` exists with 11 corrected files (all four columns complete including Rationale), 7 excluded files with rationale, and PRE-01/02/03 marked Complete. Structural defect remediated and re-verified. |
| DGA-02 | PASS | All seven section 12.4 documents verified per section 12.6 evidence table. `DecisionRecord` present in normative text; `DecisionObject` only in historical/migration notes. `PHASE2_ACCEPTANCE_GATES.md` identifies ADH-2026-014 and ADH-2026-015 as joint controlling handoffs (ADH-2026-015 superseding only the earlier six-scope statement). |
| DGA-03 | PASS | ADH-2026-015 supersedes the ADH-2026-014 six-scope statement and approves the seven-value `ScopeKind`/`AuditEvent` correction (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance`) via `metadata.scopeRef`; Human approval: Approved, Sanjeev Kumar, 2026-07-27. |
| DGA-04 | PASS | Section 13.7 confirms DEP-01 through DEP-04 SATISFIED; FEATURE-0012-approval-review.md confirms "Final feature-review status: Approved". |

**APPROVED_FOR_DESIGN: READY_FOR_REVIEWER.**

DGA-01 through DGA-04 all pass after remediation and independent re-verification. This requirements document asserts all evidence is present and structurally valid. An independent reviewer must validate DGA-01 through DGA-04 and return the APPROVED_FOR_DESIGN token before design generation may proceed.

### 12.9 Acceptance criteria for baseline corrections

- BC-AC-01: Every normative document listed in section 12.4 contains `DecisionRecord` (not `DecisionObject`) in all normative references to the common decision envelope. (CONTRACT_NOW — gates APPROVED_FOR_DESIGN)
- BC-AC-02: The Phase 2 spine reflects the seven canonical `ScopeKind` values for `AuditEvent` (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance`) per ADH-2026-015. (CONTRACT_NOW — gates APPROVED_FOR_DESIGN)
- BC-AC-03: A terminology reconciliation traceability record exists listing all corrected files and rationale. (CONTRACT_NOW — gates APPROVED_FOR_DESIGN)
- BC-AC-04: `AuditEvent` compatibility fixtures exist for all seven canonical `ScopeKind` values (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance`), and separate regression fixtures cover the six pre-existing FEATURE-0012 values (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`). (CONTRACT_NOW — gates APPROVED_FOR_CURSOR)
- BC-AC-05: FEATURE-0012 conformance evidence demonstrates no regression from the AuditEvent scope correction. (CONTRACT_NOW — gates APPROVED_FOR_CURSOR)
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
- DQ-04: Canonical serialization for content digests — RFC 8785 JCS or alternative? Exact fields covered by the digest.
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

ADH-2026-014 and ADH-2026-015 are the joint controlling handoffs for this feature. ADH-2026-015 supersedes only the earlier six-scope `AuditEvent` statement. This section is authoritative for scope vocabulary and supersedes the historical six-scope statement retained in section 1.5. AC-24 now directly expresses the ADH-2026-015 seven-value contract. All other ADH-2026-014 requirements, scope classifications, and non-goals in this document are preserved unchanged.

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
- SV-AC-11: The Dependency Contract Reconciliation must address every dimension defined in architecture section 29.5: exact dependency version/files, enums, required fields, cardinality, ownership, serialization, validation, security semantics, Go bindings, fixtures, migration classification, and approval reference. An incomplete reconciliation blocks design approval. (CONTRACT_NOW — gates design approval)
- SV-AC-12: The compatibility evidence must cite the exact FEATURE-0012 schema, Go-binding, validator, and fixture files inspected, and must record the six existing values, the additive `ServiceInstance` value, serialization compatibility, migration classification, affected consumers, and regression obligations. (CONTRACT_NOW — gates design approval)
- SV-AC-13: Any inherited dependency delta not explicitly approved by ADH-2026-014 or ADH-2026-015 must be recorded as `ARCHITECTURE_DECISION_REQUIRED` and escalated. Kiro and Cursor must not resolve unapproved deltas or choose mappings, aliases, removals, additional scopes, restrictions, or applicability exclusions. (CONTRACT_NOW)
- SV-AC-14: The compatibility evidence must not claim that schemas, Go bindings, validators, fixtures, or tests have already been updated or passed. Its status must remain pending implementation evidence, and design/tasks/implementation remain gated on fresh independent requirements review. (CONTRACT_NOW)

### 17.4 Migration classification

The `ScopeKind` change is an additive, backward-compatible extension: `ServiceInstance` is added; `Provider` is retained; no value is removed, renamed, mapped, restricted, or semantically reinterpreted. The `AuditEvent` permitted-scope set is expanded from Organization-only to all seven canonical values.

### 17.5 Classification and preservation

This section is classified `CONTRACT_NOW`. It preserves all existing scope classifications (section 11) and non-goals (section 5) unchanged. It does not authorize any runtime capability, persistence, or product selection. Synchronization of schemas, bindings, validators, and fixtures is an obligation for a later, separately authorized implementation stage.

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
