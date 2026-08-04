---
doc_type: decision_record
title: Geography Governance Sovereignty Independent
status: Accepted
phase: 2R
ai_load_priority: important
---

# Geography Governance Sovereignty Independent

## Status

Accepted

## Context

Physical geography, governance scope, and sovereignty could be conflated as a single location hierarchy.

## Decision

Physical geography, governance scope, and sovereignty are independent dimensions. Geography alone never proves sovereignty.

## Consequences

HostingLocation is a registered fact. Sovereignty requires separate DecisionRecord evidence. Governance is resolved independently through EffectiveGovernanceContext.

## Migration

Existing ProviderLocation becomes HostingLocation registered fact.

## Conformance Evidence

Drift checks verify geography alone does not prove sovereignty in any active authority.

## Source ADH Reference

ADH-2026-024
