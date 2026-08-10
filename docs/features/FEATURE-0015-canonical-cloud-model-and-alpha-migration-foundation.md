# FEATURE-0015: Canonical Cloud Model and Alpha Migration Foundation

| Field | Value |
|-------|-------|
| Status | Approved scope (correction under ADH-2026-042) |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase | 2R |
| Order | 5 (first unimplemented Phase 2R feature) |
| Depends On | FEATURE-0011, FEATURE-0012, FEATURE-0013, FEATURE-0014 |
| Depended On By | FEATURE-0016, FEATURE-0021, FEATURE-0022 |
| Architecture Boundary | docs/architecture/FEATURE-0015-canonical-cloud-model-and-alpha-migration-foundation.md |
| Controlling Decisions | DEC-0037, DEC-0041, DEC-0042, DEC-0054, DEC-0058 |
| Controlling Handoffs | ADH-2026-020, ADH-2026-024, ADH-2026-025, ADH-2026-037, ADH-2026-041, ADH-2026-042 |

---

## 1. Feature Summary

Establish the seven-scope canonical cloud model identity layer and one-way alpha migration foundation. Register CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, InfrastructureStack, and ExecutionTarget (identity only) as managed resources. Implement CanonicalMigrationPlan and CanonicalMigrationRecord to prove deterministic dry-run cutover from alpha to canonical model with no dual write/authority.

---

## 2. Scope Classification

### 2.1 REQUIREMENTS-Owned Observable Behavior

| ID | Behavior | DEC/ADH | VS0 IDs |
|----|----------|---------|---------|
| REQ-F15-01 | Register CloudPlatform with owner Organization reference; validate uniqueness, scope, immutable identity | DEC-0037; ADH-020 | VS0-SCHEMA-008 |
| REQ-F15-02 | Register CloudProvider with operating markets; validate scope and identity | DEC-0037; ADH-020 | VS0-SCHEMA-009 |
| REQ-F15-03 | Register CloudProviderParticipation linking one CloudPlatform and one CloudProvider; enforce uniqueness per platform+provider pair; lifecycle Pending→Active→Suspended→Terminating→Terminated | DEC-0054; ADH-037 | VS0-SCHEMA-010, VS0-STATE-001 |
| REQ-F15-04 | Register topology chain: HostingLocation → Datacenter → FaultDomain → InfrastructureStack under CloudProvider scope; validate containment refs | DEC-0041; ADH-024 | VS0-SCHEMA-011..014 |
| REQ-F15-05 | Register ExecutionTarget identity (spec fields only: infrastructureStackRef, participationRef, targetClass) with cross-reference validation; do NOT introduce, persist, or initialize any status fields (qualification, availability, maintenanceEpoch, factSetRef, observedGeneration, conditions) — those are introducedBy and activatedBy FEATURE-0016 | DEC-0042; ADH-025 | VS0-SCHEMA-015 |
| REQ-F15-06 | Create signed CanonicalMigrationPlan with classification mappings; enforce immutability after signature (FINAL on persistence per VS0-STATE-010); corrections create a new linked superseding plan | DEC-0058; ADH-041 | VS0-SCHEMA-060, VS0-STATE-010 |
| REQ-F15-07 | Append CanonicalMigrationRecord for each migration milestone; each record is FINAL and immutable on persistence; records form an ordered milestone sequence linked to their plan and predecessor | DEC-0058; ADH-041 | VS0-SCHEMA-061, VS0-STATE-011 |
| REQ-F15-08 | Deterministic dry-run: convert alpha Provider/topology fixtures to canonical representations; prove signed backup and verified restore; report zero unresolved authorities, zero ambiguous references, no dual writers | DEC-0058; ADH-041 | VS0-CF-MIG01, VS0-CF-F15-08 |
| REQ-F15-09 | Cross-provider isolation: deny cross-provider target reference; use safe RESOURCE_NOT_FOUND without existence disclosure | DEC-0037,0054 | VS0-CF-X03 |
| REQ-F15-10 | Use only the seven canonical scope kinds globally; each FEATURE-0015 resource accepts only its registry-declared scope subset | DEC-0037 | VS0-SCHEMA-001,008..015,060..061 |
| REQ-F15-11 | Writer enforcement: only cloud-provider-admin may write topology spec; only cloud-platform-admin may write CloudPlatform spec | DEC-0037 | VS0-WRITER-002,003 |
| REQ-F15-12 | Existing FEATURE-0012 Problem Details codes used for all errors; no new top-level codes | — | VS0-SCHEMA-004 |
| REQ-F15-13 | Separate migration authorities: approved-migration-plan-publisher alone persists signed CanonicalMigrationPlan; migration-controller alone appends CanonicalMigrationRecord under that approved plan | DEC-0058; ADH-041 | VS0-WRITER-020,021 |

### 2.2 DESIGN-Delegated Mechanics

- In-memory registry data structures for cloud model resources
- Handler wiring and router setup
- Validation plumbing (deterministic, pure)
- Migration milestone append validator and controller internals
- UID generation strategy for split transforms
- Dry-run report format

### 2.3 Downstream Fields Visible Only as Negative Boundaries

These fields appear in the shared Slice 0 registry only so FEATURE-0015 can prove it does not activate them. They are not FEATURE-0015 requirements or persistence fields.

| Element | Introduced and Activated By | FEATURE-0015 Rule |
|---------|-----------------------------|------------------|
| ExecutionTarget `status.qualification`, `status.availability`, `status.maintenanceEpoch`, `status.factSetRef`, `status.observedGeneration`, `status.conditions` | FEATURE-0016 | Must not initialize, store, default, validate, or write |

### 2.4 EXCLUDED (Must NOT Implement)

| Concept | Reason | Reference |
|---------|--------|-----------|
| ResourcePool | Superseded | DEC-0042 |
| ProviderCapability | Superseded | DEC-0042 |
| Generic Provider (combined) | Superseded | DEC-0037 |
| SovrunnInstallation | No Phase 2R owner | DEC-0053 deferred |
| CloudProviderParticipation `spec.providerSelectionModes` | Introduced and activated only by FEATURE-0021 | VS0-SCHEMA-010 field ownership |
| Platform lifecycle resources | No Phase 2R owner | DEC-0053 deferred |
| Target qualification writes | FEATURE-0016 | DEC-0042 |
| NormalizedTargetFactSet | FEATURE-0016 | DEC-0036 |
| CloudEnrollment | FEATURE-0021 | DEC-0038 |
| Provider-selection intent evaluation | FEATURE-0021 | — |
| ServiceClass | Superseded | DEC-0049 |
| EffectivePolicyContext | Superseded | DEC-0050 |
| Six-scope vocabulary | Superseded | DEC-0037 |
| Real provisioning | Phase 3 | PHASE2R non-goals |

---

## 3. Acceptance Criteria

| ID | Criterion | Conformance ID |
|----|-----------|----------------|
| AC-F15-01 | CloudPlatform, CloudProvider, topology chain CRUD works with validation and their registry-declared scope subsets | VS0-CF-F15-01, VS0-CF-F15-02 |
| AC-F15-02 | CloudProviderParticipation lifecycle, dual delegated acceptance guard, and pair uniqueness are enforced | VS0-CF-F15-03, VS0-CF-F15-04, VS0-CF-F15-09 |
| AC-F15-03 | ExecutionTarget identity (spec) registered; no status fields introduced or persisted by FEATURE-0015 | VS0-CF-F15-01 |
| AC-F15-04 | Cross-provider target reference denied safely | VS0-CF-X03 |
| AC-F15-05 | Migration dry-run passes deterministically with zero unresolved authorities | VS0-CF-MIG01 |
| AC-F15-06 | Migration records are append-only, FINAL, and linked to plan with milestone/predecessor chain | VS0-CF-MIG02 |
| AC-F15-07 | Every output passes the canonical stale-concept anti-drift gate | Anti-drift |
| AC-F15-08 | Writer enforcement: unauthorized writes to topology/platform specs, status, and system-owned fields denied; cross-provider references deny safely | VS0-CF-F15-05, VS0-CF-F15-06, VS0-CF-X03 |
| AC-F15-09 | All errors use FEATURE-0012 Problem Details with existing codes | — |
| AC-F15-10 | Invalid milestone order, dual authority and unresolved migrated references fail closed with exact registered codes/violations and no cutover side effects | VS0-CF-MIGF01..MIGF03 |
| AC-F15-11 | Migration controller cannot create or mutate CanonicalMigrationPlan; plan publisher cannot append CanonicalMigrationRecord; plan mutation is rejected | VS0-CF-MIG01,MIG02,VS0-CF-F15-07,VS0-CF-F15-10 |

---

## 4. Migration Milestone Sequence (VS0-STATE-011)

VS0-STATE-011 represents the **CanonicalMigrationRecord append-only milestone sequence**. The CanonicalMigrationPlan itself is a signed immutable record (FINAL per VS0-STATE-010). Migration progress is proven by appending successive FINAL CanonicalMigrationRecord instances, each linked to the plan and its predecessor record.

**Milestones** (each is a new FINAL CanonicalMigrationRecord appended to the chain):

```text
InventoryValidated
  → DryRunPassed
  → WriteFrozen
  → BackupVerified
  → Transformed
  → ReferencesVerified
  → CutoverActivated
  → ConformancePassed
  → Completed
```

**CanonicalMigrationRecord milestone fields:**

| Field | Description |
|-------|-------------|
| `record.milestone` | Enum: InventoryValidated, DryRunPassed, WriteFrozen, BackupVerified, Transformed, ReferencesVerified, CutoverActivated, ConformancePassed, Completed |
| `record.stage` | Integer (1..9) corresponding to milestone ordinal position |
| `record.predecessorRef` | TypedRef<CanonicalMigrationRecord> (nil for first milestone; uid-pinned for subsequent) |

**Guards:** Each milestone record may only be appended if the preceding milestone record (by stage ordinal) exists and is FINAL for the same plan. `InventoryValidated` requires zero unclassified alpha records. `DryRunPassed` requires zero transform errors. `WriteFrozen` requires zero active legacy writers. `BackupVerified` requires both `record.signedBackupEvidenceRef` and `record.restoreVerificationEvidenceRef`; each is an opaque, externally verifiable evidence reference, and no cryptographic algorithm is selected by this feature. Missing either is rejected with VALIDATION_FAILED/422 and `VS0_MIGRATION_BACKUP_RESTORE_UNVERIFIED`. Invalid milestone order is rejected with CONFLICT/409 and `VS0_MIGRATION_STATE_INVALID`.

**Immutability:** Neither the plan nor any prior record mutates. Progress is proven exclusively by new append-only records. Corrections create linked records with a supersession reference, never mutation.

---

## 5. Error Codes Used

| Scenario | Code | HTTP | Violation |
|----------|------|------|-----------|
| Missing/invalid auth | AUTH_REQUIRED | 401 | VS0_AUTH_REQUIRED |
| Unauthorized write | AUTHORIZATION_DENIED | 403 | — |
| Cross-provider ref denied | RESOURCE_NOT_FOUND | 404 | VS0_AUTHORIZATION_SAFE_DENIAL |
| Duplicate name in scope | ALREADY_EXISTS | 409 | — |
| Invalid participation transition | CONFLICT | 409 | VS0_PARTICIPATION_STATE_INVALID |
| Invalid migration milestone append | CONFLICT | 409 | VS0_MIGRATION_STATE_INVALID |
| BackupVerified missing signed backup or verified restore evidence | VALIDATION_FAILED | 422 | VS0_MIGRATION_BACKUP_RESTORE_UNVERIFIED |
| Duplicate CloudProviderParticipation pair | ALREADY_EXISTS | 409 | VS0_PARTICIPATION_DUPLICATE |
| Missing delegated acceptance for Pending → Active | CONFLICT | 409 | VS0_PARTICIPATION_ACCEPTANCE_MISSING |
| Client writes status | AUTHORIZATION_DENIED | 403 | VS0_STATUS_FIELD_WRITE |
| Client writes system-owned metadata | AUTHORIZATION_DENIED | 403 | VS0_SYSTEM_OWNED_FIELD_WRITE |
| Persisted migration plan is mutated | CONFLICT | 409 | VS0_MIGRATION_PLAN_IMMUTABLE |
| Plan publisher or migration controller writes the other's record kind | ALREADY_EXISTS | 409 | VS0_MIGRATION_WRITER_SEPARATION |
| Resource uses a scope outside its declared subset | VALIDATION_FAILED | 422 | VS0_SCOPE_KIND_INVALID |
| Legacy and canonical desired-state authority both active | CONFLICT | 409 | VS0_MIGRATION_DUAL_AUTHORITY |
| Unresolved canonical reference after transform | VALIDATION_FAILED | 422 | VS0_MIGRATION_UNRESOLVED_REF |
| Invalid reference chain | VALIDATION_FAILED | 422 | — |
| Stale resourceVersion | STALE_RESOURCE_VERSION | 412 | — |

---

## 6. Phase 2R Exit Evidence (FEATURE-0015 contribution)

1. Deterministic dry-run conversion of Provider/topology fixtures produces zero ambiguous authorities.
2. Old and new writers never coexist after cutover activation.
3. Provider credential isolation proven by cross-provider denial tests.
4. Seven canonical scope kinds used exclusively.
5. No ResourcePool, ProviderCapability, or generic Provider in active code/schemas.

---

## 7. Next Feature Boundary

After FEATURE-0015 completes, FEATURE-0016 may introduce and activate ExecutionTarget qualification, availability, and maintenance semantics. Until then, FEATURE-0015 must not persist or otherwise implement those fields.

---

*Executable scope only. No requirements.md, design.md, or tasks.md generation from this file. Those are produced by Kiro spec stages under .kiro/specs/.*
