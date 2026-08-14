#!/usr/bin/env python3
"""Validate and resolve Sovrunn generic feature-control manifests."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import shlex
import sys
from pathlib import Path
from typing import Any

ROOT = Path(__file__).resolve().parents[1]
CONTROL_DIR = ROOT / ".automation/features"
SCHEMA_PATH = ROOT / ".automation/schemas/feature-control.schema.json"
STAGES = ("architecture", "requirements", "design", "tasks", "implementation", "review")
AI_STAGES = {"requirements", "design", "tasks", "implementation", "review"}
ALLOWED_VERIFY_TARGETS = {"fmt", "vet", "test", "test-race", "build", "ff-verify", "ff-guardrails"}
REQUIRED_GO_CONTEXT = {
    "docs/engineering/go-coding-guardrails.md",
    "docs/engineering/go-version-standard.md",
    "docs/engineering/go-observability-standard.md",
}
STOP_CONDITIONS = {
    "ARCHITECTURE_DECISION_REQUIRED",
    "REQUIREMENT_CLARIFICATION_REQUIRED",
    "BOUNDARY_CHANGE_REQUIRED",
    "DEPENDENCY_APPROVAL_REQUIRED",
    "SECURITY_REVIEW_REQUIRED",
}
REQUIRED_TOP_LEVEL = {
    "schema_version",
    "feature",
    "ownership",
    "context",
    "guardrails",
    "execution",
    "models",
    "closeout",
}


class ControlError(ValueError):
    pass


def fail(message: str) -> None:
    raise SystemExit(f"ERROR: {message}")


def control_path(feature: str) -> Path:
    return CONTROL_DIR / f"{feature}.control.json"


def repo_path(raw: str, *, field: str, must_exist: bool = True) -> Path:
    if not isinstance(raw, str) or not raw:
        raise ControlError(f"{field} must be a non-empty repository-relative path")
    path = Path(raw)
    if path.is_absolute() or ".." in path.parts or path.parts[0] == ".git":
        raise ControlError(f"{field} is not a safe repository-relative path: {raw}")
    resolved = (ROOT / path).resolve()
    try:
        resolved.relative_to(ROOT.resolve())
    except ValueError as exc:
        raise ControlError(f"{field} escapes the repository: {raw}") from exc
    if must_exist and not resolved.exists():
        raise ControlError(f"{field} does not exist: {raw}")
    return resolved


def require_keys(obj: dict[str, Any], required: set[str], *, field: str) -> None:
    missing = sorted(required - set(obj))
    extra = sorted(set(obj) - required)
    if missing or extra:
        details = []
        if missing:
            details.append("missing=" + ",".join(missing))
        if extra:
            details.append("extra=" + ",".join(extra))
        raise ControlError(f"{field} has invalid keys: {' '.join(details)}")


def require_string_list(value: Any, *, field: str, allow_empty: bool = True) -> list[str]:
    if not isinstance(value, list) or any(not isinstance(item, str) or not item for item in value):
        raise ControlError(f"{field} must be a list of non-empty strings")
    if not allow_empty and not value:
        raise ControlError(f"{field} must not be empty")
    if len(value) != len(set(value)):
        raise ControlError(f"{field} contains duplicates")
    return value


def feature_number(feature_id: str) -> int:
    match = re.fullmatch(r"FEATURE-(\d{4})", feature_id)
    if not match:
        raise ControlError(f"invalid feature id: {feature_id}")
    return int(match.group(1))


def load(path: Path) -> dict[str, Any]:
    try:
        data = json.loads(path.read_text())
    except FileNotFoundError as exc:
        raise ControlError(f"control manifest not found: {path.relative_to(ROOT)}") from exc
    except json.JSONDecodeError as exc:
        raise ControlError(f"invalid JSON in {path.relative_to(ROOT)}: {exc}") from exc
    if not isinstance(data, dict):
        raise ControlError("control manifest must be a JSON object")
    return data


def validate(data: dict[str, Any], *, expected_feature: str = "") -> None:
    if set(data) != REQUIRED_TOP_LEVEL:
        missing = sorted(REQUIRED_TOP_LEVEL - set(data))
        extra = sorted(set(data) - REQUIRED_TOP_LEVEL)
        detail = []
        if missing:
            detail.append("missing=" + ",".join(missing))
        if extra:
            detail.append("extra=" + ",".join(extra))
        raise ControlError("invalid top-level keys: " + " ".join(detail))
    if data["schema_version"] != "1.0":
        raise ControlError("schema_version must be 1.0")
    if not SCHEMA_PATH.is_file():
        raise ControlError("feature-control schema is missing")

    feature = data["feature"]
    if not isinstance(feature, dict):
        raise ControlError("feature must be an object")
    require_keys(
        feature,
        {"id", "slug", "title", "tier", "phase", "phase_branch", "feature_branch", "feature_file", "architecture", "handoffs", "spec_path"},
        field="feature",
    )
    current_number = feature_number(feature["id"])
    if expected_feature and feature["id"] != expected_feature:
        raise ControlError(f"manifest feature id {feature['id']} does not match {expected_feature}")
    if feature["tier"] not in {"A", "B", "C"}:
        raise ControlError("feature.tier must be A, B, or C")
    if not isinstance(feature["phase"], int) or feature["phase"] < 1:
        raise ControlError("feature.phase must be a positive integer")
    if not re.fullmatch(r"[a-z0-9]+(?:-[a-z0-9]+)*", feature["slug"]):
        raise ControlError("feature.slug must be lowercase kebab-case")
    for key in ("feature_file", "architecture"):
        repo_path(feature[key], field=f"feature.{key}")
    require_string_list(feature["handoffs"], field="feature.handoffs", allow_empty=False)
    for index, path in enumerate(feature["handoffs"]):
        repo_path(path, field=f"feature.handoffs[{index}]")
    repo_path(feature["spec_path"], field="feature.spec_path", must_exist=False)

    ownership = data["ownership"]
    require_keys(ownership, {"owned_resources", "previous_features", "excluded_features", "forbidden_concepts"}, field="ownership")
    require_string_list(ownership["owned_resources"], field="ownership.owned_resources", allow_empty=False)
    require_string_list(ownership["forbidden_concepts"], field="ownership.forbidden_concepts")
    if not isinstance(ownership["previous_features"], list):
        raise ControlError("ownership.previous_features must be an array")
    if feature["phase"] == 2 and current_number > 11 and not ownership["previous_features"]:
        raise ControlError("Phase 2 features after FEATURE-0011 require explicit previous-feature context")
    previous_ids: set[str] = set()
    for index, dependency in enumerate(ownership["previous_features"]):
        if not isinstance(dependency, dict):
            raise ControlError(f"ownership.previous_features[{index}] must be an object")
        require_keys(dependency, {"id", "relationship", "context"}, field=f"ownership.previous_features[{index}]")
        dep_number = feature_number(dependency["id"])
        if dep_number >= current_number:
            raise ControlError(f"previous feature {dependency['id']} must precede {feature['id']}")
        if dependency["id"] in previous_ids:
            raise ControlError(f"duplicate previous feature: {dependency['id']}")
        previous_ids.add(dependency["id"])
        if not isinstance(dependency["relationship"], str) or not dependency["relationship"].strip():
            raise ControlError(f"previous feature {dependency['id']} requires a relationship")
        validate_stage_paths(dependency["context"], field=f"ownership.previous_features[{index}].context")
        if not any(dependency["context"].values()):
            raise ControlError(f"previous feature {dependency['id']} requires at least one context path")

    excluded_ids: set[str] = set()
    for index, excluded in enumerate(ownership["excluded_features"]):
        if not isinstance(excluded, dict):
            raise ControlError(f"ownership.excluded_features[{index}] must be an object")
        require_keys(excluded, {"id", "concepts"}, field=f"ownership.excluded_features[{index}]")
        feature_number(excluded["id"])
        if excluded["id"] in previous_ids:
            raise ControlError(f"feature cannot be both previous dependency and excluded: {excluded['id']}")
        if excluded["id"] in excluded_ids:
            raise ControlError(f"duplicate excluded feature: {excluded['id']}")
        excluded_ids.add(excluded["id"])
        require_string_list(excluded["concepts"], field=f"ownership.excluded_features[{index}].concepts", allow_empty=False)

    context = data["context"]
    require_keys(context, {"always", "stages", "budgets"}, field="context")
    for index, path in enumerate(require_string_list(context["always"], field="context.always")):
        repo_path(path, field=f"context.always[{index}]")
    validate_stage_paths(context["stages"], field="context.stages")
    if not isinstance(context["budgets"], dict):
        raise ControlError("context.budgets must be an object")
    if set(context["budgets"]) != AI_STAGES:
        raise ControlError("context.budgets must define exactly requirements, design, tasks, implementation, and review")
    for stage, budget in context["budgets"].items():
        if stage not in STAGES or not isinstance(budget, int) or budget < 1024:
            raise ControlError(f"invalid context budget for {stage}")

    guardrails = data["guardrails"]
    require_keys(guardrails, {"kiro", "cursor", "stop_conditions"}, field="guardrails")
    require_keys(guardrails["kiro"], {"writable_by_stage", "forbidden_outputs"}, field="guardrails.kiro")
    validate_stage_paths(guardrails["kiro"]["writable_by_stage"], field="guardrails.kiro.writable_by_stage", must_exist=False)
    writable = guardrails["kiro"]["writable_by_stage"]
    if set(writable) != {"requirements", "design", "tasks"} or any(len(writable[stage]) != 1 for stage in writable):
        raise ControlError("Kiro writable paths must define exactly one output for requirements, design, and tasks")
    for index, path in enumerate(require_string_list(guardrails["kiro"]["forbidden_outputs"], field="guardrails.kiro.forbidden_outputs")):
        repo_path(path, field=f"guardrails.kiro.forbidden_outputs[{index}]", must_exist=False)
    require_keys(guardrails["cursor"], {"go_context", "require_complete_receipt", "verify_commands"}, field="guardrails.cursor")
    go_context = require_string_list(guardrails["cursor"]["go_context"], field="guardrails.cursor.go_context", allow_empty=False)
    if not REQUIRED_GO_CONTEXT.issubset(go_context):
        raise ControlError("Cursor context must include the canonical Go coding, version, and observability standards")
    for index, path in enumerate(go_context):
        repo_path(path, field=f"guardrails.cursor.go_context[{index}]")
    if guardrails["cursor"]["require_complete_receipt"] is not True:
        raise ControlError("Cursor COMPLETE receipt must be required")
    verify_commands = require_string_list(guardrails["cursor"]["verify_commands"], field="guardrails.cursor.verify_commands", allow_empty=False)
    for raw in verify_commands:
        parts = shlex.split(raw)
        if len(parts) < 2 or parts[0] != "make" or parts[1] not in ALLOWED_VERIFY_TARGETS:
            raise ControlError(f"unsafe Cursor verification command: {raw}")
        if any(part != f"FEATURE={feature['id']}" for part in parts[2:]):
            raise ControlError(f"unsafe Cursor verification argument: {raw}")
    stops = set(require_string_list(guardrails["stop_conditions"], field="guardrails.stop_conditions", allow_empty=False))
    if stops != STOP_CONDITIONS:
        raise ControlError("guardrails.stop_conditions must contain the complete approved stop-condition set")

    execution = data["execution"]
    require_keys(execution, {"human_gates", "tasks_per_run", "require_clean_tree", "commit_per_task", "max_stage_revisions"}, field="execution")
    if execution["human_gates"] != ["architecture", "executable_plan", "final"]:
        raise ControlError("human gates must be architecture, executable_plan, final")
    if execution["require_clean_tree"] is not True or execution["commit_per_task"] is not True:
        raise ControlError("clean tree and per-task commits are mandatory")
    if not isinstance(execution["tasks_per_run"], int) or not 1 <= execution["tasks_per_run"] <= 4:
        raise ControlError("execution.tasks_per_run must be between 1 and 4")
    if not isinstance(execution["max_stage_revisions"], int) or not 1 <= execution["max_stage_revisions"] <= 3:
        raise ControlError("execution.max_stage_revisions must be between 1 and 3")

    models = data["models"]
    expected_models = {
        "architecture": "architecture",
        "requirements": "requirements",
        "design": "design",
        "tasks": "tasks",
        "review": "review",
        "implementation": "implementation",
        "mechanical": "deterministic_first",
    }
    if models != expected_models:
        raise ControlError("models must use the approved centralized role names")

    closeout = data["closeout"]
    require_keys(closeout, {"mode", "metadata_targets", "commit", "push"}, field="closeout")
    if closeout["mode"] != "prepare_only" or closeout["commit"] is not False or closeout["push"] is not False:
        raise ControlError("closeout must be prepare_only and must not commit or push")
    for index, path in enumerate(require_string_list(closeout["metadata_targets"], field="closeout.metadata_targets", allow_empty=False)):
        repo_path(path, field=f"closeout.metadata_targets[{index}]")


def validate_stage_paths(value: Any, *, field: str, must_exist: bool = True) -> None:
    if not isinstance(value, dict):
        raise ControlError(f"{field} must be an object")
    unknown = sorted(set(value) - set(STAGES))
    if unknown:
        raise ControlError(f"{field} has unknown stages: {', '.join(unknown)}")
    for stage, paths in value.items():
        for index, path in enumerate(require_string_list(paths, field=f"{field}.{stage}")):
            repo_path(path, field=f"{field}.{stage}[{index}]", must_exist=must_exist)


def resolve_context(data: dict[str, Any], stage: str, *, review_stage: str | None = None) -> list[Path]:
    if stage not in STAGES:
        raise ControlError(f"unsupported stage: {stage}")
    if stage == "review" and review_stage not in {"requirements", "design", "tasks"}:
        raise ControlError("review context requires review_stage=requirements, design, or tasks")
    if stage != "review" and review_stage is not None:
        raise ControlError("review_stage is valid only for review context")
    feature = data["feature"]
    raw_paths: list[str] = list(data["context"]["always"])
    raw_paths.extend([feature["feature_file"], feature["architecture"], *feature["handoffs"]])
    for dependency in data["ownership"]["previous_features"]:
        raw_paths.extend(dependency["context"].get(stage, []))
    raw_paths.extend(data["context"]["stages"].get(stage, []))
    if stage == "implementation":
        raw_paths.extend(data["guardrails"]["cursor"]["go_context"])
    spec = feature["spec_path"]
    spec_stage = review_stage if stage == "review" else stage
    if spec_stage in {"design", "tasks", "implementation"}:
        raw_paths.append(f"{spec}/requirements.md")
    if spec_stage in {"tasks", "implementation"}:
        raw_paths.append(f"{spec}/design.md")
    if spec_stage == "implementation":
        raw_paths.append(f"{spec}/tasks.md")
    seen: set[str] = set()
    paths: list[Path] = []
    for index, raw in enumerate(raw_paths):
        if raw in seen:
            continue
        seen.add(raw)
        paths.append(repo_path(raw, field=f"resolved_context[{index}]"))
    return paths


def context_manifest(
    data: dict[str, Any], stage: str, *, review_stage: str | None = None
) -> dict[str, Any]:
    paths = resolve_context(data, stage, review_stage=review_stage)
    files = []
    total_bytes = 0
    for path in paths:
        payload = path.read_bytes()
        relative = str(path.relative_to(ROOT))
        total_bytes += len(payload)
        files.append(
            {
                "path": relative,
                "bytes": len(payload),
                "lines": len(payload.splitlines()),
                "estimated_tokens": (len(payload) + 3) // 4,
                "sha256": hashlib.sha256(payload).hexdigest(),
            }
        )
    budget = data["context"]["budgets"].get(stage)
    if budget is None:
        raise ControlError(f"no context budget configured for stage {stage}")
    if total_bytes > budget:
        largest = sorted(files, key=lambda item: item["bytes"], reverse=True)[:5]
        detail = ", ".join(f"{item['path']}={item['bytes']}" for item in largest)
        raise ControlError(f"{stage} context exceeds budget {total_bytes}>{budget}; largest: {detail}")
    return {
        "schema_version": "1.0",
        "feature": data["feature"]["id"],
        "stage": stage,
        "budget_bytes": budget,
        "total_bytes": total_bytes,
        "estimated_tokens": (total_bytes + 3) // 4,
        "files": files,
    }


def prompt_fragment(data: dict[str, Any], manifest: dict[str, Any], stage: str) -> str:
    ownership = data["ownership"]
    lines = [
        "## Automation-controlled context and boundaries",
        "",
        "Read only the exact semantic context files listed below. Do not search for or load adjacent feature specifications unless this manifest lists them.",
        "",
    ]
    lines.extend(f"- `{item['path']}`" for item in manifest["files"])
    lines.extend(
        [
            "",
            f"Context budget: {manifest['total_bytes']} of {manifest['budget_bytes']} bytes (approximately {manifest['estimated_tokens']} tokens).",
            "",
            "Owned resources:",
        ]
    )
    lines.extend(f"- `{resource}`" for resource in ownership["owned_resources"])
    lines.extend(["", "Previous-feature dependencies are context only; their contracts remain owned by those features:"])
    lines.extend(f"- `{item['id']}`: {item['relationship']}" for item in ownership["previous_features"])
    lines.extend(["", "Excluded adjacent-feature concepts:"])
    for item in ownership["excluded_features"]:
        lines.append(f"- `{item['id']}`: {', '.join(item['concepts'])}")
    lines.extend(["", "Fail closed with exactly one approved stop condition when the stage cannot proceed:"])
    lines.extend(f"- `{condition}`" for condition in data["guardrails"]["stop_conditions"])
    if stage in {"requirements", "design", "tasks"}:
        writable = data["guardrails"]["kiro"]["writable_by_stage"].get(stage, [])
        lines.extend(["", f"During the `{stage}` stage, modify only:"])
        lines.extend(f"- `{path}`" for path in writable)
        lines.extend(["", "Do not generate or revise later-stage artifacts, source code, schemas, automation, or architecture."])
    return "\n".join(lines) + "\n"


def validate_state_identity(data: dict[str, Any]) -> None:
    feature_id = data["feature"]["id"]
    path = ROOT / f".automation/state/{feature_id}.json"
    if not path.is_file():
        raise ControlError(f"feature state does not exist: {path.relative_to(ROOT)}")
    state = json.loads(path.read_text())
    expected = {
        "feature_id": data["feature"]["id"],
        "slug": data["feature"]["slug"],
        "title": data["feature"]["title"],
        "phase_branch": data["feature"]["phase_branch"],
        "feature_branch": data["feature"]["feature_branch"],
        "spec_path": data["feature"]["spec_path"],
    }
    mismatches = [f"{key}={state.get(key)!r}, expected {value!r}" for key, value in expected.items() if state.get(key) != value]
    if mismatches:
        raise ControlError("feature state/control identity mismatch: " + "; ".join(mismatches))


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("command", choices=("validate", "context", "prompt-fragment", "path"))
    parser.add_argument("--feature", required=True)
    parser.add_argument("--manifest", help="explicit manifest path for tests or migration validation")
    parser.add_argument("--stage", choices=STAGES)
    parser.add_argument("--review-stage", choices=("requirements", "design", "tasks"))
    parser.add_argument("--output")
    args = parser.parse_args()
    path = repo_path(args.manifest, field="manifest") if args.manifest else control_path(args.feature)
    try:
        data = load(path)
        validate(data, expected_feature=args.feature)
        if not args.manifest:
            validate_state_identity(data)
        if args.command == "validate":
            print(f"PASS: {path.relative_to(ROOT)}")
            return
        if args.command == "path":
            print(path.relative_to(ROOT))
            return
        if not args.stage:
            fail("--stage is required")
        manifest = context_manifest(data, args.stage, review_stage=args.review_stage)
        if args.command == "context":
            rendered = json.dumps(manifest, indent=2, sort_keys=True) + "\n"
        else:
            rendered = prompt_fragment(data, manifest, args.stage)
        if args.output:
            output = repo_path(args.output, field="output", must_exist=False)
            output.parent.mkdir(parents=True, exist_ok=True)
            output.write_text(rendered)
            print(output.relative_to(ROOT))
        else:
            print(rendered, end="")
    except ControlError as exc:
        fail(str(exc))


if __name__ == "__main__":
    main()
