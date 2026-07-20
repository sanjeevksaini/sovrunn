---
doc_type: decision
decision_id: DEC-0036
title: Adapter Boundaries Before External Integration
status: accepted
phase: 2
related_docs:
  - docs/architecture/adapter-boundary-model.md
ai_load_priority: important
---

# DEC-0036: Adapter Boundaries Before External Integration

## Status

Accepted

## Context

Sovrunn will depend on external reusable systems for policy, identity, secrets, operations, observability, events, and persistence. Directly coupling core logic to one implementation would make later replacement and provider-neutral operation difficult.

## Decision

Adapter boundaries are required before integrating external engines expected to evolve or be replaced.

Phase 2 defines adapter boundaries for:

```text
PolicyEngineAdapter
IdentityProviderAdapter
SecretProviderAdapter
OperationEngineAdapter
ObservabilityAdapter
EventBusAdapter
Repository interfaces
```

## Consequences

- Core code depends on Sovrunn contracts, not vendor-specific or tool-specific APIs.
- Initial bootstrap implementations are allowed only behind adapters.
- Integrating OPA, Keycloak, Vault, Temporal, OpenTelemetry, Kafka/Redpanda, PostgreSQL, or YugabyteDB later should not require rewriting core decisions.

## Related

- docs/architecture/adapter-boundary-model.md
- FEATURE-0016
- DEC-0026
