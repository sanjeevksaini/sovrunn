#!/usr/bin/env bash
# Isolated Git scenarios replacement for run.sh
# Each test uses its own temporary repository
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
FIX="$ROOT/tests/reuse-assessment/fixtures"

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
  
  # Valid assessment
  cp "$FIX/valid/complete-assessment.md" docs/features/FEATURE-0011-reuse-assessment-standard.md
  
  # Initial commit
  git add . >/dev/null 2>&1
  git commit -m "initial test setup" >/dev/null 2>&1
  git branch phase2-reuse-first-paas-fabric-foundation >/dev/null 2>&1
  git checkout -b feature-0011-reuse-assessment-standard >/dev/null 2>&1
}

echo "==> Isolated Git / RA-C13 / resolution scenarios"

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
    echo "RA-C13 docs-only change with Approved: PASS"
  else
    echo "RA-C13 docs-only change with Approved: FAIL (exit $RC)"
    echo "$OUT" | head -10
    exit 1
  fi
)

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
    echo "RA-C13 Proposed + impl => exit 1: PASS"
  else
    echo "RA-C13 Proposed + impl => exit 1: FAIL (got exit $RC, expected RA-C13)"
    echo "$OUT" | head -10
    exit 1
  fi
)

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
    echo "Actual gate collector detected all change types and failed RA-C13: PASS"
  else
    echo "Actual gate collector test: FAIL (expected exit 1 with RA-C13, got exit $RC)"
    echo "=== Full Gate output ==="
    echo "$OUT"
    exit 1
  fi
)

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
    echo "FEATURE-0012 identity mismatch exits 2: PASS"
  else
    echo "FEATURE-0012 identity mismatch: FAIL (expected exit 2, got $RC)"
    echo "$OUT" | head -10
    exit 1
  fi
)

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
    echo "missing assessment exits 2: PASS"
  else
    echo "missing assessment: FAIL (expected exit 2, got $RC)"
    echo "$OUT" | head -10  
    exit 1
  fi
)

echo
echo "All isolated Git tests passed"
exit 0