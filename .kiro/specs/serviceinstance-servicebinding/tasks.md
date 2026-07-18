# Tasks: FEATURE-0008 ServiceInstance and ServiceBinding

## Task 1: ServiceInstance and ServiceBinding Resource Structs

### Objective
Define the ServiceInstance and ServiceBinding resource structs, constants, and operation type extensions in `internal/resources`.

### Files
- `internal/resources/serviceinstance.go` (new)
- `internal/resources/servicebinding.go` (new)
- `internal/resources/operation.go` (modify — add fields and constants)

### Notes
- Follow canonical `metadata/spec/status` shape.
- ServiceInstanceSpec has: OrganizationRef, OrganizationUnitRef (omitempty), TenantRef, ProjectRef, ServiceClassRef, ServicePlanRef, Parameters (map[string]string, omitempty).
- ServiceInstanceStatus has: Phase, Message (omitempty).
- ServiceBindingSpec has: ServiceInstanceRef, ConsumerRef (*ConsumerRef), BindingType.
- ConsumerRef has: Kind, Name.
- ServiceBindingStatus has: Phase, Message (omitempty), SecretRef (omitempty).
- Add constants: ServiceInstanceAPIVersion, ServiceInstanceKind, ServiceBindingAPIVersion, ServiceBindingKind, BindingTypeCredentials.
- Add to operation.go: ServiceInstanceName and ServiceBindingName fields in OperationSpec (omitempty).
- Add operation type constants: OpCreateServiceInstance, OpUpdateServiceInstance, OpDeleteServiceInstance, OpCreateServiceBinding, OpDeleteServiceBinding.

### Tests
- Compile check only (structs have no behavior). Verified by `go vet ./...`.

### Acceptance Criteria
- Both structs compile with correct JSON tags.
- OperationSpec has new optional fields.
- New operation type constants exist.
- `go vet ./...` passes.
- `go test ./...` passes (existing tests unbroken).

### Commit Message
```
feat(resources): add ServiceInstance and ServiceBinding structs and operation constants
```

---

## Task 2: ServiceInstance Validation

### Objective
Implement pure validation functions for ServiceInstance resources.

### Files
- `internal/validation/serviceinstance.go` (new)
- `internal/validation/serviceinstance_test.go` (new)

### Notes
- `ValidateServiceInstance(si resources.ServiceInstance) []resources.FieldError` — pure, no I/O.
- `ValidateServiceInstancePathSegment(name string) []resources.FieldError` — validates DNS-label for path segment.
- Validate: metadata.name (required, DNS-label, ≤63), spec.organizationRef (required, DNS-label, ≤63), spec.organizationUnitRef (optional, if non-empty must be DNS-label ≤63), spec.tenantRef, spec.projectRef, spec.serviceClassRef, spec.servicePlanRef (all required, DNS-label, ≤63).
- No validation of spec.parameters in Phase 1.
- Reuse existing `IsValidDNSLabel` helper from validation package.

### Tests
- Valid ServiceInstance → no errors.
- Missing/empty metadata.name → error with field "metadata.name".
- Invalid DNS-label names (uppercase, special chars, leading/trailing hyphen, >63 chars) → error.
- Missing each required spec field → correct field error.
- Optional organizationUnitRef empty → accepted.
- Optional organizationUnitRef invalid → error with field "spec.organizationUnitRef".
- Parameters nil or empty map → accepted.
- ValidateServiceInstancePathSegment valid → no errors.
- ValidateServiceInstancePathSegment invalid → error.

### Acceptance Criteria
- All validation tests pass.
- `go vet ./...` passes.
- `go test ./internal/validation/...` passes.

### Commit Message
```
feat(validation): add ServiceInstance validation with tests
```

---

## Task 3: ServiceInstance Validation Property Tests

### Objective
Add property-based tests for ServiceInstance validation using `testing/quick`.

### Files
- `internal/validation/serviceinstance_property_test.go` (new)

### Notes
- Use `testing/quick` with `Config{MaxCount: 100}`.
- Property: valid DNS-label names accepted by ValidateServiceInstance.
- Property: arbitrary invalid strings rejected.
- Tag each test: `// Feature: serviceinstance-servicebinding, Property N: <title>`.
- Follow existing pattern from `serviceclass_property_test.go`.

### Tests
- Property tests pass with MaxCount 100.

### Acceptance Criteria
- Property tests compile and pass.
- `go test ./internal/validation/...` passes.

### Commit Message
```
test(validation): add ServiceInstance validation property tests
```

---

## Task 4: ServiceBinding Validation

### Objective
Implement pure validation functions for ServiceBinding resources.

### Files
- `internal/validation/servicebinding.go` (new)
- `internal/validation/servicebinding_test.go` (new)

### Notes
- `ValidateServiceBinding(sb resources.ServiceBinding) []resources.FieldError` — pure, no I/O.
- `ValidateServiceBindingPathSegment(name string) []resources.FieldError` — validates DNS-label for path segment.
- Validate: metadata.name (required, DNS-label, ≤63), spec.serviceInstanceRef (required, DNS-label, ≤63), spec.consumerRef (required, non-nil), spec.consumerRef.kind (required, non-empty, no enum restriction), spec.consumerRef.name (required, DNS-label, ≤63), spec.bindingType (must be "credentials").
- Reuse existing `IsValidDNSLabel` helper.

### Tests
- Valid ServiceBinding → no errors.
- Missing/empty metadata.name → error.
- Invalid DNS-label names → error.
- Missing spec.serviceInstanceRef → error.
- Nil spec.consumerRef → error with field "spec.consumerRef".
- Empty spec.consumerRef.kind → error with field "spec.consumerRef.kind".
- Invalid spec.consumerRef.name → error with field "spec.consumerRef.name".
- bindingType "credentials" → accepted.
- bindingType "endpoint" or any other value → error.
- Any non-empty consumerRef.kind (e.g., "Application", "Job") → accepted.
- ValidateServiceBindingPathSegment valid/invalid → correct results.

### Acceptance Criteria
- All validation tests pass.
- `go test ./internal/validation/...` passes.

### Commit Message
```
feat(validation): add ServiceBinding validation with tests
```

---

## Task 5: ServiceBinding Validation Property Tests

### Objective
Add property-based tests for ServiceBinding validation using `testing/quick`.

### Files
- `internal/validation/servicebinding_property_test.go` (new)

### Notes
- Use `testing/quick` with `Config{MaxCount: 100}`.
- Property: valid DNS-label names accepted.
- Property: arbitrary invalid strings rejected.
- Property: valid bindingType "credentials" accepted; invalid rejected.
- Tag each test: `// Feature: serviceinstance-servicebinding, Property N: <title>`.

### Tests
- Property tests pass with MaxCount 100.

### Acceptance Criteria
- Property tests compile and pass.
- `go test ./internal/validation/...` passes.

### Commit Message
```
test(validation): add ServiceBinding validation property tests
```

---

## Task 6: ServiceInstance Registry

### Objective
Implement the in-memory ServiceInstance registry with CRUD, filters, and counter methods.

### Files
- `internal/registry/serviceinstance_registry.go` (new)
- `internal/registry/serviceinstance_registry_test.go` (new)

### Notes
- Define `ServiceInstanceRegistryIface` interface with: CreateServiceInstance, GetServiceInstance, ListServiceInstances (tenantRef, projectRef filters), UpdateServiceInstance, DeleteServiceInstance, CountByServicePlan, CountByProject.
- Define `ServiceInstanceLookup` interface (single method: GetServiceInstance).
- Implement `ServiceInstanceRegistry` struct: `sync.RWMutex`, `map[string]resources.ServiceInstance`.
- Deep-copy on Create/Get/List/Update (including Parameters map, Labels, Annotations).
- UpdateServiceInstance preserves: apiVersion, kind, status, all immutable spec fields (organizationRef, organizationUnitRef, tenantRef, projectRef, serviceClassRef, servicePlanRef). Replaces: parameters, labels, annotations, displayName.
- ListServiceInstances: sorted by metadata.name ascending; filters by tenantRef/projectRef (AND logic); empty → non-nil slice.
- CountByServicePlan: count where spec.serviceClassRef AND spec.servicePlanRef match.
- CountByProject: count where all four governance refs match exactly (including empty OU).
- Sentinel errors: ErrAlreadyExists, ErrNotFound (reuse existing from registry.go).
- All methods accept context.Context as first param.

### Tests
- Create stores; Get returns deep copy.
- Duplicate → ErrAlreadyExists; original unchanged.
- Get missing → ErrNotFound.
- List sorted; empty → non-nil [].
- List with tenantRef filter; with projectRef; with both (AND); with none.
- Update mutable fields only.
- Update preserves stored status unchanged.
- Update missing → ErrNotFound.
- Delete removes entry.
- Delete missing → ErrNotFound.
- CountByServicePlan correct count.
- CountByServicePlan no false positives (same plan name, different ServiceClass).
- CountByProject correct count (all four refs).
- CountByProject no false positives (same project name, different tenant).
- CountByProject with empty OU does not match non-empty OU.
- Duplicate name across different governance refs → ErrAlreadyExists (global uniqueness).
- Deep-copy immutability: mutating returned value does not affect registry.

### Acceptance Criteria
- All registry unit tests pass.
- `go test ./internal/registry/...` passes.
- `go vet ./...` passes.

### Commit Message
```
feat(registry): add ServiceInstance in-memory registry with tests
```

---

## Task 7: ServiceInstance Registry Race and Property Tests

### Objective
Add concurrency race tests and property-based tests for ServiceInstance registry.

### Files
- `internal/registry/serviceinstance_registry_race_test.go` (new)
- `internal/registry/serviceinstance_registry_property_test.go` (new)

### Notes
- Race test: 10+ goroutines performing concurrent Create/Get/List/Update/Delete/CountByServicePlan/CountByProject. No race reports under `go test -race`.
- Property tests: Create/Get round-trip preserves data; List sort invariant; deep-copy immutability; duplicate-create idempotent error; filter correctness; CountByServicePlan correctness; CountByProject correctness.
- Use `testing/quick` with `Config{MaxCount: 100}`.

### Tests
- `go test -race ./internal/registry/...` passes with no race reports.
- Property tests pass.

### Acceptance Criteria
- Race test produces no data race reports.
- Property tests pass.
- `go test -race ./internal/registry/...` passes.

### Commit Message
```
test(registry): add ServiceInstance registry race and property tests
```

---

## Task 8: ServiceBinding Registry

### Objective
Implement the in-memory ServiceBinding registry with CRUD, filter, and counter methods.

### Files
- `internal/registry/servicebinding_registry.go` (new)
- `internal/registry/servicebinding_registry_test.go` (new)

### Notes
- Define `ServiceBindingRegistryIface` interface: CreateServiceBinding, GetServiceBinding, ListServiceBindings (serviceInstanceRef filter), DeleteServiceBinding, CountByServiceInstance.
- Define `ServiceBindingInstanceBlocker` interface: `CountByServiceInstance(ctx, instanceName) (int, error)`.
- Implement `ServiceBindingRegistry` struct: `sync.RWMutex`, `map[string]resources.ServiceBinding`.
- Deep-copy on Create/Get/List (including ConsumerRef pointer, Labels, Annotations).
- DeleteServiceBinding returns `error` only (no resource returned).
- ListServiceBindings: sorted by metadata.name, filter by serviceInstanceRef; empty → non-nil slice.
- CountByServiceInstance: count where spec.serviceInstanceRef matches.
- Sentinel errors: ErrAlreadyExists, ErrNotFound.
- All methods accept context.Context.

### Tests
- Create stores; Get returns deep copy.
- Duplicate → ErrAlreadyExists.
- Duplicate name referencing different ServiceInstance → ErrAlreadyExists (global uniqueness).
- Get missing → ErrNotFound.
- List sorted; empty → non-nil [].
- List with serviceInstanceRef filter; without filter.
- Delete removes.
- Delete missing → ErrNotFound.
- CountByServiceInstance correct count.
- Deep-copy immutability (including ConsumerRef pointer).

### Acceptance Criteria
- All registry unit tests pass.
- `go test ./internal/registry/...` passes.
- `go vet ./...` passes.

### Commit Message
```
feat(registry): add ServiceBinding in-memory registry with tests
```

---

## Task 9: ServiceBinding Registry Race and Property Tests

### Objective
Add concurrency race tests and property-based tests for ServiceBinding registry.

### Files
- `internal/registry/servicebinding_registry_race_test.go` (new)
- `internal/registry/servicebinding_registry_property_test.go` (new)

### Notes
- Race test: 10+ goroutines performing concurrent Create/Get/List/Delete/CountByServiceInstance. No race reports under `go test -race`.
- Property tests: Create/Get round-trip; List sort invariant; deep-copy immutability; duplicate-create idempotent error; filter correctness; CountByServiceInstance correctness.
- Use `testing/quick` with `Config{MaxCount: 100}`.

### Tests
- `go test -race ./internal/registry/...` passes.
- Property tests pass.

### Acceptance Criteria
- Race test produces no data race reports.
- Property tests pass.
- `go test -race ./internal/registry/...` passes.

### Commit Message
```
test(registry): add ServiceBinding registry race and property tests
```

---

## Task 10: Lookup Interfaces (ProjectLookup, ServicePlanLookup, CapabilityLookup)

### Objective
Add narrow lookup interfaces needed by ServiceInstance handler to existing registry files.

### Files
- `internal/registry/project_registry.go` (modify — add ProjectLookup interface)
- `internal/registry/serviceplan_registry.go` (modify — add ServicePlanLookup interface)
- `internal/registry/capability_registry.go` (modify — add CapabilityLookup interface)

### Notes
- `ProjectLookup` interface: `GetProject(ctx, orgName, ouName, tenantName, name string) (resources.Project, error)`. Satisfied by existing `*ProjectRegistry`.
- `ServicePlanLookup` interface: `GetServicePlan(ctx, serviceClassName, name string) (resources.ServicePlan, error)`. Satisfied by existing `*ServicePlanRegistry`.
- `CapabilityLookup` interface: `HasActiveCapabilityForServiceClass(ctx, serviceClassRef string) (bool, error)`. Satisfied by a new implementation (Task 11).
- Place each interface in the respective resource's registry file (resource-side placement).

### Tests
- Compile check (interface declarations only). Existing tests must still pass.

### Acceptance Criteria
- Interfaces compile.
- Existing `*ProjectRegistry` satisfies `ProjectLookup`.
- Existing `*ServicePlanRegistry` satisfies `ServicePlanLookup`.
- `go vet ./...` passes.
- `go test ./...` passes (no regressions).

### Commit Message
```
feat(registry): add ProjectLookup, ServicePlanLookup, and CapabilityLookup interfaces
```

---

## Task 11: CapabilityLookup Implementation

### Objective
Implement the concrete `CapabilityLookupImpl` that checks whether an active Capability exists for a ServiceClass.

### Files
- `internal/registry/capability_lookup.go` (new)
- `internal/registry/capability_lookup_test.go` (new)

### Notes
- `CapabilityLookupImpl` struct holds a `CapabilityRegistryIface`.
- `NewCapabilityLookup(reg CapabilityRegistryIface) *CapabilityLookupImpl`.
- `HasActiveCapabilityForServiceClass(ctx, serviceClassRef) (bool, error)`: calls `ListCapabilities(ctx, "", serviceClassRef)`, iterates, returns true if any capability has `status.phase == "Active"` AND `spec.supported == true`.
- Encapsulates "active" definition; ServiceInstance handler only sees `(bool, error)`.

### Tests
- No capabilities registered → returns false.
- Capability registered but inactive (phase != "Active") → returns false.
- Capability registered but unsupported (supported == false) → returns false.
- Capability registered, active, supported → returns true.
- Multiple capabilities, only one active+supported → returns true.
- Registry error propagated → returns error.

### Acceptance Criteria
- All tests pass.
- `go test ./internal/registry/...` passes.
- `go vet ./...` passes.

### Commit Message
```
feat(registry): add CapabilityLookupImpl with tests
```

---

## Task 12: ServicePlan Instance Blocker

### Objective
Implement the ServicePlan delete-blocker that prevents deletion while ServiceInstances reference the plan.

### Files
- `internal/registry/serviceplan_instance_blocker.go` (new)
- `internal/registry/serviceplan_instance_blocker_test.go` (new)

### Notes
- `ServicePlanInstanceBlocker` interface: `BlockedByServicePlanInstances(ctx, serviceClassName, planName string) ([]BlockedBy, error)`.
- `ServiceInstancePlanBlockerChecker` struct wraps `ServiceInstanceRegistryIface`.
- Implementation: calls `CountByServicePlan(ctx, serviceClassName, planName)`; if count > 0, returns `[]BlockedBy{{Kind: "ServiceInstance", Count: count}}`.
- Reuse existing `BlockedBy` struct from `internal/registry/blocker.go`.

### Tests
- No instances referencing plan → returns nil (not blocked).
- One instance referencing plan → returns blocked with kind "ServiceInstance", count 1.
- Multiple instances → correct count.
- Different ServiceClass with same plan name → not counted (no false positive).

### Acceptance Criteria
- All blocker tests pass.
- `go test ./internal/registry/...` passes.

### Commit Message
```
feat(registry): add ServicePlan instance delete-blocker with tests
```

---

## Task 13: Project Instance Blocker

### Objective
Implement the Project delete-blocker that prevents deletion while ServiceInstances exist under the project.

### Files
- `internal/registry/project_instance_blocker.go` (new)
- `internal/registry/project_instance_blocker_test.go` (new)

### Notes
- `ProjectInstanceBlocker` interface: `BlockedByProjectInstances(ctx, orgName, ouName, tenantName, projectName string) ([]BlockedBy, error)`.
- `ServiceInstanceProjectBlockerChecker` struct wraps `ServiceInstanceRegistryIface`.
- Implementation: calls `CountByProject(ctx, orgName, ouName, tenantName, projectName)`; if count > 0, returns `[]BlockedBy{{Kind: "ServiceInstance", Count: count}}`.
- Reuse existing `BlockedBy` struct.

### Tests
- No instances under project → returns nil (not blocked).
- One instance under project → returns blocked with kind "ServiceInstance", count 1.
- Multiple instances → correct count.
- Different tenant with same project name → not counted (no false positive).
- Empty OU vs non-empty OU → correct isolation.

### Acceptance Criteria
- All blocker tests pass.
- `go test ./internal/registry/...` passes.

### Commit Message
```
feat(registry): add Project instance delete-blocker with tests
```

---

## Task 14: ServiceInstance Decode Function

### Objective
Implement safe JSON decoding for ServiceInstance requests.

### Files
- `internal/api/serviceinstance_decode.go` (new)

### Notes
- `safeDecodeServiceInstance(r *http.Request) (resources.ServiceInstance, error)` following existing decode pattern.
- 1 MiB limit via `http.MaxBytesReader`.
- `json.Decoder` with `DisallowUnknownFields`.
- Reject if top-level `status` key is present in raw JSON (check via map decode or existing pattern).
- Return structured errors for: oversized body, malformed JSON, unknown fields, empty body, status key present.
- Follow existing pattern from `serviceplan_decode.go` or `capability_decode.go`.

### Tests
- Verified via handler tests in Task 17. No separate test file needed (matches existing pattern).

### Acceptance Criteria
- Decode function compiles.
- `go vet ./...` passes.
- `go test ./...` passes (no regressions).

### Commit Message
```
feat(api): add ServiceInstance safe JSON decode function
```

---

## Task 15: ServiceBinding Decode Function

### Objective
Implement safe JSON decoding for ServiceBinding requests.

### Files
- `internal/api/servicebinding_decode.go` (new)

### Notes
- `safeDecodeServiceBinding(r *http.Request) (resources.ServiceBinding, error)` following existing decode pattern.
- 1 MiB limit via `http.MaxBytesReader`.
- `json.Decoder` with `DisallowUnknownFields`.
- Reject if top-level `status` key is present (applies to POST only; ServiceBinding has no PUT).
- Return structured errors for: oversized body, malformed JSON, unknown fields, empty body, status key present.
- Follow existing decode patterns.

### Tests
- Verified via handler tests in Task 18. No separate test file needed (matches existing pattern).

### Acceptance Criteria
- Decode function compiles.
- `go vet ./...` passes.
- `go test ./...` passes (no regressions).

### Commit Message
```
feat(api): add ServiceBinding safe JSON decode function
```

---

## Task 16: ServiceInstance Handler

### Objective
Implement the ServiceInstance HTTP handler with full CRUD, reference validation, capability warning, and operation emission.

### Files
- `internal/api/serviceinstance_handler.go` (new)

### Notes
- `ServiceInstanceHandler` struct with fields: registry (ServiceInstanceRegistryIface), orgLookup, ouLookup, tenantLookup, projectLookup, serviceClassLookup, servicePlanLookup, capabilityLookup, bindingBlocker (ServiceBindingInstanceBlocker), emitter (OperationEmitter), logger (*log.Logger).
- Constructor: `NewServiceInstanceHandler(...)`.
- `HandleCollection` dispatches POST/GET on `/v1/service-instances`.
- `HandleItem` trims prefix `/v1/service-instances/`, splits on `/`, expects exactly 1 non-empty segment or returns 404.
- **Create flow**: decode → validate → org lookup → OU lookup (if non-empty) → tenant lookup → project lookup → serviceClass lookup → servicePlan lookup (verifies plan belongs to class) → capability warning (log only) → set server fields → registry create → emit operation → 201.
- **Get flow**: validate path → registry get → 200 or 404.
- **List flow**: read query params tenantRef/projectRef → registry list → `{"items": [...]}`.
- **Update flow**: validate path → decode (status rejection) → identity check (path==body name) → field validate → get stored → immutability check (all 6 fields) → registry update → emit operation → 200.
- **Delete flow**: validate path → binding blocker check → registry delete → emit operation → 204.
- Error mapping follows design.md error table.
- Operation emission is nil-safe (uses existing emitOperation helper).
- Capability warning: if HasActiveCapabilityForServiceClass returns false, log structured warning but proceed.

### Tests
- Verified in Tasks 17 and 19. This task creates the handler code only.

### Acceptance Criteria
- Handler compiles.
- `go vet ./...` passes.
- `go test ./...` passes (no regressions from existing tests).

### Commit Message
```
feat(api): add ServiceInstance HTTP handler with CRUD and reference validation
```

---

## Task 17: ServiceInstance Handler Tests

### Objective
Comprehensive handler tests for ServiceInstance CRUD, validation, references, immutability, and delete-blocking.

### Files
- `internal/api/serviceinstance_handler_test.go` (new)

### Notes
- Use mock/fake implementations of all lookup interfaces and registries (in-test or test helpers).
- Follow existing handler test patterns from `serviceplan_handler_test.go` or `project_handler_test.go`.
- Test cases:
  - POST 201 valid; response has apiVersion, kind, status.phase="Ready", status.message.
  - POST 409 duplicate name.
  - POST 409 duplicate name across different governance refs (global uniqueness).
  - POST 400: invalid fields, status key present, bad JSON, unknown field, oversized body.
  - POST 400: missing Organization ref (not found).
  - POST 400: missing Tenant ref (not found).
  - POST 400: missing Project ref (not found).
  - POST 400: missing ServiceClass ref (not found).
  - POST 400: missing ServicePlan ref (not found).
  - POST 400: ServicePlan not matching ServiceClass.
  - POST 400: governance hierarchy inconsistency.
  - GET 200 found; 404 missing; 400 invalid path segment.
  - GET list sorted; empty `{"items": []}`.
  - GET list with query params filtering.
  - PUT 200 success (mutable fields updated).
  - PUT 200 preserves stored status unchanged.
  - PUT 404 missing resource.
  - PUT 400 path/body name mismatch.
  - PUT 400 status key present → VALIDATION_FAILED field=status.
  - PUT 400 immutable field change (one test per immutable field: organizationRef, organizationUnitRef, tenantRef, projectRef, serviceClassRef, servicePlanRef).
  - DELETE 204 success.
  - DELETE 404 missing.
  - DELETE 409 blocked by ServiceBindings.
  - Wrong path shape (extra segments) → 404.

### Acceptance Criteria
- All handler tests pass.
- `go test ./internal/api/...` passes.

### Commit Message
```
test(api): add ServiceInstance handler tests
```

---

## Task 18: ServiceBinding Handler

### Objective
Implement the ServiceBinding HTTP handler with Create/Get/List/Delete and PUT→405.

### Files
- `internal/api/servicebinding_handler.go` (new)
- `internal/api/servicebinding_handler_test.go` (new)

### Notes
- `ServiceBindingHandler` struct with fields: registry (ServiceBindingRegistryIface), instanceLookup (ServiceInstanceLookup), emitter (OperationEmitter).
- Constructor: `NewServiceBindingHandler(...)`.
- `HandleCollection` dispatches POST/GET on `/v1/service-bindings`.
- `HandleItem` trims prefix `/v1/service-bindings/`, splits, expects 1 non-empty segment or 404.
- **Create flow**: decode → validate → instanceLookup.GetServiceInstance (400 if not found) → set server fields (apiVersion, kind, status.phase="Ready", status.secretRef="stub-secret-ref") → registry create → emit operation (serviceInstanceName set to referenced instance) → 201.
- **Get flow**: validate path → registry get → 200 or 404.
- **List flow**: read query param serviceInstanceRef → registry list → `{"items": [...]}`.
- **Delete flow**: validate path → registry get (capture serviceInstanceRef for emission) → registry delete → emit operation → 204.
- **PUT → 405**: return METHOD_NOT_ALLOWED with message "ServiceBinding does not support update; delete and recreate instead", regardless of resource existence.
- Test cases:
  - POST 201 valid; response has secretRef="stub-secret-ref".
  - POST 409 duplicate.
  - POST 409 duplicate name referencing different ServiceInstance (global uniqueness).
  - POST 400 validation errors, status key, bad JSON, unknown field.
  - POST 400 missing ServiceInstance ref (not found).
  - GET 200/404/400 (invalid path segment).
  - GET list with serviceInstanceRef filter; empty list.
  - PUT → 405 METHOD_NOT_ALLOWED.
  - DELETE 204/404.
  - Wrong path shape → 404.

### Acceptance Criteria
- Handler compiles and all tests pass.
- `go test ./internal/api/...` passes.

### Commit Message
```
feat(api): add ServiceBinding handler with tests
```

---

## Task 19: ServiceInstance Operation Emission Tests

### Objective
Verify operation emission behavior for ServiceInstance lifecycle actions.

### Files
- `internal/api/serviceinstance_emission_test.go` (new)

### Notes
- Follow existing pattern from `plugin_emission_test.go` or `capability_emission_test.go`.
- Test cases:
  - Successful create emits OpCreateServiceInstance with correct resourceKind, resourceName, serviceInstanceName.
  - Successful update emits OpUpdateServiceInstance.
  - Successful delete emits OpDeleteServiceInstance.
  - ServiceInstance operation records include governance fields (organizationName, tenantName, projectName).
  - Failed validation → no emission.
  - Duplicate create → no emission.
  - Missing reference → no emission.
  - Not-found on update/delete → no emission.
  - Delete-blocked → no emission.
  - Nil emitter → no panic, no emission.

### Acceptance Criteria
- All emission tests pass.
- `go test ./internal/api/...` passes.

### Commit Message
```
test(api): add ServiceInstance operation emission tests
```

---

## Task 20: ServiceBinding Operation Emission Tests

### Objective
Verify operation emission behavior for ServiceBinding lifecycle actions.

### Files
- `internal/api/servicebinding_emission_test.go` (new)

### Notes
- Follow existing emission test patterns.
- Test cases:
  - Successful create emits OpCreateServiceBinding with correct resourceKind, resourceName, serviceInstanceName (referenced instance), serviceBindingName.
  - Successful delete emits OpDeleteServiceBinding with serviceInstanceName and serviceBindingName.
  - Failed validation → no emission.
  - Duplicate create → no emission.
  - Missing ServiceInstance ref → no emission.
  - Not-found on delete → no emission.
  - Nil emitter → no panic, no emission.

### Acceptance Criteria
- All emission tests pass.
- `go test ./internal/api/...` passes.

### Commit Message
```
test(api): add ServiceBinding operation emission tests
```

---

## Task 21: Capability Warning Test

### Objective
Verify that ServiceInstance creation logs a warning when no active Capability exists for the ServiceClass, but does not block creation.

### Files
- `internal/api/serviceinstance_handler_test.go` (modify — add capability warning test cases)

### Notes
- Test: ServiceInstance created without matching active Capability → creation succeeds (201), warning logged.
- Test: ServiceInstance created with matching active Capability → creation succeeds (201), no warning.
- The warning is observability-only (server log). Use a test logger or bytes.Buffer to capture and assert log output.
- This is an additive test within the existing handler test file.

### Acceptance Criteria
- Capability warning test passes.
- Creation is not blocked in either case.
- `go test ./internal/api/...` passes.

### Commit Message
```
test(api): add ServiceInstance capability warning tests
```

---

## Task 22: Wire ServicePlan Delete-Blocker into ServicePlan Handler

### Objective
Inject the new ServicePlanInstanceBlocker into the existing ServicePlan handler so that ServicePlan deletion is blocked while ServiceInstances reference it.

### Files
- `internal/api/serviceplan_handler.go` (modify — add instanceBlocker field and call in Delete)
- `internal/api/serviceplan_handler_test.go` (modify — add delete-blocking test cases)

### Notes
- Add `instanceBlocker ServicePlanInstanceBlocker` (or use the interface from `serviceplan_instance_blocker.go`) field to ServicePlanHandler struct.
- In the Delete method, after existing child-blocker checks, call `instanceBlocker.BlockedByServicePlanInstances(ctx, serviceClassName, planName)`. If blocked → 409 DELETE_BLOCKED with message identifying "ServiceInstance" as blocking kind.
- If instanceBlocker is nil, skip check (backward compatibility for isolated tests).
- Add test cases:
  - ServicePlan delete with referencing ServiceInstances → 409 DELETE_BLOCKED.
  - ServicePlan delete with zero referencing ServiceInstances → normal delete (204).

### Acceptance Criteria
- ServicePlan delete-blocking works correctly.
- Existing ServicePlan tests still pass.
- `go test ./internal/api/...` passes.

### Commit Message
```
feat(api): wire ServiceInstance delete-blocker into ServicePlan handler
```

---

## Task 23: Wire Project Delete-Blocker into Project Handler

### Objective
Inject the new ProjectInstanceBlocker into the existing Project handler so that Project deletion is blocked while ServiceInstances exist under it.

### Files
- `internal/api/project_handler.go` (modify — add instanceBlocker field and call in Delete)
- `internal/api/project_handler_test.go` (modify — add delete-blocking test cases)

### Notes
- Add `instanceBlocker ProjectInstanceBlocker` (or use the interface from `project_instance_blocker.go`) field to ProjectHandler struct.
- In the Delete method, after existing child-blocker checks, call `instanceBlocker.BlockedByProjectInstances(ctx, orgName, ouName, tenantName, projectName)`. If blocked → 409 DELETE_BLOCKED with message identifying "ServiceInstance" as blocking kind.
- If instanceBlocker is nil, skip check (backward compatibility for isolated tests).
- Add test cases:
  - Project delete with ServiceInstances → 409 DELETE_BLOCKED.
  - Project delete with zero ServiceInstances → normal delete (204).

### Acceptance Criteria
- Project delete-blocking works correctly.
- Existing Project tests still pass.
- `go test ./internal/api/...` passes.

### Commit Message
```
feat(api): wire ServiceInstance delete-blocker into Project handler
```

---

## Task 24: Server Route Registration

### Objective
Register ServiceInstance and ServiceBinding routes in the server and accept new handlers.

### Files
- `internal/server/server.go` (modify — add handler fields, register routes)

### Notes
- Add `ServiceInstanceHandler` and `ServiceBindingHandler` fields to server config or New() params.
- Register routes:
  - `/v1/service-instances` and `/v1/service-instances/` → ServiceInstanceHandler.HandleCollection
  - `/v1/service-instances/` (item paths) → ServiceInstanceHandler.HandleItem
  - `/v1/service-bindings` and `/v1/service-bindings/` → ServiceBindingHandler.HandleCollection
  - `/v1/service-bindings/` (item paths) → ServiceBindingHandler.HandleItem
- Apply existing middleware (contentType, requestID, logging) to new routes.
- Follow existing routing pattern used by other handlers.
- Do NOT modify existing handler registrations or break existing routes.

### Tests
- Verified via server_test.go and integration in Task 25. May need to update server_test.go constructor calls to pass nil for new handler params.

### Acceptance Criteria
- Server compiles with new routes.
- Existing server tests still pass (may need nil handler params).
- `go vet ./...` passes.
- `go test ./...` passes.

### Commit Message
```
feat(server): register ServiceInstance and ServiceBinding routes
```

---

## Task 25: Main Wiring

### Objective
Wire all new registries, blockers, lookup implementations, and handlers in `cmd/sovrunn-api/main.go`.

### Files
- `cmd/sovrunn-api/main.go` (modify)

### Notes
- Instantiate `ServiceInstanceRegistry`.
- Instantiate `ServiceBindingRegistry`.
- Instantiate `CapabilityLookupImpl` with capability registry.
- Instantiate `ServiceInstancePlanBlockerChecker` with ServiceInstance registry.
- Instantiate `ServiceInstanceProjectBlockerChecker` with ServiceInstance registry.
- Instantiate `ServiceInstanceHandler` with all lookup interfaces:
  - orgLookup: org registry
  - ouLookup: OU registry
  - tenantLookup: tenant registry
  - projectLookup: project registry
  - serviceClassLookup: serviceClass registry
  - servicePlanLookup: servicePlan registry
  - capabilityLookup: capabilityLookupImpl
  - bindingBlocker: ServiceBinding registry (satisfies ServiceBindingInstanceBlocker)
  - emitter: operation emitter
  - logger: server logger
- Instantiate `ServiceBindingHandler` with:
  - registry: ServiceBinding registry
  - instanceLookup: ServiceInstance registry (satisfies ServiceInstanceLookup)
  - emitter: operation emitter
- Pass `instanceBlocker` to ServicePlan handler constructor (new param).
- Pass `instanceBlocker` to Project handler constructor (new param).
- Pass ServiceInstance and ServiceBinding handlers to server `New()`.

### Tests
- Application compiles and starts: `go build ./cmd/sovrunn-api`.
- `go vet ./...` passes.
- `go test ./...` passes.

### Acceptance Criteria
- Full application compiles and links without errors.
- All existing and new tests pass.
- `go build ./cmd/sovrunn-api` produces a working binary.

### Commit Message
```
feat(main): wire ServiceInstance and ServiceBinding handlers and dependencies
```

---

## Task 26: Final Verification and Cleanup

### Objective
Run full verification suite, ensure guardrails are met, clean up artifacts, and confirm clean git status.

### Files
- No new files. Verification and cleanup only.

### Notes
- Run final Docker verification:
  ```bash
  docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./... && go build ./cmd/sovrunn-api'
  ```
- Run guardrails:
  - `rm -f sovrunn-api`
  - `rm -rf bin`
  - Verify no `TODO(FEATURE-0008)` under `internal/` or `cmd/`.
  - Verify no `internal/api` import of `internal/server`.
  - `git status` clean (no untracked or modified files beyond expected commits).
- Verify all acceptance criteria from requirements.md are satisfied:
  - ServiceInstance CRUD works end-to-end.
  - ServiceBinding CRUD works end-to-end (PUT→405).
  - Reference validation blocks invalid references.
  - Governance hierarchy consistency enforced.
  - Catalog reference consistency enforced.
  - Immutability enforced on PUT.
  - Delete-blocking works (ServiceInstance←ServiceBinding, ServicePlan←ServiceInstance, Project←ServiceInstance).
  - Operation emission on all successful mutations.
  - Capability warning logged (no block).
  - Global name uniqueness enforced.
  - No race conditions.

### Tests
- All tests pass under Docker with race detector:
  ```bash
  docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./... && go build ./cmd/sovrunn-api'
  ```

### Acceptance Criteria
- Docker verification passes completely.
- `rm -f sovrunn-api` — no binary left.
- `rm -rf bin` — no bin directory left.
- `grep -r "TODO(FEATURE-0008)" internal/ cmd/` returns nothing.
- `grep -r '"github.com/.*internal/server"' internal/api/` returns nothing (no internal/api import of internal/server).
- `git status` is clean.
- Feature is complete.

### Commit Message
```
chore: final verification and cleanup for FEATURE-0008
```
