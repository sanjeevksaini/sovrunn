---
doc_type: decision
decision_id: DEC-0028
title: Policy Engine Adapter
status: accepted
phase: 2
related_rfc: RFC-0025
ai_load_priority: important
---

# DEC-0028: Policy Engine Adapter

## Status

Accepted

## Context

Sovrunn needs policy evaluation for governance, security, data placement, cost guardrails, entitlement, and placement decisions. Mature policy engines such as OPA and Cedar should be reused, but Phase 2 policy inputs and decision shapes are still being stabilized.

## Decision

Policy logic must go through `PolicyEngineAdapter`.

Sovrunn must not embed custom governance/security/data-placement policy rules directly in handlers, registries, or the placement engine.

Phase 2 defines:

```text
PolicyEvaluationRequest
PolicyEvaluationResult
PolicyInput
PolicyContext
PolicyBundleRef
PolicyDecisionReason
PolicyEngineAdapter
OPA adapter placeholder
Cedar adapter placeholder
```

Phase 2 does not implement full OPA/Cedar integration.

## Consequences

- Policy implementation can evolve without rewriting Sovrunn core.
- OPA is the preferred first real policy adapter candidate after the Phase 2 abstraction is stable.
- Cedar may be evaluated later for authorization/IAM-style decisions.
- Feature designs must document whether policy evaluation is required.

## Related

- RFC-0025
- docs/architecture/policy-evaluation-abstraction.md
- FEATURE-0017
