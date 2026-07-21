#!/usr/bin/env bash
# Sovrunn Feature Gate
# Authoritative Git change-set discovery and FEATURE-0011+ reuse-assessment orchestration.
set -euo pipefail

FEATURE="${1:-${FEATURE:-}}"
PHASE_BRANCH="${PHASE_BRANCH:-phase2-reuse-first-paas-fabric-foundation}"

if [[ -z "$FEATURE" ]]; then
  echo "ERROR: feature id is required"
  echo "Usage: ./scripts/feature-gate.sh FEATURE-0011"
  exit 2
fi

if [[ ! "$FEATURE" =~ ^FEATURE-[0-9]{4}$ ]]; then
  echo "ERROR: invalid feature id: $FEATURE"
  echo "Expected format: FEATURE-0011"
  exit 2
fi

FEATURE_NUM_RAW="${FEATURE#FEATURE-}"
FEATURE_NUM=$((10#$FEATURE_NUM_RAW))

LEGACY_PHASE1=false
if (( FEATURE_NUM <= 10 )); then
  LEGACY_PHASE1=true
fi

echo "==> Sovrunn Feature Gate: $FEATURE"

if [[ "$LEGACY_PHASE1" == "true" ]]; then
  echo "INFO: $FEATURE is Phase 1 legacy baseline; using legacy validation mode"
else
  echo "INFO: $FEATURE is Phase 2+; using strict AOS validation mode"
fi

fail() {
  echo "FAIL: $1"
  exit 1
}

fail_config() {
  echo "FAIL(config): $1"
  exit 2
}

pass() {
  echo "PASS: $1"
}

require_file() {
  local file="$1"
  [[ -f "$file" ]] || fail "Missing required file: $file"
  pass "Found $file"
}

require_contains() {
  local file="$1"
  local pattern="$2"
  local label="$3"
  grep -qi "$pattern" "$file" || fail "$label missing in $file"
  pass "$label present in $file"
}

# Resolve exactly one FEATURE_INDEX.md row and its Kiro Slug for FEATURE-0011+.
# No grep-based Kiro-directory fallback. No FEATURE-0011 path fallback for later features.
resolve_feature_paths() {
  local feature="$1"
  local index="docs/features/FEATURE_INDEX.md"
  [[ -f "$index" ]] || fail_config "Missing feature index: $index"

  local result
  result="$(python3 - "$feature" "$index" <<'PY'
import re
import sys
from pathlib import Path

feature = sys.argv[1]
index = Path(sys.argv[2])
text = index.read_text(encoding="utf-8")
rows = []
for line in text.splitlines():
    if not line.strip().startswith("|"):
        continue
    cells = [c.strip().replace("`", "") for c in line.strip().strip("|").split("|")]
    if len(cells) < 5:
        continue
    if cells[0] == "Feature" or set(cells[0]) <= {"-"}:
        continue
    if cells[0] == feature:
        rows.append(cells)

if len(rows) == 0:
    print("MISSING")
    sys.exit(0)
if len(rows) > 1:
    print("DUPLICATE")
    sys.exit(0)

cells = rows[0]
# Expected: Feature | Name | Phase | Scope | Kiro Slug | Purpose
slug = cells[4] if len(cells) > 4 else ""
if not slug or slug.lower() in {"kiro slug", "-"}:
    print("MALFORMED")
    sys.exit(0)
if not re.fullmatch(r"[a-z0-9]+(?:-[a-z0-9]+)*", slug):
    print("MALFORMED")
    sys.exit(0)
print(f"OK|{slug}")
PY
)"

  case "$result" in
    MISSING)
      fail_config "FEATURE_INDEX.md has no row for $feature"
      ;;
    DUPLICATE)
      fail_config "FEATURE_INDEX.md has duplicate rows for $feature"
      ;;
    MALFORMED)
      fail_config "FEATURE_INDEX.md row for $feature has a missing or malformed Kiro Slug"
      ;;
    OK\|*)
      KIRO_SLUG="${result#OK|}"
      ;;
    *)
      fail_config "unexpected feature-index resolution result for $feature"
      ;;
  esac

  KIRO_SPEC_DIR=".kiro/specs/${KIRO_SLUG}"
  ASSESSMENT_PATH="docs/features/${feature}-${KIRO_SLUG}.md"
  REQUIREMENTS_PATH="${KIRO_SPEC_DIR}/requirements.md"
  DESIGN_PATH="${KIRO_SPEC_DIR}/design.md"
  TASKS_PATH="${KIRO_SPEC_DIR}/tasks.md"
}

# Authoritative Git change-set discovery (committed + staged + unstaged + untracked).
# Fails closed: requires phase branch and merge-base to exist.
collect_changed_files() {
  local tmp
  tmp="$(mktemp)"
  if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    fail_config "not inside a git work tree; cannot collect changed files"
  fi

  local phase_ref=""
  local merge_base=""
  
  # Require either local phase branch or its remote-tracking equivalent
  if git rev-parse --verify "$PHASE_BRANCH" >/dev/null 2>&1; then
    phase_ref="$PHASE_BRANCH"
  elif git rev-parse --verify "origin/${PHASE_BRANCH}" >/dev/null 2>&1; then
    phase_ref="origin/${PHASE_BRANCH}"
  else
    fail_config "phase branch '${PHASE_BRANCH}' not found (local or origin/${PHASE_BRANCH})"
  fi

  # Require merge-base to succeed and return non-empty result
  if ! merge_base="$(git merge-base HEAD "$phase_ref" 2>/dev/null)"; then
    fail_config "git merge-base failed between HEAD and $phase_ref"
  fi
  if [[ -z "$merge_base" ]]; then
    fail_config "git merge-base returned empty result for HEAD and $phase_ref"
  fi

  # Collect changes with explicit status handling, no failure swallowing
  local committed_out staged_out unstaged_out untracked_out
  
  # Execute each Git command separately with explicit error handling
  if ! committed_out="$(git diff --name-only "${merge_base}..HEAD" 2>&1)"; then
    rm -f "$tmp"
    fail_config "git diff committed changes failed: $committed_out"
  fi
  
  if ! staged_out="$(git diff --name-only --cached 2>&1)"; then
    rm -f "$tmp"
    fail_config "git diff staged changes failed: $staged_out"
  fi
  
  if ! unstaged_out="$(git diff --name-only 2>&1)"; then
    rm -f "$tmp"
    fail_config "git diff unstaged changes failed: $unstaged_out"
  fi
  
  if ! untracked_out="$(git ls-files --others --exclude-standard 2>&1)"; then
    rm -f "$tmp"
    fail_config "git ls-files untracked failed: $untracked_out"
  fi
  
  # Combine outputs with deterministic normalization, de-duplication, and sorting
  {
    echo "$committed_out"
    echo "$staged_out"
    echo "$unstaged_out"
    echo "$untracked_out"
  } | sed '/^$/d' | sort -u >"$tmp"

  CHANGED_FILES_LIST="$tmp"
}

require_file docs/context/CURRENT_ARCHITECTURE_BASELINE.md
require_file docs/context/ARCHITECTURE_VERSION.md
require_file docs/governance/REVIEW_GATES.md
require_file docs/engineering/go-version-standard.md
require_file docs/engineering/go-observability-standard.md
require_file docs/diagrams/structurizr/workspace.dsl
require_file docs/features/FEATURE_INDEX.md

KIRO_SPEC_DIR=""
KIRO_SLUG=""
ASSESSMENT_PATH=""
REQUIREMENTS_PATH=""
DESIGN_PATH=""
TASKS_PATH=""
CHANGED_FILES_LIST=""

if [[ "$LEGACY_PHASE1" == "true" ]]; then
  # Legacy path resolution may use index slug when present, with historical grep fallback.
  FEATURE_DOC_FILES=()
  while IFS= read -r file; do
    FEATURE_DOC_FILES+=("$file")
  done < <(find docs/features -maxdepth 1 -type f -name "${FEATURE}*.md" 2>/dev/null | sort || true)

  if [[ -f docs/features/FEATURE_INDEX.md ]]; then
    set +e
    KIRO_SLUG="$(
      python3 - "$FEATURE" <<'PY'
import sys
from pathlib import Path
feature = sys.argv[1]
text = Path("docs/features/FEATURE_INDEX.md").read_text(encoding="utf-8")
for line in text.splitlines():
    if not line.strip().startswith("|"):
        continue
    cells = [c.strip().replace("`", "") for c in line.strip().strip("|").split("|")]
    if cells and cells[0] == feature and len(cells) > 4:
        print(cells[4])
        break
PY
    )"
    slug_rc=$?
    set -e
    if [[ $slug_rc -ne 0 ]]; then
      KIRO_SLUG=""
    fi
  fi
  if [[ -n "$KIRO_SLUG" && -d ".kiro/specs/$KIRO_SLUG" ]]; then
    KIRO_SPEC_DIR=".kiro/specs/$KIRO_SLUG"
  elif [[ -d ".kiro/specs" ]]; then
    KIRO_SPEC_DIR="$(grep -Ril "$FEATURE" .kiro/specs/*/requirements.md 2>/dev/null | head -n 1 | xargs dirname 2>/dev/null || true)"
  fi

  if [[ -z "$KIRO_SPEC_DIR" && ${#FEATURE_DOC_FILES[@]} -eq 0 ]]; then
    fail "No feature documentation found for $FEATURE. Expected .kiro/specs/<slug>/ or docs/features/${FEATURE}*.md."
  fi

  if [[ -n "$KIRO_SPEC_DIR" ]]; then
    echo "==> Checking Kiro spec: $KIRO_SPEC_DIR"
    require_file "$KIRO_SPEC_DIR/requirements.md"
    require_file "$KIRO_SPEC_DIR/design.md"
    require_file "$KIRO_SPEC_DIR/tasks.md"
    echo "INFO: skipping strict AOS Kiro section checks for Phase 1 legacy spec"
  fi

  if [[ ${#FEATURE_DOC_FILES[@]} -gt 0 ]]; then
    echo "==> Checking feature docs"
    printf '%s\n' "${FEATURE_DOC_FILES[@]}"
    echo "INFO: skipping strict feature-doc Acceptance Criteria check for Phase 1 legacy docs"
  fi
else
  echo "==> Resolving feature paths from docs/features/FEATURE_INDEX.md"
  resolve_feature_paths "$FEATURE"
  echo "INFO: Kiro slug=$KIRO_SLUG"
  echo "INFO: assessment=$ASSESSMENT_PATH"
  echo "INFO: kiro=$KIRO_SPEC_DIR"

  # Configuration errors for missing resolved paths (exit 2)
  [[ -f "$REQUIREMENTS_PATH" ]] || fail_config "Missing resolved requirements: $REQUIREMENTS_PATH"
  [[ -f "$DESIGN_PATH" ]] || fail_config "Missing resolved design: $DESIGN_PATH"
  [[ -f "$TASKS_PATH" ]] || fail_config "Missing resolved tasks: $TASKS_PATH"
  [[ -f "$ASSESSMENT_PATH" ]] || fail_config "Missing resolved assessment: $ASSESSMENT_PATH"

  # Active-feature identity checks
  require_contains "$REQUIREMENTS_PATH" "$FEATURE" "Active feature identity"
  require_contains "$DESIGN_PATH" "$FEATURE" "Active feature identity"
  require_contains "$TASKS_PATH" "$FEATURE" "Active feature identity"
  require_contains "$ASSESSMENT_PATH" "$FEATURE" "Active feature identity"

  # Stage labels
  require_contains "$REQUIREMENTS_PATH" "Stage: Requirements" "Requirements stage"
  require_contains "$DESIGN_PATH" "Stage: Design" "Design stage"
  require_contains "$TASKS_PATH" "Stage: Tasks" "Tasks stage"

  # Reuse summary + non-goals + controlling ADH
  require_contains "$REQUIREMENTS_PATH" "Reuse Assessment" "Reuse Assessment"
  require_contains "$REQUIREMENTS_PATH" "Acceptance Criteria" "Acceptance Criteria"
  if ! grep -qiE "Non-goals|Out of scope|non-goals" "$REQUIREMENTS_PATH"; then
    fail "Non-goals missing in $REQUIREMENTS_PATH"
  fi
  pass "Non-goals present in $REQUIREMENTS_PATH"
  require_contains "$ASSESSMENT_PATH" "Feature-level reuse summary" "Feature-level reuse summary"
  require_contains "$ASSESSMENT_PATH" "ADH-" "Controlling ADH reference"

  # Applicability-aware design checks:
  # FEATURE-0011 is governance-only and must not be forced to invent runtime
  # Architecture Drift / Observability / Security design sections.
  # Later runtime features keep those checks.
  if [[ "$FEATURE" == "FEATURE-0011" ]]; then
    echo "INFO: FEATURE-0011 governance-only — skipping runtime Architecture Drift/Observability/Security heading checks"
    require_contains "$DESIGN_PATH" "Non-goals" "Non-goals"
    require_contains "$DESIGN_PATH" "ADH-2026-011" "Controlling ADH reference"
  else
    require_contains "$DESIGN_PATH" "Architecture Drift" "Architecture Drift Checks"
    require_contains "$DESIGN_PATH" "Observability" "Observability"
    require_contains "$DESIGN_PATH" "Security" "Security"
    require_contains "$DESIGN_PATH" "Non-goals" "Non-goals"
  fi

  echo "==> Collecting authoritative changed-file list"
  collect_changed_files
  pass "Changed-file list collected ($(wc -l <"$CHANGED_FILES_LIST" | tr -d ' ') files)"

  if [[ -x scripts/reuse-assessment-check.sh ]]; then
    echo "==> Running reuse-assessment validator (strict)"
    set +e
    bash scripts/reuse-assessment-check.sh "$FEATURE" \
      --assessment "$ASSESSMENT_PATH" \
      --mode strict \
      --changed-files "$CHANGED_FILES_LIST" \
      --requirements "$REQUIREMENTS_PATH" \
      --design "$DESIGN_PATH" \
      --tasks "$TASKS_PATH"
    rc=$?
    set -e
    rm -f "$CHANGED_FILES_LIST"
    
    if [[ $rc -eq 2 ]]; then
      fail_config "reuse-assessment validator configuration error"
    elif [[ $rc -eq 1 ]]; then
      fail "reuse-assessment validator failed"
    fi
    pass "reuse-assessment validator"
  else
    rm -f "$CHANGED_FILES_LIST"
    fail_config "scripts/reuse-assessment-check.sh is missing or not executable"
  fi
  rm -f "$CHANGED_FILES_LIST"
fi

echo "==> Checking generated artifacts are not staged"
if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  if git diff --cached --name-only | grep -E '(^docs/generated-prompts/|^site/|^\.automation/generated-prompts/|^\.automation/logs/|^\.automation/reviews/|^docs/context/SOVRUNN_CONTEXT_PACK\.generated\.md$)' >/dev/null; then
    fail "Generated prompt/site/log/review artifacts are staged"
  fi
  pass "Generated artifacts are not staged"
else
  echo "WARN: not inside git work tree; skipping staged artifact check"
fi

if [[ -x scripts/structurizr-check.sh ]]; then
  echo "==> Running Structurizr workspace check"
  scripts/structurizr-check.sh
  pass "Structurizr workspace check"
fi

if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  if git diff --name-only HEAD | grep -E '^(docs/architecture/|docs/rfc/|docs/decisions/|docs/context/CURRENT_ARCHITECTURE_BASELINE.md)' >/dev/null; then
    if ! git diff --name-only HEAD | grep -E '^docs/diagrams/structurizr/workspace.dsl$' >/dev/null; then
      echo "WARN: architecture docs changed but Structurizr workspace.dsl was not updated"
      echo "WARN: if the change affects system/container/plugin/external relationships, update workspace.dsl"
    fi
  fi
fi

if compgen -G "*.go" >/dev/null || find cmd internal pkg api -type f -name '*.go' 2>/dev/null | grep -q .; then
  echo "==> Running gofmt check"
  UNFORMATTED="$(gofmt -l . | grep '\.go$' || true)"
  if [[ -n "$UNFORMATTED" ]]; then
    echo "$UNFORMATTED"
    fail "Go files need formatting"
  fi
  pass "gofmt check"

  echo "==> Running go test ./..."
  go test ./...
  pass "go test ./..."

  echo "==> Running go test -race ./..."
  go test -race ./...
  pass "go test -race ./..."

  if command -v golangci-lint >/dev/null 2>&1; then
    echo "==> Running golangci-lint"
    golangci-lint run ./...
    pass "golangci-lint"
  else
    echo "WARN: golangci-lint not installed; skipping"
  fi

  if command -v gosec >/dev/null 2>&1; then
    echo "==> Running gosec"
    gosec ./...
    pass "gosec"
  else
    echo "WARN: gosec not installed; skipping"
  fi
fi

if [[ -x scripts/phase1-consistency-check.sh ]]; then
  if [[ -f go.mod && -d internal ]]; then
    echo "==> Running phase1 consistency check"
    scripts/phase1-consistency-check.sh
    pass "phase1 consistency check"
  else
    echo "WARN: go.mod/internal not found; skipping phase1 consistency check in docs-only archive"
  fi
fi

if [[ -x scripts/phase2-scope-check.sh ]]; then
  if [[ "$LEGACY_PHASE1" == "true" ]]; then
    echo "INFO: skipping Phase 2 scope check for Phase 1 legacy feature"
  else
    echo "==> Running phase2 scope check"
    scripts/phase2-scope-check.sh "$FEATURE"
    pass "phase2 scope check"
  fi
fi

# Exact final-review marker. Only "Final feature-review status: Approved" satisfies merge approval.
REVIEW_FILE="docs/reviews/feature-gates/${FEATURE}-approval-review.md"
if [[ -f "$REVIEW_FILE" ]]; then
  echo "==> Checking final feature-review status in $REVIEW_FILE"
  FINAL_STATUS="$(python3 - "$REVIEW_FILE" <<'PY'
import re
import sys
from pathlib import Path
text = Path(sys.argv[1]).read_text(encoding="utf-8")
matches = re.findall(r"(?im)^\s*Final feature-review status:\s*(.+?)\s*$", text)
if not matches:
    print("MISSING")
else:
    print(matches[-1].strip())
PY
)"
  if [[ "$FINAL_STATUS" == "MISSING" ]]; then
    fail "Final feature-review status field missing or malformed in $REVIEW_FILE"
  fi
  if [[ "$FINAL_STATUS" != "Approved" ]]; then
    fail "Final feature-review status is '${FINAL_STATUS}' (required: Approved). Assessment decision status or other Approved mentions do not satisfy final merge approval."
  fi
  pass "Final feature-review status: Approved"
else
  if [[ "$LEGACY_PHASE1" == "true" ]]; then
    echo "WARN: $REVIEW_FILE not found"
    echo "WARN: strict team mode should require final approval review before merge"
  else
    echo "WARN: $REVIEW_FILE not found"
    echo "WARN: automated checks may pass, but final merge approval requires Final feature-review status: Approved"
  fi
fi

echo
echo "SUCCESS: $FEATURE passed Sovrunn feature gate"
