---
doc_type: independent_security_review
feature: FEATURE-0018
status: rejected-blocking-findings
reviewer: Codex isolated security architecture reviewer
review_date: 2026-08-29
updated: 2026-08-29
---

# FEATURE-0018 Independent Security Review Record

## 1. Integrity and independence statement

This review was performed by **Codex isolated security architecture reviewer**
on a fresh fork with no conversation history. The reviewer did not author or
reconcile the reviewed payload and did not defer to its prepared conclusions.

Reviewed manifest:

```text
docs/reviews/architecture-readiness/FEATURE-0018-final-architecture-approval-manifest.md
SHA-256 c447e522937846ad490fc1f5c9a3198e2dd4eb129444cd699c8c436b526282d7
```

The manifest hash independently matched the expected value. All 31 payload
files listed by that manifest independently passed SHA-256 verification before
review. The review record is intentionally excluded from the manifest payload,
as stated by the manifest, so completing this record does not create a circular
digest dependency.

Current outcome: **REJECT — BLOCKING SECURITY AND RECONCILIATION FINDINGS**.

## 2. Review inputs and boundary

The review covered the exact 31-file manifest payload, including ACR-2026-002,
DEC-0060, the indivisible ADH-2026-070 core and Appendices A/B, the proposed
canonical model/catalog/glossary reconciliation, FEATURE-0013 and FEATURE-0016
registrations, VS-000 registries/traceability, FEATURE-0018 architecture and
feature contracts, control manifest, reuse evidence, architecture digest,
industry matrix, leakage proof, standards mapping, and threat/abuse ledger.

The review assessed only architecture-stage security. It generated no
requirements, design, tasks, implementation, or architecture semantics.

## 3. Mandatory reviewer questions

| Area | Independent finding |
|---|---|
| Authorization algebra | Grant union, guardrail intersection, deny/indeterminate override, exact assignment applicability, and non-bearer results are directionally sound. The delegation ceiling is incomplete for resource-narrowed authority, however: it requires granted actions and target-scope authority but never defines resource-reach containment. See F18-SEC-002. |
| Scope/target | Ordinary authorization uses exact `ExactResource`, `CreateParent`, or `ScopeOnly` registrations and fails closed on malformed/ambiguous bindings. Delegated grant publication can still enlarge exact resource reach because the ceiling lacks a resource-set comparison. See F18-SEC-002. |
| Identity/group | Stable issuer/subject identity, direct canonical group membership, provenance/freshness, Guest expiry, suspension/revocation, and raw-claim rejection are sound at contract level. Existing Workload/System access after responsible-party inactivity and existing group access after owner inactivity are explicit availability tradeoffs and remain residual risks. |
| Delegation/SoD | Approval conflict sets correctly include direct and indirect beneficiaries, group owners/members, and responsible parties at decision and effect publication. Delegation resource reach is incomplete, and access-review `Retain` lacks equivalent beneficiary conflict checks. See F18-SEC-002 and F18-SEC-003. |
| Privileged/break-glass | Eligibility, activation deadline, time bounds, no-grace expiry, atomic activation, assurance age, break-glass phishing resistance, immediate audit, and retrospective review are strong. Normal-JIT phishing-resistance requirements conflict across the payload. See F18-SEC-004. |
| Approval/effect atomicity | Approval, ordinary assignment, privileged activation, and exception publication are correctly separated and fail closed conceptually. Cross-store/production atomicity remains downstream proof; no real persistence is authorized here. |
| Writers/concurrency | The normative sole RoleAssignment publisher conflicts with the VS-000 writer registry, which names a delegated administrator as writer of `RoleAssignment.spec`. This is a direct multiple-writer architecture defect. See F18-SEC-001. |
| Review/exception | Snapshot immutability, stale/incomplete usage qualification, exact remediation, and exception narrowing/overlap/revocation are generally complete. Standing certification can be self-dealt by ordinary direct holders, AccessGroup beneficiaries/owners, or responsible parties. See F18-SEC-003. |
| Evidence/privacy | Protected decision/audit provenance, safe denial ordering, redaction, audit-before-publication, and no-bearer results are well formed. Stale readiness mirrors contradict the normative contract and falsely report reconciliation, making the evidence package unsafe as an approval input. See F18-SEC-005. |
| Native IAM | The Sovrunn/native-IAM intersection and mutual non-substitution of Allow are correct. FEATURE-0018 does not mitigate a compromised future adapter identity making direct native calls; that remains a named future integration risk and requires separate least-privilege, credential, correlation, and revocation proof. |
| Architecture leakage | No provider-native IAM resource, credential, effect, real adapter, effective-governance resolver, or later-feature implementation is authorized. The stale digest/matrix semantics and contradictory writer registry are architecture-evidence leakage/conflict. See F18-SEC-001 and F18-SEC-005. |
| Simplicity | The six progressive-disclosure journeys preserve distinct internal resources without facade authority or requester-authored provenance. This simplicity boundary is acceptable and does not offset the blocking control defects. |

## 4. Concrete findings

| ID | Severity | Classification | Finding and evidence | Security impact | Required correction |
|---|---|---|---|---|---|
| F18-SEC-001 | High | Blocking | F18-RD-02/F18-RD-08, DEC-0060 and the canonical model make the RoleAssignment controller the sole publisher/spec/status writer, but `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` (`VS0-WRITER-008`) assigns `RoleAssignment.spec` to `delegated-identity-governance-admin`. The leakage proof nevertheless marks the writer conflict reconciled. | A direct administrator/spec-writer path can bypass immutable materialization intent, current dependency/privilege revalidation, atomic audit publication, and the sole-writer race boundary. Two incompatible authorities also make implementation nondeterministic. | Reconcile VS-000 so callers/controllers submit only the exact accepted intent and the RoleAssignment controller alone publishes canonical spec/status/revoke/replace/expiry effects. Update every affected writer ledger and rerun the 36-conflict proof. |
| F18-SEC-002 | High | Blocking | F18-RD-08 makes `resourceRef` a security-relevant exact narrowing restriction. F18-RD-09's grantor ceiling requires `roleassignment.grant`, every granted action, target-scope authority, validity, and optional delegation authority, but does not require the grantor's authority for each granted action to contain the proposal's exact resource reach. Matrix C05 and threat F18-THR-04 likewise test action/scope/duration/delegation but omit resource reach. | A grantor whose action is limited to resource A can satisfy the documented action/scope ceiling and grant the action for resource B or for the whole scope. This is a confused-deputy privilege escalation through delegated administration. | Define a deterministic per-action reach-containment rule: a resource-narrowed proposal requires independently applicable grantor authority for that exact resource; a scope-wide proposal requires scope-wide authority. No restrictions from different assignments may synthesize broader reach. Add resource A-to-B and resource-to-scope delegation negatives. |
| F18-SEC-003 | High | Blocking | F18-RD-16 prohibits only reviewing one's own *privileged* access. Privileged roles are TimeBound and cannot use the Standing path, while a `StandingCertification + Retain` for ordinary access advances `nextReviewDueAt`. Nothing bars a direct Human holder, a beneficiary/member or owner of an AccessGroup-held assignment, or a Workload/System responsible party from retaining the authority from which they benefit. Matrix E02 claims reviewer independence while encoding the same narrower privileged-only rule. | A beneficiary can repeatedly self-certify standing ordinary access, including sensitive read/decision capabilities classified Ordinary, defeating separation of duties and access-certification independence. | Apply an explicit review conflict set to every disposition that retains or replaces effective authority and to due-time advancement. At minimum exclude the direct holder, current beneficiary group owner/members, and Workload/System responsible party; define deterministic behavior for changed conflicts at decision/remediation publication and add direct/indirect self-review tests. |
| F18-SEC-004 | Medium | Blocking | Normal-JIT assurance is not deterministic across the payload. F18-RD-14 requires current AAL2 for approval-path JIT, while F18-RD-15 adds `phishingResistant=true` specifically for break-glass. The proposed canonical contract catalog instead describes every `PrivilegedAccessRequest` as requiring strong phishing-resistant assurance, and the standards mapping also describes privileged activation that way. | Implementations can either weaken the canonical/catalog claim or reject valid normal JIT requests that satisfy the normative appendix. Assurance policy and conformance cannot be independently proven from contradictory requirements. | Decide explicitly whether normal JIT requires phishing resistance, align F18-RD-14, the canonical model/catalog, matrix D09/D14, standards mapping and fixtures, and regenerate the exact approval payload. |
| F18-SEC-005 | Medium | Blocking | The readiness evidence is not actually reconciled. The architecture digest still gives `AuthorizationResult` four outcomes (`Allow | Deny | RequiresApproval | Indeterminate`) while F18-RD-02/F18-RD-09 and the canonical model close it to `Allow | Deny`; the digest describes broad FEATURE-0017 subject use while F18-RD-11 permits only one exact PrivilegedAccessRequest action and three scopes; and the industry matrix's explicit in-scope inventory omits `AccessGroup` even though conflict 27 and the leakage proof claim that omission was reconciled. | A downstream reviewer or generator can consume stale security semantics, and the final manifest's prepared claim that all 36 conflicts are reconciled is false. This violates the package's stop-on-conflict rule and invalidates it as final approval evidence. | Correct or retire every stale mirror, explicitly add AccessGroup to the matrix boundary, re-audit all 36 conflicts against exact F18-RD authority, and regenerate the manifest/digests before renewed independent review. |

All five findings are blocking. DEC-0060 and the FEATURE-0018 architecture gate
must remain pending.

## 5. Required attack and abuse analysis

| Attack / abuse case | Independent result |
|---|---|
| Confused deputy | **Fails review.** Exact authorization targets are strong, but delegated grant publication lacks resource-reach containment and the VS-000 writer contradiction permits a bypass path (F18-SEC-001/002). |
| Privilege escalation | **Fails review.** Resource A authority can be interpreted as sufficient to delegate resource B/scope-wide authority; ordinary standing access can be self-retained (F18-SEC-002/003). |
| Stale/replayed approval | Covered conceptually by immutable proposal digests, bounded expiry, current evidence at effect publication, fresh authorization, exact linkage and idempotency. Approval remains valid after an approver later loses eligibility by explicit admission-evidence design; the maximum 24-hour window is a residual risk. |
| Group-membership race | Current direct membership and group eligibility are rechecked for authorization and approval conflict at publication; suspension/revocation fail closed. Linearizable implementation and race proof remain mandatory. |
| Revoke/use race | Evaluation-instant, non-bearer results and current dependency checks support fail-closed behavior. Repository/effect-boundary linearization remains executable proof. |
| Suspension/use race | Membership/group/role suspension makes dependent grants immediately ineligible in the model; exact boundary race proof remains mandatory. |
| Activation-deadline race | Correctly requires the complete assignment/request/review atomic publication instant strictly before the deadline; at/after deadline denies without grace. |
| Approval self-dealing — direct | ApprovalRequest conflict handling is sound at decision and effect publication. AccessReview ordinary standing certification is not (F18-SEC-003). |
| Approval self-dealing — indirect | ApprovalRequest covers group members/owners, responsible parties and referenced beneficiaries. AccessReview does not apply the equivalent set (F18-SEC-003). |
| Audit failure | Required protected obligation acceptance blocks result/effect publication and completed idempotency results. This is fail closed but makes audit availability a control-plane availability dependency. |
| Idempotency collision | Actor, operation, exact target, scope and digest binding plus changed-content conflict and current replay checks are adequate at architecture level. |
| Cross-tenant/scope enumeration | Structural/local validation before safe resolution, safe 404/denial behavior and protected projections adequately prevent existence disclosure at contract level. |
| Evidence tampering/redaction | Immutable pins, protected FEATURE-0013 evidence, system-owned provenance and redaction are sound; cryptographic/storage protection remains downstream implementation/operations proof. |
| Compromised adapter identity | Not controlled by FEATURE-0018. A future adapter credential with direct native reach can bypass Sovrunn even though the intended execution path intersects both decisions. The future owner must prove least privilege, credential isolation/rotation, egress/call-path control, native log correlation and rapid revocation. |
| Native/Sovrunn policy disagreement | Correctly resolves by intersection: Allow/Allow is required; either Deny blocks. No layer may translate the other layer's Allow into authority. |
| Break-glass abuse | Individual eligibility, phishing-resistant recent assurance, one-hour maximum, exact scope, no renewal, immediate audit, atomic linked retrospective review and expiry are strong. Normal-JIT assurance ambiguity must still be corrected (F18-SEC-004). |

## 6. Residual risks after blocking corrections

These are not approval conditions because the current disposition is Reject;
they are risks a renewed review must re-evaluate:

- executable controls, failure injection and concurrency/race proof do not yet
  exist;
- identity lifecycle, authenticators and trusted group provisioning remain
  external and require a future adapter security assessment;
- native-IAM credentials and compromised-adapter containment remain future
  execution-owner responsibilities;
- existing access intentionally does not automatically fail solely because an
  AccessGroup owner or Workload/System responsible party becomes inactive;
- production Workload/System standing access may use escalation rather than
  automatic ineligibility on overdue review;
- publication-time delegation is intentionally non-cascading after later
  grantor revocation; and
- protected audit availability is a deliberate fail-closed dependency and a
  potential denial-of-service target.

## 7. Validation evidence

Read-only validation performed against the reviewed payload:

- manifest SHA-256 matched
  `c447e522937846ad490fc1f5c9a3198e2dd4eb129444cd699c8c436b526282d7`;
- manifest payload count was exactly 31 and every listed SHA-256 returned `OK`;
- F18-RD-01 through F18-RD-24 occurred exactly once across Appendices A/B;
- `.automation/features/FEATURE-0018.control.json` parsed successfully;
- `VS-000-contract-registry.yaml` parsed successfully with aliases enabled;
- `make arch-handoff-check` passed structural checks and correctly warned that
  the handoff is not Approved; and
- `git diff --check` passed.

An initial Ruby YAML helper invocation used an unavailable convenience method;
the validation was rerun with `YAML.safe_load(File.read(...), aliases: true)`
and passed. This tooling correction did not change repository content.

## 8. Reviewer disposition and signature

| Field | Value |
|---|---|
| Reviewer identity and role | Codex isolated security architecture reviewer |
| Independence basis | fresh fork with no conversation history; did not author the reviewed payload |
| Review date | 2026-08-29 |
| Exact commit/digests reviewed | Manifest SHA-256 `c447e522937846ad490fc1f5c9a3198e2dd4eb129444cd699c8c436b526282d7`; all 31 manifest-listed payload digests independently verified |
| Findings | F18-SEC-001 through F18-SEC-005; five blocking findings (three High, two Medium) |
| Required corrections | Reconcile sole RoleAssignment writer; close delegation resource reach; enforce access-review beneficiary independence; resolve JIT assurance; repair stale/false reconciliation evidence; regenerate manifest and obtain renewed independent review |
| Residual risk | High while blocking findings remain; Medium until executable, identity-adapter, native-IAM and race evidence exists after correction |
| Disposition | **Reject** |
| Signature/evidence reference | Signed 2026-08-29 by Codex isolated security architecture reviewer against manifest SHA-256 `c447e522937846ad490fc1f5c9a3198e2dd4eb129444cd699c8c436b526282d7` |

This rejection is security-review evidence only. It does not modify the
FEATURE-0018 architecture or authorize any downstream stage.
