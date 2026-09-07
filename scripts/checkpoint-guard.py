#!/usr/bin/env python3
"""Deterministic Feature Factory implementation-checkpoint guard (ADH-2026-077).

Once a feature has a committed Cursor task, its control manifest declares an
`implementation_checkpoint`. While that checkpoint is active this guard fails
closed for the delivery actions ADH-2026-077 prohibits:

- normal `requirements`/`design` Kiro stages and their LLM review routing;
- any stage-state overwrite that would leave the checkpoint; and
- any Cursor execution where the frozen design or requirements hash changed
  without an approved semantic-architecture reentry.

The guard never inspects or alters product semantics. It reads the control
manifest checkpoint block, the feature state, and the frozen spec hashes only.
A feature without a `checkpoint` block is unaffected: every mode passes.
"""

from __future__ import annotations

import argparse
import hashlib
import json
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def fail(message: str) -> None:
    print(f"CHECKPOINT_BLOCK: {message}", file=sys.stderr)
    raise SystemExit(3)


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def load_control(feature: str) -> dict:
    path = ROOT / f".automation/features/{feature}.control.json"
    if not path.is_file():
        raise SystemExit(f"ERROR: missing control manifest: {path}")
    return json.loads(path.read_text())


def load_state(feature: str) -> dict:
    path = ROOT / f".automation/state/{feature}.json"
    if not path.is_file():
        raise SystemExit(f"ERROR: missing feature state: {path}")
    return json.loads(path.read_text())


def checkpoint_of(control: dict) -> dict | None:
    checkpoint = control.get("checkpoint")
    if checkpoint is None:
        return None
    if checkpoint.get("active") is not True:
        return None
    return checkpoint


def spec_path(control: dict) -> Path:
    return ROOT / control["feature"]["spec_path"]


def assert_frozen(control: dict, checkpoint: dict) -> None:
    spec = spec_path(control)
    design = spec / "design.md"
    requirements = spec / "requirements.md"
    if not design.is_file():
        fail(f"frozen design.md missing: {design.relative_to(ROOT)}")
    if not requirements.is_file():
        fail(f"frozen requirements.md missing: {requirements.relative_to(ROOT)}")
    design_sha = sha256(design)
    requirements_sha = sha256(requirements)
    if design_sha != checkpoint["frozen_design_sha256"]:
        fail(
            "frozen design.md changed without approved semantic-architecture reentry "
            f"(got {design_sha}, expected {checkpoint['frozen_design_sha256']})"
        )
    if requirements_sha != checkpoint["frozen_requirements_sha256"]:
        fail(
            "frozen requirements.md changed without approved semantic-architecture reentry "
            f"(got {requirements_sha}, expected {checkpoint['frozen_requirements_sha256']})"
        )


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    parser.add_argument(
        "--mode",
        required=True,
        choices=(
            "status",
            "assert-stage-allowed",
            "assert-review-allowed",
            "assert-frozen",
            "assert-cursor-allowed",
        ),
    )
    parser.add_argument("--stage")
    parser.add_argument(
        "--reentry-handoff",
        default="",
        help="approved semantic-architecture reentry handoff id that permits leaving the checkpoint",
    )
    args = parser.parse_args()

    control = load_control(args.feature)
    checkpoint = checkpoint_of(control)

    if args.mode == "status":
        if checkpoint is None:
            print(f"INACTIVE: {args.feature} has no active implementation checkpoint")
        else:
            print(
                f"ACTIVE: {args.feature} implementation checkpoint "
                f"({checkpoint['handoff']}); last_committed_task="
                f"{checkpoint['last_committed_task']} current_task={checkpoint['current_task']}"
            )
        return

    if checkpoint is None:
        # No active checkpoint: this control has no delivery restriction.
        print(f"PASS: {args.feature} has no active checkpoint; {args.mode} unrestricted")
        return

    # An explicit, approved semantic-architecture reentry may lift the block.
    reentry = args.reentry_handoff.strip()
    if reentry:
        # Reentry is only honored when it names a real, approved handoff file and
        # is not the checkpoint handoff itself. This mirrors ADH-2026-077 §2:
        # leaving the checkpoint requires a separately approved handoff.
        if reentry == checkpoint["handoff"]:
            fail("reentry handoff must differ from the checkpoint handoff itself")
        handoff_glob = list(
            (ROOT / "docs/reviews/architecture-decision-handoffs").glob(f"{reentry}-*.md")
        )
        if not handoff_glob:
            fail(f"reentry handoff not found: {reentry}")
        print(f"PASS: approved semantic-architecture reentry {reentry} lifts checkpoint for {args.mode}")
        return

    if args.mode == "assert-stage-allowed":
        if not args.stage:
            raise SystemExit("ERROR: --stage is required for assert-stage-allowed")
        if args.stage in set(checkpoint["blocked_stages"]):
            fail(
                f"stage '{args.stage}' is blocked by the active implementation checkpoint "
                f"({checkpoint['handoff']}); use a task-only amendment "
                f"(scripts/checkpoint-task-amendment.sh) or an approved semantic-architecture reentry"
            )
        print(f"PASS: stage '{args.stage}' permitted under checkpoint {checkpoint['handoff']}")
        return

    if args.mode == "assert-review-allowed":
        if not args.stage:
            raise SystemExit("ERROR: --stage is required for assert-review-allowed")
        if args.stage in set(checkpoint["blocked_review_stages"]):
            fail(
                f"'{args.stage}' review routing is blocked by the active implementation checkpoint "
                f"({checkpoint['handoff']}); post-checkpoint {args.stage} LLM review is not a delivery authority"
            )
        print(f"PASS: review stage '{args.stage}' permitted under checkpoint {checkpoint['handoff']}")
        return

    if args.mode == "assert-frozen":
        assert_frozen(control, checkpoint)
        print(f"PASS: frozen design/requirements intact under checkpoint {checkpoint['handoff']}")
        return

    if args.mode == "assert-cursor-allowed":
        # Cursor may resume only at the amended task strictly greater than the
        # last committed task, and only when the frozen baseline is intact.
        assert_frozen(control, checkpoint)
        state = load_state(args.feature)
        current_task = state.get("current_task")
        try:
            current_task_num = int(current_task)
        except (TypeError, ValueError):
            fail(f"state current_task is not a task number: {current_task!r}")
        if current_task_num <= int(checkpoint["last_committed_task"]):
            fail(
                f"Cursor may not resume at task {current_task_num}; it must be greater than "
                f"last_committed_task={checkpoint['last_committed_task']}"
            )
        tasks = spec_path(control) / "tasks.md"
        tasks_sha = sha256(tasks)
        if state.get("checkpoint_task_amendment_approval_token") != "APPROVED_FOR_CURSOR":
            fail(
                "Cursor requires an explicit checkpoint task-amendment approval; "
                "run scripts/checkpoint-task-amendment-approve.py after preflight"
            )
        if state.get("checkpoint_task_amendment_tasks_sha256") != tasks_sha:
            fail(
                "checkpoint task-amendment approval does not match current tasks.md; "
                "rerun preflight and obtain a fresh approval"
            )
        if state.get("tasks_approval_token") != "APPROVED_FOR_CURSOR" or state.get("tasks_approved_sha256") != tasks_sha:
            fail("Cursor requires a matching approved tasks.md digest")
        print(
            f"PASS: Cursor resume at task {current_task_num} permitted under checkpoint {checkpoint['handoff']}"
        )
        return


if __name__ == "__main__":
    main()
