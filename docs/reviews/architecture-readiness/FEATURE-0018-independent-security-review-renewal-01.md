---
doc_type: independent_security_review
feature: FEATURE-0018
renewal: 01
status: rejected-blocking-findings
reviewer: Codex isolated security architecture reviewer 02
review_date: 2026-08-30
updated: 2026-08-30
---

# FEATURE-0018 Independent Security Review Renewal 01

## 1. Disposition, integrity, and independence

**Disposition: REJECT — BLOCKING SECURITY AND PAYLOAD-RECONCILIATION
FINDINGS REMAIN.**

This review was performed by **Codex isolated security architecture reviewer
02**. Independence basis: **fresh fork with no conversation history; did not
author the original or corrected reviewed payload**. The prior rejection was
read as historical evidence, but its conclusions and the corrected package's
prepared claims were not accepted without independent verification.

Reviewed manifest:

```text
docs/reviews/architecture-readiness/FEATURE-0018-final-architecture-approval-manifest.md
SHA-256 015b3508a31ee038e08a268bb217132eba3d14373dfcced2b7387719d82d31b7
```

The manifest hash matched the required value. Its table contained exactly 31
payload files, and the current SHA-256 of every listed file matched its pinned
digest. The review record is intentionally outside that table, so this signed
record does not create a circular digest dependency.

The corrected normative contract is materially stronger than the rejected
payload, but an exact-manifest review cannot approve only its strongest files.
The manifest also pins a contradictory architecture digest, and the corrected
access-review conflict set leaves one indirect beneficiary path unclosed.

## 2. Review boundary and method

The complete 31-file payload was reviewed, including the proposed decision and
change request; the indivisible ADH-2026-070 core and both normative appendices;
canonical resource, contract, glossary, phase, context, registry, and
traceability changes; FEATURE-0013 and FEATURE-0016 registrations; the
FEATURE-0018 architecture and feature contracts; control manifest; reuse
evidence; digest; industry matrix; leakage proof; standards mapping; and threat
and abuse ledger.

The assessment covered zero-trust correctness, authorization algebra,
scope/target/resource reach, identity and group behavior, delegated
administration, separation of duties, JIT and break-glass assurance,
approval/effect atomicity, writer ownership and concurrency, access review,
exception processing, evidence and privacy, native-IAM separation, simplicity,
architecture leakage, completeness, contradictions, and abuse cases.

This record creates no requirements, design, tasks, implementation, or
architecture semantics and modifies no reviewed payload file.

## 3. Independent security assessment

| Area | Independent determination |
|---|---|
| Zero trust and authorization algebra | The normative contract correctly requires current identity and dependency state, independently applicable grants, mandatory-guardrail intersection, deny/indeterminate fail-closed behavior, exact target binding, current evidence, and non-bearer `Allow | Deny` results. Authorization is re-evaluated at sensitive effect publication. The lower-precedence digest nevertheless publishes contradictory membership and target-input algebra; see F18-REN-001. |
| Scope, target, and resource reach | `ExactResource`, `CreateParent`, and `ScopeOnly` are closed discriminated modes. A resource-narrowed assignment contributes only to its exact UID and compatible target kind. Delegated publication now requires a complete per-action witness for scope, target mode, exact resource reach, validity, and delegation capability; partial assignments cannot be synthesized. F18-SEC-002 is closed. |
| Identity and group behavior | Stable issuer/subject identities, trusted direct membership, provenance/freshness, Guest expiry, group suspension, owner accountability, and raw-claim rejection are sound at architecture level. Groups cannot authenticate and privileged roles cannot be group-held. Workload/System beneficiaries reached through a group expose an access-review conflict gap; see F18-REN-002. |
| Delegated administration and SoD | Grantor reach containment and approval conflict checks are substantially complete. Approval SoD includes requester, direct beneficiary, group owner/current members, and the responsible Human for beneficiary Workload/System principals, with re-evaluation at decision and effect publication. Access-review SoD does not apply the same compounded indirect resolution; see F18-REN-002. |
| JIT and break-glass | Normal JIT and break-glass activation now have one non-weakenable floor: current, provider-neutral, phishing-resistant AAL2-or-higher AssuranceEvidence no more than 15 minutes old. Break-glass additionally has individual pre-authorization, a one-hour ceiling, no mutation renewal, immediate protected audit, and retrospective review within 24 hours. F18-SEC-004 is closed. |
| Approval/effect atomicity | Approval admission and downstream effect publication are distinct. Current authorization, evidence, dependency, SoD, concurrency, and audit-obligation acceptance are rechecked at the effect boundary. Privileged request activation and RoleAssignment publication succeed together or neither publishes. Production/storage realization remains downstream proof. |
| Single writer and concurrency | F18-RD-02/F18-RD-08, the architecture contract, and VS0-WRITER-008 now make the RoleAssignment controller the sole canonical spec/status/lifecycle/effect writer; workflows submit immutable intents only. Optimistic concurrency, idempotency binding, exact replacement/revocation, and first-valid-decision rules are appropriately specified. The exact payload's digest retains the former split-writer table, so end-to-end F18-SEC-001 closure is not established; see F18-REN-001. |
| Access review | Immutable snapshots, kind-specific dispositions, explicit stale/incomplete evidence, exact-item remediation, and atomic Retain/Replace publication are strong. `ExpiredIncomplete` is correctly a campaign state in the normative appendix. The responsible Human of a Workload/System member benefiting through an AccessGroup-held assignment can still certify that authority; see F18-REN-002. |
| Exception processing | Exceptions are immutable, time-bounded, exact-control and exact-subject/scope grants produced only after approval. Exceptionability, allowed effect, narrowing, overlap, current evidence, authoritative expiry, and revocation all fail closed. Exceptions cannot override identity, scope isolation, immutable audit, authorization, or other non-exceptionable controls. |
| Evidence and privacy | Decision/audit pins, provenance, safe resolution order, redaction, non-disclosure, audit-before-publication, and no-secret/no-bearer boundaries are appropriate. The contradictory digest and overclaimed leakage proof make the assembled approval evidence unreliable; see F18-REN-001. |
| Native IAM separation | Sovrunn and native IAM are an intersection: both must allow and either may deny. Provider objects, credentials, permission translation, and native effects remain outside FEATURE-0018. Future adapter compromise and credential reach remain explicit integration risks, not reasons to weaken this boundary. |
| Simplicity | Six progressive-disclosure journeys keep routine authorization invisible and preserve distinct internal resources without adding a facade or generic condition language. This is acceptable at architecture stage. Stale digest vocabulary would undermine that simplicity if consumed downstream. |
| Completeness and leakage | No future IdP, workflow engine, persistent store, real adapter, native effect, or FEATURE-0020 governance resolver is pulled into scope. However, the manifest pins mutually inconsistent writer, authorization, review, and lifecycle descriptions while the leakage proof claims complete reconciliation. The approval package is therefore not complete or contradiction-free. |

## 4. Prior-finding closure determination

| Prior finding | Renewal determination | Evidence |
|---|---|---|
| F18-SEC-001 — conflicting RoleAssignment writers | **NOT CLOSED end-to-end** | The normative appendix, feature architecture, canonical changes, and VS0-WRITER-008 now correctly assign all canonical RoleAssignment writes to the RoleAssignment controller. However, the exact-manifest architecture digest still assigns spec to `Authorized grant controller` and status to `Authorization controller`. The prior correction required every affected writer ledger to be reconciled, and the package claims that reconciliation. See F18-REN-001. |
| F18-SEC-002 — delegation ceiling omitted resource reach | **CLOSED at architecture level** | F18-RD-09 now requires one independently applicable, delegable witness for each proposed action covering the exact scope, target-binding mode, optional exact resource, validity, and delegation capability. Resource A cannot delegate resource B or scope-wide reach, and cross-assignment synthesis is prohibited. The registry and negative-case inventory mirror this rule. |
| F18-SEC-003 — access-review self-certification | **NOT CLOSED** | Direct Human holders, AccessGroup owners/current direct members, and direct Workload/System holders' responsible Humans are excluded for Retain/Replace. The set does not exclude the responsible Human for a Workload/System principal that is itself a member/beneficiary of an AccessGroup-held RoleAssignment. See F18-REN-002. |
| F18-SEC-004 — contradictory normal-JIT assurance | **CLOSED at architecture level** | F18-RD-14/F18-RD-15 and their mirrors uniformly require recent phishing-resistant AAL2-or-higher evidence for both normal JIT and break-glass activation; policy may strengthen but cannot weaken it. |
| F18-SEC-005 — stale or falsely reconciled approval evidence | **NOT CLOSED** | The exact-manifest digest retains multiple superseded semantics, while the leakage proof and manifest describe the payload as reconciled. Examples include writers, membership algebra, omitted `ScopeOnly`, qualitative Standing policy, review evidence/state, and lifecycle vocabulary. See F18-REN-001. |

Result: **two of five prior findings are closed; three remain open in the exact
payload**. The corrected package therefore cannot pass renewal review.

## 5. Blocking findings

### F18-REN-001 — Exact payload contains stale, contradictory security semantics

**Severity:** High

**Classification:** Blocking; reconciliation integrity, single-writer safety,
authorization determinism, and architecture leakage

The manifest pins `FEATURE-0018-architecture-digest.md`, yet that document still
contains superseded or internally contradictory semantics, including:

- section 2 first lists `AccessGroup` in the inherited canonical inventory and
  then states that the inherited inventory contains no access-group resource;
- section 4.1 assigns `RoleAssignment` specification to an `Authorized grant
  controller` and status/result to an `Authorization controller`, contradicting
  the sole RoleAssignment-controller rule in F18-RD-02/F18-RD-08 and corrected
  VS0-WRITER-008;
- F18-RD-08 permits production Workload/System Standing access using the
  qualitative phrase `narrowly scoped`, whereas the normative contract makes
  exact holder/scope/target registrations authoritative and gives qualitative
  narrowness no authorization meaning;
- F18-RD-09 unconditionally requires active Organization or CloudProvider
  Membership, contradicting the normative Platform/CloudPlatform exemption;
- its `AuthorizationInput` names only exact-resource and create-parent inputs,
  omitting the closed `ScopeOnly` mode;
- F18-RD-16 describes a `usageEvidenceRef` even though the normative contract
  uses an embedded `UsageEvidenceSummary` and expressly creates no usage-
  evidence resource, and it lists `ExpiredIncomplete` as a review outcome even
  though the normative contract makes it a campaign state rather than an item
  disposition;
- F18-RD-18 describes `Superseded` as a lifecycle state, while the normative
  contract treats supersession as informational and closes the lifecycle as
  Draft to Published to Retired; and
- the access-review journey repeats the stale `usageEvidenceRef` contract.

The ADH acceptance text and industry-matrix evidence-package section also still
name the scenario inventory only through F18-SCN-43, although the corrected
matrix adds F18-SCN-44 and F18-SCN-45. These are not harmless editorial
differences: they can produce a second writer, an invalid membership gate,
omitted target-mode handling, or the wrong review/lifecycle model. The leakage
proof's claim that all conflicts are reconciled and the manifest's prepared
consistency claim are consequently false for the pinned bytes.

**Required correction:** reconcile or explicitly retire every stale digest and
traceability statement against its controlling F18-RD rule; re-audit the entire
36-item leakage ledger rather than only these examples; make scenario inventory
counts consistent; regenerate all affected digests and the final manifest; and
obtain another independent review of the new exact payload.

### F18-REN-002 — Nested group/workload beneficiary can self-certify access

**Severity:** High

**Classification:** Blocking; separation of duties and indirect-beneficiary
access-review bypass

F18-RD-12's approval conflict set correctly reaches the responsible Human for a
beneficiary Workload/System principal, including when the beneficiary is
resolved indirectly. F18-RD-16's review conflict set is narrower. It includes:

1. a direct Human RoleAssignment holder;
2. the owner and each current direct principal member of an AccessGroup holder;
   and
3. the responsible Human for a Workload/System *holder*.

An AccessGroup may directly contain Workload/System principals. For an
AccessGroup-held RoleAssignment, the workload is a member but its responsible
Human is neither that direct member nor the RoleAssignment holder. Because a
reviewer must be a Human, that responsible Human can remain eligible to Retain
or Replace the group authority exercised by the workload they control, and can
advance the review due date at the same publication boundary. The generic
statement that a beneficiary may not self-certify does not define this missing
resolution, and the matrix/threat-ledger shorthand does not add normative
semantics.

**Required correction:** for reviewed and proposed assignments, define the
conflict set to include the responsible Human for every current Workload/System
beneficiary, including a beneficiary reached through an AccessGroup. Recompute
that set at disposition acceptance and at the atomic Retain/Replace/due-
advancement publication boundary. Add a deterministic
AccessGroup-to-Workload/System-to-responsible-Human negative case, align the
matrix and threat ledger, regenerate the manifest, and obtain renewed review.

## 6. Abuse-case disposition

| Attack or abuse case | Result |
|---|---|
| Confused deputy in delegated grant | **Controlled in the normative correction.** Per-action complete-reach witnesses and no cross-assignment synthesis close resource A-to-B and resource-to-scope escalation. The stale digest remains an unsafe alternate description. |
| Direct privilege escalation | Exact role version, action registry, target mode, scope/resource containment, validity, classification, current authorization, and non-enlarging controller behavior fail closed. |
| Stale or replayed approval | Exact proposal digest, bounded evidence expiry, current checks at effect publication, idempotency binding, and changed-content conflicts are sound. Later approver ineligibility does not revoke an already admitted approval by explicit design and remains a bounded residual risk. |
| Group-membership or ownership race | Authorization and approval re-evaluate current group state; review Retain/Replace re-evaluates its stated conflict set. The nested responsible-Human omission remains exploitable (F18-REN-002). Linearizable repository proof is downstream. |
| Revoke/use and suspend/use races | Non-bearer decisions, evaluation-time current-state requirements, no-grace expiry, and authorization dependency checks support fail-closed semantics. Executable linearization remains required. |
| Activation-deadline race | The complete request/assignment/audit publication must occur before the authoritative deadline; retry cannot extend validity. |
| Direct approval self-dealing | Controlled at approval decision and downstream publication. Direct Human access-review self-certification is now controlled. |
| Indirect approval/review self-dealing | Approval resolves group and responsible-party beneficiaries. Review misses the responsible Human behind a Workload/System group member (F18-REN-002). |
| Audit failure or evidence tampering | Protected audit-obligation acceptance precedes authoritative publication; immutable pins and system ownership are sound. This creates a deliberate audit-availability denial-of-service dependency. |
| Idempotency collision | Actor, operation, exact target/scope, request content, and digest binding distinguish same-content replay from conflicting reuse. |
| Cross-tenant or cross-scope enumeration | Validation/resolution order, scope-local lookup, safe denial/not-found behavior, and redacted projections are suitable at contract level. |
| Exception enlargement | Exact control/subject/scope, allowed-effect intersection, non-exceptionable invariants, current evidence, bounded time, overlap rules, and immutable grant semantics prevent silent enlargement. |
| Compromised native adapter | Outside FEATURE-0018's execution boundary. Future least-privileged identities, credential isolation/rotation, egress/call-path control, correlation, and rapid revocation must be proven. |
| Native-IAM disagreement | Controlled by intersection: Sovrunn Allow and native Allow are both required; either Deny blocks. Neither decision substitutes for the other. |
| Break-glass abuse | Individual pre-authorization, fresh phishing-resistant assurance, exact scope, one-hour maximum, no renewal mutation, immediate audit, expiry, revocation, and retrospective review are strong at architecture level. |

## 7. Residual risks and non-blocking observations

These are not approval conditions because the current disposition is Reject.
They must be re-evaluated after the blocking corrections:

- executable conformance, storage atomicity, race, retry, and failure-injection
  evidence do not yet exist;
- identity lifecycle, authenticators, trusted provisioning, and native adapter
  credential security remain future integration responsibilities;
- existing effective access intentionally does not fail solely because an
  AccessGroup owner or Workload/System responsible Human becomes inactive;
- overdue production Workload/System Standing review may enter escalation
  rather than automatic ineligibility under a registered exception;
- publication-time delegation intentionally does not cascade after later
  grantor revocation;
- protected audit availability is a fail-closed control-plane dependency and
  denial-of-service target; and
- FEATURE-0016's currently registered create/retire action granularity is an
  explicit bounded limitation until that owning feature changes it.

These risks are visible and bounded by the proposed architecture; none excuses
F18-REN-001 or F18-REN-002.

## 8. Read-only validation evidence

The following relevant commands were executed without modifying reviewed
payload files:

```text
shasum -a 256 docs/reviews/architecture-readiness/FEATURE-0018-final-architecture-approval-manifest.md
=> 015b3508a31ee038e08a268bb217132eba3d14373dfcced2b7387719d82d31b7

ruby -rdigest <manifest-table SHA-256 verifier>
=> payload_count=31
=> payload_digest_verification=PASS

ruby -rjson <control-manifest parser>
=> json_parse=PASS

ruby -ryaml <VS-000 registry safe-load with aliases enabled>
=> yaml_parse=PASS

ruby <F18-RD heading uniqueness check across Appendices A/B>
=> F18-RD-01 through F18-RD-24 each occurred exactly once
=> decision_heading_uniqueness=PASS

make arch-handoff-check HANDOFF=docs/reviews/architecture-decision-handoffs/ADH-2026-070-feature-0018-governance-iam-approval-exception-foundation.md
=> PASS (with the expected warning that the handoff is not Approved)

scripts/reuse-assessment-check.sh FEATURE-0018 --assessment docs/reviews/reuse-assessments/FEATURE-0018-approval-evidence.md --mode strict --skip-rac03 --skip-rac13
=> PASS

make phase2-scope-check FEATURE=FEATURE-0018
=> PASS

make phase2r-drift-check
=> PASS (15 checks)

git diff --check
=> PASS
```

`mkdocs build --strict` was not rerun because it writes generated output and
this review was authorized to create only this review record. The manifest's
prepared build claim was treated as a prepared claim, not independent proof.
No Structurizr-level boundary or flow was changed or proposed by this review.

## 9. Reviewer signature

| Field | Value |
|---|---|
| Reviewer identity | Codex isolated security architecture reviewer 02 |
| Independence basis | fresh fork with no conversation history; did not author the original or corrected reviewed payload |
| Review date | 2026-08-30 |
| Exact manifest reviewed | SHA-256 `015b3508a31ee038e08a268bb217132eba3d14373dfcced2b7387719d82d31b7` |
| Payload verification | All 31 manifest-listed file digests independently matched |
| Prior-finding closure | F18-SEC-002 and F18-SEC-004 closed; F18-SEC-001, F18-SEC-003, and F18-SEC-005 not closed end-to-end |
| Renewal findings | F18-REN-001 and F18-REN-002; two High blocking findings |
| Disposition | **Reject** |
| Signature/evidence reference | Signed 2026-08-30 by Codex isolated security architecture reviewer 02 against manifest SHA-256 `015b3508a31ee038e08a268bb217132eba3d14373dfcced2b7387719d82d31b7` |

This rejection is independent security-review evidence only. It does not alter
the reviewed architecture, approve reuse, or authorize requirements, design,
tasks, implementation, or any downstream gate.
