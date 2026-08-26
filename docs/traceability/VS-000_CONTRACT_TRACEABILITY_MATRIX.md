# VS-000 Slice 0 Contract Traceability Matrix

| Field | Value |
|-------|-------|
| Status | Experimental |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Controlling Adoption | ADH-2026-042 |
| Machine Source | `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` |

---

## Authority Chain

| Priority | Authority | Path |
|----------|-----------|------|
| 1 | Active architecture baseline | `docs/context/CURRENT_ARCHITECTURE_BASELINE.md` |
| 2 | Canonical model and catalog | `docs/architecture/canonical/` |
| 3 | Accepted decisions and controlling handoffs | `docs/decisions/DECISION_INDEX.md`; ADH-2026-042; ADH-2026-064 (ExecutionTarget ownership correction) |
| 4 | VS-000 core skeleton | `docs/architecture/vertical-slices/VS-000-core-skeleton.md` |
| 5 | VS-000 contract specification | `docs/architecture/vertical-slices/VS-000-contract-specification.md` |
| 6 | VS-000 contract registry YAML | `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` |
| 7 | Approved feature stage and implementation | `.kiro/specs/`; `internal/` packages |

---

## Schema Mapping (VS0-SCHEMA-001..061; VS0-SCHEMA-057/060/061 retired tombstones)

| Registry ID | Contract Identity | Owner | Controlling DEC/ADH | Repository Authority | Conformance |
|-------------|-------------------|-------|---------------------|---------------------|-------------|
| VS0-SCHEMA-001 | F12:ObjectMeta | FEATURE-0012 | DEC-0026/0027/0036; ADH-2026-012/013 | FEATURE-0012 schema | VS0-CF-HP01 |
| VS0-SCHEMA-002 | F12:TypedRef | FEATURE-0012 | DEC-0026/0027/0036; ADH-2026-012/013 | FEATURE-0012 schema | VS0-CF-HP01 |
| VS0-SCHEMA-003 | F12:Condition | FEATURE-0012 | DEC-0026/0027/0036; ADH-2026-012/013 | FEATURE-0012 schema | VS0-CF-HP01 |
| VS0-SCHEMA-004 | F12:Problem | FEATURE-0012 | DEC-0026/0027/0036; ADH-2026-012/013 | FEATURE-0012 schema | VS0-CF-F01..F20 |
| VS0-SCHEMA-005 | DecisionProfile | FEATURE-0013 | DEC-0026/0036; ADH-2026-017 | FEATURE-0013 schema | VS0-CF-F07..F09 |
| VS0-SCHEMA-006 | DecisionRecord | FEATURE-0013 | DEC-0026/0036; ADH-2026-017 | FEATURE-0013 schema | VS0-CF-HP01,T01 |
| VS0-SCHEMA-007 | AuditEvent | FEATURE-0013 | DEC-0026/0036; ADH-2026-017 | FEATURE-0013 schema | VS0-CF-HP01,T01 |
| VS0-SCHEMA-008 | CloudPlatform | FEATURE-0015 | DEC-0037, ADH-2026-042/045/046/047/048/050/051/052 | canonical data model | VS0-CF-F15-01,02,07,12,13,14,18,19,25,26,27,29,30,31,32,33 |
| VS0-SCHEMA-009 | CloudProvider | FEATURE-0015 | DEC-0037, ADH-2026-042/045/046/047/048/050/051/052 | canonical data model | VS0-CF-F15-01,02,07,12,13,14,18,19,25,26,27,29,30,31,32,33 |
| VS0-SCHEMA-010 | CloudProviderParticipation | FEATURE-0015 | DEC-0054; ADH-2026-037/042/045/046/047/048/050/051 | canonical data model | VS0-CF-F15-01..04,08..10,15..21,25,26,27,28,30,31 |
| VS0-SCHEMA-011 | HostingLocation | FEATURE-0015 | DEC-0041; ADH-2026-024/042/045/046/047/048/050/051/052 | canonical data model | VS0-CF-F15-01,02,07,12,13,14,18,19,25,26,27,29,30,31,33 |
| VS0-SCHEMA-012 | Datacenter | FEATURE-0015 | DEC-0041; ADH-2026-024/042/045/046/047/050/051/052 | canonical data model | VS0-CF-F15-01,02,07,12,13,14,18,19,25,26,27,30,31,33 |
| VS0-SCHEMA-013 | FaultDomain | FEATURE-0015 | DEC-0041; ADH-2026-024/042/045/046/047/050/051/052 | canonical data model | VS0-CF-F15-01,02,07,12,13,14,18,19,25,26,27,30,31,33 |
| VS0-SCHEMA-014 | InfrastructureStack | FEATURE-0015 | DEC-0041/0042; ADH-2026-025/042/045/046/047/050/051/052 | canonical data model | VS0-CF-F15-01,02,07,12,13,14,18,19,25,26,27,30,31,33 |
| VS0-SCHEMA-015 | ExecutionTarget | FEATURE-0016 | DEC-0042/0057; ADH-2026-025/040/042/045/058 | canonical data model | VS0-CF-F10,X03,F16-01..122 |
| VS0-SCHEMA-016 | NormalizedTargetFactSet | FEATURE-0016 | DEC-0036/0042/0057; ADH-2026-040/042/058 | canonical data model | VS0-CF-F16-20..23,67..69,90 |
| VS0-SCHEMA-017 | TargetQualificationResult | FEATURE-0016 | DEC-0036/0042/0057; ADH-2026-040/042/058 | canonical data model | VS0-CF-F16-20..23,67..69,90 |
| VS0-SCHEMA-018 | PolicyEvaluationRequest | FEATURE-0017 | DEC-0028/0043; ADH-2026-026/042/067/068/069 | `docs/architecture/policy-evaluation-abstraction.md` | FEATURE-0017 architecture §11, cases 1–24 |
| VS0-SCHEMA-019 | PolicyEvaluationResult | FEATURE-0017 | DEC-0028/0043; ADH-2026-026/042/067/068/069 | `docs/architecture/policy-evaluation-abstraction.md` | FEATURE-0017 architecture §11, cases 1–24 |
| VS0-SCHEMA-020 | PrincipalRef | FEATURE-0018 | DEC-0050; ADH-2026-033/042 | canonical data model | VS0-CF-F01,F02 |
| VS0-SCHEMA-021 | Membership | FEATURE-0018 | DEC-0050; ADH-2026-033/042 | canonical data model | VS0-CF-F01,F02 |
| VS0-SCHEMA-022 | RoleDefinition | FEATURE-0018 | DEC-0050; ADH-2026-033/042 | canonical data model | VS0-CF-F01,F02 |
| VS0-SCHEMA-023 | RoleAssignment | FEATURE-0018 | DEC-0050; ADH-2026-033/042 | canonical data model | VS0-CF-F01,F02,X01..X03 |
| VS0-SCHEMA-024 | GovernanceProfile | FEATURE-0018 | DEC-0050, ADH-2026-042 | canonical data model | VS0-CF-F07 |
| VS0-SCHEMA-025 | SovereigntyProfile | FEATURE-0019 | DEC-0041/0043/0055; ADH-2026-024/026/038/042 | canonical data model | VS0-CF-F08 |
| VS0-SCHEMA-026 | RegulatoryPolicyBundle | FEATURE-0019 | DEC-0041/0043/0055; ADH-2026-024/026/038/042 | canonical data model | VS0-CF-F08 |
| VS0-SCHEMA-027 | SovereigntyFactSet | FEATURE-0019 | DEC-0041/0043/0055; ADH-2026-024/026/038/042 | canonical data model | VS0-CF-F08 |
| VS0-SCHEMA-028 | EvidenceRecord | FEATURE-0019 | DEC-0041/0043/0055; ADH-2026-024/026/038/042 | canonical data model | VS0-CF-F08,L01 |
| VS0-SCHEMA-029 | ProfileAssignment | FEATURE-0020 | DEC-0050, ADH-2026-042 | canonical data model | VS0-CF-F07 |
| VS0-SCHEMA-030 | EffectiveGovernanceContext | FEATURE-0020 | DEC-0050, ADH-2026-042 | canonical data model | VS0-CF-HP01,F07 |
| VS0-SCHEMA-031 | CloudEnrollment | FEATURE-0021 | DEC-0038/0039/0047; ADH-2026-021/022/030/037/042 | canonical data model | VS0-CF-F03,F04 |
| VS0-SCHEMA-032 | EntitlementPackage | FEATURE-0021 | DEC-0038/0039/0047; ADH-2026-021/022/030/037/042 | canonical data model | VS0-CF-F04 |
| VS0-SCHEMA-033 | ServiceEntitlement | FEATURE-0021 | DEC-0038/0039/0047; ADH-2026-021/022/030/037/042 | canonical data model | VS0-CF-F04 |
| VS0-SCHEMA-034 | QuotaPolicy | FEATURE-0021 | DEC-0038/0039/0047; ADH-2026-021/022/030/037/042 | canonical data model | VS0-CF-F05 |
| VS0-SCHEMA-035 | QuotaReservationRequest | FEATURE-0021 | DEC-0038/0039/0047; ADH-2026-021/022/030/037/042 | canonical data model | VS0-CF-F05,I01,I02 |
| VS0-SCHEMA-036 | ServicePortfolio | FEATURE-0022 | DEC-0040/0044/0049; ADH-2026-023/027/032/042 | canonical data model | VS0-CF-HP01 |
| VS0-SCHEMA-037 | ServiceTypeDefinition | FEATURE-0022 | DEC-0040/0044/0049; ADH-2026-023/027/032/042 | canonical data model | VS0-CF-HP01,F06 |
| VS0-SCHEMA-038 | ServiceOffering | FEATURE-0022 | DEC-0040/0044/0049; ADH-2026-023/027/032/042 | canonical data model | VS0-CF-HP01,F06 |
| VS0-SCHEMA-039 | ServicePlan | FEATURE-0022 | DEC-0040/0044/0049; ADH-2026-023/027/032/042 | canonical data model | VS0-CF-HP01,F06 |
| VS0-SCHEMA-040 | ServicePlacementProfile | FEATURE-0022 | DEC-0040/0044/0049; ADH-2026-023/027/032/042 | canonical data model | VS0-CF-F09,F11 |
| VS0-SCHEMA-041 | ServiceRuntimeProfile | FEATURE-0022 | DEC-0040/0044/0049; ADH-2026-023/027/032/042 | canonical data model | VS0-CF-HP01 |
| VS0-SCHEMA-042 | ServiceRequirementSet | FEATURE-0022 | DEC-0040/0044/0049; ADH-2026-023/027/032/042 | canonical data model | VS0-CF-HP01 |
| VS0-SCHEMA-043 | SovereigntyDecisionProfile | FEATURE-0023 | DEC-0043/0045/0046; ADH-2026-026/028/029/042 | canonical data model | VS0-CF-F08 |
| VS0-SCHEMA-044 | PlacementDecisionProfile | FEATURE-0023 | DEC-0043/0045/0046; ADH-2026-026/028/029/042 | canonical data model | VS0-CF-F09,F11 |
| VS0-SCHEMA-045 | ServicePlacement | FEATURE-0023 | DEC-0043/0045/0046; ADH-2026-026/028/029/042 | canonical data model | VS0-CF-HP01,F11 |
| VS0-SCHEMA-046 | ServiceDeploymentPlan | FEATURE-0024 | DEC-0029; ADH-2026-042 | canonical data model | VS0-CF-HP01,F14 |
| VS0-SCHEMA-047 | Operation | FEATURE-0012 | DEC-0026/0027/0036; ADH-2026-012/013/042 | FEATURE-0012 schema | VS0-CF-HP01,F14..F16 |
| VS0-SCHEMA-048 | PluginExecution | FEATURE-0024 | DEC-0029; ADH-2026-042 | canonical data model | VS0-CF-HP01,F15,F16,Z01 |
| VS0-SCHEMA-049 | DecisionOperationContext | FEATURE-0025 | DEC-0043; ADH-2026-026/042 | canonical data model | VS0-CF-F19 |
| VS0-SCHEMA-050 | ServiceInstance | FEATURE-0007 | ADH-2026-042 Slice 0 adoption | retained Phase 1 contract | VS0-CF-HP01,F03..F16,D01 |
| VS0-SCHEMA-051 | SecretRef | FEATURE-0008 | DEC-0051; ADH-2026-034/042 | canonical data model | VS0-CF-HP01,F17,L01 |
| VS0-SCHEMA-052 | ServiceBinding | FEATURE-0008+FEATURE-0024 | DEC-0051; ADH-2026-034/042 | canonical data model | VS0-CF-HP01,F17,F20,D01 |
| VS0-SCHEMA-053 | Organization | FEATURE-0002 | FEATURE-0012 grammar; ADH-2026-042 Slice 0 adoption | retained Phase 1 contract | VS0-CF-HP01,X01 |
| VS0-SCHEMA-054 | Tenant | FEATURE-0005 | FEATURE-0012 grammar; ADH-2026-042 Slice 0 adoption | retained Phase 1 contract | VS0-CF-HP01,X01 |
| VS0-SCHEMA-055 | Project | FEATURE-0006 | FEATURE-0012 grammar; ADH-2026-042 Slice 0 adoption | retained Phase 1 contract | VS0-CF-HP01,X01,X02 |
| VS0-SCHEMA-056 | ServiceRegion | FEATURE-0022 | DEC-0037/0041/0054; ADH-2026-024/037/042 | canonical data model | VS0-CF-HP01,F04,F06 |
| VS0-SCHEMA-057 | SovrunnInstallation | **Removed-NoPhase2ROwner** | DEC-0054; ADH-2026-037/042 | No Phase 2R feature owns platform lifecycle; removed 2026-08-04 | — |
| VS0-SCHEMA-058 | PluginDefinition | FEATURE-0012 retained | DEC-0029; ADH-2026-042 Slice 0 adoption | FEATURE-0012 schema | VS0-CF-HP01,Z01 |
| VS0-SCHEMA-059 | AdapterConfiguration | FEATURE-0012 retained | DEC-0036; ADH-2026-042 Slice 0 adoption | FEATURE-0012 schema | VS0-CF-HP01,Z01 |
| VS0-SCHEMA-060 | CanonicalMigrationPlan | **Retired-CanonicalBootstrap** | DEC-0059 supersedes DEC-0058; ADH-2026-045 | No live alpha state to convert; permanently retired tombstone | — |
| VS0-SCHEMA-061 | CanonicalMigrationRecord | **Retired-CanonicalBootstrap** | DEC-0059 supersedes DEC-0058; ADH-2026-045 | No live alpha state to convert; permanently retired tombstone | — |

---

## Writer Mapping (VS0-WRITER-001..019; VS0-WRITER-020/021 retired tombstones)

| Registry ID | Writer Domain | Owner | Authority | Conformance Evidence |
|-------------|---------------|-------|-----------|---------------------|
| VS0-WRITER-001 | api-server | FEATURE-0012 | canonical contract catalog | VS0-CF-HP01 |
| VS0-WRITER-002 | cloud-platform-admin-or-product-publisher | FEATURE-0015+FEATURE-0022 | DEC-0037, DEC-0049 | VS0-CF-F15-01,05,06,07 |
| VS0-WRITER-003 | cloud-provider-admin | FEATURE-0015 | DEC-0037 | VS0-CF-F15-01,05,06,07,X03 |
| VS0-WRITER-004 | api-server (participation create/actions; scheduler invokes api-server expiry transition, not a separate writer) | FEATURE-0015 | DEC-0037, ADH-2026-045/046/047 | VS0-CF-F15-03,04,08,09,10,15,16,17,18,19,20,28 |
| VS0-WRITER-005 | kind-registered-controller (F0015 resolves explicitly to api-server; ADH-2026-046 decision 1) | FEATURE-0012 | canonical contract catalog | VS0-CF-HP01,F15,F16 |
| VS0-WRITER-006 | ExecutionTargetLifecycleService (sole committer of ExecutionTarget status and internal FactSet/Result records) | FEATURE-0016 | DEC-0036/0042/0057; ADH-2026-058 | VS0-CF-F10,Z01,F16-01..122 |
| VS0-WRITER-007 | authorized-governance-publisher | FEATURE-0018+FEATURE-0019 | DEC-0041/0043/0050/0055 | VS0-CF-HP01,F07 |
| VS0-WRITER-008 | delegated-identity-governance-admin | FEATURE-0018 | DEC-0050 | VS0-CF-F01,F02,X01..X03 |
| VS0-WRITER-009 | evidence-collector | FEATURE-0019 | DEC-0041/0043/0055 | VS0-CF-F08,L01 |
| VS0-WRITER-010 | deterministic-resolver | FEATURE-0020 | DEC-0050 | VS0-CF-HP01,F07 |
| VS0-WRITER-011 | authorized-decision-or-audit-producer | FEATURE-0013+FEATURE-0023 | ADH-2026-017; DEC-0043 | VS0-CF-HP01,T01 |
| VS0-WRITER-012 | deployment-planner | FEATURE-0024 | DEC-0029 | VS0-CF-HP01,F14 |
| VS0-WRITER-013 | operation-api-on-acceptance | FEATURE-0012 | canonical contract catalog | VS0-CF-HP01,F14 |
| VS0-WRITER-014 | operation-controller | FEATURE-0012 | canonical contract catalog | VS0-CF-HP01,F15,F16 |
| VS0-WRITER-015 | operation-controller (PluginExecution spec) | FEATURE-0024 | DEC-0029 | VS0-CF-HP01,Z01 |
| VS0-WRITER-016 | fake-plugin-executor | FEATURE-0024 | DEC-0029 | VS0-CF-HP01,F15,F16,Z01 |
| VS0-WRITER-017 | projection-controller | FEATURE-0023+FEATURE-0025 | DEC-0043/0046 | VS0-CF-HP01,F18,F19 |
| VS0-WRITER-018 | authorized-project-consumer | FEATURE-0007+FEATURE-0008 | DEC-0051; ADH-2026-042 Slice 0 adoption | VS0-CF-HP01,F03..F06,I01,I02 |
| VS0-WRITER-019 | binding-controller | FEATURE-0008+FEATURE-0024 | DEC-0051 | VS0-CF-HP01,F17 |
| VS0-WRITER-020 | **Retired-CanonicalBootstrap** (was migration-controller) | — | DEC-0059 supersedes DEC-0058; ADH-2026-045 | — |
| VS0-WRITER-021 | **Retired-CanonicalBootstrap** (was approved-migration-plan-publisher) | — | DEC-0059 supersedes DEC-0058; ADH-2026-045 | — |

---

## State Machine Mapping (VS0-STATE-001..010; VS0-STATE-011 retired tombstone)

| Registry ID | Kind | Owner | Authority | Conformance Evidence |
|-------------|------|-------|-----------|---------------------|
| VS0-STATE-001 | CloudProviderParticipation (Pending/Active/Rejected/Withdrawn/Expired/Suspended/Terminating/Terminated; independent platformSuspended/providerSuspended holds) | FEATURE-0015 | DEC-0037,0054; ADH-2026-042/045/046/047 | VS0-CF-F15-03,04,08,09,10,15,16,17,18,19,20,28 |
| VS0-STATE-002 | CloudEnrollment | FEATURE-0021 | DEC-0038; ADH-2026-021/037/042 | VS0-CF-F03 |
| VS0-STATE-003 | VersionedDefinition | FEATURE-0022 | DEC-0044; ADH-2026-027/042 | VS0-CF-F06 |
| VS0-STATE-004 | ExecutionTarget (Active/Retired lifecycle x Unqualified/Qualified/Rejected/Indeterminate qualification; Qualifying is a lifecycle-service-only in-flight reservation, never persisted or projected) | FEATURE-0016 | DEC-0042/0057; ADH-2026-025/040/042/045/058 | VS0-CF-F10,F16-01..122 |
| VS0-STATE-005 | QuotaReservation | FEATURE-0021 | DEC-0039; ADH-2026-022/042 | VS0-CF-F05,I01,I02 |
| VS0-STATE-006 | ServiceInstance | FEATURE-0007 | ADH-2026-042 Slice 0 adoption | VS0-CF-HP01,F15,F16,D01 |
| VS0-STATE-007 | ServiceBinding | FEATURE-0008+FEATURE-0024 | DEC-0051; ADH-2026-034/042 | VS0-CF-HP01,F17,F20,D01 |
| VS0-STATE-008 | Operation | FEATURE-0012 | ADH-2026-012/013/042 | VS0-CF-HP01,F14..F16 |
| VS0-STATE-009 | PluginExecution | FEATURE-0024 | DEC-0029; ADH-2026-042 | VS0-CF-HP01,F15,F16,Z01 |
| VS0-STATE-010 | ImmutableSlice0Records | FEATURE-0013 plus owning Slice 0 features | ADH-2026-017/042 | VS0-CF-HP01,T01 |
| VS0-STATE-011 | **Retired-CanonicalBootstrap** (was CanonicalMigrationRecord milestone sequence) | — | DEC-0059 supersedes DEC-0058; ADH-2026-045 | — |

---

## Failure/Error Mapping (VS0-F01..F20 → VS0-CF-F01..F20)

| Failure ID | Conformance ID | Code | HTTP | Violation | Owner | Controlling Decision |
|------------|----------------|------|------|-----------|-------|---------------------|
| VS0-F01 | VS0-CF-F01 | AUTH_REQUIRED | 401 | VS0_AUTH_REQUIRED | FEATURE-0018 | DEC-0050 |
| VS0-F02 | VS0-CF-F02 | RESOURCE_NOT_FOUND | 404 | VS0_AUTHORIZATION_SAFE_DENIAL | FEATURE-0018 | DEC-0050 |
| VS0-F03 | VS0-CF-F03 | VALIDATION_FAILED | 422 | VS0_ENROLLMENT_INACTIVE | FEATURE-0021 | DEC-0038 |
| VS0-F04 | VS0-CF-F04 | AUTHORIZATION_DENIED | 403 | VS0_ENTITLEMENT_MISSING | FEATURE-0021 | DEC-0039 |
| VS0-F05 | VS0-CF-F05 | CONFLICT | 409 | VS0_QUOTA_EXHAUSTED | FEATURE-0021 | DEC-0039 |
| VS0-F06 | VS0-CF-F06 | VALIDATION_FAILED | 422 | VS0_PLAN_NOT_PUBLISHED | FEATURE-0022 | DEC-0049 |
| VS0-F07 | VS0-CF-F07 | CONFLICT | 409 | VS0_GOVERNANCE_CONFLICT | FEATURE-0020 | DEC-0050 |
| VS0-F08 | VS0-CF-F08 | VALIDATION_FAILED | 422 | VS0_EVIDENCE_MISSING_OR_STALE | FEATURE-0019 | DEC-0041/0043/0055 |
| VS0-F09 | VS0-CF-F09 | VALIDATION_FAILED | 422 | VS0_PARTICIPATION_UNAVAILABLE | FEATURE-0023 | DEC-0043/0045 |
| VS0-F10 | VS0-CF-F10 | STALE_RESOURCE_VERSION | 412 | VS0_TARGET_EPOCH_STALE | FEATURE-0016 | DEC-0042/0057; ADH-2026-058 (exact local proof: VS0-CF-F16-29) |
| VS0-F11 | VS0-CF-F11 | VALIDATION_FAILED | 422 | VS0_PLACEMENT_NO_CANDIDATE | FEATURE-0023 | DEC-0043/0045 |
| VS0-F12 | VS0-CF-F12 | (none/202) | 202 | (none) | FEATURE-0024 | DEC-0029 |
| VS0-F13 | VS0-CF-F13 | CONFLICT | 409 | VS0_IDEMPOTENCY_PAYLOAD_MISMATCH | FEATURE-0024 | DEC-0029 |
| VS0-F14 | VS0-CF-F14 | STALE_RESOURCE_VERSION | 412 | VS0_GENERATION_STALE | FEATURE-0024 | DEC-0029 |
| VS0-F15 | VS0-CF-F15 | (none/200) | 200 | VS0_PLUGIN_FAILURE | FEATURE-0024 | DEC-0029 |
| VS0-F16 | VS0-CF-F16 | (none/200) | 200 | VS0_VERIFICATION_FAILURE | FEATURE-0024 | DEC-0029 |
| VS0-F17 | VS0-CF-F17 | CONFLICT | 409 | VS0_BINDING_NOT_READY | FEATURE-0008+FEATURE-0024 | DEC-0051 |
| VS0-F18 | VS0-CF-F18 | INTERNAL_ERROR | 500 | VS0_DISCLOSURE_PROHIBITED | FEATURE-0026 | ADH-2026-042 |
| VS0-F19 | VS0-CF-F19 | AUTHORIZATION_DENIED | 403 | VS0_AI_AUTHORITY_PROHIBITED | FEATURE-0025 | DEC-0043 |
| VS0-F20 | VS0-CF-F20 | (none/202) | 202 | VS0_DELETE_BINDING_ACTIVE | FEATURE-0024 | DEC-0029/0051 |

---

## FEATURE-0015 Migration Failure Mapping — Retired (Canonical Bootstrap)

Per DEC-0059 (supersedes DEC-0058) and ADH-2026-045, there is no migration failure mapping. `VS0-MIG-F01`, `VS0-MIG-F02`, and `VS0-MIG-F03` are permanently retired; their IDs must never be reused.

---

## Additional Conformance (HP01, X01..X03, L01, Z01, T01, I01..I02, D01, FEATURE-0015 local)

| Conformance ID | Class | Owner | Gate | Controlling Decision | Registry Evidence |
|----------------|-------|-------|------|---------------------|-------------------|
| VS0-CF-HP01 | Positive integration | FEATURE-0026 | slice0 | ADH-2026-042 | All schemas, writers, states |
| VS0-CF-X01 | Cross-organization isolation | FEATURE-0018 | security | DEC-0050 | VS0-WRITER-008, VS0-SCHEMA-023 |
| VS0-CF-X02 | Cross-project isolation | FEATURE-0018 | security | DEC-0050 | VS0-WRITER-008, VS0-SCHEMA-023 |
| VS0-CF-X03 | Cross-provider isolation | FEATURE-0015 | security | DEC-0037 | VS0-WRITER-003, VS0-SCHEMA-015 |
| VS0-CF-L01 | Leakage prevention | FEATURE-0026 | release-blocker | ADH-2026-042 | VS0-SCHEMA-051, redaction rules |
| VS0-CF-Z01 | Zero external effect | FEATURE-0024 | release-blocker | DEC-0029 | VS0-WRITER-006,016, VS0-SCHEMA-048 |
| VS0-CF-T01 | Correlation traceability | FEATURE-0026 | observability | ADH-2026-042 | VS0-SCHEMA-006,007,047,048 |
| VS0-CF-I01 | Idempotency same-digest | FEATURE-0024 | race | DEC-0029 | VS0-SCHEMA-035,050, VS0-STATE-005 |
| VS0-CF-I02 | Idempotency diff-digest | FEATURE-0024 | race | DEC-0029 | VS0-SCHEMA-035,050, VS0-STATE-005 |
| VS0-CF-D01 | Deletion ordering | FEATURE-0024 | lifecycle | DEC-0029/0051 | VS0-STATE-006,007, VS0-WRITER-019 |
| VS0-CF-F15-01 | Canonical resource registration | FEATURE-0015 | feature | DEC-0037/0041/0042 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-02 | Per-resource scope subset denial | FEATURE-0015 | feature | DEC-0037 | VS0-SCHEMA-001,008..014 |
| VS0-CF-F15-03 | Participation delegated acceptance | FEATURE-0015 | feature | DEC-0054; ADH-2026-037/042/045 | VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-04 | Participation pair uniqueness | FEATURE-0015 | feature | DEC-0054 | VS0-SCHEMA-010 |
| VS0-CF-F15-05 | Client status-write denial | FEATURE-0015 | security | FEATURE-0012 status contract | VS0-WRITER-005 |
| VS0-CF-F15-06 | System-owned metadata-write denial | FEATURE-0015 | security | FEATURE-0012 grammar | VS0-WRITER-001 |
| VS0-CF-F15-07 | PATCH-only update surface (no PUT/DELETE) | FEATURE-0015 | feature | DEC-0059; ADH-2026-045 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-08 | Independent suspension hold isolation | FEATURE-0015 | feature | DEC-0054; ADH-2026-045 | VS0-STATE-001, VS0-SCHEMA-010 |
| VS0-CF-F15-09 | Unauthorized CloudPlatform accept action leaves Pending participation unchanged | FEATURE-0015 | feature | DEC-0054; ADH-2026-045 | VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-10 | Suspend/resume denial in terminal/Pending/Terminating phases | FEATURE-0015 | feature | DEC-0054; ADH-2026-045 | VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-11 | Scope reference UID mismatch | FEATURE-0015 | feature | DEC-0037/0054; ADH-2026-042/043 | VS0-SCHEMA-008,010 |
| VS0-CF-F15-12 | CloudPlatform-root creation prerequisite | FEATURE-0015 | feature | ADH-2026-045/046 | VS0-SCHEMA-008..010 |
| VS0-CF-F15-13 | Valid PATCH accepted-fields-only mutation | FEATURE-0015 | feature | ADH-2026-045/046 | VS0-SCHEMA-008,009,011..014 |
| VS0-CF-F15-14 | Invalid PATCH denial (media type, immutable/status/scope/identity/relationship/owner-registration) | FEATURE-0015 | feature | ADH-2026-045/046 | VS0-SCHEMA-008,009,011..014, VS0-WRITER-002,003,005 |
| VS0-CF-F15-15 | Participation creation initial state | FEATURE-0015 | feature | DEC-0054; ADH-2026-045/046 | VS0-SCHEMA-010, VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-16 | Scheduler expiry transition | FEATURE-0015 | feature | ADH-2026-045/046 | VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-17 | Existing-participation If-Match precondition denial | FEATURE-0015 | feature | ADH-2026-046 decision 2 | VS0-WRITER-004 |
| VS0-CF-F15-18 | Target-bound, currently-authorized Idempotency-Key replay original-result return across all seven collection creates and eight participation item actions | FEATURE-0015 | idempotency | ADH-2026-045/046 decision 2/050/054 | VS0-SCHEMA-008..014, VS0-WRITER-002,003,004 |
| VS0-CF-F15-19 | Target-bound, currently-authorized Idempotency-Key reuse with changed digest across all seven collection creates and eight participation item actions | FEATURE-0015 | idempotency | ADH-2026-045/046 decision 2/050/054 | VS0-SCHEMA-008..014, VS0-WRITER-002,003,004 |
| VS0-CF-F15-20 | Explicit participation transition authorization and invalid-source-state denial | FEATURE-0015 | feature | DEC-0054; ADH-2026-045/046 | VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-21 | Read/list authorization (collection, scoped, empty, direct GET) | FEATURE-0015 | security | ADH-2026-045/046 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-22 | Exact forged-grant carriers and 401/403/400 precedence | FEATURE-0015 | security | ADH-2026-045/046/054 | VS0-WRITER-002,003,004 |
| VS0-CF-F15-23 | Topology reference ordering (safe denial vs authorized mismatch) | FEATURE-0015 | security | ADH-2026-043/045/046 | VS0-WRITER-003 |
| VS0-CF-F15-24 | Audit atomicity across F0015 mutations; exact durable-audit boundary for authenticated authorization/safe denials; required-AuditEvent-append failure returns INTERNAL_ERROR/500 (not DEPENDENCY_UNAVAILABLE) | FEATURE-0015 | observability | ADH-2026-043/046/049/051 | VS0-SCHEMA-007 |
| VS0-CF-F15-25 | F0015 route/action surface closure (no F0016/migration leakage) | FEATURE-0015 | feature | DEC-0059; ADH-2026-045/046 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-26 | Per-kind collection-create request boundary | FEATURE-0015 | feature | ADH-2026-047 decision 1/5 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-27 | Scope derivation single-source proof | FEATURE-0015 | feature | ADH-2026-047 decision 2/5 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-28 | Participation collection-create body versus empty item-action body | FEATURE-0015 | feature | ADH-2026-047 decision 4/5 | VS0-SCHEMA-010, VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-29 | ISO-3166-1 alpha-2 assigned-code rejection (syntactically valid but unassigned value denied) | FEATURE-0015 | feature | ADH-2026-048 decision 1 | VS0-SCHEMA-008,009,011 |
| VS0-CF-F15-30 | Idempotency-Key/If-Match/item-action-body malformed-input outcomes | FEATURE-0015 | feature | ADH-2026-048 decision 2 | VS0-SCHEMA-010, VS0-WRITER-004 |
| VS0-CF-F15-31 | Missing/invalid authentication on any F0015 owned method/path pattern (AUTH_REQUIRED local proof; does not replace inherited VS0-CF-F01) | FEATURE-0015 | security | ADH-2026-051 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-32 | CloudPlatform/CloudProvider metadata.name duplicate rejection (ALREADY_EXISTS local proof) | FEATURE-0015 | feature | ADH-2026-052 | VS0-SCHEMA-008,009 |
| VS0-CF-F15-33 | Registry-declared schema-constraint violation on collection create, excluding duplicate name (F15-32), unassigned ISO code (F15-29), and malformed/prohibited header/body (F15-30) | FEATURE-0015 | feature | ADH-2026-052 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-34 | HEAD matched to a GET pattern returns 405 without increasing registrations | FEATURE-0015 | feature | ADH-2026-055 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-35 | Forged header denial precedes unsupported PATCH media validation | FEATURE-0015 | security | ADH-2026-055 | VS0-SCHEMA-004,007 |
| VS0-CF-F15-36 | Bounded completed-idempotency retention, deterministic eviction, and waiter cancellation | FEATURE-0015 | idempotency | ADH-2026-055 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-37 | RFC 7396 staged PATCH/null behavior and deterministic LIST ordering | FEATURE-0015 | feature | ADH-2026-055 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-38 | PATCH current-version conditional update and lock recheck | FEATURE-0015 | feature | ADH-2026-056 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-39 | Exact scoped read actions and safe GET/LIST outcomes | FEATURE-0015 | security | ADH-2026-056 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-40 | Completed replay response allowlist and Aborted re-reservation | FEATURE-0015 | idempotency | ADH-2026-056 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-41 | Mutable-spec writer-boundary denial | FEATURE-0015 | security | DEC-0037; ADH-2026-057 | VS0-WRITER-002,003 |

## FEATURE-0016 Conformance (VS0-CF-F16-01..128)

FEATURE-0016 replaces the placeholder `VS0-SCHEMA-015..017`/`VS0-STATE-004` definitions with the approved ADH-2026-058 ExecutionTarget executable contract: closed create fields and derived CloudProvider scope, five exact routes, a synthetic-only observer, a lifecycle-service-only Qualifying reservation that is never persisted or projected, a response-only `effectiveAvailability` projection, and the complete fail-closed proof annex below. `VS0-CF-F10` and `VS0-CF-X03` remain the cross-feature/shared references; every row below is the exact F0016-local proof case per ADH-2026-058 §2/§3.

ADH-2026-063 adds the collection-create phase-one/strict-classification precedence correction and its four new sequential local conformance rows `VS0-CF-F16-123..126`. Its safe precedence invariant is: the bounded single body read returns the existing `REQUEST_TOO_LARGE`/400 first for an oversized body; otherwise phase one extracts only `spec.cloudProviderParticipationRef.uid` and `spec.infrastructureStackRef.uid` and never classifies, canonicalizes, digests, or reserves a body; derived-scope `executiontarget.write` authorization and safe backing access are evaluated before the retained same bytes are strictly classified exactly once. Existing rows `VS0-CF-F16-08`, `VS0-CF-F16-09`, `VS0-CF-F16-51`, `VS0-CF-F16-52`, and `VS0-CF-F16-53` are annotated as precedence anchors for these new rows; their existing case names and observable meanings are unchanged.

ADH-2026-065 makes the create classification and phase-one reference conformance exact and extends the FEATURE-0016 local range to `VS0-CF-F16-128`. `VS0-CF-F16-125` is now the exact non-family malformed-JSON classification case returning the existing 400 `MALFORMED_REQUEST` (anchored on `VS0-CF-F16-51`), and the added `VS0-CF-F16-127` is the exact non-family duplicate-top-level-member classification case returning the existing 400 `DUPLICATE_FIELD` (anchored on `VS0-CF-F16-52`); the two are distinct exact cases, not a mixed equivalence family. The added `VS0-CF-F16-128` is the fail-closed phase-one extraction case: when the existing bounded body read succeeds but phase one cannot extract exactly one syntactically usable UID at both `spec.cloudProviderParticipationRef.uid` and `spec.infrastructureStackRef.uid`, the create returns the existing audited 403 `AUTHORIZATION_DENIED` (anchored on `VS0-CF-F16-08`) before derived-scope authorization and safe backing access, disclosing no body-classification detail and no backing-resource existence, and performing no strict classification, canonicalization, digest, idempotency reservation, observer invocation, mutation, publication, or completion; it adds no violation or top-level Problem code. `VS0-CF-F16-123`, `VS0-CF-F16-124`, and `VS0-CF-F16-126` remain unchanged. `VS0-CF-F16-125`, `VS0-CF-F16-127`, and `VS0-CF-F16-128` map to `REQ-F16-05` and its affected acceptance proof. The controlling FEATURE-0016 semantic authorities are ADH-2026-058, ADH-2026-060, ADH-2026-061, ADH-2026-063, ADH-2026-064, ADH-2026-065, and ADH-2026-066; ADH-2026-058 is not the sole current semantic authority.

ADH-2026-066 adds no route, schema, state, violation, or local conformance row. It records the F0016-owned `BackingAccessProvider` and its private, read-only FEATURE-0015 store-backed lease: paired immutable backing inputs and caller-relative safe access are obtained under one store lock; a qualification's final recheck compares only the registered InfrastructureStack-generation and viability-fingerprint fences; and the lease remains held through F0016 audit-before-publication. FEATURE-0015 retains ownership of `CloudProviderParticipation` and `InfrastructureStack`; the bridge introduces no FEATURE-0015 public behavior.

| Conformance ID | Class | Owner | Gate | Controlling Decision | Registry Evidence |
|----------------|-------|-------|------|---------------------|-------------------|
| VS0-CF-F16-01 | Valid create | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-02 | Authorized create with one unknown field | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-03 | Authorized create with client status | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-04 | Authorized create with client metadata.uid | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-05 | Create/action with a missing Idempotency-Key | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-06 | POST with non-application/json media | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-07 | Create with missing authentication | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-08 | Authenticated create without executiontarget.write | FEATURE-0016 | feature | ADH-2026-058; ADH-2026-063 (precedence anchor) | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-09 | Inaccessible create backing reference | FEATURE-0016 | feature | ADH-2026-058; ADH-2026-063 (precedence anchor) | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-10 | Authorized ineffective participation | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-11 | Authorized non-Active InfrastructureStack | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-12 | Authorized cross-CloudProvider participation/stack graph | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-13 | Create with a duplicate derived name | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-14 | Authorized LIST with an arbitrary Idempotency-Key | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-15 | LIST without applicable executiontarget.read | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-16 | Authorized item GET with an arbitrary Idempotency-Key | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-17 | Inaccessible direct item GET | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-18 | Qualify without executiontarget.qualify | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-19 | Inaccessible direct qualify action | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-20 | Qualify, four Supported facts | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-21 | Qualify, any Unsupported fact | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-22 | Qualify with a missing target-bound fixture | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-23 | Qualify from Qualified receives one malformed fact while all captured fences remain current | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-24 | Action with a missing If-Match | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-25 | Action valid If-Match with non-zero body | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-26 | Qualify Retired target | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-27 | Qualify Maintenance target | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-28 | Different-key qualify while Qualifying | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-29 | Referenced InfrastructureStack generation changes at qualification commit | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-30 | Viability fingerprint changes at qualification commit | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-31 | Same-key completed replay after current authorization and safe access | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-32 | Same-key replay requester loses its grant while target remains safe-accessible | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-33 | Same namespace with changed digest | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-34 | Completed record expires after 24 hours | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-35 | Owner panic after reservation | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-36 | Valid retire from any Active combination | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-37 | Authenticated retire without executiontarget.write | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-38 | New-key retire after retirement | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-39 | Valid create required AuditEvent append fails | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-40 | Create authorization-denial AuditEvent append fails | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-41 | Fact expiry with successful required AuditEvent append | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-42 | Current fenced maintenance-entry trigger | FEATURE-0016 | feature | ADH-2026-058; ADH-2026-061 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-43 | Fenced maintenance clear | FEATURE-0016 | feature | ADH-2026-058; ADH-2026-061 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-44 | HEAD on the item GET path | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-45 | Successful GET with any Accept value | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-46 | Process restart | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-47 | Same-key replay target is no longer safe-accessible | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-48 | Problem response with any Accept value | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-49 | Same principal/key/digest action on another target UID | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-50 | Same principal/key/digest create in another derived CloudProvider scope | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-51 | Malformed JSON request body | FEATURE-0016 | feature | ADH-2026-058; ADH-2026-063 (precedence anchor) | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-52 | Duplicate top-level JSON member | FEATURE-0016 | feature | ADH-2026-058; ADH-2026-063 (precedence anchor) | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-53 | Oversized request body | FEATURE-0016 | feature | ADH-2026-058; ADH-2026-063 (precedence anchor) | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-54 | Missing required create field | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-55 | Invalid typed reference shape | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-56 | targetClass other than synthetic-iaas | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-57 | Authorized create with adapterAuthorityRef | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-58 | Create/action with an empty Idempotency-Key | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-59 | Create/action with a malformed Idempotency-Key | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-60 | Create/action with an overlong Idempotency-Key | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-61 | LIST with missing authentication | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-62 | Item GET with missing authentication | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-63 | Qualify with missing authentication | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-64 | Retire with missing authentication | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-65 | Create with a live duplicate backing tuple | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-66 | Inaccessible direct retire action | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-67 | Qualify receives a duplicate named fact while all captured fences remain current | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-68 | Qualify receives a missing named fact while all captured fences remain current | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-69 | Qualify receives malformed fact provenance while all captured fences remain current | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-70 | Action with a malformed If-Match | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-71 | Action with a weak If-Match | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-72 | Action with wildcard If-Match | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-73 | Action with multiple If-Match headers | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-74 | Action with a syntactically valid stale If-Match | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-75 | Captured maintenance epoch differs at qualification commit and no active current-Maintenance marker exists | FEATURE-0016 | feature | ADH-2026-058; ADH-2026-060 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-76 | Completed record is selected for deterministic capacity eviction | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-77 | Server shutdown after reservation | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-78 | Owner cancellation after reservation | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-79 | Fact expiry with required AuditEvent append failure | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-80 | Stale fenced maintenance-entry trigger | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-81 | Same-key waiter wakes after its owner aborts on a stale InfrastructureStack-generation qualification | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-82 | Required qualification AuditEvent append failure after reservation while all captured fences remain current | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-83 | Stale fenced maintenance-clear trigger | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-84 | Authorized create with a SecretRef extension field | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-85 | Qualify Maintenance target with syntactically valid stale If-Match | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-86 | Qualify with malformed If-Match and non-zero body | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-87 | Same-key completed qualify after target ETag changes | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-88 | Retire wins while qualify is in flight | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-89 | Maintenance entry wins while qualify is in flight | FEATURE-0016 | feature | ADH-2026-058; ADH-2026-060; ADH-2026-061 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-90 | Qualify with fixture-declared logical timeout | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-91 | Authorized create with client metadata.resourceVersion | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-92 | Authorized create with client metadata.generation | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-93 | Create with invalid authentication | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-94 | LIST with invalid authentication | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-95 | Item GET with invalid authentication | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-96 | Qualify with invalid authentication | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-97 | Retire with invalid authentication | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-98 | Create with a name reserved by a Retired target | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-99 | Valid retirement required AuditEvent append fails | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-100 | Current maintenance-entry required AuditEvent append fails | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-101 | Current maintenance-clear required AuditEvent append fails | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-102 | Inaccessible direct item GET required AuditEvent append fails | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-103 | Same-key replay authorization denial AuditEvent append fails | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-104 | Same-key replay safe-denial AuditEvent append fails | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-105 | Collection create with an If-Match header | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-106 | Authenticated URI-query equivalence family: POST collection, GET collection, GET item, POST qualify, and POST retire, each with `?x=1` | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-107 | Authenticated item GET with a malformed canonical UID segment | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-108 | Item GET with trailing slash | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-109 | Authenticated qualify action with a malformed canonical UID segment | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-110 | Authenticated retire action with a malformed canonical UID segment | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-111 | Direct item GET of a safe-accessible target without executiontarget.read | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-112 | Denial AuditEvent append-failure equivalence family: LIST without read; direct GET without read; qualify without qualify grant; retire without write grant; inaccessible direct qualify; inaccessible direct retire | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-113 | HEAD on collection path | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-114 | HEAD equivalence family: qualify action and retire action paths | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-115 | Trailing-slash equivalence family: POST collection, GET collection, POST qualify action, and POST retire action | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-116 | PUT on collection path | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-117 | PUT on item path | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-118 | PUT equivalence family: qualify action and retire action paths | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-119 | Fresh Qualified target backing becomes non-viable | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-120 | Fresh Qualified target backing returns viable before fact expiry | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-121 | Qualify starts from Active/Qualified with current fences | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-122 | Current-fence observer fault after qualification starts from Active/Qualified | FEATURE-0016 | feature | ADH-2026-058 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-123 | Unauthorized create (no executiontarget.write) with a malformed/duplicate body over safe-accessible backing | FEATURE-0016 | feature | ADH-2026-063 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-124 | Inaccessible-backing create with a malformed/duplicate body | FEATURE-0016 | feature | ADH-2026-063 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-125 | Authorized create over safe-accessible backing with a malformed-JSON retained body strictly classified to MALFORMED_REQUEST (exact non-family case, distinct from F16-127) | FEATURE-0016 | feature | ADH-2026-065 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-126 | Oversized collection-create body rejected before phase-one extraction, authorization, or safe access | FEATURE-0016 | feature | ADH-2026-063 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-127 | Authorized create over safe-accessible backing with a duplicate-top-level-member retained body strictly classified to DUPLICATE_FIELD (exact non-family case, distinct from F16-125) | FEATURE-0016 | feature | ADH-2026-065 | VS0-SCHEMA-015..017, VS0-STATE-004 |
| VS0-CF-F16-128 | Missing/unextractable required phase-one reference (not exactly one usable UID at both spec.cloudProviderParticipationRef.uid and spec.infrastructureStackRef.uid) fails closed to existing audited AUTHORIZATION_DENIED before derived-scope authorization and safe backing access, with no body/backing disclosure or classification/canonicalization/digest/reservation/observer/mutation/publication/completion | FEATURE-0016 | feature | ADH-2026-065 | VS0-SCHEMA-015..017, VS0-STATE-004 |



Per DEC-0059 (supersedes DEC-0058) and ADH-2026-045, `VS0-CF-MIG01`, `VS0-CF-MIG02`, `VS0-CF-MIGF01`, `VS0-CF-MIGF02`, and `VS0-CF-MIGF03` are permanently retired; their IDs must never be reused. FEATURE-0015 has no migration proof because there is no live alpha state to convert.

---

## Repository Artifacts

| Artifact | Path |
|----------|------|
| VS-000 charter | `docs/architecture/vertical-slices/VS-000-core-skeleton.md` |
| VS-000 contract specification | `docs/architecture/vertical-slices/VS-000-contract-specification.md` |
| VS-000 contract registry YAML | `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` |
| VS-000 contract traceability | `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md` |
| Kiro architecture steering | `.kiro/steering/slice0-contract.md` |
| Validation script | `scripts/vs000-contract-check.py` |
| Validation target | `make vs000-contract-check` |
| Feature gate target | `make ff-feature-gate FEATURE=FEATURE-0026` |

---

## Anti-Orphan Rules

1. Every future requirement must map to: ADH/DEC controlling decision + owner FEATURE + registry schema/writer/state/error ID + conformance ID.
2. No registry ID may be renumbered or reused after initial assignment.
3. Any change to schema, writer, state, or error contract must update registry YAML, this traceability matrix, and conformance tests in the same change.
4. A missing mapping for any requirement blocks generation with `ARCHITECTURE_DECISION_REQUIRED`.
5. Traceability completeness is validated by `make vs000-contract-check`.

---

## Runtime Conformance Status

No runtime conformance test implementation exists yet. All conformance IDs (VS0-CF-HP01, VS0-CF-F01..F20, VS0-CF-X01..X03, VS0-CF-L01, VS0-CF-Z01, VS0-CF-T01, VS0-CF-I01..I02, VS0-CF-D01, VS0-CF-F15-01..41, VS0-CF-F16-01..122) are exact test contracts owned by their respective feature tasks. FEATURE-0026 provides the integration proof that exercises the cross-feature Slice 0 contracts end-to-end in the synthetic profile; FEATURE-0015 and FEATURE-0016 each own their local conformance cases only (no migration conformance exists under DEC-0059). FEATURE-0015 owns exactly 22 logical endpoint paths and registers exactly 35 explicit Go 1.22 `http.ServeMux` method/path patterns (ADH-2026-051). FEATURE-0016 owns exactly five explicit Go 1.22 `http.ServeMux` registrations plus a pre-ServeMux transport-only method/path guard (ADH-2026-058 clause 3).

---

*Slice 0 only. No future-slice content. No FEATURE-0015 stage generation. No new architecture decisions.*
