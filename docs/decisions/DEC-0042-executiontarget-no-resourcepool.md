---
doc_type: decision_record
title: ExecutionTarget No ResourcePool
status: Accepted
supersedes: DEC-0032, DEC-0033
phase: 2R
ai_load_priority: important
---

# ExecutionTarget No ResourcePool

## Status

Accepted

Supersedes: DEC-0032, DEC-0033

## Context

ResourcePool was defined as mandatory placement boundary (DEC-0032) and ProviderCapability as mandatory compatibility boundary (DEC-0033). These synthetic pools pull core toward infrastructure ownership and duplicate backend schedulers.

## Decision

Qualified ExecutionTarget is the actionable realization boundary. No mandatory ResourcePool or provider-wide ProviderCapability exists in core.

## Consequences

Placement evaluates qualified ExecutionTarget candidates, not synthetic pools. ProviderCapability and ResourcePool are not active mandatory concepts.

## Migration

DEC-0032 and DEC-0033 are superseded. RFC-0024 placement spine based on ResourcePool/ProviderCapability is replaced.

## Conformance Evidence

No active Phase 2R feature requires ResourcePool or ProviderCapability.

## Source ADH Reference

ADH-2026-025
