---
doc_type: rfc
title: RFC-0025 Policy Evaluation Abstraction
status: draft
phase: 2
ai_load_priority: high
ai_summary: RFC for OPA/Cedar-ready policy adapter model.
---

# RFC-0025: Policy Evaluation Abstraction

See `docs/architecture/policy-evaluation-abstraction.md`.

## Decision

Sovrunn must not build a custom policy engine as core logic. It must define a policy evaluation abstraction and prioritize OPA/Cedar-compatible adapters.
