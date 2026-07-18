# Sovrunn Feature Factory Pack

This pack gives Sovrunn a repeatable **Feature Factory** workflow.

It starts with **Level 1 automation** and is **Level 2 ready**:

- Level 1: generate prompts, manage feature state, run verification, guardrails, commits, and PR creation.
- Level 2-ready: script boundaries are prepared for Kiro CLI and Cursor CLI calls once local CLI authentication is confirmed.

The goal is to remove the user as the copy-paste bottleneck while keeping safe human approval gates.

## Install

From your Sovrunn repo root:

```bash
unzip sovrunn_feature_factory_pack.zip -d /tmp/sovrunn-feature-factory-pack
cp -R /tmp/sovrunn-feature-factory-pack/* .
chmod +x scripts/*.sh scripts/*.py
```

If you already have a `Makefile`, do not overwrite it blindly. Append the contents of `Makefile.feature-factory` to your repo `Makefile`.

Also append `.automation/templates/gitignore-feature-factory.snippet` to your repo `.gitignore`.

## Standard workflow

### 1. Start a feature

```bash
make ff-start FEATURE=FEATURE-0007 SLUG=plugin-capability-registry TITLE="Plugin and Capability Registry"
```

This creates:

```text
.automation/state/FEATURE-0007.json
.automation/features/FEATURE-0007.yaml
.kiro/specs/plugin-capability-registry/
docs/generated-prompts/FEATURE-0007/
```

It also creates the feature branch from `phase1-foundation` by default.

### 2. Generate Kiro prompts

```bash
make ff-prompt-requirements FEATURE=FEATURE-0007
make ff-prompt-design FEATURE=FEATURE-0007
make ff-prompt-tasks FEATURE=FEATURE-0007
```

Default behavior now runs Kiro CLI headlessly via `scripts/kiro-stage.sh`. Set `FEATURE_FACTORY_KIRO_MODE=prompt` only when you intentionally want to copy the prompt into Kiro manually.

### 3. Review each spec stage

```bash
make ff-review STAGE=requirements FEATURE=FEATURE-0007
make ff-review STAGE=design FEATURE=FEATURE-0007
make ff-review STAGE=tasks FEATURE=FEATURE-0007
```

This creates a review prompt that can be pasted into ChatGPT or a future reviewer agent.

### 4. Commit spec files

```bash
make ff-commit-spec FEATURE=FEATURE-0007
```

### 5. Generate Cursor task prompt

```bash
make ff-cursor-task FEATURE=FEATURE-0007 TASK=1.1
```

For Level 1, copy this prompt into Cursor.
Default behavior now runs Cursor CLI headlessly via `scripts/cursor-task.sh`. Set `FEATURE_FACTORY_CURSOR_MODE=prompt` only when you intentionally want to copy the prompt into Cursor manually.

### 6. Verify and commit task

```bash
make ff-verify
make ff-guardrails FEATURE=FEATURE-0007
make ff-commit-task FEATURE=FEATURE-0007 TASK=1.1 MESSAGE="feat(...): ..."
```

### 7. Final verification and PR

```bash
make ff-final FEATURE=FEATURE-0007
make ff-pr FEATURE=FEATURE-0007
```

## Human gates retained

The Feature Factory intentionally keeps these manual gates:

1. Requirements approval
2. Design approval
3. Tasks approval
4. PR merge approval

Everything between those gates can be automated safely.

## Safety principles

- Never deploy automatically.
- Never merge automatically.
- Never delete branches before merge.
- Never skip verification.
- Never continue after verification failure.
- Never commit build artifacts such as `sovrunn-api` or `bin/`.
- Never allow `internal/api` to import `internal/server`.
- Never leave `TODO(FEATURE-XXXX)` at final verification.


## v2 additions: Level 2 reviewer approval and Kiro Decision Policy

This pack includes Level 2-ready automation:

- `scripts/reviewer-openai.py` — optional OpenAI Responses API reviewer adapter that writes strict JSON review output.
- `scripts/reviewer-stage.sh` — prompt mode or auto mode review runner.
- `scripts/approve-stage.sh` — validates exact approval tokens and advances feature state.
- `scripts/spec-flow.sh` — runs requirements/design/tasks review flow with Kiro pauses.
- `scripts/kiro-decision.sh` — answers predictable Kiro decisions or pauses for architecture review.
- `docs/prompts/reviewer/approval-review.prompt.md` — strict JSON reviewer prompt.
- `docs/automation/KIRO_DECISION_POLICY.md` — human-readable Kiro decision policy.
- `.automation/kiro-decision-policy.yaml` — machine-readable policy starter.

Automated reviewer approval uses exact tokens only:

- `APPROVED_FOR_DESIGN`
- `APPROVED_FOR_TASKS`
- `APPROVED_FOR_CURSOR`

Anything else stops the flow.

## v3 additions: model recommendation policy

This pack includes a model recommendation layer:

- `.automation/model-policy.json`
- `docs/automation/MODEL_SELECTION_POLICY.md`
- `scripts/model-recommend.py`
- `scripts/model-record.sh`

Generated Kiro and Cursor prompts now include three prioritized LLM recommendations and require the tool output to include a `Model Execution Report` showing the actual selected model, effort setting, and fallback status.

The policy is configurable. Update `.automation/model-policy.json` if your Kiro or Cursor account uses different model IDs.

## Kiro/Cursor CLI plug-in mode

Headless mode is now the default. To run Kiro CLI headlessly explicitly:

```bash
export KIRO_API_KEY=ksk_xxxxxxxx
FEATURE_FACTORY_KIRO_MODE=auto make -f Makefile.feature-factory ff-kiro-stage FEATURE=FEATURE-0007 STAGE=requirements
```

To run Cursor CLI headlessly:

```bash
FEATURE_FACTORY_CURSOR_MODE=auto make -f Makefile.feature-factory ff-cursor-task FEATURE=FEATURE-0007 TASK=1.1
```

Both scripts log outputs under `.automation/logs/<FEATURE>/` and record model usage under `.automation/model-usage/<FEATURE>/usage.jsonl`.

## v3.1 default headless mode

This pack now defaults to headless execution instead of manual copy/paste:

```bash
FEATURE_FACTORY_KIRO_MODE=auto
FEATURE_FACTORY_CURSOR_MODE=auto
FEATURE_FACTORY_REVIEW_MODE=auto
```

Use prompt/manual mode only as an explicit fallback:

```bash
FEATURE_FACTORY_KIRO_MODE=prompt make -f Makefile.feature-factory ff-kiro-stage FEATURE=FEATURE-0007 STAGE=requirements
FEATURE_FACTORY_CURSOR_MODE=prompt make -f Makefile.feature-factory ff-cursor-task FEATURE=FEATURE-0007 TASK=1.1
```
