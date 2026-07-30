#!/usr/bin/env python3
"""Build the compact founder review packet for requirements plus design."""

from __future__ import annotations

import argparse
import hashlib
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    args = parser.parse_args()
    state = json.loads((ROOT / f".automation/state/{args.feature}.json").read_text())
    control = json.loads((ROOT / f".automation/features/{args.feature}.control.json").read_text())
    spec = ROOT / state["spec_path"]
    requirements, design = spec / "requirements.md", spec / "design.md"
    review_dir = ROOT / f".automation/reviews/{args.feature}"
    reviews = {}
    deltas = {}
    for stage in ("requirements", "design"):
        reviews[stage] = json.loads((review_dir / f"{stage}.review.json").read_text())
        delta_path = review_dir / f"{stage}.semantic-delta.json"
        deltas[stage] = json.loads(delta_path.read_text()) if delta_path.is_file() else {"classification": "UNAVAILABLE"}
        if reviews[stage].get("status") != "APPROVED":
            raise SystemExit(f"ERROR: {stage} review is not approved")
    packet = {
        "feature": args.feature,
        "human_gate": "executable_plan",
        "requirements_sha256": sha256(requirements),
        "design_sha256": sha256(design),
        "machine_reviews": reviews,
        "semantic_deltas": deltas,
        "owned_resources": control["ownership"]["owned_resources"],
        "previous_features": control["ownership"]["previous_features"],
        "excluded_features": control["ownership"]["excluded_features"],
        "founder_checklist": [
            "Requirements state observable behavior without implementation choices.",
            "Design resolves mechanics without adding product semantics.",
            "Previous-feature contracts remain canonically owned and reused.",
            "Adjacent-feature and future-phase behavior remains excluded.",
            "Security, compatibility, observability, testing, and versioning are sufficient.",
            "No unresolved decision is delegated to tasks or Cursor.",
        ],
    }
    out_dir = ROOT / f".automation/reports/{args.feature}"
    out_dir.mkdir(parents=True, exist_ok=True)
    json_path = out_dir / "executable-plan-review.json"
    md_path = out_dir / "executable-plan-review.md"
    json_path.write_text(json.dumps(packet, indent=2, sort_keys=True) + "\n")
    lines = [
        f"# {args.feature} Executable Plan Review",
        "",
        f"- Requirements digest: `{packet['requirements_sha256']}`",
        f"- Design digest: `{packet['design_sha256']}`",
        f"- Requirements review: `{reviews['requirements']['status']}` — {reviews['requirements']['summary']}",
        f"- Design review: `{reviews['design']['status']}` — {reviews['design']['summary']}",
        f"- Requirements delta: `{deltas['requirements']['classification']}`",
        f"- Design delta: `{deltas['design']['classification']}`",
        "",
        "## Owned resources",
        "",
        *(f"- `{item}`" for item in packet["owned_resources"]),
        "",
        "## Previous-feature dependencies",
        "",
        *(f"- `{item['id']}` — {item['relationship']}" for item in packet["previous_features"]),
        "",
        "## Excluded adjacent features",
        "",
        *(f"- `{item['id']}` — {', '.join(item['concepts'])}" for item in packet["excluded_features"]),
        "",
        "## Founder checklist",
        "",
        *(f"- [ ] {item}" for item in packet["founder_checklist"]),
        "",
    ]
    md_path.write_text("\n".join(lines))
    print(md_path.relative_to(ROOT))


if __name__ == "__main__":
    main()
