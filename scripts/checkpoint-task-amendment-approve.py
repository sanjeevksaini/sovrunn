#!/usr/bin/env python3
"""Record explicit human approval for an ADH-2026-077 task-only amendment.

This command is intentionally separate from normal stage approval.  It binds a
preflight-passing amendment for uncommitted tasks to the current ``tasks.md``
digest without reopening requirements/design review or replacing the frozen
semantic baseline.
"""

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


def fail(message: str) -> None:
    raise SystemExit(f"ERROR: {message}")


def load_json(path: Path, label: str) -> dict:
    if not path.is_file():
        fail(f"missing {label}: {path.relative_to(ROOT)}")
    return json.loads(path.read_text())


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    parser.add_argument("--approved-by", required=True)
    parser.add_argument("--approve", action="store_true")
    args = parser.parse_args()

    if not args.approve:
        fail("explicit --approve is required")
    approved_by = args.approved_by.strip()
    if not approved_by:
        fail("--approved-by must not be empty")

    control_path = ROOT / f".automation/features/{args.feature}.control.json"
    state_path = ROOT / f".automation/state/{args.feature}.json"
    control = load_json(control_path, "feature control")
    state = load_json(state_path, "feature state")
    checkpoint = control.get("checkpoint")
    if not isinstance(checkpoint, dict) or checkpoint.get("active") is not True:
        fail("active implementation checkpoint is required")

    subprocess.run(
        [str(ROOT / "scripts/feature-control.py"), "validate", "--feature", args.feature],
        cwd=ROOT,
        check=True,
    )
    subprocess.run(
        [str(ROOT / "scripts/checkpoint-task-amendment.sh"), "--feature", args.feature],
        cwd=ROOT,
        check=True,
    )

    spec = ROOT / control["feature"]["spec_path"]
    tasks = spec / "tasks.md"
    task_sha = sha256(tasks)
    # The checkpoint's current_task is its entry point, not a permanent cursor
    # position.  A task-only amendment may be approved after later tasks have
    # completed, so preserve the state's next task rather than resetting work
    # back to the checkpoint entry point.
    try:
        resume_task = int(str(state.get("current_task")))
        state_last_committed = int(str(state.get("last_committed_task")))
    except (TypeError, ValueError):
        fail("state current_task and last_committed_task must be task numbers")
    checkpoint_last_committed = int(checkpoint["last_committed_task"])
    if resume_task <= state_last_committed:
        fail("state current_task must be greater than state last_committed_task")
    if state_last_committed < checkpoint_last_committed:
        fail("state last_committed_task cannot precede the checkpoint baseline")

    now = datetime.now(timezone.utc).isoformat()
    state.update(
        {
            "checkpoint_task_amendment_approval_token": "APPROVED_FOR_CURSOR",
            "checkpoint_task_amendment_approved_by": approved_by,
            "checkpoint_task_amendment_approved_at": now,
            "checkpoint_task_amendment_tasks_sha256": task_sha,
            "tasks_approval_source": "checkpoint_human_approval",
            "tasks_approval_token": "APPROVED_FOR_CURSOR",
            "tasks_approved_sha256": task_sha,
            "tasks_approved_at": now,
            "current_stage": "cursor",
            "current_task": str(resume_task),
            "status": "checkpoint_task_amendment_approved",
            "human_gate_required": False,
            "updated_at": now,
        }
    )
    state_path.write_text(json.dumps(state, indent=2, sort_keys=True) + "\n")
    print(
        f"PASS: {args.feature} checkpoint task amendment approved by {approved_by}; "
        f"tasks_sha256={task_sha}; resume_task={resume_task}"
    )


if __name__ == "__main__":
    main()
