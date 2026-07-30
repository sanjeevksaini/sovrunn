#!/usr/bin/env python3
"""Execute approved FEATURE-0014 tasks sequentially with fail-closed validation."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
TASKS = ROOT / ".kiro/specs/provider-neutral-resource-model/tasks.md"
STATE = ROOT / ".automation/state/FEATURE-0014.json"
ALWAYS_ALLOWED = {".automation/state/FEATURE-0014.json"}


def run(args: list[str], *, capture: bool = False) -> str:
    print("+", " ".join(args), flush=True)
    result = subprocess.run(args, cwd=ROOT, text=True, capture_output=capture)
    if result.returncode:
        if capture:
            print(result.stdout, end="")
            print(result.stderr, end="", file=sys.stderr)
        raise SystemExit(result.returncode)
    return result.stdout if capture else ""


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def blocks() -> dict[int, str]:
    text = TASKS.read_text()
    found = {}
    for match in re.finditer(r"(?ms)^## Task (\d+)\b.*?(?=^## Task \d+\b|^## Requirement-to-task ledger)", text):
        found[int(match.group(1))] = match.group(0).rstrip()
    return found


def paths(block: str) -> list[str]:
    line = next(line for line in block.splitlines() if line.startswith("Files:"))
    return re.findall(r"`([^`]+)`", line)


def commit_message(block: str) -> str:
    line = next(line for line in block.splitlines() if line.startswith("Commit message:"))
    match = re.search(r"`([^`]+)`", line)
    if not match:
        raise SystemExit("task commit message is missing")
    return match.group(1)


def changed_paths() -> list[str]:
    output = run(["git", "status", "--porcelain"], capture=True)
    return [line[3:] for line in output.splitlines() if len(line) > 3]


def allowed(path: str, rules: list[str]) -> bool:
    if path in ALWAYS_ALLOWED:
        return True
    return any(path == rule or (rule.endswith("/") and path.startswith(rule)) for rule in rules)


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--start-task", type=int, default=1)
    parser.add_argument("--stop-after", type=int)
    args = parser.parse_args()

    state = json.loads(STATE.read_text())
    if state.get("current_stage") != "cursor" or state.get("tasks_approval_token") != "APPROVED_FOR_CURSOR":
        raise SystemExit("FEATURE-0014 is not approved for cursor execution")
    if state.get("tasks_approved_sha256") != sha256(TASKS):
        raise SystemExit("tasks.md differs from its approved digest")
    plan = blocks()
    expected_tasks = list(range(1, 11)) + list(range(12, 23))
    if sorted(plan) != expected_tasks:
        raise SystemExit("expected Tasks 1–10 and 12–22; Task 11 is intentionally absent")
    last = state.get("last_committed_task")
    if last is None:
        expected_start = expected_tasks[0]
    else:
        last_id = int(last)
        try:
            expected_start = expected_tasks[expected_tasks.index(last_id) + 1]
        except (ValueError, IndexError):
            raise SystemExit(f"no approved successor after committed task {last_id}")
    if args.start_task != expected_start:
        raise SystemExit(f"sequence violation: start task must be {expected_start}, got {args.start_task}")

    initial = changed_paths()
    unexpected = [p for p in initial if p not in ALWAYS_ALLOWED]
    if unexpected:
        raise SystemExit(f"working tree is not ready; commit or resolve: {unexpected}")

    for task_id in sorted(plan):
        if task_id < args.start_task:
            continue
        if args.stop_after is not None and task_id > args.stop_after:
            break
        block = plan[task_id]
        writable = paths(block)
        print(f"\n=== FEATURE-0014 Task {task_id} ===", flush=True)
        run(["make", "ff-cursor-task-auto", "FEATURE=FEATURE-0014", f"TASK={task_id}"])
        completed_state = json.loads(STATE.read_text())
        if completed_state.get("status") != f"cursor_task_{task_id}_completed":
            raise SystemExit(
                f"Task {task_id} has no verified COMPLETE receipt; "
                f"status={completed_state.get('status')!r}"
            )
        changed = changed_paths()
        outside = [p for p in changed if not allowed(p, writable)]
        if outside:
            raise SystemExit(f"Task {task_id} changed out-of-scope paths: {outside}")
        run(["make", "ff-guardrails", "FEATURE=FEATURE-0014"])
        run(["make", "feature-0014-cursor-boundary-check", f"TASK={task_id}"])
        if task_id == 22:
            run([
                "make",
                "feature-0014-architecture-boundary-check",
                "STAGE=tasks",
                "MODE=execution",
            ])
        run(["make", "ff-verify", "FEATURE=FEATURE-0014"])
        run(["make", "ff-commit-task", "FEATURE=FEATURE-0014", f"TASK={task_id}", f"MESSAGE={commit_message(block)}"])

    print("FEATURE-0014 requested task sequence completed.")


if __name__ == "__main__":
    main()
