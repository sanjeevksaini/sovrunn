---
doc_type: architecture
title: Policy Evaluation Abstraction
status: draft
phase: 2
ai_load_priority: always
ai_summary: OPA/Cedar-ready policy evaluation abstraction. Sovrunn must not hardcode policy rules into handlers or placement logic.
---

# Policy Evaluation Abstraction

## Purpose

Sovrunn should reuse mature policy engines rather than building a custom policy language or hardcoding rules into Go handlers.

## Core Contract

```text
PolicyEvaluationRequest + PolicyContext -> PolicyEvaluationResult
```

## Resources

- PolicyInput
- PolicyContext
- PolicyBundleRef
- PolicyEvaluationRequest
- PolicyEvaluationResult
- PolicyDecisionReason
- PolicyEngineAdapter

## Result Shape

```yaml
allowed: true
reasonCodes: []
humanReadableReasons: []
policyReferences: []
matchedRules: []
riskLevel: low
suggestedActions: []
evidence: []
```

## Reuse Direction

- OPA/Rego: preferred early policy engine adapter for governance, placement, security, and data policy.
- Cedar: evaluate for authorization-style questions such as principal/action/resource/context.
- Go bootstrap evaluator: allowed only for local tests or temporary MVP bootstrap behind `PolicyEngineAdapter`.

## Rules

- No hardcoded policy decisions in API handlers.
- No hardcoded policy decisions in placement engine.
- Sovrunn builds EffectivePolicyContext and Decision objects.
- Policy engines evaluate policy and return structured results.
