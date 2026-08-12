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
F15_CF_IDS = [f"F15-{i:02d}" for i in range(1,26)]
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
    exp_cf={f"VS0-CF-{c}" for c in CF_IDS + F15_CF_IDS}; found_cf=set(); seen_cf=set()
    for c in confs:
        cid=c.get("id","?")
        if cid in seen_cf: e(f"Dup conformance: {cid}")
        seen_cf.add(cid); found_cf.add(cid)
        for f in ("owner","inputs","expectedState","expectedError","expectedSideEffects","gate"):
            if f not in c: e(f"{cid}: missing {f}")
    for x in exp_cf-found_cf: e(f"Missing conformance: {x}")
    for x in found_cf-exp_cf: e(f"Unexpected conformance: {x}")
    retired_confs = reg.get("retiredConformance", [])
    ids_retired_cf = [x.get("id") for x in retired_confs]
    if ids_retired_cf != RETIRED_CONFORMANCE_IDS: e(f"retiredConformance must contain exactly {RETIRED_CONFORMANCE_IDS} in order")
    for tomb in retired_confs:
        if not tomb.get("removalReason") or not tomb.get("retiredBy"): e(f"{tomb.get('id','?')}: retirement traceability incomplete")
    # Traceability
    if TRACE_PATH.exists():
        txt=TRACE_PATH.read_text()
        all_ids=(exp_s+exp_retired+exp_w+exp_sm+exp_f+[f"VS0-CF-{c}" for c in CF_IDS + F15_CF_IDS])
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
    execution=schema_by_kind.get("ExecutionTarget",{}).get("fieldOwnership",{})
    for field in ("spec.infrastructureStackRef","spec.participationRef","spec.targetClass","status.qualification","status.availability","status.maintenanceEpoch","status.factSetRef","status.observedGeneration","status.conditions"):
        if execution.get(field) != {"introducedBy":"FEATURE-0016","activatedBy":"FEATURE-0016"}: e(f"ExecutionTarget.{field} must be FEATURE-0016-owned")
    if schema_by_kind.get("ExecutionTarget",{}).get("owner") != "FEATURE-0016": e("ExecutionTarget must be owned by FEATURE-0016 in its entirety (DEC-0059/ADH-2026-045)")
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
