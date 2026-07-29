# Tasks Document

Feature: FEATURE-0013 — Decision Record and AuditEvent Standard
Phase: Phase 2 — Reuse-First PaaS Fabric Foundation
Stage: Tasks

## 1. Metadata

| Field | Value |
|---|---|
| Feature ID | FEATURE-0013 |
| Feature title | Decision Record and AuditEvent Standard |
| Phase | Phase 2 — Reuse-First PaaS Fabric Foundation |
| Spec stage | Tasks |
| Branch | feature-0013-decision-record-and-auditevent-standard |
| Depends on | FEATURE-0012 (API, Resource Naming, Status, and Validation Standard) |
| Controlling handoff | ADH-2026-017 (Approved) |
| Canonical architecture | docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md |
| Controlling design | .kiro/specs/decision-object-and-auditevent-standard/design.md |
| Controlling requirements | .kiro/specs/decision-object-and-auditevent-standard/requirements.md |
| Go module | github.com/sanjeevksaini/sovrunn |
| Go toolchain | Go 1.22; stdlib + `gopkg.in/yaml.v3` only |
| Status | Draft — pending review |

These tasks are generated only from the approved consolidated architecture, the
approved requirements, and the approved design. They do not use
`design.superseded-patched-draft.md` or ADH-2026-014/015/016 as semantic input;
those handoffs are historical provenance only. Only `CONTRACT_NOW` design
components produce implementation tasks. `INVARIANT_FOR_LATER` and `DEFERRED`
rows produce no source task beyond documentation and traceability.

## 2. Reuse summary

Reproduced verbatim from the approved architecture section 14 and requirements;
field meanings are owned by `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`.
Tasks introduce no new disposition.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 API/resource grammar | Extend | Reuse every shared profile, metadata, reference, boundary, validation, Problem Details, schema, and compatibility primitive; add only decision/audit domain semantics and approved per-kind constraints. | Approved | ADH-2026-017; ADH-2026-012; RFC-0022 |
| DecisionRecord, DecisionProfile, EvaluationResult, and AuditEvent domain contract | Build | No mature external model supplies Sovrunn's complete sovereign, provider-neutral authority, projection, scope, audit-obligation, and downstream-adoption semantics. | Approved | ADH-2026-017; DEC-0026; RFC-0023 |
| Bounded decision-requirements and composition concepts | Reuse | Reuse applicable OMG DMN decision-requirements concepts without selecting, embedding, or recreating a DMN engine or workflow language. | Approved | ADH-2026-017; RFC-0023 |
| Audit transport and operational correlation | Reuse | CloudEvents may be an optional transport mapping and OpenTelemetry may carry operational correlation; neither becomes the canonical record or audit authority. | Approved | ADH-2026-017; RFC-0023 |
| Supply-chain evidence references | Reuse | Reuse DSSE/in-toto and SPDX/CycloneDX concepts by typed reference without duplicating their schemas or implementing signing. | Approved | ADH-2026-017; ADR-F13-002 |
| Canonicalization, digest, signing, and cryptographic verification | Reuse | Preserve algorithm-agile carriers and evaluate mature standards later; no algorithm, covered-field set, product, or service is selected in FEATURE-0013. | Deferred | ADH-2026-017; ADR-F13-002 |
| Policy/evaluator and workflow runtimes | Wrap | CEL, OPA, Cedar, durable-execution, and workflow products are later candidates behind owning-feature adapters or ports; FEATURE-0013 defines no wrapper or runtime implementation. | Deferred | ADH-2026-017; DEC-0036 |

## 3. Guardrails and source authority

Authority order (highest first): the consolidated architecture main body
(sections 27.8 and 27.9 are closed and normative), ADH-2026-017, the approved
requirements, then the FEATURE-0011/FEATURE-0012 locked dependency files.

### 3.1 Implementation-class rule

- Every requirement, design component, and task declares exactly one class:
  `CONTRACT_NOW`, `INVARIANT_FOR_LATER`, or `DEFERRED`.
- Only `CONTRACT_NOW` items produce implementation tasks.
- `INVARIANT_FOR_LATER` and `DEFERRED` rows produce no source task; they receive
  documentation and traceability placeholders only.
- Tasks must not introduce a new architecture decision or reinterpret a closed
  one. If a task cannot proceed without changing a closed architecture semantic,
  stop and report `ARCHITECTURE_DECISION_REQUIRED`.

### 3.2 Section 27.9 anti-wandering controls (mandatory)

- `SecurityExceptionRef` is structural evidence only. No security-exception
  approval service, approval workflow, authority model, issuance, revocation, or
  runtime enforcement is created.
- Strategy metadata is never an independent fail-open authority. Any strategy
  fail-open handling is effective only through a valid governing-profile
  `SecurityExceptionRef`; otherwise the strategy fails closed.
- `DecisionRequest` is calling-domain-owned. No shared, canonical, common,
  persisted, public, or registry-backed DecisionRequest schema, Go type,
  resource, or persistence model is created. It conforms only to the inherited
  FEATURE-0012 `TransientRequestResult` boundary.
- Completed synchronous handling returns a `DecisionRecord` value/reference;
  accepted asynchronous handling returns/references the existing FEATURE-0012
  `Operation` (`LongRunningOperation`) value/reference; rejection returns
  inherited Problem Details.
- No `PendingDecision`, `DeferredDecision`, `DecisionPending`, or other
  pending-decision or non-final response envelope or status is implemented.
- `metadata.scopeRef` is the sole scope authority. Tasks introduce no second
  scope source, no scope-kind alias, no owner-derived scope, no subject-derived
  scope, and no additional audit scope field.
- ServiceInstance is a typed subject/reference and is never a ScopeKind
  (F13-SCOPE-004); its `metadata.scopeRef` is the containing Project.
- FEATURE-0013 reuses the FEATURE-0012 baseline workflow exactly:
  `api/schemas/baseline/BASELINE_MANIFEST.json` and
  `api/schemas/baseline/BASELINE_APPROVALS.json`. Do not create any parallel
  schema-baseline manifest, schema-diff directory, or schema-approval directory.
  Reuse only the FEATURE-0012 baseline manifest and approvals files.
- The FEATURE-0013 AuditEvent payload/linkage extension lives in the
  `internal/decision` domain package; `apiconform` remains binding/conformance
  support only and never becomes the canonical domain owner.
- Versioned registries use composite `(id, version)` keys; Go maps are
  derivative views of the JSON bundle.
- A present-but-empty optional `TrustCarrier` is rejected as invalid (RID-08);
  a present-but-empty carrier never satisfies provenance, trust, identity, or
  evidence.
- Projection conformance is structural only. No projection/redaction runtime or
  projected-output generation is implemented.
- Matrix E worksheet tasks may stage architecture-owned values and blank
  human-only fields only. Do not calculate, infer, downgrade, or upgrade
  residual risk.
- JSON pre-scan uses JSON-token scanning; YAML pre-scan uses `yaml.Node` before
  normalization, reusing `apivalid.StrictDecode`.
- `GraphEdge` carries a stable edge ID distinct from its `from`/`to` endpoints;
  graph reference/version validation belongs to bundle resolution before the
  in-memory graph build.

### 3.3 Inherited limit hierarchy (tasks must preserve and test)

FEATURE-0012 platform/schema ceilings are the absolute outer bounds.
DecisionProfile ceilings must remain within those inherited FEATURE-0012
ceilings. Evaluator-specific ceilings may only narrow the DecisionProfile and
must never widen it. The effective limit is the most restrictive applicable
value. Missing, unknown, incomparable, or above-ceiling mandatory limits fail
closed.

### 3.4 Immutability and erasure boundary (tasks must preserve and test)

`DecisionRecord` and `AuditEvent` are append-only. Correction, supersession,
revocation, retention, and lawful erasure never rewrite the canonical immutable
record in place. The only alternatives are separately controlled payload
custody, key destruction, or a linked tombstone/redaction record. FEATURE-0013
implements no erasure, key destruction, or cryptographic process.

### 3.5 Public error contract (tasks must implement and test)

FEATURE-0013 validation failures inherit the FEATURE-0012 Problem Details
envelope with `type` = `urn:sovrunn:problem:validation-failed`, HTTP `422`,
top-level `code` = `VALIDATION_FAILED`, an RFC 6901 `violations[].field`
pointer, one closed `violations[].code`, and a redactable message. No task
creates or reinterprets a public code, `Problem.type`, top-level code, HTTP
meaning, or code prefix. The closed `violations[].code` registry is exactly:

- `DECISION_PROFILE_UNKNOWN`, `DECISION_PROFILE_VERSION_UNSUPPORTED`,
  `DECISION_PROFILE_INACTIVE`, `DECISION_PROFILE_SCHEMA_INVALID`,
  `DECISION_PROFILE_LIMIT_INVALID`;
- `DECISION_EVALUATION_TYPE_UNSUPPORTED`, `DECISION_EVALUATION_RESULT_INVALID`,
  `DECISION_EVALUATION_SCOPE_MISMATCH`, `DECISION_EVALUATION_LIMIT_EXCEEDED`;
- `DECISION_COMPOSITION_STRATEGY_UNSUPPORTED`,
  `DECISION_COMPOSITION_GRAPH_INVALID`, `DECISION_COMPOSITION_CYCLE`,
  `DECISION_COMPOSITION_LIMIT_EXCEEDED`, `DECISION_COMPOSITION_INPUT_CONFLICT`;
- `DECISION_OBLIGATION_UNKNOWN`, `DECISION_OBLIGATION_INVALID`,
  `DECISION_OBLIGATION_UNSUPPORTED_MANDATORY`;
- `DECISION_TRUST_REQUIRED`, `DECISION_TRUST_UNKNOWN`, `DECISION_TRUST_EXPIRED`,
  `DECISION_TRUST_REVOKED`, `DECISION_TRUST_MISMATCH`;
- `DECISION_SCOPE_REQUIRED`, `DECISION_SCOPE_INVALID`, `DECISION_SCOPE_CONFLICT`,
  `DECISION_SCOPE_MISMATCH`, `DECISION_SCOPE_WIDENING`;
- `DECISION_RELATIONSHIP_KIND_INVALID`, `DECISION_RELATIONSHIP_TARGET_INVALID`,
  `DECISION_RELATIONSHIP_CYCLE`, `DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED`,
  `DECISION_RELATIONSHIP_CONFLICT`.

### 3.6 Risk namespaces

FEATURE-0012 owns `F12-R01` through `F12-R16`. FEATURE-0013 owns `F13-R01`
through `F13-R31`. No task merges, renumbers, aliases, inherits, or ordinally
maps identifiers between the two namespaces.

## 4. Implementation tasks — Go domain value types (`internal/decision`)

All tasks in this section are class `CONTRACT_NOW` and follow
`docs/engineering/go-coding-guardrails.md`. Each type anonymously embeds
`apimeta.TypeMeta` and carries `apimeta.ObjectMeta` where design section 7
requires it. Optional singular references use `*apimeta.TypedRef` + `omitempty`;
required collections are non-pointer slices; closed vocabularies are
string-backed named types with `Valid()` methods (RID-02, closure 11).

### T-001 — Package doc and closed vocabularies

- Class: `CONTRACT_NOW`.
- Design: §5, §6.3, §7.5; RID-01. Requirements: F13-FORM-001/002, F13-AUTH-001,
  F13-ADJ-001, F13-AUDIT-002, F13-IMMUT-001/002, F13-SEC-001. Architecture:
  §6.2–6.6, §12.4, §27.9(7); AD-003, AD-030, AD-031.
- Objective: create the package and the exact closed vocabularies as
  string-backed named types, each with `Valid()` and a stable-order `All…()`.
- Files: `internal/decision/doc.go` (package doc + FEATURE-0013 domain-ownership
  note, closure 7), `internal/decision/vocab.go` (`DecisionForm`,
  `DecisionAuthority`, `AdjudicationOutcome`, `AuditOutcome`, `RelationshipKind`,
  `Sensitivity`, `Finality`, `EffectiveState`).
- Notes: reproduce architecture values exactly; add no value. `Finality` is
  `FINAL` only; `EffectiveState` is projection-only (`EFFECTIVE`, `SUPERSEDED`,
  `REVOKED`).
- Tests: `internal/decision/vocab_test.go` — `Valid()` accepts each listed value
  and rejects unknown values; `All…()` returns the exact stable-ordered set.
- Acceptance: every vocabulary matches design §7.5 exactly; no extra value.
- Commit: `feat(decision): add FEATURE-0013 closed decision vocabularies`.

### T-002 — DecisionRecord value types

- Class: `CONTRACT_NOW`.
- Design: §7.1, §7.9; RID-02. Requirements: F13-OBJ-001, F13-IMMUT-001,
  F13-AUDIT-001, F13-SCALE-001/002, F13-SCOPE-002. Architecture: §5.4, §6.1;
  AD-001, AD-004, AD-036.
- Objective: define `DecisionRecord`, `DecisionBody`, `DecisionResultBody`,
  `DecisionRationale`, `Obligation`, `Correlation`, `CompositionRef`,
  `SovereigntyCarrier`, `ProfileRef`.
- Files: `internal/decision/decisionrecord.go`.
- Notes: `apiVersion=governance.sovrunn.io/v1alpha1`, `kind=DecisionRecord`;
  `metadata.scopeRef` is the sole scope authority; `Correlation.AuditEventRefs`
  is required non-empty (`minItems: 1`); `Finality` is always `FINAL` on a
  persisted record. Append-only; no mutable `spec`/`status`.
- Tests: `internal/decision/decisionrecord_test.go` — JSON round-trip;
  `omitempty` behavior for optional refs; required fields present.
- Acceptance: struct tags and requiredness match design §7.6/§7.9.
- Commit: `feat(decision): add DecisionRecord value types`.

### T-003 — EvaluationResult value types

- Class: `CONTRACT_NOW`.
- Design: §7.2. Requirements: F13-OBJ-004, F13-EVAL-001/002/003/004.
  Architecture: §7.1; AD-005, AD-037.
- Objective: define `EvaluationResult` and evaluator identity/provenance
  carriers (evaluator identity+version, input-snapshot reference or opaque
  integrity carrier, one result, timing envelope, explicit
  success/non-match/indeterminate/timeout/error semantics, structural trust
  state).
- Files: `internal/decision/evaluationresult.go`.
- Notes: `EmbeddedValue` when captured in a record; `TransientRequestResult`
  when exchanged. Scope is derived from the containing record; evaluator-supplied
  scope metadata is evidence only. No independent durable lifecycle.
- Tests: `internal/decision/evaluationresult_test.go` — round-trip; optional
  integrity-carrier omission; explicit result-status enumeration.
- Acceptance: fields match design §7.2; no independent scope field.
- Commit: `feat(decision): add EvaluationResult value types`.

### T-004 — DecisionProfile and SecurityExceptionRef value types

- Class: `CONTRACT_NOW`.
- Design: §7.3, §7.7; RID-08 context. Requirements: F13-OBJ-003, F13-PROF-001,
  F13-PROF-002, F13-PROF-005. Architecture: §6.8, §27.9(1); AD-029, AD-043.
- Objective: define `DecisionProfile`, `ProfileSpec`, limit descriptors,
  projection specs (`ProjectionRule`), and the bounded `SecurityExceptionRef`.
- Files: `internal/decision/profile.go`.
- Notes: `VersionedDefinition` shape; definition content in `spec`.
  `SecurityExceptionRef` is a structural approval-evidence reference only
  (apiVersion/kind/name/uid + bounded metadata). Tasks create no exception
  approval workflow, issuance, or revocation. `CompensatingControls` is a bounded
  non-empty `[]string`; no control object or `ControlRef` type is added.
  `ExceptionScopeRef` is exception-evidence scope only and never overrides the
  canonical `metadata.scopeRef`.
- Tests: `internal/decision/profile_test.go` — round-trip; `SecurityExceptionRef`
  requiredness (uid pinning, non-empty controls and covered modes).
- Acceptance: shapes match design §7.3/§7.7; no approval-workflow field invented.
- Commit: `feat(decision): add DecisionProfile and SecurityExceptionRef types`.

### T-005 — Trust and identity carriers

- Class: `CONTRACT_NOW`.
- Design: §7.9 (`TrustCarrier`, `SemanticDecisionIdentity`); RID-08.
  Requirements: F13-TRUST-001/002, F13-SCALE-001/002. Architecture: §12.3;
  AD-022, AD-036.
- Objective: define algorithm-agile `TrustCarrier` (structural trust-state +
  canonicalization/digest/covered-field/signature identifiers), the retry
  idempotency key, and `SemanticDecisionIdentity` (basis, input refs,
  canonical-scope descriptor derived from `metadata.scopeRef`, opaque digest).
- Files: `internal/decision/trust.go`, `internal/decision/identity.go`.
- Notes: no algorithm is selected, computed, signed, or verified. A
  present-but-empty optional `TrustCarrier` is invalid (RID-08).
  `CanonicalScope` is derived, not an independent authority.
- Tests: `internal/decision/trust_test.go`, `internal/decision/identity_test.go`
  — carrier round-trip; retry key distinct from semantic identity.
- Acceptance: fields match design §7.9; no cryptographic behavior.
- Commit: `feat(decision): add algorithm-agile trust and identity carriers`.

### T-006 — Sensitivity vocabulary and classification floor map

- Class: `CONTRACT_NOW`.
- Design: §7.5, §12. Requirements: F13-SEC-001/002/003. Architecture: §12.4;
  AD-019.
- Objective: define the four-value `Sensitivity` type and the
  `DataClassification`→sensitivity-floor map (raise-only, fail-closed on
  unmapped/unknown/contradictory/below-floor).
- Files: `internal/decision/sensitivity.go`.
- Notes: the four core values coexist with FEATURE-0012 `DataClassification` and
  never lower it. Structural only; no content-scanning engine.
- Tests: `internal/decision/sensitivity_test.go` — floor map is raise-only;
  below-floor fails; unmapped fails.
- Acceptance: mapping matches design §12 exactly.
- Commit: `feat(decision): add sensitivity vocabulary and floor mapping`.

### T-007 — AuditEvent FEATURE-0013 extension value types

- Class: `CONTRACT_NOW`.
- Design: §7.4, §7.9 (`AuditLinkage`). Requirements: F13-AUDIT-001/002/003,
  F13-COMPAT-009. Architecture: §10, §27.9(7); AD-015, AD-026.
- Objective: define the FEATURE-0013-owned AuditEvent payload/linkage extension
  value types (`AuditLinkage`: decisionRef, optional obligation identity, reused
  `Correlation`) in the `internal/decision` domain package.
- Files: `internal/decision/auditevent_ext.go`.
- Notes: the FEATURE-0012 `AuditEvent` base envelope is preserved; base fields
  are not copied into a divergent parallel type. Audit outcomes remain the closed
  `Succeeded`/`Denied`/`Failed` set.
- Tests: `internal/decision/auditevent_ext_test.go` — linkage round-trip;
  reuse of the FEATURE-0012 correlation carrier.
- Acceptance: extension types live in `internal/decision`; no base-field
  duplication.
- Commit: `feat(decision): add FEATURE-0013 AuditEvent linkage extension types`.

### T-008 — Bounded in-memory decision graph (`internal/decision/graph`)

- Class: `CONTRACT_NOW`.
- Design: §5, §7.8; RID-04, closures 3/4. Requirements: F13-COMP-003/004.
  Architecture: §7.4, §27.9(3), §27.9(4); AD-006, AD-040.
- Objective: define `Definition`, `Node`, `GraphEdge` (package-local `Edge`),
  `ParallelGroup`, `Bounds`, the `Build(def, bounds)` adjacency builder, and
  deterministic `Accounting()` (depth/width/fan-out/critical path).
- Files: `internal/decision/graph/graph.go`.
- Notes: `GraphEdge` has a stable edge ID distinct from its `from`/`to`
  endpoints; duplicate `(from,to)` pairs are invalid unless the profile permits
  labelled multi-edges by edge ID; approved node kinds are exactly `input`,
  `evaluation`, `composition`, `decision`; orientation is dependency→dependent.
  `Build` receives an already-resolved definition and performs no registry or
  version lookup (closure 4). Accounting counts each node/edge once.
- Tests: `internal/decision/graph/graph_test.go` — unique IDs, distinct edge ID,
  parallel-group independence, deterministic accounting, single-node
  empty-edge case.
- Acceptance: structure and rules match design §7.8; builder does no lookup.
- Commit: `feat(decision): add bounded in-memory decision graph builder`.

### T-009 — Declarative bundle and `(id, version)` registries (`internal/decision/bundle`)

- Class: `CONTRACT_NOW`.
- Design: §6.3, §8, §10; RID-03, closures 8/13. Requirements: F13-PROF-004,
  F13-PROF-007/008, F13-COMPAT-005. Architecture: §6.8, §10, §27.9(8);
  AD-025, AD-038.
- Objective: define `Bundle`, `BundleView`, `(id, version)`-keyed registries,
  duplicate detection, `Load(data, mode)` via `apivalid.StrictDecode`, and
  `ResolveGraph(ref)` performing graph reference/version resolution before build.
- Files: `internal/decision/bundle/bundle.go`, `internal/decision/bundle/load.go`.
- Notes: the JSON bundle is the semantic source of truth; Go maps are derivative
  views. Versioned registries are keyed by composite `(id, version)`. Offline
  verification only: no network fetch, digest computation, or signature
  verification. Unknown/expired/revoked/mismatched bundles fail closed with no
  fallback to older versions.
- Tests: `internal/decision/bundle/bundle_test.go`,
  `internal/decision/bundle/load_test.go` — duplicate `(id,version)` detection,
  lookup by composite key, JSON/YAML decode equivalence via `StrictDecode`,
  graph reference/version resolution before build.
- Acceptance: keys and resolution match design §10; no runtime registry service.
- Commit: `feat(decision): add declarative bundle and versioned registries`.

### T-010 — Pure composition strategies (`internal/decision/compose/strategy.go`)

- Class: `CONTRACT_NOW`.
- Design: §6.3, §8; RID-06, RID-10, closure 2. Requirements: F13-COMP-001/002.
  Architecture: §7.2, §27.9(2); AD-005, AD-043.
- Objective: define the pure `Strategy func(StrategyInput) (StrategyOutput,
  *apiproblem.Problem)` interface, registered strategy metadata, and the
  fail-open gate that is effective only through a valid governing-profile
  `SecurityExceptionRef` (otherwise fail closed).
- Files: `internal/decision/compose/strategy.go`.
- Notes: strategies perform no I/O and are never an independent fail-open
  authority. No approval workflow is created here.
- Tests: `internal/decision/compose/strategy_test.go` — purity (no I/O),
  fail-closed default, fail-open effective only with valid exception evidence.
- Acceptance: matches design §8/RID-10; strategy metadata carries no approval
  authority.
- Commit: `feat(decision): add pure composition strategy interface`.

### T-011 — Deterministic replay kernel (`internal/decision/compose/replay.go`)

- Class: `CONTRACT_NOW`.
- Design: §6.3, §8, §13; RID-06. Requirements: F13-EVAL-002, F13-SCALE-005.
  Architecture: §7.2, §9; AD-012, AD-037.
- Objective: define `Replay(captured, strategy, in)` that consumes captured
  evaluator output and never reruns an evaluator; holds no session state and
  performs no persistence, delivery, locking, caching, or scaling.
- Files: `internal/decision/compose/replay.go`.
- Notes: deterministic pure function of versioned inputs and explicit
  fixture-local bounded configuration.
- Tests: `internal/decision/compose/replay_test.go` — identical output from
  identical captured input; evaluator never invoked.
- Acceptance: matches design §8; pure and deterministic.
- Commit: `feat(decision): add deterministic composition replay kernel`.

## 5. Schema and binding tasks

All schema tasks are class `CONTRACT_NOW` and are authorized by design §6.1/§6.2
only. Schemas are JSON Schema with `$ref` composition and `x-sovrunn-profile`
annotations. No parallel baseline/diff/approval hierarchy is created.

### T-012 — `_common` sub-schemas

- Class: `CONTRACT_NOW`.
- Design: §6.1, §7.6–7.9; RID-03, closures 1/3/9/11. Requirements: F13-SEC-001,
  F13-TRUST-001, F13-IMMUT-002, F13-SCALE-002, F13-PROF-002, F13-COMP-003.
  Architecture: §6.8, §12.3/12.4, §27.9. Class per file: `CONTRACT_NOW`.
- Objective: create the shared sub-schemas.
- Files: `api/schemas/_common/decision-linkage.json`,
  `api/schemas/_common/decision-profile-ref.json`,
  `api/schemas/_common/trust-carrier.json`,
  `api/schemas/_common/security-exception-ref.json`,
  `api/schemas/_common/semantic-decision-identity.json`,
  `api/schemas/_common/decision-relationship.json`,
  `api/schemas/_common/sensitivity.json`,
  `api/schemas/_common/decision-graph.json`.
- Notes: enforce requiredness per design §7.6 (`required`, `minItems`,
  `nullable`); `security-exception-ref.json` carries apiVersion/kind/name/uid +
  bounded metadata only; `trust-carrier.json` is structural only; the graph
  schema gives node/edge IDs with the edge ID distinct from `(from,to)`.
- Tests: covered by TypeBinding verification (T-016) and conformance (§6).
- Acceptance: each schema matches design §6.1/§7.x; no invented field.
- Commit: `feat(schemas): add FEATURE-0013 _common decision sub-schemas`.

### T-013 — Canonical object and bundle schemas

- Class: `CONTRACT_NOW`.
- Design: §6.1, §7.1–7.3, §10; RID-03. Requirements: F13-OBJ-001/003/004,
  F13-PROF-004. Architecture: §5.4, §6.1, §6.8; AD-025, AD-029.
- Objective: create `api/schemas/decision-record.json`
  (`x-sovrunn-profile: ImmutableRecord`, `x-sovrunn-boundary: governance-only`;
  non-Platform `metadata.scopeRef` requires `uid`),
  `api/schemas/decision-profile.json` (`VersionedDefinition`; content in `spec`),
  `api/schemas/evaluation-result.json` (`EmbeddedValue`/`TransientRequestResult`),
  and `api/schemas/decision-profile-bundle.json` (declarative registry;
  `(id,version)` keyed entries + algorithm-agile trust carriers).
- Files: the four schemas above.
- Notes: reference `_common` sub-schemas via `$ref`. No new resource profile is
  introduced.
- Tests: TypeBinding verification (T-016) and conformance (§6).
- Acceptance: profiles and boundaries match design §6.1.
- Commit: `feat(schemas): add DecisionRecord/Profile/EvaluationResult/bundle schemas`.

### T-014 — AuditEvent schema expansion (additive)

- Class: `CONTRACT_NOW`.
- Design: §6.2, §7.4; closure 6. Requirements: F13-AUDIT-003, F13-COMPAT-007/008.
  Architecture: §15.3, §27.9(6); AD-015.
- Objective: expand `api/schemas/audit-event.json` `x-sovrunn-allowed-scopes`
  from `[Organization]` to the six governance scopes and add optional additive
  decision-linkage properties via `$ref: _common/decision-linkage.json`.
- Files: `api/schemas/audit-event.json`.
- Notes: preserve the base `ImmutableRecord` envelope and every existing field;
  no field removed, renamed, or reinterpreted. The Organization-scoped alpha
  fixture must remain valid (regression).
- Tests: conformance regression (F13-COMPAT-07) and six-scope acceptance
  (F13-COMPAT-08) in §6.
- Acceptance: change is additive/expanded only; base preserved.
- Commit: `feat(schemas): expand AuditEvent scopes and add decision linkage`.

### T-015 — Baseline manifest and approvals update (reuse FEATURE-0012 mechanism)

- Class: `CONTRACT_NOW`.
- Design: §6.2, §14; closure 6. Requirements: F13-COMPAT-001/006. Architecture:
  §15, §27.9(6); AD-015, AD-044.
- Objective: recompute digests for changed/added schemas using the existing
  FEATURE-0012 mechanism and record the approval-controlled baseline update.
- Files: `api/schemas/baseline/BASELINE_MANIFEST.json`,
  `api/schemas/baseline/BASELINE_APPROVALS.json`.
- Notes: reuse the exact FEATURE-0012 baseline workflow. Do not create any
  parallel schema-baseline manifest, schema-diff directory, or schema-approval
  directory. Reuse only the FEATURE-0012 baseline manifest and approvals files.
  This task stages the digest recomputation only; it does not mark human
  approval complete or advance any stage status.
- Tests: existing baseline verification tooling passes for the recomputed
  digests.
- Acceptance: only `BASELINE_MANIFEST.json`/`BASELINE_APPROVALS.json` change;
  no parallel hierarchy is introduced.
- Commit: `chore(schemas): recompute baseline digests for FEATURE-0013 schemas`.

### T-016 — apiconform TypeBinding registration

- Class: `CONTRACT_NOW`.
- Design: §6.4; RID-01. Requirements: F13-COMPAT-001. Architecture: §17;
  AD-028, AD-045.
- Objective: append FEATURE-0013 `TypeBinding` entries mapping each Go value type
  to its schema in `internal/apiconform/bindings.go` (registration only).
- Files: `internal/apiconform/bindings.go`.
- Notes: `apiconform` remains binding/conformance support only; it does not
  become the canonical domain owner and must not import `decision/validate` in a
  way that creates a cycle. Use `apischema.VerifyGoTypeAgainstSchema` for Go↔schema
  agreement.
- Tests: `internal/apiconform/bindings_test.go` — `VerifyGoTypeAgainstSchema`
  passes for every FEATURE-0013 binding (requiredness/nullable/minItems).
- Acceptance: bindings registered; Go↔schema agreement holds.
- Commit: `feat(apiconform): register FEATURE-0013 type bindings`.

## 6. Validator tasks (`internal/decision/validate`)

All validators are class `CONTRACT_NOW`, pure or bounded in-memory, fail-closed,
and emit the inherited FEATURE-0012 Problem Details envelope with exactly one
closed `violations[].code` and an RFC 6901 pointer. No task introduces a new
public code, `Problem.type`, top-level code, HTTP meaning, or prefix.

### T-017 — Closed violation-code constants

- Class: `CONTRACT_NOW`.
- Design: §6.3, §9.2. Requirements: F13-ERR-001/002/003/004. Architecture:
  §15.1, §15.2; AD-045.
- Objective: define the closed `DECISION_*` `violations[].code` constants in
  `internal/decision/validate/codes.go`, mirroring architecture §15.2 exactly.
- Files: `internal/decision/validate/codes.go`.
- Notes: the constant set is exactly the 32 codes in §3.5; no code is added,
  renamed, or aliased.
- Tests: `internal/decision/validate/codes_test.go` — the constant set equals
  the architecture-owned closed registry.
- Acceptance: exactly the closed registry; no invented code.
- Commit: `feat(decision): add closed DECISION_* violation-code constants`.

### T-018 — Scope validation pass

- Class: `CONTRACT_NOW`.
- Design: §9.1 step 2, §9.2, §7.9. Requirements: F13-SCOPE-001…008,
  F13-SCALE-002. Architecture: §6.1, §27.8; AD-015.
- Objective: implement `ValidateScope` (sole-authority, canonical Platform
  absence, six-value set, ServiceInstance-as-subject, and
  `semanticIdentity.canonicalScope` matches the computed `metadata.scopeRef`
  value).
- Files: `internal/decision/validate/scope.go`.
- Notes: map failures to `DECISION_SCOPE_REQUIRED`, `DECISION_SCOPE_INVALID`,
  `DECISION_SCOPE_CONFLICT`, `DECISION_SCOPE_MISMATCH`, `DECISION_SCOPE_WIDENING`.
  ServiceInstance is a subject, never a ScopeKind.
- Tests: `internal/decision/validate/scope_test.go` — each mapped code at its
  RFC 6901 pointer; canonical-scope descriptor mismatch → `DECISION_SCOPE_MISMATCH`.
- Acceptance: matches design §9.1/§9.2.
- Commit: `feat(decision): add scope validation pass`.

### T-019 — Profile validation pass and limit hierarchy

- Class: `CONTRACT_NOW`.
- Design: §7.3, §7.7, §9.1 step 3. Requirements: F13-PROF-001…008,
  F13-EVAL-02/03/04. Architecture: §6.8, §27.9(1); AD-029, AD-040, AD-043.
- Objective: implement `ValidateProfile` (identity/lifecycle/schema/limits and
  `SecurityExceptionRef` structural checks) applying the inherited limit
  hierarchy: FEATURE-0012 ceilings are outer bounds, profile ceilings remain
  within them, evaluator ceilings only narrow, most restrictive applicable value
  wins.
- Files: `internal/decision/validate/profile.go`.
- Notes: `SecurityExceptionRef` failures map to `DECISION_PROFILE_SCHEMA_INVALID`
  at the exact pointers in design §7.7. Limit failures map to
  `DECISION_PROFILE_LIMIT_INVALID`. No approval workflow is implemented.
- Tests: `internal/decision/validate/profile_test.go` — limit hierarchy
  boundaries; `SecurityExceptionRef` uid/scope/expiry/controls; unknown/inactive/
  version-unsupported profiles.
- Acceptance: matches design §7.3/§7.7/§9; codes correct.
- Commit: `feat(decision): add profile validation and limit hierarchy`.

### T-020 — Evaluation validation pass

- Class: `CONTRACT_NOW`.
- Design: §9.1 step 4, §9.2. Requirements: F13-EVAL-001…005. Architecture: §7.1;
  AD-005.
- Objective: implement `ValidateEvaluation` (evaluator type/result/limit/scope
  derivation) mapping to `DECISION_EVALUATION_TYPE_UNSUPPORTED`,
  `DECISION_EVALUATION_RESULT_INVALID`, `DECISION_EVALUATION_SCOPE_MISMATCH`
  (no existence disclosure), `DECISION_EVALUATION_LIMIT_EXCEEDED`.
- Files: `internal/decision/validate/evaluation.go`.
- Notes: evaluator scope is derived from the record; evaluator metadata is
  evidence only.
- Tests: `internal/decision/validate/evaluation_test.go` — each mapped code;
  scope mismatch without disclosure; evaluator-widening rejected.
- Acceptance: matches design §9.
- Commit: `feat(decision): add evaluation validation pass`.

### T-021 — Composition and graph validation pass

- Class: `CONTRACT_NOW`.
- Design: §7.8, §9.1 step 5, §9.2. Requirements: F13-COMP-001…004.
  Architecture: §7.2, §7.4, §27.9(3)/(4); AD-006.
- Objective: implement composition validation (strategy support; graph
  resolution in bundle → in-memory build → bounds/cycle/limit/input-conflict)
  with the deterministic pointer-base rule from design §7.8.
- Files: `internal/decision/validate/composition.go`.
- Notes: map to `DECISION_COMPOSITION_STRATEGY_UNSUPPORTED`,
  `DECISION_COMPOSITION_GRAPH_INVALID`, `DECISION_COMPOSITION_CYCLE`,
  `DECISION_COMPOSITION_LIMIT_EXCEEDED`, `DECISION_COMPOSITION_INPUT_CONFLICT`.
  Graph version/reference mismatch is resolved in `bundle.ResolveGraph` before
  build.
- Tests: `internal/decision/validate/composition_test.go` — graph structural
  failures at graph-local and bundle-embedded pointers; cycle; bound exceeded.
- Acceptance: matches design §7.8/§9.
- Commit: `feat(decision): add composition and graph validation pass`.

### T-022 — Trust validation pass (empty-state rule)

- Class: `CONTRACT_NOW`.
- Design: §9.1 step 6, RID-08. Requirements: F13-TRUST-001…005. Architecture:
  §12.3; AD-022.
- Objective: implement `ValidateTrust` (carrier presence/syntax/state; required
  vs optional; present-but-empty rejected).
- Files: `internal/decision/validate/trust.go`.
- Notes: map to `DECISION_TRUST_REQUIRED`, `DECISION_TRUST_UNKNOWN`,
  `DECISION_TRUST_EXPIRED`, `DECISION_TRUST_REVOKED`, `DECISION_TRUST_MISMATCH`.
  No cryptographic verification; structural/opaque state only.
- Tests: `internal/decision/validate/trust_test.go` — required-absent fails;
  present-empty fails; opaque expected≠observed fails closed.
- Acceptance: matches design §9/RID-08.
- Commit: `feat(decision): add trust validation pass`.

### T-023 — Obligation validation pass

- Class: `CONTRACT_NOW`.
- Design: §9.1 step 7, §9.2. Requirements: F13-OBLIG-001/002/003. Architecture:
  §6.10; AD-032, AD-043.
- Objective: implement `ValidateObligations` (vocabulary/payload/mandatory
  support) mapping to `DECISION_OBLIGATION_UNKNOWN`, `DECISION_OBLIGATION_INVALID`,
  `DECISION_OBLIGATION_UNSUPPORTED_MANDATORY`.
- Files: `internal/decision/validate/obligation.go`.
- Notes: an unsupported mandatory obligation makes an ALLOWED decision
  unenforceable.
- Tests: `internal/decision/validate/obligation_test.go` — each mapped code.
- Acceptance: matches design §9.
- Commit: `feat(decision): add obligation validation pass`.

### T-024 — Relationship validation pass and conflict matrix

- Class: `CONTRACT_NOW`.
- Design: §9.1 step 8, §9.3; closure 5. Requirements: F13-IMMUT-002/004/005.
  Architecture: §6.5, §6.7, §27.9(5); AD-004, AD-034.
- Objective: implement the deterministic relationship conflict matrix using only
  the five approved relationship codes and the algorithm order in design §9.3.
- Files: `internal/decision/validate/relationship.go`.
- Notes: append-only; effective state is a projection only and is never written
  back. Map to `DECISION_RELATIONSHIP_KIND_INVALID`,
  `DECISION_RELATIONSHIP_TARGET_INVALID`, `DECISION_RELATIONSHIP_CYCLE`,
  `DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED`, `DECISION_RELATIONSHIP_CONFLICT`.
- Tests: `internal/decision/validate/relationship_test.go` — each governed case
  in the design §9.3 table at its deciding step and pointer.
- Acceptance: matches design §9.3 exactly.
- Commit: `feat(decision): add relationship validation and conflict matrix`.

### T-025 — Sensitivity and projection validation pass

- Class: `CONTRACT_NOW`.
- Design: §9.1 step 9, §12; closure 10. Requirements: F13-SEC-002…007.
  Architecture: §12.4; AD-019.
- Objective: implement sensitivity ceiling/floor/content-category checks and
  projection pointer target/overlap checks (resolve against the profile-declared
  canonical view; overlap = equality + containment).
- Files: `internal/decision/validate/sensitivity.go`.
- Notes: projection conformance is structural only; no redaction execution or
  projected-output generation. Below-floor/unmapped classification fails closed.
- Tests: `internal/decision/validate/sensitivity_test.go` — ceiling boundary;
  floor cannot be lowered; overlapping includes/excludes rejected.
- Acceptance: matches design §12; no runtime redaction.
- Commit: `feat(decision): add sensitivity and projection validation pass`.

### T-026 — Record and AuditEvent orchestration (fixed pass order)

- Class: `CONTRACT_NOW`.
- Design: §8, §9.1. Requirements: F13-OBJ-007, F13-ERR-001, F13-AUDIT-001.
  Architecture: §8, §9; AD-009.
- Objective: implement `ValidateDecisionRecord`, `ValidateAuditEvent`, and
  `ValidateRelationship` orchestration running the fixed nine-pass order
  (decode/structural → scope → profile → evaluation → composition → trust →
  obligation → relationship → sensitivity/projection).
- Files: `internal/decision/validate/validate.go`.
- Notes: decode/structural failures return FEATURE-0012 top-level codes (not a
  `DECISION_*` code). Completion returns `DecisionRecord`; accepted async work
  references the existing FEATURE-0012 `Operation`; rejection returns Problem
  Details. No pending-decision envelope is produced.
- Tests: `internal/decision/validate/validate_test.go` — pass order and
  short-circuit behavior; decode failure uses inherited top-level code.
- Acceptance: matches design §9.1 pass order.
- Commit: `feat(decision): add record/audit validation orchestration`.

## 7. Fixture and conformance tasks

All fixtures are class `CONTRACT_NOW`, deterministic, and pure or bounded
in-memory. They authorize no persistence, networking, queues, workers,
cryptographic execution, AI runtime, or external integration. Coverage is
counted by scenario ID; every ID (parent, suffix, additional case) gets an
explicit coverage-matrix row.

### T-027 — Executable conformance checks (`internal/apiconform`)

- Class: `CONTRACT_NOW`.
- Design: §6.4, §13; RID-07. Requirements: F13-CONF-001…005. Architecture: §17;
  AD-028.
- Objective: add `internal/apiconform/decision_checks.go` with executable
  FEATURE-0013 conformance checks that call exported `internal/decision/validate`
  functions and reuse `apivalid.StrictDecode` and `VerifyGoTypeAgainstSchema`.
- Files: `internal/apiconform/decision_checks.go`.
- Notes: checks call validators one-directionally; `apiconform` stays a
  conformance adapter, not the domain owner.
- Tests: `internal/apiconform/decision_checks_test.go` — each check exercises its
  bound fixtures.
- Acceptance: checks resolve every registered scenario ID.
- Commit: `feat(apiconform): add FEATURE-0013 executable conformance checks`.

### T-028 — Coverage matrix

- Class: `CONTRACT_NOW`.
- Design: §6.4, §13, §18.4; RID-07. Requirements: F13-CONF-002/004.
  Architecture: §17; AD-028.
- Objective: add `tests/conformance/decision_coverage_matrix.go` mapping every
  scenario ID (F13-CF-01…28 and all mandatory suffixes; F13-SCOPE-01…11;
  F13-EVAL-01…08; F13-SEC-01…08; F13-TRUST-01…04; F13-COMPAT-01…09) to at least
  one fixture.
- Files: `tests/conformance/decision_coverage_matrix.go`.
- Notes: one physical fixture may cover several IDs, but every ID has its own
  row. Suffixes must not be merged, omitted, renumbered, or invented.
- Tests: `tests/conformance/decision_coverage_matrix_test.go` — every required
  ID resolves to ≥1 fixture.
- Acceptance: full ID coverage per architecture §17/§18.4.
- Commit: `test(conformance): add FEATURE-0013 coverage matrix`.

### T-029 — Positive conformance fixtures

- Class: `CONTRACT_NOW`.
- Design: §6.4, §13; RID-07. Requirements: F13-CONF-001, F13-SCALE-004.
  Architecture: §17; AD-028.
- Objective: add positive fixtures under
  `tests/conformance/fixtures/decision/positive/` for the positive scenarios
  (including the six governance-scope AuditEvent cases, the retained
  Organization-scoped regression, and a ServiceInstance-subject-under-Project
  case).
- Files: `tests/conformance/fixtures/decision/positive/*.json`.
- Notes: fixtures use explicit fixture-local profile values; they invent no
  FEATURE-0013 defaults, fallbacks, or hidden maximums.
- Tests: exercised via T-027/T-028; each positive fixture validates cleanly.
- Acceptance: positive coverage rows resolve.
- Commit: `test(conformance): add FEATURE-0013 positive fixtures`.

### T-030 — Negative conformance fixtures

- Class: `CONTRACT_NOW`.
- Design: §6.4, §9.2, §13; RID-07. Requirements: F13-CONF-001…004.
  Architecture: §15.2, §17; AD-028, AD-045.
- Objective: add negative fixtures under
  `tests/conformance/fixtures/negative/decision/` covering the compound-family
  suffixes and additional scope/eval/sec/trust/compat cases; each returns the
  exact mapped closed `violations[].code` at its RFC 6901 pointer.
- Files: `tests/conformance/fixtures/negative/decision/*.json`.
- Notes: negative fixtures must return exactly one closed code from §3.5; no new
  code is introduced.
- Tests: exercised via T-027/T-028; each negative fixture returns its mapped
  code.
- Acceptance: negative coverage rows resolve with correct codes.
- Commit: `test(conformance): add FEATURE-0013 negative fixtures`.

## 8. Documentation and evidence tasks

### T-031 — Invariant and deferred documentation notes

- Class: documentation-only (records `INVARIANT_FOR_LATER` and `DEFERRED`
  obligations; produces no runtime).
- Design: §6.5, §11, §12. Requirements: F13-AUDIT-005/006/007, F13-AUTHZ-004/005/006,
  F13-SCALE-006/007, F13-TRUST-004. Architecture: §9, §10, §12.1, §16.
- Objective: record the `INVARIANT_FOR_LATER` obligations
  (AD-007/008/010/013/014/035/039/041) and `DEFERRED` items (AD-011, ADR-F13-002)
  as documentation in `internal/decision/doc.go` invariant notes.
- Files: `internal/decision/doc.go` (documentation additions only).
- Notes: these carry the atomic durable-acceptance invariant, every-hop
  authN/authZ duty, break-glass follow-up-review duty, and deferred cryptographic
  selection as text only; no runtime is created.
- Tests: none (documentation only); `go vet`/`gofmt` clean.
- Acceptance: invariant/deferred obligations documented, not implemented.
- Commit: `docs(decision): record FEATURE-0013 invariant and deferred notes`.

### T-032 — Matrix E worksheet staging

- Class: governance-only (evidence staging; no residual-risk calculation).
- Design: §6.5; closure 12. Requirements: F13-GUARD-001/002. Architecture:
  §18.5, §18.6, §27.9(12).
- Objective: create
  `docs/reviews/feature-gates/FEATURE-0013-matrix-e-worksheet.md` staging the
  architecture-owned risk IDs (F13-R01…F13-R31), inherent ratings, controls,
  evidence pointers, target residual ratings, treatment, owner, and trigger,
  with blank human-only disposition fields.
- Files: `docs/reviews/feature-gates/FEATURE-0013-matrix-e-worksheet.md`.
- Notes: stage architecture-owned values and blank human fields only; do not
  calculate, infer, downgrade, or upgrade residual risk; do not mark human
  approval complete, advance stage status, or create review verdicts.
- Tests: none (governance document).
- Acceptance: worksheet stages values with blank human-only fields.
- Commit: `docs(feature-gate): stage FEATURE-0013 Matrix E worksheet`.

## 9. Verification tasks

### T-033 — Final verification, guardrails, cleanup, and clean git

- Class: `CONTRACT_NOW` (verification).
- Design: §14. Requirements: F13-CONF-003, F13-GUARD-003. Architecture: §27.7;
  AD-028.
- Objective: run full verification, prove the import graph is acyclic and free
  of prohibited runtime/provider/crypto imports, clean artifacts, and confirm a
  clean tree.
- Files: none created; may add `internal/decision/imports_test.go` asserting the
  acyclic package DAG and absence of network/database/queue/workflow/crypto/
  provider imports.
- Standard per-task verification:

  ```bash
  docker run --rm -v "$PWD":/src -w /src ${GO_DOCKER_IMAGE:-golang:1.22} \
    sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'
  ```

- Final verification:

  ```bash
  docker run --rm -v "$PWD":/src -w /src ${GO_DOCKER_IMAGE:-golang:1.22} \
    sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && \
    go test -race ./... && go build ./cmd/sovrunn-api'
  ```

- FEATURE-0013 architecture boundary gate:

  ```bash
  python3 scripts/feature-0013-architecture-boundary-check.py \
    --feature FEATURE-0013 --stage tasks --mode post
  ```

- Guardrails and cleanup:

  ```bash
  rm -f sovrunn-api
  rm -rf bin
  ```

  - No `TODO(FEATURE-0013)` under `internal` or `cmd`.
  - No `internal/api` import of `internal/server`.
  - `git status` is clean.

- Tests: `go test -race ./...` passes; import-graph test passes; boundary gate
  passes.
- Acceptance: all verification commands pass; tree clean; no prohibited imports.
- Commit: `test(decision): finalize FEATURE-0013 verification and guardrails`.

## 10. Explicit non-tasks and out-of-scope items

These are intentionally not implemented (`INVARIANT_FOR_LATER`, `DEFERRED`, or
non-goals). No source task is generated for them.

- No shared, canonical, or common DecisionRequest schema/type/registry
  entry/resource/persistence model is created; DecisionRequest is
  calling-domain-owned and conforms only to inherited `TransientRequestResult`
  (F13-OBJ-006; AD-002).
- No `PendingDecision`, `DeferredDecision`, or `DecisionPending` envelope or
  status is created; there is no pending-decision or non-final response
  contract. Accepted async handling references the existing FEATURE-0012
  `Operation` (`LongRunningOperation`) only (F13-OBJ-008; AD-009).
- The generic `Operation` contract and `LongRunningOperation` payload remain
  owned by FEATURE-0012/ADH-2026-013 and are never redefined (Non-goal 15).
- No second scope source, scope-kind alias, owner-derived scope, subject-derived
  scope, or additional audit scope field is created; `metadata.scopeRef` is the
  sole authority (F13-SCOPE-001).
- No security-exception approval service, approval workflow, authority model,
  issuance, or revocation is created; `SecurityExceptionRef` is structural
  evidence only (AD-043; §27.9(1)).
- Do not create any parallel schema-baseline manifest, schema-diff directory, or
  schema-approval directory; reuse only the FEATURE-0012 baseline manifest and
  approvals files per the FEATURE-0012 baseline workflow (§27.9(6)).
- No production decision/authorization/policy/placement/audit/evidence/identity/
  registry/synchronization/cryptographic/AI service, database, queue, broker,
  workflow runtime, worker, controller, cache, lock, or lease is created
  (Non-goals 1–7; AD-011, AD-028).
- No erasure, key destruction, cryptographic process, canonicalization, digest
  computation, signing, or verification is implemented (F13-TRUST-004; ADR-F13-002).
- No projection/redaction runtime or projected-output generation is implemented
  (F13-SEC-007; §27.9(10)).
- Residual risk is not calculated, inferred, downgraded, or upgraded; the Matrix
  E worksheet stages architecture-owned values and blank human fields only
  (§27.9(12)).
- No human approval is marked complete, no stage status is advanced, and no
  review verdict is created by these tasks.

## 11. Task-to-architecture traceability matrix

Task-ID mapping to design sections, requirement IDs, and architecture decision
IDs. Only `CONTRACT_NOW` tasks appear; T-031 is documentation-only and T-032 is
governance-only.

| Task | Design | Requirements | Architecture / AD |
|---|---|---|---|
| T-001 | §5, §6.3, §7.5 | F13-FORM-001/002, F13-AUTH-001, F13-ADJ-001, F13-AUDIT-002, F13-IMMUT-001/002, F13-SEC-001 | §6.2–6.6; AD-003, AD-030, AD-031 |
| T-002 | §7.1, §7.9 | F13-OBJ-001, F13-IMMUT-001, F13-AUDIT-001, F13-SCALE-001/002 | §5.4, §6.1; AD-001, AD-004, AD-036 |
| T-003 | §7.2 | F13-OBJ-004, F13-EVAL-001/002/003/004 | §7.1; AD-005, AD-037 |
| T-004 | §7.3, §7.7 | F13-OBJ-003, F13-PROF-001/002/005 | §6.8, §27.9(1); AD-029, AD-043 |
| T-005 | §7.9 | F13-TRUST-001/002, F13-SCALE-001/002 | §12.3; AD-022, AD-036 |
| T-006 | §7.5, §12 | F13-SEC-001/002/003 | §12.4; AD-019 |
| T-007 | §7.4, §7.9 | F13-AUDIT-001/002/003, F13-COMPAT-009 | §10, §27.9(7); AD-015, AD-026 |
| T-008 | §5, §7.8 | F13-COMP-003/004 | §7.4, §27.9(3)/(4); AD-006, AD-040 |
| T-009 | §6.3, §8, §10 | F13-PROF-004/007/008, F13-COMPAT-005 | §6.8, §10, §27.9(8); AD-025, AD-038 |
| T-010 | §6.3, §8 | F13-COMP-001/002 | §7.2, §27.9(2); AD-005, AD-043 |
| T-011 | §6.3, §8, §13 | F13-EVAL-002, F13-SCALE-005 | §7.2, §9; AD-012, AD-037 |
| T-012 | §6.1, §7.6–7.9 | F13-SEC-001, F13-TRUST-001, F13-IMMUT-002, F13-SCALE-002, F13-PROF-002, F13-COMP-003 | §6.8, §12.3/12.4, §27.9 |
| T-013 | §6.1, §7.1–7.3, §10 | F13-OBJ-001/003/004, F13-PROF-004 | §5.4, §6.1, §6.8; AD-025, AD-029 |
| T-014 | §6.2, §7.4 | F13-AUDIT-003, F13-COMPAT-007/008 | §15.3, §27.9(6); AD-015 |
| T-015 | §6.2, §14 | F13-COMPAT-001/006 | §15, §27.9(6); AD-015, AD-044 |
| T-016 | §6.4 | F13-COMPAT-001 | §17; AD-028, AD-045 |
| T-017 | §6.3, §9.2 | F13-ERR-001/002/003/004 | §15.1, §15.2; AD-045 |
| T-018 | §9.1, §9.2, §7.9 | F13-SCOPE-001…008, F13-SCALE-002 | §6.1, §27.8; AD-015 |
| T-019 | §7.3, §7.7, §9.1 | F13-PROF-001…008, F13-EVAL-02/03/04 | §6.8, §27.9(1); AD-029, AD-040, AD-043 |
| T-020 | §9.1, §9.2 | F13-EVAL-001…005 | §7.1; AD-005 |
| T-021 | §7.8, §9.1, §9.2 | F13-COMP-001…004 | §7.2, §7.4, §27.9(3)/(4); AD-006 |
| T-022 | §9.1, RID-08 | F13-TRUST-001…005 | §12.3; AD-022 |
| T-023 | §9.1, §9.2 | F13-OBLIG-001/002/003 | §6.10; AD-032, AD-043 |
| T-024 | §9.1, §9.3 | F13-IMMUT-002/004/005 | §6.5, §6.7, §27.9(5); AD-004, AD-034 |
| T-025 | §9.1, §12 | F13-SEC-002…007 | §12.4; AD-019 |
| T-026 | §8, §9.1 | F13-OBJ-007, F13-ERR-001, F13-AUDIT-001 | §8, §9; AD-009 |
| T-027 | §6.4, §13 | F13-CONF-001…005 | §17; AD-028 |
| T-028 | §6.4, §13, §18.4 | F13-CONF-002/004 | §17; AD-028 |
| T-029 | §6.4, §13 | F13-CONF-001, F13-SCALE-004 | §17; AD-028 |
| T-030 | §6.4, §9.2, §13 | F13-CONF-001…004 | §15.2, §17; AD-028, AD-045 |
| T-031 | §6.5, §11, §12 | F13-AUDIT-005/006/007, F13-AUTHZ-004/005/006, F13-SCALE-006/007, F13-TRUST-004 | §9, §10, §12.1, §16 |
| T-032 | §6.5 | F13-GUARD-001/002 | §18.5, §18.6, §27.9(12) |
| T-033 | §14 | F13-CONF-003, F13-GUARD-003 | §27.7; AD-028 |

## Architecture traceability

Controlling source:
`docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`,
approved through ADH-2026-017. Implementation classes are architecture-owned and
reproduced from section 19. Only `CONTRACT_NOW` produces implementation tasks;
`INVARIANT_FOR_LATER` and `DEFERRED` rows produce no source task. Subsection
coverage: section 5.4, section 6.1, section 6.8, section 7.1, section 12.3,
section 12.4, section 27.8, and section 27.9 are honored throughout.

### Architecture sections 1–28

| Section | Title | Disposition |
|---|---|---|
| 1 | Decision status and authority | governance-only |
| 2 | Problem statement | translated |
| 3 | Scope | translated |
| 4 | Architectural principles | translated |
| 5 | Canonical object model | translated (T-001…T-007) |
| 6 | Extensible decision semantics | translated (T-001, T-004, T-012, T-013) |
| 7 | Evaluation and composition model | translated (T-003, T-008…T-011, T-021) |
| 8 | Execution model | translated / invariant (T-026, T-031) |
| 9 | Statelessness, scalability, and cost controls | translated / invariant / deferred (T-011, T-031) |
| 10 | Audit consistency | translated / invariant (T-007, T-014, T-031) |
| 11 | Sovereign deployment architecture | translated (T-005, T-009) |
| 12 | Security, privacy, and trust boundaries | translated (T-005, T-006, T-022, T-025) |
| 13 | AI-first architecture | translated (T-025) |
| 14 | Reuse and standards strategy | translated (§2 reuse summary) |
| 15 | Compatibility and evolution | translated (T-014…T-017, T-030) |
| 16 | Observability and operability | invariant (T-031) |
| 17 | Conformance model | translated (T-027…T-030) |
| 18 | Matrix E v2 — architecture risk register | governance-only (T-032) |
| 19 | Decision scorecard | governance-only (classes reproduced below) |
| 20 | Required revisions before approval | governance-only |
| 21 | Important deferred decisions | deferred |
| 22 | Architecture acceptance criteria | governance-only |
| 23 | Reassessment triggers | governance-only |
| 24 | Normative local references | governance-only |
| 25 | External review inputs | governance-only |
| 26 | Foundation-principle validation | governance-only |
| 27 | Anti-overengineering guardrails | translated (T-033) |
| 28 | Lightweight downstream adoption contract | translated (T-032) |

### Architecture decisions AD-001 through AD-045

Exactly one entry per decision; the class is reproduced from architecture
section 19.

| Decision | Class | Disposition |
|---|---|---|
| AD-001 | CONTRACT_NOW | translated (T-002; §7.1) |
| AD-002 | CONTRACT_NOW | translated (§10; no shared DecisionRequest type) |
| AD-003 | CONTRACT_NOW | translated (T-001; §7.5) |
| AD-004 | CONTRACT_NOW | translated (T-002, T-024; §7.1/§9.3) |
| AD-005 | CONTRACT_NOW | translated (T-003, T-010; §8) |
| AD-006 | CONTRACT_NOW | translated (T-008, T-021; §7.8) |
| AD-007 | INVARIANT_FOR_LATER | invariant; documentation-only (T-031) |
| AD-008 | INVARIANT_FOR_LATER | invariant; reference kernel only (T-031) |
| AD-009 | CONTRACT_NOW | translated (T-026; Operation handoff) |
| AD-010 | INVARIANT_FOR_LATER | invariant; documentation-only (T-031) |
| AD-011 | DEFERRED | deferred; runtime/vendor selection (§15) |
| AD-012 | CONTRACT_NOW | translated (T-011; pure kernel) |
| AD-013 | INVARIANT_FOR_LATER | invariant; retry-key carrier only (T-005) |
| AD-014 | INVARIANT_FOR_LATER | invariant; atomic-acceptance carrier only (T-031) |
| AD-015 | CONTRACT_NOW | translated (T-007, T-014, T-018; scope + AuditEvent) |
| AD-016 | CONTRACT_NOW | translated (T-004; sovereignty carrier) |
| AD-017 | CONTRACT_NOW | translated (T-009; offline bundle) |
| AD-018 | CONTRACT_NOW | translated (T-009; sync carriers) |
| AD-019 | CONTRACT_NOW | translated (T-025; projection/minimization) |
| AD-020 | CONTRACT_NOW | translated (T-025; AI projection) |
| AD-021 | CONTRACT_NOW | translated (T-025, T-029; AI-off) |
| AD-022 | CONTRACT_NOW | translated (T-005, T-022; trust carriers) |
| AD-023 | CONTRACT_NOW | translated (T-033; provider neutrality) |
| AD-024 | CONTRACT_NOW | translated (T-002; references/digests) |
| AD-025 | CONTRACT_NOW | translated (T-009, T-013; static registries) |
| AD-026 | CONTRACT_NOW | translated (T-007; telemetry-not-audit) |
| AD-027 | CONTRACT_NOW | translated (T-004; jurisdiction as profile data) |
| AD-028 | CONTRACT_NOW | translated (T-027…T-030, T-033; contract-only + gates) |
| AD-029 | CONTRACT_NOW | translated (T-004, T-013; profile extension) |
| AD-030 | CONTRACT_NOW | translated (T-001; closed forms) |
| AD-031 | CONTRACT_NOW | translated (T-001; authority) |
| AD-032 | CONTRACT_NOW | translated (T-023; authorization profile) |
| AD-033 | CONTRACT_NOW | translated (T-032; risk-based recording) |
| AD-034 | CONTRACT_NOW | translated (T-024; relationship/effective-state) |
| AD-035 | INVARIANT_FOR_LATER | invariant; async materialization (T-031) |
| AD-036 | CONTRACT_NOW | translated (T-005; semantic identity vs retry key) |
| AD-037 | CONTRACT_NOW | translated (T-011; replay) |
| AD-038 | CONTRACT_NOW | translated (T-009; version-pinned bundle/trust state) |
| AD-039 | INVARIANT_FOR_LATER | invariant; local authz cache (T-031) |
| AD-040 | CONTRACT_NOW | translated (T-008, T-019; bounds/limit codes) |
| AD-041 | INVARIANT_FOR_LATER | invariant; tiered crypto documentation (T-031) |
| AD-042 | CONTRACT_NOW | translated (T-009; disconnected sync contract) |
| AD-043 | CONTRACT_NOW | translated (T-004, T-010, T-019; SecurityExceptionRef) |
| AD-044 | CONTRACT_NOW | translated (T-015, T-032; adoption gate) |
| AD-045 | CONTRACT_NOW | translated (T-017, T-030; public code binding) |

### Matrix E risks F13-R01 through F13-R31

Controls only; final residual acceptance remains a human gate (architecture
§18.4–18.6). Staged by T-032 without residual-risk calculation. F13-R04,
F13-R08, and F13-R13 remain explicitly open at High target residual.

| Risk | Disposition | Controlling design/task |
|---|---|---|
| F13-R01 | controls translated | §10 governed registration (T-009, T-030) |
| F13-R02 | controls translated | §7.3 versioned result schemas (T-013, T-030) |
| F13-R03 | controls translated | §12 AI non-authoritative (T-025) |
| F13-R04 | controls translated + invariant | §12 policy-context provenance carrier (T-005, T-031) |
| F13-R05 | controls translated | §7.1 retry vs semantic identity (T-005) |
| F13-R06 | controls translated + invariant | §11 atomic-acceptance carrier (T-031) |
| F13-R07 | controls translated | §11 local obligation acceptance (T-031) |
| F13-R08 | controls translated | §10 bounds; §9 limit codes (T-008, T-019) |
| F13-R09 | controls translated | §8 captured-output replay (T-011) |
| F13-R10 | controls translated | §12 projection/minimization (T-025) |
| F13-R11 | controls translated | §9 obligation codes (T-023) |
| F13-R12 | controls translated + invariant | §12 local-artifact invariant (T-031) |
| F13-R13 | controls translated + invariant | §12 authz freshness/revocation invariant (T-031) |
| F13-R14 | controls translated | §12 fail-closed + SecurityExceptionRef (T-004, T-019) |
| F13-R15 | controls translated | §9.3 relationship matrix (T-024) |
| F13-R16 | controls translated | §12 minimization/redaction negatives (T-025, T-030) |
| F13-R17 | controls translated | §12 tombstone/retention/legal-hold carriers (T-004) |
| F13-R18 | controls translated | §15 no crypto service; no-network gate (T-033) |
| F13-R19 | controls translated | §10 origin-aware manifests/replay (T-009, T-030) |
| F13-R20 | controls translated | §10 no last-writer-wins; conflict evidence (T-030) |
| F13-R21 | controls translated | §17 provider-neutrality; import fixtures (T-033) |
| F13-R22 | controls translated | §5/§15 contract-only; boundary gate (T-033) |
| F13-R23 | controls translated | §10 bounded payloads; cold-reference carriers (T-002) |
| F13-R24 | controls translated | §10 version-pinned bundles + trust state (T-009) |
| F13-R25 | controls translated | §7.1 trusted-time provenance carrier (T-005) |
| F13-R26 | controls translated | §9 scope mismatch/no-disclosure (T-018, T-020) |
| F13-R27 | controls translated | §9.2 closed public registry; baseline gate (T-017, T-015) |
| F13-R28 | controls translated | §11 six-scope AuditEvent regression (T-014, T-029) |
| F13-R29 | controls translated | §7.1 DecisionRecord naming; baseline checks (T-002, T-015) |
| F13-R30 | controls translated | §10 tiered profiles; minimal-path budget (T-009) |
| F13-R31 | controls translated | §16 one gate, no registry framework (T-032) |

### Conformance IDs

All conformance IDs are class `CONTRACT_NOW` and mapped to fixtures by T-028; the
compound-family suffixes are covered by T-029/T-030 and the coverage matrix.

Scenario families (F13-CF-01 through F13-CF-28): F13-CF-01, F13-CF-02,
F13-CF-03, F13-CF-04, F13-CF-05, F13-CF-06, F13-CF-07, F13-CF-08, F13-CF-09,
F13-CF-10, F13-CF-11, F13-CF-12, F13-CF-13, F13-CF-14, F13-CF-15, F13-CF-16,
F13-CF-17, F13-CF-18, F13-CF-19, F13-CF-20, F13-CF-21, F13-CF-22, F13-CF-23,
F13-CF-24, F13-CF-25, F13-CF-26, F13-CF-27, F13-CF-28.

Scope cases: F13-SCOPE-01, F13-SCOPE-02, F13-SCOPE-03, F13-SCOPE-04,
F13-SCOPE-05, F13-SCOPE-06, F13-SCOPE-07, F13-SCOPE-08, F13-SCOPE-09,
F13-SCOPE-10, F13-SCOPE-11.

Evaluation cases: F13-EVAL-01, F13-EVAL-02, F13-EVAL-03, F13-EVAL-04,
F13-EVAL-05, F13-EVAL-06, F13-EVAL-07, F13-EVAL-08.

Security cases: F13-SEC-01, F13-SEC-02, F13-SEC-03, F13-SEC-04, F13-SEC-05,
F13-SEC-06, F13-SEC-07, F13-SEC-08.

Trust cases: F13-TRUST-01, F13-TRUST-02, F13-TRUST-03, F13-TRUST-04.

Compatibility cases: F13-COMPAT-01, F13-COMPAT-02, F13-COMPAT-03, F13-COMPAT-04,
F13-COMPAT-05, F13-COMPAT-06, F13-COMPAT-07, F13-COMPAT-08, F13-COMPAT-09.
