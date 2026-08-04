---
doc_type: decision_record
title: Immutable Versioned Definitions
status: accepted
phase: 2R
ai_load_priority: important
---

# Immutable Versioned Definitions

## Status

Accepted

## Context

Product, runtime, governance, sovereignty, and policy definitions could be mutable, causing inconsistent decisions from changed inputs.

## Decision

Published product, runtime, governance, sovereignty, and policy definitions are immutable by version. Once published, a version cannot be altered.

## Consequences

ServiceTypeDefinition, ServiceRuntimeProfile, GovernanceProfile, SovereigntyProfile versions are frozen after publish. New versions are new resources.

## Migration

Existing mutable profiles gain version immutability. Phase 2R fixtures enforce this.

## Conformance Evidence

Tests verify published versions cannot be altered.

## Source ADH Reference

ADH-2026-027
