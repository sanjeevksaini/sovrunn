#!/usr/bin/env python3
"""Fail closed when FEATURE-0013 downstream artifacts reopen architecture."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
from pathlib import Path


FEATURE = "FEATURE-0013"
ARCHITECTURE = Path(
    "docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md"
)
SPEC_DIR = Path(".kiro/specs/decision-object-and-auditevent-standard")
AGENT = Path(".kiro/agents/sovrunn-spec.md")
PROMPT_TEMPLATES = (
    Path("docs/prompts/kiro/requirements.prompt.md"),
    Path("docs/prompts/kiro/design.prompt.md"),
    Path("docs/prompts/kiro/tasks.prompt.md"),
)
STATE = Path(".automation/state/FEATURE-0013.json")
KIRO_STAGE = Path("scripts/kiro-stage.sh")
SPEC_FLOW = Path("scripts/spec-flow.sh")

ACTIVE_CONTEXT = (
    Path("docs/context/CURRENT_ARCHITECTURE_BASELINE.md"),
    Path("docs/context/CURRENT_DECISION_SUMMARY.md"),
    Path("docs/context/CURRENT_PHASE_CONTEXT.md"),
    Path("docs/features/FEATURE_INDEX.md"),
    Path("docs/phase2/PHASE2_ARCHITECTURE_SPINE.md"),
    Path("docs/phase2/PHASE2_ACCEPTANCE_GATES.md"),
    Path("docs/phase2/PHASE2_FEATURE_SEQUENCE.md"),
    Path("docs/rfc/RFC-0023-decision-and-audit-standard.md"),
    Path("docs/traceability/FEATURE_TRACEABILITY_MATRIX.md"),
)

DEPENDENCY_PATHS = (
    Path("docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md"),
    Path("docs/features/FEATURE-0011-reuse-assessment-standard.md"),
    Path("docs/reviews/feature-gates/FEATURE-0011-approval-review.md"),
    Path("docs/reviews/architecture-decision-handoffs/ADH-2026-011-feature-0011-reuse-assessment-standard.md"),
    Path("docs/architecture/api-resource-standard.md"),
    Path("docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md"),
    Path(".kiro/specs/api-resource-naming-status-and-validation-standard/requirements.md"),
    Path(".kiro/specs/api-resource-naming-status-and-validation-standard/design.md"),
    Path("docs/reviews/feature-gates/FEATURE-0012-approval-review.md"),
    Path("docs/reviews/architecture-decision-handoffs/ADH-2026-012-feature-0012-api-resource-standard.md"),
    Path("docs/reviews/architecture-decision-handoffs/ADH-2026-013-operation-allowed-scopes.md"),
    Path("api/schemas/_common/scope-ref.json"),
    Path("api/schemas/audit-event.json"),
    Path("internal/apimeta/scope.go"),
)

REQUIRED_DECISION_IDS = {f"AD-{index:03d}" for index in range(1, 46)}
DECISION_IMPLEMENTATION_CLASSES = {
    **{f"AD-{index:03d}": "CONTRACT_NOW" for index in range(1, 7)},
    "AD-007": "INVARIANT_FOR_LATER",
    "AD-008": "INVARIANT_FOR_LATER",
    "AD-009": "CONTRACT_NOW",
    "AD-010": "INVARIANT_FOR_LATER",
    "AD-011": "DEFERRED",
    "AD-012": "CONTRACT_NOW",
    "AD-013": "INVARIANT_FOR_LATER",
    "AD-014": "INVARIANT_FOR_LATER",
    **{f"AD-{index:03d}": "CONTRACT_NOW" for index in range(15, 35)},
    "AD-035": "INVARIANT_FOR_LATER",
    "AD-036": "CONTRACT_NOW",
    "AD-037": "CONTRACT_NOW",
    "AD-038": "CONTRACT_NOW",
    "AD-039": "INVARIANT_FOR_LATER",
    "AD-040": "CONTRACT_NOW",
    "AD-041": "INVARIANT_FOR_LATER",
    **{f"AD-{index:03d}": "CONTRACT_NOW" for index in range(42, 46)},
}
IMPLEMENTATION_CLASSES = {"CONTRACT_NOW", "INVARIANT_FOR_LATER", "DEFERRED"}
REQUIRED_F12_RISK_IDS = {f"F12-R{index:02d}" for index in range(1, 17)}
REQUIRED_RISK_IDS = {f"F13-R{index:02d}" for index in range(1, 32)}
REQUIRED_CONFORMANCE_IDS = {
    *(f"F13-CF-{index:02d}" for index in range(1, 29)),
    *(f"F13-SCOPE-{index:02d}" for index in range(1, 12)),
    *(f"F13-EVAL-{index:02d}" for index in range(1, 9)),
    *(f"F13-SEC-{index:02d}" for index in range(1, 9)),
    *(f"F13-TRUST-{index:02d}" for index in range(1, 5)),
    *(f"F13-COMPAT-{index:02d}" for index in range(1, 10)),
}
REQUIRED_ARCHITECTURE_SECTIONS = set(range(1, 29))

STAGE_FILES = {
    "requirements": SPEC_DIR / "requirements.md",
    "design": SPEC_DIR / "design.md",
    "tasks": SPEC_DIR / "tasks.md",
}

REQUIRED_AGENT_RESOURCES = (
    "docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md",
    "docs/architecture/api-resource-standard.md",
    "ADH-2026-017-feature-0013-consolidated-architecture.md",
    "FEATURE-0011-reuse-assessment-standard.md",
    "FEATURE-0011-approval-review.md",
    "PHASE2_REUSE_ASSESSMENT_STANDARD.md",
    "FEATURE-0012-api-resource-naming-status-and-validation-standard.md",
    "api-resource-naming-status-and-validation-standard/requirements.md",
    "api-resource-naming-status-and-validation-standard/design.md",
    "FEATURE-0012-approval-review.md",
    "ADH-2026-013-operation-allowed-scopes.md",
)

REQUIRED_ARCHITECTURE_CONTROLS = (
    "### 5.4 FEATURE-0012 resource-profile and boundary inheritance",
    "### 27.8 Closed architecture boundary for downstream stages",
    "### 27.9 Downstream design closure controls for clean regeneration",
    "SecurityExceptionRef",
    "FEATURE-0012 baseline workflow reuse",
    "AuditEvent package ownership",
    "Versioned registry keys",
    "TrustCarrier empty-state rule",
    "Matrix E automation boundary",
    "Scope pre-scan and decode reuse",
    "## 29. Architecture decision history",
    "FEATURE-0013 deliberately defines no new universal numeric decision-domain",
    "F13-SCOPE-01",
    "F13-EVAL-01",
    "F13-SEC-01",
    "F13-TRUST-01",
    "F13-COMPAT-01",
    "six-value FEATURE-0012 `ScopeKind` vocabulary",
    "typed subject under Project scope",
    "`DecisionRequest` is a conceptual input role, not a fifth FEATURE-0013 shared",
    "Transition to `Operation` is the sole accepted",
    "remain owned exclusively by FEATURE-0012",
    "The two namespaces must not be",
    "owned by FEATURE-0012/ADH-2026-013",
    "every AD-001 through AD-045 implementation class is fixed by section 19",
    "| ID | Architectural decision | Implementation class |",
    "A later production implementation is constrained by the following",
    "manifest is not cryptographically signed or verified by FEATURE-0013.",
    "bounded in-memory graph representation",
)

FORBIDDEN_GLOBAL = (
    (r"\bDecisionObject\b", "obsolete DecisionObject terminology"),
    (r"\bAuditScope\b", "parallel AuditScope model"),
    (r"(?:top-level|parallel|second)\s+scopeRef", "second scope authority"),
    (r"scopeRef\s+and\s+(?:a\s+)?(?:scope|scopeKind)", "dual scope authority"),
    (r"(?:seven[- ]value|seven canonical).{0,80}\bScopeKind\b", "superseded seven-value ScopeKind model"),
)

FORBIDDEN_DESIGN_QUESTION = (
    (r"static Go registry.*(?:or|versus).*declarative", "registry source of truth"),
    (r"per-profile.*(?:or|versus).*per-evaluator", "limit ownership"),
    (r"(?:whether|should).*scope.*(?:propagat|authorit|source)", "scope authority"),
    (r"(?:choose|select|decide).*(?:canonicali[sz]ation|digest|signature|cryptographic)", "deferred cryptographic selection"),
    (r"(?:choose|select|decide).*(?:resource profile|ImmutableRecord|VersionedDefinition)", "resource-profile assignment"),
    (r"(?:choose|select|decide).*(?:error family|Problem Details envelope)", "public error architecture"),
    (r"(?:choose|select|decide).*(?:DataClassification|sensitivity).*(?:map|mapping)", "classification mapping"),
    (r"(?:whether|should|choose|select|decide).*(?:shared|canonical|common).{0,40}DecisionRequest", "DecisionRequest ownership"),
    (r"(?:whether|should|choose|select|decide).*(?:pending|non-final).{0,40}(?:response|status|envelope)", "asynchronous response contract"),
)

REUSE_SUMMARY_HEADER = (
    "| Capability / decision unit | Disposition | Rationale | Decision status "
    "| Controlling reference |"
)
REUSE_DISPOSITIONS = {"Reuse", "Wrap", "Extend", "Build"}
REUSE_STATUSES = {"Proposed", "Approved", "Deferred", "Rejected", "Superseded"}

PUBLIC_VIOLATION_CODES = {
    "DECISION_PROFILE_UNKNOWN",
    "DECISION_PROFILE_VERSION_UNSUPPORTED",
    "DECISION_PROFILE_INACTIVE",
    "DECISION_PROFILE_SCHEMA_INVALID",
    "DECISION_PROFILE_LIMIT_INVALID",
    "DECISION_EVALUATION_TYPE_UNSUPPORTED",
    "DECISION_EVALUATION_RESULT_INVALID",
    "DECISION_EVALUATION_SCOPE_MISMATCH",
    "DECISION_EVALUATION_LIMIT_EXCEEDED",
    "DECISION_COMPOSITION_STRATEGY_UNSUPPORTED",
    "DECISION_COMPOSITION_GRAPH_INVALID",
    "DECISION_COMPOSITION_CYCLE",
    "DECISION_COMPOSITION_LIMIT_EXCEEDED",
    "DECISION_COMPOSITION_INPUT_CONFLICT",
    "DECISION_OBLIGATION_UNKNOWN",
    "DECISION_OBLIGATION_INVALID",
    "DECISION_OBLIGATION_UNSUPPORTED_MANDATORY",
    "DECISION_TRUST_REQUIRED",
    "DECISION_TRUST_UNKNOWN",
    "DECISION_TRUST_EXPIRED",
    "DECISION_TRUST_REVOKED",
    "DECISION_TRUST_MISMATCH",
    "DECISION_SCOPE_REQUIRED",
    "DECISION_SCOPE_INVALID",
    "DECISION_SCOPE_CONFLICT",
    "DECISION_SCOPE_MISMATCH",
    "DECISION_SCOPE_WIDENING",
    "DECISION_RELATIONSHIP_KIND_INVALID",
    "DECISION_RELATIONSHIP_TARGET_INVALID",
    "DECISION_RELATIONSHIP_CYCLE",
    "DECISION_RELATIONSHIP_CHAIN_LIMIT_EXCEEDED",
    "DECISION_RELATIONSHIP_CONFLICT",
}


def fail(message: str) -> None:
    print(f"FAIL: {message}", file=sys.stderr)
    raise SystemExit(1)


def require_file(path: Path) -> str:
    if not path.is_file():
        fail(f"missing required file: {path}")
    return path.read_text()


def sha256_text(text: str) -> str:
    return hashlib.sha256(text.encode()).hexdigest()


def validate_architecture_inventory(text: str) -> None:
    inventories = (
        ("decision", REQUIRED_DECISION_IDS, set(re.findall(r"\bAD-\d{3}\b", text))),
        ("risk", REQUIRED_RISK_IDS, set(re.findall(r"\bF13-R\d{2}\b", text))),
        (
            "conformance",
            REQUIRED_CONFORMANCE_IDS,
            set(
                re.findall(
                    r"\bF13-(?:CF|SCOPE|EVAL|SEC|TRUST|COMPAT)-\d{2}\b",
                    text,
                )
            ),
        ),
    )
    for label, expected, found in inventories:
        missing = sorted(expected - found)
        extra = sorted(found - expected)
        if missing or extra:
            fail(
                f"architecture {label} inventory mismatch; "
                f"missing={missing}, extra={extra}"
            )

    found_sections = {
        int(value)
        for value in re.findall(r"^##\s+(\d+)\.\s+", text, re.M)
        if int(value) <= 28
    }
    if found_sections != REQUIRED_ARCHITECTURE_SECTIONS:
        fail(
            "architecture top-level section inventory mismatch; "
            f"missing={sorted(REQUIRED_ARCHITECTURE_SECTIONS - found_sections)}, "
            f"extra={sorted(found_sections - REQUIRED_ARCHITECTURE_SECTIONS)}"
        )


def validate_architecture_decision_classes(text: str) -> None:
    """Keep every AD implementation class owned and fixed by architecture."""
    header = "| ID | Architectural decision | Implementation class |"
    if header not in text:
        fail("architecture decision scorecard lacks Implementation class")

    found: dict[str, str] = {}
    for line in text.splitlines():
        if not re.match(r"^\|\s*AD-\d{3}\s*\|", line):
            continue
        cells = [cell.strip() for cell in line.split("|")[1:-1]]
        if len(cells) != 10:
            fail("architecture decision scorecard row must have exactly ten columns")
        decision_id, implementation_class = cells[0], cells[2]
        if decision_id in found:
            fail(f"duplicate architecture decision scorecard row: {decision_id}")
        if implementation_class not in IMPLEMENTATION_CLASSES:
            fail(
                f"architecture decision {decision_id} has invalid implementation "
                f"class: {implementation_class}"
            )
        found[decision_id] = implementation_class

    if set(found) != REQUIRED_DECISION_IDS:
        fail(
            "architecture decision-class inventory mismatch; "
            f"missing={sorted(REQUIRED_DECISION_IDS - set(found))}, "
            f"extra={sorted(set(found) - REQUIRED_DECISION_IDS)}"
        )
    if found != DECISION_IMPLEMENTATION_CLASSES:
        mismatches = {
            decision_id: {
                "expected": DECISION_IMPLEMENTATION_CLASSES[decision_id],
                "actual": found[decision_id],
            }
            for decision_id in sorted(REQUIRED_DECISION_IDS)
            if found[decision_id] != DECISION_IMPLEMENTATION_CLASSES[decision_id]
        }
        fail(f"architecture decision implementation-class mismatch: {mismatches}")


def dependency_lock_from_handoff(text: str) -> dict[str, str]:
    match = re.search(
        r"^dependency_content_sha256:\s*\n"
        r"(?P<rows>(?:  \"[^\"]+\": [0-9a-f]{64}\s*\n)+)",
        text,
        re.M,
    )
    if not match:
        fail("ADH-2026-017 is missing a valid dependency_content_sha256 map")
    rows: dict[str, str] = {}
    for path, digest in re.findall(
        r'^  "([^"]+)": ([0-9a-f]{64})\s*$', match.group("rows"), re.M
    ):
        if path in rows:
            fail(f"duplicate dependency digest entry: {path}")
        rows[path] = digest
    return rows


def validate_dependency_lock(handoff: str) -> None:
    locked = dependency_lock_from_handoff(handoff)
    expected_paths = {str(path) for path in DEPENDENCY_PATHS}
    locked_paths = set(locked)
    if locked_paths != expected_paths:
        fail(
            "dependency digest path set mismatch; "
            f"missing={sorted(expected_paths - locked_paths)}, "
            f"extra={sorted(locked_paths - expected_paths)}"
        )
    for path in DEPENDENCY_PATHS:
        actual = sha256_text(require_file(path))
        if locked[str(path)] != actual:
            fail(f"approved dependency digest mismatch: {path}")


def validate_dependency_lifecycle() -> None:
    reuse_standard = require_file(DEPENDENCY_PATHS[0])
    feature_0011 = require_file(DEPENDENCY_PATHS[1])
    approval_0011 = require_file(DEPENDENCY_PATHS[2])
    feature_0012 = require_file(DEPENDENCY_PATHS[5])
    if not re.search(r"^status:\s*approved\s*$", reuse_standard, re.M):
        fail("FEATURE-0011 canonical reuse standard is not status: approved")
    if not re.search(r"^status:\s*implemented\s*$", feature_0011, re.M):
        fail("FEATURE-0011 feature lifecycle is not status: implemented")
    if not re.search(r"^Final feature-review status:\s*Approved\s*$", approval_0011, re.M):
        fail("FEATURE-0011 final human feature review is not Approved")
    if re.search(r"^Final feature-review status:\s*Pending\s*$", approval_0011, re.M):
        fail("FEATURE-0011 approval evidence contains an active Pending status")
    if not re.search(r"^status:\s*implemented\s*$", feature_0012, re.M):
        fail("FEATURE-0012 feature lifecycle is not status: implemented")


def validate_dependency_semantics() -> None:
    """Verify the locked FEATURE-0011/0012 contract, not only its bytes."""
    reuse_standard = require_file(
        Path("docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md")
    )
    if REUSE_SUMMARY_HEADER not in reuse_standard:
        fail("FEATURE-0011 reuse standard lacks the mandatory five-column summary")

    scope_schema = json.loads(
        require_file(Path("api/schemas/_common/scope-ref.json"))
    )
    expected_scopes = [
        "Platform",
        "Organization",
        "OrganizationUnit",
        "Tenant",
        "Project",
        "Provider",
    ]
    actual_scopes = scope_schema.get("properties", {}).get("kind", {}).get("enum")
    if actual_scopes != expected_scopes:
        fail(
            "FEATURE-0012 ScopeKind is not the locked six-value ordered vocabulary; "
            f"actual={actual_scopes}"
        )

    scope_go = require_file(Path("internal/apimeta/scope.go"))
    for value in expected_scopes:
        if f'ScopeKind = "{value}"' not in scope_go:
            fail(f"FEATURE-0012 Go ScopeKind binding omits {value}")
    if re.search(r'Scope\w+\s+ScopeKind\s*=\s*"ServiceInstance"', scope_go):
        fail("FEATURE-0012 Go ScopeKind incorrectly includes ServiceInstance")
    if "canonical scope form" not in scope_go or "maps an explicit" not in scope_go:
        fail("FEATURE-0012 Go binding omits canonical nil-Platform normalization")

    audit_schema = json.loads(require_file(Path("api/schemas/audit-event.json")))
    if audit_schema.get("x-sovrunn-profile") != "ImmutableRecord":
        fail("FEATURE-0012 AuditEvent is not bound to ImmutableRecord")
    if audit_schema.get("x-sovrunn-allowed-scopes") != ["Organization"]:
        fail("FEATURE-0012 alpha AuditEvent baseline is not Organization-only")

    architecture_0012 = require_file(Path("docs/architecture/api-resource-standard.md"))
    for marker in (
        "ManagedResource",
        "ObservedExternalResource",
        "VersionedDefinition",
        "ImmutableRecord",
        "LongRunningOperation",
        "TransientRequestResult",
        "EmbeddedValue",
        "ListEnvelope",
        "scopeRef",
        "ownerRef",
    ):
        if marker not in architecture_0012:
            fail(f"FEATURE-0012 architecture dependency omits {marker}")
    found_f12_risks = set(re.findall(r"\bF12-R\d{2}\b", architecture_0012))
    if found_f12_risks != REQUIRED_F12_RISK_IDS:
        fail(
            "FEATURE-0012 risk inventory mismatch; "
            f"missing={sorted(REQUIRED_F12_RISK_IDS - found_f12_risks)}, "
            f"extra={sorted(found_f12_risks - REQUIRED_F12_RISK_IDS)}"
        )

    design_0012 = require_file(
        Path(".kiro/specs/api-resource-naming-status-and-validation-standard/design.md")
    )
    if "Operation and AuditEvent domain payloads belong to FEATURE-0013" in design_0012:
        fail("FEATURE-0012 design retains ambiguous Operation ownership")
    for marker in (
        "The generic `Operation` contract and its `LongRunningOperation` payload remain",
        "owned by FEATURE-0012 and ADH-2026-013",
        "FEATURE-0013 reuses that contract for",
        "asynchronous handoff and does not redefine its payload",
    ):
        if marker not in design_0012:
            fail(f"FEATURE-0012 design omits Operation ownership closure: {marker}")


def validate_approved_reuse_rows(text: str) -> None:
    lines = text.splitlines()
    start = lines.index(REUSE_SUMMARY_HEADER)
    for line in lines[start + 2 :]:
        if not line.strip().startswith("|"):
            break
        cells = [cell.strip() for cell in line.strip().split("|")[1:-1]]
        if len(cells) != 5:
            continue
        status = cells[3]
        if status == "Proposed":
            fail("approved architecture contains a Proposed reuse decision")
        if status not in {"Approved", "Deferred", "Rejected", "Superseded"}:
            fail(f"approved architecture has non-final reuse status: {status}")


def validate_complete_traceability(text: str, artifact: str) -> None:
    traceability = section(text, r"^##\s+Architecture traceability\s*$")
    if not traceability:
        fail(f"{artifact} must contain a level-two Architecture traceability section")

    found_decisions = set(re.findall(r"\bAD-\d{3}\b", traceability))
    found_risks = set(re.findall(r"\bF13-R\d{2}\b", traceability))
    found_conformance = set(
        re.findall(
            r"\bF13-(?:CF|SCOPE|EVAL|SEC|TRUST|COMPAT)-\d{2}\b",
            traceability,
        )
    )
    for label, expected, found in (
        ("decision", REQUIRED_DECISION_IDS, found_decisions),
        ("risk", REQUIRED_RISK_IDS, found_risks),
        ("conformance", REQUIRED_CONFORMANCE_IDS, found_conformance),
    ):
        missing = sorted(expected - found)
        extra = sorted(found - expected)
        if missing or extra:
            fail(
                f"{artifact} {label} traceability mismatch; "
                f"missing={missing}, extra={extra}"
            )

    for number in sorted(REQUIRED_ARCHITECTURE_SECTIONS):
        prose_reference = re.search(
            rf"(?:§|section\s+){number}(?!\d)", traceability, re.I
        )
        table_reference = re.search(
            rf"^\|\s*{number}\s*\|", traceability, re.M
        )
        if not (prose_reference or table_reference):
            fail(f"{artifact} traceability omits architecture section {number}")

    traceability_lines = traceability.splitlines()
    for decision_id, expected_class in DECISION_IMPLEMENTATION_CLASSES.items():
        entry_pattern = re.compile(
            rf"^(?:\|\s*`?{re.escape(decision_id)}`?\s*\||"
            rf"[-*]\s*`?{re.escape(decision_id)}`?\s*(?::|-))"
        )
        entries = [line for line in traceability_lines if entry_pattern.search(line.strip())]
        if len(entries) != 1:
            fail(
                f"{artifact} traceability must contain exactly one explicit "
                f"entry for {decision_id}; found {len(entries)}"
            )
        classes = {
            value for value in IMPLEMENTATION_CLASSES if value in entries[0]
        }
        if classes != {expected_class}:
            fail(
                f"{artifact} traceability class mismatch for {decision_id}; "
                f"expected={expected_class}, found={sorted(classes)}"
            )


def validate_reuse_summary(text: str, artifact: str) -> None:
    """Enforce the FEATURE-0011 mandatory feature-level summary contract."""
    lines = text.splitlines()
    header_indexes = [
        index for index, line in enumerate(lines) if line.strip() == REUSE_SUMMARY_HEADER
    ]
    if len(header_indexes) != 1:
        fail(
            f"{artifact} must contain exactly one FEATURE-0011 reuse summary; "
            f"found {len(header_indexes)}"
        )

    header_index = header_indexes[0]
    if header_index + 1 >= len(lines) or not re.fullmatch(
        r"\|(?:\s*:?-{3,}:?\s*\|){5}", lines[header_index + 1].strip()
    ):
        fail(f"{artifact} reuse summary has an invalid five-column separator")

    rows: list[list[str]] = []
    for line in lines[header_index + 2 :]:
        stripped = line.strip()
        if not stripped.startswith("|"):
            break
        cells = [cell.strip() for cell in stripped.split("|")[1:-1]]
        if len(cells) != 5:
            fail(f"{artifact} reuse summary row does not have exactly five columns")
        rows.append(cells)

    if not rows:
        fail(f"{artifact} reuse summary contains no capability rows")

    for cells in rows:
        capability, disposition, rationale, decision_status, controlling_reference = cells
        if not all((capability, rationale, controlling_reference)):
            fail(f"{artifact} reuse summary contains an empty required cell")
        if disposition not in REUSE_DISPOSITIONS:
            fail(
                f"{artifact} reuse summary has invalid disposition: {disposition}"
            )
        if decision_status not in REUSE_STATUSES:
            fail(
                f"{artifact} reuse summary has invalid decision status: "
                f"{decision_status}"
            )


def validate_limit_layering(text: str, artifact: str) -> None:
    """Preserve FEATURE-0012 ceiling >= profile ceiling >= evaluator ceiling."""
    inherited_outer_ceiling = re.search(
        r"FEATURE-0012.{0,160}(?:ceiling|bound).{0,80}(?:absolute|outer)"
        r"|FEATURE-0012.{0,160}(?:absolute|outer).{0,80}(?:ceiling|bound)",
        text,
        re.I | re.S,
    )
    profile_layer = re.search(
        r"(?:DecisionProfile|profile).{0,180}(?:within|below|narrow|must not "
        r"exceed|may not exceed).{0,180}(?:FEATURE-0012|platform|schema).{0,80}"
        r"(?:ceiling|bound)"
        r"|(?:FEATURE-0012|platform|schema).{0,120}(?:ceiling|bound).{0,180}"
        r"(?:DecisionProfile|profile)",
        text,
        re.I | re.S,
    )
    evaluator_narrows = re.search(
        r"evaluator.{0,140}(?:may|must|can).{0,50}(?:only\s+)?narrow"
        r"|evaluator.{0,140}narrow.{0,80}(?:never|must not|may not|cannot)"
        r".{0,50}widen",
        text,
        re.I | re.S,
    )
    most_restrictive = re.search(
        r"most restrictive applicable (?:value|limit|bound)", text, re.I
    )
    if not inherited_outer_ceiling:
        fail(f"{artifact} omits FEATURE-0012 as the absolute outer limit ceiling")
    if not profile_layer:
        fail(f"{artifact} omits the DecisionProfile limit layer")
    if not evaluator_narrows:
        fail(f"{artifact} does not constrain evaluators to narrowing profile limits")
    if not most_restrictive:
        fail(f"{artifact} omits the most-restrictive-applicable-limit rule")


def validate_public_error_contract(text: str, artifact: str) -> None:
    """Keep public FEATURE-0013 error semantics owned by architecture."""
    required_bindings = (
        "urn:sovrunn:problem:validation-failed",
        "VALIDATION_FAILED",
        "violations[].code",
        "RFC 6901",
    )
    for binding in required_bindings:
        if binding not in text:
            fail(f"{artifact} omits public error binding: {binding}")
    if not re.search(r"(?:HTTP\s*)?422", text):
        fail(f"{artifact} omits HTTP 422 from the validation error binding")

    found_codes = set(
        re.findall(
            r"\bDECISION_(?:PROFILE|EVALUATION|COMPOSITION|OBLIGATION|TRUST|"
            r"SCOPE|RELATIONSHIP)_[A-Z][A-Z0-9_]*\b",
            text,
        )
    )
    unknown = sorted(found_codes - PUBLIC_VIOLATION_CODES)
    missing = sorted(PUBLIC_VIOLATION_CODES - found_codes)
    if unknown:
        fail(f"{artifact} invents public violation codes: {', '.join(unknown)}")
    if missing:
        fail(f"{artifact} omits public violation codes: {', '.join(missing)}")


def validate_immutable_erasure_boundary(text: str, artifact: str) -> None:
    """Prevent lawful erasure from becoming in-place record mutation."""
    if not re.search(
        r"(?:DecisionRecord|AuditEvent).{0,100}append-only"
        r"|append-only.{0,100}(?:DecisionRecord|AuditEvent)",
        text,
        re.I | re.S,
    ):
        fail(f"{artifact} omits append-only DecisionRecord/AuditEvent semantics")
    if not re.search(
        r"(?:erasure|retention).{0,160}(?:never|must not|shall not).{0,100}"
        r"(?:rewrite|mutate|update)"
        r"|(?:never|must not|shall not).{0,100}(?:rewrite|mutate|update)"
        r".{0,160}(?:erasure|retention|canonical)",
        text,
        re.I | re.S,
    ):
        fail(f"{artifact} does not prohibit in-place mutation for erasure")
    alternatives = {
        "separately controlled payload custody": r"separately\s+controlled\s+payload\s+custody",
        "key destruction": r"(?:cryptographic\s+)?key\s+destruction",
        "linked tombstone/redaction record": r"linked\s+tombstone(?:/redaction)?\s+record",
    }
    for name, pattern in alternatives.items():
        if not re.search(pattern, text, re.I):
            fail(f"{artifact} omits immutable erasure alternative: {name}")
    if not re.search(
        r"(?:implements? no|does not implement|must not implement|defer(?:red|s)?)"
        r".{0,100}(?:erasure|cryptographic process|key destruction)",
        text,
        re.I | re.S,
    ):
        fail(f"{artifact} may authorize an erasure runtime in FEATURE-0013")


def validate_request_response_boundary(text: str, artifact: str) -> None:
    """Prevent downstream invention of request or pending-response contracts."""
    required = (
        "DecisionRequest",
        "calling-domain-owned",
        "TransientRequestResult",
        "DecisionRecord",
        "Operation",
        "Problem Details",
    )
    for marker in required:
        if marker not in text:
            fail(f"{artifact} omits request/outcome boundary marker: {marker}")

    if not re.search(
        r"(?:no|not|must not|does not|without).{0,100}"
        r"(?:shared|canonical|common).{0,40}`?DecisionRequest`?"
        r"|`?DecisionRequest`?.{0,100}(?:no|not|must not|does not|without)"
        r".{0,60}(?:shared|canonical|common)",
        text,
        re.I | re.S,
    ):
        fail(f"{artifact} does not prohibit a shared FEATURE-0013 DecisionRequest")

    if not re.search(
        r"(?:no|not|must not|does not|without).{0,100}"
        r"(?:pending|non-final).{0,40}(?:response|envelope|status)"
        r"|(?:pending|non-final).{0,40}(?:response|envelope|status).{0,100}"
        r"(?:no|not|must not|does not|without)",
        text,
        re.I | re.S,
    ):
        fail(f"{artifact} does not prohibit a pending-decision response contract")

    if re.search(r"\btyped\s+non-final\s+response\b", text, re.I):
        fail(f"{artifact} permits an unspecified typed non-final response")

    negation = re.compile(
        r"\b(?:no|not|never|must not|shall not|prohibit(?:ed|s)?|"
        r"forbid(?:den|s)?|reject(?:ed|s)?|omit(?:ted|s)?)\b",
        re.I,
    )
    pending_name = re.compile(
        r"\b(?:PendingDecision|DeferredDecision|DecisionPending)"
        r"(?:Response|Status)?\b",
        re.I,
    )
    pending_definition = re.compile(
        r"\b(?:create|define|implement|generate|add)\b.{0,80}"
        r"\b(?:PendingDecision|DeferredDecision|DecisionPending)"
        r"(?:Response|Status)?\b",
        re.I,
    )
    request_definition = re.compile(
        r"\b(?:create|define|implement|generate|add)\b.{0,80}"
        r"\b(?:shared|canonical|common)\b.{0,40}\bDecisionRequest\b",
        re.I,
    )
    lines = text.splitlines()
    section_heading = ""
    for line_number, line in enumerate(lines, 1):
        if re.match(r"^##\s+", line):
            section_heading = line
        context = " ".join(lines[max(0, line_number - 3) : line_number])
        exclusion_context = bool(
            negation.search(context)
            or re.search(
                r"\b(?:non-goals?|prohibited|excluded)\b",
                section_heading,
                re.I,
            )
        )
        if pending_definition.search(line) and not negation.search(line):
            fail(
                f"{artifact} introduces a pending-decision contract at line "
                f"{line_number}"
            )
        if pending_name.search(line) and not exclusion_context:
            fail(
                f"{artifact} introduces a pending-decision contract at line "
                f"{line_number}"
            )
        if request_definition.search(line) and not negation.search(line):
            fail(
                f"{artifact} introduces a shared DecisionRequest at line "
                f"{line_number}"
            )


def validate_risk_namespace(text: str, artifact: str) -> None:
    """Keep FEATURE-0012 and FEATURE-0013 risk identifiers disjoint."""
    forbidden_claims = (
        r"\bFEATURE-0013\b.{0,120}\b(?:owns?|inherits?|preserves?|adopts?)\b"
        r".{0,120}\bF12-R\d{2}\b",
        r"\bF12-R\d{2}\b.{0,120}\b(?:owned by|belongs to)\s+FEATURE-0013\b",
        r"\bF12-R\d{2}\b.{0,120}\b(?:becomes?|maps? to|aliases?)\b"
        r".{0,80}\bF13-R\d{2}\b",
    )
    for pattern in forbidden_claims:
        match = re.search(pattern, text, re.I | re.S)
        if match:
            excerpt = " ".join(match.group(0).split())
            fail(
                f"{artifact} claims FEATURE-0013 ownership or inheritance of "
                f"F12 risk identifiers: {excerpt}"
            )


def validate_closed_runtime_wording(text: str, artifact: str) -> None:
    """Reject legacy wording that reopens runtime, signing, or graph storage."""
    forbidden_literals = (
        "Synchronization is explicit, signed,",
        "origin trust-domain identity, signed manifests",
        "graph storage shape",
    )
    for literal in forbidden_literals:
        if literal in text:
            fail(f"{artifact} retains architecture-leaking wording: {literal}")

    positive_runtime = re.compile(
        r"\b(?:create|define|design|implement|add|use|select)\b.{0,100}"
        r"\b(?:persistent graph|graph database|graph repository|"
        r"signed synchronization|signed manifest)\b",
        re.I,
    )
    negation = re.compile(
        r"\b(?:no|not|never|must not|shall not|prohibit(?:ed|s)?|"
        r"defer(?:red|s)?|later invariant|unauthori[sz]ed)\b",
        re.I,
    )
    lines = text.splitlines()
    for line_number, line in enumerate(lines, 1):
        context = " ".join(lines[max(0, line_number - 2) : line_number + 1])
        if positive_runtime.search(line) and not negation.search(context):
            fail(
                f"{artifact} may authorize prohibited runtime/signing/graph "
                f"storage at line {line_number}"
            )



def validate_known_forbidden_paths(text: str, artifact: str) -> None:
    """Prevent regenerated artifacts from introducing parallel workflows."""
    forbidden_paths = (
        "api/schemas/SCHEMA_BASELINE_MANIFEST.json",
        "api/schemas/diffs/",
        "api/schemas/approvals/",
    )
    negation = re.compile(
        r"\b(?:no|not|never|must not|shall not|prohibit(?:ed|s)?|"
        r"without|do not|does not|is not|are not)\b",
        re.I,
    )
    lines = text.splitlines()
    for line_number, line in enumerate(lines, 1):
        for forbidden_path in forbidden_paths:
            if forbidden_path not in line:
                continue
            context = " ".join(lines[max(0, line_number - 2) : line_number + 1])
            if not negation.search(context):
                fail(
                    f"{artifact} introduces forbidden parallel baseline path "
                    f"{forbidden_path} at line {line_number}"
                )


def validate_feature_0013_tasks_guardrails(text: str) -> None:
    """Make tasks generation fail closed before Cursor can wander."""
    required_markers = (
        "ADH-2026-017",
        "section 27.9",
        "SecurityExceptionRef",
        "GraphEdge",
        "BASELINE_MANIFEST.json",
        "BASELINE_APPROVALS.json",
        "metadata.scopeRef",
        "calling-domain-owned",
        "DecisionRequest",
        "pending-decision",
        "CONTRACT_NOW",
        "INVARIANT_FOR_LATER",
        "DEFERRED",
        "Architecture traceability",
    )
    for marker in required_markers:
        if marker not in text:
            fail(f"tasks guardrail missing required marker: {marker}")

    forbidden_task_patterns = (
        (
            r"(?:create|implement|add|define|generate).{0,80}"
            r"(?:shared|canonical|common).{0,40}DecisionRequest",
            "shared DecisionRequest implementation",
        ),
        (
            r"(?:create|implement|add|define|generate).{0,80}"
            r"(?:PendingDecision|DeferredDecision|DecisionPending|pending[- ]decision)",
            "pending-decision response implementation",
        ),
        (
            r"(?:create|implement|add|define|generate).{0,80}"
            r"(?:AuditScope|second scope|parallel scope)",
            "parallel scope authority implementation",
        ),
        (
            r"(?:create|implement|add|define|generate).{0,80}"
            r"(?:SCHEMA_BASELINE_MANIFEST|api/schemas/diffs|api/schemas/approvals)",
            "parallel baseline workflow implementation",
        ),
        (
            r"(?:create|implement|add|define|generate).{0,80}"
            r"(?:approval workflow|security exception workflow|exception approval service)",
            "runtime approval workflow implementation",
        ),
        (
            r"(?:calculate|derive|infer|recompute|downgrade|upgrade).{0,80}"
            r"(?:residual risk|Matrix E risk|risk level)",
            "automated Matrix E residual-risk calculation",
        ),
    )
    negation = re.compile(
        r"\b(?:no|not|never|must not|shall not|prohibit(?:ed|s)?|"
        r"without|do not|does not|is not|are not)\b",
        re.I,
    )
    lines = text.splitlines()
    for line_number, line in enumerate(lines, 1):
        context = " ".join(lines[max(0, line_number - 2) : line_number + 1])
        for pattern, meaning in forbidden_task_patterns:
            if re.search(pattern, line, re.I) and not negation.search(context):
                fail(f"tasks may plan {meaning} at line {line_number}")

    if "Only `CONTRACT_NOW`" not in text and "only `CONTRACT_NOW`" not in text:
        fail("tasks must state that only CONTRACT_NOW items produce implementation tasks")
    if not re.search(r"no\s+source\s+task|produce\s+no\s+source\s+task", text, re.I):
        fail("tasks must state INVARIANT_FOR_LATER/DEFERRED rows produce no source task")

def architecture_preflight(*, require_approval: bool = True) -> None:
    architecture = require_file(ARCHITECTURE)
    for marker in REQUIRED_ARCHITECTURE_CONTROLS:
        if marker not in architecture:
            fail(f"architecture closure control missing: {marker}")
    validate_reuse_summary(architecture, "architecture")
    validate_limit_layering(architecture, "architecture")
    validate_public_error_contract(architecture, "architecture")
    validate_immutable_erasure_boundary(architecture, "architecture")
    validate_request_response_boundary(architecture, "architecture")
    validate_risk_namespace(architecture, "architecture")
    validate_closed_runtime_wording(architecture, "architecture")
    validate_architecture_inventory(architecture)
    validate_architecture_decision_classes(architecture)
    validate_dependency_lifecycle()
    validate_dependency_semantics()

    handoffs = sorted(
        Path("docs/reviews/architecture-decision-handoffs").glob(
            "ADH-2026-017-feature-0013-*.md"
        )
    )
    if len(handoffs) != 1:
        fail("exactly one ADH-2026-017 FEATURE-0013 consolidation handoff is required")
    handoff = handoffs[0].read_text()
    validate_dependency_lock(handoff)
    validate_risk_namespace(handoff, "ADH-2026-017")

    agent = require_file(AGENT)
    for resource in REQUIRED_AGENT_RESOURCES:
        if resource not in agent:
            fail(f"Kiro agent is missing controlling resource: {resource}")
    for marker in (
        "`LongRunningOperation` payload remain owned by FEATURE-0012/ADH-2026-013",
        "Every AD-001 through",
        "AD-045 implementation class is owned by architecture section 19",
        "Only `CONTRACT_NOW` may produce FEATURE-0013",
    ):
        if marker not in agent:
            fail(f"Kiro agent is missing closed architecture instruction: {marker}")
    for historical in ("ADH-2026-014-*", "ADH-2026-015-*", "ADH-2026-016-*"):
        if f"file://docs/reviews/architecture-decision-handoffs/{historical}" in agent:
            fail(f"Kiro agent loads historical handoff as a resource: {historical}")
    for forbidden_resource in (
        "file://docs/architecture/**/*.md",
        "file://docs/phase2/*.md",
        "observability-and-audit-baseline.md",
        "organization-governance.md",
    ):
        if forbidden_resource in agent:
            fail(f"Kiro agent loads an over-broad or conflicting resource: {forbidden_resource}")

    for script_path in (KIRO_STAGE, SPEC_FLOW):
        script = require_file(script_path)
        if 'KIRO_AGENT="sovrunn-spec"' not in script:
            fail(f"FEATURE-0013 Kiro agent is not mandatory in {script_path}")
        if "FEATURE-0013 requires KIRO_AGENT=sovrunn-spec" not in script:
            fail(f"FEATURE-0013 agent override does not fail closed in {script_path}")

    for prompt_path in PROMPT_TEMPLATES:
        prompt = require_file(prompt_path)
        for marker in (
            "Architecture traceability",
            "ARCHITECTURE_DECISION_REQUIRED",
            "DecisionRecord",
            "AD-001",
            "AD-045",
            "F13-R01",
            "F13-R31",
            "F13-CF-01",
            "F13-COMPAT-09",
            "F12-R01",
            "calling-domain-owned",
            "TransientRequestResult",
            "pending",
            "Implementation class is architecture-owned",
            "owned by FEATURE-0012/ADH-2026-013",
        ):
            if marker not in prompt:
                fail(f"{prompt_path} control missing: {marker}")

    for path in ACTIVE_CONTEXT:
        context = require_file(path)
        if "ADH-2026-017" not in context:
            fail(f"active context does not reference consolidated handoff: {path}")
        if re.search(r"joint controlling", context, re.I):
            fail(f"active context still declares multiple controlling handoffs: {path}")
        if re.search(r"(?:seven[- ]value|seven canonical).{0,80}\bScopeKind\b", context, re.I):
            fail(f"active context still declares the superseded scope model: {path}")

    state = json.loads(require_file(STATE))
    if state.get("title") != "Decision Record and AuditEvent Standard":
        fail("FEATURE-0013 state has stale title")

    if not require_approval:
        return

    # Approval is checked last. A pending review therefore still exercises
    # every deterministic architecture, dependency, context, prompt, agent,
    # and traceability control before stopping at the human gate.
    if not re.search(
        r"^status:\s*approved-for-kiro-requirements\s*$", architecture, re.M
    ):
        fail("consolidated architecture is not approved-for-kiro-requirements")
    if (
        "Human architecture review status: **APPROVED_FOR_KIRO_REQUIREMENTS**."
        not in architecture
    ):
        fail("fresh human architecture approval is not recorded")
    if re.search(r"Consolidation (?:reviewer|decision date):\s*pending", architecture):
        fail("consolidation reviewer or decision date is pending")
    validate_approved_reuse_rows(architecture)

    if not re.search(
        r"(?:Approval status|Status):\s*(?:\*\*)?Approved(?:\*\*)?", handoff, re.I
    ):
        fail(f"consolidation handoff is not Approved: {handoffs[0]}")
    digest_match = re.search(
        r"^architecture_content_sha256:\s*([0-9a-f]{64})\s*$", handoff, re.M
    )
    if not digest_match:
        fail("consolidation handoff is missing the approved architecture SHA-256")
    architecture_digest = hashlib.sha256(architecture.encode()).hexdigest()
    if digest_match.group(1) != architecture_digest:
        fail("consolidation handoff architecture SHA-256 does not match")


def section(text: str, heading_pattern: str) -> str:
    match = re.search(heading_pattern, text, re.I | re.M)
    if not match:
        return ""
    start = match.end()
    next_heading = re.search(r"^##\s+", text[start:], re.M)
    end = start + next_heading.start() if next_heading else len(text)
    return text[start:end]


def post_generation(stage: str) -> None:
    architecture_preflight()
    target = STAGE_FILES[stage]
    text = require_file(target)
    validate_reuse_summary(text, stage)
    validate_limit_layering(text, stage)
    validate_public_error_contract(text, stage)
    validate_immutable_erasure_boundary(text, stage)
    validate_request_response_boundary(text, stage)
    validate_risk_namespace(text, stage)
    validate_closed_runtime_wording(text, stage)
    validate_known_forbidden_paths(text, stage)
    if stage == "tasks":
        validate_feature_0013_tasks_guardrails(text)

    if "SUPERSEDED DRAFT" in text:
        fail(f"active {stage} artifact is a superseded draft")
    for pattern, meaning in FORBIDDEN_GLOBAL:
        if re.search(pattern, text, re.I):
            fail(f"{stage} reopens or contains {meaning}")
    for line_number, line in enumerate(text.splitlines(), 1):
        service_scope = re.search(
            r"\bServiceInstance\b\s+(?:is|as|becomes?|may\s+be|can\s+be|"
            r"is\s+added\s+as).{0,30}\b(?:governance\s+)?(?:scope|ScopeKind)\b",
            line,
            re.I,
        )
        explicit_rejection = re.search(
            r"\b(?:not|never|reject(?:ed|s)?|invalid|prohibit(?:ed|s)?)\b",
            line,
            re.I,
        )
        if service_scope and not explicit_rejection:
            fail(
                f"{stage} treats ServiceInstance as a governance scope at line "
                f"{line_number}"
            )

    validate_complete_traceability(text, stage)
    traceability = section(text, r"^##\s+Architecture traceability\s*$")
    if "FEATURE-0013-decision-record-and-auditevent-standard.md" not in traceability:
        fail(f"{stage} traceability does not cite the consolidated architecture")
    for required_section in ("5.4", "6.1", "6.8", "7.1", "9", "12.3", "12.4", "15", "17", "27.8"):
        prose_reference = re.search(
            rf"(?:§|section\s+){re.escape(required_section)}\b",
            traceability,
            re.I,
        )
        table_reference = re.search(
            rf"^\|\s*{re.escape(required_section)}\s*\|",
            traceability,
            re.M,
        )
        if not (prose_reference or table_reference):
            fail(f"{stage} traceability omits architecture section {required_section}")

    questions = section(text, r"^##\s+.*Design Questions.*$")
    for pattern, meaning in FORBIDDEN_DESIGN_QUESTION:
        if questions and re.search(pattern, questions, re.I | re.S):
            fail(f"design questions reopen closed architecture: {meaning}")

    positive_crypto_action = re.compile(
        r"(?:shall|must|will|to)\s+"
        r"(?:implement|execute|compute|generate|verify|select|choose).{0,100}"
        r"(?:RFC\s*8785|canonicali[sz]ation algorithm|content[- ]to[- ]digest|"
        r"signature algorithm|signature verification|cryptographic (?:product|service))",
        re.I,
    )
    negation = re.compile(
        r"\b(?:not|must not|shall not|prohibit(?:ed|s)?|defer(?:red|s)?|"
        r"unauthori[sz]ed|until ADR-F13-002)\b",
        re.I,
    )
    for line_number, line in enumerate(text.splitlines(), 1):
        if positive_crypto_action.search(line) and not negation.search(line):
            fail(
                f"{stage} may authorize cryptography deferred by ADR-F13-002 "
                f"at line {line_number}"
            )

    if stage in {"requirements", "design", "tasks"}:
        for scope_class in ("CONTRACT_NOW", "INVARIANT_FOR_LATER", "DEFERRED"):
            if scope_class not in text:
                fail(f"{stage} is missing scope class {scope_class}")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--feature", required=True)
    parser.add_argument("--stage", choices=tuple(STAGE_FILES), default="requirements")
    parser.add_argument(
        "--mode", choices=("readiness", "pre", "post"), default="post"
    )
    args = parser.parse_args()

    if args.feature != FEATURE:
        print(f"SKIP: architecture-boundary check does not apply to {args.feature}")
        return
    if args.mode == "readiness":
        architecture_preflight(require_approval=False)
    elif args.mode == "pre":
        architecture_preflight()
    else:
        post_generation(args.stage)
    print(f"PASS: {FEATURE} architecture boundary ({args.mode}, {args.stage})")


if __name__ == "__main__":
    main()
