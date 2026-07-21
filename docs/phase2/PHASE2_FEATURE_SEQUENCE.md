---
doc_type: feature_sequence
title: Phase 2 Feature Sequence
status: draft
phase: 2
ai_load_priority: always
ai_summary: Dependency-safe Phase 2 feature sequence for Me + AI chain execution.
---

# Phase 2 Feature Sequence

| Order | Feature | Depends On | Output |
|---:|---|---|---|
| 1 | FEATURE-0011 Reuse Assessment Standard | Phase 1 baseline | Reuse/Wrap/Extend/Build gate for every feature. |
| 2 | FEATURE-0012 API, Resource Naming, Status, and Validation Standard | FEATURE-0011 | Common resource shape and API boundary classification. |
| 3 | FEATURE-0013 Decision Object and AuditEvent Standard | FEATURE-0012 | Common decision/audit schema. |
| 4 | FEATURE-0014 Provider-Neutral Resource Model | FEATURE-0012 | Provider, location, datacenter, failure domain, IaaS stack. |
| 5 | FEATURE-0015 ResourcePool and ProviderCapability Model | FEATURE-0014 | Placement and compatibility boundary. |
| 6 | FEATURE-0016 Adapter Boundary Foundation | FEATURE-0012, FEATURE-0013 | Interfaces/placeholders for mature OSS reuse. |
| 7 | FEATURE-0017 Policy Evaluation Abstraction | FEATURE-0013, FEATURE-0016 | OPA/Cedar-ready policy evaluation contract. |
| 8 | FEATURE-0018 GovernanceProfile and SecurityProfile Foundation | FEATURE-0017 | Policy profile models. |
| 9 | FEATURE-0019 DataPlacementPolicy and CostGuardrail Minimal Foundation | FEATURE-0017 | Data/cost policy input models. |
| 10 | FEATURE-0020 ProfileAssignment and EffectivePolicyContext | FEATURE-0018, FEATURE-0019 | Resolved effective policy context. |
| 11 | FEATURE-0021 Minimal ServiceEntitlement and Quota Placeholder | FEATURE-0020 | Tenant/project service allowance checks. |
| 12 | FEATURE-0022 ServiceRuntimeProfile Foundation | FEATURE-0015 | ServicePlan-to-capability requirements. |
| 13 | FEATURE-0023 PlacementRequest and PlacementDecision v0 | FEATURE-0013, FEATURE-0015, FEATURE-0020, FEATURE-0021, FEATURE-0022 | Explainable allow/deny placement simulation. |
| 14 | FEATURE-0024 Plugin Taxonomy Foundation | FEATURE-0016, FEATURE-0022 | Provider/service/runtime plugin taxonomy and manifest/profile models. |
| 15 | FEATURE-0025 AI-Readable Decision Context | FEATURE-0013, FEATURE-0023 | Structured AI explanation input/output shape. |
| 16 | FEATURE-0026 Phase 2 Integration Demo | all Phase 2 features | End-to-end placement simulation demo. |

## Rule

Do not start `PlacementDecision` before ResourcePool, ProviderCapability, EffectivePolicyContext, ServiceEntitlement, and ServiceRuntimeProfile exist.
