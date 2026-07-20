---
doc_type: phase_scope
title: Phase 2 Scope
status: draft
phase: 2
ai_load_priority: always
ai_summary: Defines Phase 2 as reuse-first model, decision, audit, adapter, and placement simulation foundation. No real provisioning.
---

# Phase 2 Scope

## Purpose

Phase 2 establishes the reusable architecture spine for Sovrunn's provider-neutral sovereign PaaS fabric.

Phase 2 builds:

```text
models
standards
adapter boundaries
policy evaluation abstraction
decision objects
audit events
placement simulation
AI-readable explanations
```

Phase 2 does not build real provider provisioning, real database provisioning, real autoscaling, real failover, or real autonomous operations.

## In Scope

- Reuse Assessment Standard
- API, resource, status, validation, and boundary classification standard
- Decision Object and AuditEvent Standard
- Provider-neutral resource model
- ResourcePool and ProviderCapability model
- Adapter Boundary Foundation
- Policy Evaluation Abstraction
- GovernanceProfile and SecurityProfile
- DataPlacementPolicy and CostGuardrail minimal foundation
- ProfileAssignment and EffectivePolicyContext
- Minimal ServiceEntitlement and Quota placeholder
- ServiceRuntimeProfile
- PlacementRequest and PlacementDecision v0
- Plugin Taxonomy Foundation
- AI-Readable Decision Context
- Integration simulation

## Out of Scope

- Real cloud/provider provisioning
- Real PostgreSQL runtime provisioning
- Full OPA/Cedar integration
- Full Keycloak/Vault/Temporal/OpenTelemetry/Kafka integration
- Full multi-provider placement execution
- Production plugin sandbox
- DR, autoscaling, cost, compliance, and AI autonomy execution


## Roadmap Context

Phase 2 is the current execution scope. Future features are available only as scope references in:

```text
docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md
docs/features/FEATURE_INDEX.md
```

Do not implement future-phase roadmap placeholders during Phase 2. Use them only to avoid architectural dead ends and to preserve adapter boundaries for later reuse.

## Acceptance

Phase 2 is complete when Sovrunn can simulate a governed PostgreSQL placement request, explain why it is allowed or denied, and record a decision/audit event using provider-neutral resource and policy context.
