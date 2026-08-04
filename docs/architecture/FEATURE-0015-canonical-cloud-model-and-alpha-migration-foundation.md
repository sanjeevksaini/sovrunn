# FEATURE-0015 Architecture Boundary: Canonical Cloud Model and Alpha Migration Foundation

| Field | Value |
|-------|-------|
| Status | Approved boundary (correction under ADH-2026-042) |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Controlling Decisions | DEC-0037, DEC-0041, DEC-0042, DEC-0054, DEC-0058 |
| Controlling Handoffs | ADH-2026-020, ADH-2026-024, ADH-2026-025, ADH-2026-037, ADH-2026-041, consolidated ADH-2026-042 |
| Phase | 2R |
| Depends On | FEATURE-0011 (reuse), FEATURE-0012 (grammar/errors), FEATURE-0013 (decision/audit), FEATURE-0014 (alpha model) |

---

## 1. Purpose

Define the closed architecture boundary for FEATURE-0015 so that requirements, design, and tasks cannot absorb future-feature fields, stale FEATURE-0014 concepts, or unowned resource activation.

---

## 2. FEATURE-0015 Owned Resources and Activation Boundaries

FEATURE-0015 owns only the portions shown in the activation-boundary column:

| Resource | Registry ID | Scope | Profile | Activation Boundary |
|----------|------------|-------|---------|---------------------|
| CloudPlatform | VS0-SCHEMA-008 | Organization | ManagedResource | full identity + phase status |
| CloudProvider | VS0-SCHEMA-009 | Platform, Organization | ManagedResource | full identity + phase status |
| CloudProviderParticipation | VS0-SCHEMA-010 | CloudPlatform | ManagedResource | full lifecycle; state machine VS0-STATE-001 |
| HostingLocation | VS0-SCHEMA-011 | CloudProvider | ManagedResource | full identity + phase status |
| Datacenter | VS0-SCHEMA-012 | CloudProvider | ManagedResource | full identity + phase status |
| FaultDomain | VS0-SCHEMA-013 | CloudProvider | ManagedResource | full identity + phase status |
| InfrastructureStack | VS0-SCHEMA-014 | CloudProvider | ManagedResource | full identity + phase status |
| ExecutionTarget | VS0-SCHEMA-015 | CloudProvider | ManagedResource | **identity and spec only**; all `status.*` fields (qualification, availability, maintenanceEpoch, factSetRef, observedGeneration, conditions) are introducedBy AND activatedBy **FEATURE-0016** — FEATURE-0015 does not introduce, persist, or initialize them |
| CanonicalMigrationPlan | VS0-SCHEMA-060 | Platform | ImmutableRecord | FINAL on persistence; signed immutable plan |
| CanonicalMigrationRecord | VS0-SCHEMA-061 | Platform | ImmutableRecord | full; append-only milestone evidence |

---

## 3. Field Ownership Boundaries

### 3.1 ExecutionTarget field ownership

| Field | introducedBy | activatedBy | Writer |
|-------|-------------|-------------|--------|
| `metadata.*` | FEATURE-0015 | FEATURE-0015 | api-server (VS0-WRITER-001) |
| `spec.infrastructureStackRef` | FEATURE-0015 | FEATURE-0015 | cloud-provider-admin (VS0-WRITER-003) |
| `spec.participationRef` | FEATURE-0015 | FEATURE-0015 | cloud-provider-admin (VS0-WRITER-003) |
| `spec.targetClass` | FEATURE-0015 | FEATURE-0015 | cloud-provider-admin (VS0-WRITER-003) |
| `status.qualification` | **FEATURE-0016** | **FEATURE-0016** | fake-adapter-and-qualification-controller (VS0-WRITER-006) |
| `status.availability` | **FEATURE-0016** | **FEATURE-0016** | fake-adapter-and-qualification-controller (VS0-WRITER-006) |
| `status.maintenanceEpoch` | **FEATURE-0016** | **FEATURE-0016** | fake-adapter-and-qualification-controller (VS0-WRITER-006) |
| `status.factSetRef` | **FEATURE-0016** | **FEATURE-0016** | fake-adapter-and-qualification-controller (VS0-WRITER-006) |
| `status.observedGeneration` | **FEATURE-0016** | **FEATURE-0016** | target-lifecycle-controller (VS0-WRITER-005) |
| `status.conditions` | **FEATURE-0016** | **FEATURE-0016** | target-lifecycle-controller (VS0-WRITER-005) |

### 3.2 CloudProviderParticipation.spec.providerSelectionModes field ownership

| Field | introducedBy | activatedBy | Writer |
|-------|-------------|-------------|--------|
| `spec.providerSelectionModes` | **FEATURE-0021** | **FEATURE-0021** | delegated-participation-contract-authority (VS0-WRITER-004) |

FEATURE-0015 must NOT store, validate, default, or mention this field as CONTRACT_ONLY behavior. It is wholly owned by FEATURE-0021.

---

## 4. FEATURE-0016 Delegated Activation

FEATURE-0016 owns (both introduces and activates):
- All ExecutionTarget status fields: qualification, availability, maintenanceEpoch, factSetRef, observedGeneration, conditions
- Target qualification, availability and maintenance-epoch state machine (VS0-STATE-004)
- NormalizedTargetFactSet (VS0-SCHEMA-016)
- TargetQualificationResult (VS0-SCHEMA-017)
- Writing `ExecutionTarget.status.*` qualification/availability/epoch fields

ServiceRegion remains wholly owned by FEATURE-0022. FEATURE-0016 is only a declared prerequisite for FEATURE-0022's `spec.executionTargetRefs` and `status.availability`; it does not introduce or activate either field.

---

## 5. FEATURE-0021 Delegated: Provider-Selection Intent

FEATURE-0021 owns:
- `CloudProviderParticipation.spec.providerSelectionModes` field (introducedBy/activatedBy FEATURE-0021)
- `providerSelectionModes` vocabulary activation in placement evaluation context
- `ServiceInstance.spec.providerPreferenceRef` evaluation semantics
- Provider-selection intent evaluation and routing logic

FEATURE-0015 does NOT introduce, store, validate, default, or reference `providerSelectionModes` in any active behavior.

---

## 6. Removed from Slice 0: SovrunnInstallation

`SovrunnInstallation` (former VS0-SCHEMA-057) is **removed from FEATURE-0015 scope** and from Slice 0 active schemas because:

- No Phase 2R feature owns platform lifecycle (DEC-0053 scope is deferred to Phase 3+).
- The canonical model defines SovrunnInstallation as platform-operator-facing lifecycle; it requires `PlatformLifecyclePolicy`, `SovrunnRelease`, `PlatformLifecyclePlan`, and external `PlatformLifecycleAgent` — none of which are Phase 2R scope.
- The VS-000 fixture graph referenced it for provider isolation proof, but DEC-0054's one-provider-per-participation invariant is already proven by `CloudProviderParticipation` uniqueness.

**Retirement**: VS0-SCHEMA-057 is retired as a tombstone in the retiredSchemas section of the registry. The ID is permanently consumed and must not be reused. The fixture graph uses `CloudProviderParticipation` as the participation boundary proof.

---

## 7. CanonicalMigrationPlan and CanonicalMigrationRecord

### 7.1 Authority

- DEC-0058: alpha migration is one signed cutover with no dual write/authority and immutable history preservation.
- ADH-2026-041: Controlling handoff for migration semantics.
- Canonical data model §16.1: Defines migration sequence and classification rules.

### 7.2 CanonicalMigrationPlan (VS0-SCHEMA-060)

| Attribute | Value |
|-----------|-------|
| Profile | ImmutableRecord |
| Scope | Platform |
| Boundary | platform-operator-facing |
| Purpose | Signed, immutable plan that defines the complete alpha-to-canonical cutover classification mappings, scope vocabulary mapping, and resource transforms |
| Immutability | FINAL on persistence (VS0-STATE-010); no field may change after signing; corrections create a new linked plan superseding the previous |
| Writer | approved-migration-plan-publisher (VS0-WRITER-021); the migration controller is forbidden from authoring or approving its plan |

### 7.3 CanonicalMigrationRecord (VS0-SCHEMA-061)

| Attribute | Value |
|-----------|-------|
| Profile | ImmutableRecord |
| Scope | Platform |
| Boundary | platform-operator-facing |
| Purpose | Append-only milestone evidence: each record proves one migration milestone was reached; records form an ordered sequence linked to their plan and predecessor record |
| Milestone Sequence (VS0-STATE-011) | InventoryValidated → DryRunPassed → WriteFrozen → BackupVerified → Transformed → ReferencesVerified → CutoverActivated → ConformancePassed → Completed |
| Fields | record.milestone (enum), record.stage (integer 1..9), record.predecessorRef (TypedRef<CanonicalMigrationRecord>, nil for first), plus source/target/classification/transform/digest fields |
| Immutability | Append-only; each record is FINAL on persistence (VS0-STATE-010); corrections create linked records |
| Writer | migration-controller (VS0-WRITER-020) |

---

## 8. Explicitly Excluded (FEATURE-0015 Must Not Implement)

| Concept | Reason | Owner |
|---------|--------|-------|
| ResourcePool | Superseded by DEC-0042 | None (removed) |
| ProviderCapability | Superseded by DEC-0042 | None (removed) |
| Generic Provider | Superseded by DEC-0037 | None (split into CloudPlatform + CloudProvider) |
| Target qualification/availability writes | FEATURE-0016 activation | FEATURE-0016 |
| NormalizedTargetFactSet | FEATURE-0016 | FEATURE-0016 |
| Provider-selection intent evaluation | FEATURE-0021 | FEATURE-0021 |
| CloudEnrollment | FEATURE-0021 | FEATURE-0021 |
| SovrunnInstallation lifecycle | No Phase 2R owner | Deferred |
| Platform lifecycle (SovrunnRelease, PlatformLifecyclePolicy, PlatformLifecyclePlan) | DEC-0053, no Phase 2R owner | Deferred |
| ServiceClass as canonical catalog | Superseded by DEC-0049 | None |
| EffectivePolicyContext | Superseded by DEC-0050 | None |
| Six-scope vocabulary | Superseded by DEC-0037 | None |

---

## 9. Classification Semantics

| Element | Classification |
|---------|---------------|
| REQUIREMENTS-owned | Observable behavior that acceptance tests verify: resource CRUD, validation, state transitions, error codes, writer enforcement, cross-provider isolation, migration dry-run determinism |
| DESIGN-delegated | Internal mechanics: registry data structures, handler wiring, validation plumbing, in-memory store shape |
| Downstream negative boundary | ExecutionTarget status fields owned entirely by FEATURE-0016; visible here only to forbid FEATURE-0015 activation |
| EXCLUDED | Concepts explicitly not in scope per §8 |

---

## 10. Traceability Mapping

| Owned Item | DEC/ADH | VS0 Schema | VS0 Writer | VS0 State | VS0 Conformance |
|------------|---------|------------|------------|-----------|-----------------|
| CloudPlatform | DEC-0037; ADH-020/042 | VS0-SCHEMA-008 | VS0-WRITER-002,003 | — | VS0-CF-HP01 |
| CloudProvider | DEC-0037; ADH-020/042 | VS0-SCHEMA-009 | VS0-WRITER-003 | — | VS0-CF-HP01 |
| CloudProviderParticipation | DEC-0054; ADH-037/042 | VS0-SCHEMA-010 | VS0-WRITER-004 | VS0-STATE-001 | VS0-CF-HP01,F09 |
| HostingLocation | DEC-0041; ADH-024/042 | VS0-SCHEMA-011 | VS0-WRITER-003 | — | VS0-CF-HP01 |
| Datacenter | DEC-0041; ADH-024/042 | VS0-SCHEMA-012 | VS0-WRITER-003 | — | VS0-CF-HP01 |
| FaultDomain | DEC-0041; ADH-024/042 | VS0-SCHEMA-013 | VS0-WRITER-003 | — | VS0-CF-HP01 |
| InfrastructureStack | DEC-0041,0042; ADH-025/042 | VS0-SCHEMA-014 | VS0-WRITER-003 | — | VS0-CF-HP01 |
| ExecutionTarget (identity) | DEC-0042; ADH-025/042 | VS0-SCHEMA-015 | VS0-WRITER-003 | — | VS0-CF-HP01 |
| CanonicalMigrationPlan | DEC-0058; ADH-041/042 | VS0-SCHEMA-060 | VS0-WRITER-021 | VS0-STATE-010 | VS0-CF-MIG01 |
| CanonicalMigrationRecord | DEC-0058; ADH-041/042 | VS0-SCHEMA-061 | VS0-WRITER-020 | VS0-STATE-010,011 | VS0-CF-MIG02 |
| Cross-provider isolation | DEC-0037,0054 | VS0-SCHEMA-015 | VS0-WRITER-003 | — | VS0-CF-X03 |

---

## 11. Anti-Drift Rules

1. FEATURE-0015 requirements/design/tasks must NOT reference ResourcePool, ProviderCapability, generic Provider, ServiceClass, EffectivePolicyContext, or six-scope vocabulary.
2. FEATURE-0015 must NOT introduce, initialize, persist, validate, default, or write ExecutionTarget `status.*` fields — FEATURE-0016 owns their introduction and activation.
3. FEATURE-0015 must NOT evaluate provider-selection intent — that is FEATURE-0021 scope.
4. FEATURE-0015 must NOT implement platform lifecycle (SovrunnInstallation, SovrunnRelease, PlatformLifecyclePolicy, PlatformLifecyclePlan).
5. Migration state machine (VS0-STATE-011) is the accepted canonical sequence from DEC-0058/data-model §16.1; FEATURE-0015 must not invent additional states.
6. The migration controller may append CanonicalMigrationRecord evidence under an approved signed plan; it must never author, approve, replace, or mutate CanonicalMigrationPlan.

---

*End of architecture boundary.*
