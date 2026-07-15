---
doc_type: ai_prompt
title: AI Implementation Prompt for Phase 1
status: draft
phase: 1
ai_load_priority: always
ai_summary: Prompt to give to an AI coding assistant when implementing Phase 1.
---

# AI Implementation Prompt for Phase 1

You are implementing Sovrunn Phase 1 Platform Core.

## Load First

Read:

```text
docs/foundation/constitution.md
docs/decisions/DECISION_INDEX.md
docs/glossary.md
docs/architecture/platform-core.md
docs/architecture/organization-governance.md
docs/ai/AI_FEATURE_FACTORY.md
docs/features/FEATURE_SEQUENCE.md
docs/engineering/go-coding-guardrails.md
```

Then load the specific feature file being implemented.

## Rules

- Implement one feature at a time.
- Do not expand scope beyond the feature file.
- Use Go.
- Use in-memory registry first.
- Keep API simple and deterministic.
- Add tests for validation and registry behavior.
- Add curl examples when useful.
- Do not introduce a database unless explicitly requested.
- Do not introduce Kubernetes CRDs yet.
- Do not introduce plugin execution yet.
- Do not introduce AI agent execution yet.
- Do not introduce UI.

## Output Required Per Feature

For each feature, produce:

```text
1. package/file changes
2. API handlers
3. registry methods
4. validation rules
5. tests
6. demo curl commands
7. acceptance checklist
```

## Quality Bar

Code must be:

```text
simple
explicit
testable
idiomatic Go
low dependency
easy for future CRD/storage migration
```
