#!/usr/bin/env python3
"""Fail-closed feature contract coverage checker.

This checker is intentionally independent of generated Kiro artifacts.  It
checks whether a feature's approved registry contract contains enough exact,
local evidence to generate requirements, design, and tasks without choosing
observable behavior in a later stage.

The framework accepts a feature argument so each feature can register the
contract evidence appropriate to its boundary. Route-owning features validate
route catalogs; route-free in-process features validate their explicit
no-route contract instead.
"""

from __future__ import annotations

import argparse
import hashlib
import re
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
F17_ARCH = ROOT / "docs/architecture/policy-evaluation-abstraction.md"

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

F16_REQUIRED_CF_IDS = [f"VS0-CF-F16-{i:02d}" for i in range(1, 129)]
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

F17_DIGEST_PREIMAGE = (
    '{"action":"service.read","candidateRefs":[],"profileRefs":'
    '[{"apiVersion":"iam.sovrunn.io/v1alpha1","kind":"RoleDefinition",'
    '"name":"reader","uid":"role-001"}],"schema":'
    '"sovrunn.policy-evaluation-request/v1","subjectRef":'
    '{"apiVersion":"services.sovrunn.io/v1alpha1","kind":"ServiceInstance",'
    '"name":"reporting-api","uid":"service-001"}}'
)
F17_DIGEST = "daca15fd0310c46b45d5aff9bbe4f1a5dedd4b788c44c3780838a8be40a56103"


def load_registry() -> dict:
    try:
        return yaml.safe_load(REGISTRY.read_text())
    except FileNotFoundError:
        raise SystemExit(f"FAIL: missing registry: {REGISTRY.relative_to(ROOT)}")


def by_id(registry: dict) -> dict[str, dict]:
    return {str(case.get("id")): case for case in registry.get("conformance", [])}


def schemas_by_id(registry: dict) -> dict[str, dict]:
    return {str(schema.get("id")): schema for schema in registry.get("schemas", [])}


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
    require(errors, f16_126 is not None and f16_126.get("expectedError") == "REQUEST_TOO_LARGE",
            "ADH063: VS0-CF-F16-126 must exist and be REQUEST_TOO_LARGE (ADH-2026-063)")
    if f16_126 is not None:
        require(errors, "before phase-one extraction, authorization, or safe access" in str(f16_126.get("expectedSideEffects", "")).lower(),
                "ADH063: VS0-CF-F16-126 must state REQUEST_TOO_LARGE before phase-one extraction, authorization, or safe access (ADH-2026-063)")
    require(errors, "adh-2026-063" in arch_lower,
            "ADH063: F0016 architecture must reference ADH-2026-063 as a controlling correction")

    # ADH-2026-065: split the divergent strict-classification family into exact
    # non-family cases and add the fail-closed missing/unextractable phase-one
    # reference denial. F16-125 is malformed-JSON only; F16-127 is duplicate
    # top-level member only; F16-128 is the phase-one extraction denial.
    f16_127 = cases.get("VS0-CF-F16-127")
    f16_128 = cases.get("VS0-CF-F16-128")
    require(errors, f16_125 is not None and f16_125.get("expectedError") == "MALFORMED_REQUEST",
            "ADH065: VS0-CF-F16-125 must exist and be MALFORMED_REQUEST (exact non-family malformed-JSON case) (ADH-2026-065)")
    if f16_125 is not None:
        f125_effects = str(f16_125.get("expectedSideEffects", "")).lower()
        require(errors, "duplicate" not in f125_effects,
                "ADH065: VS0-CF-F16-125 must not name a duplicate case; it is the exact malformed-JSON classification case (ADH-2026-065)")
        require(errors, "malformed" in f125_effects and "malformed_request" in f125_effects and "strict single classification" in f125_effects,
                "ADH065: VS0-CF-F16-125 must state a single strict classification returning MALFORMED_REQUEST for malformed JSON (ADH-2026-065)")
    require(errors, f16_127 is not None and f16_127.get("expectedError") == "DUPLICATE_FIELD",
            "ADH065: VS0-CF-F16-127 must exist and be DUPLICATE_FIELD (exact non-family duplicate-top-level-member case) (ADH-2026-065)")
    if f16_127 is not None:
        f127_effects = str(f16_127.get("expectedSideEffects", "")).lower()
        require(errors, "malformed" not in f127_effects,
                "ADH065: VS0-CF-F16-127 must not name a malformed case; it is the exact duplicate-top-level-member classification case (ADH-2026-065)")
        require(errors, "duplicate" in f127_effects and "duplicate_field" in f127_effects and "strict single classification" in f127_effects,
                "ADH065: VS0-CF-F16-127 must state a single strict classification returning DUPLICATE_FIELD for a duplicate top-level member (ADH-2026-065)")
    require(errors, f16_128 is not None and f16_128.get("expectedError") == "AUTHORIZATION_DENIED",
            "ADH065: VS0-CF-F16-128 must exist and be the existing audited AUTHORIZATION_DENIED phase-one extraction denial (ADH-2026-065)")
    if f16_128 is not None:
        f128_text = (str(f16_128.get("inputs", "")) + " " + str(f16_128.get("expectedSideEffects", ""))).lower()
        require(errors, "exactly one syntactically usable uid" in f128_text and "both" in f128_text
                and "spec.cloudproviderparticipationref.uid" in f128_text and "spec.infrastructurestackref.uid" in f128_text,
                "ADH065: VS0-CF-F16-128 must state phase one cannot extract exactly one syntactically usable UID at both required reference paths (ADH-2026-065)")
        require(errors, "audited" in f128_text and "403" in f128_text
                and "no body-classification detail" in f128_text and "no backing-resource existence" in f128_text,
                "ADH065: VS0-CF-F16-128 must state the existing audited AUTHORIZATION_DENIED/403 with no body/backing disclosure (ADH-2026-065)")
        require(errors, all(tok in f128_text for tok in (
                "strict classification", "canonicalization", "digest", "reservation",
                "observer", "mutation", "publication", "completion")),
                "ADH065: VS0-CF-F16-128 must state no strict classification/canonicalization/digest/reservation/observer/mutation/publication/completion (ADH-2026-065)")
    require(errors, "adh-2026-065" in arch_lower,
            "ADH065: F0016 architecture must reference ADH-2026-065 as a controlling correction")
    require(errors, "before phase-one extraction, authorization, or safe access" in arch_lower,
            "ADH063: F0016 architecture must state the oversized body is rejected before phase-one extraction, authorization, or safe access (ADH-2026-063)")
    require(errors, "never classifies, canonicalizes, digests, or reserves a body" in arch_lower,
            "ADH063: F0016 architecture must state phase one never classifies/canonicalizes/digests/reserves a body (ADH-2026-063)")
    require(errors, "only then are the retained same bytes strictly classified" in arch_lower,
            "ADH063: F0016 architecture must state authorization and safe access precede the single strict classification (ADH-2026-063)")
    # ADH-2026-066 adds only an internal F0016 backing-access compatibility
    # bridge. It must not change the public contract, but its ownership and
    # registered-fence limits must remain explicit.
    normalized_arch = " ".join(arch_lower.split())
    for marker in (
        "adh-2026-066", "backingaccessprovider", "private feature-0015 store-backed read lease",
        "infrastructurestack-generation and viability-fingerprint fences",
        "participation generation is not a fence",
    ):
        require(errors, marker in normalized_arch,
                f"ADH066: F0016 architecture must state {marker!r} for the private backing-access bridge")


def check_f0017(registry: dict, errors: list[str]) -> None:
    """Validate the route-free FEATURE-0017 generation contract.

    FEATURE-0017 is an in-process transient seam. Its readiness proof is the
    exact request/result registry shape plus the sole architecture authority;
    a public route catalog would itself be architecture drift.
    """
    schemas = schemas_by_id(registry)
    request = schemas.get("VS0-SCHEMA-018")
    result = schemas.get("VS0-SCHEMA-019")

    require(errors, request is not None,
            "SCHEMA: missing VS0-SCHEMA-018 PolicyEvaluationRequest")
    if request is not None:
        require(errors, request.get("identity") == "policy.sovrunn.io/v1alpha1/PolicyEvaluationRequest",
                "SCHEMA: VS0-SCHEMA-018 identity must be PolicyEvaluationRequest")
        require(errors, request.get("owner") == "FEATURE-0017",
                "SCHEMA: VS0-SCHEMA-018 owner must be FEATURE-0017")
        require(errors, request.get("profile") == "TransientRequestResult",
                "SCHEMA: VS0-SCHEMA-018 must remain a TransientRequestResult")
        require(errors, request.get("boundary") == "internal-engine-facing",
                "SCHEMA: VS0-SCHEMA-018 must remain internal-engine-facing")
        require(errors, request.get("scopes") == ["Project", "CloudPlatform", "CloudProvider"],
                "SCHEMA: VS0-SCHEMA-018 scopes must be Project, CloudPlatform, CloudProvider")
        require(errors, request.get("required") == [
            "subjectRef:TypedRef(required)",
            "action:string(1..63)",
            "profileRefs:TypedRef[](1..32,uid-pinned)",
            "requestId:string(1..128)",
        ], "SCHEMA: VS0-SCHEMA-018 required fields/bounds do not match CDG-F17-01")
        require(errors, request.get("optional") == [
            "contextRef:TypedRef<EffectiveGovernanceContext>(uid-pinned)",
            "candidateRefs:TypedRef[](max64)",
        ], "SCHEMA: VS0-SCHEMA-018 contextRef/candidateRefs must be the exact optional fields")
        require(errors, request.get("retention") == "none",
                "SCHEMA: VS0-SCHEMA-018 retention must be none")

    require(errors, result is not None,
            "SCHEMA: missing VS0-SCHEMA-019 PolicyEvaluationResult")
    if result is not None:
        require(errors, result.get("identity") == "policy.sovrunn.io/v1alpha1/PolicyEvaluationResult",
                "SCHEMA: VS0-SCHEMA-019 identity must be PolicyEvaluationResult")
        require(errors, result.get("owner") == "FEATURE-0017",
                "SCHEMA: VS0-SCHEMA-019 owner must be FEATURE-0017")
        require(errors, result.get("profile") == "TransientRequestResult",
                "SCHEMA: VS0-SCHEMA-019 must remain a TransientRequestResult")
        require(errors, result.get("required") == [
            "outcome:enum[Allow,Deny,Indeterminate,RequiresApproval]",
            "reasonCodes:string[](1..32,unique,sorted,pattern=^[A-Z][A-Z0-9_]{0,62}$)",
            "inputDigest:string(64,lowercase-sha256-hex)",
            "evaluatedAt:RFC3339",
        ], "SCHEMA: VS0-SCHEMA-019 required fields/outcomes do not match CDG-F17-03")
        require(errors, result.get("optional") == ["obligations:object[](max32;absent-in-phase2r-fake)"],
                "SCHEMA: VS0-SCHEMA-019 must retain obligations as its sole optional field")
        require(errors, result.get("retention") == "decision-input",
                "SCHEMA: VS0-SCHEMA-019 retention must be decision-input")

    if not F17_ARCH.exists():
        require(errors, False, f"ARCHITECTURE: missing {F17_ARCH.relative_to(ROOT)}")
        return

    arch = F17_ARCH.read_text()
    arch_lower = " ".join(arch.lower().split())

    groups = re.findall(r"^### (CDG-F17-\d{2})\b", arch, flags=re.MULTILINE)
    require(errors, groups == [f"CDG-F17-{i:02d}" for i in range(1, 7)],
            f"ARCHITECTURE: expected exactly CDG-F17-01..06, found {groups}")

    calculated = hashlib.sha256(F17_DIGEST_PREIMAGE.encode("utf-8")).hexdigest()
    require(errors, calculated == F17_DIGEST,
            "DIGEST: checker fixture does not match the approved SHA-256 vector")
    require(errors, F17_DIGEST_PREIMAGE in arch,
            "DIGEST: architecture must contain the exact no-newline JCS preimage")
    require(errors, F17_DIGEST in arch,
            "DIGEST: architecture must contain the exact lowercase SHA-256 result")
    for marker in (
        "rfc 8785 json canonicalization scheme",
        "the adapter does not return a complete `policyevaluationresult`",
        "a normalized conclusion containing only `outcome` and `reasoncodes`",
        "perform exact immutable lookup by approved v1 input digest",
        "contain no conditional that interprets action, subject, profile, governance context",
        "no public route, store, controller, idempotency repository",
        "production adapter selection is explicitly outside feature-0017",
        "zero network, opa, cedar, kubernetes, cloudprovider, database, identity",
        "it neither receives nor resolves a `decisionprofile`",
        "later adopting domain validates the constructed `evaluationresult`",
    ):
        require(errors, marker in arch_lower,
                f"ARCHITECTURE: missing required FEATURE-0017 boundary marker {marker!r}")

    conformance_match = re.search(
        r"^## \d+\. Required local conformance\s*$([\s\S]*?)^## \d+\.",
        arch,
        flags=re.MULTILINE,
    )
    require(errors, conformance_match is not None,
            "CONFORMANCE: architecture must contain the local conformance section")
    if conformance_match is not None:
        numbers = [int(value) for value in re.findall(
            r"^(\d+)\.", conformance_match.group(1), flags=re.MULTILINE
        )]
        require(errors, numbers == list(range(1, 25)),
                f"CONFORMANCE: expected exact local inventory 1..24, found {numbers}")


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--feature", required=True)
    args = parser.parse_args()
    if args.feature not in ("FEATURE-0015", "FEATURE-0016", "FEATURE-0017"):
        raise SystemExit(f"FAIL: no feature contract is configured for {args.feature}")

    errors: list[str] = []
    registry = load_registry()
    if args.feature == "FEATURE-0015":
        check_f0015(registry, errors)
    elif args.feature == "FEATURE-0016":
        check_f0016(registry, errors)
    else:
        check_f0017(registry, errors)
    if errors:
        print(f"FAIL: {args.feature} feature-contract-check — {len(errors)} error(s)")
        for error in errors:
            print(f"  ✗ {error}")
        print("Kiro requirements/design/tasks generation is BLOCKED until the feature contract closes.")
        raise SystemExit(1)
    if args.feature == "FEATURE-0015":
        print("PASS: FEATURE-0015 feature contract closes route, audit, idempotency, authentication, validation-proof, and registration evidence")
    elif args.feature == "FEATURE-0016":
        print("PASS: FEATURE-0016 feature contract closes route, violation-code, conformance-catalog (VS0-CF-F16-01..128), ADH-2026-063 create phase-one/strict-classification precedence, ADH-2026-065 exact classification/phase-one-reference evidence, and the ADH-2026-066 private backing-access bridge")
    else:
        print("PASS: FEATURE-0017 feature contract closes request/result schemas, six decision groups, exact digest vector, adapter/fake boundaries, 24-case local conformance inventory, structural-only FEATURE-0013 mapping, and zero public route/store/controller")


if __name__ == "__main__":
    main()
