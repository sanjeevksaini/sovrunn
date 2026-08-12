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
- FEATURE-0015 local conformance uses `VS0-CF-F15-01..F15-30` and `VS0-CF-X03`.
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
- Registry, traceability matrix, and test contract update atomically in the same change.
- Run `make vs000-contract-check` and `make phase2r-drift-check` before marking complete.
- Run `make feature-0015-architecture-readiness` before FEATURE-0015 requirements generation.

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
