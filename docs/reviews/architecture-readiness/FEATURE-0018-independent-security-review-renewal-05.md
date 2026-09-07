---
doc_type: independent_security_review
feature: FEATURE-0018
renewal: 05
status: passed-no-blocking-findings
reviewer: Codex fresh non-author security-review subagent renewal-05
review_date: 2026-08-30
updated: 2026-08-30
---

# FEATURE-0018 Independent Security Review Renewal 05

## 1. Disposition, integrity, and independence

**Disposition: PASS — NO BLOCKING SECURITY FINDINGS IN THE EXACT REVIEWED
ARCHITECTURE PAYLOAD.**

This review was performed by **Codex fresh non-author security-review subagent
renewal-05** as a distinct delegated review task. The reviewer did not author,
reconcile, approve, or semantically modify the reviewed architecture payload.
The reviewer created only this evidence record. Repository bytes and read-only
terminal verification were treated as the source of truth. The original
security review and renewals 01 through 04 were treated as immutable rejection
evidence and were not modified.

Reviewed manifest:

```text
docs/reviews/architecture-readiness/FEATURE-0018-final-architecture-approval-manifest.md
SHA-256 9f69ea2fb89ba04aa253a37a0cf5f63124753c2b1e9cf08cf0a5c028369d8efc
```

The manifest contains exactly **32 payload files**. Its SHA-256 and every
listed payload digest independently matched before review. They were verified
again after this record was written and remained unchanged. This review record
is intentionally excluded from the manifest table, so its creation does not
create a circular evidence dependency.

The regenerated payload closes F18-REN4-001 by making every approval,
privileged-requester, and reviewer `EligibilityRef` exactly one current Human
PrincipalRef. RoleDefinition, RoleAssignment, AccessGroupRef, Membership,
claims, and every other relationship are invalid eligibility sources. The
policy or rule remains the explicit eligibility authority; its applicability
and the actor's separately authorized canonical action independently determine
scope and target reach. The architecture no longer has the role-assignment or
later-role-reference authority-laundering path identified by renewal-04.

This pass is independent security-review evidence for the exact manifest only.
It does not approve ACR-2026-002 or DEC-0060, promote any decision status,
record human approval, approve reuse, or authorize requirements, design, tasks,
or implementation.

## 2. Exact review boundary and method

The review covered the complete 32-file manifest payload: ACR-2026-002 and
DEC-0060; the indivisible ADH-2026-070 core and Appendices A/B; canonical model,
catalog, glossary, context, sequence, roadmap, registry, writer, and
traceability changes; FEATURE-0013 and FEATURE-0016 dependency registrations;
the FEATURE-0018 architecture and feature contracts; control manifest; reuse
evidence; digest; industry matrix; leakage proof; standards mapping;
threat/abuse ledger; and the pinned Phase 2R drift checker.

The method was:

1. independently verify the manifest SHA-256 and all 32 payload digests;
2. read the original review and renewals 01 through 04 and map every prior
   High and Medium finding to current normative and mirrored evidence;
3. trace every admitted `EligibilityRef` backward through all state that can
   establish current satisfaction and forward into approval, review, normal
   JIT, and break-glass effects;
4. retest role assignment, later policy/rule publication, AccessGroup
   Membership, trusted provisioning, concurrency, scope, target, action, and
   separation-of-duties eligibility-laundering paths;
5. retest direct and membership-enabled delegated authority for exact action,
   scope, target mode, resource reach, validity interval, and delegation reach,
   including the second-publisher ordering;
6. reassess RoleAssignment writer ownership, action target bindings, approval
   and effect atomicity, access-review beneficiary expansion, native-IAM
   intersection, validation order, idempotency, protected audit-before-
   publication, safe denial, user simplicity, and future-feature leakage; and
7. run safe read-only repository validation and reverify payload integrity.

The review generated no requirements, design, tasks, implementation, new
architecture semantics, or change to a manifest-pinned file.

## 3. Renewal-04 blocker retest

### 3.1 F18-REN4-001 — RoleDefinition-based EligibilityRef

**Retest result: PASS; CLOSED.**

The exact payload has one closed eligibility type and one satisfaction rule:

```text
EligibilityRef = exactly one Human PrincipalRef
satisfied       = exact current Human principal equality
```

F18-RD-02 closes the embedded-value schema. F18-RD-12, F18-RD-14, and
F18-RD-16 apply it uniformly to approvers, privileged requesters, and
reviewers. F18-RD-18 permits only non-empty Human PrincipalRef lists in the
three rule families. The glossary, canonical model/catalog, VS-000 registry,
architecture, digest, feature contract, leakage ledger, standards mapping,
threat ledger, and traceability mirrors agree.

The correction closes both exploit orderings from renewal-04:

- assigning an ordinary role to oneself or an accomplice cannot change any
  approval, review, JIT, or break-glass eligible set; and
- publishing a policy or rule after a role has acquired many holders cannot
  qualify any of those holders through the role.

`RoleDefinition`, `RoleAssignment`, `AccessGroupRef`, `Membership`, external
claims, nested groups, dynamic groups, and other relationships fail structural
eligibility validation. F18-SCN-49 and F18-RD-22 require self/accomplice,
already-broad-role, later-rule, unsupported-kind, scope/target/action, and
concurrent publication negatives.

The remaining exact Human eligibility is not action authority. The selected
policy or rule provides bounded applicability, and the actor must independently
possess the exact registered canonical action at the requested scope/target.
Eligibility and action are re-evaluated at the authoritative decision and
effect/remediation boundaries. Approval and access review additionally apply
their complete beneficiary-Human separation-of-duties sets.

## 4. Prior-finding closure matrix

| Prior finding | Renewal-05 status | Independent evidence |
|---|---|---|
| F18-SEC-001 — conflicting RoleAssignment writers | **CLOSED** | F18-RD-02/08, VS0-WRITER-008, canonical model/catalog, traceability, digest, and leakage proof make the RoleAssignment controller the sole specification, lifecycle, status, revocation, replacement, expiry, and review-effect publisher. Other controllers submit exact immutable intents only. |
| F18-SEC-002 — delegated grant omitted resource reach | **CLOSED** | F18-RD-09 requires independent per-action witnesses covering exact scope, target mode, resource reach, validity, and delegation. ExactResource A cannot authorize B or scope-wide reach; restrictions cannot be synthesized across assignments. |
| F18-SEC-003 — ordinary Standing access self-certification | **CLOSED** | F18-RD-16 expands the exact current/proposed beneficiary-Human set across direct holders, AccessGroup owners and direct members, and responsible Humans for Workload/System principals; it rechecks at decision/effect publication and blocks Retain, Replace, and due advancement on conflict or unresolved expansion. |
| F18-SEC-004 — normal-JIT phishing-resistance conflict | **CLOSED** | F18-RD-14/15 and active mirrors uniformly require current phishing-resistant AAL2-or-higher evidence no more than fifteen minutes old for JIT and break-glass activation. |
| F18-SEC-005 — stale reconciliation evidence | **CLOSED** | The compact digest and active mirrors agree on the sole writer, Membership behavior, all three target modes, review evidence/state, lifecycle, and the 49-scenario inventory. |
| F18-REN-001 — stale contradictory digest/mirrors | **CLOSED** | The digest is a compact non-normative index pinned to the exact ADH core and appendices; reported stale semantics are absent from active authority. |
| F18-REN-002 — AccessGroup/Workload responsible-Human review gap | **CLOSED** | Beneficiary expansion includes the accountable Human behind each current direct Workload/System member of a beneficiary AccessGroup and fails closed for authority-preserving decisions. |
| F18-REN2-001 — AccessGroup Membership bypasses grantor ceiling | **CLOSED END-TO-END** | Every assignment-effect expansion derives a coherent `MembershipEnabledGrantEnvelope`; non-empty envelopes require Membership authority, `roleassignment.grant`, and non-synthesized per-action scope/target/resource/validity/delegation witnesses. Empty groups and Membership cannot manufacture workflow eligibility. A later group assignment independently passes the ordinary ceiling. |
| F18-REN2-002 — DEC-0058 remains Accepted | **CLOSED** | Decision Index, Decision Traceability Matrix, and Phase 2R rebaseline identify DEC-0058 as `Superseded by DEC-0059`; drift CHECK 16 passes. |
| F18-REN3-001 — AccessGroup Membership manufactures eligibility | **CLOSED** | AccessGroupRef, Membership, group-held assignments, group-derived assignments, trusted provisioning, and raw claims are invalid eligibility sources. Membership changes only independently assigned group authority through the exact envelope. |
| F18-REN3-002 — reviewer eligibility undefined | **CLOSED** | F18-RD-16/18 define a non-empty Human PrincipalRef OR list, exact equality, rule applicability, separate `accessreview.decide`, fail-closed invalid evidence, decision/remediation rechecks, and post-eligibility beneficiary conflict. |
| F18-REN4-001 — RoleDefinition eligibility bypasses delegation | **CLOSED** | All eligibility uses are Human PrincipalRef-only. Role grant and later policy/rule reference cannot qualify role holders; F18-SCN-49 and F18-RD-22 require both-ordering and race negatives. |

Result: **every prior High and Medium finding is closed for its exact reported
defect, and regression review found no new blocking security finding.**

## 5. Whole-model security assessment

| Area | Renewal-05 determination |
|---|---|
| Authorization algebra | **PASS.** Independently applicable grant union is constrained by mandatory guardrail intersection, exact action/target/scope/resource binding, current dependencies, Deny/Indeterminate fail closure, and non-bearer results. Policy, approval, exception, Membership, controller identity, and native IAM do not create Sovrunn action authority. |
| Delegated authority | **PASS.** Direct grants and Membership-enabled group effects require complete, non-synthesized action/scope/target/resource/validity/delegation witnesses. TimeBound authority cannot enable a longer or Standing effect. Later grantor loss intentionally requires explicit revoke/policy/review response rather than hidden cascading revocation. |
| Human workflow eligibility | **PASS.** The exact Human PrincipalRef-only schema removes relationship-derived and second-publisher eligibility channels. Policy/rule applicability plus separate action authorization determines scope and target; rechecks and SoD occur at authoritative boundaries. |
| Identity and Membership | **PASS.** Issuer/subject identity, Human/Workload/System kinds, provenance/freshness, Guest expiry, direct groups, responsibility, contextual Membership, and raw-claim rejection are coherent. Membership itself grants no action or workflow eligibility. |
| RoleAssignment ownership | **PASS.** One canonical controller publishes every assignment specification and lifecycle/status/effect. Grant, approval, privileged-access, and review controllers submit only exact immutable intents and cannot enlarge or publish them. |
| Approval and separation of duties | **PASS.** Exact policy pinning, bounded stages/quorum/expiry, one effective decision per principal, current eligibility, beneficiary conflict expansion, and separate approval/effect atomic boundaries prevent approval from becoming bearer or write authority. |
| JIT and break-glass | **PASS.** Individual Human holder, exact rule/role/scope, separate submit authority, current eligibility, fresh phishing-resistant assurance, activation deadline, TimeBound assignment, one-hour break-glass maximum, audit, and retrospective review are coherent. Bypass authority exists only in the exact privileged-access rule and never overrides authentication, authorization, expiry, or audit. |
| Access review | **PASS.** Exact reviewer eligibility, separate action, immutable snapshots and usage evidence, complete direct/indirect beneficiary conflict expansion, exact remediation, and exact-version Standing certification prevent self-certification and stale remediation. |
| Exceptions and governance | **PASS.** Exception evidence is typed, narrowing-only, exact-scope, time-bounded, non-overlapping, immutable, and cannot grant a role or override non-exceptionable controls. GovernanceProfile activates only five FEATURE-0018 component families and performs no FEATURE-0020 resolution. |
| Target bindings | **PASS.** Every canonical action has explicit discriminated `ExactResource`, `CreateParent`, or `ScopeOnly` registrations. Missing, wildcard, mixed, unknown, or ambiguous registrations fail closed. FEATURE-0016 remains semantic owner of its three ExecutionTarget actions. |
| Validation, concurrency, and idempotency | **PASS at architecture stage.** Structural and local validation precede safe resolution; current authority/evidence precedes concurrency/idempotency and protected evidence acceptance. Replays bind actor, operation, target, scope, and digest, recheck current authorization, and publish no completed replay on audit failure. Executable linearization proof remains downstream. |
| Evidence integrity and privacy | **PASS.** Required DecisionRecord/AuditEvent obligations precede authoritative result or effect publication; exact versions and provenance are retained; raw tokens, claims, secrets, confidential policy inputs, and inaccessible target detail stay out of unsafe projections. |
| Native-IAM boundary | **PASS.** A Sovrunn Allow is necessary but insufficient for an ExecutionTarget effect. The separately owned adapter/controller identity must also pass native IAM. Either layer may deny and neither may override the other; provider-native IAM objects and credentials remain outside FEATURE-0018. |
| Simplicity | **PASS at architecture stage.** Routine authorization is silent; six risk-proportional journeys hide schema, versions, evidence, and controllers. Human-only eligibility removes a second entitlement system while preserving exact auditability. |
| Architecture leakage | **PASS.** No IdP/token validation, SCIM adapter, external/nested/dynamic-group authorization, arbitrary ABAC/workflow language, real policy/workflow engine, native IAM translation/credential/effect, persistence, FEATURE-0020 resolver, later-feature implementation, requirements, design, or tasks are introduced. |

## 6. Adversarial regression summary

| Attack or failure path | Renewal-05 result |
|---|---|
| Role grant creates approver/reviewer/JIT/break-glass eligibility | Rejected structurally; roles and assignments are invalid `EligibilityRef` kinds. |
| Later policy/rule references broadly assigned role | Impossible in the closed schema; a rule lists exact Humans only. |
| AccessGroup member addition activates durable authority | Non-empty assignment-effect envelope must pass complete temporal and reach dominance; a TimeBound administrator cannot enable a Standing effect. |
| Empty AccessGroup manufactures workflow eligibility | Impossible; Membership and all group relationships are invalid eligibility sources. |
| Trusted provisioner bypasses Sovrunn authority | Rejected; provisioner is a canonical System principal subject to the same operation and envelope checks. |
| Resource-A grant delegates resource B or scope-wide access | Rejected by exact per-action resource-reach containment and no witness synthesis. |
| Workflow controller publishes RoleAssignment directly | Rejected by sole-writer contract; it can submit only an exact immutable intent. |
| Reviewer certifies direct or indirect benefit | Retain/Replace/due advancement fail after complete current beneficiary-Human expansion; independently authorized Revoke remains available. |
| Stale or replayed approval performs a changed effect | Rejected by immutable subject/policy/digest linkage, bounded evidence, fresh effect authorization, exact idempotency, and separate atomic publication. |
| Break-glass bypasses authentication, action authority, expiry, or audit | Rejected; only approval is bypassable and only by the exact prepublished rule. |
| Audit service failure leaves authorization-changing state | Rejected; required obligation failure publishes neither authoritative state nor completed replay. |
| Native cloud Allow substitutes for Sovrunn Allow | Rejected by the two-plane intersection boundary. |
| Sovrunn Allow overrides native Deny | Rejected by the two-plane intersection boundary. |

## 7. Residual risks and non-blocking observations

No architecture-stage blocking finding remains. The following residual risks
must remain visible to final approval and later implementation review:

- exact-Human eligibility lists intentionally make policy/rule publication a
  high-impact governance authority; assignments for the corresponding publish
  actions must remain tightly scoped, TimeBound, audited, and independently
  reviewed under the approved model;
- implementation, persistence, storage atomicity, linearizability, race,
  retry, and failure-injection evidence does not yet exist;
- real identity lifecycle, authenticator, trusted-provisioning, external-group
  synchronization, and native-adapter credential boundaries require later
  independent review by their owning features;
- existing effective access intentionally does not fail solely because an
  AccessGroup owner or Workload/System responsible Human later becomes
  inactive; protected escalation and explicit suspension/revocation remain
  required;
- production Workload/System Standing access may escalate rather than become
  automatically ineligible when its exact registered overdue rule permits it;
- publication-time delegation intentionally does not cascade after later
  grantor authority loss; explicit assignment revocation, policy suspension,
  or review is required;
- protected audit availability is a fail-closed control-plane dependency and a
  denial-of-service target;
- FEATURE-0016's combined `executiontarget.write` create/retire granularity is
  an explicit bounded dependency limitation; and
- FEATURE-0020 must later prove correct effective profile selection without
  changing the FEATURE-0018-local rule validation semantics reviewed here.

These are explicit operational, dependency, or implementation-stage risks;
none reopens a current architecture contradiction or hidden authorization path
in the exact payload.

## 8. Read-only validation evidence

The following checks were run against the exact payload. Except for creation of
this review record, they made no repository change:

```text
manifest SHA-256
=> 9f69ea2fb89ba04aa253a37a0cf5f63124753c2b1e9cf08cf0a5c028369d8efc

manifest-table SHA-256 verifier
=> rows=32 matched=32 mismatched=0 before review
=> rows=32 matched=32 mismatched=0 after review record creation

immutable prior-review SHA-256 values
=> original   bae9a3dacadbad7a59434f215b43be0639f0490ec65782f3548faea5f962f701
=> renewal-01 1d550a2814a8804e0ca8912b53e6dd4a3a0432c89ff7f3353daccbb989a01859
=> renewal-02 e0ce7ec14d3876e04fafaa82a97301a7e4ad6e25bfdd34ec33cce88e66e4afb8
=> renewal-03 0f16272e89b44bd38011c2dfe75a59cf7bf575e02c72ae4018713bc3f60fb5c0
=> renewal-04 7502474ea74617233aca6d0da6402b8aef7b98416a384afeeee1aba99a463a92

control-manifest JSON parse and VS-000 registry YAML safe-load
=> PASS

F18-RD heading uniqueness across Appendices A/B
=> 24 headings; F18-RD-01 through F18-RD-24 exactly once; PASS

industry scenario inventory
=> 49 unique identifiers; F18-SCN-01 through F18-SCN-49; PASS

architecture handoff structural check
=> PASS with expected warning that the handoff is not marked Approved

strict reuse-assessment structure check
=> PASS; approval remains a separate human gate

make phase2-scope-check FEATURE=FEATURE-0018
=> PASS

make phase2r-drift-check
=> PASS; all 16 checks

make structurizr-check
=> workspace present; local CLI unavailable, syntax validation skipped

git diff --check
=> PASS
```

`mkdocs build --strict` was not rerun because it writes generated output and
this independent task was authorized to create only this review record. The
manifest's prepared build result was treated as a prepared claim, not
independent proof.

The review record's own final SHA-256 is reported with the reviewer handoff
rather than embedded here, because embedding its digest would be self-
referential.

## 9. Reviewer disposition and signature

| Field | Value |
|---|---|
| Reviewer identity and role | Codex fresh non-author security-review subagent renewal-05 |
| Independence basis | Distinct delegated review task; did not author, reconcile, approve, or semantically modify the reviewed payload |
| Review date | 2026-08-30 |
| Exact payload reviewed | Manifest SHA-256 `9f69ea2fb89ba04aa253a37a0cf5f63124753c2b1e9cf08cf0a5c028369d8efc`; all 32 listed digests verified before and after review |
| Renewal-04 result | F18-REN4-001 closed by exact Human PrincipalRef-only eligibility; every earlier High/Medium finding remains closed for its exact defect |
| New blocking findings | None |
| Residual risk | Medium until executable identity, concurrency, atomicity, audit-failure, and native-IAM boundary evidence passes later implementation and integration review |
| Disposition | **Pass — no blocking security findings in the exact reviewed architecture payload** |
| Signature/evidence reference | Signed 2026-08-30 by Codex fresh non-author security-review subagent renewal-05 against manifest SHA-256 `9f69ea2fb89ba04aa253a37a0cf5f63124753c2b1e9cf08cf0a5c028369d8efc` |

This pass is security-review evidence only. It does not change FEATURE-0018
architecture, approve ACR-2026-002 or DEC-0060, authorize requirements/design/
tasks/implementation, or permit modification of the reviewed manifest payload
without manifest regeneration and renewed review.
