#!/usr/bin/env python3
"""Validate one FEATURE-0014 implementation task against its approved boundary."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
TASKS = ROOT / ".kiro/specs/provider-neutral-resource-model/tasks.md"
STATE = ROOT / ".automation/state/FEATURE-0014.json"


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--task", required=True, type=int)
    args = parser.parse_args()
    manifest_path = ROOT / f"docs/generated-prompts/FEATURE-0014/cursor-task-{args.task}.context.json"
    if not manifest_path.is_file():
        raise SystemExit(f"FAIL: missing task context manifest: {manifest_path}")
    state = json.loads(STATE.read_text())
    manifest = json.loads(manifest_path.read_text())
    errors: list[str] = []
    if state.get("tasks_approval_token") != "APPROVED_FOR_CURSOR":
        errors.append("missing APPROVED_FOR_CURSOR")
    if state.get("tasks_approved_sha256") != sha256(TASKS):
        errors.append("approved tasks digest differs from tasks.md")
    if manifest.get("tasks_sha256") != sha256(TASKS) or manifest.get("task") != str(args.task):
        errors.append("task manifest identity/digest mismatch")
    for item in manifest.get("context", []):
        candidate = ROOT / item["path"]
        if not candidate.is_file() or sha256(candidate) != item.get("sha256"):
            errors.append(f"context changed or missing: {item['path']}")

    rules = manifest.get("writable_paths", [])
    status = subprocess.check_output(["git", "status", "--porcelain"], cwd=ROOT, text=True)
    changed = [line[3:] for line in status.splitlines() if len(line) > 3]
    allowed_always = {".automation/state/FEATURE-0014.json"}
    outside = [p for p in changed if p not in allowed_always and not any(
        p == rule or (rule.endswith("/") and p.startswith(rule)) for rule in rules)]
    if outside:
        errors.append(f"out-of-task changed paths: {outside}")

    forbidden = re.compile(r"\b(?:IaaSStack|ProviderRegion|ResourcePool|ProviderCapability|DecisionRecord|AuditEvent)\b")
    mechanism = re.compile(r"(?i)\b(?:type|var)\s+\w*(?:Repository|Persistence|StateProvider|Resolver|Adapter|Controller|Reconciler)\b")
    for rel in changed:
        path = ROOT / rel
        if not path.is_file() or path.suffix not in {".go", ".json"} or rel.endswith("_test.go") or "/fixtures/" in rel:
            continue
        text = path.read_text(errors="replace")
        if forbidden.search(text):
            errors.append(f"adjacent/superseded concept in production artifact: {rel}")
        if mechanism.search(text):
            errors.append(f"excluded production mechanism in: {rel}")

    if errors:
        for error in errors:
            print(f"FAIL: {error}")
        raise SystemExit(1)
    print(f"PASS: FEATURE-0014 Cursor Task {args.task} approved digest and context hashes")
    print("PASS: changed paths are confined to the task Files allowlist")
    print("PASS: production artifacts exclude adjacent kinds and runtime mechanisms")


if __name__ == "__main__":
    main()
