# Implementation Plan: FEATURE-0002 OrganizationUnit Resource and Registry

## Overview

Implement OrganizationUnit as a child governance resource under Organization.
Composite identity: `spec.organizationName/metadata.name`. Storage-only registry.
Parent existence checks in API layer. Go 1.21 stdlib only.

## Cursor Execution Rule

Implement **one numbered task at a time**. After each task:
1. Run the verification commands listed in the task.
2. Confirm all pass with zero errors before moving to the next task.
3. Do not begin a subsequent task until the current task's verification succeeds.
4. If a task fails verification, fix within the current task before proceeding.

## Tasks

- [ ] 1. OrganizationUnit resource model
  - [ ] 1.1 Create OrganizationUnit structs and constants
    **Objective:** Define the OrganizationUnit, OrganizationUnitSpec, and OrganizationUnitStatus structs with JSON tags and constants.
    **Files to create:**
    - `internal/resources/organizationunit.go` (or `organization_unit.go` — match existing repo naming style)
    **Implementation notes:**
    - Reuse existing `Metadata` struct from `organization.go`
    - `OrganizationUnitSpec`: `OrganizationName` (`json:"organizationName"`), `Description` (`json:"description,omitempty"`)
    - `OrganizationUnitStatus`: `Phase` (`json:"phase"`), `Message` (`json:"message,omitempty"`)
    - Constants: `OUAPIVersion = "platform.sovrunn.io/v1alpha1"`, `OUKind = "OrganizationUnit"`
    - Reuse existing phase constants (`PhaseActive`, etc.)
    **Tests required:** None (struct definitions only; tested via registry/handler tests later)
    **Acceptance criteria:**
    - `make fmt` and `make vet` pass
    - Structs compile and are importable from other packages
    **Verification:** `make fmt && make vet`

---

- [ ] 2. Validation functions and tests
  - [ ] 2.1 Implement ValidateOrganizationUnit and ValidateOUPathSegments
    **Objective:** Create pure validation functions for OrganizationUnit fields and path segments.
    **Files to create:**
    - `internal/validation/organizationunit.go` (or `organization_unit.go` — match repo style)
    - `internal/validation/organizationunit_test.go`
    - `internal/validation/organizationunit_property_test.go`
    **Implementation notes:**
    - `ValidateOrganizationUnit(ou resources.OrganizationUnit) []resources.FieldError` — context-free
    - `ValidateOUPathSegments(orgName, name string) []resources.FieldError` — context-free
    - Reuse existing `dnsLabelRe` and `validateName` from `organization.go` (same package)
    - Path validation field mapping:
      - invalid name segment → `error.field = "metadata.name"`
      - invalid orgName segment → `error.field = "spec.organizationName"`
    - `spec.organizationName`: required, DNS-label format, max 63 chars
    **Tests required:**
    - Unit tests: valid names accepted, empty name rejected, uppercase rejected, spaces rejected, >63 chars rejected, leading/trailing hyphen rejected, empty spec.organizationName rejected, invalid spec.organizationName format rejected, valid spec.organizationName passes
    - Property tests (testing/quick, Config{MaxCount: 100}):
      - Property 1: valid DNS-label names accepted
      - Property 2: invalid/arbitrary strings rejected
    **Acceptance criteria:**
    - All validation unit tests pass
    - Property tests pass with 100 iterations
    - `make fmt && make vet && make test` pass
    **Verification:** `make fmt && make vet && make test`

---

- [ ] 3. OrganizationUnit registry interface and implementation
  - [ ] 3.1 Implement OrganizationUnitRegistryIface and OrganizationUnitRegistry
    **Objective:** Create the thread-safe in-memory registry with composite key, deep copies, and CountByOrganization.
    **Files to create:**
    - `internal/registry/ou_registry.go`
    **Implementation notes:**
    - Interface: `OrganizationUnitRegistryIface` with 6 methods (Create, Get, List, Update, Delete, CountByOrganization)
    - `CreateOrganizationUnit(ctx, ou) (resources.OrganizationUnit, error)` — returns stored deep copy
    - `UpdateOrganizationUnit(ctx, orgName, name, ou) (resources.OrganizationUnit, error)` — returns updated deep copy
    - Composite key: `compositeKey(orgName, name string) string` → `orgName + "/" + name`
    - `sync.RWMutex` protection: RLock for Get/List/Count, Lock for Create/Update/Delete
    - `deepCopyOrganizationUnit` duplicates Labels and Annotations maps
    - List sorted: `spec.organizationName` ascending, then `metadata.name` ascending
    - Update preserves: Metadata.Name, Spec.OrganizationName, Status, APIVersion, Kind
    - Sentinel errors: reuse `ErrNotFound`, `ErrAlreadyExists` from `registry.go`
    - No package-level global state; constructor `NewOrganizationUnitRegistry()`
    - Also define `OrganizationLookup` interface in this file:
      ```
      type OrganizationLookup interface {
          GetOrganization(ctx context.Context, name string) (resources.Organization, error)
      }
      ```
    **Tests required:** None in this task (tested in Task 4)
    **Acceptance criteria:**
    - `make fmt && make vet` pass
    - Registry struct compiles and exports correct interface
    **Verification:** `make fmt && make vet`

---

- [ ] 4. Registry unit, property, and race tests
  - [ ] 4.1 Write OrganizationUnit registry tests
    **Objective:** Comprehensive test coverage for the OrganizationUnit registry.
    **Files to create:**
    - `internal/registry/ou_registry_test.go`
    - `internal/registry/ou_registry_property_test.go`
    - `internal/registry/ou_registry_race_test.go`
    **Tests required:**
    - Unit tests: Create stores OU, duplicate composite key → ErrAlreadyExists, same name under different orgs succeeds, Get by composite key, Get non-existent → ErrNotFound, List sorted by orgName then name, List empty → [], Update modifies mutable fields, Update non-existent → ErrNotFound, Delete removes OU, Delete non-existent → ErrNotFound, CountByOrganization correct count
    - Property tests (testing/quick, Config{MaxCount: 100}):
      - Property 3: Create/Get round trip preserves data
      - Property 4: Deep copy immutability (mutating returned maps doesn't affect registry)
      - Property 5: List sorted by orgName then name
      - Property 6: Update preserves immutable system fields
      - Property 7: Duplicate composite key returns ErrAlreadyExists (idempotent error + entry preservation)
    - Race test: 10+ goroutines performing concurrent Create/Get/List/Update/Delete/CountByOrganization
    **Acceptance criteria:**
    - All unit tests pass
    - All property tests pass (100 iterations each)
    - Race test passes with `go test -race`
    - `make fmt && make vet && make test && go test -race ./internal/registry/...` pass
    **Verification:** `make fmt && make vet && make test && go test -race ./internal/registry/...`

---

- [ ] 5. OrganizationUnit blocker for Organization delete
  - [ ] 5.1 Implement OUChildBlockerChecker
    **Objective:** Implement ChildBlockerChecker that blocks Organization deletion when OUs exist.
    **Files to create:**
    - `internal/registry/ou_blocker.go`
    **Implementation notes:**
    - `OUChildBlockerChecker` struct holds `OrganizationUnitRegistryIface`
    - `NewOUChildBlockerChecker(ouReg OrganizationUnitRegistryIface) *OUChildBlockerChecker`
    - `BlockedByChildren(ctx, orgName) ([]BlockedBy, error)` calls `CountByOrganization`
    - Returns `[]BlockedBy{{Kind: "OrganizationUnit", Count: n}}` when n > 0
    - Returns `nil, nil` when count is 0
    **Tests required:** None in this task (blocker tested in integration Task 10)
    **Acceptance criteria:**
    - `make fmt && make vet` pass
    - OUChildBlockerChecker satisfies ChildBlockerChecker interface (compile check)
    **Verification:** `make fmt && make vet`

---

- [ ] 6. Safe JSON decoder for OrganizationUnit
  - [ ] 6.1 Implement safeDecodeOrganizationUnit
    **Objective:** Create the safe JSON decode function for OrganizationUnit request bodies.
    **Files to create:**
    - `internal/api/ou_decode.go`
    **Implementation notes:**
    - Apply `http.MaxBytesReader(w, r.Body, 1<<20)` inside this function
    - Read body bytes, check for `*http.MaxBytesError` → errBodyTooLarge
    - Empty body → errEmptyBody
    - Decode into `map[string]json.RawMessage`; if "status" key present → errStatusFieldPresent
    - Decode into typed `OrganizationUnit` struct with `DisallowUnknownFields()`
    - Map decode errors: syntax → errMalformedJSON, unknown field → errUnknownField
    - Decode errors map to HTTP 400 or HTTP 413 only. HTTP 415 (Unsupported Media Type) is handled by `contentTypeMiddleware`, not by `safeDecodeOrganizationUnit`.
    - Reuse existing error sentinel values from `decode.go` if available
    **Tests required:** Tested via handler tests in Task 8
    **Acceptance criteria:**
    - `make fmt && make vet` pass
    - Function compiles with correct signature
    **Verification:** `make fmt && make vet`

---

- [ ] 7. OU HTTP handler (Create/Get/List/Update/Delete)
  - [ ] 7.1 Implement OUHandler with all CRUD methods
    **Objective:** Create the HTTP handler for OrganizationUnit CRUD endpoints.
    **Files to create:**
    - `internal/api/ou_handler.go`
    **Files to modify:**
    - `internal/api/response.go` — move `writeValidationErrors` here if currently in `org_handler.go`
    **Implementation notes:**
    - `OUHandler` struct: takes `OrganizationUnitRegistryIface` + `OrganizationLookup`
    - `NewOUHandler(reg, orgLookup) *OUHandler`
    - `HandleCollection(w, r)`: POST → Create, GET → List
    - `HandleItem(w, r)`: extract path via `strings.TrimPrefix` + `strings.Split`
      - `parts := strings.Split(remainder, "/")` — reject if `len(parts) != 2` or empty parts → 404
      - Then validate segments for DNS-label → 400 if invalid
      - Dispatch: GET → Get, PUT → Update, DELETE → Delete
    - Create handler flow: decode → validate → orgLookup → force server fields → registry.Create → writeJSON(201, created)
    - Update handler: validate path → decode → require body name + orgName match path → validate → registry.Update → writeJSON(200, updated)
    - Delete handler: validate path → registry.Delete → 204
    - PUT requires body.Metadata.Name AND body.Spec.OrganizationName present and matching path
    - Operation placeholders: `// TODO(FEATURE-0005): emit Operation record — type: CreateOrganizationUnit` (etc.)
    - Reuse `writeError`, `writeJSON`, `writeValidationErrors` from `response.go`
    - Do NOT duplicate validation error response logic
    **Tests required:** Tested in Task 8
    **Acceptance criteria:**
    - `make fmt && make vet` pass
    - Handler compiles; all methods implemented per design pseudocode
    **Verification:** `make fmt && make vet`

---

- [ ] 8. Handler tests
  - [ ] 8.1 Write OU HTTP handler tests
    **Objective:** Comprehensive HTTP handler tests for all OrganizationUnit endpoints.
    **Files to create:**
    - `internal/api/ou_handler_test.go`
    **Tests required:**
    Using `net/http/httptest`:
    - POST 201 valid with existing parent Organization
    - POST 409 duplicate composite key
    - POST 400 invalid metadata.name
    - POST 400 missing spec.organizationName
    - POST 400 non-existent parent Organization
    - POST 400 status key present (test with {}, null, {"phase":""})
    - POST 400 bad JSON
    - POST 413 oversized body (>1 MiB)
    - POST 415 wrong Content-Type
      > **Note:** 415 is produced by `contentTypeMiddleware`, not by `OUHandler` itself.
      > Either wrap the handler with the full middleware chain (`requestID → logging → contentType`)
      > in this test, or defer the 415 test to Task 9/10 where server route wiring is tested.
      > Direct `OUHandler` tests alone are not sufficient for 415.
    - GET 200 existing resource with full shape
    - GET 404 missing resource
    - GET 400 invalid path segments
    - GET 404 bare /v1/organization-units/ with no segments
    - GET 404 single segment only (e.g., /v1/organization-units/nic)
    - GET list 200 sorted by orgName then name
    - GET list 200 empty → {"items": []}
    - PUT 200 valid update
    - PUT 404 missing resource
    - PUT 400 metadata.name absent in body
    - PUT 400 metadata.name mismatch with path
    - PUT 400 spec.organizationName absent in body
    - PUT 400 spec.organizationName mismatch with path
    - PUT 400 status key present
    - DELETE 204 existing resource
    - DELETE 404 missing resource
    - DELETE 400 invalid path segments
    **Acceptance criteria:**
    - All handler tests pass
    - `make fmt && make vet && make test && go test -race ./internal/api/...` pass
    **Verification:** `make fmt && make vet && make test && go test -race ./internal/api/...`

---

- [ ] 9. Server and main.go wiring
  - [ ] 9.1 Wire OUHandler, registry, and blocker into server
    **Objective:** Register OU routes, inject dependencies, replace NoopChildBlockerChecker.
    **Files to modify:**
    - `internal/server/server.go` — update `New()` signature to accept `*api.OUHandler`; register OU routes
    - `cmd/sovrunn-api/main.go` — wire `OrganizationUnitRegistry`, `OUChildBlockerChecker`, `OUHandler`
    **Implementation notes:**
    - Updated `server.New` signature: `New(cfg, org *api.OrgHandler, ou *api.OUHandler, bootstrap, readiness) *Server`
    - Register routes (Go 1.21 compatible):
      ```
      mux.Handle("/v1/organization-units", chain(http.HandlerFunc(ou.HandleCollection)))
      mux.Handle("/v1/organization-units/", chain(http.HandlerFunc(ou.HandleItem)))
      ```
    - Middleware chain remains: requestID → logging → contentType → handler
    - main.go wiring:
      ```
      ouRegistry := registry.NewOrganizationUnitRegistry()
      ouBlocker := registry.NewOUChildBlockerChecker(ouRegistry)
      orgHandler := api.NewOrgHandler(orgRegistry, ouBlocker)  // replaces NoopChildBlockerChecker
      ouHandler := api.NewOUHandler(ouRegistry, orgRegistry)   // orgRegistry satisfies OrganizationLookup
      srv := server.New(cfg, orgHandler, ouHandler, bootstrapHandler, readiness)
      ```
    **Tests required:** Tested via existing server tests and integration test in Task 10
    **Acceptance criteria:**
    - Server compiles and starts with both Organization and OrganizationUnit routes
    - `make fmt && make vet && make test && make build` pass
    - `make run` starts; `/healthz`, `/readyz`, `/version` still work
    **Verification:** `make fmt && make vet && make test && make build`

---

- [ ] 10. Organization delete blocker integration tests
  - [ ] 10.1 Write blocker integration tests
    **Objective:** Verify Organization delete is blocked when OrganizationUnits reference it.
    **Files to create:**
    - `internal/registry/ou_blocker_test.go` — registry-level blocker unit tests (CountByOrganization, BlockedByChildren logic)
    - `internal/api/ou_handler_test.go` (or a server-level test file) — HTTP-level delete-blocked tests that exercise `OrgHandler` + `OUChildBlockerChecker` wiring together
    > **Note:** Registry-level tests verify the blocker returns correct `BlockedBy` values.
    > HTTP-level tests verify the full request path: `DELETE /v1/organizations/{name}` → OrgHandler → blocker → 409/204.
    **Tests required:**
    - Test: Create Organization "nic", create OU "ministry-health" under "nic", attempt DELETE /v1/organizations/nic → 409 DELETE_BLOCKED with "OrganizationUnit" in message
    - Test: Create Organization "empty-org" with no OUs, DELETE /v1/organizations/empty-org → 204
    - Test: CountByOrganization returns 0 for org with no OUs
    - Test: CountByOrganization returns correct count for org with multiple OUs
    **Acceptance criteria:**
    - Blocker correctly blocks Organization delete when OUs exist
    - Blocker allows Organization delete when no OUs reference it
    - `make fmt && make vet && make test` pass
    **Verification:** `make fmt && make vet && make test`

---

- [ ] 11. Final verification and cleanup
  - [ ] 11.1 Run full verification suite
    **Objective:** Confirm all tests pass, no race conditions, binary builds.
    **Verification commands:**
    ```
    make fmt
    make vet
    make test
    go test -race ./...
    make build
    ```
    **Acceptance criteria:**
    - Zero formatting diffs
    - Zero vet errors
    - Zero test failures
    - Zero race conditions
    - Binary builds at bin/sovrunn-api
    - Runtime smoke test: `make run` then verify /healthz, /readyz, /version, POST/GET/PUT/DELETE /v1/organization-units
    **Report when done:**
    - Files created/modified
    - New exported symbols
    - Validation rules implemented
    - Tests written (unit + property + race + blocker)
    - Non-goals not implemented

## Notes

- Each task references the approved design.md for implementation details
- Property tests use `testing/quick` with `Config{MaxCount: 100}`
- Each property test tagged: `// Feature: organizationunit-resource, Property N: <title>`
- Module path: `github.com/sanjeevksaini/sovrunn`
- Go 1.21 minimum — no Go 1.22 wildcard routing
- No new external dependencies required
- File naming: match existing repo style (underscores vs no-underscores)

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1"] },
    { "id": 1, "tasks": ["2.1"] },
    { "id": 2, "tasks": ["3.1"] },
    { "id": 3, "tasks": ["4.1", "5.1"] },
    { "id": 4, "tasks": ["6.1"] },
    { "id": 5, "tasks": ["7.1"] },
    { "id": 6, "tasks": ["8.1"] },
    { "id": 7, "tasks": ["9.1"] },
    { "id": 8, "tasks": ["10.1"] },
    { "id": 9, "tasks": ["11.1"] }
  ]
}
```
