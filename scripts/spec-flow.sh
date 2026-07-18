#!/usr/bin/env bash
set -euo pipefail
source "$(dirname "$0")/common.sh"

FEATURE=""
MODE="${FEATURE_FACTORY_REVIEW_MODE:-auto}"
KIRO_MODE="${FEATURE_FACTORY_KIRO_MODE:-auto}"
MAX_REVISIONS="${FEATURE_FACTORY_MAX_REVISIONS:-3}"

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
cd "$(repo_root)"
ensure_feature_state "$FEATURE"
SPEC_PATH=$(get_feature_value "$FEATURE" spec_path)

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

  local prompt_text
  prompt_text="$(cat "$prompt_file")"
  local args=(chat --no-interactive --effort "$effort")
  if [[ -n "${KIRO_AGENT:-}" ]]; then
    args+=(--agent "$KIRO_AGENT")
  fi
  if [[ "${KIRO_TRUST_ALL_TOOLS:-0}" == "1" ]]; then
    args+=(--trust-all-tools)
  else
    args+=(--trust-tools="${KIRO_TRUST_TOOLS:-read,grep,write}")
  fi

  info "Running Kiro CLI headlessly for $FEATURE $stage revision"
  set +e
  "$KIRO_BIN" "${args[@]}" "$prompt_text" 2>&1 | tee "$log_file"
  local status=${PIPESTATUS[0]}
  set -e
  [[ $status -eq 0 ]] || fail "Kiro CLI revision failed for $stage. See $log_file"
  [[ -f "$expected_doc" ]] || fail "missing expected Kiro output after revision: $expected_doc"
}

run_stage() {
  local stage="$1" file_target="$2"
  info "Generating $stage with Kiro CLI headless mode: $KIRO_MODE"
  ./scripts/kiro-stage.sh --feature "$FEATURE" --stage "$stage" --mode "$KIRO_MODE"
  [[ -f "$file_target" ]] || fail "missing expected Kiro output: $file_target"

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
      fail "$stage review blocked or failed"
    fi
  done
}

run_stage requirements "$SPEC_PATH/requirements.md"
run_stage design "$SPEC_PATH/design.md"
run_stage tasks "$SPEC_PATH/tasks.md"
info "Spec flow complete. Status:"
./scripts/feature-state.py get --feature "$FEATURE"
