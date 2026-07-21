---
doc_type: decision
decision_id: DEC-0026
title: Reuse Before Build
status: accepted
phase: 2
related_rfc: RFC-0021
ai_load_priority: important
---

# DEC-0026: Reuse Before Build

## Status

Accepted

## Context

Sovrunn is a sovereign PaaS platform that must integrate with mature cloud-native and open-source foundations instead of rebuilding them. Rebuilding policy engines, workflow engines, identity providers, secret stores, observability stacks, database operators, or traffic systems would slow the MVP and increase maintenance risk.

## Decision

Sovrunn follows reuse-before-build as a core engineering rule.

For every feature, teams must classify implementation as:

```text
Reuse / Wrap / Extend / Build
```

Sovrunn should build the platform-specific control-plane intelligence:

```text
governance
policy context
placement decisions
plugin contracts
operation tracking
audit/evidence
AI-readable explanations
customer/provider experience
```

Sovrunn should reuse mature foundations such as Kubernetes, OPA/Cedar, Keycloak/Dex, Vault/External Secrets, OpenTelemetry, Prometheus/Grafana, Temporal/Argo, PostgreSQL operators, Helm, GitOps tools, and service mesh/traffic components where suitable.

## Consequences

- Every feature requires a Reuse Assessment.
- Architecture reviews must reject unnecessary reinvention.
- Adapter boundaries are required when integrating reusable foundations that may evolve or be replaced.
- Building custom infrastructure is allowed only when it is core Sovrunn differentiation and explicitly approved.

## Related

- RFC-0021
- DEC-0036
- docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md
