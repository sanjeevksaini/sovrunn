---
doc_type: architecture
title: FEATURE-0013 Decision Record and AuditEvent Standard
status: approved-for-kiro-requirements
phase: 2
feature_id: FEATURE-0013
ai_load_priority: always
ai_summary: Provider-neutral, sovereign-ready architecture for extensible DecisionProfiles, immutable DecisionRecords, typed results, atomic evaluations, deterministic composition, audit accountability, and synchronous or asynchronous correlation without selecting a runtime.
---

# FEATURE-0013 — Decision Record and AuditEvent Standard

## 1. Decision status and authority

This document is the approved consolidated architecture from which Kiro may
generate FEATURE-0013 requirements. It is not an implementation specification
and does not authorize design, tasks, source code, or runtime execution; those
remain subject to their normal independent gates.

Human architecture review status: **APPROVED_FOR_KIRO_REQUIREMENTS**.

- Consolidation reviewer: Sanjeev Kumar
- Consolidation decision date: 2026-07-28
- Sole controlling handoff: ADH-2026-017
- Historical predecessor handoffs: ADH-2026-014, ADH-2026-015, and
  ADH-2026-016; provenance only after ADH-2026-017 approval

The main body consolidates every active FEATURE-0013 architecture decision and
is the sole normative architecture text for downstream requirements generation.
ADH-2026-017 is the single approved envelope for the exact consolidated
document. Historical handoffs and evidence preserve provenance but do not
create parallel semantics or downstream inputs.

The consolidated decision register contains AD-001 through AD-045. AD-001
through AD-044 preserve and reconcile the previous handoff decisions; AD-045
closes public error binding and code ownership at architecture level. The
register covers the controlled baseline corrections, contract-only Phase 2
scope, anti-overengineering guardrails, lightweight downstream adoption
contract, and the approved treatment, verification, ownership, and reassessment
plans for Matrix E F13-R01 through F13-R31. This consolidated document requires
fresh human approval and does not accept final residual risk. Final per-risk
residual acceptance remains pending implementation evidence and human semantic
review. F13-R04, F13-R08, and F13-R13 remain explicitly open at target
High residual level and require reassessment before their corresponding
production capabilities are enabled.

### 1.1 Change classification

This proposal is an **extension with two explicit corrections** to the approved
Phase 2 spine:

1. `DecisionRecord` replaces the ambiguous common name `DecisionObject`.
2. `AuditEvent` scope is expanded from the FEATURE-0012 Organization-only alpha
   profile to all six canonical FEATURE-0012 governance `ScopeKind` values:
   `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, and
   `Provider`. `metadata.scopeRef` is the sole governance-scope authority.
   `ServiceInstance` remains a Project-scoped resource and is identified by a
   typed `subjectRef` or `subjectRefs`; it is not a `ScopeKind`. This replaces
   the ServiceInstance-as-scope portion of ADH-2026-015 while preserving its
   correct retention of `Provider` and its single-scope-authority rule.

ADH-2026-017 approves both corrections and every other FEATURE-0013
architecture decision as one content-bound package. Requirements regeneration
is authorized; design, tasks, and implementation retain their separate gates.
The affected spine, baseline, RFC, feature-index, traceability, and compatibility
updates remain required before implementation proceeds.

## 2. Problem statement

Sovrunn must represent a trivial allow/deny decision and a complex,
hierarchical, multi-evaluator decision without changing the meaning of its core
objects. The representation must remain:

- deterministic and explainable;
- auditable without treating logs as audit records;
- provider, policy-engine, database, queue, and workflow neutral;
- safe for horizontally scaled and intermittently connected deployments;
- deployable in public, private, government-community, edge, and air-gapped
  environments;
- usable by humans, conventional software, and authorized AI consumers;
- bounded in latency, cost, storage, and information disclosure.

## 3. Scope

FEATURE-0013 owns only:

- the common decision and atomic-evaluation vocabulary;
- the versioned `DecisionProfile` extension mechanism and semantic registries;
- immutable `DecisionRecord` and `AuditEvent` contracts;
- deterministic composition and bounded dependency-graph semantics;
- decision, evaluation, audit, operation, trace, and request correlation;
- sovereign security, projection, portability, and evidence requirements;
- conformance fixtures, schemas, validation, and compatibility rules.

It does not implement a policy engine, placement engine, orchestration runtime,
queue, database, evidence store, identity provider, provider adapter, plugin
execution, AI agent, or production audit service.

FEATURE-0013 owns only the normalized `DecisionRecord`, `DecisionProfile`,
`EvaluationResult`, and `AuditEvent` contracts, their validation, and the
adapter-facing *data* requirements those contracts impose. It does **not**
define `EvaluatorAdapter`, `PolicyEngineAdapter`, `IdentityProviderAdapter`,
`SecretProviderAdapter`, `ObservabilityAdapter`, or any other shared adapter
interface. Those adapter interfaces and their implementations are owned by
FEATURE-0016 or the later owning feature. Where such adapter names appear in
this document they are **illustrative future consumers/producers only** and must
generate no interface or implementation artifact in FEATURE-0013. This preserves
the DEC-0036 adapter-boundary intent: evaluator-, policy-, identity-, secret-,
and observability-native data still enters only through adapters, but
FEATURE-0013 defines the normalized data contract that crosses that boundary,
not the adapter interface itself.

`DecisionRequest` is a conceptual input role, not a fifth FEATURE-0013 shared
contract. Its calling domain owns its request schema and lifecycle. When such a
request is exchanged, it must reuse the FEATURE-0012 `TransientRequestResult`
profile, common typed references, scope rules, decoding limits, validation, and
Problem Details contract. FEATURE-0013 defines no canonical shared
`DecisionRequest` schema, Go type, registry entry, persistence model, or public
API resource. A later owning feature may define a domain-specific request
without changing FEATURE-0013 when it conforms to these inherited boundaries.

## 4. Architectural principles

1. **Facts, evaluations, decisions, recommendations, actions, and accountability
   are distinct.**
2. **AI is a governed consumer, never the decision authority.**
3. **The authoritative decision is deterministic and reproducible.**
4. **The core is stateless; durable state is external and replaceable.**
5. **Sovereignty is policy plus technical enforcement, not a provider label.**
6. **Every graph, collection, payload, timeout, retry, and retention is bounded.**
7. **Authoritative composition is orchestrated; reactions may be choreographed.**
8. **Portable contracts precede runtime and vendor selection.**
9. **Disclosure is least-privilege; completeness is preserved in governance
   custody, not in every projection.**
10. **Decision families extend profiles and typed results, never the common
    envelope.**
11. **Reuse is mandatory, but external native types never become Sovrunn core
    domain types.**

## 5. Canonical object model

### 5.1 Objects and ownership

| Object | Meaning | Lifecycle owner | Authoritative? |
|---|---|---|---|
| `DecisionRequest` | Conceptual bounded input role with context references and constraints; not a FEATURE-0013 shared schema or type | Calling domain | No |
| `EffectivePolicyContext` | Resolved, read-only hierarchical governance input | FEATURE-0020 | Input authority |
| `EvaluationResult` | Atomic, evaluator-specific normalized finding | Evaluator adapter/domain | Evidence, not final decision |
| `DecisionProfile` | Versioned contract defining one decision family's form, authority, inputs, typed result, rationale, obligations, projections, and limits | FEATURE-0013 registry; adopting domain supplies profiles | Definition authority |
| `DecisionRecord` | Immutable governed conclusion, with explicit authority, rationale, inputs, strategy, and typed result | Decision service/domain | Profile-declared |
| `DecisionRationale` | Structured reasons, alternatives, obligations, warnings, and corrective suggestions embedded in the decision | Decision service/domain | Part of decision |
| `AuditEvent` | Immutable accountability record of actor/action/subject/outcome | Audit boundary | Yes |
| `Operation` | Lifecycle of an asynchronous action or attempted change | Operation controller | Yes for operation state |
| `DecisionContext` | Sanitized, policy-governed projection for AI and human explanation | FEATURE-0025 | No |

### 5.2 Why `DecisionRecord`

`DecisionObject` says only that the value is an object. `DecisionResult` is too
easy to confuse with atomic evaluator output and may imply a transient response.
`DecisionAudit` incorrectly merges decision authority with accountability, while
`DecisionExplain` describes only a projection. `DecisionRecord` communicates a
durable, immutable, explicitly authorized governed conclusion and remains distinct from
`EvaluationResult`, `AuditEvent`, and `DecisionContext`.

Specialized records such as `PlacementDecision` may compose the common
`DecisionRecord` profile but must not redefine common semantics. New decision
families register a `DecisionProfile` and typed result schema; they do not add
fields to or fork the common envelope.

### 5.3 Required semantic distinctions

- An `EvaluationResult` reports one evaluator's finding.
- A `DecisionProfile` defines the governed semantics and schema of one decision
  family without changing the common record.
- A `DecisionRecord` captures one governed authoritative, advisory,
  recommendation, or simulation conclusion from zero or more evaluations and
  governed facts.
- An `AuditEvent` records accountability for an occurrence.
- An `Operation` tracks asynchronous progress.
- A `DecisionContext` is a filtered explanation projection.
- Telemetry describes system behavior and is never an audit authority.

### 5.4 FEATURE-0012 resource-profile and boundary inheritance

FEATURE-0013 does not choose new object shapes. Each externally exchanged
object inherits exactly one FEATURE-0012 resource profile and boundary rule:

| FEATURE-0013 object | FEATURE-0012 profile | Canonical boundary and persistence rule |
|---|---|---|
| `DecisionProfile` and semantic-registry bundle | `VersionedDefinition` | Governance-only publication contract; published versions are immutable; runtime services are deferred |
| `DecisionRequest` | `TransientRequestResult` | Calling-domain-owned internal-engine-facing input; FEATURE-0013 defines no shared schema or type and it is never persisted merely to satisfy resource grammar |
| `EvaluationResult` | `EmbeddedValue` when captured inside a record; `TransientRequestResult` when exchanged independently | Adapter/internal-engine-facing evidence; no independent durable lifecycle is created by FEATURE-0013 |
| `DecisionRecord` | `ImmutableRecord` | Governance-only canonical record; append-only; customer, operator, auditor, regulator, provider, and AI access occurs through authorized projections |
| `AuditEvent` | `ImmutableRecord` | Governance-only canonical accountability record; append-only; telemetry and transport events are not substitutes |
| `Operation` | `LongRunningOperation` | Inherited unchanged from FEATURE-0012/ADH-2026-013; owns asynchronous lifecycle, not decision meaning |
| `DecisionContext` | `TransientRequestResult` | Authorized audience-specific projection; non-authoritative and not persisted merely to satisfy resource grammar |

`DecisionRecord` and `AuditEvent` use FEATURE-0012's `ImmutableRecord` shape:
type metadata, common `metadata`, and one immutable record payload. They do not
use mutable desired-state `spec`, observed `status`, or a second scope field.
The payload fields shown in section 6.1 are the record payload and are not an
alternative resource grammar. `DecisionProfile` uses the FEATURE-0012
`VersionedDefinition` shape; its definition content is `spec`, while publication
status, if present, follows FEATURE-0012 ownership and mutability rules.

Canonical governance records are never made customer-facing merely by exposing
their schema. Audience-specific disclosure is a projection decision governed by
the profile, sensitivity mapping, purpose, authorization, and minimization
rules. Requirements and design may not select a different resource profile,
introduce a parallel envelope, or turn an embedded/transient value into a
persistent resource without a new architecture decision.

## 6. Extensible decision semantics

### 6.1 Stable common envelope

Every decision uses the same `DecisionRecord` envelope. It identifies the
versioned decision profile, primary form, authority, purpose, request, actors,
subjects, scope, inputs, evaluations, composition, typed result, rationale,
obligations, provenance, validity, correlation, and audit links.

`DecisionRecord` scope identity is expressed solely through the FEATURE-0012
`metadata.scopeRef` field. `metadata.scopeRef` is the sole
logical and serialized scope authority; a top-level `scopeRef` field, a
parallel scope enum, a duplicated scope attribute, or any second scope source
is removed and prohibited. The canonical `Platform` form is an absent/nil
`metadata.scopeRef` (inherited from FEATURE-0012 `NormalizeScope` and
`CanonicalScopeIdentity`) where the applicable contract permits
`Platform`; every non-Platform scope carries a non-nil `metadata.scopeRef` with
a canonical `ScopeKind` and UID. A record that declares scope through more than
one source, or through any source other than `metadata.scopeRef`, fails
validation. A canonical nil `metadata.scopeRef` that resolves to `Platform` is
not a contradictory or duplicate scope.

The envelope does not contain every domain's result fields. `DecisionProfile`
selects a separately versioned typed-result schema. An adopting feature may add
a profile and result schema without revising FEATURE-0013 when it conforms to
the extension contract and existing semantic registries.

```yaml
apiVersion: governance.sovrunn.io/v1alpha1
kind: DecisionRecord
metadata:
  name: placement-decision-123
  uid: decision-123
  # scopeRef is omitted only for canonical Platform scope. A non-Platform
  # record carries the one canonical FEATURE-0012 ScopeRef object here.
  scopeRef:
    apiVersion: core.sovrunn.io/v1alpha1
    kind: Project
    name: payments-production
    uid: project-123
record:
  profileRef:
    name: PlacementDecision
    version: 1.0.0
  form: SELECTION
  authority: AUTHORITATIVE
  purpose: workload-placement
  subjectRefs:
    - apiVersion: platform.sovrunn.io/v1alpha1
      kind: ServiceInstance
      name: postgres-primary
      uid: service-456
  evaluationResultRefs: []
  result:
    finality: FINAL
    typedResult: {}
    rationale: {}
    obligations: []
  provenance: {}
  validity: {}
  correlation:
    auditEventRefs: []
```

The example follows FEATURE-0012 reference identity and absence rules:

- `Project` retains its approved `core.sovrunn.io/v1alpha1` API identity;
- `ServiceInstance` retains its existing
  `platform.sovrunn.io/v1alpha1` compatibility identity until a separately
  approved migration changes it;
- every present typed reference contains `apiVersion`, `kind`, and `name`, with
  optional immutable `uid`; and
- absent optional singular references are omitted. They are not serialized as
  `null` unless a future approved canonical schema explicitly declares a
  nullable union. Empty reference collections may be represented as empty
  arrays when the owning profile permits them.

### 6.2 Decision forms

The closed, slowly evolving primary-form vocabulary is:

- `ADJUDICATION` — allow, deny, or require approval;
- `SELECTION` — choose one or more feasible candidates;
- `RANKING` — order candidates without necessarily selecting one;
- `CLASSIFICATION` — assign a governed category;
- `RESOLUTION` — derive an effective value from multiple inputs;
- `ALLOCATION` — distribute bounded capacity or entitlement;
- `PLAN` — produce a proposed bounded sequence or target state;
- `ASSESSMENT` — determine risk, compliance, sufficiency, or compatibility;
- `RECOMMENDATION` — provide an explicitly non-authoritative suggestion.

A profile declares exactly one primary form and may include registered
secondary form facets. For example, placement can adjudicate feasibility and
then select a target, while its primary form remains `SELECTION`.

### 6.3 Decision authority

Authority is independent of decision family and form:

- `AUTHORITATIVE` — controls a Sovrunn outcome;
- `ADVISORY` — informs another authoritative decision;
- `RECOMMENDATION` — proposes an action for governed review;
- `SIMULATION` — counterfactual, dry-run, or what-if result.

Only an authorized profile and actor may produce an authoritative record. AI
output is never authoritative. Advisory, recommendation, and simulation records
must be visibly distinguishable in every projection and must not be consumed as
enforcement authority.

### 6.4 Adjudication outcome

Only profiles with an adjudication facet use the closed vocabulary:

- `ALLOWED`
- `DENIED`
- `REQUIRES_APPROVAL`

It is mandatory for authorization, admission, governance, approval, and similar
profiles, but not for pure selection, ranking, classification, resolution,
allocation, plan, assessment, or recommendation profiles. Those profiles carry
their answer in the registered typed result.

`UNKNOWN`, runtime failure, timeout, evaluator error, and insufficient evidence
are evaluation or production states—not business adjudications. A registered
fail-safe strategy may deterministically map them to denial or required approval
and must record that mapping.

### 6.5 Finality and effective state

Every successfully persisted `DecisionRecord` is immutable and `FINAL`. Pending,
running, retrying, and waiting belong exclusively to `Operation`; they are never
decision states.

Correction, supersession, and revocation create a new final record whose typed
relationship declares `CORRECTS`, `SUPERSEDES`, or `REVOKES` and references its
predecessor. The predecessor is never mutated. A read model may calculate
`effectiveState: EFFECTIVE | SUPERSEDED | REVOKED`, but calculated effective
state is projection metadata and is not stored by rewriting the canonical
record. Relationship chains are bounded, cycle-free, and deterministically
resolved.

### 6.6 Audit outcomes

Audit outcomes are separate from decision outcomes:

- `Succeeded`
- `Denied`
- `Failed`

For example, successfully producing a `DENIED` adjudication is an audit outcome
of `Succeeded`, not `Denied`. Audit `Denied` means the attempted audited action
was refused.

### 6.7 Immutability and correction

Final `DecisionRecord` and `AuditEvent` resources are append-only. Corrections
or revocations create new records linked by typed `supersedes`, `corrects`, or
`revokes` references. Chains are bounded and cycle-free. Retention or lawful
erasure never rewrites the canonical immutable record. A later approved
architecture may use separately controlled payload custody, cryptographic key
destruction, or a linked tombstone/redaction record that preserves identity and
chain integrity where law permits. FEATURE-0013 defines only the structural
state and linkage contract; it implements no erasure or cryptographic process.

### 6.8 DecisionProfile contract

Every registered profile declares:

- stable name, family, semantic version, owner, and lifecycle status;
- primary form, allowed secondary facets, and allowed authority levels;
- input, typed-result, rationale, and obligation schemas;
- accepted evaluation types and composition strategies;
- required scope, actor, subject, policy, evidence, and audit semantics;
- determinism, failure, timeout, and insufficient-evidence behavior;
- validity, expiry, correction, supersession, and revocation rules;
- sensitivity, residency, retention, legal-hold, and projection profiles;
- idempotency identity and bounded cost, graph, latency, and payload limits;
- compatibility policy, fixtures, owner, and reassessment triggers.

For security, authorization, admission, sovereignty, compliance, privileged
access, data movement, and break-glass profiles, fail-closed or
`REQUIRES_APPROVAL` is mandatory. Fail-open is prohibited unless a separately
approved, time-bounded `SecurityExceptionRef` is present in the governing
profile. `SecurityExceptionRef` is a structural approval-evidence reference,
not a runtime workflow: it must identify the approved exception artifact by
`apiVersion`, `kind`, `name`, and immutable `uid`; identify the exception scope,
owner, approving authority, compensating controls, effective interval, audit
treatment, reassessment trigger, and covered failure modes; and be validated
only as bounded data. FEATURE-0013 does not approve, issue, revoke, execute, or
store security exceptions. Missing, expired, out-of-scope, owner-mismatched,
malformed, or non-uid-pinned exception evidence makes fail-open invalid and the
profile fails closed. Design and tasks must not invent alternate approval
fields or workflows.

Profiles and semantic registries must be distributable and verifiable offline.
Registration requires conformance review; accepting arbitrary profile URIs at
runtime would create an ungoverned extension and is prohibited.

The canonical registry representation is a declarative, versioned, provider-
and language-neutral JSON bundle governed by JSON Schema. Go registry types and
static in-memory maps are derivative implementation views and must not become
the source of semantic truth. A bundle contains profile and semantic-registry
entries plus algorithm-agile integrity and trust carriers; it does not select
or execute canonicalization, digest, signing, or verification algorithms.

FEATURE-0012 platform/schema validation ceilings are absolute outer safety
bounds. Within those inherited ceilings, profile declarations own the maximum
permitted sensitivity, content categories, field count, bytes, depth, value
types, graph, latency, and payload bounds for their decision family. Evaluator
registrations may narrow those profile limits for one evaluator type but may
never widen, replace, or contradict either the profile or FEATURE-0012 ceiling.
The effective bound is the most restrictive applicable value. Missing, unknown,
incomparable, or above-ceiling mandatory bounds fail validation.

Registry publication and approval are later control-plane operations. Runtime
workers may eventually load locally available, versioned profile bundles and
pin profile name, version, schema identity, and opaque trust-root identity in a
decision. A synchronous decision must not require a network lookup. Before
ADR-F13-002, FEATURE-0013 validates only static bundle structure, identifier
syntax, carrier presence, opaque expected-versus-observed states, and fail-
closed trust-state transitions. Unknown, expired, revoked, mismatched, or
unverified required trust states fail closed; no cryptographic-validity claim
is made.

FEATURE-0013 defines `DecisionProfile` schemas, declarative bundle format,
governance, structural validation, and conformance only. It does not implement
or authorize a production registry, distribution service, control plane,
canonicalizer, digester, signer, verifier, key service, or revocation service.

### 6.9 Anticipated decision families

The extension model must support at least these families without changing the
common envelope:

| Domain | Illustrative profiles |
|---|---|
| Identity and access | Authorization, authentication-assurance, entitlement, admission, separation-of-duties |
| Governance and compliance | Governance, compliance, exception/waiver, approval, evidence-sufficiency |
| Security and trust | Security-risk, artifact-trust, supply-chain, secret-disclosure, break-glass |
| Placement and capacity | Placement, provider/resource selection, scheduling, capacity allocation, quota |
| Sovereignty and data | Residency, data-movement, disclosure/export, retention/disposition |
| Service lifecycle | Plan resolution, configuration resolution, lifecycle transition, provisioning strategy, upgrade/migration |
| Economics | Cost, budget, optimization, charge-policy assessment |
| Availability | Scaling, traffic routing, failover, resilience and recovery strategy |
| Compatibility | Version compatibility, capability matching, migration readiness |
| AI governance | AI-use authorization, model selection, context disclosure, non-authoritative recommendation and remediation |

This table validates extensibility; it does not authorize implementation of
these domains in FEATURE-0013.

### 6.10 Authorization specialization

`AuthorizationDecision` is a profile of `DecisionRecord`, not a separate core
object. It requires an `ADJUDICATION` facet and captures principal, delegation,
action, resource, scope, assurance, policy context, evaluations, rationale,
validity, and enforceable obligations.

Authentication, identity resolution, policy authoring, policy evaluation, and
enforcement remain outside FEATURE-0013. An enforcement point must enforce all
mandatory obligations or deny the action. An unknown, invalid, expired, or
unsupported mandatory obligation makes an `ALLOWED` decision unenforceable and
therefore must result in denial at that enforcement point. Recording policy is risk-based:
sensitive, denied, privileged, cross-scope, break-glass, export, and policy-
change decisions require durable records; high-volume low-risk checks may use a
governed aggregation profile rather than persisting every check as a full
decision.

A low-risk check that is not persisted as a `DecisionRecord` is an ephemeral
authorization evaluation or enforcement event, not a durable Sovrunn decision.
Its audit evidence remains governed by the authorization profile. Denials,
privileged access, cross-scope or cross-tenant access, break-glass access, data
disclosure/export, policy or privilege changes, human approval, and regulated
actions always require a durable `DecisionRecord` and audit obligation.

FEATURE-0013 does not require a remote central authorization call for every
request. A later enforcement architecture may use local policy decision points,
version-pinned context, later-approved verified short-lived decision artifacts,
or safe caches
when the profile defines audience, scope, expiry, revocation, freshness, and
fail-closed behavior. Aggregation may reduce durable `DecisionRecord` volume but
must not suppress audit evidence required by the applicable profile.

## 7. Evaluation and composition model

### 7.1 Atomic evaluation

An atomic evaluation has one registered evaluator type, one bounded input
snapshot reference or opaque integrity carrier, one result, one timing envelope, and explicit success,
non-match, indeterminate, timeout, or error semantics. Evaluator-native payloads
remain behind adapters; normalized fields enter the common record.

Evaluation reproducibility and decision reproducibility are distinct. An
external, probabilistic, or AI evaluator may not reproduce the same output when
rerun. Its normalized result must therefore capture evaluator identity and
version, governed input snapshot reference and optional opaque integrity carrier, returned output, relevant execution
configuration, evaluation time, trust boundary, safety/policy filters, and
structural trust state. Authoritative composition is deterministic from these
captured results and must not rerun an external or AI evaluator during replay or
audit. A profile may prohibit nondeterministic evaluators entirely.

`EvaluationResult` does not establish an independent governance scope. Its
accepted scope is derived from the containing `DecisionRecord`'s canonical
`metadata.scopeRef` and the profile's declared cross-scope rules. If an
evaluator supplies scope metadata for structural comparison, that metadata is
evidence only: it must equal or be an explicitly permitted narrowing of the
record scope and must never become a second scope authority. Profile rules own
the permitted cross-scope relationship; evaluator registrations may narrow but
not widen it. Cross-scope validation occurs before composition, and a mismatch,
unknown scope, unauthorized widening, or inaccessible reference fails closed
without disclosing whether the referenced target exists. Design may choose the
validation-pass organization but may not change these authority semantics.

### 7.2 Composite decision

A composite decision combines bounded evaluation results using a versioned,
registered strategy. Every strategy declares:

- stable identifier and version;
- accepted input types;
- ordering and precedence;
- short-circuit behavior;
- missing, timeout, conflict, and error behavior within the governing
  profile's fail-closed/default posture;
- any requested fail-open behavior only as a reference to the profile's valid
  `SecurityExceptionRef`; strategy metadata itself is never an independent
  fail-open authority or approval source;
- deterministic output mapping;
- resource budgets and maximum fan-out;
- explanation and evidence rules.

Fail-open remains subject to the profile restrictions and structural
`SecurityExceptionRef` contract in section 6.8 and can never be selected merely
as a runtime fallback, by strategy metadata alone, or by an implementation
choice.

The chosen strategy identifier, version, ordered input identities, optional
opaque integrity carriers, and output are recorded so the result can be
reproduced without requiring digest computation in FEATURE-0013.

### 7.3 Hierarchical evaluation

Governance inheritance is resolved before final decision composition into one
`EffectivePolicyContext`. The decision layer must not independently traverse
Platform, Organization, OrganizationUnit, Tenant, Project, or Provider trees.
The resolved context includes provenance, precedence, conflict resolution,
effective time, and versioned identity with an optional opaque integrity
carrier. This prevents every decision domain from creating a
different hidden inheritance engine.

A valid immutable context snapshot may be resolved ahead of the request and
cached by scope, version, effective-time window, and complete semantic identity.
An optional opaque integrity carrier may accompany the snapshot. The synchronous
path must not traverse the governance hierarchy or call a central resolver when
a policy-valid local snapshot is available. Profiles define maximum staleness,
invalidation, revocation, and fail-closed behavior.

### 7.4 Decision dependency graph

Complex composition uses a bounded, directed acyclic decision-requirements
graph inspired by DMN concepts. It is not a general workflow language.

The graph must have:

- stable node and edge identifiers;
- typed nodes for input, evaluation, composition, and decision;
- a single root decision whose authority is explicit;
- deterministic topological ordering;
- declared parallelizable groups;
- maximum nodes, depth, width, fan-out, payload, and wall-clock budget;
- maximum remote calls, unique authorities, bytes fetched, concurrent
  evaluations, per-evaluator timeout, retries, cross-jurisdiction calls, and
  cost units;
- cycle, duplicate, missing dependency, and incompatible-version rejection;
- recorded graph definition/version and optional opaque integrity carrier.

The graph may describe dependency and composition only. It cannot execute
provisioning, human tasks, compensations, arbitrary code, or unbounded loops.
Short-circuiting must cancel or disregard unnecessary work without changing
deterministic results, leaking cross-scope information, or leaving unbounded
background activity.

## 8. Execution model

### 8.1 Synchronous path

Simple and bounded decisions should complete synchronously when all required
inputs are available within the declared latency budget. The path does not
require a queue or durable workflow engine. Its externally exchanged outcome
uses only existing FEATURE-0012/0013 contracts:

- completion returns a `DecisionRecord` value or canonical typed reference;
- accepted asynchronous handling returns an existing `Operation` value or
  canonical typed reference conforming to FEATURE-0012 `LongRunningOperation`;
  and
- rejection, malformed input, or inability to accept either path returns the
  inherited FEATURE-0012 Problem Details envelope.

FEATURE-0013 defines no `PendingDecision`, `DeferredDecision`,
`DecisionPending`, pending-decision status vocabulary, or other non-final
decision response envelope. Transition to `Operation` is the sole accepted
asynchronous handoff; it does not make the decision mutable and does not imply
that a final `DecisionRecord` already exists.

### 8.2 Asynchronous path

When evaluation needs human approval, unavailable external evidence, long
running work, or a budget exceeding the synchronous contract, an `Operation`
owns lifecycle, retries, deadlines, cancellation, and status. The eventual
`DecisionRecord` links to the `Operation`; the decision itself never becomes a
mutable workflow-state object.

### 8.3 Orchestration and choreography

- **Orchestration is required** for authoritative evaluation ordering,
  composition, retry policy, deadlines, human approval, and finalization.
- **Choreography is allowed** after durable decision finalization for
  notifications, indexing, telemetry, cache invalidation, and other
  non-authoritative reactions.
- An event consumer must never silently alter or replace the authoritative
  decision.

### 8.4 Runtime neutrality

FEATURE-0013 does not select a durable-execution system, workflow engine,
message broker, queue, event-streaming platform, or database. A later runtime
decision must satisfy the contracts here and demonstrate sovereign operability,
horizontal scaling, cost, skills, portability, disconnected operation,
recovery, and licensing fit.

### 8.5 Data-plane locality invariant

All profile, schema, strategy, policy-context, trust, and evaluator material
needed for a bounded synchronous decision must be locally available,
version-pinned, and structurally verifiable. Later approved cryptographic
profiles may additionally require digest verification. Profile publication, policy resolution,
schema distribution, evidence export, audit indexing, AI inference,
notarization, and regulator replication are control-plane or asynchronous
activities unless a specific profile explicitly makes one part of the decision
input and budgets its failure behavior.

No remote registry, audit index, evidence exporter, AI model, identity control
plane, cryptographic notary, or vendor control plane may become an implicit
universal synchronous dependency. When a mandatory input is missing, the only
permitted outcomes are: a final fail-closed `DecisionRecord` when the profile
defines a deterministic adjudication; transition to the existing `Operation`
contract when asynchronous handling is accepted; or inherited FEATURE-0012
Problem Details when neither path can be accepted. No other non-final response
type or status is permitted, and processing never waits synchronously without
a bound.

## 9. Statelessness, scalability, and cost controls

Implementation classification for this section is closed by architecture:

- `CONTRACT_NOW`: schema fields for retry identity and the complete semantic
  decision-identity descriptor, profile-declared finite bounds, pure or bounded
  in-memory conformance, and documentation of the later invariants;
- `INVARIANT_FOR_LATER`: stateless interchangeable workers, external durable
  state behind ports, compare-and-create semantics, effective-once behavior,
  bounded delivery/retry/backpressure/cache behavior, and runtime resource
  budgets; and
- `DEFERRED`: every concrete store, transaction, lock, lease, queue,
  dead-letter processor, cache, worker topology, scaling policy, and production
  numeric default.

The `CONTRACT_NOW` reference kernel is a pure, deterministic function of
versioned inputs, registered strategy, and explicit fixture-local bounded
configuration. It holds no authoritative session state and performs no
production persistence, delivery, locking, caching, or scaling behavior.

A later production implementation is constrained by the following
`INVARIANT_FOR_LATER` obligations:

- a caller-supplied idempotency key scoped to caller and request purpose for
  retry identity;
- a separate complete semantic decision identity descriptor covering scope,
  purpose, profile name and version, authority, canonical input snapshot,
  `EffectivePolicyContext`,
  composition strategy/version, evaluator contract versions, and relevant
  evaluation-time or validity epoch;
- compare-and-create or equivalent concurrency control;
- at-least-once delivery without duplicate authoritative decisions;
- bounded retries with jitter and dead-letter/quarantine handling;
- backpressure, admission control, per-scope quotas, and load shedding;
- deterministic caching only by the complete semantic decision identity, with
  audience, scope, expiry, revocation, and freshness enforcement;
- no dependence on sticky sessions;
- scale-to-zero where deployment policy permits;
- hot-path summaries with cold evidence referenced by typed identity and,
  where externally supplied, an opaque integrity carrier;
- batch and parallel evaluation only within declared resource budgets.

Exactly-once transport is not claimed. In a later production implementation,
business-level effective-once behavior comes from idempotency, immutable
identity, and atomic persistence boundaries. FEATURE-0013 implements only the
carrier fields, validation, and deterministic conformance needed to preserve
that future invariant.

FEATURE-0013 deliberately defines no new universal numeric decision-domain
defaults or system-wide graph ceiling. Existing FEATURE-0012 decoder, schema,
and validation ceilings remain inherited absolute outer bounds. Every
conforming profile must declare finite values for all domain limits it uses
within those ceilings, and evaluator registrations may only narrow them. A
later platform-governance feature may approve additional deployment-wide
decision ceilings. Until then, requirements, design, fixtures, and code must
not invent a FEATURE-0013 default, fallback, or hidden maximum; conformance uses
explicit fixture-local profile values solely to test boundary behavior.

## 10. Audit consistency

Implementation classification for this section is closed by architecture:

- `CONTRACT_NOW`: immutable `DecisionRecord`/`AuditEvent` identity and linkage,
  the audit-obligation carrier, correlation fields, validation, and pure or
  bounded in-memory partial-failure/reconciliation conformance;
- `INVARIANT_FOR_LATER`: atomic durable acceptance, asynchronous
  materialization, bounded delivery, recovery, quarantine, and reconciliation;
  and
- `DEFERRED`: the transaction, outbox, event-store, database, queue, index,
  exporter, SIEM, replication, and regulator-delivery mechanisms.

A final decision must not be reported as durably recorded until:

1. its immutable `DecisionRecord`; and
2. its required audit obligation

are atomically accepted by the same local authoritative persistence boundary.
Durable audit-obligation acceptance is sufficient for the synchronous response;
full `AuditEvent` materialization, indexing, export, notification, SIEM delivery,
cross-site replication, or regulator delivery may proceed asynchronously.

The atomic boundary preallocates the immutable `AuditEvent` identity and stores
that identity in both the `DecisionRecord` reference and audit obligation.
Asynchronous materialization must create the `AuditEvent` using the same
identity and must never require mutation of the `DecisionRecord`.

The concrete local mechanism—single transaction, transactional outbox,
event-sourced append, or another atomic pattern—is deferred. A later
implementation's synchronous path must not depend on a remote audit service or
index and must expose recovery, bounded delivery, poison-record quarantine, and
reconciliation for partial downstream failure. FEATURE-0013 represents and
tests these failure semantics without implementing any such runtime mechanism.
No implementation may fabricate durable success from a telemetry log.

Every meaningful decision links to at least one `AuditEvent`. A single audit
event may reference a decision, operation, actor, subject, resource, request,
trace, and parent event without copying unrestricted payloads.

## 11. Sovereign deployment architecture

### 11.1 Sovereignty dimensions

Sovrunn treats sovereignty as independently enforceable dimensions:

| Dimension | Required architectural control |
|---|---|
| Data | Residency tags, permitted jurisdictions, data-classification-aware routing, export controls |
| Operational | In-jurisdiction administrators, local break-glass, local monitoring and support, operator-owned runbooks |
| Technical | Deployable without foreign control plane, local identity/KMS/DNS/NTP/registry dependencies, open export formats |
| Legal | Policy references, lawful basis/authority, retention class, legal-hold and disclosure records |
| Cryptographic | Customer/government-controlled keys, rotation, revocation, crypto-agility, optional HSM custody |
| Software supply chain | Offline-verifiable signed artifacts, SBOM, provenance, vulnerability policy, mirrored registries |
| AI/model | Model/data locality, approved model registry, prompt/output classification, no undeclared external inference |

### 11.2 Deployment profiles

The same canonical contracts must operate in:

- connected public or sovereign cloud;
- private cloud or Government Community Cloud;
- on-premises department or regulated-enterprise cloud;
- intermittently connected edge/local zone;
- fully disconnected or air-gapped deployment.

Disconnected profiles use local trust roots, identity, time source, artifact
registry, schema registry, policy bundles, keys, observability, audit store, and
AI inference if AI is enabled. FEATURE-0013 defines only the static,
origin-aware synchronization manifest, policy-authorization fields, ordering,
duplicate/replay markers, quarantine states, and explicit conflict evidence as
`CONTRACT_NOW`. Resumable synchronization execution and reconciliation are
`INVARIANT_FOR_LATER`. Manifest signing, signature verification, trust-root
cryptography, and all algorithms/products/services are `DEFERRED` under
ADR-F13-002. Synchronization is never a hidden foreign control-plane
dependency.

Export and import use collision-resistant record identifiers, origin trust-
domain identity, algorithm-agile signature carriers, replay markers, duplicate
detection, per-origin ordering, and explicit opaque trust-root identity. The
manifest is not cryptographically signed or verified by FEATURE-0013. Unknown
profiles, origins, or opaque trust states—including unknown, absent, expired,
revoked, mismatched, or unverified states—fail closed or are represented as
quarantined according to the profile. Immutable records are never merged using
last-writer-wins; cross-origin conflicts produce explicit evidence and a
governed reconciliation decision. Runtime quarantine and reconciliation are
later invariants, not FEATURE-0013 services.

### 11.3 Provider neutrality and deployment portability

No canonical object contains a cloud provider, infrastructure platform,
container orchestrator, AI service, product identifier, region code, or
provider-native resource type as a required semantic. Provider-specific
topology, capacity, jurisdiction, trust, and capability information enters only
through adapters, normalized capability profiles, and typed references.

The same decision and audit model must remain portable across public cloud,
private cloud, Government Community Cloud, regulated sovereign cloud, local
cloud providers, department-operated infrastructure, and on-premises
environments. Provider branding or self-attestation does not establish
sovereign conformance; conformance is demonstrated through verifiable
capabilities, deployment controls, evidence, and policy enforcement.

### 11.4 Government and regulated controls

The contracts carry or reference, without embedding secrets:

- governing jurisdiction and residency constraints;
- data classification and handling caveats;
- authority/purpose and policy version;
- tenant/department/mission boundary;
- retention class, legal hold, and disposition status;
- actor identity, workload identity, delegation chain, and authentication
  assurance;
- key domain and cryptographic proof references;
- approved disclosure and export policy;
- opaque integrity and trusted-time provenance carriers, with verification
  deferred to later approved cryptographic architecture;
- emergency/break-glass use with mandatory reason and follow-up review.

Jurisdiction-specific values belong in versioned profiles and policy
registries, not hard-coded core enums. This allows adoption across sovereign
states and regulated sectors.

## 12. Security, privacy, and trust boundaries

### 12.1 Zero-trust rules

- Authenticate and authorize every human, workload, evaluator, and service hop.
- Bind authorization to scope, purpose, action, and resource—not network
  location alone.
- Record delegation and impersonation explicitly.
- Reject unsigned or untrusted evaluator/strategy definitions where signing is
  required by the deployment profile.
- Treat evaluator output, reason text, external evidence, and AI content as
  untrusted data.
- Enforce schema, size, depth, character, URI, and reference-kind limits.

**Security validation boundary.** FEATURE-0013 owns only
structural schema validation, typed classification, declared bounds, projection
restrictions, and deterministic conformance. FEATURE-0013 does **not** claim
comprehensive semantic detection of secrets, credentials, personal data,
malware, or policy violations in arbitrary content. Producers must redact
prohibited content before submission. Comprehensive semantic content scanning
and runtime content inspection belong to later approved security features.
Kiro and Cursor must not invent secret-pattern, credential-pattern, PII,
malware, DLP, or content-scanning engines in FEATURE-0013. Where this document
describes rejecting a "secret", "credential", "prompt", or "unrestricted tenant
content" projection (for example the AC-13.7 fixtures), it means structural and
typed conformance against declared bounds, prohibited content-category rules,
and sensitivity ceilings — not a semantic content-scanning capability.

### 12.2 Confidentiality and minimization

The canonical governance record may retain protected evidence references but
must not indiscriminately duplicate personal data, secrets, credentials,
provider tokens, raw policy documents, model prompts, or unrestricted tenant
content. Records prefer stable references, digests, classifications, and
purpose-bound projections.

Different authorized projections serve customer, operator, auditor, regulator,
provider, and AI consumers. Redaction is field- and reason-aware; it must not
change the final outcome or falsely imply that omitted evidence did not exist.

### 12.3 Integrity and non-repudiation

The contract carries algorithm-agile canonicalization, digest, signature, and
trust metadata so later approved deployment profiles can require cryptographic
assurance without changing the common envelope. FEATURE-0013 does not compute a
digest over content, canonicalize content, sign, verify, timestamp, notarize,
enforce WORM retention, or claim that an opaque carrier proves content existed
or was considered.

Production deployment profiles may later require local or external signing,
trusted timestamps, append-only/WORM retention, Merkle batch proofs, or
notarization. Those capabilities, their algorithms, covered fields, products,
outage behavior, rotation, revocation, and latency policy require separately
approved architecture and implementation.

**Structural trust versus cryptographic verification boundary.**
Algorithm-agile carrier fields and structural trust metadata are `CONTRACT_NOW`:
the canonicalization-method identifier, the digest-algorithm identifier, the
covered-field descriptor, the signature-algorithm identifier, and the
structural trust-state fields are versioned data the contract carries without
selecting any value. The canonicalization algorithm/profile, the exact set of
digest-covered fields, the signature algorithm, and cryptographic
product/service selection are **DEFERRED** and remain blocked by **ADR-F13-002**.
Actual cryptographic verification, signing, HSM, notary, trusted-time, key,
revocation-distribution, and WORM services are not FEATURE-0013 artifacts.
Structural fixtures may use opaque deterministic test assertions to validate
state and failure semantics but must not claim cryptographic validity or select
an algorithm. When a profile requires verified trust, an unknown, absent,
expired, revoked, mismatched, or unverified trust state fails closed. RFC 8785
is illustrative only and cannot be selected by requirements, design, tasks,
fixtures, or code without a later approved decision.

### 12.4 Provider-neutral sensitivity classification vocabulary

Sensitivity classification uses a single closed, ordered, provider-neutral core
vocabulary:

```text
PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED
```

This sensitivity value is a decision/profile handling constraint; it does not
replace or redefine FEATURE-0012 `DataClassification`. FEATURE-0012
`DataClassification` continues to classify externally exchanged object data and
metadata using its seven canonical values. FEATURE-0013 sensitivity controls
the maximum handling and projection exposure permitted inside a
`DecisionProfile`, `EvaluationResult`, rationale, evidence, or typed result.

The two dimensions coexist and must not be silently converted. The canonical
minimum mapping is:

| FEATURE-0012 `DataClassification` | FEATURE-0013 sensitivity floor |
|---|---|
| `Public` | `PUBLIC` |
| `Customer-visible` | `CONFIDENTIAL` |
| `Tenant-confidential` | `CONFIDENTIAL` |
| `Operator-confidential` | `CONFIDENTIAL` |
| `Internal` | `INTERNAL` |
| `Sensitive` | `RESTRICTED` |
| `Secret-reference-only` | `RESTRICTED` |

A profile may raise this floor but may never lower it. The original FEATURE-
0012 classification remains present; the mapped sensitivity does not erase or
rewrite it. When multiple classifications apply, validation uses the most
restrictive result. An unmapped, unknown, contradictory, or lower-than-floor
value fails closed. Jurisdiction-specific mappings may further restrict
handling but cannot weaken either canonical dimension.

- The four core values are ordered from least to most sensitive. No fifth core
  value and no fewer than these four are permitted.
- Jurisdiction-specific classification labels and caveats (for example national
  government classification schemes or sector-specific handling markings) are
  versioned profile mappings onto this core vocabulary, not additional core
  enum values. Preserving jurisdiction-specific mappings as profile data keeps
  the contract provider-neutral and avoids enum proliferation.
- A profile that permits captured evaluator output (or any other captured
  classified content) must declare all of the following mandatory,
  profile-bound rules:
  - a sensitivity **ceiling** (the maximum core value the profile permits);
  - allowed and prohibited **content-category** rules;
  - a maximum **field** count;
  - a maximum **byte** size;
  - a maximum nesting **depth**;
  - the allowed value **types**.
- Missing or unknown mandatory profile values invalidate the profile.
- A captured value classified above the declared ceiling fails validation.

This vocabulary and its profile-bound rules are structural and typed conformance
only; they are not a semantic content-scanning capability (see section 12.1).

## 13. AI-first architecture

AI-first means AI can safely understand and assist with the system; it does not
mean AI controls the system.

Every decision provides structured reason codes, policy/evidence references,
alternatives, obligations, risk classification, corrective suggestions, and
confidence/indeterminacy where meaningful. Human-readable explanations are
derived from structured facts and identify their language and generator
version.

AI consumers receive only an authorized `DecisionContext` projection. They may:

- explain a decision;
- summarize evidence and obligations;
- identify missing information or anomalies;
- suggest a new request, policy change, or remediation for normal approval.

They may not mutate records, invent evidence, suppress audit, widen disclosure,
bypass approval, or execute a suggestion. Any AI-generated recommendation is
labelled with model/provider/version, prompt-template version, data boundary,
confidence where valid, and human-review status. Core decisions remain
available when AI is disabled or disconnected.

## 14. Reuse and standards strategy

This is the single mandatory FEATURE-0011 feature-level reuse summary for the
FEATURE-0013 architecture. Field meanings and controlled vocabularies are owned
only by `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 API/resource grammar | Extend | Reuse every shared profile, metadata, reference, boundary, validation, Problem Details, schema, and compatibility primitive; add only decision/audit domain semantics and approved per-kind constraints. | Approved | ADH-2026-017; ADH-2026-012; RFC-0022 |
| DecisionRecord, DecisionProfile, EvaluationResult, and AuditEvent domain contract | Build | No mature external model supplies Sovrunn's complete sovereign, provider-neutral authority, projection, scope, audit-obligation, and downstream-adoption semantics. The common contract is Sovrunn differentiation; external concepts remain inputs rather than native core types. | Approved | ADH-2026-017; DEC-0026; RFC-0023 |
| Bounded decision-requirements and composition concepts | Reuse | Reuse applicable OMG DMN decision-requirements concepts without selecting, embedding, or recreating a DMN engine or workflow language. | Approved | ADH-2026-017; RFC-0023 |
| Audit transport and operational correlation | Reuse | CloudEvents may be an optional transport mapping and OpenTelemetry may carry operational correlation; neither becomes the canonical record or audit authority. | Approved | ADH-2026-017; RFC-0023 |
| Supply-chain evidence references | Reuse | Reuse DSSE/in-toto and SPDX/CycloneDX concepts by typed reference without duplicating their schemas or implementing signing. | Approved | ADH-2026-017; ADR-F13-002 |
| Canonicalization, digest, signing, and cryptographic verification | Reuse | Preserve algorithm-agile carriers and evaluate mature standards later; RFC 8785 is illustrative only and no algorithm, covered-field set, product, or service is selected in FEATURE-0013. | Deferred | ADH-2026-017; ADR-F13-002 |
| Policy/evaluator and workflow runtimes | Wrap | CEL, OPA, Cedar, durable-execution, and workflow products are later candidates behind owning-feature adapters or ports; FEATURE-0013 defines no wrapper or runtime implementation. | Deferred | ADH-2026-017; DEC-0036 |

Capability-level assessments required by the canonical reuse standard must be
completed for the Sovrunn-owned domain contract and any later runtime candidate
before source implementation. They must reuse this summary and must not invent
new dispositions in requirements or design.

## 15. Compatibility and evolution

- Every canonical schema has explicit version and compatibility rules.
- Additive optional fields require semantic defaults and projection safety.
- Decision form, authority, adjudication, typed-result, reason-code, obligation,
  strategy, and digest semantics cannot change silently.
- Readers reject unsupported major versions and preserve unknown safe extension
  fields only under an approved extension mechanism.
- Registries for decision profile, form, authority, evaluator type, strategy,
  reason code, obligation, event type, evidence type, and projection are
  versioned and governable offline.

### 15.1 Public Problem Details binding

Public validation failures inherit the FEATURE-0012 RFC 9457 Problem Details
contract, stable codes, RFC 6901 JSON Pointer field paths, safe-denial, bounded
violations, and redaction rules. A FEATURE-0013 structural or semantic
validation failure has exactly this public binding:

| Problem member | Required public value or rule |
|---|---|
| `type` | `urn:sovrunn:problem:validation-failed` |
| `status` | `422` |
| `code` | `VALIDATION_FAILED` |
| `violations[].field` | RFC 6901 JSON Pointer to the rejected field or nearest safe parent |
| `violations[].code` | Exactly one code from the closed registry in section 15.2 |
| `violations[].message` | Human-readable, non-authoritative, redactable, and never parsed by clients |

The `DECISION_PROFILE_`, `DECISION_EVALUATION_`,
`DECISION_COMPOSITION_`, `DECISION_OBLIGATION_`, `DECISION_TRUST_`,
`DECISION_SCOPE_`, and `DECISION_RELATIONSHIP_` prefixes are public namespaces
only for `violations[].code`. They are not `Problem.type` URI families,
top-level `Problem.code` values, HTTP-status selectors, or internal-only
categories. Existing FEATURE-0012 top-level codes and HTTP meanings continue to
handle malformed input, authentication, authorization, absence, conflict,
stale state, internal failure, and temporary dependency unavailability.

### 15.2 Initial closed public violation-code registry

| Public family | Allowed `violations[].code` values | Architectural meaning |
|---|---|---|
| Decision profile | `DECISION_PROFILE_UNKNOWN`, `DECISION_PROFILE_VERSION_UNSUPPORTED`, `DECISION_PROFILE_INACTIVE`, `DECISION_PROFILE_SCHEMA_INVALID`, `DECISION_PROFILE_LIMIT_INVALID` | Profile identity, lifecycle, schema, or declared bounds are invalid |
| Evaluation | `DECISION_EVALUATION_TYPE_UNSUPPORTED`, `DECISION_EVALUATION_RESULT_INVALID`, `DECISION_EVALUATION_SCOPE_MISMATCH`, `DECISION_EVALUATION_LIMIT_EXCEEDED` | Captured evaluation cannot conform to the selected profile and record |
| Composition | `DECISION_COMPOSITION_STRATEGY_UNSUPPORTED`, `DECISION_COMPOSITION_GRAPH_INVALID`, `DECISION_COMPOSITION_CYCLE`, `DECISION_COMPOSITION_LIMIT_EXCEEDED`, `DECISION_COMPOSITION_INPUT_CONFLICT` | Strategy or bounded deterministic graph contract is invalid |
| Obligation | `DECISION_OBLIGATION_UNKNOWN`, `DECISION_OBLIGATION_INVALID`, `DECISION_OBLIGATION_UNSUPPORTED_MANDATORY` | Obligation vocabulary, payload, or mandatory support contract is invalid |
| Trust | `DECISION_TRUST_REQUIRED`, `DECISION_TRUST_UNKNOWN`, `DECISION_TRUST_EXPIRED`, `DECISION_TRUST_REVOKED`, `DECISION_TRUST_MISMATCH` | Required structural trust carrier or fail-closed trust state is invalid; no cryptographic-validity claim is implied |
| Scope | `DECISION_SCOPE_REQUIRED`, `DECISION_SCOPE_INVALID`, `DECISION_SCOPE_CONFLICT`, `DECISION_SCOPE_MISMATCH`, `DECISION_SCOPE_WIDENING` | Canonical scope is missing, invalid, duplicated, mismatched, or impermissibly widened |
| Relationship | `DECISION_RELATIONSHIP_KIND_INVALID`, `DECISION_RELATIONSHIP_TARGET_INVALID`, `DECISION_RELATIONSHIP_CYCLE`, `DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED`, `DECISION_RELATIONSHIP_CONFLICT` | Correction, supersession, revocation, or other governed relationship is invalid |

These names are public compatibility contracts owned by architecture. A new
code, rename, removal, changed meaning, different prefix, or remapping to a
different top-level problem requires an approved architecture and compatibility
change. Requirements may enumerate testable cases and map each case to an
existing code and JSON Pointer; they must not create, rename, alias, or assign
semantics to public codes. Design may map the approved codes to validators and
tasks may implement those mappings, but neither may introduce a parallel error
envelope, alternate `Problem.type`, new HTTP meaning, or ungoverned code.

### 15.3 Other compatibility rules

- Migrations are reversible where practical and preserve old evidence/digests.
- `AuditEvent` expansion requires fixtures for all six canonical FEATURE-0012
  governance `ScopeKind` values (`Platform`, `Organization`,
  `OrganizationUnit`, `Tenant`, `Project`, `Provider`), continued regression
  for the existing Organization-scoped alpha fixture, and a ServiceInstance
  subject fixture whose `metadata.scopeRef` is the containing Project.

Every inherited FEATURE-0011 or FEATURE-0012 vocabulary, shared type, error
contract, schema convention, or validation rule must have one dependency-
reconciliation record before design approval. That record identifies exact
versions and files; compares enums, requiredness, cardinality, ownership,
serialization, validation, security semantics, bindings, fixtures, and
migration impact; and classifies every difference as added, removed, renamed,
mapped, restricted, expanded, or semantically reinterpreted. Kiro and design
must escalate any unapproved delta and may not invent aliases, mappings,
restrictions, applicability exclusions, or compatibility behavior.

## 16. Observability and operability

Metrics and traces expose latency, queueing, evaluation fan-out, cache result,
composition strategy, failure class, retry, saturation, and audit-reconciliation
health without exposing decision evidence or personal data as labels.

Logs are structured, redacted, locally routable, retention-controlled, and
correlated by opaque identifiers. Health must distinguish evaluator, evidence,
store, audit, identity, key, time, and synchronization dependencies. Operators
must be able to diagnose and recover the system without vendor remote access.

## 17. Conformance model

The architecture requires executable positive and negative fixtures for the
following scenarios. These scenarios map exactly and
one-to-one to the stable conformance IDs `F13-CF-01` through `F13-CF-28` in
their existing order; the IDs must not be merged, renumbered, or omitted.
Additional `AC-13.7`, scope, security, and compatibility cases use separately
named IDs and do not alter this canonical count. One fixture may cover multiple
scenario IDs, but every scenario must have an explicit coverage-matrix entry;
coverage is counted by scenario ID, not by fixture-file count.

The 28 IDs are stable scenario families. Every independently mandatory clause
inside a family has a stable lowercase suffix in textual order (`.a`, `.b`,
`.c`, and so on). For example, F13-CF-08 has `.a` timeout, `.b` missing
evidence, `.c` evaluator failure, and `.d` policy conflict. The same rule is
mandatory for compound families F13-CF-04, 08, 11, 12, 16, 20, 26, and 28;
coverage of the parent requires coverage of every suffix. Requirements may
state expected results and error codes but must not merge, omit, renumber, or
invent the architectural cases.

The following additional stable cases are also mandatory and do not change the
canonical family count:

| Case range | Mandatory atomic cases |
|---|---|
| `F13-SCOPE-01` | Positive canonical Platform scope; serialized `metadata.scopeRef` is absent |
| `F13-SCOPE-02` | Positive canonical Organization scope |
| `F13-SCOPE-03` | Positive canonical OrganizationUnit scope |
| `F13-SCOPE-04` | Positive canonical Tenant scope |
| `F13-SCOPE-05` | Positive canonical Project scope |
| `F13-SCOPE-06` | Positive canonical Provider scope |
| `F13-SCOPE-07` | Absent scope rejected when Platform is not permitted |
| `F13-SCOPE-08` | Explicit Platform input normalizes to canonical absent serialized form |
| `F13-SCOPE-09` | Parallel, duplicate, or second scope source rejected |
| `F13-SCOPE-10` | Out-of-vocabulary scope, including ServiceInstance, rejected |
| `F13-SCOPE-11` | ServiceInstance subject accepted through typed `subjectRef` under its Project `metadata.scopeRef`; subject/scope mismatch rejected |
| `F13-EVAL-01` | Profile sensitivity and projection limits accepted at their declared bounds |
| `F13-EVAL-02` | Evaluator registration narrows profile limits |
| `F13-EVAL-03` | Evaluator widening of profile limits rejected |
| `F13-EVAL-04` | Missing, unknown, or incomparable mandatory limits rejected |
| `F13-EVAL-05` | Prohibited typed content category or value type rejected structurally |
| `F13-EVAL-06` | Evaluation scope mismatch or unauthorized widening rejected without existence disclosure |
| `F13-EVAL-07` | Opaque expected-versus-observed integrity state mismatch fails closed |
| `F13-EVAL-08` | Arbitrary text is not subjected to lexical secret, credential, PII, malware, or DLP scanning |
| `F13-SEC-01` | Sensitivity ceiling boundary |
| `F13-SEC-02` | Field-count boundary |
| `F13-SEC-03` | Byte-size boundary |
| `F13-SEC-04` | Nesting-depth boundary |
| `F13-SEC-05` | Allowed-value-type boundary |
| `F13-SEC-06` | Prohibited-category boundary |
| `F13-SEC-07` | FEATURE-0012 classification-to-sensitivity floor cannot be lowered |
| `F13-SEC-08` | Unknown, unmapped, contradictory, or below-floor classification fails closed |
| `F13-TRUST-01` | Algorithm-agile carrier presence and identifier syntax validated structurally |
| `F13-TRUST-02` | Covered-field descriptor presence validated without interpreting or selecting it |
| `F13-TRUST-03` | Unknown, absent, expired, revoked, mismatched, or unverified required trust state fails closed |
| `F13-TRUST-04` | No canonicalization, content-to-digest computation, signing, verification, notarization, or cryptographic-validity claim executes before ADR-F13-002 |
| `F13-COMPAT-01` | FEATURE-0012 Platform ScopeKind regression |
| `F13-COMPAT-02` | FEATURE-0012 Organization ScopeKind regression |
| `F13-COMPAT-03` | FEATURE-0012 OrganizationUnit ScopeKind regression |
| `F13-COMPAT-04` | FEATURE-0012 Tenant ScopeKind regression |
| `F13-COMPAT-05` | FEATURE-0012 Project ScopeKind regression |
| `F13-COMPAT-06` | FEATURE-0012 Provider ScopeKind regression |
| `F13-COMPAT-07` | Existing Organization-scoped FEATURE-0012 AuditEvent remains valid |
| `F13-COMPAT-08` | AuditEvent accepts all six existing governance scopes without changing the shared ScopeKind vocabulary |
| `F13-COMPAT-09` | FEATURE-0012 metadata, typed-reference, validation-layer, Problem Details, extension, and ServiceInstance-as-Project-scoped-resource semantics remain unchanged |

Every listed parent, suffix, and additional ID requires one coverage-matrix row.
A physical fixture may cover several rows, but file count never substitutes for
semantic coverage. All fixtures are deterministic, pure or bounded in-memory
contract tests; runtime-oriented scenarios validate only inputs, outputs, and
failure semantics and do not authorize persistence, networking, queues,
workers, cryptographic execution, AI runtime, or external integration.

| Conformance ID | Scenario |
|---|---|
| F13-CF-01 | simplest synchronous atomic authorization decision |
| F13-CF-02 | simplest deterministic denial with actionable rationale and obligations |
| F13-CF-03 | pure selection with a typed result and no adjudication facet |
| F13-CF-04 | ranking, classification, resolution, allocation, plan, assessment, recommendation, advisory, and simulation profiles |
| F13-CF-05 | human-approval asynchronous handoff returns an existing `Operation` value/reference; the eventual `DecisionRecord` links to that `Operation`, with no FEATURE-0013-specific pending-response envelope |
| F13-CF-06 | maximally complex but bounded hierarchical composite decision |
| F13-CF-07 | parallel evaluations with deterministic aggregation |
| F13-CF-08 | timeout, missing evidence, evaluator failure, and policy conflict |
| F13-CF-09 | idempotent replay and concurrent duplicate requests |
| F13-CF-10 | partial decision/audit persistence failure and reconciliation |
| F13-CF-11 | supersession, correction, revocation, cycle, and chain-limit behavior |
| F13-CF-12 | unauthorized profile, authority, projection, obligation bypass, and cross-scope reference rejection |
| F13-CF-13 | offline/air-gapped static profile-bundle structure, opaque trust carriers, and fail-closed trust states; no signing or verification execution |
| F13-CF-14 | export/import across two conforming provider implementations |
| F13-CF-15 | AI-disabled operation and AI projection redaction |
| F13-CF-16 | graph size, depth, width, fan-out, timeout, and payload budget rejection |
| F13-CF-17 | registration of a new decision family without a common-envelope change |
| F13-CF-18 | immutable supersession and revocation with calculated effective state |
| F13-CF-19 | local atomic decision/audit-obligation acceptance while downstream audit delivery is unavailable |
| F13-CF-20 | idempotency retry identity versus semantic decision identity across policy, profile, strategy, authority, evaluator-version, and validity changes |
| F13-CF-21 | deterministic replay from captured nondeterministic or AI evaluation output without rerunning the evaluator |
| F13-CF-22 | synchronous contract behavior with registry and policy resolver unavailable but valid local version-pinned bundles and opaque trust state present; no network or cryptographic execution |
| F13-CF-23 | structural rejection of unknown, expired, revoked, or opaque expected-versus-observed mismatched profile bundles without digest computation |
| F13-CF-24 | enforcement denial for unknown or unsupported mandatory obligations |
| F13-CF-25 | local authorization artifact cache hit, expiry, revocation, scope mismatch, and stale-policy failure |
| F13-CF-26 | cancellation and total remote-call, byte, concurrency, retry, jurisdiction, time, and cost budgets for composite graphs |
| F13-CF-27 | structural representation and fail-closed state semantics for a later local-digest/batched-signing/notarization profile during remote-service outage; no cryptographic execution |
| F13-CF-28 | bounded in-memory disconnected export/import contract semantics for replay, duplicate, ordering, opaque unknown trust root, revoked origin, and explicit conflict reconciliation; no synchronization runtime |

## 18. Matrix E v2 — architecture risk register

### 18.1 Compatibility with FEATURE-0012

Matrix E v2 reuses FEATURE-0012's risk-governance principles and evidence
lifecycle, not its identifiers. FEATURE-0012 identifiers `F12-R01` through
`F12-R16` remain owned exclusively by FEATURE-0012. FEATURE-0013 defines the
separate `F13-R01` through `F13-R31` namespace. The two namespaces must not be
merged, renumbered, aliased, inherited, or mapped by ordinal position.
FEATURE-0013 carries its own primary controls, detection/response evidence,
automated evidence staging, human-only residual acceptance, owner, corrective
path, reassessment trigger, category, affected architecture decisions,
inherent and residual assessment, treatment, control owner, later-feature
owner, acceptance expiry, and per-risk status.

Kiro and implementation automation must preserve `F13-Rxx` identifiers. They
must not renumber, merge, delete, downgrade, or accept risks. Automation may
stage evidence; only an authorized human may record residual-risk disposition.

### 18.2 Assessment method

Likelihood and impact use an ordinal 1–5 scale:

| Score | Likelihood | Impact |
|---:|---|---|
| 1 | Rare | Negligible |
| 2 | Unlikely | Minor |
| 3 | Possible | Material |
| 4 | Likely | Major |
| 5 | Almost certain | Critical or platform-wide |

Impact is the highest credible effect across correctness, security, privacy,
sovereignty, availability, compliance, compatibility, operations, financial
cost, or reputation. `score = likelihood × impact`: 1–4 Low, 5–9 Medium,
10–16 High, and 17–25 Critical. Numeric scoring supports consistency; written
rationale remains authoritative.

Residual acceptance rules:

- Low: feature reviewer may accept with evidence.
- Medium: named residual owner, corrective path, and trigger are mandatory.
- High: architecture owner plus relevant security/privacy/operations owner;
  acceptance is time-bounded.
- Critical: implementation approval is prohibited until reduced or governed by
  an explicit exceptional decision.
- Critical correctness or security risk cannot be accepted.
- Unowned, unverified, expired, or ownerless-deferred risk remains open.

Treatments are `AVOID`, `MITIGATE`, `TRANSFER`, `ACCEPT`, or `DEFER`. `DEFER`
requires a later owner and trigger; it does not mean accepted.

### 18.3 Matrix E v2 summary

The residual levels below are architecture targets after the stated controls,
not human acceptance.

| ID | Category | Risk scenario | ADs | Inherent L×I | FEATURE-0013 controls | Verification/detection | Target residual | Treatment | Residual/later owner | Reassessment trigger | Final residual status |
|---|---|---|---|---:|---|---|---:|---|---|---|---|
| F13-R01 | Governance/compatibility | Later features create overlapping `DecisionProfile` semantics | AD-029, AD-030 | 4×4 High | Governed registration, ownership, reuse review, compatibility | Duplicate-profile and new-family fixtures | 2×3 Medium | MITIGATE | Architecture owner | New profile family or stable promotion | PENDING_HUMAN_REVIEW |
| F13-R02 | Correctness | Typed results become an ungoverned property bag or override the envelope | AD-002, AD-029 | 4×4 High | Versioned result schemas; semantic-override prohibition | Schema and negative conformance | 2×4 Medium | MITIGATE | FEATURE-0013 owner | New result extension mechanism | PENDING_HUMAN_REVIEW |
| F13-R03 | Security/correctness | Advisory, recommendation, simulation, or AI output is consumed as authority | AD-020, AD-031 | 4×5 Critical | Explicit authority; projection marking; enforcement validation | Authority-confusion negative tests and audit | 1×5 Medium | MITIGATE | Security architecture; enforcement owner | New AI/execution consumer | PENDING_HUMAN_REVIEW |
| F13-R04 | Correctness | Stale `EffectivePolicyContext` produces an invalid decision | AD-007, AD-038 | 4×5 Critical | Versioned identity, optional opaque integrity carrier, effective time, freshness, revocation, fail-closed rules | Cache expiry/revocation/staleness tests without digest computation | 2×5 High | MITIGATE | Effective-context owner | Production policy resolver or new inheritance rule | PENDING_HUMAN_REVIEW |
| F13-R05 | Correctness | Weak retry identity or cache identity reuses the wrong decision | AD-013, AD-036 | 4×5 Critical | Separate retry key and complete semantic identity descriptor; later optional digest profile | Policy/profile/strategy/time mutation tests | 1×5 Medium | MITIGATE | Decision persistence owner | Semantic-identity fields or canonical digest profile change | PENDING_HUMAN_REVIEW |
| F13-R06 | Audit/compliance | A decision is durable without its audit obligation | AD-014, AD-035 | 3×5 High | Same-boundary atomic acceptance and preallocated audit identity | Partial-failure and reconciliation fixtures | 1×5 Medium | MITIGATE | Audit persistence owner | Production persistence selection | PENDING_HUMAN_REVIEW |
| F13-R07 | Availability/latency | Remote audit processing blocks synchronous decisions | AD-008, AD-035 | 4×4 High | Local obligation acceptance; asynchronous downstream processing | Remote-audit outage and backlog tests | 1×4 Low | AVOID | Audit implementation owner | External audit/SIEM integration | PENDING_HUMAN_REVIEW |
| F13-R08 | Cost/availability | Composite graphs cause fan-out, latency, retry, or cost explosion | AD-006, AD-040 | 4×5 Critical | Bounds on graph, calls, bytes, authorities, concurrency, retries, time, jurisdiction, and cost | Limit, cancellation, load, and adversarial fixtures | 2×5 High | MITIGATE | Decision runtime owner | Limit increase or new remote evaluator | PENDING_HUMAN_REVIEW |
| F13-R09 | Correctness/AI | Nondeterministic evaluator output cannot be replayed or explained | AD-005, AD-037 | 4×4 High | Captured normalized output, provenance, input identity, and opaque integrity carrier | Replay without evaluator access or digest computation | 2×3 Medium | MITIGATE | Evaluator/AI owner | New probabilistic evaluator | PENDING_HUMAN_REVIEW |
| F13-R10 | Security/AI | AI receives unauthorized evidence or sensitive canonical data | AD-019, AD-020 | 4×5 Critical | Authorized `DecisionContext`, minimization, classification, projection policies | Cross-scope/redaction and prompt-leak fixtures | 1×5 Medium | MITIGATE | AI-context and security owners | New model, projection, or jurisdiction | PENDING_HUMAN_REVIEW |
| F13-R11 | Security | Enforcement silently ignores an unknown mandatory obligation | AD-032, AD-043 | 3×5 High | Unknown/unsupported obligation makes allow unenforceable | Obligation capability negative tests | 1×5 Medium | AVOID | Enforcement owner | New obligation vocabulary | PENDING_HUMAN_REVIEW |
| F13-R12 | Latency/availability | Central authorization or decision service becomes a universal bottleneck | AD-008, AD-039 | 4×5 Critical | Local verified artifacts, policy snapshots, caches, bounded freshness | Central-service outage and horizontal-load tests | 2×4 Medium | MITIGATE | Authorization implementation owner | Production enforcement integration | PENDING_HUMAN_REVIEW |
| F13-R13 | Security | Revoked or stale cached authorization remains enforceable | AD-039, AD-043 | 3×5 High | Audience/scope/expiry/revocation/freshness validation; fail closed | Cache expiry, revocation, scope-mismatch tests | 2×5 High | MITIGATE | Authorization/security owner | Revocation architecture selection | PENDING_HUMAN_REVIEW |
| F13-R14 | Security/availability | Fail-open behavior exposes protected resources during dependency failure | AD-005, AD-043 | 3×5 High | Fail-closed high-risk profiles; bounded approved `SecurityExceptionRef` structural evidence | Failure injection, exception-reference, scope, expiry, and owner-mismatch checks | 1×5 Medium | AVOID | Security architecture owner | Any fail-open request | PENDING_HUMAN_REVIEW |
| F13-R15 | Correctness | Correction, supersession, or revocation yields ambiguous effective state | AD-004, AD-034 | 3×4 High | Immutable linked effects; bounded cycle-free deterministic projection | Chain, cycle, fork, and projection fixtures | 1×4 Low | MITIGATE | Read-model owner | New relationship/effective-state semantics | PENDING_HUMAN_REVIEW |
| F13-R16 | Privacy/security | Decision, rationale, evidence, error, or projection leaks restricted data | AD-019, AD-024 | 4×5 Critical | References/digests, minimization, classification, redaction, purpose projections | Secret/PII/cross-scope negative tests | 1×5 Medium | MITIGATE | Security/privacy owner | New evidence or projection type | PENDING_HUMAN_REVIEW |
| F13-R17 | Compliance/integrity | Erasure or retention action breaks evidence integrity or legal hold | AD-004, AD-022 | 3×5 High | Structural tombstone and retention/legal-hold contracts now; cryptographic erasure later | Structural hold, erasure-state, opaque-carrier, and tombstone tests | 2×4 Medium | MITIGATE | Evidence/retention owner | Regulated deployment or retention change | PENDING_HUMAN_REVIEW |
| F13-R18 | Latency/availability | Remote HSM, time, signing, or notarization becomes a hot-path dependency | AD-022, AD-041 | 4×4 High | No cryptographic service in FEATURE-0013; later tiered local/batch/asynchronous assurance | No-network/no-crypto gate now; outage and batch-proof tests in owning later feature | 1×4 Low | AVOID | Security infrastructure owner | High-assurance profile adoption | PENDING_HUMAN_REVIEW |
| F13-R19 | Sovereignty/security | Disconnected import accepts forged, replayed, or untrusted records | AD-018, AD-042 | 3×5 High | Static origin-aware manifests, algorithm-agile signature carriers, replay/duplicate markers, ordering, opaque trust states, and fail-closed quarantine semantics now; cryptographic signing and verification later under ADR-F13-002 | Structural replay/duplicate/ordering/unknown-or-revoked-trust-state fixtures now; forgery and signature-verification evidence in the later owning feature | 1×5 Medium | MITIGATE | Synchronization/security owner | First disconnected deployment | PENDING_HUMAN_REVIEW |
| F13-R20 | Correctness | Disconnected origins create conflicting authoritative decisions | AD-018, AD-042 | 3×5 High | Immutable origin identity; no last-writer-wins; governed reconciliation decision | Concurrent-origin conflict fixtures | 2×4 Medium | MITIGATE | Multi-site owner | Active/active or disconnected write enablement | PENDING_HUMAN_REVIEW |
| F13-R21 | Portability | Provider/runtime-native identifiers or schemas leak into the core | AD-011, AD-023 | 4×4 High | Adapter-only native data; normalized capabilities; typed references | Import/schema lint and portability fixtures | 1×4 Low | AVOID | Adapter architecture owner | First provider/runtime adapter | PENDING_HUMAN_REVIEW |
| F13-R22 | Scope/governance | Production registry, workflow, persistence, provider, or AI runtime leaks into Phase 2 | AD-011, AD-028 | 4×4 High | Contract-only scope, non-goals, feature gates | Changed-file, dependency, and no-side-effect gates | 1×4 Low | AVOID | FEATURE-0013 owner | Any runtime dependency proposal | PENDING_HUMAN_REVIEW |
| F13-R23 | Cost | Durable decisions, audit, evidence, or projections create uncontrolled storage growth | AD-024, AD-033 | 4×4 High | Risk-based recording, bounded payloads, hot summaries, referenced cold evidence, retention profiles | Volume/limit/retention tests and metrics | 2×4 Medium | MITIGATE | Audit/evidence owner | Retention increase or production volume | PENDING_HUMAN_REVIEW |
| F13-R24 | Supply chain/security | Profile, schema, strategy, evaluator, or trust bundle is compromised | AD-025, AD-038 | 3×5 High | Version-pinned bundles, opaque trust-root identity, structural revocation/quarantine/provenance now; cryptographic verification later | Structural tamper-state/revocation fixtures now; signature verification in owning later feature | 1×5 Medium | MITIGATE | Supply-chain security owner | Signing-root or distribution change | PENDING_HUMAN_REVIEW |
| F13-R25 | Correctness/security | Clock uncertainty invalidates ordering, expiry, freshness, or replay protection | AD-018, AD-022 | 3×5 High | Trusted-time provenance, validity epoch, profile outage behavior | Skew, rollback, expiry-boundary tests | 2×4 Medium | MITIGATE | Runtime/security owner | Disconnected or multi-site deployment | PENDING_HUMAN_REVIEW |
| F13-R26 | Security/privacy | Cross-scope references disclose or influence another tenant or governance domain | AD-015, AD-019 | 3×5 High | Scope authorization, no-existence disclosure, typed references, projections | Cross-scope and safe-denial fixtures | 1×5 Medium | MITIGATE | API/security owner | New scope/reference kind | PENDING_HUMAN_REVIEW |
| F13-R27 | Compatibility/governance | Later features silently redefine common decision or public error semantics | AD-029, AD-030, AD-045 | 4×5 Critical | Stable envelope, profile conformance, closed public error registry, architecture change control | Baseline/schema/error semantic-diff gate | 1×5 Medium | AVOID | Architecture owner | Common-envelope, closed-vocabulary, or public-error change | PENDING_HUMAN_REVIEW |
| F13-R28 | Compatibility | Expanded AuditEvent allowed scopes or ServiceInstance subject handling breaks FEATURE-0012 consumers | AD-015 | 3×4 High | Retain the six-value ScopeKind vocabulary, use typed ServiceInstance subject references, preserve Organization fixtures, and require schema-diff review | Six-scope AuditEvent compatibility, Organization regression, and Project-scope/ServiceInstance-subject tests | 1×4 Low | MITIGATE | FEATURE-0013/API owner | Stable API promotion | PENDING_HUMAN_REVIEW |
| F13-R29 | Traceability | `DecisionRecord` rename breaks spine, documents, schemas, or consumer identity | AD-001 | 3×4 High | Architecture clarification, coordinated traceability update, aliases/migration if required | Repository-wide identity and baseline checks | 1×4 Low | MITIGATE | Architecture owner | Rename approval or stable promotion | PENDING_HUMAN_REVIEW |
| F13-R30 | Cost/sovereignty | High-assurance controls make small or local deployments unaffordable | AD-016, AD-017, AD-041 | 3×4 High | Tiered conformance and assurance profiles; no universal remote dependencies | Minimal/regulated/sovereign profile cost fixtures | 2×3 Medium | MITIGATE | Deployment architecture owner | New mandatory assurance control | PENDING_HUMAN_REVIEW |
| F13-R31 | Governance/cost | Downstream adoption governance becomes duplicative bureaucracy or a premature registry framework | AD-044 | 4×3 High | One normative contract, one short per-feature section, one lightweight gate; trigger-based escalation only | Changed-file/gate review; absence of manifests, registry, index, and separate approval workflow | 1×3 Low | AVOID | Architecture governance owner | Repeated semantic drift or structured-manifest trigger in section 28.8 | PENDING_HUMAN_REVIEW |

### 18.4 Architecture-stage treatment approval

- Reviewer: Sanjeev Kumar
- Decision date: 2026-07-27
- Decision: `APPROVE_TREATMENT` for F13-R01 through F13-R31
- Scope: approve the risks as requirements inputs together with their proposed
  controls, verification plans, target residual levels, ownership roles, and
  reassessment triggers
- Explicitly open High target-residual risks: F13-R04, F13-R08, F13-R13
- Final residual-risk decision: `PENDING_HUMAN_REVIEW`

This approval authorizes Kiro to preserve the risk controls as requirements. It
does not assert that controls have been implemented, does not accept residual
risk, and does not authorize a production capability. Final residual decisions
require implementation evidence and the per-risk record below.

### 18.5 Final per-risk residual acceptance record

At final feature review, the human reviewer must complete one record per risk;
architecture treatment approval and a global feature approval do not accept
unreviewed residual risks.

```yaml
risk_id: F13-Rxx
decision: ACCEPT | REJECT | DEFER
residual_likelihood: 1..5
residual_impact: 1..5
residual_level: LOW | MEDIUM | HIGH | CRITICAL
residual_rationale: <human rationale>
owner: <accountable human role or identity>
later_feature: <feature or null>
corrective_path: <action if control fails or risk materializes>
reassessment_trigger: <event forcing re-review>
acceptance_expiry: <date, milestone, or trigger>
reviewer: <human identity>
date: <YYYY-MM-DD>
evidence_reviewed: []
```

Automation must generate a worksheet populated with risk definitions, controls,
verification evidence, and blank human-only fields. Automated PASS, architecture
approval, reuse approval, or prior-feature acceptance is not residual-risk
acceptance.

### 18.6 Risk traceability lifecycle

```text
architecture.md Matrix E source
  -> Kiro requirements preserve IDs and mandate controls
  -> design maps controls to components and verification
  -> tasks implement controls and evidence collection
  -> automation stages immutable evidence pointers
  -> human records per-risk residual disposition
  -> final feature approval summarizes, but does not replace, those records
```

## 19. Decision scorecard

Scores are 1 (poor) to 5 (excellent). Security and correctness are veto
criteria: a decision scoring below 4 in either cannot be approved without a
recorded exception. “Cost” means lifecycle cost, not merely initial build cost.

The implementation class is an architecture-owned disposition, not a choice
for requirements or design. Each decision has exactly one class. A
`CONTRACT_NOW` row authorizes only the schemas, carriers, validation, fixtures,
documentation, or bounded in-memory reference behavior stated in sections 3,
17, and 27; it never authorizes a production runtime. Later runtime invariants
and deferred selections referenced in a recommendation remain governed by the
explicit section boundary or ADR and do not change the row's class.

| ID | Architectural decision | Implementation class | Long. | Correct. | Cost | Security | AI-first | Reuse | Recommendation |
|---|---|---|---:|---:|---:|---:|---:|---:|---|
| AD-001 | Rename common `DecisionObject` to `DecisionRecord` | CONTRACT_NOW | 5 | 5 | 5 | 4 | 5 | 5 | Approve with spine/traceability update |
| AD-002 | Keep calling-domain requests, evaluations, decisions, audit, operations, and projections distinct; define no shared FEATURE-0013 DecisionRequest type | CONTRACT_NOW | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-003 | Adjudication outcomes are optional, facet-specific, and closed | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve; keep errors and typed results separate |
| AD-004 | Immutable final records with linked corrections | CONTRACT_NOW | 5 | 5 | 4 | 5 | 5 | 5 | Approve structural linkage; erasure mechanisms require later architecture |
| AD-005 | Deterministic versioned composition strategies | CONTRACT_NOW | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-006 | Bounded DMN-like dependency DAG, not workflow language | CONTRACT_NOW | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-007 | Runtime hierarchy resolution produces one `EffectivePolicyContext` | INVARIANT_FOR_LATER | 5 | 5 | 5 | 5 | 5 | 5 | Approve as later invariant; FEATURE-0013 records only the typed reference and provenance contract |
| AD-008 | Synchronous fast path without mandatory queue | INVARIANT_FOR_LATER | 5 | 5 | 5 | 4 | 4 | 5 | Approve with strict budgets |
| AD-009 | `Operation` owns asynchronous lifecycle and is the sole non-final asynchronous handoff; no pending-decision envelope | CONTRACT_NOW | 5 | 5 | 4 | 5 | 5 | 5 | Approve contract use and conformance only |
| AD-010 | Orchestrate authority; choreograph reactions | INVARIANT_FOR_LATER | 5 | 5 | 4 | 5 | 5 | 5 | Approve as later runtime invariant |
| AD-011 | Runtime/vendor selection deferred | DEFERRED | 5 | 5 | 5 | 5 | 4 | 5 | Define later selection gate |
| AD-012 | Pure stateless reference decision kernel | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve only bounded in-memory conformance behavior |
| AD-013 | Effective-once via idempotency, not exactly-once claims | INVARIANT_FOR_LATER | 5 | 5 | 5 | 5 | 4 | 5 | Approve as later persistence/delivery invariant |
| AD-014 | Decision and audit durable-acceptance invariant | INVARIANT_FOR_LATER | 5 | 5 | 4 | 5 | 5 | 5 | Approve with later persistence pattern |
| AD-015 | Retain the six FEATURE-0012 governance ScopeKind values; expand AuditEvent to all six; represent ServiceInstance as a typed subject under Project scope | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve through the consolidated ADH-2026-017 |
| AD-016 | Sovereignty as seven enforceable profile dimensions | CONTRACT_NOW | 5 | 5 | 4 | 5 | 5 | 5 | Approve structural profile contract |
| AD-017 | Connected, private, edge, and air-gapped conformance profiles | CONTRACT_NOW | 5 | 5 | 3 | 5 | 5 | 5 | Approve static profiles; tier conformance to control cost |
| AD-018 | Static local-dependency and synchronization carrier contracts without signing or runtime synchronization | CONTRACT_NOW | 5 | 5 | 3 | 5 | 5 | 5 | Approve structural contract; cryptographic execution deferred |
| AD-019 | Authorized projections and data minimization | CONTRACT_NOW | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-020 | AI consumes `DecisionContext`, never canonical unrestricted data | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve projection constraints; no AI runtime |
| AD-021 | AI optional; deterministic contracts and conformance work without it | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve AI-disabled conformance |
| AD-022 | Algorithm-agile integrity/signature carriers and structural trust states | CONTRACT_NOW | 5 | 5 | 4 | 5 | 5 | 5 | Approve carriers only; all selections and execution deferred under ADR-F13-002 |
| AD-023 | Provider-neutral core with adapter-owned native identifiers | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-024 | Hot-summary and referenced-cold-evidence carrier contract | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve structural references; storage and retention runtime deferred |
| AD-025 | Versioned offline-capable static semantic registries | CONTRACT_NOW | 5 | 5 | 4 | 5 | 5 | 5 | Approve; no registry service |
| AD-026 | Telemetry is never audit authority | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve contract separation and conformance |
| AD-027 | Jurisdiction-specific controls supplied as profiles, not core enums | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-028 | Contract-only Phase 2 and its feature gates | CONTRACT_NOW | 5 | 5 | 5 | 5 | 4 | 5 | Approve; production runtime remains prohibited |
| AD-029 | Stable envelope plus versioned `DecisionProfile` and typed result | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve; primary extensibility mechanism |
| AD-030 | Closed primary-form vocabulary with registered secondary facets | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-031 | Authority is orthogonal: authoritative, advisory, recommendation, simulation | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-032 | Authorization is a profile; identity, evaluation, and enforcement remain outside | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve structural profile only |
| AD-033 | Risk-based durable-recording and aggregation profile semantics | CONTRACT_NOW | 5 | 5 | 5 | 5 | 4 | 5 | Approve contract and conformance; no persistence runtime |
| AD-034 | Immutable final records; supersession/revocation are linked effects and projected state | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-035 | Local atomic decision plus audit-obligation acceptance; downstream audit is asynchronous | INVARIANT_FOR_LATER | 5 | 5 | 5 | 5 | 5 | 5 | Approve as later persistence/delivery invariant; current linkage is governed by AD-014 contract fields |
| AD-036 | Separate retry idempotency from complete semantic decision identity | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve carrier fields and validation |
| AD-037 | Deterministic composition replays captured evaluator output without evaluator rerun | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve bounded pure conformance |
| AD-038 | Static local version-pinned profile, policy-context, schema, strategy, and opaque trust-state contract | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve structural conformance; cryptographic verification is later |
| AD-039 | Production authorization may use local verified artifacts/caches with bounded freshness and revocation | INVARIANT_FOR_LATER | 5 | 5 | 5 | 5 | 4 | 5 | Approve as later authorization/cache invariant; no cache runtime now |
| AD-040 | Profile and conformance bounds cover calls, bytes, concurrency, retries, jurisdictions, time, and cost | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve declared bounds and rejection fixtures |
| AD-041 | Cryptographic assurance is tiered; remote notary/HSM is not a universal hot-path dependency | INVARIANT_FOR_LATER | 5 | 5 | 5 | 5 | 5 | 5 | Preserve as later invariant; selections and execution remain ADR-F13-002 deferred |
| AD-042 | Static disconnected synchronization contract uses origin-aware manifests and explicit conflict evidence without signing | CONTRACT_NOW | 5 | 5 | 4 | 5 | 5 | 5 | Approve structural contract; synchronization runtime and signing are later |
| AD-043 | High-risk profiles fail closed; fail-open requires a bounded approved exception | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve profile validation and conformance |
| AD-044 | Downstream adoption uses one normative contract, one short per-feature section, and one lightweight gate | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve; structured manifests remain trigger-based |
| AD-045 | Bind FEATURE-0013 validation to inherited Problem Details and an architecture-owned closed public violation-code registry | CONTRACT_NOW | 5 | 5 | 5 | 5 | 5 | 5 | Approve; requirements map cases but cannot create public codes |

### 19.1 Aggregate assessment

| Criterion | Average | Assessment |
|---|---:|---|
| Longevity | 5.0 | Strong: semantic contracts outlive runtime products |
| Correctness | 5.0 | Strong if hierarchy, composition, and failure semantics remain explicit |
| Cost-effectiveness | 4.6 | Strong; local data-plane inputs, asynchronous downstream audit, tiered cryptography, typed profiles, and risk-based recording avoid latency and excess storage |
| Security | 4.9 | Strong after projection, trust, custody, and audit controls are included |
| AI-first | 4.9 | Strong and safe: structured typed context with explicit non-authoritative modes |
| Reusability | 5.0 | Strong across domains, providers, runtimes, and jurisdictions |

## 20. Required revisions before approval

The prior draft direction was sound but incomplete for a sovereign architecture.
This section records architecture rationale and provenance; it is not an
independent implementation backlog. The implementation class of every decision
is fixed by the section 19 scorecard and the explicit section 9–12 boundaries.
No item below authorizes a production runtime, product selection, or
cryptographic execution. This version makes the following necessary revisions:

1. Treat sovereignty as data, operational, technical, legal, cryptographic,
   supply-chain, and AI/model control—not residency alone.
2. Add connected, edge, disconnected, and air-gapped deployment profiles.
3. Add local trust, identity, keys, time, registry, audit, and observability
   dependency contracts, with signed synchronization retained only as a later
   invariant pending approved cryptographic architecture.
4. Add data minimization, legal hold, retention/disposition, lawful-redaction,
   payload-custody, and integrity-preserving tombstone contracts now; retain
   cryptographic-erasure execution as a later approved invariant.
5. Add workload identity, delegation, break-glass, purpose, and assurance.
6. Add software provenance, SBOM, offline verification, and crypto-agility.
7. Add AI/model sovereignty and make AI optional for core correctness.
8. Add explicit budgets, quotas, admission control, backpressure, and tiered
   conformance to make the design affordable for small local providers.
9. Add provider-to-provider export/import conformance to prevent sovereign
   lock-in.
10. Keep exact legal retention periods and provider/runtime products out of the
    core architecture; express them through governed deployment profiles.
11. Replace the universal allow/deny assumption with decision forms, optional
    adjudication facets, and typed results.
12. Add orthogonal authority so recommendations, simulations, and advisory
    records cannot be mistaken for enforceable decisions.
13. Add the governed, versioned `DecisionProfile` extension mechanism so later
    domains do not change the common envelope.
14. Include authorization as a profile while keeping authentication, policy
    evaluation, identity integration, and enforcement in their owning features.
15. Make immutable finality unambiguous: later records carry correction,
    supersession, or revocation effects; projections calculate effective state.
16. Make decision and audit-obligation acceptance local and atomic while
    allowing downstream audit processing to remain asynchronous.
17. Separate retry idempotency from complete semantic decision identity.
18. Define deterministic replay from captured evaluator results without
    requiring nondeterministic or AI evaluators to reproduce output.
19. Prohibit implicit remote control-plane dependencies on the synchronous path.
20. Add structural contracts and conformance for later local authorization,
    profile-bundle, policy-context, evidence-fan-out, cryptographic, and
    disconnected-synchronization safeguards without implementing their
    production runtimes.

## 21. Important deferred decisions

The following must be decided later, not smuggled into FEATURE-0013:

- production decision/audit persistence and transaction pattern;
- workflow/orchestration runtime and queue/broker;
- evidence store and WORM/notarization implementation;
- identity provider, key manager/HSM, policy engine, and secrets platform;
- India government security-classification vocabulary and exact retention
  schedules;
- concrete cloud-provider, infrastructure-platform, container-orchestrator,
  accelerator, or AI-service adapters;
- production AI model, inference runtime, vector store, or agent framework;
- exact canonicalization, signature, and trusted-time profiles;
- cross-site active/active topology, consensus, RPO, and RTO.

## 22. Architecture acceptance criteria

This architecture may become the controlling FEATURE-0013 input when the human
architecture owner records:

1. approval of `DecisionRecord` terminology and affected spine updates;
2. approval of the six-value FEATURE-0012 `ScopeKind` inheritance, AuditEvent
   acceptance of all six governance scopes, and typed ServiceInstance subject
   model through ADH-2026-017;
3. approval of the durable decision/audit obligation;
4. confirmation of runtime-neutral Phase 2 scope;
5. acceptance or revision of AD-001 through AD-045;
6. acceptance of the sovereign deployment profile model;
7. architecture-stage `APPROVE_TREATMENT`, `REJECT`, or `DEFER` decisions are
   recorded for Matrix E F13-R01 through F13-R31;
8. controls, target residual levels, ownership roles, corrective paths,
   reassessment triggers, reviewer, and date are recorded as requirements
   inputs without claiming implementation evidence;
9. F13-R04, F13-R08, and F13-R13 remain explicitly open with later owners,
   production-entry gates, and reassessment triggers;
10. final per-risk residual acceptance remains a mandatory human gate after
    implementation evidence and is not a prerequisite for requirements
    generation;
11. Kiro requirements, design, and tasks pass every anti-overengineering
    guardrail in section 27;
12. the lightweight downstream-adoption contract in section 28 is preserved
    without introducing a separate registry or approval workflow;
13. approval of the FEATURE-0012 `DataClassification` to FEATURE-0013
    sensitivity-floor mapping in section 12.4;
14. approval of the declarative registry, profile/evaluator ownership,
    evaluation-scope, error-family, and atomic conformance decisions now
    consolidated in the main body; and
15. every accepted FEATURE-0013 reuse-summary row has `Decision status` set to
    `Approved` with the consolidation approval as human evidence; deferred rows
    remain `Deferred`, and no `Proposed` row remains in the approved document;
16. confirmation that `DecisionRequest` remains a calling-domain-owned
    conceptual input with no FEATURE-0013 shared schema/type, and that the only
    asynchronous handoff is the inherited `Operation` value/reference; no
    pending-decision response envelope is authorized; and
17. a consolidation handoff and traceability updates that supersede the prior
    requirements authorization without altering historical approvals.

## 23. Reassessment triggers

Re-review is mandatory before stable API promotion, production runtime
selection, external regulated deployment, new jurisdiction, cross-border data
movement, autonomous AI capability, cryptographic algorithm/profile change,
provider adapter adoption, active/active multi-site operation, evidence signing,
or any relaxation of fail-closed, projection, audit, or immutability rules.

## 24. Normative local references

- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/architecture/api-resource-standard.md`
- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
- `docs/features/FEATURE-0011-reuse-assessment-standard.md`
- `docs/reviews/feature-gates/FEATURE-0011-approval-review.md`
- `.kiro/specs/api-resource-naming-status-and-validation-standard/requirements.md`
- `.kiro/specs/api-resource-naming-status-and-validation-standard/design.md`
- `docs/reviews/feature-gates/FEATURE-0012-approval-review.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-011-feature-0011-reuse-assessment-standard.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-012-feature-0012-api-resource-standard.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-013-operation-allowed-scopes.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-017-feature-0013-consolidated-architecture.md`
- `docs/rfc/RFC-0023-decision-and-audit-standard.md`

## 25. External review inputs

- MeitY government-cloud contractual guidance: data and derived results must
  follow empanelment residency terms, with preservation and e-discovery needs.
- MeitY CSP empanelment material: public cloud, VPC, and Government Community
  Cloud are distinct deployment models with audited controls.
- Sovereign deployment capability evidence is evaluated through provider-
  neutral conformance profiles; provider marketing claims and product names
  are not architecture inputs or canonical domain dependencies.
- RFC 8785, CloudEvents, OpenTelemetry, DMN, SPDX/CycloneDX, DSSE, and in-toto
  are reuse candidates subject to the FEATURE-0013 reuse assessment.

## 26. Validation against Sovrunn foundation principles

| Sovrunn principle or spine invariant | Result | FEATURE-0013 evidence |
|---|---|---|
| Reuse before build | PASS | Section 14 reuses API/resource grammar and evaluates established decision, telemetry, canonicalization, supply-chain, policy, and durable-execution concepts before custom implementation |
| Provider-neutral core | PASS | Sections 3, 8.4, and 11.3 exclude provider-native semantics and defer concrete adapters |
| Adapter before integration | PASS | Evaluator-native results remain behind adapters; identity, policy, provider, evidence, cryptographic, AI, and runtime products remain outside the core |
| Engine-neutral policy contracts | PASS | `EvaluationResult`, `EffectivePolicyContext`, profiles, and composition carry no policy-engine-native types |
| Customer boundary protection | PASS | Canonical governance records and purpose-specific projections prevent provider/internal detail leakage to customer APIs |
| Decision before execution | PASS | `DecisionRecord` is separate from action; `Operation` owns asynchronous lifecycle; the decision graph cannot provision or execute arbitrary work |
| Decision, audit, and operation remain distinct | PASS WITH APPROVED TERMINOLOGY CHANGE | Section 5 preserves the three meanings while proposing `DecisionRecord` as the clearer replacement for `DecisionObject` |
| Structured explainability | PASS | Typed rationale, reasons, alternatives, obligations, provenance, evidence references, and authorized `DecisionContext` projections are mandatory |
| No side effects in Phase 2 contracts | PASS | Sections 3, 7.4, 8.4, and 20 prohibit runtime/provider/plugin execution and defer production persistence and orchestration choices |
| AI is consumer, not authority | PASS | Sections 6.3 and 13 make AI non-authoritative, projected, policy-governed, optional, and unable to mutate or execute |
| Earlier features own shared concepts | PASS WITH CONTROLLED EXTENSION | FEATURE-0012 grammar and six-value ScopeKind vocabulary are inherited; AuditEvent expands its allowed-scope subset and ServiceInstance remains a Project-scoped typed subject |
| Audit linkage | PASS | Section 10 requires local atomic decision and audit-obligation acceptance and prevents telemetry from becoming audit authority |
| Effective context is the resolution boundary | PASS | Section 7.3 resolves hierarchy once into version-pinned `EffectivePolicyContext` with optional opaque integrity carriers and prohibits hidden traversal by decision domains |
| Provider-neutral capability matching | PASS | Decision profiles and typed results reference normalized capabilities and typed references rather than provider resources |
| Sovereign and disconnected fit | PASS | Section 11 covers data, operations, technology, law, structural trust, later cryptographic assurance, supply chain, AI/model locality, local dependencies, synchronization contracts, and portability |
| Stateless horizontal scalability | PASS | Sections 8.5 and 9 require interchangeable workers, local version-pinned inputs, bounded work, semantic idempotency, backpressure, and no sticky sessions |
| Phase 2 scope discipline | PASS | Architecture defines contracts and conformance only; production engines, stores, providers, operations, and AI runtimes remain deferred |

### 26.1 Validation conclusion

The regenerated architecture conforms to the approved Phase 2 direction. The
only intentional baseline changes are the controlled `DecisionObject` to
`DecisionRecord` terminology correction and the FEATURE-0012 `AuditEvent` scope
correction. Neither is effective until recorded human architecture approval.

The synchronous data plane has no mandatory remote registry, policy resolver,
audit index, evidence exporter, AI model, identity control plane, cryptographic
notary, provider control plane, or vendor service dependency. Remaining latency
and throughput choices are explicit, profile-bounded implementation decisions
rather than hidden architectural hot paths.

## 27. Kiro anti-overengineering guardrails

### 27.1 Mandatory scope classification

Every generated requirement, design component, and task must declare exactly
one implementation class:

- `CONTRACT_NOW` — schema, vocabulary, static registry/bundle format,
  validation, compatibility rule, fixture, conformance check, or documentation
  implemented by FEATURE-0013;
- `INVARIANT_FOR_LATER` — normative constraint that a later implementation must
  satisfy, with no FEATURE-0013 production runtime implementation;
- `DEFERRED` — product, runtime, topology, numerical SLO, or operational choice
  that is not implemented and requires a later approved feature.

An item without a class is invalid. `INVARIANT_FOR_LATER` and `DEFERRED` items
must not generate implementation tasks except documentation, traceability, and
future conformance placeholders that do not execute production behavior.

### 27.2 Allowed FEATURE-0013 artifacts

Kiro may generate only the minimum artifacts needed for:

- canonical schemas and shared types for `DecisionRecord`, `DecisionProfile`,
  `EvaluationResult`, `DecisionRationale`, and approved `AuditEvent` linkage;
- static, versioned decision-form, authority, relationship, reason, obligation,
  evaluator, strategy, evidence, event, and projection registry formats;
- strict validation, stable errors, algorithm-agile canonicalization/digest
  carrier contracts without computation, and compatibility checks;
- simple, denial, selection, authorization, composite, replay, correction,
  projection, scope, limit, and extension conformance fixtures;
- in-memory or pure-function reference composition sufficient to prove contract
  semantics, only if required by executable conformance;
- documentation, traceability, Matrix E evidence templates, and feature gates;
- the normative downstream-adoption section and lightweight adoption gate in
  section 28.

Every changed or generated artifact must map to a `CONTRACT_NOW` requirement.
An artifact justified only by a future scenario is prohibited.
`DecisionRequest` is intentionally absent from the allowed shared-contract
artifacts: downstream stages may test conformance of calling-domain inputs but
must not create a common `DecisionRequest` schema, Go type, registry entry,
resource, or persistence contract. They also must not create a pending-decision
response type; asynchronous handoff reuses `Operation`.

### 27.3 Prohibited FEATURE-0013 implementation

Kiro must not create or select:

- a production decision, authorization, policy, placement, audit, evidence,
  identity, profile-registry, synchronization, cryptographic, or AI service;
- a database schema, migration, repository implementation, transaction/outbox
  implementation, event store, WORM store, cache, distributed lock, or lease;
- a queue, broker, workflow/durable-execution runtime, scheduler, controller,
  worker pool, background reconciler, retry daemon, or dead-letter processor;
- a network registry client, remote policy resolver, external evidence fetcher,
  HSM/notary client, provider adapter, model client, vector store, or agent;
- a general-purpose rules engine, policy language, workflow DSL, DAG execution
  engine, expression interpreter, dynamic plugin system, reflection-driven
  framework, or code generator beyond established repository schema tooling;
- production authorization tokens, cache invalidation, revocation distribution,
  multi-site consensus, disconnected synchronization, or conflict-resolution
  service;
- provider-, runtime-, database-, broker-, identity-, cryptographic-, or
  model-specific dependencies in the core contract packages.

Conceptual examples, interfaces, ports, or test doubles must not perform these
behaviors or create a disguised production framework.

### 27.4 Simplicity and necessity tests

Before adding any abstraction, Kiro must answer all of the following in the
design traceability:

1. Which accepted `CONTRACT_NOW` requirement cannot be satisfied without it?
2. Which existing FEATURE-0012 primitive or repository pattern was evaluated?
3. Why is a concrete type, pure function, schema, or fixture insufficient?
4. Does it add a runtime dependency, lifecycle, persistence concern, network
   hop, background activity, extension point, or operational burden?
5. Can it be deleted while all approved conformance fixtures still pass?

If question 1 has no exact requirement ID, or question 5 is yes, the abstraction
must not be added. Reuse, direct composition, and explicit code are preferred
over a generic framework.

### 27.5 Complexity budget

- The simple atomic decision path must remain representable without a graph,
  orchestrator, registry network call, operation, AI, evidence store, or runtime
  service.
- The complex fixture proves bounded contract composition; it does not justify
  a production graph executor.
- A single implementation must not create parallel core models for simple,
  composite, authorization, placement, AI, or sovereign decisions.
- Optional assurance fields and profiles must not become mandatory costs for the
  minimal conformance profile.
- No speculative abstraction may be introduced for only one hypothetical later
  consumer.
- Limits must be finite and configurable by approved profile, but Kiro must not
  invent production defaults when architecture has not approved them.
- Unknown product choices, owners, limits, retention periods, SLOs, or legal
  values must be marked `ARCHITECTURE_DECISION_REQUIRED`; they must not be
  guessed.

### 27.6 Stop and escalation conditions

Kiro must stop and request architecture review rather than expand scope when:

- a requirement appears to need production persistence, networking, background
  processing, distributed coordination, or external integration;
- a future decision family cannot conform without changing the common envelope;
- a closed vocabulary requires a new semantic value;
- deterministic contract proof appears to require a workflow or policy engine;
- FEATURE-0012 semantics must be changed beyond the two approved corrections;
- a High or Critical Matrix E risk would be silently accepted, downgraded, or
  transferred without its required human decision;
- customer-facing contracts would expose internal provider, evaluator,
  composition, evidence, or governance machinery.

### 27.7 Required anti-overengineering verification

The FEATURE-0013 gate must verify:

- every requirement/task has one scope class;
- every implementation artifact maps to `CONTRACT_NOW`;
- prohibited runtime dependencies and packages are absent;
- no network, persistence, worker, controller, queue, cache, or external-service
  implementation was introduced;
- simple conformance works through direct bounded evaluation;
- complex conformance remains a fixture/reference proof, not a runtime engine;
- deferred decisions remain visibly deferred;
- changed-file and dependency review confirms Phase 2 contract-only scope;
- a human semantic review confirms that Kiro did not convert architecture
  invariants into premature runtime implementation.

### 27.8 Closed architecture boundary for downstream stages

Requirements, design, and tasks must treat the following as closed
architecture and may not reopen them as design questions:

- `DecisionRecord`, `EvaluationResult`, `AuditEvent`, `Operation`, and
  `DecisionContext` have the distinct meanings and owners in section 5;
- `DecisionRequest` is a calling-domain-owned conceptual input role, has no
  FEATURE-0013 shared schema/type, and only inherits the FEATURE-0012
  `TransientRequestResult` and common validation/reference boundaries;
- `metadata.scopeRef` is the sole scope authority, including canonical absent
  Platform serialization and the six-value FEATURE-0012 `ScopeKind` vocabulary;
- `ServiceInstance` is a typed subject under Project scope, never a
  `ScopeKind`; `ownerRef` remains lifecycle containment only;
- the FEATURE-0012 `DataClassification` and FEATURE-0013 sensitivity dimensions
  coexist under the mapping and fail-closed rules in section 12.4;
- declarative JSON bundles governed by JSON Schema are the canonical profile
  and semantic-registry representation; Go maps and types are derivative;
- FEATURE-0012 platform/schema ceilings are absolute outer bounds; within them,
  profiles own narrower maximum projection, sensitivity, graph, latency, cost,
  and payload limits, and evaluator registrations may only narrow the profile;
- `EvaluationResult` scope is derived from the containing record and profile;
  evaluator metadata is evidence, never a second scope authority;
- public failures inherit the FEATURE-0012 Problem Details envelope and use
  only the FEATURE-0013 error families defined in section 15;
- completed synchronous handling returns a `DecisionRecord` value/reference,
  accepted asynchronous handling returns an existing `Operation`
  value/reference, and rejected or unaccepted handling returns inherited
  Problem Details; no FEATURE-0013-specific pending-response envelope exists;
- the generic `Operation` contract and `LongRunningOperation` payload remain
  owned by FEATURE-0012/ADH-2026-013; FEATURE-0013 references but never
  redefines them;
- every AD-001 through AD-045 implementation class is fixed by section 19;
  downstream stages reproduce those exact classes and only `CONTRACT_NOW`
  produces FEATURE-0013 implementation artifacts;
- conformance coverage uses the stable parent, suffix, and additional IDs in
  section 17;
- fail-open is available only through the section 6.8 `SecurityExceptionRef`
  structural approval-evidence contract and never through an implementation-
  chosen shortcut;
- pre-ADR-F13-002 trust behavior is structural and opaque only; cryptographic
  selection and execution remain deferred; and
- FEATURE-0013 defines no adapter interface or production runtime service.

Design may decide only representation and implementation mechanics that do not
change those semantics: Go field types, JSON Schema composition, package and
import boundaries, bounded in-memory graph representation, validation-pass
organization, pure-function interfaces, and physical fixture file format.
Persistent graph storage, a graph database/schema, repository abstraction, or
runtime graph executor is not delegated to design and is prohibited in
FEATURE-0013. Requirements translate the closed architecture into testable
obligations. Tasks implement only an approved design. If a downstream stage
cannot proceed without changing a closed item, it must stop with
`ARCHITECTURE_DECISION_REQUIRED`.


### 27.9 Downstream design closure controls for clean regeneration

The following controls close the architecture gaps that previously caused
downstream design revisions to leak semantic decisions into requirements,
design, or tasks. Kiro, reviewer prompts, requirements, design, tasks, and
Cursor must treat this section as normative architecture.

1. **Fail-open exception representation.** `AD-043` remains `CONTRACT_NOW` only
   for structural profile validation and conformance of the section 6.8
   `SecurityExceptionRef` contract. The approved exception identity is the
   immutable typed reference tuple (`apiVersion`, `kind`, `name`, `uid`) plus
   bounded metadata fields for scope, owner, approving authority, compensating
   controls, effective interval, audit treatment, reassessment trigger, and
   covered failure modes. FEATURE-0013 validates shape, requiredness, scope
   compatibility, expiry against captured/effective decision time, and uid
   pinning. It does not implement exception approval, revocation, workflow,
   persistence, authorization, or runtime enforcement. A fail-open declaration
   without a valid `SecurityExceptionRef` fails profile validation; design must
   not remove the approved exception path or invent alternate approval evidence.

2. **Strategy failure posture boundary.** A composition strategy may describe
   missing, timeout, conflict, and error handling, but it is not an independent
   authority for fail-open. Any strategy-level fail-open behavior is effective
   only when the governing profile has a valid `SecurityExceptionRef`; otherwise
   it fails closed. Strategy metadata must not create a second approval path.

3. **Graph identity and accounting.** Every graph node and every graph edge has
   a stable, non-empty, profile-bounded identifier. Edge identity is distinct
   from the `(from, to)` endpoint pair; duplicate edge identifiers are invalid,
   and duplicate endpoint pairs are invalid unless a profile explicitly permits
   labelled multi-edges by identifier. Edge orientation is from dependency to
   dependent. Parallel-group independence means no direct or transitive
   dependency path exists between group members. Depth, width, fan-out,
   concurrency, payload, and weighted critical-path calculations must be
   deterministic, count each node/edge once, and use the most restrictive
   FEATURE-0012/platform, profile, and evaluator bounds.

4. **Graph version ownership.** Graph reference and version compatibility are
   validated during graph resolution against the declarative bundle before an
   in-memory graph is built. The bounded in-memory graph builder receives an
   already resolved graph definition and must not perform registry/version
   lookup or invent compatibility policy.

5. **Relationship conflict matrix.** Relationship validation must define a
   deterministic matrix for same-target, shared-ancestor, correction,
   supersession, revocation, fork, current-vs-historical, and chain-convergence
   cases. Cycles emit `DECISION_RELATIONSHIP_CYCLE`; chain excess emits
   `DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED`; incompatible relationship
   combinations emit `DECISION_RELATIONSHIP_CONFLICT`. Requirements enumerate
   cases; design maps them; tasks implement them. No downstream stage may invent
   new public relationship codes.

6. **FEATURE-0012 baseline workflow reuse.** FEATURE-0013 schema evolution,
   AuditEvent compatibility, and schema-diff evidence reuse the exact
   FEATURE-0012 baseline mechanism: `api/schemas/baseline/BASELINE_MANIFEST.json`,
   `api/schemas/baseline/BASELINE_APPROVALS.json`, digest recomputation,
   approval-controlled baseline updates, and protected review. FEATURE-0013 must
   not introduce `api/schemas/SCHEMA_BASELINE_MANIFEST.json`, `api/schemas/diffs/`,
   `api/schemas/approvals/`, or another parallel baseline/diff/approval
   hierarchy without a new approved architecture decision.

7. **AuditEvent package ownership.** FEATURE-0012 owns the generic
   `ImmutableRecord` envelope, base metadata, strict decoding, TypeBinding, and
   conformance machinery. FEATURE-0013 owns the domain AuditEvent payload/linkage
   extension. Canonical FEATURE-0013 AuditEvent value types must live in a
   FEATURE-0013 domain/shared package, while `apiconform` remains limited to
   conformance adapters, schema binding registration, and executable checks. A
   conformance package must not become the canonical domain owner and FEATURE-0013
   must not copy FEATURE-0012 base fields into a divergent parallel type.

8. **Versioned registry keys.** Declarative semantic registries are keyed by the
   tuple `(id, version)` when versioned; registries whose values are explicitly
   unversioned must say so in the profile schema. Duplicate detection, bundle
   loading, BundleView lookup signatures, fixture identities, and traceability
   must use the same key rule. Go maps are derivative views of the JSON bundle,
   not semantic authorities.

9. **TrustCarrier empty-state rule.** A trust or integrity carrier that is
   absent is different from one that is present but structurally empty. When a
   carrier is required by profile, bundle, evaluation provenance, synchronization,
   or import/export semantics, absence or structural emptiness fails closed with
   the appropriate existing trust/evaluation code. When a carrier is optional, a
   present but structurally empty value must either be rejected as invalid input
   or normalized to absent by a single documented rule; it must never satisfy
   provenance, required trust, identity generation, or compatibility evidence.

10. **Projection pointer target and overlap.** Projection pointers are resolved
    against the canonical schema/view declared by the governing profile, not
    against fields that happen to appear in one example record. Pointer overlap
    includes exact equality and ancestor/descendant containment. Unknown pointers,
    overlapping includes/excludes, audience mismatch, and prohibited-category
    exposure fail deterministically with existing projection/profile/evaluation
    codes; requirements and design must not create executable redaction or
    projected-output generation.

11. **Requiredness for public contracts.** For every public schema or externally
    exchanged contract, architecture requires downstream design to enumerate
    requiredness for every string, boolean, number, object, pointer, map, and
    slice/array. Required-present slices must state whether empty is allowed;
    nullable optional references must identify whether `null`, absence, or both
    are canonical. Go tags, JSON Schema `required`/`minItems`/`nullable`, loaders,
    validators, fixtures, and error pointers must agree.

12. **Matrix E automation boundary.** Automation may stage architecture-owned
    risk IDs, inherent ratings, controls, evidence pointers, target residual
    ratings, treatment, owner, trigger, and blank human-disposition fields. It
    must not calculate, downgrade, upgrade, infer, or replace residual risk
    levels from test coverage. Human-only disposition fields remain blank until
    recorded human review.

13. **Scope pre-scan and decode reuse.** FEATURE-0013 must reuse exported
    FEATURE-0012 decode/pre-scan primitives where they exist and the package DAG
    permits. JSON duplicate/parallel-key checks occur on JSON tokens. YAML
    duplicate/parallel-key checks occur on the `yaml.Node` AST before
    JSON-compatible normalization. If no exported FEATURE-0012 primitive is
    available for a required check, FEATURE-0013 may implement a local pure
    pre-scan only within the approved package import graph and must not claim
    zero allocation when YAML AST parsing is performed.

14. **Clean regeneration rule.** A regenerated design must treat this
    consolidated architecture as source of truth and must not reuse superseded or
    patched design drafts as semantic input. If any item above cannot be mapped
    to concrete design mechanics without changing architecture, the design must
    stop with `ARCHITECTURE_DECISION_REQUIRED` and must not proceed to tasks.

## 28. Lightweight downstream adoption contract

### 28.1 Purpose

Later features must reuse FEATURE-0013 consistently without creating an
adoption-management system. The initial mechanism is deliberately limited to:

1. this one normative contract in FEATURE-0013;
2. one short `FEATURE-0013 Adoption` section in each directly dependent or
   semantically adopting feature's architecture document; and
3. one lightweight feature-gate check.

FEATURE-0013 does not require separate per-feature YAML manifests, a central
adoption registry, a generated adoption index, a separate adoption approval, or
a runtime adoption service. Those artifacts are prohibited until a reassessment
trigger in section 28.8 is met and a human approves the additional governance
cost.

### 28.2 Applicability decision

A feature must include the short adoption section when it directly depends on
FEATURE-0013 or produces, consumes, specializes, projects, audits, or enforces
any of the following:

- `EvaluationResult`;
- `DecisionProfile`;
- `DecisionRecord` or a specialized decision;
- decision rationale or obligations;
- decision-linked `AuditEvent`;
- decision-linked `Operation`;
- `DecisionContext` or another decision projection;
- a decision composition strategy.

The section classifies the feature as `APPLICABLE` or `NOT_APPLICABLE`. A
`NOT_APPLICABLE` declaration requires a concise rationale and must not be used
when the feature performs any behavior listed above.

### 28.3 Adoption roles

An applicable feature declares only the roles it performs:

- `EVALUATION_PRODUCER`;
- `DECISION_PRODUCER`;
- `DECISION_CONSUMER`;
- `AUDIT_PRODUCER`;
- `PROJECTION_CONSUMER`;
- `OPERATION_OWNER`;
- `COMPOSITE_ADOPTER`.

The roles are documentation and conformance classifications, not runtime
resources or services.

### 28.4 Required per-feature section

An applicable feature uses this minimal architecture section:

```markdown
## FEATURE-0013 Adoption

Applicability: APPLICABLE

Roles:
- DECISION_PRODUCER
- AUDIT_PRODUCER

Profiles introduced or reused:
- <profile name and version>

Contracts reused:
- <FEATURE-0013 objects/vocabularies>

Inherited material risks:
- <applicable F13 risk IDs and disposition in this feature>

Exceptions:
- None

Conformance evidence:
- <fixtures/checks>
```

A non-adopting feature uses only:

```markdown
## FEATURE-0013 Adoption

Applicability: NOT_APPLICABLE

Rationale: <why the feature produces or consumes none of the governed concepts>
```

The feature may add concise explanation, but Kiro must not expand this section
into a separate adoption specification unless section 28.8 is triggered.

### 28.5 Mandatory downstream rules

An applicable feature must:

- reference an exact FEATURE-0013 contract version;
- reuse the common envelope without redefining its semantics;
- register and version each new `DecisionProfile` and typed result;
- declare form, authority, failure behavior, obligations, audit linkage, scope,
  projection, identity, and applicable bounds;
- preserve evaluation, decision, audit, operation, and projection distinctions;
- identify applicable inherited High risks, risks whose controls it implements,
  and new risks it introduces;
- provide positive and negative conformance evidence;
- record any exception through the feature's existing architecture approval.

A downstream feature must not create a second common envelope, treat evaluator
output or AI output as automatic authority, use operation state as a decision,
use telemetry as audit, ignore mandatory obligations, or locally expand closed
FEATURE-0013 vocabularies.

### 28.6 Selective risk inheritance

A downstream feature does not copy all Matrix E risks. It references only:

- risks directly applicable to its behavior;
- applicable High target-residual risks;
- risks whose controls it implements or changes;
- risks transferred to it by an earlier feature; and
- new implementation risks it introduces.

Each referenced risk is marked `APPLICABLE`, `NOT_APPLICABLE` with rationale,
`MITIGATED_BY_FEATURE`, or `TRANSFERRED_FORWARD`. Acceptance of a FEATURE-0013
architecture risk does not automatically accept implementation-specific
residual risk in the downstream feature.

### 28.7 Lightweight gate and unified approval

The ordinary feature gate checks only that:

- the section exists when applicability rules require it;
- applicability and roles are stated;
- profiles are named and versioned;
- reused contracts are identified;
- applicable material risks and exceptions are explicit;
- conformance evidence is identified; and
- the common envelope and closed vocabularies were not redefined.

Adoption is reviewed within the downstream feature's existing architecture and
semantic review. There is no separate adoption approval, reviewer, workflow, or
approval document.

### 28.8 Reassessment triggers for structured manifests

Separate machine-readable manifests or a generated central index may be
proposed later only when observed repository needs justify them, such as:

- at least three or four active features define independently versioned
  decision profiles;
- duplicate or conflicting profiles appear;
- manual review repeatedly misses semantic drift;
- cross-feature compatibility automation requires structured data;
- obligations are produced and enforced by different feature owners; or
- stable/external API promotion requires machine-readable discovery.

Meeting a trigger permits evaluation; it does not automatically approve the
additional framework. The proposal must demonstrate that a small extension of
the existing feature gate is insufficient, quantify maintenance and review
cost, and receive human architecture approval.

## 29. Architecture decision history

This section records provenance only. It does not contain separate normative
requirements or override the consolidated architecture above.

- `ADH-2026-014` established the initial FEATURE-0013 architecture baseline.
- `ADH-2026-015` attempted to reconcile the `AuditEvent` scope model by adding
  ServiceInstance to `ScopeKind`; the consolidated architecture retains its
  single-scope-authority intent but replaces that semantic choice with a typed
  ServiceInstance subject under Project scope.
- `ADH-2026-016` clarified `DecisionRecord` scope authority, sensitivity,
  structural security validation, trust boundaries, conformance identifiers,
  evidence disposition, and terminology migration.
- `ADH-2026-017` is the approved single replacement handoff for this complete
  consolidated architecture. When approved, it supersedes ADH-2026-014,
  ADH-2026-015, and ADH-2026-016 in full as normative inputs while retaining
  them as historical provenance.
- The consolidated main body is authoritative; the handoff is its approval
  envelope and binds approval to the exact architecture content.
- If the main body conflicts with superseded wording retained in historical
  handoffs or evidence, the consolidated main body controls.
- Future semantic changes require a new approved architecture decision; Kiro,
  design, tasks, and implementation must not reinterpret historical wording.
