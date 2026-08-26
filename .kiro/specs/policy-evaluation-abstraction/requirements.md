# FEATURE-0017 Policy Evaluation Abstraction — Requirements

## 1. Identity and stage

| Field | Value |
|---|---|
| Feature | FEATURE-0017 — Policy Evaluation Abstraction |
| Stage | Requirements |
| Kiro slug | `policy-evaluation-abstraction` |
| Phase / order | Phase 2R / 7 |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Sole architecture authority | `docs/architecture/policy-evaluation-abstraction.md` |
| Controlling handoffs | ADH-2026-067 (closure), ADH-2026-068 (structural mapper), ADH-2026-069 (timing/field/cardinality) |
| Direct dependencies | FEATURE-0013, FEATURE-0016 |
| Reused prior authority | FEATURE-0011, FEATURE-0012, FEATURE-0015 |
| Owned schemas | `VS0-SCHEMA-018` (`PolicyEvaluationRequest`), `VS0-SCHEMA-019` (`PolicyEvaluationResult`) |
| Public routes / persistence / controller / real engine | None |

This document owns observable behavior only: intent, actors, scenarios,
invariants, validation outcomes, security/privacy, compatibility, non-goals,
and acceptance criteria. It does not own packages, files, routes, storage,
algorithms, libraries, internal interfaces, or task decomposition. It
translates the closed architecture (CDG-F17-01 through CDG-F17-06) into
verifiable statements. It creates no independent architecture authority and
reopens no closed decision. Any conflict between this document and the sole
architecture authority stops work with `ARCHITECTURE_DECISION_REQUIRED`.

### 1.1 FEATURE-0017 Reuse Assessment reference

The canonical FEATURE-0017 feature-level Reuse Assessment is owned by
`docs/features/FEATURE-0017-policy-evaluation-abstraction.md` and governed by
FEATURE-0011. This requirements stage consumes that approved assessment by
reference; it does not duplicate, extend, or reinterpret its capability
dispositions.

The approved boundary remains unchanged: reuse FEATURE-0012/0013 foundations,
RFC 8785 JCS, SHA-256, and the injected UTC time-source pattern; build only the
Sovrunn evaluation port, boundary, pure mapper, and deterministic in-process
fake; and defer wrapping any real OPA, Cedar, or other engine to a later
approved feature. This subsection is generation-control metadata, not an
independent requirement or architecture authority.

## 2. Purpose and use cases

FEATURE-0017 supplies a deterministic, in-process, engine-neutral policy
evaluation seam. One trusted in-process operation validates and identifies a
`PolicyEvaluationRequest`, invokes the single injected replaceable
`PolicyEngineAdapter` at most once—and exactly once only after every
pre-invocation check succeeds—normalizes the adapter conclusion or failure, and
returns transient audit-ready evidence without embedding, selecting, or
executing a real policy engine.

Actors (all in-process callers; no external principal or route is defined by
this feature):

| Actor | Interaction |
|---|---|
| In-process caller (later owning domain) | Submits one `PolicyEvaluationRequest` to the evaluation operation and, separately, invokes the pure mapper on a successful result. |
| Evaluation boundary | Owns validation, canonicalization, digest, cancellation/deadline precedence, timing capture, adapter-conclusion validation, and canonical result construction. |
| `PolicyEngineAdapter` port | Receives normalized semantic input plus the precomputed digest; returns only a normalized conclusion (`outcome` + `reasonCodes`) or `AdapterFailure`. |
| Deterministic in-process fake | The sole Phase 2R adapter implementation: digest-keyed, immutable, policy-intelligence-free conformance machinery. |
| Injected UTC time source | Dependency-injection pattern supplying deterministic `evaluatedAt`; not a resource, store, or service. |

Primary use cases:

1. A later owning domain obtains normalized, deterministic evaluator evidence
   from a single stable seam without embedding policy logic (see
   `Non-goals and adjacent-feature exclusions` for the excluded owning
   semantics).
2. A caller deterministically identifies equivalent inputs by the
   FEATURE-0017-specific v1 input digest for fixture lookup and replay.
3. A caller separately and purely maps one successful `PolicyEvaluationResult`
   into a FEATURE-0013 `EvaluationResult`, leaving authoritative adoption,
   profile/scope validation, and publication to the later owning domain.

FEATURE-0017 does not inspect business facts to derive a policy conclusion; the
adapter presents a configured/evaluated conclusion and the boundary presents
it unchanged.

## 3. Terms introduced by this feature

Only the terms below are introduced or realized by FEATURE-0017. Terms owned by
prior features are consumed by exact reference and are never redefined here (see
`Architecture/decision/risk traceability by exact ID`).

| Term | Meaning (owned by FEATURE-0017) |
|---|---|
| `PolicyEvaluationRequest` | `VS0-SCHEMA-018` realization: immutable, Tenant-confidential, internal-engine-facing reference envelope valid at Project, CloudPlatform, or CloudProvider scope. |
| `PolicyEvaluationResult` | `VS0-SCHEMA-019` realization: immutable, Tenant-confidential, transient evaluator evidence — not execution authority. |
| `PolicyEngineAdapter` | One engine-neutral in-process port behind the evaluation boundary that returns only a normalized conclusion or `AdapterFailure`. |
| Evaluation boundary | The FEATURE-0017-owned in-process operation surrounding the port; owner of validation, canonicalization, digest, precedence, timing, conclusion validation, and canonical result construction. |
| Deterministic in-process fake | The Phase 2R digest-keyed adapter implementation with immutable configuration and no domain-policy interpretation. |
| Normalized conclusion | The adapter's returned pair `{outcome, reasonCodes}` with no obligations in the Phase 2R fake. |
| Normalized non-result category | One of `RequestInvalid`, `Canceled`, `DeadlineExceeded`, `AdapterFailure`, `InvalidAdapterResult`; never a fifth outcome and never a `PolicyEvaluationResult`. |
| v1 input digest | The FEATURE-0017-specific lowercase 64-char hex SHA-256 over RFC 8785 JCS bytes of the normalized logical object; not a DecisionRecord, evidence, signing, idempotency, or resource digest. |
| Transient success return pairing | The private in-process return of a complete `PolicyEvaluationResult` together with the exact boundary-captured FEATURE-0013 timing envelope; not a canonical type, schema, resource, or store. |
| Pure FEATURE-0013 mapper | The separate, caller-invoked, non-persisting operation that structurally maps one successful result into a FEATURE-0013 `EvaluationResult`. |

## 4. Normative requirements and acceptance scenarios

Each approved requirement below has exactly one normative detail heading, in
approved order (REQ-F17-01 through REQ-F17-06). The exact observable behavior
remains owned solely by the referenced CDG in the sole architecture authority;
statements below translate that behavior and add no semantics. Every approved
acceptance case (`AC-F17-01` through `AC-F17-24`) is mapped explicitly outside
the `Canonical acceptance ledger` in the per-requirement acceptance-scenario
tables and consolidated in `AC coverage mapping`.

### REQ-F17-01 — Minimal request boundary (CDG-F17-01)

The request contract is: `subjectRef` (required structurally valid `TypedRef`);
`action` (required opaque string, length one to 63); `contextRef` (optional;
when present a UID-pinned `TypedRef<EffectiveGovernanceContext>`); `profileRefs`
(required set of one to 32 UID-pinned structurally valid references);
`candidateRefs` (optional set of at most 64 structurally valid references); and
`requestId` (required correlation string, length one to 128).

Observable obligations of the evaluation boundary:

- reject unknown fields and malformed references;
- reject duplicate `profileRefs` and duplicate `candidateRefs`;
- treat both reference lists as semantic sets with caller order insignificant;
- perform no existence lookup and no later-feature domain validation;
- treat `requestId` as correlation only, not semantic identity or an
  idempotency key;
- add no `evaluationClass`; and
- invoke no adapter after request validation failure.

`contextRef` optionality is permanent at the generic seam; a present reference
is validated and digested but never resolved. `EffectiveGovernanceContext`
resolution ownership is excluded (see
`Non-goals and adjacent-feature exclusions`).

| Acceptance case | Observable expectation |
|---|---|
| AC-F17-01 | A valid request at every registered bound (min/max `action` length, one and 32 `profileRefs`, zero and 64 `candidateRefs`, one and 128 `requestId` length) is accepted through validation. |
| AC-F17-02 | An unknown field or a malformed reference is rejected. |
| AC-F17-03 | A request that fails validation invokes no adapter. |
| AC-F17-04 | Absent `contextRef` and a valid present UID-pinned `contextRef` are both accepted. |

### REQ-F17-02 — Exact v1 canonicalization and digest (CDG-F17-02)

The normalized logical object contains: exact `schema` value
`sovrunn.policy-evaluation-request/v1`; `subjectRef`; exact `action`;
`contextRef` only when present; sorted `profileRefs`; and `candidateRefs`
always present, with omitted input represented as `[]`. `requestId` is
excluded. No external prefix or suffix is added.

Every normalized reference contains `apiVersion`, `kind`, and `name`; `uid` is
included only when present and valid under existing reference semantics.
`contextRef` and every `profileRefs` entry require `uid`; `subjectRef` and
`candidateRefs` retain their registered reference semantics. `profileRefs` and
`candidateRefs` sort by the tuple `(apiVersion, kind, name, uid-or-empty)` using
bytewise UTF-8 lexical ordering. Duplicate identity is equality of normalized
reference values.

Exact absent/empty handling:

- omitted `contextRef` is valid absence and is absent from the normalized
  object;
- `contextRef: null`, `{}`, or a partial reference is invalid;
- omitted `candidateRefs` and `candidateRefs: []` are the same empty set and
  normalize to `"candidateRefs":[]`;
- `candidateRefs: null` is invalid; and
- empty required strings, empty required objects, and an empty `profileRefs`
  set fail validation before digest calculation.

Digest calculation: (1) serialize the normalized object with RFC 8785 JCS;
(2) SHA-256 hash exactly those UTF-8 JCS bytes; (3) encode as lowercase
64-character hexadecimal. The byte input contains no byte-order mark, trailing
newline, or external prefix/suffix. Invalid or duplicate input produces no
digest and no adapter invocation. A future semantic field requires an approved
v2 digest contract; v1 identity is never silently changed.

The required no-newline v1 conformance preimage is exactly:

```json
{"action":"service.read","candidateRefs":[],"profileRefs":[{"apiVersion":"iam.sovrunn.io/v1alpha1","kind":"RoleDefinition","name":"reader","uid":"role-001"}],"schema":"sovrunn.policy-evaluation-request/v1","subjectRef":{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance","name":"reporting-api","uid":"service-001"}}
```

Its required lowercase SHA-256 digest is exactly:

```text
daca15fd0310c46b45d5aff9bbe4f1a5dedd4b788c44c3780838a8be40a56103
```

| Acceptance case | Observable expectation |
|---|---|
| AC-F17-05 | `null`, empty, and partial field rules (including `contextRef: null`/`{}`/partial, `candidateRefs: null`, empty required strings/objects, empty `profileRefs`) fail validation before digest calculation. |
| AC-F17-06 | Duplicate `profileRefs` or duplicate `candidateRefs` are rejected, never silently deduplicated. |
| AC-F17-07 | Reordered `profileRefs`/`candidateRefs` inputs produce the same digest. |
| AC-F17-08 | A semantic input change produces a different digest. |
| AC-F17-09 | A `requestId`-only change preserves the digest. |
| AC-F17-10 | The exact RFC 8785/SHA-256 test vector above reproduces the required lowercase digest. |

### REQ-F17-03 — Adapter, result, timing, and failure normalization (CDG-F17-03)

The evaluation boundary owns request validation and normalization, RFC 8785 JCS
canonicalization and SHA-256 digest, cancellation/deadline precedence,
invocation-start and valid-completion timing capture, adapter-conclusion
validation, and complete canonical `PolicyEvaluationResult` construction. The
adapter receives normalized semantic input plus the precomputed digest and
returns only a normalized conclusion (`outcome` + `reasonCodes`) or
`AdapterFailure`; it cannot supply or override `inputDigest`, `evaluatedAt`,
FEATURE-0013 timing, DecisionRecord identity, or AuditEvent identity.

The only valid outcomes are `Allow`, `Deny`, `Indeterminate`, and
`RequiresApproval`; exit-criteria `unknown` maps to canonical `Indeterminate`.
That mapping reconciles exit-criteria vocabulary only: a literal adapter
outcome outside the closed four-value set is invalid and produces
`InvalidAdapterResult` with no result.
`Allow` is not execution authority; `RequiresApproval` creates no approval;
`Deny` and `Indeterminate` are evaluator evidence, and `Indeterminate` remains
distinct from adapter failure.

Reason codes contain one to 32 unique lexically sorted values, each one to 63
characters matching `^[A-Z][A-Z0-9_]{0,62}$`. The boundary rejects invalid or
duplicate reason codes, lexically sorts valid codes, and supplies `inputDigest`
and `evaluatedAt`. Phase 2R reason codes are internal conformance evidence and
must not contain raw CloudProvider-native or customer-sensitive detail. The
optional `obligations` field is absent in the Phase 2R fake.

`evaluatedAt` is supplied only after valid adapter completion using the injected
deterministic in-process UTC time source and records chronology only. The
normalized non-result categories are exactly `RequestInvalid`, `Canceled`,
`DeadlineExceeded`, `AdapterFailure`, and `InvalidAdapterResult`; none produce a
`PolicyEvaluationResult` and none is a fifth outcome.

| Acceptance case | Observable expectation |
|---|---|
| AC-F17-11 | All four valid outcomes (`Allow`, `Deny`, `Indeterminate`, `RequiresApproval`) are producible and correctly normalized. |
| AC-F17-15 | An invalid adapter conclusion is rejected as `InvalidAdapterResult` with no result produced. |
| AC-F17-17 | Reason-code grammar enforcement, duplicate rejection, lexical sorting, and Phase 2R redaction (no CloudProvider-native/customer-sensitive detail) hold. |
| AC-F17-18 | The Phase 2R fake produces no obligations. |

### REQ-F17-04 — Deterministic fake (CDG-F17-04)

The single Phase 2R fake uses immutable configuration for one instance, performs
exact immutable lookup by approved v1 input digest, supports configured `Allow`,
`Deny`, `Indeterminate`, `RequiresApproval`, and `AdapterFailure` behavior,
returns `AdapterFailure` (not `Indeterminate`) for a missing digest, fails
construction on duplicate digests / malformed digests / malformed configured
conclusions, contains no time, and contains no conditional that interprets
action, subject, profile, governance context, country, CloudProvider,
candidate, IAM, sovereignty, entitlement, quota, or placement meaning.

Observable behavior is independent of fixture serialization and construction
mechanics (design-owned): a valid digest mapped to a valid conclusion returns
that outcome/reason-code set; a valid digest mapped to adapter failure or an
unconfigured valid digest returns `AdapterFailure`; a duplicate/malformed
configuration fails construction.

| Acceptance case | Observable expectation |
|---|---|
| AC-F17-12 | A missing fixture returns `AdapterFailure`, and a configured adapter failure returns `AdapterFailure`. |
| AC-F17-13 | Duplicate-digest, malformed-digest, or malformed-conclusion configuration fails fake construction. |

### REQ-F17-05 — Transient FEATURE-0013 evidence linkage (CDG-F17-05)

On success the evaluation operation returns one complete
`PolicyEvaluationResult` together with the exact boundary-captured FEATURE-0013
timing envelope (`startedAt`, `completedAt`; `durationMs` absent). This pairing
is transient in-process return metadata with no canonical type, API identity,
schema, resource kind, persistence, public projection, retention, or lifecycle.
On any normalized non-result failure the operation returns neither a result nor
a success timing envelope.

FEATURE-0013 mapping is a separate, caller-invoked, pure operation; it is never
automatic and never publishes a DecisionRecord or AuditEvent. Mapper input is:
one complete valid `PolicyEvaluationResult`; immutable evaluator identity and
version; the exact returned FEATURE-0013 timing envelope; and structural
FEATURE-0013 input identity consisting of a snapshot reference, an integrity
carrier, or both. At least one structurally valid input-identity carrier is
required. The mapper does not receive or resolve a `DecisionProfile`, does not
receive DecisionRecord scope, and validates only structural input validity plus
the closed mapping; malformed structural input fails mapping without changing
the valid policy result.

Closed status mapping: `Allow`, `Deny`, and `RequiresApproval` map to `SUCCESS`;
`Indeterminate` maps to `INDETERMINATE`. For all four outcomes the complete
`PolicyEvaluationResult` is embedded in `EvaluationResult.result`, and
`EvaluationResult.evaluatedAt`, `PolicyEvaluationResult.evaluatedAt`, and the
timing envelope's `completedAt` are identical (`evaluatedAt` equals
`completedAt`; `startedAt` is captured immediately before invocation).

For every successfully mapped outcome FEATURE-0017 populates exactly:
`evaluator` (from the caller-supplied immutable identity and version);
`inputSnapshotRef`, `inputIntegrity`, or both, exactly as supplied after
structural validation; `resultStatus` (closed mapping above); `result` (the
complete `PolicyEvaluationResult`); `timing.startedAt` and `timing.completedAt`
(from the returned envelope); and `evaluatedAt` (identical to both
`PolicyEvaluationResult.evaluatedAt` and `timing.completedAt`). The mapper
output always omits `timing.durationMs`, `executionConfig`, `trustBoundary`,
`safetyPolicyFilters`, `structuralTrustState`, and `scopeEvidence`; it never
conditionally populates them. The mapper performs no time-source read,
reference lookup, profile lookup, DecisionRecord-scope validation, DecisionRecord
write, AuditEvent write, persistence, or external effect. The later adopting
domain owns validation against its resolved `DecisionProfile` and DecisionRecord
scope before authoritative adoption. The `decision-input` retention
classification permits retention through a later adopted DecisionRecord only and
authorizes no FEATURE-0017 result repository.

| Acceptance case | Observable expectation |
|---|---|
| AC-F17-19 | Successful result/timing transport and separate pure FEATURE-0013 mapping for all four results, with the exact populated fields and exact omitted fields required by CDG-F17-05, hold. |
| AC-F17-20 | Malformed structural mapper input fails mapping with no mutation or side effect and no change to the valid policy result. |
| AC-F17-21 | Non-result interactions (validation failure, cancellation, deadline, adapter failure, invalid adapter output) produce no mapping. |

### REQ-F17-06 — In-process operation, precedence, proof, and governance (CDG-F17-06)

One trusted in-process operation accepts one request; success returns one valid
result plus its exact boundary timing envelope, and one normalized non-result
failure returns neither. FEATURE-0017 has no public route, store, controller,
idempotency repository, adapter selector, routing rule, mutable fixture reload,
automatic retry, production engine, or external call; production adapter
selection is explicitly excluded.

The adapter is invoked at most once, and exactly once only when the
already-observed cancellation/deadline check, request validation, digest
calculation, and immediate pre-invocation cancellation/deadline recheck all
succeed. Cancellation or deadline observed before invocation produces zero
invocations; late adapter output after cancellation/deadline is discarded.
Immutable fixture configuration makes equivalent concurrent calls independent
and deterministic under the same injected time source and configuration.

Observable precedence is exact and must not be reordered:

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

| Acceptance case | Observable expectation |
|---|---|
| AC-F17-14 | Cancellation, deadline, zero-invocation-before-call, at-most-once invocation, and late-output discard all hold per the exact precedence. |
| AC-F17-16 | A fixed injected time source yields deterministic replay of `evaluatedAt` and results for equivalent inputs. |
| AC-F17-22 | No FEATURE-0017 DecisionRecord, AuditEvent, route, controller, or store exists. |
| AC-F17-23 | Concurrent equivalent calls are deterministic and independent. |
| AC-F17-24 | Zero network, OPA, Cedar, Kubernetes, CloudProvider, database, identity, provisioning, or other external effects occur. |

### AC coverage mapping

Every approved acceptance case is mapped to exactly one primary requirement;
supporting requirements are listed where the architecture spreads the concern
across decision groups. This mapping, together with the `Canonical acceptance
ledger`, is the completeness claim for the 24-case inventory.

| Acceptance case | Primary requirement | Supporting requirement(s) |
|---|---|---|
| AC-F17-01 | REQ-F17-01 | REQ-F17-02 |
| AC-F17-02 | REQ-F17-01 | — |
| AC-F17-03 | REQ-F17-01 | REQ-F17-06 |
| AC-F17-04 | REQ-F17-01 | REQ-F17-02 |
| AC-F17-05 | REQ-F17-02 | REQ-F17-01 |
| AC-F17-06 | REQ-F17-02 | REQ-F17-01 |
| AC-F17-07 | REQ-F17-02 | — |
| AC-F17-08 | REQ-F17-02 | — |
| AC-F17-09 | REQ-F17-02 | — |
| AC-F17-10 | REQ-F17-02 | — |
| AC-F17-11 | REQ-F17-03 | — |
| AC-F17-12 | REQ-F17-04 | REQ-F17-03 |
| AC-F17-13 | REQ-F17-04 | — |
| AC-F17-14 | REQ-F17-06 | REQ-F17-03 |
| AC-F17-15 | REQ-F17-03 | REQ-F17-06 |
| AC-F17-16 | REQ-F17-06 | REQ-F17-03 |
| AC-F17-17 | REQ-F17-03 | — |
| AC-F17-18 | REQ-F17-03 | REQ-F17-04 |
| AC-F17-19 | REQ-F17-05 | REQ-F17-03 |
| AC-F17-20 | REQ-F17-05 | — |
| AC-F17-21 | REQ-F17-05 | REQ-F17-06 |
| AC-F17-22 | REQ-F17-06 | — |
| AC-F17-23 | REQ-F17-06 | — |
| AC-F17-24 | REQ-F17-06 | — |

### Canonical requirement ledger

The rows below are copied exactly and once from Section 5.1 of the FEATURE-0017
feature scope contract. They are subordinate generation-control mirrors of the
six canonical architecture decision groups and create no independent authority.
Their ID-to-meaning mapping is immutable; no row is repurposed, renumbered,
merged, split, omitted, or paraphrased.

| Requirement ID | Architecture group | Exact title |
|---|---|---|
| REQ-F17-01 | CDG-F17-01 | Minimal request boundary |
| REQ-F17-02 | CDG-F17-02 | Exact v1 canonicalization and digest |
| REQ-F17-03 | CDG-F17-03 | Adapter, result, timing, and failure normalization |
| REQ-F17-04 | CDG-F17-04 | Deterministic fake |
| REQ-F17-05 | CDG-F17-05 | Transient FEATURE-0013 evidence linkage |
| REQ-F17-06 | CDG-F17-06 | In-process operation, precedence, proof, and governance |

### Canonical acceptance ledger

The rows below are copied exactly and once from Section 5.2 of the FEATURE-0017
feature scope contract (itself a verbatim mirror of the architecture's 24-item
local conformance inventory in Section 10 of the sole authority). Their order
and text are verbatim; they create no independent acceptance authority and may
not be repurposed, renumbered, merged, split, or paraphrased.

| Acceptance ID | Exact acceptance item |
|---|---|
| AC-F17-01 | valid request and every registered bound; |
| AC-F17-02 | unknown-field and malformed-reference rejection; |
| AC-F17-03 | invalid request invokes no adapter; |
| AC-F17-04 | optional absent and valid-present `contextRef`; |
| AC-F17-05 | `null`, empty, and partial field rules; |
| AC-F17-06 | duplicate reference rejection; |
| AC-F17-07 | list-order-independent digest; |
| AC-F17-08 | semantic input change alters digest; |
| AC-F17-09 | `requestId`-only change preserves digest; |
| AC-F17-10 | exact RFC 8785/SHA-256 test vector; |
| AC-F17-11 | all four valid outcomes; |
| AC-F17-12 | missing fixture and configured adapter failure; |
| AC-F17-13 | invalid fixture construction failure; |
| AC-F17-14 | cancellation, deadline, zero-invocation-before-call, at-most-once invocation, and late-output discard; |
| AC-F17-15 | invalid adapter conclusion rejection; |
| AC-F17-16 | fixed-time-source deterministic replay; |
| AC-F17-17 | reason-code grammar, duplicate rejection, sorting, and redaction; |
| AC-F17-18 | absence of obligations in the Phase 2R fake; |
| AC-F17-19 | successful result/timing transport and separate pure FEATURE-0013 mapping for all four results, with the exact populated fields and exact omitted fields required by CDG-F17-05; |
| AC-F17-20 | malformed structural mapper-input failure with no mutation or side effect; |
| AC-F17-21 | no mapping for non-result interactions; |
| AC-F17-22 | no FEATURE-0017 DecisionRecord, AuditEvent, route, controller, or store; |
| AC-F17-23 | concurrent equivalent-call determinism; and |
| AC-F17-24 | zero network, OPA, Cedar, Kubernetes, CloudProvider, database, identity, provisioning, or other external effects. |

## 5. Security, privacy, compatibility, and operational requirements

Security and privacy:

- `PolicyEvaluationRequest` and `PolicyEvaluationResult` are Tenant-confidential,
  internal-engine-facing, and redacted per `VS0-SCHEMA-018`/`VS0-SCHEMA-019`;
  neither is a public projection.
- Phase 2R reason codes are internal conformance evidence and must not carry
  raw CloudProvider-native or customer-sensitive detail (REQ-F17-03).
- No secret value, credential, CloudProvider-native identifier, or engine
  diagnostic enters any FEATURE-0017 contract; references remain opaque.
- The seam performs no authentication, authorization, or principal resolution
  (excluded; see `Non-goals and adjacent-feature exclusions`). It is a trusted
  in-process operation with no network-exposed endpoint, so no route-level
  access control is introduced or implied.
- `requestId` is correlation metadata only, never an identity or idempotency
  key.

Privacy of identity/scope: `subjectRef`, `profileRefs`, `candidateRefs`, and
`contextRef` are structural references; the seam performs no existence lookup,
no resolution, and no domain interpretation of their meaning.

Compatibility:

- FEATURE-0012 `TypedRef`, transient-contract, structural-validation, and
  redaction foundations are reused by exact reference and never redefined.
- FEATURE-0013 `EvaluationResult`, `TimingEnvelope`, evaluator identity/version,
  and structural input-identity carriers are reused by exact reference for the
  pure mapper. `DecisionProfile`, `DecisionRecord`, and `AuditEvent` remain
  FEATURE-0013-owned downstream authorities; FEATURE-0017 neither receives nor
  constructs them and performs no adoption or publication.
- The v1 digest identity is fixed; a future semantic field requires an approved
  v2 digest contract and never silently changes v1 (REQ-F17-02).
- The `PolicyEngineAdapter` port is stable so a later approved real adapter can
  be wrapped behind it without changing this contract.

Operational:

- Deterministic, in-process, side-effect-free Phase 2R behavior; the injected
  UTC time source is a dependency-injection pattern, not a platform service.
- No persistence, controller, reconciliation loop, retry, or background process.
- Immutable fixture configuration guarantees concurrent-call determinism
  (REQ-F17-06).

## 6. Edge cases

| Edge case | Required observable outcome | Requirement |
|---|---|---|
| Unknown field or malformed reference present | Reject; no adapter invocation. | REQ-F17-01, REQ-F17-06 |
| Duplicate `profileRefs` or `candidateRefs` | Reject; never silently deduplicate. | REQ-F17-01, REQ-F17-02 |
| `contextRef` absent | Valid; absent from normalized object and digest. | REQ-F17-01, REQ-F17-02 |
| `contextRef: null`, `{}`, or partial | Invalid; validation fails before digest. | REQ-F17-02 |
| `candidateRefs` omitted vs `[]` | Same empty set; normalize to `"candidateRefs":[]`. | REQ-F17-02 |
| `candidateRefs: null` | Invalid. | REQ-F17-02 |
| Empty required string/object or empty `profileRefs` | Validation fails before digest. | REQ-F17-02 |
| Reordered reference lists | Identical digest. | REQ-F17-02 |
| `requestId`-only change | Identical digest. | REQ-F17-02 |
| Cancellation/deadline observed before invocation | Zero adapter invocations; normalized `Canceled`/`DeadlineExceeded`; no result. | REQ-F17-06 |
| Cancellation/deadline during invocation | Late output discarded; no result. | REQ-F17-06 |
| Missing/unconfigured digest in fake | `AdapterFailure` (not `Indeterminate`). | REQ-F17-04 |
| Duplicate/malformed fake configuration | Fake construction fails. | REQ-F17-04 |
| Invalid or duplicate reason codes from adapter | `InvalidAdapterResult`; no result. | REQ-F17-03 |
| Exit-criteria vocabulary `unknown` | Reconciled to canonical `Indeterminate`; this is not an adapter-input alias. | REQ-F17-03 |
| Literal adapter outcome outside the four-value closed set | `InvalidAdapterResult`; no result. | REQ-F17-03 |
| Malformed structural mapper input | Mapping fails; valid policy result unchanged; no side effect. | REQ-F17-05 |
| Non-result interaction submitted to mapper | No mapping; no manufactured result. | REQ-F17-05 |

## 7. Non-goals and adjacent-feature exclusions

This section consolidates prohibited/stale and adjacent-feature exclusions.
Authority-required names may also appear in normative requirements, the
verbatim canonical acceptance ledger, and traceability only to state
non-ownership, negative behavior, or downstream ownership; those mentions do
not activate an excluded capability.

FEATURE-0017 non-goals (from the sole authority, Sections 3.2, 7, and 12):

- authentication, PrincipalRef creation, membership, role, authorization
  resolution, approval workflow, and exception grants;
- governance, sovereignty, assignment, inheritance, exception, approval,
  entitlement, quota, placement, ranking, selection, or execution semantics;
- policy/profile publication, profile-kind taxonomy, bundle distribution,
  loading, or hot reload;
- `EffectiveGovernanceContext` construction, resolution, or lookup;
- DecisionRecord or AuditEvent publication/persistence and any request/result
  store;
- customer-safe or CloudProvider-safe reason/explanation projection;
- public API routes, handlers, controllers, idempotency repositories, or
  persistent storage;
- production adapter selection, registry, routing, retry, or topology; and
- real OPA, Cedar, Casbin, OpenFGA, SpiceDB, Cerbos, Kubernetes, identity,
  database, network, CloudProvider SDK, provisioning, or plugin execution.

Prohibited stale vocabulary (must not enter active requirements): boolean
`allowed`, retired `EffectivePolicyContext`, and stale DEC-0028 type
placeholders `PolicyInput`, `PolicyContext`, `PolicyBundleRef`,
`PolicyDecisionReason`, and OPA/Cedar Go placeholders.

Adjacent-feature exclusions (owned elsewhere; FEATURE-0017 must not implement,
import ownership of, or pull forward any of these):

| Feature | Excluded concepts |
|---|---|
| FEATURE-0018 | PrincipalRef creation, Membership, RoleDefinition, RoleAssignment, authorization resolution, ApprovalRequest, ExceptionGrant. |
| FEATURE-0019 | SovereigntyProfile semantics, RegulatoryPolicyBundle semantics, SovereigntyFactSet, EvidenceRecord, governed-dimension interpretation. |
| FEATURE-0020 | `EffectiveGovernanceContext` construction, assignment, inheritance, conflict resolution, exception resolution. |
| FEATURE-0021 | Entitlement, quota, provider-selection intent. |
| FEATURE-0022 | ServicePlacementProfile semantics, service requirements, ServiceRegion. |
| FEATURE-0023 | DecisionProfile adoption compatibility, DecisionRecord publication, candidate-set evaluation, per-candidate outcomes, ranking, placement, selection. |
| FEATURE-0024 | Production adapter selection, plugin execution, real realization, provisioning, execution authority. |
| FEATURE-0025 | Customer-safe explanation, CloudProvider-safe explanation, AI-readable projection. |
| FEATURE-0026 | Cross-feature orchestration, Slice 0 integration ownership. |

Negative acceptance cases enforcing these exclusions: AC-F17-18 (no
obligations), AC-F17-21 (no mapping for non-result), AC-F17-22 (no
DecisionRecord/AuditEvent/route/controller/store), and AC-F17-24 (zero external
effects). Prior-feature contracts (FEATURE-0011/0012/0013/0015/0016) are
consumed by exact reference only; their ownership is never copied, renamed,
specialized, or transferred.

## 8. Architecture/decision/risk traceability by exact ID

Sole architecture authority: `docs/architecture/policy-evaluation-abstraction.md`.

| ID | Type | Relationship to FEATURE-0017 requirements |
|---|---|---|
| ADH-2026-067 | Handoff (closure) | Establishes CDG-F17-01..06; source of all six normative detail headings. |
| ADH-2026-068 | Handoff (correction) | Structural-only mapper; narrows REQ-F17-05 mapper input and validation. |
| ADH-2026-069 | Handoff (correction) | Success timing transport, exact mapper populated/omitted fields, at-most-once invocation cardinality; refines REQ-F17-05 and REQ-F17-06. |
| CDG-F17-01 | Decision group | REQ-F17-01. |
| CDG-F17-02 | Decision group | REQ-F17-02. |
| CDG-F17-03 | Decision group | REQ-F17-03. |
| CDG-F17-04 | Decision group | REQ-F17-04. |
| CDG-F17-05 | Decision group | REQ-F17-05. |
| CDG-F17-06 | Decision group | REQ-F17-06. |
| VS0-SCHEMA-018 | Schema | `PolicyEvaluationRequest` realization (REQ-F17-01, REQ-F17-02). |
| VS0-SCHEMA-019 | Schema | `PolicyEvaluationResult` realization (REQ-F17-03, REQ-F17-05). |
| DEC-0028 | Decision | Policy-engine adapter boundary; all policy logic crosses `PolicyEngineAdapter`; no custom engine; OPA preferred later candidate, Cedar later candidate; reconciled, not reopened. |
| DEC-0036 | Decision | Sovrunn adapter boundary before external integration; vendor-native types stay outside core. |
| DEC-0043 | Decision | FEATURE-0013 owns `EvaluationResult`, `TimingEnvelope`, evaluator identity/version, structural input-identity carriers, `DecisionProfile`, `DecisionRecord`, and `AuditEvent`. FEATURE-0017 reuses only the EvaluationResult-related structures required by its pure mapper; it neither receives nor constructs `DecisionProfile`, `DecisionRecord`, or `AuditEvent`. |
| DEC-0050 | Decision | FEATURE-0020 is sole `EffectiveGovernanceContext` resolver/writer; excluded here. |
| DEC-0026 | Decision | FEATURE-0012 common contract foundations reused unchanged. |
| RFC-0025 | RFC | Policy-evaluation-abstraction RFC; reconciled to the engine-neutral contract; subordinate to the sole architecture authority. |

Risk traceability: the sole architecture authority defines no discrete numbered
risk register for FEATURE-0017. Its risk posture is expressed as the explicit
deferrals (Section 12) and architecture-leak controls (Section 9) of the sole
authority, which map to the `Non-goals and adjacent-feature exclusions` section
and to negative acceptance cases AC-F17-18, AC-F17-21, AC-F17-22, and AC-F17-24.
No new risk ID is invented here.

### Zero-inventory traceability

The following registry and public-error families have an explicit zero
FEATURE-0017 inventory. A zero inventory means that this feature neither owns
nor activates an identifier in that family; it does not erase the
authority-required negative references elsewhere in this document.

| Contract or registry family | FEATURE-0017 inventory | Requirement interpretation |
|---|---|---|
| `VS0-WRITER-*` IDs | None owned or activated. | FEATURE-0017 defines no writer contract, persistence owner, or publication path. |
| `VS0-STATE-*` IDs | None owned or activated. | FEATURE-0017 defines no resource lifecycle or state transition. |
| `VS0-CF-*` IDs | None owned or activated. | `AC-F17-01..24` are feature-local conformance evidence; the exact ledger below therefore has no data row. |
| FEATURE-0012 public Problem codes | None owned, introduced, or activated. | The trusted in-process seam defines no public Problem response, HTTP status mapping, or route error contract. |
| FEATURE-0012 Problem violation IDs or FEATURE-0017-local violation IDs | None owned, introduced, or activated. | FEATURE-0017 introduces no public or registry-addressable violation identifier. |

The five normalized non-result categories—`RequestInvalid`, `Canceled`,
`DeadlineExceeded`, `AdapterFailure`, and `InvalidAdapterResult`—are internal
evaluation-operation categories only. They are not FEATURE-0012 public Problem
codes, Problem violation IDs, HTTP status mappings, or additional policy
outcomes.

### Exact conformance semantics ledger

FEATURE-0017's feature-level acceptance inventory is `AC-F17-01..24` (the
architecture's 24-item local conformance list), not a `VS0-CF-*` range. The
VS-000 contract registry defines no `VS0-CF-*` conformance case owned by
FEATURE-0017. This requirements document therefore references no `VS0-CF-*` ID
as FEATURE-0017 behavior.

| ID | owner | inputs | expectedState | expectedError | expectedSideEffects | gate |
|---|---|---|---|---|---|---|

## 9. Design questions explicitly delegated by architecture

The sole architecture authority (Section 9, "Architecture-leak prevention")
delegates the following, and only the following, to the design stage. These are
implementation-mechanics questions, not semantic or ownership questions:

- concrete Go package layout, private carrier type(s), interface shape,
  constructors, and immutable fixture mechanics;
- the private representation of the transient success return pairing (private
  receipt, tuple, or equivalent private carrier) with no canonical contract;
- fixture file format and serialization; and
- goroutine/concurrency mechanics that preserve the exact observable precedence
  and at-most-once invocation without reordering or adding retries.

No semantic, contract, ownership, digest-identity, precedence, result/failure
vocabulary, or field-population question is delegated; those are closed by
CDG-F17-01..06. Illustrative architecture examples (action strings, reference
kinds, reason codes, fixture identities) are non-normative and must not become
required fixtures or business rules by inference.

## 10. Completeness and unresolved-decision report

Completeness:

- Every approved requirement (REQ-F17-01..06) has exactly one normative detail
  heading in approved order, each mapped to its CDG.
- Every approved acceptance case (AC-F17-01..24) is copied verbatim once in the
  `Canonical acceptance ledger` and mapped to a primary requirement in the
  `AC coverage mapping` and the per-requirement acceptance-scenario tables.
- Owned schemas `VS0-SCHEMA-018` and `VS0-SCHEMA-019` are referenced by exact
  ID with no field, bound, enum, or classification change.
- `Non-goals and adjacent-feature exclusions` consolidates adjacent-feature
  exclusions and stale vocabulary. Authority-required mentions elsewhere are
  limited to non-ownership, negative acceptance behavior, or traceability and
  do not activate those concepts.
- Prior-feature contracts are consumed by exact reference only; no inherited
  contract is duplicated, renamed, specialized, or re-owned.
- No `VS0-CF-*` case is claimed as FEATURE-0017 behavior; the exact conformance
  semantics ledger intentionally contains no data row.

Verification performed: the full file was re-read against the sole architecture
authority and the three controlling handoffs. No design choice, future-feature
semantic, duplicated inherited contract, added route/store/controller/engine, or
unresolved architecture decision was introduced. The six decision groups are
all closed; no clarification was required from any loaded authority.

Unresolved decisions: none. No approved stop condition
(`ARCHITECTURE_DECISION_REQUIRED`, `REQUIREMENT_CLARIFICATION_REQUIRED`,
`BOUNDARY_CHANGE_REQUIRED`, `DEPENDENCY_APPROVAL_REQUIRED`,
`SECURITY_REVIEW_REQUIRED`) was triggered.

STAGE_STATUS: COMPLETE
