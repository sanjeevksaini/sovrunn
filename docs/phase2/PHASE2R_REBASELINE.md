# Sovrunn Phase 2R Rebaseline

**Status:** Approved target rebaseline; repository merge pending
**Date:** 4 August 2026
**Baseline input:** `ARCH-2026.07-PHASE2-START` at commit `7bced1068f2837ead9709e78d3fc33f8e6e032c5`
**Proposed successor baseline:** `ARCH-2026.08-PHASE2R-CANONICAL`
**Controlling approval:** `ARCH-APPROVAL-2026-004` and consolidated `ADH-2026-042`

## 1. Rebaseline rule

FEATURE-0001 through FEATURE-0014 remain completed implementation history. FEATURE-0015 through FEATURE-0026 are unimplemented placeholders and are replaced in place so durable feature IDs remain sequential.

The old FEATURE-0015 requirements must not be generated. Mandatory `ResourcePool` and `ProviderCapability` are removed from the active Phase 2 design. Existing `DEC-0032` and `DEC-0033` become superseded.

Phase 2R remains a contract and side-effect-free simulation phase. It may use deterministic in-memory fakes to prove integration, but it must not call real provider, Kubernetes, OpenShift or PostgreSQL APIs.

## 2. Accepted decision mapping

| DEC | Controlling handoff | Accepted decision |
|---|---|---|
| DEC-0037 | ADH-2026-020 | Retire ambiguous Provider by classifying records into CloudPlatform ownership and/or CloudProvider operation; introduce seven canonical scopes |
| DEC-0038 | ADH-2026-021 | CloudEnrollment joins customer Organization to CloudPlatform; provider participation remains separate |
| DEC-0039 | ADH-2026-022 | Entitlement answers what may be consumed; quota independently answers how much |
| DEC-0040 | ADH-2026-023 | Industry and domain clouds are versioned ServicePortfolios, not new core provider kinds |
| DEC-0041 | ADH-2026-024 | Physical geography, governance scope and sovereignty are independent dimensions |
| DEC-0042 | ADH-2026-025 | Qualified ExecutionTarget is the actionable realization boundary; no mandatory ResourcePool or provider-wide capability exists in core |
| DEC-0043 | ADH-2026-026 | SovereigntyAssessment and PlacementDecision use registered FEATURE-0013 DecisionRecord profiles |
| DEC-0044 | ADH-2026-027 | Published product, runtime, governance, sovereignty and policy definitions are immutable by version |
| DEC-0045 | ADH-2026-028 | One service may be realized across an atomic set of execution targets with explicit cross-target constraints |
| DEC-0046 | ADH-2026-029 | ServicePlacement is a safe immutable customer projection, not canonical topology access |
| DEC-0047 | ADH-2026-030 | Personal onboarding creates a Personal Organization, default Tenant and default Project |
| DEC-0048 | ADH-2026-031 | New countries, sectors, services and backends enter through governed data, contracts, plugins and adapters—not core conditionals |
| DEC-0049 | ADH-2026-032 | ServiceOffering is the CloudPlatform product and references a reusable ServiceTypeDefinition; global ServiceClass is migration input |
| DEC-0050 | ADH-2026-033 | Governance controls compose into one EffectiveGovernanceContext; sovereignty remains separately assessed |
| DEC-0051 | ADH-2026-034 | ServiceBinding is the per-consumer, separately revocable, SecretRef-only service-access boundary |
| DEC-0052 | ADH-2026-035 | ServiceRelationshipDefinition and ServiceRelationship model implementation-neutral cross-service semantics and lifecycle |
| DEC-0053 | ADH-2026-036 | Sovrunn platform lifecycle uses one external accepted-intent authority, Operation-first activation and an independently recoverable agent |
| DEC-0054 | ADH-2026-037 | CloudPlatform ownership, CloudProvider participation and one-provider-per-installation isolation are distinct boundaries |
| DEC-0055 | ADH-2026-038 | Platform sovereignty evaluates one SovrunnInstallation and every dependency able to control, observe, change, decrypt or recover it |
| DEC-0056 | ADH-2026-039 | Release compatibility is directed; recovery actions and lifecycle-reference contracts are explicit and pinned |
| DEC-0057 | ADH-2026-040 | External maintenance has explicit notice authority, target lifecycle ownership, epochs, fences and mandatory requalification |
| DEC-0058 | ADH-2026-041 | Alpha migration is one signed cutover with no dual write/authority and immutable history preservation |

## 3. Rebaselined FEATURE-0015–0026 sequence

| Order | Feature | Scope and sole ownership | Depends on | Phase 2R exit evidence |
|---:|---|---|---|---|
| 1 | **FEATURE-0015 Canonical Cloud Model and Alpha Migration Foundation** | Seven-scope migration; CloudPlatform; CloudProvider; CloudProviderParticipation; HostingLocation; Datacenter; FaultDomain; externally operated InfrastructureStack; ExecutionTarget identity; CanonicalMigrationPlan/Record; no ResourcePool/ProviderCapability | FEATURE-0011–0014; ADH-020/024/025/037/041 | Deterministic dry-run conversion of Provider/topology fixtures; zero ambiguous authorities; old and new writers never coexist; provider credential isolation |
| 2 | **FEATURE-0016 Adapter Boundary and ExecutionTarget Qualification** | Core adapter interfaces; normalized target facts; target qualification and availability axes; target-scoped SecretRefs; fake adapter | FEATURE-0015; DEC-0036 | One fake target qualifies and one is denied; no provider-native fields in canonical/customer contracts |
| 3 | **FEATURE-0017 Policy Evaluation Abstraction** | PolicyEvaluationRequest/Result, PolicyEngineAdapter and DecisionRecord linkage; deterministic bootstrap fake only | FEATURE-0013, FEATURE-0016 | Allow, deny, unknown and adapter-failure fixtures; no embedded custom policy engine |
| 4 | **FEATURE-0018 Governance, IAM, Approval and Exception Foundation** | GovernanceProfile, Membership, RoleDefinition, RoleAssignment, PrivilegedAccessRequest, AccessReview, ApprovalPolicy, ApprovalRequest and ExceptionGrant | FEATURE-0012, FEATURE-0017 | Least privilege, scoped assignment, JIT expiry, separation of duties, no persona-based authorization |
| 5 | **FEATURE-0019 Sovereignty Facts, Evidence and Policy Foundation** | SovereigntyProfile, RegulatoryPolicyBundle, SovereigntyFactSet, EvidenceRecord and governed dimension registry | FEATURE-0013, FEATURE-0016 | Fresh/stale/missing/tampered evidence fixtures; geography alone never proves sovereignty |
| 6 | **FEATURE-0020 Assignment and Effective Governance Resolution** | ProfileAssignment and immutable EffectiveGovernanceContext with deterministic inheritance/conflict/exception resolution | FEATURE-0018, FEATURE-0019 | Exact input versions, non-weakenable controls, conflict and expired-exception tests |
| 7 | **FEATURE-0021 CloudEnrollment, Personal Onboarding, Entitlement and Quota** | CloudEnrollment; Personal Organization/default Tenant/default Project; EntitlementPackage; ServiceEntitlement; QuotaPolicy and reservation; provider-selection intent vocabulary | FEATURE-0015, FEATURE-0020 | Simple personal and enterprise onboarding; entitlement and quota independently deny; customer/provider boundaries remain independent |
| 8 | **FEATURE-0022 Service Product, Runtime and Requirement Foundation** | ServicePortfolio; ServiceTypeDefinition; CloudPlatform-scoped ServiceOffering/ServicePlan migration; ServicePlacementProfile; ServiceRuntimeProfile; immutable ServiceRequirementSet | FEATURE-0015, FEATURE-0016, FEATURE-0021; ADH-023/027/032 | Synthetic service added without core schema change; published versions immutable; old catalog references map deterministically |
| 9 | **FEATURE-0023 Sovereignty and Placement Decision v0** | Registered sovereignty/placement DecisionProfiles; authoritative DecisionRecords; candidate-set evaluation; provider selection; ServicePlacement safe projection | FEATURE-0019–0022 | Selected, denied, requires-approval and indeterminate flows; exact evidence/policy versions; no sensitive topology leakage |
| 10 | **FEATURE-0024 Plugin Taxonomy and Synthetic Execution Boundary** | Plugin/adapter roles; manifests; compatibility; protected handles; minimal PluginExecution contract; deterministic fake plugin/adapter harness only | FEATURE-0016, FEATURE-0022, FEATURE-0023 | Fake execution cannot bypass final decision, Operation, target authority or SecretRef boundary; no external side effects |
| 11 | **FEATURE-0025 AI-Readable Decision and Operation Context** | Bounded audience-safe explanation projection derived from decisions, evidence, operations and audit | FEATURE-0013, FEATURE-0023, FEATURE-0024 | AI explains authorized facts but cannot decide, approve, change plans or execute; leakage tests pass |
| 12 | **FEATURE-0026 Slice 0 Integration and Conformance Demo** | Cross-feature orchestration of one synthetic service through fake execution, safe placement and fake binding | All Phase 2R features | VS-000 definition of done and conformance matrix pass |

## 4. Dependency graph

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

FEATURE-0019 may proceed in parallel with FEATURE-0017/0018 after FEATURE-0016. No other dependency is relaxed without an approved amendment.

## 5. Phase 2R control flow

```text
Synthetic ServiceInstance request
  → authenticate and authorize through fake adapter
  → CloudEnrollment, ServiceEntitlement and Quota reservation
  → EffectiveGovernanceContext
  → immutable ServiceRequirementSet
  → qualified fake ExecutionTarget candidates
  → sovereignty DecisionRecord
  → placement DecisionRecord
  → immutable ServiceDeploymentPlan
  → accepted Operation
  → fake PluginExecution
  → verified synthetic readiness
  → customer-safe ServicePlacement
  → fake per-consumer ServiceBinding with SecretRef only
  → correlated AuditEvents and safe AI-readable explanation
```

The fake execution harness is conformance machinery. It performs no provider or workload side effects and does not move Phase 3 execution into Phase 2R.

## 6. Explicit non-goals

- Real AWS, OCI, OpenStack, VMware, Kubernetes or OpenShift calls.
- Real PostgreSQL provisioning or lifecycle logic.
- Production identity, policy, secret, workflow, observability or evidence engines.
- Capacity scheduling, ResourcePool inventory or provider-wide capability truth.
- Multi-target HA realization, failover, migration or data movement execution.
- Platform installation/upgrade implementation; Phase 2R only preserves its contracts and topology decisions for the parallel production stream.
- Billing, marketplace, UI and autonomous AI execution.
- Stable API promotion.

## 7. Phase 2R exit criteria

1. ADH-2026-042, ACR-2026-001 and DEC-0037–0058 are accepted in the repository.
2. Current baseline, architecture version, glossary, phase scope, spine, sequence, roadmap, feature index and Kiro context agree.
3. Completed feature history is retained; active summaries point to the coordinated migration.
4. FEATURE-0015 migration dry run is deterministic and reports zero unresolved authorities or references.
5. No active desired-state API supports old and canonical kinds/scopes simultaneously.
6. Every Phase 2R feature passes its reuse, architecture, drift, quality and human gates.
7. Slice 0 passes positive, denial, stale-input, authorization, idempotency, fencing, redaction and no-side-effect conformance.
8. DecisionRecord and AuditEvent semantics remain owned by FEATURE-0013 and are not redefined.
9. Customer contracts contain no provider-native objects, raw secrets or protected handles.
10. Phase 3 readiness review approves replacement of fake execution with one real PostgreSQL path.

## 8. Phase 3 handoff boundary

Phase 3 begins only after Phase 2R exit. It retains FEATURE-0027–0034 but reinterprets them as:

1. production-shaped PluginExecution attempt/transport contract;
2. Operation orchestration controller;
3. PostgreSQL management plugin;
4. OpenShift execution adapter;
5. mature PostgreSQL operator integration;
6. service realization pipeline;
7. workload identity and SecretRef binding;
8. PostgreSQL end-to-end conformance demo.

Phase 3 does not redesign the customer, governance, sovereignty, decision or target contracts proven by Slice 0.
