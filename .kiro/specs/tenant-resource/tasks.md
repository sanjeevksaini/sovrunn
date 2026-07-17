# Implementation Plan: FEATURE-0003 Tenant Resource and Registry

## Overview

Implement Tenant as the primary isolation boundary under OrganizationUnit.
Composite identity: `organizationName/organizationUnitName/name`. Storage-only registry.
Parent OU lookup in API layer. OU delete extended with Tenant blocker. Go 1.21 stdlib only.

## Cursor Execution Rule

Implement **one numbered task at a time**. After each task:
1. Run the verification commands listed in the task.
2. Confirm all pass with zero errors before moving to the next task.
3. Do not begin a subsequent task until the current task's verification succeeds.
4. If a task fails verification, fix within the current task before proceeding.

## Tasks

- [ ] 1. Add Tenant resource model
  - [ ] 1.1 Create Tenant structs and constants
    **Objective:** Define Tenant, TenantSpec, TenantStatus structs and constants.
    **Files to create:**
    - `internal/resources/tenant.go`
    **Implementation notes:**
    - Reuse existing `Metadata` struct
    - `TenantSpec`: `OrganizationName` (`json:"organizationName"`), `OrganizationUnitName` (`json:"organizationUnitName"`), `Description` (`json:"description,omitempty"`)
    - `TenantStatus`: `Phase` (`json:"phase"`), `Message` (`json:"message,omitempty"`)
    - Constants: `TenantAPIVersion = "platform.sovrunn.io/v1alpha1"`, `TenantKind = "Tenant"`
    - Reuse existing phase constants
    **Tests required:** None (tested via registry/handler tests)
    **Completion criteria:**
    - `make fmt` passes; package compiles
    **Verification:** `make fmt && make vet`

---

- [ ] 2. Add Tenant validation
  - [ ] 2.1 Implement ValidateTenant and ValidateTenantPathSegments with tests
    **Objective:** Pure validation functions for Tenant fields and path segments.
    **Files to create:**
    - `internal/validation/tenant.go`
    - `internal/validation/tenant_test.go`
    - `internal/validation/tenant_property_test.go`
    **Implementation notes:**
    - `ValidateTenant(t resources.Tenant) []resources.FieldError` — context-free
    - `ValidateTenantPathSegments(orgName, ouName, name string) []resources.FieldError` — context-free
    - Rules: metadata.name, spec.organizationName, spec.organizationUnitName each DNS-label 1–63
    - Reuse `dnsLabelRe` and `validateName` from validation package
    - Path field mapping: orgName → `spec.organizationName`, ouName → `spec.organizationUnitName`, name → `metadata.name`
    - No I/O, no registry lookup
    **Tests required:**
    - Unit: valid names accepted; empty/invalid/long metadata.name rejected; empty/invalid/long spec.organizationName rejected; empty/invalid/long spec.organizationUnitName rejected; path validation maps fields correctly
    - Property (testing/quick, Config{MaxCount: 100}): P1 valid DNS-labels accepted; P2 invalid strings rejected
    **Completion criteria:**
    - All validation tests pass
    - `make fmt && make vet && make test` pass
    **Verification:** `make fmt && make vet && make test`

---

- [ ] 3. Add Tenant registry
  - [ ] 3.1 Implement TenantRegistry with unit, property, and race tests
    **Objective:** Thread-safe in-memory registry with composite key and deep copies.
    **Files to create:**
    - `internal/registry/tenant_registry.go`
    - `internal/registry/tenant_registry_test.go`
    - `internal/registry/tenant_registry_property_test.go`
    - `internal/registry/tenant_registry_race_test.go`
    **Implementation notes:**
    - `TenantRegistryIface` with 6 methods (Create, Get, List, Update, Delete, CountByOrganizationUnit)
    - `CreateTenant(ctx, t) (resources.Tenant, error)` — returns stored deep copy
    - `UpdateTenant(ctx, t) (resources.Tenant, error)` — derives composite key from submitted resource; preserves stored APIVersion, Kind, Status, Metadata.Name, Spec.OrganizationName, Spec.OrganizationUnitName; replaces only Metadata.DisplayName, Metadata.Labels, Metadata.Annotations, Spec.Description
    - `tenantCompositeKey(orgName, ouName, name)` → `orgName + "/" + ouName + "/" + name`
    - `sync.RWMutex`: RLock for Get/List/Count, Lock for Create/Update/Delete
    - `deepCopyTenant` duplicates Labels and Annotations maps
    - List sorted: organizationName asc → organizationUnitName asc → name asc
    - Define `OrganizationUnitLookup` interface here
    - Reuse `ErrNotFound`, `ErrAlreadyExists`; no global state; no dependency on other registries
    **Tests required:**
    - Unit: Create stores; duplicate → ErrAlreadyExists (original unchanged); same name under different OUs succeeds; Get by composite key; Get non-existent → ErrNotFound; List sorted; empty List → non-nil empty slice; Update mutable fields only; Update non-existent → ErrNotFound; Delete removes; Delete non-existent → ErrNotFound; CountByOrganizationUnit correct
    - Property (testing/quick): P3 round-trip; P4 sort invariant; P5 deep-copy immutability; P6 duplicate idempotent error
    - Race: 10+ goroutines mixed Create/Get/List/Update/Delete/CountByOrganizationUnit
    **Completion criteria:**
    - All tests pass; `make fmt && make vet && make test && go test -race ./internal/registry/...` pass
    **Verification:** `make fmt && make vet && make test && go test -race ./internal/registry/...`

---

- [ ] 4. Add Tenant child blocker for OrganizationUnit delete
  - [ ] 4.1 Implement TenantChildBlockerChecker and OUChildBlocker interface
    **Objective:** Blocker that prevents OU deletion when Tenants reference it.
    **Files to create:**
    - `internal/registry/tenant_blocker.go`
    - `internal/registry/tenant_blocker_test.go`
    **Implementation notes:**
    - Define `OUChildBlocker` interface: `BlockedByOUChildren(ctx, orgName, ouName string) ([]BlockedBy, error)`
    - `TenantChildBlockerChecker` holds `TenantRegistryIface`
    - `NewTenantChildBlockerChecker(reg) *TenantChildBlockerChecker`
    - `BlockedByOUChildren` calls `CountByOrganizationUnit`; count > 0 → `[]BlockedBy{{Kind: "Tenant", Count: count}}`; count 0 → nil; propagate registry error
    **Tests required:**
    - Zero tenants → no blockers
    - One+ tenants → blocking kind "Tenant" with correct count
    - Registry error propagates
    **Completion criteria:**
    - `make fmt && make vet && make test` pass
    **Verification:** `make fmt && make vet && make test`

---

- [ ] 5. Extend OUHandler with optional OUChildBlocker
  - [ ] 5.1 Add blocker dependency to OUHandler
    **Objective:** Wire the Tenant blocker into the OU delete path.
    **Files to modify:**
    - `internal/api/ou_handler.go`
    - `internal/api/ou_handler_test.go` (and any existing NewOUHandler call sites)
    **Implementation notes:**
    - Add `blocker registry.OUChildBlocker` field to OUHandler
    - Update `NewOUHandler(reg, orgLookup, blocker)` signature
    - Allow `blocker` nil → OU delete proceeds as before (FEATURE-0002 behavior)
    - If `blocker` non-nil, DELETE consults `BlockedByOUChildren` before deleting
    - If blockers returned → HTTP 409 DELETE_BLOCKED, message identifying "Tenant"
    - Do NOT create generic blocker framework; do NOT rewrite OrgHandler
    - Update existing tests impacted by the new NewOUHandler signature
    **Tests required:**
    - nil blocker keeps existing delete behavior (204)
    - blocker returning Tenant blocks OU delete → 409 DELETE_BLOCKED
    - blocker returning empty allows delete → 204
    **Completion criteria:**
    - `make fmt && make vet && make test && go test -race ./internal/api/...` pass
    **Verification:** `make fmt && make vet && make test && go test -race ./internal/api/...`

---

- [ ] 6. Add Tenant safe JSON decoder
  - [ ] 6.1 Implement safeDecodeTenant
    **Objective:** Safe JSON decode function for Tenant request bodies.
    **Files to create:**
    - `internal/api/tenant_decode.go`
    **Implementation notes:**
    - `http.MaxBytesReader(w, r.Body, 1<<20)` inside the function
    - Read body; `*http.MaxBytesError` → errBodyTooLarge; empty → errEmptyBody
    - Decode into `map[string]json.RawMessage`; "status" key present → errStatusFieldPresent
    - Typed decode with `DisallowUnknownFields()`; unknown → errUnknownField; syntax/type → errMalformedJSON
    - Do not echo raw body
    - 415 remains contentTypeMiddleware responsibility, not the decoder
    - Reuse existing error sentinels from `decode.go` if available
    **Tests required:** Tested via handler tests in Task 7
    **Completion criteria:**
    - `make fmt && make vet` pass; function compiles
    **Verification:** `make fmt && make vet`

---

- [ ] 7. Add Tenant HTTP handler
  - [ ] 7.1 Implement TenantHandler with all CRUD methods and tests
    **Objective:** HTTP handler for Tenant CRUD endpoints.
    **Files to create:**
    - `internal/api/tenant_handler.go`
    - `internal/api/tenant_handler_test.go`
    **Implementation notes:**
    - `TenantHandler` takes `TenantRegistryIface` + `OrganizationUnitLookup`
    - `NewTenantHandler(reg, ouLookup)`
    - `HandleCollection`: POST → Create, GET → List, else 405
    - `HandleItem`: `strings.TrimPrefix` + `strings.Split`; exactly 3 non-empty segments else 404; invalid DNS segment → 400; dispatch GET/PUT/DELETE, else 405
    - Path field mapping: orgName → spec.organizationName, ouName → spec.organizationUnitName, name → metadata.name
    - Create: decode → validate → parent OU lookup (missing → 400 field=spec.organizationUnitName, message includes full parent ref) → set apiVersion/kind/status.phase=Active → registry.CreateTenant → duplicate → 409 → `// TODO(FEATURE-0005): emit Operation record — type: CreateTenant` → 201
    - Get: validate path → registry.GetTenant → 404 → 200
    - List: registry.ListTenants → `{"items": items}` ([] when empty, only items field)
    - Update: validate path → decode → require body metadata.name/spec.organizationName/spec.organizationUnitName present and matching path → validate → registry.UpdateTenant(ctx, tenant) → 404 → `// TODO(FEATURE-0005): emit Operation record — type: UpdateTenant` → 200
    - Delete: validate path → registry.DeleteTenant → 404 → `// TODO(FEATURE-0005): emit Operation record — type: DeleteTenant` → 204 no body
    - Reuse shared writeError/writeJSON/writeValidationErrors
    **Tests required:**
    - POST 201 valid; POST 409 duplicate; POST 400 (invalid fields, non-existent parent, status key, bad JSON); POST 413 oversized body
    - GET 200/404/400; GET wrong path shape → 404
    - List sorted; list empty → []
    - PUT 200/404; PUT 400 (name mismatch, orgName mismatch, ouName mismatch, status key)
    - DELETE 204/404/400
    **Completion criteria:**
    - `make fmt && make vet && make test && go test -race ./internal/api/...` pass
    **Verification:** `make fmt && make vet && make test && go test -race ./internal/api/...`

---

- [ ] 8. Wire Tenant routes into server
  - [ ] 8.1 Update server.New and route registration
    **Objective:** Register Tenant routes and require TenantHandler.
    **Files to modify:**
    - `internal/server/server.go`
    - `internal/server/server_test.go`
    **Implementation notes:**
    - Update `server.New` to accept non-nil `*api.TenantHandler`
    - Register `/v1/tenants` → HandleCollection, `/v1/tenants/` → HandleItem
    - Keep middleware order: requestID → logging → contentType → handler
    - Update tests constructing `server.New` to provide a TenantHandler fixture
    - Add route tests if the existing pattern supports them
    **Completion criteria:**
    - `make fmt && make vet && make test && make build` pass
    **Verification:** `make fmt && make vet && make test && make build`

---

- [ ] 9. Wire Tenant registry, handler, and blocker in main.go
  - [ ] 9.1 Update main.go wiring
    **Objective:** Construct and inject all Tenant dependencies.
    **Files to modify:**
    - `cmd/sovrunn-api/main.go`
    **Implementation notes:**
    ```
    tenantRegistry := registry.NewTenantRegistry()
    tenantBlocker := registry.NewTenantChildBlockerChecker(tenantRegistry)
    ouHandler := api.NewOUHandler(ouRegistry, orgRegistry, tenantBlocker)
    tenantHandler := api.NewTenantHandler(tenantRegistry, ouRegistry)
    srv := server.New(cfg, orgHandler, ouHandler, tenantHandler, bootstrapHandler, readiness)
    ```
    - Update any existing NewOUHandler call sites for the new signature
    **Completion criteria:**
    - `make fmt && make vet && make test && make build` pass
    - `make run` starts; existing endpoints and new `/v1/tenants` routes reachable
    **Verification:** `make fmt && make vet && make test && make build`

---

- [ ] 10. Add integration tests for OU delete blocked by Tenant
  - [ ] 10.1 Write OU-delete-blocked integration tests
    **Objective:** Verify end-to-end OU deletion blocking via Tenant.
    **Files to create/modify:**
    - `internal/api/tenant_handler_test.go` or `internal/server/server_test.go`
    **Tests required:**
    - Create org → create OU → create Tenant under OU → DELETE OU → 409 DELETE_BLOCKED with "Tenant" in message
    - Create org → create OU → no Tenant → DELETE OU → 204
    **Completion criteria:**
    - `make fmt && make vet && make test` pass
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
- Each property test tagged: `// Feature: tenant-resource, Property N: <title>`
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
