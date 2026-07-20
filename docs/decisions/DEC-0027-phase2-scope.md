---
doc_type: decision
decision_id: DEC-0027
title: Phase 2 Scope
status: accepted
phase: 2
related_docs:
  - docs/phase2/PHASE2_SCOPE.md
  - docs/architecture/development-phases.md
ai_load_priority: important
---

# DEC-0027: Phase 2 Scope

## Status

Accepted

## Context

Sovrunn Phase 2 must create the foundation for provider-neutral, reuse-first PaaS governance without turning into real provisioning or production integration work too early.

## Decision

Phase 2 builds only:

```text
model foundation
decision foundation
audit foundation
policy context foundation
adapter boundaries
plugin taxonomy
placement simulation
AI-readable decision context
```

Phase 2 does not build:

```text
real provider provisioning
real PostgreSQL runtime provisioning
full OPA/Cedar integration
full Keycloak/Vault/Temporal integration
global traffic execution
autoscaling execution
billing engine
full compliance engine
autonomous AI operations
```

## Consequences

- Phase 2 features must pass `scripts/phase2-scope-check.sh`.
- Future-phase concepts may appear as Non-goals, Out of Scope, Deferred, or Roadmap Placeholder only.
- Phase 3 owns the first executable PostgreSQL plugin chain.

## Related

- docs/phase2/PHASE2_SCOPE.md
- docs/phase2/PHASE2_FEATURE_SEQUENCE.md
- docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md
