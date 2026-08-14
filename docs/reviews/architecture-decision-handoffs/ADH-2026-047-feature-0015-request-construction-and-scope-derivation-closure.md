# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-047
- Date: 2026-08-12
- Source discussion: FEATURE-0015 contract-executability audit
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Close FEATURE-0015 request construction, scope derivation, and create outcomes

## Summary

This final FEATURE-0015 closure resolves the remaining pre-requirements architecture ambiguities: topology-name mutability, collection-create payloads, server-derived scope, participation create-time fields, and observable create/PATCH outcomes. It adds a fail-closed contract-executability audit before requirements generation. It does not alter the canonical model, feature order, Phase 2R scope, canonical-bootstrap decision, or the required exact requirements/acceptance/conformance ledgers.

## Classification

Correction and clarification of ADH-2026-045/046 execution semantics.

## Existing approved baseline

ADH-2026-045 establishes direct canonical bootstrap, F0015 ownership through InfrastructureStack, PATCH-only managed-resource updates, a CloudProviderParticipation action lifecycle, and API-server status ownership. ADH-2026-046 establishes API-server-only F0015 status writers, separate participation collection-create versus existing-item preconditions, and F15-01 through F15-25 local proof coverage.

The registry declares per-kind fields, scopes, mutability, and initial status, but the controlling package has not fully stated which party supplies each create-time field or how HostingLocation receives its CloudProvider scope. ADH-2026-045's phrase that topology "Name and description are PATCHable" conflicts with immutable identity in the registry.

## Decision

### 1. Closed collection-create request contract

For every F0015 collection `POST`, the client supplies only `metadata.name`, the required F0015-owned `spec` fields, and the optional F0015-owned `spec` fields listed for that kind below. The client must not supply `metadata.uid`, generation, resourceVersion, timestamps, `metadata.scopeRef`, any `status` field, a field owned by another feature, or an unknown field. The API server assigns identity/version/timestamps, `metadata.scopeRef`, initial status, and required FEATURE-0013 AuditEvent evidence.

| Kind | Client-required create fields | Client-optional create fields | Server-assigned outcome |
|---|---|---|---|
| CloudPlatform | `metadata.name`; `spec.ownerRegistration.legalName`; `spec.ownerRegistration.registrationIdentifier`; `spec.ownerRegistration.jurisdictionCode` | `metadata.displayName`; `spec.description` | deployment Platform-root `scopeRef`; `status.phase=Active`; `201` resource response |
| CloudProvider | `metadata.name`; non-empty `spec.operatingMarkets[]` | `metadata.displayName`; `spec.displayName` | deployment Platform-root `scopeRef`; `status.phase=Active`; `201` resource response |
| HostingLocation | `metadata.name`; `spec.countryCode`; `spec.locality` | `spec.administrativeAreaCode`; `spec.description` | CloudProvider `scopeRef` derived under decision 2; `status.phase=Active`; `201` resource response |
| Datacenter | `metadata.name`; `spec.hostingLocationRef` | `spec.description` | CloudProvider `scopeRef` derived from the resolved parent; `status.phase=Active`; `201` resource response |
| FaultDomain | `metadata.name`; `spec.datacenterRef` | `spec.description` | CloudProvider `scopeRef` derived from the resolved parent; `status.phase=Active`; `201` resource response |
| InfrastructureStack | `metadata.name`; `spec.faultDomainRef` | `spec.description` | CloudProvider `scopeRef` derived from the resolved parent; `status.phase=Active`; `201` resource response |

All create routes require `Idempotency-Key` under ADH-2026-045/046. A successful PATCH returns `200` with the updated resource. Existing FEATURE-0012 Problem Details semantics govern rejected unknown fields, server-owned fields, unsupported media types, immutable-field attempts, and other invalid input.

### 2. Scope derivation

`CloudPlatform` and `CloudProvider` receive the immutable deployment Platform-root scope from the API server. `CloudProviderParticipation` receives its immutable CloudPlatform scope from its resolved `spec.cloudPlatformRef`; the referenced CloudPlatform must exist and its UID must equal the resulting `metadata.scopeRef.uid`.

`HostingLocation` receives its immutable CloudProvider scope from the CloudProvider UID bound to the authenticated, server-resolved `topology.write` grant. A caller cannot supply or choose `metadata.scopeRef`. `Datacenter`, `FaultDomain`, and `InfrastructureStack` receive their immutable CloudProvider scope from their resolved immutable parent reference. Reference resolution and authorization/non-disclosure occur before structural validation; all derived scopes must match or the established safe-denial/validation outcome applies.

### 3. Topology immutability correction

`metadata.name` is immutable identity for every F0015 resource, including topology resources. For `HostingLocation`, `Datacenter`, `FaultDomain`, and `InfrastructureStack`, only `spec.description` is PATCHable. Parent reference, scope, identity, all metadata system fields, and status are immutable or server-owned. ADH-2026-045's topology sentence is corrected from "Name and description are PATCHable" to "Description is PATCHable; name is immutable identity."

### 4. Participation collection create body

`POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations` is a collection-create request, not an item lifecycle action. Its body requires:

- `metadata.name`
- `spec.cloudPlatformRef` (UID-pinned)
- `spec.cloudProviderRef` (UID-pinned)
- `spec.environment` with the only F0015-allowed value, `development`

It accepts no optional participation `spec` fields in F0015. In particular, `spec.providerSelectionModes` and `spec.permittedHostingLocationRefs` are FEATURE-0021-introduced and FEATURE-0021-activated fields; F0015 neither accepts, stores, defaults, validates, nor exposes them. The registry must record this ownership for both fields.

The API server derives `metadata.scopeRef` from `spec.cloudPlatformRef`, assigns `status.phase=Pending`, assigns `status.platformSuspended=false` and `status.providerSuspended=false`, and sets `status.requestExpiresAt=createdAt + 7 days`. It returns `201` with the created participation. The POST requires `Idempotency-Key` only and rejects/does not accept `If-Match`.

The item action routes (`:accept`, `:reject`, `:withdraw`, `:suspend`, `:resume`, `:request-release`, `:accept-release`, `:decline-release`) have an empty JSON body. Each existing-item action requires `If-Match` and `Idempotency-Key` under ADH-2026-046. The scheduler expiry is an internal API-server transition and has no public request body.

### 5. Executability audit and proof

Before a requirements-stage Kiro invocation, the F0015 readiness gate must run a deterministic contract-executability audit. The audit must fail closed unless every F0015 schema field is classified as exactly one of request-required, request-optional, server-assigned, action-only, deferred-to-owner-feature, or forbidden; every route states its request-body shape, required headers, success result, error family, status effect, audit effect, and local conformance; and all scopes have a single derivation source.

Add and trace these exact F0015-local conformance cases:

| ID | Required proof |
|---|---|
| VS0-CF-F15-26 | Per-kind collection-create request boundary: only the approved client fields are accepted; server-owned/status/scope/unknown/deferred fields are rejected; each successful create returns `201`, initializes its exact status, and has one required AuditEvent. |
| VS0-CF-F15-27 | Scope derivation: Platform-root, participation-from-CloudPlatform, HostingLocation-from-grant, and descendant-from-parent scopes are server derived; client-supplied/mismatched scope never selects another authority. |
| VS0-CF-F15-28 | Participation collection-create body versus empty item-action body: create accepts exactly the decision-4 fields, initializes Pending/false holds/seven-day expiry, rejects F0021 fields, and existing item actions accept only empty JSON body plus their required headers. |

The architecture authority must map REQ/AC rows to F15-26 through F15-28 where applicable. The readiness checker must validate the matrix and directly inspect all active ADH-045/046/047, registry, specification, feature, architecture, traceability, and generated requirements authorities. It must reject stale ADH-045 wording, absent/ambiguous create fields, a client-supplied scope, unowned participation fields, missing create statuses, or any requirement that omits the approved create/PATCH result.

### 6. Canonical ledgers remain mandatory

The requirements-stage document must retain exactly one verbatim Canonical requirement ledger, Canonical acceptance ledger, and Exact conformance semantics ledger. These are required by the manifest-controlled semantic checker and provide the approved exact traceability; they must not be replaced by compact citations. Compact mappings may supplement but never replace the mandatory ledgers.

## Rationale

This is the last contract-executability closure before requirements regeneration. It prevents design and tasks from inventing request payloads, scope source, initialization, field ownership, route body rules, or test behavior. It reuses FEATURE-0012 HTTP/metadata/error contracts and FEATURE-0013 audit contracts rather than introducing a new runtime model.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse FEATURE-0012 metadata, response, HTTP status, validation, Problem Details, and concurrency semantics.
  - Reuse FEATURE-0013 AuditEvent atomicity, correlation, and redaction semantics.
  - Extend existing registry, traceability, readiness, semantic-validation, and contract-check infrastructure.
- Sovrunn-owned responsibility summary:
  - Compose the exact F0015 resource-create, scope-derivation, field-ownership, and local-conformance behavior.
- Non-goals summary:
  - No federation, no provider integration, no migration runtime, no ExecutionTarget behavior, no FEATURE-0021 visibility/selection behavior, no Go implementation, and no requirements/design/tasks edit in this architecture update.

## Phase impact

- Current phase allowed: Yes.
- Current phase boundary impact: closes existing FEATURE-0015 canonical-bootstrap contracts only; no future feature enters F0015 scope.

## Conflict check

- Conflicts with accepted DEC/RFC: No.
- Corrected wording: ADH-2026-045 topology-name mutability and its overbroad statement that all participation actions require `If-Match`.
- Resolution required: atomic architecture, registry, traceability, and guardrail correction under this approved handoff. No new DEC, RFC, baseline, or feature-sequence change is required.

## Required action

- Update architecture doc and feature authority
- Update ADH-2026-045 annotation
- Update registry and contract specification
- Update traceability and closure matrices
- Update F0015 readiness, semantic, and registry validation checks
- Update control manifest/human digest only where necessary
- Regenerate requirements only after all architecture validations pass

## Impacted files

- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`
- `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md`
- `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`
- `docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-045-canonical-bootstrap-no-alpha-runtime-migration.md`
- `scripts/feature-0015-architecture-readiness-check.py`
- `scripts/kiro-semantic-check.py`
- `scripts/vs000-contract-check.py`
- `.automation/features/FEATURE-0015.control.json` only if controlled authority context/digests require it
- `.kiro/steering/slice0-contract.md` only if the pre-generation audit rule needs an explicit steering requirement

Inspect but do not change unless a direct consistency update is indispensable:

- `docs/decisions/DEC-0059-canonical-bootstrap-replaces-alpha-migration.md`
- `docs/decisions/DECISION_INDEX.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/diagrams/structurizr/workspace.dsl`

Do not modify `.kiro/specs/**/requirements.md`, `design.md`, `tasks.md`, or Go source as part of this architecture update.

## Impacted features

- FEATURE-0015: executable request-construction and scope-derivation closure.
- FEATURE-0021: retains sole ownership of participation provider-selection and permitted-hosting-location fields; no behavior change.
- FEATURE-0016 through FEATURE-0026: no behavior change; protected from early F0015 activation.

## Acceptance criteria for Kiro update

- [ ] All decisions in this handoff are expressed consistently across controlling authorities.
- [ ] Every F0015 collection-create kind has a closed client/server field matrix, server-derived scope, initial status, and `201` success outcome.
- [ ] Topology name is immutable; topology description is the only PATCHable topology spec field.
- [ ] Participation collection create and empty-body item actions are explicitly distinct.
- [ ] Both participation deferred fields are FEATURE-0021-owned and F0015-rejected.
- [ ] F15-01 through F15-28 are registered, exact, F0015-owned, and fully traceable.
- [ ] The executability audit fails before requirements generation on every ambiguity described in decision 5.
- [ ] Required exact canonical ledgers remain mandatory in requirements semantic validation.
- [ ] `make feature-0015-architecture-readiness`, `make vs000-contract-check`, `make phase2r-drift-check`, `git diff --check`, `mkdocs build --strict`, and `make structurizr-check` pass.
- [ ] No Go source or Kiro requirements/design/tasks artifact changed.

## Explicit instructions to Kiro

- Do not introduce new decisions beyond this handoff.
- Do not update the architecture baseline, feature sequence, or DEC-0059.
- Do not modify Go source or Kiro requirements/design/tasks artifacts.
- Apply active-authority, traceability, and guardrail changes atomically.
- Stop and report a specific conflict if any approved decision cannot be represented exactly.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-12
- Notes: This is the final F0015 architecture closure before requirements regeneration. The executability audit must pass before Kiro is permitted to generate requirements, design, or tasks.
