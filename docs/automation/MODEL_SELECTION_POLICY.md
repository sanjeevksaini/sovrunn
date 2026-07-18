# Sovrunn Model Selection Policy

This policy tells Feature Factory how to recommend three prioritized LLM models for every Kiro stage and Cursor task.

## Goals

1. ChatGPT/reviewer recommends models based on stage and task demand.
2. Kiro and Cursor prompts include exactly three prioritized model recommendations.
3. Kiro and Cursor outputs must report which model was actually used.
4. If a recommended model is unavailable, use the next model in the priority list and report the fallback.

## Important limitation

Model availability changes by product, account, plan, region, and workspace governance. The recommendation list is a policy preference, not proof of availability. Always use the available model that best matches the highest-priority recommendation.

## Standard Kiro recommendations

### requirements.md

1. Claude Opus 4.8 — High effort
2. GPT-5.6 Sol — High effort
3. Claude Sonnet 5 — High effort

Use requirements recommendations when the work is about scope, user stories, acceptance criteria, non-goals, and architecture boundary control.

### design.md

1. Claude Opus 4.8 — XHigh effort
2. GPT-5.6 Sol — High effort
3. Claude Sonnet 5 — High effort

Use design recommendations when the work is about architecture, interfaces, data models, API behavior, edge cases, and implementation strategy.

### tasks.md

1. GPT-5.6 Sol — High effort
2. Claude Opus 4.8 — High effort
3. GPT-5.6 Terra — Medium effort

Use tasks recommendations when the work is about task decomposition, commit boundaries, verification commands, and Cursor task prompts.

### revision prompts

1. Claude Opus 4.8 — High effort
2. GPT-5.6 Sol — High effort
3. Claude Sonnet 5 — Medium effort

Use revision recommendations when Kiro must update an existing requirements/design/tasks document.

## Standard Cursor recommendations

### complex_code

1. GPT-5.6 Sol — High effort
2. Claude Opus 4.8 — High effort
3. Claude Sonnet 5 — High effort

Use for handlers, server wiring, registries, cross-file changes, refactors, concurrency, and architecture-sensitive implementation.

### routine_code

1. GPT-5.6 Terra — Medium effort
2. Claude Sonnet 5 — Medium effort
3. GPT-5.6 Luna — Medium effort

Use for small focused implementation tasks with limited risk.

### tests

1. GPT-5.6 Terra — Medium effort
2. Claude Sonnet 5 — Medium effort
3. GPT-5.6 Luna — Low effort

Use for unit tests, property tests, race tests, table-driven tests, and test-only tasks.

### debug_fix

1. Claude Opus 4.8 — XHigh effort
2. GPT-5.6 Sol — High effort
3. Claude Sonnet 5 — High effort

Use when verification fails or the task is explicitly a bug/debug/failure fix.

## Required execution report

Every Kiro and Cursor output must include this block:

```text
Model Execution Report:
- Tool: Kiro | Cursor
- Stage or task: <requirements | design | tasks | Task N.N>
- Recommended priority list: <copied from prompt>
- Selected model: <actual model selected>
- Effort/reasoning setting: <actual setting if visible>
- Fallback used: yes/no
- Fallback reason: <unavailable/cost/latency/user override/none>
```

If the tool cannot report the actual model, it must write:

```text
Selected model: unknown/not visible in tool output
```

## Governance rule

Do not let the model choose a fourth option unless all three recommended models are unavailable or blocked by governance. If a fourth model is used, the execution report must explain why.
