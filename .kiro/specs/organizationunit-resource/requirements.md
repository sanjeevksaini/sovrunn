# Requirements Document

## Introduction

FEATURE-0002 implements the `OrganizationUnit` resource as a delegated governance boundary under an `Organization` in Sovrunn. It is the second resource in the Phase 1 governance hierarchy:

```
Organization
  -> OrganizationUnit
      -> Tenant (future)
          -> Project (future)
```

OrganizationUnit depends on FEATURE-0001 (Organization Resource and Registry) being fully implemented. It reuses the existing project skeleton, API server, middleware chain, health endpoints, configuration, and error infrastructure established in FEATURE-0001.

The key differentiator from Organization is the parent-child relationship: every OrganizationUnit belongs to exactly one Organization, identified by `spec.organizationName`. Identity is composite — `organizationName + metadata.name` — meaning the same `metadata.name` may exist under different Organizations without conflict.

This feature covers only `OrganizationUnit`. It does not implement `Tenant`, `Project`, nested OrganizationUnit hierarchies, persistent storage, Kubernetes CRDs, ServiceOps, the Operation framework implementation, UI, AI agent execution, or SDE runtime transformation.

## Glossary

- **OrganizationUnit**: Delegated governance boundary under an Organization. Example: `ministry-health` under Organization `nic`. Defined in `docs/glossary.md`.
- **Organization**: Top-level administrative and governance boundary in Sovrunn. Parent of OrganizationUnit. Implemented in FEATURE-0001.
- **Registry**: In-memory, thread-safe store for platform resources. Backed by a `sync.RWMutex`-protected map. Replaceable in future phases.
- **Metadata**: Identity and classification fields on a resource: `name`, `displayName`, `labels`, `annotations`.
- **Spec**: Desired state fields of a resource. Set by users or GitOps.
- **Status**: Observed state fields of a resource. System-owned. Not accepted from user input.
- **DNS-label name**: A resource name that matches `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`. Lowercase, hyphen-separated, no spaces. 1 to 63 characters.
- **Composite Key**: The identity key for OrganizationUnit: `organizationName/name`. Uniqueness is enforced per this composite.
- **APIError**: Structured JSON error body with `code`, `message`, and optional `field`/`details`.
- **Error Code**: Stable string error code from the set: `VALIDATION_FAILED`, `RESOURCE_NOT_FOUND`, `RESOURCE_ALREADY_EXISTS`, `DELETE_BLOCKED`, `INTERNAL_ERROR`.
- **Server**: The `net/http`-based HTTP server established in FEATURE-0001.
- **Validator**: The package or functions that deterministically validate OrganizationUnit resource fields.
- **ChildBlockerChecker**: An interface injected into the Organization delete handler that checks whether an Organization has any blocking child resources. FEATURE-0002 plugs the OrganizationUnit registry into this interface.
- **OUHandler**: HTTP handler for OrganizationUnit CRUD endpoints.
- **OrganizationUnitRegistry**: In-memory, thread-safe store for OrganizationUnit resources keyed by composite key `organizationName/name`.

## Requirements

---

### Requirement 1: OrganizationUnit Resource Shape

**User Story:** As a platform operator, I want the OrganizationUnit resource to follow the canonical `metadata/spec/status` shape, so that it is consistent with all Phase 1 resources and can be evolved toward Kubernetes-compatible desired-state reconciliation.

#### Acceptance Criteria

1. THE `OrganizationUnit` struct SHALL include `APIVersion`, `Kind`, `Metadata`, `Spec`, and `Status` fields with JSON tags `apiVersion`, `kind`, `metadata`, `spec`, and `status` respectively.
2. THE `Metadata` struct SHALL use JSON tag `json:"name"` (no omitempty) for `Name`, and `json:"displayName,omitempty"`, `json:"labels,omitempty"`, `json:"annotations,omitempty"` for the optional fields.
3. THE `OrganizationUnitSpec` struct SHALL include `OrganizationName` (`json:"organizationName"`, required) and `Description` (`json:"description,omitempty"`).
4. THE `OrganizationUnitStatus` struct SHALL include `Phase` (`json:"phase"`) with valid values `"Active"`, `"Inactive"`, `"Deleting"`, `"Failed"`, and `Message` (`json:"message,omitempty"`).
5. WHEN the Server returns a successful OrganizationUnit response (HTTP 200 or 201), THE Server SHALL set `apiVersion` to `"platform.sovrunn.io/v1alpha1"` and `kind` to `"OrganizationUnit"` regardless of what the client submitted.
6. IF the top-level JSON request body on a POST or PUT request contains the key `status` — regardless of value (`{}`, `null`, `{"phase": ""}`, or any other value) — THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "status"`.
7. WHEN `POST /v1/organization-units` succeeds with HTTP 201, THE Server SHALL set `status.phase` to `"Active"` both in the stored registry entry and in the 201 response body. A successful POST means all validation, parent Organization lookup, duplicate checks, and registry creation succeeded. If any validation or creation error occurs, THE Server SHALL return an error response and SHALL NOT store the OrganizationUnit.

---

### Requirement 2: OrganizationUnit Name Validation

**User Story:** As a platform operator, I want deterministic OrganizationUnit name validation, so that only safe, consistent, and DNS-compatible names are accepted and stored.

#### Acceptance Criteria

1. IF `metadata.name` is absent or an empty string, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"` and skip further name validation rules.
2. IF `metadata.name` is non-empty and does not match the pattern `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"` indicating an invalid DNS-label format.
3. IF `metadata.name` is non-empty and its length exceeds 63 characters, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"` indicating the name is too long.
4. IF `metadata.name` is non-empty and matches `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$` and has length between 1 and 63 inclusive, THEN THE Validator SHALL return a result containing zero validation errors for the name field.
5. THE Validator SHALL be callable as a pure function without starting the HTTP server or instantiating the registry.

---

### Requirement 3: Parent Organization Reference Validation

**User Story:** As a platform operator, I want `spec.organizationName` to be validated and to reference an existing Organization, so that OrganizationUnits cannot be created without a valid governance parent.

#### Acceptance Criteria

1. IF `spec.organizationName` is absent or an empty string, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.organizationName"` indicating the field is required.
2. IF `spec.organizationName` is non-empty and does not match the DNS-label pattern `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$` or exceeds 63 characters, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.organizationName"` indicating an invalid DNS-label format.
3. WHEN a valid `POST /v1/organization-units` request is received and all field-level validation passes, THE Server SHALL verify that the Organization identified by `spec.organizationName` exists in the Organization registry before storing the OrganizationUnit.
4. IF the Organization identified by `spec.organizationName` does not exist, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "spec.organizationName"` and `error.message` indicating the referenced Organization does not exist. Other field-level validation errors (metadata.name format, spec.organizationName format) SHALL also prevent storage even when the parent Organization exists.
5. THE `spec.organizationName` field SHALL be immutable after creation. THE Server SHALL NOT allow changing `spec.organizationName` via PUT requests (enforced by requiring path `organizationName` to match body `spec.organizationName`).

---

### Requirement 4: Composite Identity and Uniqueness

**User Story:** As a platform operator, I want OrganizationUnit identity to be composite (organizationName + name), so that I can use the same OrganizationUnit name under different Organizations without conflict.

#### Acceptance Criteria

1. THE OrganizationUnitRegistry SHALL use the composite key `organizationName/name` as the unique identity for each OrganizationUnit.
2. WHEN two OrganizationUnits have the same `metadata.name` but different `spec.organizationName` values, THE OrganizationUnitRegistry SHALL store both without conflict.
3. WHEN a `POST /v1/organization-units` request has a `spec.organizationName` and `metadata.name` combination that already exists in the registry, THE Server SHALL return HTTP 409 with `error.code = "RESOURCE_ALREADY_EXISTS"`.
4. THE GET, PUT, and DELETE endpoints SHALL use the path segments `{organizationName}/{name}` to identify a specific OrganizationUnit by its composite key.

---

### Requirement 5: REST API Endpoints

**User Story:** As a platform operator, I want CRUD REST endpoints for OrganizationUnit, so that I can manage delegated governance boundaries through the API.

#### Acceptance Criteria

1. THE Server SHALL register the route `POST /v1/organization-units` for creating an OrganizationUnit.
2. THE Server SHALL register the route `GET /v1/organization-units` for listing all OrganizationUnits.
3. THE Server SHALL register the route `GET /v1/organization-units/{organizationName}/{name}` for retrieving a single OrganizationUnit by composite key.
4. THE Server SHALL register the route `PUT /v1/organization-units/{organizationName}/{name}` for updating a single OrganizationUnit by composite key.
5. THE Server SHALL register the route `DELETE /v1/organization-units/{organizationName}/{name}` for deleting a single OrganizationUnit by composite key.
6. THE Server SHALL use Go 1.21-compatible routing: register a collection pattern `/v1/organization-units` and a subtree pattern `/v1/organization-units/`, extracting `organizationName` and `name` path segments via `strings.TrimPrefix` and `strings.SplitN`, not Go 1.22 wildcard syntax.

---

### Requirement 6: Create OrganizationUnit — POST /v1/organization-units

**User Story:** As a platform operator, I want to create an OrganizationUnit via REST API, so that I can register a new delegated governance boundary under an existing Organization.

#### Acceptance Criteria

1. WHEN a valid `POST /v1/organization-units` request is received, THE Server SHALL validate the request body, verify the parent Organization exists, store the OrganizationUnit in the Registry, and return HTTP 201 with a JSON body containing the full OrganizationUnit resource including server-set `apiVersion`, `kind`, and `status.phase = "Active"`. THE Server SHALL store the OrganizationUnit only on the successful create path that produces HTTP 201. If validation, parent lookup, duplicate check, registry create, or response preparation fails, THE Server SHALL return an error response and SHALL NOT store the resource.
2. WHEN the request body contains a composite key (`spec.organizationName` + `metadata.name`) that already exists in the Registry, THE Server SHALL return HTTP 409 with `error.code = "RESOURCE_ALREADY_EXISTS"`.
3. WHEN the request body fails field validation (invalid `metadata.name` or invalid/missing `spec.organizationName`), THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`, `error.message` describing the failure, and `error.field` set to the dot-separated JSON path of the first invalid field.
4. WHEN the Organization referenced by `spec.organizationName` does not exist in the Organization registry, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "spec.organizationName"`.
5. WHEN the request body is not valid JSON, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.message` indicating malformed JSON.
6. WHEN the request `Content-Type` is not `application/json`, THE Server SHALL return HTTP 415 with `error.code = "VALIDATION_FAILED"`.
7. IF the top-level JSON request body contains the key `status` — regardless of value — THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "status"` without storing the resource.
8. THE Server SHALL reject unknown fields in the request body via `json.Decoder.DisallowUnknownFields()`.

---

### Requirement 7: Get OrganizationUnit — GET /v1/organization-units/{organizationName}/{name}

**User Story:** As a platform operator, I want to retrieve a single OrganizationUnit by its composite key, so that I can inspect its current metadata, spec, and status.

#### Acceptance Criteria

1. WHEN `GET /v1/organization-units/{organizationName}/{name}` is requested and the OrganizationUnit exists, THE Server SHALL return HTTP 200 with a JSON body containing the full OrganizationUnit resource shape (`apiVersion`, `kind`, `metadata`, `spec`, `status`).
2. WHEN `GET /v1/organization-units/{organizationName}/{name}` is requested and the OrganizationUnit does not exist, THE Server SHALL return HTTP 404 with `error.code = "RESOURCE_NOT_FOUND"`.
3. IF the `{organizationName}` or `{name}` path segment does not match the DNS-label pattern `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$` or exceeds 63 characters, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` without performing a registry lookup.
4. THE Server SHALL return a deep copy of the OrganizationUnit (including independent copies of `Metadata.Labels` and `Metadata.Annotations` maps); mutations to the returned struct by the caller SHALL NOT affect the state stored in the Registry.

---

### Requirement 8: List OrganizationUnits — GET /v1/organization-units

**User Story:** As a platform operator, I want to list all OrganizationUnits, so that I can see all registered delegated governance boundaries across all Organizations.

#### Acceptance Criteria

1. WHEN `GET /v1/organization-units` is requested, THE Server SHALL return HTTP 200 with a JSON body containing an `items` array of all stored OrganizationUnits.
2. WHEN no OrganizationUnits are stored, THE Server SHALL return HTTP 200 with `{"items": []}`.
3. WHEN `GET /v1/organization-units` is requested, THE Server SHALL return OrganizationUnits sorted in ascending lexicographic order of `spec.organizationName` first, then ascending lexicographic order of `metadata.name` within the same Organization.
4. THE response body top-level object SHALL contain only the `items` array field; no additional registry-internal fields SHALL appear at the top level.
5. WHEN the registry encounters an unexpected internal error during list, THE Server SHALL return HTTP 500 with `error.code = "INTERNAL_ERROR"`.

---

### Requirement 9: Update OrganizationUnit — PUT /v1/organization-units/{organizationName}/{name}

**User Story:** As a platform operator, I want to update mutable fields of an OrganizationUnit, so that I can change its description, display name, labels, and annotations over time.

#### Acceptance Criteria

1. WHEN a valid `PUT /v1/organization-units/{organizationName}/{name}` request is received and the OrganizationUnit exists, THE Server SHALL replace `metadata.displayName`, `metadata.labels`, `metadata.annotations`, and `spec.description` with the values from the request body, preserve `metadata.name`, `spec.organizationName`, and all `status` fields from the stored entry unchanged, and return HTTP 200 with the updated OrganizationUnit as a JSON body.
2. WHEN `PUT /v1/organization-units/{organizationName}/{name}` is requested and the OrganizationUnit does not exist, THE Server SHALL return HTTP 404 with `error.code = "RESOURCE_NOT_FOUND"`.
3. THE Server SHALL allow updates to `metadata.displayName`, `metadata.labels`, `metadata.annotations`, and `spec.description`.
4. THE Server SHALL treat `metadata.name` and `spec.organizationName` as immutable on update.
5. THE Server SHALL require `metadata.name` to be present and non-empty in the PUT request body. IF `metadata.name` is absent or empty, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "metadata.name"`. IF `metadata.name` is present but differs from the `{name}` path segment, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "metadata.name"`.
6. THE Server SHALL require `spec.organizationName` to be present and non-empty in the PUT request body. IF `spec.organizationName` is absent or empty, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "spec.organizationName"`. IF `spec.organizationName` is present but differs from the `{organizationName}` path segment, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "spec.organizationName"`.
7. IF the top-level JSON request body contains the key `status` — regardless of value — THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "status"`.
8. WHEN the request body is not valid JSON, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` regardless of whether the target resource exists.
9. WHEN the `Content-Type` header is not `application/json`, THE Server SHALL return HTTP 415 with `error.code = "VALIDATION_FAILED"` regardless of whether the target resource exists.
10. THE registry update operation SHALL preserve the stored `spec.organizationName` and SHALL NOT move an OrganizationUnit between Organizations.

---

### Requirement 10: Delete OrganizationUnit — DELETE /v1/organization-units/{organizationName}/{name}

**User Story:** As a platform operator, I want to delete an OrganizationUnit, so that I can remove delegated governance boundaries that are no longer needed.

#### Acceptance Criteria

1. WHEN `DELETE /v1/organization-units/{organizationName}/{name}` is requested and the OrganizationUnit exists, THE Server SHALL remove the OrganizationUnit from the Registry and return HTTP 204 with no body.
2. WHEN `DELETE /v1/organization-units/{organizationName}/{name}` is requested and the OrganizationUnit does not exist, THE Server SHALL return HTTP 404 with `error.code = "RESOURCE_NOT_FOUND"`.
3. IF the `{organizationName}` or `{name}` path segment does not match the DNS-label pattern or exceeds 63 characters, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`.
4. THE delete handler SHALL include a clearly commented placeholder at the point where an Operation record will be created once FEATURE-0005 is implemented: `// TODO(FEATURE-0005): emit Operation record — type: DeleteOrganizationUnit`.
5. Future Tenant child-resource blockers for OrganizationUnit deletion are out of scope for FEATURE-0002. The delete handler SHALL NOT check for child resources in this feature.

---

### Requirement 11: Organization Deletion Blocker Integration

**User Story:** As a platform operator, I want Organization deletion to be blocked when OrganizationUnits reference it, so that governance boundaries are not accidentally removed while delegated units depend on them.

#### Acceptance Criteria

1. WHEN `DELETE /v1/organizations/{name}` is requested and one or more OrganizationUnits have `spec.organizationName` equal to `{name}`, THE Server SHALL return HTTP 409 with `error.code = "DELETE_BLOCKED"` and `error.message` identifying `"OrganizationUnit"` as the blocking resource kind.
2. THE OrganizationUnit blocker SHALL implement the `ChildBlockerChecker` interface defined in FEATURE-0001 so that it can be injected into the existing Organization delete handler without modifying the handler logic.
3. THE blocker SHALL query the OrganizationUnitRegistry to count OrganizationUnits referencing the Organization being deleted and return a `BlockedBy{Kind: "OrganizationUnit", Count: n}` entry when `n > 0`.
4. WHEN the Organization has zero OrganizationUnits referencing it, THE blocker SHALL return an empty blocking set and the Organization delete SHALL proceed as defined in FEATURE-0001.

---

### Requirement 12: In-Memory Thread-Safe OrganizationUnit Registry

**User Story:** As a developer, I want a thread-safe in-memory OrganizationUnit registry, so that concurrent API requests can safely read and write OrganizationUnit state without data races.

#### Acceptance Criteria

1. THE OrganizationUnitRegistry SHALL store OrganizationUnits in a `map[string]resources.OrganizationUnit` protected by a `sync.RWMutex`, using the composite key `organizationName/name` as the map key.
2. THE OrganizationUnitRegistry SHALL use a read lock (`RLock`/`RUnlock`) for `GetOrganizationUnit` and `ListOrganizationUnits` operations.
3. THE OrganizationUnitRegistry SHALL use a write lock (`Lock`/`Unlock`) for `CreateOrganizationUnit`, `UpdateOrganizationUnit`, and `DeleteOrganizationUnit` operations.
4. THE OrganizationUnitRegistry SHALL NOT return a reference to its internal map; `GetOrganizationUnit` SHALL return a deep copy of the stored `OrganizationUnit` struct, including independent copies of `Metadata.Labels` and `Metadata.Annotations` maps.
5. WHEN `ListOrganizationUnits` is called, THE OrganizationUnitRegistry SHALL return a new `[]resources.OrganizationUnit` slice whose elements are deep copies of stored OrganizationUnits (including independent copies of maps).
6. THE OrganizationUnitRegistry SHALL accept `context.Context` as the first parameter on all public methods.
7. THE OrganizationUnitRegistry SHALL be instantiable as a zero-value struct or via a constructor; no package-level global variables SHALL hold registry state.
8. WHEN run under `go test -race ./...` with at least 10 concurrent goroutines performing mixed reads and writes, THE OrganizationUnitRegistry SHALL produce no data race reports.
9. THE OrganizationUnitRegistry SHALL return a typed error value (or sentinel error) distinguishing `RESOURCE_NOT_FOUND` from `RESOURCE_ALREADY_EXISTS` so that HTTP handlers can map registry errors to the correct APIError code without inspecting error message strings.
10. WHEN `ListOrganizationUnits` is called, THE OrganizationUnitRegistry SHALL return results sorted by `spec.organizationName` ascending, then `metadata.name` ascending (deterministic ordering).
11. THE OrganizationUnitRegistry SHALL remain storage-only and SHALL NOT directly depend on OrganizationRegistry. Parent existence checks SHALL be performed in the API/service layer before calling `CreateOrganizationUnit` or `UpdateOrganizationUnit`.

---

### Requirement 13: Safe JSON Decoding with Request Body Limit

**User Story:** As a security-conscious operator, I want all OrganizationUnit JSON decoding to follow the same safe decoding pattern as Organization, so that the API server is protected against oversized bodies and unknown fields.

#### Acceptance Criteria

1. THE Server SHALL apply `http.MaxBytesReader` with a limit of 1 MiB (1,048,576 bytes) to every OrganizationUnit request body before decoding.
2. THE Server SHALL use `json.Decoder.DisallowUnknownFields()` when decoding OrganizationUnit request bodies. Unknown fields SHALL be rejected automatically by the decoder.
3. IF JSON decoding returns an error because the body exceeds 1 MiB, THE Server SHALL return HTTP 413 with `error.code = "VALIDATION_FAILED"`. The body size check via `http.MaxBytesReader` takes priority over syntax errors — if a body both exceeds 1 MiB and contains malformed JSON, HTTP 413 is returned.
4. IF JSON decoding returns a syntax error, a type-mismatch error, or an unknown-field error (and the body does not exceed 1 MiB), THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and a message indicating the nature of the problem.
5. THE Server SHALL NOT echo back the raw request body in any error response.
6. THE `status` key SHALL be rejected explicitly before the typed decode regardless of its value.

---

### Requirement 14: Request Context Usage

**User Story:** As a developer, I want all OrganizationUnit registry and validation functions to accept `context.Context`, so that request deadlines, cancellations, and request-scoped values propagate correctly through the call stack.

#### Acceptance Criteria

1. THE OrganizationUnitRegistry methods `CreateOrganizationUnit`, `GetOrganizationUnit`, `ListOrganizationUnits`, `UpdateOrganizationUnit`, and `DeleteOrganizationUnit` SHALL each accept `context.Context` as their first parameter.
2. Pure validation functions MAY remain context-free if they do not perform I/O, registry lookups, or cancellation-aware work.
3. THE Server SHALL pass `r.Context()` to all OrganizationUnit registry method calls within HTTP handlers.
4. THE Server SHALL NOT store `context.Context` in any struct field.
5. THE Server SHALL NOT pass `context.Background()` or `context.TODO()` inside HTTP handlers in place of the request context.

---

### Requirement 15: Operation Framework Boundary (Placeholder Only)

**User Story:** As an architect, I want the FEATURE-0002 implementation to define a clear operation boundary, so that FEATURE-0005 can wire in Operation records for mutating actions without requiring API handler refactoring.

#### Acceptance Criteria

1. THE `internal/api/` handlers for `CreateOrganizationUnit`, `UpdateOrganizationUnit`, and `DeleteOrganizationUnit` SHALL include a clearly commented placeholder at the point where an Operation record will be created once FEATURE-0005 is implemented.
2. THE placeholder SHALL NOT import the `internal/operation/` package or call any Operation creation code; it SHALL be a code comment only.
3. THE placeholder comment SHALL identify the expected Operation type: `// TODO(FEATURE-0005): emit Operation record — type: CreateOrganizationUnit`, `UpdateOrganizationUnit`, or `DeleteOrganizationUnit` respectively.

---

### Requirement 16: Tests for Validation, Registry, and HTTP Handlers

**User Story:** As a developer, I want comprehensive tests for OrganizationUnit validation, the in-memory registry, and the HTTP handlers, so that correctness is proven by the test suite and regressions are caught automatically.

#### Acceptance Criteria

1. THE `internal/validation/` package SHALL include tests that verify: (a) valid OrganizationUnit names are accepted, (b) empty metadata.name is rejected, (c) names with uppercase letters are rejected, (d) names with spaces are rejected, (e) names exceeding 63 characters are rejected, (f) names with leading or trailing hyphens are rejected, (g) empty spec.organizationName is rejected, (h) invalid spec.organizationName format is rejected, (i) valid spec.organizationName passes format validation.
2. THE `internal/registry/` package SHALL include tests that verify: (a) Create stores an OrganizationUnit, (b) duplicate composite key Create returns `RESOURCE_ALREADY_EXISTS` error, (c) same metadata.name under different Organizations succeeds, (d) Get returns the stored OrganizationUnit by composite key, (e) Get for a non-existent composite key returns `RESOURCE_NOT_FOUND` error, (f) List returns all stored OrganizationUnits sorted by organizationName then name, (g) List on empty registry returns an empty slice, (h) Update modifies mutable fields, (i) Update of a non-existent OrganizationUnit returns `RESOURCE_NOT_FOUND` error, (j) Delete removes the OrganizationUnit, (k) Delete of a non-existent OrganizationUnit returns `RESOURCE_NOT_FOUND` error.
3. THE `internal/api/` package SHALL include HTTP handler tests using `net/http/httptest` that verify: (a) POST returns 201 for valid input with existing parent Organization, (b) POST returns 409 for duplicate composite key, (c) POST returns 400 for invalid metadata.name, (d) POST returns 400 for missing spec.organizationName, (e) POST returns 400 for non-existent parent Organization, (f) POST returns 400 for user-authored status field, (g) GET returns 200 for existing resource with full resource shape, (h) GET returns 404 for missing resource, (i) GET returns 400 for invalid path segments, (j) GET list returns 200 with `items` array sorted by organizationName then name, (k) PUT returns 200 for valid update, (l) PUT returns 404 for missing resource, (m) PUT returns 400 when metadata.name in body differs from path, (n) PUT returns 400 when spec.organizationName in body differs from path, (o) DELETE returns 204 for existing resource, (p) DELETE returns 404 for missing resource.
4. THE test suite SHALL include an Organization deletion blocker test that verifies: (a) deleting an Organization with OrganizationUnits returns 409 DELETE_BLOCKED, (b) deleting an Organization with zero OrganizationUnits returns 204.
5. WHEN `go test -race ./...` is executed with at least 10 concurrent goroutines performing mixed OrganizationUnit registry operations, THE test suite SHALL produce no data race reports.
6. ALL tests SHALL be deterministic and SHALL NOT depend on external services, network access, or filesystem state beyond the repository root.

---

### Requirement 17: Property-Based Tests

**User Story:** As a developer, I want property-based tests using `testing/quick` for OrganizationUnit validation and registry, so that edge cases are discovered by randomized input generation.

#### Acceptance Criteria

1. THE validation package SHALL include property tests using `testing/quick` that generate valid DNS-label names (matching `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$` with length 1–63) and verify they are accepted by the validator with zero errors for that field.
2. THE validation package SHALL include property tests using `testing/quick` that generate arbitrary strings and verify that strings outside the accepted DNS-label domain are rejected, while valid generated DNS-label names are accepted.
3. THE registry package SHALL include a property test using `testing/quick` that verifies: FOR ALL valid OrganizationUnit resources, Create followed by Get with the same composite key returns the equivalent resource (round-trip property).
4. THE registry package SHALL include a property test using `testing/quick` that verifies: FOR ALL sequences of Create operations, ListOrganizationUnits returns items sorted by organizationName ascending then name ascending (sort invariant).
5. THE registry package SHALL include a property test using `testing/quick` that verifies: Create is idempotent in error — creating the same composite key twice returns `RESOURCE_ALREADY_EXISTS` on the second call and the original entry remains unchanged (verify both error code and entry preservation).

---

### Requirement 18: Concurrency Test

**User Story:** As a developer, I want a dedicated concurrency test with at least 10 goroutines performing mixed OrganizationUnit registry operations, so that thread safety is proven under contention.

#### Acceptance Criteria

1. THE registry test suite SHALL include a concurrency test that launches at least 10 goroutines performing a mix of Create, Get, List, Update, and Delete operations concurrently on the OrganizationUnitRegistry.
2. WHEN run with `go test -race ./internal/registry/...`, THE concurrency test SHALL produce no data race reports.
3. THE concurrency test SHALL verify that no panic occurs and all operations return expected errors or success values. Success is a guaranteed outcome when no panics occur and all operations complete without data races.

---

### Requirement 19: Non-Goals

**User Story:** As an architect, I want clear boundaries on what FEATURE-0002 does NOT implement, so that scope is controlled and future features are not prematurely introduced.

#### Acceptance Criteria

1. THE implementation SHALL NOT implement Tenant, Project, or nested OrganizationUnit hierarchies.
2. THE implementation SHALL NOT implement persistent storage, Kubernetes CRDs, or any durable backend.
3. THE implementation SHALL NOT implement authentication, authorization, or RBAC.
4. THE implementation SHALL NOT implement the Operation framework (beyond placeholder comments).
5. THE implementation SHALL NOT implement ServiceOps, plugin execution, or AI agent execution.
6. THE implementation SHALL NOT implement UI, SDE runtime transformation, or billing.
7. THE implementation SHALL NOT implement child-resource blocker checks on OrganizationUnit delete (future Tenant blockers are out of scope).
8. THE implementation SHALL NOT implement `delegatedAdminGroups` or any governance policy fields beyond `spec.description`.

---

## Design Questions

1. **Parent Organization existence check on Create**: The requirement specifies returning `VALIDATION_FAILED` when the parent Organization does not exist. This is a cross-registry validation (OUHandler must query the OrganizationRegistry). The design should clarify whether to inject `OrganizationRegistryIface` into the OUHandler or into the OrganizationUnitRegistry. Recommendation: inject into the handler to keep the registry layer storage-only.

2. **Composite blocker for Organization delete**: FEATURE-0001 defined `ChildBlockerChecker` with a single `BlockedByChildren(ctx, orgName)` method. FEATURE-0002 must supply an implementation that queries the OrganizationUnitRegistry. The design should clarify how to compose multiple blockers when FEATURE-0003 (Tenant) adds another one — whether to use a composite/chain pattern or replace the single checker with a list.
