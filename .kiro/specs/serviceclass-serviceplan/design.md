# Design Document — FEATURE-0006 ServiceClass and ServicePlan

## Overview

FEATURE-0006 implements `ServiceClass` and `ServicePlan` as **global platform catalog resources**.
A ServiceClass defines a service type (PostgreSQL, Redis, …); a ServicePlan defines a tier/shape
under a ServiceClass. These are catalog definitions only — they do NOT provision, bind, or execute.

**Global scope:** neither resource is scoped to Organization/OU/Tenant/Project. ServiceClass
identity is `metadata.name`; ServicePlan identity is the composite `serviceClassName/name`.

**Resolved design decisions:**

1. **Parent lookup interface** — A narrow `ServiceClassLookup` interface is defined in
   `internal/registry` (mirroring `OrganizationUnitLookup`/`TenantLookup`). `ServicePlanHandler`
   depends on the interface; `ServicePlanRegistry` does NOT depend on `ServiceClassRegistry`. The
   existing `*ServiceClassRegistry` satisfies the interface.

2. **Delete-blocker interface** — `ServiceClassChildBlocker` with
   `BlockedByServiceClassChildren(ctx, serviceClassName) ([]registry.BlockedBy, error)`, implemented
   by a ServicePlan-backed checker using `CountByServiceClass`. No generic blocker framework.

3. **Parameters value type** — `spec.parameters` is `map[string]string` in Phase 1. This keeps
   parameters simple, deterministic, and less likely to carry nested secret-bearing structures.

**Scope boundary:** No ServiceInstance, ServiceBinding, Plugin/Capability registry, ServiceOps,
provisioning, persistence, auth, UI, AI, or marketplace workflow.

**Architectural constraint:** `internal/api` MUST NOT import `internal/server`; request ID is read
via the existing API-local helper.

---

## Architecture

```
cmd/sovrunn-api → internal/server → internal/api → internal/registry
                                          │              │
                                          │              ├─ ServiceClassRegistry (storage-only)
                                          │              ├─ ServicePlanRegistry (storage-only)
                                          │              └─ ServicePlanChildBlockerChecker
                                          └─ ServiceClassHandler / ServicePlanHandler
                                             (emit Operations via FEATURE-0005 emitter)
                        internal/resources (ServiceClass, ServicePlan models + constants)
```

- ServiceClass/ServicePlan handlers reuse: safe decoding, `writeError`/`writeJSON`, the nil-safe
  `OperationEmitter`, and the API-local `requestIDFromContext`.
- ServicePlan handler holds a `ServiceClassLookup` for parent verification.
- ServiceClass handler holds a `ServiceClassChildBlocker` for delete blocking.

---

## Files to Create/Modify

### New Files

| File | Purpose |
|---|---|
| `internal/resources/serviceclass.go` | ServiceClass, spec/status, category/lifecycle constants, kind |
| `internal/resources/serviceplan.go` | ServicePlan, spec/status, tier constants, kind |
| `internal/registry/serviceclass_registry.go` | ServiceClassRegistryIface, ServiceClassRegistry, ServiceClassLookup |
| `internal/registry/serviceplan_registry.go` | ServicePlanRegistryIface, ServicePlanRegistry, CountByServiceClass |
| `internal/registry/serviceclass_blocker.go` | ServiceClassChildBlocker iface, ServicePlanChildBlockerChecker |
| `internal/validation/serviceclass.go` | ValidateServiceClass, ValidateServiceClassPathSegment |
| `internal/validation/serviceplan.go` | ValidateServicePlan, ValidateServicePlanPathSegments |
| `internal/api/serviceclass_decode.go` | safeDecodeServiceClass |
| `internal/api/serviceplan_decode.go` | safeDecodeServicePlan |
| `internal/api/serviceclass_handler.go` | ServiceClassHandler |
| `internal/api/serviceplan_handler.go` | ServicePlanHandler |
| plus `_test.go` files for each package | unit/property/race/handler tests |

### Modified Files

| File | Change |
|---|---|
| `internal/resources/operation.go` | Add 6 op type constants + `ServiceClassName`/`ServicePlanName` OperationSpec fields + `ServiceClassKind`/`ServicePlanKind` |
| `internal/server/server.go` | Accept `*api.ServiceClassHandler`, `*api.ServicePlanHandler`; register routes |
| `internal/server/server_test.go` | Update constructor tests; add route tests |
| `cmd/sovrunn-api/main.go` | Wire registries, blocker, lookup, handlers, emitter |

---

## Data Models

### internal/resources/serviceclass.go

```go
package resources

type ServiceClass struct {
    APIVersion string             `json:"apiVersion"`
    Kind       string             `json:"kind"`
    Metadata   Metadata           `json:"metadata"`
    Spec       ServiceClassSpec   `json:"spec"`
    Status     ServiceClassStatus `json:"status"`
}

type ServiceClassSpec struct {
    DisplayName     string   `json:"displayName,omitempty"`
    Description     string   `json:"description,omitempty"`
    Category        string   `json:"category"`
    Provider        string   `json:"provider,omitempty"`
    Lifecycle       string   `json:"lifecycle"`
    DefaultPlanName string   `json:"defaultPlanName,omitempty"`
    Tags            []string `json:"tags,omitempty"`
}

type ServiceClassStatus struct {
    Phase   string `json:"phase"`
    Message string `json:"message,omitempty"`
}

const (
    ServiceClassKind = "ServiceClass"
)

// ServiceClass category constants.
const (
    CategoryDatabase      = "Database"
    CategoryCache         = "Cache"
    CategoryObjectStorage = "ObjectStorage"
    CategoryStream        = "Stream"
    CategoryGateway       = "Gateway"
    CategoryFunction      = "Function"
    CategoryAnalytics     = "Analytics"
    CategoryOther         = "Other"
)

// Lifecycle constants (shared by ServiceClass and ServicePlan).
const (
    LifecyclePreview    = "Preview"
    LifecycleActive     = "Active"
    LifecycleDeprecated = "Deprecated"
    LifecycleRetired    = "Retired"
)
```

### internal/resources/serviceplan.go

```go
package resources

type ServicePlan struct {
    APIVersion string            `json:"apiVersion"`
    Kind       string            `json:"kind"`
    Metadata   Metadata          `json:"metadata"`
    Spec       ServicePlanSpec   `json:"spec"`
    Status     ServicePlanStatus `json:"status"`
}

type ServicePlanSpec struct {
    ServiceClassName string            `json:"serviceClassName"`
    DisplayName      string            `json:"displayName,omitempty"`
    Description      string            `json:"description,omitempty"`
    Tier             string            `json:"tier"`
    Lifecycle        string            `json:"lifecycle"`
    Parameters       map[string]string `json:"parameters,omitempty"`
    Tags             []string          `json:"tags,omitempty"`
}

type ServicePlanStatus struct {
    Phase   string `json:"phase"`
    Message string `json:"message,omitempty"`
}

const (
    ServicePlanKind = "ServicePlan"
)

// ServicePlan tier constants.
const (
    TierDev        = "Dev"
    TierSmall      = "Small"
    TierMedium     = "Medium"
    TierLarge      = "Large"
    TierProduction = "Production"
    TierCustom     = "Custom"
)
```

### internal/resources/operation.go — additions

```go
// New operation type constants (FEATURE-0006).
const (
    OpCreateServiceClass = "CreateServiceClass"
    OpUpdateServiceClass = "UpdateServiceClass"
    OpDeleteServiceClass = "DeleteServiceClass"
    OpCreateServicePlan  = "CreateServicePlan"
    OpUpdateServicePlan  = "UpdateServicePlan"
    OpDeleteServicePlan  = "DeleteServicePlan"
)

// OperationSpec gains two optional catalog-reference fields:
//   ServiceClassName string `json:"serviceClassName,omitempty"`
//   ServicePlanName  string `json:"servicePlanName,omitempty"`
// (ServiceClassKind and ServicePlanKind are defined in their resource files.)
```

Reuses the existing `Metadata` struct and phase constants (`PhaseActive`, etc.).

---

## Components and Interfaces

| Component | Package | Responsibility |
|---|---|---|
| `ServiceClass` / `ServicePlan` models | `internal/resources` | Data shape, category/tier/lifecycle/kind constants |
| `ServiceClassRegistryIface` / `ServiceClassRegistry` | `internal/registry` | Storage-only in-memory ServiceClass store |
| `ServicePlanRegistryIface` / `ServicePlanRegistry` | `internal/registry` | Storage-only in-memory ServicePlan store (composite key) |
| `ServiceClassLookup` interface | `internal/registry` | Narrow parent-existence lookup for ServicePlan handler |
| `ServiceClassChildBlocker` interface | `internal/registry` | Narrow delete-blocker boundary for ServiceClass delete |
| `ServicePlanChildBlockerChecker` | `internal/registry` | Implements the blocker via `CountByServiceClass` |
| `ValidateServiceClass` / `ValidateServicePlan` | `internal/validation` | Pure, context-free field validation |
| `safeDecodeServiceClass` / `safeDecodeServicePlan` | `internal/api` | Safe JSON decoding (1 MiB, DisallowUnknownFields, status rejection) |
| `ServiceClassHandler` / `ServicePlanHandler` | `internal/api` | CRUD handlers; emit Operations via FEATURE-0005 emitter |
| `OperationEmitter` (reused) | `internal/api` | Emission boundary consumed by both handlers (nil-safe) |

- Both handlers depend on registry interfaces, never on concrete registries directly for cross-resource access.
- `ServicePlanHandler` holds a `ServiceClassLookup` for parent verification.
- `ServiceClassHandler` holds a nil-safe `ServiceClassChildBlocker` for delete blocking.
- `internal/api` MUST NOT import `internal/server`; request ID is read via the existing API-local `requestIDFromContext`.

Detailed interface and struct definitions follow in the per-component sections below.

---

## Validation Design

### internal/validation/serviceclass.go

```go
// ValidateServiceClass validates all user-authored fields. Context-free, no I/O.
func ValidateServiceClass(sc resources.ServiceClass) []resources.FieldError

// ValidateServiceClassPathSegment validates the single URL path segment.
func ValidateServiceClassPathSegment(name string) []resources.FieldError
```

**Rules:**
- `metadata.name` DNS-label (1–63) → `error.field = "metadata.name"`.
- `spec.category` required and ∈ {Database, Cache, ObjectStorage, Stream, Gateway, Function, Analytics, Other} → `error.field = "spec.category"`.
- `spec.lifecycle` required and ∈ {Preview, Active, Deprecated, Retired} → `error.field = "spec.lifecycle"`.
- IF `spec.defaultPlanName` is non-empty, it MUST be a valid DNS-label → `error.field = "spec.defaultPlanName"` (existence NOT verified).
- `spec.displayName`, `spec.description`, `spec.provider`, `spec.tags` optional; not format-validated.

### internal/validation/serviceplan.go

```go
// ValidateServicePlan validates all user-authored fields. Context-free, no I/O.
func ValidateServicePlan(sp resources.ServicePlan) []resources.FieldError

// ValidateServicePlanPathSegments validates the two URL path segments.
func ValidateServicePlanPathSegments(serviceClassName, name string) []resources.FieldError
```

**Rules:**
- `metadata.name` DNS-label (1–63) → `error.field = "metadata.name"`.
- `spec.serviceClassName` required, DNS-label (1–63) → `error.field = "spec.serviceClassName"`.
- `spec.tier` required and ∈ {Dev, Small, Medium, Large, Production, Custom} → `error.field = "spec.tier"`.
- `spec.lifecycle` required and ∈ {Preview, Active, Deprecated, Retired} → `error.field = "spec.lifecycle"`.
- `spec.parameters` forbidden-key check (see below) → `error.field = "spec.parameters"`.

**Path validation field mapping:**
- invalid `serviceClassName` segment → `spec.serviceClassName`
- invalid `name` segment → `metadata.name`

Both validators reuse `dnsLabelRe`/`validateName` from the validation package (same package, unexported).
Enum membership checks use small package-level sets or switch statements over the resource constants.

### Forbidden Parameter Key Logic (Requirement 4.6)

```go
// forbiddenParamSubstrings are matched case-insensitively against each
// parameter KEY. The plain substring "key" is intentionally NOT listed;
// only the composite secret-bearing phrases below trigger rejection.
var forbiddenParamSubstrings = []string{
    "password", "secret", "token", "credential", "auth",
    "apikey", "accesskey", "secretkey", "privatekey",
}

// isForbiddenParamKey lowercases the key once, then checks containment.
func isForbiddenParamKey(key string) bool {
    lk := strings.ToLower(key)
    for _, s := range forbiddenParamSubstrings {
        if strings.Contains(lk, s) {
            return true
        }
    }
    return false
}
```

- Matching is case-insensitive (`strings.ToLower` on the key, lowercase substrings).
- `apiKey`, `accessKey`, `secretKey`, `privateKey` are rejected because their lowercased forms
  (`apikey`, `accesskey`, `secretkey`, `privatekey`) are listed.
- Benign keys such as `regionKey`, `masterKeyCount`, or plain `key` are NOT rejected — the bare
  substring `key` is not in the forbidden list.
- Any offending key yields a single `FieldError` with `Field = "spec.parameters"`.

---

## Safe JSON Decoding

### internal/api/serviceclass_decode.go / serviceplan_decode.go

```go
func safeDecodeServiceClass(w http.ResponseWriter, r *http.Request) (resources.ServiceClass, error)
func safeDecodeServicePlan(w http.ResponseWriter, r *http.Request) (resources.ServicePlan, error)
```

**Sequence (identical to the FEATURE-0003/0004 decoder pattern):**
1. `r.Body = http.MaxBytesReader(w, r.Body, 1<<20)` — inside the function.
2. Read body; `*http.MaxBytesError` → errBodyTooLarge (413).
3. Empty body → errEmptyBody (400).
4. Decode into `map[string]json.RawMessage`; if `status` key present → errStatusFieldPresent (400).
5. Typed decode with `DisallowUnknownFields()` into the target struct.
6. Unknown field → errUnknownField (400); syntax/type → errMalformedJSON (400).
7. Return the decoded resource.

Do not echo raw body. HTTP 415 is handled by `contentTypeMiddleware`, not here. Reuses the existing
error sentinels from the shared decoder helpers.

---

## ServiceClassRegistry Design

### internal/registry/serviceclass_registry.go

```go
type ServiceClassRegistryIface interface {
    CreateServiceClass(ctx context.Context, sc resources.ServiceClass) (resources.ServiceClass, error)
    GetServiceClass(ctx context.Context, name string) (resources.ServiceClass, error)
    ListServiceClasses(ctx context.Context) ([]resources.ServiceClass, error)
    UpdateServiceClass(ctx context.Context, sc resources.ServiceClass) (resources.ServiceClass, error)
    DeleteServiceClass(ctx context.Context, name string) error
}

type ServiceClassRegistry struct {
    mu    sync.RWMutex
    store map[string]resources.ServiceClass // key: metadata.name
}

func NewServiceClassRegistry() *ServiceClassRegistry {
    return &ServiceClassRegistry{store: make(map[string]resources.ServiceClass)}
}

// ServiceClassLookup is the narrow parent-existence interface consumed by
// ServicePlanHandler. The concrete *ServiceClassRegistry satisfies it.
type ServiceClassLookup interface {
    GetServiceClass(ctx context.Context, name string) (resources.ServiceClass, error)
}
```

### Method Behavior

| Method | Lock | Returns | Error |
|---|---|---|---|
| CreateServiceClass | Lock | deep copy of stored | ErrAlreadyExists if name exists |
| GetServiceClass | RLock | deep copy | ErrNotFound if absent |
| ListServiceClasses | RLock | sorted slice of deep copies | nil |
| UpdateServiceClass | Lock | deep copy of updated | ErrNotFound if absent |
| DeleteServiceClass | Lock | — | ErrNotFound if absent |

**List sort order:** `Metadata.Name` ascending.

**deepCopyServiceClass:** copies the `Tags` slice and `Metadata.Labels`/`Metadata.Annotations`
maps so callers cannot mutate stored state.

**UpdateServiceClass behavior:** derives the key from `sc.Metadata.Name`, looks up the stored
entry, preserves stored `APIVersion`, `Kind`, `Status`, `Metadata.Name`, and replaces only mutable
fields (`Metadata.Labels`, `Metadata.Annotations`, `spec.displayName`, `spec.description`,
`spec.category`, `spec.provider`, `spec.lifecycle`, `spec.defaultPlanName`, `spec.tags`). Returns a
deep copy of the updated stored ServiceClass.

Reuses sentinel errors `ErrNotFound` and `ErrAlreadyExists` from `internal/registry/registry.go`.
No dependency on other registries; no package-level global state.

---

## ServicePlanRegistry Design

### internal/registry/serviceplan_registry.go

```go
type ServicePlanRegistryIface interface {
    CreateServicePlan(ctx context.Context, sp resources.ServicePlan) (resources.ServicePlan, error)
    GetServicePlan(ctx context.Context, serviceClassName, name string) (resources.ServicePlan, error)
    ListServicePlans(ctx context.Context) ([]resources.ServicePlan, error)
    UpdateServicePlan(ctx context.Context, sp resources.ServicePlan) (resources.ServicePlan, error)
    DeleteServicePlan(ctx context.Context, serviceClassName, name string) error
    CountByServiceClass(ctx context.Context, serviceClassName string) (int, error)
}

type ServicePlanRegistry struct {
    mu    sync.RWMutex
    store map[string]resources.ServicePlan // key: "serviceClassName/name"
}

func NewServicePlanRegistry() *ServicePlanRegistry {
    return &ServicePlanRegistry{store: make(map[string]resources.ServicePlan)}
}

func servicePlanCompositeKey(serviceClassName, name string) string {
    return serviceClassName + "/" + name
}
```

### Method Behavior

| Method | Lock | Returns | Error |
|---|---|---|---|
| CreateServicePlan | Lock | deep copy of stored | ErrAlreadyExists if composite key exists |
| GetServicePlan | RLock | deep copy | ErrNotFound if absent |
| ListServicePlans | RLock | sorted slice of deep copies | nil |
| UpdateServicePlan | Lock | deep copy of updated | ErrNotFound if absent |
| DeleteServicePlan | Lock | — | ErrNotFound if absent |
| CountByServiceClass | RLock | int count | nil |

**List sort order:** `Spec.ServiceClassName` ascending, then `Metadata.Name` ascending.

**Composite identity:** the same `metadata.name` under two different `spec.serviceClassName`
values yields two distinct keys and both are stored without conflict.

**deepCopyServicePlan:** copies the `Parameters` map and `Tags` slice plus
`Metadata.Labels`/`Metadata.Annotations` maps so callers cannot mutate stored state.

**UpdateServicePlan behavior:** derives the composite key from
`sp.Spec.ServiceClassName/sp.Metadata.Name`, looks up the stored entry, preserves stored
`APIVersion`, `Kind`, `Status`, `Metadata.Name`, `Spec.ServiceClassName`, and replaces only mutable
fields (`Metadata.Labels`, `Metadata.Annotations`, `spec.displayName`, `spec.description`,
`spec.tier`, `spec.lifecycle`, `spec.parameters`, `spec.tags`). Never moves a plan between classes.
Returns a deep copy of the updated stored ServicePlan.

**CountByServiceClass:** RLock; iterates entries; counts those where `Spec.ServiceClassName`
matches. Used only by the delete blocker.

Reuses sentinel errors `ErrNotFound` and `ErrAlreadyExists`. No dependency on other registries.

---

## ServiceClass Delete Blocker Design

### internal/registry/serviceclass_blocker.go

```go
type ServiceClassChildBlocker interface {
    BlockedByServiceClassChildren(ctx context.Context, serviceClassName string) ([]BlockedBy, error)
}

type ServicePlanChildBlockerChecker struct {
    servicePlanRegistry ServicePlanRegistryIface
}

func NewServicePlanChildBlockerChecker(reg ServicePlanRegistryIface) *ServicePlanChildBlockerChecker

func (c *ServicePlanChildBlockerChecker) BlockedByServiceClassChildren(
    ctx context.Context, serviceClassName string,
) ([]BlockedBy, error) {
    count, err := c.servicePlanRegistry.CountByServiceClass(ctx, serviceClassName)
    if err != nil { return nil, err }
    if count > 0 {
        return []BlockedBy{{Kind: "ServicePlan", Count: count}}, nil
    }
    return nil, nil
}
```

Reuses the existing `BlockedBy` type. No generic blocker framework. Lifecycle state (including
`Retired`) does not exempt a ServicePlan from blocking its parent's deletion.

---

## ServiceClass HTTP Handler

### internal/api/serviceclass_handler.go

```go
type ServiceClassHandler struct {
    registry registry.ServiceClassRegistryIface
    blocker  registry.ServiceClassChildBlocker // nil-safe (delete blocking)
    emitter  OperationEmitter                  // nil-safe (FEATURE-0005)
}

func NewServiceClassHandler(
    reg registry.ServiceClassRegistryIface,
    blocker registry.ServiceClassChildBlocker,
    emitter OperationEmitter,
) *ServiceClassHandler

func (h *ServiceClassHandler) HandleCollection(w http.ResponseWriter, r *http.Request) // POST/GET
func (h *ServiceClassHandler) HandleItem(w http.ResponseWriter, r *http.Request)       // GET/PUT/DELETE
```

### HandleItem path parsing (Go 1.21, exactly one segment)

```
remainder := strings.TrimPrefix(r.URL.Path, "/v1/service-classes/")
parts := strings.Split(remainder, "/")
if len(parts) != 1 || parts[0] == "" → 404
name := parts[0]
// dispatch GET → Get, PUT → Update, DELETE → Delete, else 405
```

### Create Flow

```
1. safeDecodeServiceClass (MaxBytesReader, status detection, DisallowUnknownFields)
2. decode error → writeError(400/413)
3. ValidateServiceClass(sc) → []FieldError; if errors → writeValidationErrors
4. Force: apiVersion, kind = ServiceClass, status.phase = Active
5. created, err := registry.CreateServiceClass(ctx, sc)
6. ErrAlreadyExists → writeError(409, RESOURCE_ALREADY_EXISTS)
7. emitOperation(ctx, h.emitter, OperationSpec{Type: OpCreateServiceClass,
     ResourceKind: ServiceClassKind, ResourceName: created.Metadata.Name,
     ServiceClassName: created.Metadata.Name, RequestID: requestIDFromContext(ctx)})
8. writeJSON(201, created)
```

### Get / List Flow

```
Get:  ValidateServiceClassPathSegment(name); GetServiceClass; ErrNotFound → 404; 200
List: ListServiceClasses; error → 500; writeJSON(200, {"items": items})  // [] when empty
```

### Update Flow

```
1. ValidateServiceClassPathSegment(name); if errors → writeValidationErrors
2. safeDecodeServiceClass(w, r); decode error → writeError(400/413)
3. body.Metadata.Name present and == name, else 400 (field="metadata.name")
4. ValidateServiceClass(sc); if errors → writeValidationErrors
5. updated, err := registry.UpdateServiceClass(ctx, sc)
6. ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
7. emitOperation(ctx, h.emitter, OperationSpec{Type: OpUpdateServiceClass,
     ResourceKind: ServiceClassKind, ResourceName: updated.Metadata.Name,
     ServiceClassName: updated.Metadata.Name, RequestID: requestIDFromContext(ctx)})
8. writeJSON(200, updated)
```

### Delete Flow

```
1. ValidateServiceClassPathSegment(name); if errors → writeValidationErrors
2. if h.blocker != nil:
     blockers, err := h.blocker.BlockedByServiceClassChildren(ctx, name)
     err → writeError(500, INTERNAL_ERROR)
     blockers non-empty → writeError(409, DELETE_BLOCKED, message "ServicePlan")
3. err := registry.DeleteServiceClass(ctx, name)
4. ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
5. emitOperation(ctx, h.emitter, OperationSpec{Type: OpDeleteServiceClass,
     ResourceKind: ServiceClassKind, ResourceName: name,
     ServiceClassName: name, RequestID: requestIDFromContext(ctx)})
6. w.WriteHeader(204)
```

**Nil-safe blocker:** if `blocker` is nil, delete proceeds without child checks (preserves isolated
handler test compatibility). Production wiring injects `ServicePlanChildBlockerChecker`.

---

## ServicePlan HTTP Handler

### internal/api/serviceplan_handler.go

```go
type ServicePlanHandler struct {
    registry           registry.ServicePlanRegistryIface
    serviceClassLookup registry.ServiceClassLookup
    emitter            OperationEmitter // nil-safe (FEATURE-0005)
}

func NewServicePlanHandler(
    reg registry.ServicePlanRegistryIface,
    serviceClassLookup registry.ServiceClassLookup,
    emitter OperationEmitter,
) *ServicePlanHandler
```

### HandleItem path parsing (Go 1.21, exactly two segments)

```
remainder := strings.TrimPrefix(r.URL.Path, "/v1/service-plans/")
parts := strings.Split(remainder, "/")
if len(parts) != 2 || parts[0] == "" || parts[1] == "" → 404
serviceClassName, name := parts[0], parts[1]
// dispatch GET → Get, PUT → Update, DELETE → Delete, else 405
```

### Create Flow

```
1. safeDecodeServicePlan (MaxBytesReader, status detection, DisallowUnknownFields)
2. decode error → writeError(400/413)
3. ValidateServicePlan(sp) → []FieldError; if errors → writeValidationErrors
4. serviceClassLookup.GetServiceClass(ctx, sp.Spec.ServiceClassName)
5. ErrNotFound → writeError(400, VALIDATION_FAILED, field="spec.serviceClassName",
     message="parent ServiceClass not found: <serviceClassName>")
6. Force: apiVersion, kind = ServicePlan, status.phase = Active
7. created, err := registry.CreateServicePlan(ctx, sp)
8. ErrAlreadyExists → writeError(409, RESOURCE_ALREADY_EXISTS)
9. emitOperation(ctx, h.emitter, OperationSpec{Type: OpCreateServicePlan,
     ResourceKind: ServicePlanKind, ResourceName: created.Metadata.Name,
     ServiceClassName: created.Spec.ServiceClassName, ServicePlanName: created.Metadata.Name,
     RequestID: requestIDFromContext(ctx)})
10. writeJSON(201, created)
```

### Get / List Flow

```
Get:  ValidateServicePlanPathSegments(serviceClassName, name); GetServicePlan; ErrNotFound → 404; 200
List: ListServicePlans; error → 500; writeJSON(200, {"items": items})  // [] when empty
```

### Update Flow

```
1. ValidateServicePlanPathSegments(serviceClassName, name); if errors → writeValidationErrors
2. safeDecodeServicePlan(w, r); decode error → writeError(400/413)
3. body.Spec.ServiceClassName present and == serviceClassName, else 400 (field="spec.serviceClassName")
4. body.Metadata.Name present and == name, else 400 (field="metadata.name")
5. ValidateServicePlan(sp); if errors → writeValidationErrors
6. serviceClassLookup.GetServiceClass(ctx, serviceClassName); ErrNotFound →
     writeError(400, VALIDATION_FAILED, field="spec.serviceClassName")
7. updated, err := registry.UpdateServicePlan(ctx, sp)
8. ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
9. emitOperation(ctx, h.emitter, OperationSpec{Type: OpUpdateServicePlan,
     ResourceKind: ServicePlanKind, ResourceName: updated.Metadata.Name,
     ServiceClassName: updated.Spec.ServiceClassName, ServicePlanName: updated.Metadata.Name,
     RequestID: requestIDFromContext(ctx)})
10. writeJSON(200, updated)
```

### Delete Flow

```
1. ValidateServicePlanPathSegments(serviceClassName, name); if errors → writeValidationErrors
2. err := registry.DeleteServicePlan(ctx, serviceClassName, name)
3. ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
4. emitOperation(ctx, h.emitter, OperationSpec{Type: OpDeleteServicePlan,
     ResourceKind: ServicePlanKind, ResourceName: name,
     ServiceClassName: serviceClassName, ServicePlanName: name, RequestID: requestIDFromContext(ctx)})
5. w.WriteHeader(204)
```

**Parent lookup:** ServicePlan create and update both verify the parent ServiceClass exists via the
narrow `ServiceClassLookup`; the ServicePlanRegistry never depends on the ServiceClassRegistry.
`spec.serviceClassName` is immutable — a plan is never moved between classes.

---

## Operation Emission Behavior

Both handlers reuse the FEATURE-0005 nil-safe `OperationEmitter` and the `emitOperation` helper.
Emission occurs only AFTER a successful mutating registry action and BEFORE writing the response.
Emission never affects the primary response; a nil emitter is skipped, and emitter errors are
swallowed.

| Action | Operation type | resourceKind | serviceClassName | servicePlanName |
|---|---|---|---|---|
| ServiceClass create | OpCreateServiceClass | ServiceClass | class name | — |
| ServiceClass update | OpUpdateServiceClass | ServiceClass | class name | — |
| ServiceClass delete | OpDeleteServiceClass | ServiceClass | class name | — |
| ServicePlan create | OpCreateServicePlan | ServicePlan | parent class name | plan name |
| ServicePlan update | OpUpdateServicePlan | ServicePlan | parent class name | plan name |
| ServicePlan delete | OpDeleteServicePlan | ServicePlan | parent class name | plan name |

No Operation is emitted on failed validation, duplicate create, missing parent, not-found, or
delete-blocked cases (Requirement 12.8, 17.1).

---

## Server and main.go Wiring

### Updated server.New signature

```go
func New(
    cfg config.Config,
    org *api.OrgHandler,
    ou *api.OUHandler,
    tenant *api.TenantHandler,
    project *api.ProjectHandler,
    operation *api.OperationHandler,
    serviceClass *api.ServiceClassHandler, // NEW
    servicePlan *api.ServicePlanHandler,   // NEW
    bootstrap *api.BootstrapHandler,
    readiness *health.ReadinessState,
) *Server
```

`server.New` SHALL register ServiceClass and ServicePlan routes when non-nil handlers are provided.
Production `main.go` wiring SHALL provide non-nil ServiceClass and ServicePlan handlers.

### Route registration (Go 1.21)

```go
mux.Handle("/v1/service-classes", chain(http.HandlerFunc(serviceClass.HandleCollection)))
mux.Handle("/v1/service-classes/", chain(http.HandlerFunc(serviceClass.HandleItem)))
mux.Handle("/v1/service-plans", chain(http.HandlerFunc(servicePlan.HandleCollection)))
mux.Handle("/v1/service-plans/", chain(http.HandlerFunc(servicePlan.HandleItem)))
```

Middleware order unchanged: requestID → logging → contentType → handler.

### main.go wiring

```go
serviceClassRegistry := registry.NewServiceClassRegistry()
servicePlanRegistry := registry.NewServicePlanRegistry()
serviceClassBlocker := registry.NewServicePlanChildBlockerChecker(servicePlanRegistry)

serviceClassHandler := api.NewServiceClassHandler(serviceClassRegistry, serviceClassBlocker, emitter)
servicePlanHandler := api.NewServicePlanHandler(servicePlanRegistry, serviceClassRegistry, emitter)

srv := server.New(cfg, orgHandler, ouHandler, tenantHandler, projectHandler,
    operationHandler, serviceClassHandler, servicePlanHandler, bootstrapHandler, readiness)
```

Reuses the single `emitter` created for FEATURE-0005; no duplicate Operation registry.

---

## Error Handling

### Error Mapping Table

| Condition | HTTP | error.code | error.field |
|---|---|---|---|
| ErrAlreadyExists (create) | 409 | RESOURCE_ALREADY_EXISTS | — |
| ErrNotFound (get/update/delete) | 404 | RESOURCE_NOT_FOUND | — |
| []FieldError non-empty | 400 | VALIDATION_FAILED | first field |
| Invalid category | 400 | VALIDATION_FAILED | spec.category |
| Invalid tier | 400 | VALIDATION_FAILED | spec.tier |
| Invalid lifecycle | 400 | VALIDATION_FAILED | spec.lifecycle |
| Invalid defaultPlanName | 400 | VALIDATION_FAILED | spec.defaultPlanName |
| Forbidden parameter key | 400 | VALIDATION_FAILED | spec.parameters |
| JSON syntax/type error | 400 | VALIDATION_FAILED | — |
| Body exceeds 1 MiB | 413 | VALIDATION_FAILED | — |
| Content-Type mismatch | 415 | VALIDATION_FAILED | — (middleware) |
| Status key in body | 400 | VALIDATION_FAILED | status |
| Parent ServiceClass not found (plan create/update) | 400 | VALIDATION_FAILED | spec.serviceClassName |
| metadata.name absent/mismatch (PUT) | 400 | VALIDATION_FAILED | metadata.name |
| spec.serviceClassName absent/mismatch (plan PUT) | 400 | VALIDATION_FAILED | spec.serviceClassName |
| ServiceClass has ServicePlans (delete) | 409 | DELETE_BLOCKED | — |
| Wrong item path segment count | 404 | RESOURCE_NOT_FOUND | — |
| Blocker/registry unexpected error | 500 | INTERNAL_ERROR | — |

Emission-path errors are NEVER mapped to client responses; `emitOperation` swallows them.

---

## Security and Privacy Constraints

- ServiceClass/ServicePlan store only catalog metadata; NO secrets, tokens, credentials, passwords.
- `spec.parameters` keys are rejected when they match the Requirement 4.6 forbidden substrings
  (case-insensitive); the plain substring `key` alone is allowed.
- The server SHALL NOT store raw request bodies and SHALL NOT echo raw bodies in error responses.
- The server SHALL NOT log secrets or raw request bodies.
- Operation records carry only non-sensitive reference fields (type, kind, names, request ID).

---

## Correctness Properties

Each correctness property is implemented as a `testing/quick` test with `Config{MaxCount: 100}` and
tagged `// Feature: serviceclass-serviceplan, Property N: <title>`.

### Property 1: ValidateServiceClass accepts valid inputs

For any DNS-label `metadata.name` with a valid category and lifecycle, `ValidateServiceClass`
returns no `FieldError`.

**Validates: Requirements 3.1, 3.2, 3.3**

### Property 2: ValidateServiceClass rejects invalid names

For any string that is not a valid DNS-label, `ValidateServiceClass` returns a `FieldError` with
`Field = "metadata.name"`.

**Validates: Requirements 3.1**

### Property 3: ValidateServicePlan rejects forbidden parameter keys

For any parameter key containing a Requirement 4.6 secret-bearing substring (case-insensitive),
`ValidateServicePlan` returns a `FieldError` with `Field = "spec.parameters"`; benign keys like
`regionKey` and plain `key` are accepted.

**Validates: Requirements 4.6, 13.2**

### Property 4: ServiceClass Create/Get round trip preserves data

For any valid ServiceClass, `CreateServiceClass` then `GetServiceClass` returns an equal resource.

**Validates: Requirements 5.2, 5.3**

### Property 5: ServicePlan composite identity is stable

For any valid ServicePlan, `CreateServicePlan` then `GetServicePlan` by `serviceClassName/name`
returns an equal resource, and the same name under a different class does not collide.

**Validates: Requirements 6.1, 6.2**

### Property 6: List ordering is deterministic

`ListServiceClasses` is sorted by name; `ListServicePlans` is sorted by serviceClassName then name.

**Validates: Requirements 5.7, 6.8**

### Property 7: Registries return deep copies

Mutating a value returned by any Create/Get/List/Update call never affects stored registry state.

**Validates: Requirements 5.2, 6.3**

### Property 8: Duplicate create never overwrites

A second `CreateServiceClass`/`CreateServicePlan` with an existing key returns `ErrAlreadyExists`
and leaves the stored entry unchanged.

**Validates: Requirements 5.6, 6.7**

---

## Testing Strategy

### Validation tests
- `serviceclass_test.go`: valid names accepted; empty/invalid/long `metadata.name` rejected;
  invalid/empty category and lifecycle rejected; non-DNS-label `defaultPlanName` rejected; empty
  `defaultPlanName` accepted; path segment validation.
- `serviceplan_test.go`: valid inputs accepted; empty/invalid `metadata.name`/`serviceClassName`
  rejected; invalid/empty tier and lifecycle rejected; forbidden parameter keys rejected
  (case-insensitive, including `apiKey`, `ACCESSKEY`); benign keys (`regionKey`, `key`) accepted;
  path segment field mapping.
- Property tests: `serviceclass_property_test.go`, `serviceplan_property_test.go` (Properties 1–3).

### Registry tests
- `serviceclass_registry_test.go`: Create, Get, List sorted, Update preserves immutables, Delete,
  duplicate → ErrAlreadyExists (original unchanged), not-found errors, empty list → non-nil `[]`.
- `serviceplan_registry_test.go`: same plus `CountByServiceClass`, composite-key isolation (same
  name under different classes), two-level sort order.
- Property tests: Properties 4–8 across both registries.
- Race tests: `serviceclass_registry_race_test.go`, `serviceplan_registry_race_test.go` — 10+
  goroutines mixed CRUD (+ CountByServiceClass), zero race reports.

### Blocker tests
- `serviceclass_blocker_test.go`: `CountByServiceClass` correct; `BlockedByServiceClassChildren`
  returns a `ServicePlan` blocker when count > 0; returns nil when 0; registry error propagates;
  Retired-lifecycle plan still blocks.

### Handler tests
- `serviceclass_handler_test.go`: POST 201/409/400 (invalid fields, status key, bad JSON, unknown
  field), 413 oversized; GET 200/404/400; wrong path shape → 404; list sorted/empty; PUT 200/404/400
  (name mismatch); DELETE 204/404; DELETE 409 DELETE_BLOCKED with a plan present; DELETE 204 with
  zero plans; nil emitter/blocker no-panic.
- `serviceplan_handler_test.go`: POST 201/409/400 (invalid fields, missing parent, status key, bad
  JSON), 413 oversized; GET 200/404/400; wrong path shape → 404; list sorted/empty; PUT 200/404/400
  (name or serviceClassName mismatch, missing parent on update); DELETE 204/404; nil emitter no-panic.

### Operation emission tests
- Table-driven across both resources × create/update/delete: correct Operation type, resourceKind,
  serviceClassName/servicePlanName recorded; failed actions emit nothing; emission failure does not
  change the primary response.

### Server tests
- `server_test.go` (modified): constructor with ServiceClass/ServicePlan handler fixtures; route
  registration for all four patterns.

### Verification
- `make fmt && make vet && make test && go test -race ./... && make build`
- If host Go is unavailable, use the Docker fallback:
  `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./... && go build ./cmd/sovrunn-api'`

---

## Non-Goals

1. NO ServiceInstance or ServiceBinding resources.
2. NO service/datastore provisioning, Kubernetes operators, Crossplane, Kratix, or GitOps.
3. NO Plugin registry, Capability registry, or capability matching.
4. NO ServiceOps execution or plugin execution.
5. NO async workflows, approval flows, queues, or background workers.
6. NO pricing/billing engine, quota enforcement, subscription model, or marketplace publishing.
7. NO persistence/database storage, auth/RBAC, AI automation, or UI.
8. NO verification that `spec.defaultPlanName` references an existing plan.
9. NO generic delete-blocker framework (narrow per-parent interface only).
10. NO list filtering or pagination.
11. NO Go 1.22 wildcard routing; no new external dependencies.
12. NO secrets stored in catalog resources or logs.

---

## Design Questions (Resolved)

| # | Question | Resolution |
|---|---|---|
| 1 | Parent lookup interface | Narrow `ServiceClassLookup` in `internal/registry`; existing `*ServiceClassRegistry` satisfies it; `ServicePlanRegistry` never depends on `ServiceClassRegistry` |
| 2 | Delete-blocker interface shape | `ServiceClassChildBlocker.BlockedByServiceClassChildren(ctx, serviceClassName)`; implemented by `ServicePlanChildBlockerChecker` via `CountByServiceClass`; no generic framework |
| 3 | Parameters value type | `map[string]string` in Phase 1 to avoid nested secret-bearing structures |
