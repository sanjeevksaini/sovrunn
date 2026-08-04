---
doc_type: architecture
title: Organization Governance Architecture
status: draft
phase: 1
ai_load_priority: phase
ai_summary: Defines organization-first governance, policy inheritance, tenant isolation, and centralized control for NIC-like deployments.
---

# Organization Governance Architecture

## 1. Purpose

This document defines the Organization, OrganizationUnit, Tenant, and Project governance model.

It supports large institutional deployments such as a national cloud provider serving many government departments as isolated tenants.

## 2. Canonical Hierarchy

```text
Organization
  -> OrganizationUnit
      -> Tenant
          -> Project
              -> ServiceInstance
```

## 3. Governance Responsibilities

| Layer | Responsibilities |
|---|---|
| Organization | Central baseline for identity, policy, audit, security, logs, backup, archival, retention, catalog entitlement, quota framework, compliance. |
| OrganizationUnit | Delegated administration, stricter policies, OU-level quotas, OU-level approvals, OU-level views. |
| Tenant | Isolated service consumption, tenant users/groups, tenant quotas, service instances, service bindings. |
| Project | Environment/workload separation, project service instances, project-level policy additions. |
| ServiceInstance | Actual managed PaaS service. |

## 4. Organization-Level Capabilities

Organization owns:

- identity federation baseline,
- RBAC baseline,
- security baseline,
- compliance baseline,
- audit aggregation,
- logs aggregation,
- backup policy standards,
- archival policy standards,
- retention policy standards,
- approved service catalog,
- approved plugin catalog,
- quota framework,
- incident visibility,
- risk dashboard.

## 5. OrganizationUnit-Level Capabilities

OrganizationUnit may define:

- delegated admins,
- OU-specific quotas,
- stricter policy sets,
- department-specific approval workflows,
- OU-specific service plan entitlement,
- OU-specific backup retention,
- OU-level audit views,
- OU-level cost/showback.

OrganizationUnit must not weaken Organization baseline.

## 6. Tenant-Level Capabilities

Tenant may define:

- tenant admins,
- project list,
- service instance requests,
- service bindings,
- tenant-level quota allocation,
- tenant-level observability view,
- tenant-level secrets references.

Tenant must operate within Organization and OrganizationUnit boundaries.

## 7. Project-Level Capabilities

Project represents environment or workload grouping.

Examples:

- dev,
- test,
- staging,
- production,
- mission-critical-app.

Project may hold:

- ServiceInstances,
- ServiceBindings,
- local labels,
- environment-level policies that strengthen inherited policies.

## 8. Policy Inheritance

Policy resolution:

```text
effectivePolicy =
  Organization baseline
  + OrganizationUnit policy
  + Tenant policy
  + Project policy
```

Rules:

- lower-level policy may strengthen upper policy,
- lower-level policy must not weaken upper policy,
- conflicts must fail fast,
- every policy decision must be auditable.

## 9. Tenant Isolation Profiles

Supported profiles:

| Profile | Use Case |
|---|---|
| namespace | Low-risk dev/test tenants. |
| vCluster | Medium/high isolation tenants needing Kubernetes API boundary. |
| dedicated-cluster | Highest isolation tenants or critical workloads. |

## 10. Identity Model

Sovrunn should not build an identity provider.

Sovrunn should integrate with identity providers such as Keycloak using OIDC/SAML.

Sovrunn owns:

- IdentityProviderRef,
- group mapping,
- role mapping,
- Organization admin mapping,
- OrganizationUnit admin mapping,
- Tenant admin mapping,
- service account model.

## 11. Audit Model

Sovrunn must normalize audit events from:

- Sovrunn API,
- ServiceOps operations,
- plugin actions,
- policy decisions,
- identity events,
- Kubernetes events,
- GitOps events,
- backup/restore events,
- SDE admin events.

Phase 1 conceptual audit inputs (historical):

| Field | Required | Description |
|---|---:|---|
| eventId | yes | Unique event ID. |
| timestamp | yes | Event time. |
| organizationId | yes | Organization scope. |
| organizationUnitId | no | OU scope. |
| tenantId | no | Tenant scope. |
| projectId | no | Project scope. |
| actor | yes | Requester. |
| action | yes | Action performed. |
| resourceType | yes | Target resource type. |
| resourceId | yes | Target resource ID. |
| decision | yes | Allowed, denied, failed, observed. |
| reason | no | Machine-readable reason. |
| correlationId | yes | Trace correlation ID. |

For Phase 2 and later, this table is input history, not a parallel public
`AuditEvent` schema. The canonical contract is the FEATURE-0012
`ImmutableRecord` envelope specialized by FEATURE-0013: type metadata, common
`metadata`, and one immutable `record` payload. `metadata.scopeRef` is the sole
governance-scope authority and uses exactly the six FEATURE-0012 `ScopeKind`
values. Organization, OrganizationUnit, Tenant, and Project identities shown
above are represented through the canonical scope and typed references rather
than independent top-level scope fields. A ServiceInstance is a typed subject
under its governing Project scope, never a `ScopeKind`. Requirements, design,
tasks, and implementation must not recreate the Phase 1 conceptual columns as
a second scope model.

## 12. Backup and Archival Governance

Organization may define default:

- backup frequency,
- backup retention,
- archival retention,
- legal hold rules,
- restore approval workflow,
- backup failure alerting.

ServiceOps plugins execute actual service-specific backup and restore.

## 13. Failure Modes

| Failure | Expected Behavior |
|---|---|
| OrganizationUnit references missing Organization | Reject with `ParentNotFound`. |
| Tenant references missing OrganizationUnit | Reject with `ParentNotFound`. |
| Lower policy weakens baseline | Reject with `PolicyWeakensBaseline`. |
| Tenant exceeds quota | Reject or fail Operation with `QuotaExceeded`. |
| Tenant requests unauthorized ServicePlan | Reject with `EntitlementDenied`. |
| Isolation profile unsupported | Reject with `IsolationProfileUnsupported`. |

## 14. Acceptance Criteria

- Organization hierarchy can be represented.
- OrganizationUnit cannot exist without Organization.
- Tenant cannot exist without OrganizationUnit.
- Project cannot exist without Tenant.
- Policy inheritance algorithm is defined.
- Policy weakening is rejected.
- Isolation profile is explicit.
- AuditEvent contains hierarchy scope.
- AI and generated code use canonical terms.

## 15. AI Implementation Guidance

- Do not use synonyms such as Department, Agency, Workspace, or Account in code unless mapped to canonical terms.
- Use OrganizationUnit, not Department.
- Use Project, not Workspace or Environment.
- Do not implement full policy engine in Phase 1.
- Define clean interfaces for future Kyverno/OPA integration.
- Add negative tests for missing parents and policy weakening.

## Phase 2 Governance Extension

Phase 2R shifts from generic `PolicySet` language toward explicit governance, sovereignty, and policy-context objects under the canonical model (ADH-2026-042):

- GovernanceProfile,
- SovereigntyProfile,
- DataPlacementPolicy,
- CostGuardrail,
- ProfileAssignment,
- EffectiveGovernanceContext (DEC-0050),
- PolicyEvaluationRequest,
- PolicyEvaluationResult,
- DecisionRecord with sovereignty profile (DEC-0043),
- DecisionRecord with placement profile (DEC-0043),
- ServicePlacement (safe customer projection, DEC-0046).

Policy evaluation must go through PolicyEngineAdapter so OPA/Cedar or other engines can be reused later.

## Phase 2R Governance Update

**ADH-2026-042 / DEC-0050**

Phase 2R governance uses:

- **Seven canonical scope kinds:** Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider (DEC-0037).
- **EffectiveGovernanceContext:** single resolved governance composition from profiles, with non-weakenable controls and deterministic conflict resolution (DEC-0050).
- **Sovereignty:** separately assessed through evidence-backed DecisionRecord profiles, not conflated with governance composition (DEC-0055).
- **Personal Organization:** automatic onboarding with default Tenant and default Project (DEC-0047).
- **EntitlementPackage and QuotaPolicy:** independent evaluation (DEC-0039).

The canonical governance model is in `docs/architecture/canonical/sovrunn-finalized-data-model.md`.
