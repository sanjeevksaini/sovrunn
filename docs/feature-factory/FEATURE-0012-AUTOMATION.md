# FEATURE-0012 automated implementation flow

`feature-0012-flow.py` is a feature-specific, resumable orchestrator for the
approved FEATURE-0012 task graph. It replaces the generic integer-only
`task-flow.sh` for this feature.

## Safety model

The flow:

- reads exact task IDs and dependency waves from `tasks.md`, including IDs such
  as `6.5a` and `7a.1`;
- skips completed tasks and resumes at the first pending dependency-ready task;
- requires the approved feature branch, `APPROVED_FOR_CURSOR`, a clean working
  tree, Docker, and Cursor CLI;
- runs one Cursor task at a time;
- rejects changes outside a task-family allowlist;
- rejects changes to approved requirements, design, and architecture;
- rejects dependency changes after Task 1.1;
- allows only the current task checkbox to change, then rolls up the parent
  checkbox when all children are complete;
- runs Go 1.22 Docker verification and repository guardrails before every task
  commit;
- commits one validated task at a time using an exact staged path set;
- executes checkpoints 3, 7, and 13 without Cursor;
- reserves checkpoint 18 for one explicit final human-reviewed invocation;
- stops fail-closed on any Cursor, scope, validation, checkpoint, or commit
  failure;
- never writes a human implementation approval token;
- finishes unattended work in `PENDING_HUMAN_REVIEW`;
- after the manual final checkpoint, remains gated in `PENDING_HUMAN_REVIEW`.

Machine gate logs are written under
`.automation/logs/FEATURE-0012/flow/` and remain ignored runtime evidence.

## Before starting

Finish and commit any manually running task first. The working tree must be
clean. In particular, when Task 1.2 is being completed in Cursor Studio, commit
Task 1.2 before starting this flow. The script will then resume at Task 1.3.

Commit the automation changes separately from a feature leaf task, for example:

```bash
chmod +x scripts/feature-0012-flow.py
python3 scripts/feature-0012-flow.py --self-test
git diff --check
git add \
  Makefile \
  scripts/feature-0012-flow.py \
  scripts/task-flow.sh \
  scripts/verify.sh \
  scripts/final-verify.sh \
  docs/feature-factory/FEATURE-0012-AUTOMATION.md
git commit -m "chore(feature-factory): add safe FEATURE-0012 flow"
```

## Inspect the resume plan

```bash
make ff-feature-0012-plan
```

## Canary run

Run exactly one pending step and stop:

```bash
CONFIRM_FEATURE_0012_AUTORUN=YES \
MAX_STEPS=1 \
make ff-feature-0012-run
```

Review the resulting commit, then resume with the same command without
`MAX_STEPS`.

## Full resumable run

```bash
CONFIRM_FEATURE_0012_AUTORUN=YES \
make ff-feature-0012-run
```

The same command resumes after a machine failure once the failure has been
reviewed, corrected, and committed or reverted to a clean tree.

Optional controls:

```bash
CONFIRM_FEATURE_0012_AUTORUN=YES \
START_TASK=8.1 \
STOP_AFTER=9.5 \
make ff-feature-0012-run
```

`START_TASK` is accepted only when all earlier dependency waves are complete.
`STOP_AFTER` and `MAX_STEPS` are safe pause controls, not task-skipping tools.

## Failure handling

On failure the flow:

1. stops immediately;
2. records the failed task and error in feature state and ignored machine logs;
3. leaves the task working tree untouched for diagnosis;
4. does not commit the failed task;
5. does not continue to later tasks.

Inspect with:

```bash
git status --short
git diff --check
make ff-state FEATURE=FEATURE-0012
cat .automation/logs/FEATURE-0012/flow/<task-id>.json
```

Do not blindly rerun while the tree is dirty. Review or revert the failed task,
return to a clean tree, and then resume.

## Final manual review and checkpoint

The unattended run stops after Task 17.3 and sets:

```text
status: PENDING_HUMAN_REVIEW
human_gate_required: true
automation_flow_status: awaiting_final_checkpoint_review
```

Review the complete implementation and residual-risk evidence. Then run the
final checkpoint explicitly:

```bash
CONFIRM_FEATURE_0012_FINAL_CHECKPOINT=YES \
make ff-feature-0012-final-checkpoint
```

Checkpoint 18 runs the full test, race, vet, guardrail, and feature-gate
barrier, marks Task 18 complete, and commits the checkpoint. It then leaves the
feature in:

```text
status: PENDING_HUMAN_REVIEW
human_gate_required: true
automation_flow_status: final_checkpoint_passed_pending_approval
```

A human records the final implementation approval separately. The orchestrator
never writes that approval token or accepts residual risk.
