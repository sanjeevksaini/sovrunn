---
doc_type: architecture_readiness_evidence
feature: FEATURE-0018
status: architecture-approved
decision: DEC-0060-accepted
updated: 2026-08-30
---

# FEATURE-0018 Architecture Leakage Proof and Reconciliation Ledgers

## 1. Evidence status and authority

This is an evidence mirror, not a semantic authority. F18-RD-01 through
F18-RD-24 in the exact ADH-2026-070 package control every value. `RECONCILED`
below means the approved repository application contains an aligned reference
or transcription; it does not mean executable implementation proof exists or a
later Feature Factory gate passed.

Overall architecture-leakage result: **PASS — RENEWAL-05 FOUND NO BLOCKING
FINDINGS; FINAL ARCHITECTURE APPROVAL RECORDED**.

## 2. Contract and profile ledger

| Contract | Profile | Independent identity/lifecycle | Scope summary | Sole semantic owner |
|---|---|---|---|---|
| PrincipalRef | EmbeddedValue | No | None | FEATURE-0018 reference semantics; external identity lifecycle |
| AccessGroup | ManagedResource | Yes | Organization, CloudProvider | FEATURE-0018 |
| Membership | ManagedResource | Yes | Organization, CloudProvider | FEATURE-0018 |
| RoleDefinition | VersionedDefinition | Yes | Platform, Organization, CloudProvider | FEATURE-0018 |
| RoleAssignment | ManagedResource | Yes | Seven canonical ScopeKinds | FEATURE-0018 |
| PrivilegedAccessRequest | LongRunningOperation | Yes | Seven canonical ScopeKinds | FEATURE-0018 |
| AccessReview | LongRunningOperation | Yes | Seven canonical ScopeKinds | FEATURE-0018 |
| ApprovalPolicy | VersionedDefinition | Yes | Platform, Organization, CloudProvider | FEATURE-0018 |
| ApprovalRequest | LongRunningOperation | Yes | Seven canonical ScopeKinds | FEATURE-0018 |
| ExceptionGrant | ImmutableRecord | Yes | Seven canonical ScopeKinds | FEATURE-0018 evidence; FEATURE-0020 effective application |
| GovernanceProfile | VersionedDefinition | Yes | Platform, Organization, CloudProvider | FEATURE-0018 v1 envelope/components; FEATURE-0020 resolution |

Closed non-resource values: RoleHolderRef, AssuranceEvidence,
AuthorizationInput, AuthorizationResult, ApprovalRequirement,
EligibilityRef, RoleAssignmentProposal, MembershipEnabledGrantEnvelope, ExceptionProposal,
UsageEvidenceSummary, GovernanceApplicability and ActionTargetBinding. They have no endpoint,
independent persistence/lifecycle or separate user journey.

### 2.1 Exact ledger index

To prevent an evidence copy from becoming a second authority, exact value
tables stay in whole F18-RD groups. This index proves every required ledger has
one controlling location and one review mirror.

| Required ledger | Exact controlling payload | Review mirror / reconciliation | Status |
|---|---|---|---|
| Resource/supporting inventory and profiles | F18-RD-02 Appendix B | Section 2; canonical catalog; VS-000 schemas | RECONCILED |
| Allowed scopes/reference compatibility/containment | F18-RD-02 Appendix B | Canonical catalog scope paragraph; matrix C06 | RECONCILED |
| Spec/status/result/field writer ownership | F18-RD-02 Appendix B plus resource RDs | Section 3; canonical catalog | RECONCILED |
| Operation/mutability/lifecycle/effective projection | F18-RD-02 Appendix B; F18-RD-04/05/07/08/12..18 Appendix A/B | Canonical catalog; matrix B..F | RECONCILED |
| Canonical actions/owner/classification/target binding | F18-RD-06 Appendix B | FEATURE-0016 additive table; matrix C01/C11..C14 | RECONCILED |
| Role classification/assignment applicability | F18-RD-07..09 Appendix A | Canonical model/catalog; matrix C | RECONCILED |
| Approval/privilege state and evidence flow | F18-RD-10/12..15 Appendix A | Matrix D; digest journeys | RECONCILED |
| Review and exception state/evidence flow | F18-RD-16/17 Appendix A | Matrix E; digest journeys | RECONCILED |
| Governance component applicability | F18-RD-18 Appendix B | Canonical catalog; matrix F | RECONCILED |
| Validation/concurrency/idempotency/terminal-retry | F18-RD-17/20/21 Appendix A/B | Digest F18-RD-21; matrix G/H | RECONCILED |
| Decision profiles and AuditEvent taxonomy | F18-RD-19/20 Appendix B | FEATURE-0013 adopting registration; matrix G | RECONCILED |
| FEATURE-0017 adoption | F18-RD-10/11 Appendix A | Feature/architecture dependency boundary; matrix G06..G09 | RECONCILED |
| Six progressive journeys/projections | F18-RD-24 Appendix A | Digest journey table; architecture section 12 | RECONCILED |
| Local conformance/deterministic vectors | F18-RD-22 Appendix B | Industry F18-IVM/F18-SCN inventory | ARCHITECTURE INVENTORY COMPLETE; EXECUTION PENDING |
| Standards/control/shared responsibility | F18-RD-23 Appendix B | FEATURE-0018 standards mapping | ARCHITECTURE EVIDENCE REVIEWED |
| Independent security review | F18-RD-23 Appendix B | Original and renewal-01 through renewal-04 rejection records retained. Renewal-05 verified F18-REN4-001 closed, regression-tested every prior finding and found no new blocker. | PASSED_BY_RENEWAL-05 |

## 3. Writer and materialization ledger

| Effect | Sole canonical publisher/status writer | Permitted upstream input | Forbidden competing writer |
|---|---|---|---|
| AccessGroup state | AccessGroup controller | Authorized administrator intent | Provisioner, assignment or approval controller |
| Membership state | Membership controller | Exact administrator or canonical System-provisioner intent; non-empty AccessGroup assignment-effect expansion also supplies the complete current grant-envelope authority evidence | Raw IdP claim, source-system trust, adapter, administrator, provisioner or authorization evaluator writing canonical state |
| Definition publication/lifecycle | Corresponding publication controller | Authorized draft/publish/lifecycle intent | Caller-authored status or scheduler |
| RoleAssignment publication/status | RoleAssignment controller | Exact accepted grant, privileged or review materialization intent | Grant, privileged, approval or review workflow controller |
| PrivilegedAccessRequest state | Privileged-access controller | Accepted requester/admin intent and current evidence | Approval controller or assignment controller |
| ApprovalRequest state/result | Approval controller | Exact request materialized by originating controller | Requester after acceptance or downstream effect controller |
| AccessReview state/result | Review controller | Authorized campaign or mandatory pinned-rule intent | Assignment controller |
| ExceptionGrant/linked revocation evidence | Exception controller | Exact approved proposal/revocation intent | Approval controller or FEATURE-0020 resolver |

Approval publication, assignment publication, privileged activation and
exception publication remain separate F18-RD-20 atomic boundaries.

## 4. Authorization and dependency ledger

| Boundary | Closed rule | Leakage sentinel |
|---|---|---|
| Grant composition | Union of independently applicable grants | No cross-assignment synthesis |
| Guardrails | Intersection; any deny/missing/ambiguous dependency denies | No policy/approval/exception-created permission |
| Target binding | ExactResource, CreateParent or ScopeOnly discriminated registration | No wildcard, mixed or inferred binding |
| AccessGroup | Current direct canonical membership only | No external claim, nesting or dynamic expression |
| AccessGroup membership expansion | Exact current group-held action/reach/validity envelope; Membership authority plus `roleassignment.grant` and one non-synthesized per-action witness; Standing effect requires Standing witnesses | No JIT-to-Standing laundering, trusted-provisioner bypass or stale membership/assignment snapshot |
| Human workflow eligibility | Non-empty OR list of exact Human PrincipalRefs only; separate policy/rule applicability and action authority determine scope/target reach; eligibility is rechecked at decision/effect publication | No RoleDefinition, RoleAssignment, AccessGroupRef, Membership, claim, other relationship, empty-list, inactive/stale Human, scope/target mismatch or manufactured eligibility |
| Validity | Standing or TimeBound; exact no-grace boundaries | No implicit mode or mutable renewal |
| Delegation | Publication-time grantor ceiling with one independently applicable scope/target/resource/time/delegation witness per action; no dependent v1 chains | No cross-assignment synthesis or grant beyond current exact reach |
| Review independence | Direct Human holder, direct Workload/System responsible Human, and, for AccessGroup-held authority, the owner, direct Human members and responsible Humans of direct Workload/System members are excluded at atomic disposition/effect publication; unresolved expansion fails closed for Retain/Replace without blocking Revoke | No direct, indirect or group-workload alias self-certification |
| Evidence | Exact versions/provenance; authorization result is non-bearer | No stale request/result replay as authority |
| Native effect | Sovrunn Allow intersects independent native-IAM Allow | No provider-native type/credential/effect in FEATURE-0018 |
| Effective governance | FEATURE-0020 only | No hierarchy, resolution or exception application in FEATURE-0018 |

## 5. Thirty-six-conflict reconciliation ledger

| ID | Reconciliation evidence in proposed application | Status |
|---:|---|---|
| 1 | Canonical model/catalog, glossary, sequence and VS-000 now include AccessGroup, holder/supporting values and direct group lifecycle boundary. | RECONCILED |
| 2 | Membership now supports Organization, CloudProvider or direct same-scope AccessGroup context without making AccessGroup a ScopeKind. | RECONCILED |
| 3 | RoleAssignment uses RoleHolderRef, exact role version, registered scope/optional resource and required dependency pins by F18-RD-08. | RECONCILED |
| 4 | Stale principal/condition/targetScope fields were removed; exact action-target and assignment restrictions point to F18-RD-02/06/08/09. | RECONCILED |
| 5 | Standing and TimeBound plus exact no-grace TimeBound boundaries are canonical/glossary terms. | RECONCILED |
| 6 | RoleDefinition now separates supersession, emergency suspension/restoration and terminal retirement. | RECONCILED |
| 7 | Membership provenance/freshness, Guest expiry, System/Workload responsibility and non-authorizing behavior are recorded; executable proof remains downstream. | RECONCILED |
| 8 | GovernanceProfile v1 is limited to five FEATURE-0018 component families; FEATURE-0020 resolution and future-domain ownership are explicit. | RECONCILED |
| 9 | ApprovalPolicy, ApprovalRequest and ExceptionGrant profiles, writers/lifecycles and boundaries are in canonical catalog/feature authority. | RECONCILED |
| 10 | FEATURE-0013 architecture records exactly three proposed adopting profiles and explicitly excludes an AccessReview profile. | RECONCILED |
| 11 | FEATURE-0017 is referenced unchanged through F18-RD-10/11; no F17 contract was modified. | RECONCILED |
| 12 | VS-000 stale Service principal, membership grant refs, missing scopes and broad GovernanceProfile fields were corrected; future semantics remain excluded. | RECONCILED |
| 13 | Canonical catalog records exact initial scope sets, generic versioned extensibility and core-change boundary for new scope/containment. | RECONCILED |
| 14 | Canonical/feature authorities record registered action ownership, monotonic privilege, fail-safe ambiguity and provenance. | RECONCILED |
| 15 | ApprovalRequirement is operation-local; privilegedAccessRule.breakGlassAllowed is sole bypass authority; ApprovalPolicy never bypasses. | RECONCILED |
| 16 | AccessReview modes, immutable snapshot/usage evidence, dispositions/timing and due-advancement semantics resolve through F18-RD-16 and matrix E03/E15. | RECONCILED |
| 17 | F18-RD-02/08, canonical catalog and VS0-WRITER-008 name RoleAssignment controller as sole spec/status/lifecycle/effect writer; VS0-WRITER-022 isolates Membership and VS0-WRITER-023 preserves future ProfileAssignment authority. | RECONCILED |
| 18 | ExceptionGrant exact interval/subject/control, immutable expiry/revocation evidence, overlap/narrowing and FEATURE-0020 boundary resolve through F18-RD-17. | RECONCILED |
| 19 | Privileged request mode, approval mapping, activation deadline/duration, readiness, break-glass and review boundary resolve through F18-RD-14/15 and matrix D. | RECONCILED |
| 20 | Canonical catalog/ADH register closed operations, action classes, three binding variants and the sole assignment writer. | RECONCILED |
| 21 | AssuranceEvidence is a closed provider-neutral transient carrier; F18-RD-14/15 apply the same phishing-resistant AAL2-or-higher, fifteen-minute activation floor to JIT and break-glass. | RECONCILED |
| 22 | Resource-specific lifecycle remains indexed to exact F18-RDs; supersession, campaigns and item dispositions are no longer conflated in summaries. | RECONCILED |
| 23 | Governance components and exception control-owned registrations are closed by F18-RD-17/18 and represented in canonical catalog. | RECONCILED |
| 24 | FEATURE-0013 registration references the exact F18-RD-20 taxonomy and pre-publication obligation boundary without copying envelope ownership. | RECONCILED |
| 25 | F18-RD-21 and matrix H01 order structural/local validation before safe reference resolution; the compact digest delegates the exact order to that authority. | RECONCILED |
| 26 | Matrix C11/C17 explicitly cover ExactResource, CreateParent and ScopeOnly. | RECONCILED |
| 27 | Matrix Section 2.1 explicitly lists AccessGroup; group scenarios include canonical direct membership and exclude external/nested/dynamic authorization. | RECONCILED |
| 28 | Matrix H01 now uses the exact F18-RD-21 order. | RECONCILED |
| 29 | Matrix D09/D14 name the uniform JIT/break-glass phishing-resistant activation floor and privilegedAccessRule as sole break-glass bypass owner. | RECONCILED |
| 30 | Matrix B06/B13 and scenarios cover relationship uniqueness, Guest terminal expiry, synchronized recoverable staleness and inactive-owner privilege-reduction recovery. | RECONCILED |
| 31 | Matrix C05/C15 and F18-RD-08/09/18 use publication-time per-action resource-reach dominance and deterministic standing applicability rather than qualitative scope or cross-assignment synthesis. | RECONCILED |
| 32 | Matrix D06/D11 and E02 cover indirect approval/review beneficiaries, including responsible Humans behind direct Workload/System members of beneficiary AccessGroups, and atomic activation strictly before activationDeadline. | RECONCILED |
| 33 | Matrix E03 separates mode/state/disposition; F18-RD-17 closes exception terminal/retry and narrowing algebra. | RECONCILED |
| 34 | Matrix F02 distinguishes supersession pins from terminal retirement/current-evidence unavailability. | RECONCILED |
| 35 | Canonical architecture and matrix C14/H10 record Sovrunn/native-IAM intersection without provider implementation. | RECONCILED |
| 36 | Matrix H08 pins ASVS 5.0 L2/L3 architecture applicability and process is ADH → approval → application → final evidence review. | RECONCILED |

### 5.1 Independent-review correction ledger

The first independent review remains immutable rejection evidence. Renewal-01
closed F18-SEC-002 and F18-SEC-004. Renewal-02 independently closed every
F18-SEC and F18-REN finding for its reported defect. No row below overrides a
review record.

| Finding | Incorporated correction | Current gate state |
|---|---|---|
| F18-SEC-001 | F18-RD-02/08 and VS0-WRITER-008 make the RoleAssignment controller the sole spec/status/lifecycle/effect writer; Membership and future ProfileAssignment authority are separated into VS0-WRITER-022/023. | CLOSED_BY_RENEWAL-02 |
| F18-SEC-002 | F18-RD-09 requires `roleassignment.grant` over complete reach and one independent per-action scope/target/resource/time/delegation witness; cross-assignment synthesis is prohibited. | CLOSED_BY_RENEWAL-01 |
| F18-SEC-003 | F18-RD-16 applies a complete `ReviewBeneficiaryHumanSet`, including responsible Humans behind direct Workload/System group members, to Retain/Replace/due advancement at atomic decision/effect publication. | CLOSED_BY_RENEWAL-02 |
| F18-SEC-004 | F18-RD-14/15 apply one phishing-resistant AAL2-or-higher, maximum-fifteen-minute activation floor to normal JIT and break-glass. | CLOSED_BY_RENEWAL-01 |
| F18-SEC-005 | The digest is now a compact non-normative index, and the scenario, writer, membership, target, Standing, review-evidence/state, and lifecycle mirrors are reconciled to normative F18-RD authority. | CLOSED_BY_RENEWAL-02 |

### 5.2 Renewal-01 correction ledger

The renewal-01 review remains immutable rejection evidence. Renewal-02 verified
both findings closed for their reported defects.

| Finding | Incorporated correction | Current gate state |
|---|---|---|
| F18-REN-001 | The semantic digest was replaced by a compact non-normative readiness index. It distinguishes inherited and proposed inventories, points all exact contracts to the ADH, contains only explicit consistency sentinels, and uses the complete scenario inventory with F18-SCN-38 retired. | CLOSED_BY_RENEWAL-02 |
| F18-REN-002 | F18-RD-16 derives a bounded `ReviewBeneficiaryHumanSet` that expands an AccessGroup holder through each current direct Workload/System member to its current responsible Human; current/proposed replacement sets are united, unresolved evidence fails closed for Retain/Replace, and Revoke remains available. Matrix, standards, threat, feature, canonical and VS mirrors agree. | CLOSED_BY_RENEWAL-02 |

### 5.3 Renewal-02 correction ledger

Renewal-02 remains immutable rejection evidence. These approved corrections
require renewal-03 verification against the regenerated exact manifest.

| Finding | Incorporated correction | Current gate state |
|---|---|---|
| F18-REN2-001 | F18-RD-02/04/05/09/10/11/20/21/22 derive and validate `MembershipEnabledGrantEnvelope` for every AccessGroup assignment-effect expansion; renewal-03 verified the exact group-held RoleAssignment laundering path closed. The separate workflow-eligibility gap is tracked as F18-REN3-001 below. | CLOSED_EXACT_PATH_BY_RENEWAL-03 |
| F18-REN2-002 | DEC-0058 is `Superseded by DEC-0059` in active traceability, obsolete ownership/validation is historicalized, and `phase2r-drift-check` now compares active Decision Index/traceability status plus rebaseline replacement identity. | CLOSED_BY_RENEWAL-03 |

### 5.4 Renewal-03 correction ledger

Renewal-03 remains immutable rejection evidence. Renewal-04 verified both
findings closed and also closed F18-REN2-001 end-to-end.

| Finding | Incorporated correction | Current gate state |
|---|---|---|
| F18-REN3-001 | F18-RD-02/04/05/09/10/12/14/16/18/22 defined a bounded direct-Human eligibility contract, prohibited AccessGroupRef and group-derived qualification, kept eligibility separate from canonical action authorization, and rechecked it at decision/effect publication. Membership creates only assignment effects covered by the envelope. | CLOSED_BY_RENEWAL-04 |
| F18-REN3-002 | F18-RD-16/18 closed a non-empty OR reviewer list, exact admitted kinds for the reviewed payload, scope/target compatibility, separate `accessreview.decide`, fail-closed invalid evidence, decision/remediation recheck and post-eligibility beneficiary conflict. | CLOSED_BY_RENEWAL-04 |

### 5.5 Renewal-04 correction ledger

Renewal-04 remains immutable rejection evidence. Renewal-05 verified its
approved correction against the regenerated exact manifest.

| Finding | Incorporated correction | Current gate state |
|---|---|---|
| F18-REN4-001 | F18-RD-02/12/14/16/18/22 restrict every `EligibilityRef` to exactly one Human PrincipalRef satisfied by principal equality. RoleDefinition, RoleAssignment, AccessGroupRef, Membership, claims and all other relationships are structurally invalid eligibility sources. Policy/rule applicability plus the independently authorized canonical action determine scope and target reach. F18-SCN-49 and dependent mirrors cover role-grant and later-rule-publication laundering. | CLOSED_BY_RENEWAL-05 |

## 6. Future-feature and provider leakage sentinels

The final documentation/drift gate must find no FEATURE-0018 authority for:

- IdentityProvider/SCIM/OIDC/WebAuthn/SPIFFE adapter implementation or raw token;
- external/nested/dynamic group authorization;
- arbitrary condition/ABAC/workflow language or a real policy/workflow engine;
- provider-native principal, role, policy, permission, credential or IAM effect;
- ProfileAssignment, hierarchy/conflict/exception resolution or
  EffectiveGovernanceContext publication;
- sovereignty, placement, entitlement, quota, service, plugin, explanation or
  orchestration semantics;
- database, network, Kubernetes/provider SDK, provisioning or deployment model;
- requirements, design, tasks or implementation produced by this application.

## 7. Remaining delivery gates

- separate human authorization before requirements generation;
- separate design, tasks and implementation approvals;
- later executable F18-RD-22 proof and FEATURE-0018 feature gate.
