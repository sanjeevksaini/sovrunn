# Sovrunn Feature Factory Operating Model

## Purpose

The Feature Factory automates the repeatable part of Sovrunn development while keeping human gates for architecture decisions and merge approval.

## State machine

```text
NEW → BRANCH_CREATED → REQUIREMENTS_GENERATED → REQUIREMENTS_APPROVED → DESIGN_GENERATED → DESIGN_APPROVED → TASKS_GENERATED → TASKS_APPROVED → TASK_RUNNING → TASK_VERIFIED → TASK_COMMITTED → FINAL_VERIFIED → PR_CREATED → READY_TO_MERGE
```

## Human gates

- Requirements approval
- Design approval
- Tasks approval
- PR merge approval

## Level 1

Scripts now run Kiro CLI and Cursor CLI headlessly by default, while still generating prompts/logs for auditability. Manual prompt mode remains available as an explicit fallback.

## Level 2-ready

The same scripts expose integration points for CLI-based Kiro, reviewer, and Cursor execution.

## Guardrails

- Repo-local only.
- No deployment.
- No automatic merge.
- No credentials.
- No skipped tests.
- No forbidden package imports.
- No build artifacts committed.

## Level 2 Reviewer Approval

Feature Factory now supports strict reviewer JSON approval.

Approval tokens:

- requirements.md -> `APPROVED_FOR_DESIGN`
- design.md -> `APPROVED_FOR_TASKS`
- tasks.md -> `APPROVED_FOR_CURSOR`

Explicit manual prompt fallback:

```bash
make -f Makefile.feature-factory ff-review FEATURE=FEATURE-0007 STAGE=requirements
# paste generated prompt into reviewer, save JSON to .automation/reviews/FEATURE-0007/requirements.review.json
make -f Makefile.feature-factory ff-approve-requirements FEATURE=FEATURE-0007
```

Automated reviewer mode with OpenAI adapter:

```bash
export OPENAI_API_KEY="..."
export FEATURE_FACTORY_REVIEW_MODE=auto
export FEATURE_FACTORY_REVIEWER_MODEL="gpt-5"
make -f Makefile.feature-factory ff-review-auto FEATURE=FEATURE-0007 STAGE=requirements
make -f Makefile.feature-factory ff-approve-requirements FEATURE=FEATURE-0007
```

Full spec flow with Kiro CLI headless execution and automatic reviewer approval:

```bash
export OPENAI_API_KEY="..."
export FEATURE_FACTORY_REVIEW_MODE=auto
make -f Makefile.feature-factory ff-spec-flow FEATURE=FEATURE-0007
```

`scripts/spec-flow.sh` now calls Kiro CLI headlessly for requirements, design, tasks, and reviewer-requested revisions. Set `FEATURE_FACTORY_KIRO_MODE=prompt` only to force manual pauses.

## Kiro Decision Policy

Kiro decisions are controlled by:

- `docs/automation/KIRO_DECISION_POLICY.md`
- `.automation/kiro-decision-policy.yaml`
- `scripts/kiro-decision.sh`

Example:

```bash
make -f Makefile.feature-factory ff-kiro-decision \
  FEATURE=FEATURE-0007 \
  STAGE=requirements \
  QUESTION="Should I implement code now?"
```

The default policy is:

- Auto-accept predictable workflow decisions.
- Auto-reject implementation and common scope-expansion decisions.
- Pause for architecture-impacting decisions.

## Model recommendation and execution reporting

Feature Factory v3 adds a model recommendation layer.

For every Kiro stage and Cursor task, generated prompts include three prioritized LLM model recommendations. The recommendations are derived from `.automation/model-policy.json`.

### Kiro stages

- `requirements` uses the requirements profile.
- `design` uses the design profile.
- `tasks` uses the tasks profile.
- `revision` uses the revision profile.

### Cursor tasks

Cursor task prompts classify the task from `tasks.md` text:

- test/property/race/verification tasks -> `tests`
- bug/failure/debug tasks -> `debug_fix`
- handler/server/registry/interface/refactor tasks -> `complex_code`
- otherwise -> `routine_code`

### Execution report requirement

Every Kiro and Cursor response must include:

```text
Model Execution Report:
- Tool: Kiro | Cursor
- Stage or task: <requirements | design | tasks | Task N.N>
- Recommended priority list: <copied from prompt>
- Selected model: <actual model selected>
- Effort/reasoning setting: <actual setting if visible>
- Fallback used: yes/no
- Fallback reason: <unavailable/cost/latency/user override/none>
```

If the tool cannot show the selected model in output, the response must say `unknown/not visible in tool output`.

### Recording actual model usage

After a Kiro/Cursor execution, record the selected model:

```bash
make -f Makefile.feature-factory ff-model-record \
  FEATURE=FEATURE-0007 \
  TOOL=cursor \
  TASK=1.1 \
  SELECTED_MODEL="GPT-5.6 Terra" \
  EFFORT="Medium"
```

This appends one JSONL entry to:

```text
.automation/model-usage/FEATURE-0007/usage.jsonl
```

### Manual recommendation lookup

```bash
make -f Makefile.feature-factory ff-model-recommend TOOL=kiro STAGE=requirements
make -f Makefile.feature-factory ff-model-recommend TOOL=cursor TASK=1.1 TASKS_PATH=.kiro/specs/plugin-capability-registry/tasks.md
```

## Automated revision loop

`review-and-route-stage.sh` turns review output into routing behavior.

Reviewer JSON status handling:

- `APPROVED` with exact token -> calls `approve-stage.sh` and advances state.
- `NEEDS_REVISION` -> writes `.automation/generated-prompts/<FEATURE>/<stage>.revision.prompt.md` and exits with code 2.
- `BLOCKED` -> stops the flow and marks the stage blocked.

The revision prompt is stage-locked:

- requirements revision may revise `requirements.md` only.
- design revision may revise `design.md` only.
- tasks revision may revise `tasks.md` only.

Default maximum revision attempts is 3 and can be overridden with:

```bash
FEATURE_FACTORY_MAX_REVISIONS=5 make -f Makefile.feature-factory ff-spec-flow FEATURE=FEATURE-0007
```

Manual one-stage routing:

```bash
make -f Makefile.feature-factory ff-review-route FEATURE=FEATURE-0007 STAGE=requirements
```

If exit code is 2, paste the generated revision prompt into Kiro, wait for the file to be updated, then rerun the command.

## CLI integration: Kiro and Cursor

The Feature Factory can now run Kiro and Cursor from the terminal. Kiro CLI and Cursor CLI headless execution are now the defaults. Prompt/manual mode is retained only as an explicit fallback.

### Kiro CLI

Kiro headless mode is enabled with:

```bash
FEATURE_FACTORY_KIRO_MODE=auto ./scripts/kiro-stage.sh --feature FEATURE-0007 --stage requirements
```

Required:

```bash
export KIRO_API_KEY=ksk_xxxxxxxx
```

Useful options:

```bash
export KIRO_EFFORT=high               # low|medium|high|xhigh|max
export KIRO_TRUST_TOOLS=read,grep,write
export KIRO_TRUST_ALL_TOOLS=1          # use only in trusted repo-local automation
export KIRO_SELECTED_MODEL="Claude Opus 4.8"  # recorded in model usage log
export KIRO_AGENT=sovrunn-spec-agent   # optional custom Kiro agent
```

Kiro CLI currently runs the prompt with `kiro-cli chat --no-interactive`. The model priority list is embedded in the prompt and the selected model is recorded in `.automation/model-usage/<FEATURE>/usage.jsonl`. If Kiro exposes a stable per-run `--model` flag in your installed version, set it through your Kiro configuration or adapt `scripts/kiro-stage.sh` locally.

### Cursor CLI

Cursor headless mode is enabled with:

```bash
FEATURE_FACTORY_CURSOR_MODE=auto ./scripts/cursor-task.sh --feature FEATURE-0007 --task 1.1
```

Useful options:

```bash
export CURSOR_AGENT_BIN=cursor-agent
export CURSOR_SELECTED_MODEL=gpt-5.6-terra   # optional override
export CURSOR_OUTPUT_FORMAT=text
export FEATURE_FACTORY_CURSOR_VERIFY=1       # run scripts/verify.sh after Cursor returns
```

`cursor-task.sh` renders the task prompt, reads the model recommendations, then tries the recommended Cursor models in priority order using `cursor-agent -p --model <model>`. If the first model fails or is unavailable, it tries the next recommended model and records fallback usage.

Logs are written to:

```text
.automation/logs/<FEATURE>/kiro-<stage>.log
.automation/logs/<FEATURE>/cursor-task-<TASK>.log
```

Model usage is recorded to:

```text
.automation/model-usage/<FEATURE>/usage.jsonl
```

## v3.1 update: default headless Kiro/Cursor mode

This pack has been updated so the normal path no longer depends on manual copy/paste into Kiro or Cursor.

Defaults:

```bash
FEATURE_FACTORY_KIRO_MODE=auto
FEATURE_FACTORY_CURSOR_MODE=auto
FEATURE_FACTORY_REVIEW_MODE=auto
```

Fallback manual mode is still available for debugging:

```bash
FEATURE_FACTORY_KIRO_MODE=prompt make -f Makefile.feature-factory ff-kiro-stage FEATURE=FEATURE-0007 STAGE=requirements
FEATURE_FACTORY_CURSOR_MODE=prompt make -f Makefile.feature-factory ff-cursor-task FEATURE=FEATURE-0007 TASK=1.1
```

The spec flow now executes Kiro headlessly for initial stage generation and for reviewer-requested revisions.

## Phase 2 Reuse and Drift Gates

Every FEATURE-0011-and-later feature must include a reuse assessment that
conforms to the canonical standard:

`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

Do not duplicate or redefine the assessment field schema in this document.
Populate the feature-level reuse summary and capability-level assessments
using the canonical fields and controlled vocabularies.

Architecture drift checks:

- no provider-specific hardcoding in core,
- no Kubernetes-only assumptions in core,
- no PostgreSQL lifecycle logic in core placement engine,
- no custom policy engine embedded in handlers,
- no raw secret storage,
- no customer-facing IaaS leakage,
- explainable decision object,
- defined audit behavior,
- preserved adapter boundaries.
