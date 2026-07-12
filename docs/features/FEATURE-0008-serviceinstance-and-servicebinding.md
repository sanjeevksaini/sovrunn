---
doc_type: feature
id: FEATURE-0008
title: ServiceInstance and ServiceBinding
status: draft
phase: 1
depends_on: [FEATURE-0004, FEATURE-0006, FEATURE-0007]
ai_load_priority: feature
ai_summary: Implements service consumption resources without real provisioning.
---

# FEATURE-0008 ServiceInstance and ServiceBinding

## 1. Objective

Implement the core service consumption resources.

`ServiceInstance` represents a tenant/project-scoped requested service.

`ServiceBinding` represents connection/use relationship between a consumer and a ServiceInstance.

No real infrastructure provisioning is required in Phase 1.

## 2. ServiceInstance Resource

```json
{
  "apiVersion": "platform.sovrunn.io/v1alpha1",
  "kind": "ServiceInstance",
  "metadata": {
    "name": "nhm-prod-postgres"
  },
  "spec": {
    "organizationRef": "nic",
    "organizationUnitRef": "ministry-health",
    "tenantRef": "national-health-mission",
    "projectRef": "production",
    "serviceClassRef": "datastore.postgresql",
    "servicePlanRef": "postgres-small-ha",
    "parameters": {
      "databaseName": "nhm"
    }
  },
  "status": {
    "phase": "Ready",
    "message": "Registered only; no real provisioning in Phase 1"
  }
}
```

## 3. ServiceBinding Resource

```json
{
  "apiVersion": "platform.sovrunn.io/v1alpha1",
  "kind": "ServiceBinding",
  "metadata": {
    "name": "nhm-app-postgres-binding"
  },
  "spec": {
    "serviceInstanceRef": "nhm-prod-postgres",
    "consumerRef": {
      "kind": "Application",
      "name": "nhm-app"
    },
    "bindingType": "credentials"
  },
  "status": {
    "phase": "Ready",
    "secretRef": "stub-secret-ref"
  }
}
```

## 4. API

| Method | Path |
|---|---|
| POST | `/v1/service-instances` |
| GET | `/v1/service-instances/{name}` |
| GET | `/v1/service-instances?tenantRef=...&projectRef=...` |
| PUT | `/v1/service-instances/{name}` |
| DELETE | `/v1/service-instances/{name}` |
| POST | `/v1/service-bindings` |
| GET | `/v1/service-bindings/{name}` |
| GET | `/v1/service-bindings?serviceInstanceRef=...` |
| DELETE | `/v1/service-bindings/{name}` |

## 5. Validation

- Organization must exist.
- OrganizationUnit, if set, must exist.
- Tenant must exist and match Organization/OU.
- Project must exist and match Tenant.
- ServiceClass must exist.
- ServicePlan must exist and reference ServiceClass.
- At least one active Capability should exist for the ServiceClass, but Phase 1 may warn instead of blocking.
- ServiceInstance delete fails if ServiceBindings exist.

## 6. Phase 1 Behavior

ServiceInstance creation does not provision real infrastructure.

It validates references, stores the resource, creates Operation record, and sets status to `Ready`.

## 7. Acceptance Criteria

- Create ServiceInstance under Project.
- Reject invalid Tenant/Project.
- Reject invalid ServicePlan.
- Create ServiceBinding to ServiceInstance.
- Delete ServiceInstance blocked while binding exists.
- Operations are recorded.
