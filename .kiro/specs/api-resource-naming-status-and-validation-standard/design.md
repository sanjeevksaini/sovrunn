# Design Document

Feature: FEATURE-0012 — API, Resource Naming, Status, and Validation Standard
Phase: Phase 2 — Reuse-First PaaS Fabric Foundation (initial adoption phase; the standard applies across phases)
Stage: Design
Controlling handoff: ADH-2026-012 (Approved), ADH-2026-013 (Approved)
Canonical architecture: docs/architecture/api-resource-standard.md
Canonical reuse standard: docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md
Depends on: FEATURE-0011 (Reuse Assessment Standard)

## Overview

FEATURE-0012 delivers the Sovrunn-owned, provider-neutral API and resource
grammar as a set of shared Go primitives, machine-readable schema
conventions, strict decoding/validation helpers, conformance fixtures, a
Phase 1 compatibility report, and feature-gate checks. It is a
standard-and-conformance foundation, not a domain runtime.

The architecture horizon is broad and cross-phase; the FEATURE-0012
implementation scope is deliberately narrow (F12-IMPL-001, F12-IMPL-002).
This design converts the approved architecture baseline
(`docs/architecture/api-resource-standard.md`) and the normative
requirements (`F12-*`) into a concrete implementation design without
introducing any later-feature domain behavior, runtime, or persistence.

This design delivers exactly these implementation elements:

1. Shared, versioned type-metadata, common-metadata, scope, reference,
   condition, and problem primitives (Go packages under `internal/`).
2. Machine-readable canonical schemas (JSON Schema 2020-12 documents) with
   Sovrunn `x-sovrunn-*` annotations as the single source of truth per
   contract.
3. Strict, operation-aware decoding (pure JSON and YAML decoders separated
   from the HTTP adapter) and layered validation helpers (unknown/duplicate-
   field rejection, bounded-subset structural schema validation, deterministic
   defaulting, stable codes, JSON Pointer paths).
4. RFC 9457 Problem Details error transport with Sovrunn extensions and a
   stable HTTP status mapping.
5. Eight representative contract fixtures proving grammar fit without
   domain execution.
6. A Phase 1 compatibility report and an executable Phase 1 coverage check.
7. Executable fitness functions wired into the FEATURE-0012 feature gate.

No provider, plugin, policy, placement, audit, or provisioning service is
created (F12-IMPL-002). No new runtime HTTP routes are added; the route
grammar is delivered as documentation plus a route-form validation helper.

This design resolves the twelve deferred design questions in requirements
section 9 within the approved contracts. Where a required semantic decision
were unavailable from approved sources, this design would halt and report
`ARCHITECTURE_DECISION_REQUIRED` (F12-GOV-001); ADH-2026-013 resolved the
previously missing exact Operation scope enumeration; no other halt was
required.

## Controlling inputs

This design is derived only from the approved controlling inputs and
introduces no new architecture decision. ADH-2026-013 resolved the
previously missing exact Operation scope enumeration as a bounded
clarification.

| Input | Role in this design |
|---|---|
| requirements.md (Approved for design) | Normative acceptance criteria (`F12-*`). |
| docs/architecture/api-resource-standard.md | Approved architecture baseline; matrices A–E and fitness functions. |
| ADH-2026-012 (Approved) | Controlling handoff; approves the Extend disposition and standard scope. |
| ADH-2026-013 (Approved) | Controlling handoff; resolves the canonical Operation allowed-scope enumeration and target-scope equality invariant. |
| RFC-0022 (DEC-0026, DEC-0027, DEC-0036) | Reuse-first and adapter-boundary basis. |
| docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md | Canonical reuse assessment schema (fields not redefined here). |
| docs/engineering/go-coding-guardrails.md | Go implementation guardrails. |
| docs/engineering/go-version-standard.md | Go version standard (minimum Go 1.22). |
| docs/engineering/ai-context-loading-standard.md | Context loading discipline. |
| Existing Phase 1 packages (`internal/resources`, `internal/api`, `internal/validation`, `internal/registry`) | Existing patterns extended by the shared grammar. |

## Reuse summary (feature-level)

Field definitions are owned by
`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` and are not redefined
here. The approved feature-level summary is carried unchanged from the
requirements.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| Sovrunn API/resource meta-model and conformance foundation | Extend | Extend mature HTTP, OpenAPI 3.1/JSON Schema 2020-12, RFC 9457 Problem Details, RFC 6901 JSON Pointer, ETag/`If-Match` concurrency, and selected Kubernetes API conventions with Sovrunn-owned sovereign scope, boundary, ownership, compatibility, and conformance rules. | Approved | ADH-2026-012; RFC-0022 (DEC-0026, DEC-0027, DEC-0036) |

Sovrunn owns the resource-profile taxonomy, API/naming conventions, common
metadata and identity, scope and reference semantics, API-boundary
classification, ownership and mutability rules, status/condition grammar,
validation and error contracts, provider-neutrality constraints,
compatibility policy, conformance rules, and reassessment triggers.
HTTP semantics, OpenAPI 3.1, JSON Schema 2020-12, RFC 9457, RFC 6901,
ETag/`If-Match`, and selected Kubernetes API conventions are the reused or
extended responsibility. Adapter required: No — FEATURE-0012 defines
contract grammar and performs no external integration; later integrations
must use DEC-0036-compliant adapter boundaries. The capability-level
assessment appears in the "Reuse assessment (capability-level)" section.

## Resolved decisions

These are the binding design decisions. Each stays inside the approved
baseline and cites the requirement(s) it satisfies. None introduces a new
architecture decision; none expands FEATURE-0012 implementation scope.

| ID | Decision | Satisfies |
|---|---|---|
| D-01 | Canonical contract source of truth is a JSON Schema 2020-12 document per contract, under `api/schemas/`; Go types, docs, fixtures, and SDKs are derivative and checked for consistency. Structural validation of the canonical schemas is executable via the bounded-subset validator in D-01a, and derivative Go-type consistency is executable via the reflection-based check in D-01b — neither relies on fixture round-tripping and annotation checks alone as proof. | F12-NAMING-005; F12-VALIDATION-001(4) |
| D-01b | Derivative Go-type consistency with the canonical schemas is executable. A `TypeBinding` registry (`internal/apischema`) maps each canonical schema to its derivative Go type, and a reflection-based `VerifyGoTypeAgainstSchema` check verifies, for the supported schema subset, that the Go type matches the schema: property names, JSON tags, required-versus-optional (`omitempty`/pointer) fields, primitive types, arrays/maps, embedded (promoted) fields, enum-backed named types, and `additionalProperties` behavior. A mismatch fails the feature gate. Fixture round-tripping remains useful supporting evidence but is NOT complete proof that derivative Go types match the canonical schemas. | F12-NAMING-005; F12-VALIDATION-001(4); F12-VERIFY-001(13) |
| D-01a | Structural schema validation uses a Sovrunn-owned, explicitly bounded JSON Schema 2020-12 **supported subset** validator built on the standard library. Every canonical schema is scanned before use and any keyword outside the declared supported subset is **executably rejected (fail-closed)** with a stable code, so no schema constraint is ever silently unenforced. This is a deliberate, bounded Build (not a partial generic engine); a full generic JSON Schema engine for arbitrary documents remains deferred and, if later required, needs an approved dependency decision. The Reuse/Wrap/Extend/Build assessment for this choice is in the Validation section. | F12-NAMING-005; F12-VALIDATION-001(4); F12-VERIFY-001(13) |
| D-02 | Shared grammar lives in small, single-responsibility packages under `internal/` (`apimeta`, `apiref`, `apicond`, `apiproblem`, `apivalid`, `apischema`, `apiconform`); none imports `internal/api` or `internal/server`. | F12-IMPL-001; go-coding-guardrails §4 |
| D-03 | Decoding is provided as **pure functions separated from the HTTP adapter**: `DecodeJSON` (stdlib `encoding/json`, `DisallowUnknownFields` + token-scan duplicate-key detector) and `DecodeYAML`. YAML is treated as a **strict JSON-compatible input representation only** (D-03a), never as full YAML; it is normalized to a JSON-compatible value and then decoded through the same `DecodeJSON` path used for JSON input, ensuring identical unknown-field rejection, FieldPolicy enforcement, and stable error-code mapping. The HTTP adapter selects a decoder by media type; accepted request media types are exactly `application/json` and `application/yaml` (with `application/x-yaml` and `text/yaml` treated as `application/yaml`); all other media types map to 415. | F12-VALIDATION-001(2)/002; F12-ERROR-002 |
| D-03a | YAML input is a **strict JSON-compatible subset**. Before structural validation the YAML path MUST: (1) require exactly one YAML document (zero or multiple documents rejected); (2) require string mapping keys; (3) reject aliases, anchors, merge keys (`<<`), custom/explicit tags, non-finite numbers (`.nan`, `.inf`), and YAML-only timestamp/binary coercions; (4) perform an explicit `yaml.Node` safety and duplicate-key pass; (5) normalize the accepted YAML node to a JSON-compatible value; (6) marshal that normalized value to JSON bytes; (7) pass those JSON bytes through the same `DecodeJSON` path (unknown-field rejection, FieldPolicy enforcement, same destination Go type, same stable error-code and JSON Pointer mapping). No direct yaml.v3 typed decoding is performed after normalization. YAML struct tags are not required for validation equivalence because the normalized representation is decoded through the JSON path. JSON and YAML representations of the same object MUST produce **equivalent validation results**. | F12-VALIDATION-001(2)/(4); F12-VALIDATION-002 |
| D-04 | Validation is a nine-layer ordered pipeline. Layer 7 performs **structural** reference/kind/scope checks only; **caller-specific** cross-tenant/cross-organization authorization and no-existence-disclosure move to layer 8, which is defined as a small adopter-owned authorization interface plus a uniform safe-denial mapping (no policy engine is implemented here). Layer 9 (later-feature capability/runtime) is reserved. | F12-VALIDATION-001/004; F12-SCOPE-002; F12-SEC-004; F12-IMPL-002 |
| D-05 | Errors use RFC 9457 Problem Details with Sovrunn extensions (`code`, `requestId`, `violations[]`); violation `field` is an RFC 6901 JSON Pointer. | F12-ERROR-001/002/003 |
| D-06 | Finite platform limits are set as reviewed configuration defaults in an `apivalid.Limits` struct (values in the Data models section). | F12-VALIDATION-007; F12-LIST-002 |
| D-07 | `uid` is generated from `crypto/rand` as an opaque 128-bit value that is **collision-resistant, not collision-proof**; adopting storage MUST perform a uniqueness/collision check on persist and reject or regenerate on the (astronomically unlikely) collision. `resourceVersion` is an opaque string; both are treated as opaque by clients. No external UUID library. | F12-META-004; F12-UPDATE-002 |
| D-08 | The Sovrunn extension registry contains exactly five registered extensions: `x-sovrunn-profile`, `x-sovrunn-boundary`, `x-sovrunn-allowed-scopes`, `x-sovrunn-stability`, and `x-sovrunn-field-policy`. These are JSON Schema extension keywords validated by a fitness function against controlled vocabularies. No wildcard `x-sovrunn-*` namespace is allowed; unknown `x-sovrunn-*` extensions fail closed. `x-sovrunn-field-policy` is a strictly validated property-level object carrying or inheriting exactly: classification, authorizedWriter, authorizedReaders, mutability, retention, redaction, residency, auditRequired. | F12-NAMING-006; F12-VERIFY-001(1); F12-SEC-001 |
| D-09 | New adopting APIs use `/apis/<group>/<version>/<plural-kebab>` routes; Phase 1 routes are retained unchanged; FEATURE-0012 adds no runtime routes, only a route-form validator. | F12-NAMING-004; F12-COMPAT-003 |
| D-10 | Optimistic concurrency is a reusable `If-Match`/`resourceVersion` comparison helper mapping stale writes to 412; enforcement belongs to adopting features. | F12-UPDATE-002 |
| D-11 | A schema-diff gate classifies changes between a stored baseline snapshot and the current schema per the change-classification table, wired into the feature gate. The baseline is **immutable except through an approval-controlled workflow**. `BASELINE_MANIFEST.json` records integrity digests and detects silent baseline edits, but it is an **integrity mechanism, not an independently unforgeable approval mechanism**: because a committer can edit a baseline file and its manifest digest together, changing both in the same commit MUST NOT be sufficient to approve a baseline change. A baseline change **fails unless accompanied by recorded approval evidence** (`api/schemas/baseline/BASELINE_APPROVALS.json` or equivalent) that contains the exact old and new digests, the approving ADH or approval token, the reviewer, and the date; the gate recomputes digests and matches them to this evidence. The actual human governance boundary is **protected review / CODEOWNERS (or equivalent)** on the baseline and its approval record, which is where a human authorizes the change. | F12-EVOLVE-002; F12-VERIFY-001(10) |
| D-12 | The boundary ledger has a **machine-readable representation** (`docs/api/boundary-ledger.yaml`) governed by a strict internal ledger schema; `docs/api/BOUNDARY_LEDGER.md` is a derivative human view generated from it. A fitness function parses the ledger strictly and asserts that every declared boundary carries all F12-LEDGER-001 categories and that every boundary category present in the schemas has a ledger entry. | F12-LEDGER-001; F12-VERIFY-001 |
| D-13 | The Phase 1 compatibility report (`docs/api/PHASE1_COMPATIBILITY_REPORT.md`) records conforming behavior, exceptions, and migration candidates; a check asserts full Phase 1 coverage. | F12-COMPAT-001/002 |
| D-14 | Target Go version is 1.22 per the version standard; only the standard library and the already-present `gopkg.in/yaml.v3` are used. | Hard constraint; go-version-standard |
| D-15 | Strict decoding is **operation-aware** via a `DecodeMode`/`FieldPolicy` covering at least create request, full replacement, status update, internal object, and read representation. Customer mutation requests (create, replace) reject unauthorized system-owned and `status` fields; status-update, internal-object, read-representation, and fixture decoding accept them under the Matrix C2 ownership rules. Unconditional status/system-field rejection is removed. | F12-VALIDATION-002; F12-META-002; F12-OWNER-002 |
| D-16 | `ScopeRef` conforms to the common typed-reference contract by carrying `apiVersion`, `kind`, `name`, and optional `uid` through the shared typed-reference base (`apimeta.TypedRef`, re-exported as `apiref.TypedRef`), with `kind` constrained to the Matrix B scope kinds. To preserve the documented package layering (`apimeta` is stdlib-only; `apiref` -> `apimeta`), the base lives in `apimeta` so `ScopeRef`/`OwnerRef` embed it without an import cycle. **Canonical platform-scope representation (resolves the absent-vs-explicit ambiguity):** the single canonical stored and emitted form of platform scope is an **absent (nil) `scopeRef`**. An explicit `ScopeRef` with `Kind == "Platform"` is an accepted *input alternate* but is **normalized to nil during layer-5 defaulting, before identity, authorization, concurrency, persistence, and output processing**, so downstream stages and emitted representations only ever see the canonical nil form. A nil/normalized platform scope is valid only for resources whose `x-sovrunn-allowed-scopes` includes `Platform`. For the `API group + kind + scope UID + name` uniqueness rule (F12-SCOPE-002), platform-scoped resources use a reserved **platform-scope identity sentinel** `apimeta.PlatformScopeUID` (constant value `"platform"`, which is not a valid generated `uid`), so their identity tuple is well-defined without a scope object. This aligns `scopeRef` to F12-REF-001 and introduces no architecture exception. | F12-REF-001; F12-SCOPE-002 |
| D-17 | **Operation allowed scopes (ADH-2026-013).** The canonical generic Operation contract declares exactly six allowed scopes: Platform, Organization, OrganizationUnit, Tenant, Project, Provider. Operation.scopeRef MUST equal the resolved canonical governance scope of Operation.targetRef. For a platform-scoped target, Operation.scopeRef is canonically nil. For a non-platform target, Operation.scopeRef identifies the target's governance scope by UID. Operation.ownerRef MAY represent lifecycle containment but MUST NOT replace scopeRef or act as a governance or security scope. A target/scope mismatch is rejected with a stable validation code and an RFC 6901 JSON Pointer path. The six-value allowed-scope list does not grant authorization; target-kind constraints, caller authorization, and no-existence-disclosure rules remain mandatory. | F12-SCOPE-002; F12-REF-001; F12-FIXTURE-001; F12-FIXTURE-002 |

## Architecture

FEATURE-0012 introduces a horizontal "grammar layer" of shared primitives
that later features compose, plus a conformance layer that enforces the
grammar. The layer sits below `internal/api` handlers and does not depend on
the HTTP server or any domain registry.

```text
Canonical schemas (api/schemas/*.json, JSON Schema 2020-12)
        |  (source of truth; x-sovrunn-* annotations)
        v
Shared grammar primitives (internal/apimeta, apiref, apicond, apiproblem)
        |
        v
Strict decode + layered validation (internal/apivalid)
        |
        v
Schema metadata + diff + route-form (internal/apischema)
        |
        v
Conformance: fitness functions + fixtures + Phase 1 coverage
(internal/apiconform, tests/conformance, scripts/api-conformance-check.sh)
        |
        v
Feature gate (make ff-feature-gate FEATURE=FEATURE-0012)
```

Import direction (enforced; no cycles, and honoring the hard constraint
that `internal/api` MUST NOT import `internal/server`):

```text
internal/apimeta   -> (stdlib only)
internal/apiref    -> apimeta
internal/apicond   -> (stdlib only)
internal/apiproblem-> (stdlib only)
internal/apivalid  -> apimeta, apiref, apicond, apiproblem
internal/apischema -> apimeta (+ stdlib json)
internal/apiconform-> apimeta, apiref, apicond, apiproblem, apivalid, apischema
internal/api       -> may adopt the above; MUST NOT import internal/server
```

The grammar primitives are provider-neutral and free of Kubernetes-only or
PostgreSQL-specific assumptions. No package embeds a provider SDK type, and
no placement, policy, or provisioning logic is present. This preserves the
adapter boundaries (DEC-0036): provider/plugin/adapter-native data can only
be expressed behind `adapter-facing`/`plugin-facing` schemas that later
features own, never inside these core primitives.

The eight required fixtures (Matrix D) are represented as canonical schema
documents plus decoded Go fixtures; they prove the grammar can model future
scenarios (Project, ResourcePool, DiscoveredDatabase, PluginDefinition,
AdapterConfiguration, PlacementEvaluationRequest, Operation, AuditEvent)
without implementing any of their runtime behavior.

## Files

New and changed files. Test files are listed in the Testing Strategy
section. No file adds a runtime route or a domain service.

### Shared grammar primitives (Go)

```text
internal/apimeta/typemeta.go       # TypeMeta: apiVersion, kind; group/version parsing
internal/apimeta/objectmeta.go     # ObjectMeta: name, uid, displayName, scopeRef,
                                   #   labels, annotations, generation, resourceVersion,
                                   #   createdAt, updatedAt; ownership/mutability tags
internal/apimeta/reference.go      # TypedRef shared base (apiVersion/kind/name/uid); stdlib-only
internal/apimeta/scope.go          # ScopeRef (embeds TypedRef; Matrix B kind), OwnerRef, ScopeKind consts
internal/apimeta/profile.go        # Profile, Boundary, Stability, DataClassification enums
internal/apimeta/uid.go            # crypto/rand opaque UID generation; opaque helpers
internal/apiref/reference.go       # TypedRef alias of apimeta.TypedRef + Refs; name/uid consistency
internal/apiref/constraints.go     # allowed kinds/scopes/direction constraint helpers
internal/apicond/condition.go      # Condition type + Set/Get/merge helpers; status enum
internal/apiproblem/problem.go     # RFC 9457 Problem + Violation + Sovrunn extensions
internal/apiproblem/codes.go       # stable ErrorCode + violation-code registry
internal/apiproblem/httpmap.go     # failure-class -> HTTP status mapping
internal/apivalid/decode.go        # pure DecodeJSON (encoding/json, DisallowUnknownFields,
                                   #   duplicate-key token scan) + DecodeYAML (yaml.v3 node parse,
                                   #   YAML-safety rejection, normalize to JSON, then DecodeJSON);
                                   #   no HTTP dependency
internal/apivalid/httpdecode.go    # HTTP adapter: MaxBytes, media-type selection -> DecodeJSON/DecodeYAML
internal/apivalid/fieldpolicy.go   # DecodeMode + FieldPolicy (create/replace/status/internal/read)
internal/apivalid/pipeline.go      # nine-layer ordered validation pipeline + Result
internal/apivalid/structural.go    # StructuralValidator interface (apivalid MUST NOT import apischema)
internal/apivalid/limits.go        # Limits config struct + reviewed defaults
internal/apivalid/concurrency.go   # If-Match/resourceVersion comparison -> 412 helper
internal/apivalid/authz.go         # layer-8 adopter-owned ScopeAuthorizer +
                                   #   AuthorizedResolver + AuthorizedTargetScopeResolver
                                   #   interfaces + safe-denial map
                                   #   (authorize-before-lookup / combined-resolver /
                                   #   target-scope-resolver contract)
internal/apischema/annotations.go  # x-sovrunn-* keyword parsing + vocabulary checks
internal/apischema/subsetvalidate.go # bounded JSON Schema 2020-12 subset validator; rejects
                                   #   unsupported keywords fail-closed (D-01a)
internal/apischema/typebinding.go  # TypeBinding registry + reflection-based
                                   #   VerifyGoTypeAgainstSchema Go-type consistency check (D-01b)
internal/apischema/route.go        # /apis/<group>/<version>/<plural-kebab> route-form validator
internal/apischema/diff.go         # schema-diff change classifier + baseline-integrity check
internal/apiconform/structural.go  # StructuralValidator adapter: apischema -> apiproblem.Violation translation
internal/apiconform/fitness.go     # executable fitness functions (F12-VERIFY-001 checks)
internal/apiconform/fixtures.go    # fixture loader + Matrix D scenario assertions
internal/apiconform/compat.go      # Phase 1 coverage assertion for the report
```

### Canonical schemas, fixtures, and baselines (data)

```text
api/schemas/project.json                    # ManagedResource / customer-facing / Tenant
api/schemas/resource-pool.json              # ManagedResource / operator-facing / Provider
api/schemas/discovered-database.json        # ObservedExternalResource / adapter-facing
api/schemas/plugin-definition.json          # VersionedDefinition / plugin-facing / Platform
api/schemas/adapter-configuration.json      # ManagedResource / adapter-facing / Provider
api/schemas/placement-evaluation-request.json # TransientRequestResult / internal-engine
api/schemas/operation.json                  # LongRunningOperation / plugin/operator / Platform,Organization,OrganizationUnit,Tenant,Project,Provider
api/schemas/audit-event.json                # ImmutableRecord / governance-only
api/schemas/_common/*.json                  # shared metadata/reference/condition/problem sub-schemas
api/schemas/baseline/*.json                 # frozen snapshots for the schema-diff gate
api/schemas/baseline/BASELINE_MANIFEST.json # integrity digests for baselines (tamper detection, not approval)
api/schemas/baseline/BASELINE_APPROVALS.json # recorded baseline-change approval evidence
                                   #   (old/new digests, approving ADH/token, reviewer, date)
tests/conformance/fixtures/*.json           # valid + invalid instances per contract
```

### Documentation, scripts, gate

```text
docs/api/boundary-ledger.yaml               # machine-readable boundary ledger, source of truth (F12-LEDGER-001)
docs/api/BOUNDARY_LEDGER.md                 # derivative human view generated from boundary-ledger.yaml
docs/api/PHASE1_COMPATIBILITY_REPORT.md     # Phase 1 compatibility report (F12-COMPAT-001)
scripts/api-conformance-check.sh            # runs fitness functions + schema-diff + coverage
scripts/feature-gate.sh                     # extended to invoke api-conformance-check for FEATURE-0012
```

## Data models

All shared types use explicit JSON tags in lowerCamelCase (F12-NAMING-002),
`omitempty` only where absence is meaningful, and never expose secret
values (F12-SEC-003). Types are illustrative signatures; field docs carry
ownership/mutability annotations mirroring Matrix C2.

### Type identity and common metadata (internal/apimeta)

```go
// TypeMeta identifies the contract of an externally exchanged object.
type TypeMeta struct {
    APIVersion string `json:"apiVersion"` // <domain>.sovrunn.io/{v1alpha1|v1beta1|v1}
    Kind       string `json:"kind"`       // singular PascalCase
}

// ObjectMeta is the applicable subset of common metadata for persistent
// resources. Ownership/mutability is enforced by validation, not by type.
type ObjectMeta struct {
    Name            string            `json:"name"`                      // creator, immutable
    UID             string            `json:"uid,omitempty"`             // Sovrunn, immutable, opaque
    DisplayName     string            `json:"displayName,omitempty"`     // owner, mutable; not identity
    ScopeRef        *ScopeRef         `json:"scopeRef,omitempty"`        // creator on create; normally immutable
    Labels          map[string]string `json:"labels,omitempty"`         // bounded; no secrets
    Annotations     map[string]string `json:"annotations,omitempty"`    // namespaced; bounded; no secrets
    Generation      int64             `json:"generation,omitempty"`      // system-only
    ResourceVersion string            `json:"resourceVersion,omitempty"` // system-only, opaque
    CreatedAt       string            `json:"createdAt,omitempty"`       // system-only, UTC RFC 3339
    UpdatedAt       string            `json:"updatedAt,omitempty"`       // system-only, UTC RFC 3339
}
```

### Scope, ownership, and enumerations (internal/apimeta)

The **shared typed-reference base** is defined in `apimeta` — the lowest,
stdlib-only package — so that `ScopeRef`, `OwnerRef`, and the `apiref`
domain-specific references can all build on one base without breaking the
documented import direction (`apimeta` -> stdlib only; `apiref` -> `apimeta`).
`apiref` re-exports this base as `apiref.TypedRef` (a type alias) and adds
the constraint/collection helpers; no cycle is introduced.

```go
// TypedRef is the common typed-reference base (F12-REF-001): apiVersion,
// kind, name, and optional immutable uid. It is stdlib-only and lives in
// apimeta so both scope/owner references and apiref aliases embed it.
type TypedRef struct {
    APIVersion string `json:"apiVersion"`
    Kind       string `json:"kind"`
    Name       string `json:"name"`
    UID        string `json:"uid,omitempty"` // optional immutable; must agree with name
}

type ScopeKind string // Matrix B — the only valid scopeRef.kind values

const (
    ScopePlatform         ScopeKind = "Platform"
    ScopeOrganization     ScopeKind = "Organization"
    ScopeOrganizationUnit ScopeKind = "OrganizationUnit"
    ScopeTenant           ScopeKind = "Tenant"
    ScopeProject          ScopeKind = "Project"
    ScopeProvider         ScopeKind = "Provider"
)

// ScopeRef is the immutable primary security/governance ownership reference.
// It is NOT a location and NOT a lifecycle-containment reference.
//
// ScopeRef conforms to the common typed-reference contract (F12-REF-001) by
// carrying apiVersion, kind, name, and optional immutable uid through the
// shared TypedRef base. Its Kind is additionally constrained to the Matrix B
// scope kinds by validation (see below); embedding the shared base keeps
// scopeRef a first-class typed reference rather than an architecture exception.
type ScopeRef struct {
    // Anonymous embedding: encoding/json promotes the embedded fields
    // (apiVersion, kind, name, uid) with no tag. DecodeYAML normalizes
    // accepted YAML to JSON and then uses DecodeJSON, so no YAML inline
    // tag participates in typed decoding.
    TypedRef
    // Kind MUST be one of the Matrix B ScopeKind values; enforced in validation.
    // authorization resolves by uid, not name.
}

// CANONICAL PLATFORM SCOPE (F12-SCOPE-002, D-16): the single canonical stored
// and emitted form of platform scope is an ABSENT (nil) *ScopeRef. An explicit
// ScopeRef with Kind == "Platform" is an accepted INPUT ALTERNATE only; it is
// normalized to nil during layer-5 defaulting BEFORE identity, authorization,
// concurrency, persistence, and output processing, so no downstream stage or
// emitted response ever observes two forms of platform scope. A nil/normalized
// platform scope is valid ONLY when the schema's x-sovrunn-allowed-scopes
// includes "Platform". For any resource whose allowed scopes do not include
// Platform, a nil scopeRef is rejected and exactly one immutable primary
// scopeRef is required.
//
// PlatformScopeUID is the reserved platform-scope identity sentinel used by the
// "API group + kind + scope UID + name" uniqueness rule for platform-scoped
// resources; it is a fixed constant that can never collide with a generated uid.
const PlatformScopeUID = "platform"

// NormalizeScope returns the canonical scope form: it maps an explicit
// Kind=="Platform" ScopeRef to nil and leaves all other scopes unchanged. It
// runs during layer-5 defaulting so identity, authorization, concurrency,
// persistence, and output all operate on the canonical representation.
func NormalizeScope(s *ScopeRef) *ScopeRef

// OwnerRef expresses resource-local lifecycle containment only (F12-OWNER-001).
// It MUST NOT be used as a scopeRef.kind or as a security/governance scope.
// Like all references it uses the shared typed-reference base.
type OwnerRef struct {
    TypedRef // apiVersion, kind, name, optional uid; JSON promotes embedded fields
}
```

```go
type Profile string // Matrix A — exactly one per externally exchanged object

const (
    ProfileManagedResource         Profile = "ManagedResource"
    ProfileObservedExternalResource Profile = "ObservedExternalResource"
    ProfileVersionedDefinition     Profile = "VersionedDefinition"
    ProfileImmutableRecord         Profile = "ImmutableRecord"
    ProfileLongRunningOperation    Profile = "LongRunningOperation"
    ProfileTransientRequestResult  Profile = "TransientRequestResult"
    ProfileEmbeddedValue           Profile = "EmbeddedValue"
    ProfileListEnvelope            Profile = "ListEnvelope"
)

type Boundary string // Matrix C1

const (
    BoundaryCustomerFacing       Boundary = "customer-facing"
    BoundaryOperatorFacing       Boundary = "operator-facing"
    BoundaryInternalEngineFacing Boundary = "internal-engine-facing"
    BoundaryAdapterFacing        Boundary = "adapter-facing"
    BoundaryPluginFacing         Boundary = "plugin-facing"
    BoundaryGovernanceOnly       Boundary = "governance-only"
)

type Stability string // maturity / compatibility expectation

const (
    StabilityAlpha  Stability = "alpha"
    StabilityBeta   Stability = "beta"
    StabilityStable Stability = "stable"
)

type DataClassification string // F12-SEC-002 — exactly seven values

const (
    ClassPublic              DataClassification = "Public"
    ClassCustomerVisible     DataClassification = "Customer-visible"
    ClassTenantConfidential  DataClassification = "Tenant-confidential"
    ClassOperatorConfidential DataClassification = "Operator-confidential"
    ClassInternal            DataClassification = "Internal"
    ClassSensitive           DataClassification = "Sensitive"
    ClassSecretReferenceOnly DataClassification = "Secret-reference-only"
)
```

### References (internal/apiref)

```go
// TypedRef is re-exported from apimeta as the common typed-reference base
// (defined in apimeta to avoid an import cycle; see the Scope/ownership
// block). Singular fields end in Ref; collections end in Refs. Provider-native
// IDs MUST NOT act as core refs. apiref adds the Constraint and collection
// (Refs) helpers on top of this base.
type TypedRef = apimeta.TypedRef

// Constraint restricts a reference field's allowed kinds, scopes, and
// direction. Public schemas SHOULD expose domain-specific aliases built on
// this base (F12-REF-004).
type Constraint struct {
    AllowedKinds  []string
    AllowedScopes []ScopeKind
    Direction     Direction // Inbound, Outbound, or Bidirectional
}
```

### Conditions (internal/apicond)

```go
type ConditionStatus string // "True" | "False" | "Unknown"

// Condition is a stable, machine-readable current-fact observation.
// It is NOT event history (F12-STATUS-002/004).
type Condition struct {
    Type               string          `json:"type"`               // stable PascalCase
    Status             ConditionStatus `json:"status"`
    Reason             string          `json:"reason"`             // stable PascalCase
    Message            string          `json:"message,omitempty"`  // human-readable; not a contract
    ObservedGeneration int64           `json:"observedGeneration"` // evaluated desired-state generation
    LastTransitionTime string          `json:"lastTransitionTime"` // changes only on status change
}
```

### Problem details (internal/apiproblem)

```go
// Problem is an RFC 9457 Problem Details response with Sovrunn extensions.
type Problem struct {
    Type       string      `json:"type"`
    Title      string      `json:"title"`
    Status     int         `json:"status"`
    Detail     string      `json:"detail,omitempty"`
    Instance   string      `json:"instance,omitempty"`
    Code       string      `json:"code"`                 // stable machine contract
    RequestID  string      `json:"requestId,omitempty"`
    Violations []Violation `json:"violations,omitempty"`
}

// Violation identifies one invalid field by RFC 6901 JSON Pointer.
type Violation struct {
    Field   string `json:"field"`   // JSON Pointer, e.g. /spec/storage/sizeGiB
    Code    string `json:"code"`    // stable violation code
    Message string `json:"message"` // human-readable; redacted of sensitive detail
}
```

### List envelope (internal/apimeta)

```go
// ListEnvelope is the paginated collection response profile.
//
// Embedding note (F12-NAMING-005 correctness): the standard library
// encoding/json does NOT honor a `json:",inline"` tag. An anonymous embedded
// struct without a JSON tag already promotes its fields (apiVersion, kind)
// to the enclosing JSON object, which is the intended flat representation.
// DecodeYAML does not perform direct yaml.v3 typed decoding: accepted YAML
// is normalized to JSON and decoded through DecodeJSON, so YAML struct tags
// are neither required nor used by the typed decoder.
type ListEnvelope[T any] struct {
    TypeMeta // JSON fields promoted through anonymous embedding
    Items    []T  `json:"items"`
    Page     Page `json:"page"`
}

type Page struct {
    NextPageToken string `json:"nextPageToken,omitempty"` // opaque; no offsets/provider detail
    // total count intentionally optional and omitted by default (F12-LIST-002)
}
```

### Finite limits (internal/apivalid) — reviewed platform configuration

Per F12-VALIDATION-007 the architecture authorizes design to set initial
finite limits as reviewed configuration. These defaults are overridable via
validated configuration; they are not runtime-provider-specific.

```go
type Limits struct {
    MaxObjectBytes        int // 1_048_576  (1 MiB; matches existing 1<<20)
    MaxNestingDepth       int // 32
    MaxLabels             int // 64
    MaxLabelKeyChars      int // 63
    MaxLabelValueChars    int // 253
    MaxAnnotationsBytes   int // 262_144  (256 KiB total)
    MaxConditions         int // 32
    MaxReferencesPerField int // 64
    MaxViolations         int // 100
    DefaultPageSize       int // 50
    MaxPageSize           int // 200
}
```

| Limit | Default | Rationale |
|---|---:|---|
| Object size | 1 MiB | Consistent with the existing Phase 1 `MaxBytesReader` bound. |
| Nesting depth | 32 | Prevents pathological documents; ample for known contracts. |
| Labels | 64 | Bounded, indexed classification. |
| Label key/value | 63 / 253 | DNS-style key and bounded value. |
| Annotations total | 256 KiB | Bounded non-indexed metadata. |
| Conditions | 32 | Conditions are current facts, not history. |
| References per field | 64 | Bounds fan-out and resolution cost. |
| Violations returned | 100 | Bounds error payloads. |
| Page size (default/max) | 50 / 200 | Safe portal paging without leaking offsets. |

Exceeding any limit fails deterministically with a stable code and JSON
Pointer path (F12-VALIDATION-006/007, F12-LIST-002).

## Components and Interfaces

Interfaces are small and context-aware where request-scoped
(go-coding-guardrails §6). Layer 8 implements only the shared orchestration,
interfaces, configuration checks, and safe-denial mapping; adopter-specific
authorization and resolution behavior remains out of scope. Layer 9 is
reserved for later-feature runtime (F12-IMPL-002).

### Strict decoding (internal/apivalid)

Decoding is split into **pure functions** (no HTTP dependency) and a thin
**HTTP adapter**. The pure functions are directly testable offline and are
reused for fixture decoding (F12-VALIDATION-005).

```go
// DecodeMode makes decoding operation-aware (F12-VALIDATION-002, F12-META-002,
// F12-OWNER-002). It selects the FieldPolicy that governs which ownership
// classes of fields are accepted or rejected.
type DecodeMode int
const (
    ModeCreateRequest     DecodeMode = iota // customer/operator create: reject system-owned + status
    ModeReplaceRequest                      // full replacement: reject status + immutable system fields
    ModeStatusUpdate                        // authorized controller: accept status, reject spec mutation
    ModeInternalObject                      // internal/system producer: accept system-owned fields
    ModeReadRepresentation                  // decode a stored/response object: accept all fields
)

// FieldPolicy resolves, for a DecodeMode, which field ownership classes
// (per Matrix C2: creator, system, spec-owner, status-owner) are permitted.
// Customer mutation modes (create, replace) reject unauthorized system/status
// fields; internal, status-update, read, and fixture decoding accept them
// under the correct ownership rules. There is no unconditional rejection.
type FieldPolicy struct {
    Mode              DecodeMode
    AllowStatus       bool
    AllowSystemOwned  bool // uid, generation, resourceVersion, timestamps
    AllowSpecMutation bool
}
func PolicyFor(mode DecodeMode) FieldPolicy

// DecodeJSON is pure: safe JSON decoding via encoding/json with
// DisallowUnknownFields plus a token-scan duplicate-key detector, applying the
// FieldPolicy. No HTTP dependency. Errors carry a stable code + JSON Pointer.
func DecodeJSON(data []byte, lim Limits, pol FieldPolicy, dst any) *apiproblem.Problem

// DecodeYAML is pure and treats YAML as a STRICT JSON-COMPATIBLE input
// representation only (D-03a). Using gopkg.in/yaml.v3 for safe syntax-tree
// parsing and normalization only, it, in order: (1) requires exactly one YAML
// document; (2) requires string mapping keys; (3) rejects aliases, anchors,
// merge keys (<<), custom/explicit tags, non-finite numbers (.nan/.inf), and
// YAML-only timestamp/binary coercions; (4) runs an explicit yaml.Node safety
// + duplicate-key pass; (5) normalizes the accepted YAML node to a
// JSON-compatible value; (6) marshals that normalized value to JSON bytes;
// (7) passes those JSON bytes through the same DecodeJSON path (unknown-field
// rejection, FieldPolicy enforcement, same destination Go type, same stable
// error-code and JSON Pointer mapping). No direct yaml.v3 typed decoding is
// performed after normalization. YAML struct tags are not required for
// validation equivalence because the normalized representation is decoded
// through the JSON path. JSON and YAML forms of the same object produce
// EQUIVALENT validation results. No HTTP dependency. Errors carry a stable
// code + JSON Pointer.
func DecodeYAML(data []byte, lim Limits, pol FieldPolicy, dst any) *apiproblem.Problem

// StrictDecode is the HTTP adapter (internal/apivalid/httpdecode.go). It
// enforces the size limit (http.MaxBytesReader), selects a decoder by media
// type — application/json and application/yaml only (application/x-yaml and
// text/yaml are treated as application/yaml); any other type maps to 415 —
// and delegates to DecodeJSON/DecodeYAML. It never mutates input semantics.
func StrictDecode(w http.ResponseWriter, r *http.Request, lim Limits, mode DecodeMode, dst any) *apiproblem.Problem
```

### Layered validation pipeline (internal/apivalid)

```go
// Layer is one ordered validation stage (F12-VALIDATION-001).
type Layer int
const (
    LayerHTTPContent Layer = iota + 1 // 1 HTTP/content/size
    LayerDecode                       // 2 safe decode
    LayerFieldHygiene                 // 3 duplicate/unknown-field rejection
    LayerStructural                   // 4 structural schema validation
    LayerDefaulting                   // 5 deterministic defaulting
    LayerSemantic                     // 6 semantic validation
    LayerReference                    // 7 STRUCTURAL reference/kind/scope validation
    LayerAuthorization                // 8 caller-specific authz + no-existence-disclosure (adopter-owned interface)
    LayerCapabilityRuntime            // 9 later-feature runtime — RESERVED
)

// Result keeps failure classes distinguishable (F12-VALIDATION-004).
//
// Problem is the safe client-facing failure contract (e.g. 500 INTERNAL_ERROR
// when a StructuralValidator is nil or returns an error). Err is internal
// diagnostic context and MUST NOT be serialized or exposed to callers.
//
// Binding rules:
//   - On a nil StructuralValidator or validator error at layer 4:
//       FailedAt = LayerStructural;
//       Problem  = a 500 INTERNAL_ERROR Problem;
//       Err      = the internal cause (nil validator or validator error);
//       no success result is returned;
//       layers 5 through 7 do not execute.
//   - Ordinary structural violations populate Violations and the normal 422
//     VALIDATION_FAILED mapping; Problem and Err are nil.
//   - Successful validation: Violations is empty, Problem is nil, Err is nil.
type Result struct {
    Violations []apiproblem.Violation
    FailedAt   Layer
    Problem    *apiproblem.Problem // safe client-facing failure; nil on success or ordinary violations
    Err        error               // internal diagnostic; MUST NOT be serialized or exposed
}

// Input carries everything the pipeline requires for a single validation
// invocation. StructuralValidator MUST be non-nil for any pipeline run that
// processes an external object; a nil validator at layer 4 causes the
// pipeline to stop with Result.Problem set to 500 INTERNAL_ERROR and
// Result.Err recording the internal cause.
//
// LAYER-8 CONFIGURATION MATRIX:
//
// When OperationScope is non-nil, exactly one of Path A or Path B MUST be
// completely configured. Configuring both TargetScope and TargetScopeResolver
// is invalid. Configuring neither is invalid. Missing TargetRef is invalid.
// Missing Caller is invalid. An incomplete path configuration is invalid.
// Every invalid layer-8 configuration stops at LayerAuthorization with
// Result.FailedAt = LayerAuthorization, Result.Problem = 500 INTERNAL_ERROR,
// Result.Err = internal configuration cause, no target lookup, no success
// result, and no silent skip.
//
// Path A — authoritative target scope derivable without lookup:
//   Required: OperationScope, TargetRef, TargetScope, Authorizer, Caller.
//   TargetScope MUST come from an authoritative route binding or previously
//   validated trusted context, not an arbitrary unverified caller value.
//   ScopeAuthorizer runs before any target lookup. Denial maps through
//   SafeDenial. Allow permits CheckOperationTargetScopeMatch.
//
// Path B — target scope requires authorized lookup:
//   Required: OperationScope, TargetRef, TargetScopeResolver, Caller.
//   A separate ScopeAuthorizer is not required because
//   AuthorizedTargetScopeResolver performs combined authorized resolution.
//   available=false maps through SafeDenial. available=true permits
//   CheckOperationTargetScopeMatch.
//
// Generic non-Operation validation (OperationScope is nil and no layer-8
// capability was requested) MAY stop successfully after layer 7.
type Input struct {
    Data        []byte
    Mode        DecodeMode
    SchemaID    string
    Validator   StructuralValidator // required; nil = unavailable = 500
    Authorizer  ScopeAuthorizer     // required for Path A; nil when using Path B only
    Caller      *CallerContext      // required when OperationScope is non-nil
    Dst         any                 // decode target

    // Operation target-scope equality (layer 8). These fields are governed
    // by the layer-8 configuration matrix above. When OperationScope is
    // non-nil, exactly one of Path A (TargetScope set) or Path B
    // (TargetScopeResolver set) MUST be completely configured; incomplete or
    // contradictory configuration is a 500 INTERNAL_ERROR at
    // LayerAuthorization (no silent skip, no target lookup).
    OperationScope      *ScopeIdentity                // Operation's canonical scope for the equality check
    TargetRef           *apiref.TypedRef              // Operation.targetRef; required when OperationScope is non-nil
    TargetScope         *ScopeIdentity                // pre-derived authoritative target scope (Path A)
    TargetScopeResolver AuthorizedTargetScopeResolver // combined authorized resolver (Path B)
}

// Validate runs layers 1..7 deterministically and safely offline
// (F12-VALIDATION-005). Layer 7 is STRUCTURAL only (well-formed references,
// allowed kinds/scopes/direction, name/uid agreement). At layer 4, if
// Input.Validator is nil OR if Validator.Validate returns a non-nil error,
// the pipeline MUST stop at LayerStructural, set Result.Problem to a 500
// INTERNAL_ERROR Problem, set Result.Err to the internal cause, and MUST NOT
// execute layers 5 through 7.
//
// Layer 8 behavior is governed by the layer-8 configuration matrix (see
// Input). When OperationScope is non-nil, exactly one of Path A or Path B
// MUST be completely configured; an invalid configuration stops at
// LayerAuthorization with Result.Problem = 500 INTERNAL_ERROR,
// Result.Err = internal configuration cause, no target lookup, no success
// result, and no silent skip. A valid Path A runs ScopeAuthorizer before
// lookup; denial maps through SafeDenial; Allow permits
// CheckOperationTargetScopeMatch. A valid Path B uses
// AuthorizedTargetScopeResolver; available=false maps through SafeDenial;
// available=true permits CheckOperationTargetScopeMatch.
//
// Generic non-Operation validation (OperationScope is nil and no layer-8
// capability was requested) MAY stop successfully after layer 7.
// Layer 9 is reserved.
func Validate(ctx context.Context, in Input, lim Limits) Result
```

### Layer 8: caller-specific authorization (adopter-owned interface)

FEATURE-0012 defines a small authorization boundary and a uniform safe-denial
mapping; it implements **no policy engine** (F12-IMPL-002). Cross-tenant and
cross-organization decisions and no-existence-disclosure behavior are layer 8,
not layer 7.

**Response mapping alone is insufficient.** Mapping every denied cross-scope
outcome to an identical 404 hides existence in the *response body, code, and
title*, but it does NOT by itself hide existence in **control flow and
timing**: an implementation that looks a resource up first and only then
authorizes can still leak existence through observable latency or side
effects (e.g. a slow "found then denied" path versus a fast "absent" path).
The grammar therefore defines an **adopter contract** for the ordering, not
just the response shape:

- **Authorize-before-lookup (preferred, when scope is derivable):** when the
  target scope can be derived from the request (route scope, supplied
  `scopeRef`, or caller context) without reading the object, the adopter MUST
  run `ScopeAuthorizer.Authorize` first and only perform resolution/lookup
  after an `Allow`. A cross-scope denial then never triggers a lookup, so no
  existence-dependent work occurs.
- **Combined authorized resolver (when scope is only knowable after lookup):**
  when the target scope cannot be derived without reading the object, the
  adopter MUST use a single combined resolver that returns one uniform
  "unavailable" outcome for BOTH "not found" and "found but unauthorized",
  hiding which of resolution or authorization failed. The resolver MUST NOT
  branch observably between the two cases.
- **No existence-dependent fast paths:** adopters MUST NOT add caches,
  early returns, short-circuits, or differing error/log/audit paths that make
  the "exists but inaccessible" case observably different from the "absent"
  case. Denied and absent outcomes take structurally equivalent paths.

This is a **path- and response-equivalence** contract, not a perfect
constant-time guarantee: FEATURE-0012 does not claim constant-time execution
or defend against fine-grained microarchitectural timing attacks. It requires
that denial and absence follow the same code path and produce identical
responses, and provides response/path-equivalence tests (below) rather than
statistical timing proofs. `SafeDenial` supplies the uniform response; the
ordering contract above supplies the equivalent control flow.

```go
// ScopeIdentity is a canonical value representation of a governance scope,
// usable for authorization comparison without requiring a full ScopeRef
// pointer. It avoids the nil-vs-non-nil ambiguity of *ScopeRef for
// platform scope and enables direct equality comparison.
type ScopeIdentity struct {
    Kind apimeta.ScopeKind
    UID  string // PlatformScopeUID for platform; target scope UID otherwise
}

// CanonicalScopeIdentity converts a *ScopeRef to a ScopeIdentity:
//   - nil scopeRef -> ScopeIdentity{Kind: ScopePlatform, UID: PlatformScopeUID}
//   - non-platform scopeRef -> ScopeIdentity{Kind: ref.Kind, UID: ref.UID}
func CanonicalScopeIdentity(s *apimeta.ScopeRef) ScopeIdentity

// CallerContext is the minimal request-scoped identity/scope the authorizer
// needs. It is provided by the adopting feature, not resolved here.
type CallerContext struct {
    Scopes []ScopeIdentity // scopes the caller is entitled to
}

// Decision is the outcome of an authorization check. Denials are mapped
// uniformly so inaccessible objects are never disclosed (F12-SEC-004).
type Decision int
const (
    Allow            Decision = iota
    DenyNotDisclosed          // cross-scope: MUST be indistinguishable from "not found"
    DenyKnown                 // in-scope authorization denial where existence is already known
)

// ScopeAuthorizer is implemented by adopting features. FEATURE-0012 ships the
// interface and the safe-denial mapping only — no concrete authorizer. When
// the target scope is derivable from the request without reading the object,
// adopters MUST call Authorize BEFORE lookup so a cross-scope denial performs
// no existence-dependent work.
type ScopeAuthorizer interface {
    Authorize(ctx context.Context, caller CallerContext, target apiref.TypedRef, targetScope ScopeIdentity) Decision
}

// AuthorizedResolver is the required contract when the target scope is only
// knowable AFTER reading the object. Adopters implement Resolve so that a
// missing object and a present-but-unauthorized object return the SAME uniform
// unavailable outcome (found=false with no leaked detail), hiding whether
// resolution or authorization failed. Implementations MUST NOT branch
// observably (timing, side effects, logs, audit) between the two cases.
// FEATURE-0012 defines the contract and the equivalence tests; adopters supply
// the implementation (no policy engine here).
type AuthorizedResolver interface {
    Resolve(ctx context.Context, caller CallerContext, target apiref.TypedRef) (obj any, found bool)
}

// SafeDenial maps a Decision to a uniform Problem. DenyNotDisclosed always maps
// to an identical 404 RESOURCE_NOT_FOUND (same code, title, and message) so
// that "exists but inaccessible" and "does not exist" are indistinguishable in
// the RESPONSE (F12-SCOPE-002, F12-SEC-004). DenyKnown maps to 403
// AUTHORIZATION_DENIED. NOTE: this response mapping alone is NOT sufficient for
// no-existence-disclosure — control-flow/timing equivalence is the adopter's
// responsibility via authorize-before-lookup or AuthorizedResolver (see the
// adopter contract above). This is a path/response-equivalence guarantee, not
// a perfect constant-time guarantee.
//
// SafeDenial is owned and tested under internal/apivalid/authz.go, not
// internal/apiproblem.
func SafeDenial(d Decision) *apiproblem.Problem
```

#### Operation target-scope equality (layer 8, adopter-owned)

Layers 1 through 7 are offline: layer 7 validates Operation scope syntax,
allowed-scope membership, targetRef shape, and UID requirements. It does NOT
resolve the target from external state and does NOT claim to perform
target-scope equality as an offline check.

An adopter-owned **authorized target-scope resolver** operates at layer 8.
It MUST preserve the authorize-before-lookup / combined-resolver
no-existence-disclosure contract defined above: unauthorized or unavailable
targets use `SafeDenial` and MUST NOT disclose a scope mismatch.

```go
// AuthorizedTargetScopeResolver resolves the canonical governance scope of
// an Operation's target reference WITH authorization. It is the adopter-owned
// contract that layer 8 uses for Operation target-scope equality when the
// target scope requires lookup (path B).
//
// Binding semantics:
//   - available=false is the single uniform result for BOTH:
//       (a) target absent;
//       (b) target present but unauthorized.
//     Those two cases MUST follow the same path and map through SafeDenial.
//     No target scope or mismatch detail is returned when available=false.
//   - available=true means an authorized canonical ScopeIdentity is
//     available and the caller may proceed to CheckOperationTargetScopeMatch.
//
// Implementations MUST NOT branch observably (timing, side effects, logs,
// audit) between the absent and unauthorized cases.
//
// Owned by internal/apivalid/authz.go. FEATURE-0012 defines the interface;
// adopters supply the implementation (no policy engine here).
type AuthorizedTargetScopeResolver interface {
    ResolveAuthorizedTargetScope(
        ctx context.Context,
        caller CallerContext,
        target apiref.TypedRef,
    ) (scope ScopeIdentity, available bool)
}
```

Two approved no-existence-disclosure paths for Operation target-scope
equality:

**Path A — authoritative target scope derivable without lookup:**

Required Input fields: OperationScope, TargetRef, TargetScope, Authorizer,
Caller.

1. TargetScope MUST come from an authoritative route binding or previously
   validated trusted context, not an arbitrary unverified caller value.
2. Run `ScopeAuthorizer.Authorize` before any target lookup; a
   `DenyNotDisclosed` maps through `SafeDenial`.
3. After `Allow`, compare using `CheckOperationTargetScopeMatch`.

**Path B — target scope requires authorized lookup:**

Required Input fields: OperationScope, TargetRef, TargetScopeResolver,
Caller.

1. A separate ScopeAuthorizer is not required because
   `AuthorizedTargetScopeResolver` performs combined authorized resolution.
2. Use `AuthorizedTargetScopeResolver.ResolveAuthorizedTargetScope`.
3. `available=false` maps through `SafeDenial`; no mismatch is disclosed.
4. `available=true` permits the pure comparison via
   `CheckOperationTargetScopeMatch`.

**Fail-closed configuration rules:**

When OperationScope is non-nil, exactly one of Path A or Path B MUST be
completely configured. Configuring both TargetScope and TargetScopeResolver
is invalid. Configuring neither is invalid. Missing TargetRef is invalid.
Missing Caller is invalid. An incomplete Path A or Path B configuration is
invalid. Every invalid layer-8 configuration stops at LayerAuthorization
with: Result.FailedAt = LayerAuthorization, Result.Problem = 500
INTERNAL_ERROR, Result.Err = internal configuration cause, no target lookup,
no success result, no silent skip.

Generic non-Operation validation (OperationScope is nil and no layer-8
capability was requested) MAY stop successfully after layer 7.

After an authorized target scope is available, a pure helper compares the
Operation's canonical `ScopeIdentity` with the target's canonical
`ScopeIdentity`:

```go
// CheckOperationTargetScopeMatch compares the Operation's canonical scope
// with the resolved target's canonical scope. If they differ, it returns a
// violation with code OPERATION_TARGET_SCOPE_MISMATCH and JSON Pointer
// /metadata/scopeRef. Returns nil when scopes match.
func CheckOperationTargetScopeMatch(
    opScope ScopeIdentity,
    targetScope ScopeIdentity,
) *apiproblem.Violation
```

Invariants:

- Unauthorized or unavailable targets MUST use SafeDenial and MUST NOT
  disclose a scope mismatch.
- Only after authorized target resolution succeeds does the pure comparator
  execute.
- The mismatch violation uses code `OPERATION_TARGET_SCOPE_MISMATCH` and
  JSON Pointer `/metadata/scopeRef`.

### Structural validation bridge (internal/apivalid, internal/apiconform)

The package import direction requires that `apivalid` MUST NOT import
`apischema`. Layer 4 (structural schema validation) therefore uses a
`StructuralValidator` interface owned by `apivalid`; `apiconform` implements
the adapter that delegates to `apischema` and translates package-local
`apischema.SchemaIssue` values to `apiproblem.Violation` values.

```go
// StructuralValidator is the layer-4 structural validation contract owned
// by apivalid. Implementations live in apiconform (using apischema) so
// apivalid never imports apischema directly.
//
// The Validate method returns both violations and an error. The error
// signals validator unavailability or configuration failure (e.g. schema
// not found, registry misconfiguration). Violations represent ordinary
// schema-mismatch findings. This separation makes failure distinguishable
// from clean validation:
//
//   - err != nil: structural validation was UNAVAILABLE; the pipeline MUST
//     stop at LayerStructural, set Result.Problem to a 500 INTERNAL_ERROR
//     Problem, set Result.Err to the internal cause, return no success
//     result, and MUST NOT execute layers 5 through 7.
//   - err == nil, len(violations) > 0: ordinary schema violations.
//   - err == nil, len(violations) == 0: instance is structurally valid.
//
// A nil StructuralValidator in Input at layer 4 is treated identically to
// a non-nil validator returning an error: the pipeline stops, sets
// Result.Problem to 500 INTERNAL_ERROR, sets Result.Err, and never
// executes layers 5 through 7.
//
// Primitive unit tests MAY call individual primitive functions directly or
// inject a deterministic stub StructuralValidator. No full pipeline
// invocation for an external object may omit the validator.
type StructuralValidator interface {
    Validate(instance any, schemaID string) ([]apiproblem.Violation, error)
}
```

### Bounded schema-subset validation (internal/apischema)

```go
// SchemaIssue is a package-local diagnostic value representing a schema or
// instance validation finding. apischema MUST NOT import apiproblem; the
// translation from SchemaIssue to apiproblem.Violation is owned by
// apiconform.
type SchemaIssue struct {
    Path    string // RFC 6901 JSON Pointer
    Code    string // stable machine-readable code
    Message string // human-readable detail
}
```

```go
// SupportedKeywords is the explicit, bounded JSON Schema 2020-12 subset the
// Sovrunn validator supports. FEATURE-0012 supports exactly:
//   $schema, $id, $ref, title, description, type, properties, required,
//   enum, items, additionalProperties, minLength, maxLength, minimum,
//   maximum, pattern, default, examples
// plus the registered x-sovrunn-* extension keywords.
// $defs is explicitly prohibited in FEATURE-0012. Shared schemas use
// approved relative $ref values targeting api/schemas/_common only.
// The schema walker is context-aware: property names under properties are
// identifiers, not keywords; extension-object fields are validated by the
// registered extension contract; unsupported actual schema keywords fail
// closed; no unsupported constraint is silently ignored.
var SupportedKeywords map[string]struct{}

// ValidateSchemaSupport scans a canonical schema and returns issues for
// ANY keyword outside SupportedKeywords. It is FAIL-CLOSED: an unsupported
// keyword is rejected, never ignored, so no constraint is silently unenforced
// (D-01a, F12-NAMING-005).
func ValidateSchemaSupport(schema []byte) []SchemaIssue

// ValidateInstance structurally validates a decoded instance against a
// canonical schema using only the supported subset. Callers MUST first pass
// ValidateSchemaSupport. This is layer 4 (structural) of the pipeline.
func ValidateInstance(schema []byte, instance any) []SchemaIssue

// VerifyBaselineIntegrity checks each api/schemas/baseline/* file against the
// digests in BASELINE_MANIFEST.json; a mismatch fails the schema-diff gate so
// a silent baseline edit is detected. This is an INTEGRITY check only — the
// manifest is not an independently unforgeable approval, since a committer can
// change a baseline and its digest together (D-11).
func VerifyBaselineIntegrity(manifestPath, baselineDir string) error

// VerifyBaselineApproval enforces that any baseline change is APPROVED, not
// merely digest-consistent. When a baseline file's digest differs from the
// prior recorded digest, it requires matching approval evidence in
// BASELINE_APPROVALS.json containing the exact old/new digests, the approving
// ADH or approval token, the reviewer, and the date; a change without matching
// evidence fails the gate. Changing the baseline and its manifest in one commit
// is NOT sufficient. The human governance boundary is protected review /
// CODEOWNERS on the baseline and its approval record (D-11, F12-EVOLVE-002).
func VerifyBaselineApproval(approvalsPath, manifestPath, baselineDir string) error
```

### Canonical-schema-to-Go consistency (internal/apischema, D-01b)

Derivative Go-type consistency with the canonical schemas is made
**executable** rather than asserted by fixture round-tripping alone. A
`TypeBinding` registry maps each canonical schema to the derivative Go type
that represents it, and a reflection-based check verifies the type against
the schema for the supported subset.

```go
// TypeBinding maps one canonical schema to its derivative Go type. Bindings
// are registered for the eight canonical contracts and the _common
// sub-schemas so every derivative type is checked (D-01, D-01b).
type TypeBinding struct {
    SchemaPath string       // e.g. api/schemas/project.json
    GoType     reflect.Type // e.g. reflect.TypeOf(Project{})
}

// TypeBindings is the registry of all schema-to-Go bindings checked by the
// feature gate. Adding a canonical schema without a binding fails the gate.
var TypeBindings []TypeBinding

// VerifyGoTypeAgainstSchema verifies, via reflection, that a derivative Go
// type matches its canonical schema for the SUPPORTED subset. It checks at
// least: property names, JSON tags (lowerCamelCase, omitempty), required vs
// optional fields (required -> no omitempty / non-pointer; optional ->
// omitempty or pointer), primitive types (string/number/integer/boolean),
// arrays (slices) and maps (objects with additionalProperties), embedded
// (promoted) fields, enum-backed named types (named string/int types whose
// schema declares enum), and additionalProperties behavior (open maps vs
// closed structs). Any mismatch is returned as a stable-coded violation and
// fails the feature gate. Callers MUST first pass ValidateSchemaSupport so
// only supported-subset schemas are checked.
//
// Fixture round-tripping remains useful supporting evidence of consistency
// but is explicitly NOT treated as complete proof that the Go type matches
// the schema; VerifyGoTypeAgainstSchema is the authoritative consistency
// check (D-01b, F12-NAMING-005).
func VerifyGoTypeAgainstSchema(schema []byte, goType reflect.Type) []SchemaIssue
```

### Concurrency, references, conditions, schema

```go
// CheckIfMatch maps a stale resourceVersion to HTTP 412 (F12-UPDATE-002).
// Returns nil when If-Match matches or is absent for an unprotected update.
func CheckIfMatch(ifMatch, currentResourceVersion string) *apiproblem.Problem

// ValidateRef enforces allowed kinds/scopes/direction and name/uid agreement.
func (c Constraint) ValidateRef(ref apiref.TypedRef, path string) []apiproblem.Violation

// SetCondition upserts a condition and only advances LastTransitionTime on a
// status change (F12-STATUS-003).
func SetCondition(conds []apicond.Condition, cond apicond.Condition, now time.Time) []apicond.Condition

// Annotations reads and validates the five registered x-sovrunn-* schema
// keywords (x-sovrunn-profile, x-sovrunn-boundary, x-sovrunn-allowed-scopes,
// x-sovrunn-stability, x-sovrunn-field-policy). Unknown x-sovrunn-* extensions
// fail closed. No wildcard x-sovrunn-* namespace is allowed.
func ReadAnnotations(schema []byte) (SchemaMeta, []SchemaIssue)

// ClassifyChange returns Compatible | Breaking | ReviewRequired for a diff.
func ClassifyChange(oldSchema, newSchema []byte) []Change

// ValidateRoute enforces /apis/<group>/<version>/<plural-kebab> (F12-NAMING-004).
func ValidateRoute(path string) error
```

## Validation

Validation is deterministic, explicit, and testable without starting the
server (go-coding-guardrails §12). The ordered layers (F12-VALIDATION-001)
and their behavior:

1. **HTTP/content/size** — `http.MaxBytesReader(w, r.Body, Limits.MaxObjectBytes)`;
   accept exactly `application/json` and `application/yaml` (with
   `application/x-yaml` and `text/yaml` treated as `application/yaml`), all
   other media types → 415. An oversized body maps to **one deterministic
   status: 400 `REQUEST_TOO_LARGE`** (the "malformed request" class of the
   approved F12-ERROR-002 baseline). The approved baseline does not define
   413; introducing a distinct 413 would be an architecture clarification and
   is recorded as a founder-approval item (see Error Handling), not applied
   silently.
2. **Safe decode** — pure `DecodeJSON`/`DecodeYAML` (selected by media type)
   with a single body read. JSON uses `json.Unmarshal` into
   `map[string]json.RawMessage` for pre-checks then `json.Decoder` with
   `DisallowUnknownFields`. YAML is treated as a **strict JSON-compatible
   input representation only** (D-03a): using `gopkg.in/yaml.v3` for safe
   syntax-tree parsing and normalization only, it requires exactly one YAML
   document, requires string mapping keys, rejects aliases, anchors, merge
   keys (`<<`), custom/explicit tags, non-finite numbers (`.nan`/`.inf`),
   and YAML-only timestamp/binary coercions, runs an explicit `yaml.Node`
   safety and duplicate-key pass, normalizes the accepted node to a
   JSON-compatible value, marshals that value to JSON bytes, then passes
   those bytes through the same `DecodeJSON` path (unknown-field rejection,
   FieldPolicy enforcement, same destination Go type, same stable error-code
   and JSON Pointer mapping). No direct yaml.v3 typed decoding is performed
   after normalization. YAML struct tags are not required for validation
   equivalence because the normalized representation is decoded through the
   JSON path. Because the final typed decoding is identical, **JSON and YAML
   representations of the same object MUST produce equivalent validation
   results**. This extends the existing `internal/api/decode.go` pattern and
   adds the YAML path required by F12-VALIDATION-001(2).
3. **Field hygiene (operation-aware)** — reject unknown fields and duplicate
   keys with `UNKNOWN_FIELD` / `DUPLICATE_FIELD` codes. Acceptance of
   `status` and system-owned fields is governed by the `FieldPolicy` for the
   `DecodeMode` (D-15): customer mutation modes (`ModeCreateRequest`,
   `ModeReplaceRequest`) reject unauthorized `status`/system-owned fields;
   `ModeStatusUpdate`, `ModeInternalObject`, `ModeReadRepresentation`, and
   fixture decoding accept them under the Matrix C2 ownership rules. There is
   no unconditional status/system-field rejection (F12-VALIDATION-002,
   F12-META-002, F12-OWNER-002).
4. **Structural** — validate the decoded instance against its canonical JSON
   Schema 2020-12 document using the **bounded-subset validator** (via the
   `StructuralValidator` interface; `apiconform` delegates to
   `apischema.ValidateInstance`, D-01a): required shape for the object's
   profile (Matrix A invariants), finite structural limits (nesting, counts).
   Before any schema is used, `apischema.ValidateSchemaSupport` rejects any
   keyword outside the supported subset (fail-closed), so structural
   validation genuinely enforces the canonical schema rather than relying on
   fixture round-tripping alone (F12-NAMING-005, F12-VALIDATION-001(4)).
   If `Input.Validator` is nil or returns a non-nil error, the pipeline
   stops at `LayerStructural`, sets `Result.Problem` to a 500
   `INTERNAL_ERROR` Problem, sets `Result.Err` to the internal cause,
   returns no success result, and MUST NOT execute layers 5 through 7.
5. **Defaulting** — apply documented, versioned, deterministic defaults
   only (F12-VALIDATION-003); defaulting never depends on external state.
   This layer also **normalizes scope to the canonical form** via
   `apimeta.NormalizeScope`: an explicit `Kind == "Platform"` `scopeRef`
   input alternate is mapped to the canonical absent (nil) form before
   identity, authorization, concurrency, persistence, and output processing
   (D-16, F12-SCOPE-002).
6. **Semantic** — naming rules (F12-NAMING-002), enum membership, scope-kind
   validity (Matrix B), condition grammar, phase/condition consistency where
   both are present (F12-STATUS-005).
7. **Reference/kind/scope (STRUCTURAL only)** — `Constraint.ValidateRef`:
   allowed kinds and scopes, direction, name/uid agreement, and scopeRef
   well-formedness. Scope has already been normalized to the canonical
   absent (nil) platform form in layer 5, so this layer validates a single
   canonical representation and applies the platform-scope identity sentinel
   (`apimeta.PlatformScopeUID`) for the uniqueness tuple
   (F12-REF-002/003, F12-SCOPE-002). For Operation, this layer validates
   scope syntax, allowed-scope membership (six values), targetRef shape,
   and UID requirements; it does NOT resolve the target from external state
   and does NOT perform target-scope equality comparison. This layer performs
   no caller-specific authorization and makes no cross-tenant access decision.
8. **Authorization / no-existence-disclosure (adopter-owned interface)** —
   caller-specific cross-tenant and cross-organization decisions run here via
   the adopter-supplied `ScopeAuthorizer`; denials use the uniform
   `SafeDenial` mapping so inaccessible targets are indistinguishable from
   absent ones (F12-SCOPE-002, F12-SEC-004). For Operation, layer-8
   behavior is governed by the layer-8 configuration matrix (see Input):
   when OperationScope is non-nil, exactly one of Path A (ScopeAuthorizer +
   pre-derived TargetScope) or Path B (AuthorizedTargetScopeResolver) MUST
   be completely configured; an invalid configuration stops at
   LayerAuthorization with Result.Problem = 500 INTERNAL_ERROR. A valid
   path resolves the target scope and invokes
   `CheckOperationTargetScopeMatch`; unauthorized or unavailable targets use
   SafeDenial and MUST NOT disclose a scope mismatch. FEATURE-0012 provides
   the interface and mapping only and implements no policy engine
   (F12-IMPL-002). Generic non-Operation validation (OperationScope is nil
   and no layer-8 capability was requested) MAY stop successfully after
   layer 7.
9. **Later-feature capability/runtime** — RESERVED; declared so adopting
   features slot in without reordering.

Rules enforced across layers:

- Unknown fields and duplicate keys fail by default (F12-VALIDATION-002).
- Structural, semantic, reference, authorization, and policy failures remain
  distinguishable via `Result.FailedAt` and per-violation codes
  (F12-VALIDATION-004).
- Validation is safe offline; layers 1–7 require no external state
  (F12-VALIDATION-005).
- Every failure carries a stable code and a JSON Pointer path; messages
  redact unauthorized/sensitive detail (F12-VALIDATION-006).
- All limits are finite and configurable (F12-VALIDATION-007).

Profile selection (F12-PROFILE-001) is validated against the schema's
`x-sovrunn-profile`; an object that fits no approved profile is rejected and
surfaces reassessment trigger F12-TRIGGER-010 rather than being force-fit.

### Canonical schema validation approach (Reuse / Wrap / Extend / Build)

F12-NAMING-005 and F12-VALIDATION-001(4) require that the canonical JSON
Schema 2020-12 document is the enforced source of truth and that structural
validation is executable. Fixture round-tripping and `x-sovrunn-*` annotation
checks alone do not satisfy this, so exactly one enforceable approach is
selected. The three candidate approaches were assessed:

| Approach | Assessment | Decision |
|---|---|---|
| **Reuse** a mature JSON Schema 2020-12 implementation | Fullest coverage, but every mature Go implementation is an external dependency, contradicting the hard constraint that only the standard library and the already-present `gopkg.in/yaml.v3` are used (D-14); adopting one requires an approved dependency decision. | Not selected now (would need founder-approved dependency decision). |
| **Build** generated validators/types from the canonical schemas | Removes runtime schema interpretation, but requires a code-generation toolchain (an added build dependency) and generated artifacts that must themselves be consistency-checked. | Not selected now. |
| **Build** an explicitly bounded supported-subset validator (stdlib) that fail-closed rejects unsupported keywords | Deterministic, offline-capable, provider-neutral, zero new dependency, and honest: the supported subset is explicit and any out-of-subset keyword is executably rejected so nothing is silently unenforced. Bounded scope avoids building a partial generic engine by stealth. | **Selected (D-01a).** |

Selected approach: **Build — bounded supported subset with fail-closed
rejection of unsupported keywords** (`apischema.ValidateSchemaSupport` +
`apischema.ValidateInstance`). This is a narrow, declared Build that satisfies
the reuse-first guardrail because the meta-model still *extends* JSON Schema
2020-12 as the canonical format; only a bounded, explicit validation subset is
implemented in-house, and a full generic engine remains deferred behind an
approved dependency decision. The eight canonical schemas are authored within
this subset; if a future contract needs a keyword outside the subset, the
schema is rejected until either the subset is extended by approved change or a
mature implementation is adopted by approved dependency decision.

## API and handler design

FEATURE-0012 adds no runtime routes and creates no domain handlers
(F12-IMPL-001/002). It defines the API grammar that adopting features use
and provides reusable helpers that a handler composes:

- **Route form** — new adopting APIs use
  `/apis/<group>/<version>/<plural-kebab>`, with scoped collections nesting
  the parent scope, e.g.
  `/apis/core.sovrunn.io/v1alpha1/tenants/{tenant}/projects` and
  `/apis/fabric.sovrunn.io/v1alpha1/resource-pools`. `ValidateRoute`
  enforces the form; unversioned public endpoints are rejected
  (F12-NAMING-004). Phase 1 routes (`/organizations`, `/tenants`, etc.) are
  retained unchanged until a separately approved migration
  (F12-COMPAT-003).
- **Normative operations** — create, get, list, full replace, delete
  (F12-UPDATE-001). PATCH, watch/change-stream, and status-update paths are
  reserved: documented and mapped to a stable "not implemented" contract,
  with no route registered (F12-UPDATE-003, design question 11).
- **Handler skeleton (reference, for adopters)** — decode → validate →
  (authorize, adopter) → execute → record operation → record audit →
  respond, mirroring go-coding-guardrails §8/§19. FEATURE-0012 supplies the
  decode, validate, concurrency, and problem-response helpers only.
- **Concurrency** — full replace uses opaque `resourceVersion` via ETag;
  protected updates require `If-Match`; `CheckIfMatch` maps stale writes to
  412 without overwriting state (F12-UPDATE-002).
- **Lists** — responses use `ListEnvelope` with opaque `page.nextPageToken`,
  deterministic ordering, and bounded page size; tokens never expose
  offsets, provider details, or authorization context (F12-LIST-001/002/003).

Because this feature registers no routes, the hard constraint that
`internal/api` MUST NOT import `internal/server` is trivially preserved; the
new primitive packages sit below `internal/api` and import neither `api`
nor `server`.

## Registry and storage design

FEATURE-0012 introduces no persistence and no registry (F12-IMPL-002). It
generalizes the Phase 1 in-memory, metadata/spec/status, stable-error-code
model into shared contracts without changing storage
(F12-COMPAT-005). Storage-related concepts are handled as contracts, not
implementations:

- `resourceVersion` is defined as an opaque string; its concrete generation
  and the production storage/indexing implementation are deferred
  (design question 5; architecture §13).
- `uid` generation is provided as a stateless `crypto/rand` helper producing
  an opaque 128-bit **collision-resistant** value (go-coding-guardrails §22).
  Collision resistance is not a collision-proof guarantee: the probability is
  astronomically low but non-zero, so adopting storage MUST perform a
  uniqueness/collision check when persisting a new object and reject or
  regenerate on a detected collision. The algorithm remains an opaque
  implementation detail (F12-META-004).
- The grammar keeps storage replaceable and horizontally scalable: no
  package holds process-global mutable state, and adopters bind their own
  registry behind their own interfaces.

## Operation and audit behavior

FEATURE-0012 does not create an operation or audit service (F12-IMPL-002);
Operation and AuditEvent domain payloads belong to FEATURE-0013
(F12-STATUS-004). This feature instead guarantees the grammar can carry the
correlation data those features require:

- The `Operation` fixture uses the `LongRunningOperation` profile and can
  represent target, action, requester, idempotency/correlation, progress,
  retryability, and terminal result — without executing anything
  (F12-PROFILE-002, Matrix D). The Operation schema declares exactly six
  allowed scopes (Platform, Organization, OrganizationUnit, Tenant, Project,
  Provider) per D-17. Operation.scopeRef MUST equal the resolved canonical
  governance scope of Operation.targetRef; for a platform-scoped target,
  scopeRef is canonically nil; a target/scope mismatch is rejected with a
  stable validation code and JSON Pointer path.
- The `AuditEvent` fixture uses the `ImmutableRecord` profile (append-only,
  linked corrections) and can carry actor, request, operation, subject, and
  version references, so auditing links these without storing history in
  current resource status (F12-SEC-007, F12-STATUS-004).
- Conditions are current facts only; historical decisions/actions are
  explicitly out of status and routed to FEATURE-0013 records.
- `requestId` is a first-class Problem extension, and `ObjectMeta` carries
  identity/version fields, so adopters can correlate actor → request →
  operation → subject → version → decision.

The design preserves the Phase 1 audit posture (structured, secret-free,
non-blocking) without adding new audit emission in this feature.

## Correctness Properties

A property is a characteristic or behavior that should hold true across all
valid executions of the grammar — a formal statement about what the shared
primitives, decoders, validators, and conformance checks must do. These
properties identify the key invariants that tasks and property-based tests
MUST preserve. Each property is universally quantified and cites the
requirement(s) and design decision(s) it validates. Property tests run a
minimum of 100 iterations and are tagged
`Feature: api-resource-naming-status-and-validation-standard, Property N: ...`.

### Property 1: Fail-closed schema support

*For any* canonical schema document, if the schema contains any keyword
outside the declared supported subset, then `ValidateSchemaSupport` rejects
it with a stable code; no schema containing an unsupported keyword is ever
accepted or partially enforced. Equivalently, a schema passes support
validation *if and only if* every keyword it uses is in the supported subset.

Traceability: F12-NAMING-005, F12-VALIDATION-001(4); D-01a.

**Validates: Requirements 4.1, 4.9**

### Property 2: JSON/YAML validation equivalence

*For any* object expressible in the strict JSON-compatible subset, decoding
and validating its JSON form and its strict-YAML form produce equivalent
results (same accept/reject outcome and, on success, the same normalized
representation). *For any* YAML input using an alias, anchor, merge key,
custom tag, multiple documents, a non-string mapping key, a non-finite
number, or a YAML-only timestamp/binary coercion, decoding is rejected.

Traceability: F12-VALIDATION-001(2), F12-VALIDATION-002; D-03a.

**Validates: Requirements 4.9**

### Property 3: Operation-aware field ownership

*For any* object carrying `status` or system-owned fields, decoding under a
customer mutation mode (`ModeCreateRequest`, `ModeReplaceRequest`) is
rejected, while decoding the same object under `ModeStatusUpdate`,
`ModeInternalObject`, or `ModeReadRepresentation` is accepted under the
Matrix C2 ownership rules. Field acceptance is a deterministic function of
the `DecodeMode` alone.

Traceability: F12-VALIDATION-002, F12-META-002, F12-OWNER-002; D-15.

**Validates: Requirements 4.3, 4.7, 4.9**

### Property 4: Canonical platform scope

*For any* object whose schema allows `Platform`, the object with an absent
`scopeRef` and the same object with an explicit `Kind == "Platform"`
`scopeRef` normalize (via `NormalizeScope`) to the identical canonical form,
and thereafter produce identical identity tuples (using
`PlatformScopeUID`), validation outcomes, and emitted output. *For any*
object whose schema does not allow `Platform`, a nil/normalized platform
scope is rejected.

Traceability: F12-SCOPE-002, F12-REF-001; D-16.

**Validates: Requirements 4.4, 4.5**

### Property 5: Safe-denial path and response equivalence

*For any* denied cross-scope access, the "exists but inaccessible" outcome
and the "absent" outcome produce byte-identical `SafeDenial` Problem
responses (404 `RESOURCE_NOT_FOUND`) and take the same control-flow path
through the adopter contract (authorize-before-lookup, or a combined
`AuthorizedResolver` whose "not found" and "found but unauthorized" cases are
indistinguishable), with no existence-dependent fast path. This is a
path/response-equivalence property, not a constant-time guarantee.

Traceability: F12-SEC-004, F12-SCOPE-002; D-04.

**Validates: Requirements 4.4, 7.4**

### Property 6: Controlled baseline updates

*For any* change to a stored baseline schema, the schema-diff gate fails
unless the change is accompanied by recorded approval evidence whose old and
new digests match the actual baseline digests and that carries an approving
ADH/token, reviewer, and date. Co-editing a baseline file and its manifest
digest in the same commit, without matching approval evidence, is never
sufficient to pass the gate.

Traceability: F12-EVOLVE-002, F12-VERIFY-001(10); D-11.

**Validates: Requirements 4.14, 4.16**

### Property 7: Provider neutrality of core

*For any* core or customer-facing schema and any core primitive
(`apimeta`/`apiref`/`apicond`/`apiproblem`/`apivalid`), no provider-native
identifier, provider SDK type, or provider-specific field is present;
provider-native data is expressible only behind `adapter-facing`/
`plugin-facing` schemas. The provider-neutrality fitness check fails on any
violation.

Traceability: F12-SCOPE-001, F12-SEC-006, F12-BOUNDARY-001, F12-REF-003.

**Validates: Requirements 4.4, 4.5, 4.6, 7.6**

### Property 8: Condition transition semantics

*For any* sequence of condition upserts, `SetCondition` advances
`lastTransitionTime` if and only if the condition's `status` changed, and
leaves other conditions unchanged; `status` is always one of
`True`/`False`/`Unknown`, and `type`/`reason` remain stable PascalCase
identifiers. Conditions represent current facts only and never accumulate
event history.

Traceability: F12-STATUS-002, F12-STATUS-003; D-04.

**Validates: Requirements 4.8**

### Property 9: Derivative Go-type / schema consistency

*For any* registered `TypeBinding`, `VerifyGoTypeAgainstSchema` accepts the
derivative Go type *if and only if* it matches its canonical schema across
the supported subset (property names, JSON tags, required vs optional,
primitives, arrays/maps, embedded fields, enum-backed types, and
`additionalProperties` behavior); any deliberate mismatch is rejected.
Fixture round-tripping is supporting evidence only, not proof.

Traceability: F12-NAMING-005, F12-VALIDATION-001(4), F12-VERIFY-001(13); D-01b.

**Validates: Requirements 4.1, 4.9, 4.16**

### Property 10: Concurrency staleness

*For any* pair of resource versions, `CheckIfMatch` returns a 412
`STALE_RESOURCE_VERSION` Problem exactly when a supplied `If-Match` does not
equal the current `resourceVersion`, and nil (allow) when they match or when
no protection is required; a stale write never overwrites current state.

Traceability: F12-UPDATE-002; D-10.

**Validates: Requirements 4.11**

### Property 11: Operation target-scope equality

*For any* Operation, its `scopeRef` MUST equal the resolved canonical
governance scope of its `targetRef`. For a platform-scoped target,
Operation.scopeRef is canonically nil (ScopeIdentity{Platform, PlatformScopeUID}).
For a non-platform target, Operation.scopeRef identifies the target's
governance scope by UID. An Operation whose scope differs from its target's
resolved scope is rejected with violation code `OPERATION_TARGET_SCOPE_MISMATCH`
and JSON Pointer `/metadata/scopeRef`. Unauthorized or unavailable targets
use SafeDenial and MUST NOT disclose a scope mismatch. The six-value
allowed-scope list (Platform, Organization, OrganizationUnit, Tenant, Project,
Provider) does not grant authorization. Layer 7 validates scope syntax,
allowed-scope membership, targetRef shape, and UID requirements offline;
target-scope resolution and equality comparison execute at layer 8 via an
adopter-owned authorized target-scope resolver.

Generated positive and negative tests for Property 11:

- each of the six scopes with matching target (positive);
- platform nil scope with a platform-scoped target (positive);
- matching non-platform UID (positive);
- kind mismatch (negative: OPERATION_TARGET_SCOPE_MISMATCH, /metadata/scopeRef);
- UID mismatch (negative: OPERATION_TARGET_SCOPE_MISMATCH, /metadata/scopeRef);
- unavailable target (negative: SafeDenial 404, no mismatch disclosed);
- unauthorized target (negative: SafeDenial 404, no mismatch disclosed);
- path A: complete configuration (OperationScope, TargetRef, TargetScope,
  Authorizer, Caller) with ScopeAuthorizer Allow then match (positive);
- path A: pre-derived TargetScope with ScopeAuthorizer DenyNotDisclosed (SafeDenial);
- path B: complete configuration (OperationScope, TargetRef,
  TargetScopeResolver, Caller) with available=true then match (positive);
- path B: AuthorizedTargetScopeResolver available=false (SafeDenial, no mismatch);
- OperationScope non-nil with neither TargetScope nor TargetScopeResolver
  configured (500 INTERNAL_ERROR at LayerAuthorization, no target lookup);
- both TargetScope and TargetScopeResolver configured (500 INTERNAL_ERROR
  at LayerAuthorization, no target lookup);
- missing Caller when OperationScope is non-nil (500 INTERNAL_ERROR at
  LayerAuthorization);
- missing TargetRef when OperationScope is non-nil (500 INTERNAL_ERROR at
  LayerAuthorization);
- incomplete Path A: TargetScope set but Authorizer nil (500 INTERNAL_ERROR
  at LayerAuthorization);
- incomplete Path B: TargetScopeResolver set but Caller nil (500
  INTERNAL_ERROR at LayerAuthorization);
- generic non-Operation validation (OperationScope nil) stops successfully
  after layer 7 (positive);
- exact OPERATION_TARGET_SCOPE_MISMATCH code assertion;
- exact /metadata/scopeRef path assertion.

Traceability: F12-SCOPE-002, F12-REF-001, F12-FIXTURE-001; D-17.

**Validates: Requirements 4.4, 4.5**

## Error Handling

Errors use RFC 9457 Problem Details with Sovrunn extensions
(F12-ERROR-001). `apiproblem.httpmap` provides the single baseline mapping
(F12-ERROR-002):

| Failure class | HTTP | Example `code` |
|---|---:|---|
| Malformed request | 400 | `MALFORMED_REQUEST`, `UNKNOWN_FIELD`, `DUPLICATE_FIELD`, `REQUEST_TOO_LARGE` |
| Authentication required | 401 | `AUTH_REQUIRED` (reserved for adopters) |
| Authorization denied | 403 | `AUTHORIZATION_DENIED` (reserved for adopters) |
| Not found | 404 | `RESOURCE_NOT_FOUND` |
| Lifecycle/uniqueness conflict | 409 | `CONFLICT`, `ALREADY_EXISTS`, `DELETE_BLOCKED` |
| Stale resource version | 412 | `STALE_RESOURCE_VERSION` |
| Unsupported media type | 415 | `UNSUPPORTED_MEDIA_TYPE` |
| Structurally/semantically invalid | 422 | `VALIDATION_FAILED` (+ `violations[]`) |
| Internal failure | 500 | `INTERNAL_ERROR` |
| Temporary dependency failure | 503 | `DEPENDENCY_UNAVAILABLE` |

Rules:

- Problem `type`, `code`, and violation codes are stable machine contracts;
  human `message`/`detail` MAY evolve and MUST NOT be parsed by clients
  (F12-ERROR-003).
- Responses MUST NOT expose credentials, stack traces, raw provider errors,
  sensitive policy inputs, or inaccessible resource details (F12-ERROR-004).
- An **oversized request body** maps to exactly one deterministic status:
  400 `REQUEST_TOO_LARGE` (the malformed-request class of the approved
  F12-ERROR-002 baseline). The earlier ambiguous "413/400" wording is
  removed. **Founder-approval item:** the approved baseline (F12-ERROR-002)
  does not define 413; if a distinct 413 "Payload Too Large" is later
  desired, that is an architecture clarification requiring founder approval
  and a requirements update, not a silent design change.
- Inaccessible objects are indistinguishable from absent ones in the
  response — the layer-8 `SafeDenial` mapping routes every `DenyNotDisclosed`
  outcome to an identical 404 `RESOURCE_NOT_FOUND` (same code, title, and
  message), so no existence disclosure occurs via error text or code.
  Response mapping alone is not sufficient: the adopter contract
  (authorize-before-lookup or combined `AuthorizedResolver`, no
  existence-dependent fast paths) supplies the equivalent control flow, and
  the path/response-equivalence tests verify it. This is a path/response
  equivalence guarantee, not a perfect constant-time guarantee. In-scope
  authorization denials where existence is already known to the caller use
  403 `AUTHORIZATION_DENIED` (F12-SEC-004, F12-SCOPE-002).
- This generalizes, and stays consistent with, the Phase 1
  `resources.ErrorCode` registry and `{"error": ...}` envelope; the Phase 1
  envelope is recorded as a compatibility exception with a migration
  candidate to Problem Details in the Phase 1 compatibility report
  (F12-COMPAT-002/005).

## Security and privacy

The shared grammar is the enforcement point for the security requirements
(F12-SEC-001..009):

- **Field classification** — every schema field crossing a boundary
  declares or inherits data classification, authorized writer/readers,
  mutability, retention, redaction, residency implication, and audit
  requirement. `x-sovrunn-*` annotations plus per-field classification carry
  this; a fitness function asserts presence (F12-SEC-001).
- **Closed classification set** — exactly the seven `DataClassification`
  values; any other value is rejected (F12-SEC-002).
- **Least privilege / no raw secrets** — strict decoding, typed
  scope/reference rules, and bounded inputs are mandatory; secret values are
  representable only through typed secret references, never in metadata,
  labels, annotations, status, errors, or audit messages. A fitness function
  scans schemas and fixtures for secret-like patterns and fails on any
  (F12-SEC-003, F12-EXT-002, Matrix E F12-R10).
- **No existence disclosure** — caller-specific cross-scope denial is a
  layer-8 decision (adopter-owned `ScopeAuthorizer`); the uniform
  `SafeDenial` mapping routes `DenyNotDisclosed` to an identical 404 so error
  text and codes never reveal whether an inaccessible target exists
  (F12-SEC-004). Response mapping alone is insufficient: adopters MUST also
  make control flow equivalent by authorizing before lookup when scope is
  derivable, or using the combined `AuthorizedResolver` so "not found" and
  "found but unauthorized" return one uniform unavailable outcome; adopters
  MUST NOT introduce existence-dependent fast paths. This is a path- and
  response-equivalence guarantee, not a perfect constant-time guarantee.
  Layer 7 remains structural only and makes no access decision.
- **Provenance and freshness** — `ObservedExternalResource` fixtures require
  source, observation time, and freshness; a disconnected dependency
  produces explicit stale/unknown/absent status, not fabricated accuracy
  (F12-SEC-005, F12-STATUS-003).
- **Provider replaceability** — provider choice, location, and ownership
  scope are separable; provider replacement requires no customer-contract
  change; provider-native identifiers are prohibited from customer/core
  contracts (F12-SEC-006, F12-SCOPE-001, F12-REF-003).
- **Auditability without status history** — identity, version, and
  reference fields support actor/request/operation/subject/version linking
  for FEATURE-0013 (F12-SEC-007).
- **Governed AI** — AI is only ever a consumer of an authorized,
  boundary-filtered view relying on stable structured codes; it receives no
  bypass boundary and must not scrape human messages. The grammar exposes
  stable `code`/`reason` values precisely so explanations need not parse
  prose (F12-SEC-008, F12-BOUNDARY-002).
- **Sovereignty / offline** — schemas are open and locally processable;
  structural validation and contract interpretation require no external
  availability, supporting on-premise, disconnected, and air-gapped
  deployments (F12-SEC-009, F12-VALIDATION-005).

Boundary separation (F12-BOUNDARY-003) is enforced by preferring separate
boundary-specific schema views over hidden-field redaction; a single schema
serving audiences with differing trust/sensitivity fails the boundary
fitness function.

## Testing Strategy

Every primitive has unit tests; the grammar has conformance tests. Tests are
deterministic and require no external services (go-coding-guardrails §30).

### Unit tests (per package)

```text
internal/apimeta   — type-identity parsing; naming rules; scope-kind validity;
                     ScopeRef typed-ref conformance (apiVersion/kind/name/uid);
                     canonical platform scope: nil is the stored/emitted form,
                     accepted only when allowed and rejected otherwise, and
                     NormalizeScope maps an explicit Platform scopeRef to nil;
                     PlatformScopeUID sentinel used in the identity tuple;
                     ownerRef-as-scope rejection; embedded TypeMeta/TypedRef JSON
                     promotion (no json:",inline"); UID opacity + collision-resistant
                     generation.
internal/apiref    — allowed kind/scope/direction; name/uid agreement and mismatch;
                     provider-native-id rejection; collection Refs bounds.
internal/apicond   — status enum; PascalCase type/reason; LastTransitionTime only
                     changes on status change; conditions-not-history.
internal/apiproblem— Problem shape; stable codes; violation JSON Pointer; httpmap
                     incl. oversize -> 400 REQUEST_TOO_LARGE.
internal/apivalid  — pure DecodeJSON and DecodeYAML (unknown/duplicate rejection,
                     yaml.v3 safety + normalize-to-JSON + DecodeJSON reuse,
                     duplicate-key rejection, size); YAML strict
                     JSON-compatible rules (single document, string keys, reject
                     aliases/anchors/merge keys/custom tags/non-finite numbers/YAML-only
                     timestamp+binary coercions) and JSON/YAML validation equivalence for
                     the same object (identical typed values, identical error codes and
                     JSON Pointers, identical unknown-field rejection, identical
                     operation-aware FieldPolicy, YAML-only constructs rejected before
                     JSON normalization); media-type selection incl. 415 for unsupported
                     types; DecodeMode/FieldPolicy
                     (create/replace reject status+system; status-update/internal/read/
                     fixture accept under ownership rules); nine-layer ordering with
                     layer 7 structural vs layer 8 authorization split;
                     StructuralValidator unavailability (nil validator or error) stops
                     pipeline at LayerStructural with Result.Problem=500 INTERNAL_ERROR,
                     Result.Err=internal cause, and never executes layers 5–7;
                     layer-8 configuration matrix tests (complete Path A, complete
                     Path B, OperationScope with neither path configured, both paths
                     configured, missing Caller, missing TargetRef, incomplete Path A,
                     incomplete Path B — each invalid configuration stops at
                     LayerAuthorization with 500 INTERNAL_ERROR; generic non-Operation
                     validation with nil OperationScope stops successfully after
                     layer 7);
                     ScopeAuthorizer + SafeDenial (owned here in
                     authz.go) no-existence-disclosure uniform 404; ScopeIdentity +
                     CanonicalScopeIdentity (nil -> Platform/PlatformScopeUID);
                     CheckOperationTargetScopeMatch (positive: each of six scopes,
                     platform nil, matching non-platform UID; negative: kind mismatch,
                     UID mismatch -> OPERATION_TARGET_SCOPE_MISMATCH /metadata/scopeRef;
                     unavailable target -> SafeDenial; unauthorized target -> SafeDenial);
                     AuthorizedTargetScopeResolver contract (stub: available=false for
                     absent and unauthorized both map through SafeDenial identically;
                     available=true feeds CheckOperationTargetScopeMatch; path A via
                     pre-derived TargetScope; path B via resolver);
                     limits; CheckIfMatch 412; offline safety.
internal/apischema — x-sovrunn-* presence/vocabulary; ValidateSchemaSupport
                     fail-closed rejection of unsupported keywords (returns
                     []SchemaIssue, not apiproblem.Violation); ValidateInstance
                     structural validation against the bounded subset (returns
                     []SchemaIssue); VerifyGoTypeAgainstSchema Go-type consistency
                     across the supported subset (property names, JSON tags,
                     required vs optional, primitives, arrays/maps, embedded fields,
                     enum-backed types, additionalProperties) with mismatches
                     returning SchemaIssue; ReadAnnotations returns SchemaIssue;
                     route-form; schema-diff classification (add optional/required,
                     remove/rename, narrow enum, add enum value, target kind/scope
                     change); VerifyBaselineIntegrity against BASELINE_MANIFEST.json;
                     VerifyBaselineApproval requiring recorded approval evidence
                     (old/new digests, ADH/token, reviewer, date) so co-editing
                     baseline + manifest alone is not sufficient to approve.
```

### Conformance tests (tests/conformance, internal/apiconform)

- **Positive** — each of the eight fixtures decodes (in both JSON and YAML
  form), structurally validates against its canonical schema via the
  bounded-subset validator, and matches it; profile/boundary/scope/stability
  annotations present and valid (F12-FIXTURE-002, F12-VERIFY-001(1)). Fixture
  decoding uses `ModeReadRepresentation`/`ModeInternalObject` so legitimate
  `status`/system fields are accepted; a status-update fixture decodes under
  `ModeStatusUpdate`.
- **Operation-aware fields** — the same object is accepted under internal/
  read/status-update modes but rejected under create/replace modes when it
  carries `status` or system-owned fields, and vice-versa (D-15,
  F12-META-002, F12-OWNER-002).
- **Schema-support gate** — a schema using a keyword outside the supported
  subset is rejected fail-closed by `ValidateSchemaSupport`, proving no
  constraint is silently unenforced (D-01a, F12-NAMING-005).
- **Canonical platform scope** — an object supplied with an explicit
  `Kind == "Platform"` `scopeRef` and the same object supplied with an
  absent `scopeRef` normalize (via `NormalizeScope`) to the identical
  canonical form and produce identical identity, validation, and emitted
  output; the `PlatformScopeUID` sentinel yields a well-defined identity
  tuple (D-16, F12-SCOPE-002).
- **JSON/YAML equivalence** — for each fixture, the JSON form and the
  strict JSON-compatible YAML form decode and validate to equivalent results;
  YAML aliases, anchors, merge keys, custom tags, multiple documents,
  non-string keys, and non-finite numbers are each rejected (D-03a,
  F12-VALIDATION-001(2)).
- **Negative** — unknown field, duplicate key (JSON and YAML), unauthorized
  client-authored status/system field under a mutation mode, name/uid
  mismatch, invalid scope kind, `ownerRef` supplied where a `scopeRef.kind`
  is expected, nil scopeRef on a resource that does not allow Platform,
  oversized body (→ 400 `REQUEST_TOO_LARGE`), over-nested/over-count inputs,
  unsupported media type (→ 415), and unversioned route all fail with the
  correct code and JSON Pointer.
- **Boundary** — object/metadata/condition/violation/reference/page limits
  reject at the boundary (F12-VERIFY-001(11)).
- **Security** — secret-like values in metadata/labels/annotations/status/
  extensions are rejected; a cross-tenant/cross-organization reference is
  denied at layer 8 via a stub `ScopeAuthorizer` and the `SafeDenial`
  mapping, with the "exists but inaccessible" and "absent" cases producing
  an identical 404 (no existence disclosure). A **response/path-equivalence**
  test asserts that both cases return byte-identical Problem responses AND
  take the same control-flow path through a stub `AuthorizedResolver` (same
  branch, no existence-dependent fast path, no extra side effect/log/audit);
  the authorize-before-lookup path is exercised so a cross-scope denial
  performs no lookup. These are equivalence tests, not constant-time timing
  proofs. No provider SDK/native type is embedded in core/customer schemas
  (F12-VERIFY-001(2,6,7), F12-SEC-004).
- **Compatibility** — schema-diff detects breaking changes against
  `api/schemas/baseline/*`; a tampered baseline file fails
  `VerifyBaselineIntegrity` against `BASELINE_MANIFEST.json`, and a baseline
  change with a matching (co-edited) manifest but NO recorded approval
  evidence fails `VerifyBaselineApproval`, proving a baseline change cannot be
  self-approved by editing baseline + manifest in one commit; the
  diff gate cannot be silently bypassed; `VerifyGoTypeAgainstSchema` proves
  each derivative Go type matches its canonical schema across the supported
  subset (a deliberately mismatched Go tag/type fails the check, not merely a
  fixture round-trip) (F12-VERIFY-001(10,13)); Matrix D scenarios each map
  to a fixture with the required proof (F12-FIXTURE-001).
- **Boundary ledger** — `docs/api/boundary-ledger.yaml` parses strictly and
  every declared boundary carries all F12-LEDGER-001 categories; a boundary
  present in a schema without a ledger entry fails the check (F12-LEDGER-001,
  F12-VERIFY-001).
- **Absence of runtime** — a check asserts no provider/plugin/policy/
  placement/audit/provisioning execution is present (F12-VERIFY-001(14),
  F12-IMPL-002).

### Matrix D coverage

The `internal/apiconform` scenario table maps all seventeen Matrix D
scenarios to a fixture + required-proof assertion, so the conformance suite
fails if any scenario becomes unrepresentable. The "Future provisioning
executes" scenario uses the Operation fixture with all six allowed scopes
(Platform, Organization, OrganizationUnit, Tenant, Project, Provider) and
asserts target-scope equality per D-17, including generated positive and
negative tests: each of the six scopes with matching target; platform nil
scope; matching non-platform UID; kind mismatch
(OPERATION_TARGET_SCOPE_MISMATCH, /metadata/scopeRef); UID mismatch;
unavailable target (SafeDenial, no mismatch disclosed); unauthorized target
(SafeDenial, no mismatch disclosed).

## Verification

Verification is split into automatable checks and recorded human review, per
F12-VERIFY-002 and F12-VERIFY-003. Not every acceptance criterion is
machine-verifiable.

### Automated (F12-VERIFY-002)

```bash
make fmt
make test          # unit + conformance tests, incl. -race for concurrency helpers
make vet
make ff-feature-gate FEATURE=FEATURE-0012
```

`scripts/api-conformance-check.sh` (invoked by the feature gate) runs the
fifteen fitness functions (F12-VERIFY-001), the bounded-subset schema-support
gate (fail-closed on unsupported keywords, D-01a), the schema-diff gate
against `api/schemas/baseline/*` with the `BASELINE_MANIFEST.json` integrity
check (tamper detection) plus the `VerifyBaselineApproval` check requiring
recorded approval evidence (old/new digests, approving ADH/token, reviewer,
date) so a baseline change is not approvable merely by editing the baseline
and manifest together (D-11), the strict
machine-readable boundary-ledger check (all F12-LEDGER-001 categories, D-12),
and the Phase 1 compatibility coverage check. Any failure fails the gate. The
gate also runs the FEATURE-0011 reuse assessment check in strict mode.

Fitness function → requirement mapping (F12-VERIFY-001 checks 1–15) is
enumerated in `internal/apiconform/fitness.go` and asserted by test.

### Human semantic review (F12-VERIFY-003)

The following are validated by recorded human review, not by `make` alone,
and must record reviewer, decision, and date:

- architecture approvals and any granted adoption exception
  (F12-SCOPESTD-004);
- residual-risk acceptance for Matrix E (F12-RISK-001);
- correctness of boundary classifications and responsibility boundaries;
- adequacy of the Phase 1 compatibility exceptions and migration candidates.

The feature must not claim that every acceptance criterion is verified by
`make` commands alone.

### Go version and dependencies

- Target Go 1.22 per `docs/engineering/go-version-standard.md`. Design note:
  `go.mod` currently declares `go 1.21`; tasks should align it to `1.22`
  (not a requirements contradiction, so requirements.md is unchanged).
- No external dependency is added. Only the standard library and the
  already-present `gopkg.in/yaml.v3` are used (hard constraint).

## Non-goals

Consistent with F12-IMPL-001/002 and requirements section 5, this design
does not implement:

- provider or substrate models (FEATURE-0014/0015);
- adapter protocols (FEATURE-0016); adapter-native data stays behind the
  adapter-facing boundary;
- policy evaluation (FEATURE-0017+);
- DecisionObject/AuditEvent domain payloads (FEATURE-0013);
- placement behavior (FEATURE-0023);
- plugin taxonomy/execution (FEATURE-0024 and later);
- provisioning, persistence selection, workflow execution, billing,
  failover, or autonomous AI operations;
- a wholesale rewrite of Phase 1 resources or routes (the compatibility
  report records exceptions and migration candidates only);
- vendor-specific core models or provider-native core references;
- an unrestricted arbitrary extension mechanism;
- exact HTTP route migration, PATCH format, watch/change-stream protocol,
  production storage/indexing, or stable API promotion (all deferred and
  separately approved);
- a full generic JSON Schema 2020-12 validation engine for arbitrary
  documents — FEATURE-0012 implements an explicitly bounded supported-subset
  validator (D-01a) that structurally enforces the canonical schemas with the
  standard library and fail-closed rejects any keyword outside the subset, so
  no constraint is silently unenforced. A general-purpose engine for arbitrary
  documents is deliberately out of scope and, if later required, needs an
  approved dependency decision (Reuse) or an approved subset extension.

Defining the standard as cross-phase (F12-SCOPESTD-001) does not expand this
implementation scope. Later applicable features adopt the grammar under
their own approvals (F12-SCOPESTD-002/004, F12-COMPAT-006).

## Resolved design questions

The twelve deferred questions in requirements section 9 are resolved below,
each within the approved contracts. None required an
`ARCHITECTURE_DECISION_REQUIRED` halt (F12-GOV-001).

1. **Canonical schema source of truth + derivative flow** — the JSON Schema
   2020-12 document per contract under `api/schemas/` is canonical (D-01).
   Go types, docs, examples, and SDKs are derivative. Structural conformance
   is **executable**: the bounded supported-subset validator (D-01a)
   validates instances against the canonical schema and fail-closed rejects
   any unsupported keyword, so enforcement does not rely on fixture
   round-tripping or annotation checks alone. Derivative **Go-type**
   consistency is also executable: the `TypeBinding` registry plus the
   reflection-based `VerifyGoTypeAgainstSchema` check (D-01b) verifies each
   derivative Go type against its canonical schema across the supported
   subset (property names, JSON tags, required vs optional, primitives,
   arrays/maps, embedded fields, enum-backed types, and
   `additionalProperties` behavior). Round-tripping fixtures through the Go
   types with strict decoding and asserting `x-sovrunn-*` annotations against
   controlled vocabularies remains useful supporting evidence but is NOT
   treated as complete proof of Go-type/schema consistency
   (F12-NAMING-005/006, F12-VALIDATION-001(4)). The chosen approach and its
   Reuse/Wrap/Extend/Build assessment are in the Validation section.
2. **Go package/module boundaries** — `internal/apimeta`, `internal/apiref`,
   `internal/apicond`, `internal/apiproblem`, `internal/apivalid`,
   `internal/apischema`, `internal/apiconform`, with the import direction in
   the Architecture section. No package imports `internal/api` or
   `internal/server` (D-02).
3. **Strict decoding, validation ordering, error mapping, concurrency** —
   pure `DecodeJSON` (stdlib `DisallowUnknownFields` + token-scan duplicate
   detection) and `DecodeYAML`, which treats YAML as a **strict
   JSON-compatible input representation only** (D-03a): exactly one document,
   string keys, no aliases/anchors/merge keys/custom tags/non-finite numbers/
   YAML-only timestamp or binary coercions, an explicit `yaml.Node` safety +
   duplicate-key pass, normalization to a JSON-compatible value, marshaling
   to JSON bytes, then passing those bytes through the same `DecodeJSON` path
   (unknown-field rejection, FieldPolicy enforcement, same destination Go
   type, same stable error-code and JSON Pointer mapping) — so JSON and YAML
   forms of the same object produce equivalent validation results. No direct
   yaml.v3 typed decoding is performed after normalization. Both are
   separated from the `StrictDecode` HTTP adapter which selects a decoder by
   media type (`application/json`/`application/yaml` only, else 415).
   Decoding is operation-aware via `DecodeMode`/`FieldPolicy`
   (D-15). The nine-layer `Validate` pipeline keeps layer 7 structural and
   moves caller-specific authorization + no-existence-disclosure to the
   adopter-owned layer-8 `ScopeAuthorizer` + `SafeDenial` (D-04). Error
   mapping is the single Problem/httpmap baseline, with an oversized body
   mapped deterministically to 400 `REQUEST_TOO_LARGE`; the `CheckIfMatch`
   412 helper handles concurrency (D-03, D-05, D-10).
4. **Exact finite limits** — the `apivalid.Limits` struct and default table
   in Data models, treated as reviewed platform configuration (D-06,
   F12-VALIDATION-007).

5. **UID + resourceVersion representation** — `uid` is an opaque 128-bit
   **collision-resistant** value from `crypto/rand` (not collision-proof);
   adopting storage MUST perform a uniqueness/collision check on persist and
   reject or regenerate on collision. `resourceVersion` is an opaque string
   whose concrete generation is the adopting storage's responsibility.
   Clients treat both as opaque (D-07, F12-META-004).
6. **Machine-readable annotations + fitness validation** — the Sovrunn
   extension registry contains exactly five extensions: `x-sovrunn-profile`,
   `x-sovrunn-boundary`, `x-sovrunn-allowed-scopes`, `x-sovrunn-stability`,
   and `x-sovrunn-field-policy`. These are JSON Schema extension keywords
   parsed by `apischema.ReadAnnotations` and asserted by fitness function
   check 1. No wildcard `x-sovrunn-*` namespace is allowed; unknown
   `x-sovrunn-*` extensions fail closed. `x-sovrunn-field-policy` is a
   strictly validated property-level object carrying or inheriting exactly:
   classification, authorizedWriter, authorizedReaders, mutability,
   retention, redaction, residency, auditRequired (D-08).
7. **HTTP route form + Phase 1 coexistence** —
   `/apis/<group>/<version>/<plural-kebab>` with nested scope for scoped
   collections, validated by `ValidateRoute`; Phase 1 routes retained
   unchanged until a separately approved migration. FEATURE-0012 registers
   no routes (D-09, F12-NAMING-004, F12-COMPAT-003).
8. **Schema-diff gate** — `apischema.ClassifyChange` compares current
   schemas to `api/schemas/baseline/*` and classifies each change
   (compatible / breaking / review-required) per the change-classification
   table; wired into `scripts/api-conformance-check.sh` and the feature gate.
   The baseline is **immutable except through an approval-controlled
   workflow**. `VerifyBaselineIntegrity` checks each baseline file against the
   digests in `api/schemas/baseline/BASELINE_MANIFEST.json` and detects a
   silent baseline edit, but the manifest is an **integrity mechanism, not an
   independently unforgeable approval mechanism** — editing a baseline file
   and its manifest digest in the same commit MUST NOT be sufficient to
   approve the change. `VerifyBaselineApproval` additionally requires recorded
   approval evidence (`api/schemas/baseline/BASELINE_APPROVALS.json` or
   equivalent) carrying the exact old/new digests, the approving ADH or
   approval token, the reviewer, and the date; a baseline change without
   matching approval evidence fails the gate. The human governance boundary
   is **protected review / CODEOWNERS (or equivalent)** on the baseline and
   its approval record (D-11, F12-EVOLVE-002).
9. **Boundary ledger representation + sync** — the source of truth is the
   **machine-readable** `docs/api/boundary-ledger.yaml`, governed by a strict
   internal ledger schema; `docs/api/BOUNDARY_LEDGER.md` is a derivative human
   view. A fitness function parses the ledger strictly and asserts that every
   declared boundary carries all F12-LEDGER-001 categories (purpose, owner,
   producers, consumers, allowed/prohibited data, authorization, audit,
   observability, failure behavior, versioning, replacement path, migration
   path, reassessment trigger) and that every boundary present in a schema
   has a ledger entry (D-12).
10. **Testing strategy + eight fixtures** — the Testing Strategy section
    defines positive/negative/boundary/security/compatibility suites and the
    eight fixtures with Matrix D scenario coverage (F12-FIXTURE-001/002).
11. **Reserving PATCH / watch / status-update** — reserved by documentation
    and by declared-but-unimplemented validation layers and a stable
    "not implemented" contract; no route or handler is added
    (F12-UPDATE-003).
12. **Phase 1 compatibility report capture** —
    `docs/api/PHASE1_COMPATIBILITY_REPORT.md` records, per Phase 1 contract,
    conforming behavior, explicit exceptions, and migration candidates;
    `apiconform.compat` asserts the report covers every required Phase 1
    resource and endpoint and never triggers a rewrite (D-13,
    F12-COMPAT-001/002).

## Reuse assessment (capability-level)

Field definitions are owned by
`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` and are not redefined
here. This assessment populates the canonical fields for the single
significant decision unit.

### Identity

- Feature identity: FEATURE-0012
- Capability / decision-unit identity: Sovrunn API/resource meta-model and
  conformance foundation
- Assessment owner: Sovrunn Architecture Owner

### Classification

- Disposition: Extend
- Decision status: Approved

### Analysis

- Assessment scope: shared API/resource grammar primitives, schema
  conventions, validation helpers, error contract, and conformance tooling.
- Candidate category: API/resource meta-model standards and conventions.
- Mature candidates / applicable standards: Kubernetes API conventions,
  OpenAPI 3.1, JSON Schema 2020-12, RFC 9457 Problem Details, RFC 6901 JSON
  Pointer, HTTP semantics, ETag/`If-Match` optimistic concurrency.
- Relevant candidate strengths: proven declarative desired/observed model,
  machine-readable schemas, standard error transport, well-understood
  concurrency and field-path semantics.
- Material candidate constraints: none imposes Sovrunn's sovereign scope,
  boundary, ownership, provider-neutrality, or conformance rules; using them
  unmodified would risk provider/customer coupling and Kubernetes overfit.
- Rationale: extend mature standards with Sovrunn-owned constraints rather
  than build a new meta-model or fork an engine.
- Selected foundation or approach: Extend the listed standards behind
  Sovrunn-owned grammar, primitives, and executable conformance.

### Boundary

- Sovrunn-owned responsibility: profile taxonomy; naming/identity; common
  metadata; scope/reference semantics; boundary classification; ownership
  and mutability; status/condition grammar; validation stages; error codes;
  redaction; extension governance; provider-neutrality; compatibility,
  migration, conformance, and reassessment policy.
- Reused or external responsibility: generic HTTP, OpenAPI/JSON Schema,
  Problem Details, JSON Pointer, ETag semantics, and selected Kubernetes
  conventions.
- Data crossing the boundary: schema documents and validated object
  instances only; no provider-native or secret data.
- Control crossing the boundary: none at runtime; FEATURE-0012 performs no
  external integration.
- Adapter required: No.
- Adapter rationale: FEATURE-0012 defines contract grammar and conformance
  only and performs no external system integration; there is no external
  engine to wrap. Later integrations that touch external systems must
  introduce DEC-0036-compliant adapter boundaries and are out of scope here.
- Adapter or contract identifier: none.
- Vendor-native types allowed: No.

### Suitability

- Sovereignty and deployment fit: schemas are open and locally processable;
  structural validation and contract interpretation require no external
  service, supporting on-premise, disconnected, and air-gapped deployments
  (F12-SEC-009).
- Security and trust: strict decoding, data classification, redaction,
  secret-reference-only rules, bounded inputs, and no-existence-disclosure
  denial (F12-SEC-001..004).
- Operational and supportability: standard library only; deterministic,
  testable helpers; executable feature-gate conformance.
- Licensing and supply-chain: no new dependency; reused standards are open
  and permissively usable; no third-party code introduced.
- Portability and provider-neutrality impact: provider SDK/native types are
  prohibited from core/customer contracts; provider choice, location, and
  ownership scope are separable (F12-SCOPE-001, F12-SEC-006).

### Phase and scope

- Allowed in current phase: Yes.
- Current-phase work: shared primitives, canonical schema conventions,
  strict decode/validation helpers, eight fixtures, Phase 1 compatibility
  report, and feature-gate fitness functions (F12-IMPL-001).
- Deferred work: HTTP route migration, PATCH/watch/status-update protocols,
  production storage/indexing, stable API promotion, and any real provider/
  plugin/adapter/policy integration (architecture §13).
- Explicit non-goals: see the Non-goals section (F12-IMPL-002).
- Exit or migration boundary: Phase 1 contracts coexist unchanged and
  migrate only through separately approved features; breaking changes follow
  the maturity/compatibility policy (F12-COMPAT-003/004, F12-EVOLVE-001).
- Phase 2 non-goal acknowledgement: Phase 2 remains a model, standard,
  decision, audit, adapter-boundary, and simulation foundation; no real
  runtime integration or later-phase execution is authorized here.

### Risk mitigation

- Applicable architecture risks: Matrix E, F12-R01..F12-R16.
- Preventive controls: profile matrix and provider-neutral rules; separate
  typed scope/reference/boundary concepts; one-owner-per-field/condition
  rule; strict decoding; data classification and secret-reference-only;
  registered namespaced extensions; maturity/compatibility policy; finite
  limits and opaque pagination; authorized boundary-filtered AI views.
- Detection controls: the fifteen fitness functions, schema-diff gate,
  boundary-ledger sync check, import/schema lint for provider types,
  security scan for secret-like values, and negative/security conformance
  tests.
- Corrective path: reclassify or re-scope via a new Architecture Decision
  Handoff; move leaked fields behind adapter/plugin boundaries; add a new
  API version and migration plan for breaking changes; record an approved
  adoption exception (F12-SCOPESTD-004) where a provision cannot be adopted.
- Residual risk: low; the standard is documentation plus shared primitives
  and executable conformance, with no runtime integration.
- Replacement risk: Low.
- Reassessment triggers: F12-TRIGGER-001..013 (e.g. first real provider
  adapter, first data-path plugin, stable API promotion, regulated
  workloads, an object that cannot select an approved profile, or any
  backward-incompatible migration).

### Traceability

- Related DEC / RFC / ADH references: ADH-2026-012; ADH-2026-013; RFC-0022;
  DEC-0026, DEC-0027, DEC-0036; RFC-0021 and DEC-0026/DEC-0036 via
  FEATURE-0011.
- Linked acceptance criteria: all `F12-*` requirements; the fitness-function
  → requirement map in `internal/apiconform/fitness.go`.
- Validation and review evidence: `make fmt`/`test`/`vet`, the FEATURE-0012
  feature gate, and recorded human semantic review (F12-VERIFY-002/003).

### Human-approval evidence

- Approving role: Sovrunn Architecture Owner.
- Approval date: 2026-07-22.
- Approved reference: ADH-2026-012 (Approved); architecture baseline
  approval recorded in `docs/architecture/api-resource-standard.md`.
- Scope of approval: applies to the recorded Extend disposition and the
  responsibility boundary above. This approval authorizes the design stage;
  tasks and implementation retain their separate approval gates.
- Amendment reference: ADH-2026-013 (Approved, 2026-07-23); resolves the
  canonical Operation allowed-scope enumeration and target-scope equality
  invariant as a bounded clarification.
- Amendment note: ADH-2026-013 does not constitute reapproval of the
  amended design itself; design reapproval retains its own gate.

### Nested capability assessment: Bounded JSON Schema 2020-12 structural validator

This is an **additional, nested capability-level assessment** for one decision
unit *inside* the Extend foundation above. It does NOT replace the feature-level
Extend disposition (which remains Approved for the meta-model and conformance
foundation); it records the narrower Build choice for the structural validator
component only (D-01a).

- **Capability / decision-unit identity:** Bounded JSON Schema 2020-12
  structural validator (`apischema.ValidateSchemaSupport` +
  `apischema.ValidateInstance`).
- **Disposition: Build.** (Nested within the feature-level Extend; the
  meta-model still extends JSON Schema 2020-12 as the canonical format.)
- **Decision status:** Approved (within the D-01a design decision under
  ADH-2026-012; no new architecture decision).
- **Why mature reuse was not selected now:** every mature Go JSON Schema
  2020-12 implementation is an external dependency, contradicting the hard
  constraint that only the standard library and the already-present
  `gopkg.in/yaml.v3` are used (D-14); adopting one requires a founder-approved
  dependency decision. Deferred, not rejected.
- **Why code generation was not selected now:** generating validators/types
  from the canonical schemas requires a code-generation toolchain (an added
  build dependency) and generated artifacts that themselves need
  consistency-checking; it adds tooling and supply-chain surface without a
  current need. Deferred.
- **Sovrunn-owned responsibility:** the explicit supported-keyword subset;
  fail-closed rejection of any out-of-subset keyword; structural validation of
  instances against the supported subset; stable violation codes and JSON
  Pointer paths; the bound itself (what is and is not supported).
- **Reused-standards responsibility:** the JSON Schema 2020-12 vocabulary and
  semantics for the supported keywords, and RFC 6901 JSON Pointer for paths —
  reused as specifications, not as an imported engine.
- **Adapter required:** No — the validator is an in-process pure-function
  component over decoded instances; there is no external engine to wrap.
- **Non-goals:** a full generic JSON Schema 2020-12 engine for arbitrary
  documents; remote `$ref` resolution; unsupported keywords (e.g. `allOf`/
  `anyOf`/`oneOf`/`if`/`then`/`else`, `$dynamicRef`, format assertion as
  validation); schema authoring outside the supported subset. Any of these
  remains deferred behind an approved dependency decision or an approved
  subset extension.
- **Security controls:** fail-closed by construction (an unrecognized keyword
  is rejected, never silently ignored, so no constraint is unenforced);
  offline/deterministic evaluation with no external fetch; bounded nesting and
  input size; no execution of schema-embedded code.
- **Maintenance controls:** the supported subset is a single declared list
  asserted by test; adding a keyword is a reviewed change with new tests;
  canonical schemas are authored within the subset and any out-of-subset
  keyword fails CI until the subset is extended by approved change.
- **Replacement risk: Low.** The validator sits behind
  `apischema.ValidateInstance`/`ValidateSchemaSupport`; if a mature engine is
  later approved, it can be substituted behind the same interface without
  changing callers or the canonical schemas.
- **Reassessment triggers:** a canonical contract needs a keyword outside the
  supported subset; a security defect is found in the bounded validator;
  maintenance cost of the subset becomes excessive; or a mature JSON Schema
  2020-12 dependency is approved by a founder dependency decision. Any trigger
  reopens Reuse-vs-Build for this component.

## Architecture drift checks

Each Phase 2 drift gate is addressed by this design:

- **No provider-specific hardcoding in core** — core primitives
  (`apimeta`/`apiref`/`apicond`/`apiproblem`/`apivalid`) contain no provider
  names, SDK types, or provider-native fields; fitness function check 2
  fails on any provider SDK/native import or embed in core/customer schemas
  (F12-VERIFY-001(2), F12-R04).
- **No Kubernetes-only assumptions in core** — Kubernetes conventions are
  extended, not required; the profile matrix and provider-neutral rules
  prevent Kubernetes overfit, proven by the future-scenario fixtures
  (F12-R01).
- **No PostgreSQL lifecycle logic in core placement engine** — no placement
  engine and no PostgreSQL/datastore lifecycle logic exist in this feature;
  the `PlacementEvaluationRequest` fixture is a `TransientRequestResult`
  contract only, with no evaluation behavior (F12-IMPL-002).
- **No custom policy engine embedded in handlers** — no handlers and no
  policy engine are created. The layer-8 authorization decision boundary
  (D-04) defines only a small adopter-owned `ScopeAuthorizer` interface and
  a uniform `SafeDenial` mapping; the authorization decision logic is
  supplied by adopting features, not implemented here. This drift gate
  concerns the layer-8 authorization decision (D-04), not the `ScopeRef`
  typed-reference contract (D-16) (D-04, F12-IMPL-002).
- **No raw secret storage** — secret values are representable only through
  typed secret references; a security fitness function rejects secret-like
  values in metadata/labels/annotations/status/extensions (F12-SEC-003,
  F12-R10).
- **No customer-facing IaaS leakage** — provider-native identifiers and
  internals are prohibited from customer/core contracts; boundary
  classification and check 2 enforce this; customer provider/pool selection
  uses approved domain contracts (F12-SCOPE-002, F12-BOUNDARY-001).
- **Explainable decision object** — the error/problem contract exposes
  stable `code`/`reason` values and JSON Pointer violation paths, and the
  `AuditEvent`/`Operation` fixtures carry structured references, so
  decisions and denials are explainable from structured context without
  message scraping (F12-ERROR-001/003, F12-SEC-008).
- **Defined audit behavior** — audit linkage (actor/request/operation/
  subject/version) is representable without storing history in current
  status; historical records route to FEATURE-0013 (F12-SEC-007,
  F12-STATUS-004).
- **Preserved adapter boundaries** — adapter-facing and plugin-facing
  boundaries keep provider handles and provider data behind versioned
  contracts; no external integration is performed here, and later
  integrations must use DEC-0036-compliant adapters (F12-BOUNDARY-001,
  F12-REF-003, DEC-0036).

## Traceability summary

- Architecture topics and matrices A–E: covered by the primitives, schema
  conventions, fitness functions, and fixtures described above.
- Requirement coverage: the Resolved decisions table, Validation, Error
  Handling, Security and privacy, Testing Strategy, Correctness Properties,
  and Verification sections collectively cite the `F12-*` identifiers they
  satisfy; requirements section 10's coverage matrix remains the
  authoritative index.
- No new architecture decision is introduced; ADH-2026-013 resolved the
  previously missing exact Operation scope enumeration as a bounded
  clarification (D-17). No canonical term is renamed; the feature sequence
  is unchanged; no future-phase feature is implemented.

This document is the Design stage for FEATURE-0012. Tasks and implementation
remain unauthorized until their separate human approval tokens are issued.
