#!/usr/bin/env python3
"""Validate that all Kiro spec stages retain their approved content digests."""

from __future__ import annotations

import argparse
import hashlib
import json
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def sha256(path: Path) -> str:
    if not path.is_file():
        raise ValueError(f"missing specification file: {path.relative_to(ROOT)}")
    return hashlib.sha256(path.read_bytes()).hexdigest()


def combined_digest(paths: list[Path]) -> str:
    digest = hashlib.sha256()
    for path in paths:
        digest.update(str(path.relative_to(ROOT)).encode())
        digest.update(b"\0")
        digest.update(path.read_bytes())
        digest.update(b"\0")
    return digest.hexdigest()


def validate_approvals(state: dict[str, object], requirements: Path, design: Path, tasks: Path) -> None:
    checks = {
        "requirements_approval_token": "APPROVED_FOR_DESIGN",
        "design_approval_token": "APPROVED_FOR_TASKS",
        "executable_plan_approval_token": "APPROVED_EXECUTABLE_PLAN",
        "tasks_approval_token": "APPROVED_FOR_CURSOR",
    }
    for key, expected in checks.items():
        if state.get(key) != expected:
            raise ValueError(f"{key} must be {expected}")
    if state.get("current_stage") != "cursor":
        raise ValueError("approved specification must be at the cursor stage")
    if state.get("requirements_approved_sha256") != sha256(requirements):
        raise ValueError("requirements.md differs from its approved digest")
    if state.get("design_approved_sha256") != sha256(design):
        raise ValueError("design.md differs from its approved digest")
    if state.get("executable_plan_approved_sha256") != combined_digest([requirements, design]):
        raise ValueError("requirements/design differ from the founder-approved executable plan")
    if state.get("tasks_approved_sha256") != sha256(tasks):
        raise ValueError("tasks.md differs from its approved digest")


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    args = parser.parse_args()
    subprocess.run(
        [str(ROOT / "scripts/feature-control.py"), "validate", "--feature", args.feature],
        cwd=ROOT,
        check=True,
    )
    state = json.loads((ROOT / f".automation/state/{args.feature}.json").read_text())
    spec = ROOT / str(state["spec_path"])
    try:
        validate_approvals(state, spec / "requirements.md", spec / "design.md", spec / "tasks.md")
    except ValueError as exc:
        raise SystemExit(f"ERROR: {exc}") from exc
    print(f"PASS: {args.feature} approved specification digests")


if __name__ == "__main__":
    main()
