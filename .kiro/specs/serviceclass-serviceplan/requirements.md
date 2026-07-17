# Requirements Document

## Introduction

FEATURE-0006 introduces `ServiceClass` and `ServicePlan` as the Phase 1 service catalog
foundation in Sovrunn. These are control-plane catalog definitions: a `ServiceClass` defines a
type of managed service (e.g. PostgreSQL, Redis, ObjectStorage, SDEGateway), and a `ServicePlan`
defines an offering/shape/tier under a ServiceClass (e.g. PostgreSQL dev-small, Redis cache-small).

These resources define what services **could** be requested later. They do NOT provision, bind, or
execute anything. ServiceInstance, ServiceBinding, Plugin/Capability registry, and ServiceOps
execution are future features and are out of scope.

**Scope of resources:** ServiceClass and ServicePlan are **global platform catalog resources** in
FEATURE-0006. They are NOT scoped to, owned by, or nested under Organization, OrganizationUnit,
Tenant, or Project. Their identities do not include any governance-hierarchy path segments:
a ServiceClass is identified by `metadata.name` alone, and a ServicePlan by the composite
`serviceClassName/name`.

This feature depends on FEATURE-0001 through FEATURE-0005. It reuses the existing project skeleton,
in-memory registry patterns, safe JSON decoding, structured errors, and the FEATURE-0005 Operation
emission mechanism.

## Glossary

- **ServiceClass**: A catalog definition of a type of managed service. Identity: `metadata.name`.
- **ServicePlan**: A catalog offering/tier under a ServiceClass. Identity: composite `serviceClassName/name`.
- **Category**: ServiceClass classification (Database, Cache, ObjectStorage, Stream, Gateway, Function, Analytics, Other).
- **Lifecycle**: Catalog lifecycle state (Preview, Active, Deprecated, Retired).
- **Tier**: ServicePlan sizing/shape (Dev, Small, Medium, Large, Production, Custom).
- **Registry**: In-memory, thread-safe store backed by `sync.RWMutex`, consistent with existing registries.
- **Error Code**: Stable string from: `VALIDATION_FAILED`, `RESOURCE_NOT_FOUND`, `RESOURCE_ALREADY_EXISTS`, `DELETE_BLOCKED`, `INTERNAL_ERROR`.
- **OperationEmitter**: The FEATURE-0005 mechanism for recording lifecycle actions.
- **DNS-label name**: Matches `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`, 1–63 chars.

## Requirements

---

### Requirement 1: ServiceClass Resource Shape

**User Story:** As a catalog operator, I want ServiceClass to follow the canonical `metadata/spec/status` shape, so that it is consistent with all Phase 1 resources.

#### Acceptance Criteria

1. THE `ServiceClass` struct SHALL include `APIVersion`, `Kind`, `Metadata`, `Spec`, and `Status` with JSON tags `apiVersion`, `kind`, `metadata`, `spec`, `status`.
2. THE Server SHALL set `apiVersion` to `"platform.sovrunn.io/v1alpha1"` and `kind` to `"ServiceClass"` on all successful responses regardless of client input.
3. THE `ServiceClassSpec` SHALL include `DisplayName` (`json:"displayName,omitempty"`), `Description` (`json:"description,omitempty"`), `Category` (`json:"category"`), `Provider` (`json:"provider,omitempty"`), `Lifecycle` (`json:"lifecycle"`), `DefaultPlanName` (`json:"defaultPlanName,omitempty"`), and `Tags` (`json:"tags,omitempty"`, `[]string`).
4. THE `ServiceClassStatus` SHALL include `Phase` (`json:"phase"`) and `Message` (`json:"message,omitempty"`).
5. THE `metadata.name` SHALL be the ServiceClass name (DNS-label).
6. IF the top-level request body contains the key `status` — regardless of value — THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "status"`.

---

### Requirement 2: ServicePlan Resource Shape

**User Story:** As a catalog operator, I want ServicePlan to follow the canonical shape and reference its parent ServiceClass.

#### Acceptance Criteria

1. THE `ServicePlan` struct SHALL include `APIVersion`, `Kind`, `Metadata`, `Spec`, and `Status` with the standard JSON tags.
2. THE Server SHALL set `apiVersion` to `"platform.sovrunn.io/v1alpha1"` and `kind` to `"ServicePlan"` on all successful responses.
3. THE `ServicePlanSpec` SHALL include `ServiceClassName` (`json:"serviceClassName"`), `DisplayName` (`json:"displayName,omitempty"`), `Description` (`json:"description,omitempty"`), `Tier` (`json:"tier"`), `Lifecycle` (`json:"lifecycle"`), `Parameters` (`json:"parameters,omitempty"`, `map[string]string`), and `Tags` (`json:"tags,omitempty"`, `[]string`).
4. THE `ServicePlanStatus` SHALL include `Phase` (`json:"phase"`) and `Message` (`json:"message,omitempty"`).
5. THE `metadata.name` SHALL be the ServicePlan name (DNS-label); identity is composite `serviceClassName/name`.
6. IF the top-level request body contains the key `status` — regardless of value — THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "status"`.

---

### Requirement 3: ServiceClass Validation

**User Story:** As a catalog operator, I want deterministic ServiceClass validation, so that only well-formed catalog entries are stored.

#### Acceptance Criteria

1. IF `metadata.name` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"`.
2. THE `spec.category` SHALL be one of: `Database`, `Cache`, `ObjectStorage`, `Stream`, `Gateway`, `Function`, `Analytics`, `Other`. Any other value → `FieldError` with `Field = "spec.category"`.
3. THE `spec.lifecycle` SHALL be one of: `Preview`, `Active`, `Deprecated`, `Retired`. Any other value → `FieldError` with `Field = "spec.lifecycle"`.
4. `spec.category` and `spec.lifecycle` SHALL be required (empty → the respective FieldError).
5. THE Validator SHALL be a pure function (context-free, no I/O, no registry lookup).
6. IF `spec.defaultPlanName` is non-empty, THEN it SHALL be a valid DNS-label (1–63 chars); otherwise → `FieldError` with `Field = "spec.defaultPlanName"`. FEATURE-0006 does NOT verify that the referenced default plan actually exists.
7. `spec.displayName`, `spec.description`, `spec.provider`, and `spec.tags` are optional and not format-validated in Phase 1 beyond basic type correctness.

---

### Requirement 4: ServicePlan Validation

**User Story:** As a catalog operator, I want deterministic ServicePlan validation.

#### Acceptance Criteria

1. IF `metadata.name` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "metadata.name"`.
2. IF `spec.serviceClassName` is absent, empty, not DNS-label, or exceeds 63 chars, THEN THE Validator SHALL return a `FieldError` with `Field = "spec.serviceClassName"`.
3. THE `spec.tier` SHALL be one of: `Dev`, `Small`, `Medium`, `Large`, `Production`, `Custom`. Any other value → `FieldError` with `Field = "spec.tier"`.
4. THE `spec.lifecycle` SHALL be one of: `Preview`, `Active`, `Deprecated`, `Retired`. Any other value → `FieldError` with `Field = "spec.lifecycle"`.
5. `spec.serviceClassName`, `spec.tier`, and `spec.lifecycle` SHALL be required.
6. THE Validator SHALL reject any `spec.parameters` key that, case-insensitively, contains any of these secret-bearing substrings: `password`, `secret`, `token`, `credential`, `auth`, `apiKey`, `accessKey`, `secretKey`, `privateKey`. Such a key → `FieldError` with `Field = "spec.parameters"`. The plain substring `key` by itself SHALL NOT trigger rejection unless it appears as part of one of the secret-bearing phrases above (e.g. `apiKey`, `accessKey`, `secretKey`, `privateKey`), so benign keys like `masterKeyCount` are not falsely rejected while `regionKey` is allowed and `apiKey` is rejected.
7. THE Validator SHALL be a pure function (context-free, no I/O, no registry lookup).

---

### Requirement 5: In-Memory ServiceClass Registry

**User Story:** As a developer, I want a thread-safe in-memory ServiceClass registry consistent with existing registries.

#### Acceptance Criteria

1. THE ServiceClassRegistry SHALL store entries in a `map[string]resources.ServiceClass` protected by `sync.RWMutex`, keyed by `metadata.name`.
2. THE registry SHALL return deep copies on Create/Get/List/Update (including Tags slice and any maps).
3. `CreateServiceClass` and `UpdateServiceClass` SHALL return `(resources.ServiceClass, error)`.
4. `UpdateServiceClass` SHALL preserve `metadata.name`, `status`, `apiVersion`, `kind` from the stored entry; replace only mutable fields.
5. THE registry SHALL accept `context.Context` as the first parameter on all public methods; no package-level global state.
6. THE registry SHALL return sentinel errors distinguishing `RESOURCE_NOT_FOUND` from `RESOURCE_ALREADY_EXISTS`.
7. `ListServiceClasses` SHALL return a non-nil slice sorted by `metadata.name` ascending.
8. THE registry SHALL be storage-only; no dependency on other registries.
9. THE registry SHALL produce no data race reports under `go test -race` with 10+ goroutines.

---

### Requirement 6: In-Memory ServicePlan Registry

**User Story:** As a developer, I want a thread-safe in-memory ServicePlan registry with composite identity.

#### Acceptance Criteria

1. THE ServicePlanRegistry SHALL store entries in a `map[string]resources.ServicePlan` protected by `sync.RWMutex`, keyed by composite `serviceClassName/name`.
2. WHEN two ServicePlans have the same `metadata.name` but different `spec.serviceClassName`, THE registry SHALL store both without conflict.
3. THE registry SHALL return deep copies on Create/Get/List/Update (including Parameters map and Tags slice).
4. `CreateServicePlan` and `UpdateServicePlan` SHALL return `(resources.ServicePlan, error)`.
5. `UpdateServicePlan` SHALL preserve `metadata.name`, `spec.serviceClassName`, `status`, `apiVersion`, `kind` from the stored entry; replace only mutable fields.
6. THE registry SHALL accept `context.Context` as the first parameter on all public methods; no package-level global state.
7. THE registry SHALL return sentinel errors distinguishing `RESOURCE_NOT_FOUND` from `RESOURCE_ALREADY_EXISTS`.
8. `ListServicePlans` SHALL return a non-nil slice sorted by `spec.serviceClassName` ascending, then `metadata.name` ascending.
9. THE registry SHALL include `CountByServiceClass(ctx, serviceClassName) (int, error)` for delete-blocker use.
10. THE registry SHALL be storage-only; no dependency on other registries.
11. THE registry SHALL produce no data race reports under `go test -race` with 10+ goroutines.

---

### Requirement 7: ServiceClass REST API

**User Story:** As a catalog operator, I want CRUD endpoints for ServiceClass.

#### Acceptance Criteria

1. THE Server SHALL register `POST /v1/service-classes`, `GET /v1/service-classes`, `GET/PUT/DELETE /v1/service-classes/{name}`.
2. THE Server SHALL use Go 1.21-compatible routing: `/v1/service-classes` and `/v1/service-classes/` patterns; item path has exactly ONE non-empty segment. Wrong segment count → HTTP 404.
3. IF the `{name}` segment is not a valid DNS-label, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "metadata.name"` without a registry lookup.
4. `POST` valid → 201 with full resource, `status.phase = "Active"`, server-set apiVersion/kind. Duplicate name → 409 `RESOURCE_ALREADY_EXISTS`.
5. `GET` item exists → 200; missing → 404 `RESOURCE_NOT_FOUND`. `GET` list → `{"items": [...]}` sorted by name; empty → `{"items": []}`.
6. `PUT` → 200 on success; missing → 404. `DELETE` → 204 no body; missing → 404.
7. Safe JSON decoding (1 MiB limit, DisallowUnknownFields, status rejection); 415 handled by contentTypeMiddleware; bad JSON/oversized/unknown-field per existing decoder patterns.

---

### Requirement 8: ServicePlan REST API

**User Story:** As a catalog operator, I want CRUD endpoints for ServicePlan under its ServiceClass.

#### Acceptance Criteria

1. THE Server SHALL register `POST /v1/service-plans`, `GET /v1/service-plans`, `GET/PUT/DELETE /v1/service-plans/{serviceClassName}/{name}`.
2. THE Server SHALL use Go 1.21-compatible routing: `/v1/service-plans` and `/v1/service-plans/` patterns; item path has exactly TWO non-empty segments (`serviceClassName`, `name`). Wrong segment count → HTTP 404.
3. Path validation field mapping: invalid `serviceClassName` segment → `error.field = "spec.serviceClassName"`; invalid `name` segment → `error.field = "metadata.name"`. Invalid segment → 400 `VALIDATION_FAILED` without a registry lookup.
4. `POST` valid → 201 with full resource, `status.phase = "Active"`. Duplicate composite key → 409 `RESOURCE_ALREADY_EXISTS`.
5. `GET` item exists → 200; missing → 404. `GET` list → `{"items": [...]}` sorted by serviceClassName then name; empty → `{"items": []}`.
6. `PUT` → 200 on success; missing → 404. `DELETE` → 204 no body; missing → 404.
7. Safe JSON decoding identical to ServiceClass (1 MiB, DisallowUnknownFields, status rejection); 415 via contentTypeMiddleware.

---

### Requirement 9: Parent Relationship (ServicePlan → ServiceClass)

**User Story:** As a catalog operator, I want ServicePlans to require an existing parent ServiceClass.

#### Acceptance Criteria

1. WHEN a valid `POST /v1/service-plans` request passes field validation, THE Server SHALL verify the parent ServiceClass named `spec.serviceClassName` exists before storing.
2. IF the parent ServiceClass does not exist, THEN THE Server SHALL return HTTP 400 with `error.code = "VALIDATION_FAILED"` and `error.field = "spec.serviceClassName"` and a message identifying the missing ServiceClass.
3. THE parent existence check SHALL occur in the API/service layer via a narrow lookup interface; the ServicePlanRegistry SHALL NOT depend on the ServiceClassRegistry.
4. `spec.serviceClassName` SHALL be immutable after creation. PUT SHALL require body `spec.serviceClassName` to match the path segment.
5. WHEN a `PUT` update passes validation and identity checks, THE Server SHALL verify the parent ServiceClass still exists; if absent, THE Server SHALL return HTTP 400 `VALIDATION_FAILED` with `error.field = "spec.serviceClassName"`.

---

### Requirement 10: ServiceClass Delete Blocking

**User Story:** As a catalog operator, I want ServiceClass deletion blocked while ServicePlans exist under it.

#### Acceptance Criteria

1. WHEN `DELETE /v1/service-classes/{name}` is requested and one or more ServicePlans have `spec.serviceClassName == name`, THE Server SHALL return HTTP 409 with `error.code = "DELETE_BLOCKED"` and a message identifying `"ServicePlan"` as the blocking kind.
2. THE blocker SHALL query the ServicePlanRegistry via `CountByServiceClass` and block when count > 0.
3. WHEN the ServiceClass has zero ServicePlans, THE delete SHALL proceed and return 204.
4. THE blocker SHALL use a narrow blocker interface consistent with existing patterns; no generic blocker framework.

---

### Requirement 11: Update Identity Consistency

**User Story:** As a catalog operator, I want PUT requests to enforce path/body identity matching.

#### Acceptance Criteria

1. FOR `PUT /v1/service-classes/{name}`, THE Server SHALL require body `metadata.name` present and equal to the path `{name}`. Absent, empty, or mismatched → HTTP 400 `VALIDATION_FAILED`, `error.field = "metadata.name"`.
2. FOR `PUT /v1/service-plans/{serviceClassName}/{name}`, THE Server SHALL require body `spec.serviceClassName` == path `{serviceClassName}` and body `metadata.name` == path `{name}`. Mismatch/absent → HTTP 400 `VALIDATION_FAILED` with the appropriate field (`spec.serviceClassName` or `metadata.name`).
3. THE registry update SHALL preserve immutable identity fields and SHALL NOT move a ServicePlan between ServiceClasses.
4. THE Server SHALL preserve or reset server-owned `apiVersion`, `kind`, and `status` on update; client input SHALL NOT change server-owned or immutable fields.
5. Mutable fields include: `metadata.labels` and `metadata.annotations` (where the existing `Metadata` struct supports them), `spec.displayName`, `spec.description`, `spec.category`(ServiceClass)/`spec.tier`(ServicePlan), `spec.lifecycle`, `spec.tags`, `spec.parameters`(ServicePlan), `spec.provider`/`spec.defaultPlanName`(ServiceClass). `metadata.displayName` is NOT a mutable field (display naming lives in `spec.displayName`).

---

### Requirement 12: Operation Emission (FEATURE-0005 Integration)

**User Story:** As a catalog operator, I want catalog lifecycle actions recorded as Operations.

#### Acceptance Criteria

1. AFTER a successful ServiceClass create/update/delete, THE Server SHALL emit an Operation of type `CreateServiceClass`, `UpdateServiceClass`, or `DeleteServiceClass` respectively, with `resourceKind = "ServiceClass"`.
2. AFTER a successful ServicePlan create/update/delete, THE Server SHALL emit an Operation of type `CreateServicePlan`, `UpdateServicePlan`, or `DeleteServicePlan` respectively, with `resourceKind = "ServicePlan"`.
3. THE resource kind constants (`ServiceClassKind`, `ServicePlanKind`) SHALL be added to `internal/resources` following the existing pattern.
4. FEATURE-0006 SHALL extend the FEATURE-0005 operation type constants with: `CreateServiceClass`, `UpdateServiceClass`, `DeleteServiceClass`, `CreateServicePlan`, `UpdateServicePlan`, `DeleteServicePlan`.
5. FEATURE-0006 SHALL extend `OperationSpec` with two optional catalog-reference fields: `serviceClassName` (`json:"serviceClassName,omitempty"`) and `servicePlanName` (`json:"servicePlanName,omitempty"`). These are used only for ServiceClass and ServicePlan Operation records and remain empty/omitted for other resource kinds.
6. ServiceClass Operation records SHALL set `serviceClassName` to the class name. ServicePlan Operation records SHALL set `serviceClassName` to the parent class name and `servicePlanName` to the plan name.
7. Emission SHALL use the nil-safe FEATURE-0005 emitter; emission failure SHALL NOT affect the primary API response.
8. THE Server SHALL NOT emit an Operation on failed validation, duplicate create, missing parent, not-found, or delete-blocked cases.

---

### Requirement 13: Security and Privacy

**User Story:** As a security-conscious operator, I want catalog resources to never store secrets.

#### Acceptance Criteria

1. ServiceClass and ServicePlan SHALL NOT store secrets, tokens, credentials, or passwords.
2. `spec.parameters` SHALL be catalog metadata only. Forbidden secret-bearing parameter key patterns are exactly those defined in Requirement 4.6, matched case-insensitively. The plain substring `key` by itself SHALL NOT be rejected unless it appears as part of `apiKey`, `accessKey`, `secretKey`, or `privateKey`.
3. THE Server SHALL NOT store raw request bodies and SHALL NOT echo raw bodies in error responses.
4. THE Server SHALL NOT log secrets or raw bodies.

---

### Requirement 14: Tests

**User Story:** As a developer, I want comprehensive tests for both resources.

#### Acceptance Criteria

1. Validation unit tests: valid names accepted; invalid/empty/long names rejected; invalid category/tier/lifecycle rejected; forbidden parameter keys rejected (case-insensitive).
2. Registry unit tests (both registries): Create stores; duplicate → ErrAlreadyExists (original unchanged); Get by key; missing → ErrNotFound; List sorted; empty → non-nil `[]`; Update mutable fields only; Update missing → ErrNotFound; Delete removes; Delete missing → ErrNotFound; ServicePlan CountByServiceClass correct; same plan name under different classes succeeds.
3. Handler tests (both resources): POST 201/409/400 (invalid fields, status key, bad JSON, unknown field); POST 413 oversized; ServicePlan POST 400 missing parent; GET 200/404/400; wrong path shape → 404; list sorted/empty; PUT 200/404/400 (identity mismatch); DELETE 204/404.
4. Delete-blocking tests: ServiceClass delete with ServicePlans → 409 DELETE_BLOCKED; with zero plans → 204.
5. Operation emission tests: successful create/update/delete of each resource records the correct Operation type and resource kind; failed actions emit nothing; emission failure does not change the primary response (table-driven where practical).
6. `go test -race ./...` with 10+ goroutines produces no race reports.
7. ALL tests deterministic; no external dependencies.

---

### Requirement 15: Property-Based Tests

**User Story:** As a developer, I want property tests for validation and registries.

#### Acceptance Criteria

1. Validation package property tests (`testing/quick`, `Config{MaxCount: 100}`): valid DNS-label names accepted; arbitrary invalid strings rejected.
2. Registry property tests: Create/Get round-trip preserves data; List sort invariant; deep-copy immutability; duplicate-create idempotent error (original unchanged).
3. Each property test tagged `// Feature: serviceclass-serviceplan, Property N: <title>`.

---

### Requirement 16: Non-Goals

**User Story:** As an architect, I want clear scope boundaries.

#### Acceptance Criteria

1. NO ServiceInstance or ServiceBinding.
2. NO service provisioning, datastore provisioning, Kubernetes operators, Crossplane, Kratix, or GitOps.
3. NO Plugin registry, Capability registry, or capability matching.
4. NO ServiceOps execution or plugin execution.
5. NO async workflows, approval flows, queues, or background workers.
6. NO pricing/billing engine, quota enforcement, tenant subscription model, or marketplace publishing.
7. NO persistence/database storage, auth/RBAC, AI automation, or UI.
8. NO Go 1.22 wildcard routing; no new external dependencies.
9. NO secrets stored in catalog resources.

---

### Requirement 17: Edge Cases

**User Story:** As a developer, I want edge-case behavior defined.

#### Acceptance Criteria

1. WHEN a mutating action fails (validation, duplicate, missing parent, not-found, delete-blocked), THE Server SHALL NOT emit an Operation.
2. WHEN a ServicePlan item path has fewer or more than two segments, THE Server SHALL return HTTP 404.
3. WHEN a ServiceClass item path has more than one segment, THE Server SHALL return HTTP 404.
4. WHEN `spec.parameters` is nil or empty, THE ServicePlan SHALL be accepted (parameters are optional).
5. WHEN the same ServicePlan `metadata.name` is created under two different ServiceClasses, both SHALL succeed (composite identity).
6. WHEN deleting a ServiceClass referenced by a ServicePlan whose lifecycle is `Retired`, THE delete SHALL still be blocked (lifecycle does not exempt blocking).
7. WHEN a nil emitter is used (isolated handler tests), emission SHALL be skipped gracefully without panic.

---

## Compatibility with FEATURE-0005 Operation Emission

FEATURE-0006 reuses the FEATURE-0005 `OperationEmitter` interface and nil-safe `emitOperation`
helper. ServiceClass and ServicePlan handlers receive the emitter via constructor injection
(nil-safe), consistent with Org/OU/Tenant/Project handlers. Six new operation type constants
(`CreateServiceClass`…`DeleteServicePlan`) and two resource kind constants (`ServiceClassKind`,
`ServicePlanKind`) are added to `internal/resources`. Emission occurs only after a successful
mutating action, uses the request context for the request ID, and never alters the primary response.

## Design Questions

> These questions are OPEN and MUST be resolved in `design.md` before implementation.

1. **Parent lookup interface:** Define a narrow `ServiceClassLookup` interface (mirroring
   `OrganizationUnitLookup`/`TenantLookup`) injected into the ServicePlan handler, satisfied by the
   existing `ServiceClassRegistry`. Confirm placement in `internal/api` vs `internal/registry`.

2. **Delete-blocker interface shape:** ServiceClass delete takes a single path segment, so the
   blocker signature is `BlockedByServiceClassChildren(ctx, serviceClassName) ([]BlockedBy, error)`
   implemented by a ServicePlan-backed checker. Confirm naming consistency with existing blockers.

3. **Parameters value type:** `map[string]string` is proposed for `spec.parameters`. Confirm whether
   a structured value type is needed later; Phase 1 stays with string values to avoid secret-bearing
   nested structures.
