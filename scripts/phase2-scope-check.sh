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

echo "==> Checking Phase 2 scope boundaries for $FEATURE"

FILES=()
while IFS= read -r f; do FILES+=("$f"); done < <(find docs/features -maxdepth 1 -type f -name "${FEATURE}*.md" 2>/dev/null | sort || true)

if [[ -d ".kiro/specs" ]]; then
  while IFS= read -r f; do FILES+=("$f"); done < <(grep -Ril "$FEATURE" .kiro/specs/*/*.md 2>/dev/null | sort || true)
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
    print(f"FAIL: potential Phase 2 scope violations for {feature}")
    print("The scope checker ignores phrases inside Non-goals/Out of Scope/Deferred sections and explicit negated statements.")
    for path, line_no, section, pattern, line in violations:
        print(f"{path}:{line_no}: section='{section}' pattern='{pattern}'")
        print(f"  {line}")
    sys.exit(1)

print("PASS: Phase 2 scope boundaries")
PY
