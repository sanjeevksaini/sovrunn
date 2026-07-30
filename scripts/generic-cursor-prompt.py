#!/usr/bin/env python3
"""Render one guarded, token-bounded Cursor task prompt from feature control."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def task_blocks(text: str) -> dict[int, str]:
    blocks: dict[int, str] = {}
    pattern = re.compile(r"(?ms)^## Task (\d+)\b.*?(?=^## Task \d+\b|^## [^#]|\Z)")
    for match in pattern.finditer(text):
        blocks[int(match.group(1))] = match.group(0).rstrip()
    return blocks


def writable_paths(block: str) -> list[str]:
    lines = [line for line in block.splitlines() if line.startswith("Files:")]
    if len(lines) != 1:
        raise SystemExit("ERROR: task must contain exactly one Files: line")
    paths = re.findall(r"`([^`]+)`", lines[0])
    if not paths or len(paths) != len(set(paths)):
        raise SystemExit("ERROR: task Files: line must contain unique repository paths")
    for path in paths:
        parsed = Path(path)
        if parsed.is_absolute() or ".." in parsed.parts or parsed.parts[0] == ".git":
            raise SystemExit(f"ERROR: unsafe writable path: {path}")
    return paths


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    parser.add_argument("--task", required=True, type=int)
    parser.add_argument("--output", required=True)
    parser.add_argument("--manifest", required=True)
    args = parser.parse_args()

    state_path = ROOT / f".automation/state/{args.feature}.json"
    state = json.loads(state_path.read_text())
    control_path = ROOT / f".automation/features/{args.feature}.control.json"
    if state.get("control_approved_sha256") != sha256(control_path):
        raise SystemExit("ERROR: feature control differs from the founder-approved digest")
    if state.get("current_stage") != "cursor" or state.get("tasks_approval_token") != "APPROVED_FOR_CURSOR":
        raise SystemExit("ERROR: feature is not approved for Cursor")
    tasks_path = ROOT / state["spec_path"] / "tasks.md"
    if state.get("tasks_approved_sha256") != sha256(tasks_path):
        raise SystemExit("ERROR: tasks.md differs from its approved digest")
    blocks = task_blocks(tasks_path.read_text())
    if args.task not in blocks:
        raise SystemExit(f"ERROR: Task {args.task} not found")
    block = blocks[args.task]
    writable = writable_paths(block)

    manifest_path = ROOT / args.manifest
    subprocess.check_call(
        [
            str(ROOT / "scripts/feature-control.py"),
            "context",
            "--feature",
            args.feature,
            "--stage",
            "implementation",
            "--output",
            str(manifest_path.relative_to(ROOT)),
        ],
        cwd=ROOT,
    )
    manifest = json.loads(manifest_path.read_text())
    manifest.update(
        {
            "task": args.task,
            "task_block_sha256": hashlib.sha256(block.encode()).hexdigest(),
            "tasks_approved_sha256": state["tasks_approved_sha256"],
            "writable_paths": writable,
        }
    )
    manifest_path.write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n")
    fragment = subprocess.check_output(
        [
            str(ROOT / "scripts/feature-control.py"),
            "prompt-fragment",
            "--feature",
            args.feature,
            "--stage",
            "implementation",
        ],
        cwd=ROOT,
        text=True,
    )
    try:
        model = subprocess.check_output(
            [
                str(ROOT / "scripts/model-recommend.py"),
                "--tool",
                "cursor",
                "--task",
                str(args.task),
                "--tasks-path",
                str(tasks_path.relative_to(ROOT)),
            ],
            cwd=ROOT,
            text=True,
        )
    except subprocess.CalledProcessError as exc:
        raise SystemExit("ERROR: model recommendation failed") from exc
    template = (ROOT / "docs/prompts/cursor/generic-task.prompt.md").read_text()
    rendered = (
        template.replace("{{FEATURE_ID}}", args.feature)
        .replace("{{TASK_ID}}", str(args.task))
        .replace("{{MODEL_RECOMMENDATIONS}}", model.strip())
        .replace("{{CONTROL_FRAGMENT}}", fragment.strip())
        .replace("{{TASK_BLOCK}}", block)
        .replace("{{WRITABLE_PATHS}}", "\n".join(f"- `{path}`" for path in writable))
    )
    output = ROOT / args.output
    output.parent.mkdir(parents=True, exist_ok=True)
    output.write_text(rendered.rstrip() + "\n")
    print(output.relative_to(ROOT))


if __name__ == "__main__":
    main()
