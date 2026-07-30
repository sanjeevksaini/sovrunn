#!/usr/bin/env python3
"""Produce a conservative, token-efficient semantic-delta review input."""

from __future__ import annotations

import argparse
import difflib
import hashlib
import json
import re
from pathlib import Path

ID_RE = re.compile(r"\b(?:F\d{2}-(?:REQ|AD|R|DD|DQ|NG|CF|SEC|COMPAT)-\d+|ADH-\d{4}-\d+|FEATURE-\d{4}|DEC-\d{4})\b")
NORMATIVE_RE = re.compile(r"\b(?:MUST(?: NOT)?|SHALL(?: NOT)?|REQUIRED|PROHIBITED|ONLY|EXACTLY)\b", re.IGNORECASE)
HEADING_RE = re.compile(r"^#{1,6}\s+(.+?)\s*$")


def digest(text: str) -> str:
    return hashlib.sha256(text.encode()).hexdigest()


def mechanical_normalize(text: str) -> str:
    """Normalize only changes proven not to alter Markdown semantics."""
    return "\n".join(line.rstrip(" \t") for line in text.splitlines()).rstrip("\n") + "\n"


def lines_matching(text: str, pattern: re.Pattern[str]) -> list[str]:
    return [line.strip() for line in text.splitlines() if pattern.search(line)]


def headings(text: str) -> set[str]:
    return {match.group(1) for line in text.splitlines() if (match := HEADING_RE.match(line))}


def identifiers(text: str) -> set[str]:
    return set(ID_RE.findall(text))


def report(baseline_path: Path | None, current_path: Path) -> dict[str, object]:
    current = current_path.read_text()
    if baseline_path is None:
        return {
            "classification": "INITIAL_REVIEW",
            "baseline": None,
            "current": str(current_path),
            "current_sha256": digest(current),
            "semantic_review_required": True,
            "summary": "No prior reviewed snapshot exists; perform a complete semantic review.",
        }
    baseline = baseline_path.read_text()
    if baseline == current:
        classification = "NO_CHANGE"
        semantic_required = False
    elif mechanical_normalize(baseline) == mechanical_normalize(current):
        classification = "MECHANICAL_ONLY"
        semantic_required = False
    else:
        classification = "SEMANTIC_REVIEW_REQUIRED"
        semantic_required = True
    old_ids, new_ids = identifiers(baseline), identifiers(current)
    old_headings, new_headings = headings(baseline), headings(current)
    old_normative = lines_matching(baseline, NORMATIVE_RE)
    new_normative = lines_matching(current, NORMATIVE_RE)
    diff = list(difflib.unified_diff(baseline.splitlines(), current.splitlines(), n=0))
    added = sum(1 for line in diff if line.startswith("+") and not line.startswith("+++"))
    removed = sum(1 for line in diff if line.startswith("-") and not line.startswith("---"))
    return {
        "classification": classification,
        "baseline": str(baseline_path),
        "current": str(current_path),
        "baseline_sha256": digest(baseline),
        "current_sha256": digest(current),
        "semantic_review_required": semantic_required,
        "changed_lines": {"added": added, "removed": removed},
        "identifiers": {"added": sorted(new_ids - old_ids), "removed": sorted(old_ids - new_ids)},
        "headings": {"added": sorted(new_headings - old_headings), "removed": sorted(old_headings - new_headings)},
        "normative_lines": {
            "added": sorted(set(new_normative) - set(old_normative))[:40],
            "removed": sorted(set(old_normative) - set(new_normative))[:40],
            "truncated": len(set(new_normative) ^ set(old_normative)) > 80,
        },
        "summary": (
            "Only trailing whitespace or final-newline changes detected."
            if classification == "MECHANICAL_ONLY"
            else "Review the compact identifier, heading, and normative-line delta; inspect the full diff only where necessary."
        ),
    }


def render_markdown(data: dict[str, object]) -> str:
    lines = [
        "## Semantic delta",
        "",
        f"- Classification: `{data['classification']}`",
        f"- Semantic review required: `{str(data['semantic_review_required']).lower()}`",
        f"- Summary: {data['summary']}",
    ]
    if "changed_lines" in data:
        changed = data["changed_lines"]
        lines.append(f"- Changed lines: +{changed['added']} / -{changed['removed']}")
        for section in ("identifiers", "headings"):
            values = data[section]
            lines.extend(
                [
                    f"- {section.title()} added: {', '.join(values['added']) or 'none'}",
                    f"- {section.title()} removed: {', '.join(values['removed']) or 'none'}",
                ]
            )
        normative = data["normative_lines"]
        lines.extend(["", "Normative lines added:"])
        lines.extend(f"- {line}" for line in normative["added"] or ["none"])
        lines.extend(["", "Normative lines removed:"])
        lines.extend(f"- {line}" for line in normative["removed"] or ["none"])
    return "\n".join(lines) + "\n"


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--baseline")
    parser.add_argument("--current", required=True)
    parser.add_argument("--json-out")
    parser.add_argument("--markdown-out")
    args = parser.parse_args()
    baseline = Path(args.baseline) if args.baseline else None
    current = Path(args.current)
    if baseline and not baseline.is_file():
        raise SystemExit(f"ERROR: baseline not found: {baseline}")
    if not current.is_file():
        raise SystemExit(f"ERROR: current file not found: {current}")
    data = report(baseline, current)
    rendered = json.dumps(data, indent=2, sort_keys=True) + "\n"
    if args.json_out:
        path = Path(args.json_out)
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(rendered)
    else:
        print(rendered, end="")
    if args.markdown_out:
        path = Path(args.markdown_out)
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(render_markdown(data))


if __name__ == "__main__":
    main()
