#!/usr/bin/env python3
"""Render one fail-closed FEATURE-0014 Cursor task prompt and hashed context."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SPEC = ROOT / ".kiro/specs/provider-neutral-resource-model"
TASKS = SPEC / "tasks.md"
DESIGN = SPEC / "design.md"
STATE = ROOT / ".automation/state/FEATURE-0014.json"
TEMPLATE = ROOT / "docs/prompts/cursor/feature-0014-task.prompt.md"

CORE = [
    "AGENTS.md",
    "docs/engineering/go-coding-guardrails.md",
    "docs/engineering/go-version-standard.md",
    "docs/architecture/api-resource-standard.md",
    ".kiro/specs/provider-neutral-resource-model/design.md",
]

VALIDATOR_REUSE_CONTEXT = [
    "internal/validation/providerlocation.go",
    "internal/validation/providerlocation_test.go",
    "internal/apivalid/structural.go",
    "internal/apivalid/stage_semantic.go",
    "internal/apivalid/limits.go",
    "internal/apivalid/translate.go",
    "internal/apiref/constraints.go",
    "internal/apimeta/scope.go",
]

CONTEXT_BY_TASK = {
    **{n: ["internal/resources/organization.go", "internal/apimeta/objectmeta.go", "internal/apimeta/typemeta.go", "internal/apimeta/scope.go", "internal/apiref/reference.go"] for n in range(1, 6)},
    **{n: ["api/schemas/project.json", "api/schemas/_common/object-meta.json", "api/schemas/_common/scope-ref.json", "api/schemas/_common/typed-ref.json", "api/schemas/_common/condition.json"] for n in range(6, 11)},
    12: ["internal/resources/providerlocation.go", "internal/validation/organization.go", "internal/apivalid/structural.go", "internal/apivalid/limits.go"],
    13: ["internal/resources/provider.go", "api/schemas/provider.json", *VALIDATOR_REUSE_CONTEXT],
    14: ["internal/resources/providerdatacenter.go", "api/schemas/provider-datacenter.json", *VALIDATOR_REUSE_CONTEXT],
    15: ["internal/resources/datacenterfailuredomain.go", "api/schemas/datacenter-failure-domain.json", *VALIDATOR_REUSE_CONTEXT],
    16: ["internal/resources/infrastructurestack.go", "api/schemas/infrastructure-stack.json", *VALIDATOR_REUSE_CONTEXT],
    17: ["internal/resources/provider.go", "internal/resources/providerlocation.go", "internal/resources/providerdatacenter.go", "internal/resources/datacenterfailuredomain.go", "internal/resources/infrastructurestack.go", "internal/apiref/constraints.go", "internal/apimeta/scope.go"],
    18: ["internal/resources/provider.go", "internal/resources/providerlocation.go", "internal/resources/providerdatacenter.go", "internal/resources/datacenterfailuredomain.go", "internal/resources/infrastructurestack.go", "internal/validation/topology.go", "internal/apicond/condition.go", "internal/apimeta/uid.go"],
    19: ["internal/resources/provider.go", "internal/resources/providerlocation.go", "internal/resources/providerdatacenter.go", "internal/resources/datacenterfailuredomain.go", "internal/resources/infrastructurestack.go", "api/schemas/provider.json", "api/schemas/provider-location.json", "api/schemas/provider-datacenter.json", "api/schemas/datacenter-failure-domain.json", "api/schemas/infrastructure-stack.json", "internal/apiconform/bindings.go", "internal/apiconform/bindings_test.go", "internal/apiconform/schemaregistry.go", "internal/apiconform/imports_test.go"],
    20: ["internal/resources/provider.go", "internal/resources/providerlocation.go", "internal/resources/providerdatacenter.go", "internal/resources/datacenterfailuredomain.go", "internal/resources/infrastructurestack.go", "internal/validation/topology.go", "internal/validation/completeness.go", "internal/apiconform/fixtures.go", "internal/apiconform/fixtures_valid_test.go", "tests/conformance/fixtures/project.json"],
    21: ["internal/validation/topology.go", "internal/validation/completeness.go", "internal/apiconform/feature0014_bindings.go", "internal/apiconform/fixtures.go", "internal/apiconform/fixtures_negative_test.go", "tests/conformance/fixtures/negative/invalid-scope-kind.json"],
    22: ["Makefile", "scripts/verify.sh", "scripts/guardrails.sh", "scripts/feature-0014-boundary-check.py"],
}


def digest(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def task_block(task_id: int) -> str:
    text = TASKS.read_text()
    match = re.search(rf"(?ms)^## Task {task_id}\b.*?(?=^## Task \d+\b|^## Requirement-to-task ledger)", text)
    if not match:
        raise SystemExit(f"task {task_id} not found")
    return match.group(0).rstrip()


def writable_paths(block: str) -> list[str]:
    line = next((line for line in block.splitlines() if line.startswith("Files:")), "")
    paths = re.findall(r"`([^`]+)`", line)
    if not paths and "Files: none" not in line:
        raise SystemExit("task Files label has no explicit path")
    return paths


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--task", required=True, type=int)
    parser.add_argument("--output", required=True)
    parser.add_argument("--manifest", required=True)
    args = parser.parse_args()

    state = json.loads(STATE.read_text())
    if state.get("tasks_approval_token") != "APPROVED_FOR_CURSOR":
        raise SystemExit("missing APPROVED_FOR_CURSOR")
    if state.get("tasks_approved_sha256") != digest(TASKS):
        raise SystemExit("approved tasks digest does not match tasks.md")
    if state.get("current_stage") != "cursor":
        raise SystemExit("FEATURE-0014 is not at cursor stage")

    block = task_block(args.task)
    writable = writable_paths(block)
    relpaths = list(dict.fromkeys(CORE + CONTEXT_BY_TASK.get(args.task, [])))
    missing = [p for p in relpaths if not (ROOT / p).is_file()]
    if missing:
        raise SystemExit(f"context files missing: {missing}")

    manifest_path = Path(args.manifest)
    manifest_path.parent.mkdir(parents=True, exist_ok=True)
    manifest = {
        "feature": "FEATURE-0014",
        "task": str(args.task),
        "tasks_sha256": digest(TASKS),
        "writable_paths": writable,
        "context": [{"path": p, "sha256": digest(ROOT / p)} for p in relpaths],
    }
    manifest_path.write_text(json.dumps(manifest, indent=2) + "\n")

    chunks = []
    for rel in relpaths:
        chunks.append(f"### `{rel}`\n\n```text\n{(ROOT / rel).read_text()}\n```")
    writable_md = "\n".join(f"- `{p}`" for p in writable) if writable else "- No source files (verification-only task)."
    try:
        manifest_label = str(manifest_path.resolve().relative_to(ROOT))
    except ValueError:
        manifest_label = str(manifest_path)
    rendered = (TEMPLATE.read_text()
        .replace("{{TASK_ID}}", str(args.task))
        .replace("{{WRITABLE_PATHS}}", writable_md)
        .replace("{{TASK_BLOCK}}", block)
        .replace("{{CONTEXT_MANIFEST}}", manifest_label)
        .replace("{{CONTEXT_CONTENT}}", "\n\n".join(chunks)))
    Path(args.output).write_text(rendered)
    print(args.output)


if __name__ == "__main__":
    main()
