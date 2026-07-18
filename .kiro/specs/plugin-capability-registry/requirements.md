# Requirements Document

## Introduction

FEATURE-0007 introduces `Plugin` and `Capability` as Phase 1 plugin registry primitives in
Sovrunn. A `Plugin` declares an implementation unit that performs lifecycle operations for a
service family or provider. A `Capability` declares a specific lifecycle action that a plugin
supports for a given ServiceClass.

These resources record **what a plugin claims to support**. They do NOT execute plugins, load
dynamic code, call external runtimes, or provision infrastructure. Plugin execution, ServiceOps
orchestration, and runtime plugin loading are future capabilities.

**Scope:** Plugin and Capability are **global platform registry resources**. They are NOT scoped
to Organization, OrganizationUnit, Tenant, or Project. A Plugin is identified by `metadata.name`
alone. A Capability is identified by `metadata.name` alone.

**Relationship to catalog:** A Plugin declares `serviceClassRefs` referencing existing
ServiceClasses (FEATURE-0006). A Capability references a Plugin via `pluginRef` and a ServiceClass
via `serviceClassRef`. This establishes:

```text
Plugin
  -> declares support for one or more ServiceClasses
  -> Capability entries declare individual lifecycle operations
```

This feature depends on FEATURE-0001 through FEATURE-0006. It reuses the existing project
skeleton, in-memory registry patterns, safe JSON decoding, structured errors, and the FEATURE-0005
Operation emission mechanism.

## Glossary

- **Plugin**: Implementation unit that performs lifecycle operations for a service family or provider. Identity: `metadata.name`.
- **Capability**: Declared lifecycle action supported by a plugin for a specific ServiceClass. Identity: `metadata.name`.
- **PluginType**: Classification of the plugin family. Allowed values: `dStoreOps`, `cacheOps`, `streamOps`, `objectOps`, `gatewayOps`, `faasOps`, `lbOps`, `k8sOps`, `bigDataOps`, `sdeOps`. Note: `k8sOps` is accepted for Phase 1; `docs/glossary.md` will be updated in this feature to include it.
- **DeploymentMode**: How the plugin executes at runtime. Allowed Phase 1 value: `compiled-in`. Future values: `sidecar`, `remote`.
- **CapabilityOperation**: Lifecycle action a capability declares. Allowed values: `Validate`, `Plan`, `Provision`, `Configure`, `Bind`, `Observe`, `Scale`, `Upgrade`, `Backup`, `Restore`, `RotateCredentials`, `Unbind`, `Delete`.
- **serviceClassRefs**: List of ServiceClass names a Plugin claims to support.
- **pluginRef**: Reference from a Capability to its parent Plugin.
- **serviceClassRef**: Reference from a Capability to the ServiceClass the operation applies to.
- **Registry**: In-memory, thread-safe store backed by `sync.RWMutex`, consistent with existing registries.
- **Error Code**: Stable string from: `VALIDATION_FAILED`, `RESOURCE_NOT_FOUND`, `RESOURCE_ALREADY_EXISTS`, `DELETE_BLOCKED`, `METHOD_NOT_ALLOWED`, `INTERNAL_ERROR`.
- **OperationEmitter**: The FEATURE-0005 mechanism for recording lifecycle actions.
- **DNS-label name**: Matches `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`, 1–63 chars.

## Requirements

---

### Requirement 1: Plugin Resource Shape

**User Story:** As a platform operator, I want Plugin to follow the canonical `metadata/spec/status` shape, so that it is consistent with all Phase 1 resources.

#### Acceptance Criteria

1. THE `Plugin` struct SHALL include `APIVersion`, `Kind`, `Metadata`, `Spec`, and `Status` with JSON tags `apiVersion`, `kind`, `metadata`, `spec`, `status`.
2. THE Server SHALL set `apiVersion` to `"platform.sovrunn.io/v1alpha1"` and `kind` to `"Plugin"` on all successful responses regardless of client input.
3. THE `PluginSpec` SHALL include `PluginType` (`json:"pluginType"`), `Version` (`json:"version"`), `ServiceClassRefs` (`json:"serviceClassRefs"`, `[]string`), `DeploymentMode` (`json:"deploymentMode"`), `Description` (`json:"description,omitempty"`), and `Tags` (`json:"tags,omitempty"`, `[]string`).
4. THE `PluginStatus` SHALL include `Phase` (`json:"phase"`) and `Message` (`json:"message,omitempty"`).
5. THE `metadata.name` SHALL be the Plugin name (DNS-label).
6. IF the top-level request body contains the key `status` — regardless of value — THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "status"`.

---

### Requirement 2: Capability Resource Shape

**User Story:** As a platform operator, I want Capability to follow the canonical `metadata/spec/status` shape and reference its parent Plugin and ServiceClass.

#### Acceptance Criteria

1. THE `Capability` struct SHALL include `APIVersion`, `Kind`, `Metadata`, `Spec`, and `Status` with the standard JSON tags.
2. THE Server SHALL set `apiVersion` to `"platform.sovrunn.io/v1alpha1"` and `kind` to `"Capability"` on all successful responses.
3. THE `CapabilitySpec` SHALL include `PluginRef` (`json:"pluginRef"`), `ServiceClassRef` (`json:"serviceClassRef"`), `Operation` (`json:"operation"`), `Supported` (`json:"supported"`, `bool`), and `Description` (`json:"description,omitempty"`).
4. THE `CapabilityStatus` SHALL include `Phase` (`json:"phase"`) and `Message` (`json:"message,omitempty"`).
5. THE `metadata.name` SHALL be the Capability name (DNS-label).
6. IF the top-level request body contains the key `status` — regardless of value — THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "status"`.

---

### Requirement 3: Plugin Validation

**User Story:** As a platform operator, I want deterministic Plugin validation, so that only well-formed plugin entries are stored.

#### Acceptance Criteria

1. IF `metadata.name` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"`.
2. THE `spec.pluginType` SHALL be one of: `dStoreOps`, `cacheOps`, `streamOps`, `objectOps`, `gatewayOps`, `faasOps`, `lbOps`, `k8sOps`, `bigDataOps`, `sdeOps`. Any other value → `FieldError` with `Field = "spec.pluginType"`.
3. THE `spec.version` SHALL be required and non-empty. Empty or absent → `FieldError` with `Field = "spec.version"`.
4. THE `spec.serviceClassRefs` SHALL be required, non-nil, and contain at least one entry. Empty or nil → `FieldError` with `Field = "spec.serviceClassRefs"`.
5. EACH entry in `spec.serviceClassRefs` SHALL be a valid DNS-label (1–63 chars). Invalid entry → `FieldError` with `Field = "spec.serviceClassRefs"`.
6. THE `spec.deploymentMode` SHALL be one of: `compiled-in`. Any other value → `FieldError` with `Field = "spec.deploymentMode"`. Future values (`sidecar`, `remote`) are not accepted in Phase 1.
7. `spec.description` and `spec.tags` are optional and not format-validated beyond basic type correctness.
8. THE Validator SHALL be a pure function (context-free, no I/O, no registry lookup).

---

### Requirement 4: Capability Validation

**User Story:** As a platform operator, I want deterministic Capability validation.

#### Acceptance Criteria

1. IF `metadata.name` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"`.
2. IF `spec.pluginRef` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.pluginRef"`.
3. IF `spec.serviceClassRef` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.serviceClassRef"`.
4. THE `spec.operation` SHALL be one of: `Validate`, `Plan`, `Provision`, `Configure`, `Bind`, `Observe`, `Scale`, `Upgrade`, `Backup`, `Restore`, `RotateCredentials`, `Unbind`, `Delete`. Any other value → `FieldError` with `Field = "spec.operation"`.
5. `spec.pluginRef`, `spec.serviceClassRef`, and `spec.operation` SHALL be required.
6. `spec.supported` defaults to `false` if absent (Go zero-value); no validation error for absent boolean.
7. `spec.description` is optional and not format-validated.
8. THE Validator SHALL be a pure function (context-free, no I/O, no registry lookup).

---

### Requirement 5: In-Memory Plugin Registry

**User Story:** As a developer, I want a thread-safe in-memory Plugin registry consistent with existing registries.

#### Acceptance Criteria

1. THE PluginRegistry SHALL store entries in a `map[string]resources.Plugin` protected by `sync.RWMutex`, keyed by `metadata.name`.
2. THE registry SHALL return deep copies on Create/Get/List/Update (including ServiceClassRefs slice, Tags slice).
3. `CreatePlugin` and `UpdatePlugin` SHALL return `(resources.Plugin, error)`.
4. `UpdatePlugin` SHALL preserve `metadata.name`, `status`, `apiVersion`, `kind` from the stored entry; replace only mutable fields.
5. THE registry SHALL accept `context.Context` as the first parameter on all public methods; no package-level global state.
6. THE registry SHALL return sentinel errors distinguishing `RESOURCE_NOT_FOUND` from `RESOURCE_ALREADY_EXISTS`.
7. `ListPlugins` SHALL return a non-nil slice sorted by `metadata.name` ascending.
8. THE registry SHALL be storage-only; no dependency on other registries.
9. THE registry SHALL produce no data race reports under `go test -race` with 10+ goroutines.

---

### Requirement 6: In-Memory Capability Registry

**User Story:** As a developer, I want a thread-safe in-memory Capability registry with filtering support.

#### Acceptance Criteria

1. THE CapabilityRegistry SHALL store entries in a `map[string]resources.Capability` protected by `sync.RWMutex`, keyed by `metadata.name`.
2. THE registry SHALL return deep copies on Create/Get/List.
3. `CreateCapability` SHALL return `(resources.Capability, error)`.
4. `DeleteCapability(ctx, name)` SHALL return `error` only. This matches prior registry patterns and the HTTP 204 no-body response.
5. THE registry SHALL accept `context.Context` as the first parameter on all public methods; no package-level global state.
6. THE registry SHALL return sentinel errors distinguishing `RESOURCE_NOT_FOUND` from `RESOURCE_ALREADY_EXISTS`.
7. `ListCapabilities` SHALL return a non-nil slice sorted by `metadata.name` ascending.
8. `ListCapabilities` SHALL accept optional filters: `pluginRef` and `serviceClassRef`. When both are provided, results match both (AND logic). When neither is provided, all entries are returned.
9. THE registry SHALL include `CountByPlugin(ctx, pluginName) (int, error)` for delete-blocker use.
10. THE registry SHALL be storage-only; no dependency on other registries.
11. THE registry SHALL produce no data race reports under `go test -race` with 10+ goroutines.

---

### Requirement 7: Plugin REST API

**User Story:** As a platform operator, I want CRUD endpoints for Plugin.

#### Acceptance Criteria

1. THE Server SHALL register `POST /v1/plugins`, `GET /v1/plugins`, `GET/PUT/DELETE /v1/plugins/{name}`.
2. THE Server SHALL use Go 1.21-compatible routing: `/v1/plugins` and `/v1/plugins/` patterns; item path has exactly ONE non-empty segment. Wrong segment count → HTTP 404.
3. IF the `{name}` segment is not a valid DNS-label, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "metadata.name"` without a registry lookup.
4. `POST` valid → 201 with full resource, `status.phase = "Active"`, server-set apiVersion/kind. Duplicate name → 409 `RESOURCE_ALREADY_EXISTS`.
5. `GET` item exists → 200; missing → 404 `RESOURCE_NOT_FOUND`. `GET` list → `{"items": [...]}` sorted by name; empty → `{"items": []}`.
6. `PUT` → 200 on success; missing → 404. `DELETE` → 204 no body; missing → 404.
7. Safe JSON decoding (1 MiB limit, DisallowUnknownFields, status rejection); 415 handled by contentTypeMiddleware; bad JSON/oversized/unknown-field per existing decoder patterns.

---

### Requirement 8: Capability REST API

**User Story:** As a platform operator, I want CRUD endpoints for Capability with filtering.

#### Acceptance Criteria

1. THE Server SHALL register `POST /v1/capabilities`, `GET /v1/capabilities`, `GET/DELETE /v1/capabilities/{name}`.
2. THE Server SHALL use Go 1.21-compatible routing: `/v1/capabilities` and `/v1/capabilities/` patterns; item path has exactly ONE non-empty segment. Wrong segment count → HTTP 404.
3. IF the `{name}` segment is not a valid DNS-label, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "metadata.name"` without a registry lookup.
4. `POST` valid → 201 with full resource, `status.phase = "Active"`, server-set apiVersion/kind. Duplicate name → 409 `RESOURCE_ALREADY_EXISTS`.
5. `GET` item exists → 200; missing → 404 `RESOURCE_NOT_FOUND`.
6. `GET` list SHALL accept optional query parameters `pluginRef` and `serviceClassRef` for filtering. Result → `{"items": [...]}` sorted by name; empty → `{"items": []}`.
7. `DELETE` → 204 no body; missing → 404.
8. Capability does NOT support `PUT` (update). Capabilities are immutable after creation; to change, delete and recreate.
9. IF `PUT /v1/capabilities/{name}` is received, THE Server SHALL return HTTP 405 Method Not Allowed with `error.code = "METHOD_NOT_ALLOWED"` and `error.message` indicating that Capability does not support update, regardless of whether a Capability with that name exists.
10. Safe JSON decoding (1 MiB limit, DisallowUnknownFields, status rejection); 415 handled by contentTypeMiddleware; bad JSON/oversized/unknown-field per existing decoder patterns.

---

### Requirement 9: Reference Validation (Plugin → ServiceClass)

**User Story:** As a platform operator, I want Plugin serviceClassRefs validated against existing ServiceClasses.

#### Acceptance Criteria

1. WHEN a valid `POST /v1/plugins` request passes field validation, THE Server SHALL verify each entry in `spec.serviceClassRefs` exists as a registered ServiceClass.
2. IF any referenced ServiceClass does not exist, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "spec.serviceClassRefs"` and a message identifying the missing ServiceClass name(s).
3. THE reference check SHALL occur in the API/service layer via a narrow lookup interface; the PluginRegistry SHALL NOT depend on the ServiceClassRegistry.
4. WHEN a `PUT` update passes validation and identity checks, THE Server SHALL verify all `spec.serviceClassRefs` still exist; if any are absent, THE Server SHALL return HTTP 400 `VALIDATION_FAILED` with `error.field = "spec.serviceClassRefs"`.

---

### Requirement 10: Reference Validation (Capability → Plugin and ServiceClass)

**User Story:** As a platform operator, I want Capability references validated against existing Plugin and ServiceClass resources.

#### Acceptance Criteria

1. WHEN a valid `POST /v1/capabilities` request passes field validation, THE Server SHALL verify the Plugin named `spec.pluginRef` exists.
2. IF the referenced Plugin does not exist, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "spec.pluginRef"` and a message identifying the missing Plugin.
3. WHEN a valid `POST /v1/capabilities` request passes field validation, THE Server SHALL verify the ServiceClass named `spec.serviceClassRef` exists.
4. IF the referenced ServiceClass does not exist, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "spec.serviceClassRef"` and a message identifying the missing ServiceClass.
5. THE reference checks SHALL occur in the API/service layer via narrow lookup interfaces; the CapabilityRegistry SHALL NOT depend on other registries.
6. WHEN a valid `POST /v1/capabilities` request passes field validation and Plugin reference validation, THE Server SHALL verify that `spec.serviceClassRef` is present in the referenced Plugin's `spec.serviceClassRefs`. IF the ServiceClass is NOT in the Plugin's `serviceClassRefs`, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"`, `error.field = "spec.serviceClassRef"`, and a message stating that the ServiceClass is not declared by the referenced Plugin.

---

### Requirement 11: Plugin Delete Blocking

**User Story:** As a platform operator, I want Plugin deletion blocked while Capabilities still reference it.

#### Acceptance Criteria

1. WHEN `DELETE /v1/plugins/{name}` is requested and one or more Capabilities have `spec.pluginRef == name`, THE Server SHALL return HTTP 409 with `error.code = "DELETE_BLOCKED"` and a message identifying `"Capability"` as the blocking kind.
2. THE blocker SHALL query the CapabilityRegistry via `CountByPlugin` and block when count > 0.
3. WHEN the Plugin has zero Capabilities, THE delete SHALL proceed and return 204.
4. THE blocker SHALL use a narrow blocker interface consistent with existing patterns; no generic blocker framework.

---

### Requirement 12: Update Identity Consistency (Plugin)

**User Story:** As a platform operator, I want PUT requests to enforce path/body identity matching for Plugin.

#### Acceptance Criteria

1. FOR `PUT /v1/plugins/{name}`, THE Server SHALL require body `metadata.name` present and equal to the path `{name}`. Absent, empty, or mismatched → HTTP 400 `VALIDATION_FAILED`, `error.field = "metadata.name"`.
2. THE registry update SHALL preserve immutable identity fields and SHALL NOT rename a Plugin.
3. THE Server SHALL preserve or reset server-owned `apiVersion`, `kind`, and `status` on update; client input SHALL NOT change server-owned or immutable fields.
4. Mutable fields include: `metadata.labels`, `metadata.annotations`, `spec.pluginType`, `spec.version`, `spec.serviceClassRefs`, `spec.deploymentMode`, `spec.description`, `spec.tags`.

---

### Requirement 13: Operation Emission (FEATURE-0005 Integration)

**User Story:** As a platform operator, I want plugin registry lifecycle actions recorded as Operations.

#### Acceptance Criteria

1. AFTER a successful Plugin create/update/delete, THE Server SHALL emit an Operation of type `CreatePlugin`, `UpdatePlugin`, or `DeletePlugin` respectively, with `resourceKind = "Plugin"`.
2. AFTER a successful Capability create/delete, THE Server SHALL emit an Operation of type `CreateCapability` or `DeleteCapability` respectively, with `resourceKind = "Capability"`.
3. THE resource kind constants (`PluginKind`, `CapabilityKind`) SHALL be added to `internal/resources` following the existing pattern.
4. FEATURE-0007 SHALL extend the FEATURE-0005 operation type constants with: `CreatePlugin`, `UpdatePlugin`, `DeletePlugin`, `CreateCapability`, `DeleteCapability`.
5. FEATURE-0007 SHALL extend `OperationSpec` with two optional plugin-reference fields: `pluginName` (`json:"pluginName,omitempty"`) and `capabilityName` (`json:"capabilityName,omitempty"`). These field names are LOCKED (matching the established pattern of `serviceClassName`/`servicePlanName` from FEATURE-0006). They are used only for Plugin and Capability Operation records and remain empty/omitted for other resource kinds.
6. Plugin Operation records SHALL set `pluginName` to the plugin name. Capability Operation records SHALL set `pluginName` to the referenced plugin name and `capabilityName` to the capability name.
7. Emission SHALL use the nil-safe FEATURE-0005 emitter; emission failure SHALL NOT affect the primary API response.
8. THE Server SHALL NOT emit an Operation on failed validation, duplicate create, missing reference, not-found, or delete-blocked cases.

---

### Requirement 14: Security and Privacy

**User Story:** As a security-conscious operator, I want plugin registry resources to never store secrets.

#### Acceptance Criteria

1. Plugin and Capability SHALL NOT store secrets, tokens, credentials, or passwords.
2. Plugin `spec.serviceClassRefs` and Capability `spec.pluginRef`/`spec.serviceClassRef` are references only — no credential material.
3. THE Server SHALL NOT store raw request bodies and SHALL NOT echo raw bodies in error responses.
4. THE Server SHALL NOT log secrets or raw bodies.
5. Plugin `spec.description` and Capability `spec.description` SHALL NOT be used to store credential-bearing content; no format validation is enforced in Phase 1, but documentation SHALL warn against it.

---

### Requirement 15: Tests

**User Story:** As a developer, I want comprehensive tests for Plugin and Capability resources.

#### Acceptance Criteria

1. Validation unit tests: valid names accepted; invalid/empty/long names rejected; invalid pluginType rejected; invalid deploymentMode rejected; invalid operation rejected; empty serviceClassRefs rejected; invalid serviceClassRef entries rejected; missing required fields rejected.
2. Registry unit tests (both registries): Create stores; duplicate → ErrAlreadyExists (original unchanged); Get by key; missing → ErrNotFound; List sorted; empty → non-nil `[]`; Plugin Update mutable fields only; Update missing → ErrNotFound; Delete removes; Delete missing → ErrNotFound; Capability CountByPlugin correct; Capability List with filters (pluginRef, serviceClassRef, both, neither).
3. Handler tests (both resources): POST 201/409/400 (invalid fields, status key, bad JSON, unknown field); POST 413 oversized; Plugin POST 400 missing ServiceClass ref; Capability POST 400 missing Plugin ref; Capability POST 400 missing ServiceClass ref; Capability POST 400 ServiceClass not in Plugin serviceClassRefs; GET 200/404/400; wrong path shape → 404; list sorted/empty; list with query filters; Plugin PUT 200/404/400 (identity mismatch); DELETE 204/404; Capability PUT → 405 METHOD_NOT_ALLOWED.
4. Delete-blocking tests: Plugin delete with Capabilities → 409 DELETE_BLOCKED; with zero Capabilities → 204.
5. Operation emission tests: successful create/update/delete of Plugin and create/delete of Capability records the correct Operation type and resource kind; failed actions emit nothing; emission failure does not change the primary response.
6. `go test -race ./...` with 10+ goroutines produces no race reports.
7. ALL tests deterministic; no external dependencies.

---

### Requirement 16: Property-Based Tests

**User Story:** As a developer, I want property tests for validation and registries.

#### Acceptance Criteria

1. Validation package property tests (`testing/quick`, `Config{MaxCount: 100}`): valid DNS-label names accepted; arbitrary invalid strings rejected; valid enum values accepted; invalid enum values rejected.
2. Registry property tests: Create/Get round-trip preserves data; List sort invariant; deep-copy immutability; duplicate-create idempotent error (original unchanged); Capability filter correctness.
3. Each property test tagged `// Feature: plugin-capability-registry, Property N: <title>`.

---

### Requirement 17: Non-Goals

**User Story:** As an architect, I want clear scope boundaries.

#### Acceptance Criteria

1. NO plugin execution, runtime loading, or dynamic plugin discovery.
2. NO ServiceOps orchestration, workflow engine, or plugin lifecycle management.
3. NO remote plugin communication (gRPC, HTTP callbacks, webhooks).
4. NO plugin marketplace, versioning rules, or compatibility matrix enforcement.
5. NO ServiceInstance or ServiceBinding (FEATURE-0008).
6. NO Kubernetes CRDs, operators, or GitOps controllers.
7. NO authentication, authorization, RBAC, or tenant-scoped plugin visibility.
8. NO async workflows, approval flows, queues, or background workers.
9. NO persistence/database storage, UI, or billing.
10. NO Go 1.22 wildcard routing; no new external dependencies.
11. NO conformance testing framework or plugin certification.
12. NO capability matching or resolution for ServiceInstance provisioning.

---

### Requirement 18: Edge Cases

**User Story:** As a developer, I want edge-case behavior defined.

#### Acceptance Criteria

1. WHEN a mutating action fails (validation, duplicate, missing reference, not-found, delete-blocked), THE Server SHALL NOT emit an Operation.
2. WHEN a Plugin item path has more than one segment, THE Server SHALL return HTTP 404.
3. WHEN a Capability item path has more than one segment, THE Server SHALL return HTTP 404.
4. WHEN `spec.serviceClassRefs` contains duplicate ServiceClass names, THE Server SHALL accept the request (duplicates are harmless metadata; no deduplication enforced in Phase 1).
5. WHEN a Plugin's `spec.serviceClassRefs` references a ServiceClass that is later deleted, THE Plugin remains valid (referential integrity is checked at write time only, not continuously enforced).
6. WHEN a Capability's `spec.pluginRef` references a Plugin, the referenced Plugin cannot be deleted via the API while the Capability exists (Requirement 11 blocks the delete). Therefore, under normal API operations, a Capability's `spec.pluginRef` always points to an existing Plugin. Referential integrity is enforced at write time (Capability creation) and protected at delete time (Plugin deletion blocked). No background reconciliation is performed.
7. WHEN a nil emitter is used (isolated handler tests), emission SHALL be skipped gracefully without panic.
8. WHEN `PUT /v1/capabilities/{name}` is received, THE Server SHALL return HTTP 405 Method Not Allowed with `error.code = "METHOD_NOT_ALLOWED"` (Capability does not support update).
9. WHEN the same `metadata.name` is used for both a Plugin and a Capability, THE Server SHALL allow it (Plugin and Capability are distinct resource kinds with separate registries).
10. WHEN `spec.supported` is `false` on a Capability, THE Server SHALL accept and store it (a declared-unsupported capability is still valid registry metadata).

---

## Compatibility with Completed Phase 1 Features

### FEATURE-0001 through FEATURE-0004 (Organization, OrganizationUnit, Tenant, Project)

Plugin and Capability are global platform resources. They do not reference or depend on the
governance hierarchy. No Organization, OrganizationUnit, Tenant, or Project scoping is required.

### FEATURE-0005 (Operation Resource)

FEATURE-0007 reuses the FEATURE-0005 `OperationEmitter` interface and nil-safe `emitOperation`
helper. Plugin and Capability handlers receive the emitter via constructor injection (nil-safe),
consistent with all other resource handlers. Five new operation type constants
(`CreatePlugin`, `UpdatePlugin`, `DeletePlugin`, `CreateCapability`, `DeleteCapability`) and two
resource kind constants (`PluginKind`, `CapabilityKind`) are added to `internal/resources`.
Emission occurs only after a successful mutating action, uses the request context for the
request ID, and never alters the primary response.

### FEATURE-0006 (ServiceClass and ServicePlan)

Plugin `spec.serviceClassRefs` and Capability `spec.serviceClassRef` reference ServiceClass
resources registered in the FEATURE-0006 ServiceClassRegistry. Reference validation uses a narrow
`ServiceClassLookup` interface injected into the handlers, not direct registry dependency. This
maintains the existing dependency-injection and interface-segregation patterns.

## Design Questions

> These questions are OPEN and MUST be resolved in `design.md` before implementation.

1. **ServiceClass lookup interface:** Define a narrow `ServiceClassLookup` interface (consistent with existing lookup patterns like `OrganizationLookup`, `TenantLookup`) injected into Plugin and Capability handlers. Confirm interface placement (`internal/api` vs `internal/registry`) and whether a single interface serves both Plugin and Capability handlers.

2. **Plugin lookup interface:** Define a narrow `PluginLookup` interface injected into the Capability handler for `spec.pluginRef` validation. Confirm whether this is a separate interface or combined with the delete-blocker interface.

3. **Delete-blocker interface shape:** Plugin delete takes a single path segment, so the blocker signature is `BlockedByPluginChildren(ctx, pluginName) ([]BlockedBy, error)` implemented by a Capability-backed checker. Confirm naming consistency with existing blockers (e.g., `BlockedByServiceClassChildren`).

4. **Capability immutability rationale:** The API contract shows no `PUT` for Capability. Confirm that delete-and-recreate is the correct update pattern, or whether a future phase may add `PUT`.

5. **Capability query parameter filtering:** Confirm query parameter names (`pluginRef`, `serviceClassRef`) and whether additional filters (e.g., `operation`) are needed for Phase 1 or deferred.

6. **ServiceClassRefs validation scope:** When updating a Plugin, all `spec.serviceClassRefs` are re-validated. Confirm whether only changed refs should be validated or the full list (full list is simpler and safer for Phase 1).

7. ~~**Operation spec fields:**~~ RESOLVED: `pluginName` and `capabilityName` are locked as the canonical field names, consistent with `serviceClassName`/`servicePlanName` from FEATURE-0006.

8. ~~**k8sOps plugin type:**~~ RESOLVED: `k8sOps` is accepted for Phase 1. The glossary (`docs/glossary.md` section 6) will be updated in this feature to include it.

