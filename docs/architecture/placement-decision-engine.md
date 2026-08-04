---
doc_type: architecture
title: Placement Decision Engine
status: draft
phase: 2
ai_load_priority: always
ai_summary: Defines PlacementRequest and PlacementDecision v0 for capability-driven, policy-aware placement simulation.
---

# Placement Decision Engine

## Purpose

The placement engine determines whether a requested service can run on a qualified ExecutionTarget safely, compliantly, and explainably.

## Inputs

- ServiceInstance request
- ServiceRuntimeProfile
- EffectiveGovernanceContext
- ServiceEntitlement / quota placeholder
- DataPlacementPolicy
- CostGuardrail
- Qualified ExecutionTarget candidates (Phase 2R; formerly ResourcePool/ProviderCapability)
- PolicyEvaluationResult

## Outputs

- PlacementDecision
- selected target
- rejected alternatives
- reason codes
- suggested actions
- audit event reference
- AI-readable DecisionContext

## Phase 2 Scope

Phase 2 placement is simulation-only. It must not provision infrastructure.

## Rule

`ServicePlan` remains customer-facing. `ServiceRuntimeProfile` bridges the customer-facing service plan to provider/runtime capability requirements.

## Phase 2R Migration Note

**ADH-2026-042 / DEC-0042**

The placement engine no longer evaluates ResourcePool inventory or ProviderCapability matching. Under the canonical model:

- Inputs use qualified ExecutionTarget candidates from adapter-provided normalized facts (FEATURE-0016).
- Policy context uses EffectiveGovernanceContext (DEC-0050), not EffectivePolicyContext.
- Sovereignty is a separate DecisionRecord profile (DEC-0043), not a policy evaluation input.
- Customer-safe ServicePlacement (DEC-0046) replaces direct topology exposure.
- ServiceRequirementSet (immutable, DEC-0044) replaces ServiceRuntimeProfile as the bridge to placement.

Canonical terminology: implementation-neutral, not provider-neutral, as the core adjective.

See `docs/architecture/canonical/sovrunn-finalized-data-model.md` for the authoritative placement model.
