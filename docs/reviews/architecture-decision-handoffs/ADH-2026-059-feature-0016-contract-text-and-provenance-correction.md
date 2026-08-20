# ADH-2026-059 — FEATURE-0016 Contract Text and Provenance Correction

## Metadata

- Handoff ID: ADH-2026-059
- Date: 2026-08-20
- Source discussion: Codex architecture review
- Related feature: FEATURE-0016
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Correct FEATURE-0016 registry outcome text and FactSet/Result provenance terminology.

## Summary

This is a correction to the approved ADH-2026-058 package, not a new product
decision. It repairs six truncated registry `expectedState` values; reconciles
the persisted FactSet `factVersion` field with observer provenance; makes the
registered Result fence set exact; and removes a requirements-stage copy of
the inherited FEATURE-0012 Problem-code mapping. The mandatory exact 122-row
FEATURE-0016 conformance ledger remains required.

## Classification

- Correction

## Existing approved baseline

The approved Phase 2R baseline establishes ExecutionTarget as the qualified
realization boundary and keeps FEATURE-0012 as the sole owner of top-level
Problem code/HTTP mappings. ADH-2026-058 establishes VS0-SCHEMA-015..017,
VS0-STATE-004, and VS0-CF-F16-01..122 as FEATURE-0016's executable contract.

Relevant references:

- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`
- ADH-2026-058

## Decision or proposed decision

1. Correct the following `expectedState` values in the registry, without
   changing their inputs, errors, violations, side effects, gates, or behavior:

   | Case | Correct `expectedState` |
   |---|---|
   | VS0-CF-F16-16 | `200 ETag safe projection without facts, result, observer, or handle` |
   | VS0-CF-F16-20 | `200 Qualified/effectiveAvailability Available` |
   | VS0-CF-F16-21 | `200 Rejected/effectiveAvailability Unavailable` |
   | VS0-CF-F16-22 | `200 Indeterminate/effectiveAvailability Unavailable` |
   | VS0-CF-F16-36 | `200 Retired/Unqualified/effectiveAvailability Unavailable, ETag, refs cleared, AuditEvent, tuple release, completed replay` |
   | VS0-CF-F16-90 | `200 Indeterminate/effectiveAvailability Unavailable` |

2. `VS0-SCHEMA-016.factVersion=const[v1]` is the persisted FactSet version
   field. `factSchemaVersion=v1`, when named, is fixed observer provenance
   that explains the observer's emitted fact schema; it is not an additional
   persisted FactSet field.

3. `VS0-SCHEMA-017` records exactly its registered fences:
   `infrastructureStackGeneration` and `maintenanceEpoch`, alongside its
   target/fact-set references. `viabilityFingerprint` belongs only to
   `VS0-SCHEMA-016`; Result references the FactSet rather than duplicating it.

4. The requirements artifact must retain the exact 122-row registry ledger
   for VS0-CF-F16-01..122. It must not republish FEATURE-0012's top-level
   Problem-code/HTTP mapping; it cites FEATURE-0012 and the registry's
   `problemCodes` section instead.

## Rationale

The current malformed state text is an accidental truncation in the canonical
registry. The FactSet/Result wording otherwise invites design to add an
unregistered persisted field or duplicate a FactSet-only fence. Correcting the
canonical sources before requirements approval prevents architecture decisions
from leaking into design.

## Reuse-before-build assessment

- Disposition: Reuse
- Summary of mature candidates / applicable standards:
  - FEATURE-0012 owns Problem Details and HTTP mapping semantics.
  - FEATURE-0013 owns AuditEvent semantics.
- Sovrunn-owned responsibility summary:
  - FEATURE-0016 owns only its local observations, qualifications, and
    conformance evidence.
- Non-goals summary:
  - No new fields, routes, Problem codes, state transitions, adapter behavior,
    persistence mechanism, or external dependency.

## Phase impact

- Current phase allowed? Yes
- Current phase boundary impact:
  - Correction only; it preserves the Phase 2R side-effect-free synthetic
    observation boundary.

## Conflict check

- Conflicts with accepted DEC/RFC? No
- Conflicting decisions, if any: None
- Resolution required: None

## Required action

- Update architecture doc
- Update Kiro requirements.md
- Update traceability matrix
- Update feature gate/checks

## Impacted files

Kiro must first read the Architecture Operating System sources required by
`AGENTS.md`. Apart from those governance reads, it may read only:

- this handoff;
- `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`;
- `docs/features/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`;
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`;
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`;
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`;
- `.kiro/specs/adapter-boundary-and-executiontarget-qualification/requirements.md`;
- `scripts/feature-0016-architecture-readiness-check.py`;
- `scripts/feature-contract-check.py`;
- `scripts/vs000-contract-check.py`.

If approved, Kiro may write only:

- this handoff;
- `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`;
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`;
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`;
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`;
- `.kiro/specs/adapter-boundary-and-executiontarget-qualification/requirements.md`;
- `scripts/feature-0016-architecture-readiness-check.py`;
- `scripts/feature-contract-check.py`;
- `scripts/vs000-contract-check.py`.

Kiro must not modify the control manifest, design, tasks, Go source, schemas,
formal models, state files, baseline/DEC/RFC records, or any other file.

## Impacted features

- FEATURE-0016: corrected local conformance and requirements terminology.
- FEATURE-0012: cited only as the fixed owner of Problem semantics.

## Acceptance criteria for Kiro update

- [ ] The six registry states exactly match this handoff.
- [ ] Registry, architecture, contract specification, and traceability use
  `factVersion` as the persisted FactSet field and do not add a second field.
- [ ] Result text names only the registered generation and maintenance fences;
  `viabilityFingerprint` remains FactSet-only.
- [ ] The 122 F0016 conformance rows remain local exact evidence.
- [ ] Requirements cite rather than copy FEATURE-0012 Problem mappings.
- [ ] Readiness, feature-contract, VS-000, Phase 2R, semantic, whitespace,
  documentation, and Structurizr checks pass.

## Explicit instructions to Kiro

- Do not introduce any decision beyond these corrections.
- Preserve the approved five-route surface, state machine, audit, idempotency,
  security, and Phase 2R exclusions.
- Do not remove or weaken the exact 122-row requirements conformance ledger.
- Stop if a correction cannot be expressed within this write boundary.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-20
- Notes: Approved correction; no new product behavior.

This Architecture Decision Handoff is ready for Kiro validation and repo update only after human approval.
