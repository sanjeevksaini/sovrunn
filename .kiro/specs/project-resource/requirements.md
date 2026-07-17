# Requirements Document

## Introduction

FEATURE-0004 implements the `Project` resource as the workload/environment grouping boundary
under a Tenant in Sovrunn. It is the fourth resource in the Phase 1 governance hierarchy:

```
Organization
  -> OrganizationUnit
      -> Tenant
          -> Project
              -> ServiceInstance (future)
```

Project depends on FEATURE-0001 (Organization), FEATURE-0002 (OrganizationUnit), and
FEATURE-0003 (Tenant) being fully implemented. Every Project belongs to exactly one Tenant,
identified by the composite parent reference `spec.organizationName` +
`spec.organizationUnitName` + `spec.tenantName`.

Identity is a four-part composite: `organizationName/organizationUnitName/tenantName/name`.
The same `metadata.name` may exist under different Tenants without conflict.

This feature covers only `Project`. It does not implement `ServiceInstance`, nested Project
hierarchies, persistent storage, Kubernetes CRDs, ServiceOps, the Operation framework
implementation, UI, AI agent execution, or SDE runtime transformation.

## Glossary

- **Project**: Workload/environment grouping boundary under a Tenant.
  Example: `prod` under Tenant `payments` in OrganizationUnit `ministry-health`, Organization `nic`.
- **Organization / OrganizationUnit / Tenant**: Parent governance resources (FEATURE-0001/0002/0003).
- **Registry**: In-memory, thread-safe store backed by `sync.RWMutex`-protected map.
- **Composite Key**: The identity key for Project: `organizationName/organizationUnitName/tenantName/name`.
- **DNS-label name**: Matches `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`. 1 to 63 characters.
- **APIError**: Structured JSON error body with `code`, `message`, optional `field`/`details`.
- **Error Code**: Stable string from: `VALIDATION_FAILED`, `RESOURCE_NOT_FOUND`, `RESOURCE_ALREADY_EXISTS`, `DELETE_BLOCKED`, `INTERNAL_ERROR`.
- **ProjectRegistry**: In-memory store for Project resources keyed by composite key.
- **ProjectHandler**: HTTP handler for Project CRUD endpoints.
- **TenantLookup**: Narrow interface for verifying parent Tenant existence in the API layer.
- **TenantChildBlocker**: Interface consulted by TenantHandler before deleting a Tenant; blocks deletion when Projects reference that Tenant.

## Requirements

---

### Requirement 1: Project Resource Shape

**User Story:** As a platform operator, I want the Project resource to follow the canonical `metadata/spec/status` shape, so that it is consistent with all Phase 1 resources.

#### Acceptance Criteria

1. THE `Project` struct SHALL include `APIVersion`, `Kind`, `Metadata`, `Spec`, and `Status` fields with JSON tags `apiVersion`, `kind`, `metadata`, `spec`, and `status` respectively.
2. THE `Metadata` struct SHALL use `json:"name"` (no omitempty) for `Name`, and `json:"displayName,omitempty"`, `json:"labels,omitempty"`, `json:"annotations,omitempty"` for optional fields.
3. THE `ProjectSpec` struct SHALL include `OrganizationName` (`json:"organizationName"`), `OrganizationUnitName` (`json:"organizationUnitName"`), `TenantName` (`json:"tenantName"`), and `Description` (`json:"description,omitempty"`).
4. THE `ProjectStatus` struct SHALL include `Phase` (`json:"phase"`) with valid values `"Active"`, `"Inactive"`, `"Deleting"`, `"Failed"`, and `Message` (`json:"message,omitempty"`).
5. WHEN the Server returns a successful Project response (HTTP 200 or 201), THE Server SHALL set `apiVersion` to `"platform.sovrunn.io/v1alpha1"` and `kind` to `"Project"` regardless of client input.
6. IF the top-level JSON request body on POST or PUT contains the key `status` — regardless of value (`{}`, `null`, `{"phase":""}`, or any other value) — THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "status"`.
7. WHEN `POST /v1/projects` succeeds with HTTP 201, THE Server SHALL set `status.phase` to `"Active"` both in the stored entry and in the response body. A successful POST means all validation, parent lookup, duplicate checks, and registry creation succeeded. If any error occurs, THE Server SHALL return an error and SHALL NOT store the Project.

---

### Requirement 2: Project Name Validation

**User Story:** As a platform operator, I want deterministic Project name validation, so that only safe, DNS-compatible names are accepted.

#### Acceptance Criteria

1. IF `metadata.name` is absent or empty, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"`.
2. IF `metadata.name` does not match `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"`.
3. IF `metadata.name` exceeds 63 characters, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"`.
4. IF `metadata.name` is valid (1–63 chars, DNS-label pattern), THEN THE Validator SHALL return zero errors for that field.
5. THE Validator SHALL be a pure function (context-free, no I/O, no registry lookup).

---

### Requirement 3: Parent Reference Validation

**User Story:** As a platform operator, I want `spec.organizationName`, `spec.organizationUnitName`, and `spec.tenantName` validated and referencing an existing Tenant, so that Projects cannot be created without a valid parent.

#### Acceptance Criteria

1. IF `spec.organizationName` is absent, empty, not DNS-label, or exceeds 63 characters, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.organizationName"`.
2. IF `spec.organizationUnitName` is absent, empty, not DNS-label, or exceeds 63 characters, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.organizationUnitName"`.
3. IF `spec.tenantName` is absent, empty, not DNS-label, or exceeds 63 characters, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.tenantName"`.
4. WHEN a valid `POST /v1/projects` request is received and field-level validation passes, THE Server SHALL verify that the Tenant identified by `spec.organizationName/spec.organizationUnitName/spec.tenantName` exists before storing the Project.
5. IF the parent Tenant does not exist, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`, `error.field = "spec.tenantName"`, and `error.message` identifying the missing parent reference `spec.organizationName/spec.organizationUnitName/spec.tenantName`.
6. `spec.organizationName`, `spec.organizationUnitName`, and `spec.tenantName` SHALL be immutable after creation. PUT requests must include all three fields matching the path segments.

---

### Requirement 4: Composite Identity and Uniqueness

**User Story:** As a platform operator, I want Project identity to be composite, so that the same project name can exist under different Tenants.

#### Acceptance Criteria

1. THE ProjectRegistry SHALL use the composite key `organizationName/organizationUnitName/tenantName/name` as the unique identity.
2. WHEN two Projects have the same `metadata.name` but different parent paths, THE ProjectRegistry SHALL store both without conflict.
3. WHEN a POST request has a composite key that already exists, THE Server SHALL return HTTP 409 with `error.code = "RESOURCE_ALREADY_EXISTS"`.
4. GET, PUT, and DELETE endpoints SHALL use path segments `{organizationName}/{organizationUnitName}/{tenantName}/{name}` to identify a Project.

---

### Requirement 5: REST API Endpoints and Path Parsing

**User Story:** As a platform operator, I want CRUD REST endpoints for Project.

#### Acceptance Criteria

1. THE Server SHALL register `POST /v1/projects` for creating a Project.
2. THE Server SHALL register `GET /v1/projects` for listing all Projects.
3. THE Server SHALL register `GET /v1/projects/{organizationName}/{organizationUnitName}/{tenantName}/{name}` for retrieving a Project.
4. THE Server SHALL register `PUT /v1/projects/{organizationName}/{organizationUnitName}/{tenantName}/{name}` for updating a Project.
5. THE Server SHALL register `DELETE /v1/projects/{organizationName}/{organizationUnitName}/{tenantName}/{name}` for deleting a Project.
6. THE Server SHALL use Go 1.21-compatible routing: register `/v1/projects` and `/v1/projects/` patterns. THE Server SHALL parse item paths using `strings.TrimPrefix` and `strings.Split`, requiring exactly four non-empty path segments: `organizationName`, `organizationUnitName`, `tenantName`, and `name`. If the path has fewer or more than four segments, or any segment is empty, THE Server SHALL return HTTP 404. If exactly four segments exist but any segment fails DNS-label validation, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`.
7. Path validation errors SHALL map to public fields: invalid `organizationName` → `spec.organizationName`, invalid `organizationUnitName` → `spec.organizationUnitName`, invalid `tenantName` → `spec.tenantName`, invalid `name` → `metadata.name`.
8. Unsupported methods on either route SHALL return HTTP 405.

---

### Requirement 6: Create Project — POST /v1/projects

**User Story:** As a platform operator, I want to create a Project via REST API.

#### Acceptance Criteria

1. WHEN a valid POST is received, THE Server SHALL validate the body first, verify the parent Tenant exists, store the Project, and return HTTP 201 with full resource including server-set `apiVersion`, `kind`, and `status.phase = "Active"`. THE Server SHALL store only on the successful path that produces HTTP 201.
2. WHEN the composite key already exists, THE Server SHALL return HTTP 409 with `error.code = "RESOURCE_ALREADY_EXISTS"`.
3. WHEN field validation fails, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field` identifying the first invalid field.
4. WHEN the parent Tenant does not exist, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`, `error.field = "spec.tenantName"`, and `error.message` identifying `spec.organizationName/spec.organizationUnitName/spec.tenantName`.
5. WHEN the body is not valid JSON, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`.
6. WHEN `Content-Type` is not `application/json`, THE Server SHALL return HTTP 415 with `error.code = "VALIDATION_FAILED"`.
7. IF the body contains the key `status`, THE Server SHALL return HTTP 400 with `error.field = "status"` without storing.
8. THE Server SHALL reject unknown fields via `json.Decoder.DisallowUnknownFields()`.

---

### Requirement 7: Get Project — GET /v1/projects/{organizationName}/{organizationUnitName}/{tenantName}/{name}

**User Story:** As a platform operator, I want to retrieve a Project by its composite key.

#### Acceptance Criteria

1. WHEN the Project exists, THE Server SHALL return HTTP 200 with the full resource shape.
2. WHEN the Project does not exist, THE Server SHALL return HTTP 404 with `error.code = "RESOURCE_NOT_FOUND"`.
3. IF any path segment is invalid (not DNS-label or >63 chars), THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` without a registry lookup, using the field mapping in Requirement 5.7.
4. THE Server SHALL return a deep copy; mutations by the caller SHALL NOT affect registry state.

---

### Requirement 8: List Projects — GET /v1/projects

**User Story:** As a platform operator, I want to list all Projects.

#### Acceptance Criteria

1. THE Server SHALL return HTTP 200 with `{"items": [...]}`.
2. WHEN no Projects are stored, THE Server SHALL return `{"items": []}` (non-nil empty slice).
3. THE Server SHALL return Projects sorted: `spec.organizationName` asc, then `spec.organizationUnitName` asc, then `spec.tenantName` asc, then `metadata.name` asc.
4. THE response top-level object SHALL contain only the `items` field.
5. On internal error, THE Server SHALL return HTTP 500 with `error.code = "INTERNAL_ERROR"`.

---

### Requirement 9: Update Project — PUT /v1/projects/{organizationName}/{organizationUnitName}/{tenantName}/{name}

**User Story:** As a platform operator, I want to update mutable Project fields.

#### Acceptance Criteria

1. WHEN a valid PUT is received and the Project exists, THE Server SHALL replace `metadata.displayName`, `metadata.labels`, `metadata.annotations`, and `spec.description` from the body, preserve `metadata.name`, `spec.organizationName`, `spec.organizationUnitName`, `spec.tenantName`, and `status` from the stored entry, and return HTTP 200 with the updated Project.
2. WHEN the Project does not exist, THE Server SHALL return HTTP 404 with `error.code = "RESOURCE_NOT_FOUND"`.
3. THE Server SHALL require `metadata.name`, `spec.organizationName`, `spec.organizationUnitName`, and `spec.tenantName` to be present in the PUT body and matching the corresponding path segments. Absent, empty, or mismatched values SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and the appropriate `error.field`.
4. `spec.organizationName`, `spec.organizationUnitName`, `spec.tenantName`, and `metadata.name` SHALL be immutable.
5. IF the body contains the key `status`, THE Server SHALL return HTTP 400 with `error.field = "status"`.
6. WHEN the body is not valid JSON, THE Server SHALL return HTTP 400 regardless of resource existence.
7. WHEN `Content-Type` is not `application/json`, THE Server SHALL return HTTP 415 regardless of resource existence.
8. THE Server SHALL preserve or reset server-owned `apiVersion` and `kind` on update and SHALL NOT allow client input to change server-owned or immutable fields.

---

### Requirement 10: Delete Project — DELETE /v1/projects/{organizationName}/{organizationUnitName}/{tenantName}/{name}

**User Story:** As a platform operator, I want to delete a Project.

#### Acceptance Criteria

1. WHEN the Project exists, THE Server SHALL remove it and return HTTP 204 with no body.
2. WHEN the Project does not exist, THE Server SHALL return HTTP 404 with `error.code = "RESOURCE_NOT_FOUND"`.
3. IF any path segment is invalid, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`.
4. THE delete handler SHALL include a placeholder: `// TODO(FEATURE-0005): emit Operation record — type: DeleteProject`.
5. Project delete child-resource blockers (future ServiceInstance blockers) are out of scope for FEATURE-0004.

---

### Requirement 11: Tenant Deletion Blocker Integration

**User Story:** As a platform operator, I want Tenant deletion to be blocked when Projects reference it.

#### Acceptance Criteria

1. WHEN `DELETE /v1/tenants/{organizationName}/{organizationUnitName}/{tenantName}` is requested and one or more Projects reference that Tenant, THE Server SHALL return HTTP 409 with `error.code = "DELETE_BLOCKED"` and `error.message` identifying `"Project"` as the blocking kind.
2. THE Project blocker SHALL implement a narrow blocker interface similar to `OUChildBlocker`. It SHALL NOT be a generic blocker framework.
3. THE blocker SHALL query the ProjectRegistry via `CountByTenant(ctx, orgName, ouName, tenantName)` and return a blocking entry when count > 0.
4. WHEN the Tenant has zero Projects referencing it, THE blocker SHALL return an empty blocking set and the delete SHALL proceed.
5. FEATURE-0004 SHALL extend TenantHandler deletion behavior so that `DELETE /v1/tenants/{organizationName}/{organizationUnitName}/{tenantName}` consults an injected Project child blocker before deleting a Tenant. This is an intentional FEATURE-0004 change to FEATURE-0003's Tenant delete path, because Project child blockers were out of scope for FEATURE-0003. The design may modify the TenantHandler interface if needed to inject the Project blocker.
6. TenantHandler SHALL allow the blocker to be nil. If nil, Tenant delete proceeds without child-blocker checks. Production FEATURE-0004 wiring SHALL inject the Project blocker.

---

### Requirement 12: In-Memory Thread-Safe Project Registry

**User Story:** As a developer, I want a thread-safe in-memory Project registry.

#### Acceptance Criteria

1. THE ProjectRegistry SHALL store Projects in a `map[string]resources.Project` protected by `sync.RWMutex`, using the composite key `organizationName/organizationUnitName/tenantName/name`.
2. THE ProjectRegistry SHALL use RLock for Get/List/Count and Lock for Create/Update/Delete.
3. THE ProjectRegistry SHALL return deep copies on Create/Get/List/Update (including Labels and Annotations maps).
4. `CreateProject` SHALL return `(resources.Project, error)` — the stored deep copy.
5. `UpdateProject` SHALL return `(resources.Project, error)` — derive the composite key from the submitted resource; preserve stored `APIVersion`, `Kind`, `Status`, `Metadata.Name`, `Spec.OrganizationName`, `Spec.OrganizationUnitName`, `Spec.TenantName`; replace only `Metadata.DisplayName`, `Metadata.Labels`, `Metadata.Annotations`, `Spec.Description`.
6. THE ProjectRegistry SHALL accept `context.Context` as the first parameter on all public methods.
7. THE ProjectRegistry SHALL be instantiable via constructor; no package-level global state.
8. THE ProjectRegistry SHALL produce no data race reports under `go test -race` with 10+ concurrent goroutines.
9. THE ProjectRegistry SHALL return typed sentinel errors distinguishing `RESOURCE_NOT_FOUND` from `RESOURCE_ALREADY_EXISTS` (reuse existing registry sentinel errors).
10. THE ProjectRegistry SHALL return results sorted by `spec.organizationName`, then `spec.organizationUnitName`, then `spec.tenantName`, then `metadata.name`.
11. THE ProjectRegistry SHALL include `CountByTenant(ctx, orgName, ouName, tenantName) (int, error)` for deletion blocker use.
12. THE ProjectRegistry SHALL remain storage-only and SHALL NOT depend on OrganizationRegistry, OrganizationUnitRegistry, or TenantRegistry.

---

### Requirement 13: Safe JSON Decoding

**User Story:** As a security-conscious operator, I want safe JSON decoding for Project requests.

#### Acceptance Criteria

1. THE Server SHALL apply `http.MaxBytesReader` with 1 MiB limit inside the decode function.
2. THE Server SHALL use `json.Decoder.DisallowUnknownFields()`.
3. Body exceeding 1 MiB SHALL return HTTP 413. Size check takes priority over syntax errors.
4. Syntax, type-mismatch, or unknown-field errors SHALL return HTTP 400.
5. THE Server SHALL NOT echo raw request body in error responses.
6. THE `status` key SHALL be rejected before the typed decode, regardless of value.
7. HTTP 415 Unsupported Media Type is handled by `contentTypeMiddleware`, not by the Project decode function.

---

### Requirement 14: Request Context Usage

**User Story:** As a developer, I want registry methods to accept `context.Context`.

#### Acceptance Criteria

1. THE ProjectRegistry methods SHALL each accept `context.Context` as their first parameter.
2. Pure validation functions MAY remain context-free.
3. THE Server SHALL pass `r.Context()` to all Project registry calls within HTTP handlers.
4. THE Server SHALL NOT store `context.Context` in any struct field.
5. THE Server SHALL NOT use `context.Background()` or `context.TODO()` inside handlers.

---

### Requirement 15: Operation Framework Boundary (Placeholder Only)

**User Story:** As an architect, I want operation boundary placeholders for future FEATURE-0005 wiring.

#### Acceptance Criteria

1. THE handlers for CreateProject, UpdateProject, and DeleteProject SHALL include these exact placeholder comments:
   - `// TODO(FEATURE-0005): emit Operation record — type: CreateProject`
   - `// TODO(FEATURE-0005): emit Operation record — type: UpdateProject`
   - `// TODO(FEATURE-0005): emit Operation record — type: DeleteProject`
2. THE placeholder SHALL NOT import or call any Operation package code.

---

### Requirement 16: Tests

**User Story:** As a developer, I want comprehensive tests for Project validation, registry, and handlers.

#### Acceptance Criteria

1. THE validation package SHALL include unit tests for: valid names accepted; empty/invalid/long metadata.name rejected; empty/invalid/long spec.organizationName, spec.organizationUnitName, spec.tenantName rejected; path validation field mapping correct.
2. THE registry package SHALL include unit tests for: Create stores; duplicate key → ErrAlreadyExists (original unchanged); same name under different Tenants succeeds; Get by composite key; Get non-existent → ErrNotFound; List sorted; empty List → non-nil empty slice; Update mutable fields only; Update non-existent → ErrNotFound; Delete removes; Delete non-existent → ErrNotFound; CountByTenant correct counts.
3. THE api package SHALL include HTTP handler tests for: POST 201 valid; POST 409 duplicate; POST 400 (invalid fields, non-existent parent, status key, bad JSON); POST 413 oversized body; GET 200/404/400; GET wrong path shape → 404; list sorted; list empty → []; PUT 200/404; PUT 400 (name mismatch, orgName mismatch, ouName mismatch, tenantName mismatch, status); DELETE 204/404/400.
4. THE test suite SHALL include Tenant deletion blocker tests: delete Tenant with Projects → 409; delete Tenant with zero Projects → 204.
5. `go test -race ./...` with 10+ concurrent goroutines SHALL produce no race reports.
6. ALL tests SHALL be deterministic with no external dependencies.

---

### Requirement 17: Property-Based Tests

**User Story:** As a developer, I want property-based tests for Project validation and registry.

#### Acceptance Criteria

1. THE validation package SHALL include property tests using `testing/quick` that generate valid DNS-label names and verify acceptance.
2. THE validation package SHALL include property tests using `testing/quick` that generate arbitrary strings and verify rejection of invalid names.
3. THE registry package SHALL include a round-trip property: Create then Get returns equivalent resource.
4. THE registry package SHALL include a sort property: List returns items in correct four-level composite order.
5. THE registry package SHALL include a deep-copy property: mutating returned maps does not affect stored state.
6. THE registry package SHALL include an idempotent-error property: duplicate Create returns ErrAlreadyExists and original entry is unchanged.

---

### Requirement 18: Concurrency Test

**User Story:** As a developer, I want a concurrency stress test for the Project registry.

#### Acceptance Criteria

1. THE test suite SHALL include a concurrency test with 10+ goroutines performing mixed Create/Get/List/Update/Delete/CountByTenant operations.
2. WHEN run with `go test -race`, THE test SHALL produce no race reports.
3. THE test SHALL verify no panics occur and operations return expected results.

---

### Requirement 19: Non-Goals

**User Story:** As an architect, I want clear scope boundaries.

#### Acceptance Criteria

1. THE implementation SHALL NOT implement ServiceInstance or nested Project hierarchies.
2. THE implementation SHALL NOT implement persistent storage or Kubernetes CRDs.
3. THE implementation SHALL NOT implement authentication, authorization, or RBAC.
4. THE implementation SHALL NOT implement the Operation framework beyond placeholder comments.
5. THE implementation SHALL NOT implement ServiceOps, plugin execution, or AI agent execution.
6. THE implementation SHALL NOT implement UI, SDE runtime transformation, or billing.
7. THE implementation SHALL NOT implement child-resource blockers on Project delete (future ServiceInstance blockers are out of scope).
8. THE implementation SHALL NOT introduce a generic blocker framework.
9. THE implementation SHALL NOT introduce Go 1.22 wildcard routing.

---

## Design Questions

1. **Parent lookup interface:** The design should define a narrow `TenantLookup` interface (mirroring `OrganizationUnitLookup` from FEATURE-0003) injected into ProjectHandler for parent existence checks. The existing `*TenantRegistry` should satisfy it via `GetTenant(ctx, orgName, ouName, name)`.

2. **Tenant child blocker injection into TenantHandler:** FEATURE-0003 Requirement 19.7 excluded Project child blockers from the Tenant delete path. FEATURE-0004 must introduce this blocker. The design should define the cleanest way to inject a Project child blocker into TenantHandler — e.g., adding a nil-safe blocker parameter to `NewTenantHandler`, mirroring how FEATURE-0003 added `OUChildBlocker` to `OUHandler`.

3. **Four-segment path parsing:** Resolved — Requirement 5.6 specifies exact four-segment parsing using `strings.Split` with `len(parts) == 4` check, consistent with FEATURE-0002/0003 patterns.
