# FEATURE-0015 — Canonical Cloud Model Foundation: Requirements

## 1. Identity and stage

| Field | Value |
|-------|-------|
| Feature | FEATURE-0015 — Canonical Cloud Model Foundation |
| Spec | `.kiro/specs/canonical-cloud-model-and-alpha-migration/` |
| Stage | Requirements |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase | 2R |
| Owner scope | Requirements-owned observable behavior only (intent, actors, scenarios, invariants, validation outcomes, security/privacy, compatibility, non-goals, acceptance) |
| Architecture authority | `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`; `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md` |
| Controlling Decisions | DEC-0037, DEC-0041, DEC-0042, DEC-0054, DEC-0059 |
| Controlling Handoffs | ADH-2026-020, ADH-2026-024, ADH-2026-025, ADH-2026-037, ADH-2026-042, ADH-2026-043 (non-migration portions), ADH-2026-045, ADH-2026-046, ADH-2026-047, ADH-2026-048, ADH-2026-049, ADH-2026-050, ADH-2026-051, ADH-2026-052, ADH-2026-053, ADH-2026-054, ADH-2026-055, ADH-2026-056, ADH-2026-057 |
| Depends On | FEATURE-0011 (reuse governance), FEATURE-0012 (Problem Details/ObjectMeta/TypedRef/Condition/API resource standard), FEATURE-0013 (DecisionRecord/AuditEvent), FEATURE-0014 (retained repository asset only) |
| Depended On By | FEATURE-0016, FEATURE-0021, FEATURE-0022 |
| Owned resources | CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, InfrastructureStack |

This stage modifies only this file. It defines observable behavior; it does not own packages, files, routes, storage, algorithms, libraries, internal interfaces, or task decomposition. All normative content is a one-to-one translation of the approved architecture; no decision is reopened or extended.

---

## 2. Purpose and use cases

FEATURE-0015 establishes the seven-scope canonical cloud model identity layer through direct canonical bootstrap. It registers CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, and InfrastructureStack as managed resources created directly (not migrated), and it ends at InfrastructureStack. Per DEC-0059 (superseding DEC-0058) and ADH-2026-045, there is no runtime migration of the alpha model; FEATURE-0001 through FEATURE-0014 are retained repository assets, not live state requiring conversion.

Primary actors and use cases:

| Actor | Use case |
|-------|----------|
| Bootstrap principal (server-resolved) | Creates the first CloudPlatform (Platform-root) and CloudProvider resources under deterministic, server-resolved grants. |
| cloud-platform-admin | PATCHes CloudPlatform mutable fields; accepts/rejects/requests-release-of participations; sets/clears the platform suspension hold. |
| cloud-provider-admin | Creates topology (HostingLocation → Datacenter → FaultDomain → InfrastructureStack) and PATCHes its mutable fields; requests/withdraws participation; accepts/declines release; sets/clears the provider suspension hold. |
| Deterministic scheduler (system actor) | Invokes the api-server's participation expiry transition at `requestExpiresAt`. |
| api-server | Sole status writer for all FEATURE-0015-owned resources; derives scope; enforces validation, idempotency, and audit atomicity. |

Out-of-scope actors and consumer-facing visibility (customer/Organization eligibility, individual authorization) are governed by other features and are enumerated in §7.

---

## 3. Terms introduced by this feature

| Term | Meaning within FEATURE-0015 |
|------|------------------------------|
| CloudPlatform | Platform-scoped product-ownership resource with immutable `spec.ownerRegistration`. |
| CloudProvider | Platform-scoped infrastructure-operation resource with `spec.operatingMarkets[]`. |
| CloudProviderParticipation | CloudPlatform-scoped resource linking one CloudPlatform and one CloudProvider; carries the participation lifecycle and two independent suspension holds. |
| HostingLocation, Datacenter, FaultDomain, InfrastructureStack | CloudProvider-scoped topology chain, each with an immutable parent reference (topology ends at InfrastructureStack). |
| Independent suspension holds | `platformSuspended` and `providerSuspended`; each settable/clearable only by its respective administrator; effective Active requires both false. |
| Effective Suspended | Derived phase when either hold is true. |
| Deployment Platform-root scope | The server-derived Platform scope assigned to CloudPlatform and CloudProvider. |
| Assigned ISO-3166-1 alpha-2 value | A code present in the fixed, repository-owned, version-pinned static dataset (as of 2026-08-12); a syntactically valid but unassigned code is rejected. |
| Closed collection-create request contract | Per-kind boundary of exactly the client-required and client-optional fields; all other fields are rejected. |
| Bootstrap grant | Server-resolved, deterministic authorization; never accepted from header or body; no persisted roles/memberships/assignments. |

Contracts introduced by prior features (Problem Details, ObjectMeta, TypedRef, Condition, DecisionRecord, AuditEvent, the reuse governance contract) are consumed by exact reference only and are not redefined here.

---

## 4. Normative requirements and acceptance scenarios

Each approved requirement REQ-F15-01 through REQ-F15-24 has exactly one normative detail heading, presented in approved order. The exact requirement text is preserved verbatim in the `Canonical requirement ledger` (§4.25). Each approved acceptance criterion AC-F15-01 through AC-F15-20 is mapped explicitly in the acceptance-scenario coverage table (§4.26); the exact AC text is preserved verbatim in the `Canonical acceptance ledger` (§4.27).

### 4.1 REQ-F15-01 — CloudPlatform registration and identity

CloudPlatform is a Platform-scoped resource with a unique `metadata.name` and immutable `spec.ownerRegistration` (legalName, registrationIdentifier, jurisdictionCode only; no user/Organization/role/grant reference). The system validates uniqueness, scope, and immutable identity. Local proof: VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-12, VS0-CF-F15-13, VS0-CF-F15-14, VS0-CF-F15-32 (name uniqueness), VS0-CF-F15-33 (other schema constraints).

### 4.2 REQ-F15-02 — CloudProvider registration and identity

CloudProvider is a Platform-scoped resource with a unique `metadata.name` and a non-empty, unique, uppercase-ISO-3166-1-alpha2 `spec.operatingMarkets[]`. The system validates scope and identity. Local proof: VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-12, VS0-CF-F15-13, VS0-CF-F15-14, VS0-CF-F15-32 (name uniqueness), VS0-CF-F15-33 (operatingMarkets and other schema constraints), VS0-CF-F15-29 (assigned ISO-code semantics).

### 4.3 REQ-F15-03 — CloudProviderParticipation registration and lifecycle

CloudProviderParticipation is a CloudPlatform-scoped resource linking one CloudPlatform and one CloudProvider. Exactly one non-terminal participation may exist per platform+provider pair. Lifecycle: Pending→Active/Rejected/Withdrawn/Expired, Active↔Suspended, Active/Suspended→Terminating→Terminated (VS0-STATE-001). Local proof: VS0-CF-F15-01, VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-15, VS0-CF-F15-16 (scheduler expiry), VS0-CF-F15-20.

### 4.4 REQ-F15-04 — Topology chain under CloudProvider scope

The topology chain HostingLocation (uppercase ISO 3166-1 alpha-2 `countryCode`, non-empty `locality`, optional ISO 3166-2 `administrativeAreaCode`) → Datacenter → FaultDomain → InfrastructureStack is registered under CloudProvider scope; each descendant has an immutable parent TypedRef. Containment references and same-provider UID preservation through the chain are validated. Local proof: VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-13, VS0-CF-F15-14, VS0-CF-F15-23 (safe-denial versus authorized topology mismatch only), VS0-CF-F15-33 (registered topology field/reference validation).

### 4.5 REQ-F15-05 — Feature ends at InfrastructureStack

FEATURE-0015 introduces no ExecutionTarget identity, schema, route, status, writer, conformance, or target lifecycle behavior; ExecutionTarget belongs in its entirety to FEATURE-0016. Local proof: VS0-CF-F15-01, VS0-CF-F15-25.

### 4.6 REQ-F15-06 — Participation lifecycle actions and independent suspension

Participation exposes exactly eight public item-action routes (accept, reject, withdraw, suspend, resume, request-release, accept-release, decline-release); expiry is an internal api-server transition invoked by the deterministic scheduler, not a public participation action. Suspended is a derived effective state carried by two independent holds (platformSuspended, providerSuspended), each settable/clearable only by its respective administrator; clearing one hold never reactivates while the other remains true. Terminal states are Rejected, Withdrawn, Expired, Terminated. Local proof: VS0-CF-F15-08, VS0-CF-F15-10, VS0-CF-F15-15, VS0-CF-F15-16 (scheduler expiry), VS0-CF-F15-20.

### 4.7 REQ-F15-07 — No mutable participation spec; header preconditions

FEATURE-0015 exposes no mutable CloudProviderParticipation spec fields. Create (participation.request) requires `Idempotency-Key` only and no `If-Match`. Every action on an existing participation (accept, reject, withdraw, suspend, resume, request-release, accept-release, decline-release) requires both `If-Match` and `Idempotency-Key`. Local proof: VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-15, VS0-CF-F15-17, VS0-CF-F15-18, VS0-CF-F15-19, VS0-CF-F15-20.

### 4.8 REQ-F15-08 — PATCH-only update surface

Mutable CloudPlatform, CloudProvider, and topology updates use PATCH with `application/merge-patch+json` only; PUT and DELETE are not exposed and return 405. PATCH requires `If-Match` carrying the current `metadata.resourceVersion`; a missing, malformed, or stale value returns STALE_RESOURCE_VERSION/412 with no mutation, AuditEvent, or idempotency record. Local proof: VS0-CF-F15-07, VS0-CF-F15-13, VS0-CF-F15-14, VS0-CF-F15-38.

### 4.9 REQ-F15-09 — Cross-provider isolation

A cross-provider target reference is denied with safe RESOURCE_NOT_FOUND and no existence disclosure. Local proof: VS0-CF-X03, VS0-CF-F15-23.

### 4.10 REQ-F15-10 — Seven canonical scope kinds

Only the seven canonical scope kinds are used globally. CloudPlatform and CloudProvider accept only Platform scope; CloudProviderParticipation only CloudPlatform scope; topology resources only CloudProvider scope. Local proof: VS0-CF-F15-01, VS0-CF-F15-02.

### 4.11 REQ-F15-11 — Writer enforcement

Only cloud-provider-admin may write CloudProvider or topology spec; only cloud-platform-admin may write CloudPlatform spec (consumed by reference to VS0-WRITER-002 and VS0-WRITER-003). Local proof: VS0-CF-F15-41.

### 4.12 REQ-F15-12 — Existing Problem Details codes only

All Problem Details errors use existing FEATURE-0012 codes; no new top-level codes are introduced. Pure HTTP transport outcomes, including HTTP 405, are not top-level Problem codes. This is an anti-drift/non-runtime obligation proven collectively by every local case that emits Problem Details; it is not a standalone runtime-observable scenario.

### 4.13 REQ-F15-13 — No migration plan/record/controller/cutover

No CanonicalMigrationPlan, CanonicalMigrationRecord, migration controller, or cutover state machine is implemented or reintroduced; FEATURE-0001–0014 are retained repository assets. This is an explicit anti-drift/non-runtime verification gate proven by absence (retired tombstones VS0-SCHEMA-060/061, VS0-WRITER-020/021, VS0-STATE-011 never reused), excluded from the runtime proof-case mapping per ADH-2026-046 decision 3.

### 4.14 REQ-F15-14 — Scope reference UID invariant

`CloudProviderParticipation.spec.cloudPlatformRef.uid` must equal its CloudPlatform `metadata.scopeRef.uid`. CloudPlatform itself has no Organization reference to validate (immutable `spec.ownerRegistration` replaces `ownerOrganizationRef`). Mismatch returns VALIDATION_FAILED (422) with VS0_SCOPE_REFERENCE_MISMATCH. Local proof: VS0-CF-F15-11.

### 4.15 REQ-F15-15 — Audit evidence, atomicity, and append-failure outcome

After authentication and current authorization are evaluated, a redacted FEATURE-0013 AuditEvent (with request correlation, no secret or inaccessible-resource disclosure) is produced exactly for: successful resource collection create or PATCH; successful participation create/action or deterministic scheduler expiry; an authenticated client attempt to write api-server-owned status or identity/metadata; a CloudPlatform-root creation denial; an authenticated collection LIST without a read grant; and an authenticated inaccessible cross-provider/safe-denial reference. The exact forged-grant paths are the exception: REQ-F15-17 defines their ADH-2026-054 precedence, under which the reserved header or valid duplicate-free reserved body member is denied and audited before current authorization. No durable AuditEvent is produced for missing/invalid authentication, malformed/prohibited body/header/field input other than that exact forged-grant carrier, unsupported method/media type, stale If-Match, same-key replay, changed-digest idempotency conflict, pair-uniqueness conflict, or an invalid participation source state. The mutation and its required AuditEvent are one atomic outcome; a failure to append a required AuditEvent returns the inherited INTERNAL_ERROR (500) Problem Details response with safe non-disclosure, and the resource mutation, lifecycle/status transition, and idempotency completion are not published; DEPENDENCY_UNAVAILABLE is not used for this outcome. Local proof: VS0-CF-F15-24.

### 4.16 REQ-F15-16 — CloudPlatform-root requirement

A CloudProvider, HostingLocation, Datacenter, FaultDomain, InfrastructureStack, or CloudProviderParticipation create is denied until at least one CloudPlatform exists, returning CONFLICT (409) with VS0_CLOUDPLATFORM_ROOT_REQUIRED. Local proof: VS0-CF-F15-12.

### 4.17 REQ-F15-17 — Bootstrap authorization boundary

Grants are server-resolved and deterministic only; grants are never accepted from an HTTP header or request body. Only `X-Sovrunn-Bootstrap-Grant` and valid duplicate-free top-level `bootstrapGrant` are audited forged-grant carriers: authentication precedes header denial; malformed/oversized/duplicate bodies retain existing 400 outcomes; the valid body carrier is denied before authorization/reference access; all other unknown body members remain 400. FEATURE-0015 persists no roles, memberships, or assignments. Local proof: VS0-CF-F15-22.

### 4.18 REQ-F15-18 — Idempotency semantics

POST create/actions require `Idempotency-Key`. An item action binds its namespace to the concrete target UID; current authentication, authorization, and safe target/reference access precede any completed replay. The same principal+route+target (when applicable)+key+digest returns the original result before current-version/lifecycle validation; a different digest with the same key returns CONFLICT (409) with VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH. A replay returns only the stored successful status/body and `Content-Type: application/json`; each request receives fresh request correlation. Only successful collection creates and successful participation actions complete a replay record; validation, uniqueness, lifecycle, stale-version, audit failure, panic, and shutdown abort it. Local proof: VS0-CF-F15-18, VS0-CF-F15-19, VS0-CF-F15-40.

### 4.19 REQ-F15-19 — Closed collection-create request contract

For every FEATURE-0015 collection POST, the client supplies only `metadata.name`, the required kind-specific spec fields, and the registered optional kind-specific spec fields. The server rejects `metadata.uid`/generation/resourceVersion/timestamps, `metadata.scopeRef`, any status field, a field owned by another feature, or an unknown field. Every successful create returns 201 with its exact initial status; a successful PATCH returns 200 with the updated resource. Local proof: VS0-CF-F15-26.

### 4.20 REQ-F15-20 — Single-source scope derivation

Each FEATURE-0015 resource kind has exactly one server-side scope-derivation source: CloudPlatform and CloudProvider from the deployment Platform-root; CloudProviderParticipation from its resolved `spec.cloudPlatformRef` (UID-equality enforced); HostingLocation from the CloudProvider UID bound to the authenticated, server-resolved `topology.write` grant; Datacenter/FaultDomain/InfrastructureStack from their resolved immutable parent reference. A client cannot supply or select `metadata.scopeRef` for any FEATURE-0015 resource. Local proof: VS0-CF-F15-27.

### 4.21 REQ-F15-21 — Participation collection-create body and item-action body

The participation collection-create body is exactly `metadata.name`, `spec.cloudPlatformRef` (UID-pinned), `spec.cloudProviderRef` (UID-pinned), and `spec.environment` (`development` only). `spec.providerSelectionModes` and `spec.permittedHostingLocationRefs` are FEATURE-0021-introduced and FEATURE-0021-activated and are neither accepted, stored, defaulted, validated, nor exposed by FEATURE-0015. Item action routes are distinct from create and carry only an empty JSON body plus their required headers. Local proof: VS0-CF-F15-28.

### 4.22 REQ-F15-22 — ISO-3166-1 alpha-2 assigned-code semantics

An ISO-3166-1 alpha-2 value (CloudPlatform `spec.ownerRegistration.jurisdictionCode`, CloudProvider `spec.operatingMarkets[]`, HostingLocation `spec.countryCode`) means an assigned code from a fixed, repository-owned, version-pinned static dataset (as of 2026-08-12). A syntactically valid but unassigned code (e.g. ZZ) is rejected with VALIDATION_FAILED (422). HostingLocation `spec.administrativeAreaCode` remains syntax-plus-country-prefix validation only, with no ISO-3166-2 membership dataset. Local proof: VS0-CF-F15-29.

### 4.23 REQ-F15-23 — Malformed/prohibited-input outcomes

Reusing only existing FEATURE-0012 top-level Problem codes: a missing/empty/malformed/over-length `Idempotency-Key` on a route that requires it returns MALFORMED_REQUEST (400); an `If-Match` header on a participation collection create returns MALFORMED_REQUEST (400); a missing/malformed/non-current `If-Match` on an existing-participation action returns STALE_RESOURCE_VERSION (412); a non-empty item-action body without valid duplicate-free top-level `bootstrapGrant` returns MALFORMED_REQUEST (400). None of these produce a mutation, idempotency record, or AuditEvent unless an already-approved audited-denial rule independently applies. Local proof: VS0-CF-F15-30.

### 4.24 REQ-F15-24 — Authentication on every owned method/path pattern

Missing or invalid authentication on any FEATURE-0015 owned method/path pattern returns AUTH_REQUIRED (401); no resource mutation, lifecycle/status write, idempotency record, or AuditEvent is produced. Local proof: VS0-CF-F15-31.

### 4.25 Canonical requirement ledger

The following table is an exact, one-to-one copy of the approved REQ table (`docs/features/FEATURE-0015-canonical-cloud-model-foundation.md` §2.1). Columns and cell text are preserved verbatim; no row is renumbered, merged, split, omitted, or paraphrased.

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

### 4.26 Acceptance-scenario coverage mapping

Every approved AC is mapped here to its detailed acceptance scenario and its FEATURE-0015-local conformance case(s), outside the canonical acceptance ledger. Anti-drift/non-runtime items are labelled per ADH-2026-046 decision 3 and carry no runtime proof case.

| AC ID | Acceptance scenario (observable outcome) | FEATURE-0015-local proof case(s) |
|-------|-------------------------------------------|----------------------------------|
| AC-F15-01 | CloudPlatform, CloudProvider, and topology-chain collection create, GET/LIST, and permitted PATCH behavior succeed with validation at their registry-declared scope subsets; PUT and DELETE return 405; the feature ends at InfrastructureStack. | VS0-CF-F15-01, VS0-CF-F15-02, VS0-CF-F15-38, VS0-CF-F15-39 |
| AC-F15-02 | Participation lifecycle with provider request plus CloudPlatform acceptance guard is enforced; pair uniqueness is enforced. | VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-09 |
| AC-F15-03 | No ExecutionTarget identity/schema/route/status/writer/conformance/target-lifecycle behavior is introduced. | VS0-CF-F15-01, VS0-CF-F15-25 |
| AC-F15-04 | A cross-provider target reference is denied safely (no existence disclosure). | VS0-CF-X03, VS0-CF-F15-23 |
| AC-F15-05 | No migration plan/record/controller/cutover state machine exists. **Anti-drift/non-runtime verification gate** (excluded per ADH-2026-046 decision 3). | — |
| AC-F15-06 | PUT/DELETE return 405 for CloudPlatform, CloudProvider, and topology; PATCH accepts only `application/merge-patch+json`. | VS0-CF-F15-07 |
| AC-F15-07 | Every output passes the canonical stale-concept anti-drift gate. **Anti-drift/non-runtime verification gate** (excluded per ADH-2026-046 decision 3). | — |
| AC-F15-08 | Unauthorized writes to CloudPlatform, CloudProvider, or topology specs, status, and system-owned fields are denied; cross-provider references deny safely. | VS0-CF-F15-05, VS0-CF-F15-06, VS0-CF-F15-41, VS0-CF-X03 |
| AC-F15-09 | All Problem Details errors use FEATURE-0012 codes from its closed set. **Contract-reuse claim** proven by the closed registry set: each non-null `expectedError` in §4.28 is an existing FEATURE-0012 code; pure transport outcomes recorded in `expectedSideEffects`, including HTTP 405, introduce no top-level Problem code. Not a standalone runtime scenario. | Closed-set proof: §4.28 exact conformance semantics ledger |
| AC-F15-10 | Independent holds behave correctly; suspend/resume denied in terminal/Pending/Terminating phases; clearing one hold does not reactivate while the other remains true. | VS0-CF-F15-08, VS0-CF-F15-10 |
| AC-F15-11 | Participation exposes no mutable spec fields; create requires Idempotency-Key only; existing-participation lifecycle changes require If-Match and Idempotency-Key. | VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-15, VS0-CF-F15-17, VS0-CF-F15-18, VS0-CF-F15-19 |
| AC-F15-12 | Scope-reference UID invariant enforced; mismatch rejected with VS0_SCOPE_REFERENCE_MISMATCH. | VS0-CF-F15-11 |
| AC-F15-13 | Lifecycle/hold-changing actions, resource create/PATCH, and security denials produce correlated, redacted AuditEvent with no secrets. | VS0-CF-F15-24 |
| AC-F15-14 | CloudPlatform-root requirement enforced with CONFLICT/409 VS0_CLOUDPLATFORM_ROOT_REQUIRED until at least one CloudPlatform exists. | VS0-CF-F15-12 |
| AC-F15-15 | Every collection-create kind enforces its closed field boundary; server-owned/status/scope/unknown/deferred fields rejected; each successful create returns 201 with exact initial status. | VS0-CF-F15-26 |
| AC-F15-16 | Every resource kind derives scope from exactly one server-side source; a client-supplied or mismatched scope is never honored. | VS0-CF-F15-27 |
| AC-F15-17 | Participation create accepts exactly the four decision-4 fields, initializes Pending/false holds/seven-day expiry, and rejects FEATURE-0021-owned fields; item actions accept only an empty JSON body plus required headers. | VS0-CF-F15-28 |
| AC-F15-18 | A syntactically valid but unassigned ISO-3166-1 alpha-2 value is rejected with VALIDATION_FAILED; administrativeAreaCode remains syntax-plus-country-prefix only. | VS0-CF-F15-29 |
| AC-F15-19 | A malformed/missing Idempotency-Key, an If-Match on participation create, and an ordinary non-empty participation item-action body that lacks the reserved top-level `bootstrapGrant` member in a syntactically valid duplicate-free body each return MALFORMED_REQUEST with no mutation/idempotency record/AuditEvent. The narrow ADH-2026-054 exception—a valid duplicate-free body containing that reserved member—returns audited AUTHORIZATION_DENIED instead. | VS0-CF-F15-30; VS0-CF-F15-22 |
| AC-F15-20 | Missing or invalid authentication on any owned method/path pattern returns AUTH_REQUIRED with no mutation, lifecycle/status write, idempotency record, or AuditEvent. | VS0-CF-F15-31 |
| Idempotency replay response | A completed same-key/same-digest replay returns the stored successful status/body and `Content-Type: application/json` only, with fresh request correlation. | VS0-CF-F15-40 |

### 4.27 Canonical acceptance ledger

The following table is an exact, one-to-one copy of the approved acceptance-criteria table (`docs/features/FEATURE-0015-canonical-cloud-model-foundation.md` §3). Columns and cell text are preserved verbatim; no row is renumbered, merged, split, omitted, or paraphrased.

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

### 4.28 Exact conformance semantics ledger

For every referenced `VS0-CF-*` ID, the registry fields ID, owner, inputs, expectedState, expectedError, expectedViolation (where registered), expectedSideEffects, and gate are copied exactly from `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`. A conformance case proves only its registered scenario; a conformance case is never used as a substitute for a state-machine ID, and a downstream case is never reinterpreted as this feature's behavior.

| ID | owner | inputs | expectedState | expectedError | expectedViolation | expectedSideEffects | gate |
|----|-------|--------|---------------|---------------|-------------------|---------------------|------|
| VS0-CF-X03 | FEATURE-0015 | cross-provider-target-ref | unchanged | RESOURCE_NOT_FOUND | VS0_AUTHORIZATION_SAFE_DENIAL | exactly one redacted FEATURE-0013 AuditEvent with request correlation and no target-existence disclosure; a required-append failure returns INTERNAL_ERROR (500) with no publication (ADH-2026-051) | security |
| VS0-CF-F15-01 | FEATURE-0015 | authorized creation for each FEATURE-0015 kind at its registry-declared scope subset | resource persisted with uid-pinned references; FEATURE-0015 ends at InfrastructureStack | null | — | no ExecutionTarget contract, invariant, or route introduced by FEATURE-0015 | feature |
| VS0-CF-F15-02 | FEATURE-0015 | FEATURE-0015 resource creation with a scope kind outside that kind's registry-declared subset | unchanged | VALIDATION_FAILED | VS0_SCOPE_KIND_INVALID | none | feature |
| VS0-CF-F15-03 | FEATURE-0015 | Pending CloudProviderParticipation created by the provider's request, with valid refs, followed by the CloudPlatform's accept action | participation Active | null | — | controller-owned phase only | feature |
| VS0-CF-F15-04 | FEATURE-0015 | second non-terminal CloudProviderParticipation for an existing CloudPlatform and CloudProvider pair | unchanged | ALREADY_EXISTS | VS0_PARTICIPATION_DUPLICATE | no second participation persisted | feature |
| VS0-CF-F15-05 | FEATURE-0015 | authenticated client write to a FEATURE-0015 managed-resource status field | unchanged | AUTHORIZATION_DENIED | VS0_STATUS_FIELD_WRITE | exactly one redacted FEATURE-0013 AuditEvent with request correlation; a required-append failure returns INTERNAL_ERROR (500) with no publication (ADH-2026-051) | security |
| VS0-CF-F15-06 | FEATURE-0015 | authenticated client write to api-server-owned identity or metadata field | unchanged | AUTHORIZATION_DENIED | VS0_SYSTEM_OWNED_FIELD_WRITE | exactly one redacted FEATURE-0013 AuditEvent with request correlation; a required-append failure returns INTERNAL_ERROR (500) with no publication (ADH-2026-051) | security |
| VS0-CF-F15-07 | FEATURE-0015 | PUT or DELETE request against CloudPlatform, CloudProvider, or a topology resource | unchanged | null | — | 405 returned; only PATCH with application/merge-patch+json is accepted for mutable resources | feature |
| VS0-CF-F15-08 | FEATURE-0015 | resume action clearing only the acting party's own hold while the other hold remains true | participation remains effectively Suspended | null | — | only the acting administrator's own hold changes; the other hold is untouched | feature |
| VS0-CF-F15-09 | FEATURE-0015 | caller without participation.accept.platform posts the existing accept action for a Pending CloudProviderParticipation | participation remains Pending | AUTHORIZATION_DENIED | — | no Active transition, lifecycle/status write, AuditEvent, or idempotency completion | feature |
| VS0-CF-F15-10 | FEATURE-0015 | suspend or resume action attempted while participation phase is Pending, Terminating, Rejected, Withdrawn, Expired, or Terminated | unchanged | CONFLICT | VS0_PARTICIPATION_STATE_INVALID | no hold change | feature |
| VS0-CF-F15-11 | FEATURE-0015 | CloudProviderParticipation with spec.cloudPlatformRef.uid != metadata.scopeRef.uid | unchanged | VALIDATION_FAILED | VS0_SCOPE_REFERENCE_MISMATCH | none | feature |
| VS0-CF-F15-12 | FEATURE-0015 | CloudProvider, HostingLocation, Datacenter, FaultDomain, InfrastructureStack, or CloudProviderParticipation create attempted before any CloudPlatform exists | unchanged | CONFLICT | VS0_CLOUDPLATFORM_ROOT_REQUIRED | no resource persisted; exactly one redacted FEATURE-0013 AuditEvent with request correlation for the denial; a required-append failure returns INTERNAL_ERROR (500) with no publication (ADH-2026-051) | feature |
| VS0-CF-F15-13 | FEATURE-0015 | valid PATCH request using application/merge-patch+json against an approved mutable field of CloudPlatform, CloudProvider, or a topology resource | only the approved mutable field changes; resourceVersion increments | null | — | no other field changes | feature |
| VS0-CF-F15-14 | FEATURE-0015 | PATCH request using an unsupported media type, or a PATCH attempting to change an immutable, status, scope, identity, relationship-reference, or owner-registration field of CloudPlatform, CloudProvider, or a topology resource | unchanged | null | — | unsupported media type returns 415 UNSUPPORTED_MEDIA_TYPE; immutable/scope/identity/relationship-reference/owner-registration field changes return VALIDATION_FAILED (422) VS0_PATCH_IMMUTABLE_FIELD; status field changes return AUTHORIZATION_DENIED (403) VS0_STATUS_FIELD_WRITE; the resource is left unchanged in every case | feature |
| VS0-CF-F15-15 | FEATURE-0015 | CloudProviderParticipation created by a CloudProvider administrator's request action with valid refs | status.phase=Pending; status.platformSuspended=false; status.providerSuspended=false; status.requestExpiresAt=createdAt + 7 days | null | — | no Active relationship permitted before the CloudPlatform's accept action | feature |
| VS0-CF-F15-16 | FEATURE-0015 | deterministic system scheduler invoking the api-server's expiry transition for a still-Pending participation at requestExpiresAt | participation Expired | null | — | idempotent; exactly one correlated AuditEvent naming the system actor and scheduler execution ID | feature |
| VS0-CF-F15-17 | FEATURE-0015 | missing, malformed, or non-current If-Match on an action against an existing CloudProviderParticipation | unchanged | STALE_RESOURCE_VERSION | — | no lifecycle/status write; no additional AuditEvent | feature |
| VS0-CF-F15-18 | FEATURE-0015 | same authenticated principal, registered method/path pattern, Idempotency-Key, and canonical request digest replayed against any of the seven FEATURE-0015 collection creates (CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, InfrastructureStack), or against one concrete CloudProviderParticipation UID on any of the eight existing participation item-action routes (accept, reject, withdraw, suspend, resume, request-release, accept-release, decline-release), after current server-resolved authorization and safe target/reference access succeed | original result returned before current-version and lifecycle-source-state validation | null | — | no duplicate resource, transition, or AuditEvent; no completed replay for another concrete item target; a revoked, narrowed, or otherwise insufficient current grant denies without stored-result disclosure | idempotency |
| VS0-CF-F15-19 | FEATURE-0015 | same authenticated principal, registered method/path pattern, Idempotency-Key, and concrete CloudProviderParticipation UID when applicable, with a changed canonical request digest against any of the seven FEATURE-0015 collection creates (CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, InfrastructureStack) or any of the eight existing participation item-action routes (accept, reject, withdraw, suspend, resume, request-release, accept-release, decline-release), after current server-resolved authorization and safe target/reference access succeed | original result retained | CONFLICT | VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH | no stored-result disclosure when current authorization or safe access fails | idempotency |
| VS0-CF-F15-20 | FEATURE-0015 | reject, withdraw, request-release, accept-release, and decline-release actions invoked by their exact authorized action grant against a participation in the state required for each transition, and the same actions attempted from an invalid source state | each authorized action reaches its registered target state; decline-release resolves to the correct hold-derived Active or Suspended state; an action attempted from an invalid source state leaves the participation unchanged | CONFLICT | VS0_PARTICIPATION_STATE_INVALID | no state change for an invalid-source-state attempt | feature |
| VS0-CF-F15-21 | FEATURE-0015 | authenticated collection LIST without a read grant, LIST with a scoped read grant covering a subset of resources, LIST with a scoped read grant covering no visible resources, and direct GET or reference resolution without access | unchanged | AUTHORIZATION_DENIED | — | collection LIST without a read grant returns 403 without resource detail and exactly one redacted FEATURE-0013 AuditEvent with request correlation; a required-append failure returns INTERNAL_ERROR (500) with no publication (ADH-2026-051); LIST with a scoped read grant returns only visible resources; LIST with a scoped read grant and no visible resources returns 200 with an empty collection; unauthorized direct GET or reference resolution returns 404 RESOURCE_NOT_FOUND/VS0_AUTHORIZATION_SAFE_DENIAL | security |
| VS0-CF-F15-22 | FEATURE-0015 | an authenticated request presenting X-Sovrunn-Bootstrap-Grant with any value, or presenting a top-level bootstrapGrant JSON member with any value in a syntactically valid duplicate-free body, attempting to expand or substitute for the server-resolved bootstrap grant | unchanged | AUTHORIZATION_DENIED | — | only the exact reserved header/member receive this outcome; header detection follows authentication and precedes body decode; malformed/oversized/duplicate body checks precede body-member detection; body-member detection precedes current authorization, root/reference access, idempotency, and normal unknown-field classification; both claims are ignored and undisclosed; only the server-resolved action/scope/resourceUID-restricted grant authorizes its intended F0015 operation; exactly one redacted FEATURE-0013 AuditEvent with request correlation; a required-append failure returns INTERNAL_ERROR (500) with no publication (ADH-2026-054) | security |
| VS0-CF-F15-23 | FEATURE-0015 | an authenticated inaccessible cross-provider topology reference, and a caller authorized for both resources presenting a structurally mismatched provider graph | unchanged | RESOURCE_NOT_FOUND for the inaccessible cross-provider reference; VALIDATION_FAILED for the authorized structurally mismatched provider graph | VS0_AUTHORIZATION_SAFE_DENIAL for the inaccessible reference; VS0_TOPOLOGY_PROVIDER_MISMATCH for the authorized mismatch | exactly one redacted FEATURE-0013 AuditEvent with request correlation and no target-existence disclosure for the inaccessible-reference case; a required-append failure returns INTERNAL_ERROR (500) with no publication (ADH-2026-051) | security |
| VS0-CF-F15-24 | FEATURE-0015 | successful resource collection create/PATCH, successful participation create/action/expiry, authenticated authorization denial (status/identity-metadata write, root-creation denial, LIST-without-read-grant, forged-grant attempt), and authenticated safe denial, including a simulated required-AuditEvent-append failure during one such outcome | the outcome and its required AuditEvent are one atomic publication; on append failure neither the resource mutation, lifecycle/status transition, nor idempotency completion is published | INTERNAL_ERROR | — | exactly one redacted FEATURE-0013 AuditEvent with prescribed correlation per event; a required-AuditEvent-append failure leaves the associated outcome unpublished and returns INTERNAL_ERROR (500) with safe non-disclosure; DEPENDENCY_UNAVAILABLE is not used for this outcome; an already-appended AuditEvent may remain after a later in-process publication failure (ADH-2026-049, extended by ADH-2026-051) | observability |
| VS0-CF-F15-25 | FEATURE-0015 | the complete F0015 owned route/method surface, participation action request bodies, and a scan for ExecutionTarget or retired migration routes/schemas/writers/conformance | owned collections/items expose only POST/GET/LIST/PATCH; participation action requests carry an empty JSON body | null | — | no F0016 ExecutionTarget route/schema/writer/conformance and no retired migration route/schema/writer/conformance is exposed by FEATURE-0015 | feature |
| VS0-CF-F15-26 | FEATURE-0015 | per-kind collection-create request boundary | Per-kind collection-create request boundary: only the approved client fields are accepted; server-owned/status/scope/unknown/deferred fields are rejected; each successful create returns 201, initializes its exact status, and has one required AuditEvent. | null | — | Per-kind collection-create request boundary: only the approved client fields are accepted; server-owned/status/scope/unknown/deferred fields are rejected; each successful create returns 201, initializes its exact status, and has one required AuditEvent. | feature |
| VS0-CF-F15-27 | FEATURE-0015 | scope derivation | Scope derivation: Platform-root, participation-from-CloudPlatform, HostingLocation-from-grant, and descendant-from-parent scopes are server derived; client-supplied/mismatched scope never selects another authority. | null | — | Scope derivation: Platform-root, participation-from-CloudPlatform, HostingLocation-from-grant, and descendant-from-parent scopes are server derived; client-supplied/mismatched scope never selects another authority. | feature |
| VS0-CF-F15-28 | FEATURE-0015 | participation collection-create body versus empty item-action body | Participation collection-create body versus empty item-action body: create accepts exactly the decision-4 fields, initializes Pending/false holds/seven-day expiry, rejects F0021 fields, and existing item actions accept only empty JSON body plus their required headers. | null | — | Participation collection-create body versus empty item-action body: create accepts exactly the decision-4 fields, initializes Pending/false holds/seven-day expiry, rejects F0021 fields, and existing item actions accept only empty JSON body plus their required headers. | feature |
| VS0-CF-F15-29 | FEATURE-0015 | a syntactically valid but unassigned ISO-3166-1 alpha-2 value (e.g. ZZ) supplied for CloudPlatform spec.ownerRegistration.jurisdictionCode, CloudProvider spec.operatingMarkets[], or HostingLocation spec.countryCode, checked against the fixed repository-owned version-pinned assigned-code dataset (as of 2026-08-12); and a syntactically valid HostingLocation spec.administrativeAreaCode with the correct two-letter country prefix (no ISO-3166-2 membership dataset) | unchanged | VALIDATION_FAILED | — | none | feature |
| VS0-CF-F15-30 | FEATURE-0015 | a missing, empty, malformed, or over-length Idempotency-Key on a route that requires it; an If-Match header supplied on a CloudProviderParticipation collection-create request; and a non-empty CloudProviderParticipation item-action body that does not contain the reserved top-level bootstrapGrant member in a syntactically valid duplicate-free body | unchanged | MALFORMED_REQUEST | — | no mutation, no idempotency record, and no AuditEvent unless an already-approved audited-denial rule independently applies; a valid duplicate-free body with bootstrapGrant is governed instead by VS0-CF-F15-22 (ADH-2026-054) | feature |
| VS0-CF-F15-31 | FEATURE-0015 | missing or invalid authentication on any F0015 owned method/path pattern | unchanged | AUTH_REQUIRED | — | no mutation, no lifecycle/status write, no idempotency record, no AuditEvent | security |
| VS0-CF-F15-32 | FEATURE-0015 | an authorized CloudPlatform or CloudProvider collection create whose metadata.name duplicates an existing resource of the same kind in that kind's server-derived Platform-root scope | unchanged | ALREADY_EXISTS | — | no duplicate resource persisted | feature |
| VS0-CF-F15-33 | FEATURE-0015 | an authorized FEATURE-0015 collection create that, after closed-contract classification, violates a registry-declared semantic schema constraint in VS0-SCHEMA-008..014 other than a duplicate name, an unassigned ISO-3166-1 alpha-2 code covered by VS0-CF-F15-29, or a malformed/prohibited header/body condition covered by VS0-CF-F15-30 | unchanged | VALIDATION_FAILED | — | no resource persisted, no idempotency completion, no AuditEvent unless an already-approved audited-denial rule independently applies | feature |
| VS0-CF-F15-34 | FEATURE-0015 | a HEAD request matched by Go 1.22 http.ServeMux to any registered FEATURE-0015 GET pattern | unchanged | null | — | the handler checks the actual request method and returns HTTP 405; this transport outcome introduces no top-level Problem code; no AuditEvent, mutation, lifecycle/status write, or idempotency record; exactly 22 logical paths and 35 explicit method/path registrations remain unchanged | feature |
| VS0-CF-F15-35 | FEATURE-0015 | an authenticated request with a present X-Sovrunn-Bootstrap-Grant header and an otherwise unsupported PATCH media type | unchanged | AUTHORIZATION_DENIED | — | header denial precedes media and body validation; exactly one redacted FEATURE-0013 AuditEvent with request correlation; no mutation, lifecycle/status write, or idempotency record; without the reserved header, unsupported PATCH media type remains UNSUPPORTED_MEDIA_TYPE before body decode | security |
| VS0-CF-F15-36 | FEATURE-0015 | a completed FEATURE-0015 idempotency record at 24 hours after completion or selected for deterministic completed-record capacity eviction, and a cancelled waiter on an InFlight reservation | expired or evicted completed record is unavailable for replay and the next request is processed as new; cancellation detaches only its waiter; InFlight record remains until completion or abort | null | — | completed records retained for 24 hours and capped at 10000 per process; insertion evicts earliest expiry then lexical namespace; every abort wakes remaining waiters and leaves no completed record; no durable or external persistence dependency | idempotency |
| VS0-CF-F15-37 | FEATURE-0015 | an RFC 7396 PATCH of a FEATURE-0015 mutable field including null, and a successful authorized collection LIST | PATCH stages a clone, classifies fields before merge, removes only optional mutable fields with null, rejects null for required mutable or immutable/server-owned/deferred fields, then validates post-merge; LIST is ascending metadata.uid | VALIDATION_FAILED for null required mutable field; existing immutable/server-owned/deferred field outcomes otherwise | — | current resourceVersion comparison, required AuditEvent append, and publication occur only after post-merge validation; deterministic ordering does not broaden read visibility | feature |
| VS0-CF-F15-38 | FEATURE-0015 | a FEATURE-0015 PATCH missing, malformed, or carrying a stale If-Match resourceVersion, and a matching If-Match rechecked under the publication lock | missing/malformed/stale request unchanged; matching request may publish only after lock recheck | STALE_RESOURCE_VERSION | — | no mutation, AuditEvent, or idempotency record for a missing/malformed/stale token; no client body field added | feature |
| VS0-CF-F15-39 | FEATURE-0015 | GET or LIST for CloudPlatform, CloudProvider, topology, or CloudProviderParticipation with and without its exact server-resolved read action and scope/UID boundary | authorized LIST remains ascending metadata.uid and filtered to authorized scope; authorized GET returns only its accessible target | AUTHORIZATION_DENIED for LIST without applicable read action; RESOURCE_NOT_FOUND for inaccessible direct GET | — | cloudplatform.read and cloudprovider.read use Platform-root scope; topology.read uses CloudProvider scope; participation.read uses derived CloudPlatform scope; LIST denial is audited and inaccessible GET has safe non-disclosure | security |
| VS0-CF-F15-40 | FEATURE-0015 | a completed same-key/same-digest replay, process termination before 24 hours, or an Aborted idempotency namespace before re-reservation | replay returns stored successful status/body only while record remains; after process termination, expiry, eviction, or atomic Aborted removal the request is processed as new | null | — | replay returns Content-Type application/json only; every received request receives fresh request/correlation identity; requestId, correlation, ETag, Location, and other entity/transport headers are not replayed; effective retention is lesser of 24 hours and process lifetime | idempotency |
| VS0-CF-F15-41 | FEATURE-0015 | an authenticated principal without the exact CloudPlatform administrator writer boundary PATCHing CloudPlatform spec.description, or without the exact CloudProvider administrator writer boundary PATCHing CloudProvider spec.displayName/spec.operatingMarkets or topology spec.description | unchanged | AUTHORIZATION_DENIED | — | writer authorization denies before semantic PATCH processing and publication; no resource mutation, status/lifecycle update, no AuditEvent, and no idempotency record are published; inaccessible cross-provider reference safe-denial remains governed by VS0-CF-X03 | security |

---

## 5. Security, privacy, compatibility, and operational requirements

### 5.1 Security and authorization

- Authentication is required on every owned method/path pattern; missing/invalid authentication returns AUTH_REQUIRED (401) with no side effect (REQ-F15-24; VS0-CF-F15-31).
- Authorization uses server-resolved, deterministic bootstrap grants only; grant claims in an HTTP header or request body are ignored and never expand or substitute for the server-resolved grant (REQ-F15-17; VS0-CF-F15-22). FEATURE-0015 persists no roles, memberships, or assignments; identity-governance ownership remains with FEATURE-0018 by reference.
- Writer boundaries are consumed by reference to VS0-WRITER-002 (cloud-platform-admin CloudPlatform.spec), VS0-WRITER-003 (cloud-provider-admin CloudProvider/topology spec), VS0-WRITER-004 (participation actions, api-server), and VS0-WRITER-005 (status, api-server sole status writer for all owned kinds). Clients never write status (VS0-CF-F15-05) or api-server-owned identity/metadata (VS0-CF-F15-06).
- Cross-provider and inaccessible references use safe denial (RESOURCE_NOT_FOUND / VS0_AUTHORIZATION_SAFE_DENIAL) with no existence disclosure (REQ-F15-09; VS0-CF-X03, VS0-CF-F15-23).
- Read/list authorization boundary (normative, not a separately numbered REQ/AC): an authenticated collection LIST without a read grant is denied (403) without resource detail and is audited; LIST with a scoped read grant returns only visible resources; LIST with a scoped read grant covering no visible resources returns 200 with an empty collection; and an unauthorized direct GET or reference resolution returns safe RESOURCE_NOT_FOUND / VS0_AUTHORIZATION_SAFE_DENIAL. Local proof: VS0-CF-F15-21.

### 5.2 Privacy, classification, and redaction

- AuditEvents are redacted with request correlation and no secret or inaccessible-resource disclosure (REQ-F15-15). No credential values or protected handles appear in any audit record; references are UID-pinned identifiers only.
- Resource classification and redaction follow the registry (CloudPlatform Internal with ownerRegistration never customer-visible; CloudProvider/topology/participation Provider-confidential). Non-Active/non-visible participation phases are never end-user visible; customer/Organization visibility eligibility is out of scope (see §7).

### 5.3 Compatibility and contract reuse

- All errors reuse existing FEATURE-0012 Problem Details codes; no new top-level Problem code or violation code is introduced (REQ-F15-12; REQ-F15-23 clarification).
- FEATURE-0012 ObjectMeta, TypedRef, and Condition; FEATURE-0013 DecisionRecord and AuditEvent; and the FEATURE-0011 reuse governance contract are consumed by exact reference only and are not redefined, renamed, specialized, or transferred.
- Update surface is PATCH-only (`application/merge-patch+json`); PUT and DELETE return 405 (REQ-F15-08). Successful create returns 201; successful PATCH returns 200 (REQ-F15-19).

### 5.4 Operational requirements

- Idempotency: POST create/actions require `Idempotency-Key`; replay semantics and reuse-mismatch conflict are observable per REQ-F15-18 (VS0-CF-F15-18, VS0-CF-F15-19), covering all seven collection creates and all eight existing-participation item-action routes. The public participation surface has exactly eight item-action routes (accept, reject, withdraw, suspend, resume, request-release, accept-release, decline-release).
- Completed idempotency records retain for 24 hours and are capped at 10,000 per process; deterministic completed-only eviction means an expired/evicted key has no replay guarantee, while InFlight reservations are never evicted. A cancelled waiter detaches only itself and every abort wakes remaining waiters (VS0-CF-F15-36).
- Successful LIST output is ascending `metadata.uid`; exact read actions are `cloudplatform.read`, `cloudprovider.read`, `topology.read`, and `participation.read` at their registered derived scopes with optional GET target-UID narrowing. PATCH uses staged RFC 7396 field classification and post-merge validation; `null` removes only a mutable optional field, and its current `If-Match` is rechecked under the publication lock (VS0-CF-F15-37, VS0-CF-F15-38, VS0-CF-F15-39).
- Expiry is an internal api-server transition invoked by the deterministic scheduler at `requestExpiresAt`, not a public participation action; the deterministic scheduler is a system actor, not an additional writer or controller authority (VS0-CF-F15-16).
- FEATURE-0015 introduces no durable or external persistence dependency, which is why DEPENDENCY_UNAVAILABLE (503) is never used for the audit-append-failure outcome (REQ-F15-15).
- Mutation and its required AuditEvent are one atomic outcome; audit-append failure leaves the resource/lifecycle/idempotency unpublished and returns INTERNAL_ERROR (500) (REQ-F15-15; VS0-CF-F15-24).

---

## 6. Edge cases

| Edge case | Observable outcome | Proof |
|-----------|--------------------|-------|
| Create of CloudProvider/topology/participation before any CloudPlatform exists | CONFLICT (409) VS0_CLOUDPLATFORM_ROOT_REQUIRED; no persistence; audited denial | VS0-CF-F15-12 |
| Second non-terminal participation for the same (CloudPlatform, CloudProvider) pair | ALREADY_EXISTS (409) VS0_PARTICIPATION_DUPLICATE; no second participation; no duplicate audit | VS0-CF-F15-04 |
| Caller without `participation.accept.platform` posts accept for Pending participation | AUTHORIZATION_DENIED (403); remains Pending with no lifecycle/status write, AuditEvent, or idempotency completion | VS0-CF-F15-09 |
| Resume clearing only one hold while the other remains true | Remains effectively Suspended; only the acting hold changes | VS0-CF-F15-08 |
| Suspend/resume while phase is Pending/Terminating/terminal | CONFLICT (409) VS0_PARTICIPATION_STATE_INVALID; no hold change | VS0-CF-F15-10 |
| Action from an invalid source state | Unchanged; CONFLICT (409) VS0_PARTICIPATION_STATE_INVALID | VS0-CF-F15-20 |
| Missing/malformed/non-current If-Match on existing-participation action | STALE_RESOURCE_VERSION (412); no lifecycle/status write; no additional AuditEvent | VS0-CF-F15-17 |
| If-Match supplied on participation collection create | MALFORMED_REQUEST (400); no mutation | VS0-CF-F15-30 |
| Missing/empty/malformed/over-length Idempotency-Key where required | MALFORMED_REQUEST (400); no mutation/idempotency record/AuditEvent | VS0-CF-F15-30 |
| Ordinary non-empty JSON body on a participation item action without the reserved top-level `bootstrapGrant` member in a syntactically valid duplicate-free body | MALFORMED_REQUEST (400); no mutation | VS0-CF-F15-30 |
| Syntactically valid duplicate-free participation item-action body containing top-level `bootstrapGrant` | AUTHORIZATION_DENIED (403); exactly one redacted AuditEvent; denied before current authorization/reference access | VS0-CF-F15-22 |
| Same key + same digest replay | Original result returned before version checking; no duplicate | VS0-CF-F15-18 |
| Same key + changed digest | CONFLICT (409) VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH; original retained | VS0-CF-F15-19 |
| cloudPlatformRef.uid ≠ scopeRef.uid | VALIDATION_FAILED (422) VS0_SCOPE_REFERENCE_MISMATCH | VS0-CF-F15-11 |
| Scope kind outside a kind's declared subset | VALIDATION_FAILED (422) VS0_SCOPE_KIND_INVALID | VS0-CF-F15-02 |
| Client supplies metadata.scopeRef, server-owned, status, unknown, or another feature's field on create | Rejected under the closed create contract | VS0-CF-F15-26, VS0-CF-F15-27 |
| Syntactically valid but unassigned ISO-3166-1 alpha-2 value (e.g. ZZ) | VALIDATION_FAILED (422) | VS0-CF-F15-29 |
| PATCH with unsupported media type or immutable/identity/scope/relationship/owner-registration field change | 415 UNSUPPORTED_MEDIA_TYPE, or 422 VS0_PATCH_IMMUTABLE_FIELD, or 403 VS0_STATUS_FIELD_WRITE; resource unchanged | VS0-CF-F15-14 |
| PUT or DELETE against CloudPlatform, CloudProvider, or a topology resource | 405; only PATCH accepted | VS0-CF-F15-07 |
| HEAD matched by ServeMux to a registered GET pattern | 405; no AuditEvent or registration added | VS0-CF-F15-34 |
| Authenticated forged header plus unsupported PATCH media type | AUTHORIZATION_DENIED (403) before media/body validation; audited | VS0-CF-F15-35 |
| LIST without read grant / with scoped grant / with no visible resources / unauthorized GET | 403 (audited); filtered results; 200 empty; 404 safe-denial respectively | VS0-CF-F15-21 |
| Header/body forged grant claim | Ignored; AUTHORIZATION_DENIED; audited without forged-claim disclosure | VS0-CF-F15-22 |
| Scheduler expiry of a still-Pending participation at requestExpiresAt | Expired; idempotent; one correlated AuditEvent naming the system actor | VS0-CF-F15-16 |
| Required AuditEvent append fails during an audited outcome | INTERNAL_ERROR (500); outcome unpublished; no DEPENDENCY_UNAVAILABLE | VS0-CF-F15-24 |

---

## 7. Non-goals and adjacent-feature exclusions

Prohibited concepts may be named when expressing exclusions, canonical ledgers, conformance evidence, or traceability, but must never become active FEATURE-0015 behavior; active requirements and acceptance cases refer to this exclusion section by reference and do not repeat prohibited names as active behavior.

### 7.1 Prohibited stale concepts (must not appear as active behavior)

ResourcePool (DEC-0042), ProviderCapability (DEC-0042), generic Provider as a combined owner/operator (DEC-0037), ServiceClass as canonical catalog (DEC-0049), EffectivePolicyContext (DEC-0050), the six-scope vocabulary (DEC-0037), SovrunnInstallation as an active Slice 0 resource (DEC-0053 deferred; VS0-SCHEMA-057 tombstone), and CanonicalMigrationPlan/CanonicalMigrationRecord/migration controller/cutover state machine (DEC-0059; ADH-2026-045; retired tombstones VS0-SCHEMA-060/061, VS0-WRITER-020/021, VS0-STATE-011 — never reused).

### 7.2 Adjacent-feature exclusions (enforceable non-goals)

| Excluded concept | Owner | Negative acceptance |
|------------------|-------|---------------------|
| ExecutionTarget (identity, schema, routes, status, writer, conformance, target lifecycle), NormalizedTargetFactSet, TargetQualificationResult, target qualification/availability/maintenanceEpoch writes | FEATURE-0016 | VS0-CF-F15-01, VS0-CF-F15-25 prove none is introduced (REQ-F15-05; AC-F15-03) |
| PolicyEvaluationRequest/Result, PolicyEngineAdapter | FEATURE-0017 | Not referenced; scope-kind and route surface bounded by VS0-CF-F15-02, VS0-CF-F15-25 |
| GovernanceProfile, Membership, RoleDefinition, RoleAssignment; inherited AUTH_REQUIRED proof VS0-CF-F01 | FEATURE-0018 | FEATURE-0015 persists no roles/memberships/assignments (REQ-F15-17); AUTH_REQUIRED proven locally by VS0-CF-F15-31, not by inherited VS0-CF-F01 |
| SovereigntyProfile, SovereigntyFactSet, EvidenceRecord | FEATURE-0019 | Not introduced |
| ProfileAssignment, EffectiveGovernanceContext | FEATURE-0020 | Not introduced |
| CloudEnrollment, EntitlementPackage, ServiceEntitlement, QuotaPolicy, provider-selection intent evaluation; `spec.providerSelectionModes`, `spec.permittedHostingLocationRefs`; Organization visibility eligibility | FEATURE-0021 | Provider-selection fields rejected on participation create (REQ-F15-21; VS0-CF-F15-28); visibility eligibility not stored |
| ServiceTypeDefinition, ServiceOffering, ServicePlan, ServiceRuntimeProfile, ServiceRequirementSet, ServiceRegion | FEATURE-0022 | Not introduced |
| DecisionProfiles, ServicePlacement, PlacementDecision; placement proof VS0-CF-F09 | FEATURE-0023 | Not introduced; VS0-CF-F09 is an integration reference only, never local evidence |
| PluginExecution, ServiceDeploymentPlan, fake harness | FEATURE-0024 | Not introduced |
| DecisionOperationContext, AI-readable projection | FEATURE-0025 | Not introduced |
| Slice 0 integration demo, cross-feature orchestration; happy-path proof VS0-CF-HP01, trace proof VS0-CF-T01 | FEATURE-0026 | Not introduced; VS0-CF-HP01 and VS0-CF-T01 are integration references only, never local evidence |
| Federation, remote deployment resource, signed remote offer, cross-control-plane replication; platform lifecycle (SovrunnRelease, PlatformLifecyclePolicy, PlatformLifecyclePlan) | No Phase 2R owner | Explicitly excluded (ADH-2026-045 decision 12; DEC-0053 deferred) |
| Real provisioning | Phase 3 | Excluded per PHASE2R non-goals |

FEATURE-0014's alpha model (Provider, ProviderLocation, ProviderDatacenter, DatacenterFailureDomain, InfrastructureStack) is a retained repository asset assessed for reuse under FEATURE-0011; it is not runtime-migrated, and ResourcePool/ProviderCapability are not adopted.

---

## 8. Architecture/decision/risk traceability by exact ID

### 8.1 Owned resources → schema/writer/state/conformance

| Owned item | DEC/ADH | VS0 Schema | VS0 Writer | VS0 State | FEATURE-0015-local conformance |
|------------|---------|------------|------------|-----------|--------------------------------|
| CloudPlatform | DEC-0037; ADH-020/042/045/046/047/050/052 | VS0-SCHEMA-008 | VS0-WRITER-002,005 | — | VS0-CF-F15-01,02,07,11,12,13,14,18,19,21,24,25,26,27,32,33 |
| CloudProvider | DEC-0037; ADH-020/042/045/046/047/050/052 | VS0-SCHEMA-009 | VS0-WRITER-003,005 | — | VS0-CF-F15-01,02,07,12,13,14,18,19,21,24,25,26,27,29,32,33 |
| CloudProviderParticipation | DEC-0054; ADH-037/042/045/046/047/050 | VS0-SCHEMA-010 | VS0-WRITER-004,005 | VS0-STATE-001 | VS0-CF-F15-01,02,03,04,08,09,10,11,12,15,16,17,18,19,20,24,25,26,27,28 |
| HostingLocation | DEC-0041; ADH-024/042/045/046/047/050/052 | VS0-SCHEMA-011 | VS0-WRITER-003,005 | — | VS0-CF-F15-01,02,07,12,13,14,18,19,21,24,25,26,27,29,33 |
| Datacenter | DEC-0041; ADH-024/042/045/046/047/050/052 | VS0-SCHEMA-012 | VS0-WRITER-003,005 | — | VS0-CF-F15-01,02,07,12,13,14,18,19,21,24,25,26,27,33 |
| FaultDomain | DEC-0041; ADH-024/042/045/046/047/050/052 | VS0-SCHEMA-013 | VS0-WRITER-003,005 | — | VS0-CF-F15-01,02,07,12,13,14,18,19,21,24,25,26,27,33 |
| InfrastructureStack | DEC-0041,0042; ADH-025/042/045/046/047/050/052 | VS0-SCHEMA-014 | VS0-WRITER-003,005 | — | VS0-CF-F15-01,02,07,12,13,14,18,19,21,24,25,26,27,33 |
| Cross-provider isolation | DEC-0037,0054; ADH-043/045 | — | VS0-WRITER-003 | — | VS0-CF-X03, VS0-CF-F15-23 |
| Scope-reference integrity | DEC-0037,0054; ADH-043/045/047 | VS0-SCHEMA-008,010 | VS0-WRITER-002,004 | — | VS0-CF-F15-11, VS0-CF-F15-27 |
| Status writer resolution (api-server sole writer) | ADH-2026-046 decision 1 | VS0-SCHEMA-008..014 | VS0-WRITER-005 | — | VS0-CF-F15-12,15,16 |
| Participation create-vs-existing preconditions; all-route idempotency coverage | ADH-2026-046 decision 2; ADH-2026-050; ADH-2026-054 | VS0-SCHEMA-008..014 | VS0-WRITER-002,003,004 | VS0-STATE-001 | VS0-CF-F15-15,17,18,19 |
| Bootstrap-grant boundary | ADH-2026-045; ADH-2026-046 decision 3; ADH-2026-054 | — | VS0-WRITER-002,003,004 | — | VS0-CF-F15-22 |
| Audit atomicity; required-append-failure mapping; durable-audit boundary | DEC-0059; ADH-2026-043 (preserved),046,049,051,054 | VS0-SCHEMA-007 | VS0-WRITER-011 (F0013) | — | VS0-CF-F15-05,06,12,16,21,22,23,24, VS0-CF-X03 |
| Authentication local proof for every owned method/path pattern | ADH-2026-051 | VS0-SCHEMA-008..014 | — | — | VS0-CF-F15-31 |
| Exact route model (22 logical paths; 35 registrations; ServeMux-compatible participation item actions) | ADH-2026-051; ADH-2026-053 | VS0-SCHEMA-008..014 | — | — | VS0-CF-F15-25 |
| Closed collection-create request contract | ADH-2026-047 decision 1 | VS0-SCHEMA-008..014 | VS0-WRITER-002,003,004 | — | VS0-CF-F15-26 |
| Scope derivation single-source proof | ADH-2026-047 decision 2 | VS0-SCHEMA-008..014 | VS0-WRITER-002,003,004 | — | VS0-CF-F15-27 |
| Topology immutability correction | ADH-2026-047 decision 3 | VS0-SCHEMA-011..014 | VS0-WRITER-003 | — | VS0-CF-F15-13,14 |
| Participation collection-create body vs empty item-action body | ADH-2026-047 decision 4 | VS0-SCHEMA-010 | VS0-WRITER-004 | VS0-STATE-001 | VS0-CF-F15-28 |
| ISO-3166-1 alpha-2 assigned-code semantics | ADH-2026-048 decision 1 | VS0-SCHEMA-008,009,011 | — | — | VS0-CF-F15-29 |
| Idempotency-Key/If-Match/item-action-body malformed-input outcomes | ADH-2026-048 decision 2; ADH-2026-054 | VS0-SCHEMA-010 | VS0-WRITER-004 | — | VS0-CF-F15-30 |
| CloudPlatform/CloudProvider name-uniqueness proof | ADH-2026-052 | VS0-SCHEMA-008,009 | VS0-WRITER-002,003 | — | VS0-CF-F15-32 |
| Registry-declared schema-constraint validation proof (excluding duplicate name, unassigned ISO code, and malformed/prohibited header/body) | ADH-2026-052 | VS0-SCHEMA-008..014 | VS0-WRITER-002,003 | — | VS0-CF-F15-33 |
| Executable HEAD, combined-input precedence, bounded idempotency, PATCH, and deterministic LIST closure | ADH-2026-055 | VS0-SCHEMA-008..014 | VS0-WRITER-002,003,004 | VS0-STATE-001 | VS0-CF-F15-34,35,36,37 |
| PATCH conditional update, exact read actions, and replay response closure | ADH-2026-056 | VS0-SCHEMA-008..014 | VS0-WRITER-002,003,004 | VS0-STATE-001 | VS0-CF-F15-38,39,40 |
| Mutable-spec writer-boundary denial and acceptance-proof closure | ADH-2026-057 | VS0-SCHEMA-008,009,011..014 | VS0-WRITER-002,003 | — | VS0-CF-F15-41 |

### 8.2 REQ → local proof (from architecture §10.1)

REQ-F15-01→VS0-CF-F15-01,02,12,13,14,32,33; REQ-F15-02→VS0-CF-F15-01,02,12,13,14,29,32,33; REQ-F15-03→VS0-CF-F15-01,03,04,15,16,20; REQ-F15-04→VS0-CF-F15-01,02,13,14,23,33; REQ-F15-05→VS0-CF-F15-01,25; REQ-F15-06→VS0-CF-F15-08,10,15,16,20; REQ-F15-07→VS0-CF-F15-03,04,15,17,18,19,20; REQ-F15-08→VS0-CF-F15-07,13,14,38; REQ-F15-09→VS0-CF-X03, VS0-CF-F15-23; REQ-F15-10→VS0-CF-F15-01,02; REQ-F15-11→VS0-CF-F15-41; REQ-F15-12→anti-drift/non-runtime (collective error-emitting cases); REQ-F15-13→anti-drift/non-runtime verification gate; REQ-F15-14→VS0-CF-F15-11; REQ-F15-15→VS0-CF-F15-24; REQ-F15-16→VS0-CF-F15-12; REQ-F15-17→VS0-CF-F15-22,39; REQ-F15-18→VS0-CF-F15-18,19,40; REQ-F15-19→VS0-CF-F15-26; REQ-F15-20→VS0-CF-F15-27; REQ-F15-21→VS0-CF-F15-28; REQ-F15-22→VS0-CF-F15-29; REQ-F15-23→VS0-CF-F15-30; REQ-F15-24→VS0-CF-F15-31. The unnumbered read/list authorization boundary (§5.1) traces to VS0-CF-F15-21,39.

### 8.3 Error codes (consumed by reference from FEATURE-0012 closed set)

AUTH_REQUIRED/401, AUTHORIZATION_DENIED/403, RESOURCE_NOT_FOUND/404, ALREADY_EXISTS/409, CONFLICT/409, MALFORMED_REQUEST/400, UNSUPPORTED_MEDIA_TYPE/415, STALE_RESOURCE_VERSION/412, VALIDATION_FAILED/422, INTERNAL_ERROR/500; HTTP 405 is a transport outcome for unsupported methods and has no top-level Problem code; 201 create; 200 PATCH. Slice 0 violation codes (VS0_PARTICIPATION_STATE_INVALID, VS0_PARTICIPATION_DUPLICATE, VS0_STATUS_FIELD_WRITE, VS0_SYSTEM_OWNED_FIELD_WRITE, VS0_SCOPE_KIND_INVALID, VS0_PATCH_IMMUTABLE_FIELD, VS0_CLOUDPLATFORM_ROOT_REQUIRED, VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH, VS0_SCOPE_REFERENCE_MISMATCH, VS0_TOPOLOGY_PROVIDER_MISMATCH, VS0_AUTHORIZATION_SAFE_DENIAL, VS0_AUTH_REQUIRED) appear only in `violations[].code`, never as top-level Problem codes. No new top-level or violation code is introduced.

### 8.4 Risk/anti-drift traceability

Architecture boundary §11 rules 1–12 are enforced through their applicable §7 exclusions/non-goals and §4.26 acceptance scenarios, including the explicit anti-drift/non-runtime cases where applicable. Rule 13 is a proof-integrity control, not a non-goal or acceptance scenario: REQ-F15-01, REQ-F15-02, and REQ-F15-04 require the exact validation evidence; VS0-CF-F15-32 proves name uniqueness and VS0-CF-F15-33 proves other registry-declared schema constraints in §4.28; and §8.2 records the corresponding REQ-to-local-proof mappings.

---

## 9. Design questions explicitly delegated by architecture

The architecture explicitly delegates the following as deterministic design mechanics (DESIGN-delegated per architecture §2.2, §9, and §7.12). They are recorded here as design-stage questions only; they are not semantic or ownership questions and introduce no new behavior. FEATURE-0015 introduces no durable or external persistence dependency (the only requirements-owned observable persistence boundary); the implementation data structures and coordination mechanics that satisfy this boundary are design-owned, not selected here:

1. Registry data structures for the cloud-model resources, consistent with the no-durable-or-external-persistence-dependency boundary.
2. Handler wiring and Go 1.22 `http.ServeMux` registration for the exact route model in architecture §7.13 (22 logical paths; 35 explicit method/path patterns), using no path-only or catch-all registration, reflection, or auto-registration. The approved complete `{uid}` wildcard segment remains required where the route model declares it.
3. Deterministic validation plumbing (pure functions).
4. Bootstrap grant resolution internals.
5. Idempotency reservation-state mechanics and waiter/abort/release handling on validation failure, stale version, audit failure, or recovered panic; idempotency records written only for completed replayable outcomes consistent with the approved audit policy.
6. Audit/mutation coordination mechanics ensuring a resource/idempotency change is not published until its required audit append succeeds, consistent with the no-durable-or-external-persistence-dependency boundary.

No design question reopens a semantic or ownership decision; all are bounded by the approved contract.

---

## 10. Completeness and unresolved-decision report

- Requirements ledger: all 24 approved REQ rows copied verbatim in §4.25; each has exactly one normative detail heading in §4.1–§4.24 in approved order.
- Acceptance ledger: all 20 approved AC rows copied verbatim in §4.27; each mapped explicitly in §4.26 with its FEATURE-0015-local proof case(s); anti-drift/non-runtime items (AC-F15-05, AC-F15-07) and contract-reuse items (AC-F15-09) are labelled and carry no runtime proof case per ADH-2026-046 decision 3.
- Conformance semantics ledger: §4.28 copies exact registry fields for VS0-CF-X03 and VS0-CF-F15-01 through VS0-CF-F15-41. Downstream-owned cases (VS0-CF-HP01, VS0-CF-F09, VS0-CF-F01, VS0-CF-T01) are cited only as non-owning integration references in §7 and are never used as local acceptance evidence.
- Inherited contracts (FEATURE-0011 reuse governance; FEATURE-0012 Problem Details/ObjectMeta/TypedRef/Condition; FEATURE-0013 DecisionRecord/AuditEvent) are consumed by exact reference; none is redefined or transferred.
- No downstream status field, metric, side effect, controller, adapter/plugin execution, or runtime proof is imported. Prohibited concepts may be named when expressing exclusions, canonical ledgers, conformance evidence, or traceability, but must never become active FEATURE-0015 behavior.
- No orphan IDs: every referenced REQ, AC, DEC, ADH, VS0-SCHEMA, VS0-WRITER, VS0-STATE, VS0-CF, and violation/Problem code exists in the loaded authorities.
- Unresolved architecture decisions: none. No `ARCHITECTURE_DECISION_REQUIRED`, `REQUIREMENT_CLARIFICATION_REQUIRED`, `BOUNDARY_CHANGE_REQUIRED`, `DEPENDENCY_APPROVAL_REQUIRED`, or `SECURITY_REVIEW_REQUIRED` condition was triggered; all requirements map one-to-one to existing registry entries and approved decisions.

---

## Model Execution Report:
- Tool: kiro
- Stage or task: FEATURE-0015 Requirements (`.kiro/specs/canonical-cloud-model-and-alpha-migration/requirements.md`)
- Recommended priority list: 1) Architecture-heavy (`claude-opus-4.8`), effort: high; 2) Fallback (`claude-sonnet-4.5`), effort: medium
- Selected model: claude-opus-4.8
- Effort/reasoning setting: high
- Fallback used: no
- Fallback reason: none
STAGE_STATUS: COMPLETE
