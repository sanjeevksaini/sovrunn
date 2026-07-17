# Implementation Plan: FEATURE-0006 ServiceClass and ServicePlan

## Overview

Implement `ServiceClass` and `ServicePlan` as **global platform catalog resources**. Catalog
definitions only — no provisioning, binding, or execution. In-memory storage-only registries;
composite `serviceClassName/name` identity for ServicePlan; ServiceClass delete blocked while child
ServicePlans exist; catalog lifecycle Operations emitted via the FEATURE-0005 nil-safe emitter.

## Cursor Execution Rule

Implement **one numbered task at a time**. After each task:
1. Run the verification commands listed in the task.
2. Confirm all pass with zero errors before moving to the next task.
3. **Do not start the next task** until the current task's verification succeeds.
4. If a task fails verification, fix within the current task before proceeding.

## Global Constraints

- Go 1.21-compatible only; no Go 1.22 wildcard routing; no external dependencies.
- `internal/api` MUST NOT import `internal/server`.
- ServiceClass and ServicePlan are global catalog resources — NOT scoped to Organization,
  OrganizationUnit, Tenant, or Project.
- Registries are storage-only; no cross-registry dependencies. Parent existence checks happen in the
  handler via the narrow `ServiceClassLookup` interface.
- Reuse the existing FEATURE-0005 `OperationEmitter` and nil-safe `emitOperation` helper.
- Operation emission failure MUST NOT change the primary API response.
- No secrets stored in catalog resources; forbidden `spec.parameters` keys per Requirement 4.6.
- No new logger fields on handlers solely for emission.
- Existing handler tests must pass with a nil emitter and nil blocker.

## Tasks

- [ ] 1. Add ServiceClass and ServicePlan resource models/constants
  - [ ] 1.1 Create ServiceClass and ServicePlan structs and constants
    **Files to create:** `internal/resources/serviceclass.go`, `internal/resources/serviceplan.go`
    **Implementation notes:**
    - `ServiceClass` (APIVersion, Kind, Metadata, Spec, Status) with standard JSON tags.
    - `ServiceClassSpec`: DisplayName, Description (omitempty), Category, Provider (omitempty),
      Lifecycle, DefaultPlanName (omitempty), Tags (omitempty, `[]string`).
    - `ServiceClassStatus`: Phase, Message (omitempty).
    - `ServicePlan` (APIVersion, Kind, Metadata, Spec, Status).
    - `ServicePlanSpec`: ServiceClassName, DisplayName (omitempty), Description (omitempty), Tier,
      Lifecycle, Parameters (omitempty, `map[string]string`), Tags (omitempty, `[]string`).
    - `ServicePlanStatus`: Phase, Message (omitempty).
    - Constants: `ServiceClassKind = "ServiceClass"`, `ServicePlanKind = "ServicePlan"`.
    - Category constants (8): Database, Cache, ObjectStorage, Stream, Gateway, Function, Analytics,
      Other. Lifecycle constants (4): Preview, Active, Deprecated, Retired. Tier constants (6):
      Dev, Small, Medium, Large, Production, Custom.
    - Reuse existing `Metadata` struct and `PhaseActive` phase constant.
    **Tests to run:** none (covered by later layers).
    **Acceptance criteria:** structs and constants compile; JSON tags match Requirements 1–2.
    **Commit message:** `feat(resources): add ServiceClass and ServicePlan models (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 1.1, 1.3, 1.4, 2.1, 2.3, 2.4, 12.3_

---

- [ ] 2. Add Operation constants and OperationSpec catalog reference fields
  - [ ] 2.1 Extend operation.go with catalog types and reference fields
    **Files to modify:** `internal/resources/operation.go`
    **Implementation notes:**
    - Add 6 operation type constants: `OpCreateServiceClass`, `OpUpdateServiceClass`,
      `OpDeleteServiceClass`, `OpCreateServicePlan`, `OpUpdateServicePlan`, `OpDeleteServicePlan`.
    - Add two optional `OperationSpec` fields: `ServiceClassName string json:"serviceClassName,omitempty"`
      and `ServicePlanName string json:"servicePlanName,omitempty"`.
    - Do NOT change existing fields, tags, or constants.
    **Tests to run:** `make fmt && make vet && make test` (existing Operation tests still pass).
    **Acceptance criteria:** new constants/fields present; existing FEATURE-0005 tests unaffected.
    **Commit message:** `feat(resources): add catalog Operation types and spec fields (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 12.4, 12.5, 12.6_

---

- [ ] 3. Add ServiceClass validation
  - [ ] 3.1 Implement ValidateServiceClass and path-segment validation
    **Files to create:** `internal/validation/serviceclass.go`
    **Implementation notes:**
    - `ValidateServiceClass(sc resources.ServiceClass) []resources.FieldError` — pure, context-free.
    - `metadata.name` DNS-label (1–63) → `metadata.name`.
    - `spec.category` required, ∈ 8-value set → `spec.category`.
    - `spec.lifecycle` required, ∈ 4-value set → `spec.lifecycle`.
    - If `spec.defaultPlanName` non-empty, must be DNS-label → `spec.defaultPlanName` (existence NOT
      verified).
    - `ValidateServiceClassPathSegment(name string) []resources.FieldError`.
    - Reuse `dnsLabelRe`/`validateName` from the validation package; enum checks via switch/sets.
    **Tests to run:** `make fmt && make vet` (unit tests added in task 5).
    **Acceptance criteria:** valid inputs pass; invalid category/lifecycle/name/defaultPlanName
    produce the correct field errors.
    **Commit message:** `feat(validation): add ServiceClass validation (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 7.3_

---

- [ ] 4. Add ServicePlan validation
  - [ ] 4.1 Implement ValidateServicePlan, path segments, and forbidden-key logic
    **Files to create:** `internal/validation/serviceplan.go`
    **Implementation notes:**
    - `ValidateServicePlan(sp resources.ServicePlan) []resources.FieldError` — pure, context-free.
    - `metadata.name` DNS-label → `metadata.name`; `spec.serviceClassName` required DNS-label →
      `spec.serviceClassName`.
    - `spec.tier` required, ∈ 6-value set → `spec.tier`; `spec.lifecycle` required, ∈ 4-value set →
      `spec.lifecycle`.
    - `forbiddenParamSubstrings` = {password, secret, token, credential, auth, apikey, accesskey,
      secretkey, privatekey}; `isForbiddenParamKey` lowercases the key and checks `strings.Contains`.
      Plain `key` NOT rejected; `regionKey`/`masterKeyCount` allowed; `apiKey` rejected. Offending
      key → single `FieldError` field `spec.parameters`.
    - `ValidateServicePlanPathSegments(serviceClassName, name string)`: invalid serviceClassName →
      `spec.serviceClassName`; invalid name → `metadata.name`.
    **Tests to run:** `make fmt && make vet` (unit tests added in task 5).
    **Acceptance criteria:** forbidden keys rejected case-insensitively; benign keys accepted.
    **Commit message:** `feat(validation): add ServicePlan validation (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 8.3, 13.2_

---

- [ ] 5. Add validation unit/property tests
  - [ ] 5.1 ServiceClass and ServicePlan validation tests
    **Files to create:** `internal/validation/serviceclass_test.go`,
    `internal/validation/serviceplan_test.go`,
    `internal/validation/serviceclass_property_test.go`,
    `internal/validation/serviceplan_property_test.go`
    **Implementation notes:**
    - ServiceClass: valid accepted; empty/invalid/long `metadata.name` rejected; invalid/empty
      category and lifecycle rejected; non-DNS-label `defaultPlanName` rejected; empty
      `defaultPlanName` accepted; path-segment validation.
    - ServicePlan: valid accepted; empty/invalid name and serviceClassName rejected; invalid/empty
      tier and lifecycle rejected; forbidden keys rejected (`apiKey`, `ACCESSKEY`, `Secret`); benign
      keys accepted (`regionKey`, `key`, `masterKeyCount`); path-segment field mapping.
    - Property tests use `testing/quick` `Config{MaxCount: 100}`, tagged
      `// Feature: serviceclass-serviceplan, Property N: <title>` (Properties 1–3).
    **Tests to run:** `make fmt && make vet && make test`
    **Acceptance criteria:** all validation unit and property tests pass deterministically.
    **Commit message:** `test(validation): ServiceClass/ServicePlan validation tests (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 14.1, 15.1, 15.3_

---

- [ ] 6. Add ServiceClass registry
  - [ ] 6.1 Implement ServiceClassRegistry and ServiceClassLookup
    **Files to create:** `internal/registry/serviceclass_registry.go`
    **Implementation notes:**
    - `ServiceClassRegistryIface`: Create/Get/List/Update/Delete (all ctx-first).
    - `ServiceClassRegistry`: `sync.RWMutex` + `map[string]resources.ServiceClass` keyed by
      `metadata.name`; `NewServiceClassRegistry()`; no global state.
    - `deepCopyServiceClass` copies Tags slice and Metadata Labels/Annotations maps.
    - Create → `ErrAlreadyExists` on dup; Get → `ErrNotFound`; List → non-nil, sorted by name;
      Update preserves APIVersion/Kind/Status/Metadata.Name, replaces mutable fields, `ErrNotFound`
      if absent; Delete → `ErrNotFound` if absent.
    - `ServiceClassLookup` interface (`GetServiceClass`) in the registry package; `*ServiceClassRegistry`
      satisfies it. Storage-only; no dependency on other registries.
    **Tests to run:** `make fmt && make vet` (tests added in task 8).
    **Acceptance criteria:** methods behave per design; returns deep copies.
    **Commit message:** `feat(registry): add ServiceClass registry (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7, 5.8_

---

- [ ] 7. Add ServicePlan registry
  - [ ] 7.1 Implement ServicePlanRegistry with composite identity and CountByServiceClass
    **Files to create:** `internal/registry/serviceplan_registry.go`
    **Implementation notes:**
    - `ServicePlanRegistryIface`: Create/Get/List/Update/Delete + `CountByServiceClass` (ctx-first).
    - `ServicePlanRegistry`: `sync.RWMutex` + `map[string]resources.ServicePlan` keyed by composite
      `serviceClassName/name` via `servicePlanCompositeKey`; `NewServicePlanRegistry()`; no globals.
    - `deepCopyServicePlan` copies Parameters map, Tags slice, Metadata Labels/Annotations maps.
    - Same name under different classes stored without conflict (composite key).
    - Create → `ErrAlreadyExists`; Get → `ErrNotFound`; List → non-nil, sorted by serviceClassName
      then name; Update preserves APIVersion/Kind/Status/Metadata.Name/Spec.ServiceClassName (never
      moves a plan between classes), `ErrNotFound` if absent; Delete → `ErrNotFound` if absent.
    - `CountByServiceClass` RLock, counts entries matching `Spec.ServiceClassName`. Storage-only.
    **Tests to run:** `make fmt && make vet` (tests added in task 8).
    **Acceptance criteria:** composite identity, sort order, and CountByServiceClass correct.
    **Commit message:** `feat(registry): add ServicePlan registry (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 6.8, 6.9, 6.10_

---

- [ ] 8. Add registry unit/property/race tests
  - [ ] 8.1 ServiceClass and ServicePlan registry tests
    **Files to create:** `internal/registry/serviceclass_registry_test.go`,
    `internal/registry/serviceplan_registry_test.go`,
    `internal/registry/serviceclass_registry_property_test.go`,
    `internal/registry/serviceplan_registry_property_test.go`,
    `internal/registry/serviceclass_registry_race_test.go`,
    `internal/registry/serviceplan_registry_race_test.go`
    **Implementation notes:**
    - Unit: Create stores; duplicate → `ErrAlreadyExists` (original unchanged); Get by key; missing
      → `ErrNotFound`; List sorted; empty → non-nil `[]`; Update mutable-only; Update missing →
      `ErrNotFound`; Delete removes; Delete missing → `ErrNotFound`; ServicePlan
      `CountByServiceClass` correct; same plan name under different classes succeeds.
    - Property (Properties 4–8): Create/Get round-trip; sort invariants; deep-copy immutability;
      duplicate-create idempotent error. Tag `// Feature: serviceclass-serviceplan, Property N`.
    - Race: 10+ goroutines mixed CRUD (+ CountByServiceClass), zero race reports.
    **Tests to run:** `make fmt && make vet && make test && go test -race ./internal/registry/...`
    **Acceptance criteria:** all pass; no race reports.
    **Commit message:** `test(registry): ServiceClass/ServicePlan registry tests (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 5.9, 6.11, 14.2, 15.2, 15.3_

---

- [ ] 9. Add ServiceClass delete blocker
  - [ ] 9.1 Implement ServiceClassChildBlocker and ServicePlanChildBlockerChecker
    **Files to create:** `internal/registry/serviceclass_blocker.go`
    **Implementation notes:**
    - `ServiceClassChildBlocker` interface:
      `BlockedByServiceClassChildren(ctx, serviceClassName) ([]BlockedBy, error)`.
    - `ServicePlanChildBlockerChecker` holds a `ServicePlanRegistryIface`;
      `NewServicePlanChildBlockerChecker(reg)`.
    - Implementation calls `CountByServiceClass`; count > 0 → `[]BlockedBy{{Kind: "ServicePlan",
      Count: count}}`; else nil. Propagate registry error.
    - Reuse existing `BlockedBy` type. No generic blocker framework. Lifecycle (incl. Retired) does
      not exempt blocking.
    **Tests to run:** `make fmt && make vet` (tests added in task 10).
    **Acceptance criteria:** blocks when plans exist; nil when none; error propagates.
    **Commit message:** `feat(registry): add ServiceClass delete blocker (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 10.1, 10.2, 10.3, 10.4, 17.6_

---

- [ ] 10. Add blocker tests
  - [ ] 10.1 ServiceClassChildBlocker unit tests
    **Files to create:** `internal/registry/serviceclass_blocker_test.go`
    **Implementation notes:**
    - `CountByServiceClass` correct across seeded plans.
    - `BlockedByServiceClassChildren` returns a `ServicePlan` blocker when count > 0; nil when 0.
    - Registry error propagates.
    - Retired-lifecycle plan still blocks.
    **Tests to run:** `make fmt && make vet && make test`
    **Acceptance criteria:** all blocker tests pass.
    **Commit message:** `test(registry): ServiceClass delete blocker tests (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 10.1, 10.2, 10.3, 14.4, 17.6_

---

- [ ] 11. Add ServiceClass safe decoder
  - [ ] 11.1 Implement safeDecodeServiceClass
    **Files to create:** `internal/api/serviceclass_decode.go`
    **Implementation notes:**
    - Signature `safeDecodeServiceClass(w, r) (resources.ServiceClass, error)`.
    - `http.MaxBytesReader(w, r.Body, 1<<20)` inside the function; `*http.MaxBytesError` → 413.
    - Empty body → 400; decode into `map[string]json.RawMessage`, `status` key → 400; typed decode
      with `DisallowUnknownFields()`; unknown field/syntax/type → 400.
    - Reuse existing decoder error sentinels; do not echo raw body; 415 handled by middleware.
    **Tests to run:** `make fmt && make vet` (exercised by handler tests in task 14).
    **Acceptance criteria:** decoder matches the FEATURE-0003/0004 pattern.
    **Commit message:** `feat(api): add ServiceClass safe decoder (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 1.6, 7.7_

---

- [ ] 12. Add ServicePlan safe decoder
  - [ ] 12.1 Implement safeDecodeServicePlan
    **Files to create:** `internal/api/serviceplan_decode.go`
    **Implementation notes:**
    - Signature `safeDecodeServicePlan(w, r) (resources.ServicePlan, error)`.
    - Identical sequence to ServiceClass: 1 MiB `MaxBytesReader`, empty-body 400, `status`-key 400,
      `DisallowUnknownFields()`, unknown/syntax/type 400, oversized 413.
    - Reuse existing decoder error sentinels; do not echo raw body; 415 handled by middleware.
    **Tests to run:** `make fmt && make vet` (exercised by handler tests in task 16).
    **Acceptance criteria:** decoder matches the shared pattern.
    **Commit message:** `feat(api): add ServicePlan safe decoder (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 2.6, 8.7_

---

- [ ] 13. Add ServiceClass HTTP handler
  - [ ] 13.1 Implement ServiceClassHandler with nil-safe blocker and emitter
    **Files to create:** `internal/api/serviceclass_handler.go`
    **Implementation notes:**
    - Struct fields: `registry registry.ServiceClassRegistryIface`,
      `blocker registry.ServiceClassChildBlocker` (nil-safe),
      `emitter OperationEmitter` (nil-safe). `NewServiceClassHandler(reg, blocker, emitter)`.
    - `HandleCollection`: POST → Create, GET → List, else 405.
    - `HandleItem`: TrimPrefix `/v1/service-classes/`, split; exactly ONE non-empty segment else 404;
      GET/PUT/DELETE dispatch else 405.
    - Create: decode → validate → force apiVersion/kind/status.phase=Active → CreateServiceClass →
      `ErrAlreadyExists` 409 → emitOperation(OpCreateServiceClass, ServiceClassKind,
      ServiceClassName=name) → 201.
    - Get: path validation → GetServiceClass → 404 → 200. List: 200 `{"items": []}` sorted/empty.
    - Update: path validation → decode → body `metadata.name` present & == path (else 400
      metadata.name) → validate → UpdateServiceClass → 404 → emit OpUpdateServiceClass → 200.
    - Delete: path validation → if blocker != nil and blocked → 409 DELETE_BLOCKED ("ServicePlan")
      (blocker error → 500) → DeleteServiceClass → 404 → emit OpDeleteServiceClass → 204.
    - Use API-local `requestIDFromContext`; do NOT import `internal/server`.
    **Tests to run:** `make fmt && make vet` (tests in task 14).
    **Acceptance criteria:** flows and status codes match the design.
    **Commit message:** `feat(api): add ServiceClass HTTP handler (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 10.1, 11.1, 11.4, 11.5, 12.1, 17.3, 17.7_

---

- [ ] 14. Add ServiceClass handler tests
  - [ ] 14.1 ServiceClassHandler tests
    **Files to create:** `internal/api/serviceclass_handler_test.go`
    **Implementation notes:**
    - POST 201; 409 duplicate; 400 invalid fields, status key, bad JSON, unknown field; 413 oversized.
    - GET 200/404; GET 400 invalid `{name}` segment; wrong path shape → 404; list sorted/empty.
    - PUT 200/404; 400 metadata.name mismatch/absent.
    - DELETE 204; 404 missing; 409 DELETE_BLOCKED with a child plan present; 204 with zero plans.
    - nil emitter and nil blocker do not panic.
    **Tests to run:** `make fmt && make vet && make test`
    **Acceptance criteria:** all handler tests pass.
    **Commit message:** `test(api): ServiceClass handler tests (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 14.3, 14.4, 17.3, 17.7_

---

- [ ] 15. Add ServicePlan HTTP handler
  - [ ] 15.1 Implement ServicePlanHandler with parent lookup and nil-safe emitter
    **Files to create:** `internal/api/serviceplan_handler.go`
    **Implementation notes:**
    - Struct fields: `registry registry.ServicePlanRegistryIface`,
      `serviceClassLookup registry.ServiceClassLookup`, `emitter OperationEmitter` (nil-safe).
      `NewServicePlanHandler(reg, serviceClassLookup, emitter)`.
    - `HandleCollection`: POST → Create, GET → List, else 405.
    - `HandleItem`: TrimPrefix `/v1/service-plans/`, split; exactly TWO non-empty segments else 404;
      GET/PUT/DELETE dispatch else 405.
    - Create: decode → validate → `GetServiceClass(spec.serviceClassName)`; `ErrNotFound` → 400
      VALIDATION_FAILED field `spec.serviceClassName` → force apiVersion/kind/status.phase=Active →
      CreateServicePlan → `ErrAlreadyExists` 409 → emit OpCreateServicePlan (ServiceClassName=parent,
      ServicePlanName=name) → 201.
    - Get: path validation → GetServicePlan → 404 → 200. List: 200 sorted/empty.
    - Update: path validation → decode → body `spec.serviceClassName` == path (else 400) and
      `metadata.name` == path (else 400) → validate → verify parent still exists (else 400
      spec.serviceClassName) → UpdateServicePlan → 404 → emit OpUpdateServicePlan → 200.
    - Delete: path validation → DeleteServicePlan → 404 → emit OpDeleteServicePlan → 204.
    - Use API-local `requestIDFromContext`; do NOT import `internal/server`.
    **Tests to run:** `make fmt && make vet` (tests in task 16).
    **Acceptance criteria:** flows, parent checks, and status codes match the design.
    **Commit message:** `feat(api): add ServicePlan HTTP handler (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6, 9.1, 9.2, 9.3, 9.4, 9.5, 11.2, 11.3, 12.2, 17.2, 17.4, 17.5_

---

- [ ] 16. Add ServicePlan handler tests
  - [ ] 16.1 ServicePlanHandler tests
    **Files to create:** `internal/api/serviceplan_handler_test.go`
    **Implementation notes:**
    - POST 201; 409 duplicate composite key; 400 invalid fields, missing parent, status key, bad
      JSON, unknown field; 413 oversized.
    - GET 200/404; GET 400 invalid segment (field mapping); wrong path shape → 404; list sorted/empty.
    - PUT 200/404; 400 metadata.name or serviceClassName mismatch; 400 parent missing on update.
    - DELETE 204/404. Same plan name under two classes both succeed.
    - nil emitter does not panic.
    **Tests to run:** `make fmt && make vet && make test`
    **Acceptance criteria:** all handler tests pass.
    **Commit message:** `test(api): ServicePlan handler tests (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 14.3, 17.2, 17.4, 17.5, 17.7_

---

- [ ] 17. Add Operation emission tests for ServiceClass and ServicePlan
  - [ ] 17.1 Table-driven emission integration tests
    **Files to create:** `internal/api/serviceclass_emission_test.go`,
    `internal/api/serviceplan_emission_test.go`
    **Implementation notes:**
    - Inject a stub/real emitter; verify successful create/update/delete of each resource records the
      correct Operation type, `resourceKind`, `serviceClassName` (and `servicePlanName` for plans).
    - Failed actions (validation, duplicate, missing parent, not-found, delete-blocked) emit nothing.
    - Simulated emission failure does NOT change the primary response (201/200/204 unchanged).
    - Table-driven where practical.
    **Tests to run:** `make fmt && make vet && make test`
    **Acceptance criteria:** emission behavior matches Requirement 12; no failure leaks to response.
    **Commit message:** `test(api): catalog Operation emission tests (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 12.1, 12.2, 12.6, 12.7, 12.8, 14.5, 17.1_

---

- [ ] 18. Wire ServiceClass and ServicePlan routes into server
  - [ ] 18.1 Extend server.New and register routes
    **Files to modify:** `internal/server/server.go`
    **Implementation notes:**
    - Add `serviceClass *api.ServiceClassHandler` and `servicePlan *api.ServicePlanHandler`
      parameters to `server.New` (after `operation`, before `bootstrap`).
    - `server.New` SHALL register ServiceClass and ServicePlan routes when non-nil handlers are
      provided:
      `mux.Handle("/v1/service-classes", chain(http.HandlerFunc(serviceClass.HandleCollection)))`,
      `mux.Handle("/v1/service-classes/", chain(http.HandlerFunc(serviceClass.HandleItem)))`,
      `mux.Handle("/v1/service-plans", chain(http.HandlerFunc(servicePlan.HandleCollection)))`,
      `mux.Handle("/v1/service-plans/", chain(http.HandlerFunc(servicePlan.HandleItem)))`.
    - Middleware order unchanged: requestID → logging → contentType → handler.
    **Tests to run:** `make fmt && make vet && make build` (server tests updated in task 20).
    **Acceptance criteria:** routes registered; compiles.
    **Commit message:** `feat(server): register ServiceClass/ServicePlan routes (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 7.1, 8.1_

---

- [ ] 19. Wire registries, blocker, lookup, handlers, and emitter in main.go
  - [ ] 19.1 Production wiring in main.go
    **Files to modify:** `cmd/sovrunn-api/main.go`
    **Implementation notes:**
    - `serviceClassRegistry := registry.NewServiceClassRegistry()`.
    - `servicePlanRegistry := registry.NewServicePlanRegistry()`.
    - `serviceClassBlocker := registry.NewServicePlanChildBlockerChecker(servicePlanRegistry)`.
    - `serviceClassHandler := api.NewServiceClassHandler(serviceClassRegistry, serviceClassBlocker, emitter)`.
    - `servicePlanHandler := api.NewServicePlanHandler(servicePlanRegistry, serviceClassRegistry, emitter)`.
    - Pass both handlers into `server.New` in the correct position. Reuse the single FEATURE-0005
      `emitter`; do NOT create a duplicate Operation registry. Production wiring SHALL provide
      non-nil ServiceClass and ServicePlan handlers.
    **Tests to run:** `make fmt && make vet && make build && make run` then
    `curl -sf http://127.0.0.1:8080/readyz`.
    **Acceptance criteria:** server boots; catalog routes reachable.
    **Commit message:** `feat(cmd): wire catalog registries and handlers (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 7.1, 8.1, 9.3, 10.2, 12.7_

---

- [ ] 20. Add/update server route tests
  - [ ] 20.1 Update server constructor and add catalog route tests
    **Files to modify:** `internal/server/server_test.go`
    **Implementation notes:**
    - Update all `server.New` call sites for the new handler parameters (use real or fixture
      handlers).
    - Add route registration tests for `/v1/service-classes`, `/v1/service-classes/`,
      `/v1/service-plans`, `/v1/service-plans/` (collection + item paths reachable).
    - Confirm wrong-shape item paths return 404 through the router.
    **Tests to run:** `make fmt && make vet && make test`
    **Acceptance criteria:** server tests pass; routes verified.
    **Commit message:** `test(server): catalog route registration tests (FEATURE-0006)`
    **Do not start the next task until this passes.**
    _Requirements: 7.1, 7.2, 8.1, 8.2_

---

- [ ] 21. Final verification and cleanup
  - [ ] 21.1 Full verification, guardrail checks, and clean tree
    **Files to modify:** none (verification only).
    **Implementation notes / commands:**
    - `make fmt && make vet && make test && go test -race ./... && make build`.
    - Docker fallback if host Go is unavailable:
      `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./... && go build ./cmd/sovrunn-api'`.
    - `rm -rf bin` if the build created a `bin/` directory.
    - Confirm no FEATURE-0006 TODO placeholders remain:
      `! grep -R "TODO(FEATURE-0006)" internal cmd`.
    - Confirm `internal/api` does NOT import `internal/server`:
      `! grep -R "internal/server" internal/api`.
    - `git status` must be clean after commit (all changes committed, no stray artifacts).
    **Tests to run:** all commands above pass with zero errors.
    **Acceptance criteria:** all builds/tests/race pass; no TODO markers; no forbidden import;
    clean working tree.
    **Commit message:** `chore(feature-0006): final verification and cleanup`
    **Do not proceed past this task.**
    _Requirements: 14.6, 14.7, 16.8, 16.9_

---

## Notes

- Module path: `github.com/sanjeevksaini/sovrunn`
- Go 1.21 minimum — no Go 1.22 wildcard routing
- No new external dependencies
- ServiceClass and ServicePlan are global platform catalog resources (not scoped to the governance hierarchy)
- Registries are storage-only; parent existence checks live in the handler via `ServiceClassLookup`
- Delete blocking uses the narrow `ServiceClassChildBlocker`; no generic blocker framework
- Operation emission reuses the FEATURE-0005 nil-safe emitter; failures never affect the primary response
- Correctness Properties (from design.md) tagged `// Feature: serviceclass-serviceplan, Property N: <title>`

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1", "2.1"] },
    { "id": 1, "tasks": ["3.1", "4.1", "6.1", "7.1"] },
    { "id": 2, "tasks": ["5.1", "8.1", "9.1", "11.1", "12.1"] },
    { "id": 3, "tasks": ["10.1", "13.1", "15.1"] },
    { "id": 4, "tasks": ["14.1", "16.1", "17.1", "18.1"] },
    { "id": 5, "tasks": ["19.1"] },
    { "id": 6, "tasks": ["20.1"] },
    { "id": 7, "tasks": ["21.1"] }
  ]
}
```

**Critical path:** 1 → 6/7 → 8/9 → 13/15 → 18 → 19 → 20 → 21.
