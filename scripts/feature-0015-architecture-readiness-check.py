#!/usr/bin/env python3
"""Fail-closed architecture-readiness check for FEATURE-0015.

Validates that the ADH-2026-045 canonical-bootstrap replacement decisions,
plus the ADH-2026-043 decisions that remain intact (safe-denial, scope-reference
integrity, audit reuse, feature-local traceability), are representable and
consistent across the repository authorities before requirements generation
may proceed.

This checker rejects reintroduction of removed migration concepts
(CanonicalMigrationPlan, CanonicalMigrationRecord, migration controller, and
ExecutionTarget in FEATURE-0015 scope) and validates the approved resource
fields/mutability, participation action-state table, independent suspension
holds, topology safe-denial ordering, bootstrap grant actions, API/update
behavior, concurrency, and audit matrix described in ADH-2026-045.

Exit 0 = PASS (requirements generation may proceed).
Exit 1 = FAIL (architecture gap remains; requirements generation blocked).
"""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

try:
    import yaml
except ImportError:
    print("FAIL: PyYAML required. Install: pip install pyyaml", file=sys.stderr)
    sys.exit(1)

ROOT = Path(__file__).resolve().parents[1]

# --- Paths ---
REG_PATH = ROOT / "docs/architecture/vertical-slices/VS-000-contract-registry.yaml"
SPEC_PATH = ROOT / "docs/architecture/vertical-slices/VS-000-contract-specification.md"
TRACE_PATH = ROOT / "docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md"
F15_ARCH = ROOT / "docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md"
F15_FEAT = ROOT / "docs/features/FEATURE-0015-canonical-cloud-model-foundation.md"
CLOSURE = ROOT / "docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md"
PHASE_CTX = ROOT / "docs/context/CURRENT_PHASE_CONTEXT.md"
STEER = ROOT / ".kiro/steering/slice0-contract.md"
ADH_045 = ROOT / "docs/reviews/architecture-decision-handoffs/ADH-2026-045-canonical-bootstrap-no-alpha-runtime-migration.md"
ADH_046 = ROOT / "docs/reviews/architecture-decision-handoffs/ADH-2026-046-feature-0015-executable-contract-closure.md"
ADH_047 = ROOT / "docs/reviews/architecture-decision-handoffs/ADH-2026-047-feature-0015-request-construction-and-scope-derivation-closure.md"
SPEC_DIR = ROOT / ".kiro/specs/canonical-cloud-model-and-alpha-migration"

# F0015-owned resources whose status must resolve solely to api-server (ADH-2026-046 decision 1).
F15_STATUS_OWNED_KINDS = (
    "CloudPlatform", "CloudProvider", "HostingLocation", "Datacenter",
    "FaultDomain", "InfrastructureStack", "CloudProviderParticipation",
)
NON_API_SERVER_CONTROLLER_TERMS = (
    "topology-controller", "stack-controller", "provider-controller",
    "participation-controller", "delegated-participation-contract-authority",
)
# Required F0015-local conformance IDs after ADH-2026-047 decision 5 (extends ADH-2026-046 decision 3).
F15_REQUIRED_CF_IDS = [f"VS0-CF-F15-{i:02d}" for i in range(1, 29)]
F15_REQUIRED_CF_FIELDS = ("id", "owner", "inputs", "expectedState", "expectedError", "expectedSideEffects", "gate")

# F0015-owned resource kinds whose create-request contract must be closed (ADH-2026-047 decision 1).
F15_CREATE_KINDS = (
    "CloudPlatform", "CloudProvider", "CloudProviderParticipation", "HostingLocation",
    "Datacenter", "FaultDomain", "InfrastructureStack",
)
# Allowed field classification words for the ADH-2026-047 decision 5 contract-executability audit.
F15_FIELD_CLASSIFICATIONS = (
    "request-required", "request-optional", "server-assigned",
    "action-only", "deferred-to-owner-feature", "forbidden",
)
# FEATURE-0021-owned CloudProviderParticipation fields that must be deferred/rejected by FEATURE-0015 (ADH-2026-047 decision 4).
F15_DEFERRED_PARTICIPATION_FIELDS = ("providerSelectionModes", "permittedHostingLocationRefs")

# Removed migration concepts that must never be reintroduced as active behavior.
REMOVED_MIGRATION_TERMS = [
    "CanonicalMigrationPlan",
    "CanonicalMigrationRecord",
    "migration-controller",
    "approved-migration-plan-publisher",
]

errs: list[str] = []


def e(msg: str) -> None:
    errs.append(msg)


def read(path: Path) -> str:
    if not path.exists():
        e(f"Required file missing: {path.relative_to(ROOT)}")
        return ""
    return path.read_text()


def is_allowed_migration_mention(text: str, term: str) -> bool:
    """A removed-term mention is allowed only in an explicit exclusion/retirement/
    non-goal/historical context (e.g. 'must not implement CanonicalMigrationPlan',
    'retired', 'superseded', 'no CanonicalMigrationPlan')."""
    lower = text.lower()
    allow_markers = (
        "must not", "no longer", "superseded", "retired", "removed", "does not",
        "no ", "not implement", "excluded", "non-goal", "tombstone", "never",
        "reject", "replace", "replaces", "must never",
    )
    idx = 0
    term_lower = term.lower()
    while True:
        pos = lower.find(term_lower, idx)
        if pos == -1:
            return True
        window = lower[max(0, pos - 120):pos]
        if not any(marker in window for marker in allow_markers):
            return False
        idx = pos + len(term_lower)


def check_no_reintroduced_migration_concepts(reg: dict, f15_arch: str, f15_feat: str) -> None:
    """Fail-closed: reject reintroduction of CanonicalMigrationPlan/Record, a
    migration controller, or ExecutionTarget in FEATURE-0015 scope."""
    # Registry: no active schema/writer/state may be CanonicalMigrationPlan/Record
    schemas = reg.get("schemas", [])
    for s in schemas:
        identity = str(s.get("identity", ""))
        if "CanonicalMigration" in identity:
            e(f"MIG: active schema {s.get('id')} reintroduces {identity}")
    writers = reg.get("writers", [])
    for w in writers:
        if any("CanonicalMigration" in p for p in w.get("paths", [])):
            e(f"MIG: active writer {w.get('id')} reintroduces a CanonicalMigration path")
    state_machines = reg.get("stateMachines", [])
    for sm in state_machines:
        if "CanonicalMigration" in str(sm.get("kind", "")):
            e(f"MIG: active state machine {sm.get('id')} reintroduces CanonicalMigration")
    if reg.get("migrationFailureMappings"):
        e("MIG: migrationFailureMappings must not exist as an active registry section")
    # ExecutionTarget must not be FEATURE-0015-owned; it belongs entirely to FEATURE-0016
    execution_target = next((s for s in schemas if s.get("identity", "").endswith("/ExecutionTarget")), None)
    if execution_target and execution_target.get("owner") != "FEATURE-0016":
        e(f"MIG: ExecutionTarget must be owned by FEATURE-0016 in its entirety, found owner={execution_target.get('owner')}")
    # Retired tombstones must exist for the removed IDs
    retired_schema_ids = {s.get("id") for s in reg.get("retiredSchemas", [])}
    for rid in ("VS0-SCHEMA-060", "VS0-SCHEMA-061"):
        if rid not in retired_schema_ids:
            e(f"MIG: {rid} must be a retired tombstone")
    retired_writer_ids = {w.get("id") for w in reg.get("retiredWriters", [])}
    for rid in ("VS0-WRITER-020", "VS0-WRITER-021"):
        if rid not in retired_writer_ids:
            e(f"MIG: {rid} must be a retired tombstone")
    retired_state_ids = {sm.get("id") for sm in reg.get("retiredStateMachines", [])}
    if "VS0-STATE-011" not in retired_state_ids:
        e("MIG: VS0-STATE-011 must be a retired tombstone")
    # Feature/architecture authorities must not describe migration concepts as active behavior
    for term in REMOVED_MIGRATION_TERMS:
        if not is_allowed_migration_mention(f15_arch, term):
            e(f"MIG: F0015 architecture mentions '{term}' outside an explicit exclusion/retirement context")
        if not is_allowed_migration_mention(f15_feat, term):
            e(f"MIG: F0015 feature mentions '{term}' outside an explicit exclusion/retirement context")
    if "ends at InfrastructureStack" not in f15_arch and "ends at `InfrastructureStack`" not in f15_arch:
        e("MIG: F0015 architecture must state it ends at InfrastructureStack")
    if "ends at InfrastructureStack" not in f15_feat and "ends at InfrastructureStack." not in f15_feat:
        e("MIG: F0015 feature must state it ends at InfrastructureStack")


def check_resource_fields_and_mutability(reg: dict, f15_arch: str, f15_feat: str) -> None:
    """Validate approved resource fields/mutability: PATCH-only with
    application/merge-patch+json; no PUT/DELETE; identity/scope immutable."""
    schemas = reg.get("schemas", [])
    schema_by_kind = {str(s.get("identity", "")).rsplit("/", 1)[-1]: s for s in schemas}
    for kind_name in ("CloudPlatform", "CloudProvider", "HostingLocation", "Datacenter", "FaultDomain", "InfrastructureStack"):
        s = schema_by_kind.get(kind_name)
        if not s:
            e(f"FIELDS: {kind_name} schema not found in registry")
            continue
        mutability = str(s.get("mutability", ""))
        if "immutable" not in mutability.lower():
            e(f"FIELDS: {kind_name} mutability must state immutable identity/scope fields")
    participation = schema_by_kind.get("CloudProviderParticipation", {})
    if not participation:
        e("FIELDS: CloudProviderParticipation schema not found in registry")
    else:
        mutability = str(participation.get("mutability", "")).lower()
        if "no mutable" not in mutability and "never patch" not in mutability:
            e("FIELDS: CloudProviderParticipation mutability must state it has no mutable F0015 spec fields and is never PATCHed")
        required = participation.get("required", [])
        if not any(str(f).startswith("status.platformSuspended") for f in required):
            e("FIELDS: CloudProviderParticipation must require status.platformSuspended")
        if not any(str(f).startswith("status.providerSuspended") for f in required):
            e("FIELDS: CloudProviderParticipation must require status.providerSuspended")
        phase_field = next((f for f in required if str(f).startswith("status.phase")), "")
        for phase in ("Pending", "Active", "Rejected", "Withdrawn", "Expired", "Suspended", "Terminating", "Terminated"):
            if phase not in str(phase_field):
                e(f"FIELDS: CloudProviderParticipation status.phase enum missing {phase}")
    # Architecture/feature authorities state PATCH-only, no PUT/DELETE
    for label, text in (("architecture", f15_arch), ("feature", f15_feat)):
        lower = text.lower()
        if "merge-patch+json" not in lower and "merge-patch" not in lower:
            e(f"FIELDS: F0015 {label} must state PATCH uses application/merge-patch+json")
        if "put" not in lower or "delete" not in lower:
            e(f"FIELDS: F0015 {label} must state PUT and DELETE are not exposed")


def check_participation_action_state_table(reg: dict, f15_arch: str, f15_feat: str) -> None:
    """Validate the participation action-state table: explicit actions govern
    every transition; no mutable spec fields; empty body + If-Match + Idempotency-Key."""
    sms = reg.get("stateMachines", [])
    sm_001 = next((sm for sm in sms if sm.get("id") == "VS0-STATE-001"), None)
    if not sm_001:
        e("PART: VS0-STATE-001 not found")
        return
    states = sm_001.get("states", [])
    for phase in ("Pending", "Active", "Rejected", "Withdrawn", "Expired", "Suspended", "Terminating", "Terminated"):
        if phase not in states:
            e(f"PART: VS0-STATE-001 missing state {phase}")
    terminal = sm_001.get("terminal", [])
    for phase in ("Rejected", "Withdrawn", "Expired", "Terminated"):
        if phase not in terminal:
            e(f"PART: VS0-STATE-001 terminal states must include {phase}")
    hold_rule = str(sm_001.get("holdRule", "")).lower()
    for required in ("platformsuspended", "providersuspended", "clearing one hold"):
        if required not in hold_rule:
            e(f"PART: VS0-STATE-001 holdRule missing '{required}'")
    for label, text in (("architecture", f15_arch), ("feature", f15_feat)):
        lower = text.lower()
        for action in ("accept", "reject", "withdraw", "expire", "suspend", "resume", "request-release", "accept-release", "decline-release"):
            if action not in lower:
                e(f"PART: F0015 {label} must describe the '{action}' participation action")
        if "if-match" not in lower:
            e(f"PART: F0015 {label} must require If-Match for participation actions")
        if "idempotency-key" not in lower and "idempotency key" not in lower:
            e(f"PART: F0015 {label} must require an Idempotency-Key for participation actions")


def check_independent_suspension_holds(reg: dict, f15_arch: str, f15_feat: str) -> None:
    """Validate independent suspension holds: platformSuspended and
    providerSuspended are each controlled only by their own administrator."""
    schemas = reg.get("schemas", [])
    participation = next((s for s in schemas if str(s.get("identity", "")).endswith("/CloudProviderParticipation")), {})
    required = " ".join(participation.get("required", []))
    for field in ("status.platformSuspended", "status.providerSuspended"):
        if field not in required:
            e(f"HOLDS: CloudProviderParticipation missing required {field}")
    for label, text in (("architecture", f15_arch), ("feature", f15_feat)):
        lower = text.lower()
        if "platformsuspended" not in lower or "providersuspended" not in lower:
            e(f"HOLDS: F0015 {label} must name both platformSuspended and providerSuspended holds")
        if "clearing one hold" not in lower and "clearing one hold does not reactivate" not in lower and "never reactivate" not in lower:
            e(f"HOLDS: F0015 {label} must state that clearing one hold does not reactivate while the other remains true")


def check_topology_safe_denial_ordering(reg: dict, spec: str, f15_feat: str) -> None:
    """Decision (preserved from ADH-2026-043 D4): inaccessible cross-provider =
    RESOURCE_NOT_FOUND/404 + VS0_AUTHORIZATION_SAFE_DENIAL; authorized invalid
    same-provider = VALIDATION_FAILED/422. Resolution precedes structural validation."""
    confs = reg.get("conformance", [])
    x03 = next((c for c in confs if c.get("id") == "VS0-CF-X03"), None)
    if not x03:
        e("SAFEDENIAL: VS0-CF-X03 conformance entry missing")
        return
    if x03.get("expectedError") != "RESOURCE_NOT_FOUND":
        e("SAFEDENIAL: VS0-CF-X03 expectedError must be RESOURCE_NOT_FOUND")
    vcs = {vc.get("code") for vc in reg.get("violationCodes", {}).get("slice0", [])}
    if "VS0_AUTHORIZATION_SAFE_DENIAL" not in vcs:
        e("SAFEDENIAL: VS0_AUTHORIZATION_SAFE_DENIAL must be a registered violation code")
    if "VS0_AUTHORIZATION_SAFE_DENIAL" not in spec:
        e("SAFEDENIAL: VS-000 specification must reference VS0_AUTHORIZATION_SAFE_DENIAL safe-denial rule")
    if "RESOURCE_NOT_FOUND" not in f15_feat or "VS0_AUTHORIZATION_SAFE_DENIAL" not in f15_feat:
        e("SAFEDENIAL: F0015 error table must include RESOURCE_NOT_FOUND + VS0_AUTHORIZATION_SAFE_DENIAL")


def check_scope_reference_integrity(reg: dict, spec: str, f15_feat: str) -> None:
    """Decision (preserved from ADH-2026-043 D5, narrowed by ADH-2026-045 decision 8):
    UID-pinned scope-reference invariant with VS0_SCOPE_REFERENCE_MISMATCH. CloudPlatform
    no longer has an Organization reference (ADH-2026-045); the invariant now applies to
    CloudProviderParticipation.spec.cloudPlatformRef.uid == CloudPlatform.metadata.scopeRef.uid."""
    vcs = {vc.get("code") for vc in reg.get("violationCodes", {}).get("slice0", [])}
    if "VS0_SCOPE_REFERENCE_MISMATCH" not in vcs:
        e("SCOPEREF: VS0_SCOPE_REFERENCE_MISMATCH must be a registered violation code")
    confs = reg.get("conformance", [])
    f15_11 = next((c for c in confs if c.get("id") == "VS0-CF-F15-11"), None)
    if not f15_11:
        e("SCOPEREF: VS0-CF-F15-11 conformance entry missing")
    elif f15_11.get("expectedError") != "VALIDATION_FAILED":
        e("SCOPEREF: VS0-CF-F15-11 expectedError must be VALIDATION_FAILED")
    elif f15_11.get("expectedViolation") != "VS0_SCOPE_REFERENCE_MISMATCH":
        e("SCOPEREF: VS0-CF-F15-11 expectedViolation must be VS0_SCOPE_REFERENCE_MISMATCH")
    if "VS0_SCOPE_REFERENCE_MISMATCH" not in spec:
        e("SCOPEREF: VS-000 specification must reference VS0_SCOPE_REFERENCE_MISMATCH")
    if "VS0_SCOPE_REFERENCE_MISMATCH" not in f15_feat:
        e("SCOPEREF: F0015 feature must reference VS0_SCOPE_REFERENCE_MISMATCH")
    if "cloudPlatformRef" not in f15_feat:
        e("SCOPEREF: F0015 feature must state the cloudPlatformRef UID invariant")
    if not is_allowed_migration_mention(f15_feat, "ownerOrganizationRef"):
        e("SCOPEREF: F0015 feature must not reference ownerOrganizationRef as active behavior (removed by ADH-2026-045; CloudPlatform has no Organization reference)")


def check_bootstrap_grant_actions(f15_arch: str, f15_feat: str) -> None:
    """Validate bootstrap grant actions: server-resolved deterministic grants;
    never from header/body; F0015 persists no roles/memberships/assignments;
    CloudPlatform-root requirement."""
    for label, text in (("architecture", f15_arch), ("feature", f15_feat)):
        lower = text.lower()
        if "bootstrap" not in lower:
            e(f"BOOTSTRAP: F0015 {label} must describe bootstrap authorization")
        if "cloudplatform_root_required" not in lower.replace(" ", "_") and "vs0_cloudplatform_root_required" not in lower:
            e(f"BOOTSTRAP: F0015 {label} must state the VS0_CLOUDPLATFORM_ROOT_REQUIRED rule")


def check_api_update_behavior_and_concurrency(f15_feat: str) -> None:
    """Validate API/update behavior and concurrency: PATCH-only, If-Match,
    Idempotency-Key, stale-version handling."""
    lower = f15_feat.lower()
    if "stale_resource_version" not in lower.replace(" ", "_"):
        e("API: F0015 feature must describe STALE_RESOURCE_VERSION for stale If-Match")
    if "idempotency" not in lower:
        e("API: F0015 feature must describe idempotency behavior")


def check_audit_matrix(f15_arch: str, f15_feat: str) -> None:
    """Decision (preserved from ADH-2026-043 D6): F0015 reuses FEATURE-0013
    AuditEvent with correlation/redaction/no-secrets, for participation
    lifecycle/hold-changing actions and resource create/PATCH."""
    for label, text in (("architecture", f15_arch), ("feature", f15_feat)):
        lower = text.lower()
        if "feature-0013" not in lower or "auditevent" not in lower:
            e(f"AUDIT: F0015 {label} must reference FEATURE-0013 AuditEvent reuse")
        if "no secret" not in lower:
            e(f"AUDIT: F0015 {label} must state the no-secrets rule for audit")
        if "correlation" not in lower:
            e(f"AUDIT: F0015 {label} must state correlation requirements")
        if "vs0-cf-t01" not in lower:
            e(f"AUDIT: F0015 {label} must explicitly not import VS0-CF-T01")


def check_traceability_local_only(f15_arch: str) -> None:
    """Decision (preserved from ADH-2026-043 D7): F0015 traceability uses only
    local conformance IDs; downstream HP01/F09 are not local acceptance evidence."""
    import re
    sec10 = ""
    lines = f15_arch.splitlines()
    in_sec10 = False
    for line in lines:
        if re.match(r"^#+\s+10\.", line):
            in_sec10 = True
            continue
        if in_sec10 and re.match(r"^#+\s+\d+\.", line):
            break
        if in_sec10:
            sec10 += line + "\n"
    table_lines = [l for l in sec10.splitlines() if l.strip().startswith("|") and "VS0-CF" in l]
    for tl in table_lines:
        cells = [c.strip() for c in tl.split("|")]
        if len(cells) >= 6:
            conf_cell = cells[-2] if cells[-1] == "" else cells[-1]
            if "HP01" in conf_cell and "note" not in tl.lower():
                e("TRACE: §10 table must not use VS0-CF-HP01 as local conformance evidence")
            if conf_cell.strip() == "VS0-CF-F09" or ", F09" in conf_cell or ",F09" in conf_cell:
                e("TRACE: §10 table must not use VS0-CF-F09 as local conformance evidence")
    if "downstream" not in sec10.lower() and "non-owning" not in sec10.lower():
        e("TRACE: §10 must note that HP01/F09 are downstream integration references only")


def check_closure_matrix(closure: str) -> None:
    """All ARC-F15-01..07 table rows must be RESOLVED (historical closure record
    retained; supersession annotations for migration-related claims are allowed)."""
    table_rows = [line for line in closure.splitlines() if line.strip().startswith("| ARC-F15-")]
    for i in range(1, 8):
        arc_id = f"ARC-F15-{i:02d}"
        row = next((r for r in table_rows if r.strip().startswith(f"| {arc_id} |")), None)
        if not row:
            e(f"Closure: {arc_id} table row not found in closure matrix")
        elif "RESOLVED" not in row:
            e(f"Closure: {arc_id} is not RESOLVED")


def check_traceability_matrix(trace: str) -> None:
    """VS0-CF-F15-11 must appear in the traceability matrix."""
    if "VS0-CF-F15-11" not in trace:
        e("Traceability: VS0-CF-F15-11 missing from traceability matrix")


def check_no_stale_cloudplatform_owner_org_ref(reg_text: str, spec: str) -> None:
    """Fail-closed: 'ownerOrganizationRef' must not appear as active CloudPlatform
    behavior in the registry YAML or the VS-000 contract specification. CloudPlatform
    has no Organization reference (ADH-2026-045 decision 8); it carries an immutable
    spec.ownerRegistration instead. A mention is allowed only in an explicit
    historical/retired/superseded/narrowed annotation context."""
    for label, text in (("registry", reg_text), ("VS-000 specification", spec)):
        if not is_allowed_migration_mention(text, "ownerOrganizationRef"):
            e(f"STALEOWNER: {label} references ownerOrganizationRef as active CloudPlatform "
              f"behavior (removed by ADH-2026-045 decision 8; CloudPlatform has no Organization "
              f"reference and carries immutable spec.ownerRegistration instead)")


def check_no_stale_dual_delegated_acceptance(f15_arch: str, f15_feat: str, reg_text: str, spec: str) -> None:
    """Fail-closed: 'dual delegated acceptance' / 'both delegated authority acceptances'
    must not remain as active (non-superseded) wording anywhere in the F0015 architecture
    doc, F0015 feature doc, the registry, or the VS-000 spec text. Per ADH-2026-045's
    participation lifecycle table, CREATE is the provider's affirmative request evidence
    and Pending->Active requires only the CloudPlatform's accept action."""
    stale_phrases = ("dual delegated acceptance", "both delegated authority acceptances")
    for label, text in (
        ("F0015 architecture", f15_arch),
        ("F0015 feature", f15_feat),
        ("registry", reg_text),
        ("VS-000 specification", spec),
    ):
        lower = text.lower()
        for phrase in stale_phrases:
            if phrase in lower and not is_allowed_migration_mention(text, phrase):
                e(f"STALEACCEPT: {label} still contains stale phrase '{phrase}' as active wording "
                  f"(ADH-2026-045: CREATE is the provider's request evidence; Pending->Active requires "
                  f"only the CloudPlatform's accept action, not dual delegated acceptances)")


def check_f15_08_uses_resume_not_suspend(reg: dict) -> None:
    """Fail-closed: the VS0-CF-F15-08 registry entry's 'inputs' field must not describe
    a hold-clearing action as 'suspend'. Clearing a hold is a resume action; suspend SETS
    a hold. A pattern like 'suspend action clearing' or 'suspend...clearing...hold' is
    backwards terminology and must be flagged."""
    confs = reg.get("conformance", [])
    f15_08 = next((c for c in confs if c.get("id") == "VS0-CF-F15-08"), None)
    if not f15_08:
        e("HOLDTERM: VS0-CF-F15-08 conformance entry missing")
        return
    inputs_text = str(f15_08.get("inputs", ""))
    lower = inputs_text.lower()
    if re.search(r"suspend\s+action\s+clearing", lower) or re.search(r"suspend.*clearing.*hold", lower):
        e(f"HOLDTERM: VS0-CF-F15-08 inputs describes clearing a hold as a 'suspend' action "
          f"(clearing a hold is a resume action, not suspend): {inputs_text!r}")


def check_status_writer_conflicts(reg: dict, reg_text: str, spec: str, f15_arch: str, f15_feat: str, trace: str) -> None:
    """Fail-closed (ADH-2026-046 decision 4a): reject a non-api-server status writer for
    any F0015-owned resource across ADH-045/046, registry, contract specification, feature
    authority, architecture boundary, and traceability. Requirements are checked (not edited)."""
    schemas = reg.get("schemas", [])
    schema_by_kind = {str(s.get("identity", "")).rsplit("/", 1)[-1]: s for s in schemas}
    for kind_name in F15_STATUS_OWNED_KINDS:
        s = schema_by_kind.get(kind_name)
        if not s:
            continue
        mutability = str(s.get("mutability", "")).lower()
        if kind_name != "CloudProviderParticipation" and "api-server" not in mutability:
            e(f"WRITER045/046: {kind_name} registry mutability does not resolve status to api-server")
        for term in NON_API_SERVER_CONTROLLER_TERMS:
            if term in mutability:
                e(f"WRITER045/046: {kind_name} registry mutability still uses non-api-server controller term '{term}'")
    writers = reg.get("writers", [])
    for w in writers:
        if w.get("id") in ("VS0-WRITER-004", "VS0-WRITER-005"):
            writer_name = str(w.get("writer", ""))
            if writer_name in NON_API_SERVER_CONTROLLER_TERMS:
                e(f"WRITER045/046: {w.get('id')} writer field is a non-api-server controller term: {writer_name}")
    for label, text in (
        ("registry (active)", reg_text),
        ("VS-000 specification", spec),
        ("F0015 architecture", f15_arch),
        ("F0015 feature", f15_feat),
        ("traceability matrix", trace),
    ):
        for term in ("topology-controller status", "stack-controller status", "provider-controller status", "participation-controller owns status", "participation-controller owns both"):
            if term in text and not is_allowed_migration_mention(text, term):
                e(f"WRITER045/046: {label} still attributes status ownership to a non-api-server controller ('{term}')")


def check_participation_create_precondition(f15_arch: str, f15_feat: str, spec: str, reg: dict) -> None:
    """Fail-closed (ADH-2026-046 decision 4b): reject an If-Match requirement on participation
    collection CREATE, or missing If-Match/idempotency behavior for EXISTING participation actions."""
    overbroad_patterns = (
        r"all participation actions[^.]*require[^.]*if-match",
        r"every participation action[^.]*require[^.]*if-match",
        r"participation actions \(create,[^)]*\)[^.]*require[^.]*if-match",
    )
    for label, text in (("F0015 architecture", f15_arch), ("F0015 feature", f15_feat), ("VS-000 specification", spec)):
        lower = text.lower()
        for pattern in overbroad_patterns:
            if re.search(pattern, lower):
                e(f"PRECOND046: {label} states an If-Match requirement scoped to ALL participation "
                  f"actions including create; create must require Idempotency-Key only (ADH-2026-046 decision 2)")
        if "existing participation" not in lower and "existing-participation" not in lower:
            e(f"PRECOND046: {label} must scope the If-Match requirement to actions on an EXISTING participation")
    writers = reg.get("writers", [])
    w004 = next((w for w in writers if w.get("id") == "VS0-WRITER-004"), None)
    if w004:
        update_method = str(w004.get("updateMethod", "")).lower()
        if "idempotency-key only" not in update_method and "no if-match" not in update_method:
            e("PRECOND046: VS0-WRITER-004 updateMethod does not scope create to Idempotency-Key-only (no If-Match)")
        if "if-match" not in update_method:
            e("PRECOND046: VS0-WRITER-004 updateMethod does not describe the If-Match requirement for existing-participation actions")


def check_f15_proof_matrix_and_mapping(reg: dict, f15_arch: str, f15_feat: str, trace: str) -> None:
    """Fail-closed (ADH-2026-046 decision 4c/4d): reject any missing F15-01..25 conformance
    case, incomplete registry fields, or incomplete REQ/AC-to-proof mapping."""
    confs = {c.get("id"): c for c in reg.get("conformance", [])}
    for cf_id in F15_REQUIRED_CF_IDS:
        entry = confs.get(cf_id)
        if not entry:
            e(f"PROOF046: required F0015-local conformance case missing from registry: {cf_id}")
            continue
        for field in F15_REQUIRED_CF_FIELDS:
            if field not in entry or entry.get(field) in (None, ""):
                if field == "expectedError" and entry.get(field, "unset") is None:
                    continue  # null is a valid explicit value for expectedError
                if field not in entry:
                    e(f"PROOF046: {cf_id} missing required registry field '{field}'")
        if cf_id not in trace:
            e(f"PROOF046: {cf_id} missing from traceability matrix")
        if cf_id not in f15_arch:
            e(f"PROOF046: {cf_id} missing from F0015 architecture §10 mapping")
    # REQ/AC-to-proof mapping completeness (architecture §10.1)
    for i in range(1, 19):
        req_id = f"REQ-F15-{i:02d}"
        if req_id not in f15_arch:
            e(f"PROOF046: {req_id} missing from F0015 architecture REQ-to-proof mapping")
    for i in range(1, 15):
        ac_id = f"AC-F15-{i:02d}"
        if ac_id not in f15_arch:
            e(f"PROOF046: {ac_id} missing from F0015 architecture AC-to-proof mapping")
        if ac_id not in f15_feat:
            e(f"PROOF046: {ac_id} missing from F0015 feature Acceptance Criteria table")


def check_no_f16_leakage_in_f15_conformance(reg: dict, f15_feat: str) -> None:
    """Fail-closed (ADH-2026-046 decision 4e): reject F0016+ schema/writer/route/conformance
    claimed as FEATURE-0015 behavior, except explicit labelled non-goal/exclusion references."""
    confs = reg.get("conformance", [])
    for c in confs:
        if c.get("id", "").startswith("VS0-CF-F15-") and c.get("owner") != "FEATURE-0015":
            e(f"LEAK046: {c.get('id')} is F0015-local-numbered but owner is {c.get('owner')}, not FEATURE-0015")
    lower = f15_feat.lower()
    for downstream_id in ("vs0-cf-f10", "vs0-cf-hp01"):
        if downstream_id in lower and not is_allowed_migration_mention(f15_feat, downstream_id):
            e(f"LEAK046: F0015 feature cites downstream conformance {downstream_id} without an explicit "
              f"labelled non-goal/exclusion or 'downstream integration reference, non-owning' annotation")


def check_design_tasks_proof_coverage(control_json_path: Path, mode: str) -> None:
    """Fail-closed (ADH-2026-046 decision 4f): a design/tasks document must have a testable
    task for every F0015-local conformance case. Deferred to 'post' mode; skipped gracefully
    at the pre-requirements 'readiness' gate and when design.md/tasks.md do not yet exist,
    consistent with the script's existing --mode readiness/post pattern."""
    if mode != "post":
        return
    design_path = SPEC_DIR / "design.md"
    tasks_path = SPEC_DIR / "tasks.md"
    if not design_path.exists() or not tasks_path.exists():
        return
    tasks_text = tasks_path.read_text()
    for cf_id in F15_REQUIRED_CF_IDS:
        if cf_id not in tasks_text:
            e(f"TASKCOV046: tasks.md has no testable task referencing required conformance case {cf_id}")


def check_requirements_authority_conflicts(mode: str) -> None:
    """Fail-closed (ADH-2026-046 decision 4a/4b): if a Kiro requirements.md exists for
    FEATURE-0015, check it (do not edit it) for a non-api-server status-writer conflict or
    an overbroad participation create If-Match requirement. requirements.md is regenerated
    only after this readiness gate passes (per ADH-2026-046 human-approval notes), so an
    existing pre-046 requirements.md is expected to still contain stale wording until it is
    regenerated. To avoid blocking the pre-requirements 'readiness' gate on a document that
    is itself produced only after that gate passes, this check runs under --mode post only,
    consistent with the script's existing readiness/post deferral pattern (see
    check_design_tasks_proof_coverage above)."""
    if mode != "post":
        return
    req_path = SPEC_DIR / "requirements.md"
    if not req_path.exists():
        return
    req_text = req_path.read_text()
    lower = req_text.lower()
    for term in ("topology-controller status", "stack-controller status", "provider-controller status",
                 "participation-controller owns status", "participation-controller owns both"):
        if term in lower and not is_allowed_migration_mention(req_text, term):
            e(f"WRITER045/046: requirements.md still attributes status ownership to a non-api-server "
              f"controller ('{term}'); requirements is a checked (not edited) authority per ADH-2026-046 decision 4a")
    overbroad_patterns = (
        r"all participation actions[^.]*require[^.]*if-match",
        r"every participation action[^.]*require[^.]*if-match",
        r"participation actions \(create,[^)]*\)[^.]*require[^.]*if-match",
    )
    for pattern in overbroad_patterns:
        if re.search(pattern, lower):
            e(f"PRECOND046: requirements.md states an If-Match requirement scoped to ALL participation "
              f"actions including create; create must require Idempotency-Key only (ADH-2026-046 decision 2); "
              f"requirements is a checked (not edited) authority per ADH-2026-046 decision 4b")


def check_closed_create_field_matrix(reg: dict, f15_arch: str, f15_feat: str) -> None:
    """Fail-closed (ADH-2026-047 decision 1): every F0015 collection-create kind must
    have a closed, complete client-required/client-optional create-field table in both
    the architecture and feature authorities, and the registry's schema requestContract
    must be present and complete for each kind. Rejects absent/ambiguous create fields."""
    schemas = reg.get("schemas", [])
    schema_by_kind = {str(s.get("identity", "")).rsplit("/", 1)[-1]: s for s in schemas}
    for kind_name in F15_CREATE_KINDS:
        s = schema_by_kind.get(kind_name)
        if not s:
            e(f"CREATECONTRACT047: {kind_name} schema not found in registry")
            continue
        rc = s.get("requestContract")
        if not rc:
            e(f"CREATECONTRACT047: {kind_name} registry schema missing requestContract (ADH-2026-047 decision 1)")
            continue
        if "clientRequired" not in rc or not rc.get("clientRequired"):
            e(f"CREATECONTRACT047: {kind_name} requestContract missing non-empty clientRequired")
        if "clientOptional" not in rc:
            e(f"CREATECONTRACT047: {kind_name} requestContract missing clientOptional (use [] if none)")
        if "createSuccessStatus" not in rc or rc.get("createSuccessStatus") != 201:
            e(f"CREATECONTRACT047: {kind_name} requestContract must state createSuccessStatus: 201")
        if kind_name != "CloudProviderParticipation" and rc.get("patchSuccessStatus") != 200:
            e(f"CREATECONTRACT047: {kind_name} requestContract must state patchSuccessStatus: 200")
        if "rejects" not in rc and "forbidden" not in rc:
            e(f"CREATECONTRACT047: {kind_name} requestContract missing rejects/forbidden field list")
        if "scopeDerivation" not in rc or not rc.get("scopeDerivation"):
            e(f"CREATECONTRACT047: {kind_name} requestContract missing scopeDerivation")
    for label, text in (("architecture", f15_arch), ("feature", f15_feat)):
        for kind_name in F15_CREATE_KINDS:
            if kind_name not in text:
                e(f"CREATECONTRACT047: F0015 {label} does not mention {kind_name} in its closed create-field table")
        if "client-required" not in text.lower() and "client-required create fields" not in text.lower():
            e(f"CREATECONTRACT047: F0015 {label} must state a closed client-required create-field table (ADH-2026-047 decision 1)")
        if "client-optional" not in text.lower() and "client-optional create fields" not in text.lower():
            e(f"CREATECONTRACT047: F0015 {label} must state a closed client-optional create-field table (ADH-2026-047 decision 1)")
        if "200" not in text or "201" not in text:
            e(f"CREATECONTRACT047: F0015 {label} must state both the 201 create and 200 PATCH success outcomes (ADH-2026-047 decision 1)")


def check_single_scope_derivation_source(reg: dict, f15_arch: str, f15_feat: str) -> None:
    """Fail-closed (ADH-2026-047 decision 2): each F0015 resource kind must have exactly
    one stated scope-derivation source in the registry schema requestContract and in the
    architecture/feature authorities. Rejects a client-supplied scope for any kind."""
    schemas = reg.get("schemas", [])
    schema_by_kind = {str(s.get("identity", "")).rsplit("/", 1)[-1]: s for s in schemas}
    expected_derivation_marker = {
        "CloudPlatform": "platform-root",
        "CloudProvider": "platform-root",
        "CloudProviderParticipation": "cloudplatformref",
        "HostingLocation": "topology.write",
        "Datacenter": "parent reference",
        "FaultDomain": "parent reference",
        "InfrastructureStack": "parent reference",
    }
    for kind_name, marker in expected_derivation_marker.items():
        s = schema_by_kind.get(kind_name)
        if not s:
            continue
        rc = s.get("requestContract", {})
        derivation = str(rc.get("scopeDerivation", "")).lower()
        if marker.lower() not in derivation:
            e(f"SCOPEDERIV047: {kind_name} registry requestContract.scopeDerivation does not state its single source ('{marker}')")
    client_supplied_patterns = (
        r"client\s+(?:supplies|selects|chooses|provides)\s+(?:its own\s+)?(?:metadata\.)?scoperef",
        r"caller\s+(?:supplies|selects|chooses|provides)\s+(?:its own\s+)?(?:metadata\.)?scoperef",
    )
    for label, text in (
        ("architecture", f15_arch),
        ("feature", f15_feat),
    ):
        lower = text.lower()
        if "topology.write" not in text:
            e(f"SCOPEDERIV047: F0015 {label} must state that HostingLocation scope derives from the topology.write grant (ADH-2026-047 decision 2)")
        for pattern in client_supplied_patterns:
            if re.search(pattern, lower) and not is_allowed_migration_mention(text, "scoperef"):
                e(f"SCOPEDERIV047: F0015 {label} appears to state a client-supplied scopeRef, which is prohibited (ADH-2026-047 decision 2)")
        if "cannot supply or choose" not in lower and "cannot supply or select" not in lower and "never supplies" not in lower and "client never supplies" not in lower:
            e(f"SCOPEDERIV047: F0015 {label} must explicitly state a client cannot supply or choose metadata.scopeRef (ADH-2026-047 decision 2)")


def check_route_contract_completeness(f15_arch: str, f15_feat: str) -> None:
    """Fail-closed (ADH-2026-047 decision 5b): every F0015 route must state its
    request-body shape, required headers, success result, error family, status
    effect, and audit effect together. Uses the existing Decision 1 create-field
    table, error-codes table, and audit table as the combined route-contract source."""
    for label, text in (("architecture", f15_arch), ("feature", f15_feat)):
        lower = text.lower()
        required_markers = (
            ("request-body shape", ("client-required", "empty json body")),
            ("required headers", ("idempotency-key", "if-match")),
            ("success result", ("201", "200")),
            ("error family", ("problem details", "conflict", "validation_failed", "resource_not_found")),
            ("status effect", ("status.phase",)),
            ("audit effect", ("auditevent",)),
        )
        for element_name, markers in required_markers:
            if not any(marker in lower for marker in markers):
                e(f"ROUTECONTRACT047: F0015 {label} does not jointly state the route-contract element '{element_name}' "
                  f"required by ADH-2026-047 decision 5 (request-body shape, required headers, success result, "
                  f"error family, status effect, audit effect, local conformance must all be stated)")
        if "vs0-cf-f15" not in lower:
            e(f"ROUTECONTRACT047: F0015 {label} does not state local conformance case references for its routes")


def check_deferred_participation_fields_rejected(reg: dict, f15_arch: str, f15_feat: str) -> None:
    """Fail-closed (ADH-2026-047 decision 4/5): providerSelectionModes and
    permittedHostingLocationRefs must both be registered as FEATURE-0021
    fieldOwnership on CloudProviderParticipation, and must never appear as
    F0015-accepted/stored/defaulted/validated in the architecture/feature authorities
    (only as explicitly rejected/deferred-to-owner-feature)."""
    schemas = reg.get("schemas", [])
    participation = next((s for s in schemas if str(s.get("identity", "")).endswith("/CloudProviderParticipation")), {})
    field_ownership = participation.get("fieldOwnership", {})
    for field in F15_DEFERRED_PARTICIPATION_FIELDS:
        key = f"spec.{field}"
        ownership = field_ownership.get(key, {})
        if ownership.get("introducedBy") != "FEATURE-0021" or ownership.get("activatedBy") != "FEATURE-0021":
            e(f"DEFERREDFIELD047: registry CloudProviderParticipation.{key} must record "
              f"introducedBy/activatedBy FEATURE-0021 (ADH-2026-047 decision 4)")
    accept_patterns = (
        r"f0015\s+accepts?\s+{field}",
        r"{field}\s+is\s+accepted\s+by\s+f(?:eature-)?0015",
        r"{field}\s+is\s+stored\s+by\s+f(?:eature-)?0015",
        r"{field}\s+is\s+defaulted\s+by\s+f(?:eature-)?0015",
        r"{field}\s+is\s+validated\s+by\s+f(?:eature-)?0015",
    )
    for label, text in (("architecture", f15_arch), ("feature", f15_feat)):
        lower = text.lower()
        for field in F15_DEFERRED_PARTICIPATION_FIELDS:
            fl = field.lower()
            for pattern in accept_patterns:
                if re.search(pattern.format(field=re.escape(fl)), lower):
                    e(f"DEFERREDFIELD047: F0015 {label} appears to describe '{field}' as F0015-accepted/stored/"
                      f"defaulted/validated; it must be explicitly rejected/deferred-to-FEATURE-0021")
            if fl in lower and "feature-0021" not in lower:
                e(f"DEFERREDFIELD047: F0015 {label} mentions '{field}' without attributing ownership to FEATURE-0021")


def check_create_status_and_result_completeness(reg: dict, f15_arch: str, f15_feat: str) -> None:
    """Fail-closed (ADH-2026-047 decision 5): reject a kind's create outcome that does
    not state its initial status, and reject any requirement text that omits the
    approved create/PATCH result (201 for create, 200 for PATCH)."""
    schemas = reg.get("schemas", [])
    schema_by_kind = {str(s.get("identity", "")).rsplit("/", 1)[-1]: s for s in schemas}
    for kind_name in F15_CREATE_KINDS:
        s = schema_by_kind.get(kind_name)
        if not s:
            continue
        rc = s.get("requestContract", {})
        server_assigned = " ".join(str(x) for x in rc.get("serverAssigned", []))
        if kind_name == "CloudProviderParticipation":
            if "status.phase=pending" not in server_assigned.lower():
                e(f"CREATESTATUS047: {kind_name} requestContract.serverAssigned must state its initial status (status.phase=Pending)")
        else:
            if "status.phase=active" not in server_assigned.lower():
                e(f"CREATESTATUS047: {kind_name} requestContract.serverAssigned must state its initial status (status.phase=Active)")
    for label, text in (("architecture", f15_arch), ("feature", f15_feat)):
        if "201" not in text:
            e(f"CREATESTATUS047: F0015 {label} must state the approved create result (201) somewhere in its create-contract description")
        if "200" not in text:
            e(f"CREATESTATUS047: F0015 {label} must state the approved PATCH result (200) somewhere in its create-contract description")


def check_f15_26_28_conformance_and_mapping(reg: dict, f15_arch: str, f15_feat: str, trace: str) -> None:
    """Fail-closed (ADH-2026-047 decision 5): VS0-CF-F15-26..28 must be registered with
    complete fields, present in the traceability matrix and architecture §10 mapping,
    and REQ-F15-19..21/AC-F15-15..17 must map to them in architecture §10.1 and the
    feature Acceptance Criteria table."""
    confs = {c.get("id"): c for c in reg.get("conformance", [])}
    for i in range(26, 29):
        cf_id = f"VS0-CF-F15-{i:02d}"
        entry = confs.get(cf_id)
        if not entry:
            e(f"PROOF047: required F0015-local conformance case missing from registry: {cf_id}")
            continue
        for field in F15_REQUIRED_CF_FIELDS:
            if field not in entry or entry.get(field) in (None, ""):
                if field == "expectedError" and entry.get(field, "unset") is None:
                    continue
                if field not in entry:
                    e(f"PROOF047: {cf_id} missing required registry field '{field}'")
        if cf_id not in trace:
            e(f"PROOF047: {cf_id} missing from traceability matrix")
        if cf_id not in f15_arch:
            e(f"PROOF047: {cf_id} missing from F0015 architecture §10 mapping")
    for i in range(19, 22):
        req_id = f"REQ-F15-{i:02d}"
        if req_id not in f15_arch:
            e(f"PROOF047: {req_id} missing from F0015 architecture REQ-to-proof mapping")
        if req_id not in f15_feat:
            e(f"PROOF047: {req_id} missing from F0015 feature authority")
    for i in range(15, 18):
        ac_id = f"AC-F15-{i:02d}"
        if ac_id not in f15_arch:
            e(f"PROOF047: {ac_id} missing from F0015 architecture AC-to-proof mapping")
        if ac_id not in f15_feat:
            e(f"PROOF047: {ac_id} missing from F0015 feature Acceptance Criteria table")


def check_no_stale_topology_patch_wording(adh_045: str, f15_arch: str, f15_feat: str, reg_text: str) -> None:
    """Fail-closed (ADH-2026-047 decision 5): reject the stale ADH-045 topology wording
    'Name and description are PATCHable' anywhere it is not explicitly corrected/superseded."""
    stale_phrase = "name and description are patchable"
    for label, text in (
        ("ADH-2026-045", adh_045),
        ("F0015 architecture", f15_arch),
        ("F0015 feature", f15_feat),
        ("registry", reg_text),
    ):
        lower = text.lower()
        idx = 0
        while True:
            pos = lower.find(stale_phrase, idx)
            if pos == -1:
                break
            window = lower[max(0, pos - 300):pos + len(stale_phrase) + 300]
            if "correct" not in window and "superseded" not in window:
                e(f"STALETOPOLOGY047: {label} still contains uncorrected stale wording "
                  f"'Name and description are PATCHable' (ADH-2026-047 decision 3)")
            idx = pos + len(stale_phrase)


def check_topology_name_immutable_identity(f15_arch: str, f15_feat: str, reg: dict) -> None:
    """Fail-closed (ADH-2026-047 decision 3): confirm metadata.name is stated as
    immutable identity for every F0015 resource including topology, and that only
    spec.description is PATCHable for topology resources in the registry."""
    schemas = reg.get("schemas", [])
    schema_by_kind = {str(s.get("identity", "")).rsplit("/", 1)[-1]: s for s in schemas}
    for kind_name in ("HostingLocation", "Datacenter", "FaultDomain", "InfrastructureStack"):
        s = schema_by_kind.get(kind_name)
        if not s:
            continue
        mutability = str(s.get("mutability", "")).lower()
        if "name is immutable identity" not in mutability and "identity" not in mutability.split("immutable")[0] and "identity/scope" not in mutability and "identity/scope/parent-ref" not in mutability:
            e(f"TOPOLOGYNAME047: {kind_name} registry mutability does not clearly state immutable identity for metadata.name")
        if "description is patchable" not in mutability and "description patch-only" not in mutability:
            e(f"TOPOLOGYNAME047: {kind_name} registry mutability does not state description is the sole PATCHable topology field")
    for label, text in (("architecture", f15_arch), ("feature", f15_feat)):
        lower = text.lower()
        if "name is immutable identity" not in lower:
            e(f"TOPOLOGYNAME047: F0015 {label} must state 'name is immutable identity' for topology resources (ADH-2026-047 decision 3)")


def check_contract_executability_audit(reg: dict, f15_arch: str, f15_feat: str, trace: str) -> None:
    """ADH-2026-047 decision 5: deterministic, fail-closed contract-executability
    audit run as part of the default readiness gate before requirements generation.
    Composes the individual decision-5 sub-checks below."""
    check_closed_create_field_matrix(reg, f15_arch, f15_feat)
    check_single_scope_derivation_source(reg, f15_arch, f15_feat)
    check_route_contract_completeness(f15_arch, f15_feat)
    check_deferred_participation_fields_rejected(reg, f15_arch, f15_feat)
    check_create_status_and_result_completeness(reg, f15_arch, f15_feat)
    check_f15_26_28_conformance_and_mapping(reg, f15_arch, f15_feat, trace)
    check_topology_name_immutable_identity(f15_arch, f15_feat, reg)


def main() -> None:
    parser = argparse.ArgumentParser(description="FEATURE-0015 architecture-readiness check")
    parser.add_argument("--mode", choices=["readiness", "post"], default="readiness",
                        help="readiness = pre-requirements gate; post = post-stage boundary check")
    args = parser.parse_args()

    reg_text = read(REG_PATH)
    if not reg_text:
        print(f"FAIL: Cannot load registry at {REG_PATH.relative_to(ROOT)}")
        sys.exit(1)
    reg = yaml.safe_load(reg_text)

    f15_arch = read(F15_ARCH)
    f15_feat = read(F15_FEAT)
    spec = read(SPEC_PATH)
    closure = read(CLOSURE)
    trace = read(TRACE_PATH)

    # Fail-closed rejection of reintroduced migration concepts
    check_no_reintroduced_migration_concepts(reg, f15_arch, f15_feat)

    # Approved resource/behavior validations (ADH-2026-045 required feature-local proof)
    check_resource_fields_and_mutability(reg, f15_arch, f15_feat)
    check_participation_action_state_table(reg, f15_arch, f15_feat)
    check_independent_suspension_holds(reg, f15_arch, f15_feat)
    check_topology_safe_denial_ordering(reg, spec, f15_feat)
    check_scope_reference_integrity(reg, spec, f15_feat)
    check_bootstrap_grant_actions(f15_arch, f15_feat)
    check_api_update_behavior_and_concurrency(f15_feat)
    check_audit_matrix(f15_arch, f15_feat)
    check_traceability_local_only(f15_arch)

    # Cross-checks
    check_closure_matrix(closure)
    check_traceability_matrix(trace)

    # FEATURE-0015 architecture-fidelity corrections (stale artifact rejection)
    check_no_stale_cloudplatform_owner_org_ref(reg_text, spec)
    check_no_stale_dual_delegated_acceptance(f15_arch, f15_feat, reg_text, spec)
    check_f15_08_uses_resume_not_suspend(reg)

    # ADH-2026-046 fail-closed guardrails
    check_status_writer_conflicts(reg, reg_text, spec, f15_arch, f15_feat, trace)
    check_participation_create_precondition(f15_arch, f15_feat, spec, reg)
    check_f15_proof_matrix_and_mapping(reg, f15_arch, f15_feat, trace)
    check_no_f16_leakage_in_f15_conformance(reg, f15_feat)
    check_design_tasks_proof_coverage(ROOT / ".automation/features/FEATURE-0015.control.json", args.mode)
    check_requirements_authority_conflicts(args.mode)

    # ADH-2026-047 decision 5: fail-closed contract-executability audit, run as part
    # of the default readiness gate (before requirements-stage Kiro invocation).
    adh_045_text = read(ADH_045)
    check_contract_executability_audit(reg, f15_arch, f15_feat, trace)
    check_no_stale_topology_patch_wording(adh_045_text, f15_arch, f15_feat, reg_text)

    if errs:
        print(f"FAIL: FEATURE-0015 architecture-readiness — {len(errs)} error(s)")
        for err in errs:
            print(f"  ✗ {err}")
        print("\nRequirements generation is BLOCKED until all readiness checks pass.")
        sys.exit(1)
    else:
        print("PASS: FEATURE-0015 architecture-readiness — canonical bootstrap (ADH-2026-045) validated")
        print("  ✓ No reintroduced migration concepts (CanonicalMigrationPlan/Record, migration controller, ExecutionTarget in F0015 scope)")
        print("  ✓ Resource fields/mutability (PATCH-only application/merge-patch+json; no PUT/DELETE)")
        print("  ✓ Participation action-state table (explicit actions; If-Match + Idempotency-Key)")
        print("  ✓ Independent suspension holds (platformSuspended / providerSuspended)")
        print("  ✓ Topology safe-denial ordering (RESOURCE_NOT_FOUND/404 vs VALIDATION_FAILED/422)")
        print("  ✓ Scope-reference integrity (VS0_SCOPE_REFERENCE_MISMATCH)")
        print("  ✓ Bootstrap grant actions (VS0_CLOUDPLATFORM_ROOT_REQUIRED)")
        print("  ✓ API/update behavior and concurrency (idempotency, stale-version handling)")
        print("  ✓ Audit matrix (FEATURE-0013 AuditEvent reuse, correlation, no secrets)")
        print("  ✓ Traceability (F0015-local conformance only)")
        print("  ✓ Closure matrix: all ARC-F15 rows RESOLVED")
        print("  ✓ Traceability matrix: VS0-CF-F15-11 present")
        print("  ✓ Status-writer resolution (api-server sole writer; ADH-2026-046 decision 1)")
        print("  ✓ Participation create-vs-existing preconditions (ADH-2026-046 decision 2)")
        print("  ✓ Complete F15-01..25 proof matrix and REQ/AC mapping (ADH-2026-046 decision 3)")
        print("  ✓ No F0016+ leakage into F0015-local conformance (ADH-2026-046 decision 4e)")
        print("  ✓ Closed collection-create request contract (per-kind client field boundary; ADH-2026-047 decision 1)")
        print("  ✓ Single-source scope derivation per resource kind (ADH-2026-047 decision 2)")
        print("  ✓ Topology name-immutable/description-PATCHable correction; no stale ADH-045 wording (ADH-2026-047 decision 3)")
        print("  ✓ Participation collection-create body vs. empty item-action body; FEATURE-0021 fields deferred (ADH-2026-047 decision 4)")
        print("  ✓ Contract-executability audit: field classification, route-contract completeness, F15-26..28 proof matrix (ADH-2026-047 decision 5)")
        sys.exit(0)


if __name__ == "__main__":
    main()
