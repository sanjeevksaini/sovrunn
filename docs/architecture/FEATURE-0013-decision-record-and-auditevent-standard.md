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

This document is the approved architecture input from which Kiro may generate
FEATURE-0013 requirements only. It is not an implementation specification and
does not authorize design, tasks, source code, or runtime execution.

Human architecture review status: **APPROVED_FOR_KIRO_REQUIREMENTS**.

- Reviewer: Sanjeev Kumar
- Decision date: 2026-07-27
- Controlling handoffs: ADH-2026-014 and ADH-2026-015 (joint)

> **Scope vocabulary clarification (ADH-2026-015, Approved 2026-07-27).**
> ADH-2026-014 and ADH-2026-015 are the joint controlling handoffs for
> FEATURE-0013. ADH-2026-015 supersedes **only** the earlier six-scope
> `AuditEvent` statement in ADH-2026-014. All other ADH-2026-014 decisions
> remain controlling and unchanged. Under ADH-2026-015, `ScopeKind` is the
> single canonical shared scope vocabulary with exactly seven values
> (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`,
> `Provider`, `ServiceInstance`); `ServiceInstance` is additive and `Provider`
> is retained. `AuditEvent` initially permits all seven values and uses
> `metadata.scopeRef` as its sole canonical scope identity. See section 29 for
> the full clarification, the value-by-value compatibility table, and the
> mandatory Dependency Contract Reconciliation framework. Where any earlier
> statement in this document refers to "six" `AuditEvent` scopes, section 29 is
> authoritative and controls.

The reviewer approves AD-001 through AD-044, the two controlled baseline
corrections, contract-only Phase 2 scope, anti-overengineering guardrails,
lightweight downstream adoption contract, and the proposed treatment,
verification, ownership, and reassessment plans for Matrix E F13-R01 through
F13-R31. This architecture approval does not accept final residual risk. Final
per-risk residual acceptance remains pending implementation evidence and human
semantic review. F13-R04, F13-R08, and F13-R13 remain explicitly open at target
High residual level and require reassessment before their corresponding
production capabilities are enabled.

### 1.1 Change classification

This proposal is an **extension with two explicit corrections** to the approved
Phase 2 spine:

1. `DecisionRecord` replaces the ambiguous common name `DecisionObject`.
2. `AuditEvent` scope is expanded from the FEATURE-0012 Organization-only alpha
   profile to all seven canonical `ScopeKind` values (`Platform`,
   `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`,
   `ServiceInstance`) per ADH-2026-015, using `metadata.scopeRef` as the sole
   canonical scope identity. `ServiceInstance` is additive to the FEATURE-0012
   shared vocabulary and `Provider` is retained. The earlier ADH-2026-014
   statement that expanded `AuditEvent` to "six governance scopes already
   established by FEATURE-0012" is superseded and was factually incorrect:
   FEATURE-0012 established the shared `ScopeKind` vocabulary, not a six-scope
   `AuditEvent` profile. Section 29 (ADH-2026-015) is authoritative.

ADH-2026-014 approves the terminology correction, and ADH-2026-015 approves the
seven-value `ScopeKind`/`AuditEvent` scope correction, for FEATURE-0013
requirements generation.
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
| `DecisionRequest` | Bounded request to decide, with context references and constraints | Calling domain | No |
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

## 6. Extensible decision semantics

### 6.1 Stable common envelope

Every decision uses the same `DecisionRecord` envelope. It identifies the
versioned decision profile, primary form, authority, purpose, request, actors,
subjects, scope, inputs, evaluations, composition, typed result, rationale,
obligations, provenance, validity, correlation, and audit links.

The envelope does not contain every domain's result fields. `DecisionProfile`
selects a separately versioned typed-result schema. An adopting feature may add
a profile and result schema without revising FEATURE-0013 when it conforms to
the extension contract and existing semantic registries.

```yaml
kind: DecisionRecord
apiVersion: sovrunn.io/v1alpha1
profileRef:
  name: PlacementDecision
  version: 1.0.0
form: SELECTION
authority: AUTHORITATIVE
purpose: workload-placement
requestRef: null
actorRef: null
subjectRefs: []
scopeRef: null
effectivePolicyContextRef: null
inputSnapshotRef: null
evaluationResultRefs: []
composition:
  strategyRef: null
  graphRef: null
result:
  finality: FINAL
  adjudication: null
  typedResult: {}
  rationale: {}
  obligations: []
provenance: {}
validity: {}
correlation:
  operationRef: null
  auditEventRefs: []
```

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
erasure may remove protected payloads only through an approved cryptographic
erasure/redaction process that preserves a verifiable tombstone and chain
integrity where law permits.

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
approved, time-bounded security exception identifies scope, owner,
compensating controls, expiry, audit treatment, and reassessment trigger.

Profiles and semantic registries must be distributable and verifiable offline.
Registration requires conformance review; accepting arbitrary profile URIs at
runtime would create an ungoverned extension and is prohibited.

Registry publication and approval are control-plane operations. Runtime workers
load signed, versioned profile bundles locally and pin profile name, version,
schema digest, and trust-root identity in each decision. A synchronous decision
must not require a network lookup to the profile or semantic registry. Unknown,
revoked, expired, or digest-mismatched profiles are rejected or quarantined.

FEATURE-0013 defines `DecisionProfile` schemas, governance, signed bundle
format, validation, and conformance only. It does not implement or authorize a
production network registry, distribution service, or control plane in Phase 2.
Production publication, replication, revocation distribution, and runtime
registry services are deferred to separately approved features.

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
digest-pinned context, verified short-lived decision artifacts, or safe caches
when the profile defines audience, scope, expiry, revocation, freshness, and
fail-closed behavior. Aggregation may reduce durable `DecisionRecord` volume but
must not suppress audit evidence required by the applicable profile.

## 7. Evaluation and composition model

### 7.1 Atomic evaluation

An atomic evaluation has one registered evaluator type, one bounded input
snapshot or digest, one result, one timing envelope, and explicit success,
non-match, indeterminate, timeout, or error semantics. Evaluator-native payloads
remain behind adapters; normalized fields enter the common record.

Evaluation reproducibility and decision reproducibility are distinct. An
external, probabilistic, or AI evaluator may not reproduce the same output when
rerun. Its normalized result must therefore capture evaluator identity and
version, governed input snapshot or digest, returned output, relevant execution
configuration, evaluation time, trust boundary, safety/policy filters, and
integrity digest. Authoritative composition is deterministic from these
captured results and must not rerun an external or AI evaluator during replay or
audit. A profile may prohibit nondeterministic evaluators entirely.

### 7.2 Composite decision

A composite decision combines bounded evaluation results using a versioned,
registered strategy. Every strategy declares:

- stable identifier and version;
- accepted input types;
- ordering and precedence;
- short-circuit behavior;
- missing, timeout, conflict, and error behavior;
- fail-open or fail-closed posture;
- deterministic output mapping;
- resource budgets and maximum fan-out;
- explanation and evidence rules.

Fail-open remains subject to the profile restrictions in section 6.8 and can
never be selected merely as a runtime fallback.

The chosen strategy identifier, version, ordered input digests, and output are
recorded so the result can be reproduced.

### 7.3 Hierarchical evaluation

Governance inheritance is resolved before final decision composition into one
`EffectivePolicyContext`. The decision layer must not independently traverse
Platform, Organization, OrganizationUnit, Tenant, Project, or Provider trees.
The resolved context includes provenance, precedence, conflict resolution,
effective time, and digest. This prevents every decision domain from creating a
different hidden inheritance engine.

A valid immutable context snapshot may be resolved ahead of the request and
cached by scope, version, effective-time window, and digest. The synchronous
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
- recorded graph definition/version or immutable digest.

The graph may describe dependency and composition only. It cannot execute
provisioning, human tasks, compensations, arbitrary code, or unbounded loops.
Short-circuiting must cancel or disregard unnecessary work without changing
deterministic results, leaking cross-scope information, or leaving unbounded
background activity.

## 8. Execution model

### 8.1 Synchronous path

Simple and bounded decisions should complete synchronously when all required
inputs are available within the declared latency budget. The path does not
require a queue or durable workflow engine. It returns either a final
`DecisionRecord` reference/value or a typed non-final response indicating that
asynchronous handling is required.

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
version-pinned, and digest-verifiable. Profile publication, policy resolution,
schema distribution, evidence export, audit indexing, AI inference,
notarization, and regulator replication are control-plane or asynchronous
activities unless a specific profile explicitly makes one part of the decision
input and budgets its failure behavior.

No remote registry, audit index, evidence exporter, AI model, identity control
plane, cryptographic notary, or vendor control plane may become an implicit
universal synchronous dependency. A missing mandatory input causes a typed
non-final response, deterministic fail-closed adjudication, or transition to an
`Operation` according to the profile; it never causes an unbounded synchronous
wait.

## 9. Statelessness, scalability, and cost controls

The decision kernel is a pure, deterministic function of versioned inputs,
registered strategy, and bounded configuration. Workers hold no authoritative
local session state between requests and are interchangeable behind ordinary
load balancing.

Durable state, deduplication, locks/leases, input snapshots, and audit records
are external services accessed through ports. Implementations must support:

- a caller-supplied idempotency key scoped to caller and request purpose for
  retry identity;
- a separate semantic decision digest covering scope, purpose, profile name and
  version, authority, canonical input snapshot, `EffectivePolicyContext`,
  composition strategy/version, evaluator contract versions, and relevant
  evaluation-time or validity epoch;
- compare-and-create or equivalent concurrency control;
- at-least-once delivery without duplicate authoritative decisions;
- bounded retries with jitter and dead-letter/quarantine handling;
- backpressure, admission control, per-scope quotas, and load shedding;
- deterministic caching only by the complete semantic decision digest, with
  audience, scope, expiry, revocation, and freshness enforcement;
- no dependence on sticky sessions;
- scale-to-zero where deployment policy permits;
- hot-path summaries with cold evidence referenced by digest;
- batch and parallel evaluation only within declared resource budgets.

Exactly-once transport is not claimed. Business-level effective-once behavior
comes from idempotency, immutable identity, and atomic persistence boundaries.

## 10. Audit consistency

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
event-sourced append, or another atomic pattern—is deferred. The synchronous
path must not depend on a remote audit service or index. Implementations must
expose recovery, bounded delivery, poison-record quarantine, and reconciliation
for partial downstream failure and must never fabricate success from a
telemetry log.

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
AI inference if AI is enabled. Synchronization is explicit, signed,
policy-authorized, resumable, and conflict-detecting; it is never a hidden
foreign control-plane dependency.

Export and import use collision-resistant record identifiers, origin trust-
domain identity, signed manifests, replay protection, duplicate detection,
per-origin ordering, and explicit trust-root version. Unknown profiles, keys,
origins, or revoked trust roots are quarantined. Immutable records are never
merged using last-writer-wins; cross-origin conflicts produce explicit evidence
and a governed reconciliation decision.

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
- integrity digest and trusted timestamp provenance;
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

Canonical serialization and content digests are supported. Deployment profiles
may require signatures, trusted timestamps, append-only/WORM retention, Merkle
batch proofs, or external notarization. The core contract remains
algorithm-agile and does not mandate one national or vendor cryptographic
service.

Local digest computation is permitted on every record. Per-record remote HSM,
trusted-time, or notarization calls are not universal hot-path requirements.
Profiles may require local signing, per-record signing, batched Merkle signing,
or asynchronous external notarization according to assurance level. They must
declare latency, outage, key-rotation, revocation, and fail-closed behavior.

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

| Foundation | Decision |
|---|---|
| FEATURE-0012 API/resource grammar | Extend; do not create a parallel metadata, reference, error, or validation grammar |
| OMG DMN concepts | Reuse bounded decision-requirements concepts; do not select a DMN engine |
| CloudEvents | Optional transport mapping only; never the canonical audit or decision record |
| OpenTelemetry | Correlation and operational telemetry only; never audit authority |
| RFC 8785 JCS | Candidate canonicalization profile for digests; exact covered fields decided in design |
| DSSE/in-toto and SPDX/CycloneDX concepts | Reuse for signed supply-chain evidence references; do not duplicate their schemas |
| CEL/OPA/Cedar | Future evaluator candidates behind adapters; no engine-native core types |
| Workflow/durable execution systems | Later runtime candidates behind operation/orchestration ports |

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
- Migrations are reversible where practical and preserve old evidence/digests.
- Per ADH-2026-015, `AuditEvent` expansion requires fixtures for all seven
  canonical `ScopeKind` values (`Platform`, `Organization`, `OrganizationUnit`,
  `Tenant`, `Project`, `Provider`, `ServiceInstance`), separate regression
  fixtures for the six pre-existing FEATURE-0012 `ScopeKind` values (`Platform`,
  `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`), and a
  documented alpha compatibility decision.

## 16. Observability and operability

Metrics and traces expose latency, queueing, evaluation fan-out, cache result,
composition strategy, failure class, retry, saturation, and audit-reconciliation
health without exposing decision evidence or personal data as labels.

Logs are structured, redacted, locally routable, retention-controlled, and
correlated by opaque identifiers. Health must distinguish evaluator, evidence,
store, audit, identity, key, time, and synchronization dependencies. Operators
must be able to diagnose and recover the system without vendor remote access.

## 17. Conformance model

The architecture requires executable positive and negative fixtures for:

1. simplest synchronous atomic authorization decision;
2. simplest deterministic denial with actionable rationale and obligations;
3. pure selection with a typed result and no adjudication facet;
4. ranking, classification, resolution, allocation, plan, assessment,
   recommendation, advisory, and simulation profiles;
5. human-approval asynchronous decision linked to an `Operation`;
6. maximally complex but bounded hierarchical composite decision;
7. parallel evaluations with deterministic aggregation;
8. timeout, missing evidence, evaluator failure, and policy conflict;
9. idempotent replay and concurrent duplicate requests;
10. partial decision/audit persistence failure and reconciliation;
11. supersession, correction, revocation, cycle, and chain-limit behavior;
12. unauthorized profile, authority, projection, obligation bypass, and
    cross-scope reference rejection;
13. offline/air-gapped profile registry, trust, and signed-bundle operation;
14. export/import across two conforming provider implementations;
15. AI-disabled operation and AI projection redaction;
16. graph size, depth, width, fan-out, timeout, and payload budget rejection;
17. registration of a new decision family without a common-envelope change;
18. immutable supersession and revocation with calculated effective state;
19. local atomic decision/audit-obligation acceptance while downstream audit
    delivery is unavailable;
20. idempotency retry identity versus semantic decision identity across policy,
    profile, strategy, authority, evaluator-version, and validity changes;
21. deterministic replay from captured nondeterministic or AI evaluation output
    without rerunning the evaluator;
22. synchronous operation with profile registry and policy resolver unavailable
    but valid local digest-pinned bundles present;
23. rejection of unknown, expired, revoked, or digest-mismatched profile bundles;
24. enforcement denial for unknown or unsupported mandatory obligations;
25. local authorization artifact cache hit, expiry, revocation, scope mismatch,
    and stale-policy failure;
26. cancellation and total remote-call, byte, concurrency, retry, jurisdiction,
    time, and cost budgets for composite graphs;
27. local digest plus batched signing/notarization during remote cryptographic-
    service outage;
28. disconnected export/import replay, duplicate, ordering, unknown trust root,
    revoked origin, and explicit conflict reconciliation.

## 18. Matrix E v2 — architecture risk register

### 18.1 Compatibility with FEATURE-0012

Matrix E v2 preserves FEATURE-0012's stable risk identifiers, primary controls,
detection/response evidence, automated evidence staging, human-only residual
acceptance, owner, corrective path, and reassessment trigger. It adds explicit
category, affected architecture decisions, inherent and residual assessment,
treatment, separate detection and response, control owner, later-feature owner,
acceptance expiry, and per-risk status.

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
| F13-R04 | Correctness | Stale `EffectivePolicyContext` produces an invalid decision | AD-007, AD-038 | 4×5 Critical | Digest, effective time, freshness, revocation, fail-closed rules | Cache expiry/revocation/staleness tests | 2×5 High | MITIGATE | Effective-context owner | Production policy resolver or new inheritance rule | PENDING_HUMAN_REVIEW |
| F13-R05 | Correctness | Weak retry identity or cache identity reuses the wrong decision | AD-013, AD-036 | 4×5 Critical | Separate retry key and complete semantic digest | Policy/profile/strategy/time mutation tests | 1×5 Medium | MITIGATE | Decision persistence owner | Canonical digest profile change | PENDING_HUMAN_REVIEW |
| F13-R06 | Audit/compliance | A decision is durable without its audit obligation | AD-014, AD-035 | 3×5 High | Same-boundary atomic acceptance and preallocated audit identity | Partial-failure and reconciliation fixtures | 1×5 Medium | MITIGATE | Audit persistence owner | Production persistence selection | PENDING_HUMAN_REVIEW |
| F13-R07 | Availability/latency | Remote audit processing blocks synchronous decisions | AD-008, AD-035 | 4×4 High | Local obligation acceptance; asynchronous downstream processing | Remote-audit outage and backlog tests | 1×4 Low | AVOID | Audit implementation owner | External audit/SIEM integration | PENDING_HUMAN_REVIEW |
| F13-R08 | Cost/availability | Composite graphs cause fan-out, latency, retry, or cost explosion | AD-006, AD-040 | 4×5 Critical | Bounds on graph, calls, bytes, authorities, concurrency, retries, time, jurisdiction, and cost | Limit, cancellation, load, and adversarial fixtures | 2×5 High | MITIGATE | Decision runtime owner | Limit increase or new remote evaluator | PENDING_HUMAN_REVIEW |
| F13-R09 | Correctness/AI | Nondeterministic evaluator output cannot be replayed or explained | AD-005, AD-037 | 4×4 High | Captured normalized output, provenance, input and integrity digest | Replay without evaluator access | 2×3 Medium | MITIGATE | Evaluator/AI owner | New probabilistic evaluator | PENDING_HUMAN_REVIEW |
| F13-R10 | Security/AI | AI receives unauthorized evidence or sensitive canonical data | AD-019, AD-020 | 4×5 Critical | Authorized `DecisionContext`, minimization, classification, projection policies | Cross-scope/redaction and prompt-leak fixtures | 1×5 Medium | MITIGATE | AI-context and security owners | New model, projection, or jurisdiction | PENDING_HUMAN_REVIEW |
| F13-R11 | Security | Enforcement silently ignores an unknown mandatory obligation | AD-032, AD-043 | 3×5 High | Unknown/unsupported obligation makes allow unenforceable | Obligation capability negative tests | 1×5 Medium | AVOID | Enforcement owner | New obligation vocabulary | PENDING_HUMAN_REVIEW |
| F13-R12 | Latency/availability | Central authorization or decision service becomes a universal bottleneck | AD-008, AD-039 | 4×5 Critical | Local verified artifacts, policy snapshots, caches, bounded freshness | Central-service outage and horizontal-load tests | 2×4 Medium | MITIGATE | Authorization implementation owner | Production enforcement integration | PENDING_HUMAN_REVIEW |
| F13-R13 | Security | Revoked or stale cached authorization remains enforceable | AD-039, AD-043 | 3×5 High | Audience/scope/expiry/revocation/freshness validation; fail closed | Cache expiry, revocation, scope-mismatch tests | 2×5 High | MITIGATE | Authorization/security owner | Revocation architecture selection | PENDING_HUMAN_REVIEW |
| F13-R14 | Security/availability | Fail-open behavior exposes protected resources during dependency failure | AD-005, AD-043 | 3×5 High | Fail-closed high-risk profiles; bounded approved exception | Failure injection and exception-expiry checks | 1×5 Medium | AVOID | Security architecture owner | Any fail-open request | PENDING_HUMAN_REVIEW |
| F13-R15 | Correctness | Correction, supersession, or revocation yields ambiguous effective state | AD-004, AD-034 | 3×4 High | Immutable linked effects; bounded cycle-free deterministic projection | Chain, cycle, fork, and projection fixtures | 1×4 Low | MITIGATE | Read-model owner | New relationship/effective-state semantics | PENDING_HUMAN_REVIEW |
| F13-R16 | Privacy/security | Decision, rationale, evidence, error, or projection leaks restricted data | AD-019, AD-024 | 4×5 Critical | References/digests, minimization, classification, redaction, purpose projections | Secret/PII/cross-scope negative tests | 1×5 Medium | MITIGATE | Security/privacy owner | New evidence or projection type | PENDING_HUMAN_REVIEW |
| F13-R17 | Compliance/integrity | Erasure or retention action breaks evidence integrity or legal hold | AD-004, AD-022 | 3×5 High | Cryptographic erasure, tombstone, retention/legal-hold profile | Hold, erasure, digest, and tombstone tests | 2×4 Medium | MITIGATE | Evidence/retention owner | Regulated deployment or retention change | PENDING_HUMAN_REVIEW |
| F13-R18 | Latency/availability | Remote HSM, time, signing, or notarization becomes a hot-path dependency | AD-022, AD-041 | 4×4 High | Local digest; tiered local/batch/asynchronous assurance | Trust-service outage and batch-proof tests | 1×4 Low | AVOID | Security infrastructure owner | High-assurance profile adoption | PENDING_HUMAN_REVIEW |
| F13-R19 | Sovereignty/security | Disconnected import accepts forged, replayed, or untrusted records | AD-018, AD-042 | 3×5 High | Signed manifests, origin trust, replay protection, ordering, quarantine | Forgery/replay/unknown-root fixtures | 1×5 Medium | MITIGATE | Synchronization/security owner | First disconnected deployment | PENDING_HUMAN_REVIEW |
| F13-R20 | Correctness | Disconnected origins create conflicting authoritative decisions | AD-018, AD-042 | 3×5 High | Immutable origin identity; no last-writer-wins; governed reconciliation decision | Concurrent-origin conflict fixtures | 2×4 Medium | MITIGATE | Multi-site owner | Active/active or disconnected write enablement | PENDING_HUMAN_REVIEW |
| F13-R21 | Portability | Provider/runtime-native identifiers or schemas leak into the core | AD-011, AD-023 | 4×4 High | Adapter-only native data; normalized capabilities; typed references | Import/schema lint and portability fixtures | 1×4 Low | AVOID | Adapter architecture owner | First provider/runtime adapter | PENDING_HUMAN_REVIEW |
| F13-R22 | Scope/governance | Production registry, workflow, persistence, provider, or AI runtime leaks into Phase 2 | AD-011, AD-028 | 4×4 High | Contract-only scope, non-goals, feature gates | Changed-file, dependency, and no-side-effect gates | 1×4 Low | AVOID | FEATURE-0013 owner | Any runtime dependency proposal | PENDING_HUMAN_REVIEW |
| F13-R23 | Cost | Durable decisions, audit, evidence, or projections create uncontrolled storage growth | AD-024, AD-033 | 4×4 High | Risk-based recording, bounded payloads, hot summaries, referenced cold evidence, retention profiles | Volume/limit/retention tests and metrics | 2×4 Medium | MITIGATE | Audit/evidence owner | Retention increase or production volume | PENDING_HUMAN_REVIEW |
| F13-R24 | Supply chain/security | Profile, schema, strategy, evaluator, or trust bundle is compromised | AD-025, AD-038 | 3×5 High | Signed digest-pinned bundles, trust roots, revocation, quarantine, provenance | Signature/tamper/revocation fixtures | 1×5 Medium | MITIGATE | Supply-chain security owner | Signing-root or distribution change | PENDING_HUMAN_REVIEW |
| F13-R25 | Correctness/security | Clock uncertainty invalidates ordering, expiry, freshness, or replay protection | AD-018, AD-022 | 3×5 High | Trusted-time provenance, validity epoch, profile outage behavior | Skew, rollback, expiry-boundary tests | 2×4 Medium | MITIGATE | Runtime/security owner | Disconnected or multi-site deployment | PENDING_HUMAN_REVIEW |
| F13-R26 | Security/privacy | Cross-scope references disclose or influence another tenant or governance domain | AD-015, AD-019 | 3×5 High | Scope authorization, no-existence disclosure, typed references, projections | Cross-scope and safe-denial fixtures | 1×5 Medium | MITIGATE | API/security owner | New scope/reference kind | PENDING_HUMAN_REVIEW |
| F13-R27 | Compatibility/governance | Later features silently redefine common decision semantics | AD-029, AD-030 | 4×5 Critical | Stable envelope, profile conformance, architecture change control | Baseline/schema semantic-diff gate | 1×5 Medium | AVOID | Architecture owner | Common-envelope or closed-vocabulary change | PENDING_HUMAN_REVIEW |
| F13-R28 | Compatibility | Seven-value `ScopeKind`/`AuditEvent` correction (ADH-2026-015) breaks FEATURE-0012 consumers | AD-015 | 3×4 High | Explicit correction approval, versioning, seven-value fixtures, six-value regression fixtures, migration evidence | Seven-value `ScopeKind`/`AuditEvent` compatibility and six-value regression baseline tests | 1×4 Low | MITIGATE | FEATURE-0013/API owner | Stable API promotion | PENDING_HUMAN_REVIEW |
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

| ID | Architectural decision | Long. | Correct. | Cost | Security | AI-first | Reuse | Recommendation |
|---|---|---:|---:|---:|---:|---:|---:|---|
| AD-001 | Rename common `DecisionObject` to `DecisionRecord` | 5 | 5 | 5 | 4 | 5 | 5 | Approve with spine/traceability update |
| AD-002 | Keep evaluations, decisions, audit, operations, and projections distinct | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-003 | Adjudication outcomes are optional, facet-specific, and closed | 5 | 5 | 5 | 5 | 5 | 5 | Approve; keep errors and typed results separate |
| AD-004 | Immutable final records with linked corrections | 5 | 5 | 4 | 5 | 5 | 5 | Approve with lawful-erasure profile |
| AD-005 | Deterministic versioned composition strategies | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-006 | Bounded DMN-like dependency DAG, not workflow language | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-007 | Resolve hierarchy once into `EffectivePolicyContext` | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-008 | Synchronous fast path without mandatory queue | 5 | 5 | 5 | 4 | 4 | 5 | Approve with strict budgets |
| AD-009 | `Operation` owns asynchronous lifecycle | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-010 | Orchestrate authority; choreograph reactions | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-011 | Runtime/vendor selection deferred | 5 | 5 | 5 | 5 | 4 | 5 | Approve; define later selection gate |
| AD-012 | Pure stateless decision kernel | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-013 | Effective-once via idempotency, not exactly-once claims | 5 | 5 | 5 | 5 | 4 | 5 | Approve |
| AD-014 | Decision and audit durable-acceptance invariant | 5 | 5 | 4 | 5 | 5 | 5 | Approve with later persistence pattern |
| AD-015 | Extend `AuditEvent` to all seven canonical `ScopeKind` values (`Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`, `ServiceInstance`) per ADH-2026-015, using `metadata.scopeRef` as sole scope identity | 5 | 5 | 4 | 5 | 5 | 5 | Approve with FEATURE-0012 compatibility correction (ADH-2026-015) |
| AD-016 | Sovereignty as seven enforceable dimensions | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-017 | Connected, private, edge, and air-gapped profiles | 5 | 5 | 3 | 5 | 5 | 5 | Approve; tier conformance to control cost |
| AD-018 | Local control-plane dependencies and signed synchronization | 5 | 5 | 3 | 5 | 5 | 5 | Approve; advanced profile only |
| AD-019 | Authorized projections and data minimization | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-020 | AI consumes `DecisionContext`, never canonical unrestricted data | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-021 | AI optional; deterministic core works without it | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-022 | Content digests and crypto-agility; signatures profile-driven | 5 | 5 | 4 | 5 | 5 | 5 | Approve; choose exact profile in design |
| AD-023 | Provider-neutral core with adapter-owned native identifiers | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-024 | Hot summaries plus referenced cold evidence | 5 | 5 | 5 | 5 | 5 | 5 | Approve with retention/integrity rules |
| AD-025 | Versioned offline-capable semantic registries | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-026 | Telemetry is never audit authority | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-027 | Jurisdiction-specific controls supplied as profiles, not core enums | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-028 | Contract-only Phase 2; runtime implementation deferred | 5 | 5 | 5 | 5 | 4 | 5 | Approve |
| AD-029 | Stable envelope plus versioned `DecisionProfile` and typed result | 5 | 5 | 5 | 5 | 5 | 5 | Approve; primary extensibility mechanism |
| AD-030 | Closed primary-form vocabulary with registered secondary facets | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-031 | Authority is orthogonal: authoritative, advisory, recommendation, simulation | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-032 | Authorization is a profile; identity, evaluation, and enforcement remain outside | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-033 | Risk-based durable recording for high-volume authorization | 5 | 5 | 5 | 5 | 4 | 5 | Approve with governed aggregation profile |
| AD-034 | Immutable final records; supersession/revocation are linked effects and projected state | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-035 | Local atomic decision plus audit-obligation acceptance; downstream audit is asynchronous | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-036 | Separate retry idempotency from complete semantic decision identity | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-037 | Deterministic composition replays captured evaluator output without evaluator rerun | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-038 | Synchronous data plane uses local digest-pinned profiles, policy context, schema, strategy, and trust | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-039 | Authorization permits local verified artifacts/caches with bounded freshness and revocation | 5 | 5 | 5 | 5 | 4 | 5 | Approve |
| AD-040 | Composite work is bounded by calls, bytes, concurrency, retries, jurisdictions, time, and cost | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-041 | Cryptographic assurance is tiered; remote notary/HSM is not a universal hot-path dependency | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-042 | Disconnected synchronization uses signed origin-aware manifests and explicit reconciliation | 5 | 5 | 4 | 5 | 5 | 5 | Approve |
| AD-043 | High-risk profiles fail closed; fail-open requires a bounded approved exception | 5 | 5 | 5 | 5 | 5 | 5 | Approve |
| AD-044 | Downstream adoption uses one normative contract, one short per-feature section, and one lightweight gate | 5 | 5 | 5 | 5 | 5 | 5 | Approve; structured manifests remain trigger-based |

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
This version makes the following necessary revisions:

1. Treat sovereignty as data, operational, technical, legal, cryptographic,
   supply-chain, and AI/model control—not residency alone.
2. Add connected, edge, disconnected, and air-gapped deployment profiles.
3. Add local trust, identity, keys, time, registry, audit, and observability
   dependencies with signed synchronization.
4. Add data minimization, legal hold, retention/disposition, lawful redaction,
   cryptographic erasure, and integrity-preserving tombstones.
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
20. Add local authorization, profile-bundle, policy-context, evidence-fan-out,
    cryptographic, and disconnected-synchronization safeguards.

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
2. approval of the seven-value `ScopeKind`/`AuditEvent` compatibility correction
   (ADH-2026-015), permitting all seven canonical values via `metadata.scopeRef`;
3. approval of the durable decision/audit obligation;
4. confirmation of runtime-neutral Phase 2 scope;
5. acceptance or revision of AD-001 through AD-044;
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
13. an architecture clarification/decision record and traceability updates.

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
| Earlier features own shared concepts | PASS WITH CONTROLLED CORRECTION | FEATURE-0012 grammar is extended rather than duplicated; the seven-value `ScopeKind`/`AuditEvent` correction (ADH-2026-015) remains subject to explicit compatibility approval |
| Audit linkage | PASS | Section 10 requires local atomic decision and audit-obligation acceptance and prevents telemetry from becoming audit authority |
| Effective context is the resolution boundary | PASS | Section 7.3 resolves hierarchy once into digest-pinned `EffectivePolicyContext` and prohibits hidden traversal by decision domains |
| Provider-neutral capability matching | PASS | Decision profiles and typed results reference normalized capabilities and typed references rather than provider resources |
| Sovereign and disconnected fit | PASS | Section 11 covers data, operations, technology, law, cryptography, supply chain, AI/model locality, local dependencies, signed synchronization, and portability |
| Stateless horizontal scalability | PASS | Sections 8.5 and 9 require interchangeable workers, local digest-pinned inputs, bounded work, semantic idempotency, backpressure, and no sticky sessions |
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
- strict validation, stable errors, canonical serialization/digest contracts,
  and compatibility checks;
- simple, denial, selection, authorization, composite, replay, correction,
  projection, scope, limit, and extension conformance fixtures;
- in-memory or pure-function reference composition sufficient to prove contract
  semantics, only if required by executable conformance;
- documentation, traceability, Matrix E evidence templates, and feature gates;
- the normative downstream-adoption section and lightweight adoption gate in
  section 28.

Every changed or generated artifact must map to a `CONTRACT_NOW` requirement.
An artifact justified only by a future scenario is prohibited.

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

## 29. ADH-2026-015 scope vocabulary clarification and dependency reconciliation

### 29.1 Controlling authority

ADH-2026-014 and ADH-2026-015 are the **joint controlling handoffs** for
FEATURE-0013. ADH-2026-015 (Approved, 2026-07-27, Sanjeev Kumar) supersedes
**only** the earlier six-scope `AuditEvent` statement in ADH-2026-014. Every
other ADH-2026-014 decision — including the provider-neutral, contract-only,
statelessness, security, sovereignty, anti-overengineering, Matrix E, and
downstream-adoption decisions — remains controlling and unchanged.

Where any earlier text in this document (for example sections 1.1, 7.3, and 15)
refers to "six" `AuditEvent` scopes, this section is authoritative and the
seven-value canonical vocabulary controls.

### 29.2 Approved scope vocabulary decision

1. `ScopeKind` is the single canonical shared vocabulary for scope identity
   across Sovrunn.
2. Its exact canonical values are: `Platform`, `Organization`,
   `OrganizationUnit`, `Tenant`, `Project`, `Provider`, and `ServiceInstance`.
3. `ServiceInstance` is an **additive** extension of the FEATURE-0012 shared
   vocabulary; `Provider` is **retained unchanged**.
4. `AuditEvent` initially permits **all seven** canonical values.
5. `AuditEvent` uses `metadata.scopeRef` as its **sole canonical scope
   identity**.
6. No separate `AuditEvent` scope enum or independently mutable scope field is
   permitted anywhere in schemas, bindings, or validators.
7. Existing FEATURE-0012 serialized `ScopeKind` values and existing
   Organization-scoped `AuditEvent`s remain valid.

### 29.3 metadata.scopeRef is the sole scope authority

For `AuditEvent`, `metadata.scopeRef` is the single authoritative expression of
scope identity. A second `AuditEvent` scope source — a parallel scope enum, a
duplicated scope field, or an independently mutable scope attribute — is
prohibited. Any record that declares scope through more than one source, or
through any source other than `metadata.scopeRef`, must be rejected as a
contradictory or duplicate scope.

### 29.4 Value-by-value FEATURE-0012 versus FEATURE-0013 compatibility

Every inherited controlled-vocabulary value is classified using the required
taxonomy: **added, removed, renamed, mapped, restricted, expanded, or
semantically reinterpreted**.

| Value | FEATURE-0012 `ScopeKind` | FEATURE-0013 (ADH-2026-015) | Classification | Notes |
|---|---|---|---|---|
| `Platform` | Present | Present | Unchanged | Canonical stored form remains absent/nil (FEATURE-0012 D-16). |
| `Organization` | Present | Present | Unchanged | Existing Organization-scoped `AuditEvent`s remain valid. |
| `OrganizationUnit` | Present | Present | Unchanged | — |
| `Tenant` | Present | Present | Unchanged | — |
| `Project` | Present | Present | Unchanged | — |
| `Provider` | Present | Present | Unchanged (retained) | ADH-2026-014's six-scope `AuditEvent` list omitted `Provider`; ADH-2026-015 retains it and permits it for `AuditEvent`. |
| `ServiceInstance` | Absent | Present | **Added** (additive) | Added to the shared vocabulary to reconcile ADH-2026-014's `AuditEvent` requirement with the FEATURE-0012 vocabulary. |

`AuditEvent` profile-level difference:

| Contract element | FEATURE-0012 current | FEATURE-0013 (ADH-2026-015) | Classification |
|---|---|---|---|
| `AuditEvent` permitted scopes (`x-sovrunn-allowed-scopes`) | `["Organization"]` (alpha) | All seven canonical values | **Expanded** |
| `AuditEvent` scope identity field | `metadata.scopeRef` only | `metadata.scopeRef` only (confirmed sole authority) | Unchanged / confirmed |

Overall migration classification: **additive, backward-compatible extension**.
No value is removed, renamed, mapped, restricted, or semantically reinterpreted.

Detailed inspected-file evidence is recorded in
`docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md`.

### 29.5 Mandatory Dependency Contract Reconciliation framework

Whenever FEATURE-0013 (or any feature) inherits a controlled vocabulary, shared
type, or contract from a dependency, a Dependency Contract Reconciliation must
be recorded before design approval. The reconciliation must cover **every** of
the following dimensions explicitly:

1. **Exact dependency version and files** — the dependency's contract version
   and the exact schema, binding, validator, and fixture files inspected.
2. **Enums** — every controlled-vocabulary value, classified per section 29.6.
3. **Required fields** — which fields are required and whether requiredness
   changed.
4. **Cardinality** — how many of each element are permitted; any change.
5. **Ownership** — which feature owns the shared vocabulary or contract.
6. **Serialization** — the impact on existing serialized values and forms.
7. **Validation** — accepted values, rejected values, and rejection rules.
8. **Security semantics** — authorization/scope semantics and whether they
   changed (scope does not grant authorization).
9. **Go bindings** — the impact on generated or hand-written Go types.
10. **Fixtures** — positive and negative fixture obligations.
11. **Migration classification** — additive, backward-compatible, or breaking.
12. **Approval reference** — the controlling ADH/DEC/RFC approval.

A reconciliation that leaves any dimension unaddressed is incomplete and blocks
design approval.

### 29.6 Controlled-vocabulary difference classification requirement

Every inherited controlled-vocabulary difference must be classified as exactly
one of: **added**, **removed**, **renamed**, **mapped**, **restricted**,
**expanded**, or **semantically reinterpreted**. A difference that cannot be
classified with confidence, or that would be `removed`, `renamed`, `mapped`,
`restricted`, or `semantically reinterpreted`, is not permitted without a
separately approved architecture decision.

### 29.7 Prohibition on agent-resolved dependency deltas

Kiro and Cursor must not resolve any unapproved dependency delta between an
inherited contract and the FEATURE-0013 contract. They must not choose
mappings, aliases, removals, additional scopes, restrictions, or applicability
exclusions. Any dependency delta not explicitly approved by ADH-2026-014 or
ADH-2026-015 must be recorded as `ARCHITECTURE_DECISION_REQUIRED` and escalated
to the Architecture Owner. No implementation may proceed from an unapproved
dependency delta.

### 29.8 Synchronization obligation

Canonical schemas, Go bindings, validators, fixtures, compatibility evidence,
architecture, requirements, and design must eventually be synchronized to the
seven-value canonical vocabulary. This synchronization is an obligation, not a
completed state. It must occur only during a separately authorized
implementation stage after requirements receive fresh independent approval.

### 29.9 Preservation of unchanged decisions

ADH-2026-015 changes only the `AuditEvent` scope statement. All provider-
neutral, contract-only, statelessness, security, sovereignty, anti-
overengineering, Matrix E, and downstream-adoption decisions in this document
remain in force. Sections 1 through 28, AD-001 through AD-044, and Matrix E
F13-R01 through F13-R31 are preserved except for the clarified `AuditEvent`
scope vocabulary controlled by this section.
