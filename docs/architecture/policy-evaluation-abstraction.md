---
doc_type: architecture
feature: FEATURE-0017
title: Policy Evaluation Abstraction
status: approved
phase: 2R
baseline: ARCH-2026.08-PHASE2R-CANONICAL
controlling_handoff: ADH-2026-067
clarifying_handoff: ADH-2026-068
corrective_handoff: ADH-2026-069
controlling_payload_sha256: 2da0d06d2e4b0f9c8afb8387b7effc64a79d9d2b022e5bbb570f1ea9795ebf9e
clarifying_payload_sha256: 6c85554764edacf500dc2376ad10ed379d7c7b9873dee86b92d892da066ec3d6
corrective_payload_sha256: 38e038f7de4329b2e5012fefa0f5b092a1d1283d499b59eb37f02049b87b614b
updated: 2026-08-25
ai_load_priority: always
ai_summary: Sole FEATURE-0017 architecture decision authority for deterministic, in-process, engine-neutral policy evaluation with a deterministic fake adapter and pure FEATURE-0013 evidence mapping.
---

# FEATURE-0017 Policy Evaluation Abstraction

## 1. Authority and status

This document is the sole FEATURE-0017 architecture decision authority. It
applies the six-group closure approved through `ADH-2026-067`, the
structural-only mapper correction approved through `ADH-2026-068`, and the
success-timing, exact mapper-field, and invocation-cardinality correction
approved through `ADH-2026-069`. Their exact authority-payload digests are
pinned symmetrically in this document's frontmatter.

No other FEATURE-0017 document may introduce, repeat as independent authority,
or reinterpret the six decision groups in this document. The feature scope
contract, Kiro specifications, task plans, implementation, examples, and
review artifacts must point here and remain subordinate to it.

This document controls FEATURE-0017-specific conflicts; the active architecture
baseline otherwise remains unchanged. Mutable feature-factory stage and gate
status belongs to the FEATURE-0017 automation state, not this architecture.

## 2. Architecture outcome

FEATURE-0017 creates a small executable policy-evaluation seam, not a policy
system or policy engine.

```text
PolicyEvaluationRequest
  -> structural validation and normalization
  -> RFC 8785 JCS v1 input identity
  -> SHA-256 inputDigest
  -> one in-process PolicyEngineAdapter invocation
  -> normalized conclusion or non-result failure
  -> boundary-owned PolicyEvaluationResult construction
  -> transient success return of PolicyEvaluationResult
     together with boundary-captured FEATURE-0013 TimingEnvelope

complete PolicyEvaluationResult
  + exact returned boundary TimingEnvelope
  + caller-supplied FEATURE-0013 adoption metadata
  -> separate pure mapping operation
  -> FEATURE-0013 EvaluationResult
```

Phase 2R supplies only one deterministic fake adapter. FEATURE-0017 does not
inspect business facts to derive a policy conclusion. A later approved real
adapter may evaluate policy and normalize engine-native output behind the same
port, but real-engine selection, integration, and execution are outside this
feature.

## 3. Ownership boundary

### 3.1 FEATURE-0017 owns

- the `PolicyEvaluationRequest` contract realization;
- the `PolicyEvaluationResult` contract realization;
- one engine-neutral `PolicyEngineAdapter` port;
- the evaluation boundary around that port;
- structural request validation and normalization;
- FEATURE-0017 v1 RFC 8785 JCS canonicalization;
- SHA-256 input-digest calculation;
- cancellation and deadline precedence;
- adapter-conclusion and canonical-result validation;
- canonical `PolicyEvaluationResult` construction;
- one deterministic in-process fake adapter;
- normalized result-versus-failure semantics;
- transient transport of successful result plus boundary-captured timing;
- one separate caller-invoked pure mapper to FEATURE-0013
  `EvaluationResult`;
- transient evidence and retention interpretation; and
- FEATURE-0017-local deterministic and zero-external-effect proof.

### 3.2 FEATURE-0017 does not own

- upstream identity, authorization, governance, policy/profile, sovereignty,
  or `EffectiveGovernanceContext` semantics and resolution;
- ExecutionTarget, candidate, placement, entitlement, quota, service
  requirement, or execution-authority semantics;
- authoritative DecisionRecord or AuditEvent publication, persistence, or
  customer-safe and CloudProvider-safe explanation projection;
- production adapter selection, routing, retry, topology, real-engine
  execution, or external integration; or
- a custom policy language or policy engine.

Sections 8 and 12 provide the feature-by-feature ownership map and exhaustive
deferral inventory.

## 4. Component-role map

This section is orientation only. The exact normative contracts are stated once
in CDG-F17-01 through CDG-F17-06.

| Component | Role | Normative authority |
|---|---|---|
| `PolicyEvaluationRequest` | Immutable, Tenant-confidential, internal-engine-facing reference envelope at Project, CloudPlatform, or CloudProvider scope. | CDG-F17-01 and CDG-F17-02 |
| `PolicyEngineAdapter` | Stable in-process, engine-neutral port behind the evaluation boundary. | CDG-F17-03 and CDG-F17-06 |
| `PolicyEvaluationResult` | Immutable, Tenant-confidential, transient evaluator evidence rather than execution authority. | CDG-F17-03 |
| Deterministic fake | Digest-keyed conformance adapter with immutable behavior and no domain-policy interpretation. | CDG-F17-04 |
| FEATURE-0013 mapper | Separate pure structural mapping from successful transient evidence to `EvaluationResult`; authoritative adoption remains downstream. | CDG-F17-05 |

## 5. Consolidated decision register

CDG-F17-01 through CDG-F17-06 are the complete FEATURE-0017 architecture
closure. There is no seventh group and no open decision within these groups.

### CDG-F17-01 — Minimal request boundary

The request contract is:

- `subjectRef`: required structurally valid `TypedRef`;
- `action`: required opaque string, length one to 63;
- `contextRef`: optional; when present, a UID-pinned
  `TypedRef<EffectiveGovernanceContext>`;
- `profileRefs`: required set of one to 32 UID-pinned structurally valid
  references;
- `candidateRefs`: optional set of at most 64 structurally valid references;
  and
- `requestId`: required correlation string, length one to 128.

The evaluation boundary must:

- reject unknown fields and malformed references;
- reject duplicate `profileRefs` and `candidateRefs`;
- treat both reference lists as semantic sets with caller order insignificant;
- perform no existence lookup or later-feature domain validation;
- treat `requestId` as correlation only; and
- invoke no adapter after request validation failure.

No `evaluationClass` is added. FEATURE-0018 may omit context for
context-independent authorization. FEATURE-0023 or another later owner may
require it for a governed call. FEATURE-0020 remains the sole
`EffectiveGovernanceContext` resolver/writer.

### CDG-F17-02 — Exact v1 canonicalization and digest

The digest is FEATURE-0017-specific. It is not a universal DecisionRecord,
evidence, signing, idempotency, or Sovrunn resource digest.

The normalized logical object contains:

- exact `schema` value `sovrunn.policy-evaluation-request/v1`;
- `subjectRef`;
- exact `action`;
- `contextRef` only when present;
- sorted `profileRefs`; and
- `candidateRefs`, always present, with omitted input represented as `[]`.

`requestId` is excluded. The `schema` field is the in-object domain/version
separator; no external prefix or suffix is added.

Every normalized reference contains `apiVersion`, `kind`, and `name`. `uid` is
included only when present and valid under existing reference semantics.
`contextRef` and every `profileRefs` entry require `uid`; `subjectRef` and
`candidateRefs` retain their registered reference semantics.

`profileRefs` and `candidateRefs` sort by the tuple
`(apiVersion, kind, name, uid-or-empty)` using bytewise UTF-8 lexical ordering.
Duplicate identity is equality of normalized reference values. Duplicates are
rejected rather than silently deduplicated.

Absent and empty handling is exact:

- omitted `contextRef` is valid absence and is absent from the normalized
  object;
- `contextRef: null`, `{}`, or a partial reference is invalid;
- omitted `candidateRefs` and `candidateRefs: []` are the same empty set and
  normalize to `"candidateRefs":[]`;
- `candidateRefs: null` is invalid; and
- empty required strings, empty required objects, and an empty `profileRefs`
  set fail validation before digest calculation.

Digest calculation is:

1. serialize the normalized logical object using RFC 8785 JSON
   Canonicalization Scheme;
2. hash exactly those UTF-8 JCS bytes with SHA-256; and
3. encode the hash as lowercase 64-character hexadecimal.

The byte input contains no byte-order mark, trailing newline, external domain
prefix, or suffix. Invalid or duplicate input produces no digest and no
adapter invocation. A future semantic field requires an approved v2 digest
contract rather than silently changing v1 identity.

The required no-newline conformance preimage is:

```json
{"action":"service.read","candidateRefs":[],"profileRefs":[{"apiVersion":"iam.sovrunn.io/v1alpha1","kind":"RoleDefinition","name":"reader","uid":"role-001"}],"schema":"sovrunn.policy-evaluation-request/v1","subjectRef":{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance","name":"reporting-api","uid":"service-001"}}
```

Its required lowercase SHA-256 digest is:

```text
daca15fd0310c46b45d5aff9bbe4f1a5dedd4b788c44c3780838a8be40a56103
```

Conformance must also prove absent and valid-present `contextRef`, reordered
lists producing the same digest, semantic changes producing different digests,
and `requestId`-only changes preserving the digest.

### CDG-F17-03 — Adapter, result, timing, and failure normalization

The evaluation boundary owns:

- request validation and normalization;
- RFC 8785 JCS canonicalization and SHA-256 digest calculation;
- cancellation/deadline precedence;
- invocation-start and valid-completion timing;
- adapter-conclusion validation; and
- complete canonical `PolicyEvaluationResult` construction.

The adapter receives normalized semantic input plus the precomputed digest. The
adapter does not return a complete `PolicyEvaluationResult`; it returns a
normalized conclusion containing only `outcome` and `reasonCodes`, or
`AdapterFailure`. Phase 2R adapter output contains no obligations. The adapter
cannot control canonical result metadata.

The evaluation boundary rejects invalid or duplicate reason codes, lexically
sorts valid codes, and supplies `inputDigest` and `evaluatedAt`.

The only valid result outcomes are:

```text
Allow
Deny
Indeterminate
RequiresApproval
```

The boundary presents the adapter's conclusion; it does not derive a business
outcome itself. Exit-criteria prose `unknown` maps to canonical
`Indeterminate`.

Reason codes contain one to 32 unique lexically sorted values. Each value is
one to 63 characters and matches:

```text
^[A-Z][A-Z0-9_]{0,62}$
```

Phase 2R reason codes are internal conformance evidence and must not contain
raw CloudProvider-native or customer-sensitive detail.

The optional registered `obligations` field is absent in the Phase 2R fake.
Typed obligation semantics and enforcement ownership remain deferred.

The boundary supplies `evaluatedAt` only after valid adapter completion using
an injected deterministic in-process UTC time source. The time source is a
dependency-injection pattern, not a resource, store, network service, or
FEATURE-0017-only platform service. `evaluatedAt` records chronology only;
production time-dependent policy semantics remain deferred.

The normalized non-result categories are exactly:

```text
RequestInvalid
Canceled
DeadlineExceeded
AdapterFailure
InvalidAdapterResult
```

Request validation failure, cancellation, deadline, adapter failure, and
invalid adapter conclusion produce no `PolicyEvaluationResult`. No non-result
interaction may manufacture a valid result. A non-result category is not a
fifth outcome.

### CDG-F17-04 — Deterministic fake

The fake must:

- use immutable configuration for one instance;
- perform exact immutable lookup by approved v1 input digest;
- support configured `Allow`, `Deny`, `Indeterminate`,
  `RequiresApproval`, and `AdapterFailure` behavior;
- return `AdapterFailure`, not `Indeterminate`, for a missing digest;
- fail construction for duplicate digests, malformed digests, or malformed
  configured conclusions;
- contain no time and let the evaluation boundary supply result time; and
- contain no conditional that interprets action, subject, profile, governance
  context, country, CloudProvider, candidate, IAM, sovereignty, entitlement,
  quota, or placement meaning.

The observable fake behavior is independent of fixture serialization and Go
construction mechanics:

| Configured condition | Observable behavior |
|---|---|
| Valid digest mapped to a valid conclusion | Return that configured outcome and reason-code set to the evaluation boundary. |
| Valid digest mapped to adapter failure | Return `AdapterFailure`. |
| No configuration for a valid digest | Return `AdapterFailure`. |
| Duplicate digest configuration | Fake construction fails. |
| Malformed digest or configured conclusion | Fake construction fails. |

The fake does not explain why a real domain policy would reach the conclusion.
Fixture format, constructors, packages, concrete structs, and method signatures
belong to design.

### CDG-F17-05 — Transient FEATURE-0013 evidence linkage

On success, the evaluation operation returns one complete
`PolicyEvaluationResult` together with the exact boundary-captured FEATURE-0013
timing envelope. The pair is transient in-process return metadata with no
canonical type, API identity, schema, resource kind, persistence, public
projection, retention, or independent lifecycle. Design may use a private
receipt, tuple, or equivalent private carrier. On a normalized non-result
failure, the operation returns neither a result nor a success timing envelope.

FEATURE-0013 mapping is a separate caller-invoked pure operation. Mapping is
never automatic and never publishes a DecisionRecord or AuditEvent.

Mapper input is:

- one complete valid `PolicyEvaluationResult`;
- immutable evaluator identity and version;
- the exact FEATURE-0013 timing envelope returned with the successful
  evaluation; and
- structural FEATURE-0013 input identity consisting of a snapshot reference,
  an integrity carrier, or both.

At least one structurally valid FEATURE-0013 input-identity carrier is required.
Malformed structural mapper input causes mapping failure without changing the
valid policy result.

Mapping is exact:

| Policy outcome | FEATURE-0013 `resultStatus` |
|---|---|
| `Allow` | `SUCCESS` |
| `Deny` | `SUCCESS` |
| `RequiresApproval` | `SUCCESS` |
| `Indeterminate` | `INDETERMINATE` |

For all four outcomes, the complete `PolicyEvaluationResult` must be embedded
in `EvaluationResult.result`. `EvaluationResult.evaluatedAt`,
`PolicyEvaluationResult.evaluatedAt`, and the timing envelope's `completedAt`
must be identical.

The evaluation boundary captures `startedAt` immediately before adapter
invocation and `completedAt` after valid adapter output. `evaluatedAt` equals
`completedAt`. The mapper consumes the exact returned timing envelope; it does
not recapture time, use hidden mutable state, read a store, or reconstruct
`startedAt` from `evaluatedAt`.

For every successfully mapped outcome, FEATURE-0017 populates exactly:

- `evaluator` from the caller-supplied immutable identity and version;
- `inputSnapshotRef`, `inputIntegrity`, or both, exactly as supplied after
  structural validation;
- `resultStatus` using the closed mapping above;
- `result` with the complete `PolicyEvaluationResult`;
- `timing.startedAt` and `timing.completedAt` from the returned boundary timing
  envelope; and
- `evaluatedAt`, identical to both `PolicyEvaluationResult.evaluatedAt` and
  `timing.completedAt`.

FEATURE-0017 mapper output always omits `timing.durationMs`, `executionConfig`,
`trustBoundary`, `safetyPolicyFilters`, `structuralTrustState`, and
`scopeEvidence`. The mapper does not inspect a `DecisionProfile` or
conditionally populate those omitted fields. A later adopting domain that
requires an omitted field must use its own approved `EvaluationResult`
construction/adoption path or obtain a separately approved FEATURE-0017 mapper
extension.

The mapper performs no time-source read, reference lookup, profile lookup,
DecisionRecord-scope validation, DecisionRecord write, AuditEvent write,
persistence, or external effect. It neither receives nor resolves a
`DecisionProfile`. The later adopting domain validates the constructed
`EvaluationResult` against its resolved `DecisionProfile` and DecisionRecord
scope before authoritative adoption.

Request validation failure, cancellation, deadline, adapter failure, and
invalid adapter output have no `PolicyEvaluationResult` and are not mapped.
Later material-decision owners may independently record FEATURE-0013 `TIMEOUT`
or `ERROR` evidence under their own profiles.

The `decision-input` retention classification means a valid result may be
retained through a later DecisionRecord when adopted. It does not authorize a
FEATURE-0017 result repository.

### CDG-F17-06 — In-process operation, precedence, proof, and governance

One trusted in-process operation accepts one request. Success returns one valid
result together with its exact boundary timing envelope; one normalized
non-result failure returns neither. The success pairing is private transient
metadata, not a new canonical contract. FEATURE-0017 has no public route,
store, controller, idempotency repository, adapter selector, routing rule,
mutable fixture reload, automatic retry, production engine, or external call.
Production adapter selection is explicitly outside FEATURE-0017.

The adapter is invoked at most once. It is invoked exactly once only when the
already-observed cancellation/deadline check, request validation, digest
calculation, and immediate pre-invocation cancellation/deadline recheck all
succeed. Cancellation or deadline observed before invocation produces zero
adapter invocations.

Late adapter output after cancellation or deadline is discarded. Immutable
fixture configuration makes equivalent concurrent calls independent and
deterministic under the same injected time-source/configuration.

Observable precedence is exact:

```text
1. already-observed cancellation or expired deadline
2. request structural/reference/bound validation
3. canonicalization and SHA-256 digest
4. cancellation/deadline recheck before invocation
5. capture invocation start and perform exactly one fake-adapter invocation
6. cancellation/deadline observed during invocation; discard late output
7. normalized adapter failure
8. normalized adapter-conclusion validation
9. capture completion after valid adapter output, construct the complete
   canonical result with evaluatedAt equal to completion, and return transient
   evidence
```

Design may choose internal mechanics but must not reorder observable
precedence, add retries, or manufacture a result for a non-result interaction.

Only these six decision groups exist. The reuse summary, dependency views,
deferrals, validation contract, and leakage controls below are normative
applications of these groups, not additional decision groups.

## 6. Reuse-before-build disposition

This is the single FEATURE-0017 feature-level reuse summary required by the
Phase 2 reuse-assessment standard.

| Capability | Disposition | Decision status | Rationale | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 `TypedRef`, transient contract, validation, and redaction foundations | Reuse | Approved | Existing canonical foundations satisfy the common contract need. | DEC-0026 and ADH-012 |
| FEATURE-0013 `EvaluationResult` | Reuse | Approved | Decision evidence remains owned by FEATURE-0013. | ADH-017 and DEC-0043 |
| RFC 8785 JCS | Reuse | Approved | A stable standard supplies deterministic JSON canonicalization without a Sovrunn-specific algorithm. | ADH-2026-067 |
| SHA-256 | Reuse | Approved | A standard digest supplies deterministic v1 input identity. | ADH-2026-067 |
| Sovrunn `PolicyEngineAdapter` port | Build | Approved | No external component owns the canonical Sovrunn policy-evaluation contract. | DEC-0028, DEC-0036, ADH-2026-067 |
| Evaluation boundary and pure FEATURE-0013 mapper | Build | Approved | Sovrunn owns validation, normalization, evidence construction, success timing transport, and canonical linkage while engine behavior remains outside core. | ADH-2026-067, ADH-2026-068, and ADH-2026-069 |
| Deterministic in-process fake | Build | Approved | Phase 2R needs executable conformance without selecting or embedding a real policy engine. | Phase 2R scope and ADH-2026-067 |
| Injected deterministic UTC time-source pattern | Reuse | Approved | In-process dependency injection supplies deterministic timing without a new platform service. | ADH-2026-067 |
| Real OPA, Cedar, or other engine adapter | Wrap | Deferred | A later feature requires a fresh candidate assessment; no real engine is selected or allowed here. | DEC-0028 and DEC-0036 |

For the Build units, Reuse is rejected because no external component owns the
Sovrunn canonical seam. Wrap is rejected because no real engine is selected in
FEATURE-0017. Extend is rejected because there is no approved external
implementation to extend. These Build dispositions do not authorize building
a policy engine.

## 7. Accepted-decision reconciliation

ADH-2026-067 is an extension containing corrective reconciliation. It does not
replace DEC-0028 or DEC-0036.

ADH-2026-068 narrows mapper validation to structural inputs and closed mapping
semantics. ADH-2026-069 closes success timing transport, exact mapper-field
population/omission, and invocation cardinality. It narrowly supersedes the
ADH-2026-067 sentence that allowed this mapper to populate other optional
FEATURE-0013 fields when an adopting profile required them; all other
ADH-2026-067 decisions remain intact.

Preserved decisions are:

- all policy logic crosses `PolicyEngineAdapter`;
- Sovrunn does not build a custom policy engine;
- OPA remains the preferred first later real candidate;
- Cedar remains a later authorization-oriented candidate; and
- FEATURE-0017 selects or integrates neither candidate.

Stale DEC-0028 inventory maps as follows:

| Stale term or placeholder | FEATURE-0017 reconciliation |
|---|---|
| `PolicyInput` | No separate canonical type; semantic input is carried by `PolicyEvaluationRequest`. |
| `PolicyContext` | No separate canonical type; optional `contextRef<EffectiveGovernanceContext>` carries the reference. |
| `PolicyBundleRef` | Represented by bounded opaque `profileRefs`. |
| `PolicyDecisionReason` | Represented by normalized `reasonCodes`. |
| OPA/Cedar adapter placeholders | Documentation-only future integration slots; no Go placeholder, dependency, selector, or Phase 2R implementation. |

Boolean `allowed` and retired `EffectivePolicyContext` are not active
FEATURE-0017 vocabulary. Any additional canonical type or real-adapter choice
requires a later approved decision and fresh reuse assessment.

## 8. Earlier and later feature boundaries

| Feature | Relationship | Boundary |
|---|---|---|
| FEATURE-0011 | Reuse-before-build discipline | Controls assessment format and later real-engine reassessment. |
| FEATURE-0012 | Common contract foundation | Reuse `TypedRef`, transient contract, validation, and redaction; do not redefine them. |
| FEATURE-0013 | Direct dependency and decision/audit owner | FEATURE-0017 supplies transient evidence and a pure mapper only. |
| FEATURE-0015 | CloudPlatform and CloudProvider reference foundation | No CloudProvider store, credentials, native identifiers, or API access. |
| FEATURE-0016 | Direct dependency and ExecutionTarget qualification owner | No private-store access, qualification, lifecycle, facts, locks, or mutation. |
| FEATURE-0018 | Later IAM/governance/approval consumer | Owns Principal/IAM mapping, domain requiredness, approvals, and audit behavior. |
| FEATURE-0019 | Later sovereignty policy/evidence owner | Profile, fact, and evidence semantics remain opaque to FEATURE-0017. |
| FEATURE-0020 | EffectiveGovernanceContext resolver/writer | FEATURE-0017 accepts an optional structural reference only. |
| FEATURE-0021 | Enrollment, entitlement, and quota owner | No entitlement or quota policy is hidden in the fake. |
| FEATURE-0022 | Product/runtime/requirement owner | No product or placement-profile interpretation occurs here. |
| FEATURE-0023 | Authoritative sovereignty/placement decision owner | Owns required context, candidate semantics, evidence adoption, composition, and DecisionRecord publication. |
| FEATURE-0024 | Synthetic execution owner | A raw FEATURE-0017 `Allow` never authorizes execution. |
| FEATURE-0025 | Safe explanation owner | Raw inputs, fixture codes, and engine diagnostics are not customer/AI projections. |
| FEATURE-0026 | Slice 0 integration owner | Uses the fake without treating fixture conclusions as production policy. |
| Later production-engine feature | Real adapter and engine-selection owner | Performs a fresh candidate assessment and preserves this port. |

FEATURE-0017 may reuse common public immutable references. It must not import
or call FEATURE-0015/0016 private stores or bridges, hold their locks, obtain
CloudProvider-native facts or credentials, observe or mutate ExecutionTarget
lifecycle, requalify a target, copy backing resolution, or turn an opaque
candidate reference into placement authority.

## 9. Architecture-leak prevention

Architecture fixes data ownership, observable validation, ordering, results,
failures, digest identity, and side-effect boundaries. It does not fix Go
packages, concrete types, method signatures, fixture serialization, goroutine
mechanics, or production topology.

| Stage | May do | Must not do |
|---|---|---|
| Requirements | Translate each decision group into observable requirements, acceptance criteria, and FEATURE-0017-local proof. | Add domain policy, profile taxonomy, routes, stores, engines, downstream semantics, or implementation mechanics. |
| Design | Choose minimal private packages, carriers, interfaces, constructors, and immutable fixture mechanics. | Change request fields, digest identity, precedence, result/failure vocabulary, evidence ownership, or phase scope. |
| Tasks | Decompose approved design into exact FEATURE-0017 paths and proof. | Repair architecture, modify earlier-feature ownership, or add future integrations. |
| Cursor | Implement approved tasks and tests only. | Add policy conditionals, real engines, selectors, persistence, public APIs, external effects, or cross-feature mutation. |

Examples are non-normative. Illustrative action strings, reference kinds,
reason codes, fixture identities, packages, constructors, structs, methods, and
file formats cannot become requirements unless separately approved.

No requirement, design, task, or implementation artifact may reopen a closed
decision or fill a deferred area by inference.

## 10. Required local conformance

FEATURE-0017-local proof must cover:

1. valid request and every registered bound;
2. unknown-field and malformed-reference rejection;
3. invalid request invokes no adapter;
4. optional absent and valid-present `contextRef`;
5. `null`, empty, and partial field rules;
6. duplicate reference rejection;
7. list-order-independent digest;
8. semantic input change alters digest;
9. `requestId`-only change preserves digest;
10. exact RFC 8785/SHA-256 test vector;
11. all four valid outcomes;
12. missing fixture and configured adapter failure;
13. invalid fixture construction failure;
14. cancellation, deadline, zero-invocation-before-call, at-most-once
    invocation, and late-output discard;
15. invalid adapter conclusion rejection;
16. fixed-time-source deterministic replay;
17. reason-code grammar, duplicate rejection, sorting, and redaction;
18. absence of obligations in the Phase 2R fake;
19. successful result/timing transport and separate pure FEATURE-0013 mapping
    for all four results, with the exact populated fields and exact omitted
    fields required by CDG-F17-05;
20. malformed structural mapper-input failure with no mutation or side effect;
21. no mapping for non-result interactions;
22. no FEATURE-0017 DecisionRecord, AuditEvent, route, controller, or store;
23. concurrent equivalent-call determinism; and
24. zero network, OPA, Cedar, Kubernetes, CloudProvider, database, identity,
    provisioning, or other external effects.

## 11. Non-normative examples

### 11.1 Request with governance context

```yaml
subjectRef:
  apiVersion: services.sovrunn.io/v1alpha1
  kind: ServiceInstance
  name: postal-postgresql
  uid: service-instance-fixture-001
action: service-placement.evaluate
contextRef:
  apiVersion: governance.sovrunn.io/v1alpha1
  kind: EffectiveGovernanceContext
  name: fixture-context-a
  uid: fixture-context-a-001
profileRefs:
  - apiVersion: sovereignty.sovrunn.io/v1alpha1
    kind: SovereigntyProfile
    name: fixture-sovereignty-profile-a
    uid: sovereignty-profile-fixture-001
candidateRefs:
  - apiVersion: execution.sovrunn.io/v1alpha1
    kind: ExecutionTarget
    name: fixture-target-a
    uid: fixture-target-a-001
requestId: req-fixture-001
```

The fake validates, canonicalizes, digests, and looks up the input. It does not
infer identity, governance, sovereignty, CloudProvider, or placement meaning
from the names.

### 11.2 Request without governance context

```yaml
subjectRef:
  apiVersion: services.sovrunn.io/v1alpha1
  kind: ServiceInstance
  name: reporting-api
  uid: service-001
action: service.read
profileRefs:
  - apiVersion: iam.sovrunn.io/v1alpha1
    kind: RoleDefinition
    name: reader
    uid: role-001
candidateRefs: []
requestId: req-authz-fixture-001
```

Absence of `contextRef` is valid at the generic seam. It does not assert that a
later governed decision may omit context.

### 11.3 Result

```yaml
outcome: RequiresApproval
reasonCodes:
  - POLICY_EVAL_FIXTURE_APPROVAL_REQUIRED
inputDigest: 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef
evaluatedAt: 2026-08-25T10:00:00Z
```

The digest and timestamp above are illustrative. `RequiresApproval` does not
create an ApprovalRequest. `Allow` does not authorize execution. `Deny` is
evaluator evidence. `Indeterminate` is a valid inconclusive result and remains
distinct from adapter failure.

## 12. Explicit deferrals

The following must not enter FEATURE-0017 requirements, design, tasks, or
implementation:

- PrincipalRef/IAM mapping, membership, roles, authorization fact resolution,
  approval workflow, and exception grants;
- universal principal/action/resource/context authorization modeling;
- fixed policy/profile kind taxonomy;
- policy/profile publication, resolution, bundle distribution, loading, or hot
  reload;
- governance assignment, inheritance, conflict resolution, exception
  resolution, and `EffectiveGovernanceContext` construction or lookup;
- sovereignty facts, evidence, profile interpretation, and authoritative
  sovereignty assessment;
- ExecutionTarget qualification, backing access, lifecycle, normalized target
  facts, and private stores;
- candidate meaning, candidate-set and per-candidate outcomes, filtering,
  ranking, placement, and selection;
- entitlement, quota, service requirements, and execution authority;
- production time-dependent policy input;
- production reason-code or obligation registries;
- customer or CloudProvider reason projection;
- DecisionRecord or AuditEvent publication and downstream authoritative failure
  DecisionRecord policy;
- production adapter selection, registry, routing, or retry;
- real OPA, Cedar, Casbin, OpenFGA, SpiceDB, Cerbos, or other execution;
- in-process versus sidecar versus remote production topology;
- persistent result, audit, policy, or fixture stores; and
- public API routes or controllers.

Before a real adapter is approved, a fresh FEATURE-0011 reuse assessment must
compare current stable candidates, license, security, governance, policy
lifecycle, multi-tenancy, deterministic/sandbox behavior, deployment topology,
performance, explanation, and migration characteristics.
