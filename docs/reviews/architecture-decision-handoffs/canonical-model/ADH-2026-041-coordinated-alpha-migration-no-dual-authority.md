> **SUPERSESSION NOTICE (2026-08-11):** This handoff's migration runtime semantics are superseded in their entirety by **ADH-2026-045** ("Replace alpha runtime migration with a canonical bootstrap") and **DEC-0059** (supersedes DEC-0058). Sovrunn has no live control plane, customer data, or persisted alpha state to convert; FEATURE-0001–0014 are retained repository assets and reuse input, not live state requiring conversion. There is no `CanonicalMigrationPlan`, `CanonicalMigrationRecord`, migration controller, or cutover state machine in active Phase 2R authorities. This document is preserved below as an unmodified historical record — it is not rewritten or deleted.

# ADH-2026-041: Coordinated Alpha Migration and No-Dual-Authority Cutover

- **Status:** Architecture-owner approved; repository integration pending; **superseded in its entirety by ADH-2026-045 (2026-08-11)**
- **Handoff ID:** ADH-2026-041
- **Date:** 4 August 2026
- **Classification:** Breaking alpha migration policy
- **Canonical decision:** FCM-ADR-022

## Context

The finalized model changes approved alpha concepts, including Provider identity/scope, provider-prefixed topology and global ServiceClass semantics. Running old and new models as independent authorities would create split-brain authorization, catalog, placement and audit behavior.

## Decision

1. Migration is one coordinated release and cutover, not independent feature renaming.
2. Completed FEATURE-0001–0014 remain immutable history. Their live alpha data and APIs migrate through a signed `CanonicalMigrationPlan` and retained `CanonicalMigrationRecord` evidence.
3. Dual write and dual authority are prohibited. A bounded read-only legacy importer or compatibility projection is permitted only for migration and cannot accept desired-state writes.
4. Existing Provider records are classified, not blindly renamed:
   - customer-facing cloud/owner semantics become owner Organization + CloudPlatform;
   - infrastructure operator semantics become CloudProvider + CloudProviderParticipation;
   - an entity performing both roles becomes both relationships with separate identities and authority;
   - one provider participation receives one isolated SovrunnInstallation per environment.
5. ProviderLocation, ProviderDatacenter and DatacenterFailureDomain deterministically map to HostingLocation, Datacenter and FaultDomain with stable identity mapping and explicit operator/contract relationships. InfrastructureStack is retained only as an externally operated hosting fact.
6. ResourcePool and ProviderCapability placeholders are withdrawn before implementation; no runtime data migration is required for unimplemented resources.
7. Global ServiceClass splits into one reusable ServiceTypeDefinition and one or more CloudPlatform-scoped ServiceOfferings. ServicePlan, plugin, ServiceInstance, CLI and fixture references are rewritten through the migration map.
8. CloudProviderEnrollment becomes CloudEnrollment between customer Organization and CloudPlatform. Provider eligibility/assignment derives from CloudProviderParticipation and placement policy.
9. Provider scope splits into explicit CloudPlatform product/governance scope and CloudProvider operational scope. Ambiguous old scope references fail migration until classified.
10. Immutable historical DecisionRecords and AuditEvents are never rewritten in place. They remain in a read-only legacy archive or original schema representation, with append-only migration records and safe projections linking old and new UIDs.
11. One-to-one transforms preserve UID where semantics remain identical; split transforms allocate new UIDs and retain an immutable old-to-new mapping.
12. After cutover, old kind names, scope values and write endpoints are rejected with stable migration/deprecation errors. Aliases never become permanent contracts.

### Cutover state machine

```text
Draft
  → InventoryValidated
  → DryRunPassed
  → WriteFrozen
  → BackupVerified
  → Transformed
  → ReferencesVerified
  → CutoverActivated
  → ConformancePassed
  → Completed
```

Before `CutoverActivated`, failure may use the compatibility contract's permitted in-place or restore path. After the irreversible migration checkpoint, recovery follows the pinned `RestoreRequired` or `ForwardRecoveryOnly` rule; ad hoc downgrade is prohibited.

### Required gates

- complete inventory and classification with zero ambiguous authorities;
- signed backup and verified restore;
- deterministic dry run producing identical mappings on replay;
- zero unresolved or wrong-kind references;
- scope, authorization, decision, audit, catalog, topology, CLI and fixture conformance;
- no old desired-state writers or controllers running at cutover;
- source, generated schemas, docs, decisions, roadmap and traceability agree on the new vocabulary;
- post-cutover drift scan finds no old authoritative types outside retained history/migration tooling.

## Consequence

The migration may be operationally disruptive, but it is safer than long-lived aliases or dual control. Because the APIs are alpha, correctness and one authority take precedence over preserving obsolete names.
