---
doc_type: feature
id: FEATURE-0002
title: OrganizationUnit Resource
status: draft
phase: 1
depends_on: [FEATURE-0001]
ai_load_priority: feature
ai_summary: Implements delegated OrganizationUnit hierarchy under Organization.
---

# FEATURE-0002 OrganizationUnit Resource

## 1. Objective

Implement `OrganizationUnit` as a delegated governance boundary under an Organization.

## 2. Resource

Example:

```json
{
  "apiVersion": "platform.sovrunn.io/v1alpha1",
  "kind": "OrganizationUnit",
  "metadata": {
    "name": "ministry-health",
    "displayName": "Ministry of Health"
  },
  "spec": {
    "organizationRef": "nic",
    "parentOrganizationUnitRef": null,
    "description": "Health ministry OU",
    "delegatedAdminGroups": ["moh-cloud-admins"]
  },
  "status": {
    "phase": "Active"
  }
}
```

## 3. API

| Method | Path | Behavior |
|---|---|---|
| POST | `/v1/organization-units` | Create OU |
| GET | `/v1/organization-units/{name}` | Read OU |
| GET | `/v1/organization-units?organizationRef=nic` | List OUs |
| PUT | `/v1/organization-units/{name}` | Update mutable fields |
| DELETE | `/v1/organization-units/{name}` | Delete only if empty |

## 4. Validation

- `metadata.name` is required.
- `spec.organizationRef` must reference an existing Organization.
- `parentOrganizationUnitRef`, if set, must reference an OU in the same Organization.
- Name must be unique within the platform for Phase 1.
- Delete fails if OU has child OUs, Tenants, Projects, or ServiceInstances.

## 5. Registry Methods

```go
CreateOrganizationUnit(ctx, ou)
GetOrganizationUnit(ctx, name)
ListOrganizationUnits(ctx, filter)
UpdateOrganizationUnit(ctx, name, patch)
DeleteOrganizationUnit(ctx, name)
```

## 6. Acceptance Criteria

- Cannot create OU for missing Organization.
- Can create OU under Organization.
- Can create nested OU.
- Cannot create OU with parent in different Organization.
- Can list OUs by Organization.
- Delete protected when children exist.
