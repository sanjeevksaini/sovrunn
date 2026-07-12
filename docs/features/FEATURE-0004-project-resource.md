---
doc_type: feature
id: FEATURE-0004
title: Project Resource
status: draft
phase: 1
depends_on: [FEATURE-0003]
ai_load_priority: feature
ai_summary: Implements Project as a workload or environment grouping inside a Tenant.
---

# FEATURE-0004 Project Resource

## 1. Objective

Implement `Project` as a workload or environment grouping inside a Tenant.

`Tenant` is the primary isolation boundary. `Project` groups workloads such as dev, staging, production, analytics, and reporting.

## 2. Resource

```json
{
  "apiVersion": "platform.sovrunn.io/v1alpha1",
  "kind": "Project",
  "metadata": {
    "name": "production",
    "displayName": "Production"
  },
  "spec": {
    "organizationRef": "nic",
    "organizationUnitRef": "ministry-health",
    "tenantRef": "national-health-mission",
    "environmentType": "production",
    "description": "Production project for NHM"
  },
  "status": {
    "phase": "Active"
  }
}
```

## 3. API

| Method | Path | Behavior |
|---|---|---|
| POST | `/v1/projects` | Create project |
| GET | `/v1/projects/{name}` | Read project |
| GET | `/v1/projects?tenantRef=national-health-mission` | List projects |
| PUT | `/v1/projects/{name}` | Update mutable fields |
| DELETE | `/v1/projects/{name}` | Delete only if empty |

## 4. Validation

- `tenantRef` must exist.
- `organizationRef` must match tenant organization.
- `organizationUnitRef`, if set, must match tenant OU.
- `environmentType` allowed values: `dev`, `test`, `staging`, `production`, `analytics`, `reporting`, `default`.
- Delete fails if ServiceInstances or ServiceBindings exist.

## 5. Default Project Rule

Phase 1 may auto-create `default` Project when a Tenant is created, but this is optional.

Recommended: do not auto-create yet; keep behavior explicit for API clarity.

## 6. Acceptance Criteria

- Create project under Tenant.
- Reject project for missing Tenant.
- Reject mismatched Organization.
- List projects by Tenant.
- Delete protected when ServiceInstances exist.
