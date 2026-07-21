#!/usr/bin/env bash
set -euo pipefail

HANDOFF="${1:-${HANDOFF:-}}"

if [[ -z "$HANDOFF" ]]; then
  echo "ERROR: HANDOFF file is required"
  echo "Usage: ./scripts/architecture-handoff-check.sh docs/reviews/architecture-decision-handoffs/ADH-YYYY-NNN.md"
  exit 1
fi

if [[ ! -f "$HANDOFF" ]]; then
  echo "ERROR: handoff file not found: $HANDOFF"
  exit 1
fi

fail() {
  echo "FAIL: $1"
  exit 1
}

pass() {
  echo "PASS: $1"
}

require_heading() {
  local heading="$1"
  grep -qiE "^##[[:space:]]+$heading[[:space:]]*$" "$HANDOFF" || fail "Missing section: $heading"
  pass "Section present: $heading"
}

require_text() {
  local pattern="$1"
  local label="$2"
  grep -qiE "$pattern" "$HANDOFF" || fail "$label missing"
  pass "$label present"
}

echo "==> Validating Architecture Decision Handoff: $HANDOFF"

require_heading "Metadata"
require_heading "Decision title"
require_heading "Summary"
require_heading "Classification"
require_heading "Existing approved baseline"
require_heading "Decision or proposed decision"
require_heading "Rationale"
require_heading "Reuse-before-build assessment"
require_heading "Phase impact"
require_heading "Conflict check"
require_heading "Required action"
require_heading "Impacted files"
require_heading "Impacted features"
require_heading "Acceptance criteria for Kiro update"
require_heading "Explicit instructions to Kiro"
require_heading "Human approval"

require_text "Approval status:[[:space:]]*(Approved|Proposed|Rejected|Deferred)" "Approval status"
require_text "Classification" "Classification"
require_text "Reuse|Wrap|Extend|Build" "Reuse-before-build decision vocabulary"
require_text "Phase" "Phase impact"
require_text "DEC|RFC|Architecture Change Request|ACR|No repo change|Update architecture doc|Open Question" "Required action signal"

if grep -qiE "Approval status:[[:space:]]*Approved" "$HANDOFF"; then
  pass "Handoff is approved for Kiro validation"
else
  echo "WARN: handoff is not marked Approved; Kiro must not apply it without human approval"
fi

echo

echo "SUCCESS: handoff structure is valid"
