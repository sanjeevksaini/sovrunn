# ACR-2026-001: Final Canonical Model and Phase 2R Rebaseline

## Title

Adopt the service-first canonical model and rebaseline Phase 2 after FEATURE-0014.

## Change Type

Coordinated correction, replacement, extension and new decision.

## Current Approved Position

The repository baseline `ARCH-2026.07-PHASE2-START` treats:

- generic `Provider` as both governance/product and infrastructure-operator boundary;
- provider-prefixed topology as strict containment;
- `ResourcePool` as the mandatory placement boundary;
- `ProviderCapability` as the mandatory compatibility boundary;
- global `ServiceClass` as customer product identity;
- `EffectivePolicyContext` and overlapping governance/security profile kinds as the planned governance model;
- Phase 2 as FEATURE-0015 through FEATURE-0026 using those assumptions.

FEATURE-0001 through FEATURE-0014 are completed implementation history and remain immutable as delivered baselines.

## Proposed Change

Adopt the approved, content-bound `ADH-2026-020` through `ADH-2026-041` package and:

1. separate owner `Organization`, customer-facing `CloudPlatform`, independent `CloudProvider`, `CloudProviderParticipation` and customer `CloudEnrollment`;
2. classify and migrate the six-scope alpha model to seven canonical scopes: Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform and CloudProvider;
3. model `HostingLocation`, `Datacenter`, `FaultDomain` and `InfrastructureStack` as registered external facts and use qualified `ExecutionTarget` as the actionable realization boundary;
4. withdraw mandatory `ResourcePool` and `ProviderCapability` from the core and supersede `DEC-0032` and `DEC-0033`;
5. split reusable `ServiceTypeDefinition` from CloudPlatform-scoped `ServiceOffering` and versioned `ServicePlan`;
6. consolidate governance resolution into `GovernanceProfile`, `ProfileAssignment` and immutable `EffectiveGovernanceContext` while keeping sovereignty explicit;
7. use FEATURE-0013 `DecisionRecord` profiles for sovereignty and placement;
8. retain and strengthen `ServiceBinding`, add implementation-neutral service relationships and preserve adapter-only implementation details;
9. add independently recoverable Sovrunn platform lifecycle, directed release compatibility, maintenance fencing and platform-sovereignty composition;
10. perform one coordinated alpha cutover with no dual writers or permanent aliases;
11. replace the planned Phase 2 sequence from FEATURE-0015 onward with Phase 2R and prove it through Slice 0.

## Reason

The change realigns the repository with Sovrunn's constitutional service-first purpose: productize and govern PaaS services over reused execution infrastructures without becoming another infrastructure scheduler or leaking backend complexity to customers. The alpha stage is the least costly point to correct ownership, scope, catalog and placement boundaries.

## Alternatives Considered

- Continue with mandatory ResourcePool/ProviderCapability: rejected because it makes Sovrunn core responsible for a synthetic infrastructure inventory and duplicates native schedulers.
- Rename Provider only: rejected because it preserves cloud-product owner and infrastructure-operator conflation.
- Keep both old and new APIs: rejected because dual scope, catalog and placement authority creates split-brain authorization and audit behavior.
- Defer correction until after PostgreSQL MVP: rejected because later service, policy, placement and plugin contracts would build on the wrong boundary and make migration materially more expensive.

## Reuse Assessment

- **Disposition:** Extend
- **Mature candidates and applicable standards:** existing FEATURE-0011–0014 contracts; Kubernetes-style metadata/spec/status and typed references; FEATURE-0013 DecisionRecord/AuditEvent; OCI artifacts; GitOps; external identity, policy, secret, workflow, observability and service-operator systems.
- **Sovrunn-owned responsibility:** service product semantics; cloud participation and customer governance; implementation-neutral requirements; evidence-backed sovereignty; decisions; lifecycle intent; safe projections; plugin/adapter contracts and audit correlation.
- **Non-goals:** native infrastructure scheduling; custom identity/policy/secret/observability engines; provider SDK models in core; PostgreSQL HA implementation; permanent compatibility aliases; autonomous AI authority.

## Impacted Docs

At minimum:

- `docs/context/ARCHITECTURE_VERSION.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/context/CURRENT_DECISION_SUMMARY.md`
- `docs/decisions/DECISION_INDEX.md` and new `DEC-0037`–`DEC-0058`
- `docs/glossary.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/phase2/PHASE2_ACCEPTANCE_GATES.md`
- `docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md`
- `docs/features/FEATURE_INDEX.md`
- `docs/architecture/development-phases.md`
- `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md`
- `docs/traceability/DECISION_TRACEABILITY_MATRIX.md`
- Kiro context, prompts, feature controls and boundary checks for FEATURE-0015 onward.

Completed Phase 0/1 and FEATURE-0011–0014 documents remain historical implementation records. Active summaries and future specifications must point to the new baseline and migration rather than silently rewriting history.

## Impacted Features

- FEATURE-0006–0010: migration of catalog and service references; implementation history retained.
- FEATURE-0012: seven-scope alpha migration and conformance update.
- FEATURE-0013: scope fixtures migrate; DecisionRecord/AuditEvent semantics remain owned and otherwise unchanged.
- FEATURE-0014: resource data migrates/classifies; delivered feature history remains unchanged.
- FEATURE-0015–0026: fully rebaselined as Phase 2R.
- FEATURE-0027 onward: roadmap realigned to vertical slices and production outcomes.

## Backward Compatibility Impact

Breaking alpha migration. Old and new desired-state APIs never run concurrently. One-to-one semantic transforms preserve UID when valid; splits create new linked UIDs. Immutable DecisionRecord/AuditEvent history is retained rather than rewritten. Obsolete writes are rejected after cutover.

## Phase Impact

Allowed as a corrective Phase 2R rebaseline before FEATURE-0015 begins. Phase 2 remains side-effect-free except for in-memory/fake conformance execution. Real PostgreSQL and backend execution remain Phase 3.

## Security Impact

Positive: explicit owner/provider/customer boundaries, least-privileged participation, no-existence disclosure, provider credential isolation, SecretRef-only access, evidence-backed sovereignty, target fencing and no dual authorization authority. Migration requires write freeze, verified backup/restore and immutable audit evidence.

## Observability Impact

Every request, decision, migration step, plan, Operation and fake/real execution uses stable correlation and separate AuditEvent. Logs, metrics and traces remain operational telemetry and never substitute for audit or authoritative decisions.

## Decision Required

Yes. Human approval is recorded in `sovrunn-architecture-approval-record.md`; repository merge remains required.

## Recommendation

Accept and apply as one coordinated repository change. Do not generate FEATURE-0015 requirements from the old sequence.

## Approval

- **Proposed by:** Codex, from the finalized architecture discussion
- **Reviewed by:** Architecture-owner review requested through the current task
- **Approved by:** Sanjeev Kumar, through explicit instruction to complete approval and rebaseline
- **Date:** 4 August 2026
- **Repository authority:** Pending reviewed merge
