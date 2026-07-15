---
doc_type: engineering_standard
title: AI Controlled Development
status: draft
phase: 1
ai_load_priority: reference
ai_summary: Operating model for founder-controlled, spec-first, test-gated AI-assisted development.
---

# AI Controlled Development

Sovrunn uses AI-assisted development, not AI-owned development.

AI tools may accelerate planning, coding, review, debugging, and documentation, but architecture remains founder-controlled, spec-first, test-gated, and terminal-verified.

## Purpose

This document defines how AI should participate in Sovrunn engineering without causing architecture drift, scope expansion, or unverifiable implementation.

Canonical context loading rules are defined in:

```text
docs/engineering/ai-context-loading-standard.md
```

Canonical Go implementation rules are defined in:

```text
docs/engineering/go-coding-guardrails.md
```

## Development Principle

```text
Human owns architecture.
Specs define scope.
AI proposes changes.
Terminal verifies truth.
Git records history.
```

## Tool Responsibilities

```text
Kiro
  architecture/spec/task planning

Cursor
  code editing/refactoring/debugging/tests

ChatGPT
  deep architecture review, exact patches, debugging, consistency review

Terminal
  source of truth for build, test, runtime, Git, and benchmark behavior
```

No tool owns a separate architecture.

## AI Must Follow Canonical Sources

AI agents must follow:

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
```

For Go work, also follow:

```text
docs/engineering/go-coding-guardrails.md
docs/architecture/controller-reconciliation-model.md
docs/architecture/observability-and-audit-baseline.md
```

## Feature Discipline

AI must implement one feature at a time.

AI must not implement future features early.

During a feature implementation, AI should load:

```text
ALWAYS context
GO_IMPLEMENTATION context, when coding Go
current FEATURE file only
```

Do not load future feature files unless explicitly requested.

## AI May Do

AI may:

```text
summarize architecture
review design consistency
generate feature plans
generate Go code
generate tests
generate documentation
suggest exact patches
debug terminal output
explain test failures
suggest performance improvements
suggest security improvements
```

## AI Must Not Do

AI must not:

```text
invent architecture
change resource hierarchy
rename canonical terms
change accepted decisions
change feature sequence
implement future features early
remove validation
remove tests
skip terminal verification
hide uncertainty
add unapproved dependencies
introduce persistent storage in Phase 1
introduce Kubernetes CRDs in Phase 1
introduce ServiceOps execution in Phase 1
introduce AI agent execution in Phase 1
introduce UI or billing in Phase 1
```

## Scope Control

Every AI coding task must clearly state:

```text
feature ID
scope
non-goals
files expected to change
tests expected
verification commands
```

AI should stop and ask for approval if the task requires:

```text
new resource kinds
API contract changes
terminology changes
feature sequence changes
new dependencies
persistent storage
distributed systems behavior
security-sensitive behavior
```

## Verification

Terminal verification is mandatory.

For Go code:

```bash
make fmt
make test
make vet
```

For concurrency-sensitive changes:

```bash
go test -race ./...
```

For documentation:

```bash
mkdocs build --strict
```

If a command fails, AI must report the failure and should not mark the task complete.

## Required AI Completion Report

For every implementation task, AI must report:

```text
feature implemented
files changed
why each file changed
tests added
validation added
security considerations
observability considerations
performance considerations
commands run
command results
non-goals intentionally not implemented
known limitations
next feature boundary
```

## Human Review Gates

Human approval is required for:

```text
architecture changes
resource model changes
API contract changes
accepted decision changes
dependency additions
security model changes
storage model changes
multi-tenant isolation changes
AI autonomy level changes
```

## Final Rule

AI accelerates Sovrunn development.

AI does not replace architecture governance, tests, terminal verification, or human judgment.
