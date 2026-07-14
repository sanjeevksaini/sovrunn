---
doc_type: architecture
title: Controller and Reconciliation Model
status: draft
phase: 1
ai_load_priority: phase
ai_summary: Defines the desired-state and reconciliation model that Phase 1 resources should follow even before Kubernetes CRDs are introduced.
---

# Controller and Reconciliation Model

## 1. Purpose

Sovrunn Phase 1 starts with REST APIs and an in-memory registry.

However, resources must be designed so they can later move to Kubernetes CRDs, controllers, GitOps, and durable reconciliation.

The principle is:

```text
spec = desired state
status = observed state
operation = lifecycle trace
controller = reconciliation logic
```

## 2. Resource Shape

Every major resource should follow:

```json
{
  "apiVersion": "platform.sovrunn.io/v1alpha1",
  "kind": "ResourceKind",
  "metadata": {
    "name": "resource-name",
    "displayName": "Display Name",
    "labels": {},
    "annotations": {}
  },
  "spec": {},
  "status": {}
}
```

## 3. Metadata

`metadata` represents identity and classification.

User-controlled fields:

```text
name
displayName
labels
annotations
```

System-controlled fields, future:

```text
createdAt
updatedAt
generation
resourceVersion
```

## 4. Spec

`spec` represents desired state.

Examples:

```text
Tenant.spec.isolationProfile
Project.spec.environmentType
ServiceInstance.spec.serviceClassRef
ServiceInstance.spec.servicePlanRef
Plugin.spec.pluginType
Capability.spec.operation
```

Users and GitOps may author `spec`.

## 5. Status

`status` represents observed state.

Users must not directly set `status`.

Examples:

```text
status.phase
status.message
status.conditions
status.lastObservedAt
status.operationRef
```

## 6. Operation

`Operation` records lifecycle activity.

Examples:

```text
CreateOrganization
CreateTenant
CreateProject
CreateServiceInstance
DeleteServiceBinding
```

Operations are not desired state. They are lifecycle traces.

## 7. Phase 1 Reconciliation Behavior

In Phase 1, reconciliation may be synchronous and simple:

```text
API request
  -> validate request
  -> write desired resource to registry
  -> create operation record
  -> set status
  -> return response
```

## 8. Future Reconciliation Behavior

Later phases may use asynchronous controllers:

```text
API request or GitOps apply
  -> store desired state
  -> controller observes desired state
  -> controller validates dependencies
  -> controller calls ServiceOps plugin
  -> controller updates status
  -> operation records lifecycle
```

## 9. Idempotency

Controllers and lifecycle handlers must be idempotent.

Repeated execution of the same desired state should not create duplicate external resources.

Phase 1 registry methods should prepare for this by enforcing:

```text
unique resource names
deterministic validation
safe update behavior
delete-blocked behavior
clear conflict errors
```

## 10. Conditions

Future resource status may include conditions:

```json
{
  "type": "Ready",
  "status": "True",
  "reason": "ValidationSucceeded",
  "message": "Resource is ready",
  "lastTransitionTime": "2026-07-11T00:00:00Z"
}
```

Phase 1 may start with simple `phase` and `message`.

## 11. Non-Goals for Phase 1

Do not implement:

```text
Kubernetes CRDs
controller-runtime
operator reconciliation loops
GitOps sync controller
external infrastructure provisioning
plugin execution
durable event queues
```

## 12. Acceptance Criteria

Phase 1 code follows this model if:

```text
all resources have metadata/spec/status shape
status is not user-authored
mutating requests create operation records after FEATURE-0005
delete-blocked behavior protects child resources
validation is deterministic
resources can be serialized as JSON/YAML
```

## 13. Final Principle

Phase 1 may be simple internally, but the public resource model must be ready for future desired-state reconciliation.
