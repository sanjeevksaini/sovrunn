# FEATURE-0013 — Decision Record and AuditEvent Standard: Design

> **Status: SUPERSEDED DRAFT — DO NOT USE FOR TASKS OR IMPLEMENTATION.**
> This design predates ADH-2026-015 and contains the superseded six-scope and parallel AuditScope approach. It is preserved as review evidence. It must be revised against freshly approved requirements before design approval or task generation.

Stage: Design
Feature: FEATURE-0013
Phase: Phase 2 — Reuse-First PaaS Fabric Foundation
Scope: CONTRACT_NOW — schemas, vocabularies, validation, conformance fixtures, compatibility, documentation only
Controlling architecture: docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md
Controlling handoff: ADH-2026-014

## 1. Overview

This design implements the FEATURE-0013 contract-only deliverables:

- Go struct definitions for `DecisionRecord`, `DecisionProfile`, `EvaluationResult`, and six-scope `AuditEvent`
- Closed vocabulary types for decision forms, authority, adjudication outcomes, audit outcomes
- Composition strategy and bounded dependency-graph schemas
- Immutability, correction, supersession, and revocation relationship contracts
- Atomic decision/audit-obligation acceptance contract (pure-function model)
- Conformance fixtures (49+ scenarios from architecture section 17, plus AC-13.7 evaluator-output scenarios)
- Downstream adoption contract (section 28 of controlling architecture)
- JSON Schema definitions for canonical schemas
- Validation rules and error codes aligned with FEATURE-0012 patterns

No production runtime, persistence, orchestration, provider adapter, AI service, or network service is implemented. All deliverables are deterministic, pure-function, in-memory, or documentary.

## 2. Resolved Decisions

| ID | Decision | Rationale |
|---|---|---|
| DQ-01 | `typedResult` uses `json.RawMessage` | Profile-declared typed results vary per family; `json.RawMessage` defers decoding to profile-specific validators while keeping the common envelope stable. Structural validation enforces size bounds. |
| DQ-02 | `DecisionProfile` is a Go struct with embedded validation rules | Profiles are validated by Go conformance code; no external schema-validation engine dependency. Profile schemas are also published as JSON Schema for documentation and cross-language conformance. |
| DQ-03 | Bounded graph uses adjacency list with typed nodes and explicit max-bounds struct | Adjacency list is simple, deterministic, and validatable. Profile-declared finite limits are conditionally mandatory via `GraphLimits` (see section 4.4.2). |
| DQ-04 | Semantic digest: DEFERRED pending ADR-F13-003 resolution (see section 17) | RFC 8785 JCS is the architecture candidate. The exact canonicalization algorithm/version and complete covered-field set require ADR-F13-003 approval before implementation. This design defines the `SemanticDigest` struct but does NOT implement digest computation or publish a partial covered-field list. Fixtures use `PLACEHOLDER_PENDING_ADR_F13_003` values. |
| DQ-05 | `AuditEvent` uses `Metadata.ScopeRef` as canonical scope source plus `Record.Scope` as the audit-specific scope discriminator with mandatory consistency validation | See section 3.5 for full scope semantics. |
| DQ-06 | Composition strategy registration uses a static Go registry map | No external file format or runtime discovery. Strategies are registered at init-time in conformance packages. Future phases may add dynamic registry via approved features. |
| DQ-07 | Atomic acceptance is modeled as a pure function returning `AcceptanceResult` struct | Fixtures prove invariants without persistence. The contract requires both artifacts or neither; error means neither was accepted. |
| DQ-08 | Conformance fixtures use Go table-driven tests with JSON test-case data files | Go tests provide deterministic execution. JSON data files allow cross-language fixture sharing using Go standard library `encoding/json` without introducing any new dependency. |
| DQ-09 | FEATURE-0013 schemas compose FEATURE-0012 types via Go embedding and JSON Schema `$ref` | `apimeta.TypeMeta`, `apimeta.ObjectMeta`, `apimeta.TypedRef`, `apimeta.ScopeRef` are embedded. No parallel grammar. |
| DQ-10 | Profile bundle uses a signed JSON manifest with trust-root identity, profile entries, and SHA-256 content digests | JSON for standard-library parsing; SHA-256 for integrity; signature format is algorithm-agile (field declares algorithm OID). |
| DQ-11 | Error codes extend FEATURE-0012 stable error-code pattern with decision-specific codes | New codes: `PROFILE_REJECTED`, `GRAPH_LIMIT_EXCEEDED`, `OBLIGATION_UNENFORCEABLE`, `SCOPE_MISMATCH`, `CHAIN_DEPTH_EXCEEDED`, `AUTHORITY_MISMATCH`, `PROJECTION_EXCEEDED`, `DIGEST_MISMATCH`. |
| DQ-12 | Four-package layout: `internal/decision/core/`, `internal/decision/audit/`, `internal/decision/bundle/`, `internal/decision/fixtures/` | See section 3.1 for rationale. |
| DQ-13 | `authorizedOutputProjection` is a bounded `json.RawMessage` with per-profile max bytes, max fields, allowed types declared in `ProjectionBounds` struct | Bounds are profile-declared. Validation is a pure function comparing the raw projection against the profile's `ProjectionBounds`. |
| DQ-14 | Cross-scope safety validation is a separate validation pass invoked after projection extraction and before acceptance | Scope identity is propagated via the `GovernedProjection.SourceScope` field. The validation pass compares source scope against target record scope and rejects cross-tenant content disclosure. |

## 3. Architecture

### 3.1 Package Structure

The package structure resolves the parent/child import cycle identified in review. The shared decision/audit types that cross the decision ↔ audit boundary live in a `core` leaf package that both `internal/decision/` (root acceptance/validation orchestration) and `internal/decision/audit/` can import without circular dependency.

```text
internal/decision/
├── core/
│   ├── types.go              # DecisionRecord, DecisionForm, DecisionAuthority, Adjudication, Finality, Correlation, Relationship
│   ├── profile.go            # DecisionProfile, ProfileRef, ProfileLifecycle, GraphLimits, ProjectionBounds, all profile sub-types
│   ├── evaluation.go         # EvaluationResult, EvaluatorType, EvaluationOutcome, GovernedProjection, TimingEnvelope
│   ├── composition.go        # CompositionStrategy, StrategyRef, GraphNode, GraphEdge, DependencyGraph, FailPosture
│   ├── identity.go           # IdempotencyKey, SemanticDigest (declaration only, no computation)
│   ├── scope.go              # AuditScope, AuditScopeKind constants, scope validation helpers
│   ├── acceptance.go         # AuditObligation, AcceptanceResult, Acceptor interface
│   └── errors.go             # Decision-specific error codes extending FEATURE-0012 patterns
├── validation.go             # Envelope validation, vocabulary checks, size/depth/bounds, cross-scope safety, projection bounds
├── registry.go               # In-memory profile registry, composition strategy registry, ProfileRegistry/StrategyRegistry interfaces
├── audit/
│   ├── types.go              # Six-scope AuditEvent, AuditOutcome, AuditEventRecord
│   ├── acceptance.go         # Audit-obligation acceptance contract
│   ├── validation.go         # AuditEvent validation, scope consistency
│   └── boundary.go           # RequiresDurableRecord aggregation boundary contract
├── bundle/
│   ├── types.go              # ProfileBundle, BundleManifest, BundleTrustRoot, BundleEntry, BundleSignature
│   ├── verifier.go           # BundleVerifier interface and inMemoryBundleVerifier, TrustRootStore interface
│   └── digest.go             # Bundle-entry content digest computation (SHA-256) — NOT semantic decision digest
└── fixtures/
    ├── testdata/             # JSON test case files for 49+ scenarios
    └── conformance_test.go   # Table-driven conformance fixture tests
```

**Rationale for four-package layout (vs. DQ-12 from prior revision):**

The prior design placed `AuditObligation` and `AuditScope` in `internal/decision/acceptance.go` (root package) while declaring `AuditScope` only in `internal/decision/audit/types.go`. This creates an unresolvable circular import: parent cannot import child. The solution extracts all shared types to `internal/decision/core/` which is a leaf package with no internal siblings as dependencies. Both root (`internal/decision/`) and `audit/` import `core/` without cycles.

### 3.2 Package Ownership and Import Rules

| Package | Owns | May Import |
|---|---|---|
| `internal/decision/core/` | All shared decision types: DecisionRecord, profile, evaluation, composition, identity, AuditScope, AuditObligation, AcceptanceResult, Acceptor interface, error codes | `internal/apimeta`, `internal/apiproblem`, standard library |
| `internal/decision/` (root) | Validation orchestration, profile/strategy registries, interfaces | `internal/decision/core/`, `internal/apimeta`, `internal/apiproblem`, standard library |
| `internal/decision/audit/` | AuditEvent, AuditOutcome, AuditEventRecord, audit acceptance, audit validation, aggregation boundary | `internal/decision/core/`, `internal/apimeta`, `internal/apiproblem`, standard library |
| `internal/decision/bundle/` | ProfileBundle, BundleManifest, trust-root, bundle-entry content digest computation, bundle verifier | `internal/decision/core/`, `internal/apimeta`, `crypto/sha256`, `crypto/ed25519`, `encoding/hex`, `encoding/base64`, standard library |
| `internal/decision/fixtures/` | Test data files, conformance test runner | `internal/decision/core/`, `internal/decision/`, `internal/decision/audit/`, `internal/decision/bundle/`, `internal/apimeta`, `testing`, standard library |

Import direction rules:
- `internal/decision/core/` MUST NOT import any sibling or parent package in the `internal/decision/` tree
- `internal/decision/` (root) MAY import `internal/decision/core/` only (NOT `audit/`, `bundle/`, or `fixtures/`)
- `internal/decision/audit/` MAY import `internal/decision/core/` only (NOT root, `bundle/`, or `fixtures/`)
- `internal/decision/bundle/` MAY import `internal/decision/core/` only (NOT root, `audit/`, or `fixtures/`)
- `internal/decision/fixtures/` MAY import all packages in the tree (it is test-only)
- NO package in `internal/decision/` tree imports `internal/server`, `internal/api`, `internal/registry`, or any external dependency

### 3.3 Integration with FEATURE-0012

FEATURE-0013 types compose FEATURE-0012 types:

- `apimeta.TypeMeta` embedded for `apiVersion` + `kind`
- `apimeta.ObjectMeta` for metadata (name, uid, scopeRef, labels, annotations, timestamps)
- `apimeta.TypedRef` for all typed references (actorRef, subjectRefs, operationRef, auditEventRefs)
- `apimeta.ScopeRef` / `apimeta.ScopeKind` for governance scope
- `apiproblem.Problem` / `apiproblem.Violation` for validation errors
- JSON Schema `$ref` to `_common/` sub-schemas

No parallel metadata, reference, validation, or error grammar is created.

### 3.4 Canonical Scope Source for DecisionRecord

`Metadata.ScopeRef` (i.e., `apimeta.ObjectMeta.ScopeRef`) is the single canonical source of governance scope for a `DecisionRecord`. There is no separate record-level `ScopeRef` field.

Rules:
- `DecisionRecord.Metadata.ScopeRef` is the ONE canonical scope source
- Validation reads scope from `Metadata.ScopeRef`
- Audit linkage reads scope from `Metadata.ScopeRef`
- Cross-scope safety reads target scope from `Metadata.ScopeRef`
- `Metadata.ScopeRef` is `*apimeta.ScopeRef` (nullable; nil means platform scope per FEATURE-0012 D-16)
- Profile `ScopeSemantics.RequireScope` determines whether non-nil scope is mandatory
- `apiVersion`/`kind` values: `apiVersion: "decision.sovrunn.io/v1alpha1"`, `kind: "DecisionRecord"`

### 3.5 Scope Semantics: AuditEvent and ServiceInstance Compatibility

#### 3.5.1 The problem

FEATURE-0012 defines six `ScopeKind` values: `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`. The FEATURE-0013 architecture mandates AuditEvent support for: Platform, Organization, OrganizationUnit, Tenant, Project, ServiceInstance. `ServiceInstance` is NOT in FEATURE-0012's `ScopeKind` enum and `Provider` is NOT in the FEATURE-0013 architecture's audit scope list.

#### 3.5.2 Resolution: ADR-F13-004 with compliant temporary contract

**ADR-F13-004 status:** ARCHITECTURE_DECISION_REQUIRED — whether ServiceInstance should be added to FEATURE-0012's `ScopeKind` or remain a FEATURE-0013 local extension.

**Temporary compliant contract (valid until ADR-F13-004 resolves):**

1. `AuditScope` is defined in `internal/decision/core/scope.go` as a struct with a `Kind` field of type `string` (not `apimeta.ScopeKind` directly).
2. `AuditScope` validation accepts a closed set of six string values: `"Platform"`, `"Organization"`, `"OrganizationUnit"`, `"Tenant"`, `"Project"`, `"ServiceInstance"`.
3. `AuditScope` does NOT claim to compose or extend `apimeta.ScopeKind`. It is a FEATURE-0013-local type.
4. When `AuditScope.Kind` matches one of the five values that overlap with `apimeta.ScopeKind` (Platform through Project), it is semantically equivalent to the corresponding `apimeta.ScopeKind` value.
5. `ServiceInstance` is a valid `AuditScope.Kind` for FEATURE-0013 audit events only. It has no representation in `apimeta.ScopeKind` until ADR-F13-004 resolves.

#### 3.5.3 AuditEvent scope representation: canonical source and consistency

`AuditEvent` has TWO scope representations that MUST be consistent:

1. **`Metadata.ScopeRef`** (`*apimeta.ScopeRef`): The FEATURE-0012 standard governance scope. For ServiceInstance-scoped events, this field is nil (platform scope) OR references the parent governance scope (e.g., the Project that owns the ServiceInstance), per ADR-F13-004 temporary treatment.
2. **`Record.Scope`** (`AuditScope`): The FEATURE-0013 audit-specific scope discriminator with ServiceInstance support.

**Canonical source rule:** `Record.Scope` is the authoritative audit scope for FEATURE-0013 consumers. `Metadata.ScopeRef` provides FEATURE-0012 compatibility.

**Consistency validation invariant:**
- When `Record.Scope.Kind` is one of {Platform, Organization, OrganizationUnit, Tenant, Project}: `Metadata.ScopeRef` MUST be consistent — the `ScopeRef.Kind` must equal `Record.Scope.Kind` and if `Record.Scope.UID` is non-empty then `ScopeRef.UID` must match.
- When `Record.Scope.Kind` is `ServiceInstance`: `Metadata.ScopeRef` MUST reference the owning governance scope (the parent Project or Tenant) per ADR-F13-004 temporary treatment. The consistency rule is: `Metadata.ScopeRef` MUST NOT be nil (ServiceInstance events always have a parent governance scope) and its `Kind` must be one of {Tenant, Project}.
- Platform-scope consistency: when `Record.Scope.Kind` == "Platform", `Metadata.ScopeRef` MUST be nil (FEATURE-0012 canonical platform representation).

**Validation errors:**
- `SCOPE_MISMATCH` when `Record.Scope.Kind` ∈ {Org, OU, Tenant, Project} but `Metadata.ScopeRef.Kind` does not match
- `SCOPE_MISMATCH` when `Record.Scope.Kind` == "ServiceInstance" but `Metadata.ScopeRef` is nil or `Metadata.ScopeRef.Kind` ∉ {Tenant, Project}
- `SCOPE_MISMATCH` when `Record.Scope.Kind` == "Platform" but `Metadata.ScopeRef` is non-nil
- `VALIDATION_FAILED` when `Record.Scope.Kind` is not in the closed six-value set

### 3.6 Digest Taxonomy: Three Distinct Digest Contracts

This design distinguishes three separate digest concepts that the prior revision conflated:

| Digest Type | Defined In | Algorithm | Input | Verification | ADR-F13-003 Blocked? |
|---|---|---|---|---|---|
| **Bundle-entry content digest** | `internal/decision/bundle/digest.go` | SHA-256 | Raw bytes of the profile JSON file | Computed at verification time; compared to manifest-declared value | NO — fully specified and executable |
| **Protected-content integrity digest** | `EvaluationResult.IntegrityDigest` field | SHA-256 | Raw bytes of the referenced protected content | Provided by evaluator adapter; verified as non-empty + format-valid + matched against re-computation of provided content bytes | NO — fully specified (see section 3.6.2) |
| **Semantic decision digest** | `SemanticDigest` struct | Candidate: SHA-256 | Canonicalized decision identity fields (exact set TBD) | NOT IMPLEMENTED — computation deferred to ADR-F13-003 | YES — completely deferred |

#### 3.6.1 Bundle-Entry Content Digest (Fully Specified)

- **Algorithm:** SHA-256
- **Input:** The exact byte content of the profile JSON file referenced by `BundleEntry.Path`
- **Encoding:** `"sha256:" + lowercase hex(sha256(content_bytes))`
- **Verification:** At bundle verification time, the verifier reads the entry content bytes, computes `sha256(bytes)`, encodes as above, and compares string-equality to `BundleEntry.ContentDigest`. Mismatch → `DIGEST_MISMATCH` error.
- **Signature coverage:** The bundle signature covers the entire manifest JSON (excluding the `signature.value` field itself). The signed input is: the manifest JSON with `signature.value` set to empty string `""`, serialized as compact JSON with keys in declaration order. The signature is computed over these exact bytes.
- **Supported signature algorithm:** Ed25519 (via `crypto/ed25519` standard library). Algorithm field declares `"Ed25519"`. Future algorithms may be added via ADR.
- **Expiry:** Bundle is expired when `time.Now()` > `parse(signature.signedAt) + bundle.validityMs`. Expired bundles are rejected with `PROFILE_REJECTED`.

#### 3.6.2 Protected-Content Integrity Digest (Fully Specified)

This digest proves that specific protected content existed and was considered by an evaluator, without copying that content into the canonical record (AC-13.3).

- **Algorithm:** SHA-256
- **Input:** The exact byte content of the protected resource referenced by the evaluation (e.g., an evidence document, model response, policy snapshot)
- **Encoding:** `"sha256:" + lowercase hex(sha256(content_bytes))`
- **Field:** `EvaluationResult.IntegrityDigest`
- **When required:** When `EvaluationResult.InputRef` is non-nil (a reference to external content is used instead of inline content), `IntegrityDigest` MUST be non-empty and MUST use the `"sha256:<hex>"` format.
- **Verification contract (AC-13.7f, AC-13.7g):**
  - **Format validation (always):** If `InputRef` is non-nil, validate that `IntegrityDigest` is non-empty and matches the pattern `^sha256:[0-9a-f]{64}$`. Missing or malformed → `DIGEST_MISMATCH`.
  - **Content verification (when content bytes are available):** When the referenced content bytes are available to the acceptance gate (e.g., in conformance fixtures where test content is provided), compute `sha256(content_bytes)` and compare to the declared digest. Mismatch → `DIGEST_MISMATCH`.
  - **Content verification (when content bytes are NOT available):** When the acceptance gate cannot access the referenced content (deferred resolution, air-gapped source), format validation alone is enforced. This is NOT a deferred decision — it is the specified behavior for the unavailable-content case.
- **Fixture implications:**
  - F13-FIXTURE-042 (AC-13.7f): provides matching content bytes + digest → accepted
  - F13-FIXTURE-043 (AC-13.7g): provides mismatched content bytes + digest → rejected with `DIGEST_MISMATCH`
  - Both fixtures provide `contentBytes` in the test data so the verifier CAN compute and compare

#### 3.6.3 Semantic Decision Digest (Deferred — ADR-F13-003)

This digest provides semantic decision identity for deduplication and caching. It is FULLY DEFERRED.

- **Algorithm:** Candidate SHA-256, pending ADR-F13-003
- **Canonicalization:** Candidate RFC 8785 JCS, pending ADR-F13-003
- **Covered fields:** The architecture (section 9) requires coverage of: scope, purpose, profile name+version, authority, canonical input snapshot, EffectivePolicyContext, composition strategy/version, evaluator contract versions, and relevant evaluation-time or validity epoch. The EXACT field set, ordering, and canonicalization rules require ADR-F13-003 resolution.
- **What this design provides:** The `SemanticDigest` struct declaration (algorithm, digest, coveredFields fields) for schema documentation only.
- **What this design does NOT provide:** No covered-field list publication (the prior revision's `SemanticDigestCoveredFields` variable is REMOVED as it was incomplete and misleading). No computation function. No fixtures that assert specific digest values.
- **Fixture treatment:** Fixtures referencing semantic digest use `"PLACEHOLDER_PENDING_ADR_F13_003"` and validate only structural presence (non-empty string when the field is present).

## 4. Data Models

### 4.1 DecisionRecord Envelope

```go
// Package core contains all shared decision/audit types.
// Package path: internal/decision/core

// DecisionRecord is the immutable governed conclusion (AC-01, AC-02, AC-03).
// Every successfully persisted DecisionRecord is FINAL and immutable.
// Canonical scope source: Metadata.ScopeRef (section 3.4).
type DecisionRecord struct {
    apimeta.TypeMeta                         // apiVersion: "decision.sovrunn.io/v1alpha1", kind: "DecisionRecord"
    Metadata    apimeta.ObjectMeta           `json:"metadata"`
    ProfileRef  ProfileRef                   `json:"profileRef"`
    Form        DecisionForm                 `json:"form"`
    Authority   DecisionAuthority            `json:"authority"`
    Purpose     string                       `json:"purpose"`
    RequestRef  *apimeta.TypedRef            `json:"requestRef,omitempty"`
    ActorRef    *apimeta.TypedRef            `json:"actorRef,omitempty"`
    SubjectRefs []apimeta.TypedRef           `json:"subjectRefs,omitempty"`
    EffectivePolicyContextRef *apimeta.TypedRef `json:"effectivePolicyContextRef,omitempty"`
    InputSnapshotRef *apimeta.TypedRef       `json:"inputSnapshotRef,omitempty"`
    EvaluationResultRefs []apimeta.TypedRef  `json:"evaluationResultRefs,omitempty"`
    Composition *CompositionMetadata         `json:"composition,omitempty"`
    Result      DecisionResult               `json:"result"`
    Relationship *DecisionRelationship       `json:"relationship,omitempty"`
    Provenance  Provenance                   `json:"provenance,omitempty"`
    Validity    Validity                     `json:"validity,omitempty"`
    Correlation Correlation                  `json:"correlation,omitempty"`
}

// ProfileRef identifies the versioned DecisionProfile governing this record.
type ProfileRef struct {
    Name    string `json:"name"`
    Version string `json:"version"`
}

// DecisionResult carries finality, optional adjudication, typed result, rationale, and obligations.
type DecisionResult struct {
    Finality     Finality           `json:"finality"`
    Adjudication *Adjudication      `json:"adjudication,omitempty"`
    TypedResult  json.RawMessage    `json:"typedResult,omitempty"`
    Rationale    *DecisionRationale `json:"rationale,omitempty"`
    Obligations  []Obligation       `json:"obligations,omitempty"`
}

// CompositionMetadata references the strategy and graph used for composition.
type CompositionMetadata struct {
    StrategyRef *StrategyRef      `json:"strategyRef,omitempty"`
    GraphRef    *apimeta.TypedRef `json:"graphRef,omitempty"`
}
```

Note: `DecisionRecord` does NOT have a separate `ScopeRef` field. Scope is exclusively `Metadata.ScopeRef` (see section 3.4).

### 4.2 Closed Vocabularies

```go
// DecisionForm is the closed primary-form vocabulary (AC-04).
type DecisionForm string

const (
    FormAdjudication    DecisionForm = "ADJUDICATION"
    FormSelection       DecisionForm = "SELECTION"
    FormRanking         DecisionForm = "RANKING"
    FormClassification  DecisionForm = "CLASSIFICATION"
    FormResolution      DecisionForm = "RESOLUTION"
    FormAllocation      DecisionForm = "ALLOCATION"
    FormPlan            DecisionForm = "PLAN"
    FormAssessment      DecisionForm = "ASSESSMENT"
    FormRecommendation  DecisionForm = "RECOMMENDATION"
)

// DecisionAuthority is orthogonal to form (AC-05).
type DecisionAuthority string

const (
    AuthorityAuthoritative  DecisionAuthority = "AUTHORITATIVE"
    AuthorityAdvisory       DecisionAuthority = "ADVISORY"
    AuthorityRecommendation DecisionAuthority = "RECOMMENDATION"
    AuthoritySimulation     DecisionAuthority = "SIMULATION"
)

// Adjudication is the closed adjudication vocabulary (AC-06).
type Adjudication string

const (
    AdjudicationAllowed         Adjudication = "ALLOWED"
    AdjudicationDenied          Adjudication = "DENIED"
    AdjudicationRequiresApproval Adjudication = "REQUIRES_APPROVAL"
)

// Finality states that every persisted DecisionRecord is FINAL (AC-03).
type Finality string

const (
    FinalityFinal Finality = "FINAL"
)
```

### 4.3 AuditScope (Shared Type in core/)

```go
// Package: internal/decision/core

// AuditScope identifies the governance scope of a decision or audit event.
// It uses a string Kind to accommodate ServiceInstance which is not in
// FEATURE-0012's ScopeKind enum (see ADR-F13-004 in section 17).
// This type does NOT claim to compose apimeta.ScopeKind.
type AuditScope struct {
    Kind string `json:"kind"` // closed set: Platform, Organization, OrganizationUnit, Tenant, Project, ServiceInstance
    UID  string `json:"uid,omitempty"`
    Name string `json:"name,omitempty"`
}

// AuditScopeKind constants for the six architecture-mandated scopes.
const (
    AuditScopePlatform         = "Platform"
    AuditScopeOrganization     = "Organization"
    AuditScopeOrganizationUnit = "OrganizationUnit"
    AuditScopeTenant           = "Tenant"
    AuditScopeProject          = "Project"
    AuditScopeServiceInstance  = "ServiceInstance"
)

// ValidAuditScopeKinds is the closed set of valid AuditScope.Kind values.
var ValidAuditScopeKinds = []string{
    AuditScopePlatform,
    AuditScopeOrganization,
    AuditScopeOrganizationUnit,
    AuditScopeTenant,
    AuditScopeProject,
    AuditScopeServiceInstance,
}

// IsValidAuditScopeKind reports whether kind is in the closed audit scope set.
func IsValidAuditScopeKind(kind string) bool {
    for _, v := range ValidAuditScopeKinds {
        if kind == v {
            return true
        }
    }
    return false
}
```

### 4.4 EvaluationResult

```go
// Package: internal/decision/core

// EvaluationResult is the atomic evaluator-specific normalized finding (AC-11, AC-12).
type EvaluationResult struct {
    apimeta.TypeMeta                          // apiVersion: "decision.sovrunn.io/v1alpha1", kind: "EvaluationResult"
    Metadata         apimeta.ObjectMeta       `json:"metadata"`
    EvaluatorType    string                   `json:"evaluatorType"`
    EvaluatorVersion string                   `json:"evaluatorVersion"`
    InputRef         *apimeta.TypedRef        `json:"inputRef,omitempty"`
    InputDigest      string                   `json:"inputDigest,omitempty"`
    Outcome          EvaluationOutcome        `json:"outcome"`
    Result           json.RawMessage          `json:"result,omitempty"`
    Timing           TimingEnvelope           `json:"timing"`
    TrustBoundary    string                   `json:"trustBoundary,omitempty"`
    Projection       *GovernedProjection      `json:"projection,omitempty"`
    SensitivityClass string                   `json:"sensitivityClass"`
    IntegrityDigest  string                   `json:"integrityDigest,omitempty"`
}

// EvaluationOutcome is the explicit outcome semantics (AC-11).
type EvaluationOutcome string

const (
    EvalOutcomeSuccess       EvaluationOutcome = "SUCCESS"
    EvalOutcomeNonMatch      EvaluationOutcome = "NON_MATCH"
    EvalOutcomeIndeterminate EvaluationOutcome = "INDETERMINATE"
    EvalOutcomeTimeout       EvaluationOutcome = "TIMEOUT"
    EvalOutcomeError         EvaluationOutcome = "ERROR"
)

// TimingEnvelope captures evaluation time boundaries.
type TimingEnvelope struct {
    StartedAt   string `json:"startedAt"`
    CompletedAt string `json:"completedAt"`
    DurationMs  int64  `json:"durationMs,omitempty"`
}

// GovernedProjection is the bounded, authorized output projection (AC-13).
type GovernedProjection struct {
    Content             json.RawMessage  `json:"content"`
    SizeBytes           int              `json:"sizeBytes"`
    FieldCount          int              `json:"fieldCount"`
    SensitivityClass    string           `json:"sensitivityClass"`
    SourceScope         apimeta.ScopeRef `json:"sourceScope"`
    ExecutionConfigRef  *apimeta.TypedRef `json:"executionConfigRef,omitempty"`
    SafetyFilterDecls   []string         `json:"safetyFilterDecls,omitempty"`
}
```

### 4.5 DecisionProfile (Complete AC-08 Model)

```go
// Package: internal/decision/core

// DecisionProfile defines one decision family's full contract (AC-08, AC-09).
type DecisionProfile struct {
    apimeta.TypeMeta                              // apiVersion: "decision.sovrunn.io/v1alpha1", kind: "DecisionProfile"
    Metadata           apimeta.ObjectMeta         `json:"metadata"`
    Family             string                     `json:"family"`
    Version            string                     `json:"version"`
    Owner              string                     `json:"owner"`
    Lifecycle          ProfileLifecycle           `json:"lifecycle"`

    // Form and authority semantics
    PrimaryForm        DecisionForm               `json:"primaryForm"`
    AllowedFacets      []DecisionForm             `json:"allowedFacets,omitempty"`
    AllowedAuthority   []DecisionAuthority        `json:"allowedAuthority"`

    // Schemas for typed extensibility
    InputSchema        json.RawMessage            `json:"inputSchema,omitempty"`
    ResultSchema       json.RawMessage            `json:"resultSchema,omitempty"`
    RationaleSchema    json.RawMessage            `json:"rationaleSchema,omitempty"`
    ObligationSchema   json.RawMessage            `json:"obligationSchema,omitempty"`

    // Evaluator and composition
    AcceptedEvaluators []EvaluatorRegistration    `json:"acceptedEvaluators,omitempty"`
    CompositionStrategies []string                `json:"compositionStrategies,omitempty"`

    // Scope, actor, subject, policy, evidence, audit semantics (AC-08)
    ScopeSemantics     ScopeSemantics             `json:"scopeSemantics"`
    ActorSemantics     ActorSemantics             `json:"actorSemantics"`
    SubjectSemantics   SubjectSemantics           `json:"subjectSemantics"`
    PolicySemantics    PolicySemantics            `json:"policySemantics"`
    EvidenceSemantics  EvidenceSemantics          `json:"evidenceSemantics"`
    AuditSemantics     AuditSemantics             `json:"auditSemantics"`

    // Determinism, failure, timeout behavior
    Determinism        DeterminismPolicy          `json:"determinism"`
    FailureBehavior    FailureBehavior            `json:"failureBehavior"`

    // Validity, expiry, correction, supersession, revocation rules
    ValidityRules      ValidityRules              `json:"validityRules"`

    // Sensitivity, residency, retention, legal-hold, projection profiles (AC-08)
    SensitivityClass   string                     `json:"sensitivityClass"`
    ResidencyRules     ResidencyRules             `json:"residencyRules,omitempty"`
    RetentionRules     RetentionRules             `json:"retentionRules,omitempty"`
    ProjectionProfiles []ProjectionProfile        `json:"projectionProfiles,omitempty"`

    // Idempotency identity
    IdempotencyRules   IdempotencyRules           `json:"idempotencyRules"`

    // Bounded limits — conditionally mandatory (see section 4.5.2)
    GraphLimits        *GraphLimits               `json:"graphLimits,omitempty"`
    ProjectionBounds   *ProjectionBounds          `json:"projectionBounds,omitempty"`

    // Compatibility, fixtures, reassessment
    CompatibilityPolicy CompatibilityPolicy       `json:"compatibilityPolicy"`
    FixtureRefs        []string                   `json:"fixtureRefs,omitempty"`
    ReassessmentTriggers []string                 `json:"reassessmentTriggers,omitempty"`
}
```

#### 4.5.1 Profile Supporting Types

```go
// ProfileLifecycle is the lifecycle status of a DecisionProfile.
type ProfileLifecycle string

const (
    ProfileDraft      ProfileLifecycle = "Draft"
    ProfileActive     ProfileLifecycle = "Active"
    ProfileDeprecated ProfileLifecycle = "Deprecated"
    ProfileRevoked    ProfileLifecycle = "Revoked"
)

// EvaluatorRegistration declares an accepted evaluator type with constraints.
type EvaluatorRegistration struct {
    Type               string `json:"type"`
    MinVersion         string `json:"minVersion,omitempty"`
    MaxVersion         string `json:"maxVersion,omitempty"`
    TrustBoundary      string `json:"trustBoundary,omitempty"`
    RequireSigning     bool   `json:"requireSigning,omitempty"`
}

// ScopeSemantics declares how scope applies to this profile (AC-08: scope semantics).
type ScopeSemantics struct {
    AllowedScopes []string `json:"allowedScopes"` // values from closed AuditScopeKind set
    RequireScope  bool     `json:"requireScope"`
}

// ActorSemantics declares actor requirements for this profile.
type ActorSemantics struct {
    RequireActor         bool     `json:"requireActor"`
    AllowDelegation      bool     `json:"allowDelegation"`
    AllowImpersonation   bool     `json:"allowImpersonation"`
    RequireAssuranceLevel bool    `json:"requireAssuranceLevel"`
    // ARCHITECTURE_DECISION_REQUIRED: ADR-F13-005
    // Subject: Exact actor-assurance-level vocabulary
}

// SubjectSemantics declares subject constraints.
type SubjectSemantics struct {
    RequireSubject bool `json:"requireSubject"`
    MaxSubjects    int  `json:"maxSubjects"`
    AllowedKinds   []string `json:"allowedKinds,omitempty"`
}

// PolicySemantics declares policy-context requirements.
type PolicySemantics struct {
    RequireEffectivePolicyContext bool   `json:"requireEffectivePolicyContext"`
    MaxStalenessMs               int64  `json:"maxStalenessMs,omitempty"`
    FailOnMissingContext         bool   `json:"failOnMissingContext"`
}

// EvidenceSemantics declares evidence requirements.
type EvidenceSemantics struct {
    RequireEvidence      bool     `json:"requireEvidence"`
    AcceptedEvidenceKinds []string `json:"acceptedEvidenceKinds,omitempty"`
    RequireDigest        bool     `json:"requireDigest"`
}

// AuditSemantics declares audit linkage requirements.
type AuditSemantics struct {
    RequireDurableRecord  bool   `json:"requireDurableRecord"`
    AggregationEligible   bool   `json:"aggregationEligible"`
    MaxAggregationWindowMs int64 `json:"maxAggregationWindowMs,omitempty"`
    MinAuditEvidencePerWindow string `json:"minAuditEvidencePerWindow,omitempty"`
}

// ResidencyRules declares data residency constraints.
type ResidencyRules struct {
    RequireResidencyDeclaration bool     `json:"requireResidencyDeclaration"`
    AllowedJurisdictions        []string `json:"allowedJurisdictions,omitempty"`
    // ARCHITECTURE_DECISION_REQUIRED: ADR-F13-006
}

// RetentionRules declares retention and legal-hold constraints.
type RetentionRules struct {
    RequireRetentionClass bool   `json:"requireRetentionClass"`
    AllowLegalHold        bool   `json:"allowLegalHold"`
    // ARCHITECTURE_DECISION_REQUIRED: ADR-F13-007
}

// DeterminismPolicy declares whether the profile allows nondeterministic evaluators.
type DeterminismPolicy struct {
    RequireDeterministic bool   `json:"requireDeterministic"`
    AllowAI             bool   `json:"allowAI"`
    ReplayPolicy        string `json:"replayPolicy"` // "FROM_CAPTURED" or "RERUN_PROHIBITED"
}

// FailureBehavior declares how the profile handles failures.
type FailureBehavior struct {
    FailPosture          FailPosture `json:"failPosture"`
    TimeoutBehavior      string      `json:"timeoutBehavior"`
    InsufficientEvidence string      `json:"insufficientEvidence"`
}

// ValidityRules declares decision expiry and correction semantics.
type ValidityRules struct {
    DefaultValidityMs int64  `json:"defaultValidityMs,omitempty"`
    AllowCorrection   bool   `json:"allowCorrection"`
    AllowSupersession bool   `json:"allowSupersession"`
    AllowRevocation   bool   `json:"allowRevocation"`
    MaxChainDepth     int    `json:"maxChainDepth"`
}

// IdempotencyRules declares idempotency identity semantics.
type IdempotencyRules struct {
    RequireIdempotencyKey bool `json:"requireIdempotencyKey"`
}

// ProjectionProfile declares consumer-specific projection rules (AC-44).
type ProjectionProfile struct {
    Consumer       string   `json:"consumer"`
    RedactedFields []string `json:"redactedFields,omitempty"`
    Reason         string   `json:"reason,omitempty"`
}

// CompatibilityPolicy declares version compatibility rules.
type CompatibilityPolicy struct {
    MinReaderVersion string `json:"minReaderVersion,omitempty"`
    BackwardCompat   bool   `json:"backwardCompat"`
}
```

#### 4.5.2 GraphLimits: Conditional Mandatory Requirement

`GraphLimits` is conditionally mandatory on `DecisionProfile`:

**Rule:** `GraphLimits` MUST be non-nil and have all positive values when ANY of the following conditions hold:
1. `CompositionStrategies` is non-empty (profile permits composition)
2. `AcceptedEvaluators` has more than one entry (profile permits multiple evaluations)
3. Profile declares a dependency graph in any fixture or conformance scenario

**Rule:** `GraphLimits` MAY be nil when ALL of the following hold:
1. `CompositionStrategies` is empty
2. `AcceptedEvaluators` has zero or one entry
3. Profile does not reference composition or graph-based evaluation

**Validation behavior:**
- When conditions for mandatory GraphLimits are met and `GraphLimits` is nil → `VALIDATION_FAILED` with message "graphLimits required when profile permits composition or multiple evaluators"
- When `GraphLimits` is non-nil, all numeric fields must be > 0 → `VALIDATION_FAILED` if any field ≤ 0
- No system-wide default values are invented. Each profile declares its own bounds.

```go
// GraphLimits declares the profile-specific finite bounds for composition (DQ-03).
// Conditionally mandatory: required when profile permits composition or multiple evaluators.
// All fields must be positive (> 0) when the struct is present.
type GraphLimits struct {
    MaxNodes              int   `json:"maxNodes"`
    MaxDepth              int   `json:"maxDepth"`
    MaxWidth              int   `json:"maxWidth"`
    MaxFanOut             int   `json:"maxFanOut"`
    MaxPayloadBytes       int64 `json:"maxPayloadBytes"`
    MaxWallClockMs        int64 `json:"maxWallClockMs"`
    MaxRemoteCalls        int   `json:"maxRemoteCalls"`
    MaxAuthorities        int   `json:"maxAuthorities"`
    MaxConcurrentEvals    int   `json:"maxConcurrentEvals"`
    MaxPerEvaluatorTimeoutMs int64 `json:"maxPerEvaluatorTimeoutMs"`
    MaxRetries            int   `json:"maxRetries"`
    MaxCrossJurisdiction  int   `json:"maxCrossJurisdictionCalls"`
    MaxCostUnits          int64 `json:"maxCostUnits"`
}

// ProjectionBounds declares per-profile bounds for governed evaluator output (AC-13.2).
type ProjectionBounds struct {
    MaxFieldCount        int      `json:"maxFieldCount"`
    MaxTotalBytes        int      `json:"maxTotalBytes"`
    AllowedValueTypes    []string `json:"allowedValueTypes"`
    ProhibitedCategories []string `json:"prohibitedCategories"`
}
```

### 4.6 Composition and Dependency Graph

```go
// Package: internal/decision/core

// StrategyRef identifies a versioned composition strategy.
type StrategyRef struct {
    Name    string `json:"name"`
    Version string `json:"version"`
}

// CompositionStrategy declares a registered, deterministic aggregation algorithm (AC-14).
type CompositionStrategy struct {
    Name              string      `json:"name"`
    Version           string      `json:"version"`
    AcceptedInputTypes []string   `json:"acceptedInputTypes"`
    Ordering          string      `json:"ordering"`
    Precedence        string      `json:"precedence,omitempty"`
    ShortCircuit      bool        `json:"shortCircuit"`
    MissingBehavior   string      `json:"missingBehavior"`
    TimeoutBehavior   string      `json:"timeoutBehavior"`
    ConflictBehavior  string      `json:"conflictBehavior"`
    ErrorBehavior     string      `json:"errorBehavior"`
    FailPosture       FailPosture `json:"failPosture"`
    DeterministicOutput bool      `json:"deterministicOutput"`
    MaxFanOut         int         `json:"maxFanOut"`
    ExplanationRules  string      `json:"explanationRules,omitempty"`
}

// FailPosture declares the fail-closed or fail-open posture.
type FailPosture string

const (
    FailPostureClosed FailPosture = "FAIL_CLOSED"
    FailPostureOpen   FailPosture = "FAIL_OPEN"
)

// DependencyGraph is the bounded DAG for complex composition (AC-15, AC-16).
type DependencyGraph struct {
    ID         string      `json:"id"`
    Version    string      `json:"version"`
    RootNodeID string      `json:"rootNodeId"`
    Nodes      []GraphNode `json:"nodes"`
    Edges      []GraphEdge `json:"edges"`
}

// GraphNodeType classifies nodes in the dependency graph.
type GraphNodeType string

const (
    NodeInput       GraphNodeType = "INPUT"
    NodeEvaluation  GraphNodeType = "EVALUATION"
    NodeComposition GraphNodeType = "COMPOSITION"
    NodeDecision    GraphNodeType = "DECISION"
)

// GraphNode is a typed node in the bounded DAG.
type GraphNode struct {
    ID              string        `json:"id"`
    Type            GraphNodeType `json:"type"`
    EvaluatorType   string        `json:"evaluatorType,omitempty"`
    StrategyRef     *StrategyRef  `json:"strategyRef,omitempty"`
    ParallelGroup   string        `json:"parallelGroup,omitempty"`
}

// GraphEdge is a directed edge in the DAG.
type GraphEdge struct {
    FromNodeID string `json:"fromNodeId"`
    ToNodeID   string `json:"toNodeId"`
}
```

### 4.7 Rationale, Obligations, and Provenance

```go
// Package: internal/decision/core

// DecisionRationale embeds structured reasons (AC-01).
type DecisionRationale struct {
    ReasonCodes   []string `json:"reasonCodes,omitempty"`
    Reasons       []string `json:"reasons,omitempty"`
    Alternatives  []string `json:"alternatives,omitempty"`
    Warnings      []string `json:"warnings,omitempty"`
    Corrective    []string `json:"corrective,omitempty"`
}

// Obligation is an enforceable condition tied to the decision (AC-08).
type Obligation struct {
    Type       ObligationType `json:"type"`
    ID         string         `json:"id"`
    Mandatory  bool           `json:"mandatory"`
    Condition  string         `json:"condition,omitempty"`
    Deadline   string         `json:"deadline,omitempty"`
}

// ObligationType classifies obligations.
type ObligationType string

const (
    ObligationEnforce    ObligationType = "ENFORCE"
    ObligationNotify     ObligationType = "NOTIFY"
    ObligationLog        ObligationType = "LOG"
    ObligationReview     ObligationType = "REVIEW"
    ObligationExpiry     ObligationType = "EXPIRY"
)

// Provenance records creation context.
type Provenance struct {
    CreatedBy       string `json:"createdBy,omitempty"`
    CreatedAt       string `json:"createdAt,omitempty"`
    ProfileDigest   string `json:"profileDigest,omitempty"`
    StrategyDigest  string `json:"strategyDigest,omitempty"`
}

// Validity defines expiry and scope of the decision.
type Validity struct {
    ValidFrom   string `json:"validFrom,omitempty"`
    ValidUntil  string `json:"validUntil,omitempty"`
    Revocable   bool   `json:"revocable,omitempty"`
}

// Correlation links the decision to operations, audit events, and traces.
type Correlation struct {
    OperationRef   *apimeta.TypedRef  `json:"operationRef,omitempty"`
    AuditEventRefs []apimeta.TypedRef `json:"auditEventRefs,omitempty"`
    TraceID        string             `json:"traceId,omitempty"`
    SpanID         string             `json:"spanId,omitempty"`
    RequestID      string             `json:"requestId,omitempty"`
}
```

### 4.8 Relationship and Effective State

```go
// Package: internal/decision/core

// RelationshipType for correction, supersession, and revocation (AC-18).
type RelationshipType string

const (
    RelCorrects   RelationshipType = "CORRECTS"
    RelSupersedes RelationshipType = "SUPERSEDES"
    RelRevokes    RelationshipType = "REVOKES"
)

// DecisionRelationship links to a predecessor record.
// Predecessor kind constraint: predecessorRef.Kind MUST be "DecisionRecord".
type DecisionRelationship struct {
    Type           RelationshipType `json:"type"`
    PredecessorRef apimeta.TypedRef `json:"predecessorRef"`
}

// EffectiveState is projection metadata (AC-19). It is computed, not stored.
type EffectiveState string

const (
    StateEffective  EffectiveState = "EFFECTIVE"
    StateSuperseded EffectiveState = "SUPERSEDED"
    StateRevoked    EffectiveState = "REVOKED"
)
```

**Relationship resolution semantics:**
- A DecisionRecord with no `Relationship` field (nil) is a root record with `EffectiveState = EFFECTIVE`.
- A DecisionRecord with `Relationship.Type = SUPERSEDES` causes its predecessor to have computed `EffectiveState = SUPERSEDED`.
- A DecisionRecord with `Relationship.Type = REVOKES` causes its predecessor to have computed `EffectiveState = REVOKED`.
- A DecisionRecord with `Relationship.Type = CORRECTS` creates a correction chain; the predecessor retains its effective state, and the correction is the authoritative version for audit purposes.
- Chain depth is validated against `profile.ValidityRules.MaxChainDepth`. If the chain (counting all predecessor links recursively) exceeds MaxChainDepth, validation rejects with `CHAIN_DEPTH_EXCEEDED`.
- Chains must be cycle-free: a record cannot transitively reference itself as a predecessor.
- Fork resolution: when multiple records reference the same predecessor (forked chain), effective state is determined by creation timestamp ordering (latest creation wins for SUPERSEDES/REVOKES). This is a read-model projection concern, not a mutation.

**Predecessor kind constraint:** `DecisionRelationship.PredecessorRef.Kind` MUST equal `"DecisionRecord"`. Validation rejects any other kind value.

### 4.9 Six-Scope AuditEvent

```go
// Package: internal/decision/audit

// AuditEvent extends the FEATURE-0012 alpha to six governance scopes (AC-24).
// apiVersion: "governance.sovrunn.io/v1alpha1", kind: "AuditEvent"
type AuditEvent struct {
    apimeta.TypeMeta
    Metadata apimeta.ObjectMeta `json:"metadata"`
    Record   AuditEventRecord   `json:"record"`
}

// AuditEventRecord is the append-only audit payload with six-scope support.
type AuditEventRecord struct {
    ActorRef               apimeta.TypedRef   `json:"actorRef"`
    RequestID              string             `json:"requestId"`
    OperationRef           *apimeta.TypedRef  `json:"operationRef,omitempty"`
    DecisionRef            *apimeta.TypedRef  `json:"decisionRef,omitempty"`
    SubjectRef             apimeta.TypedRef   `json:"subjectRef"`
    SubjectResourceVersion string             `json:"subjectResourceVersion,omitempty"`
    Action                 string             `json:"action"`
    Outcome                AuditOutcome       `json:"outcome"`
    ReasonCode             string             `json:"reasonCode"`
    Scope                  core.AuditScope    `json:"scope"`
    ParentEventRef         *apimeta.TypedRef  `json:"parentEventRef,omitempty"`
    TraceID                string             `json:"traceId,omitempty"`
    CorrectionOfRef        *apimeta.TypedRef  `json:"correctionOfRef,omitempty"`
}

// AuditOutcome is the closed audit outcome vocabulary (AC-25).
type AuditOutcome string

const (
    AuditOutcomeSucceeded AuditOutcome = "Succeeded"
    AuditOutcomeDenied    AuditOutcome = "Denied"
    AuditOutcomeFailed    AuditOutcome = "Failed"
)
```

Note: `AuditEventRecord.Scope` uses `core.AuditScope` (imported from `internal/decision/core/`). This is the audit-specific scope discriminator. `AuditEvent.Metadata.ScopeRef` provides FEATURE-0012 compatibility. Both must satisfy the consistency invariant in section 3.5.3.

### 4.10 Atomic Acceptance Contract

```go
// Package: internal/decision/core

// AuditObligation is the preallocated audit identity accepted atomically
// with the decision (AC-22, AC-23).
type AuditObligation struct {
    AuditEventID string       `json:"auditEventId"`
    Scope        AuditScope   `json:"scope"`
    DecisionRef  apimeta.TypedRef `json:"decisionRef"`
    AcceptedAt   string       `json:"acceptedAt"`
}

// AcceptanceResult models the atomic acceptance outcome (DQ-07).
// Either both Decision and Obligation are non-nil, or both are nil and Err is set.
type AcceptanceResult struct {
    Decision   *DecisionRecord  `json:"decision,omitempty"`
    Obligation *AuditObligation `json:"obligation,omitempty"`
    Err        error            `json:"-"`
}

// Acceptor models the atomic decision/audit-obligation acceptance contract (AC-22).
// This is a pure-function contract for conformance testing, not a production service.
type Acceptor interface {
    Accept(record DecisionRecord, scope AuditScope) AcceptanceResult
}
```

Note: `AuditObligation` and `AcceptanceResult` reference `AuditScope` which is in the same `core` package. No cross-package circular dependency exists.

### 4.11 Identity and Digest (Declaration Only)

```go
// Package: internal/decision/core

// IdempotencyKey is the caller-supplied retry identity (AC-30).
type IdempotencyKey struct {
    CallerID string `json:"callerId"`
    Purpose  string `json:"purpose"`
    Key      string `json:"key"`
}

// SemanticDigest is the semantic decision identity for deduplication (AC-30).
// IMPORTANT: Digest computation is NOT implemented in this feature.
// ADR-F13-003 must resolve the exact canonicalization algorithm/version
// and complete covered-field set before implementation.
// The prior revision's SemanticDigestCoveredFields variable is REMOVED
// because it was incomplete relative to the architecture's requirements.
type SemanticDigest struct {
    Algorithm     string   `json:"algorithm"`     // candidate: "SHA-256" pending ADR-F13-003
    Digest        string   `json:"digest"`        // hex-encoded; PLACEHOLDER_PENDING_ADR_F13_003 in fixtures
    CoveredFields []string `json:"coveredFields"` // populated by ADR-F13-003 resolution; empty until then
}
```

**ADR-F13-003 scope:** The architecture (section 9) requires semantic digest coverage of: scope, purpose, profile name+version, authority, canonical input snapshot, EffectivePolicyContext, composition strategy/version, evaluator contract versions, and relevant evaluation-time or validity epoch. This design does NOT publish a partial or complete covered-field list because ADR-F13-003 must also resolve the canonicalization algorithm, field ordering, and byte representation. Publishing a partial list would be misleading and would falsely constrain future resolution.

Fixtures that reference semantic digest use `PLACEHOLDER_PENDING_ADR_F13_003` as the digest value and validate only structural presence (non-empty string when required), not correctness of the hash or completeness of covered fields.

## 5. Interfaces

### 5.1 Profile Registry Interface

```go
// Package: internal/decision/ (root)

// ProfileRegistry provides access to registered DecisionProfiles (AC-10, AC-53).
type ProfileRegistry interface {
    Get(name, version string) (*core.DecisionProfile, error)
    Register(profile core.DecisionProfile) error
    List() []core.DecisionProfile
}
```

### 5.2 Composition Strategy Registry Interface

```go
// Package: internal/decision/ (root)

// StrategyRegistry provides access to registered composition strategies (AC-14).
type StrategyRegistry interface {
    Get(name, version string) (*core.CompositionStrategy, error)
    Register(strategy core.CompositionStrategy) error
    List() []core.CompositionStrategy
}
```

### 5.3 Graph Validator Interface

```go
// Package: internal/decision/ (root)

// GraphValidator validates a DependencyGraph against profile-declared limits (AC-15, AC-16).
type GraphValidator interface {
    Validate(graph core.DependencyGraph, limits core.GraphLimits) []apiproblem.Violation
}
```

### 5.4 Projection Validator Interface

```go
// Package: internal/decision/ (root)

// ProjectionValidator validates evaluator output projections (AC-13).
type ProjectionValidator interface {
    Validate(projection core.GovernedProjection, bounds core.ProjectionBounds, targetScope apimeta.ScopeRef) []apiproblem.Violation
}
```

### 5.5 Bundle Verifier Interface

```go
// Package: internal/decision/bundle/

// BundleVerifier verifies signed profile bundles for offline operation (AC-31, AC-53).
type BundleVerifier interface {
    Verify(bundle ProfileBundle, contentProvider ContentProvider) error
}

// ContentProvider resolves bundle entry content for digest verification.
// For conformance testing, content is provided in-memory.
type ContentProvider interface {
    Content(path string) ([]byte, error)
}

// TrustRootStore provides access to known trust roots and their revocation status.
type TrustRootStore interface {
    IsKnown(id string) bool
    IsRevoked(id string) bool
    PublicKey(id string) ([]byte, error)
}
```

### 5.6 Relationship Chain Validator Interface

```go
// Package: internal/decision/ (root)

// ChainValidator validates relationship chains against profile-declared depth limits (AC-18, AC-19).
type ChainValidator interface {
    ValidateChain(record core.DecisionRecord, profile core.DecisionProfile, predecessorLookup func(apimeta.TypedRef) (*core.DecisionRecord, error)) []apiproblem.Violation
}
```

### 5.7 Protected-Content Integrity Verifier Interface

```go
// Package: internal/decision/ (root)

// IntegrityVerifier validates protected-content integrity digests (AC-13.7f, AC-13.7g).
type IntegrityVerifier interface {
    // VerifyFormat checks that the digest string matches "sha256:<64 hex chars>".
    VerifyFormat(digest string) error
    // VerifyContent computes sha256 of contentBytes and compares to declared digest.
    // Returns DIGEST_MISMATCH error on mismatch.
    VerifyContent(digest string, contentBytes []byte) error
}
```

## 6. Validation Rules

### 6.1 DecisionRecord Envelope Validation

| Field | Rule | Error Code |
|---|---|---|
| profileRef.name | Required, non-empty, max 128 chars, DNS-like pattern | VALIDATION_FAILED |
| profileRef.version | Required, semver format | VALIDATION_FAILED |
| form | Required, must be in closed vocabulary | VALIDATION_FAILED |
| authority | Required, must be in closed vocabulary | VALIDATION_FAILED |
| authority + AI origin | AI output must not be AUTHORITATIVE (AC-07) | AUTHORITY_MISMATCH |
| result.adjudication | Required only if profile declares ADJUDICATION facet (AC-06) | VALIDATION_FAILED |
| result.adjudication presence | Must be absent if profile has no ADJUDICATION facet (EC-01) | VALIDATION_FAILED |
| result.finality | Must be FINAL (AC-03) | VALIDATION_FAILED |
| metadata.scopeRef | Must match profile's ScopeSemantics.AllowedScopes; required if ScopeSemantics.RequireScope is true | SCOPE_MISMATCH |
| subjectRefs | Max count bounded by profile SubjectSemantics.MaxSubjects | VALIDATION_FAILED |
| evaluationResultRefs | Max count bounded per profile GraphLimits.MaxConcurrentEvals (when GraphLimits present) | GRAPH_LIMIT_EXCEEDED |
| obligations | Unknown/unsupported mandatory obligation → unenforceable (EC-08) | OBLIGATION_UNENFORCEABLE |
| relationship.predecessorRef.Kind | Must be "DecisionRecord" when relationship is present | REFERENCE_INVALID |
| relationship.type | Must be in closed RelationshipType vocabulary when relationship is present | VALIDATION_FAILED |

### 6.2 EvaluationResult Validation

| Field | Rule | Error Code |
|---|---|---|
| evaluatorType | Required, non-empty | VALIDATION_FAILED |
| outcome | Required, closed vocabulary | VALIDATION_FAILED |
| sensitivityClass | Required, non-empty (AC-13.7h) | VALIDATION_FAILED |
| projection.sizeBytes | Must not exceed ProjectionBounds.MaxTotalBytes (AC-13.7b) | PROJECTION_EXCEEDED |
| projection.fieldCount | Must not exceed ProjectionBounds.MaxFieldCount (AC-13.7b) | PROJECTION_EXCEEDED |
| projection.content | Must not contain secret/credential patterns (AC-13.7c) | PROJECTION_EXCEEDED |
| projection.content | Must not contain raw prompts or system instructions (AC-13.7d) | PROJECTION_EXCEEDED |
| projection.sourceScope | Must not disclose cross-tenant content (AC-13.5, AC-13.7e) | SCOPE_MISMATCH |
| integrityDigest (when inputRef non-nil) | Must be non-empty and match `^sha256:[0-9a-f]{64}$` (AC-13.7f) | DIGEST_MISMATCH |
| integrityDigest (when content bytes available) | Must match sha256 of referenced content (AC-13.7g) | DIGEST_MISMATCH |

### 6.3 Graph Validation

| Rule | Error Code |
|---|---|
| Graph is acyclic | GRAPH_LIMIT_EXCEEDED |
| Single root node of type DECISION | VALIDATION_FAILED |
| Node count ≤ limits.MaxNodes | GRAPH_LIMIT_EXCEEDED |
| Depth ≤ limits.MaxDepth | GRAPH_LIMIT_EXCEEDED |
| Width ≤ limits.MaxWidth | GRAPH_LIMIT_EXCEEDED |
| Fan-out ≤ limits.MaxFanOut | GRAPH_LIMIT_EXCEEDED |
| All node IDs unique and non-empty | VALIDATION_FAILED |
| All edges reference existing nodes | VALIDATION_FAILED |
| Deterministic topological order | VALIDATION_FAILED |

### 6.4 Profile Validation

| Rule | Error Code |
|---|---|
| Name/version unique in registry | RESOURCE_ALREADY_EXISTS |
| PrimaryForm in closed vocabulary | VALIDATION_FAILED |
| AllowedAuthority subset of closed vocabulary | VALIDATION_FAILED |
| GraphLimits mandatory when CompositionStrategies non-empty OR AcceptedEvaluators > 1 | VALIDATION_FAILED |
| GraphLimits has all positive values when present | VALIDATION_FAILED |
| Lifecycle in closed vocabulary | VALIDATION_FAILED |
| FailureBehavior.FailPosture FAIL_CLOSED for security profiles (SEC-12) | VALIDATION_FAILED |
| ValidityRules.MaxChainDepth > 0 when AllowCorrection or AllowSupersession or AllowRevocation is true | VALIDATION_FAILED |
| SubjectSemantics.MaxSubjects > 0 when RequireSubject is true | VALIDATION_FAILED |

### 6.5 Relationship Chain Validation

| Rule | Error Code |
|---|---|
| Chain is cycle-free | CHAIN_DEPTH_EXCEEDED |
| Chain depth ≤ profile.ValidityRules.MaxChainDepth | CHAIN_DEPTH_EXCEEDED |
| Predecessor exists (resolvable) | REFERENCE_INVALID |
| Predecessor kind is "DecisionRecord" | REFERENCE_INVALID |
| Relationship type is in closed vocabulary | VALIDATION_FAILED |
| Profile permits the relationship type (AllowCorrection/AllowSupersession/AllowRevocation) | VALIDATION_FAILED |

### 6.6 Cross-Scope Safety Validation (AC-13.5)

| Rule | Error Code |
|---|---|
| GovernedProjection.SourceScope.UID must equal target DecisionRecord's Metadata.ScopeRef UID (when both are non-platform scope) | SCOPE_MISMATCH |
| Platform-scoped projections may only enter platform-scoped decisions | SCOPE_MISMATCH |
| Projection content must not contain identifiers (UIDs, names) from a different scope | SCOPE_MISMATCH |

### 6.7 AuditEvent Scope Consistency Validation

| Rule | Error Code |
|---|---|
| Record.Scope.Kind must be in closed six-value set | VALIDATION_FAILED |
| When Record.Scope.Kind ∈ {Org, OU, Tenant, Project}: Metadata.ScopeRef.Kind must match | SCOPE_MISMATCH |
| When Record.Scope.Kind == "Platform": Metadata.ScopeRef must be nil | SCOPE_MISMATCH |
| When Record.Scope.Kind == "ServiceInstance": Metadata.ScopeRef must be non-nil and Kind ∈ {Tenant, Project} | SCOPE_MISMATCH |
| When Record.Scope.UID non-empty and overlapping kind: Metadata.ScopeRef.UID must match | SCOPE_MISMATCH |

### 6.8 Validation Error Behavior

All validation functions return `[]apiproblem.Violation`. An empty slice means valid. Validation errors MUST NOT expose sensitive projection data content in error messages — only field paths, bound values, and generic violation descriptions.

## 7. Registry and Storage Design

### 7.1 In-Memory Profile Registry

The profile registry is an in-memory Go map for conformance testing. It is not a production service.

```go
// Package: internal/decision/ (root)

type inMemoryProfileRegistry struct {
    mu       sync.RWMutex
    profiles map[string]*core.DecisionProfile // key: "name:version"
}
```

Registration validates:
- Profile struct completeness and field validity
- Name+version uniqueness
- Lifecycle state is Draft or Active (cannot register Deprecated/Revoked)
- GraphLimits conditional mandatory rule (section 4.5.2)
- GraphLimits positive bounds when present
- FailureBehavior conformance for security-classified profiles
- ValidityRules.MaxChainDepth > 0 when relationship types are allowed

### 7.2 In-Memory Strategy Registry

```go
// Package: internal/decision/ (root)

type inMemoryStrategyRegistry struct {
    mu         sync.RWMutex
    strategies map[string]*core.CompositionStrategy // key: "name:version"
}
```

### 7.3 Profile Bundle Format (JSON)

The signed bundle is a JSON manifest with embedded or referenced profiles. JSON is used because `encoding/json` is in the Go standard library.

```json
{
  "bundleVersion": "1.0.0",
  "trustRoot": {
    "id": "sovrunn-root-2026",
    "algorithm": "Ed25519",
    "publicKeyRef": "keys/sovrunn-root-2026.pub"
  },
  "entries": [
    {
      "name": "AuthorizationDecision",
      "version": "1.0.0",
      "contentDigest": "sha256:abcdef0123456789...",
      "path": "profiles/authorization-decision-1.0.0.json"
    }
  ],
  "signature": {
    "algorithm": "Ed25519",
    "value": "base64-encoded-signature",
    "signedAt": "2026-07-27T12:00:00Z"
  },
  "validityMs": 86400000
}
```

Bundle verification (pure function, fully specified):

1. Parse manifest JSON (standard library `encoding/json`)
2. Verify trust-root identity is known via `TrustRootStore.IsKnown(id)` → unknown → `PROFILE_REJECTED`
3. Verify trust-root is not revoked via `TrustRootStore.IsRevoked(id)` → revoked → `PROFILE_REJECTED`
4. Verify bundle not expired: `time.Now()` > `parse(signature.signedAt) + validityMs` → expired → `PROFILE_REJECTED`
5. Construct signed input: manifest JSON with `signature.value` replaced by `""`, serialized as compact JSON (no extra whitespace, keys in struct-declaration order via `encoding/json` default marshaling)
6. Retrieve public key via `TrustRootStore.PublicKey(id)`
7. Verify Ed25519 signature: `ed25519.Verify(publicKey, signedInput, base64Decode(signature.value))` → invalid → `PROFILE_REJECTED`
8. For each entry: retrieve content via `ContentProvider.Content(entry.path)`, compute `"sha256:" + hex(sha256(content))`, compare to `entry.contentDigest` → mismatch → `DIGEST_MISMATCH`
9. Return nil error on success

**Error behavior:** Bundle verification returns a single error value. On first failure, verification halts and returns the specific error. Partial success is not reported.

### 7.4 Fixture Data Format

Conformance fixture test data uses JSON format exclusively. JSON is parsed with Go standard library `encoding/json`. No YAML parsing is required for fixture data.

Fixture files are stored in `internal/decision/fixtures/testdata/*.json`.

## 8. Operation and Audit Behavior

### 8.1 Decision-to-Audit Linkage (AC-21)

Every `DecisionRecord` links to at least one `AuditEvent` via `correlation.auditEventRefs`. The audit event identity is preallocated at acceptance time (AC-23).

Flow (contract, not runtime implementation):
1. Generate immutable audit-event UID
2. Store UID in `DecisionRecord.Correlation.AuditEventRefs`
3. Store UID in `AuditObligation.AuditEventID`
4. Atomic acceptance: both records accepted or neither
5. Asynchronous materialization: creates full `AuditEvent` using the preallocated UID

### 8.2 Audit Outcome Semantics (AC-25)

| Situation | Adjudication | Audit Outcome |
|---|---|---|
| Successfully produced ALLOWED decision | ALLOWED | Succeeded |
| Successfully produced DENIED decision | DENIED | Succeeded |
| Decision production failed | N/A | Failed |
| Audited action was refused by enforcement | N/A | Denied |

### 8.3 Six-Scope AuditEvent (AC-24)

The `AuditScope` field discriminates governance scope:

| Scope Kind | Usage |
|---|---|
| Platform | Platform-level decisions (cross-org admin) |
| Organization | Organization governance decisions |
| OrganizationUnit | OU-level delegation decisions |
| Tenant | Tenant-scoped service and access decisions |
| Project | Project-scoped workload decisions |
| ServiceInstance | Instance-level lifecycle decisions |

### 8.4 Aggregation Boundary (Requirements Section 14)

The design implements the audit-recording boundary as a pure validation function:

```go
// Package: internal/decision/audit/

// RequiresDurableRecord checks the six distinguishing criteria.
// Returns true if the event must produce a durable DecisionRecord + AuditEvent.
// Returns false only when ALL six criteria for aggregation-eligible hold.
func RequiresDurableRecord(action string, outcome core.Adjudication, obligations []core.Obligation, flags EscalationFlags, profile core.DecisionProfile) bool
```

Criteria checked (all must hold for aggregation eligibility):
1. Profile classifies action as eligible for aggregation (`AuditSemantics.AggregationEligible`)
2. Outcome is ALLOWED
3. No mandatory obligations beyond standard enforcement
4. No escalation flags (denial, privilege, cross-scope, break-glass, etc.)
5. No profile-specific escalation trigger active
6. Aggregation window within profile bounds (`AuditSemantics.MaxAggregationWindowMs`)

Failure of ANY criterion requires durable record.

## 9. Error Mapping

### 9.1 New Error Codes

| Code | HTTP Status | Meaning |
|---|---|---|
| PROFILE_REJECTED | 422 | Profile is unknown, revoked, expired, or signature-invalid |
| GRAPH_LIMIT_EXCEEDED | 422 | Dependency graph exceeds declared bounds |
| OBLIGATION_UNENFORCEABLE | 422 | Unknown/unsupported mandatory obligation makes decision unenforceable |
| SCOPE_MISMATCH | 403 | Cross-scope reference violation or scope constraint failure |
| CHAIN_DEPTH_EXCEEDED | 422 | Correction/supersession/revocation chain exceeds bounded depth |
| AUTHORITY_MISMATCH | 422 | Authority level violates profile constraints (e.g., AI as AUTHORITATIVE) |
| PROJECTION_EXCEEDED | 422 | Evaluator output projection exceeds declared bounds |
| DIGEST_MISMATCH | 422 | Content digest does not match declared or computed value |

### 9.2 Reused FEATURE-0012 Error Codes

| Code | Usage in FEATURE-0013 |
|---|---|
| VALIDATION_FAILED | General field validation errors |
| RESOURCE_NOT_FOUND | Profile, strategy, or referenced record not found |
| RESOURCE_ALREADY_EXISTS | Duplicate profile registration |
| REFERENCE_INVALID | Invalid typed reference in decision or audit record |

## 10. Security and Privacy

### 10.1 Data Minimization (SEC-07, SEC-08, AC-43)

- `DecisionRecord` stores references and digests, not raw content
- `EvaluationResult.Projection` enforces `ProjectionBounds` per profile
- Secrets, credentials, provider tokens, raw prompts never in canonical records
- Redaction is mandatory before persistence for sensitive content

### 10.2 Cross-Scope Safety (SEC-14, AC-13.5)

- `ProjectionValidator` checks `sourceScope` against target record scope (via `Metadata.ScopeRef`)
- References must not disclose existence of resources in another tenant
- Safe-denial returns a generic error without existence disclosure (EC-10)

### 10.3 Authority Enforcement (AC-07, SEC-12)

- AI output is structurally prohibited from `AUTHORITATIVE` authority
- Advisory/recommendation/simulation are visibly distinguishable
- Unknown mandatory obligations force denial even on ALLOWED (EC-08)
- Fail-closed is mandatory for security-classified profiles

### 10.4 Untrusted Data (SEC-05, SEC-06, AC-42)

- All evaluator output, reason text, evidence, AI content treated as untrusted
- Schema validation enforces size, depth, character, URI, reference-kind limits
- Validation runs before acceptance into canonical records

### 10.5 Projection Consumers (SEC-09, AC-44)

Different projections serve different consumers. The `ProjectionProfile` in `DecisionProfile` declares:
- Consumer audience (customer, operator, auditor, regulator, provider, AI)
- Redaction rules per field
- Reason-aware redaction (does not falsify omission)

### 10.6 Bundle Trust (SEC-04, AC-45)

- Signed bundles use algorithm-agile signatures (Ed25519 by default, declared in `signature.algorithm`)
- Trust-root identity is explicit and verifiable
- Unknown, revoked, or expired trust roots are rejected
- Offline verification supported without network dependency
- Verification uses `crypto/ed25519` from Go standard library
- Signed input excludes the signature value field itself (section 3.6.1)

### 10.7 Validation Error Safety

Validation error messages MUST NOT expose:
- Raw projection content from `GovernedProjection.Content`
- Secret patterns detected during SEC-07 checks
- Cross-scope identifiers that triggered SCOPE_MISMATCH
- Specific content that matched prohibited categories

Error messages report: field path, applicable bound/limit, violation type, and generic description only.

## 11. Testing Strategy

### 11.1 Unit Tests

Each source file has a dedicated `_test.go` file in the same package:

| Test File | Test Focus |
|---|---|
| `internal/decision/core/types_test.go` | Vocabulary validation, envelope field constraints |
| `internal/decision/core/profile_test.go` | Profile validation, sub-type constraints, GraphLimits conditional mandatory, AC-08 completeness |
| `internal/decision/core/evaluation_test.go` | Projection bounds, outcome validation |
| `internal/decision/core/composition_test.go` | Graph validation, cycle detection, bound checks, topological sort |
| `internal/decision/core/identity_test.go` | IdempotencyKey validation, SemanticDigest struct (no computation) |
| `internal/decision/core/scope_test.go` | AuditScope validation, six-value closed set |
| `internal/decision/core/acceptance_test.go` | Atomic acceptance contract, both-or-neither invariant |
| `internal/decision/core/errors_test.go` | Error code constants, error constructors |
| `internal/decision/validation_test.go` | Envelope validation edge cases, cross-scope safety, error-message safety |
| `internal/decision/registry_test.go` | Profile and strategy registry CRUD, duplicate rejection, conditional GraphLimits enforcement |
| `internal/decision/audit/types_test.go` | Six-scope AuditEvent validation, scope consistency invariant |
| `internal/decision/audit/acceptance_test.go` | Audit obligation acceptance |
| `internal/decision/audit/boundary_test.go` | Aggregation boundary six-criteria checks |
| `internal/decision/bundle/verifier_test.go` | Bundle verification: trust root, signature, content digest, expiry |
| `internal/decision/bundle/digest_test.go` | Content digest computation and comparison |
| `internal/decision/fixtures/conformance_test.go` | All 49+ conformance fixture scenarios |

### 11.2 Conformance Fixtures (AC-50)

28 architecture-mandated scenarios plus AC-13.7 evaluator-output scenarios plus relationship scenarios. Each fixture is:
- A JSON test-case file in `internal/decision/fixtures/testdata/`
- A Go table-driven test in `internal/decision/fixtures/conformance_test.go`

Fixture format (JSON):

```json
{
  "id": "F13-FIXTURE-001",
  "name": "Simplest synchronous atomic authorization decision",
  "description": "AC-50 scenario 1: minimal authorization allowing access",
  "category": "authorization",
  "classification": "CONTRACT_NOW",
  "executionMode": "PURE_CONTRACT",
  "input": {
    "profileRef": { "name": "AuthorizationDecision", "version": "1.0.0" },
    "form": "ADJUDICATION",
    "authority": "AUTHORITATIVE"
  },
  "expected": {
    "valid": true,
    "adjudication": "ALLOWED",
    "auditOutcome": "Succeeded",
    "auditLinked": true
  }
}
```

**`executionMode` field:** Every fixture declares its execution mode:
- `PURE_CONTRACT`: Fixture is fully executable as pure-function/in-memory validation. Input data and expected output are complete. The test runner can execute and assert.
- `CONTRACT_SHAPE_ONLY`: Fixture validates the shape and contract semantics of a scenario whose full execution requires runtime components not available in FEATURE-0013. The fixture specifies exact input structures and expected validation outcomes for the parts that ARE pure contract logic. It explicitly documents which assertions require runtime and are therefore non-executing placeholder assertions.

### 11.3 Fixture Scenario Map

| # | Scenario | Fixture ID | Key ACs | Execution Mode |
|---|---|---|---|---|
| 1 | Simplest atomic authorization | F13-FIXTURE-001 | AC-01, AC-21, AC-38 | PURE_CONTRACT |
| 2 | Deterministic denial with rationale | F13-FIXTURE-002 | AC-06, AC-25 | PURE_CONTRACT |
| 3 | Pure selection (no adjudication) | F13-FIXTURE-003 | EC-01, AC-04 | PURE_CONTRACT |
| 4 | All forms (ranking, classification, etc.) | F13-FIXTURE-004..012 | AC-04 | PURE_CONTRACT |
| 5 | Human-approval async linked to Operation | F13-FIXTURE-013 | AC-27 | CONTRACT_SHAPE_ONLY |
| 6 | Maximally bounded composite | F13-FIXTURE-014 | AC-15, AC-16 | PURE_CONTRACT |
| 7 | Parallel evaluations with aggregation | F13-FIXTURE-015 | AC-14 | PURE_CONTRACT |
| 8 | Timeout/failure/conflict | F13-FIXTURE-016 | EC-02 | PURE_CONTRACT |
| 9 | Idempotent replay | F13-FIXTURE-017 | AC-30, EC-04 | CONTRACT_SHAPE_ONLY |
| 10 | Partial persistence failure | F13-FIXTURE-018 | EC-05, AC-22 | CONTRACT_SHAPE_ONLY |
| 11 | Supersession/correction/revocation | F13-FIXTURE-019 | AC-18, AC-19 | PURE_CONTRACT |
| 12 | Unauthorized rejection | F13-FIXTURE-020 | AC-41, SEC-01 | PURE_CONTRACT |
| 13 | Offline/air-gapped | F13-FIXTURE-021 | AC-33, AC-31 | PURE_CONTRACT |
| 14 | Export/import | F13-FIXTURE-022 | AC-35 | CONTRACT_SHAPE_ONLY |
| 15 | AI-disabled | F13-FIXTURE-023 | AC-49, EC-09 | PURE_CONTRACT |
| 16 | Graph limit rejection | F13-FIXTURE-024 | EC-07 | PURE_CONTRACT |
| 17 | New family registration | F13-FIXTURE-025 | AC-09 | PURE_CONTRACT |
| 18 | Local atomic acceptance | F13-FIXTURE-026 | AC-22, AC-23 | PURE_CONTRACT |
| 19 | Idempotency vs semantic identity | F13-FIXTURE-027 | AC-30, EC-04 | CONTRACT_SHAPE_ONLY |
| 20 | Deterministic replay | F13-FIXTURE-028 | AC-13.6 | PURE_CONTRACT |
| 21 | Local-bundle operation | F13-FIXTURE-029 | AC-31 | PURE_CONTRACT |
| 22 | Profile rejection | F13-FIXTURE-030 | EC-03 | PURE_CONTRACT |
| 23 | Obligation denial | F13-FIXTURE-031 | EC-08 | PURE_CONTRACT |
| 24 | Authorization cache contract | F13-FIXTURE-032 | AC-50 note | CONTRACT_SHAPE_ONLY |
| 25 | Composite budgets | F13-FIXTURE-033 | AC-15 | PURE_CONTRACT |
| 26 | Local-digest + batch signing | F13-FIXTURE-034 | SEC-11 | CONTRACT_SHAPE_ONLY |
| 27 | Disconnected import | F13-FIXTURE-035 | AC-35 | CONTRACT_SHAPE_ONLY |
| 28 | Chain depth exceeded | F13-FIXTURE-036 | EC-06 | PURE_CONTRACT |
| E1 | Conforming projection accepted | F13-FIXTURE-037 | AC-13.7a | PURE_CONTRACT |
| E2 | Projection size exceeded rejected | F13-FIXTURE-038 | AC-13.7b | PURE_CONTRACT |
| E3 | Secret pattern rejected | F13-FIXTURE-039 | AC-13.7c | PURE_CONTRACT |
| E4 | Raw prompt rejected | F13-FIXTURE-040 | AC-13.7d | PURE_CONTRACT |
| E5 | Cross-scope content rejected | F13-FIXTURE-041 | AC-13.7e | PURE_CONTRACT |
| E6 | Valid reference+digest accepted | F13-FIXTURE-042 | AC-13.7f | PURE_CONTRACT |
| E7 | Mismatched digest rejected | F13-FIXTURE-043 | AC-13.7g | PURE_CONTRACT |
| E8 | Missing sensitivity class rejected | F13-FIXTURE-044 | AC-13.7h | PURE_CONTRACT |
| R1 | Correction chain within depth | F13-FIXTURE-045 | AC-18, AC-19 | PURE_CONTRACT |
| R2 | Supersession chain within depth | F13-FIXTURE-046 | AC-18, AC-19 | PURE_CONTRACT |
| R3 | Revocation of effective record | F13-FIXTURE-047 | AC-18, AC-19 | PURE_CONTRACT |
| R4 | Chain cycle detection | F13-FIXTURE-048 | AC-19 | PURE_CONTRACT |
| R5 | Fork resolution (multiple successors) | F13-FIXTURE-049 | AC-19 | PURE_CONTRACT |

### 11.4 CONTRACT_SHAPE_ONLY Fixture Specification

Fixtures marked `CONTRACT_SHAPE_ONLY` cannot execute complete end-to-end validation as pure contract logic because they involve:
- Asynchronous operation lifecycle (F13-FIXTURE-013)
- Persistence layer semantics (F13-FIXTURE-018)
- Idempotency with persistence state (F13-FIXTURE-017, F13-FIXTURE-027)
- Cross-system synchronization (F13-FIXTURE-022, F13-FIXTURE-035)
- Runtime authorization cache state (F13-FIXTURE-032)
- Runtime cryptographic service interaction (F13-FIXTURE-034)

For each `CONTRACT_SHAPE_ONLY` fixture, the test:
1. Validates that the input data structure is schema-conforming (field types, required fields, vocabulary values)
2. Validates that the expected output structure is schema-conforming
3. Asserts specific validation outcomes for the subset of logic that IS pure contract (e.g., profile reference validation, vocabulary checks, relationship constraints)
4. Documents as comments which assertions require runtime components and WHY they cannot execute in-memory

Example for F13-FIXTURE-032 (authorization cache contract):
- **Executing:** Profile reference valid, form=ADJUDICATION valid, authority valid, scope matches profile allowedScopes, validity fields structurally valid
- **Non-executing (documented):** Cache TTL enforcement, revocation propagation timing, stale-context detection, freshness enforcement — these require a cache runtime

### 11.5 Six-Scope AuditEvent Compatibility Fixtures (AC-54)

Six fixtures proving AuditEvent works at each governance scope:

| Fixture | Scope | Purpose |
|---|---|---|
| F13-AUDIT-SCOPE-PLATFORM | Platform | Platform-admin decision audit |
| F13-AUDIT-SCOPE-ORG | Organization | Org governance audit (existing alpha baseline) |
| F13-AUDIT-SCOPE-OU | OrganizationUnit | OU delegation audit |
| F13-AUDIT-SCOPE-TENANT | Tenant | Tenant access audit |
| F13-AUDIT-SCOPE-PROJECT | Project | Project workload audit |
| F13-AUDIT-SCOPE-SI | ServiceInstance | Instance lifecycle audit |

Each fixture validates:
- `Record.Scope.Kind` matches the expected scope
- `Metadata.ScopeRef` satisfies the consistency invariant (section 3.5.3)
- All other AuditEvent fields are structurally valid
- The fixture produces no FEATURE-0012 conformance regression

### 11.6 AC-13.7f/g Fixture Specification (Protected-Content Integrity Digest)

**F13-FIXTURE-042 (AC-13.7f — valid reference+digest accepted):**
```json
{
  "id": "F13-FIXTURE-042",
  "name": "Protected-content reference with valid integrity digest accepted",
  "executionMode": "PURE_CONTRACT",
  "input": {
    "evaluationResult": {
      "inputRef": { "apiVersion": "evidence.sovrunn.io/v1alpha1", "kind": "PolicySnapshot", "name": "tenant-a-policy-2026", "uid": "uid-evidence-001" },
      "integrityDigest": "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
      "sensitivityClass": "INTERNAL"
    },
    "contentBytes": ""
  },
  "expected": {
    "valid": true,
    "digestFormatValid": true,
    "digestContentMatch": true
  },
  "notes": "contentBytes is the empty string; SHA-256 of empty input matches the declared digest. This proves format + content verification work together."
}
```

**F13-FIXTURE-043 (AC-13.7g — mismatched digest rejected):**
```json
{
  "id": "F13-FIXTURE-043",
  "name": "Protected-content reference with mismatched integrity digest rejected",
  "executionMode": "PURE_CONTRACT",
  "input": {
    "evaluationResult": {
      "inputRef": { "apiVersion": "evidence.sovrunn.io/v1alpha1", "kind": "PolicySnapshot", "name": "tenant-a-policy-2026", "uid": "uid-evidence-002" },
      "integrityDigest": "sha256:0000000000000000000000000000000000000000000000000000000000000000",
      "sensitivityClass": "INTERNAL"
    },
    "contentBytes": "actual policy content bytes that do not match the declared digest"
  },
  "expected": {
    "valid": false,
    "digestFormatValid": true,
    "digestContentMatch": false,
    "errorCode": "DIGEST_MISMATCH",
    "errorFieldPath": "integrityDigest"
  }
}
```

**Verification logic (executed by `IntegrityVerifier`):**
1. Parse `integrityDigest` — must match `^sha256:[0-9a-f]{64}$` → if not, `DIGEST_MISMATCH`
2. Compute `"sha256:" + hex(sha256(contentBytes))`
3. Compare computed digest to declared `integrityDigest` (string equality)
4. Mismatch → `DIGEST_MISMATCH` on field path `integrityDigest`

## 12. Verification

### 12.1 Build Verification

```bash
make fmt
make test
make vet
go test -race ./internal/decision/...
```

### 12.2 Conformance Verification

All conformance fixtures must pass:
```bash
go test -v ./internal/decision/fixtures/ -run TestConformance
```

### 12.3 Schema Verification

JSON Schema files validate against published schemas:
```bash
go test -v ./internal/decision/core/ -run TestSchemaBinding
```

TypeBinding checks ensure Go structs match JSON Schema definitions (reusing FEATURE-0012 `apischema.VerifyGoTypeAgainstSchema` pattern).

### 12.4 Compatibility Verification

- Six-scope AuditEvent fixtures pass without breaking FEATURE-0012 baseline
- DecisionRecord envelope TypeBinding against JSON Schema
- Profile extension does not modify common envelope fields

## 13. Files to Create

### 13.1 Go Source Files

| File | Purpose | Package | Key ACs |
|---|---|---|---|
| `internal/decision/core/types.go` | DecisionRecord, forms, authority, adjudication, finality, correlation, relationship, rationale, obligation, provenance, validity | `core` | AC-01..07, AC-18, AC-19 |
| `internal/decision/core/profile.go` | DecisionProfile, ProfileRef, lifecycle, all sub-types | `core` | AC-08..10 |
| `internal/decision/core/evaluation.go` | EvaluationResult, EvaluationOutcome, TimingEnvelope, GovernedProjection | `core` | AC-11..13 |
| `internal/decision/core/composition.go` | CompositionStrategy, StrategyRef, FailPosture, DependencyGraph, GraphNode, GraphEdge, GraphNodeType | `core` | AC-14..16 |
| `internal/decision/core/identity.go` | IdempotencyKey, SemanticDigest (declaration only) | `core` | AC-30 |
| `internal/decision/core/scope.go` | AuditScope, AuditScopeKind constants, IsValidAuditScopeKind | `core` | AC-24 |
| `internal/decision/core/acceptance.go` | Acceptor interface, AcceptanceResult, AuditObligation | `core` | AC-22, AC-23 |
| `internal/decision/core/errors.go` | Error codes, error constructors | `core` | DQ-11 |
| `internal/decision/validation.go` | ValidateDecisionRecord, ValidateEvaluationResult, ValidateGraph, ValidateProfile, ValidateProjection, ValidateRelationshipChain, ValidateCrossScopeSafety, ValidateIntegrityDigest | `decision` | AC-01..19, AC-13.7 |
| `internal/decision/registry.go` | inMemoryProfileRegistry, inMemoryStrategyRegistry, ProfileRegistry/StrategyRegistry interfaces | `decision` | AC-10, AC-53 |
| `internal/decision/audit/types.go` | AuditEvent, AuditEventRecord, AuditOutcome | `audit` | AC-24, AC-25 |
| `internal/decision/audit/acceptance.go` | Audit-obligation acceptance contract | `audit` | AC-22 |
| `internal/decision/audit/validation.go` | ValidateAuditEvent, scope consistency validation | `audit` | AC-24, 6.7 |
| `internal/decision/audit/boundary.go` | RequiresDurableRecord, EscalationFlags | `audit` | AC-40, Req §14 |
| `internal/decision/bundle/types.go` | ProfileBundle, BundleManifest, BundleTrustRoot, BundleEntry, BundleSignature | `bundle` | AC-31 |
| `internal/decision/bundle/verifier.go` | BundleVerifier, inMemoryBundleVerifier, TrustRootStore, ContentProvider | `bundle` | AC-31, AC-53 |
| `internal/decision/bundle/digest.go` | ComputeContentDigest (SHA-256), CompareContentDigest | `bundle` | SEC-10 |

### 13.2 Test Files

| File | Package | Purpose |
|---|---|---|
| `internal/decision/core/types_test.go` | `core` | Vocabulary and envelope validation tests |
| `internal/decision/core/profile_test.go` | `core` | Profile validation, GraphLimits conditional mandatory |
| `internal/decision/core/evaluation_test.go` | `core` | Projection bounds, outcome tests |
| `internal/decision/core/composition_test.go` | `core` | Graph validation, cycle detection, topological sort |
| `internal/decision/core/identity_test.go` | `core` | IdempotencyKey validation, struct completeness |
| `internal/decision/core/scope_test.go` | `core` | AuditScope validation, closed-set enforcement |
| `internal/decision/core/acceptance_test.go` | `core` | Atomic acceptance contract tests |
| `internal/decision/validation_test.go` | `decision` | Envelope validation edge cases, cross-scope safety, integrity digest, error-message safety |
| `internal/decision/registry_test.go` | `decision` | Profile and strategy registry tests |
| `internal/decision/audit/types_test.go` | `audit` | Six-scope AuditEvent tests, scope consistency |
| `internal/decision/audit/acceptance_test.go` | `audit` | Audit acceptance tests |
| `internal/decision/audit/boundary_test.go` | `audit` | Aggregation boundary tests |
| `internal/decision/bundle/verifier_test.go` | `bundle` | Bundle verification: trust root, signature, content digest, expiry |
| `internal/decision/bundle/digest_test.go` | `bundle` | Content digest computation tests |
| `internal/decision/fixtures/conformance_test.go` | `fixtures` | All 49+ conformance fixture tests |

### 13.3 JSON Schema Files

| File | Purpose |
|---|---|
| `api/schemas/decision-record.json` | DecisionRecord JSON Schema |
| `api/schemas/decision-profile.json` | DecisionProfile JSON Schema |
| `api/schemas/evaluation-result.json` | EvaluationResult JSON Schema |
| `api/schemas/_common/decision-form.json` | Closed form vocabulary |
| `api/schemas/_common/decision-authority.json` | Closed authority vocabulary |
| `api/schemas/_common/adjudication.json` | Closed adjudication vocabulary |
| `api/schemas/_common/audit-outcome.json` | Closed audit outcome vocabulary |
| `api/schemas/_common/audit-scope.json` | Six-scope audit scope |

### 13.4 Fixture Data Files

All fixtures use JSON format (standard library parseable):

| File | Purpose |
|---|---|
| `internal/decision/fixtures/testdata/f13-fixture-001-simple-auth.json` | Scenario 1 |
| `internal/decision/fixtures/testdata/f13-fixture-002-denial.json` | Scenario 2 |
| `internal/decision/fixtures/testdata/f13-fixture-003-pure-selection.json` | Scenario 3 |
| `internal/decision/fixtures/testdata/...` | Scenarios 4–28 |
| `internal/decision/fixtures/testdata/f13-fixture-037-projection-valid.json` | AC-13.7a |
| `internal/decision/fixtures/testdata/...` | AC-13.7b–h |
| `internal/decision/fixtures/testdata/f13-fixture-042-digest-valid.json` | AC-13.7f |
| `internal/decision/fixtures/testdata/f13-fixture-043-digest-mismatch.json` | AC-13.7g |
| `internal/decision/fixtures/testdata/f13-fixture-045-correction-chain.json` | Relationship: correction |
| `internal/decision/fixtures/testdata/f13-fixture-046-supersession-chain.json` | Relationship: supersession |
| `internal/decision/fixtures/testdata/f13-fixture-047-revocation.json` | Relationship: revocation |
| `internal/decision/fixtures/testdata/f13-fixture-048-chain-cycle.json` | Chain cycle detection |
| `internal/decision/fixtures/testdata/f13-fixture-049-fork-resolution.json` | Fork resolution |
| `internal/decision/fixtures/testdata/f13-audit-scope-platform.json` | Platform scope |
| `internal/decision/fixtures/testdata/f13-audit-scope-org.json` | Organization scope |
| `internal/decision/fixtures/testdata/f13-audit-scope-ou.json` | OU scope |
| `internal/decision/fixtures/testdata/f13-audit-scope-tenant.json` | Tenant scope |
| `internal/decision/fixtures/testdata/f13-audit-scope-project.json` | Project scope |
| `internal/decision/fixtures/testdata/f13-audit-scope-si.json` | ServiceInstance scope |

### 13.5 Documentation Files

| File | Purpose |
|---|---|
| `docs/features/FEATURE-0013-decision-record-and-auditevent-standard.md` | Feature documentation |

## 14. Downstream Adoption Contract (AC-55, AC-56, AC-57, AC-58)

This section implements the exact lightweight downstream-adoption structure from canonical architecture section 28. No manifests, registry, index, or separate approval workflow are introduced.

### 14.1 Applicability Decision

A feature must include the adoption section when it directly depends on FEATURE-0013 or produces, consumes, specializes, projects, audits, or enforces any of the following: `EvaluationResult`, `DecisionProfile`, `DecisionRecord`, decision rationale or obligations, decision-linked `AuditEvent`, decision-linked `Operation`, `DecisionContext`, or a decision composition strategy.

Applicability uses controlled values:
- `APPLICABLE` — the feature performs one or more of the above behaviors
- `NOT_APPLICABLE` — the feature performs none; requires concise rationale

### 14.2 Adoption Roles (Controlled Vocabulary)

- `EVALUATION_PRODUCER`
- `DECISION_PRODUCER`
- `DECISION_CONSUMER`
- `AUDIT_PRODUCER`
- `PROJECTION_CONSUMER`
- `OPERATION_OWNER`
- `COMPOSITE_ADOPTER`

### 14.3 Required Per-Feature Section (Applicable)

```markdown
## FEATURE-0013 Adoption

Applicability: APPLICABLE

Contract version: FEATURE-0013 v1alpha1 (decision.sovrunn.io/v1alpha1)

Roles:
- <one or more from controlled vocabulary>

Profiles introduced or reused:
- <profile name and version>

Contracts reused:
- <FEATURE-0013 objects/vocabularies reused>

Inherited risks (selective, with dispositions):
- F13-Rxx: <APPLICABLE | NOT_APPLICABLE | MITIGATED_BY_FEATURE | TRANSFERRED_FORWARD> — <rationale>

Exceptions:
- None | <list with rationale>

Conformance evidence:
- <fixtures/checks demonstrating compliance>
```

### 14.4 Required Per-Feature Section (Not Applicable)

```markdown
## FEATURE-0013 Adoption

Applicability: NOT_APPLICABLE

Rationale: <why the feature produces or consumes none of the governed concepts>
```

### 14.5 Feature-Gate Check (AC-57)

A lightweight gate validates:
1. Adoption section exists when applicability rules require it
2. Applicability value is `APPLICABLE` or `NOT_APPLICABLE`
3. Roles use only the controlled vocabulary
4. Contract version references exact FEATURE-0013 contract version
5. Profiles are named and versioned
6. Reused contracts are identified
7. Inherited risks reference valid `F13-Rxx` identifiers with controlled disposition
8. Conformance evidence is identified
9. Common envelope and closed vocabularies are not redefined

### 14.6 Prohibited Artifacts (AC-58)

The adoption contract explicitly DOES NOT introduce:
- Separate per-feature YAML manifests
- A central adoption registry
- A generated adoption index
- A separate adoption approval workflow
- A runtime adoption service

### 14.7 Selective Risk Inheritance

Controlled disposition values:
- `APPLICABLE` — risk applies and controls are in scope
- `NOT_APPLICABLE` — risk does not apply; rationale required
- `MITIGATED_BY_FEATURE` — feature implements the risk's control
- `TRANSFERRED_FORWARD` — risk ownership transferred to a later feature

### 14.8 Unified Approval

Adoption is reviewed within the downstream feature's existing architecture and semantic review. No separate adoption approval workflow exists.

## 15. Non-Goals (Explicit)

This design intentionally does NOT:

- NG-D01: Implement a production decision service, authorization service, or policy engine
- NG-D02: Implement database persistence, migrations, transactions, or outbox patterns
- NG-D03: Implement a queue, broker, workflow engine, scheduler, or worker pool
- NG-D04: Implement network registry clients, remote resolvers, HSM clients, or AI model clients
- NG-D05: Implement a rules engine, policy language, workflow DSL, or expression interpreter
- NG-D06: Select any specific workflow engine, broker, database, identity provider, or AI model
- NG-D07: Implement runtime authorization caching, token validation, or revocation distribution
- NG-D08: Implement HTTP API handlers or REST endpoints for decision or audit resources
- NG-D09: Expose decision or audit resources through the Phase 1 API server
- NG-D10: Import `internal/server` or `internal/api` packages
- NG-D11: Add external Go dependencies beyond standard library
- NG-D12: Invent numeric defaults for GraphLimits, retention periods, SLOs, or operational constants
- NG-D13: Implement real audit event persistence or materialization
- NG-D14: Implement provider-specific adapters or evaluator integrations
- NG-D15: Implement RFC 8785 JCS or any semantic decision digest computation function (pending ADR-F13-003)
- NG-D16: Publish a partial `SemanticDigestCoveredFields` list (deferred entirely to ADR-F13-003)

## 16. Resolved Design Questions

| DQ | Resolution | Design Section |
|---|---|---|
| DQ-01 | `json.RawMessage` for typedResult | 4.1 |
| DQ-02 | Go struct with embedded validation | 4.5 |
| DQ-03 | Adjacency list with typed nodes; profile-declared conditional bounds | 4.6, 4.5.2 |
| DQ-04 | RFC 8785 JCS candidate DEFERRED to ADR-F13-003; struct declared, no computation, no covered-field list | 4.11, 3.6.3 |
| DQ-05 | `Record.Scope` is authoritative audit scope; `Metadata.ScopeRef` for FEATURE-0012 compat; consistency validated | 3.5 |
| DQ-06 | Static Go registry map | 7.2 |
| DQ-07 | Pure function returning AcceptanceResult | 4.10, 5.1 |
| DQ-08 | Go table-driven tests with JSON data files (standard library parseable) | 11.2 |
| DQ-09 | Embedding + `$ref` composition | 3.3 |
| DQ-10 | Signed JSON manifest with SHA-256 content digests (standard library) | 7.3, 3.6.1 |
| DQ-11 | Extend FEATURE-0012 pattern with 8 new codes | 9.1 |
| DQ-12 | Four-package layout: `core/`, `audit/`, `bundle/`, `fixtures/` | 3.1 |
| DQ-13 | Bounded `json.RawMessage` with `ProjectionBounds` per profile | 4.4 |
| DQ-14 | Separate validation pass with `SourceScope` propagation | 5.4, 6.6 |

## 17. ARCHITECTURE_DECISION_REQUIRED Escalations

| ID | Subject | Context | Escalation Owner | Resolution Path | Blocking Scope |
|---|---|---|---|---|---|
| ADR-F13-001 | System-wide maximum ceiling for GraphLimits | DQ-03: if an upper bound that no profile may exceed is needed, architecture must decide it | Architecture Owner | ADH or Architecture Owner directive | Graph validation fixtures that test absolute ceiling (if any) |
| ADR-F13-002 | Default maximum relationship chain depth | AC-19, EC-06: whether a system-wide default exists beyond per-profile MaxChainDepth | Architecture Owner | ADH or Architecture Owner directive | Chain-depth fixtures use profile-declared value; system ceiling deferred |
| ADR-F13-003 | Canonical serialization exact algorithm, version, covered-field set, field ordering, and byte representation for semantic decision digests | DQ-04: RFC 8785 JCS is the candidate. Architecture section 9 requires coverage of scope, purpose, profile name+version, authority, canonical input snapshot, EffectivePolicyContext, composition strategy/version, evaluator contract versions, and evaluation-time/validity epoch. Exact specification requires ADR approval. | Architecture Owner | DEC record | Semantic digest computation, digest-based deduplication, SemanticDigestCoveredFields publication. Note: bundle-entry content digests and protected-content integrity digests are NOT blocked by this ADR. |
| ADR-F13-004 | ServiceInstance as ScopeKind in FEATURE-0012 | AC-24 requires ServiceInstance scope; FEATURE-0012 has Provider not ServiceInstance. Whether to add ServiceInstance to apimeta.ScopeKind or keep as FEATURE-0013 local constant. Temporary contract (section 3.5.2) is compliant without this resolution. | Architecture Owner | ADH with FEATURE-0012 owner | AuditScope permanent type relationship to apimeta.ScopeKind; removal of temporary AuditScope string-Kind workaround |
| ADR-F13-005 | Exact actor-assurance-level vocabulary | AC-08 actor semantics include assurance; exact level values not approved | Architecture Owner | DEC record | Actor-assurance validation fixtures |
| ADR-F13-006 | Exact jurisdiction vocabulary and residency-enforcement semantics | AC-08 residency; NG-09 prohibits inventing jurisdiction-specific values | Architecture Owner | DEC record | Jurisdiction-specific validation fixtures |
| ADR-F13-007 | Exact retention-class vocabulary and minimum/maximum periods | AC-08 retention/legal-hold; NG-09 prohibits inventing exact schedules | Architecture Owner | DEC record | Retention-class validation fixtures with real period values |

Conformance fixtures use clearly labelled placeholder values (e.g., `PLACEHOLDER_PENDING_ADR_F13_003`) that fail validation until the architecture decision supplies the approved value.

## 18. Observability and Audit Applicability Statement

This feature produces contract-only Go code. The following observability constraints apply:

1. **No production logging** is introduced.
2. **No metrics** are introduced.
3. **No traces** are introduced.
4. **No lifecycle paths** are introduced.
5. **Correlation fields are schema-only.** `Correlation.TraceID`, `Correlation.SpanID`, and `Correlation.RequestID` are data-model fields for future runtime use. They carry no runtime behavior in this feature.
6. **AuditEvent is an obligation/contract, not a log substitute.**
7. **Tests must verify that sensitive projection data is not exposed in validation errors.**

## 19. Reuse Assessment

### 19.1 Feature-Level Reuse Summary

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| Common DecisionRecord envelope and vocabulary | Build | Sovrunn-owned differentiation; no existing standard provides the exact governed-conclusion envelope | Approved | ADH-2026-014, DEC-0026 |
| DecisionProfile extension mechanism | Build | Sovrunn-owned typed extensibility for decision families | Approved | ADH-2026-014 |
| EvaluationResult normalization | Extend | Extends FEATURE-0012 grammar with evaluation-specific fields | Approved | ADH-2026-014, DEC-0036 |
| Composition and bounded dependency graph | Reuse | Reuses bounded decision-requirements concepts from DMN | Approved | ADH-2026-014 |
| AuditEvent linkage contract | Extend | Extends FEATURE-0012 alpha AuditEvent to six scopes | Approved | ADH-2026-014 |
| Canonical serialization and digest | Reuse | Reuses RFC 8785 JCS concepts for canonicalization (deferred implementation) | Approved | ADH-2026-014 |
| Telemetry correlation | Reuse | Reuses OpenTelemetry correlation concepts for trace/span (schema-only) | Approved | ADH-2026-014 |
| Event transport mapping | Reuse | Optional CloudEvents transport mapping | Approved | ADH-2026-014 |
| Supply-chain evidence references | Reuse | Reuses DSSE/in-toto and SPDX/CycloneDX concepts | Approved | ADH-2026-014 |

### 19.2 Assessment Format Version

```text
reuse_assessment_format_version: 1.0.0
```

## 20. Architecture Drift Checks

| Drift check | Status | Evidence |
|---|---|---|
| No provider-specific hardcoding in core | PASS | No provider identifiers in any schema, type, or vocabulary |
| No Kubernetes-only assumptions in core | PASS | Contracts are HTTP/API-native per FEATURE-0012 grammar; no CRD dependency |
| No PostgreSQL lifecycle logic in core placement engine | PASS | Not applicable; FEATURE-0013 does not contain placement or PostgreSQL logic |
| No custom policy engine embedded in handlers | PASS | Hierarchy resolved through EffectivePolicyContext reference; policy engines via adapters per DEC-0028 |
| No raw secret storage | PASS | Records use references and digests per SEC-07; ProjectionBounds.ProhibitedCategories enforces no secrets |
| No customer-facing IaaS leakage | PASS | Provider details enter only through adapters and typed references per AC-36 |
| Explainable decision object | PASS | Structured rationale, reason codes, obligations, projection profiles mandatory per AC-01, AC-08, AC-46 |
| Defined audit behavior | PASS | Atomic decision/audit acceptance, six scopes, separated outcomes per AC-21..AC-25 |
| Preserved adapter boundaries | PASS | Evaluator-native payloads behind adapters per AC-12; no external engine types in core; DEC-0036 referenced |
| No external dependency beyond standard library | PASS | FEATURE-0013 packages use only Go standard library |
| No YAML parsing in FEATURE-0013 code | PASS | Fixtures use JSON; bundles use JSON; profile data uses JSON |
| No circular package imports | PASS | Four-package layout with strict top-down import rules; `core/` is leaf; no sibling imports |

## 21. Risk Traceability Matrix (F13-R01 through F13-R31)

This section maps every architecture risk to design components, treatment classification, verification evidence, staging owner, and pending disposition.

| Risk ID | Category | Design Component(s) | Treatment | CONTRACT_NOW Evidence | Verification Fixture(s) | Owner |
|---|---|---|---|---|---|---|
| F13-R01 | Governance/compat | ProfileRegistry, profile validation (§7.1, §6.4) | MITIGATE | Profile uniqueness validation, duplicate-rejection test | F13-FIXTURE-025, registry_test.go | Architecture owner |
| F13-R02 | Correctness | DecisionResult.TypedResult (§4.1), ResultSchema in profile (§4.5) | MITIGATE | TypedResult size-bounded validation, schema-binding test | types_test.go, conformance_test.go | FEATURE-0013 owner |
| F13-R03 | Security/correctness | Authority vocabulary (§4.2), validation rule authority+AI (§6.1) | MITIGATE | Authority-confusion negative tests, AUTHORITY_MISMATCH error | F13-FIXTURE-020, validation_test.go | Security architecture |
| F13-R04 | Correctness | PolicySemantics.MaxStalenessMs (§4.5.1), EffectivePolicyContextRef (§4.1) | MITIGATE | Staleness field presence validation | profile_test.go | Effective-context owner |
| F13-R05 | Correctness | IdempotencyKey (§4.11), SemanticDigest struct declaration | MITIGATE | IdempotencyKey struct validation, separate identity tests | F13-FIXTURE-017, identity_test.go | Decision persistence owner |
| F13-R06 | Audit/compliance | AcceptanceResult both-or-neither (§4.10) | MITIGATE | Atomic acceptance fixture: both nil or both present | F13-FIXTURE-018, F13-FIXTURE-026, acceptance_test.go | Audit persistence owner |
| F13-R07 | Availability/latency | AuditObligation local acceptance (§4.10), no remote audit import | AVOID | No remote-audit import in package; acceptance is local pure function | acceptance_test.go | Audit implementation owner |
| F13-R08 | Cost/availability | GraphLimits (§4.5.2), graph validation (§6.3) | MITIGATE | Graph bound enforcement tests, limit-exceeded rejection | F13-FIXTURE-024, F13-FIXTURE-033, composition_test.go | Decision runtime owner |
| F13-R09 | Correctness/AI | GovernedProjection capture (§4.4), DeterminismPolicy (§4.5.1) | MITIGATE | Replay fixture uses captured projection without re-execution | F13-FIXTURE-028, evaluation_test.go | Evaluator/AI owner |
| F13-R10 | Security/AI | GovernedProjection bounds (§4.4), cross-scope safety (§6.6), ProjectionBounds (§4.5.2) | MITIGATE | Cross-scope/redaction/prompt-leak negative fixtures | F13-FIXTURE-039..041, validation_test.go | AI-context/security owners |
| F13-R11 | Security | Obligation validation (§6.1), OBLIGATION_UNENFORCEABLE error (§9.1) | AVOID | Unknown-obligation denial fixture | F13-FIXTURE-031, validation_test.go | Enforcement owner |
| F13-R12 | Latency/availability | No remote import in decision packages; bundle local verification (§7.3) | MITIGATE | Import audit: no net/* imports | bundle_test.go, drift check §20 | Authorization impl owner |
| F13-R13 | Security | Validity (§4.7), ValidityRules (§4.5.1), FailPosture (§4.5.1) | MITIGATE | Expiry/revocation field validation, fail-closed profile tests | profile_test.go | Authorization/security owner |
| F13-R14 | Security/availability | FailPosture FAIL_CLOSED enforcement for security profiles (§6.4) | AVOID | Profile validation rejects FAIL_OPEN for security-classified | profile_test.go, validation_test.go | Security architecture owner |
| F13-R15 | Correctness | DecisionRelationship (§4.8), chain validation (§6.5), EffectiveState (§4.8) | MITIGATE | Chain/cycle/fork fixtures, deterministic projection | F13-FIXTURE-019, F13-FIXTURE-045..049 | Read-model owner |
| F13-R16 | Privacy/security | ProhibitedCategories (§4.5.2), validation error safety (§6.8, §10.7) | MITIGATE | Secret/PII/cross-scope negative tests; error-message-safety tests | F13-FIXTURE-039..041, validation_test.go | Security/privacy owner |
| F13-R17 | Compliance/integrity | Validity.Revocable (§4.7), ValidityRules (§4.5.1) | MITIGATE | Chain-integrity fixtures, tombstone contract (schema only) | F13-FIXTURE-019, types_test.go | Evidence/retention owner |
| F13-R18 | Latency/availability | Bundle local verification (§7.3), no HSM import | AVOID | Bundle verification uses crypto/ed25519 locally; no remote call | bundle_test.go | Security infrastructure owner |
| F13-R19 | Sovereignty/security | Bundle trust-root verification (§7.3), unknown/revoked root rejection | MITIGATE | Forgery/unknown-root/revoked-root fixtures | F13-FIXTURE-021, bundle_test.go | Synchronization/security owner |
| F13-R20 | Correctness | No last-writer-wins in relationship model (§4.8); fork resolution by timestamp | MITIGATE | Fork-resolution fixture | F13-FIXTURE-049 | Multi-site owner |
| F13-R21 | Portability | No provider-native identifiers in types; adapter boundary (§3.3) | AVOID | Import audit, schema lint, portability drift check | Drift check §20 | Adapter architecture owner |
| F13-R22 | Scope/governance | Non-goals (§15), package imports (§3.2), contract-only scope | AVOID | Changed-file gate, dependency gate, no-side-effect check | Drift check §20, NG-D01..NG-D16 | FEATURE-0013 owner |
| F13-R23 | Cost | ProjectionBounds (§4.5.2), GraphLimits (§4.5.2) | MITIGATE | Projection-size and graph-payload limit tests | F13-FIXTURE-038, composition_test.go | Audit/evidence owner |
| F13-R24 | Supply chain/security | BundleVerifier (§5.5, §7.3), trust-root identity, signature verification | MITIGATE | Signature/tamper/revocation fixtures | F13-FIXTURE-029, bundle_test.go | Supply-chain security owner |
| F13-R25 | Correctness/security | TimingEnvelope (§4.4), Validity timestamps (§4.7) | MITIGATE | Timestamp format validation, expiry-boundary tests | evaluation_test.go, types_test.go | Runtime/security owner |
| F13-R26 | Security/privacy | Cross-scope safety validation (§6.6), SCOPE_MISMATCH error, error-message safety (§10.7) | MITIGATE | Cross-scope and safe-denial fixtures | F13-FIXTURE-041, validation_test.go | API/security owner |
| F13-R27 | Compatibility/governance | Stable envelope (§4.1), profile extension without envelope change (§4.5) | AVOID | Schema-binding test against JSON Schema | types_test.go, schema verification §12.3 | Architecture owner |
| F13-R28 | Compatibility | Six-scope AuditEvent (§4.9), compatibility fixtures (§11.5) | MITIGATE | Six-scope compatibility fixtures, FEATURE-0012 baseline tests | F13-AUDIT-SCOPE-*, audit/types_test.go | FEATURE-0013/API owner |
| F13-R29 | Traceability | DecisionRecord terminology (§4.1) | MITIGATE | Repository-wide identity check | Drift check §20 | Architecture owner |
| F13-R30 | Cost/sovereignty | Tiered conformance via profile-declared bounds (§4.5) | MITIGATE | Minimal profile fixture without remote deps | profile_test.go | Deployment architecture owner |
| F13-R31 | Governance/cost | Downstream adoption (§14): one section, one gate, no manifests | AVOID | Gate validation fixture; absence of manifests/registry/index | Drift check §20, §14.6 | Architecture governance owner |

## 22. Supporting Types

```go
// Package: internal/decision/audit/

// EscalationFlags are the flags that force a durable record (Requirements §14.3 criterion 4).
type EscalationFlags struct {
    Denial              bool `json:"denial"`
    PrivilegeEscalation bool `json:"privilegeEscalation"`
    CrossScopeAccess    bool `json:"crossScopeAccess"`
    CrossTenantAccess   bool `json:"crossTenantAccess"`
    BreakGlass          bool `json:"breakGlass"`
    DataDisclosure      bool `json:"dataDisclosure"`
    DataExport          bool `json:"dataExport"`
    PolicyChange        bool `json:"policyChange"`
    PrivilegeChange     bool `json:"privilegeChange"`
    HumanApproval       bool `json:"humanApproval"`
    RegulatedAction     bool `json:"regulatedAction"`
}
```

```go
// Package: internal/decision/bundle/

// ProfileBundle is the signed bundle manifest for offline profile distribution (AC-31).
type ProfileBundle struct {
    BundleVersion string          `json:"bundleVersion"`
    TrustRoot     BundleTrustRoot `json:"trustRoot"`
    Entries       []BundleEntry   `json:"entries"`
    Signature     BundleSignature `json:"signature"`
    ValidityMs    int64           `json:"validityMs,omitempty"` // expiry window from signedAt
}

// BundleTrustRoot identifies the signing authority.
type BundleTrustRoot struct {
    ID           string `json:"id"`
    Algorithm    string `json:"algorithm"`
    PublicKeyRef string `json:"publicKeyRef"`
}

// BundleEntry is one profile in the bundle.
type BundleEntry struct {
    Name          string `json:"name"`
    Version       string `json:"version"`
    ContentDigest string `json:"contentDigest"` // "sha256:<hex>"
    Path          string `json:"path"`
}

// BundleSignature is the algorithm-agile bundle signature.
type BundleSignature struct {
    Algorithm string `json:"algorithm"`
    Value     string `json:"value"`    // base64-encoded
    SignedAt  string `json:"signedAt"` // RFC 3339
}
```

## 23. Dependency and Import Audit Summary

### 23.1 External Dependencies Used by FEATURE-0013

| Dependency | Used? | Justification |
|---|---|---|
| `gopkg.in/yaml.v3` | NO | Exists in go.mod for other features; FEATURE-0013 uses JSON exclusively |

### 23.2 Standard Library Packages Used

| Package | Usage |
|---|---|
| `encoding/json` | All data serialization: fixtures, bundles, typed results, schemas |
| `crypto/sha256` | Bundle-entry content digest and protected-content integrity digest computation |
| `crypto/ed25519` | Bundle signature verification |
| `encoding/hex` | Digest hex encoding |
| `encoding/base64` | Signature value decoding |
| `sync` | RWMutex for in-memory registries |
| `strings`, `regexp` | Validation pattern matching (DNS names, semver, digest format) |
| `time` | Timestamp validation (RFC 3339 parsing), bundle expiry |
| `testing` | Test framework |

### 23.3 Internal Packages Imported

| Package | Imported By |
|---|---|
| `internal/apimeta` | `internal/decision/core/`, `internal/decision/audit/`, `internal/decision/bundle/` |
| `internal/apiproblem` | `internal/decision/core/`, `internal/decision/`, `internal/decision/audit/` |
| `internal/apischema` | `internal/decision/core/` (for TypeBinding conformance tests only) |
| `internal/decision/core/` | `internal/decision/`, `internal/decision/audit/`, `internal/decision/bundle/`, `internal/decision/fixtures/` |

### 23.4 Packages NOT Imported (Explicitly Prohibited)

| Package | Reason |
|---|---|
| `internal/server` | No HTTP handlers (NG-D08) |
| `internal/api` | No API registration (NG-D09) |
| `internal/registry` | Phase 1 registries are separate (NG-D10) |
| `net/*` | No network operations (NG-D04) |
| `database/*` | No persistence (NG-D02) |
| `os` | No filesystem I/O in library code (test helpers may use os for testdata) |
| `internal/decision/audit/` | NOT imported by `internal/decision/core/` or `internal/decision/` root |
| `internal/decision/bundle/` | NOT imported by `internal/decision/core/` or `internal/decision/` root |

## 24. Revision Summary: Changes from Prior Design

This section documents the material changes made to resolve the reviewer's blocking issues.

### 24.1 Package Layout Restructured (Blocking Issue 1)

**Problem:** `internal/decision/acceptance.go` used `AuditScope` from child package `internal/decision/audit/`, creating an uncompilable circular import.

**Resolution:** Introduced `internal/decision/core/` leaf package that owns all shared types including `AuditScope`, `AuditObligation`, and `AcceptanceResult`. Both the root `internal/decision/` and `internal/decision/audit/` import `core/` without cycles. Added `internal/decision/bundle/` for bundle-specific logic. Strict import direction rules prevent any circular dependency.

### 24.2 AuditEvent Scope Semantics Unified (Blocking Issue 2)

**Problem:** Contradictory dual scope grammar — `apimeta.ScopeKind` in metadata vs local string-based `AuditScope` with ServiceInstance, no consistency validation.

**Resolution:** Defined `Record.Scope` (`AuditScope`) as the authoritative audit scope for FEATURE-0013 consumers. `Metadata.ScopeRef` provides FEATURE-0012 compatibility. Mandatory consistency validation invariant (section 3.5.3) defines exact rules per scope kind including ServiceInstance's temporary treatment. ADR-F13-004 explicitly blocks permanent resolution. The design does NOT claim `AuditScope` composes `apimeta.ScopeKind`.

### 24.3 Digest Taxonomy Separated (Blocking Issue 3)

**Problem:** Bundle-entry content digests, protected-content integrity digests, and semantic decision digests were conflated. ADR-F13-003 was overapplied to block ALL digest verification.

**Resolution:** Section 3.6 defines three distinct digest contracts. Bundle-entry content digests (SHA-256, fully specified) and protected-content integrity digests (SHA-256, fully specified with format + content verification) are FULLY EXECUTABLE. Only semantic decision digests are deferred to ADR-F13-003. Each has its own algorithm, input, encoding, verification logic, and error behavior.

### 24.4 AC-13.7g Made Executable (Blocking Issue 4)

**Problem:** Design claimed to reject mismatched integrity digests (AC-13.7g) while simultaneously deferring digest matching.

**Resolution:** Protected-content integrity digest verification is fully specified (section 3.6.2). When content bytes are available (as they are in conformance fixtures), the verifier computes SHA-256 and compares. Fixture F13-FIXTURE-043 provides explicit mismatched `contentBytes` and expects `DIGEST_MISMATCH` rejection. The `IntegrityVerifier` interface (section 5.7) provides both format and content verification methods.

### 24.5 SemanticDigestCoveredFields Removed (Blocking Issue 5)

**Problem:** Prior design published a partial `SemanticDigestCoveredFields` list that omitted architecture-required inputs (purpose, EffectivePolicyContext, evaluator contract versions, validity epoch).

**Resolution:** The `SemanticDigestCoveredFields` variable is entirely removed (NG-D16). The `SemanticDigest` struct retains a `CoveredFields []string` field for future use but it is empty until ADR-F13-003 resolves. No partial or misleading list is published.

### 24.6 GraphLimits Made Conditionally Mandatory (Blocking Issue 6)

**Problem:** `GraphLimits` was declared optional despite architecture requiring finite bounds for profiles permitting composition.

**Resolution:** Section 4.5.2 defines the conditional mandatory rule: `GraphLimits` is required when `CompositionStrategies` is non-empty OR `AcceptedEvaluators` has more than one entry. It remains optional for simple single-evaluator profiles. Validation (section 6.4) enforces this rule. No system-wide numeric defaults are invented.

### 24.7 Conformance Fixtures Classified (Blocking Issue 7)

**Problem:** Several conformance scenarios (authorization cache, export/import, batch signing) implied runtime behavior that cannot be validated as pure contract logic.

**Resolution:** Every fixture now declares `executionMode`: either `PURE_CONTRACT` (fully executable in-memory) or `CONTRACT_SHAPE_ONLY` (validates structure and schema-compliant subset; documents which assertions require runtime). Section 11.4 specifies exactly what `CONTRACT_SHAPE_ONLY` fixtures execute and what they document as non-executing.

---

End of design.
