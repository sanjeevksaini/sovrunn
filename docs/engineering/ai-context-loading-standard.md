---
doc_type: engineering_standard
title: Sovrunn AI Context Loading Standard
status: draft
phase: 1
ai_load_priority: always
ai_summary: Canonical standard and authoritative registry for deciding which Sovrunn repository files AI coding and planning agents should load for each task type.
---

# Sovrunn AI Context Loading Standard

## 1. Purpose

This document is the **canonical source of truth** for selecting AI context files in Sovrunn.

It defines:

```text
which files are always loaded
which files are task-specific
which files are feature-specific
which files are role-specific
which files are reference-only
which files must not be loaded
```

It applies to:

```text
Kiro
Cursor
ChatGPT
Claude Code
other AI coding or review agents
```

The goal is to give AI agents enough context to stay correct without loading so much material that they become slow, confused, repetitive, or tempted to implement future phases early.

## 2. Authoritative Priority Rule

This file owns context priority and classification.

```text
docs/engineering/ai-context-loading-standard.md
  = authoritative source of truth for AI context priority and loading classification
```

Individual file front matter such as:

```yaml
ai_load_priority: always
```

is only a **local hint**.

If local file metadata conflicts with this document, this document wins.

AI agents must not treat individual file front matter as more authoritative than the registry in this document.

## 3. Core Principle

Do not load the whole repository by default.

Use:

```text
small always-load core
+ task-specific context
+ current feature file
+ implementation-specific guardrails
```

Avoid:

```text
loading all feature files
loading all RFCs
loading archived notes
loading generated site output
loading old SDE-only documents for Phase 1 platform implementation
```

AI must load the smallest set of documents that can answer or implement the current task correctly.

## 4. Canonical Ownership

If files conflict, canonical ownership wins:

```text
AGENTS.md
  short master AI operating contract

docs/engineering/ai-context-loading-standard.md
  authoritative context classification and loading matrix

docs/engineering/go-coding-guardrails.md
  canonical Go implementation rules

docs/ai/AI_CONTEXT_GUIDE.md
  short pointer to canonical standards

docs/engineering/context-engineering-standard.md
  high-level context-engineering principle only

docs/engineering/ai-controlled-development.md
  operating model for AI-assisted development

.cursor/rules/*
  short Cursor enforcement rules

.kiro/steering/*
  short Kiro steering rules

docs/glossary.md
  canonical terminology

docs/decisions/DECISION_INDEX.md
  accepted decisions

docs/api/API_CONTRACT_PHASE1.md
  API behavior

docs/resource-specs/RESOURCE_MODEL_PHASE1.md
  resource model

current docs/features/FEATURE-xxxx file
  current feature scope
```

## 5. Source of Truth Hierarchy

When files conflict, use this order:

```text
1. AGENTS.md
2. docs/foundation/constitution.md
3. docs/decisions/DECISION_INDEX.md
4. docs/glossary.md
5. docs/features/FEATURE_SEQUENCE.md
6. docs/resource-specs/RESOURCE_MODEL_PHASE1.md
7. docs/api/API_CONTRACT_PHASE1.md
8. current feature file
9. relevant engineering or architecture standard
10. role-specific Kiro/Cursor/ChatGPT prompt files
11. reference RFCs and historical docs
```

Archived documents are not source of truth unless explicitly reactivated.

## 6. Context Categories

Sovrunn uses these context categories:

```text
always
go-implementation
feature
role-kiro
role-cursor
role-chatgpt
reference
archive
do-not-load
```

Category meanings:

| Category | Meaning |
|---|---|
| `always` | Load for every AI task unless the task is trivial and does not need repo context |
| `go-implementation` | Load for Go coding, tests, registry, API, validation, security, performance, or observability work |
| `feature` | Load only the current feature file |
| `role-kiro` | Load for Kiro planning and task breakdown |
| `role-cursor` | Load for Cursor code editing/refactoring/debugging |
| `role-chatgpt` | Load only when asking ChatGPT for that specific review, patch, debug, or implementation task |
| `reference` | Load only when directly relevant |
| `archive` | Historical only; not source of truth |
| `do-not-load` | Never load into AI reasoning context |

## 7. Authoritative Context File Registry

This table is authoritative.

If file-local metadata differs from this table, this table wins.

| File or Pattern | Category | Owner | Purpose |
|---|---|---|---|
| `AGENTS.md` | `always` | master AI contract | Short operating contract for all AI tools |
| `README.md` | `always` | repo identity | Sovrunn-first positioning and onboarding |
| `docs/foundation/constitution.md` | `always` | platform principles | Platform boundaries, principles, and non-goals |
| `docs/decisions/DECISION_INDEX.md` | `always` | accepted decisions | Decision registry and accepted architecture choices |
| `docs/glossary.md` | `always` | terminology | Canonical terms |
| `docs/features/FEATURE_SEQUENCE.md` | `always` | feature order | Phase 1 feature sequence |
| `docs/resource-specs/RESOURCE_MODEL_PHASE1.md` | `always` | resource model | Phase 1 resource model |
| `docs/api/API_CONTRACT_PHASE1.md` | `always` | API behavior | Phase 1 API contract |
| `docs/engineering/ai-context-loading-standard.md` | `always` | context loading | Authoritative context registry and loading matrix |
| `docs/engineering/go-coding-guardrails.md` | `go-implementation` | Go implementation | Go coding rules and guardrails |
| `docs/architecture/controller-reconciliation-model.md` | `go-implementation` | controller model | Reconciliation and lifecycle behavior |
| `docs/architecture/observability-and-audit-baseline.md` | `go-implementation` | observability/audit | Logs, request IDs, audit, and operation expectations |
| `docs/features/FEATURE-*.md` | `feature` | feature scope | Current feature scope, non-goals, and acceptance criteria |
| `.kiro/steering/product.md` | `role-kiro` | Kiro product steering | Product planning context |
| `.kiro/steering/architecture.md` | `role-kiro` | Kiro architecture steering | Architecture planning context |
| `.kiro/steering/engineering.md` | `role-kiro` | Kiro engineering steering | Engineering planning context |
| `.cursor/rules/sovrunn-constitution.mdc` | `role-cursor` | Cursor rule | Constitution enforcement |
| `.cursor/rules/sovrunn-architecture.mdc` | `role-cursor` | Cursor rule | Architecture enforcement |
| `.cursor/rules/sovrunn-go-style.mdc` | `role-cursor` | Cursor rule | Go coding enforcement |
| `.cursor/rules/sovrunn-testing.mdc` | `role-cursor` | Cursor rule | Test enforcement |
| `.cursor/rules/sovrunn-ai-boundaries.mdc` | `role-cursor` | Cursor rule | AI boundary enforcement |
| `docs/prompts/CHATGPT_REVIEW_PROMPT.md` | `role-chatgpt` | ChatGPT prompt | Review prompt template |
| `docs/prompts/CHATGPT_PATCH_PROMPT.md` | `role-chatgpt` | ChatGPT prompt | Patch prompt template |
| `docs/prompts/CHATGPT_DEBUG_PROMPT.md` | `role-chatgpt` | `ChatGPT prompt` | Debug prompt template |
| `docs/prompts/AI_IMPLEMENTATION_PROMPT_PHASE1.md` | `role-chatgpt` | implementation prompt | Phase 1 implementation prompt template |
| `docs/foundation/vision.md` | `reference` | product vision | High-level product direction |
| `docs/foundation/philosophy.md` | `reference` | product philosophy | Design and product philosophy |
| `docs/foundation/PHASE1_REVIEW_STATUS.md` | `reference` | review status | Phase 1 review notes |
| `docs/architecture/platform-core.md` | `reference` | platform architecture | Platform core architecture |
| `docs/architecture/organization-governance.md` | `reference` | governance architecture | Organization governance architecture |
| `docs/architecture/development-phases.md` | `role-kiro` | roadmap | Development phase planning |
| `docs/architecture/gitops-desired-state-model.md` | `reference` | GitOps architecture | Desired-state/GitOps model |
| `docs/engineering/ai-controlled-development.md` | `reference` | AI development model | AI-assisted development operating model |
| `docs/engineering/context-engineering-standard.md` | `reference` | context principle | High-level context engineering principle |
| `docs/ai/AI_CONTEXT_GUIDE.md` | `reference` | AI guide | Short pointer to canonical standards |
| `docs/ai/AI_DOC_AUTHORING_STANDARD.md` | `reference` | doc authoring | Documentation authoring standard |
| `docs/ai/AI_FEATURE_FACTORY.md` | `reference` | feature authoring | AI-assisted feature authoring workflow |
| `docs/rfc/*.md` | `reference` | RFCs | Architecture decisions and deeper design |
| `docs/archive/**` | `archive` | history | Historical/superseded notes only |
| `.git/**` | `do-not-load` | Git internals | Never load into AI reasoning |
| `site/**` | `do-not-load` | generated docs | Generated output |
| `__MACOSX/**` | `do-not-load` | macOS artifact | ZIP artifact |
| `*.DS_Store` | `do-not-load` | macOS artifact | Local metadata |
| `docs/.obsidian/**` | `do-not-load` | editor metadata | Local Obsidian metadata |
| `node_modules/**` | `do-not-load` | dependencies | External packages |
| `vendor/**` | `do-not-load` | dependencies | External packages |
| `bin/**` | `do-not-load` | build output | Compiled artifacts |
| `dist/**` | `do-not-load` | build output | Generated artifacts |
| `coverage/**` | `do-not-load` | test output | Coverage reports |

## 8. ALWAYS Loading Set

Load for every AI coding, planning, review, or debugging task:

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

Why these are always loaded:

```text
They prevent architecture drift.
They define canonical terms.
They define Phase 1 boundaries.
They define the Phase 1 resource model.
They define API behavior.
They define accepted decisions.
They stop AI from implementing future features early.
They define context-selection behavior.
```

Do not add files to ALWAYS casually. The always-load set must stay small.

## 9. GO_IMPLEMENTATION Loading Set

Load when the task involves Go code, API handlers, validation, registry, server lifecycle, tests, performance, observability, security, or runtime behavior:

```text
docs/engineering/go-coding-guardrails.md
docs/architecture/controller-reconciliation-model.md
docs/architecture/observability-and-audit-baseline.md
```

Use for:

```text
writing Go code
reviewing Go code
debugging Go tests
designing HTTP handlers
designing registry behavior
designing validation
reviewing latency or performance
reviewing security guardrails
reviewing logs, audit, or operation behavior
```

## 10. FEATURE Loading Rule

Load only the current feature file.

Examples:

```text
docs/features/FEATURE-0001-organization-resource-and-registry.md
docs/features/FEATURE-0002-organizationunit-resource.md
docs/features/FEATURE-0003-tenant-resource.md
docs/features/FEATURE-0004-project-resource.md
docs/features/FEATURE-0005-operation-resource.md
docs/features/FEATURE-0006-serviceclass-and-serviceplan.md
docs/features/FEATURE-0007-plugin-and-capability-registry.md
docs/features/FEATURE-0008-serviceinstance-and-servicebinding.md
docs/features/FEATURE-0009-api-server-health-readiness.md
docs/features/FEATURE-0010-basic-cli-api-demo-flow.md
```

Rule:

```text
Load FEATURE_SEQUENCE.md plus the current feature file.
Do not load all feature files during implementation.
```

Why:

```text
FEATURE_SEQUENCE.md gives the roadmap.
The current feature file gives the implementation boundary.
Future feature files can cause AI agents to implement scope too early.
```

## 11. ROLE-SPECIFIC Loading Rules

### Kiro Planning

Load:

```text
ALWAYS
role-kiro files
current FEATURE file, when planning a feature
```

Use for:

```text
creating feature specs
breaking features into tasks
planning implementation order
reviewing roadmap
checking architectural consistency
reviewing product scope
```

### Cursor Coding

Load:

```text
ALWAYS
role-cursor files
GO_IMPLEMENTATION, when coding Go
current FEATURE file
```

Cursor must not load every document in the repository by default.

### ChatGPT Tasks

Load only the prompt matching the task:

```text
CHATGPT_REVIEW_PROMPT.md
  for architecture, code, or doc review

CHATGPT_PATCH_PROMPT.md
  for exact file-level patches

CHATGPT_DEBUG_PROMPT.md
  for terminal output, test failures, runtime logs, and API debugging

AI_IMPLEMENTATION_PROMPT_PHASE1.md
  for full Phase 1 feature implementation prompts
```

Do not load all prompt files at once.

## 12. REFERENCE Loading Rule

Reference files are useful but should not be loaded by default.

Load reference files only when directly relevant to the task.

Examples:

```text
docs/foundation/vision.md
docs/foundation/philosophy.md
docs/architecture/platform-core.md
docs/architecture/organization-governance.md
docs/architecture/gitops-desired-state-model.md
docs/engineering/ai-controlled-development.md
docs/engineering/context-engineering-standard.md
docs/ai/AI_CONTEXT_GUIDE.md
docs/ai/AI_DOC_AUTHORING_STANDARD.md
docs/ai/AI_FEATURE_FACTORY.md
docs/rfc/*.md
```

Use when:

```text
reviewing high-level architecture
writing or editing RFCs
changing AI workflow
changing GitOps model
reviewing product philosophy
reviewing documentation standards
checking long-term design direction
```

## 13. ARCHIVE and DO_NOT_LOAD Rules

Archived files are historical reference only.

```text
docs/archive/
old SDE-only notes
temporary update reports
superseded design notes
```

Rules:

```text
Do not load archive files by default.
Do not treat archive files as source of truth.
Load archive files only when the user explicitly asks for history or comparison.
```

Never load these into AI reasoning context:

```text
.git/
site/
__MACOSX/
.DS_Store
docs/.DS_Store
docs/.obsidian/
node_modules/
vendor/
tmp/
dist/
bin/
coverage/
build output
compiled binaries
```

## 14. Loading Matrix by Task Type

### 14.1 Feature Implementation in Go

Load:

```text
ALWAYS
GO_IMPLEMENTATION
current FEATURE file
```

Optional:

```text
role-cursor files, when using Cursor
AI_IMPLEMENTATION_PROMPT_PHASE1.md, when asking ChatGPT for implementation guidance
```

Do not load:

```text
future feature files
all RFCs
archive
generated site
old SDE-only notes
```

### 14.2 FEATURE-0001 Implementation Example

Load:

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
docs/features/FEATURE-0001-organization-resource-and-registry.md
```

Do not load:

```text
FEATURE-0002 through FEATURE-0010
RFC-0020 ServiceOps Plugin Framework
SDE runtime docs
GitOps model
AI feature factory
archive notes
```

### 14.3 Architecture Review

Load:

```text
ALWAYS
docs/foundation/vision.md
docs/foundation/philosophy.md
docs/architecture/platform-core.md
docs/architecture/organization-governance.md
docs/architecture/development-phases.md
```

Optional:

```text
docs/architecture/controller-reconciliation-model.md
docs/architecture/gitops-desired-state-model.md
docs/rfc/*.md
docs/prompts/CHATGPT_REVIEW_PROMPT.md
```

### 14.4 API Design

Load:

```text
ALWAYS
current FEATURE file
```

For Go implementation, also load:

```text
GO_IMPLEMENTATION
```

Do not load unrelated RFCs or future features.

### 14.5 Debugging Terminal Errors

Load:

```text
ALWAYS
docs/prompts/CHATGPT_DEBUG_PROMPT.md
current FEATURE file
terminal output
changed files
```

If the error is Go-related, also load:

```text
docs/engineering/go-coding-guardrails.md
```

Do not load broad architecture unless the failure is architectural.

### 14.6 Documentation Edit

Load:

```text
ALWAYS
target document
related source-of-truth document
```

If editing architecture docs, add:

```text
docs/foundation/vision.md
docs/foundation/philosophy.md
docs/architecture/platform-core.md
```

If editing AI docs, add:

```text
docs/engineering/ai-controlled-development.md
docs/engineering/context-engineering-standard.md
```

### 14.7 RFC Authoring

Load:

```text
ALWAYS
docs/ai/AI_DOC_AUTHORING_STANDARD.md
docs/foundation/vision.md
docs/foundation/philosophy.md
related architecture docs
related existing RFCs
```

Do not load unrelated feature implementation files unless they are directly affected.

## 15. Rules for Adding New Files

Individual file front matter is optional local metadata, not canonical priority.

If used, recommended fields are:

```yaml
---
doc_type: engineering_standard
title: Example Title
status: draft
phase: 1
ai_summary: One-line summary for AI context selection.
---
```

Avoid overusing `ai_load_priority` in individual files.

If `ai_load_priority` is present, it is only a hint. The authoritative registry in this file wins.

When adding a new file, update the Authoritative Context File Registry in this document if the file should be visible to AI tools.

## 16. Rules for Promoting Files to ALWAYS

A file may be added to ALWAYS only if it satisfies all of these:

```text
It is required for almost every AI task.
It prevents major architecture drift.
It is stable enough to be source of truth.
It is not too long.
It does not duplicate another always-load file.
It does not encourage implementation of future features.
```

Do not add files to ALWAYS just because they are useful.

Useful files usually belong in REFERENCE.

## 17. Token Optimization Rules

AI agents should reduce token load by following these rules:

```text
Load summaries first when available.
Load current feature file instead of all features.
Load the API contract instead of every API example.
Load the resource model instead of repeating resource definitions.
Load reference files only when directly relevant.
Avoid archive files unless requested.
Avoid generated docs site output.
Avoid binary/build files.
```

Do not duplicate long rules across many files. Use pointers to canonical documents.

## 18. Confusion Avoidance Rules

Avoid creating multiple files that say the same thing differently.

If two files conflict, update the non-canonical file to point to the canonical one.

Tool-specific files should be short wrappers, not duplicate source-of-truth documents.

## 19. Clean AI Review ZIP Rule

When creating a ZIP for AI review, exclude non-context files:

```bash
zip -r sovrunn-clean.zip . \\
  -x ".git/*" \\
  -x "site/*" \\
  -x "__MACOSX/*" \\
  -x "*.DS_Store" \\
  -x "docs/.obsidian/*" \\
  -x "node_modules/*" \\
  -x "vendor/*" \\
  -x "bin/*" \\
  -x "dist/*" \\
  -x "coverage/*"
```

Do not delete `.git/` from the actual working repository.

## 20. Final Recommendation

Keep ALWAYS small.

The default Go coding set should usually be:

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
current FEATURE file
```

This gives AI agents enough context to implement correctly without creating architecture drift, duplication, or unnecessary token load.
