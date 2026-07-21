---
doc_type: decision_index
title: Sovrunn Decision Index
status: active
phase: 2
ai_load_priority: always
ai_summary: Canonical accepted and proposed decisions that AI must follow during Sovrunn development.
---

# Sovrunn Decision Index

## 1. Purpose

This file is the compact decision index for Sovrunn.

AI agents must check this file before generating architecture, code, tests, or documentation.

## 2. Status Values

| Status | Meaning |
|---|---|
| Proposed | Not yet accepted |
| Accepted | Must be followed |
| Implemented | Accepted and implemented in code or docs |
| Validated | Implemented and verified by feature/phase gate |
| Superseded | Replaced by a newer decision |
| Deprecated | Still historical, but not recommended for new work |
| Rejected | Must not be implemented |

## 3.1 Decision Lifecycle Rule

A decision moves through this lifecycle when applicable:

```text
IDEA -> PROPOSED -> ACCEPTED -> IMPLEMENTED -> VALIDATED -> SUPERSEDED / DEPRECATED
```

Chat discussions and generated prompts do not create accepted decisions. Accepted decisions require explicit approval and index updates.

## 3. Decisions

| ID | Decision | Area | Status | Reference |
|---|---|---|---|---|
| DEC-0001 | SDE SQL hot path must remain in-process and avoid synchronous network calls between runtime stages. | SDE Runtime | Accepted | ADR-0001 |
| DEC-0002 | Protocol plugins translate client protocol behavior into SDE semantic requests; they do not own planning or execution. | SDE Runtime | Accepted | RFC-SDE |
| DEC-0003 | Engine plugins execute against downstream datastores and must not know client protocol details. | SDE Runtime | Accepted | RFC-SDE |
| DEC-0004 | Control Plane must not be called synchronously for every SQL query. | SDE Runtime | Accepted | ADR-0001 |
| DEC-0005 | dStoreOps plugins belong to Datastore Management Plane, not SQL data path. | Management Plane | Accepted | ADR-0001 |
| DEC-0006 | Data Plane should scale horizontally as mostly stateless replicas. | Runtime | Accepted | Constitution |
| DEC-0007 | Go is the primary implementation language; Rust is optional later for proven hot modules. | Engineering | Accepted | ADR-FUTURE |
| DEC-0008 | Data-path plugins start compiled-in; remote/dynamic runtime plugins are deferred. | Plugin Model | Accepted | ADR-0001 |
| DEC-0009 | Platform Core Skeleton should be built before expanding PostgreSQL compatibility. | Roadmap | Accepted | development-phases.md |
| DEC-0010 | Control Plane, Management Planes, and Data Plane may deploy as separate containers. | Deployment | Accepted | platform-core.md |
| DEC-0011 | Foundation services start embedded inside Control Plane but package-isolated. | Platform Core | Accepted | platform-core.md |
| DEC-0012 | PostgreSQL dStoreOps plugin starts compiled into Datastore Management Plane behind a plugin interface. | Datastore Management | Accepted | RFC-0020 |
| DEC-0013 | SDE should evolve as part of a plugin ecosystem for sovereign data runtime dependencies. | SDE | Accepted | vision.md |
| DEC-0014 | Sovrunn is designed from day one as a cloud-native sovereign PaaS platform, not only an SDE product. | Product | Accepted | vision.md |
| DEC-0015 | SDE is a major interoperable data platform capability inside Sovrunn, not the entire platform. | Product | Accepted | vision.md |
| DEC-0016 | Sovrunn’s purpose is to fill sovereign cloud-native capability gaps for local providers, colocation providers, government platforms, and on-prem enterprises. | Product | Accepted | vision.md |
| DEC-0017 | Development priority is Platform Core Skeleton, then Organization Governance, then ServiceOps, then PostgreSQL PaaS, then more services and SDE. | Roadmap | Accepted | development-phases.md |
| DEC-0018 | Sovrunn must support Organization and OrganizationUnit layers above Tenant for large institutional deployments. | Governance | Accepted | RFC-0012 |
| DEC-0019 | Governance, audit, logs, backup, archival, identity, and security baselines must be centrally managed at Organization level. | Governance | Accepted | RFC-0012 |
| DEC-0020 | OrganizationUnit policies may strengthen but must not weaken Organization baseline policies. | Policy | Accepted | RFC-0012 |
| DEC-0021 | Tenant isolation must support namespace, vCluster, and dedicated-cluster profiles. | Multi-Tenancy | Accepted | RFC-0012 |
| DEC-0022 | Sovrunn should use open-source identity, policy, observability, GitOps, and secrets systems while providing unified abstractions above them. | Platform Stack | Accepted | constitution.md |
| DEC-0023 | All meaningful platform changes must create or link to an Operation. | Operations | Accepted | platform-core.md |
| DEC-0024 | AI agents must operate through governed Sovrunn tools and must not bypass policy, tenant boundaries, approvals, or audit. | AI | Accepted | constitution.md |
| DEC-0025 | All Sovrunn design files are platform source files optimized for human review and AI reasoning. | Documentation | Accepted | AI_DOC_AUTHORING_STANDARD.md |
| DEC-0026 | Sovrunn follows reuse-before-build as a core engineering rule. | Architecture | Accepted | DEC-0026 / RFC-0021 |
| DEC-0027 | Phase 2 builds model, decision, audit, policy-context, adapter, and plugin-taxonomy foundation only. | Roadmap | Accepted | DEC-0027 / PHASE2_SCOPE |
| DEC-0028 | Policy logic uses PolicyEngineAdapter; Sovrunn must not embed custom policy rules in handlers or placement. | Policy | Accepted | DEC-0028 / RFC-0025 |
| DEC-0029 | Provider/Substrate Plugin, PaaS Service Management Plane Plugin, and PaaS Service Runtime Plugin are separate plugin planes. | Plugin Model | Accepted | DEC-0029 / RFC-0027 |
| DEC-0030 | MVP-001 is Governed PostgreSQL PaaS Placement and Provisioning on one substrate. | MVP | Accepted | DEC-0030 / RFC-0029 |
| DEC-0031 | ServicePlan remains customer-facing; ServiceRuntimeProfile bridges to infrastructure and runtime requirements. | Service Catalog | Accepted | RFC-0028 |
| DEC-0032 | ResourcePool is the placement boundary. | Placement | Accepted | RFC-0024 |
| DEC-0033 | ProviderCapability is the compatibility boundary. | Placement | Accepted | RFC-0024 |
| DEC-0034 | PlacementDecision is required before provisioning. | Placement | Accepted | RFC-0026 |
| DEC-0035 | Customer-facing APIs must not expose low-level IaaS complexity by default. | API | Accepted | api-resource-standard.md |
| DEC-0036 | Adapter boundaries are required before integrating external engines expected to evolve or be replaced. | Architecture | Accepted | DEC-0036 / adapter-boundary-model.md |

## 4. AI Usage Notes

When implementing features:

- cite relevant DEC IDs in code comments only when useful,
- include DEC IDs in RFC related decisions,
- never contradict an Accepted decision,
- create a new ADR if an implementation requires changing an Accepted decision.
