#!/usr/bin/env python3
"""Build the hash-bound FEATURE-0018 semantic design-review TARGET projection.

Authorized by ADH-2026-076 (FEATURE-0018 semantic-review authority realignment),
clause 3 ("Hash-bound semantic review package"). This is the missed routing
implementation: the FEATURE-0018 design review must render a controlled semantic
projection of raw ``design.md`` as the reviewer target instead of the ~4,000-line
raw design, while ``semantic-delta.py`` continues to hash the raw ``design.md``.

The projection retains the semantic architecture, ownership, routes, validation,
security, traceability, non-goals and classification content of the raw design and
omits only the contract-governed executable-mechanics restatements (which the
hash-bound Design Mechanics Contract is the sole authority for, per ADH-2026-075)
and the non-normative revision history.

The projection is bound to three deterministic proofs so the reviewer cannot be
handed a stale or unverified target:

1. the raw ``design.md`` SHA-256 (the exact bytes ``semantic-delta.py`` also hashes);
2. the Design Mechanics Contract SHA-256; and
3. a successful ``make feature-0018-design-contract-check`` receipt (the checker
   verifies internal consistency and the contract's own hash binding to design.md).

If the deterministic contract check does not pass, or either source hash cannot be
read, the build fails closed and produces no projection. This keeps the raw design
mechanics from ever overriding the contract at review time.

Standard library only. Deterministic (fixed retained ranges, stable ordering, no
network access).
"""

from __future__ import annotations

import argparse
import hashlib
import json
import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SPEC_DIR = ROOT / ".kiro/specs/governance-iam-approval-exception-foundation"
DESIGN_PATH = SPEC_DIR / "design.md"
CONTRACT_PATH = SPEC_DIR / "design-mechanics-contract.json"
CONTRACT_CHECK = ROOT / "scripts/feature0018-design-contract-check.py"

OUTPUT_PATH = ROOT / ".automation/context-projections/FEATURE-0018.design-review-target.md"
MANIFEST_PATH = ROOT / ".automation/context-projections/FEATURE-0018.design-review-target.manifest.json"


def checkpoint_active() -> bool:
    """True when the FEATURE-0018 implementation checkpoint is active.

    ADH-2026-077 §4/§5: after the checkpoint, the hash-bound design-review target
    projection is no longer a design-review authority. The frozen design.md
    baseline governs and no design review runs, so this builder is inert.
    """
    control_path = ROOT / ".automation/features/FEATURE-0018.control.json"
    if not control_path.is_file():
        return False
    try:
        control = json.loads(control_path.read_text())
    except (OSError, json.JSONDecodeError):
        return False
    checkpoint = control.get("checkpoint")
    return isinstance(checkpoint, dict) and checkpoint.get("active") is True


# Retained SEMANTIC line ranges (1-based, inclusive) from raw design.md.
#
# Each range is one contiguous span of semantic architecture / ownership / routes /
# validation / security / traceability / non-goals / classification content. The
# gaps between the ranges are the OMITTED contract-governed executable-mechanics
# restatements and non-normative history:
#
#   576..903   DD-13 + DD-14 (contract-governed static-checker/publication mechanics)
#   1443..2590 §5.4 In-memory atomic publication protocol (contract-governed mechanics)
#   2831..3018 §7.2 Unit and property tests (contract-governed focused-mechanics tests)
#   3472..EOF  §10.4 Final self-verification + Model Execution Report (non-normative history)
#
# The retained ranges are validated against the raw design's own headings so the
# projection breaks loudly (fails closed) if design.md is restructured, rather than
# silently omitting or including the wrong content.
RETAINED_RANGES: list[dict[str, object]] = [
    {
        "label": "Front matter, section index, §1 identity/boundary, §2 decisions DD-01..DD-12",
        "start": 1,
        "end": 575,
    },
    {
        "label": "§3 components/ownership/writers/topology, §4 data/API/route ledger, §5.1-§5.3 validation and error mapping",
        "start": 904,
        "end": 1442,
    },
    {
        "label": "§5.5 audited release, §6 security/privacy/observability/compat, FEATURE-0013 Adoption, §7.1 local conformance",
        "start": 2591,
        "end": 2830,
    },
    {
        "label": "§7.3 executable conformance gating, §8 classification, §9 traceability, §10.1-§10.3 non-goals/absence/unresolved",
        "start": 3019,
        "end": 3471,
    },
]

# Section-boundary anchors that must remain intact for the fixed ranges to be
# correct. Each entry is (1-based line, exact prefix). If design.md is restructured
# so that any anchor moves, the build fails closed rather than emit a mis-sliced
# projection.
BOUNDARY_ANCHORS: list[tuple[int, str]] = [
    (576, "### DD-13 "),
    (762, "### DD-14 "),
    (904, "## 3. Components and repository paths"),
    (1329, "## 5. Validation and deterministic error behavior"),
    (1443, "### 5.4 In-memory atomic publication protocol"),
    (2591, "### 5.5 Provisional authorization and post-stage-10 audited release"),
    (2831, "### 7.2 Unit and property tests"),
    (3019, "### 7.3 Executable conformance gating"),
    (3472, "### 10.4 Final self-verification"),
    (3717, "## Model Execution Report"),
]


class BuildError(Exception):
    """Raised for a fail-closed projection-build violation."""


def digest(payload: bytes) -> str:
    return hashlib.sha256(payload).hexdigest()


def run_contract_check() -> str:
    """Run the deterministic design-contract checker; fail closed if it does not pass."""
    proc = subprocess.run(
        [sys.executable, str(CONTRACT_CHECK)],
        cwd=str(ROOT),
        capture_output=True,
        text=True,
    )
    receipt = (proc.stdout or "").strip()
    if proc.returncode != 0 or "PASS:" not in receipt:
        raise BuildError(
            "design-mechanics contract check did not pass; refusing to build the "
            "review target projection (raw design mechanics must never override the "
            f"contract). rc={proc.returncode} stdout={receipt!r} "
            f"stderr={(proc.stderr or '').strip()!r}"
        )
    # Record the single authoritative PASS line as the receipt.
    for line in receipt.splitlines():
        if line.startswith("PASS:"):
            return line.strip()
    return receipt


def verify_boundaries(lines: list[str]) -> None:
    for line_no, prefix in BOUNDARY_ANCHORS:
        if line_no > len(lines):
            raise BuildError(
                f"design.md is shorter than expected; missing anchor at line {line_no} "
                f"({prefix!r}). Refusing to emit a mis-sliced projection."
            )
        actual = lines[line_no - 1]
        if not actual.startswith(prefix):
            raise BuildError(
                f"design.md structure drifted: expected line {line_no} to start with "
                f"{prefix!r} but found {actual!r}. Retained ranges are no longer valid; "
                "update RETAINED_RANGES/BOUNDARY_ANCHORS under ADH-2026-076 before rebuilding."
            )


def build(design_bytes: bytes, contract_bytes: bytes, receipt: str) -> bytes:
    design_text = design_bytes.decode("utf-8")
    lines = design_text.splitlines()
    verify_boundaries(lines)

    raw_design_sha = digest(design_bytes)
    contract_sha = digest(contract_bytes)

    header = [
        "---",
        "doc_type: ai_design_review_target_projection",
        "feature: FEATURE-0018",
        "stage: design-review-target",
        "authority: non-authoritative-semantic-projection-of-raw-design",
        "controlling_handoff: ADH-2026-076",
        "generated: true",
        "---",
        "",
        "# FEATURE-0018 Design-Review Target Projection (semantic, hash-bound)",
        "",
        "This is a non-authoritative, mechanically generated semantic projection of",
        "the raw `design.md` submitted for review, produced under ADH-2026-076 clause 3 as the",
        "controlled FEATURE-0018 design-review **target**. It retains the semantic",
        "architecture, ownership, routes, validation, security, traceability, non-goals",
        "and classification content and omits only the contract-governed executable",
        "mechanics (for which the hash-bound Design Mechanics Contract is the sole",
        "authority under ADH-2026-075) and the non-normative revision history.",
        "",
        "The raw `design.md` and the Design Mechanics Contract remain authoritative.",
        "This projection is bound to their exact bytes and to a successful deterministic",
        "contract-check receipt; a source-hash mismatch or a failing contract check",
        "blocks regeneration, so the raw design mechanics can never override the contract",
        "at review time.",
        "",
        "## Hash-bound review-target binding",
        "",
        f"- Raw design source: `{DESIGN_PATH.relative_to(ROOT)}`",
        f"- Raw design SHA-256 (also hashed by `semantic-delta.py`): `{raw_design_sha}`",
        f"- Design Mechanics Contract source: `{CONTRACT_PATH.relative_to(ROOT)}`",
        f"- Design Mechanics Contract SHA-256: `{contract_sha}`",
        f"- Deterministic contract-check receipt (`make feature-0018-design-contract-check`): `{receipt}`",
        "",
        "## Omitted from this semantic target (authoritative elsewhere)",
        "",
        "- Contract-governed executable mechanics — DD-13, DD-14, §5.4 in-memory atomic",
        "  publication protocol, and §7.2 focused-mechanics unit/property tests — are the",
        "  Design Mechanics Contract's authority (ADH-2026-075) and are supplied to the",
        "  reviewer as the hash-bound contract in the controlled design-review context,",
        "  not restated here.",
        "- Non-normative history — §10.4 Final self-verification and the Model Execution",
        "  Report — is omitted; it asserts no product semantics.",
        "",
        "---",
        "",
        "## Retained semantic design content (exact excerpts from raw `design.md`)",
        "",
    ]

    body: list[str] = []
    for rng in RETAINED_RANGES:
        start = int(rng["start"])
        end = int(rng["end"])
        excerpt = "\n".join(lines[start - 1 : end])
        body.extend(
            [
                f"<!-- BEGIN RETAINED SEMANTIC EXCERPT: raw design.md lines {start}-{end} "
                f"({rng['label']}) -->",
                excerpt,
                f"<!-- END RETAINED SEMANTIC EXCERPT: raw design.md lines {start}-{end} -->",
                "",
            ]
        )

    rendered = "\n".join(header + body).rstrip("\n") + "\n"
    return rendered.encode("utf-8")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--design", default=str(DESIGN_PATH))
    parser.add_argument("--contract", default=str(CONTRACT_PATH))
    parser.add_argument("--output", default=str(OUTPUT_PATH))
    parser.add_argument("--manifest", default=str(MANIFEST_PATH))
    parser.add_argument(
        "--check-only",
        action="store_true",
        help="verify an existing projection is current without rewriting it",
    )
    args = parser.parse_args()

    if checkpoint_active():
        print(
            "SKIP: FEATURE-0018 implementation checkpoint active (ADH-2026-077); "
            "the design-review target projection is not a design-review authority "
            "and is not built. The frozen design.md baseline governs."
        )
        return 0

    design_path = Path(args.design)
    contract_path = Path(args.contract)
    output_path = Path(args.output)
    manifest_path = Path(args.manifest)

    try:
        if not design_path.is_file():
            raise BuildError(f"design.md not found: {design_path}")
        if not contract_path.is_file():
            raise BuildError(f"design-mechanics contract not found: {contract_path}")

        design_bytes = design_path.read_bytes()
        contract_bytes = contract_path.read_bytes()
        receipt = run_contract_check()
        payload = build(design_bytes, contract_bytes, receipt)

        manifest = {
            "schema_version": "1.0",
            "feature": "FEATURE-0018",
            "stage": "design-review-target",
            "controlling_handoff": "ADH-2026-076",
            "projection": str(output_path.relative_to(ROOT)),
            "projection_bytes": len(payload),
            "projection_estimated_tokens": (len(payload) + 3) // 4,
            "projection_sha256": digest(payload),
            "binding": {
                "raw_design": str(design_path.relative_to(ROOT)),
                "raw_design_sha256": digest(design_bytes),
                "mechanics_contract": str(contract_path.relative_to(ROOT)),
                "mechanics_contract_sha256": digest(contract_bytes),
                "contract_check_receipt": receipt,
            },
            "retained_ranges": [
                {"label": r["label"], "start": r["start"], "end": r["end"]}
                for r in RETAINED_RANGES
            ],
            "omitted": [
                "DD-13 (contract-governed static-checker mechanics)",
                "DD-14 (contract-governed publication-instant mechanics)",
                "§5.4 in-memory atomic publication protocol (contract-governed mechanics)",
                "§7.2 focused-mechanics unit/property tests (contract-governed mechanics)",
                "§10.4 Final self-verification (non-normative history)",
                "Model Execution Report (non-normative history)",
            ],
        }
        manifest_bytes = (json.dumps(manifest, indent=2, sort_keys=True) + "\n").encode("utf-8")
    except BuildError as exc:
        print(f"FAIL: {exc}", file=sys.stderr)
        return 1

    if args.check_only:
        if not output_path.is_file() or output_path.read_bytes() != payload:
            print(
                "FAIL: FEATURE-0018 design-review target projection is stale or missing; "
                "run scripts/feature0018-design-review-target.py to regenerate",
                file=sys.stderr,
            )
            return 1
        if not manifest_path.is_file() or manifest_path.read_bytes() != manifest_bytes:
            print(
                "FAIL: FEATURE-0018 design-review target manifest is stale or missing; "
                "run scripts/feature0018-design-review-target.py to regenerate",
                file=sys.stderr,
            )
            return 1
        print(
            "PASS: FEATURE-0018 design-review target projection is current "
            f"(sha256={manifest['projection_sha256']}, bytes={manifest['projection_bytes']})"
        )
        return 0

    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_bytes(payload)
    manifest_path.parent.mkdir(parents=True, exist_ok=True)
    manifest_path.write_bytes(manifest_bytes)
    print(f"projection={output_path.relative_to(ROOT)}")
    print(f"manifest={manifest_path.relative_to(ROOT)}")
    print(f"sha256={manifest['projection_sha256']}")
    print(f"bytes={manifest['projection_bytes']}")
    print(f"raw_design_sha256={manifest['binding']['raw_design_sha256']}")
    print(f"mechanics_contract_sha256={manifest['binding']['mechanics_contract_sha256']}")
    print(f"contract_check_receipt={receipt}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
