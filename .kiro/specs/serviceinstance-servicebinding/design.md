# Design: FEATURE-0008 ServiceInstance and ServiceBinding

## Overview

FEATURE-0008 adds the two service consumption resources that complete the Phase 1 resource model: `ServiceInstance` (a tenant/project-scoped requested service) and `ServiceBinding` (the consumption relationship between a consumer and a ServiceInstance). Both resources follow the canonical `metadata/spec/status` shape, use in-memory registries, emit Operations on successful mutations, and integrate with the full governance hierarchy and catalog via narrow lookup interfaces.

No real infrastructure is provisioned, no credentials are generated, no plugins are executed. ServiceInstance creation validates all governance and catalog references, logs a warning if no active Capability exists for the ServiceClass, sets `status.phase = "Ready"`, and returns. ServiceBinding stores a stub `secretRef`.

## Resolved Design Decisions

### DQ1: Governance Hierarchy Consistency Validation

The handler validates lineage step by step:

1. **Organization existence**: `OrganizationLookup.GetOrganization(ctx, spec.organizationRef)` — must exist.
2. **OrganizationUnit existence + parent match** (when `spec.organizationUnitRef` is non-empty): `OrganizationUnitLookup.GetOrganizationUnit(ctx, spec.organizationRef, spec.organizationUnitRef)` — the existing interface already uses the composite key `(orgName, name)`. If `ErrNotFound`, reject.
3. **Tenant existence + lineage match**: `TenantLookup.GetTenant(ctx, spec.organizationRef, spec.organizationUnitRef, spec.tenantRef)` — the existing `TenantLookup` interface uses the composite key `(orgName, ouName, name)`. A successful return proves the Tenant belongs to the correct Organization and OrganizationUnit.
4. **Project existence + parent match**: `ProjectLookup.GetProject(ctx, spec.organizationRef, spec.organizationUnitRef, spec.tenantRef, spec.projectRef)` — the existing `ProjectRegistryIface.GetProject` uses the composite key `(orgName, ouName, tenantName, name)`. A successful return proves the Project belongs to the correct Tenant.

The lookup interfaces return the full resource struct. FEATURE-0008 only needs existence proof (no field inspection beyond what the composite key already guarantees), but returning the full struct is consistent with the existing interface signatures.

### DQ2: Lookup Interfaces Placement

Each `{Resource}Lookup` interface is declared in the registry file of the resource whose implementation satisfies it (colocated with the provider). This aligns with the pattern established by `ServiceClassLookup` and `PluginLookup`.

**Existing lookup interfaces (actual file locations in the repository):**
- `OrganizationLookup` in `internal/registry/ou_registry.go` (legacy consumer-side placement)
- `OrganizationUnitLookup` in `internal/registry/tenant_registry.go` (legacy consumer-side placement)
- `TenantLookup` in `internal/registry/project_registry.go` (legacy consumer-side placement)
- `ServiceClassLookup` in `internal/registry/serviceclass_registry.go` (resource-side placement)
- `PluginLookup` in `internal/registry/plugin_registry.go` (resource-side placement)

**New interfaces added by FEATURE-0008 (resource-side placement):**
- `ProjectLookup` in `internal/registry/project_registry.go` — satisfied by `*ProjectRegistry`
- `ServicePlanLookup` in `internal/registry/serviceplan_registry.go` — satisfied by `*ServicePlanRegistry`
- `ServiceInstanceLookup` in `internal/registry/serviceinstance_registry.go` — satisfied by `*ServiceInstanceRegistry`
- `CapabilityLookup` in `internal/registry/capability_registry.go` — satisfied by `CapabilityLookupImpl`

Naming convention: `{ResourceKind}Lookup` with a single method (e.g., `Get{ResourceKind}` or `Has...`).

### DQ3: ServicePlan Delete-Blocker Extension

A new blocker interface `ServicePlanInstanceBlocker` is defined in `internal/registry/serviceplan_instance_blocker.go`:

```go
type ServicePlanInstanceBlocker interface {
    BlockedByServicePlanInstances(ctx context.Context, serviceClassName, planName string) ([]BlockedBy, error)
}
```

A concrete implementation `ServiceInstancePlanBlockerChecker` wraps `ServiceInstanceRegistryIface.CountByServicePlan`. Both the interface and the implementation live in `internal/registry/serviceplan_instance_blocker.go`. This is injected into the existing `ServicePlanHandler` as an additional blocker field, consistent with how `PluginChildBlocker` is injected into `PluginHandler`.

### DQ4: Project Delete-Blocker Extension

A new blocker interface `ProjectInstanceBlocker` is defined in `internal/registry/project_instance_blocker.go`:

```go
type ProjectInstanceBlocker interface {
    BlockedByProjectInstances(ctx context.Context, orgName, ouName, tenantName, projectName string) ([]BlockedBy, error)
}
```

A concrete implementation `ServiceInstanceProjectBlockerChecker` wraps `ServiceInstanceRegistryIface.CountByProject`. This is injected into the existing `ProjectHandler` as an additional blocker field.

### DQ5: ServiceBinding Immutability

Confirmed: ServiceBinding does NOT support PUT. The server returns HTTP 405 `METHOD_NOT_ALLOWED` for any PUT request on `/v1/service-bindings/{name}`, regardless of resource existence. Delete-and-recreate is the correct update pattern.

### DQ6: Capability Warning Interface

Interface defined in `internal/registry/capability_registry.go` (colocated with `CapabilityRegistryIface`):

```go
type CapabilityLookup interface {
    HasActiveCapabilityForServiceClass(ctx context.Context, serviceClassRef string) (bool, error)
}
```

The concrete implementation `CapabilityLookupImpl` lives in `internal/registry/capability_lookup.go` and queries the `CapabilityRegistryIface.ListCapabilities` method with `serviceClassRef` filter. A Capability is considered "active" when:
- `status.phase == "Active"` (the server sets this on create)
- `spec.supported == true`

Both conditions must hold. This logic is encapsulated inside the `CapabilityLookupImpl`; the ServiceInstance handler only sees `(bool, error)`.

### DQ7: Stub secretRef Value

Confirmed: `"stub-secret-ref"` is the Phase 1 placeholder. It clearly communicates that no real secret exists.

### DQ8: Filter Query Parameters

Confirmed:
- ServiceInstance list: `tenantRef` and `projectRef` query parameters (AND logic when both present).
- ServiceBinding list: `serviceInstanceRef` query parameter.
- Additional filters (`serviceClassRef`, `servicePlanRef`) are deferred to future phases.

### DQ9: ServiceInstance Governance Field Immutability

**Policy (Phase 1): All governance fields are immutable after creation.**

The following fields in `ServiceInstanceSpec` are immutable:
- `organizationRef`
- `organizationUnitRef`
- `tenantRef`
- `projectRef`
- `serviceClassRef`
- `servicePlanRef`

On PUT, the handler compares each of these fields against the stored entry. If any differ, the handler returns 400 `VALIDATION_FAILED` with the changed field name. This eliminates cross-lineage reparenting risk entirely and is consistent with the immutability patterns established by Tenant (tenantRef/projectRef immutability) and Project resources.

Rationale: reparenting a ServiceInstance across Organizations, OrganizationUnits, Tenants, or Projects would require complex lineage re-validation, audit trail complications, and policy re-evaluation. Phase 1 avoids this complexity. Future phases may introduce controlled migration operations as explicit lifecycle actions with full audit trails, rather than field-level mutations.

**Mutable fields on PUT:**
- `spec.parameters` (user-supplied instance configuration)
- `metadata.labels`
- `metadata.annotations`
- `metadata.displayName`

### DQ10: Empty OrganizationUnitRef Semantics

`spec.organizationUnitRef` is optional. When empty (zero-value `""`), the ServiceInstance belongs directly to the Organization without an intermediate OrganizationUnit.

**Composite key representation:** The empty string is a valid component in composite keys. For example, `CountByProject` uses the four-tuple `(orgRef, ouRef, tenantRef, projectRef)` where `ouRef` may be `""`. Two ServiceInstances with the same `orgRef`, `tenantRef`, and `projectRef` but different `ouRef` values (one empty, one non-empty) are NOT considered to be in the same project scope — they belong to different lineage paths.

**Lookup behavior with empty OU:**
- `OrganizationUnitLookup.GetOrganizationUnit` is SKIPPED when `organizationUnitRef` is empty. No validation error is returned for empty OU.
- `TenantLookup.GetTenant(ctx, orgName, "", tenantName)` uses empty string as the OU component. The underlying TenantRegistry uses composite key `"orgName//tenantName"` (double slash represents empty OU). This must match how the Tenant was originally created.
- `ProjectLookup.GetProject(ctx, orgName, "", tenantName, projectName)` similarly passes empty OU through to the composite key.
- `CountByProject(ctx, orgRef, "", tenantRef, projectRef)` counts only instances where `spec.organizationUnitRef` is also empty.

**Required tests:**
- Create ServiceInstance with empty OU → succeeds if matching Tenant/Project exist with empty OU.
- Create ServiceInstance with non-empty OU → validates OU existence.
- `CountByProject` with empty OU does NOT match instances that have a non-empty OU (even if other refs are identical).
- `CountByProject` with non-empty OU does NOT match instances that have empty OU.
- Lineage validation rejects cross-path: if Tenant was created under OU "finance", a ServiceInstance with empty OU referencing that Tenant fails.

### DQ11: Name Uniqueness Scope

**ServiceInstance and ServiceBinding names are globally unique across the entire API server.**

The `metadata.name` field is the sole key for both registries. Two ServiceInstances cannot have the same name even if they belong to different Organizations, Tenants, or Projects. Two ServiceBindings cannot have the same name even if they reference different ServiceInstances.

This is consistent with the existing behavior of all Phase 1 resources (Organization, Tenant, Project, ServiceClass, Plugin) where `metadata.name` is the global unique identifier.

**Registry enforcement:** `CreateServiceInstance` and `CreateServiceBinding` check the map keyed by `metadata.name` and return `ErrAlreadyExists` (mapped to 409 `RESOURCE_ALREADY_EXISTS`) if the name is already taken, regardless of governance refs.

**Handler error mapping:** Duplicate name on POST → 409 `RESOURCE_ALREADY_EXISTS`. The error message does not expose the existing resource's governance refs (information leakage prevention).

**Required tests:**
- Create ServiceInstance "my-db" under Project A → succeeds.
- Create another ServiceInstance "my-db" under Project B (different org/tenant/project) → 409 `RESOURCE_ALREADY_EXISTS`.
- Create ServiceBinding "my-binding" for Instance X → succeeds.
- Create another ServiceBinding "my-binding" for Instance Y (different instance) → 409 `RESOURCE_ALREADY_EXISTS`.
- The 409 response body does NOT include details about the conflicting resource's governance refs.

## Architecture

### Component Interaction (Create ServiceInstance flow)

```text
Client
  → POST /v1/service-instances
    → contentTypeMiddleware (415 if not application/json)
    → requestIDMiddleware (assigns X-Sovrunn-Request-ID)
    → loggingMiddleware
    → ServiceInstanceHandler.Create
        1. safeDecodeServiceInstance (1 MiB, status rejection, DisallowUnknownFields)
        2. validation.ValidateServiceInstance (pure, field-level)
        3. OrganizationLookup.GetOrganization (existence)
        4. OrganizationUnitLookup.GetOrganizationUnit (existence + org match)
        5. TenantLookup.GetTenant (existence + lineage)
        6. ProjectLookup.GetProject (existence + tenant match)
        7. ServiceClassLookup.GetServiceClass (existence)
        8. ServicePlanLookup.GetServicePlan (existence + serviceClass match)
        9. CapabilityLookup.HasActiveCapabilityForServiceClass (warning only)
       10. Set apiVersion, kind, status.phase="Ready", status.message
       11. Registry.CreateServiceInstance
       12. emitOperation (nil-safe)
       13. writeJSON 201
```

### Component Interaction (Create ServiceBinding flow)

```text
Client
  → POST /v1/service-bindings
    → middleware chain
    → ServiceBindingHandler.Create
        1. safeDecodeServiceBinding (1 MiB, status rejection, DisallowUnknownFields)
        2. validation.ValidateServiceBinding (pure, field-level)
        3. ServiceInstanceLookup.GetServiceInstance (existence)
        4. Set apiVersion, kind, status.phase="Ready", status.secretRef="stub-secret-ref"
        5. Registry.CreateServiceBinding
        6. emitOperation (nil-safe)
        7. writeJSON 201
```

## Files

### New Files

| File | Responsibility |
|------|----------------|
| `internal/resources/serviceinstance.go` | ServiceInstance struct, constants |
| `internal/resources/servicebinding.go` | ServiceBinding struct, constants |
| `internal/registry/serviceinstance_registry.go` | ServiceInstanceRegistryIface, ServiceInstanceLookup, in-memory ServiceInstanceRegistry |
| `internal/registry/servicebinding_registry.go` | ServiceBindingRegistryIface, ServiceBindingInstanceBlocker, in-memory ServiceBindingRegistry |
| `internal/registry/serviceplan_instance_blocker.go` | ServicePlanInstanceBlocker interface + ServiceInstancePlanBlockerChecker impl |
| `internal/registry/project_instance_blocker.go` | ProjectInstanceBlocker interface + ServiceInstanceProjectBlockerChecker impl |
| `internal/registry/capability_lookup.go` | CapabilityLookupImpl (HasActiveCapabilityForServiceClass concrete implementation) |
| `internal/validation/serviceinstance.go` | ValidateServiceInstance, ValidateServiceInstancePathSegment |
| `internal/validation/servicebinding.go` | ValidateServiceBinding, ValidateServiceBindingPathSegment |
| `internal/api/serviceinstance_handler.go` | ServiceInstanceHandler (CRUD) |
| `internal/api/servicebinding_handler.go` | ServiceBindingHandler (Create/Get/List/Delete, PUT→405) |
| `internal/api/serviceinstance_decode.go` | safeDecodeServiceInstance |
| `internal/api/servicebinding_decode.go` | safeDecodeServiceBinding |

### Test Files

| File | Coverage |
|------|----------|
| `internal/validation/serviceinstance_test.go` | Unit tests for ValidateServiceInstance |
| `internal/validation/servicebinding_test.go` | Unit tests for ValidateServiceBinding |
| `internal/validation/serviceinstance_property_test.go` | Property tests |
| `internal/validation/servicebinding_property_test.go` | Property tests |
| `internal/registry/serviceinstance_registry_test.go` | CRUD, filters, CountByServicePlan, CountByProject |
| `internal/registry/servicebinding_registry_test.go` | CRUD, filters, CountByServiceInstance |
| `internal/registry/serviceinstance_registry_race_test.go` | 10+ goroutines race test |
| `internal/registry/servicebinding_registry_race_test.go` | 10+ goroutines race test |
| `internal/registry/serviceinstance_registry_property_test.go` | Property tests |
| `internal/registry/servicebinding_registry_property_test.go` | Property tests |
| `internal/registry/serviceplan_instance_blocker_test.go` | Blocker unit test |
| `internal/registry/project_instance_blocker_test.go` | Blocker unit test |
| `internal/registry/capability_lookup_test.go` | HasActiveCapabilityForServiceClass tests |
| `internal/api/serviceinstance_handler_test.go` | Handler tests (CRUD, validation, references, blocker) |
| `internal/api/servicebinding_handler_test.go` | Handler tests (Create/Get/List/Delete, PUT→405) |
| `internal/api/serviceinstance_emission_test.go` | Operation emission tests |
| `internal/api/servicebinding_emission_test.go` | Operation emission tests |

### Modified Files

| File | Change |
|------|--------|
| `internal/resources/operation.go` | Add `ServiceInstanceName`, `ServiceBindingName` to OperationSpec; add operation type constants |
| `internal/server/server.go` | Register ServiceInstance and ServiceBinding routes; accept new handlers in `New()` |
| `internal/api/serviceplan_handler.go` | Add `instanceBlocker` field; call it in Delete |
| `internal/api/project_handler.go` | Add `instanceBlocker` field; call it in Delete |
| `internal/registry/project_registry.go` | Add `ProjectLookup` interface declaration |
| `internal/registry/serviceplan_registry.go` | Add `ServicePlanLookup` interface declaration |
| `internal/registry/capability_registry.go` | Add `CapabilityLookup` interface declaration |
| `cmd/sovrunn-api/main.go` | Wire new registries, blockers, lookup implementations, and handlers |

## Data Models

### ServiceInstance

```go
// internal/resources/serviceinstance.go
package resources

type ServiceInstance struct {
    APIVersion string                `json:"apiVersion"`
    Kind       string                `json:"kind"`
    Metadata   Metadata              `json:"metadata"`
    Spec       ServiceInstanceSpec   `json:"spec"`
    Status     ServiceInstanceStatus `json:"status"`
}

type ServiceInstanceSpec struct {
    OrganizationRef     string            `json:"organizationRef"`
    OrganizationUnitRef string            `json:"organizationUnitRef,omitempty"`
    TenantRef           string            `json:"tenantRef"`
    ProjectRef          string            `json:"projectRef"`
    ServiceClassRef     string            `json:"serviceClassRef"`
    ServicePlanRef      string            `json:"servicePlanRef"`
    Parameters          map[string]string `json:"parameters,omitempty"`
}

type ServiceInstanceStatus struct {
    Phase   string `json:"phase"`
    Message string `json:"message,omitempty"`
}

const (
    ServiceInstanceAPIVersion = "platform.sovrunn.io/v1alpha1"
    ServiceInstanceKind       = "ServiceInstance"
)
```

### ServiceBinding

```go
// internal/resources/servicebinding.go
package resources

type ServiceBinding struct {
    APIVersion string               `json:"apiVersion"`
    Kind       string               `json:"kind"`
    Metadata   Metadata             `json:"metadata"`
    Spec       ServiceBindingSpec   `json:"spec"`
    Status     ServiceBindingStatus `json:"status"`
}

type ServiceBindingSpec struct {
    ServiceInstanceRef string       `json:"serviceInstanceRef"`
    ConsumerRef        *ConsumerRef `json:"consumerRef"`
    BindingType        string       `json:"bindingType"`
}

type ConsumerRef struct {
    Kind string `json:"kind"`
    Name string `json:"name"`
}

type ServiceBindingStatus struct {
    Phase     string `json:"phase"`
    Message   string `json:"message,omitempty"`
    SecretRef string `json:"secretRef,omitempty"`
}

const (
    ServiceBindingAPIVersion = "platform.sovrunn.io/v1alpha1"
    ServiceBindingKind       = "ServiceBinding"
)

// Allowed binding types (Phase 1).
const (
    BindingTypeCredentials = "credentials"
)
```

### OperationSpec Extension

```go
// Added to internal/resources/operation.go (OperationSpec struct)
ServiceInstanceName string `json:"serviceInstanceName,omitempty"`
ServiceBindingName  string `json:"serviceBindingName,omitempty"`
```

### New Operation Type Constants

```go
// Added to internal/resources/operation.go
const (
    OpCreateServiceInstance = "CreateServiceInstance"
    OpUpdateServiceInstance = "UpdateServiceInstance"
    OpDeleteServiceInstance = "DeleteServiceInstance"
    OpCreateServiceBinding  = "CreateServiceBinding"
    OpDeleteServiceBinding  = "DeleteServiceBinding"
)
```

## Interfaces

### Registry Interfaces

```go
// internal/registry/serviceinstance_registry.go

type ServiceInstanceRegistryIface interface {
    CreateServiceInstance(ctx context.Context, si resources.ServiceInstance) (resources.ServiceInstance, error)
    GetServiceInstance(ctx context.Context, name string) (resources.ServiceInstance, error)
    ListServiceInstances(ctx context.Context, tenantRef, projectRef string) ([]resources.ServiceInstance, error)
    UpdateServiceInstance(ctx context.Context, name string, si resources.ServiceInstance) (resources.ServiceInstance, error)
    DeleteServiceInstance(ctx context.Context, name string) error
    CountByServicePlan(ctx context.Context, serviceClassRef, servicePlanRef string) (int, error)
    CountByProject(ctx context.Context, organizationRef, organizationUnitRef, tenantRef, projectRef string) (int, error)
}

// ServiceInstanceLookup is declared here (colocated with ServiceInstanceRegistryIface
// which satisfies it). Used by ServiceBindingHandler.
type ServiceInstanceLookup interface {
    GetServiceInstance(ctx context.Context, name string) (resources.ServiceInstance, error)
}
```

```go
// internal/registry/servicebinding_registry.go

type ServiceBindingRegistryIface interface {
    CreateServiceBinding(ctx context.Context, sb resources.ServiceBinding) (resources.ServiceBinding, error)
    GetServiceBinding(ctx context.Context, name string) (resources.ServiceBinding, error)
    ListServiceBindings(ctx context.Context, serviceInstanceRef string) ([]resources.ServiceBinding, error)
    DeleteServiceBinding(ctx context.Context, name string) error
    CountByServiceInstance(ctx context.Context, instanceName string) (int, error)
}

// ServiceBindingInstanceBlocker is declared here (colocated with
// ServiceBindingRegistryIface which satisfies it). Used by ServiceInstanceHandler
// to check for referencing bindings before delete.
type ServiceBindingInstanceBlocker interface {
    CountByServiceInstance(ctx context.Context, instanceName string) (int, error)
}
```

```go
// internal/registry/project_registry.go (added by FEATURE-0008)

// ProjectLookup is a narrow interface for verifying parent Project existence.
// The existing *ProjectRegistry satisfies it via GetProject.
type ProjectLookup interface {
    GetProject(ctx context.Context, orgName, ouName, tenantName, name string) (resources.Project, error)
}
```

```go
// internal/registry/serviceplan_registry.go (added by FEATURE-0008)

// ServicePlanLookup is a narrow interface for verifying ServicePlan existence
// and its association with a ServiceClass. The existing *ServicePlanRegistry
// satisfies it via GetServicePlan.
type ServicePlanLookup interface {
    GetServicePlan(ctx context.Context, serviceClassName, name string) (resources.ServicePlan, error)
}
```

```go
// internal/registry/capability_registry.go (added by FEATURE-0008)

// CapabilityLookup is a narrow interface for checking whether an active
// Capability exists for a given ServiceClass. Satisfied by CapabilityLookupImpl
// in capability_lookup.go.
type CapabilityLookup interface {
    HasActiveCapabilityForServiceClass(ctx context.Context, serviceClassRef string) (bool, error)
}
```

### Blocker Interfaces

```go
// internal/registry/serviceplan_instance_blocker.go

type ServicePlanInstanceBlocker interface {
    BlockedByServicePlanInstances(ctx context.Context, serviceClassName, planName string) ([]BlockedBy, error)
}

type ServiceInstancePlanBlockerChecker struct {
    siRegistry ServiceInstanceRegistryIface
}

func (c *ServiceInstancePlanBlockerChecker) BlockedByServicePlanInstances(
    ctx context.Context, serviceClassName, planName string,
) ([]BlockedBy, error) {
    count, err := c.siRegistry.CountByServicePlan(ctx, serviceClassName, planName)
    if err != nil {
        return nil, err
    }
    if count > 0 {
        return []BlockedBy{{Kind: "ServiceInstance", Count: count}}, nil
    }
    return nil, nil
}
```

```go
// internal/registry/project_instance_blocker.go

type ProjectInstanceBlocker interface {
    BlockedByProjectInstances(ctx context.Context, orgName, ouName, tenantName, projectName string) ([]BlockedBy, error)
}

type ServiceInstanceProjectBlockerChecker struct {
    siRegistry ServiceInstanceRegistryIface
}

func (c *ServiceInstanceProjectBlockerChecker) BlockedByProjectInstances(
    ctx context.Context, orgName, ouName, tenantName, projectName string,
) ([]BlockedBy, error) {
    count, err := c.siRegistry.CountByProject(ctx, orgName, ouName, tenantName, projectName)
    if err != nil {
        return nil, err
    }
    if count > 0 {
        return []BlockedBy{{Kind: "ServiceInstance", Count: count}}, nil
    }
    return nil, nil
}
```

### CapabilityLookup Implementation

```go
// internal/registry/capability_lookup.go

type CapabilityLookupImpl struct {
    capRegistry CapabilityRegistryIface
}

func NewCapabilityLookup(reg CapabilityRegistryIface) *CapabilityLookupImpl {
    return &CapabilityLookupImpl{capRegistry: reg}
}

func (l *CapabilityLookupImpl) HasActiveCapabilityForServiceClass(
    ctx context.Context, serviceClassRef string,
) (bool, error) {
    caps, err := l.capRegistry.ListCapabilities(ctx, "", serviceClassRef)
    if err != nil {
        return false, err
    }
    for _, c := range caps {
        if c.Status.Phase == resources.PhaseActive && c.Spec.Supported {
            return true, nil
        }
    }
    return false, nil
}
```

## Validation

### ValidateServiceInstance

Pure function. No I/O, no registry lookups.

```go
func ValidateServiceInstance(si resources.ServiceInstance) []resources.FieldError
```

Checks:
1. `metadata.name` — required, DNS-label, ≤63 chars.
2. `spec.organizationRef` — required, DNS-label, ≤63 chars.
3. `spec.organizationUnitRef` — optional; if non-empty, DNS-label, ≤63 chars.
4. `spec.tenantRef` — required, DNS-label, ≤63 chars.
5. `spec.projectRef` — required, DNS-label, ≤63 chars.
6. `spec.serviceClassRef` — required, DNS-label, ≤63 chars.
7. `spec.servicePlanRef` — required, DNS-label, ≤63 chars.
8. `spec.parameters` — no validation in Phase 1 (type correctness enforced by JSON decoder).

Field names in errors use the JSON path: `"metadata.name"`, `"spec.organizationRef"`, etc.

### ValidateServiceInstancePathSegment

```go
func ValidateServiceInstancePathSegment(name string) []resources.FieldError
```

Validates the single `{name}` path segment as a DNS-label before registry lookup.

### ValidateServiceBinding

Pure function.

```go
func ValidateServiceBinding(sb resources.ServiceBinding) []resources.FieldError
```

Checks:
1. `metadata.name` — required, DNS-label, ≤63 chars.
2. `spec.serviceInstanceRef` — required, DNS-label, ≤63 chars.
3. `spec.consumerRef` — required (non-nil).
4. `spec.consumerRef.kind` — required (non-empty string, no enum restriction).
5. `spec.consumerRef.name` — required, DNS-label, ≤63 chars.
6. `spec.bindingType` — required, must be `"credentials"`.

### ValidateServiceBindingPathSegment

```go
func ValidateServiceBindingPathSegment(name string) []resources.FieldError
```

Validates the single `{name}` path segment as a DNS-label before registry lookup.

## API / Handler Design

### ServiceInstanceHandler

```go
type ServiceInstanceHandler struct {
    registry           registry.ServiceInstanceRegistryIface
    orgLookup          registry.OrganizationLookup
    ouLookup           registry.OrganizationUnitLookup
    tenantLookup       registry.TenantLookup
    projectLookup      registry.ProjectLookup
    serviceClassLookup registry.ServiceClassLookup
    servicePlanLookup  registry.ServicePlanLookup
    capabilityLookup   registry.CapabilityLookup
    bindingBlocker     registry.ServiceBindingInstanceBlocker
    emitter            OperationEmitter
    logger             *log.Logger
}
```

Constructor: `NewServiceInstanceHandler(...)`. All lookup fields are required (non-nil). `emitter` and `logger` may be nil.

**Routes:**
- `POST /v1/service-instances` → Create
- `GET /v1/service-instances` → List (accepts `?tenantRef=...&projectRef=...`)
- `GET /v1/service-instances/{name}` → Get
- `PUT /v1/service-instances/{name}` → Update
- `DELETE /v1/service-instances/{name}` → Delete

**HandleCollection** dispatches POST/GET. **HandleItem** trims prefix, splits on `/`, expects exactly 1 non-empty segment or returns 404.

#### Create Flow

1. `safeDecodeServiceInstance` — 1 MiB, status rejection, DisallowUnknownFields.
2. `validation.ValidateServiceInstance` — pure field validation.
3. Reference validation (in order):
   - `orgLookup.GetOrganization(ctx, spec.organizationRef)` → 400 if not found.
   - If `spec.organizationUnitRef` non-empty: `ouLookup.GetOrganizationUnit(ctx, spec.organizationRef, spec.organizationUnitRef)` → 400 if not found.
   - `tenantLookup.GetTenant(ctx, spec.organizationRef, spec.organizationUnitRef, spec.tenantRef)` → 400 if not found.
   - `projectLookup.GetProject(ctx, spec.organizationRef, spec.organizationUnitRef, spec.tenantRef, spec.projectRef)` → 400 if not found.
   - `serviceClassLookup.GetServiceClass(ctx, spec.serviceClassRef)` → 400 if not found.
   - `servicePlanLookup.GetServicePlan(ctx, spec.serviceClassRef, spec.servicePlanRef)` → 400 if not found (also verifies plan belongs to class via composite key).
4. Capability warning: `capabilityLookup.HasActiveCapabilityForServiceClass(ctx, spec.serviceClassRef)` — if false, log warning; do not block.
5. Set server-owned fields: `apiVersion`, `kind`, `status.phase = "Ready"`, `status.message = "Registered only; no real provisioning in Phase 1"`.
6. `registry.CreateServiceInstance` → 409 if duplicate.
7. `emitOperation` with `OpCreateServiceInstance`.
8. Return 201 with full resource.

#### Update Flow

1. Validate path segment (DNS-label).
2. `safeDecodeServiceInstance` — status rejection.
3. Identity check: body `metadata.name` must be present and match path `{name}`.
4. `validation.ValidateServiceInstance` — pure field validation.
5. Immutability check: get stored entry; compare ALL governance and catalog fields: `spec.organizationRef`, `spec.organizationUnitRef`, `spec.tenantRef`, `spec.projectRef`, `spec.serviceClassRef`, `spec.servicePlanRef`. If any differ → 400 `VALIDATION_FAILED` with the changed field name.
6. `registry.UpdateServiceInstance` — preserves `apiVersion`, `kind`, `status` from stored entry; preserves all immutable spec fields from stored entry; updates mutable fields only: `spec.parameters`, `metadata.labels`, `metadata.annotations`, `metadata.displayName`.
7. `emitOperation` with `OpUpdateServiceInstance`.
8. Return 200 with updated resource.

#### Delete Flow

1. Validate path segment.
2. Binding blocker: `bindingBlocker.CountByServiceInstance(ctx, name)` → 409 `DELETE_BLOCKED` if count > 0.
3. `registry.DeleteServiceInstance` → 404 if not found.
4. `emitOperation` with `OpDeleteServiceInstance`.
5. Return 204.

#### List Flow

1. Read query params `tenantRef`, `projectRef`.
2. `registry.ListServiceInstances(ctx, tenantRef, projectRef)`.
3. Return `{"items": [...]}` (empty array if none).

### ServiceBindingHandler

```go
type ServiceBindingHandler struct {
    registry       registry.ServiceBindingRegistryIface
    instanceLookup registry.ServiceInstanceLookup
    emitter        OperationEmitter
}
```

**Routes:**
- `POST /v1/service-bindings` → Create
- `GET /v1/service-bindings` → List (accepts `?serviceInstanceRef=...`)
- `GET /v1/service-bindings/{name}` → Get
- `PUT /v1/service-bindings/{name}` → 405 Method Not Allowed
- `DELETE /v1/service-bindings/{name}` → Delete

#### Create Flow

1. `safeDecodeServiceBinding` — 1 MiB, status rejection, DisallowUnknownFields.
2. `validation.ValidateServiceBinding` — pure field validation.
3. `instanceLookup.GetServiceInstance(ctx, spec.serviceInstanceRef)` → 400 if not found.
4. Set server-owned fields: `apiVersion`, `kind`, `status.phase = "Ready"`, `status.secretRef = "stub-secret-ref"`.
5. `registry.CreateServiceBinding` → 409 if duplicate.
6. `emitOperation` with `OpCreateServiceBinding` (sets `serviceInstanceName` to the referenced instance name).
7. Return 201 with full resource.

#### Delete Flow

1. Validate path segment.
2. Get binding from registry first (to capture `serviceInstanceRef` for operation emission).
3. `registry.DeleteServiceBinding` → 404 if not found.
4. `emitOperation` with `OpDeleteServiceBinding`.
5. Return 204.

#### PUT → 405

The `HandleItem` switch returns `METHOD_NOT_ALLOWED` with message "ServiceBinding does not support update; delete and recreate instead", regardless of resource existence.

## Registry / Storage Design

### ServiceInstanceRegistry

In-memory, `map[string]resources.ServiceInstance`, keyed by `metadata.name`, protected by `sync.RWMutex`.

- **CreateServiceInstance**: deep-copy input, store. Return deep-copy of stored. Duplicate → `ErrAlreadyExists`.
- **GetServiceInstance**: read lock, return deep-copy. Missing → `ErrNotFound`.
- **ListServiceInstances(ctx, tenantRef, projectRef)**: read lock, iterate, filter if params non-empty (AND logic), sort by `metadata.name`, return slice. Empty → non-nil `[]resources.ServiceInstance{}`.
- **UpdateServiceInstance(name, si)**: write lock, check existence → `ErrNotFound` if missing. Preserve: `apiVersion`, `kind`, `status.phase`, `status.message` from stored. Preserve (immutable): `spec.organizationRef`, `spec.organizationUnitRef`, `spec.tenantRef`, `spec.projectRef`, `spec.serviceClassRef`, `spec.servicePlanRef` from stored. Replace (mutable): `spec.parameters`, `metadata.labels`, `metadata.annotations`, `metadata.displayName`. Deep-copy on store and return.
- **DeleteServiceInstance**: write lock, check existence → `ErrNotFound`. Remove from map.
- **CountByServicePlan(ctx, serviceClassRef, servicePlanRef)**: read lock, iterate, count entries where `spec.serviceClassRef == serviceClassRef` AND `spec.servicePlanRef == servicePlanRef`.
- **CountByProject(ctx, orgRef, ouRef, tenantRef, projectRef)**: read lock, iterate, count entries where all four governance refs match exactly (including empty string for `ouRef` — an empty `ouRef` only matches instances where `spec.organizationUnitRef` is also empty).

Deep copy duplicates `Metadata.Labels`, `Metadata.Annotations`, and `Spec.Parameters` maps.

### ServiceBindingRegistry

In-memory, `map[string]resources.ServiceBinding`, keyed by `metadata.name`, protected by `sync.RWMutex`.

- **CreateServiceBinding**: deep-copy input, store. Return deep-copy. Duplicate → `ErrAlreadyExists`.
- **GetServiceBinding**: read lock, deep-copy. Missing → `ErrNotFound`.
- **ListServiceBindings(ctx, serviceInstanceRef)**: read lock, iterate, filter if param non-empty, sort by `metadata.name`. Empty → non-nil `[]resources.ServiceBinding{}`.
- **DeleteServiceBinding**: write lock, check existence → `ErrNotFound`. Remove from map. Returns `error` only (no body on 204).
- **CountByServiceInstance(ctx, instanceName)**: read lock, count entries where `spec.serviceInstanceRef == instanceName`.

Deep copy duplicates `Metadata.Labels`, `Metadata.Annotations` maps. `ConsumerRef` is a pointer — deep copy creates a new struct.

### ServiceBindingInstanceBlocker

The ServiceInstance delete handler needs to check for referencing ServiceBindings. The `ServiceBindingInstanceBlocker` interface is declared in `internal/registry/servicebinding_registry.go` (colocated with `ServiceBindingRegistryIface`):

```go
// internal/registry/servicebinding_registry.go

type ServiceBindingInstanceBlocker interface {
    CountByServiceInstance(ctx context.Context, instanceName string) (int, error)
}
```

The existing `ServiceBindingRegistry` already satisfies this interface via its `CountByServiceInstance` method. The ServiceInstance handler receives it as a constructor dependency named `bindingBlocker`.

## Operation / Audit Behavior

### Operation Emission Rules

1. Emit only AFTER a successful mutating action (create/update/delete).
2. Failed validation, duplicate, missing reference, not-found, delete-blocked → no Operation emitted.
3. Emission uses `emitOperation(ctx, emitter, spec)` — nil-safe, swallows errors.
4. Operation records include `requestId` from context.

### ServiceInstance Operation Fields

| Action | Type | ResourceKind | ResourceName | ServiceInstanceName | ServiceBindingName |
|--------|------|-------------|--------------|--------------------|--------------------|
| Create | `CreateServiceInstance` | `ServiceInstance` | instance name | instance name | (empty) |
| Update | `UpdateServiceInstance` | `ServiceInstance` | instance name | instance name | (empty) |
| Delete | `DeleteServiceInstance` | `ServiceInstance` | instance name | instance name | (empty) |

Additionally, governance fields (`OrganizationName`, `OrganizationUnitName`, `TenantName`, `ProjectName`) are populated from the ServiceInstance spec.

### ServiceBinding Operation Fields

| Action | Type | ResourceKind | ResourceName | ServiceInstanceName | ServiceBindingName |
|--------|------|-------------|--------------|--------------------|--------------------|
| Create | `CreateServiceBinding` | `ServiceBinding` | binding name | referenced instance name | binding name |
| Delete | `DeleteServiceBinding` | `ServiceBinding` | binding name | referenced instance name | binding name |

## Error Mapping

| Condition | HTTP Status | Error Code | Field |
|-----------|-------------|------------|-------|
| Invalid metadata.name | 400 | VALIDATION_FAILED | metadata.name |
| Invalid spec.organizationRef | 400 | VALIDATION_FAILED | spec.organizationRef |
| Invalid spec.tenantRef | 400 | VALIDATION_FAILED | spec.tenantRef |
| Invalid spec.projectRef | 400 | VALIDATION_FAILED | spec.projectRef |
| Invalid spec.serviceClassRef | 400 | VALIDATION_FAILED | spec.serviceClassRef |
| Invalid spec.servicePlanRef | 400 | VALIDATION_FAILED | spec.servicePlanRef |
| Invalid spec.organizationUnitRef | 400 | VALIDATION_FAILED | spec.organizationUnitRef |
| Invalid spec.serviceInstanceRef | 400 | VALIDATION_FAILED | spec.serviceInstanceRef |
| Invalid spec.consumerRef | 400 | VALIDATION_FAILED | spec.consumerRef |
| Invalid spec.consumerRef.kind | 400 | VALIDATION_FAILED | spec.consumerRef.kind |
| Invalid spec.consumerRef.name | 400 | VALIDATION_FAILED | spec.consumerRef.name |
| Invalid spec.bindingType | 400 | VALIDATION_FAILED | spec.bindingType |
| Status key in body | 400 | VALIDATION_FAILED | status |
| Body too large | 413 | VALIDATION_FAILED | (empty) |
| Malformed JSON | 400 | VALIDATION_FAILED | (empty) |
| Unknown field | 400 | VALIDATION_FAILED | (empty) |
| Empty body | 400 | VALIDATION_FAILED | (empty) |
| Wrong content type | 415 | VALIDATION_FAILED | (empty) |
| Organization not found (ref) | 400 | VALIDATION_FAILED | spec.organizationRef |
| OU not found or mismatch (ref) | 400 | VALIDATION_FAILED | spec.organizationUnitRef |
| Tenant not found or mismatch (ref) | 400 | VALIDATION_FAILED | spec.tenantRef |
| Project not found or mismatch (ref) | 400 | VALIDATION_FAILED | spec.projectRef |
| ServiceClass not found (ref) | 400 | VALIDATION_FAILED | spec.serviceClassRef |
| ServicePlan not found or mismatch (ref) | 400 | VALIDATION_FAILED | spec.servicePlanRef |
| ServiceInstance not found (ref for binding) | 400 | VALIDATION_FAILED | spec.serviceInstanceRef |
| Immutable field changed on update | 400 | VALIDATION_FAILED | (changed field: spec.organizationRef, spec.organizationUnitRef, spec.tenantRef, spec.projectRef, spec.serviceClassRef, or spec.servicePlanRef) |
| Path/body name mismatch | 400 | VALIDATION_FAILED | metadata.name |
| Duplicate name (create) | 409 | RESOURCE_ALREADY_EXISTS | (empty) |
| Resource not found (get/update/delete) | 404 | RESOURCE_NOT_FOUND | (empty) |
| Delete blocked (ServiceBindings exist) | 409 | DELETE_BLOCKED | (empty) |
| Delete blocked (ServiceInstances exist for Plan) | 409 | DELETE_BLOCKED | (empty) |
| Delete blocked (ServiceInstances exist for Project) | 409 | DELETE_BLOCKED | (empty) |
| PUT on ServiceBinding | 405 | METHOD_NOT_ALLOWED | (empty) |
| Wrong path segment count | 404 | RESOURCE_NOT_FOUND | (empty) |
| Internal error | 500 | INTERNAL_ERROR | (empty) |

## Security and Privacy

1. **No secrets in parameters**: `spec.parameters` is a plain `map[string]string`. Documentation warns against storing secrets. No enforcement in Phase 1 (the ServicePlan validator's forbidden-key check does NOT apply to ServiceInstance parameters in Phase 1, as those parameters are user-supplied instance configuration, not catalog-level definitions).
2. **Stub secretRef**: `"stub-secret-ref"` is a static placeholder. No credential generation occurs.
3. **No raw body storage**: decoded into typed structs only; raw bytes discarded after decode.
4. **No raw body echo**: error responses never include the original request body.
5. **No parameter logging**: handler logs include resource name, kind, and reference names only. Parameter map values are never logged.
6. **No internal detail exposure**: error messages use generic safe text (e.g., "internal error") for 500s.
7. **Body size limit**: 1 MiB via `http.MaxBytesReader`.
8. **Content-type enforcement**: `contentTypeMiddleware` rejects non-`application/json` for mutating endpoints.

## Testing Strategy

### Validation Unit Tests (serviceinstance_test.go, servicebinding_test.go)

- Valid ServiceInstance/ServiceBinding accepted (no errors).
- Missing/empty `metadata.name` → error.
- Invalid DNS-label names (uppercase, special chars, leading/trailing hyphen, >63 chars) → error.
- Missing required spec fields → error on each.
- Optional `organizationUnitRef` empty → accepted.
- Optional `organizationUnitRef` invalid → error.
- Optional `parameters` nil or empty map → accepted.
- `bindingType` valid (`credentials`) → accepted.
- `bindingType` invalid → error.
- `consumerRef` nil → error.
- `consumerRef.kind` empty → error.
- `consumerRef.name` invalid → error.
- `consumerRef.kind` any non-empty string → accepted.

### Validation Property Tests (serviceinstance_property_test.go, servicebinding_property_test.go)

- Valid DNS-label names accepted (testing/quick, MaxCount: 100).
- Arbitrary invalid strings rejected.
- Valid bindingType accepted; invalid rejected.

### Registry Unit Tests

**ServiceInstance:**
- Create stores; Get returns deep copy.
- Duplicate → `ErrAlreadyExists`; original unchanged.
- Get missing → `ErrNotFound`.
- List sorted by name; empty → non-nil `[]`.
- List with `tenantRef` filter.
- List with `projectRef` filter.
- List with both filters (AND).
- List with no filters (all returned).
- Update mutable fields only (parameters, labels, annotations, displayName).
- Update preserves stored status unchanged.
- Update missing → `ErrNotFound`.
- Delete removes entry.
- Delete missing → `ErrNotFound`.
- `CountByServicePlan` correct count (matching both serviceClassRef and servicePlanRef).
- `CountByServicePlan` no false positives when same planName exists under different ServiceClasses.
- `CountByProject` correct count (matching all four hierarchy refs).
- `CountByProject` no false positives when same projectName exists under different tenants.
- `CountByProject` with empty `organizationUnitRef` does NOT match instances with non-empty OU.
- `CountByProject` with non-empty `organizationUnitRef` does NOT match instances with empty OU.
- Duplicate name across different governance refs → `ErrAlreadyExists` (global name uniqueness).
- Deep-copy immutability: mutating returned value does not affect registry.

**ServiceBinding:**
- Create stores; Get returns deep copy.
- Duplicate → `ErrAlreadyExists`.
- Duplicate name referencing different ServiceInstance → `ErrAlreadyExists` (global name uniqueness).
- Get missing → `ErrNotFound`.
- List sorted; empty → non-nil `[]`.
- List with `serviceInstanceRef` filter.
- List with no filter (all returned).
- Delete removes.
- Delete missing → `ErrNotFound`.
- `CountByServiceInstance` correct count.
- Deep-copy immutability.

### Registry Race Tests

10+ goroutines performing concurrent Create/Get/List/Update/Delete/Count operations. No race reports under `go test -race`.

### Registry Property Tests

- Create/Get round-trip preserves data.
- List sort invariant.
- Deep-copy immutability.
- Duplicate-create idempotent error.
- Filter correctness.
- CountByServicePlan correctness (no false positives).
- CountByProject correctness (no false positives).
- CountByServiceInstance correctness.

### Handler Tests

**ServiceInstance:**
- POST 201 with valid payload; response has apiVersion, kind, status.phase="Ready".
- POST 409 duplicate name.
- POST 400: invalid fields, status key present, bad JSON, unknown field, oversized body.
- POST 400: missing Organization ref, Tenant ref, Project ref, ServiceClass ref, ServicePlan ref.
- POST 400: ServicePlan not matching ServiceClass.
- POST 400: governance hierarchy inconsistency.
- GET 200 found; 404 missing; 400 invalid path segment.
- GET list sorted; empty `{"items": []}`.
- GET list with query params filtering.
- PUT 200 success.
- PUT 404 missing resource.
- PUT 400 path/body name mismatch.
- PUT 400 status key present → `VALIDATION_FAILED` field=status.
- PUT 400 immutable field change (organizationRef, organizationUnitRef, tenantRef, projectRef, serviceClassRef, servicePlanRef) — one test per field.
- PUT 200 preserves stored status unchanged.
- PUT 200 updates mutable fields (parameters, labels, annotations, displayName).
- DELETE 204 success.
- DELETE 404 missing.
- DELETE 409 blocked by ServiceBindings.
- POST 409 duplicate name across different governance refs (global uniqueness).
- Wrong path shape (extra segments) → 404.

**ServiceBinding:**
- POST 201 valid; response has secretRef="stub-secret-ref".
- POST 409 duplicate.
- POST 409 duplicate name referencing different ServiceInstance (global uniqueness).
- POST 400 validation errors, status key, bad JSON.
- POST 400 missing ServiceInstance ref.
- GET 200/404/400.
- GET list with serviceInstanceRef filter.
- PUT → 405 METHOD_NOT_ALLOWED.
- DELETE 204/404.
- Wrong path shape → 404.

### Delete-Blocking Tests

- ServiceInstance delete with ServiceBindings → 409.
- ServiceInstance delete with zero bindings → 204.
- ServicePlan delete with referencing ServiceInstances → 409.
- ServicePlan delete with zero → normal delete.
- Project delete with ServiceInstances → 409.
- Project delete with zero → normal delete.

### Operation Emission Tests

- Successful create/update/delete of ServiceInstance emits correct Operation type.
- Successful create/delete of ServiceBinding emits correct Operation type.
- ServiceBinding emission includes `serviceInstanceName`.
- Failed actions (validation, duplicate, missing ref, not-found, blocked) → no emission.
- Nil emitter → no panic, no emission.

### Capability Warning Test

- ServiceInstance created without matching active Capability → warning logged, creation succeeds.
- ServiceInstance created with matching active Capability → no warning, creation succeeds.

## Verification Commands

```bash
make fmt
make test
make vet
go test -race ./...
```

## Non-Goals

Explicitly NOT implemented in this feature:

1. Real infrastructure provisioning.
2. Actual credential generation or secret backend integration.
3. Plugin execution or ServiceOps orchestration.
4. Async provisioning workflows.
5. Kubernetes CRDs or operators.
6. Authentication, authorization, or RBAC.
7. Quota or resource limit enforcement.
8. Billing, metering, or cost tracking.
9. UI portal.
10. Persistent database storage.
11. Go 1.22 wildcard routing or new external dependencies.
12. SDE-specific behavior.
13. Multi-cluster federation.
14. Service mesh integration.
15. Additional filter query parameters beyond `tenantRef`/`projectRef`/`serviceInstanceRef`.

## Resolved Design Questions Summary

| # | Question | Resolution |
|---|----------|------------|
| DQ1 | Governance hierarchy validation logic | Step-by-step composite-key lookups (org→OU→tenant→project); existing interfaces return full resource |
| DQ2 | Lookup interface placement | Each `{Resource}Lookup` in the resource's own registry file (colocated with provider); new: ProjectLookup in project_registry.go, ServicePlanLookup in serviceplan_registry.go, ServiceInstanceLookup in serviceinstance_registry.go, CapabilityLookup in capability_registry.go |
| DQ3 | ServicePlan delete-blocker extension | `ServicePlanInstanceBlocker` interface + impl in `serviceplan_instance_blocker.go`; injected into ServicePlanHandler |
| DQ4 | Project delete-blocker extension | `ProjectInstanceBlocker` interface + impl in `project_instance_blocker.go`; injected into ProjectHandler |
| DQ5 | ServiceBinding immutability | Confirmed: no PUT; 405 returned; delete-and-recreate is the pattern |
| DQ6 | Capability warning interface | `CapabilityLookup` declared in `capability_registry.go`; impl in `capability_lookup.go`; active = phase Active AND supported=true |
| DQ7 | Stub secretRef value | `"stub-secret-ref"` confirmed |
| DQ8 | Filter query parameters | ServiceInstance: `tenantRef`, `projectRef` (AND); ServiceBinding: `serviceInstanceRef`; no additional filters in Phase 1 |
| DQ9 | ServiceInstance governance immutability | All governance fields (organizationRef, organizationUnitRef, tenantRef, projectRef, serviceClassRef, servicePlanRef) are immutable after creation in Phase 1 |
| DQ10 | Empty organizationUnitRef semantics | Empty string is valid; skips OU lookup; passes through as empty component in composite keys; exact-match semantics in CountByProject |
| DQ11 | Name uniqueness scope | Global across the API server for both ServiceInstance and ServiceBinding; duplicate name rejected regardless of governance refs |
