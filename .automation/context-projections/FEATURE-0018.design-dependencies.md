---
doc_type: ai_context_projection
feature: FEATURE-0018
stage: design
authority: non-authoritative-exact-excerpts
generated: true
---

# FEATURE-0018 Design Context Projection

This is a non-authoritative, mechanically generated projection of exact
approved-source excerpts. The FEATURE architecture and controlling handoff
package remain authoritative. A source-hash mismatch blocks regeneration.

## Source: `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`

Approved SHA-256: `325c20278da09149dabe7f4d0cee583269f78959a37294a98eb5bceeb13f6097`

### Canonical DecisionRecord, DecisionProfile and FEATURE-0018 adoption contracts

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
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

### 6.9.1 FEATURE-0018 adopting-domain registration

Under accepted DEC-0060, FEATURE-0018 registers exactly
`authorization-decision/v1`, `approval-decision/v1`, and
`exception-decision/v1` through this document's existing DecisionProfile
extension mechanism. FEATURE-0013 retains envelope, validity, authority,
projection, obligation, publication and audit-linkage ownership. AccessReview
remains a FEATURE-0018 LongRunningOperation and is not a fourth DecisionProfile.

The exact input, result, validity and mandatory-emission semantics are
F18-RD-19; the audit-before-publication boundary and exact event taxonomy are
F18-RD-20. This registration does not change the FEATURE-0013 envelope or
authorize any downstream effect. Executable use remains gated by separately
approved FEATURE-0018 requirements, design, tasks and implementation.

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
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

### Evaluation composition, execution neutrality and data locality

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
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
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

### Audit consistency contract

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
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
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

### Security, privacy and trust boundaries

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
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
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

### Compatibility and observability adoption constraints

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
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
<!-- END EXACT APPROVED-SOURCE EXCERPT -->

### Lightweight downstream adoption contract

<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->
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
<!-- END EXACT APPROVED-SOURCE EXCERPT -->
