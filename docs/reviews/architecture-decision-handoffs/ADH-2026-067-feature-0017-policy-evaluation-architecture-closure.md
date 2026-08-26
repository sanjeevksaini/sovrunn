# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-067
- Date: 2026-08-25
- Source discussion: Codex architecture session; consolidated human-approved
  FEATURE-0017 architecture digest
- Related feature: FEATURE-0017
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved
- Authority payload SHA-256:
  `2da0d06d2e4b0f9c8afb8387b7effc64a79d9d2b022e5bbb570f1ea9795ebf9e`
- Authority payload rule: SHA-256 of the exact UTF-8 bytes between the
  `authority-payload:start` and `authority-payload:end` marker lines, excluding
  both marker lines and with the final newline before the end marker included.

<!-- authority-payload:start -->
## Decision title

FEATURE-0017 deterministic policy-evaluation abstraction architecture closure

## Summary

FEATURE-0017 establishes a small, deterministic, in-process policy-evaluation
seam for Phase 2R. It defines the engine-neutral request and result boundary,
one `PolicyEngineAdapter` port, an evaluation boundary that owns validation,
canonicalization, digest and result construction, one deterministic fake, and
a pure mapping into FEATURE-0013 `EvaluationResult`. It does not implement or
select a real policy engine, interpret domain policy, persist evaluation
requests/results, publish DecisionRecords or AuditEvents, or execute any
external operation.

This handoff carries exactly six decision groups. It corrects stale contract
inventory and makes existing accepted adapter-first, decision-first,
audit-first architecture executable without replacing the current baseline.

## Classification

Extension, including corrective reconciliation of stale FEATURE-0017
artifacts. The handoff adds compatible implementation-neutral closure detail
for digest identity, adapter/result ownership, failure precedence, timing,
fake behavior, and FEATURE-0013 mapping while correcting inconsistent contract
inventory. It is not a replacement or new baseline decision and requires no
Architecture Change Request or baseline update.

## Existing approved baseline

- `ARCH-2026.08-PHASE2R-CANONICAL` is the active baseline.
- Phase 2R permits an in-process, deterministic policy-evaluation abstraction
  with DecisionRecord linkage and a bootstrap fake, but no real OPA/Cedar or
  external policy-engine execution.
- DEC-0028 requires all policy logic to cross `PolicyEngineAdapter`, prohibits
  policy rules in handlers and placement logic, prefers OPA as the first later
  real candidate, and permits later Cedar evaluation.
- DEC-0036 requires a Sovrunn adapter boundary before external integration and
  keeps vendor-native types outside core.
- DEC-0043 and FEATURE-0013 own `DecisionRecord`, `DecisionProfile`,
  `EvaluationResult`, and `AuditEvent` semantics.
- DEC-0050 and FEATURE-0020 own immutable
  `EffectiveGovernanceContext` resolution. Authoritative governed decisions
  may require an exact context, but the generic FEATURE-0017 seam must not own
  or resolve that context.
- FEATURE-0016 owns ExecutionTarget qualification, normalized target facts,
  backing access, target lifecycle, and its private stores.
- `VS0-SCHEMA-018` currently makes `contextRef` mandatory and therefore
  conflicts with FEATURE-0018 preceding FEATURE-0020 and with generic
  authorization requests that do not require governance context.
- The draft policy architecture and RFC still contain stale `PolicyInput`,
  `PolicyContext`, `PolicyBundleRef`, `PolicyDecisionReason`, boolean
  `allowed`, and retired `EffectivePolicyContext` vocabulary.

Relevant authority:

- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/decisions/DEC-0028-policy-engine-adapter.md`
- `docs/decisions/DEC-0036-adapter-boundaries.md`
- DEC-0043 and DEC-0050 in `docs/decisions/DECISION_INDEX.md`
- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- `docs/reviews/architecture-readiness/FEATURE-0017-architecture-digest.md`

## Decision or proposed decision

The following six groups are the complete FEATURE-0017 architecture closure.
No other FEATURE-0017 decision group is open or authorized.

### CDG-F17-01 — Minimal request boundary

`PolicyEvaluationRequest` remains an immutable, Tenant-confidential,
internal-engine-facing `TransientRequestResult` valid at Project,
CloudPlatform, or CloudProvider scope.

Its boundary is:

- `subjectRef`: required structurally valid `TypedRef`;
- `action`: required opaque Sovrunn action string, length one to 63;
- `contextRef`: optional; when present, a UID-pinned
  `TypedRef<EffectiveGovernanceContext>`;
- `profileRefs`: required set of one to 32 UID-pinned structurally valid
  references;
- `candidateRefs`: optional set of at most 64 structurally valid references;
  and
- `requestId`: required correlation string, length one to 128.

The boundary must:

- reject unknown fields, malformed references, duplicate `profileRefs`, and
  duplicate `candidateRefs`;
- treat both reference lists as semantic sets with caller order insignificant;
- perform no existence lookup or later-feature domain validation;
- treat `requestId` as correlation metadata, not semantic identity or an
  idempotency key; and
- invoke no adapter after request validation failure.

The optionality of `contextRef` is permanent at the generic seam, not a
temporary exception until FEATURE-0020. FEATURE-0018 may omit it for
context-independent authorization. FEATURE-0023 or another later owning
contract may require it for a governed call. FEATURE-0017 validates and
digests a present reference but never resolves it. No `evaluationClass` is
added.

FEATURE-0017 does not define PrincipalRef, Membership, RoleAssignment,
authorization facts, permitted profile-kind taxonomy, profile publication,
EffectiveGovernanceContext construction, candidate qualification,
per-candidate meaning, candidate selection, or production policy/context
materialization.

### CDG-F17-02 — Exact v1 canonicalization and digest

The digest is FEATURE-0017-specific. It is not a universal DecisionRecord,
evidence, signing, idempotency, or Sovrunn resource digest.

The normalized logical object contains:

- exact `schema` value `sovrunn.policy-evaluation-request/v1`;
- `subjectRef`;
- exact `action`;
- `contextRef` only when present;
- sorted `profileRefs`; and
- `candidateRefs`, always present in the normalized object, with omitted input
  represented as `[]`.

`requestId` is excluded. The `schema` field supplies in-object domain/version
separation; no external prefix or suffix is used.

Every normalized reference contains `apiVersion`, `kind`, and `name`. `uid` is
included only when present and valid under the reference's existing registry
semantics. `contextRef` and every `profileRefs` entry require `uid`;
`subjectRef` and `candidateRefs` preserve their registered reference
semantics.

Reference lists sort by `(apiVersion, kind, name, uid-or-empty)` using bytewise
UTF-8 lexical ordering. Duplicate identity is equality of normalized reference
values; duplicates are rejected and never silently deduplicated.

Absent and empty handling is exact:

- omitted `contextRef` is valid absence and the field is absent from the
  normalized object;
- `contextRef: null`, `{}`, or a partial reference is invalid;
- omitted `candidateRefs` and `candidateRefs: []` are the same empty set and
  normalize to `"candidateRefs":[]`;
- `candidateRefs: null` is invalid; and
- empty required strings, empty required objects, and an empty `profileRefs`
  set fail validation before digest calculation.

Digest calculation is:

1. serialize the normalized logical object with RFC 8785 JSON Canonicalization
   Scheme;
2. hash exactly those UTF-8 JCS bytes with SHA-256; and
3. encode as lowercase 64-character hexadecimal.

The byte input contains no byte-order mark, trailing newline, external domain
prefix, or suffix. Invalid or duplicate input produces no digest and no
adapter invocation. A future semantic field requires an approved v2 digest
contract rather than silently changing v1 identity.

The required no-newline v1 conformance preimage is:

```json
{"action":"service.read","candidateRefs":[],"profileRefs":[{"apiVersion":"iam.sovrunn.io/v1alpha1","kind":"RoleDefinition","name":"reader","uid":"role-001"}],"schema":"sovrunn.policy-evaluation-request/v1","subjectRef":{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance","name":"reporting-api","uid":"service-001"}}
```

Its required lowercase SHA-256 digest is:

```text
daca15fd0310c46b45d5aff9bbe4f1a5dedd4b788c44c3780838a8be40a56103
```

Conformance must additionally prove absent and valid-present `contextRef`,
reordered lists producing the same digest, semantic changes producing
different digests, and `requestId`-only changes preserving the digest.

### CDG-F17-03 — Adapter, result, timing, and failure normalization

The FEATURE-0017 evaluation boundary owns:

- request validation and normalization;
- RFC 8785 canonicalization and SHA-256 input-digest calculation;
- cancellation/deadline precedence;
- invocation start and valid-completion timing;
- adapter-conclusion validation;
- canonical `PolicyEvaluationResult` construction.

FEATURE-0017 also owns the separate pure FEATURE-0013 mapper defined by
CDG-F17-05. The evaluation operation does not invoke that mapper automatically.

`PolicyEngineAdapter` receives normalized semantic input plus the precomputed
digest. `requestId` may accompany the invocation as correlation metadata but
must not affect adapter behavior or result identity.

The adapter returns exactly one of:

- a normalized policy conclusion containing only `outcome` and
  `reasonCodes`; or
- `AdapterFailure`.

Phase 2R adapter output contains no obligations. Invocation and conclusion
carriers are design-private implementation details, not canonical types or
resources. The adapter cannot supply or override `inputDigest`, `evaluatedAt`,
FEATURE-0013 timing, DecisionRecord identity, or AuditEvent identity.

The boundary validates the conclusion, rejects invalid or duplicate reason
codes, lexically sorts valid codes, and constructs the complete canonical
result with the precomputed digest and boundary-owned time.

`PolicyEvaluationResult` has exactly four normalized outcomes:

- `Allow`;
- `Deny`;
- `Indeterminate`; and
- `RequiresApproval`.

The boundary presents the adapter's configured/evaluated conclusion; it does
not derive an outcome from business facts. `unknown` maps to the canonical
`Indeterminate` vocabulary. `RequiresApproval` is evaluator evidence and does
not create or grant approval. `Allow` is not execution authority.

Reason codes must contain one to 32 unique lexically sorted values. Each value
is one to 63 characters and matches `^[A-Z][A-Z0-9_]{0,62}$`. Phase 2R codes
are internal conformance evidence and must not contain raw CloudProvider-native
or customer-sensitive detail.

The optional registered `obligations` field is absent in the Phase 2R fake.
Typed obligation semantics and enforcement ownership remain deferred.

The evaluation boundary supplies `evaluatedAt` only after valid adapter
completion, using an injected deterministic in-process UTC time source. The
time source is a dependency-injection pattern, not a resource, store, network
service, or FEATURE-0017-only platform service. `evaluatedAt` provides result
chronology only; time-dependent production policy semantics remain deferred.

The normalized non-result categories are exactly:

```text
RequestInvalid
Canceled
DeadlineExceeded
AdapterFailure
InvalidAdapterResult
```

Request validation failure, cancellation, deadline, adapter failure, and an
invalid adapter conclusion produce no `PolicyEvaluationResult`. No non-result
interaction may manufacture `Allow` or another valid result. A non-result
category is not a fifth policy outcome.

### CDG-F17-04 — Deterministic fake

FEATURE-0017 builds one deterministic in-process fake as conformance machinery.
It is not a policy engine and contains no policy intelligence.

The fake must:

- use immutable configuration for one instance;
- look up exact approved v1 input digests only;
- support configured `Allow`, `Deny`, `Indeterminate`,
  `RequiresApproval`, and `AdapterFailure` behavior;
- return `AdapterFailure`, not `Indeterminate`, for a missing digest;
- fail construction for duplicate digests, malformed digests, or malformed
  configured conclusions;
- contain no time, letting the evaluation boundary supply result time; and
- contain no conditional that interprets action, subject, profile, governance
  context, country, CloudProvider, candidate, IAM, sovereignty, entitlement,
  quota, or placement meaning.

Normative observable behavior is independent of fixture serialization or Go
construction mechanics:

| Configured condition | Observable behavior |
|---|---|
| Valid digest mapped to a valid conclusion | Return that configured outcome and reason-code set to the evaluation boundary. |
| Valid digest mapped to adapter failure | Return `AdapterFailure`. |
| No configuration for a valid digest | Return `AdapterFailure`. |
| Duplicate digest configuration | Fake construction fails. |
| Malformed digest or configured conclusion | Fake construction fails. |

The fake does not explain why a real domain policy would reach the configured
conclusion. Fixture file format, constructors, packages, concrete structs, and
method signatures belong to design.

### CDG-F17-05 — Transient FEATURE-0013 evidence linkage

FEATURE-0017 persists neither request nor result, publishes no DecisionRecord
or AuditEvent, and owns no result repository.

The evaluation operation returns `PolicyEvaluationResult` or a normalized
non-result failure. FEATURE-0013 mapping is a separate, caller-invoked pure
operation. Mapping is never automatic and never publishes a DecisionRecord or
AuditEvent.

FEATURE-0017 owns a pure, non-persisting mapper whose inputs are:

- one complete valid `PolicyEvaluationResult`;
- immutable evaluator identity and version;
- the evaluation-boundary-captured timing envelope; and
- valid FEATURE-0013 input-identity adoption metadata under the adopting
  `DecisionProfile`.

One valid FEATURE-0013 input-identity arrangement must be supplied under the
adopting profile: snapshot reference, integrity carrier, or both when that
profile permits both. Malformed or profile-incompatible adoption metadata
makes mapping fail without changing the already-valid policy result.

Mapping is exact:

| Policy outcome | FEATURE-0013 result status |
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
`completedAt`. `durationMs` is absent in Phase 2R. Other optional FEATURE-0013
fields remain absent unless an approved adopting profile later requires them.

The mapper performs no clock read, reference lookup, profile lookup,
DecisionRecord write, AuditEvent write, persistence, or external effect.
Request validation failure, cancellation, deadline, adapter failure, and
invalid adapter output have no `PolicyEvaluationResult` and are not mapped.
Later material-decision owners may independently record FEATURE-0013 `TIMEOUT`
or `ERROR` evidence under their own profiles.

The `decision-input` retention classification means a valid result may be
retained through a later DecisionRecord when adopted. It does not authorize a
FEATURE-0017 result store. The later authoritative domain owns adoption,
DecisionRecord publication, AuditEvent obligations, and retention.

### CDG-F17-06 — In-process operation, precedence, proof, and governance

One trusted in-process operation accepts one request and yields one valid
result or one normalized non-result failure. FEATURE-0017 has no public route,
store, controller, idempotency repository, adapter selector, routing rule,
mutable fixture reload, automatic retry, production engine, or external call.
The selected production adapter is explicitly outside FEATURE-0017.

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
6. cancellation/deadline observed during invocation; discard late result
7. normalized adapter failure
8. normalized adapter-conclusion validation
9. capture completion after valid adapter output, construct the complete
   canonical result with evaluatedAt equal to completion, and return transient
   evidence
```

Design may choose internal mechanics but must not reorder observable
precedence, add retries, or manufacture a result for a non-result interaction.

FEATURE-0017-local conformance must prove:

1. valid request and all registered bounds;
2. malformed request invokes no adapter;
3. optional absent and valid-present `contextRef`;
4. duplicate rejection and list-order-independent digest;
5. semantic input changes alter the digest;
6. `requestId`-only changes preserve the digest;
7. all four valid outcomes;
8. missing fixture and configured adapter failure;
9. cancellation/deadline and late-output discard;
10. malformed adapter conclusion rejection;
11. fixed-time-source deterministic replay;
12. reason-code validation, sorting, and redaction;
13. pure valid-result mapping to FEATURE-0013 evidence;
14. no FEATURE-0017 DecisionRecord, AuditEvent, or store; and
15. zero network, OPA, Cedar, Kubernetes, CloudProvider, database, or identity
    effects.

Only CDG-F17-01 through CDG-F17-06 are FEATURE-0017 decision groups. The reuse
summary, accepted-decision reconciliation, phase/conflict analysis, file
allowlist, and leakage restrictions elsewhere in this handoff are normative
clauses of CDG-F17-06, not additional decision groups.

## Rationale

The closure provides the smallest executable value needed from FEATURE-0017:
later features can call one stable, deterministic seam and consume normalized,
auditable evaluator evidence without embedding policy logic or prematurely
selecting an engine. Exact input identity enables deterministic fake lookup and
replay. Boundary-owned result construction prevents adapters from forging
canonical metadata. Pure FEATURE-0013 mapping preserves decision/audit
ownership. Optional governance context avoids a dependency loop with
FEATURE-0020 while allowing later governed callers to require an exact context.

The architecture remains extensible because a later real engine is wrapped
behind the same port, while engine-native types, topology, and lifecycle stay
outside Sovrunn core.

## Reuse-before-build assessment

Canonical standard:
`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`.

This is the single FEATURE-0017 feature-level reuse summary to carry into the
architecture contract and later Kiro artifacts.

| Capability | Disposition | Decision status | Rationale | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 `TypedRef`, transient contract, validation, and redaction foundations | Reuse | Approved | Existing canonical foundations satisfy the common contract need. | DEC-0026 and ADH-012 |
| FEATURE-0013 `EvaluationResult` | Reuse | Approved | Decision evidence remains owned by FEATURE-0013. | ADH-017 and DEC-0043 |
| RFC 8785 JCS | Reuse | Approved | A stable standard supplies deterministic JSON canonicalization without a Sovrunn-specific algorithm. | ADH-2026-067 after approval |
| SHA-256 | Reuse | Approved | A standard digest supplies deterministic v1 input identity. | ADH-2026-067 after approval |
| Sovrunn `PolicyEngineAdapter` port | Build | Approved | No external component owns the canonical Sovrunn policy-evaluation contract. | DEC-0028, DEC-0036, and ADH-2026-067 after approval |
| Evaluation boundary and pure FEATURE-0013 mapper | Build | Approved | Sovrunn owns validation, normalization, evidence construction, and canonical linkage while engine behavior remains outside core. | ADH-2026-067 after approval |
| Deterministic in-process fake | Build | Approved | Phase 2R needs executable conformance without selecting or embedding a real policy engine. | Phase 2R scope and ADH-2026-067 after approval |
| Injected deterministic UTC time-source pattern | Reuse | Approved | In-process dependency injection supplies deterministic timing without a new platform service. | ADH-2026-067 after approval |
| Real OPA, Cedar, or other engine adapter | Wrap | Deferred | A later feature requires a fresh candidate assessment; no real engine is selected or allowed here. | DEC-0028 and DEC-0036 |

For the Build units, Reuse is rejected because no external component owns the
Sovrunn canonical seam. Wrap is rejected because no real engine is selected in
FEATURE-0017. Extend is rejected because there is no approved external
implementation to extend. These Build dispositions do not authorize building
a policy engine.

Mature real-engine candidates must be reassessed when a later real adapter is
proposed. That assessment must compare current stable options, license,
security, governance, policy lifecycle, multi-tenancy, deterministic/sandbox
behavior, deployment topology, performance, explanation, and migration cost.

## Phase impact

- Current phase allowed? Yes.
- Current phase boundary impact: Phase 2R gains an in-process, deterministic,
  side-effect-free policy-evaluation seam and fake only.
- No real provider, Kubernetes, identity, database, network, OPA, Cedar, or
  other policy-engine execution is introduced.
- No feature-sequence change is required.
- FEATURE-0017 remains after FEATURE-0013 and FEATURE-0016 and before
  FEATURE-0018. Optional `contextRef` prevents FEATURE-0020 from becoming a
  prerequisite for the generic seam.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- DEC-0028 and DEC-0036 are preserved: all policy logic crosses the adapter;
  no custom policy engine is built; OPA remains the preferred first later real
  candidate and Cedar a later authorization-oriented candidate.
- DEC-0043 is preserved: FEATURE-0013 remains the sole owner of common
  decision/audit contracts and FEATURE-0017 only supplies transient evidence
  plus a pure mapper.
- DEC-0050 is preserved: FEATURE-0020 remains sole
  `EffectiveGovernanceContext` resolver/writer, while later authoritative
  governed decisions may require the exact context.
- Correction required: `VS0-SCHEMA-018` must move `contextRef` from required to
  optional while retaining UID pinning when present.
- Clarification required: `VS0-SCHEMA-019` must record exact reason-code and
  lowercase SHA-256 grammar.
- Stale DEC-0028 inventory reconciliation:
  - `PolicyInput` becomes no separate canonical type; semantic input is in
    `PolicyEvaluationRequest`;
  - `PolicyContext` becomes no separate canonical type; optional
    `contextRef<EffectiveGovernanceContext>` carries the reference;
  - `PolicyBundleRef` is represented by `profileRefs`;
  - `PolicyDecisionReason` is represented by `reasonCodes`; and
  - OPA/Cedar placeholders are documentation-only future integration slots,
    not Go types or Phase 2R implementations.
- The draft boolean `allowed` and retired `EffectivePolicyContext` vocabulary
  must be removed from active FEATURE-0017 architecture.
- Resolution required: approved ADH application with DEC-0028, RFC-0025,
  architecture, glossary, registry, and traceability reconciliation. No ACR,
  replacement DEC, baseline update, or sequence amendment is required.

## Required action

- Update the existing policy architecture document into the sole FEATURE-0017
  architecture decision authority containing the exact six-group closure.
- Create the FEATURE-0017 feature scope contract containing scope,
  dependencies, the required reuse summary, acceptance boundary, non-goals,
  and a pointer to the sole architecture authority. It must not repeat,
  reinterpret, or become a parallel authority for the six decision groups.
- Correct DEC-0028 and RFC-0025 stale inventory while preserving their accepted
  adapter-first and deferred-engine direction.
- Correct `VS0-SCHEMA-018` and clarify `VS0-SCHEMA-019` exactly as specified.
- Reconcile glossary and direct FEATURE-0017/DEC-0028/VS-000 traceability.
- Update the architecture digest with the approved handoff identity and Kiro
  application status after application succeeds.
- Run the exact validation commands listed below.

Do not update `CURRENT_ARCHITECTURE_BASELINE.md`, phase scope, feature
sequence, roadmap, canonical semantic model, canonical contract catalog,
Structurizr DSL, requirements, design, tasks, or Go implementation.

## Impacted files

Kiro write allowlist after human approval:

| Path | Required update |
|---|---|
| `docs/architecture/policy-evaluation-abstraction.md` | Replace the stale draft with the sole dedicated FEATURE-0017 architecture authority containing all six decision groups and explicit non-goals. |
| `docs/features/FEATURE-0017-policy-evaluation-abstraction.md` | Create a feature scope contract with scope, dependencies, its required reuse summary, acceptance boundary, non-goals, and a pointer to the sole architecture authority. Do not duplicate or reinterpret the six decision groups. |
| `docs/decisions/DEC-0028-policy-engine-adapter.md` | Reconcile stale type inventory and document the deterministic fake/current versus real-adapter/deferred boundary without changing the accepted engine direction. |
| `docs/rfc/RFC-0025-policy-evaluation-abstraction.md` | Replace the stale draft body with the approved engine-neutral contract and Phase 2R/future-engine boundary; update status consistently with the approved handoff. |
| `docs/glossary.md` | Correct `PolicyEvaluationRequest`, `PolicyEvaluationResult`, and `PolicyEngineAdapter` definitions so the adapter returns conclusion/failure and the boundary constructs the canonical result. |
| `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` | Move UID-pinned `contextRef` to optional in VS0-SCHEMA-018; add exact reason-code and lowercase SHA-256 grammar to VS0-SCHEMA-019; preserve all other registered ownership and bounds. |
| `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md` | Add ADH-2026-067 provenance to VS0-SCHEMA-018/019 and point repository authority to the dedicated FEATURE-0017 architecture. |
| `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md` | Record FEATURE-0017 architecture-ready status and ADH-2026-067 control after successful application; do not claim implementation readiness. |
| `docs/traceability/DECISION_TRACEABILITY_MATRIX.md` | Record the reconciled DEC-0028 architecture/RFC and ADH-2026-067 validation state. |
| `docs/reviews/architecture-readiness/FEATURE-0017-architecture-digest.md` | Record approved handoff identity and canonical-application completion; retain the digest as review history, not parallel authority. |

Kiro read/verify-only boundary:

- `docs/context/ARCHITECTURE_VERSION.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/context/CURRENT_DECISION_SUMMARY.md`
- `docs/context/OPEN_QUESTIONS.md`
- `docs/governance/ARCHITECTURE_CHANGE_CONTROL.md`
- `docs/governance/REVIEW_GATES.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/phase2/PHASE2R_REBASELINE.md`
- `docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md`
- `docs/features/FEATURE_INDEX.md`
- `docs/architecture/canonical/sovrunn-finalized-data-model.md`
- `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`
- `docs/architecture/vertical-slices/VS-000-core-skeleton.md`
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`
- `docs/diagrams/structurizr/workspace.dsl`

No Structurizr update is required because this correction changes no system
boundary, major container, external integration, deployment topology, plugin
plane, or major cross-container dynamic flow. If Kiro determines otherwise, it
must stop for human review rather than expanding the write allowlist.

Exact validation commands after Kiro application are:

```bash
make arch-handoff-check HANDOFF=docs/reviews/architecture-decision-handoffs/ADH-2026-067-feature-0017-policy-evaluation-architecture-closure.md
mkdocs build --strict
make feature-contract-check FEATURE=FEATURE-0017
make vs000-contract-check
make phase2-scope-check FEATURE=FEATURE-0017
```

The phase-scope check currently reports completed FEATURE-0012 Kiro artifacts
that retain historical `ResourcePool` vocabulary. That is a pre-existing
checker-scope defect under DEC-0059, not FEATURE-0017 authority. Kiro must
record `PRE_EXISTING_CHECKER_SCOPE_DEFECT`, must not modify FEATURE-0012
artifacts through this handoff, and must require a separate checker-maintenance
correction that recognizes retained historical artifacts before FEATURE-0017
requirements generation.

## Impacted features

- FEATURE-0012: reused unchanged for `TypedRef`, transient contract,
  validation, and redaction foundations.
- FEATURE-0013: reused unchanged; its mapping target and no-writer boundary are
  clarified, but its contracts and implementation are not modified.
- FEATURE-0015: terminology/reference foundation only; no CloudProvider store,
  credentials, native identifiers, lifecycle, or implementation are touched.
- FEATURE-0016: direct dependency and opaque candidate-reference foundation;
  no private store, qualification, lifecycle, backing, lock, or implementation
  is touched.
- FEATURE-0017: architecture contract and direct traceability are created or
  reconciled. No implementation or Kiro specification is generated in this
  pass.
- FEATURE-0018: later consumer may omit governance context for
  context-independent authorization and owns Principal/IAM/approval semantics.
- FEATURE-0019: later sovereignty profile/fact/evidence owner; its kinds remain
  opaque references to FEATURE-0017.
- FEATURE-0020: remains sole EffectiveGovernanceContext resolver/writer.
- FEATURE-0023: later authoritative decision owner may require context,
  interpret candidates, adopt evidence, and publish DecisionRecords.
- FEATURE-0024 through FEATURE-0026: future consumers/integration only; no
  behavior is pulled forward.

## Acceptance criteria for Kiro update

- [ ] Handoff authority-payload digest matches the human-approved SHA-256.
- [ ] Handoff is validated against the Architecture Operating System before
      edits.
- [ ] One sole dedicated FEATURE-0017 architecture authority contains exactly
      CDG-F17-01 through CDG-F17-06 and no parallel decision inventory.
- [ ] The FEATURE-0017 feature scope contract points to the sole architecture
      authority and does not repeat or reinterpret its six decision groups.
- [ ] `contextRef` is optional at the generic seam and UID-pinned when present;
      no `evaluationClass` is introduced.
- [ ] Exact RFC 8785 JCS/SHA-256 v1 normalization, absence rules, ordering,
      duplicate rejection, byte input, encoding, and verified test vector are
      preserved verbatim.
- [ ] The evaluation boundary and adapter have non-overlapping ownership:
      adapter conclusion/failure only; boundary validation, digest, timing,
      result construction, and FEATURE-0013 mapping.
- [ ] All four outcomes, five non-result categories, reason-code grammar,
      obligations absence, and timestamp ownership are exact.
- [ ] The fake is digest-configured conformance machinery with no domain policy
      intelligence, real engine, selector, route, store, or external effect.
- [ ] FEATURE-0013 mapping embeds the complete result for all four outcomes,
      preserves exact timing equality, maps statuses exactly, and performs no
      lookup, clock read, write, persistence, or external effect.
- [ ] The evaluation operation returns only a `PolicyEvaluationResult` or
      normalized non-result failure; FEATURE-0013 mapping remains a separate,
      caller-invoked pure operation and is never automatic.
- [ ] Cancellation, validation, digest, invocation, failure, result validation,
      timing, and return precedence remain exact.
- [ ] Exactly one feature-level reuse summary uses only controlled Reuse/Wrap/
      Extend/Build dispositions and Approved/Deferred statuses.
- [ ] DEC-0028 and RFC-0025 are reconciled without selecting or implementing a
      real engine and without weakening the adapter boundary.
- [ ] `VS0-SCHEMA-018/019`, glossary, architecture, feature authority, and
      traceability agree.
- [ ] CloudProvider terminology is used consistently; no generic Provider or
      provider-native type enters the active contract.
- [ ] No FEATURE-0013, FEATURE-0015, or FEATURE-0016 implementation/ownership is
      reopened.
- [ ] No requirements, design, tasks, Go code, baseline, sequence, roadmap,
      canonical model/catalog, or Structurizr file is modified.
- [ ] The exact validation command set in this handoff is run and its results
      are recorded. The known phase-scope checker defect is not hidden or
      repaired by editing retained FEATURE-0012 artifacts.

## Explicit instructions to Kiro

- Treat this handoff—not chat history or examples—as the complete update
  authority after human approval.
- Keep `docs/architecture/policy-evaluation-abstraction.md` as the sole
  FEATURE-0017 architecture decision authority. The feature scope contract may
  summarize scope and reuse but must not duplicate or reinterpret the six
  decision groups.
- Apply exactly CDG-F17-01 through CDG-F17-06. Do not add an architecture
  choice, profile taxonomy, subject model, candidate semantics, obligation
  model, persistence model, route, engine selector, retry policy, or real
  engine.
- Keep the architecture deterministic, in-process, Phase 2R-scoped,
  audit-first, engine-neutral, and without real external policy-engine
  execution.
- Use CloudProvider terminology consistently.
- Treat examples as non-normative. Do not turn illustrative action strings,
  reference kinds, reason codes, fixture identities, packages, constructors,
  structs, methods, or file formats into architecture requirements.
- Do not introduce engine-native OPA, Rego, Cedar, Kubernetes, identity,
  database, network, or CloudProvider SDK types.
- Do not generate or update `.kiro` requirements, design, or tasks in this
  architecture pass.
- Do not modify Go source, tests, generated prompts, or implementation control
  files.
- Do not update the active architecture baseline, phase sequence, roadmap,
  canonical model/catalog, or Structurizr DSL.
- Stop with `ARCHITECTURE_DECISION_REQUIRED` if applying the six groups needs
  any observable behavior not stated here.
- Stop with `AUTHORITY_CONFLICT` if a controlling accepted source cannot be
  reconciled as a correction without a replacement/new decision path.
<!-- authority-payload:end -->

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-25
- Notes: ADH-2026-067 and authority-payload SHA-256
  `2da0d06d2e4b0f9c8afb8387b7effc64a79d9d2b022e5bbb570f1ea9795ebf9e`
  are approved for Kiro architecture-only validation and canonical repository
  update. This approval does not authorize requirements, design, tasks, or
  implementation in the same pass.

This Architecture Decision Handoff is ready for Kiro validation and repo update only after human approval.
