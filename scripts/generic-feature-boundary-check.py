#!/usr/bin/env python3
"""Fail-closed boundary validation for manifest-controlled feature tasks."""

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


def changed_paths() -> list[str]:
    result = subprocess.run(
        ["git", "status", "--porcelain=v1", "--untracked-files=all"],
        cwd=ROOT,
        text=True,
        capture_output=True,
        check=True,
    )
    return [line[3:] for line in result.stdout.splitlines() if len(line) > 3]


def allowed(path: str, rules: list[str], feature: str) -> bool:
    if path == f".automation/state/{feature}.json":
        return True
    return any(path == rule or (rule.endswith("/") and path.startswith(rule)) for rule in rules)


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    parser.add_argument("--task", required=True, type=int)
    parser.add_argument("--mode", choices=("context", "changed", "all"), default="all")
    args = parser.parse_args()
    subprocess.run(
        [str(ROOT / "scripts/feature-control.py"), "validate", "--feature", args.feature],
        cwd=ROOT,
        check=True,
    )
    state = json.loads((ROOT / f".automation/state/{args.feature}.json").read_text())
    control_path = ROOT / f".automation/features/{args.feature}.control.json"
    if state.get("control_approved_sha256") != sha256(control_path):
        raise SystemExit("FAIL: feature control differs from the founder-approved digest")
    if state.get("current_stage") != "cursor" or state.get("tasks_approval_token") != "APPROVED_FOR_CURSOR":
        raise SystemExit("FAIL: feature is not approved for Cursor")
    context_path = ROOT / state["generated_prompt_path"] / f"cursor-task-{args.task}.context.json"
    if not context_path.is_file():
        raise SystemExit(f"FAIL: task context manifest missing: {context_path.relative_to(ROOT)}")
    context = json.loads(context_path.read_text())
    if context.get("feature") != args.feature or context.get("task") != args.task:
        raise SystemExit("FAIL: task context identity mismatch")
    tasks = ROOT / state["spec_path"] / "tasks.md"
    if context.get("tasks_approved_sha256") != state.get("tasks_approved_sha256") or sha256(tasks) != state.get("tasks_approved_sha256"):
        raise SystemExit("FAIL: approved tasks digest mismatch")
    if args.mode in {"context", "all"}:
        for item in context.get("files", []):
            path = ROOT / item["path"]
            if not path.is_file() or sha256(path) != item["sha256"]:
                raise SystemExit(f"FAIL: context changed or missing: {item['path']}")
        print(f"PASS: {args.feature} Task {args.task} approved context hashes")
    if args.mode in {"changed", "all"}:
        writable = context.get("writable_paths") or []
        outside = [path for path in changed_paths() if not allowed(path, writable, args.feature)]
        if outside:
            raise SystemExit("FAIL: changed paths outside task allowlist: " + ", ".join(outside))
        control = json.loads((ROOT / f".automation/features/{args.feature}.control.json").read_text())
        forbidden = control["ownership"]["forbidden_concepts"]
        production = [
            path for path in changed_paths()
            if path.startswith(("internal/", "cmd/", "api/schemas/"))
            and not path.endswith("_test.go")
            and (ROOT / path).is_file()
        ]
        violations = []
        for relative in production:
            text = (ROOT / relative).read_text(errors="replace")
            for concept in forbidden:
                if re.search(rf"(?i)(?<![A-Za-z0-9_]){re.escape(concept)}(?![A-Za-z0-9_])", text):
                    violations.append(f"{relative}: {concept}")
        if violations:
            raise SystemExit("FAIL: forbidden adjacent-feature concepts in production changes: " + "; ".join(violations))
        print(f"PASS: {args.feature} Task {args.task} changed paths and forbidden-concept boundary")


if __name__ == "__main__":
    main()
