---
doc_type: rfc
title: RFC-0026 Placement Decision Engine v0
status: draft
phase: 2
ai_load_priority: high
ai_summary: RFC for first placement decision engine scope.
---

# RFC-0026: Placement Decision Engine v0

See `docs/architecture/placement-decision-engine.md`.

## Decision

Phase 2 placement is simulation-only. Phase 3 consumes PlacementDecision to trigger one executable PostgreSQL plugin chain.

## Phase 2R Migration Note

**ADH-2026-042 / DEC-0042:** Placement no longer evaluates ResourcePool inventory or ProviderCapability matching. Phase 2R placement uses qualified ExecutionTarget candidates from FEATURE-0016 adapter qualification and FEATURE-0013 DecisionRecord profiles. The PlacementDecision contract remains valid; its inputs change from ResourcePool/ProviderCapability to qualified ExecutionTarget candidates with explicit cross-target constraints (DEC-0045). Customer-safe ServicePlacement projection (DEC-0046) replaces direct topology exposure.
