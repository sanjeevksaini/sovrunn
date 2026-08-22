# Slice 0 Contract Anti-Drift Steering

Mandatory Kiro steering for FEATURE-0015 through FEATURE-0026 stages contributing to VS-000.

## 1. Context Load Order

Load in this exact sequence before any VS-000 contributing stage:

```text
1. AGENTS.md
2. docs/context/CURRENT_ARCHITECTURE_BASELINE.md
3. docs/architecture/canonical/sovrunn-finalized-data-model.md
4. docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md
5. docs/phase2/PHASE2R_REBASELINE.md
6. docs/phase2/PHASE2_FEATURE_SEQUENCE.md
7. docs/architecture/vertical-slices/VS-000-core-skeleton.md
8. docs/architecture/vertical-slices/VS-000-contract-specification.md
9. docs/architecture/vertical-slices/VS-000-contract-registry.yaml
10. docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md
11. Current FEATURE-xxxx file (requirements, design, tasks)
12. docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md (when working on FEATURE-0015)
13. FEATURE-0012 and FEATURE-0013 authorities (when schemas/errors/decisions consumed)
```

## 2. Authority Precedence

```text
active baseline > canonical model/catalog > accepted DEC/ADH > VS-000 charter > contract spec > registry YAML > approved feature stage > implementation
```

Conflict stop: if any loaded authority contradicts a higher-precedence authority, emit `ARCHITECTURE_DECISION_REQUIRED` and halt. Do not resolve by interpretation.

## 3. Slice0Only Boundary

- Registry narrows canonical semantics; it cannot broaden them.
- No future slice (VS-001+) kinds, routes, profiles, or state machines.
- No PostgreSQL provisioning, OpenShift integration, native network objects, provider SDK calls, or real external effects.
- All plugin execution is synthetic (fake executor, `externalCallCount=0`).
- Seven canonical scope kinds only: Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider.

## 4. One Feature / One Stage

- Implement exactly one feature, one stage at a time.
- Never generate FEATURE-0015 or any later stage without explicit gate approval for the preceding stage.
- VS-000 charter and registry define cross-feature acceptance criteria; they never substitute for the current feature's requirements, design, or tasks.
- Feature gate (`make ff-feature-gate`) must pass before advancing.

## 5. Schema Anti-Drift

- Every resource kind or DTO used maps to exactly one `VS0-SCHEMA-<NNN>` registry entry.
- No field, type, enum value, default, limit, reference, scope, or projection change outside the registry.
- FEATURE-0012 and FEATURE-0013 schemas are reused by reference only; no redefinition, no parallel envelope.
- No unregistered resource kind, API route, or decision profile.
- Seven scope kinds only; no additional scopes.
- A resource `owner` must match `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`; dependencies never transfer ownership.
- `fieldOwnership.introducedBy` and `fieldOwnership.activatedBy` override the containing resource owner for that field.
- A field owned by a later feature must not be initialized, persisted, defaulted, validated, or exposed by an earlier feature.
- ExecutionTarget `status.*` is introduced and activated only by FEATURE-0016.
- ADH-2026-058 atomically replaces the FEATURE-0016 `VS0-SCHEMA-015..017`/`VS0-STATE-004` placeholders: ExecutionTarget create accepts exactly `metadata.name`, `spec.cloudProviderParticipationRef.uid`, `spec.infrastructureStackRef.uid`, and immutable `spec.targetClass=synthetic-iaas`; scope derives solely from the participation reference with no separate persisted/projected `scopeRef`; there is no active persisted `status.availability`, no `Draining`, and no publicly projected `Qualifying` — `Qualifying` is a lifecycle-service-only in-flight reservation and `effectiveAvailability` (Available/Unavailable/Maintenance) is a response-only projection with no independent writer and no ETag effect.
- FEATURE-0016 exposes exactly five Go 1.22 `http.ServeMux` registrations (collection POST/GET, item GET, qualify POST, retire POST) behind one pre-ServeMux transport-only method/path guard that supplies every unmatched-method/trailing-slash outcome with no Problem body, authentication, audit, idempotency, or lifecycle effect; no PATCH, PUT, DELETE, HEAD, watch, filter, pagination, fact, result, or adapter-selection route exists (ADH-2026-058 clause 3).
- FEATURE-0016's sole observer is `sovrunn.synthetic-iaas-observer/v1`: target-bound, clock-driven, external-effect-free; no real adapter, credential, external call, Crossplane dependency, placement, provisioning, plugin execution, or customer IaaS exposure is active (ADH-2026-058 clauses 6,9,10).
- `ExecutionTargetLifecycleService` is the sole committer of ExecutionTarget status and internal `NormalizedTargetFactSet`/`TargetQualificationResult` records; idempotency applies only to ExecutionTarget create, qualify, and retire; GET/LIST never touch idempotency (ADH-2026-058 clauses 5,7,8).
- ADH-2026-060: at ExecutionTarget qualification commit, the Maintenance-race outcomes are mutually exclusive and ordered — an active current-Maintenance marker always selects `CONFLICT`/409 + `VS0_TARGET_MAINTENANCE` first; only when no active marker exists does a changed captured maintenance epoch select `STALE_RESOURCE_VERSION`/412 + `VS0_TARGET_EPOCH_STALE`; a changed referenced InfrastructureStack generation remains the independent epoch-stale case. `VS0-CF-F16-75` is the inactive-marker epoch-stale case; `VS0-CF-F16-89` is the active-marker Maintenance-wins case.
- ADH-2026-061: every successful ExecutionTarget Maintenance entry clears both current `NormalizedTargetFactSet`/`TargetQualificationResult` links and persists `Active/Unqualified`, whether or not a qualification is in flight — it never preserves the last completed conclusion. If a qualification is in flight, Maintenance entry additionally aborts only that qualification and its linked idempotency reservation and wakes only its waiters (never all reservations for the target); the qualifying caller keeps the existing `CONFLICT`/409 + `VS0_TARGET_MAINTENANCE` outcome with no qualification result, completion, or qualification AuditEvent. If no qualification is in flight, no idempotency reservation is cancelled, and the Maintenance entry is still audited. `VS0-CF-F16-42`, `VS0-CF-F16-43`, and `VS0-CF-F16-89` keep their existing exact observable semantics unchanged.
- ADH-2026-063: for ExecutionTarget collection create (`POST /apis/execution.sovrunn.io/v1alpha1/execution-targets`), the safe precedence invariant is ordered — the bounded single body read returns the existing `REQUEST_TOO_LARGE`/400 first when the body exceeds the existing limit, before phase-one extraction, authorization, or safe access; otherwise the exact bytes are retained and phase one extracts only `spec.cloudProviderParticipationRef.uid` and `spec.infrastructureStackRef.uid`, emitting no body-classification Problem; derived-scope `executiontarget.write` authorization and safe backing access precede any body classification; only then are the retained same bytes strictly classified exactly once, then graph/reference validation, idempotency, audit, and publication. Phase one never classifies, canonicalizes, digests, or reserves a body; invalid, unknown, and duplicate bodies are never digested or reserved. Action `If-Match` header-form precedence is preserved. `VS0-CF-F16-123..126` are the exact local proof rows.
- ADH-2026-065: for ExecutionTarget collection create, the create classification and phase-one reference conformance cases are made exact without altering ADH-2026-063 precedence. `VS0-CF-F16-125` is the exact non-family malformed-JSON classification case returning the existing `MALFORMED_REQUEST`/400 (as `VS0-CF-F16-51`); the new `VS0-CF-F16-127` is the exact non-family duplicate-top-level-member classification case returning the existing `DUPLICATE_FIELD`/400 (as `VS0-CF-F16-52`); `VS0-CF-F16-125` and `VS0-CF-F16-127` are distinct exact non-family classification cases, not a mixed equivalence family. The new `VS0-CF-F16-128` is the fail-closed phase-one extraction case: when the existing bounded body read succeeds but phase one cannot extract exactly one syntactically usable UID at both `spec.cloudProviderParticipationRef.uid` and `spec.infrastructureStackRef.uid`, the create returns the existing audited `AUTHORIZATION_DENIED`/403 before derived-scope authorization and safe backing access, discloses no body-classification detail and no backing-resource existence, and performs no strict classification, canonicalization, digest, idempotency reservation, observer invocation, mutation, publication, or completion, adding no violation or top-level Problem code. `VS0-CF-F16-123`, `VS0-CF-F16-124`, and `VS0-CF-F16-126` remain unchanged; the FEATURE-0016 local conformance range extends to `VS0-CF-F16-128`. No route, field, state, writer, dependency, top-level Problem code, or external effect is added.
- ServiceRegion and all its fields are introduced and activated only by FEATURE-0022; its dependency on FEATURE-0016 does not transfer activation to FEATURE-0016.
- CloudProviderParticipation `spec.providerSelectionModes` is introduced and activated only by FEATURE-0021.
- VS0-SCHEMA-057 is a permanent retired tombstone; SovrunnInstallation is not an active Slice 0 resource and ID 057 must never be reused.

## 6. Writer Anti-Drift

- Every mutable path maps to exactly one `VS0-WRITER-<NNN>` registry entry.
- Clients never write `status`.
- No competing condition producers for the same condition type.
- No plan/decision/execution authority collapse (planner ≠ decision-service ≠ executor).
- No secret values in any writer path; `SecretRef` identifier only.
- Conflicts fail closed per registered conflict code.
- CanonicalMigrationPlan/Record, the migration-controller, and the approved-migration-plan-publisher are retired (DEC-0059 supersedes DEC-0058; ADH-2026-045 canonical bootstrap). VS0-WRITER-020 and VS0-WRITER-021 are permanently retired tombstones and must never be reused or reintroduced as active writers.

## 7. State Anti-Drift

- Exact `VS0-STATE-<NNN>` transitions and guards are authoritative.
- Milestones are observable checkpoints, not lifecycle phases.
- DecisionRecord is always persisted as `FINAL`; no intermediate record states.
- Operation owns async lifecycle: `Pending → Running → Succeeded | Failed | Cancelled`.
- Immutable records are append-only; corrections create linked records, never mutate.
- Invalid transitions are rejected with the registered error/violation.
- VS0-STATE-011 (formerly the CanonicalMigrationRecord milestone sequence) is retired (DEC-0059 supersedes DEC-0058; ADH-2026-045 canonical bootstrap) and must never be reused or reintroduced.

## 8. Error Anti-Drift

- Only exact FEATURE-0012 top-level codes and their HTTP/URN mappings are permitted.
- Slice 0 violation codes appear only in `violations[].code`, never as top-level Problem codes.
- HTTP status codes are not Problem `code` values; they are transport metadata.
- Safe denial: error responses never confirm existence of inaccessible resources.
- No `429` or quota-specific HTTP code; quota exhaustion uses `CONFLICT`/409 + `VS0_QUOTA_EXHAUSTED`.

## 9. Conformance and Traceability

- Every requirement maps to: ADH/DEC authority + owner + schema/writer/state/error IDs + `VS0-CF-<ID>`.
- Failure mappings `VS0-F01..F20` have one-to-one conformance cases `VS0-CF-F01..F20`.
- Migration conformance (`VS0-MIG-F01..F03`, `VS0-CF-MIG01..MIG02`, `VS0-CF-MIGF01..MIGF03`) is retired (DEC-0059 supersedes DEC-0058; ADH-2026-045 canonical bootstrap) and must never be reused or reintroduced.
- FEATURE-0015 local conformance uses `VS0-CF-F15-01..F15-41` and `VS0-CF-X03`.
- FEATURE-0015 status writes for `CloudPlatform`, `CloudProvider`, `HostingLocation`, `Datacenter`, `FaultDomain`, `InfrastructureStack`, and `CloudProviderParticipation` resolve to `api-server` as the sole status writer; no separate topology/stack/provider/participation controller exists (ADH-2026-046 decision 1).
- `CloudProviderParticipation` create requires `Idempotency-Key` only (no `If-Match`); every action on an existing participation requires both `If-Match` and `Idempotency-Key` (ADH-2026-046 decision 2).
- FEATURE-0015 traceability must NOT use downstream-owned conformance (HP01, F09) as local acceptance evidence; those are integration references only.
- Every F0015 collection-create kind has a closed client-required/client-optional field boundary; the client must never supply `metadata.scopeRef`, `status.*`, identity/version/timestamp fields, a field owned by another feature, or an unknown field; every scope is server-derived from exactly one registered source per kind (Platform-root, participation-from-`cloudPlatformRef`, HostingLocation-from-`topology.write` grant, descendant-from-parent-reference); a successful create returns `201` and a successful PATCH returns `200` (ADH-2026-047 decisions 1-2).
- Topology `metadata.name` is immutable identity; only `spec.description` is PATCHable for `HostingLocation`, `Datacenter`, `FaultDomain`, and `InfrastructureStack` ("Description is PATCHable; name is immutable identity" corrects ADH-2026-045's stale "Name and description are PATCHable" wording; ADH-2026-047 decision 3).
- `CloudProviderParticipation` collection create accepts exactly `metadata.name`, `spec.cloudPlatformRef`, `spec.cloudProviderRef`, and `spec.environment` (`development` only); `spec.providerSelectionModes` and `spec.permittedHostingLocationRefs` are FEATURE-0021-introduced and FEATURE-0021-activated and must never be accepted, stored, defaulted, validated, or exposed by FEATURE-0015 (ADH-2026-047 decision 4).
- Before a requirements-stage Kiro invocation, the FEATURE-0015 readiness gate must run a deterministic, fail-closed contract-executability audit (field classification, route-contract completeness, single-scope-derivation-source, stale-wording/absent-create-field/client-supplied-scope/unowned-participation-field/missing-create-status rejection) as part of `make feature-0015-architecture-readiness` (ADH-2026-047 decision 5).
- For FEATURE-0015, an ISO-3166-1 alpha-2 value means an assigned code from a fixed, repository-owned, version-pinned static dataset (as of 2026-08-12), not merely a syntactically valid code; a syntactically valid but unassigned code is rejected with `VALIDATION_FAILED`/422; `administrativeAreaCode` remains syntax-plus-country-prefix validation only, with no ISO-3166-2 membership dataset (ADH-2026-048 decision 1).
- A missing/empty/malformed/over-length `Idempotency-Key` on a required route, an `If-Match` supplied on a `CloudProviderParticipation` collection create, and a non-empty JSON body on a participation item action each return `MALFORMED_REQUEST`/400 using only existing FEATURE-0012 codes; a missing/malformed/non-current `If-Match` on an existing-participation action remains `STALE_RESOURCE_VERSION`/412; none of these produce a mutation, idempotency record, or AuditEvent beyond an already-approved audited denial (ADH-2026-048 decision 2).
- The FEATURE-0015 readiness gate must additionally prove, for every owned route, the complete pipeline (body/headers → authn/authz → root/reference/safe-denial ordering → structural/semantic validation → exact Problem outcome → state/status result → audit effect → idempotency/race outcome → local conformance test), failing closed on an unspecified observable input outcome, missing per-route mapping, unassigned ISO code accepted, non-empty participation-action body accepted, or a route-count/pattern mismatch against the explicit 22-route design (ADH-2026-048 decision 4).
- For FEATURE-0015, a failure to append the FEATURE-0013 AuditEvent required for a mutation returns the inherited `INTERNAL_ERROR`/500 Problem Details response; the mutation and its idempotency completion are not published; `DEPENDENCY_UNAVAILABLE`/503 must never be used for this outcome (ADH-2026-049).
- `VS0-CF-F15-18` (Idempotency-Key replay) and `VS0-CF-F15-19` (Idempotency-Key reuse with changed digest) must each explicitly cover all seven FEATURE-0015 collection creates (`CloudPlatform`, `CloudProvider`, `CloudProviderParticipation`, `HostingLocation`, `Datacenter`, `FaultDomain`, `InfrastructureStack`) and all eight existing-participation create/action routes (`accept`, `reject`, `withdraw`, `suspend`, `resume`, `request-release`, `accept-release`, `decline-release`); their existing IDs, expected states, `CONFLICT`, and `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH` behavior must never be narrowed or reinterpreted (ADH-2026-050).
- ADH-2026-054 closes forged-grant and replay isolation: only `X-Sovrunn-Bootstrap-Grant` and valid duplicate-free top-level `bootstrapGrant` are audited forged-grant carriers with the exact 401/403/400 precedence; item-action idempotency binds the concrete `{uid}` target and rechecks current authorization/safe access before replay; only successful collection creates/actions complete records, while validation, uniqueness, lifecycle, stale-version, audit, panic, and shutdown abort and wake waiters.
- FEATURE-0015 produces exactly one redacted FEATURE-0013 AuditEvent, after authentication/current authorization are evaluated, for: a successful resource collection create or PATCH; a successful participation create/action or deterministic scheduler expiry; an authenticated client attempt to write api-server-owned status or identity/metadata; a CloudPlatform-root creation denial; an authenticated collection LIST without a read grant; an authenticated header/body forged-grant attempt; and an authenticated inaccessible cross-provider/safe-denial reference. FEATURE-0015 produces no durable AuditEvent for missing/invalid authentication, malformed/prohibited body/header/field input, unsupported method/media type, stale `If-Match`, same-key replay, changed-digest idempotency conflict, pair-uniqueness conflict, or an invalid participation source state (ADH-2026-051).
- FEATURE-0015 owns exactly 22 logical endpoint paths (seven collection, seven item, eight `CloudProviderParticipation` action) and registers exactly 35 explicit Go 1.22 `http.ServeMux` method/path patterns (14 collection LIST/create + 12 PATCHable-item GET/PATCH + 1 participation-item GET-only + 8 participation actions); no path-only or catch-all handler registration, reflection, or internal HTTP-method dispatch is permitted. The declared complete `{uid}` segment remains a permitted ServeMux wildcard (ADH-2026-051,053).
- The eight `CloudProviderParticipation` action registrations use `/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/<action>`, where `{uid}` is a complete Go 1.22 `http.ServeMux` wildcard segment read with `r.PathValue("uid")`. The eight `<action>` values are `accept`, `reject`, `withdraw`, `suspend`, `resume`, `request-release`, `accept-release`, and `decline-release`. The retired `/{uid}:<action>` form is not an allowed FEATURE-0015 route form (ADH-2026-053).
- `VS0-CF-F15-31` proves that missing or invalid authentication on any F0015 owned method/path pattern returns `AUTH_REQUIRED`/401 with no mutation, lifecycle/status write, idempotency record, or AuditEvent; it does not replace or modify inherited FEATURE-0018 `VS0-CF-F01` and introduces no new top-level Problem or violation code (ADH-2026-051).
- `VS0-CF-F15-32` is the exact local proof that an authorized CloudPlatform or CloudProvider collection create whose `metadata.name` duplicates an existing resource of the same kind in that kind's server-derived Platform-root scope returns `ALREADY_EXISTS`/409 with no duplicate resource persisted; `VS0-CF-F15-33` is the exact local proof that an authorized FEATURE-0015 collection create violating a registry-declared semantic schema constraint in `VS0-SCHEMA-008..014` (other than a duplicate name, an unassigned ISO-3166-1 alpha-2 code, or a malformed/prohibited header/body condition) returns `VALIDATION_FAILED`/422 with no resource persisted, no idempotency completion, and no AuditEvent unless an already-approved audited-denial rule independently applies. `VS0-CF-F15-01`, `02`, `12`, `13`, `14`, and `23` must never be represented as proof of name-duplicate rejection or general schema-constraint failure; `VS0-CF-F15-21` must never be falsely assigned to an unrelated REQ or AC identifier (ADH-2026-052).
- ADH-2026-056: every mutable-resource PATCH requires current `If-Match` and returns `STALE_RESOURCE_VERSION`/412 if missing, malformed, or stale; exact read actions are `cloudplatform.read`, `cloudprovider.read`, `topology.read`, and `participation.read` at their registered scopes; replay returns stored successful status/body plus `Content-Type: application/json` only, with fresh request correlation and effective retention limited by both 24 hours and process lifetime. `VS0-CF-F15-38..40` are the exact local proof cases.
- ADH-2026-057: a wrong-administrator mutable-spec PATCH is denied `AUTHORIZATION_DENIED` before semantic PATCH processing or publication with no mutation, lifecycle/status update, AuditEvent, or idempotency record. `VS0-CF-F15-41` is its sole exact local proof case.
- Registry, traceability matrix, and test contract update atomically in the same change.
- ADH-2026-055: HEAD returns 405 despite ServeMux GET matching; the forged header precedes media/body validation; completed idempotency retains for 24 hours with deterministic 10,000-record completed-only eviction; RFC 7396 PATCH and ascending-UID LIST are mandatory.
- Run `make vs000-contract-check` and `make phase2r-drift-check` before marking complete.
- Run `make feature-0015-architecture-readiness`, `make feature-contract-check FEATURE=FEATURE-0015`, and `make feature-0015-formal-check` before FEATURE-0015 requirements generation.

## 10. Prohibited Active Concepts

Never use in active Slice 0 artifacts:

```text
ResourcePool
ProviderCapability
generic Provider as combined owner/operator
ServiceClass as canonical catalog
EffectivePolicyContext
provider-neutral (as adjective)
six-scope authority
SovrunnInstallation as active Slice 0 resource
```

## 11. Stop Conditions

Emit `ARCHITECTURE_DECISION_REQUIRED` and halt when:

- A semantic choice is missing from loaded authorities.
- Conflicting owners exist for the same path or condition.
- A schema, state transition, or error code is required but unregistered.
- A requirement cannot be mapped to existing registry entries.

Kiro may resolve ordinary mechanics (CRUD wiring, validation plumbing, projection assembly) already explicit in the registry. Kiro must never invent semantics.

## 12. Generated Output Contract

Every stage output (requirements, design, or tasks) must include:

- **Architecture Traceability section**: lists consumed ADH/DEC/RFC, schema IDs, writer IDs, state IDs, error codes.
- **Conformance Mapping section**: maps each acceptance criterion to `VS0-CF-<ID>`.
- No orphan IDs (every referenced ID must exist in registry).
- IDs are never renumbered or reused across stages.

## 13. Implementation Exclusion

This steering governs spec/design/task generation only. No Go code or runtime test implementation derives from this steering alone. Implementation proceeds only from approved feature tasks.
