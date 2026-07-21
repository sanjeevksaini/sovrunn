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

The placement engine determines whether a requested service can run on a provider ResourcePool safely, compliantly, and explainably.

## Inputs

- ServiceInstance request
- ServiceRuntimeProfile
- EffectivePolicyContext
- ServiceEntitlement / quota placeholder
- DataPlacementPolicy
- CostGuardrail
- Provider / ResourcePool / ProviderCapability
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
