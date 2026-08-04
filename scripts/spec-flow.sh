#!/opt/homebrew/bin/bash
set -euo pipefail
source "$(dirname "$0")/common.sh"

FEATURE=""
MODE="${FEATURE_FACTORY_REVIEW_MODE:-auto}"
KIRO_MODE="${FEATURE_FACTORY_KIRO_MODE:-auto}"
MAX_REVISIONS="${FEATURE_FACTORY_MAX_REVISIONS:-5}"
RESUME="${FEATURE_FACTORY_RESUME:-1}"
STARTED_EPOCH="$(date +%s)"
KIRO_INVOCATIONS=0
FLOW_STATUS="FAILED"
ACTION_REQUIRED="Inspect the terminal output and .automation logs; correct the failure, then rerun make ff-spec-flow."

while [[ $# -gt 0 ]]; do
  case "$1" in
    --feature) FEATURE="$2"; shift 2;;
    --mode) MODE="$2"; shift 2;;
    --kiro-mode) KIRO_MODE="$2"; shift 2;;
    --max-revisions) MAX_REVISIONS="$2"; shift 2;;
    *) fail "unknown arg: $1";;
  esac
done

[[ -n "$FEATURE" ]] || fail "--feature required"
[[ "$MODE" == "auto" ]] || fail "unattended spec flow requires reviewer mode auto"
[[ "$KIRO_MODE" == "auto" ]] || fail "unattended spec flow requires Kiro mode auto"
configure_reviewer_adapter
info "Specification reviewer adapter: $FEATURE_FACTORY_REVIEWER_CMD"
command -v kiro-cli >/dev/null 2>&1 || command -v kiro >/dev/null 2>&1 || fail "Kiro CLI is not available"
cd "$(repo_root)"
ensure_feature_state "$FEATURE"
if [[ -f ".automation/features/${FEATURE}.control.json" && -z "${FEATURE_FACTORY_MAX_REVISIONS:-}" ]]; then
  MAX_REVISIONS="$(python3 - "$FEATURE" <<'PY'
import json, sys
from pathlib import Path
control = json.loads(Path(f'.automation/features/{sys.argv[1]}.control.json').read_text())
print(control['execution']['max_stage_revisions'])
PY
)"
fi
if [[ "$FEATURE" == "FEATURE-0013" ]]; then
  if [[ -n "${KIRO_AGENT:-}" && "${KIRO_AGENT}" != "sovrunn-spec" ]]; then
    fail "FEATURE-0013 requires KIRO_AGENT=sovrunn-spec; refusing override: ${KIRO_AGENT}"
  fi
  export KIRO_AGENT="sovrunn-spec"
fi
SPEC_PATH=$(get_feature_value "$FEATURE" spec_path)
BASELINE_LINES="$(find "$SPEC_PATH" -maxdepth 1 -type f \( -name 'requirements.md' -o -name 'design.md' -o -name 'tasks.md' \) -exec cat {} + 2>/dev/null | wc -l | tr -d ' ')"

write_report() {
  ./scripts/spec-flow-report.py \
    --feature "$FEATURE" \
    --spec-path "$SPEC_PATH" \
    --status "$FLOW_STATUS" \
    --started-epoch "$STARTED_EPOCH" \
    --kiro-invocations "$KIRO_INVOCATIONS" \
    --baseline-lines "$BASELINE_LINES" \
    --action "$ACTION_REQUIRED" || true
}
trap write_report EXIT

run_kiro_prompt_file() {
  local stage="$1" prompt_file="$2" expected_doc="$3"
  if [[ "$KIRO_MODE" == "prompt" || "$KIRO_MODE" == "manual" ]]; then
    info "Prompt/manual Kiro mode: run $prompt_file in Kiro, then press Enter."
    read -r -p "Press Enter after Kiro has updated $expected_doc... " _
    [[ -f "$expected_doc" ]] || fail "missing expected Kiro output: $expected_doc"
    return 0
  fi

  KIRO_BIN="${KIRO_CLI_BIN:-}"
  if [[ -z "$KIRO_BIN" ]]; then
    if command -v kiro-cli >/dev/null 2>&1; then
      KIRO_BIN="kiro-cli"
    elif command -v kiro >/dev/null 2>&1; then
      KIRO_BIN="kiro"
    else
      fail "Kiro CLI not found in PATH. Install/configure Kiro CLI or set KIRO_CLI_BIN."
    fi
  fi
  command -v "$KIRO_BIN" >/dev/null 2>&1 || fail "Kiro CLI command not found: $KIRO_BIN"
  if [[ -z "${KIRO_API_KEY:-}" && "${KIRO_ALLOW_LOCAL_AUTH:-1}" != "1" ]]; then
    fail "KIRO_API_KEY is required when KIRO_ALLOW_LOCAL_AUTH is not enabled"
  fi

  mkdir -p ".automation/logs/$FEATURE"
  local log_file=".automation/logs/$FEATURE/kiro-${stage}-revision.log"
  local effort="${KIRO_EFFORT:-high}"
  effort="$(echo "$effort" | tr '[:upper:]' '[:lower:]')"
  case "$effort" in low|medium|high|xhigh|max) ;; *) effort="high";; esac
  local model="${KIRO_MODEL:-}"
  if [[ "$FEATURE" == "FEATURE-0013" && -z "$model" ]]; then
    model="claude-opus-4.8"
  fi

  local prompt_text
  prompt_text="$(cat "$prompt_file")"
  local args=(chat --no-interactive --effort "$effort")
  if [[ -n "$model" ]]; then
    args+=(--model "$model")
  fi
  if [[ -n "${KIRO_AGENT:-}" ]]; then
    args+=(--agent "$KIRO_AGENT")
  fi
  if [[ "${KIRO_TRUST_ALL_TOOLS:-0}" == "1" ]]; then
    args+=(--trust-all-tools)
  else
    args+=(--trust-tools="${KIRO_TRUST_TOOLS:-read,grep,write}")
  fi

  info "Running Kiro CLI headlessly for $FEATURE $stage revision"
  if [[ -f ".automation/features/${FEATURE}.control.json" ]]; then
    PYTHONDONTWRITEBYTECODE=1 python3 \
      ./scripts/generic-kiro-boundary-check.py \
      --feature "$FEATURE" --stage "$stage" --mode pre
  fi
  KIRO_INVOCATIONS=$((KIRO_INVOCATIONS + 1))
  set +e
  "$KIRO_BIN" "${args[@]}" "$prompt_text" 2>&1 | tee "$log_file"
  local status=${PIPESTATUS[0]}
  set -e
  [[ $status -eq 0 ]] || fail "Kiro CLI revision failed for $stage. See $log_file"
  if [[ -f ".automation/features/${FEATURE}.control.json" ]]; then
    ./scripts/receipt-check.py --log "$log_file" --kind stage || \
      fail "Kiro revision did not produce exactly one COMPLETE receipt. See $log_file"
    PYTHONDONTWRITEBYTECODE=1 python3 \
      ./scripts/generic-kiro-boundary-check.py \
      --feature "$FEATURE" --stage "$stage" --mode post
  fi
  [[ -f "$expected_doc" ]] || fail "missing expected Kiro output after revision: $expected_doc"
  ./scripts/feature-state.py set --feature "$FEATURE" --key status --value "${stage}_generated" >/dev/null
}

run_stage() {
  local stage="$1" file_target="$2"
  if [[ "$RESUME" == "1" && -f "$file_target" ]]; then
    info "Resume mode: using existing $file_target"
  else
    info "Generating $stage with Kiro CLI headless mode: $KIRO_MODE"
    KIRO_INVOCATIONS=$((KIRO_INVOCATIONS + 1))
    ./scripts/kiro-stage.sh --feature "$FEATURE" --stage "$stage" --mode "$KIRO_MODE"
  fi
  [[ -f "$file_target" ]] || fail "missing expected Kiro output: $file_target"

  local pending_status
  pending_status="$(python3 - "$FEATURE" <<'PY'
import json, sys
from pathlib import Path
print(json.loads(Path(f'.automation/state/{sys.argv[1]}.json').read_text()).get('status', ''))
PY
)"
  if [[ "$pending_status" == "${stage}_revision_required" ]]; then
    local pending_prompt=".automation/generated-prompts/$FEATURE/${stage}.revision.prompt.md"
    [[ -f "$pending_prompt" ]] || fail "missing pending revision prompt: $pending_prompt"
    info "Resume mode: applying pending $stage revision before re-review"
    run_kiro_prompt_file "$stage" "$pending_prompt" "$file_target"
  fi

  while true; do
    set +e
    ./scripts/review-and-route-stage.sh --feature "$FEATURE" --stage "$stage" --mode "$MODE" --max-revisions "$MAX_REVISIONS"
    rc=$?
    set -e
    if [[ $rc -eq 0 ]]; then
      info "$stage approved."
      break
    elif [[ $rc -eq 2 ]]; then
      revision_prompt=".automation/generated-prompts/$FEATURE/${stage}.revision.prompt.md"
      [[ -f "$revision_prompt" ]] || fail "missing generated revision prompt: $revision_prompt"
      info "Revision required. Running Kiro CLI headlessly with $revision_prompt"
      run_kiro_prompt_file "$stage" "$revision_prompt" "$file_target"
    else
      FLOW_STATUS="HUMAN_ACTION_REQUIRED"
      ACTION_REQUIRED="Stage $stage is blocked or its automation failed. Review .automation/reviews/$FEATURE/${stage}.review.json and .automation/logs/$FEATURE before resuming."
      fail "$stage review blocked or failed"
    fi
  done
}

CONTROL_FILE=".automation/features/${FEATURE}.control.json"
if [[ -f "$CONTROL_FILE" ]]; then
  PLAN_TOKEN="$(python3 - "$FEATURE" <<'PY'
import json, sys
from pathlib import Path
state = json.loads(Path(f'.automation/state/{sys.argv[1]}.json').read_text())
print(state.get('executable_plan_approval_token', ''))
PY
)"
  if [[ "$PLAN_TOKEN" != "APPROVED_EXECUTABLE_PLAN" ]]; then
    run_stage requirements "$SPEC_PATH/requirements.md"
    run_stage design "$SPEC_PATH/design.md"
    ./scripts/executable-plan-report.py --feature "$FEATURE"
    ./scripts/feature-state.py set --feature "$FEATURE" --key current_stage --value executable_plan >/dev/null
    ./scripts/feature-state.py set --feature "$FEATURE" --key status --value executable_plan_human_review_required >/dev/null
    ./scripts/feature-state.py set --feature "$FEATURE" --key human_gate_required --value true >/dev/null
    FLOW_STATUS="HUMAN_ACTION_REQUIRED"
    ACTION_REQUIRED="Review approved requirements and design together, then run make ff-approve-human-gate FEATURE=$FEATURE GATE=executable_plan APPROVED_BY='<name>'."
    info "$ACTION_REQUIRED"
    exit 3
  fi
  run_stage tasks "$SPEC_PATH/tasks.md"
else
  run_stage requirements "$SPEC_PATH/requirements.md"
  run_stage design "$SPEC_PATH/design.md"
  run_stage tasks "$SPEC_PATH/tasks.md"
fi
FLOW_STATUS="COMPLETE"
ACTION_REQUIRED=""
info "Spec flow complete. Status:"
./scripts/feature-state.py get --feature "$FEATURE"
