# Design Document — FEATURE-0007 Plugin and Capability Registry

## Overview

FEATURE-0007 implements `Plugin` and `Capability` as **global platform registry resources**.
A Plugin declares an implementation unit that performs lifecycle operations for a service family
or provider. A Capability declares a specific lifecycle action that a plugin supports for a given
ServiceClass.

**Global scope:** neither resource is scoped to the Organization/OU/Tenant/Project hierarchy.
Plugin identity is `metadata.name`; Capability identity is `metadata.name`. Both are simple
(non-composite) keys.

**Registry-only:** these resources record what a plugin claims to support. They do NOT execute
plugins, load dynamic code, call external runtimes, or provision infrastructure.

**Relationship to FEATURE-0006 catalog:** Plugin declares `serviceClassRefs` referencing existing
ServiceClasses. Capability references a Plugin via `pluginRef` and a ServiceClass via
`serviceClassRef`. Reference validation occurs at write time only.

## Resolved Design Decisions

| # | Question | Resolution |
|---|---|---|
| 1 | ServiceClass lookup interface | A single `ServiceClassLookup` interface (already defined in `internal/registry` from FEATURE-0006) serves both Plugin and Capability handlers. The existing `*ServiceClassRegistry` satisfies it via `GetServiceClass`. No new interface needed for ServiceClass lookup. |
| 2 | Plugin lookup interface | A new narrow `PluginLookup` interface is defined in `internal/registry/plugin_registry.go`. It exposes `GetPlugin(ctx, name) (resources.Plugin, error)`. The existing `*PluginRegistry` satisfies it. `PluginLookup` is injected into the CapabilityHandler for `spec.pluginRef` validation. It is separate from the delete-blocker interface. |
| 3 | Delete-blocker interface shape | `PluginChildBlocker` with `BlockedByPluginChildren(ctx, pluginName) ([]registry.BlockedBy, error)`, implemented by `CapabilityChildBlockerChecker` via `CountByPlugin`. Naming follows the established pattern (`BlockedByServiceClassChildren`, `BlockedByOUChildren`). No generic blocker framework. |
| 4 | Capability immutability rationale | Capability does not support PUT. Capabilities are immutable after creation; to change, delete and recreate. This simplifies the model — a capability declaration either exists or does not. A future phase may add PUT if versioned capability evolution is needed. |
| 5 | Capability query parameter filtering | Phase 1 supports `pluginRef` and `serviceClassRef` as optional query parameters on `GET /v1/capabilities`. The `operation` filter is deferred to a future phase. |
| 6 | ServiceClassRefs validation scope | On Plugin update, the full `spec.serviceClassRefs` list is re-validated against the ServiceClassRegistry. This is simpler and safer than tracking deltas. |
| 7 | Operation spec fields | RESOLVED in requirements: `pluginName` and `capabilityName` are locked as the canonical field names. |
| 8 | k8sOps plugin type | RESOLVED in requirements: `k8sOps` is accepted for Phase 1. |

## Architecture

```
cmd/sovrunn-api → internal/server → internal/api → internal/registry
                                          │              │
                                          │              ├─ PluginRegistry (storage-only)
                                          │              ├─ CapabilityRegistry (storage-only)
                                          │              └─ CapabilityChildBlockerChecker
                                          └─ PluginHandler / CapabilityHandler
                                             (emit Operations via FEATURE-0005 emitter)
                        internal/resources (Plugin, Capability models + constants)
```

- Plugin/Capability handlers reuse: safe decoding, `writeError`/`writeJSON`, the nil-safe
  `OperationEmitter`, and the API-local `requestIDFromContext`.
- Plugin handler holds a `ServiceClassLookup` for `serviceClassRefs` verification and a
  `PluginChildBlocker` for delete blocking.
- Capability handler holds a `PluginLookup` for `pluginRef` verification and a
  `ServiceClassLookup` for `serviceClassRef` verification. It also reads the resolved Plugin's
  `serviceClassRefs` to verify that the capability's `serviceClassRef` is declared by the plugin.
- `internal/api` MUST NOT import `internal/server`; request ID is read via the existing API-local
  `requestIDFromContext`.

## Files to Create/Modify

### New Files

| File | Purpose |
|---|---|
| `internal/resources/plugin.go` | Plugin struct, PluginSpec, PluginStatus, PluginType/DeploymentMode constants, kind constant |
| `internal/resources/capability.go` | Capability struct, CapabilitySpec, CapabilityStatus, CapabilityOperation constants, kind constant |
| `internal/registry/plugin_registry.go` | PluginRegistryIface, PluginRegistry, PluginLookup interface |
| `internal/registry/capability_registry.go` | CapabilityRegistryIface, CapabilityRegistry, CountByPlugin |
| `internal/registry/plugin_blocker.go` | PluginChildBlocker interface, CapabilityChildBlockerChecker |
| `internal/validation/plugin.go` | ValidatePlugin, ValidatePluginPathSegment |
| `internal/validation/capability.go` | ValidateCapability, ValidateCapabilityPathSegment |
| `internal/api/plugin_decode.go` | safeDecodePlugin |
| `internal/api/capability_decode.go` | safeDecodeCapability |
| `internal/api/plugin_handler.go` | PluginHandler (CRUD) |
| `internal/api/capability_handler.go` | CapabilityHandler (Create, Get, List, Delete; PUT → 405) |
| plus `_test.go` files for each package | unit/property/race/handler/emission tests |

### Modified Files

| File | Change |
|---|---|
| `internal/resources/errors.go` | Add `ErrCodeMethodNotAllowed ErrorCode = "METHOD_NOT_ALLOWED"` |
| `internal/resources/operation.go` | Add 5 op type constants + `PluginName`/`CapabilityName` OperationSpec fields + `PluginKind`/`CapabilityKind` resource kind constants |
| `internal/server/server.go` | Accept `*api.PluginHandler`, `*api.CapabilityHandler`; register routes |
| `internal/server/server_test.go` | Update constructor tests; add route tests |
| `cmd/sovrunn-api/main.go` | Wire registries, blocker, lookups, handlers |

## Data Models

### internal/resources/plugin.go

```go
package resources

// Plugin is a global platform registry resource declaring an implementation
// unit that performs lifecycle operations for a service family or provider.
// Identity: metadata.name (simple key). Follows canonical metadata/spec/status.
type Plugin struct {
    APIVersion string       `json:"apiVersion"`
    Kind       string       `json:"kind"`
    Metadata   Metadata     `json:"metadata"`
    Spec       PluginSpec   `json:"spec"`
    Status     PluginStatus `json:"status"`
}

type PluginSpec struct {
    PluginType       string   `json:"pluginType"`
    Version          string   `json:"version"`
    ServiceClassRefs []string `json:"serviceClassRefs"`
    DeploymentMode   string   `json:"deploymentMode"`
    Description      string   `json:"description,omitempty"`
    Tags             []string `json:"tags,omitempty"`
}

type PluginStatus struct {
    Phase   string `json:"phase"`
    Message string `json:"message,omitempty"`
}

const (
    PluginAPIVersion = "platform.sovrunn.io/v1alpha1"
    PluginKind       = "Plugin"
)

// PluginType constants.
const (
    PluginTypeDStoreOps  = "dStoreOps"
    PluginTypeCacheOps   = "cacheOps"
    PluginTypeStreamOps  = "streamOps"
    PluginTypeObjectOps  = "objectOps"
    PluginTypeGatewayOps = "gatewayOps"
    PluginTypeFaasOps    = "faasOps"
    PluginTypeLBOps      = "lbOps"
    PluginTypeK8sOps     = "k8sOps"
    PluginTypeBigDataOps = "bigDataOps"
    PluginTypeSdeOps     = "sdeOps"
)

// DeploymentMode constants. Phase 1 accepts only CompiledIn.
const (
    DeploymentModeCompiledIn = "compiled-in"
)
```

### internal/resources/capability.go

```go
package resources

// Capability is a global platform registry resource declaring a specific
// lifecycle action supported by a plugin for a given ServiceClass.
// Identity: metadata.name (simple key). Follows canonical metadata/spec/status.
type Capability struct {
    APIVersion string           `json:"apiVersion"`
    Kind       string           `json:"kind"`
    Metadata   Metadata         `json:"metadata"`
    Spec       CapabilitySpec   `json:"spec"`
    Status     CapabilityStatus `json:"status"`
}

type CapabilitySpec struct {
    PluginRef       string `json:"pluginRef"`
    ServiceClassRef string `json:"serviceClassRef"`
    Operation       string `json:"operation"`
    Supported       bool   `json:"supported"`
    Description     string `json:"description,omitempty"`
}

type CapabilityStatus struct {
    Phase   string `json:"phase"`
    Message string `json:"message,omitempty"`
}

const (
    CapabilityAPIVersion = "platform.sovrunn.io/v1alpha1"
    CapabilityKind       = "Capability"
)

// CapabilityOperation constants.
const (
    CapOpValidate          = "Validate"
    CapOpPlan              = "Plan"
    CapOpProvision         = "Provision"
    CapOpConfigure         = "Configure"
    CapOpBind              = "Bind"
    CapOpObserve           = "Observe"
    CapOpScale             = "Scale"
    CapOpUpgrade           = "Upgrade"
    CapOpBackup            = "Backup"
    CapOpRestore           = "Restore"
    CapOpRotateCredentials = "RotateCredentials"
    CapOpUnbind            = "Unbind"
    CapOpDelete            = "Delete"
)
```

### internal/resources/operation.go — additions

```go
// Plugin and Capability operation type constants (FEATURE-0007).
const (
    OpCreatePlugin     = "CreatePlugin"
    OpUpdatePlugin     = "UpdatePlugin"
    OpDeletePlugin     = "DeletePlugin"
    OpCreateCapability = "CreateCapability"
    OpDeleteCapability = "DeleteCapability"
)

// OperationSpec gains two optional plugin-reference fields:
//   PluginName     string `json:"pluginName,omitempty"`
//   CapabilityName string `json:"capabilityName,omitempty"`
// These are used only for Plugin and Capability Operation records.
```

### internal/resources/errors.go — addition

```go
const (
    ErrCodeMethodNotAllowed ErrorCode = "METHOD_NOT_ALLOWED"
)
```

Reuses the existing `Metadata` struct and phase constants (`PhaseActive`, etc.).

---

## Components and Interfaces

| Component | Package | Responsibility |
|---|---|---|
| `Plugin` / `Capability` models | `internal/resources` | Data shape, type/mode/operation constants, kind constants |
| `PluginRegistryIface` / `PluginRegistry` | `internal/registry` | Storage-only in-memory Plugin store |
| `CapabilityRegistryIface` / `CapabilityRegistry` | `internal/registry` | Storage-only in-memory Capability store with filtering |
| `PluginLookup` interface | `internal/registry` | Narrow existence-check for Capability handler |
| `ServiceClassLookup` interface (existing) | `internal/registry` | Narrow existence-check reused by Plugin and Capability handlers |
| `PluginChildBlocker` interface | `internal/registry` | Narrow delete-blocker boundary for Plugin delete |
| `CapabilityChildBlockerChecker` | `internal/registry` | Implements the blocker via `CountByPlugin` |
| `ValidatePlugin` / `ValidateCapability` | `internal/validation` | Pure, context-free field validation |
| `safeDecodePlugin` / `safeDecodeCapability` | `internal/api` | Safe JSON decoding (1 MiB, DisallowUnknownFields, status rejection) |
| `PluginHandler` | `internal/api` | CRUD handler; emits Operations via FEATURE-0005 emitter |
| `CapabilityHandler` | `internal/api` | Create/Get/List/Delete handler; PUT → 405; emits Operations |
| `OperationEmitter` (reused) | `internal/api` | Emission boundary consumed by both handlers (nil-safe) |

---

## Validation Design

### internal/validation/plugin.go

```go
// ValidatePlugin validates all user-authored Plugin fields. Context-free, no I/O.
func ValidatePlugin(p resources.Plugin) []resources.FieldError

// ValidatePluginPathSegment validates the single URL path segment (name).
func ValidatePluginPathSegment(name string) []resources.FieldError
```

**Rules:**
- `metadata.name` DNS-label (1–63) → `error.field = "metadata.name"`.
- `spec.pluginType` required, ∈ {dStoreOps, cacheOps, streamOps, objectOps, gatewayOps, faasOps, lbOps, k8sOps, bigDataOps, sdeOps} → `error.field = "spec.pluginType"`.
- `spec.version` required, non-empty → `error.field = "spec.version"`.
- `spec.serviceClassRefs` required, non-nil, len ≥ 1 → `error.field = "spec.serviceClassRefs"`.
- Each entry in `spec.serviceClassRefs` must be a valid DNS-label (1–63) → `error.field = "spec.serviceClassRefs"`.
- `spec.deploymentMode` required, ∈ {compiled-in} → `error.field = "spec.deploymentMode"`.
- `spec.description` and `spec.tags` optional; not format-validated.

### internal/validation/capability.go

```go
// ValidateCapability validates all user-authored Capability fields. Context-free, no I/O.
func ValidateCapability(c resources.Capability) []resources.FieldError

// ValidateCapabilityPathSegment validates the single URL path segment (name).
func ValidateCapabilityPathSegment(name string) []resources.FieldError
```

**Rules:**
- `metadata.name` DNS-label (1–63) → `error.field = "metadata.name"`.
- `spec.pluginRef` required, DNS-label (1–63) → `error.field = "spec.pluginRef"`.
- `spec.serviceClassRef` required, DNS-label (1–63) → `error.field = "spec.serviceClassRef"`.
- `spec.operation` required, ∈ {Validate, Plan, Provision, Configure, Bind, Observe, Scale, Upgrade, Backup, Restore, RotateCredentials, Unbind, Delete} → `error.field = "spec.operation"`.
- `spec.supported` defaults to `false` (Go zero-value); no validation error for absent boolean.
- `spec.description` optional; not format-validated.

Both validators reuse `validateName` (DNS-label check) from the validation package.
Enum membership checks use switch statements over the resource constants.

---

## Safe JSON Decoding

### internal/api/plugin_decode.go

```go
func safeDecodePlugin(w http.ResponseWriter, r *http.Request) (resources.Plugin, error)
```

### internal/api/capability_decode.go

```go
func safeDecodeCapability(w http.ResponseWriter, r *http.Request) (resources.Capability, error)
```

**Sequence (identical to existing decoder pattern):**
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

## PluginRegistry Design

### internal/registry/plugin_registry.go

```go
type PluginRegistryIface interface {
    CreatePlugin(ctx context.Context, p resources.Plugin) (resources.Plugin, error)
    GetPlugin(ctx context.Context, name string) (resources.Plugin, error)
    ListPlugins(ctx context.Context) ([]resources.Plugin, error)
    UpdatePlugin(ctx context.Context, p resources.Plugin) (resources.Plugin, error)
    DeletePlugin(ctx context.Context, name string) error
}

// PluginLookup is the narrow interface for verifying Plugin existence.
// Injected into CapabilityHandler. The concrete *PluginRegistry satisfies it.
type PluginLookup interface {
    GetPlugin(ctx context.Context, name string) (resources.Plugin, error)
}

type PluginRegistry struct {
    mu    sync.RWMutex
    store map[string]resources.Plugin // key: metadata.name
}

func NewPluginRegistry() *PluginRegistry {
    return &PluginRegistry{store: make(map[string]resources.Plugin)}
}
```

### Method Behavior

| Method | Lock | Returns | Error |
|---|---|---|---|
| CreatePlugin | Lock | deep copy of stored | ErrAlreadyExists if name exists |
| GetPlugin | RLock | deep copy | ErrNotFound if absent |
| ListPlugins | RLock | sorted slice of deep copies | nil |
| UpdatePlugin | Lock | deep copy of updated | ErrNotFound if absent |
| DeletePlugin | Lock | — | ErrNotFound if absent |

**List sort order:** `Metadata.Name` ascending.

**deepCopyPlugin:** copies the `ServiceClassRefs` slice, `Tags` slice, and
`Metadata.Labels`/`Metadata.Annotations` maps so callers cannot mutate stored state.

**UpdatePlugin behavior:** derives the key from `p.Metadata.Name`, looks up the stored entry,
preserves stored `APIVersion`, `Kind`, `Status`, `Metadata.Name`, and replaces only mutable fields:
- `Metadata.Labels`
- `Metadata.Annotations`
- `Spec.PluginType`
- `Spec.Version`
- `Spec.ServiceClassRefs`
- `Spec.DeploymentMode`
- `Spec.Description`
- `Spec.Tags`

Returns a deep copy of the updated stored Plugin.

Reuses sentinel errors `ErrNotFound` and `ErrAlreadyExists` from `internal/registry`.
No dependency on other registries; no package-level global state.

---

## CapabilityRegistry Design

### internal/registry/capability_registry.go

```go
type CapabilityRegistryIface interface {
    CreateCapability(ctx context.Context, c resources.Capability) (resources.Capability, error)
    GetCapability(ctx context.Context, name string) (resources.Capability, error)
    ListCapabilities(ctx context.Context, pluginRef, serviceClassRef string) ([]resources.Capability, error)
    DeleteCapability(ctx context.Context, name string) error
    CountByPlugin(ctx context.Context, pluginName string) (int, error)
}

type CapabilityRegistry struct {
    mu    sync.RWMutex
    store map[string]resources.Capability // key: metadata.name
}

func NewCapabilityRegistry() *CapabilityRegistry {
    return &CapabilityRegistry{store: make(map[string]resources.Capability)}
}
```

### Method Behavior

| Method | Lock | Returns | Error |
|---|---|---|---|
| CreateCapability | Lock | deep copy of stored | ErrAlreadyExists if name exists |
| GetCapability | RLock | deep copy | ErrNotFound if absent |
| ListCapabilities | RLock | sorted, filtered slice of deep copies | nil |
| DeleteCapability | Lock | — | ErrNotFound if absent |
| CountByPlugin | RLock | int count | nil |

**List sort order:** `Metadata.Name` ascending.

**ListCapabilities filtering:** accepts `pluginRef` and `serviceClassRef` as filter strings. When
both are non-empty, results match both (AND logic). When either is empty, it is not applied. When
neither is provided, all entries are returned. Returns a non-nil empty slice when no matches.

**deepCopyCapability:** copies `Metadata.Labels`/`Metadata.Annotations` maps. `CapabilitySpec`
contains only scalar fields so no slice/map copy is needed for spec.

**CountByPlugin:** RLock; iterates entries; counts those where `Spec.PluginRef` matches. Used only
by the delete blocker.

**No Update method:** Capability is immutable after creation; the registry does not expose an
update method.

Reuses sentinel errors `ErrNotFound` and `ErrAlreadyExists`. No dependency on other registries.

---

## Plugin Delete Blocker Design

### internal/registry/plugin_blocker.go

```go
// PluginChildBlocker is injected into the Plugin delete path.
// It reports child resources that block deleting a specific Plugin.
type PluginChildBlocker interface {
    BlockedByPluginChildren(ctx context.Context, pluginName string) ([]BlockedBy, error)
}

// CapabilityChildBlockerChecker implements PluginChildBlocker for Capability
// resources. It queries the CapabilityRegistry to determine whether any
// Capabilities reference the Plugin being deleted.
type CapabilityChildBlockerChecker struct {
    capabilityRegistry CapabilityRegistryIface
}

func NewCapabilityChildBlockerChecker(reg CapabilityRegistryIface) *CapabilityChildBlockerChecker

func (c *CapabilityChildBlockerChecker) BlockedByPluginChildren(
    ctx context.Context, pluginName string,
) ([]BlockedBy, error) {
    count, err := c.capabilityRegistry.CountByPlugin(ctx, pluginName)
    if err != nil { return nil, err }
    if count > 0 {
        return []BlockedBy{{Kind: "Capability", Count: count}}, nil
    }
    return nil, nil
}
```

Reuses the existing `BlockedBy` type. No generic blocker framework. Consistent with
`ServicePlanChildBlockerChecker` pattern from FEATURE-0006.

---

## Plugin HTTP Handler Design

### internal/api/plugin_handler.go

```go
type PluginHandler struct {
    registry           registry.PluginRegistryIface
    serviceClassLookup registry.ServiceClassLookup  // verifies serviceClassRefs
    blocker            registry.PluginChildBlocker   // nil-safe (delete blocking)
    emitter            OperationEmitter              // nil-safe (FEATURE-0005)
}

func NewPluginHandler(
    reg registry.PluginRegistryIface,
    serviceClassLookup registry.ServiceClassLookup,
    blocker registry.PluginChildBlocker,
    emitter OperationEmitter,
) *PluginHandler

func (h *PluginHandler) HandleCollection(w http.ResponseWriter, r *http.Request) // POST/GET
func (h *PluginHandler) HandleItem(w http.ResponseWriter, r *http.Request)       // GET/PUT/DELETE
```

### HandleItem path parsing (Go 1.21, exactly one segment)

```
remainder := strings.TrimPrefix(r.URL.Path, "/v1/plugins/")
parts := strings.Split(remainder, "/")
if len(parts) != 1 || parts[0] == "" → 404
name := parts[0]
// dispatch GET → Get, PUT → Update, DELETE → Delete, else 405
```

### Create Flow

```
1. safeDecodePlugin (MaxBytesReader, status detection, DisallowUnknownFields)
2. decode error → writeError(400/413)
3. ValidatePlugin(p) → []FieldError; if errors → writeValidationErrors
4. Reference validation: for each ref in p.Spec.ServiceClassRefs:
     serviceClassLookup.GetServiceClass(ctx, ref)
     if ErrNotFound → writeError(400, VALIDATION_FAILED, field="spec.serviceClassRefs",
       message="referenced ServiceClass not found: <ref>")
5. Force: apiVersion, kind = Plugin, status.phase = Active
6. created, err := registry.CreatePlugin(ctx, p)
7. ErrAlreadyExists → writeError(409, RESOURCE_ALREADY_EXISTS)
8. emitOperation(ctx, h.emitter, OperationSpec{Type: OpCreatePlugin,
     ResourceKind: PluginKind, ResourceName: created.Metadata.Name,
     PluginName: created.Metadata.Name, RequestID: requestIDFromContext(ctx)})
9. writeJSON(201, created)
```

### Get / List Flow

```
Get:  ValidatePluginPathSegment(name); GetPlugin; ErrNotFound → 404; 200
List: ListPlugins; error → 500; writeJSON(200, {"items": items})  // [] when empty
```

### Update Flow

```
1. ValidatePluginPathSegment(name); if errors → writeValidationErrors
2. safeDecodePlugin(w, r); decode error → writeError(400/413)
3. body.Metadata.Name present and == name, else 400 (field="metadata.name")
4. ValidatePlugin(p); if errors → writeValidationErrors
5. Reference validation: for each ref in p.Spec.ServiceClassRefs:
     serviceClassLookup.GetServiceClass(ctx, ref)
     if ErrNotFound → writeError(400, VALIDATION_FAILED, field="spec.serviceClassRefs",
       message="referenced ServiceClass not found: <ref>")
6. updated, err := registry.UpdatePlugin(ctx, p)
7. ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
8. emitOperation(ctx, h.emitter, OperationSpec{Type: OpUpdatePlugin,
     ResourceKind: PluginKind, ResourceName: updated.Metadata.Name,
     PluginName: updated.Metadata.Name, RequestID: requestIDFromContext(ctx)})
9. writeJSON(200, updated)
```

### Delete Flow

```
1. ValidatePluginPathSegment(name); if errors → writeValidationErrors
2. if h.blocker != nil:
     blockers, err := h.blocker.BlockedByPluginChildren(ctx, name)
     err → writeError(500, INTERNAL_ERROR)
     blockers non-empty → writeError(409, DELETE_BLOCKED, message "deletion blocked by Capability resources")
3. err := registry.DeletePlugin(ctx, name)
4. ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
5. emitOperation(ctx, h.emitter, OperationSpec{Type: OpDeletePlugin,
     ResourceKind: PluginKind, ResourceName: name,
     PluginName: name, RequestID: requestIDFromContext(ctx)})
6. w.WriteHeader(204)
```

**Nil-safe blocker:** if `blocker` is nil, delete proceeds without child checks (preserves isolated
handler test compatibility). Production wiring injects `CapabilityChildBlockerChecker`.

---

## Capability HTTP Handler Design

### internal/api/capability_handler.go

```go
type CapabilityHandler struct {
    registry           registry.CapabilityRegistryIface
    pluginLookup       registry.PluginLookup           // verifies pluginRef
    serviceClassLookup registry.ServiceClassLookup     // verifies serviceClassRef
    emitter            OperationEmitter                 // nil-safe (FEATURE-0005)
}

func NewCapabilityHandler(
    reg registry.CapabilityRegistryIface,
    pluginLookup registry.PluginLookup,
    serviceClassLookup registry.ServiceClassLookup,
    emitter OperationEmitter,
) *CapabilityHandler

func (h *CapabilityHandler) HandleCollection(w http.ResponseWriter, r *http.Request) // POST/GET
func (h *CapabilityHandler) HandleItem(w http.ResponseWriter, r *http.Request)       // GET/DELETE; PUT → 405
```

### HandleItem path parsing (Go 1.21, exactly one segment)

```
remainder := strings.TrimPrefix(r.URL.Path, "/v1/capabilities/")
parts := strings.Split(remainder, "/")
if len(parts) != 1 || parts[0] == "" → 404
name := parts[0]
// dispatch GET → Get, PUT → write 405 METHOD_NOT_ALLOWED, DELETE → Delete, else 405
```

### HandleItem PUT specifically

```
case http.MethodPut:
    writeError(w, r, http.StatusMethodNotAllowed, resources.ErrCodeMethodNotAllowed,
        "Capability does not support update; delete and recreate instead", "", "")
```

This returns 405 with `error.code = "METHOD_NOT_ALLOWED"` regardless of whether a Capability
with that name exists (no existence check performed).

### Create Flow

```
1. safeDecodeCapability (MaxBytesReader, status detection, DisallowUnknownFields)
2. decode error → writeError(400/413)
3. ValidateCapability(c) → []FieldError; if errors → writeValidationErrors
4. Plugin reference validation:
     plugin, err := pluginLookup.GetPlugin(ctx, c.Spec.PluginRef)
     if ErrNotFound → writeError(400, VALIDATION_FAILED, field="spec.pluginRef",
       message="referenced Plugin not found: <pluginRef>")
5. ServiceClass reference validation:
     _, err := serviceClassLookup.GetServiceClass(ctx, c.Spec.ServiceClassRef)
     if ErrNotFound → writeError(400, VALIDATION_FAILED, field="spec.serviceClassRef",
       message="referenced ServiceClass not found: <serviceClassRef>")
6. ServiceClass-in-Plugin validation:
     check c.Spec.ServiceClassRef ∈ plugin.Spec.ServiceClassRefs
     if NOT present → writeError(400, VALIDATION_FAILED, field="spec.serviceClassRef",
       message="ServiceClass <ref> is not declared by Plugin <pluginRef>")
7. Force: apiVersion, kind = Capability, status.phase = Active
8. created, err := registry.CreateCapability(ctx, c)
9. ErrAlreadyExists → writeError(409, RESOURCE_ALREADY_EXISTS)
10. emitOperation(ctx, h.emitter, OperationSpec{Type: OpCreateCapability,
      ResourceKind: CapabilityKind, ResourceName: created.Metadata.Name,
      PluginName: created.Spec.PluginRef, CapabilityName: created.Metadata.Name,
      RequestID: requestIDFromContext(ctx)})
11. writeJSON(201, created)
```

### Get / List Flow

```
Get:  ValidateCapabilityPathSegment(name); GetCapability; ErrNotFound → 404; 200
List: read query params r.URL.Query().Get("pluginRef"), r.URL.Query().Get("serviceClassRef")
      ListCapabilities(ctx, pluginRef, serviceClassRef); error → 500
      writeJSON(200, {"items": items})  // [] when empty
```

### Delete Flow

```
1. ValidateCapabilityPathSegment(name); if errors → writeValidationErrors
2. err := registry.DeleteCapability(ctx, name)
3. ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
4. emitOperation(ctx, h.emitter, OperationSpec{Type: OpDeleteCapability,
     ResourceKind: CapabilityKind, ResourceName: name,
     PluginName: <looked up from stored capability before delete? see note>,
     CapabilityName: name, RequestID: requestIDFromContext(ctx)})
5. w.WriteHeader(204)
```

**Delete emission note:** To populate `PluginName` in the Operation record, the handler must
retrieve the Capability (via `GetCapability`) BEFORE calling `DeleteCapability`. If the get
succeeds, use `cap.Spec.PluginRef` as `PluginName`. If the get fails (race condition, should not
happen under normal single-server operation), emit without `PluginName`.

Revised delete flow:

```
1. ValidateCapabilityPathSegment(name); if errors → writeValidationErrors
2. cap, err := registry.GetCapability(ctx, name)
   ErrNotFound → writeError(404, RESOURCE_NOT_FOUND); return
3. err := registry.DeleteCapability(ctx, name)
   ErrNotFound → writeError(404, RESOURCE_NOT_FOUND); return  // unlikely race
4. emitOperation(ctx, h.emitter, OperationSpec{Type: OpDeleteCapability,
     ResourceKind: CapabilityKind, ResourceName: name,
     PluginName: cap.Spec.PluginRef, CapabilityName: name,
     RequestID: requestIDFromContext(ctx)})
5. w.WriteHeader(204)
```

---

## Operation Emission Behavior

Both handlers reuse the FEATURE-0005 nil-safe `OperationEmitter` and the `emitOperation` helper.
Emission occurs only AFTER a successful mutating registry action and BEFORE writing the response.
Emission never affects the primary response; a nil emitter is skipped, and emitter errors are
swallowed.

| Action | Operation type | resourceKind | pluginName | capabilityName |
|---|---|---|---|---|
| Plugin create | OpCreatePlugin | Plugin | plugin name | — |
| Plugin update | OpUpdatePlugin | Plugin | plugin name | — |
| Plugin delete | OpDeletePlugin | Plugin | plugin name | — |
| Capability create | OpCreateCapability | Capability | referenced plugin name | capability name |
| Capability delete | OpDeleteCapability | Capability | referenced plugin name | capability name |

No Operation is emitted on failed validation, duplicate create, missing reference, not-found,
delete-blocked, or method-not-allowed cases (Requirements 13.8, 18.1).

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
    serviceClass *api.ServiceClassHandler,
    servicePlan *api.ServicePlanHandler,
    plugin *api.PluginHandler,         // NEW
    capability *api.CapabilityHandler,  // NEW
    bootstrap *api.BootstrapHandler,
    readiness *health.ReadinessState,
) *Server
```

`server.New` SHALL register Plugin and Capability routes when non-nil handlers are provided.

### Route registration (Go 1.21)

```go
if plugin != nil {
    mux.Handle("/v1/plugins", chain(http.HandlerFunc(plugin.HandleCollection)))
    mux.Handle("/v1/plugins/", chain(http.HandlerFunc(plugin.HandleItem)))
}
if capability != nil {
    mux.Handle("/v1/capabilities", chain(http.HandlerFunc(capability.HandleCollection)))
    mux.Handle("/v1/capabilities/", chain(http.HandlerFunc(capability.HandleItem)))
}
```

Middleware order unchanged: requestID → logging → contentType → handler.

### main.go wiring

```go
pluginRegistry := registry.NewPluginRegistry()
capabilityRegistry := registry.NewCapabilityRegistry()
pluginBlocker := registry.NewCapabilityChildBlockerChecker(capabilityRegistry)

pluginHandler := api.NewPluginHandler(pluginRegistry, serviceClassRegistry, pluginBlocker, emitter)
capabilityHandler := api.NewCapabilityHandler(capabilityRegistry, pluginRegistry, serviceClassRegistry, emitter)

srv := server.New(cfg, orgHandler, ouHandler, tenantHandler, projectHandler,
    operationHandler, serviceClassHandler, servicePlanHandler,
    pluginHandler, capabilityHandler, bootstrapHandler, readiness)
```

Reuses the single `emitter` and `serviceClassRegistry` created earlier; no duplicate registries.

---

## Error Handling

### Error Mapping Table

| Condition | HTTP | error.code | error.field |
|---|---|---|---|
| ErrAlreadyExists (create) | 409 | RESOURCE_ALREADY_EXISTS | — |
| ErrNotFound (get/update/delete) | 404 | RESOURCE_NOT_FOUND | — |
| []FieldError non-empty | 400 | VALIDATION_FAILED | first field |
| Invalid pluginType | 400 | VALIDATION_FAILED | spec.pluginType |
| Invalid deploymentMode | 400 | VALIDATION_FAILED | spec.deploymentMode |
| Invalid operation (capability) | 400 | VALIDATION_FAILED | spec.operation |
| Empty/invalid serviceClassRefs | 400 | VALIDATION_FAILED | spec.serviceClassRefs |
| Empty/invalid pluginRef | 400 | VALIDATION_FAILED | spec.pluginRef |
| Empty/invalid serviceClassRef | 400 | VALIDATION_FAILED | spec.serviceClassRef |
| Missing version | 400 | VALIDATION_FAILED | spec.version |
| JSON syntax/type error | 400 | VALIDATION_FAILED | — |
| Body exceeds 1 MiB | 413 | VALIDATION_FAILED | — |
| Content-Type mismatch | 415 | VALIDATION_FAILED | — (middleware) |
| Status key in body | 400 | VALIDATION_FAILED | status |
| Referenced ServiceClass not found (plugin create/update) | 400 | VALIDATION_FAILED | spec.serviceClassRefs |
| Referenced Plugin not found (capability create) | 400 | VALIDATION_FAILED | spec.pluginRef |
| Referenced ServiceClass not found (capability create) | 400 | VALIDATION_FAILED | spec.serviceClassRef |
| ServiceClass not declared by Plugin (capability create) | 400 | VALIDATION_FAILED | spec.serviceClassRef |
| metadata.name absent/mismatch (plugin PUT) | 400 | VALIDATION_FAILED | metadata.name |
| Plugin has Capabilities (delete) | 409 | DELETE_BLOCKED | — |
| Capability PUT attempt | 405 | METHOD_NOT_ALLOWED | — |
| Wrong item path segment count | 404 | RESOURCE_NOT_FOUND | — |
| Invalid path segment (not DNS-label) | 400 | VALIDATION_FAILED | metadata.name |
| Blocker/registry unexpected error | 500 | INTERNAL_ERROR | — |

Emission-path errors are NEVER mapped to client responses; `emitOperation` swallows them.

---

## Security and Privacy Constraints

- Plugin and Capability store only registry metadata; NO secrets, tokens, credentials, or passwords.
- `spec.serviceClassRefs`, `spec.pluginRef`, `spec.serviceClassRef` are name references only — no
  credential material.
- `spec.description` fields SHALL NOT be used to store credential-bearing content; no format
  validation is enforced in Phase 1, but documentation warns against it.
- The server SHALL NOT store raw request bodies and SHALL NOT echo raw bodies in error responses.
- The server SHALL NOT log secrets or raw request bodies.
- Operation records carry only non-sensitive reference fields (type, kind, names, request ID).

---

## Testing Strategy

### Validation tests
- `plugin_test.go`: valid names accepted; empty/invalid/long `metadata.name` rejected; invalid
  pluginType rejected; empty version rejected; empty/nil serviceClassRefs rejected; invalid entries
  in serviceClassRefs rejected; invalid deploymentMode rejected; path segment validation.
- `capability_test.go`: valid inputs accepted; empty/invalid `metadata.name` rejected;
  empty/invalid `spec.pluginRef` rejected; empty/invalid `spec.serviceClassRef` rejected; invalid
  operation rejected; path segment validation.

### Property tests (testing/quick, Config{MaxCount: 100})
- `plugin_property_test.go`:
  - Property 1: valid DNS-label names with valid enum values accepted.
  - Property 2: arbitrary invalid strings rejected for name.
  - Property 3: valid pluginType values accepted; invalid values rejected.
  - Property 4: valid deploymentMode values accepted; invalid values rejected.
- `capability_property_test.go`:
  - Property 5: valid DNS-label names with valid enum operation accepted.
  - Property 6: arbitrary invalid strings rejected for name/pluginRef/serviceClassRef.
  - Property 7: valid operation values accepted; invalid values rejected.

Each test tagged `// Feature: plugin-capability-registry, Property N: <title>`.

### Registry tests
- `plugin_registry_test.go`: Create stores; duplicate → ErrAlreadyExists (original unchanged);
  Get by key; missing → ErrNotFound; List sorted; empty → non-nil `[]`; Update mutable fields only;
  Update preserves immutable fields; Update missing → ErrNotFound; Delete removes; Delete missing →
  ErrNotFound; deep-copy immutability (mutating returned value does not change stored state).
- `capability_registry_test.go`: Create stores; duplicate → ErrAlreadyExists (original unchanged);
  Get by key; missing → ErrNotFound; List sorted; empty → non-nil `[]`; Delete removes; Delete
  missing → ErrNotFound; CountByPlugin correct; List with pluginRef filter; List with
  serviceClassRef filter; List with both filters (AND); List with no filters (all); deep-copy
  immutability.

### Registry property tests
- `plugin_registry_property_test.go`: Create/Get round-trip preserves data; List sort invariant;
  deep-copy immutability; duplicate-create idempotent error.
- `capability_registry_property_test.go`: Create/Get round-trip preserves data; List sort
  invariant; deep-copy immutability; duplicate-create idempotent error; filter correctness.

### Race tests
- `plugin_registry_race_test.go`: 10+ goroutines mixed CRUD, zero race reports.
- `capability_registry_race_test.go`: 10+ goroutines mixed Create/Get/List/Delete/CountByPlugin,
  zero race reports.

### Blocker tests
- `plugin_blocker_test.go`: `CountByPlugin` correct; `BlockedByPluginChildren` returns a
  `Capability` blocker when count > 0; returns nil when 0; registry error propagates.

### Handler tests
- `plugin_handler_test.go`: POST 201/409/400 (invalid fields, status key, bad JSON, unknown field);
  POST 413 oversized; POST 400 missing ServiceClass ref; GET 200/404/400; wrong path shape → 404;
  list sorted/empty; PUT 200/404/400 (identity mismatch, invalid fields, missing ServiceClass ref);
  DELETE 204/404; DELETE 409 DELETE_BLOCKED with capabilities present; DELETE 204 with zero
  capabilities; nil emitter/blocker no-panic.
- `capability_handler_test.go`: POST 201/409/400 (invalid fields, status key, bad JSON, unknown
  field); POST 413 oversized; POST 400 missing Plugin ref; POST 400 missing ServiceClass ref;
  POST 400 ServiceClass not in Plugin serviceClassRefs; GET 200/404/400; wrong path shape → 404;
  list sorted/empty; list with query filters (pluginRef, serviceClassRef, both, neither);
  DELETE 204/404; PUT → 405 METHOD_NOT_ALLOWED; nil emitter no-panic.

### Delete-blocking integration tests
- Plugin delete with Capabilities → 409 DELETE_BLOCKED.
- Plugin delete with zero Capabilities → 204.

### Operation emission tests
- `plugin_emission_test.go`: Table-driven across create/update/delete: correct Operation type,
  resourceKind, pluginName recorded; failed actions emit nothing; emission failure does not change
  the primary response.
- `capability_emission_test.go`: Table-driven across create/delete: correct Operation type,
  resourceKind, pluginName, capabilityName recorded; failed actions emit nothing; emission failure
  does not change the primary response.

### Server tests
- `server_test.go` (modified): constructor with Plugin/Capability handler fixtures; route
  registration for all four patterns (`/v1/plugins`, `/v1/plugins/`, `/v1/capabilities`,
  `/v1/capabilities/`).

---

## Verification

```bash
make fmt && make vet && make test && go test -race ./... && make build
```

If host Go is unavailable, use the Docker fallback:

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c \
  'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./... && go build ./cmd/sovrunn-api'
```

---

## Non-Goals

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
13. NO list pagination.
14. NO `operation` filter on Capability list (deferred).

---

## Edge Cases (Design Clarifications)

| # | Scenario | Behavior |
|---|---|---|
| 1 | Mutating action fails (validation, duplicate, missing ref, not-found, delete-blocked) | No Operation emitted |
| 2 | Plugin item path has more than one segment | HTTP 404 |
| 3 | Capability item path has more than one segment | HTTP 404 |
| 4 | `spec.serviceClassRefs` contains duplicate ServiceClass names | Accepted (no deduplication enforced) |
| 5 | Plugin's referenced ServiceClass later deleted | Plugin remains valid (write-time check only) |
| 6 | Capability's `spec.pluginRef` references a Plugin | Plugin cannot be deleted while Capability exists (blocker enforces) |
| 7 | Nil emitter (isolated handler tests) | Emission skipped gracefully without panic |
| 8 | `PUT /v1/capabilities/{name}` | HTTP 405 METHOD_NOT_ALLOWED regardless of existence |
| 9 | Same `metadata.name` used for both a Plugin and a Capability | Allowed (separate registries) |
| 10 | `spec.supported = false` on Capability | Accepted and stored (valid registry metadata) |

---

## Resolved Design Questions Summary

All design questions from the requirements are resolved above. No open questions remain before
implementation.
