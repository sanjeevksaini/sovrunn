# FEATURE-0015 Architecture Boundary: Canonical Cloud Model Foundation

| Field | Value |
|-------|-------|
| Status | Approved boundary (replacement under ADH-2026-045; corrected/clarified under ADH-2026-046; closed under ADH-2026-047; renamed/re-scoped, no alpha migration) |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Controlling Decisions | DEC-0037, DEC-0041, DEC-0042, DEC-0054, DEC-0059 |
| Controlling Handoffs | ADH-2026-020, ADH-2026-024, ADH-2026-025, ADH-2026-037, consolidated ADH-2026-042, ADH-2026-043 (non-migration portions), ADH-2026-045, ADH-2026-046, ADH-2026-047, ADH-2026-048, ADH-2026-049, ADH-2026-050 |
| Phase | 2R |
| Depends On | FEATURE-0011 (reuse), FEATURE-0012 (grammar/errors), FEATURE-0013 (decision/audit), FEATURE-0014 (alpha model — retained repository asset only) |

---

## 1. Purpose

Define the closed architecture boundary for FEATURE-0015 so that requirements, design, and tasks cannot absorb future-feature fields, stale FEATURE-0014 concepts, or unowned resource activation.

---

## 2. FEATURE-0015 Owned Resources and Activation Boundaries

FEATURE-0015 owns only the portions shown in the activation-boundary column. FEATURE-0015 ends at `InfrastructureStack`; it introduces no `ExecutionTarget` contract, invariant, or route (`ExecutionTarget` belongs in its entirety to FEATURE-0016 per DEC-0059/ADH-2026-045):

| Resource | Registry ID | Scope | Profile | Activation Boundary |
|----------|------------|-------|---------|---------------------|
| CloudPlatform | VS0-SCHEMA-008 | Platform | ManagedResource | full identity + phase status; immutable `spec.ownerRegistration`; PATCH-only (`application/merge-patch+json`) for `spec.description` |
| CloudProvider | VS0-SCHEMA-009 | Platform | ManagedResource | full identity + phase status; PATCH-only for `spec.displayName`/`spec.operatingMarkets` |
| CloudProviderParticipation | VS0-SCHEMA-010 | CloudPlatform | ManagedResource | full lifecycle; state machine VS0-STATE-001; no mutable spec fields — lifecycle changes occur only through explicit actions |
| HostingLocation | VS0-SCHEMA-011 | CloudProvider | ManagedResource | full identity + phase status; PATCH-only for approved mutable fields |
| Datacenter | VS0-SCHEMA-012 | CloudProvider | ManagedResource | full identity + phase status; PATCH-only for approved mutable fields |
| FaultDomain | VS0-SCHEMA-013 | CloudProvider | ManagedResource | full identity + phase status; PATCH-only for approved mutable fields |
| InfrastructureStack | VS0-SCHEMA-014 | CloudProvider | ManagedResource | full identity + phase status; PATCH-only for approved mutable fields |

---

## 3. Field Ownership Boundaries

### 3.1 ExecutionTarget — not a FEATURE-0015 contract

`ExecutionTarget` (VS0-SCHEMA-015), including its identity, schema, routes, status, writer, conformance, and target lifecycle, belongs in its entirety to FEATURE-0016 (DEC-0059; ADH-2026-045). FEATURE-0015 introduces no `ExecutionTarget` identity, schema, route, status, writer, conformance, or target lifecycle behavior, and ends at `InfrastructureStack`.

### 3.2 CloudProviderParticipation deferred FEATURE-0021 field ownership (extended by ADH-2026-047 decision 4)

| Field | introducedBy | activatedBy | Writer |
|-------|-------------|-------------|--------|
| `spec.providerSelectionModes` | **FEATURE-0021** | **FEATURE-0021** | delegated-participation-contract-authority (VS0-WRITER-004) |
| `spec.permittedHostingLocationRefs` | **FEATURE-0021** | **FEATURE-0021** | delegated-participation-contract-authority (VS0-WRITER-004) |

FEATURE-0015 must NOT store, validate, default, or mention either field as CONTRACT_ONLY behavior. Both are wholly owned by FEATURE-0021 and are rejected if present on a FEATURE-0015 participation create request (ADH-2026-047 decision 4).

The registry's `fieldOwnership.introducedBy` and `activatedBy` rule governs feature activation boundaries. It does not alter the source-of-truth precedence order: the registry remains the machine-readable expression of the higher-precedence Slice 0 specification.

---

## 4. FEATURE-0016 Delegated Activation

FEATURE-0016 owns (both introduces and activates), in full:
- ExecutionTarget identity, schema, routes, and all status fields (qualification, availability, maintenanceEpoch, factSetRef, observedGeneration, conditions)
- Target qualification, availability and maintenance-epoch state machine (VS0-STATE-004)
- NormalizedTargetFactSet (VS0-SCHEMA-016)
- TargetQualificationResult (VS0-SCHEMA-017)
- Target registration, status, adapter facts, qualification, availability, maintenance, and the complete target lifecycle

ServiceRegion remains wholly owned by FEATURE-0022. FEATURE-0016 is only a declared prerequisite for FEATURE-0022's `spec.executionTargetRefs` and `status.availability`; it does not introduce or activate either field.

---

## 5. FEATURE-0021 Delegated: Provider-Selection Intent

FEATURE-0021 owns:
- `CloudProviderParticipation.spec.providerSelectionModes` field (introducedBy/activatedBy FEATURE-0021)
- `CloudProviderParticipation.spec.permittedHostingLocationRefs` field (introducedBy/activatedBy FEATURE-0021; ADH-2026-047 decision 4)
- `providerSelectionModes` vocabulary activation in placement evaluation context
- `ServiceInstance.spec.providerPreferenceRef` evaluation semantics
- Provider-selection intent evaluation and routing logic

FEATURE-0015 does NOT introduce, store, validate, default, or reference `providerSelectionModes` or `permittedHostingLocationRefs` in any active behavior; both fields are explicitly rejected on FEATURE-0015 participation create requests.

---

## 6. Removed from Slice 0: SovrunnInstallation

`SovrunnInstallation` (former VS0-SCHEMA-057) is **removed from FEATURE-0015 scope** and from Slice 0 active schemas because:

- No Phase 2R feature owns platform lifecycle (DEC-0053 scope is deferred to Phase 3+).
- The canonical model defines SovrunnInstallation as platform-operator-facing lifecycle; it requires `PlatformLifecyclePolicy`, `SovrunnRelease`, `PlatformLifecyclePlan`, and external `PlatformLifecycleAgent` — none of which are Phase 2R scope.
- The VS-000 fixture graph referenced it for provider isolation proof, but DEC-0054's one-provider-per-participation invariant is already proven by `CloudProviderParticipation` uniqueness.

**Retirement**: VS0-SCHEMA-057 is retired as a tombstone in the retiredSchemas section of the registry. The ID is permanently consumed and must not be reused. The fixture graph uses `CloudProviderParticipation` as the participation boundary proof.

---

## 7. Canonical Bootstrap (No Alpha Runtime Migration)

### 7.1 Authority

- DEC-0059 (supersedes DEC-0058): the first executable control-plane release exposes canonical resource contracts only; FEATURE-0001–0014 are retained repository assets, not live state requiring conversion.
- ADH-2026-045: Replace alpha runtime migration with a canonical bootstrap.
- ADH-2026-041 is superseded in its entirety (migration runtime semantics). ADH-2026-043 decisions 1–3 and 6 (migration plan/record behavior) are superseded; its safe-denial, scope-reference integrity, audit reuse, and feature-local traceability decisions remain intact. ADH-2026-044 is superseded in its entirety.

### 7.2 What FEATURE-0015 Does Not Implement

FEATURE-0015 has no `CanonicalMigrationPlan`, `CanonicalMigrationRecord`, migration controller, migration milestone state machine, or migration failure mapping. It implements no alpha import, fixture conversion as a product capability, migration plans/records, migration-controller authority, migration milestones, write-freeze/cutover behavior, backup/restore migration evidence, or global migration completion. `VS0-SCHEMA-060`, `VS0-SCHEMA-061`, `VS0-WRITER-020`, `VS0-WRITER-021`, and `VS0-STATE-011` are permanently retired tombstones; their IDs must never be reused.

If historical fixture comparison remains useful during development, it is a bounded developer/test utility only — not a platform resource, runtime protocol, or migration authority.

### 7.3 Direct Canonical Creation

FEATURE-0015 owns direct creation of `CloudPlatform`, `CloudProvider`, `CloudProviderParticipation`, `HostingLocation`, `Datacenter`, `FaultDomain`, and `InfrastructureStack` as canonical resources. There is no classification/transform step, no alpha-to-canonical mapping, and no dependency on FEATURE-0014 alpha records at runtime. FEATURE-0014's alpha resources (Provider, ProviderLocation, ProviderDatacenter, DatacenterFailureDomain, InfrastructureStack) remain retained repository history assessed for reuse under FEATURE-0011.

### 7.4 Participation Lifecycle and Independent Suspension Holds

`CloudProviderParticipation.status` is api-server-owned and persists only current lifecycle facts: `phase`, `platformSuspended`, `providerSuspended`, retained `requestExpiresAt`, and standard system status fields. It does not duplicate transition history, which is durable FEATURE-0013 AuditEvent evidence. Both holds initialize to `false`. The FEATURE-0015 API server is the sole status writer for `CloudPlatform`, `CloudProvider`, `HostingLocation`, `Datacenter`, `FaultDomain`, `InfrastructureStack`, and `CloudProviderParticipation`; there is no independent topology, stack, provider, or participation controller (ADH-2026-046 decision 1). The deterministic scheduler is a system actor that invokes the api-server's expiry transition; it is not an additional status writer or controller authority.

| Action | Actor | From → To | Rule |
|---|---|---|---|
| Create participation | CloudProvider administrator | absent → Pending | `requestExpiresAt = createdAt + 7 days` |
| Accept | CloudPlatform administrator | Pending → Active | Platform's action completes the agreement |
| Reject | CloudPlatform administrator | Pending → Rejected | Terminal |
| Withdraw | CloudProvider administrator | Pending → Withdrawn | Terminal |
| Expire | deterministic scheduler | Pending → Expired | At `requestExpiresAt`; terminal, idempotent |
| Suspend | CloudPlatform or CloudProvider administrator | Active or Suspended → effective Suspended | Actor sets only its own hold true |
| Resume | CloudPlatform or CloudProvider administrator | Suspended → Active only when both holds false | Actor clears only its own hold |
| Request release | CloudPlatform administrator | Active or Suspended → Terminating | Holds retained while pending |
| Accept release | CloudProvider administrator | Terminating → Terminated | Terminal |
| Decline release | CloudProvider administrator | Terminating → prior effective Active/Suspended | Returns to Active only if both holds false |

Suspend/resume are denied with `VS0_PARTICIPATION_STATE_INVALID` while phase is Pending, Terminating, Rejected, Withdrawn, Expired, or Terminated. Exactly one non-terminal participation may exist for a `(cloudPlatformUID, cloudProviderUID)` pair; a new request after a terminal record receives a new UID. FEATURE-0015 exposes no mutable participation spec fields — lifecycle changes occur only through these explicit actions. Clearing one hold (a resume action) never reactivates the participation while the other party's hold remains true; both holds must be false for the effective state to become Active again. Each hold-changing action produces correlated AuditEvent evidence naming the acting party and resulting effective state. An Active participation is necessary but not sufficient for customer visibility; the CloudPlatform separately governs which enrolled Organizations may see or select it (FEATURE-0021/FEATURE-0018 own that eligibility and authorization).

All participation actions have an empty JSON body and require `Idempotency-Key`. `POST` create (`participation.request`) creates an absent resource: it requires `Idempotency-Key` only and does not require or accept `If-Match`; its atomic uniqueness key is the non-terminal `(cloudPlatformUID, cloudProviderUID)` pair — one concurrent create succeeds, the other returns `ALREADY_EXISTS` (409) `VS0_PARTICIPATION_DUPLICATE`, and no duplicate resource or audit event is stored. PATCH and every action on an *existing* participation (accept, reject, withdraw, suspend, resume, request-release, accept-release, decline-release) require both `If-Match` (the current `resourceVersion`) and `Idempotency-Key` (ADH-2026-046 decision 2). A missing, malformed, or non-current `If-Match` on an existing-participation action returns 412 (`STALE_RESOURCE_VERSION`), performs no lifecycle/status write, and emits no additional AuditEvent. The same principal, route, key, and request digest replayed returns the original result before version checking; a changed digest with the same key returns `CONFLICT` (409) `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`.

### 7.5 PATCH-Only Update Surface

FEATURE-0015 mutable resource updates use `PATCH` with `application/merge-patch+json` only. `PUT` and `DELETE` are not exposed and return 405. PATCH is available only for `CloudPlatform`, `CloudProvider`, and topology resources (`HostingLocation`, `Datacenter`, `FaultDomain`, `InfrastructureStack`); identity, scope, relationship references, and the embedded owner registration remain immutable. Because `CloudProviderParticipation` has no mutable F0015 spec fields, it exposes lifecycle actions rather than PATCH.

### 7.6 CloudPlatform-Root Requirement and Bootstrap Authorization

CloudPlatform exists first. A CloudProvider, `HostingLocation`, `Datacenter`, `FaultDomain`, `InfrastructureStack`, or `CloudProviderParticipation` create is denied until at least one CloudPlatform exists in the authoritative deployment control plane, returning CONFLICT (409) with `VS0_CLOUDPLATFORM_ROOT_REQUIRED`. CloudProvider exists before its topology and before participation.

The authenticated request context receives server-resolved, deterministic bootstrap grants; grants are never accepted from an HTTP header or request body, and FEATURE-0015 persists no roles, memberships, or assignments. The server-configured bootstrap principal receives `cloudplatform.write` and `cloudprovider.write` at the deployment Platform-root scope to create the first Platform-scoped resources; after creation the same action may be narrowed to one target resource UID for PATCH. `topology.write` is CloudProvider-UID scoped; participation actions are scoped to their existing CloudPlatform or CloudProvider target.

### 7.7 Closed Collection-Create Request Contract (ADH-2026-047 decision 1)

For every F0015 collection `POST`, the client supplies only `metadata.name`, the required F0015-owned `spec` fields, and the optional F0015-owned `spec` fields listed for that kind below. The client must not supply `metadata.uid`, generation, resourceVersion, timestamps, `metadata.scopeRef`, any `status` field, a field owned by another feature, or an unknown field — such requests are rejected. The API server assigns identity/version/timestamps, `metadata.scopeRef`, initial status, and required FEATURE-0013 AuditEvent evidence.

| Kind | Client-required create fields | Client-optional create fields | Server-assigned outcome |
|---|---|---|---|
| CloudPlatform | `metadata.name`; `spec.ownerRegistration.legalName`; `spec.ownerRegistration.registrationIdentifier`; `spec.ownerRegistration.jurisdictionCode` | `metadata.displayName`; `spec.description` | deployment Platform-root `scopeRef`; `status.phase=Active`; `201` resource response |
| CloudProvider | `metadata.name`; non-empty `spec.operatingMarkets[]` | `metadata.displayName`; `spec.displayName` | deployment Platform-root `scopeRef`; `status.phase=Active`; `201` resource response |
| HostingLocation | `metadata.name`; `spec.countryCode`; `spec.locality` | `spec.administrativeAreaCode`; `spec.description` | CloudProvider `scopeRef` derived under §7.8; `status.phase=Active`; `201` resource response |
| Datacenter | `metadata.name`; `spec.hostingLocationRef` | `spec.description` | CloudProvider `scopeRef` derived from the resolved parent; `status.phase=Active`; `201` resource response |
| FaultDomain | `metadata.name`; `spec.datacenterRef` | `spec.description` | CloudProvider `scopeRef` derived from the resolved parent; `status.phase=Active`; `201` resource response |
| InfrastructureStack | `metadata.name`; `spec.faultDomainRef` | `spec.description` | CloudProvider `scopeRef` derived from the resolved parent; `status.phase=Active`; `201` resource response |

All create routes require `Idempotency-Key` under ADH-2026-045/046. A successful PATCH returns `200` with the updated resource. Existing FEATURE-0012 Problem Details semantics govern rejected unknown fields, server-owned fields, unsupported media types, immutable-field attempts, and other invalid input.

### 7.8 Scope Derivation (ADH-2026-047 decision 2)

Each F0015 resource kind has exactly one server-side scope-derivation source; a client never supplies or selects `metadata.scopeRef`:

| Resource kind | Scope derivation source (single authority) |
|---|---|
| CloudPlatform | Immutable deployment Platform-root scope, server-derived. |
| CloudProvider | Immutable deployment Platform-root scope, server-derived. |
| CloudProviderParticipation | Immutable CloudPlatform scope derived from resolved `spec.cloudPlatformRef`; the referenced CloudPlatform must exist and its UID must equal the resulting `metadata.scopeRef.uid`. |
| HostingLocation | Immutable CloudProvider scope derived from the CloudProvider UID bound to the authenticated, server-resolved `topology.write` grant. A caller cannot supply or choose `metadata.scopeRef`. |
| Datacenter | Immutable CloudProvider scope derived from its resolved immutable parent reference (`spec.hostingLocationRef`). |
| FaultDomain | Immutable CloudProvider scope derived from its resolved immutable parent reference (`spec.datacenterRef`). |
| InfrastructureStack | Immutable CloudProvider scope derived from its resolved immutable parent reference (`spec.faultDomainRef`). |

Reference resolution and authorization/non-disclosure occur before structural validation; derived-scope mismatches trigger the established safe-denial/validation outcome (`VS0_AUTHORIZATION_SAFE_DENIAL` for an inaccessible reference; `VS0_TOPOLOGY_PROVIDER_MISMATCH` or `VS0_SCOPE_REFERENCE_MISMATCH` for an authorized structural mismatch).

### 7.9 Topology Immutability Correction (ADH-2026-047 decision 3)

`metadata.name` is immutable identity for every F0015 resource, including topology resources (`HostingLocation`, `Datacenter`, `FaultDomain`, `InfrastructureStack`). For those topology resources, only `spec.description` is PATCHable. Parent reference, scope, identity, all metadata system fields, and status are immutable or server-owned.

**Correction:** ADH-2026-045's topology sentence "Name and description are PATCHable" is corrected to **"Description is PATCHable; name is immutable identity."** This corrects wording only; the registry's topology mutability fields (VS0-SCHEMA-011 through VS0-SCHEMA-014) already stated identity/scope immutability and PATCH-only `spec.description` and required no field-level change.

### 7.10 Participation Collection-Create Body (ADH-2026-047 decision 4)

`POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations` is a collection-create request, not an item lifecycle action. Its body requires exactly:

- `metadata.name`
- `spec.cloudPlatformRef` (UID-pinned)
- `spec.cloudProviderRef` (UID-pinned)
- `spec.environment` with the only F0015-allowed value, `development`

It accepts no optional participation `spec` fields in F0015. In particular, `spec.providerSelectionModes` and `spec.permittedHostingLocationRefs` are FEATURE-0021-introduced and FEATURE-0021-activated fields (VS0-SCHEMA-010 `fieldOwnership`); F0015 neither accepts, stores, defaults, validates, nor exposes them.

The API server derives `metadata.scopeRef` from `spec.cloudPlatformRef`, assigns `status.phase=Pending`, `status.platformSuspended=false`, `status.providerSuspended=false`, and sets `status.requestExpiresAt=createdAt + 7 days`. It returns `201` with the created participation. The POST requires `Idempotency-Key` only and rejects/does not accept `If-Match`.

The item action routes (`:accept`, `:reject`, `:withdraw`, `:suspend`, `:resume`, `:request-release`, `:accept-release`, `:decline-release`) have an empty JSON body. Each existing-item action requires `If-Match` and `Idempotency-Key` under ADH-2026-046. The scheduler expiry is an internal API-server transition and has no public request body.

### 7.11 Audit, Correlation, and Redaction Requirements (ADH-2026-043 decision 6, preserved)

FEATURE-0015 reuses FEATURE-0013 AuditEvent. The following lifecycle and security events produce audit evidence:

| Event | Audit Evidence | Correlation |
|-------|----------------|-------------|
| CloudProviderParticipation lifecycle change (create/accept/reject/withdraw/expire/suspend/resume/request-release/accept-release/decline-release) | AuditEvent with subjectRef, actor, action, resulting effective state, timestamp | participationRef UID, requestId |
| Resource create/PATCH (CloudPlatform, CloudProvider, topology) | AuditEvent with resource UID, actor, result | resource UID, requestId |
| Cross-provider safe-denial (security) | AuditEvent with denied actor, denied action; no target existence disclosed | requestId |

Rules:
- Safe projection/redaction: audit events use the same safe-denial principle — no existence disclosure for cross-provider references.
- No secrets: no credential values, protected handles, or secret material in any audit record.
- Intentional correlation: each audit record links to its subject, actor, and governing participation through UID-pinned references.
- The mutation and its required AuditEvent are one atomic outcome; audit failure leaves the resource/lifecycle unchanged.
- A failure to append the required FEATURE-0013 AuditEvent for a FEATURE-0015 mutation returns the inherited FEATURE-0012 `INTERNAL_ERROR` Problem Details response (HTTP 500); the resource mutation and its idempotency completion are not published. `DEPENDENCY_UNAVAILABLE` (503) is not used for this outcome because FEATURE-0015 introduces no durable or external persistence dependency (ADH-2026-049).
- FEATURE-0015 does not import FEATURE-0026's integration-only trace conformance (VS0-CF-T01).

### 7.12 ISO-3166 Assigned-Code Semantics and Malformed-Input Outcomes (ADH-2026-048)

For FEATURE-0015, an ISO-3166-1 alpha-2 value (`CloudPlatform.spec.ownerRegistration.jurisdictionCode`, `CloudProvider.spec.operatingMarkets[]`, `HostingLocation.spec.countryCode`) means an **assigned** code from a fixed, repository-owned, version-pinned static dataset. The dataset contains the assigned ISO-3166-1 alpha-2 code list as of `2026-08-12`; it is compiled into or bundled with the repository, requires no network request, and introduces no third-party runtime dependency. A syntactically valid but unassigned value such as `ZZ` is rejected with the existing `VALIDATION_FAILED` (422) Problem Details contract. `HostingLocation.spec.administrativeAreaCode` remains syntax-plus-country-prefix validation only: it must match the existing ISO-3166-2 shape when present and share the same two-letter prefix as `spec.countryCode`; F0015 introduces no ISO-3166-2 membership dataset.

Four exact header/body outcomes reuse only existing FEATURE-0012 top-level Problem codes; no new top-level code or violation code is introduced:

| Input | Exact outcome |
|---|---|
| Required `Idempotency-Key` missing, empty, malformed, or over the shared length limit | `MALFORMED_REQUEST` / 400 |
| `If-Match` supplied on a participation collection-create request | `MALFORMED_REQUEST` / 400 |
| `If-Match` missing, malformed, or non-current on an existing-participation action | `STALE_RESOURCE_VERSION` / 412 (preserved ADH-2026-046 rule) |
| Participation item action with any non-empty JSON body | `MALFORMED_REQUEST` / 400 |

A malformed/prohibited input under this rule produces no mutation, no idempotency record, and no AuditEvent unless an already-approved audited-denial rule independently applies; this clarification does not broaden the ADH-2026-046/047 audit-denial list.

The following are deterministic design mechanics that the F0015 design must define exactly, without seeking new architecture authority: idempotency reservation states `InFlight`/`Completed`/`Aborted` with waiter wake-up on both terminal states and abort/release of an in-flight reservation on validation failure, stale version, audit failure, or recovered panic; idempotency records written only for completed replayable outcomes consistent with the approved audit policy; precise in-memory audit/mutation coordination wording stating that a resource/idempotency change is not published until its required audit append succeeds, with no durable cross-store transaction claim; and Go 1.22 `http.ServeMux` registration using explicit method/path patterns and `Request.PathValue` with seven collection, seven item, and eight participation-action registrations (22 total), using no wildcard, reflection, or auto-registration.

---

## 8. Explicitly Excluded (FEATURE-0015 Must Not Implement)

| Concept | Reason | Owner |
|---------|--------|-------|
| ResourcePool | Superseded by DEC-0042 | None (removed) |
| ProviderCapability | Superseded by DEC-0042 | None (removed) |
| Generic Provider | Superseded by DEC-0037 | None (split into CloudPlatform + CloudProvider) |
| CanonicalMigrationPlan, CanonicalMigrationRecord, migration controller, cutover state machine | Superseded by DEC-0059 (canonical bootstrap; no runtime alpha migration) | None (removed) |
| ExecutionTarget (identity, schema, routes, status, writer, conformance, target lifecycle) | Owned entirely by FEATURE-0016 | FEATURE-0016 |
| Target qualification/availability writes | FEATURE-0016 activation | FEATURE-0016 |
| NormalizedTargetFactSet | FEATURE-0016 | FEATURE-0016 |
| Provider-selection intent evaluation | FEATURE-0021 | FEATURE-0021 |
| CloudEnrollment | FEATURE-0021 | FEATURE-0021 |
| Individual-user/Organization visibility eligibility for participation | FEATURE-0018/0021 own authorization; not stored on CloudProviderParticipation | FEATURE-0018, FEATURE-0021 |
| Federation, remote deployment resource, signed remote offer, cross-control-plane replication | No Phase 2R owner; explicitly excluded by ADH-2026-045 decision 12 | None |
| SovrunnInstallation lifecycle | No Phase 2R owner | Deferred |
| Platform lifecycle (SovrunnRelease, PlatformLifecyclePolicy, PlatformLifecyclePlan) | DEC-0053, no Phase 2R owner | Deferred |
| ServiceClass as canonical catalog | Superseded by DEC-0049 | None |
| EffectivePolicyContext | Superseded by DEC-0050 | None |
| Six-scope vocabulary | Superseded by DEC-0037 | None |

---

## 9. Classification Semantics

| Element | Classification |
|---------|---------------|
| REQUIREMENTS-owned | Observable behavior that acceptance tests verify: resource CRUD, validation, state transitions, error codes, writer enforcement, cross-provider isolation, participation lifecycle and independent suspension holds, PATCH-only update surface |
| DESIGN-delegated | Internal mechanics: registry data structures, handler wiring, validation plumbing, in-memory store shape |
| Downstream negative boundary | ExecutionTarget owned entirely by FEATURE-0016; visible here only to forbid FEATURE-0015 activation |
| EXCLUDED | Concepts explicitly not in scope per §8 |

---

## 10. Traceability Mapping

Downstream-owned conformance (`VS0-CF-HP01`, `VS0-CF-F09`) is never used as FEATURE-0015 local acceptance evidence; they are non-owning, future integration references only, noted separately below.

| Owned Item | DEC/ADH | VS0 Schema | VS0 Writer | VS0 State | VS0 Conformance (FEATURE-0015 local) |
|------------|---------|------------|------------|-----------|--------------------------------------|
| CloudPlatform | DEC-0037; ADH-020/042/045/046/047/050 | VS0-SCHEMA-008 | VS0-WRITER-002,005 | — | VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-07, VS0-CF-F15-11, VS0-CF-F15-12, VS0-CF-F15-13, VS0-CF-F15-14, VS0-CF-F15-18, VS0-CF-F15-19, VS0-CF-F15-21, VS0-CF-F15-24, VS0-CF-F15-25, VS0-CF-F15-26, VS0-CF-F15-27 |
| CloudProvider | DEC-0037; ADH-020/042/045/046/047/050 | VS0-SCHEMA-009 | VS0-WRITER-003,005 | — | VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-07, VS0-CF-F15-12, VS0-CF-F15-13, VS0-CF-F15-14, VS0-CF-F15-18, VS0-CF-F15-19, VS0-CF-F15-21, VS0-CF-F15-24, VS0-CF-F15-25, VS0-CF-F15-26, VS0-CF-F15-27 |
| CloudProviderParticipation | DEC-0054; ADH-037/042/045/046/047/050 | VS0-SCHEMA-010 | VS0-WRITER-004,005 | VS0-STATE-001 | VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-08, VS0-CF-F15-09, VS0-CF-F15-10, VS0-CF-F15-11, VS0-CF-F15-12, VS0-CF-F15-15, VS0-CF-F15-16, VS0-CF-F15-17, VS0-CF-F15-18, VS0-CF-F15-19, VS0-CF-F15-20, VS0-CF-F15-24, VS0-CF-F15-25, VS0-CF-F15-26, VS0-CF-F15-27, VS0-CF-F15-28 |
| HostingLocation | DEC-0041; ADH-024/042/045/046/047/050 | VS0-SCHEMA-011 | VS0-WRITER-003,005 | — | VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-07, VS0-CF-F15-12, VS0-CF-F15-13, VS0-CF-F15-14, VS0-CF-F15-18, VS0-CF-F15-19, VS0-CF-F15-21, VS0-CF-F15-24, VS0-CF-F15-25, VS0-CF-F15-26, VS0-CF-F15-27 |
| Datacenter | DEC-0041; ADH-024/042/045/046/047/050 | VS0-SCHEMA-012 | VS0-WRITER-003,005 | — | VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-07, VS0-CF-F15-12, VS0-CF-F15-13, VS0-CF-F15-14, VS0-CF-F15-18, VS0-CF-F15-19, VS0-CF-F15-21, VS0-CF-F15-24, VS0-CF-F15-25, VS0-CF-F15-26, VS0-CF-F15-27 |
| FaultDomain | DEC-0041; ADH-024/042/045/046/047/050 | VS0-SCHEMA-013 | VS0-WRITER-003,005 | — | VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-07, VS0-CF-F15-12, VS0-CF-F15-13, VS0-CF-F15-14, VS0-CF-F15-18, VS0-CF-F15-19, VS0-CF-F15-21, VS0-CF-F15-24, VS0-CF-F15-25, VS0-CF-F15-26, VS0-CF-F15-27 |
| InfrastructureStack | DEC-0041,0042; ADH-025/042/045/046/047/050 | VS0-SCHEMA-014 | VS0-WRITER-003,005 | — | VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-07, VS0-CF-F15-12, VS0-CF-F15-13, VS0-CF-F15-14, VS0-CF-F15-18, VS0-CF-F15-19, VS0-CF-F15-21, VS0-CF-F15-24, VS0-CF-F15-25, VS0-CF-F15-26, VS0-CF-F15-27 |
| Cross-provider isolation | DEC-0037,0054; ADH-043/045 | — | VS0-WRITER-003 | — | VS0-CF-X03, VS0-CF-F15-23 |
| Scope-reference integrity | DEC-0037,0054; ADH-043/045/047 | VS0-SCHEMA-008,010 | VS0-WRITER-002,004 | — | VS0-CF-F15-11, VS0-CF-F15-27 |
| Status writer resolution (api-server sole writer) | ADH-2026-046 decision 1 | VS0-SCHEMA-008..010 | VS0-WRITER-005 | — | VS0-CF-F15-12, VS0-CF-F15-15, VS0-CF-F15-16 |
| Participation create-versus-existing preconditions; all-route idempotency replay and mismatch coverage | ADH-2026-046 decision 2; ADH-2026-050 | VS0-SCHEMA-008..014 | VS0-WRITER-002,003,004 | VS0-STATE-001 | VS0-CF-F15-15, VS0-CF-F15-17, VS0-CF-F15-18, VS0-CF-F15-19 |
| Bootstrap-grant boundary | ADH-2026-045; ADH-2026-046 decision 3 | — | VS0-WRITER-002,003,004 | — | VS0-CF-F15-22 |
| Audit atomicity; required-AuditEvent-append failure mapping | DEC-0059; ADH-2026-043 (preserved),046,049 | VS0-SCHEMA-007 | VS0-WRITER-011 (F0013) | — | VS0-CF-F15-24 |
| Closed collection-create request contract | ADH-2026-047 decision 1 | VS0-SCHEMA-008..014 | VS0-WRITER-002,003,004 | — | VS0-CF-F15-26 |
| Scope derivation single-source proof | ADH-2026-047 decision 2 | VS0-SCHEMA-008..014 | VS0-WRITER-002,003,004 | — | VS0-CF-F15-27 |
| Topology immutability correction (name immutable; description PATCHable) | ADH-2026-047 decision 3 | VS0-SCHEMA-011..014 | VS0-WRITER-003 | — | VS0-CF-F15-13, VS0-CF-F15-14 |
| Participation collection-create body vs. empty item-action body | ADH-2026-047 decision 4 | VS0-SCHEMA-010 | VS0-WRITER-004 | VS0-STATE-001 | VS0-CF-F15-28 |
| ISO-3166-1 alpha-2 assigned-code semantics | ADH-2026-048 decision 1 | VS0-SCHEMA-008,009,011 | — | — | VS0-CF-F15-29 |
| Idempotency-Key/If-Match/item-action-body malformed-input outcomes | ADH-2026-048 decision 2 | VS0-SCHEMA-010 | VS0-WRITER-004 | — | VS0-CF-F15-30 |

Downstream integration references (non-owning, FEATURE-0015 does not claim this evidence as local): `VS0-CF-HP01` (FEATURE-0026 cross-feature integration), `VS0-CF-F09` (FEATURE-0023 placement).

### 10.1 REQ-F15 and AC-F15 to F0015-local-proof mapping (ADH-2026-046 decision 3; extended by ADH-2026-047 decision 5)

Every `REQ-F15-01` through `REQ-F15-23` and every `AC-F15-01` through `AC-F15-19` maps to at least one F0015-local conformance case, except the explicit anti-drift/non-runtime verification items labelled below. No row cites a FEATURE-0016+ conformance case as required local evidence.

| REQ ID | F0015-local proof case(s) |
|--------|----------------------------|
| REQ-F15-01 | VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-12, VS0-CF-F15-13, VS0-CF-F15-14 |
| REQ-F15-02 | VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-12, VS0-CF-F15-13, VS0-CF-F15-14 |
| REQ-F15-03 | VS0-CF-F15-01, VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-15, VS0-CF-F15-20 |
| REQ-F15-04 | VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-13, VS0-CF-F15-14, VS0-CF-F15-23 |
| REQ-F15-05 | VS0-CF-F15-01, VS0-CF-F15-25 |
| REQ-F15-06 | VS0-CF-F15-08, VS0-CF-F15-10, VS0-CF-F15-15, VS0-CF-F15-20 |
| REQ-F15-07 | VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-15, VS0-CF-F15-17, VS0-CF-F15-18, VS0-CF-F15-19, VS0-CF-F15-20 |
| REQ-F15-08 | VS0-CF-F15-07, VS0-CF-F15-13, VS0-CF-F15-14 |
| REQ-F15-09 | VS0-CF-X03, VS0-CF-F15-23 |
| REQ-F15-10 | VS0-CF-F15-01, VS0-CF-F15-02 |
| REQ-F15-11 | VS0-CF-F15-05, VS0-CF-F15-06, VS0-CF-F15-14 |
| REQ-F15-12 | Anti-drift/non-runtime — no new top-level Problem code is introduced; proven collectively by every error-emitting F0015-local case (VS0-CF-F15-02, VS0-CF-F15-05, VS0-CF-F15-06, VS0-CF-F15-07, VS0-CF-X03, VS0-CF-F15-12, VS0-CF-F15-14, VS0-CF-F15-17, VS0-CF-F15-19, VS0-CF-F15-20), each of which uses only an existing FEATURE-0012 code. Not a standalone runtime-observable scenario. |
| REQ-F15-13 | **Explicit anti-drift/non-runtime verification gate.** "No `CanonicalMigrationPlan`/`CanonicalMigrationRecord`/migration controller/cutover state machine exists" is an anti-drift claim proven by absence (retired tombstones VS0-SCHEMA-060/061, VS0-WRITER-020/021, VS0-STATE-011 never reused), not a runtime-observable proof case. Excluded from the proof-case mapping per ADH-2026-046 decision 3. |
| REQ-F15-14 | VS0-CF-F15-11 |
| REQ-F15-15 | VS0-CF-F15-24 |
| REQ-F15-16 | VS0-CF-F15-12 |
| REQ-F15-17 | VS0-CF-F15-22 |
| REQ-F15-18 | VS0-CF-F15-18, VS0-CF-F15-19 |
| REQ-F15-19 | VS0-CF-F15-26 |
| REQ-F15-20 | VS0-CF-F15-27 |
| REQ-F15-21 | VS0-CF-F15-28 |
| REQ-F15-22 | VS0-CF-F15-29 |
| REQ-F15-23 | VS0-CF-F15-30 |

| AC ID | F0015-local proof case(s) |
|-------|----------------------------|
| AC-F15-01 | VS0-CF-F15-01, VS0-CF-F15-02 |
| AC-F15-02 | VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-09 |
| AC-F15-03 | VS0-CF-F15-01, VS0-CF-F15-25 |
| AC-F15-04 | VS0-CF-X03, VS0-CF-F15-23 |
| AC-F15-05 | **Explicit anti-drift/non-runtime verification gate.** Same rationale as REQ-F15-13; excluded from the proof-case mapping. |
| AC-F15-06 | VS0-CF-F15-07 |
| AC-F15-07 | **Explicit anti-drift/non-runtime verification gate.** The canonical stale-concept anti-drift gate is a documentation/registry consistency check, not a runtime-observable scenario. Excluded from the proof-case mapping. |
| AC-F15-08 | VS0-CF-F15-05, VS0-CF-F15-06, VS0-CF-X03 |
| AC-F15-09 | Contract-reuse claim proven collectively by every error-emitting F0015-local case (VS0-CF-F15-02, VS0-CF-F15-05, VS0-CF-F15-06, VS0-CF-F15-07, VS0-CF-X03, VS0-CF-F15-12, VS0-CF-F15-14, VS0-CF-F15-17, VS0-CF-F15-19, VS0-CF-F15-20), each of which uses only an existing FEATURE-0012 Problem Details code. Not a standalone runtime-observable scenario. |
| AC-F15-10 | VS0-CF-F15-08, VS0-CF-F15-10 |
| AC-F15-11 | VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-15, VS0-CF-F15-17, VS0-CF-F15-18, VS0-CF-F15-19 |
| AC-F15-12 | VS0-CF-F15-11 |
| AC-F15-13 | VS0-CF-F15-24 |
| AC-F15-14 | VS0-CF-F15-12 |
| AC-F15-15 | VS0-CF-F15-26 |
| AC-F15-16 | VS0-CF-F15-27 |
| AC-F15-17 | VS0-CF-F15-28 |
| AC-F15-18 | VS0-CF-F15-29 |
| AC-F15-19 | VS0-CF-F15-30 |

---

## 11. Anti-Drift Rules

1. FEATURE-0015 requirements/design/tasks must NOT introduce or use any superseded concept listed in §8 as active behavior. Those names may appear only in the explicitly labelled exclusion/non-goal context that explains why they are prohibited.
2. FEATURE-0015 must NOT introduce, initialize, persist, validate, default, or write any `ExecutionTarget` identity, schema, route, status, writer, conformance, or target lifecycle field — FEATURE-0016 owns them in their entirety.
3. FEATURE-0015 must NOT evaluate provider-selection intent — that is FEATURE-0021 scope.
4. FEATURE-0015 must NOT implement platform lifecycle (SovrunnInstallation, SovrunnRelease, PlatformLifecyclePolicy, PlatformLifecyclePlan).
5. FEATURE-0015 must NOT reintroduce `CanonicalMigrationPlan`, `CanonicalMigrationRecord`, a migration controller, or any cutover state machine (DEC-0059; ADH-2026-045). `VS0-SCHEMA-060`, `VS0-SCHEMA-061`, `VS0-WRITER-020`, `VS0-WRITER-021`, and `VS0-STATE-011` are permanently retired and must never be reused.
6. FEATURE-0015 exposes mutable updates only through `PATCH` with `application/merge-patch+json`; it must never expose `PUT` or `DELETE`.
7. FEATURE-0015 must not store individual-user or Organization-level participation visibility eligibility on `CloudProviderParticipation`; that belongs to FEATURE-0021 (Organization eligibility) and FEATURE-0018 (individual authorization).
8. FEATURE-0015 must not accept, store, default, validate, or expose a client-supplied `metadata.scopeRef` for any owned resource kind; every scope is server-derived from exactly one registered source per §7.8 (ADH-2026-047 decision 2).
9. FEATURE-0015 must not describe topology `metadata.name` as PATCHable; only `spec.description` is PATCHable for topology resources (ADH-2026-047 decision 3).
10. FEATURE-0015 must not accept, store, default, validate, or expose `spec.providerSelectionModes` or `spec.permittedHostingLocationRefs` on a `CloudProviderParticipation` create request; both are FEATURE-0021-introduced and FEATURE-0021-activated (ADH-2026-047 decision 4).
11. FEATURE-0015 must treat an ISO-3166-1 alpha-2 value as valid only if it is an assigned code from the fixed, repository-owned, version-pinned dataset; a syntactically valid but unassigned code must be rejected with VALIDATION_FAILED, and `administrativeAreaCode` must never require an ISO-3166-2 membership dataset (ADH-2026-048 decision 1). FEATURE-0015 must not introduce a new top-level Problem code or violation code for the four malformed-input outcomes in §7.12 (ADH-2026-048 decision 2).

---

*End of architecture boundary.*
