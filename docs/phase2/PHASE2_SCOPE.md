---
doc_type: phase_scope
title: Phase 2R Scope
status: approved
phase: 2R
ai_load_priority: always
ai_summary: Defines Phase 2R as canonical-model, implementation-neutral, simulation-only PaaS fabric foundation. No real provisioning. No ResourcePool or ProviderCapability.
---

# Phase 2R Scope

## Purpose

Phase 2R establishes the canonical, implementation-neutral PaaS fabric foundation after the ADH-2026-042 model adoption.

Phase 2R builds:

```text
canonical cloud model and alpha migration
adapter boundaries and target qualification
policy evaluation abstraction
governance, IAM, approval, and exception foundation
sovereignty facts, evidence, and policy
effective governance resolution
enrollment, entitlement, and quota
service product, runtime, and requirement foundation
sovereignty and placement decisions
plugin taxonomy and synthetic execution boundary
AI-readable decision context
Slice 0 integration and conformance
```

Phase 2R does not build real provider provisioning, real database provisioning, real autoscaling, real failover, or real autonomous operations.

## In Scope

- Canonical Cloud Model and Alpha Migration Foundation (seven scopes, CloudPlatform, CloudProvider, ExecutionTarget)
- Adapter Boundary and ExecutionTarget Qualification
- Policy Evaluation Abstraction with DecisionRecord linkage
- Governance, IAM, Approval, and Exception Foundation
- Sovereignty Facts, Evidence, and Policy Foundation
- Assignment and Effective Governance Resolution (EffectiveGovernanceContext)
- CloudEnrollment, Personal Onboarding, Entitlement, and Quota
- Service Product, Runtime, and Requirement Foundation (ServiceTypeDefinition, ServiceOffering, ServicePlan, ServiceRequirementSet)
- Sovereignty and Placement Decision v0 (DecisionRecord profiles)
- Plugin Taxonomy and Synthetic Execution Boundary
- AI-Readable Decision and Operation Context
- Slice 0 Integration and Conformance Demo (VS-000)

## Out of Scope

- Real cloud/provider provisioning
- Real PostgreSQL runtime provisioning
- Full OPA/Cedar integration
- Full Keycloak/Vault/Temporal/OpenTelemetry/Kafka integration
- Full multi-provider placement execution
- Production plugin sandbox
- DR, autoscaling, cost, compliance, and AI autonomy execution
- Mandatory ResourcePool or ProviderCapability
- Dual-authority coexistence between alpha and canonical models
- Native AWS, OCI, OpenStack, Kubernetes, or OpenShift objects in customer/core schemas

## Superseded Concepts

The following are not active Phase 2R targets:

- ResourcePool as placement boundary (DEC-0032 superseded)
- ProviderCapability as compatibility boundary (DEC-0033 superseded)
- Generic Provider as combined owner/operator
- Global ServiceClass as canonical catalog concept
- EffectivePolicyContext (replaced by EffectiveGovernanceContext)

## Roadmap Context

Phase 2R is the current execution scope. Future features are available only as scope references in:

```text
docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md
docs/features/FEATURE_INDEX.md
```

Do not implement future-phase roadmap placeholders during Phase 2R. Use them only to avoid architectural dead ends and to preserve adapter boundaries for later reuse.

## Acceptance

Phase 2R is complete when:

1. VS-000 definition of done and conformance matrix pass.
2. Sovereignty and placement simulation demonstrates selected, denied, requires-approval, and indeterminate outcomes.
3. Customer contracts contain no provider-native objects, raw secrets, or protected handles.
4. No active desired-state API supports old and canonical kinds/scopes simultaneously.
5. DecisionRecord and AuditEvent semantics remain owned by FEATURE-0013.
6. Phase 3 readiness review approves replacement of fake execution with one real PostgreSQL path.
