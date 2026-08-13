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
| 3 | Accepted decisions and controlling handoff | `docs/decisions/DECISION_INDEX.md`; ADH-2026-042 |
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
| VS0-SCHEMA-008 | CloudPlatform | FEATURE-0015 | DEC-0037, ADH-2026-042/045/046/047/048/050/051 | canonical data model | VS0-CF-F15-01,02,07,12,13,14,18,19,25,26,27,29,30,31 |
| VS0-SCHEMA-009 | CloudProvider | FEATURE-0015 | DEC-0037, ADH-2026-042/045/046/047/048/050/051 | canonical data model | VS0-CF-F15-01,02,07,12,13,14,18,19,25,26,27,29,30,31 |
| VS0-SCHEMA-010 | CloudProviderParticipation | FEATURE-0015 | DEC-0054; ADH-2026-037/042/045/046/047/048/050/051 | canonical data model | VS0-CF-F15-01..04,08..10,15..21,25,26,27,28,30,31 |
| VS0-SCHEMA-011 | HostingLocation | FEATURE-0015 | DEC-0041; ADH-2026-024/042/045/046/047/048/050/051 | canonical data model | VS0-CF-F15-01,02,07,12,13,14,18,19,25,26,27,29,30,31 |
| VS0-SCHEMA-012 | Datacenter | FEATURE-0015 | DEC-0041; ADH-2026-024/042/045/046/047/050/051 | canonical data model | VS0-CF-F15-01,02,07,12,13,14,18,19,25,26,27,30,31 |
| VS0-SCHEMA-013 | FaultDomain | FEATURE-0015 | DEC-0041; ADH-2026-024/042/045/046/047/050/051 | canonical data model | VS0-CF-F15-01,02,07,12,13,14,18,19,25,26,27,30,31 |
| VS0-SCHEMA-014 | InfrastructureStack | FEATURE-0015 | DEC-0041/0042; ADH-2026-025/042/045/046/047/050/051 | canonical data model | VS0-CF-F15-01,02,07,12,13,14,18,19,25,26,27,30,31 |
| VS0-SCHEMA-015 | ExecutionTarget | FEATURE-0016 | DEC-0042/0057; ADH-2026-025/040/042/045 | canonical data model | VS0-CF-F10,X03 |
| VS0-SCHEMA-016 | NormalizedTargetFactSet | FEATURE-0016 | DEC-0036/0042/0057; ADH-2026-040/042 | canonical data model | VS0-CF-F10 |
| VS0-SCHEMA-017 | TargetQualificationResult | FEATURE-0016 | DEC-0036/0042/0057; ADH-2026-040/042 | canonical data model | VS0-CF-F10 |
| VS0-SCHEMA-018 | PolicyEvaluationRequest | FEATURE-0017 | DEC-0028/0043; ADH-2026-026/042 | canonical data model | VS0-CF-F07..F09 |
| VS0-SCHEMA-019 | PolicyEvaluationResult | FEATURE-0017 | DEC-0028/0043; ADH-2026-026/042 | canonical data model | VS0-CF-F07..F09 |
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
| VS0-WRITER-006 | fake-adapter-and-qualification-controller | FEATURE-0016 | DEC-0036/0042/0057 | VS0-CF-F10,Z01 |
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
| VS0-STATE-004 | ExecutionTarget | FEATURE-0016 | DEC-0042/0057; ADH-2026-025/040/042/045 | VS0-CF-F10 |
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
| VS0-F10 | VS0-CF-F10 | STALE_RESOURCE_VERSION | 412 | VS0_TARGET_EPOCH_STALE | FEATURE-0016 | DEC-0042/0057 |
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
| VS0-CF-F15-09 | Participation activation denial without both delegations | FEATURE-0015 | feature | DEC-0054; ADH-2026-037/042 | VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-10 | Suspend/resume denial in terminal/Pending/Terminating phases | FEATURE-0015 | feature | DEC-0054; ADH-2026-045 | VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-11 | Scope reference UID mismatch | FEATURE-0015 | feature | DEC-0037/0054; ADH-2026-042/043 | VS0-SCHEMA-008,010 |
| VS0-CF-F15-12 | CloudPlatform-root creation prerequisite | FEATURE-0015 | feature | ADH-2026-045/046 | VS0-SCHEMA-008..010 |
| VS0-CF-F15-13 | Valid PATCH accepted-fields-only mutation | FEATURE-0015 | feature | ADH-2026-045/046 | VS0-SCHEMA-008,009,011..014 |
| VS0-CF-F15-14 | Invalid PATCH denial (media type, immutable/status/scope/identity/relationship/owner-registration) | FEATURE-0015 | feature | ADH-2026-045/046 | VS0-SCHEMA-008,009,011..014, VS0-WRITER-002,003,005 |
| VS0-CF-F15-15 | Participation creation initial state | FEATURE-0015 | feature | DEC-0054; ADH-2026-045/046 | VS0-SCHEMA-010, VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-16 | Scheduler expiry transition | FEATURE-0015 | feature | ADH-2026-045/046 | VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-17 | Existing-participation If-Match precondition denial | FEATURE-0015 | feature | ADH-2026-046 decision 2 | VS0-WRITER-004 |
| VS0-CF-F15-18 | Idempotency-Key replay original-result return across all seven collection creates and eight participation create/actions | FEATURE-0015 | idempotency | ADH-2026-045/046 decision 2/050 | VS0-SCHEMA-008..014, VS0-WRITER-002,003,004 |
| VS0-CF-F15-19 | Idempotency-Key reuse with changed digest across all seven collection creates and eight participation create/actions | FEATURE-0015 | idempotency | ADH-2026-045/046 decision 2/050 | VS0-SCHEMA-008..014, VS0-WRITER-002,003,004 |
| VS0-CF-F15-20 | Explicit participation transition authorization and invalid-source-state denial | FEATURE-0015 | feature | DEC-0054; ADH-2026-045/046 | VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-21 | Read/list authorization (collection, scoped, empty, direct GET) | FEATURE-0015 | security | ADH-2026-045/046 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-22 | Bootstrap-grant boundary (header/body claims ignored) | FEATURE-0015 | security | ADH-2026-045/046 | VS0-WRITER-002,003,004 |
| VS0-CF-F15-23 | Topology reference ordering (safe denial vs authorized mismatch) | FEATURE-0015 | security | ADH-2026-043/045/046 | VS0-WRITER-003 |
| VS0-CF-F15-24 | Audit atomicity across F0015 mutations; exact durable-audit boundary for authenticated authorization/safe denials; required-AuditEvent-append failure returns INTERNAL_ERROR/500 (not DEPENDENCY_UNAVAILABLE) | FEATURE-0015 | observability | ADH-2026-043/046/049/051 | VS0-SCHEMA-007 |
| VS0-CF-F15-25 | F0015 route/action surface closure (no F0016/migration leakage) | FEATURE-0015 | feature | DEC-0059; ADH-2026-045/046 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-26 | Per-kind collection-create request boundary | FEATURE-0015 | feature | ADH-2026-047 decision 1/5 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-27 | Scope derivation single-source proof | FEATURE-0015 | feature | ADH-2026-047 decision 2/5 | VS0-SCHEMA-008..014 |
| VS0-CF-F15-28 | Participation collection-create body versus empty item-action body | FEATURE-0015 | feature | ADH-2026-047 decision 4/5 | VS0-SCHEMA-010, VS0-STATE-001, VS0-WRITER-004 |
| VS0-CF-F15-29 | ISO-3166-1 alpha-2 assigned-code rejection (syntactically valid but unassigned value denied) | FEATURE-0015 | feature | ADH-2026-048 decision 1 | VS0-SCHEMA-008,009,011 |
| VS0-CF-F15-30 | Idempotency-Key/If-Match/item-action-body malformed-input outcomes | FEATURE-0015 | feature | ADH-2026-048 decision 2 | VS0-SCHEMA-010, VS0-WRITER-004 |
| VS0-CF-F15-31 | Missing/invalid authentication on any F0015 owned method/path pattern (AUTH_REQUIRED local proof; does not replace inherited VS0-CF-F01) | FEATURE-0015 | security | ADH-2026-051 | VS0-SCHEMA-008..014 |

## FEATURE-0015 Migration Conformance — Retired (Canonical Bootstrap)

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

No runtime conformance test implementation exists yet. All conformance IDs (VS0-CF-HP01, VS0-CF-F01..F20, VS0-CF-X01..X03, VS0-CF-L01, VS0-CF-Z01, VS0-CF-T01, VS0-CF-I01..I02, VS0-CF-D01, VS0-CF-F15-01..31) are exact test contracts owned by their respective feature tasks. FEATURE-0026 provides the integration proof that exercises the cross-feature Slice 0 contracts end-to-end in the synthetic profile; FEATURE-0015 owns its local conformance cases only (no migration conformance exists under DEC-0059). FEATURE-0015 owns exactly 22 logical endpoint paths and registers exactly 35 explicit Go 1.22 `http.ServeMux` method/path patterns (ADH-2026-051).

---

*Slice 0 only. No future-slice content. No FEATURE-0015 stage generation. No new architecture decisions.*
