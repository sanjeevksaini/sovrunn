---
doc_type: independent_security_review
feature: FEATURE-0018
renewal: 02
status: rejected-blocking-findings
reviewer: Codex independent security architecture reviewer renewal-02
review_date: 2026-08-30
updated: 2026-08-30
---

# FEATURE-0018 Independent Security Review Renewal 02

## 1. Disposition, integrity, and independence

**Disposition: REJECT — TWO NEW BLOCKING FINDINGS.**

This review was performed by **Codex independent security architecture reviewer
renewal-02** as a distinct non-author review task. The reviewer did not author,
approve, reconcile, or semantically modify the reviewed payload. Repository
bytes and read-only terminal validation were the source of truth; no discussion
history or prepared conclusion was accepted as architecture evidence. The
original review and renewal-01 were read only as immutable prior evidence and
were not modified.

Reviewed manifest:

```text
docs/reviews/architecture-readiness/FEATURE-0018-final-architecture-approval-manifest.md
SHA-256 7d2e70bb1ce367d5499d033f01969887a2cbfb518ffb71efc736e74812d8533d
```

The manifest contained exactly **31 payload files**. Its SHA-256 and every
listed payload digest independently matched before review. They were verified
again after this evidence record was written and remained unchanged. This
review record is intentionally excluded from the manifest table, so its
creation does not create a circular evidence dependency.

The regenerated payload closes F18-REN-001 and F18-REN-002 and, with them, all
five original F18-SEC findings. Those closures do not make the whole model
approvable: an AccessGroup membership-publication path can still turn delegated
membership administration into authority the actor does not hold, and a
manifest-pinned traceability mirror still marks superseded DEC-0058 as
Accepted. See F18-REN2-001 and F18-REN2-002.

## 2. Review boundary and method

The review covered the exact 31-file manifest payload: the ACR and proposed
DEC; the indivisible ADH-2026-070 core and Appendices A/B; the canonical model,
catalog, glossary, phase, context, sequence, roadmap, registry, and traceability
changes; FEATURE-0013 and FEATURE-0016 registrations; the FEATURE-0018
architecture and feature contracts; control manifest; reuse evidence; digest;
industry matrix; leakage proof; standards mapping; and threat/abuse ledger.

The assessment independently examined zero-trust authorization algebra,
identity and direct-group semantics, writer ownership, delegated reach,
approval, JIT and break-glass, access-review beneficiary expansion, exception
and governance boundaries, provider-native IAM intersection, validation order,
concurrency and idempotency, privacy and audit, simplicity, completeness,
architecture leakage, stale mirrors, and evidence consistency. It created no
requirements, design, tasks, implementation, or architecture semantics.

## 3. Independent security assessment

| Area | Renewal-02 determination |
|---|---|
| Zero trust and authorization algebra | Grant union across independently applicable assignments, mandatory-guardrail intersection, exact target binding, current dependencies, deny/indeterminate fail-closed behavior, and evaluation-instant non-bearer results are coherent. A separate grant-producing path remains outside the grantor ceiling: creating/reactivating direct AccessGroup Membership can activate an existing group RoleAssignment for a new beneficiary without action-reach dominance. See F18-REN2-001. |
| Scope, target, and delegated reach | `ExactResource`, `CreateParent`, and `ScopeOnly` form a closed discriminated registry. RoleAssignment delegation now requires `roleassignment.grant` over complete reach and one complete, independently applicable witness for each proposed action; resource A cannot delegate resource B or scope-wide authority. F18-SEC-002 remains closed. The same dominance is not applied when group membership activates that assignment for another principal. |
| Identity and group semantics | Stable issuer/subject identity, closed Human/Workload/System kinds, direct canonical membership, provenance/freshness, Guest expiry, group suspension/retirement, owner accountability, and raw-claim rejection are sound. However, the authorized administrator/trusted-provisioner membership intent is itself an effective access-enablement boundary when the target group already holds an assignment, and its authority is not bounded by those actions. See F18-REN2-001. |
| RoleAssignment sole writer | F18-RD-02/F18-RD-08, the compact digest, canonical catalog, VS0-WRITER-008, and traceability now consistently make the RoleAssignment controller the sole canonical spec/status/lifecycle/effect publisher. Grant, approval, privileged-access, and review controllers submit exact immutable intents only. F18-SEC-001 and the writer portion of F18-REN-001 are closed. |
| Approval and separation of duties | Approval stages, exact eligibility, quorum, bounded expiry, current conflict-set checks at decision and effect publication, and approval/effect separation are appropriate. Direct and indirect beneficiaries include group members/owners and responsible Humans. Group membership administration can nevertheless change who receives effective ordinary authority or becomes eligible through a role-holding group without the equivalent grantor ceiling. |
| JIT and break-glass | Normal JIT and break-glass share a non-weakenable, provider-neutral, phishing-resistant AAL2-or-higher activation floor with evidence age at most 15 minutes. Individual eligibility, activation deadline, no-grace TimeBound assignment, one-hour break-glass ceiling, immediate audit, and linked retrospective review are coherent. F18-SEC-004 remains closed. |
| Access review | F18-RD-16 now derives `ReviewBeneficiaryHumanSet` from a direct Human holder, a Workload/System holder's responsible Human, or an AccessGroup owner/direct Human members/responsible Humans of all direct Workload/System members. Current and proposed replacement sets are united, rechecked at atomic decision/effect publication, and unresolved expansion blocks Retain/Replace without blocking Revoke. F18-SEC-003 and F18-REN-002 are closed. |
| Exceptions and governance | Exception evidence is exact, immutable, time-bounded, non-overlapping, typed, and narrowing-only; it cannot override non-exceptionable identity, scope, authorization, or audit invariants and never grants a role. FEATURE-0020 retains effective exception/profile resolution. GovernanceProfile v1 stays within five FEATURE-0018 component families. |
| Provider-native IAM | Sovrunn authorization and native IAM are a true intersection: both must allow and either may deny. No native principal, policy, credential, translation, or effect enters FEATURE-0018. Compromise containment for a future adapter identity remains a named later integration risk. |
| Validation, concurrency, and idempotency | The ordered boundary performs structural/local validation before safe reference resolution, then current evidence and authorization before concurrency/idempotency, audit obligations, and publication. Exact replay binding, first-valid terminal transitions, no completed replay on audit failure, and separate atomic approval/activation/exception boundaries are architecture-level sound; executable linearization remains downstream proof. |
| Privacy and audit | Raw assertions, factors, claims, secrets, provider credentials, confidential policy inputs, inaccessible target detail, and evaluator diagnostics are excluded from unsafe errors/logs/projections. Protected audit-obligation acceptance precedes authorization results and mutations. Audit availability is therefore intentionally a fail-closed availability dependency. |
| Simplicity | Six progressive-disclosure journeys keep routine authorization and system-owned versions/provenance hidden without inventing a facade resource. The resource separation is proportionate at architecture stage. |
| Completeness, leakage, and evidence | F18-REN-001's stale FEATURE-0018 digest/matrix defects are corrected, and the original stale-phrase scan is clean in active non-historical files. The manifest-pinned Decision Traceability Matrix still publishes DEC-0058 and DEC-0059 as simultaneously Accepted despite DEC-0059 superseding DEC-0058. See F18-REN2-002. No real adapter, persistence, provider execution, FEATURE-0020 resolver, or later-feature implementation leaks into FEATURE-0018. |

## 4. Prior-finding closure status

| Prior finding | Renewal-02 status | Independent evidence |
|---|---|---|
| F18-SEC-001 — conflicting RoleAssignment writers | **CLOSED** | F18-RD-02/08, the canonical catalog, compact digest, VS0-WRITER-008, and VS-000 traceability consistently assign canonical RoleAssignment spec/status/lifecycle/effects to the RoleAssignment controller; other authorities submit immutable intents only. |
| F18-SEC-002 — delegation omitted resource reach | **CLOSED** | F18-RD-09 requires `roleassignment.grant` over the complete proposal reach and one independent per-action witness covering scope, target mode, resource reach, validity, and delegation; cross-assignment synthesis and resource-A-to-B/scope widening are prohibited. |
| F18-SEC-003 — access-review self-certification | **CLOSED** | F18-RD-16's exact `ReviewBeneficiaryHumanSet` excludes direct holders, group owner/direct Human members, direct Workload/System responsible Humans, and responsible Humans behind direct Workload/System group members at decision/effect publication. |
| F18-SEC-004 — inconsistent normal-JIT assurance | **CLOSED** | F18-RD-14/15 and active mirrors uniformly require current phishing-resistant AAL2-or-higher assurance no more than 15 minutes old for normal JIT and break-glass activation. |
| F18-SEC-005 — stale or falsely reconciled readiness evidence | **CLOSED for the enumerated FEATURE-0018 defects** | The digest is now a compact non-normative index; it correctly records sole writer, contextual Membership, all three target modes, embedded UsageEvidenceSummary, campaign/item vocabulary, informational supersession, and F18-SCN-01..45 with F18-SCN-38 retired. A distinct traceability inconsistency is recorded separately as F18-REN2-002. |
| F18-REN-001 — stale contradictory digest and scenario/writer mirrors | **CLOSED** | The regenerated digest no longer restates the obsolete 1,200-line contract, pins all exact semantics to the ADH, and carries aligned security sentinels and the complete scenario inventory. ADH core/appendix digests match the architecture pins. |
| F18-REN-002 — AccessGroup-to-Workload/System responsible-Human review gap | **CLOSED** | F18-RD-16 explicitly expands AccessGroup-held authority through every current direct Workload/System member to its current responsible Human, unions reviewed/replacement sets, fails closed on unresolved expansion, and preserves independently authorized Revoke. Matrix E02, F18-SCN-45, threat F18-THR-07, standards, canonical, feature, and VS mirrors agree. |

Result: **all seven prior findings are closed for their exact reported defects**.
The new findings below arise from independent reassessment of the complete
regenerated payload, not from carrying forward an earlier conclusion.

## 5. New blocking findings

### F18-REN2-001 — AccessGroup membership publication bypasses grantor reach

**Severity:** High

**Blocking:** Yes

**Classification:** Authorization bypass; delegated privilege escalation;
group-membership confused deputy; incomplete grant-producing-operation model

**Evidence:**

- F18-RD-04/F18-RD-05 make an active direct AccessGroup Membership immediately
  eligible to consume every otherwise-applicable RoleAssignment held by that
  group.
- F18-RD-06 defines `membership.write` as the Privileged action that can create
  trusted belonging/group intent at an Organization or CloudProvider parent.
  It does not bind that writer to the actions, exact resources, scopes,
  duration, or delegation capability carried by the target group's existing
  assignments.
- F18-RD-09's mandatory grantor ceiling is defined for a proposed
  RoleAssignment. It does not run when a Membership is created, reactivated,
  or refreshed into authorization eligibility.
- F18-RD-10's closed grant-producing operation derivation covers RoleAssignment
  grant/proposal, PrivilegedAccessRequest, and ExceptionProposal controllers;
  it gives an AccessGroup membership change no ApprovalRequirement or equivalent
  current-effect ceiling.
- F18-RD-20 correctly recognizes Membership publication as
  authorization-changing and requires audit, but audit does not constrain the
  authority being enabled.
- The industry matrix tests membership-only denial and RoleAssignment grantor
  reach separately, but contains no negative in which an actor with only a
  temporary `membership.write` grant adds themself or an accomplice to a group
  whose ordinary-role actions they do not hold.

**Security impact:** A principal may receive a bounded JIT assignment containing
`membership.write`, add themself to an AccessGroup that already holds durable
Ordinary authority outside the principal's current action/resource reach, and
retain that group-derived access after the JIT administrator assignment expires.
The same path can change role-based approver or reviewer eligibility. A
compromised trusted provisioner can likewise activate all current group-held
authority for a selected principal even though raw external claims were
designed not to be authorization. This launders temporary or narrowly delegated
membership administration into durable authority and bypasses the
RoleAssignment grantor ceiling.

**Recommendation:** Treat every AccessGroup Membership create, reactivate, and
freshness/provenance transition that makes a principal eligible as a
grant-changing publication boundary. Define one deterministic, atomic rule that
derives the complete current group-held assignment reach becoming available and
requires the actor/provisioner to possess a non-synthesized per-action ceiling
over that reach, or requires an exact approved group-membership policy with
equivalent separation of duties and bounded trusted-provisioner authority.
Re-evaluate current group assignments, membership source, owner, actor,
beneficiary, scope/resource reach, validity, and approval at publication;
changed or unresolved inputs must fail closed. Add self-add, accomplice-add,
JIT-expiry laundering, approver/reviewer-group, group-assignment race, and
trusted-provisioner overreach scenarios. This is an architecture correction
requiring renewed exact-payload approval; requirements must not invent it.

### F18-REN2-002 — Manifest-pinned traceability keeps superseded DEC-0058 Accepted

**Severity:** Medium

**Blocking:** Yes

**Classification:** Reconciliation integrity; stale authority mirror;
architecture leakage and evidence inconsistency

**Evidence:**

- `docs/traceability/DECISION_TRACEABILITY_MATRIX.md` marks
  `DEC-0058 Alpha Migration No Dual Authority` as `Accepted` and immediately
  marks DEC-0059 as `Accepted`.
- The higher-authority `docs/decisions/DECISION_INDEX.md` marks DEC-0058
  `Superseded by DEC-0059`.
- `docs/phase2/PHASE2R_REBASELINE.md` and the current architecture baseline
  likewise state that DEC-0059 replaced the signed-cutover/runtime-migration
  model with canonical bootstrap and no runtime alpha migration.
- The traceability matrix is one of the 31 manifest-pinned files, and the
  prepared Phase 2R drift result did not detect this contradictory status.

**Security/architecture impact:** A downstream reviewer, generator, or gate can
interpret both the signed-cutover migration authority and the no-migration
bootstrap authority as active. That can reintroduce a migration controller or
cutover-state design expressly prohibited by DEC-0059 and invalidates the exact
payload's claim that authoritative mirrors are reconciled.

**Recommendation:** Change the DEC-0058 traceability status to
`Superseded by DEC-0059`, remove or clearly historicalize obsolete validation
and feature ownership, and add an accepted/superseded-status consistency check
between the Decision Index, Phase 2R rebaseline, and traceability matrix.
Regenerate every affected digest and obtain another exact-manifest review.

## 6. Abuse-case disposition

| Attack or abuse case | Result |
|---|---|
| Confused deputy in RoleAssignment delegation | **Controlled.** Complete per-action witnesses and no cross-assignment synthesis close resource/scope/duration/delegation overreach. |
| Confused deputy in group membership administration | **Not controlled.** Membership publication can activate an existing group assignment without equivalent action-reach dominance or mandatory approval. F18-REN2-001. |
| Temporary-administrator privilege laundering | **Not controlled.** A JIT `membership.write` holder can potentially create durable group-derived Ordinary access that survives expiry of the administrative JIT assignment. F18-REN2-001. |
| Raw external group claim | Controlled as a direct input: it grants nothing. A trusted provisioner remains a high-impact boundary and currently lacks a group-effective-reach ceiling. |
| RoleAssignment multiple-writer bypass | Controlled by the sole RoleAssignment controller and exact immutable intents. |
| Direct and compounded-indirect review self-dealing | Controlled, including AccessGroup to direct Workload/System member to responsible Human and reviewed/replacement set union. |
| Stale/replayed approval | Controlled by exact proposal digest, bounded expiry, current publication checks, replay binding, and immutable provenance after effect publication. |
| Group membership/review race | Review conflict expansion is rechecked at publication. Membership activation against concurrent group assignments lacks the independent reach/approval rule identified in F18-REN2-001. |
| Revoke/use, suspend/use, and expiry races | Architecture requires current dependencies, authoritative time, non-bearer decisions, no grace, and atomic publication; executable linearization remains downstream proof. |
| Break-glass abuse | Individual eligibility, fresh phishing-resistant assurance, exact scope, one-hour maximum, no mutation renewal, immediate audit, independent expiry, and linked retrospective review are strong. |
| Exception enlargement | Exact typed narrowing, minimum/union/intersection composition, local non-overlap, bounded time, and non-exceptionable invariants control silent enlargement. |
| Audit failure | Fail closed before authoritative result/effect publication; deliberate availability/denial-of-service risk remains. |
| Cross-scope enumeration | Structural/local validation, safe resolution, inherited 404/denial behavior, redaction, and protected evidence are suitable at contract level. |
| Native-IAM disagreement | Controlled by two-plane intersection; neither Allow substitutes and either Deny blocks. |
| Compromised future adapter identity | Outside FEATURE-0018 execution. Later proof must cover least privilege, credential isolation/rotation, call-path and egress constraints, native-log correlation, and rapid revocation. |

## 7. Residual risks and non-blocking observations

The following remain material even after correcting the blocking findings:

- implementation, storage atomicity, linearizability, race, retry, and failure-
  injection evidence does not yet exist;
- real identity lifecycle, authenticators, trusted provisioning, external group
  synchronization, and native-adapter credential security require later
  independent assessment;
- existing effective access intentionally does not fail solely because an
  AccessGroup owner or Workload/System responsible Human later becomes inactive;
- production Workload/System Standing access may escalate rather than become
  automatically ineligible on overdue review under its exact registered rule;
- publication-time delegation intentionally does not cascade when the grantor
  later loses authority;
- approval evidence admitted within its at-most-24-hour window is not revoked
  solely by later approver ineligibility, although conflict and effect evidence
  are rechecked as specified;
- protected audit availability is a fail-closed control-plane dependency and a
  denial-of-service target;
- FEATURE-0016's combined `executiontarget.write` create/retire granularity is
  an explicit bounded dependency limitation; and
- the control manifest correctly remains pending and points to renewal-01 until
  a later controlled approval-metadata reconciliation; this rejection does not
  authorize that reconciliation.

## 8. Read-only validation evidence

The following checks were run against the exact payload. Except for this review
record, they made no repository change:

```text
shasum -a 256 docs/reviews/architecture-readiness/FEATURE-0018-final-architecture-approval-manifest.md
=> 7d2e70bb1ce367d5499d033f01969887a2cbfb518ffb71efc736e74812d8533d

manifest-table SHA-256 verifier
=> payload_count=31
=> failures=0 before review
=> failures=0 after review record creation

ruby -rjson JSON.parse(.automation/features/FEATURE-0018.control.json)
=> PASS

ruby -ryaml YAML.safe_load(VS-000-contract-registry.yaml, aliases: true)
=> PASS

F18-RD heading uniqueness across Appendices A/B
=> 24 headings; F18-RD-01 through F18-RD-24 exactly once; PASS

ADH core conflict ledger and leakage-proof reconciliation ledger
=> IDs 1 through 36 in exact order, each once, in both ledgers; PASS

ADH content pins
=> core f35357b03cfe0fdb9ccc3d14ed4a3330d7f367dd681ba12290fccd3a1d2865f7; PASS
=> Appendix A 1efa8768c1eaa04c43addbd19579adf55cb79ed91939ee7adbaf709fc8db8fd4; PASS
=> Appendix B 83f55f6a9b29e618e8cc7f0d68e6792b4aff396ea6040b37783f8ccf656aab24; PASS

active-file stale-phrase scans
=> no obsolete split RoleAssignment writer, usageEvidenceRef, four-outcome
   AuthorizationResult, missing ScopeOnly, F18-SCN-01..43 inventory, qualitative
   Standing narrowness, item-level ExpiredIncomplete, or Superseded lifecycle
   mirror found; PASS for the F18-REN-001 phrase set
=> DEC-0058 status contradiction found; F18-REN2-002

make arch-handoff-check HANDOFF=docs/reviews/architecture-decision-handoffs/ADH-2026-070-feature-0018-governance-iam-approval-exception-foundation.md
=> PASS with expected warning that the handoff is Proposed, not Approved

scripts/reuse-assessment-check.sh FEATURE-0018 --assessment docs/features/FEATURE-0018-governance-iam-approval-exception-foundation.md --mode strict --skip-rac03 --skip-rac13
=> PASS

make phase2-scope-check FEATURE=FEATURE-0018
=> PASS

make phase2r-drift-check
=> PASS (15 checks; the current script does not detect F18-REN2-002)

make feature-contract-check FEATURE=FEATURE-0018
=> NOT CONFIGURED; command reports no FEATURE-0018 contract profile; not counted as a pass

make structurizr-check
=> workspace-presence PASS; local CLI unavailable, syntax validation skipped;
   no Structurizr-level boundary or flow change is proposed by this review

mkdocs build --strict --site-dir <mktemp /tmp/f18-renewal02-mkdocs.*>
=> PASS; temporary output removed

git diff --check
=> PASS
```

An initial reuse-check invocation accidentally targeted the separate approval-
evidence record, which is not the feature assessment and therefore failed the
assessment schema. It was rerun against the correct manifest-pinned feature
assessment and passed. The correction changed no file.

Go format, test, vet, race, runtime, and the final feature gate are not
applicable at this architecture-only stage: no requirements, design, tasks, or
implementation are authorized.

## 9. Reviewer signature

| Field | Value |
|---|---|
| Reviewer identity | Codex independent security architecture reviewer renewal-02 |
| Independence basis | Distinct non-author review task; did not author, reconcile, or approve payload; repository bytes and terminal evidence only; prior reviews treated as immutable evidence |
| Review date | 2026-08-30 |
| Exact manifest reviewed | SHA-256 `7d2e70bb1ce367d5499d033f01969887a2cbfb518ffb71efc736e74812d8533d` |
| Payload verification | All 31 manifest-listed SHA-256 digests matched before review and after evidence creation |
| Prior-finding closure | F18-SEC-001 through F18-SEC-005 and F18-REN-001/F18-REN-002 closed for their exact reported defects |
| New findings | F18-REN2-001 High blocking; F18-REN2-002 Medium blocking |
| Residual risk | High while membership-effective-reach and traceability-authority defects remain; Medium until executable and future identity/native-adapter proof exists after correction |
| Disposition | **Reject** |
| Signature/evidence reference | Signed 2026-08-30 against manifest SHA-256 `7d2e70bb1ce367d5499d033f01969887a2cbfb518ffb71efc736e74812d8533d` |

This rejection is independent security-review evidence only. It does not alter
the reviewed architecture, approve reuse, promote DEC-0060/ACR-2026-002,
authorize approval-metadata changes, or authorize requirements, design, tasks,
implementation, or a downstream feature gate.
