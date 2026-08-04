#!/opt/homebrew/bin/bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

FEATURE=""
STAGE=""
MODE="${FEATURE_FACTORY_KIRO_MODE:-auto}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --feature) FEATURE="$2"; shift 2;;
    --stage) STAGE="$2"; shift 2;;
    --mode) MODE="$2"; shift 2;;
    *) fail "unknown arg: $1";;
  esac
done

[[ -n "$FEATURE" ]] || fail "--feature required"
[[ -n "$STAGE" ]] || fail "--stage required"
case "$STAGE" in requirements|design|tasks) ;; *) fail "unsupported stage: $STAGE";; esac
case "$MODE" in prompt|manual|auto) ;; *) fail "unsupported mode: $MODE";; esac

cd "$(repo_root)"
ensure_feature_state "$FEATURE"
CONTROL_FILE=".automation/features/${FEATURE}.control.json"

if [[ -f "$CONTROL_FILE" ]]; then
  PYTHONDONTWRITEBYTECODE=1 python3 \
    ./scripts/generic-kiro-boundary-check.py \
    --feature "$FEATURE" --stage "$STAGE" --mode pre
fi

if [[ "$FEATURE" == "FEATURE-0013" ]]; then
  if [[ -n "${KIRO_AGENT:-}" && "${KIRO_AGENT}" != "sovrunn-spec" ]]; then
    fail "FEATURE-0013 requires KIRO_AGENT=sovrunn-spec; refusing override: ${KIRO_AGENT}"
  fi
  KIRO_AGENT="sovrunn-spec"
  PYTHONDONTWRITEBYTECODE=1 python3 \
    ./scripts/feature-0013-architecture-boundary-check.py \
    --feature "$FEATURE" --stage "$STAGE" --mode pre
fi

if [[ "$FEATURE" == "FEATURE-0014" ]]; then
  PYTHONDONTWRITEBYTECODE=1 python3 \
    ./scripts/feature-0014-boundary-check.py \
    --feature "$FEATURE" --stage "$STAGE" --mode pre
fi

# Render the prompt first. render-prompt.py prints the generated file path.
PROMPT_PATH="$(./scripts/render-prompt.py --feature "$FEATURE" --stage "$STAGE" | tail -n 1)"
[[ -f "$PROMPT_PATH" ]] || fail "generated prompt not found: $PROMPT_PATH"

if [[ -f "$CONTROL_FILE" ]]; then
  PYTHONDONTWRITEBYTECODE=1 python3 \
    ./scripts/generic-kiro-boundary-check.py \
    --feature "$FEATURE" --stage "$STAGE" --mode prompt
fi

if [[ "$FEATURE" == "FEATURE-0014" ]]; then
  PYTHONDONTWRITEBYTECODE=1 python3 \
    ./scripts/feature-0014-boundary-check.py \
    --feature "$FEATURE" --stage "$STAGE" --mode prompt
fi

SPEC_PATH="$(get_feature_value "$FEATURE" spec_path)"
EXPECTED_DOC="$SPEC_PATH/$STAGE.md"
if [[ "$STAGE" == "revision" ]]; then
  EXPECTED_DOC="$SPEC_PATH/requirements.md"
fi

./scripts/feature-state.py set --feature "$FEATURE" --key current_stage --value "$STAGE" >/dev/null
./scripts/feature-state.py set --feature "$FEATURE" --key status --value "kiro_prompt_generated" >/dev/null

if [[ "$MODE" == "prompt" || "$MODE" == "manual" ]]; then
  info "Kiro prompt generated: $PROMPT_PATH"
  info "Prompt/manual mode: paste this prompt into Kiro. Default is auto/headless; set FEATURE_FACTORY_KIRO_MODE=prompt only when you intentionally want manual mode."
  exit 0
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
# Some Kiro installations use logged-in local credentials instead of KIRO_API_KEY.
if [[ -z "${KIRO_API_KEY:-}" && "${KIRO_ALLOW_LOCAL_AUTH:-1}" != "1" ]]; then
  fail "KIRO_API_KEY is required when KIRO_ALLOW_LOCAL_AUTH is not enabled"
fi

mkdir -p ".automation/logs/$FEATURE"
LOG_FILE=".automation/logs/$FEATURE/kiro-${STAGE}.log"
MODEL_JSON=".automation/logs/$FEATURE/kiro-${STAGE}.model.json"
./scripts/model-recommend.py --tool kiro --stage "$STAGE" --format json > "$MODEL_JSON"

read -r REC_MODEL REC_LABEL REC_EFFORT < <(python3 - "$MODEL_JSON" <<'PY'
import json, sys
j=json.load(open(sys.argv[1]))
r=j['recommendations'][0]
print(r.get('model',''), r.get('label',''), r.get('effort','medium'))
PY
)

KIRO_MODEL="${KIRO_MODEL:-}"
if [[ -f "$CONTROL_FILE" && -z "$KIRO_MODEL" ]]; then
  KIRO_MODEL="$REC_MODEL"
fi
if [[ "$FEATURE" == "FEATURE-0013" && -z "$KIRO_MODEL" ]]; then
  KIRO_MODEL="claude-opus-4.8"
fi
SELECTED_MODEL="${KIRO_SELECTED_MODEL:-${KIRO_MODEL:-$REC_LABEL}}"
EFFORT_RAW="${KIRO_EFFORT:-$REC_EFFORT}"
EFFORT="$(echo "$EFFORT_RAW" | tr '[:upper:]' '[:lower:]')"
case "$EFFORT" in low|medium|high|xhigh|max) ;; *) EFFORT="high";; esac

PROMPT_TEXT="$(cat "$PROMPT_PATH")"

KIRO_ARGS=(chat --no-interactive --effort "$EFFORT")
if [[ -n "$KIRO_MODEL" ]]; then
  KIRO_ARGS+=(--model "$KIRO_MODEL")
fi
if [[ -n "${KIRO_AGENT:-}" ]]; then
  KIRO_ARGS+=(--agent "$KIRO_AGENT")
fi
if [[ "${KIRO_TRUST_ALL_TOOLS:-0}" == "1" ]]; then
  KIRO_ARGS+=(--trust-all-tools)
else
  KIRO_ARGS+=(--trust-tools="${KIRO_TRUST_TOOLS:-read,grep,write}")
fi

info "Running Kiro CLI for $FEATURE $STAGE"
info "Recommended model priority is embedded in the prompt. Recording selected model as: $SELECTED_MODEL"
info "Kiro CLI effort: $EFFORT"
set +e
"$KIRO_BIN" "${KIRO_ARGS[@]}" "$PROMPT_TEXT" 2>&1 | tee "$LOG_FILE"
STATUS=${PIPESTATUS[0]}
set -e

./scripts/model-record.sh \
  --feature "$FEATURE" \
  --tool kiro \
  --stage "$STAGE" \
  --selected-model "$SELECTED_MODEL" \
  --effort "$EFFORT" \
  --fallback-used "${KIRO_FALLBACK_USED:-no}" \
  --fallback-reason "${KIRO_FALLBACK_REASON:-none}" >/dev/null || true

if [[ $STATUS -ne 0 ]]; then
  ./scripts/feature-state.py set --feature "$FEATURE" --key status --value "kiro_${STAGE}_failed" >/dev/null || true
  fail "Kiro CLI failed for $STAGE. See $LOG_FILE"
fi

if [[ -f "$CONTROL_FILE" ]]; then
  if ! ./scripts/receipt-check.py --log "$LOG_FILE" --document "$EXPECTED_DOC" --kind stage; then
    ./scripts/feature-state.py set --feature "$FEATURE" --key status --value "kiro_${STAGE}_blocked" >/dev/null || true
    fail "Kiro stage did not produce exactly one COMPLETE receipt; refusing stage acceptance. See $LOG_FILE"
  fi
fi

if [[ ! -f "$EXPECTED_DOC" ]]; then
  echo "WARNING: expected Kiro output file not found: $EXPECTED_DOC" >&2
  echo "Check Kiro log: $LOG_FILE" >&2
else
  info "Kiro output detected: $EXPECTED_DOC"
fi

if [[ "$FEATURE" == "FEATURE-0013" ]]; then
  PYTHONDONTWRITEBYTECODE=1 python3 \
    ./scripts/feature-0013-architecture-boundary-check.py \
    --feature "$FEATURE" --stage "$STAGE" --mode post
fi

if [[ "$FEATURE" == "FEATURE-0014" ]]; then
  PYTHONDONTWRITEBYTECODE=1 python3 \
    ./scripts/feature-0014-boundary-check.py \
    --feature "$FEATURE" --stage "$STAGE" --mode post
fi

if [[ -f "$CONTROL_FILE" ]]; then
  PYTHONDONTWRITEBYTECODE=1 python3 \
    ./scripts/generic-kiro-boundary-check.py \
    --feature "$FEATURE" --stage "$STAGE" --mode post
fi

if [[ -f "$CONTROL_FILE" && "${FEATURE_FACTORY_DEFER_SEMANTIC_CHECK:-0}" != "1" ]]; then
  SEMANTIC_REVISION_PROMPT=".automation/generated-prompts/$FEATURE/${STAGE}.semantic-revision.prompt.md"
  if ! ./scripts/kiro-semantic-check.py \
    --feature "$FEATURE" \
    --stage "$STAGE" \
    --write-revision-prompt "$SEMANTIC_REVISION_PROMPT"; then
    ./scripts/feature-state.py set --feature "$FEATURE" --key status --value "${STAGE}_revision_required" >/dev/null || true
    fail "Kiro semantic guardrails rejected $STAGE. Apply $SEMANTIC_REVISION_PROMPT, then rerun the stage."
  fi
fi

./scripts/feature-state.py set --feature "$FEATURE" --key status --value "${STAGE}_generated" >/dev/null
info "Kiro stage completed. Log: $LOG_FILE"
