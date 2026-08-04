# Sovrunn Product Steering

## Product Identity

Sovrunn is a sovereign cloud-native PaaS platform.

It provides an opinionated product layer above proven open-source and open-standard infrastructure.

## SDE Positioning

SDE is a major differentiated capability inside Sovrunn. SDE is not the entire platform.

## Target Customers

```text
government-scale platforms
local cloud providers
local colocation providers
regulated enterprises
on-prem/cloud enterprises
system integrators
```

## Product Goals

Sovrunn should provide:

```text
organization-first governance
multi-tenant service consumption
service catalog (ServiceTypeDefinition + ServiceOffering + ServicePlan)
service plans (immutable, versioned)
Service Management Plane registry
ServiceOps plugin framework
capability registry
operation framework
policy inheritance (EffectiveGovernanceContext)
sovereignty evidence and assessment
audit aggregation
backup and archival governance
cloud management across multiple sovereign datacenter locations
AI-assisted operations
SDE
```

## Phase 2R Goal

Phase 2R builds the canonical model foundation:

```text
CloudPlatform, CloudProvider, CloudProviderParticipation, CloudEnrollment
HostingLocation, Datacenter, FaultDomain, InfrastructureStack, ExecutionTarget
ServiceTypeDefinition, ServiceOffering, ServicePlan, ServiceRequirementSet
EffectiveGovernanceContext, SovereigntyProfile, EvidenceRecord
DecisionRecord profiles (sovereignty, placement, governance)
EntitlementPackage, QuotaPolicy
PluginExecution (synthetic), ServicePlacement (safe projection)
ServiceBinding (SecretRef-only, per-consumer)
VS-000 Slice 0 conformance demo
```

## Phase 2R Non-Goals

Do not build yet:

```text
production UI
billing engine
marketplace
multi-cluster federation implementation
persistent database storage
Kubernetes CRDs
GitOps controller
real ServiceOps plugin execution
real datastore provisioning
AI agent execution
SDE transformation
mandatory ResourcePool or ProviderCapability
native provider objects in customer schemas
```

## Product Rule

Expose simple, governed service consumption. Do not expose raw Kubernetes complexity to tenants. Do not expose topology, credentials, or protected handles to customers.
