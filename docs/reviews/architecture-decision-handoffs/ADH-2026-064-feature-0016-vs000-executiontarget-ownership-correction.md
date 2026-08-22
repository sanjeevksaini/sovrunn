# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-064
- Date: 2026-08-21
- Source discussion: Codex architecture session
- Related feature: FEATURE-0016
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Correct stale VS-000 skeleton ownership of ExecutionTarget

## Summary

`VS-000-core-skeleton.md` assigns `ExecutionTarget (identity)` to
FEATURE-0015, contradicting the higher-precedence canonical model, contract
registry, traceability matrix, feature authority, and FEATURE-0016 control
manifest. This correction changes that skeleton row only: FEATURE-0015 owns the
canonical cloud-model resources through `InfrastructureStack`; FEATURE-0016
owns ExecutionTarget, NormalizedTargetFactSet, and TargetQualificationResult in
their entirety.

No F0016 route, field, state, conformance case, or implementation changes.

## Classification

Correction.

## Existing approved baseline

- `docs/architecture/canonical/sovrunn-finalized-data-model.md` assigns
  ExecutionTarget identity, schema, routes, status, writer, conformance, and
  lifecycle entirely to FEATURE-0016.
- `VS0-SCHEMA-015..017`, `VS0-WRITER-006`, and `VS0-STATE-004` are all owned by
  FEATURE-0016 in the active registry and traceability matrix.
- FEATURE-0015 remains the read-only prior-feature authority for
  CloudProviderParticipation and InfrastructureStack.
- VS-000 core skeleton line 101 is the sole stale ownership statement.

## Decision

The VS-000 core skeleton's Cloud model row lists only:

`CloudPlatform`, `CloudProvider`, `CloudProviderParticipation`,
`HostingLocation`, `Datacenter`, `FaultDomain`, and `InfrastructureStack` as
FEATURE-0015-owned.

The Integration row lists `ExecutionTarget`, normalized target facts,
qualification, and the synthetic observer boundary as FEATURE-0016-owned.

FEATURE-0016 consumes CloudProviderParticipation and InfrastructureStack
read-only for backing resolution, authorization/safe access, viability, and
fence capture. It neither transfers nor modifies their FEATURE-0015 ownership.

## Rationale

The stale row halted the strict Kiro architecture-update workflow before
ADH-063 could be applied. Aligning it removes a competing lower-precedence
authority and prevents a future F0015 plan from accidentally reclaiming
ExecutionTarget.

## Reuse-before-build assessment

- Disposition: Clarify existing reuse.
- FEATURE-0015 CloudProviderParticipation and InfrastructureStack are reused
  read-only by FEATURE-0016.
- FEATURE-0012, FEATURE-0013, and FEATURE-0015 contracts retain their current
  ownership.
- No real adapter, credential, external call, persistence, controller, plugin,
  placement, provisioning, or ExecutionTarget subtype is introduced.

## Phase impact

- Current phase allowed? Yes.
- Current phase boundary impact: none; this removes stale documentation only.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting authority: `VS-000-core-skeleton.md` Cloud model row.
- Resolution required: approved ADH application only; no DEC, RFC, or baseline
  update.

## Required action

- Update VS-000 core skeleton ownership wording.
- Update directly dependent VS-000 traceability/check assertions only where
  they inspect that row.
- Add ADH-064 to FEATURE-0016 handoff evidence.
- Add a fail-closed assertion that FEATURE-0015 never owns ExecutionTarget.

## Strict Kiro read/write boundary

| Path | Access | Purpose |
|---|---|---|
| `docs/reviews/architecture-decision-handoffs/ADH-2026-064-feature-0016-vs000-executiontarget-ownership-correction.md` | read | Controlling approved decision |
| `docs/architecture/canonical/sovrunn-finalized-data-model.md` | read | Higher-precedence ownership authority |
| `docs/architecture/vertical-slices/VS-000-core-skeleton.md` | write | Remove stale FEATURE-0015 ExecutionTarget ownership |
| `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md` | write | Add ADH-064 provenance only if the ownership evidence needs it |
| `.automation/features/FEATURE-0016.control.json` | write | Add ADH-064 only to `feature.handoffs` |
| `scripts/vs000-contract-check.py` | write | Fail closed if skeleton assigns ExecutionTarget to FEATURE-0015 |
| `scripts/feature-0016-architecture-readiness-check.py` | write | Fail closed on stale F0015 ownership wording |

## Acceptance criteria for Kiro update

- [ ] The VS-000 skeleton assigns ExecutionTarget only to FEATURE-0016.
- [ ] FEATURE-0015 retains only read-only CloudProviderParticipation and
  InfrastructureStack ownership at this boundary.
- [ ] No F0016 route, field, state, conformance case, or source file changes.
- [ ] Registry and traceability continue to identify FEATURE-0016 as the sole
  owner of VS0-SCHEMA-015..017, VS0-WRITER-006, and VS0-STATE-004.
- [ ] Fail-closed checks reject any future skeleton assignment of
  ExecutionTarget to FEATURE-0015.

## Explicit instructions to Kiro

- Apply this approved correction before loading or applying ADH-063.
- Do not modify requirements, design, tasks, Go source, tests, formal models,
  routes, fields, states, or conformance catalog semantics.
- Do not update `CURRENT_ARCHITECTURE_BASELINE.md`.
- Stop with `ARCHITECTURE_DECISION_REQUIRED` if a second active F0015 ownership
  claim is found outside this allowlist.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-21
- Notes: Founder approved the ownership-only VS-000 skeleton correction.
