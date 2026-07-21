#!/usr/bin/env bash
# Prove each Git change source through actual feature gate
# Four separate isolated repositories, one change type each.
# Each scenario uses its own subshell + trap (never overwrite one global EXIT trap).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
FIX="$ROOT/tests/reuse-assessment/fixtures"

create_gate_test_repo() {
  local tmp_repo="$1"

  mkdir -p "$tmp_repo"
  cd "$tmp_repo"

  git -c init.templateDir= -c core.hooksPath=/dev/null init >/dev/null 2>&1
  git config user.name "Test User"
  git config user.email "test@example.com"

  mkdir -p docs/phase2 docs/features docs/decisions docs/rfc \
    docs/reviews/architecture-decision-handoffs docs/context docs/governance \
    docs/engineering docs/diagrams/structurizr \
    .kiro/specs/reuse-assessment-standard scripts

  cp "$ROOT/docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md" docs/phase2/
  cp "$ROOT/scripts/reuse-assessment-check.sh" scripts/
  cp "$ROOT/scripts/phase2-scope-check.sh" scripts/ 2>/dev/null || true
  cp "$ROOT/scripts/feature-gate.sh" scripts/
  chmod +x scripts/*.sh 2>/dev/null || true

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

  echo "# DEC-0026" > docs/decisions/DEC-0026-reuse-before-build.md
  echo "# DEC-0036" > docs/decisions/DEC-0036-adapter-boundaries.md
  echo "# RFC-0021" > docs/rfc/RFC-0021-reuse-first-architecture.md
  cat > docs/reviews/architecture-decision-handoffs/ADH-2026-011-feature-0011-reuse-assessment-standard.md <<'ADH'
# Architecture Decision Handoff

- **Approval status:** Approved
- **Related feature:** FEATURE-0011

Feature: FEATURE-0011
Disposition: Extend

Approved by: Test Owner
Date: 2026-07-21
ADH

  cat > docs/features/FEATURE_INDEX.md <<'IDX'
| Feature | Name | Phase | Scope | Kiro Slug | Purpose |
|---|---|---|---|---|---|
| FEATURE-0011 | Reuse Assessment Standard | Phase 2 | Executable | `reuse-assessment-standard` | reuse standard |
IDX

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

  cp "$FIX/valid/complete-assessment.md" docs/features/FEATURE-0011-reuse-assessment-standard.md
  python3 - <<'PY'
from pathlib import Path
p = Path("docs/features/FEATURE-0011-reuse-assessment-standard.md")
text = p.read_text(encoding="utf-8")
text = text.replace("| Decision status | Approved |", "| Decision status | Proposed |")
p.write_text(text, encoding="utf-8")
PY

  git add . >/dev/null 2>&1
  git commit -m "initial test setup" >/dev/null 2>&1
  git branch phase2-reuse-first-paas-fabric-foundation >/dev/null 2>&1
  git checkout -b feature-0011-reuse-assessment-standard >/dev/null 2>&1
}

run_one_gate_change_source() {
  local test_name="$1"
  local change_type="$2"

  (
    local test_repo
    test_repo="$(mktemp -d)"
    trap 'rm -rf "$test_repo"' EXIT

    create_gate_test_repo "$test_repo"

    case "$change_type" in
      committed)
        mkdir -p internal/committed
        echo "package committed" > internal/committed/committed.go
        git add internal/committed/committed.go
        git commit -m "add committed implementation" >/dev/null 2>&1
        ;;
      staged)
        mkdir -p internal/staged
        echo "package staged" > internal/staged/staged.go
        git add internal/staged/staged.go
        ;;
      unstaged)
        mkdir -p internal/unstaged
        echo "package unstaged" > internal/unstaged/unstaged.go
        git add internal/unstaged/unstaged.go
        git commit -m "add file for unstaged test" >/dev/null 2>&1
        echo "// modified" >> internal/unstaged/unstaged.go
        ;;
      untracked)
        mkdir -p internal/untracked
        echo "package untracked" > internal/untracked/untracked.go
        ;;
      *)
        echo "FAIL: unknown change type $change_type"
        exit 1
        ;;
    esac

    set +e
    OUT="$(bash scripts/feature-gate.sh FEATURE-0011 2>&1)"
    RC=$?
    set -e

    if [[ "$RC" -eq 1 ]] && echo "$OUT" | grep -q "RA-C13"; then
      if ! echo "$OUT" | grep -q "CONFIG"; then
        echo "$test_name: PASS (gate exit 1, RA-C13 detected, no CONFIG errors)"
      else
        echo "$test_name: FAIL (detected CONFIG error)"
        echo "$OUT" | grep CONFIG || true
        exit 1
      fi
    else
      echo "$test_name: FAIL (expected gate exit 1 with RA-C13, got exit $RC)"
      echo "=== Gate output sample ==="
      echo "$OUT" | tail -10
      exit 1
    fi
  )
}

echo "==> Testing actual gate collector with each Git change source"

run_one_gate_change_source "Committed implementation change" "committed"
run_one_gate_change_source "Staged implementation change" "staged"
run_one_gate_change_source "Unstaged implementation change" "unstaged"
run_one_gate_change_source "Untracked implementation change" "untracked"

echo
echo "All gate collector tests passed - each change source properly detected"
echo "Temporary-repository cleanup: each scenario used an isolated subshell trap"
exit 0
