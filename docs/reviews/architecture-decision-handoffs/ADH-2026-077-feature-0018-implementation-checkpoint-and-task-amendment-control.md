# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-077
- Date: 2026-09-07
- Source discussion: ChatGPT Project / Codex
- Related feature: FEATURE-0018 — Governance, IAM, Approval and Exception Foundation
- Related phase: Phase 2R
- Author: ChatGPT / Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

FEATURE-0018 implementation checkpoint and task-amendment control

## Summary

FEATURE-0018 has committed Cursor Tasks 1–3, but the Feature Factory permits a
later Kiro requirements or design prompt to overwrite the global feature stage.
It also treats missing implementation paths or Go-level mechanics as reasons to
reopen the whole design. This replacement corrects the delivery process without
changing IAM product behavior: it freezes the approved semantic design baseline,
introduces an implementation checkpoint after the first committed Cursor task,
allows bounded amendments only for uncommitted tasks, and moves Go-level API and
checker detail to implementation evidence rather than recurring LLM design review.

## Classification

- Replacement

## Existing approved baseline

`ARCH-2026.08-PHASE2R-CANONICAL`, DEC-0060, and ADH-2026-070 through
ADH-2026-074 remain the FEATURE-0018 product-semantic authority. In particular,
the Task-1–3 checkpoint was built against the approved design blob:

```text
commit: 7d4447d7c1d07c37da4af738dfb399adaa9f4ebc
path:   .kiro/specs/governance-iam-approval-exception-foundation/design.md
sha256: 4eb84a0002d8e58de466806b2dcd18085608a0a62d97cc78684b7e98210d92a8
```

That blob includes the approved ADH-2026-073
`exceptionproposal.decided` correction. The approved requirements hash is:

```text
aed8f347ea154c30292cdd0d8f43841fad6d419b8cf789becacc8b57d9c2a9e4
```

Tasks 1–3 are committed. The feature state records `last_committed_task: "3"`
and `current_task: "4"`. The later unapproved working `design.md` hash
`bc516a50f2705b8ec168b8ed21476ddd88667b37de8e0654933329bd67915219`
is not an additional product-semantic authority.

ADH-2026-075 and ADH-2026-076 correctly sought to reduce duplicated mechanics
prose, but their post-checkpoint application permitted a Design Mechanics Contract
to become an LLM-review target for exact Go API shapes. That is not an appropriate
delivery authority after Cursor implementation has begun.

Relevant baseline references:

- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/decisions/DECISION_INDEX.md`
- `docs/decisions/DEC-0060-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-073-feature-0018-exceptionproposal-terminal-decision-audit-event.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-075-feature-0018-design-normalization-and-contract-check.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-076-feature-0018-semantic-review-authority-realignment.md`

## Decision or proposed decision

### 1. Freeze the FEATURE-0018 semantic design baseline

The approved design blob at commit `7d4447d` and SHA-256 `4eb84a...92a8` is
the frozen FEATURE-0018 semantic design baseline for implementation. It remains
paired with the approved requirements hash above and the approved product-semantic
ADHs through ADH-2026-074.

Kiro must restore or derive the implementation-facing semantic design from that
baseline only. It must not infer a new IAM model from the later working design
draft. Before replacement, Kiro must record the working draft hash and a concise
classified delta ledger; it must preserve no duplicate mechanics prose solely for
historical recovery.

The frozen baseline may be reopened only for a human-approved handoff that changes
one or more of: product resource, action, route, writer, public error, event
taxonomy, authorization/approval/review/exception semantics, dependency ownership,
or the Phase 2R boundary.

### 2. Establish an implementation checkpoint

Once a feature has a committed Cursor task, its Feature Factory state enters
`implementation_checkpoint`. FEATURE-0018 enters this state immediately with:

```text
lastCommittedTask = 3
currentTask       = 4
frozenDesignSha256 = 4eb84a...92a8
frozenRequirementsSha256 = aed8f...9e4
```

While an implementation checkpoint is active:

- normal `requirements` and `design` Kiro stages and their LLM review routes must
  fail closed;
- `kiro-stage.sh`, `review-and-route-stage.sh`, `approve-stage.sh`, and Cursor
  prompt generation must preserve the checkpoint rather than overwrite
  `current_stage`;
- a caller must choose either a bounded task-plan amendment or a separately
  approved semantic architecture handoff before leaving the checkpoint; and
- neither an unapproved `design.md` edit nor an LLM design finding may invalidate
  completed tasks by itself.

An explicit architecture-reentry command may be provided only when it requires an
approved handoff identifier and records the semantic impact on committed tasks.

### 3. Define the only permitted repair path for Task 4–10

An `implementation_task_amendment` may change only an uncommitted task numbered
strictly greater than `lastCommittedTask`. It may correct its writable paths,
package/import ordering, test locations, verification commands, dependency edges,
or bounded implementation mechanics already required by the frozen design and
requirements. It must not alter product semantics.

Before the amended task plan can be approved for Cursor, a deterministic task
preflight must prove that each affected task:

- has writable paths for every production package and test file it must create or
  modify;
- consumes only approved earlier-wave outputs or paths included in the same
  executable task;
- has no forward dependency on a later executable task;
- preserves the completed Task 1–3 paths and acceptance evidence; and
- invokes the applicable implementation checks.

The Task 4 amendment must specifically retain its supporting-value model paths and
the exact dependency/preflight evidence that prevented its earlier Cursor block.
It is a task-only correction. It must not trigger a design review.

The amendment requires a fresh task-plan and executable-plan approval bound to the
amended `tasks.md` hash, then returns directly to `cursor` at Task 4.

### 4. Put Go-level detail in implementation proof

Exact Go interfaces, method selectors, transaction helper names, `StateEditor`
method surfaces, lock implementation details, factory shapes, and static-checker
fixtures are implementation mechanics. They are defined and proven by the code,
Go compiler, focused unit tests, race tests, and the FEATURE-0018 static
architecture checker during Tasks 4–10.

They must not be used to reopen `design.md` when the frozen product semantics,
requirements, and task-plan boundary already cover the behavior. A deterministic
mechanics checker may remain as an implementation check, but it is not an LLM
design-review authority and must not prescribe a second product-design source.

### 5. Replace the post-checkpoint portions of ADH-2026-075 and ADH-2026-076

For FEATURE-0018 after the Task-1 commit, this handoff replaces:

- ADH-2026-075 section 2's requirement that the Design Mechanics Contract be the
  authoritative complete Go API/method-surface representation before further
  design review; and
- ADH-2026-076 sections 1–4 insofar as they permit a design review to raise a
  Go-level mechanics completeness finding after an implementation checkpoint.

ADH-2026-075 and ADH-2026-076 remain valid as historical normalization and
review-context evidence. They do not authorize another FEATURE-0018 design
iteration after this checkpoint.

### 6. Final review happens after implementation, once

After Task 10 and all implementation checks pass, run one independent final review
against the frozen semantic baseline, approved requirements, approved task-plan
amendments, implementation diff, and feature-gate evidence. That review may assess
actual product-semantic divergence and security defects. It must not request a
design rewrite for an alternative Go API representation.

## Rationale

The repeated design cycles did not discover new IAM architecture. They exposed a
missing distinction between semantic architecture changes and implementation-plan
repairs. Task 4's missing model writable paths should have been caught by a
pre-Cursor package/path preflight. Instead, the Feature Factory allowed global
stage rewrites, and ADH-2026-075 caused design review to inspect compiler-level
details repeatedly. Freezing the known approved semantic baseline, preserving the
Cursor checkpoint, and using deterministic implementation evidence restores a
short, test-gated delivery path.

## Reuse-before-build assessment

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Extend the existing Feature Factory hash, state, approval, task-context and
    Cursor path-allowlist controls.
  - Reuse Go compiler, `go test`, `go test -race`, existing Feature Factory tests,
    and the repository static-checker pattern for implementation mechanics.
- Sovrunn-owned responsibility summary:
  - FEATURE-0018 owns only its checkpoint metadata, bounded task-plan amendment,
    deterministic preflight and implementation evidence.
- Non-goals summary:
  - No IAM product change; no resource, route, writer, action, error, event,
    dependency, provider-native effect, persistence, external I/O, or Phase 2R
    change.

## Phase impact

- Current phase allowed? Yes
- If not current phase, target phase: Not applicable
- Current phase boundary impact:
  - Feature Factory governance and deterministic validation only. No external I/O,
    durable persistence, provider effect, or execution-target effect is introduced.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting decisions, if any:
  - ADH-2026-075 and ADH-2026-076 contain post-checkpoint review behavior that
    this handoff expressly replaces. Their product-semantic protections remain.
- Resolution required:
  - Human approval of this handoff, then Kiro application. No DEC, RFC, ACR or
    architecture-baseline update is required.

## Required action

- Update Feature Factory checkpoint state/control/schema and deterministic tests
- Update Kiro-stage, review-route, approval and Cursor-prompt guards
- Create/update bounded task-amendment and preflight commands
- Restore/freeze the approved FEATURE-0018 semantic design baseline and record a
  classified working-draft delta ledger
- Update FEATURE-0018 tasks.md only through the checkpoint-aware Task 4 amendment
- Update feature gate/checks
- Do not update requirements.md
- Do not modify completed Task 1–3 implementation
- Do not generate or review a new FEATURE-0018 design before Task 10

## Impacted files

- `docs/reviews/architecture-decision-handoffs/ADH-2026-077-feature-0018-implementation-checkpoint-and-task-amendment-control.md`
- `.automation/state/FEATURE-0018.json`
- `.automation/features/FEATURE-0018.control.json`
- `.automation/schemas/feature-control.schema.json` only if needed for checkpoint
  metadata validation
- `.kiro/specs/governance-iam-approval-exception-foundation/design.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/tasks.md`
- `scripts/kiro-stage.sh`
- `scripts/review-and-route-stage.sh`
- `scripts/approve-stage.sh`
- `scripts/generic-cursor-prompt.py`
- `scripts/cursor-task.sh`
- `scripts/feature-control.py` and new bounded preflight/amendment scripts only as
  required
- `scripts/feature-gate.sh` and `Makefile` only as needed for the new checks
- `tests/feature_factory/` for deterministic checkpoint, amendment and preflight
  regression tests
- `docs/reviews/architecture-readiness/FEATURE-0018-implementation-checkpoint-delta-ledger.md`

## Impacted features

- FEATURE-0018: delivery-process correction and Task 4–10 execution control only.
- Other features: no immediate behavior change. The generic Feature Factory
  checkpoint mechanism may be reusable after FEATURE-0018 proof succeeds, but is
  not adopted for another feature through this handoff.

## Acceptance criteria for Kiro update

- [ ] Validate this handoff against the Architecture Operating System.
- [ ] Preserve the frozen approved design and requirements hashes in
  FEATURE-0018 checkpoint metadata.
- [ ] Preserve all Task 1–3 commits and reject an amendment that changes their
  paths, task blocks or accepted evidence.
- [ ] Make an active checkpoint reject ordinary Kiro requirements/design stages,
  design review routing and stage-state overwrite.
- [ ] Make Cursor reject execution when the frozen semantic design or requirements
  hash changes without approved semantic architecture reentry.
- [ ] Provide a deterministic, task-only amendment path for Task 4–10.
- [ ] Fail the amendment preflight on a missing writable path, later-task import,
  invalid wave dependency or modification of a committed task.
- [ ] Require fresh hash-bound task/executable-plan approval after a permitted
  amendment, then restore `current_stage: cursor` and `current_task: 4`.
- [ ] Keep Go-level mechanics evidence in implementation checks and prohibit an
  LLM design revision for that category after the checkpoint.
- [ ] Run focused Feature Factory tests, design semantic guardrail, VS-000 check,
  Phase 2R drift check and the applicable feature-gate checks.

## Explicit instructions to Kiro

- Do not infer or change FEATURE-0018 IAM behavior.
- Do not change requirements, DEC-0060, ADH-2026-070 through ADH-2026-074,
  canonical architecture, resources, actions, routes, writers, errors, event
  taxonomy, dependency ownership, or Phase 2R scope.
- Do not rewrite the frozen semantic design from the current working draft.
- Do not use a design review to correct an implementation path, Go signature,
  static-checker selector, compiler error, test failure or task dependency.
- Stop with `ARCHITECTURE_DECISION_REQUIRED` if a Task 4–10 repair changes product
  semantics or affects the committed Task 1–3 behavioral contract.
- Do not modify `internal/govaccess/` under this handoff. Cursor resumes code work
  only after the amended task plan is separately approved.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-09-07
- Notes: Approved to stop post-checkpoint FEATURE-0018 design iteration and resume
  delivery through a frozen semantic baseline, deterministic task amendment and
  implementation evidence.
