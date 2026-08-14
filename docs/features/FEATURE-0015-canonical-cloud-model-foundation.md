---
doc_type: feature
id: FEATURE-0015
title: Canonical Cloud Model Foundation
status: approved
phase: 2R
reuse_assessment_format_version: 1.0.0
canonical_architecture: docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md
kiro_slug: canonical-cloud-model-and-alpha-migration
---

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
| Controlling Handoffs | ADH-2026-020, ADH-2026-024, ADH-2026-025, ADH-2026-037, ADH-2026-042, ADH-2026-043 (non-migration portions), ADH-2026-045, ADH-2026-046, ADH-2026-047, ADH-2026-048, ADH-2026-049, ADH-2026-050, ADH-2026-051, ADH-2026-052, ADH-2026-053, ADH-2026-054, ADH-2026-055, ADH-2026-056, ADH-2026-057 |

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
| REQ-F15-15 | Audit evidence: after authentication and current authorization are evaluated, a redacted FEATURE-0013 AuditEvent (with request correlation and no secret or inaccessible-resource disclosure) is produced exactly for: successful resource collection create or PATCH; successful participation create/action or deterministic scheduler expiry; an authenticated client attempt to write api-server-owned status or identity/metadata; a CloudPlatform-root creation denial; an authenticated collection LIST without a read grant; an authenticated `X-Sovrunn-Bootstrap-Grant` header or valid duplicate-free top-level `bootstrapGrant` body member; and an authenticated inaccessible cross-provider/safe-denial reference. No durable AuditEvent is produced for missing/invalid authentication, malformed/prohibited body/header/field input other than that exact forged-grant carrier, unsupported method/media type, stale If-Match, same-key replay, changed-digest idempotency conflict, pair-uniqueness conflict, or an invalid participation source state (redacted logs/metrics under the existing observability baseline are not AuditEvents). Mutation and its required AuditEvent are one atomic outcome; a failure to append a required AuditEvent returns the inherited INTERNAL_ERROR (500) Problem Details response with safe non-disclosure, and the resource mutation, lifecycle/status transition, and idempotency completion are not published (unpublished); DEPENDENCY_UNAVAILABLE is not used for this outcome | DEC-0059; ADH-043 (preserved),049,051,054 | VS0-SCHEMA-007 |
| REQ-F15-16 | Creation order: a CloudProvider, HostingLocation, Datacenter, FaultDomain, InfrastructureStack, or CloudProviderParticipation create is denied until at least one CloudPlatform exists, returning CONFLICT (409) with VS0_CLOUDPLATFORM_ROOT_REQUIRED | ADH-045 | — |
| REQ-F15-17 | Bootstrap authorization: server-resolved, deterministic grants only; grants never accepted from header/body; only X-Sovrunn-Bootstrap-Grant and valid duplicate-free top-level bootstrapGrant are audited forged-grant carriers with the ADH-2026-054 precedence; F0015 persists no roles, memberships, or assignments | ADH-045,054 | — |
| REQ-F15-18 | Idempotency: POST create/actions require Idempotency-Key; item actions bind the namespace to the concrete target UID; current authorization and safe access precede a completed replay; same principal+route+target (when applicable)+key+digest returns original result; different digest with same key returns CONFLICT/409 with VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH | ADH-045,054 | — |
| REQ-F15-19 | Closed collection-create request contract: for every F0015 collection POST, the client supplies only metadata.name, the required kind-specific spec fields, and the registered optional kind-specific spec fields; the server rejects metadata.uid/generation/resourceVersion/timestamps, metadata.scopeRef, any status field, a field owned by another feature, or an unknown field; every successful create returns 201 with its exact initial status; a successful PATCH returns 200 with the updated resource | ADH-047 decision 1 | VS0-SCHEMA-008..014 |
| REQ-F15-20 | Scope derivation: each F0015 resource kind has exactly one server-side scope-derivation source — CloudPlatform and CloudProvider from the deployment Platform-root; CloudProviderParticipation from its resolved spec.cloudPlatformRef (UID-equality enforced); HostingLocation from the CloudProvider UID bound to the authenticated, server-resolved topology.write grant; Datacenter/FaultDomain/InfrastructureStack from their resolved immutable parent reference; a client cannot supply or select metadata.scopeRef for any F0015 resource | ADH-047 decision 2 | VS0-SCHEMA-008..014 |
| REQ-F15-21 | Participation collection-create body is exactly metadata.name, spec.cloudPlatformRef (UID-pinned), spec.cloudProviderRef (UID-pinned), and spec.environment (development only); spec.providerSelectionModes and spec.permittedHostingLocationRefs are FEATURE-0021-introduced and FEATURE-0021-activated and are neither accepted, stored, defaulted, validated, nor exposed by F0015; item action routes are distinct from create and carry only an empty JSON body plus their required headers | ADH-047 decision 4 | VS0-SCHEMA-010 |
| REQ-F15-22 | An ISO-3166-1 alpha-2 value (CloudPlatform spec.ownerRegistration.jurisdictionCode, CloudProvider spec.operatingMarkets[], HostingLocation spec.countryCode) means an assigned code from a fixed, repository-owned, version-pinned static dataset (as of 2026-08-12); a syntactically valid but unassigned code (e.g. ZZ) is rejected with VALIDATION_FAILED (422); HostingLocation spec.administrativeAreaCode remains syntax-plus-country-prefix validation only, with no ISO-3166-2 membership dataset | ADH-048 decision 1 | VS0-SCHEMA-008,009,011 |
| REQ-F15-23 | Four exact malformed/prohibited-input outcomes, reusing only existing FEATURE-0012 top-level Problem codes: a missing/empty/malformed/over-length Idempotency-Key on a route that requires it returns MALFORMED_REQUEST (400); an If-Match header on a CloudProviderParticipation collection create returns MALFORMED_REQUEST (400); a missing/malformed/non-current If-Match on an existing-participation action returns STALE_RESOURCE_VERSION (412, preserved ADH-046 rule); a non-empty item-action body without valid duplicate-free top-level bootstrapGrant returns MALFORMED_REQUEST (400); none produce a mutation, idempotency record, or AuditEvent unless an already-approved audited-denial rule independently applies | ADH-048 decision 2,054 | VS0-SCHEMA-010 |
| REQ-F15-24 | Missing or invalid authentication on any F0015 owned method/path pattern returns AUTH_REQUIRED (401); no resource mutation, lifecycle/status write, idempotency record, or AuditEvent is produced | ADH-051 | VS0-SCHEMA-008..014 |

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
| CloudProviderParticipation `spec.providerSelectionModes` | Introduced and activated only by FEATURE-0021; rejected on F0015 participation create | VS0-SCHEMA-010 field ownership; ADH-047 decision 4 |
| CloudProviderParticipation `spec.permittedHostingLocationRefs` | Introduced and activated only by FEATURE-0021; rejected on F0015 participation create | VS0-SCHEMA-010 field ownership; ADH-047 decision 4 |
| Client-supplied `metadata.scopeRef` for any F0015 resource | Scope is always server-derived from exactly one registered source | ADH-047 decision 2 |
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
| AC-F15-15 | Every F0015 collection-create kind enforces its closed client-required/client-optional field boundary; server-owned/status/scope/unknown/deferred fields are rejected; each successful create returns 201 and initializes its exact status | VS0-CF-F15-26 |
| AC-F15-16 | Every F0015 resource kind derives its scope from exactly one server-side source (Platform-root, participation-from-cloudPlatformRef, HostingLocation-from-grant, descendant-from-parent); a client-supplied or mismatched scope is never honored | VS0-CF-F15-27 |
| AC-F15-17 | Participation collection-create body accepts exactly metadata.name/spec.cloudPlatformRef/spec.cloudProviderRef/spec.environment, initializes Pending/false holds/seven-day expiry, and rejects FEATURE-0021-owned fields; existing-item participation actions accept only an empty JSON body plus their required headers | VS0-CF-F15-28 |
| AC-F15-18 | A syntactically valid but unassigned ISO-3166-1 alpha-2 value is rejected with VALIDATION_FAILED against the fixed assigned-code dataset for jurisdictionCode, operatingMarkets, and countryCode; administrativeAreaCode validation remains syntax-plus-country-prefix only | VS0-CF-F15-29 |
| AC-F15-19 | A malformed/missing Idempotency-Key, an If-Match on participation create, and a non-empty participation item-action body each return MALFORMED_REQUEST with no mutation, idempotency record, or AuditEvent beyond an already-approved audited denial | VS0-CF-F15-30 |
| AC-F15-20 | Missing or invalid authentication on any F0015 owned method/path pattern returns AUTH_REQUIRED with no mutation, lifecycle/status write, idempotency record, or AuditEvent | VS0-CF-F15-31 |

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

### 4.5 Closed Collection-Create Request Contract (ADH-2026-047 decision 1)

For every F0015 collection `POST`, the client supplies only `metadata.name`, the required F0015-owned `spec` fields, and the registered optional F0015-owned `spec` fields for that kind:

| Kind | Client-required create fields | Client-optional create fields | Server-assigned outcome |
|---|---|---|---|
| CloudPlatform | `metadata.name`; `spec.ownerRegistration.legalName`; `spec.ownerRegistration.registrationIdentifier`; `spec.ownerRegistration.jurisdictionCode` | `metadata.displayName`; `spec.description` | deployment Platform-root `scopeRef`; `status.phase=Active`; `201` |
| CloudProvider | `metadata.name`; non-empty `spec.operatingMarkets[]` | `metadata.displayName`; `spec.displayName` | deployment Platform-root `scopeRef`; `status.phase=Active`; `201` |
| HostingLocation | `metadata.name`; `spec.countryCode`; `spec.locality` | `spec.administrativeAreaCode`; `spec.description` | CloudProvider `scopeRef` from §4.6; `status.phase=Active`; `201` |
| Datacenter | `metadata.name`; `spec.hostingLocationRef` | `spec.description` | CloudProvider `scopeRef` from resolved parent; `status.phase=Active`; `201` |
| FaultDomain | `metadata.name`; `spec.datacenterRef` | `spec.description` | CloudProvider `scopeRef` from resolved parent; `status.phase=Active`; `201` |
| InfrastructureStack | `metadata.name`; `spec.faultDomainRef` | `spec.description` | CloudProvider `scopeRef` from resolved parent; `status.phase=Active`; `201` |

The client must not supply `metadata.uid`, generation, resourceVersion, timestamps, `metadata.scopeRef`, any `status` field, a field owned by another feature, or an unknown field; such requests are rejected. All create routes require `Idempotency-Key`. A successful PATCH returns `200` with the updated resource.

### 4.6 Scope Derivation (ADH-2026-047 decision 2)

Each F0015 resource kind has exactly one server-side scope-derivation source; a client never supplies or selects `metadata.scopeRef`:

- **CloudPlatform, CloudProvider**: immutable deployment Platform-root scope, server-derived.
- **CloudProviderParticipation**: immutable CloudPlatform scope derived from resolved `spec.cloudPlatformRef`; the referenced CloudPlatform must exist and its UID must equal the resulting `metadata.scopeRef.uid`.
- **HostingLocation**: immutable CloudProvider scope derived from the CloudProvider UID bound to the authenticated, server-resolved `topology.write` grant. A caller cannot supply or choose `metadata.scopeRef`.
- **Datacenter, FaultDomain, InfrastructureStack**: immutable CloudProvider scope derived from their resolved immutable parent reference.

Reference resolution and authorization/non-disclosure occur before structural validation; derived-scope mismatches trigger the established safe-denial/validation outcome.

### 4.7 Topology Immutability Correction (ADH-2026-047 decision 3)

`metadata.name` is immutable identity for every F0015 resource, including topology resources (`HostingLocation`, `Datacenter`, `FaultDomain`, `InfrastructureStack`). For those topology resources, only `spec.description` is PATCHable. **Correction:** ADH-2026-045's topology sentence "Name and description are PATCHable" is corrected to "Description is PATCHable; name is immutable identity."

### 4.8 Participation Collection-Create Body (ADH-2026-047 decision 4)

`POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations` requires exactly `metadata.name`, `spec.cloudPlatformRef` (UID-pinned), `spec.cloudProviderRef` (UID-pinned), and `spec.environment` (`development` only). It accepts no optional participation `spec` fields; `spec.providerSelectionModes` and `spec.permittedHostingLocationRefs` are FEATURE-0021-owned and are rejected if present. The server derives `metadata.scopeRef` from `spec.cloudPlatformRef`, assigns `status.phase=Pending`, `status.platformSuspended=false`, `status.providerSuspended=false`, `status.requestExpiresAt=createdAt+7 days`, and returns `201`. This create requires `Idempotency-Key` only; `If-Match` is rejected/not accepted.

The eight public existing-participation item-action routes are:

- `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/accept`
- `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/reject`
- `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/withdraw`
- `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/suspend`
- `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/resume`
- `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/request-release`
- `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/accept-release`
- `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/decline-release`

In each path, `{uid}` is a complete Go 1.22 `http.ServeMux` wildcard segment, read through `r.PathValue("uid")`. Each item-action request has an empty JSON body and requires both `If-Match` and `Idempotency-Key`. Scheduler expiry is an internal API-server transition with no public request body.

### 4.9 ISO-3166 Assigned-Code Semantics and Malformed-Input Outcomes (ADH-2026-048)

For FEATURE-0015, an ISO-3166-1 alpha-2 value (`CloudPlatform.spec.ownerRegistration.jurisdictionCode`, `CloudProvider.spec.operatingMarkets[]`, `HostingLocation.spec.countryCode`) means an **assigned** code from a fixed, repository-owned, version-pinned static dataset (as of `2026-08-12`), not merely a syntactically valid two-letter code. A syntactically valid but unassigned code (e.g. `ZZ`) is rejected with `VALIDATION_FAILED` (422). `HostingLocation.spec.administrativeAreaCode` remains syntax-plus-country-prefix validation only (matching the existing ISO-3166-2 shape with the same two-letter prefix as `spec.countryCode`); F0015 introduces no ISO-3166-2 membership dataset.

Four exact header/body outcomes reuse only existing FEATURE-0012 top-level Problem codes; no new top-level code or violation code is introduced:

| Input | Exact outcome |
|---|---|
| Required `Idempotency-Key` missing, empty, malformed, or over the shared length limit | `MALFORMED_REQUEST` / 400 |
| `If-Match` supplied on a participation collection-create request | `MALFORMED_REQUEST` / 400 |
| `If-Match` missing, malformed, or non-current on an existing-participation action | `STALE_RESOURCE_VERSION` / 412 (preserved ADH-2026-046 rule) |
| Participation item action with any non-empty JSON body | `MALFORMED_REQUEST` / 400 |

A malformed/prohibited input under this rule produces no mutation, no idempotency record, and no AuditEvent unless an already-approved audited-denial rule independently applies; this clarification does not broaden the ADH-2026-046/047 audit-denial list.

FEATURE-0015 owns exactly 22 logical endpoint paths (seven collection paths, seven item paths, and eight CloudProviderParticipation action paths) and registers exactly 35 explicit Go 1.22 `http.ServeMux` method/path patterns: `GET` LIST and `POST` create for each of the seven collections (14); `GET` and `PATCH` for each of the six PATCHable resource items (12); `GET` only for the CloudProviderParticipation item, which has no PATCH surface (1); and `POST` for each of the eight participation actions (8) — 14+12+1+8 = 35. No path-only or catch-all handler registration, reflection, or internal HTTP-method dispatch is used; the approved complete `{uid}` wildcard segment is retained for item and action paths. Unsupported methods return the established method-not-allowed behavior and do not create a new endpoint path. (ADH-2026-051)

ADH-2026-055 closes executable behavior without adding a route or grant: HEAD returns 405 despite ServeMux GET matching; authenticated `X-Sovrunn-Bootstrap-Grant` denial precedes media/body validation; completed idempotency records retain for 24 hours and are capped at 10,000 per process with deterministic completed-only eviction; RFC 7396 PATCH uses staged post-merge validation; and successful LIST is ascending `metadata.uid`.

ADH-2026-056 closes remaining observable mechanics: every mutable-resource PATCH requires current `If-Match` and returns `STALE_RESOURCE_VERSION`/412 when missing, malformed, or stale; reads use `cloudplatform.read`, `cloudprovider.read`, `topology.read`, or `participation.read` at their exact derived scope with optional GET target-UID narrowing; and replay returns only its stored successful status/body plus `Content-Type: application/json`, with fresh request correlation on every request.

The exact FEATURE-0015-local proof cases for these closures are VS0-CF-F15-34 (HEAD), VS0-CF-F15-35 (forged header/media precedence), VS0-CF-F15-36 (retention/eviction/waiters), VS0-CF-F15-37 (PATCH/LIST), VS0-CF-F15-38 (PATCH conditional update), VS0-CF-F15-39 (read actions), and VS0-CF-F15-40 (replay response).

ADH-2026-057 adds VS0-CF-F15-41 as the sole exact local proof of wrong-administrator mutable-spec PATCH denial across the CloudPlatform and CloudProvider/topology writer boundaries. It preserves the distinct status, system-owned-field, immutable-field, and cross-provider safe-denial proofs.

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
| Caller without `participation.accept.platform` posts accept for Pending participation | AUTHORIZATION_DENIED | 403 | — |
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
| ISO-3166-1 alpha-2 value is syntactically valid but unassigned | VALIDATION_FAILED | 422 | — |
| Idempotency-Key missing, empty, malformed, or over-length on a route that requires it | MALFORMED_REQUEST | 400 | — |
| If-Match supplied on a CloudProviderParticipation collection create | MALFORMED_REQUEST | 400 | — |
| CloudProviderParticipation item action with a non-empty JSON body | MALFORMED_REQUEST | 400 | — |
| Required FEATURE-0013 AuditEvent append fails for a FEATURE-0015 mutation | INTERNAL_ERROR | 500 | — |
| Successful collection create | — (201) | 201 | — |
| Successful PATCH | — (200) | 200 | — |

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

---

## Feature-level reuse summary

This approved assessment uses the canonical format in
`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`. It records existing
authority and does not change the approved FEATURE-0015 architecture.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| FEATURE-0012 API/resource grammar and FEATURE-0013 AuditEvent contract | Reuse | The feature consumes the established metadata, Problem Details, validation-stage, and durable-audit contracts without redefining them. | Approved | DEC-0026; ADH-2026-045 |
| Canonical seven-resource cloud model and participation lifecycle | Build | No existing internal or external model supplies the approved scope, lifecycle, safe-denial, writer, and idempotency semantics as one contract. | Approved | DEC-0037; DEC-0054; DEC-0059; ADH-2026-045 |
| Provider-native execution and adapters | Wrap | Any provider-native translation remains behind a future adapter boundary and is not activated by FEATURE-0015. | Deferred | DEC-0036; ADH-2026-045 |

## Identity

| Field | Value |
|---|---|
| Feature identity | FEATURE-0015 |
| Capability or decision-unit identity | Canonical seven-resource cloud model and participation lifecycle |
| Assessment owner | Sanjeev Kumar, Sovrunn Architecture Owner |

## Classification

| Field | Value |
|---|---|
| Disposition | Build |
| Decision status | Approved |

## Analysis

| Field | Value |
|---|---|
| Assessment scope | Canonical CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, and InfrastructureStack contracts, their bounded HTTP surface, validation, lifecycle, and audit/idempotency behavior. |
| Candidate category | Existing Sovrunn contracts, cloud inventory domain models, and provider SDK/API models. |
| Mature candidates / applicable standards | FEATURE-0011 reuse governance; FEATURE-0012 API/resource grammar and RFC 9457 Problem Details; FEATURE-0013 AuditEvent; JSON Schema; RFC 7396 merge patch; ISO-3166 code standards. |
| Relevant candidate strengths | Inherited contracts provide stable metadata, validation, Problem Details, concurrency, audit, and conformance conventions. |
| Material candidate constraints | Existing alpha and provider-native models do not satisfy the approved seven-scope, safe-denial, participation-lifecycle, or no-runtime-migration boundary. |
| Rationale | Build only the F0015-owned canonical semantics and reuse inherited cross-cutting contracts exactly. |
| Selected foundation or approach | Seven typed canonical resources with FEATURE-0012 StageSet validation, FEATURE-0013 AuditEvent publication, and no external provider integration. |
| Why Reuse is insufficient | No reusable model provides the complete approved canonical resource ownership and lifecycle semantics. |
| Why Wrap is insufficient | Wrapping a provider-native model would introduce the deferred adapter/execution boundary. |
| Why Extend is insufficient | Extending retained alpha assets would violate DEC-0059's no-alpha-runtime-migration rule. |
| Protected Sovrunn differentiation and long-term ownership | Sovrunn owns canonical scope, resource, lifecycle, safe-denial, validation, and audit/idempotency semantics; prior features retain ownership of their common contracts. |

## Boundary

| Field | Value |
|---|---|
| Sovrunn-owned responsibility | The seven F0015 resource contracts, lifecycle transitions, request validation, authorization enforcement, and local conformance behavior. |
| Reused or extended responsibility | FEATURE-0011 reuse governance; FEATURE-0012 metadata, Problem Details, API validation pipeline, and resource conventions; FEATURE-0013 AuditEvent contract. |
| Responsibility/control boundary | F0015 composes inherited contracts but does not redefine them; later adapter, execution, IAM, eligibility, persistence, and provisioning responsibilities remain outside the feature. |
| Data crossing the boundary | Typed references, ObjectMeta, request context/grants, Problem Details, and redacted AuditEvent evidence. |
| Control crossing the boundary | Server-resolved authorization and local api-server transitions only; no provider calls, external persistence, or controller ownership. |
| Adapter required | No |
| Adapter rationale | The current feature creates no external integration; future provider translation is separately owned and must satisfy DEC-0036. |
| Adapter or contract identifier | none |
| Vendor-native types allowed | No |

## Suitability

| Field | Value |
|---|---|
| Sovereignty and deployment fit | Local canonical contracts and in-process phase boundary support disconnected and sovereign deployments without a cloud-vendor dependency. |
| Security and trust | Server-resolved grants, exact scope derivation, safe denial, writer enforcement, redacted audit evidence, and closed request contracts are mandatory. |
| Operational and supportability | Deterministic validation, registered Problem outcomes, conformance tests, and bounded idempotency retention make behavior diagnosable. |
| Licensing and supply-chain | Standard-library and repository-owned implementation; no new vendor SDK or external service dependency. |
| Portability and provider-neutrality impact | Canonical resources avoid provider-native runtime types and preserve a later replaceable adapter boundary. |

## Phase and scope

| Field | Value |
|---|---|
| Allowed in current phase | Yes |
| Current-phase work | Implement and verify the approved seven-resource canonical model, bounded API routes, validation, participation lifecycle, idempotency, and audit behavior. |
| Deferred work | Provider adapters, ExecutionTarget, external/durable persistence, real provisioning, customer eligibility, governance IAM, policy, and plugin execution. |
| Explicit non-goals | No alpha runtime migration, provider-native runtime objects, external provider calls, real provisioning, customer-facing selection, or future-feature semantics. |
| Exit or migration boundary | A change to canonical resource ownership, inherited contract ownership, or the no-migration boundary requires an approved ADH and compatibility review. |
| Phase 2 non-goal acknowledgement | Phase 2R authorizes no durable external persistence, real infrastructure provisioning, plugin execution, or future-feature activation. |

## Risk mitigation

| Field | Value |
|---|---|
| Applicable architecture risks | Contract drift, scope/authorization leakage, unsafe reference disclosure, idempotency replay errors, audit publication failure, and future-feature leakage. |
| Residual risk | Medium; residual risk is bounded by deterministic checks, formal models, conformance coverage, and final feature-gate evidence. |
| Replacement risk | Medium |
| Reassessment triggers | Any approved change to FEATURE-0012/0013 contracts, canonical resource ownership, route/lifecycle semantics, provider integration boundary, or a failing conformance/formal invariant. |

### Risk-control matrix

| Risk | Preventive control | Detection control | Corrective path |
|---|---|---|---|
| Contract drift | Registry/spec/traceability single-source checks | `make feature-contract-check FEATURE=FEATURE-0015` | Stop progression; correct the controlling authority through approved governance. |
| Scope or authorization leakage | Server-resolved grants, exact scope derivation, writer rules | Local conformance and authorization tests | Deny safely, repair the owning handler/validator, and rerun contract tests. |
| Reference disclosure | Safe-denial ordering and target-scope checks | F0015 safe-denial conformance cases | Correct ordering before release; do not expose existence. |
| Idempotency or audit publication failure | Target-bound idempotency and append-before-publication | TLC models, race tests, and F0015 conformance tests | Abort staged work, return the approved outcome, and investigate before retrying. |
| Future-feature leakage | Explicit exclusions and Phase 2R drift guardrails | `make phase2r-drift-check` | Remove the leaked behavior or raise an approved architecture change. |

## Traceability

| Field | Value |
|---|---|
| Related DEC / RFC / ADH references | DEC-0026; DEC-0036; DEC-0037; DEC-0041; DEC-0042; DEC-0054; DEC-0059; ADH-2026-045 through ADH-2026-057. |
| Linked acceptance criteria | AC-F15-01 through AC-F15-20 and the exact local VS0-CF-F15-01 through VS0-CF-F15-41 conformance cases. |
| Validation and review evidence | FEATURE-0015 architecture readiness, feature-contract check, TLC lifecycle/idempotency-audit models, VS-000 contract check, Phase 2R drift check, requirements/design/tasks approvals, and final feature gate. |

## Human-approval evidence

| Field | Value |
|---|---|
| Approving person or role | Sanjeev Kumar, Sovrunn Architecture Owner |
| Approval date | 2026-08-14 |
| Approved ADH or assessment-review reference | ADH-2026-045 |
| Structured approval-evidence record | docs/reviews/reuse-assessments/FEATURE-0015-approval-evidence.md |
| Approval applies to | FEATURE-0015 reuse-assessment governance and final feature-gate readiness; it does not alter architecture, Kiro stage, Cursor, or release approvals. |
