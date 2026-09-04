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
TASK_HEADING = r"(?:Task\s+(\d+)\s*:|TASK-[A-Za-z0-9]+-(\d+)\s+[—:-])"


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
    for match in re.finditer(
        rf"(?ms)^### {TASK_HEADING}.*?(?=^### {TASK_HEADING}|^## [^#]|\Z)", text
    ):
        task = int(match.group(1) or match.group(2))
        if task in found:
            raise ValueError(f"duplicate task heading: {task}")
        found[task] = match.group(0).rstrip()
    return found


def commit_message(block: str) -> str | None:
    if re.search(r"(?m)^\*\*Commit message:\*\*\s+None\b", block):
        return None
    matches = list(
        re.finditer(r"(?ms)^\*\*Commit message:\*\*\s*\n```[^\n]*\n(.*?)^```\s*$", block)
    )
    if not matches:
        raise SystemExit("ERROR: task must contain one fenced Commit message: section")
    lines = [line.strip() for line in matches[0].group(1).splitlines() if line.strip()]
    if not lines:
        raise SystemExit("ERROR: task commit message must have a non-empty subject")
    return lines[0]


def section_paths(block: str, headings: tuple[str, ...]) -> list[str]:
    labels = "|".join(re.escape(heading) for heading in headings)
    matches = list(
        re.finditer(
            rf"(?ms)^\*\*(?:{labels}):\*\*\s*\n(.*?)(?=^\*\*[^\n]+:\*\*|^---\s*$|^#{{1,6}}\s|\Z)",
            block,
        )
    )
    if not matches:
        rendered = "/".join(headings)
        raise SystemExit(f"ERROR: task must contain a {rendered}: section")
    paths = []
    for match in matches:
        for line in match.group(1).splitlines():
            if not re.match(r"^\s*-\s+", line):
                continue
            for path in re.findall(r"`([^`]+)`", line):
                if (
                    path.startswith(("internal/", "cmd/", "api/", "tests/", "scripts/", "docs/", ".automation/"))
                    and (Path(path).suffix or path.endswith("/"))
                ):
                    paths.append(path)
    return paths


def task_writable_paths(block: str) -> list[str]:
    paths = section_paths(block, ("Writable paths", "Included writable paths"))
    paths.extend(section_paths(block, ("Tests", "Included tests")))
    return list(dict.fromkeys(paths))


def commit_task_ids(plan: dict[int, str]) -> list[int]:
    return [task for task in sorted(plan) if commit_message(plan[task]) is not None]


def verification_checkpoint_ids(plan: dict[int, str]) -> list[int]:
    return [task for task in sorted(plan) if commit_message(plan[task]) is None]


def verify_command(raw: str, feature: str) -> list[str]:
    parts = shlex.split(raw)
    if len(parts) < 2 or parts[0] != "make" or parts[1] not in ALLOWED_MAKE_TARGETS:
        raise SystemExit(f"ERROR: verification command is not allowlisted: {raw}")
    for part in parts[2:]:
        if part != f"FEATURE={feature}":
            raise SystemExit(f"ERROR: unsafe verification argument: {part}")
    return parts


def changed_paths() -> list[str]:
    output = run(["git", "status", "--porcelain=v1", "--untracked-files=all"], capture=True)
    return [line[3:] for line in output.splitlines() if len(line) > 3]


def allowed_resume_paths(state_path: Path, plan: dict[int, str], batches: list[list[int]]) -> set[str]:
    """Allow an interrupted first task to resume only inside its declared scope."""
    if not batches or not batches[0]:
        raise ValueError("no selected task for resume scope")
    task = batches[0][0]
    return {str(state_path.relative_to(ROOT)), *task_writable_paths(plan[task])}


def task_batches(
    ordered: list[int],
    last_raw: object,
    requested_start: int | None,
    stop_after: int | None,
    batch_limit: int,
    run_all: bool,
) -> list[list[int]]:
    if not ordered:
        raise ValueError("no numbered tasks found")
    if last_raw in (None, ""):
        expected_start = ordered[0]
    else:
        last = int(str(last_raw))
        if last not in ordered:
            raise ValueError(f"committed task {last} is not in the approved plan")
        if ordered.index(last) + 1 >= len(ordered):
            return []
        expected_start = ordered[ordered.index(last) + 1]
    start = requested_start if requested_start is not None else expected_start
    if start != expected_start:
        raise ValueError(f"sequence violation: expected Task {expected_start}, got {start}")
    selected = [task for task in ordered if task >= start]
    if stop_after is not None:
        selected = [task for task in selected if task <= stop_after]
    if not selected:
        raise ValueError("selected task batch is empty")
    if not run_all:
        selected = selected[:batch_limit]
    return [selected[index : index + batch_limit] for index in range(0, len(selected), batch_limit)]


def run_verification_checkpoint(feature: str, control: dict[str, object], task: int) -> None:
    """Run only manifest-controlled verification for a non-commit checkpoint.

    Checkpoints do not invoke Cursor, so they deliberately have no generated
    ``cursor-task-<n>.context.json``.  The generic Cursor boundary checker is
    therefore inapplicable here; the feature gate below supplies the
    checkpoint's clean-tree and approved-package verification.
    """
    print(f"\n=== {feature} verification checkpoint Task {task} ===", flush=True)
    guardrails = control["guardrails"]
    if not isinstance(guardrails, dict):
        raise SystemExit("ERROR: control guardrails must be an object")
    cursor = guardrails["cursor"]
    if not isinstance(cursor, dict):
        raise SystemExit("ERROR: control cursor guardrails must be an object")
    commands = cursor["verify_commands"]
    if not isinstance(commands, list):
        raise SystemExit("ERROR: control verify_commands must be a list")
    for command in commands:
        if not isinstance(command, str):
            raise SystemExit("ERROR: control verification command must be text")
        run(verify_command(command, feature))
    run(["make", "ff-guardrails", f"FEATURE={feature}"])
    run(["make", "ff-feature-gate", f"FEATURE={feature}"])


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    parser.add_argument("--start-task", type=int)
    parser.add_argument("--stop-after", type=int)
    parser.add_argument(
        "--all-tasks",
        action="store_true",
        help="run all remaining tasks while preserving manifest-sized verification batches",
    )
    args = parser.parse_args()
    run([str(ROOT / "scripts/feature-control.py"), "validate", "--feature", args.feature])
    control = json.loads((ROOT / f".automation/features/{args.feature}.control.json").read_text())
    state_path = ROOT / f".automation/state/{args.feature}.json"
    state = json.loads(state_path.read_text())
    if state.get("current_stage") != "cursor" or state.get("tasks_approval_token") != "APPROVED_FOR_CURSOR":
        raise SystemExit("ERROR: feature is not approved for Cursor execution")
    tasks_path = ROOT / state["spec_path"] / "tasks.md"
    plan = blocks(tasks_path.read_text())
    ordered = commit_task_ids(plan)
    checkpoints = verification_checkpoint_ids(plan)
    for task, block in plan.items():
        paths = task_writable_paths(block)
        if commit_message(block) is not None and not paths:
            raise SystemExit(f"ERROR: Task {task} must declare writable paths and tests")
    last_raw = state.get("last_committed_task")
    batch_limit = control["execution"]["tasks_per_run"]
    try:
        batches = task_batches(
            ordered,
            last_raw,
            args.start_task,
            args.stop_after,
            batch_limit,
            args.all_tasks,
        )
    except (TypeError, ValueError) as exc:
        raise SystemExit(f"ERROR: {exc}") from exc
    if not batches:
        if checkpoints:
            rendered = ", ".join(str(task) for task in checkpoints)
            print(
                f"{args.feature} has all approved commit tasks completed; "
                f"running verification-only Task {rendered}"
            )
            for checkpoint in checkpoints:
                run_verification_checkpoint(args.feature, control, checkpoint)
            print(f"{args.feature} controlled task run complete")
            return
        if args.all_tasks:
            print(f"{args.feature} already has all approved tasks committed")
            return
        raise SystemExit(f"ERROR: no approved successor after committed task {last_raw}")
    allowed_dirty = allowed_resume_paths(state_path, plan, batches)
    unexpected = [path for path in changed_paths() if path not in allowed_dirty]
    if unexpected:
        raise SystemExit("ERROR: working tree is not ready: " + ", ".join(unexpected))

    for selected in batches:
        print(f"\n=== {args.feature} controlled batch {selected} ===", flush=True)
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
            message = commit_message(plan[task])
            if message is None:
                raise SystemExit(f"ERROR: Task {task} is a verification checkpoint, not a commit task")
            run(
                [
                    "make",
                    "ff-commit-task",
                    f"FEATURE={args.feature}",
                    f"TASK={task}",
                    f"MESSAGE={message}",
                ]
            )
        print(f"{args.feature} task batch complete: {selected}")
    completed_final_commit = batches[-1][-1] == ordered[-1]
    if completed_final_commit:
        for task in checkpoints:
            run_verification_checkpoint(args.feature, control, task)
    if checkpoints:
        rendered = ", ".join(str(task) for task in checkpoints)
        print(
            f"{args.feature} controlled task run complete; Task {rendered} is verification-only "
            "and must not be sent to Cursor or committed"
        )
    else:
        print(f"{args.feature} controlled task run complete")


if __name__ == "__main__":
    main()
