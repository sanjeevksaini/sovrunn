# ADH-2026-025: ExecutionTarget Without a Mandatory ResourcePool

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-025
- **Date:** 2 August 2026
- **Classification:** Product boundary and roadmap correction
- **Canonical decision:** FCM-ADR-006

## Context

Sovrunn reuses OpenShift, Kubernetes, OpenStack, AWS, OCI and local infrastructure. These platforms already own capacity grouping and detailed scheduling. A mandatory ResourcePool and provider-wide capability model would pull Sovrunn toward infrastructure inventory and scheduler duplication rather than managed-service delivery.

## Decision

1. `InfrastructureStack` registers an existing infrastructure platform when the realization boundary is infrastructure-backed.
2. `ExecutionTarget` is the approved adapter-addressable realization boundary against which a ServiceDeploymentPlan or lifecycle action can execute.
3. An ExecutionTarget may represent an infrastructure platform, external managed-service API, edge control plane, remote service endpoint or future federated endpoint; InfrastructureStack is therefore not universally mandatory.
4. Sovrunn selects targets; native platforms retain detailed scheduling and allocation where applicable.
5. No mandatory `ResourcePool`, `ProviderCapability` or `ResourcePoolCapability` exists in the core.
6. Optional backend partitions may be referenced through a placement profile or protected adapter contract when a concrete provider need exists.
7. Compatibility is established from typed ServiceRequirementSet, ServiceTypeDefinition/runtime compatibility, normalized target facts and adapter qualification—not a provider-wide claim.

## Alternatives rejected

- **InfrastructureStack as the only target:** too coarse when one stack exposes separately governed execution boundaries.
- **Mandatory ResourcePool:** duplicates native allocation models and implies capacity ownership.
- **Provider-wide capability:** cannot say where a service can actually execute.
- **Require InfrastructureStack for every target:** excludes brokered APIs, edge control planes and future federation even though all are valid adapter-addressable realization boundaries.

## Consequences

- FEATURE-0015/0016 must be replanned around target qualification, facts and adapters.
- Capacity guarantees, reservations or backend segmentation require explicit later contracts.
- Target status must not overclaim current capacity or placement success.

## Security and sovereignty

Credentials are typed secret references, target-scoped and short-lived where possible. Customers never select target IDs or receive protected handles for normal managed-service operations.

## Conformance evidence

- Same runtime requirements qualify on two implementation types.
- A brokered external service target qualifies and executes without an InfrastructureStack while preserving the customer ServiceInstance contract.
- Unsupported target is rejected with typed reasons.
- Backend native pool/namespace identifiers remain adapter-facing.
- Target deletion is denied while active placement records depend on it.
- No core schema contains mandatory ResourcePool or provider-wide capability.

## Reassessment triggers

Three independent providers demonstrate a placement/capacity requirement that cannot be represented through ExecutionTarget, profiles, facts, reservations and adapter contracts without loss of correctness.
