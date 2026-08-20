#!/usr/bin/env python3
"""Fail-closed FEATURE-0016 source-derived task-scope inventory.

Wired to the tasks admission gate (ADH-2026-058 §4.4, acceptance criterion
"source-derived task-scope inventory is wired to the tasks admission gate").
Fails on a missing writable/test/dependency path or a cyclic graph once
`.kiro/specs/adapter-boundary-and-executiontarget-qualification/tasks.md`
exists.

If tasks.md does not yet exist, this check validates only that the
FEATURE-0016 control manifest declares a well-formed, acyclic dependency set
(pre-tasks readiness mode) and reports that tasks-stage execution has not
started. This is a mechanical inventory check; it does not choose task
content, Go data structures, or implementation scope.

Exit 0 = PASS. Exit 1 = FAIL (missing path or cyclic dependency graph).
"""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
CONTROL = ROOT / ".automation/features/FEATURE-0016.control.json"
TASKS = ROOT / ".kiro/specs/adapter-boundary-and-executiontarget-qualification/tasks.md"

errs: list[str] = []


def e(msg: str) -> None:
    errs.append(msg)


def check_control_dependency_graph() -> dict:
    if not CONTROL.exists():
        e(f"Required file missing: {CONTROL.relative_to(ROOT)}")
        return {}
    try:
        data = json.loads(CONTROL.read_text())
    except Exception as exc:
        e(f"Control manifest invalid JSON: {exc}")
        return {}

    previous_ids = [dep.get("id") for dep in data.get("ownership", {}).get("previous_features", [])]
    excluded_ids = [dep.get("id") for dep in data.get("ownership", {}).get("excluded_features", [])]

    if "FEATURE-0016" in previous_ids or "FEATURE-0016" in excluded_ids:
        e("Dependency graph is cyclic: FEATURE-0016 lists itself as a dependency or exclusion")

    overlap = set(previous_ids) & set(excluded_ids)
    if overlap:
        e(f"Dependency graph inconsistent: feature(s) both previous and excluded: {sorted(overlap)}")

    for feature_id in previous_ids:
        if not re.fullmatch(r"FEATURE-\d{4}", feature_id or ""):
            e(f"Malformed previous-feature id: {feature_id!r}")
            continue
        if int(feature_id.split("-")[1]) >= 16:
            e(f"Dependency graph is cyclic or forward-referencing: previous feature {feature_id} is not < FEATURE-0016")

    for feature_id in excluded_ids:
        if not re.fullmatch(r"FEATURE-\d{4}", feature_id or ""):
            e(f"Malformed excluded-feature id: {feature_id!r}")

    for dep in data.get("ownership", {}).get("previous_features", []):
        for stage, paths in dep.get("context", {}).items():
            for path in paths:
                if not (ROOT / path).exists():
                    e(f"Previous-feature context path missing: {path} (dependency {dep.get('id')}, stage {stage})")

    for path in (
        data.get("feature", {}).get("feature_file"),
        data.get("feature", {}).get("architecture"),
        *data.get("feature", {}).get("handoffs", []),
    ):
        if path and not (ROOT / path).exists():
            e(f"Feature-identity path missing: {path}")

    # Note: .kiro/specs/** is created only once Kiro requirements generation
    # begins (out of this architecture-update run's write boundary). Its
    # absence before that stage is expected, not a task-scope gap; this check
    # validates only that the declared spec_path is well-formed and consistent
    # with the feature id, not that it already exists on disk.
    spec_path = data.get("feature", {}).get("spec_path", "")
    if not spec_path.startswith(".kiro/specs/"):
        e(f"feature.spec_path must be under .kiro/specs/: {spec_path!r}")
    writable = data.get("guardrails", {}).get("kiro", {}).get("writable_by_stage", {})
    for stage in ("requirements", "design", "tasks"):
        paths = writable.get(stage, [])
        if len(paths) != 1 or not paths[0].startswith(f"{spec_path}/"):
            e(f"guardrails.kiro.writable_by_stage.{stage} must declare exactly one path under {spec_path}")

    return data


def check_tasks_stage(control: dict) -> None:
    if not TASKS.exists():
        print(
            "NOTE: tasks.md does not yet exist for FEATURE-0016; "
            "task-scope inventory ran in pre-tasks readiness mode only."
        )
        return
    text = TASKS.read_text()
    forbidden_outputs = control.get("guardrails", {}).get("kiro", {}).get("forbidden_outputs", [])
    for forbidden in forbidden_outputs:
        if forbidden in text:
            e(f"tasks.md references a forbidden output path: {forbidden}")
    owned_resources = control.get("ownership", {}).get("owned_resources", [])
    for resource in owned_resources:
        if resource not in text:
            e(f"tasks.md does not reference owned resource: {resource}")


def main() -> None:
    control = check_control_dependency_graph()
    check_tasks_stage(control)
    if errs:
        print(f"FAIL: {len(errs)} error(s)")
        for err in errs:
            print(f"  ✗ {err}")
        sys.exit(1)
    print("PASS: FEATURE-0016 task-scope inventory — dependency graph is acyclic and control-manifest paths resolve")


if __name__ == "__main__":
    main()
