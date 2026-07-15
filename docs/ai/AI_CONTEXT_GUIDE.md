---
doc_type: ai_guide
title: AI Context Guide
status: draft
phase: 1
ai_summary: Short guide pointing AI agents to the canonical context-loading and implementation standards.
---

# AI Context Guide

AI agents working on Sovrunn must follow the canonical context-loading rules in:

```text
docs/engineering/ai-context-loading-standard.md
```

That file is the authoritative source of truth for context priority and classification.

Individual file front matter is only a local hint.

If local metadata conflicts with `ai-context-loading-standard.md`, the central standard wins.

## Go Implementation

For Go implementation, AI agents must also follow:

```text
docs/engineering/go-coding-guardrails.md
```

## Core Rule

Do not load all repository files by default.

Use the smallest complete context set for the task:

```text
ALWAYS
+ current FEATURE file
+ GO_IMPLEMENTATION when coding Go
+ role-specific prompt or steering files when needed
```

## Default Go Coding Context

For a Go feature implementation, load:

```text
AGENTS.md
README.md
docs/foundation/constitution.md
docs/decisions/DECISION_INDEX.md
docs/glossary.md
docs/features/FEATURE_SEQUENCE.md
docs/resource-specs/RESOURCE_MODEL_PHASE1.md
docs/api/API_CONTRACT_PHASE1.md
docs/engineering/ai-context-loading-standard.md
docs/engineering/go-coding-guardrails.md
docs/architecture/controller-reconciliation-model.md
docs/architecture/observability-and-audit-baseline.md
current FEATURE-xxxx file
```

Do not load future feature files unless explicitly requested.

## Canonical Ownership

If files conflict, canonical ownership wins:

```text
context loading
  docs/engineering/ai-context-loading-standard.md

Go coding
  docs/engineering/go-coding-guardrails.md

terminology
  docs/glossary.md

accepted decisions
  docs/decisions/DECISION_INDEX.md

API behavior
  docs/api/API_CONTRACT_PHASE1.md

resource model
  docs/resource-specs/RESOURCE_MODEL_PHASE1.md

feature scope
  current docs/features/FEATURE-xxxx file

master AI operating contract
  AGENTS.md
```

## Archive Rule

Archive and historical files are not source of truth.

Do not load:

```text
docs/archive/
old SDE-only notes
temporary update reports
generated site/
.git/
.DS_Store
__MACOSX/
```

unless the user explicitly asks for history or comparison.

This file is a short pointer guide. It does not override `ai-context-loading-standard.md`.
