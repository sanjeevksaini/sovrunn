# Sovrunn Architecture Steering

## Source Documents

```text
docs/foundation/vision.md
docs/foundation/philosophy.md
docs/foundation/constitution.md
docs/decisions/DECISION_INDEX.md
docs/glossary.md
docs/architecture/platform-core.md
docs/architecture/organization-governance.md
docs/architecture/controller-reconciliation-model.md
docs/architecture/observability-and-audit-baseline.md
docs/architecture/gitops-desired-state-model.md
docs/resource-specs/RESOURCE_MODEL_PHASE1.md
docs/api/API_CONTRACT_PHASE1.md
```

## Architecture Layers

```text
Users and Access Channels
  -> Portal, CLI, API, GitOps, SDKs, AI Assistant

Organization Management Layer
  -> Organization, OrganizationUnit, Tenant, Project

Cloud Management Plane
  -> API server, resource registry, service catalog,
     Service Management Plane registry, plugin registry,
     capability registry, operation framework

Service Management Planes
  -> datastore, cache, object storage, stream, gateway,
     load balancer, FaaS, big data, SDE

Execution Substrate
  -> Kubernetes, operators, GitOps, policy, observability,
     identity, secrets, storage, networking
```

## Resource Model

```text
metadata = identity and classification
spec     = desired state
status   = observed state
operation = lifecycle trace
```

## Hierarchy

```text
Organization
  -> OrganizationUnit
      -> Tenant
          -> Project
              -> ServiceInstance
                  -> ServiceBinding
```

Catalog:

```text
ServiceClass
  -> ServicePlan
```

Plugin registry:

```text
Plugin
  -> Capability
```

## Phase 1 Storage

Use in-memory registry only.

Do not introduce durable storage, Kubernetes CRDs, or database persistence in Phase 1 unless explicitly approved.

## Review Rule

If a proposed change modifies resource shape, feature sequence, terminology, architecture boundaries, or Phase 1 scope, stop and ask for founder approval.
