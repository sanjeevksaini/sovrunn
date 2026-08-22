#!/usr/bin/env python3
"""Fail-closed architecture-readiness check for FEATURE-0016.

Validates that ADH-2026-058's atomic replacement of the FEATURE-0016
placeholder registry definitions (VS0-SCHEMA-015..017, VS0-STATE-004) is
representable and consistent across the repository authorities before
requirements generation may proceed.

This checker rejects reintroduction of the removed placeholder concepts
(persisted status.availability, Draining, publicly projected Qualifying,
any alias/migration/legacy-compatibility behavior for them) and validates
the approved five-route surface, closed create-field boundary, sole
lifecycle-service writer, sole synthetic observer, idempotency scope, audit
matrix, and closed local violation set described in ADH-2026-058.

Exit 0 = PASS (requirements generation may proceed).
Exit 1 = FAIL (architecture gap remains; requirements generation blocked).
"""

from __future__ import annotations

import sys
from pathlib import Path

try:
    import yaml
except ImportError:
    print("FAIL: PyYAML required. Install: pip install pyyaml", file=sys.stderr)
    sys.exit(1)

ROOT = Path(__file__).resolve().parents[1]

REG_PATH = ROOT / "docs/architecture/vertical-slices/VS-000-contract-registry.yaml"
TRACE_PATH = ROOT / "docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md"
F16_ARCH = ROOT / "docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md"
F16_FEAT = ROOT / "docs/features/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md"
CLOSURE = ROOT / "docs/reviews/architecture-readiness/FEATURE-0016-architecture-closure-matrix.md"
REUSE = ROOT / "docs/reviews/reuse-assessments/FEATURE-0016-approval-evidence.md"
STEER = ROOT / ".kiro/steering/slice0-contract.md"
CHARTER = ROOT / "docs/architecture/vertical-slices/VS-000-core-skeleton.md"
CONTROL = ROOT / ".automation/features/FEATURE-0016.control.json"
ADH_058 = ROOT / "docs/reviews/architecture-decision-handoffs/ADH-2026-058-feature-0016-adapter-and-executiontarget-executable-contract-closure.md"
ADH_066 = ROOT / "docs/reviews/architecture-decision-handoffs/ADH-2026-066-feature-0016-coherent-backing-access-compatibility-bridge.md"

F16_REQUIRED_CF_IDS = [f"VS0-CF-F16-{i:02d}" for i in range(1, 129)]

# Removed placeholder concepts that must never reappear as active FEATURE-0016 behavior.
REMOVED_PLACEHOLDER_TERMS = [
    "status.availability",
    "Draining",
]

REMOVED_ALIAS_FIELD_TERMS = [
    "spec.participationRef",
]

F16_ROUTES = [
    ("POST", "/apis/execution.sovrunn.io/v1alpha1/execution-targets"),
    ("GET", "/apis/execution.sovrunn.io/v1alpha1/execution-targets"),
    ("GET", "/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}"),
    ("POST", "/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify"),
    ("POST", "/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire"),
]

F16_CLOSED_VIOLATIONS = [
    "VS0_TARGET_RETIRED",
    "VS0_TARGET_MAINTENANCE",
    "VS0_TARGET_QUALIFICATION_IN_PROGRESS",
    "VS0_TARGET_EPOCH_STALE",
    "VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE",
    "VS0_EXECUTION_TARGET_STACK_UNAVAILABLE",
    "VS0_EXECUTION_TARGET_SCOPE_MISMATCH",
    "VS0_EXECUTION_TARGET_VIABILITY_STALE",
]

errs: list[str] = []


def e(msg: str) -> None:
    errs.append(msg)


def read(path: Path) -> str:
    if not path.exists():
        e(f"Required file missing: {path.relative_to(ROOT)}")
        return ""
    return path.read_text()


def is_allowed_removed_mention(text: str, term: str) -> bool:
    """A removed-concept mention is allowed only in an explicit exclusion/
    removal/non-goal context (e.g. 'no active persisted status.availability',
    'Draining is removed'), on a line naming a different resource's field
    (e.g. ServiceRegion.status.availability, which is FEATURE-0022 scope and
    not the removed ExecutionTarget placeholder), or in a table cell whose
    row already carries a removal marker anywhere in that row/line."""
    term_lower = term.lower()
    allow_markers = (
        "no active", "must not", "no longer", "removed", "does not", "never",
        "not active", "no persisted", "excluded", "non-goal", "must never",
        "placeholder", "prohibited", "rejected", "no ",
    )
    other_resource_markers = ("serviceregion",)
    heading_allow_markers = ("removed", "excluded", "non-goal", "prohibited")
    lines = text.splitlines()
    for i, line in enumerate(lines):
        lower_line = line.lower()
        if term_lower not in lower_line:
            continue
        if any(marker in lower_line for marker in other_resource_markers):
            continue
        if any(marker in lower_line for marker in allow_markers):
            continue
        # Heading-aware fallback: find the nearest preceding heading and check
        # whether it establishes a removal/exclusion/non-goal context.
        heading = ""
        for j in range(i, -1, -1):
            if lines[j].strip().startswith("#"):
                heading = lines[j].lower()
                break
        if any(marker in heading for marker in heading_allow_markers):
            continue
        return False
    return True


def check_no_reintroduced_placeholders(reg: dict, f16_arch: str, f16_feat: str) -> None:
    schemas = reg.get("schemas", [])
    execution_target = next((s for s in schemas if s.get("identity", "").endswith("/ExecutionTarget")), None)
    if not execution_target:
        e("PLACEHOLDER: VS0-SCHEMA-015 (ExecutionTarget) not found in registry")
        return
    required = " ".join(execution_target.get("required", []))
    if "status.availability" in required:
        e("PLACEHOLDER: VS0-SCHEMA-015 must not require a persisted status.availability field (ADH-2026-058)")
    if "Draining" in required:
        e("PLACEHOLDER: VS0-SCHEMA-015 must not reference Draining as an active state (ADH-2026-058)")
    field_ownership = execution_target.get("fieldOwnership", {})
    if "status.availability" in field_ownership:
        e("PLACEHOLDER: VS0-SCHEMA-015 fieldOwnership must not own status.availability (ADH-2026-058)")
    if "spec.participationRef" in field_ownership:
        e("PLACEHOLDER: VS0-SCHEMA-015 fieldOwnership must not retain the renamed spec.participationRef alias (ADH-2026-058)")

    state_machines = reg.get("stateMachines", [])
    sm_004 = next((sm for sm in state_machines if sm.get("id") == "VS0-STATE-004"), None)
    if not sm_004:
        e("PLACEHOLDER: VS0-STATE-004 not found in registry")
    else:
        states = sm_004.get("states", [])
        for state in states:
            if "Draining" in state:
                e(f"PLACEHOLDER: VS0-STATE-004 must not contain a Draining state, found {state!r} (ADH-2026-058)")
            if "Qualifying" in state:
                e(f"PLACEHOLDER: VS0-STATE-004 must not persist Qualifying as a public state row, found {state!r} (ADH-2026-058 clause 7)")
        in_flight = str(sm_004.get("inFlightReservation", ""))
        if "never a projected resource status" not in in_flight and "never persisted or projected" not in in_flight.lower():
            e("PLACEHOLDER: VS0-STATE-004 must state Qualifying is a lifecycle-service-only in-flight reservation, never persisted or projected (ADH-2026-058 clause 7)")

    for label, text in (("architecture", f16_arch), ("feature", f16_feat)):
        if not text:
            continue
        for term in REMOVED_PLACEHOLDER_TERMS:
            if not is_allowed_removed_mention(text, term):
                e(f"PLACEHOLDER: F0016 {label} mentions '{term}' outside an explicit removal/exclusion context")
        if "publicly projected qualifying" not in text.lower() and "never a projected resource status" not in text.lower() and "never persisted or projected" not in text.lower():
            e(f"PLACEHOLDER: F0016 {label} must explicitly state Qualifying is never publicly projected (ADH-2026-058)")


def check_five_route_surface(f16_arch: str, f16_feat: str) -> None:
    # The architecture authority is the exact route source of truth; the feature
    # authority may reference it by clause rather than restating every path,
    # but must still assert the five-route/one-guard closure and prohibition.
    if f16_arch:
        lower = f16_arch.lower()
        for method, path in F16_ROUTES:
            if path.lower() not in lower:
                e(f"ROUTES: F0016 architecture must state the exact route {method} {path} (ADH-2026-058 clause 3)")
        if "no patch" not in lower:
            e("ROUTES: F0016 architecture must explicitly state no PATCH/PUT/DELETE/HEAD public route exists (ADH-2026-058 clause 3)")
        if "pre-servemux" not in lower:
            e("ROUTES: F0016 architecture must describe the pre-ServeMux transport-only method/path guard (ADH-2026-058 clause 3)")
    if f16_feat:
        lower = f16_feat.lower()
        if "exactly five" not in lower and "five routes" not in lower and "five explicit" not in lower:
            e("ROUTES: F0016 feature must state the exact five-route closure (ADH-2026-058 clause 3)")
        if "pre-servemux" not in lower and "transport-only" not in lower:
            e("ROUTES: F0016 feature must reference the pre-ServeMux transport-only guard (ADH-2026-058 clause 3)")


def check_closed_violation_set(reg: dict, f16_arch: str) -> None:
    vcs = {vc.get("code") for vc in reg.get("violationCodes", {}).get("slice0", [])}
    for code in F16_CLOSED_VIOLATIONS:
        if code not in vcs:
            e(f"VIOLATIONS: {code} must be a registered violation code (ADH-2026-058 clause 9)")
        if f16_arch and code not in f16_arch:
            e(f"VIOLATIONS: F0016 architecture must reference {code} (ADH-2026-058 clause 9)")


def check_sole_writer_and_observer(reg: dict, f16_arch: str) -> None:
    writers = reg.get("writers", [])
    writer_006 = next((w for w in writers if w.get("id") == "VS0-WRITER-006"), None)
    if not writer_006:
        e("WRITER: VS0-WRITER-006 not found in registry")
    else:
        if writer_006.get("writer") != "ExecutionTargetLifecycleService":
            e(f"WRITER: VS0-WRITER-006 writer must be ExecutionTargetLifecycleService, found {writer_006.get('writer')!r} (ADH-2026-058 clause 7)")
        paths = writer_006.get("paths", [])
        if not any("ExecutionTarget.status" in p for p in paths):
            e("WRITER: VS0-WRITER-006 must own ExecutionTarget.status (ADH-2026-058 clause 7)")
    if f16_arch and "sovrunn.synthetic-iaas-observer/v1" not in f16_arch:
        e("OBSERVER: F0016 architecture must name the sole observer sovrunn.synthetic-iaas-observer/v1 (ADH-2026-058 clause 6)")


def check_maintenance_race_mutual_exclusivity(reg: dict, f16_arch: str) -> None:
    """ADH-2026-060: F16-75 (inactive-marker epoch-stale) and F16-89
    (active-marker Maintenance-wins) must be mutually exclusive and never
    reversed. F16-75 must not say Maintenance entry wins; F16-89 must."""
    confs = reg.get("conformance", [])
    by_id = {c.get("id"): c for c in confs}
    f16_75 = by_id.get("VS0-CF-F16-75")
    f16_89 = by_id.get("VS0-CF-F16-89")
    if not f16_75:
        e("MAINTENANCE-RACE: VS0-CF-F16-75 not found in registry (ADH-2026-060)")
    if not f16_89:
        e("MAINTENANCE-RACE: VS0-CF-F16-89 not found in registry (ADH-2026-060)")
    if f16_75 and f16_89:
        f75_inputs = str(f16_75.get("inputs", "")).lower()
        f89_inputs = str(f16_89.get("inputs", "")).lower()
        if "maintenance entry wins" in f75_inputs or "maintenance-entry wins" in f75_inputs:
            e("MAINTENANCE-RACE: VS0-CF-F16-75 must not say Maintenance entry wins (ADH-2026-060); reversed predicate")
        if "no active current-maintenance marker" not in f75_inputs:
            e("MAINTENANCE-RACE: VS0-CF-F16-75 input must state no active current-Maintenance marker exists (ADH-2026-060)")
        if "maintenance entry wins" not in f89_inputs:
            e("MAINTENANCE-RACE: VS0-CF-F16-89 input must state Maintenance entry wins (ADH-2026-060)")
        if f16_75.get("expectedError") != "STALE_RESOURCE_VERSION" or f16_75.get("expectedViolation") != "VS0_TARGET_EPOCH_STALE":
            e("MAINTENANCE-RACE: VS0-CF-F16-75 must be 412 STALE_RESOURCE_VERSION/VS0_TARGET_EPOCH_STALE (ADH-2026-060)")
        if f16_89.get("expectedError") != "CONFLICT" or f16_89.get("expectedViolation") != "VS0_TARGET_MAINTENANCE":
            e("MAINTENANCE-RACE: VS0-CF-F16-89 must be 409 CONFLICT/VS0_TARGET_MAINTENANCE (ADH-2026-060)")
    if f16_arch and "adh-2026-060" not in f16_arch.lower():
        e("MAINTENANCE-RACE: F0016 architecture must reference ADH-2026-060's mutually exclusive ordering")


def check_maintenance_entry_link_clearing(f16_arch: str, f16_feat: str) -> None:
    """ADH-2026-061: every successful Maintenance entry must clear both
    current FactSet/Result links and persist Active/Unqualified, whether or
    not a qualification is in flight. Reject wording that preserves a
    completed conclusion or current links on Maintenance entry, and reject
    an over-broad claim that all reservations for a target are aborted."""
    retained_conclusion_markers = (
        "preserves the last completed",
        "preserve the last completed",
        "retains the last completed",
        "retain the last completed",
        "preserves the last committed conclusion",
        "otherwise it preserves",
    )
    over_broad_abort_markers = (
        "aborts all reservations",
        "abort all reservations",
        "aborts every reservation",
        "abort every reservation",
        "all reservations for the target",
        "all reservations for that target",
        "all target reservations",
    )
    for label, text in (("architecture", f16_arch), ("feature", f16_feat)):
        if not text:
            continue
        lower = text.lower()
        for marker in retained_conclusion_markers:
            if marker in lower:
                e(f"MAINTENANCE-LINKS: F0016 {label} retains prohibited retained-conclusion wording {marker!r} (ADH-2026-061)")
        for marker in over_broad_abort_markers:
            if marker in lower:
                e(f"MAINTENANCE-LINKS: F0016 {label} contains an over-broad abort-scope claim {marker!r} (ADH-2026-061)")
    # The unconditional link-clearing rule is normative only in the architecture
    # authority (ADH-2026-061's write allowlist does not include the tabular
    # docs/features/FEATURE-0016 ID-summary file).
    if f16_arch and "clears both current" not in f16_arch.lower() and "clear both current" not in f16_arch.lower():
        e("MAINTENANCE-LINKS: F0016 architecture must state every successful Maintenance entry clears both current FactSet/Result links (ADH-2026-061)")


def check_adh_063_create_precedence(reg: dict, f16_arch: str, steer: str) -> None:
    """ADH-2026-063: the ExecutionTarget collection-create phase-one and
    strict-classification precedence correction adds the four new sequential
    local conformance rows VS0-CF-F16-123..126. ADH-2026-065 then makes the
    classification and phase-one reference cases exact: it refines
    VS0-CF-F16-125 to malformed-JSON only, adds VS0-CF-F16-127 (duplicate
    top-level member only) and VS0-CF-F16-128 (fail-closed missing/unextractable
    required phase-one reference denial). Fail closed if any row is missing or
    deviates, if the authoritative text permits strict classification before
    authorization/safe access, if phase one permits body
    classification/digest/reservation, if the oversized body is not rejected
    first, if F16-125 names a duplicate case, if F16-127 names a malformed case,
    or if F16-128 does not pin the exact fail-closed extraction denial."""
    by_id = {c.get("id"): c for c in reg.get("conformance", [])}
    f123 = by_id.get("VS0-CF-F16-123")
    f124 = by_id.get("VS0-CF-F16-124")
    f125 = by_id.get("VS0-CF-F16-125")
    f126 = by_id.get("VS0-CF-F16-126")
    if not f123 or f123.get("expectedError") != "AUTHORIZATION_DENIED":
        e("ADH063: VS0-CF-F16-123 must exist and be AUTHORIZATION_DENIED (ADH-2026-063)")
    elif "never classifies, canonicalizes, digests, or reserves" not in str(f123.get("expectedSideEffects", "")).lower():
        e("ADH063: VS0-CF-F16-123 must state phase one never classifies/canonicalizes/digests/reserves the body (ADH-2026-063)")
    if not f124 or f124.get("expectedError") != "RESOURCE_NOT_FOUND" or f124.get("expectedViolation") != "VS0_AUTHORIZATION_SAFE_DENIAL":
        e("ADH063: VS0-CF-F16-124 must exist and be safe RESOURCE_NOT_FOUND + VS0_AUTHORIZATION_SAFE_DENIAL (ADH-2026-063)")
    elif "never classifies, canonicalizes, digests, or reserves" not in str(f124.get("expectedSideEffects", "")).lower():
        e("ADH063: VS0-CF-F16-124 must state phase one never classifies/canonicalizes/digests/reserves the body (ADH-2026-063)")
    if not f126 or f126.get("expectedError") != "REQUEST_TOO_LARGE":
        e("ADH063: VS0-CF-F16-126 must exist and be REQUEST_TOO_LARGE (ADH-2026-063)")
    elif "before phase-one extraction, authorization, or safe access" not in str(f126.get("expectedSideEffects", "")).lower():
        e("ADH063: VS0-CF-F16-126 must state REQUEST_TOO_LARGE before phase-one extraction, authorization, or safe access (ADH-2026-063)")

    # ADH-2026-065: split the strict-classification family into exact non-family
    # cases and add the fail-closed missing/unextractable phase-one reference
    # denial. F16-125 is malformed-JSON only; F16-127 is duplicate-top-level
    # member only; F16-128 is the phase-one extraction denial.
    f127 = by_id.get("VS0-CF-F16-127")
    f128 = by_id.get("VS0-CF-F16-128")
    if not f125 or f125.get("expectedError") != "MALFORMED_REQUEST":
        e("ADH065: VS0-CF-F16-125 must exist and be MALFORMED_REQUEST (exact non-family malformed-JSON case) (ADH-2026-065)")
    else:
        f125_fx = str(f125.get("expectedSideEffects", "")).lower()
        if "duplicate" in f125_fx:
            e("ADH065: VS0-CF-F16-125 must not name a duplicate case; it is the exact malformed-JSON classification case (ADH-2026-065)")
        elif "malformed" not in f125_fx or "malformed_request" not in f125_fx or "strict single classification" not in f125_fx:
            e("ADH065: VS0-CF-F16-125 must state a single strict classification returning MALFORMED_REQUEST for malformed JSON (ADH-2026-065)")
    if not f127 or f127.get("expectedError") != "DUPLICATE_FIELD":
        e("ADH065: VS0-CF-F16-127 must exist and be DUPLICATE_FIELD (exact non-family duplicate-top-level-member case) (ADH-2026-065)")
    else:
        f127_fx = str(f127.get("expectedSideEffects", "")).lower()
        if "malformed" in f127_fx:
            e("ADH065: VS0-CF-F16-127 must not name a malformed case; it is the exact duplicate-top-level-member classification case (ADH-2026-065)")
        elif "duplicate" not in f127_fx or "duplicate_field" not in f127_fx or "strict single classification" not in f127_fx:
            e("ADH065: VS0-CF-F16-127 must state a single strict classification returning DUPLICATE_FIELD for a duplicate top-level member (ADH-2026-065)")
    if not f128 or f128.get("expectedError") != "AUTHORIZATION_DENIED":
        e("ADH065: VS0-CF-F16-128 must exist and be the existing audited AUTHORIZATION_DENIED phase-one extraction denial (ADH-2026-065)")
    else:
        f128_fx = (str(f128.get("inputs", "")) + " " + str(f128.get("expectedSideEffects", ""))).lower()
        if not ("exactly one syntactically usable uid" in f128_fx and "both" in f128_fx
                and "spec.cloudproviderparticipationref.uid" in f128_fx and "spec.infrastructurestackref.uid" in f128_fx):
            e("ADH065: VS0-CF-F16-128 must state phase one cannot extract exactly one syntactically usable UID at both required reference paths (ADH-2026-065)")
        elif not ("audited" in f128_fx and "403" in f128_fx
                  and "no body-classification detail" in f128_fx and "no backing-resource existence" in f128_fx):
            e("ADH065: VS0-CF-F16-128 must state the existing audited AUTHORIZATION_DENIED/403 with no body/backing disclosure (ADH-2026-065)")
        elif not all(tok in f128_fx for tok in ("strict classification", "canonicalization", "digest",
                     "reservation", "observer", "mutation", "publication", "completion")):
            e("ADH065: VS0-CF-F16-128 must state no strict classification/canonicalization/digest/reservation/observer/mutation/publication/completion (ADH-2026-065)")
    if f16_arch:
        lower = f16_arch.lower()
        if "adh-2026-063" not in lower:
            e("ADH063: F0016 architecture must reference ADH-2026-063 as a controlling correction")
        if "before phase-one extraction, authorization, or safe access" not in lower:
            e("ADH063: F0016 architecture must reject an oversized body before phase-one extraction, authorization, or safe access (ADH-2026-063)")
        if "never classifies, canonicalizes, digests, or reserves a body" not in lower:
            e("ADH063: F0016 architecture must state phase one never classifies/canonicalizes/digests/reserves a body (ADH-2026-063)")
        if "only then are the retained same bytes strictly classified" not in lower:
            e("ADH063: F0016 architecture must state authorization and safe access precede the single strict classification (ADH-2026-063)")
        if "adh-2026-065" not in lower:
            e("ADH065: F0016 architecture must reference ADH-2026-065 as a controlling correction")
    if steer and "adh-2026-063" not in steer.lower():
        e("ADH063: .kiro/steering/slice0-contract.md must record the ADH-2026-063 safe precedence invariant")
    if steer and "adh-2026-065" not in steer.lower():
        e("ADH065: .kiro/steering/slice0-contract.md must record the ADH-2026-065 exact classification/phase-one-reference invariant")


def check_conformance_completeness(reg: dict) -> None:
    confs = reg.get("conformance", [])
    found = {c.get("id") for c in confs}
    missing = [cid for cid in F16_REQUIRED_CF_IDS if cid not in found]
    if missing:
        e(f"CONFORMANCE: missing {len(missing)} required VS0-CF-F16 case(s): {missing[:10]}{'...' if len(missing) > 10 else ''}")
    f16_cases = [c for c in confs if c.get("id", "").startswith("VS0-CF-F16-")]
    for c in f16_cases:
        for field in ("id", "owner", "inputs", "expectedState", "expectedError", "expectedSideEffects", "gate"):
            if field not in c:
                e(f"CONFORMANCE: {c.get('id', '?')} missing required field {field!r}")
        if c.get("owner") != "FEATURE-0016":
            e(f"CONFORMANCE: {c.get('id')} owner must be FEATURE-0016, found {c.get('owner')!r}")


def check_vs000_executiontarget_ownership(charter: str) -> None:
    """ADH-2026-064: the VS-000 charter must not reclaim ExecutionTarget
    for FEATURE-0015; FEATURE-0016 consumes only the F0015 participation and
    stack contracts read-only."""
    if not charter:
        return
    cloud_row = next((line for line in charter.splitlines() if line.startswith("| Cloud model |")), "")
    expected = ("| Cloud model | CloudPlatform, CloudProvider, CloudProviderParticipation, "
                "HostingLocation, Datacenter, FaultDomain, InfrastructureStack | FEATURE-0015 |")
    if cloud_row != expected:
        e("OWNERSHIP: VS-000 Cloud model row must assign only CloudPlatform through InfrastructureStack to FEATURE-0015 (ADH-2026-064)")
    if "ExecutionTarget" in cloud_row:
        e("OWNERSHIP: VS-000 Cloud model row must never assign ExecutionTarget to FEATURE-0015 (ADH-2026-064)")
    integration_row = next((line for line in charter.splitlines() if line.startswith("| Integration |")), "")
    for marker in ("ExecutionTarget", "normalized target facts", "qualification", "synthetic observer boundary", "CloudProviderParticipation and InfrastructureStack read-only"):
        if marker not in integration_row or not integration_row.endswith("| FEATURE-0016 |"):
            e("OWNERSHIP: VS-000 Integration row must assign " + repr(marker) + " to FEATURE-0016 (ADH-2026-064)")


def check_traceability(trace: str) -> None:
    if not trace:
        return
    missing = [cid for cid in F16_REQUIRED_CF_IDS if cid not in trace]
    if missing:
        e(f"TRACEABILITY: traceability matrix missing {len(missing)} VS0-CF-F16 reference(s): {missing[:10]}{'...' if len(missing) > 10 else ''}")
    for schema_id in ("VS0-SCHEMA-015", "VS0-SCHEMA-016", "VS0-SCHEMA-017", "VS0-STATE-004", "VS0-WRITER-006"):
        if schema_id not in trace:
            e(f"TRACEABILITY: traceability matrix missing {schema_id}")
    if "ADH-2026-066" not in trace:
        e("TRACEABILITY: traceability matrix must record the ADH-2026-066 backing-access bridge")


def check_adh_066_backing_access(f16_arch: str, f16_feat: str, closure: str) -> None:
    if not ADH_066.exists():
        e(f"ADH066: handoff file missing: {ADH_066.relative_to(ROOT)}")
        return
    handoff = ADH_066.read_text()
    if "Approval status: Approved" not in handoff:
        e("ADH066: handoff must record 'Approval status: Approved'")
    normalized_arch = " ".join(f16_arch.split()).lower()
    required_arch_markers = (
        "ADH-2026-066", "BackingAccessProvider", "private FEATURE-0015 store-backed read lease",
        "InfrastructureStack-generation and viability-fingerprint fences", "participation generation is not a fence",
        "FEATURE-0015 backing lease, then the FEATURE-0016 lifecycle-service mutex",
        "prospective admission control",
    )
    for marker in required_arch_markers:
        if marker.lower() not in normalized_arch:
            e(f"ADH066: F0016 architecture missing backing-access marker: {marker}")
    if "ADH-2026-066" not in f16_feat or "BackingAccessProvider" not in f16_feat:
        e("ADH066: F0016 feature authority must record the F0016-owned backing-access bridge")
    if "ARC-F16-15" not in closure or "ADH-2026-066" not in closure:
        e("ADH066: closure matrix must record ARC-F16-15 as the resolved backing-access bridge")


def check_closure_matrix(closure: str) -> None:
    if not closure:
        return
    table_rows = [line for line in closure.splitlines() if line.strip().startswith("| ARC-F16-")]
    for i in range(1, 11):
        arc_id = f"ARC-F16-{i:02d}"
        row = next((r for r in table_rows if r.strip().startswith(f"| {arc_id} |")), None)
        if not row:
            e(f"CLOSURE: {arc_id} table row not found in closure matrix")
        elif "RESOLVED" not in row and "OUT OF SCOPE" not in row and "NOTED" not in row:
            e(f"CLOSURE: {arc_id} is not RESOLVED or explicitly noted out-of-scope")


def check_control_manifest_exists() -> None:
    if not CONTROL.exists():
        e(f"CONTROL: manifest missing: {CONTROL.relative_to(ROOT)}")
        return
    import json
    try:
        data = json.loads(CONTROL.read_text())
    except Exception as exc:
        e(f"CONTROL: manifest invalid JSON: {exc}")
        return
    owned = set(data.get("ownership", {}).get("owned_resources", []))
    expected = {"ExecutionTarget", "NormalizedTargetFactSet", "TargetQualificationResult"}
    if owned != expected:
        e(f"CONTROL: owned_resources mismatch: {sorted(owned ^ expected)}")
    # ADH-2026-063: the control manifest must record the ADH-2026-063 handoff.
    handoffs = data.get("feature", {}).get("handoffs", [])
    if not any("ADH-2026-063" in str(h) for h in handoffs):
        e("CONTROL: feature.handoffs must include the ADH-2026-063 handoff path (ADH-2026-063)")
    # ADH-2026-065: the control manifest must record the ADH-2026-065 handoff.
    if not any("ADH-2026-065" in str(h) for h in handoffs):
        e("CONTROL: feature.handoffs must include the ADH-2026-065 handoff path (ADH-2026-065)")
    if not any("ADH-2026-066" in str(h) for h in handoffs):
        e("CONTROL: feature.handoffs must include the ADH-2026-066 handoff path (ADH-2026-066)")


def check_reuse_assessment_exists() -> None:
    if not REUSE.exists():
        e(f"REUSE: reuse assessment missing: {REUSE.relative_to(ROOT)}")
        return
    text = REUSE.read_text()
    if "ADH-2026-058" not in text:
        e("REUSE: reuse assessment must reference ADH-2026-058")
    if "Approved" not in text:
        e("REUSE: reuse assessment must record Approved status")


def check_steering_updated(steer: str) -> None:
    if not steer:
        return
    if "ADH-2026-058" not in steer:
        e("STEERING: .kiro/steering/slice0-contract.md must reference ADH-2026-058")
    if "ExecutionTargetLifecycleService" not in steer:
        e("STEERING: .kiro/steering/slice0-contract.md must name ExecutionTargetLifecycleService as the sole F0016 status/record writer")
    if "sovrunn.synthetic-iaas-observer/v1" not in steer:
        e("STEERING: .kiro/steering/slice0-contract.md must name the sole observer sovrunn.synthetic-iaas-observer/v1")


def check_adh_058_approved() -> None:
    if not ADH_058.exists():
        e(f"ADH: handoff file missing: {ADH_058.relative_to(ROOT)}")
        return
    text = ADH_058.read_text()
    if "Approval status: Approved" not in text:
        e("ADH: ADH-2026-058 must record 'Approval status: Approved'")


def main() -> None:
    if not REG_PATH.exists():
        print(f"FAIL: registry not found: {REG_PATH.relative_to(ROOT)}")
        sys.exit(1)
    reg = yaml.safe_load(REG_PATH.read_text())
    f16_arch = read(F16_ARCH)
    f16_feat = read(F16_FEAT)
    trace = read(TRACE_PATH)
    closure = read(CLOSURE)
    steer = read(STEER)
    charter = read(CHARTER)

    check_adh_058_approved()
    check_no_reintroduced_placeholders(reg, f16_arch, f16_feat)
    check_five_route_surface(f16_arch, f16_feat)
    check_closed_violation_set(reg, f16_arch)
    check_sole_writer_and_observer(reg, f16_arch)
    check_maintenance_race_mutual_exclusivity(reg, f16_arch)
    check_maintenance_entry_link_clearing(f16_arch, f16_feat)
    check_adh_063_create_precedence(reg, f16_arch, steer)
    check_adh_066_backing_access(f16_arch, f16_feat, closure)
    check_conformance_completeness(reg)
    check_vs000_executiontarget_ownership(charter)
    check_traceability(trace)
    check_closure_matrix(closure)
    check_control_manifest_exists()
    check_reuse_assessment_exists()
    check_steering_updated(steer)

    if errs:
        print(f"FAIL: {len(errs)} error(s)")
        for err in errs:
            print(f"  ✗ {err}")
        sys.exit(1)
    print(
        "PASS: FEATURE-0016 architecture readiness — ADH-2026-058 placeholder replacement, "
        "five-route surface, sole writer/observer, closed violation set, 128-case conformance "
        "(incl. ADH-2026-063 create phase-one/strict-classification precedence rows VS0-CF-F16-123..126 "
        "and ADH-2026-065 exact classification/phase-one-reference rows VS0-CF-F16-125/127/128; "
        "ADH-2026-066 coherent principal-aware backing access), "
        "traceability, closure matrix, control manifest, reuse assessment, and steering are consistent"
    )


if __name__ == "__main__":
    main()
