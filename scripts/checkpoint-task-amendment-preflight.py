#!/usr/bin/env python3
"""Deterministic task-only amendment preflight (ADH-2026-077 §3).

While an implementation checkpoint is active, the only permitted repair path
for uncommitted tasks is an `implementation_task_amendment`. This preflight
proves, deterministically and fail-closed, that a proposed amended `tasks.md`:

- changes only tasks numbered strictly greater than `last_committed_task`
  (committed Task 1..last_committed_task blocks must be byte-identical to the
  approved/committed baseline);
- gives every affected task writable paths for each production package/test it
  must create or modify;
- has no forward dependency on a later executable task;
- preserves the completed committed-task paths and acceptance evidence; and
- invokes the applicable implementation checks (verification commands).

It never inspects or changes product semantics. It compares the working
`tasks.md` against the committed baseline blob and the frozen requirements/
design hashes recorded in the checkpoint. On success it prints the amended
`tasks.md` sha256 so a fresh hash-bound task/executable-plan approval can be
recorded separately (this script never records approval and never mutates
state).
"""

from __future__ import annotations

import argparse
import hashlib
import importlib.util
import json
import re
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

# Reuse the exact task-block and writable-path parsing used to build Cursor
# prompts so the preflight and the executor share one grammar.
_spec = importlib.util.spec_from_file_location(
    "generic_cursor_prompt", ROOT / "scripts/generic-cursor-prompt.py"
)
_gcp = importlib.util.module_from_spec(_spec)
_spec.loader.exec_module(_gcp)  # type: ignore[union-attr]

DEP_LINE = re.compile(r"(?m)^\*\*Dependencies:\*\*\s*(.+)$")
VERIFY_HEADING = re.compile(r"(?m)^\*\*Verification commands:\*\*\s*$")
ACCEPTANCE_HEADING = re.compile(r"(?m)^\*\*Acceptance criteria:\*\*\s*$")


def fail(message: str) -> None:
    print(f"AMENDMENT_PREFLIGHT_FAIL: {message}", file=sys.stderr)
    raise SystemExit(1)


def sha256_bytes(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def load_control(feature: str) -> dict:
    path = ROOT / f".automation/features/{feature}.control.json"
    if not path.is_file():
        raise SystemExit(f"ERROR: missing control manifest: {path}")
    return json.loads(path.read_text())


def committed_tasks_blob(feature: str, control: dict) -> str:
    """Return the committed baseline tasks.md content.

    The committed baseline is the tasks.md at the last commit that recorded the
    approved task-plan (HEAD of the feature branch). We read it from git to
    guarantee the completed Task blocks are compared against real committed
    bytes rather than the working tree.
    """
    spec = control["feature"]["spec_path"]
    rel = f"{spec}/tasks.md"
    result = subprocess.run(
        ["git", "show", f"HEAD:{rel}"],
        cwd=ROOT,
        text=True,
        capture_output=True,
    )
    if result.returncode != 0:
        fail(f"cannot read committed baseline tasks.md from HEAD:{rel}: {result.stderr.strip()}")
    return result.stdout


def task_dependency_numbers(block: str, unit_to_task: dict[str, int]) -> set[int]:
    """Resolve the executable-task numbers a task depends on.

    Dependencies are expressed as TASK-F18-NN atomic-unit ids. We map each unit
    id to the executable Task number that delivers it and return that set,
    excluding the task's own executable number.
    """
    match = DEP_LINE.search(block)
    if not match:
        return set()
    line = match.group(1)
    deps: set[int] = set()
    for unit in re.findall(r"TASK-[A-Z0-9]+-\d+", line):
        owner = unit_to_task.get(unit)
        if owner is not None:
            deps.add(owner)
    return deps


def build_unit_to_task(text: str, blocks: dict[int, str]) -> dict[str, int]:
    """Map each TASK-<unit> id to the executable Task number whose block
    introduces it via a 'Consolidated atomic task units' declaration, falling
    back to the first block that names it."""
    unit_to_task: dict[str, int] = {}
    consolidated = re.compile(
        r"(?m)^\*\*Consolidated atomic task units[^:]*:\*\*\s*(.+)$"
    )
    for task_num, block in sorted(blocks.items()):
        match = consolidated.search(block)
        if match:
            for unit in re.findall(r"TASK-[A-Z0-9]+-\d+", match.group(1)):
                unit_to_task.setdefault(unit, task_num)
    # Fallback: a task that declares no consolidated-units line still owns any
    # unit id it is the first to mention outside a Dependencies line.
    for task_num, block in sorted(blocks.items()):
        for line in block.splitlines():
            if line.startswith("**Dependencies:**"):
                continue
            for unit in re.findall(r"TASK-[A-Z0-9]+-\d+", line):
                unit_to_task.setdefault(unit, task_num)
    return unit_to_task


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    args = parser.parse_args()

    control = load_control(args.feature)
    checkpoint = control.get("checkpoint")
    if checkpoint is None or checkpoint.get("active") is not True:
        fail(f"{args.feature} has no active implementation checkpoint; amendment preflight does not apply")

    last_committed = int(checkpoint["last_committed_task"])
    spec = ROOT / control["feature"]["spec_path"]
    tasks_path = spec / "tasks.md"
    if not tasks_path.is_file():
        fail(f"missing tasks.md: {tasks_path.relative_to(ROOT)}")

    # 1. Frozen design/requirements must be intact before any amendment.
    guard = subprocess.run(
        [sys.executable, str(ROOT / "scripts/checkpoint-guard.py"), "--feature", args.feature, "--mode", "assert-frozen"],
        cwd=ROOT,
    )
    if guard.returncode != 0:
        fail("frozen design/requirements changed; a task amendment cannot proceed")

    working_text = tasks_path.read_text()
    baseline_text = committed_tasks_blob(args.feature, control)

    working_blocks = _gcp.task_blocks(working_text)
    baseline_blocks = _gcp.task_blocks(baseline_text)

    # 2. Committed Task 1..last_committed blocks must be byte-identical.
    for task_num in range(1, last_committed + 1):
        if task_num not in baseline_blocks:
            fail(f"committed baseline is missing Task {task_num}")
        if task_num not in working_blocks:
            fail(f"amended tasks.md removed committed Task {task_num}")
        if working_blocks[task_num] != baseline_blocks[task_num]:
            fail(
                f"amended tasks.md modified committed Task {task_num}; "
                f"tasks 1..{last_committed} are frozen and must not change"
            )

    # 3. Only uncommitted tasks (> last_committed) may differ.
    changed_tasks = sorted(
        num
        for num in working_blocks
        if baseline_blocks.get(num) != working_blocks.get(num)
    )
    for num in changed_tasks:
        if num <= last_committed:
            fail(f"Task {num} is committed and must not be amended")

    unit_to_task = build_unit_to_task(working_text, working_blocks)

    # 4. Per-affected-task deterministic proof.
    for num in changed_tasks:
        block = working_blocks[num]
        # 4a. Writable paths for every production package/test file.
        try:
            writable = _gcp.writable_paths(block)
        except SystemExit as exc:  # writable_paths fails closed with a message
            fail(f"Task {num}: {exc}")
        if not writable:
            fail(f"Task {num}: no writable paths declared")
        # 4b. No forward dependency on a later executable task.
        deps = task_dependency_numbers(block, unit_to_task)
        forward = sorted(d for d in deps if d > num)
        if forward:
            fail(
                f"Task {num}: forward dependency on later task(s) {forward}; "
                f"a task may consume only earlier-wave outputs or same-task paths"
            )
        # 4c. Verification commands must be present (applicable implementation checks).
        if not VERIFY_HEADING.search(block):
            fail(f"Task {num}: missing Verification commands: section")
        if not ACCEPTANCE_HEADING.search(block):
            fail(f"Task {num}: missing Acceptance criteria: section")

    amended_sha = sha256_bytes(tasks_path.read_bytes())
    print("PASS: task-only amendment preflight")
    print(f"last_committed_task={last_committed}")
    print(f"amended_tasks={sorted(changed_tasks)}")
    print(f"amended_tasks_sha256={amended_sha}")
    print(
        "Record a fresh hash-bound task/executable-plan approval bound to this sha256, "
        "then return current_stage=cursor and current_task="
        f"{checkpoint['current_task']}."
    )


if __name__ == "__main__":
    main()
