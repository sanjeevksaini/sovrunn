# Design Document — FEATURE-0004 Project Resource and Registry

## Overview

FEATURE-0004 implements `Project` as the workload/environment grouping boundary under a Tenant.
Hierarchy: Organization → OrganizationUnit → Tenant → **Project** → ServiceInstance (future).

**Key decisions:**
- Composite identity: `organizationName/organizationUnitName/tenantName/name`
- `spec.organizationName`, `spec.organizationUnitName`, `spec.tenantName` immutable after create
- ProjectRegistry is storage-only — no dependency on other registries
- Parent Tenant existence check in API layer via `TenantLookup` interface
- Tenant delete extended with Project child blocker (intentional FEATURE-0004 change)
- Go 1.21 routing: `/v1/projects` and `/v1/projects/` patterns
- `http.MaxBytesReader` inside `safeDecodeProject`, not middleware
- Middleware order: requestID → logging → contentType → handler
- `CreateProject` and `UpdateProject` return `(resources.Project, error)`

**Scope boundary:** Project only. No ServiceInstance, persistence, CRDs, auth, ServiceOps, UI, AI, SDE.

---

## Resource Model

### internal/resources/project.go

```go
package resources

type Project struct {
    APIVersion string        `json:"apiVersion"`
    Kind       string        `json:"kind"`
    Metadata   Metadata      `json:"metadata"`
    Spec       ProjectSpec   `json:"spec"`
    Status     ProjectStatus `json:"status"`
}

type ProjectSpec struct {
    OrganizationName     string `json:"organizationName"`
    OrganizationUnitName string `json:"organizationUnitName"`
    TenantName           string `json:"tenantName"`
    Description          string `json:"description,omitempty"`
}

type ProjectStatus struct {
    Phase   string `json:"phase"`
    Message string `json:"message,omitempty"`
}

const (
    ProjectAPIVersion = "platform.sovrunn.io/v1alpha1"
    ProjectKind       = "Project"
)
```

Reuses existing `Metadata` struct and `PhaseActive`/`PhaseInactive`/`PhaseDeleting`/`PhaseFailed` constants.

---

## Files to Create/Modify

### New Files

| File | Purpose |
|---|---|
| `internal/resources/project.go` | Project, ProjectSpec, ProjectStatus, constants |
| `internal/registry/project_registry.go` | ProjectRegistryIface, ProjectRegistry, deepCopyProject, TenantLookup |
| `internal/registry/project_blocker.go` | ProjectChildBlockerChecker, TenantChildBlocker interface |
| `internal/validation/project.go` | ValidateProject, ValidateProjectPathSegments |
| `internal/api/project_handler.go` | ProjectHandler, HandleCollection, HandleItem |
| `internal/api/project_decode.go` | safeDecodeProject |
| `internal/validation/project_test.go` | Validation unit tests |
| `internal/validation/project_property_test.go` | Validation property tests |
| `internal/registry/project_registry_test.go` | Registry unit tests |
| `internal/registry/project_registry_property_test.go` | Registry property tests |
| `internal/registry/project_registry_race_test.go` | Concurrency stress test |
| `internal/registry/project_blocker_test.go` | Blocker unit tests |
| `internal/api/project_handler_test.go` | Handler tests |

### Modified Files

| File | Change |
|---|---|
| `internal/api/tenant_handler.go` | Add nil-safe `TenantChildBlocker` dependency for Tenant delete |
| `internal/api/tenant_handler_test.go` | Update NewTenantHandler call sites; add blocker tests |
| `internal/server/server.go` | Accept `*api.ProjectHandler`, register `/v1/projects` routes |
| `internal/server/server_test.go` | Update server constructor tests; add route tests |
| `cmd/sovrunn-api/main.go` | Wire ProjectRegistry, ProjectBlocker, ProjectHandler; inject blocker into TenantHandler |

---

## ProjectRegistry Design

### Interface

```go
// internal/registry/project_registry.go

type ProjectRegistryIface interface {
    CreateProject(ctx context.Context, p resources.Project) (resources.Project, error)
    GetProject(ctx context.Context, orgName, ouName, tenantName, name string) (resources.Project, error)
    ListProjects(ctx context.Context) ([]resources.Project, error)
    UpdateProject(ctx context.Context, p resources.Project) (resources.Project, error)
    DeleteProject(ctx context.Context, orgName, ouName, tenantName, name string) error
    CountByTenant(ctx context.Context, orgName, ouName, tenantName string) (int, error)
}
```

### Concrete Implementation

```go
type ProjectRegistry struct {
    mu    sync.RWMutex
    store map[string]resources.Project // key: "orgName/ouName/tenantName/name"
}

func NewProjectRegistry() *ProjectRegistry {
    return &ProjectRegistry{store: make(map[string]resources.Project)}
}

func projectCompositeKey(orgName, ouName, tenantName, name string) string {
    return orgName + "/" + ouName + "/" + tenantName + "/" + name
}
```

### deepCopyProject

Duplicates `Metadata.Labels` and `Metadata.Annotations` maps so callers cannot mutate stored state.
ProjectSpec has no slices, so no slice copy needed.

### Method Behavior

| Method | Lock | Returns | Error |
|---|---|---|---|
| CreateProject | Lock | deep copy of stored | ErrAlreadyExists if key exists |
| GetProject | RLock | deep copy | ErrNotFound if absent |
| ListProjects | RLock | sorted slice of deep copies | nil |
| UpdateProject | Lock | deep copy of updated | ErrNotFound if absent |
| DeleteProject | Lock | — | ErrNotFound if absent |
| CountByTenant | RLock | int count | nil |

**List sort order:** `Spec.OrganizationName` → `Spec.OrganizationUnitName` → `Spec.TenantName` → `Metadata.Name`, all ascending.

**UpdateProject behavior:** derives composite key from the submitted resource
(`p.Spec.OrganizationName/p.Spec.OrganizationUnitName/p.Spec.TenantName/p.Metadata.Name`),
looks up the existing stored Project, preserves stored `APIVersion`, `Kind`, `Status`,
`Metadata.Name`, `Spec.OrganizationName`, `Spec.OrganizationUnitName`, `Spec.TenantName`,
and replaces only `Metadata.DisplayName`, `Metadata.Labels`, `Metadata.Annotations`,
`Spec.Description`. Returns a deep copy of the updated stored Project.

**CountByTenant:** RLock; iterates entries; counts those where the three parent spec fields match.

Reuses sentinel errors `ErrNotFound` and `ErrAlreadyExists` from `internal/registry/registry.go`.

---

## Parent Tenant Lookup

```go
// internal/registry/project_registry.go (or shared interfaces file)

type TenantLookup interface {
    GetTenant(ctx context.Context, orgName, ouName, name string) (resources.Tenant, error)
}
```

The existing `*TenantRegistry` satisfies this interface via duck typing.
Before `CreateProject`, ProjectHandler calls `tenantLookup.GetTenant(ctx, orgName, ouName, tenantName)`.
If `ErrNotFound`, return HTTP 400, `VALIDATION_FAILED`, `field = "spec.tenantName"`,
message identifying `"spec.organizationName/spec.organizationUnitName/spec.tenantName"`.

---

## Validation Design

### internal/validation/project.go

```go
// ValidateProject validates all user-authored fields. Context-free, no I/O.
func ValidateProject(p resources.Project) []resources.FieldError

// ValidateProjectPathSegments validates four URL path segments.
func ValidateProjectPathSegments(orgName, ouName, tenantName, name string) []resources.FieldError
```

**Rules (all DNS-label, 1–63 chars):**
- `metadata.name` → `error.field = "metadata.name"`
- `spec.organizationName` → `error.field = "spec.organizationName"`
- `spec.organizationUnitName` → `error.field = "spec.organizationUnitName"`
- `spec.tenantName` → `error.field = "spec.tenantName"`

**Path validation field mapping:**
- invalid orgName segment → `spec.organizationName`
- invalid ouName segment → `spec.organizationUnitName`
- invalid tenantName segment → `spec.tenantName`
- invalid name segment → `metadata.name`

Reuses `dnsLabelRe` and `validateName` from the validation package (same package, unexported).

---

## Safe JSON Decoding

### internal/api/project_decode.go

```go
func safeDecodeProject(w http.ResponseWriter, r *http.Request) (resources.Project, error)
```

**Sequence:**
1. `r.Body = http.MaxBytesReader(w, r.Body, 1<<20)` — inside this function
2. Read body bytes; `*http.MaxBytesError` → errBodyTooLarge (413)
3. Empty body → errEmptyBody (400)
4. Decode into `map[string]json.RawMessage`; if "status" key present → errStatusFieldPresent (400)
5. Typed decode with `DisallowUnknownFields()` into `resources.Project`
6. Unknown field → errUnknownField (400); syntax/type → errMalformedJSON (400)
7. Return decoded Project

Do not echo raw body. HTTP 415 is handled by `contentTypeMiddleware`, not here.
Reuses existing error sentinels from the FEATURE-0003 decoder pattern.

---

## Project HTTP Handler

### internal/api/project_handler.go

```go
type ProjectHandler struct {
    registry     registry.ProjectRegistryIface
    tenantLookup registry.TenantLookup
}

func NewProjectHandler(
    reg registry.ProjectRegistryIface,
    tenantLookup registry.TenantLookup,
) *ProjectHandler
```

### HandleCollection (POST/GET)

`POST → Create`, `GET → List`, else 405.

### HandleItem (GET/PUT/DELETE)

```go
// remainder := strings.TrimPrefix(r.URL.Path, "/v1/projects/")
// parts := strings.Split(remainder, "/")
// if len(parts) != 4 || any part == "" → 404
// orgName, ouName, tenantName, name := parts[0], parts[1], parts[2], parts[3]
// dispatch GET → Get, PUT → Update, DELETE → Delete, else 405
```

### Create Flow

```
1. safeDecodeProject (MaxBytesReader, status detection, DisallowUnknownFields)
2. If decode error → writeError(400/413)
3. ValidateProject(project) → []FieldError
4. If errors → writeValidationErrors
5. tenantLookup.GetTenant(ctx, orgName, ouName, tenantName)
6. If ErrNotFound → writeError(400, VALIDATION_FAILED, field="spec.tenantName",
     message="parent not found: orgName/ouName/tenantName")
7. Force: apiVersion, kind, status.phase = Active
8. created, err := registry.CreateProject(ctx, project)
9. If ErrAlreadyExists → writeError(409, RESOURCE_ALREADY_EXISTS)
10. // TODO(FEATURE-0005): emit Operation record — type: CreateProject
11. writeJSON(201, created)
```

### Get Flow

```
1. ValidateProjectPathSegments(orgName, ouName, tenantName, name)
2. If errors → writeValidationErrors
3. project, err := registry.GetProject(ctx, orgName, ouName, tenantName, name)
4. If ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
5. writeJSON(200, project)
```

### List Flow

```
1. items, err := registry.ListProjects(ctx)
2. If error → writeError(500, INTERNAL_ERROR)
3. writeJSON(200, {"items": items})  // [] when empty
```

### Update Flow

```
1. ValidateProjectPathSegments(orgName, ouName, tenantName, name)
2. If errors → writeValidationErrors
3. safeDecodeProject(w, r)
4. If decode error → writeError(400/413)
5. body.Metadata.Name present and == name, else 400 (field="metadata.name")
6. body.Spec.OrganizationName present and == orgName, else 400 (field="spec.organizationName")
7. body.Spec.OrganizationUnitName present and == ouName, else 400 (field="spec.organizationUnitName")
8. body.Spec.TenantName present and == tenantName, else 400 (field="spec.tenantName")
9. ValidateProject(project)
10. If errors → writeValidationErrors
11. updated, err := registry.UpdateProject(ctx, project)
12. If ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
13. // TODO(FEATURE-0005): emit Operation record — type: UpdateProject
14. writeJSON(200, updated)
```

### Delete Flow

```
1. ValidateProjectPathSegments(orgName, ouName, tenantName, name)
2. If errors → writeValidationErrors
3. err := registry.DeleteProject(ctx, orgName, ouName, tenantName, name)
4. If ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
5. // TODO(FEATURE-0005): emit Operation record — type: DeleteProject
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
| Parent Tenant not found | 400 | VALIDATION_FAILED | spec.tenantName |
| metadata.name absent/mismatch | 400 | VALIDATION_FAILED | metadata.name |
| spec.organizationName absent/mismatch | 400 | VALIDATION_FAILED | spec.organizationName |
| spec.organizationUnitName absent/mismatch | 400 | VALIDATION_FAILED | spec.organizationUnitName |
| spec.tenantName absent/mismatch | 400 | VALIDATION_FAILED | spec.tenantName |
| Any unexpected error | 500 | INTERNAL_ERROR | — |

---

## Tenant Deletion Blocker Integration

### internal/registry/project_blocker.go

```go
type TenantChildBlocker interface {
    BlockedByTenantChildren(ctx context.Context, orgName, ouName, tenantName string) ([]BlockedBy, error)
}

type ProjectChildBlockerChecker struct {
    projectRegistry ProjectRegistryIface
}

func NewProjectChildBlockerChecker(reg ProjectRegistryIface) *ProjectChildBlockerChecker

func (c *ProjectChildBlockerChecker) BlockedByTenantChildren(
    ctx context.Context, orgName, ouName, tenantName string,
) ([]BlockedBy, error) {
    count, err := c.projectRegistry.CountByTenant(ctx, orgName, ouName, tenantName)
    if err != nil { return nil, err }
    if count > 0 {
        return []BlockedBy{{Kind: "Project", Count: count}}, nil
    }
    return nil, nil
}
```

### TenantHandler Modification

```go
// internal/api/tenant_handler.go — FEATURE-0004 change

type TenantHandler struct {
    registry registry.TenantRegistryIface
    ouLookup registry.OrganizationUnitLookup
    blocker  registry.TenantChildBlocker  // NEW: FEATURE-0004, may be nil
}

func NewTenantHandler(
    reg registry.TenantRegistryIface,
    ouLookup registry.OrganizationUnitLookup,
    blocker registry.TenantChildBlocker,  // NEW: FEATURE-0004
) *TenantHandler
```

**Nil blocker behavior:** TenantHandler SHALL allow `blocker` to be nil. If nil, Tenant delete
proceeds without child-blocker checks. Production FEATURE-0004 wiring injects
`ProjectChildBlockerChecker`. This preserves FEATURE-0003 test compatibility.

### Updated TenantHandler.Delete

```
1. ValidateTenantPathSegments(orgName, ouName, name)
2. If errors → writeValidationErrors
3. If blocker != nil:
     blockers, err := blocker.BlockedByTenantChildren(ctx, orgName, ouName, name)
     If blockers non-empty → writeError(409, DELETE_BLOCKED, message "Project")
4. registry.DeleteTenant(ctx, orgName, ouName, name)
5. If ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
6. // TODO(FEATURE-0005): emit Operation record — type: DeleteTenant
7. w.WriteHeader(204)
```

Do not create a generic blocker framework. Do not rewrite OUHandler or OrgHandler.

---

## Server and main.go Wiring

### Updated server.New signature

```go
func New(
    cfg config.Config,
    org *api.OrgHandler,
    ou *api.OUHandler,
    tenant *api.TenantHandler,
    project *api.ProjectHandler,  // NEW
    bootstrap *api.BootstrapHandler,
    readiness *health.ReadinessState,
) *Server
```

`server.New` SHALL require a non-nil `ProjectHandler` in FEATURE-0004 production wiring and tests.

### Route registration (Go 1.21)

```go
mux.Handle("/v1/projects", chain(http.HandlerFunc(project.HandleCollection)))
mux.Handle("/v1/projects/", chain(http.HandlerFunc(project.HandleItem)))
```

Middleware order unchanged: requestID → logging → contentType → handler.

### Updated main.go wiring

```go
projectRegistry := registry.NewProjectRegistry()
projectBlocker := registry.NewProjectChildBlockerChecker(projectRegistry)

// TenantHandler now receives project blocker
tenantHandler := api.NewTenantHandler(tenantRegistry, ouRegistry, projectBlocker)

// ProjectHandler receives project registry + tenant lookup
projectHandler := api.NewProjectHandler(projectRegistry, tenantRegistry)

srv := server.New(cfg, orgHandler, ouHandler, tenantHandler, projectHandler, bootstrapHandler, readiness)
```

Update any existing NewTenantHandler call sites for the new signature.

---

## Test Design

### Validation Tests

| File | Coverage |
|---|---|
| `project_test.go` | valid names; empty/invalid/long metadata.name, orgName, ouName, tenantName rejected; path field mapping |
| `project_property_test.go` | P1 valid DNS-labels accepted; P2 invalid strings rejected |

### Registry Tests

| File | Coverage |
|---|---|
| `project_registry_test.go` | Create, Get, List sorted, Update preserves immutables, Delete, CountByTenant, duplicate → ErrAlreadyExists, same name different Tenants, not found errors |
| `project_registry_property_test.go` | P3 round-trip; P4 four-level sort invariant; P5 deep-copy immutability; P6 idempotent duplicate error |
| `project_registry_race_test.go` | 10+ goroutines mixed CRUD + CountByTenant, zero race reports |

### Blocker Tests

| File | Coverage |
|---|---|
| `project_blocker_test.go` | CountByTenant correct; BlockedByTenantChildren returns Project blocker when count > 0; returns nil when 0; registry error propagates |

### Handler Tests

| File | Coverage |
|---|---|
| `project_handler_test.go` | POST 201/409/400(invalid fields, missing parent, status key, bad JSON, oversized body); GET 200/404/400; GET wrong path shape → 404; list sorted/empty; PUT 200/404/400(name/orgName/ouName/tenantName mismatch, status); DELETE 204/404/400 |

### Tenant Blocker Integration Tests

| File | Coverage |
|---|---|
| `tenant_handler_test.go` (modified) | Tenant delete with Projects → 409 DELETE_BLOCKED ("Project"); Tenant delete with zero Projects → 204; nil blocker preserves FEATURE-0003 behavior |

### Server Tests

| File | Coverage |
|---|---|
| `server_test.go` (modified) | server.New with ProjectHandler fixture; `/v1/projects` and `/v1/projects/` route registration |

### Property Test Tags

Each property test tagged: `// Feature: project-resource, Property N: <title>`

| N | Title |
|---|---|
| 1 | ValidateProject accepts valid DNS-label names |
| 2 | ValidateProject rejects invalid names |
| 3 | Create/Get round trip preserves data |
| 4 | List returns projects in correct four-level composite sort order |
| 5 | Registry returns deep copies (mutations don't affect stored state) |
| 6 | Duplicate create returns ErrAlreadyExists and original unchanged |

---

## Non-Goals

1. ServiceInstance resource or nested Project hierarchies
2. Persistent storage, Kubernetes CRDs
3. Authentication, authorization, RBAC
4. Operation framework implementation (beyond TODO comments)
5. ServiceOps, plugin execution, AI agent execution
6. UI, SDE runtime transformation, billing
7. Project delete child-resource blockers (future ServiceInstance blockers)
8. Generic blocker framework
9. List filtering or pagination
10. Go 1.22 wildcard routing

---

## Design Questions (Resolved)

| # | Question | Resolution |
|---|---|---|
| 1 | Parent lookup interface | Use `TenantLookup` interface; existing `*TenantRegistry` satisfies it |
| 2 | Project blocker injection into TenantHandler | Add nil-safe `TenantChildBlocker` to `TenantHandler` via updated `NewTenantHandler`; inject `ProjectChildBlockerChecker` |
| 3 | Four-segment path parsing | `strings.Split` with `len(parts) == 4` check (resolved by Req 5.6) |
