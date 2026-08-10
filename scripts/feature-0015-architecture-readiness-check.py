#!/usr/bin/env python3
"""Fail-closed architecture-readiness check for FEATURE-0015.

Validates that all seven ADH-2026-043 decisions plus the ADH-2026-044 run-local
migration-completion decision are representable and consistent across the
repository authorities before requirements generation may proceed.

Exit 0 = PASS (requirements generation may proceed).
Exit 1 = FAIL (architecture gap remains; requirements generation blocked).
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

# --- Paths ---
REG_PATH = ROOT / "docs/architecture/vertical-slices/VS-000-contract-registry.yaml"
SPEC_PATH = ROOT / "docs/architecture/vertical-slices/VS-000-contract-specification.md"
TRACE_PATH = ROOT / "docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md"
F15_ARCH = ROOT / "docs/architecture/FEATURE-0015-canonical-cloud-model-and-alpha-migration-foundation.md"
F15_FEAT = ROOT / "docs/features/FEATURE-0015-canonical-cloud-model-and-alpha-migration-foundation.md"
CLOSURE = ROOT / "docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md"
PHASE_CTX = ROOT / "docs/context/CURRENT_PHASE_CONTEXT.md"
STEER = ROOT / ".kiro/steering/slice0-contract.md"

errs: list[str] = []


def e(msg: str) -> None:
    errs.append(msg)


def read(path: Path) -> str:
    if not path.exists():
        e(f"Required file missing: {path.relative_to(ROOT)}")
        return ""
    return path.read_text()


def check_decision_1_correction_model(reg: dict, f15_arch: str, f15_feat: str) -> None:
    """Decision 1: CanonicalMigrationRecord has no direct correction link;
    correction requires a superseding plan and new migration run."""
    # Registry: VS0-SCHEMA-061 must not have a correctionRef or supersessionRef field
    schemas = reg.get("schemas", [])
    schema_061 = next((s for s in schemas if s.get("id") == "VS0-SCHEMA-061"), None)
    if not schema_061:
        e("D1: VS0-SCHEMA-061 (CanonicalMigrationRecord) not found in registry")
        return
    all_fields = " ".join(schema_061.get("required", []) + schema_061.get("optional", []))
    if "correctionRef" in all_fields or "supersessionRef" in all_fields:
        e("D1: CanonicalMigrationRecord must not have correctionRef or supersessionRef fields")
    # Architecture doc must state the correction model
    if "superseding CanonicalMigrationPlan" not in f15_arch:
        e("D1: F0015 architecture must state corrections require a superseding plan")
    if "never altered or re-ordered" not in f15_arch:
        e("D1: F0015 architecture must state prior records are never altered or re-ordered")
    # Feature file
    if "corrections create a new linked superseding plan" not in f15_feat:
        e("D1: F0015 feature file must describe correction via superseding plan")


def check_decision_2_draft(reg: dict, f15_arch: str, f15_feat: str) -> None:
    """Decision 2: Draft is pre-persistence preparation; signed plan precedes
    persisted 1..9 milestone sequence."""
    sms = reg.get("stateMachines", [])
    sm_011 = next((sm for sm in sms if sm.get("id") == "VS0-STATE-011"), None)
    if not sm_011:
        e("D2: VS0-STATE-011 not found")
        return
    milestones = sm_011.get("milestones", [])
    if milestones and milestones[0] != "InventoryValidated":
        e("D2: VS0-STATE-011 must begin at InventoryValidated, not Draft")
    if "Draft" in milestones:
        e("D2: Draft must not appear in the persisted milestone sequence")
    # Architecture doc
    if "pre-persistence" not in f15_arch.lower():
        e("D2: F0015 architecture must state Draft is pre-persistence preparation")
    # Feature file
    if "pre-persistence" not in f15_feat.lower():
        e("D2: F0015 feature must state Draft is pre-persistence preparation")


def check_decision_3_inventory(f15_arch: str, f15_feat: str) -> None:
    """Decision 3: F0015 owns provider/topology transforms only; defers others;
    final cutover is FEATURE-0026."""
    if "FEATURE-0026" not in f15_arch or "cutover" not in f15_arch.lower():
        e("D3: F0015 architecture must attribute final cutover to FEATURE-0026")
    if "defers transform" not in f15_feat.lower() and "defers transform" not in f15_arch.lower():
        e("D3: F0015 must defer non-provider transforms to later features")
    # Must classify but not execute catalog/enrollment/governance transforms
    for domain in ("ServiceClass", "CloudEnrollment", "Governance"):
        if domain not in f15_arch and domain not in f15_feat:
            e(f"D3: F0015 must classify {domain} in the signed plan (deferred)")


def check_decision_4_safe_denial(reg: dict, spec: str, f15_feat: str) -> None:
    """Decision 4: inaccessible cross-provider = RESOURCE_NOT_FOUND/404 + VS0_AUTHORIZATION_SAFE_DENIAL;
    authorized invalid same-provider = VALIDATION_FAILED/422."""
    # Check VS0-CF-X03 conformance entry
    confs = reg.get("conformance", [])
    x03 = next((c for c in confs if c.get("id") == "VS0-CF-X03"), None)
    if not x03:
        e("D4: VS0-CF-X03 conformance entry missing")
        return
    if x03.get("expectedError") != "RESOURCE_NOT_FOUND":
        e("D4: VS0-CF-X03 expectedError must be RESOURCE_NOT_FOUND")
    # VS0_AUTHORIZATION_SAFE_DENIAL in violation codes
    vcs = {vc.get("code") for vc in reg.get("violationCodes", {}).get("slice0", [])}
    if "VS0_AUTHORIZATION_SAFE_DENIAL" not in vcs:
        e("D4: VS0_AUTHORIZATION_SAFE_DENIAL must be a registered violation code")
    # Spec mentions safe denial semantics
    if "VS0_AUTHORIZATION_SAFE_DENIAL" not in spec:
        e("D4: VS-000 specification must reference VS0_AUTHORIZATION_SAFE_DENIAL safe-denial rule")
    # Feature file error table
    if "RESOURCE_NOT_FOUND" not in f15_feat or "VS0_AUTHORIZATION_SAFE_DENIAL" not in f15_feat:
        e("D4: F0015 error table must include RESOURCE_NOT_FOUND + VS0_AUTHORIZATION_SAFE_DENIAL")


def check_decision_5_scope_integrity(reg: dict, spec: str, f15_feat: str) -> None:
    """Decision 5: UID-pinned scope-reference invariants with VS0_SCOPE_REFERENCE_MISMATCH."""
    # Violation code registered
    vcs = {vc.get("code") for vc in reg.get("violationCodes", {}).get("slice0", [])}
    if "VS0_SCOPE_REFERENCE_MISMATCH" not in vcs:
        e("D5: VS0_SCOPE_REFERENCE_MISMATCH must be a registered violation code")
    # Conformance VS0-CF-F15-11
    confs = reg.get("conformance", [])
    f15_11 = next((c for c in confs if c.get("id") == "VS0-CF-F15-11"), None)
    if not f15_11:
        e("D5: VS0-CF-F15-11 conformance entry missing")
    elif f15_11.get("expectedError") != "VALIDATION_FAILED":
        e("D5: VS0-CF-F15-11 expectedError must be VALIDATION_FAILED")
    elif f15_11.get("expectedViolation") != "VS0_SCOPE_REFERENCE_MISMATCH":
        e("D5: VS0-CF-F15-11 expectedViolation must be VS0_SCOPE_REFERENCE_MISMATCH")
    # Spec
    if "VS0_SCOPE_REFERENCE_MISMATCH" not in spec:
        e("D5: VS-000 specification must reference VS0_SCOPE_REFERENCE_MISMATCH")
    # Feature file
    if "VS0_SCOPE_REFERENCE_MISMATCH" not in f15_feat:
        e("D5: F0015 feature must reference VS0_SCOPE_REFERENCE_MISMATCH")
    if "ownerOrganizationRef" not in f15_feat:
        e("D5: F0015 feature must state ownerOrganizationRef UID invariant")


def check_decision_6_audit(f15_arch: str, f15_feat: str) -> None:
    """Decision 6: F0015 reuses FEATURE-0013 AuditEvent with correlation/redaction/no-secrets."""
    if "FEATURE-0013" not in f15_feat or "AuditEvent" not in f15_feat:
        e("D6: F0015 feature must reference FEATURE-0013 AuditEvent reuse")
    if "no secrets" not in f15_feat.lower() and "no secret" not in f15_feat.lower():
        e("D6: F0015 feature must state no-secrets rule for audit")
    if "correlation" not in f15_feat.lower():
        e("D6: F0015 feature must state correlation requirements")
    if "VS0-CF-T01" not in f15_feat:
        e("D6: F0015 feature must explicitly not import VS0-CF-T01")
    # Architecture doc
    if "AuditEvent" not in f15_arch:
        e("D6: F0015 architecture must reference AuditEvent obligations")
    if "no secrets" not in f15_arch.lower() and "no secret" not in f15_arch.lower():
        e("D6: F0015 architecture must state no-secrets rule")


def check_decision_7_traceability(f15_arch: str) -> None:
    """Decision 7: F0015 traceability uses only local conformance IDs;
    downstream HP01/F09 are not local acceptance evidence."""
    # §10 must not use VS0-CF-HP01 or VS0-CF-F09 as FEATURE-0015 acceptance evidence
    # The traceability table is after "## 10." heading
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

    # Check that the table column header mentions "local"
    if "FEATURE-0015 local" not in sec10 and "F0015 local" not in sec10:
        e("D7: §10 traceability must use F0015-local conformance IDs")
    # VS0-CF-HP01 and VS0-CF-F09 must not be in the table's conformance column
    # They can appear in a note below explaining they are downstream
    table_lines = [l for l in sec10.splitlines() if l.strip().startswith("|") and "VS0-CF" in l]
    for tl in table_lines:
        cells = [c.strip() for c in tl.split("|")]
        # Last populated cell is conformance
        if len(cells) >= 6:
            conf_cell = cells[-2] if cells[-1] == "" else cells[-1]
            if "HP01" in conf_cell and "note" not in tl.lower():
                e("D7: §10 table must not use VS0-CF-HP01 as local conformance evidence")
            if conf_cell.strip() == "VS0-CF-F09" or ", F09" in conf_cell or ",F09" in conf_cell:
                e("D7: §10 table must not use VS0-CF-F09 as local conformance evidence")

    # Must have a note about downstream
    if "downstream" not in sec10.lower() and "non-owning" not in sec10.lower():
        e("D7: §10 must note that HP01/F09 are downstream integration references only")


def check_adh_044_run_cardinality(reg: dict, spec: str, f15_arch: str, f15_feat: str) -> None:
    """ADH-2026-044: migration completion is run-local under a plan-declared runKey.

    Fail-closed: rejects missing runKey linkage, missing per-(planRef.uid, runKey)
    cardinality, missing run-local Completed, or a global-completion claim by F0015.
    """
    schemas = reg.get("schemas", [])
    plan = next((s for s in schemas if s.get("id") == "VS0-SCHEMA-060"), None)
    record = next((s for s in schemas if s.get("id") == "VS0-SCHEMA-061"), None)
    if not plan:
        e("D044: VS0-SCHEMA-060 (CanonicalMigrationPlan) not found in registry")
    else:
        plan_text = " ".join(plan.get("required", [])) + " " + str(plan.get("validation", ""))
        if "runKey" not in plan_text:
            e("D044: CanonicalMigrationPlan must declare a plan-level runKey per resourceTransforms entry")
    if not record:
        e("D044: VS0-SCHEMA-061 (CanonicalMigrationRecord) not found in registry")
    else:
        if not any(str(f).startswith("record.runKey:") for f in record.get("required", [])):
            e("D044: CanonicalMigrationRecord must require record.runKey")
        if "runKey" not in str(record.get("mutability", "")):
            e("D044: CanonicalMigrationRecord mutability must bind records to a plan-declared runKey")

    sms = reg.get("stateMachines", [])
    sm_011 = next((sm for sm in sms if sm.get("id") == "VS0-STATE-011"), {})
    state_text = " ".join(str(sm_011.get(k, "")) for k in ("description", "ordering", "immutabilityRule"))
    if "runKey" not in state_text:
        e("D044: VS0-STATE-011 must express per-(planRef.uid, runKey) cardinality")
    if "planRef.uid" not in state_text:
        e("D044: VS0-STATE-011 must key the milestone chain by planRef.uid and runKey")
    guard_text = " ".join(sm_011.get("guards", [])) + " " + str(sm_011.get("ordering", ""))
    if "seals" not in guard_text.lower() and "run-local" not in state_text.lower():
        e("D044: VS0-STATE-011 must state that Completed seals only its run (run-local)")

    confs = reg.get("conformance", [])
    mig01 = next((c for c in confs if c.get("id") == "VS0-CF-MIG01"), {})
    mig01_text = " ".join(str(mig01.get(k, "")) for k in ("inputs", "expectedState", "expectedSideEffects"))
    if "provider-topology" not in mig01_text:
        e("D044: VS0-CF-MIG01 must prove the provider-topology run, not global completion")
    if "FEATURE-0026" not in mig01_text:
        e("D044: VS0-CF-MIG01 must defer global all-run completion to FEATURE-0026")

    # Architecture authority
    if "runKey" not in f15_arch:
        e("D044: F0015 architecture must define plan-declared runKey semantics")
    if "run-local" not in f15_arch.lower() and "seals that run" not in f15_arch.lower():
        e("D044: F0015 architecture must state Completed is run-local")
    if "provider-topology" not in f15_arch.lower():
        e("D044: F0015 architecture must scope F0015 to the provider-topology run")
    # CloudPlatform writer traceability must be VS0-WRITER-002 only
    for line in f15_arch.splitlines():
        if line.strip().startswith("| CloudPlatform ") and "VS0-WRITER-002,003" in line:
            e("D044: CloudPlatform writer traceability must be VS0-WRITER-002 only")

    # Feature authority
    if "runKey" not in f15_feat:
        e("D044: F0015 feature must define plan-declared runKey semantics")
    if "provider-topology" not in f15_feat.lower():
        e("D044: F0015 feature must scope F0015 to the provider-topology run")
    if "FEATURE-0026" not in f15_feat:
        e("D044: F0015 feature must attribute global completion proof to FEATURE-0026")


def check_closure_matrix(closure: str) -> None:
    """All ARC-F15-01..07 must be RESOLVED."""
    for i in range(1, 8):
        arc_id = f"ARC-F15-{i:02d}"
        if arc_id not in closure:
            e(f"Closure: {arc_id} not found in closure matrix")
        elif "RESOLVED" not in closure.split(arc_id)[1].split("\n")[0]:
            e(f"Closure: {arc_id} is not RESOLVED")


def check_traceability_matrix(trace: str) -> None:
    """VS0-CF-F15-11 must appear in the traceability matrix."""
    if "VS0-CF-F15-11" not in trace:
        e("Traceability: VS0-CF-F15-11 missing from traceability matrix")


def main() -> None:
    parser = argparse.ArgumentParser(description="FEATURE-0015 architecture-readiness check")
    parser.add_argument("--mode", choices=["readiness", "post"], default="readiness",
                        help="readiness = pre-requirements gate; post = post-stage boundary check")
    parser.parse_args()

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

    # Run all seven decision checks
    check_decision_1_correction_model(reg, f15_arch, f15_feat)
    check_decision_2_draft(reg, f15_arch, f15_feat)
    check_decision_3_inventory(f15_arch, f15_feat)
    check_decision_4_safe_denial(reg, spec, f15_feat)
    check_decision_5_scope_integrity(reg, spec, f15_feat)
    check_decision_6_audit(f15_arch, f15_feat)
    check_decision_7_traceability(f15_arch)
    check_adh_044_run_cardinality(reg, spec, f15_arch, f15_feat)

    # Cross-checks
    check_closure_matrix(closure)
    check_traceability_matrix(trace)

    if errs:
        print(f"FAIL: FEATURE-0015 architecture-readiness — {len(errs)} error(s)")
        for err in errs:
            print(f"  ✗ {err}")
        print("\nRequirements generation is BLOCKED until all readiness checks pass.")
        sys.exit(1)
    else:
        print("PASS: FEATURE-0015 architecture-readiness — all 7 ADH-2026-043 decisions validated")
        print("  ✓ D1: Correction model (superseding plan, no direct record link)")
        print("  ✓ D2: Draft representation (pre-persistence, not a milestone)")
        print("  ✓ D3: Migration inventory boundary (F0015 provider/topology only)")
        print("  ✓ D4: Safe-denial semantics (RESOURCE_NOT_FOUND + VS0_AUTHORIZATION_SAFE_DENIAL)")
        print("  ✓ D5: Scope-reference UID invariants (VS0_SCOPE_REFERENCE_MISMATCH)")
        print("  ✓ D6: Audit evidence (FEATURE-0013 reuse, correlation, no secrets)")
        print("  ✓ D7: Traceability (F0015-local conformance only)")
        print("  ✓ D044: Run-local migration completion (plan-declared runKey, per-(planRef.uid,runKey) chain, FEATURE-0026 global proof)")
        print("  ✓ Closure matrix: all ARC-F15 rows RESOLVED")
        print("  ✓ Traceability matrix: VS0-CF-F15-11 present")
        sys.exit(0)


if __name__ == "__main__":
    main()
