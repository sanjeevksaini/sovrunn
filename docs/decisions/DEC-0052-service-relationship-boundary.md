---
doc_type: decision_record
title: Service Relationship Boundary
status: Accepted
phase: 2R
ai_load_priority: important
---

# Service Relationship Boundary

## Status

Accepted

## Context

Cross-service dependencies (e.g., application depends on database, or service depends on cache) need explicit modeling without creating implicit runtime coupling.

## Decision

ServiceRelationshipDefinition and ServiceRelationship model implementation-neutral cross-service semantics and lifecycle.

## Consequences

Service dependencies are explicit, typed, and lifecycle-aware. They do not create runtime coupling or bypass individual service governance.

## Migration

No existing cross-service relationships need migration; this establishes the forward contract.

## Conformance Evidence

Phase 2R models service relationships in synthetic scenarios without runtime coupling.

## Source ADH Reference

ADH-2026-035
