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
- ServiceTypeDefinition, ServiceOffering, and ServicePlan
- GovernanceProfile and EffectiveGovernanceContext
- SovereigntyProfile and EvidenceRecord
- ServiceRuntimeProfile and ServiceRequirementSet
- CloudPlatform/CloudProvider/ExecutionTarget
- PlacementDecision (DecisionRecord profile)
- AuditEvent
- PostgreSQL Management Plane Plugin v0
- Kubernetes/Local Substrate Plugin v0
- PostgreSQL Runtime Plugin v0
- ServiceBinding with SecretRef/CredentialRef
- AI-readable explanation

## Reuse Direction

Actual PostgreSQL runtime should be delegated to mature components such as CloudNativePG, Crunchy Postgres Operator, or a Helm chart for earliest local MVP.

Sovrunn must not build PostgreSQL HA/failover/backup internals in MVP-001.

## Phase 2R Canonical Model Context

**ADH-2026-042 / ARCH-2026.08-PHASE2R-CANONICAL**

The MVP continues to target governed PostgreSQL PaaS placement and provisioning on one substrate. Phase 2R establishes the canonical contracts that the MVP must use:

- **Service catalog:** ServiceTypeDefinition + CloudPlatform-scoped ServiceOffering + versioned ServicePlan (DEC-0049).
- **Governance:** EffectiveGovernanceContext resolved from composed profiles (DEC-0050).
- **Placement:** Qualified ExecutionTarget candidates with DecisionRecord profiles (DEC-0042, DEC-0043).
- **Sovereignty:** Evidence-backed SovereigntyAssessment through registered DecisionProfile (DEC-0055).
- **Binding:** Per-consumer, SecretRef-only ServiceBinding (DEC-0051).
- **Execution:** PluginExecution through governed plugin/adapter boundary (Phase 2R: fake; Phase 3: real).
- **Customer projection:** Safe, immutable ServicePlacement with no topology leakage (DEC-0046).
- **Reference flow:** `docs/architecture/reference-flows/postgresql-end-to-end.md`.

The canonical reference flow for the MVP is defined in the PostgreSQL end-to-end reference flow document. Phase 3 executes this flow with real PostgreSQL on one substrate.
