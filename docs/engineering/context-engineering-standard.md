---
doc_type: engineering_standard
title: Context Engineering Standard
status: draft
phase: 1
ai_summary: High-level context-engineering principle. Canonical loading rules are defined in ai-context-loading-standard.md.
---

# Context Engineering Standard

Sovrunn uses context engineering to keep AI-assisted development accurate, scoped, and token-efficient.

This file defines the high-level principle only.

Canonical context selection rules are defined in:

```text
docs/engineering/ai-context-loading-standard.md
```

Go implementation guardrails are defined in:

```text
docs/engineering/go-coding-guardrails.md
```

This file does not override the canonical loading matrix.

## Principle

Load the minimum complete context required for the task.

```text
Enough context to stay correct.
Not so much context that the AI becomes slow, confused, repetitive, or tempted to implement future phases early.
```

## Authoritative Classification

Context priority and classification are defined centrally in:

```text
docs/engineering/ai-context-loading-standard.md
```

Individual file front matter is only a local hint.

If local metadata conflicts with the central standard, the central standard wins.

## What Context Engineering Prevents

Context engineering prevents:

```text
architecture drift
terminology drift
future feature leakage
duplicate implementation
scope expansion
token waste
AI hallucination from missing files
AI confusion from historical notes
```

## Context Selection Rule

Use:

```text
small always-load core
+ current feature file
+ implementation-specific guardrails
+ role-specific prompt or steering file when needed
```

Avoid:

```text
loading the entire repository
loading all feature files
loading all RFCs
loading archive files
loading generated docs site output
loading stale SDE-only notes for Phase 1 platform implementation
```

## Canonical Ownership

If documents conflict, use this ownership model:

```text
AGENTS.md
  short master operating contract

docs/engineering/ai-context-loading-standard.md
  authoritative context classification and loading matrix

docs/engineering/go-coding-guardrails.md
  canonical Go implementation rules

docs/ai/AI_CONTEXT_GUIDE.md
  short pointer to canonical standards

docs/engineering/context-engineering-standard.md
  high-level principle only

.cursor/rules/*
  short Cursor enforcement rules

.kiro/steering/*
  short Kiro steering rules
```

## Non-Goal

This file is not a replacement for:

```text
docs/engineering/ai-context-loading-standard.md
docs/engineering/go-coding-guardrails.md
AGENTS.md
```

It should remain short and conceptual.
