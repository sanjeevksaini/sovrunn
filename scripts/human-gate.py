#!/usr/bin/env python3
"""Record explicit founder approvals for the three generic human gates."""

from __future__ import annotations

import argparse
import hashlib
import json
import subprocess
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def combined_digest(paths: list[Path]) -> str:
    digest = hashlib.sha256()
    for path in paths:
        digest.update(str(path.relative_to(ROOT)).encode())
        digest.update(b"\0")
        digest.update(path.read_bytes())
        digest.update(b"\0")
    return digest.hexdigest()


def save(path: Path, state: dict[str, object]) -> None:
    state["updated_at"] = datetime.now(timezone.utc).isoformat()
    path.write_text(json.dumps(state, indent=2, sort_keys=True) + "\n")


def approved_review(feature: str, stage: str, token: str) -> None:
    path = ROOT / f".automation/reviews/{feature}/{stage}.review.json"
    if not path.is_file():
        raise SystemExit(f"ERROR: missing machine review: {path.relative_to(ROOT)}")
    review = json.loads(path.read_text())
    if review.get("status") != "APPROVED" or review.get("approval_token") != token:
        raise SystemExit(f"ERROR: {stage} machine review is not approved with {token}")


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    parser.add_argument("--gate", required=True, choices=("architecture", "executable_plan", "final"))
    parser.add_argument("--approved-by", required=True)
    parser.add_argument("--approve", action="store_true")
    args = parser.parse_args()
    if not args.approve:
        raise SystemExit("ERROR: explicit --approve is required")
    subprocess.run(
        [str(ROOT / "scripts/feature-control.py"), "validate", "--feature", args.feature], cwd=ROOT, check=True
    )
    control = json.loads((ROOT / f".automation/features/{args.feature}.control.json").read_text())
    control_path = ROOT / f".automation/features/{args.feature}.control.json"
    state_path = ROOT / f".automation/state/{args.feature}.json"
    state = json.loads(state_path.read_text())
    now = datetime.now(timezone.utc).isoformat()
    if args.gate == "architecture":
        paths = [
            ROOT / control["feature"]["feature_file"],
            ROOT / control["feature"]["architecture"],
            *(ROOT / path for path in control["feature"]["handoffs"]),
        ]
        state.update(
            {
                "architecture_approval_token": "APPROVED_ARCHITECTURE",
                "architecture_approved_by": args.approved_by,
                "architecture_approved_at": now,
                "architecture_approved_sha256": combined_digest(paths),
                "control_approved_sha256": sha256(control_path),
                "current_stage": "requirements",
                "status": "architecture_approved",
                "human_gate_required": False,
            }
        )
    elif args.gate == "executable_plan":
        if state.get("control_approved_sha256") != sha256(control_path):
            raise SystemExit("ERROR: feature control changed after architecture approval")
        approved_review(args.feature, "requirements", "APPROVED_FOR_DESIGN")
        approved_review(args.feature, "design", "APPROVED_FOR_TASKS")
        spec = ROOT / state["spec_path"]
        requirements, design = spec / "requirements.md", spec / "design.md"
        packet_path = ROOT / f".automation/reports/{args.feature}/executable-plan-review.json"
        if not packet_path.is_file():
            raise SystemExit("ERROR: executable-plan review packet is missing")
        packet = json.loads(packet_path.read_text())
        if packet.get("requirements_sha256") != sha256(requirements) or packet.get("design_sha256") != sha256(design):
            raise SystemExit("ERROR: executable-plan review packet does not match current requirements/design")
        if state.get("requirements_approved_sha256") != sha256(requirements):
            raise SystemExit("ERROR: requirements differ from machine-approved digest")
        if state.get("design_approved_sha256") != sha256(design):
            raise SystemExit("ERROR: design differs from machine-approved digest")
        state.update(
            {
                "executable_plan_approval_token": "APPROVED_EXECUTABLE_PLAN",
                "executable_plan_approved_by": args.approved_by,
                "executable_plan_approved_at": now,
                "executable_plan_approved_sha256": combined_digest([requirements, design]),
                "current_stage": "tasks",
                "status": "executable_plan_approved",
                "human_gate_required": False,
            }
        )
    else:
        if state.get("control_approved_sha256") != sha256(control_path):
            raise SystemExit("ERROR: feature control changed after architecture approval")
        if state.get("tasks_approval_token") != "APPROVED_FOR_CURSOR":
            raise SystemExit("ERROR: tasks are not approved for Cursor")
        subprocess.run([str(ROOT / "scripts/final-verify.sh"), "--feature", args.feature], cwd=ROOT, check=True)
        state.update(
            {
                "final_approval_token": "APPROVED_FOR_MERGE",
                "final_approved_by": args.approved_by,
                "final_approved_at": now,
                "final_approved_commit": subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=ROOT, text=True).strip(),
                "current_stage": "final",
                "status": "approved_for_merge",
                "human_gate_required": False,
            }
        )
    save(state_path, state)
    print(f"PASS: {args.feature} {args.gate} approved by {args.approved_by}")


if __name__ == "__main__":
    main()
