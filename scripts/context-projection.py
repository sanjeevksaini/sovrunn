#!/usr/bin/env python3
"""Generate a compact, source-hash-bound AI context projection."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]


def digest(payload: bytes) -> str:
    return hashlib.sha256(payload).hexdigest()


def repo_path(raw: str) -> Path:
    path = Path(raw)
    if path.is_absolute() or ".." in path.parts or not path.parts:
        raise SystemExit(f"ERROR: unsafe repository path: {raw}")
    resolved = (ROOT / path).resolve()
    try:
        resolved.relative_to(ROOT.resolve())
    except ValueError as exc:
        raise SystemExit(f"ERROR: path escapes repository: {raw}") from exc
    return resolved


def select_lines(lines: list[str], selection: dict[str, object]) -> tuple[str, str]:
    kind = selection.get("kind")
    label = str(selection.get("label", kind))
    if kind == "line_range":
        start = selection.get("start")
        end = selection.get("end")
        if not isinstance(start, int) or not isinstance(end, int) or start < 1 or end < start or end > len(lines):
            raise SystemExit(f"ERROR: invalid line range {start}..{end} for {label}")
        return label, "\n".join(lines[start - 1 : end]) + "\n"
    if kind == "matching_lines":
        raw_pattern = selection.get("pattern")
        if not isinstance(raw_pattern, str) or not raw_pattern:
            raise SystemExit(f"ERROR: missing matching-lines pattern for {label}")
        pattern = re.compile(raw_pattern)
        matched = [line for line in lines if pattern.search(line)]
        if not matched:
            raise SystemExit(f"ERROR: selector produced no lines for {label}")
        return label, "\n".join(matched) + "\n"
    raise SystemExit(f"ERROR: unsupported selector kind: {kind}")


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--spec", required=True)
    args = parser.parse_args()

    spec_path = repo_path(args.spec)
    spec = json.loads(spec_path.read_text())
    if spec.get("schema_version") != "1.0":
        raise SystemExit("ERROR: projection schema_version must be 1.0")

    output_path = repo_path(str(spec["output"]))
    manifest_path = repo_path(str(spec["manifest"]))
    rendered = [
        "---",
        "doc_type: ai_context_projection",
        f"feature: {spec['feature']}",
        f"stage: {spec['stage']}",
        "authority: non-authoritative-exact-excerpts",
        "generated: true",
        "---",
        "",
        f"# {spec['feature']} {str(spec['stage']).title()} Context Projection",
        "",
        "This is a non-authoritative, mechanically generated projection of exact",
        "approved-source excerpts. The FEATURE architecture and controlling handoff",
        "package remain authoritative. A source-hash mismatch blocks regeneration.",
        "",
    ]
    source_records: list[dict[str, object]] = []

    for source in spec["sources"]:
        source_path = repo_path(str(source["path"]))
        payload = source_path.read_bytes()
        actual = digest(payload)
        expected = str(source["sha256"])
        if actual != expected:
            raise SystemExit(
                f"ERROR: approved source hash mismatch for {source['path']}: {actual} != {expected}"
            )
        lines = payload.decode("utf-8").splitlines()
        selector_records = []
        rendered.extend(
            [
                f"## Source: `{source['path']}`",
                "",
                f"Approved SHA-256: `{actual}`",
                "",
            ]
        )
        for selection in source["selections"]:
            label, excerpt = select_lines(lines, selection)
            rendered.extend(
                [
                    f"### {label}",
                    "",
                    "<!-- BEGIN EXACT APPROVED-SOURCE EXCERPT -->",
                    excerpt.rstrip("\n"),
                    "<!-- END EXACT APPROVED-SOURCE EXCERPT -->",
                    "",
                ]
            )
            selector_records.append({**selection, "excerpt_sha256": digest(excerpt.encode("utf-8"))})
        source_records.append(
            {
                "path": source["path"],
                "sha256": actual,
                "bytes": len(payload),
                "selections": selector_records,
            }
        )

    output_payload = ("\n".join(rendered).rstrip() + "\n").encode("utf-8")
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_bytes(output_payload)
    manifest = {
        "schema_version": "1.0",
        "feature": spec["feature"],
        "stage": spec["stage"],
        "projection": str(output_path.relative_to(ROOT)),
        "projection_bytes": len(output_payload),
        "projection_estimated_tokens": (len(output_payload) + 3) // 4,
        "projection_sha256": digest(output_payload),
        "sources": source_records,
    }
    manifest_path.parent.mkdir(parents=True, exist_ok=True)
    manifest_path.write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n")
    print(f"projection={output_path.relative_to(ROOT)}")
    print(f"manifest={manifest_path.relative_to(ROOT)}")
    print(f"sha256={manifest['projection_sha256']}")
    print(f"bytes={manifest['projection_bytes']}")
    print(f"estimated_tokens={manifest['projection_estimated_tokens']}")


if __name__ == "__main__":
    main()
