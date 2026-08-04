#!/opt/homebrew/bin/bash
set -euo pipefail

fail() { echo "ERROR: $*" >&2; exit 1; }
info() { echo "==> $*"; }
repo_root() { git rev-parse --show-toplevel 2>/dev/null || fail "not inside a git repository"; }
require_clean_tree() {
  if ! git diff --quiet || ! git diff --cached --quiet; then
    fail "working tree has uncommitted tracked changes"
  fi
}
feature_branch_from_id_slug() {
  local feature="$1" slug="$2" lower
  lower=$(echo "$feature" | tr '[:upper:]' '[:lower:]')
  echo "${lower}-${slug}"
}
state_file() { echo ".automation/state/${1}.json"; }
get_feature_value() { ./scripts/feature-state.py get-value --feature "$1" --key "$2"; }
ensure_feature_state() { test -f "$(state_file "$1")" || fail "missing feature state for $1; run feature-start first"; }

configure_reviewer_adapter() {
  if [[ -n "${FEATURE_FACTORY_REVIEWER_CMD:-}" ]]; then
    case "$FEATURE_FACTORY_REVIEWER_CMD" in
      *reviewer-openai.py) FEATURE_FACTORY_REVIEWER_RAW_SUFFIX="${FEATURE_FACTORY_REVIEWER_RAW_SUFFIX:-openai.raw.json}";;
      *reviewer-codex.sh) FEATURE_FACTORY_REVIEWER_RAW_SUFFIX="${FEATURE_FACTORY_REVIEWER_RAW_SUFFIX:-codex.raw.json}";;
      *) FEATURE_FACTORY_REVIEWER_RAW_SUFFIX="${FEATURE_FACTORY_REVIEWER_RAW_SUFFIX:-reviewer.raw.json}";;
    esac
    export FEATURE_FACTORY_REVIEWER_CMD FEATURE_FACTORY_REVIEWER_RAW_SUFFIX
    return
  fi

  case "${OPENAI_API_KEY:-}" in
    ""|"..."|"<"*|YOUR_*|REPLACE_*) REVIEWER_OPENAI_KEY_AVAILABLE=0;;
    *) REVIEWER_OPENAI_KEY_AVAILABLE=1;;
  esac
  if [[ "$REVIEWER_OPENAI_KEY_AVAILABLE" == "1" ]]; then
    FEATURE_FACTORY_REVIEWER_CMD="./scripts/reviewer-openai.py"
    FEATURE_FACTORY_REVIEWER_RAW_SUFFIX="openai.raw.json"
  elif command -v codex >/dev/null 2>&1; then
    FEATURE_FACTORY_REVIEWER_CMD="./scripts/reviewer-codex.sh"
    FEATURE_FACTORY_REVIEWER_RAW_SUFFIX="codex.raw.json"
  else
    fail "no reviewer adapter available: configure a real OPENAI_API_KEY, install/sign in to Codex CLI, or set FEATURE_FACTORY_REVIEWER_CMD"
  fi
  export FEATURE_FACTORY_REVIEWER_CMD FEATURE_FACTORY_REVIEWER_RAW_SUFFIX
}
