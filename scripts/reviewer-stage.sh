#!/opt/homebrew/bin/bash
set -euo pipefail
source "$(dirname "$0")/common.sh"
FEATURE=""; STAGE=""; MODE="${FEATURE_FACTORY_REVIEW_MODE:-auto}"
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
cd "$(repo_root)"; ensure_feature_state "$FEATURE"
SPEC_PATH=$(get_feature_value "$FEATURE" spec_path)
TITLE=$(get_feature_value "$FEATURE" title)
OUT_DIR=".automation/reviews/$FEATURE"; mkdir -p "$OUT_DIR"
PROMPT_OUT="$OUT_DIR/${STAGE}-review.prompt.md"
REVIEW_OUT="$OUT_DIR/${STAGE}.review.json"
RAW_OUT="$OUT_DIR/${STAGE}.openai.raw.json"
case "$STAGE" in
  requirements) TARGET="$SPEC_PATH/requirements.md";;
  design) TARGET="$SPEC_PATH/design.md";;
  tasks) TARGET="$SPEC_PATH/tasks.md";;
  *) fail "stage must be requirements, design, or tasks";;
esac
[[ -f "$TARGET" ]] || fail "missing file to review: $TARGET"
python3 - "$FEATURE" "$TITLE" "$STAGE" "$TARGET" "$PROMPT_OUT" <<'PYREVIEW'
from pathlib import Path
import sys
feature, title, stage, target, out = sys.argv[1:]
template_path = Path('docs/prompts/reviewer/approval-review.prompt.md')
if not template_path.exists():
    template_path = Path('docs/prompts/reviewer/spec-review.prompt.md')
template = template_path.read_text()
content = Path(target).read_text()
rendered = (template
    .replace('{{FEATURE_ID}}', feature)
    .replace('{{TITLE}}', title)
    .replace('{{STAGE}}', stage)
    .replace('{{TARGET_PATH}}', target)
    .replace('{{DOCUMENT_CONTENT}}', content))
Path(out).write_text(rendered)
print(out)
PYREVIEW
info "Review prompt generated: $PROMPT_OUT"
case "$MODE" in
  prompt)
    info "Prompt mode: paste $PROMPT_OUT into ChatGPT/reviewer, then save strict JSON to $REVIEW_OUT."
    ;;
  auto)
    REVIEWER_CMD="${FEATURE_FACTORY_REVIEWER_CMD:-./scripts/reviewer-openai.py}"
    info "Auto review mode: running $REVIEWER_CMD"
    "$REVIEWER_CMD" --prompt "$PROMPT_OUT" --out "$REVIEW_OUT" --raw-out "$RAW_OUT"
    info "Review JSON written: $REVIEW_OUT"
    ;;
  *)
    fail "--mode must be prompt or auto"
    ;;
esac
