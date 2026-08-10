# FEATURE-0015 Architecture Boundary: Canonical Cloud Model and Alpha Migration Foundation

| Field | Value |
|-------|-------|
| Status | Approved boundary (correction under ADH-2026-042) |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Controlling Decisions | DEC-0037, DEC-0041, DEC-0042, DEC-0054, DEC-0058 |
| Controlling Handoffs | ADH-2026-020, ADH-2026-024, ADH-2026-025, ADH-2026-037, ADH-2026-041, consolidated ADH-2026-042, ADH-2026-043, ADH-2026-044 |
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

The registry's `fieldOwnership.introducedBy` and `activatedBy` rule governs feature activation boundaries. It does not alter the source-of-truth precedence order: the registry remains the machine-readable expression of the higher-precedence Slice 0 specification.

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
- ADH-2026-044: Migration completion is run-local under plan-declared `runKey`; FEATURE-0015 owns the provider-topology run; FEATURE-0026 proves global completion.
- Canonical data model §16.1: Defines migration sequence and classification rules.

### 7.2 CanonicalMigrationPlan (VS0-SCHEMA-060)

| Attribute | Value |
|-----------|-------|
| Profile | ImmutableRecord |
| Scope | Platform |
| Boundary | platform-operator-facing |
| Purpose | Signed, immutable plan that defines the complete alpha-to-canonical cutover classification mappings, scope vocabulary mapping, and resource transforms |
| Run declaration (ADH-2026-044) | `record.resourceTransforms[]` declares a unique immutable `runKey` per transform domain; every appended CanonicalMigrationRecord must bind to one plan-declared runKey |
| Immutability | FINAL on persistence (VS0-STATE-010); no field may change after signing; corrections create a new linked plan superseding the previous |
| Writer | approved-migration-plan-publisher (VS0-WRITER-021); the migration controller is forbidden from authoring or approving its plan |

### 7.3 CanonicalMigrationRecord (VS0-SCHEMA-061)

| Attribute | Value |
|-----------|-------|
| Profile | ImmutableRecord |
| Scope | Platform |
| Boundary | platform-operator-facing |
| Purpose | Append-only milestone evidence: each record proves one migration milestone was reached; records form an ordered sequence linked to their plan and predecessor record |
| Run cardinality (ADH-2026-044) | Each record carries a required immutable `record.runKey` matching one plan-declared runKey; the strict 1..9 milestone chain is ordered only within one `(planRef.uid, runKey)` run; `Completed` seals that run and is not a global-completion claim |
| Milestone Sequence (VS0-STATE-011) | InventoryValidated → DryRunPassed → WriteFrozen → BackupVerified → Transformed → ReferencesVerified → CutoverActivated → ConformancePassed → Completed |
| Fields | record.planRef (uid-pinned), record.runKey (string 1..63, required, immutable), record.milestone (enum), record.stage (integer 1..9), record.predecessorRef (TypedRef<CanonicalMigrationRecord>, nil for the first milestone of a run), plus source/target/classification/transform/digest fields |
| Immutability | Append-only; each record is FINAL on persistence (VS0-STATE-010); no direct correction or supersession link on records |
| Correction Model | A CanonicalMigrationRecord has no correction/supersession link. A correction requires a superseding CanonicalMigrationPlan and a new linked migration run; retained prior records are never altered or re-ordered (ADH-2026-043 decision 1) |
| Writer | migration-controller (VS0-WRITER-020) |

For the `BackupVerified` milestone only, `record.signedBackupEvidenceRef` and `record.restoreVerificationEvidenceRef` are required evidence references. They prove the signed backup and verified restore gates from ADH-2026-041 without selecting a cryptographic algorithm or implementing a backup service in Phase 2R.

### 7.4 Draft Representation (ADH-2026-043 decision 2)

`Draft` is pre-persistence preparation of a migration plan. The signed immutable CanonicalMigrationPlan precedes the persisted CanonicalMigrationRecord sequence. The 1..9 milestone evidence sequence begins at `InventoryValidated`. Draft has no representation in the append-only record chain or in the immutable record state machine (VS0-STATE-010/011).

### 7.5 Migration Inventory and Feature Boundary (ADH-2026-043 decision 3)

FEATURE-0015 owns the signed global inventory/classification plan and executes provider/topology migration only. The following table classifies every alpha record family for FEATURE-0015:

| Alpha Record Family | F0015 Action | Notes |
|---------------------|--------------|-------|
| Provider (combined) | Classifies and transforms now | Splits into CloudPlatform + CloudProvider |
| ProviderLocation | Classifies and transforms now | Maps to HostingLocation |
| ProviderDatacenter | Classifies and transforms now | Maps to Datacenter |
| DatacenterFailureDomain | Classifies and transforms now | Maps to FaultDomain |
| InfrastructureStack | Classifies and transforms now | Identity preserved, new scope |
| ExecutionTarget (identity) | Classifies and transforms now | Spec only; status deferred to FEATURE-0016 |
| ServiceClass | Classifies in signed plan; defers transform | FEATURE-0022 owns ServiceTypeDefinition/ServiceOffering |
| CloudEnrollment (alpha) | Classifies in signed plan; defers transform | FEATURE-0021 owns enrollment |
| Governance/policy alpha refs | Classifies in signed plan; defers transform | FEATURE-0018/0019/0020 own transforms |
| Placement/decision alpha refs | Classifies in signed plan; defers transform | FEATURE-0023 owns transforms |
| Plugin/operation alpha refs | Classifies in signed plan; defers transform | FEATURE-0024 owns transforms |

FEATURE-0015 appends only the `provider-topology` plan-declared run and reaches a run-local `Completed` for that run only. FEATURE-0015 cannot claim global cutover completion; each later domain owner appends only its assigned plan-declared run, and final all-domain completion, cutover, and conformance are proven only by FEATURE-0026 integration evidence after every required run is sealed (ADH-2026-044).

### 7.6 Audit, Correlation, and Redaction Requirements (ADH-2026-043 decision 6)

FEATURE-0015 reuses FEATURE-0013 AuditEvent. The following lifecycle and security events produce audit evidence:

| Event | Audit Evidence | Correlation |
|-------|----------------|-------------|
| CloudProviderParticipation state change | AuditEvent with subjectRef, actor, transition, timestamp | participationRef UID |
| CanonicalMigrationPlan publication | AuditEvent with planRef, publisher actor, signedAt | planRef UID |
| Each migration milestone record append | AuditEvent with milestoneRef, stage, predecessorRef | planRef UID + milestone ordinal |
| Migration cutover decision | AuditEvent with cutover activation evidence | planRef UID |
| Cross-provider safe-denial (security) | AuditEvent with denied actor, denied action; no target existence disclosed | requestId |

Rules:
- Safe projection/redaction: audit events use the same safe-denial principle — no existence disclosure for cross-provider references.
- No secrets: no credential values, protected handles, or secret material in any audit record.
- Intentional correlation: each audit record links to its subject, actor, and governing plan or participation through UID-pinned references.
- FEATURE-0015 does not import FEATURE-0026's integration-only trace conformance (VS0-CF-T01).

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

| Owned Item | DEC/ADH | VS0 Schema | VS0 Writer | VS0 State | VS0 Conformance (FEATURE-0015 local) |
|------------|---------|------------|------------|-----------|--------------------------------------|
| CloudPlatform | DEC-0037; ADH-020/042/043 | VS0-SCHEMA-008 | VS0-WRITER-002 | — | VS0-CF-F15-01,F15-02,F15-11 |
| CloudProvider | DEC-0037; ADH-020/042/043 | VS0-SCHEMA-009 | VS0-WRITER-003 | — | VS0-CF-F15-01,F15-02 |
| CloudProviderParticipation | DEC-0054; ADH-037/042/043 | VS0-SCHEMA-010 | VS0-WRITER-004 | VS0-STATE-001 | VS0-CF-F15-01,F15-02,F15-03,F15-04,F15-09,F15-11 |
| HostingLocation | DEC-0041; ADH-024/042/043 | VS0-SCHEMA-011 | VS0-WRITER-003 | — | VS0-CF-F15-01,F15-02 |
| Datacenter | DEC-0041; ADH-024/042/043 | VS0-SCHEMA-012 | VS0-WRITER-003 | — | VS0-CF-F15-01,F15-02 |
| FaultDomain | DEC-0041; ADH-024/042/043 | VS0-SCHEMA-013 | VS0-WRITER-003 | — | VS0-CF-F15-01,F15-02 |
| InfrastructureStack | DEC-0041,0042; ADH-025/042/043 | VS0-SCHEMA-014 | VS0-WRITER-003 | — | VS0-CF-F15-01,F15-02 |
| ExecutionTarget (identity) | DEC-0042; ADH-025/042/043 | VS0-SCHEMA-015 | VS0-WRITER-003 | — | VS0-CF-F15-01,F15-02,F15-05,F15-06 |
| CanonicalMigrationPlan | DEC-0058; ADH-041/042/043 | VS0-SCHEMA-060 | VS0-WRITER-021 | VS0-STATE-010 | VS0-CF-MIG01,F15-07,F15-10 |
| CanonicalMigrationRecord | DEC-0058; ADH-041/042/043 | VS0-SCHEMA-061 | VS0-WRITER-020 | VS0-STATE-010,011 | VS0-CF-MIG02,F15-08,MIGF01..MIGF03 |
| Cross-provider isolation | DEC-0037,0054; ADH-043 | VS0-SCHEMA-015 | VS0-WRITER-003 | — | VS0-CF-X03 |
| Scope-reference integrity | DEC-0037,0054; ADH-043 | VS0-SCHEMA-008,010 | VS0-WRITER-002,004 | — | VS0-CF-F15-11 |

Note: VS0-CF-HP01 (FEATURE-0026 integration) and VS0-CF-F09 (FEATURE-0023 placement) are downstream integration references only; they are not FEATURE-0015 local acceptance evidence.

---

## 11. Anti-Drift Rules

1. FEATURE-0015 requirements/design/tasks must NOT introduce or use any superseded concept listed in §8 as active behavior. Those names may appear only in the explicitly labelled exclusion/non-goal context that explains why they are prohibited.
2. FEATURE-0015 must NOT introduce, initialize, persist, validate, default, or write ExecutionTarget `status.*` fields — FEATURE-0016 owns their introduction and activation.
3. FEATURE-0015 must NOT evaluate provider-selection intent — that is FEATURE-0021 scope.
4. FEATURE-0015 must NOT implement platform lifecycle (SovrunnInstallation, SovrunnRelease, PlatformLifecyclePolicy, PlatformLifecyclePlan).
5. Migration state machine (VS0-STATE-011) is the accepted canonical sequence from DEC-0058/data-model §16.1; FEATURE-0015 must not invent additional states.
6. The migration controller may append CanonicalMigrationRecord evidence under an approved signed plan; it must never author, approve, replace, or mutate CanonicalMigrationPlan.

---

*End of architecture boundary.*
