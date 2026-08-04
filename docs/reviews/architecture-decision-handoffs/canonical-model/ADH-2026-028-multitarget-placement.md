# ADH-2026-028: Multi-Target Service Placement

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-028
- **Date:** 2 August 2026
- **Classification:** Placement and orchestration invariant
- **Canonical decision:** FCM-ADR-009

## Context

A managed service may require runtime, primary data, replicas, backup and disaster recovery in different failure domains or execution environments. Selecting each independently can violate latency, co-location, resilience, connectivity or sovereignty constraints.

## Decision

1. Placement selects a compatible set of one or more ExecutionTargets.
2. ServiceRequirementSet expresses component roles and cross-component relationships.
3. Placement evaluates individual target eligibility and set-level constraints.
4. The placement DecisionRecord identifies candidates, selected targets, component-role assignments, reasons, constraints and obligations.
5. ServiceDeploymentPlan is generated only from a final authoritative placement decision.
6. Connectivity is never inferred from topology; an explicit validated connectivity contract is required when the service depends on it.
7. Multi-target realization is distinct from multi-service composition. An optional versioned ServiceCompositionDefinition declares service-component dependencies, lifecycle ordering, typed outputs and compensation.
8. A component becomes a child ServiceInstance only when it has independent customer ownership, entitlement, charging, access or lifecycle; implementation-only components remain plugin-internal.

## Alternatives rejected

- **Exactly one target per service:** cannot represent independent backup/DR or composed services.
- **Independent placement per component:** misses set-level constraints.
- **Customer names targets:** leaks internals and makes portability impossible.
- **Represent every internal component as ServiceInstance:** exposes implementation detail and makes simple services operationally noisy.

## Consequences

- Placement is a constraint-set problem, not a single-target filter.
- Partial failure and compensation must be modeled in the deployment plan.
- MVP may use bounded profiles and small candidate sets while preserving the general contract.
- Composite products require explicit child lifecycle and compensation; implicit cascade is prohibited.

## Security and sovereignty

Every selected target and data movement relationship must satisfy the effective governance and sovereignty context. Protected target details remain internal.

## Conformance evidence

- One-target and two-target success cases.
- Individually eligible targets rejected when set-level latency/connectivity fails.
- Independent backup-site and separate-fault-domain cases.
- Deterministic reasons for rejected candidates.
- No topology-derived connectivity assumption.
- Composite test distinguishes child ServiceInstances from plugin-internal components and proves typed ServiceBinding exchange.

## Reassessment triggers

Scale or latency evidence requires a different solver interface, or a service family needs dynamic per-request component placement not expressible by bounded typed requirements.
