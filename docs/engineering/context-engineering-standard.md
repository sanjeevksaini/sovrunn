---
doc_type: engineering_standard
title: Context Engineering Standard
status: draft
phase: 1
ai_load_priority: always
ai_summary: Defines how Sovrunn coding agents should load task-specific context to reduce architecture drift and avoid token pollution.
---

# Context Engineering Standard

## 1. Purpose

Sovrunn has a large architecture surface. AI coding tools must receive the right context, not all context.

The principle is:

```text
Load the minimum complete context required for the task.
Do not load unrelated future architecture.
```

## 2. Context Tiers

### Tier 1: Always Load

```text
docs/foundation/constitution.md
docs/decisions/DECISION_INDEX.md
docs/glossary.md
docs/ai/AI_CONTEXT_GUIDE.md
docs/engineering/ai-controlled-development.md
docs/engineering/context-engineering-standard.md
```

### Tier 2: Phase 1 Load

```text
docs/features/FEATURE_SEQUENCE.md
docs/resource-specs/RESOURCE_MODEL_PHASE1.md
docs/api/API_CONTRACT_PHASE1.md
docs/architecture/platform-core.md
docs/architecture/organization-governance.md
docs/architecture/controller-reconciliation-model.md
docs/architecture/observability-and-audit-baseline.md
docs/architecture/gitops-desired-state-model.md
docs/engineering/go-style.md
docs/engineering/package-layout.md
docs/engineering/testing-standard.md
docs/engineering/security-checklist.md
```

### Tier 3: Feature Load

Load exactly one feature file at a time:

```text
docs/features/FEATURE-0001-organization-resource-and-registry.md
docs/features/FEATURE-0002-organizationunit-resource.md
docs/features/FEATURE-0003-tenant-resource.md
...
```

Do not load all feature files unless the task is cross-feature planning.

## 3. Do Not Load by Default

Do not load these during Phase 1 platform core coding unless explicitly requested:

```text
SDE transformation RFCs
datastore migration architecture
AI agent execution design
OPA/Gatekeeper implementation
OpenTelemetry collector deployment
Kubernetes CRD generation
Argo CD or Flux sync-controller implementation
marketplace design
billing engine design
UI portal design
multi-cluster federation implementation
```

## 4. Prompt Pattern

Every AI coding prompt should say:

```text
Implement FEATURE-xxxx only.
Use the loaded docs as source of truth.
Do not invent resource kinds.
Do not expand scope.
Ask before changing architecture.
Return code changes, tests, and acceptance validation.
```

## 5. Context Refresh Rule

Before each new feature:

```text
clear previous feature-specific context
reload Tier 1 and Tier 2
load only the new feature file
```

## 6. Conflict Rule

If documents conflict, priority order is:

```text
constitution.md
DECISION_INDEX.md
glossary.md
current feature file
RESOURCE_MODEL_PHASE1.md
API_CONTRACT_PHASE1.md
engineering standards
older notes
```

If conflict remains, stop and ask for human decision.

## 7. Output Discipline

AI output should include:

```text
files changed
why each file changed
tests added
commands to run
acceptance criteria satisfied
known limitations
next feature boundary
```

## 8. Final Principle

Good AI coding comes from good context boundaries.
