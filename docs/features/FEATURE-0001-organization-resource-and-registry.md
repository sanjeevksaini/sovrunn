---
doc_type: feature
id: FEATURE-0001
title: Organization Resource and Registry
status: draft
phase: 1
depends_on: []
ai_load_priority: feature
ai_summary: Implements the top-level Organization resource and in-memory registry.
---

# FEATURE-0001 Organization Resource and Registry
> **Phase 1 baseline note:** This document remains valid as the Phase 1 baseline. Phase 2 extends Sovrunn with reuse-first architecture, adapter boundaries, provider-neutral resource modeling, policy evaluation abstraction, decision/audit standards, plugin taxonomy, and governed placement. Do not treat this Phase 1 document as the complete Phase 2 scope.

## 1. Objective

Implement `Organization` as the top-level governance and ownership resource in Sovrunn.

## 2. Resource

`Organization` represents an enterprise, government body, cloud provider, local provider, or large customer operating Sovrunn.

Example:

```json
{
  "apiVersion": "platform.sovrunn.io/v1alpha1",
  "kind": "Organization",
  "metadata": {
    "name": "nic",
    "displayName": "National Informatics Centre",
    "labels": {
      "country": "in"
    }
  },
  "spec": {
    "description": "Central government cloud organization",
    "sovereignLocations": ["in-delhi-1", "in-mumbai-1"],
    "defaultPolicyProfile": "org-baseline"
  },
  "status": {
    "phase": "Active"
  }
}
```

## 3. API

| Method | Path | Behavior |
|---|---|---|
| POST | `/v1/organizations` | Create organization |
| GET | `/v1/organizations/{name}` | Read organization |
| GET | `/v1/organizations` | List organizations |
| PUT | `/v1/organizations/{name}` | Update mutable fields |
| DELETE | `/v1/organizations/{name}` | Delete only if empty |

## 4. Validation

- `metadata.name` is required.
- Name must be DNS-label compatible.
- Name must be globally unique.
- `metadata.displayName` is optional.
- `spec.sovereignLocations` is optional in Phase 1 but should accept a list.
- Delete must fail if Organization has OrganizationUnits, Tenants, Projects, or ServiceInstances.

## 5. Registry

Create an in-memory registry with:

```go
CreateOrganization(ctx, org)
GetOrganization(ctx, name)
ListOrganizations(ctx)
UpdateOrganization(ctx, name, patch)
DeleteOrganization(ctx, name)
```

## 6. Operation Behavior

Creation may immediately return success in Phase 1, but must create an Operation record once FEATURE-0005 exists.

Before FEATURE-0005, return direct response.

After FEATURE-0005, create:

```text
Operation type: CreateOrganization
target: Organization/nic
phase: Succeeded
```

## 7. Acceptance Criteria

- Can create an Organization.
- Duplicate create fails.
- Invalid name fails.
- Can get/list organizations.
- Can update description/display name/labels.
- Delete non-existing org returns 404.
- Delete org with children returns conflict.
- Unit tests cover create/get/list/update/delete.
