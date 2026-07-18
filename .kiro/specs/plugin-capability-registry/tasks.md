# Implementation Plan: FEATURE-0007 Plugin and Capability Registry

## Overview

Implement `Plugin` and `Capability` as **global platform registry resources**. Registry
declarations only — no execution, runtime loading, or dynamic plugin discovery. In-memory
storage-only registries; simple `metadata.name` identity for both; Plugin delete blocked while
child Capabilities exist; Capability immutable (no PUT); lifecycle Operations emitted via the
FEATURE-0005 nil-safe emitter; reference validation at write time via narrow lookup interfaces.

## Cursor Execution Rule

Implement **one numbered task at a time**. After each task:
1. Run the verification commands listed in the task.
2. Confirm all pass with zero errors before moving to the next task.
3. **Do not start the next task** until the current task's verification succeeds.
4. If a task fails verification, fix within the current task before proceeding.

## Global Constraints

- Go 1.21-compatible only; no Go 1.22 wildcard routing; no external dependencies.
- `internal/api` MUST NOT import `internal/server`.
- Plugin and Capability are global platform resources — NOT scoped to Organization,
  OrganizationUnit, Tenant, or Project.
- Registries are storage-only; no cross-registry dependencies. Reference existence checks happen
  in the handler via narrow lookup interfaces (`ServiceClassLookup`, `PluginLookup`).
- Reuse the existing FEATURE-0005 `OperationEmitter` and nil-safe `emitOperation` helper.
- Operation emission failure MUST NOT change the primary API response.
- No secrets stored in Plugin or Capability resources.
- Existing handler tests must pass with a nil emitter and nil blocker.
- Capability does NOT support PUT (update); `PUT /v1/capabilities/{name}` → 405.
- Module path: `github.com/sanjeevksaini/sovrunn`.

## Tasks

- [ ] 1. Add Plugin and Capability resource models and constants
  - [ ] 1.1 Create Plugin and Capability structs, type/mode/operation constants, kind constants
    **Files to create:** `internal/resources/plugin.go`, `internal/resources/capability.go`
    **Implementation notes:**
    - `Plugin` (APIVersion, Kind, Metadata, Spec, Status) with standard JSON tags.
    - `PluginSpec`: PluginType, Version, ServiceClassRefs (`[]string`), DeploymentMode,
      Description (omitempty), Tags (omitempty, `[]string`).
    - `PluginStatus`: Phase, Message (omitempty).
    - `Capability` (APIVersion, Kind, Metadata, Spec, Status) with standard JSON tags.
    - `CapabilitySpec`: PluginRef, ServiceClassRef, Operation, Supported (`bool`),
      Description (omitempty).
    - `CapabilityStatus`: Phase, Message (omitempty).
    - Constants: `PluginKind = "Plugin"`, `CapabilityKind = "Capability"`,
      `PluginAPIVersion = "platform.sovrunn.io/v1alpha1"`,
      `CapabilityAPIVersion = "platform.sovrunn.io/v1alpha1"`.
    - PluginType constants (10): dStoreOps, cacheOps, streamOps, objectOps, gatewayOps,
      faasOps, lbOps, k8sOps, bigDataOps, sdeOps.
    - DeploymentMode constant (1): compiled-in.
    - CapabilityOperation constants (13): Validate, Plan, Provision, Configure, Bind,
      Observe, Scale, Upgrade, Backup, Restore, RotateCredentials, Unbind, Delete.
    - Reuse existing `Metadata` struct and `PhaseActive` phase constant.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** structs and constants compile; JSON tags match Requirements 1–2.
    **Commit message:** `feat(resources): add Plugin and Capability models (FEATURE-0007)`
    _Requirements: 1.1–1.5, 2.1–2.5_

---

- [ ] 2. Add Operation constants, spec fields, and METHOD_NOT_ALLOWED error code
  - [ ] 2.1 Extend operation.go and errors.go with FEATURE-0007 additions
    **Files to modify:** `internal/resources/operation.go`, `internal/resources/errors.go`
    **Implementation notes:**
    - Add 5 operation type constants: `OpCreatePlugin`, `OpUpdatePlugin`, `OpDeletePlugin`,
      `OpCreateCapability`, `OpDeleteCapability`.
    - Add two optional `OperationSpec` fields: `PluginName string json:"pluginName,omitempty"`
      and `CapabilityName string json:"capabilityName,omitempty"`.
    - Add `ErrCodeMethodNotAllowed ErrorCode = "METHOD_NOT_ALLOWED"` to errors.go.
    - Do NOT change existing fields, tags, or constants.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** new constants/fields present; existing tests unaffected.
    **Commit message:** `feat(resources): add plugin Operation types, spec fields, METHOD_NOT_ALLOWED (FEATURE-0007)`
    _Requirements: 13.3–13.6, 8.9_

---

- [ ] 3. Add Plugin validation
  - [ ] 3.1 Implement ValidatePlugin and ValidatePluginPathSegment
    **Files to create:** `internal/validation/plugin.go`
    **Implementation notes:**
    - `ValidatePlugin(p resources.Plugin) []resources.FieldError` — pure, context-free.
    - `metadata.name` DNS-label (1–63) → field `metadata.name`.
    - `spec.pluginType` required, ∈ {dStoreOps, cacheOps, streamOps, objectOps, gatewayOps,
      faasOps, lbOps, k8sOps, bigDataOps, sdeOps} → field `spec.pluginType`.
    - `spec.version` required, non-empty → field `spec.version`.
    - `spec.serviceClassRefs` required, non-nil, len ≥ 1 → field `spec.serviceClassRefs`.
    - Each entry in `spec.serviceClassRefs` must be DNS-label (1–63) → field `spec.serviceClassRefs`.
    - `spec.deploymentMode` required, ∈ {compiled-in} → field `spec.deploymentMode`.
    - `spec.description` and `spec.tags` optional; not format-validated.
    - `ValidatePluginPathSegment(name string) []resources.FieldError` — DNS-label check on name.
    - Reuse `validateName`/`dnsLabelRe` from the validation package; enum checks via switch.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** valid inputs pass; invalid pluginType/deploymentMode/name/version/
    serviceClassRefs produce the correct field errors.
    **Commit message:** `feat(validation): add Plugin validation (FEATURE-0007)`
    _Requirements: 3.1–3.8_

---

- [ ] 4. Add Capability validation
  - [ ] 4.1 Implement ValidateCapability and ValidateCapabilityPathSegment
    **Files to create:** `internal/validation/capability.go`
    **Implementation notes:**
    - `ValidateCapability(c resources.Capability) []resources.FieldError` — pure, context-free.
    - `metadata.name` DNS-label (1–63) → field `metadata.name`.
    - `spec.pluginRef` required, DNS-label (1–63) → field `spec.pluginRef`.
    - `spec.serviceClassRef` required, DNS-label (1–63) → field `spec.serviceClassRef`.
    - `spec.operation` required, ∈ {Validate, Plan, Provision, Configure, Bind, Observe, Scale,
      Upgrade, Backup, Restore, RotateCredentials, Unbind, Delete} → field `spec.operation`.
    - `spec.supported` defaults to `false` (Go zero-value); no validation error for absent boolean.
    - `spec.description` optional; not format-validated.
    - `ValidateCapabilityPathSegment(name string) []resources.FieldError` — DNS-label check.
    - Reuse `validateName`/`dnsLabelRe` from the validation package; enum checks via switch.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** valid inputs pass; invalid pluginRef/serviceClassRef/operation/name
    produce the correct field errors.
    **Commit message:** `feat(validation): add Capability validation (FEATURE-0007)`
    _Requirements: 4.1–4.8_

---

- [ ] 5. Add validation unit and property tests
  - [ ] 5.1 Plugin and Capability validation tests
    **Files to create:** `internal/validation/plugin_test.go`,
    `internal/validation/capability_test.go`,
    `internal/validation/plugin_property_test.go`,
    `internal/validation/capability_property_test.go`
    **Implementation notes:**
    - Plugin: valid accepted; empty/invalid/long `metadata.name` rejected; invalid/empty
      pluginType rejected; empty version rejected; empty/nil serviceClassRefs rejected; invalid
      entries in serviceClassRefs rejected; invalid deploymentMode rejected; path-segment
      validation.
    - Capability: valid accepted; empty/invalid/long `metadata.name` rejected; empty/invalid
      `spec.pluginRef` rejected; empty/invalid `spec.serviceClassRef` rejected; invalid
      operation rejected; path-segment validation.
    - Property tests use `testing/quick` `Config{MaxCount: 100}`, tagged
      `// Feature: plugin-capability-registry, Property N: <title>`:
      - Property 1: valid DNS-label names with valid enum values accepted (Plugin).
      - Property 2: arbitrary invalid strings rejected for Plugin name.
      - Property 3: valid pluginType values accepted; invalid values rejected.
      - Property 4: valid deploymentMode values accepted; invalid values rejected.
      - Property 5: valid DNS-label names with valid enum operation accepted (Capability).
      - Property 6: arbitrary invalid strings rejected for Capability name/pluginRef/serviceClassRef.
      - Property 7: valid operation values accepted; invalid values rejected.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** all validation unit and property tests pass deterministically.
    **Commit message:** `test(validation): Plugin/Capability validation tests (FEATURE-0007)`
    _Requirements: 15.1, 16.1, 16.2_

---

- [ ] 6. Add Plugin registry
  - [ ] 6.1 Implement PluginRegistry, PluginRegistryIface, and PluginLookup interface
    **Files to create:** `internal/registry/plugin_registry.go`
    **Implementation notes:**
    - `PluginRegistryIface`: Create/Get/List/Update/Delete (all ctx-first).
    - `PluginRegistry`: `sync.RWMutex` + `map[string]resources.Plugin` keyed by
      `metadata.name`; `NewPluginRegistry()`; no global state.
    - `PluginLookup` interface: `GetPlugin(ctx context.Context, name string) (resources.Plugin, error)`.
      The concrete `*PluginRegistry` satisfies it. Defined in same file.
    - `deepCopyPlugin`: copies ServiceClassRefs slice, Tags slice, Metadata.Labels and
      Metadata.Annotations maps.
    - Create → `ErrAlreadyExists` on dup; Get → `ErrNotFound`; List → non-nil, sorted by name;
      Update preserves APIVersion/Kind/Status/Metadata.Name, replaces mutable fields
      (Labels, Annotations, PluginType, Version, ServiceClassRefs, DeploymentMode, Description,
      Tags); Delete → `ErrNotFound` if absent.
    - Storage-only; no dependency on other registries; no package-level global state.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** methods behave per design; returns deep copies.
    **Commit message:** `feat(registry): add Plugin registry (FEATURE-0007)`
    _Requirements: 5.1–5.9_

---

- [ ] 7. Add Capability registry
  - [ ] 7.1 Implement CapabilityRegistry, CapabilityRegistryIface, and CountByPlugin
    **Files to create:** `internal/registry/capability_registry.go`
    **Implementation notes:**
    - `CapabilityRegistryIface`: Create/Get/List/Delete + `CountByPlugin` (all ctx-first).
      No Update method (Capability is immutable).
    - `CapabilityRegistry`: `sync.RWMutex` + `map[string]resources.Capability` keyed by
      `metadata.name`; `NewCapabilityRegistry()`; no global state.
    - `deepCopyCapability`: copies Metadata.Labels and Metadata.Annotations maps.
      CapabilitySpec contains only scalar fields so no slice/map copy needed for spec.
    - Create → `ErrAlreadyExists` on dup; Get → `ErrNotFound`; Delete → `ErrNotFound` if absent.
    - `ListCapabilities(ctx, pluginRef, serviceClassRef string) ([]resources.Capability, error)`:
      when both non-empty → AND filter; when either empty → not applied; returns non-nil,
      sorted by `metadata.name`.
    - `CountByPlugin(ctx, pluginName string) (int, error)`: RLock; iterate and count where
      `Spec.PluginRef == pluginName`. Used by the delete blocker.
    - Storage-only; no dependency on other registries; no package-level global state.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** methods behave per design; returns deep copies; filtering works.
    **Commit message:** `feat(registry): add Capability registry (FEATURE-0007)`
    _Requirements: 6.1–6.11_

---

- [ ] 8. Add registry unit, property, and race tests
  - [ ] 8.1 Plugin and Capability registry tests
    **Files to create:** `internal/registry/plugin_registry_test.go`,
    `internal/registry/capability_registry_test.go`,
    `internal/registry/plugin_registry_property_test.go`,
    `internal/registry/capability_registry_property_test.go`,
    `internal/registry/plugin_registry_race_test.go`,
    `internal/registry/capability_registry_race_test.go`
    **Implementation notes:**
    - Plugin unit: Create stores; duplicate → `ErrAlreadyExists` (original unchanged); Get by key;
      missing → `ErrNotFound`; List sorted; empty → non-nil `[]`; Update mutable fields only;
      Update preserves immutable fields; Update missing → `ErrNotFound`; Delete removes;
      Delete missing → `ErrNotFound`; deep-copy immutability (mutating returned value does not
      change stored state).
    - Capability unit: Create stores; duplicate → `ErrAlreadyExists` (original unchanged);
      Get by key; missing → `ErrNotFound`; List sorted; empty → non-nil `[]`; Delete removes;
      Delete missing → `ErrNotFound`; CountByPlugin correct; List with pluginRef filter;
      List with serviceClassRef filter; List with both filters (AND); List with no filters (all);
      deep-copy immutability.
    - Property tests (Properties 8–12): Create/Get round-trip preserves data; List sort invariant;
      deep-copy immutability; duplicate-create idempotent error; Capability filter correctness.
      Tag `// Feature: plugin-capability-registry, Property N: <title>`.
    - Race tests: 10+ goroutines mixed CRUD (+ CountByPlugin for Capability), zero race reports.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./internal/registry/...'`
    **Acceptance criteria:** all pass; no race reports.
    **Commit message:** `test(registry): Plugin/Capability registry tests (FEATURE-0007)`
    _Requirements: 5.9, 6.11, 15.2, 16.1, 16.2_

---

- [ ] 9. Add Plugin delete blocker
  - [ ] 9.1 Implement PluginChildBlocker interface and CapabilityChildBlockerChecker
    **Files to create:** `internal/registry/plugin_blocker.go`
    **Implementation notes:**
    - `PluginChildBlocker` interface:
      `BlockedByPluginChildren(ctx context.Context, pluginName string) ([]BlockedBy, error)`.
    - `CapabilityChildBlockerChecker` holds a `CapabilityRegistryIface`;
      `NewCapabilityChildBlockerChecker(reg CapabilityRegistryIface)`.
    - Implementation calls `CountByPlugin`; count > 0 → `[]BlockedBy{{Kind: "Capability",
      Count: count}}`; else nil. Propagate registry error.
    - Reuse existing `BlockedBy` type. No generic blocker framework. Consistent with
      `ServicePlanChildBlockerChecker` pattern from FEATURE-0006.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** blocks when capabilities exist; nil when none; error propagates.
    **Commit message:** `feat(registry): add Plugin delete blocker (FEATURE-0007)`
    _Requirements: 11.1–11.4_

---

- [ ] 10. Add blocker tests
  - [ ] 10.1 PluginChildBlocker unit tests
    **Files to create:** `internal/registry/plugin_blocker_test.go`
    **Implementation notes:**
    - `CountByPlugin` correct across seeded capabilities.
    - `BlockedByPluginChildren` returns a `Capability` blocker when count > 0; nil when 0.
    - Registry error propagates.
    - Multiple capabilities referencing same plugin returns correct count.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** all blocker tests pass.
    **Commit message:** `test(registry): Plugin delete blocker tests (FEATURE-0007)`
    _Requirements: 11.1–11.4, 15.4_

---

- [ ] 11. Add Plugin safe decoder
  - [ ] 11.1 Implement safeDecodePlugin
    **Files to create:** `internal/api/plugin_decode.go`
    **Implementation notes:**
    - Signature: `safeDecodePlugin(w http.ResponseWriter, r *http.Request) (resources.Plugin, error)`.
    - `http.MaxBytesReader(w, r.Body, 1<<20)` inside the function.
    - `*http.MaxBytesError` → errBodyTooLarge (413).
    - Empty body → errEmptyBody (400).
    - Decode into `map[string]json.RawMessage`; if `status` key present → errStatusFieldPresent (400).
    - Typed decode with `DisallowUnknownFields()` into `resources.Plugin`.
    - Unknown field → errUnknownField (400); syntax/type → errMalformedJSON (400).
    - Reuse existing decoder error sentinels; do not echo raw body; 415 handled by middleware.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** decoder matches the existing FEATURE-0003/0004/0006 pattern.
    **Commit message:** `feat(api): add Plugin safe decoder (FEATURE-0007)`
    _Requirements: 1.6, 7.7_

---

- [ ] 12. Add Capability safe decoder
  - [ ] 12.1 Implement safeDecodeCapability
    **Files to create:** `internal/api/capability_decode.go`
    **Implementation notes:**
    - Signature: `safeDecodeCapability(w http.ResponseWriter, r *http.Request) (resources.Capability, error)`.
    - Identical sequence to Plugin decoder: 1 MiB `MaxBytesReader`, empty-body 400,
      `status`-key 400, `DisallowUnknownFields()`, unknown/syntax/type 400, oversized 413.
    - Reuse existing decoder error sentinels; do not echo raw body; 415 handled by middleware.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** decoder matches the shared pattern.
    **Commit message:** `feat(api): add Capability safe decoder (FEATURE-0007)`
    _Requirements: 2.6, 8.10_

---

- [ ] 13. Add Plugin HTTP handler
  - [ ] 13.1 Implement PluginHandler with reference validation, nil-safe blocker and emitter
    **Files to create:** `internal/api/plugin_handler.go`
    **Implementation notes:**
    - Struct fields: `registry registry.PluginRegistryIface`,
      `serviceClassLookup registry.ServiceClassLookup` (verifies serviceClassRefs),
      `blocker registry.PluginChildBlocker` (nil-safe),
      `emitter OperationEmitter` (nil-safe).
      `NewPluginHandler(reg, serviceClassLookup, blocker, emitter)`.
    - `HandleCollection`: POST → Create, GET → List, else 405.
    - `HandleItem`: TrimPrefix `/v1/plugins/`, split; exactly ONE non-empty segment else 404;
      GET/PUT/DELETE dispatch else 405.
    - Create: decode → validate → for each ref in serviceClassRefs verify via
      serviceClassLookup.GetServiceClass (ErrNotFound → 400 VALIDATION_FAILED
      field="spec.serviceClassRefs") → force apiVersion/kind/status.phase=Active →
      CreatePlugin → `ErrAlreadyExists` 409 → emitOperation(OpCreatePlugin, PluginKind,
      PluginName=name) → 201.
    - Get: path validation → GetPlugin → 404 → 200.
    - List: ListPlugins → 200 `{"items": [...]}` sorted; empty → `{"items": []}`.
    - Update: path validation → decode → body `metadata.name` present & == path (else 400
      metadata.name) → validate → verify all serviceClassRefs → UpdatePlugin → 404 →
      emit OpUpdatePlugin → 200.
    - Delete: path validation → if blocker != nil and blocked → 409 DELETE_BLOCKED
      ("deletion blocked by Capability resources") (blocker error → 500) → DeletePlugin →
      404 → emit OpDeletePlugin → 204.
    - Use API-local `requestIDFromContext`; do NOT import `internal/server`.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** flows and status codes match the design; nil blocker/emitter safe.
    **Commit message:** `feat(api): add Plugin HTTP handler (FEATURE-0007)`
    _Requirements: 7.1–7.7, 9.1–9.4, 11.1–11.4, 12.1–12.4, 13.1, 13.7, 13.8_

---

- [ ] 14. Add Plugin handler tests
  - [ ] 14.1 PluginHandler unit tests
    **Files to create:** `internal/api/plugin_handler_test.go`
    **Implementation notes:**
    - POST 201; 409 duplicate; 400 invalid fields (missing name, invalid pluginType, empty version,
      empty serviceClassRefs, invalid deploymentMode); 400 status key; 400 bad JSON; 400 unknown
      field; 413 oversized; 400 missing ServiceClass ref.
    - GET 200/404; GET 400 invalid `{name}` path segment; wrong path shape (multi-segment) → 404;
      list sorted/empty.
    - PUT 200/404; 400 metadata.name mismatch/absent; 400 invalid fields; 400 missing
      ServiceClass ref on update.
    - DELETE 204/404; DELETE 409 DELETE_BLOCKED with capabilities present; DELETE 204 with zero
      capabilities.
    - nil emitter and nil blocker do not panic.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** all handler tests pass.
    **Commit message:** `test(api): Plugin handler tests (FEATURE-0007)`
    _Requirements: 15.3, 15.4, 15.5_

---

- [ ] 15. Add Capability HTTP handler
  - [ ] 15.1 Implement CapabilityHandler with reference validation and nil-safe emitter
    **Files to create:** `internal/api/capability_handler.go`
    **Implementation notes:**
    - Struct fields: `registry registry.CapabilityRegistryIface`,
      `pluginLookup registry.PluginLookup` (verifies pluginRef),
      `serviceClassLookup registry.ServiceClassLookup` (verifies serviceClassRef),
      `emitter OperationEmitter` (nil-safe).
      `NewCapabilityHandler(reg, pluginLookup, serviceClassLookup, emitter)`.
    - `HandleCollection`: POST → Create, GET → List, else 405.
    - `HandleItem`: TrimPrefix `/v1/capabilities/`, split; exactly ONE non-empty segment else 404;
      GET → Get, PUT → 405 METHOD_NOT_ALLOWED ("Capability does not support update; delete and
      recreate instead"), DELETE → Delete, else 405.
    - Create: decode → validate → pluginLookup.GetPlugin(pluginRef) (ErrNotFound → 400
      field="spec.pluginRef") → serviceClassLookup.GetServiceClass(serviceClassRef) (ErrNotFound →
      400 field="spec.serviceClassRef") → verify serviceClassRef ∈ plugin.Spec.ServiceClassRefs
      (not found → 400 field="spec.serviceClassRef" message="ServiceClass <ref> is not declared
      by Plugin <pluginRef>") → force apiVersion/kind/status.phase=Active → CreateCapability →
      `ErrAlreadyExists` 409 → emitOperation(OpCreateCapability, CapabilityKind,
      PluginName=pluginRef, CapabilityName=name) → 201.
    - Get: path validation → GetCapability → 404 → 200.
    - List: read query params `pluginRef` and `serviceClassRef` → ListCapabilities(ctx,
      pluginRef, serviceClassRef) → 200 `{"items": [...]}` sorted; empty → `{"items": []}`.
    - Delete: path validation → GetCapability (404 if not found) → DeleteCapability (404 if
      race, unlikely) → emitOperation(OpDeleteCapability, CapabilityKind,
      PluginName=cap.Spec.PluginRef, CapabilityName=name) → 204.
    - Use API-local `requestIDFromContext`; do NOT import `internal/server`.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** flows and status codes match the design; PUT → 405; nil emitter safe.
    **Commit message:** `feat(api): add Capability HTTP handler (FEATURE-0007)`
    _Requirements: 8.1–8.10, 10.1–10.6, 13.2, 13.5–13.8, 18.8_

---

- [ ] 16. Add Capability handler tests
  - [ ] 16.1 CapabilityHandler unit tests
    **Files to create:** `internal/api/capability_handler_test.go`
    **Implementation notes:**
    - POST 201; 409 duplicate; 400 invalid fields (missing name, invalid pluginRef, invalid
      serviceClassRef, invalid operation); 400 status key; 400 bad JSON; 400 unknown field;
      413 oversized; 400 missing Plugin ref; 400 missing ServiceClass ref; 400 ServiceClass not
      in Plugin serviceClassRefs.
    - GET 200/404; GET 400 invalid `{name}` path segment; wrong path shape (multi-segment) → 404;
      list sorted/empty; list with query filters (pluginRef only, serviceClassRef only, both,
      neither).
    - DELETE 204/404.
    - PUT → 405 METHOD_NOT_ALLOWED (regardless of existence).
    - nil emitter does not panic.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** all handler tests pass.
    **Commit message:** `test(api): Capability handler tests (FEATURE-0007)`
    _Requirements: 15.3, 15.4, 15.5_

---

- [ ] 17. Add Operation emission tests for Plugin and Capability
  - [ ] 17.1 Table-driven emission integration tests
    **Files to create:** `internal/api/plugin_emission_test.go`,
    `internal/api/capability_emission_test.go`
    **Implementation notes:**
    - Inject a stub/real emitter; verify successful create/update/delete of Plugin records the
      correct Operation type (`OpCreatePlugin`, `OpUpdatePlugin`, `OpDeletePlugin`),
      `resourceKind = "Plugin"`, `pluginName` = plugin name.
    - Verify successful create/delete of Capability records the correct Operation type
      (`OpCreateCapability`, `OpDeleteCapability`), `resourceKind = "Capability"`,
      `pluginName` = referenced plugin name, `capabilityName` = capability name.
    - Failed actions (validation, duplicate, missing reference, not-found, delete-blocked,
      method-not-allowed) emit nothing.
    - Simulated emission failure does NOT change the primary response (201/200/204 unchanged).
    - Table-driven where practical.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** emission behavior matches Requirements 13; no failure leaks to response.
    **Commit message:** `test(api): Plugin/Capability Operation emission tests (FEATURE-0007)`
    _Requirements: 13.1–13.8, 15.5_

---

- [ ] 18. Wire Plugin and Capability routes into server
  - [ ] 18.1 Extend server.New and register routes
    **Files to modify:** `internal/server/server.go`
    **Implementation notes:**
    - Add `plugin *api.PluginHandler` and `capability *api.CapabilityHandler` parameters to
      `server.New` (after `servicePlan`, before `bootstrap`).
    - Register Plugin and Capability routes when non-nil handlers are provided:
      `mux.Handle("/v1/plugins", chain(http.HandlerFunc(plugin.HandleCollection)))`,
      `mux.Handle("/v1/plugins/", chain(http.HandlerFunc(plugin.HandleItem)))`,
      `mux.Handle("/v1/capabilities", chain(http.HandlerFunc(capability.HandleCollection)))`,
      `mux.Handle("/v1/capabilities/", chain(http.HandlerFunc(capability.HandleItem)))`.
    - Middleware order unchanged: requestID → logging → contentType → handler.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** routes registered; compiles; existing server tests pass
    (update constructor call sites in tests).
    **Commit message:** `feat(server): register Plugin/Capability routes (FEATURE-0007)`
    _Requirements: 7.1, 7.2, 8.1, 8.2_

---

- [ ] 19. Wire registries, blocker, lookups, handlers, and emitter in main.go
  - [ ] 19.1 Production wiring in main.go
    **Files to modify:** `cmd/sovrunn-api/main.go`
    **Implementation notes:**
    - `pluginRegistry := registry.NewPluginRegistry()`.
    - `capabilityRegistry := registry.NewCapabilityRegistry()`.
    - `pluginBlocker := registry.NewCapabilityChildBlockerChecker(capabilityRegistry)`.
    - `pluginHandler := api.NewPluginHandler(pluginRegistry, serviceClassRegistry, pluginBlocker, emitter)`.
    - `capabilityHandler := api.NewCapabilityHandler(capabilityRegistry, pluginRegistry, serviceClassRegistry, emitter)`.
    - Pass both handlers into `server.New` in the correct position.
    - Reuse the single FEATURE-0005 `emitter` and the existing `serviceClassRegistry`;
      do NOT create duplicates.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go build ./cmd/sovrunn-api'`
    **Acceptance criteria:** server compiles and boots; plugin/capability routes reachable.
    **Commit message:** `feat(cmd): wire plugin registries and handlers (FEATURE-0007)`
    _Requirements: 7.1, 8.1, 9.3, 10.5, 11.2, 13.7_

---

- [ ] 20. Add/update server route tests
  - [ ] 20.1 Update server constructor and add plugin/capability route tests
    **Files to modify:** `internal/server/server_test.go`
    **Implementation notes:**
    - Update all `server.New` call sites for the new handler parameters (use real or fixture
      handlers).
    - Add route registration tests for `/v1/plugins`, `/v1/plugins/`, `/v1/capabilities`,
      `/v1/capabilities/` (collection + item paths reachable).
    - Confirm wrong-shape item paths return 404 through the router.
    **Tests to run:** `docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'`
    **Acceptance criteria:** server tests pass; routes verified.
    **Commit message:** `test(server): Plugin/Capability route registration tests (FEATURE-0007)`
    _Requirements: 7.1, 7.2, 8.1, 8.2_

---

- [ ] 21. Final verification and cleanup
  - [ ] 21.1 Full verification, guardrail checks, and clean tree
    **Files to modify:** none (verification only).
    **Implementation notes / commands:**
    - Full Docker verification:
      ```bash
      docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./... && go build ./cmd/sovrunn-api'
      ```
    - Artifact cleanup:
      ```bash
      rm -f sovrunn-api
      rm -rf bin
      ```
    - Confirm no FEATURE-0007 TODO placeholders remain:
      ```bash
      ! grep -rn "TODO(FEATURE-0007)" internal cmd
      ```
    - Confirm `internal/api` does NOT import `internal/server`:
      ```bash
      ! grep -rn "internal/server" internal/api
      ```
    - `git status` must be clean after commit (all changes committed, no stray artifacts).
    **Tests to run:** all commands above pass with zero errors.
    **Acceptance criteria:** all builds/tests/race pass; no TODO markers; no forbidden import;
    clean working tree; no binary artifacts.
    **Commit message:** `chore(feature-0007): final verification and cleanup`
    **Do not proceed past this task.**
    _Requirements: 15.6, 15.7, 17.1–17.12_

---

## Notes

- Module path: `github.com/sanjeevksaini/sovrunn`
- Go 1.21 minimum — no Go 1.22 wildcard routing
- No new external dependencies
- Plugin and Capability are global platform resources (not scoped to the governance hierarchy)
- Registries are storage-only; reference existence checks live in the handler via narrow interfaces
- Plugin delete blocking uses the narrow `PluginChildBlocker`; no generic blocker framework
- Capability is immutable (no PUT/update); PUT returns 405 METHOD_NOT_ALLOWED
- Operation emission reuses the FEATURE-0005 nil-safe emitter; failures never affect the primary response
- Capability list supports optional `pluginRef` and `serviceClassRef` query parameter filters
- Same `metadata.name` allowed for both a Plugin and a Capability (separate registries)
- Duplicate entries in `spec.serviceClassRefs` are accepted (no deduplication enforced)
- Correctness Properties tagged `// Feature: plugin-capability-registry, Property N: <title>`

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

**Critical path:** 1 → 6/7 → 9 → 13/15 → 18 → 19 → 20 → 21.
