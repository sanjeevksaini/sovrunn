#!/usr/bin/env python3
"""Fail-closed semantic checks for manifest-controlled Kiro spec stages."""

from __future__ import annotations

import argparse
import json
import re
import subprocess
import sys
from pathlib import Path
from typing import Any

try:
    import yaml
except ImportError:  # pragma: no cover - repository validation installs PyYAML
    yaml = None


ROOT = Path(__file__).resolve().parents[1]
STAGES = ("requirements", "design", "tasks")
RECEIPT = "STAGE_STATUS: COMPLETE"
REQ_ID = re.compile(r"REQ-[A-Z0-9]+-\d+")
AC_ID = re.compile(r"AC-[A-Z0-9]+-\d+")
CF_ID = re.compile(r"VS0-CF-[A-Z]+\d+")
VS_ID = re.compile(r"VS0-(?:SCHEMA|WRITER|STATE|CF)-[A-Z0-9-]+")


def table_cells(line: str) -> tuple[str, ...] | None:
    stripped = line.strip()
    if not (stripped.startswith("|") and stripped.endswith("|")):
        return None
    return tuple(cell.strip() for cell in stripped[1:-1].split("|"))


def rows_by_id(text: str, pattern: re.Pattern[str]) -> dict[str, list[tuple[str, ...]]]:
    rows: dict[str, list[tuple[str, ...]]] = {}
    for line in text.splitlines():
        cells = table_cells(line)
        if not cells or not pattern.fullmatch(cells[0]):
            continue
        rows.setdefault(cells[0], []).append(cells)
    return rows


def section(text: str, title: str) -> str | None:
    lines = text.splitlines()
    start = None
    level = 0
    wanted = title.casefold()
    for index, line in enumerate(lines):
        heading = re.match(r"^(#{1,6})\s+(.+?)\s*$", line)
        heading_title = heading.group(2).strip() if heading else ""
        heading_title = re.sub(r"^\d+(?:\.\d+)*[.)]?\s+", "", heading_title)
        if heading and heading_title.casefold() == wanted:
            start = index + 1
            level = len(heading.group(1))
            break
    if start is None:
        return None
    end = len(lines)
    for index in range(start, len(lines)):
        heading = re.match(r"^(#{1,6})\s+", lines[index])
        if heading and len(heading.group(1)) <= level:
            end = index
            break
    return "\n".join(lines[start:end])


def scalar(value: Any) -> str:
    if value is None:
        return "null"
    if isinstance(value, bool):
        return "true" if value else "false"
    return str(value)


def expand_cf_ranges(text: str) -> set[str]:
    found = set(CF_ID.findall(text))
    range_pattern = re.compile(r"VS0-CF-([A-Z]+)(\d+)\.\.([A-Z]*)(\d+)")
    for match in range_pattern.finditer(text):
        left_prefix, left_number, right_prefix, right_number = match.groups()
        prefix = right_prefix or left_prefix
        if prefix != left_prefix:
            continue
        width = max(len(left_number), len(right_number))
        for number in range(int(left_number), int(right_number) + 1):
            found.add(f"VS0-CF-{prefix}{number:0{width}d}")
    return found


def registry_conformance(root: Path) -> dict[str, dict[str, Any]]:
    if yaml is None:
        raise RuntimeError("PyYAML is required for Kiro semantic checks")
    registry_path = root / "docs/architecture/vertical-slices/VS-000-contract-registry.yaml"
    data = yaml.safe_load(registry_path.read_text())
    entries = data.get("conformance", [])
    return {entry["id"]: entry for entry in entries}


def expected_cf_row(entry: dict[str, Any]) -> tuple[str, ...]:
    return (
        scalar(entry.get("id")),
        scalar(entry.get("owner")),
        scalar(entry.get("inputs")),
        scalar(entry.get("expectedState")),
        scalar(entry.get("expectedError")),
        scalar(entry.get("expectedSideEffects")),
        scalar(entry.get("gate")),
    )


def require_exact_ledger(
    errors: list[str],
    source_rows: dict[str, list[tuple[str, ...]]],
    target_text: str,
    title: str,
) -> None:
    body = section(target_text, title)
    if body is None:
        errors.append(f"missing exact section heading: {title}")
        return
    target_rows = rows_by_id(body, REQ_ID if title == "Canonical requirement ledger" else AC_ID)
    for item_id, rows in source_rows.items():
        if len(rows) != 1:
            errors.append(f"approved feature authority defines {item_id} {len(rows)} times; expected once")
            continue
        if target_rows.get(item_id) != rows:
            errors.append(
                f"{title} must contain the approved {item_id} row exactly once with unchanged columns"
            )
    extra = sorted(set(target_rows) - set(source_rows))
    if extra:
        errors.append(f"{title} invents IDs not present in the feature authority: {', '.join(extra)}")


def check_requirements(
    errors: list[str],
    feature: str,
    feature_text: str,
    target_text: str,
    conformance: dict[str, dict[str, Any]],
) -> None:
    source_requirements = rows_by_id(feature_text, REQ_ID)
    source_acceptance = rows_by_id(feature_text, AC_ID)
    if not source_requirements or not source_acceptance:
        errors.append("feature authority must define tabular REQ and AC ledgers")
        return

    require_exact_ledger(
        errors, source_requirements, target_text, "Canonical requirement ledger"
    )
    require_exact_ledger(
        errors, source_acceptance, target_text, "Canonical acceptance ledger"
    )

    for item_id in source_requirements:
        count = len(
            re.findall(
                rf"^#{{2,6}}\s+{re.escape(item_id)}(?:\s|—|-|$)",
                target_text,
                flags=re.MULTILINE,
            )
        )
        if count != 1:
            errors.append(f"{item_id} must have exactly one normative detail heading; found {count}")
    for item_id in source_acceptance:
        count = len(re.findall(rf"\b{re.escape(item_id)}\b", target_text))
        if count < 2:
            errors.append(
                f"{item_id} must appear in the canonical ledger and at least one detailed acceptance/coverage mapping"
            )

    target_cf_ids = expand_cf_ranges(target_text)
    source_cf_ids = expand_cf_ranges(feature_text)
    unknown = sorted(target_cf_ids - set(conformance))
    if unknown:
        errors.append(f"unknown VS-000 conformance IDs: {', '.join(unknown)}")

    ledger = section(target_text, "Exact conformance semantics ledger")
    if ledger is None:
        errors.append("missing exact section heading: Exact conformance semantics ledger")
    else:
        ledger_rows = rows_by_id(ledger, CF_ID)
        for cf_id in sorted(target_cf_ids & set(conformance)):
            expected = expected_cf_row(conformance[cf_id])
            if ledger_rows.get(cf_id) != [expected]:
                errors.append(
                    f"Exact conformance semantics ledger must copy all seven registry fields for {cf_id} exactly once"
                )

    for cf_id in sorted(target_cf_ids & set(conformance)):
        entry = conformance[cf_id]
        if entry.get("owner") != feature and cf_id not in source_cf_ids:
            errors.append(
                f"{cf_id} is owned by {entry.get('owner')} and is not authorized by the feature authority"
            )
        semantics = " ".join(
            scalar(entry.get(key)).casefold()
            for key in ("inputs", "expectedState", "expectedError", "expectedSideEffects")
        )
        if entry.get("owner") == feature or cf_id.startswith("VS0-CF-HP"):
            continue
        for paragraph in re.split(r"\n\s*\n", target_text):
            if cf_id not in paragraph or "transition" not in paragraph.casefold():
                continue
            if re.search(r"invalid.{0,40}transition|transition.{0,40}invalid", paragraph, re.I):
                if re.search(
                    r"(?:does|do|must|is)\s+not.{0,100}(?:own|prove|govern|enforce|authority)",
                    paragraph,
                    re.I | re.S,
                ):
                    continue
                if "transition" not in semantics:
                    errors.append(
                        f"{cf_id} is used as transition evidence but its registry semantics do not govern a transition"
                    )
                    break

    for cf_id, entry in conformance.items():
        side_effect = scalar(entry.get("expectedSideEffects"))
        if (
            entry.get("owner") != feature
            and cf_id not in source_cf_ids
            and side_effect.casefold() not in {"null", "none", "no side effects"}
            and len(side_effect) >= 12
            and side_effect in target_text
        ):
            errors.append(
                f"downstream side-effect contract leaked from {cf_id} ({entry.get('owner')}): {side_effect}"
            )


def check_coverage(
    errors: list[str], feature_text: str, target_text: str, stage: str
) -> None:
    requirements = rows_by_id(feature_text, REQ_ID)
    acceptance = rows_by_id(feature_text, AC_ID)
    title = "Canonical coverage ledger"
    ledger = section(target_text, title)
    if ledger is None:
        errors.append(f"missing exact section heading: {title}")
        return
    for item_id in [*requirements, *acceptance]:
        count = len(re.findall(rf"\b{re.escape(item_id)}\b", ledger))
        if count != 1:
            errors.append(f"{stage} {title} must map {item_id} exactly once; found {count}")


def check_tasks(errors: list[str], target_text: str, tier: str) -> None:
    matches = list(
        re.finditer(r"^(#{2,6})\s+Task\s+(\d+)\b.*$", target_text, re.MULTILINE)
    )
    if tier == "A" and not 6 <= len(matches) <= 12:
        errors.append(f"Tier A tasks must contain 6..12 Task headings; found {len(matches)}")
    for index, match in enumerate(matches):
        end = matches[index + 1].start() if index + 1 < len(matches) else len(target_text)
        body = target_text[match.start():end]
        task_id = match.group(2)
        for label in (
            "Writable paths:",
            "Tests:",
            "Verification commands:",
            "Acceptance criteria:",
            "Commit message:",
        ):
            count = body.count(label)
            if count != 1:
                errors.append(f"Task {task_id} must contain exactly one '{label}' field; found {count}")


def check_stage(root: Path, control: dict[str, Any], stage: str) -> list[str]:
    errors: list[str] = []
    feature = control["feature"]["id"]
    feature_path = root / control["feature"]["feature_file"]
    target_path = root / control["feature"]["spec_path"] / f"{stage}.md"
    if not target_path.is_file():
        return [f"missing Kiro stage output: {target_path.relative_to(root)}"]

    feature_text = feature_path.read_text()
    target_text = target_path.read_text()
    if target_text.count(RECEIPT) != 1:
        errors.append(f"stage document must contain exactly one '{RECEIPT}' receipt")
    if re.search(r"\b(?:TODO|TBD|FIXME)\b", target_text, re.I):
        errors.append("stage document contains unresolved TODO/TBD/FIXME markers")

    known_registry_ids = set(
        re.findall(
            r"\bid:\s*(VS0-(?:SCHEMA|WRITER|STATE|CF)-[A-Z0-9-]+)",
            (root / "docs/architecture/vertical-slices/VS-000-contract-registry.yaml").read_text(),
        )
    )
    unknown_registry_ids = sorted(set(VS_ID.findall(target_text)) - known_registry_ids)
    if unknown_registry_ids:
        errors.append(f"unknown VS-000 registry IDs: {', '.join(unknown_registry_ids)}")

    conformance = registry_conformance(root)
    if stage == "requirements":
        check_requirements(errors, feature, feature_text, target_text, conformance)
    else:
        check_coverage(errors, feature_text, target_text, stage)
        if stage == "design":
            for classification in ("IMPLEMENT", "CONTRACT_ONLY/NO_TASK", "EXCLUDED"):
                if classification not in target_text:
                    errors.append(f"design is missing required classification: {classification}")
        if stage == "tasks":
            check_tasks(errors, target_text, control["feature"]["tier"])
    return errors


def write_revision_prompt(path: Path, feature: str, stage: str, errors: list[str]) -> None:
    target = f"{stage}.md"
    lines = [
        f"Revise {target} only for {feature}. Do not modify any other file or advance stages.",
        "",
        "These deterministic semantic guardrail failures must all be corrected:",
    ]
    lines.extend(f"{index}. {error}" for index, error in enumerate(errors, 1))
    lines.extend(
        [
            "",
            "Anti-drift rules:",
            "- The approved feature REQ and AC tables are immutable semantic ledgers.",
            "- Copy their rows exactly into the required canonical ledgers; do not repurpose, renumber, merge, split, or paraphrase ledger rows.",
            "- Copy each referenced conformance row's owner, inputs, expectedState, expectedError, expectedSideEffects, and gate exactly from the VS-000 registry.",
            "- A conformance case proves only its registered scenario; a state-machine ID alone governs transitions.",
            "- Keep stale and excluded terms inside explicit Non-goals/Out of Scope sections and do not repeat them in active acceptance text.",
            "- Do not import downstream fields, metrics, side effects, or runtime behavior.",
            "- Preserve exactly one STAGE_STATUS: COMPLETE receipt after reading the complete revised file.",
            "",
            f"Before returning, run: ./scripts/kiro-semantic-check.py --feature {feature} --stage {stage}",
        ]
    )
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text("\n".join(lines) + "\n")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--feature", required=True)
    parser.add_argument("--stage", choices=STAGES, required=True)
    parser.add_argument("--write-revision-prompt")
    args = parser.parse_args()

    control_path = ROOT / ".automation/features" / f"{args.feature}.control.json"
    try:
        control = json.loads(control_path.read_text())
        errors = check_stage(ROOT, control, args.stage)
    except (OSError, ValueError, RuntimeError) as exc:
        print(f"FAIL: Kiro semantic guardrail could not run: {exc}", file=sys.stderr)
        raise SystemExit(1) from exc

    scope = subprocess.run(
        ["./scripts/phase2-scope-check.sh", args.feature],
        cwd=ROOT,
        check=False,
        capture_output=True,
        text=True,
    )
    if scope.returncode:
        detail = (scope.stdout + scope.stderr).strip()
        errors.append(
            "phase2 scope check failed; keep excluded/stale concepts only in explicitly allowed non-goal context"
            + (f"\n{detail}" if detail else "")
        )
    elif scope.stdout.strip():
        print(scope.stdout.strip())

    if errors:
        print(f"FAIL: {args.feature} {args.stage} semantic guardrails")
        for error in errors:
            print(f"- {error}")
        if args.write_revision_prompt:
            path = ROOT / args.write_revision_prompt
            write_revision_prompt(path, args.feature, args.stage, errors)
            try:
                print(path.relative_to(ROOT))
            except ValueError:
                print(path)
        raise SystemExit(1)

    print(f"PASS: {args.feature} {args.stage} semantic guardrails")


if __name__ == "__main__":
    main()
