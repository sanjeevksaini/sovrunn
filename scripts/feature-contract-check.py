#!/usr/bin/env python3
"""Fail-closed feature contract coverage checker.

This checker is intentionally independent of generated Kiro artifacts.  It
checks whether a feature's approved registry contract contains enough exact,
local evidence to generate requirements, design, and tasks without choosing
observable behavior in a later stage.

The framework accepts a feature argument so future features can add their own
route catalog.  FEATURE-0015 is the first enforced catalog.
"""

from __future__ import annotations

import argparse
import sys
from pathlib import Path

try:
    import yaml
except ImportError:
    print("FAIL: PyYAML required. Install: pip install pyyaml", file=sys.stderr)
    sys.exit(1)


ROOT = Path(__file__).resolve().parents[1]
REGISTRY = ROOT / "docs/architecture/vertical-slices/VS-000-contract-registry.yaml"
F15_ARCH = ROOT / "docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md"
F15_FEATURE = ROOT / "docs/features/FEATURE-0015-canonical-cloud-model-foundation.md"
F16_ARCH = ROOT / "docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md"
F16_FEATURE = ROOT / "docs/features/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md"

F15_KINDS = (
    "CloudPlatform", "CloudProvider", "CloudProviderParticipation", "HostingLocation",
    "Datacenter", "FaultDomain", "InfrastructureStack",
)
F15_ACTIONS = (
    "accept", "reject", "withdraw", "suspend", "resume",
    "request-release", "accept-release", "decline-release",
)
# Seven collection paths + seven item paths + eight action paths.  Participation
# has no PATCH surface, hence Go 1.22 needs 14 + 13 + 8 = 35 method patterns.
F15_LOGICAL_PATHS = 22
F15_METHOD_PATTERNS = 35

F16_REQUIRED_CF_IDS = [f"VS0-CF-F16-{i:02d}" for i in range(1, 127)]
F16_ROUTES = (
    ("POST", "/apis/execution.sovrunn.io/v1alpha1/execution-targets"),
    ("GET", "/apis/execution.sovrunn.io/v1alpha1/execution-targets"),
    ("GET", "/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}"),
    ("POST", "/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify"),
    ("POST", "/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire"),
)
F16_CLOSED_VIOLATIONS = (
    "VS0_TARGET_RETIRED", "VS0_TARGET_MAINTENANCE", "VS0_TARGET_QUALIFICATION_IN_PROGRESS",
    "VS0_TARGET_EPOCH_STALE", "VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE",
    "VS0_EXECUTION_TARGET_STACK_UNAVAILABLE", "VS0_EXECUTION_TARGET_SCOPE_MISMATCH",
    "VS0_EXECUTION_TARGET_VIABILITY_STALE",
)


def load_registry() -> dict:
    try:
        return yaml.safe_load(REGISTRY.read_text())
    except FileNotFoundError:
        raise SystemExit(f"FAIL: missing registry: {REGISTRY.relative_to(ROOT)}")


def by_id(registry: dict) -> dict[str, dict]:
    return {str(case.get("id")): case for case in registry.get("conformance", [])}


def has_audit_effect(case: dict) -> bool:
    return "auditevent" in str(case.get("expectedSideEffects", "")).lower()


def require(errors: list[str], condition: bool, message: str) -> None:
    if not condition:
        errors.append(message)


def check_f0015(registry: dict, errors: list[str]) -> None:
    cases = by_id(registry)

    writer_boundary = cases.get("VS0-CF-F15-41")
    require(errors, writer_boundary is not None,
            "WRITER_BOUNDARY: missing VS0-CF-F15-41 (ADH-2026-057)")
    if writer_boundary:
        contract = " ".join(str(writer_boundary.get(field, "")) for field in (
            "inputs", "expectedState", "expectedError", "expectedSideEffects", "gate",
        )).lower()
        for marker in (
            "cloudplatform spec.description", "cloudprovider spec.displayname/spec.operatingmarkets",
            "topology spec.description", "authorization_denied", "before semantic patch processing",
            "no resource mutation", "no auditevent", "no idempotency record",
        ):
            require(errors, marker in contract,
                    f"WRITER_BOUNDARY: F15-41 must state {marker!r} (ADH-2026-057)")

    # Every collection create/action is covered by the already-approved
    # idempotency replay and changed-digest rules.
    for case_id in ("VS0-CF-F15-18", "VS0-CF-F15-19"):
        case = cases.get(case_id)
        require(errors, case is not None, f"IDEMPOTENCY: missing {case_id}")
        if not case:
            continue
        inputs = str(case.get("inputs", "")).lower()
        for kind in F15_KINDS:
            require(errors, kind.lower() in inputs,
                    f"IDEMPOTENCY: {case_id} must cover {kind} collection POST")
        for action in F15_ACTIONS:
            require(errors, action in inputs,
                    f"IDEMPOTENCY: {case_id} must cover participation action {action}")
    if cases.get("VS0-CF-F15-18"):
        require(errors, cases["VS0-CF-F15-18"].get("expectedError") is None,
                "IDEMPOTENCY: F15-18 replay must preserve null expectedError")
    if cases.get("VS0-CF-F15-19"):
        require(errors, cases["VS0-CF-F15-19"].get("expectedError") == "CONFLICT",
                "IDEMPOTENCY: F15-19 changed-digest replay must return CONFLICT")
        require(errors,
                cases["VS0-CF-F15-19"].get("expectedViolation") == "VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH",
                "IDEMPOTENCY: F15-19 must use VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH")

    # ADH-2026-054 closes the cross-target/revoked-grant replay gap.  The
    # same all-route cases must also state the concrete item target and the
    # current server-side authorization/safe-access checks that precede a
    # completed replay.
    for case_id in ("VS0-CF-F15-18", "VS0-CF-F15-19"):
        case = cases.get(case_id)
        if not case:
            continue
        inputs = str(case.get("inputs", "")).lower()
        for marker in (
            "concrete cloudproviderparticipation uid",
            "current server-resolved authorization",
            "safe target/reference access",
        ):
            require(errors, marker in inputs,
                    f"IDEMPOTENCY: {case_id} must state {marker!r} (ADH-2026-054)")

    forged = cases.get("VS0-CF-F15-22")
    require(errors, forged is not None,
            "FORGED_GRANT: missing VS0-CF-F15-22 (ADH-2026-054)")
    if forged:
        inputs = str(forged.get("inputs", "")).lower()
        effects = str(forged.get("expectedSideEffects", "")).lower()
        for marker in ("x-sovrunn-bootstrap-grant", "bootstrapgrant", "syntactically valid duplicate-free"):
            require(errors, marker in inputs,
                    f"FORGED_GRANT: F15-22 must state {marker!r} carrier semantics (ADH-2026-054)")
        for marker in ("header detection follows authentication", "body-member detection precedes current authorization"):
            require(errors, marker in effects,
                    f"FORGED_GRANT: F15-22 must state {marker!r} precedence (ADH-2026-054)")

    malformed = cases.get("VS0-CF-F15-30")
    require(errors, malformed is not None,
            "FORGED_GRANT: missing VS0-CF-F15-30 (ADH-2026-054)")
    if malformed:
        inputs = str(malformed.get("inputs", "")).lower()
        effects = str(malformed.get("expectedSideEffects", "")).lower()
        require(errors, "does not contain the reserved top-level bootstrapgrant" in inputs,
                "FORGED_GRANT: F15-30 must exclude the reserved body carrier (ADH-2026-054)")
        require(errors, "governed instead by vs0-cf-f15-22" in effects,
                "FORGED_GRANT: F15-30 must delegate the reserved body carrier to F15-22 (ADH-2026-054)")

    # The active F0015 audit rule says authorization and safe denials produce
    # evidence.  Every local case expressing one of those outcomes must carry
    # the same side effect; otherwise no downstream stage can know the rule.
    audit_rule = cases.get("VS0-CF-F15-24")
    require(errors, audit_rule is not None, "AUDIT: missing VS0-CF-F15-24")
    if audit_rule:
        inputs = str(audit_rule.get("inputs", "")).lower()
        requires_denial_audit = "authorization denial" in inputs and "safe denial" in inputs
        if requires_denial_audit:
            audited_denials = (
                "VS0-CF-F15-05", "VS0-CF-F15-06", "VS0-CF-F15-21",
                "VS0-CF-F15-22", "VS0-CF-F15-23", "VS0-CF-X03",
            )
            for case_id in audited_denials:
                case = cases.get(case_id)
                require(errors, case is not None, f"AUDIT: missing {case_id}")
                if case:
                    require(errors, has_audit_effect(case),
                            f"AUDIT: {case_id} is an authorization/safe denial but lacks required AuditEvent side effect")

    # Authentication is an explicitly required F0015 route outcome and must
    # have local conformance; inherited FEATURE-0018 evidence cannot prove it.
    auth_case = cases.get("VS0-CF-F15-31")
    require(errors, auth_case is not None,
            "AUTHN: missing F0015-local VS0-CF-F15-31 for AUTH_REQUIRED on every owned route")
    if auth_case:
        require(errors, auth_case.get("owner") == "FEATURE-0015",
                "AUTHN: F15-31 must be owned by FEATURE-0015")
        require(errors, auth_case.get("expectedError") == "AUTH_REQUIRED",
                "AUTHN: F15-31 must map missing/invalid authentication to AUTH_REQUIRED")
        require(errors, "no mutation" in str(auth_case.get("expectedSideEffects", "")).lower(),
                "AUTHN: F15-31 must state no mutation side effect")

    # ADH-2026-052: CloudPlatform/CloudProvider name-uniqueness and general
    # registry-declared schema-constraint validation must each have exact local
    # proof; the older cases already covering scope-subset/root/PATCH/topology-
    # ordering scenarios cannot stand in for these two absent scenarios.
    name_uniqueness_case = cases.get("VS0-CF-F15-32")
    require(errors, name_uniqueness_case is not None,
            "VALIDATIONPROOF: missing F0015-local VS0-CF-F15-32 for CloudPlatform/CloudProvider name-uniqueness rejection")
    if name_uniqueness_case:
        require(errors, name_uniqueness_case.get("owner") == "FEATURE-0015",
                "VALIDATIONPROOF: F15-32 must be owned by FEATURE-0015")
        require(errors, name_uniqueness_case.get("expectedError") == "ALREADY_EXISTS",
                "VALIDATIONPROOF: F15-32 must map a duplicate metadata.name to ALREADY_EXISTS")

    schema_constraint_case = cases.get("VS0-CF-F15-33")
    require(errors, schema_constraint_case is not None,
            "VALIDATIONPROOF: missing F0015-local VS0-CF-F15-33 for registry-declared schema-constraint validation")
    if schema_constraint_case:
        require(errors, schema_constraint_case.get("owner") == "FEATURE-0015",
                "VALIDATIONPROOF: F15-33 must be owned by FEATURE-0015")
        require(errors, schema_constraint_case.get("expectedError") == "VALIDATION_FAILED",
                "VALIDATIONPROOF: F15-33 must map a registry-declared schema-constraint violation to VALIDATION_FAILED")
        f33_inputs = str(schema_constraint_case.get("inputs", "")).lower()
        for excluded in ("duplicate name", "unassigned", "malformed"):
            require(errors, excluded in f33_inputs,
                    f"VALIDATIONPROOF: F15-33 inputs must explicitly exclude '{excluded}' outcomes covered by other cases")

    # Route terminology must distinguish externally visible endpoint paths from
    # Go 1.22 method-qualified ServeMux registrations.  This catches impossible
    # registration arithmetic before requirements/design generation.
    authority_text = F15_ARCH.read_text() + "\n" + F15_FEATURE.read_text()
    logical = f"{F15_LOGICAL_PATHS} logical"
    method_patterns = f"{F15_METHOD_PATTERNS} explicit"
    require(errors, logical in authority_text.lower(),
            f"ROUTING: authority must state {F15_LOGICAL_PATHS} logical endpoint paths")
    require(errors, method_patterns in authority_text.lower(),
            f"ROUTING: authority must state {F15_METHOD_PATTERNS} explicit Go 1.22 method/path registrations")
    require(errors, "no path-only" in authority_text.lower() or "no internal method" in authority_text.lower(),
            "ROUTING: authority must prohibit path-only internal method dispatch")


def check_f0016(registry: dict, errors: list[str]) -> None:
    """FEATURE-0016 route catalog per ADH-2026-058.

    This checker validates only the route/conformance/violation closure that
    is this checker's concern; the full architecture-readiness surface
    (placeholder removal, sole writer/observer, closure matrix, control
    manifest, reuse assessment, steering) is separately and more completely
    validated by scripts/feature-0016-architecture-readiness-check.py.
    """
    cases = by_id(registry)

    missing_cf = [cid for cid in F16_REQUIRED_CF_IDS if cid not in cases]
    require(errors, not missing_cf,
            f"CONFORMANCE: missing {len(missing_cf)} required VS0-CF-F16 case(s): {missing_cf[:10]}")
    for case in cases.values():
        if not str(case.get("id", "")).startswith("VS0-CF-F16-"):
            continue
        require(errors, case.get("owner") == "FEATURE-0016",
                f"CONFORMANCE: {case.get('id')} owner must be FEATURE-0016, found {case.get('owner')!r}")

    if not F16_ARCH.exists():
        require(errors, False, f"ROUTES: missing architecture authority {F16_ARCH.relative_to(ROOT)}")
        return
    arch_lower = F16_ARCH.read_text().lower()
    for method, path in F16_ROUTES:
        require(errors, path.lower() in arch_lower,
                f"ROUTES: F0016 architecture must state route {method} {path} (ADH-2026-058 clause 3)")
    require(errors, "no patch" in arch_lower,
            "ROUTES: F0016 architecture must state no PATCH/PUT/DELETE/HEAD public route exists (ADH-2026-058 clause 3)")
    require(errors, "pre-servemux" in arch_lower,
            "ROUTES: F0016 architecture must describe the pre-ServeMux transport-only method/path guard (ADH-2026-058 clause 3)")

    vcs = {vc.get("code") for vc in registry.get("violationCodes", {}).get("slice0", [])}
    for code in F16_CLOSED_VIOLATIONS:
        require(errors, code in vcs, f"VIOLATIONS: {code} must be a registered violation code (ADH-2026-058 clause 9)")

    if F16_FEATURE.exists():
        feat_lower = F16_FEATURE.read_text().lower()
        require(errors, "exactly five" in feat_lower or "five routes" in feat_lower or "five explicit" in feat_lower,
                "ROUTES: F0016 feature must state the exact five-route closure (ADH-2026-058 clause 3)")

    # ADH-2026-060: F16-75 (inactive-marker epoch-stale) and F16-89 (active-marker
    # Maintenance-wins) must be mutually exclusive and must not be reversed.
    f16_75 = cases.get("VS0-CF-F16-75")
    f16_89 = cases.get("VS0-CF-F16-89")
    require(errors, f16_75 is not None, "MAINTENANCE-RACE: VS0-CF-F16-75 not found in registry (ADH-2026-060)")
    require(errors, f16_89 is not None, "MAINTENANCE-RACE: VS0-CF-F16-89 not found in registry (ADH-2026-060)")
    if f16_75 is not None and f16_89 is not None:
        f75_inputs = str(f16_75.get("inputs", "")).lower()
        f89_inputs = str(f16_89.get("inputs", "")).lower()
        require(errors, "maintenance entry wins" not in f75_inputs and "maintenance-entry wins" not in f75_inputs,
                "MAINTENANCE-RACE: VS0-CF-F16-75 must not say Maintenance entry wins (ADH-2026-060); that is the F16-89 case")
        require(errors, "no active current-maintenance marker" in f75_inputs or "no active current-maintenance marker exists" in f75_inputs,
                "MAINTENANCE-RACE: VS0-CF-F16-75 input must state no active current-Maintenance marker exists (ADH-2026-060)")
        require(errors, "maintenance entry wins" in f89_inputs,
                "MAINTENANCE-RACE: VS0-CF-F16-89 input must state Maintenance entry wins (ADH-2026-060)")
        require(errors, f16_75.get("expectedError") == "STALE_RESOURCE_VERSION",
                f"MAINTENANCE-RACE: VS0-CF-F16-75 expectedError must be STALE_RESOURCE_VERSION, found {f16_75.get('expectedError')!r} (ADH-2026-060)")
        require(errors, f16_75.get("expectedViolation") == "VS0_TARGET_EPOCH_STALE",
                f"MAINTENANCE-RACE: VS0-CF-F16-75 expectedViolation must be VS0_TARGET_EPOCH_STALE, found {f16_75.get('expectedViolation')!r} (ADH-2026-060)")
        require(errors, f16_89.get("expectedError") == "CONFLICT",
                f"MAINTENANCE-RACE: VS0-CF-F16-89 expectedError must be CONFLICT, found {f16_89.get('expectedError')!r} (ADH-2026-060)")
        require(errors, f16_89.get("expectedViolation") == "VS0_TARGET_MAINTENANCE",
                f"MAINTENANCE-RACE: VS0-CF-F16-89 expectedViolation must be VS0_TARGET_MAINTENANCE, found {f16_89.get('expectedViolation')!r} (ADH-2026-060)")

    # ADH-2026-061: F16-42 (Maintenance entry clears links), F16-43 (Maintenance
    # clear), and F16-89 (Maintenance-wins during in-flight qualification) must
    # keep their exact pre-existing observable semantics; the registry itself is
    # not the source of the fixed contradiction, so this checker only verifies
    # it remains unaltered by this correction.
    f16_42 = cases.get("VS0-CF-F16-42")
    f16_43 = cases.get("VS0-CF-F16-43")
    require(errors, f16_42 is not None, "MAINTENANCE-LINKS: VS0-CF-F16-42 not found in registry (ADH-2026-061)")
    require(errors, f16_43 is not None, "MAINTENANCE-LINKS: VS0-CF-F16-43 not found in registry (ADH-2026-061)")
    if f16_42 is not None:
        f42_effects = str(f16_42.get("expectedSideEffects", "")).lower()
        require(errors, "current factset/result links clear" in f42_effects,
                "MAINTENANCE-LINKS: VS0-CF-F16-42 must state current FactSet/Result links clear (ADH-2026-061)")
    if f16_43 is not None:
        f43_effects = str(f16_43.get("expectedSideEffects", "")).lower()
        require(errors, "current factset/result links clear" in f43_effects,
                "MAINTENANCE-LINKS: VS0-CF-F16-43 must state current FactSet/Result links clear (ADH-2026-061)")
    if f16_89 is not None:
        f89_effects = str(f16_89.get("expectedSideEffects", "")).lower()
        require(errors, "no qualification result, completion, or qualification auditevent" in f89_effects,
                "MAINTENANCE-LINKS: VS0-CF-F16-89 must keep its no-result/no-completion/no-AuditEvent outcome unchanged (ADH-2026-061)")

    # ADH-2026-063: the four new collection-create phase-one/strict-classification
    # rows must exist and must not deviate; the architecture authority must pin
    # the ordered precedence (oversized first; phase one never classifies/digests/
    # reserves; authorization/safe access precede strict classification).
    f16_123 = cases.get("VS0-CF-F16-123")
    f16_124 = cases.get("VS0-CF-F16-124")
    f16_125 = cases.get("VS0-CF-F16-125")
    f16_126 = cases.get("VS0-CF-F16-126")
    require(errors, f16_123 is not None and f16_123.get("expectedError") == "AUTHORIZATION_DENIED",
            "ADH063: VS0-CF-F16-123 must exist and be AUTHORIZATION_DENIED with no body classification (ADH-2026-063)")
    if f16_123 is not None:
        require(errors, "never classifies, canonicalizes, digests, or reserves" in str(f16_123.get("expectedSideEffects", "")).lower(),
                "ADH063: VS0-CF-F16-123 must state phase one never classifies/canonicalizes/digests/reserves the body (ADH-2026-063)")
    require(errors, f16_124 is not None and f16_124.get("expectedError") == "RESOURCE_NOT_FOUND"
            and (f16_124 or {}).get("expectedViolation") == "VS0_AUTHORIZATION_SAFE_DENIAL",
            "ADH063: VS0-CF-F16-124 must exist and be safe RESOURCE_NOT_FOUND + VS0_AUTHORIZATION_SAFE_DENIAL (ADH-2026-063)")
    if f16_124 is not None:
        require(errors, "never classifies, canonicalizes, digests, or reserves" in str(f16_124.get("expectedSideEffects", "")).lower(),
                "ADH063: VS0-CF-F16-124 must state phase one never classifies/canonicalizes/digests/reserves the body (ADH-2026-063)")
    require(errors, f16_125 is not None and f16_125.get("expectedError") == "MALFORMED_REQUEST",
            "ADH063: VS0-CF-F16-125 must exist and be MALFORMED_REQUEST (with DUPLICATE_FIELD as the applicable family member) (ADH-2026-063)")
    if f16_125 is not None:
        f125_effects = str(f16_125.get("expectedSideEffects", "")).lower()
        require(errors, "strict single classification" in f125_effects and "duplicate_field" in f125_effects,
                "ADH063: VS0-CF-F16-125 must state a single strict classification returning MALFORMED_REQUEST or DUPLICATE_FIELD (ADH-2026-063)")
    require(errors, f16_126 is not None and f16_126.get("expectedError") == "REQUEST_TOO_LARGE",
            "ADH063: VS0-CF-F16-126 must exist and be REQUEST_TOO_LARGE (ADH-2026-063)")
    if f16_126 is not None:
        require(errors, "before phase-one extraction, authorization, or safe access" in str(f16_126.get("expectedSideEffects", "")).lower(),
                "ADH063: VS0-CF-F16-126 must state REQUEST_TOO_LARGE before phase-one extraction, authorization, or safe access (ADH-2026-063)")
    require(errors, "adh-2026-063" in arch_lower,
            "ADH063: F0016 architecture must reference ADH-2026-063 as a controlling correction")
    require(errors, "before phase-one extraction, authorization, or safe access" in arch_lower,
            "ADH063: F0016 architecture must state the oversized body is rejected before phase-one extraction, authorization, or safe access (ADH-2026-063)")
    require(errors, "never classifies, canonicalizes, digests, or reserves a body" in arch_lower,
            "ADH063: F0016 architecture must state phase one never classifies/canonicalizes/digests/reserves a body (ADH-2026-063)")
    require(errors, "only then are the retained same bytes strictly classified" in arch_lower,
            "ADH063: F0016 architecture must state authorization and safe access precede the single strict classification (ADH-2026-063)")


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    args = parser.parse_args()
    if args.feature not in ("FEATURE-0015", "FEATURE-0016"):
        raise SystemExit(f"FAIL: no route catalog is configured for {args.feature}")

    errors: list[str] = []
    registry = load_registry()
    if args.feature == "FEATURE-0015":
        check_f0015(registry, errors)
    else:
        check_f0016(registry, errors)
    if errors:
        print(f"FAIL: {args.feature} feature-contract-check — {len(errors)} error(s)")
        for error in errors:
            print(f"  ✗ {error}")
        print("Kiro requirements/design/tasks generation is BLOCKED until the feature contract closes.")
        raise SystemExit(1)
    if args.feature == "FEATURE-0015":
        print("PASS: FEATURE-0015 feature contract closes route, audit, idempotency, authentication, validation-proof, and registration evidence")
    else:
        print("PASS: FEATURE-0016 feature contract closes route, violation-code, conformance-catalog (VS0-CF-F16-01..126), and ADH-2026-063 create phase-one/strict-classification precedence evidence")


if __name__ == "__main__":
    main()
