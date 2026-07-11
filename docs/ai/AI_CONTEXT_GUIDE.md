---
doc_type: ai_guide
title: AI Context Guide
status: draft
phase: 0
ai_load_priority: always
ai_summary: Tells AI coding agents which Sovrunn files to read and how to use them before implementing features.
---

# AI Context Guide

## 1. Purpose

This file tells AI agents how to work on Sovrunn.

AI must not invent architecture.

AI must load the correct context, follow accepted decisions, and keep generated code aligned with the platform constitution.

## 2. Required Always-Load Files

For most Sovrunn development tasks, load:

```text
docs/foundation/vision.md
docs/foundation/philosophy.md
docs/foundation/constitution.md
docs/decisions/DECISION_INDEX.md
docs/glossary.md
docs/ai/AI_DOC_AUTHORING_STANDARD.md
docs/ai/AI_FEATURE_FACTORY.md
```

## 3. Phase-Specific Files

For Phase 1 work, also load:

```text
docs/architecture/development-phases.md
docs/architecture/platform-core.md
docs/architecture/organization-governance.md
docs/engineering/go-style.md
docs/engineering/package-layout.md
docs/engineering/testing-standard.md
docs/engineering/security-checklist.md
docs/rfc/RFC-0012-organization-tenant-governance-model.md
docs/rfc/RFC-0020-serviceops-plugin-framework.md
```

## 4. AI Operating Rules

AI must:

- use canonical terms from `glossary.md`,
- follow `constitution.md`,
- follow `DECISION_INDEX.md`,
- follow relevant RFCs,
- follow engineering standards,
- keep scope within current phase,
- ask for missing context when needed,
- state assumptions explicitly,
- generate tests with code,
- update docs when behavior changes.

AI must not:

- invent new architecture,
- create future-phase features unless requested,
- bypass Operation model,
- bypass policy/audit concepts,
- introduce uncontrolled AI actions,
- expose secrets,
- implement custom infrastructure where open-source substrate is intended,
- use synonyms for core terms.

## 5. Feature Task Context Bundle

Every feature task should define:

```text
Always Load
Phase Load
Feature Load
Engineering Load
Existing Code Paths
Expected Output Files
Acceptance Criteria
```

## 6. AI Coding Sequence

Use this sequence:

```text
1. Read required context files.
2. Summarize applicable decisions.
3. Confirm in-scope and out-of-scope.
4. Identify resource/API changes.
5. Identify code paths.
6. Generate or update tests.
7. Generate implementation.
8. Run tests where possible.
9. Update docs.
10. Report changed files and acceptance status.
```

## 7. Constitutional Checks

Before finishing a task, AI must verify:

- organization-first governance respected,
- tenant isolation respected,
- policy inheritance not weakened,
- Operation created for lifecycle changes,
- audit fields preserved,
- no secrets exposed,
- no AI in latency-sensitive hot path,
- no future-phase scope creep,
- failure modes explicit.

## 8. Missing Context Rule

If a required RFC, spec, or standard is missing, AI should not silently invent it.

AI should respond with:

```text
Missing required context: <file>.
Assumption if proceeding: <assumption>.
Risk: <risk>.
```

## 9. Output Format for AI Work

AI should return:

```text
Summary
Files changed
Tests added
Acceptance criteria status
Open questions
Risks
Next recommended task
```
