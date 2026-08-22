# VS-000 Slice 0 Exact Contract Specification

| Field | Value |
|-------|-------|
| Status | Experimental |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Controlling Adoption | ADH-2026-042 |
| Machine-Readable Registry | `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` |

## Scope Statement

This document is the normative specification for the exact minimal synthetic Slice 0 integration profile only.

It is **not**: a full schema for FEATURE-0015 through FEATURE-0026; not a feature requirements, design, or tasks artifact; not a future slice definition.

The machine-readable normative registry (`VS-000-contract-registry.yaml`) is the authoritative source for field-level detail. This specification defines the rules that govern that registry.

FEATURE-0012 (common grammar, Problem Details) and FEATURE-0013 (DecisionRecord, AuditEvent) schemas are reused by reference. Parallel envelopes for the same concerns are forbidden.

---

## 1. Authority and Precedence

1. Canonical data model (`docs/architecture/canonical/sovrunn-finalized-data-model.md`) is semantic owner of all resource kinds.
2. Canonical contract catalog (`docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`) is semantic owner of all API contracts.
3. This specification defines the Slice 0 narrowing profile; it cannot broaden canonical semantics.
4. The registry YAML is the machine-readable expression of this specification.
5. Precedence: active baseline > canonical model/catalog > accepted DEC/ADH records > VS-000 charter > this specification > registry YAML > approved feature stage > implementation code.
6. Seven canonical scope kinds apply: Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider.

---

## 2. Schema-Profile Rules

Every kind or DTO exchanged or persisted within Slice 0 must have a registry entry containing:

| Attribute | Requirement |
|-----------|-------------|
| Registry ID | Unique `VS0-SCHEMA-<NNN>` identifier |
| Identity | Primary key fields and uniqueness constraint |
| Profile | Slice 0 narrowing profile name |
| Boundary | Which layer owns creation/mutation |
| Scope | One of seven canonical scope kinds |
| Fields | Name, type, requiredness, enums, defaults, limits |
| References | Foreign-key relationships to other registry entries |
| Classification/Redaction | Sensitivity label; redaction rule for projection |
| Mutability/Retention | Immutable-after-create, mutable fields, retention class |
| Validation | Deterministic validation rule references |
| Projection | Safe customer-facing projection subset |

FEATURE-0012 common grammar (request/response envelope, pagination, error shape) and FEATURE-0013 DecisionRecord/AuditEvent schemas are incorporated by reference to their canonical definitions. No Slice 0 entry may redefine them.

Slice 0 narrowing may restrict cardinality, constrain enums, or tighten limits. It must not add fields, relax constraints, or introduce semantics absent from the canonical model.

---

## 3. Exact Writer Rules

| Rule ID | Rule |
|---------|------|
| VS0-W01 | Every mutable field has exactly one authoritative writer domain. |
| VS0-W02 | Clients (Portal, CLI, API consumers) never write `status`. |
| VS0-W03 | The deployment planner records an already-authorized target set; it never selects targets, executes steps, or changes decisions. |
| VS0-W04 | Decision services never mutate their inputs; they emit DecisionRecords. |
| VS0-W05 | Adapters write only normalized observations or execution results. |
| VS0-W06 | Projection endpoints emit read-only safe subsets; no write path. |
| VS0-W07 | Binding writers validate consumer entitlement before credential issuance. |
| VS0-W08 | Conflicts between concurrent writers fail closed (reject, do not merge). |

Writer domain table (normative rows in registry under `VS0-WRITER-<NNN>`):

| Writer Domain | Writes To | Never Writes |
|---------------|-----------|--------------|
| Client | `spec` of owned resources | `status`, DecisionRecord |
| Controller | `status`, Operation state | `spec` |
| Deployment planner | Immutable plan from final decisions and pinned inputs | Target selection, execution results, credentials |
| Decision Service | DecisionRecord | Input resources, `status` |
| Adapter | Observations, execution results | `spec`, DecisionRecord |
| Projection controller | Append-only ServicePlacement and transient explanation projection | Canonical decisions, provider topology, secret values |

---

## 4. Exact State-Machine Rules

1. Registry transitions (under `VS0-STATE-<NNN>`) are authoritative for allowed state changes and guards.
2. Milestones are observable checkpoints; they are not lifecycle phases.
3. Immutable Slice 0 records are persisted directly as `FINAL` and are append-only thereafter; `Created` is not a persisted record state.
4. Operation owns `Pending → Running → Succeeded | Failed | Cancelled`; pending, retrying and waiting are never DecisionRecord states.
5. Deletion of a ServiceInstance revokes all associated ServiceBindings before completing.
6. Deletion retains immutable audit/decision records; it does not purge history.
7. Invalid transitions are rejected; the system does not silently coerce state.

State-machine summary (exact definitions in registry):

| Kind Group | States | Terminal | Registry Prefix |
|------------|--------|----------|-----------------|
| Participation | Pending, Active, Rejected, Withdrawn, Expired, Suspended, Terminating, Terminated | Rejected, Withdrawn, Expired, Terminated | VS0-STATE-001 |
| CloudEnrollment | Pending, Active, Suspended, Terminating, Terminated | Terminated | VS0-STATE-002 |
| Published definitions | Draft, Published, Retired | Retired | VS0-STATE-003 |
| Target qualification/availability | Orthogonal qualified/available combinations in the registry | None | VS0-STATE-004 |
| Quota reservation | None, Reserved, Committed, Released | Released | VS0-STATE-005 |
| ServiceInstance | Pending, Provisioning, Ready, Failed, Deleting, Deleted | Deleted | VS0-STATE-006 |
| ServiceBinding | Pending, Bound, Revoking, Revoked, Failed | Revoked | VS0-STATE-007 |
| Operation | Pending, Running, Succeeded, Failed, Cancelled | Succeeded, Failed, Cancelled | VS0-STATE-008 |
| PluginExecution | Pending, Running, Succeeded, Failed | Succeeded, Failed | VS0-STATE-009 |
| Immutable records | FINAL | FINAL | VS0-STATE-010 |

---

## 5. Errors

### RFC 9457 Problem Details Members

All error responses use RFC 9457 Problem Details with the FEATURE-0012 Sovrunn extensions. Required contract members are `type`, `title`, `status`, `detail`, `instance`, `code` and `requestId`; `violations` is present when field-level details apply.

### Top-Level Error Code Mapping

Exact closed set (normative rows in the registry under `problemCodes`):

| Codes | HTTP |
|-------|------|
| `MALFORMED_REQUEST`, `UNKNOWN_FIELD`, `DUPLICATE_FIELD`, `REQUEST_TOO_LARGE` | 400 |
| `AUTH_REQUIRED` | 401 |
| `AUTHORIZATION_DENIED` | 403 |
| `RESOURCE_NOT_FOUND` | 404 |
| `CONFLICT`, `ALREADY_EXISTS`, `DELETE_BLOCKED` | 409 |
| `STALE_RESOURCE_VERSION` | 412 |
| `UNSUPPORTED_MEDIA_TYPE` | 415 |
| `VALIDATION_FAILED` | 422 |
| `INTERNAL_ERROR` | 500 |
| `DEPENDENCY_UNAVAILABLE` | 503 |

For every code, `type` is exactly `urn:sovrunn:problem:<lowercase-kebab-code>`. Slice 0 does not add a top-level code or HTTP mapping. Quota exhaustion is a `CONFLICT`/409 with violation `VS0_QUOTA_EXHAUSTED`; inaccessible cross-scope resources use safe `RESOURCE_NOT_FOUND`/404.

### Validation Detail Violation Codes

Validation details use only FEATURE-0012 shared violations, the inherited FEATURE-0013 decision violations, or the closed Slice 0 entries under registry `violationCodes.slice0`. A violation carries its exact RFC 6901 JSON Pointer or `null` when a safe field path cannot be disclosed.

### Inherited Schemas

FEATURE-0013 registry defines AuditEvent and DecisionRecord error semantics. They are inherited, not duplicated.

### Safety Rules

- Error responses must not confirm existence of resources the caller cannot access (safe denial).
- An inaccessible cross-provider reference returns RESOURCE_NOT_FOUND (404) with violation VS0_AUTHORIZATION_SAFE_DENIAL and no existence disclosure. An authorized but structurally invalid same-provider reference returns VALIDATION_FAILED (422). VS0-CF-X03 records the safe-denial violation exactly.
- CloudProviderParticipation `spec.cloudPlatformRef.uid` must equal its CloudPlatform `metadata.scopeRef.uid`. CloudPlatform has no Organization reference to validate; it carries an immutable `spec.ownerRegistration` instead (ADH-2026-045 decision 8). A mismatch on the participation reference returns VALIDATION_FAILED (422) with VS0_SCOPE_REFERENCE_MISMATCH.
- `POST` create of a CloudProviderParticipation requires `Idempotency-Key` only; it does not require or accept `If-Match` because it targets an absent resource, keyed atomically by the non-terminal `(cloudPlatformUID, cloudProviderUID)` pair. PATCH and every action on an existing participation require both `If-Match` and `Idempotency-Key` (ADH-2026-045/046).
- Idempotency uses authenticated principal, registered method/path pattern, concrete target UID when the route has `{uid}`, and `Idempotency-Key`. Current authorization and safe access succeed before a completed replay is returned; only successful collection creates and successful participation actions complete records. Validation, uniqueness, lifecycle, stale-version, audit failure, panic, and shutdown abort and wake waiters (ADH-2026-054).
- ADH-2026-055 closes remaining executable behavior: a HEAD request that ServeMux matches to a GET pattern returns 405 without adding a registration; an authenticated `X-Sovrunn-Bootstrap-Grant` denial precedes media/body validation, while unsupported PATCH media precedes body-carrier decoding when that header is absent; completed idempotency records retain for 24 hours, are capped at 10,000 per process, and evict earliest-expiry then lexical namespace; PATCH is staged RFC 7396 with classified null/post-merge validation; successful LIST is ascending `metadata.uid`.
- ADH-2026-056 closes remaining conditional-update/read/replay behavior: every mutable-resource PATCH requires current `If-Match` and returns `STALE_RESOURCE_VERSION`/412 when missing, malformed, or stale, with its version rechecked under the publication lock; reads use `cloudplatform.read`, `cloudprovider.read`, `topology.read`, or `participation.read` at their registered derived scope with optional direct-GET UID narrowing; completed replay returns only stored successful status/body plus `Content-Type: application/json`, every request receives fresh request correlation, and effective replay retention is the lesser of 24 hours and process lifetime.
- ADH-2026-057 closes mutable-spec writer-enforcement proof: `VS0-CF-F15-41` denies a principal outside the exact CloudPlatform or CloudProvider writer boundary with `AUTHORIZATION_DENIED` before semantic PATCH processing or publication, with no mutation, lifecycle/status update, AuditEvent, or idempotency record. This does not alter cross-provider safe-denial.
- FEATURE-0015 status fields for `CloudPlatform`, `CloudProvider`, `HostingLocation`, `Datacenter`, `FaultDomain`, `InfrastructureStack`, and `CloudProviderParticipation` are written solely by `api-server`; no separate topology, stack, provider, or participation controller authority exists. The deterministic system scheduler is a system actor that invokes the api-server's expiry transition and is not an additional status writer (ADH-2026-046 decision 1).
- No error `detail` or `title` string is parsed programmatically by clients; codes and types are the contract.
- For every F0015 collection `POST`, the client supplies only `metadata.name`, the required F0015-owned `spec` fields, and the optional F0015-owned `spec` fields registered for that kind (VS0-SCHEMA-008 through VS0-SCHEMA-014 `requestContract`). The client must not supply `metadata.uid`, generation, resourceVersion, timestamps, `metadata.scopeRef`, any `status` field, a field owned by another feature, or an unknown field; such requests are rejected. The API server assigns identity/version/timestamps, `metadata.scopeRef`, initial status, and required FEATURE-0013 AuditEvent evidence. A successful create returns `201` with the created resource; a successful PATCH returns `200` with the updated resource (ADH-2026-047 decision 1).
- Scope derivation is server-side and has a single source per resource kind (ADH-2026-047 decision 2): `CloudPlatform` and `CloudProvider` receive the immutable deployment Platform-root scope; `CloudProviderParticipation` receives its immutable CloudPlatform scope derived from its resolved `spec.cloudPlatformRef` (the referenced CloudPlatform must exist and its UID must equal the resulting `metadata.scopeRef.uid`); `HostingLocation` receives its immutable CloudProvider scope from the CloudProvider UID bound to the authenticated, server-resolved `topology.write` grant (a caller cannot supply or choose `metadata.scopeRef`); `Datacenter`, `FaultDomain`, and `InfrastructureStack` receive their immutable CloudProvider scope from their resolved immutable parent reference. Reference resolution and authorization/non-disclosure occur before structural validation; derived-scope mismatches trigger the established safe-denial/validation outcome.
- `metadata.name` is immutable identity for every F0015 resource, including topology resources (`HostingLocation`, `Datacenter`, `FaultDomain`, `InfrastructureStack`). For topology resources, only `spec.description` is PATCHable: "Description is PATCHable; name is immutable identity" (ADH-2026-047 decision 3, correcting ADH-2026-045's superseded "Name and description are PATCHable" wording).
- `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations` (collection create) requires exactly `metadata.name`, `spec.cloudPlatformRef` (UID-pinned), `spec.cloudProviderRef` (UID-pinned), and `spec.environment` (only allowed value `development`). It accepts no optional participation `spec` fields in F0015; `spec.providerSelectionModes` and `spec.permittedHostingLocationRefs` are FEATURE-0021-introduced and FEATURE-0021-activated fields (VS0-SCHEMA-010 `fieldOwnership`) that F0015 neither accepts, stores, defaults, validates, nor exposes. The server derives `metadata.scopeRef` from `spec.cloudPlatformRef`, assigns `status.phase=Pending`, `status.platformSuspended=false`, `status.providerSuspended=false`, and `status.requestExpiresAt=createdAt+7 days`, and returns `201`. This create requires `Idempotency-Key` only and rejects/does not accept `If-Match`.

  The eight existing-participation item-action paths are:

  - `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/accept`
  - `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/reject`
  - `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/withdraw`
  - `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/suspend`
  - `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/resume`
  - `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/request-release`
  - `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/accept-release`
  - `POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/decline-release`

  In each path, `{uid}` is a complete Go 1.22 `http.ServeMux` wildcard segment, read through `r.PathValue("uid")`. Each action has an empty JSON body and requires both `If-Match` and `Idempotency-Key`. Scheduler expiry is an internal API-server transition with no public request body (ADH-2026-047 decision 4; ADH-2026-053).
- For FEATURE-0015, an ISO-3166-1 alpha-2 value (`CloudPlatform.spec.ownerRegistration.jurisdictionCode`, `CloudProvider.spec.operatingMarkets[]`, `HostingLocation.spec.countryCode`) means an **assigned** code from a fixed, repository-owned, version-pinned static dataset (as of `2026-08-12`), not merely a syntactically valid two-letter code. A syntactically valid but unassigned code (e.g. `ZZ`) is rejected with `VALIDATION_FAILED` (422). `HostingLocation.spec.administrativeAreaCode` remains syntax-plus-country-prefix validation only; F0015 introduces no ISO-3166-2 membership dataset (ADH-2026-048 decision 1).
- Four exact malformed/prohibited-input outcomes reuse only existing FEATURE-0012 top-level Problem codes; no new top-level Problem code or violation code is introduced (ADH-2026-048 decision 2):
  - A required `Idempotency-Key` that is missing, empty, malformed, or exceeds the shared length limit returns `MALFORMED_REQUEST` (400).
  - An `If-Match` header supplied on a `CloudProviderParticipation` collection-create request returns `MALFORMED_REQUEST` (400).
  - An `If-Match` header that is missing, malformed, or non-current on an action against an existing participation returns `STALE_RESOURCE_VERSION` (412), preserving the ADH-2026-046 rule.
  - A `CloudProviderParticipation` item action accepts only EOF/zero bytes as empty. Whitespace, `{}`, or another non-empty body returns `MALFORMED_REQUEST` (400), except that a valid duplicate-free top-level `bootstrapGrant` member follows its existing audited `AUTHORIZATION_DENIED` outcome.
  - A malformed/prohibited input under this rule produces no mutation, no idempotency record, and no AuditEvent unless an already-approved audited-denial rule independently applies; this clarification does not broaden the ADH-2026-046/047 audit-denial list.
- A failure to append the FEATURE-0013 AuditEvent required for a FEATURE-0015 mutation returns the inherited `INTERNAL_ERROR` (500) Problem Details response; the resource mutation and its idempotency completion are not published. `DEPENDENCY_UNAVAILABLE` (503) is not used for this outcome because FEATURE-0015 introduces no durable or external persistence dependency (ADH-2026-049, extended to audited denials by ADH-2026-051).
- After authentication and current authorization are evaluated, FEATURE-0015 produces exactly one redacted FEATURE-0013 AuditEvent, with request correlation and no secret or inaccessible-resource disclosure, for: a successful resource collection create or PATCH; a successful participation create/action or deterministic scheduler expiry; an authenticated client attempt to write api-server-owned status or identity/metadata; a CloudPlatform-root creation denial; an authenticated collection LIST without a read grant; an authenticated `X-Sovrunn-Bootstrap-Grant` header or valid duplicate-free top-level `bootstrapGrant` body member; and an authenticated inaccessible cross-provider/safe-denial reference. Missing/invalid authentication remains 401 with no AuditEvent; the reserved header is denied before body decode; malformed/oversized/duplicate bodies retain their 400 no-AuditEvent outcomes; and the valid reserved body member is denied before authorization/reference access. All other unknown body fields remain 400 with no AuditEvent (ADH-2026-054).
- Missing or invalid authentication on any FEATURE-0015 owned method/path pattern returns `AUTH_REQUIRED` (401) with no resource mutation, lifecycle/status write, idempotency record, or AuditEvent; this is proven by the FEATURE-0015-local `VS0-CF-F15-31` case and does not replace or modify the inherited FEATURE-0018 `VS0-CF-F01` (ADH-2026-051).
- ADH-2026-058 atomically replaces the FEATURE-0016 `VS0-SCHEMA-015..017`/`VS0-STATE-004` placeholders with one closed ExecutionTarget executable contract: exactly five Go 1.22 `http.ServeMux` registrations (collection POST/GET, item GET, qualify POST, retire POST) behind one pre-ServeMux transport-only method/path guard that supplies every unmatched-method/trailing-slash outcome with no Problem body, authentication, audit, idempotency, or lifecycle effect; the sole observer `sovrunn.synthetic-iaas-observer/v1` (target-bound, clock-driven, external-effect-free); `ExecutionTargetLifecycleService` as the sole committer of ExecutionTarget status and internal `NormalizedTargetFactSet`/`TargetQualificationResult` records; a lifecycle-service-only `Qualifying` in-flight reservation that is never persisted or publicly projected; and a response-only `effectiveAvailability` (Available/Unavailable/Maintenance) projection with no persisted `status.availability`, no `Draining`, and no independent writer or ETag effect. Idempotency applies only to ExecutionTarget create, qualify, and retire; GET/LIST never touch idempotency. The eight closed local violation codes are `VS0_TARGET_RETIRED`, `VS0_TARGET_MAINTENANCE`, `VS0_TARGET_QUALIFICATION_IN_PROGRESS`, `VS0_TARGET_EPOCH_STALE`, `VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE`, `VS0_EXECUTION_TARGET_STACK_UNAVAILABLE`, `VS0_EXECUTION_TARGET_SCOPE_MISMATCH`, and `VS0_EXECUTION_TARGET_VIABILITY_STALE`; FEATURE-0012 owns all top-level Problem codes. No new top-level code or HTTP mapping is introduced. Full clause-by-clause detail is normative in `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md` and the registry; this specification does not duplicate it.
- ADH-2026-060 makes the FEATURE-0016 qualification-commit Maintenance-race outcomes mutually exclusive and ordered: an active current-Maintenance marker always selects `CONFLICT`/409 + `VS0_TARGET_MAINTENANCE` first, regardless of whether the captured maintenance epoch also changed; only when no active marker exists does a changed captured maintenance epoch select `STALE_RESOURCE_VERSION`/412 + `VS0_TARGET_EPOCH_STALE`; a changed referenced InfrastructureStack generation remains the independent `STALE_RESOURCE_VERSION`/412 + `VS0_TARGET_EPOCH_STALE` case. No top-level Problem code, field, route, state, or dependency changes.
- ADH-2026-061 corrects the FEATURE-0016 Maintenance-entry rule: every successful Maintenance entry clears both current `NormalizedTargetFactSet`/`TargetQualificationResult` links and persists `Active/Unqualified`, whether or not a qualification is in flight. If a qualification is in flight, Maintenance entry additionally aborts only that qualification and its linked idempotency reservation and wakes only its waiters; the qualifying caller keeps the existing `CONFLICT`/409 + `VS0_TARGET_MAINTENANCE` outcome with no qualification result, completion, or qualification AuditEvent. If no qualification is in flight, no idempotency reservation is cancelled. `VS0-CF-F16-42`, `VS0-CF-F16-43`, and `VS0-CF-F16-89` keep their existing exact observable semantics. No top-level Problem code, field, route, state, or dependency changes.
- ADH-2026-063 corrects the FEATURE-0016 `POST /apis/execution.sovrunn.io/v1alpha1/execution-targets` create pipeline: the bounded single body read returns the existing `REQUEST_TOO_LARGE`/400 first for an oversized body before phase-one extraction, authorization, or safe access; otherwise the exact bytes are retained and phase one extracts only `spec.cloudProviderParticipationRef.uid` and `spec.infrastructureStackRef.uid`, emitting no body-classification Problem and never classifying, canonicalizing, digesting, or reserving a body; derived-scope `executiontarget.write` authorization and safe backing access are then evaluated; only then are the retained same bytes strictly classified exactly once, followed by graph/reference validation, idempotency, audit, and publication. Invalid, unknown, and duplicate bodies are never digested or reserved. Action `If-Match` header-form precedence is preserved. Four new sequential local conformance rows `VS0-CF-F16-123..126` cover unauthorized-with-malformed/duplicate (existing audited 403 `AUTHORIZATION_DENIED`), inaccessible-backing-with-malformed/duplicate (existing audited safe 404 `RESOURCE_NOT_FOUND` + `VS0_AUTHORIZATION_SAFE_DENIAL`), authorized strict classification (existing 400 `MALFORMED_REQUEST`/`DUPLICATE_FIELD`), and oversized body (existing 400 `REQUEST_TOO_LARGE`); each with no body classification/mutation/idempotency-digest/reservation/completion beyond the applicable existing outcome. No top-level Problem code, field, route, state, or dependency changes.
- ADH-2026-065 makes the FEATURE-0016 create classification and phase-one reference cases exact without altering ADH-2026-063 precedence. `VS0-CF-F16-125` is now the exact non-family malformed-JSON classification case returning the existing 400 `MALFORMED_REQUEST` (as `VS0-CF-F16-51`), and the new `VS0-CF-F16-127` is the exact non-family duplicate-top-level-member classification case returning the existing 400 `DUPLICATE_FIELD` (as `VS0-CF-F16-52`); `VS0-CF-F16-125` and `VS0-CF-F16-127` are distinct exact non-family cases, not a mixed equivalence family. The new `VS0-CF-F16-128` is the fail-closed phase-one extraction case: when the existing bounded body read succeeds but phase one cannot extract exactly one syntactically usable UID at both `spec.cloudProviderParticipationRef.uid` and `spec.infrastructureStackRef.uid`, the create returns the existing audited 403 `AUTHORIZATION_DENIED` before derived-scope authorization and safe backing access, discloses no body-classification detail and no backing-resource existence, and performs no strict classification, canonicalization, digest, idempotency reservation, observer invocation, mutation, publication, or completion. `VS0-CF-F16-123`, `VS0-CF-F16-124`, and `VS0-CF-F16-126` remain unchanged. The FEATURE-0016 local conformance range extends to `VS0-CF-F16-128`. No route, field, state, writer, dependency, top-level Problem code, or external effect changes.

---

## 6. Conformance

### Conformance Classes

Each conformance requirement has a stable registry identifier. `VS0-CF-HP01` owns the end-to-end happy path; `VS0-CF-F01` through `VS0-CF-F20` map one-to-one to the charter failures; `X`, `L`, `Z`, `T`, `I` and `D` cases cover scope, leakage, external effects, traceability, concurrency and deletion ordering. Feature-local pre-integration evidence uses `VS0-CF-F<feature>-<case>` and must not reuse a downstream-owned scenario. FEATURE-0015 has no migration proof: per DEC-0059 (canonical bootstrap replaces alpha runtime migration), there is no `CanonicalMigrationPlan`/`CanonicalMigrationRecord`, and `VS0-CF-MIG01..MIG02`/`VS0-CF-MIGF01..MIGF03` are retired.

| Registry range | Class | Verification |
|----------------|-------|--------------|
| `VS0-CF-HP01` | Positive integration | All 20 charter steps complete deterministically. |
| `VS0-CF-F01..F20` | Required failures | Exact state, Problem/status code and side-effect assertions for every charter failure. |
| `VS0-CF-X01..X03` | Scope isolation | Cross-Organization, cross-Project and cross-provider references deny safely. |
| `VS0-CF-L01` | Leakage prevention | Responses, logs, audit and explanation contain no prohibited data. |
| `VS0-CF-Z01` | Zero external effect | Fake execution records `externalCallCount=0`. |
| `VS0-CF-T01` | Correlation | Request, decisions, plan, operation, executions and audit form one complete trace graph. |
| `VS0-CF-I01..I02` | Idempotency and races | Same payload converges; different payload conflicts; one quota/operation wins. |
| `VS0-CF-D01` | Deletion ordering | Binding revocation precedes cleanup, quota release and instance finalization. |
| `VS0-CF-F15-01..F15-41` | FEATURE-0015 local evidence | Canonical resource/scope validation, participation lifecycle actions and independent suspension holds, writer denial, PATCH-only update surface (no PUT/DELETE), scope reference UID invariants, CloudPlatform-root precondition, PATCH accept/deny cases, participation create-versus-existing-action preconditions and idempotency, read/list authorization, bootstrap-grant boundary, topology reference ordering, audit atomicity, route/action surface closure, per-kind collection-create request boundary (F15-26), scope derivation single-source proof (F15-27), participation collection-create body versus empty item-action body (F15-28, ADH-2026-047), ISO-3166-1 alpha-2 assigned-code rejection (F15-29), the four header/body malformed-input outcomes (F15-30, ADH-2026-048), missing/invalid authentication on any owned method/path pattern (F15-31, ADH-2026-051), CloudPlatform/CloudProvider name-uniqueness rejection (F15-32), registry-declared schema-constraint validation (F15-33, ADH-2026-052), HEAD, combined-input precedence, bounded idempotency, RFC 7396 PATCH, and deterministic LIST closure (F15-34..37, ADH-2026-055), PATCH conditional update, exact read actions, and replay response semantics (F15-38..40, ADH-2026-056), and exact mutable-spec writer-boundary denial (F15-41, ADH-2026-057). FEATURE-0015 owns exactly 22 logical endpoint paths and registers exactly 35 explicit Go 1.22 `http.ServeMux` method/path patterns, with no path-only or catch-all handler registration, reflection, or internal HTTP-method dispatch; the declared complete `{uid}` segment remains a permitted ServeMux wildcard (ADH-2026-051,053). |
| `VS0-CF-F16-01..F16-128` | FEATURE-0016 local evidence | ADH-2026-058 closes the complete ExecutionTarget executable contract: closed create-field boundary and derived-scope/tuple-uniqueness validation; the exact five-route/pre-ServeMux-guard transport surface (HEAD, trailing-slash, and every other unmatched method); request grammar (Idempotency-Key, media type, If-Match header-form versus current-version comparison, zero-byte action bodies); executable request precedence; idempotency namespace/digest/replay/eviction/waiter rules scoped to create, qualify, and retire only; the sole synthetic observer's four named facts and deterministic missing-fixture/timeout outcomes; `ExecutionTargetLifecycleService` as sole state/record committer with the lifecycle-service-only `Qualifying` reservation and the `effectiveAvailability` truth table; Retire/Maintenance priority over expiry and over an in-flight qualification; the audit/replay/shutdown matrix; and the eight closed local violation codes. ADH-2026-063 adds `VS0-CF-F16-123..126`: the collection-create bounded body read returns `REQUEST_TOO_LARGE`/400 first for an oversized body, otherwise phase one extracts only the two allowed UID references without classifying/canonicalizing/digesting/reserving a body, authorization and safe access precede the single strict classification of the retained same bytes, and unauthorized, inaccessible-backing, and authorized-classification outcomes each reuse their existing Problem outcome with no extra classification/mutation/idempotency effect. ADH-2026-065 makes the classification and phase-one reference cases exact: `VS0-CF-F16-125` is malformed JSON only (existing 400 `MALFORMED_REQUEST`) and the added `VS0-CF-F16-127` is a duplicate top-level member only (existing 400 `DUPLICATE_FIELD`) — distinct exact non-family classification cases — while the added `VS0-CF-F16-128` fails closed when phase one cannot extract exactly one syntactically usable UID at both `spec.cloudProviderParticipationRef.uid` and `spec.infrastructureStackRef.uid`, returning the existing audited 403 `AUTHORIZATION_DENIED` before derived-scope authorization and safe backing access with no body/backing disclosure and no strict classification, canonicalization, digest, reservation, observer invocation, mutation, publication, or completion; `VS0-CF-F16-123`, `VS0-CF-F16-124`, and `VS0-CF-F16-126` remain unchanged. No shared or downstream case (`VS0-CF-F10`, `VS0-CF-X03`) is counted as F0016-local proof. |

Each conformance test specifies exact expected state, expected error (code + HTTP + type), and expected side effects.

---

## 7. Kiro Anti-Drift Contract

When implementing any Slice 0 feature stage, Kiro must:

1. Load: canonical model, canonical catalog, VS-000 charter, this specification, registry YAML, Slice 0 traceability matrix, FEATURE traceability and DECISION_INDEX.
2. Implement one feature at a time, one stage at a time.
3. Never bypass feature-stage approval gates.
4. Never use stale/prohibited concepts (see §9).
5. Never invent schema fields, state transitions, or error codes not in the registry.
6. Any change to schema, state machine, or error contract must update registry, traceability, and conformance tests in the same change.
7. If a genuine semantic gap is encountered (question cannot be answered by loaded authorities), stop and emit `ARCHITECTURE_DECISION_REQUIRED` with the gap description. Do not guess.

Violation of any rule above constitutes architecture drift and blocks merge.

---

## 8. Repository Traceability and Validation

### Traceability Paths

```text
docs/architecture/vertical-slices/VS-000-contract-specification.md  (this file)
docs/architecture/vertical-slices/VS-000-contract-registry.yaml     (machine-readable registry)
docs/architecture/vertical-slices/VS-000-core-skeleton.md            (skeleton design)
docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md             (Slice 0 contract linkage)
docs/traceability/FEATURE_TRACEABILITY_MATRIX.md                     (feature linkage)
docs/traceability/DECISION_TRACEABILITY_MATRIX.md                    (decision linkage)
docs/decisions/DECISION_INDEX.md                                     (accepted decisions)
```

### Validation Target

```bash
make vs000-contract-check
```

This target validates:
- Registry YAML parses without error.
- All `VS0-SCHEMA-*` entries have the required profile attributes.
- All `VS0-STATE-*` transitions reference declared states.
- Every Problem code matches the FEATURE-0012 closed set and HTTP/type mapping.
- Every `VS0-F01..F20` maps to exactly one `VS0-CF-F*` case.
- No prohibited terms appear in active registry entries.
- No field exists in implementation that is absent from registry.

---

## 9. Explicit Non-Goals

This specification does not:

- Define new resource kinds or API routes.
- Serve as requirements, design, or tasks for any FEATURE.
- Specify future slices (VS-001+).
- Replace the canonical model or catalog as semantic owner.
- Define runtime behavior beyond synthetic Slice 0 scope.
- Mandate persistent storage, Kubernetes CRDs, or real provisioning.
- Prescribe UI, billing, marketplace, or multi-cluster behavior.

### Prohibited Terms (stale/superseded — must not appear in active Slice 0 artifacts)

| Term | Superseded By | Reference |
|------|---------------|-----------|
| ResourcePool | (removed) | DEC-0042 |
| ProviderCapability | (removed) | DEC-0042 |
| Generic Provider (combined) | CloudPlatform + CloudProvider | DEC-0037 |
| ServiceClass (as canonical catalog) | ServiceTypeDefinition + ServiceOffering | DEC-0049 |
| EffectivePolicyContext | EffectiveGovernanceContext | DEC-0050 |
| Six-scope vocabulary | Seven canonical scopes | DEC-0037 |
| SovrunnInstallation (active Slice 0) | CloudProviderParticipation boundary proof | No Phase 2R owner; DEC-0053 deferred |
| CanonicalMigrationPlan, CanonicalMigrationRecord, migration controller | Direct canonical resource creation (FEATURE-0015) | DEC-0059 (supersedes DEC-0058); ADH-2026-045 |

---

## 10. Architecture Decision Required Rule

Ordinary mechanics (CRUD, state transitions, validation, projection, error mapping) must be exact and fully specified in the registry. Implementation must not require interpretation.

When a genuine semantic gap exists — a question that cannot be resolved by the canonical model, catalog, accepted decisions, or this specification — the following applies:

1. The gap is recorded in registry `architectureDecisionRequired` with a stable ID, question and affected features.
2. Downstream code generation, test generation, and implementation for the affected scope stops.
3. An `ARCHITECTURE_DECISION_REQUIRED` marker is emitted.
4. Resolution requires a new DEC, RFC, or ADH through the Architecture Operating System.
5. After resolution, the registry entry is updated, conformance tests are added, and implementation may proceed.

No agent may fill a semantic gap with assumptions, defaults, or analogies from other platforms.

---

*End of specification.*
