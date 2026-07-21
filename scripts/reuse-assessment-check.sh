#!/usr/bin/env bash
# reuse_assessment_format_version=1.0.0
# Sovrunn Reuse Assessment validator (FEATURE-0011)
# Exit codes: 0 pass, 1 validation failure, 2 usage/config/internal error
set -euo pipefail

usage() {
  cat <<'EOF' >&2
Usage: scripts/reuse-assessment-check.sh FEATURE-NNNN [options]

Options:
  --repo-root PATH           Repository root (default: cwd)
  --assessment PATH          Assessment markdown path (required in strict mode)
  --mode strict|legacy       Validation mode (default: strict for FEATURE-0011+)
  --changed-files PATH       File containing newline-separated changed paths (RA-C13)
  --requirements PATH        Kiro requirements.md for the active feature
  --design PATH              Kiro design.md for the active feature
  --tasks PATH               Kiro tasks.md for the active feature
  --skip-rac03               Skip RA-C03 scope check invocation
  --skip-rac13               Skip RA-C13 approval enforcement (assessment-only validation)
  -h, --help                 Show help

Exit codes:
  0  validation passed
  1  validation failure
  2  usage, configuration, or internal error
EOF
}

FEATURE=""
REPO_ROOT=""
ASSESSMENT=""
MODE=""
CHANGED_FILES=""
REQUIREMENTS=""
DESIGN=""
TASKS=""
SKIP_RAC03=0
SKIP_RAC13=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)
      usage
      exit 2
      ;;
    --repo-root)
      REPO_ROOT="${2:-}"
      shift 2
      ;;
    --assessment)
      ASSESSMENT="${2:-}"
      shift 2
      ;;
    --mode)
      MODE="${2:-}"
      shift 2
      ;;
    --changed-files)
      CHANGED_FILES="${2:-}"
      shift 2
      ;;
    --requirements)
      REQUIREMENTS="${2:-}"
      shift 2
      ;;
    --design)
      DESIGN="${2:-}"
      shift 2
      ;;
    --tasks)
      TASKS="${2:-}"
      shift 2
      ;;
    --skip-rac03)
      SKIP_RAC03=1
      shift
      ;;
    --skip-rac13)
      SKIP_RAC13=1
      shift
      ;;
    FEATURE-[0-9][0-9][0-9][0-9])
      if [[ -n "$FEATURE" ]]; then
        echo "ERROR: multiple feature identifiers supplied" >&2
        exit 2
      fi
      FEATURE="$1"
      shift
      ;;
    *)
      if [[ -z "$FEATURE" && "$1" =~ ^FEATURE-[0-9]{4}$ ]]; then
        FEATURE="$1"
        shift
      elif [[ -z "$FEATURE" && "$1" =~ ^FEATURE- ]]; then
        echo "ERROR: malformed feature identifier: $1" >&2
        exit 2
      else
        echo "ERROR: unknown argument: $1" >&2
        usage
        exit 2
      fi
      ;;
  esac
done

if [[ -z "$FEATURE" ]]; then
  echo "ERROR: feature identifier is required" >&2
  usage
  exit 2
fi

if [[ ! "$FEATURE" =~ ^FEATURE-[0-9]{4}$ ]]; then
  echo "ERROR: invalid feature identifier: $FEATURE" >&2
  exit 2
fi

FEATURE_NUM_RAW="${FEATURE#FEATURE-}"
FEATURE_NUM=$((10#$FEATURE_NUM_RAW))

if [[ -z "$REPO_ROOT" ]]; then
  REPO_ROOT="$(pwd)"
fi
REPO_ROOT="$(cd "$REPO_ROOT" && pwd)"

if [[ -z "$MODE" ]]; then
  if (( FEATURE_NUM <= 10 )); then
    MODE="legacy"
  else
    MODE="strict"
  fi
fi

if [[ "$MODE" != "strict" && "$MODE" != "legacy" ]]; then
  echo "ERROR: --mode must be strict or legacy" >&2
  exit 2
fi

if [[ "$MODE" == "legacy" ]]; then
  echo "INFO: $FEATURE legacy mode — strict reuse-assessment checks skipped"
  exit 0
fi

if [[ -z "$ASSESSMENT" ]]; then
  echo "ERROR: --assessment is required in strict mode" >&2
  exit 2
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export REPO_ROOT FEATURE ASSESSMENT MODE CHANGED_FILES REQUIREMENTS DESIGN TASKS SKIP_RAC03 SKIP_RAC13 SCRIPT_DIR

python3 - "$REPO_ROOT" "$FEATURE" "$ASSESSMENT" "$MODE" "$CHANGED_FILES" "$REQUIREMENTS" "$DESIGN" "$TASKS" "$SKIP_RAC03" "$SKIP_RAC13" "$SCRIPT_DIR" <<'PY'
import os
import re
import subprocess
import sys
from pathlib import Path

REPO_ROOT = Path(sys.argv[1]).resolve()
FEATURE = sys.argv[2]
ASSESSMENT = sys.argv[3]
MODE = sys.argv[4]
CHANGED_FILES = sys.argv[5]
REQUIREMENTS = sys.argv[6]
DESIGN = sys.argv[7]
TASKS = sys.argv[8]
SKIP_RAC03 = sys.argv[9] == "1"
SKIP_RAC13 = sys.argv[10] == "1"
SCRIPT_DIR = Path(sys.argv[11]).resolve()

CANONICAL_REL = "docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md"
CANONICAL_PATH = REPO_ROOT / CANONICAL_REL
VERSION_RE = re.compile(r"^\d+\.\d+\.\d+$")
FEATURE_RE = re.compile(r"^FEATURE-[0-9]{4}$")
CONCEPTUAL_LABEL = "Conceptual example — illustrative only and outside execution scope"
SHELL_VERSION_MARKER = "# reuse_assessment_format_version="

DISPOSITIONS = {"Reuse", "Wrap", "Extend", "Build"}
DECISION_STATUSES = {"Proposed", "Approved", "Deferred", "Rejected", "Superseded"}
YES_NO = {"Yes", "No"}
REPLACEMENT_RISKS = {"Low", "Medium", "High"}

REQUIRED_SECTIONS = [
    "Feature-level reuse summary",
    "Identity",
    "Classification",
    "Analysis",
    "Boundary",
    "Suitability",
    "Phase and scope",
    "Risk mitigation",
    "Traceability",
]

# Field labels matched case-insensitively against table rows or bold/list forms
REQUIRED_FIELDS = [
    "Feature identity",
    "Capability or decision-unit identity",
    "Assessment owner",
    "Disposition",
    "Decision status",
    "Assessment scope",
    "Candidate category",
    "Mature candidates / applicable standards",
    "Relevant candidate strengths",
    "Material candidate constraints",
    "Rationale",
    "Selected foundation or approach",
    "Sovrunn-owned responsibility",
    "Reused or extended responsibility",  # alias: "Reused or external responsibility"
    "Data crossing the boundary",
    "Control crossing the boundary",
    "Adapter required",
    "Adapter rationale",
    "Adapter or contract identifier",
    "Vendor-native types allowed",
    "Sovereignty and deployment fit",
    "Security and trust",
    "Operational and supportability",
    "Licensing and supply-chain",
    "Portability and provider-neutrality impact",
    "Allowed in current phase",
    "Current-phase work",
    "Deferred work",
    "Explicit non-goals",
    "Exit or migration boundary",
    "Phase 2 non-goal acknowledgement",
    "Risk mitigation",  # Can be "Applicable architecture risks" or "Risk-control matrix"
    "Residual risk",
    "Replacement risk",
    "Reassessment triggers",
    "Related DEC / RFC / ADH references",
    "Linked acceptance criteria",
    "Validation and review evidence",
]

RISK_CONTROL_FIELDS = [
    "Preventive controls",
    "Detection controls",
    "Corrective path",
]

TRACEABILITY_FIELDS = [
    "Related DEC / RFC / ADH references",
    "Linked acceptance criteria",
    "Validation and review evidence",
]

OPERATIONAL_LINK_TARGETS = [
    "docs/prompts/kiro/requirements.prompt.md",
    "docs/prompts/kiro/design.prompt.md",
    "docs/prompts/kiro/tasks.prompt.md",
    "docs/prompts/cursor/task.prompt.md",
    "docs/prompts/reviewer/spec-review.prompt.md",
    "docs/prompts/reviewer/approval-review.prompt.md",
    "docs/automation/KIRO_DECISION_POLICY.md",
    "docs/automation/FEATURE_FACTORY.md",
    "docs/ai/AI_FEATURE_FACTORY.md",
    "docs/templates/ARCHITECTURE_DECISION_HANDOFF.md",
    "docs/templates/ARCHITECTURE_CHANGE_REQUEST.md",
    "docs/templates/RFC_TEMPLATE.md",
    "docs/templates/FEATURE_REVIEW_TEMPLATE.md",
]

IMPLEMENTATION_PREFIXES = ("cmd/", "internal/", "pkg/", "api/", "scripts/")

diagnostics = []
config_error = False


def emit(rule_id, layer, path, section, message, severity="error", guidance=""):
    diagnostics.append(
        {
            "rule_id": rule_id,
            "layer": layer,
            "feature": FEATURE,
            "path": str(path),
            "section": section or "",
            "message": message,
            "severity": severity,
            "guidance": guidance,
        }
    )


def read_text(path: Path):
    try:
        return path.read_text(encoding="utf-8")
    except OSError as exc:
        return None, str(exc)


def resolve_assessment_path():
    p = Path(ASSESSMENT)
    if not p.is_absolute():
        p = REPO_ROOT / p
    return p.resolve()


def extract_front_matter_version(text: str):
    if not text.startswith("---"):
        return None
    end = text.find("\n---", 3)
    if end < 0:
        return None
    block = text[3:end]
    for line in block.splitlines():
        m = re.match(r"^\s*reuse_assessment_format_version:\s*['\"]?([^'\"\s]+)", line)
        if m:
            return m.group(1).strip()
    return None


def extract_shell_version(text: str):
    for line in text.splitlines():
        if line.startswith(SHELL_VERSION_MARKER):
            return line.split("=", 1)[1].strip()
    return None


def field_value_map(text: str):
    """Extract Field|Value table rows and '- Field: Value' forms."""
    values = {}
    # Markdown tables: | Field | Value |
    for m in re.finditer(
        r"^\|\s*([^|]+?)\s*\|\s*([^|]*?)\s*\|",
        text,
        flags=re.MULTILINE,
    ):
        key = m.group(1).strip()
        val = m.group(2).strip()
        if key.lower() in {"field", "---", ""}:
            continue
        if set(key) <= {"-"}:
            continue
        values[key] = val
    # Bold / list forms
    for m in re.finditer(
        r"^\s*[-*]\s*\*?\*?([^:*\n]+?)\*?\*?\s*:\s*(.+)$",
        text,
        flags=re.MULTILINE,
    ):
        key = m.group(1).strip()
        val = m.group(2).strip()
        values.setdefault(key, val)
    return values


def find_field(values, *names):
    lower_map = {k.lower(): (k, v) for k, v in values.items()}
    for name in names:
        hit = lower_map.get(name.lower())
        if hit:
            return hit[1]
    # partial contains
    for name in names:
        for k, v in values.items():
            if name.lower() in k.lower():
                return v
    return None


def has_section(text: str, title: str) -> bool:
    """Check if section exists as actual Markdown heading (no prose fallback)."""
    pat = re.compile(rf"^\s{{0,3}}#{{1,6}}\s+.*{re.escape(title)}\s*$", re.IGNORECASE | re.MULTILINE)
    return bool(pat.search(text))


def has_field_label(text: str, label: str) -> bool:
    # Markdown table field: | Field | Value |
    if re.search(rf"\|\s*{re.escape(label)}\s*\|", text, flags=re.IGNORECASE):
        return True
    # Markdown heading: ## Field or ### Field or #### Field
    if re.search(rf"^\s*#{{1,6}}\s+.*{re.escape(label)}\s*$", text, flags=re.IGNORECASE | re.MULTILINE):
        return True
    # List item with colon: - Field: or * Field:
    if re.search(rf"^\s*[-*]\s+\*?\*?{re.escape(label)}\*?\*?\s*:", text, flags=re.IGNORECASE | re.MULTILINE):
        return True
    return False


def has_heading_with_content(text: str, label: str) -> bool:
    """Check if heading exists with non-empty content before next heading of equal/higher level."""
    # Find the heading (exact match, no trailing whitespace in pattern)
    heading_pattern = r"^(\s*)(#{1,6})\s+" + re.escape(label) + r"\s*$"
    heading_match = re.search(heading_pattern, text, re.IGNORECASE | re.MULTILINE)
    if not heading_match:
        return False
    
    heading_level = len(heading_match.group(2))  # Number of # characters
    start_pos = heading_match.end()
    
    # Look for content until next heading of same or higher level
    remaining_text = text[start_pos:]
    
    # Split into lines and check for content
    lines = remaining_text.split('\n')
    content_found = False
    
    for line in lines:
        # Check if this is a heading of same or higher level (fewer #'s)
        heading_check = re.match(r'^(\s*)(#{1,6})\s+', line)
        if heading_check:
            next_level = len(heading_check.group(2))
            if next_level <= heading_level:
                break
        
        # Check for non-empty content (excluding just whitespace)
        if line.strip():
            content_found = True
            break
    
    return content_found


def has_populated_field(text: str, label: str) -> bool:
    """Check if field has both valid syntax and actual non-empty content."""
    # First try table/list field parsing
    values = field_value_map(text)
    field_val = find_field(values, label)
    if field_val is not None and str(field_val).strip():
        return True
    
    # Then try heading-form with content
    return has_heading_with_content(text, label)


# --- Canonical version resolution (config errors => exit 2) ---
canonical_version = None
if not CANONICAL_PATH.is_file():
    emit(
        "CONFIG",
        "config",
        CANONICAL_REL,
        "reuse_assessment_format_version",
        "canonical reuse assessment standard file is missing",
        guidance="Restore docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md",
    )
    config_error = True
else:
    try:
        raw = CANONICAL_PATH.read_text(encoding="utf-8")
    except (OSError, UnicodeError) as exc:
        emit(
            "CONFIG",
            "config",
            CANONICAL_REL,
            "canonical file",
            f"canonical reuse assessment standard file unreadable: {exc}",
            guidance="Restore readable docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md",
        )
        config_error = True
        canonical_version = None
    else:
        canonical_version = extract_front_matter_version(raw)
    if not canonical_version:
        emit(
            "CONFIG",
            "config",
            CANONICAL_REL,
            "reuse_assessment_format_version",
            "canonical reuse_assessment_format_version is missing",
            guidance="Add reuse_assessment_format_version to canonical front matter",
        )
        config_error = True
    elif not VERSION_RE.match(canonical_version):
        emit(
            "CONFIG",
            "config",
            CANONICAL_REL,
            "reuse_assessment_format_version",
            f"canonical reuse_assessment_format_version is malformed: {canonical_version}",
            guidance="Use semantic version MAJOR.MINOR.PATCH",
        )
        config_error = True

if config_error:
    for d in sorted(diagnostics, key=lambda x: (x["path"], x["section"], x["rule_id"])):
        print(
            f"{d['severity']}|{d['rule_id']}|{d['layer']}|{d['feature']}|{d['path']}|{d['section']}|{d['message']}|{d['guidance']}"
        )
    sys.exit(2)

assessment_path = resolve_assessment_path()
if not assessment_path.is_file():
    emit(
        "CONFIG",
        "config",
        str(assessment_path),
        "assessment",
        "assessment artifact is missing",
        guidance="Supply a resolvable --assessment path for the active feature",
    )
    for d in sorted(diagnostics, key=lambda x: (x["path"], x["section"], x["rule_id"])):
        print(
            f"{d['severity']}|{d['rule_id']}|{d['layer']}|{d['feature']}|{d['path']}|{d['section']}|{d['message']}|{d['guidance']}"
        )
    sys.exit(2)

try:
    assessment_text = assessment_path.read_text(encoding="utf-8")
except (OSError, UnicodeError) as exc:
    emit(
        "CONFIG",
        "config",
        str(assessment_path),
        "assessment",
        f"assessment artifact unreadable: {exc}",
    )
    for d in sorted(diagnostics, key=lambda x: (x["path"], x["section"], x["rule_id"])):
        print(
            f"{d['severity']}|{d['rule_id']}|{d['layer']}|{d['feature']}|{d['path']}|{d['section']}|{d['message']}|{d['guidance']}"
        )
    sys.exit(2)

rel_assessment = str(assessment_path.relative_to(REPO_ROOT)) if str(assessment_path).startswith(str(REPO_ROOT)) else str(assessment_path)
values = field_value_map(assessment_text)

# RA-C10 feature identifier + active-feature identity match
feature_identity = find_field(values, "Feature identity") or ""
if not FEATURE_RE.match(FEATURE):
    emit("RA-C10", "consistency", rel_assessment, "Feature identity", f"invalid feature identifier: {FEATURE}")
elif feature_identity and not FEATURE_RE.match(feature_identity):
    emit("RA-C10", "consistency", rel_assessment, "Feature identity", f"invalid feature identifier: {feature_identity}")
elif feature_identity and feature_identity != FEATURE:
    emit(
        "CONFIG",
        "config",
        rel_assessment,
        "Feature identity",
        f"assessment feature identity '{feature_identity}' does not match active feature '{FEATURE}'",
        guidance="Supply the active feature's own assessment path",
    )
    for d in sorted(diagnostics, key=lambda x: (x["path"], x["section"], x["rule_id"])):
        print(
            f"{d['severity']}|{d['rule_id']}|{d['layer']}|{d['feature']}|{d['path']}|{d['section']}|{d['message']}|{d['guidance']}"
        )
    sys.exit(2)

# RA-C09 version on assessment
assess_version = extract_front_matter_version(assessment_text)
if not assess_version:
    emit(
        "RA-C09",
        "consistency",
        rel_assessment,
        "reuse_assessment_format_version",
        "required assessment version is missing",
        guidance=f"Declare reuse_assessment_format_version: {canonical_version}",
    )
elif not VERSION_RE.match(assess_version):
    emit(
        "RA-C09",
        "consistency",
        rel_assessment,
        "reuse_assessment_format_version",
        f"assessment version is malformed: {assess_version}",
    )
elif assess_version != canonical_version:
    emit(
        "RA-C09",
        "consistency",
        rel_assessment,
        "reuse_assessment_format_version",
        f"assessment version {assess_version} does not match canonical {canonical_version}",
    )

# RA-C09 validator shell marker
validator_script = SCRIPT_DIR / "reuse-assessment-check.sh"
if validator_script.is_file():
    vtext = validator_script.read_text(encoding="utf-8", errors="replace")
    vver = extract_shell_version(vtext)
    if not vver:
        emit(
            "RA-C09",
            "consistency",
            "scripts/reuse-assessment-check.sh",
            "reuse_assessment_format_version",
            "validator version marker is missing",
        )
    elif not VERSION_RE.match(vver):
        emit(
            "RA-C09",
            "consistency",
            "scripts/reuse-assessment-check.sh",
            "reuse_assessment_format_version",
            f"validator version marker is malformed: {vver}",
        )
    elif vver != canonical_version:
        emit(
            "RA-C09",
            "consistency",
            "scripts/reuse-assessment-check.sh",
            "reuse_assessment_format_version",
            f"validator version {vver} does not match canonical {canonical_version}",
        )

# Layer 1 structural: sections
for section in REQUIRED_SECTIONS:
    if not has_section(assessment_text, section):
        emit("RA-S01", "structural", rel_assessment, section, f"required section missing: {section}")

# Layer 1 structural: fields
# Risk-control and traceability fields use RA-S09 / RA-S10 instead of RA-S02.
_special = {f.lower() for f in RISK_CONTROL_FIELDS + TRACEABILITY_FIELDS}
for field in REQUIRED_FIELDS:
    if field.lower() in _special:
        continue
    # Special handling for risk mitigation (can be either format)
    if field == "Risk mitigation":
        has_risks = (has_field_label(assessment_text, "Applicable architecture risks") or 
                     has_field_label(assessment_text, "Risk-control matrix"))
        if not has_risks:
            emit("RA-S02", "structural", rel_assessment, field, 
                f"required risk mitigation section missing (use 'Applicable architecture risks' or 'Risk-control matrix')")
        continue
    if field == "Reused or extended responsibility":
        if not (
            has_populated_field(assessment_text, "Reused or extended responsibility")
            or has_populated_field(assessment_text, "Reused or external responsibility")
        ):
            emit(
                "RA-S02",
                "structural",
                rel_assessment,
                field,
                "required field missing or empty: Reused or extended/external responsibility",
            )
        continue
    if not has_populated_field(assessment_text, field):
        emit("RA-S02", "structural", rel_assessment, field, f"required field missing or empty: {field}")

disposition = find_field(values, "Disposition")
if disposition is not None:
    if disposition not in DISPOSITIONS:
        emit("RA-S03", "structural", rel_assessment, "Disposition", f"invalid disposition: {disposition}")

decision_status = find_field(values, "Decision status")
if decision_status is not None:
    if decision_status not in DECISION_STATUSES:
        emit("RA-S04", "structural", rel_assessment, "Decision status", f"invalid decision status: {decision_status}")

adapter_required = find_field(values, "Adapter required")
if adapter_required is not None:
    if adapter_required not in YES_NO:
        emit("RA-S05", "structural", rel_assessment, "Adapter required", f"invalid adapter value: {adapter_required}")

allowed_phase = find_field(values, "Allowed in current phase")
if allowed_phase is not None:
    if allowed_phase not in YES_NO:
        emit("RA-S06", "structural", rel_assessment, "Allowed in current phase", f"invalid phase value: {allowed_phase}")

vendor_native = find_field(values, "Vendor-native types allowed")
if vendor_native is not None:
    ok = vendor_native == "No" or vendor_native.startswith("Approved exception")
    if not ok:
        emit(
            "RA-S07",
            "structural",
            rel_assessment,
            "Vendor-native types allowed",
            f"invalid vendor-native-types value: {vendor_native}",
        )

# Replacement risk (table or heading form)
replacement_risk = find_field(values, "Replacement risk")
if replacement_risk is None or str(replacement_risk).strip() not in REPLACEMENT_RISKS:
    m = re.search(
        r"Replacement risk\s*\n+\s*(Low|Medium|High|\S+)\b",
        assessment_text,
        flags=re.IGNORECASE,
    )
    if m:
        token = m.group(1)
        capped = token[:1].upper() + token[1:].lower()
        if capped in REPLACEMENT_RISKS:
            replacement_risk = capped
        else:
            replacement_risk = token
    elif has_field_label(assessment_text, "Replacement risk"):
        replacement_risk = ""
if has_field_label(assessment_text, "Replacement risk"):
    if replacement_risk not in REPLACEMENT_RISKS:
        emit(
            "RA-S08",
            "structural",
            rel_assessment,
            "Replacement risk",
            f"invalid replacement-risk value: {replacement_risk or '<missing>'}",
        )

# RA-S09: Risk control fields (only required if using legacy format, not matrix format)
uses_matrix = has_field_label(assessment_text, "Risk-control matrix")
if not uses_matrix:
    for field in RISK_CONTROL_FIELDS:
        if not has_field_label(assessment_text, field):
            emit("RA-S09", "structural", rel_assessment, field, f"missing risk-control field: {field}")

for field in TRACEABILITY_FIELDS:
    if not has_populated_field(assessment_text, field):
        emit("RA-S10", "structural", rel_assessment, field, f"missing or empty traceability field: {field}")

# RA-C01 adapter rationale
adapter_rationale = find_field(values, "Adapter rationale")
if adapter_required in YES_NO:
    if not adapter_rationale or not adapter_rationale.strip():
        emit(
            "RA-C01",
            "consistency",
            rel_assessment,
            "Adapter rationale",
            "adapter rationale is mandatory for Adapter required Yes and No",
        )

# RA-C02 DEC-0036 for adapter-related
# Adapter-related when Adapter required is Yes, or disposition is Reuse/Wrap/Extend with external engine language.
adapter_related = adapter_required == "Yes"
refs_blob = find_field(values, "Related DEC / RFC / ADH references") or ""
refs_text = refs_blob + "\n" + assessment_text
if adapter_related and "DEC-0036" not in refs_text:
    emit(
        "RA-C02",
        "consistency",
        rel_assessment,
        "Related DEC / RFC / ADH references",
        "adapter-related assessment must reference DEC-0036",
    )

# RA-C04 conceptual example label
if re.search(r"conceptual example", assessment_text, flags=re.IGNORECASE):
    if CONCEPTUAL_LABEL not in assessment_text:
        emit(
            "RA-C04",
            "consistency",
            rel_assessment,
            "Conceptual example",
            "conceptual example lacks the exact required label",
            guidance=f'Use exactly: {CONCEPTUAL_LABEL}',
        )

# RA-C05 DEC/RFC/ADH existence
ref_ids = set(re.findall(r"\b(DEC-\d{4}|RFC-\d{4}|ADH-\d{4}-\d{3})\b", refs_text))
for ref in sorted(ref_ids):
    found = False
    if ref.startswith("DEC-"):
        candidates = list((REPO_ROOT / "docs/decisions").glob(f"{ref}*.md"))
        found = any(candidates)
    elif ref.startswith("RFC-"):
        candidates = list((REPO_ROOT / "docs/rfc").glob(f"{ref}*.md"))
        found = any(candidates)
    elif ref.startswith("ADH-"):
        candidates = list((REPO_ROOT / "docs/reviews/architecture-decision-handoffs").glob(f"{ref}*.md"))
        found = any(candidates)
    if not found:
        emit(
            "RA-C05",
            "consistency",
            rel_assessment,
            "Related DEC / RFC / ADH references",
            f"referenced record not found: {ref}",
        )

# RA-C06 adapter/contract identifier
adapter_id = find_field(values, "Adapter or contract identifier")
if adapter_required == "No":
    if adapter_id is None or not str(adapter_id).strip():
        emit(
            "RA-C06",
            "consistency",
            rel_assessment,
            "Adapter or contract identifier",
            "adapter or contract identifier is mandatory; use reserved literal none when Adapter required is No",
        )
    elif str(adapter_id).strip() != "none":
        emit(
            "RA-C06",
            "consistency",
            rel_assessment,
            "Adapter or contract identifier",
            "when Adapter required is No, identifier must be the reserved literal none",
        )
elif adapter_required == "Yes":
    if adapter_id is None or not str(adapter_id).strip() or str(adapter_id).strip() == "none":
        emit(
            "RA-C06",
            "consistency",
            rel_assessment,
            "Adapter or contract identifier",
            "adapter or contract identifier is mandatory when Adapter required is Yes",
        )

# RA-C07 Build triple
if disposition == "Build":
    for label in (
        "Why Reuse is insufficient",
        "Why Wrap is insufficient",
        "Why Extend is insufficient",
        "Protected Sovrunn differentiation and long-term ownership",
    ):
        if not has_field_label(assessment_text, label):
            emit(
                "RA-C07",
                "consistency",
                rel_assessment,
                label,
                f"Build disposition missing required field: {label}",
            )

# RA-C08 risk triple — verify each risk has preventive, detection, and corrective controls
has_risks_heading = has_field_label(assessment_text, "Applicable architecture risks")
has_matrix_heading = has_field_label(assessment_text, "Risk-control matrix")
risk_section_found = has_risks_heading or has_matrix_heading

if risk_section_found:
    risks_found = []
    
    # Check for risk-control matrix format first
    if has_matrix_heading:
        # Must have content after the heading, not just the heading itself
        if not has_heading_with_content(assessment_text, "Risk-control matrix"):
            emit("RA-C08", "consistency", rel_assessment, "Risk-control matrix",
                "empty Risk-control matrix heading - matrix must contain at least one risk row")
        else:
            # Parse the matrix table
            table_match = re.search(
                r"\|\s*Risk\s*\|\s*Preventive control\s*\|\s*Detection control\s*\|\s*Corrective path\s*\|.*?\n((?:\|[^|\n]+\|[^|\n]+\|[^|\n]+\|[^|\n]+\|.*?\n?)+)",
                assessment_text,
                flags=re.IGNORECASE | re.MULTILINE | re.DOTALL
            )
            
            if table_match:
                table_content = table_match.group(1)
                for line in table_content.split('\n'):
                    if line.strip() and '|' in line and not set(line.strip()) <= {'-', '|', ' '}:
                        cells = [cell.strip() for cell in line.strip('|').split('|')]
                        if len(cells) >= 4:
                            risk_name = cells[0].strip()
                            preventive = cells[1].strip()
                            detection = cells[2].strip()
                            corrective = cells[3].strip()
                            
                            if risk_name and not risk_name.lower() in {"risk", "---"}:
                                risks_found.append(risk_name)
                                
                                if not preventive:
                                    emit("RA-C08", "consistency", rel_assessment, risk_name, 
                                        f"risk '{risk_name}' missing preventive control")
                                if not detection:
                                    emit("RA-C08", "consistency", rel_assessment, risk_name,
                                        f"risk '{risk_name}' missing detection control")
                                if not corrective:
                                    emit("RA-C08", "consistency", rel_assessment, risk_name,
                                        f"risk '{risk_name}' missing corrective path")
            
            if not risks_found:
                emit("RA-C08", "consistency", rel_assessment, "Risk-control matrix",
                    "risk matrix header found but no valid risk rows - matrix must contain at least one complete risk")
    
    # Legacy format with separate control sections
    elif has_risks_heading:
        if not has_populated_field(assessment_text, "Applicable architecture risks"):
            emit("RA-C08", "consistency", rel_assessment, "Applicable architecture risks",
                "empty Applicable architecture risks - must list specific risks")
        else:
            # Require all control section headings for legacy format
            for label in RISK_CONTROL_FIELDS:
                if not has_field_label(assessment_text, label):
                    emit("RA-C08", "consistency", rel_assessment, label,
                        f"risk section present but missing control heading: {label}")
            
            # Assume risks exist if heading has content (legacy validation)
            risks_found = ["legacy format risks"]

# RA-C11 Phase 2 scope acknowledgement
ack = find_field(values, "Phase 2 non-goal acknowledgement") or ""
if not ack.strip() and not has_field_label(assessment_text, "Phase 2 non-goal acknowledgement"):
    emit(
        "RA-C11",
        "consistency",
        rel_assessment,
        "Phase 2 non-goal acknowledgement",
        "missing Phase 2 scope acknowledgement",
    )
elif "phase 2" not in assessment_text.lower() or "non-goal" not in assessment_text.lower():
    if not ack.strip():
        emit(
            "RA-C11",
            "consistency",
            rel_assessment,
            "Phase 2 non-goal acknowledgement",
            "missing Phase 2 scope acknowledgement",
        )

# RA-C12 duplicated schema in operational artifacts
schema_markers = [
    "### Existing mature solutions",
    "### Decision\nReuse / Wrap / Extend / Build",
    "### Adapter boundary required?",
]
for rel in OPERATIONAL_LINK_TARGETS:
    p = REPO_ROOT / rel
    if not p.is_file():
        continue
    text = p.read_text(encoding="utf-8", errors="replace")
    # Prohibited: fenced markdown that redefines the old mini-schema fields together
    if "### Existing mature solutions" in text and "### Adapter boundary required?" in text:
        emit(
            "RA-C12",
            "consistency",
            rel,
            "schema",
            "operational artifact redefines the canonical field schema instead of referencing it",
            guidance="Replace duplicated schema with a link to docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md",
        )

# RA-C14 canonical-source reference
for rel in OPERATIONAL_LINK_TARGETS:
    p = REPO_ROOT / rel
    if not p.is_file():
        continue
    text = p.read_text(encoding="utf-8", errors="replace")
    if "PHASE2_REUSE_ASSESSMENT_STANDARD.md" not in text:
        emit(
            "RA-C14",
            "consistency",
            rel,
            "canonical source",
            "required artifact does not reference docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md",
        )

# RA-C03 via authoritative phase2-scope-check.sh
if not SKIP_RAC03:
    scope_script = SCRIPT_DIR / "phase2-scope-check.sh"
    if scope_script.is_file():
        env = os.environ.copy()
        # Run from repo root; script discovers feature docs itself
        try:
            proc = subprocess.run(
                ["bash", str(scope_script), FEATURE],
                cwd=str(REPO_ROOT),
                capture_output=True,
                text=True,
            )
            if proc.returncode not in (0, 1):
                emit(
                    "CONFIG",
                    "config",
                    "scripts/phase2-scope-check.sh",
                    "RA-C03",
                    f"scope checker returned unexpected exit {proc.returncode}",
                )
            elif proc.returncode == 1:
                emit(
                    "RA-C03",
                    "consistency",
                    rel_assessment,
                    "Phase 2 scope",
                    "future-integration or blocked Phase 2 scope content appears outside allowed headings",
                    guidance="Move blocked phrases under deferred work, non-goals, or future-phase headings",
                )
        except OSError as exc:
            emit(
                "CONFIG",
                "config",
                "scripts/phase2-scope-check.sh",
                "RA-C03",
                f"failed to invoke scope checker: {exc}",
            )

# RA-C13 approval enforcement
# Skip if explicitly requested for assessment-only validation
if not SKIP_RAC13:
    # In strict mode, changed-file list is required for complete validation
    if MODE == "strict" and not CHANGED_FILES:
        emit(
            "CONFIG",
            "config",
            "RA-C13",
            "changed-files",
            "strict mode requires gate-supplied changed-file list for RA-C13 evaluation",
            guidance="Feature gate must supply --changed-files in strict mode, or use --skip-rac13 for assessment-only validation",
        )
        for d in sorted(diagnostics, key=lambda x: (x["path"], x["section"], x["rule_id"])):
            print(
                f"{d['severity']}|{d['rule_id']}|{d['layer']}|{d['feature']}|{d['path']}|{d['section']}|{d['message']}|{d['guidance']}"
            )
        sys.exit(2)

    impl_paths = []
    if CHANGED_FILES:
        cf = Path(CHANGED_FILES)
        if not cf.is_file():
            emit(
                "CONFIG",
                "config",
                str(cf),
                "changed-files",
                "changed-file list path is missing",
            )
            for d in sorted(diagnostics, key=lambda x: (x["path"], x["section"], x["rule_id"])):
                print(
                    f"{d['severity']}|{d['rule_id']}|{d['layer']}|{d['feature']}|{d['path']}|{d['section']}|{d['message']}|{d['guidance']}"
                )
            sys.exit(2)
        try:
            changed = [ln.strip() for ln in cf.read_text(encoding="utf-8").splitlines() if ln.strip()]
        except (OSError, UnicodeError) as exc:
            emit(
                "CONFIG",
                "config",
                str(cf),
                "changed-files",
                f"changed-file list unreadable: {exc}",
            )
            for d in sorted(diagnostics, key=lambda x: (x["path"], x["section"], x["rule_id"])):
                print(
                    f"{d['severity']}|{d['rule_id']}|{d['layer']}|{d['feature']}|{d['path']}|{d['section']}|{d['message']}|{d['guidance']}"
                )
            sys.exit(2)
        for path in changed:
            norm = path.lstrip("./")
            if any(norm == pref[:-1] or norm.startswith(pref) for pref in IMPLEMENTATION_PREFIXES):
                impl_paths.append(norm)

        if impl_paths:
            # Require Approved + human-approval evidence
            status = decision_status or find_field(values, "Decision status")
            if status != "Approved":
                emit(
                    "RA-C13",
                    "consistency",
                    rel_assessment,
                    "Decision status",
                    "implementation-attempt paths present without Approved decision status",
                    guidance=f"Implementation paths: {', '.join(impl_paths[:5])}",
                )
            else:
                # Structured approval-evidence record (RA-C13)
                def normalize_evidence_value(text: str) -> str:
                    """Trim, collapse whitespace, strip presentation list markers; preserve case."""
                    if text is None:
                        return ""
                    parts = []
                    for line in str(text).splitlines():
                        line = line.strip()
                        line = re.sub(r"^[-*]\s+", "", line)
                        line = re.sub(r"^[-*]\s+", "", line)
                        if line:
                            parts.append(line)
                    collapsed = " ".join(parts)
                    return re.sub(r"\s+", " ", collapsed).strip()

                def parse_evidence_field(text: str, *names: str) -> str:
                    """Parse Field: value from list, bold-list, table, or plain lines (body only)."""
                    # Prefer body fields after front matter
                    body = text
                    if text.startswith("---"):
                        # skip first front-matter block
                        rest = text.split("\n", 1)[1]
                        for i, line in enumerate(rest.splitlines(keepends=True)):
                            if line.rstrip("\r\n") == "---":
                                body = rest[sum(len(x) for x in rest.splitlines(keepends=True)[: i + 1]) :]
                                break
                    for name in names:
                        patterns = [
                            rf"(?im)^\s*[-*]\s+\*?\*?{re.escape(name)}\*?\*?\s*:\s*(.+?)\s*$",
                            rf"(?im)^\s*\*?\*?{re.escape(name)}\*?\*?\s*:\s*(.+?)\s*$",
                            rf"(?im)^\|\s*{re.escape(name)}\s*\|\s*([^|\n]+?)\s*\|",
                        ]
                        for pat in patterns:
                            m = re.search(pat, body)
                            if m:
                                return m.group(1).strip().strip("`")
                    return ""

                def resolve_evidence_path(assessment_text: str, values: dict) -> str:
                    ref = find_field(
                        values,
                        "Structured approval-evidence record",
                        "Approval evidence record",
                        "Reuse assessment approval evidence",
                    ) or ""
                    ref = ref.strip().strip("`").strip()
                    if not ref:
                        m = re.search(
                            r"(docs/reviews/reuse-assessments/FEATURE-\d{4}-approval-evidence\.md)",
                            assessment_text,
                        )
                        if m:
                            ref = m.group(1).strip()
                    return ref

                def emit_config_and_exit(path: str, section: str, message: str, guidance: str = ""):
                    emit("CONFIG", "config", path, section, message, guidance=guidance)
                    for d in sorted(diagnostics, key=lambda x: (x["path"], x["section"], x["rule_id"])):
                        print(
                            f"{d['severity']}|{d['rule_id']}|{d['layer']}|{d['feature']}|{d['path']}|{d['section']}|{d['message']}|{d['guidance']}"
                        )
                    sys.exit(2)

                def parse_evidence_front_matter(text: str):
                    """Parse exactly one opening/closing --- front-matter block at file start."""
                    if not (text.startswith("---\n") or text.startswith("---\r\n")):
                        return None, "missing opening front-matter delimiter"
                    rest = text.split("\n", 1)[1]
                    fm_lines = []
                    closed = False
                    body_started_at = None
                    offset = 0
                    for line in rest.splitlines(keepends=True):
                        if not closed and line.rstrip("\r\n") == "---":
                            closed = True
                            body_started_at = offset + len(line)
                            break
                        if not closed:
                            fm_lines.append(line.rstrip("\r\n"))
                        offset += len(line)
                    if not closed:
                        return None, "unterminated front-matter block"
                    # Reject a second top-level --- front-matter style block immediately after? not required
                    keys = {}
                    for raw in fm_lines:
                        if not raw.strip() or raw.strip().startswith("#"):
                            continue
                        if ":" not in raw:
                            return None, f"malformed front-matter line: {raw}"
                        key, val = raw.split(":", 1)
                        key = key.strip()
                        val = val.strip()
                        if not key:
                            return None, "empty front-matter key"
                        if key in keys:
                            return None, f"duplicate front-matter key: {key}"
                        if not val:
                            return None, f"empty front-matter value for key: {key}"
                        keys[key] = val
                    required_keys = [
                        "feature",
                        "evidence_type",
                        "approval_status",
                        "approval_date",
                        "approving_role",
                        "assessment_format_version",
                    ]
                    missing = [k for k in required_keys if k not in keys]
                    if missing:
                        return None, "missing front-matter keys: " + ", ".join(missing)
                    return keys, ""

                evidence_ref = resolve_evidence_path(assessment_text, values)
                expected_evidence_ref = f"docs/reviews/reuse-assessments/{FEATURE}-approval-evidence.md"
                if not evidence_ref:
                    emit(
                        "RA-C13",
                        "consistency",
                        rel_assessment,
                        "Human-approval evidence",
                        "Approved status lacks structured approval-evidence record reference",
                        guidance=f"Reference {expected_evidence_ref}",
                    )
                else:
                    cleaned_ref = evidence_ref.strip().strip("`").strip()

                    # Absolute / traversal / outside-repo => CONFIG exit 2; never read external files
                    if cleaned_ref.startswith("/") or re.match(r"^[A-Za-z]:[\\/]", cleaned_ref):
                        emit_config_and_exit(
                            cleaned_ref,
                            "approval-evidence",
                            "structured approval-evidence path must be repository-relative, not absolute",
                        )
                    if ".." in Path(cleaned_ref).parts:
                        emit_config_and_exit(
                            cleaned_ref,
                            "approval-evidence",
                            "structured approval-evidence path must not contain '..' traversal",
                        )

                    if cleaned_ref != expected_evidence_ref:
                        emit(
                            "RA-C13",
                            "consistency",
                            rel_assessment,
                            "Human-approval evidence",
                            "structured approval-evidence reference must equal "
                            f"{expected_evidence_ref}",
                            guidance=f"assessment referenced '{cleaned_ref}'",
                        )
                    else:
                        try:
                            repo_root_resolved = REPO_ROOT.resolve()
                            evidence_path = (REPO_ROOT / cleaned_ref).resolve()
                        except OSError as exc:
                            emit_config_and_exit(
                                cleaned_ref,
                                "approval-evidence",
                                f"structured approval-evidence path resolution failed: {exc}",
                            )
                        try:
                            evidence_path.relative_to(repo_root_resolved)
                        except ValueError:
                            emit_config_and_exit(
                                cleaned_ref,
                                "approval-evidence",
                                "structured approval-evidence path resolves outside the repository",
                            )
                        if not evidence_path.is_file():
                            emit_config_and_exit(
                                cleaned_ref,
                                "approval-evidence",
                                "structured approval-evidence file is missing or unresolvable",
                            )
                        try:
                            evidence_text = evidence_path.read_text(encoding="utf-8")
                        except (OSError, UnicodeError) as exc:
                            emit_config_and_exit(
                                cleaned_ref,
                                "approval-evidence",
                                f"structured approval-evidence file unreadable: {exc}",
                            )
                        if not evidence_text.strip():
                            emit_config_and_exit(
                                cleaned_ref,
                                "approval-evidence",
                                "structured approval-evidence file is empty or unparseable",
                            )

                        fm_keys, fm_err = parse_evidence_front_matter(evidence_text)
                        if fm_keys is None:
                            emit_config_and_exit(
                                cleaned_ref,
                                "approval-evidence",
                                f"structured approval-evidence front matter invalid: {fm_err}",
                            )

                        # Body fields (must not fall back to front-matter keys)
                        body_feature = normalize_evidence_value(
                            parse_evidence_field(evidence_text, "Feature")
                        )
                        body_evidence_type = normalize_evidence_value(
                            parse_evidence_field(evidence_text, "Evidence type")
                        )
                        body_status = normalize_evidence_value(
                            parse_evidence_field(evidence_text, "Approval status")
                        )
                        body_date = normalize_evidence_value(
                            parse_evidence_field(evidence_text, "Approval date")
                        )
                        body_approver = normalize_evidence_value(
                            parse_evidence_field(evidence_text, "Approver or approving role")
                        )
                        body_format = normalize_evidence_value(
                            parse_evidence_field(evidence_text, "Assessment format version")
                        )

                        fm_pairs = [
                            ("feature", normalize_evidence_value(fm_keys["feature"]), body_feature),
                            (
                                "evidence_type",
                                normalize_evidence_value(fm_keys["evidence_type"]),
                                body_evidence_type,
                            ),
                            (
                                "approval_status",
                                normalize_evidence_value(fm_keys["approval_status"]),
                                body_status,
                            ),
                            (
                                "approval_date",
                                normalize_evidence_value(fm_keys["approval_date"]),
                                body_date,
                            ),
                            (
                                "approving_role",
                                normalize_evidence_value(fm_keys["approving_role"]),
                                body_approver,
                            ),
                            (
                                "assessment_format_version",
                                normalize_evidence_value(fm_keys["assessment_format_version"]),
                                body_format,
                            ),
                        ]
                        fm_body_mismatches = [name for name, left, right in fm_pairs if left != right]
                        if fm_body_mismatches:
                            emit(
                                "RA-C13",
                                "consistency",
                                rel_assessment,
                                "Human-approval evidence",
                                "structured approval-evidence front matter conflicts with body: "
                                + ", ".join(fm_body_mismatches),
                            )
                        else:
                            # Active assessment fields
                            active_feature = FEATURE
                            active_assessment_path = normalize_evidence_value(rel_assessment)
                            active_format = normalize_evidence_value(
                                find_field(values, "reuse_assessment_format_version")
                                or assess_version
                                or ""
                            )
                            active_disposition = normalize_evidence_value(
                                find_field(values, "Disposition") or ""
                            )
                            active_sovrunn = normalize_evidence_value(
                                find_field(values, "Sovrunn-owned responsibility") or ""
                            )
                            active_reused = normalize_evidence_value(
                                find_field(
                                    values,
                                    "Reused or extended responsibility",
                                    "Reused or external responsibility",
                                )
                                or ""
                            )
                            active_boundary = normalize_evidence_value(
                                find_field(values, "Responsibility/control boundary") or ""
                            )
                            active_adh = ""
                            adh_m = re.search(r"\b(ADH-\d{4}-\d{3})\b", assessment_text)
                            if adh_m:
                                active_adh = adh_m.group(1)
                            controlling = find_field(
                                values,
                                "Approved ADH or assessment-review reference",
                                "Controlling ADH",
                                "controlling_handoff",
                            ) or ""
                            adh_from_field = re.search(r"\b(ADH-\d{4}-\d{3})\b", controlling)
                            if adh_from_field:
                                active_adh = adh_from_field.group(1)
                            active_approver = normalize_evidence_value(
                                find_field(values, "Approving person or role") or ""
                            )
                            active_approval_date = normalize_evidence_value(
                                find_field(values, "Approval date") or ""
                            )

                            # Evidence body fields used for assessment comparison
                            ev_feature = body_feature
                            ev_assessment_path = normalize_evidence_value(
                                parse_evidence_field(evidence_text, "Assessment artifact")
                            )
                            ev_format = body_format
                            ev_disposition = normalize_evidence_value(
                                parse_evidence_field(evidence_text, "Disposition")
                            )
                            ev_sovrunn = normalize_evidence_value(
                                parse_evidence_field(evidence_text, "Sovrunn-owned responsibility")
                            )
                            ev_reused = normalize_evidence_value(
                                parse_evidence_field(
                                    evidence_text, "Reused or extended responsibility"
                                )
                            )
                            ev_boundary = normalize_evidence_value(
                                parse_evidence_field(
                                    evidence_text, "Responsibility/control boundary"
                                )
                            )
                            ev_adh = normalize_evidence_value(
                                parse_evidence_field(evidence_text, "Controlling ADH")
                            )
                            ev_status = body_status
                            ev_approver = body_approver
                            ev_date = body_date

                            required_evidence_fields = {
                                "Feature": ev_feature,
                                "Assessment artifact": ev_assessment_path,
                                "Assessment format version": ev_format,
                                "Disposition": ev_disposition,
                                "Sovrunn-owned responsibility": ev_sovrunn,
                                "Reused or extended responsibility": ev_reused,
                                "Responsibility/control boundary": ev_boundary,
                                "Controlling ADH": ev_adh,
                                "Approval status": ev_status,
                                "Approver or approving role": ev_approver,
                                "Approval date": ev_date,
                            }
                            missing_ev = [k for k, v in required_evidence_fields.items() if not v]
                            if missing_ev:
                                emit(
                                    "RA-C13",
                                    "consistency",
                                    rel_assessment,
                                    "Human-approval evidence",
                                    "structured approval-evidence record missing required fields: "
                                    + ", ".join(missing_ev),
                                )
                            elif ev_status != "Approved":
                                emit(
                                    "RA-C13",
                                    "consistency",
                                    rel_assessment,
                                    "Human-approval evidence",
                                    f"structured approval-evidence status is '{ev_status}' (required: Approved)",
                                )
                            elif not active_approver or not active_approval_date:
                                emit(
                                    "RA-C13",
                                    "consistency",
                                    rel_assessment,
                                    "Human-approval evidence",
                                    "assessment missing Approving person or role and/or Approval date required for RA-C13 evidence comparison",
                                )
                            else:
                                mismatches = []
                                comparisons = [
                                    ("feature identifier", active_feature, ev_feature),
                                    (
                                        "assessment artifact path",
                                        active_assessment_path,
                                        ev_assessment_path,
                                    ),
                                    ("assessment format version", active_format, ev_format),
                                    ("disposition", active_disposition, ev_disposition),
                                    ("Sovrunn-owned responsibility", active_sovrunn, ev_sovrunn),
                                    (
                                        "reused or extended responsibility",
                                        active_reused,
                                        ev_reused,
                                    ),
                                    (
                                        "responsibility/control boundary",
                                        active_boundary,
                                        ev_boundary,
                                    ),
                                    ("controlling ADH", active_adh, ev_adh),
                                    ("approver or approving role", active_approver, ev_approver),
                                    ("approval date", active_approval_date, ev_date),
                                ]
                                for label, left, right in comparisons:
                                    if left != right:
                                        mismatches.append(label)

                                adh_files = (
                                    list(
                                        (
                                            REPO_ROOT
                                            / "docs/reviews/architecture-decision-handoffs"
                                        ).glob(f"{ev_adh}-*.md")
                                    )
                                    if ev_adh
                                    else []
                                )
                                if not adh_files:
                                    emit(
                                        "RA-C13",
                                        "consistency",
                                        rel_assessment,
                                        "Human-approval evidence",
                                        f"controlling ADH record not found for {ev_adh}",
                                    )
                                elif mismatches:
                                    emit(
                                        "RA-C13",
                                        "consistency",
                                        rel_assessment,
                                        "Human-approval evidence",
                                        "structured approval-evidence does not match assessment: "
                                        + ", ".join(mismatches),
                                    )

# Emit diagnostics
errors = [d for d in diagnostics if d["severity"] == "error"]
ordered = sorted(diagnostics, key=lambda x: (x["path"], x["section"], x["rule_id"]))
for d in ordered:
    print(
        f"{d['severity']}|{d['rule_id']}|{d['layer']}|{d['feature']}|{d['path']}|{d['section']}|{d['message']}|{d['guidance']}"
    )

# Config-class diagnostics force exit 2
if any(d["rule_id"] == "CONFIG" for d in diagnostics):
    sys.exit(2)

if errors:
    sys.exit(1)

print(f"PASS: reuse assessment validation for {FEATURE}")
sys.exit(0)
PY
