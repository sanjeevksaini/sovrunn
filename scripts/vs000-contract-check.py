#!/usr/bin/env python3
"""VS-000 Contract Registry deterministic checker. Exits 1 on any failure."""
import json, os, re, sys, subprocess
from pathlib import Path
try:
    import yaml
except ImportError:
    print("FAIL: PyYAML required. Install: pip install pyyaml", file=sys.stderr); sys.exit(1)

REPO = Path(__file__).resolve().parent.parent
REG_PATH = REPO/"docs/architecture/vertical-slices/VS-000-contract-registry.yaml"
TRACE_PATH = REPO/"docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md"
STEER_PATH = REPO/".kiro/steering/slice0-contract.md"
CHARTER_PATH = REPO/"docs/architecture/vertical-slices/VS-000-core-skeleton.md"
SPEC_PATH = REPO/"docs/architecture/vertical-slices/VS-000-contract-specification.md"
SEQUENCE_PATH = REPO/"docs/phase2/PHASE2_FEATURE_SEQUENCE.md"
F15_ARCH_PATH = REPO/"docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md"
F15_FEATURE_PATH = REPO/"docs/features/FEATURE-0015-canonical-cloud-model-foundation.md"
F15_CONTROL_PATH = REPO/".automation/features/FEATURE-0015.control.json"
KIRO_AGENT_PATH = REPO/".kiro/agents/sovrunn-spec.md"
SCOPES = {"Platform","Organization","OrganizationUnit","Tenant","Project","CloudPlatform","CloudProvider"}
F12_CODES = {"MALFORMED_REQUEST":(400,"malformed-request"),"UNKNOWN_FIELD":(400,"unknown-field"),
    "DUPLICATE_FIELD":(400,"duplicate-field"),"REQUEST_TOO_LARGE":(400,"request-too-large"),
    "AUTH_REQUIRED":(401,"auth-required"),"AUTHORIZATION_DENIED":(403,"authorization-denied"),
    "RESOURCE_NOT_FOUND":(404,"resource-not-found"),"CONFLICT":(409,"conflict"),
    "ALREADY_EXISTS":(409,"already-exists"),"DELETE_BLOCKED":(409,"delete-blocked"),
    "STALE_RESOURCE_VERSION":(412,"stale-resource-version"),"UNSUPPORTED_MEDIA_TYPE":(415,"unsupported-media-type"),
    "VALIDATION_FAILED":(422,"validation-failed"),"INTERNAL_ERROR":(500,"internal-error"),
    "DEPENDENCY_UNAVAILABLE":(503,"dependency-unavailable")}
INVENTED = {"NOT_FOUND","FORBIDDEN","UNAUTHORIZED","QUOTA_EXCEEDED","INTERNAL","UNAVAILABLE"}
PROHIBITED = ["ResourcePool","ProviderCapability","generic Provider as combined owner/operator",
    "ServiceClass as canonical catalog","EffectivePolicyContext","provider-neutral","six-scope authority",
    "CanonicalMigrationPlan","CanonicalMigrationRecord","migration-controller","approved-migration-plan-publisher"]
CF_IDS = (["HP01"]+[f"F{i:02d}" for i in range(1,21)]
    +["X01","X02","X03","L01","Z01","T01","I01","I02","D01"])
F15_CF_IDS = [f"F15-{i:02d}" for i in range(1,42)]
F16_CF_IDS = [f"F16-{i:02d}" for i in range(1,129)]
F15_OWNED = {
    "CloudPlatform", "CloudProvider", "CloudProviderParticipation", "HostingLocation",
    "Datacenter", "FaultDomain", "InfrastructureStack",
}
F15_CONTROL_EXPECTED_OWNED = F15_OWNED
RETIRED_SCHEMA_IDS = ["VS0-SCHEMA-057", "VS0-SCHEMA-060", "VS0-SCHEMA-061"]
RETIRED_WRITER_IDS = ["VS0-WRITER-020", "VS0-WRITER-021"]
RETIRED_STATE_IDS = ["VS0-STATE-011"]
RETIRED_CONFORMANCE_IDS = ["VS0-CF-MIG01", "VS0-CF-MIG02", "VS0-CF-MIGF01", "VS0-CF-MIGF02", "VS0-CF-MIGF03"]
RETIRED_MIGRATION_FAILURE_IDS = ["VS0-MIG-F01", "VS0-MIG-F02", "VS0-MIG-F03"]
errs = []
def e(m): errs.append(m)

def kind(schema):
    return str(schema.get("identity", "")).rsplit("/", 1)[-1].rsplit(":", 1)[-1]

def run():
    if not REG_PATH.exists(): e(f"Registry not found: {REG_PATH}"); return
    reg = yaml.safe_load(REG_PATH.read_text())
    # Metadata
    m = reg.get("metadata",{})
    for k in ("id","version","status","scope","baseline","controllingHandoff"):
        if not m.get(k): e(f"metadata.{k} missing")
    if m.get("scope")!="Slice0Only": e(f"metadata.scope != Slice0Only")
    if set(m.get("scopeKinds",[]))!=SCOPES: e("metadata.scopeKinds != seven canonical scopes")
    # Schemas
    schemas = reg.get("schemas",[])
    exp_s = [f"VS0-SCHEMA-{i:03d}" for i in range(1,62) if f"VS0-SCHEMA-{i:03d}" not in RETIRED_SCHEMA_IDS]
    ids_s = [s["id"] for s in schemas]
    if ids_s!=exp_s: e(f"Active schema IDs must be 001..061 minus retired IDs, in order (got {len(ids_s)})")
    retired = reg.get("retiredSchemas", [])
    exp_retired = RETIRED_SCHEMA_IDS
    ids_retired = [s.get("id") for s in retired]
    if ids_retired != exp_retired: e(f"retiredSchemas must contain exactly {exp_retired} tombstones in order")
    if sorted(ids_s+ids_retired) != [f"VS0-SCHEMA-{i:03d}" for i in range(1,62)]:
        e("Active plus retired schema IDs must consume exactly 001..061")
    for tomb in retired:
        tid = tomb.get("id", "?")
        if tomb.get("status") != "tombstone": e(f"{tid}: retired status must be tombstone")
        if tomb.get("owner"): e(f"{tid}: retired tombstone must not have an active owner")
        if not tomb.get("removalReason") or not tomb.get("retiredBy"): e(f"{tid}: retirement traceability incomplete")
    seen=set()
    for s in schemas:
        sid=s.get("id","?")
        if sid in seen: e(f"Dup schema: {sid}")
        seen.add(sid)
        for f in ("identity","owner","profile","boundary"):
            if not s.get(f): e(f"{sid}: missing {f}")
        sc=s.get("scopes")
        if not s.get("externalRef") and not s.get("required"): e(f"{sid}: need externalRef or required")
        if isinstance(sc,list) and sc:
            for v in sc:
                if v not in SCOPES: e(f"{sid}: bad scope '{v}'")
    # Writers
    writers = reg.get("writers",[])
    exp_w = [f"VS0-WRITER-{i:03d}" for i in range(1,20)]
    ids_w = [w["id"] for w in writers]
    if ids_w!=exp_w: e(f"Writer IDs not contiguous 001..019")
    pmap={}
    for w in writers:
        wid=w.get("id","?")
        if not w.get("paths"): e(f"{wid}: paths empty")
        if not w.get("writer"): e(f"{wid}: writer empty")
        if "forbidden" not in w: e(f"{wid}: no forbidden")
        if "conflict" not in w: e(f"{wid}: no conflict")
        for p in w.get("paths",[]):
            if p in pmap: e(f"Path '{p}' in both {pmap[p]} and {wid}")
            pmap[p]=wid
    retired_writers = reg.get("retiredWriters", [])
    ids_retired_w = [w.get("id") for w in retired_writers]
    if ids_retired_w != RETIRED_WRITER_IDS: e(f"retiredWriters must contain exactly {RETIRED_WRITER_IDS} in order")
    for tomb in retired_writers:
        if not tomb.get("removalReason") or not tomb.get("retiredBy"): e(f"{tomb.get('id','?')}: retirement traceability incomplete")
    # State machines
    sms = reg.get("stateMachines",[])
    exp_sm = [f"VS0-STATE-{i:03d}" for i in range(1,11)]
    ids_sm = [x["id"] for x in sms]
    if ids_sm!=exp_sm: e("StateMachine IDs not contiguous 001..010")
    retired_sms = reg.get("retiredStateMachines", [])
    ids_retired_sm = [x.get("id") for x in retired_sms]
    if ids_retired_sm != RETIRED_STATE_IDS: e(f"retiredStateMachines must contain exactly {RETIRED_STATE_IDS} in order")
    for tomb in retired_sms:
        if not tomb.get("removalReason") or not tomb.get("retiredBy"): e(f"{tomb.get('id','?')}: retirement traceability incomplete")
    for sm in sms:
        mid=sm.get("id","?")
        sts=sm.get("states",[]); ini=sm.get("initial"); terms=sm.get("terminal",[])
        if not sts: e(f"{mid}: no states")
        if ini and ini not in sts: e(f"{mid}: initial not in states")
        for t in terms:
            if t not in sts: e(f"{mid}: terminal '{t}' not in states")
        for tr in sm.get("transitions",[]):
            if not isinstance(tr,str): continue
            arrow=tr.split(":",1)[0]
            if "->" in arrow:
                sp,tgt=arrow.split("->",1)
                for src in sp.split("|"):
                    if src.strip() not in sts: e(f"{mid}: src '{src.strip()}' not in states")
                if tgt.strip() not in sts: e(f"{mid}: tgt '{tgt.strip()}' not in states")
        if mid=="VS0-STATE-010":
            if not sm.get("appendOnly"): e("VS0-STATE-010: appendOnly must be true")
            if sts!=["FINAL"]: e("VS0-STATE-010: states must be [FINAL]")
        if mid=="VS0-STATE-001":
            hold_rule = str(sm.get("holdRule","")).lower()
            for required in ("platformsuspended", "providersuspended", "clearing one hold"):
                if required not in hold_rule: e(f"VS0-STATE-001: holdRule missing '{required}'")
    # Problem codes
    pcs = reg.get("problemCodes",[])
    if len(pcs)!=15: e(f"problemCodes count {len(pcs)} != 15")
    found_pc={pc["code"] for pc in pcs}
    for pc in pcs:
        c,h,t=pc.get("code"),pc.get("http"),pc.get("type","")
        if c not in F12_CODES: e(f"Unknown problem code: {c}"); continue
        eh,slug=F12_CODES[c]; exp_urn=f"urn:sovrunn:problem:{slug}"
        if h!=eh: e(f"{c}: http {h}!={eh}")
        if t!=exp_urn: e(f"{c}: type mismatch")
        if t and not re.match(r"^urn:sovrunn:problem:[a-z][a-z0-9-]*$",t): e(f"{c}: URN not kebab")
    for c in F12_CODES:
        if c not in found_pc: e(f"Missing problem code: {c}")
    # Violation codes
    vcs = reg.get("violationCodes",{}).get("slice0",[])
    vc_set=set()
    for vc in vcs:
        code=vc.get("code","")
        if not code.startswith("VS0_"): e(f"Violation not VS0_ prefixed: {code}")
        if code in vc_set: e(f"Dup violation: {code}")
        vc_set.add(code)
        if "field" not in vc: e(f"{code}: missing field (use null)")
    # Failure mappings
    fms = reg.get("failureMappings",[])
    exp_f=[f"VS0-F{i:02d}" for i in range(1,21)]
    ids_f=[f["id"] for f in fms]
    if ids_f!=exp_f: e("failureMapping IDs != F01..F20")
    for fm in fms:
        fid=fm.get("id","?"); num=fid.replace("VS0-F","")
        if fm.get("conformance")!=f"VS0-CF-F{num}": e(f"{fid}: conformance mismatch")
        code=fm.get("code"); viol=fm.get("violation")
        if code and code not in found_pc: e(f"{fid}: code '{code}' not in problemCodes")
        if code and code in F12_CODES:
            eh,slug=F12_CODES[code]
            if fm.get("http")!=eh: e(f"{fid}: http mismatch")
            if fm.get("type")!=f"urn:sovrunn:problem:{slug}": e(f"{fid}: type mismatch")
        if viol and viol not in vc_set: e(f"{fid}: violation '{viol}' not in slice0 codes")
    # Migration failure mappings are retired: canonical bootstrap replaces alpha runtime migration (DEC-0059).
    retired_mig_fms = reg.get("retiredMigrationFailureMappings", [])
    ids_retired_mig_fms = [x.get("id") for x in retired_mig_fms]
    if ids_retired_mig_fms != RETIRED_MIGRATION_FAILURE_IDS: e(f"retiredMigrationFailureMappings must contain exactly {RETIRED_MIGRATION_FAILURE_IDS} in order")
    for tomb in retired_mig_fms:
        if not tomb.get("removalReason") or not tomb.get("retiredBy"): e(f"{tomb.get('id','?')}: retirement traceability incomplete")
    if reg.get("migrationFailureMappings"): e("migrationFailureMappings must not exist as an active section (DEC-0059)")
    # Conformance
    confs = reg.get("conformance",[])
    exp_cf={f"VS0-CF-{c}" for c in CF_IDS + F15_CF_IDS + F16_CF_IDS}; found_cf=set(); seen_cf=set()
    for c in confs:
        cid=c.get("id","?")
        if cid in seen_cf: e(f"Dup conformance: {cid}")
        seen_cf.add(cid); found_cf.add(cid)
        for f in ("owner","inputs","expectedState","expectedError","expectedSideEffects","gate"):
            if f not in c: e(f"{cid}: missing {f}")
    for x in exp_cf-found_cf: e(f"Missing conformance: {x}")
    # ADH-2026-049: FEATURE-0015 required-AuditEvent-append failure must map to the
    # inherited INTERNAL_ERROR/500 outcome, never DEPENDENCY_UNAVAILABLE/503.
    f15_24 = next((c for c in confs if c.get("id") == "VS0-CF-F15-24"), None)
    if not f15_24:
        e("VS0-CF-F15-24 conformance entry missing (ADH-2026-049)")
    else:
        if f15_24.get("expectedError") != "INTERNAL_ERROR":
            e(f"VS0-CF-F15-24 expectedError must be INTERNAL_ERROR (ADH-2026-049), found {f15_24.get('expectedError')!r}")
        side_effects_049 = str(f15_24.get("expectedSideEffects", "")).lower()
        if "dependency_unavailable" in side_effects_049 and "not used" not in side_effects_049:
            e("VS0-CF-F15-24 expectedSideEffects must not assign DEPENDENCY_UNAVAILABLE (ADH-2026-049)")
    # ADH-2026-050: VS0-CF-F15-18/19 must each explicitly cover all seven FEATURE-0015
    # collection creates and all eight existing-participation create/action routes,
    # while preserving their existing IDs, expected states, and violation/code semantics.
    f15_18 = next((c for c in confs if c.get("id") == "VS0-CF-F15-18"), None)
    f15_19 = next((c for c in confs if c.get("id") == "VS0-CF-F15-19"), None)
    if not f15_18:
        e("VS0-CF-F15-18 conformance entry missing (ADH-2026-050)")
    if not f15_19:
        e("VS0-CF-F15-19 conformance entry missing (ADH-2026-050)")
    if f15_18 and f15_19:
        f15_all_route_kinds = ("CloudPlatform", "CloudProvider", "CloudProviderParticipation",
                                "HostingLocation", "Datacenter", "FaultDomain", "InfrastructureStack")
        f15_all_route_actions = ("accept", "reject", "withdraw", "suspend", "resume",
                                  "request-release", "accept-release", "decline-release")
        for case_id, entry in (("VS0-CF-F15-18", f15_18), ("VS0-CF-F15-19", f15_19)):
            inputs_lower = str(entry.get("inputs", "")).lower()
            for kind_name in f15_all_route_kinds:
                if kind_name.lower() not in inputs_lower:
                    e(f"{case_id} inputs must name {kind_name} in the all-route scope (ADH-2026-050)")
            for action in f15_all_route_actions:
                if action not in inputs_lower:
                    e(f"{case_id} inputs must name the '{action}' participation action in the all-route scope (ADH-2026-050)")
        if f15_18.get("expectedError") is not None:
            e("VS0-CF-F15-18 expectedError must remain null (ADH-2026-050)")
        if f15_19.get("expectedError") != "CONFLICT":
            e(f"VS0-CF-F15-19 expectedError must remain CONFLICT (ADH-2026-050), found {f15_19.get('expectedError')!r}")
        if f15_19.get("expectedViolation") != "VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH":
            e(f"VS0-CF-F15-19 expectedViolation must remain VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH (ADH-2026-050), found {f15_19.get('expectedViolation')!r}")
    # ADH-2026-054: replay is bound to the concrete item target and is allowed
    # only after current authorization and safe target/reference access.
    for case_id, entry in (("VS0-CF-F15-18", f15_18), ("VS0-CF-F15-19", f15_19)):
        if not entry:
            continue
        inputs_lower = str(entry.get("inputs", "")).lower()
        for marker in ("concrete cloudproviderparticipation uid", "current server-resolved authorization", "safe target/reference access"):
            if marker not in inputs_lower:
                e(f"{case_id} inputs must state {marker!r} (ADH-2026-054)")
    f15_22 = next((c for c in confs if c.get("id") == "VS0-CF-F15-22"), None)
    f15_30 = next((c for c in confs if c.get("id") == "VS0-CF-F15-30"), None)
    if not f15_22:
        e("VS0-CF-F15-22 conformance entry missing (ADH-2026-054)")
    else:
        inputs_054 = str(f15_22.get("inputs", "")).lower()
        effects_054 = str(f15_22.get("expectedSideEffects", "")).lower()
        for marker in ("x-sovrunn-bootstrap-grant", "bootstrapgrant", "syntactically valid duplicate-free"):
            if marker not in inputs_054:
                e(f"VS0-CF-F15-22 inputs must state {marker!r} (ADH-2026-054)")
        for marker in ("header detection follows authentication", "body-member detection precedes current authorization"):
            if marker not in effects_054:
                e(f"VS0-CF-F15-22 expectedSideEffects must state {marker!r} (ADH-2026-054)")
    if not f15_30:
        e("VS0-CF-F15-30 conformance entry missing (ADH-2026-054)")
    else:
        inputs_054 = str(f15_30.get("inputs", "")).lower()
        effects_054 = str(f15_30.get("expectedSideEffects", "")).lower()
        if "does not contain the reserved top-level bootstrapgrant" not in inputs_054:
            e("VS0-CF-F15-30 inputs must exclude the reserved bootstrapGrant carrier (ADH-2026-054)")
        if "governed instead by vs0-cf-f15-22" not in effects_054:
            e("VS0-CF-F15-30 expectedSideEffects must delegate reserved bootstrapGrant to F15-22 (ADH-2026-054)")
    # ADH-2026-051: the exact durable-audit boundary requires an AuditEvent side effect
    # for every authenticated authorization/safe-denial category it names, and VS0-CF-F15-31
    # must prove the new AUTH_REQUIRED local case without altering the inherited F01 case.
    audited_denial_ids = ("VS0-CF-F15-05", "VS0-CF-F15-06", "VS0-CF-F15-21", "VS0-CF-F15-22", "VS0-CF-F15-23", "VS0-CF-X03")
    for cid in audited_denial_ids:
        entry = next((c for c in confs if c.get("id") == cid), None)
        if not entry:
            e(f"{cid} conformance entry missing (ADH-2026-051)")
            continue
        if "auditevent" not in str(entry.get("expectedSideEffects", "")).lower():
            e(f"{cid} expectedSideEffects must state exactly one redacted FEATURE-0013 AuditEvent (ADH-2026-051)")
    f15_31 = next((c for c in confs if c.get("id") == "VS0-CF-F15-31"), None)
    if not f15_31:
        e("VS0-CF-F15-31 conformance entry missing (ADH-2026-051)")
    else:
        if f15_31.get("owner") != "FEATURE-0015": e("VS0-CF-F15-31 owner must be FEATURE-0015 (ADH-2026-051)")
        if f15_31.get("expectedError") != "AUTH_REQUIRED": e("VS0-CF-F15-31 expectedError must be AUTH_REQUIRED (ADH-2026-051)")
        if "no mutation" not in str(f15_31.get("expectedSideEffects", "")).lower():
            e("VS0-CF-F15-31 expectedSideEffects must state no mutation (ADH-2026-051)")
    # ADH-2026-052: FEATURE-0015 must have exact local proof cases for CloudPlatform/CloudProvider
    # name-uniqueness rejection (F15-32) and general registry-declared schema-constraint validation
    # (F15-33), and F15-33 must explicitly exclude duplicate name, unassigned ISO code, and
    # malformed/prohibited header/body outcomes already covered by other cases.
    f15_32 = next((c for c in confs if c.get("id") == "VS0-CF-F15-32"), None)
    if not f15_32:
        e("VS0-CF-F15-32 conformance entry missing (ADH-2026-052)")
    else:
        if f15_32.get("owner") != "FEATURE-0015": e("VS0-CF-F15-32 owner must be FEATURE-0015 (ADH-2026-052)")
        if f15_32.get("expectedError") != "ALREADY_EXISTS": e("VS0-CF-F15-32 expectedError must be ALREADY_EXISTS (ADH-2026-052)")
        inputs_32 = str(f15_32.get("inputs", "")).lower()
        for required in ("cloudplatform", "cloudprovider", "metadata.name", "duplicate"):
            if required not in inputs_32:
                e(f"VS0-CF-F15-32 inputs must name '{required}' (ADH-2026-052)")
    f15_33 = next((c for c in confs if c.get("id") == "VS0-CF-F15-33"), None)
    if not f15_33:
        e("VS0-CF-F15-33 conformance entry missing (ADH-2026-052)")
    else:
        if f15_33.get("owner") != "FEATURE-0015": e("VS0-CF-F15-33 owner must be FEATURE-0015 (ADH-2026-052)")
        if f15_33.get("expectedError") != "VALIDATION_FAILED": e("VS0-CF-F15-33 expectedError must be VALIDATION_FAILED (ADH-2026-052)")
        inputs_33 = str(f15_33.get("inputs", "")).lower()
        for required in ("duplicate name", "unassigned", "malformed"):
            if required not in inputs_33:
                e(f"VS0-CF-F15-33 inputs must explicitly exclude '{required}' outcomes covered by other cases (ADH-2026-052)")
    for x in found_cf-exp_cf: e(f"Unexpected conformance: {x}")
    retired_confs = reg.get("retiredConformance", [])
    ids_retired_cf = [x.get("id") for x in retired_confs]
    if ids_retired_cf != RETIRED_CONFORMANCE_IDS: e(f"retiredConformance must contain exactly {RETIRED_CONFORMANCE_IDS} in order")
    for tomb in retired_confs:
        if not tomb.get("removalReason") or not tomb.get("retiredBy"): e(f"{tomb.get('id','?')}: retirement traceability incomplete")
    # Traceability
    if TRACE_PATH.exists():
        txt=TRACE_PATH.read_text()
        all_ids=(exp_s+exp_retired+exp_w+exp_sm+exp_f+[f"VS0-CF-{c}" for c in CF_IDS + F15_CF_IDS + F16_CF_IDS])
        for a in all_ids:
            if a not in txt: e(f"Traceability missing: {a}")
        for p in ("VS-000-contract-registry.yaml","VS-000-contract-specification.md","VS-000_CONTRACT_TRACEABILITY_MATRIX.md"):
            if p not in txt: e(f"Traceability missing path: {p}")
    else: e("Traceability matrix not found")
    # Steering
    if STEER_PATH.exists():
        st=STEER_PATH.read_text()
        for p in ("VS-000-contract-registry.yaml","VS-000-contract-specification.md","VS-000_CONTRACT_TRACEABILITY_MATRIX.md"):
            if p not in st: e(f"Steering missing path: {p}")
        if "ARCHITECTURE_DECISION_REQUIRED" not in st: e("Steering missing ARCHITECTURE_DECISION_REQUIRED")
        if "one feature" not in st.lower() or "one stage" not in st.lower(): e("Steering missing one-feature/one-stage")
        if "429" not in st: e("Steering missing no-429 rule")
        for term in PROHIBITED:
            if term not in st: e(f"Steering missing prohibited: {term}")
    else: e("Steering not found")
    # Prohibited in active registry (retired/tombstone sections are historical record, not active)
    active_reg = {k: v for k, v in reg.items() if not k.startswith("retired") and k != "prohibitedTerms"}
    active = yaml.dump(active_reg, default_flow_style=False)
    for term in PROHIBITED:
        if term.lower() in active.lower(): e(f"Registry active data has prohibited: '{term}'")
    # Spec checks
    if SPEC_PATH.exists():
        sp=SPEC_PATH.read_text()
        for inv in INVENTED:
            for ln in sp.split("\n"):
                cells=[c.strip() for c in ln.split("|")]
                if inv in cells: e(f"Spec has invented standalone code: {inv}"); break
        for c in F12_CODES:
            if c not in sp: e(f"Spec missing F12 code: {c}")
    else: e("Spec not found")
    # Charter
    if CHARTER_PATH.exists():
        ch=CHARTER_PATH.read_text()
        if "VS-000-contract-specification.md" not in ch:
            e("Charter missing exact contract specification pointer")
        for i in range(1,21):
            fid=f"VS0-F{i:02d}"
            cnt=len(re.findall(rf"^\|\s*{re.escape(fid)}\s*\|",ch,re.MULTILINE))
            if cnt==0: e(f"Charter missing {fid}")
            elif cnt>1: e(f"Charter has {fid} {cnt}x (need 1)")
        cloud_row=next((ln for ln in ch.splitlines() if re.match(r"^\|\s*Cloud model\s*\|",ln)), "")
        expected_cloud_contracts="CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, InfrastructureStack"
        if not cloud_row: e("Charter missing Cloud model ownership row")
        elif f"| {expected_cloud_contracts} | FEATURE-0015 |" not in cloud_row: e("Charter Cloud model row must assign only the seven FEATURE-0015 resources through InfrastructureStack")
        if "ExecutionTarget" in cloud_row: e("Charter must not assign ExecutionTarget to FEATURE-0015 (ADH-2026-064)")
        integration_row=next((ln for ln in ch.splitlines() if re.match(r"^\|\s*Integration\s*\|",ln)), "")
        for required in ("ExecutionTarget", "normalized target facts", "qualification", "synthetic observer boundary"):
            if required not in integration_row or not integration_row.rstrip().endswith("| FEATURE-0016 |"):
                e(f"Charter Integration row must assign {required} to FEATURE-0016 (ADH-2026-064)")
    else: e("Charter not found")
    # Feature ownership is checked against the authoritative Phase 2 sequence and
    # the executable feature-control manifest, preventing prose-only reassignment.
    schema_by_kind={kind(s):s for s in schemas}
    f15_actual={k for k,s in schema_by_kind.items() if s.get("owner")=="FEATURE-0015"}
    if f15_actual != F15_OWNED: e(f"FEATURE-0015 registry ownership mismatch: {sorted(f15_actual ^ F15_OWNED)}")
    if schema_by_kind.get("ServiceRegion",{}).get("owner") != "FEATURE-0022": e("ServiceRegion must be owned by FEATURE-0022")
    if "SovrunnInstallation" in schema_by_kind: e("SovrunnInstallation must not be an active Slice 0 schema")
    participation=schema_by_kind.get("CloudProviderParticipation", {})
    pso=participation.get("fieldOwnership",{}).get("spec.providerSelectionModes",{})
    if pso != {"introducedBy":"FEATURE-0021", "activatedBy":"FEATURE-0021"}: e("providerSelectionModes must be introduced and activated by FEATURE-0021")
    if any(str(x).startswith("spec.providerSelectionModes:") for x in participation.get("required",[])): e("providerSelectionModes must remain optional before FEATURE-0021")
    phlr=participation.get("fieldOwnership",{}).get("spec.permittedHostingLocationRefs",{})
    if phlr != {"introducedBy":"FEATURE-0021", "activatedBy":"FEATURE-0021"}: e("permittedHostingLocationRefs must be introduced and activated by FEATURE-0021 (ADH-2026-047 decision 4)")
    if any(str(x).startswith("spec.permittedHostingLocationRefs:") for x in participation.get("required",[])): e("permittedHostingLocationRefs must remain optional before FEATURE-0021")
    execution=schema_by_kind.get("ExecutionTarget",{}).get("fieldOwnership",{})
    for field in ("spec.cloudProviderParticipationRef","spec.infrastructureStackRef","spec.targetClass","status.lifecycle","status.qualification","status.maintenanceEpoch","status.observedGeneration","status.factSetRef","status.qualificationResultRef"):
        if execution.get(field) != {"introducedBy":"FEATURE-0016","activatedBy":"FEATURE-0016"}: e(f"ExecutionTarget.{field} must be FEATURE-0016-owned")
    if "status.availability" in schema_by_kind.get("ExecutionTarget",{}).get("fieldOwnership",{}): e("ExecutionTarget must not have an active persisted status.availability field (ADH-2026-058)")
    if schema_by_kind.get("ExecutionTarget",{}).get("owner") != "FEATURE-0016": e("ExecutionTarget must be owned by FEATURE-0016 in its entirety (DEC-0059/ADH-2026-045)")
    # ADH-2026-060: F16-75 (inactive-marker epoch-stale) and F16-89 (active-marker
    # Maintenance-wins) must be mutually exclusive and never reversed.
    f16_by_id={str(c.get("id")):c for c in reg.get("conformance",[]) if str(c.get("id","")).startswith("VS0-CF-F16-")}
    f16_75=f16_by_id.get("VS0-CF-F16-75"); f16_89=f16_by_id.get("VS0-CF-F16-89")
    if not f16_75: e("VS0-CF-F16-75 not found in registry (ADH-2026-060)")
    if not f16_89: e("VS0-CF-F16-89 not found in registry (ADH-2026-060)")
    if f16_75 and f16_89:
        f75_in=str(f16_75.get("inputs","")).lower(); f89_in=str(f16_89.get("inputs","")).lower()
        if "maintenance entry wins" in f75_in: e("VS0-CF-F16-75 must not say Maintenance entry wins; reversed predicate (ADH-2026-060)")
        if "no active current-maintenance marker" not in f75_in: e("VS0-CF-F16-75 input must state no active current-Maintenance marker exists (ADH-2026-060)")
        if "maintenance entry wins" not in f89_in: e("VS0-CF-F16-89 input must state Maintenance entry wins (ADH-2026-060)")
        if f16_75.get("expectedError")!="STALE_RESOURCE_VERSION" or f16_75.get("expectedViolation")!="VS0_TARGET_EPOCH_STALE": e("VS0-CF-F16-75 must be 412 STALE_RESOURCE_VERSION/VS0_TARGET_EPOCH_STALE (ADH-2026-060)")
        if f16_89.get("expectedError")!="CONFLICT" or f16_89.get("expectedViolation")!="VS0_TARGET_MAINTENANCE": e("VS0-CF-F16-89 must be 409 CONFLICT/VS0_TARGET_MAINTENANCE (ADH-2026-060)")
    # ADH-2026-061: F16-42/43/89 registry evidence (link-clearing and abort
    # scope) must remain intact and unaltered by this correction; this
    # checker is not the source of the fixed contradiction, only a guard that
    # the reconciled feature authority did not drift from already-correct
    # registry semantics.
    f16_42=f16_by_id.get("VS0-CF-F16-42"); f16_43=f16_by_id.get("VS0-CF-F16-43")
    if not f16_42: e("VS0-CF-F16-42 not found in registry (ADH-2026-061)")
    if not f16_43: e("VS0-CF-F16-43 not found in registry (ADH-2026-061)")
    if f16_42 and "current factset/result links clear" not in str(f16_42.get("expectedSideEffects","")).lower():
        e("VS0-CF-F16-42 must state current FactSet/Result links clear (ADH-2026-061)")
    if f16_43 and "current factset/result links clear" not in str(f16_43.get("expectedSideEffects","")).lower():
        e("VS0-CF-F16-43 must state current FactSet/Result links clear (ADH-2026-061)")
    if f16_89 and "no qualification result, completion, or qualification auditevent" not in str(f16_89.get("expectedSideEffects","")).lower():
        e("VS0-CF-F16-89 must keep its no-result/no-completion/no-AuditEvent outcome unchanged (ADH-2026-061)")
    # ADH-2026-063: the four new collection-create phase-one/strict-classification
    # rows must exist and must not deviate from their exact approved outcomes.
    f16_123=f16_by_id.get("VS0-CF-F16-123"); f16_124=f16_by_id.get("VS0-CF-F16-124")
    f16_125=f16_by_id.get("VS0-CF-F16-125"); f16_126=f16_by_id.get("VS0-CF-F16-126")
    if not f16_123: e("VS0-CF-F16-123 not found in registry (ADH-2026-063)")
    elif f16_123.get("expectedError")!="AUTHORIZATION_DENIED": e("VS0-CF-F16-123 must be AUTHORIZATION_DENIED (ADH-2026-063)")
    elif "never classifies, canonicalizes, digests, or reserves" not in str(f16_123.get("expectedSideEffects","")).lower(): e("VS0-CF-F16-123 must state phase one never classifies/canonicalizes/digests/reserves the body (ADH-2026-063)")
    if not f16_124: e("VS0-CF-F16-124 not found in registry (ADH-2026-063)")
    elif f16_124.get("expectedError")!="RESOURCE_NOT_FOUND" or f16_124.get("expectedViolation")!="VS0_AUTHORIZATION_SAFE_DENIAL": e("VS0-CF-F16-124 must be safe 404 RESOURCE_NOT_FOUND + VS0_AUTHORIZATION_SAFE_DENIAL (ADH-2026-063)")
    elif "never classifies, canonicalizes, digests, or reserves" not in str(f16_124.get("expectedSideEffects","")).lower(): e("VS0-CF-F16-124 must state phase one never classifies/canonicalizes/digests/reserves the body (ADH-2026-063)")
    if not f16_126: e("VS0-CF-F16-126 not found in registry (ADH-2026-063)")
    elif f16_126.get("expectedError")!="REQUEST_TOO_LARGE": e("VS0-CF-F16-126 must be REQUEST_TOO_LARGE (ADH-2026-063)")
    elif "before phase-one extraction, authorization, or safe access" not in str(f16_126.get("expectedSideEffects","")).lower(): e("VS0-CF-F16-126 must state REQUEST_TOO_LARGE before phase-one extraction, authorization, or safe access (ADH-2026-063)")
    # ADH-2026-065: F16-125 is malformed-JSON only; F16-127 is duplicate
    # top-level member only; F16-128 is the fail-closed missing/unextractable
    # phase-one reference denial. Each is an exact non-family case.
    f16_127=f16_by_id.get("VS0-CF-F16-127"); f16_128=f16_by_id.get("VS0-CF-F16-128")
    if not f16_125: e("VS0-CF-F16-125 not found in registry (ADH-2026-065)")
    elif f16_125.get("expectedError")!="MALFORMED_REQUEST": e("VS0-CF-F16-125 must be MALFORMED_REQUEST (exact non-family malformed-JSON case) (ADH-2026-065)")
    else:
        f125=str(f16_125.get("expectedSideEffects","")).lower()
        if "duplicate" in f125: e("VS0-CF-F16-125 must not name a duplicate case; it is the exact malformed-JSON classification case (ADH-2026-065)")
        elif "malformed" not in f125 or "malformed_request" not in f125 or "strict single classification" not in f125: e("VS0-CF-F16-125 must state strict single classification returning MALFORMED_REQUEST for malformed JSON (ADH-2026-065)")
    if not f16_127: e("VS0-CF-F16-127 not found in registry (ADH-2026-065)")
    elif f16_127.get("expectedError")!="DUPLICATE_FIELD": e("VS0-CF-F16-127 must be DUPLICATE_FIELD (exact non-family duplicate-top-level-member case) (ADH-2026-065)")
    else:
        f127=str(f16_127.get("expectedSideEffects","")).lower()
        if "malformed" in f127: e("VS0-CF-F16-127 must not name a malformed case; it is the exact duplicate-top-level-member classification case (ADH-2026-065)")
        elif "duplicate" not in f127 or "duplicate_field" not in f127 or "strict single classification" not in f127: e("VS0-CF-F16-127 must state strict single classification returning DUPLICATE_FIELD for a duplicate top-level member (ADH-2026-065)")
    if not f16_128: e("VS0-CF-F16-128 not found in registry (ADH-2026-065)")
    elif f16_128.get("expectedError")!="AUTHORIZATION_DENIED": e("VS0-CF-F16-128 must be the existing audited AUTHORIZATION_DENIED phase-one extraction denial (ADH-2026-065)")
    else:
        f128=(str(f16_128.get("inputs",""))+" "+str(f16_128.get("expectedSideEffects",""))).lower()
        if not("exactly one syntactically usable uid" in f128 and "both" in f128 and "spec.cloudproviderparticipationref.uid" in f128 and "spec.infrastructurestackref.uid" in f128):
            e("VS0-CF-F16-128 must state phase one cannot extract exactly one syntactically usable UID at both required reference paths (ADH-2026-065)")
        elif not("audited" in f128 and "403" in f128 and "no body-classification detail" in f128 and "no backing-resource existence" in f128):
            e("VS0-CF-F16-128 must state the existing audited AUTHORIZATION_DENIED/403 with no body/backing disclosure (ADH-2026-065)")
        elif not all(tok in f128 for tok in ("strict classification","canonicalization","digest","reservation","observer","mutation","publication","completion")):
            e("VS0-CF-F16-128 must state no strict classification/canonicalization/digest/reservation/observer/mutation/publication/completion (ADH-2026-065)")
    # ADH-2026-066 introduces no registry semantics. It only records an
    # internal compatibility bridge, whose traceability must not be allowed to
    # drift into a FEATURE-0015 ownership or a new conformance requirement.
    trace = TRACE_PATH.read_text() if TRACE_PATH.exists() else ""
    if "ADH-2026-066" not in trace:
        e("TRACEABILITY: ADH-2026-066 backing-access bridge must be recorded in the VS-000 traceability matrix")
    if any("ParticipationGeneration" in str(case) for case in f16_by_id.values()):
        e("ADH-2026-066: ParticipationGeneration must not become a registered F0016 conformance fence")
    region=schema_by_kind.get("ServiceRegion",{}).get("fieldOwnership",{})
    for field in ("spec.displayName","spec.hostingLocationRefs"):
        if region.get(field) != {"introducedBy":"FEATURE-0022","activatedBy":"FEATURE-0022"}: e(f"ServiceRegion.{field} must be FEATURE-0022-owned")
    for field in ("spec.executionTargetRefs","status.availability"):
        if region.get(field) != {"introducedBy":"FEATURE-0022","activatedBy":"FEATURE-0022","dependsOn":"FEATURE-0016"}: e(f"ServiceRegion.{field} ownership/dependency mismatch")
    for k in ("CanonicalMigrationPlan","CanonicalMigrationRecord"):
        if k in schema_by_kind: e(f"{k} must not be an active Slice 0 schema (DEC-0059)")
    if SEQUENCE_PATH.exists():
        seq=SEQUENCE_PATH.read_text()
        f15_line=next((ln for ln in seq.splitlines() if re.match(r"^\|\s*5\s*\|\s*FEATURE-0015\b",ln)), "")
        f22_line=next((ln for ln in seq.splitlines() if re.match(r"^\|\s*12\s*\|\s*FEATURE-0022\b",ln)), "")
        for resource in F15_OWNED:
            display="ExecutionTarget identity" if resource=="ExecutionTarget" else resource
            if display not in f15_line: e(f"PHASE2_FEATURE_SEQUENCE FEATURE-0015 missing {display}")
        for leaked in ("ServiceRegion","SovrunnInstallation","qualification","availability"):
            if leaked in f15_line: e(f"PHASE2_FEATURE_SEQUENCE leaks {leaked} into FEATURE-0015")
        if "ServiceRegion" not in f22_line: e("PHASE2_FEATURE_SEQUENCE FEATURE-0022 missing ServiceRegion")
    else: e("PHASE2_FEATURE_SEQUENCE not found")
    if F15_CONTROL_PATH.exists():
        try: control=json.loads(F15_CONTROL_PATH.read_text())
        except Exception as ex: control={}; e(f"FEATURE-0015 control manifest invalid JSON: {ex}")
        owned=set(control.get("ownership",{}).get("owned_resources",[]))
        exp_control=F15_CONTROL_EXPECTED_OWNED
        if owned != exp_control: e(f"FEATURE-0015 control owned_resources mismatch: {sorted(owned ^ exp_control)}")
        excluded={x.get("id") for x in control.get("ownership",{}).get("excluded_features",[])}
        if excluded != {f"FEATURE-{i:04d}" for i in range(16,27)}: e("FEATURE-0015 control must exclude exactly FEATURE-0016..0026")
        for path in control.get("feature",{}).get("handoffs",[]):
            if not (REPO/path).exists(): e(f"FEATURE-0015 control handoff missing: {path}")
        fc=subprocess.run([sys.executable,"scripts/feature-control.py","validate","--feature","FEATURE-0015","--manifest",str(F15_CONTROL_PATH.relative_to(REPO))],capture_output=True,text=True,cwd=REPO,timeout=30)
        if fc.returncode: e(f"feature-control validation failed: {(fc.stdout+fc.stderr).strip()}")
    else: e("FEATURE-0015 control manifest not found")
    for p,label in ((F15_ARCH_PATH,"architecture"),(F15_FEATURE_PATH,"feature")):
        if not p.exists(): e(f"FEATURE-0015 {label} file not found")
    if KIRO_AGENT_PATH.exists():
        agent=KIRO_AGENT_PATH.read_text()
        for required in (str(F15_ARCH_PATH.relative_to(REPO)), str(F15_FEATURE_PATH.relative_to(REPO)), str(F15_CONTROL_PATH.relative_to(REPO))):
            if f"file://{required}" not in agent: e(f"Kiro specification agent missing FEATURE-0015 resource: {required}")
    else: e("Kiro specification agent not found")
    if F15_FEATURE_PATH.exists():
        feature_text=F15_FEATURE_PATH.read_text()
        if "ServiceRegion" in feature_text: e("FEATURE-0015 feature requirements must not mention ServiceRegion")
        for leak in ("status fields initialized", "stores 0", "field stored"):
            if leak in feature_text.lower(): e(f"FEATURE-0015 feature requirements contain activation leakage: {leak}")
    # No .kiro/specs for FEATURE-0015
    sd=REPO/".kiro"/"specs"
    if sd.exists():
        for p in sd.rglob("*"):
            n=p.name.lower().replace("-","").replace("_","")
            if "feature0015" in n or "0015" in n: e(f"FEATURE-0015 spec exists: {p.relative_to(REPO)}")
    # No Go changes
    try:
        r1=subprocess.run(["git","diff","--name-only","HEAD"],capture_output=True,text=True,cwd=REPO,timeout=10)
        r2=subprocess.run(["git","diff","--cached","--name-only"],capture_output=True,text=True,cwd=REPO,timeout=10)
        for f in (r1.stdout+r2.stdout).strip().split("\n"):
            if f.strip().endswith(".go"): e(f"Go file changed: {f.strip()}")
    except Exception: pass
    # architectureDecisionRequired
    adr=reg.get("architectureDecisionRequired")
    if adr is None: e("architectureDecisionRequired missing (must be list)")
    elif not isinstance(adr,list): e(f"architectureDecisionRequired not a list")
    elif adr:
        for gap in adr: e(f"Unresolved gap: {gap}")

run()
if errs:
    print(f"FAIL: {len(errs)} error(s)")
else:
    reg=yaml.safe_load(REG_PATH.read_text())
    print(
        "PASS: 0 error(s); "
        f"active-schemas={len(reg.get('schemas', []))}, "
        f"retired-schemas={len(reg.get('retiredSchemas', []))}, "
        f"writers={len(reg.get('writers', []))}, "
        f"state-machines={len(reg.get('stateMachines', []))}, "
        f"problem-codes={len(reg.get('problemCodes', []))}, "
        f"failure-mappings={len(reg.get('failureMappings', []))}, "
        f"conformance={len(reg.get('conformance', []))}"
    )
for x in errs: print(f"  ✗ {x}")
sys.exit(1 if errs else 0)
