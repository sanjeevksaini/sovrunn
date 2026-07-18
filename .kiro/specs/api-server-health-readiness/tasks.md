# FEATURE-0009: API Server Health and Readiness — Tasks

## Task 1: Extend ReadinessState with reason tracking

### Objective
Add reason constants, atomic reason pointer, typed setter methods, and zero-value safety to `ReadinessState`.

### Files
- `internal/health/readiness.go`

### Notes
- Add exported constants `ReasonInitializing = "initializing"` and `ReasonShuttingDown = "shutting_down"`.
- Add `reason atomic.Pointer[string]` field to the struct.
- Add `NewReadinessState() *ReadinessState` constructor that stores a pointer to `ReasonInitializing`.
- Add `SetInitializing()` — stores `ReasonInitializing`, sets `ready=false`.
- Add `SetShuttingDown()` — stores `ReasonShuttingDown`, sets `ready=false`.
- Add `Reason() string` — returns `ReasonInitializing` when pointer is nil and not ready (zero-value safe), stored string when pointer is non-nil and not ready, `""` when ready.
- Modify `SetReady(true)` to clear the reason pointer (store nil).
- `SetReady(false)` does NOT change the reason (callers should use typed helpers).
- Do NOT add a generic `SetReason(string)` method.
- Zero value `&ReadinessState{}` must be fully usable without calling the constructor.

### Tests
- Add to `internal/health/readiness_test.go`:
  - Zero-value `Reason()` returns `"initializing"` when not ready.
  - `NewReadinessState()` defaults to not ready with reason `"initializing"`.
  - `SetReady(true)` clears reason → `Reason()` returns `""`.
  - `SetReady(false)` after `SetReady(true)` → `Reason()` returns `"initializing"` (nil pointer → default).
  - `SetShuttingDown()` → `IsReady() == false`, `Reason() == "shutting_down"`.
  - `SetInitializing()` → `IsReady() == false`, `Reason() == "initializing"`.
  - `SetShuttingDown()` then `SetReady(true)` → `Reason() == ""`.

### Acceptance criteria
- All existing tests in `internal/health/` still pass.
- New tests pass.
- `go vet ./internal/health/...` clean.
- `gofmt -l ./internal/health/` produces no output.

### Verification
```bash
docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'
```

### Commit message
```
feat(health): add reason tracking to ReadinessState

Add ReasonInitializing and ReasonShuttingDown constants with typed
setter methods. ReadinessState is zero-value safe — Reason() returns
"initializing" even on an uninitialized struct.

FEATURE-0009
```

---

## Task 2: Add methodGET middleware

### Objective
Implement the `methodGET` middleware in `internal/server/middleware.go` that rejects non-GET requests with HTTP 405.

### Files
- `internal/server/middleware.go`

### Notes
- Add function `methodGET(next http.Handler) http.Handler`.
- If method is GET, call `next.ServeHTTP(w, r)`.
- For HEAD: set `Allow: GET`, `Content-Type: application/json`, write status 405, do NOT write body.
- For all other methods (POST, PUT, DELETE, PATCH, etc.): set `Allow: GET`, `Content-Type: application/json`, write status 405, write `resources.APIErrorEnvelope` JSON body with code `resources.ErrCodeMethodNotAllowed` and message `"only GET is supported"`.
- Import `encoding/json` and `github.com/sanjeevksaini/sovrunn/internal/resources` (resources is already imported in this file).
- Use `json.NewEncoder(w).Encode(...)` for the error body (consistent with `writeErrorBody` in server.go).
- Do NOT move or refactor existing middleware functions.

### Tests
- Add to `internal/server/middleware_test.go`:
  - `methodGET` allows GET → handler called, no 405.
  - `methodGET` rejects POST → 405, `Allow: GET` header, `Content-Type: application/json`, APIErrorEnvelope body with code `"METHOD_NOT_ALLOWED"`.
  - `methodGET` rejects PUT → 405, `Allow: GET` header, APIErrorEnvelope body.
  - `methodGET` rejects DELETE → 405, `Allow: GET` header, APIErrorEnvelope body.
  - `methodGET` rejects PATCH → 405, `Allow: GET` header, APIErrorEnvelope body.
  - `methodGET` rejects HEAD → 405, `Allow: GET` header, `Content-Type: application/json`, empty body (assert `rec.Body.Len() == 0`).

### Acceptance criteria
- All existing middleware tests still pass.
- New tests pass.
- `go vet ./internal/server/...` clean.
- `gofmt -l ./internal/server/` produces no output.

### Verification
```bash
docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'
```

### Commit message
```
feat(server): add methodGET middleware for bootstrap endpoints

Rejects non-GET requests with 405 Method Not Allowed. HEAD receives
status and headers only (no body per HTTP semantics). Other methods
receive APIErrorEnvelope JSON body.

FEATURE-0009
```

---

## Task 3: Wire methodGET into bootstrapChain and update shutdown path

### Objective
Update `internal/server/server.go` to compose `methodGET` inside `bootstrapChain` and use `SetShuttingDown()` in the shutdown path.

### Files
- `internal/server/server.go`

### Notes
- Change `bootstrapChain` composition from:
  ```go
  bootstrapChain := func(h http.Handler) http.Handler {
      return requestIDMiddleware(loggingMiddleware(logger)(h))
  }
  ```
  to:
  ```go
  bootstrapChain := func(h http.Handler) http.Handler {
      return requestIDMiddleware(loggingMiddleware(logger)(methodGET(h)))
  }
  ```
- Canonical ordering: requestIDMiddleware → loggingMiddleware → methodGET → handler.
- In `Shutdown()` method: replace `s.readiness.SetReady(false)` with `s.readiness.SetShuttingDown()`.
- In the `Start()` error branch (`case err := <-errCh`): keep `s.readiness.SetReady(false)` — the zero value already has reason `"initializing"`, so this is correct (no need to call `SetInitializing()` since reason pointer is already nil → defaults to "initializing").
- Do NOT change the `Server.New()` function signature.
- Do NOT change the route registration pattern (the `mux.HandleFunc` closures remain).
- Do NOT change the `chain` (non-bootstrap) composition.

### Tests
- Existing tests in `internal/server/server_test.go` must still pass (they exercise route registration and readiness behavior).
- No new tests in this task (method enforcement is tested end-to-end in Task 5).

### Acceptance criteria
- All existing server tests pass.
- `go vet ./internal/server/...` clean.
- `gofmt -l ./internal/server/` produces no output.
- `bootstrapChain` now includes `methodGET`.
- `Shutdown()` calls `SetShuttingDown()`.

### Verification
```bash
docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'
```

### Commit message
```
feat(server): wire methodGET into bootstrapChain, use SetShuttingDown

Bootstrap endpoints now enforce GET-only via methodGET middleware.
Shutdown path uses SetShuttingDown() for explicit reason tracking.

FEATURE-0009
```

---

## Task 4: Update BootstrapHandler with response structs and reason support

### Objective
Replace `map[string]string` responses with named structs. Update `Readyz` to include the reason field in 503 responses. Fix `"not-ready"` → `"not_ready"`.

### Files
- `internal/api/bootstrap_handler.go`

### Notes
- Add response structs (unexported, in this file):
  ```go
  type healthResponse struct {
      Status string `json:"status"`
  }
  type readyResponse struct {
      Status string `json:"status"`
  }
  type notReadyResponse struct {
      Status  string `json:"status"`
      Message string `json:"message"`
  }
  type versionResponse struct {
      Name    string `json:"name"`
      Version string `json:"version"`
      Phase   string `json:"phase"`
      Status  string `json:"status"`
  }
  ```
- Add constants:
  ```go
  const versionPhase  = "1"
  const versionStatus = "alpha"
  ```
- Update `Healthz`: use `healthResponse{Status: "ok"}`.
- Update `Readyz`:
  - Ready: use `readyResponse{Status: "ready"}`.
  - Not ready: use `notReadyResponse{Status: "not_ready", Message: h.readiness.Reason()}`.
  - Note: status changes from `"not-ready"` (hyphen) to `"not_ready"` (underscore).
- Update `Version`: use `versionResponse{Name: "sovrunn-api", Version: buildVersion, Phase: versionPhase, Status: versionStatus}`.
- Remove method-level checks from handlers (method enforcement is now in `methodGET` middleware). Handlers only handle GET logic.
- Do NOT remove `buildVersion` variable.
- Do NOT change the handler method signatures.

### Tests
- Update `internal/api/bootstrap_handler_test.go`:
  - Fix `TestBootstrapHandler_Readyz_NotReady` to expect `"not_ready"` (underscore) and assert `message` field equals `"initializing"`.
  - Verify `Content-Type: application/json` on all responses.
  - Add test: readyz with `SetShuttingDown()` returns `"not_ready"` with message `"shutting_down"`.

### Acceptance criteria
- All tests in `internal/api/` pass.
- Response shapes match the design document exactly.
- `"not-ready"` string no longer appears in handler code or tests.
- `go vet ./internal/api/...` clean.
- `gofmt -l ./internal/api/` produces no output.

### Verification
```bash
docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'
```

### Commit message
```
feat(api): use named response structs, add reason to readyz 503

Replace map[string]string with typed structs for compile-time safety.
Readyz 503 now includes message field from ReadinessState.Reason().
Fix status value from "not-ready" to "not_ready" per contract.

FEATURE-0009
```

---

## Task 5: Add integration tests for method enforcement and shutdown

### Objective
Add integration-level tests verifying end-to-end method enforcement on bootstrap endpoints and shutdown readiness transition.

### Files
- `internal/server/server_test.go`

### Notes
- Add test: `POST /healthz` via full server handler returns 405 with `Allow: GET` header and APIErrorEnvelope body.
- Add test: `PUT /readyz` via full server handler returns 405 with `Allow: GET` header and APIErrorEnvelope body.
- Add test: `DELETE /version` via full server handler returns 405 with `Allow: GET` header and APIErrorEnvelope body.
- Add test: `HEAD /healthz` via full server handler returns 405 with `Allow: GET` header, `Content-Type: application/json`, and empty body.
- Add test: start server on random port, poll `/readyz` every 50ms until 200 (max 5s timeout), then verify `/healthz` also returns 200.
- Add test: after readyz returns 200, call `Shutdown()`, then verify `/readyz` returns 503 with `{"status":"not_ready","message":"shutting_down"}`.
- Use `newTestServer()` helper for httptest-style tests.
- For live server tests, use `net.Listen("tcp", "127.0.0.1:0")` to get a free port, then construct the server with that port.
- Use deterministic poll-until-ready (no fixed sleeps).
- Tests must pass under `go test -race ./internal/server/...`.

### Tests
Self-contained in this task (integration tests in `internal/server/server_test.go`).

### Acceptance criteria
- All new and existing server tests pass.
- No data races under `-race`.
- Method enforcement is verified end-to-end through the full middleware chain.
- Shutdown transition is verified with proper reason string.
- `go vet ./internal/server/...` clean.

### Verification
```bash
docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'
```

### Commit message
```
test(server): add integration tests for method enforcement and shutdown

Verify 405 responses through full bootstrapChain. Verify readyz
transitions from 503→200 on startup and 200→503 on shutdown with
correct reason strings.

FEATURE-0009
```

---

## Task 6: Update Makefile with VERSION and LDFLAGS

### Objective
Add build-time version injection to the Makefile via `-ldflags`.

### Files
- `Makefile`

### Notes
- Add variables near the top (after `APP_NAME` and `CONFIG`):
  ```makefile
  VERSION ?= dev
  MODULE  = github.com/sanjeevksaini/sovrunn
  LDFLAGS = -X '$(MODULE)/internal/api.buildVersion=$(VERSION)'
  ```
- Update the `build` target from:
  ```makefile
  build:
  	mkdir -p bin
  	go build -o bin/$(APP_NAME) ./cmd/sovrunn-api
  ```
  to:
  ```makefile
  build:
  	mkdir -p bin
  	go build -ldflags "$(LDFLAGS)" -o bin/$(APP_NAME) ./cmd/sovrunn-api
  ```
- Do NOT change any other Makefile targets.
- Do NOT add a separate `release` target.
- The injected variable is `internal/api.buildVersion` (already declared in `bootstrap_handler.go`).
- Default `VERSION=dev` means `go build` without override produces `"dev"` in `/version`.

### Tests
- No Go tests. Verify manually:
  ```bash
  make build
  ./bin/sovrunn-api --config configs/sovrunn-api.local.yaml &
  curl -s http://127.0.0.1:8080/version | grep '"version":"dev"'
  kill %1
  make build VERSION=0.1.0
  ./bin/sovrunn-api --config configs/sovrunn-api.local.yaml &
  curl -s http://127.0.0.1:8080/version | grep '"version":"0.1.0"'
  kill %1
  ```

### Acceptance criteria
- `make build` succeeds with ldflags.
- `make build VERSION=0.1.0` injects `0.1.0` into the binary.
- Default build produces `"dev"` in `/version` response.
- No other Makefile targets are broken.

### Verification
```bash
docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'
```

### Commit message
```
build: add VERSION ldflags injection to Makefile build target

VERSION defaults to "dev". Override with make build VERSION=x.y.z.
Injects into internal/api.buildVersion for /version endpoint.

FEATURE-0009
```

---

## Task 7: Final verification, guardrails, and cleanup

### Objective
Run full verification suite, enforce guardrails, clean artifacts, and confirm clean git status.

### Files
- No source file changes (verification-only task).

### Notes
- Run the final Docker verification command.
- Run guardrail checks.
- Clean up any build artifacts.
- Confirm no TODO(FEATURE-0009) markers remain under `internal/` or `cmd/`.
- Confirm `internal/api` does NOT import `internal/server`.
- Confirm git status is clean (no untracked or modified files after commit).

### Tests
All tests pass including race detection.

### Acceptance criteria

1. Final Docker verification passes:
```bash
docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./... && go build ./cmd/sovrunn-api'
```

2. Guardrails pass:
```bash
rm -f sovrunn-api
rm -rf bin
```

3. No TODO markers for this feature:
```bash
! grep -r 'TODO(FEATURE-0009)' internal/ cmd/
```

4. No forbidden import:
```bash
! grep -r '"github.com/sanjeevksaini/sovrunn/internal/server"' internal/api/
```

5. Git status clean:
```bash
git status --porcelain
```
Must produce no output.

### Verification
```bash
docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./... && go build ./cmd/sovrunn-api'
```

### Commit message
```
N/A — no source changes. Verification-only task.
```
