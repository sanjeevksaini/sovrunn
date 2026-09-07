#!/usr/bin/env python3
"""Prepare verified feature closeout metadata without committing or pushing."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

ROOT = Path(__file__).resolve().parents[1]


def git(*args: str) -> str:
    return subprocess.check_output(["git", *args], cwd=ROOT, text=True).strip()


def gh_pr(number: int) -> dict[str, Any]:
    raw = subprocess.check_output(
        [
            "gh",
            "pr",
            "view",
            str(number),
            "--json",
            "number,state,mergedAt,mergeCommit,baseRefName,headRefName,title",
        ],
        cwd=ROOT,
        text=True,
    )
    return json.loads(raw)


def replace_unique(text: str, pattern: str, replacement: str, *, label: str, flags: int = 0) -> str:
    updated, count = re.subn(pattern, replacement, text, flags=flags)
    if count != 1:
        raise SystemExit(f"ERROR: {label} expected one match, found {count}")
    return updated


def update_frontmatter(path: Path, values: dict[str, str]) -> None:
    text = path.read_text()
    if not text.startswith("---\n"):
        raise SystemExit(f"ERROR: missing front matter: {path.relative_to(ROOT)}")
    end = text.find("\n---\n", 4)
    if end < 0:
        raise SystemExit(f"ERROR: unterminated front matter: {path.relative_to(ROOT)}")
    front, body = text[4:end], text[end + 5 :]
    lines = front.splitlines()
    for key, value in values.items():
        matches = [index for index, line in enumerate(lines) if line.startswith(f"{key}:")]
        if len(matches) > 1:
            raise SystemExit(f"ERROR: duplicate front-matter key {key}")
        rendered = f"{key}: {value}"
        if matches:
            lines[matches[0]] = rendered
        else:
            lines.append(rendered)
    path.write_text("---\n" + "\n".join(lines) + "\n---\n" + body)


def validate_frontmatter(path: Path) -> None:
    text = path.read_text()
    if not text.startswith("---\n") or text.find("\n---\n", 4) < 0:
        raise SystemExit(f"ERROR: invalid front matter: {path.relative_to(ROOT)}")


def table_row(text: str, feature: str) -> tuple[str, list[str]]:
    matches = [line for line in text.splitlines() if re.match(rf"^\|\s*{re.escape(feature)}\s*\|", line)]
    if len(matches) != 1:
        raise SystemExit(f"ERROR: expected one {feature} table row, found {len(matches)}")
    line = matches[0]
    columns = [part.strip() for part in line.strip().strip("|").split("|")]
    return line, columns


def update_feature_index(path: Path, feature: str, evidence: str) -> None:
    text = path.read_text()
    line, columns = table_row(text, feature)
    if len(columns) != 6:
        raise SystemExit("ERROR: unexpected FEATURE_INDEX table shape")
    columns[3] = "Implemented and Merged"
    if evidence not in columns[5]:
        columns[5] = columns[5].rstrip(". ") + "; " + evidence + "."
    replacement = "| " + " | ".join(columns) + " |"
    path.write_text(text.replace(line, replacement, 1))


def update_traceability(path: Path, feature: str, gate: str) -> None:
    text = path.read_text()
    line, columns = table_row(text, feature)
    columns[2] = "Implemented and Merged"
    if len(columns) == 6:
        # Phase 2R's canonical matrix is Feature, Phase, Status, Decisions,
        # Controlling Source, Notes. Preserve its authority columns and make
        # the final-gate evidence explicit in Notes.
        columns[5] = gate
    elif len(columns) == 8:
        # Retain support for the legacy generic matrix used by older features.
        columns[6] = gate
    else:
        raise SystemExit("ERROR: unsupported traceability table shape; expected 6 or 8 columns")
    replacement = "| " + " | ".join(columns) + " |"
    path.write_text(text.replace(line, replacement, 1))


def update_phase_context(path: Path, data: dict[str, Any], evidence: str, date: str) -> None:
    feature = data["feature"]
    text = path.read_text()
    phase_goals = re.findall(r"(?m)^## (Phase [A-Za-z0-9]+) Goal\s*$", text)
    if len(phase_goals) != 1:
        raise SystemExit("ERROR: CURRENT_PHASE_CONTEXT must contain exactly one Phase <name> Goal heading")
    phase_heading = phase_goals[0]
    status_line = f"{feature['id']} status: implemented and merged; {evidence}."
    pattern = rf"(?m)^{re.escape(feature['id'])} status:.*$"
    if re.search(pattern, text):
        text = re.sub(pattern, status_line, text, count=1)
    else:
        anchor = f"## {phase_heading} Goal"
        if anchor not in text:
            raise SystemExit(f"ERROR: CURRENT_PHASE_CONTEXT missing {phase_heading} Goal")
        text = text.replace(anchor, status_line + "\n\n" + anchor, 1)
    row = f"| {feature['id']} {feature['title']} | Implemented and merged; {evidence} | Final feature gate passed {date} |"
    completed_heading = f"## {phase_heading} Completed Features"
    next_heading = f"## {phase_heading} Next Planned Feature"
    if completed_heading not in text or next_heading not in text:
        raise SystemExit("ERROR: CURRENT_PHASE_CONTEXT missing completed/next sections")
    before_completed, remainder = text.split(completed_heading, 1)
    completed_body, after_completed = remainder.split(next_heading, 1)
    completed_pattern = rf"(?m)^\|\s*{re.escape(feature['id'])}\s+.*$"
    if re.search(completed_pattern, completed_body):
        completed_body = re.sub(completed_pattern, row, completed_body, count=1)
    else:
        completed_body = completed_body.rstrip() + "\n" + row + "\n\n"
    next_body, separator, rest = after_completed.partition("\n## ")
    next_body = re.sub(completed_pattern, "", next_body)
    next_id = f"FEATURE-{int(feature['id'].split('-')[1]) + 1:04d}"
    index_text = (ROOT / "docs/features/FEATURE_INDEX.md").read_text()
    try:
        _, next_columns = table_row(index_text, next_id)
    except SystemExit:
        next_columns = []
    data_rows = [
        line for line in next_body.splitlines()
        if line.startswith("|") and not re.match(r"^\|\s*(?:Feature|---)", line)
    ]
    for line in data_rows:
        next_body = next_body.replace(line, "")
    if next_columns:
        next_row = f"| {next_id} {next_columns[1]} | Architecture not started | Pending architecture decision handoff |"
        next_body = next_body.rstrip() + "\n" + next_row + "\n"
    text = before_completed + completed_heading + completed_body + next_heading + next_body
    if separator:
        text += separator + rest
    path.write_text(text)


def ownership_summary(data: dict[str, Any]) -> str:
    owned = ", ".join(data["ownership"]["owned_resources"])
    excluded = "; ".join(
        f"{item['id']} ({', '.join(item['concepts'])})" for item in data["ownership"]["excluded_features"]
    )
    return f"It owns {owned}. Adjacent ownership remains excluded: {excluded}."


def update_architecture_baseline(path: Path, data: dict[str, Any], evidence: str) -> None:
    feature = data["feature"]
    handoffs = ", ".join(f"`{Path(item).stem.split('-feature-')[0]}`" for item in feature["handoffs"])
    replacement = (
        f"Completed and merged: `{feature['id']}: {feature['title']}` — {evidence}. "
        f"The approved architecture `{feature['architecture']}` and handoff package {handoffs} remain controlling. "
        + ownership_summary(data)
    )
    text = path.read_text()
    pattern = rf"(?m)^(?:Approved|Completed)[^\n]*`{re.escape(feature['id'])}:?[^\n]*$"
    if re.search(pattern, text):
        text = re.sub(pattern, replacement, text, count=1)
    else:
        phase2r_section = "## Implemented Phase 2R Features"
        phase2r_anchor = "\n## Phase 2R Next Planned Feature"
        if phase2r_section in text:
            text = text.replace(phase2r_section, phase2r_section + "\n\n" + replacement, 1)
        elif phase2r_anchor in text:
            text = text.replace(
                phase2r_anchor,
                "\n" + phase2r_section + "\n\n" + replacement + "\n" + phase2r_anchor,
                1,
            )
        else:
            anchor = "\n## Approved Consolidated"
            if anchor not in text:
                raise SystemExit("ERROR: architecture baseline has no supported closeout insertion anchor")
            text = text.replace(anchor, "\n" + replacement + anchor, 1)
    path.write_text(text)


def update_decision_summary(path: Path, data: dict[str, Any], evidence: str) -> None:
    feature = data["feature"]
    replacement = (
        f"- {feature['id']} is implemented and merged; {evidence}. "
        f"Its approved architecture and handoff package remain controlling. "
        + ownership_summary(data)
    )
    text = path.read_text()
    # A feature decision is commonly a wrapped Markdown list item. Replace the
    # entire item, not just its first physical line, so stale continuation text
    # cannot survive a closeout update.
    pattern = rf"(?m)^- {re.escape(feature['id'])}\b.*(?:\n(?:[ \t]+.*))*"
    if re.search(pattern, text):
        text = re.sub(pattern, replacement, text, count=1)
    else:
        anchor = "## Accepted Decisions\n"
        if anchor not in text:
            raise SystemExit("ERROR: decision summary missing accepted-decisions section")
        text = text.replace(anchor, anchor + "\n" + replacement + "\n", 1)

    stage_heading = f"## Approved {feature['id']} stage boundary"
    if stage_heading in text:
        closeout_heading = f"## {feature['id']} implementation closeout"
        closeout = (
            f"{closeout_heading}\n\n"
            f"{feature['id']} is implemented and merged; {evidence}. Its final human approval "
            "and feature gate are recorded. Its approved architecture and accepted handoffs remain controlling. "
            + ownership_summary(data)
        )
        text = replace_unique(
            text,
            rf"(?ms)^{re.escape(stage_heading)}\n\n.*?(?=^## |\Z)",
            closeout + "\n",
            label="FEATURE-0018 stage-boundary section",
        )
    path.write_text(text)


def preflight_closeout_targets(feature_path: Path, targets: dict[str, Path], feature: str) -> None:
    """Validate every closeout source before any --apply mutation occurs."""
    required = {
        "FEATURE_INDEX.md",
        "CURRENT_PHASE_CONTEXT.md",
        "CURRENT_ARCHITECTURE_BASELINE.md",
        "CURRENT_DECISION_SUMMARY.md",
        "FEATURE_TRACEABILITY_MATRIX.md",
        "VS-000_CONTRACT_TRACEABILITY_MATRIX.md",
    }
    missing = sorted(required - targets.keys())
    if missing:
        raise SystemExit(f"ERROR: closeout metadata targets missing: {', '.join(missing)}")
    validate_frontmatter(feature_path)
    _, index_columns = table_row(targets["FEATURE_INDEX.md"].read_text(), feature)
    if len(index_columns) != 6:
        raise SystemExit("ERROR: unexpected FEATURE_INDEX table shape")
    phase_text = targets["CURRENT_PHASE_CONTEXT.md"].read_text()
    if len(re.findall(r"(?m)^## Phase [A-Za-z0-9]+ Goal\s*$", phase_text)) != 1:
        raise SystemExit("ERROR: CURRENT_PHASE_CONTEXT must contain exactly one Phase <name> Goal heading")
    phase_heading = re.search(r"(?m)^## (Phase [A-Za-z0-9]+) Goal\s*$", phase_text).group(1)
    if f"## {phase_heading} Next Planned Feature" not in phase_text:
        raise SystemExit(f"ERROR: CURRENT_PHASE_CONTEXT missing {phase_heading} next-feature section")
    baseline_text = targets["CURRENT_ARCHITECTURE_BASELINE.md"].read_text()
    if not any(
        anchor in baseline_text
        for anchor in ("## Implemented Phase 2R Features", "## Phase 2R Next Planned Feature", "## Approved Consolidated")
    ):
        raise SystemExit("ERROR: architecture baseline has no supported closeout insertion anchor")
    if "## Accepted Decisions\n" not in targets["CURRENT_DECISION_SUMMARY.md"].read_text():
        raise SystemExit("ERROR: decision summary missing accepted-decisions section")
    _, trace_columns = table_row(targets["FEATURE_TRACEABILITY_MATRIX.md"].read_text(), feature)
    if len(trace_columns) not in (6, 8):
        raise SystemExit("ERROR: unsupported traceability table shape; expected 6 or 8 columns")
    # The VS-000 matrix is an authoritative contract registry, not a feature
    # lifecycle table. It is a required read-only evidence target at closeout.
    vs000_text = targets["VS-000_CONTRACT_TRACEABILITY_MATRIX.md"].read_text()
    if "## Schema Mapping" not in vs000_text or feature not in vs000_text:
        raise SystemExit("ERROR: VS-000 traceability matrix lacks FEATURE-0018 evidence")


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    parser.add_argument("--pr", required=True, type=int)
    parser.add_argument("--apply", action="store_true")
    args = parser.parse_args()
    subprocess.run(
        [str(ROOT / "scripts/feature-control.py"), "validate", "--feature", args.feature],
        cwd=ROOT,
        check=True,
    )
    control = json.loads((ROOT / f".automation/features/{args.feature}.control.json").read_text())
    control_path = ROOT / f".automation/features/{args.feature}.control.json"
    state_path = ROOT / f".automation/state/{args.feature}.json"
    state = json.loads(state_path.read_text())
    if state.get("control_approved_sha256") != hashlib.sha256(control_path.read_bytes()).hexdigest():
        raise SystemExit("ERROR: feature control differs from the founder-approved digest")
    if state.get("final_approval_token") != "APPROVED_FOR_MERGE" or not state.get("final_approved_commit"):
        raise SystemExit("ERROR: closeout requires the explicit final human approval")
    pr = gh_pr(args.pr)
    feature = control["feature"]
    if pr.get("state") != "MERGED" or not pr.get("mergedAt") or not pr.get("mergeCommit", {}).get("oid"):
        raise SystemExit("ERROR: PR is not merged with a recorded merge commit")
    if pr["baseRefName"] != feature["phase_branch"] or pr["headRefName"] != feature["feature_branch"]:
        raise SystemExit("ERROR: PR branch identity does not match feature control")
    branch = git("branch", "--show-current")
    if branch != feature["phase_branch"]:
        raise SystemExit(f"ERROR: closeout must run on {feature['phase_branch']}, current={branch}")
    if git("status", "--porcelain"):
        raise SystemExit("ERROR: closeout requires a clean working tree")
    commit = pr["mergeCommit"]["oid"]
    if subprocess.run(["git", "merge-base", "--is-ancestor", commit, "HEAD"], cwd=ROOT).returncode:
        raise SystemExit("ERROR: merged PR commit is not an ancestor of the current branch")
    if subprocess.run(
        ["git", "merge-base", "--is-ancestor", state["final_approved_commit"], "HEAD"], cwd=ROOT
    ).returncode:
        raise SystemExit("ERROR: final human-approved commit is not present on the phase branch")
    merged_at = datetime.fromisoformat(pr["mergedAt"].replace("Z", "+00:00"))
    date = merged_at.date().isoformat()
    short = commit[:7]
    evidence = f"PR #{pr['number']} merged {date} as commit `{short}` into `{feature['phase_branch']}`"
    targets = {Path(item).name: ROOT / item for item in control["closeout"]["metadata_targets"]}
    plan = {
        "feature": args.feature,
        "pr": pr["number"],
        "merged_at": date,
        "merge_commit": commit,
        "phase_branch": feature["phase_branch"],
        "mode": "prepare_only",
        "will_commit": False,
        "will_push": False,
        "targets": [feature["feature_file"], f".automation/state/{args.feature}.json", *control["closeout"]["metadata_targets"]],
    }
    print(json.dumps(plan, indent=2, sort_keys=True))
    if not args.apply:
        print("DRY RUN: rerun with --apply to prepare metadata edits; commit and push remain manual.")
        return
    feature_path = ROOT / feature["feature_file"]
    preflight_closeout_targets(feature_path, targets, args.feature)
    update_frontmatter(
        feature_path,
        {
            "status": "implemented_and_merged",
            "ai_summary": f"Implemented and merged through PR #{pr['number']} as commit {short} on {date}; approved architecture and handoffs remain controlling.",
            "merged_pr": f"#{pr['number']}",
            "merged_at": date,
            "merged_commit": commit,
        },
    )
    update_feature_index(targets["FEATURE_INDEX.md"], args.feature, evidence)
    update_phase_context(targets["CURRENT_PHASE_CONTEXT.md"], control, evidence, date)
    update_architecture_baseline(targets["CURRENT_ARCHITECTURE_BASELINE.md"], control, evidence)
    update_decision_summary(targets["CURRENT_DECISION_SUMMARY.md"], control, evidence)
    update_traceability(targets["FEATURE_TRACEABILITY_MATRIX.md"], args.feature, f"Final feature gate passed; {evidence}")
    state.update(
        {
            "current_stage": "merged",
            "current_task": "",
            "status": "implemented_and_merged",
            "merged_at": date,
            "merged_commit": commit,
            "merged_into": feature["phase_branch"],
            "merged_pr": f"#{pr['number']}",
            "updated_at": datetime.now(timezone.utc).isoformat(),
        }
    )
    state_path.write_text(json.dumps(state, indent=2, sort_keys=True) + "\n")
    subprocess.run(["git", "diff", "--check"], cwd=ROOT, check=True)
    print("Prepared closeout metadata. No commit or push was performed.")
    print("Review with: git diff --check && git diff")
    print("Commit manually after review.")


if __name__ == "__main__":
    main()
