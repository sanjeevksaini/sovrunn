#!/usr/bin/env python3
"""Validate exactly one terminal-formatted Kiro or Cursor completion receipt."""

from __future__ import annotations

import argparse
import re
from pathlib import Path

ANSI_ESCAPE = re.compile(r"\x1b(?:\[[0-?]*[ -/]*[@-~]|\][^\x07]*(?:\x07|\x1b\\))")
PATTERNS = {
    "stage": re.compile(r"STAGE_STATUS: (?:COMPLETE|BLOCKED(?: [A-Z_]+)?)\Z"),
    "task": re.compile(r"TASK_STATUS: (?:COMPLETE|BLOCKED)\Z"),
}
EXPECTED = {
    "stage": "STAGE_STATUS: COMPLETE",
    "task": "TASK_STATUS: COMPLETE",
}


def normalized_lines(raw: str) -> list[str]:
    return [ANSI_ESCAPE.sub("", line).removesuffix("\r") for line in raw.splitlines()]


def find_receipts(raw: str, kind: str) -> list[str]:
    pattern = PATTERNS[kind]
    return [line for line in normalized_lines(raw) if pattern.fullmatch(line)]


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--log", required=True)
    parser.add_argument(
        "--document",
        help="stage document fallback when the CLI summary omits an exact receipt",
    )
    parser.add_argument("--kind", required=True, choices=tuple(PATTERNS))
    args = parser.parse_args()
    path = Path(args.log)
    if not path.is_file():
        raise SystemExit(f"ERROR: receipt log does not exist: {path}")
    receipts = find_receipts(path.read_text(errors="replace"), args.kind)
    expected = EXPECTED[args.kind]
    if len(receipts) == 1 and receipts[0] == expected:
        print(f"PASS: exactly one {expected} in CLI log")
        return
    if receipts:
        raise SystemExit(
            f"ERROR: expected exactly one {args.kind} receipt, found {len(receipts)}: {receipts}"
        )
    if not args.document:
        raise SystemExit(
            f"ERROR: expected exactly one {args.kind} receipt, found 0: []"
        )
    document = Path(args.document)
    if not document.is_file():
        raise SystemExit(f"ERROR: receipt document does not exist: {document}")
    document_receipts = find_receipts(
        document.read_text(errors="replace"), args.kind
    )
    if document_receipts != [expected]:
        raise SystemExit(
            "ERROR: CLI log omitted its receipt and document fallback did not "
            f"contain exactly one completion receipt: {document_receipts}"
        )
    print(f"PASS: exactly one {expected} in stage document fallback")


if __name__ == "__main__":
    main()
