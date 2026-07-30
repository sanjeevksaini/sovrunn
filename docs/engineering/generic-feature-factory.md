---
doc_type: engineering_standard
title: Generic Manifest-Controlled Feature Factory
status: approved
phase: 2
ai_load_priority: important
ai_summary: Defines the manifest-controlled, three-human-gate workflow used for FEATURE-0015 and later features.
---

# Generic Manifest-Controlled Feature Factory

## Purpose

This workflow removes feature-specific prompt and orchestration code without
weakening architecture, semantic review, implementation boundaries, tests, or
founder approval. It applies when `.automation/features/<FEATURE>.control.json`
exists. Completed legacy features remain on their historical workflows.

The JSON Schema is `.automation/schemas/feature-control.schema.json`; the
starter is `.automation/templates/feature-control.template.json`.

## Three human gates

```text
Architecture boundary and handoff
  -> founder APPROVED_ARCHITECTURE
  -> Kiro requirements + independent machine review
  -> Kiro design + independent machine review
  -> founder APPROVED_EXECUTABLE_PLAN
  -> Kiro tasks + independent machine review
  -> Cursor vertical task batches + deterministic verification
  -> founder APPROVED_FOR_MERGE
  -> PR and merge
  -> prepared closeout edits; manual commit and push
```

Requirements, design, and tasks remain separate artifacts. The requirements
and design approvals between human gates are independent machine reviews, not
founder approvals. Any semantic conflict routes back to the stage that owns it.

## Control manifest ownership

The manifest owns execution control only:

- feature identity and tier;
- owned resources;
- relevant previous-feature dependencies and their exact context;
- excluded adjacent features and forbidden concepts;
- stage-specific context allowlists and byte budgets;
- Kiro writable outputs;
- Cursor Go context, verification commands, and COMPLETE receipt;
- batch size and human gates;
- centralized model roles;
- closeout targets and the mandatory prepare-only policy.

It does not own product or architecture semantics. Those remain in the approved
feature architecture and handoff package.

The architecture human gate records the SHA-256 digest of the manifest. Every
later Kiro, Cursor, PR, and closeout entry point fails closed if the manifest
changes. An intentional manifest change therefore requires a fresh architecture
approval; execution policy cannot drift silently during a feature.

## Context and token rules

Every stage receives the smallest sufficient set:

```text
small canonical core
+ current feature architecture and handoffs
+ exact previous-feature contracts consumed by this stage
+ stage-specific implementation or review evidence
```

The resolver deduplicates paths, verifies SHA-256 hashes, calculates bytes and
estimated tokens, and fails when the configured stage budget is exceeded. A
previous feature is context only where its relationship and files are explicit;
its contracts remain owned by that previous feature. Future feature specs are
not loaded merely for awareness.

Kiro prompts contain no FEATURE-0013/0014 special-case prose for controlled
future features. Cursor receives only the requested task block rather than the
whole task file, plus exact context and writable paths.

## Approved model routing

| Work | Primary route | Reasoning |
|---|---|---|
| Architecture and handoff | Human-guided architecture workspace | High; no autonomous gate transition |
| Kiro requirements | `claude-opus-4.8` | High |
| Kiro design | `claude-opus-4.8` | High |
| Kiro tasks | `claude-sonnet-4.5` | Medium |
| Independent semantic review | `gpt-5.5-thinking` | High for requirements/design/implementation; medium for tasks |
| Cursor implementation | `cursor-grok-4.5-high-fast` | High |
| Mechanical normalization and metadata | Deterministic scripts first | No model unless a semantic stop is raised |

Fallback order remains centralized in `.automation/model-policy.json`; feature
manifests name roles and cannot silently substitute models.

## Kiro anti-wandering controls

For requirements, design, and tasks:

- exactly one stage output is writable;
- later-stage artifacts and source code are forbidden;
- context is hash-bound to the rendered prompt;
- inherited contracts must be referenced, not copied or redefined;
- adjacent-feature concepts remain excluded;
- unresolved semantics stop with one approved stop condition;
- Kiro must produce exactly one `STAGE_STATUS: COMPLETE` receipt;
- changed paths are validated after generation.

## Cursor anti-wandering and Go controls

Cursor receives the canonical Go coding, version, and observability standards,
the approved spec package, relevant previous-feature implementation context,
one task block, and its exact writable paths. It must produce exactly one
`TASK_STATUS: COMPLETE` receipt. Context hashes, approved task digest, changed
paths, forbidden concepts, tests, vet, race checks where configured, and build
remain deterministic gates.

## Semantic delta

Each independent stage review stores the reviewed document. A later revision
gets a compact delta containing changed stable IDs, headings, normative lines,
and hashes. Only exact equality or trailing-whitespace/final-newline differences
can bypass semantic review. Every other change is
`SEMANTIC_REVIEW_REQUIRED`; scripts never claim semantic equivalence from
keywords alone.

## Closeout safety

`feature-closeout.py` verifies the merged PR through `gh`, confirms branch
identity and merge ancestry, and prepares the configured metadata edits. It
never commits or pushes. The founder reviews the diff and runs the printed Git
commands manually.

## Commands

```bash
make ff-feature-control FEATURE=FEATURE-0015

make ff-approve-human-gate \
  FEATURE=FEATURE-0015 GATE=architecture APPROVED_BY="Sanjeev Kumar"

make ff-spec-flow FEATURE=FEATURE-0015

make ff-approve-human-gate \
  FEATURE=FEATURE-0015 GATE=executable_plan APPROVED_BY="Sanjeev Kumar"

make ff-spec-flow FEATURE=FEATURE-0015

START_TASK=1 make ff-controlled-run FEATURE=FEATURE-0015

make ff-approve-human-gate \
  FEATURE=FEATURE-0015 GATE=final APPROVED_BY="Sanjeev Kumar"

make ff-pr FEATURE=FEATURE-0015

make ff-closeout FEATURE=FEATURE-0015 PR=<number>
APPLY=1 make ff-closeout FEATURE=FEATURE-0015 PR=<number>
```

The final two commands only inspect and prepare closeout. Commit and push remain
manual.
