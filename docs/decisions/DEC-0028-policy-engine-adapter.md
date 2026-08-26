---
doc_type: decision
decision_id: DEC-0028
title: Policy Engine Adapter
status: accepted
phase: 2
related_rfc: RFC-0025
ai_load_priority: important
updated: 2026-08-25
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

## Phase 2R Reconciliation

ADH-2026-067, ADH-2026-068, and ADH-2026-069 reconcile the original Phase 2
inventory with the sole FEATURE-0017 architecture authority. The accepted
adapter-first decision remains unchanged.

The active FEATURE-0017 realization consists only of:

```text
PolicyEvaluationRequest
PolicyEvaluationResult
PolicyEngineAdapter
evaluation boundary
deterministic in-process fake
structural-only FEATURE-0013 mapper
transient success timing metadata for that mapper
```

The remaining original inventory is historical or represented through the
active contracts:

- `PolicyInput` is not a separate canonical type; normalized semantic input is
  carried by `PolicyEvaluationRequest` at the evaluation boundary.
- `PolicyContext` is not a separate canonical type; the request carries an
  optional UID-pinned `contextRef<EffectiveGovernanceContext>`.
- `PolicyBundleRef` is represented structurally by bounded `profileRefs`.
- `PolicyDecisionReason` is represented by normalized `reasonCodes`.
- OPA and Cedar adapter placeholders are documentation-only future candidate
  slots, not Go types, dependencies, selectors, or Phase 2R implementations.

FEATURE-0017 selects and executes no real policy engine. Any later real-engine
selection requires a fresh FEATURE-0011 reuse assessment and a separately
approved feature decision.

## Consequences

- Policy implementation can evolve without rewriting Sovrunn core.
- OPA is the preferred first real policy adapter candidate after the Phase 2 abstraction is stable.
- Cedar may be evaluated later for authorization/IAM-style decisions.
- Feature designs must document whether policy evaluation is required.

## Related

- RFC-0025
- docs/architecture/policy-evaluation-abstraction.md
- ADH-2026-067
- ADH-2026-068
- ADH-2026-069
- FEATURE-0017
