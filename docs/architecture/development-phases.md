---
doc_type: architecture
title: Sovrunn Development Phases
status: draft
phase: 0
ai_load_priority: always
ai_summary: Splits the Sovrunn end-state into development phases and defines scope boundaries for Phase 1.
---

# Sovrunn Development Phases

## 1. Purpose

This document splits the Sovrunn end-state architecture into development phases.

AI agents must use this file to avoid implementing future-phase features too early.

## 2. Phase Summary

| Phase | Name | Primary Goal |
|---:|---|---|
| 0 | Foundation and AI Development System | Make project AI-developable without architecture drift. |
| 1 | Platform Core Skeleton | Build the core Sovrunn resource grammar. |
| 2 | Organization Governance Layer | Add organization-first governance and inheritance. |
| 3 | ServiceOps Framework | Create repeatable plugin model. |
| 4 | PostgreSQL PaaS MVP | First customer-testable service. |
| 5 | Observability, Audit, Backup, Security | Make the platform design-partner credible. |
| 6 | Multi-Service PaaS Beta | Add 2–3 more service families. |
| 7 | AI Plane MVP | Add read-only and draft-generating AI assistance. |
| 8 | SDE as Managed Service | Bring SDE into Sovrunn as a managed service. |
| 9 | Pilot Readiness | Package for customer lab and paid proof-of-concept. |
| 10 | Production Beta | Controlled production for friendly customers. |
| 11 | Advanced SDE and Transformation | Add SIR, migration, and transformation capabilities. |

## 3. Phase 0: Foundation and AI Development System

### Goal

Create the minimum documentation and standards required before Phase 1 coding.

### Deliverables

- foundation docs,
- decision index,
- glossary,
- architecture docs,
- AI authoring standard,
- AI context guide,
- AI feature factory,
- engineering standards,
- RFC-0012,
- RFC-0020.

### Acceptance

AI can understand:

- what Sovrunn is,
- what Sovrunn is not,
- what to build in Phase 1,
- what not to build,
- canonical terms,
- constitutional rules,
- coding standards,
- testing standards.

## 4. Phase 1: Platform Core Skeleton

### Goal

Build the core Sovrunn platform grammar.

### In Scope

Resources:

- Organization,
- OrganizationUnit,
- Tenant,
- Project,
- ServiceClass,
- ServicePlan,
- ServiceInstance,
- ServiceBinding,
- Plugin,
- Capability,
- Operation.

Capabilities:

- API server skeleton,
- in-memory or simple persistent registry,
- basic validation,
- health/readiness,
- basic CLI/API,
- operation status model,
- plugin/capability registration model.

### Out of Scope

- full UI,
- billing,
- marketplace,
- multi-cluster federation,
- advanced AI,
- SDE transformation,
- PostgreSQL operator integration,
- Keycloak production integration,
- Vault production integration,
- advanced observability.

### Acceptance

- Organization can be created and retrieved.
- OrganizationUnit can be created under Organization.
- Tenant can be created under OrganizationUnit.
- Project can be created under Tenant.
- ServiceClass and ServicePlan can be registered.
- Plugin and Capability can be registered.
- ServiceInstance request can create an Operation.
- Operation status can be tracked.
- All validation failures return explicit reasons.

## 5. Phase 2: Organization Governance Layer

### Goal

Add NIC-like organization governance.

### In Scope

- PolicySet,
- QuotaProfile,
- SecurityProfile,
- BackupProfile,
- ArchivalProfile,
- AuditProfile,
- IdentityProviderRef,
- IsolationProfile,
- Entitlement,
- policy inheritance resolver,
- effective policy calculator,
- audit event model.

### Acceptance

- Organization policy applies to all descendants.
- OrganizationUnit policy may strengthen baseline.
- OrganizationUnit policy cannot weaken baseline.
- Tenant requests are checked against entitlement, quota, and policy.
- Denied requests fail fast with clear reasons.

## 6. Phase 3: ServiceOps Framework

### Goal

Create repeatable plugin model.

### In Scope

- ServiceOps SDK,
- PluginManifest schema,
- CapabilityDeclaration schema,
- Operation handler contract,
- conformance test framework,
- plugin registry,
- capability registry.

### Acceptance

- plugin can declare supported operations,
- plugin can be registered,
- ServiceInstance can resolve to a plugin,
- Operation can call a plugin handler,
- plugin can pass conformance tests.

## 7. Phase 4: PostgreSQL PaaS MVP

### Goal

First customer-testable managed service.

### In Scope

- `datastore.postgresql` ServiceClass,
- PostgreSQL ServicePlans,
- `postgres.dStoreOps` plugin,
- integration with existing PostgreSQL operator,
- ServiceBinding with credentials,
- basic observe/status,
- delete flow,
- basic backup hook.

### Acceptance

- Tenant can request PostgreSQL.
- Sovrunn provisions PostgreSQL through existing operator.
- Credentials are generated.
- ServiceBinding is created.
- Operation status is tracked.
- Audit event is recorded.
- Policy and quota are enforced.

## 8. Phase 5: Observability, Audit, Backup, Security

### Goal

Make platform credible for design partners.

### In Scope

- OpenTelemetry integration,
- Prometheus/Grafana dashboard templates,
- Loki/log profile integration,
- AuditEvent aggregation,
- BackupStatus,
- RestoreRequest,
- SecurityProfile enforcement,
- Kyverno/OPA integration,
- External Secrets/Vault integration.

## 9. Phase 6: Multi-Service PaaS Beta

### Goal

Prove Sovrunn is not only a PostgreSQL wrapper.

### Recommended Services

1. MinIO objectOps.
2. Redis/Valkey or Dragonfly cacheOps.
3. Envoy/Kong/MetalLB gatewayOps or lbOps.

## 10. Phase 7: AI Plane MVP

### Goal

Make Sovrunn AI-first without unsafe automation.

### In Scope

- AI Gateway,
- Tool Registry,
- RAG over docs/DEC/ADR/RFC/runbooks,
- Read-only Platform Assistant,
- Operation Explainer,
- ServiceInstance Draft Generator,
- Plugin Manifest Generator,
- Plan Validator.

Allowed AI levels:

- explain,
- recommend,
- draft.

No direct production mutation.

## 11. Phase 8: SDE as Managed Service

### Goal

Bring SDE into Sovrunn as a managed PaaS service.

### In Scope

- `dataengine.sde` ServiceClass,
- SDE ServicePlan,
- SDE Gateway deployment,
- DataPlane registration,
- binding to PostgreSQL ServiceInstance,
- basic timing metrics.

## 12. Phase 9: Pilot Readiness

### Goal

Package Sovrunn for customer testing and paid proof-of-concept.

### In Scope

- installer,
- Helm charts,
- GitOps install profile,
- upgrade procedure,
- backup/restore runbook,
- tenant onboarding guide,
- admin guide,
- troubleshooting guide,
- demo scripts.

## 13. Phase 10: Production Beta

### Goal

Controlled production for friendly customers.

### In Scope

- HA control plane,
- persistent metadata store,
- multi-cluster baseline,
- OIDC/Keycloak integration,
- production secrets profile,
- policy packs,
- plugin conformance suite,
- upgrade automation,
- tenant isolation profiles,
- support runbooks.

## 14. Phase 11: Advanced SDE and Transformation

### Goal

Differentiate beyond service management.

### In Scope

- SIR canonical model,
- Capability Analyzer,
- TransformationMapping,
- Migration Controller,
- PostgreSQL workload analyzer,
- hybrid routing,
- explicit reject model,
- AI migration advisor.

## 15. Phase 1 Coding Rule

Do not start PostgreSQL provisioning before the platform grammar exists.

Do not start SDE transformation before platform credibility exists.

Do not start marketplace, billing, or advanced UI before pilot readiness.
