# Sovrunn Final Canonical Model — Architecture Handoff Index

**Package status:** Architecture-owner approved for controlled adoption; repository integration and migration pending
**Date:** 4 August 2026
**Allocated range:** `ADH-2026-020` through `ADH-2026-041`
**Repository baseline:** branch `phase2-reuse-first-paas-fabric-foundation`, commit `7bced1068f2837ead9709e78d3fc33f8e6e032c5`

These Architecture Decision Handoffs formalize the consequential boundaries in the final canonical model. The identifiers continue the repository sequence after approved `ADH-2026-019`. The architecture owner approved the package on 4 August 2026 through `ARCH-APPROVAL-2026-004`. The decisions become repository-authoritative only after the controlled integration is reviewed and merged.

The [repository baseline review](../sovrunn-final-canonical-baseline-review.md) records compatibility, conflicts and migration gates. The package is exercised by the [PostgreSQL end-to-end reference and conformance flow](../sovrunn-postgresql-reference-flow.md).

| Handoff | Decision | Baseline effect |
|---|---|---|
| [ADH-2026-020](ADH-2026-020-cloudprovider-identity-and-scope.md) | Retire ambiguous `Provider` through role classification | Adds distinct CloudPlatform and CloudProvider scopes under ADH-037/041 |
| [ADH-2026-021](ADH-2026-021-enrollment-relationship.md) | Cloud product and customer models meet through `CloudEnrollment` | Keeps provider participation separate from enrollment and authorization |
| [ADH-2026-022](ADH-2026-022-entitlement-and-quota.md) | Entitlement and quota remain independent | Refines planned FEATURE-0021 |
| [ADH-2026-023](ADH-2026-023-industry-cloud-portfolios.md) | Industry clouds are configured portfolios | Adds CloudPlatform product/catalog boundary |
| [ADH-2026-024](ADH-2026-024-geography-and-sovereignty.md) | Geography and sovereignty remain separate | Replaces provider-prefixed geographic semantics in FEATURE-0014 |
| [ADH-2026-025](ADH-2026-025-executiontarget-no-resourcepool.md) | `ExecutionTarget` is actionable; no mandatory `ResourcePool` | Replaces DEC-0032/0033 and reframes FEATURE-0015/0016 |
| [ADH-2026-026](ADH-2026-026-decisionrecord-adoption.md) | Sovereignty and placement reuse `DecisionRecord` | Extends FEATURE-0013 adoption into planned decision features |
| [ADH-2026-027](ADH-2026-027-immutable-versioned-definitions.md) | Published definitions are immutable by version | Extends the shared definition contract |
| [ADH-2026-028](ADH-2026-028-multitarget-placement.md) | One service may select multiple targets | Extends planned FEATURE-0023 and deployment planning |
| [ADH-2026-029](ADH-2026-029-customer-safe-placement.md) | Customer placement is a safe projection | Adds customer visibility/security boundary |
| [ADH-2026-030](ADH-2026-030-personal-organization.md) | Individuals receive default governance resources | Extends Phase 1 onboarding and IAM behavior |
| [ADH-2026-031](ADH-2026-031-extensibility-by-data.md) | Countries, sectors, services and backends extend through governed contracts | Strengthens extension and compatibility governance |
| [ADH-2026-032](ADH-2026-032-serviceoffering-catalog-evolution.md) | Split reusable `ServiceTypeDefinition` from provider `ServiceOffering` | Migrates FEATURE-0006/0007/0008/0010 catalog references |
| [ADH-2026-033](ADH-2026-033-governance-context-profile-consolidation.md) | Governance profiles compose into one effective context | Rebaselines planned FEATURE-0018–0021 and one FEATURE-0013 term |
| [ADH-2026-034](ADH-2026-034-servicebinding-access-boundary.md) | `ServiceBinding` is the managed-service access boundary | Retains and strengthens FEATURE-0008 for the PostgreSQL MVP |
| [ADH-2026-035](ADH-2026-035-service-relationship-boundary.md) | Service relationships are first-class and implementation-neutral | Adds reusable relationship contract, runtime resource and anti-drift boundary |
| [ADH-2026-036](ADH-2026-036-platform-lifecycle-boundary.md) | Sovrunn platform lifecycle is independently recoverable | Adds installation, release, policy, immutable plan, single external lifecycle authority, Operation-first atomic activation, external agent and infrastructure-maintenance boundary |
| [ADH-2026-037](ADH-2026-037-cloudplatform-provider-participation-installation-isolation.md) | One owner cloud may contract multiple isolated CloudProviders | Adds CloudPlatform, provider participation, CloudEnrollment, provider selection and one-provider-per-installation boundary |
| [ADH-2026-038](ADH-2026-038-platform-sovereignty-composition.md) | Platform sovereignty evaluates the installation and all controlling dependencies | Extends the shared sovereignty model with PlatformDependencySnapshot and a mandatory placement prerequisite |
| [ADH-2026-039](ADH-2026-039-release-recovery-compatibility-and-lifecycle-references.md) | Release compatibility is directed and recovery semantics are explicit | Adds compatibility dimensions, safe recovery modes, irreversible checkpoints and typed lifecycle reference contracts |
| [ADH-2026-040](ADH-2026-040-infrastructure-maintenance-authority-and-fencing.md) | External maintenance uses explicit authority, target state and fencing | Adds InfrastructureMaintenanceNotice, separate availability/qualification axes, epochs and requalification rules |
| [ADH-2026-041](ADH-2026-041-coordinated-alpha-migration-no-dual-authority.md) | Canonical alpha migration is coordinated and never dual-authoritative | Defines classification, mapping, cutover, historical integrity and post-cutover rejection rules |

## Completion status

- [x] Review the package against the current repository baseline.
- [x] Allocate collision-free repository `ADH-YYYY-NNN` identifiers.
- [x] Add missing decisions exposed by the baseline review: catalog evolution, governance-context consolidation and service binding.
- [x] Add the implementation-neutral service-relationship boundary exposed by the VM-to-PostgreSQL scenario.
- [x] Add the independently recoverable Sovrunn platform lifecycle and external infrastructure-maintenance boundary.
- [x] Trace the managed-service package through the canonical model and PostgreSQL reference flow.
- [x] Define separate upgrade, rollback, restore and infrastructure-maintenance conformance evidence for the platform lifecycle boundary.
- [x] Define the single external lifecycle source of truth and atomic Operation-first desired-state activation semantics.
- [x] Separate CloudPlatform ownership from CloudProvider participation and bind each installation to one provider participation.
- [x] Define evidence-backed platform sovereignty across all controlling dependencies.
- [x] Define directed release compatibility, recovery modes and typed lifecycle references.
- [x] Define infrastructure-maintenance announcement, target state ownership and fencing.
- [x] Define coordinated alpha migration with no dual write or dual authority.
- [x] Obtain architecture-owner approval.
- [ ] Integrate the handoffs and migrations into the repository.
- [x] Rebaseline the Phase 2R+ feature sequence from the accepted target model.

## Approval and integration sequence

1. Approve all 22 handoffs as one coordinated architecture amendment; do not partially adopt conflicting ownership, scope, catalog, topology, sovereignty, relationship, migration or platform-lifecycle models.
2. Integrate the shared vocabulary, scope and topology amendments first.
3. Migrate product catalog and governance contracts while the affected APIs remain alpha.
4. Replace the ResourcePool/ProviderCapability plan with ExecutionTarget qualification and adapter contracts.
5. Adopt sovereignty and placement DecisionRecord profiles, multi-target placement and customer-safe projections.
6. Implement and validate independent Sovrunn upgrade, rollback and restore before the production MVP gate.
7. Execute the coordinated alpha migration with one write authority and preserved historical evidence.
8. Run the PostgreSQL conformance flow before accepting feature requirements as implementation-ready.
