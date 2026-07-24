#!/opt/homebrew/bin/bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

FEATURE=""
TASK=""
MODE="${FEATURE_FACTORY_CURSOR_MODE:-auto}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --feature) FEATURE="$2"; shift 2;;
    --task) TASK="$2"; shift 2;;
    --mode) MODE="$2"; shift 2;;
    *) fail "unknown arg: $1";;
  esac
done

[[ -n "$FEATURE" ]] || fail "--feature required"
[[ -n "$TASK" ]] || fail "--task required"
case "$MODE" in prompt|manual|auto) ;; *) fail "unsupported mode: $MODE";; esac

cd "$(repo_root)"
ensure_feature_state "$FEATURE"

SPEC_PATH="$(get_feature_value "$FEATURE" spec_path)"
OUT_DIR="$(get_feature_value "$FEATURE" generated_prompt_path)"
mkdir -p "$OUT_DIR"
OUT="$OUT_DIR/cursor-task-${TASK}.prompt.md"
TASKS_PATH="$SPEC_PATH/tasks.md"
[[ -f "$TASKS_PATH" ]] || fail "missing tasks.md: $TASKS_PATH"

python3 - "$FEATURE" "$TASK" "$SPEC_PATH" "$OUT" <<'PYCURSOR'
from pathlib import Path
import subprocess, sys
feature, task, spec_path, out = sys.argv[1:]
template = Path('docs/prompts/cursor/task.prompt.md').read_text()
tasks_path = str(Path(spec_path, 'tasks.md'))
tasks = Path(tasks_path).read_text()
try:
    model_recs = subprocess.check_output(['./scripts/model-recommend.py','--tool','cursor','--task',task,'--tasks-path',tasks_path], text=True)
except Exception as e:
    model_recs = f'Model recommendation unavailable: {e}'
rendered = (template
  .replace('{{FEATURE_ID}}', feature)
  .replace('{{TASK_ID}}', task)
  .replace('{{SPEC_PATH}}', spec_path)
  .replace('{{REQUIREMENTS_PATH}}', f'{spec_path}/requirements.md')
  .replace('{{DESIGN_PATH}}', f'{spec_path}/design.md')
  .replace('{{TASKS_PATH}}', tasks_path)
  .replace('{{MODEL_RECOMMENDATIONS}}', model_recs)
  .replace('{{TASKS_CONTENT}}', tasks))
Path(out).write_text(rendered)
print(out)
PYCURSOR

./scripts/feature-state.py set --feature "$FEATURE" --key current_task --value "$TASK" >/dev/null
./scripts/feature-state.py set --feature "$FEATURE" --key status --value "cursor_prompt_generated" >/dev/null
info "Cursor task prompt generated: $OUT"

if [[ -n "${FEATURE_FACTORY_ALLOWED_PATHS:-}" ]]; then
  {
    cat <<'EOF'

## Automation-enforced path scope

Every changed file must match one of the allowed rules below. Do not
create, modify, rename, or delete files outside this scope.

Reuse existing verification helpers when possible. Do not add a test in
another package unless that path or directory is explicitly listed.

Allowed path rules:
EOF

    while IFS= read -r allowed_file; do
      [[ -n "$allowed_file" ]] || continue
      printf -- '- `%s`\n' "$allowed_file"
    done <<<"$FEATURE_FACTORY_ALLOWED_PATHS"

    cat <<'EOF'

Before reporting completion, run `git diff --name-only HEAD` and confirm
that every changed path is allowed.
EOF
  } >> "$OUT"
fi

if [[ "$MODE" == "prompt" || "$MODE" == "manual" ]]; then
  info "Prompt/manual mode: paste this into Cursor. Default is auto/headless; set FEATURE_FACTORY_CURSOR_MODE=prompt only when you intentionally want manual mode."
  exit 0
fi

CURSOR_BIN="${CURSOR_AGENT_BIN:-}"
if [[ -z "$CURSOR_BIN" ]]; then
  if command -v cursor-agent >/dev/null 2>&1; then
    CURSOR_BIN="$(command -v cursor-agent)"
  elif command -v agent >/dev/null 2>&1; then
    CURSOR_BIN="$(command -v agent)"
  else
    fail "Cursor Agent CLI not found. Install cursor-agent or set CURSOR_AGENT_BIN to its absolute path."
  fi
elif [[ "$CURSOR_BIN" != */* ]]; then
  CURSOR_BIN="$(command -v "$CURSOR_BIN")" ||
    fail "Cursor Agent CLI command not found: $CURSOR_BIN"
fi

[[ -x "$CURSOR_BIN" ]] ||
  fail "Cursor Agent CLI is not executable: $CURSOR_BIN"

case "$(basename "$CURSOR_BIN")" in
  cursor|Cursor)
    fail "refusing desktop Cursor launcher; use cursor-agent or set CURSOR_AGENT_BIN"
    ;;
esac

CURSOR_HELP="$("$CURSOR_BIN" --help 2>&1 || true)"
for required_option in --output-format --model --force; do
  grep -q -- "$required_option" <<<"$CURSOR_HELP" ||
    fail "$CURSOR_BIN is not a compatible Cursor Agent CLI: missing $required_option"
done

mkdir -p ".automation/logs/$FEATURE"
LOG_FILE=".automation/logs/$FEATURE/cursor-task-${TASK}.log"
MODEL_JSON=".automation/logs/$FEATURE/cursor-task-${TASK}.model.json"
./scripts/model-recommend.py --tool cursor --task "$TASK" --tasks-path "$TASKS_PATH" --format json > "$MODEL_JSON"

mapfile -t MODEL_ROWS < <(python3 - "$MODEL_JSON" <<'PY'
import json, sys
j=json.load(open(sys.argv[1]))
for r in j['recommendations']:
    print('\t'.join([r.get('model',''), r.get('label',''), r.get('effort','medium')]))
PY
)

PROMPT_TEXT="$(cat "$OUT")"
OUTPUT_FORMAT="${CURSOR_OUTPUT_FORMAT:-text}"
VERIFY_AFTER="${FEATURE_FACTORY_CURSOR_VERIFY:-1}"

MODELS_TO_TRY=()
if [[ -n "${CURSOR_SELECTED_MODEL:-}" ]]; then
  MODELS_TO_TRY+=("${CURSOR_SELECTED_MODEL}")
fi
for row in "${MODEL_ROWS[@]}"; do
  IFS=$'\t' read -r model label effort <<< "$row"
  [[ -n "$model" ]] && MODELS_TO_TRY+=("$model")
done

SUCCESS=0
SELECTED_MODEL=""
SELECTED_EFFORT=""
FALLBACK_USED="no"
FALLBACK_REASON="none"
: > "$LOG_FILE"

for idx in "${!MODELS_TO_TRY[@]}"; do
  model="${MODELS_TO_TRY[$idx]}"
  [[ -n "$model" ]] || continue
  if [[ "$idx" != "0" ]]; then
    FALLBACK_USED="yes"
    FALLBACK_REASON="previous_model_failed_or_unavailable"
  fi
  SELECTED_MODEL="$model"
  # Find effort from recommendation list where possible.
  SELECTED_EFFORT="Medium"
  for row in "${MODEL_ROWS[@]}"; do
    IFS=$'\t' read -r rec_model rec_label rec_effort <<< "$row"
    if [[ "$rec_model" == "$model" || "$rec_label" == "$model" ]]; then
      SELECTED_EFFORT="$rec_effort"
      break
    fi
  done
  info "Running Cursor CLI for $FEATURE Task $TASK with model: $model"
  {
    echo "=== Cursor task $TASK model=$model output_format=$OUTPUT_FORMAT ==="
    date -u '+%Y-%m-%dT%H:%M:%SZ'
  } >> "$LOG_FILE"
  set +e
  "$CURSOR_BIN" -p --trust --force --output-format "$OUTPUT_FORMAT" --model "$model" "$PROMPT_TEXT" 2>&1 | tee -a "$LOG_FILE"
  STATUS=${PIPESTATUS[0]}
  set -e
  if [[ $STATUS -eq 0 ]]; then
    SUCCESS=1
    break
  fi
  echo "Cursor CLI failed with model $model, trying next recommended model if available." | tee -a "$LOG_FILE" >&2
done

./scripts/model-record.sh \
  --feature "$FEATURE" \
  --tool cursor \
  --task "$TASK" \
  --selected-model "${SELECTED_MODEL:-unknown}" \
  --effort "${SELECTED_EFFORT:-unknown}" \
  --fallback-used "$FALLBACK_USED" \
  --fallback-reason "$FALLBACK_REASON" >/dev/null || true

if [[ "$SUCCESS" != "1" ]]; then
  ./scripts/feature-state.py set --feature "$FEATURE" --key status --value "cursor_task_${TASK}_failed" >/dev/null || true
  fail "Cursor CLI failed for all recommended models. See $LOG_FILE"
fi

if [[ "$VERIFY_AFTER" == "1" ]]; then
  info "Running post-Cursor verification"
  ./scripts/verify.sh | tee -a "$LOG_FILE"
fi

./scripts/feature-state.py set --feature "$FEATURE" --key status --value "cursor_task_${TASK}_completed" >/dev/null
info "Cursor task completed. Log: $LOG_FILE"
