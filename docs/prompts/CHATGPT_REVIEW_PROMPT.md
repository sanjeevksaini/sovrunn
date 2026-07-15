# ChatGPT Review Prompt

Use this prompt when asking ChatGPT to review Sovrunn architecture, docs, or code.

## Prompt

You are reviewing Sovrunn Phase 1 Platform Core.

Use these as source of truth:

```text
AGENTS.md
docs/foundation/constitution.md
docs/decisions/DECISION_INDEX.md
docs/glossary.md
docs/engineering/ai-controlled-development.md
docs/engineering/context-engineering-standard.md
docs/architecture/controller-reconciliation-model.md
docs/architecture/observability-and-audit-baseline.md
docs/architecture/gitops-desired-state-model.md
docs/features/FEATURE_SEQUENCE.md
docs/resource-specs/RESOURCE_MODEL_PHASE1.md
docs/api/API_CONTRACT_PHASE1.md
current feature file
```

Review for architecture consistency, scope control, terminology correctness, missing validation, missing tests, API contract mismatch, security concerns, observability/audit gaps, and future-feature leakage.

Return:

```text
summary
critical issues
minor issues
recommended patch
tests to run
go/no-go recommendation
```
