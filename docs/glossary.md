---
doc_type: glossary
title: Sovrunn Glossary
status: active
phase: 2R
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
| Personal Organization | Automatically created Organization for personal onboarding with default Tenant and default Project (DEC-0047). |
| IsolationProfile | Defines tenant isolation mode: namespace, vCluster, or dedicated cluster. |

## 4. Canonical Cloud Model Terms (Phase 2R)

| Term | Definition |
|---|---|
| CloudPlatform | Product-ownership boundary for a sovereign cloud offering. Owns ServiceOfferings and customer enrollment. |
| CloudProvider | Infrastructure/operations supply boundary. Operates physical infrastructure and participates in CloudPlatform installations. |
| CloudProviderParticipation | Relationship between a CloudProvider and a CloudPlatform installation. Separate from customer enrollment. |
| CloudEnrollment | Relationship that joins a customer Organization to a CloudPlatform (DEC-0038). |
| SovrunnInstallation | One deployed instance of the Sovrunn platform, with exactly one CloudProvider for credential isolation (DEC-0054). |
| HostingLocation | Registered geographic/sovereignty fact. Independent of governance scope (DEC-0041). |
| Datacenter | Physical site boundary. |
| FaultDomain | Failure-isolation boundary inside a datacenter. |
| InfrastructureStack | Externally operated substrate stack contained by exactly one FaultDomain; technology is descriptive and not identity or capability. |
| ExecutionTarget | Qualified actionable realization boundary where services may execute (DEC-0042). No mandatory ResourcePool or ProviderCapability. |

## 5. Service Catalog Terms (Canonical)

| Term | Definition |
|---|---|
| ServiceTypeDefinition | Reusable, implementation-neutral definition of a service type. Shared across CloudPlatforms. Immutable by version (DEC-0044). |
| ServiceOffering | CloudPlatform-scoped product that references a ServiceTypeDefinition (DEC-0049). |
| ServicePlan | Versioned, immutable customer-facing plan under a ServiceOffering. |
| ServicePortfolio | Versioned collection of ServiceOfferings for an industry or domain (DEC-0040). |
| ServiceRuntimeProfile | Runtime/capability requirements bridge between customer-facing plan and infrastructure. Immutable by version. |
| ServiceRequirementSet | Immutable requirements for placement and execution. |
| ServiceInstance | A provisioned instance of a service. Lifecycle-contained by its owner scope. |
| ServiceBinding | Per-consumer, separately revocable, SecretRef-only service-access boundary (DEC-0051). |
| ServicePlacement | Safe, immutable customer projection of where a service runs. No topology or credential leakage (DEC-0046). |
| ServiceDeploymentPlan | Immutable plan referencing ExecutionTargets and cross-target constraints. |
| ServiceRelationshipDefinition | Type definition for cross-service semantics (DEC-0052). |
| ServiceRelationship | Instance of a cross-service dependency with lifecycle awareness (DEC-0052). |
| SecretRef | Reference to a secret stored in an approved secret backend. Secret values must not be exposed in normal API responses. |

## 6. Service Catalog Terms (Migration Input)

| Term | Definition | Migration |
|---|---|---|
| ServiceClass | Former global type concept. | Maps to ServiceTypeDefinition + CloudPlatform-scoped ServiceOffering (DEC-0049). No active canonical authority. |

## 7. Plugin Terms

| Term | Definition |
|---|---|
| Plugin | Implementation unit that performs lifecycle operations for a service family or provider. |
| Capability | Declared ability of a plugin or service. Example: `Backup`, `Restore`, `Scale`. |
| PluginManifest | Document that declares plugin name, kind, version, capabilities, dependencies, and deployment mode. |
| PluginExecution | Bounded execution of a plugin action within the governed contract (Phase 2R: fake/synthetic only). |
| ConformanceTest | Test that validates a plugin implements its declared contract correctly. |
| ProtectedHandle | Internal plugin/adapter handle that must not appear in customer APIs or customer-facing projections. |

## 8. ServiceOps Plugin Families

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

## 9. Operation Terms

| Term | Definition |
|---|---|
| Operation | Auditable asynchronous lifecycle action. Example: provision, backup, restore, delete. |
| OperationStatus | Current state of an Operation. |
| AuditEvent | Normalized audit record for platform, plugin, policy, identity, and service actions. |
| CorrelationID | Identifier used to trace a request across components. |

## 10. Governance Terms (Canonical)

| Term | Definition |
|---|---|
| EffectiveGovernanceContext | Single resolved governance composition from governance, security, data, and cost profiles (DEC-0050). Non-weakenable by default. |
| GovernanceProfile | Set of governance controls applicable to a scope. |
| SecurityProfile | Security baseline applied to resources. |
| SovereigntyProfile | Sovereignty requirements and dimensional evaluation criteria (DEC-0055). |
| SovereigntyFactSet | Registered facts about dependencies, jurisdictions, and control paths for sovereignty evidence. |
| EvidenceRecord | Proof that a control or requirement is satisfied. Fresh/stale/missing/tampered states (DEC-0055). |
| RegulatoryPolicyBundle | Regulatory requirements that may reference governance or sovereignty controls. |
| ProfileAssignment | Assignment of a governance, security, sovereignty, or cost profile to a scope. |
| EntitlementPackage | Collection of service entitlements granted to a scope (DEC-0039). |
| ServiceEntitlement | Permission to consume a specific ServiceOffering or ServicePlan (DEC-0039). |
| QuotaPolicy | Resource limits and quotas, independently evaluated from entitlement (DEC-0039). |
| DataPlacementPolicy | Policy describing where data may be stored, processed, replicated, backed up, or cached. |
| CostGuardrail | Cost or budget boundary considered before placement/scaling/execution. |
| ExceptionGrant | Scoped, time-bound, approved exception to a governance control. |

## 11. Decision and Policy Terms

| Term | Definition |
|---|---|
| DecisionRecord | Common immutable governed-conclusion envelope for Sovrunn decisions. Uses registered DecisionProfile extensions (DEC-0043). |
| DecisionProfile | Typed extension of DecisionRecord: sovereignty, placement, governance, compliance, etc. |
| EvaluationResult | Atomic result within a DecisionRecord: allow, deny, unknown, error. |
| SovereigntyAssessment | DecisionRecord profile for sovereignty evaluation (DEC-0043). |
| PlacementDecision | DecisionRecord profile for placement evaluation (DEC-0043). |
| PolicyEvaluationRequest | Engine-neutral request sent to a policy engine adapter. |
| PolicyEvaluationResult | Engine-neutral result returned from a policy engine adapter. |
| PolicyEngineAdapter | Adapter boundary for OPA, Cedar, or other policy engines. |

## 12. Scope Terms

| Term | Definition |
|---|---|
| ScopeKind | Seven-value governance vocabulary: Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider (DEC-0037). |
| metadata.scopeRef | Sole scope authority on any scoped resource. |

## 13. IAM and Access Terms

| Term | Definition |
|---|---|
| Membership | Organization/tenant membership of a principal. |
| RoleDefinition | Named set of permissions. |
| RoleAssignment | Scoped binding of a principal to a role. |
| PrivilegedAccessRequest | JIT request for elevated access. |
| AccessReview | Periodic review of role assignments. |
| ApprovalPolicy | Conditions under which an action requires approval. |
| ApprovalRequest | Request for human approval before proceeding. |

## 14. Adapter Terms

| Term | Definition |
|---|---|
| Provider/Substrate Plugin | Plugin that executes or validates infrastructure actions. |
| PaaS Service Management Plane Plugin | Plugin that plans service lifecycle and topology. |
| PaaS Service Runtime Plugin | Plugin that configures/checks actual service runtime and bindings. |
| PluginTrustProfile | Trust and certification metadata for plugins. |
| PluginCredentialPolicy | Policy governing plugin credentials and permissions. |
| CredentialRef | Reference to credentials stored outside Sovrunn core records. |
| OperationEngineAdapter | Adapter boundary for operation/workflow engines. |
| ObservabilityAdapter | Adapter boundary for metrics, logs, traces, and health systems. |
| RepositoryAdapter | Persistence boundary for resource stores. |

## 15. Platform Lifecycle Terms

| Term | Definition |
|---|---|
| CanonicalMigrationPlan | Controlling plan for alpha-to-canonical model migration (DEC-0058). |
| CanonicalMigrationRecord | Auditable record of a completed migration step. |
| MaintenanceEpoch | Provider-declared maintenance window with explicit notice and fencing (DEC-0057). |
| MaintenanceFence | Prevention of execution on targets undergoing maintenance until requalified (DEC-0057). |

## 16. SDE Terms

| Term | Definition |
|---|---|
| Protocol Plugin | SDE plugin that understands a client protocol such as PostgreSQL or MySQL. |
| Semantic Request | Protocol-neutral representation of client intent. |
| SIR | SDE Intermediate Representation; semantic intent model, not a PostgreSQL-only AST. |
| Capability Analyzer | SDE component that checks whether target datastore can support a request. |
| TransformationMapping | Mapping specification for supported source-to-target datastore transformations. |
| Engine Plugin | SDE plugin that executes against a target datastore. |
| Hybrid Routing | SDE mode where some requests are routed to one backend and some to another based on capability and mapping. |

## 17. AI Terms

| Term | Definition |
|---|---|
| AI Plane | Sovrunn layer for AI gateway, agents, tools, RAG, validation, approval, and audit. |
| AI Gateway | Model/provider abstraction and policy enforcement point for AI calls. |
| Tool Registry | Registry of governed tools AI agents may call. |
| RAG Knowledge Base | Indexed Sovrunn docs, decisions, RFCs, specs, runbooks, and state used for AI retrieval. |
| Plan Validator | Validator that checks AI-generated plans against schema, policy, quota, security, and capability rules. |

## 18. Retired/Migration-Input Terms

These terms have no active canonical authority and appear only as migration inputs or implementation history:

| Term | Former Definition | Replaced By |
|---|---|---|
| Provider | Former combined product-owner/operator boundary | CloudPlatform + CloudProvider (DEC-0037) |
| ProviderLocation | Former geography boundary | HostingLocation (DEC-0041) |
| ProviderDatacenter | Former physical site boundary | Datacenter |
| DatacenterFailureDomain | Former failure domain | FaultDomain |
| ResourcePool | Former placement boundary | ExecutionTarget (DEC-0042) |
| ProviderCapability | Former compatibility boundary | Target qualification via adapter (DEC-0042) |
| CapabilityStatus | Former capability state | Target availability axes |
| EffectivePolicyContext | Former resolved context term | EffectiveGovernanceContext (DEC-0050) |
| ServiceClass | Former global service type | ServiceTypeDefinition + ServiceOffering (DEC-0049) |
| PlacementCandidate | Former candidate concept | Qualified ExecutionTarget candidates |
| PlacementRequest | Former request concept | ServiceRequirementSet + qualified candidates |
