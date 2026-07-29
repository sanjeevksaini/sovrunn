# FEATURE-0013 Matrix E Worksheet — Residual-risk staging

Feature: FEATURE-0013 — Decision Record and AuditEvent Standard
Requirements: F13-GUARD-001, F13-GUARD-002 (governance staging only)
Architecture: §18.3, §18.5, §18.6, §27.9(12)
Task: T-032

Status: **STAGED — PENDING_HUMAN_REVIEW**

This worksheet stages architecture-owned Matrix E values and blank
human-only disposition fields. Automation and coding agents must **not**:

- calculate, infer, downgrade, upgrade, or replace residual risk levels;
- mark human approval complete;
- advance stage status;
- create review verdicts;
- accept residual risk.

Architecture-stage treatment approval (`APPROVE_TREATMENT`, 2026-07-27)
authorizes these risks as requirements inputs. It does **not** accept final
residual risk. Final residual decisions require human completion of the
per-risk record in §18.5.

Source register:
`docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
§18.3–§18.6.

Explicitly open High target-residual risks (architecture-owned; not
human-accepted): **F13-R04**, **F13-R08**, **F13-R13**.

FEATURE-0012 identifiers `F12-R01`–`F12-R16` are out of scope. This worksheet
uses only `F13-R01`–`F13-R31`.

---

## Architecture-owned staging table

Values below are copied from architecture §18.3. Target residual is the
architecture target after stated controls, **not** human acceptance.
`Final residual status` remains `PENDING_HUMAN_REVIEW` for every row.

| ID | Category | Risk scenario | ADs | Inherent L×I | FEATURE-0013 controls | Verification / detection | Evidence pointers | Target residual | Treatment | Residual / later owner | Reassessment trigger | Final residual status |
|---|---|---|---|---|---|---|---|---|---|---|---|---|
| F13-R01 | Governance/compatibility | Later features create overlapping `DecisionProfile` semantics | AD-029, AD-030 | 4×4 High | Governed registration, ownership, reuse review, compatibility | Duplicate-profile and new-family fixtures | arch §10; T-009, T-030 | 2×3 Medium | MITIGATE | Architecture owner | New profile family or stable promotion | PENDING_HUMAN_REVIEW |
| F13-R02 | Correctness | Typed results become an ungoverned property bag or override the envelope | AD-002, AD-029 | 4×4 High | Versioned result schemas; semantic-override prohibition | Schema and negative conformance | arch §7.3; T-013, T-030 | 2×4 Medium | MITIGATE | FEATURE-0013 owner | New result extension mechanism | PENDING_HUMAN_REVIEW |
| F13-R03 | Security/correctness | Advisory, recommendation, simulation, or AI output is consumed as authority | AD-020, AD-031 | 4×5 Critical | Explicit authority; projection marking; enforcement validation | Authority-confusion negative tests and audit | arch §12; T-025 | 1×5 Medium | MITIGATE | Security architecture; enforcement owner | New AI/execution consumer | PENDING_HUMAN_REVIEW |
| F13-R04 | Correctness | Stale `EffectivePolicyContext` produces an invalid decision | AD-007, AD-038 | 4×5 Critical | Versioned identity, optional opaque integrity carrier, effective time, freshness, revocation, fail-closed rules | Cache expiry/revocation/staleness tests without digest computation | arch §12; T-005, T-031 | 2×5 High | MITIGATE | Effective-context owner | Production policy resolver or new inheritance rule | PENDING_HUMAN_REVIEW |
| F13-R05 | Correctness | Weak retry identity or cache identity reuses the wrong decision | AD-013, AD-036 | 4×5 Critical | Separate retry key and complete semantic identity descriptor; later optional digest profile | Policy/profile/strategy/time mutation tests | arch §7.1; T-005 | 1×5 Medium | MITIGATE | Decision persistence owner | Semantic-identity fields or canonical digest profile change | PENDING_HUMAN_REVIEW |
| F13-R06 | Audit/compliance | A decision is durable without its audit obligation | AD-014, AD-035 | 3×5 High | Same-boundary atomic acceptance and preallocated audit identity | Partial-failure and reconciliation fixtures | arch §11; T-031 | 1×5 Medium | MITIGATE | Audit persistence owner | Production persistence selection | PENDING_HUMAN_REVIEW |
| F13-R07 | Availability/latency | Remote audit processing blocks synchronous decisions | AD-008, AD-035 | 4×4 High | Local obligation acceptance; asynchronous downstream processing | Remote-audit outage and backlog tests | arch §11; T-031 | 1×4 Low | AVOID | Audit implementation owner | External audit/SIEM integration | PENDING_HUMAN_REVIEW |
| F13-R08 | Cost/availability | Composite graphs cause fan-out, latency, retry, or cost explosion | AD-006, AD-040 | 4×5 Critical | Bounds on graph, calls, bytes, authorities, concurrency, retries, time, jurisdiction, and cost | Limit, cancellation, load, and adversarial fixtures | arch §10, §9; T-008, T-019 | 2×5 High | MITIGATE | Decision runtime owner | Limit increase or new remote evaluator | PENDING_HUMAN_REVIEW |
| F13-R09 | Correctness/AI | Nondeterministic evaluator output cannot be replayed or explained | AD-005, AD-037 | 4×4 High | Captured normalized output, provenance, input identity, and opaque integrity carrier | Replay without evaluator access or digest computation | arch §8; T-011 | 2×3 Medium | MITIGATE | Evaluator/AI owner | New probabilistic evaluator | PENDING_HUMAN_REVIEW |
| F13-R10 | Security/AI | AI receives unauthorized evidence or sensitive canonical data | AD-019, AD-020 | 4×5 Critical | Authorized `DecisionContext`, minimization, classification, projection policies | Cross-scope/redaction and prompt-leak fixtures | arch §12; T-025 | 1×5 Medium | MITIGATE | AI-context and security owners | New model, projection, or jurisdiction | PENDING_HUMAN_REVIEW |
| F13-R11 | Security | Enforcement silently ignores an unknown mandatory obligation | AD-032, AD-043 | 3×5 High | Unknown/unsupported obligation makes allow unenforceable | Obligation capability negative tests | arch §9; T-023 | 1×5 Medium | AVOID | Enforcement owner | New obligation vocabulary | PENDING_HUMAN_REVIEW |
| F13-R12 | Latency/availability | Central authorization or decision service becomes a universal bottleneck | AD-008, AD-039 | 4×5 Critical | Local verified artifacts, policy snapshots, caches, bounded freshness | Central-service outage and horizontal-load tests | arch §12; T-031 | 2×4 Medium | MITIGATE | Authorization implementation owner | Production enforcement integration | PENDING_HUMAN_REVIEW |
| F13-R13 | Security | Revoked or stale cached authorization remains enforceable | AD-039, AD-043 | 3×5 High | Audience/scope/expiry/revocation/freshness validation; fail closed | Cache expiry, revocation, scope-mismatch tests | arch §12; T-031 | 2×5 High | MITIGATE | Authorization/security owner | Revocation architecture selection | PENDING_HUMAN_REVIEW |
| F13-R14 | Security/availability | Fail-open behavior exposes protected resources during dependency failure | AD-005, AD-043 | 3×5 High | Fail-closed high-risk profiles; bounded approved `SecurityExceptionRef` structural evidence | Failure injection, exception-reference, scope, expiry, and owner-mismatch checks | arch §12; T-004, T-019 | 1×5 Medium | AVOID | Security architecture owner | Any fail-open request | PENDING_HUMAN_REVIEW |
| F13-R15 | Correctness | Correction, supersession, or revocation yields ambiguous effective state | AD-004, AD-034 | 3×4 High | Immutable linked effects; bounded cycle-free deterministic projection | Chain, cycle, fork, and projection fixtures | arch §9.3; T-024 | 1×4 Low | MITIGATE | Read-model owner | New relationship/effective-state semantics | PENDING_HUMAN_REVIEW |
| F13-R16 | Privacy/security | Decision, rationale, evidence, error, or projection leaks restricted data | AD-019, AD-024 | 4×5 Critical | References/digests, minimization, classification, redaction, purpose projections | Secret/PII/cross-scope negative tests | arch §12; T-025, T-030 | 1×5 Medium | MITIGATE | Security/privacy owner | New evidence or projection type | PENDING_HUMAN_REVIEW |
| F13-R17 | Compliance/integrity | Erasure or retention action breaks evidence integrity or legal hold | AD-004, AD-022 | 3×5 High | Structural tombstone and retention/legal-hold contracts now; cryptographic erasure later | Structural hold, erasure-state, opaque-carrier, and tombstone tests | arch §12; T-004 | 2×4 Medium | MITIGATE | Evidence/retention owner | Regulated deployment or retention change | PENDING_HUMAN_REVIEW |
| F13-R18 | Latency/availability | Remote HSM, time, signing, or notarization becomes a hot-path dependency | AD-022, AD-041 | 4×4 High | No cryptographic service in FEATURE-0013; later tiered local/batch/asynchronous assurance | No-network/no-crypto gate now; outage and batch-proof tests in owning later feature | arch §15; T-033 | 1×4 Low | AVOID | Security infrastructure owner | High-assurance profile adoption | PENDING_HUMAN_REVIEW |
| F13-R19 | Sovereignty/security | Disconnected import accepts forged, replayed, or untrusted records | AD-018, AD-042 | 3×5 High | Static origin-aware manifests, algorithm-agile signature carriers, replay/duplicate markers, ordering, opaque trust states, and fail-closed quarantine semantics now; cryptographic signing and verification later under ADR-F13-002 | Structural replay/duplicate/ordering/unknown-or-revoked-trust-state fixtures now; forgery and signature-verification evidence in the later owning feature | arch §10; T-009, T-030 | 1×5 Medium | MITIGATE | Synchronization/security owner | First disconnected deployment | PENDING_HUMAN_REVIEW |
| F13-R20 | Correctness | Disconnected origins create conflicting authoritative decisions | AD-018, AD-042 | 3×5 High | Immutable origin identity; no last-writer-wins; governed reconciliation decision | Concurrent-origin conflict fixtures | arch §10; T-030 | 2×4 Medium | MITIGATE | Multi-site owner | Active/active or disconnected write enablement | PENDING_HUMAN_REVIEW |
| F13-R21 | Portability | Provider/runtime-native identifiers or schemas leak into the core | AD-011, AD-023 | 4×4 High | Adapter-only native data; normalized capabilities; typed references | Import/schema lint and portability fixtures | arch §17; T-033 | 1×4 Low | AVOID | Adapter architecture owner | First provider/runtime adapter | PENDING_HUMAN_REVIEW |
| F13-R22 | Scope/governance | Production registry, workflow, persistence, provider, or AI runtime leaks into Phase 2 | AD-011, AD-028 | 4×4 High | Contract-only scope, non-goals, feature gates | Changed-file, dependency, and no-side-effect gates | arch §5, §15; T-033 | 1×4 Low | AVOID | FEATURE-0013 owner | Any runtime dependency proposal | PENDING_HUMAN_REVIEW |
| F13-R23 | Cost | Durable decisions, audit, evidence, or projections create uncontrolled storage growth | AD-024, AD-033 | 4×4 High | Risk-based recording, bounded payloads, hot summaries, referenced cold evidence, retention profiles | Volume/limit/retention tests and metrics | arch §10; T-002 | 2×4 Medium | MITIGATE | Audit/evidence owner | Retention increase or production volume | PENDING_HUMAN_REVIEW |
| F13-R24 | Supply chain/security | Profile, schema, strategy, evaluator, or trust bundle is compromised | AD-025, AD-038 | 3×5 High | Version-pinned bundles, opaque trust-root identity, structural revocation/quarantine/provenance now; cryptographic verification later | Structural tamper-state/revocation fixtures now; signature verification in owning later feature | arch §10; T-009 | 1×5 Medium | MITIGATE | Supply-chain security owner | Signing-root or distribution change | PENDING_HUMAN_REVIEW |
| F13-R25 | Correctness/security | Clock uncertainty invalidates ordering, expiry, freshness, or replay protection | AD-018, AD-022 | 3×5 High | Trusted-time provenance, validity epoch, profile outage behavior | Skew, rollback, expiry-boundary tests | arch §7.1; T-005 | 2×4 Medium | MITIGATE | Runtime/security owner | Disconnected or multi-site deployment | PENDING_HUMAN_REVIEW |
| F13-R26 | Security/privacy | Cross-scope references disclose or influence another tenant or governance domain | AD-015, AD-019 | 3×5 High | Scope authorization, no-existence disclosure, typed references, projections | Cross-scope and safe-denial fixtures | arch §9; T-018, T-020 | 1×5 Medium | MITIGATE | API/security owner | New scope/reference kind | PENDING_HUMAN_REVIEW |
| F13-R27 | Compatibility/governance | Later features silently redefine common decision or public error semantics | AD-029, AD-030, AD-045 | 4×5 Critical | Stable envelope, profile conformance, closed public error registry, architecture change control | Baseline/schema/error semantic-diff gate | arch §9.2; T-017, T-015 | 1×5 Medium | AVOID | Architecture owner | Common-envelope, closed-vocabulary, or public-error change | PENDING_HUMAN_REVIEW |
| F13-R28 | Compatibility | Expanded AuditEvent allowed scopes or ServiceInstance subject handling breaks FEATURE-0012 consumers | AD-015 | 3×4 High | Retain the six-value ScopeKind vocabulary, use typed ServiceInstance subject references, preserve Organization fixtures, and require schema-diff review | Six-scope AuditEvent compatibility, Organization regression, and Project-scope/ServiceInstance-subject tests | arch §11; T-014, T-029 | 1×4 Low | MITIGATE | FEATURE-0013/API owner | Stable API promotion | PENDING_HUMAN_REVIEW |
| F13-R29 | Traceability | `DecisionRecord` rename breaks spine, documents, schemas, or consumer identity | AD-001 | 3×4 High | Architecture clarification, coordinated traceability update, aliases/migration if required | Repository-wide identity and baseline checks | arch §7.1; T-002, T-015 | 1×4 Low | MITIGATE | Architecture owner | Rename approval or stable promotion | PENDING_HUMAN_REVIEW |
| F13-R30 | Cost/sovereignty | High-assurance controls make small or local deployments unaffordable | AD-016, AD-017, AD-041 | 3×4 High | Tiered conformance and assurance profiles; no universal remote dependencies | Minimal/regulated/sovereign profile cost fixtures | arch §10; T-009 | 2×3 Medium | MITIGATE | Deployment architecture owner | New mandatory assurance control | PENDING_HUMAN_REVIEW |
| F13-R31 | Governance/cost | Downstream adoption governance becomes duplicative bureaucracy or a premature registry framework | AD-044 | 4×3 High | One normative contract, one short per-feature section, one lightweight gate; trigger-based escalation only | Changed-file/gate review; absence of manifests, registry, index, and separate approval workflow | arch §16; T-032 | 1×3 Low | AVOID | Architecture governance owner | Repeated semantic drift or structured-manifest trigger in section 28.8 | PENDING_HUMAN_REVIEW |

---

## Human-only disposition fields (blank until human review)

Per architecture §18.5, the human reviewer must complete one record per risk.
Architecture treatment approval and a global feature approval do **not** accept
unreviewed residual risks.

All fields below are intentionally blank. Do not fill them by automation.

### Global human decision fields

| Field | Value (human only) |
|---|---|
| Reviewer identity | |
| Review date | |
| Residual-risk decision summary | |
| Decision notes | |

### Per-risk residual acceptance record (template)

Duplicate and complete one block per risk ID (`F13-R01` … `F13-R31`). Leave
blank until a human records disposition.

```yaml
risk_id: F13-Rxx
decision:            # ACCEPT | REJECT | DEFER
residual_likelihood: # 1..5
residual_impact:     # 1..5
residual_level:      # LOW | MEDIUM | HIGH | CRITICAL
residual_rationale:
owner:
later_feature:
corrective_path:
reassessment_trigger:
acceptance_expiry:
reviewer:
date:
evidence_reviewed: []
```

### Per-risk human disposition checklist

| ID | Human decision | Residual L | Residual I | Residual level | Owner (human) | Later feature | Corrective path | Acceptance expiry | Reviewer | Date | Evidence reviewed |
|---|---|---|---|---|---|---|---|---|---|---|---|
| F13-R01 | | | | | | | | | | | |
| F13-R02 | | | | | | | | | | | |
| F13-R03 | | | | | | | | | | | |
| F13-R04 | | | | | | | | | | | |
| F13-R05 | | | | | | | | | | | |
| F13-R06 | | | | | | | | | | | |
| F13-R07 | | | | | | | | | | | |
| F13-R08 | | | | | | | | | | | |
| F13-R09 | | | | | | | | | | | |
| F13-R10 | | | | | | | | | | | |
| F13-R11 | | | | | | | | | | | |
| F13-R12 | | | | | | | | | | | |
| F13-R13 | | | | | | | | | | | |
| F13-R14 | | | | | | | | | | | |
| F13-R15 | | | | | | | | | | | |
| F13-R16 | | | | | | | | | | | |
| F13-R17 | | | | | | | | | | | |
| F13-R18 | | | | | | | | | | | |
| F13-R19 | | | | | | | | | | | |
| F13-R20 | | | | | | | | | | | |
| F13-R21 | | | | | | | | | | | |
| F13-R22 | | | | | | | | | | | |
| F13-R23 | | | | | | | | | | | |
| F13-R24 | | | | | | | | | | | |
| F13-R25 | | | | | | | | | | | |
| F13-R26 | | | | | | | | | | | |
| F13-R27 | | | | | | | | | | | |
| F13-R28 | | | | | | | | | | | |
| F13-R29 | | | | | | | | | | | |
| F13-R30 | | | | | | | | | | | |
| F13-R31 | | | | | | | | | | | |

---

## Risk traceability lifecycle (architecture §18.6)

```text
architecture.md Matrix E source
  -> Kiro requirements preserve IDs and mandate controls
  -> design maps controls to components and verification
  -> tasks implement controls and evidence collection
  -> automation stages immutable evidence pointers   ← this worksheet (T-032)
  -> human records per-risk residual disposition     ← blank above
  -> final feature approval summarizes, but does not replace, those records
```

---

## Explicit non-claims

- This worksheet does **not** accept residual risk for any `F13-Rxx` ID.
- Automated PASS, architecture `APPROVE_TREATMENT`, reuse approval, or
  prior-feature acceptance is **not** residual-risk acceptance.
- Target residual ratings above are architecture-owned staging values only.
- F13-R04, F13-R08, and F13-R13 remain explicitly open at High target residual
  pending later owners and production-entry gates.
- No erasure, key destruction, signing, cryptographic verification, AI runtime,
  production runtime, or persistence mechanism is authorized by this document.
