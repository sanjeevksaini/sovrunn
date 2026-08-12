# FEATURE-0015: Canonical Cloud Model Foundation

| Field | Value |
|-------|-------|
| Status | Approved scope (replacement under ADH-2026-045; corrected/clarified under ADH-2026-046; renamed/re-scoped, no alpha migration) |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase | 2R |
| Order | 5 (first unimplemented Phase 2R feature) |
| Depends On | FEATURE-0011, FEATURE-0012, FEATURE-0013, FEATURE-0014 (retained repository asset only) |
| Depended On By | FEATURE-0016, FEATURE-0021, FEATURE-0022 |
| Architecture Boundary | docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md |
| Controlling Decisions | DEC-0037, DEC-0041, DEC-0042, DEC-0054, DEC-0059 |
| Controlling Handoffs | ADH-2026-020, ADH-2026-024, ADH-2026-025, ADH-2026-037, ADH-2026-042, ADH-2026-043 (non-migration portions), ADH-2026-045, ADH-2026-046 |

---

## 1. Feature Summary

Establish the seven-scope canonical cloud model identity layer through direct canonical bootstrap. Register CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, and InfrastructureStack as managed resources, created directly rather than migrated. FEATURE-0015 ends at InfrastructureStack; ExecutionTarget belongs in its entirety to FEATURE-0016. Per DEC-0059 (superseding DEC-0058), there is no `CanonicalMigrationPlan`, `CanonicalMigrationRecord`, migration controller, or cutover state machine — FEATURE-0001 through FEATURE-0014 are retained repository assets, not live state requiring conversion.

---

## 2. Scope Classification

### 2.1 REQUIREMENTS-Owned Observable Behavior

| ID | Behavior | DEC/ADH | VS0 IDs |
|----|----------|---------|---------|
| REQ-F15-01 | Register CloudPlatform as a Platform-scoped resource with a unique `metadata.name` and immutable `spec.ownerRegistration` (legalName, registrationIdentifier, jurisdictionCode only; no user/Organization/role/grant reference); validate uniqueness, scope, immutable identity | DEC-0037; ADH-020,045 | VS0-SCHEMA-008 |
| REQ-F15-02 | Register CloudProvider as a Platform-scoped resource with unique `metadata.name` and non-empty, unique, uppercase-ISO-3166-1-alpha2 `spec.operatingMarkets[]`; validate scope and identity | DEC-0037; ADH-020 | VS0-SCHEMA-009 |
| REQ-F15-03 | Register CloudProviderParticipation as a CloudPlatform-scoped resource linking one CloudPlatform and one CloudProvider; enforce exactly one non-terminal participation per platform+provider pair; lifecycle Pending→Active/Rejected/Withdrawn/Expired, Active↔Suspended, Active/Suspended→Terminating→Terminated | DEC-0054; ADH-037,045 | VS0-SCHEMA-010, VS0-STATE-001 |
| REQ-F15-04 | Register topology chain under CloudProvider scope: HostingLocation (uppercase ISO 3166-1 alpha-2 `countryCode`, non-empty `locality`, optional ISO 3166-2 `administrativeAreaCode`) → Datacenter → FaultDomain → InfrastructureStack, each with an immutable parent TypedRef; validate containment refs and same-provider UID preservation through the chain | DEC-0041; ADH-024,045 | VS0-SCHEMA-011..014 |
| REQ-F15-05 | FEATURE-0015 ends at InfrastructureStack. It introduces no ExecutionTarget identity, schema, route, status, writer, conformance, or target lifecycle behavior — ExecutionTarget belongs in its entirety to FEATURE-0016 | DEC-0059; ADH-045 | — |
| REQ-F15-06 | CloudProviderParticipation lifecycle: explicit request/accept/reject/withdraw/expire actions; Suspended is a derived effective state carried by two independent holds (platformSuspended, providerSuspended), each settable/clearable only by its respective administrator; clearing one hold never reactivates while the other remains true; terminal states are Rejected, Withdrawn, Expired, Terminated | DEC-0054; ADH-037,045 | VS0-SCHEMA-010, VS0-STATE-001 |
| REQ-F15-07 | FEATURE-0015 exposes no mutable CloudProviderParticipation spec fields; create (participation.request) requires Idempotency-Key only, no If-Match; every action on an existing participation (accept, reject, withdraw, suspend, resume, request-release, accept-release, decline-release) requires both If-Match and Idempotency-Key | DEC-0054; ADH-045,046 | VS0-SCHEMA-010, VS0-STATE-001 |
| REQ-F15-08 | Mutable CloudPlatform, CloudProvider, and topology resource updates use PATCH with application/merge-patch+json only; PUT and DELETE are not exposed and return 405 | DEC-0059; ADH-045 | VS0-SCHEMA-008..014 |
| REQ-F15-09 | Cross-provider isolation: deny cross-provider target reference; use safe RESOURCE_NOT_FOUND without existence disclosure | DEC-0037,0054 | VS0-CF-X03 |
| REQ-F15-10 | Use only the seven canonical scope kinds globally; CloudPlatform and CloudProvider accept only Platform scope, CloudProviderParticipation only CloudPlatform scope, and topology resources only CloudProvider scope | DEC-0037 | VS0-SCHEMA-001,008..014 |
| REQ-F15-11 | Writer enforcement: only cloud-provider-admin may write topology spec; only cloud-platform-admin may write CloudPlatform spec | DEC-0037 | VS0-WRITER-002,003 |
| REQ-F15-12 | Existing FEATURE-0012 Problem Details codes used for all errors; no new top-level codes | — | VS0-SCHEMA-004 |
| REQ-F15-13 | No CanonicalMigrationPlan, CanonicalMigrationRecord, migration controller, or cutover state machine is implemented or reintroduced; FEATURE-0001–0014 are retained repository assets, not live state requiring conversion | DEC-0059; ADH-045 | (retired: VS0-SCHEMA-060,061; VS0-WRITER-020,021; VS0-STATE-011) |
| REQ-F15-14 | Scope reference UID invariant: CloudProviderParticipation `spec.cloudPlatformRef.uid` must equal its CloudPlatform `metadata.scopeRef.uid`. CloudPlatform itself has no Organization reference to validate (ADH-045 decision 8: immutable `spec.ownerRegistration` replaces `ownerOrganizationRef`). Mismatch returns VALIDATION_FAILED (422) with VS0_SCOPE_REFERENCE_MISMATCH | DEC-0037,0054; ADH-043 (preserved),045 | VS0-SCHEMA-010, VS0-CF-F15-11 |
| REQ-F15-15 | Audit evidence: participation lifecycle/hold-changing actions and resource create/PATCH, and safe-denial/security denials produce AuditEvent (reusing FEATURE-0013) with intentional correlation fields, safe projection/redaction, and no secrets; mutation and its required AuditEvent are one atomic outcome | DEC-0059; ADH-043 (preserved) | VS0-SCHEMA-007 |
| REQ-F15-16 | Creation order: a CloudProvider, HostingLocation, Datacenter, FaultDomain, InfrastructureStack, or CloudProviderParticipation create is denied until at least one CloudPlatform exists, returning CONFLICT (409) with VS0_CLOUDPLATFORM_ROOT_REQUIRED | ADH-045 | — |
| REQ-F15-17 | Bootstrap authorization: server-resolved, deterministic grants only; grants never accepted from header/body; F0015 persists no roles, memberships, or assignments | ADH-045 | — |
| REQ-F15-18 | Idempotency: POST create/actions require Idempotency-Key; same principal+route+key+digest returns original result; different digest with same key returns CONFLICT/409 with VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH | ADH-045 | — |

### 2.2 DESIGN-Delegated Mechanics

- In-memory registry data structures for cloud model resources
- Handler wiring and router setup
- Validation plumbing (deterministic, pure)
- Bootstrap grant resolution internals

### 2.3 Downstream Concepts Visible Only as Negative Boundaries

These concepts appear in the shared Slice 0 registry only so FEATURE-0015 can prove it does not activate them. They are not FEATURE-0015 requirements or persistence fields.

| Element | Introduced and Activated By | FEATURE-0015 Rule |
|---------|-----------------------------|------------------|
| ExecutionTarget (identity, schema, routes, status, writer, conformance, target lifecycle) | FEATURE-0016 | Must not introduce, persist, default, validate, or write any part of the ExecutionTarget contract |

### 2.4 EXCLUDED (Must NOT Implement)

| Concept | Reason | Reference |
|---------|--------|-----------|
| ResourcePool | Superseded | DEC-0042 |
| ProviderCapability | Superseded | DEC-0042 |
| Generic Provider (combined) | Superseded | DEC-0037 |
| CanonicalMigrationPlan, CanonicalMigrationRecord, migration controller, cutover state machine | Superseded by canonical bootstrap; no runtime alpha migration | DEC-0059; ADH-045 |
| ExecutionTarget (any part) | Owned entirely by FEATURE-0016 | ADH-045 decision 9 |
| SovrunnInstallation | No Phase 2R owner | DEC-0053 deferred |
| CloudProviderParticipation `spec.providerSelectionModes` | Introduced and activated only by FEATURE-0021 | VS0-SCHEMA-010 field ownership |
| Platform lifecycle resources | No Phase 2R owner | DEC-0053 deferred |
| Target qualification writes | FEATURE-0016 | DEC-0042 |
| NormalizedTargetFactSet | FEATURE-0016 | DEC-0036 |
| CloudEnrollment | FEATURE-0021 | DEC-0038 |
| Provider-selection intent evaluation | FEATURE-0021 | — |
| End-user CloudProvider projection; participation visibility eligibility | FEATURE-0021 (Organization eligibility), FEATURE-0018 (individual authorization) | ADH-045 |
| Federation, remote deployment resource, signed remote offer, cross-control-plane replication | Explicitly excluded single-authority onboarding model | ADH-045 decision 12 |
| ServiceClass | Superseded | DEC-0049 |
| EffectivePolicyContext | Superseded | DEC-0050 |
| Six-scope vocabulary | Superseded | DEC-0037 |
| Real provisioning | Phase 3 | PHASE2R non-goals |

---

## 3. Acceptance Criteria

| ID | Criterion | Conformance ID |
|----|-----------|----------------|
| AC-F15-01 | CloudPlatform, CloudProvider, topology chain CRUD works with validation and their registry-declared scope subsets; FEATURE-0015 ends at InfrastructureStack | VS0-CF-F15-01, VS0-CF-F15-02 |
| AC-F15-02 | CloudProviderParticipation lifecycle, provider request plus CloudPlatform acceptance guard, and pair uniqueness are enforced | VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-09 |
| AC-F15-03 | No ExecutionTarget identity, schema, route, status, writer, conformance, or target lifecycle behavior is introduced by FEATURE-0015 | VS0-CF-F15-01, VS0-CF-F15-25 |
| AC-F15-04 | Cross-provider target reference denied safely | VS0-CF-X03, VS0-CF-F15-23 |
| AC-F15-05 | No CanonicalMigrationPlan, CanonicalMigrationRecord, migration controller, or cutover state machine exists in active FEATURE-0015 behavior | — (anti-drift/non-runtime; excluded per ADH-2026-046 decision 3) |
| AC-F15-06 | PUT and DELETE return 405 for CloudPlatform, CloudProvider, and topology resources; PATCH accepts only application/merge-patch+json | VS0-CF-F15-07 |
| AC-F15-07 | Every output passes the canonical stale-concept anti-drift gate | Anti-drift (non-runtime; excluded per ADH-2026-046 decision 3) |
| AC-F15-08 | Writer enforcement: unauthorized writes to topology/platform specs, status, and system-owned fields denied; cross-provider references deny safely | VS0-CF-F15-05, VS0-CF-F15-06, VS0-CF-X03 |
| AC-F15-09 | All errors use FEATURE-0012 Problem Details with existing codes | — (contract reuse; proven collectively by every error-emitting local case) |
| AC-F15-10 | Independent suspension holds behave correctly: each administrator may set/clear only its own hold; clearing one hold does not reactivate while the other remains true; suspend/resume denied in terminal/Pending/Terminating phases | VS0-CF-F15-08, VS0-CF-F15-10 |
| AC-F15-11 | CloudProviderParticipation exposes no mutable spec fields; create requires Idempotency-Key only; lifecycle changes on an existing participation occur only through explicit actions requiring If-Match and Idempotency-Key | VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-15, VS0-CF-F15-17, VS0-CF-F15-18, VS0-CF-F15-19 |
| AC-F15-12 | Scope reference UID invariant enforced: cloudPlatformRef UID equals scopeRef UID for CloudProviderParticipation; mismatch rejected with VS0_SCOPE_REFERENCE_MISMATCH | VS0-CF-F15-11 |
| AC-F15-13 | Participation lifecycle/hold-changing actions, resource create/PATCH, and security denials produce correlated, redacted AuditEvent with no secrets | VS0-CF-F15-24 |
| AC-F15-14 | CloudPlatform-root requirement enforced: create of CloudProvider/topology/participation denied with CONFLICT/409 VS0_CLOUDPLATFORM_ROOT_REQUIRED until at least one CloudPlatform exists | VS0-CF-F15-12 |

---

## 4. Canonical Bootstrap (No Alpha Runtime Migration)

Per DEC-0059 (superseding DEC-0058) and ADH-2026-045, FEATURE-0015 has no `CanonicalMigrationPlan`, `CanonicalMigrationRecord`, migration controller, or cutover state machine. `VS0-SCHEMA-060`, `VS0-SCHEMA-061`, `VS0-WRITER-020`, `VS0-WRITER-021`, and `VS0-STATE-011` are permanently retired tombstones and must never be reused.

Sovrunn has no live control plane, customer data, production API estate, or persisted alpha state to convert. FEATURE-0001 through FEATURE-0014 remain retained repository assets and reuse input under FEATURE-0011; they create no active runtime data-migration obligation. FEATURE-0015 creates `CloudPlatform`, `CloudProvider`, `CloudProviderParticipation`, `HostingLocation`, `Datacenter`, `FaultDomain`, and `InfrastructureStack` directly as canonical resources — there is no classification/transform step and no dependency on FEATURE-0014 alpha records at runtime.

If historical fixture comparison remains useful during development, it is a bounded developer/test utility only: it is neither a platform resource nor a runtime protocol, accepts no user data, has no active API, and creates no migration authority.

### 4.1 CloudProviderParticipation Lifecycle and Independent Suspension (VS0-STATE-001)

`CloudProviderParticipation.status` is api-server-owned and persists only current lifecycle facts: `phase`, `platformSuspended`, `providerSuspended`, retained `requestExpiresAt`, and standard system status fields — it does not duplicate transition history, which is durable FEATURE-0013 AuditEvent evidence. Both holds initialize to `false`. The FEATURE-0015 API server is the sole status writer for all F0015-owned resources including `CloudProviderParticipation`; there is no separate participation controller. The deterministic scheduler is a system actor that invokes the api-server's expiry transition, not an additional writer authority (ADH-2026-046 decision 1).

```text
absent
  → Pending (create; participation.request)
      → Active (accept; participation.accept.platform)
      → Rejected (reject; participation.reject.platform) [terminal]
      → Withdrawn (withdraw; participation.withdraw.provider) [terminal]
      → Expired (scheduler at requestExpiresAt) [terminal]
Active/Suspended
  → effective Suspended (suspend; either administrator sets only its own hold true)
  → Active (resume; actor clears only its own hold; both holds must be false)
  → Terminating (request-release; participation.request-release.platform)
Terminating
  → Terminated (accept-release; participation.accept-release.provider) [terminal]
  → prior effective Active/Suspended (decline-release; participation.decline-release.provider)
```

For accepted participation, the server stores `phase=Active` exactly when both holds are false and `phase=Suspended` when either hold is true. Pending, Terminating, and every terminal phase override the hold-derived phase; the holds remain retained but cannot be changed outside Active or Suspended. Suspend/resume are denied with the participation-state conflict outcome (CONFLICT/409, `VS0_PARTICIPATION_STATE_INVALID`) while phase is Pending, Terminating, Rejected, Withdrawn, Expired, or Terminated. Exactly one non-terminal participation may exist for a `(cloudPlatformUID, cloudProviderUID)` pair; a new request after a terminal record receives a new UID.

FEATURE-0015 exposes no mutable participation spec fields; all lifecycle changes occur only through the explicit actions above. Every action has an empty JSON body and requires `Idempotency-Key`. Create (`participation.request`) targets an absent resource: it requires `Idempotency-Key` only and does not require or accept `If-Match`; its atomic uniqueness key is the non-terminal `(cloudPlatformUID, cloudProviderUID)` pair (ADH-2026-046 decision 2). Every action on an *existing* participation additionally requires `If-Match`; the same key and request returns the original result before version checking, while a new key with a stale version returns 412 (STALE_RESOURCE_VERSION) with no lifecycle/status write and no additional AuditEvent. Each hold-changing action produces correlated AuditEvent evidence naming the acting party and resulting effective state.

An Active `CloudProviderParticipation` is necessary but not sufficient for customer visibility. The CloudPlatform decides which enrolled customer Organizations may see or select an Active participation (FEATURE-0021 owns Organization eligibility; FEATURE-0018 owns individual authorization); this is not stored on `CloudProviderParticipation`. Pending, Rejected, Withdrawn, Expired, Suspended, Terminating, and Terminated participations are never end-user visible.

Audit correlation for participation/hold-changing actions and resource create/PATCH reuses FEATURE-0013 AuditEvent with intentional correlation fields (requestId, resource/participation UID), safe projection/redaction, and no secrets. This feature does not import FEATURE-0026's integration-only trace conformance (VS0-CF-T01).

### 4.2 PATCH-Only Update Surface

FEATURE-0015 mutable resource updates use `PATCH` with `application/merge-patch+json`. `PUT` and `DELETE` are not exposed by FEATURE-0015 and return 405. PATCH is available only for `CloudPlatform`, `CloudProvider`, and topology resources (`HostingLocation`, `Datacenter`, `FaultDomain`, `InfrastructureStack`). Identity, scope, relationship references, and the embedded owner registration remain immutable. Because `CloudProviderParticipation` has no mutable F0015 spec fields, it exposes lifecycle actions rather than PATCH.

### 4.3 CloudPlatform-Root Requirement

CloudPlatform exists first. A CloudProvider, HostingLocation, Datacenter, FaultDomain, InfrastructureStack, or CloudProviderParticipation create is denied until at least one CloudPlatform exists, returning CONFLICT (409) with `VS0_CLOUDPLATFORM_ROOT_REQUIRED`. CloudProvider exists before its topology and before participation.

### 4.4 Bootstrap Authorization

The authenticated request context receives server-resolved, deterministic bootstrap grants; grants are never accepted from an HTTP header or request body, and FEATURE-0015 persists no roles, memberships, or assignments. The server-configured bootstrap principal receives `cloudplatform.write` and `cloudprovider.write` at the deployment Platform-root scope to create the first Platform-scoped resources. `topology.write` is CloudProvider-UID scoped; participation actions are scoped to their existing CloudPlatform or CloudProvider target.

---

## 5. Error Codes Used

| Scenario | Code | HTTP | Violation |
|----------|------|------|-----------|
| Missing/invalid auth | AUTH_REQUIRED | 401 | VS0_AUTH_REQUIRED |
| Unauthorized write | AUTHORIZATION_DENIED | 403 | — |
| Cross-provider ref denied | RESOURCE_NOT_FOUND | 404 | VS0_AUTHORIZATION_SAFE_DENIAL |
| Duplicate name in scope | ALREADY_EXISTS | 409 | — |
| Invalid participation transition | CONFLICT | 409 | VS0_PARTICIPATION_STATE_INVALID |
| Duplicate CloudProviderParticipation pair | ALREADY_EXISTS | 409 | VS0_PARTICIPATION_DUPLICATE |
| Missing delegated acceptance for Pending → Active | CONFLICT | 409 | VS0_PARTICIPATION_ACCEPTANCE_MISSING |
| Client writes status | AUTHORIZATION_DENIED | 403 | VS0_STATUS_FIELD_WRITE |
| Client writes system-owned metadata | AUTHORIZATION_DENIED | 403 | VS0_SYSTEM_OWNED_FIELD_WRITE |
| Resource uses a scope outside its declared subset | VALIDATION_FAILED | 422 | VS0_SCOPE_KIND_INVALID |
| PUT or DELETE requested against a FEATURE-0015 resource | — | 405 | — |
| Immutable field patched | VALIDATION_FAILED | 422 | VS0_PATCH_IMMUTABLE_FIELD |
| Create before any CloudPlatform exists | CONFLICT | 409 | VS0_CLOUDPLATFORM_ROOT_REQUIRED |
| Idempotency key reused with a changed request | CONFLICT | 409 | VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH |
| Invalid reference chain | VALIDATION_FAILED | 422 | — |
| Stale resourceVersion / stale If-Match | STALE_RESOURCE_VERSION | 412 | — |
| Scope reference UID mismatch (cloudPlatformRef UID ≠ scopeRef UID) | VALIDATION_FAILED | 422 | VS0_SCOPE_REFERENCE_MISMATCH |

---

## 6. Phase 2R Exit Evidence (FEATURE-0015 contribution)

1. CloudPlatform, CloudProvider, and topology resources are created directly with no migration classification/transform step.
2. Provider credential isolation proven by cross-provider denial tests.
3. Seven canonical scope kinds used exclusively.
4. No ResourcePool, ProviderCapability, generic Provider, or CanonicalMigrationPlan/Record in active code/schemas.
5. CloudProviderParticipation independent suspension holds and lifecycle actions behave deterministically.

---

## 7. Next Feature Boundary

FEATURE-0016 owns ExecutionTarget registration and its complete target lifecycle (identity, schema, routes, status, writer, conformance) in its entirety; FEATURE-0015 introduces none of it.

---

*Executable scope only. No requirements.md, design.md, or tasks.md generation from this file. Those are produced by Kiro spec stages under .kiro/specs/.*
