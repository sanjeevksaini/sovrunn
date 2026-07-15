# Implementation Plan: FEATURE-0001 Organization Resource and Registry

## Overview

Implement the Organization resource as the top-level governance boundary in Sovrunn, bootstrapping the full Go project skeleton with in-memory registry, validation, HTTP handlers, middleware, and graceful shutdown. All code uses Go 1.21 stdlib `net/http` only (no third-party router). The single permitted external dependency is `gopkg.in/yaml.v3` for config parsing.

Implement FEATURE-0001 strictly from the approved design.md and tasks.md.

Follow the Cursor Execution Rule:
- Implement one numbered task at a time.
- After each task, run the verification commands listed for that task.
- Do not start the next numbered task until the current task passes verification.
- Do not implement future resources or future Phase 1 features early.
- Stop and report if design.md and tasks.md conflict.

## Tasks

- [ ] 1. Go module initialization, Makefile verification, and minimal main.go skeleton
  - [ ] 1.1 Initialize Go module and create minimal compilable main.go
    - Create `go.mod` with module path `github.com/sanjeevksaini/sovrunn` and Go 1.21
    - Create `cmd/sovrunn-api/main.go` with a minimal `main()` that prints "sovrunn-api starting" and exits (placeholder for wiring)
    - Do NOT create doc.go placeholders — packages will be created naturally when their implementation tasks are reached
    - Create `tests/integration/` directory as a placeholder for future integration tests (empty or with a single .gitkeep)
    - Verify: `make fmt`, `make vet`, `make test`, `make build` all pass with exit code 0
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 2.1_

- [ ] 2. Resource types, error types, and constants
  - [ ] 2.1 Implement internal/resources package
    - Create `internal/resources/organization.go`: `Organization`, `Metadata`, `OrganizationSpec`, `OrganizationStatus` structs with exact JSON tags from design.md; phase constants (`PhaseActive`, `PhaseInactive`, `PhaseDeleting`, `PhaseFailed`); API version and kind constants (`OrgAPIVersion`, `OrgKind`)
    - Create `internal/resources/errors.go`: `ErrorCode` type, all 5 error code constants, `APIError`, `APIErrorEnvelope`, `FieldError` structs
    - Verify: `make fmt`, `make vet`, `make test` pass
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 13.1, 13.2_

- [ ] 3. Configuration loading
  - [ ] 3.1 Implement internal/config package
    - Create `internal/config/config.go`: `Config`, `ServerConfig`, `LogConfig`, `RegistryConfig` structs with YAML tags; `Load(path string) (Config, error)` function using `gopkg.in/yaml.v3`; `Addr()` method
    - Add `gopkg.in/yaml.v3` to `go.mod` and run `go mod tidy`
    - Validate required fields in `Load()`: port > 0, host non-empty; return error on invalid/missing config
    - Ensure `configs/sovrunn-api.local.yaml` matches the expected schema
    - Write unit tests in `internal/config/config_test.go`: valid config loads, missing file errors, invalid port errors, missing host errors
    - Verify: `make fmt`, `make vet`, `make test` pass
    - _Requirements: 2.8, 17.1, 17.2, 17.3_

- [ ] 4. Health/readiness state
  - [ ] 4.1 Implement internal/health package
    - Create `internal/health/readiness.go`: `ReadinessState` struct with `atomic.Bool`, `SetReady(bool)`, `IsReady() bool` methods
    - Write unit tests in `internal/health/readiness_test.go`: default is not ready, SetReady(true) makes IsReady() return true, SetReady(false) reverts
    - Verify: `make fmt`, `make vet`, `make test` pass
    - _Requirements: 4.2, 4.3_

- [ ] 5. Validation logic with property tests
  - [ ] 5.1 Implement internal/validation package with unit and property tests
    - Create `internal/validation/organization.go`: package-level `dnsLabelRe = regexp.MustCompile(...)`, `ValidateOrganization(ctx, org) []FieldError`, `ValidateNamePath(ctx, name) []FieldError`
    - Validation rules: empty name → error; len > 63 → error; regex mismatch → error; valid name → no errors
    - Both functions accept `context.Context` as first parameter
    - Create `internal/validation/organization_test.go`: unit tests for valid names (a, a1, a-b, 63-char max), invalid names (empty, uppercase, spaces, leading hyphen, trailing hyphen, >63 chars, special chars, single hyphen)
    - Create `internal/validation/organization_property_test.go`: Property 1 and Property 2 tests using `testing/quick` with `Config{MaxCount: 100}`
    - Property 1: generate invalid names (empty, uppercase, spaces, >63, leading/trailing hyphen) → at least one FieldError with Field="metadata.name"
    - Property 2: generate valid DNS-label names (1–63 chars matching regex) → no FieldError with Field="metadata.name"
    - Each property test function has comment: `// Feature: organization-resource-registry, Property N: <title>`
    - Verify: `make fmt`, `make vet`, `make test`, `go test -race ./internal/validation/...` pass
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 15.1, 18.1_

- [ ] 6. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 7. In-memory registry with property tests and race test
  - [ ] 7.1 Implement internal/registry package with unit, property, and concurrency tests
    - Create `internal/registry/registry.go`: `OrganizationRegistryIface` interface (5 methods with `context.Context` first param), `ErrNotFound`, `ErrAlreadyExists` sentinel errors
    - Create `internal/registry/org_registry.go`: `OrganizationRegistry` struct (`sync.RWMutex` + `map[string]resources.Organization`), `NewOrganizationRegistry()`, `deepCopyOrganization()`, all 5 CRUD methods with deep copies on every path
    - `UpdateOrganization` signature: `(ctx, name string, org Organization) (Organization, error)` — preserves metadata.name, status, apiVersion, kind from stored entry
    - Create `internal/registry/blocker.go`: `BlockedBy` struct, `ChildBlockerChecker` interface, `NoopChildBlockerChecker` struct
    - Create `internal/registry/org_registry_test.go`: unit tests for all scenarios in design.md Testing Strategy table (create stores, duplicate → ErrAlreadyExists, get exists, get missing → ErrNotFound, list empty, list sorted, update mutable fields, update missing → ErrNotFound, delete exists, delete missing → ErrNotFound)
    - Create `internal/registry/org_registry_property_test.go`: Property 3 (create/get round trip), Property 4 (value copy immutability), Property 5 (list sorted order), Property 6 (update preserves system fields) — all using `testing/quick` with `Config{MaxCount: 100}`
    - Create `internal/registry/org_registry_race_test.go`: stress test launching 10+ goroutines performing concurrent mixed CRUD operations on the same registry instance
    - Each property test function has comment: `// Feature: organization-resource-registry, Property N: <title>`
    - Verify: `make fmt`, `make vet`, `make test`, `go test -race ./internal/registry/...` pass
    - _Requirements: 2.4, 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7, 7.8, 7.9, 12.4, 12.5, 15.1, 18.2_

- [ ] 8. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 9. HTTP handlers (OrgHandler, BootstrapHandler, helpers)
  - [ ] 9.1 Implement response helpers and safe decode logic
    - Create `internal/api/response.go`: `writeError(w, r, status, code, message, field, details)` and `writeJSON(w, r, status, v)` helper functions; always set Content-Type and X-Sovrunn-Request-ID header
    - Create `internal/api/decode.go`: `safeDecodeOrganization(w, r) (Organization, error)` — applies `http.MaxBytesReader(w, r.Body, 1<<20)` inside, reads body bytes, checks for "status" key via `map[string]json.RawMessage` decode (reject if key present regardless of value), then typed decode with `DisallowUnknownFields()`; sentinel error types for caller mapping (errBodyTooLarge, errStatusFieldPresent, errMalformedJSON, errEmptyBody, errUnknownField)
    - Verify: `make fmt`, `make vet` pass (no tests yet — helpers tested via handler tests)
    - _Requirements: 13.1, 14.1, 14.2, 14.3, 14.4, 14.5_
  - [ ] 9.2 Implement BootstrapHandler
    - Create `internal/api/bootstrap_handler.go`: `BootstrapHandler` struct (takes `Config` + `*ReadinessState`), `NewBootstrapHandler()`, `Healthz`, `Readyz`, `Version` handlers; version handler returns `name`, `version`, `phase`, `status` fields
    - Create `internal/api/bootstrap_handler_test.go`: tests for healthz → 200, readyz → 200 when ready, readyz → 503 when not ready, version → 200 with expected fields
    - Verify: `make fmt`, `make vet`, `make test`, `go test -race ./internal/api/...` pass
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7_
  - [ ] 9.3 Implement OrgHandler (Create, Get, List, Update, Delete)
    - Create `internal/api/org_handler.go`: `OrgHandler` struct (takes `OrganizationRegistryIface` + `ChildBlockerChecker` interfaces), `NewOrgHandler()`, `HandleCollection(w, r)` dispatching POST→Create and GET→List, `HandleItem(w, r)` dispatching GET→Get, PUT→Update, DELETE→Delete; name extraction via `strings.TrimPrefix(r.URL.Path, "/v1/organizations/")`; bare `/v1/organizations/` with empty name → 404
    - Implement Create, Get, List, Update, Delete per design.md handler pseudocode exactly
    - PUT handler: require body.Metadata.Name present and equal to path name; empty/absent → 400
    - Add `// TODO(FEATURE-0005): emit Operation record — type: CreateOrganization` (and Update/Delete) comments
    - Verify: `make fmt`, `make vet` pass
    - _Requirements: 2.6, 5.5, 5.6, 5.7, 8.1–8.8, 9.1–9.5, 10.1–10.5, 11.1–11.7, 12.1–12.3, 15.2, 16.1–16.5, 19.1–19.4_
  - [ ] 9.4 Implement OrgHandler tests
    - Create `internal/api/org_handler_test.go`: HTTP handler tests using `net/http/httptest` covering all scenarios: POST 201 valid, POST 409 duplicate, POST 400 invalid name, POST 400 status field present ({}, null, {"phase":""}), POST 400 bad JSON, POST 413 oversized body, POST 415 wrong content-type, GET 200 exists (full shape), GET 404 missing, GET 400 invalid path name, GET list 200 sorted items, GET list 200 empty, PUT 200 valid update, PUT 404 missing, PUT 400 name mismatch, PUT 400 name absent in body, PUT 400 status field, DELETE 204 success, DELETE 404 missing, DELETE 400 invalid path name
    - Verify: `make fmt`, `make vet`, `make test`, `go test -race ./internal/api/...` pass
    - _Requirements: 13.3–13.7, 18.3_

- [ ] 10. Server lifecycle, middleware, and route registration
  - [ ] 10.1 Implement internal/server package
    - Create `internal/server/middleware.go`: `requestIDMiddleware` (crypto/rand 16-byte hex ID, reads/writes X-Sovrunn-Request-ID, stores in context), `contentTypeMiddleware` (rejects non-application/json Content-Type on POST/PUT/PATCH with 415), `loggingMiddleware` (captures status code, logs request_id/method/path/status_code/latency_ms, error_code on failure)
    - Create `internal/server/server.go`: `Server` struct, `New(cfg, orgHandler, bootstrapHandler, readiness)` — creates ServeMux, registers routes with middleware chain `requestID → contentType → logging`, registers bootstrap routes without contentType check; `Start()` — binds listener, sets readiness true, blocks on signal; `Shutdown(timeout)` — graceful drain
    - Route registration (Go 1.21 compatible): `mux.Handle("/v1/organizations", chain(orgHandler.HandleCollection))`, `mux.Handle("/v1/organizations/", chain(orgHandler.HandleItem))`, `mux.HandleFunc("/healthz", ...)`, `mux.HandleFunc("/readyz", ...)`, `mux.HandleFunc("/version", ...)`
    - Create `internal/server/middleware_test.go`: tests for request ID generation when absent, request ID propagation when present, content-type rejection for POST without application/json, content-type pass-through for GET, logging middleware writes structured log line
    - Verify: `make fmt`, `make vet`, `make test`, `go test -race ./internal/server/...` pass
    - _Requirements: 2.7, 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 15.3, 15.4, 16.4, 18.4_

- [ ] 11. Checkpoint — Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 12. Application wiring and integration verification
  - [ ] 12.1 Wire cmd/sovrunn-api/main.go and run full verification
    - Update `cmd/sovrunn-api/main.go`: parse `--config` flag (default `configs/sovrunn-api.local.yaml`), call `config.Load()`, construct `registry.NewOrganizationRegistry()`, construct `api.NewOrgHandler(reg, registry.NoopChildBlockerChecker{})`, construct `api.NewBootstrapHandler(cfg, &readiness)`, construct `server.New(cfg, orgHandler, bootstrapHandler, &readiness)`, call `server.Start()`
    - No business logic in main.go — wiring only
    - Verify full pipeline: `make fmt`, `make vet`, `make test`, `go test -race ./...`, `make build` all pass
    - Verify runtime: `make run` starts the server, then manually confirm `curl http://127.0.0.1:8080/healthz` → `{"status":"ok"}`, `curl http://127.0.0.1:8080/readyz` → `{"status":"ready"}`, `curl http://127.0.0.1:8080/version` → version JSON
    - _Requirements: 1.3, 1.4, 1.5, 1.6, 2.2, 3.1, 3.2, 17.1, 17.2, 17.3, 17.4, 17.5, 17.6_

- [ ] 13. Final checkpoint — Implementation summary and verification report
  - Run: `make fmt`, `make vet`, `make test`, `go test -race ./...`, `make build`
  - Verify no formatting diffs, no vet errors, no test failures, no race conditions, binary built successfully
  - Produce a final implementation report listing:
    - Files created
    - New exported symbols
    - Validation rules implemented
    - Tests written (unit + property counts)
    - Security considerations addressed
    - Known limitations
    - Non-goals intentionally not implemented
    - All verification commands run and their results
    - Acceptance criteria satisfied per requirements.md

## Cursor Execution Rule

Implement **one numbered task at a time**. After each task:
1. Run the verification commands listed in the task (`make fmt`, `make vet`, `make test`, etc.)
2. Confirm all pass with zero errors before moving to the next task.
3. Do not begin a subsequent task until the current task's verification commands succeed.
4. If a task fails verification, fix the issue within the current task before proceeding.

## Notes

- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests (Properties 1–6) validate universal correctness properties and are bundled with the code they test
- Unit tests validate specific examples and edge cases
- The concurrency stress test ensures zero data races under `go test -race`
- The sole external dependency is `gopkg.in/yaml.v3` for YAML config parsing
- Module path: `github.com/sanjeevksaini/sovrunn`
- Go 1.21 minimum — no Go 1.22 wildcard routing patterns
- `safeDecodeOrganization` handles body size limiting internally (no bodyLimitMiddleware)
- `UpdateOrganization` returns `(resources.Organization, error)` — handler uses returned value directly

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1"] },
    { "id": 1, "tasks": ["2.1", "4.1"] },
    { "id": 2, "tasks": ["3.1", "5.1"] },
    { "id": 3, "tasks": ["7.1"] },
    { "id": 4, "tasks": ["9.1"] },
    { "id": 5, "tasks": ["9.2", "9.3"] },
    { "id": 6, "tasks": ["9.4"] },
    { "id": 7, "tasks": ["10.1"] },
    { "id": 8, "tasks": ["12.1"] }
  ]
}
```
