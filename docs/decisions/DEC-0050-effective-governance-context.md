---
doc_type: decision_record
title: Effective Governance Context
status: Accepted
phase: 2R
ai_load_priority: important
---

# Effective Governance Context

## Status

Accepted

## Context

Multiple governance controls (governance profile, security profile, data placement, cost guardrail) needed composition into one resolved context, but sovereignty must remain separately assessed.

## Decision

Governance controls compose into one EffectiveGovernanceContext. Sovereignty remains separately assessed through DecisionRecord evidence.

## Consequences

EffectiveGovernanceContext is the single resolved governance term. Sovereignty is not weakened by governance composition. Controls are non-weakenable by default.

## Migration

Existing EffectivePolicyContext references migrate to EffectiveGovernanceContext.

## Conformance Evidence

EffectiveGovernanceContext is the only new canonical resolved-governance term. Non-weakenable control tests pass.

## Source ADH Reference

ADH-2026-033
