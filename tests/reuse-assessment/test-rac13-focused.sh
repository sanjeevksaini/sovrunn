#!/usr/bin/env bash
# Focused RA-C13 structured approval-evidence tests (FEATURE-0011)
# Each mismatch test starts from matching artifacts and introduces one defect.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
export SOVRUNN_ROOT="$ROOT"
FIX="$ROOT/tests/reuse-assessment/fixtures"
CANON_EVIDENCE="$ROOT/docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md"
PASS=0
FAIL=0

assert_pass() {
  PASS=$((PASS + 1))
  echo "PASS: $1"
}

assert_fail_msg() {
  FAIL=$((FAIL + 1))
  echo "FAIL: $1"
  echo "  $2"
}

setup_repo() {
  local dir="$1"
  mkdir -p "$dir"
  cd "$dir"
  git -c init.templateDir= -c core.hooksPath=/dev/null init >/dev/null 2>&1
  git config user.name "Harness"
  git config user.email "harness@sovrunn.test"
  git config core.hooksPath /dev/null

  mkdir -p docs/phase2 docs/features docs/decisions docs/rfc \
    docs/reviews/architecture-decision-handoffs docs/reviews/reuse-assessments scripts

  cp "$ROOT/docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md" docs/phase2/
  cp "$ROOT/scripts/reuse-assessment-check.sh" scripts/
  chmod +x scripts/*.sh

  echo "# DEC-0026" > docs/decisions/DEC-0026-reuse-before-build.md
  echo "# DEC-0036" > docs/decisions/DEC-0036-adapter-boundaries.md
  echo "# RFC-0021" > docs/rfc/RFC-0021-reuse-first-architecture.md
  cat > docs/reviews/architecture-decision-handoffs/ADH-2026-011-feature-0011-reuse-assessment-standard.md <<'ADH'
# Architecture Decision Handoff
- **Related feature:** FEATURE-0011
- **Approval status:** Approved
ADH

  cp "$FIX/valid/complete-assessment.md" docs/features/FEATURE-0011-reuse-assessment-standard.md
  # Assessment path in fixture must match evidence Assessment artifact
  # complete-assessment already references FEATURE-0011 path and evidence file.
  cp "$CANON_EVIDENCE" docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md

  mkdir -p scripts/impl
  echo '#!/bin/bash' > scripts/impl/test-tool.sh
  printf '%s\n' "scripts/impl/test-tool.sh" > changed-files.txt

  git add -A >/dev/null 2>&1
  git commit -q -m "base"
}

run_check() {
  bash scripts/reuse-assessment-check.sh FEATURE-0011 \
    --assessment docs/features/FEATURE-0011-reuse-assessment-standard.md \
    --changed-files changed-files.txt --skip-rac03 2>&1
}

expect_pass() {
  local name="$1"
  local out rc
  set +e
  out="$(run_check)"
  rc=$?
  set -e
  if [[ "$rc" -eq 0 ]] && ! grep -q '|CONFIG|' <<<"$out"; then
    assert_pass "$name"
  else
    assert_fail_msg "$name" "expected exit 0, got $rc: $(echo "$out" | head -3)"
  fi
}

expect_rac13() {
  local name="$1"
  local out rc
  set +e
  out="$(run_check)"
  rc=$?
  set -e
  if [[ "$rc" -eq 1 ]] && grep -q '|RA-C13|' <<<"$out" && ! grep -q '|CONFIG|' <<<"$out"; then
    assert_pass "$name"
  else
    assert_fail_msg "$name" "expected RA-C13 exit 1, got $rc: $(echo "$out" | head -5)"
  fi
}

expect_config() {
  local name="$1"
  local out rc
  set +e
  out="$(run_check)"
  rc=$?
  set -e
  if [[ "$rc" -eq 2 ]] && grep -q '|CONFIG|' <<<"$out"; then
    assert_pass "$name"
  else
    assert_fail_msg "$name" "expected CONFIG exit 2, got $rc: $(echo "$out" | head -5)"
  fi
}

mutate_evidence() {
  local field="$1"
  local new_value="$2"
  python3 - "$field" "$new_value" <<'PY'
import re, sys
from pathlib import Path
field, new_value = sys.argv[1], sys.argv[2]
path = Path("docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md")
text = path.read_text(encoding="utf-8")
patterns = [
    re.compile(rf"(^[\s*]*\*?\s*{re.escape(field)}\s*:\s*).+$", re.IGNORECASE | re.MULTILINE),
]
snake = field.lower().replace(" ", "_").replace("/", "_")
patterns.append(re.compile(rf"(^{re.escape(snake)}\s*:\s*).+$", re.IGNORECASE | re.MULTILINE))
# Front-matter aliases used by the approved evidence document
aliases = {
    "Approver or approving role": ["approving_role"],
    "Approval date": ["approval_date"],
    "Approval status": ["approval_status"],
    "Feature": ["feature"],
    "Assessment format version": ["assessment_format_version"],
}
for alias in aliases.get(field, []):
    patterns.append(re.compile(rf"(^{re.escape(alias)}\s*:\s*).+$", re.IGNORECASE | re.MULTILINE))
n_total = 0
for pat in patterns:
    text, n = pat.subn(rf"\g<1>{new_value}", text)
    n_total += n
if n_total == 0:
    raise SystemExit(f"field not found for mutation: {field}")
path.write_text(text, encoding="utf-8")
PY
}

remove_evidence_field_line() {
  local field="$1"
  python3 - "$field" <<'PY'
import re, sys
from pathlib import Path
field = sys.argv[1]
path = Path("docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md")
text = path.read_text(encoding="utf-8")
# Remove only the body list field; keep front matter intact so FM/body conflict or
# missing body field yields RA-C13 (not CONFIG from missing FM keys).
end = text.find("\n---\n", 3)
if end < 0:
    raise SystemExit("front matter not found")
head, body = text[: end + 5], text[end + 5 :]
body2, n = re.subn(
    rf"(?m)^\s*\*\s*{re.escape(field)}\s*:.*\n?",
    "",
    body,
    count=1,
)
if n != 1:
    raise SystemExit(f"body field not found: {field}")
path.write_text(head + body2, encoding="utf-8")
PY
}

echo "==> Focused RA-C13 structured evidence tests"

RESULTS="$(mktemp)"
trap 'rm -f "$RESULTS"' EXIT

record() {
  echo "$1" >>"$RESULTS"
}

# 1. Exact matching evidence passes
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 0 ]] && ! grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: exact matching evidence passes"
    record PASS
  else
    echo "FAIL: exact matching evidence passes (exit $rc)"
    record FAIL
    exit 1
  fi
)

run_mismatch() {
  local name="$1"
  local field="$2"
  local value="$3"
  (
    repo="$(mktemp -d)"
    trap "rm -rf '$repo'" EXIT
    setup_repo "$repo"
    mutate_evidence "$field" "$value"
    set +e
    out="$(run_check)"; rc=$?; set -e
    if [[ "$rc" -eq 1 ]] && grep -q '|RA-C13|' <<<"$out" && ! grep -q '|CONFIG|' <<<"$out"; then
      echo "PASS: $name"
      record PASS
    else
      echo "FAIL: $name (exit $rc): $(echo "$out" | head -3)"
      record FAIL
      exit 1
    fi
  )
}

run_mismatch "mismatched feature fails RA-C13" "Feature" "FEATURE-0012"
run_mismatch "mismatched assessment path fails RA-C13" "Assessment artifact" "docs/features/FEATURE-0012-other.md"
run_mismatch "mismatched format version fails RA-C13" "Assessment format version" "9.9.9"
run_mismatch "mismatched disposition fails RA-C13" "Disposition" "Reuse"
run_mismatch "mismatched Sovrunn-owned responsibility fails RA-C13" "Sovrunn-owned responsibility" "Completely different Sovrunn-owned responsibility text."
run_mismatch "mismatched reused/extended responsibility fails RA-C13" "Reused or extended responsibility" "Completely different reused responsibility text."
run_mismatch "mismatched responsibility/control boundary fails RA-C13" "Responsibility/control boundary" "Completely different control boundary text."
run_mismatch "unapproved evidence fails RA-C13" "Approval status" "Proposed"
run_mismatch "mismatched approver fails RA-C13" "Approver or approving role" "Different Approver Role"
run_mismatch "mismatched approval date fails RA-C13" "Approval date" "2099-01-01"

set_assessment_evidence_ref() {
  local new_ref="$1"
  python3 - "$new_ref" <<'PY'
import sys
from pathlib import Path
new_ref = sys.argv[1]
p = Path("docs/features/FEATURE-0011-reuse-assessment-standard.md")
text = p.read_text(encoding="utf-8")
# Replace structured evidence table cell and prose references
old = "docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md"
if old not in text:
    raise SystemExit("evidence ref not found in assessment")
text = text.replace(old, new_ref)
p.write_text(text, encoding="utf-8")
PY
}

# Exact authoritative path acceptance (already covered by exact matching; assert explicitly)
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 0 ]]; then
    echo "PASS: exact active-feature evidence path passes"
    record PASS
  else
    echo "FAIL: exact path acceptance (exit $rc)"
    record FAIL
    exit 1
  fi
)

# Alternate relative filename
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  mkdir -p docs/reviews/reuse-assessments
  cp docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md \
    docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence-alt.md
  set_assessment_evidence_ref "docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence-alt.md"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 1 ]] && grep -q '|RA-C13|' <<<"$out" && ! grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: alternate relative evidence path rejected RA-C13"
    record PASS
  else
    echo "FAIL: alternate relative path (exit $rc): $(echo "$out" | head -3)"
    record FAIL
    exit 1
  fi
)

# Cross-feature evidence path
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  mkdir -p docs/reviews/reuse-assessments
  cp docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md \
    docs/reviews/reuse-assessments/FEATURE-0012-approval-evidence.md
  set_assessment_evidence_ref "docs/reviews/reuse-assessments/FEATURE-0012-approval-evidence.md"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 1 ]] && grep -q '|RA-C13|' <<<"$out" && ! grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: cross-feature evidence path rejected RA-C13"
    record PASS
  else
    echo "FAIL: cross-feature path (exit $rc): $(echo "$out" | head -3)"
    record FAIL
    exit 1
  fi
)

# Absolute path => CONFIG exit 2
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  abs="$repo/docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md"
  set_assessment_evidence_ref "$abs"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 2 ]] && grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: absolute evidence path rejected CONFIG exit 2"
    record PASS
  else
    echo "FAIL: absolute path (exit $rc): $(echo "$out" | head -3)"
    record FAIL
    exit 1
  fi
)

# .. traversal => CONFIG exit 2
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  set_assessment_evidence_ref "docs/reviews/reuse-assessments/../reuse-assessments/FEATURE-0011-approval-evidence.md"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 2 ]] && grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: traversal evidence path rejected CONFIG exit 2"
    record PASS
  else
    echo "FAIL: traversal path (exit $rc): $(echo "$out" | head -3)"
    record FAIL
    exit 1
  fi
)

# Path resolving outside repository => CONFIG exit 2
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  # Symlink-style outside path via absolute /tmp reference already covered;
  # use explicit absolute outside path.
  outside="$(mktemp)"
  cp docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md "$outside"
  set_assessment_evidence_ref "$outside"
  set +e
  out="$(run_check)"; rc=$?; set -e
  rm -f "$outside"
  if [[ "$rc" -eq 2 ]] && grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: outside-repository evidence path rejected CONFIG exit 2"
    record PASS
  else
    echo "FAIL: outside path (exit $rc): $(echo "$out" | head -3)"
    record FAIL
    exit 1
  fi
)

mutate_fm_only() {
  local key="$1"
  local value="$2"
  python3 - "$key" "$value" <<'PY'
import re, sys
from pathlib import Path
key, value = sys.argv[1], sys.argv[2]
path = Path("docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md")
text = path.read_text(encoding="utf-8")
if not text.startswith("---\n"):
    raise SystemExit("no front matter")
end = text.find("\n---\n", 3)
fm, body = text[4:end], text[end + 5 :]
fm2, n = re.subn(rf"(?m)^{re.escape(key)}\s*:.*$", f"{key}: {value}", fm, count=1)
if n != 1:
    raise SystemExit(f"front-matter key not found: {key}")
path.write_text("---\n" + fm2 + "\n---\n" + body, encoding="utf-8")
PY
}

mutate_body_only() {
  local field="$1"
  local value="$2"
  python3 - "$field" "$value" <<'PY'
import re, sys
from pathlib import Path
field, value = sys.argv[1], sys.argv[2]
path = Path("docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md")
text = path.read_text(encoding="utf-8")
# Only mutate body list field after closing front matter
end = text.find("\n---\n", 3)
fm, body = text[: end + 5], text[end + 5 :]
body2, n = re.subn(
    rf"(?m)^(\s*\*\s*{re.escape(field)}\s*:\s*).+$",
    rf"\g<1>{value}",
    body,
    count=1,
)
if n != 1:
    raise SystemExit(f"body field not found: {field}")
path.write_text(fm + body2, encoding="utf-8")
PY
}

# Matching front matter and body already covered by exact match suite entry

# Conflicting approval status (FM vs body)
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  mutate_fm_only "approval_status" "Proposed"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 1 ]] && grep -q '|RA-C13|' <<<"$out" && ! grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: conflicting approval status FM/body fails RA-C13"
    record PASS
  else
    echo "FAIL: conflicting approval status (exit $rc): $(echo "$out" | head -3)"
    record FAIL
    exit 1
  fi
)

# Conflicting approver
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  mutate_fm_only "approving_role" "Other Role"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 1 ]] && grep -q '|RA-C13|' <<<"$out" && ! grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: conflicting approver FM/body fails RA-C13"
    record PASS
  else
    echo "FAIL: conflicting approver (exit $rc): $(echo "$out" | head -3)"
    record FAIL
    exit 1
  fi
)

# Conflicting approval date
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  mutate_fm_only "approval_date" "2099-01-01"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 1 ]] && grep -q '|RA-C13|' <<<"$out" && ! grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: conflicting approval date FM/body fails RA-C13"
    record PASS
  else
    echo "FAIL: conflicting approval date (exit $rc): $(echo "$out" | head -3)"
    record FAIL
    exit 1
  fi
)

# Conflicting feature
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  mutate_fm_only "feature" "FEATURE-0012"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 1 ]] && grep -q '|RA-C13|' <<<"$out" && ! grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: conflicting feature FM/body fails RA-C13"
    record PASS
  else
    echo "FAIL: conflicting feature (exit $rc): $(echo "$out" | head -3)"
    record FAIL
    exit 1
  fi
)

# Conflicting format version
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  mutate_fm_only "assessment_format_version" "9.9.9"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 1 ]] && grep -q '|RA-C13|' <<<"$out" && ! grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: conflicting format version FM/body fails RA-C13"
    record PASS
  else
    echo "FAIL: conflicting format version (exit $rc): $(echo "$out" | head -3)"
    record FAIL
    exit 1
  fi
)

# Unterminated front matter => CONFIG exit 2
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  python3 - <<'PY'
from pathlib import Path
p = Path("docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md")
text = p.read_text(encoding="utf-8")
# Remove closing ---
text = text.replace("\n---\n\n#", "\n\n#", 1)
p.write_text(text, encoding="utf-8")
PY
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 2 ]] && grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: unterminated front matter CONFIG exit 2"
    record PASS
  else
    echo "FAIL: unterminated front matter (exit $rc): $(echo "$out" | head -3)"
    record FAIL
    exit 1
  fi
)

# Duplicate front-matter key => CONFIG exit 2
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  python3 - <<'PY'
from pathlib import Path
p = Path("docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md")
text = p.read_text(encoding="utf-8")
text = text.replace(
    "feature: FEATURE-0011\n",
    "feature: FEATURE-0011\nfeature: FEATURE-0011\n",
    1,
)
p.write_text(text, encoding="utf-8")
PY
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 2 ]] && grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: duplicate front-matter key CONFIG exit 2"
    record PASS
  else
    echo "FAIL: duplicate front-matter key (exit $rc): $(echo "$out" | head -3)"
    record FAIL
    exit 1
  fi
)

# Missing approver
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  remove_evidence_field_line "Approver or approving role"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 1 ]] && grep -q '|RA-C13|' <<<"$out" && ! grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: missing approver fails RA-C13"
    record PASS
  else
    echo "FAIL: missing approver (exit $rc)"
    record FAIL
    exit 1
  fi
)

# Missing approval date
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  remove_evidence_field_line "Approval date"
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 1 ]] && grep -q '|RA-C13|' <<<"$out" && ! grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: missing approval date fails RA-C13"
    record PASS
  else
    echo "FAIL: missing approval date (exit $rc)"
    record FAIL
    exit 1
  fi
)

# Approval details only in assessment
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  python3 - <<'PY'
from pathlib import Path
p = Path("docs/features/FEATURE-0011-reuse-assessment-standard.md")
text = p.read_text(encoding="utf-8")
text = text.replace("| Structured approval-evidence record | `docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md` |\n", "")
text = text.replace("docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md", "NO-EVIDENCE-REF")
p.write_text(text, encoding="utf-8")
PY
  rm -f docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 1 ]] && grep -q '|RA-C13|' <<<"$out" && ! grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: approval details only in assessment fail RA-C13"
    record PASS
  else
    echo "FAIL: approval-only-in-assessment (exit $rc)"
    record FAIL
    exit 1
  fi
)

# Missing evidence file exits 2
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  rm -f docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 2 ]] && grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: missing evidence file exits 2"
    record PASS
  else
    echo "FAIL: missing evidence file (exit $rc)"
    record FAIL
    exit 1
  fi
)

# Invalid UTF-8
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  printf 'feature: FEATURE-0011\n\xff\xfe invalid' > docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 2 ]] && grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: invalid UTF-8 evidence exits 2"
    record PASS
  else
    echo "FAIL: invalid UTF-8 (exit $rc)"
    record FAIL
    exit 1
  fi
)

# Empty evidence
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  : > docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md
  set +e
  out="$(run_check)"; rc=$?; set -e
  if [[ "$rc" -eq 2 ]] && grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: empty evidence exits 2"
    record PASS
  else
    echo "FAIL: empty evidence (exit $rc)"
    record FAIL
    exit 1
  fi
)

# Evidence artifact structure: front matter + required headings
(
  set +e
  out="$(python3 - <<'PY'
from pathlib import Path
import os
import sys
root = os.environ["SOVRUNN_ROOT"]
text = Path(root, "docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md").read_text(encoding="utf-8")
if not text.startswith("---\n"):
    print("missing opening front-matter delimiter")
    sys.exit(1)
end = text.find("\n---\n", 3)
if end < 0:
    print("missing closing front-matter delimiter")
    sys.exit(1)
fm = text[4:end]
required_fm = [
    "feature: FEATURE-0011",
    "evidence_type: reuse-assessment-approval",
    "approval_status: Approved",
    "approval_date: 2026-07-21",
    "approving_role: Sovrunn project owner",
    "assessment_format_version: 1.0.0",
]
for line in required_fm:
    if line not in fm:
        print(f"missing front-matter field: {line}")
        sys.exit(1)
required_headings = [
    "# FEATURE-0011 Reuse Assessment Approval Evidence",
    "## Evidence identity",
    "## Approved decision",
    "## Approved responsibility boundary",
    "## Controlling decision",
    "## Comparison contract",
    "## Status separation",
]
for h in required_headings:
    if h not in text:
        print(f"missing heading: {h}")
        sys.exit(1)
print("ok")
sys.exit(0)
PY
)"
  rc=$?
  set -e
  if [[ "$rc" -eq 0 ]]; then
    echo "PASS: maintained evidence artifact has valid front matter and required headings"
    record PASS
  else
    echo "FAIL: evidence artifact structure: $out"
    record FAIL
    exit 1
  fi
)

# Strict missing-changed-files
(
  repo="$(mktemp -d)"
  trap "rm -rf '$repo'" EXIT
  setup_repo "$repo"
  set +e
  out="$(bash scripts/reuse-assessment-check.sh FEATURE-0011 \
    --assessment docs/features/FEATURE-0011-reuse-assessment-standard.md \
    --mode strict --skip-rac03 2>&1)"
  rc=$?
  set -e
  if [[ "$rc" -eq 2 ]] && grep -q '|CONFIG|' <<<"$out"; then
    echo "PASS: strict missing-changed-files exits 2"
    record PASS
  else
    echo "FAIL: strict missing-changed-files (exit $rc)"
    record FAIL
    exit 1
  fi
)

PASS=$(grep -c '^PASS$' "$RESULTS" || true)
FAIL=$(grep -c '^FAIL$' "$RESULTS" || true)
echo
echo "Focused RA-C13 evidence tests: PASS=$PASS FAIL=$FAIL"
if [[ "$FAIL" -ne 0 ]]; then
  exit 1
fi
exit 0
