#!/usr/bin/env python3
"""Enforce manifest context and single-file writes for Kiro spec stages."""

from __future__ import annotations

import argparse
import hashlib
import json
import subprocess
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


def changed_paths() -> list[str]:
    result = subprocess.run(
        ["git", "status", "--porcelain"], cwd=ROOT, text=True, capture_output=True, check=True
    )
    return [line[3:] for line in result.stdout.splitlines() if len(line) > 3]


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    parser.add_argument("--stage", required=True, choices=("requirements", "design", "tasks"))
    parser.add_argument("--mode", required=True, choices=("pre", "prompt", "post"))
    args = parser.parse_args()
    subprocess.run(
        [str(ROOT / "scripts/feature-control.py"), "validate", "--feature", args.feature], cwd=ROOT, check=True
    )
    state_path = ROOT / f".automation/state/{args.feature}.json"
    state = json.loads(state_path.read_text())
    control = json.loads((ROOT / f".automation/features/{args.feature}.control.json").read_text())
    control_path = ROOT / f".automation/features/{args.feature}.control.json"
    if state.get("control_approved_sha256") != sha256(control_path):
        raise SystemExit("FAIL: feature control differs from the founder-approved digest")
    spec = ROOT / control["feature"]["spec_path"]
    if args.stage == "requirements":
        architecture_paths = [
            ROOT / control["feature"]["feature_file"],
            ROOT / control["feature"]["architecture"],
            *(ROOT / path for path in control["feature"]["handoffs"]),
        ]
        if state.get("architecture_approval_token") != "APPROVED_ARCHITECTURE" or state.get("architecture_approved_sha256") != combined_digest(architecture_paths):
            raise SystemExit("FAIL: current architecture package lacks explicit human approval")
    elif args.stage == "design":
        requirements = spec / "requirements.md"
        if state.get("requirements_approval_token") != "APPROVED_FOR_DESIGN" or not requirements.is_file() or state.get("requirements_approved_sha256") != sha256(requirements):
            raise SystemExit("FAIL: current requirements are not machine-approved for design")
    else:
        requirements, design = spec / "requirements.md", spec / "design.md"
        if state.get("executable_plan_approval_token") != "APPROVED_EXECUTABLE_PLAN" or state.get("executable_plan_approved_sha256") != combined_digest([requirements, design]):
            raise SystemExit("FAIL: current requirements/design package lacks explicit human executable-plan approval")
    writable = set(control["guardrails"]["kiro"]["writable_by_stage"].get(args.stage, []))
    if len(writable) != 1:
        raise SystemExit(f"FAIL: {args.stage} must have exactly one writable output")
    # Earlier approved spec files may remain uncommitted until the complete
    # requirements/design/tasks package is committed. Their approval digests
    # above pin their content; they are not additional Kiro-writable outputs.
    allowed = set(writable)
    if args.stage in {"design", "tasks"}:
        allowed.add(str((spec / "requirements.md").relative_to(ROOT)))
    if args.stage == "tasks":
        allowed.add(str((spec / "design.md").relative_to(ROOT)))
    allowed.add(str(state_path.relative_to(ROOT)))
    if args.mode == "pre":
        dirty = changed_paths()
        unexpected = [path for path in dirty if path not in allowed]
        if unexpected:
            raise SystemExit("FAIL: working tree has out-of-stage changes: " + ", ".join(unexpected))
        print(f"PASS: {args.feature} {args.stage} pre-stage boundary")
        return
    context_path = ROOT / state["generated_prompt_path"] / f"{args.stage}.context.json"
    if not context_path.is_file():
        raise SystemExit(f"FAIL: context manifest missing: {context_path.relative_to(ROOT)}")
    context = json.loads(context_path.read_text())
    if context.get("feature") != args.feature or context.get("stage") != args.stage:
        raise SystemExit("FAIL: context manifest identity mismatch")
    for item in context.get("files", []):
        path = ROOT / item["path"]
        if not path.is_file() or sha256(path) != item["sha256"]:
            raise SystemExit(f"FAIL: context changed or missing: {item['path']}")
    if args.mode == "post":
        target = ROOT / next(iter(writable))
        if not target.is_file():
            raise SystemExit(f"FAIL: expected stage output missing: {target.relative_to(ROOT)}")
        unexpected = [path for path in changed_paths() if path not in allowed]
        if unexpected:
            raise SystemExit("FAIL: Kiro changed paths outside stage allowlist: " + ", ".join(unexpected))
    print(f"PASS: {args.feature} {args.stage} {args.mode} boundary")


if __name__ == "__main__":
    main()
