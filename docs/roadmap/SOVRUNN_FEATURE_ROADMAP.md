---
doc_type: roadmap
title: Sovrunn Feature Roadmap
status: draft
phase: roadmap
ai_load_priority: important
ai_summary: Roadmap-level feature placeholders for all Sovrunn phases. Phase 2 is executable detail; Phase 3 is MVP planning detail; Phase 4+ are scope references and must be revalidated after Phase 2 and Phase 3.
---

# Sovrunn Feature Roadmap

## 1. Purpose

This document records the full Sovrunn feature roadmap at a scope-reference level so AI agents, Kiro, Cursor, and reviewers understand the long-term product direction.

This file is **not** a detailed design for every feature.

```text
Phase 2: detailed enough for immediate execution.
Phase 3: detailed enough for MVP planning.
Phase 4+: roadmap placeholders only.
```

Before starting Phase 4 or any later phase, this roadmap must be reviewed and re-baselined using the outcomes of Phase 2, Phase 3, customer feedback, technical learning, and market validation.

## 2. Roadmap Governance Rules

```text
Do not generate detailed Kiro specs for all future features upfront.
Do not treat future placeholders as committed implementation design.
Do not pull future-phase behavior into Phase 2 or Phase 3 unless a formal decision approves it.
Use this roadmap to avoid architectural dead ends, not to overbuild early phases.
```

Each feature still requires a feature-specific Architecture Contract, Reuse Assessment, Kiro requirements/design/tasks, Cursor implementation, automated review, tests, and final architecture acceptance before development.

## 3. Current MVP Anchor

The current MVP anchor is:

```text
MVP-001: Governed PostgreSQL PaaS Placement and Provisioning on one substrate.
```

This MVP validates Sovrunn's core value proposition:

```text
Customer expresses a PaaS outcome.
Provider exposes infrastructure capability.
Sovrunn governs, places, provisions through reusable components, explains, and audits.
```

## 4. Phase 2: Reuse-First PaaS Fabric Foundation

Phase 2 features are executable Phase 2 scope and should be developed in order unless a formal architecture decision changes the order.

| Feature | Name | Scope Level | Purpose |
|---|---|---|---|
| FEATURE-0011 | Reuse Assessment Standard | Executable | Require Reuse / Wrap / Extend / Build decision for every feature. |
| FEATURE-0012 | API, Resource Naming, Status, and Validation Standard | Executable | Establish API/resource conventions, status, conditions, references, validation, and API boundary classification. |
| FEATURE-0013 | Decision Object and AuditEvent Standard | Executable | Define common decision and audit event structure. |
| FEATURE-0014 | Provider-Neutral Resource Model | Executable | Define Provider, ProviderLocation/Region, ProviderDatacenter, DatacenterFailureDomain, and IaaSStack. |
| FEATURE-0015 | ResourcePool and ProviderCapability Model | Executable | Define ResourcePool as placement boundary and ProviderCapability as compatibility boundary. |
| FEATURE-0016 | Adapter Boundary Foundation | Executable | Define adapter interfaces for policy, identity, secrets, operations, observability, events, and repositories. |
| FEATURE-0017 | Policy Evaluation Abstraction | Executable | Define OPA/Cedar-ready policy input, evaluation result, and engine adapter contracts. |
| FEATURE-0018 | GovernanceProfile and SecurityProfile Foundation | Executable | Define governance and security profile objects. |
| FEATURE-0019 | DataPlacementPolicy and CostGuardrail Minimal Foundation | Executable | Define minimal data residency, movement, and cost guardrail inputs. |
| FEATURE-0020 | ProfileAssignment and EffectivePolicyContext | Executable | Resolve effective policy context for Organization, Tenant, Project, and ServiceInstance requests. |
| FEATURE-0021 | Minimal ServiceEntitlement and Quota Placeholder | Executable | Validate that a tenant/project may request a ServiceClass/ServicePlan. |
| FEATURE-0022 | ServiceRuntimeProfile Foundation | Executable | Map customer-facing ServicePlan to runtime/capability requirements. |
| FEATURE-0023 | PlacementRequest and PlacementDecision v0 | Executable | Evaluate resource pools against runtime, policy, entitlement, and capability requirements. |
| FEATURE-0024 | Plugin Taxonomy Foundation | Executable | Define plugin types and boundaries for provider, service management, runtime, traffic, backup, observability, security, evidence, and AI operations. |
| FEATURE-0025 | AI-Readable Decision Context | Executable | Create structured explanation context for allowed/denied decisions. |
| FEATURE-0026 | Phase 2 Integration Demo | Executable | Demonstrate provider/resource/policy/runtime/placement/audit/explanation simulation. |

## 5. Phase 3: First Executable PaaS Plugin Chain

Phase 3 is the first real execution phase. It remains narrow and proves one plugin chain.

| Feature | Name | Scope Level | Purpose |
|---|---|---|---|
| FEATURE-0027 | Plugin Execution Contract v0 | MVP Planning | Define plugin execution request/result/status and operation linkage. |
| FEATURE-0028 | Operation Controller v0 | MVP Planning | Track operation steps, status, retry placeholder, approval placeholder, and audit linkage behind OperationEngineAdapter. |
| FEATURE-0029 | PostgreSQL Management Plane Plugin v0 | MVP Planning | Plan PostgreSQL service lifecycle using wrappers around mature PostgreSQL runtime tooling. |
| FEATURE-0030 | Kubernetes / Local Substrate Plugin v0 | MVP Planning | Execute one local/k3s/Kubernetes path through Kubernetes APIs, Helm, or operator CR wrappers. |
| FEATURE-0031 | PostgreSQL Runtime Plugin v0 | MVP Planning | Wrap runtime create/readiness/binding/endpoint/delete behavior through reused PostgreSQL operator or Helm flow. |
| FEATURE-0032 | ServiceInstance Provisioning v0 | MVP Planning | Convert approved PlacementDecision into Operation and plugin execution. |
| FEATURE-0033 | ServiceBinding and SecretRef Integration | MVP Planning | Create binding using SecretRef/CredentialRef without storing raw credentials in Sovrunn. |
| FEATURE-0034 | Phase 3 End-to-End MVP Demo | MVP Planning | Demonstrate governed PostgreSQL placement and provisioning on one substrate. |

## 6. Phase 4: Customer-Testable MVP Hardening

Phase 4 turns the Phase 3 vertical slice into a customer/design-partner validation package.

| Feature | Name | Scope Level | Purpose |
|---|---|---|---|
| FEATURE-0035 | Persistent Registry Backend | Roadmap Placeholder | Add persistent storage behind repository interfaces. |
| FEATURE-0036 | Minimal Auth/RBAC Adapter Integration | Roadmap Placeholder | Integrate minimal authorization path without building custom IAM. |
| FEATURE-0037 | Customer Demo CLI/API Flow | Roadmap Placeholder | Package the MVP flow for customer-facing API/CLI demonstration. |
| FEATURE-0038 | Integration Test Suite | Roadmap Placeholder | Add deterministic end-to-end tests for allowed/denied/provision/binding/audit flows. |
| FEATURE-0039 | Security/Lint/Gosec/Race Gate | Roadmap Placeholder | Formalize security and quality gates for MVP release. |
| FEATURE-0040 | Pilot Demo Packaging | Roadmap Placeholder | Package local demo, sample configs, runbooks, and validation guide. |
| FEATURE-0041 | Customer Feedback Capture Template | Roadmap Placeholder | Capture structured customer feedback against MVP hypotheses. |

## 7. Phase 5: Provider / Plugin Framework and Certification

Phase 5 should be revalidated after MVP feedback. It expands from one executable plugin chain into a safer plugin ecosystem.

| Feature | Name | Scope Level | Purpose |
|---|---|---|---|
| FEATURE-0042 | Plugin Manifest Validation | Roadmap Placeholder | Validate plugin metadata, types, versions, and declared capabilities. |
| FEATURE-0043 | Plugin TrustProfile Enforcement | Roadmap Placeholder | Enforce trust boundaries and allowed plugin capabilities. |
| FEATURE-0044 | Plugin CredentialPolicy Integration | Roadmap Placeholder | Connect plugin credentials to SecretRef/CredentialRef and approved secret providers. |
| FEATURE-0045 | Provider Capability Validation Workflow | Roadmap Placeholder | Move capabilities from declared to validated/certified/degraded/disabled. |
| FEATURE-0046 | Plugin Certification Test Harness | Roadmap Placeholder | Run conformance tests for provider, service management, and runtime plugins. |
| FEATURE-0047 | Provider Onboarding Workflow | Roadmap Placeholder | Guide provider/MSP through provider, location, IaaSStack, ResourcePool, and capability onboarding. |
| FEATURE-0048 | Plugin Versioning and Compatibility Checks | Roadmap Placeholder | Manage plugin compatibility with Sovrunn API and resource versions. |
| FEATURE-0049 | Plugin Health and Degradation Model | Roadmap Placeholder | Represent plugin health, degraded capability states, and disabled execution paths. |

## 8. Phase 6: Resilience, Traffic, and Data Movement

Phase 6 adds cross-location thinking after the base MVP is proven.

| Feature | Name | Scope Level | Purpose |
|---|---|---|---|
| FEATURE-0050 | ResilienceGroup Foundation | Roadmap Placeholder | Model execution-location groupings for HA, DR, failover, and cloudbursting. |
| FEATURE-0051 | DRProfile Foundation | Roadmap Placeholder | Define customer-facing recovery objectives and modes. |
| FEATURE-0052 | ReplicationPolicy Foundation | Roadmap Placeholder | Model sync/async/semi-sync/snapshot/backup-copy replication choices. |
| FEATURE-0053 | NetworkConnectivityProfile | Roadmap Placeholder | Model network latency, private connectivity, and routing capability between locations. |
| FEATURE-0054 | GlobalTrafficPolicy Foundation | Roadmap Placeholder | Model traffic routing, failover, weighted routing, and active/passive modes. |
| FEATURE-0055 | TrafficDecision Foundation | Roadmap Placeholder | Explain traffic routing decisions and denials. |
| FEATURE-0056 | FailoverDecision Foundation | Roadmap Placeholder | Govern failover approval, risk, data lag, and recovery status. |
| FEATURE-0057 | DataMovementDecision v1 | Roadmap Placeholder | Evaluate and explain data movement across locations/providers. |
| FEATURE-0058 | Cross-Location Placement Simulation | Roadmap Placeholder | Simulate cross-location placement without full production execution. |

## 9. Phase 7: Autoscaling, Capacity, Cost, and Spot

Phase 7 adds governed scaling and cost-awareness.

| Feature | Name | Scope Level | Purpose |
|---|---|---|---|
| FEATURE-0059 | AutoscalingPolicy Foundation | Roadmap Placeholder | Model scale triggers, safe actions, approval modes, and limits. |
| FEATURE-0060 | CapacityPolicy and CapacityClass | Roadmap Placeholder | Model on-demand, reserved, dedicated, spot, burstable, and committed capacity. |
| FEATURE-0061 | ResourcePool Capacity Model | Roadmap Placeholder | Track capacity availability, reservation, and exhaustion signals. |
| FEATURE-0062 | CostEstimate Foundation | Roadmap Placeholder | Estimate cost impact of provisioning/scaling decisions through reusable cost sources. |
| FEATURE-0063 | CostGuardrail v1 | Roadmap Placeholder | Enforce cost limits, approval thresholds, and budget risk reasons. |
| FEATURE-0064 | ScalingDecision Foundation | Roadmap Placeholder | Explain allowed/denied scaling decisions and alternatives. |
| FEATURE-0065 | SpotInterruptionPolicy | Roadmap Placeholder | Model interruption handling, drain, replacement, and fallback. |
| FEATURE-0066 | InterruptionEvent Model | Roadmap Placeholder | Record and react to spot/preemptible interruption events. |
| FEATURE-0067 | Autoscaling Simulation Demo | Roadmap Placeholder | Demonstrate safe scale-out decisioning without broad production automation. |

## 10. Phase 8: Compliance Evidence and Sovereign Assurance

Phase 8 adds evidence and control mapping, not legal certification claims.

| Feature | Name | Scope Level | Purpose |
|---|---|---|---|
| FEATURE-0068 | ComplianceProfile Foundation | Roadmap Placeholder | Model compliance intent and associated control objectives. |
| FEATURE-0069 | ControlObjective and ControlMapping | Roadmap Placeholder | Map policies, evidence, and controls. |
| FEATURE-0070 | EvidenceRecord v1 | Roadmap Placeholder | Store proof that a control was evaluated or satisfied. |
| FEATURE-0071 | EvidenceCollectorAdapter | Roadmap Placeholder | Wrap existing evidence sources instead of building a GRC system. |
| FEATURE-0072 | ComplianceDecision Foundation | Roadmap Placeholder | Explain compliance-related allow/deny/warn decisions. |
| FEATURE-0073 | ExceptionRecord Model | Roadmap Placeholder | Track scoped, approved, time-bound exceptions. |
| FEATURE-0074 | AttestationReport Foundation | Roadmap Placeholder | Export technical assurance reports for operators/customers. |
| FEATURE-0075 | AuditExport Foundation | Roadmap Placeholder | Export audit trails for customer and regulator review. |
| FEATURE-0076 | Sovereign Assurance Demo | Roadmap Placeholder | Demonstrate evidence-backed sovereign policy enforcement. |

## 11. Phase 9: AI-Assisted Operations

Phase 9 adds recommendation and controlled autonomy. AI must not bypass policy.

| Feature | Name | Scope Level | Purpose |
|---|---|---|---|
| FEATURE-0077 | AIOperationRecommendation | Roadmap Placeholder | Generate structured recommendations from decisions, operations, health, and evidence. |
| FEATURE-0078 | RunbookPlan Foundation | Roadmap Placeholder | Generate and store operator-readable runbook plans. |
| FEATURE-0079 | RemediationPlan Foundation | Roadmap Placeholder | Propose remediation actions with risk and approval context. |
| FEATURE-0080 | RiskAssessment Model | Roadmap Placeholder | Attach risk level and rationale to recommendations. |
| FEATURE-0081 | HumanApprovalGate | Roadmap Placeholder | Require human/operator approval for risky actions. |
| FEATURE-0082 | AutonomyPolicy | Roadmap Placeholder | Define what AI-assisted automation may do automatically, with approval, or never. |
| FEATURE-0083 | AI Operation Memory Foundation | Roadmap Placeholder | Store operational learnings without hiding decisions or bypassing audit. |
| FEATURE-0084 | AI Explanation and Recommendation API | Roadmap Placeholder | Expose explanations and recommendations through API/CLI/portal. |
| FEATURE-0085 | AI-Assisted Operations Demo | Roadmap Placeholder | Demonstrate AI recommendation without autonomous execution bypass. |

## 12. Phase 10: Multi-Service PaaS Beta

Phase 10 broadens services only after PostgreSQL MVP validation.

| Feature | Name | Scope Level | Purpose |
|---|---|---|---|
| FEATURE-0086 | Redis / Dragonfly Service Plugin | Roadmap Placeholder | Add cache service through service management/runtime plugin wrappers. |
| FEATURE-0087 | Object Storage Service Plugin | Roadmap Placeholder | Add object storage service abstraction and lifecycle wrapper. |
| FEATURE-0088 | Kafka / Streaming Service Plugin | Roadmap Placeholder | Add streaming service through mature operator/tool wrappers. |
| FEATURE-0089 | Vector Database Service Plugin | Roadmap Placeholder | Add vector database service through reusable engines. |
| FEATURE-0090 | AI Inference Service Plugin | Roadmap Placeholder | Add governed inference service placement and runtime wrapper. |
| FEATURE-0091 | Multi-Service Catalog Experience | Roadmap Placeholder | Improve service catalog for multiple PaaS offerings. |
| FEATURE-0092 | Service Entitlement v2 | Roadmap Placeholder | Expand entitlement and quota beyond MVP placeholders. |
| FEATURE-0093 | Service Dependency Graph | Roadmap Placeholder | Model dependencies between PaaS services. |
| FEATURE-0094 | Multi-Service Provisioning Demo | Roadmap Placeholder | Demonstrate multiple service provisioning flows. |
| FEATURE-0095 | Multi-Service Lifecycle Validation | Roadmap Placeholder | Validate upgrade, backup, binding, and delete flows across services. |

## 13. Phase 11: SDE as Managed Service

Phase 11 brings Sovrunn Data Engine into Sovrunn as a governed managed service.

| Feature | Name | Scope Level | Purpose |
|---|---|---|---|
| FEATURE-0096 | SDE ServiceClass | Roadmap Placeholder | Model SDE as a first-class managed service. |
| FEATURE-0097 | SDE ServiceRuntimeProfile | Roadmap Placeholder | Define SDE runtime requirements and capabilities. |
| FEATURE-0098 | PostgreSQL Wire Gateway Service Plugin | Roadmap Placeholder | Wrap SDE gateway deployment and lifecycle. |
| FEATURE-0099 | Metadata Store Integration | Roadmap Placeholder | Integrate SDE metadata store requirements. |
| FEATURE-0100 | Object Storage Offload Integration | Roadmap Placeholder | Integrate SDE blob/object offload requirements. |
| FEATURE-0101 | Cache Integration | Roadmap Placeholder | Integrate SDE cache requirements without storing blob payloads in cache. |
| FEATURE-0102 | SDE PlacementDecision | Roadmap Placeholder | Extend placement decisioning for SDE-specific topology. |
| FEATURE-0103 | SDE ServiceInstance Provisioning | Roadmap Placeholder | Provision SDE as a governed managed service. |
| FEATURE-0104 | SDE Observability and Audit | Roadmap Placeholder | Add SDE-specific health, operation, and audit context. |
| FEATURE-0105 | SDE MVP Demo | Roadmap Placeholder | Demonstrate SDE inside Sovrunn after platform MVP validation. |

## 14. Phase 12: Production Beta / Enterprise Readiness

Phase 12 hardens the platform for controlled production beta.

| Feature | Name | Scope Level | Purpose |
|---|---|---|---|
| FEATURE-0106 | Production Multi-Tenant Control Plane | Roadmap Placeholder | Harden control plane for production tenancy and isolation. |
| FEATURE-0107 | Upgrade and Migration Framework | Roadmap Placeholder | Support resource, schema, and plugin version migration. |
| FEATURE-0108 | Backup and Restore for Control Plane | Roadmap Placeholder | Protect Sovrunn control plane state. |
| FEATURE-0109 | HA Control Plane Deployment | Roadmap Placeholder | Deploy Sovrunn control plane in highly available mode. |
| FEATURE-0110 | Tenant Isolation Hardening | Roadmap Placeholder | Validate tenant boundaries and least privilege. |
| FEATURE-0111 | Security Threat Model Validation | Roadmap Placeholder | Complete security validation and threat model closure. |
| FEATURE-0112 | Load and Scale Testing | Roadmap Placeholder | Validate control plane and MVP service scalability. |
| FEATURE-0113 | Chaos and Failure Testing | Roadmap Placeholder | Validate failure scenarios and recovery behavior. |
| FEATURE-0114 | Supportability and Diagnostics | Roadmap Placeholder | Add logs, diagnostics, support bundles, and runbooks. |
| FEATURE-0115 | Production Beta Release Gate | Roadmap Placeholder | Define beta acceptance, known limitations, and release criteria. |

## 15. Rebaseline Rule

Before starting Phase 4 and later phases, run a rebaseline review:

```text
1. What did Phase 2 prove or disprove?
2. What did Phase 3 prove or disprove?
3. What did customer/design-partner validation prove or disprove?
4. Which roadmap placeholders should be kept, split, renamed, moved, or removed?
5. Which mature OSS components should now be selected for adapters?
6. Which architectural assumptions remain valid?
```

Do not treat this roadmap as frozen beyond Phase 3.
