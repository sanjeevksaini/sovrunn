---
doc_type: resource_spec
title: Phase 1 Resource Model
status: draft
phase: 1
ai_load_priority: always
ai_summary: Consolidated resource model for Sovrunn Phase 1 platform grammar.
---

# Phase 1 Resource Model
> **Phase 1 baseline note:** This document remains valid as the Phase 1 baseline. Phase 2 extends Sovrunn with reuse-first architecture, adapter boundaries, provider-neutral resource modeling, policy evaluation abstraction, decision/audit standards, plugin taxonomy, and governed placement. Do not treat this Phase 1 document as the complete Phase 2 scope.

## Resource Hierarchy

```text
Organization
  -> OrganizationUnit
      -> Tenant
          -> Project
              -> ServiceInstance
                  -> ServiceBinding
```

## Catalog Hierarchy

```text
ServiceClass
  -> ServicePlan
```

## Plugin Hierarchy

```text
Plugin
  -> Capability
```

## Operation Model

```text
Operation
  -> records create/update/delete/lifecycle activity
```

## Common Metadata

All resources should share:

```json
{
  "apiVersion": "platform.sovrunn.io/v1alpha1",
  "kind": "ResourceKind",
  "metadata": {
    "name": "resource-name",
    "displayName": "Optional Display Name",
    "labels": {},
    "annotations": {}
  },
  "spec": {},
  "status": {
    "phase": "Active",
    "message": ""
  }
}
```

## Common Status Phases

For governance/catalog resources:

```text
Active
Inactive
Deleting
Failed
```

For Operations:

```text
Pending
Running
Succeeded
Failed
Cancelled
```

For ServiceInstance/ServiceBinding:

```text
Pending
Ready
Failed
Deleting
```

## Name Rules

- Lowercase DNS label preferred.
- Use hyphen separator.
- No spaces.
- Must be unique per resource kind in Phase 1.

## Phase 1 Storage

Use in-memory maps.

Future storage may use PostgreSQL, Kubernetes CRDs, or another durable backend.
