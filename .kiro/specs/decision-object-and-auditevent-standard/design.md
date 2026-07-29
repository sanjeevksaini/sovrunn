# Design Document

Feature: FEATURE-0013 — Decision Record and AuditEvent Standard
Phase: Phase 2 — Reuse-First PaaS Fabric Foundation
Stage: Design

## Metadata

| Field | Value |
|---|---|
| Feature ID | FEATURE-0013 |
| Feature title | Decision Record and AuditEvent Standard |
| Phase | Phase 2 — Reuse-First PaaS Fabric Foundation |
| Spec stage | Design |
| Branch | feature-0013-decision-record-and-auditevent-standard |
| Depends on | FEATURE-0012 (API, Resource Naming, Status, and Validation Standard) |
| Controlling handoff | ADH-2026-017 (Approved) |
| Canonical architecture | docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md |
| Approved requirements | .kiro/specs/decision-object-and-auditevent-standard/requirements.md |
| Go module | github.com/sanjeevksaini/sovrunn |
| Status | Draft — pending review |

This design is a clean regeneration from the approved consolidated architecture
and the approved requirements. It does not reuse, patch, continue, or copy from
`design.superseded-patched-draft.md`, and it does not use ADH-2026-014/015/016
or any prior FEATURE-0013 design review or revision prompt as a semantic source.
Per architecture section 27.9(14), architecture is the sole source of truth.

## 1. Overview

FEATURE-0013 delivers a **contract-only** decision and audit vocabulary for
Sovrunn. It defines the immutable `DecisionRecord` envelope, the versioned
`DecisionProfile` extension mechanism, the atomic `EvaluationResult` contract,
the FEATURE-0013 `AuditEvent` payload/linkage extension, deterministic bounded
composition semantics, structural trust carriers, sovereign portability
carriers, the closed public violation-code binding, and executable conformance
fixtures.

The design resolves only the representation and implementation mechanics that
architecture section 27.8 explicitly delegates to design: Go field types, JSON
Schema composition, package/import boundaries, bounded in-memory graph
representation, validation-pass organization, pure-function interfaces, and
physical fixture file format. It reproduces every AD-001 through AD-045
implementation class exactly (architecture section 19) and produces
implementation artifacts only for `CONTRACT_NOW` decisions. `INVARIANT_FOR_LATER`
and `DEFERRED` items receive documentation and traceability placeholders only
and generate no runtime.

The design creates no runtime service, persistence, queue, worker, workflow,
network client, cryptographic execution, AI runtime, adapter interface, policy
engine, or production registry. All behavior is expressed as canonical schemas,
Go value types, pure/bounded in-memory validation and composition functions, a
declarative JSON registry bundle, and deterministic fixtures.

`DecisionRequest` is calling-domain-owned. FEATURE-0013 defines no shared,
canonical, or common `DecisionRequest` schema, Go type, registry entry, resource,
persistence model, or public API resource; the request shape belongs to each
calling domain. Completed synchronous handling returns a `DecisionRecord`
value/reference; accepted asynchronous handling returns the existing FEATURE-0012
`Operation` (`LongRunningOperation`) value/reference; rejection returns inherited
Problem Details. No FEATURE-0013 pending-decision response envelope or non-final
response contract is defined.

The canonical `DecisionRecord` and `AuditEvent` objects are append-only, per
architecture section 6.7: **DecisionRecord and AuditEvent are append-only.**
Correction, supersession, revocation, retention, and lawful erasure never
rewrite the canonical immutable record in place. The only alternatives are
separately controlled payload custody, key destruction, or a linked
tombstone/redaction record, and FEATURE-0013 implements no erasure, key
destruction, or cryptographic process.

### 1.1 Design goals

1. Bind every FEATURE-0013 object to exactly one inherited FEATURE-0012 resource
   profile without creating a parallel grammar (F13-OBJ-001…008).
2. Make `metadata.scopeRef` the sole scope authority with canonical absent
   Platform serialization and the six-value FEATURE-0012 `ScopeKind` set
   (F13-SCOPE-001…008).
3. Represent the closed decision-form, authority, adjudication, audit-outcome,
   relationship, sensitivity, and public violation-code vocabularies as
   FEATURE-0013-owned Go constants derived from architecture, with no new value.
4. Provide deterministic, pure, bounded in-memory validation and composition
   sufficient to prove contract semantics through fixtures.
5. Reuse FEATURE-0012 decode/pre-scan, Problem Details, TypeBinding, and
   baseline-manifest machinery exactly; never fork it.

### 1.2 Inherited limit hierarchy

All FEATURE-0013 limit resolution obeys the exact inherited-ceiling hierarchy
fixed by architecture; the design applies it without reinterpretation:

- FEATURE-0012 platform/schema ceilings are absolute outer bounds.
- DecisionProfile ceilings must remain within those inherited ceilings.
- Evaluator-specific ceilings may only narrow the DecisionProfile.
- The effective limit is the most restrictive applicable value.

Missing, unknown, incomparable, or above-ceiling mandatory limits fail closed.
This hierarchy is directional (inherited outer bound → profile → evaluator) and
is not a symmetric "most restrictive wins" comparison between independent peers.

## 2. Controlling inputs

Authority order (highest first), per the regeneration prompt and ADH-2026-017:

1. `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
   main body — controls all semantics; sections 27.8 and 27.9 are closed and
   normative.
2. `docs/reviews/architecture-decision-handoffs/ADH-2026-017-feature-0013-consolidated-architecture.md`
   — approval and exact content binding (`architecture_content_sha256`,
   `dependency_content_sha256`). Sole active controlling FEATURE-0013 handoff.
3. `.kiro/specs/decision-object-and-auditevent-standard/requirements.md`
   — approved stage input; architecture wins on any conflict.
4. FEATURE-0011 and FEATURE-0012 locked dependency files referenced by
   ADH-2026-017 `dependency_content_sha256` — control inherited contracts.
5. Repo engineering guardrails (Go 1.22, stdlib + `gopkg.in/yaml.v3` only, no
   new third-party runtime dependency) where they do not conflict with the
   FEATURE-0013 architecture.

Historical only (not used as design source): ADH-2026-014/015/016; old
FEATURE-0013 requirements/design drafts, reviews, and revision prompts;
`design.superseded-patched-draft.md`; generated reports.

Inherited FEATURE-0012 sources materially used by this design:
`internal/apimeta` (TypeMeta, ObjectMeta, TypedRef, ScopeRef, OwnerRef,
`ScopeKind`, `Profile`, `DataClassification`, `NormalizeScope`,
`CanonicalScopeIdentity`), `internal/apivalid` (`StrictDecode`, `Limits`,
`DecodeMode`, `FieldPolicy`, JSON-token duplicate scan, pre-normalization
`yaml.Node` duplicate rejection), `internal/apiproblem` (`Problem`, `Violation`),
`internal/apischema` (`TypeBinding`, `VerifyGoTypeAgainstSchema`, diff,
subset validation, annotations, baseline support), `internal/apiconform`
(`TypeBindings`, contract types, fixtures, executable fitness checks),
`api/schemas/_common/*`, `api/schemas/audit-event.json`, and
`api/schemas/baseline/BASELINE_MANIFEST.json` / `BASELINE_APPROVALS.json`.

## 3. Feature-level reuse summary

Reproduced verbatim from the approved architecture section 14 and requirements;
FEATURE-0011 canonical columns. Field meanings are owned by
`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`. Design introduces no new
disposition.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 API/resource grammar | Extend | Reuse every shared profile, metadata, reference, boundary, validation, Problem Details, schema, and compatibility primitive; add only decision/audit domain semantics and approved per-kind constraints. | Approved | ADH-2026-017; ADH-2026-012; RFC-0022 |
| DecisionRecord, DecisionProfile, EvaluationResult, and AuditEvent domain contract | Build | No mature external model supplies Sovrunn's complete sovereign, provider-neutral authority, projection, scope, audit-obligation, and downstream-adoption semantics. The common contract is Sovrunn differentiation; external concepts remain inputs rather than native core types. | Approved | ADH-2026-017; DEC-0026; RFC-0023 |
| Bounded decision-requirements and composition concepts | Reuse | Reuse applicable OMG DMN decision-requirements concepts without selecting, embedding, or recreating a DMN engine or workflow language. | Approved | ADH-2026-017; RFC-0023 |
| Audit transport and operational correlation | Reuse | CloudEvents may be an optional transport mapping and OpenTelemetry may carry operational correlation; neither becomes the canonical record or audit authority. | Approved | ADH-2026-017; RFC-0023 |
| Supply-chain evidence references | Reuse | Reuse DSSE/in-toto and SPDX/CycloneDX concepts by typed reference without duplicating their schemas or implementing signing. | Approved | ADH-2026-017; ADR-F13-002 |
| Canonicalization, digest, signing, and cryptographic verification | Reuse | Preserve algorithm-agile carriers and evaluate mature standards later; no algorithm, covered-field set, product, or service is selected in FEATURE-0013. | Deferred | ADH-2026-017; ADR-F13-002 |
| Policy/evaluator and workflow runtimes | Wrap | CEL, OPA, Cedar, durable-execution, and workflow products are later candidates behind owning-feature adapters or ports; FEATURE-0013 defines no wrapper or runtime implementation. | Deferred | ADH-2026-017; DEC-0036 |

## 4. Resolved implementation decisions

These are the design-owned representation/mechanics decisions delegated by
architecture section 27.8. Each closes an open design question (DQ) or a section
27.9 closure control. None changes a closed architecture semantic.

- **RID-01 (DQ-13-03, closure 7): Package boundaries.** Canonical FEATURE-0013
  domain value types live in a new `internal/decision` package tree owned by
  FEATURE-0013. `internal/apiconform` remains limited to TypeBinding
  registration, conformance adapters, and executable checks; it does not become
  the canonical domain owner and does not receive divergent copies of
  FEATURE-0012 base fields. Class: `CONTRACT_NOW`.

- **RID-02 (DQ-13-01): Go representation.** Every canonical object is a plain Go
  struct that anonymously embeds `apimeta.TypeMeta` and carries `apimeta.ObjectMeta`
  for `metadata`, mirroring FEATURE-0012 `internal/apiconform/contracts.go`.
  Optional singular references use `*apimeta.TypedRef` / `omitempty`; required
  collections use non-pointer slices; closed vocabularies are string-backed named
  types with `Valid()` methods. Class: `CONTRACT_NOW`.

- **RID-03 (DQ-13-02, closure 6+8): Bundle and schema layout.** The canonical
  registry is a single declarative JSON bundle schema
  (`api/schemas/decision-profile-bundle.json`) governed by JSON Schema, with
  per-object schemas for `DecisionRecord`, `DecisionProfile`, and
  `EvaluationResult`, plus `_common` sub-schemas. Schema evolution, AuditEvent
  compatibility, and schema-diff evidence reuse the exact FEATURE-0012 baseline
  mechanism (`api/schemas/baseline/BASELINE_MANIFEST.json`,
  `BASELINE_APPROVALS.json`, digest recomputation, approval-controlled updates,
  protected review). No parallel `api/schemas/SCHEMA_BASELINE_MANIFEST.json`,
  `api/schemas/diffs/`, or `api/schemas/approvals/` hierarchy is introduced.
  Class: `CONTRACT_NOW`.

- **RID-04 (DQ-13-04, closure 3+4): Graph representation.** The bounded
  decision-requirements DAG is an in-memory adjacency structure built by
  `internal/decision/graph` from an **already-resolved** graph definition; the
  builder performs no registry or version lookup. Graph reference/version
  compatibility is validated earlier during bundle resolution
  (`internal/decision/bundle`). Class: `CONTRACT_NOW`.

- **RID-05 (DQ-13-05): Validation-pass order.** Validators run in a fixed,
  documented order (section 9). Class: `CONTRACT_NOW`.

- **RID-06 (DQ-13-06): Composition interface.** Composition strategies are pure
  functions `func(StrategyInput) (StrategyOutput, *apiproblem.Problem)` with no
  I/O; deterministic replay consumes captured evaluator output and never reruns
  an evaluator. Class: `CONTRACT_NOW`.

- **RID-07 (DQ-13-07, DQ-13-09): Fixture layout.** Positive fixtures live under
  `tests/conformance/fixtures/decision/positive/`; negative fixtures under
  `tests/conformance/fixtures/negative/decision/` (extending the existing
  `apiconform.ConformanceNegativeFixturesDir` convention). A single
  coverage-matrix source (`tests/conformance/decision_coverage_matrix.go`)
  maps every scenario ID to fixtures. Class: `CONTRACT_NOW`.

- **RID-08 (closure 9): TrustCarrier empty-state rule (design-delegated
  choice).** Architecture delegates one documented rule for an optional,
  present-but-structurally-empty carrier: reject-as-invalid **or**
  normalize-to-absent. This design selects **reject-as-invalid**: an optional
  trust/integrity carrier that is present but structurally empty is a malformed
  input and fails with the appropriate existing `DECISION_TRUST_*` /
  `DECISION_EVALUATION_*` code. An absent optional carrier is valid. A
  present-but-empty carrier never satisfies provenance, required trust, identity
  generation, or compatibility evidence. A required carrier that is absent or
  structurally empty fails closed. This is a mechanics choice explicitly
  permitted by closure 9; it adds no new code or vocabulary. Class: `CONTRACT_NOW`.

- **RID-09 (closure 13): Decode/pre-scan reuse.** FEATURE-0013 decoding reuses
  `apivalid.StrictDecode`, which performs JSON-token duplicate/parallel-key
  detection on JSON input and pre-normalization `yaml.Node` AST duplicate
  rejection on YAML input. FEATURE-0013 adds no local decoder and makes no
  zero-allocation claim for the YAML AST path. Class: `CONTRACT_NOW`.

- **RID-10 (closure 2): Strategy fail-open boundary.** A composition strategy
  may describe missing/timeout/conflict/error handling but is never an
  independent fail-open authority. Strategy-level fail-open is effective only
  when the governing profile carries a valid `SecurityExceptionRef`; otherwise
  the strategy fails closed. Class: `CONTRACT_NOW`.

## 5. Package / import architecture

New FEATURE-0013 packages under `internal/decision` form a small, acyclic DAG.
Arrows point in the allowed import direction. FEATURE-0012 packages are reused
unchanged.

```
apimeta  apiproblem  apivalid  apischema        (FEATURE-0012, reused unchanged)
   ^          ^          ^          ^
   |          |          |          |
internal/decision            (domain value types + closed vocabularies)
   ^        ^        ^
   |        |        |
   |        |        +-- internal/decision/graph   (bounded in-memory DAG + accounting)
   |        +----------- internal/decision/bundle  (declarative JSON bundle, BundleView,
   |                                                 (id,version) registry, graph resolution)
   +-------------------- internal/decision/compose (pure strategies + deterministic replay)
                                 ^        ^
                                 |        |
                    internal/decision/validate  (validation passes -> apiproblem.Problem)
                                 ^
                                 |
             internal/apiconform (TypeBinding registration + conformance adapters + checks)
                                 ^
                                 |
             tests/conformance   (fixtures + coverage matrix)
```

Import rules (verified by the anti-overengineering gate, section 27.7 /
`imports_test.go` pattern):

- `internal/decision` (domain) imports only `apimeta` and stdlib. It owns the
  canonical AuditEvent payload/linkage extension value types (closure 7).
- `internal/decision/bundle` imports `decision`, `apivalid`, `apischema`,
  `apiproblem`. It owns `(id, version)` keyed registries, duplicate detection,
  and graph reference/version resolution (closure 4, 8).
- `internal/decision/graph` imports `decision`, `apiproblem`. It receives a
  resolved graph definition; it performs no registry/version lookup (closure 4).
- `internal/decision/compose` imports `decision`, `graph`, `apiproblem`. Pure
  functions only.
- `internal/decision/validate` imports `decision`, `bundle`, `graph`, `compose`,
  `apimeta`, `apivalid`, `apiproblem`. It emits `apiproblem.Problem` with the
  closed `DECISION_*` violation codes.
- `internal/apiconform` gains only FEATURE-0013 `TypeBinding` entries and
  executable checks/fixtures; it never becomes the canonical domain owner and
  must not import `decision/validate` in a way that creates a cycle (checks call
  exported validators one-directionally).
- No FEATURE-0013 package imports a network, database, queue, workflow,
  cryptographic, or provider package (none exists; none is created).

Each package satisfies the section 27.4 necessity test: it maps to a specific
`CONTRACT_NOW` requirement (domain types → AC-13.1/13.3/13.4; bundle → AC-13.5;
graph → AC-13.6; compose → AC-13.6/13.9 kernel; validate → AC-13.13 error
mapping), reuses a FEATURE-0012 primitive where one exists (decode, Problem,
TypeBinding, baseline), and can be deleted only by removing the fixtures that
prove its `CONTRACT_NOW` requirement.

## 6. Files to create / change

All entries are `CONTRACT_NOW` unless marked otherwise. No file is created for
an `INVARIANT_FOR_LATER` or `DEFERRED` decision except documentation/traceability.

### 6.1 Schemas (create)

| Path | Purpose | Profile / notes |
|---|---|---|
| `api/schemas/decision-record.json` | Canonical `DecisionRecord` | `x-sovrunn-profile: ImmutableRecord`, `x-sovrunn-boundary: governance-only`; non-Platform `metadata.scopeRef` requires `uid` (F13-SCOPE-002) |
| `api/schemas/decision-profile.json` | Canonical `DecisionProfile` | `x-sovrunn-profile: VersionedDefinition`; definition content in `spec` |
| `api/schemas/evaluation-result.json` | Canonical `EvaluationResult` | `x-sovrunn-profile: EmbeddedValue` (in-record) / `TransientRequestResult` (exchanged) |
| `api/schemas/decision-profile-bundle.json` | Declarative registry bundle | Canonical semantic-registry representation; `(id,version)` keyed entries + algorithm-agile trust carriers |
| `api/schemas/_common/decision-linkage.json` | AuditEvent payload/linkage extension + decision↔audit correlation | Additive optional; referenced by `audit-event.json` and `decision-record.json` |
| `api/schemas/_common/decision-profile-ref.json` | `profileRef` (name + semver) | Typed profile reference |
| `api/schemas/_common/trust-carrier.json` | Algorithm-agile canonicalization/digest/signature/trust-state carrier | Structural only; no algorithm selected |
| `api/schemas/_common/security-exception-ref.json` | Bounded `SecurityExceptionRef` approval-evidence reference | apiVersion/kind/name/uid + bounded metadata (closure 1) |
| `api/schemas/_common/semantic-decision-identity.json` | Complete semantic decision-identity descriptor | Distinct from retry idempotency key |
| `api/schemas/_common/decision-relationship.json` | `CORRECTS`/`SUPERSEDES`/`REVOKES` typed relationship | Bounded, cycle-free |
| `api/schemas/_common/sensitivity.json` | `PUBLIC<INTERNAL<CONFIDENTIAL<RESTRICTED` + profile-bound rules | Separate from `DataClassification` |
| `api/schemas/_common/decision-graph.json` | Bounded graph definition (nodes/edges/bounds) | Node/edge IDs, edge id distinct from (from,to) (closure 3) |

### 6.2 Schemas (change through the FEATURE-0012 baseline mechanism)

| Path | Change | Delta class (locked reconciliation) |
|---|---|---|
| `api/schemas/audit-event.json` | Expand `x-sovrunn-allowed-scopes` from `[Organization]` to the six governance scopes; add optional additive decision-linkage properties via `$ref: _common/decision-linkage.json` | expanded + added (approved, ADH-2026-017) |
| `api/schemas/baseline/BASELINE_MANIFEST.json` | Recompute digests for changed/added files using the existing mechanism | reuse |
| `api/schemas/baseline/BASELINE_APPROVALS.json` | Record the approval-controlled baseline update for the AuditEvent expansion and new schemas | reuse |

The base FEATURE-0012 `AuditEvent` envelope, `ImmutableRecord` shape, and all
existing base fields are preserved unchanged (no field removed, renamed, or
reinterpreted). No parallel baseline/diff/approval hierarchy is created
(closure 6).

### 6.3 Go packages (create)

| Path | Contents |
|---|---|
| `internal/decision/doc.go` | Package doc + ownership note (closure 7) |
| `internal/decision/vocab.go` | Closed vocabularies: `DecisionForm`, `DecisionAuthority`, `AdjudicationOutcome`, `AuditOutcome`, `RelationshipKind`, `Sensitivity`, `Finality`, `EffectiveState`, each with `Valid()` and stable-order `All…()` |
| `internal/decision/decisionrecord.go` | `DecisionRecord`, `DecisionResultBody`, `DecisionRationale`, `Obligation`, `Correlation` |
| `internal/decision/evaluationresult.go` | `EvaluationResult`, evaluator identity/provenance carriers |
| `internal/decision/profile.go` | `DecisionProfile`, `ProfileSpec`, limit descriptors, projection specs, `SecurityExceptionRef` |
| `internal/decision/auditevent_ext.go` | FEATURE-0013 AuditEvent payload/linkage extension value types (closure 7) |
| `internal/decision/trust.go` | Algorithm-agile trust/integrity carriers + structural trust-state enum |
| `internal/decision/identity.go` | Retry idempotency key + semantic decision-identity descriptor |
| `internal/decision/sensitivity.go` | Sensitivity vocabulary + `DataClassification`→sensitivity-floor map |
| `internal/decision/bundle/bundle.go` | `Bundle`, `BundleView`, `(id,version)` registries, duplicate detection, graph reference/version resolution |
| `internal/decision/bundle/load.go` | Bundle decode via `apivalid.StrictDecode` (closure 13) |
| `internal/decision/graph/graph.go` | In-memory DAG builder from resolved definition; deterministic accounting (closure 3, 4) |
| `internal/decision/compose/strategy.go` | Pure strategy interface, registered strategy metadata, fail-open-only-via-`SecurityExceptionRef` gate (closure 2) |
| `internal/decision/compose/replay.go` | Deterministic replay reference kernel (AD-012, AD-037) |
| `internal/decision/validate/*.go` | Validation passes (section 9), each emitting `apiproblem.Problem` with `DECISION_*` codes |
| `internal/decision/validate/codes.go` | FEATURE-0013 closed `violations[].code` constants (mirrors architecture section 15.2; no new code) |

### 6.4 Conformance / apiconform (change)

| Path | Change |
|---|---|
| `internal/apiconform/bindings.go` | Append FEATURE-0013 `TypeBinding` entries for the new schemas (registration only) |
| `internal/apiconform/decision_checks.go` (new) | Executable FEATURE-0013 conformance checks calling exported `decision/validate` functions |
| `tests/conformance/fixtures/decision/positive/*.json` | Positive fixtures |
| `tests/conformance/fixtures/negative/decision/*.json` | Negative fixtures |
| `tests/conformance/decision_coverage_matrix.go` (new) | Scenario-ID → fixture coverage matrix |

### 6.5 Documentation / traceability placeholders only (no runtime)

| Path | Purpose | Class |
|---|---|---|
| `docs/reviews/feature-gates/FEATURE-0013-matrix-e-worksheet.md` (staged) | Matrix E worksheet: staged risk IDs, controls, evidence pointers, blank human-only disposition fields (closure 12) | Governance-only |
| `internal/decision/doc.go` invariant notes | Record `INVARIANT_FOR_LATER` obligations (AD-007/008/010/013/014/035/039/041) and `DEFERRED` items (AD-011, ADR-F13-002) as documentation | Documentation |

## 7. Data models

All models embed FEATURE-0012 primitives and use its serialization rules. Closed
vocabularies reproduce architecture values exactly; none is added.

### 7.1 DecisionRecord (`ImmutableRecord`)

```go
// internal/decision
type DecisionRecord struct {
    apimeta.TypeMeta                 // apiVersion=governance.sovrunn.io/v1alpha1, kind=DecisionRecord
    Metadata apimeta.ObjectMeta `json:"metadata"` // scopeRef is the sole scope authority
    Record   DecisionBody       `json:"record"`
}

type DecisionBody struct {
    ProfileRef        ProfileRef            `json:"profileRef"`
    Form              DecisionForm          `json:"form"`
    Authority         DecisionAuthority     `json:"authority"`
    Purpose           string                `json:"purpose"`
    Adjudication      *AdjudicationOutcome  `json:"adjudication,omitempty"` // only for adjudication facets
    SubjectRefs       []apimeta.TypedRef    `json:"subjectRefs,omitempty"`
    EvaluationResults []EvaluationResult    `json:"evaluationResults,omitempty"` // EmbeddedValue
    Composition       *CompositionRef       `json:"composition,omitempty"`
    Result            DecisionResultBody    `json:"result"`
    Relationship      *DecisionRelationship `json:"relationship,omitempty"`
    RetryKey          string                `json:"retryKey,omitempty"`   // idempotency (F13-SCALE-001)
    SemanticIdentity  SemanticDecisionIdentity `json:"semanticIdentity"`  // F13-SCALE-002
    Trust             *TrustCarrier         `json:"trust,omitempty"`
    Sovereignty       *SovereigntyCarrier   `json:"sovereignty,omitempty"`
    Correlation       Correlation           `json:"correlation"` // auditEventRefs, operationRef, traceRef
    Finality          Finality              `json:"finality"`    // always FINAL when persisted
}
```

- `Finality` is always `FINAL` for a persisted record (F13-IMMUT-001); pending/
  running/retrying belong to `Operation` only.
- `Correlation.AuditEventRefs` links to ≥1 preallocated `AuditEvent` identity
  (F13-AUDIT-001, F13-AUDIT-006 carrier).
- Non-Platform `metadata.scopeRef` requires `uid` (F13-SCOPE-002, locked
  reconciliation delta).

**Immutability (append-only).** The canonical `DecisionRecord` and `AuditEvent`
objects are append-only (architecture section 6.7): **DecisionRecord and
AuditEvent are append-only.** Correction, supersession, revocation, retention,
and lawful erasure never rewrite the canonical immutable record in place. The
only alternatives are separately controlled payload custody, key destruction, or
a linked tombstone/redaction record, and FEATURE-0013 implements no erasure, key
destruction, or cryptographic process (F13-IMMUT-001/002/003/007, AD-004).

### 7.2 EvaluationResult (`EmbeddedValue` / `TransientRequestResult`)

Captures one registered evaluator's finding: evaluator identity+version, one
bounded input-snapshot reference or opaque integrity carrier, one result, one
timing envelope, explicit success/non-match/indeterminate/timeout/error
semantics, execution config, evaluation time, trust boundary, safety/policy
filters, and structural trust state (F13-EVAL-001/002). Scope is derived from the
containing record; evaluator-supplied scope metadata is evidence only
(F13-EVAL-003/004). No independent durable lifecycle (F13-OBJ-004).

### 7.3 DecisionProfile (`VersionedDefinition`)

`spec` carries the full profile declaration required by F13-PROF-001: name,
family, semver, owner, lifecycle status, primary form, secondary facets, allowed
authority levels, input/typed-result/rationale/obligation schemas, accepted
evaluation types and composition strategies, scope/actor/subject/policy/evidence/
audit semantics, determinism/failure/timeout/insufficient-evidence behavior,
validity/expiry/correction/supersession/revocation rules, sensitivity/residency/
retention/legal-hold/projection profiles, idempotency identity, bounded limits,
and compatibility policy with fixtures/owner/reassessment triggers. Profile
limit descriptors are resolved under the inherited limit hierarchy of section
1.2: FEATURE-0012 platform/schema ceilings are absolute outer bounds;
DecisionProfile ceilings must remain within those inherited ceilings;
evaluator-specific ceilings may only narrow the DecisionProfile; the effective
limit is the most restrictive applicable value. Security-class
profiles require fail-closed or `REQUIRES_APPROVAL`; fail-open requires a valid
`SecurityExceptionRef` (F13-PROF-002, closure 1).

### 7.4 AuditEvent FEATURE-0013 extension

The FEATURE-0012 `AuditEvent` base (`internal/apiconform` contract type +
`api/schemas/audit-event.json`) is preserved. FEATURE-0013 adds
`internal/decision` value types for the decision-linkage extension
(decisionRef, auditObligation identity, correlation) and expands allowed scopes
to six. Audit outcomes remain the closed `Succeeded`/`Denied`/`Failed`
vocabulary (F13-AUDIT-002). Base fields are not copied into a divergent parallel
type (closure 7).

### 7.5 Closed vocabularies (reproduced, not extended)

| Go type | Values (exact) | Requirement |
|---|---|---|
| `DecisionForm` | ADJUDICATION, SELECTION, RANKING, CLASSIFICATION, RESOLUTION, ALLOCATION, PLAN, ASSESSMENT, RECOMMENDATION | F13-FORM-001 |
| `DecisionAuthority` | AUTHORITATIVE, ADVISORY, RECOMMENDATION, SIMULATION | F13-AUTH-001 |
| `AdjudicationOutcome` | ALLOWED, DENIED, REQUIRES_APPROVAL | F13-ADJ-001 |
| `AuditOutcome` | Succeeded, Denied, Failed | F13-AUDIT-002 |
| `RelationshipKind` | CORRECTS, SUPERSEDES, REVOKES | F13-IMMUT-002 |
| `Sensitivity` | PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED | F13-SEC-001 |
| `Finality` | FINAL | F13-IMMUT-001 |
| `EffectiveState` (projection only) | EFFECTIVE, SUPERSEDED, REVOKED | F13-IMMUT-003 |

### 7.6 Requiredness enumeration (closure 11)

For every public/externally-exchanged contract, requiredness is enumerated per
field kind so that Go tags, JSON Schema `required`/`minItems`/`nullable`,
loaders, validators, fixtures, and error pointers agree.

| Field kind | Rule |
|---|---|
| string | Required → non-empty (`minLength: 1`), no Go `omitempty`; optional → Go `omitempty`, absent = unset |
| boolean | Modeled as explicit required bool, or `*bool` when tri-state (true/false/unset) is contract-meaningful; documented per field |
| number | Required → present + within profile bounds; optional → `omitempty` |
| object | Required → present, validated recursively; optional singular nested object → pointer + `omitempty` |
| pointer (typed ref) | Required singular ref → non-pointer struct; optional singular ref → `*apimeta.TypedRef` + `omitempty`; **null is not a second canonical form** — absence is canonical unless a future approved nullable union is declared |
| map | Required → present; optional → `omitempty`; keys validated; `(id,version)` maps use the composite-key rule (closure 8) |
| slice / array | Required-present slice → declared **empty-allowed** or **non-empty** (`minItems: 1`) explicitly per field; optional → `omitempty`, absent ≠ empty |

Concrete applications: `DecisionBody.SubjectRefs` is optional (empty allowed
where the profile permits); `Correlation.AuditEventRefs` is required non-empty
(`minItems: 1`, F13-AUDIT-001); `metadata.scopeRef` optional (absent = Platform)
but `uid` required when present and non-Platform (F13-SCOPE-002); optional trust
carriers follow RID-08 (present-empty rejected).

### 7.7 SecurityExceptionRef (fail-open evidence carrier, closure 1)

`SecurityExceptionRef` is a bounded structural approval-evidence reference carried
inside a `DecisionProfile` failure-behavior declaration (architecture sections 6.8
and 27.9(1)). It is evidence only: it never approves, issues, revokes, executes,
stores, authorizes, or workflows an exception.

```go
// internal/decision
type SecurityExceptionRef struct {
    ApprovalRef           apimeta.TypedRef  `json:"approvalRef"`            // apiVersion/kind/name/uid; uid mandatory
    ExceptionScopeRef     *apimeta.ScopeRef `json:"exceptionScopeRef,omitempty"` // exception-evidence scope only; absent = Platform where allowed
    OwnerRef              apimeta.TypedRef  `json:"ownerRef"`               // uid required
    ApprovingAuthorityRef apimeta.TypedRef  `json:"approvingAuthorityRef"`  // uid required
    CompensatingControls  []string          `json:"compensatingControls"`   // required non-empty; bounded non-empty control identifiers only
    EffectiveFrom         string            `json:"effectiveFrom"`          // UTC RFC3339
    EffectiveUntil        string            `json:"effectiveUntil"`         // UTC RFC3339; after EffectiveFrom and >= captured/effective decision time
    AuditTreatment        string            `json:"auditTreatment"`         // bounded registry identifier; structural only
    ReassessmentTrigger   string            `json:"reassessmentTrigger"`    // bounded registry identifier
    CoveredFailureModes   []string          `json:"coveredFailureModes"`    // non-empty; profile-declared failure modes only
    Purpose               string            `json:"purpose"`                // non-empty bounded string
}
```

- `CompensatingControls` is a required, non-empty `[]string` of bounded, non-empty
  control identifiers. FEATURE-0013 introduces no control object, control schema,
  or `ControlRef`/control-ref type; the field is exactly a bounded string
  identifier list and nothing else.
- `ExceptionScopeRef` is **exception-evidence scope only**. It never overrides,
  duplicates, replaces, widens, or becomes an alternative to the canonical
  `DecisionRecord`/`AuditEvent` `metadata.scopeRef`, which remains the sole scope
  authority (§1.1, §9.1 Scope pass). It exists solely to prove that the approved
  exception evidence is compatible with the governing profile/record scope. It
  follows the same canonical Platform absent/nil convention for the exception
  artifact where Platform is allowed (absent = Platform), per the section 7.6
  scope rule. It is validated against the governing profile/record scope for
  compatibility only.

**Validation ownership.** All checks belong to profile validation and fail closed
with the existing `DECISION_PROFILE_SCHEMA_INVALID` code at exact RFC 6901 JSON
Pointer paths (no new code is introduced). The container path
`/spec/failureBehavior/securityExceptionRef` is the profile-schema-defined anchor;
the field pointers below are the deterministic pattern:

| Failure | `violations[].code` | JSON Pointer |
|---|---|---|
| Malformed / missing `SecurityExceptionRef` when fail-open is declared | `DECISION_PROFILE_SCHEMA_INVALID` | `/spec/failureBehavior/securityExceptionRef` |
| Missing `approvalRef.uid` | `DECISION_PROFILE_SCHEMA_INVALID` | `/spec/failureBehavior/securityExceptionRef/approvalRef/uid` |
| Owner mismatch vs governing profile/record | `DECISION_PROFILE_SCHEMA_INVALID` | `/spec/failureBehavior/securityExceptionRef/ownerRef` |
| Scope mismatch vs governing profile/record scope | `DECISION_PROFILE_SCHEMA_INVALID` | `/spec/failureBehavior/securityExceptionRef/exceptionScopeRef` |
| Expired / non-positive interval (`effectiveUntil` ≤ `effectiveFrom`, or before captured/effective decision time) | `DECISION_PROFILE_SCHEMA_INVALID` | `/spec/failureBehavior/securityExceptionRef/effectiveUntil` |
| Empty `compensatingControls` | `DECISION_PROFILE_SCHEMA_INVALID` | `/spec/failureBehavior/securityExceptionRef/compensatingControls` |
| Empty `coveredFailureModes` | `DECISION_PROFILE_SCHEMA_INVALID` | `/spec/failureBehavior/securityExceptionRef/coveredFailureModes` |
| Unknown covered mode (not profile-declared) | `DECISION_PROFILE_SCHEMA_INVALID` | `/spec/failureBehavior/securityExceptionRef/coveredFailureModes/{index}` |

- `SecurityExceptionRef` is structural evidence only; it never approves, issues,
  revokes, executes, stores, or workflows exceptions (AD-043, F13-PROF-002).
- Strategy-level fail-open (section 8 / RID-10) is valid **only** through a valid
  governing-profile `SecurityExceptionRef`; otherwise the strategy fails closed
  (closure 2).

### 7.8 Composition graph model (closure 3/4)

The bounded in-memory composition graph is modeled in package
`internal/decision/graph` as the concrete types below (referred to in prose as
`GraphDefinition`, `GraphNode`, `GraphEdge`, `ParallelGroup`, `GraphBounds`; the
package-local names avoid stutter and match the `graph.Definition`/`graph.Bounds`
signatures in section 8). Orientation is dependency → dependent.

```go
// internal/decision/graph
type Definition struct { // GraphDefinition
    Ref            string          `json:"ref"`                       // graph identity/reference; required non-empty
    Version        string          `json:"version"`                   // (ref,version) resolved in bundle.ResolveGraph before BuildGraph
    Nodes          []Node          `json:"nodes"`                     // required non-empty
    Edges          []Edge          `json:"edges"`                     // required-present; empty allowed only for a single-node decision graph
    ParallelGroups []ParallelGroup `json:"parallelGroups,omitempty"`  // optional
    Bounds         Bounds          `json:"bounds"`                    // required
}

type Node struct { // GraphNode
    ID           string `json:"id"`                     // required non-empty bounded identifier
    Kind         string `json:"kind"`                   // required; one of input | evaluation | composition | decision
    PayloadBytes int64  `json:"payloadBytes,omitempty"` // bounded payload accounting; validated against Bounds.MaxPayloadBytes
}

type Edge struct { // GraphEdge
    ID   string `json:"id"`   // required non-empty bounded identifier; distinct from From/To
    From string `json:"from"` // required; references an existing node ID (dependency)
    To   string `json:"to"`   // required; references an existing node ID (dependent)
}

type ParallelGroup struct {
    ID      string   `json:"id"`      // required non-empty bounded identifier
    NodeIDs []string `json:"nodeIds"` // required non-empty; each references an existing node ID; members must be mutually independent
}

type Bounds struct { // GraphBounds
    MaxNodes        int   `json:"maxNodes"`
    MaxEdges        int   `json:"maxEdges"`
    MaxDepth        int   `json:"maxDepth"`
    MaxWidth        int   `json:"maxWidth"`
    MaxFanOut       int   `json:"maxFanOut"`
    MaxPayloadBytes int64 `json:"maxPayloadBytes"`
}
```

Structural rules (deterministic):

- Node IDs are unique; edge IDs are unique; an edge ID is distinct from its own
  `from`/`to` values.
- Duplicate `(from,to)` pairs are invalid unless the governing profile explicitly
  permits labelled multi-edges distinguished by edge ID.
- Approved node roles are exactly `input`, `evaluation`, `composition`, `decision`;
  no other role is accepted (no new vocabulary introduced).
- A valid graph root/terminal shape is required; the empty-edge case is legal only
  when the sole node's `Kind` is `decision`.
- `Build` receives an already-resolved `Definition` plus the effective `Bounds`
  resolved from `Definition.Bounds` under the inherited ceiling hierarchy
  (section 1.2); the builder performs no registry/version lookup (closure 4).

**Validation ownership and `DECISION_*` mapping.** All structural failures map to
the existing composition codes at exact RFC 6901 pointers (no new code).

**Pointer-base rule (deterministic).** A graph definition may appear either as a
standalone artifact or embedded inside the profile bundle. Validators must use one
deterministic pointer convention:

- A standalone `decision-graph.json` uses root-relative graph-local pointers such
  as `/nodes/0/id` and `/edges/0/id`.
- A bundle-embedded graph definition uses the caller-supplied base path prefixed
  to the graph-local pointer, i.e. `/spec/graphs/{index}/nodes/0/id` and
  `/spec/graphs/{index}/edges/0/id`.
- Validators preserve the caller-supplied base path and append the graph-local
  pointer; they never rewrite or discard the base. The table below lists the
  **graph-local** pointer examples (the standalone/root-relative form); for a
  bundle-embedded graph the same graph-local suffix is appended after
  `/spec/graphs/{index}`.

This rule changes no graph semantics and introduces no new error code:

| Failure | `violations[].code` | Graph-local JSON Pointer example |
|---|---|---|
| Missing node ID | `DECISION_COMPOSITION_GRAPH_INVALID` | `/nodes/0/id` |
| Duplicate node ID | `DECISION_COMPOSITION_GRAPH_INVALID` | `/nodes/3/id` |
| Missing edge ID | `DECISION_COMPOSITION_GRAPH_INVALID` | `/edges/0/id` |
| Duplicate edge ID | `DECISION_COMPOSITION_GRAPH_INVALID` | `/edges/2/id` |
| Unknown endpoint (`from`/`to` not an existing node) | `DECISION_COMPOSITION_GRAPH_INVALID` | `/edges/1/to` |
| Duplicate disallowed `(from,to)` pair | `DECISION_COMPOSITION_GRAPH_INVALID` | `/edges/4` |
| Invalid node kind | `DECISION_COMPOSITION_GRAPH_INVALID` | `/nodes/2/kind` |
| Invalid graph root / terminal shape | `DECISION_COMPOSITION_GRAPH_INVALID` | `/nodes` |
| Invalid parallel group (unknown member / non-independent) | `DECISION_COMPOSITION_GRAPH_INVALID` | `/parallelGroups/0/nodeIds/1` |
| Malformed bounds | `DECISION_COMPOSITION_GRAPH_INVALID` | `/bounds/maxNodes` |
| Cycle | `DECISION_COMPOSITION_CYCLE` | `/edges` |
| Exceeded bounds (nodes/edges/depth/width/fan-out/payload) | `DECISION_COMPOSITION_LIMIT_EXCEEDED` | `/bounds/maxDepth` |
| Conflicting inputs | `DECISION_COMPOSITION_INPUT_CONFLICT` | `/nodes/1` |
| Graph version/reference mismatch | resolved in `bundle.ResolveGraph` (before `BuildGraph`) | — |

Graph version/reference compatibility is validated in `bundle.ResolveGraph`
before `BuildGraph`; `BuildGraph` only ever receives an already-resolved
`Definition` (closure 4).

### 7.9 Supporting data-model shapes (compact)

Concrete field-level shapes for every schema/common type the design creates.
`SecurityExceptionRef` (§7.7) and the graph types (§7.8) are defined above and are
not repeated here; the remaining common types are:

```go
// internal/decision
type ProfileRef struct {
    Name    string `json:"name"`    // profile id; required non-empty
    Version string `json:"version"` // semver; required; (name,version) composite key (closure 8)
}

type SemanticDecisionIdentity struct {
    Basis          string   `json:"basis"`            // profile-declared identity basis id; required non-empty (F13-SCALE-002)
    InputRefs      []string `json:"inputRefs"`        // bounded ordered input identity refs; required non-empty
    CanonicalScope string   `json:"canonicalScope"`   // canonical scope descriptor derived from metadata.scopeRef; required non-empty (F13-SCALE-002)
    Digest         string   `json:"digest,omitempty"` // opaque structural identity carrier; no algorithm selected/executed
}

type TrustCarrier struct {
    State                  string `json:"state"`                            // structural trust-state; required when carrier present
    CanonicalizationMethod string `json:"canonicalizationMethod,omitempty"` // structural id; no execution
    DigestAlgorithm        string `json:"digestAlgorithm,omitempty"`        // structural id; no computation
    CoveredFields          string `json:"coveredFields,omitempty"`          // covered-field descriptor
    SignatureAlgorithm     string `json:"signatureAlgorithm,omitempty"`     // structural id; no signing/verification
    // present-but-empty optional carrier is rejected (RID-08, closure 9)
}

type DecisionRelationship struct {
    Kind      RelationshipKind `json:"kind"`      // CORRECTS | SUPERSEDES | REVOKES
    TargetRef apimeta.TypedRef `json:"targetRef"` // predecessor; uid required; orientation new record -> predecessor
}

type AuditLinkage struct { // FEATURE-0013 AuditEvent extension value types (§7.4)
    DecisionRef        apimeta.TypedRef  `json:"decisionRef"`                  // uid required
    AuditObligationRef *apimeta.TypedRef `json:"auditObligationRef,omitempty"` // optional obligation identity
    Correlation        Correlation       `json:"correlation"`                  // reused FEATURE-0012 carrier; never redefined
}

type ProjectionRule struct {
    Audience string   `json:"audience"`           // profile-declared audience id; required
    Includes []string `json:"includes,omitempty"` // RFC 6901 pointers into the profile-declared canonical view
    Excludes []string `json:"excludes,omitempty"` // RFC 6901 pointers; overlap (equality/containment) with Includes is invalid
    // resolves against the canonical view, not one record; no redaction execution (closure 10)
}
```

`Correlation` (auditEventRefs, operationRef, traceRef), `CompositionRef`,
`SovereigntyCarrier`, `Finality`, and `Obligation` remain as declared in §7.1/§7.5
and are not restated here. The closed vocabulary types (`DecisionForm`,
`DecisionAuthority`, `AdjudicationOutcome`, `AuditOutcome`, `RelationshipKind`,
`Sensitivity`, `Finality`, `EffectiveState`) are string-backed constants per §7.5.

**`SemanticDecisionIdentity.CanonicalScope` (F13-SCALE-002).** Semantic identity
covers scope through the `CanonicalScope` descriptor string. It is:

- **Derived, not authoritative.** Computed deterministically from the normalized
  `metadata.scopeRef` and the profile scope rules. It is never an independent scope
  authority and never overrides, replaces, or competes with `metadata.scopeRef`,
  which remains the sole scope authority (§1.1, §9.1 Scope pass).
- **Platform form.** When `metadata.scopeRef` is absent/nil (Platform), the
  descriptor is the canonical Platform descriptor derived from that absence.
- **Non-Platform form.** When `metadata.scopeRef` is present, the descriptor
  includes the `ScopeKind` and `uid` from `metadata.scopeRef`.
- **Addressable.** It is a validation-addressable identity field at
  `/record/semanticIdentity/canonicalScope`.
- **Verified in the Scope pass.** The Scope validation pass (§9.1 step 2) verifies
  that the supplied `canonicalScope` descriptor matches the value computed from the
  normalized `metadata.scopeRef` and profile scope rules. A mismatch maps to the
  existing `DECISION_SCOPE_MISMATCH` code (no new `DECISION_*` code). This code is
  chosen because the failure is a derived-descriptor-vs-computed-value mismatch,
  which is exactly the "derived scope ≠ expected scope" semantics that
  `DECISION_SCOPE_MISMATCH` already owns (§9.2, F13-SCOPE-008 family);
  `DECISION_SCOPE_CONFLICT` is reserved for a genuine second/parallel/duplicate
  scope *source*, and using it here would wrongly imply `canonicalScope` is a
  competing scope authority, which it is not.

## 8. Interfaces / functions

All exported functions are pure or bounded in-memory; none performs I/O beyond
reading local bundle/fixture bytes through `apivalid.StrictDecode`.

```go
// internal/decision/validate
func ValidateDecisionRecord(rec decision.DecisionRecord, view bundle.BundleView) *apiproblem.Problem
func ValidateAuditEvent(ev apiconform.AuditEvent, ext decision.AuditLinkage, view bundle.BundleView) *apiproblem.Problem
func ValidateProfile(p decision.DecisionProfile) *apiproblem.Problem
func ValidateEvaluation(e decision.EvaluationResult, recScope apimeta.ScopeIdentity, p decision.DecisionProfile) *apiproblem.Problem
func ValidateScope(m apimeta.ObjectMeta, allowed []apimeta.ScopeKind, platformAllowed bool) *apiproblem.Problem
func ValidateRelationship(r decision.DecisionRelationship, chain []apimeta.TypedRef, maxChain int) *apiproblem.Problem
func ValidateTrust(c *decision.TrustCarrier, required bool) *apiproblem.Problem
func ValidateObligations(obs []decision.Obligation, supported map[string]bool) *apiproblem.Problem

// internal/decision/bundle
func Load(data []byte, mode apivalid.DecodeMode) (Bundle, *apiproblem.Problem) // reuses StrictDecode
func (v BundleView) Profile(id string, version string) (decision.DecisionProfile, bool) // (id,version) key
func (b Bundle) ResolveGraph(ref GraphRef) (graph.Definition, *apiproblem.Problem)      // version resolution here

// internal/decision/graph
func Build(def Definition, bounds Bounds) (Graph, *apiproblem.Problem) // no registry lookup
func (g Graph) Accounting() Metrics                                    // deterministic depth/width/fan-out/critical path

// internal/decision/compose
type Strategy func(in StrategyInput) (StrategyOutput, *apiproblem.Problem) // pure
func Replay(captured []decision.EvaluationResult, s Strategy, in StrategyInput) (StrategyOutput, *apiproblem.Problem)
```

## 9. Validation passes and exact error mapping

Validation is a fixed, deterministic, fail-closed pipeline. Passes short-circuit
on the first blocking failure within a pass but collect bounded `violations[]`
where architecture permits multiple. Every public failure returns the inherited
FEATURE-0012 Problem Details envelope: `type=urn:sovrunn:problem:validation-failed`,
`status=422`, `code=VALIDATION_FAILED`, RFC 6901 `violations[].field` pointer,
one closed `violations[].code`, and a redactable non-authoritative message
(F13-ERR-001/003). No new `Problem.type`, top-level code, HTTP meaning, or code
prefix is introduced (F13-ERR-004, AD-045). FEATURE-0012 top-level codes continue
to own malformed input, authN/authZ, absence, conflict, stale state, internal
failure, and dependency unavailability.

### 9.1 Pass order

1. **Decode / structural** — `apivalid.StrictDecode` (JSON-token duplicate +
   pre-normalization `yaml.Node` duplicate rejection, closure 13). Malformed input
   → FEATURE-0012 top-level codes (not a `DECISION_*` code).
2. **Scope** — sole-authority, canonical Platform absence, six-value set,
   ServiceInstance-as-subject, and `semanticIdentity.canonicalScope` descriptor
   matches the computed `metadata.scopeRef`-derived value (§7.9).
3. **Profile** — identity/lifecycle/schema/limits; limit resolution follows the
   inherited limit hierarchy (section 1.2): FEATURE-0012 platform/schema ceilings
   are absolute outer bounds; DecisionProfile ceilings must remain within those
   inherited ceilings; evaluator-specific ceilings may only narrow the
   DecisionProfile; the effective limit is the most restrictive applicable value.
4. **Evaluation** — evaluator type/result/limit/scope derivation.
5. **Composition** — strategy support; graph resolution (bundle) → graph build
   (in-memory) → bounds/cycle/limit/input-conflict.
6. **Trust** — carrier presence/syntax/state (RID-08 empty-state rule).
7. **Obligation** — vocabulary/payload/mandatory-support.
8. **Relationship** — deterministic conflict matrix (closure 5).
9. **Sensitivity/projection** — ceiling/floor/content-category/pointer overlap.

### 9.2 Exact code mapping

| Case | `violations[].code` | Owning requirement |
|---|---|---|
| Missing scope when Platform not permitted | `DECISION_SCOPE_REQUIRED` | F13-SCOPE-001; edge 1 |
| Out-of-vocabulary scope kind (incl. ServiceInstance) | `DECISION_SCOPE_INVALID` | F13-SCOPE-007; edge 6 |
| Second/parallel/duplicate scope source | `DECISION_SCOPE_CONFLICT` | F13-SCOPE-005; edge 2 |
| ServiceInstance subject scopeRef ≠ containing Project | `DECISION_SCOPE_MISMATCH` | F13-SCOPE-008 |
| `semanticIdentity.canonicalScope` descriptor ≠ computed `metadata.scopeRef`-derived value | `DECISION_SCOPE_MISMATCH` | F13-SCALE-002; §7.9 |
| Impermissible scope widening | `DECISION_SCOPE_WIDENING` | AC-13.2 |
| Unknown profile | `DECISION_PROFILE_UNKNOWN` | F13-PROF-007 |
| Unsupported profile version (no fallback) | `DECISION_PROFILE_VERSION_UNSUPPORTED` | F13-PROF-007; edge 16 |
| Inactive profile | `DECISION_PROFILE_INACTIVE` | F13-PROF-007 |
| Schema-invalid profile | `DECISION_PROFILE_SCHEMA_INVALID` | F13-PROF-007 |
| Missing/unknown/incomparable/above-ceiling limit | `DECISION_PROFILE_LIMIT_INVALID` | F13-PROF-006; edge 10 |
| Unsupported evaluator type | `DECISION_EVALUATION_TYPE_UNSUPPORTED` | F13-EVAL-005 |
| Invalid evaluation result | `DECISION_EVALUATION_RESULT_INVALID` | F13-EVAL-005 |
| Evaluation scope mismatch / unauthorized widening (no existence disclosure) | `DECISION_EVALUATION_SCOPE_MISMATCH` | F13-EVAL-004; edge 3 |
| Evaluator limit exceeds profile | `DECISION_EVALUATION_LIMIT_EXCEEDED` | F13-EVAL-005; edge 11 |
| Unsupported composition strategy | `DECISION_COMPOSITION_STRATEGY_UNSUPPORTED` | F13-COMP-002 |
| Invalid graph | `DECISION_COMPOSITION_GRAPH_INVALID` | F13-COMP-002 |
| Graph cycle | `DECISION_COMPOSITION_CYCLE` | F13-COMP-002; edge 4 |
| Graph bound exceeded | `DECISION_COMPOSITION_LIMIT_EXCEEDED` | F13-COMP-002; edge 5 |
| Composition input conflict | `DECISION_COMPOSITION_INPUT_CONFLICT` | F13-COMP-002 |
| Unknown obligation | `DECISION_OBLIGATION_UNKNOWN` | F13-OBLIG-003 |
| Invalid obligation | `DECISION_OBLIGATION_INVALID` | F13-OBLIG-003 |
| Unsupported mandatory obligation → allow unenforceable | `DECISION_OBLIGATION_UNSUPPORTED_MANDATORY` | F13-OBLIG-002/003; edge 14 |
| Required trust absent | `DECISION_TRUST_REQUIRED` | F13-TRUST-003 |
| Unknown trust state | `DECISION_TRUST_UNKNOWN` | F13-TRUST-003 |
| Expired trust | `DECISION_TRUST_EXPIRED` | F13-TRUST-003 |
| Revoked trust | `DECISION_TRUST_REVOKED` | F13-TRUST-003 |
| Trust expected≠observed | `DECISION_TRUST_MISMATCH` | F13-TRUST-003; edge 13 |
| Invalid relationship kind | `DECISION_RELATIONSHIP_KIND_INVALID` | F13-IMMUT-005 |
| Invalid relationship target | `DECISION_RELATIONSHIP_TARGET_INVALID` | F13-IMMUT-005 |
| Relationship cycle | `DECISION_RELATIONSHIP_CYCLE` | F13-IMMUT-004; edge 7 |
| Relationship chain limit exceeded | `DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED` | F13-IMMUT-004; edge 8 |
| Incompatible relationship combination | `DECISION_RELATIONSHIP_CONFLICT` | closure 5 |

### 9.3 Relationship conflict matrix (closure 5)

The relationship pass is fully deterministic. It validates one record's
`record.relationship` (`kind` + `targetRef`, orientation new record → predecessor)
against the already-established relationship chain/graph, using only the five
approved codes (`DECISION_RELATIONSHIP_KIND_INVALID`,
`DECISION_RELATIONSHIP_TARGET_INVALID`, `DECISION_RELATIONSHIP_CYCLE`,
`DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED`, `DECISION_RELATIONSHIP_CONFLICT`).
No new `RelationshipKind` value or `DECISION_*` code is introduced.

**Deterministic definitions.**

- *Current successor* of predecessor `T` via kind `K`: a non-superseded,
  non-revoked record whose `relationship` targets `T` with kind `K`.
- *Fork*: one predecessor has ≥2 distinct successor records. Structurally
  permitted.
- *Current-state conflict*: ≥2 successors are simultaneously effective/current for
  the same predecessor (ambiguous head).
- *Terminal revocation*: once `T` is `REVOKES`-linked, its chain is terminal; any
  later relationship targeting `T` or its revoked chain conflicts.

**Deterministic algorithm order** (short-circuits at the first failing step;
`kind` vocabulary validity is a structural precondition checked before this pass
and maps to `DECISION_RELATIONSHIP_KIND_INVALID` at `/record/relationship/kind`):

1. **Target validation** — `targetRef` present, resolvable, uid-pinned, not
   self-referential → else `DECISION_RELATIONSHIP_TARGET_INVALID`
   (`/record/relationship/targetRef`).
2. **Cycle check** → `DECISION_RELATIONSHIP_CYCLE`
   (`/record/relationship/targetRef`).
3. **Chain length** → `DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED`
   (`/record/relationship`).
4. **Same-target duplicate** → `DECISION_RELATIONSHIP_CONFLICT`
   (`/record/relationship`).
5. **Terminal revocation conflict** → `DECISION_RELATIONSHIP_CONFLICT`
   (`/record/relationship/targetRef`).
6. **Fork / current-state conflict** → `DECISION_RELATIONSHIP_CONFLICT`
   (`/record/relationship`).
7. **Projection** — no error; computes `EFFECTIVE`/`SUPERSEDED`/`REVOKED`
   (never stored).

**Governed cases.**

| Case | Valid? | `violations[].code` (if invalid) | JSON Pointer owner | Deciding step |
|---|---|---|---|---|
| Kind not in {CORRECTS, SUPERSEDES, REVOKES} | invalid | `DECISION_RELATIONSHIP_KIND_INVALID` | `/record/relationship/kind` | precondition |
| Missing / unresolvable / self-referential target | invalid | `DECISION_RELATIONSHIP_TARGET_INVALID` | `/record/relationship/targetRef` | 1 |
| Same target, same kind — adding the relationship does not create more than one effective/current successor after deterministic projection | valid — recorded as a fork successor; projection selects the single current head | — | — | 6→7 |
| Same target, same kind — adding the relationship would create more than one effective/current successor after deterministic projection | invalid | `DECISION_RELATIONSHIP_CONFLICT` | `/record/relationship` | 6 |
| Same target, different authoritative kinds (e.g. `CORRECTS` + `SUPERSEDES`) | invalid — incompatible kinds on one target | `DECISION_RELATIONSHIP_CONFLICT` | `/record/relationship` | 4 |
| Chain convergence (shared ancestor, acyclic) | valid | — | — | 2→7 |
| Fork: one predecessor → multiple distinct successors | valid (structurally permitted) | — | — | 6 |
| Multiple current successors (fork with >1 effective head) | invalid | `DECISION_RELATIONSHIP_CONFLICT` | `/record/relationship` | 6 |
| Current vs historical projection | valid — projection computes `EFFECTIVE`/`SUPERSEDED`/`REVOKED`; not stored | — | — | 7 |
| Revocation of a live target (terminal behavior) | valid `REVOKES` link; target chain becomes terminal | — | — | 6→7 |
| Correction after revocation (`CORRECTS` a revoked target/chain) | invalid — terminal | `DECISION_RELATIONSHIP_CONFLICT` | `/record/relationship/targetRef` | 5 |
| Supersession after revocation (`SUPERSEDES` a revoked target/chain) | invalid — terminal | `DECISION_RELATIONSHIP_CONFLICT` | `/record/relationship/targetRef` | 5 |
| Re-revocation of an already-revoked target | invalid — terminal | `DECISION_RELATIONSHIP_CONFLICT` | `/record/relationship/targetRef` | 5 |
| Cycle via correction / supersession / revocation | invalid | `DECISION_RELATIONSHIP_CYCLE` | `/record/relationship/targetRef` | 2 |
| Chain length exceeded | invalid | `DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED` | `/record/relationship` | 3 |

In all valid cases the predecessor record is never mutated (append-only,
section 6.7); effective state is a projection only and is never written back to a
canonical record.

## 10. Registry / bundle design

- The canonical registry is the declarative JSON bundle
  (`decision-profile-bundle.json`), governed by JSON Schema. Go maps are
  derivative views (F13-PROF-004); the bundle is the semantic source of truth.
- Versioned registries (profile, form, authority, evaluator type, strategy,
  reason code, obligation, event type, evidence type, projection) are keyed by
  the composite `(id, version)` tuple; any explicitly unversioned registry
  states so in the profile schema (closure 8). Duplicate detection, bundle
  loading, `BundleView` lookup signatures, fixture identities, and traceability
  all use the same key rule.
- Graph reference/version compatibility is resolved in `bundle.ResolveGraph`
  before an in-memory graph is built; the graph builder receives an
  already-resolved definition and performs no registry/version lookup
  (closure 4).
- Bundles are offline-verifiable: structure, identifier syntax, carrier presence,
  opaque expected-vs-observed states, and fail-closed trust transitions only;
  no network fetch, no digest computation, no signing/verification (F13-PROF-008,
  AD-025, AD-038). Unknown/expired/revoked/mismatched bundles fail closed
  (F13-CF-23; edge 16 — no fallback to older version).
- Trust/integrity carriers are algorithm-agile structural fields
  (canonicalization-method, digest-algorithm, covered-field descriptor,
  signature-algorithm, structural trust-state); no value is selected
  (AD-022, ADR-F13-002). Empty-state handling per RID-08.

## 11. Audit / Operation behavior

- **Audit linkage (`CONTRACT_NOW`).** Every meaningful decision carries ≥1
  `Correlation.AuditEventRefs` entry (F13-AUDIT-001). Audit outcomes are the
  closed `Succeeded`/`Denied`/`Failed` set, separate from decision outcomes
  (F13-AUDIT-002). Telemetry is never audit authority (F13-AUDIT-004, AD-026).
  `AuditEvent` accepts all six governance scopes; the Organization-scoped alpha
  fixture is retained as regression; a ServiceInstance-subject-under-Project
  fixture is added (F13-AUDIT-003).
- **Atomic durable-acceptance (`INVARIANT_FOR_LATER`).** The design carries and
  documents (does not implement) the invariant that a decision is not reported
  durable until its `DecisionRecord` and audit obligation are atomically accepted
  by the same local boundary, with the `AuditEvent` identity preallocated and
  stored in both, and asynchronous materialization reusing that identity without
  mutating the record (F13-AUDIT-005/006/007, AD-014, AD-035). Only the carrier
  fields and partial-failure/reconciliation conformance are `CONTRACT_NOW`.
- **Operation handoff (`CONTRACT_NOW` for reference use).** `DecisionRequest` is
  calling-domain-owned; FEATURE-0013 defines no shared, canonical, or common
  `DecisionRequest` schema, Go type, registry entry, resource, persistence model,
  or public API resource. Completed synchronous
  handling returns a `DecisionRecord` value/reference; accepted asynchronous
  handling returns the existing FEATURE-0012 `Operation` (`LongRunningOperation`)
  value/reference; rejection returns inherited Problem Details (F13-OBJ-007,
  AD-009). No FEATURE-0013 pending-decision response envelope or non-final
  response contract is defined; no `PendingDecision`/`DeferredDecision`/
  `DecisionPending` or other non-final envelope is defined (F13-OBJ-008). The
  generic `Operation` contract is referenced by typed reference and never
  redefined (Non-goal 15).

## 12. Security / privacy

- **Structural boundary only.** FEATURE-0013 validates schema, typed
  classification, declared bounds, projection restrictions, and deterministic
  conformance. It implements no secret/credential/PII/malware/DLP/content-scanning
  engine; producers must redact before submission (F13-SEC-006; F13-EVAL-08
  proves arbitrary text is not lexically scanned).
- **Sensitivity + classification coexistence.** The four-value `Sensitivity`
  vocabulary and the seven-value FEATURE-0012 `DataClassification` coexist; the
  floor map (Public→PUBLIC, Customer-visible/Tenant-confidential/
  Operator-confidential→CONFIDENTIAL, Internal→INTERNAL, Sensitive/
  Secret-reference-only→RESTRICTED) may be raised, never lowered; unmapped/
  unknown/contradictory/below-floor fails closed (F13-SEC-002/003; F13-SEC-07/08).
- **Projection (closure 10).** Projection pointers resolve against the canonical
  schema/view declared by the governing profile, not one example record. Overlap
  includes exact equality and ancestor/descendant containment; unknown pointers,
  overlapping includes/excludes, audience mismatch, and prohibited-category
  exposure fail deterministically with existing projection/profile/evaluation
  codes. No redaction execution or projected-output generation is implemented
  (F13-SEC-007). AI receives only an authorized projection and is never
  authoritative; recommendations are labelled (F13-SEC-008); core works AI-off
  (F13-SEC-009).
- **Fail-open (closure 1).** `SecurityExceptionRef` is a bounded structural
  approval-evidence reference (apiVersion/kind/name/uid + scope, owner, approving
  authority, compensating controls, effective interval, audit treatment,
  reassessment trigger, covered failure modes). FEATURE-0013 validates shape,
  requiredness, scope compatibility, expiry against captured/effective decision
  time, uid pinning, and fail-closed-when-invalid only. It does not approve,
  issue, revoke, execute, store, authorize, or workflow exceptions
  (F13-PROF-002, AD-043). Strategy-level fail-open is effective only through a
  valid profile `SecurityExceptionRef` (closure 2).
- **Runtime authN/authZ + break-glass follow-up (`INVARIANT_FOR_LATER`).** The
  every-hop authentication/authorization duty and break-glass follow-up-review
  duty are carried and documented for a later owning feature; the mandatory
  break-glass reason carrier itself is `CONTRACT_NOW` under F13-SOV-007
  (F13-AUTHZ-005/006, AD-039).
- **Erasure.** FEATURE-0013 defines only structural tombstone/custody/retention/
  legal-hold linkage carriers; it implements no erasure, key destruction, or
  cryptographic process (F13-IMMUT-007, AD-004).

## 13. Testing / conformance

- All fixtures are deterministic, pure or bounded in-memory contract tests; none
  authorizes persistence, networking, queues, workers, cryptographic execution,
  AI runtime, or external integration (F13-CONF-003).
- Coverage is counted by scenario ID; one physical fixture may cover several IDs,
  but every ID (parent, suffix, and additional case) has an explicit
  coverage-matrix row (F13-CONF-002/004). The matrix lives in
  `tests/conformance/decision_coverage_matrix.go`.
- The reference kernel (`compose.Replay` + pure strategies) proves contract
  semantics without a runtime engine (F13-CONF-005, AD-012, AD-037). Complex
  composite scenarios (F13-CF-06/07/16/26) are fixture/reference proofs, not
  executors.
- Positive and negative fixtures are registered through `apiconform`
  (`TypeBindings` + executable checks), reusing `VerifyGoTypeAgainstSchema` for
  Go↔schema agreement and `apivalid.StrictDecode` for decode equivalence.
- Compound families require every architecture-ordered suffix
  (F13-CF-04.a–.i, 08.a–.d, 11.a–.e, 12.a–.e, 16.a–.f, 20.a–.f, 26.a–.h,
  28.a–.f); additional cases F13-SCOPE-01…11, F13-EVAL-01…08, F13-SEC-01…08,
  F13-TRUST-01…04, F13-COMPAT-01…09 each get a row (section 18.4).

## 14. Verification

Design-stage verification obligations (executed at implementation, not now):

1. `go build ./...` and `go test ./...` pass; `internal/decision` tree and
   `apiconform` additions compile with stdlib + `yaml.v3` only.
2. `apischema.VerifyGoTypeAgainstSchema` passes for every FEATURE-0013 TypeBinding
   (Go↔schema requiredness/nullable/minItems agreement, closure 11).
3. Baseline digests recomputed and approved through the existing FEATURE-0012
   mechanism; `BASELINE_MANIFEST.json`/`BASELINE_APPROVALS.json` updated; no
   parallel hierarchy (closure 6).
4. Import-graph test confirms the acyclic package DAG and absence of prohibited
   runtime/provider/crypto imports (section 27.7).
5. Every conformance ID row in the coverage matrix resolves to ≥1 fixture; every
   negative fixture returns the exact mapped `violations[].code`.
6. The `feature-0013-architecture-boundary-check.py` gate passes (contract-only,
   one scope class per artifact, no invented codes/vocab/scope/workflow).

## 15. Non-goals

Reproduced from requirements (architecture sections 3, 27.3). FEATURE-0013
creates none of the following:

1. Production decision/authorization/policy/placement/audit/evidence/identity/
   profile-registry/synchronization/cryptographic/AI service.
2. Database schema, migration, repository, transaction/outbox, event store, WORM
   store, cache, distributed lock, or lease.
3. Queue, broker, workflow/durable-execution runtime, scheduler, controller,
   worker pool, reconciler, retry daemon, or dead-letter processor.
4. Network registry client, remote policy resolver, evidence fetcher, HSM/notary
   client, provider adapter, model client, vector store, or agent.
5. Rules engine, policy language, workflow DSL, DAG execution engine, expression
   interpreter, dynamic plugin system, or code generator beyond repo schema
   tooling.
6. Production authorization tokens, cache invalidation, revocation distribution,
   multi-site consensus, disconnected synchronization, or conflict-resolution
   service.
7. Provider/runtime/database/broker/identity/cryptographic/model-specific
   dependencies in core packages.
8. A shared `DecisionRequest` schema, Go type, registry entry, resource, or
   persistence model.
9. Any pending-decision response envelope.
10. Adapter interfaces (owned by FEATURE-0016).
11. Cryptographic algorithm selection, digest computation, signing, verification,
    timestamping, notarization, WORM enforcement, HSM integration, or
    canonicalization execution (ADR-F13-002).
12. Production numeric defaults, system-wide graph ceilings, retention periods,
    or SLO targets not approved by architecture.
13. India-specific classification vocabulary or retention schedules.
14. Concrete provider/runtime/database/broker/identity/AI adapters.
15. The `Operation` contract itself (FEATURE-0012/ADH-2026-013).

## 16. Resolved design questions

| DQ | Resolution | Reference |
|---|---|---|
| DQ-13-01 | RID-02: struct layout embedding `apimeta.TypeMeta`/`ObjectMeta`, `*TypedRef`+omitempty for optional refs, string-backed vocab types | Section 7 |
| DQ-13-02 | RID-03: per-object schemas + `_common` sub-schemas + single bundle schema; JSON Schema `$ref` composition | Sections 6.1, 10 |
| DQ-13-03 | RID-01: `internal/decision` tree owns domain; `apiconform` limited to binding/checks | Section 5 |
| DQ-13-04 | RID-04: adjacency-list DAG built from resolved definition; deterministic accounting; no lookup in builder | Sections 5, 10; closure 3/4 |
| DQ-13-05 | RID-05: fixed 9-pass order | Section 9.1 |
| DQ-13-06 | RID-06: pure strategy func + `Replay`, no evaluator rerun | Section 8 |
| DQ-13-07 | RID-07: fixtures under `tests/conformance/fixtures/decision/{positive,negative}` + matrix file | Section 6.4 |
| DQ-13-08 | Closed — dependency reconciliation is locked in requirements; design renders already-decided evidence only, no reopen/reformat | Requirements reconciliation section |
| DQ-13-09 | RID-07: coverage-matrix file enumerates all families, suffixes, and additional cases | Sections 13, 18.4 |
| DQ-13-10 | Feature-gate check verifies section existence, applicability, roles, named/versioned profiles, reused contracts, risks, exceptions, evidence, and non-redefinition — one lightweight gate, no registry/index/approval workflow | Section 11 of arch 28.7; F13-ADOPT-004/005 |

## 17. Architecture drift checks

| Check | Status |
|---|---|
| No provider-specific hardcoding in core | PASS — provider-native data only via adapters/typed refs (AD-023) |
| No Kubernetes-only assumptions | PASS — no CRD/operator/K8s-native type |
| No PostgreSQL/placement lifecycle logic | N/A — FEATURE-0013 defines no placement/provisioning |
| No custom policy engine in handlers | PASS — policy evaluation stays behind later adapters |
| No raw secret storage | PASS — references/digests preferred; no secret embedding |
| No customer-facing IaaS leakage | PASS — governance-only boundary; no provider identifiers |
| Explainable DecisionRecord | PASS — structured rationale/reasons/alternatives/obligations |
| Defined audit behavior | PASS — linkage, obligation carrier, six-scope expansion |
| Preserved adapter boundaries | PASS — no adapter interface created |
| No new public code/vocab/scope/error/approval workflow | PASS — all reproduced from architecture |
| No parallel FEATURE-0012 baseline workflow | PASS — reuses `BASELINE_MANIFEST.json`/`BASELINE_APPROVALS.json` (closure 6) |
| Superseded/patched draft not used as source | PASS — clean regeneration (closure 14) |

## Architecture traceability

Controlling source:
`docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`,
approved through ADH-2026-017. Every design component cites its architecture
source and reproduces the architecture-owned implementation class. Only
`CONTRACT_NOW` produces artifacts.

### 18.1 Architecture sections 1–28 and 27.9

| Section | Title | Disposition in design | Design component |
|---|---|---|---|
| 1 | Decision status and authority | Governance-only | Metadata; controlling inputs (§2) |
| 2 | Problem statement | Translated | Overview (§1) |
| 3 | Scope | Translated | Non-goals (§15); package boundaries (§5) |
| 4 | Architectural principles | Translated | Distributed across §5–§13 |
| 5 | Canonical object model | Translated | Data models (§7) |
| 5.4 | FEATURE-0012 profile/boundary inheritance | Translated | §7 profile bindings; §6.1 schemas |
| 6 | Extensible decision semantics | Translated | §7 vocabularies; §10 registry |
| 6.1 | Stable common envelope | Translated | §7.1 DecisionRecord |
| 6.2–6.4 | Forms, authority, adjudication | Translated | §7.5 vocabularies |
| 6.5 | Finality/effective state | Translated | §7.1; §9.3 |
| 6.6 | Audit outcomes | Translated | §11 |
| 6.7 | Immutability/correction | Translated | §7.1; §9.3; §12 erasure |
| 6.8 | DecisionProfile contract | Translated | §7.3; §12 fail-open |
| 6.9 | Anticipated families | Governance-only | Extensibility validated, no impl |
| 6.10 | Authorization specialization | Translated | §7.3 (profile), §11 |
| 7 | Evaluation and composition | Translated | §7.2; §8; §10 |
| 7.1 | Atomic evaluation | Translated | §7.2 |
| 7.2 | Composite decision | Translated | §8 compose; §10.RID-10 |
| 7.3 | Hierarchical evaluation | Invariant | AD-007 documentation (§12) |
| 7.4 | Decision dependency graph | Translated | §5/§10 graph (closure 3/4) |
| 8 | Execution model | Translated / Invariant | §11 Operation handoff |
| 8.5 | Data-plane locality | Translated | §10 offline bundle |
| 9 | Statelessness/scalability/cost | Translated / Invariant / Deferred | §7.1 identity carriers; §13 kernel |
| 10 | Audit consistency | Translated / Invariant | §11 |
| 11 | Sovereign deployment | Translated | §7.1 sovereignty carrier; §10 |
| 12 | Security/privacy/trust | Translated | §12 |
| 12.1 | Runtime authN/authZ + break-glass | Invariant | §12 (F13-AUTHZ-005/006) |
| 12.3 | Integrity/non-repudiation | Deferred (exec) / Translated (carriers) | §10 trust; ADR-F13-002 |
| 12.4 | Sensitivity classification | Translated | §7.5; §12 floor map |
| 13 | AI-first | Translated | §12 (AI projection, labelling, AI-off) |
| 14 | Reuse and standards | Translated | §3 reuse summary |
| 15 | Compatibility and evolution | Translated | §6.2 baseline; §9 error binding |
| 15.1/15.2 | Public Problem Details binding + closed registry | Translated | §9.2 code mapping |
| 15.3 | Dependency reconciliation | Closed | Locked in requirements; not reopened (DQ-13-08) |
| 16 | Observability/operability | Invariant | Documentation; no runtime |
| 17 | Conformance model | Translated | §13; §18.4 |
| 18 | Matrix E risk register | Governance-only | §18.3 controls; §6.5 worksheet (closure 12) |
| 19 | Decision scorecard | Governance-only | Classes reproduced §18.2 |
| 20 | Required revisions | Governance-only | Historical rationale |
| 21 | Deferred decisions | Deferred | §15 non-goals |
| 22 | Architecture acceptance criteria | Governance-only | Approval evidence |
| 23 | Reassessment triggers | Governance-only | Later gates |
| 24 | Normative local references | Governance-only | §2 controlling inputs |
| 25 | External review inputs | Governance-only | Context |
| 26 | Foundation-principle validation | Governance-only | §17 drift checks |
| 27 | Anti-overengineering guardrails | Translated | §5 necessity tests; §14 verification |
| 27.1 | Scope classification | Translated | Every component labelled |
| 27.2 | Allowed artifacts | Translated | §6 files |
| 27.3 | Prohibited implementation | Translated | §15 non-goals; §17 |
| 27.4 | Simplicity/necessity tests | Translated | §5 package justification |
| 27.5 | Complexity budget | Translated | §13 kernel; single core model |
| 27.6 | Stop/escalation | Translated | §19 (no ADR required) |
| 27.7 | Anti-overengineering verification | Translated | §14 |
| 27.8 | Closed architecture boundary | Translated | Honored throughout; only mechanics resolved (§4) |
| 27.9 | Downstream design closure controls | Translated | §18.5 (each control mapped) |
| 28 | Lightweight downstream adoption | Translated | §16 DQ-13-10 |

### 18.2 AD-001 through AD-045 (architecture-owned class reproduced exactly)

| AD | Class | Design disposition | Component |
|---|---|---|---|
| AD-001 | CONTRACT_NOW | Translated | §7.1 DecisionRecord kind/naming |
| AD-002 | CONTRACT_NOW | Translated | §15 no DecisionRequest type |
| AD-003 | CONTRACT_NOW | Translated | §7.5 adjudication vocab |
| AD-004 | CONTRACT_NOW | Translated | §7.1; §9.3; §12 erasure carriers |
| AD-005 | CONTRACT_NOW | Translated | §8 strategies |
| AD-006 | CONTRACT_NOW | Translated | §5/§10 graph |
| AD-007 | INVARIANT_FOR_LATER | Invariant | §12 doc (EffectivePolicyContext ref only) |
| AD-008 | INVARIANT_FOR_LATER | Invariant | §11 sync path (kernel only) |
| AD-009 | CONTRACT_NOW | Translated | §11 Operation handoff |
| AD-010 | INVARIANT_FOR_LATER | Invariant | §6.5 documentation |
| AD-011 | DEFERRED | Deferred | §15 non-goals |
| AD-012 | CONTRACT_NOW | Translated | §8/§13 pure kernel |
| AD-013 | INVARIANT_FOR_LATER | Invariant | §7.1 retry key carrier |
| AD-014 | INVARIANT_FOR_LATER | Invariant | §11 atomic acceptance (carrier only) |
| AD-015 | CONTRACT_NOW | Translated | §7 scope + AuditEvent expansion |
| AD-016 | CONTRACT_NOW | Translated | §7.1 sovereignty carrier |
| AD-017 | CONTRACT_NOW | Translated | §10 offline bundle profiles |
| AD-018 | CONTRACT_NOW | Translated | §10 sync/import carriers |
| AD-019 | CONTRACT_NOW | Translated | §12 projection/minimization |
| AD-020 | CONTRACT_NOW | Translated | §12 AI projection |
| AD-021 | CONTRACT_NOW | Translated | §12 AI-off |
| AD-022 | CONTRACT_NOW | Translated | §10 trust carriers (no selection) |
| AD-023 | CONTRACT_NOW | Translated | §17 provider neutrality |
| AD-024 | CONTRACT_NOW | Translated | §12 references/digests carriers |
| AD-025 | CONTRACT_NOW | Translated | §10 static registries |
| AD-026 | CONTRACT_NOW | Translated | §11 telemetry-not-audit |
| AD-027 | CONTRACT_NOW | Translated | §7.1 jurisdiction as profile data |
| AD-028 | CONTRACT_NOW | Translated | §13/§14 contract-only + gates |
| AD-029 | CONTRACT_NOW | Translated | §7.3 profile extension |
| AD-030 | CONTRACT_NOW | Translated | §7.5 closed forms |
| AD-031 | CONTRACT_NOW | Translated | §7.5 authority |
| AD-032 | CONTRACT_NOW | Translated | §7.3/§11 authorization profile |
| AD-033 | CONTRACT_NOW | Translated | §11 risk-based recording |
| AD-034 | CONTRACT_NOW | Translated | §9.3 relationship/effective-state |
| AD-035 | INVARIANT_FOR_LATER | Invariant | §11 async materialization |
| AD-036 | CONTRACT_NOW | Translated | §7.1 semantic identity vs retry key |
| AD-037 | CONTRACT_NOW | Translated | §8 replay |
| AD-038 | CONTRACT_NOW | Translated | §10 version-pinned bundle/trust state |
| AD-039 | INVARIANT_FOR_LATER | Invariant | §12 local authz cache invariant |
| AD-040 | CONTRACT_NOW | Translated | §10 bounds; §9 limit codes |
| AD-041 | INVARIANT_FOR_LATER | Invariant | §12 tiered crypto documentation |
| AD-042 | CONTRACT_NOW | Translated | §10 disconnected sync contract |
| AD-043 | CONTRACT_NOW | Translated | §12 fail-open/SecurityExceptionRef |
| AD-044 | CONTRACT_NOW | Translated | §16 adoption gate |
| AD-045 | CONTRACT_NOW | Translated | §9.2 public code binding |

### 18.3 F13-R01 through F13-R31 (controls mapped to design components)

Controls only; final residual acceptance remains a human gate (arch 18.4–18.6).
F13-R04, F13-R08, F13-R13 remain explicitly open at High target residual with
later production-entry gates.

| Risk | Design component implementing/carrying controls |
|---|---|
| F13-R01 | §10 governed registration/versioning; §13 duplicate-profile & new-family fixtures |
| F13-R02 | §7.3 versioned result schemas; §9 negative conformance |
| F13-R03 | §7.5 authority; §12 AI non-authoritative; authority-confusion fixtures |
| F13-R04 (High) | §12 EffectivePolicyContext ref/provenance carrier; §7.1 semantic identity; later gate |
| F13-R05 | §7.1 retry key vs semantic identity; F13-CF-20 fixtures |
| F13-R06 | §11 atomic-acceptance carrier; partial-failure/reconciliation fixtures |
| F13-R07 | §11 local obligation acceptance; async downstream |
| F13-R08 (High) | §10 bounds; §9 limit codes; §13 limit/cancellation fixtures; later gate |
| F13-R09 | §8 captured-output replay; F13-CF-21 |
| F13-R10 | §12 projection/minimization; redaction/prompt-leak fixtures |
| F13-R11 | §9 obligation codes; F13-CF-24 |
| F13-R12 | §12 local-artifact invariant (documented) |
| F13-R13 (High) | §12 authz freshness/revocation invariant; F13-CF-25; later gate |
| F13-R14 | §12 fail-closed + SecurityExceptionRef checks |
| F13-R15 | §9.3 relationship matrix; chain/cycle/fork fixtures |
| F13-R16 | §12 minimization/redaction; secret/PII/cross-scope negatives |
| F13-R17 | §12 tombstone/retention/legal-hold carriers |
| F13-R18 | §15 no crypto service; §12 no-network/no-crypto gate |
| F13-R19 | §10 origin-aware manifests/replay/duplicate/ordering; F13-CF-28 |
| F13-R20 | §10 no last-writer-wins; conflict evidence; F13-CF-28.f |
| F13-R21 | §17 provider-neutrality; import/portability fixtures |
| F13-R22 | §5/§15 contract-only; boundary-check gate |
| F13-R23 | §10 bounded payloads; §7.1 hot-summary/cold-reference carriers |
| F13-R24 | §10 version-pinned bundles + structural trust/revocation |
| F13-R25 | §7.1 trusted-time provenance carrier; skew/expiry fixtures |
| F13-R26 | §9 scope mismatch/no-existence-disclosure; cross-scope fixtures |
| F13-R27 | §9.2 closed public registry; §6.2 baseline diff gate |
| F13-R28 | §11 six-scope AuditEvent; Organization regression; ServiceInstance-subject fixtures |
| F13-R29 | §7.1 DecisionRecord naming; repo identity/baseline checks |
| F13-R30 | §10 tiered profiles; §5 minimal-path budget |
| F13-R31 | §16 one gate, no registry/index/approval workflow |

### 18.4 Stable conformance IDs (enumerated individually)

Families F13-CF-01 … F13-CF-28, every mandatory suffix, and every additional
case each get a coverage-matrix row. Design component: `tests/conformance`
fixtures + `decision_coverage_matrix.go`, exercised through `apiconform`. Class:
`CONTRACT_NOW` (AD-028).

Base families: F13-CF-01, F13-CF-02, F13-CF-03, F13-CF-04, F13-CF-05, F13-CF-06,
F13-CF-07, F13-CF-08, F13-CF-09, F13-CF-10, F13-CF-11, F13-CF-12, F13-CF-13,
F13-CF-14, F13-CF-15, F13-CF-16, F13-CF-17, F13-CF-18, F13-CF-19, F13-CF-20,
F13-CF-21, F13-CF-22, F13-CF-23, F13-CF-24, F13-CF-25, F13-CF-26, F13-CF-27,
F13-CF-28.

Compound suffixes (each its own row):
- F13-CF-04.a, .b, .c, .d, .e, .f, .g, .h, .i (ranking, classification,
  resolution, allocation, plan, assessment, recommendation, advisory, simulation)
- F13-CF-08.a, .b, .c, .d (timeout, missing evidence, evaluator failure, policy
  conflict)
- F13-CF-11.a, .b, .c, .d, .e (supersession, correction, revocation, cycle,
  chain-limit)
- F13-CF-12.a, .b, .c, .d, .e (unauthorized profile, authority, projection,
  obligation bypass, cross-scope)
- F13-CF-16.a, .b, .c, .d, .e, .f (graph size, depth, width, fan-out, timeout,
  payload budget)
- F13-CF-20.a, .b, .c, .d, .e, .f (policy, profile, strategy, authority,
  evaluator-version, validity)
- F13-CF-26.a, .b, .c, .d, .e, .f, .g, .h (cancellation, remote-call, byte,
  concurrency, retry, jurisdiction, time, cost budgets)
- F13-CF-28.a, .b, .c, .d, .e, .f (replay, duplicate, ordering, unknown trust
  root, revoked origin, conflict reconciliation)

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

### 18.5 Section 27.9 closure control resolution (explicit)

| Control | Design resolution | Component |
|---|---|---|
| 1. Fail-open exception representation | `SecurityExceptionRef` bounded structural type; validate shape/requiredness/scope-compat/expiry/uid-pinning; no approval/revoke/workflow | §12; RID-08 not applicable; §6.1 schema |
| 2. Strategy failure posture boundary | Strategy fail-open effective only via valid profile `SecurityExceptionRef`; else fail closed | RID-10; §8 |
| 3. Graph identity/accounting | Stable node/edge IDs; edge id distinct from (from,to); orientation dependency→dependent; parallel-group independence; deterministic once-counted metrics | §5/§10 graph; RID-04 |
| 4. Graph version ownership | Version resolved in `bundle.ResolveGraph`; builder no lookup | RID-04; §8/§10 |
| 5. Relationship conflict matrix | Deterministic matrix using five approved codes | §9.3 |
| 6. FEATURE-0012 baseline reuse | Reuse `BASELINE_MANIFEST.json`/`BASELINE_APPROVALS.json`; no parallel hierarchy | §6.2; RID-03 |
| 7. AuditEvent package ownership | `internal/decision` owns extension value types; `apiconform` limited to binding/checks | RID-01; §5; §7.4 |
| 8. Versioned registry keys | `(id, version)` keying across duplicate detection/load/lookup/fixtures/traceability | §10; RID-03 |
| 9. TrustCarrier empty-state | Design chooses reject-as-invalid for present-but-empty optional carriers | RID-08; §7.6 |
| 10. Projection pointer target/overlap | Resolve against profile-declared canonical view; overlap = equality + containment; no redaction execution | §12 |
| 11. Requiredness for public contracts | Full per-field-kind enumeration; Go tags/JSON Schema/validators/fixtures agree | §7.6 |
| 12. Matrix E automation boundary | Worksheet stages IDs/controls/evidence + blank human fields; no calculation | §6.5 |
| 13. Scope pre-scan/decode reuse | Reuse `apivalid.StrictDecode` (JSON token dup + pre-normalization YAML AST dup); no zero-alloc claim | RID-09; §9.1 |
| 14. Clean regeneration | Architecture as sole source; superseded/patched draft unused; no ADR required | §1; §2; §19 |

## 19. ARCHITECTURE_DECISION_REQUIRED items

None. Every section 27.9 closure control and every `CONTRACT_NOW` requirement
maps to concrete design mechanics without changing a closed architecture
semantic. The single design-delegated choice (closure 9 trust empty-state) was
resolved within the options architecture explicitly permits (RID-08). No new
public vocabulary, `DECISION_*` code, Problem Details type, HTTP meaning,
`ScopeKind` value, approval workflow, or parallel baseline hierarchy is
introduced. Design does not proceed to tasks.
