---
doc_type: glossary
title: Sovrunn Glossary
status: draft
phase: 0
ai_load_priority: always
ai_summary: Canonical terminology for Sovrunn. AI must use these terms consistently and avoid synonyms unless explicitly mapped.
---

# Sovrunn Glossary

## 1. Purpose

This glossary defines canonical Sovrunn terms.

AI agents must use these terms consistently.

Rule:

```text
One concept = one canonical term.
```

## 2. Core Product Terms

| Term | Definition |
|---|---|
| Sovrunn | The complete AI-first, open-standard, cloud-native sovereign PaaS platform. |
| Sovrunn Data Engine | Interoperable data platform capability inside Sovrunn. |
| SDE | Abbreviation for Sovrunn Data Engine. |
| Cloud Management Plane | Sovrunn layer that owns API server, registries, catalog, operations, policy integration, and service binding. |
| Organization Management Layer | Sovrunn layer that owns Organization, OrganizationUnit, Tenant, Project, governance, policy inheritance, audit, backup, archival, identity, and security baselines. |
| Service Management Plane | Domain-specific management layer for a family of services, such as datastore, cache, object storage, gateway, or SDE. |
| ServiceOps | Generic plugin and lifecycle framework for managing PaaS services. |

## 3. Organization Terms

| Term | Definition |
|---|---|
| Organization | Top-level administrative and governance boundary. Example: NIC. |
| OrganizationUnit | Delegated governance boundary under an Organization. Example: Ministry of Health. |
| Tenant | Isolated service consumption boundary under an OrganizationUnit. |
| Project | Environment or workload grouping under a Tenant. Example: dev, test, staging, production. |
| IsolationProfile | Defines tenant isolation mode: namespace, vCluster, or dedicated cluster. |

## 4. Service Catalog Terms

| Term | Definition |
|---|---|
| ServiceClass | A type of managed service offered by Sovrunn. Example: `datastore.postgresql`. |
| ServicePlan | A predefined plan for a ServiceClass. Example: `small`, `medium`, `ha`. |
| ServiceInstance | A provisioned instance of a ServiceClass using a ServicePlan. |
| ServiceBinding | A binding that exposes connection details or credentials to a consumer. |
| SecretRef | Reference to a secret stored in an approved secret backend. Secret values must not be exposed in normal API responses. |

## 5. Plugin Terms

| Term | Definition |
|---|---|
| Plugin | Implementation unit that performs lifecycle operations for a service family or provider. |
| Capability | Declared ability of a plugin or service. Example: `Backup`, `Restore`, `Scale`. |
| PluginManifest | Document that declares plugin name, kind, version, capabilities, dependencies, and deployment mode. |
| ConformanceTest | Test that validates a plugin implements its declared contract correctly. |

## 6. ServiceOps Plugin Families

| Term | Definition |
|---|---|
| dStoreOps | Datastore operations plugin family. |
| cacheOps | Cache operations plugin family. |
| objectOps | Object storage operations plugin family. |
| streamOps | Streaming/messaging operations plugin family. |
| gatewayOps | API gateway operations plugin family. |
| lbOps | Load balancer operations plugin family. |
| faasOps | FaaS/serverless operations plugin family. |
| bigDataOps | Big data processing operations plugin family. |
| sdeOps | SDE service operations plugin family. |

## 7. Operation Terms

| Term | Definition |
|---|---|
| Operation | Auditable asynchronous lifecycle action. Example: provision, backup, restore, delete. |
| OperationStatus | Current state of an Operation. |
| AuditEvent | Normalized audit record for platform, plugin, policy, identity, and service actions. |
| CorrelationID | Identifier used to trace a request across components. |

## 8. Governance Terms

| Term | Definition |
|---|---|
| PolicySet | Collection of policy rules applied at Organization, OrganizationUnit, Tenant, Project, or ServiceInstance level. |
| EffectivePolicy | Resolved policy after applying inheritance. |
| QuotaProfile | Resource limits and quotas for an OrganizationUnit, Tenant, or Project. |
| SecurityProfile | Security baseline applied to resources. |
| BackupProfile | Backup requirements and schedules. |
| ArchivalProfile | Archival and retention requirements. |
| AuditProfile | Audit collection and retention requirements. |
| Entitlement | Permission to consume a ServiceClass, ServicePlan, plugin, or capability. |

## 9. SDE Terms

| Term | Definition |
|---|---|
| Protocol Plugin | SDE plugin that understands a client protocol such as PostgreSQL or MySQL. |
| Semantic Request | Protocol-neutral representation of client intent. |
| SIR | SDE Intermediate Representation; semantic intent model, not a PostgreSQL-only AST. |
| Capability Analyzer | SDE component that checks whether target datastore can support a request. |
| TransformationMapping | Mapping specification for supported source-to-target datastore transformations. |
| Engine Plugin | SDE plugin that executes against a target datastore. |
| Hybrid Routing | SDE mode where some requests are routed to one backend and some to another based on capability and mapping. |

## 10. AI Terms

| Term | Definition |
|---|---|
| AI Plane | Sovrunn layer for AI gateway, agents, tools, RAG, validation, approval, and audit. |
| AI Gateway | Model/provider abstraction and policy enforcement point for AI calls. |
| Tool Registry | Registry of governed tools AI agents may call. |
| RAG Knowledge Base | Indexed Sovrunn docs, decisions, RFCs, specs, runbooks, and state used for AI retrieval. |
| Plan Validator | Validator that checks AI-generated plans against schema, policy, quota, security, and capability rules. |
