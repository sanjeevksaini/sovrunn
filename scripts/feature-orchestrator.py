#!/usr/bin/env python3
"""Run manifest-controlled Cursor tasks in small, verified, committed batches."""

from __future__ import annotations

import argparse
import json
import re
import shlex
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
ALLOWED_MAKE_TARGETS = {"fmt", "vet", "test", "test-race", "build", "ff-verify", "ff-guardrails"}


def run(args: list[str], *, capture: bool = False) -> str:
    print("+", " ".join(args), flush=True)
    result = subprocess.run(args, cwd=ROOT, text=True, capture_output=capture)
    if result.returncode:
        if capture:
            print(result.stdout, end="")
            print(result.stderr, end="")
        raise SystemExit(result.returncode)
    return result.stdout if capture else ""


def blocks(text: str) -> dict[int, str]:
    found: dict[int, str] = {}
    for match in re.finditer(r"(?ms)^## Task (\d+)\b.*?(?=^## Task \d+\b|^## [^#]|\Z)", text):
        found[int(match.group(1))] = match.group(0).rstrip()
    return found


def commit_message(block: str) -> str:
    lines = [line for line in block.splitlines() if line.startswith("Commit message:")]
    if len(lines) != 1:
        raise SystemExit("ERROR: task must contain exactly one Commit message: line")
    match = re.fullmatch(r"Commit message:\s*`([^`]+)`\s*", lines[0])
    if not match:
        raise SystemExit("ERROR: task commit message must be one backtick-delimited value")
    return match.group(1)


def verify_command(raw: str, feature: str) -> list[str]:
    parts = shlex.split(raw)
    if len(parts) < 2 or parts[0] != "make" or parts[1] not in ALLOWED_MAKE_TARGETS:
        raise SystemExit(f"ERROR: verification command is not allowlisted: {raw}")
    for part in parts[2:]:
        if part != f"FEATURE={feature}":
            raise SystemExit(f"ERROR: unsafe verification argument: {part}")
    return parts


def changed_paths() -> list[str]:
    output = run(["git", "status", "--porcelain"], capture=True)
    return [line[3:] for line in output.splitlines() if len(line) > 3]


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    parser.add_argument("--start-task", type=int)
    parser.add_argument("--stop-after", type=int)
    args = parser.parse_args()
    run([str(ROOT / "scripts/feature-control.py"), "validate", "--feature", args.feature])
    control = json.loads((ROOT / f".automation/features/{args.feature}.control.json").read_text())
    state_path = ROOT / f".automation/state/{args.feature}.json"
    state = json.loads(state_path.read_text())
    if state.get("current_stage") != "cursor" or state.get("tasks_approval_token") != "APPROVED_FOR_CURSOR":
        raise SystemExit("ERROR: feature is not approved for Cursor execution")
    tasks_path = ROOT / state["spec_path"] / "tasks.md"
    plan = blocks(tasks_path.read_text())
    ordered = sorted(plan)
    if not ordered:
        raise SystemExit("ERROR: no numbered tasks found")
    last_raw = state.get("last_committed_task")
    if last_raw in (None, ""):
        expected_start = ordered[0]
    else:
        last = int(last_raw)
        if last not in ordered or ordered.index(last) + 1 >= len(ordered):
            raise SystemExit(f"ERROR: no approved successor after committed task {last}")
        expected_start = ordered[ordered.index(last) + 1]
    start = args.start_task if args.start_task is not None else expected_start
    if start != expected_start:
        raise SystemExit(f"ERROR: sequence violation: expected Task {expected_start}, got {start}")
    batch_limit = control["execution"]["tasks_per_run"]
    selected = [task for task in ordered if task >= start]
    if args.stop_after is not None:
        selected = [task for task in selected if task <= args.stop_after]
    selected = selected[:batch_limit]
    if not selected:
        raise SystemExit("ERROR: selected task batch is empty")
    allowed_dirty = {str(state_path.relative_to(ROOT))}
    unexpected = [path for path in changed_paths() if path not in allowed_dirty]
    if unexpected:
        raise SystemExit("ERROR: working tree is not ready: " + ", ".join(unexpected))

    for task in selected:
        print(f"\n=== {args.feature} Task {task} ===", flush=True)
        run(["make", "ff-cursor-task-auto", f"FEATURE={args.feature}", f"TASK={task}"])
        completed = json.loads(state_path.read_text())
        if completed.get("status") != f"cursor_task_{task}_completed":
            raise SystemExit(f"ERROR: Task {task} lacks a verified COMPLETE receipt")
        run(
            [
                str(ROOT / "scripts/generic-feature-boundary-check.py"),
                "--feature",
                args.feature,
                "--task",
                str(task),
                "--mode",
                "all",
            ]
        )
        for command in control["guardrails"]["cursor"]["verify_commands"]:
            run(verify_command(command, args.feature))
        run(["make", "ff-guardrails", f"FEATURE={args.feature}"])
        run(
            [
                "make",
                "ff-commit-task",
                f"FEATURE={args.feature}",
                f"TASK={task}",
                f"MESSAGE={commit_message(plan[task])}",
            ]
        )
    print(f"{args.feature} task batch complete: {selected}")


if __name__ == "__main__":
    main()
