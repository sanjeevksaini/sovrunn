# FEATURE-0009: API Server Health and Readiness — Design

## Overview

This design formalizes and hardens the existing `/healthz`, `/readyz`, and `/version` endpoints
in the sovrunn-api server. The bootstrap implementation already exists in
`internal/api/bootstrap_handler.go` and `internal/health/readiness.go`. This feature adds:

- HTTP method enforcement (405 for non-GET) as a dedicated middleware layer
- Structured not-ready responses with fixed reason constants (`ReasonInitializing`, `ReasonShuttingDown`)
- Proper `Allow: GET` header on 405 responses
- HEAD-specific behavior (no body in 405 response per HTTP semantics)
- Consistent use of `resources.APIErrorEnvelope` for error conditions
- Graceful shutdown readiness interaction (ReadinessState → false before drain)
- Build-time version injection via Makefile ldflags (`-X 'github.com/sanjeevksaini/sovrunn/internal/api.buildVersion=$(VERSION)'`)
- Zero-value-safe `ReadinessState` (Reason() returns `"initializing"` even on an uninitialized struct)
- Comprehensive test coverage including method enforcement and shutdown behavior

No new packages are introduced. No external dependencies are added.

---

## Resolved Design Decisions

### 1. Method enforcement placement

Method enforcement is implemented as a **dedicated middleware function** in `internal/server/middleware.go`
named `methodGET`. The bootstrap route chain is composed as:

```go
bootstrapChain := func(h http.Handler) http.Handler {
    return requestIDMiddleware(loggingMiddleware(logger)(methodGET(h)))
}
```

Canonical ordering (outermost to innermost):
```text
requestIDMiddleware → loggingMiddleware → methodGET → endpoint handler
```

This means `loggingMiddleware` wraps `methodGET`, so 405 responses produced by `methodGET`
are captured in structured logs (status_code, latency_ms). The request-ID middleware is
outermost so that all responses—including 405—include the `X-Sovrunn-Request-ID` header.

Rationale: A dedicated middleware keeps handler methods clean and focused on domain logic. It is
reusable across all three bootstrap endpoints without duplication. Inline checks in each handler
would violate the DRY principle and mix concerns.

### 2. Shutdown race window

The current `sync/atomic.Bool` approach in `ReadinessState` is **sufficient** for Phase 1.

The `Server.Shutdown` method calls `s.readiness.SetShuttingDown()` before
`s.httpServer.Shutdown(ctx)`. The atomic stores are visible to concurrent readers immediately.
`http.Server.Shutdown` then stops accepting new connections and waits for in-flight requests
to complete within the timeout. This means:

- Any `/readyz` request arriving after `SetShuttingDown()` will see `false`.
- Any `/readyz` request in-flight during the atomic store may see either value (benign race
  on behavior, not on memory safety). The next probe resolves it.

No `http.Server.RegisterOnShutdown` callback is needed. The explicit `SetShuttingDown()` call
already provides the required ordering guarantee.

### 3. Version field extensibility

The `/version` response uses a **fixed struct** (`versionResponse`) with exported fields. Only
the four required fields (`name`, `version`, `phase`, `status`) are included. No optional or
dynamic fields are added. Future extensions (e.g. `goVersion`, `buildTime`) will be additive
struct fields — adding a JSON field is backward-compatible for clients that ignore unknown keys.

### 4. Makefile ldflags target

The existing `build` target in the Makefile is **modified** to include ldflags. A separate
`release` target is not needed in Phase 1. The `VERSION` variable defaults to `dev` and can
be overridden at invocation time (`make build VERSION=1.0.0`).

The injected variable is `github.com/sanjeevksaini/sovrunn/internal/api.buildVersion`, which
is declared in `internal/api/bootstrap_handler.go` as:
```go
var buildVersion = "dev"
```

The Makefile already defines `APP_NAME=sovrunn-api`. The output binary is `bin/$(APP_NAME)`.

The `phase` and `status` fields in the `/version` response are sourced from string constants
defined in `internal/api/bootstrap_handler.go`:
```go
const versionPhase  = "1"
const versionStatus = "alpha"
```

These are compile-time constants for Phase 1 and do not come from config or ldflags.

### 5. Readyz not-ready response shape

The 503 response from `/readyz` uses a **domain-specific shape**:
```json
{"status":"not_ready","message":"initializing"}
```

This is NOT wrapped in `resources.APIErrorEnvelope` because it represents an expected
operational state, not an error condition. Only true errors (405 Method Not Allowed) use the
error envelope on these endpoints.

### 6. Readyz message for pre-startup vs post-shutdown

`ReadinessState` is extended with a **reason field** (atomic string pointer) that holds the
current not-ready reason. The package defines exactly two exported constants:

```go
const ReasonInitializing = "initializing"
const ReasonShuttingDown = "shutting_down"
```

**Zero-value safety**: `Reason()` returns `ReasonInitializing` when `IsReady()` is false and
no explicit reason has been set (i.e., the atomic pointer is nil). This ensures the struct is
fully usable at its zero value without requiring the constructor. A `NewReadinessState()`
constructor is still provided for clarity but is not required for correctness.

**Fixed-reason enforcement**: `SetReason` is unexported. Instead, two explicit methods enforce
the contract:
- `SetInitializing()` — stores `ReasonInitializing` and sets ready=false.
- `SetShuttingDown()` — stores `ReasonShuttingDown` and sets ready=false.

`SetReady(true)` clears the reason pointer. This guarantees the not-ready message is always
one of the two documented constants.

The `Readyz` handler reads both the boolean and the reason atomically.

---

## Architecture

### Component Interaction

```text
HTTP Request
  │
  ├─ requestIDMiddleware  (assigns/propagates X-Sovrunn-Request-ID)  [outermost]
  │
  ├─ loggingMiddleware    (logs method, path, status_code, latency_ms, request_id)
  │
  ├─ methodGET middleware (rejects non-GET → 405 + Allow: GET)       [innermost wrapper]
  │
  └─ endpoint handler     (Healthz | Readyz | Version)
```

Canonical composition in code:
```go
bootstrapChain := func(h http.Handler) http.Handler {
    return requestIDMiddleware(loggingMiddleware(logger)(methodGET(h)))
}
```

### Package Responsibilities

| Package | Role in FEATURE-0009 |
|---------|---------------------|
| `internal/server` | Registers routes, composes middleware chain, owns shutdown sequence |
| `internal/api` | `BootstrapHandler` with `Healthz`, `Readyz`, `Version` methods |
| `internal/health` | `ReadinessState` atomic state with reason tracking |
| `internal/resources` | `APIErrorEnvelope`, `ErrorCode` constants (existing, no changes) |
| `internal/config` | `Config` struct (existing, no changes) |
| `internal/requestctx` | Request-ID context helpers (existing, no changes) |

### Dependency Direction

```text
cmd/sovrunn-api → internal/server → internal/api → internal/health
                                  → internal/api → internal/config
                                  → internal/api → internal/resources
                                  → internal/api → internal/requestctx
```

`internal/api` does NOT import `internal/server`. This constraint is preserved.

---

## Files Changed

| File | Action | Purpose |
|------|--------|---------|
| `internal/health/readiness.go` | Modify | Add `ReasonInitializing`, `ReasonShuttingDown` constants; add `reason atomic.Pointer[string]`; add `SetInitializing`, `SetShuttingDown`, `Reason` methods; make zero-value safe |
| `internal/health/readiness_test.go` | Modify | Add tests for reason tracking, zero-value safety, typed helpers |
| `internal/api/bootstrap_handler.go` | Modify | Update `Readyz` to include reason in 503 body; fix status value from `"not-ready"` to `"not_ready"`; add `versionPhase`/`versionStatus` constants; extract `versionResponse` struct |
| `internal/api/bootstrap_handler_test.go` | Modify | Expand tests: method enforcement, HEAD handling, 503 reason strings, Content-Type assertions |
| `internal/server/middleware.go` | Modify | Add `methodGET` middleware function |
| `internal/server/server.go` | Modify | Wire `methodGET` into bootstrapChain composition (inside logging); use `SetShuttingDown()` in shutdown path |
| `internal/server/server_test.go` | Modify | Add integration tests for method enforcement and shutdown behavior |
| `Makefile` | Modify | Add `VERSION`, `MODULE`, `LDFLAGS` variables; inject `-X '$(MODULE)/internal/api.buildVersion=$(VERSION)'` in `build` target |

No new files are created. No packages are added.

---

## Data Models

### ReadinessState (modified)

```go
// Package: internal/health

const ReasonInitializing = "initializing"
const ReasonShuttingDown = "shutting_down"

type ReadinessState struct {
    ready  atomic.Bool
    reason atomic.Pointer[string]
}
```

**Zero-value behavior**: When the struct is created as `&ReadinessState{}`, `ready` is `false`
and `reason` is nil. `Reason()` detects the nil pointer and returns `ReasonInitializing` when
`IsReady()` is false, making the zero value fully usable without a constructor.

`NewReadinessState()` explicitly stores a pointer to `ReasonInitializing` for documentation
clarity, but the behavior is identical to the zero value.

### versionResponse (new struct in internal/api)

```go
type versionResponse struct {
    Name    string `json:"name"`
    Version string `json:"version"`
    Phase   string `json:"phase"`
    Status  string `json:"status"`
}
```

### notReadyResponse (new struct in internal/api)

```go
type notReadyResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}
```

### healthResponse (new struct in internal/api)

```go
type healthResponse struct {
    Status string `json:"status"`
}
```

### readyResponse (new struct in internal/api)

```go
type readyResponse struct {
    Status string `json:"status"`
}
```

Using named structs instead of `map[string]string` provides compile-time safety, avoids
allocation of map internals per request, and documents the response contract in code.

---

## Interfaces

### ReadinessState API (modified)

```go
// NewReadinessState creates a ReadinessState with reason="initializing".
// The zero value is also usable — Reason() returns "initializing" when not ready
// and reason is unset.
func NewReadinessState() *ReadinessState

// SetReady marks the server as ready and clears the reason.
func (s *ReadinessState) SetReady(v bool)

// SetInitializing sets ready=false and reason=ReasonInitializing.
func (s *ReadinessState) SetInitializing()

// SetShuttingDown sets ready=false and reason=ReasonShuttingDown.
func (s *ReadinessState) SetShuttingDown()

// IsReady returns whether the server has completed initialization.
func (s *ReadinessState) IsReady() bool

// Reason returns the current not-ready reason. Returns ReasonInitializing
// when not ready and no explicit reason has been stored (zero-value safe).
// Returns "" when ready.
func (s *ReadinessState) Reason() string
```

### BootstrapHandler methods (unchanged signatures)

```go
func (h *BootstrapHandler) Healthz(w http.ResponseWriter, r *http.Request)
func (h *BootstrapHandler) Readyz(w http.ResponseWriter, r *http.Request)
func (h *BootstrapHandler) Version(w http.ResponseWriter, r *http.Request)
```

### methodGET middleware (new in internal/server/middleware.go)

```go
// methodGET rejects non-GET requests with 405 Method Not Allowed.
// For HEAD: responds with 405, Allow: GET, Content-Type: application/json, no body.
// For other methods: responds with 405, Allow: GET, APIErrorEnvelope JSON body.
// Uses resources.ErrCodeMethodNotAllowed ("METHOD_NOT_ALLOWED") for the error code.
func methodGET(next http.Handler) http.Handler
```

This middleware is composed as the innermost layer of `bootstrapChain`, wrapping the endpoint
handler directly. It does NOT wrap `loggingMiddleware` — logging wraps it so that 405
responses are logged.

---

## Validation

No resource validation is involved. The endpoints are simple read-only probes.

Input validation is limited to:
- HTTP method check (enforced by `methodGET` middleware)
- Unknown query parameters: ignored (no validation error)
- Request body on GET: ignored (no validation error)

---

## API / Handler Design

### GET /healthz

**Request**: No body, no required headers.

**Response 200**:
```json
{"status":"ok"}
```

**Behavior**: Always returns 200 if the process is alive. No registry or subsystem checks.

### GET /readyz

**Request**: No body, no required headers.

**Response 200** (when ready):
```json
{"status":"ready"}
```

**Response 503** (when not ready):
```json
{"status":"not_ready","message":"initializing"}
```
or
```json
{"status":"not_ready","message":"shutting_down"}
```

**Behavior**: Reads `ReadinessState.IsReady()`. If false, reads `ReadinessState.Reason()`
to populate the `message` field. The reason is always one of the two fixed constants
(`ReasonInitializing` or `ReasonShuttingDown`), enforced by the typed helper methods.

### GET /version

**Request**: No body, no required headers.

**Response 200**:
```json
{"name":"sovrunn-api","version":"dev","phase":"1","status":"alpha"}
```

**Behavior**: Returns static build metadata. The `version` field is set via ldflags at build
time; defaults to `"dev"`.

### Non-GET methods (all three endpoints)

**Response 405** (non-HEAD methods):
```json
{"error":{"code":"METHOD_NOT_ALLOWED","message":"only GET is supported"}}
```
Headers: `Allow: GET`, `Content-Type: application/json`

**Response 405** (HEAD method):
No body. Headers: `Allow: GET`, `Content-Type: application/json`, status 405.

---

## Registry / Storage Design

Not applicable. Health/readiness/version endpoints do not interact with the resource registry.
The only indirect dependency is that readiness checks confirm the registry subsystem initialized
during server startup (handled by the existing `SetReady(true)` call in `Server.Start()`).

---

## Operation / Audit Behavior

Health, readiness, and version requests do **NOT** emit Operation records. They are
non-mutating infrastructure probes, not governance actions. This is consistent with the
compatibility table in requirements (FEATURE-0005 compatibility note).

Structured logging via `loggingMiddleware` provides observability:
```text
request_id=<id> method=GET path=/healthz status_code=200 latency_ms=0
request_id=<id> method=POST path=/healthz status_code=405 latency_ms=0
```

No additional audit hooks are needed.

---

## Error Mapping

| Condition | HTTP Status | Response Shape |
|-----------|-------------|----------------|
| GET /healthz, process alive | 200 | `{"status":"ok"}` |
| GET /readyz, ready | 200 | `{"status":"ready"}` |
| GET /readyz, not ready | 503 | `{"status":"not_ready","message":"..."}` |
| GET /version | 200 | `{"name":"...","version":"...","phase":"...","status":"..."}` |
| Non-GET (except HEAD) on any bootstrap endpoint | 405 | `{"error":{"code":"METHOD_NOT_ALLOWED","message":"only GET is supported"}}` |
| HEAD on any bootstrap endpoint | 405 | No body; headers only |

Error codes used:
- `METHOD_NOT_ALLOWED` — existing constant `resources.ErrCodeMethodNotAllowed` in
  `internal/resources/errors.go` (verified: `ErrorCode = "METHOD_NOT_ALLOWED"`)

No new error codes are introduced.

---

## Security and Privacy

1. No internal IPs, hostnames, file paths, goroutine counts, or memory stats exposed.
2. `/version` returns only `name`, `version`, `phase`, `status`. No Git SHA, build machine,
   or CI identifiers.
3. `/readyz` 503 message contains only `"initializing"` or `"shutting_down"`. No subsystem
   names, error details, or stack traces.
4. No authentication required on any bootstrap endpoint in Phase 1.
5. `Authorization` headers are not read or validated by these endpoints.
6. Responses are small fixed-size JSON — not usable as amplification vectors.
7. No request bodies or authorization headers are logged.

---

## Testing Strategy

### Unit Tests (internal/api/bootstrap_handler_test.go)

| Test Case | Assertion |
|-----------|-----------|
| GET /healthz returns 200 | Status 200, body `{"status":"ok"}`, Content-Type header |
| GET /readyz when ready | Status 200, body `{"status":"ready"}` |
| GET /readyz when not ready (initializing) | Status 503, body `{"status":"not_ready","message":"initializing"}` |
| GET /readyz when not ready (shutting_down) | Status 503, body `{"status":"not_ready","message":"shutting_down"}` |
| GET /version | Status 200, all four fields present, name=sovrunn-api |
| POST /healthz (via method middleware) | Status 405, APIErrorEnvelope body, Allow header |
| PUT /readyz (via method middleware) | Status 405, APIErrorEnvelope body, Allow header |
| DELETE /version (via method middleware) | Status 405, APIErrorEnvelope body, Allow header |
| HEAD /healthz (via method middleware) | Status 405, Allow header, Content-Type header, NO body |
| PATCH /healthz (via method middleware) | Status 405, APIErrorEnvelope body, Allow header |

### Unit Tests (internal/health/readiness_test.go)

| Test Case | Assertion |
|-----------|-----------|
| Zero-value ReadinessState defaults to not ready | `IsReady() == false` |
| Zero-value ReadinessState Reason() returns "initializing" | `Reason() == "initializing"` (zero-value safe) |
| NewReadinessState defaults to not ready | `IsReady() == false` |
| NewReadinessState Reason() returns "initializing" | `Reason() == "initializing"` |
| SetReady(true) clears reason | `IsReady() == true`, `Reason() == ""` |
| SetReady(false) after SetReady(true) preserves nil reason → "initializing" | `IsReady() == false`, `Reason() == "initializing"` |
| SetShuttingDown sets ready=false and reason | `IsReady() == false`, `Reason() == "shutting_down"` |
| SetInitializing sets ready=false and reason | `IsReady() == false`, `Reason() == "initializing"` |
| SetShuttingDown then SetReady(true) clears reason | `IsReady() == true`, `Reason() == ""` |

### Unit Tests (internal/server/middleware_test.go)

| Test Case | Assertion |
|-----------|-----------|
| methodGET allows GET | Handler is called, no 405 |
| methodGET rejects POST | 405, Allow: GET, APIErrorEnvelope body |
| methodGET rejects HEAD | 405, Allow: GET, Content-Type header, empty body |
| methodGET rejects PUT | 405, Allow: GET, APIErrorEnvelope body |
| methodGET rejects DELETE | 405, Allow: GET, APIErrorEnvelope body |
| methodGET rejects PATCH | 405, Allow: GET, APIErrorEnvelope body |

### Integration Tests (internal/server/server_test.go)

| Test Case | Assertion |
|-----------|-----------|
| Server starts, /readyz returns 200 (poll with timeout) | Poll /readyz every 50ms, assert 200 within 5s |
| After readyz 200, /healthz returns 200 | Assert /healthz 200 |
| Shutdown sets readyz to 503 | Call Shutdown(), assert /readyz 503 with "shutting_down" (from SetShuttingDown) |
| Method enforcement end-to-end | POST /healthz via server returns 405 |

### Race Detection

All tests pass under `go test -race ./...`. The `ReadinessState` uses `atomic.Bool` and
`atomic.Pointer[string]` which are race-safe by design.

---

## Verification Commands

```bash
make fmt
make test
make vet
go test -race ./...
make build
make build VERSION=0.1.0
./bin/sovrunn-api --config configs/sovrunn-api.local.yaml &
curl -s http://127.0.0.1:8080/healthz
curl -s http://127.0.0.1:8080/readyz
curl -s http://127.0.0.1:8080/version
curl -s -X POST http://127.0.0.1:8080/healthz
curl -s -I -X HEAD http://127.0.0.1:8080/healthz
kill %1
```

`make build` without `VERSION` produces a binary with `buildVersion="dev"`.
`make build VERSION=0.1.0` produces a binary with `buildVersion="0.1.0"`.

---

## Non-Goals (not implemented in this feature)

1. Deep health checks (database, external service reachability)
2. Prometheus `/metrics` endpoint
3. Authentication on health endpoints
4. Custom/dynamic readiness check registration
5. Rate limiting on health endpoints
6. HTML health page
7. Distributed readiness coordination
8. OpenTelemetry health instrumentation
9. Startup probe endpoint (`/startupz`)
10. TLS termination
11. Verbose readiness mode (`/readyz?verbose=true`)

---

## Resolved Design Questions

### Q1: Method enforcement placement

**Answer**: Dedicated `methodGET` middleware in `internal/server/middleware.go`. Composed as
the innermost layer of `bootstrapChain`:
```go
bootstrapChain := func(h http.Handler) http.Handler {
    return requestIDMiddleware(loggingMiddleware(logger)(methodGET(h)))
}
```
`loggingMiddleware` wraps `methodGET`, so 405 responses are captured in structured logs.
This keeps handler methods focused on domain logic.

### Q2: Verbose readiness mode

**Answer**: Deferred. Not in scope for FEATURE-0009.

### Q3: Shutdown race window

**Answer**: Current `atomic.Bool` approach is sufficient. `Server.Shutdown` calls
`SetShuttingDown()` (which stores ready=false + reason="shutting_down" atomically in sequence)
before `httpServer.Shutdown(ctx)`. The atomic stores provide immediate visibility to concurrent
readers. No `RegisterOnShutdown` callback needed.

### Q4: Version field extensibility

**Answer**: Fixed struct (`versionResponse`). Future fields are additive. No dynamic field
maps.

### Q5: Makefile ldflags target

**Answer**: Modify existing `build` target. Add `VERSION ?= dev` variable. The module path
for the injected variable is `github.com/sanjeevksaini/sovrunn/internal/api.buildVersion`.
The existing `APP_NAME=sovrunn-api` variable defines the binary output name. Inject via:
```makefile
VERSION ?= dev
MODULE  = github.com/sanjeevksaini/sovrunn
LDFLAGS = -X '$(MODULE)/internal/api.buildVersion=$(VERSION)'

build:
	mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) ./cmd/sovrunn-api
```

---

## Implementation Notes

### Changes to existing code

1. **`internal/health/readiness.go`**: Add `reason atomic.Pointer[string]` field. Add
   exported constants `ReasonInitializing = "initializing"` and `ReasonShuttingDown = "shutting_down"`.
   Add `NewReadinessState()` constructor that stores a pointer to `ReasonInitializing`.
   Add `SetInitializing()` and `SetShuttingDown()` methods that set ready=false and store
   the respective reason constant. Add `Reason() string` method that returns `ReasonInitializing`
   when the pointer is nil and `IsReady()` is false (zero-value safety), the stored string
   when the pointer is non-nil and `IsReady()` is false, or `""` when `IsReady()` is true.
   Modify `SetReady(true)` to clear the reason pointer. Remove the need for a generic
   `SetReason(string)` method.

2. **`internal/api/bootstrap_handler.go`**: Add constants `versionPhase = "1"` and
   `versionStatus = "alpha"` as the source for `/version` response fields. Change `Readyz`
   to read reason from `ReadinessState.Reason()` and return
   `notReadyResponse{Status: "not_ready", Message: reason}`. Change status string from
   `"not-ready"` to `"not_ready"` (underscore, per requirements). Replace `map[string]string`
   with named response structs. The `buildVersion` variable (`var buildVersion = "dev"`)
   remains in this file and is the ldflags injection target.

3. **`internal/server/middleware.go`**: Add `methodGET(next http.Handler) http.Handler`.
   For HEAD: write 405 status + headers only. For other non-GET: write 405 +
   `APIErrorEnvelope` JSON body using `resources.ErrCodeMethodNotAllowed` (verified existing
   constant with value `"METHOD_NOT_ALLOWED"`).

4. **`internal/server/server.go`**: Change `bootstrapChain` composition to:
   `requestIDMiddleware(loggingMiddleware(logger)(methodGET(h)))`. This ensures logging
   captures 405 responses. Change all call sites that create `ReadinessState` to use
   `health.NewReadinessState()` (or accept the zero value, which is equally correct). In
   `Shutdown`, call `s.readiness.SetShuttingDown()` instead of `s.readiness.SetReady(false)`.
   In the error branch of `Start`, call `s.readiness.SetInitializing()` or just
   `s.readiness.SetReady(false)` (reason is already "initializing" from zero value).

5. **`Makefile`**: Add `VERSION ?= dev`, `MODULE = github.com/sanjeevksaini/sovrunn`, and
   `LDFLAGS = -X '$(MODULE)/internal/api.buildVersion=$(VERSION)'`. Update `build` target
   to pass `-ldflags "$(LDFLAGS)"`. The existing `APP_NAME=sovrunn-api` is used for the
   output binary name.

### Breaking change note

The `/readyz` 503 body changes from `{"status":"not-ready"}` to
`{"status":"not_ready","message":"initializing"}`. This is an intentional contract fix
per the requirements. The old hyphenated form was a bootstrap placeholder, not a documented
contract.

### Backward compatibility

- `/healthz` response is unchanged.
- `/version` response is unchanged (same four fields, same values).
- `/readyz` 200 response is unchanged (`{"status":"ready"}`).
- New: 405 responses on all three endpoints (previously returned 200 with handler content
  regardless of method).

---

## Go 1.21 Compatibility

- `sync/atomic.Bool` — available since Go 1.19.
- `sync/atomic.Pointer[string]` — available since Go 1.19 (generics in 1.18, atomic.Pointer
  in 1.19).
- `signal.NotifyContext` — available since Go 1.16.
- No features requiring Go 1.22+ are used.
