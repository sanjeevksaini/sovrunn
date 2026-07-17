# Design Document — FEATURE-0005 Operation Resource

## Overview

FEATURE-0005 implements `Operation` as the Phase 1 lifecycle/audit record for API-driven
mutating control-plane actions (create/update/delete on Organization, OrganizationUnit,
Tenant, Project). Operations are emitted **synchronously** after a successful mutating action
and stored in an in-memory registry. This is NOT a workflow engine.

**Resolved design decisions:**

1. **OperationEmitter interface placement** — A narrow `OperationEmitter` interface is defined in
   `internal/api` and consumed by the mutating handlers. Handlers depend on the interface, never
   on the concrete `OperationRegistry`. A thin adapter implements the emitter using the registry.

2. **ID generation location** — Operation IDs are generated in the **emitter layer** using
   `crypto/rand`. The `OperationRegistry` is storage-only and requires `metadata.name` to be set.
   On `ErrAlreadyExists` collision, the emitter retries a bounded number of times (5). Persistent
   collision failure is treated as an emission failure and does NOT affect the primary response.

3. **Phase 1 status phase** — Emitted Operations use `Succeeded` only. `Pending`, `Running`, and
   `Failed` remain reserved in the model for future async workflows. No async behavior is implemented.

**Scope boundary:** Operation record + read API + synchronous emission only. No async engine,
queue, workers, approvals, persistence, auth, ServiceOps, or workflow framework.

**Architectural constraint:** `internal/api` MUST NOT import `internal/server`. The request ID is
read via an API-local helper (see "Request ID helper" below), never by importing the server
package. This prevents an import cycle and keeps the handler layer independent of server wiring.

---

## Architecture

The Operation feature spans three layers, consistent with existing resources:

```
cmd/sovrunn-api  →  internal/server  →  internal/api  →  internal/registry
                                             │                  │
                                             └─ OperationEmitter (adapter) ─┘
                                        internal/resources (Operation model)
```

- **internal/resources** — the Operation data model (storage-agnostic value types).
- **internal/registry** — `OperationRegistry` (storage-only) + `registryEmitter` adapter.
- **internal/api** — `OperationHandler` (read API) + `OperationEmitter` interface consumed by the
  mutating handlers. `internal/api` MUST NOT import `internal/server`.
- **internal/server** — route registration and middleware chain (requestID → logging → contentType).

Data flow for emission: a mutating handler completes its registry action, then calls the nil-safe
`emitOperation` helper, which delegates to the injected `OperationEmitter`; the adapter generates an
ID and stores the Operation in the `OperationRegistry`. Read flow: `OperationHandler` reads directly
from the `OperationRegistry`.

---

## Files to Create/Modify

### New Files

| File | Purpose |
|---|---|
| `internal/resources/operation.go` | Operation, OperationSpec, OperationStatus, type & phase constants |
| `internal/registry/operation_registry.go` | OperationRegistryIface, OperationRegistry, deepCopyOperation |
| `internal/api/operation_emitter.go` | OperationEmitter interface, registryEmitter adapter, ID generator |
| `internal/api/operation_handler.go` | OperationHandler, HandleCollection, HandleItem |
| `internal/registry/operation_registry_test.go` | Registry unit tests |
| `internal/registry/operation_registry_race_test.go` | Concurrency stress test |
| `internal/api/operation_emitter_test.go` | Emitter unit tests (ID gen, retry, nil-safety) |
| `internal/api/operation_handler_test.go` | Handler tests |

### Modified Files

| File | Change |
|---|---|
| `internal/api/org_handler.go` | Inject nil-safe OperationEmitter; replace TODO markers with emission |
| `internal/api/ou_handler.go` | Same |
| `internal/api/tenant_handler.go` | Same |
| `internal/api/project_handler.go` | Same |
| `internal/api/*_handler_test.go` | Update constructor call sites (nil emitter) |
| `internal/server/server.go` | Accept `*api.OperationHandler`, register `/v1/operations` routes |
| `internal/server/server_test.go` | Update server constructor tests; add route tests |
| `cmd/sovrunn-api/main.go` | Wire OperationRegistry, emitter, OperationHandler into handlers/server |

---

## Components and Interfaces

| Component | Package | Responsibility |
|---|---|---|
| `Operation` model | `internal/resources` | Data shape, type/phase constants |
| `OperationRegistryIface` / `OperationRegistry` | `internal/registry` | Storage-only in-memory store |
| `OperationEmitter` interface | `internal/api` | Emission boundary consumed by handlers |
| `registryEmitter` adapter | `internal/registry` (or `internal/api`) | Generates ID, stores Operation via registry |
| `OperationHandler` | `internal/api` | Read API (`GET /v1/operations`, `GET /v1/operations/{name}`) |
| mutating handlers | `internal/api` | Emit Operations after successful actions (nil-safe) |

Detailed interface and struct definitions follow in the Data Models and per-component sections below.

---

## Data Models

### internal/resources/operation.go

```go
package resources

type Operation struct {
    APIVersion string          `json:"apiVersion"`
    Kind       string          `json:"kind"`
    Metadata   Metadata        `json:"metadata"`
    Spec       OperationSpec   `json:"spec"`
    Status     OperationStatus `json:"status"`
}
```

```go
// OperationSpec is set at emission time. It records only non-sensitive
// resource-reference fields and the operation type. No raw bodies or secrets.
type OperationSpec struct {
    Type                 string `json:"type"`
    ResourceKind         string `json:"resourceKind"`
    ResourceName         string `json:"resourceName"`
    OrganizationName     string `json:"organizationName,omitempty"`
    OrganizationUnitName string `json:"organizationUnitName,omitempty"`
    TenantName           string `json:"tenantName,omitempty"`
    ProjectName          string `json:"projectName,omitempty"`
    Actor                string `json:"actor"`
    RequestID            string `json:"requestId,omitempty"`
}

// OperationStatus is system-owned. Phase 1 sets Phase=Succeeded only;
// other phases are reserved for future async workflows.
type OperationStatus struct {
    Phase       string `json:"phase"`
    Message     string `json:"message,omitempty"`
    CreatedAt   string `json:"createdAt"`   // RFC3339 UTC
    UpdatedAt   string `json:"updatedAt"`   // RFC3339 UTC
    CompletedAt string `json:"completedAt,omitempty"` // RFC3339 UTC, terminal only
}

const (
    OperationAPIVersion = "platform.sovrunn.io/v1alpha1"
    OperationKind       = "Operation"
)

// Operation phase constants (all four reserved; Phase 1 emits Succeeded).
const (
    OperationPhasePending   = "Pending"
    OperationPhaseRunning   = "Running"
    OperationPhaseSucceeded = "Succeeded"
    OperationPhaseFailed    = "Failed"
)

// Operation type constants.
const (
    OpCreateOrganization     = "CreateOrganization"
    OpUpdateOrganization     = "UpdateOrganization"
    OpDeleteOrganization     = "DeleteOrganization"
    OpCreateOrganizationUnit = "CreateOrganizationUnit"
    OpUpdateOrganizationUnit = "UpdateOrganizationUnit"
    OpDeleteOrganizationUnit = "DeleteOrganizationUnit"
    OpCreateTenant           = "CreateTenant"
    OpUpdateTenant           = "UpdateTenant"
    OpDeleteTenant           = "DeleteTenant"
    OpCreateProject          = "CreateProject"
    OpUpdateProject          = "UpdateProject"
    OpDeleteProject          = "DeleteProject"
)
```

`CreatedAt`/`UpdatedAt`/`CompletedAt` are RFC3339 UTC strings for JSON stability and to avoid
`time.Time` marshaling surprises. Reuses the existing `Metadata` struct (Operation ID in `Name`).

---

## OperationRegistry Design (storage-only)

### internal/registry/operation_registry.go

```go
type OperationRegistryIface interface {
    CreateOperation(ctx context.Context, op resources.Operation) (resources.Operation, error)
    GetOperation(ctx context.Context, id string) (resources.Operation, error)
    ListOperations(ctx context.Context) ([]resources.Operation, error)
}

type OperationRegistry struct {
    mu    sync.RWMutex
    store map[string]resources.Operation // key: Operation ID (metadata.name)
}

func NewOperationRegistry() *OperationRegistry {
    return &OperationRegistry{store: make(map[string]resources.Operation)}
}
```

### Method Behavior

| Method | Lock | Behavior |
|---|---|---|
| CreateOperation | Lock | Require `op.Metadata.Name != ""`. IF empty → return a concrete non-nil error (`ErrMissingOperationID`, a new sentinel in the registry package) and DO NOT store. IF key already exists → return `ErrAlreadyExists` and DO NOT overwrite. Else store deep copy, return deep copy. |
| GetOperation | RLock | Return deep copy; `ErrNotFound` if absent. |
| ListOperations | RLock | Return non-nil slice of deep copies, sorted by `status.createdAt` asc, then ID asc. |

- **Storage-only, no ID generation:** the registry NEVER generates Operation IDs — the caller
  (emitter) must supply `metadata.name`. Empty `metadata.name` is a caller error and returns
  `ErrMissingOperationID` (non-nil), never a silently-generated ID.
- Reuses `ErrNotFound` and `ErrAlreadyExists` sentinels from `internal/registry/registry.go`, plus a
  new `ErrMissingOperationID` sentinel. No dependency on other registries. No global state.
- `deepCopyOperation` copies the struct; OperationSpec/Status are all value types (strings), so a
  plain struct copy suffices. If `Metadata.Labels`/`Annotations` are populated, copy those maps too.

---

## Operation ID Generator Design

### internal/api/operation_emitter.go

```go
// newOperationID returns a URL-safe, path-segment-safe unique token
// using crypto/rand (16 bytes, hex-encoded → 32 hex chars).
func newOperationID() (string, error) {
    var b [16]byte
    if _, err := rand.Read(b[:]); err != nil {
        return "", err
    }
    return hex.EncodeToString(b[:]), nil
}
```

- IDs are opaque server-generated tokens — NOT DNS-label validated.
- Generated in the emitter layer, never in the registry.

---

## OperationEmitter Design

### Interface (consumed by handlers)

```go
// internal/api/operation_emitter.go

// OperationEmitter records a control-plane Operation after a successful
// mutating action. Implementations must be safe to call with a nil receiver
// check by the caller (handlers guard against a nil emitter).
type OperationEmitter interface {
    Emit(ctx context.Context, spec resources.OperationSpec) error
}
```

### Adapter Implementation

The adapter does NOT introduce a new logger dependency into the mutating handlers. It may hold an
optional logger ONLY if the existing codebase already exposes an API-safe logger pattern (e.g. the
same logger already used by `loggingMiddleware`). If no such shared logger is readily available in
`internal/api`, the emitter records failures by returning an error to the caller, and the nil-safe
`emitOperation` helper swallows it without logging. Do NOT add logger fields to OrgHandler,
OUHandler, TenantHandler, or ProjectHandler solely for Operation emission.

```go
type registryEmitter struct {
    registry registry.OperationRegistryIface
    // logger is OPTIONAL: set only if an existing API-safe logger pattern exists.
    // If nil, emission failures are surfaced via the returned error only.
    logger *log.Logger
}

// NewRegistryEmitter constructs the adapter. logger may be nil.
func NewRegistryEmitter(reg registry.OperationRegistryIface, logger *log.Logger) *registryEmitter

// Emit generates an ID, sets server-owned fields, and stores the Operation.
// Bounded retry (5 attempts) on ErrAlreadyExists collision.
func (e *registryEmitter) Emit(ctx context.Context, spec resources.OperationSpec) error {
    now := time.Now().UTC().Format(time.RFC3339)
    spec.Actor = "system"
    for attempt := 0; attempt < 5; attempt++ {
        id, err := newOperationID()
        if err != nil {
            return err
        }
        op := resources.Operation{
            APIVersion: resources.OperationAPIVersion,
            Kind:       resources.OperationKind,
            Metadata:   resources.Metadata{Name: id},
            Spec:       spec,
            Status: resources.OperationStatus{
                Phase:       resources.OperationPhaseSucceeded,
                CreatedAt:   now,
                UpdatedAt:   now,
                CompletedAt: now,
            },
        }
        _, err = e.registry.CreateOperation(ctx, op)
        if errors.Is(err, registry.ErrAlreadyExists) {
            continue // ID collision — retry with a new ID
        }
        return err // nil on success, or a non-collision error
    }
    return errOperationIDExhausted
}
```

- `Actor` is forced to `"system"` in Phase 1.
- `RequestID` is expected to already be set on `spec` by the caller (from request context).
- Phase 1 always sets `Succeeded` + `CompletedAt`.

### Handler Emission Helper (nil-safe, no handler logger dependency)

Each mutating handler holds an `emitter OperationEmitter` field that MAY be nil. A small helper
centralizes nil-safety and failure isolation. It takes NO logger parameter — handlers do not gain a
logger field for emission. Any failure logging is the emitter adapter's concern (and only if an
existing API-safe logger pattern is available):

```go
// emitOperation records an Operation but never affects the primary response.
// It is nil-safe and swallows emitter errors so the primary handler flow is unaffected.
func emitOperation(ctx context.Context, emitter OperationEmitter, spec resources.OperationSpec) {
    if emitter == nil {
        return // isolated handler tests may omit the emitter
    }
    _ = emitter.Emit(ctx, spec) // failures are swallowed; emitter may log internally
}
```

The `registryEmitter.Emit` adapter MAY log a safe warning internally IF an existing API-safe logger
is available; otherwise it simply returns the error, which `emitOperation` swallows. Per Requirement
6, complex panic-recovery is not required. Emission occurs only after the mutating registry action
succeeds.

---

## Handler Emission Integration

Each mutating handler gains an `emitter OperationEmitter` field (nil-safe) and its constructor
gains an `emitter` parameter. At the point where `// TODO(FEATURE-0005): emit Operation record`
currently sits — AFTER the registry mutation succeeds and BEFORE writing the response — the
handler calls `emitOperation(...)`.

### Example: OrgHandler.Create

```
1. ... decode, validate, registry.CreateOrganization succeeds → created ...
2. emitOperation(ctx, h.emitter, resources.OperationSpec{
       Type:             resources.OpCreateOrganization,
       ResourceKind:     resources.OrganizationKind,
       ResourceName:     created.Metadata.Name,
       OrganizationName: created.Metadata.Name,
       RequestID:        requestIDFromContext(ctx),
   })
3. writeJSON(201, created)
```

### Request ID helper (API-local, no internal/server import)

`internal/api` MUST NOT import `internal/server`. Use the existing API-safe request ID context
helper if one already lives in an API-accessible package. If the request ID context key/helper
currently lives only in `internal/server`, define an API-local `requestIDFromContext(ctx)` in
`internal/api` that reads the same context value, following the existing middleware pattern (the
middleware writes the request ID into the context; the API layer reads it via its own helper). Do
NOT introduce an import cycle or import `internal/server` from `internal/api`.

### Spec field population by resource

Use the existing resource kind constants from `internal/resources` (add any missing constant only
if it fits the existing pattern):

| Handler | ResourceKind constant | ResourceName | Parent ref fields set |
|---|---|---|---|
| Organization | `resources.OrganizationKind` | org name | organizationName = org name |
| OrganizationUnit | `resources.OrganizationUnitKind` | OU name | organizationName, organizationUnitName |
| Tenant | `resources.TenantKind` | tenant name | organizationName, organizationUnitName, tenantName |
| Project | `resources.ProjectKind` | project name | organizationName, organizationUnitName, tenantName, projectName |

Operation type is the matching constant (`OpCreate*`, `OpUpdate*`, `OpDelete*`). For delete, the
reference fields are taken from the path segments; for create/update, from the created/updated
resource. All 12 TODO markers are replaced. Post-implementation, a repo-wide grep for
`TODO(FEATURE-0005): emit Operation record` MUST return zero matches (Requirement 5.7).

---

## Operation HTTP Handler Design

### internal/api/operation_handler.go

```go
type OperationHandler struct {
    registry registry.OperationRegistryIface
}

func NewOperationHandler(reg registry.OperationRegistryIface) *OperationHandler

func (h *OperationHandler) HandleCollection(w http.ResponseWriter, r *http.Request) // GET=List, POST=405, else 405
func (h *OperationHandler) HandleItem(w http.ResponseWriter, r *http.Request)       // GET=Get, else 405
```

### Path parsing (Go 1.21)

```
remainder := strings.TrimPrefix(r.URL.Path, "/v1/operations/")
if remainder == "" || remainder == r.URL.Path → 404   // bare path
if strings.Contains(remainder, "/") → 404             // extra segments, e.g. /v1/operations/{name}/extra
id := remainder   // opaque token, no DNS-label validation
// If remainder contains an additional "/" (extra segments) → 404
```

### List Flow (GET /v1/operations)

```
1. items, err := registry.ListOperations(ctx)
2. If error → writeError(500, INTERNAL_ERROR)
3. writeJSON(200, {"items": items})  // [] when empty
```

### Get Flow (GET /v1/operations/{name})

```
1. Extract id from path; bare/empty or extra segments → 404
2. op, err := registry.GetOperation(ctx, id)
3. If ErrNotFound → writeError(404, RESOURCE_NOT_FOUND)
4. If other error → writeError(500, INTERNAL_ERROR)
5. writeJSON(200, op)
```

### POST /v1/operations

`HandleCollection` returns HTTP 405 Method Not Allowed for POST (and any non-GET method).
No client-facing Operation creation endpoint exists in Phase 1.

---

## Error Handling

### Error Mapping Table

| Condition | HTTP | error.code |
|---|---|---|
| ErrNotFound (Get) | 404 | RESOURCE_NOT_FOUND |
| Bare/malformed/extra-segment path (Get item) | 404 | RESOURCE_NOT_FOUND |
| POST or unsupported method | 405 | (method not allowed) |
| Registry internal error (List/Get) | 500 | INTERNAL_ERROR |

Registry-layer errors: `ErrMissingOperationID` (empty ID on Create), `ErrAlreadyExists` (duplicate
ID), `ErrNotFound` (Get miss). These are consumed by the emitter/handler, not surfaced directly.

Emission-path errors are NEVER mapped to client responses (Requirement 6): a nil emitter is skipped,
and emitter errors (including ID-collision exhaustion) are swallowed by `emitOperation`.

---

## Server and main.go Wiring

### server.New signature (extended)

```go
func New(
    cfg config.Config,
    org *api.OrgHandler,
    ou *api.OUHandler,
    tenant *api.TenantHandler,
    project *api.ProjectHandler,
    operation *api.OperationHandler,  // NEW
    bootstrap *api.BootstrapHandler,
    readiness *health.ReadinessState,
) *Server
```

`server.New` SHALL require a non-nil `*api.OperationHandler`. Register:

```go
mux.Handle("/v1/operations", chain(http.HandlerFunc(operation.HandleCollection)))
mux.Handle("/v1/operations/", chain(http.HandlerFunc(operation.HandleItem)))
```

Middleware order unchanged: requestID → logging → contentType → handler.

### main.go wiring

```go
operationRegistry := registry.NewOperationRegistry()
emitter := api.NewRegistryEmitter(operationRegistry, logger)

orgHandler := api.NewOrgHandler(orgRegistry, ouBlocker, emitter)
ouHandler := api.NewOUHandler(ouRegistry, orgRegistry, tenantBlocker, emitter)
tenantHandler := api.NewTenantHandler(tenantRegistry, ouRegistry, projectBlocker, emitter)
projectHandler := api.NewProjectHandler(projectRegistry, tenantRegistry, emitter)
operationHandler := api.NewOperationHandler(operationRegistry)

srv := server.New(cfg, orgHandler, ouHandler, tenantHandler, projectHandler,
    operationHandler, bootstrapHandler, readiness)
```

Each mutating handler constructor gains a trailing `emitter OperationEmitter` parameter.
Existing handler tests pass `nil` for the emitter (nil-safe emission).

---

## Emission Failure Behavior

- Emission runs only after the mutating registry action succeeds.
- `emitOperation` isolates failures: a nil emitter is skipped; an emitter error is logged as a
  safe warning and swallowed. The primary HTTP response (201/200/204) is returned unchanged.
- Bounded ID-collision retry (5 attempts) inside the emitter; exhaustion returns
  `errOperationIDExhausted`, which is logged and swallowed by `emitOperation`.
- No panic-recovery machinery is required in Phase 1.

## Logging Behavior for Emission Failures

- Log level `warn`, single structured line with fields: `type`, `resourceKind`, `resourceName`,
  `requestId`, `error`.
- NEVER log secrets, tokens, authorization headers, or raw request bodies.
- Successful emission is not logged at warn level (optional debug logging only).

## Security and Privacy Constraints

- Operation records store only: type, resource kind/name, parent-path reference fields, actor,
  request ID, phase, timestamps.
- NO raw request bodies, secrets, tokens, credentials, or authorization header values in
  Operations or logs.
- `RequestID` is the tracing ID only.

---

## Testing Strategy

### Registry tests (`operation_registry_test.go`)
- Create requires non-empty `metadata.name`; empty → error, not stored.
- Create stores; duplicate ID → `ErrAlreadyExists`, original unchanged.
- Get by ID; missing → `ErrNotFound`.
- List sorted by `createdAt` then ID; empty → non-nil `[]`.
- Deep-copy immutability: mutating returned value does not affect stored state.

### Race test (`operation_registry_race_test.go`)
- 10+ goroutines mixed Create/Get/List; zero race reports under `go test -race`.

### Emitter tests (`operation_emitter_test.go`)
- `newOperationID` returns unique, non-empty, hex tokens.
- Emit sets Actor=system, Phase=Succeeded, CreatedAt/UpdatedAt/CompletedAt.
- Collision retry: stub registry returning `ErrAlreadyExists` N times then success.
- Exhaustion after 5 collisions returns error (verified via stub).
- `emitOperation` with nil emitter is a no-op (no panic).
- `emitOperation` swallows emitter errors (primary flow unaffected).

### Handler tests (`operation_handler_test.go`)
- GET list 200 sorted / empty `[]`.
- GET item 200 / 404 missing / 404 bare path (`/v1/operations/`).
- GET item 404 for extra segments (`/v1/operations/{name}/extra`) — explicit test case.
- POST → 405.

### Registry tests — empty ID case
- `CreateOperation` with empty `metadata.name` returns `ErrMissingOperationID` and does NOT store
  (a subsequent `ListOperations` returns no entry for it).

### Emission integration tests (implement last)
- Implement these AFTER the resource model, registry, emitter, and handler tests are stable, since
  they depend on all of those layers being in place.
- Table-driven across Organization, OrganizationUnit, Tenant, Project × create/update/delete:
  verify an Operation is recorded with correct type, resource reference, actor=system,
  phase=Succeeded (Requirement 13.5).
- Simulated emission failure does NOT change the primary handler response (Requirement 13.6).

### Verification
- `make fmt && make vet && make test && go test -race ./... && make build`
- Repo-wide grep for `TODO(FEATURE-0005): emit Operation record` returns zero matches.

---

## Correctness Properties

These correctness properties are optional to implement as `testing/quick` tests in Phase 1. If
added, each test function carries a `Feature: operation-resource` comment tag with its number and title.

### Property 1: Emitted Operations are well-formed

Any Operation stored via the emitter has a non-empty ID and `status.phase = "Succeeded"`.

**Validates: Requirements 3.2, 5.3, 10.1**

### Property 2: List ordering is deterministic

`ListOperations` output is always sorted by `status.createdAt` ascending, then Operation ID ascending.

**Validates: Requirements 4.11, 7.3**

### Property 3: Reads return deep copies

`GetOperation` returns a deep copy: mutating the returned value never affects stored registry state.

**Validates: Requirements 4.3, 8.3**

### Property 4: Create never overwrites

`CreateOperation` with an existing ID never overwrites the stored entry and returns `ErrAlreadyExists`.

**Validates: Requirements 4.5**

---

## Non-Goals

1. Async workflow engine, background queue, or background workers.
2. Retries beyond the bounded ID-collision retry inside the emitter.
3. Step execution engine, approval workflows.
4. ServiceOps execution, GitOps integration.
5. Observability pipeline, persistent database storage, Kubernetes CRDs/controllers.
6. Authentication, authorization, RBAC.
7. AI automation, billing, UI.
8. Distributed locking, multi-node coordination.
9. Client-facing Operation creation endpoint (POST returns 405).
10. Go 1.22 wildcard routing; no new external dependencies.
11. Storing raw request bodies or secrets in Operations or logs.

---

## Design Questions (Resolved)

| # | Question | Resolution |
|---|---|---|
| 1 | OperationEmitter interface placement | Narrow `OperationEmitter` interface in `internal/api`, consumed by handlers; `registryEmitter` adapter implements it over the storage-only registry. Handlers never depend on the concrete registry. |
| 2 | ID generation location | Generated in the emitter via `crypto/rand`; registry is storage-only and requires `metadata.name`; bounded retry (5) on `ErrAlreadyExists`; exhaustion = emission failure (swallowed). |
| 3 | Phase 1 status phase | Emit `Succeeded` only; `Pending`/`Running`/`Failed` reserved in the model for future async workflows; no async behavior implemented. |
