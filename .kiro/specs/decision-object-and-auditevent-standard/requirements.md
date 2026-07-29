---
reuse_assessment_format_version: 1.0.0
---

# Requirements Document

Feature: FEATURE-0013 — Decision Record and AuditEvent Standard
Phase: Phase 2 — Reuse-First PaaS Fabric Foundation
Stage: Requirements

## Metadata

| Field | Value |
|---|---|
| Feature ID | FEATURE-0013 |
| Feature title | Decision Record and AuditEvent Standard |
| Phase | Phase 2 — Reuse-First PaaS Fabric Foundation |
| Spec stage | Requirements |
| Branch | feature-0013-decision-record-and-auditevent-standard |
| Depends on | FEATURE-0012 (API, Resource Naming, Status, and Validation Standard) |
| Controlling handoff | ADH-2026-017 (Approved) |
| Canonical architecture | docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md |
| Canonical reuse standard | docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md |
| Classification | Build / Extend (mixed; see reuse summary) |
| Status | Draft — pending review |

## FEATURE-0013 Reuse Assessment summary

Feature identity: FEATURE-0013 — Decision Record and AuditEvent Standard.

This summary is populated from the approved consolidated architecture
(`docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`)
and ADH-2026-017. Field definitions are owned by
`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` and are not redefined here.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 API/resource grammar | Extend | Reuse every shared profile, metadata, reference, boundary, validation, Problem Details, schema, and compatibility primitive; add only decision/audit domain semantics and approved per-kind constraints. | Approved | ADH-2026-017; ADH-2026-012; RFC-0022 |
| DecisionRecord, DecisionProfile, EvaluationResult, and AuditEvent domain contract | Build | No mature external model supplies Sovrunn's complete sovereign, provider-neutral authority, projection, scope, audit-obligation, and downstream-adoption semantics. The common contract is Sovrunn differentiation; external concepts remain inputs rather than native core types. | Approved | ADH-2026-017; DEC-0026; RFC-0023 |
| Bounded decision-requirements and composition concepts | Reuse | Reuse applicable OMG DMN decision-requirements concepts without selecting, embedding, or recreating a DMN engine or workflow language. | Approved | ADH-2026-017; RFC-0023 |
| Audit transport and operational correlation | Reuse | CloudEvents may be an optional transport mapping and OpenTelemetry may carry operational correlation; neither becomes the canonical record or audit authority. | Approved | ADH-2026-017; RFC-0023 |
| Supply-chain evidence references | Reuse | Reuse DSSE/in-toto and SPDX/CycloneDX concepts by typed reference without duplicating their schemas or implementing signing. | Approved | ADH-2026-017; ADR-F13-002 |
| Canonicalization, digest, signing, and cryptographic verification | Reuse | Preserve algorithm-agile carriers and evaluate mature standards later; no algorithm, covered-field set, product, or service is selected in FEATURE-0013. | Deferred | ADH-2026-017; ADR-F13-002 |
| Policy/evaluator and workflow runtimes | Wrap | CEL, OPA, Cedar, durable-execution, and workflow products are later candidates behind owning-feature adapters or ports; FEATURE-0013 defines no wrapper or runtime implementation. | Deferred | ADH-2026-017; DEC-0036 |

## Introduction

FEATURE-0013 defines the common decision and audit vocabulary for Sovrunn.
It establishes the immutable `DecisionRecord` envelope, the versioned
`DecisionProfile` extension mechanism, the atomic `EvaluationResult`
contract, and the governed `AuditEvent` accountability record.

These contracts enable Sovrunn to represent a trivial allow/deny decision
and a complex, hierarchical, multi-evaluator decision using the same
governed envelope. Authoritative decision composition MUST be deterministic
and reproducible: the authoritative conclusion is a pure function of
captured, versioned inputs — profile name and version,
`EffectivePolicyContext`, the registered composition strategy and version,
evaluator contract versions, the canonical input snapshot, and the validity
epoch — and recomputes identically from those captured inputs. A
nondeterministic or AI evaluator MAY itself be non-reproducible; its output
is captured once, and the authoritative decision is composed and replayed
from that captured output without rerunning the evaluator. Determinism is
therefore a property of the authoritative decision and its composition, not a
claim that every evaluator is internally deterministic. Captured
nondeterministic or AI evaluator output remains permitted evidence and never
becomes authoritative. Every authoritative decision remains explainable,
auditable, provider-neutral, sovereignty-ready, and bounded in latency, cost,
storage, and information disclosure. Advisory, recommendation, and simulation
authority levels remain permitted and MUST NOT be consumed as enforcement
authority. All ADH-2026-017 scope, ownership, error, trust, non-goal, and
deferred boundaries stated below are preserved unchanged.

FEATURE-0013 inherits FEATURE-0012's API/resource grammar including resource
profiles, type metadata, common metadata, the six governance `ScopeKind`
values (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`,
`Provider`), typed references, boundaries, ownership, validation, Problem
Details, and compatibility rules. It extends this grammar only with
decision/audit domain semantics.

Phase 2 scope is contract-only: schemas, vocabularies, static registry
formats, validation, compatibility rules, and deterministic conformance
fixtures. No production persistence, queue, worker, workflow, provider
integration, cryptographic service, policy engine, AI runtime, or other
production decision/audit service is authorized.

## Glossary

Terms below are introduced or constrained by this feature. Canonical
platform terminology remains owned by `docs/glossary.md`.

| Term | Meaning in this feature |
|---|---|
| `DecisionRecord` | Immutable governed conclusion with explicit authority, rationale, inputs, composition strategy, typed result, and audit linkage. Uses FEATURE-0012 `ImmutableRecord` profile and is the sole active canonical name. |
| `DecisionProfile` | Versioned contract defining one decision family's form, authority, inputs, typed result, rationale, obligations, projections, and limits. Uses FEATURE-0012 `VersionedDefinition` profile. |
| `EvaluationResult` | Atomic, evaluator-specific normalized finding. Uses `EmbeddedValue` when captured inside a record; `TransientRequestResult` when exchanged independently. |
| `AuditEvent` | Immutable accountability record of actor/action/subject/outcome. Uses FEATURE-0012 `ImmutableRecord` profile. |
| `DecisionRationale` | Structured reasons, alternatives, obligations, warnings, and corrective suggestions embedded in a decision. |
| `DecisionContext` | Sanitized, policy-governed projection for AI and human explanation. Uses FEATURE-0012 `TransientRequestResult` profile. Owned by FEATURE-0025. |
| `DecisionRequest` | Conceptual bounded input role with context references and constraints; calling-domain-owned, not a FEATURE-0013 shared schema or type. |
| Decision form | Primary shape of a decision: ADJUDICATION, SELECTION, RANKING, CLASSIFICATION, RESOLUTION, ALLOCATION, PLAN, ASSESSMENT, or RECOMMENDATION. |
| Decision authority | Whether a record is AUTHORITATIVE, ADVISORY, RECOMMENDATION, or SIMULATION. |
| Adjudication outcome | Closed vocabulary for adjudication facets: ALLOWED, DENIED, REQUIRES_APPROVAL. |
| Finality | Every successfully persisted `DecisionRecord` is immutable and FINAL. Pending, running, and retrying belong to `Operation`. |
| Composition strategy | Versioned, registered algorithm for combining bounded evaluation results into a composite decision. |
| Sensitivity | Provider-neutral handling constraint: PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED. Separate from FEATURE-0012 `DataClassification`. |
| Structural trust | Algorithm-agile carrier fields and trust-state metadata. FEATURE-0013 validates structural presence and opaque state only; cryptographic verification is deferred under ADR-F13-002. |
| Semantic decision identity | Complete descriptor covering scope, purpose, profile, authority, canonical input snapshot, EffectivePolicyContext, composition strategy, evaluator versions, and validity epoch. Distinct from retry idempotency key. |

## User stories

- As a Sovrunn feature author building placement, governance, authorization,
  compliance, or AI-explanation capabilities, I want a stable common
  `DecisionRecord` envelope and `DecisionProfile` extension mechanism, so
  that new decision families register without changing the common contract.

- As a platform security reviewer, I want every authoritative decision to
  carry explicit authority, rationale, scope, actor, subject, and audit
  linkage, so that governance and accountability are traceable.

- As a platform operator, I want `DecisionRecord` and `AuditEvent` to be
  immutable append-only records with linked corrections, so that audit
  evidence is never silently rewritten.

- As an architecture owner, I want decision forms, authority levels, and
  adjudication outcomes as closed vocabularies, so that new decision families
  extend the model through profiles rather than ad-hoc fields.

- As a platform engineer, I want atomic evaluation results captured with
  evaluator identity, version, normalized output, and provenance, so that
  deterministic replay does not require rerunning external or AI evaluators.

- As a tenant administrator, I want decisions scoped through
  `metadata.scopeRef` using the six FEATURE-0012 governance scopes, so that
  cross-scope access is denied by default without existence disclosure.

- As an AI consumer, I want an authorized `DecisionContext` projection with
  structured reason codes, alternatives, and obligations, so that I can
  explain decisions without accessing unrestricted canonical data.

- As a sovereign deployment operator, I want all profile, schema, strategy,
  and trust material locally available for synchronous decisions, so that no
  remote registry or control plane becomes a universal hot-path dependency.

- As a compliance officer, I want decisions and audit events to carry
  retention class, legal-hold status, jurisdiction, and data classification,
  so that erasure and retention obligations are governable.

- As an enforcement point implementor, I want mandatory obligations
  validated as supportable before an ALLOWED decision is enforceable, so
  that unknown obligations force denial rather than silent bypass.

- As a feature author adopting FEATURE-0013, I want a lightweight
  downstream-adoption contract with one normative section per feature, so
  that consistency is enforced without a separate registry or approval
  workflow.

## Acceptance criteria

Requirement strength uses MUST / MUST NOT / SHOULD / MAY. Every normative
criterion carries a stable `F13-<TOPIC>-<NNN>` identifier. Implementation
class for each group is copied exactly from the architecture section 19
scorecard and the section 9–12 boundaries.

### AC-13.1 Object model and resource-profile inheritance

Implementation class: `CONTRACT_NOW` (AD-001, AD-002, AD-015, AD-029).

1. **F13-OBJ-001**: `DecisionRecord` MUST use the FEATURE-0012
   `ImmutableRecord` profile shape: type metadata, common `metadata`, and
   one immutable `record` payload. It MUST NOT use mutable `spec`/`status`.

2. **F13-OBJ-002**: `AuditEvent` MUST use the FEATURE-0012
   `ImmutableRecord` profile shape.

3. **F13-OBJ-003**: `DecisionProfile` MUST use the FEATURE-0012
   `VersionedDefinition` profile shape. Its definition content is `spec`;
   publication status follows FEATURE-0012 ownership and mutability rules.

4. **F13-OBJ-004**: `EvaluationResult` MUST use `EmbeddedValue` when
   captured inside a record and `TransientRequestResult` when exchanged
   independently. It MUST NOT have an independent durable lifecycle created
   by FEATURE-0013.

5. **F13-OBJ-005**: `DecisionContext` MUST use `TransientRequestResult`.
   It is non-authoritative and MUST NOT be persisted merely to satisfy
   resource grammar.

6. **F13-OBJ-006**: FEATURE-0013 MUST NOT define a shared
   `DecisionRequest` schema, Go type, registry entry, resource, or
   persistence model. `DecisionRequest` remains a calling-domain-owned
   conceptual input role conforming to inherited FEATURE-0012
   `TransientRequestResult` boundaries.

7. **F13-OBJ-007**: Completed synchronous handling MUST return a
   `DecisionRecord` value or canonical typed reference. Accepted
   asynchronous handling MUST return the existing FEATURE-0012 `Operation`
   value or reference conforming to `LongRunningOperation`. Rejected or
   unaccepted handling MUST return inherited Problem Details.

8. **F13-OBJ-008**: FEATURE-0013 MUST NOT define `PendingDecision`,
   `DeferredDecision`, `DecisionPending`, a pending-decision status, or any
   other non-final response envelope.

### AC-13.2 Scope authority and governance scopes

Implementation class: `CONTRACT_NOW` (AD-015).

1. **F13-SCOPE-001**: `metadata.scopeRef` MUST be the sole governance-scope
   authority for `DecisionRecord` and `AuditEvent`. No top-level `scopeRef`
   field, parallel scope enum, duplicated scope attribute, or second scope
   source is permitted.

2. **F13-SCOPE-002**: The canonical `Platform` scope MUST be represented as
   absent/nil `metadata.scopeRef`. Every non-Platform scope MUST carry a
   non-nil `metadata.scopeRef` with a canonical `ScopeKind` and UID.

3. **F13-SCOPE-003**: The allowed `ScopeKind` values are the six
   FEATURE-0012 governance scopes: `Platform`, `Organization`,
   `OrganizationUnit`, `Tenant`, `Project`, `Provider`. No additional
   scope kind (including `ServiceInstance`) is permitted.

4. **F13-SCOPE-004**: A `ServiceInstance` MUST be identified through a typed
   `subjectRef` or `subjectRefs` field, with `metadata.scopeRef` set to the
   containing Project. ServiceInstance is never a `ScopeKind`.

5. **F13-SCOPE-005**: A record declaring scope through more than one source,
   or through any source other than `metadata.scopeRef`, MUST fail
   validation with `DECISION_SCOPE_CONFLICT`.

6. **F13-SCOPE-006**: An explicit Platform input MUST normalize to the
   canonical absent serialized `metadata.scopeRef` form. A canonical nil
   `metadata.scopeRef` resolving to Platform is not a contradictory or
   duplicate scope.

7. **F13-SCOPE-007**: An out-of-vocabulary scope value MUST be rejected
   with `DECISION_SCOPE_INVALID`.

8. **F13-SCOPE-008**: A ServiceInstance subject whose `metadata.scopeRef`
   does not match its containing Project MUST be rejected with
   `DECISION_SCOPE_MISMATCH`.

### AC-13.3 Decision forms, authority, and adjudication

Implementation class: `CONTRACT_NOW` (AD-003, AD-030, AD-031).

1. **F13-FORM-001**: The closed primary-form vocabulary MUST be:
   `ADJUDICATION`, `SELECTION`, `RANKING`, `CLASSIFICATION`, `RESOLUTION`,
   `ALLOCATION`, `PLAN`, `ASSESSMENT`, `RECOMMENDATION`.

2. **F13-FORM-002**: A profile MUST declare exactly one primary form and
   MAY include registered secondary form facets.

3. **F13-AUTH-001**: The authority vocabulary MUST be: `AUTHORITATIVE`,
   `ADVISORY`, `RECOMMENDATION`, `SIMULATION`.

4. **F13-AUTH-002**: Only an authorized profile and actor MAY produce an
   authoritative record. AI output MUST NOT be authoritative.

5. **F13-AUTH-003**: Advisory, recommendation, and simulation records MUST
   be visibly distinguishable in every projection and MUST NOT be consumed
   as enforcement authority.

6. **F13-ADJ-001**: The adjudication-outcome vocabulary MUST be: `ALLOWED`,
   `DENIED`, `REQUIRES_APPROVAL`. It is mandatory for profiles with an
   adjudication facet.

7. **F13-ADJ-002**: `UNKNOWN`, runtime failure, timeout, evaluator error,
   and insufficient evidence MUST NOT be adjudication outcomes. They are
   evaluation or production states. A registered fail-safe strategy MAY
   deterministically map them to denial or `REQUIRES_APPROVAL` and MUST
   record that mapping.

### AC-13.4 Immutability, finality, and relationships

Implementation class: `CONTRACT_NOW` (AD-004, AD-034).

1. **F13-IMMUT-001**: Every successfully persisted `DecisionRecord` MUST be
   immutable and `FINAL`. Pending, running, retrying, and waiting belong
   exclusively to `Operation`.

2. **F13-IMMUT-002**: Correction, supersession, and revocation MUST create
   a new final record whose typed relationship declares `CORRECTS`,
   `SUPERSEDES`, or `REVOKES` and references its predecessor. The
   predecessor MUST NOT be mutated.

3. **F13-IMMUT-003**: A read model MAY calculate
   `effectiveState: EFFECTIVE | SUPERSEDED | REVOKED`, but calculated
   effective state is projection metadata and MUST NOT be stored by
   rewriting the canonical record.

4. **F13-IMMUT-004**: Relationship chains MUST be bounded, cycle-free, and
   deterministically resolved. A cycle MUST be rejected with
   `DECISION_RELATIONSHIP_CYCLE`. A chain exceeding the declared limit MUST
   be rejected with `DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED`.

5. **F13-IMMUT-005**: An invalid relationship kind MUST be rejected with
   `DECISION_RELATIONSHIP_KIND_INVALID`. An invalid target MUST be rejected
   with `DECISION_RELATIONSHIP_TARGET_INVALID`.

6. **F13-IMMUT-006**: `AuditEvent` MUST be append-only with the same
   correction and linkage semantics as `DecisionRecord`.

7. **F13-IMMUT-007**: Retention or lawful erasure MUST NOT rewrite the
   canonical immutable record. Structural alternatives are separately
   controlled payload custody, later key destruction, and linked
   tombstone/redaction records. FEATURE-0013 defines the structural state
   and linkage contract only; it MUST NOT implement erasure or any
   cryptographic process.

### AC-13.5 DecisionProfile extension mechanism

Implementation class: `CONTRACT_NOW` (AD-029, AD-030).

1. **F13-PROF-001**: Every `DecisionProfile` MUST declare: stable name,
   family, semantic version, owner, lifecycle status, primary form, allowed
   secondary facets, allowed authority levels, input/typed-result/rationale/
   obligation schemas, accepted evaluation types and composition strategies,
   scope/actor/subject/policy/evidence/audit semantics, determinism/failure/
   timeout/insufficient-evidence behavior, validity/expiry/correction/
   supersession/revocation rules, sensitivity/residency/retention/legal-hold/
   projection profiles, idempotency identity and bounded limits, and
   compatibility policy with fixtures/owner/reassessment triggers.

2. **F13-PROF-002**: For security/authorization/admission/sovereignty/
   compliance/privileged-access/data-movement/break-glass profiles,
   fail-closed or `REQUIRES_APPROVAL` MUST be mandatory. Fail-open MUST be
   prohibited unless a separately approved, time-bounded security exception
   with scope, owner, compensating controls, expiry, audit treatment, and
   reassessment trigger exists.

3. **F13-PROF-003**: Registration MUST require conformance review. Accepting
   arbitrary profile URIs at runtime MUST be prohibited.

4. **F13-PROF-004**: The canonical registry representation MUST be a
   declarative, versioned, provider- and language-neutral JSON bundle
   governed by JSON Schema. Go registry types and static in-memory maps are
   derivative implementation views and MUST NOT become the source of
   semantic truth.

5. **F13-PROF-005**: FEATURE-0012 platform/schema validation ceilings MUST
   be absolute outer bounds. Within those inherited ceilings, profile
   declarations own narrower maximum limits. Evaluator registrations MAY
   narrow those profile limits but MUST NOT widen, replace, or contradict
   the profile or FEATURE-0012 ceiling. The effective bound MUST be the
   most restrictive applicable value.

6. **F13-PROF-006**: Missing, unknown, incomparable, or above-ceiling
   mandatory bounds MUST fail validation with
   `DECISION_PROFILE_LIMIT_INVALID`.

7. **F13-PROF-007**: An unknown profile MUST be rejected with
   `DECISION_PROFILE_UNKNOWN`. An unsupported version MUST be rejected
   with `DECISION_PROFILE_VERSION_UNSUPPORTED`. An inactive profile MUST
   be rejected with `DECISION_PROFILE_INACTIVE`. A schema-invalid profile
   MUST be rejected with `DECISION_PROFILE_SCHEMA_INVALID`.

8. **F13-PROF-008**: Profiles and semantic registries MUST be distributable
   and verifiable offline. FEATURE-0013 validates only static bundle
   structure, identifier syntax, carrier presence, opaque
   expected-versus-observed states, and fail-closed trust-state transitions.

### AC-13.6 Evaluation and composition

Implementation class: `CONTRACT_NOW` (AD-005, AD-006, AD-037).

1. **F13-EVAL-001**: An atomic evaluation MUST have one registered evaluator
   type, one bounded input snapshot reference or opaque integrity carrier,
   one result, one timing envelope, and explicit success/non-match/
   indeterminate/timeout/error semantics.

2. **F13-EVAL-002**: Nondeterministic or AI evaluator results MUST capture
   evaluator identity and version, governed input snapshot reference and
   optional opaque integrity carrier, returned output, relevant execution
   configuration, evaluation time, trust boundary, safety/policy filters,
   and structural trust state. Composition MUST replay captured output
   without rerunning the evaluator.

3. **F13-EVAL-003**: `EvaluationResult` MUST NOT establish an independent
   governance scope. Its accepted scope MUST be derived from the containing
   `DecisionRecord`'s canonical `metadata.scopeRef` and the profile's
   declared cross-scope rules.

4. **F13-EVAL-004**: Evaluator-supplied scope metadata is evidence only. It
   MUST equal or be an explicitly permitted narrowing of the record scope.
   It MUST NOT become a second scope authority. A mismatch, unknown scope,
   unauthorized widening, or inaccessible reference MUST fail closed with
   `DECISION_EVALUATION_SCOPE_MISMATCH` without disclosing whether the
   referenced target exists.

5. **F13-EVAL-005**: An unsupported evaluator type MUST be rejected with
   `DECISION_EVALUATION_TYPE_UNSUPPORTED`. An invalid result MUST be
   rejected with `DECISION_EVALUATION_RESULT_INVALID`. A limit exceeding
   the profile MUST be rejected with
   `DECISION_EVALUATION_LIMIT_EXCEEDED`.

6. **F13-COMP-001**: A composite decision MUST combine bounded evaluation
   results using a versioned, registered strategy declaring: stable
   identifier and version, accepted input types, ordering and precedence,
   short-circuit behavior, missing/timeout/conflict/error behavior within the governing profile boundary, deterministic output mapping, resource
   budgets and maximum fan-out, and explanation/evidence rules. Any requested fail-open behavior MUST be governed only by the architecture-defined SecurityExceptionRef contract; strategy metadata is never an independent fail-open authority or approval source.

7. **F13-COMP-002**: An unsupported strategy MUST be rejected with
   `DECISION_COMPOSITION_STRATEGY_UNSUPPORTED`. An invalid graph MUST be
   rejected with `DECISION_COMPOSITION_GRAPH_INVALID`. A cycle MUST be
   rejected with `DECISION_COMPOSITION_CYCLE`. A limit exceeded MUST be
   rejected with `DECISION_COMPOSITION_LIMIT_EXCEEDED`. An input conflict
   MUST be rejected with `DECISION_COMPOSITION_INPUT_CONFLICT`.

8. **F13-COMP-003**: The bounded decision-requirements graph MUST have:
   stable node/edge identifiers, typed nodes, single root decision, maximum
   nodes/depth/width/fan-out/payload/wall-clock budget, maximum remote
   calls/unique authorities/bytes/concurrent evaluations/per-evaluator
   timeout/retries/cross-jurisdiction calls/cost units, cycle/duplicate/
   missing-dependency/incompatible-version rejection, and recorded graph
   definition/version with optional opaque integrity carrier.

9. **F13-COMP-004**: The graph MUST NOT execute provisioning, human tasks,
   compensations, arbitrary code, or unbounded loops.

10. **F13-COMP-005**: Governance hierarchy MUST be resolved before final
    decision composition into one `EffectivePolicyContext`. The decision
    layer MUST NOT independently traverse governance trees.

### AC-13.7 Sensitivity, security, and projection

Implementation class: `CONTRACT_NOW` (AD-019, AD-020, AD-021, AD-024).

1. **F13-SEC-001**: The sensitivity vocabulary MUST be:
   `PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED`. No fewer and no more
   than these four core values are permitted.

2. **F13-SEC-002**: FEATURE-0013 sensitivity MUST NOT replace or redefine
   FEATURE-0012 `DataClassification`. The two dimensions coexist. The
   canonical minimum mapping from FEATURE-0012 classification to
   FEATURE-0013 sensitivity floor MUST be:
   - `Public` → `PUBLIC`
   - `Customer-visible` → `CONFIDENTIAL`
   - `Tenant-confidential` → `CONFIDENTIAL`
   - `Operator-confidential` → `CONFIDENTIAL`
   - `Internal` → `INTERNAL`
   - `Sensitive` → `RESTRICTED`
   - `Secret-reference-only` → `RESTRICTED`

3. **F13-SEC-003**: A profile MAY raise the sensitivity floor but MUST NOT
   lower it. The original FEATURE-0012 classification MUST remain present.
   An unmapped, unknown, contradictory, or below-floor value MUST fail
   closed.

4. **F13-SEC-004**: A profile permitting captured evaluator output MUST
   declare: a sensitivity ceiling, allowed and prohibited content-category
   rules, a maximum field count, a maximum byte size, a maximum nesting
   depth, and the allowed value types. Missing or unknown mandatory
   profile values MUST invalidate the profile.

5. **F13-SEC-005**: Captured values classified above the declared ceiling
   MUST fail validation.

6. **F13-SEC-006**: FEATURE-0013 owns only structural schema validation,
   typed classification, declared bounds, projection restrictions, and
   deterministic conformance. It MUST NOT claim comprehensive semantic
   detection of secrets, credentials, personal data, malware, or policy
   violations. Producers MUST redact prohibited content before submission.

7. **F13-SEC-007**: Different authorized projections MUST serve customer,
   operator, auditor, regulator, provider, and AI consumers. Redaction
   MUST be field- and reason-aware and MUST NOT change the final outcome
   or falsely imply that omitted evidence did not exist.

8. **F13-SEC-008**: AI consumers MUST receive only an authorized
   `DecisionContext` projection. They MUST NOT mutate records, invent
   evidence, suppress audit, widen disclosure, bypass approval, or execute
   suggestions. Any AI-generated recommendation MUST be labelled with
   model/provider/version, prompt-template version, data boundary,
   confidence where valid, and human-review status.

9. **F13-SEC-009**: Core decisions MUST remain available when AI is disabled
   or disconnected.

### AC-13.8 Audit consistency and obligation

Implementation class for schema/linkage/validation: `CONTRACT_NOW` (AD-014,
AD-026, AD-035). Implementation class for durable-acceptance semantics:
`INVARIANT_FOR_LATER`.

1. **F13-AUDIT-001**: Every meaningful decision MUST link to at least one
   `AuditEvent`. A single audit event MAY reference a decision, operation,
   actor, subject, resource, request, trace, and parent event without
   copying unrestricted payloads.

2. **F13-AUDIT-002**: Audit outcomes MUST be: `Succeeded`, `Denied`,
   `Failed`. These are separate from decision outcomes. Successfully
   producing a `DENIED` adjudication is an audit outcome of `Succeeded`.

3. **F13-AUDIT-003**: `AuditEvent` MUST accept all six canonical
   FEATURE-0012 governance `ScopeKind` values. It MUST retain the existing
   Organization-scoped alpha fixture as regression evidence. A
   ServiceInstance subject fixture MUST use its Project `metadata.scopeRef`.

4. **F13-AUDIT-004**: Telemetry MUST NOT be treated as audit authority.
   Transport events are not substitutes for `AuditEvent`.

5. **F13-AUDIT-005**: The contract MUST carry the atomic durable-acceptance
   invariant: a final decision MUST NOT be reported as durably recorded
   until its immutable `DecisionRecord` and its required audit obligation
   are atomically accepted by the same local authoritative persistence
   boundary. This is an `INVARIANT_FOR_LATER` that FEATURE-0013 carries,
   documents, and conformance-tests without implementing a production
   persistence mechanism.

6. **F13-AUDIT-006**: The atomic boundary MUST preallocate the immutable
   `AuditEvent` identity and store that identity in both the
   `DecisionRecord` reference and audit obligation. Asynchronous
   materialization MUST use the same identity and MUST NOT require mutation
   of the `DecisionRecord`. This is an `INVARIANT_FOR_LATER`.

7. **F13-AUDIT-007**: Durable audit-obligation acceptance is sufficient for
   the synchronous response; full materialization, indexing, export,
   notification, SIEM delivery, cross-site replication, or regulator
   delivery MAY proceed asynchronously. This is an `INVARIANT_FOR_LATER`.

### AC-13.9 Statelessness, scalability, and cost controls

Implementation class for carrier fields and schema: `CONTRACT_NOW` (AD-012,
AD-013, AD-036, AD-040). Implementation class for production
scaling/delivery: `INVARIANT_FOR_LATER`. Implementation class for production
store/queue/topology: `DEFERRED`.

1. **F13-SCALE-001**: The schema MUST carry a caller-supplied idempotency
   key scoped to caller and request purpose (retry identity) as a distinct
   field from the complete semantic decision identity descriptor.

2. **F13-SCALE-002**: The semantic decision identity descriptor MUST cover:
   scope, purpose, profile name and version, authority, canonical input
   snapshot, `EffectivePolicyContext`, composition strategy/version,
   evaluator contract versions, and relevant evaluation-time or validity
   epoch.

3. **F13-SCALE-003**: FEATURE-0013 MUST NOT define new universal numeric
   decision-domain defaults or system-wide graph ceilings. Existing
   FEATURE-0012 ceilings remain inherited outer bounds. Every conforming
   profile MUST declare finite values for all domain limits it uses within
   those ceilings.

4. **F13-SCALE-004**: Conformance fixtures MUST use explicit fixture-local
   profile values to test boundary behavior and MUST NOT invent FEATURE-0013
   defaults, fallbacks, or hidden maximums.

5. **F13-SCALE-005**: The `CONTRACT_NOW` reference kernel MUST be a pure,
   deterministic function of versioned inputs, registered strategy, and
   explicit fixture-local bounded configuration. It MUST hold no
   authoritative session state and MUST perform no production persistence,
   delivery, locking, caching, or scaling behavior.

6. **F13-SCALE-006**: The following are `INVARIANT_FOR_LATER` carried and
   documented by FEATURE-0013: stateless interchangeable workers, external
   durable state behind ports, compare-and-create semantics, effective-once
   behavior, bounded delivery/retry/backpressure/cache behavior, runtime
   resource budgets, no sticky sessions, and scale-to-zero permission.

7. **F13-SCALE-007**: The following are `DEFERRED`: every concrete store,
   transaction, lock, lease, queue, dead-letter processor, cache, worker
   topology, scaling policy, and production numeric default.

### AC-13.10 Obligations

Implementation class: `CONTRACT_NOW` (AD-032, AD-043).

1. **F13-OBLIG-001**: An enforcement point MUST enforce all mandatory
   obligations or deny the action.

2. **F13-OBLIG-002**: An unknown, invalid, expired, or unsupported
   mandatory obligation MUST make an `ALLOWED` decision unenforceable
   and therefore MUST result in denial at that enforcement point.

3. **F13-OBLIG-003**: Unknown obligations MUST be rejected with
   `DECISION_OBLIGATION_UNKNOWN`. Invalid obligations MUST be rejected
   with `DECISION_OBLIGATION_INVALID`. Unsupported mandatory obligations
   MUST be rejected with `DECISION_OBLIGATION_UNSUPPORTED_MANDATORY`.

### AC-13.11 Structural trust and integrity carriers

Implementation class: `CONTRACT_NOW` for carrier fields and structural
trust-state semantics (AD-022). Implementation class for cryptographic
execution: `DEFERRED` (ADR-F13-002).

1. **F13-TRUST-001**: The contract MUST carry algorithm-agile fields:
   canonicalization-method identifier, digest-algorithm identifier,
   covered-field descriptor, signature-algorithm identifier, and structural
   trust-state fields.

2. **F13-TRUST-002**: FEATURE-0013 MUST validate only structural presence,
   identifier syntax, carrier presence, opaque expected-versus-observed
   states, and fail-closed trust-state transitions.

3. **F13-TRUST-003**: When a profile requires verified trust, an unknown,
   absent, expired, revoked, mismatched, or unverified trust state MUST
   fail closed with the appropriate `DECISION_TRUST_*` code:
   `DECISION_TRUST_REQUIRED`, `DECISION_TRUST_UNKNOWN`,
   `DECISION_TRUST_EXPIRED`, `DECISION_TRUST_REVOKED`, or
   `DECISION_TRUST_MISMATCH`.

4. **F13-TRUST-004**: FEATURE-0013 MUST NOT compute a digest over content,
   canonicalize content, sign, verify, timestamp, notarize, enforce WORM
   retention, or claim that an opaque carrier proves content existed or
   was considered. All algorithm selections, covered-field sets, signing,
   verification, product, and service selections are `DEFERRED` under
   ADR-F13-002.

5. **F13-TRUST-005**: Structural fixtures MUST use opaque deterministic
   test assertions to validate state and failure semantics and MUST NOT
   claim cryptographic validity or select an algorithm.

### AC-13.12 Sovereign deployment contracts

Implementation class: `CONTRACT_NOW` for static profile/carrier contracts
(AD-016, AD-017, AD-018, AD-023, AD-042). Implementation class for
synchronization runtime: `INVARIANT_FOR_LATER`.

1. **F13-SOV-001**: Sovereignty MUST be treated as seven enforceable
   dimensions: data, operational, technical, legal, cryptographic, software
   supply chain, and AI/model.

2. **F13-SOV-002**: The same canonical contracts MUST operate in: connected
   public/sovereign cloud, private cloud/GCC, on-premises/regulated
   enterprise, intermittently connected edge, and fully disconnected/
   air-gapped deployment.

3. **F13-SOV-003**: No canonical object MUST contain a cloud provider,
   infrastructure platform, container orchestrator, AI service, product
   identifier, region code, or provider-native resource type as a required
   semantic.

4. **F13-SOV-004**: Export and import MUST use collision-resistant record
   identifiers, origin trust-domain identity, algorithm-agile signature
   carriers, replay markers, duplicate detection, per-origin ordering, and
   explicit opaque trust-root identity.

5. **F13-SOV-005**: The manifest MUST NOT be cryptographically signed or
   verified by FEATURE-0013. Unknown profiles, origins, or opaque trust
   states MUST fail closed or be represented as quarantined.

6. **F13-SOV-006**: Immutable records MUST NOT be merged using
   last-writer-wins. Cross-origin conflicts MUST produce explicit evidence
   and a governed reconciliation decision.

7. **F13-SOV-007**: The contracts MUST carry or reference (without embedding
   secrets): governing jurisdiction, residency constraints, data
   classification, handling caveats, authority/purpose/policy version,
   tenant/department/mission boundary, retention class, legal hold,
   disposition status, actor/workload/delegation/assurance identity, key
   domain and cryptographic proof references, approved disclosure/export
   policy, opaque integrity/trusted-time provenance carriers, and
   emergency/break-glass use with mandatory reason.

### AC-13.13 Public error contract

Implementation class: `CONTRACT_NOW` (AD-045).

1. **F13-ERR-001**: FEATURE-0013 structural or semantic validation failures
   MUST inherit the FEATURE-0012 RFC 9457 Problem Details contract with:
   - `type`: `urn:sovrunn:problem:validation-failed`
   - `status`: `422`
   - `code`: `VALIDATION_FAILED`
   - `violations[].field`: RFC 6901 JSON Pointer to the rejected field or
     nearest safe parent
   - `violations[].code`: exactly one code from the closed registry below
   - `violations[].message`: human-readable, non-authoritative, redactable

2. **F13-ERR-002**: The closed public `violations[].code` registry MUST be:

   **Decision profile family:**
   `DECISION_PROFILE_UNKNOWN`,
   `DECISION_PROFILE_VERSION_UNSUPPORTED`,
   `DECISION_PROFILE_INACTIVE`,
   `DECISION_PROFILE_SCHEMA_INVALID`,
   `DECISION_PROFILE_LIMIT_INVALID`

   **Evaluation family:**
   `DECISION_EVALUATION_TYPE_UNSUPPORTED`,
   `DECISION_EVALUATION_RESULT_INVALID`,
   `DECISION_EVALUATION_SCOPE_MISMATCH`,
   `DECISION_EVALUATION_LIMIT_EXCEEDED`

   **Composition family:**
   `DECISION_COMPOSITION_STRATEGY_UNSUPPORTED`,
   `DECISION_COMPOSITION_GRAPH_INVALID`,
   `DECISION_COMPOSITION_CYCLE`,
   `DECISION_COMPOSITION_LIMIT_EXCEEDED`,
   `DECISION_COMPOSITION_INPUT_CONFLICT`

   **Obligation family:**
   `DECISION_OBLIGATION_UNKNOWN`,
   `DECISION_OBLIGATION_INVALID`,
   `DECISION_OBLIGATION_UNSUPPORTED_MANDATORY`

   **Trust family:**
   `DECISION_TRUST_REQUIRED`,
   `DECISION_TRUST_UNKNOWN`,
   `DECISION_TRUST_EXPIRED`,
   `DECISION_TRUST_REVOKED`,
   `DECISION_TRUST_MISMATCH`

   **Scope family:**
   `DECISION_SCOPE_REQUIRED`,
   `DECISION_SCOPE_INVALID`,
   `DECISION_SCOPE_CONFLICT`,
   `DECISION_SCOPE_MISMATCH`,
   `DECISION_SCOPE_WIDENING`

   **Relationship family:**
   `DECISION_RELATIONSHIP_KIND_INVALID`,
   `DECISION_RELATIONSHIP_TARGET_INVALID`,
   `DECISION_RELATIONSHIP_CYCLE`,
   `DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED`,
   `DECISION_RELATIONSHIP_CONFLICT`

3. **F13-ERR-003**: These prefixes are public namespaces only for
   `violations[].code`. They MUST NOT be used as `Problem.type` URI
   families, top-level `Problem.code` values, HTTP-status selectors, or
   internal-only categories.

4. **F13-ERR-004**: A new code, rename, removal, changed meaning, different
   prefix, or remapping to a different top-level problem MUST require an
   approved architecture and compatibility change.

5. **F13-ERR-005**: Requirements MAY enumerate testable cases and map each
   case to an existing code and JSON Pointer. Requirements MUST NOT create,
   rename, alias, or assign new semantics to public codes.

### AC-13.14 Compatibility and evolution

Implementation class: `CONTRACT_NOW` (AD-015, AD-044).

1. **F13-COMPAT-001**: Every canonical schema MUST have explicit version
   and compatibility rules.

2. **F13-COMPAT-002**: Additive optional fields MUST require semantic
   defaults and projection safety.

3. **F13-COMPAT-003**: Decision form, authority, adjudication, typed-result,
   reason-code, obligation, strategy, and digest semantics MUST NOT change
   silently.

4. **F13-COMPAT-004**: Readers MUST reject unsupported major versions and
   MUST preserve unknown safe extension fields only under an approved
   extension mechanism.

5. **F13-COMPAT-005**: Registries for decision profile, form, authority,
   evaluator type, strategy, reason code, obligation, event type, evidence
   type, and projection MUST be versioned and governable offline.

6. **F13-COMPAT-006**: Every inherited FEATURE-0011 or FEATURE-0012
   vocabulary, shared type, error contract, schema convention, or validation
   rule MUST have one dependency-reconciliation record before design
   approval. The record MUST identify exact versions and files, compare
   enums/requiredness/cardinality/ownership/serialization/validation/
   security semantics/bindings/fixtures/migration impact, and classify
   every difference. This record is completed in this requirements document
   in the "FEATURE-0011 and FEATURE-0012 dependency reconciliation
   (ADH-2026-017 locked)" section, bound to the exact ADH-2026-017
   `dependency_content_sha256` paths and hashes and covering all section 15.3
   dimensions with every delta classified. Design MUST NOT create, reopen,
   defer, or reformat it.

7. **F13-COMPAT-007**: The six FEATURE-0012 ScopeKind values MUST remain
   unchanged. `AuditEvent` MUST accept all six governance scopes without
   changing the shared ScopeKind vocabulary.

8. **F13-COMPAT-008**: The existing Organization-scoped FEATURE-0012
   AuditEvent fixture MUST remain valid as regression evidence.

9. **F13-COMPAT-009**: FEATURE-0012 metadata, typed-reference,
   validation-layer, Problem Details, extension, and
   ServiceInstance-as-Project-scoped-resource semantics MUST remain
   unchanged.

### AC-13.15 Downstream adoption contract

Implementation class: `CONTRACT_NOW` (AD-044).

1. **F13-ADOPT-001**: Later features MUST include one short
   `FEATURE-0013 Adoption` section when they directly depend on FEATURE-0013
   or produce/consume/specialize/project/audit/enforce any governed concept.

2. **F13-ADOPT-002**: The section MUST classify the feature as `APPLICABLE`
   or `NOT_APPLICABLE`. A `NOT_APPLICABLE` declaration MUST include a
   concise rationale.

3. **F13-ADOPT-003**: An applicable feature MUST declare roles from:
   `EVALUATION_PRODUCER`, `DECISION_PRODUCER`, `DECISION_CONSUMER`,
   `AUDIT_PRODUCER`, `PROJECTION_CONSUMER`, `OPERATION_OWNER`,
   `COMPOSITE_ADOPTER`.

4. **F13-ADOPT-004**: The feature-gate check MUST verify section existence,
   applicability, roles, named/versioned profiles, reused contracts,
   applicable risks, exceptions, conformance evidence, and that the common
   envelope and closed vocabularies were not redefined.

5. **F13-ADOPT-005**: FEATURE-0013 MUST NOT require separate per-feature
   YAML manifests, a central adoption registry, a generated adoption index,
   a separate adoption approval, or a runtime adoption service. Those are
   prohibited until a section 28.8 reassessment trigger is met.

### AC-13.16 Conformance model

Implementation class: `CONTRACT_NOW` (AD-028).

1. **F13-CONF-001**: Executable positive and negative conformance fixtures
   MUST cover all F13-CF-01 through F13-CF-28 scenario families, and MUST
   cover every mandatory lowercase suffix (in architecture section 17
   textual order) for the compound families as an individual coverage
   obligation. Per the architecture, coverage of a compound parent requires
   coverage of every suffix. The mandatory suffixes are enumerated
   individually in the "Compound conformance suffixes" traceability
   subsection and are: F13-CF-04.a–.i (ranking, classification, resolution,
   allocation, plan, assessment, recommendation, advisory, simulation);
   F13-CF-08.a–.d (timeout, missing evidence, evaluator failure, policy
   conflict); F13-CF-11.a–.e (supersession, correction, revocation, cycle,
   chain-limit); F13-CF-12.a–.e (unauthorized profile, authority, projection,
   obligation bypass, cross-scope); F13-CF-16.a–.f (graph size, depth, width,
   fan-out, timeout, payload budget); F13-CF-20.a–.f (policy, profile,
   strategy, authority, evaluator-version, validity); F13-CF-26.a–.h
   (cancellation, remote-call, byte, concurrency, retry, jurisdiction, time,
   cost budgets); and F13-CF-28.a–.f (replay, duplicate, ordering, unknown
   trust root, revoked origin, conflict reconciliation). These suffixes MUST
   NOT be merged, omitted, renumbered, or invented.

2. **F13-CONF-002**: Additional conformance cases F13-SCOPE-01 through
   F13-SCOPE-11, F13-EVAL-01 through F13-EVAL-08, F13-SEC-01 through
   F13-SEC-08, F13-TRUST-01 through F13-TRUST-04, and F13-COMPAT-01
   through F13-COMPAT-09 MUST each have explicit coverage-matrix entries.

3. **F13-CONF-003**: All fixtures MUST be deterministic, pure or bounded
   in-memory contract tests. They MUST NOT authorize persistence,
   networking, queues, workers, cryptographic execution, AI runtime, or
   external integration.

4. **F13-CONF-004**: Coverage is counted by scenario ID, not by fixture-file
   count. A physical fixture MAY cover several IDs, but every ID MUST have
   an explicit coverage-matrix row.

5. **F13-CONF-005**: The `CONTRACT_NOW` reference kernel MUST be sufficient
   to prove contract semantics through bounded in-memory or pure-function
   reference composition. Complex conformance MUST remain a
   fixture/reference proof, not a runtime engine.

### AC-13.17 Authorization specialization

Implementation class: `CONTRACT_NOW` for structural profile only (AD-032,
AD-033, AD-039, AD-043).

1. **F13-AUTHZ-001**: `AuthorizationDecision` MUST be a profile of
   `DecisionRecord`, not a separate core object. It MUST require an
   `ADJUDICATION` facet.

2. **F13-AUTHZ-002**: Sensitive, denied, privileged, cross-scope,
   break-glass, export, and policy-change decisions MUST require durable
   records. High-volume low-risk checks MAY use a governed aggregation
   profile.

3. **F13-AUTHZ-003**: An ephemeral low-risk check not persisted as a
   `DecisionRecord` MUST NOT be a durable Sovrunn decision. Its audit
   evidence remains governed by the authorization profile.

4. **F13-AUTHZ-004**: FEATURE-0013 MUST NOT require a remote central
   authorization call for every request. Local policy decision points,
   version-pinned context, short-lived decision artifacts (with profile-
   defined audience/scope/expiry/revocation/freshness/fail-closed behavior),
   and safe caches are later permitted patterns. This is an
   `INVARIANT_FOR_LATER`.

5. **F13-AUTHZ-005**: Every human, workload, evaluator, and service hop in a
   decision or audit flow MUST be authenticated and authorized, bound to
   scope, purpose, action, and resource. This translates architecture section
   12.1's every-hop authentication/authorization duty and is an
   `INVARIANT_FOR_LATER`: FEATURE-0013 carries and documents it for a later
   owning feature and MUST NOT implement any runtime authentication,
   authorization, session, or enforcement mechanism. It generates no
   FEATURE-0013 implementation task beyond documentation and traceability.

6. **F13-AUTHZ-006**: Break-glass or emergency access MUST require follow-up
   review after the fact. This translates architecture section 12.1's
   break-glass follow-up-review duty and is an `INVARIANT_FOR_LATER`:
   FEATURE-0013 carries and documents it for a later owning feature and MUST
   NOT implement the review workflow, notification, or enforcement. It
   generates no FEATURE-0013 implementation task beyond documentation and
   traceability. The mandatory break-glass reason carrier itself remains
   `CONTRACT_NOW` and is owned by F13-SOV-007.

### AC-13.18 Anti-overengineering guardrails

Implementation class: `CONTRACT_NOW` (AD-028).

1. **F13-GUARD-001**: Every generated requirement, design component, and
   task MUST declare exactly one implementation class: `CONTRACT_NOW`,
   `INVARIANT_FOR_LATER`, or `DEFERRED`. An item without a class is invalid.

2. **F13-GUARD-002**: `INVARIANT_FOR_LATER` and `DEFERRED` items MUST NOT
   generate implementation tasks except documentation, traceability, and
   future conformance placeholders that do not execute production behavior.

3. **F13-GUARD-003**: FEATURE-0013 MUST NOT create or select any of the
   prohibited artifacts listed in architecture section 27.3 (production
   services, databases, queues, workflows, network clients, runtimes,
   policy engines, cryptographic services, provider-specific dependencies).

4. **F13-GUARD-004**: Before adding any abstraction, design MUST answer the
   five simplicity/necessity tests from architecture section 27.4. If
   question 1 has no exact requirement ID or question 5 is yes, the
   abstraction MUST NOT be added.

5. **F13-GUARD-005**: The simple atomic decision path MUST remain
   representable without a graph, orchestrator, registry network call,
   operation, AI, evidence store, or runtime service.

6. **F13-GUARD-006**: Unknown product choices, owners, limits, retention
   periods, SLOs, or legal values MUST be marked
   `ARCHITECTURE_DECISION_REQUIRED`; they MUST NOT be guessed.

## Non-goals

The following are explicitly excluded from FEATURE-0013 scope:

1. Production decision, authorization, policy, placement, audit, evidence,
   identity, profile-registry, synchronization, cryptographic, or AI
   service.
2. Database schema, migration, repository implementation, transaction/outbox,
   event store, WORM store, cache, distributed lock, or lease.
3. Queue, broker, workflow/durable-execution runtime, scheduler, controller,
   worker pool, background reconciler, retry daemon, or dead-letter processor.
4. Network registry client, remote policy resolver, external evidence
   fetcher, HSM/notary client, provider adapter, model client, vector store,
   or agent.
5. General-purpose rules engine, policy language, workflow DSL, DAG execution
   engine, expression interpreter, dynamic plugin system, or code generator.
6. Production authorization tokens, cache invalidation, revocation
   distribution, multi-site consensus, disconnected synchronization service,
   or conflict-resolution service.
7. Provider-, runtime-, database-, broker-, identity-, cryptographic-, or
   model-specific dependencies in core contract packages.
8. A shared `DecisionRequest` schema, Go type, registry entry, resource,
   or persistence model.
9. Any pending-decision response envelope (`PendingDecision`,
   `DeferredDecision`, `DecisionPending`, or similar).
10. Adapter interfaces (`EvaluatorAdapter`, `PolicyEngineAdapter`,
    `IdentityProviderAdapter`, `SecretProviderAdapter`,
    `ObservabilityAdapter`, etc.) — these are owned by FEATURE-0016.
11. Cryptographic algorithm selection, digest computation, signing,
    verification, trusted timestamping, notarization, WORM enforcement, HSM
    integration, or canonicalization execution (all deferred under
    ADR-F13-002).
12. Production numeric defaults, system-wide graph ceilings, retention
    periods, or SLO targets not yet approved by architecture.
13. India-specific government security-classification vocabulary and exact
    retention schedules.
14. Concrete provider/runtime/database/broker/identity/AI adapters.
15. The `Operation` contract itself — owned by FEATURE-0012/ADH-2026-013;
    FEATURE-0013 references but does not redefine it.

## Edge cases

The following are non-normative worked examples that illustrate
already-normative requirements. Each example is owned by the acceptance-
criteria requirement cited below, which holds the normative force and the
architecture-owned implementation class; the conformance ID (or, where no
dedicated conformance ID exists, the owning validator requirement) proves it.
These
examples add no new obligation and generate no implementation task beyond the
cited requirement. All cited owning requirements are `CONTRACT_NOW`.

| Example | Owning requirement (class) | Conformance |
|---|---|---|
| 1 | F13-SCOPE-001 (AC-13.2, `CONTRACT_NOW`) | F13-SCOPE-07 |
| 2 | F13-SCOPE-005 (AC-13.2, `CONTRACT_NOW`) | F13-SCOPE-09 |
| 3 | F13-EVAL-004 (AC-13.6, `CONTRACT_NOW`) | F13-EVAL-06 |
| 4 | F13-COMP-002 (AC-13.6, `CONTRACT_NOW`) | F13-COMP-002 (composition-cycle validator; no dedicated CF ID) |
| 5 | F13-COMP-002 (AC-13.6, `CONTRACT_NOW`) | F13-CF-16.a, F13-CF-16.b, F13-CF-16.c, F13-CF-16.d, F13-CF-16.f |
| 6 | F13-SCOPE-003 (AC-13.2, `CONTRACT_NOW`) | F13-SCOPE-10 |
| 7 | F13-IMMUT-004 (AC-13.4, `CONTRACT_NOW`) | F13-CF-11.d |
| 8 | F13-IMMUT-004 (AC-13.4, `CONTRACT_NOW`) | F13-CF-11.e |
| 9 | F13-SEC-005 (AC-13.7, `CONTRACT_NOW`) | F13-SEC-01 |
| 10 | F13-PROF-006 (AC-13.5, `CONTRACT_NOW`) | F13-EVAL-04 |
| 11 | F13-EVAL-005 (AC-13.6, `CONTRACT_NOW`) | F13-EVAL-03 |
| 12 | F13-SCOPE-006 (AC-13.2, `CONTRACT_NOW`) | F13-SCOPE-08 |
| 13 | F13-TRUST-003 (AC-13.11, `CONTRACT_NOW`) | F13-TRUST-03 |
| 14 | F13-OBLIG-002 (AC-13.10, `CONTRACT_NOW`) | F13-CF-24 |
| 15 | F13-EVAL-002 (AC-13.6, `CONTRACT_NOW`) | F13-CF-21 |
| 16 | F13-PROF-007 (AC-13.5, `CONTRACT_NOW`) | F13-CF-23 |
| 17 | F13-SOV-005 (AC-13.12, `CONTRACT_NOW`) | F13-CF-28.d |
| 18 | F13-SEC-008 (AC-13.7, `CONTRACT_NOW`) | F13-CF-15 |

1. A `DecisionRecord` with `metadata.scopeRef` absent and the profile does
   not permit Platform scope — MUST reject with `DECISION_SCOPE_REQUIRED`.

2. A `DecisionRecord` with both `metadata.scopeRef` and a top-level
   `scopeRef` — MUST reject with `DECISION_SCOPE_CONFLICT`.

3. An evaluator reporting a scope wider than the containing record's scope
   — MUST reject with `DECISION_EVALUATION_SCOPE_MISMATCH` without
   existence disclosure.

4. A composition graph containing a cycle — MUST reject with
   `DECISION_COMPOSITION_CYCLE`.

5. A composition graph exceeding declared maximum nodes, depth, width,
   fan-out, or payload — MUST reject with
   `DECISION_COMPOSITION_LIMIT_EXCEEDED`.

6. An `AuditEvent` with `ScopeKind: ServiceInstance` — MUST reject with
   `DECISION_SCOPE_INVALID`. ServiceInstance is a subject, not a scope.

7. A relationship chain that forms a cycle through supersession,
   correction, or revocation — MUST reject with
   `DECISION_RELATIONSHIP_CYCLE`.

8. A relationship chain exceeding declared maximum length — MUST reject
   with `DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED`.

9. A captured evaluation value classified above the profile's declared
   sensitivity ceiling — MUST fail validation.

10. A profile declaring a limit wider than the FEATURE-0012 ceiling — MUST
    reject with `DECISION_PROFILE_LIMIT_INVALID`.

11. An evaluator registration declaring a limit wider than its owning
    profile — MUST reject with `DECISION_EVALUATION_LIMIT_EXCEEDED`.

12. An explicit Platform scope in a request normalizes to absent
    `metadata.scopeRef`; the normalized form MUST NOT trigger
    `DECISION_SCOPE_CONFLICT`.

13. A trust carrier with opaque expected state not matching observed state
    — MUST fail closed with `DECISION_TRUST_MISMATCH`.

14. A mandatory obligation type not supported by the enforcement point
    — MUST reject with `DECISION_OBLIGATION_UNSUPPORTED_MANDATORY` and
    MUST make the `ALLOWED` decision unenforceable.

15. A nondeterministic evaluator produces different output when re-queried
    — composition MUST use the captured normalized output from the original
    evaluation and MUST NOT rerun the evaluator.

16. An unknown profile version in a static bundle — MUST reject with
    `DECISION_PROFILE_VERSION_UNSUPPORTED` and MUST NOT fall back to an
    older version.

17. A cross-origin import with an unknown trust-root identity — MUST fail
    closed or quarantine; MUST NOT accept via last-writer-wins.

18. An AI consumer requesting unrestricted canonical data — MUST receive
    only the authorized `DecisionContext` projection; request MUST be denied
    without disclosing the canonical record.

## Security and privacy requirements

This section is a non-normative restatement of security and privacy
obligations already owned by the acceptance criteria above. Each obligation
cites its owning requirement and carries exactly one architecture-owned
implementation class. Structural carrier and validation obligations are
`CONTRACT_NOW` and are implemented by FEATURE-0013. Runtime authentication,
authorization, and follow-up-review invariants are `INVARIANT_FOR_LATER`:
they are carried and documented but produce no FEATURE-0013 implementation
work. This section creates no additional obligation.

| # | Obligation | Kind | Owning requirement | Class |
|---|---|---|---|---|
| 1 | Canonical records (`DecisionRecord`, `AuditEvent`) do not indiscriminately duplicate personal data, secrets, credentials, tokens, raw policy, model prompts, or unrestricted content; they prefer references, digests, classifications, and purpose-bound projections | Structural carrier | F13-SEC-006; AD-024 | `CONTRACT_NOW` |
| 2 | Every human, workload, evaluator, and service hop is authenticated and authorized, bound to scope, purpose, action, and resource | Runtime authN/authZ | F13-AUTHZ-005; AD-039 | `INVARIANT_FOR_LATER` |
| 3 | Delegation and impersonation are recorded explicitly through the decision record's identity carrier fields | Structural carrier | F13-SOV-007; AD-027 | `CONTRACT_NOW` |
| 4 | Evaluator output, reason text, external evidence, and AI content are treated as untrusted, subject to schema, size, depth, character, URI, and reference-kind limits | Structural validation | F13-SEC-004; F13-SEC-006 | `CONTRACT_NOW` |
| 5 | Cross-scope references neither disclose nor influence another tenant or governance domain; inaccessible references fail without existence disclosure | Structural fail-closed | F13-EVAL-004 | `CONTRACT_NOW` |
| 6 | Break-glass or emergency access carries a mandatory reason | Structural carrier | F13-SOV-007 | `CONTRACT_NOW` |
| 7 | Break-glass or emergency access requires follow-up review | Runtime review duty | F13-AUTHZ-006 | `INVARIANT_FOR_LATER` |
| 8 | Retention class, legal-hold status, and disposition are carried without implementing the erasure or retention mechanism | Structural carrier | F13-IMMUT-007; F13-SOV-007 | `CONTRACT_NOW` |
| 9 | No secret-pattern, credential-pattern, PII, malware, DLP, or content-scanning engine is implemented; security validation is structural and typed conformance only | Structural boundary | F13-SEC-006 | `CONTRACT_NOW` |
| 10 | AI-generated recommendations are labelled with model/provider/version, prompt-template version, data boundary, confidence, and human-review status; AI is never authoritative | Structural labelling | F13-SEC-008; F13-AUTH-002 | `CONTRACT_NOW` |
| 11 | Records support field-aware redaction for audience projections without changing the final outcome | Structural projection | F13-SEC-007 | `CONTRACT_NOW` |

## Compatibility with already completed Phase 1 features

This section is a non-normative restatement of inherited-compatibility
obligations already owned by AC-13.14 (F13-COMPAT-001 through F13-COMPAT-009,
all `CONTRACT_NOW`) and by the completed "FEATURE-0011 and FEATURE-0012
dependency reconciliation (ADH-2026-017 locked)" section. FEATURE-0013
introduces no breaking changes to Phase 1 resources or endpoints. Each item
cites its owning requirement; the section creates no additional obligation.

| # | Compatibility constraint | Owning requirement | Class |
|---|---|---|---|
| 1 | Phase 1 `Operation` resource (FEATURE-0005) is unchanged; FEATURE-0013 references it by typed reference and does not redefine its schema, status, lifecycle, or API | F13-COMPAT-009; Non-goal 15 | `CONTRACT_NOW` |
| 2 | Phase 1 governance resources (Organization, OrganizationUnit, Tenant, Project) retain their API identities and are referenced only through FEATURE-0012 `scopeRef` with canonical `apiVersion` and `kind` | F13-COMPAT-007; F13-COMPAT-009 | `CONTRACT_NOW` |
| 3 | Phase 1 ServiceInstance retains its `platform.sovrunn.io/v1alpha1` identity, is referenced through `subjectRef`/`subjectRefs`, and remains a Project-scoped subject | F13-COMPAT-009 | `CONTRACT_NOW` |
| 4 | Phase 1 health and readiness endpoints are unchanged | F13-COMPAT-009 | `CONTRACT_NOW` |
| 5 | No Phase 1 route, wire format, or existing error envelope is modified | F13-COMPAT-009; F13-ERR-003 | `CONTRACT_NOW` |
| 6 | FEATURE-0012 shared grammar (metadata, typed references, Problem Details, validation) is inherited and extended, not replaced | F13-COMPAT-009 | `CONTRACT_NOW` |
| 7 | The six FEATURE-0012 `ScopeKind` values are inherited verbatim; no ScopeKind is added or removed | F13-COMPAT-007 | `CONTRACT_NOW` |

## Design questions to resolve later in design.md

1. **DQ-13-01**: Go type representation for `DecisionRecord`,
   `DecisionProfile`, `EvaluationResult`, and `AuditEvent` — struct layout,
   field types, JSON tags, and JSON Schema composition approach.

2. **DQ-13-02**: Physical JSON Schema file organization for the declarative
   profile/registry bundle format.

3. **DQ-13-03**: Package and import boundaries for decision-domain types,
   validation, and conformance.

4. **DQ-13-04**: Bounded in-memory graph representation for the
   decision-requirements DAG (data structure, traversal, short-circuit).

5. **DQ-13-05**: Validation-pass organization — ordering of scope, profile,
   evaluation, composition, trust, and obligation validators.

6. **DQ-13-06**: Pure-function interfaces for composition strategies and
   deterministic replay.

7. **DQ-13-07**: Physical fixture file format and coverage-matrix layout.

8. **DQ-13-08**: Not an open design question. The FEATURE-0011/FEATURE-0012
   dependency reconciliation is completed and locked in requirements (see the
   "FEATURE-0011 and FEATURE-0012 dependency reconciliation (ADH-2026-017
   locked)" section). It is closed and not delegated to design. Per
   F13-COMPAT-006 and that section, design MUST NOT copy, reformat, reopen,
   postpone, or re-classify the record or any completed comparison dimension
   or delta, and has no presentation, reformatting, or layout permission over
   it. Any future change requires a new approved architecture handoff.

9. **DQ-13-09**: Conformance test organization for the 28 scenario families,
   their suffixes, and additional scope/eval/sec/trust/compat cases.

10. **DQ-13-10**: Feature-gate check implementation for the downstream
    adoption contract.

## FEATURE-0011 and FEATURE-0012 dependency reconciliation (ADH-2026-017 locked)

Implementation class: `CONTRACT_NOW` (AD-015, AD-044, AD-045).

This section is the completed pre-design dependency-reconciliation record
required by architecture section 15.3 and F13-COMPAT-006. It is performed now,
in requirements, against the exact FEATURE-0011 and FEATURE-0012 sources locked
by `ADH-2026-017` `dependency_content_sha256`. Design MUST NOT create, reopen,
defer, or reformat this record; it may only render already-decided evidence.
Every delta is classified with the section 15.3 vocabulary: added, removed,
renamed, mapped, restricted, expanded, or semantically reinterpreted. All
deltas recorded here — including the approved restricted per-kind non-Platform
`scopeRef` `uid` requirement (see the `metadata.scopeRef` row below and
F13-SCOPE-002) — are approved under ADH-2026-017, so no unapproved delta
exists; any future delta requires a new approved handoff.

### Locked reconciliation sources

| Locked path | Locked `sha256` (ADH-2026-017) |
|---|---|
| `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` | `f9bd5ad9cf0aa10b06175204ed388637e64fe740fadc121d30983e00f0786623` |
| `docs/features/FEATURE-0011-reuse-assessment-standard.md` | `ad1f551a9f6f9710dd1e2a9e0ea0a0189e6240716c3fe016a44aa4f32bc34454` |
| `docs/reviews/feature-gates/FEATURE-0011-approval-review.md` | `934ac514cbe54dec22ae6d293fda3dddafb2e6af37eddf262e79b33fa06d7393` |
| `docs/reviews/architecture-decision-handoffs/ADH-2026-011-feature-0011-reuse-assessment-standard.md` | `78c28407bba69116e1c16b2a33ca5ae04ab55db5625425deca5ec3f23ed8946e` |
| `docs/architecture/api-resource-standard.md` | `b39e353324ba41c20dda07bd1c5e77bda2bbd103f2c6392f1ec1e5743cb505df` |
| `docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md` | `47c47462b81e75c2ab63fbcfba49001f5611d081b838ed881f666da20fc61b0d` |
| `.kiro/specs/api-resource-naming-status-and-validation-standard/requirements.md` | `b9e068f0d782f38f028ad3daa835d58c6ad1824bb6e70a33bb4d363b09a95d06` |
| `.kiro/specs/api-resource-naming-status-and-validation-standard/design.md` | `2e0cead4e8a82735c10bc79ba6133377b0ad6c51528f80b52fe7e03cabc50b0e` |
| `docs/reviews/feature-gates/FEATURE-0012-approval-review.md` | `6d649372550eed3c07fe5080d542073ea9a2e66b9f863d6cad64d7a81b8f2169` |
| `docs/reviews/architecture-decision-handoffs/ADH-2026-012-feature-0012-api-resource-standard.md` | `a1613acb79974efa5958fe5f5daef24a0c488c30fb9e6a39a42b8f2ae25c98d2` |
| `docs/reviews/architecture-decision-handoffs/ADH-2026-013-operation-allowed-scopes.md` | `0eef8061b631681111bea5a81a4da36281892e23bd743a3d36a0d3a9c8a40f8c` |
| `api/schemas/_common/scope-ref.json` | `53abad345e304970afb8bddb157106fe54c107849b31135fd53def50b1177f7f` |
| `api/schemas/audit-event.json` | `c60dd8e0f9b5e5388ba1982ae316d619ea3ea23fb560b0b933d6b998f94a5ef0` |
| `internal/apimeta/scope.go` | `f8b0eaaf622ff4299ad520b249f8785bb73df67396ee46dc354db07ae67ab63f` |

A `dependency_content_sha256` mismatch on any locked path invalidates this
reconciliation and downstream authorization until a new human architecture
decision updates the lock.

### Per-contract reconciliation

Each row compares the inherited contract across enums, requiredness,
cardinality, ownership, serialization, validation, security semantics,
bindings, fixtures, and migration impact, then classifies the FEATURE-0013
delta.

| Inherited contract | Source (locked) | Enum / values | Requiredness & cardinality | Ownership | Serialization | Validation | Security semantics | Bindings | Fixtures | Migration impact | FEATURE-0013 delta class |
|---|---|---|---|---|---|---|---|---|---|---|---|
| `ScopeKind` vocabulary | `scope-ref.json`; `internal/apimeta/scope.go` | Closed six-value set `Platform, Organization, OrganizationUnit, Tenant, Project, Provider` (stable enumeration order is documentation ordering, not semantic ordering or precedence) | Non-Platform scope requires exactly one canonical `ScopeKind` value | FEATURE-0012 | Enum string in `metadata.scopeRef` | FEATURE-0012 enum validation | Governance-scope isolation | Go `ScopeKind` constants | Existing scope fixtures reused | None | Reused unchanged; closed six-value set preserved; no value added, removed, renamed, or reordered as semantics |
| `metadata.scopeRef` authority & Platform form | `scope-ref.json`; `internal/apimeta/scope.go` | Platform = absent `scopeRef`; explicit `null` is not a second form | FEATURE-0012 `scopeRef` requires `apiVersion`, `kind`, and `name`; `uid` optional. Sole scope source; exactly one | FEATURE-0012 | Absence encodes Platform; canonical nil normalization | Rejects duplicate or second scope source | Prevents cross-scope disclosure | Go canonical-scope normalization | Platform and non-Platform fixtures reused | None | Reused unchanged; sole authority and required `apiVersion`/`kind`/`name` preserved. Approved restricted per-kind delta (ADH-2026-017): every non-Platform `DecisionRecord`/`AuditEvent` `scopeRef` MUST additionally carry `uid`, narrowing the FEATURE-0012 optional `uid` to required for these kinds (owned by F13-SCOPE-002). Narrowing only; no widening, and the FEATURE-0012 optional-`uid` contract is unchanged for all other kinds |
| Common `metadata` / `ObjectMeta` | `api-resource-standard.md`; FEATURE-0012 requirements | Identity and classification fields | Per FEATURE-0012 | FEATURE-0012 | JSON object | FEATURE-0012 metadata validation | Owner and identity governance | ObjectMeta binding | Metadata fixtures reused | None | Reused unchanged |
| Typed references (`scopeRef`, `ownerRef`, `subjectRef`) | `api-resource-standard.md` | Typed reference shape | Optional or required per use | FEATURE-0012 | Typed reference JSON | Reference-kind validation | No existence disclosure | Reference primitives | Reference fixtures reused | None | `subjectRef`/`subjectRefs` subject-identity use added; primitive unchanged |
| Resource profiles (`ImmutableRecord`, `VersionedDefinition`, `TransientRequestResult`, `EmbeddedValue`, `LongRunningOperation`, `ListEnvelope`) | `api-resource-standard.md` | Profile set | Profile-specific | FEATURE-0012 | Profile-defined | Profile validation | Profile boundaries | `x-sovrunn-profile` | Profile fixtures reused | None | Reused unchanged; `DecisionRecord`/`AuditEvent` bind `ImmutableRecord`, `DecisionProfile` binds `VersionedDefinition` |
| Problem Details error contract | `api-resource-standard.md`; FEATURE-0012 requirements | `type=urn:sovrunn:problem:validation-failed`, `status=422`, `code=VALIDATION_FAILED` | Required envelope | FEATURE-0012 envelope; FEATURE-0013 owns `violations[].code` names | RFC 9457 JSON; RFC 6901 pointers | Inherited validation-failed binding | Redactable, non-authoritative messages | Problem Details binding | Error fixtures reused | None to envelope | Envelope reused unchanged; closed `DECISION_*` `violations[].code` registry expanded (architecture-owned, AD-045) |
| Validation limit layering | `api-resource-standard.md` | Platform and schema ceilings | Absolute outer bounds | FEATURE-0012 ceilings; FEATURE-0013 profile and evaluator layers | N/A | Most-restrictive-applicable-value rule | Bounded cost and abuse control | Ceiling checks | Limit fixtures reused | None | FEATURE-0012 ceilings reused unchanged; profile and evaluator layers restricted (narrower) beneath them |
| `DataClassification` | `api-resource-standard.md` | FEATURE-0012 classification values | Per FEATURE-0012 | FEATURE-0012 | Enum string | FEATURE-0012 validation | Handling constraint | Classification binding | Classification fixtures reused | None | Unchanged; a separate `Sensitivity` dimension is added and mapped from it (floor only, never lowered) |
| `AuditEvent` base type and schema | `audit-event.json` | `x-sovrunn-profile: ImmutableRecord`; alpha `x-sovrunn-allowed-scopes: [Organization]` | Immutable record payload | FEATURE-0012 base type and schema preserved | ImmutableRecord JSON | FEATURE-0012 schema validation | Accountability record | `x-sovrunn-profile`, `x-sovrunn-allowed-scopes` | Organization-scoped alpha fixture retained as regression | Additive only | Preservation: the inherited `ImmutableRecord` envelope and the FEATURE-0012 `AuditEvent` base type and schema are preserved unchanged (no field removed, renamed, or reinterpreted). Approved additive deltas, distinct from that preservation: (i) the `x-sovrunn-allowed-scopes` annotation is expanded from `[Organization]` to the six governance scopes (expanded), and (ii) a FEATURE-0013-owned payload and linkage extension is added (added) |
| `Operation` contract and `LongRunningOperation` payload | `ADH-2026-013-operation-allowed-scopes.md`; FEATURE-0012 design | Operation states | Per FEATURE-0012 | FEATURE-0012 and ADH-2026-013 | FEATURE-0012 JSON | FEATURE-0012 validation | Lifecycle trace | Operation binding | Operation fixtures reused | None | Reused unchanged; referenced by typed reference for async handoff, never redefined |
| FEATURE-0011 reuse-assessment format | `PHASE2_REUSE_ASSESSMENT_STANDARD.md`; `FEATURE-0011-*` | Five-column summary; dispositions `Reuse/Wrap/Extend/Build` | One conforming feature-level summary required | FEATURE-0011 | Markdown table | Feature-gate validation | Governance evidence | Reuse-summary header | Reuse-summary fixture reused | None | Reused unchanged; FEATURE-0013 supplies one conforming summary |

### Reconciliation outcome

- No inherited FEATURE-0011 or FEATURE-0012 vocabulary, shared type, error
  contract, schema convention, or validation rule is removed, renamed, or
  semantically reinterpreted by FEATURE-0013.
- The only deltas are: (a) `subjectRef`/`subjectRefs` subject-identity use
  (added); (b) the `Sensitivity` dimension mapped from `DataClassification`
  (added and mapped); (c) the closed `DECISION_*` `violations[].code` registry
  (expanded within the unchanged Problem Details envelope); (d) profile and
  evaluator limit layers (restricted beneath inherited ceilings); (e)
  `AuditEvent` scope acceptance (expanded to the six governance scopes) plus
  its FEATURE-0013-owned payload and linkage extension (added), with the
  FEATURE-0012 `AuditEvent` `ImmutableRecord` envelope, base schema, and type
  preserved unchanged; and (f) the non-Platform `scopeRef` `uid` requirement
  for `DecisionRecord`/`AuditEvent` (an approved restricted per-kind delta
  narrowing the FEATURE-0012 optional `uid` to required for these kinds,
  owned by F13-SCOPE-002; no other kind is affected).
- The existing Organization-scoped FEATURE-0012 `AuditEvent` alpha fixture
  remains valid regression evidence; the six governance-scope fixtures and a
  ServiceInstance-subject-under-Project fixture are additive.
- All deltas are additive or restricting-only within inherited bounds and are
  backward compatible; no Phase 1 route, wire format, or existing error
  envelope changes. `ServiceInstance` is not a `ScopeKind`; it stays a
  Project-scoped typed subject.

## Architecture traceability

Controlling source: `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`, approved through `ADH-2026-017`.

This section maps every requirement family, retained design question, and
conformance obligation to the controlling architecture section. It contains
every top-level architecture section (1–28), every decision ID (AD-001
through AD-045), every risk ID (F13-R01 through F13-R31), and every stable
conformance ID. Each entry states whether it is:

- **Translated** — requirements above convert it to testable obligations;
- **Invariant** — preserved as a later invariant, carried and documented;
- **Deferred** — requires a later approved feature or decision;
- **Delegated to design** — design.md resolves representation/mechanics;
- **Governance-only** — approval, provenance, or review metadata.

### Architecture sections

| Section | Title | Disposition | Requirement families |
|---|---|---|---|
| 1 | Decision status and authority | Governance-only | Metadata; approval provenance |
| 2 | Problem statement | Translated | Introduction; all AC groups |
| 3 | Scope | Translated | AC-13.1 F13-OBJ-006; Non-goals |
| 4 | Architectural principles | Translated | All AC groups (principles are distributed across requirements) |
| 5 | Canonical object model | Translated | AC-13.1, AC-13.2, AC-13.3, AC-13.4 |
| 5.4 | FEATURE-0012 resource-profile and boundary inheritance | Translated | AC-13.1 F13-OBJ-001 through F13-OBJ-008 |
| 6 | Extensible decision semantics | Translated | AC-13.3, AC-13.4, AC-13.5 |
| 6.1 | Stable common envelope | Translated | AC-13.1 F13-OBJ-001, AC-13.2 |
| 6.5–6.7 | Finality, audit outcomes, immutability | Translated | AC-13.4 F13-IMMUT-001 through F13-IMMUT-007 |
| 6.8 | DecisionProfile contract | Translated | AC-13.5 F13-PROF-001 through F13-PROF-008 |
| 7 | Evaluation and composition model | Translated | AC-13.6 |
| 7.1 | Atomic evaluation | Translated | AC-13.6 F13-EVAL-001 through F13-EVAL-005 |
| 8 | Execution model | Translated | AC-13.1 F13-OBJ-007, F13-OBJ-008; AC-13.9 |
| 9 | Statelessness, scalability, and cost controls | Translated / Invariant | AC-13.9 F13-SCALE-001 through F13-SCALE-007 |
| 10 | Audit consistency | Translated / Invariant | AC-13.8 F13-AUDIT-001 through F13-AUDIT-007 |
| 11 | Sovereign deployment architecture | Translated | AC-13.12 F13-SOV-001 through F13-SOV-007 |
| 12 | Security, privacy, and trust boundaries | Translated | AC-13.7, AC-13.11, Security section |
| 12.1 | Runtime authentication/authorization and break-glass follow-up review duties | Translated / Invariant | AC-13.17 F13-AUTHZ-005, F13-AUTHZ-006 (carried as `INVARIANT_FOR_LATER`) |
| 12.3 | Integrity and non-repudiation | Translated / Deferred | AC-13.11 F13-TRUST-001 through F13-TRUST-005 |
| 12.4 | Provider-neutral sensitivity classification | Translated | AC-13.7 F13-SEC-001 through F13-SEC-006 |
| 13 | AI-first architecture | Translated | AC-13.7 F13-SEC-008, F13-SEC-009 |
| 14 | Reuse and standards strategy | Translated | Reuse summary |
| 15 | Compatibility and evolution | Translated | AC-13.13, AC-13.14 |
| 16 | Observability and operability | Invariant | Carried for later implementation; no FEATURE-0013 runtime |
| 17 | Conformance model | Translated | AC-13.16 |
| 18 | Matrix E v2 — architecture risk register | Governance-only | Risk controls preserved as requirements inputs |
| 19 | Decision scorecard | Governance-only | Implementation classes reproduced below |
| 20 | Required revisions before approval | Governance-only | Historical rationale |
| 21 | Important deferred decisions | Deferred | Non-goals |
| 22 | Architecture acceptance criteria | Governance-only | Approval evidence |
| 23 | Reassessment triggers | Governance-only | Later review gates |
| 24 | Normative local references | Governance-only | Context loading |
| 25 | External review inputs | Governance-only | Context loading |
| 26 | Validation against Sovrunn foundation principles | Governance-only | Conformance evidence |
| 27 | Kiro anti-overengineering guardrails | Translated | AC-13.18 |
| 27.8 | Closed architecture boundary | Translated | All AC groups (closed decisions) |
| 28 | Lightweight downstream adoption contract | Translated | AC-13.15 |

### Decision IDs (AD-001 through AD-045)

| ID | Decision | Implementation class | Disposition | Requirement reference |
|---|---|---|---|---|
| AD-001 | Enforce canonical `DecisionRecord` terminology | `CONTRACT_NOW` | Translated | AC-13.1 F13-OBJ-001; Glossary |
| AD-002 | Keep requests, evaluations, decisions, audit, operations, projections distinct; no shared DecisionRequest type | `CONTRACT_NOW` | Translated | AC-13.1 F13-OBJ-006 |
| AD-003 | Adjudication outcomes optional, facet-specific, closed | `CONTRACT_NOW` | Translated | AC-13.3 F13-ADJ-001, F13-ADJ-002 |
| AD-004 | Immutable final records with linked corrections | `CONTRACT_NOW` | Translated | AC-13.4 F13-IMMUT-001 through F13-IMMUT-007 |
| AD-005 | Deterministic versioned composition strategies | `CONTRACT_NOW` | Translated | AC-13.6 F13-COMP-001, F13-COMP-002 |
| AD-006 | Bounded DMN-like dependency DAG, not workflow language | `CONTRACT_NOW` | Translated | AC-13.6 F13-COMP-003, F13-COMP-004 |
| AD-007 | Runtime hierarchy resolution produces one EffectivePolicyContext | `INVARIANT_FOR_LATER` | Invariant | AC-13.6 F13-COMP-005 (reference only) |
| AD-008 | Synchronous fast path without mandatory queue | `INVARIANT_FOR_LATER` | Invariant | AC-13.9 F13-SCALE-005 (reference kernel only) |
| AD-009 | Operation owns async lifecycle; no pending-decision envelope | `CONTRACT_NOW` | Translated | AC-13.1 F13-OBJ-007, F13-OBJ-008 |
| AD-010 | Orchestrate authority; choreograph reactions | `INVARIANT_FOR_LATER` | Invariant | Carried in documentation |
| AD-011 | Runtime/vendor selection deferred | `DEFERRED` | Deferred | Non-goals |
| AD-012 | Pure stateless reference decision kernel | `CONTRACT_NOW` | Translated | AC-13.9 F13-SCALE-005 |
| AD-013 | Effective-once via idempotency, not exactly-once claims | `INVARIANT_FOR_LATER` | Invariant | AC-13.9 F13-SCALE-006 |
| AD-014 | Decision and audit durable-acceptance invariant | `INVARIANT_FOR_LATER` | Invariant | AC-13.8 F13-AUDIT-005, F13-AUDIT-006 |
| AD-015 | Retain six FEATURE-0012 ScopeKind; expand AuditEvent; ServiceInstance as typed subject | `CONTRACT_NOW` | Translated | AC-13.2, AC-13.8 F13-AUDIT-003, AC-13.14 |
| AD-016 | Sovereignty as seven enforceable profile dimensions | `CONTRACT_NOW` | Translated | AC-13.12 F13-SOV-001 |
| AD-017 | Connected, private, edge, and air-gapped conformance profiles | `CONTRACT_NOW` | Translated | AC-13.12 F13-SOV-002 |
| AD-018 | Static local-dependency and synchronization carrier contracts without signing | `CONTRACT_NOW` | Translated | AC-13.12 F13-SOV-004 through F13-SOV-006 |
| AD-019 | Authorized projections and data minimization | `CONTRACT_NOW` | Translated | AC-13.7 F13-SEC-007 |
| AD-020 | AI consumes DecisionContext, never canonical unrestricted data | `CONTRACT_NOW` | Translated | AC-13.7 F13-SEC-008 |
| AD-021 | AI optional; deterministic contracts work without it | `CONTRACT_NOW` | Translated | AC-13.7 F13-SEC-009 |
| AD-022 | Algorithm-agile integrity/signature carriers and structural trust states | `CONTRACT_NOW` | Translated | AC-13.11 F13-TRUST-001 through F13-TRUST-005 |
| AD-023 | Provider-neutral core with adapter-owned native identifiers | `CONTRACT_NOW` | Translated | AC-13.12 F13-SOV-003 |
| AD-024 | Hot-summary and referenced-cold-evidence carrier contract | `CONTRACT_NOW` | Translated | AC-13.7 F13-SEC-007 (reference/digest preference) |
| AD-025 | Versioned offline-capable static semantic registries | `CONTRACT_NOW` | Translated | AC-13.5 F13-PROF-004, F13-PROF-008 |
| AD-026 | Telemetry is never audit authority | `CONTRACT_NOW` | Translated | AC-13.8 F13-AUDIT-004 |
| AD-027 | Jurisdiction-specific controls supplied as profiles, not core enums | `CONTRACT_NOW` | Translated | AC-13.12 F13-SOV-007 |
| AD-028 | Contract-only Phase 2 and its feature gates | `CONTRACT_NOW` | Translated | AC-13.16, AC-13.18 |
| AD-029 | Stable envelope plus versioned DecisionProfile and typed result | `CONTRACT_NOW` | Translated | AC-13.5 F13-PROF-001 |
| AD-030 | Closed primary-form vocabulary with registered secondary facets | `CONTRACT_NOW` | Translated | AC-13.3 F13-FORM-001, F13-FORM-002 |
| AD-031 | Authority is orthogonal: authoritative, advisory, recommendation, simulation | `CONTRACT_NOW` | Translated | AC-13.3 F13-AUTH-001 through F13-AUTH-003 |
| AD-032 | Authorization is a profile; identity, evaluation, enforcement outside | `CONTRACT_NOW` | Translated | AC-13.17 F13-AUTHZ-001 |
| AD-033 | Risk-based durable-recording and aggregation profile semantics | `CONTRACT_NOW` | Translated | AC-13.17 F13-AUTHZ-002, F13-AUTHZ-003 |
| AD-034 | Immutable final records; supersession/revocation are linked effects | `CONTRACT_NOW` | Translated | AC-13.4 F13-IMMUT-002, F13-IMMUT-003 |
| AD-035 | Local atomic decision+audit-obligation acceptance; downstream audit async | `INVARIANT_FOR_LATER` | Invariant | AC-13.8 F13-AUDIT-005 through F13-AUDIT-007 |
| AD-036 | Separate retry idempotency from semantic decision identity | `CONTRACT_NOW` | Translated | AC-13.9 F13-SCALE-001, F13-SCALE-002 |
| AD-037 | Deterministic composition replays captured evaluator output | `CONTRACT_NOW` | Translated | AC-13.6 F13-EVAL-002 |
| AD-038 | Static local version-pinned profile/policy/schema/strategy and trust-state contract | `CONTRACT_NOW` | Translated | AC-13.5 F13-PROF-008 |
| AD-039 | Production authorization may use local verified artifacts/caches with bounded freshness | `INVARIANT_FOR_LATER` | Invariant | AC-13.17 F13-AUTHZ-004, F13-AUTHZ-005 |
| AD-040 | Profile and conformance bounds cover calls, bytes, concurrency, retries, jurisdictions, time, cost | `CONTRACT_NOW` | Translated | AC-13.5 F13-PROF-005, AC-13.6 F13-COMP-003 |
| AD-041 | Cryptographic assurance is tiered; remote notary/HSM not a universal hot-path dependency | `INVARIANT_FOR_LATER` | Invariant | AC-13.11 F13-TRUST-004 |
| AD-042 | Static disconnected synchronization contract with origin-aware manifests | `CONTRACT_NOW` | Translated | AC-13.12 F13-SOV-004 through F13-SOV-006 |
| AD-043 | High-risk profiles fail closed; fail-open requires bounded exception | `CONTRACT_NOW` | Translated | AC-13.5 F13-PROF-002, AC-13.10 |
| AD-044 | Downstream adoption uses one normative contract, one short section, one gate | `CONTRACT_NOW` | Translated | AC-13.15 |
| AD-045 | Bind FEATURE-0013 validation to inherited Problem Details and closed violation-code registry | `CONTRACT_NOW` | Translated | AC-13.13 |

### Risk IDs (F13-R01 through F13-R31)

All risks are preserved as requirements inputs with their architecture-
approved controls. Final residual acceptance requires implementation evidence
and human review (see architecture section 18.4–18.6).

| ID | Category | Disposition | Controlling requirement |
|---|---|---|---|
| F13-R01 | Governance/compatibility | Translated (controls) | AC-13.5 F13-PROF-003; AC-13.14 F13-COMPAT-001 |
| F13-R02 | Correctness | Translated (controls) | AC-13.5 F13-PROF-001; AC-13.3 |
| F13-R03 | Security/correctness | Translated (controls) | AC-13.3 F13-AUTH-002, F13-AUTH-003; AC-13.7 F13-SEC-008 |
| F13-R04 | Correctness (High residual) | Translated (controls) + Invariant | AC-13.6 F13-COMP-005; AC-13.9 F13-SCALE-002; later production gate |
| F13-R05 | Correctness | Translated (controls) | AC-13.9 F13-SCALE-001, F13-SCALE-002 |
| F13-R06 | Audit/compliance | Translated (controls) + Invariant | AC-13.8 F13-AUDIT-005, F13-AUDIT-006 |
| F13-R07 | Availability/latency | Translated (controls) | AC-13.8 F13-AUDIT-007 |
| F13-R08 | Cost/availability (High residual) | Translated (controls) | AC-13.6 F13-COMP-003; AC-13.9 F13-SCALE-003; later production gate |
| F13-R09 | Correctness/AI | Translated (controls) | AC-13.6 F13-EVAL-002 |
| F13-R10 | Security/AI | Translated (controls) | AC-13.7 F13-SEC-008 |
| F13-R11 | Security | Translated (controls) | AC-13.10 F13-OBLIG-002, F13-OBLIG-003 |
| F13-R12 | Latency/availability | Translated (controls) + Invariant | AC-13.17 F13-AUTHZ-004 |
| F13-R13 | Security (High residual) | Translated (controls) + Invariant | AC-13.17 F13-AUTHZ-004; later production gate |
| F13-R14 | Security/availability | Translated (controls) | AC-13.5 F13-PROF-002; AC-13.10 |
| F13-R15 | Correctness | Translated (controls) | AC-13.4 F13-IMMUT-002 through F13-IMMUT-005 |
| F13-R16 | Privacy/security | Translated (controls) | AC-13.7 F13-SEC-007; Security section |
| F13-R17 | Compliance/integrity | Translated (controls) | AC-13.4 F13-IMMUT-007 |
| F13-R18 | Latency/availability | Translated (controls) | AC-13.11 F13-TRUST-004 |
| F13-R19 | Sovereignty/security | Translated (controls) | AC-13.12 F13-SOV-004, F13-SOV-005 |
| F13-R20 | Correctness | Translated (controls) | AC-13.12 F13-SOV-006 |
| F13-R21 | Portability | Translated (controls) | AC-13.12 F13-SOV-003 |
| F13-R22 | Scope/governance | Translated (controls) | AC-13.18 F13-GUARD-003; Non-goals |
| F13-R23 | Cost | Translated (controls) | AC-13.9 F13-SCALE-003, F13-SCALE-004 |
| F13-R24 | Supply chain/security | Translated (controls) | AC-13.5 F13-PROF-008; AC-13.11 |
| F13-R25 | Correctness/security | Translated (controls) | AC-13.12 F13-SOV-007 (trusted-time provenance) |
| F13-R26 | Security/privacy | Translated (controls) | AC-13.2 F13-SCOPE-008; AC-13.6 F13-EVAL-004 |
| F13-R27 | Compatibility/governance | Translated (controls) | AC-13.13 F13-ERR-004; AC-13.14 F13-COMPAT-003 |
| F13-R28 | Compatibility | Translated (controls) | AC-13.8 F13-AUDIT-003; AC-13.14 F13-COMPAT-007 through F13-COMPAT-009 |
| F13-R29 | Traceability | Translated (controls) | AC-13.1 F13-OBJ-001; Glossary (DecisionRecord rename) |
| F13-R30 | Cost/sovereignty | Translated (controls) | AC-13.12 F13-SOV-002; AC-13.18 F13-GUARD-005 |
| F13-R31 | Governance/cost | Translated (controls) | AC-13.15 F13-ADOPT-005 |

Explicitly open High target-residual risks requiring reassessment before
production: F13-R04, F13-R08, F13-R13. These carry mandatory later-owner
and production-entry gates in addition to requirements controls.

### Conformance IDs (F13-CF-01 through F13-CF-28)

| ID | Scenario | Disposition |
|---|---|---|
| F13-CF-01 | Simplest synchronous atomic authorization decision | Translated (AC-13.16) |
| F13-CF-02 | Simplest deterministic denial with rationale and obligations | Translated (AC-13.16) |
| F13-CF-03 | Pure selection with typed result, no adjudication facet | Translated (AC-13.16) |
| F13-CF-04 | Ranking, classification, resolution, allocation, plan, assessment, recommendation, advisory, simulation profiles | Translated (AC-13.16); compound suffixes mandatory |
| F13-CF-05 | Human-approval async handoff returns Operation; eventual DecisionRecord links to it | Translated (AC-13.16) |
| F13-CF-06 | Maximally complex but bounded hierarchical composite decision | Translated (AC-13.16) |
| F13-CF-07 | Parallel evaluations with deterministic aggregation | Translated (AC-13.16) |
| F13-CF-08 | Timeout (.a), missing evidence (.b), evaluator failure (.c), policy conflict (.d) | Translated (AC-13.16); compound suffixes mandatory |
| F13-CF-09 | Idempotent replay and concurrent duplicate requests | Translated (AC-13.16) |
| F13-CF-10 | Partial decision/audit persistence failure and reconciliation | Translated (AC-13.16) |
| F13-CF-11 | Supersession (.a), correction (.b), revocation (.c), cycle (.d), chain-limit (.e) | Translated (AC-13.16); compound suffixes mandatory |
| F13-CF-12 | Unauthorized profile (.a), authority (.b), projection (.c), obligation bypass (.d), cross-scope (.e) | Translated (AC-13.16); compound suffixes mandatory |
| F13-CF-13 | Offline/air-gapped static profile-bundle, trust carriers, fail-closed states | Translated (AC-13.16) |
| F13-CF-14 | Export/import across two conforming implementations | Translated (AC-13.16) |
| F13-CF-15 | AI-disabled operation and AI projection redaction | Translated (AC-13.16) |
| F13-CF-16 | Graph size (.a), depth (.b), width (.c), fan-out (.d), timeout (.e), payload budget (.f) | Translated (AC-13.16); compound suffixes mandatory |
| F13-CF-17 | Registration of new decision family without common-envelope change | Translated (AC-13.16) |
| F13-CF-18 | Immutable supersession and revocation with calculated effective state | Translated (AC-13.16) |
| F13-CF-19 | Local atomic decision/audit-obligation while downstream audit unavailable | Translated (AC-13.16) |
| F13-CF-20 | Idempotency retry identity vs semantic decision identity across policy/profile/strategy/authority/evaluator/validity changes | Translated (AC-13.16); compound suffixes mandatory |
| F13-CF-21 | Deterministic replay from captured nondeterministic/AI output | Translated (AC-13.16) |
| F13-CF-22 | Synchronous contract with registry/resolver unavailable but valid local bundles present | Translated (AC-13.16) |
| F13-CF-23 | Structural rejection of unknown/expired/revoked/mismatched profile bundles | Translated (AC-13.16) |
| F13-CF-24 | Enforcement denial for unknown/unsupported mandatory obligations | Translated (AC-13.16) |
| F13-CF-25 | Local authorization artifact cache hit, expiry, revocation, scope mismatch, stale-policy failure | Translated (AC-13.16) |
| F13-CF-26 | Cancellation and total remote-call/byte/concurrency/retry/jurisdiction/time/cost budgets | Translated (AC-13.16); compound suffixes mandatory |
| F13-CF-27 | Structural representation for later digest/signing/notarization during remote-service outage | Translated (AC-13.16) |
| F13-CF-28 | Bounded in-memory disconnected export/import: replay (.a), duplicate (.b), ordering (.c), unknown trust root (.d), revoked origin (.e), conflict reconciliation (.f) | Translated (AC-13.16); compound suffixes mandatory |

### Compound conformance suffixes (mandatory coverage rows)

Per architecture section 17, each compound F13-CF family below has stable
lowercase suffixes in textual order. Each suffix is its own traceability and
coverage obligation; coverage of the parent requires coverage of every
suffix. All are Translated under AC-13.16 (F13-CONF-001).

| Suffix ID | Parent | Mandatory clause (textual order) | Disposition |
|---|---|---|---|
| F13-CF-04.a | F13-CF-04 | ranking profile | Translated (AC-13.16) |
| F13-CF-04.b | F13-CF-04 | classification profile | Translated (AC-13.16) |
| F13-CF-04.c | F13-CF-04 | resolution profile | Translated (AC-13.16) |
| F13-CF-04.d | F13-CF-04 | allocation profile | Translated (AC-13.16) |
| F13-CF-04.e | F13-CF-04 | plan profile | Translated (AC-13.16) |
| F13-CF-04.f | F13-CF-04 | assessment profile | Translated (AC-13.16) |
| F13-CF-04.g | F13-CF-04 | recommendation profile | Translated (AC-13.16) |
| F13-CF-04.h | F13-CF-04 | advisory profile | Translated (AC-13.16) |
| F13-CF-04.i | F13-CF-04 | simulation profile | Translated (AC-13.16) |
| F13-CF-08.a | F13-CF-08 | timeout | Translated (AC-13.16) |
| F13-CF-08.b | F13-CF-08 | missing evidence | Translated (AC-13.16) |
| F13-CF-08.c | F13-CF-08 | evaluator failure | Translated (AC-13.16) |
| F13-CF-08.d | F13-CF-08 | policy conflict | Translated (AC-13.16) |
| F13-CF-11.a | F13-CF-11 | supersession | Translated (AC-13.16) |
| F13-CF-11.b | F13-CF-11 | correction | Translated (AC-13.16) |
| F13-CF-11.c | F13-CF-11 | revocation | Translated (AC-13.16) |
| F13-CF-11.d | F13-CF-11 | cycle | Translated (AC-13.16) |
| F13-CF-11.e | F13-CF-11 | chain-limit | Translated (AC-13.16) |
| F13-CF-12.a | F13-CF-12 | unauthorized profile | Translated (AC-13.16) |
| F13-CF-12.b | F13-CF-12 | unauthorized authority | Translated (AC-13.16) |
| F13-CF-12.c | F13-CF-12 | unauthorized projection | Translated (AC-13.16) |
| F13-CF-12.d | F13-CF-12 | obligation bypass | Translated (AC-13.16) |
| F13-CF-12.e | F13-CF-12 | cross-scope reference | Translated (AC-13.16) |
| F13-CF-16.a | F13-CF-16 | graph size | Translated (AC-13.16) |
| F13-CF-16.b | F13-CF-16 | depth | Translated (AC-13.16) |
| F13-CF-16.c | F13-CF-16 | width | Translated (AC-13.16) |
| F13-CF-16.d | F13-CF-16 | fan-out | Translated (AC-13.16) |
| F13-CF-16.e | F13-CF-16 | timeout | Translated (AC-13.16) |
| F13-CF-16.f | F13-CF-16 | payload budget | Translated (AC-13.16) |
| F13-CF-20.a | F13-CF-20 | policy change | Translated (AC-13.16) |
| F13-CF-20.b | F13-CF-20 | profile change | Translated (AC-13.16) |
| F13-CF-20.c | F13-CF-20 | strategy change | Translated (AC-13.16) |
| F13-CF-20.d | F13-CF-20 | authority change | Translated (AC-13.16) |
| F13-CF-20.e | F13-CF-20 | evaluator-version change | Translated (AC-13.16) |
| F13-CF-20.f | F13-CF-20 | validity change | Translated (AC-13.16) |
| F13-CF-26.a | F13-CF-26 | cancellation | Translated (AC-13.16) |
| F13-CF-26.b | F13-CF-26 | remote-call budget | Translated (AC-13.16) |
| F13-CF-26.c | F13-CF-26 | byte budget | Translated (AC-13.16) |
| F13-CF-26.d | F13-CF-26 | concurrency budget | Translated (AC-13.16) |
| F13-CF-26.e | F13-CF-26 | retry budget | Translated (AC-13.16) |
| F13-CF-26.f | F13-CF-26 | jurisdiction budget | Translated (AC-13.16) |
| F13-CF-26.g | F13-CF-26 | time budget | Translated (AC-13.16) |
| F13-CF-26.h | F13-CF-26 | cost budget | Translated (AC-13.16) |
| F13-CF-28.a | F13-CF-28 | replay | Translated (AC-13.16) |
| F13-CF-28.b | F13-CF-28 | duplicate | Translated (AC-13.16) |
| F13-CF-28.c | F13-CF-28 | ordering | Translated (AC-13.16) |
| F13-CF-28.d | F13-CF-28 | unknown trust root | Translated (AC-13.16) |
| F13-CF-28.e | F13-CF-28 | revoked origin | Translated (AC-13.16) |
| F13-CF-28.f | F13-CF-28 | conflict reconciliation | Translated (AC-13.16) |

### Scope conformance IDs (F13-SCOPE-01 through F13-SCOPE-11)

| ID | Case | Disposition |
|---|---|---|
| F13-SCOPE-01 | Positive canonical Platform scope; metadata.scopeRef absent | Translated (AC-13.2 F13-SCOPE-002) |
| F13-SCOPE-02 | Positive canonical Organization scope | Translated (AC-13.2 F13-SCOPE-003) |
| F13-SCOPE-03 | Positive canonical OrganizationUnit scope | Translated (AC-13.2 F13-SCOPE-003) |
| F13-SCOPE-04 | Positive canonical Tenant scope | Translated (AC-13.2 F13-SCOPE-003) |
| F13-SCOPE-05 | Positive canonical Project scope | Translated (AC-13.2 F13-SCOPE-003) |
| F13-SCOPE-06 | Positive canonical Provider scope | Translated (AC-13.2 F13-SCOPE-003) |
| F13-SCOPE-07 | Absent scope rejected when Platform not permitted | Translated (AC-13.2 F13-SCOPE-001; edge case 1) |
| F13-SCOPE-08 | Explicit Platform input normalizes to canonical absent form | Translated (AC-13.2 F13-SCOPE-006) |
| F13-SCOPE-09 | Parallel, duplicate, or second scope source rejected | Translated (AC-13.2 F13-SCOPE-005) |
| F13-SCOPE-10 | Out-of-vocabulary scope (including ServiceInstance) rejected | Translated (AC-13.2 F13-SCOPE-007) |
| F13-SCOPE-11 | ServiceInstance subject accepted through subjectRef under Project scopeRef; mismatch rejected | Translated (AC-13.2 F13-SCOPE-004, F13-SCOPE-008) |

### Evaluation conformance IDs (F13-EVAL-01 through F13-EVAL-08)

| ID | Case | Disposition |
|---|---|---|
| F13-EVAL-01 | Profile sensitivity and projection limits accepted at declared bounds | Translated (AC-13.7 F13-SEC-004) |
| F13-EVAL-02 | Evaluator registration narrows profile limits | Translated (AC-13.5 F13-PROF-005) |
| F13-EVAL-03 | Evaluator widening of profile limits rejected | Translated (AC-13.5 F13-PROF-005; edge case 11) |
| F13-EVAL-04 | Missing, unknown, or incomparable mandatory limits rejected | Translated (AC-13.5 F13-PROF-006) |
| F13-EVAL-05 | Prohibited typed content category or value type rejected structurally | Translated (AC-13.7 F13-SEC-004, F13-SEC-006) |
| F13-EVAL-06 | Evaluation scope mismatch or unauthorized widening rejected without existence disclosure | Translated (AC-13.6 F13-EVAL-004) |
| F13-EVAL-07 | Opaque expected-vs-observed integrity state mismatch fails closed | Translated (AC-13.11 F13-TRUST-003) |
| F13-EVAL-08 | Arbitrary text not subjected to lexical secret/PII/malware scanning | Translated (AC-13.7 F13-SEC-006) |

### Security conformance IDs (F13-SEC-01 through F13-SEC-08)

| ID | Case | Disposition |
|---|---|---|
| F13-SEC-01 | Sensitivity ceiling boundary | Translated (AC-13.7 F13-SEC-004, F13-SEC-005) |
| F13-SEC-02 | Field-count boundary | Translated (AC-13.7 F13-SEC-004) |
| F13-SEC-03 | Byte-size boundary | Translated (AC-13.7 F13-SEC-004) |
| F13-SEC-04 | Nesting-depth boundary | Translated (AC-13.7 F13-SEC-004) |
| F13-SEC-05 | Allowed-value-type boundary | Translated (AC-13.7 F13-SEC-004) |
| F13-SEC-06 | Prohibited-category boundary | Translated (AC-13.7 F13-SEC-004) |
| F13-SEC-07 | FEATURE-0012 classification-to-sensitivity floor cannot be lowered | Translated (AC-13.7 F13-SEC-003) |
| F13-SEC-08 | Unknown, unmapped, contradictory, or below-floor classification fails closed | Translated (AC-13.7 F13-SEC-003) |

### Trust conformance IDs (F13-TRUST-01 through F13-TRUST-04)

| ID | Case | Disposition |
|---|---|---|
| F13-TRUST-01 | Algorithm-agile carrier presence and identifier syntax validated structurally | Translated (AC-13.11 F13-TRUST-001, F13-TRUST-002) |
| F13-TRUST-02 | Covered-field descriptor presence validated without interpreting or selecting it | Translated (AC-13.11 F13-TRUST-002) |
| F13-TRUST-03 | Unknown, absent, expired, revoked, mismatched, or unverified required trust state fails closed | Translated (AC-13.11 F13-TRUST-003) |
| F13-TRUST-04 | No canonicalization, content-to-digest computation, signing, verification, notarization, or cryptographic-validity claim executes before ADR-F13-002 | Translated (AC-13.11 F13-TRUST-004) |

### Compatibility conformance IDs (F13-COMPAT-01 through F13-COMPAT-09)

| ID | Case | Disposition |
|---|---|---|
| F13-COMPAT-01 | FEATURE-0012 Platform ScopeKind regression | Translated (AC-13.14 F13-COMPAT-007) |
| F13-COMPAT-02 | FEATURE-0012 Organization ScopeKind regression | Translated (AC-13.14 F13-COMPAT-007) |
| F13-COMPAT-03 | FEATURE-0012 OrganizationUnit ScopeKind regression | Translated (AC-13.14 F13-COMPAT-007) |
| F13-COMPAT-04 | FEATURE-0012 Tenant ScopeKind regression | Translated (AC-13.14 F13-COMPAT-007) |
| F13-COMPAT-05 | FEATURE-0012 Project ScopeKind regression | Translated (AC-13.14 F13-COMPAT-007) |
| F13-COMPAT-06 | FEATURE-0012 Provider ScopeKind regression | Translated (AC-13.14 F13-COMPAT-007) |
| F13-COMPAT-07 | Existing Organization-scoped FEATURE-0012 AuditEvent remains valid | Translated (AC-13.8 F13-AUDIT-003; AC-13.14 F13-COMPAT-008) |
| F13-COMPAT-08 | AuditEvent accepts all six governance scopes without changing ScopeKind vocabulary | Translated (AC-13.8 F13-AUDIT-003; AC-13.14 F13-COMPAT-007) |
| F13-COMPAT-09 | FEATURE-0012 metadata, typed-reference, validation-layer, Problem Details, extension, ServiceInstance-as-Project-scoped-resource semantics unchanged | Translated (AC-13.14 F13-COMPAT-009) |

## Architecture drift checks

| Check | Status |
|---|---|
| No provider-specific hardcoding in core | PASS — all provider-native identifiers enter only through adapters/typed references |
| No Kubernetes-only assumptions in core | PASS — no CRD, operator, or K8s-native type in FEATURE-0013 contracts |
| No PostgreSQL lifecycle logic in core placement engine | N/A — FEATURE-0013 defines no placement or provisioning logic |
| No custom policy engine embedded in handlers | PASS — policy evaluation remains behind adapters (FEATURE-0016+) |
| No raw secret storage | PASS — records prefer stable references and digests; no secret embedding |
| No customer-facing IaaS leakage | PASS — canonical objects contain no provider region/product/IaaS identifiers |
| Explainable DecisionRecord | PASS — structured rationale, reason codes, alternatives, and obligations are mandatory |
| Defined audit behavior | PASS — AC-13.8 defines audit linkage, obligation, and scope expansion |
| Preserved adapter boundaries | PASS — FEATURE-0013 defines no adapter interface; evaluator-native data stays behind adapters |
