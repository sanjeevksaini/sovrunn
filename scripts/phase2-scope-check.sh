#!/usr/bin/env bash
set -euo pipefail

FEATURE="${1:-${FEATURE:-}}"

if [[ -z "$FEATURE" ]]; then
  echo "ERROR: feature id is required"
  exit 1
fi

case "$FEATURE" in
  FEATURE-0011|FEATURE-0012|FEATURE-0013|FEATURE-0014|FEATURE-0015|FEATURE-0016|FEATURE-0017|FEATURE-0018|FEATURE-0019|FEATURE-0020|FEATURE-0021|FEATURE-0022|FEATURE-0023|FEATURE-0024|FEATURE-0025|FEATURE-0026)
    ;;
  *)
    echo "WARN: $FEATURE is not a Phase 2 feature; skipping Phase 2 scope check"
    exit 0
    ;;
esac

echo "==> Checking Phase 2R scope boundaries for $FEATURE"

if [[ "$FEATURE" == "FEATURE-0014" ]] \
  && grep -q '"status": "implemented_and_merged"' .automation/state/FEATURE-0014.json \
  && grep -q '^status: historical-implementation$' docs/architecture/provider-neutral-resource-model.md \
  && grep -q 'ARCH-2026.08-PHASE2R-CANONICAL' docs/context/CURRENT_ARCHITECTURE_BASELINE.md; then
  echo "PASS: FEATURE-0014 is immutable merged alpha-model history under the Phase 2R canonical rebaseline"
  exit 0
fi

FILES=()
while IFS= read -r f; do FILES+=("$f"); done < <(find docs/features -maxdepth 1 -type f -name "${FEATURE}*.md" 2>/dev/null | sort || true)

CONTROL_MANIFEST=".automation/features/${FEATURE}.control.json"
if [[ -f "$CONTROL_MANIFEST" ]]; then
  SPEC_PATH="$(python3 - "$CONTROL_MANIFEST" <<'PY'
import json
import sys
from pathlib import Path

print(json.loads(Path(sys.argv[1]).read_text())["feature"]["spec_path"])
PY
)"
  if [[ -d "$SPEC_PATH" ]]; then
    while IFS= read -r f; do FILES+=("$f"); done < <(find "$SPEC_PATH" -maxdepth 1 -type f -name '*.md' | sort)
  fi
elif [[ -d ".kiro/specs" ]]; then
  # Resolve only the active feature's exact Kiro slug. A repository-wide grep
  # also selects earlier specs that merely mention this feature as a future
  # boundary, which caused retained FEATURE-0012 history to be judged as if it
  # were active FEATURE-0017 scope.
  SPEC_PATH="$(python3 - "$FEATURE" <<'PY'
import re
import sys
from pathlib import Path

feature = sys.argv[1]
index = Path("docs/features/FEATURE_INDEX.md")
if not index.is_file():
    raise SystemExit(0)

for line in index.read_text(encoding="utf-8").splitlines():
    if not re.match(rf"^\|\s*{re.escape(feature)}\s*\|", line):
        continue
    columns = [column.strip() for column in line.strip().strip("|").split("|")]
    if len(columns) < 5:
        raise SystemExit(0)
    slug = columns[4].strip("`")
    if re.fullmatch(r"[a-z0-9]+(?:-[a-z0-9]+)*", slug):
        print(f".kiro/specs/{slug}")
    raise SystemExit(0)
PY
)"
  if [[ -n "$SPEC_PATH" && -d "$SPEC_PATH" ]]; then
    while IFS= read -r f; do FILES+=("$f"); done < <(find "$SPEC_PATH" -maxdepth 1 -type f -name '*.md' | sort)
  fi
fi

if [[ ${#FILES[@]} -eq 0 ]]; then
  echo "WARN: no feature docs found for scope text scan"
  exit 0
fi

python3 - "$FEATURE" "${FILES[@]}" <<'PY'
import re
import sys
from pathlib import Path

feature = sys.argv[1]
files = [Path(p) for p in sys.argv[2:]]

blocked_patterns = [
    "real provider provisioning",
    "production provider provisioning",
    "real postgresql provisioning",
    "postgresql ha controller",
    "full opa integration",
    "full cedar integration",
    "full keycloak integration",
    "full vault integration",
    "temporal integration",
    "argo workflows integration",
    "global traffic execution",
    "autoscaling execution",
    "failover execution",
    "billing engine",
    "compliance engine",
    "autonomous remediation",
]

# Phase 2R canonical model: stale concepts must not appear as active requirements
phase2r_blocked_active_patterns = [
    "resourcepool",
    "providercapability",
    "effectivepolicycontext",
]

allowed_section_markers = [
    "non-goals",
    "non goals",
    "out of scope",
    "deferred",
    "not approved",
    "future phase",
    "later phase",
    "phase 2 does not",
]

blocked_section_markers = [
    "requirements",
    "acceptance criteria",
    "design",
    "implementation",
    "tasks",
    "scope",
    "in scope",
    "must",
    "shall",
]

violations = []

for path in files:
    try:
        lines = path.read_text(encoding="utf-8", errors="ignore").splitlines()
    except OSError:
        continue

    current_section = ""
    for idx, line in enumerate(lines, start=1):
        heading = re.match(r"^\s{0,3}#{1,6}\s+(.*)$", line)
        if heading:
            current_section = heading.group(1).strip().lower()

        lower = line.lower()
        for pattern in blocked_patterns:
            if pattern in lower:
                in_allowed_section = any(marker in current_section for marker in allowed_section_markers)
                negated_inline = any(token in lower for token in [
                    "do not", "does not", "must not", "should not", "not implement",
                    "deferred", "out of scope", "non-goal", "future phase", "later phase"
                ])

                # Allow explicit non-goal/deferred statements. These are good documentation, not violations.
                if in_allowed_section or negated_inline:
                    continue

                violations.append((str(path), idx, current_section or "<no heading>", pattern, line.strip()))

if violations:
    print(f"FAIL: potential Phase 2R scope violations for {feature}")
    print("The scope checker ignores phrases inside Non-goals/Out of Scope/Deferred sections and explicit negated statements.")
    for path, line_no, section, pattern, line in violations:
        print(f"{path}:{line_no}: section='{section}' pattern='{pattern}'")
        print(f"  {line}")
    sys.exit(1)

# Phase 2R stale concept check: ResourcePool, ProviderCapability, EffectivePolicyContext
# must not appear as active requirements (only in supersession/historical/non-goal context)
stale_violations = []
for path in files:
    try:
        lines = path.read_text(encoding="utf-8", errors="ignore").splitlines()
    except OSError:
        continue

    current_section = ""
    for idx, line in enumerate(lines, start=1):
        heading = re.match(r"^\s{0,3}#{1,6}\s+(.*)$", line)
        if heading:
            current_section = heading.group(1).strip().lower()

        lower = line.lower()
        for pattern in phase2r_blocked_active_patterns:
            if pattern in lower:
                in_allowed_section = any(marker in current_section for marker in allowed_section_markers)
                negated_inline = any(token in lower for token in [
                    "do not", "does not", "must not", "should not", "not implement",
                    "no resourcepool", "no providercapability", "supersed", "replaced",
                    "formerly", "migration", "historical", "retired", "dec-0042", "dec-0050"
                ])
                if in_allowed_section or negated_inline:
                    continue
                stale_violations.append((str(path), idx, current_section or "<no heading>", pattern, line.strip()))

if stale_violations:
    print(f"FAIL: Phase 2R stale concept violations for {feature}")
    print("ResourcePool, ProviderCapability, and EffectivePolicyContext must not appear as active requirements.")
    for path, line_no, section, pattern, line in stale_violations:
        print(f"{path}:{line_no}: section='{section}' stale='{pattern}'")
        print(f"  {line}")
    sys.exit(1)

print("PASS: Phase 2 scope boundaries")
PY
