#!/usr/bin/env bash
# Phase 2R Canonical Model Drift Check
# Deterministic enforcement per ADH-2026-042 manifest section 7.
# Scans all active authority files from manifest sections 5.1–5.5.
# ALL checks are hard failures. No WARN-then-pass.

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
FAIL=0

fail() { echo "  FAIL: $1"; FAIL=$((FAIL + 1)); }
pass() { echo "  PASS"; }

echo "=== Phase 2R Canonical Drift Check ==="
echo ""

# ---------------------------------------------------------------------------
# Active authority file sets (manifest sections 5.1–5.5)
# ---------------------------------------------------------------------------

# 5.1 Context files
CONTEXT_FILES=(
  "$REPO_ROOT/docs/context/ARCHITECTURE_VERSION.md"
  "$REPO_ROOT/docs/context/CURRENT_ARCHITECTURE_BASELINE.md"
  "$REPO_ROOT/docs/context/CURRENT_PHASE_CONTEXT.md"
  "$REPO_ROOT/docs/context/CURRENT_DECISION_SUMMARY.md"
  "$REPO_ROOT/docs/context/OPEN_QUESTIONS.md"
  "$REPO_ROOT/docs/context/SOVRUNN_CONTEXT_PACK.md"
  "$REPO_ROOT/docs/context/CHATGPT_ARCHITECTURE_SESSION_PROMPT.md"
)

# 5.2 Architecture and vocabulary
ARCH_FILES=(
  "$REPO_ROOT/docs/glossary.md"
  "$REPO_ROOT/docs/architecture/development-phases.md"
  "$REPO_ROOT/docs/architecture/platform-core.md"
  "$REPO_ROOT/docs/architecture/placement-decision-engine.md"
  "$REPO_ROOT/docs/architecture/plugin-taxonomy-and-boundaries.md"
  "$REPO_ROOT/docs/architecture/adapter-boundary-model.md"
  "$REPO_ROOT/docs/architecture/organization-governance.md"
  "$REPO_ROOT/docs/architecture/controller-reconciliation-model.md"
  "$REPO_ROOT/docs/architecture/gitops-desired-state-model.md"
  "$REPO_ROOT/docs/architecture/vertical-slices/VS-000-contract-specification.md"
  "$REPO_ROOT/docs/api/BOUNDARY_LEDGER.md"
  "$REPO_ROOT/docs/api/boundary-ledger.yaml"
)

# 5.3 Phase, roadmap, traceability, MVP
PHASE_FILES=(
  "$REPO_ROOT/docs/phase2/PHASE2_SCOPE.md"
  "$REPO_ROOT/docs/phase2/PHASE2_ARCHITECTURE_SPINE.md"
  "$REPO_ROOT/docs/phase2/PHASE2_FEATURE_SEQUENCE.md"
  "$REPO_ROOT/docs/phase2/PHASE2_ACCEPTANCE_GATES.md"
  "$REPO_ROOT/docs/phase2/PHASE2_EXECUTION_STRATEGY.md"
  "$REPO_ROOT/docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md"
  "$REPO_ROOT/docs/features/FEATURE_INDEX.md"
  "$REPO_ROOT/docs/traceability/FEATURE_TRACEABILITY_MATRIX.md"
  "$REPO_ROOT/docs/traceability/DECISION_TRACEABILITY_MATRIX.md"
  "$REPO_ROOT/docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md"
  "$REPO_ROOT/docs/mvp/MVP_001_GOVERNED_POSTGRESQL_PAAS.md"
  "$REPO_ROOT/docs/mvp/MVP_001_CUSTOMER_DEMO_SCENARIOS.md"
  "$REPO_ROOT/docs/mvp/MVP_001_OUT_OF_SCOPE.md"
)

# 5.4 Diagrams
DIAGRAM_FILES=(
  "$REPO_ROOT/docs/diagrams/structurizr/workspace.dsl"
)

# 5.5 Kiro and feature-factory controls
KIRO_FILES=(
  "$REPO_ROOT/.kiro/agents/sovrunn-spec.md"
  "$REPO_ROOT/.kiro/steering/architecture.md"
  "$REPO_ROOT/.kiro/steering/product.md"
  "$REPO_ROOT/.kiro/steering/engineering.md"
  "$REPO_ROOT/.kiro/steering/slice0-contract.md"
)

# Combined set for stale-term scanning
ALL_ACTIVE=("${CONTEXT_FILES[@]}" "${ARCH_FILES[@]}" "${PHASE_FILES[@]}" "${DIAGRAM_FILES[@]}" "${KIRO_FILES[@]}")

# Files that must contain exact feature titles
TITLE_FILES=(
  "$REPO_ROOT/docs/phase2/PHASE2_FEATURE_SEQUENCE.md"
  "$REPO_ROOT/docs/features/FEATURE_INDEX.md"
  "$REPO_ROOT/docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md"
  "$REPO_ROOT/docs/architecture/development-phases.md"
  "$REPO_ROOT/docs/phase2/PHASE2_ARCHITECTURE_SPINE.md"
)

# ---------------------------------------------------------------------------
# Allowed-context detection (heading-aware and file-specific)
# ---------------------------------------------------------------------------

# Returns 0 (true = allowed) if the line is in a historical/negation/out-of-scope context.
# Arguments: $1=file path, $2=line number, $3=line content
is_allowed_context() {
  local file="$1" lineno="$2" line="$3"
  local lower
  lower=$(echo "$line" | tr '[:upper:]' '[:lower:]')

  # --- File-specific allowances ---

  # provider-neutral-resource-model.md is explicitly historical
  case "$file" in
    *provider-neutral-resource-model.md) return 0 ;;
  esac

  # FEATURE_INDEX.md: completed/merged features are historical records
  if [[ "$file" == *FEATURE_INDEX.md ]]; then
    if echo "$line" | grep -qE "Implemented|Merged|Completed"; then
      return 0
    fi
  fi

  # Phase 1 section in development-phases.md (lines describing what Phase 1 built)
  if [[ "$file" == *development-phases.md ]]; then
    # Phase 1 purpose block mentions ServiceClass historically
    if echo "$lower" | grep -qE "^build core resource grammar"; then
      return 0
    fi
    # Completed features table entries
    if echo "$line" | grep -qE "Completed|Merged|alpha model"; then
      return 0
    fi
  fi

  # platform-core.md is marked historical-phase1
  if [[ "$file" == *platform-core.md ]]; then
    # The whole file is historical; check front matter marker
    if head -10 "$file" | grep -qi "historical"; then
      return 0
    fi
  fi

  # organization-governance.md: Phase 1 audit table and pre-Phase2R sections are historical input
  if [[ "$file" == *organization-governance.md ]]; then
    if echo "$lower" | grep -qE "input history|for phase 2 and later.*input history|canonical contract is"; then
      return 0
    fi
  fi

  # --- Inline negation/supersession patterns ---
  # These patterns indicate the line is explicitly denying, superseding, or restricting the term.
  if echo "$lower" | grep -qE "supersed|no active|not a canonical|migration input|is not a canonical|replaced by|retired|historical|do not use|formerly|not use|must not|does not|should not|no longer|is superseded|out of scope|non.goal|replaces "; then
    return 0
  fi

  # "No mandatory X" or "no X" patterns — but NOT bare "mandatory" alone
  if echo "$lower" | grep -qE "^-.*no (resourcepool|providercapability|mandatory)"; then
    return 0
  fi
  if echo "$lower" | grep -qE "no.*(resourcepool|providercapability|provider-wide)"; then
    return 0
  fi
  if echo "$lower" | grep -qE "(resourcepool|providercapability).*(supersed|not|removed|no longer|is not|as active)"; then
    return 0
  fi

  # DEC-0042 / DEC-0049 / DEC-0050 reference on same line = migration context
  if echo "$lower" | grep -qE "dec-004[0-9]|dec-005[0-9]"; then
    return 0
  fi

  # Explicit non-goal/out-of-scope/prohibited list item
  if echo "$lower" | grep -qE "^- .*(not include|not build|must not|non-goal|out of scope|prohibited)"; then
    return 0
  fi

  # --- Heading-aware: determine if current section is allowed ---
  # Read backwards from the line to find the nearest heading
  local heading
  heading=$(head -n "$lineno" "$file" | grep -E "^#{1,6} " | tail -1 | tr '[:upper:]' '[:lower:]' || true)

  if echo "$heading" | grep -qE "non.goal|out of scope|supersed|retired|migration|historical|phase 1|completed|not approved|prohibited|drift|avoid"; then
    return 0
  fi

  # Check for "Avoid:" or "Do not:" label preceding the line in the same block (within 10 lines)
  local preceding
  preceding=$(head -n "$lineno" "$file" | tail -15 | tr '[:upper:]' '[:lower:]')
  if echo "$preceding" | grep -qE "^avoid:|^do not"; then
    return 0
  fi

  # Check for "input history" context in surrounding paragraph (within 5 lines before)
  local nearby
  nearby=$(sed -n "$((lineno > 5 ? lineno - 5 : 1)),${lineno}p" "$file" | tr '[:upper:]' '[:lower:]')
  if echo "$nearby" | grep -qE "input history|this table is input history"; then
    return 0
  fi

  # If none of the above matched, the context is NOT allowed
  return 1
}

# ---------------------------------------------------------------------------
# CHECK 1: Exact twelve FEATURE-0015–0026 titles across five authority files
# ---------------------------------------------------------------------------
echo "1. Checking all twelve FEATURE-0015–0026 titles across five files..."
EXPECTED_TITLES=(
  "Canonical Cloud Model Foundation"
  "Adapter Boundary and ExecutionTarget Qualification"
  "Policy Evaluation Abstraction"
  "Governance, IAM, Approval and Exception Foundation"
  "Sovereignty Facts, Evidence and Policy Foundation"
  "Assignment and Effective Governance Resolution"
  "CloudEnrollment, Personal Onboarding, Entitlement and Quota"
  "Service Product, Runtime and Requirement Foundation"
  "Sovereignty and Placement Decision v0"
  "Plugin Taxonomy and Synthetic Execution Boundary"
  "AI-Readable Decision and Operation Context"
  "Slice 0 Integration and Conformance Demo"
)
TITLE_FAIL=0
for title in "${EXPECTED_TITLES[@]}"; do
  for f in "${TITLE_FILES[@]}"; do
    if [ -f "$f" ] && ! grep -qF "$title" "$f" 2>/dev/null; then
      echo "  FAIL: \"$title\" not found in $(basename "$f")"
      TITLE_FAIL=1
      FAIL=$((FAIL + 1))
    fi
  done
done
[ "$TITLE_FAIL" -eq 0 ] && pass

# ---------------------------------------------------------------------------
# CHECK 2: Exactly seven canonical scope kinds in glossary
# ---------------------------------------------------------------------------
echo "2. Checking exactly seven canonical scope kinds..."
if grep -qF "Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider" "$REPO_ROOT/docs/glossary.md"; then
  pass
else
  fail "Seven scope kinds string not found in glossary"
fi

# ---------------------------------------------------------------------------
# CHECK 3: CloudEnrollment and CloudProviderParticipation distinct
# ---------------------------------------------------------------------------
echo "3. Checking CloudEnrollment and CloudProviderParticipation are distinct..."
if grep -q "CloudEnrollment" "$REPO_ROOT/docs/glossary.md" && \
   grep -q "CloudProviderParticipation" "$REPO_ROOT/docs/glossary.md"; then
  pass
else
  fail "CloudEnrollment or CloudProviderParticipation missing from glossary"
fi

# ---------------------------------------------------------------------------
# CHECK 4: No active stale ResourcePool / ProviderCapability
# ---------------------------------------------------------------------------
echo "4. Checking no active authority uses ResourcePool/ProviderCapability as target..."
RP_FAIL=0
for f in "${ALL_ACTIVE[@]}"; do
  [ -f "$f" ] || continue
  while IFS= read -r match; do
    lineno="${match%%:*}"
    line="${match#*:}"
    if ! is_allowed_context "$f" "$lineno" "$line"; then
      echo "  FAIL [$(basename "$f"):$lineno]: $line"
      RP_FAIL=1
      FAIL=$((FAIL + 1))
    fi
  done < <(grep -n -E "ResourcePool|ProviderCapability" "$f" 2>/dev/null || true)
done
[ "$RP_FAIL" -eq 0 ] && pass

# ---------------------------------------------------------------------------
# CHECK 5: No active generic Provider as combined owner/operator
# ---------------------------------------------------------------------------
echo "5. Checking no active generic Provider as combined concept..."
# Only check glossary for "| Provider |" as an active table entry
PROV_FAIL=0
if grep -q "^| Provider |" "$REPO_ROOT/docs/glossary.md"; then
  # Must be in the retired/migration section
  section_heading=$(grep -B 20 "^| Provider |" "$REPO_ROOT/docs/glossary.md" | grep -E "^## " | tail -1 || true)
  if echo "$section_heading" | grep -qiE "retired|migration"; then
    : # ok
  else
    fail "Provider appears as active term in glossary (not in retired section)"
    PROV_FAIL=1
  fi
fi
[ "$PROV_FAIL" -eq 0 ] && pass

# ---------------------------------------------------------------------------
# CHECK 6: No active ServiceClass authority
# ---------------------------------------------------------------------------
echo "6. Checking ServiceClass has no active canonical authority..."
if grep -qE "ServiceClass.*(No active|no active|migration input)" "$REPO_ROOT/docs/glossary.md"; then
  pass
else
  fail "ServiceClass may still have active canonical authority in glossary"
fi

# ---------------------------------------------------------------------------
# CHECK 7: No active EffectivePolicyContext
# ---------------------------------------------------------------------------
echo "7. Checking no active EffectivePolicyContext in authorities..."
EPC_FAIL=0
for f in "${ALL_ACTIVE[@]}"; do
  [ -f "$f" ] || continue
  while IFS= read -r match; do
    lineno="${match%%:*}"
    line="${match#*:}"
    if ! is_allowed_context "$f" "$lineno" "$line"; then
      echo "  FAIL [$(basename "$f"):$lineno]: $line"
      EPC_FAIL=1
      FAIL=$((FAIL + 1))
    fi
  done < <(grep -n "EffectivePolicyContext" "$f" 2>/dev/null || true)
done
[ "$EPC_FAIL" -eq 0 ] && pass

# ---------------------------------------------------------------------------
# CHECK 8: No active provider-neutral as canonical adjective
# ---------------------------------------------------------------------------
echo "8. Checking no active provider-neutral as canonical adjective..."
PN_FAIL=0
for f in "${ALL_ACTIVE[@]}"; do
  [ -f "$f" ] || continue
  while IFS= read -r match; do
    lineno="${match%%:*}"
    line="${match#*:}"
    if ! is_allowed_context "$f" "$lineno" "$line"; then
      echo "  FAIL [$(basename "$f"):$lineno]: $line"
      PN_FAIL=1
      FAIL=$((FAIL + 1))
    fi
  done < <(grep -n "provider-neutral" "$f" 2>/dev/null || true)
done
[ "$PN_FAIL" -eq 0 ] && pass

# ---------------------------------------------------------------------------
# CHECK 9: No fewer-than-seven active scope vocabulary
# ---------------------------------------------------------------------------
echo "9. Checking no active six-scope vocabulary..."
SIX_FAIL=0
for f in "${ALL_ACTIVE[@]}"; do
  [ -f "$f" ] || continue
  while IFS= read -r match; do
    lineno="${match%%:*}"
    line="${match#*:}"
    if ! is_allowed_context "$f" "$lineno" "$line"; then
      echo "  FAIL [$(basename "$f"):$lineno]: $line"
      SIX_FAIL=1
      FAIL=$((FAIL + 1))
    fi
  done < <(grep -n -iE "six.*(scope|ScopeKind|governance voc)" "$f" 2>/dev/null || true)
done
[ "$SIX_FAIL" -eq 0 ] && pass

# ---------------------------------------------------------------------------
# CHECK 10: Customer/core schemas have no provider-native objects gate
# ---------------------------------------------------------------------------
echo "10. Checking provider-native exclusion gate exists..."
if grep -qF "no provider-native" "$REPO_ROOT/docs/phase2/PHASE2_ACCEPTANCE_GATES.md" && \
   grep -qF "provider-native" "$REPO_ROOT/docs/phase2/PHASE2_ARCHITECTURE_SPINE.md"; then
  pass
else
  fail "Provider-native exclusion gate not found in acceptance gates or spine"
fi

# ---------------------------------------------------------------------------
# CHECK 11: No Go code changed
# ---------------------------------------------------------------------------
echo "11. Checking no Go code changed..."
GO_CHANGES=$(git -C "$REPO_ROOT" diff --name-only HEAD 2>/dev/null | grep '\.go$' || true)
GO_UNTRACKED=$(git -C "$REPO_ROOT" ls-files --others --exclude-standard 2>/dev/null | grep '\.go$' || true)
if [ -z "$GO_CHANGES" ] && [ -z "$GO_UNTRACKED" ]; then
  pass
else
  fail "Go files changed or added: $GO_CHANGES $GO_UNTRACKED"
fi

# ---------------------------------------------------------------------------
# CHECK 12: No generated prompts or site output changed
# ---------------------------------------------------------------------------
echo "12. Checking no generated prompts or site output changed..."
GEN_CHANGES=$(git -C "$REPO_ROOT" diff --name-only HEAD 2>/dev/null | grep -E '^site/|^docs/generated' || true)
GEN_UNTRACKED=$(git -C "$REPO_ROOT" ls-files --others --exclude-standard 2>/dev/null | grep -E '^site/|^docs/generated' || true)
if [ -z "$GEN_CHANGES" ] && [ -z "$GEN_UNTRACKED" ]; then
  pass
else
  fail "Generated prompts or site output changed: $GEN_CHANGES $GEN_UNTRACKED"
fi

# ---------------------------------------------------------------------------
# CHECK 13: Architecture baseline consistent
# ---------------------------------------------------------------------------
echo "13. Checking architecture baseline consistency..."
BASELINE="ARCH-2026.08-PHASE2R-CANONICAL"
BASELINE_CHECK_FILES=(
  "$REPO_ROOT/docs/context/ARCHITECTURE_VERSION.md"
  "$REPO_ROOT/docs/context/CURRENT_ARCHITECTURE_BASELINE.md"
  "$REPO_ROOT/docs/context/CURRENT_PHASE_CONTEXT.md"
  "$REPO_ROOT/.kiro/steering/architecture.md"
  "$REPO_ROOT/docs/phase2/PHASE2_ARCHITECTURE_SPINE.md"
)
BL_FAIL=0
for f in "${BASELINE_CHECK_FILES[@]}"; do
  if [ -f "$f" ] && ! grep -qF "$BASELINE" "$f"; then
    echo "  FAIL: Baseline not found in $(basename "$f")"
    BL_FAIL=1
    FAIL=$((FAIL + 1))
  fi
done
[ "$BL_FAIL" -eq 0 ] && pass

# ---------------------------------------------------------------------------
# CHECK 14: DecisionRecord remains shared envelope
# ---------------------------------------------------------------------------
echo "14. Checking DecisionRecord is the shared envelope..."
if grep -qF "DecisionRecord" "$REPO_ROOT/docs/glossary.md" && \
   grep -qF "DecisionProfile" "$REPO_ROOT/docs/glossary.md"; then
  pass
else
  fail "DecisionRecord/DecisionProfile not found in glossary"
fi

# ---------------------------------------------------------------------------
# CHECK 15: Architecture spine does not contain old pre-canonical body
# ---------------------------------------------------------------------------
echo "15. Checking architecture spine is canonical (no old body)..."
SPINE="$REPO_ROOT/docs/phase2/PHASE2_ARCHITECTURE_SPINE.md"
SPINE_FAIL=0
# The old spine had these unique markers that must not exist in the rewritten version
for marker in \
  "ResourcePool catalogue" \
  "ProviderCapabilities" \
  "FEATURE-0015 ResourcePool" \
  "FEATURE-0020 ProfileAssignment and EffectivePolicyContext" \
  "FEATURE-0026 Phase 2 Integration Demo" \
  "P2-C09 — ResourcePool Is the Placement Boundary" \
  "P2-C10 — Capabilities Drive Compatibility" \
  "six-value" \
  "six governance scopes"; do
  if grep -qF "$marker" "$SPINE" 2>/dev/null; then
    echo "  FAIL: Old spine marker found: \"$marker\""
    SPINE_FAIL=1
    FAIL=$((FAIL + 1))
  fi
done
[ "$SPINE_FAIL" -eq 0 ] && pass

# ---------------------------------------------------------------------------
# CHECK 16: Active decision-status mirrors agree with the Decision Index
# ---------------------------------------------------------------------------
echo "16. Checking decision status and supersession consistency..."
DECISION_INDEX="$REPO_ROOT/docs/decisions/DECISION_INDEX.md"
DECISION_TRACE="$REPO_ROOT/docs/traceability/DECISION_TRACEABILITY_MATRIX.md"
REBASELINE="$REPO_ROOT/docs/phase2/PHASE2R_REBASELINE.md"
DEC_STATUS_FAIL=0

if ! awk -F'|' '
  function trim(s) {
    gsub(/^[[:space:]]+|[[:space:]]+$/, "", s)
    return s
  }
  function normalized(s) {
    s = trim(s)
    if (s ~ /^Accepted([[:space:]]|$|\()/) return "Accepted"
    if (s ~ /^Proposed([[:space:]]|$|\()/) return "Proposed"
    return s
  }
  FNR == NR && /^\| DEC-[0-9]+/ {
    id = trim($2)
    index_status[id] = normalized($5)
    next
  }
  FNR != NR && /^\| DEC-[0-9]+/ {
    cell = trim($2)
    split(cell, parts, /[[:space:]]+/)
    id = parts[1]
    trace_status = normalized($3)
    if (!(id in index_status)) {
      printf "  FAIL: %s exists in traceability but not Decision Index\n", id
      failed = 1
    } else if (trace_status != index_status[id]) {
      printf "  FAIL: %s status mismatch: Decision Index=%s; traceability=%s\n", id, index_status[id], trace_status
      failed = 1
    }
  }
  END { exit failed }
' "$DECISION_INDEX" "$DECISION_TRACE"; then
  DEC_STATUS_FAIL=1
  FAIL=$((FAIL + 1))
fi

while IFS='|' read -r _ dec _ _ status _; do
  dec=$(echo "$dec" | xargs)
  status=$(echo "$status" | xargs)
  case "$status" in
    "Superseded by "*)
      replacement=${status#Superseded by }
      if grep -qE "^\|[[:space:]]*$dec[[:space:]]*\|" "$REBASELINE" && \
         ! grep -E "^\|[[:space:]]*$dec[[:space:]]*\|" "$REBASELINE" | grep -qF "Superseded by $replacement"; then
        fail "$dec supersession to $replacement is not mirrored by PHASE2R_REBASELINE"
        DEC_STATUS_FAIL=1
      fi
      ;;
  esac
done < <(grep -E '^\| DEC-[0-9]+' "$DECISION_INDEX")

[ "$DEC_STATUS_FAIL" -eq 0 ] && pass

# ---------------------------------------------------------------------------
# RESULT
# ---------------------------------------------------------------------------
echo ""
if [ "$FAIL" -eq 0 ]; then
  echo "=== ALL DRIFT CHECKS PASSED ==="
  exit 0
else
  echo "=== DRIFT CHECKS FAILED ($FAIL issues) ==="
  exit 1
fi
