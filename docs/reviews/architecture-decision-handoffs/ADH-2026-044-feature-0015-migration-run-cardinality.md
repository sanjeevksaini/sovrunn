> **SUPERSESSION NOTICE (2026-08-11):** This handoff is superseded in its entirety by **ADH-2026-045** ("Replace alpha runtime migration with a canonical bootstrap") and **DEC-0059** (supersedes DEC-0058). This handoff exists solely to govern `CanonicalMigrationPlan`/`CanonicalMigrationRecord` run cardinality; since those resources no longer exist as active Phase 2R behavior, its run-cardinality decision is moot. This document is preserved below as an unmodified historical record — it is not rewritten or deleted.

# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-044
- Date: 2026-08-10
- Source discussion: FEATURE-0015 requirements review resolution
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved; **superseded in its entirety by ADH-2026-045 (2026-08-11)**

## Decision title

Make CanonicalMigrationRecord completion run-local under a signed global plan

## Summary

A signed CanonicalMigrationPlan declares stable migration `runKey` values in its transform inventory. CanonicalMigrationRecords form one strict 1..9 milestone chain per `(planRef.uid, runKey)`. `Completed` seals that run only. FEATURE-0015 owns the provider-topology run; FEATURE-0026 alone proves every required plan run complete and the global cutover/conformance result.

## Classification

Correction

## Existing approved baseline

DEC-0058 and ADH-2026-041 require a signed global migration plan, immutable history, and no dual authority. ADH-2026-043 correctly limits FEATURE-0015 transforms to provider/topology but did not define record-chain cardinality.

## Decision or proposed decision

- `CanonicalMigrationPlan.record.resourceTransforms[]` declares a unique immutable `runKey` for each transform domain.
- `CanonicalMigrationRecord.record.runKey` is required and must match one plan-declared runKey.
- `VS0-STATE-011` ordering is strict only within the same `(planRef.uid, runKey)` pair.
- A `Completed` record seals that run; it is not a global migration completion claim.
- FEATURE-0015 may append only the `provider-topology` run. Later domain owners append only their assigned plan-declared runs.
- FEATURE-0026 proves global completion only after every required run is sealed and its integration conformance passes.
- F0015 `VS0-CF-MIG01` proves deterministic provider-topology run completion, with no legacy desired-state writers for that run.

## Rationale

The decision makes deferred transforms and immutable evidence representable without reopening a sealed chain or allowing FEATURE-0015 to claim all-domain cutover.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse existing signed plan, TypedRef, immutable record, VS-000 state-machine, and AuditEvent contracts.
- Sovrunn-owned responsibility summary:
  - Plan-declared migration-run identity and per-run conformance.
- Non-goals summary:
  - A new migration runtime, durable store, parallel audit envelope, or early execution of deferred domains.

## Phase impact

- Current phase allowed? Yes.
- Current phase boundary impact:
  - No Phase 2R scope expansion; this assigns evidence ownership only.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting decisions, if any: None.
- Resolution required: Approved Architecture Decision Handoff; no baseline replacement.

## Required action

- Update architecture doc
- Update phase scope doc
- Update Kiro requirements.md
- Update traceability matrix
- Update feature gate/checks

## Impacted files

- `docs/features/FEATURE-0015-canonical-cloud-model-and-alpha-migration-foundation.md`
- `docs/architecture/FEATURE-0015-canonical-cloud-model-and-alpha-migration-foundation.md`
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`
- `.automation/features/FEATURE-0015.control.json`
- `scripts/feature-0015-architecture-readiness-check.py`
- `.kiro/specs/canonical-cloud-model-and-alpha-migration/requirements.md`

## Impacted features

- FEATURE-0015: provider-topology migration run only.
- FEATURE-0016–0025: assigned deferred runs only when their feature is active.
- FEATURE-0026: global migration-completion/integration proof.

## Acceptance criteria for Kiro update

- [ ] Plan and record schemas define immutable runKey linkage.
- [ ] State and conformance define per-run cardinality and run-local Completed.
- [ ] F0015 does not claim global completion.
- [ ] F0026 owns all-run/global completion proof.
- [ ] Readiness gate rejects missing runKey/cardinality semantics.
- [ ] Requirements are regenerated only after all authorities agree.

## Explicit instructions to Kiro

- Apply only this approved handoff and ADH-2026-043.
- Do not modify Go code or execute deferred migrations.
- Correct CloudPlatform writer traceability to VS0-WRITER-002 only.
- Keep stale terms only in explicit exclusion/non-goal sections.
- Do not remove required canonical REQ/AC or exact conformance ledgers.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-10
- Notes: Explicit approval of the run-local Completed decision.
