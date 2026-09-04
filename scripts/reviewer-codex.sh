#!/opt/homebrew/bin/bash
set -euo pipefail

source "$(dirname "$0")/common.sh"

PROMPT=""
OUT=""
RAW_OUT=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --prompt) PROMPT="$2"; shift 2;;
    --out) OUT="$2"; shift 2;;
    --raw-out) RAW_OUT="$2"; shift 2;;
    *) fail "unknown arg: $1";;
  esac
done
[[ -n "$PROMPT" ]] || fail "--prompt required"
[[ -n "$OUT" ]] || fail "--out required"
[[ -n "$RAW_OUT" ]] || fail "--raw-out required"

cd "$(repo_root)"
[[ -f "$PROMPT" ]] || fail "review prompt does not exist: $PROMPT"
SCHEMA=".automation/schemas/reviewer-response.schema.json"
[[ -f "$SCHEMA" ]] || fail "reviewer response schema does not exist: $SCHEMA"

CODEX_BIN="$(resolve_codex_bin || true)"
[[ -n "$CODEX_BIN" && -x "$CODEX_BIN" ]] || \
  fail "Codex CLI is unavailable; install/sign in to Codex CLI/Desktop or configure FEATURE_FACTORY_REVIEWER_CMD"

mkdir -p "$(dirname "$OUT")" "$(dirname "$RAW_OUT")"
TMP_OUT="$(mktemp "${TMPDIR:-/tmp}/sovrunn-codex-review.XXXXXX")"
trap 'rm -f "$TMP_OUT"' EXIT

read -r FEATURE STAGE < <(python3 - "$PROMPT" <<'PY'
import re, sys
from pathlib import Path
text = Path(sys.argv[1]).read_text()
feature = re.search(r"(?m)^Feature:\s*\n\s*(FEATURE-\d{4})\b", text)
stage = re.search(r"(?m)^Stage:\s*\n\s*(requirements|design|tasks)\s*$", text)
if not feature or not stage:
    raise SystemExit("ERROR: reviewer prompt identity is missing")
print(feature.group(1), stage.group(1))
PY
)
LOG_DIR=".automation/logs/${FEATURE}"
LOG_FILE="${LOG_DIR}/codex-${STAGE}-review.log"
mkdir -p "$LOG_DIR"

CODEX_ARGS=(
  exec
  --ephemeral
  --sandbox read-only
  --output-schema "$SCHEMA"
  --output-last-message "$TMP_OUT"
  --color never
  --cd "$PWD"
)
if [[ -n "${CODEX_REVIEWER_MODEL:-}" ]]; then
  CODEX_ARGS+=(--model "$CODEX_REVIEWER_MODEL")
fi
if [[ -n "${CODEX_REVIEWER_REASONING_EFFORT:-}" ]]; then
  CODEX_ARGS+=(
    --config
    "model_reasoning_effort=\"${CODEX_REVIEWER_REASONING_EFFORT}\""
  )
fi

info "Running independent Codex CLI review for $FEATURE $STAGE (read-only sandbox)"
set +e
env -u OPENAI_API_KEY "$CODEX_BIN" "${CODEX_ARGS[@]}" - < "$PROMPT" 2>&1 | tee "$LOG_FILE"
CODEX_STATUS=${PIPESTATUS[0]}
set -e
[[ "$CODEX_STATUS" -eq 0 ]] || fail "Codex reviewer failed. See $LOG_FILE"

python3 - "$PROMPT" "$TMP_OUT" "$OUT" "$RAW_OUT" <<'PY'
import importlib.util
import json
import sys
from pathlib import Path

prompt_path, raw_path, out_path, audit_path = map(Path, sys.argv[1:])
module_path = Path("scripts/reviewer-openai.py")
spec = importlib.util.spec_from_file_location("reviewer_contract", module_path)
module = importlib.util.module_from_spec(spec)
assert spec.loader
spec.loader.exec_module(module)
stage = module.extract_stage(prompt_path.read_text())
try:
    review = json.loads(raw_path.read_text())
except json.JSONDecodeError as exc:
    raise SystemExit(f"ERROR: Codex reviewer output was not strict JSON: {exc}") from exc
review = module.normalize_review_routing(review)
module.validate_review(review, stage)
rendered = json.dumps(review, indent=2, sort_keys=True) + "\n"
out_path.write_text(rendered)
audit_path.write_text(rendered)
print(out_path)
PY

info "Codex review JSON written: $OUT"
