# Requirements Document

## Introduction

FEATURE-0008 introduces `ServiceInstance` and `ServiceBinding` as Phase 1 service consumption
primitives in Sovrunn. A `ServiceInstance` represents a tenant/project-scoped requested service.
A `ServiceBinding` represents the consumption relationship between a consumer (application or
workload) and a ServiceInstance, exposing connection details or credentials.

These resources complete the Phase 1 service consumption model. They do NOT provision real
infrastructure, execute plugins, call external systems, or manage actual credentials. In Phase 1,
ServiceInstance creation validates references, stores the resource, records an Operation, and sets
status to `Ready`. ServiceBinding stores a stub secret reference.

**Scope:** ServiceInstance is scoped to `organizationRef + organizationUnitRef + tenantRef +
projectRef`. ServiceBinding references a ServiceInstance. Both are identified by `metadata.name`
which must be globally unique per resource kind in Phase 1 (consistent with all other Phase 1
resources).

**Relationship to governance hierarchy:** ServiceInstance references Organization (required),
OrganizationUnit (optional), Tenant (required), and Project (required). All references must
exist and be consistent with the governance hierarchy.

**Relationship to catalog:** ServiceInstance references a ServiceClass and a ServicePlan. The
ServicePlan must reference the specified ServiceClass. ServicePlan identity is composite
(`serviceClassName/name`) per FEATURE-0006; therefore, `spec.servicePlanRef` stores only the
plan name, and the combination of `spec.serviceClassRef` + `spec.servicePlanRef` forms the
fully-qualified plan reference used for validation and counters.

**Relationship to plugin registry:** At least one active Capability should exist for the
referenced ServiceClass; Phase 1 may warn (log) instead of blocking creation. The capability
check is encapsulated behind a narrow interface; implementation details of how capability
status is determined remain internal to the FEATURE-0007 registry.

**Relationship to operations:** All successful mutating actions emit Operation records via the
FEATURE-0005 OperationEmitter.

This feature depends on FEATURE-0001 through FEATURE-0007. It reuses the existing project
skeleton, in-memory registry patterns, safe JSON decoding, structured errors, and the
FEATURE-0005 Operation emission mechanism.

## Glossary

- **ServiceInstance**: A provisioned instance of a ServiceClass using a ServicePlan, scoped to a Project within the governance hierarchy. Identity: `metadata.name`.
- **ServiceBinding**: A binding that exposes connection details or credentials to a consumer for a specific ServiceInstance. Identity: `metadata.name`.
- **serviceClassRef**: Reference from a ServiceInstance to the ServiceClass being consumed. Value: plain `metadata.name` of the ServiceClass (DNS-label).
- **servicePlanRef**: Reference from a ServiceInstance to the ServicePlan being used. Value: plain `metadata.name` of the ServicePlan (DNS-label). The fully-qualified plan identity is the composite `serviceClassRef/servicePlanRef`; both fields together form the complete reference.
- **organizationRef**: Reference from a ServiceInstance to the owning Organization. Value: plain `metadata.name` of the Organization (DNS-label).
- **organizationUnitRef**: Optional reference from a ServiceInstance to the owning OrganizationUnit. Value: plain `metadata.name` of the OrganizationUnit (DNS-label).
- **tenantRef**: Reference from a ServiceInstance to the owning Tenant. Value: plain `metadata.name` of the Tenant (DNS-label). The fully-qualified tenant identity is determined by its governance hierarchy (org/ou/tenant); the ServiceInstance carries `organizationRef` + `organizationUnitRef` + `tenantRef` to disambiguate.
- **projectRef**: Reference from a ServiceInstance to the owning Project. Value: plain `metadata.name` of the Project (DNS-label). The fully-qualified project identity is determined by its governance hierarchy (org/ou/tenant/project); the ServiceInstance carries all ancestor refs to disambiguate.
- **serviceInstanceRef**: Reference from a ServiceBinding to the target ServiceInstance.
- **consumerRef**: Structured reference identifying the consumer of a ServiceBinding. Contains `kind` and `name`.
- **bindingType**: Classification of what a ServiceBinding exposes. Phase 1 allowed value: `credentials`. Future values may include `endpoint`, `config`.
- **secretRef**: Stub reference in ServiceBinding status indicating where credentials would be stored. Phase 1 uses a placeholder value.
- **parameters**: Optional key-value map on ServiceInstance spec for instance-specific configuration. Not validated against schema in Phase 1.
- **Registry**: In-memory, thread-safe store backed by `sync.RWMutex`, consistent with existing registries.
- **Error Code**: Stable string from: `VALIDATION_FAILED`, `RESOURCE_NOT_FOUND`, `RESOURCE_ALREADY_EXISTS`, `DELETE_BLOCKED`, `METHOD_NOT_ALLOWED`, `INTERNAL_ERROR`.
- **OperationEmitter**: The FEATURE-0005 mechanism for recording lifecycle actions.
- **DNS-label name**: Matches `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`, 1–63 chars.

## Requirements

---

### Requirement 1: ServiceInstance Resource Shape

**User Story:** As a platform operator, I want ServiceInstance to follow the canonical `metadata/spec/status` shape, so that it is consistent with all Phase 1 resources.

#### Acceptance Criteria

1. THE `ServiceInstance` struct SHALL include `APIVersion`, `Kind`, `Metadata`, `Spec`, and `Status` with JSON tags `apiVersion`, `kind`, `metadata`, `spec`, `status`.
2. THE Server SHALL set `apiVersion` to `"platform.sovrunn.io/v1alpha1"` and `kind` to `"ServiceInstance"` on all successful responses regardless of client input.
3. THE `ServiceInstanceSpec` SHALL include `OrganizationRef` (`json:"organizationRef"`), `OrganizationUnitRef` (`json:"organizationUnitRef,omitempty"`), `TenantRef` (`json:"tenantRef"`), `ProjectRef` (`json:"projectRef"`), `ServiceClassRef` (`json:"serviceClassRef"`), `ServicePlanRef` (`json:"servicePlanRef"`), and `Parameters` (`json:"parameters,omitempty"`, `map[string]string`).
4. THE `ServiceInstanceStatus` SHALL include `Phase` (`json:"phase"`) and `Message` (`json:"message,omitempty"`).
5. THE `metadata.name` SHALL be the ServiceInstance name (DNS-label).
6. IF the top-level request body contains the key `status` — regardless of value — on any mutating request (POST or PUT), THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "status"`.

---

### Requirement 2: ServiceBinding Resource Shape

**User Story:** As a platform operator, I want ServiceBinding to follow the canonical `metadata/spec/status` shape and reference its parent ServiceInstance.

#### Acceptance Criteria

1. THE `ServiceBinding` struct SHALL include `APIVersion`, `Kind`, `Metadata`, `Spec`, and `Status` with the standard JSON tags.
2. THE Server SHALL set `apiVersion` to `"platform.sovrunn.io/v1alpha1"` and `kind` to `"ServiceBinding"` on all successful responses.
3. THE `ServiceBindingSpec` SHALL include `ServiceInstanceRef` (`json:"serviceInstanceRef"`), `ConsumerRef` (`json:"consumerRef"`, struct with `Kind` and `Name` fields), and `BindingType` (`json:"bindingType"`).
4. THE `ConsumerRef` struct SHALL include `Kind` (`json:"kind"`) and `Name` (`json:"name"`).
5. THE `ServiceBindingStatus` SHALL include `Phase` (`json:"phase"`), `Message` (`json:"message,omitempty"`), and `SecretRef` (`json:"secretRef,omitempty"`).
6. THE `metadata.name` SHALL be the ServiceBinding name (DNS-label).
7. IF the top-level request body contains the key `status` — regardless of value — on POST, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "status"`. (ServiceBinding has no PUT, so this applies only to POST.)

---

### Requirement 3: ServiceInstance Validation

**User Story:** As a platform operator, I want deterministic ServiceInstance validation, so that only well-formed instances are stored.

#### Acceptance Criteria

1. IF `metadata.name` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"`.
2. IF `spec.organizationRef` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.organizationRef"`.
3. IF `spec.tenantRef` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.tenantRef"`.
4. IF `spec.projectRef` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.projectRef"`.
5. IF `spec.serviceClassRef` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.serviceClassRef"`.
6. IF `spec.servicePlanRef` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.servicePlanRef"`.
7. `spec.organizationUnitRef` is optional. IF present and not empty, it SHALL be validated as a DNS-label (1–63 chars). Invalid → `FieldError` with `Field = "spec.organizationUnitRef"`.
8. `spec.parameters` is optional. No key/value format validation in Phase 1 beyond basic type correctness (must be `map[string]string` if present).
9. THE Validator SHALL be a pure function (context-free, no I/O, no registry lookup).

---

### Requirement 4: ServiceBinding Validation

**User Story:** As a platform operator, I want deterministic ServiceBinding validation.

#### Acceptance Criteria

1. IF `metadata.name` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"`.
2. IF `spec.serviceInstanceRef` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.serviceInstanceRef"`.
3. IF `spec.consumerRef` is absent or nil, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.consumerRef"`.
4. IF `spec.consumerRef.kind` is absent or empty, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.consumerRef.kind"`.
5. IF `spec.consumerRef.name` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.consumerRef.name"`.
6. THE `spec.bindingType` SHALL be one of: `credentials`. Any other value → `FieldError` with `Field = "spec.bindingType"`. Future values (`endpoint`, `config`) are not accepted in Phase 1.
7. `spec.consumerRef.kind` is not enum-validated in Phase 1 (any non-empty string is accepted). Future phases may restrict to known kinds.
8. THE Validator SHALL be a pure function (context-free, no I/O, no registry lookup).

---

### Requirement 5: In-Memory ServiceInstance Registry

**User Story:** As a developer, I want a thread-safe in-memory ServiceInstance registry consistent with existing registries.

#### Acceptance Criteria

1. THE ServiceInstanceRegistry SHALL store entries in a `map[string]resources.ServiceInstance` protected by `sync.RWMutex`, keyed by `metadata.name`.
2. THE registry SHALL return deep copies on Create/Get/List/Update (including Parameters map, all string fields).
3. `CreateServiceInstance` and `UpdateServiceInstance` SHALL return `(resources.ServiceInstance, error)`.
4. `UpdateServiceInstance` SHALL preserve `metadata.name`, `status` (both `phase` and `message` from the stored entry unchanged), `apiVersion`, `kind` from the stored entry; replace only mutable spec fields and metadata annotations/labels.
5. THE registry SHALL accept `context.Context` as the first parameter on all public methods; no package-level global state.
6. THE registry SHALL return sentinel errors distinguishing `RESOURCE_NOT_FOUND` from `RESOURCE_ALREADY_EXISTS`.
7. `ListServiceInstances` SHALL return a non-nil slice sorted by `metadata.name` ascending.
8. `ListServiceInstances` SHALL accept optional filters: `tenantRef` and `projectRef`. When both are provided, results match both (AND logic). When neither is provided, all entries are returned. Filter values are plain `metadata.name` strings matching the corresponding `spec.tenantRef` and `spec.projectRef` fields stored on each ServiceInstance.
9. THE ServiceInstance registry SHALL be storage-only; no dependency on other registries.
10. THE registry SHALL include `CountByServicePlan(ctx context.Context, serviceClassRef string, servicePlanRef string) (int, error)` returning the count of ServiceInstances whose `spec.serviceClassRef` matches `serviceClassRef` AND `spec.servicePlanRef` matches `servicePlanRef`. This uses the fully-qualified plan identity (serviceClass + planName) to avoid false positives when plan names are reused across different ServiceClasses. Used by the FEATURE-0006 ServicePlan delete-blocker (Requirement 14A).
11. THE registry SHALL include `CountByProject(ctx context.Context, organizationRef string, organizationUnitRef string, tenantRef string, projectRef string) (int, error)` returning the count of ServiceInstances whose `spec.organizationRef` matches `organizationRef` AND `spec.organizationUnitRef` matches `organizationUnitRef` AND `spec.tenantRef` matches `tenantRef` AND `spec.projectRef` matches `projectRef`. This uses the fully-qualified project identity to avoid false positives when project names are reused across different tenants. Used by the FEATURE-0004 Project delete-blocker (Requirement 14A).
12. THE registry SHALL produce no data race reports under `go test -race` with 10+ goroutines.

---

### Requirement 6: In-Memory ServiceBinding Registry

**User Story:** As a developer, I want a thread-safe in-memory ServiceBinding registry with filtering support.

#### Acceptance Criteria

1. THE ServiceBindingRegistry SHALL store entries in a `map[string]resources.ServiceBinding` protected by `sync.RWMutex`, keyed by `metadata.name`.
2. THE registry SHALL return deep copies on Create/Get/List.
3. `CreateServiceBinding` SHALL return `(resources.ServiceBinding, error)`.
4. `DeleteServiceBinding(ctx, name)` SHALL return `error` only. This matches prior registry patterns and the HTTP 204 no-body response.
5. THE registry SHALL accept `context.Context` as the first parameter on all public methods; no package-level global state.
6. THE registry SHALL return sentinel errors distinguishing `RESOURCE_NOT_FOUND` from `RESOURCE_ALREADY_EXISTS`.
7. `ListServiceBindings` SHALL return a non-nil slice sorted by `metadata.name` ascending.
8. `ListServiceBindings` SHALL accept optional filter: `serviceInstanceRef`. When provided, results match that filter. When absent, all entries are returned.
9. THE registry SHALL include `CountByServiceInstance(ctx, instanceName) (int, error)` for delete-blocker use on ServiceInstance.
10. THE registry SHALL be storage-only; no dependency on other registries.
11. THE registry SHALL produce no data race reports under `go test -race` with 10+ goroutines.

---

### Requirement 7: ServiceInstance REST API

**User Story:** As a platform operator, I want CRUD endpoints for ServiceInstance.

#### Acceptance Criteria

1. THE Server SHALL register `POST /v1/service-instances`, `GET /v1/service-instances`, `GET/PUT/DELETE /v1/service-instances/{name}`.
2. THE Server SHALL use Go 1.21-compatible routing: `/v1/service-instances` and `/v1/service-instances/` patterns; item path has exactly ONE non-empty segment. Wrong segment count → HTTP 404.
3. IF the `{name}` segment is not a valid DNS-label, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "metadata.name"` without a registry lookup.
4. `POST` valid → 201 with full resource, `status.phase = "Ready"`, `status.message = "Registered only; no real provisioning in Phase 1"`, server-set apiVersion/kind. Duplicate name → 409 `RESOURCE_ALREADY_EXISTS`.
5. `GET` item exists → 200; missing → 404 `RESOURCE_NOT_FOUND`. `GET` list → `{"items": [...]}` sorted by name; empty → `{"items": []}`.
6. `GET` list SHALL accept optional query parameters `tenantRef` and `projectRef` for filtering.
7. `PUT` → 200 on success; missing → 404. `DELETE` → 204 no body; missing → 404.
8. Safe JSON decoding (1 MiB limit, DisallowUnknownFields, status rejection on POST and PUT); HTTP 415 handled by contentTypeMiddleware with `error.code = "VALIDATION_FAILED"` and message "content type must be application/json" (consistent with existing middleware behavior); bad JSON/oversized/unknown-field per existing decoder patterns.

---

### Requirement 8: ServiceBinding REST API

**User Story:** As a platform operator, I want CRUD endpoints for ServiceBinding with filtering.

#### Acceptance Criteria

1. THE Server SHALL register `POST /v1/service-bindings`, `GET /v1/service-bindings`, `GET/DELETE /v1/service-bindings/{name}`.
2. THE Server SHALL use Go 1.21-compatible routing: `/v1/service-bindings` and `/v1/service-bindings/` patterns; item path has exactly ONE non-empty segment. Wrong segment count → HTTP 404.
3. IF the `{name}` segment is not a valid DNS-label, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "metadata.name"` without a registry lookup.
4. `POST` valid → 201 with full resource, `status.phase = "Ready"`, `status.secretRef = "stub-secret-ref"`, server-set apiVersion/kind. Duplicate name → 409 `RESOURCE_ALREADY_EXISTS`.
5. `GET` item exists → 200; missing → 404 `RESOURCE_NOT_FOUND`.
6. `GET` list SHALL accept optional query parameter `serviceInstanceRef` for filtering. Result → `{"items": [...]}` sorted by name; empty → `{"items": []}`.
7. `DELETE` → 204 no body; missing → 404.
8. ServiceBinding does NOT support `PUT` (update). Bindings are immutable after creation; to change, delete and recreate.
9. IF `PUT /v1/service-bindings/{name}` is received, THE Server SHALL return HTTP 405 Method Not Allowed with `error.code = "METHOD_NOT_ALLOWED"` and `error.message` indicating that ServiceBinding does not support update, regardless of whether a ServiceBinding with that name exists.
10. Safe JSON decoding (1 MiB limit, DisallowUnknownFields, status rejection on POST); HTTP 415 handled by contentTypeMiddleware with `error.code = "VALIDATION_FAILED"` and message "content type must be application/json" (consistent with existing middleware behavior); bad JSON/oversized/unknown-field per existing decoder patterns.

---

### Requirement 9: Reference Validation (ServiceInstance → Governance Hierarchy)

**User Story:** As a platform operator, I want ServiceInstance governance references validated against existing resources.

#### Acceptance Criteria

1. WHEN a valid `POST /v1/service-instances` request passes field validation, THE Server SHALL verify `spec.organizationRef` exists as a registered Organization.
2. IF the referenced Organization does not exist, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "spec.organizationRef"`.
3. IF `spec.organizationUnitRef` is present and non-empty, THE Server SHALL verify it exists as a registered OrganizationUnit AND that its `spec.organizationRef` matches the ServiceInstance's `spec.organizationRef`. Mismatch → HTTP 400 `VALIDATION_FAILED`, `error.field = "spec.organizationUnitRef"`.
4. THE Server SHALL verify `spec.tenantRef` exists as a registered Tenant AND that the Tenant's organization lineage is consistent with the ServiceInstance's `spec.organizationRef` (and `spec.organizationUnitRef` if present). Inconsistency → HTTP 400 `VALIDATION_FAILED`, `error.field = "spec.tenantRef"`.
5. THE Server SHALL verify `spec.projectRef` exists as a registered Project AND that the Project's `spec.tenantRef` matches the ServiceInstance's `spec.tenantRef`. Mismatch → HTTP 400 `VALIDATION_FAILED`, `error.field = "spec.projectRef"`.
6. THE reference checks SHALL occur in the API/service layer via narrow lookup interfaces; the ServiceInstanceRegistry SHALL NOT depend on other registries directly.
7. WHEN a `PUT` update passes field validation, identity checks, and immutability checks, THE Server SHALL re-verify mutable governance references (`spec.organizationRef`, `spec.organizationUnitRef`) still exist and are consistent with the immutable `spec.tenantRef`; if any are absent or inconsistent, THE Server SHALL return HTTP 400 `VALIDATION_FAILED` with the appropriate field. Immutable references (`spec.tenantRef`, `spec.projectRef`) are not re-validated since they cannot change.

---

### Requirement 10: Reference Validation (ServiceInstance → Catalog)

**User Story:** As a platform operator, I want ServiceInstance catalog references validated against existing ServiceClass and ServicePlan.

#### Acceptance Criteria

1. WHEN a valid `POST /v1/service-instances` request passes field validation, THE Server SHALL verify `spec.serviceClassRef` exists as a registered ServiceClass.
2. IF the referenced ServiceClass does not exist, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "spec.serviceClassRef"`.
3. THE Server SHALL verify `spec.servicePlanRef` exists as a registered ServicePlan AND that the ServicePlan's `spec.serviceClassRef` matches the ServiceInstance's `spec.serviceClassRef`. Mismatch → HTTP 400 `VALIDATION_FAILED`, `error.field = "spec.servicePlanRef"`.
4. IF the referenced ServicePlan does not exist, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "spec.servicePlanRef"`.
5. ON `PUT` update, catalog references (`spec.serviceClassRef`, `spec.servicePlanRef`) are immutable (per Requirement 14.4). The server rejects changes to these fields before catalog re-validation would occur. No catalog re-verification is needed on update.
6. THE reference checks SHALL use narrow lookup interfaces consistent with existing patterns.

---

### Requirement 11: Capability Warning (ServiceInstance → Plugin Registry)

**User Story:** As a platform operator, I want ServiceInstance creation to warn if no active Capability exists for the ServiceClass, but not block creation in Phase 1.

#### Acceptance Criteria

1. WHEN a valid `POST /v1/service-instances` request passes all reference validation, THE Server SHALL check whether at least one active Capability exists for the ServiceInstance's `spec.serviceClassRef` by calling `HasActiveCapabilityForServiceClass(ctx, serviceClassRef) (bool, error)` on the injected `CapabilityLookup` interface.
2. IF no active Capability exists (returns `false`), THE Server SHALL log a structured warning with the ServiceInstance name, ServiceClass name, and a message indicating no active capability is registered.
3. THE Server SHALL NOT block ServiceInstance creation based on missing Capability in Phase 1.
4. THE Server SHALL NOT include the warning in the API response body. The warning is observability-only (server log).
5. This check is informational only. Future phases may enforce capability requirements.
6. The definition of "active capability" is encapsulated within the FEATURE-0007 `CapabilityLookup` implementation; FEATURE-0008 does not reference internal Capability fields directly. Implementation details of how capability status is determined (e.g., which fields are checked) remain behind the interface boundary.

---

### Requirement 12: Reference Validation (ServiceBinding → ServiceInstance)

**User Story:** As a platform operator, I want ServiceBinding references validated against existing ServiceInstances.

#### Acceptance Criteria

1. WHEN a valid `POST /v1/service-bindings` request passes field validation, THE Server SHALL verify `spec.serviceInstanceRef` exists as a registered ServiceInstance.
2. IF the referenced ServiceInstance does not exist, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "spec.serviceInstanceRef"`.
3. THE reference check SHALL use a narrow lookup interface; the ServiceBindingRegistry SHALL NOT depend on the ServiceInstanceRegistry directly.

---

### Requirement 13: ServiceInstance Delete Blocking

**User Story:** As a platform operator, I want ServiceInstance deletion blocked while ServiceBindings still reference it.

#### Acceptance Criteria

1. WHEN `DELETE /v1/service-instances/{name}` is requested and one or more ServiceBindings have `spec.serviceInstanceRef == name`, THE Server SHALL return HTTP 409 with `error.code = "DELETE_BLOCKED"` and a message identifying `"ServiceBinding"` as the blocking kind.
2. THE blocker SHALL query the ServiceBindingRegistry via `CountByServiceInstance` and block when count > 0.
3. WHEN the ServiceInstance has zero ServiceBindings, THE delete SHALL proceed and return 204.
4. THE blocker SHALL use a narrow blocker interface consistent with existing patterns; no generic blocker framework.

---

### Requirement 14: Update Identity and Mutability (ServiceInstance)

**User Story:** As a platform operator, I want PUT requests to enforce path/body identity matching and clear mutability rules for ServiceInstance.

#### Acceptance Criteria

1. FOR `PUT /v1/service-instances/{name}`, THE Server SHALL require body `metadata.name` present and equal to the path `{name}`. Absent, empty, or mismatched → HTTP 400 `VALIDATION_FAILED`, `error.field = "metadata.name"`.
2. THE registry update SHALL preserve immutable identity fields and SHALL NOT rename a ServiceInstance.
3. THE Server SHALL preserve server-owned fields on update: `apiVersion`, `kind`, and `status` (both `phase` and `message`) are copied from the stored entry unchanged. Client-supplied values for these fields in the request body SHALL be ignored (status key is rejected per Requirement 1.6; apiVersion/kind are overwritten by server).
4. **Immutable spec fields** (cannot be changed after creation): `spec.tenantRef`, `spec.projectRef`, `spec.serviceClassRef`, `spec.servicePlanRef`. IF a PUT request supplies a different value for any immutable spec field compared to the stored entry, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field` set to the changed field (e.g., `"spec.tenantRef"`).
5. **Mutable spec fields** (can be changed on update): `spec.organizationRef`, `spec.organizationUnitRef`, `spec.parameters`.
6. **Mutable metadata fields** (can be changed on update): `metadata.labels`, `metadata.annotations`.
7. IF mutable spec fields are changed, all applicable reference validations (governance hierarchy) SHALL be re-executed on the new values. Specifically, if `spec.organizationRef` or `spec.organizationUnitRef` changes, the new values must be validated for existence and consistency with the (immutable) `spec.tenantRef`.
8. THE status-key rejection (Requirement 1.6) applies to PUT requests: if the top-level request body contains the key `status`, THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "status"` before any other processing.

---

### Requirement 14A: Cross-Feature Delete-Blockers (ServicePlan and Project)

**User Story:** As a platform operator, I want ServicePlan and Project deletions blocked while ServiceInstances reference them, ensuring referential integrity.

#### Acceptance Criteria

1. THE ServiceInstanceRegistry SHALL expose `CountByServicePlan(ctx context.Context, serviceClassRef string, servicePlanRef string) (int, error)` returning the count of ServiceInstances whose `spec.serviceClassRef` matches `serviceClassRef` AND `spec.servicePlanRef` matches `servicePlanRef`.
2. THE ServiceInstanceRegistry SHALL expose `CountByProject(ctx context.Context, organizationRef string, organizationUnitRef string, tenantRef string, projectRef string) (int, error)` returning the count of ServiceInstances whose governance refs match all four parameters.
3. WHEN `DELETE /v1/service-plans/{serviceClassName}/{name}` is requested (FEATURE-0006 handler) and `CountByServicePlan(ctx, serviceClassName, name)` returns count > 0, THE Server SHALL return HTTP 409 with `error.code = "DELETE_BLOCKED"` and a message identifying `"ServiceInstance"` as the blocking kind.
4. WHEN `DELETE /v1/projects/{org}/{ou}/{tenant}/{name}` is requested (FEATURE-0004 handler) and `CountByProject(ctx, org, ou, tenant, name)` returns count > 0, THE Server SHALL return HTTP 409 with `error.code = "DELETE_BLOCKED"` and a message identifying `"ServiceInstance"` as the blocking kind.
5. FEATURE-0008 implementation SHALL wire these blocker checks into the existing ServicePlan and Project delete handlers via narrow blocker interfaces, consistent with existing child-blocker patterns (e.g., `ServicePlanChildBlockerChecker`).
6. WHEN the ServicePlan has zero referencing ServiceInstances and the Project has zero ServiceInstances, THE respective deletes SHALL proceed normally.

---

### Requirement 15: Operation Emission (FEATURE-0005 Integration)

**User Story:** As a platform operator, I want service consumption lifecycle actions recorded as Operations.

#### Acceptance Criteria

1. AFTER a successful ServiceInstance create/update/delete, THE Server SHALL emit an Operation of type `CreateServiceInstance`, `UpdateServiceInstance`, or `DeleteServiceInstance` respectively, with `resourceKind = "ServiceInstance"`.
2. AFTER a successful ServiceBinding create/delete, THE Server SHALL emit an Operation of type `CreateServiceBinding` or `DeleteServiceBinding` respectively, with `resourceKind = "ServiceBinding"`.
3. THE resource kind constants (`ServiceInstanceKind`, `ServiceBindingKind`) SHALL be added to `internal/resources` following the existing pattern.
4. FEATURE-0008 SHALL extend the FEATURE-0005 operation type constants with: `CreateServiceInstance`, `UpdateServiceInstance`, `DeleteServiceInstance`, `CreateServiceBinding`, `DeleteServiceBinding`.
5. FEATURE-0008 SHALL extend `OperationSpec` with two optional fields: `serviceInstanceName` (`json:"serviceInstanceName,omitempty"`) and `serviceBindingName` (`json:"serviceBindingName,omitempty"`). These names are final and consistent with the established naming convention (`serviceClassName`, `servicePlanName`, `pluginName`, `capabilityName`). These fields are used only for ServiceInstance and ServiceBinding Operation records and remain empty/omitted for other resource kinds.
6. ServiceInstance Operation records SHALL set `serviceInstanceName` to the instance name. ServiceBinding Operation records SHALL set `serviceInstanceName` to the referenced instance name and `serviceBindingName` to the binding name.
7. Emission SHALL use the nil-safe FEATURE-0005 emitter; emission failure SHALL NOT affect the primary API response.
8. THE Server SHALL NOT emit an Operation on failed validation, duplicate create, missing reference, not-found, or delete-blocked cases.

---

### Requirement 16: Security and Privacy

**User Story:** As a security-conscious operator, I want service consumption resources to handle credentials safely.

#### Acceptance Criteria

1. ServiceInstance `spec.parameters` SHALL NOT be used to store secrets, tokens, passwords, or credentials. Documentation SHALL warn against it. No enforcement in Phase 1.
2. ServiceBinding `status.secretRef` is a stub reference (`"stub-secret-ref"`) in Phase 1. It does NOT point to a real secret backend. No actual credential generation or storage occurs.
3. THE Server SHALL NOT store raw request bodies and SHALL NOT echo raw bodies in error responses.
4. THE Server SHALL NOT log secrets, parameters values, or raw bodies.
5. Future phases will integrate with a secret backend for real credential management. Phase 1 only establishes the secretRef field structure.
6. THE Server SHALL NOT expose internal implementation details in error messages.

---

### Requirement 17: Tests

**User Story:** As a developer, I want comprehensive tests for ServiceInstance and ServiceBinding resources.

#### Acceptance Criteria

1. Validation unit tests: valid names accepted; invalid/empty/long names rejected; missing required fields rejected; invalid organizationUnitRef rejected; valid bindingType accepted; invalid bindingType rejected; missing consumerRef fields rejected.
2. Registry unit tests (both registries): Create stores; duplicate → ErrAlreadyExists (original unchanged); Get by key; missing → ErrNotFound; List sorted; empty → non-nil `[]`; ServiceInstance Update mutable fields only; Update preserves stored status unchanged; Update rejects immutable field change (returns error); Update missing → ErrNotFound; Delete removes; Delete missing → ErrNotFound; ServiceBinding CountByServiceInstance correct; ServiceInstance List with filters (tenantRef, projectRef, both, neither); ServiceBinding List with filter (serviceInstanceRef); ServiceInstance CountByServicePlan correct (matching both serviceClassRef and servicePlanRef, no false positives across different ServiceClasses); ServiceInstance CountByProject correct (matching all four hierarchy refs, no false positives across different tenants).
3. Handler tests (both resources): POST 201/409/400 (invalid fields, status key, bad JSON, unknown field); POST 413 oversized; ServiceInstance POST 400 missing Organization ref; ServiceInstance POST 400 missing Tenant ref; ServiceInstance POST 400 missing Project ref; ServiceInstance POST 400 missing ServiceClass ref; ServiceInstance POST 400 missing ServicePlan ref; ServiceInstance POST 400 ServicePlan not matching ServiceClass; ServiceInstance POST 400 governance hierarchy inconsistency; ServiceBinding POST 400 missing ServiceInstance ref; GET 200/404/400; wrong path shape → 404; list sorted/empty; list with query filters; ServiceInstance PUT 200/404/400 (identity mismatch); ServiceInstance PUT 400 with status key present → VALIDATION_FAILED field=status; ServiceInstance PUT 400 immutable field change (tenantRef/projectRef/serviceClassRef/servicePlanRef) → VALIDATION_FAILED with field; ServiceInstance PUT 200 preserves stored status unchanged; DELETE 204/404; ServiceBinding PUT → 405 METHOD_NOT_ALLOWED.
4. Delete-blocking tests: ServiceInstance delete with ServiceBindings → 409 DELETE_BLOCKED; with zero ServiceBindings → 204; ServicePlan delete with referencing ServiceInstances → 409 DELETE_BLOCKED; ServicePlan delete with zero referencing ServiceInstances → normal delete; Project delete with ServiceInstances → 409 DELETE_BLOCKED; Project delete with zero ServiceInstances → normal delete.
5. Operation emission tests: successful create/update/delete of ServiceInstance and create/delete of ServiceBinding records the correct Operation type and resource kind; failed actions emit nothing; emission failure does not change the primary response.
6. Capability warning test: ServiceInstance created without matching Capability → warning logged (no block).
7. `go test -race ./...` with 10+ goroutines produces no race reports.
8. ALL tests deterministic; no external dependencies.

---

### Requirement 18: Property-Based Tests

**User Story:** As a developer, I want property tests for validation and registries.

#### Acceptance Criteria

1. Validation package property tests (`testing/quick`, `Config{MaxCount: 100}`): valid DNS-label names accepted; arbitrary invalid strings rejected; valid bindingType accepted; invalid bindingType rejected.
2. Registry property tests: Create/Get round-trip preserves data; List sort invariant; deep-copy immutability; duplicate-create idempotent error (original unchanged); ServiceInstance filter correctness; ServiceBinding filter correctness; CountByServiceInstance correctness; CountByServicePlan correctness (no false positives when same planName exists under different ServiceClasses); CountByProject correctness (no false positives when same projectName exists under different tenants).
3. Each property test tagged `// Feature: serviceinstance-servicebinding, Property N: <title>`.

---

### Requirement 19: Non-Goals

**User Story:** As an architect, I want clear scope boundaries.

#### Acceptance Criteria

1. NO real infrastructure provisioning, container creation, or database deployment.
2. NO actual credential generation, secret storage, or secret backend integration.
3. NO plugin execution, ServiceOps orchestration, or workflow engine.
4. NO async provisioning workflows, queues, or background workers.
5. NO Kubernetes CRDs, operators, or GitOps controllers.
6. NO authentication, authorization, RBAC, or policy enforcement on service consumption.
7. NO quota enforcement or resource limit checking.
8. NO billing, metering, or cost tracking.
9. NO UI portal or dashboard.
10. NO persistence/database storage.
11. NO Go 1.22 wildcard routing; no new external dependencies.
12. NO SDE-specific behavior or transformation logic.
13. NO multi-cluster federation or cross-cluster service binding.
14. NO service mesh integration or network policy creation.

---

### Requirement 20: Edge Cases

**User Story:** As a developer, I want edge-case behavior defined.

#### Acceptance Criteria

1. WHEN a mutating action fails (validation, duplicate, missing reference, not-found, delete-blocked), THE Server SHALL NOT emit an Operation.
2. WHEN a ServiceInstance item path has more than one segment, THE Server SHALL return HTTP 404.
3. WHEN a ServiceBinding item path has more than one segment, THE Server SHALL return HTTP 404.
4. WHEN `spec.parameters` contains an empty map `{}`, THE Server SHALL accept the request (empty parameters is valid).
5. WHEN `spec.parameters` is absent/nil, THE Server SHALL accept the request (parameters are optional).
6. WHEN `spec.organizationUnitRef` is absent or empty string, THE Server SHALL skip OrganizationUnit validation (field is optional).
7. Governance resource deletions (Organization, Tenant, Project) are blocked while dependent ServiceInstances exist (per Requirement 14A and existing child-blocker patterns in FEATURE-0001 through FEATURE-0004). Referential integrity is checked at write time only, not continuously enforced; however, parent delete APIs actively prevent deletion while children reference them. A ServiceInstance's governance references cannot become dangling while the platform enforces delete-blockers.
8. Catalog resource deletions: ServicePlan delete is blocked while ServiceInstances reference it (per Requirement 14A). ServiceClass delete is blocked by ServicePlans (existing FEATURE-0006 behavior). A ServiceInstance's catalog references cannot become dangling while the platform enforces delete-blockers. Referential integrity is checked at write time only.
9. WHEN a nil emitter is used (isolated handler tests), emission SHALL be skipped gracefully without panic.
10. WHEN `PUT /v1/service-bindings/{name}` is received, THE Server SHALL return HTTP 405 Method Not Allowed with `error.code = "METHOD_NOT_ALLOWED"` (ServiceBinding does not support update).
11. WHEN the same `metadata.name` is used for both a ServiceInstance and a ServiceBinding, THE Server SHALL allow it (ServiceInstance and ServiceBinding are distinct resource kinds with separate registries).
12. WHEN `spec.consumerRef.kind` is any non-empty string (e.g., `"Application"`, `"Job"`, `"CronJob"`), THE Server SHALL accept it. No enum restriction in Phase 1.
13. WHEN a ServiceInstance is created with a ServicePlan whose `spec.serviceClassRef` does not match the ServiceInstance's `spec.serviceClassRef`, THE Server SHALL reject with HTTP 400 `VALIDATION_FAILED`, `error.field = "spec.servicePlanRef"`.

---

## Compatibility with Completed Phase 1 Features

### FEATURE-0001 (Organization Resource and Registry)

ServiceInstance `spec.organizationRef` references an Organization. The handler validates existence
via a narrow `OrganizationLookup` interface (consistent with existing patterns from FEATURE-0002
through FEATURE-0004). Organization cannot be deleted while dependent resources exist (existing
child-blocking from OrganizationUnit).

### FEATURE-0002 (OrganizationUnit Resource)

ServiceInstance `spec.organizationUnitRef` (optional) references an OrganizationUnit. The handler
validates existence and verifies that the OrganizationUnit belongs to the same Organization as
specified in `spec.organizationRef`. Uses a narrow `OrganizationUnitLookup` interface.

### FEATURE-0003 (Tenant Resource)

ServiceInstance `spec.tenantRef` references a Tenant. The handler validates existence and verifies
governance hierarchy consistency (Tenant must belong to the correct Organization/OrganizationUnit
lineage). Uses a narrow `TenantLookup` interface.

### FEATURE-0004 (Project Resource)

ServiceInstance `spec.projectRef` references a Project. The handler validates existence and
verifies that the Project belongs to the specified Tenant. Uses a narrow `ProjectLookup`
interface. Project cannot be deleted while ServiceInstances exist under it (this delete-blocking
is a new behavior added by FEATURE-0008).

### FEATURE-0005 (Operation Resource)

FEATURE-0008 reuses the FEATURE-0005 `OperationEmitter` interface and nil-safe `emitOperation`
helper. ServiceInstance and ServiceBinding handlers receive the emitter via constructor injection
(nil-safe), consistent with all other resource handlers. Five new operation type constants
(`CreateServiceInstance`, `UpdateServiceInstance`, `DeleteServiceInstance`,
`CreateServiceBinding`, `DeleteServiceBinding`) and two resource kind constants
(`ServiceInstanceKind`, `ServiceBindingKind`) are added to `internal/resources`. Emission occurs
only after a successful mutating action, uses the request context for the request ID, and never
alters the primary response.

### FEATURE-0006 (ServiceClass and ServicePlan)

ServiceInstance `spec.serviceClassRef` references a ServiceClass registered via FEATURE-0006.
ServiceInstance `spec.servicePlanRef` references a ServicePlan and the handler verifies that the
ServicePlan's `spec.serviceClassRef` matches the ServiceInstance's `spec.serviceClassRef`.
Reference validation uses narrow `ServiceClassLookup` and `ServicePlanLookup` interfaces injected
into the handler, not direct registry dependency.

Additionally, ServicePlan delete-blocking (FEATURE-0006) must be extended: a ServicePlan cannot
be deleted if any ServiceInstance references it via `spec.servicePlanRef`. This adds a new
delete-blocker to the ServicePlan delete handler.

### FEATURE-0007 (Plugin and Capability Registry)

FEATURE-0008 uses a narrow `CapabilityLookup` interface to check whether at least one active
Capability exists for a ServiceClass. This is used for the informational warning (Requirement 11)
and does not block ServiceInstance creation. The interface provides the method
`HasActiveCapabilityForServiceClass(ctx, serviceClassRef) (bool, error)`. The definition of
"active" is encapsulated within the FEATURE-0007 implementation; FEATURE-0008 does not reference
internal Capability struct fields (such as status booleans) directly.

## Design Questions

> These questions are OPEN and MUST be resolved in `design.md` before implementation.

1. **Governance hierarchy consistency validation:** Define the exact logic for validating that Tenant belongs to the correct Organization (and OrganizationUnit if specified), and that Project belongs to the specified Tenant. Confirm whether the lookup interfaces need to return the full resource or just existence + parent refs.

2. **Lookup interfaces placement:** Confirm whether lookup interfaces for Organization, OrganizationUnit, Tenant, Project, ServiceClass, ServicePlan, ServiceInstance, and Capability are placed in `internal/api` or `internal/registry`. Confirm naming consistency with existing interfaces (e.g., `OrganizationLookup`, `TenantLookup`).

3. **ServicePlan delete-blocker extension:** FEATURE-0006 ServicePlan delete must now also check for referencing ServiceInstances. The interface shape is decided: `CountByServicePlan(ctx, serviceClassRef, servicePlanRef) (int, error)` on the ServiceInstanceRegistry (see Requirement 14A). The caller passes both the `serviceClassName` and `name` from the delete path. Design.md must define whether this is wired as a new narrow blocker interface injected into the ServicePlan handler or via an adapter pattern consistent with existing `ServicePlanChildBlockerChecker`.

4. **Project delete-blocker extension:** FEATURE-0004 Project delete must now also check for ServiceInstances under it. The interface shape is decided: `CountByProject(ctx, organizationRef, organizationUnitRef, tenantRef, projectRef) (int, error)` on the ServiceInstanceRegistry (see Requirement 14A). The caller passes all four path segments from the delete path. Design.md must define the wiring approach consistent with the existing Project delete handler.

5. **ServiceBinding immutability rationale:** The API contract shows no `PUT` for ServiceBinding. Confirm that delete-and-recreate is the correct update pattern.

6. **Capability warning interface:** The normative interface is `HasActiveCapabilityForServiceClass(ctx, serviceClassRef) (bool, error)`. Design.md must define the concrete implementation (which Capability fields determine "active") within the FEATURE-0007 registry boundary. FEATURE-0008 only depends on the interface, not on internal field semantics.

7. **Stub secretRef value:** Confirm that `"stub-secret-ref"` is the correct placeholder value for Phase 1 ServiceBinding status.secretRef, or whether a more descriptive value is preferred.

8. **Filter query parameters:** Confirm query parameter names for list filtering: `tenantRef` and `projectRef` for ServiceInstance; `serviceInstanceRef` for ServiceBinding. Confirm whether additional filters (e.g., `serviceClassRef`, `servicePlanRef`) are needed for Phase 1 or deferred.

### Resolved Questions (decided in requirements, no longer open)

- **ServiceInstance update scope (formerly Q5):** Resolved in Requirement 14. `spec.tenantRef`, `spec.projectRef`, `spec.serviceClassRef`, `spec.servicePlanRef` are immutable after creation. `spec.organizationRef`, `spec.organizationUnitRef`, `spec.parameters` are mutable.
- **Operation spec fields naming (formerly Q8):** Resolved in Requirement 15. `serviceInstanceName` and `serviceBindingName` are the final canonical field names, consistent with `serviceClassName`, `servicePlanName`, `pluginName`, `capabilityName`.
