# Requirements Document

## Introduction

FEATURE-0003 implements the `Tenant` resource as the primary isolation and security boundary
in Sovrunn. It is the third resource in the Phase 1 governance hierarchy:

```
Organization
  -> OrganizationUnit
      -> Tenant
          -> Project (future)
```

Tenant depends on FEATURE-0001 (Organization) and FEATURE-0002 (OrganizationUnit) being fully
implemented. Every Tenant belongs to exactly one OrganizationUnit, identified by the composite
parent reference `spec.organizationName` + `spec.organizationUnitName`.

Identity is a three-part composite: `organizationName/organizationUnitName/name`.
The same `metadata.name` may exist under different OrganizationUnits without conflict.

This feature covers only `Tenant`. It does not implement `Project`, nested hierarchies,
persistent storage, Kubernetes CRDs, ServiceOps, the Operation framework implementation,
UI, AI agent execution, or SDE runtime transformation.

## Glossary

- **Tenant**: Primary isolation and security boundary under an OrganizationUnit.
  Example: `prod-tenant` under OrganizationUnit `ministry-health` in Organization `nic`.
- **Organization**: Top-level governance boundary. Implemented in FEATURE-0001.
- **OrganizationUnit**: Delegated governance boundary under Organization. Implemented in FEATURE-0002.
- **Registry**: In-memory, thread-safe store. Backed by `sync.RWMutex`-protected map.
- **Composite Key**: The identity key for Tenant: `organizationName/organizationUnitName/name`.
- **DNS-label name**: Matches `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`. 1 to 63 characters.
- **APIError**: Structured JSON error body with `code`, `message`, optional `field`/`details`.
- **Error Code**: Stable string from: `VALIDATION_FAILED`, `RESOURCE_NOT_FOUND`, `RESOURCE_ALREADY_EXISTS`, `DELETE_BLOCKED`, `INTERNAL_ERROR`.
- **TenantRegistry**: In-memory store for Tenant resources keyed by composite key.
- **TenantHandler**: HTTP handler for Tenant CRUD endpoints.

## Requirements

---

### Requirement 1: Tenant Resource Shape

**User Story:** As a platform operator, I want the Tenant resource to follow the canonical `metadata/spec/status` shape, so that it is consistent with all Phase 1 resources.

#### Acceptance Criteria

1. THE `Tenant` struct SHALL include `APIVersion`, `Kind`, `Metadata`, `Spec`, and `Status` fields with JSON tags `apiVersion`, `kind`, `metadata`, `spec`, and `status` respectively.
2. THE `Metadata` struct SHALL use `json:"name"` (no omitempty) for `Name`, and `json:"displayName,omitempty"`, `json:"labels,omitempty"`, `json:"annotations,omitempty"` for optional fields.
3. THE `TenantSpec` struct SHALL include `OrganizationName` (`json:"organizationName"`), `OrganizationUnitName` (`json:"organizationUnitName"`), and `Description` (`json:"description,omitempty"`).
4. THE `TenantStatus` struct SHALL include `Phase` (`json:"phase"`) with valid values `"Active"`, `"Inactive"`, `"Deleting"`, `"Failed"`, and `Message` (`json:"message,omitempty"`).
5. WHEN the Server returns a successful Tenant response (HTTP 200 or 201), THE Server SHALL set `apiVersion` to `"platform.sovrunn.io/v1alpha1"` and `kind` to `"Tenant"` regardless of client input.
6. IF the top-level JSON request body on POST or PUT contains the key `status` — regardless of value (`{}`, `null`, `{"phase":""}`, or any other value) — THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "status"`.
7. WHEN `POST /v1/tenants` succeeds with HTTP 201, THE Server SHALL set `status.phase` to `"Active"` both in the stored entry and in the response body. A successful POST means all validation, parent lookup, duplicate checks, and registry creation succeeded. If any error occurs, THE Server SHALL return an error and SHALL NOT store the Tenant.

---

### Requirement 2: Tenant Name Validation

**User Story:** As a platform operator, I want deterministic Tenant name validation, so that only safe, DNS-compatible names are accepted.

#### Acceptance Criteria

1. IF `metadata.name` is absent or empty, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"`.
2. IF `metadata.name` does not match `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"`.
3. IF `metadata.name` exceeds 63 characters, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"`.
4. IF `metadata.name` is valid (1–63 chars, DNS-label pattern), THEN THE Validator SHALL return zero errors for that field.
5. THE Validator SHALL be a pure function (context-free, no I/O).

---

### Requirement 3: Parent Reference Validation

**User Story:** As a platform operator, I want `spec.organizationName` and `spec.organizationUnitName` validated and referencing existing resources, so that Tenants cannot be created without a valid parent.

#### Acceptance Criteria

1. IF `spec.organizationName` is absent or empty, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.organizationName"`.
2. IF `spec.organizationName` does not match the DNS-label pattern or exceeds 63 characters, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.organizationName"`.
3. IF `spec.organizationUnitName` is absent or empty, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.organizationUnitName"`.
4. IF `spec.organizationUnitName` does not match the DNS-label pattern or exceeds 63 characters, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.organizationUnitName"`.
5. WHEN a valid `POST /v1/tenants` request is received and field-level validation passes, THE Server SHALL verify that the OrganizationUnit identified by `spec.organizationName/spec.organizationUnitName` exists before storing the Tenant.
6. IF the parent OrganizationUnit does not exist, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`, `error.field = "spec.organizationUnitName"`, and `error.message` identifying the missing parent reference `spec.organizationName/spec.organizationUnitName`.
7. `spec.organizationName` and `spec.organizationUnitName` SHALL be immutable after creation. PUT requests must include both fields matching the path segments.

---

### Requirement 4: Composite Identity and Uniqueness

**User Story:** As a platform operator, I want Tenant identity to be composite, so that the same tenant name can exist under different OrganizationUnits.

#### Acceptance Criteria

1. THE TenantRegistry SHALL use the composite key `organizationName/organizationUnitName/name` as the unique identity.
2. WHEN two Tenants have the same `metadata.name` but different parent paths, THE TenantRegistry SHALL store both without conflict.
3. WHEN a POST request has a composite key that already exists, THE Server SHALL return HTTP 409 with `error.code = "RESOURCE_ALREADY_EXISTS"`.
4. GET, PUT, and DELETE endpoints SHALL use path segments `{organizationName}/{organizationUnitName}/{name}` to identify a Tenant.

---

### Requirement 5: REST API Endpoints

**User Story:** As a platform operator, I want CRUD REST endpoints for Tenant.

#### Acceptance Criteria

1. THE Server SHALL register `POST /v1/tenants` for creating a Tenant.
2. THE Server SHALL register `GET /v1/tenants` for listing all Tenants.
3. THE Server SHALL register `GET /v1/tenants/{organizationName}/{organizationUnitName}/{name}` for retrieving a Tenant.
4. THE Server SHALL register `PUT /v1/tenants/{organizationName}/{organizationUnitName}/{name}` for updating a Tenant.
5. THE Server SHALL register `DELETE /v1/tenants/{organizationName}/{organizationUnitName}/{name}` for deleting a Tenant.
6. THE Server SHALL use Go 1.21-compatible routing: register `/v1/tenants` and `/v1/tenants/` patterns. THE Server SHALL parse item paths using `strings.TrimPrefix` and `strings.Split`, requiring exactly three non-empty path segments: `organizationName`, `organizationUnitName`, and `name`. If the path has fewer or more than three segments, THE Server SHALL return HTTP 404. If exactly three segments exist but any segment fails DNS-label validation, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`.

---

### Requirement 6: Create Tenant — POST /v1/tenants

**User Story:** As a platform operator, I want to create a Tenant via REST API.

#### Acceptance Criteria

1. WHEN a valid POST is received, THE Server SHALL validate, verify parent exists, store the Tenant, and return HTTP 201 with full resource including server-set `apiVersion`, `kind`, and `status.phase = "Active"`. THE Server SHALL store only on the successful path that produces HTTP 201.
2. WHEN the composite key already exists, THE Server SHALL return HTTP 409 with `error.code = "RESOURCE_ALREADY_EXISTS"`.
3. WHEN field validation fails, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field` identifying the first invalid field.
4. WHEN the parent OrganizationUnit does not exist, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`, `error.field = "spec.organizationUnitName"`, and `error.message` identifying the missing parent reference `spec.organizationName/spec.organizationUnitName`.
5. WHEN the body is not valid JSON, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`.
6. WHEN `Content-Type` is not `application/json`, THE Server SHALL return HTTP 415 with `error.code = "VALIDATION_FAILED"`.
7. IF the body contains the key `status`, THE Server SHALL return HTTP 400 with `error.field = "status"` without storing.
8. THE Server SHALL reject unknown fields via `json.Decoder.DisallowUnknownFields()`.

---

### Requirement 7: Get Tenant — GET /v1/tenants/{organizationName}/{organizationUnitName}/{name}

**User Story:** As a platform operator, I want to retrieve a Tenant by its composite key.

#### Acceptance Criteria

1. WHEN the Tenant exists, THE Server SHALL return HTTP 200 with the full resource shape.
2. WHEN the Tenant does not exist, THE Server SHALL return HTTP 404 with `error.code = "RESOURCE_NOT_FOUND"`.
3. IF any path segment is invalid (not DNS-label or >63 chars), THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` without a registry lookup. Path validation errors SHALL map to these public fields: invalid `organizationName` segment → `error.field = "spec.organizationName"`, invalid `organizationUnitName` segment → `error.field = "spec.organizationUnitName"`, invalid `name` segment → `error.field = "metadata.name"`.
4. THE Server SHALL return a deep copy; mutations by the caller SHALL NOT affect registry state.

---

### Requirement 8: List Tenants — GET /v1/tenants

**User Story:** As a platform operator, I want to list all Tenants.

#### Acceptance Criteria

1. THE Server SHALL return HTTP 200 with `{"items": [...]}`.
2. WHEN no Tenants are stored, THE Server SHALL return `{"items": []}`.
3. THE Server SHALL return Tenants sorted: `spec.organizationName` ascending, then `spec.organizationUnitName` ascending, then `metadata.name` ascending.
4. THE response top-level object SHALL contain only the `items` field.
5. On internal error, THE Server SHALL return HTTP 500 with `error.code = "INTERNAL_ERROR"`.

---

### Requirement 9: Update Tenant — PUT /v1/tenants/{organizationName}/{organizationUnitName}/{name}

**User Story:** As a platform operator, I want to update mutable Tenant fields.

#### Acceptance Criteria

1. WHEN a valid PUT is received and the Tenant exists, THE Server SHALL replace `metadata.displayName`, `metadata.labels`, `metadata.annotations`, and `spec.description` from the body, preserve `metadata.name`, `spec.organizationName`, `spec.organizationUnitName`, and `status` from stored entry, and return HTTP 200 with the updated Tenant.
2. WHEN the Tenant does not exist, THE Server SHALL return HTTP 404 with `error.code = "RESOURCE_NOT_FOUND"`.
3. THE Server SHALL require `metadata.name`, `spec.organizationName`, and `spec.organizationUnitName` to be present in the PUT body and matching the corresponding path segments. Absent, empty, or mismatched values SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and the appropriate `error.field`.
4. IF the body contains the key `status`, THE Server SHALL return HTTP 400 with `error.field = "status"`.
5. WHEN the body is not valid JSON, THE Server SHALL return HTTP 400 regardless of resource existence.
6. WHEN `Content-Type` is not `application/json`, THE Server SHALL return HTTP 415 regardless of resource existence.
7. THE registry update SHALL preserve stored `spec.organizationName` and `spec.organizationUnitName` and SHALL NOT move a Tenant between parents.
8. THE Server SHALL preserve or reset server-owned `apiVersion` and `kind` on update and SHALL NOT allow client input to change them.

---

### Requirement 10: Delete Tenant — DELETE /v1/tenants/{organizationName}/{organizationUnitName}/{name}

**User Story:** As a platform operator, I want to delete a Tenant.

#### Acceptance Criteria

1. WHEN the Tenant exists, THE Server SHALL remove it and return HTTP 204 with no body.
2. WHEN the Tenant does not exist, THE Server SHALL return HTTP 404 with `error.code = "RESOURCE_NOT_FOUND"`.
3. IF any path segment is invalid, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`.
4. THE delete handler SHALL include a placeholder: `// TODO(FEATURE-0005): emit Operation record — type: DeleteTenant`.
5. Future Project child-resource blockers are out of scope for FEATURE-0003.

---

### Requirement 11: OrganizationUnit Deletion Blocker Integration

**User Story:** As a platform operator, I want OrganizationUnit deletion to be blocked when Tenants reference it.

#### Acceptance Criteria

1. WHEN `DELETE /v1/organization-units/{orgName}/{ouName}` is requested and one or more Tenants reference that OrganizationUnit, THE Server SHALL return HTTP 409 with `error.code = "DELETE_BLOCKED"` and `error.message` identifying `"Tenant"` as the blocking kind.
2. THE Tenant blocker SHALL implement the same `ChildBlockerChecker` pattern used for Organization deletion blocking.
3. THE blocker SHALL query the TenantRegistry to count Tenants referencing the OrganizationUnit and return a blocking entry when count > 0.
4. WHEN the OrganizationUnit has zero Tenants referencing it, THE blocker SHALL return an empty blocking set and the delete SHALL proceed.
5. FEATURE-0003 SHALL extend OUHandler deletion behavior so that `DELETE /v1/organization-units/{organizationName}/{organizationUnitName}` consults an injected child blocker before deleting an OrganizationUnit. This is an intentional FEATURE-0003 change to FEATURE-0002's OrganizationUnit delete path, because Tenant child blockers were explicitly out of scope for FEATURE-0002 (Requirement 19.7 of FEATURE-0002). The design may modify the OUHandler interface if needed to inject the Tenant blocker.

---

### Requirement 12: In-Memory Thread-Safe Tenant Registry

**User Story:** As a developer, I want a thread-safe in-memory Tenant registry.

#### Acceptance Criteria

1. THE TenantRegistry SHALL store Tenants in a `map[string]resources.Tenant` protected by `sync.RWMutex`, using the composite key `organizationName/organizationUnitName/name`.
2. THE TenantRegistry SHALL use RLock for Get/List/Count and Lock for Create/Update/Delete.
3. THE TenantRegistry SHALL return deep copies on Get/List/Update (including Labels and Annotations maps).
4. `CreateTenant` SHALL return `(resources.Tenant, error)` — the stored deep copy.
5. `UpdateTenant` SHALL return `(resources.Tenant, error)` — the updated deep copy.
6. THE TenantRegistry SHALL accept `context.Context` as the first parameter on all public methods.
7. THE TenantRegistry SHALL be instantiable via constructor; no package-level global state.
8. THE TenantRegistry SHALL produce no data race reports under `go test -race` with 10+ concurrent goroutines.
9. THE TenantRegistry SHALL return typed sentinel errors distinguishing `RESOURCE_NOT_FOUND` from `RESOURCE_ALREADY_EXISTS`.
10. THE TenantRegistry SHALL return results sorted by `spec.organizationName`, then `spec.organizationUnitName`, then `metadata.name` (deterministic ordering).
11. THE TenantRegistry SHALL include `CountByOrganizationUnit(ctx, orgName, ouName) (int, error)` for deletion blocker use.
12. THE TenantRegistry SHALL remain storage-only and SHALL NOT depend on OrganizationUnitRegistry or OrganizationRegistry.

---

### Requirement 13: Safe JSON Decoding

**User Story:** As a security-conscious operator, I want safe JSON decoding for Tenant requests.

#### Acceptance Criteria

1. THE Server SHALL apply `http.MaxBytesReader` with 1 MiB limit inside the decode function.
2. THE Server SHALL use `json.Decoder.DisallowUnknownFields()`.
3. Body exceeding 1 MiB SHALL return HTTP 413. Size check takes priority over syntax errors.
4. Syntax, type-mismatch, or unknown-field errors SHALL return HTTP 400.
5. THE Server SHALL NOT echo raw request body in error responses.
6. THE `status` key SHALL be rejected before the typed decode.
7. HTTP 415 Unsupported Media Type is handled by `contentTypeMiddleware`, not by the Tenant safe decode function.

---

### Requirement 14: Request Context Usage

**User Story:** As a developer, I want registry methods to accept `context.Context`.

#### Acceptance Criteria

1. THE TenantRegistry methods SHALL each accept `context.Context` as their first parameter.
2. Pure validation functions MAY remain context-free.
3. THE Server SHALL pass `r.Context()` to all Tenant registry calls within HTTP handlers.
4. THE Server SHALL NOT store `context.Context` in any struct field.
5. THE Server SHALL NOT use `context.Background()` or `context.TODO()` inside handlers.

---

### Requirement 15: Operation Framework Boundary (Placeholder Only)

**User Story:** As an architect, I want operation boundary placeholders for future FEATURE-0005 wiring.

#### Acceptance Criteria

1. THE handlers for CreateTenant, UpdateTenant, and DeleteTenant SHALL include these exact placeholder comments:
   - `// TODO(FEATURE-0005): emit Operation record — type: CreateTenant`
   - `// TODO(FEATURE-0005): emit Operation record — type: UpdateTenant`
   - `// TODO(FEATURE-0005): emit Operation record — type: DeleteTenant`
2. THE placeholder SHALL NOT import or call any Operation package code.

---

### Requirement 16: Tests

**User Story:** As a developer, I want comprehensive tests for Tenant validation, registry, and handlers.

#### Acceptance Criteria

1. THE validation package SHALL include unit tests for: valid names accepted, empty/invalid/long names rejected, empty/invalid spec.organizationName rejected, empty/invalid spec.organizationUnitName rejected.
2. THE registry package SHALL include unit tests for: Create stores, duplicate key → ErrAlreadyExists, same name under different parents succeeds, Get by composite key, Get non-existent → ErrNotFound, List sorted, List empty, Update mutable fields, Update non-existent → ErrNotFound, Delete removes, Delete non-existent → ErrNotFound, CountByOrganizationUnit correct counts.
3. THE api package SHALL include HTTP handler tests for: POST 201 valid, POST 409 duplicate, POST 400 invalid fields, POST 400 non-existent parent, POST 400 status key, GET 200/404/400, GET list sorted/empty, PUT 200/404/400 (name mismatch, orgName mismatch, ouName mismatch, status), DELETE 204/404/400.
4. THE test suite SHALL include OU deletion blocker tests: delete OU with Tenants → 409, delete OU with zero Tenants → 204.
5. `go test -race ./...` with 10+ concurrent goroutines SHALL produce no race reports.
6. ALL tests SHALL be deterministic with no external dependencies.

---

### Requirement 17: Property-Based Tests

**User Story:** As a developer, I want property-based tests for Tenant validation and registry.

#### Acceptance Criteria

1. THE validation package SHALL include property tests using `testing/quick` that generate valid DNS-label names and verify acceptance.
2. THE validation package SHALL include property tests using `testing/quick` that generate arbitrary strings and verify rejection of invalid names.
3. THE registry package SHALL include a round-trip property: Create then Get returns equivalent resource.
4. THE registry package SHALL include a sort property: List returns items in correct composite order.
5. THE registry package SHALL include a deep-copy property: mutating returned maps does not affect stored state.
6. THE registry package SHALL include an idempotent-error property: duplicate Create returns ErrAlreadyExists and original entry is unchanged.

---

### Requirement 18: Concurrency Test

**User Story:** As a developer, I want a concurrency stress test for the Tenant registry.

#### Acceptance Criteria

1. THE test suite SHALL include a concurrency test with 10+ goroutines performing mixed Create/Get/List/Update/Delete/CountByOrganizationUnit operations.
2. WHEN run with `go test -race`, THE test SHALL produce no race reports.
3. THE test SHALL verify no panics occur and operations return expected results.

---

### Requirement 19: Non-Goals

**User Story:** As an architect, I want clear scope boundaries.

#### Acceptance Criteria

1. THE implementation SHALL NOT implement Project or nested Tenant hierarchies.
2. THE implementation SHALL NOT implement persistent storage or Kubernetes CRDs.
3. THE implementation SHALL NOT implement authentication, authorization, or RBAC.
4. THE implementation SHALL NOT implement the Operation framework beyond placeholders.
5. THE implementation SHALL NOT implement ServiceOps, plugin execution, or AI agent execution.
6. THE implementation SHALL NOT implement UI, SDE runtime transformation, or billing.
7. THE implementation SHALL NOT implement child-resource blockers on Tenant delete (future Project blockers are out of scope).

---

## Design Questions

1. **Parent lookup interface:** The design should clarify whether to inject an `OrganizationUnitLookup` interface (similar to `OrganizationLookup` in FEATURE-0002) into the TenantHandler for parent existence checks. Recommendation: use a narrow interface `OrganizationUnitLookup` that the existing `OrganizationUnitRegistry` satisfies.

2. **Tenant child blocker injection into OUHandler:** FEATURE-0002 Requirement 19.7 explicitly excluded Tenant child blockers from the OU delete path. FEATURE-0003 must now introduce this blocker. The design should define the cleanest way to inject a Tenant child blocker into OUHandler — whether by adding a `ChildBlockerChecker` parameter to `OUHandler` (mirroring how `OrgHandler` accepts one), or by another composition pattern.

3. **Three-segment path parsing:** Resolved — Requirement 5.6 specifies exact three-segment parsing using `strings.Split` with `len(parts) == 3` check, consistent with FEATURE-0002's two-segment approach.
