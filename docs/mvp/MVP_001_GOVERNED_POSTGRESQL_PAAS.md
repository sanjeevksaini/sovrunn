---
doc_type: mvp
title: MVP-001 Governed PostgreSQL PaaS
status: draft
phase: 3
ai_load_priority: always
ai_summary: Defines the first customer-testable Sovrunn MVP: governed PostgreSQL placement and provisioning on one substrate.
---

# MVP-001: Governed PostgreSQL PaaS Placement and Provisioning

## Statement

```text
Governed PostgreSQL PaaS Placement and Provisioning on one substrate,
with explainable decisions, security/data policy, audit events,
plugin-chain execution, and AI-readable explanations.
```

## Customer Promise

A customer requests a PostgreSQL service as a PaaS outcome. Sovrunn validates governance, security, entitlement, data placement, service runtime needs, provider capability, and audit requirements before provisioning.

## In Scope

- Organization/Tenant/Project context
- ServiceClass and ServicePlan
- GovernanceProfile and SecurityProfile
- DataPlacementPolicy
- ServiceRuntimeProfile
- Provider/ResourcePool/ProviderCapability
- PlacementDecision
- AuditEvent
- PostgreSQL Management Plane Plugin v0
- Kubernetes/Local Substrate Plugin v0
- PostgreSQL Runtime Plugin v0
- ServiceBinding with SecretRef/CredentialRef
- AI-readable explanation

## Reuse Direction

Actual PostgreSQL runtime should be delegated to mature components such as CloudNativePG, Crunchy Postgres Operator, or a Helm chart for earliest local MVP.

Sovrunn must not build PostgreSQL HA/failover/backup internals in MVP-001.
