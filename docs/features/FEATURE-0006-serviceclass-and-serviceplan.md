---
doc_type: feature
id: FEATURE-0006
title: ServiceClass and ServicePlan
status: draft
phase: 1
depends_on: [FEATURE-0005]
ai_load_priority: feature
ai_summary: Implements service catalog primitives: ServiceClass and ServicePlan.
---

# FEATURE-0006 ServiceClass and ServicePlan

## 1. Objective

Implement the basic service catalog.

`ServiceClass` defines what type of service Sovrunn can offer.

`ServicePlan` defines an approved configuration/capacity/policy shape for that service.

## 2. ServiceClass Resource

```json
{
  "apiVersion": "platform.sovrunn.io/v1alpha1",
  "kind": "ServiceClass",
  "metadata": {
    "name": "datastore.postgresql",
    "displayName": "PostgreSQL"
  },
  "spec": {
    "category": "datastore",
    "description": "Managed PostgreSQL datastore",
    "requiredCapabilities": ["Provision", "Bind", "Observe", "Backup", "Restore", "Delete"]
  },
  "status": {
    "phase": "Active"
  }
}
```

## 3. ServicePlan Resource

```json
{
  "apiVersion": "platform.sovrunn.io/v1alpha1",
  "kind": "ServicePlan",
  "metadata": {
    "name": "postgres-small-ha"
  },
  "spec": {
    "serviceClassRef": "datastore.postgresql",
    "description": "Small HA PostgreSQL plan",
    "tier": "small",
    "highAvailability": true,
    "capacity": {
      "cpu": "2",
      "memory": "4Gi",
      "storage": "100Gi"
    },
    "backupProfileRef": "standard-30d"
  },
  "status": {
    "phase": "Active"
  }
}
```

## 4. API

| Method | Path |
|---|---|
| POST | `/v1/service-classes` |
| GET | `/v1/service-classes/{name}` |
| GET | `/v1/service-classes` |
| PUT | `/v1/service-classes/{name}` |
| DELETE | `/v1/service-classes/{name}` |
| POST | `/v1/service-plans` |
| GET | `/v1/service-plans/{name}` |
| GET | `/v1/service-plans?serviceClassRef=datastore.postgresql` |
| PUT | `/v1/service-plans/{name}` |
| DELETE | `/v1/service-plans/{name}` |

## 5. Validation

- ServiceClass name is unique.
- ServiceClass category allowed values: `datastore`, `cache`, `object-storage`, `stream`, `gateway`, `load-balancer`, `faas`, `big-data`, `sde`, `other`.
- ServicePlan must reference existing ServiceClass.
- ServicePlan cannot be deleted if ServiceInstances use it.

## 6. Acceptance Criteria

- Can register PostgreSQL ServiceClass.
- Can register postgres-small-ha ServicePlan.
- Cannot create plan for missing ServiceClass.
- Cannot delete ServiceClass with plans.
- Operations are recorded.
