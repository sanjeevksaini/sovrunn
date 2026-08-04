# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-042
- Date: 2026-08-04
- Source discussion: Sovrunn Architecture Governor discussion and finalized canonical package
- Related feature: FEATURE-0015 through FEATURE-0026
- Related phase: Phase 2R
- Author: Codex draft from content-bound approved architecture
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Adopt the final canonical model and replace the post-FEATURE-0014 Phase 2 plan with Phase 2R.

## Summary

Adopt `ADH-2026-020` through `ADH-2026-041` as one content-bound semantic package, supersede mandatory ResourcePool/ProviderCapability placement, separate CloudPlatform product ownership from CloudProvider supply participation, introduce qualified ExecutionTarget realization, migrate catalog/governance/sovereignty/lifecycle contracts coherently, and rebaseline FEATURE-0015 through FEATURE-0026 around the synthetic VS-000 vertical slice. Completed FEATURE-0001 through FEATURE-0014 remain immutable delivery history; their active alpha data and dependent references migrate through one controlled cutover.

## Classification

Correction, Replacement, Extension and New decision applied as one coordinated architecture change.

## Existing approved baseline

The active repository baseline is `ARCH-2026.07-PHASE2-START` at reviewed commit `7bced1068f2837ead9709e78d3fc33f8e6e032c5`. FEATURE-0011 through FEATURE-0014 are implemented and merged. The current accepted sequence makes ResourcePool and ProviderCapability mandatory and retains generic Provider, global ServiceClass and overlapping policy/profile assumptions.

Relevant baseline references:

- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/context/CURRENT_DECISION_SUMMARY.md`
- `docs/decisions/DECISION_INDEX.md`
- `docs/glossary.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md`
- `docs/features/FEATURE_INDEX.md`
- `ADH-2026-017`, `ADH-2026-018` and `ADH-2026-019`

Approved package content is pinned by `sovrunn-canonical-adoption-content-manifest.sha256`. The controlling model is `sovrunn-finalized-data-model.md` version 1.7 and its companion contract catalog version 1.7.

## Decision or proposed decision

1. Accept the semantic decisions in ADH-2026-020 through ADH-2026-041 without independent downstream reinterpretation.
2. Add DEC-0037 through DEC-0058 according to `sovrunn-phase2r-rebaseline.md`; mark DEC-0032 and DEC-0033 Superseded.
3. Adopt `CloudPlatform`, `CloudProvider`, `CloudProviderParticipation` and `CloudEnrollment` as separate boundaries.
4. Migrate the alpha scope vocabulary to Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform and CloudProvider.
5. Treat HostingLocation, Datacenter, FaultDomain and InfrastructureStack as registered external facts; use qualified ExecutionTarget as the actionable realization boundary.
6. Do not introduce mandatory ResourcePool, ProviderCapability or native scheduler semantics into core.
7. Split reusable ServiceTypeDefinition from CloudPlatform-scoped ServiceOffering and versioned ServicePlan.
8. Resolve governance through immutable EffectiveGovernanceContext; keep sovereignty as a separate evidence-backed DecisionRecord profile.
9. Preserve FEATURE-0013 DecisionRecord/AuditEvent ownership and strengthen ServiceBinding as SecretRef-only per-consumer access.
10. Preserve implementation-neutral relationship, service lifecycle and platform lifecycle boundaries defined by the approved package.
11. Execute one CanonicalMigrationPlan with write freeze, verified backup/restore, deterministic mapping, no dual authority and preserved immutable history.
12. Replace FEATURE-0015 through FEATURE-0026 with the Phase 2R sequence and use VS-000 as its cross-feature acceptance slice.

## Rationale

Sovrunn creates value by transforming reused infrastructure into simple, governed and continuously provable cloud-native services. Mandatory synthetic infrastructure pools and provider-wide capability truth pull core toward infrastructure ownership, duplicate backend schedulers and leak the wrong abstraction upward. Correcting the alpha model before FEATURE-0015 prevents more expensive catalog, governance, placement, plugin and PostgreSQL rework later.

## Reuse-before-build assessment

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - completed FEATURE-0011–0014 contracts and conformance machinery;
  - FEATURE-0013 DecisionRecord/AuditEvent;
  - open metadata/spec/status, typed-reference, OCI, GitOps and external-system adapter patterns;
  - mature external identity, policy, secret, workflow, observability and service-operator systems.
- Sovrunn-owned responsibility summary:
  - cloud product and participation semantics;
  - customer governance and entitlement;
  - implementation-neutral requirements, decisions, lifecycle and safe projections;
  - sovereignty evidence composition and plugin/adapter contracts.
- Non-goals summary:
  - native scheduling or provider inventory replacement;
  - custom IAM, policy, secret, workflow, observability or database engines;
  - provider-native types in canonical/customer APIs;
  - autonomous AI authority.

## Phase impact

- Current phase allowed? Yes
- Target phase: Corrective Phase 2R before FEATURE-0015
- Current phase boundary impact:
  - Phase 2 remains model/decision/audit/adapter and deterministic simulation only.
  - Real provider, OpenShift and PostgreSQL execution remains Phase 3.
  - FEATURE-0015 through FEATURE-0026 are renamed and rescoped; their IDs remain stable.

## Conflict check

- Conflicts with accepted DEC/RFC? Yes
- Conflicting decisions:
  - DEC-0032 ResourcePool is the placement boundary.
  - DEC-0033 ProviderCapability is the compatibility boundary.
  - RFC-0024 and active Phase 2 documents encode the superseded topology/placement spine.
- Resolution required:
  - ACR-2026-001;
  - DEC-0032/0033 supersession and DEC-0037–0058 adoption;
  - RFC and baseline updates;
  - coordinated alpha migration and Phase 2R rebaseline.

## Required action

- Create Architecture Change Request ACR-2026-001.
- Update architecture docs and current baseline.
- Update DEC records/index and affected RFCs.
- Update phase scope, architecture spine, sequence, roadmap and feature index.
- Update glossary, traceability, Kiro context/prompts/controls and feature gates.
- Add the canonical model, contract catalog, reference flow, rebaseline and VS-000 charter to repository architecture docs.
- Generate FEATURE-0015 requirements only after repository readiness checks pass.
- Do not modify Go code in this architecture-update change.

## Impacted files

Exact paths and actions are defined in `sovrunn-repository-integration-manifest.md`. Minimum authorities include:

- `docs/context/ARCHITECTURE_VERSION.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/context/CURRENT_DECISION_SUMMARY.md`
- `docs/decisions/DECISION_INDEX.md`
- `docs/glossary.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/phase2/PHASE2_ACCEPTANCE_GATES.md`
- `docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md`
- `docs/features/FEATURE_INDEX.md`
- `docs/architecture/development-phases.md`
- architecture/decision traceability matrices;
- Kiro steering, agent resources, prompts, automation manifests and boundary checks.

## Impacted features

- FEATURE-0006–0010: migration inputs only; implementation history retained.
- FEATURE-0011: reuse standard retained.
- FEATURE-0012: scope vocabulary and conformance migration.
- FEATURE-0013: scope fixtures migrate; DecisionRecord/AuditEvent semantics retained.
- FEATURE-0014: data and active summaries migrate; completed history retained.
- FEATURE-0015–0026: replaced by approved Phase 2R sequence.
- FEATURE-0027–0034: aligned to the real PostgreSQL Slice 1 without redesigning Phase 2R contracts.
- Later roadmap: reclassified by outcome and revalidated after vertical-slice evidence.

## Acceptance criteria for Kiro update

- [ ] Content-manifest hashes match before import.
- [ ] ACR-2026-001 and DEC-0037–0058 are added and indexed.
- [ ] DEC-0032 and DEC-0033 are marked Superseded with replacement references.
- [ ] Architecture version becomes `ARCH-2026.08-PHASE2R-CANONICAL` only after all affected authorities agree.
- [ ] No active authority retains generic Provider as the combined product/operator concept.
- [ ] Exactly seven canonical scope kinds are used after cutover.
- [ ] CloudEnrollment targets CloudPlatform; CloudProviderParticipation remains separate.
- [ ] No active Phase 2 feature requires ResourcePool or ProviderCapability.
- [ ] Completed feature history is retained and distinguished from active canonical contracts.
- [ ] ServiceClass migration has one mapping to ServiceTypeDefinition and CloudPlatform-scoped ServiceOffering.
- [ ] EffectiveGovernanceContext is the only new canonical resolved-governance term.
- [ ] Sovereignty and placement use FEATURE-0013 DecisionRecord profiles.
- [ ] Customer APIs contain no provider-native object, target credential, raw secret or protected handle.
- [ ] Phase 2R sequence and dependencies match `sovrunn-phase2r-rebaseline.md`.
- [ ] VS-000 is imported as cross-feature acceptance context, not as a substitute feature spec.
- [ ] Kiro context loads ADH-2026-042 and the approved canonical docs, not chat history or mutable duplicate summaries.
- [ ] Required repository documentation, handoff, drift, traceability and strict MkDocs checks pass.
- [ ] Git diff contains no Go implementation or generated build output.

## Explicit instructions to Kiro

- Validate the content digests before applying the package.
- Treat ADH-2026-042 as the single controlling repository adoption handoff; preserve ADH-2026-020–041 as content-bound decision provenance.
- Apply the repository update atomically across all authority files; do not leave mixed old/new active terminology.
- Do not rewrite completed FEATURE-0001–0014 history as though it had originally implemented the new model.
- Do not create a compatibility alias or dual writer for old and new scopes/resources.
- Do not generate FEATURE-0015 requirements until repository alignment and tooling preflight pass.
- Generate one feature and one Kiro stage at a time.
- Stop with `REPOSITORY_CONTEXT_NOT_READY` for incomplete repository integration and `ARCHITECTURE_DECISION_REQUIRED` for any semantic gap.
- Do not implement code in this architecture update.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-04
- Notes: Approval is explicitly directed in the current task. Repository authority begins only after the reviewed Git integration is merged.

This approved Architecture Decision Handoff is ready for Kiro validation and the controlled repository update.
