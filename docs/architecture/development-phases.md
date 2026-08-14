---
doc_type: architecture
title: Sovrunn Development Phases
status: approved
phase: 2R
ai_load_priority: always
ai_summary: 'Reuse-first phased plan from Phase 0/1 through Phase 2R canonical model foundation, Phase 3 executable PostgreSQL plugin chain, and later capabilities.'
---

# Sovrunn Development Phases

## 1. Purpose

This document defines the current Sovrunn phased delivery plan after the product scope evolved from a Phase 1 platform skeleton into a reuse-first, AI-first, sovereign cloud-native PaaS platform.

AI agents must use this file to avoid implementing future-phase features too early.

## 2. Operating Rule

```text
Design ahead.
Implement narrowly.
Reuse mature open-source systems.
Build Sovrunn-specific governance, decisions, audit, orchestration, plugin contracts, and AI-readable context.
```

## 3. Current Execution Model

Phase 2R and Phase 3 assume:

```text
Human owner + ChatGPT = architecture contract and final acceptance
Kiro = feature requirements, design, and tasks
Cursor = implementation and tests
Automated reviewer = quality, security, and architecture drift checks
```

## 4. Phase Summary

| Phase | Name | Primary Goal |
|---:|---|---|
| 0 | Foundation and AI Development System | Make project AI-developable without architecture drift. |
| 1 | Platform Core Skeleton | Build core Sovrunn resource grammar. |
| 2R | Canonical Model PaaS Fabric Foundation | Build implementation-neutral model, adapter, governance, sovereignty, decision, audit, plugin taxonomy, and placement simulation foundation under the canonical model. |
| 3 | First Executable PaaS Plugin Chain | Execute one governed PostgreSQL provisioning path on one substrate by wrapping mature components. |
| 4 | Customer-Testable MVP Hardening | Package the PostgreSQL PaaS MVP for design-partner/customer validation. |
| 5 | Provider/Plugin Framework and Certification | Formalize provider, service, runtime, traffic, backup, evidence, and observability plugin contracts. |
| 6 | Resilience, Traffic, and Data-Movement Foundation | Add cross-location models, DR profiles, replication policy, traffic decisions, and data movement controls. |
| 7 | Autoscaling, Capacity, Cost, and Spot Foundation | Add governed capacity, cost, scaling, and spot decision models. |
| 8 | Compliance Evidence and Sovereign Assurance | Add evidence records, control mappings, collectors, and attestation outputs. |
| 9 | AI-Assisted Operations | Add AI recommendations, runbooks, risk assessment, and controlled approval gates. |
| 10 | Multi-Service PaaS Beta | Add more service classes after PostgreSQL MVP validation. |
| 11 | SDE as Managed Service | Bring SDE into Sovrunn as a governed managed service. |
| 12 | Production Beta | Harden for controlled production with friendly customers. |

## 5. Phase 0: Foundation and AI Development System

Status: complete or baseline.

Purpose:

```text
Create source-of-truth docs, AI development workflow, decisions, glossary, and feature factory.
```

## 6. Phase 1: Platform Core Skeleton

Status: complete or near complete.

Purpose:

```text
Build core resource grammar: Organization, OrganizationUnit, Tenant, Project, ServiceClass, ServicePlan, ServiceInstance, ServiceBinding, Operation, Plugin, Capability, API health/readiness, and demo flow.
```

Phase 1 documents remain valid as baseline records. They are not the complete Phase 2 scope.

## 7. Phase 2R: Canonical Model PaaS Fabric Foundation

### Goal

Phase 2R makes Sovrunn excellent at:

```text
implementation-neutral modeling
governance and sovereignty resolution
deciding and explaining
auditing with immutable records
creating adapter boundaries
qualifying execution targets
conformance simulation
```

Phase 2R must not perform real provider provisioning, real PostgreSQL runtime provisioning, or any external side effects.

### In Scope

- reuse assessment standard,
- API/resource standard,
- decision and audit standard (DecisionRecord with registered profiles),
- canonical cloud model foundation via direct bootstrap (seven scope kinds, CloudPlatform, CloudProvider; no alpha runtime migration),
- adapter boundary and ExecutionTarget qualification,
- policy evaluation abstraction with DecisionRecord linkage,
- governance, IAM, approval, and exception foundation,
- sovereignty facts, evidence, and policy foundation,
- effective governance resolution (EffectiveGovernanceContext),
- CloudEnrollment, personal onboarding, entitlement, and quota,
- service product, runtime, and requirement foundation (ServiceTypeDefinition, ServiceOffering, ServicePlan, ServiceRequirementSet),
- sovereignty and placement decision v0 (DecisionRecord profiles),
- plugin taxonomy and synthetic execution boundary,
- AI-readable decision and operation context,
- Slice 0 integration and conformance demo (VS-000).

### Out of Scope

- real AWS/Azure/VMware/OpenStack/OpenShift provisioning,
- real PostgreSQL operator integration,
- full OPA/Cedar/Keycloak/Vault/Temporal integration,
- real autoscaling,
- real failover,
- real global traffic management,
- full billing,
- full compliance engine,
- autonomous AI operations,
- production-grade plugin sandbox,
- mandatory capacity scheduling or provider-wide capability truth,
- native provider objects in customer or core schemas,
- runtime alpha data migration, a migration controller, or a cutover state machine (DEC-0059).

### Phase 2R Feature Sequence

#### Completed Features (Implementation History)

| Feature | Name | Purpose |
|---|---|---|
| FEATURE-0011 | Reuse Assessment Standard | Force every feature to decide Reuse / Wrap / Extend / Build. |
| FEATURE-0012 | API, Resource Naming, Status, and Validation Standard | Establish resource conventions, status, conditions, references, validation, and API boundary classification. |
| FEATURE-0013 | Decision Record and AuditEvent Standard | Define DecisionRecord envelope, DecisionProfile extension, AuditEvent, and scope authority. Controlled by ADH-2026-017. |
| FEATURE-0014 | Provider-Neutral Resource Model | Define Provider, ProviderLocation, ProviderDatacenter, DatacenterFailureDomain, and InfrastructureStack. Alpha model; controlled by ADH-2026-018. |

#### Phase 2R Rebaselined Features

| Feature | Name | Purpose |
|---|---|---|
| FEATURE-0015 | Canonical Cloud Model Foundation | Seven-scope canonical bootstrap; CloudPlatform; CloudProvider; CloudProviderParticipation; HostingLocation; Datacenter; FaultDomain; InfrastructureStack (ends here; ExecutionTarget owned by FEATURE-0016); no alpha runtime migration (DEC-0059). |
| FEATURE-0016 | Adapter Boundary and ExecutionTarget Qualification | Core adapter interfaces; normalized target facts; target qualification and availability axes; target-scoped SecretRefs; fake adapter. |
| FEATURE-0017 | Policy Evaluation Abstraction | PolicyEvaluationRequest/Result; PolicyEngineAdapter; DecisionRecord linkage; deterministic bootstrap fake. |
| FEATURE-0018 | Governance, IAM, Approval and Exception Foundation | GovernanceProfile; Membership; RoleDefinition; RoleAssignment; PrivilegedAccessRequest; AccessReview; ApprovalPolicy; ApprovalRequest; ExceptionGrant. |
| FEATURE-0019 | Sovereignty Facts, Evidence and Policy Foundation | SovereigntyProfile; RegulatoryPolicyBundle; SovereigntyFactSet; EvidenceRecord; governed dimension registry. |
| FEATURE-0020 | Assignment and Effective Governance Resolution | ProfileAssignment; immutable EffectiveGovernanceContext; deterministic inheritance/conflict/exception resolution. |
| FEATURE-0021 | CloudEnrollment, Personal Onboarding, Entitlement and Quota | CloudEnrollment; Personal Organization; EntitlementPackage; ServiceEntitlement; QuotaPolicy; provider-selection intent. |
| FEATURE-0022 | Service Product, Runtime and Requirement Foundation | ServicePortfolio; ServiceTypeDefinition; ServiceOffering; ServicePlan migration; ServiceRuntimeProfile; immutable ServiceRequirementSet. |
| FEATURE-0023 | Sovereignty and Placement Decision v0 | Registered sovereignty/placement DecisionProfiles; authoritative DecisionRecords; candidate-set evaluation; ServicePlacement safe projection. |
| FEATURE-0024 | Plugin Taxonomy and Synthetic Execution Boundary | Plugin/adapter roles; manifests; compatibility; protected handles; PluginExecution contract; deterministic fake harness. |
| FEATURE-0025 | AI-Readable Decision and Operation Context | Bounded audience-safe explanation projection from decisions, evidence, operations, and audit. |
| FEATURE-0026 | Slice 0 Integration and Conformance Demo | Cross-feature orchestration of one synthetic service through fake execution, safe placement, and fake binding; VS-000 conformance. |

### Phase 2R Acceptance

- CloudPlatform, CloudProvider, and ExecutionTarget can be modeled with seven canonical scope kinds.
- ExecutionTarget qualification produces qualified/denied outcomes through adapter-provided normalized facts.
- Governance profiles compose into an immutable EffectiveGovernanceContext with non-weakenable controls.
- Sovereignty is assessed through evidence-backed DecisionRecord profiles; geography alone never proves sovereignty.
- A ServiceRequirementSet translates customer-facing ServicePlan into placement requirements.
- Placement DecisionRecord selects, denies, or marks targets as requires-approval with exact evidence/policy versions.
- ServicePlacement is a safe customer projection with no topology or credential leakage.
- ServiceBinding is per-consumer, SecretRef-only, separately revocable.
- AuditEvent is created for meaningful decisions with correlated DecisionRecords.
- AI-readable context explains authorized facts but cannot decide, approve, or execute.
- VS-000 definition of done and conformance matrix pass.
- Every feature includes a Reuse Assessment.

## 8. Phase 3: First Executable PaaS Plugin Chain

### Goal

Convert one approved placement decision into one controlled executable operation.

### MVP Chain

```text
Sovrunn Core decides.
PostgreSQL Management Plane Plugin plans.
Kubernetes/Local Substrate Plugin provisions via mature OSS.
PostgreSQL Runtime Plugin wraps actual runtime/operator/Helm behavior.
Operation tracks lifecycle.
Audit records the result.
AI explains from structured context.
```

### Phase 3 Features

| Feature | Name | Purpose |
|---|---|---|
| FEATURE-0027 | Plugin Execution Contract v0 | Define execution request/result/status and operation linkage. |
| FEATURE-0028 | Operation Controller v0 | Track operation, operation steps, status, approval placeholder, retry placeholder, and audit linkage behind OperationEngineAdapter. |
| FEATURE-0029 | PostgreSQL Management Plane Plugin v0 | Plan PostgreSQL runtime requirements and lifecycle using wrapper logic, not custom PostgreSQL HA. |
| FEATURE-0030 | Kubernetes/Local Substrate Plugin v0 | Execute one local/k3s/Kubernetes substrate path using Kubernetes APIs, Helm, or operator CR wrappers. |
| FEATURE-0031 | PostgreSQL Runtime Plugin v0 | Wrap CloudNativePG/Crunchy/Helm-based runtime actions for create, readiness, binding, endpoint, and delete. |
| FEATURE-0032 | ServiceInstance Provisioning v0 | Link PlacementDecision to Operation and plugin execution, then update ServiceInstance status. |
| FEATURE-0033 | ServiceBinding and SecretRef Integration | Create binding with credentialRef/secretRef, not raw secrets in Sovrunn resources. |
| FEATURE-0034 | Phase 3 End-to-End MVP Demo | Demonstrate governed PostgreSQL placement and provisioning on one substrate. |

### Phase 3 Acceptance

- A customer can request PostgreSQL basic/small service.
- Sovrunn evaluates placement using Phase 2 decision contracts.
- Sovrunn creates an Operation.
- Sovrunn delegates runtime creation to a reused PostgreSQL operator or Helm chart through plugin wrappers.
- Sovrunn creates a ServiceBinding using SecretRef/CredentialRef.
- Sovrunn records AuditEvent.
- Denied and allowed flows are explainable.

## 9. Phase 4: Customer-Testable MVP Hardening

### Goal

Make the MVP suitable for real customer/design-partner validation.

### Features

- persistent registry backend,
- minimal auth/RBAC adapter integration,
- CLI/API customer demo flow,
- integration test suite,
- security/lint/gosec/race gate,
- pilot demo packaging,
- customer feedback capture template.

### MVP Statement

```text
Governed PostgreSQL PaaS Placement and Provisioning on one substrate,
with explainable decisions, security/data policy, audit events,
plugin-chain execution, and AI-readable explanations.
```

## 10. Later Phase Feature Placeholders

Later phase features are roadmap-level placeholders, not detailed implementation design. They exist so Phase 2 and Phase 3 remain aligned with the full Sovrunn scope.

Before starting Phase 4 or any later phase, rebaseline the roadmap using Phase 2 outcomes, Phase 3 MVP results, customer/design-partner feedback, and reuse assessment findings.

The detailed all-phase roadmap is maintained in:

```text
docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md
docs/features/FEATURE_INDEX.md
```

### Phase 4: Customer-Testable MVP Hardening

```text
FEATURE-0035 Persistent Registry Backend
FEATURE-0036 Minimal Auth/RBAC Adapter Integration
FEATURE-0037 Customer Demo CLI/API Flow
FEATURE-0038 Integration Test Suite
FEATURE-0039 Security/Lint/Gosec/Race Gate
FEATURE-0040 Pilot Demo Packaging
FEATURE-0041 Customer Feedback Capture Template
```

### Phase 5: Provider / Plugin Framework and Certification

```text
FEATURE-0042 Plugin Manifest Validation
FEATURE-0043 Plugin TrustProfile Enforcement
FEATURE-0044 Plugin CredentialPolicy Integration
FEATURE-0045 Provider Capability Validation Workflow
FEATURE-0046 Plugin Certification Test Harness
FEATURE-0047 Provider Onboarding Workflow
FEATURE-0048 Plugin Versioning and Compatibility Checks
FEATURE-0049 Plugin Health and Degradation Model
```

### Phase 6: Resilience, Traffic, and Data Movement

```text
FEATURE-0050 ResilienceGroup Foundation
FEATURE-0051 DRProfile Foundation
FEATURE-0052 ReplicationPolicy Foundation
FEATURE-0053 NetworkConnectivityProfile
FEATURE-0054 GlobalTrafficPolicy Foundation
FEATURE-0055 TrafficDecision Foundation
FEATURE-0056 FailoverDecision Foundation
FEATURE-0057 DataMovementDecision v1
FEATURE-0058 Cross-Location Placement Simulation
```

### Phase 7: Autoscaling, Capacity, Cost, and Spot

```text
FEATURE-0059 AutoscalingPolicy Foundation
FEATURE-0060 CapacityPolicy and CapacityClass
FEATURE-0061 Capacity Model (future; ExecutionTarget-based after Phase 2R revalidation)
FEATURE-0062 CostEstimate Foundation
FEATURE-0063 CostGuardrail v1
FEATURE-0064 ScalingDecision Foundation
FEATURE-0065 SpotInterruptionPolicy
FEATURE-0066 InterruptionEvent Model
FEATURE-0067 Autoscaling Simulation Demo
```

### Phase 8: Compliance Evidence and Sovereign Assurance

```text
FEATURE-0068 ComplianceProfile Foundation
FEATURE-0069 ControlObjective and ControlMapping
FEATURE-0070 EvidenceRecord v1
FEATURE-0071 EvidenceCollectorAdapter
FEATURE-0072 ComplianceDecision Foundation
FEATURE-0073 ExceptionRecord Model
FEATURE-0074 AttestationReport Foundation
FEATURE-0075 AuditExport Foundation
FEATURE-0076 Sovereign Assurance Demo
```

### Phase 9: AI-Assisted Operations

```text
FEATURE-0077 AIOperationRecommendation
FEATURE-0078 RunbookPlan Foundation
FEATURE-0079 RemediationPlan Foundation
FEATURE-0080 RiskAssessment Model
FEATURE-0081 HumanApprovalGate
FEATURE-0082 AutonomyPolicy
FEATURE-0083 AI Operation Memory Foundation
FEATURE-0084 AI Explanation and Recommendation API
FEATURE-0085 AI-Assisted Operations Demo
```

### Phase 10: Multi-Service PaaS Beta

```text
FEATURE-0086 Redis / Dragonfly Service Plugin
FEATURE-0087 Object Storage Service Plugin
FEATURE-0088 Kafka / Streaming Service Plugin
FEATURE-0089 Vector Database Service Plugin
FEATURE-0090 AI Inference Service Plugin
FEATURE-0091 Multi-Service Catalog Experience
FEATURE-0092 Service Entitlement v2
FEATURE-0093 Service Dependency Graph
FEATURE-0094 Multi-Service Provisioning Demo
FEATURE-0095 Multi-Service Lifecycle Validation
```

### Phase 11: SDE as Managed Service

```text
FEATURE-0096 SDE ServiceTypeDefinition and ServiceOffering
FEATURE-0097 SDE ServiceRuntimeProfile
FEATURE-0098 PostgreSQL Wire Gateway Service Plugin
FEATURE-0099 Metadata Store Integration
FEATURE-0100 Object Storage Offload Integration
FEATURE-0101 Cache Integration
FEATURE-0102 SDE PlacementDecision
FEATURE-0103 SDE ServiceInstance Provisioning
FEATURE-0104 SDE Observability and Audit
FEATURE-0105 SDE MVP Demo
```

### Phase 12: Production Beta / Enterprise Readiness

```text
FEATURE-0106 Production Multi-Tenant Control Plane
FEATURE-0107 Upgrade and Migration Framework
FEATURE-0108 Backup and Restore for Control Plane
FEATURE-0109 HA Control Plane Deployment
FEATURE-0110 Tenant Isolation Hardening
FEATURE-0111 Security Threat Model Validation
FEATURE-0112 Load and Scale Testing
FEATURE-0113 Chaos and Failure Testing
FEATURE-0114 Supportability and Diagnostics
FEATURE-0115 Production Beta Release Gate
```

### Later-Phase Non-Execution Rule

```text
Future roadmap features may be referenced for scope awareness.
They must not be implemented during Phase 2 or Phase 3 unless a formal decision changes the phase boundary.
```

## 11. Phase Coding Rules

```text
Phase 2: model, decide, explain, audit, define adapter boundaries.
Phase 3: execute one PostgreSQL plugin chain only.
Phase 4: harden for customer validation.
Later: add production-grade engines and more providers/services.
```

Never build broad execution before the decision model and adapter boundaries are stable.
