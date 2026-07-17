# Implementation Plan: FEATURE-0004 Project Resource and Registry

## Overview

Implement Project as the workload/environment grouping boundary under Tenant.
Composite identity: `organizationName/organizationUnitName/tenantName/name`. Storage-only registry.
Parent Tenant lookup in API layer. Tenant delete extended with Project blocker. Go 1.21 stdlib only.

## Cursor Execution Rule

Implement **one numbered task at a time**. After each task:
1. Run the verification commands listed in the task.
2. Confirm all pass with zero errors before moving to the next task.
3. Do not begin a subsequent task until the current task's verification succeeds.
4. If a task fails verification, fix within the current task before proceeding.

## Tasks

- [ ] 1. Add Project resource model
  - [ ] 1.1 Create Project structs and constants
    **Objective:** Define Project, ProjectSpec, ProjectStatus structs and constants.
    **Files to create:**
    - `internal/resources/project.go`
    **Implementation notes:**
    - Reuse existing `Metadata` struct
    - `ProjectSpec`: `OrganizationName`, `OrganizationUnitName`, `TenantName`, `Description` (omitempty)
    - `ProjectStatus`: `Phase`, `Message` (omitempty)
    - Constants: `ProjectAPIVersion = "platform.sovrunn.io/v1alpha1"`, `ProjectKind = "Project"`
    - Reuse existing phase constants
    **Tests required:** None (tested via registry/handler tests)
    **Completion criteria:** `make fmt` passes; package compiles
    **Verification:** `make fmt && make vet`

---

- [ ] 2. Add Project validation
  - [ ] 2.1 Implement ValidateProject and ValidateProjectPathSegments with tests
    **Objective:** Pure validation functions for Project fields and path segments.
    **Files to create:**
    - `internal/validation/project.go`
    - `internal/validation/project_test.go`
    - `internal/validation/project_property_test.go`
    **Implementation notes:**
    - `ValidateProject(p resources.Project) []resources.FieldError` — context-free
    - `ValidateProjectPathSegments(orgName, ouName, tenantName, name string) []resources.FieldError` — context-free
    - Rules: metadata.name, spec.organizationName, spec.organizationUnitName, spec.tenantName each DNS-label 1–63
    - Reuse `dnsLabelRe` and `validateName` from validation package
    - Path field mapping: orgName → spec.organizationName, ouName → spec.organizationUnitName, tenantName → spec.tenantName, name → metadata.name
    - No I/O, no registry lookup
    **Tests required:**
    - Unit: valid names accepted; empty/invalid/long metadata.name, spec.organizationName, spec.organizationUnitName, spec.tenantName rejected; path validation maps fields correctly
    - Property (testing/quick, Config{MaxCount: 100}): P1 valid DNS-labels accepted; P2 invalid strings rejected
    **Completion criteria:** All validation tests pass; `make fmt && make vet && make test` pass
    **Verification:** `make fmt && make vet && make test`

---

- [ ] 3. Add Project registry
  - [ ] 3.1 Implement ProjectRegistry with unit, property, and race tests
    **Objective:** Thread-safe in-memory registry with four-part composite key and deep copies.
    **Files to create:**
    - `internal/registry/project_registry.go`
    - `internal/registry/project_registry_test.go`
    - `internal/registry/project_registry_property_test.go`
    - `internal/registry/project_registry_race_test.go`
    **Implementation notes:**
    - `ProjectRegistryIface` with 6 methods (Create, Get, List, Update, Delete, CountByTenant)
    - `CreateProject(ctx, p) (resources.Project, error)` — returns stored deep copy
    - `UpdateProject(ctx, p) (resources.Project, error)` — derives composite key from submitted resource; preserves stored APIVersion, Kind, Status, Metadata.Name, Spec.OrganizationName, Spec.OrganizationUnitName, Spec.TenantName; replaces only Metadata.DisplayName, Metadata.Labels, Metadata.Annotations, Spec.Description; returns deep copy
    - `projectCompositeKey(orgName, ouName, tenantName, name)` → joined with "/"
    - `sync.RWMutex`: RLock for Get/List/Count, Lock for Create/Update/Delete
    - `deepCopyProject` duplicates Labels and Annotations maps
    - List sorted: organizationName → organizationUnitName → tenantName → name, all ascending
    - Define `TenantLookup` interface here
    - Reuse `ErrNotFound`, `ErrAlreadyExists`; no global state; no dependency on other registries
    **Tests required:**
    - Unit: Create stores; duplicate → ErrAlreadyExists (original unchanged); same name under different Tenants succeeds; Get by composite key; Get non-existent → ErrNotFound; List sorted; empty List → non-nil empty slice; Update mutable fields only; Update non-existent → ErrNotFound; Delete removes; Delete non-existent → ErrNotFound; CountByTenant correct
    - Property (testing/quick): P3 round-trip; P4 four-level sort invariant; P5 deep-copy immutability; P6 duplicate idempotent error
    - Race: 10+ goroutines mixed Create/Get/List/Update/Delete/CountByTenant
    **Completion criteria:** All tests pass; `make fmt && make vet && make test && go test -race ./internal/registry/...` pass
    **Verification:** `make fmt && make vet && make test && go test -race ./internal/registry/...`

---

- [ ] 4. Add Project child blocker for Tenant delete
  - [ ] 4.1 Implement ProjectChildBlockerChecker and TenantChildBlocker interface
    **Objective:** Blocker that prevents Tenant deletion when Projects reference it.
    **Files to create:**
    - `internal/registry/project_blocker.go`
    - `internal/registry/project_blocker_test.go`
    **Implementation notes:**
    - Define `TenantChildBlocker` interface: `BlockedByTenantChildren(ctx, orgName, ouName, tenantName string) ([]BlockedBy, error)`
    - `ProjectChildBlockerChecker` holds `ProjectRegistryIface`
    - `NewProjectChildBlockerChecker(reg) *ProjectChildBlockerChecker`
    - `BlockedByTenantChildren` calls `CountByTenant`; count > 0 → `[]BlockedBy{{Kind: "Project", Count: count}}`; count 0 → nil; propagate registry error
    **Tests required:**
    - Zero projects → no blockers
    - One+ projects → blocking kind "Project" with correct count
    - Registry error propagates
    **Completion criteria:** `make fmt && make vet && make test` pass
    **Verification:** `make fmt && make vet && make test`

---

- [ ] 5. Extend TenantHandler with optional TenantChildBlocker
  - [ ] 5.1 Add blocker dependency to TenantHandler
    **Objective:** Wire the Project blocker into the Tenant delete path.
    **Files to modify:**
    - `internal/api/tenant_handler.go`
    - `internal/api/tenant_handler_test.go` (and any existing NewTenantHandler call sites)
    **Implementation notes:**
    - Add `blocker registry.TenantChildBlocker` field to TenantHandler
    - Update `NewTenantHandler(reg, ouLookup, blocker)` signature
    - Allow `blocker` nil → Tenant delete proceeds as before (FEATURE-0003 behavior)
    - If `blocker` non-nil, DELETE consults `BlockedByTenantChildren` before deleting
    - If blockers returned → HTTP 409 DELETE_BLOCKED, message identifying "Project"
    - Do NOT create generic blocker framework; do NOT rewrite OUHandler or OrgHandler
    - Update existing tests impacted by the new NewTenantHandler signature
    **Tests required:**
    - nil blocker keeps existing delete behavior (204)
    - blocker returning Project blocks Tenant delete → 409 DELETE_BLOCKED
    - blocker returning empty allows delete → 204
    **Completion criteria:** `make fmt && make vet && make test && go test -race ./internal/api/...` pass
    **Verification:** `make fmt && make vet && make test && go test -race ./internal/api/...`

---

- [ ] 6. Add Project safe JSON decoder
  - [ ] 6.1 Implement safeDecodeProject
    **Objective:** Safe JSON decode function for Project request bodies.
    **Files to create:**
    - `internal/api/project_decode.go`
    **Implementation notes:**
    - `http.MaxBytesReader(w, r.Body, 1<<20)` inside the function
    - Read body; `*http.MaxBytesError` → errBodyTooLarge; empty → errEmptyBody
    - Decode into `map[string]json.RawMessage`; "status" key present → errStatusFieldPresent
    - Typed decode with `DisallowUnknownFields()`; unknown → errUnknownField; syntax/type → errMalformedJSON
    - Do not echo raw body
    - 415 remains contentTypeMiddleware responsibility, not the decoder
    - Reuse existing error sentinels from the FEATURE-0003 decoder pattern
    **Tests required:** Tested via handler tests in Task 7
    **Completion criteria:** `make fmt && make vet` pass; function compiles
    **Verification:** `make fmt && make vet`

---

- [ ] 7. Add Project HTTP handler
  - [ ] 7.1 Implement ProjectHandler with all CRUD methods and tests
    **Objective:** HTTP handler for Project CRUD endpoints.
    **Files to create:**
    - `internal/api/project_handler.go`
    - `internal/api/project_handler_test.go`
    **Implementation notes:**
    - `ProjectHandler` takes `ProjectRegistryIface` + `TenantLookup`
    - `NewProjectHandler(reg, tenantLookup)`
    - `HandleCollection`: POST → Create, GET → List, else 405
    - `HandleItem`: `strings.TrimPrefix` + `strings.Split`; exactly 4 non-empty segments else 404; invalid DNS segment → 400; dispatch GET/PUT/DELETE, else 405
    - Path field mapping: orgName → spec.organizationName, ouName → spec.organizationUnitName, tenantName → spec.tenantName, name → metadata.name
    - Create: decode → validate → parent Tenant lookup (missing → 400 field=spec.tenantName, message includes full parent ref) → set apiVersion/kind/status.phase=Active → registry.CreateProject → duplicate → 409 → `// TODO(FEATURE-0005): emit Operation record — type: CreateProject` → 201
    - Get: validate path → registry.GetProject → 404 → 200
    - List: registry.ListProjects → `{"items": items}` ([] when empty, only items field)
    - Update: validate path → decode → require body metadata.name/spec.organizationName/spec.organizationUnitName/spec.tenantName present and matching path → validate → registry.UpdateProject(ctx, project) → 404 → `// TODO(FEATURE-0005): emit Operation record — type: UpdateProject` → 200
    - Delete: validate path → registry.DeleteProject → 404 → `// TODO(FEATURE-0005): emit Operation record — type: DeleteProject` → 204 no body
    - Reuse shared writeError/writeJSON/writeValidationErrors
    **Tests required:**
    - POST 201 valid; POST 409 duplicate; POST 400 (invalid fields, non-existent parent, status key, bad JSON, unknown field); POST 413 oversized body
    - GET 200/404/400; GET wrong path shape → 404
    - List sorted; list empty → []
    - PUT 200/404; PUT 400 (name mismatch, orgName mismatch, ouName mismatch, tenantName mismatch, status key)
    - DELETE 204/404/400
    **Completion criteria:** `make fmt && make vet && make test && go test -race ./internal/api/...` pass
    **Verification:** `make fmt && make vet && make test && go test -race ./internal/api/...`

---

- [ ] 8. Wire Project routes into server
  - [ ] 8.1 Update server.New and route registration
    **Objective:** Register Project routes and require ProjectHandler.
    **Files to modify:**
    - `internal/server/server.go`
    - `internal/server/server_test.go`
    - `cmd/sovrunn-api/main.go` (minimal compile fix only — see note)
    **Implementation notes:**
    - Update `server.New` to accept non-nil `*api.ProjectHandler`
    - Register `/v1/projects` → HandleCollection, `/v1/projects/` → HandleItem
    - Keep middleware order: requestID → logging → contentType → handler
    - Update tests constructing `server.New` to provide a ProjectHandler fixture
    - **Minimal main.go compile fix only:** create `projectRegistry := registry.NewProjectRegistry()`, `projectHandler := api.NewProjectHandler(projectRegistry, tenantRegistry)`, pass `projectHandler` into `server.New`. Do NOT wire the ProjectChildBlocker into TenantHandler here unless required for compilation — Task 9 completes blocker wiring.
    **Completion criteria:** `make fmt && make vet && make test && make build` pass
    **Verification:** `make fmt && make vet && make test && make build`

---

- [ ] 9. Wire Project registry, handler, and blocker in main.go
  - [ ] 9.1 Complete production wiring
    **Objective:** Construct and inject all Project dependencies, including the Tenant blocker.
    **Files to modify:**
    - `cmd/sovrunn-api/main.go`
    **Implementation notes:**
    ```
    projectRegistry := registry.NewProjectRegistry()
    projectBlocker := registry.NewProjectChildBlockerChecker(projectRegistry)
    tenantHandler := api.NewTenantHandler(tenantRegistry, ouRegistry, projectBlocker)
    projectHandler := api.NewProjectHandler(projectRegistry, tenantRegistry)
    srv := server.New(cfg, orgHandler, ouHandler, tenantHandler, projectHandler, bootstrapHandler, readiness)
    ```
    - Update any remaining NewTenantHandler call sites for the new signature
    **Completion criteria:** `make fmt && make vet && make test && make build` pass; `make run` reachable at `/v1/projects`
    **Verification:** `make fmt && make vet && make test && make build`

---

- [ ] 10. Add integration tests for Tenant delete blocked by Project
  - [ ] 10.1 Write Tenant-delete-blocked integration tests
    **Objective:** Verify end-to-end Tenant deletion blocking via Project.
    **Files to create/modify:**
    - `internal/api/project_handler_test.go`, `internal/api/tenant_handler_test.go`, or `internal/server/server_test.go`
    **Tests required:**
    - Create org → OU → Tenant → Project under Tenant → DELETE Tenant → 409 DELETE_BLOCKED with "Project" in message
    - Create org → OU → Tenant → no Project → DELETE Tenant → 204
    **Completion criteria:** `make fmt && make vet && make test` pass
    **Verification:** `make fmt && make vet && make test`

---

- [ ] 11. Final verification
  - [ ] 11.1 Run the full verification suite
    **Objective:** Confirm all checks pass and working tree is clean.
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
    **Cleanup:**
    ```
    rm -rf bin
    git status
    ```
    **Completion criteria:**
    - All checks pass
    - Working tree clean except intended changes
    - No generated binary committed
    - No unrelated files modified

## Notes

- Property tests use `testing/quick` with `Config{MaxCount: 100}`
- Each property test tagged: `// Feature: project-resource, Property N: <title>`
- Module path: `github.com/sanjeevksaini/sovrunn`
- Go 1.21 minimum — no Go 1.22 wildcard routing
- No new external dependencies required

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1"] },
    { "id": 1, "tasks": ["2.1"] },
    { "id": 2, "tasks": ["3.1"] },
    { "id": 3, "tasks": ["4.1"] },
    { "id": 4, "tasks": ["5.1", "6.1"] },
    { "id": 5, "tasks": ["7.1"] },
    { "id": 6, "tasks": ["8.1"] },
    { "id": 7, "tasks": ["9.1"] },
    { "id": 8, "tasks": ["10.1"] },
    { "id": 9, "tasks": ["11.1"] }
  ]
}
```
