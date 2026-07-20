---
doc_type: api_contract
title: Phase 1 REST API Contract
status: draft
phase: 1
ai_load_priority: phase
ai_summary: Defines initial REST endpoints for Phase 1 resources.
---

# Phase 1 REST API Contract
> **Phase 1 baseline note:** This document remains valid as the Phase 1 baseline. Phase 2 extends Sovrunn with reuse-first architecture, adapter boundaries, provider-neutral resource modeling, policy evaluation abstraction, decision/audit standards, plugin taxonomy, and governed placement. Do not treat this Phase 1 document as the complete Phase 2 scope.

## Base URL

```text
http://127.0.0.1:8080
```

## Health

```text
GET /healthz
GET /readyz
GET /version
```

## Organization Management

```text
POST   /v1/organizations
GET    /v1/organizations
GET    /v1/organizations/{name}
PUT    /v1/organizations/{name}
DELETE /v1/organizations/{name}

POST   /v1/organization-units
GET    /v1/organization-units
GET    /v1/organization-units/{name}
PUT    /v1/organization-units/{name}
DELETE /v1/organization-units/{name}

POST   /v1/tenants
GET    /v1/tenants
GET    /v1/tenants/{name}
PUT    /v1/tenants/{name}
DELETE /v1/tenants/{name}

POST   /v1/projects
GET    /v1/projects
GET    /v1/projects/{name}
PUT    /v1/projects/{name}
DELETE /v1/projects/{name}
```

## Operations

```text
GET /v1/operations
GET /v1/operations/{name}
```

## Catalog

```text
POST   /v1/service-classes
GET    /v1/service-classes
GET    /v1/service-classes/{name}
PUT    /v1/service-classes/{name}
DELETE /v1/service-classes/{name}

POST   /v1/service-plans
GET    /v1/service-plans
GET    /v1/service-plans/{name}
PUT    /v1/service-plans/{name}
DELETE /v1/service-plans/{name}
```

## Plugin Registry

```text
POST   /v1/plugins
GET    /v1/plugins
GET    /v1/plugins/{name}
PUT    /v1/plugins/{name}
DELETE /v1/plugins/{name}

POST   /v1/capabilities
GET    /v1/capabilities
GET    /v1/capabilities/{name}
DELETE /v1/capabilities/{name}
```

## Service Consumption

```text
POST   /v1/service-instances
GET    /v1/service-instances
GET    /v1/service-instances/{name}
PUT    /v1/service-instances/{name}
DELETE /v1/service-instances/{name}

POST   /v1/service-bindings
GET    /v1/service-bindings
GET    /v1/service-bindings/{name}
DELETE /v1/service-bindings/{name}
```

## Error Shape

```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "organizationRef does not exist",
    "details": {}
  }
}
```

## HTTP Status Rules

| Status | Meaning |
|---|---|
| 200 | Successful get/list/update |
| 201 | Successful create |
| 202 | Accepted async operation, future |
| 400 | Invalid request |
| 404 | Resource not found |
| 409 | Conflict, duplicate, or delete blocked |
| 500 | Internal error |
