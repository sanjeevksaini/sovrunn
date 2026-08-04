---
doc_type: decision_record
title: DecisionRecord Profiles for Sovereignty and Placement
status: accepted
phase: 2R
ai_load_priority: important
---

# DecisionRecord Profiles for Sovereignty and Placement

## Status

Accepted

## Context

Sovereignty assessment and placement decision need governed conclusion envelopes but should not create new envelope kinds.

## Decision

SovereigntyAssessment and PlacementDecision use registered FEATURE-0013 DecisionRecord profiles. They are DecisionProfile extensions, not separate resource kinds.

## Consequences

DecisionRecord remains the shared envelope. New profile types are data, not schema. FEATURE-0013 ownership is preserved.

## Migration

No existing sovereignty or placement records need migration; this establishes the canonical pattern for Phase 2R.

## Conformance Evidence

DecisionRecord profiles are registered and testable without new envelope kinds.

## Source ADH Reference

ADH-2026-026
