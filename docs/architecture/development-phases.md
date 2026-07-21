---
doc_type: architecture
title: Sovrunn Development Phases
status: draft
phase: 2
ai_load_priority: always
ai_summary: Reuse-first phased plan for Sovrunn from Phase 0/1 baseline through Phase 2 provider-neutral PaaS fabric foundation, Phase 3 executable PostgreSQL plugin chain, and later production capabilities.
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

Phase 2 and Phase 3 assume:

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
| 2 | Reuse-First PaaS Fabric Foundation | Build model, adapter, policy-context, decision, audit, plugin taxonomy, and placement simulation foundation. |
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

## 7. Phase 2: Reuse-First PaaS Fabric Foundation

### Goal

Phase 2 makes Sovrunn excellent at:

```text
modeling
validating
deciding
explaining
auditing
creating adapter boundaries
```

Phase 2 must not perform real provider provisioning or real PostgreSQL runtime provisioning.

### In Scope

- reuse assessment standard,
- API/resource standard,
- decision and audit standard,
- provider-neutral resource model,
- ResourcePool and ProviderCapability model,
- adapter boundaries for mature OSS reuse,
- policy evaluation abstraction,
- governance/security/data/cost policy models,
- ProfileAssignment and EffectivePolicyContext,
- minimal entitlement/quota placeholder,
- ServiceRuntimeProfile,
- PlacementRequest and PlacementDecision v0,
- plugin taxonomy foundation,
- AI-readable decision context,
- Phase 2 integration simulation.

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
- production-grade plugin sandbox.

### Phase 2 Feature Sequence

| Feature | Name | Purpose |
|---|---|---|
| FEATURE-0011 | Reuse Assessment Standard | Force every feature to decide Reuse / Wrap / Extend / Build. |
| FEATURE-0012 | API, Resource Naming, Status, and Validation Standard | Establish Kubernetes-inspired resource conventions, status, conditions, references, validation errors, and API boundary classification. |
| FEATURE-0013 | Decision Object and AuditEvent Standard | Define common decision, reason, rejected alternative, suggested action, and audit event structure. |
| FEATURE-0014 | Provider-Neutral Resource Model | Define Provider, ProviderLocation/Region, ProviderDatacenter, DatacenterFailureDomain, and IaaSStack. |
| FEATURE-0015 | ResourcePool and ProviderCapability Model | Define ResourcePool as placement boundary and ProviderCapability as compatibility boundary. |
| FEATURE-0016 | Adapter Boundary Foundation | Define adapter interfaces for policy, identity, secrets, operations, observability, events, and repositories. |
| FEATURE-0017 | Policy Evaluation Abstraction | Define PolicyEvaluationRequest/Result and OPA/Cedar-ready PolicyEngineAdapter. |
| FEATURE-0018 | GovernanceProfile and SecurityProfile Foundation | Define policy profiles and profile references. |
| FEATURE-0019 | DataPlacementPolicy and CostGuardrail Minimal Foundation | Define minimal data residency, movement, and cost guardrail inputs. |
| FEATURE-0020 | ProfileAssignment and EffectivePolicyContext | Resolve effective policy context for Organization, Tenant, Project, and ServiceInstance requests. |
| FEATURE-0021 | Minimal ServiceEntitlement and Quota Placeholder | Check whether tenant/project may request a ServiceClass/ServicePlan. |
| FEATURE-0022 | ServiceRuntimeProfile Foundation | Map customer-facing ServicePlan to required runtime capabilities. |
| FEATURE-0023 | PlacementRequest and PlacementDecision v0 | Match ServiceRuntimeProfile, effective policy, entitlement, and ResourcePool capabilities with explainable results. |
| FEATURE-0024 | Plugin Taxonomy Foundation | Define provider/substrate, service management, runtime, traffic, backup, observability, security, compliance evidence, and AI-operations plugin types. |
| FEATURE-0025 | AI-Readable Decision Context | Create structured context for AI explanations without making AI an execution dependency. |
| FEATURE-0026 | Phase 2 Integration Demo | Demonstrate provider/resource/policy/runtime/placement/audit/explanation simulation. |

### Phase 2 Acceptance

- A provider and resource pools can be modeled.
- Capabilities can be declared with status.
- Policy profiles can be assigned and resolved into an EffectivePolicyContext.
- A PostgreSQL ServiceRuntimeProfile can be mapped to required capabilities.
- PlacementDecision can return ALLOWED or DENIED with structured reasons and alternatives.
- AuditEvent is created for meaningful decisions.
- AI-readable DecisionContext can explain allowed/denied placement.
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
FEATURE-0061 ResourcePool Capacity Model
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
FEATURE-0096 SDE ServiceClass
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
