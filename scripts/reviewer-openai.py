#!/usr/bin/env python3
"""Optional OpenAI reviewer adapter for Sovrunn Feature Factory.

Reads a rendered reviewer prompt and writes strict review JSON.
Requires OPENAI_API_KEY. Model can be overridden with FEATURE_FACTORY_REVIEWER_MODEL.
"""
import argparse
import json
import os
import re
import sys
import urllib.error
import urllib.request
from pathlib import Path

SCHEMA = {
    "type": "object",
    "additionalProperties": False,
    "required": [
        "status",
        "approval_token",
        "next_stage",
        "summary",
        "blocking_issues",
        "required_changes",
        "revision_prompt",
    ],
    "properties": {
        "status": {"type": "string", "enum": ["APPROVED", "NEEDS_REVISION", "BLOCKED"]},
        "approval_token": {
            "type": "string",
            "enum": ["APPROVED_FOR_DESIGN", "APPROVED_FOR_TASKS", "APPROVED_FOR_CURSOR", "NONE"],
        },
        "next_stage": {"type": "string", "enum": ["design", "tasks", "cursor", "none"]},
        "summary": {"type": "string"},
        "blocking_issues": {"type": "array", "items": {"type": "string"}},
        "required_changes": {"type": "array", "items": {"type": "string"}},
        "revision_prompt": {"type": "string"},
    },
}

EXPECTED_BY_STAGE = {
    "requirements": ("APPROVED_FOR_DESIGN", "design"),
    "design": ("APPROVED_FOR_TASKS", "tasks"),
    "tasks": ("APPROVED_FOR_CURSOR", "cursor"),
}


def extract_stage(prompt: str) -> str:
    m = re.search(r"(?m)^Stage:\s*\n\s*([a-zA-Z0-9_-]+)\s*$", prompt)
    if not m:
        return ""
    return m.group(1).strip()


def response_text(payload: dict) -> str:
    if isinstance(payload.get("output_text"), str):
        return payload["output_text"]
    parts = []
    for item in payload.get("output", []) or []:
        for content in item.get("content", []) or []:
            text = content.get("text")
            if isinstance(text, str):
                parts.append(text)
    return "\n".join(parts).strip()


def validate_review(review: dict, stage: str) -> None:
    required = set(SCHEMA["required"])
    missing = sorted(required - set(review))
    if missing:
        raise SystemExit(f"ERROR: review JSON missing required fields: {', '.join(missing)}")
    status = review.get("status")
    token = review.get("approval_token")
    next_stage = review.get("next_stage")
    if status not in {"APPROVED", "NEEDS_REVISION", "BLOCKED"}:
        raise SystemExit(f"ERROR: invalid status: {status}")
    if token not in {"APPROVED_FOR_DESIGN", "APPROVED_FOR_TASKS", "APPROVED_FOR_CURSOR", "NONE"}:
        raise SystemExit(f"ERROR: invalid approval_token: {token}")
    if next_stage not in {"design", "tasks", "cursor", "none"}:
        raise SystemExit(f"ERROR: invalid next_stage: {next_stage}")
    if not isinstance(review.get("blocking_issues"), list):
        raise SystemExit("ERROR: blocking_issues must be an array")
    if not isinstance(review.get("required_changes"), list):
        raise SystemExit("ERROR: required_changes must be an array")
    if status == "APPROVED":
        expected = EXPECTED_BY_STAGE.get(stage)
        if expected and (token, next_stage) != expected:
            raise SystemExit(
                f"ERROR: APPROVED review has token/next_stage {(token, next_stage)}, expected {expected} for stage {stage}"
            )
        if review.get("blocking_issues") or review.get("required_changes"):
            raise SystemExit("ERROR: APPROVED review must not include blocking_issues or required_changes")
    else:
        if token != "NONE" or next_stage != "none":
            raise SystemExit("ERROR: non-approved review must use approval_token=NONE and next_stage=none")


def call_openai(prompt: str) -> dict:
    api_key = os.environ.get("OPENAI_API_KEY")
    if not api_key:
        raise SystemExit("ERROR: OPENAI_API_KEY is required for reviewer-openai.py")
    model = os.environ.get("FEATURE_FACTORY_REVIEWER_MODEL", "gpt-5")
    body = {
        "model": model,
        "store": False,
        "input": [
            {
                "role": "system",
                "content": [
                    {
                        "type": "input_text",
                        "text": "You are a strict Sovrunn Feature Factory reviewer. Return only valid JSON matching the schema.",
                    }
                ],
            },
            {"role": "user", "content": [{"type": "input_text", "text": prompt}]},
        ],
        "text": {
            "format": {
                "type": "json_schema",
                "name": "sovrunn_feature_review",
                "strict": True,
                "schema": SCHEMA,
            }
        },
    }
    req = urllib.request.Request(
        "https://api.openai.com/v1/responses",
        data=json.dumps(body).encode("utf-8"),
        headers={
            "Authorization": f"Bearer {api_key}",
            "Content-Type": "application/json",
        },
        method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=int(os.environ.get("FEATURE_FACTORY_REVIEW_TIMEOUT", "240"))) as resp:
            return json.loads(resp.read().decode("utf-8"))
    except urllib.error.HTTPError as e:
        detail = e.read().decode("utf-8", errors="replace")
        raise SystemExit(f"ERROR: OpenAI reviewer request failed: HTTP {e.code}: {detail}") from e
    except urllib.error.URLError as e:
        raise SystemExit(f"ERROR: OpenAI reviewer request failed: {e}") from e


def main() -> None:
    p = argparse.ArgumentParser()
    p.add_argument("--prompt", required=True)
    p.add_argument("--out", required=True)
    p.add_argument("--raw-out", default="")
    args = p.parse_args()

    prompt = Path(args.prompt).read_text()
    stage = extract_stage(prompt)
    if stage not in EXPECTED_BY_STAGE:
        raise SystemExit(f"ERROR: could not determine valid stage from prompt: {stage!r}")

    raw = call_openai(prompt)
    if args.raw_out:
        Path(args.raw_out).parent.mkdir(parents=True, exist_ok=True)
        Path(args.raw_out).write_text(json.dumps(raw, indent=2, sort_keys=True) + "\n")

    text = response_text(raw)
    if not text:
        raise SystemExit("ERROR: reviewer response did not contain output text")
    try:
        review = json.loads(text)
    except json.JSONDecodeError as e:
        raise SystemExit(f"ERROR: reviewer output was not valid JSON: {e}: {text[:1000]}") from e
    validate_review(review, stage)

    out = Path(args.out)
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(review, indent=2, sort_keys=True) + "\n")
    print(out)


if __name__ == "__main__":
    main()
