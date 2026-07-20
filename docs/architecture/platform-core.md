---
doc_type: architecture
title: Sovrunn Platform Core
status: draft
phase: 1
ai_load_priority: phase
ai_summary: Defines the Phase 1 Platform Core Skeleton resources, boundaries, and acceptance criteria.
---

# Sovrunn Platform Core

## 1. Purpose

Platform Core is the minimum Sovrunn runtime grammar required before service provisioning.

It defines the resources and flows that all future management planes, ServiceOps plugins, AI agents, and SDE management will use.

## 2. Goals

- Define core resources.
- Provide API skeleton.
- Provide registry/storage abstraction.
- Provide validation.
- Provide Operation model.
- Provide Plugin and Capability registration.
- Provide health/readiness.
- Enable future PostgreSQL PaaS implementation.

## 3. Non-Goals

Phase 1 must not implement:

- PostgreSQL operator integration,
- production identity integration,
- production policy engine integration,
- full UI,
- billing,
- marketplace,
- advanced AI,
- SDE transformation,
- multi-cluster federation,
- production HA control plane.

## 4. Core Resources

| Resource | Purpose |
|---|---|
| Organization | Top-level administrative and governance boundary. |
| OrganizationUnit | Delegated governance boundary under Organization. |
| Tenant | Isolated service consumption boundary. |
| Project | Environment/workload grouping under Tenant. |
| ServiceClass | Type of managed service offered by Sovrunn. |
| ServicePlan | Plan/SKU for a ServiceClass. |
| ServiceInstance | Requested or provisioned service. |
| ServiceBinding | Connection/credential binding for a ServiceInstance. |
| Plugin | Registered implementation provider. |
| Capability | Declared ability of a plugin or service. |
| Operation | Auditable asynchronous lifecycle action. |

## 5. Phase 1 Resource Relationships

```text
Organization
  -> OrganizationUnit
      -> Tenant
          -> Project
              -> ServiceInstance
                  -> ServiceBinding

ServiceClass
  -> ServicePlan

Plugin
  -> Capability

ServiceInstance
  -> Operation
```

## 6. API Behavior

Minimum API groups:

| API Group | Example Operations |
|---|---|
| Organization API | create/get/list Organization |
| OrganizationUnit API | create/get/list OrganizationUnit |
| Tenant API | create/get/list Tenant |
| Project API | create/get/list Project |
| Catalog API | create/list ServiceClass and ServicePlan |
| Plugin API | register/list Plugin and Capability |
| Service API | create/list ServiceInstance and ServiceBinding |
| Operation API | create/get/list/update Operation status |
| Health API | healthz, readyz |

## 7. Operation Model

Every meaningful platform change must create or link to an Operation.

Operation phases:

| Phase | Meaning |
|---|---|
| Pending | Operation accepted but not started. |
| Running | Operation is executing. |
| Succeeded | Operation completed successfully. |
| Failed | Operation failed with explicit reason. |
| Cancelled | Operation was cancelled. |

Minimum Operation fields:

| Field | Required | Description |
|---|---:|---|
| id | yes | Unique Operation ID. |
| action | yes | Requested action. |
| resourceRef | yes | Target resource. |
| requestedBy | yes | Actor or service. |
| scope | yes | Organization/OU/Tenant/Project scope. |
| phase | yes | Current phase. |
| reason | no | Machine-readable reason. |
| message | no | Human-readable explanation. |
| correlationId | yes | Trace correlation ID. |
| createdAt | yes | Creation timestamp. |
| updatedAt | yes | Last update timestamp. |

## 8. Validation Rules

- Organization name must be unique.
- OrganizationUnit must reference existing Organization.
- Tenant must reference existing Organization and OrganizationUnit.
- Project must reference existing Tenant.
- ServicePlan must reference existing ServiceClass.
- ServiceInstance must reference existing Project, ServiceClass, and ServicePlan.
- ServiceBinding must reference existing ServiceInstance.
- Plugin name and version must be unique.
- Capability must reference existing Plugin.
- Operation must reference a valid target resource.

## 9. Failure Modes

| Failure | Expected Behavior |
|---|---|
| Parent resource missing | Reject with `ParentNotFound`. |
| Duplicate name | Reject with `AlreadyExists`. |
| Invalid reference | Reject with `InvalidReference`. |
| Unsupported service plan | Reject with `ServicePlanUnsupported`. |
| Plugin not found | Reject with `PluginNotFound`. |
| Capability missing | Reject with `CapabilityUnsupported`. |
| Operation failure | Mark Operation `Failed` with reason and message. |

## 10. Storage

Phase 1 may use simple storage.

Allowed options:

- in-memory registry for first skeleton,
- file-backed registry for local development,
- embedded database for early persistence.

Production metadata store is future work.

## 11. Observability

Every API request should emit:

- structured log,
- correlation ID,
- operation ID when applicable,
- error reason when applicable.

Metrics are optional in Phase 1 but must not be blocked by design.

## 12. Security

Phase 1 may use simple development authentication or no-auth local mode.

The API design must preserve fields needed for future identity and RBAC.

Do not expose secret values in API responses.

## 13. Acceptance Criteria

- API server starts.
- healthz returns healthy.
- readyz returns ready.
- Organization can be created and listed.
- OrganizationUnit can be created under Organization.
- Tenant can be created under OrganizationUnit.
- Project can be created under Tenant.
- ServiceClass and ServicePlan can be registered.
- Plugin and Capability can be registered.
- ServiceInstance request creates an Operation.
- Operation status can be retrieved.
- Invalid references fail fast.

## 14. AI Implementation Guidance

- Do not implement PostgreSQL provisioning in Phase 1.
- Do not implement a UI.
- Do not add external Keycloak/Vault/Argo dependencies yet.
- Keep interfaces clean for future Kubernetes CRD support.
- Use canonical terms from glossary.md.
- Add tests for validation and failure modes.

## Phase 2 Extension

Phase 2 extends Platform Core with reuse-first provider-neutral PaaS fabric foundations:

- API/resource standard,
- decision and audit standard,
- provider-neutral resource model,
- adapter boundary model,
- policy evaluation abstraction,
- placement decision engine,
- plugin taxonomy,
- ServiceRuntimeProfile.

`ServicePlan` remains customer-facing. `ServiceRuntimeProfile` bridges the service plan to runtime and provider capability requirements.
