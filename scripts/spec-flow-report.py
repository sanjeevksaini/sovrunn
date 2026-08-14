#!/usr/bin/env python3
"""Create a deterministic report for a Feature Factory specification run."""

import argparse
import datetime as dt
import json
import os
import re
from pathlib import Path


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--feature", required=True)
    parser.add_argument("--spec-path", required=True)
    parser.add_argument("--status", required=True)
    parser.add_argument("--started-epoch", required=True, type=int)
    parser.add_argument("--kiro-invocations", required=True, type=int)
    parser.add_argument("--baseline-lines", default=0, type=int)
    parser.add_argument("--action", default="")
    args = parser.parse_args()

    review_epoch = os.environ.get("FEATURE_FACTORY_REVIEW_EPOCH", "")
    reviewer_adapter = os.environ.get("FEATURE_FACTORY_REVIEWER_CMD", "unspecified")
    if review_epoch and not re.fullmatch(
        r"[A-Za-z0-9][A-Za-z0-9._-]{0,79}", review_epoch
    ):
        raise SystemExit(f"ERROR: invalid review epoch: {review_epoch!r}")

    now = dt.datetime.now(dt.timezone.utc)
    elapsed = max(0, int(now.timestamp()) - args.started_epoch)
    spec_path = Path(args.spec_path)
    review_root = Path(".automation/reviews") / args.feature
    history = review_root / "history"

    documents = {}
    total_lines = 0
    for stage in ("requirements", "design", "tasks"):
        path = spec_path / f"{stage}.md"
        lines = len(path.read_text().splitlines()) if path.exists() else 0
        documents[stage] = {"path": str(path), "exists": path.exists(), "lines": lines}
        total_lines += lines

    usage = {"input_tokens": 0, "output_tokens": 0, "total_tokens": 0}
    api_calls = 0
    for path in sorted(history.glob("*.openai.raw.json")) if history.exists() else []:
        try:
            raw = json.loads(path.read_text())
            row = raw.get("usage") or {}
            for key in usage:
                usage[key] += int(row.get(key) or 0)
            api_calls += 1
        except (OSError, ValueError, TypeError):
            pass

    revisions = {}
    verdicts = {}
    for stage in ("requirements", "design", "tasks"):
        count_name = (
            f"{stage}.{review_epoch}.revision-count"
            if review_epoch
            else f"{stage}.revision-count"
        )
        count_path = review_root / count_name
        revisions[stage] = int(count_path.read_text().strip()) if count_path.exists() else 0
        review_path = review_root / f"{stage}.review.json"
        if review_path.exists():
            try:
                verdicts[stage] = json.loads(review_path.read_text()).get("status", "UNKNOWN")
            except (OSError, ValueError):
                verdicts[stage] = "INVALID_REVIEW_JSON"
        else:
            verdicts[stage] = "NOT_REVIEWED"

    report = {
        "feature": args.feature,
        "review_epoch": review_epoch or None,
        "reviewer_adapter": reviewer_adapter,
        "status": args.status,
        "started_at": dt.datetime.fromtimestamp(args.started_epoch, dt.timezone.utc).isoformat(),
        "finished_at": now.isoformat(),
        "elapsed_seconds": elapsed,
        "documents": documents,
        "total_spec_lines": total_lines,
        "lines_generated_or_added_this_run": max(0, total_lines - args.baseline_lines),
        "revisions": revisions,
        "review_verdicts": verdicts,
        "openai": {"api_calls": api_calls, **usage},
        "kiro": {
            "routing": "Auto",
            "invocations": args.kiro_invocations,
            "credits_used": None,
            "credits_measurement": "unavailable_from_headless_cli; use interactive /usage or enterprise usage report",
        },
        "human_action_required": args.status != "COMPLETE",
        "required_action": args.action,
    }

    out_dir = Path(".automation/reports") / args.feature
    out_dir.mkdir(parents=True, exist_ok=True)
    json_path = out_dir / "spec-flow-latest.json"
    md_path = out_dir / "spec-flow-latest.md"
    json_path.write_text(json.dumps(report, indent=2, sort_keys=True) + "\n")

    md = [
        f"# {args.feature} specification automation report",
        "",
        f"- Status: **{args.status}**",
        f"- Review epoch: {review_epoch or 'legacy/default'}",
        f"- Elapsed: {elapsed} seconds",
        f"- Kiro routing: Auto",
        f"- Reviewer adapter: `{reviewer_adapter}`",
        f"- Kiro invocations: {args.kiro_invocations}",
        "- Kiro credits: unavailable from the headless CLI",
        f"- OpenAI review calls: {api_calls}",
        f"- OpenAI input/output/total tokens: {usage['input_tokens']} / {usage['output_tokens']} / {usage['total_tokens']}",
        f"- Total specification lines: {total_lines}",
        f"- Lines generated or added this run: {max(0, total_lines - args.baseline_lines)}",
        "",
        "## Stages",
        "",
    ]
    for stage in ("requirements", "design", "tasks"):
        md.append(
            f"- {stage}: {verdicts[stage]}; {documents[stage]['lines']} lines; "
            f"{revisions[stage]} revisions"
        )
    if args.status != "COMPLETE":
        md.extend(["", "## Human action required", "", args.action or "Inspect the latest logs and review JSON."])
    md_path.write_text("\n".join(md) + "\n")
    print(json_path)
    print(md_path)


if __name__ == "__main__":
    main()
