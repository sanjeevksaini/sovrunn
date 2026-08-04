---
doc_type: decision_record
title: Release Recovery Compatibility
status: Accepted
phase: 2R
ai_load_priority: important
---

# Release Recovery Compatibility

## Status

Accepted

## Context

Release compatibility and recovery actions could be implicit or undirected.

## Decision

Release compatibility is directed; recovery actions and lifecycle-reference contracts are explicit and pinned.

## Consequences

Upgrades follow explicit compatibility direction. Recovery actions reference specific pinned versions. No implicit forward-compatibility assumptions.

## Migration

No existing release compatibility records need migration; establishes the forward contract.

## Conformance Evidence

Release and recovery contracts are explicit in Phase 2R models.

## Source ADH Reference

ADH-2026-039
