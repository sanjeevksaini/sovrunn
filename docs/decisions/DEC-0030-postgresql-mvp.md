---
doc_type: decision
decision_id: DEC-0030
title: Governed PostgreSQL MVP
status: accepted
phase: 3
related_rfc: RFC-0029
ai_load_priority: important
---

# DEC-0030: Governed PostgreSQL MVP

## Status

Accepted

## Context

Sovrunn needs a customer-testable MVP that proves its core value proposition without attempting broad multi-cloud or multi-service scope too early.

## Decision

MVP-001 is:

```text
Governed PostgreSQL PaaS Placement and Provisioning on one substrate.
```

The MVP demonstrates:

```text
service request
entitlement check
effective policy context
placement decision
operation tracking
plugin execution
service instance status
service binding
audit event
AI-readable explanation
```

The actual PostgreSQL runtime should reuse mature foundations such as CloudNativePG, Crunchy Postgres Operator, or Helm rather than building PostgreSQL HA/lifecycle from scratch.

## Consequences

- Phase 3 implements one narrow executable plugin chain.
- PostgreSQL lifecycle logic belongs in plugin wrappers, not core.
- Advanced HA, backup, failover, autoscaling, and multi-provider execution are deferred.

## Related

- RFC-0029
- docs/mvp/MVP_001_GOVERNED_POSTGRESQL_PAAS.md
- FEATURE-0027 through FEATURE-0034
