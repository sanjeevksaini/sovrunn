---
doc_type: architecture_readiness_digest
feature: FEATURE-0017
title: Policy Evaluation Abstraction
status: superseded-by-canonical-architecture
baseline: ARCH-2026.08-PHASE2R-CANONICAL
prepared: 2026-08-22
updated: 2026-08-25
---

# FEATURE-0017 Architecture Digest

> **Superseded readiness record — not architecture authority.** The sole
> FEATURE-0017 architecture decision authority is
> `docs/architecture/policy-evaluation-abstraction.md`. Do not load this digest
> to generate requirements, design, tasks, or implementation instructions.
> Statements below describe the pre-application readiness history and may use
> status language that is no longer current.

## 1. Purpose and authority boundary

This digest consolidates the approved repository starting position for
FEATURE-0017 and records the human-approved architecture closure organized
into six decision groups. These are the only FEATURE-0017 decision groups.

Superseded planning matrices are historical discussion only. They must not be
included in Kiro context, translated into requirements, or treated as open
decisions.

This digest records the pre-application reasoning that led to the approved
handoff. It is retained for audit history only and grants no authority to
generate requirements, design, tasks, or implementation instructions.

Section 5 is the single consolidated register of all human-approved
FEATURE-0017 closure decisions. Section 3 identifies inherited authority and
the approval event; Sections 2, 4, and 6 through 11 explain, trace, or verify
the register and must not be interpreted as parallel decision inventories.

Controlling authority remains, in order:

1. `ARCH-2026.08-PHASE2R-CANONICAL` and the current architecture baseline;
2. the canonical model and contract catalog;
3. accepted DEC/RFC records;
4. the Phase 2R rebaseline and architecture spine;
5. the VS-000 contract registry and traceability authority; and
6. `docs/architecture/policy-evaluation-abstraction.md`, the approved and sole
   FEATURE-0017 architecture decision authority.

Chat discussion, this digest, draft policy documents, and example handoffs do
not override those sources.

## 2. Executive architecture position

FEATURE-0017 is a small executable policy-evaluation seam, not a policy system.

Its Phase 2R value is:

```text
bounded engine-neutral request
  -> deterministic validation and v1 input identity
  -> one in-process PolicyEngineAdapter invocation
  -> deterministic fake conclusion or normalized failure
  -> evaluation-boundary canonical result construction
  -> transient canonical evaluator evidence
  -> optional later FEATURE-0013 evidence adoption
```

FEATURE-0017 owns only:

- `PolicyEvaluationRequest` and `PolicyEvaluationResult` contract realization;
- one engine-neutral `PolicyEngineAdapter` contract;
- one deterministic in-process fake implementation;
- request/result validation;
- FEATURE-0017-specific canonicalization and SHA-256 input digest;
- normalized result-versus-failure behavior;
- transient evidence linkage to FEATURE-0013; and
- FEATURE-0017-local conformance and zero-external-effect proof.

FEATURE-0017 does not own policy intelligence. It does not authenticate a
principal, resolve IAM, compose governance, interpret a profile, qualify or
select an ExecutionTarget, enforce an obligation, publish an authoritative
decision, select a production engine, or execute OPA, Cedar, Kubernetes,
CloudProvider, database, identity, or network operations.

Completeness means every observable behavior inside this boundary is closed.
It does not mean predicting every future policy input or downstream decision.

## 3. Approved starting position

### 3.1 Approved contracts and invariants

| ID | Approved position |
|---|---|
| AP-F17-01 | FEATURE-0017 is Phase 2R order 7 and directly depends on FEATURE-0013 and FEATURE-0016. The approved sequence remains authoritative until formally amended. |
| AP-F17-02 | FEATURE-0017 owns `PolicyEvaluationRequest`, `PolicyEvaluationResult`, `PolicyEngineAdapter`, DecisionRecord linkage, and a deterministic bootstrap fake only. |
| AP-F17-03 | Registered `PolicyEvaluationRequest` is an immutable, Tenant-confidential, internal-engine-facing `TransientRequestResult` for Project, CloudPlatform, or CloudProvider. It currently requires `subjectRef`, `action`, UID-pinned `contextRef<EffectiveGovernanceContext>`, one to 32 UID-pinned `profileRefs`, and `requestId`; `candidateRefs` is optional with maximum 64; retention is none. |
| AP-F17-04 | Registered `PolicyEvaluationResult` is an immutable, Tenant-confidential, internal-engine-facing `TransientRequestResult` with `Allow`, `Deny`, `Indeterminate`, or `RequiresApproval`; one to 32 reason codes; SHA-256 input digest; RFC3339 evaluation time; optional obligations maximum 32; retention classification `decision-input`. |
| AP-F17-05 | A `TransientRequestResult` has no independent resource lifecycle and must not be persisted merely to satisfy resource grammar. |
| AP-F17-06 | All policy logic crosses `PolicyEngineAdapter`; handlers, registries, governance resolution, and placement logic must not become hidden policy engines. Engine-native types remain outside Sovrunn core. |
| AP-F17-07 | Phase 2R is deterministic, in-process, and side-effect-free. FEATURE-0017 includes no real policy engine or external integration. |
| AP-F17-08 | Phase 2R exit evidence includes allow, deny, unknown-like/indeterminate, and adapter-failure behavior with no embedded custom policy engine. |
| AP-F17-09 | FEATURE-0013 solely owns `DecisionRecord`, `DecisionProfile`, `EvaluationResult`, and `AuditEvent` semantics. |
| AP-F17-10 | FEATURE-0016 solely owns ExecutionTarget qualification, normalized facts, lifecycle, backing access, and private stores. |
| AP-F17-11 | DEC-0028 records OPA and Cedar as deferred future candidates and OPA as the current preferred first candidate; neither real adapter is implemented or selected by FEATURE-0017. |
| AP-F17-12 | New countries, sectors, services, backends, and engines enter through governed data, contracts, plugins, and adapters rather than core conditionals. |

### 3.2 Resolved registry conflict pending canonical application

The current mandatory `contextRef<EffectiveGovernanceContext>` conflicts with
the approved sequence and control flow: FEATURE-0018 authorization precedes
FEATURE-0020 context resolution, and the registered context is Project-scoped
while the policy request also supports CloudPlatform and CloudProvider.

The conflict is resolved by CDG-F17-01 in the consolidated register. Kiro must
apply that approved correction to `VS0-SCHEMA-018` through the Architecture
Decision Handoff. This subsection records why reconciliation is required; it
does not restate or extend the decision.

### 3.3 Final human approval record

On 2026-08-25 the human approver accepted adapter-contract completeness,
digest reproducibility, FEATURE-0013 linkage, reuse-before-build compliance,
accepted-decision reconciliation, and architecture-leak prevention as final
closure corrections. Section 5 incorporates all six corrections into the six
existing decision groups; no additional decision group or parallel decision
list remains.

This human approval authorizes the exact Architecture Decision Handoff only.
Requirements generation remains unauthorized until the handoff receives its
exact authority-digest approval and Kiro completes canonical application.

## 4. Role of each FEATURE-0017 component

### 4.1 PolicyEvaluationRequest

`PolicyEvaluationRequest` is a bounded engine-neutral reference envelope. In
FEATURE-0017, its fields intentionally have only structural semantics:

| Input | FEATURE-0017 meaning | Explicit non-ownership |
|---|---|---|
| `subjectRef` | Required structurally valid `TypedRef` identifying the subject supplied by the trusted caller. | No authentication, PrincipalRef conversion, membership, role, or authorization semantics. FEATURE-0018 owns later IAM mapping. |
| `action` | Required opaque bounded Sovrunn action string and part of fixture/digest identity. | No action registry or action-specific policy branch in FEATURE-0017. |
| `contextRef` | Human-approved optional UID-pinned `TypedRef<EffectiveGovernanceContext>`, pending canonical Kiro application. | No existence resolution, governance composition, inheritance, conflict, or exception logic. FEATURE-0020 owns the context. |
| `profileRefs` | One to 32 UID-pinned opaque policy/profile references treated as an order-insensitive set. | No fixed profile-kind taxonomy, publication lookup, policy loading, or interpretation in FEATURE-0017. Owning later features define semantic validity. |
| `candidateRefs` | Optional bounded opaque references associated with the evaluation request. | No ExecutionTarget lookup, qualification, filtering, per-candidate outcome, ranking, or selection. FEATURE-0023 owns candidate meaning. |
| `requestId` | Trusted caller correlation identity. | Not an idempotency key and excluded from semantic input identity. |

A common UID-pinned reference is structurally shaped as:

```yaml
apiVersion: governance.sovrunn.io/v1alpha1
kind: EffectiveGovernanceContext
name: postal-postgresql-context-001
uid: egc-8f42c19a
```

For isolated FEATURE-0017 conformance, later-owned references are synthetic.
Synthetic reference validity proves only the FEATURE-0017 contract; it does not
claim that FEATURE-0018, FEATURE-0020, or FEATURE-0023 resources already exist.

### 4.2 PolicyEvaluationResult

`PolicyEvaluationResult` is normalized evaluator evidence, not an authoritative
authorization, sovereignty, placement, approval, or provisioning decision.

Valid outcomes are:

| Outcome | Meaning inside FEATURE-0017 |
|---|---|
| `Allow` | The adapter implementation returned configured evidence permitting the evaluated request. |
| `Deny` | The adapter implementation returned configured evidence prohibiting the evaluated request. |
| `Indeterminate` | Evaluation completed but produced no trustworthy allow/deny conclusion. |
| `RequiresApproval` | Evaluation evidence says automated continuation is insufficient; it is not approval. |

In Phase 2R the deterministic fake is the only outcome producer and returns a
preconfigured conclusion to the evaluation boundary; the boundary constructs
the complete canonical result. FEATURE-0017 does not inspect business facts to
derive an outcome. A future approved engine adapter will evaluate policy and
map its native answer to the same conclusion contract, but production adapter
implementation and selection are outside FEATURE-0017.

### 4.3 PolicyEngineAdapter

`PolicyEngineAdapter` is one stable in-process port. Its architecture contract
accepts normalized semantic request input plus a boundary-computed digest and
returns a normalized policy conclusion or `AdapterFailure`. The evaluation
boundary, not the adapter, constructs `PolicyEvaluationResult`. Invocation and
conclusion carriers are design-private; the port exposes no Rego, OPA SDK,
Cedar entity, HTTP, Kubernetes admission, CloudProvider SDK, or other
engine/vendor type.

FEATURE-0017 provides no production adapter selector, registry, routing rule,
or dynamic configuration. The deterministic fake is the only Phase 2R
implementation. A future approved feature may provide a concrete adapter that
implements the same port.

### 4.4 Deterministic fake

The fake is conformance machinery. It performs exact immutable fixture lookup
by the approved canonical input digest. It never evaluates IAM, country,
sovereignty, entitlement, quota, placement, CloudProvider, or customer policy.

### 4.5 DecisionRecord linkage

FEATURE-0017 returns transient evidence suitable for later mapping into
FEATURE-0013 `EvaluationResult`. It publishes no `DecisionRecord` or
`AuditEvent` and owns no repository. The later authoritative decision service
decides whether and how to adopt the evidence under its registered
`DecisionProfile`.

## 5. Consolidated human-approved decision register

This section is the only normative FEATURE-0017 closure register in this
digest. All six groups form one human-approved package and must be transferred
without addition, omission, or reinterpretation into the exact Architecture
Decision Handoff. They become canonical repository authority only after that
handoff is approved by exact digest and applied by Kiro.

| Decision group | Consolidated subject |
|---|---|
| CDG-F17-01 | Minimal request contract, optional governance context, validation ownership, and later-feature boundaries |
| CDG-F17-02 | Exact v1 normalization, RFC 8785 JCS bytes, SHA-256 digest, and conformance vectors |
| CDG-F17-03 | Evaluation-boundary ownership, adapter contract, canonical result construction, timing, outcomes, reason codes, and normalized failures |
| CDG-F17-04 | Deterministic fake construction, digest lookup, configured conclusions/failures, and zero policy intelligence |
| CDG-F17-05 | Pure FEATURE-0013 mapping, timing and result preservation, transient retention, and no-writer boundary |
| CDG-F17-06 | In-process operation, precedence, local proof, reuse disposition, accepted-decision reconciliation, and architecture-leak controls |

There are no other FEATURE-0017 closure decision groups and no open decision
inside these six groups.

### CDG-F17-01 — Minimal request boundary

Decision:

- preserve the registered request fields and bounds except make `contextRef`
  optional;
- treat `subjectRef`, `profileRefs`, and `candidateRefs` as structurally typed,
  opaque references in FEATURE-0017;
- reject unknown fields and malformed references;
- reject duplicate `profileRefs` and duplicate `candidateRefs`;
- treat both reference lists as semantic sets, not caller-significant order;
- perform no existence lookup or later-feature domain validation;
- treat `requestId` as correlation only; and
- invoke no adapter after request validation failure.

Deferred to owning later features:

- PrincipalRef, Membership, RoleAssignment, and authorization facts;
- allowed policy/profile kinds and policy publication semantics;
- EffectiveGovernanceContext creation and requiredness for governed calls;
- candidate qualification, set meaning, per-candidate evaluation, and
  selection; and
- production materialization of referenced policy/context data.

### CDG-F17-02 — FEATURE-0017 v1 canonicalization and digest

Decision:

- the digest is FEATURE-0017-specific, not a universal DecisionRecord,
  evidence, signing, or Sovrunn resource digest;
- the normalized logical object contains the exact field
  `"schema":"sovrunn.policy-evaluation-request/v1"` as its in-object domain
  separator; no prefix or suffix is added outside the JCS document;
- canonical input covers `schema`, `subjectRef`, exact `action`, optional
  `contextRef`, sorted `profileRefs`, and sorted `candidateRefs`;
- `requestId` is excluded;
- each reference contains `apiVersion`, `kind`, and `name`; `uid` is included
  only when present and valid under that reference's existing registry
  semantics;
- `contextRef` and `profileRefs` require `uid`; `subjectRef` and
  `candidateRefs` preserve their registered reference semantics;
- `profileRefs` and `candidateRefs` sort by the tuple
  `(apiVersion, kind, name, uid-or-empty)` using bytewise UTF-8 lexical order;
- duplicate references are detected by equality of their normalized reference
  values and rejected rather than silently deduplicated;
- an omitted optional `contextRef` is valid absence, while `null`, `{}`, or a
  partial reference is invalid and produces no digest;
- omitted `candidateRefs` and an explicit empty array are the same empty set
  and canonicalize as `"candidateRefs":[]`; explicit `null` is invalid;
- required strings and `profileRefs` never have an empty semantic form: empty
  values fail validation before canonicalization;
- RFC 8785 JSON Canonicalization Scheme (JCS) is used for the normalized
  logical v1 value, followed by SHA-256 and lowercase 64-character hexadecimal
  encoding; this is a FEATURE-0017-specific choice and is not inherited from
  FEATURE-0013's deferred general canonicalization work;
- the SHA-256 input is exactly the UTF-8 bytes of the RFC 8785 JCS document,
  with no byte-order mark, trailing newline, external domain prefix, or suffix;
- invalid or duplicate input produces no digest and no adapter invocation; and
- a future semantic field may require an approved v2 digest contract rather
  than changing v1 identity silently.

The following no-newline JCS document is the required v1 conformance vector:

```json
{"action":"service.read","candidateRefs":[],"profileRefs":[{"apiVersion":"iam.sovrunn.io/v1alpha1","kind":"RoleDefinition","name":"reader","uid":"role-001"}],"schema":"sovrunn.policy-evaluation-request/v1","subjectRef":{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance","name":"reporting-api","uid":"service-001"}}
```

Its required lowercase SHA-256 digest is:

```text
daca15fd0310c46b45d5aff9bbe4f1a5dedd4b788c44c3780838a8be40a56103
```

Conformance also covers absent and present `contextRef`, reordered reference
lists producing the same digest, semantic changes producing different digests,
and `requestId`-only changes preserving the digest.

### CDG-F17-03 — Result and failure normalization

Decision:

- the evaluation boundary owns request validation, normalization,
  canonicalization, SHA-256 input-digest calculation, cancellation precedence,
  boundary timing, canonical `PolicyEvaluationResult` construction, and the
  pure FEATURE-0013 mapping;
- the adapter receives normalized semantic input plus the precomputed digest;
  `requestId` may accompany the invocation as correlation metadata but must
  not affect adapter behavior or the result;
- the adapter returns either a normalized policy conclusion containing only
  `outcome` and `reasonCodes`, or `AdapterFailure`; Phase 2R adapter output does
  not contain obligations;
- the boundary validates adapter output, rejects invalid or duplicate reason
  codes, lexically sorts valid reason codes, and supplies `inputDigest` and
  `evaluatedAt` to construct the complete canonical result;
- the adapter cannot supply or override the canonical digest, result time,
  FEATURE-0013 timing, DecisionRecord identity, or AuditEvent identity;
- any invocation and conclusion carriers are design-private implementation
  details, not additional canonical types or resources;
- the only valid results are `Allow`, `Deny`, `Indeterminate`, and
  `RequiresApproval`;
- Phase 2R replaces exit-criteria prose `unknown` with canonical
  `Indeterminate`;
- request validation failure, cancellation, deadline, adapter failure, and
  invalid adapter result are non-result interactions;
- no non-result interaction may manufacture `Allow` or any other valid result;
- reason codes are one to 32 unique, lexically ordered strings, each one to 63
  characters and matching `^[A-Z][A-Z0-9_]{0,62}$`;
- Phase 2R reason codes are internal conformance codes and never raw
  CloudProvider/customer projections;
- optional `obligations` is absent in the Phase 2R fake because no typed
  obligation schema or enforcement owner is approved; and
- the FEATURE-0017 evaluation boundary, not the adapter or fixture, supplies
  `evaluatedAt` after a valid adapter completion from an injected in-process
  deterministic UTC time source;
- the time source is a consumer dependency, not a new resource, network
  service, store, or FEATURE-0017-only platform capability; and
- `evaluatedAt` provides result chronology only; production time-dependent
  policy semantics are deferred.

Minimal normalized non-result categories are:

```text
RequestInvalid
Canceled
DeadlineExceeded
AdapterFailure
InvalidAdapterResult
```

The category is part of the internal interaction contract, not a fifth policy
outcome.

### CDG-F17-04 — Deterministic fake behavior

Decision:

- fixtures are immutable for one fake instance;
- lookup uses exact approved v1 input digest only;
- fixtures locally prove `Allow`, `Deny`, `Indeterminate`,
  `RequiresApproval`, and adapter failure;
- a missing fixture returns `AdapterFailure` with internal misconfiguration
  detail, not `Indeterminate`;
- duplicate or malformed fixture configuration prevents a usable fake from
  being constructed;
- fixtures contain no time; the enclosing FEATURE-0017 evaluation boundary
  supplies result time from its injected deterministic UTC time source;
- fixture reason codes are internal and contain no CloudProvider-native or
  customer-sensitive data; and
- no conditional inspects the semantic meaning of an action, subject, profile,
  governance context, country, CloudProvider, or candidate.

The normative fake behavior is intentionally independent of any fixture file
format or Go construction API:

| Configured condition | Observable behavior |
|---|---|
| Valid digest mapped to a valid configured conclusion | Return that configured outcome and reason-code set to the evaluation boundary. |
| Valid digest mapped to configured adapter failure | Return `AdapterFailure`. |
| No configuration for a valid digest | Return `AdapterFailure`. |
| Duplicate digest configuration | Fake construction fails. |
| Malformed digest or malformed configured conclusion | Fake construction fails. |

The fake does not encode why a domain policy would reach the configured
conclusion. Fixture serialization and construction mechanics belong to design.

### CDG-F17-05 — Transient evidence and FEATURE-0013 linkage

Decision:

- FEATURE-0017 persists neither request nor result;
- FEATURE-0017 publishes no DecisionRecord or AuditEvent;
- FEATURE-0017 owns a pure, non-persisting mapping from a valid result plus the
  required adoption metadata into FEATURE-0013 `EvaluationResult`;
- mapper input is one complete `PolicyEvaluationResult`, immutable evaluator
  identity and version, a boundary-captured timing envelope, and valid
  FEATURE-0013 input-identity adoption metadata under the adopting
  `DecisionProfile`; these inputs are not new canonical resources;
- one valid FEATURE-0013 input-identity arrangement is supplied under the
  adopting profile: snapshot reference, integrity carrier, or both when that
  profile permits both; malformed or profile-incompatible adoption metadata
  makes mapping fail without changing the already-valid policy result;
- `Allow`, `Deny`, and `RequiresApproval` are successful evaluator executions
  and map to `EvaluationResult.resultStatus=SUCCESS` when adopted;
- `Indeterminate` maps to `EvaluationResult.resultStatus=INDETERMINATE` when
  adopted;
- the complete canonical policy result must be embedded in
  `EvaluationResult.result` for all four outcomes;
- `EvaluationResult.evaluatedAt`, `PolicyEvaluationResult.evaluatedAt`, and the
  timing envelope's completion time are identical;
- the evaluation boundary captures invocation start immediately before the
  adapter call and completion after valid adapter output; `evaluatedAt` equals
  that completion time, and `durationMs` remains absent in Phase 2R;
- optional FEATURE-0013 fields remain absent unless an approved adopting
  profile later requires them;
- the mapping performs no clock read, reference lookup, profile lookup,
  DecisionRecord write, AuditEvent write, or other side effect;
- request validation failure, cancellation, deadline, adapter failure, and
  invalid adapter output contain no `PolicyEvaluationResult` and are not
  mapped; later material-decision owners may independently record FEATURE-0013
  `TIMEOUT` or `ERROR` evidence under their own profiles; and
- `decision-input` retention means evidence is retained through a later
  DecisionRecord when adopted, not in a FEATURE-0017 result repository.

Example:

```text
transient PolicyEvaluationResult
  -> later FEATURE-0023 adopts it as EvaluationResult evidence
  -> FEATURE-0023 publishes its authoritative DecisionRecord
  -> DecisionRecord retention applies
  -> no independent PolicyEvaluationResult resource is stored
```

Validation failure normally creates no DecisionRecord. Cancellation before
accepted work normally creates no DecisionRecord. Calling-domain audit and
accepted-operation behavior remain owned by the later caller.

### CDG-F17-06 — In-process boundary, precedence, and local proof

Decision:

- one trusted in-process operation accepts one request and yields one valid
  result or one normalized non-result failure;
- FEATURE-0017 performs no automatic retry;
- no public route, store, controller, idempotency repository, engine selector,
  mutable fixture reload, production engine, or external call exists;
- late adapter results after cancellation/deadline are discarded; and
- immutable fixture state makes equivalent concurrent calls independent and
  deterministic under the same injected time-source/configuration.

Approved observable precedence must be:

```text
1. already-observed cancellation or expired deadline
2. request structural/reference/bound validation
3. canonicalization and SHA-256 digest
4. cancellation/deadline recheck before invocation
5. capture invocation start and perform exactly one fake-adapter invocation
6. cancellation/deadline observed during invocation; discard late result
7. normalized adapter failure
8. normalized adapter-conclusion validation
9. capture completion after valid adapter output, construct the complete
   canonical result with evaluatedAt equal to completion, and return transient
   evidence
```

This order closes observable behavior. Design may choose internal mechanics but
must not reorder precedence or add retries/results.

Minimum FEATURE-0017-local conformance proves:

1. valid request and all registered bounds;
2. malformed request invokes no adapter;
3. optional absent and valid-present `contextRef` behavior;
4. duplicate rejection and list-order-independent digest;
5. semantic input change changes digest;
6. `requestId` change alone preserves digest;
7. all four valid outcomes;
8. missing fixture and configured adapter failure;
9. cancellation/deadline and late-result discard;
10. malformed adapter conclusion rejection;
11. fixed-time-source deterministic replay;
12. reason-code validation and redaction;
13. pure successful-result mapping to FEATURE-0013 evidence;
14. no FEATURE-0017 DecisionRecord/AuditEvent/store; and
15. zero network, OPA, Cedar, Kubernetes, CloudProvider, database, or identity
    effects.

#### Consolidated reuse-before-build disposition

This is the single FEATURE-0017 feature-level reuse summary required by the
Phase 2 reuse-assessment standard.

| Capability | Disposition | Decision status | Rationale | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 `TypedRef`, transient contract, validation, and redaction foundations | Reuse | Approved | Existing canonical foundations satisfy the common contract need. | DEC-0026 and ADH-012 |
| FEATURE-0013 `EvaluationResult` | Reuse | Approved | Decision evidence remains owned by FEATURE-0013. | ADH-017 and DEC-0043 |
| RFC 8785 JCS | Reuse | Approved | A stable standard supplies deterministic JSON canonicalization without a Sovrunn-specific algorithm. | FEATURE-0017 Architecture Decision Handoff (pending identifier and Kiro application) |
| SHA-256 | Reuse | Approved | A standard digest supplies deterministic v1 input identity. | FEATURE-0017 Architecture Decision Handoff (pending identifier and Kiro application) |
| Sovrunn `PolicyEngineAdapter` port | Build | Approved | No external component owns the canonical Sovrunn policy-evaluation contract. | DEC-0028, DEC-0036, and FEATURE-0017 Architecture Decision Handoff (pending identifier and Kiro application) |
| Evaluation boundary and pure FEATURE-0013 mapper | Build | Approved | Sovrunn must own validation, normalization, evidence construction, and canonical linkage while keeping engine behavior outside core. | FEATURE-0017 Architecture Decision Handoff (pending identifier and Kiro application) |
| Deterministic in-process fake | Build | Approved | Phase 2R needs executable conformance without selecting or embedding a real policy engine. | Phase 2R baseline and FEATURE-0017 Architecture Decision Handoff (pending identifier and Kiro application) |
| Injected deterministic UTC time-source pattern | Reuse | Approved | Existing in-process dependency-injection practice satisfies deterministic timing without a new platform service. | FEATURE-0017 Architecture Decision Handoff (pending identifier and Kiro application) |
| Real OPA, Cedar, or other engine adapter | Wrap | Deferred | A real adapter requires a fresh candidate assessment and a later feature; no real engine is selected or allowed in FEATURE-0017. | DEC-0028 and DEC-0036 |

The Build choices are intentionally narrow. Reuse is rejected because no
external component owns Sovrunn's canonical port, boundary, mapper, or Phase
2R fake. Wrap is rejected for those units because no real engine is selected
or allowed in FEATURE-0017. Extend is rejected because there is no approved
external implementation to extend. These Build dispositions do not authorize
building a policy engine.

#### Accepted-decision reconciliation

The FEATURE-0017 handoff is a correction and clarification of stale contract
inventory, not a replacement of DEC-0028 or DEC-0036. It preserves these
accepted decisions:

- all policy logic crosses `PolicyEngineAdapter`;
- Sovrunn does not build a custom policy engine;
- OPA remains the preferred first real candidate and Cedar remains a later
  candidate; and
- FEATURE-0017 selects or integrates neither candidate.

The older DEC-0028 inventory is reconciled as follows:

| Older term or placeholder | Phase 2R reconciliation |
|---|---|
| `PolicyInput` | No separate canonical type; the semantic input is carried by `PolicyEvaluationRequest`. |
| `PolicyContext` | No separate canonical type; the generic boundary carries optional `contextRef<EffectiveGovernanceContext>`. |
| `PolicyBundleRef` | Represented by the bounded opaque `profileRefs` set. |
| `PolicyDecisionReason` | Represented by bounded normalized `reasonCodes`. |
| OPA and Cedar adapter placeholders | Documentation-only future integration slots; no Go placeholder, dependency, selector, or Phase 2R implementation. |

The draft policy architecture's boolean `allowed` and retired
`EffectivePolicyContext` are not active FEATURE-0017 canonical types. The
active vocabulary is `PolicyEvaluationRequest`, `PolicyEvaluationResult`,
`PolicyEngineAdapter`, and the deterministic fake role. Any additional
canonical type requires a later approved decision.

After the exact handoff is approved, Kiro must reconcile DEC-0028, draft
RFC-0025, the draft policy architecture document, `VS0-SCHEMA-018` optional
`contextRef`, `VS0-SCHEMA-019` reason-code grammar, relevant traceability, and
the dedicated FEATURE-0017 architecture authority. This closure requires no
baseline replacement or feature-sequence change.

#### Architecture-leak prevention

Only the six CDG-F17 groups are decision groups. Normative closure statements
use `must`. Examples are explicitly non-normative and cannot introduce action
taxonomies, reference kinds, fixture schemas, or downstream policy meaning.

Architecture fixes data ownership, observable validation, ordering, results,
failures, and side-effect boundaries. Design remains free to choose minimal Go
packages, interfaces, private carriers, constructors, and fixture
serialization without changing those observables. No architecture requirement
prescribes a Go method signature or concrete type name.

| Stage | May do | Must not do |
|---|---|---|
| Architecture | Carry the six observable decision groups, one schema correction, deferrals, proof, and ownership into canonical authority. | Select Go mechanics, real engines, downstream IAM/profile/candidate semantics, or production topology. |
| Requirements | Translate each canonically applied group into observable REQ/AC and FEATURE-0017-local conformance. | Add domain rules, policy kinds, routes, stores, engine choices, or later-feature behavior. |
| Design | Choose minimal internal packages and immutable fixture mechanics implementing the approved behavior. | Change fields, digest identity, precedence, results/failures, evidence ownership, or phase boundary. |
| Tasks | Decompose approved design into exact FEATURE-0017 source/test/schema/check paths. | Repair architecture, modify earlier-feature ownership, add future integrations, or broaden paths silently. |
| Cursor | Implement only approved tasks and tests. | Add policy conditionals, real engines, selectors, persistence, public APIs, or cross-feature mutations. |

The register is human-approved for the exact Architecture Decision Handoff and
pending Kiro canonical application. Its digest and precedence must not be
described as canonical repository authority until that application completes.

## 6. Possible inputs and results in Sovrunn

The examples in this section are illustrative. They do not approve action
strings, subject kinds, profile kinds, fixture identities, or downstream
policies.

### 6.1 Illustrative request

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

The fake validates, canonicalizes, digests, and looks up this exact input. It
does not infer identity, governance, sovereignty, or placement meaning from the
names.

### 6.2 Illustrative result

```yaml
outcome: RequiresApproval
reasonCodes:
  - POLICY_EVAL_FIXTURE_APPROVAL_REQUIRED
inputDigest: <64-character-lowercase-sha256>
evaluatedAt: 2026-08-24T10:00:00Z
```

`RequiresApproval` does not create an ApprovalRequest. `Allow` does not grant
execution authority. `Deny` is evaluator evidence. `Indeterminate` is a valid
inconclusive result and remains distinct from adapter failure.

### 6.3 Illustrative downstream scenarios

| Scenario | FEATURE-0017 contribution | Owning feature responsibility |
|---|---|---|
| IAM authorization | Stable adapter result/failure vocabulary | FEATURE-0018 authenticates, resolves IAM facts, defines its mapping, and owns authorization/audit behavior. |
| Sovereignty evaluation | Engine-neutral request/result and evidence digest | FEATURE-0019 owns profiles/facts/evidence; FEATURE-0023 owns authoritative assessment. |
| Governance-dependent decision | Optional EGC reference transport | FEATURE-0020 creates the immutable context; later caller requires it where appropriate. |
| Candidate/placement evaluation | Bounded opaque candidate transport | FEATURE-0016 owns qualification; FEATURE-0023 owns per-candidate meaning, selection, and DecisionRecord. |
| Privileged action | Possible `RequiresApproval` evidence | FEATURE-0018 owns approval policy, request, workflow, and grant semantics. |
| Adapter unavailable | Normalized non-result failure | Calling domain applies its registered fail-closed and audit/decision policy. |

## 7. Earlier and later feature boundaries

This section is a dependency and ownership view derived from the consolidated
decision register. It introduces no additional FEATURE-0017 decision.

| Feature | Relationship to FEATURE-0017 | Boundary |
|---|---|---|
| FEATURE-0011 | Reuse-before-build and reassessment discipline. | FEATURE-0017 classifies the fake as Build and future engines as deferred Wrap candidates. |
| FEATURE-0012 | Common `TypedRef`, transient profile, validation, and redaction grammar. | FEATURE-0017 reuses; it does not redefine common resource identity. |
| FEATURE-0013 | Direct dependency and sole decision/audit contract owner. | FEATURE-0017 supplies transient evidence only. |
| FEATURE-0015 | Canonical CloudPlatform and CloudProvider terminology/reference foundation. | No CloudProvider store, credential, native identifier, or API access. |
| FEATURE-0016 | Direct dependency and sole ExecutionTarget qualification/backing owner. | No private-store access, qualification, lifecycle, facts, locks, or mutation. |
| FEATURE-0018 | Later IAM/governance/approval consumer. | Owns subject/IAM mapping and domain requiredness; does not move IAM into FEATURE-0017. |
| FEATURE-0019 | Later sovereignty policy/evidence owner. | Owns profile/fact/evidence semantics; FEATURE-0017 sees opaque refs only. |
| FEATURE-0020 | Sole EffectiveGovernanceContext resolver/writer. | FEATURE-0017 sees an optional structural reference only. |
| FEATURE-0021 | Enrollment, entitlement, and quota owner. | No entitlement/quota policy hidden in the fake. |
| FEATURE-0022 | Product/runtime/requirement/profile owner. | No product or placement-profile interpretation in FEATURE-0017. |
| FEATURE-0023 | Authoritative sovereignty/placement decision and candidate owner. | Owns required context, candidate semantics, composition, and DecisionRecord. |
| FEATURE-0024 | Synthetic execution owner. | Raw FEATURE-0017 `Allow` never authorizes execution. |
| FEATURE-0025 | Safe explanation owner. | Raw inputs, fixture codes, and engine diagnostics are not customer/AI projections. |
| FEATURE-0026 | Slice 0 integration owner. | Cross-feature orchestration uses the fake without treating fixtures as production policy. |
| Future production-engine feature | Future real adapter and engine-selection owner. | Must perform a fresh FEATURE-0011 reuse assessment and preserve the FEATURE-0017 port. |

### 7.1 FEATURE-0015/0016 reuse exclusions

FEATURE-0017 may reuse common types and public immutable references. It must not:

- import or call FEATURE-0015/0016 private stores or compatibility bridges;
- hold their locks;
- obtain CloudProvider-native facts or credentials;
- observe or modify ExecutionTarget lifecycle;
- requalify a target;
- copy their backing-resource resolution; or
- turn an opaque candidate reference into placement authority.

## 8. Explicitly deferred architecture

The following do not block the small FEATURE-0017 feature and must not appear in
its requirements, design, tasks, or implementation:

- authentication and `PrincipalRef`/IAM mapping;
- a universal principal/action/resource/context authorization model;
- fixed policy/profile kind taxonomy;
- policy publication, resolution, bundle distribution, or hot reload;
- EffectiveGovernanceContext resolution;
- candidate-set and per-candidate placement semantics;
- production time-dependent policy input;
- production reason-code and obligation registries;
- customer or CloudProvider reason projection;
- authoritative failure DecisionRecord policy for each downstream domain;
- production adapter selection or routing;
- real OPA, Cedar, Casbin, OpenFGA, SpiceDB, Cerbos, or other execution;
- in-process versus sidecar versus remote production topology; and
- persistent result, audit, or policy stores.

OPA and Cedar may be mentioned only as inherited future candidates and explicit
non-goals. Before a real adapter is approved, a new FEATURE-0011 reuse
assessment must compare current stable candidates, security, governance,
license, deployment, deterministic/sandbox behavior, policy lifecycle,
multi-tenancy, performance, explanation, and migration characteristics.

This deferral view is derived from CDG-F17-01 through CDG-F17-06 and introduces
no additional decision.

## 9. FEATURE-0015/0016 learning gates

| Learning | FEATURE-0017 prevention |
|---|---|
| Ambiguous observable behavior leaked into requirements/design. | Use the six human-approved groups after canonical Kiro application; no group may contain `TBD` or a choice delegated to Kiro. |
| Domain ownership migrated into convenience bridges. | Opaque references only; no private store, resolver, lock, writer, or lifecycle access. |
| Downstream proof was counted as local feature proof. | CDG-F17-06 defines an exact FEATURE-0017-local conformance ledger; later VS cases are integration evidence only. |
| Error and timing precedence was closed too late. | CDG-F17-03 and CDG-F17-06 fix valid results, non-result failures, cancellation, late-result discard, and deterministic time. |
| Digest/idempotency concepts were conflated. | CDG-F17-02 excludes `requestId`; no idempotency repository exists. |
| Architecture terms drifted into generated artifacts. | Only registered canonical types and `Indeterminate` are active; retired/draft/engine-native terms are exclusions. |
| Tasks omitted path ownership and encouraged out-of-scope edits. | Exact writable/test/schema/checker path inventory is required before Cursor; F0013/F0015/F0016 remain read-only unless an approved compatibility task says otherwise. |
| Context expansion displaced controlling authority. | Kiro receives a measured control manifest; this digest remains planning input until superseded by approved authority. |

## 10. Pre-requirements approval checklist

- [x] Human approval records all six final closure corrections without adding
      decision groups.
- [ ] One approved Architecture Decision Handoff adopts or amends all six
      closure groups as one coherent package.
- [ ] Human approval records the exact authority digest.
- [ ] The registered request is updated so `contextRef` optionality matches the
      approved closure, or the proposal is explicitly rejected and replaced.
- [ ] `VS0-SCHEMA-018/019`, architecture, glossary, and traceability agree.
- [x] The digest closes the FEATURE-0017-specific v1 canonicalization and
      verifies its exact SHA-256 conformance vector.
- [x] All four outcomes and all non-result categories are unambiguous in the
      approved closure.
- [x] Fake fixture hit, miss, malformed configuration, deterministic time, and
      adapter failure behavior are exact.
- [x] Obligations are explicitly absent in the Phase 2R fake.
- [x] FEATURE-0013 mapping, no-writer boundary, and retention interpretation
      are exact.
- [x] FEATURE-0017-local conformance and zero-external-effect proof are mapped.
- [x] No IAM mapping, profile-kind taxonomy, candidate selection, production
      engine, route, store, controller, or external integration has leaked in.
- [x] OPA/Cedar remain deferred candidate context only; production selection
      requires a fresh reuse assessment.
- [ ] Kiro stage boundaries and the later Cursor writable-path allowlist are
      exact.

## 11. Exit condition

No FEATURE-0017 architecture choice remains open inside the six-group closure.
The remaining work is governance processing: exact handoff generation, exact
authority-digest approval, Kiro canonical reconciliation, and stage-gate
approval.

FEATURE-0017 may enter requirements generation only after the six-group package
is approved through the Architecture Decision Handoff flow and applied by Kiro
to dedicated FEATURE-0017 architecture authority, the contract registry,
traceability, and readiness checks.

Until then, no requirements, design, tasks, real policy engine, or
implementation work is authorized.
