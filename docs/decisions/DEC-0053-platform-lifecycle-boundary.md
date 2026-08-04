---
doc_type: decision_record
title: Platform Lifecycle Boundary
status: Accepted
phase: 2R
ai_load_priority: important
---

# Platform Lifecycle Boundary

## Status

Accepted

## Context

Platform lifecycle (install, upgrade, recovery) could be conflated with service lifecycle or scattered across ad-hoc mechanisms.

## Decision

Sovrunn platform lifecycle uses one external accepted-intent authority, Operation-first activation, and an independently recoverable agent.

## Consequences

Platform changes follow the same Operation/DecisionRecord pattern as service changes. Recovery is independent. Lifecycle has one external intent authority.

## Migration

Phase 2R preserves lifecycle contracts and topology decisions for the parallel production stream.

## Conformance Evidence

Platform lifecycle contracts are explicit and do not conflate with service lifecycle.

## Source ADH Reference

ADH-2026-036
