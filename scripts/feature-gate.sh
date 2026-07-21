#!/usr/bin/env bash
set -euo pipefail

FEATURE="${1:-${FEATURE:-}}"

if [[ -z "$FEATURE" ]]; then
  echo "ERROR: feature id is required"
  echo "Usage: ./scripts/feature-gate.sh FEATURE-0011"
  exit 1
fi

if [[ ! "$FEATURE" =~ ^FEATURE-[0-9]{4}$ ]]; then
  echo "ERROR: invalid feature id: $FEATURE"
  echo "Expected format: FEATURE-0011"
  exit 1
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

resolve_kiro_slug() {
  local feature="$1"
  local index="docs/features/FEATURE_INDEX.md"
  [[ -f "$index" ]] || return 0
  awk -F'|' -v feature="$feature" '
    $2 ~ feature {
      slug=$6
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", slug)
      gsub(/`/, "", slug)
      print slug
      exit
    }
  ' "$index"
}

require_file docs/context/CURRENT_ARCHITECTURE_BASELINE.md
require_file docs/context/ARCHITECTURE_VERSION.md
require_file docs/governance/REVIEW_GATES.md
require_file docs/engineering/go-version-standard.md
require_file docs/engineering/go-observability-standard.md
require_file docs/diagrams/structurizr/workspace.dsl
require_file docs/features/FEATURE_INDEX.md

FEATURE_DOC_FILES=()
while IFS= read -r file; do
  FEATURE_DOC_FILES+=("$file")
done < <(find docs/features -maxdepth 1 -type f -name "${FEATURE}*.md" 2>/dev/null | sort || true)

KIRO_SPEC_DIR=""
KIRO_SLUG="$(resolve_kiro_slug "$FEATURE" || true)"
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

  if [[ "$LEGACY_PHASE1" == "true" ]]; then
    echo "INFO: skipping strict AOS Kiro section checks for Phase 1 legacy spec"
  else
    require_contains "$KIRO_SPEC_DIR/requirements.md" "Reuse Assessment" "Reuse Assessment"
    require_contains "$KIRO_SPEC_DIR/requirements.md" "Acceptance Criteria" "Acceptance Criteria"
    require_contains "$KIRO_SPEC_DIR/design.md" "Architecture Drift" "Architecture Drift Checks"
    require_contains "$KIRO_SPEC_DIR/design.md" "Observability" "Observability"
    require_contains "$KIRO_SPEC_DIR/design.md" "Security" "Security"
    require_contains "$KIRO_SPEC_DIR/design.md" "Non-goals" "Non-goals"
  fi
fi

if [[ ${#FEATURE_DOC_FILES[@]} -gt 0 ]]; then
  echo "==> Checking feature docs"
  printf '%s\n' "${FEATURE_DOC_FILES[@]}"

  if [[ "$LEGACY_PHASE1" == "true" ]]; then
    echo "INFO: skipping strict feature-doc Acceptance Criteria check for Phase 1 legacy docs"
  else
    grep -qi "Acceptance Criteria" "${FEATURE_DOC_FILES[@]}" || fail "Acceptance Criteria missing in feature docs for $FEATURE"
    pass "Acceptance Criteria present"
  fi
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

REVIEW_FILE="docs/reviews/feature-gates/${FEATURE}-approval-review.md"
if [[ -f "$REVIEW_FILE" ]]; then
  require_contains "$REVIEW_FILE" "APPROVED" "Approval marker"
else
  echo "WARN: $REVIEW_FILE not found"
  echo "WARN: strict team mode should require final approval review before merge"
fi

echo
echo "SUCCESS: $FEATURE passed Sovrunn feature gate"
