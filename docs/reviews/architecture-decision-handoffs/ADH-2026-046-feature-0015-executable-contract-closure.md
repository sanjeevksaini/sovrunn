# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-046
- Date: 2026-08-12
- Source discussion: Phase 2R FEATURE-0015 architecture-completeness review
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Close FEATURE-0015 executable-contract ownership, preconditions, and proof coverage

## Summary

This correction makes the approved canonical-bootstrap scope in ADH-2026-045 executable without design-stage invention.  It reconciles the apparent F0015 status-writer conflict, distinguishes collection creation from mutations of an existing participation, and requires complete FEATURE-0015-local conformance coverage for every owned observable behavior.  It adds fail-closed consistency checks so that documentation/registry drift, unmapped requirements, future-feature leakage, and later implementation proof gaps are detected before a Kiro requirements, design, or tasks stage is accepted.

## Classification

Correction and clarification of ADH-2026-045 execution semantics; no change to the canonical model, Phase 2R scope, resource ownership, or feature sequence.

## Existing approved baseline

`ARCH-2026.08-PHASE2R-CANONICAL`, DEC-0059, and ADH-2026-045 establish direct canonical bootstrap, FEATURE-0015 ownership through InfrastructureStack, FEATURE-0016 ownership of ExecutionTarget, PATCH-only managed-resource updates, action-only participation lifecycle, server-resolved grants, durable FEATURE-0013 audit evidence, and one authoritative control plane.

ADH-2026-045 requires feature-local proof for each owned route, transition, expiry, authorization, idempotency, audit, race, and PATCH case.  The controlling registry currently contains only partial F0015-local proof and uses generic controller language that can be read as conflicting with ADH-2026-045's API-server status ownership.

## Decision

### 1. F0015 status ownership

The FEATURE-0015 API server is the sole status writer for `CloudPlatform`, `CloudProvider`, `HostingLocation`, `Datacenter`, `FaultDomain`, and `InfrastructureStack`.  It creates each of these resources with `status.phase=Active`.  FEATURE-0015 has no independent topology, stack, provider, or asynchronous domain controller.

`CloudProviderParticipation.status` is likewise server-owned in FEATURE-0015.  The deterministic scheduler is a system actor that invokes the API server's expiry transition; it is not an additional status writer or controller authority.  `VS0-WRITER-005` remains the generic client-status-write prohibition, but its registered F0015 status writer must resolve explicitly to `api-server`, not a separately implemented controller.

### 2. Participation create versus existing-participation preconditions

`POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations` creates an absent resource and therefore requires `Idempotency-Key` but **does not require or accept `If-Match`**.  Its atomic uniqueness key is the non-terminal `(cloudPlatformUID, cloudProviderUID)` pair.  One concurrent create succeeds; the other returns `ALREADY_EXISTS` (409) with `VS0_PARTICIPATION_DUPLICATE`; neither duplicate resource nor duplicate audit event is stored.

PATCH and every action on an existing participation require both `If-Match` and `Idempotency-Key`.  The precondition subject is the participation item identified by the action route.  A missing, malformed, or non-current `If-Match` returns the existing FEATURE-0012 `STALE_RESOURCE_VERSION` Problem Details contract (412), performs no lifecycle/status write, and emits no additional AuditEvent.  The same principal, operation route, key, and canonical request digest returns the original result before precondition evaluation; a changed digest using the same key returns `CONFLICT` (409) with `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`.

### 3. Required complete F0015-local proof matrix

The registry must retain F15-01 through F15-11 and add the following F0015-owned conformance cases.  Every row must provide all registry conformance fields (`id`, `owner`, `inputs`, `expectedState`, `expectedError`, `expectedViolation` where applicable, `expectedSideEffects`, and `gate`) and must be reflected in the specification and traceability matrix.

| ID | Required observable proof |
|---|---|
| VS0-CF-F15-12 | CloudPlatform-root creation prerequisite: denied dependent create has 409/`VS0_CLOUDPLATFORM_ROOT_REQUIRED`, no resource, and no audit mutation except the required denial evidence. |
| VS0-CF-F15-13 | Valid PATCH accepts only `application/merge-patch+json`, changes only each resource's approved mutable fields, and increments its resource version. |
| VS0-CF-F15-14 | Invalid PATCH behavior: unsupported media type is 415; immutable/status/scope/identity/relationship/owner-registration changes are denied with the registered immutable/system-field outcome and leave the resource unchanged. |
| VS0-CF-F15-15 | Participation creation initializes `Pending`, both holds `false`, `requestExpiresAt = createdAt + 7 days`, and permits no active relationship before the CloudPlatform accept action. |
| VS0-CF-F15-16 | Scheduler expiry at `requestExpiresAt` changes only a still-Pending participation to `Expired`, is idempotent, and writes one correlated system-actor/scheduler AuditEvent. |
| VS0-CF-F15-17 | Existing-participation precondition behavior: missing, malformed, and stale `If-Match` return 412 `STALE_RESOURCE_VERSION`, make no state change, and create no additional AuditEvent. |
| VS0-CF-F15-18 | Same principal/route/key/digest replay returns the original create/action result before version checking and creates no duplicate resource, transition, or AuditEvent. |
| VS0-CF-F15-19 | Same principal/route/key with a changed canonical request digest returns 409/`VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`, retaining the original result. |
| VS0-CF-F15-20 | Every explicit participation transition, including reject, withdraw, request-release, accept-release, and decline-release to the correct hold-derived state, is authorized by its exact action grant and leaves invalid transitions unchanged. |
| VS0-CF-F15-21 | Read/list authorization: collection read without a read grant is 403 without resource detail; list with a scoped read grant returns only visible resources and an empty 200 collection where none are visible; unauthorized direct GET/reference resolution is safe 404. |
| VS0-CF-F15-22 | Bootstrap-grant boundary: server-resolved action/scope/resource restrictions authorize only their intended F0015 operation; body/header grant claims are ignored and denied. |
| VS0-CF-F15-23 | Topology reference ordering: an inaccessible cross-provider reference is safe 404/`VS0_AUTHORIZATION_SAFE_DENIAL`; a caller authorized for both resources with a structurally mismatched provider graph receives 422/`VS0_TOPOLOGY_PROVIDER_MISMATCH`. |
| VS0-CF-F15-24 | Audit atomicity: required resource create/PATCH, participation creation/action/expiry, authorization denial, and safe denial each produce exactly one redacted FEATURE-0013 AuditEvent with prescribed correlation; audit persistence failure leaves the associated mutation unchanged. |
| VS0-CF-F15-25 | F0015 route/action surface: owned collections/items expose only approved POST/GET/LIST/PATCH methods, participation actions have empty JSON bodies, and no F0016 ExecutionTarget or retired migration route/schema/writer/conformance is exposed. |

The architecture authority shall map each REQ-F15-01 through REQ-F15-18 and AC-F15-01 through AC-F15-14 to one or more F0015-local proof cases, except an explicit anti-drift/non-runtime verification gate.  No F0015 acceptance criterion may rely on a FEATURE-0016+ conformance case.  The mapping must include the existing F15 cases and the new F15-12 through F15-25 cases above.

### 4. Fail-closed guardrails

The FEATURE-0015 readiness checker must reject:

- a status-writer conflict among ADH-045/046, registry, contract specification, feature authority, architecture boundary, requirements, and traceability;
- an `If-Match` requirement on participation collection create, or missing `If-Match`/idempotency behavior on existing participation actions;
- any missing required F0015-local conformance case or REQ/AC-to-proof mapping;
- incomplete exact registry fields for an F0015 conformance case;
- F0016+ schema/writer/route/conformance/runtime behavior claimed by FEATURE-0015, except explicit labelled non-goal/exclusion references;
- a design or tasks document that lacks a testable task for every F0015-local conformance case.

The semantic checker must preserve exact canonical ledgers and must not classify retired migration identifiers occurring solely in an explicit non-goal/tombstone section as active conformance.  It must compare exact conformance ledger fields by normalized table-header name rather than a fixed presentation-column order, while continuing to reject unknown/retired IDs in active or ledger content.

## Rationale

This is a correction to documentation and verification completeness, not a new runtime capability.  It makes the already-approved F0015 contract deterministic for requirements, design, tasks, and tests, without allowing a generator to invent controllers, create preconditions, error semantics, route behavior, audit rules, or future-feature dependencies.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse FEATURE-0012 resource version, Problem Details, HTTP, PATCH, and idempotency conventions.
  - Reuse FEATURE-0013 AuditEvent semantics and correlation/redaction contract.
  - Extend existing VS-000 registry, traceability, readiness, semantic-validation, and feature-gate infrastructure.
- Sovrunn-owned responsibility summary:
  - Define the F0015-specific ownership, lifecycle, authorization, audit, and conformance composition of reused standards.
- Non-goals summary:
  - No migration runtime, no controller framework, no Federation, no F0016 ExecutionTarget behavior, no Go implementation, and no requirements/design/tasks regeneration in this architecture update.

## Phase impact

- Current phase allowed: Yes.
- Current phase boundary impact: clarifies and verifies the existing Phase 2R FEATURE-0015 canonical-bootstrap boundary only.

## Conflict check

- Conflicts with accepted DEC/RFC: No.
- Conflicting wording corrected: ADH-045 status ownership versus generic registry/controller wording; participation collection-create precondition wording.
- Resolution required: architecture/registry/traceability correction under this approved handoff; no new DEC, RFC, ACR, baseline, or feature-sequence change.

## Required action

- Update architecture doc
- Update feature authority
- Update contract registry and contract specification
- Update traceability matrix
- Update architecture readiness, semantic, and contract-validation checks
- Update FEATURE-0015 control context/digests only if necessary after the authority changes

## Impacted files

- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`
- `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md`
- `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`
- `docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-045-canonical-bootstrap-no-alpha-runtime-migration.md` (annotate the corrected execution details; do not rewrite its approved product decision)
- `scripts/feature-0015-architecture-readiness-check.py`
- `scripts/kiro-semantic-check.py`
- `scripts/vs000-contract-check.py` where required to validate the expanded F0015 matrix
- `.automation/features/FEATURE-0015.control.json` only if controlled context/digests require it
- `.kiro/steering/slice0-contract.md` only if its F0015 proof rule must be made explicit

Files to inspect and leave unchanged unless a direct consistency update is required:

- `docs/decisions/DEC-0059-canonical-bootstrap-replaces-alpha-migration.md`
- `docs/decisions/DECISION_INDEX.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/diagrams/structurizr/workspace.dsl`

Do not modify `.kiro/specs/**/requirements.md`, `design.md`, `tasks.md`, or Go source in this update.

## Impacted features

- FEATURE-0015: complete executable-contract and conformance closure.
- FEATURE-0016 through FEATURE-0026: protected from premature F0015 ownership; no behavior change.

## Acceptance criteria for Kiro update

- [ ] All decision sections above appear consistently in every controlling F0015 authority.
- [ ] F0015 statuses resolve to API-server ownership without a separate controller authority.
- [ ] Participation collection create has no `If-Match`; existing participation mutations have the approved precondition/idempotency rules.
- [ ] F15-01 through F15-25 are registered, exact, F0015-owned, and traceable; all REQ/AC mappings are complete.
- [ ] Existing canonical-bootstrap, non-migration, ExecutionTarget, and future-feature boundaries remain intact.
- [ ] Readiness and semantic checks reject the drift classes listed above.
- [ ] `make feature-0015-architecture-readiness`, `make vs000-contract-check`, `make phase2r-drift-check`, `git diff --check`, `mkdocs build --strict`, and `make structurizr-check` pass (or Structurizr reports its established environment-only CLI warning after workspace checks pass).
- [ ] No Go source or Kiro requirements/design/tasks artifact changed.

## Explicit instructions to Kiro

- Do not introduce decisions beyond this approved handoff.
- Do not update the architecture baseline, feature sequence, or DEC-0059.
- Do not modify Go source or Kiro requirements/design/tasks artifacts.
- Apply all active-authority, traceability, and guardrail changes atomically.
- If a controlling authority cannot express any approved decision exactly, stop and report the precise conflict rather than selecting a new semantic.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-12
- Notes: Final FEATURE-0015 architecture closure before requirements regeneration.  Requirements, design, and tasks must not discover architecture choices that this handoff is intended to close.
