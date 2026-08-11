# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-045
- Date: 2026-08-11
- Source discussion: Phase 2R architecture completeness audit
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Replace alpha runtime migration with a canonical bootstrap

## Summary

FEATURE-0001 through FEATURE-0014 remain retained repository assets and create no active runtime data-migration obligation. Sovrunn has no live control plane, customer data, production API estate, or persisted alpha state to convert. The first control-plane release therefore starts directly with canonical contracts. Existing standards, generic infrastructure, and compatible implementations must be reused or extended under the FEATURE-0011 reuse assessment; only obsolete alpha domain kinds and APIs are not carried forward as active compatibility surfaces. Phase 2R does not introduce a `CanonicalMigrationPlan`, `CanonicalMigrationRecord`, migration controller, runtime converter, or cutover state machine.

## Classification

Replacement.

## Existing approved baseline

DEC-0058 and ADH-2026-041 define alpha migration as a signed cutover with no dual write/authority and immutable history preservation. ADH-2026-043 and ADH-2026-044 add FEATURE-0015 migration closure and run-cardinality semantics on that premise.

That premise is no longer appropriate: the completed alpha features are repository history rather than a live deployment with state that must be migrated. `ARCH-2026.08-PHASE2R-CANONICAL` otherwise remains the active baseline.

## Decision or proposed decision

1. Replace the DEC-0058/ADH-2026-041 runtime alpha-migration approach with a canonical bootstrap.
2. FEATURE-0001 through FEATURE-0014 remain retained repository assets. Their code, documents, tests, standards, and compatible generic infrastructure must be assessed for reuse or extension under FEATURE-0011; they do not constitute live control-plane data requiring conversion.
3. The first executable Sovrunn control plane exposes canonical resource contracts only. Alpha kinds, scope vocabulary, and desired-state endpoints are not active compatibility surfaces and must not be reintroduced as runtime APIs, writers, projections, or persistence contracts.
4. FEATURE-0015 owns direct creation of its canonical cloud-model resources. It does not implement alpha import, fixture conversion as a product capability, migration plans/records, migration controller authority, migration milestones, write-freeze/cutover behavior, backup/restore migration evidence, or global migration completion.
5. If historical fixture comparison remains useful, it is a bounded developer/test utility only. It is neither a platform resource nor a runtime protocol, accepts no user data, has no active API, and creates no migration authority.
6. The one-authority principle remains mandatory: it is achieved by beginning with canonical desired-state contracts, not by operating alpha and canonical writers during a cutover.
7. A new accepted DEC must supersede DEC-0058; ADH-2026-041, ADH-2026-043 migration decisions, and ADH-2026-044 are superseded only in their migration-specific portions. Their unrelated canonical-model, ownership, scope, audit, and anti-drift decisions remain intact.
8. An Active `CloudProviderParticipation` is necessary but not sufficient for customer visibility. The CloudPlatform decides which enrolled customer Organizations may see or select an Active participation. Individual-user visibility derives from that Organization-level eligibility and scoped authorization; it is not stored directly on `CloudProviderParticipation`. Pending, Rejected, Withdrawn, Expired, Suspended, Terminating, and Terminated participations are never end-user visible.
9. `ExecutionTarget` moves in its entirety to FEATURE-0016. FEATURE-0015 ends at `InfrastructureStack` and contains no ExecutionTarget identity, schema, route, status, writer, conformance, or target lifecycle behavior. FEATURE-0016 owns target registration, identity, status, adapter facts, qualification, availability, maintenance, and the target state machine.
10. `CloudProviderParticipation` retains `Suspended` as a derived effective lifecycle state. Its controller-owned status carries two durable, independently controlled holds: `platformSuspended` and `providerSuspended`. A CloudPlatform administrator may set or clear only `platformSuspended` through explicit suspend/resume actions; a CloudProvider administrator may set or clear only `providerSuspended` through its explicit suspend/resume actions. A participation is effectively `Active` only when it has been accepted and both holds are false; it is effectively `Suspended` when either hold is true. Clearing one hold must not reactivate a participation while the other remains true. Its terminal request states are Rejected, Withdrawn, and Expired. FEATURE-0015 exposes no mutable participation spec fields; lifecycle changes occur only through explicit actions. Each hold-changing action produces correlated AuditEvent evidence naming the acting party and resulting effective state.
11. FEATURE-0015 mutable resource updates use `PATCH` with `application/merge-patch+json`. `PUT` and `DELETE` are not exposed by FEATURE-0015. Identity, scope, relationship references, and the embedded owner registration remain immutable. Because CloudProviderParticipation has no mutable F0015 spec fields, it exposes lifecycle actions rather than PATCH.
12. FEATURE-0015 uses a single-authority onboarding model for CloudPlatform participation. The CloudPlatform's control plane is authoritative for its CloudProvider, topology, and CloudProviderParticipation records. A provider may operate Sovrunn software in its own sites, but that installation is not a federated control plane, a second authoritative canonical store, or an F0015 synchronization protocol. A provider administrator registers the provider and its topology in the CloudPlatform-authoritative control plane, then requests participation there; the CloudPlatform administrator accepts or otherwise acts through the participation protocol. FEATURE-0015 introduces no remote deployment resource, signed remote offer, cross-control-plane replication, or federation trust protocol.

## FEATURE-0015 architecture closure decisions

The following decisions are part of this approved handoff. Kiro must carry them verbatim in meaning into the controlling contract authorities; requirements, design, and tasks must not select alternatives.

### Canonical resources, scope, and update surface

| Resource or concern | Approved decision |
|---|---|
| Scope and deployment | `CloudPlatform` and `CloudProvider` are Platform-scoped canonical resources. `CloudProviderParticipation` is CloudPlatform-scoped. `HostingLocation`, `Datacenter`, `FaultDomain`, and `InfrastructureStack` are CloudProvider-scoped. CloudPlatform has no Organization scope or Organization reference. A Sovrunn deployment is the implicit runtime/control-plane boundary, not an F0015 resource or reference; all F0015 records reside in its single authoritative control plane. A CloudPlatform may use provider sites in many locations, but does not span independent authoritative control-plane deployments in F0015. Every F0015 Platform-scoped resource receives an immutable server-assigned `metadata.scopeRef` to the deployment's canonical Platform root (`core.sovrunn.io/v1alpha1`, kind `Platform`, name `sovrunn`, UID-pinned). Clients neither supply nor PATCH this reference. |
| CloudPlatform | `metadata.name` is the unique canonical name and is not duplicated in `spec`. `spec.description` is optional and PATCHable. Immutable `spec.ownerRegistration` contains exactly `legalName`, `registrationIdentifier`, and `jurisdictionCode`; it contains no user, Organization, role, or grant reference. Multiple CloudPlatforms may have the same owner registration. |
| CloudProvider | `metadata.name` is unique in the authoritative control plane. `spec.displayName` and non-empty `spec.operatingMarkets[]` are PATCHable; operating markets are unique, uppercase ISO 3166-1 alpha-2 country codes. It has no deployment reference, credential, or CloudPlatform reference. |
| Topology | Each topology resource has immutable parent TypedRef where applicable: Datacenter→HostingLocation, FaultDomain→Datacenter, InfrastructureStack→FaultDomain. `HostingLocation.spec` requires uppercase ISO 3166-1 alpha-2 `countryCode` and non-empty `locality`; optional `administrativeAreaCode`, if present, is ISO 3166-2. Name and description are PATCHable; parent references, scope, identity, and metadata UID are immutable. Resource names are unique within kind and authoritative scope. |
| Creation order | CloudPlatform exists first. A CloudProvider, HostingLocation, Datacenter, FaultDomain, InfrastructureStack, or CloudProviderParticipation create is denied until at least one CloudPlatform exists in the authoritative deployment control plane, returning CONFLICT (409) with `VS0_CLOUDPLATFORM_ROOT_REQUIRED`. CloudProvider exists before its topology and before participation. Provider topology may be registered before participation is Active, but it is not customer-visible merely because it exists. |
| Status and deletion | The API server creates CloudPlatform, CloudProvider, and topology resources with `status.phase=Active` and is their sole F0015 status writer. F0015 exposes no lifecycle action, decommission, or deletion for those resources. A client cannot PATCH `status`. |
| Update protocol | PATCH is only `application/merge-patch+json`. PATCH is available only for CloudPlatform, CloudProvider, and topology. `PUT` and `DELETE` return 405. Participation has no PATCHable spec fields. |

### Participation lifecycle and independent suspension

`CloudProviderParticipation.status` is controller-owned and persists only current lifecycle facts: server-computed/stored `phase`, `platformSuspended`, `providerSuspended`, retained `requestExpiresAt`, and the standard system status fields. It must not duplicate transition history, which is durable FEATURE-0013 AuditEvent evidence. Both suspension holds initialize to false. Holds are source facts: for accepted participation, the server stores `phase=Active` exactly when both holds are false and `phase=Suspended` when either hold is true. Pending, Terminating, and every terminal phase override the hold-derived phase; the holds remain retained but cannot be changed outside Active or Suspended.

| Action or event | Authorized actor/action | From → to | Rule |
|---|---|---|---|
| Create participation | CloudProvider administrator / `participation.request` | absent → Pending | The create request is the provider's affirmative request. Set `requestExpiresAt = createdAt + 7 days`. |
| Accept | CloudPlatform administrator / `participation.accept.platform` | Pending → Active | The platform's action completes the participation agreement. |
| Reject | CloudPlatform administrator / `participation.reject.platform` | Pending → Rejected | Terminal; no relationship becomes active. |
| Withdraw | CloudProvider administrator / `participation.withdraw.provider` | Pending → Withdrawn | Terminal; no relationship becomes active. |
| Expire | deterministic system scheduler | Pending → Expired | Occurs at `requestExpiresAt`; terminal and idempotent. |
| Suspend | CloudPlatform or CloudProvider administrator / respective suspend grant | Active or Suspended → effective Suspended | Actor sets only its own durable hold true. |
| Resume | CloudPlatform or CloudProvider administrator / respective resume grant | Suspended → Active only when both holds false | Actor clears only its own durable hold. The other hold cannot be overridden. |
| Request release | CloudPlatform administrator / `participation.request-release.platform` | Active or Suspended → Terminating | Holds are retained while release is pending. |
| Accept release | CloudProvider administrator / `participation.accept-release.provider` | Terminating → Terminated | Terminal; the relationship ends without deletion. |
| Decline release | CloudProvider administrator / `participation.decline-release.provider` | Terminating → prior effective Active or Suspended | Returns to Active only if both holds are false; otherwise returns to Suspended. |

Suspend and resume actions are denied with the established participation-state conflict outcome while phase is Pending, Terminating, Rejected, Withdrawn, Expired, or Terminated. Rejected, Withdrawn, Expired, and Terminated are retained terminal records. Exactly one non-terminal participation—Pending, Active, Suspended, or Terminating—may exist for a `(cloudPlatformUID, cloudProviderUID)` pair; a new request after a terminal record receives a new UID. All participation actions have an empty JSON body and require `If-Match` plus `Idempotency-Key`; the same key and request returns the original result before version checking, while a new key with a stale version returns 412.

### Referential integrity and safe denial

The same CloudProvider UID must be preserved through Datacenter→HostingLocation, FaultDomain→Datacenter, and InfrastructureStack→FaultDomain. Reference resolution performs authorization/non-disclosure before structural validation: a caller unable to access a referenced provider resource receives `RESOURCE_NOT_FOUND` (404) with `VS0_AUTHORIZATION_SAFE_DENIAL`; a caller authorized to see both resources but presenting a mismatched provider graph receives `VALIDATION_FAILED` (422). FEATURE-0015 has no ExecutionTarget invariant or route.

### Bootstrap authorization, API, concurrency, and audit

| Concern | Approved decision |
|---|---|
| Bootstrap authorization | Authenticated request context receives server-resolved, deterministic bootstrap grants; grants are never accepted from an HTTP header or request body and F0015 persists no roles, memberships, or assignments. A grant has an action, scopeRef, and optional `resourceUID` restriction. The server-configured bootstrap principal receives `cloudplatform.write` and `cloudprovider.write` at the deployment Platform-root scope to create the first Platform-scoped resources. After creation, the same action may be narrowed to one target resource UID for PATCH. `topology.write` is CloudProvider-UID scoped; participation actions are scoped to their existing CloudPlatform or CloudProvider target. FEATURE-0018 may replace the grant source while preserving action names/scopes. |
| Actions | `cloudplatform.{read,write}`, `cloudprovider.{read,write}`, `topology.{read,write}`, and `participation.{read,request,accept.platform,reject.platform,withdraw.provider,request-release.platform,accept-release.provider,decline-release.provider,suspend.platform,resume.platform,suspend.provider,resume.provider}`. Every grant is scope-pinned to the relevant CloudPlatform or CloudProvider. |
| Routes | Exact collections are `/apis/core.sovrunn.io/v1alpha1/cloud-platforms`, `/apis/core.sovrunn.io/v1alpha1/cloud-providers`, `/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations`, and `/apis/infrastructure.sovrunn.io/v1alpha1/{hosting-locations,datacenters,fault-domains,infrastructure-stacks}`. Every item route appends `/{uid}`. All owned resources expose POST/GET/LIST; only mutable resources expose PATCH. Participation POST actions append `:accept`, `:reject`, `:withdraw`, `:suspend`, `:resume`, `:request-release`, `:accept-release`, or `:decline-release` to its item route. |
| Read/list behavior | LIST returns only resources visible under the caller's scoped read grant and returns `200 []` when none are visible. A caller with no collection-read grant receives authorization denied (403) without resource details. Direct GET or referenced-resource resolution without access returns 404 safe denial. |
| Optimistic concurrency | PATCH and every participation action require the current FEATURE-0012 `resourceVersion` through `If-Match`; a stale value returns the FEATURE-0012 precondition failure (412) without a write. |
| Idempotency | POST create/actions require `Idempotency-Key`. Same authenticated principal, operation route, key, and canonical request digest return the original result. Reuse of the key with a different request returns CONFLICT (409) using the stable F0012 idempotency-mismatch Problem mapping. Idempotency records remain valid for the lifetime of the running Phase 2R API process. A process restart clears the in-memory resource registry and idempotency cache together; F0015 makes no cross-restart replay guarantee. |
| Races | Concurrent creation of the same non-terminal participation pair permits exactly one success; the other returns ALREADY_EXISTS (409). No duplicate participation is stored. |
| Audit | Durable FEATURE-0013 AuditEvent is required for resource create/PATCH, every participation creation/action/expiry, authorization denial, and safe denial. The mutation and required AuditEvent are one atomic outcome: audit failure leaves the resource/lifecycle unchanged and returns the established standard 5xx problem; idempotency replay creates no second AuditEvent. Each event contains actor, action, subject UID, result, timestamp, and request ID; expiry additionally carries system actor and scheduler execution ID. GET/LIST use logs/metrics only. No secret, credential, or protected handle is recorded. |
| Stable errors | Reuse FEATURE-0012 precondition failure for stale `If-Match` (412), `ALREADY_EXISTS`/`VS0_PARTICIPATION_DUPLICATE` for duplicate participation pair (409), and `CONFLICT`/`VS0_PARTICIPATION_STATE_INVALID` for invalid lifecycle actions (409). Preserve 404/`VS0_AUTHORIZATION_SAFE_DENIAL` for inaccessible references. Add `VALIDATION_FAILED` (422) violations `VS0_TOPOLOGY_PROVIDER_MISMATCH` for an authorized cross-provider topology graph and `VS0_PATCH_IMMUTABLE_FIELD` for an immutable-field patch; add `CONFLICT` (409) `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH` for changed reuse of an idempotency key. |

### Required feature-local proof

The revised registry/specification/traceability and the FEATURE-0015 readiness checker must prove every owned route, field mutability rule, topology invariant, lifecycle transition, hold combination, expiry, authorization denial, safe denial, idempotency replay/mismatch, duplicate-pair race, stale PATCH, no-DELETE boundary, and absence of migration and ExecutionTarget contracts. No F0015 acceptance case may rely on a future-feature conformance ID.

## Rationale

Building runtime migration machinery before a control plane exists would add a product surface, lifecycle, security boundary, controller authority, audit stream, and conformance burden for data that does not exist. The Phase 2R goal is to establish the canonical foundation correctly on its first executable path. Canonical bootstrap removes the current migration-record cardinality contradiction rather than encoding speculative technical debt into the platform.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Build
- Summary of mature candidates / applicable standards:
- Reuse the existing FEATURE-0012 metadata, typed-reference, validation, Problem Details, and resource-shape standards; reuse FEATURE-0013 AuditEvent only for material canonical-resource actions that actually occur.
- Reuse compatible existing registry, handler, storage, validation, optimistic-concurrency, idempotency, and test-harness infrastructure where it satisfies the canonical contract; extend it rather than duplicating it when needed.
  - No migration engine, compatibility layer, or data-conversion product capability is needed because there is no live alpha state.
- Why Reuse, Wrap, and Extend are insufficient:
  - Reusing, wrapping, or extending a migration engine would introduce an unnecessary runtime authority and an unneeded compatibility surface.
- Sovrunn-owned responsibility summary:
  - Direct canonical bootstrap, canonical-only resource/API admission, and explicit rejection of obsolete active concepts.
- Non-goals summary:
  - Runtime alpha conversion, dual read/write, compatibility aliases, migration records/plans/controllers, customer-data import, real provider execution, and Phase 3 lifecycle work.

## Phase impact

- Current phase allowed? Yes.
- Current phase boundary impact:
  - Reduces Phase 2R scope by removing speculative migration machinery. It does not move implementation into a future phase or relax any canonical security, ownership, audit, or no-dual-authority invariant.

## Conflict check

- Conflicts with accepted DEC/RFC? Yes.
- Conflicting decisions:
  - DEC-0058
  - ADH-2026-041 migration runtime semantics
  - ADH-2026-043 decisions 1–3 and 6 insofar as they govern migration plan/record behavior
  - ADH-2026-044
- Resolution required:
  - Create an accepted replacement DEC that explicitly supersedes DEC-0058, update the Decision Index and current decision summary, and annotate the affected handoffs with supersession scope. The architecture baseline remains `ARCH-2026.08-PHASE2R-CANONICAL`; it requires a controlled baseline review/update because its active migration description changes.

## Required action

- Update architecture documentation
- Update phase scope and feature sequence documentation
- Create/update DEC
- Update current architecture baseline/context and open-question status
- Update traceability matrix
- Update feature gate/checks
- Regenerate Kiro requirements only after the architecture admission gate passes

## Impacted files

Kiro must validate the complete impact set and update only the necessary controlling files, including:

- `docs/decisions/DECISION_INDEX.md`
- a new accepted DEC superseding DEC-0058
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_DECISION_SUMMARY.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/context/OPEN_QUESTIONS.md`
- `docs/architecture/canonical/sovrunn-finalized-data-model.md`
- `docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md`
- `docs/phase2/PHASE2R_REBASELINE.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/features/FEATURE_INDEX.md`
- `docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md`
- `docs/architecture/vertical-slices/INTEGRATED_DELIVERY_PLAN.md`
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`
- `docs/features/FEATURE-0015-canonical-cloud-model-and-alpha-migration-foundation.md`
- `docs/architecture/FEATURE-0015-canonical-cloud-model-and-alpha-migration-foundation.md`
- `.automation/features/FEATURE-0015.control.json`
- `scripts/feature-0015-architecture-readiness-check.py`
- `docs/reviews/architecture-readiness/PHASE2R_ARCHITECTURE_COMPLETENESS_AUDIT.md`
- `docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md` (annotate as superseded for migration-related closure claims; retain it as historical review evidence)
- `docs/diagrams/structurizr/workspace.dsl` only if it currently represents a migration controller, migration plan/record, or cutover flow.

Kiro must also inspect ADH-2026-041, ADH-2026-043, and ADH-2026-044 and preserve them as historical records with precise supersession annotations rather than silently rewriting history.

## Impacted features

- FEATURE-0015: rename/re-scope to the canonical cloud-model foundation; remove runtime alpha-migration resources, authorities, conformance, and acceptance criteria.
- FEATURE-0016: own ExecutionTarget registration and complete target lifecycle; no deferred migration-run obligation.
- FEATURE-0017 through FEATURE-0025: no deferred migration-run ownership or migration continuation obligations.
- FEATURE-0026: no global migration-completion proof; retain only applicable Slice 0 integration/conformance responsibilities.

## Acceptance criteria for Kiro update

- [ ] A replacement DEC explicitly supersedes DEC-0058 and preserves the canonical-only, one-active-authority invariant.
- [ ] All migration-specific runtime contracts are removed from active Phase 2R authorities: `CanonicalMigrationPlan`, `CanonicalMigrationRecord`, VS0-SCHEMA-060/061, VS0-WRITER-020/021, VS0-STATE-011, migration failure mappings, and VS0-CF-MIG*/migration-specific F15 cases.
- [ ] FEATURE-0015 has no runtime alpha conversion, import, compatibility, cutover, backup/restore migration, or migration-controller scope.
- [ ] FEATURE-0015 ends at InfrastructureStack and has no ExecutionTarget contract; FEATURE-0016 owns the complete ExecutionTarget contract.
- [ ] Participation lifecycle has explicit Rejected, Withdrawn, Expired, Suspended, and resume semantics; no mutable participation spec fields are present in FEATURE-0015.
- [ ] FEATURE-0015 exposes PATCH only as `application/merge-patch+json`; it exposes neither PUT nor DELETE.
- [ ] FEATURE-0015 exposes no end-user CloudProvider projection and enforces Active participation as the minimum visibility precondition; FEATURE-0021 owns Organization/enrollment eligibility and FEATURE-0018 owns individual authorization.
- [ ] Historical alpha references occur only as explicitly labelled history or bounded test-fixture context; they cannot describe active behavior.
- [ ] The FEATURE-0015 architecture admission check rejects reintroduction of removed migration concepts and validates direct canonical bootstrap plus field-staging consistency.
- [ ] The FEATURE-0015 architecture admission check validates the approved resource fields/mutability, participation action-state table, independent suspension holds, topology safe-denial ordering, bootstrap grant actions, API/update behavior, concurrency, and audit matrix.
- [ ] The Phase 2R sequence, roadmap, traceability, context files, and feature authority agree on the canonical-bootstrap boundary.
- [ ] No Go implementation is modified by this architecture update.
- [ ] Kiro does not edit the current generated requirements artifact as a substitute for authority updates; requirements are regenerated only after the revised admission gate passes.

## Explicit instructions to Kiro

- Apply only this approved replacement and required consistency changes; do not introduce a different migration or compatibility model.
- Do not modify Go code, execute any conversion, delete historical implementation, or add runtime compatibility endpoints.
- Preserve completed FEATURE-0001–0014 as retained repository assets; assess and reuse compatible standards and implementations, but do not reinterpret superseded alpha domain kinds or routes as active canonical contracts.
- Preserve ADH-2026-043 decisions unrelated to runtime migration, including safe-denial, scope-reference integrity, audit reuse, and feature-local traceability.
- Preserve the canonical data model, seven scope kinds, CloudPlatform/CloudProvider boundary, ExecutionTarget direction, FEATURE-0013 audit ownership, and all Phase 2R no-side-effect constraints.
- Stop for human review if the replacement DEC or baseline-review path cannot be completed atomically.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-11
- Notes: “Replace alpha runtime migration with a canonical bootstrap. FEATURE-0001–0014 remain immutable implementation history. They create no active runtime data-migration obligation. The first Sovrunn control-plane release exposes canonical contracts only; obsolete alpha kinds and APIs are not carried forward as active compatibility surfaces.”
