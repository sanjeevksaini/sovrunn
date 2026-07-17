# Implementation Plan: FEATURE-0005 Operation Resource

## Overview

Implement the Operation resource as the Phase 1 lifecycle/audit record for mutating control-plane
actions. Synchronous emission after successful create/update/delete; in-memory storage-only registry;
read-only API. No async engine, queue, workers, persistence, auth, or workflow framework.

## Cursor Execution Rule

Implement **one numbered task at a time**. After each task:
1. Run the verification commands listed in the task.
2. Confirm all pass with zero errors before moving to the next task.
3. **Do not start the next task** until the current task's verification succeeds.
4. If a task fails verification, fix within the current task before proceeding.

## Global Constraints

- Go 1.21-compatible only; no Go 1.22 wildcard routing; no external dependencies.
- `internal/api` MUST NOT import `internal/server`.
- OperationRegistry is storage-only and does NOT generate IDs.
- Empty `metadata.name` → `ErrMissingOperationID`; duplicate → `ErrAlreadyExists`.
- Emitter generates IDs via `crypto/rand`; retries `ErrAlreadyExists` up to 5 times.
- Phase 1 emits `Succeeded` only.
- Emission failure never changes the primary API response.
- No logger fields added to Org/OU/Tenant/Project handlers solely for emission.
- Existing handler tests must pass with a nil emitter.

## Tasks

- [ ] 1. Add Operation resource model
  - [ ] 1.1 Create Operation structs and constants
    **Files to create:** `internal/resources/operation.go`
    **Implementation details:**
    - `Operation` struct (APIVersion, Kind, Metadata, Spec, Status) with JSON tags
    - `OperationSpec`: Type, ResourceKind, ResourceName, OrganizationName/OrganizationUnitName/TenantName/ProjectName (omitempty), Actor, RequestID (omitempty)
    - `OperationStatus`: Phase, Message (omitempty), CreatedAt, UpdatedAt, CompletedAt (omitempty) — RFC3339 UTC strings
    - Constants: `OperationAPIVersion`, `OperationKind`
    - Phase constants: `OperationPhasePending/Running/Succeeded/Failed` (all defined; only Succeeded used in Phase 1)
    - 12 operation type constants (`OpCreateOrganization` … `OpDeleteProject`)
    - Confirm resource kind constants exist: `OrganizationKind`, `OrganizationUnitKind`, `TenantKind`, `ProjectKind`; add any missing one following the existing pattern
    - Reuse existing `Metadata` struct (Operation ID in `Name`)
    **Required tests:** None (covered by later layers)
    **Verification:** `make fmt && make vet`
    **Scope boundary:** Model only. No registry, emitter, handler, or wiring.
    **Do not start the next task until this passes.**

---

- [ ] 2. Add Operation registry
  - [ ] 2.1 Implement OperationRegistry (storage-only)
    **Files to create:** `internal/registry/operation_registry.go`
    **Files to modify:** `internal/registry/registry.go` (add `ErrMissingOperationID` sentinel)
    **Implementation details:**
    - `OperationRegistryIface`: CreateOperation, GetOperation, ListOperations (all ctx-first)
    - `OperationRegistry` struct: `sync.RWMutex` + `map[string]resources.Operation` keyed by ID
    - `NewOperationRegistry()` constructor; no package-level global state
    - `CreateOperation`: require non-empty `metadata.name` → else `ErrMissingOperationID` (no store); duplicate → `ErrAlreadyExists` (no overwrite); else store deep copy, return deep copy
    - `GetOperation`: RLock, deep copy, `ErrNotFound` if absent
    - `ListOperations`: RLock, non-nil slice of deep copies, sorted by `status.createdAt` asc then ID asc
    - `deepCopyOperation` helper (copy Labels/Annotations maps if present)
    - Reuse `ErrNotFound`, `ErrAlreadyExists`; storage-only; NO ID generation; no dependency on other registries
    **Required tests:** None in this task (Task 3)
    **Verification:** `make fmt && make vet`
    **Scope boundary:** Registry only. No emitter, handler, or wiring.
    **Do not start the next task until this passes.**

---

- [ ] 3. Add Operation registry tests
  - [ ] 3.1 Write registry unit and race tests
    **Files to create:**
    - `internal/registry/operation_registry_test.go`
    - `internal/registry/operation_registry_race_test.go`
    **Required tests:**
    - Create with empty `metadata.name` → `ErrMissingOperationID`, NOT stored (List shows no entry)
    - Create stores; duplicate ID → `ErrAlreadyExists`, original unchanged
    - Get by ID; missing → `ErrNotFound`
    - List sorted by createdAt then ID; empty → non-nil `[]`
    - Deep-copy immutability: mutating returned value/maps does not affect stored state
    - Race: 10+ goroutines mixed Create/Get/List, zero race reports
    **Verification:** `make fmt && make vet && make test && go test -race ./internal/registry/...`
    **Scope boundary:** Registry tests only.
    **Do not start the next task until this passes.**

---

- [ ] 4. Add Operation emitter
  - [ ] 4.1 Implement OperationEmitter interface, adapter, and ID generator
    **Files to create:** `internal/api/operation_emitter.go`
    **Implementation details:**
    - `OperationEmitter` interface: `Emit(ctx, spec resources.OperationSpec) error`
    - `newOperationID()` — `crypto/rand` 16 bytes, hex-encoded (32 chars); opaque, not DNS-validated
    - `registryEmitter` adapter holding `registry.OperationRegistryIface` and OPTIONAL `*log.Logger` (nil allowed)
    - `NewRegistryEmitter(reg, logger)` — logger may be nil
    - `Emit`: set Actor="system", timestamps (RFC3339 UTC), Phase=Succeeded, CompletedAt set; bounded retry (5) on `ErrAlreadyExists`; exhaustion → `errOperationIDExhausted`
    - `emitOperation(ctx, emitter, spec)` helper — nil-safe, swallows emitter errors, NO logger param
    - `requestIDFromContext(ctx)` — API-local helper reading the request ID context value; MUST NOT import `internal/server`
    **Required tests:** None in this task (Task 5)
    **Verification:** `make fmt && make vet`
    **Scope boundary:** Emitter + ID gen + request-ID helper only. No handler changes yet.
    **Do not start the next task until this passes.**

---

- [ ] 5. Add Operation emitter tests
  - [ ] 5.1 Write emitter unit tests
    **Files to create:** `internal/api/operation_emitter_test.go`
    **Required tests:**
    - `newOperationID` returns unique, non-empty hex tokens
    - Emit sets Actor=system, Phase=Succeeded, CreatedAt/UpdatedAt/CompletedAt populated
    - Collision retry: stub registry returns `ErrAlreadyExists` N (<5) times then success → Emit succeeds
    - Exhaustion: stub returns `ErrAlreadyExists` 5 times → Emit returns error
    - `emitOperation` with nil emitter is a no-op (no panic)
    - `emitOperation` swallows emitter errors (no error propagation)
    **Verification:** `make fmt && make vet && make test && go test -race ./internal/api/...`
    **Scope boundary:** Emitter tests only.
    **Do not start the next task until this passes.**

---

- [ ] 6. Add Operation HTTP handler
  - [ ] 6.1 Implement OperationHandler (read-only)
    **Files to create:** `internal/api/operation_handler.go`
    **Implementation details:**
    - `OperationHandler` holds `registry.OperationRegistryIface`; `NewOperationHandler(reg)`
    - `HandleCollection`: GET → List; POST → 405; any other method → 405
    - `HandleItem`: GET → Get; else 405
    - Path parsing (Go 1.21): `strings.TrimPrefix(r.URL.Path, "/v1/operations/")`; bare/empty → 404; `strings.Contains(remainder, "/")` (extra segments) → 404; id is opaque, no DNS validation
    - List: `{"items": items}` ([] when empty, only items field); registry error → 500 INTERNAL_ERROR
    - Get: 200 with full resource; `ErrNotFound` → 404 RESOURCE_NOT_FOUND; other error → 500
    - Reuse shared writeError/writeJSON
    **Required tests:** None in this task (Task 7)
    **Verification:** `make fmt && make vet`
    **Scope boundary:** Read handler only. No server/main wiring yet.
    **Do not start the next task until this passes.**

---

- [ ] 7. Add Operation handler tests
  - [ ] 7.1 Write handler tests
    **Files to create:** `internal/api/operation_handler_test.go`
    **Required tests (net/http/httptest):**
    - GET list 200 sorted; GET list 200 empty → `{"items": []}`
    - GET item 200 existing; GET item 404 missing
    - GET item 404 bare path (`/v1/operations/`)
    - GET item 404 extra segments (`/v1/operations/{name}/extra`) — explicit case
    - POST `/v1/operations` → 405
    **Verification:** `make fmt && make vet && make test && go test -race ./internal/api/...`
    **Scope boundary:** Handler tests only.
    **Do not start the next task until this passes.**

---

- [ ] 8. Wire Operation routes into server
  - [ ] 8.1 Register /v1/operations routes
    **Files to modify:**
    - `internal/server/server.go` — extend `New` to accept non-nil `*api.OperationHandler`; register `/v1/operations` → HandleCollection and `/v1/operations/` → HandleItem
    - `internal/server/server_test.go` — update constructor tests with an OperationHandler fixture; add route tests
    - `cmd/sovrunn-api/main.go` — minimal compile fix: create `operationRegistry`, `operationHandler`, pass into `server.New` (full emitter wiring in Task 9)
    **Implementation details:**
    - Keep middleware order: requestID → logging → contentType → handler
    - `internal/api` still must not import `internal/server`
    **Required tests:** server route tests for `/v1/operations` and `/v1/operations/`
    **Verification:** `make fmt && make vet && make test && make build`
    **Scope boundary:** Route wiring + minimal main.go compile fix only. No emitter injection into mutating handlers yet.
    **Do not start the next task until this passes.**

---

- [ ] 9. Inject nil-safe OperationEmitter into existing mutating handlers
  - [ ] 9.1 Add emitter dependency to Org/OU/Tenant/Project handlers
    **Files to modify:**
    - `internal/api/org_handler.go`, `internal/api/ou_handler.go`, `internal/api/tenant_handler.go`, `internal/api/project_handler.go` — add trailing `emitter OperationEmitter` field + constructor param (nil-safe); store emitter in each handler struct
    - Corresponding `*_handler_test.go` — update constructor call sites to pass `nil` emitter
    - `cmd/sovrunn-api/main.go` — only as a temporary compile fix if required (pass `nil` temporarily); full production emitter wiring happens in Task 10
    **Implementation details:**
    - Existing handler tests should pass `nil` for the emitter.
    - Do NOT emit Operations yet.
    - Do NOT remove FEATURE-0005 TODO markers yet.
    - Do NOT wire the real production emitter in main.go yet (nil temporarily if needed for compilation).
    - Do NOT add logger fields to handlers.
    - `internal/api` must not import `internal/server`.
    **Required tests:** existing handler tests pass unchanged (with nil emitter)
    **Verification:** `make fmt && make vet && make test && go test -race ./... && make build` (Docker fallback if host Go unavailable)
    **Scope boundary:** Dependency injection only. No emission calls. No marker removal. Do not start Task 10.
    **Do not start the next task until this passes.**

---

- [ ] 10. Wire Operation registry/emitter/handler in main.go
  - [ ] 10.1 Complete production wiring
    **Files to modify:** `cmd/sovrunn-api/main.go`
    **Implementation details:**
    - `operationRegistry := registry.NewOperationRegistry()`
    - `emitter := api.NewRegistryEmitter(operationRegistry, nil)` (or an existing logger if one already exists cleanly)
    - `operationHandler := api.NewOperationHandler(operationRegistry)`
    - Pass `emitter` into the OrgHandler, OUHandler, TenantHandler, and ProjectHandler constructors (params added in Task 9)
    - Pass `operationHandler` into `server.New`
    - Do NOT create duplicate operation registries
    - Do NOT change the `server.New` signature again unless a real compile issue requires it
    **Required tests:** existing tests still pass
    **Verification:** `make fmt && make vet && make test && make build` (Docker fallback if host Go unavailable)
    **Scope boundary:** main.go production wiring only. No emission calls. No marker removal. Do not start Task 11.
    **Do not start the next task until this passes.**

---

- [ ] 11. Replace FEATURE-0005 TODO markers with Operation emission
  - [ ] 11.1 Emit Operations from all 12 mutating paths
    **Files to modify:**
    - `internal/api/org_handler.go`, `ou_handler.go`, `tenant_handler.go`, `project_handler.go`
    **Implementation details:**
    - At each `// TODO(FEATURE-0005): emit Operation record` marker (after the registry mutation
      succeeds, before writing the response), call `emitOperation(ctx, h.emitter, spec)`
    - Populate `OperationSpec` per the design table: Type constant, ResourceKind constant
      (`resources.OrganizationKind`, etc.), ResourceName, parent-ref fields, RequestID via
      `requestIDFromContext(ctx)`
    - Create/Update: use created/updated resource fields; Delete: use path segments
    - Remove every TODO marker comment
    - Emission must never alter the primary response (201/200/204)
    **Required tests:** existing handler tests still pass (nil emitter path)
    **Verification:** `make fmt && make vet && make test && go test -race ./... && make build`
    **Scope boundary:** Emission calls + marker removal only. Integration tests in Task 12.
    **Do not start the next task until this passes.**

---

- [ ] 12. Add table-driven emission integration tests
  - [ ] 12.1 Verify emission across all resources and failure isolation
    **Files to create/modify:**
    - `internal/api/operation_emission_test.go` (or extend existing handler tests)
    **Required tests:**
    - Table-driven across Organization, OrganizationUnit, Tenant, Project × create/update/delete:
      after a successful action, exactly one Operation is recorded with the correct Type,
      ResourceKind, ResourceName, parent-ref fields, `actor = "system"`, `status.phase = "Succeeded"`
    - Use a real `OperationRegistry` + `registryEmitter` (or a capturing stub emitter) to assert emission
    - Failure isolation: with a stub emitter that always errors, the primary response is unchanged
      (still 201/200/204) and no client error is returned
    - No Operation is emitted when the mutating action itself fails (validation/duplicate/not-found)
    **Note:** Implement this task LAST among feature tasks — it depends on the model, registry,
    emitter, and handlers all being stable.
    **Verification:** `make fmt && make vet && make test && go test -race ./...`
    **Scope boundary:** Integration tests only.
    **Do not start the next task until this passes.**

---

- [ ] 13. Final verification and cleanup
  - [ ] 13.1 Run full verification and marker check
    **Verification commands (host Go):**
    ```
    go fmt ./...
    go vet ./...
    go test ./...
    go test -race ./...
    go build -o bin/sovrunn-api ./cmd/sovrunn-api
    ```
    **Docker fallback (if host Go unavailable):**
    ```
    docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'go fmt ./... && go vet ./... && go test ./... && go test -race ./... && go build -o bin/sovrunn-api ./cmd/sovrunn-api'
    ```
    **TODO marker check (MUST return zero matches):**
    ```
    ! grep -R "TODO(FEATURE-0005): emit Operation record" internal/api
    ```
    **Cleanup:**
    ```
    rm -rf bin
    git status
    ```
    **Completion criteria:**
    - All checks pass; zero race reports; binary builds
    - Zero matches for the TODO marker (Requirement 5.7)
    - `internal/api` does not import `internal/server`
    - Working tree clean except intended changes; no binary committed

## Notes

- Module path: `github.com/sanjeevksaini/sovrunn`
- Go 1.21 minimum — no Go 1.22 wildcard routing
- No new external dependencies
- Correctness Properties (optional, from design.md) tagged `// Feature: operation-resource, Property N: <title>`

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1"] },
    { "id": 1, "tasks": ["2.1"] },
    { "id": 2, "tasks": ["3.1"] },
    { "id": 3, "tasks": ["4.1"] },
    { "id": 4, "tasks": ["5.1"] },
    { "id": 5, "tasks": ["6.1"] },
    { "id": 6, "tasks": ["7.1"] },
    { "id": 7, "tasks": ["8.1"] },
    { "id": 8, "tasks": ["9.1"] },
    { "id": 9, "tasks": ["10.1"] },
    { "id": 10, "tasks": ["11.1"] },
    { "id": 11, "tasks": ["12.1"] },
    { "id": 12, "tasks": ["13.1"] }
  ]
}
```
