# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-055
- Date: 2026-08-13
- Source discussion: FEATURE-0015 design executable-plan review
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Close FEATURE-0015 method, precedence, idempotency-retention, PATCH, and list-order mechanics

## Summary

This correction makes the already-owned API surface executable without adding a resource, route, grant, lifecycle state, Problem code, or persistence dependency. It closes exact `HEAD`, combined forged-input, completed-idempotency retention, RFC 7396 PATCH, deterministic LIST, and route-to-grant behavior needed by design and tasks.

**Implementation correction:** `VS0-CF-F15-34.expectedError` is `null`. HEAD's 405 is an HTTP transport outcome recorded in `expectedSideEffects`, not a FEATURE-0012 top-level Problem code. This implements, rather than changes, the approved no-new-Problem-code boundary.

## Classification

Correction — missing executable semantics for existing FEATURE-0015 observable behavior.

## Decision

1. `HEAD` is not an F0015 supported method. Although Go 1.22 `http.ServeMux` can match `HEAD` to a `GET` pattern, the handler checks the actual request method and returns 405, with no AuditEvent. This preserves exactly 22 logical paths and 35 explicit registrations.
2. After route resolution and successful authentication, a present `X-Sovrunn-Bootstrap-Grant` takes precedence over all media and body validation and returns audited `AUTHORIZATION_DENIED`/403. Without that header, an unsupported PATCH media type returns 415 before body decoding; therefore a body `bootstrapGrant` carrier is considered only after valid media and a syntactically valid duplicate-free decode.
3. A completed idempotency record is retained for 24 hours from completion. The process retains at most 10,000 completed records; on insertion it deterministically evicts the completed record with the earliest expiry, breaking ties by namespace lexical order. InFlight records are never evicted. After expiry/eviction there is no replay guarantee and the request is processed as new. Cancellation detaches only that waiter; every non-completion path aborts and wakes all remaining waiters.
4. PATCH uses RFC 7396 against a staged clone. Field classification occurs before merge; `null` removes only a mutable optional field, is invalid for a required mutable field, and is rejected for immutable/server-owned/deferred fields. Post-merge schema validation, current resourceVersion comparison, required audit append, then publication apply in that order.
5. Every successful LIST is deterministically ordered by ascending `metadata.uid`. Read authorization uses the current server-resolved read grant covering the resource's derived scope and, if the grant is UID-narrowed, that UID. Existing named write/action grants and registered writers remain unchanged.

## Rationale

These rules close Go runtime defaults and bounded in-process coordination choices before tasks can encode incompatible behavior. The contract remains canonical-bootstrap only and does not depend on durable storage.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse Go 1.22 `http.ServeMux`, RFC 7396 merge-patch semantics, standard-library sorting, contexts, and SHA-256.
  - Reuse FEATURE-0012 Problem Details and FEATURE-0013 AuditEvent contracts.
- Sovrunn-owned responsibility summary:
  - Enforce actual-method, precedence, retention, replay, patch, list, and grant-matrix behavior.
- Non-goals summary:
  - No router dependency, persistent idempotency store, new grant/role, new route, or new Problem/violation code.

## Phase impact

- Current phase allowed: Yes.
- Current phase boundary impact: FEATURE-0015 executable-contract correction only.

## Conflict check

- Conflicts with accepted DEC/RFC: No.
- Resolution required: None.

## Required action

- Update architecture/feature/registry/spec/traceability/steering/control/checkers.
- Update only affected requirements and design text/ledgers.
- Do not update tasks or Go code.

## Strict read/write boundary

Only these files may be read or written:

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

## Acceptance criteria

- [ ] HEAD is 405 while 22 paths/35 registrations remain unchanged.
- [ ] Header/body/media precedence is exact and has local conformance.
- [ ] Completed idempotency retention is bounded and deterministic; InFlight is not evicted.
- [ ] PATCH null and post-merge rules are executable without new API behavior.
- [ ] LIST ordering and current read-grant boundary are deterministic.
- [ ] Architecture, requirements, design, traceability, and deterministic checks agree.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-13
- Notes: Approved as a bounded FEATURE-0015 executable-contract closure.
