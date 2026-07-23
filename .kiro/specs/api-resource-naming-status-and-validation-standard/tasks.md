# Implementation Plan: FEATURE-0012 — API, Resource Naming, Status, and Validation Standard

## Overview

This plan converts the approved design into discrete, dependency-ordered coding tasks.
Implementation language: Go 1.22 (per `docs/engineering/go-version-standard.md`).
Only the standard library and the already-present `gopkg.in/yaml.v3` are used.
No runtime routes, domain services, or later-feature execution is introduced.

Import direction (enforced, no cycles):
- `apimeta`: standard library only
- `apiref`: `apimeta` only
- `apicond`: standard library only
- `apiproblem`: standard library only
- `apivalid`: `apimeta`, `apiref`, `apicond`, `apiproblem`
- `apischema`: `apimeta` plus standard library (`apischema` MUST NOT import `apiproblem`)
- `apiconform`: may compose all grammar packages
- `apivalid` MUST NOT import `apischema`

Package boundary rules:
- `apischema` returns package-local `SchemaIssue` (not `apiproblem.Violation`)
- `apiconform` translates `SchemaIssue` → `apiproblem.Violation`
- `apivalid` uses a `StructuralValidator` interface returning `([]apiproblem.Violation, error)`
- `SafeDenial` is owned by `internal/apivalid/authz.go`, not `apiproblem`

All eleven correctness-property tests are mandatory. Each requires at least
100 generated iterations with a deterministic seed (or reports the failing
seed for reproducibility). Property tests may not be skipped for an MVP.

Controlling handoffs: ADH-2026-012 (Approved), ADH-2026-013 (Approved).

## Tasks

- [x] 1. Repository and Go 1.22 alignment
  - [x] 1.1 Update go.mod to declare `go 1.22` minimum
    - Change `go 1.21` directive to `go 1.22` in `go.mod`
    - Run `go mod tidy` to verify no breakage
    - _Design: D-14_
    - _Requirements: F12-IMPL-001, F12-VERIFY-002_
    - _Verification: `go build ./...` succeeds; `go.mod` declares 1.22_

  - [x] 1.2 Create package directories for new grammar primitives
    - Create `internal/apimeta`, `internal/apiref`, `internal/apicond`, `internal/apiproblem`, `internal/apivalid`, `internal/apischema`, `internal/apiconform`
    - Add a `doc.go` with package doc comment in each
    - _Design: D-02_
    - _Requirements: F12-IMPL-001_
    - _Verification: `go build ./internal/...` succeeds; no import cycles detected by `go vet ./...`_

  - [x] 1.3 Add provider-neutral import-direction and package-boundary enforcement test
    - Create `internal/apiconform/imports_test.go` that parses Go imports of all grammar packages and asserts:
      - `apimeta` imports only stdlib
      - `apiref` imports only `apimeta`
      - `apicond` imports only stdlib
      - `apiproblem` imports only stdlib
      - `apivalid` imports only `apimeta`, `apiref`, `apicond`, `apiproblem` (NOT `apischema`)
      - `apischema` imports only `apimeta` plus stdlib (NOT `apiproblem`)
      - `apiconform` may import all grammar packages
      - No grammar package imports `internal/api` or `internal/server`
      - No provider SDK imports in any grammar package
    - _Design: D-02; architecture provider-neutrality_
    - _Requirements: F12-IMPL-001, F12-VERIFY-001(2,14)_
    - _Verification: `go test ./internal/apiconform/...` passes_

- [ ] 2. Shared metadata, reference, condition, and problem primitives
  - [ ] 2.1 Implement TypeMeta, ObjectMeta, and profile/boundary/stability/data-classification enums (`internal/apimeta`)
    - `typemeta.go`: TypeMeta struct with `apiVersion`/`kind`; group/version parsing helpers
    - `objectmeta.go`: ObjectMeta struct with all fields per design; ownership/mutability doc comments
    - `profile.go`: Profile, Boundary, Stability, DataClassification enums with controlled vocabularies
    - _Design: D-08; data models section_
    - _Requirements: F12-NAMING-001, F12-NAMING-002, F12-META-001, F12-META-002, F12-PROFILE-001, F12-BOUNDARY-001, F12-SEC-002_
    - _Verification: `go build ./internal/apimeta` compiles; enum constants match architecture matrices A, C1_

  - [ ] 2.2 Implement TypedRef base, ScopeRef, OwnerRef, scope-kind constants, ScopeIdentity, CanonicalScopeIdentity, canonical platform scope, PlatformScopeUID, and UID generation (`internal/apimeta`)
    - `reference.go`: TypedRef struct (apiVersion/kind/name/uid) stdlib-only
    - `scope.go`: ScopeRef (embeds TypedRef), OwnerRef (embeds TypedRef), ScopeKind consts for exactly six values: Platform, Organization, OrganizationUnit, Tenant, Project, Provider
    - `scope.go`: PlatformScopeUID constant ("platform"); NormalizeScope function (explicit Platform → nil canonical form)
    - `scope.go`: ScopeIdentity struct (Kind ScopeKind, UID string) for authorization comparison without nil ambiguity
    - `scope.go`: CanonicalScopeIdentity(*ScopeRef) → ScopeIdentity: nil → {Platform, PlatformScopeUID}; non-platform → {ref.Kind, ref.UID}
    - `uid.go`: crypto/rand opaque 128-bit collision-resistant UID generation; document adopter persistence collision-check contract
    - _Design: D-16, D-07, D-17_
    - _Requirements: F12-REF-001, F12-SCOPE-002, F12-OWNER-001, F12-META-004_
    - _Verification: `go test ./internal/apimeta/...` — NormalizeScope maps explicit Platform to nil; CanonicalScopeIdentity nil → Platform/PlatformScopeUID; PlatformScopeUID not a valid generated uid; UID generation produces opaque 128-bit values_

  - [ ] 2.3 Implement ListEnvelope generic type and Page struct (`internal/apimeta`)
    - ListEnvelope[T] with TypeMeta anonymous embed (json promotion), Items, Page
    - Page struct with opaque NextPageToken
    - _Design: data models section (ListEnvelope)_
    - _Requirements: F12-LIST-001, F12-LIST-002, F12-NAMING-002_
    - _Verification: JSON marshal of ListEnvelope produces flat apiVersion/kind at top level (no nesting)_

  - [ ] 2.4 Implement TypedRef re-export and reference constraint helpers (`internal/apiref`)
    - `reference.go`: type alias `TypedRef = apimeta.TypedRef`; `Refs` collection type
    - `constraints.go`: Constraint struct (AllowedKinds, AllowedScopes, Direction); `ValidateRef` method returning package-local `RefIssue` slice (field path + stable code + message); does NOT import `apiproblem`
    - _Design: data models section (apiref)_
    - _Requirements: F12-REF-001, F12-REF-002, F12-REF-003, F12-REF-004_
    - _Verification: `go test ./internal/apiref/...` — allowed/disallowed kinds, name/uid mismatch, provider-native rejection_

  - [ ] 2.5 Implement Condition type, SetCondition, and status enum (`internal/apicond`)
    - `condition.go`: Condition struct, ConditionStatus consts (True/False/Unknown), SetCondition (upsert; LastTransitionTime changes only on status change), Get helper
    - Enforce: type/reason are PascalCase; message is informational
    - _Design: D-04; data models section (apicond)_
    - _Requirements: F12-STATUS-002, F12-STATUS-003_
    - _Verification: `go test ./internal/apicond/...` — transition time invariant; conditions-not-history_

  - [ ] 2.6 Implement RFC 9457 Problem, Violation, stable codes, and generic HTTP status mapping (`internal/apiproblem`)
    - `problem.go`: Problem struct with type/title/status/detail/instance/code/requestId/violations
    - `codes.go`: stable ErrorCode constants + violation-code registry (including OPERATION_TARGET_SCOPE_MISMATCH)
    - `httpmap.go`: failure-class to HTTP status mapping per F12-ERROR-002 baseline (generic mapping only; no Decision/SafeDenial here)
    - _Design: D-05; error handling section_
    - _Requirements: F12-ERROR-001, F12-ERROR-002, F12-ERROR-003, F12-ERROR-004_
    - _Verification: `go test ./internal/apiproblem/...` — codes match baseline table; Problem serializes to RFC 9457_

  - [ ] 2.7 Add apiref issue-to-violation translator in apivalid
    - In `internal/apivalid/translate.go`: function `RefIssuesToViolations([]apiref.RefIssue) []apiproblem.Violation` translating package-local apiref results into apiproblem violations
    - _Design: D-02 (import direction); D-04 (validation pipeline)_
    - _Requirements: F12-VALIDATION-006_
    - _Verification: unit test — RefIssue with path/code maps to equivalent Violation_

  - [ ] 2.8 Write property test for condition transition semantics (Property 8)
    - **Property 8: Condition transition semantics**
    - For any sequence of condition upserts, SetCondition advances lastTransitionTime iff status changed; status always True/False/Unknown; type/reason stable PascalCase
    - Minimum 100 generated iterations; deterministic seed or report failing seed
    - **Validates: Requirements 4.8 (F12-STATUS-002, F12-STATUS-003)**
    - _Verification: `go test ./internal/apicond/...` passes with -count=1_

- [ ] 3. Checkpoint — primitives compile and pass
  - Run: `make fmt`; `git diff --check`; `go test ./internal/apimeta/... ./internal/apiref/... ./internal/apicond/... ./internal/apiproblem/... ./internal/apivalid/... ./internal/apiconform/...`; `go test -race ./internal/apimeta/... ./internal/apiref/... ./internal/apicond/... ./internal/apiproblem/... ./internal/apivalid/... ./internal/apiconform/...`; `go vet ./internal/apimeta/... ./internal/apiref/... ./internal/apicond/... ./internal/apiproblem/... ./internal/apivalid/... ./internal/apiconform/...`
  - Covers: apimeta, apiref, apicond, apiproblem (tasks 2.1–2.6), apivalid task 2.7 (RefIssuesToViolations), apiconform task 1.3 (import-boundary test)
  - All tests must pass before proceeding to task group 4

- [ ] 4. Strict JSON/YAML decoding and operation-aware field policies
  - [ ] 4.1 Implement pure DecodeJSON with duplicate-key detection (`internal/apivalid/decode.go`)
    - encoding/json with DisallowUnknownFields + token-scan duplicate-key detector
    - Returns `*apiproblem.Problem` with stable code and JSON Pointer on failure
    - Accepts Limits and FieldPolicy parameters; no HTTP dependency
    - _Design: D-03, D-15_
    - _Requirements: F12-VALIDATION-001(2), F12-VALIDATION-002, F12-VALIDATION-006_
    - _Verification: unit tests — duplicate key rejected; unknown field rejected; stable code + JSON Pointer in error_

  - [ ] 4.2 Implement strict JSON-compatible DecodeYAML (`internal/apivalid/decode.go`)
    - Using gopkg.in/yaml.v3 with this exact ordered pipeline:
      1. yaml.Node safety parsing (syntax tree only)
      2. Reject YAML-only constructs: aliases, anchors, merge keys (<<), custom/explicit tags, non-finite numbers (.nan/.inf), YAML-only timestamp/binary coercions, multiple documents, non-string mapping keys
      3. Explicit yaml.Node duplicate-key detection pass
      4. Normalize the accepted YAML node to a JSON-compatible value
      5. Marshal that normalized value to JSON bytes
      6. Pass those JSON bytes through the same DecodeJSON path (unknown-field rejection, FieldPolicy enforcement, same destination Go type, same stable error-code and JSON Pointer mapping)
    - NO direct yaml.v3 typed decoding after normalization
    - NO dependence on YAML struct tags
    - NO KnownFields(true) on yaml.v3 decoder
    - Same FieldPolicy and Limits as DecodeJSON; same Problem return type
    - JSON and YAML representations of the same object MUST produce equivalent validation results
    - _Design: D-03a_
    - _Requirements: F12-VALIDATION-001(2), F12-VALIDATION-002_
    - _Verification: unit tests — each rejected YAML feature fails; valid YAML decodes equivalently to JSON_

  - [ ] 4.3 Implement DecodeMode, FieldPolicy, and PolicyFor (`internal/apivalid/fieldpolicy.go`)
    - DecodeMode consts: ModeCreateRequest, ModeReplaceRequest, ModeStatusUpdate, ModeInternalObject, ModeReadRepresentation
    - FieldPolicy struct: Mode, AllowStatus, AllowSystemOwned, AllowSpecMutation
    - PolicyFor(mode) returns the correct policy per Matrix C2 ownership rules
    - _Design: D-15_
    - _Requirements: F12-VALIDATION-002, F12-META-002, F12-OWNER-002_
    - _Verification: unit tests — customer modes reject status/system; internal/read modes accept them_

  - [ ] 4.4 Implement HTTP decode adapter (`internal/apivalid/httpdecode.go`)
    - StrictDecode: MaxBytesReader, media-type selection (application/json, application/yaml accepted; application/x-yaml and text/yaml treated as yaml; all others → 415), delegates to DecodeJSON/DecodeYAML
    - Oversized body → 400 REQUEST_TOO_LARGE
    - _Design: D-03; validation section layer 1_
    - _Requirements: F12-VALIDATION-001(1), F12-ERROR-002_
    - _Verification: unit tests — 415 for unsupported type; 400 for oversized; correct decoder selected_

  - [ ] 4.5 Write property test for operation-aware field ownership (Property 3)
    - **Property 3: Operation-aware field ownership**
    - For any object with status/system fields, customer mutation modes reject while status-update/internal/read modes accept; field acceptance is a deterministic function of DecodeMode
    - Minimum 100 generated iterations; deterministic seed or report failing seed
    - **Validates: Requirements 4.3, 4.7, 4.9 (F12-VALIDATION-002, F12-META-002, F12-OWNER-002)**
    - _Verification: `go test ./internal/apivalid/...` passes with -count=1_

- [ ] 5. Bounded JSON Schema supported-subset validation and StructuralValidator interface
  - [ ] 5.1 Implement SupportedKeywords set, context-aware schema walker, and ValidateSchemaSupport (`internal/apischema/subsetvalidate.go`)
    - Define the exact bounded vocabulary — these keywords only:
      $schema, $id, $ref, title, description, type, properties, required,
      enum, items, additionalProperties, minLength, maxLength, minimum,
      maximum, pattern, default, examples
    - Plus the five registered x-sovrunn-* extension keywords (task 9.1)
    - $defs is PROHIBITED
    - Shared schemas use approved relative $ref values targeting api/schemas/_common only
    - Implement a context-aware schema walker that distinguishes:
      - Keys under `properties` are property identifiers, NOT keywords (must not be rejected)
      - Registered extension-object fields (e.g. fields inside x-sovrunn-field-policy) are validated by their extension schema, not treated as core JSON Schema keywords
      - Only actual schema-position keywords outside the supported set trigger fail-closed rejection
    - ValidateSchemaSupport scans a schema document and rejects any unsupported actual keyword (fail-closed) with a stable code
    - Returns package-local `SchemaIssue` slice (NOT `apiproblem.Violation`)
    - `apischema` MUST NOT import `apiproblem`
    - _Design: D-01a_
    - _Requirements: F12-NAMING-005, F12-VALIDATION-001(4)_
    - _Verification: unit tests proving:_
      - _unsupported keyword (e.g. `oneOf`, `if/then/else`, `$defs`) is rejected_
      - _property name under `properties` is NOT rejected_
      - _extension-object fields are not treated as keywords_
      - _document metadata ($schema, $id, title, description) accepted_
      - _fail-closed: no constraint silently ignored_

  - [ ] 5.2 Implement ValidateInstance for structural validation (`internal/apischema/subsetvalidate.go`)
    - Structurally validates a decoded instance against a canonical schema using only the supported subset
    - Returns package-local `SchemaIssue` slice with stable codes and JSON Pointer paths
    - Callers must first pass ValidateSchemaSupport
    - _Design: D-01a_
    - _Requirements: F12-VALIDATION-001(4), F12-VALIDATION-006_
    - _Verification: unit tests — valid instance passes; missing required field fails; wrong type fails; enum mismatch fails_

  - [ ] 5.3 Define StructuralValidator interface in apivalid (`internal/apivalid/structural.go`)
    - Interface signature: `Validate(instance any, schemaID string) ([]apiproblem.Violation, error)`
    - Returns BOTH violations AND an error:
      - `err != nil`: structural validation UNAVAILABLE; pipeline MUST stop at LayerStructural with Result.Problem = 500 INTERNAL_ERROR, Result.Err = internal cause, layers 5–7 MUST NOT execute
      - `err == nil, len(violations) > 0`: ordinary schema violations (422)
      - `err == nil, len(violations) == 0`: structurally valid
    - A nil StructuralValidator in Input is treated identically to a non-nil validator returning an error: stop, 500, no layers 5–7
    - `apivalid` MUST NOT import `apischema`
    - _Design: D-01a, D-04; import direction_
    - _Requirements: F12-VALIDATION-001(4), F12-IMPL-001_
    - _Verification: `go build ./internal/apivalid` compiles; interface is importable from apiconform without cycle; no Result/Problem/pipeline behavior tested here (see tasks 6.6, 6.7)_

  - [ ] 5.4 Write property test for fail-closed schema support (Property 1)
    - **Property 1: Fail-closed schema support**
    - For any schema document, if it contains a keyword outside the supported subset, ValidateSchemaSupport rejects it; a schema passes iff every keyword is in the set
    - Minimum 100 generated iterations; deterministic seed or report failing seed
    - **Validates: Requirements 4.1, 4.9 (F12-NAMING-005, F12-VALIDATION-001(4))**
    - _Verification: `go test ./internal/apischema/...` passes with -count=1_

- [ ] 6. Validation pipeline, limits, concurrency, and layer-8 authorization infrastructure
  - [ ] 6.1 Implement Limits struct with reviewed defaults (`internal/apivalid/limits.go`)
    - Limits struct with all fields per design (MaxObjectBytes 1MiB, MaxNestingDepth 32, MaxLabels 64, MaxLabelKeyChars 63, MaxLabelValueChars 253, MaxAnnotationsBytes 256KiB, MaxConditions 32, MaxReferencesPerField 64, MaxViolations 100, DefaultPageSize 50, MaxPageSize 200)
    - _Design: D-06_
    - _Requirements: F12-VALIDATION-007, F12-LIST-002_
    - _Verification: compiles; defaults match design table_

  - [ ] 6.2 Implement CallerContext, Decision, ScopeAuthorizer, AuthorizedResolver, AuthorizedTargetScopeResolver, and SafeDenial (`internal/apivalid/authz.go`)
    - CallerContext struct with `Scopes []ScopeIdentity`
    - Decision enum: Allow, DenyNotDisclosed, DenyKnown
    - ScopeAuthorizer interface: `Authorize(ctx, CallerContext, apiref.TypedRef, ScopeIdentity) Decision`
    - AuthorizedResolver interface: `Resolve(ctx, CallerContext, apiref.TypedRef) (obj any, found bool)` — single uniform unavailable outcome for absent and found-but-unauthorized
    - AuthorizedTargetScopeResolver interface: `ResolveAuthorizedTargetScope(ctx, CallerContext, apiref.TypedRef) (scope ScopeIdentity, available bool)` — available=false for both absent and unauthorized (hides which failed)
    - SafeDenial function: DenyNotDisclosed → identical 404 RESOURCE_NOT_FOUND; DenyKnown → 403 AUTHORIZATION_DENIED
    - SafeDenial is OWNED here in `internal/apivalid/authz.go`, NOT in `apiproblem`
    - _Design: D-04, D-17; layer-8 section_
    - _Requirements: F12-SEC-004, F12-SCOPE-002, F12-IMPL-002_
    - _Verification: unit tests — SafeDenial produces byte-identical Problem for DenyNotDisclosed vs "absent"; DenyKnown → 403_

  - [ ] 6.3 Implement CheckOperationTargetScopeMatch (`internal/apivalid/authz.go`)
    - Pure function: `CheckOperationTargetScopeMatch(opScope ScopeIdentity, targetScope ScopeIdentity) *apiproblem.Violation`
    - If scopes match → nil (no violation)
    - If scopes differ → Violation with code `OPERATION_TARGET_SCOPE_MISMATCH` and field `/metadata/scopeRef`
    - _Design: D-17; layer-8 section_
    - _Requirements: F12-SCOPE-002, F12-REF-001_
    - _Verification: unit tests — matching scopes (all six kinds including platform) → nil; kind mismatch → violation with correct code and path; UID mismatch → violation_

  - [ ] 6.4 Implement nine-layer validation pipeline with structural fail-closed and layer-8 configuration matrix (`internal/apivalid/pipeline.go`)
    - Layer enum (1..9), Result struct:
      - `Violations []apiproblem.Violation`
      - `FailedAt Layer`
      - `Problem *apiproblem.Problem` — safe client-facing failure (500 when structural unavailable or layer-8 misconfigured)
      - `Err error` — non-serialized internal diagnostic context
    - Structural fail-closed rules (layer 4):
      - If Input.Validator is nil: stop at LayerStructural, Result.Problem = 500 INTERNAL_ERROR, Result.Err = "nil validator", layers 5–7 MUST NOT execute
      - If Validator.Validate returns non-nil error: same behavior
      - If Validator.Validate returns ordinary violations: stop at
        LayerStructural, set FailedAt = LayerStructural, populate Violations,
        leave Problem and Err nil, and do not execute layers 5–7
      - Deterministic stub validators MUST be usable for pipeline tests
    - Layer-8 configuration matrix (when OperationScope is non-nil):
      - Path A requires: OperationScope, TargetRef, TargetScope (authoritative), Authorizer, Caller
      - Path B requires: OperationScope, TargetRef, TargetScopeResolver, Caller
      - Invalid configurations stop at LayerAuthorization with Result.Problem = 500 INTERNAL_ERROR, Result.Err = internal cause, no lookup, no success, no silent skip:
        - Neither path configured (no TargetScope AND no TargetScopeResolver)
        - Both paths configured (TargetScope AND TargetScopeResolver both set)
        - Missing Caller
        - Missing TargetRef
        - Incomplete Path A (TargetScope set but Authorizer nil)
        - Incomplete Path B (TargetScopeResolver set but Caller nil)
      - Valid Path A: ScopeAuthorizer.Authorize before lookup; DenyNotDisclosed → SafeDenial; Allow → CheckOperationTargetScopeMatch
      - Valid Path B: AuthorizedTargetScopeResolver; available=false → SafeDenial; available=true → CheckOperationTargetScopeMatch
    - Generic non-Operation validation (OperationScope nil): MAY stop successfully after layer 7
    - _Design: D-04, D-17; validation section_
    - _Requirements: F12-VALIDATION-001, F12-VALIDATION-004, F12-VALIDATION-005, F12-SCOPE-002_
    - _Verification: unit tests — structural fail-closed on nil/error; all layer-8 invalid configs produce 500; valid paths produce correct outcomes_

  - [ ] 6.5 Implement If-Match/resourceVersion stale-write helper (`internal/apivalid/concurrency.go`)
    - CheckIfMatch: compares If-Match header value with current resourceVersion; stale → 412 STALE_RESOURCE_VERSION Problem; match or absent → nil
    - _Design: D-10_
    - _Requirements: F12-UPDATE-002_
    - _Verification: unit tests — stale returns 412; match returns nil; absent returns nil_

  - [ ] 6.5a Define pipeline stage interfaces and invocation contract (`internal/apivalid/stages.go`)
    - Illustrative stage contracts:
      ```
      type DefaultingStage interface {
          Apply(ctx context.Context, object any) (objectAfterDefaults any, err error)
      }
      type ValidationStage interface {
          Validate(ctx context.Context, object any) ([]apiproblem.Violation, error)
      }
      type StageSet struct {
          Defaulting DefaultingStage
          Semantic   ValidationStage
          Reference  ValidationStage
      }
      ```
    - Input explicitly carries `Stages StageSet`
    - Binding rules:
      - Defaulting returns the object used by all later layers (semantic and reference receive the defaulted object)
      - Semantic and Reference receive the defaulted object
      - Stage implementations own immutable trusted rule configuration
      - Reference-stage construction receives trusted reference constraints and allowed scopes; these are not accepted as arbitrary caller input
      - A missing required stage fails closed at its corresponding layer
      - An explicitly constructed deterministic no-op stage is still invoked and is allowed only when the contract declares that layer inapplicable
      - Stage errors set Result.Problem to 500 INTERNAL_ERROR and Result.Err to the internal cause
      - Ordinary semantic/reference findings populate Result.Violations,
        set FailedAt to the current layer, leave Problem and Err nil, and
        stop before any later layer executes
      - Use LayerDefaulting, LayerSemantic, and LayerReference as FailedAt values
    - Pipeline invocation rule: full external-object pipeline MUST NOT silently omit a requested layer; a missing required stage implementation or stage-internal error fails closed at that layer with Result.Problem = 500 INTERNAL_ERROR, Result.Err = internal cause, no later layer execution
    - Preserve package import boundaries: stages are apivalid-owned interfaces, concrete implementations may live in apivalid or apiconform
    - _Design: D-04; validation pipeline layers 5–7_
    - _Requirements: F12-VALIDATION-001, F12-VALIDATION-004, F12-VALIDATION-005_
    - _Verification: `go build ./internal/apivalid` compiles; stage interfaces importable_

  - [ ] 6.5b Implement deterministic defaulting stage (`internal/apivalid/stage_defaulting.go`)
    - Implements DefaultingStage interface (Apply returns the defaulted object used by all later layers)
    - Common defaulting: canonical Platform scope normalization via NormalizeScope (explicit Platform → nil)
    - Deterministic explicit no-op when no applicable defaulting rules exist for the input kind
    - Returns error (fails closed) if defaulting logic encounters an internal fault
    - _Design: D-04; D-16 (canonical platform scope)_
    - _Requirements: F12-VALIDATION-004, F12-SCOPE-002_
    - _Verification: unit tests — Platform scope normalized; non-platform scope unchanged; no-op for unknown kind; internal fault → error_

  - [ ] 6.5c Implement semantic-validation stage (`internal/apivalid/stage_semantic.go`)
    - Implements ValidationStage interface for the Semantic slot in StageSet
    - Receives the defaulted object from DefaultingStage
    - Common semantic validation includes: grammar-level naming rules
      (resource name regex), enum value validation, condition type/reason
      PascalCase, phase/condition coherence, the rule that ownerRef MUST NOT
      replace a required scopeRef or act as a governance/security scope, and
      finite-limit enforcement (MaxLabels, MaxConditions,
      MaxAnnotationsBytes)
    - Do not reject an object solely because ownerRef and scopeRef identify
      the same target; replacement or governance misuse is the prohibited
      behavior
    - Returns Violations for ordinary failures (422 handling)
    - Returns error (fails closed) if internal logic error
    - _Design: D-04; D-06 (limits)_
    - _Requirements: F12-VALIDATION-004, F12-VALIDATION-005, F12-NAMING-001, F12-OWNER-001_
    - _Verification: unit tests — invalid name → violation; invalid enum →
      violation; ownerRef with a missing required scopeRef still fails;
      authorization never substitutes ownerRef for scopeRef; over-limit →
      violation; internal fault → error_

  - [ ] 6.5d Implement structural reference/kind/scope-validation stage (`internal/apivalid/stage_reference.go`)
    - Implements ValidationStage interface for the Reference slot in StageSet
    - Receives the defaulted object from DefaultingStage
    - Construction receives trusted reference constraints and allowed scopes; these are not accepted as arbitrary caller input
    - Applies typed-reference constraints using apiref.ValidateRef
    - Translates RefIssues to Violations using RefIssuesToViolations (task 2.7)
    - Validates scope-kind against schema-declared allowed-scopes captured
      in immutable trusted stage configuration at construction time, not as
      arbitrary runtime Input and not by importing apischema
    - Returns Violations for ordinary failures (422 handling)
    - Returns error (fails closed) if reference constraint configuration is missing or malformed
    - _Design: D-04; D-02 (import direction)_
    - _Requirements: F12-VALIDATION-004, F12-REF-001, F12-REF-002, F12-SCOPE-002_
    - _Verification: unit tests — disallowed kind → violation; scope not in allowed set → violation; valid ref passes; missing constraint config → error_

  - [ ] 6.5e Integrate stage invocation into pipeline.go layers 5–7
    - Pipeline.Run reads Input.Stages (StageSet) and invokes Defaulting (layer 5), Semantic (layer 6), Reference (layer 7) in order after structural validation passes
    - Defaulting.Apply returns the defaulted object; Semantic.Validate and Reference.Validate receive that defaulted object
    - A nil stage in StageSet where the contract requires one: fail closed at that layer (500 INTERNAL_ERROR)
    - A stage returning error: fail closed at that layer (500 INTERNAL_ERROR, Result.Err = cause, no later layers)
    - A stage returning violations: populate Result.Violations, set
      Result.FailedAt to that layer, leave Problem and Err nil, and stop
      before any later layer executes
    - Deterministic no-op stages: allowed only when declared; pipeline does not skip the invocation
    - _Design: D-04; validation pipeline_
    - _Requirements: F12-VALIDATION-001, F12-VALIDATION-004_
    - _Verification: unit tests — nil required stage → 500; stage error →
      500 + no later layers; semantic violations stop at LayerSemantic before
      Reference; reference violations stop at LayerReference; no-op stage
      passes through_

  - [ ] 6.6 Write structural fail-closed tests (`internal/apivalid/pipeline_test.go`)
    - Test: nil StructuralValidator → Result.Problem = 500 INTERNAL_ERROR, Result.Err set, FailedAt = LayerStructural, Violations empty, layers 5–7 not executed
    - Test: StructuralValidator returns non-nil error → same behavior
    - Test: deterministic stub validator returning violations →
      Result.Violations populated, FailedAt = LayerStructural, Problem nil,
      Err nil, and layers 5–7 do not execute
    - Test: deterministic stub validator returning no violations → clean pass through layers 5–7
    - _Design: D-01a, D-04_
    - _Requirements: F12-VALIDATION-001(4), F12-VALIDATION-004_
    - _Verification: `go test ./internal/apivalid/...` passes_

  - [ ] 6.7 Write layer-8 configuration matrix tests (`internal/apivalid/pipeline_test.go`)
    - Test: OperationScope non-nil, neither TargetScope nor TargetScopeResolver → 500 INTERNAL_ERROR at LayerAuthorization, no target lookup
    - Test: OperationScope non-nil, BOTH TargetScope AND TargetScopeResolver set → 500 INTERNAL_ERROR at LayerAuthorization, no target lookup
    - Test: OperationScope non-nil, missing Caller → 500 INTERNAL_ERROR at LayerAuthorization
    - Test: OperationScope non-nil, missing TargetRef → 500 INTERNAL_ERROR at LayerAuthorization
    - Test: Incomplete Path A (TargetScope set, Authorizer nil) → 500 INTERNAL_ERROR at LayerAuthorization
    - Test: Incomplete Path B (TargetScopeResolver set, Caller nil) → 500 INTERNAL_ERROR at LayerAuthorization
    - Test: Complete Path A (OperationScope, TargetRef, TargetScope, Authorizer, Caller) with Allow → CheckOperationTargetScopeMatch runs
    - Test: Complete Path B (OperationScope, TargetRef, TargetScopeResolver, Caller) with available=true → CheckOperationTargetScopeMatch runs
    - Test: Path A with DenyNotDisclosed → SafeDenial 404, no mismatch disclosed
    - Test: Path B with available=false → SafeDenial 404, no mismatch disclosed
    - Test: Generic non-Operation (OperationScope nil) → stops successfully after layer 7
    - _Design: D-17; layer-8 configuration matrix_
    - _Requirements: F12-SCOPE-002, F12-SEC-004_
    - _Verification: `go test ./internal/apivalid/...` passes_

  - [ ] 6.7a Write layer 5–7 ordering and fail-closed tests (`internal/apivalid/pipeline_test.go`)
    - Test: nil required defaulting stage → 500 INTERNAL_ERROR at LayerDefaulting, no later layers
    - Test: defaulting stage returns error → 500 at LayerDefaulting, no semantic or reference layer runs
    - Test: nil required semantic stage → 500 INTERNAL_ERROR at LayerSemantic, no reference layer
    - Test: semantic stage returns error → 500 at LayerSemantic, no reference layer
    - Test: nil required reference stage → 500 INTERNAL_ERROR at LayerReference
    - Test: reference stage returns error → 500 at LayerReference
    - Test: defaulting no-op + semantic violations → FailedAt =
      LayerSemantic and the reference stage does not execute
    - Test: defaulting no-op + clean semantic stage + reference violations →
      FailedAt = LayerReference
    - Test: layers execute in strict order (5 before 6 before 7) verified via side-effect ordering
    - _Design: D-04; validation pipeline layers 5–7_
    - _Requirements: F12-VALIDATION-004, F12-VALIDATION-005_
    - _Verification: `go test ./internal/apivalid/...` passes_

  - [ ] 6.8 Write property test for concurrency staleness (Property 10)
    - **Property 10: Concurrency staleness**
    - For any pair of resource versions, CheckIfMatch returns 412 exactly when If-Match != current resourceVersion; nil when they match or no protection required
    - Minimum 100 generated iterations; deterministic seed or report failing seed
    - **Validates: Requirements 4.11 (F12-UPDATE-002)**
    - _Verification: `go test ./internal/apivalid/...` passes with -count=1_

  - [ ] 6.9 Write property test for canonical platform scope (Property 4 — partial)
    - **Property 4 (partial): Canonical platform scope normalization and identity**
    - For any ScopeRef, NormalizeScope maps explicit Platform to nil; CanonicalScopeIdentity(nil) → {Platform, PlatformScopeUID}; CanonicalScopeIdentity(non-platform) → {ref.Kind, ref.UID}; normalization is idempotent
    - Uses a test-local allowed-scope contract ([]ScopeKind parameter) to validate: Platform allowed → nil accepted; Platform not allowed → violation
    - Does NOT depend on canonical schemas or annotations (those do not exist yet)
    - The complete Property 4 including "Platform allowed by schema annotation" runs in task 14.5 after schemas and structural adapter exist
    - Minimum 100 generated iterations; deterministic seed or report failing seed
    - **Validates: Requirements 4.4, 4.5 (F12-SCOPE-002, F12-REF-001)**
    - _Verification: `go test ./internal/apivalid/...` passes with -count=1_

  - [ ] 6.10 Write property test for safe-denial path/response equivalence (Property 5)
    - **Property 5: Safe-denial path and response equivalence**
    - For any denied cross-scope access, exists-but-inaccessible and absent produce byte-identical SafeDenial responses (404) and take the same control-flow path; no existence-dependent fast path permitted
    - Does NOT claim perfect constant-time execution
    - Minimum 100 generated iterations; deterministic seed or report failing seed
    - **Validates: Requirements 4.4, 7.4 (F12-SEC-004, F12-SCOPE-002)**
    - _Verification: `go test ./internal/apivalid/...` passes with -count=1_

  - [ ] 6.11 Write property test for Operation target-scope equality (Property 11)
    - **Property 11: Operation target-scope equality**
    - For any Operation, scopeRef MUST equal the resolved canonical governance scope of targetRef
    - Generated positive tests: each of six scopes with matching target; platform nil scope with platform target; matching non-platform UID
    - Generated negative tests: kind mismatch → OPERATION_TARGET_SCOPE_MISMATCH at /metadata/scopeRef; UID mismatch → same; unavailable target → SafeDenial 404 (no mismatch disclosed); unauthorized target → SafeDenial 404 (no mismatch disclosed)
    - Path A tests: complete config with Allow then match; DenyNotDisclosed → SafeDenial
    - Path B tests: complete config with available=true then match; available=false → SafeDenial
    - Configuration failure tests: neither path → 500; both paths → 500; missing Caller → 500; missing TargetRef → 500; incomplete Path A → 500; incomplete Path B → 500
    - Generic non-Operation (nil OperationScope) → successful stop after layer 7
    - Minimum 100 generated iterations; deterministic seed or report failing seed
    - **Validates: Requirements 4.4, 4.5 (F12-SCOPE-002, F12-REF-001; D-17)**
    - _Verification: `go test ./internal/apivalid/...` passes with -count=1_

- [ ] 7. Checkpoint — validation infrastructure compiles and passes
  - Run: `make fmt`; `git diff --check`; `go test ./internal/apimeta/... ./internal/apiref/... ./internal/apicond/... ./internal/apiproblem/... ./internal/apivalid/... ./internal/apischema/... ./internal/apiconform/...`; `go test -race ./internal/apimeta/... ./internal/apiref/... ./internal/apicond/... ./internal/apiproblem/... ./internal/apivalid/... ./internal/apischema/... ./internal/apiconform/...`; `go vet ./internal/apimeta/... ./internal/apiref/... ./internal/apicond/... ./internal/apiproblem/... ./internal/apivalid/... ./internal/apischema/... ./internal/apiconform/...`
  - Covers: all primitive packages (tasks 2.1–2.8), apivalid pipeline + stages + authz (tasks 4–6), apischema subset validation (task 5), apiconform import-boundary test (task 1.3), structural fail-closed tests (task 6.6), layer 5–7 tests (task 6.7a), layer-8 configuration tests (task 6.7)
  - All tests must pass before proceeding to task group 7a

- [ ] 7a. Immutable schema registry and safe local $ref resolution
  - [ ] 7a.1 Implement immutable SchemaRegistry interface and in-memory implementation (`internal/apiconform/schemaregistry.go`)
    - SchemaRegistry interface: `Load(schemaID string) (schema []byte, err error)`
    - In-memory implementation for tests: accepts pre-loaded schemas by stable schema ID
    - Repository canonical implementation: loads from `api/schemas/` and `api/schemas/_common/` paths
    - No process-global mutable registry; each instance is immutable after construction
    - No network access; rejects any schema ID that would trigger a network fetch
    - _Design: D-01a, D-02_
    - _Requirements: F12-VALIDATION-001(4), F12-IMPL-001_
    - _Verification: unit tests — load known ID returns schema; load unknown ID returns error; registry is immutable after construction_

  - [ ] 7a.2 Implement safe local $ref resolver (`internal/apiconform/refresolver.go`)
    - Resolves only approved relative $ref values under `api/schemas/_common`
    - Rejects remote URI references (http://, https://, ftp://, etc.)
    - Rejects absolute filesystem paths (leading /)
    - Rejects path traversal outside `api/schemas/_common` (../ sequences, symlink escape)
    - Rejects missing targets (referenced file does not exist in registry)
    - Detects reference cycles (A → B → A)
    - Enforces a finite reference depth (configurable, default 10)
    - Returns error on any rejection so the pipeline fails closed at LayerStructural
    - _Design: D-01a_
    - _Requirements: F12-VALIDATION-001(4), F12-IMPL-001_
    - _Verification: unit tests — valid local ref resolves; missing ref → error; remote URI → error; absolute path → error; traversal → error; cycle → error; depth overflow → error_

  - [ ] 7a.3 Implement immutable StructuralValidatorConfig (`internal/apiconform/structural_config.go`)
    - Define an immutable configuration value containing SchemaRegistry and
      RefResolver
    - Its constructor rejects a nil SchemaRegistry or nil RefResolver
    - It contains no process-global mutable state
    - Task 8.2 consumes this configuration; this task does not implement or
      compile-check the adapter before task 8.2 creates it
    - Registry/configuration/ref-resolution failures ultimately cause the
      adapter to return an error, triggering fail-closed behavior at
      LayerStructural
    - _Design: D-01a, D-04_
    - _Requirements: F12-VALIDATION-001(4), F12-VALIDATION-004_
    - _Verification: unit tests — valid dependencies create a config; nil
      registry and nil resolver are rejected_

- [ ] 8. StructuralValidator adapter, schema-to-violation translation, and JSON/YAML equivalence
  - [ ] 8.1 Add apischema SchemaIssue-to-Violation translator in apiconform (`internal/apiconform/structural.go`)
    - Function `SchemaIssuesToViolations([]apischema.SchemaIssue) []apiproblem.Violation` translating package-local apischema results into apiproblem violations
    - _Design: D-02 (import direction); D-01a_
    - _Requirements: F12-VALIDATION-006_
    - _Verification: unit test — SchemaIssue with path/code maps to equivalent Violation_

  - [ ] 8.2 Implement StructuralValidator adapter in apiconform (`internal/apiconform/structural.go`)
    - Implements the `apivalid.StructuralValidator` interface: `Validate(instance any, schemaID string) ([]apiproblem.Violation, error)`
    - Constructor receives StructuralValidatorConfig from task 7a.3
      explicitly
    - Uses the configured SchemaRegistry to load schema by ID and the
      configured RefResolver for $ref resolution
    - Calls `apischema.ValidateSchemaSupport` and `apischema.ValidateInstance`
    - Translates `SchemaIssue` results to `apiproblem.Violation` using the translator from task 8.1
    - Returns error when schema is not found, registry is misconfigured, or ref-resolution fails (triggers structural fail-closed)
    - No import cycle: apiconform → apischema, apiconform → apivalid, apiconform → apiproblem (all allowed)
    - _Design: D-01a, D-02, D-04_
    - _Requirements: F12-VALIDATION-001(4), F12-VALIDATION-006_
    - _Verification: unit test — adapter rejects invalid instance and returns violations; adapter accepts valid instance; missing schema returns error; nil registry → error; ref-resolution failure → error_

  - [ ] 8.3 Implement JSON/YAML decode-only equivalence tests (`internal/apivalid/decode_equiv_test.go`)
    - For a set of representative objects, assert that JSON and strict-YAML forms produce equivalent decoded output
    - Tests validate: identical typed values, identical error codes and JSON Pointers, identical unknown-field rejection, identical FieldPolicy enforcement
    - Confirm YAML-only constructs are rejected BEFORE JSON normalization
    - These tests use decode functions only — no full pipeline invocation
    - _Design: D-03a_
    - _Requirements: F12-VALIDATION-001(2), F12-VALIDATION-002_
    - _Verification: table-driven test passes for all representative objects_

  - [ ] 8.4 Implement full-pipeline JSON/YAML equivalence tests with test-local schemas (`internal/apiconform/yaml_equiv_test.go`)
    - Uses test-local in-memory schemas (loaded via SchemaRegistry from 7a.1) and test-local objects — NOT canonical fixtures from task 14
    - For each test-local object, JSON and strict-YAML forms decode AND validate (through the StructuralValidator adapter from 8.2) to equivalent results
    - Tests cover: typed values, unknown fields, FieldPolicy, stable codes, JSON Pointer paths
    - Full canonical-fixture equivalence testing remains in task 14.5 after canonical schemas and fixtures exist
    - Depends on: tasks 7a.1, 7a.2, 8.2 (SchemaRegistry + RefResolver + StructuralValidator adapter)
    - _Design: D-03a_
    - _Requirements: F12-VALIDATION-001(2), F12-VALIDATION-002_
    - _Verification: `go test ./internal/apiconform/...` passes_

  - [ ] 8.5 Write property test for JSON/YAML validation equivalence (Property 2)
    - **Property 2: JSON/YAML validation equivalence**
    - For any object expressible in the strict JSON-compatible subset, JSON and YAML decode/validate produce equivalent results; YAML with aliases/anchors/merge keys/custom tags/multiple docs/non-string keys/non-finite numbers is rejected
    - Depends on: tasks 4.1, 4.2, 8.2 (decoders + StructuralValidator adapter)
    - Minimum 100 generated iterations; deterministic seed or report failing seed
    - **Validates: Requirements 4.9 (F12-VALIDATION-001(2), F12-VALIDATION-002)**
    - _Verification: `go test ./internal/apiconform/...` passes with -count=1_

- [ ] 9. Schema annotations, extensions, route validation, and schema-diff gate
  - [ ] 9.1 Implement x-sovrunn-* extension parsing and vocabulary checks (`internal/apischema/annotations.go`)
    - ReadAnnotations: parses exactly these five registered extensions:
      - x-sovrunn-profile
      - x-sovrunn-boundary
      - x-sovrunn-allowed-scopes
      - x-sovrunn-stability
      - x-sovrunn-field-policy
    - x-sovrunn-field-policy is a strictly validated property-level object carrying exactly: classification, authorizedWriter, authorizedReaders, mutability, retention, redaction, residency, auditRequired
    - Unknown x-sovrunn-* extensions MUST fail closed (reject, never silently ignore)
    - Validates against controlled vocabularies (Profile, Boundary, Stability enums from apimeta)
    - Returns SchemaMeta + package-local SchemaIssue slice on invalid/missing/unknown extensions
    - `apischema` MUST NOT import `apiproblem`
    - _Design: D-08_
    - _Requirements: F12-NAMING-006, F12-VERIFY-001(1), F12-SEC-001_
    - _Verification: unit tests — valid annotations pass; missing annotation fails; invalid vocabulary fails; unknown x-sovrunn-foo fails closed_

  - [ ] 9.2 Implement route-form validator (`internal/apischema/route.go`)
    - ValidateRoute enforces `/apis/<group>/<version>/<plural-kebab>` pattern
    - Rejects unversioned public endpoints
    - _Design: D-09_
    - _Requirements: F12-NAMING-004_
    - _Verification: unit tests — valid routes pass; unversioned route rejected; malformed group rejected_

  - [ ] 9.3 Implement schema-diff change classifier (`internal/apischema/diff.go`)
    - ClassifyChange: compares old/new schemas and returns change classification (Compatible, Breaking, ReviewRequired) per the change-classification table
    - _Design: D-11_
    - _Requirements: F12-EVOLVE-002_
    - _Verification: unit tests — add optional field = compatible; remove field = breaking; add enum value = review required_

  - [ ] 9.4 Implement VerifyBaselineIntegrity and VerifyBaselineApproval (`internal/apischema/diff.go`)
    - VerifyBaselineIntegrity: recomputes SHA-256 digests of baseline files and compares to BASELINE_MANIFEST.json
    - VerifyBaselineApproval: baseline changes MUST fail unless accompanied by recorded approval evidence in BASELINE_APPROVALS.json containing: exact old digest, exact new digest, approving ADH or approval token, reviewer identity, and date
    - Co-editing baseline + manifest in the same commit without recorded approval evidence is NOT sufficient
    - _Design: D-11_
    - _Requirements: F12-EVOLVE-002, F12-VERIFY-001(10)_
    - _Verification: unit tests — tampered baseline fails integrity; missing approval evidence fails; valid approval with matching digests passes_

  - [ ] 9.5 Write property test for controlled baseline updates (Property 6)
    - **Property 6: Controlled baseline updates**
    - For any baseline change, gate fails unless accompanied by recorded approval evidence with matching digests and ADH/token/reviewer/date; co-editing baseline+manifest without evidence never passes
    - Minimum 100 generated iterations; deterministic seed or report failing seed
    - **Validates: Requirements 4.14, 4.16 (F12-EVOLVE-002, F12-VERIFY-001(10))**
    - _Verification: `go test ./internal/apischema/...` passes with -count=1_

- [ ] 10. Canonical schemas and _common schemas
  - [ ] 10.1 Create _common sub-schemas (`api/schemas/_common/`)
    - `type-meta.json`, `object-meta.json`, `typed-ref.json`, `scope-ref.json`, `owner-ref.json`, `condition.json`, `problem.json`, `violation.json`, `page.json`
    - Each declares JSON Schema 2020-12 with only the approved bounded vocabulary
    - No $defs; shared definitions use $ref to _common files only
    - Every boundary-crossing property in the _common schemas explicitly
      declares x-sovrunn-field-policy with exactly: classification,
      authorizedWriter, authorizedReaders, mutability, retention, redaction,
      residency, auditRequired
    - No inheritance algorithm is used for FEATURE-0012 _common schemas
    - _Design: D-01, D-08; files section_
    - _Requirements: F12-NAMING-005, F12-NAMING-006, F12-SEC-001_
    - _Verification: ValidateSchemaSupport and field-policy completeness
      checks pass for every _common schema; missing or unknown policy fields
      fail_

  - [ ] 10.2 Create eight canonical schemas (`api/schemas/`)
    - `project.json` (ManagedResource / customer-facing / Tenant)
    - `resource-pool.json` (ManagedResource / operator-facing / Provider)
    - `discovered-database.json` (ObservedExternalResource / adapter-facing / Provider)
    - `plugin-definition.json` (VersionedDefinition / plugin-facing / Platform)
    - `adapter-configuration.json` (ManagedResource / adapter-facing / Provider)
    - `placement-evaluation-request.json` (TransientRequestResult / internal-engine-facing / Project)
    - `operation.json` (LongRunningOperation / plugin-facing / Platform, Organization, OrganizationUnit, Tenant, Project, Provider) — exactly six scopes per ADH-2026-013
    - `audit-event.json` (ImmutableRecord / governance-only / Organization)
    - Each includes x-sovrunn-profile, x-sovrunn-boundary, x-sovrunn-allowed-scopes, x-sovrunn-stability; uses only supported-subset keywords; references _common via $ref; no $defs
    - Operation schema declares `x-sovrunn-allowed-scopes: [Platform, Organization, OrganizationUnit, Tenant, Project, Provider]`
    - Every property crossing an API boundary in all eight canonical schemas MUST explicitly declare x-sovrunn-field-policy containing exactly: classification, authorizedWriter, authorizedReaders, mutability, retention, redaction, residency, auditRequired
    - Field-policy declarations are explicit per-property; no undefined inheritance algorithm is used for FEATURE-0012 schemas
    - Verification of field-policy completeness:
      - every applicable property has all eight field-policy fields
      - no unknown field-policy field exists
      - all controlled values are valid
      - unknown x-sovrunn-* extensions fail closed
    - _Design: D-01, D-08, D-17; files section; architecture Matrix D_
    - _Requirements: F12-NAMING-005, F12-NAMING-006, F12-FIXTURE-002, F12-PROFILE-001, F12-SCOPE-002, F12-SEC-001_
    - _Verification: ValidateSchemaSupport passes for all eight; ReadAnnotations extracts valid metadata; Operation schema has exactly six scopes; every applicable property has complete explicit x-sovrunn-field-policy_

  - [ ] 10.3 Create baseline snapshots and BASELINE_MANIFEST.json (`api/schemas/baseline/`)
    - Copy initial schemas to baseline directory
    - Generate BASELINE_MANIFEST.json with SHA-256 digests
    - Create empty BASELINE_APPROVALS.json (initial baseline needs no prior-approval evidence)
    - _Design: D-11_
    - _Requirements: F12-EVOLVE-002, F12-VERIFY-001(10)_
    - _Verification: VerifyBaselineIntegrity passes against the generated manifest_

- [ ] 11. Canonical-schema-to-Go TypeBinding verification
  - [ ] 11.1 Implement TypeBinding struct and VerifyGoTypeAgainstSchema (`internal/apischema/typebinding.go`)
    - TypeBinding struct: SchemaPath + reflect.Type (generic; no import of concrete contract types)
    - VerifyGoTypeAgainstSchema: reflection-based check verifying property names, JSON tags, required vs optional (omitempty/pointer), primitive types, arrays/maps, embedded fields, enum-backed types, additionalProperties behavior
    - Returns package-local `SchemaIssue` slice (NOT apiproblem.Violation); mismatches carry stable codes
    - Fixture round-tripping is supporting evidence, not complete proof of Go-type/schema agreement
    - _Design: D-01b_
    - _Requirements: F12-NAMING-005, F12-VALIDATION-001(4), F12-VERIFY-001(13)_
    - _Verification: unit test with a deliberately mismatched Go type fails; correct type passes_

- [ ] 12. Conformance-only Go contract types, boundary ledger, and TypeBinding registry
  - [ ] 12.1 Implement conformance-only Go contract types (`internal/apiconform/contracts.go`)
    - Concrete Go structs for all eight canonical schemas: Project, ResourcePool, DiscoveredDatabase, PluginDefinition, AdapterConfiguration, PlacementEvaluationRequest, Operation, AuditEvent
    - Each struct uses apimeta types (TypeMeta, ObjectMeta, ScopeRef, etc.) and explicit JSON tags
    - Operation type includes fields for targetRef, scopeRef (per D-17), ownerRef
    - These are conformance-only types proving schema fit; they do NOT implement domain behavior
    - _Design: D-01b, D-17; architecture Matrix D_
    - _Requirements: F12-FIXTURE-001, F12-FIXTURE-002, F12-NAMING-005_
    - _Verification: `go build ./internal/apiconform` compiles; types match canonical schema structure_

  - [ ] 12.2 Register TypeBindings in apiconform (`internal/apiconform/bindings.go`)
    - Populate concrete TypeBinding slice mapping each canonical schema and _common sub-schema to its corresponding Go type
    - Tests call apischema.VerifyGoTypeAgainstSchema for each binding
    - Keeps concrete bindings in apiconform so apischema does not import apiconform (no cycle)
    - _Design: D-01b_
    - _Requirements: F12-NAMING-005, F12-VERIFY-001(13)_
    - _Verification: `go test ./internal/apiconform/...` — VerifyGoTypeAgainstSchema passes for all bindings_

  - [ ] 12.3 Write property test for derivative Go-type/schema consistency (Property 9)
    - **Property 9: Derivative Go-type / schema consistency**
    - For any registered TypeBinding, VerifyGoTypeAgainstSchema accepts iff the Go type matches the schema across the supported subset; deliberate mismatches are rejected
    - Minimum 100 generated iterations; deterministic seed or report failing seed
    - **Validates: Requirements 4.1, 4.9, 4.16 (F12-NAMING-005, F12-VERIFY-001(13))**
    - _Verification: `go test ./internal/apiconform/...` passes with -count=1_

  - [ ] 12.4 Create machine-readable boundary ledger (`docs/api/boundary-ledger.yaml`)
    - YAML document with entries for each boundary: customer-facing, operator-facing, internal-engine-facing, adapter-facing, plugin-facing, governance-only
    - Each entry records: purpose, owner, producers, consumers, allowed/prohibited data, authorization, audit, observability, failure behavior, versioning, replacement path, migration path, reassessment trigger
    - _Design: D-12_
    - _Requirements: F12-LEDGER-001_
    - _Verification: strict YAML parse; all F12-LEDGER-001 categories present per boundary_

  - [ ] 12.5 Implement deterministic boundary-ledger Markdown generator (`internal/apiconform/ledgergen.go`)
    - Reads `docs/api/boundary-ledger.yaml` (YAML source of truth)
    - Generates deterministic Markdown output for `docs/api/BOUNDARY_LEDGER.md`
    - Generator is owned under `internal/apiconform` (or a bounded script callable from Go tests)
    - Output is fully determined by the YAML input — same YAML always produces byte-identical Markdown
    - _Design: D-12_
    - _Requirements: F12-LEDGER-001_
    - _Verification: `go build ./internal/apiconform` compiles; generator produces Markdown from YAML_

  - [ ] 12.6 Generate initial BOUNDARY_LEDGER.md and add synchronization test (`internal/apiconform/ledgergen_test.go`)
    - Run generator to produce `docs/api/BOUNDARY_LEDGER.md`
    - Synchronization test: parses `docs/api/boundary-ledger.yaml`, generates expected Markdown, compares byte-for-byte with `docs/api/BOUNDARY_LEDGER.md`, fails when the derivative file is stale
    - The YAML remains the sole source of truth; the Markdown is always regenerable
    - _Design: D-12_
    - _Requirements: F12-LEDGER-001_
    - _Verification: `go test ./internal/apiconform/...` — sync test passes; deliberately stale file fails_

- [ ] 13. Checkpoint — schemas, contract types, ledger, type bindings, and validation infrastructure validated
  - Run: `make fmt`; `git diff --check`; `go test ./internal/apivalid/... ./internal/apischema/... ./internal/apiconform/...`; `go test -race ./internal/apivalid/... ./internal/apischema/... ./internal/apiconform/...`; `go vet ./internal/apivalid/... ./internal/apischema/... ./internal/apiconform/...`
  - Includes: apivalid, apischema, apiconform, package import-boundary tests (task 1.3), structural fail-closed tests (task 6.6), layer-8 configuration tests (task 6.7), JSON/YAML equivalence tests (tasks 8.3, 8.4)
  - All tests must pass before proceeding to task group 14

- [ ] 14. Eight conformance fixture families and Matrix D coverage
  - [ ] 14.1 Create valid JSON fixtures for all eight contracts (`tests/conformance/fixtures/`)
    - One valid JSON file per non-Operation contract: project.json, resource-pool.json, discovered-database.json, plugin-definition.json, adapter-configuration.json, placement-evaluation-request.json, audit-event.json
    - Operation fixture family:
      - One canonical representative fixture (operation.json) using Platform scope (nil scopeRef, platform target)
      - Six explicit scope variants (operation-platform.json, operation-organization.json, operation-organizationunit.json, operation-tenant.json, operation-project.json, operation-provider.json)
      - Each variant has exactly one Operation scope with scopeRef matching targetRef scope per D-17
      - Platform variant uses canonical nil scopeRef; non-platform variants use matching scope UID
      - Mismatch cases are separate negative fixtures (task 14.3)
    - Total: exactly eight canonical contract families
    - Each proves: schema fit, boundary classification, allowed scopes, ownership, strict parsing, reference behavior, absence of later-phase execution
    - _Design: architecture Matrix D; testing strategy section; D-17_
    - _Requirements: F12-FIXTURE-001, F12-FIXTURE-002, F12-SCOPE-002_
    - _Verification: each fixture decodes under ModeReadRepresentation and validates against its canonical schema; Operation variants each pass with their declared scope_

  - [ ] 14.2 Create valid YAML equivalents for each fixture (`tests/conformance/fixtures/`)
    - Strict JSON-compatible YAML form of each fixture for JSON/YAML equivalence testing
    - _Design: D-03a; testing strategy section_
    - _Requirements: F12-VALIDATION-001(2)_
    - _Verification: YAML decode produces equivalent output to JSON decode for each fixture_

  - [ ] 14.3 Create negative/invalid fixtures (`tests/conformance/fixtures/`)
    - Invalid fixtures covering: unknown field, duplicate key (JSON + YAML), unauthorized status in create mode, name/uid mismatch, invalid scope kind, ownerRef-as-scopeRef, nil scopeRef on non-Platform resource, oversized body, over-nested doc, unsupported media type, unversioned route, YAML aliases/anchors/merge keys/custom tags/multiple docs/non-string keys/non-finite numbers, Operation with scope/target mismatch
    - _Design: testing strategy section (negative tests); D-17_
    - _Requirements: F12-VALIDATION-002, F12-VALIDATION-006, F12-OWNER-001, F12-REF-002, F12-SCOPE-002_
    - _Verification: each negative fixture rejected with expected stable code and JSON Pointer_

  - [ ] 14.4 Implement fixture loader and Matrix D scenario assertions (`internal/apiconform/fixtures.go`)
    - Fixture loader reads JSON/YAML fixtures from `tests/conformance/fixtures/`
    - Matrix D scenario table maps all seventeen scenarios to a fixture + required-proof assertion
    - The "Future provisioning executes" scenario uses Operation scope variants (one per scope) and asserts target-scope equality per D-17 for each
    - Conformance test fails if any scenario is unrepresentable
    - _Design: testing strategy section (Matrix D coverage); D-17_
    - _Requirements: F12-FIXTURE-001, F12-FIXTURE-002, F12-VERIFY-001_
    - _Verification: `go test ./internal/apiconform/...` — all 17 scenarios pass_

  - [ ] 14.5 Implement conformance test suite (`internal/apiconform/fixtures_test.go`)
    - Positive: each fixture decodes (JSON + YAML), validates against its schema via StructuralValidator adapter (task 8.2), annotations valid
    - Operation-aware: same object accepted under internal/read but rejected under create/replace
    - Platform-scope canonical form (complete Property 4): explicit Platform normalizes to nil; identity tuple uses PlatformScopeUID via CanonicalScopeIdentity; "Platform allowed by schema" annotation verified; nil scopeRef rejected when Platform is not in x-sovrunn-allowed-scopes
    - Operation target-scope equality: each Operation scope variant with matching scope passes; mismatched scope → OPERATION_TARGET_SCOPE_MISMATCH at /metadata/scopeRef
    - JSON/YAML equivalence: both canonical-fixture JSON and YAML forms produce equivalent results
    - Negative: each invalid fixture rejected with correct code/pointer
    - Boundary limits: over-limit inputs rejected at boundary values
    - Security: secret-like values in metadata rejected; cross-scope denial uses SafeDenial (from apivalid/authz.go) with byte-identical 404; response/path-equivalence tested via stub AuthorizedResolver and AuthorizedTargetScopeResolver
    - _Design: testing strategy section; D-17_
    - _Requirements: F12-FIXTURE-002, F12-VALIDATION-002, F12-SEC-003, F12-SEC-004, F12-SCOPE-002_
    - _Verification: `go test ./internal/apiconform/...` passes_

  - [ ] 14.6 Write property test for provider neutrality of core (Property 7)
    - **Property 7: Provider neutrality of core**
    - For any core/customer-facing schema and any grammar primitive, no provider-native identifier/SDK type/provider-specific field is present; provider-native data only behind adapter-facing/plugin-facing schemas
    - Minimum 100 generated iterations; deterministic seed or report failing seed
    - **Validates: Requirements 4.4, 4.5, 4.6, 7.6 (F12-SCOPE-001, F12-SEC-006, F12-BOUNDARY-001)**
    - _Verification: `go test ./internal/apiconform/...` passes with -count=1_

- [ ] 15. Phase 1 compatibility report and coverage checks
  - [ ] 15.1 Create Phase 1 compatibility report (`docs/api/PHASE1_COMPATIBILITY_REPORT.md`)
    - Cover all Phase 1 resources: Organization, OrganizationUnit, Tenant, Project, Operation, ServiceClass, ServicePlan, Plugin, Capability, ServiceInstance, ServiceBinding, health/readiness, demo-flow
    - Per contract: record conforming behavior, explicit exceptions, and migration candidates
    - Note Phase 1 routes retained unchanged; no rewrite triggered
    - _Design: D-13_
    - _Requirements: F12-COMPAT-001, F12-COMPAT-002, F12-COMPAT-003_
    - _Verification: report covers all Phase 1 resources; exceptions documented_

  - [ ] 15.2 Implement Phase 1 coverage assertion (`internal/apiconform/compat.go`)
    - Asserts every Phase 1 resource is covered in the compatibility report
    - Fails if a Phase 1 resource is missing from the report
    - _Design: D-13_
    - _Requirements: F12-COMPAT-001_
    - _Verification: `go test ./internal/apiconform/...` passes the coverage check_

- [ ] 16. Fitness functions and feature-gate integration
  - [ ] 16.1 Implement schema, metadata, ownership, and field-policy fitness functions (`internal/apiconform/fitness_schema.go`)
    - Check 1: Every external schema declares profile/boundary/stability/allowed-scopes
    - Check 1a (field-policy coverage): Every property crossing an API boundary explicitly declares x-sovrunn-field-policy with exactly: classification, authorizedWriter, authorizedReaders, mutability, retention, redaction, residency, auditRequired; no unknown policy field accepted; controlled values are valid; unknown x-sovrunn-* extensions fail closed; verifies the explicit declarations used by FEATURE-0012 schemas (no inheritance algorithm)
    - Check 3: Every mutable field and condition has one owner
    - Check 4: Unknown and duplicate fields fail
    - Check 9: Published definitions are immutable
    - _Design: architecture §8.2; testing strategy section_
    - _Requirements: F12-VERIFY-001(1,3,4,9), F12-SEC-001, F12-SEC-002_
    - _Verification: `go test ./internal/apiconform/...` — checks 1, 1a, 3, 4, 9 pass; missing policy field fails; unknown policy field fails; unknown x-sovrunn-* fails_

  - [ ] 16.2 Implement reference, scope, boundary, and security fitness functions (`internal/apiconform/fitness_ref.go`)
    - Check 2: No core/customer schema imports provider SDK/native types
    - Check 5: References constrain kinds and scopes
    - Check 6: Cross-tenant access fails without existence disclosure (using SafeDenial from apivalid/authz.go)
    - Check 7: Raw secret-like values prohibited from metadata/status/errors
    - Check 8: Externally sourced observations include provenance and freshness
    - _Design: architecture §8.2; testing strategy section_
    - _Requirements: F12-VERIFY-001(2,5,6,7,8)_
    - _Verification: `go test ./internal/apiconform/...` — checks 2, 5, 6, 7, 8 pass_

  - [ ] 16.3 Implement compatibility, schema evolution, and limits fitness functions (`internal/apiconform/fitness_compat.go`)
    - Check 10: Schema compatibility detects breaking changes (uses schema-diff gate + VerifyBaselineIntegrity + VerifyBaselineApproval)
    - Check 11: Object/metadata/condition/violation/reference/page sizes bounded
    - Check 12: Errors use stable codes and JSON Pointer paths
    - Check 13: Generated artifacts (Go types) match canonical schema (uses TypeBinding verification)
    - _Design: architecture §8.2; testing strategy section; D-01b, D-11_
    - _Requirements: F12-VERIFY-001(10,11,12,13)_
    - _Verification: `go test ./internal/apiconform/...` — checks 10, 11, 12, 13 pass_

  - [ ] 16.4 Implement runtime absence, exception governance, and traceability fitness functions (`internal/apiconform/fitness_runtime.go`)
    - Check 14: Later-feature runtime behavior is absent (no provider/plugin/policy/placement/audit/provisioning service in grammar packages; no runtime HTTP routes added)
    - Check 15: Exceptions require an approved architecture handoff
    - _Design: architecture §8.2; testing strategy section_
    - _Requirements: F12-IMPL-002, F12-VERIFY-001(14,15)_
    - _Verification: `go test ./internal/apiconform/...` — checks 14, 15 pass_

  - [ ] 16.5 Implement fitness function aggregation and boundary-ledger check (`internal/apiconform/fitness.go`)
    - Aggregation: prove all checks 1–15 are registered and executed in a single test
    - Boundary-ledger: parse `docs/api/boundary-ledger.yaml` strictly; assert every declared boundary carries all F12-LEDGER-001 categories; assert every boundary present in canonical schemas has a ledger entry
    - _Design: D-12; architecture §8.2_
    - _Requirements: F12-LEDGER-001, F12-VERIFY-001(1–15)_
    - _Verification: `go test ./internal/apiconform/...` — all fifteen pass; a boundary missing a category fails_

  - [ ] 16.6 Create api-conformance-check.sh script (`scripts/api-conformance-check.sh`)
    - Invokes `go test ./internal/apiconform/...` (fitness functions + schema-diff + coverage + baseline integrity + baseline approval + ledger)
    - Exit 0 on pass; exit 1 on failure
    - _Design: files section (scripts)_
    - _Requirements: F12-VERIFY-002_
    - _Verification: script exits 0 when all conformance tests pass_

  - [ ] 16.7 Wire FEATURE-0012 into feature-gate (`scripts/feature-gate.sh`)
    - Add FEATURE-0012 case that invokes `scripts/api-conformance-check.sh`
    - Ensure `make ff-feature-gate FEATURE=FEATURE-0012` triggers the new check
    - _Design: files section (gate)_
    - _Requirements: F12-VERIFY-002_
    - _Verification: `make ff-feature-gate FEATURE=FEATURE-0012` passes_

- [ ] 17. Baseline protected-review governance and full verification
  - [ ] 17.1 Configure CODEOWNERS for baseline protected review
    - Inspect existing CODEOWNERS for the repository owner pattern already in use
    - Update CODEOWNERS to include `api/schemas/baseline/**` using the existing owner identity; do NOT invent an owner
    - Branch-protection configuration is external evidence: collect or document available evidence; do not claim that repository code configured branch protection; if required evidence is unavailable, preserve PENDING_HUMAN_REVIEW
    - _Design: D-11_
    - _Requirements: F12-EVOLVE-002, F12-VERIFY-001(10)_
    - _Verification: CODEOWNERS lists baseline paths with existing owner; branch-protection evidence collected or marked PENDING_HUMAN_REVIEW_

  - [ ] 17.2 Run full automated verification suite
    - Execute: `make fmt`; `git diff --check`; `go test ./...`; `go test -race ./...`; `go vet ./...`; `make ff-feature-gate FEATURE=FEATURE-0012`
    - Fix any failures discovered
    - _Design: verification section (F12-VERIFY-002)_
    - _Requirements: F12-VERIFY-002, F12-IMPL-001_
    - _Verification: all commands exit 0_

  - [ ] 17.3 Collect human-review evidence and mark PENDING_HUMAN_REVIEW
    - Collect and stage for human review: test and feature-gate output, changed-file inventory, compatibility-report evidence, boundary-ledger evidence, Matrix E residual-risk review inputs
    - Mark the result PENDING_HUMAN_REVIEW and stop
    - The coding agent MUST NOT: write an approval token, mark the implementation approved, invent a reviewer, invent a review date, or accept residual risk on behalf of a human
    - _Design: verification section (F12-VERIFY-003)_
    - _Requirements: F12-VERIFY-003, F12-RISK-001_
    - _Verification: evidence collected; status marked PENDING_HUMAN_REVIEW; no approval token written by agent_

- [ ] 18. Final checkpoint — feature gate passes
  - Run: `make fmt`; `git diff --check`; `go test ./...`; `go test -race ./...`; `go vet ./...`; `make ff-feature-gate FEATURE=FEATURE-0012`
  - Confirm: no runtime routes added; no domain services created; requirements.md and design.md unchanged
  - All tests must pass; the feature gate must exit 0
  - Human semantic-review criteria (F12-VERIFY-003) are NOT verified by these commands

## Notes

- All eleven correctness-property tests are mandatory (Properties 1–11). Each requires at least 100 generated iterations with a deterministic seed (or reports the failing seed for reproducibility).
- Each task references specific design decisions (D-*) and requirement IDs (F12-*) for traceability.
- Checkpoints (tasks 3, 7, 13, 18) act as barriers between dependency waves; all prior tests must pass before proceeding.
- Checkpoint 3 covers: apimeta, apiref, apicond, apiproblem, apivalid (task 2.7), apiconform (task 1.3).
- `apischema` returns package-local `SchemaIssue` (MUST NOT import `apiproblem`); translation to `apiproblem.Violation` occurs in `apiconform` (task 8.1).
- `apiref` returns package-local `RefIssue`; translation occurs in `apivalid` (task 2.7).
- `apivalid` MUST NOT import `apischema`; layer 4 uses the injected `StructuralValidator` interface (task 5.3) returning `([]apiproblem.Violation, error)`.
- Task 5.3 defines the interface only; Result/Problem/pipeline behavior tests are in tasks 6.6, 6.7, 6.7a.
- Layers 5–7 (defaulting, semantic, reference) use DefaultingStage and ValidationStage interfaces carried via Input.Stages (StageSet) (tasks 6.5a–6.5e); a missing required stage or stage-internal error fails closed with 500 INTERNAL_ERROR.
- The concrete StructuralValidator adapter lives in `apiconform` (task 8.2), receives SchemaRegistry (task 7a.1) and RefResolver (task 7a.2) explicitly, and MUST exist before any full-pipeline external-object equivalence test (task 8.4).
- Task 8.4 uses test-local schemas; full canonical-fixture equivalence is in task 14.5.
- The immutable SchemaRegistry (task 7a) resolves only approved relative $ref under api/schemas/_common; rejects remote URIs, absolute paths, traversal, cycles, and missing targets.
- ScopeIdentity, CanonicalScopeIdentity, and PlatformScopeUID are defined in `internal/apimeta/scope.go` (task 2.2).
- Decision, ScopeAuthorizer, AuthorizedResolver, AuthorizedTargetScopeResolver, SafeDenial, and CheckOperationTargetScopeMatch are implemented in `internal/apivalid/authz.go` (tasks 6.2, 6.3); they import and use apimeta types but do NOT exist in `apiproblem`.
- Concrete contract types live in `internal/apiconform/contracts.go`; the generic TypeBinding struct and VerifyGoTypeAgainstSchema live in `internal/apischema/typebinding.go`; no import cycle is introduced.
- Operation schema and fixtures declare exactly six allowed scopes: Platform, Organization, OrganizationUnit, Tenant, Project, Provider (ADH-2026-013). Each Operation fixture variant has exactly one scope; no single fixture object simultaneously uses all six.
- The YAML decode path follows exactly: yaml.Node safety parsing → reject YAML-only constructs and duplicate keys → normalize to JSON-compatible value → marshal to JSON bytes → DecodeJSON. No direct yaml.v3 typed decoding, KnownFields(true), or YAML struct tag dependence.
- Property 4 (canonical platform scope): partial test (NormalizeScope + CanonicalScopeIdentity + test-local allowed-scope contract) runs in task 6.9; complete test including schema annotation runs in task 14.5.
- Boundary ledger: YAML is sole source of truth; deterministic Markdown generator (task 12.5) produces BOUNDARY_LEDGER.md; synchronization test (task 12.6) enforces byte-for-byte consistency.
- No task implements provider/substrate services, adapters, plugin execution, policy engines, placement behavior, audit/operation services, provisioning/persistence services, new runtime HTTP routes, PATCH/watch/status-update routes, Phase 1 route rewrites, or a general-purpose JSON Schema engine (F12-IMPL-001, F12-IMPL-002).
- Human semantic-review evidence (F12-VERIFY-003) is collected by the agent but approval is PENDING_HUMAN_REVIEW (task 17.3); the agent must not write approval tokens.
- BASELINE_MANIFEST.json is an integrity mechanism; the approval boundary is protected review (CODEOWNERS + branch protection), not the manifest alone.
- CODEOWNERS uses the existing repository owner pattern; the coding agent must not invent an owner identity.
- Branch-protection configuration is external evidence; the coding agent collects/documents available evidence but does not claim to have configured it.

## Traceability Summary

| Task Group | Design Decisions | Requirements | Verification Evidence |
|---|---|---|---|
| 1. Repository & Go 1.22 alignment | D-14, D-02 | F12-IMPL-001, F12-VERIFY-001(2,14), F12-VERIFY-002 | `go build`; import-direction + package-boundary test |
| 2. Shared primitives | D-04, D-05, D-07, D-08, D-16, D-17 | F12-NAMING-001/002/006, F12-META-001/002/004, F12-REF-001/002/003/004, F12-SCOPE-002, F12-OWNER-001, F12-STATUS-002/003, F12-ERROR-001/002/003/004, F12-SEC-002, F12-LIST-001/002, F12-PROFILE-001, F12-BOUNDARY-001, F12-VALIDATION-006 | `go test ./internal/api{meta,ref,cond,problem,valid,conform}/...` |
| 4. Strict decoding & field policies | D-03, D-03a, D-15 | F12-VALIDATION-001/002/006, F12-META-002, F12-OWNER-002, F12-ERROR-002 | `go test ./internal/apivalid/...` |
| 5. Bounded schema validation | D-01a, D-02, D-04 | F12-NAMING-005, F12-VALIDATION-001(4)/006 | `go test ./internal/apischema/...`; fail-closed rejection |
| 6. Pipeline, limits, stages, layer-8 authz | D-04, D-06, D-10, D-16, D-17 | F12-VALIDATION-001/004/005/007, F12-SEC-004, F12-SCOPE-002, F12-UPDATE-002, F12-LIST-002, F12-IMPL-002, F12-NAMING-001, F12-OWNER-001, F12-REF-001/002 | `go test ./internal/apivalid/...`; structural fail-closed; layers 5–7 fail-closed; layer-8 config matrix |
| 7a. Schema registry & ref resolution | D-01a, D-02, D-04 | F12-VALIDATION-001(4), F12-VALIDATION-004, F12-IMPL-001 | `go test ./internal/apiconform/...`; ref-resolution tests |
| 8. StructuralValidator adapter + equivalence | D-01a, D-02, D-03a, D-04 | F12-VALIDATION-001(2,4), F12-VALIDATION-002/006 | `go test ./internal/apiconform/...`; JSON/YAML equivalence (test-local schemas) |
| 9. Annotations, extensions, routes, diff | D-08, D-09, D-11 | F12-NAMING-004/006, F12-EVOLVE-002, F12-VERIFY-001(1,10), F12-SEC-001 | `go test ./internal/apischema/...` |
| 10. Canonical schemas & _common | D-01, D-08, D-11, D-17 | F12-NAMING-005/006, F12-FIXTURE-002, F12-PROFILE-001, F12-SCOPE-002, F12-EVOLVE-002 | ValidateSchemaSupport passes; Operation has six scopes |
| 11. TypeBinding verification | D-01b | F12-NAMING-005, F12-VERIFY-001(13) | `go test` — VerifyGoTypeAgainstSchema |
| 12. Contract types, ledger & bindings | D-01b, D-12, D-17 | F12-LEDGER-001, F12-NAMING-005, F12-FIXTURE-001/002, F12-VERIFY-001(13) | Ledger parses; ledger sync test passes; TypeBindings pass; contract types compile |
| 14. Conformance fixture families & Matrix D | D-03a, D-04, D-16, D-17 | F12-FIXTURE-001/002, F12-VALIDATION-001/002/006, F12-SEC-003/004, F12-SCOPE-002, F12-OWNER-001, F12-REF-002, F12-BOUNDARY-001 | `go test ./internal/apiconform/...`; 17 scenarios; Operation 6 scope variants |
| 15. Phase 1 compatibility | D-13 | F12-COMPAT-001/002/003 | Coverage assertion; report exists |
| 16. Fitness functions & gate | D-11, D-12, D-01b | F12-VERIFY-001(1–15), F12-VERIFY-002, F12-IMPL-002, F12-LEDGER-001, F12-SEC-001/002 | `make ff-feature-gate FEATURE=FEATURE-0012` exit 0 |
| 17. Governance, verification & review | D-11 | F12-VERIFY-002, F12-VERIFY-003, F12-RISK-001, F12-EVOLVE-002 | All `make` commands pass; PENDING_HUMAN_REVIEW |

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1", "1.2"], "barrier": false },
    { "id": 1, "tasks": ["1.3", "2.1"], "barrier": false },
    { "id": 2, "tasks": ["2.2", "2.3", "2.4", "2.5", "2.6"], "barrier": false },
    { "id": 3, "tasks": ["2.7", "2.8"], "barrier": false },
    { "id": 4, "tasks": ["3"], "barrier": true },
    { "id": 5, "tasks": ["4.1", "4.3", "6.1"], "barrier": false },
    { "id": 6, "tasks": ["4.2", "4.4", "4.5"], "barrier": false },
    { "id": 7, "tasks": ["5.1", "5.3", "6.2", "6.3", "6.5a"], "barrier": false },
    { "id": 8, "tasks": ["5.2", "6.4", "6.5", "6.5b", "6.5c", "6.5d"], "barrier": false },
    { "id": 9, "tasks": ["6.5e"], "barrier": false },
    { "id": 10, "tasks": ["5.4", "6.6", "6.7", "6.7a", "6.8", "6.9", "6.10", "6.11"], "barrier": false },
    { "id": 11, "tasks": ["7"], "barrier": true },
    { "id": 12, "tasks": ["7a.1", "7a.2"], "barrier": false },
    { "id": 13, "tasks": ["7a.3"], "barrier": false },
    { "id": 14, "tasks": ["8.1", "8.2", "8.3", "9.1", "9.2", "9.3"], "barrier": false },
    { "id": 15, "tasks": ["8.4", "8.5", "9.4"], "barrier": false },
    { "id": 16, "tasks": ["9.5", "10.1", "10.2"], "barrier": false },
    { "id": 17, "tasks": ["10.3", "11.1"], "barrier": false },
    { "id": 18, "tasks": ["12.1", "12.4", "12.5"], "barrier": false },
    { "id": 19, "tasks": ["12.2", "12.3", "12.6"], "barrier": false },
    { "id": 20, "tasks": ["13"], "barrier": true },
    { "id": 21, "tasks": ["14.1", "14.2", "14.3"], "barrier": false },
    { "id": 22, "tasks": ["14.4", "14.5", "14.6"], "barrier": false },
    { "id": 23, "tasks": ["15.1"], "barrier": false },
    { "id": 24, "tasks": ["15.2"], "barrier": false },
    { "id": 25, "tasks": ["16.1", "16.2", "16.3", "16.4"], "barrier": false },
    { "id": 26, "tasks": ["16.5"], "barrier": false },
    { "id": 27, "tasks": ["16.6"], "barrier": false },
    { "id": 28, "tasks": ["16.7"], "barrier": false },
    { "id": 29, "tasks": ["17.1"], "barrier": false },
    { "id": 30, "tasks": ["17.2"], "barrier": false },
    { "id": 31, "tasks": ["17.3"], "barrier": false },
    { "id": 32, "tasks": ["18"], "barrier": true }
  ]
}
```

## Counts

| Metric | Value |
|---|---|
| Top-level task groups | 19 |
| Leaf tasks | 78 |
| Waves | 33 |
| Checkpoints (barriers) | 4 (tasks 3, 7, 13, 18) |
| Property tests | 11 (Properties 1–11) |
