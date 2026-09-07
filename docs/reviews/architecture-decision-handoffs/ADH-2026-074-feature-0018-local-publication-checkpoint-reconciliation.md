# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-074
- Date: 2026-09-05
- Source discussion: ChatGPT Project / Codex
- Related feature: FEATURE-0018 — Governance, IAM, Approval and Exception Foundation
- Related phase: Phase 2R
- Author: ChatGPT / Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

FEATURE-0018 local publication and checkpoint reconciliation

## Summary

This correction makes four already-approved FEATURE-0018-local implementation
mechanics exact: abort-safe publication sequence allocation, completed-record time
flow, caller-reservation ownership, and feature-gate integration. It adds no
resource, action, route, role, policy, event, error, external dependency, durable
persistence, or Phase 2R capability. It neither changes FEATURE-0013's contracts
nor makes FEATURE-0013 a transaction or acceptance service.

## Classification

- Correction

## Existing approved baseline

`ARCH-2026.08-PHASE2R-CANONICAL`, DEC-0060, and the approved ADH-2026-070 package
remain the FEATURE-0018 product-semantic authority. ADH-2026-072 clarifies that
FEATURE-0018 owns deterministic, local, non-durable publication machinery, while
FEATURE-0013 owns only shared decision/audit carrier contracts and validation.

The current approved design already assigns `state.Store` the sole caller-reservation
surface and assigns `idempotency` only neutral idempotency value construction and
validation. It also requires one state-issued publication instant, deterministic
time, atomic publication, and abort without publication. The task plan and committed
checkpoint contain local implementation ambiguities that must be reconciled to those
authorities before Cursor resumes.

Relevant baseline references:

- `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `docs/decisions/DEC-0060-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-072-feature-0018-design-authority-and-audited-evaluation-clarification.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/design.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/tasks.md`

## Decision or proposed decision

### 1. Abort-safe publication sequence

`state` derives a provisional publication sequence of `current + 1` while holding
`stateMu`. It binds that provisional sequence into the state-issued sealed
publication context, but persists the sequence in `state.Store` only as part of a
successful aggregate commit. An abort, failed compare-and-swap, cancellation,
evidence failure, or panic-cleanup path publishes nothing and consumes no sequence.

The successful commit validates that the provisional sequence is still the next
sequence for the committed root before persisting it. This is local, deterministic
in-memory state management; it is neither durable sequencing nor a FEATURE-0013
responsibility.

### 2. Exact completed-record timestamp flow

`operation.CompletionBinding` remains a neutral package-owned value and is extended
to carry the publication instant from the sealed, state-issued publication context.
No caller, owner, evidence component, or idempotency component reads a clock for
this purpose, and no second clock read is permitted.

`idempotency.PrepareCompletedResult` consumes that neutral completion binding and
derives the fixed 24-hour completed-record deadline from its publication instant.
`idempotency` must not import `state`, `evidence`, or `clock`; the completed record
retains the derived values needed for deterministic replay and expiry behavior.

### 3. Caller-reservation checkpoint ownership

`state.Store` remains the sole owner of `InspectOrReserveCallerMutation` and
`CallerReservationLease`, as specified by the approved design. Task 6 extends the
existing `internal/govaccess/state/store.go` implementation for completed-record
retention and replay behavior. It must not create a second caller-reservation type,
method, table, or declaration in `idempotency` or another package.

`idempotency` retains only its already-approved neutral lookup, request-binding,
prepared-completed-result, and completed-record value responsibilities. This
preserves the approved acyclic package graph and prevents competing reservation
writers.

### 4. Exact feature-gate integration path

Task 10 modifies the existing `scripts/feature-gate.sh` integration point. It may
modify `Makefile` only where that existing target needs wiring. The task plan must
remove the nonexistent alternative feature-gate path and must name no replacement
script.

## Rationale

The four corrections eliminate implementation ambiguity without reopening IAM
semantics. They keep publication deterministic, prevent aborted work from creating
observable sequence gaps, preserve one authoritative timestamp for retention, keep
the sealed API's ownership graph executable, and give Cursor an exact repository
path for the final feature gate.

## Reuse-before-build assessment

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse the existing FEATURE-0018 `state`, `operation`, and `idempotency` package
    boundaries and FEATURE-0013 carrier/validation contracts unchanged.
  - Reuse the existing repository `scripts/feature-gate.sh` integration point.
- Sovrunn-owned responsibility summary:
  - FEATURE-0018 owns only deterministic, process-local, non-durable publication,
    replay-retention, and conformance mechanics for its own aggregate.
- Non-goals summary:
  - No FEATURE-0013 service or contract expansion; no database, clock service,
    provider interaction, new event, schema, caller operation, error, route, or
    future-feature behavior.

## Phase impact

- Current phase allowed? Yes
- If not current phase, target phase: Not applicable
- Current phase boundary impact:
  - The correction remains deterministic, in-memory, zero-external-I/O Phase 2R
    conformance machinery. It makes no durable-persistence claim.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting decisions, if any:
  - The committed checkpoint and task-plan path declarations are inconsistent with
    the approved FEATURE-0018 design ownership and abort semantics.
- Resolution required:
  - Human approval of this correction, then Kiro reconciliation of the design and
    task plan. No DEC, RFC, ACR, baseline update, or FEATURE-0013 change is
    required.

## Required action

- Update Kiro design.md
- Update Kiro tasks.md
- Update feature gate/checks
- Update controlled task/review context and approval hashes as required
- Do not update requirements.md
- Do not modify Go code through this handoff

## Impacted files

Kiro must inspect and, where needed, update only the following planning or
integration artifacts:

- `docs/reviews/architecture-decision-handoffs/ADH-2026-074-feature-0018-local-publication-checkpoint-reconciliation.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/design.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/tasks.md`
- `.automation/features/FEATURE-0018.control.json`
- `.automation/context-projections/FEATURE-0018.tasks-review-requirements.spec.json`
- `.automation/context-projections/FEATURE-0018.tasks-review-requirements.md`
- `.automation/context-projections/FEATURE-0018.tasks-review-requirements.manifest.json`
- `.automation/state/FEATURE-0018.json` only for approved context/spec hash reconciliation
- `scripts/feature-gate.sh`
- `Makefile` only if existing target wiring requires it

The planned implementation paths, to be changed later by Cursor only after an
approved plan, are bounded to:

- `internal/govaccess/operation/`
- `internal/govaccess/state/store.go`
- `internal/govaccess/state/` tests
- `internal/govaccess/idempotency/`

## Impacted features

- FEATURE-0018: local design/task/checkpoint reconciliation only.
- FEATURE-0013, FEATURE-0012, FEATURE-0016, FEATURE-0017, FEATURE-0019 through
  FEATURE-0026: no change.

## Acceptance criteria for Kiro update

- [ ] Confirm the handoff introduces no product-semantic change and does not modify
  FEATURE-0013 ownership or contracts.
- [ ] State in design that a publication sequence is provisional until successful
  aggregate commit, and that every non-commit termination consumes no sequence.
- [ ] Define the sealed `operation.CompletionBinding` publication-instant flow and
  make `idempotency.PrepareCompletedResult` derive the 24-hour deadline without a
  clock read or a forbidden package import.
- [ ] Preserve `state.Store` as the sole owner of caller inspection/reservation and
  lease declarations; make Task 6 extend `state/store.go` rather than duplicate
  them.
- [ ] Name `scripts/feature-gate.sh`, plus `Makefile` only if needed, as Task 10's
  exact feature-gate integration path and remove the nonexistent alternative.
- [ ] Update task ownership, dependencies, verification commands, controlled review
  inputs, and required approvals without changing REQ, AC, resource, action, route,
  error, writer, or conformance semantics.
- [ ] Do not generate tasks beyond this reconciliation and do not modify Go code.

## Explicit instructions to Kiro

- Treat this as a narrow correction to FEATURE-0018-local design and task
  executability, not a change to IAM, approval, exception, or audit semantics.
- Do not add a clock dependency to `idempotency`, and do not make `state`,
  `evidence`, or `clock` an `idempotency` import.
- Do not create a second caller-reservation declaration or move that authority from
  `state.Store`.
- Do not create a new feature-gate script, alter FEATURE-0013, or claim durable
  publication behavior.
- Do not change requirements, architecture baseline, DEC-0060, Phase 2R scope,
  resource schema, routes, errors, event taxonomy, or conformance identifiers.
- Do not modify Go code until this handoff has been approved and Kiro has reconciled
  the controlled design and task plan.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-09-05
- Notes: Approved as a narrow FEATURE-0018-local correction for abort-safe
  publication, neutral completed-record timestamp flow, caller-reservation
  ownership, and exact feature-gate integration.
