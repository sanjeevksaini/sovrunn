# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-043
- Date: 2026-08-10
- Source discussion: FEATURE-0015 architecture-readiness closure
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Close FEATURE-0015 migration, safe-denial, scope-integrity, audit, and traceability readiness gaps

## Summary

Apply the seven approved FEATURE-0015 closure decisions in `docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md`. The package preserves the Phase 2R canonical baseline, no-dual-authority migration policy, immutable history, safe reference isolation, and FEATURE-0013 audit ownership while making all F0015 obligations representable and locally testable before requirements are regenerated.

## Classification

Correction. The package corrects incomplete or mutually inconsistent F0015/VS-000 authorities and records bounded clarifications required to apply the accepted baseline.

## Existing approved baseline

The active baseline is `ARCH-2026.08-PHASE2R-CANONICAL`. DEC-0037 separates CloudPlatform and CloudProvider using seven canonical scopes. DEC-0054 keeps ownership, participation, and installation isolation distinct. DEC-0058 and ADH-2026-041 require a signed, one-way, no-dual-authority alpha cutover with immutable retained history. FEATURE-0013 remains the sole owner of DecisionRecord and AuditEvent schemas.

Relevant baseline references:

- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/decisions/DEC-0037-canonical-cloud-model-seven-scopes.md`
- `docs/decisions/DEC-0054-cloudplatform-provider-installation-isolation.md`
- `docs/decisions/DEC-0058-alpha-migration-no-dual-authority.md`
- `docs/reviews/architecture-decision-handoffs/canonical-model/ADH-2026-041-coordinated-alpha-migration-no-dual-authority.md`

## Decision or proposed decision

1. A CanonicalMigrationRecord is immutable evidence and has no direct correction/supersession link. A correction requires a superseding CanonicalMigrationPlan and a new linked migration run; retained prior records are never altered or re-ordered.
2. `Draft` is pre-persistence preparation of a migration plan. The signed immutable plan precedes the persisted CanonicalMigrationRecord sequence, which remains stages 1..9 from InventoryValidated through Completed.
3. FEATURE-0015 owns the signed global inventory/classification plan and executes provider/topology migration only. Later domain owners execute catalog, enrollment, governance, placement, plugin, and related transforms. FEATURE-0015 cannot claim global cutover completion; final all-domain cutover/conformance is FEATURE-0026 integration evidence.
4. An inaccessible cross-provider reference returns RESOURCE_NOT_FOUND (404) with VS0_AUTHORIZATION_SAFE_DENIAL and no existence disclosure. An authorized but structurally invalid same-provider reference returns VALIDATION_FAILED (422). VS0-CF-X03 records the safe-denial violation exactly.
5. CloudPlatform `spec.ownerOrganizationRef.uid` must equal its Organization `metadata.scopeRef.uid`; CloudProviderParticipation `spec.cloudPlatformRef.uid` must equal its CloudPlatform `metadata.scopeRef.uid`. A mismatch returns VALIDATION_FAILED (422) with VS0_SCOPE_REFERENCE_MISMATCH. Add an F0015-local conformance case.
6. FEATURE-0015 reuses FEATURE-0013 AuditEvent. Participation lifecycle changes, plan publication, each migration milestone, migration cutover decision, and safe-denial/security denials produce the applicable audit evidence with intentional correlation fields, safe projection/redaction, and no secrets. It does not import FEATURE-0026's integration-only trace conformance.
7. FEATURE-0015 traceability uses only its local conformance records (F15-01..11, MIG01/02, MIGF01..03, X03). Downstream HP01/F09 records are never F0015 acceptance evidence.

## Rationale

The decisions eliminate the representability conflict between immutable correction semantics and strict migration ordering, keep migration scope implementable one feature at a time, and prevent inconsistent disclosure behavior or contradictory ownership graphs. They also make audit and traceability explicit without changing the canonical model, adding a parallel audit envelope, or moving future-feature execution into FEATURE-0015.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse FEATURE-0012 typed references, Problem Details, metadata/scope grammar, and status ownership.
  - Reuse FEATURE-0013 AuditEvent and DecisionRecord envelopes.
  - Reuse the existing VS-000 registry, traceability, semantic checking, and Phase 2R in-memory execution boundary.
- Sovrunn-owned responsibility summary:
  - Canonical migration classification, immutable migration evidence, cross-provider isolation, and feature-local conformance.
- Non-goals summary:
  - Durable storage, real migration execution, provider SDK calls, catalog/enrollment/governance/placement/plugin transformation implementation, and an alternate audit system.

## Phase impact

- Current phase allowed? Yes.
- Current phase boundary impact:
  - Phase 2R remains side-effect-free and implementation-neutral. The package narrows F0015 to its foundation responsibilities and reserves all-domain cutover proof for FEATURE-0026; it introduces no future-phase runtime capability.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting decisions, if any:
  - None. The package resolves repository-authority inconsistencies under DEC-0037, DEC-0054, DEC-0058, and ADH-2026-041.
- Resolution required:
  - Approved Architecture Decision Handoff; no baseline replacement, DEC, RFC, or ACR is required.

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
- `.kiro/steering/slice0-contract.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `scripts/vs000-contract-check.py`
- `scripts/kiro-semantic-check.py`
- `Makefile`
- `docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md`
- `.kiro/specs/canonical-cloud-model-and-alpha-migration/requirements.md`

## Impacted features

- FEATURE-0015: corrected and bounded migration foundation.
- FEATURE-0016 through FEATURE-0025: plan-classified deferred migration domains only; no implementation work authorized.
- FEATURE-0026: owns all-domain cutover/integration conformance; no current implementation change.

## Acceptance criteria for Kiro update

- [ ] Handoff validated against Architecture Operating System files.
- [ ] No baseline replacement or Phase 2R scope expansion.
- [ ] CanonicalMigrationRecord correction language is representable and no longer conflicts with strict milestone ordering.
- [ ] Draft representation and the 1..9 persisted evidence sequence agree in every authority.
- [ ] Migration inventory clearly distinguishes F0015 execution from later-feature transforms and F0026 integration completion.
- [ ] Same-provider invalidity and inaccessible cross-provider safe denial have exact, registered outcomes.
- [ ] UID scope/reference invariants, error mapping, and F0015-local conformance are registered.
- [ ] F0015 audit/correlation/redaction requirements reuse FEATURE-0013 without importing downstream conformance.
- [ ] F0015 traceability contains no downstream-owned conformance as local evidence.
- [ ] A pre-requirements architecture-readiness check enforces these decisions.
- [ ] Requirements are revised only after all architecture checks pass.

## Explicit instructions to Kiro

- Apply only this approved handoff; do not add architecture decisions.
- Do not update `CURRENT_ARCHITECTURE_BASELINE.md`; this is not a baseline replacement.
- Do not modify Go code or implement migration execution.
- Keep all deferred domain transforms as plan-classified references only.
- Preserve FEATURE-0011/0012/0013 ownership and use existing Problem Details top-level codes.
- Update Structurizr only if the architecture model changes beyond these contract-level corrections; otherwise do not modify it.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-10
- Notes: Explicit instruction: "close with these decisions" after reviewing the seven proposed closure decisions.
