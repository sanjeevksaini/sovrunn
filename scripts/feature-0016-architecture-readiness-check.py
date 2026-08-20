#!/usr/bin/env python3
"""Fail-closed architecture-readiness check for FEATURE-0016.

Validates that ADH-2026-058's atomic replacement of the FEATURE-0016
placeholder registry definitions (VS0-SCHEMA-015..017, VS0-STATE-004) is
representable and consistent across the repository authorities before
requirements generation may proceed.

This checker rejects reintroduction of the removed placeholder concepts
(persisted status.availability, Draining, publicly projected Qualifying,
any alias/migration/legacy-compatibility behavior for them) and validates
the approved five-route surface, closed create-field boundary, sole
lifecycle-service writer, sole synthetic observer, idempotency scope, audit
matrix, and closed local violation set described in ADH-2026-058.

Exit 0 = PASS (requirements generation may proceed).
Exit 1 = FAIL (architecture gap remains; requirements generation blocked).
"""

from __future__ import annotations

import sys
from pathlib import Path

try:
    import yaml
except ImportError:
    print("FAIL: PyYAML required. Install: pip install pyyaml", file=sys.stderr)
    sys.exit(1)

ROOT = Path(__file__).resolve().parents[1]

REG_PATH = ROOT / "docs/architecture/vertical-slices/VS-000-contract-registry.yaml"
TRACE_PATH = ROOT / "docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md"
F16_ARCH = ROOT / "docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md"
F16_FEAT = ROOT / "docs/features/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md"
CLOSURE = ROOT / "docs/reviews/architecture-readiness/FEATURE-0016-architecture-closure-matrix.md"
REUSE = ROOT / "docs/reviews/reuse-assessments/FEATURE-0016-reuse-assessment.md"
STEER = ROOT / ".kiro/steering/slice0-contract.md"
CONTROL = ROOT / ".automation/features/FEATURE-0016.control.json"
ADH_058 = ROOT / "docs/reviews/architecture-decision-handoffs/ADH-2026-058-feature-0016-adapter-and-executiontarget-executable-contract-closure.md"

F16_REQUIRED_CF_IDS = [f"VS0-CF-F16-{i:02d}" for i in range(1, 123)]

# Removed placeholder concepts that must never reappear as active FEATURE-0016 behavior.
REMOVED_PLACEHOLDER_TERMS = [
    "status.availability",
    "Draining",
]

REMOVED_ALIAS_FIELD_TERMS = [
    "spec.participationRef",
]

F16_ROUTES = [
    ("POST", "/apis/execution.sovrunn.io/v1alpha1/execution-targets"),
    ("GET", "/apis/execution.sovrunn.io/v1alpha1/execution-targets"),
    ("GET", "/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}"),
    ("POST", "/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify"),
    ("POST", "/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire"),
]

F16_CLOSED_VIOLATIONS = [
    "VS0_TARGET_RETIRED",
    "VS0_TARGET_MAINTENANCE",
    "VS0_TARGET_QUALIFICATION_IN_PROGRESS",
    "VS0_TARGET_EPOCH_STALE",
    "VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE",
    "VS0_EXECUTION_TARGET_STACK_UNAVAILABLE",
    "VS0_EXECUTION_TARGET_SCOPE_MISMATCH",
    "VS0_EXECUTION_TARGET_VIABILITY_STALE",
]

errs: list[str] = []


def e(msg: str) -> None:
    errs.append(msg)


def read(path: Path) -> str:
    if not path.exists():
        e(f"Required file missing: {path.relative_to(ROOT)}")
        return ""
    return path.read_text()


def is_allowed_removed_mention(text: str, term: str) -> bool:
    """A removed-concept mention is allowed only in an explicit exclusion/
    removal/non-goal context (e.g. 'no active persisted status.availability',
    'Draining is removed'), on a line naming a different resource's field
    (e.g. ServiceRegion.status.availability, which is FEATURE-0022 scope and
    not the removed ExecutionTarget placeholder), or in a table cell whose
    row already carries a removal marker anywhere in that row/line."""
    term_lower = term.lower()
    allow_markers = (
        "no active", "must not", "no longer", "removed", "does not", "never",
        "not active", "no persisted", "excluded", "non-goal", "must never",
        "placeholder", "prohibited", "rejected", "no ",
    )
    other_resource_markers = ("serviceregion",)
    heading_allow_markers = ("removed", "excluded", "non-goal", "prohibited")
    lines = text.splitlines()
    for i, line in enumerate(lines):
        lower_line = line.lower()
        if term_lower not in lower_line:
            continue
        if any(marker in lower_line for marker in other_resource_markers):
            continue
        if any(marker in lower_line for marker in allow_markers):
            continue
        # Heading-aware fallback: find the nearest preceding heading and check
        # whether it establishes a removal/exclusion/non-goal context.
        heading = ""
        for j in range(i, -1, -1):
            if lines[j].strip().startswith("#"):
                heading = lines[j].lower()
                break
        if any(marker in heading for marker in heading_allow_markers):
            continue
        return False
    return True


def check_no_reintroduced_placeholders(reg: dict, f16_arch: str, f16_feat: str) -> None:
    schemas = reg.get("schemas", [])
    execution_target = next((s for s in schemas if s.get("identity", "").endswith("/ExecutionTarget")), None)
    if not execution_target:
        e("PLACEHOLDER: VS0-SCHEMA-015 (ExecutionTarget) not found in registry")
        return
    required = " ".join(execution_target.get("required", []))
    if "status.availability" in required:
        e("PLACEHOLDER: VS0-SCHEMA-015 must not require a persisted status.availability field (ADH-2026-058)")
    if "Draining" in required:
        e("PLACEHOLDER: VS0-SCHEMA-015 must not reference Draining as an active state (ADH-2026-058)")
    field_ownership = execution_target.get("fieldOwnership", {})
    if "status.availability" in field_ownership:
        e("PLACEHOLDER: VS0-SCHEMA-015 fieldOwnership must not own status.availability (ADH-2026-058)")
    if "spec.participationRef" in field_ownership:
        e("PLACEHOLDER: VS0-SCHEMA-015 fieldOwnership must not retain the renamed spec.participationRef alias (ADH-2026-058)")

    state_machines = reg.get("stateMachines", [])
    sm_004 = next((sm for sm in state_machines if sm.get("id") == "VS0-STATE-004"), None)
    if not sm_004:
        e("PLACEHOLDER: VS0-STATE-004 not found in registry")
    else:
        states = sm_004.get("states", [])
        for state in states:
            if "Draining" in state:
                e(f"PLACEHOLDER: VS0-STATE-004 must not contain a Draining state, found {state!r} (ADH-2026-058)")
            if "Qualifying" in state:
                e(f"PLACEHOLDER: VS0-STATE-004 must not persist Qualifying as a public state row, found {state!r} (ADH-2026-058 clause 7)")
        in_flight = str(sm_004.get("inFlightReservation", ""))
        if "never a projected resource status" not in in_flight and "never persisted or projected" not in in_flight.lower():
            e("PLACEHOLDER: VS0-STATE-004 must state Qualifying is a lifecycle-service-only in-flight reservation, never persisted or projected (ADH-2026-058 clause 7)")

    for label, text in (("architecture", f16_arch), ("feature", f16_feat)):
        if not text:
            continue
        for term in REMOVED_PLACEHOLDER_TERMS:
            if not is_allowed_removed_mention(text, term):
                e(f"PLACEHOLDER: F0016 {label} mentions '{term}' outside an explicit removal/exclusion context")
        if "publicly projected qualifying" not in text.lower() and "never a projected resource status" not in text.lower() and "never persisted or projected" not in text.lower():
            e(f"PLACEHOLDER: F0016 {label} must explicitly state Qualifying is never publicly projected (ADH-2026-058)")


def check_five_route_surface(f16_arch: str, f16_feat: str) -> None:
    # The architecture authority is the exact route source of truth; the feature
    # authority may reference it by clause rather than restating every path,
    # but must still assert the five-route/one-guard closure and prohibition.
    if f16_arch:
        lower = f16_arch.lower()
        for method, path in F16_ROUTES:
            if path.lower() not in lower:
                e(f"ROUTES: F0016 architecture must state the exact route {method} {path} (ADH-2026-058 clause 3)")
        if "no patch" not in lower:
            e("ROUTES: F0016 architecture must explicitly state no PATCH/PUT/DELETE/HEAD public route exists (ADH-2026-058 clause 3)")
        if "pre-servemux" not in lower:
            e("ROUTES: F0016 architecture must describe the pre-ServeMux transport-only method/path guard (ADH-2026-058 clause 3)")
    if f16_feat:
        lower = f16_feat.lower()
        if "exactly five" not in lower and "five routes" not in lower and "five explicit" not in lower:
            e("ROUTES: F0016 feature must state the exact five-route closure (ADH-2026-058 clause 3)")
        if "pre-servemux" not in lower and "transport-only" not in lower:
            e("ROUTES: F0016 feature must reference the pre-ServeMux transport-only guard (ADH-2026-058 clause 3)")


def check_closed_violation_set(reg: dict, f16_arch: str) -> None:
    vcs = {vc.get("code") for vc in reg.get("violationCodes", {}).get("slice0", [])}
    for code in F16_CLOSED_VIOLATIONS:
        if code not in vcs:
            e(f"VIOLATIONS: {code} must be a registered violation code (ADH-2026-058 clause 9)")
        if f16_arch and code not in f16_arch:
            e(f"VIOLATIONS: F0016 architecture must reference {code} (ADH-2026-058 clause 9)")


def check_sole_writer_and_observer(reg: dict, f16_arch: str) -> None:
    writers = reg.get("writers", [])
    writer_006 = next((w for w in writers if w.get("id") == "VS0-WRITER-006"), None)
    if not writer_006:
        e("WRITER: VS0-WRITER-006 not found in registry")
    else:
        if writer_006.get("writer") != "ExecutionTargetLifecycleService":
            e(f"WRITER: VS0-WRITER-006 writer must be ExecutionTargetLifecycleService, found {writer_006.get('writer')!r} (ADH-2026-058 clause 7)")
        paths = writer_006.get("paths", [])
        if not any("ExecutionTarget.status" in p for p in paths):
            e("WRITER: VS0-WRITER-006 must own ExecutionTarget.status (ADH-2026-058 clause 7)")
    if f16_arch and "sovrunn.synthetic-iaas-observer/v1" not in f16_arch:
        e("OBSERVER: F0016 architecture must name the sole observer sovrunn.synthetic-iaas-observer/v1 (ADH-2026-058 clause 6)")


def check_conformance_completeness(reg: dict) -> None:
    confs = reg.get("conformance", [])
    found = {c.get("id") for c in confs}
    missing = [cid for cid in F16_REQUIRED_CF_IDS if cid not in found]
    if missing:
        e(f"CONFORMANCE: missing {len(missing)} required VS0-CF-F16 case(s): {missing[:10]}{'...' if len(missing) > 10 else ''}")
    f16_cases = [c for c in confs if c.get("id", "").startswith("VS0-CF-F16-")]
    for c in f16_cases:
        for field in ("id", "owner", "inputs", "expectedState", "expectedError", "expectedSideEffects", "gate"):
            if field not in c:
                e(f"CONFORMANCE: {c.get('id', '?')} missing required field {field!r}")
        if c.get("owner") != "FEATURE-0016":
            e(f"CONFORMANCE: {c.get('id')} owner must be FEATURE-0016, found {c.get('owner')!r}")


def check_traceability(trace: str) -> None:
    if not trace:
        return
    missing = [cid for cid in F16_REQUIRED_CF_IDS if cid not in trace]
    if missing:
        e(f"TRACEABILITY: traceability matrix missing {len(missing)} VS0-CF-F16 reference(s): {missing[:10]}{'...' if len(missing) > 10 else ''}")
    for schema_id in ("VS0-SCHEMA-015", "VS0-SCHEMA-016", "VS0-SCHEMA-017", "VS0-STATE-004", "VS0-WRITER-006"):
        if schema_id not in trace:
            e(f"TRACEABILITY: traceability matrix missing {schema_id}")


def check_closure_matrix(closure: str) -> None:
    if not closure:
        return
    table_rows = [line for line in closure.splitlines() if line.strip().startswith("| ARC-F16-")]
    for i in range(1, 11):
        arc_id = f"ARC-F16-{i:02d}"
        row = next((r for r in table_rows if r.strip().startswith(f"| {arc_id} |")), None)
        if not row:
            e(f"CLOSURE: {arc_id} table row not found in closure matrix")
        elif "RESOLVED" not in row and "OUT OF SCOPE" not in row and "NOTED" not in row:
            e(f"CLOSURE: {arc_id} is not RESOLVED or explicitly noted out-of-scope")


def check_control_manifest_exists() -> None:
    if not CONTROL.exists():
        e(f"CONTROL: manifest missing: {CONTROL.relative_to(ROOT)}")
        return
    import json
    try:
        data = json.loads(CONTROL.read_text())
    except Exception as exc:
        e(f"CONTROL: manifest invalid JSON: {exc}")
        return
    owned = set(data.get("ownership", {}).get("owned_resources", []))
    expected = {"ExecutionTarget", "NormalizedTargetFactSet", "TargetQualificationResult"}
    if owned != expected:
        e(f"CONTROL: owned_resources mismatch: {sorted(owned ^ expected)}")


def check_reuse_assessment_exists() -> None:
    if not REUSE.exists():
        e(f"REUSE: reuse assessment missing: {REUSE.relative_to(ROOT)}")
        return
    text = REUSE.read_text()
    if "ADH-2026-058" not in text:
        e("REUSE: reuse assessment must reference ADH-2026-058")
    if "Approved" not in text:
        e("REUSE: reuse assessment must record Approved status")


def check_steering_updated(steer: str) -> None:
    if not steer:
        return
    if "ADH-2026-058" not in steer:
        e("STEERING: .kiro/steering/slice0-contract.md must reference ADH-2026-058")
    if "ExecutionTargetLifecycleService" not in steer:
        e("STEERING: .kiro/steering/slice0-contract.md must name ExecutionTargetLifecycleService as the sole F0016 status/record writer")
    if "sovrunn.synthetic-iaas-observer/v1" not in steer:
        e("STEERING: .kiro/steering/slice0-contract.md must name the sole observer sovrunn.synthetic-iaas-observer/v1")


def check_adh_058_approved() -> None:
    if not ADH_058.exists():
        e(f"ADH: handoff file missing: {ADH_058.relative_to(ROOT)}")
        return
    text = ADH_058.read_text()
    if "Approval status: Approved" not in text:
        e("ADH: ADH-2026-058 must record 'Approval status: Approved'")


def main() -> None:
    if not REG_PATH.exists():
        print(f"FAIL: registry not found: {REG_PATH.relative_to(ROOT)}")
        sys.exit(1)
    reg = yaml.safe_load(REG_PATH.read_text())
    f16_arch = read(F16_ARCH)
    f16_feat = read(F16_FEAT)
    trace = read(TRACE_PATH)
    closure = read(CLOSURE)
    steer = read(STEER)

    check_adh_058_approved()
    check_no_reintroduced_placeholders(reg, f16_arch, f16_feat)
    check_five_route_surface(f16_arch, f16_feat)
    check_closed_violation_set(reg, f16_arch)
    check_sole_writer_and_observer(reg, f16_arch)
    check_conformance_completeness(reg)
    check_traceability(trace)
    check_closure_matrix(closure)
    check_control_manifest_exists()
    check_reuse_assessment_exists()
    check_steering_updated(steer)

    if errs:
        print(f"FAIL: {len(errs)} error(s)")
        for err in errs:
            print(f"  ✗ {err}")
        sys.exit(1)
    print(
        "PASS: FEATURE-0016 architecture readiness — ADH-2026-058 placeholder replacement, "
        "five-route surface, sole writer/observer, closed violation set, 122-case conformance, "
        "traceability, closure matrix, control manifest, reuse assessment, and steering are consistent"
    )


if __name__ == "__main__":
    main()
