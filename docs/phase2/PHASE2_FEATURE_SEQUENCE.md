---
doc_type: feature_sequence
title: Phase 2R Feature Sequence
status: approved
phase: 2R
ai_load_priority: always
ai_summary: Dependency-safe Phase 2R feature sequence after canonical model adoption (ADH-2026-042).
---

# Phase 2R Feature Sequence

Controlling authority: `docs/phase2/PHASE2R_REBASELINE.md`

Architecture baseline: `ARCH-2026.08-PHASE2R-CANONICAL`

## Completed Features (Implementation History)

| Order | Feature | Status | Notes |
|---:|---|---|---|
| 1 | FEATURE-0011 Reuse Assessment Standard | Merged | Phase 2 baseline |
| 2 | FEATURE-0012 API, Resource Naming, Status, and Validation Standard | Merged (PR #14) | ADH-2026-012 |
| 3 | FEATURE-0013 Decision Record and AuditEvent Standard | Merged (PR #15) | ADH-2026-017 |
| 4 | FEATURE-0014 Provider-Neutral Resource Model | Merged (PR #16) | ADH-2026-018/019; alpha model |

## Phase 2R Rebaselined Features

| Order | Feature | Depends On | Output |
|---:|---|---|---|
| 5 | FEATURE-0015 Canonical Cloud Model Foundation | FEATURE-0011–0014; ADH-020/024/025/037/045 | Seven-scope canonical bootstrap; CloudPlatform; CloudProvider; CloudProviderParticipation; HostingLocation; Datacenter; FaultDomain; InfrastructureStack (ends here; ExecutionTarget owned by FEATURE-0016) |
| 6 | FEATURE-0016 Adapter Boundary and ExecutionTarget Qualification | FEATURE-0015; DEC-0036 | Core adapter interfaces; normalized target facts; target qualification; fake adapter |
| 7 | FEATURE-0017 Policy Evaluation Abstraction | FEATURE-0013, FEATURE-0016 | PolicyEvaluationRequest/Result; PolicyEngineAdapter; DecisionRecord linkage; deterministic fake |
| 8 | FEATURE-0018 Governance, IAM, Approval and Exception Foundation | FEATURE-0012, FEATURE-0017; FEATURE-0013 and FEATURE-0016 registrations | AccessGroup; Membership; RoleDefinition; RoleAssignment; PrivilegedAccessRequest; AccessReview; ApprovalPolicy; ApprovalRequest; ExceptionGrant; FEATURE-0018-limited GovernanceProfile |
| 9 | FEATURE-0019 Sovereignty Facts, Evidence and Policy Foundation | FEATURE-0013, FEATURE-0016 | SovereigntyProfile; RegulatoryPolicyBundle; SovereigntyFactSet; EvidenceRecord; governed dimension registry |
| 10 | FEATURE-0020 Assignment and Effective Governance Resolution | FEATURE-0018, FEATURE-0019 | ProfileAssignment; immutable EffectiveGovernanceContext; deterministic inheritance/conflict/exception resolution |
| 11 | FEATURE-0021 CloudEnrollment, Personal Onboarding, Entitlement and Quota | FEATURE-0015, FEATURE-0020 | CloudEnrollment; Personal Organization; EntitlementPackage; ServiceEntitlement; QuotaPolicy; provider-selection intent |
| 12 | FEATURE-0022 Service Product, Runtime and Requirement Foundation | FEATURE-0015, FEATURE-0016, FEATURE-0021; ADH-023/024/027/032 | ServicePortfolio; ServiceTypeDefinition; ServiceOffering; ServicePlan migration; ServicePlacementProfile; ServiceRuntimeProfile; ServiceRequirementSet; ServiceRegion |
| 13 | FEATURE-0023 Sovereignty and Placement Decision v0 | FEATURE-0019–0022 | Registered DecisionProfiles; authoritative DecisionRecords; candidate-set evaluation; ServicePlacement safe projection |
| 14 | FEATURE-0024 Plugin Taxonomy and Synthetic Execution Boundary | FEATURE-0016, FEATURE-0022, FEATURE-0023 | Plugin/adapter roles; manifests; compatibility; protected handles; PluginExecution contract; fake harness |
| 15 | FEATURE-0025 AI-Readable Decision and Operation Context | FEATURE-0013, FEATURE-0023, FEATURE-0024 | Bounded audience-safe explanation projection |
| 16 | FEATURE-0026 Slice 0 Integration and Conformance Demo | All Phase 2R features | VS-000 cross-feature orchestration; fake execution; conformance matrix |

## Dependency Graph

```text
FEATURE-0011–0014 completed baseline
              │
              ▼
FEATURE-0015 canonical migration and cloud model
              │
              ▼
FEATURE-0016 adapters and target qualification
       ┌──────┴───────────┐
       ▼                  ▼
FEATURE-0017        FEATURE-0019
       │                  │
       ▼                  │
FEATURE-0018              │
       └──────┬───────────┘
              ▼
FEATURE-0020 effective governance
              │
              ▼
FEATURE-0021 enrollment/entitlement/quota
              │
              ▼
FEATURE-0022 product/runtime/requirements
              │
              ▼
FEATURE-0023 sovereignty and placement decisions
              │
              ▼
FEATURE-0024 plugin taxonomy + synthetic execution boundary
              │
              ▼
FEATURE-0025 safe AI-readable context
              │
              ▼
FEATURE-0026 Slice 0 integration demo
```

FEATURE-0019 may proceed in parallel with FEATURE-0017/0018 after FEATURE-0016.

## Rule

Do not start FEATURE-0023 before FEATURE-0019, FEATURE-0020, FEATURE-0021, and FEATURE-0022 exist.

Do not implement features using ResourcePool, ProviderCapability, or generic Provider as active concepts. Use ExecutionTarget, CloudPlatform, and CloudProvider.

Do not implement CanonicalMigrationPlan, CanonicalMigrationRecord, a migration controller, or a cutover state machine in FEATURE-0015 (DEC-0059 supersedes DEC-0058). FEATURE-0001–0014 are retained repository assets, not live state requiring conversion.
