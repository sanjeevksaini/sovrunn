# Design Document — FEATURE-0003 Tenant Resource and Registry

## Overview

FEATURE-0003 implements `Tenant` as the primary isolation and security boundary in Sovrunn.
Hierarchy: Organization → OrganizationUnit → **Tenant** → Project (future).

**Key decisions:**
- Composite identity: `organizationName/organizationUnitName/name`
- `spec.organizationName` and `spec.organizationUnitName` immutable after create
- TenantRegistry is storage-only — no dependency on other registries
- Parent OU existence check in API layer via `OrganizationUnitLookup` interface
- OU delete extended with Tenant child blocker (intentional FEATURE-0003 change)
- Go 1.21 routing: `/v1/tenants` and `/v1/tenants/` patterns
- `http.MaxBytesReader` inside `safeDecodeTenant`, not middleware
- Middleware order: requestID → logging → contentType → handler
- `UpdateTenant` returns `(resources.Tenant, error)`
- `CreateTenant` returns `(resources.Tenant, error)`

**Scope boundary:** Tenant only. No Project, persistence, CRDs, auth, ServiceOps, UI, AI, SDE.

---

## Resource Model

### internal/resources/tenant.go

```go
package resources

type Tenant struct {
    APIVersion string       `json:"apiVersion"`
    Kind       string       `json:"kind"`
    Metadata   Metadata     `json:"metadata"`
    Spec       TenantSpec   `json:"spec"`
    Status     TenantStatus `json:"status"`
}

type TenantSpec struct {
    OrganizationName     string `json:"organizationName"`
    OrganizationUnitName string `json:"organizationUnitName"`
    Description          string `json:"description,omitempty"`
}

type TenantStatus struct {
    Phase   string `json:"phase"`
    Message string `json:"message,omitempty"`
}

const (
    TenantAPIVersion = "platform.sovrunn.io/v1alpha1"
    TenantKind       = "Tenant"
)
```

Reuses existing `Metadata` struct and `PhaseActive`/`PhaseInactive`/`PhaseDeleting`/`PhaseFailed` constants.

---

## Files to Create/Modify

### New Files

| File | Purpose |
|---|---|
| `internal/resources/tenant.go` | Tenant, TenantSpec, TenantStatus, constants |
| `internal/registry/tenant_registry.go` | TenantRegistryIface, TenantRegistry, deepCopyTenant, compositeKey |
| `internal/registry/tenant_blocker.go` | TenantChildBlockerChecker for OU delete |
| `internal/validation/tenant.go` | ValidateTenant, ValidateTenantPathSegments |
| `internal/api/tenant_handler.go` | TenantHandler, HandleCollection, HandleItem |
| `internal/api/tenant_decode.go` | safeDecodeTenant |
| `internal/validation/tenant_test.go` | Validation unit tests |
| `internal/validation/tenant_property_test.go` | Validation property tests |
| `internal/registry/tenant_registry_test.go` | Registry unit tests |
| `internal/registry/tenant_registry_property_test.go` | Registry property tests |
| `internal/registry/tenant_registry_race_test.go` | Concurrency stress test |
| `internal/registry/tenant_blocker_test.go` | Blocker unit tests |
| `internal/api/tenant_handler_test.go` | Handler tests |

### Modified Files

| File | Change |
|---|---|
| `internal/api/ou_handler.go` | Add optional `OUChildBlocker` dependency for OU delete |
| `internal/server/server.go` | Accept `*api.TenantHandler`, register `/v1/tenants` routes |
| `internal/server/server_test.go` | Update server constructor tests and route registration tests for TenantHandler |
| `cmd/sovrunn-api/main.go` | Wire TenantRegistry, TenantBlocker, TenantHandler; inject blocker into OUHandler |

---

## TenantRegistry Design

### Interface

```go
// internal/registry/tenant_registry.go

type TenantRegistryIface interface {
    CreateTenant(ctx context.Context, t resources.Tenant) (resources.Tenant, error)
    GetTenant(ctx context.Context, orgName, ouName, name string) (resources.Tenant, error)
    ListTenants(ctx context.Context) ([]resources.Tenant, error)
    UpdateTenant(ctx context.Context, t resources.Tenant) (resources.Tenant, error)
    DeleteTenant(ctx context.Context, orgName, ouName, name string) error
    CountByOrganizationUnit(ctx context.Context, orgName, ouName string) (int, error)
}
```

### Concrete Implementation

```go
type TenantRegistry struct {
    mu    sync.RWMutex
    store map[string]resources.Tenant // key: "orgName/ouName/name"
}

func NewTenantRegistry() *TenantRegistry {
    return &TenantRegistry{store: make(map[string]resources.Tenant)}
}

func tenantCompositeKey(orgName, ouName, name string) string {
    return orgName + "/" + ouName + "/" + name
}
```

### deepCopyTenant

```go
func deepCopyTenant(t resources.Tenant) resources.Tenant {
    cp := t
    if t.Metadata.Labels != nil {
        cp.Metadata.Labels = make(map[string]string, len(t.Metadata.Labels))
        for k, v := range t.Metadata.Labels { cp.Metadata.Labels[k] = v }
    }
    if t.Metadata.Annotations != nil {
        cp.Metadata.Annotations = make(map[string]string, len(t.Metadata.Annotations))
        for k, v := range t.Metadata.Annotations { cp.Metadata.Annotations[k] = v }
    }
    return cp
}
```

### Method Behavior

| Method | Lock | Returns | Error |
|---|---|---|---|
| CreateTenant | Lock | deep copy of stored | ErrAlreadyExists if key exists |
| GetTenant | RLock | deep copy | ErrNotFound if absent |
| ListTenants | RLock | sorted slice of deep copies | nil |
| UpdateTenant | Lock | deep copy of updated | ErrNotFound if absent |
| DeleteTenant | Lock | — | ErrNotFound if absent |
| CountByOrganizationUnit | RLock | int count | nil |

**List sort order:** `Spec.OrganizationName` asc → `Spec.OrganizationUnitName` asc → `Metadata.Name` asc.

**UpdateTenant behavior:** `UpdateTenant` SHALL look up the existing stored Tenant by the composite key derived from the submitted resource (`t.Spec.OrganizationName/t.Spec.OrganizationUnitName/t.Metadata.Name`). It SHALL preserve stored `APIVersion`, `Kind`, `Status`, `Metadata.Name`, `Spec.OrganizationName`, and `Spec.OrganizationUnitName`. It SHALL replace only `Metadata.DisplayName`, `Metadata.Labels`, `Metadata.Annotations`, and `Spec.Description`. This prevents client input from overwriting server-owned or immutable fields.

Reuses sentinel errors `ErrNotFound` and `ErrAlreadyExists` from `internal/registry/registry.go`.

---

## Parent OrganizationUnit Lookup

```go
// internal/registry/tenant_registry.go (or shared interfaces file)

type OrganizationUnitLookup interface {
    GetOrganizationUnit(ctx context.Context, orgName, name string) (resources.OrganizationUnit, error)
}
```

The existing `*OrganizationUnitRegistry` satisfies this interface via duck typing.
TenantHandler receives `OrganizationUnitLookup` at construction. Before `CreateTenant`,
the handler calls `ouLookup.GetOrganizationUnit(ctx, orgName, ouName)`. If `ErrNotFound`,
return HTTP 400, `VALIDATION_FAILED`, `field = "spec.organizationUnitName"`,
message identifying `"spec.organizationName/spec.organizationUnitName"`.

---

## Validation Design

### internal/validation/tenant.go

```go
// ValidateTenant validates all user-authored fields. Context-free, no I/O.
func ValidateTenant(t resources.Tenant) []resources.FieldError

// ValidateTenantPathSegments validates three URL path segments.
func ValidateTenantPathSegments(orgName, ouName, name string) []resources.FieldError
```

**Rules:**
- `metadata.name`: required, DNS-label, 1–63 → `error.field = "metadata.name"`
- `spec.organizationName`: required, DNS-label, 1–63 → `error.field = "spec.organizationName"`
- `spec.organizationUnitName`: required, DNS-label, 1–63 → `error.field = "spec.organizationUnitName"`

**Path validation field mapping:**
- invalid orgName segment → `error.field = "spec.organizationName"`
- invalid ouName segment → `error.field = "spec.organizationUnitName"`
- invalid name segment → `error.field = "metadata.name"`

Reuses `dnsLabelRe` and `validateName` from validation package (same package, unexported).

---

## Safe JSON Decoding

### internal/api/tenant_decode.go

```go
func safeDecodeTenant(w http.ResponseWriter, r *http.Request) (resources.Tenant, error)
```

**Sequence:**
1. `r.Body = http.MaxBytesReader(w, r.Body, 1<<20)` — inside this function
2. Read body bytes; `*http.MaxBytesError` → errBodyTooLarge (413)
3. Empty body → errEmptyBody (400)
4. Decode into `map[string]json.RawMessage`; if "status" key present → errStatusFieldPresent (400)
5. Typed decode with `DisallowUnknownFields()` into `resources.Tenant`
6. Unknown field → errUnknownField (400); syntax/type → errMalformedJSON (400)
7. Return decoded Tenant

HTTP 415 is handled by `contentTypeMiddleware`, not here.

---

## Tenant HTTP Handler

### internal/api/tenant_handler.go

```go
type TenantHandler struct {
    registry  registry.TenantRegistryIface
    ouLookup  registry.OrganizationUnitLookup
}

func NewTenantHandler(
    reg registry.TenantRegistryIface,
    ouLookup registry.OrganizationUnitLookup,
) *TenantHandler
```

### HandleCollection (POST/GET)

```go
func (h *TenantHandler) HandleCollection(w http.ResponseWriter, r *http.Request)
// POST → Create, GET → List, else 405
```

### HandleItem (GET/PUT/DELETE)

```go
func (h *TenantHandler) HandleItem(w http.ResponseWriter, r *http.Request)
// Path parsing:
//   remainder := strings.TrimPrefix(r.URL.Path, "/v1/tenants/")
//   parts := strings.Split(remainder, "/")
//   if len(parts) != 3 || parts[0]=="" || parts[1]=="" || parts[2]=="" → 404
//   orgName, ouName, name := parts[0], parts[1], parts[2]
// Then dispatch: GET → Get, PUT → Update, DELETE → Delete, else 405
```

### Create Flow

```
1. safeDecodeTenant (MaxBytesReader, status detection, DisallowUnknownFields)
2. If decode error → writeError(400/413)
3. ValidateTenant(tenant) → []FieldError
4. If errors → writeValidationErrors
5. ouLookup.GetOrganizationUnit(ctx, orgName, ouName)
6. If ErrNotFound → writeError(400, VALIDATION_FAILED, field="spec.organizationUnitName",
     message="parent not found: orgName/ouName")
7. Force: apiVersion, kind, status.phase = Active
8. created, err := registry.CreateTenant(ctx, tenant)
9. If ErrAlreadyExists → writeError(409, RESOURCE_ALREADY_EXISTS)
10. // TODO(FEATURE-0005): emit Operation record — type: CreateTenant
11. writeJSON(201, created)
```

### Get Flow

```
1. ValidateTenantPathSegments(orgName, ouName, name)
2. If errors → writeValidationErrors
3. tenant, err := registry.GetTenant(ctx, orgName, ouName, name)
4. If ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
5. writeJSON(200, tenant)
```

### List Flow

```
1. items, err := registry.ListTenants(ctx)
2. If error → writeError(500, INTERNAL_ERROR)
3. writeJSON(200, {"items": items})  // [] when empty
```

### Update Flow

```
1. ValidateTenantPathSegments(orgName, ouName, name)
2. If errors → writeValidationErrors
3. safeDecodeTenant(w, r)
4. If decode error → writeError(400/413)
5. body.Metadata.Name must be present and == name, else 400
6. body.Spec.OrganizationName must be present and == orgName, else 400
7. body.Spec.OrganizationUnitName must be present and == ouName, else 400
8. ValidateTenant(tenant)
9. If errors → writeValidationErrors
10. updated, err := registry.UpdateTenant(ctx, tenant)
11. If ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
12. // TODO(FEATURE-0005): emit Operation record — type: UpdateTenant
13. writeJSON(200, updated)
```

### Delete Flow

```
1. ValidateTenantPathSegments(orgName, ouName, name)
2. If errors → writeValidationErrors
3. err := registry.DeleteTenant(ctx, orgName, ouName, name)
4. If ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
5. // TODO(FEATURE-0005): emit Operation record — type: DeleteTenant
6. w.WriteHeader(204)
```

---

## Error Mapping Table

| Condition | HTTP | error.code | error.field |
|---|---|---|---|
| ErrAlreadyExists | 409 | RESOURCE_ALREADY_EXISTS | — |
| ErrNotFound | 404 | RESOURCE_NOT_FOUND | — |
| []FieldError non-empty | 400 | VALIDATION_FAILED | first field |
| JSON syntax/type error | 400 | VALIDATION_FAILED | — |
| Body exceeds 1 MiB | 413 | VALIDATION_FAILED | — |
| Content-Type mismatch | 415 | VALIDATION_FAILED | — (middleware) |
| Status key in body | 400 | VALIDATION_FAILED | status |
| Parent OU not found | 400 | VALIDATION_FAILED | spec.organizationUnitName |
| metadata.name absent/mismatch | 400 | VALIDATION_FAILED | metadata.name |
| spec.organizationName absent/mismatch | 400 | VALIDATION_FAILED | spec.organizationName |
| spec.organizationUnitName absent/mismatch | 400 | VALIDATION_FAILED | spec.organizationUnitName |
| Any unexpected error | 500 | INTERNAL_ERROR | — |

---

## OrganizationUnit Deletion Blocker Integration

### internal/registry/tenant_blocker.go

```go
type TenantChildBlockerChecker struct {
    tenantRegistry TenantRegistryIface
}

func NewTenantChildBlockerChecker(reg TenantRegistryIface) *TenantChildBlockerChecker

func (c *TenantChildBlockerChecker) BlockedByChildren(
    ctx context.Context, orgName, ouName string,
) ([]BlockedBy, error) {
    count, err := c.tenantRegistry.CountByOrganizationUnit(ctx, orgName, ouName)
    if err != nil { return nil, err }
    if count > 0 {
        return []BlockedBy{{Kind: "Tenant", Count: count}}, nil
    }
    return nil, nil
}
```

### OUHandler Modification

FEATURE-0003 extends OUHandler to accept a child blocker via the `OUChildBlocker` interface:

```go
// internal/registry/tenant_blocker.go (or shared interfaces)

type OUChildBlocker interface {
    BlockedByOUChildren(ctx context.Context, orgName, ouName string) ([]BlockedBy, error)
}
```

`TenantChildBlockerChecker` implements `OUChildBlocker`:

```go
func (c *TenantChildBlockerChecker) BlockedByOUChildren(
    ctx context.Context, orgName, ouName string,
) ([]BlockedBy, error) {
    count, err := c.tenantRegistry.CountByOrganizationUnit(ctx, orgName, ouName)
    if err != nil { return nil, err }
    if count > 0 {
        return []BlockedBy{{Kind: "Tenant", Count: count}}, nil
    }
    return nil, nil
}
```

Updated OUHandler struct:

```go
// internal/api/ou_handler.go — FEATURE-0003 change

type OUHandler struct {
    registry  registry.OrganizationUnitRegistryIface
    orgLookup registry.OrganizationLookup
    blocker   registry.OUChildBlocker  // NEW: added by FEATURE-0003, may be nil
}
```

**Nil blocker behavior:** OUHandler SHALL allow `blocker` to be nil. If `blocker` is nil, OU delete proceeds without child-blocker checks. In production FEATURE-0003 wiring, `main.go` SHALL inject `TenantChildBlockerChecker`. This preserves compatibility with existing FEATURE-0002 tests while enabling Tenant child blocking in production wiring.

### Updated NewOUHandler

```go
func NewOUHandler(
    reg registry.OrganizationUnitRegistryIface,
    orgLookup registry.OrganizationLookup,
    blocker registry.OUChildBlocker,  // NEW: FEATURE-0003
) *OUHandler
```

### Updated OUHandler.Delete

```
1. ValidateOUPathSegments(orgName, name)
2. If errors → writeValidationErrors
3. If blocker != nil:
     blockers, err := blocker.BlockedByOUChildren(ctx, orgName, name)
     If blockers non-empty → writeError(409, DELETE_BLOCKED, message "Tenant")
4. registry.DeleteOrganizationUnit(ctx, orgName, name)
5. If ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
6. // TODO(FEATURE-0005): emit Operation record — type: DeleteOrganizationUnit
7. w.WriteHeader(204)
```

---

## Server and main.go Wiring

### Updated server.New signature

```go
func New(
    cfg config.Config,
    org *api.OrgHandler,
    ou *api.OUHandler,
    tenant *api.TenantHandler,  // NEW
    bootstrap *api.BootstrapHandler,
    readiness *health.ReadinessState,
) *Server
```

### Route registration (Go 1.21)

```go
mux.Handle("/v1/tenants", chain(http.HandlerFunc(tenant.HandleCollection)))
mux.Handle("/v1/tenants/", chain(http.HandlerFunc(tenant.HandleItem)))
```

### Updated main.go wiring

```go
tenantRegistry := registry.NewTenantRegistry()
tenantBlocker := registry.NewTenantChildBlockerChecker(tenantRegistry)

// OUHandler now receives tenant blocker
ouHandler := api.NewOUHandler(ouRegistry, orgRegistry, tenantBlocker)

// TenantHandler receives tenant registry + OU lookup
tenantHandler := api.NewTenantHandler(tenantRegistry, ouRegistry)

srv := server.New(cfg, orgHandler, ouHandler, tenantHandler, bootstrapHandler, readiness)
```

**server.New requirement:** `server.New` SHALL require a non-nil `TenantHandler` in FEATURE-0003 production wiring and tests. Production `main.go` SHALL always create and pass `TenantHandler`. Tests that construct `server.New` must provide a `TenantHandler` test fixture.

---

## Test Design

### Validation Tests

| File | Coverage |
|---|---|
| `tenant_test.go` | valid names, empty/invalid/long names, empty/invalid orgName, empty/invalid ouName |
| `tenant_property_test.go` | P1: valid DNS-labels accepted; P2: invalid strings rejected |

### Registry Tests

| File | Coverage |
|---|---|
| `tenant_registry_test.go` | Create, Get, List sorted, Update preserves immutables (derives key from submitted resource), Delete, CountByOU, duplicate → ErrAlreadyExists, same name different parents succeeds, not found errors |
| `tenant_registry_property_test.go` | P3: round-trip; P4: sort invariant; P5: deep-copy immutability; P6: idempotent duplicate error |
| `tenant_registry_race_test.go` | 10+ goroutines mixed CRUD + CountByOU, zero race reports |

### Blocker Tests

| File | Coverage |
|---|---|
| `tenant_blocker_test.go` | CountByOU returns correct count; BlockedByOUChildren returns Tenant blocker when count > 0; returns nil when 0 |

### Handler Tests

| File | Coverage |
|---|---|
| `tenant_handler_test.go` | POST 201/409/400(invalid name, missing orgName, missing ouName, parent not found, status key, bad JSON, oversized body); GET 200/404/400(invalid segments, wrong count → 404); GET list sorted/empty; PUT 200/404/400(name mismatch, orgName mismatch, ouName mismatch, status); DELETE 204/404/400; OU delete blocked by Tenants → 409; OU delete allowed when 0 Tenants |

### Property Test Tags

Each property test tagged: `// Feature: tenant-resource, Property N: <title>`

| N | Title |
|---|---|
| 1 | ValidateTenant accepts valid DNS-label names |
| 2 | ValidateTenant rejects invalid names |
| 3 | Create/Get round trip preserves data |
| 4 | List returns tenants in correct composite sort order |
| 5 | Registry returns deep copies (mutations don't affect stored state) |
| 6 | Duplicate create returns ErrAlreadyExists and original unchanged |

---

## Non-Goals

1. Project resource or nested Tenant hierarchies
2. Persistent storage, Kubernetes CRDs
3. Authentication, authorization, RBAC
4. Operation framework implementation (beyond TODO comments)
5. ServiceOps, plugin execution, AI agent execution
6. UI, SDE runtime transformation, billing
7. Tenant delete child-resource blockers (future Project blockers)
8. Generic blocker framework
9. List filtering or pagination
10. Go 1.22 wildcard routing

---

## Design Questions (Resolved)

| # | Question | Resolution |
|---|---|---|
| 1 | Parent lookup interface | Use `OrganizationUnitLookup` interface; existing `*OrganizationUnitRegistry` satisfies it |
| 2 | OU delete blocker injection | Add `OUChildBlocker` interface to `OUHandler`; inject `TenantChildBlockerChecker` via updated `NewOUHandler` |
| 3 | Three-segment path parsing | `strings.Split` with `len(parts) == 3` check (resolved by Req 5.6) |
