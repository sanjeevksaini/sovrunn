---
doc_type: rfc
id: RFC-0012
title: Organization, Tenant, and Governance Model
status: draft
phase: 1
depends_on:
  - DEC-0018
  - DEC-0019
  - DEC-0020
  - DEC-0021
  - DEC-0022
  - constitution.md
ai_load_priority: feature
ai_summary: Defines Organization, OrganizationUnit, Tenant, Project, policy inheritance, isolation profiles, and audit scope.
---

# RFC-0012: Organization, Tenant, and Governance Model

## 1. Status

Draft for founder review.

## 2. Purpose

Define the organization-first governance model for Sovrunn.

This RFC establishes:

- Organization,
- OrganizationUnit,
- Tenant,
- Project,
- governance inheritance,
- tenant isolation profiles,
- audit scope,
- identity mapping placeholders.

## 3. Goals

- Support NIC-like large organization deployments.
- Provide hierarchy above Tenant.
- Enable centralized governance.
- Enable delegated OrganizationUnit administration.
- Enable isolated Tenant consumption.
- Enable Project-level environment/workload grouping.
- Preserve future identity, policy, audit, backup, archival, and security integrations.

## 4. Non-Goals

This RFC does not implement:

- full policy engine,
- Keycloak integration,
- Kubernetes RBAC integration,
- vCluster creation,
- dedicated cluster provisioning,
- billing,
- full UI,
- customer-facing portal.

## 5. Definitions

| Term | Definition |
|---|---|
| Organization | Top-level administrative and governance boundary. |
| OrganizationUnit | Delegated governance boundary under Organization. |
| Tenant | Isolated service consumption boundary. |
| Project | Environment/workload grouping inside Tenant. |
| IsolationProfile | Defines isolation mode: namespace, vCluster, dedicated cluster. |
| PolicySet | Collection of policy rules. |
| EffectivePolicy | Resolved inherited policy. |

## 6. Resource Model

### 6.1 Organization

Purpose:

- owns central governance baseline,
- owns identity baseline,
- owns policy baseline,
- owns audit and compliance baseline.

Minimum fields:

| Field | Required | Description |
|---|---:|---|
| name | yes | Unique Organization name. |
| displayName | no | Human-readable name. |
| defaultPolicySetRefs | no | Default policy sets. |
| defaultAuditProfileRef | no | Default audit profile. |
| defaultIdentityProviderRef | no | Default identity provider reference. |
| status.phase | yes | Pending, Ready, Failed. |

### 6.2 OrganizationUnit

Purpose:

- delegated governance boundary,
- maps to department/ministry/business unit,
- may strengthen Organization policies.

Minimum fields:

| Field | Required | Description |
|---|---:|---|
| name | yes | Unique name within Organization. |
| organizationRef | yes | Parent Organization. |
| displayName | no | Human-readable name. |
| delegatedAdminGroups | no | Admin group references. |
| policySetRefs | no | OU-specific policies. |
| quotaProfileRef | no | OU quota profile. |
| status.phase | yes | Pending, Ready, Failed. |

### 6.3 Tenant

Purpose:

- isolated service consumption boundary.

Minimum fields:

| Field | Required | Description |
|---|---:|---|
| name | yes | Unique name within OrganizationUnit. |
| organizationRef | yes | Parent Organization. |
| organizationUnitRef | yes | Parent OrganizationUnit. |
| isolationProfileRef | yes | Tenant isolation profile. |
| quotaProfileRef | no | Tenant quota. |
| status.phase | yes | Pending, Ready, Failed. |

### 6.4 Project

Purpose:

- environment/workload grouping inside Tenant.

Minimum fields:

| Field | Required | Description |
|---|---:|---|
| name | yes | Unique name within Tenant. |
| organizationRef | yes | Organization scope. |
| organizationUnitRef | yes | OrganizationUnit scope. |
| tenantRef | yes | Parent Tenant. |
| environment | no | dev, test, staging, production, or custom. |
| status.phase | yes | Pending, Ready, Failed. |

## 7. Policy Inheritance

Policy resolution:

```text
Organization baseline
  -> OrganizationUnit policy
      -> Tenant policy
          -> Project policy
```

Rules:

- Organization policy applies to all descendants.
- OrganizationUnit policy may strengthen Organization policy.
- OrganizationUnit policy must not weaken Organization policy.
- Tenant policy may strengthen inherited policy.
- Project policy may strengthen inherited policy.
- Weakening attempts must fail with `PolicyWeakensBaseline`.

## 8. Isolation Profiles

Supported isolation modes:

| Mode | Description |
|---|---|
| namespace | Shared cluster, namespace isolation. |
| vCluster | Virtual Kubernetes control plane inside shared host cluster. |
| dedicated-cluster | Dedicated cluster for high isolation. |

Phase 1 only models the profile.

Actual Kubernetes/vCluster provisioning is future work.

## 9. API Behavior

Minimum API behavior:

- create/get/list Organization,
- create/get/list OrganizationUnit,
- create/get/list Tenant,
- create/get/list Project.

Validation must occur before persistence.

## 10. Validation Rules

- Organization name must be unique.
- OrganizationUnit must reference existing Organization.
- Tenant must reference existing Organization and OrganizationUnit.
- Tenant OrganizationUnit must belong to the referenced Organization.
- Project must reference existing Tenant.
- Project Tenant must belong to referenced OrganizationUnit and Organization.
- IsolationProfile must be known or explicitly accepted as placeholder.
- Policy weakening must be rejected when policy model is available.

## 11. Failure Modes

| Failure | Reason Code |
|---|---|
| Parent Organization missing | ParentNotFound |
| Parent OrganizationUnit missing | ParentNotFound |
| Tenant parent mismatch | InvalidHierarchy |
| Project parent mismatch | InvalidHierarchy |
| Duplicate resource | AlreadyExists |
| Unsupported isolation profile | IsolationProfileUnsupported |
| Policy weakens baseline | PolicyWeakensBaseline |

## 12. Security and Governance

- Do not expose secret values.
- Include actor/requester where available.
- Include hierarchy scope in audit events.
- Preserve identity provider references but do not implement IdP integration in Phase 1.

## 13. Observability and Audit

Every create/update/delete should produce or be able to produce an AuditEvent with:

- Organization,
- OrganizationUnit where applicable,
- Tenant where applicable,
- Project where applicable,
- actor,
- action,
- resource,
- decision,
- reason,
- correlation ID.

## 14. Tests

Required tests:

- create Organization success,
- duplicate Organization failure,
- create OrganizationUnit success,
- OrganizationUnit missing parent failure,
- create Tenant success,
- Tenant missing parent failure,
- Tenant hierarchy mismatch failure,
- create Project success,
- Project missing parent failure,
- invalid isolation profile failure.

## 15. Acceptance Criteria

- Canonical hierarchy is implemented in models.
- Parent validation works.
- Failure reasons are explicit.
- Audit scope can include hierarchy fields.
- AI uses OrganizationUnit, not Department.
- Phase 1 does not implement full external policy engine.

## 16. Related Decisions

- DEC-0018
- DEC-0019
- DEC-0020
- DEC-0021
- DEC-0022

## 17. AI Implementation Guidance

- Implement only hierarchy resources in Phase 1.
- Do not implement Keycloak.
- Do not implement vCluster provisioning.
- Do not implement full Kyverno/OPA integration.
- Do not use synonyms in model names.
- Add unit tests for validation.
- Add negative tests for hierarchy mismatch.
