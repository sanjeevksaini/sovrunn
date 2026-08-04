# Sovrunn Architecture Steering

## Source Documents

```text
docs/foundation/vision.md
docs/foundation/philosophy.md
docs/foundation/constitution.md
docs/decisions/DECISION_INDEX.md
docs/glossary.md
docs/architecture/canonical/sovrunn-finalized-data-model.md
docs/architecture/canonical/sovrunn-final-canonical-contract-catalog.md
docs/architecture/reference-flows/postgresql-end-to-end.md
docs/architecture/platform-core.md
docs/architecture/organization-governance.md
docs/architecture/controller-reconciliation-model.md
docs/architecture/observability-and-audit-baseline.md
docs/architecture/gitops-desired-state-model.md
docs/architecture/api-resource-standard.md
docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md
docs/phase2/PHASE2R_REBASELINE.md
docs/architecture/vertical-slices/VS-000-core-skeleton.md
docs/architecture/vertical-slices/VS-000-contract-specification.md
docs/architecture/vertical-slices/VS-000-contract-registry.yaml
docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md
.kiro/steering/slice0-contract.md
```

## Architecture Baseline

`ARCH-2026.08-PHASE2R-CANONICAL`

Controlling adoption: `ADH-2026-042`, `ACR-2026-001`

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

## Canonical Cloud Model

```text
CloudPlatform (product ownership)
CloudProvider (infrastructure operation)
CloudProviderParticipation (provider participates in platform)
CloudEnrollment (customer Organization joins CloudPlatform)
SovrunnInstallation (one provider per installation)
HostingLocation, Datacenter, FaultDomain, InfrastructureStack
ExecutionTarget (qualified realization boundary)
```

## Canonical Service Model

```text
ServiceTypeDefinition (reusable, implementation-neutral)
ServiceOffering (CloudPlatform-scoped product)
ServicePlan (versioned, immutable, customer-facing)
ServiceRuntimeProfile / ServiceRequirementSet (internal bridge)
ServiceDeploymentPlan (immutable ExecutionTarget set)
ServicePlacement (safe customer projection)
ServiceBinding (per-consumer, SecretRef-only)
```

## Governance Model

```text
EffectiveGovernanceContext (resolved composition)
SovereigntyProfile / SovereigntyFactSet / EvidenceRecord
DecisionRecord with profiles (sovereignty, placement, governance)
EntitlementPackage / ServiceEntitlement
QuotaPolicy (independent from entitlement)
```

## Seven Canonical Scope Kinds

```text
Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider
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

CloudPlatform -> ServiceOffering -> ServicePlan

Plugin registry: Plugin -> Capability

## Phase 2R Storage

Use in-memory registry only.

Do not introduce durable storage, Kubernetes CRDs, or database persistence in Phase 2R unless explicitly approved.

## Superseded Concepts

Do not use as active targets:

- ResourcePool, ProviderCapability (DEC-0042)
- Generic Provider as combined concept (DEC-0037)
- Global ServiceClass as canonical catalog (DEC-0049)
- EffectivePolicyContext (use EffectiveGovernanceContext per DEC-0050)
- Six-scope vocabulary (use seven-scope per DEC-0037)
- SovrunnInstallation as active Slice 0 resource (no Phase 2R owner; DEC-0053 deferred)

## Review Rule

If a proposed change modifies resource shape, feature sequence, terminology, architecture boundaries, or Phase 2R scope, stop and ask for founder approval.
