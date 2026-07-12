---
doc_type: feature
id: FEATURE-0003
title: Tenant Resource
status: draft
phase: 1
depends_on: [FEATURE-0001, FEATURE-0002]
ai_load_priority: feature
ai_summary: Implements Tenant as the primary isolated consumption boundary.
---

# FEATURE-0003 Tenant Resource

## 1. Objective

Implement `Tenant` as the primary security, isolation, and consumption boundary in Sovrunn.

## 2. Resource

```json
{
  "apiVersion": "platform.sovrunn.io/v1alpha1",
  "kind": "Tenant",
  "metadata": {
    "name": "national-health-mission",
    "displayName": "National Health Mission"
  },
  "spec": {
    "organizationRef": "nic",
    "organizationUnitRef": "ministry-health",
    "isolationProfile": "namespace",
    "sovereignLocationRefs": ["in-delhi-1"],
    "adminGroups": ["nhm-admins"]
  },
  "status": {
    "phase": "Active"
  }
}
```

## 3. API

| Method | Path | Behavior |
|---|---|---|
| POST | `/v1/tenants` | Create tenant |
| GET | `/v1/tenants/{name}` | Read tenant |
| GET | `/v1/tenants?organizationRef=nic&organizationUnitRef=ministry-health` | List tenants |
| PUT | `/v1/tenants/{name}` | Update mutable fields |
| DELETE | `/v1/tenants/{name}` | Delete only if empty |

## 4. Validation

- `organizationRef` must exist.
- `organizationUnitRef`, if set, must exist and belong to same Organization.
- `isolationProfile` allowed values in Phase 1: `namespace`, `vcluster`, `dedicated-cluster`.
- Default isolation profile: `namespace`.
- Delete fails if Projects or ServiceInstances exist.

## 5. Acceptance Criteria

- Create tenant under Organization.
- Create tenant under OrganizationUnit.
- Reject missing Organization.
- Reject OU from another Organization.
- List tenants by Organization/OU.
- Delete protected when Projects exist.
