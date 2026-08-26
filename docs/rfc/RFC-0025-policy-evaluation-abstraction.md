---
doc_type: rfc
title: RFC-0025 Policy Evaluation Abstraction
status: approved
phase: 2R
approved_by: ADH-2026-067
clarified_by: ADH-2026-068
corrected_by: ADH-2026-069
updated: 2026-08-25
ai_load_priority: high
ai_summary: Approved engine-neutral FEATURE-0017 seam with a deterministic in-process fake and deferred real-engine integration.
---

# RFC-0025: Policy Evaluation Abstraction

Normative detail is owned solely by
`docs/architecture/policy-evaluation-abstraction.md`. This RFC summarizes the
direction and does not duplicate its six decision groups.

## Decision

Sovrunn must not build a custom policy engine or embed policy rules in core
handlers, registries, or placement logic.

All real policy logic crosses the engine-neutral `PolicyEngineAdapter`. The
adapter receives normalized semantic input plus its canonical digest and
returns either a normalized conclusion or `AdapterFailure`. The evaluation
boundary validates that conclusion and constructs the complete transient
`PolicyEvaluationResult`. Successful evaluation transports that result with
its exact boundary-captured FEATURE-0013 timing envelope as private transient
metadata for the separate pure mapper; this does not create a canonical type,
resource, route, or store.

Phase 2R provides only a deterministic in-process fake. It selects or executes
no OPA, Cedar, other policy engine, adapter selector, policy loader, external
integration, or production topology.

A later real-engine feature must perform a fresh reuse assessment and preserve
the canonical FEATURE-0017 port without exposing engine-native types in
Sovrunn core.
