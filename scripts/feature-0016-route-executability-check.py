#!/usr/bin/env python3
"""Fail-closed FEATURE-0016 route-executability audit.

Wired to the design admission gate (ADH-2026-058 §4.3, acceptance criterion
"route-executability audit is wired to the design admission gate"). Proves
that the approved five-route ExecutionTarget HTTP surface and its
pre-ServeMux transport-only guard have a complete outcome mapping before a
FEATURE-0016 design document is admitted.

If `.kiro/specs/adapter-boundary-and-executiontarget-qualification/design.md`
does not yet exist, this check validates only that the architecture-level
route/outcome contract is complete and internally consistent (pre-design
readiness mode), and reports that design-stage execution has not started.

Exit 0 = PASS. Exit 1 = FAIL (missing route outcome or precedence rule).
"""

from __future__ import annotations

import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
F16_ARCH = ROOT / "docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md"
DESIGN = ROOT / ".kiro/specs/adapter-boundary-and-executiontarget-qualification/design.md"

REQUIRED_ROUTE_PATHS = [
    "/apis/execution.sovrunn.io/v1alpha1/execution-targets",
    "/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}",
    "/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify",
    "/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire",
]

# Used only for the design-stage check, where a design document is expected to
# spell out method+path together (e.g. in a route table or handler signature).
REQUIRED_ROUTE_CLAUSES = [
    "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets",
    "GET /apis/execution.sovrunn.io/v1alpha1/execution-targets",
    "GET /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}",
    "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify",
    "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire",
]

REQUIRED_PRECEDENCE_STAGES = [
    "transport method/path guard",
    "authentication",
    "method/media/header-form validation",
    "phase-one safe references",
    "authorization",
    "safe access",
    "strict classification",
    "graph/reference validation",
    "replay/reservation",
    "current if-match version comparison",
    "lifecycle state",
    "atomic commit",
]

REQUIRED_TRANSPORT_GUARD_OUTCOMES = [
    "head on the collection path",
    "head on the item path",
    "head on either action path",
    "trailing-slash",
]

errs: list[str] = []


def e(msg: str) -> None:
    errs.append(msg)


def _normalize(text: str) -> str:
    """Collapse whitespace (including line-wrap newlines inside prose) so
    multi-line markdown clauses match a single-line search phrase."""
    return " ".join(text.split())


def check_architecture_route_contract() -> str:
    if not F16_ARCH.exists():
        e(f"Required file missing: {F16_ARCH.relative_to(ROOT)}")
        return ""
    text = F16_ARCH.read_text()
    lower = _normalize(text).lower()
    for path in REQUIRED_ROUTE_PATHS:
        if _normalize(path).lower() not in lower:
            e(f"Missing route path in architecture authority: {path}")
    for stage in REQUIRED_PRECEDENCE_STAGES:
        if stage not in lower:
            e(f"Missing precedence stage in architecture authority: {stage!r}")
    for outcome in REQUIRED_TRANSPORT_GUARD_OUTCOMES:
        if outcome not in lower:
            e(f"Missing transport-guard outcome in architecture authority: {outcome!r}")
    return text


def check_design_stage() -> None:
    if not DESIGN.exists():
        print(
            "NOTE: design.md does not yet exist for FEATURE-0016; "
            "route-executability audit ran in pre-design readiness mode only."
        )
        return
    text = DESIGN.read_text()
    lower = text.lower()
    for clause in REQUIRED_ROUTE_CLAUSES:
        if clause.lower() not in lower:
            e(f"design.md missing implementation path for route: {clause}")
    for stage in REQUIRED_PRECEDENCE_STAGES:
        if stage not in lower:
            e(f"design.md missing pipeline stage: {stage!r}")


def main() -> None:
    check_architecture_route_contract()
    check_design_stage()
    if errs:
        print(f"FAIL: {len(errs)} error(s)")
        for err in errs:
            print(f"  ✗ {err}")
        sys.exit(1)
    print("PASS: FEATURE-0016 route-executability audit — five-route surface, transport guard, and precedence pipeline are complete")


if __name__ == "__main__":
    main()
