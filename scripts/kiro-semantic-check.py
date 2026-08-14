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
# Conformance IDs have both historical compact forms (for example, MIGF01)
# and feature-local forms (for example, F15-01).  The word boundaries prevent
# a local ID from also being counted as its non-existent F15 prefix.
CF_ID = re.compile(r"\bVS0-CF-(?:[A-Z]+\d+(?:-\d+)?)\b")
VS_ID = re.compile(r"VS0-(?:SCHEMA|WRITER|STATE|CF)-[A-Z0-9-]+")


def _mention_at_position_allowed(text: str, position: int) -> bool:
    """True if the text preceding the given offset (within a bounded lookback
    window) contains an explicit exclusion/retirement/non-goal/tombstone
    marker, meaning the mention at that position is historical record rather
    than an active conformance claim."""
    lower = text.casefold()
    window = lower[max(0, position - 160):position]
    allow_markers = (
        "must not", "no longer", "superseded", "retired", "removed", "does not",
        "not implement", "excluded", "non-goal", "tombstone", "never",
        "must never", "reject",
    )
    return any(marker in window for marker in allow_markers)


def table_cells(line: str) -> tuple[str, ...] | None:
    stripped = line.strip()
    if not (stripped.startswith("|") and stripped.endswith("|")):
        return None
    return tuple(cell.strip() for cell in stripped[1:-1].split("|"))


def table_header_cells(text: str) -> tuple[str, ...] | None:
    """Return the header row of the first markdown table in text, identified by
    the standard header/separator pattern (a row immediately followed by a
    ``|---|---|...|`` separator row), rather than assuming a fixed column
    position. This makes ledger comparisons robust to column reordering."""
    lines = text.splitlines()
    for index in range(len(lines) - 1):
        header = table_cells(lines[index])
        separator = table_cells(lines[index + 1])
        if not header or not separator:
            continue
        if all(re.fullmatch(r":?-{1,}:?", cell) for cell in separator):
            return header
    return None


def normalize_header(name: str) -> str:
    return re.sub(r"[^a-z0-9]", "", name.casefold())


# Maps a normalized registry conformance table-header name to the registry
# field it represents. Matching by normalized header name (instead of a fixed
# presentation-column index) keeps the comparison correct regardless of the
# order columns are presented in a target document's table.
CF_FIELD_BY_HEADER = {
    "id": "id",
    "owner": "owner",
    "inputs": "inputs",
    "expectedstate": "expectedState",
    "expectederror": "expectedError",
    "expectedviolation": "expectedViolation",
    "expectedsideeffects": "expectedSideEffects",
    "gate": "gate",
}
CF_REQUIRED_HEADER_FIELDS = (
    "id", "owner", "inputs", "expectedState", "expectedError", "expectedSideEffects", "gate",
)


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
    return set(expand_cf_ranges_with_positions(text))


def expand_cf_ranges_with_positions(text: str) -> dict[str, int]:
    """Like expand_cf_ranges, but also records the text offset of the mention
    (direct ID or range expression) that produced each expanded ID, so callers
    can check whether that specific mention occurs in an allowed exclusion
    context (e.g. a retired-ID range spelled out in a tombstone/non-goal note)."""
    found: dict[str, int] = {}
    for match in CF_ID.finditer(text):
        found.setdefault(match.group(0), match.start())
    local_range_pattern = re.compile(r"\bVS0-CF-([A-Z]+\d+)-(\d+)\.\.(\d+)\b")
    for match in local_range_pattern.finditer(text):
        prefix, left_number, right_number = match.groups()
        width = max(len(left_number), len(right_number))
        for number in range(int(left_number), int(right_number) + 1):
            found.setdefault(f"VS0-CF-{prefix}-{number:0{width}d}", match.start())
    range_pattern = re.compile(r"VS0-CF-([A-Z]+)(\d+)\.\.([A-Z]*)(\d+)")
    for match in range_pattern.finditer(text):
        left_prefix, left_number, right_prefix, right_number = match.groups()
        prefix = right_prefix or left_prefix
        if prefix != left_prefix:
            continue
        width = max(len(left_number), len(right_number))
        for number in range(int(left_number), int(right_number) + 1):
            found.setdefault(f"VS0-CF-{prefix}{number:0{width}d}", match.start())
    return found


def req_detail_heading_matches(item_id: str, target_text: str) -> int:
    """Count normative REQ detail headings for item_id.  A valid heading is a
    Markdown heading (``##``..``######``) whose text is item_id itself,
    optionally preceded by an approved Markdown section-number prefix (for
    example ``4.1``, ``4.12``, or no numeric prefix at all) before the REQ
    ID -- for example ``### 4.1 REQ-F15-01 -- ...`` or ``### REQ-F15-01
    ...``.  Only an exact match of this pattern anchored at the start of a
    heading line counts; a duplicate heading is counted again (so callers can
    reject a count other than exactly one), and an incidental in-body mention
    of the REQ ID that is not itself a heading is never counted."""
    pattern = re.compile(
        rf"^#{{2,6}}\s+(?:\d+(?:\.\d+)*[.)]?\s+)?{re.escape(item_id)}(?:\s|—|-|$)",
        re.MULTILINE,
    )
    return len(pattern.findall(target_text))


def owned_conformance_ids(
    feature_text: str,
    authority_text: str,
    conformance: dict[str, dict[str, Any]],
    feature: str,
) -> set[str]:
    """Return the F0015-style feature-owned conformance IDs authorized for the
    Exact conformance semantics ledger.  Authorization is derived from the
    union of every conformance ID mentioned in the current feature authority
    and in the control-manifest architecture authority: a required proof
    matrix (for example ADH-2026-046's complete F15-01..25 cases) may be
    mapped in the manifest-controlled architecture authority's traceability
    section rather than repeated verbatim in the feature authority.  Only
    registry cases whose registered owner is this feature are included;
    downstream, unowned, or unknown IDs mentioned in either authority are
    never authorized by this function."""
    mentioned = expand_cf_ranges(feature_text) | expand_cf_ranges(authority_text)
    return {
        cf_id
        for cf_id in mentioned
        if cf_id in conformance and conformance[cf_id].get("owner") == feature
    }


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
        scalar(entry.get("expectedViolation", "—")),
        scalar(entry.get("expectedSideEffects")),
        scalar(entry.get("gate")),
    )


def parse_ledger_rows_by_header(ledger_text: str) -> dict[str, dict[str, str]] | None:
    """Parse the conformance ledger table using its own header row rather than
    assuming a fixed presentation-column order. Returns a mapping of
    conformance ID -> {registry field name -> cell value}, using only the
    columns the ledger's header names identify (unrecognized/reordered extra
    columns, such as expectedViolation appearing before expectedError, are
    handled correctly). Returns None if no table with an 'id' header is found.
    """
    header = table_header_cells(ledger_text)
    if header is None:
        return None
    normalized_header = [normalize_header(cell) for cell in header]
    if "id" not in normalized_header:
        return None
    id_index = normalized_header.index("id")
    field_by_index: dict[int, str] = {}
    for index, cell in enumerate(normalized_header):
        field = CF_FIELD_BY_HEADER.get(cell)
        if field:
            field_by_index[index] = field

    rows: dict[str, dict[str, str]] = {}
    lines = ledger_text.splitlines()
    started = False
    for line in lines:
        cells = table_cells(line)
        if cells is None:
            continue
        if not started:
            if cells == header:
                started = True
            continue
        if all(re.fullmatch(r":?-{1,}:?", cell) for cell in cells):
            continue
        if id_index >= len(cells) or not CF_ID.fullmatch(cells[id_index]):
            continue
        cf_id = cells[id_index]
        row_fields = {
            field: cells[index]
            for index, field in field_by_index.items()
            if index < len(cells)
        }
        rows[cf_id] = row_fields
    return rows


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
    architecture_text: str = "",
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
        count = req_detail_heading_matches(item_id, target_text)
        if count != 1:
            errors.append(f"{item_id} must have exactly one normative detail heading; found {count}")
    for item_id in source_acceptance:
        count = len(re.findall(rf"\b{re.escape(item_id)}\b", target_text))
        if count < 2:
            errors.append(
                f"{item_id} must appear in the canonical ledger and at least one detailed acceptance/coverage mapping"
            )

    # A feature authority may name a downstream conformance case only to
    # exclude it.  Authorization for the Exact conformance semantics ledger is
    # derived from the union of the current feature authority and the
    # control-manifest architecture authority (ADH-2026-046 §3): a required
    # proof matrix may be mapped in the architecture authority's traceability
    # section rather than repeated verbatim in the feature authority.  Only
    # cases owned by this feature can authorize its ledger.
    source_cf_ids = expand_cf_ranges(feature_text) | expand_cf_ranges(architecture_text)
    source_owned_cf_ids = owned_conformance_ids(feature_text, architecture_text, conformance, feature)
    target_cf_positions = expand_cf_ranges_with_positions(target_text)
    target_cf_ids = set(target_cf_positions)
    # A retired/tombstone ID mentioned solely inside an explicit non-goal,
    # exclusion, or tombstone context (e.g. "retired tombstones that must
    # never be reused: VS0-CF-MIG01..MIGF03...") is historical record, not an
    # active conformance claim, and must not be misclassified as
    # unknown/active. The mention is checked at the exact source position
    # (direct ID or range expression) that produced the expanded ID, so an ID
    # only reachable through a retired range expression is still recognized.
    unknown = sorted(
        cf_id
        for cf_id in (target_cf_ids - set(conformance))
        if not _mention_at_position_allowed(target_text, target_cf_positions[cf_id])
    )
    if unknown:
        errors.append(f"unknown VS-000 conformance IDs: {', '.join(unknown)}")

    ledger = section(target_text, "Exact conformance semantics ledger")
    if ledger is None:
        errors.append("missing exact section heading: Exact conformance semantics ledger")
    else:
        ledger_rows = parse_ledger_rows_by_header(ledger)
        if ledger_rows is None:
            errors.append(
                "Exact conformance semantics ledger table must have a recognizable 'id' header column"
            )
            ledger_rows = {}
        expected_ledger_ids = source_owned_cf_ids
        extra_ledger_ids = sorted(set(ledger_rows) - expected_ledger_ids)
        if extra_ledger_ids:
            errors.append(
                "Exact conformance semantics ledger contains cases not authorized by the feature authority: "
                + ", ".join(extra_ledger_ids)
            )
        for cf_id in sorted(expected_ledger_ids):
            entry = conformance[cf_id]
            actual_row = ledger_rows.get(cf_id)
            if actual_row is None:
                errors.append(
                    f"Exact conformance semantics ledger must copy all eight registry fields for {cf_id} exactly once"
                )
                continue
            mismatched_fields = [
                field
                for field in CF_REQUIRED_HEADER_FIELDS
                if field not in actual_row or scalar(entry.get(field)) != actual_row.get(field)
            ]
            if "expectedViolation" in entry and (
                "expectedViolation" not in actual_row
                or scalar(entry.get("expectedViolation")) != actual_row.get("expectedViolation")
            ):
                mismatched_fields.append("expectedViolation")
            if mismatched_fields:
                errors.append(
                    f"Exact conformance semantics ledger must copy all eight registry fields for {cf_id} exactly once "
                    f"(mismatched by normalized header name: {', '.join(mismatched_fields)})"
                )

    for cf_id in sorted(source_owned_cf_ids):
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
        architecture_rel = control["feature"].get("architecture", "")
        architecture_text = ""
        if architecture_rel:
            architecture_path = root / architecture_rel
            if architecture_path.is_file():
                architecture_text = architecture_path.read_text()
        check_requirements(errors, feature, feature_text, target_text, conformance, architecture_text)
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


def self_test() -> None:
    """Deterministic, offline regression coverage for the two checker defects
    corrected here: (1) the normative REQ detail heading matcher must accept
    an optional approved Markdown section-number prefix before the REQ ID
    while still rejecting duplicates and incidental in-body mentions; and (2)
    Exact conformance semantics ledger authorization must be derived from the
    union of the feature authority and the control-manifest architecture
    authority, filtered to registry cases owned by the feature."""
    failures: list[str] = []

    def check(label: str, condition: bool) -> None:
        if not condition:
            failures.append(label)

    # --- Case 1: normative REQ detail heading matcher ---
    check(
        "accepts a numbered section-number prefix (4.1)",
        req_detail_heading_matches("REQ-F15-01", "### 4.1 REQ-F15-01 — Register CloudPlatform") == 1,
    )
    check(
        "accepts a two-digit subsection number (4.12)",
        req_detail_heading_matches("REQ-F15-01", "#### 4.12 REQ-F15-01 detail") == 1,
    )
    check(
        "accepts no numeric prefix at all",
        req_detail_heading_matches("REQ-F15-01", "### REQ-F15-01 — Register CloudPlatform") == 1,
    )
    check(
        "rejects zero headings (only an incidental in-body mention)",
        req_detail_heading_matches(
            "REQ-F15-01", "This paragraph mentions REQ-F15-01 without a heading."
        )
        == 0,
    )
    check(
        "rejects duplicate headings (counts both, so callers reject count != 1)",
        req_detail_heading_matches(
            "REQ-F15-01",
            "### 4.1 REQ-F15-01 — first\n\nbody\n\n### 4.9 REQ-F15-01 — duplicate\n",
        )
        == 2,
    )
    check(
        "does not match a different REQ ID sharing a numeric prefix",
        req_detail_heading_matches("REQ-F15-01", "### 4.1 REQ-F15-02 — Register CloudProvider") == 0,
    )

    # --- Case 2: exact conformance-ledger authority union ---
    conformance = {
        "VS0-CF-F15-13": {"owner": "FEATURE-0015"},
        "VS0-CF-F15-14": {"owner": "FEATURE-0015"},
        "VS0-CF-F15-16": {"owner": "FEATURE-0015"},
        "VS0-CF-F15-20": {"owner": "FEATURE-0015"},
        "VS0-CF-F15-21": {"owner": "FEATURE-0015"},
        "VS0-CF-F15-22": {"owner": "FEATURE-0015"},
        "VS0-CF-F10": {"owner": "FEATURE-0016"},
    }
    feature_text = "Owned locally: VS0-CF-F15-01, VS0-CF-F15-02. Excluded: must not use VS0-CF-F10."
    architecture_text = (
        "| REQ-F15-08 | VS0-CF-F15-13, VS0-CF-F15-14 |\n"
        "| REQ-F15-16 | VS0-CF-F15-16 |\n"
        "| REQ-F15-07 | VS0-CF-F15-20 |\n"
        "| AC-F15-... | VS0-CF-F15-21, VS0-CF-F15-22 |\n"
    )
    owned = owned_conformance_ids(feature_text, architecture_text, conformance, "FEATURE-0015")
    check(
        "accepts F15-13, F15-14, F15-16, F15-20, F15-21, F15-22 sourced only from the architecture authority",
        {
            "VS0-CF-F15-13",
            "VS0-CF-F15-14",
            "VS0-CF-F15-16",
            "VS0-CF-F15-20",
            "VS0-CF-F15-21",
            "VS0-CF-F15-22",
        }
        <= owned,
    )
    check(
        "still rejects an unowned/downstream ID mentioned only for exclusion (VS0-CF-F10)",
        "VS0-CF-F10" not in owned,
    )
    check(
        "still rejects an unknown ID not present in the registry",
        "VS0-CF-F15-99" not in owned,
    )
    check(
        "owned set is empty when neither authority mentions any owned case",
        owned_conformance_ids("no mentions here", "no mentions here either", conformance, "FEATURE-0015")
        == set(),
    )

    if failures:
        print(f"FAIL: kiro-semantic-check self-test — {len(failures)} failure(s)")
        for failure in failures:
            print(f"  ✗ {failure}")
        raise SystemExit(1)
    print(f"PASS: kiro-semantic-check self-test — {6 + 4} checks")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--feature")
    parser.add_argument("--stage", choices=STAGES)
    parser.add_argument("--write-revision-prompt")
    parser.add_argument(
        "--self-test",
        action="store_true",
        help="run deterministic offline regression checks and exit (no --feature/--stage required)",
    )
    args = parser.parse_args()

    if args.self_test:
        self_test()
        return

    if not args.feature or not args.stage:
        print("FAIL: --feature and --stage are required unless --self-test is given", file=sys.stderr)
        raise SystemExit(1)

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
