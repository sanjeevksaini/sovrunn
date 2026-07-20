---
doc_type: decision
decision_id: DEC-0029
title: Separate Plugin Planes
status: accepted
phase: 2
related_rfc: RFC-0027
ai_load_priority: important
---

# DEC-0029: Separate Plugin Planes

## Status

Accepted

## Context

A single generic plugin category would blur infrastructure execution, service lifecycle planning, and runtime operations. That would make provider-specific and service-specific logic leak into Sovrunn core.

## Decision

Sovrunn separates plugin responsibilities into three planes:

```text
Provider/Substrate Plugin
  Executes infrastructure or substrate operations.

PaaS Service Management Plane Plugin
  Plans service lifecycle, validates service plans, and maps service intent to runtime requirements.

PaaS Service Runtime Plugin
  Creates/configures/checks/binds the actual runtime service through reused operators, Helm charts, APIs, or native tools.
```

Sovrunn core owns governance, policy context, placement, decisions, operations, audit, and AI-readable explanations.

## Consequences

- Core placement logic must not contain PostgreSQL lifecycle behavior.
- Provider plugins must not own service semantics.
- Runtime plugins must not bypass placement or policy decisions.
- Phase 3 uses this split for the governed PostgreSQL plugin chain.

## Related

- RFC-0027
- docs/architecture/plugin-taxonomy-and-boundaries.md
- FEATURE-0024
