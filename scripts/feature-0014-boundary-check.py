#!/usr/bin/env python3
"""Fail-closed architecture/readiness checks for FEATURE-0014."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
import sys
from pathlib import Path


FEATURE = "FEATURE-0014"
SLUG = "provider-neutral-resource-model"
ARCH = Path("docs/architecture/provider-neutral-resource-model.md")
ADH = Path("docs/reviews/architecture-decision-handoffs/ADH-2026-018-feature-0014-provider-neutral-resource-model.md")
SPEC_DIR = Path(".kiro/specs") / SLUG
OLD_TERMS = ("IaaSStack", "ProviderRegion", "ProviderLocation/Region", "ProviderLocation / ProviderRegion")
DECISIONS = {f"F14-AD-{i:03d}" for i in range(1, 22)}
RISKS = {f"F14-R{i:02d}" for i in range(1, 31)}
REQUIREMENTS = {f"F14-REQ-{i:02d}" for i in range(1, 32)}
KINDS = {
    "Provider",
    "ProviderLocation",
    "ProviderDatacenter",
    "DatacenterFailureDomain",
    "InfrastructureStack",
}

DESIGN_CONTEXT = {
    "AGENTS.md",
    "README.md",
    "docs/engineering/ai-context-loading-standard.md",
    "docs/foundation/constitution.md",
    "docs/decisions/DECISION_INDEX.md",
    "docs/glossary.md",
    "docs/features/FEATURE_SEQUENCE.md",
    "docs/resource-specs/RESOURCE_MODEL_PHASE1.md",
    "docs/api/API_CONTRACT_PHASE1.md",
    ".kiro/steering/product.md",
    ".kiro/steering/architecture.md",
    ".kiro/steering/engineering.md",
    "docs/engineering/go-coding-guardrails.md",
    "docs/engineering/go-version-standard.md",
    "docs/phase2/PHASE2_SCOPE.md",
    "docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md",
    "docs/architecture/api-resource-standard.md",
    str(ARCH),
    str(ADH),
    "docs/features/FEATURE-0014-provider-neutral-resource-model.md",
    str(SPEC_DIR / "requirements.md"),
}

ACTIVE_CONTEXT = [
    ARCH,
    ADH,
    Path("docs/rfc/RFC-0024-provider-neutral-resource-model.md"),
    Path("docs/glossary.md"),
    Path("docs/features/FEATURE_INDEX.md"),
    Path("docs/phase2/PHASE2_FEATURE_SEQUENCE.md"),
    Path("docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md"),
    Path("docs/architecture/development-phases.md"),
    Path("docs/context/CURRENT_PHASE_CONTEXT.md"),
    Path("docs/context/CURRENT_ARCHITECTURE_BASELINE.md"),
    Path("docs/context/CURRENT_DECISION_SUMMARY.md"),
]

REQUIRED_CONTEXT = [
    Path("AGENTS.md"),
    Path("docs/engineering/ai-context-loading-standard.md"),
    Path("docs/foundation/constitution.md"),
    Path("docs/decisions/DECISION_INDEX.md"),
    Path("docs/phase2/PHASE2_EXECUTION_STRATEGY.md"),
    Path("docs/phase2/PHASE2_ARCHITECTURE_SPINE.md"),
    Path("docs/phase2/PHASE2_SCOPE.md"),
    Path("docs/phase2/PHASE2_FEATURE_SEQUENCE.md"),
    Path("docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md"),
    Path("docs/architecture/api-resource-standard.md"),
    Path("docs/reviews/architecture-decision-handoffs/ADH-2026-012-feature-0012-api-resource-standard.md"),
    Path("docs/reviews/architecture-decision-handoffs/ADH-2026-013-operation-allowed-scopes.md"),
    Path("docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md"),
    Path("docs/reviews/architecture-decision-handoffs/ADH-2026-017-feature-0013-consolidated-architecture.md"),
    ARCH,
    ADH,
    Path("docs/features/FEATURE-0014-provider-neutral-resource-model.md"),
    Path("docs/features/FEATURE_INDEX.md"),
    Path("docs/context/CURRENT_ARCHITECTURE_BASELINE.md"),
    Path("docs/context/CURRENT_DECISION_SUMMARY.md"),
    Path("docs/glossary.md"),
]


class Check:
    def __init__(self) -> None:
        self.errors: list[str] = []
        self.passes: list[str] = []

    def require(self, condition: bool, message: str) -> None:
        (self.passes if condition else self.errors).append(message)

    def require_file(self, path: Path) -> None:
        self.require(path.is_file(), f"required file: {path}")

    def finish(self) -> None:
        for message in self.passes:
            print(f"PASS: {message}")
        if self.errors:
            for message in self.errors:
                print(f"FAIL: {message}")
            print("REPOSITORY_CONTEXT_NOT_READY")
            raise SystemExit(1)
        print("PASS: FEATURE-0014 boundary/readiness checks")


def read(path: Path) -> str:
    return path.read_text(encoding="utf-8") if path.is_file() else ""


def ids(text: str, pattern: str) -> set[str]:
    return set(re.findall(pattern, text))


def check_repository(c: Check, stage: str) -> None:
    for path in REQUIRED_CONTEXT:
        c.require_file(path)

    arch = read(ARCH)
    adh = read(ADH)
    c.require("status: approved-pending-repository-alignment" in arch or "status: approved-for-kiro-requirements" in arch,
              "architecture has an approved readiness status")
    c.require("- Approval status: Approved" in adh, "ADH-2026-018 is human-approved")
    c.require(ids(arch, r"F14-AD-\d{3}") == DECISIONS, "architecture enumerates exactly F14-AD-001..021")
    c.require(ids(arch, r"F14-R\d{2}") == RISKS, "architecture enumerates exactly F14-R01..30")
    c.require(all(kind in arch for kind in KINDS), "architecture contains exactly the approved kind vocabulary")
    c.require("Applicability: `NOT_APPLICABLE`" in arch, "FEATURE-0013 adoption is NOT_APPLICABLE")

    for path in ACTIVE_CONTEXT:
        text = read(path)
        for term in OLD_TERMS:
            c.require(term not in text, f"{path} excludes superseded term {term!r}")

    index = read(Path("docs/features/FEATURE_INDEX.md"))
    c.require(f"| {FEATURE} | Provider-Neutral Resource Model" in index and f"`{SLUG}`" in index,
              "FEATURE_INDEX resolves the canonical FEATURE-0014 slug")
    c.require("ResourcePool" not in read(Path("docs/rfc/RFC-0024-provider-neutral-resource-model.md")).split("## Decision", 1)[0],
              "RFC summary does not absorb FEATURE-0015")

    prompt_checks = {
        Path("docs/prompts/kiro/requirements.prompt.md"): "## FEATURE-0014 closed architecture boundary",
        Path("docs/prompts/kiro/design.prompt.md"): "For FEATURE-0014",
        Path("docs/prompts/kiro/feature-0014-design.prompt.md"): "## Exact context boundary",
        Path("docs/prompts/kiro/tasks.prompt.md"): "For FEATURE-0014",
        Path("docs/prompts/reviewer/spec-review.prompt.md"): "For FEATURE-0014",
        Path("docs/prompts/reviewer/approval-review.prompt.md"): "For FEATURE-0014",
        Path("scripts/reviewer-stage.sh"): "docs/architecture/provider-neutral-resource-model.md",
        Path("scripts/kiro-stage.sh"): "feature-0014-boundary-check.py",
        Path("scripts/feature-gate.sh"): "feature-0014-boundary-check.py",
        Path("Makefile"): "feature-0014-architecture-readiness",
    }
    for path, marker in prompt_checks.items():
        c.require(marker in read(path), f"{path} contains FEATURE-0014 guardrail marker")

    state = Path(f".automation/state/{FEATURE}.json")
    assessment = Path(f".automation/features/{FEATURE}.yaml")
    config = SPEC_DIR / ".config.kiro"
    for path in (state, assessment, config):
        c.require_file(path)
    if state.is_file():
        data = json.loads(read(state))
        c.require(data.get("slug") == SLUG and data.get("feature_branch") == "feature-0014-provider-neutral-resource-model",
                  "FEATURE-0014 automation state resolves slug and branch")
        c.require(data.get("current_stage") == stage,
                  f"FEATURE-0014 automation state matches requested stage {stage}")
    if config.is_file():
        c.require(f"feature_id={FEATURE}" in read(config) and f"slug={SLUG}" in read(config),
                  "Kiro config resolves FEATURE-0014 identity")


def check_manifest(c: Check, path: Path, stage: str) -> None:
    c.require_file(path)
    if not path.is_file():
        return
    data = json.loads(read(path))
    manifest_paths = {item["path"] for item in data.get("files", [])}
    c.require(str(ARCH) in manifest_paths, f"{path} includes canonical architecture")
    c.require(str(ADH) in manifest_paths, f"{path} includes ADH-2026-018")
    for old in ("ADH-2026-014", "ADH-2026-015", "ADH-2026-016"):
        c.require(not any(old in item for item in manifest_paths), f"{path} excludes {old}")
    for item in data.get("files", []):
        candidate = Path(item["path"])
        digest = hashlib.sha256(candidate.read_bytes()).hexdigest() if candidate.is_file() else ""
        c.require(digest == item.get("sha256"), f"{path} hash matches {candidate}")
    if stage == "design" and "generated-prompts" in str(path):
        c.require(manifest_paths == DESIGN_CONTEXT,
                  f"{path} is the exact minimal FEATURE-0014 design context")
        c.require(not any("FEATURE-0013" in item or "ADH-2026-017" in item for item in manifest_paths),
                  f"{path} excludes FEATURE-0013 semantic payload")
        prompt = read(Path(f"docs/generated-prompts/{FEATURE}/design.prompt.md"))
        c.require(all(f"`{item}`" in prompt for item in DESIGN_CONTEXT),
                  "rendered design prompt lists every manifest input")
        c.require("existing implementations for similar resources" not in prompt.lower(),
                  "rendered design prompt has no open-ended implementation search")
        c.require("SecurityExceptionRef" not in prompt and "GraphEdge" not in prompt,
                  "rendered design prompt excludes FEATURE-0013 design instructions")


def check_design_authorization(c: Check) -> None:
    state_path = Path(f".automation/state/{FEATURE}.json")
    requirements_path = SPEC_DIR / "requirements.md"
    if not state_path.is_file() or not requirements_path.is_file():
        c.require(False, "design requires approval state and requirements")
        return
    state = json.loads(read(state_path))
    digest = hashlib.sha256(requirements_path.read_bytes()).hexdigest()
    c.require(state.get("requirements_approval_token") == "APPROVED_FOR_DESIGN",
              "design input has APPROVED_FOR_DESIGN token")
    c.require(state.get("requirements_approved_sha256") == digest,
              "approved requirements digest matches current design input")


def check_design_changed_files(c: Check) -> None:
    result = subprocess.run(["git", "status", "--porcelain"], text=True, capture_output=True, check=True)
    permitted = {
        f".automation/state/{FEATURE}.json",
        str(SPEC_DIR / "design.md"),
    }
    changed = set()
    for line in result.stdout.splitlines():
        path = line[3:].strip()
        if " -> " in path:
            path = path.split(" -> ", 1)[1]
        changed.add(path)
    c.require(changed <= permitted,
              "design stage changed only design.md and automation state"
              + (f": {sorted(changed - permitted)}" if changed - permitted else ""))


def check_stage(c: Check, stage: str, require_output: bool) -> None:
    path = SPEC_DIR / f"{stage}.md"
    if not path.is_file():
        c.require(not require_output, f"{path} may be absent before generation")
        return
    text = read(path)
    c.require(FEATURE in text, f"{path} contains active feature identity")
    c.require(f"Stage: {stage.title()}" in text, f"{path} contains stage label")
    for term in OLD_TERMS:
        c.require(term not in text, f"{path} excludes superseded term {term!r}")
    if stage == "requirements":
        c.require(ids(text, r"F14-AD-\d{3}") == DECISIONS, "requirements trace exactly all F14 decisions")
        c.require(ids(text, r"F14-R\d{2}") == RISKS, "requirements trace exactly all F14 risks")
        lower = text.lower()
        c.require("single-owner" in lower and "overlap ledger" in lower, "requirements contain single-owner overlap ledger")
        c.require("normalization ledger" in lower and "semantic key" in lower, "requirements contain normalization ledger")
        c.require("not_applicable" in lower, "requirements preserve FEATURE-0013 NOT_APPLICABLE")
    elif stage == "design":
        c.require((SPEC_DIR / "requirements.md").is_file(), "design has approved requirements input")
        c.require(ids(text, r"F14-REQ-\d{2}") == REQUIREMENTS,
                  "design maps exactly F14-REQ-01..31")
        c.require(ids(text, r"F14-AD-\d{3}") == DECISIONS, "design traces exactly all F14 decisions")
        c.require(ids(text, r"F14-R\d{2}") == RISKS, "design traces exactly all F14 risks")
        lower = text.lower()
        c.require("field" in lower and "writer" in lower and "mutability" in lower,
                  "design contains field ownership/writer/mutability ledger")
        c.require("requirement traceability" in lower and "orphan" in lower,
                  "design contains requirement traceability and orphan report")
        c.require("risk" in lower and "control" in lower and "evidence" in lower,
                  "design maps risk controls to evidence")
        c.require("not_applicable" in lower, "design preserves FEATURE-0013 NOT_APPLICABLE")
        c.require("absence ledger" in lower, "design contains adjacent-feature absence ledger")
        normative_paragraphs = []
        for paragraph in re.split(r"\n\s*\n", text):
            # Inline-code mentions explain normative keywords; they do not make
            # the surrounding design prose normative.
            prose = re.sub(r"`[^`\n]*`", "", paragraph)
            if re.search(r"\b(?:MUST|SHALL|MUST NOT|SHALL NOT)\b", prose):
                normative_paragraphs.append(paragraph)
        uncited = [paragraph.splitlines()[0][:100] for paragraph in normative_paragraphs
                   if not re.search(r"F14-REQ-\d{2}", paragraph)]
        c.require(not uncited, "design has no uncited normative paragraph"
                  + (f": {uncited}" if uncited else ""))
    elif stage == "tasks":
        c.require((SPEC_DIR / "requirements.md").is_file() and (SPEC_DIR / "design.md").is_file(),
                  "tasks have approved requirements and design inputs")
        c.require(ids(text, r"F14-AD-\d{3}") == DECISIONS, "tasks trace exactly all F14 decisions")
        c.require(ids(text, r"F14-R\d{2}") == RISKS, "tasks trace exactly all F14 risks")


def check_changed_files(c: Check) -> None:
    result = subprocess.run(["git", "status", "--porcelain"], text=True, capture_output=True, check=True)
    allowed = (
        ".automation/features/FEATURE-0014.yaml",
        ".automation/state/FEATURE-0014.json",
        ".kiro/specs/provider-neutral-resource-model/",
        "Makefile",
        "docs/",
        "scripts/feature-0014-boundary-check.py",
        "scripts/feature-gate.sh",
        "scripts/kiro-stage.sh",
        "scripts/reviewer-stage.sh",
        "scripts/render-prompt.py",
    )
    unexpected = []
    for line in result.stdout.splitlines():
        path = line[3:].strip()
        if " -> " in path:
            path = path.split(" -> ", 1)[1]
        if not any(path == item or (item.endswith("/") and path.startswith(item)) for item in allowed):
            unexpected.append(path)
    c.require(not unexpected, "working tree has no unrelated changes" + (f": {unexpected}" if unexpected else ""))


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--feature", default=FEATURE)
    parser.add_argument("--stage", choices=("requirements", "design", "tasks"), default="requirements")
    parser.add_argument("--mode", choices=("readiness", "pre", "prompt", "post", "review"), default="readiness")
    args = parser.parse_args()
    if args.feature != FEATURE:
        raise SystemExit(f"ERROR: validator supports only {FEATURE}")

    c = Check()
    check_repository(c, args.stage)
    check_changed_files(c)
    if args.stage == "design":
        check_design_authorization(c)
    if args.mode == "prompt":
        check_manifest(c, Path(f"docs/generated-prompts/{FEATURE}/{args.stage}.context.json"), args.stage)
    elif args.mode == "review":
        check_manifest(c, Path(f".automation/reviews/{FEATURE}/{args.stage}-review.context.json"), args.stage)
        check_stage(c, args.stage, True)
    elif args.mode == "post":
        check_stage(c, args.stage, True)
        if args.stage == "design":
            check_design_changed_files(c)
    elif args.mode == "pre":
        check_stage(c, args.stage, False)
    c.finish()


if __name__ == "__main__":
    main()
