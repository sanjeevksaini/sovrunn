# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-071
- Date: 2026-08-30
- Source discussion: ChatGPT Project
- Related feature: FEATURE-0018 — Governance, IAM, Approval and Exception Foundation
- Related phase: Phase 2R
- Author: ChatGPT / Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

FEATURE-0018 conformance executability reconciliation

## Summary

This proposed correction makes the already-approved FEATURE-0018 architecture
executable by the mandatory Slice-0 stage contract. It allocates exact
FEATURE-0018-local `VS0-CF-F18-*` conformance identifiers, maps every active
`AC-F18-*` and every `REQ-F18-*` to registered local proof, and closes the
requirements-stage reconciliation defects reported by the independent reviewer.
It changes no FEATURE-0018 resource, field, action, scope, lifecycle, writer,
authorization rule, approval rule, error contract, route, persistence behavior,
external integration, or other runtime semantic. ADH-2026-070 and DEC-0060 remain
the sole semantic authority.

## Classification

- Correction

## Existing approved baseline

The active baseline is `ARCH-2026.08-PHASE2R-CANONICAL`. DEC-0060 and the
approved ADH-2026-070 three-file package close FEATURE-0018 semantics through
F18-RD-01..24. F18-RD-22 requires deterministic, FEATURE-0018-local executable
conformance and prohibits integration evidence from substituting for local
proof.

The mandatory `.kiro/steering/slice0-contract.md` contract additionally requires:

1. every requirement to map to its architecture authority, owner, applicable
   schema/writer/state/error identifiers, and an exact registered `VS0-CF-*`;
2. every acceptance criterion to map to an exact registered `VS0-CF-*`; and
3. registry, traceability, and test contracts to change atomically.

The current approved feature authority supplies exact `REQ-F18-01..24` and
`AC-F18-01..49`, with `AC-F18-38` explicitly excluded, but the Slice-0 registry
contains only four shared FEATURE-0018-owned security cases (`VS0-CF-F01`,
`VS0-CF-F02`, `VS0-CF-X01`, and `VS0-CF-X02`). That is a conformance-executability
gap, not an unclosed runtime behavior decision.

Relevant baseline references:

- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/decisions/DECISION_INDEX.md`
- `docs/decisions/DEC-0060-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md`
- `.kiro/steering/slice0-contract.md`

## Decision or proposed decision

### 1. Authority and non-semantic boundary

ADH-2026-071 authorizes only conformance identifiers, traceability mappings, and
faithful transcription of already-approved semantics. A registration below is a
stable name for existing F18-RD/F18-SCN behavior; it is not a new behavioral
source. If application would require Kiro to choose or change an observable
runtime result, Kiro must stop with `ARCHITECTURE_DECISION_REQUIRED` rather than
infer that result from this handoff.

For the registration ledgers below:

- `expectedError` names an already-approved decision/validation outcome. `Deny`
  is an `AuthorizationResult` or domain decision, not a new Problem code.
- `none` means that the registered successful variant has no error. It does not
  suppress an error required by an approved negative variant.
- `existing safe error` means the exact already-approved FEATURE-0012/F18-RD
  error selected by the cited semantic authority; it expressly authorizes no
  new top-level code, violation code, HTTP status, or precedence rule.
- all state, audit, decision, idempotency, and publication effects are exactly
  those closed by the mapped F18-RD and F18-SCN authority.

### 2. Identifier allocation and invariants

1. Register `VS0-CF-F18-01..VS0-CF-F18-37` and
   `VS0-CF-F18-39..VS0-CF-F18-54` as FEATURE-0018-owned local conformance
   identifiers.
2. `VS0-CF-F18-01..37` and `VS0-CF-F18-39..49` map one-to-one to the active
   `AC-F18-*` having the same suffix.
3. `VS0-CF-F18-38` is a permanent tombstone corresponding to excluded
   `AC-F18-38` / retired duplicate F18-SCN-38. It is never active proof and its
   identifier must never be reused.
4. `VS0-CF-F18-50..54` are requirement-only architecture/contract conformance
   cases. They close exact proof for REQ-F18-01, REQ-F18-18, REQ-F18-19,
   REQ-F18-22, and REQ-F18-23, whose obligations are not completely represented
   by one end-user acceptance scenario.
5. Existing shared cases `VS0-CF-F01`, `VS0-CF-F02`, `VS0-CF-X01`, and
   `VS0-CF-X02` remain unchanged and may be cited as supplementary shared proof.
   They do not replace a local `VS0-CF-F18-*` mapping.
6. `VS0-CF-X03` remains FEATURE-0015-owned inherited safe-denial evidence. It is
   never FEATURE-0018-local acceptance proof and must not be deferred to
   FEATURE-0018 design.
7. The already-existing global `VS0-CF-F18` identifier owned by FEATURE-0026 is
   distinct from the new hyphenated FEATURE-0018-local range
   `VS0-CF-F18-NN`; neither identifier changes ownership or meaning.

### 3. Exact active AC-local conformance registration ledger

Every row below is normative for identifier registration and traceability. Its
scenario meaning is the exact mapped F18-SCN/AC behavior already approved by
ADH-2026-070; abbreviated outcome text must not be used to weaken the full
semantic authority.

| ID | Owner | Exact input fixture | Expected state/result | Expected error/decision | Expected side effects | Gate | Exact AC |
|---|---|---|---|---|---|---|---|
| VS0-CF-F18-01 | FEATURE-0018 | Active employee, published Developer role, active project assignment; same action at unrelated project | requested project `Allow`; unrelated project unchanged | `none`; unrelated project `Deny` | exact protected authorization evidence and audit; no unrelated-scope effect | security | AC-F18-01 |
| VS0-CF-F18-02 | FEATURE-0018 | Same issuer/subject principal with changed email display attribute | principal and existing authority remain bound to issuer/subject | none | attribution retains stable principal identity | identity | AC-F18-02 |
| VS0-CF-F18-03 | FEATURE-0018 | Active assignment whose required Membership/principal becomes suspended | assignment retained but immediately ineffective | `Deny` | retained history and protected denial audit; no access effect | security | AC-F18-03 |
| VS0-CF-F18-04 | FEATURE-0018 | Revoked/expired relationship followed by a later rejoin | new relationship does not reactivate old assignment | `Deny` until newly authorized grant | retained old evidence; no silent reactivation | lifecycle | AC-F18-04 |
| VS0-CF-F18-05 | FEATURE-0018 | Workload principal with responsible Human and exact scoped assignment invokes one action | applicable action may `Allow`; unrelated action/scope denies | `none` or `Deny` by fixture | independent workload attribution and protected audit | security | AC-F18-05 |
| VS0-CF-F18-06 | FEATURE-0018 | System controller submits an owned intent and separately attempts a non-owned write | owned controller transition may publish; non-owned write unchanged | `none` or `AUTHORIZATION_DENIED` | sole writer publishes; direct competing write has no effect | writer | AC-F18-06 |
| VS0-CF-F18-07 | FEATURE-0018 | Consultant has distinct supplier and customer Guest Membership fixtures | only exact active target context is eligible until exact expiry | `none` or existing safe error | sponsorship, scope, expiry, and responsible Human remain attributable | security | AC-F18-07 |
| VS0-CF-F18-08 | FEATURE-0018 | Raw external group claim; then trusted-provisioned AccessGroup, current direct Membership, and active scoped assignment | raw claim grants nothing; canonical relationship may allow | `Deny`; then `none` when all canonical evidence applies | no native/external group authority; canonical provenance audited | security | AC-F18-08 |
| VS0-CF-F18-09 | FEATURE-0018 | Administrator delegates project action/reach within and beyond one independently applicable per-action witness | contained grant intent may publish; broader intent unchanged | `none` or `AUTHORIZATION_DENIED` | RoleAssignment controller is sole publisher; no witness synthesis | security | AC-F18-09 |
| VS0-CF-F18-10 | FEATURE-0018 | Exact FEATURE-0016 ExecutionTarget read and qualify bindings for one accessible target | read/qualify apply only to registered target binding | none | protected action/target provenance; no provider-native effect | integration-contract | AC-F18-10 |
| VS0-CF-F18-11 | FEATURE-0018 | Qualifier-only holder attempts ExecutionTarget retirement | target remains unchanged | `AUTHORIZATION_DENIED` | denial audit; no retirement | security | AC-F18-11 |
| VS0-CF-F18-12 | FEATURE-0018 | Holder requests create-only ExecutionTarget authority under combined existing write action | no invented create-only grant/action is published | existing safe error | recorded granularity limitation; no action-registry mutation | compatibility | AC-F18-12 |
| VS0-CF-F18-13 | FEATURE-0018 | Exact Human requests approved two-hour JIT production DB role before activation deadline | one TimeBound RoleAssignment becomes active then ineffective at exact expiry | none | immutable approval/activation provenance and required audit | privileged-access | AC-F18-13 |
| VS0-CF-F18-14 | FEATURE-0018 | Requester attempts to approve own privileged request | request not approved or activated | `AUTHORIZATION_DENIED` | denied decision audited; no RoleAssignment intent/publication | security | AC-F18-14 |
| VS0-CF-F18-15 | FEATURE-0018 | Listed approver loses current eligibility before decision acceptance | request remains without that accepted approval | `AUTHORIZATION_DENIED` | denial audited; no downstream effect | security | AC-F18-15 |
| VS0-CF-F18-16 | FEATURE-0018 | Approval decision races ApprovalRequest expiry | exactly one valid terminal state wins; no unauthorized assignment | `none` for winner; `CONFLICT` for loser | one terminal evidence publication; no duplicate effect | concurrency | AC-F18-16 |
| VS0-CF-F18-17 | FEATURE-0018 | Same approved JIT request activated twice with same and conflicting replay fixtures | exactly one temporary RoleAssignment; validity never extends | safe replay or `CONFLICT` | one activation publication and one completed result | idempotency | AC-F18-17 |
| VS0-CF-F18-18 | FEATURE-0018 | Eligible Human invokes approved break-glass path with current assurance evidence | individual TimeBound access lasts no more than one hour | none | immediate protected audit and retrospective review due within 24 hours | privileged-access | AC-F18-18 |
| VS0-CF-F18-19 | FEATURE-0018 | Standing certification finds stale administrator assignment and selects Revoke | exact assignment becomes revoked exactly once | none | retained immutable review, usage, decision, and remediation evidence | review | AC-F18-19 |
| VS0-CF-F18-20 | FEATURE-0018 | Reviewed assignment version changes after immutable snapshot | review becomes `Stale`; assignment unchanged | none | no remediation or review-due advancement | concurrency | AC-F18-20 |
| VS0-CF-F18-21 | FEATURE-0018 | Incident action requires both applicable DB role and bounded change-window exception | `Allow` only when both independent authorities are current | `none` or `Deny` | exact role and exception provenance; neither synthesizes the other | security | AC-F18-21 |
| VS0-CF-F18-22 | FEATURE-0018 | Proposal attempts to except registered non-exceptionable tenant isolation control | no ExceptionGrant published | `Deny` using existing registered outcome | terminal protected denial evidence and audit | security | AC-F18-22 |
| VS0-CF-F18-23 | FEATURE-0018 | Active ExceptionGrant receives valid linked early-revocation proposal | original record remains immutable and exact grant becomes ineffective at revoke instant | none | one linked `effect=Revoke` record and audit | exception | AC-F18-23 |
| VS0-CF-F18-24 | FEATURE-0018 | Proposed ExceptionGrant overlaps same control/subject/scope/interval | existing grant unchanged; new grant absent | `CONFLICT` | no merge, precedence inference, or partial publication | exception | AC-F18-24 |
| VS0-CF-F18-25 | FEATURE-0018 | FEATURE-0017 returns Allow while applicable RoleAssignment is expired | final authorization is `Deny` | `Deny` | protected decision provenance records both inputs; no effect | security | AC-F18-25 |
| VS0-CF-F18-26 | FEATURE-0018 | Compatible pinned FEATURE-0017 evaluation returns RequiresApproval | no effect; request exists only through the explicit approved flow | `RequiresApproval` | protected evaluation/requirement provenance; no implicit ApprovalRequest | policy-seam | AC-F18-26 |
| VS0-CF-F18-27 | FEATURE-0018 | FEATURE-0017 result is Indeterminate | final authorization is `Deny` | `Deny` | redacted protected decision/audit; no effect | security | AC-F18-27 |
| VS0-CF-F18-28 | FEATURE-0018 | Required AuditEvent append fails at JIT activation publication boundary | request is not Active; no RoleAssignment or completed replay published | `INTERNAL_ERROR` | atomic rollback/no publication | audit | AC-F18-28 |
| VS0-CF-F18-29 | FEATURE-0018 | Authorized end customer performs a routine low-risk action | business result returned without governance ceremony | none | identity/grants/guardrails/provenance/audit resolve internally | usability | AC-F18-29 |
| VS0-CF-F18-30 | FEATURE-0018 | Denied operation references an inaccessible target | protected state unchanged | `RESOURCE_NOT_FOUND` safe denial | redacted audit; no existence or policy-detail disclosure | security | AC-F18-30 |
| VS0-CF-F18-31 | FEATURE-0018 | Production Workload/System receives policy-permitted narrow Standing assignment | assignment effective, immediately revocable, and review-due | none | responsible Human, provenance, and review schedule retained | review | AC-F18-31 |
| VS0-CF-F18-32 | FEATURE-0018 | Guest RoleAssignment omits `notBefore` or `expiresAt` | no assignment published | `VALIDATION_FAILED` | no audit-changing publication or completed replay | contract | AC-F18-32 |
| VS0-CF-F18-33 | FEATURE-0018 | Two grants allow action while a mandatory applicable guardrail excludes it | final authorization is `Deny` | `Deny` | exact grant and guardrail provenance; no effect | security | AC-F18-33 |
| VS0-CF-F18-34 | FEATURE-0018 | Published RoleDefinition version is suspended, then validly restored | pinned assignments ineffective during suspension; other versions unaffected | `Deny` during suspension | suspension/restoration DecisionRecord and audit; no assignment rewrite | lifecycle | AC-F18-34 |
| VS0-CF-F18-35 | FEATURE-0018 | AuthorizationResult is replayed for a different principal/action/target/scope | result grants nothing and target state remains unchanged | `Deny` | fresh authorization/audit path; no bearer reuse | security | AC-F18-35 |
| VS0-CF-F18-36 | FEATURE-0018 | Federated provenance is stale or responsible Human is absent at a required boundary | privilege-increasing operation unchanged | `Deny` | safe protected denial; raw claim cannot repair canonical evidence | security | AC-F18-36 |
| VS0-CF-F18-37 | FEATURE-0018 | AccessReview usage telemetry is absent, incomplete, and complete in separate fixtures | absence is not non-use; no silent workload revocation | none | immutable coverage/usage evidence and explicit reviewer outcome | review | AC-F18-37 |
| VS0-CF-F18-39 | FEATURE-0018 | Administrator supplies holder, role version, scope, and validity; separately forges system-owned fields | valid exact intent may publish; forged request unchanged | `none` or `AUTHORIZATION_DENIED` | RoleAssignment controller supplies refs/provenance/audit as sole publisher | usability | AC-F18-39 |
| VS0-CF-F18-40 | FEATURE-0018 | Engineer submits exact time-bound privileged request | one requester-facing request status and expiry projection | none | internal request/evidence/resource separation retained | usability | AC-F18-40 |
| VS0-CF-F18-41 | FEATURE-0018 | Exact eligible approver decides privileged/exception request; includes cross-stage reuse attempt | valid decision affects only current stage; one Human cannot count in two stages | `none` or `AUTHORIZATION_DENIED` | immutable decision evidence; downstream effect remains separate | approval | AC-F18-41 |
| VS0-CF-F18-42 | FEATURE-0018 | Eligible independent reviewer selects Retain, Revoke, or Replace on exact snapshot | exact approved disposition or explicit stale outcome | none or existing safe error | usage summary, beneficiary set, decision, and remediation evidence retained | review | AC-F18-42 |
| VS0-CF-F18-43 | FEATURE-0018 | User proposes bounded exceptionable control exception | one request status/effective-period projection; grant only after approved terminal processing | none or approved terminal `Deny` | proposal/approval/grant evidence remains internally distinct | usability | AC-F18-43 |
| VS0-CF-F18-44 | FEATURE-0018 | `ExactResource(A)` grantor attempts resource B and ScopeOnly delegation | no broader RoleAssignment published | `AUTHORIZATION_DENIED` | no cross-assignment synthesis or partial grant | security | AC-F18-44 |
| VS0-CF-F18-45 | FEATURE-0018 | Direct or indirect Standing beneficiary attempts Retain/Replace self-certification; unresolved expansion fixture included | no access-preserving/replacing effect or due advancement | `REVIEWER_CONFLICT` or `REVIEW_BENEFICIARY_UNRESOLVED` | independently authorized Revoke remains available and audited | security | AC-F18-45 |
| VS0-CF-F18-46 | FEATURE-0018 | TimeBound membership administrator attempts to enable durable group-held assignment | no durable derived access; empty group addition grants nothing | `AUTHORIZATION_DENIED` | exact envelope/witness denial evidence; no Membership publication when expansion exceeds ceiling | security | AC-F18-46 |
| VS0-CF-F18-47 | FEATURE-0018 | Membership administrator uses group/relationship/role evidence to manufacture approval, JIT, break-glass, or review eligibility | eligibility remains unsatisfied | existing structural-validation or safe-denial error | no ApprovalRequest decision, activation, review effect, or membership-derived eligibility | security | AC-F18-47 |
| VS0-CF-F18-48 | FEATURE-0018 | Reviewer list is empty, indirect, stale, scope-incompatible, target-incompatible, or changes before remediation | review effect unchanged | existing F18-RD-16 safe error | eligibility and beneficiary conflict rechecked; no remediation/due advancement | security | AC-F18-48 |
| VS0-CF-F18-49 | FEATURE-0018 | Role grantor or later rule publisher attempts to make a role/group/relationship qualify for eligibility | eligibility remains exact-Human-only and unsatisfied | existing structural-validation or safe-denial error | no approval, activation, Retain/Replace, or eligibility side effect | security | AC-F18-49 |

### 4. Exact requirement-only local registrations

These cases supply explicit conformance for architecture/contract requirements
that are not exhausted by a single active end-user AC. They do not create new
acceptance criteria.

| ID | Owner | Exact input fixture | Expected state/result | Expected error/decision | Expected side effects | Gate | Exact REQ |
|---|---|---|---|---|---|---|---|
| VS0-CF-F18-50 | FEATURE-0018 | Active FEATURE-0018 contract fixture containing one future-owned/prohibited resource, field, adapter, credential flow, or provider-native IAM object | fixture rejected before evaluation/publication | existing architecture-drift or structural-validation failure | zero external I/O and no canonical/runtime publication | architecture | REQ-F18-01 |
| VS0-CF-F18-51 | FEATURE-0018 | Valid exact GovernanceProfile v1 component/applicability fixture; then unknown, ambiguous, scope-incompatible, broadened, and retired-version fixtures | valid fixture accepted as pinned rule evidence only; invalid fixtures rejected | none or existing F18-RD-18 validation failure | no profile assignment/effective-resolution effect and no silent successor migration | contract | REQ-F18-18 |
| VS0-CF-F18-52 | FEATURE-0018 | Exact three DecisionProfile registrations and mandatory authorization/approval/exception decision and AuditEvent fixtures, including required-append failure | only approved profiles/taxonomy accepted; required publication is atomic with evidence | none or `INTERNAL_ERROR` for required-append failure | exact immutable FEATURE-0013 evidence; no competing envelope or fourth profile | audit | REQ-F18-19 |
| VS0-CF-F18-53 | FEATURE-0018 | Valid and invalid immutable synthetic fixture graphs with deterministic UTC clock and external-call counter | invalid graph rejected before evaluation; valid repeated evaluation is identical | none or existing fixture-validation failure | `externalCallCount=0`; no IdP, policy-engine, workflow-engine, CloudProvider, or network I/O | conformance | REQ-F18-22 |
| VS0-CF-F18-54 | FEATURE-0018 | Standards matrix with every applicable invariant mapped to local proof, reused dependency, or named exclusion; then one unmapped invariant | complete matrix passes; incomplete matrix fails final architecture approval gate | conformance-gate failure, not a runtime Problem code | no runtime state or conformity claim created | standards | REQ-F18-23 |

### 5. Exact AC-to-conformance mapping

The mapping is closed as follows:

```text
AC-F18-01..AC-F18-37 -> VS0-CF-F18-01..VS0-CF-F18-37 (same suffix)
AC-F18-38             -> excluded; VS0-CF-F18-38 tombstone; no active proof
AC-F18-39..AC-F18-49 -> VS0-CF-F18-39..VS0-CF-F18-49 (same suffix)
```

Kiro must render the mapping as explicit individual rows in generated stage
artifacts; ranges may be used in architecture prose only.

### 6. Exact REQ-to-conformance mapping

| Requirement | Exact FEATURE-0018-local conformance IDs |
|---|---|
| REQ-F18-01 | VS0-CF-F18-50 |
| REQ-F18-02 | VS0-CF-F18-06 |
| REQ-F18-03 | VS0-CF-F18-01, VS0-CF-F18-02, VS0-CF-F18-05, VS0-CF-F18-06 |
| REQ-F18-04 | VS0-CF-F18-01, VS0-CF-F18-08 |
| REQ-F18-05 | VS0-CF-F18-01, VS0-CF-F18-03, VS0-CF-F18-04, VS0-CF-F18-07, VS0-CF-F18-08, VS0-CF-F18-32, VS0-CF-F18-36, VS0-CF-F18-46, VS0-CF-F18-47 |
| REQ-F18-06 | VS0-CF-F18-10, VS0-CF-F18-11, VS0-CF-F18-12 |
| REQ-F18-07 | VS0-CF-F18-34 |
| REQ-F18-08 | VS0-CF-F18-01, VS0-CF-F18-03, VS0-CF-F18-04, VS0-CF-F18-05, VS0-CF-F18-07, VS0-CF-F18-09, VS0-CF-F18-21, VS0-CF-F18-31, VS0-CF-F18-32, VS0-CF-F18-39, VS0-CF-F18-44 |
| REQ-F18-09 | VS0-CF-F18-01, VS0-CF-F18-03, VS0-CF-F18-05, VS0-CF-F18-08, VS0-CF-F18-09, VS0-CF-F18-25, VS0-CF-F18-29, VS0-CF-F18-33, VS0-CF-F18-35, VS0-CF-F18-44, VS0-CF-F18-46 |
| REQ-F18-10 | VS0-CF-F18-25, VS0-CF-F18-26, VS0-CF-F18-27 |
| REQ-F18-11 | VS0-CF-F18-26 |
| REQ-F18-12 | VS0-CF-F18-13, VS0-CF-F18-14, VS0-CF-F18-15, VS0-CF-F18-41, VS0-CF-F18-47, VS0-CF-F18-49 |
| REQ-F18-13 | VS0-CF-F18-13, VS0-CF-F18-15, VS0-CF-F18-16, VS0-CF-F18-41 |
| REQ-F18-14 | VS0-CF-F18-13, VS0-CF-F18-17, VS0-CF-F18-40, VS0-CF-F18-47, VS0-CF-F18-49 |
| REQ-F18-15 | VS0-CF-F18-18 |
| REQ-F18-16 | VS0-CF-F18-18, VS0-CF-F18-19, VS0-CF-F18-20, VS0-CF-F18-31, VS0-CF-F18-37, VS0-CF-F18-42, VS0-CF-F18-45, VS0-CF-F18-47, VS0-CF-F18-48, VS0-CF-F18-49 |
| REQ-F18-17 | VS0-CF-F18-21, VS0-CF-F18-22, VS0-CF-F18-23, VS0-CF-F18-24, VS0-CF-F18-43 |
| REQ-F18-18 | VS0-CF-F18-51 |
| REQ-F18-19 | VS0-CF-F18-52 |
| REQ-F18-20 | VS0-CF-F18-06, VS0-CF-F18-28, VS0-CF-F18-52 |
| REQ-F18-21 | VS0-CF-F18-16, VS0-CF-F18-17, VS0-CF-F18-27, VS0-CF-F18-30 |
| REQ-F18-22 | VS0-CF-F18-01..VS0-CF-F18-37, VS0-CF-F18-39..VS0-CF-F18-53 |
| REQ-F18-23 | VS0-CF-F18-54 |
| REQ-F18-24 | VS0-CF-F18-29, VS0-CF-F18-39, VS0-CF-F18-40, VS0-CF-F18-41, VS0-CF-F18-42, VS0-CF-F18-43 |

Kiro must expand the REQ-F18-22 ranges to explicit IDs wherever a generated
contract prohibits range notation.

### 7. Requirements-stage semantic restoration and cleanup

The following are transcription corrections from ADH-2026-070, not new
semantics:

1. Approval stages: a Human principal whose approval counts toward one stage of
   an ApprovalRequest cannot count toward another stage of that same request.
2. Accepted break-glass bypass intent: when bypass evidence becomes invalid
   before activation, the PrivilegedAccessRequest remains `Submitted` with
   readiness `Blocked`; no ApprovalRequest, RoleAssignment intent, or
   RoleAssignment is published.
3. Replay: retain current authentication, current authorization, and safe target
   visibility re-evaluation before returning a stored result; changed content
   conflicts and required audit failure publishes no state change or completed
   replay.
4. Audit atomicity: retain strict required-audit-obligation acceptance before
   every authorization-changing publication; failure publishes neither the
   state change nor a completed idempotency result.
5. AccessReview: transcribe F18-RD-16 without abbreviation that loses its exact
   mode-specific timing, start/completion conditions, immutable snapshot and
   usage-evidence requirements, exact-version `StandingCertification + Retain`
   due advancement, stale behavior, beneficiary/reviewer rechecks, and overdue
   consequences.
6. Writer presentation: list `VS0-WRITER-008` and `VS0-WRITER-022` as
   FEATURE-0018-owned. `VS0-WRITER-007` is a consumed shared writer registration.
   `VS0-WRITER-023` is FEATURE-0020-owned and may appear only under inherited or
   out-of-scope writer references, never under FEATURE-0018-owned writers.

These corrections must preserve `REQ-F18-01..24`, `AC-F18-01..49`, the excluded
status of `AC-F18-38`, and all approved titles verbatim.

### 8. Closed architecture citation ledger

Kiro must cite, not recompute, these exact authorities:

| Required design input | Sole semantic authority |
|---|---|
| FEATURE-0018 operation surface | ADH-2026-070 Appendix B, F18-RD-02 operation-surface table |
| Canonical action and target-binding registry | ADH-2026-070 Appendix B, F18-RD-06 action and discriminated target-binding tables |
| GovernanceProfile v1 components and applicability | ADH-2026-070 Appendix B, F18-RD-18 component schemas and applicability registrations |
| FEATURE-0018 AuditEvent taxonomy and publication boundary | ADH-2026-070 Appendix B, F18-RD-20 taxonomy; F18-RD-19/20 evidence rules |
| Local fixture/race/replay/redaction inventory | ADH-2026-070 Appendix B, F18-RD-22 |
| Resource semantics, validity, approval, privilege, review, and exception behavior | ADH-2026-070 Appendix A, F18-RD-02..18 and F18-RD-21 |

No generated requirement, design, or task may replace an exact table with a
newly inferred action, target mode, component, event type, state, or error.

## Rationale

Stable local identifiers turn approved prose obligations into independently
addressable evidence without reopening the authorization model. One local case
per active AC makes review and failure diagnosis simple; five requirement-only
cases prevent architecture, evidence-adoption, determinism, and standards
requirements from becoming untestable or falsely attached to unrelated user
journeys. The tombstone preserves identifier history, and explicit ownership
prevents inherited FEATURE-0015/0020 evidence from being misrepresented as
FEATURE-0018-local proof.

## Reuse-before-build assessment

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse the existing Slice-0 registry, traceability matrix, stable-ID rules,
    stage steering, deterministic fixture approach, and feature-gate machinery.
  - Reuse ADH-2026-070/DEC-0060 semantics and FEATURE-0012/0013/0016/0017
    dependency contracts by exact reference.
  - Apply the existing FEATURE-0011 reuse assessment and FEATURE-0015 local
    conformance conventions; no parallel conformance framework is created.
- Sovrunn-owned responsibility summary:
  - FEATURE-0018 owns its local case registrations, local fixtures, AC/REQ
    mappings, and proof that its approved control-plane semantics conform.
  - The generic Slice-0 registry/checker owns identifier integrity and
    cross-artifact traceability enforcement.
- Non-goals summary:
  - No new runtime resource, action, state, error, route, controller, adapter,
    persistence choice, provider-native IAM behavior, or external I/O.
  - No change to dependency-owned conformance semantics or writer ownership.
  - No design, task, or Go implementation generation in this handoff.

## Phase impact

- Current phase allowed? Yes
- If not current phase, target phase: Not applicable
- Current phase boundary impact:
  - Documentation, registry, traceability, controlled requirements generation,
    and feature-gate evidence only. Phase 2R in-memory and zero-external-I/O
    boundaries remain unchanged.

## Conflict check

- Conflicts with accepted DEC/RFC? No
- Conflicting decisions, if any:
  - None. This corrects the executable-traceability mismatch between
    F18-RD-22/approved AC ledgers and the mandatory Slice-0 generated-output
    contract.
- Resolution required:
  - Approved ADH only. No new ACR, DEC, RFC, or architecture-baseline update is
    required because runtime semantics and accepted architecture do not change.

## Required action

- Update architecture doc
- Update Kiro requirements.md
- Update traceability matrix
- Update feature gate/checks

Do not update Kiro `design.md` or `tasks.md` as part of this reconciliation.

## Impacted files

Kiro must update or inspect the following exact files atomically after human
approval:

- `docs/reviews/architecture-decision-handoffs/ADH-2026-071-feature-0018-conformance-executability-reconciliation.md`
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`
- `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `docs/features/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/requirements.md`
- `.automation/context-projections/FEATURE-0018.requirements.spec.json`
- `.automation/context-projections/FEATURE-0018.requirements.md`
- `.automation/context-projections/FEATURE-0018.requirements.manifest.json`
- `.automation/state/FEATURE-0018.json`
- `scripts/kiro-semantic-check.py`

Kiro must inspect but not mutate the immutable historical approval snapshot:

- `docs/reviews/architecture-readiness/FEATURE-0018-final-architecture-approval-manifest.md`

Kiro must not manually edit the generated reviewer result. It must be regenerated
by the controlled review route:

- `.automation/reviews/FEATURE-0018/requirements.review.json`

`FEATURE-0018.control.json` changes only if its current generic schema already
contains a governed field needed to pin this reconciliation. Kiro must not invent
a feature-specific control-schema field.

## Impacted features

- FEATURE-0018: receives exact local executable conformance and requirements
  traceability; runtime authority is unchanged.
- FEATURE-0015: `VS0-CF-X03` ownership and semantics remain unchanged and are
  consumed only as inherited evidence.
- FEATURE-0020: `VS0-WRITER-023` ownership remains unchanged; FEATURE-0018 gains
  no ProfileAssignment authority.
- FEATURE-0012, FEATURE-0013, FEATURE-0016, FEATURE-0017: dependency contracts
  remain unchanged and are consumed by exact reference only.

## Acceptance criteria for Kiro update

- [ ] Handoff validated against Architecture Operating System files.
- [ ] Human approval is recorded before any application.
- [ ] No unapproved baseline, DEC, RFC, ACR, or runtime-semantic change is made.
- [ ] Phase 2R scope and zero-external-I/O boundaries are preserved.
- [ ] Reuse-before-build section is preserved.
- [ ] All 48 active ACs map one-to-one to existing registered
  `VS0-CF-F18-*` IDs with the same suffix.
- [ ] `AC-F18-38` remains excluded and `VS0-CF-F18-38` is a non-active,
  permanently consumed tombstone.
- [ ] All 24 REQs map to at least one exact registered FEATURE-0018-local case.
- [ ] `VS0-CF-F18-50..54` remain requirement-only cases and create no ACs.
- [ ] Registry, traceability matrix, feature authority, architecture, generated
  context hashes, approval evidence, and requirements are reconciled atomically.
- [ ] `VS0-CF-X03` is not claimed as FEATURE-0018-local proof.
- [ ] `VS0-WRITER-023` is not listed as FEATURE-0018-owned.
- [ ] The six semantic transcriptions in section 7 exactly match ADH-2026-070.
- [ ] Operation, action/target-binding, GovernanceProfile, and AuditEvent inputs
  cite the closed ledgers in section 8 and are not recomputed.
- [ ] The generic semantic checker rejects any active AC without a registered
  exact conformance mapping, any REQ without a registered local mapping, any
  orphan/duplicate/reused ID, and any active use of the tombstone.
- [ ] No `design.md`, `tasks.md`, Go code, runtime test, or external integration
  is generated by this reconciliation.
- [ ] `feature-control.py validate`, the requirements semantic checker,
  documentation strict build, architecture-drift checks, applicable Structurizr
  checks, and the controlled independent requirements review all pass.

## Explicit instructions to Kiro

- Do not introduce new decisions beyond this handoff.
- Do not reinterpret abbreviated conformance outcomes; transcribe the full
  mapped ADH-2026-070 authority.
- Do not create a new runtime error or select a more specific runtime error where
  ADH-2026-070 does not already select one.
- Do not renumber REQ-F18, AC-F18, F18-SCN, existing `VS0-CF-*`, schema, writer,
  state, failure, action, or audit identifiers.
- Do not activate `VS0-CF-F18-38` or reuse its suffix.
- Do not count inherited or integration evidence as FEATURE-0018-local proof.
- Do not modify `CURRENT_ARCHITECTURE_BASELINE.md`; no baseline update is
  authorized.
- Do not move future-feature implementation into Phase 2R.
- Do not modify Go code or generate implementation work.
- Stop with `ARCHITECTURE_DECISION_REQUIRED` if applying this handoff would
  require any observable semantic choice not already closed by ADH-2026-070.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-30
- Notes: Approval authorizes Kiro to apply only the conformance and traceability
  correction above. It is not requirements-stage authorization and does not
  authorize design or tasks.

This Architecture Decision Handoff is ready for Kiro validation and repo update only after human approval.
