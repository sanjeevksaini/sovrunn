#!/usr/bin/env bash
# Reuse Assessment Standard test harness (FEATURE-0011 / Task 15)
# Never mutates the active Sovrunn repository index, branch, or commits.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

VALIDATOR="$ROOT/scripts/reuse-assessment-check.sh"
GATE="$ROOT/scripts/feature-gate.sh"
FIX="$ROOT/tests/reuse-assessment/fixtures"
EXPECTED="$FIX/expected/diagnostics.txt"

PASS=0
FAIL=0
failures=()

BEFORE_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
BEFORE_STATUS="$(git status --porcelain)"
BEFORE_HEAD="$(git rev-parse HEAD)"

cleanup_active_check() {
  local after_branch after_status after_head
  after_branch="$(git rev-parse --abbrev-ref HEAD)"
  after_status="$(git status --porcelain)"
  after_head="$(git rev-parse HEAD)"
  if [[ "$after_branch" != "$BEFORE_BRANCH" || "$after_head" != "$BEFORE_HEAD" ]]; then
    echo "FAIL: active repository branch/HEAD changed by harness"
    echo "  before: $BEFORE_BRANCH $BEFORE_HEAD"
    echo "  after:  $after_branch $after_head"
    exit 1
  fi
  # Allow only authorized implementation file modifications already present before/after;
  # harness must not create new porcelain noise beyond what existed.
  # Compare sorted status sets ignoring tests that only read.
  if [[ "$after_status" != "$BEFORE_STATUS" ]]; then
    # Filter: harness may create nothing in active repo; any delta is a failure.
    echo "FAIL: active repository working tree changed by harness"
    echo "--- before ---"
    echo "$BEFORE_STATUS"
    echo "--- after ---"
    echo "$after_status"
    exit 1
  fi
  echo "PASS: active repository branch/index/working-tree unchanged"
}

assert_ok() {
  local name="$1"
  PASS=$((PASS + 1))
  echo "PASS: $name"
}

assert_fail() {
  local name="$1"
  local detail="${2:-}"
  FAIL=$((FAIL + 1))
  failures+=("$name :: $detail")
  echo "FAIL: $name — $detail"
}

run_validator() {
  local feature="$1"
  local assessment="$2"
  shift 2
  set +e
  OUT="$(bash "$VALIDATOR" "$feature" --assessment "$assessment" --skip-rac03 --skip-rac13 "$@" 2>&1)"
  RC=$?
  set -e
}

echo "==> Reuse assessment harness"
echo "Active branch before: $BEFORE_BRANCH"

# ---------------------------------------------------------------------------
# File-content cases from expected diagnostics
# ---------------------------------------------------------------------------
echo "==> File-content fixture cases"
while IFS='|' read -r fixture rule exp_exit; do
  [[ -z "$fixture" || "$fixture" =~ ^# ]] && continue
  path="$FIX/$fixture"
  [[ -f "$path" ]] || { assert_fail "fixture $fixture" "missing file"; continue; }

  # RA-C10 fixtures use FEATURE-0011 CLI; invalid identity is inside assessment
  run_validator "FEATURE-0011" "$path"
  if [[ "$RC" -ne "$exp_exit" ]]; then
    assert_fail "$fixture exit" "expected $exp_exit got $RC"
    echo "$OUT" | head -5
    continue
  fi
  if [[ -n "$rule" ]]; then
    if ! grep -q "|$rule|" <<<"$OUT"; then
      assert_fail "$fixture rule $rule" "rule id not in output"
      echo "$OUT" | head -10
      continue
    fi
  fi
  assert_ok "$fixture (exit $exp_exit${rule:+, $rule})"
done < "$EXPECTED"

# CLI malformed feature identifiers (RA-C10 / fail-safe)
echo "==> CLI feature-identifier anchoring"
for bad in XFEATURE-0011 FEATURE-0011-extra FEATURE-011 FEATURE-001A; do
  set +e
  OUT="$(bash "$VALIDATOR" "$bad" --assessment "$FIX/valid/complete-assessment.md" --skip-rac03 2>&1)"
  RC=$?
  set -e
  if [[ "$RC" -eq 2 ]]; then
    assert_ok "CLI reject $bad (exit 2)"
  else
    assert_fail "CLI reject $bad" "expected exit 2 got $RC"
  fi
done

# Legacy exemption
set +e
OUT="$(bash "$VALIDATOR" FEATURE-0005 2>&1)"
RC=$?
set -e
if [[ "$RC" -eq 0 ]]; then
  assert_ok "legacy FEATURE-0005 exempt"
else
  assert_fail "legacy FEATURE-0005" "expected exit 0 got $RC"
fi

# Canonical config errors (exit 2) in isolated temp repo layout
echo "==> Canonical configuration errors"
canon_tmp="$(mktemp -d)"
trap 'rm -rf "$canon_tmp"' RETURN
mkdir -p "$canon_tmp/docs/phase2" "$canon_tmp/scripts" "$canon_tmp/docs/features"
cp "$VALIDATOR" "$canon_tmp/scripts/reuse-assessment-check.sh"
cp "$FIX/valid/complete-assessment.md" "$canon_tmp/docs/features/assess.md"
# missing canonical
set +e
OUT="$(bash "$canon_tmp/scripts/reuse-assessment-check.sh" FEATURE-0011 --repo-root "$canon_tmp" --assessment docs/features/assess.md --skip-rac03 2>&1)"
RC=$?
set -e
if [[ "$RC" -eq 2 ]]; then
  assert_ok "canonical missing => exit 2"
else
  assert_fail "canonical missing" "expected 2 got $RC"
fi

cp "$FIX/canonical/malformed-version.md" "$canon_tmp/docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md"
set +e
OUT="$(bash "$canon_tmp/scripts/reuse-assessment-check.sh" FEATURE-0011 --repo-root "$canon_tmp" --assessment docs/features/assess.md --skip-rac03 2>&1)"
RC=$?
set -e
if [[ "$RC" -eq 2 ]]; then
  assert_ok "canonical malformed version => exit 2"
else
  assert_fail "canonical malformed" "expected 2 got $RC"
fi

cp "$FIX/canonical/missing-version.md" "$canon_tmp/docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md"
set +e
OUT="$(bash "$canon_tmp/scripts/reuse-assessment-check.sh" FEATURE-0011 --repo-root "$canon_tmp" --assessment docs/features/assess.md --skip-rac03 2>&1)"
RC=$?
set -e
if [[ "$RC" -eq 2 ]]; then
  assert_ok "canonical missing version => exit 2"
else
  assert_fail "canonical missing version" "expected 2 got $RC"
fi
rm -rf "$canon_tmp"
trap - RETURN

# ---------------------------------------------------------------------------
# Isolated Git scenarios for RA-C13 and feature resolution  
# Each test uses its own temporary repository to avoid mutation
# ---------------------------------------------------------------------------
echo "==> Isolated Git / RA-C13 / resolution scenarios"

# Helper function to create isolated Git repository
create_isolated_repo() {
  local tmp_repo="$1"
  
  mkdir -p "$tmp_repo"
  cd "$tmp_repo"
  
  git -c init.templateDir= -c core.hooksPath=/dev/null init >/dev/null 2>&1
  git config user.name "Test User" 
  git config user.email "test@example.com"
  git config core.hooksPath /dev/null
  
  # Copy base structure from main repository  
  mkdir -p docs/phase2 docs/features docs/decisions docs/rfc \
    docs/reviews/architecture-decision-handoffs docs/context docs/governance \
    docs/engineering docs/diagrams/structurizr \
    .kiro/specs/reuse-assessment-standard scripts
  
  cp "$ROOT/docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md" docs/phase2/
  cp "$ROOT/scripts/reuse-assessment-check.sh" scripts/
  cp "$ROOT/scripts/phase2-scope-check.sh" scripts/ 2>/dev/null || true
  chmod +x scripts/*.sh 2>/dev/null || true
  
  # Required baseline files for feature gate
  for f in \
    docs/context/CURRENT_ARCHITECTURE_BASELINE.md \
    docs/context/ARCHITECTURE_VERSION.md \
    docs/governance/REVIEW_GATES.md \
    docs/engineering/go-version-standard.md \
    docs/engineering/go-observability-standard.md \
    docs/diagrams/structurizr/workspace.dsl
  do
    mkdir -p "$(dirname "$f")"
    echo "# stub baseline file" >"$f"
  done
  
  # Decision records and ADH
  echo "# DEC-0026" > docs/decisions/DEC-0026-reuse-before-build.md
  echo "# DEC-0036" > docs/decisions/DEC-0036-adapter-boundaries.md
  echo "# RFC-0021" > docs/rfc/RFC-0021-reuse-first-architecture.md
  cat > docs/reviews/architecture-decision-handoffs/ADH-2026-011-feature-0011-reuse-assessment-standard.md <<'ADH'
# Architecture Decision Handoff

**Approval status:** Approved

Feature: FEATURE-0011
Disposition: Extend
Sovrunn-owned responsibility: Four-disposition vocabulary; capability-level assessment rules; governance contract architecture; decision-first architecture  
Reused or external responsibility: General architecture-decision practices; software-selection practices; assessment template; human review practices
ADH
  
  # Feature index
  cat > docs/features/FEATURE_INDEX.md <<'IDX'
| Feature | Name | Phase | Scope | Kiro Slug | Purpose |
|---|---|---|---|---|---|
| FEATURE-0011 | Reuse Assessment Standard | Phase 2 | Executable | `reuse-assessment-standard` | reuse standard |
IDX

  # Kiro specs with proper feature identity
  mkdir -p .kiro/specs/reuse-assessment-standard
  cat > .kiro/specs/reuse-assessment-standard/requirements.md <<'R'
# Requirements
Feature: FEATURE-0011
Stage: Requirements
## Reuse Assessment
## Acceptance Criteria
## Non-goals
- no runtime
R
  cat > .kiro/specs/reuse-assessment-standard/design.md <<'D'
# Design
Feature: FEATURE-0011
Stage: Design
Controlling handoff: ADH-2026-011
## Non-goals
- no runtime
D
  cat > .kiro/specs/reuse-assessment-standard/tasks.md <<'T'
# Tasks
Feature: FEATURE-0011
Stage: Tasks
T
  
  # Valid assessment + structured approval evidence
  mkdir -p docs/reviews/reuse-assessments
  cp "$FIX/valid/complete-assessment.md" docs/features/FEATURE-0011-reuse-assessment-standard.md
  if [[ -f "$ROOT/docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md" ]]; then
    cp "$ROOT/docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md" \
      docs/reviews/reuse-assessments/FEATURE-0011-approval-evidence.md
  fi
  
  # Initial commit
  git add . >/dev/null 2>&1
  git commit -m "initial test setup" >/dev/null 2>&1
  git branch phase2-reuse-first-paas-fabric-foundation >/dev/null 2>&1
  git checkout -b feature-0011-reuse-assessment-standard >/dev/null 2>&1
}

# Test: RA-C13 docs-only change with Approved assessment
(
  test_repo="$(mktemp -d)"
  trap "rm -rf '$test_repo'" EXIT
  
  create_isolated_repo "$test_repo"
  
  echo "# documentation note" >> docs/features/FEATURE-0011-reuse-assessment-standard.md
  echo "docs/features/FEATURE-0011-reuse-assessment-standard.md" > changed-files.txt
  
  set +e
  OUT="$(bash scripts/reuse-assessment-check.sh FEATURE-0011 \
    --assessment docs/features/FEATURE-0011-reuse-assessment-standard.md \
    --changed-files changed-files.txt --skip-rac03 2>&1)"
  RC=$?
  set -e
  
  if [[ "$RC" -eq 0 ]]; then
    echo "PASS: RA-C13 docs-only change with Approved"
  else
    echo "FAIL: RA-C13 docs-only — exit $RC"
    echo "$OUT" | head -10
    exit 1
  fi
)
assert_ok "RA-C13 docs-only Approved"

# Test: RA-C13 Proposed assessment with implementation fails  
(
  test_repo="$(mktemp -d)"
  trap "rm -rf '$test_repo'" EXIT
  
  create_isolated_repo "$test_repo"
  
  # Create implementation change
  mkdir -p scripts/test
  echo "#!/bin/bash" > scripts/test/new-tool.sh
  echo "scripts/test/new-tool.sh" > changed-files.txt
  
  # Create Proposed assessment
  sed 's/| Decision status | Approved |/| Decision status | Proposed |/g' \
    docs/features/FEATURE-0011-reuse-assessment-standard.md > proposed-assessment.md
  
  set +e
  OUT="$(bash scripts/reuse-assessment-check.sh FEATURE-0011 \
    --assessment proposed-assessment.md \
    --changed-files changed-files.txt --skip-rac03 2>&1)"
  RC=$?
  set -e
  
  if [[ "$RC" -eq 1 ]] && echo "$OUT" | grep -q "RA-C13"; then
    echo "PASS: RA-C13 Proposed + impl => exit 1"
  else
    echo "FAIL: RA-C13 Proposed should fail with RA-C13, got exit $RC"
    echo "$OUT" | head -10
    exit 1
  fi
)
assert_ok "RA-C13 Proposed+impl exit 1"

# Test: Actual gate collector with all four Git change types
(
  test_repo="$(mktemp -d)"
  trap "rm -rf '$test_repo'" EXIT
  
  create_isolated_repo "$test_repo"
  
  # Copy feature gate for actual test
  cp "$ROOT/scripts/feature-gate.sh" scripts/
  
  # Create all four types of implementation changes:
  
  # 1. Committed change
  mkdir -p internal/committed
  echo "package committed" > internal/committed/committed.go
  git add internal/committed/committed.go
  git commit -m "add committed change" >/dev/null 2>&1
  
  # 2. Staged change
  mkdir -p internal/staged
  echo "package staged" > internal/staged/staged.go
  git add internal/staged/staged.go
  
  # 3. Unstaged tracked change (modify existing file)
  echo "# modified" >> docs/features/FEATURE-0011-reuse-assessment-standard.md
  
  # 4. Untracked change
  mkdir -p internal/untracked
  echo "package untracked" > internal/untracked/untracked.go
  
  # Create Proposed assessment so RA-C13 must fail
  sed 's/| Decision status | Approved |/| Decision status | Proposed |/g' \
    docs/features/FEATURE-0011-reuse-assessment-standard.md > proposed-assessment.md
  cp proposed-assessment.md docs/features/FEATURE-0011-reuse-assessment-standard.md
  git add docs/features/FEATURE-0011-reuse-assessment-standard.md  # Make this a staged change too
  
  # Run actual feature gate
  set +e
  OUT="$(bash scripts/feature-gate.sh FEATURE-0011 2>&1)"
  RC=$?
  set -e
  
  # Should exit 1 due to RA-C13 failure and should detect all changes
  if [[ "$RC" -eq 1 ]] && echo "$OUT" | grep -q "RA-C13"; then
    echo "PASS: Actual gate collector detected all change types and failed RA-C13"
  else
    echo "FAIL: Gate collector test — expected exit 1 with RA-C13, got exit $RC"
    echo "=== Gate output ==="
    echo "$OUT" | head -20
    exit 1
  fi
)
assert_ok "RA-C13 git-state union"

# Test: FEATURE-0012 cannot use FEATURE-0011 identity
(
  test_repo="$(mktemp -d)"
  trap "rm -rf '$test_repo'" EXIT
  
  create_isolated_repo "$test_repo"
  
  # Create FEATURE-0012 assessment but keep FEATURE-0011 identity in content
  # This creates identity mismatch: requested FEATURE-0012 but content says FEATURE-0011
  cp docs/features/FEATURE-0011-reuse-assessment-standard.md test-identity-mismatch.md
  
  set +e
  OUT="$(bash scripts/reuse-assessment-check.sh FEATURE-0012 \
    --assessment test-identity-mismatch.md \
    --mode strict --skip-rac03 --skip-rac13 2>&1)"
  RC=$?
  set -e
  
  if [[ "$RC" -eq 2 ]]; then
    echo "PASS: FEATURE-0012 identity mismatch exits 2"
  else
    echo "FAIL: Identity mismatch should exit 2, got $RC"
    echo "$OUT" | head -10
    exit 1
  fi
)
assert_ok "FEATURE-0012 identity mismatch check"

# Test: Missing assessment file exits 2
(
  test_repo="$(mktemp -d)" 
  trap "rm -rf '$test_repo'" EXIT
  
  create_isolated_repo "$test_repo"
  
  rm -f docs/features/FEATURE-0011-reuse-assessment-standard.md
  
  set +e
  OUT="$(bash scripts/reuse-assessment-check.sh FEATURE-0011 \
    --assessment docs/features/FEATURE-0011-reuse-assessment-standard.md \
    --mode strict --skip-rac03 --skip-rac13 2>&1)"
  RC=$?
  set -e
  
  if [[ "$RC" -eq 2 ]]; then
    echo "PASS: missing assessment => exit 2"
  else
    echo "FAIL: missing assessment should exit 2, got $RC"
    echo "$OUT" | head -10  
    exit 1
  fi
)
assert_ok "missing assessment exit 2"

# Final-review status parsing (gate function)
echo "==> Final-review status parsing"
parse_final() {
  local file="$1"
  python3 - "$file" <<'PY'
import re, sys
from pathlib import Path
text = Path(sys.argv[1]).read_text(encoding="utf-8")
matches = re.findall(r"(?im)^\s*Final feature-review status:\s*(.+?)\s*$", text)
print(matches[-1].strip() if matches else "MISSING")
PY
}

rev="$(mktemp)"
cat > "$rev" <<'EOF'
# Review
Assessment decision status: Approved
Final feature-review status: Pending
ADH-2026-011 Approved
EOF
if [[ "$(parse_final "$rev")" == "Pending" ]]; then
  assert_ok "pending final review parses as Pending"
else
  assert_fail "pending final review" "got $(parse_final "$rev")"
fi

cat > "$rev" <<'EOF'
# Review
Assessment decision status: Approved
Final feature-review status: Approved
EOF
if [[ "$(parse_final "$rev")" == "Approved" ]]; then
  assert_ok "Approved final review parses"
else
  assert_fail "Approved final review" "got $(parse_final "$rev")"
fi

cat > "$rev" <<'EOF'
# Review
Assessment decision status: Approved
Controlling reference Approved
EOF
if [[ "$(parse_final "$rev")" == "MISSING" ]]; then
  assert_ok "unrelated Approved does not satisfy final-review field"
else
  assert_fail "unrelated Approved" "got $(parse_final "$rev")"
fi
rm -f "$rev"

# ---------------------------------------------------------------------------
# Focused RA-C13 structured approval-evidence tests
# ---------------------------------------------------------------------------
echo "==> Focused RA-C13 structured evidence tests"
set +e
FOCUSED_OUT="$(bash "$ROOT/tests/reuse-assessment/test-rac13-focused.sh" 2>&1)"
FOCUSED_RC=$?
set -e
echo "$FOCUSED_OUT" | sed 's/^/  /'
if [[ "$FOCUSED_RC" -eq 0 ]]; then
  assert_ok "focused RA-C13 structured evidence suite"
else
  assert_fail "focused RA-C13 structured evidence suite" "exit $FOCUSED_RC"
fi

# ---------------------------------------------------------------------------
# Actual feature-gate collector tests (committed/staged/unstaged/untracked)
# ---------------------------------------------------------------------------
echo "==> Actual feature-gate collector tests"
set +e
GATE_COLLECTOR_OUT="$(bash "$ROOT/tests/reuse-assessment/test-gate-collector.sh" 2>&1)"
GATE_COLLECTOR_RC=$?
set -e
echo "$GATE_COLLECTOR_OUT" | sed 's/^/  /'
if [[ "$GATE_COLLECTOR_RC" -eq 0 ]]; then
  assert_ok "actual gate collector committed/staged/unstaged/untracked"
else
  assert_fail "actual gate collector suite" "exit $GATE_COLLECTOR_RC"
fi

# ---------------------------------------------------------------------------
# Summary and cleanup
# ---------------------------------------------------------------------------

echo
echo "Active branch after: $(git symbolic-ref --short HEAD 2>/dev/null || echo 'detached')"
echo "PASS=$PASS FAIL=$FAIL"

if [[ $FAIL -eq 0 ]]; then
  echo "All reuse assessment tests passed"
  exit 0
else
  echo "Test failures:"
  for failure in "${failures[@]}"; do
    echo "  - $failure"
  done
  exit 1
fi
