# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-056
- Date: 2026-08-13
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Close FEATURE-0015 PATCH concurrency, read-grant, and idempotency replay semantics

## Classification

Clarification — exact observable behavior missing from the approved PATCH, read/list, and idempotency boundaries.

## Decision

1. Every FEATURE-0015 PATCH requires `If-Match` containing the current resource `metadata.resourceVersion`. Missing, malformed, or stale `If-Match` returns existing `STALE_RESOURCE_VERSION` / HTTP 412, with no mutation, AuditEvent, or idempotency record. The handler captures the parsed value before validation and, while holding the publication lock, rechecks it against the stored resource immediately before required AuditEvent append and publication. No client body field is added.
2. The server-resolved read actions are exact: `cloudplatform.read` for CloudPlatform GET/LIST at Platform-root scope; `cloudprovider.read` for CloudProvider GET/LIST at Platform-root scope; `topology.read` for HostingLocation, Datacenter, FaultDomain, and InfrastructureStack GET/LIST at CloudProvider scope; and `participation.read` for CloudProviderParticipation GET/LIST at its derived CloudPlatform scope. A grant may additionally be narrowed to the concrete GET target UID. LIST without the applicable read action returns the existing audited 403 boundary; inaccessible direct GET remains the existing safe 404 boundary.
3. A completed idempotency replay returns exactly the stored successful HTTP status and response body with `Content-Type: application/json`. It never replays `requestId`, correlation, `ETag`, `Location`, or another entity/transport header. The server assigns a fresh request/correlation ID to every received request, including a replay. The effective completed-record lifetime is the lesser of 24 hours from completion and process lifetime. Before re-reservation, an `Aborted` record is removed atomically under the idempotency lock; it cannot be replayed.

## Non-goals

- No new resource, route, lifecycle state, role, persisted grant, storage dependency, Problem code, violation code, or cross-process replay guarantee.
- No change to participation create/action `If-Match` rules, except that PATCH now has its own exact conditional-update rule.

## Required action

- Update the registry, contract specification, traceability, FEATURE-0015 feature and architecture authorities, closure matrix, steering, control manifest, and deterministic checks atomically.
- Update only the affected requirements and design text/ledgers.
- Do not update tasks or Go code.

## Strict read/write boundary

- this handoff;
- `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`;
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`;
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`;
- `docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md`;
- `.kiro/steering/slice0-contract.md`;
- `.automation/features/FEATURE-0015.control.json`;
- `scripts/feature-0015-architecture-readiness-check.py`;
- `scripts/feature-contract-check.py`;
- `scripts/vs000-contract-check.py`;
- `.kiro/specs/canonical-cloud-model-and-alpha-migration/requirements.md`;
- `.kiro/specs/canonical-cloud-model-and-alpha-migration/design.md`.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-13
- Notes: Approved as a bounded FEATURE-0015 design-executability closure.
