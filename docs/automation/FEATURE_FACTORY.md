# Sovrunn Feature Factory Operating Model

## Purpose

The Feature Factory automates the repeatable parts of Sovrunn development while preserving human gates for architecture decisions, specification approval, implementation authorization, and merge approval.

## State machine

```text
NEW
→ BRANCH_CREATED
→ REQUIREMENTS_GENERATED
→ REQUIREMENTS_APPROVED
→ DESIGN_GENERATED
→ DESIGN_APPROVED
→ TASKS_GENERATED
→ TASKS_APPROVED
→ TASK_RUNNING
→ TASK_VERIFIED
→ TASK_COMMITTED
→ FINAL_VERIFIED
→ PR_CREATED
→ READY_TO_MERGE
```

## Human gates

- Requirements approval
- Design approval
- Tasks approval
- Implementation authorization
- PR merge approval

## Governed specification workflow

Sovrunn specifications must be generated and approved one stage at a time:

```text
requirements.md
→ APPROVED_FOR_DESIGN
→ design.md
→ APPROVED_FOR_TASKS
→ tasks.md
→ APPROVED_FOR_CURSOR
→ implementation
```

Kiro Quick Plan, or any equivalent mode that generates requirements, design,
and tasks without preserving these approval boundaries, must not be used for
governed Sovrunn features.

Automation may orchestrate the full flow, but it must:

- generate only the active stage;
- obtain and persist the exact approval token for that stage;
- stop on `NEEDS_REVISION` or `BLOCKED`;
- advance only after successful approval;
- preserve prompts, reviews, decisions, model records, and execution logs;
- never bypass requirements, design, tasks, implementation, or merge gates.

Manual stage execution remains available:

```bash
make -f Makefile.feature-factory \
  ff-kiro-stage \
  FEATURE=<FEATURE-ID> \
  STAGE=requirements
```

## Level 1

Scripts run Kiro CLI and Cursor CLI headlessly by default while generating prompts and logs for auditability. Manual prompt mode remains available as an explicit fallback.

## Level 2-ready

The same scripts expose integration points for CLI-based Kiro, reviewer, and Cursor execution.

## Guardrails

- Repository-local operations only.
- No deployment.
- No automatic merge.
- No credentials in prompts, logs, specifications, or source.
- No skipped required tests.
- No forbidden package imports.
- No committed build artifacts.
- No silent architecture or scope expansion.
- No generation of later specification stages before approval.

## Level 2 Reviewer Approval

Feature Factory supports strict reviewer JSON approval.

Approval tokens:

- `requirements.md` → `APPROVED_FOR_DESIGN`
- `design.md` → `APPROVED_FOR_TASKS`
- `tasks.md` → `APPROVED_FOR_CURSOR`

Explicit manual prompt fallback:

```bash
make -f Makefile.feature-factory \
  ff-review \
  FEATURE=FEATURE-0007 \
  STAGE=requirements

# Paste the generated prompt into the reviewer and save the JSON to:
# .automation/reviews/FEATURE-0007/requirements.review.json

make -f Makefile.feature-factory \
  ff-approve-requirements \
  FEATURE=FEATURE-0007
```

Automated reviewer mode with the OpenAI adapter:

```bash
export OPENAI_API_KEY="..."
export FEATURE_FACTORY_REVIEW_MODE=auto
export FEATURE_FACTORY_REVIEWER_MODEL="gpt-5"

make -f Makefile.feature-factory \
  ff-review-auto \
  FEATURE=FEATURE-0007 \
  STAGE=requirements

make -f Makefile.feature-factory \
  ff-approve-requirements \
  FEATURE=FEATURE-0007
```

Full governed specification flow with Kiro CLI headless execution and automated
stage-by-stage reviewer routing:

```bash
export OPENAI_API_KEY="..."
export FEATURE_FACTORY_REVIEW_MODE=auto

make -f Makefile.feature-factory \
  ff-spec-flow \
  FEATURE=FEATURE-0007
```

`ff-spec-flow` is not a Quick Plan workflow. It generates one stage at a time,
runs the reviewer for that stage, persists the review result, and advances only
when the exact approval token is present.

`scripts/spec-flow.sh` calls Kiro CLI headlessly for requirements, design,
tasks, and reviewer-requested revisions. Set
`FEATURE_FACTORY_KIRO_MODE=prompt` only to force manual pauses.

## Kiro Decision Policy

Kiro decisions are controlled by:

- `docs/automation/KIRO_DECISION_POLICY.md`
- `.automation/kiro-decision-policy.yaml`
- `scripts/kiro-decision.sh`

Example:

```bash
make -f Makefile.feature-factory \
  ff-kiro-decision \
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

For every Kiro stage and Cursor task, generated prompts include three
prioritized model recommendations. Recommendations are derived from:

```text
.automation/model-policy.json
```

### Kiro stages

- `requirements` uses the requirements profile.
- `design` uses the design profile.
- `tasks` uses the tasks profile.
- `revision` uses the revision profile.

### Cursor tasks

Cursor task prompts classify the task from `tasks.md` text:

- test, property, race, or verification tasks → `tests`
- bug, failure, or debugging tasks → `debug_fix`
- handler, server, registry, interface, or refactor tasks → `complex_code`
- all other tasks → `routine_code`

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

If the tool cannot show the selected model, the response must state:

```text
unknown/not visible in tool output
```

### Recording actual model usage

After a Kiro or Cursor execution, record the selected model:

```bash
make -f Makefile.feature-factory \
  ff-model-record \
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
make -f Makefile.feature-factory \
  ff-model-recommend \
  TOOL=kiro \
  STAGE=requirements

make -f Makefile.feature-factory \
  ff-model-recommend \
  TOOL=cursor \
  TASK=1.1 \
  TASKS_PATH=.kiro/specs/plugin-capability-registry/tasks.md
```

## Automated revision loop

`review-and-route-stage.sh` turns review output into routing behavior.

Reviewer JSON status handling:

- `APPROVED` with the exact token calls `approve-stage.sh` and advances state.
- `NEEDS_REVISION` writes
  `.automation/generated-prompts/<FEATURE>/<stage>.revision.prompt.md`
  and exits with code `2`.
- `BLOCKED` stops the flow and marks the stage blocked.

The revision prompt is stage-locked:

- requirements revision may revise `requirements.md` only;
- design revision may revise `design.md` only;
- tasks revision may revise `tasks.md` only.

The default maximum revision count is three. Override it with:

```bash
FEATURE_FACTORY_MAX_REVISIONS=5 \
make -f Makefile.feature-factory \
  ff-spec-flow \
  FEATURE=FEATURE-0007
```

Manual one-stage routing:

```bash
make -f Makefile.feature-factory \
  ff-review-route \
  FEATURE=FEATURE-0007 \
  STAGE=requirements
```

If exit code is `2`, the automated flow passes the stage-locked revision prompt
back to Kiro. In explicit prompt mode, paste the generated revision prompt into
Kiro, wait for the active-stage file to be updated, and rerun the command.

## CLI integration: Kiro and Cursor

The Feature Factory can run Kiro and Cursor from the terminal. Headless
execution is the default; prompt/manual mode is an explicit fallback.

### Kiro CLI

Kiro headless mode:

```bash
FEATURE_FACTORY_KIRO_MODE=auto \
./scripts/kiro-stage.sh \
  --feature FEATURE-0007 \
  --stage requirements
```

Required when the installed Kiro CLI uses API-key authentication:

```bash
export KIRO_API_KEY=ksk_xxxxxxxx
```

Useful options:

```bash
export KIRO_EFFORT=high                  # low|medium|high|xhigh|max
export KIRO_TRUST_TOOLS=read,grep,write
export KIRO_TRUST_ALL_TOOLS=1            # trusted repository-local automation only
export KIRO_SELECTED_MODEL="Claude Opus 4.8"
export KIRO_AGENT=sovrunn-spec-agent
```

`KIRO_SELECTED_MODEL` records the declared model in the usage log; it does not
select the Kiro model by itself. Model selection is controlled by the installed
Kiro CLI configuration or the configured Kiro agent.

Kiro CLI runs the rendered prompt through:

```text
kiro-cli chat --no-interactive
```

The model priority list is embedded in the prompt. Actual usage is recorded in:

```text
.automation/model-usage/<FEATURE>/usage.jsonl
```

### Cursor CLI

Cursor headless mode:

```bash
FEATURE_FACTORY_CURSOR_MODE=auto \
./scripts/cursor-task.sh \
  --feature FEATURE-0007 \
  --task 1.1
```

Useful options:

```bash
export CURSOR_AGENT_BIN=cursor-agent
export CURSOR_SELECTED_MODEL=gpt-5.6-terra
export CURSOR_OUTPUT_FORMAT=text
export FEATURE_FACTORY_CURSOR_VERIFY=1
```

`cursor-task.sh` renders the task prompt, reads model recommendations, and
tries the recommended Cursor models in priority order using:

```text
cursor-agent -p --model <model>
```

If the first model fails or is unavailable, the script tries the next
recommended model and records fallback usage.

Logs are written to:

```text
.automation/logs/<FEATURE>/kiro-<stage>.log
.automation/logs/<FEATURE>/cursor-task-<TASK>.log
```

Model usage is recorded in:

```text
.automation/model-usage/<FEATURE>/usage.jsonl
```

## v3.1 update: default headless Kiro/Cursor mode

The normal path does not depend on manual copy and paste into Kiro or Cursor.

Defaults:

```bash
FEATURE_FACTORY_KIRO_MODE=auto
FEATURE_FACTORY_CURSOR_MODE=auto
FEATURE_FACTORY_REVIEW_MODE=auto
```

Manual fallback mode remains available for debugging:

```bash
FEATURE_FACTORY_KIRO_MODE=prompt \
make -f Makefile.feature-factory \
  ff-kiro-stage \
  FEATURE=FEATURE-0007 \
  STAGE=requirements

FEATURE_FACTORY_CURSOR_MODE=prompt \
make -f Makefile.feature-factory \
  ff-cursor-task \
  FEATURE=FEATURE-0007 \
  TASK=1.1
```

The specification flow executes Kiro headlessly for initial stage generation
and reviewer-requested revisions while preserving each approval boundary.

## Cross-Phase Reuse Assessment Gate

Every FEATURE-0011-and-later feature across all Sovrunn phases must include a
reuse assessment conforming to the current canonical standard:

```text
docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md
```

The path remains under `docs/phase2/` until the approved controlled migration
to a cross-phase architecture location is implemented. Its applicability is
not limited to Phase 2.

Do not duplicate or redefine the assessment field schema in this document.
Populate feature-level summaries and capability-level assessments using the
canonical fields and controlled vocabularies.

## Phase 2 Architecture Drift Checks

For Phase 2 features, the architecture gates additionally require:

- no provider-specific hardcoding in core;
- no Kubernetes-only assumptions in core;
- no PostgreSQL lifecycle logic in the core placement engine;
- no custom policy engine embedded in handlers;
- no raw secret storage;
- no customer-facing IaaS leakage;
- explainable decision objects;
- defined audit behavior;
- preserved adapter boundaries;
- conformance with `docs/architecture/api-resource-standard.md` for
  FEATURE-0012-and-later resource and API contracts.
