# Current Phase Context

Current active phase: Phase 2.

Architecture baseline: `ARCH-2026.07-PHASE2-START`.

FEATURE-0013 status: implemented and merged through PR #15 on 2026-07-29. ADH-2026-017 remains the approved single replacement architecture handoff; ADH-2026-014/015/016 are historical provenance and not downstream inputs.

FEATURE-0014 status: implemented and merged through PR #16 as commit `1ed47ac` on 2026-07-30. ADH-2026-018 remains controlling, with the geographic-descriptor clarification in ADH-2026-019.

## Phase 2 Goal

Establish the reuse-first, provider-neutral PaaS fabric foundation required before executable plugin-chain development.

## Phase 2 Build Scope

- Reuse assessment standard
- API/resource standard
- Decision and AuditEvent standard (DecisionRecord envelope, FEATURE-0012 six-scope governance vocabulary, ServiceInstance as typed subject, AuditEvent via metadata.scopeRef)
- Provider-neutral resource model
- ResourcePool and ProviderCapability model
- Adapter boundary foundation
- Policy evaluation abstraction
- GovernanceProfile and SecurityProfile foundation
- DataPlacementPolicy and CostGuardrail minimal foundation
- ProfileAssignment and EffectivePolicyContext
- Minimal ServiceEntitlement and quota placeholder
- ServiceRuntimeProfile foundation
- PlacementRequest and PlacementDecision v0
- Plugin taxonomy foundation
- AI-readable decision context
- Phase 2 integration demo

## Phase 2 Completed Features

| Feature | Status | Approval |
|---|---|---|
| FEATURE-0011 Reuse Assessment Standard | Merged | Complete |
| FEATURE-0012 API, Resource Naming, Status, and Validation Standard | Implemented and merged through PR #14 as commit `a1b74fb` into `phase2-reuse-first-paas-fabric-foundation` | Final human approval 2026-07-24 |
| FEATURE-0013 Decision Record and AuditEvent Standard | Implemented and merged through PR #15 into `phase2-reuse-first-paas-fabric-foundation` | Final human/Codex review 2026-07-29 |
| FEATURE-0014 Provider-Neutral Resource Model | Implemented and merged through PR #16 as commit `1ed47ac` into `phase2-reuse-first-paas-fabric-foundation` | Final feature gate passed 2026-07-30 |

## Phase 2 Next Planned Feature

| Feature | Status | Controlling handoff |
|---|---|---|
| FEATURE-0015 ResourcePool and ProviderCapability Model | Architecture not started | Pending architecture decision handoff |


## FEATURE-0014 Merged Architecture Boundary

Later features consume FEATURE-0014 through its approved architecture, merged contracts, and this focused dependency boundary:

- FEATURE-0011 controls the mandatory reuse-assessment format and reuse-before-build gate.
- FEATURE-0012 controls API/resource grammar, metadata, references, status, validation shape, Problem Details envelope, and the six-value `ScopeKind` vocabulary.
- FEATURE-0013 controls `DecisionRecord`, `DecisionProfile`, `EvaluationResult`, `AuditEvent`, `metadata.scopeRef` as sole scope authority, and the rule that `ServiceInstance` is a typed subject, not a `ScopeKind`.
- FEATURE-0014 owns only provider-neutral substrate resources: Provider, ProviderLocation, ProviderDatacenter, DatacenterFailureDomain, and InfrastructureStack.
- FEATURE-0014 must not own `ResourcePool` or `ProviderCapability`; those belong to FEATURE-0015.
- FEATURE-0014 must not define adapter interfaces; those belong to FEATURE-0016.
- FEATURE-0014 must not introduce provider-specific runtime provisioning, plugin execution, placement decisions, policy evaluation, or new decision/audit envelopes.
- FEATURE-0014 is `NOT_APPLICABLE` under the FEATURE-0013 adoption contract and must not invent decision, audit, or operation behavior.
- Physical containment implies neither network connectivity nor isolation; explicit connectivity remains outside FEATURE-0014.

## Phase 2 Exit Criteria

- All Phase 2 features pass feature gate.
- Placement simulation demonstrates explainable allowed/denied outcomes.
- AuditEvent records are produced for required decisions.
- Adapter boundaries are present before real integrations.
- Architecture drift review is approved.
- Phase 3 readiness review is complete.
