# ADH-2026-070 Appendix A — FEATURE-0018 domain semantic contract

> Normative attachment to ADH-2026-070. This file has no independent approval
> status or application authority.

## Package authority

Upon exact joint approval, this appendix becomes one immutable part of the
ADH-2026-070 three-file package. The package consists of the core ADH, Appendix A, and
Appendix B named by the core's `Approval-package authority` section.

Approval status is inherited exclusively from the core ADH and its external
approval evidence. Human approval must pin the exact bytes of all three files
through one Git commit or independent SHA-256 digests. After approval, any byte
change to this appendix or either companion file invalidates approval for the
whole package. Kiro must not validate, approve, amend, or apply this appendix
independently. No package file has precedence; any inconsistency is a blocker
requiring correction and renewed joint approval.

## Normative allocation

Appendix A carries F18-RD-01, F18-RD-03 through F18-RD-05,
F18-RD-07 through F18-RD-17, F18-RD-21, and F18-RD-24. Appendix B carries the
remaining registry, publication, conformance, and standards-validation groups.
All cross-group references resolve through the core's exact decision index.

The decision text originated in the former single-file ADH-2026-070 payload and
includes the human-approved F18-SEC-001 through F18-SEC-005 corrections applied
after the first independent security review. It must be reviewed and approved
as part of this exact renewed three-file package.

## Normative decision groups

### F18-RD-01 — Current-feature-only authority

FEATURE-0018 decides only the contracts and behavior required for its own
Phase 2R development. It may define implementation-neutral future adapter boundaries
but does not select, design, or implement future adapters or consumers.

The following are closed FEATURE-0018 exclusions rather than open decisions:

- identity-provider selection, token processing, federation, or SCIM sync;
- external group synchronization and native provider or Kubernetes IAM;
- real policy-engine or workflow-engine selection or execution;
- ProfileAssignment, inheritance, effective governance, or cross-profile
  exception resolution, which remain FEATURE-0020-owned;
- sovereignty facts/evidence, placement, enrollment, entitlement, quota,
  billing, provisioning, infrastructure execution, or AI explanation, which
  remain with their separately sequenced owning features; and
- production persistence or deployment topology beyond Phase 2R authority.

These exclusions authorize no future resource, adapter, vendor, protocol,
credential flow, deployment boundary, requirement, task, or implementation.

### F18-RD-03 — Stable principal identity

Durable external identity is `issuer + subject`. Email, display name, username,
persona, job title, and external role label are non-authorizing attributes.
The closed principal categories are Human, Workload, and System.
FEATURE-0018 consumes an already-authenticated `PrincipalRef` and
canonical assurance evidence; it does not validate tokens or contact an
identity provider.

`AssuranceEvidence` is an operation-local value, not a resource or credential.
It binds the exact PrincipalRef, authoritative `authenticatedAt`, closed
assurance level `AAL1 | AAL2 | AAL3`, boolean `phishingResistant`, trusted
`sourceAuthority`, opaque integrity reference, and evidence expiry when one is
supplied. The evidence PrincipalRef must equal the authenticated actor and its
age is calculated only from authoritative UTC time. Missing, expired,
unrecognized, mismatched, or untrusted evidence fails closed. Raw tokens,
assertions, authentication secrets, and provider-native method payloads are
never retained. Deterministic fixtures supply this already-validated carrier;
real IdP and authenticator validation remains excluded.

### F18-RD-04 — Direct principal and scoped AccessGroup assignments

`RoleAssignment.roleHolderRef` is exactly one `PrincipalRef` or
`AccessGroupRef`. `AccessGroup`:

- is scoped to exactly one Organization or CloudProvider;
- has one responsible Human PrincipalRef owner;
- cannot authenticate;
- contains only direct `PrincipalRef` membership in v1;
- rejects nested and dynamic groups;
- becomes authorization-ineligible when suspended;
- becomes terminally authorization-ineligible when retired, prevents new
  membership and assignment, and retains immutable history;
- makes every otherwise-valid group-derived grant immediately ineffective when
  suspended or retired; and
- never derives access directly from an OIDC or other external group claim.

AccessGroup lifecycle is `Active -> Suspended -> Active` and `Active |
Suspended -> Retired`. Creation publishes Active only after authorization and
required audit-obligation acceptance. Retired is terminal.

The responsible owner must resolve to a current Human PrincipalRef at group
creation and every privilege-increasing operation. An inactive or unresolved
owner blocks adding or reactivating members, expanding the group, and creating
or expanding group-held assignments, but is not a runtime dependency for
already-effective access and does not silently disable it. Membership removal
or revocation, assignment revocation, group suspension or retirement, and
authorized owner recovery remain permitted. Owner failure triggers protected
review and escalation; owner recovery must use independently authorized
administration and cannot derive authority from the inactive owner.

Creating or reactivating a direct AccessGroup Membership, or restoring or
extending its assignment-effect interval through synchronized provenance or
freshness, is a privilege-increasing operation whenever it makes an existing
group-held RoleAssignment newly effective or effective for a longer interval.
It must pass F18-RD-09's membership-enabled grant ceiling at the Membership
controller's atomic publication boundary. Group ownership, administrator
status, trusted-provisioner status, or controller identity never bypasses that
ceiling.

An external group affects authorization only after trusted provisioning into
canonical `AccessGroup` and `Membership` resources. Privileged elevation is
always individual, approved, TimeBound, and auditable.

`AccessGroup` is a relationship target, never a `ScopeKind`. A group Membership
retains its Organization or CloudProvider `metadata.scopeRef` and carries one
exact UID-pinned `accessGroupRef`. The Membership and referenced AccessGroup
must have the same canonical owning scope. Cross-scope group membership is
rejected.

Role-holder reference compatibility is exact. An Organization-
scoped AccessGroup may hold a RoleAssignment at that Organization or its
OrganizationUnit, Tenant, or Project descendants. A CloudProvider-scoped
AccessGroup may hold a RoleAssignment only at that exact CloudProvider.
AccessGroup-held Platform and CloudPlatform assignments are prohibited in v1.
AccessGroup and group-derived assignments never supply approver, requester, or
reviewer eligibility in v1. Compatibility never grants an action, and policy
may narrow but never broaden it.

### F18-RD-05 — Membership is non-authorizing

Membership associates one principal with an Organization or CloudProvider
belonging context, or records one direct AccessGroup relationship within that
same Organization or CloudProvider scope. It never grants an action; establishes
approver, requester, reviewer, JIT, or break-glass eligibility; acts as a role;
or stores `roleAssignmentRefs`. A closed context kind distinguishes
Organization belonging, CloudProvider belonging, and AccessGroup membership.
Only AccessGroup membership carries `accessGroupRef`; the reference is a
relationship and never another scope authority. AccessGroup membership
contributes only when an independently active RoleAssignment names that exact
AccessGroup.

Membership retains immutable `sourceAuthority`, optional opaque
`sourceObjectRef`, `provisionedAt`, and, when synchronized,
`lastSynchronizedAt` and system-owned `freshUntil`. A synchronized membership
is authorization-ineligible at `freshUntil`. `membershipType` is exactly
`Standard | Guest`. Guest Membership is Human-only and requires a Human
PrincipalRef sponsor and `expiresAt`; Workload and System Memberships must be
Standard and require `responsiblePartyRef`. Raw tokens, assertions,
credentials, claims, and provider-native group schemas are prohibited.

Membership identity and uniqueness are deterministic:

```text
MembershipRelationshipKey =
  principalRef
  + contextKind
  + metadata.scopeRef
  + accessGroupRef when contextKind = AccessGroup
```

At most one non-terminal Membership may exist for a relationship key. Active
and recoverable Suspended relationships are non-terminal; creating another
record cannot bypass suspension or staleness. Membership retained lifecycle is
`Active -> Suspended -> Active` and `Active | Suspended -> Revoked`. Guest
expiry and synchronized freshness expiry are authoritative-time eligibility
projections, not mutable lifecycle transitions. `GuestExpired` is terminal for
that Membership UID and a new sponsored period requires a new Membership UID.
`Stale` is recoverable: the same synchronized Membership may become eligible
again only through an authorized update to its explicitly system-owned
freshness/provenance fields. Either boundary may emit protected audit evidence
and refresh a non-authoritative eligibility projection but never mutates the
immutable relationship specification.

For a Workload or System principal, the responsible party must resolve to a
current Human PrincipalRef when a new Membership is created and when an
AccessReview retains it. This is retained accountability evidence, not an
independent permission or a per-request authorization condition. Later
responsible-party inactivity blocks new Membership creation or review retention
and triggers protected review/escalation, but does not by itself make an
otherwise-current production Workload or System grant ineffective. Explicit
policy suspension or revocation may make it ineffective.

Membership replacement and renewal are not v1 operations. Synchronization may
update only its explicitly system-owned freshness and provenance fields. A new
guest period or changed immutable relationship requires a new Membership after
the previous relationship expires or is revoked; it never mutates the prior
record or rebinds an assignment pinned to the prior UID.

Organization belonging, CloudProvider belonging, and AccessGroup relationships
remain independently evaluated. Multiple memberships never union identity,
scope, or authority. Suspension or revocation immediately disables only the
direct or group-derived authorization that depends on that exact membership
context; it does not disable unrelated assignments in another context.
Where F18-RD-09 requires Membership for a direct-principal assignment, that
RoleAssignment pins the exact authorizing Membership UID. Suspension or
revocation of the pinned Membership makes that assignment ineffective;
revocation is permanent for that assignment. Reactivation after suspension may
restore an otherwise-current assignment, but a new or later Membership
never revives an assignment bound to a revoked Membership. Guest-assignment
`expiresAt` cannot exceed the pinned guest Membership `expiresAt`. Group-held
assignments remain bound to the AccessGroup; each principal's eligibility is
instead determined from their exact current direct AccessGroup Membership.
Because guest identity is established only through an expiring Membership,
guest RoleAssignment at Platform or CloudPlatform is prohibited in v1.

Membership remains non-authorizing even though an AccessGroup assignment-
effect expansion is grant-producing: the Membership contributes no action of its own
and can enable only independently current RoleAssignments already held by the
exact group. Creation, reactivation, or a freshness/provenance change that
restores or extends that assignment effect must pass the deterministic envelope rule
in F18-RD-09. Suspension, revocation, expiry enforcement, and other
privilege-reducing transitions remain permitted without a complete envelope
and never silently retain or enlarge access.

### F18-RD-07 — Versioned RoleDefinition

RoleDefinition contains only registered canonical Sovrunn actions; provider-
native IAM permissions and unrestricted wildcards are prohibited.

Effective lifecycle is:

```text
Draft -> Published
Published -> Suspended
Suspended -> Published
Published | Suspended -> Retired
```

Draft is mutable. Published versions are immutable. Changes create a new
version. Supersession is an informational version-lineage relationship, not an
authorization lifecycle state. Version lineage is linear in v1. Publishing a
successor atomically records exact
`supersededByRef` on the previously current published version through the
publication controller and emits the registered supersession AuditEvent.
Superseded versions receive no new assignments, requests, or profile
references, while independently valid existing exact pins remain effective.
Emergency suspension makes every assignment pinned to that exact version
ineffective immediately without mutating the definition or assignment and may
be reversed through restoration. Retirement is terminal: a retired version
receives no new references and every assignment pinned to that exact version is
immediately and permanently authorization-ineligible. Non-disruptive migration
uses supersession; retirement is not a migration mechanism. Suspension,
restoration, and retirement require current authority, justification, mandatory
`authorization-decision/v1`
DecisionRecord evidence, and protected audit-obligation acceptance before
lifecycle publication.

Each published RoleDefinition version declares an immutable base privilege
classification of `Ordinary` or `Privileged`. Its effective classification is
`Privileged` when the definition declares it, any registered action requires
privileged handling, a system-selected applicable `privilegedAccessRule`
elevates it, or any required classification input is missing or ambiguous.
Classification sources may elevate but never downgrade. FEATURE-0018 validates
but does not resolve the applicable rule. Publication validates the base and
registered-action classifications; authorization provenance records every
classification source that affected the result. A privileged role cannot
receive a Standing or AccessGroup assignment and uses an approved direct-
PrincipalRef TimeBound flow through either a PrivilegedAccessRequest or an
authorized RoleAssignmentProposal. The production Workload/System Standing
exception in F18-RD-08 applies only to effectively Ordinary roles.

### F18-RD-08 — Scoped RoleAssignment

One RoleAssignment binds:

```text
one roleHolderRef (PrincipalRef | AccessGroupRef)
+ one published RoleDefinition identity and exact version
+ metadata.scopeRef as the sole canonical target scope
+ optional exact resourceRef governed by that scope
+ exact authorizing membershipRef when F18-RD-09 requires Membership for a
  direct PrincipalRef holder
+ exact responsiblePartyRef for a direct Workload or System holder
+ explicit validity mode and mode-specific timestamps
```

`resourceRef` narrows an assignment to one exact resource and never becomes a
second scope authority. A resource-narrowed assignment contributes only to an
`ExactResource` authorization whose requested target has the same UID and a kind
permitted by the action owner's target-binding registration. It contributes
nothing to `CreateParent`, `ScopeOnly`, a different resource, or an incompatible
kind; another independently applicable assignment may still grant that request.
Assignment creation rejects when none of the pinned role version's actions can
apply to the resource kind. Holder, exact role version, scope, resource
restriction, required authorizing Membership UID, responsible party, pinned
Standing review rule, and validity are immutable;
change occurs through replacement or revocation. `membershipRef` is prohibited
for an AccessGroup holder and for Platform or CloudPlatform authorization,
where F18-RD-09 intentionally requires no Membership. `responsiblePartyRef` is
system-selected from trusted accountability evidence and must identify a
current Human PrincipalRef when the assignment is created or replaced and when
an AccessReview retains it. Where a Workload/System Membership is required, the
assignment value must equal that Membership's responsible party; at Platform
or CloudPlatform it is carried directly by the assignment. Missing, mismatched,
or non-Human responsibility fails closed at those boundaries.

The RoleAssignment controller is the sole canonical publisher and sole writer
of RoleAssignment specification, lifecycle, status, revocation, replacement,
expiry, and review-due effects. RD-13, RD-14, and RD-16 own distinct
authorization triggers that may submit only an exact immutable materialization,
revocation, replacement, or due-advancement intent; administrators and
originating workflow controllers cannot directly create, modify, revoke,
replace, expire, or publish a RoleAssignment. The RoleAssignment controller
revalidates the exact holder, role version, scope, resource restriction,
Membership, responsibility, validity, privilege, concurrency, idempotency, and
required decision/audit evidence before publishing the effect, and it cannot
enlarge the submitted intent.

Responsible-party evidence supplies accountability, not another permission or
per-request authorization condition. Later responsible-party inactivity blocks
new, replacement, renewed, or review-retained access and triggers protected
review/escalation, but does not alone make an otherwise-effective production
Workload or System assignment ineffective. Explicit policy suspension or
revocation may make it ineffective.

`Standing` has no automatic expiry, requires periodic review, and is
immediately revocable. `TimeBound` requires both `notBefore` and `expiresAt`
and becomes ineffective at exact expiry without grace. Privileged, guest, JIT,
and break-glass assignments are TimeBound. A Standing assignment is valid only
when its exact applicable access-review rule registers the holder's exact kind
in `allowedStandingHolderKinds`, the assignment ScopeKind in
`allowedStandingScopeKinds`, and, when resource-narrowed, the target kind in
`allowedStandingTargetKinds`; the role is effectively Ordinary; every
responsibility requirement passes; and periodic review applies. Production
Workload and System assignments may therefore be Standing only through this
exact registered path. Qualitative scope terms such as "narrow" have no
authorization meaning. Expiry retains the assignment and its evidence.

Every Standing assignment pins one exact system-selected published
`accessReviewRuleRef` and has system-owned `status.nextReviewDueAt`. Creation
derives the initial due time from that rule. Only a completed current
`StandingCertification` AccessReview with `Retain` for the exact current
assignment version advances the due time by the pinned rule's interval; a
successfully published `Replace` under F18-RD-16 produces a new assignment with
its own rule and due time. For Human or AccessGroup-held Standing access, reaching
the due time without current completed certification makes the assignment
authorization-ineligible until certified or replaced. The bounded production
Workload/System availability exception in F18-RD-16 remains the only v1
overdue-review exception. Review status never mutates the immutable assignment
specification.

FEATURE-0018 v1 has no user-authored RoleAssignment condition language, deny
assignment, wildcard condition, arbitrary expression, or generic ABAC
document. Scope and explicit validity are the only assignment-local
restrictions; contextual restrictions are mandatory guardrails evaluated
through the approved policy boundary.

RoleAssignment retained lifecycle is `Current -> Revoked | Expired`.
`Expired` is an automatic terminal transition only for a TimeBound assignment
at exact `expiresAt`; Standing never automatically expires. Its effective
projection is exactly `NotYetValid | Effective | InactiveDependency | Revoked |
Expired`. A RoleDefinition version that is authorization-ineligible under
F18-RD-07, an ineligible group or Membership, an overdue review, or an
unavailable mandatory guardrail produces `InactiveDependency` without mutating
immutable grant intent. Recovery is possible only when both the assignment
lifecycle remains Current and the dependency permits recovery; a retired
RoleDefinition version never recovers.

### F18-RD-09 — Deterministic scoped authorization composition

Authorization uses explicit grant union constrained by guardrail intersection:

```text
GrantSet = union(actions from every independently applicable RoleAssignment)

Allow iff:
  authenticated PrincipalRef is current
  AND every membership required by the requested scope and holder is current
  AND requested action is in GrantSet
  AND exact target and canonical scope match
  AND every contributing resourceRef matches the exact target
  AND every contributing assignment is valid
  AND every mandatory guardrail permits
  AND no policy outcome is Deny or Indeterminate
  AND any operation-local approval or exception evidence required now is current and exact
```

Direct and AccessGroup-derived assignments participate in the same union, but
each assignment must independently apply. Restrictions from different
assignments never synthesize a new grant. Policy Allow, approval, exception,
and guardrail evidence preserve or reduce authority and never create it.
`RequiresApproval` is not Allow; current authorization is re-evaluated after
approval.

Approval is admission-time evidence. It must be current when a RoleAssignment
or privileged activation is published. After publication, runtime authorization
validates the immutable approval provenance and linkage but does not require the
originating ApprovalRequest to remain current. The published RoleAssignment's
validity, dependencies, review state, and runtime guardrails determine continuing
authorization. Expired originating approval evidence never revokes or expires
an independently valid published assignment.

The AuthorizationInput target form must match the action owner's registered
target-binding mode and kind. A resource-narrowed assignment is independently
inapplicable unless F18-RD-08's exact-resource rule matches; its mismatch does
not contaminate or broaden another assignment. Missing or ambiguous registration
denies before grant union is evaluated.

No grant, invalid/expired evidence, raw external group claim, nested/dynamic
group membership, unresolved or cross-tree scope, resource mismatch,
conflicting evidence, or `Indeterminate` denies. Policy Deny overrides a role.
Explicit deny assignments are not introduced.

Membership eligibility is contextual, not universal. Authorization in an
Organization, OrganizationUnit, Tenant, or Project requires current belonging
in the governing Organization. CloudProvider authorization requires current
belonging in that exact CloudProvider. AccessGroup-derived authorization also
requires current direct Membership in the exact same-scope AccessGroup and an
authorization-eligible group. Platform and CloudPlatform authorization use
current PrincipalRef and RoleAssignment authority without inventing an
Organization or CloudProvider Membership. Suspension or revocation affects only
assignments whose authorization depends on that exact Membership.

The only inherited containment chain is Organization -> OrganizationUnit ->
Tenant -> Project. Platform, CloudPlatform, and CloudProvider do not implicitly
contain one another. Undefined containment is exact-scope only.

Delegated administration has a mandatory grantor ceiling. The proposed grant
has one exact reach derived from its registered target-binding mode, canonical
scope, and optional exact resource restriction. A grantor must have current
`roleassignment.grant` authority covering that complete proposed reach and,
for every proposed action, one independently applicable delegable authority
witness covering the same or a broader permitted reach. `ExactResource(A)`
authority covers only resource A; it covers neither resource B nor a scope-wide
grant. Scope-wide authority may cover an exact resource governed by that scope.
`CreateParent` requires authority covering the exact registered parent scope,
and `ScopeOnly` requires applicable scope-wide authority. Existing canonical
scope containment may establish broader scope reach only where F18-RD-09
explicitly defines it.

Each per-action witness must independently cover action, scope, target-binding
mode, resource reach, and proposed validity. Restrictions or reach fragments
from different assignments cannot be combined to manufacture a broader
witness. The grantor additionally requires `roleassignment.delegate` when the
proposed role contains either `roleassignment.grant` or
`roleassignment.delegate`; a non-delegable grant cannot contain either action.
Grantor authority is evaluated at RoleAssignment publication. A TimeBound
assignment cannot extend beyond its applicable witnesses' then-current temporal
authority; a Standing assignment therefore requires then-current Standing
witnesses. Later grantor suspension, expiry, or revocation does not silently
cascade to an already-published assignment. Explicit assignment revocation,
policy suspension, or access review handles that assignment. Runtime-dependent
delegation chains and cascading revocation are excluded from v1; each
downstream grant must independently pass the full grantor ceiling at its own
publication. No actor may manufacture action, scope, target mode, resource
reach, duration, or delegation capability.

An AccessGroup Membership assignment-effect expansion is a separate grant-
producing publication boundary. It includes Membership creation or
reactivation and any
synchronized freshness or provenance update that restores current eligibility
or extends its future eligibility interval. Immediately before publication,
the Membership controller derives one operation-local
`MembershipEnabledGrantEnvelope` from a coherent current snapshot. The envelope
contains every independently applicable group-held RoleAssignment that would
become effective or remain effective for a newly extended interval if the
exact Membership intent were published. Each entry preserves the exact
RoleAssignment and RoleDefinition versions, action, scope, target-binding mode,
optional exact resource reach, and only the newly enabled or extended effect
interval. That interval is the intersection of the assignment validity, Guest
expiry or synchronized `freshUntil` when present, AccessGroup eligibility, and
every other current dependency. Restrictions are never flattened or unioned
into broader authority.

Publishing a non-empty envelope requires all of the following:

1. one independently applicable authority for the exact Membership operation;
2. current `roleassignment.grant` authority covering the complete envelope;
3. for every envelope action/reach/validity tuple, one independently applicable
   delegable-authority witness satisfying the same action, scope, target mode,
   exact-resource reach, validity, and delegation rules used for a direct
   RoleAssignment grant; and
4. accepted concurrency, decision, and audit obligations under F18-RD-20 and
   F18-RD-21.

Every authority relied on must cover the full effect interval it authorizes. A
Standing group-derived effect therefore requires Standing Membership-operation,
`roleassignment.grant`, and per-action witnesses. A TimeBound administrator or
provisioner cannot enable access beyond its applicable temporal ceiling.
Witness restrictions from different assignments cannot be synthesized.
Missing, stale, ambiguous, conflicting, cross-scope, resource-incompatible, or
otherwise unresolved group, membership, assignment, role, owner,
responsibility, validity, or witness evidence fails closed and publishes no
assignment-effect expansion or completed replay.

An empty envelope grants no RoleAssignment action and requires only the
otherwise-applicable Membership operation authorization. Membership cannot
establish approval, reviewer, JIT, or break-glass eligibility under F18-RD-12,
F18-RD-14, or F18-RD-16, so an empty envelope has no separate eligibility
effect. If a group-held RoleAssignment is
published later, that RoleAssignment publication independently passes the
ordinary F18-RD-09 grantor ceiling. Membership and group-held RoleAssignment
publication each re-evaluate the relevant current member/assignment predicate
at their own atomic publication boundary against one coherent snapshot; a
concurrent change conflicts or retries rather than using a stale snapshot.
Thus whichever operation completes second observes and authorizes the effective
access it creates.

A trusted provisioner participates only as a canonical System PrincipalRef
with current, exact Sovrunn authority and the same envelope ceiling. External
claims, source-system authority, controller identity, or provisioning trust
never substitute for that authorization. Later expiry, suspension, or
revocation of the publishing actor's authority does not silently cascade to a
published Membership; explicit Membership/group/assignment suspension,
revocation, expiry, review, or a newly authorized publication controls
continuing access.

This ordinary-administration ceiling does not require a JIT or break-glass
requester to already hold the requested privileged actions. That requester must
hold `privilegedaccessrequest.submit` at the exact scope and satisfy the exact
pre-authorized eligibility rule. The privileged-access controller owns the
activation trigger, not grant authority or RoleAssignment publication: it may
submit only the exact role, holder, scope, resource restriction, responsibility,
and interval authorized by the accepted request and current approval or bypass
evidence under F18-RD-14/F18-RD-15. The RoleAssignment controller alone may
publish that intent under F18-RD-08.

`AuthorizationInput` and `AuthorizationResult` are operation-local values, not
resources, bearer tokens, sessions, or retained evidence. Input binds actor,
action, exactly one exact-resource target, create-parent scope, or scope-only
target, authoritative current time,
assurance evidence, request ID, and correlation. Result binds the exact outcome
`Allow | Deny`, evaluation instant, contributing assignment, membership, group,
and role versions, target/scope relationship, policy, approval, and exception
evidence, and stable reasons. It cannot authorize a different request. Protected
DecisionRecord/AuditEvent projections retain evidence; safe projections reveal
no inaccessible target, confidential policy input, assertion, or secret. A
validation, evaluation, or mandatory audit-obligation failure returns a safe
error and no AuthorizationResult; FEATURE-0017 `Indeterminate` is successfully
mapped to `Deny` with stable protected provenance.

### F18-RD-10 — FEATURE-0017 adoption without reinterpretation

FEATURE-0018 consumes FEATURE-0017 outcomes exactly:

| FEATURE-0017 outcome | FEATURE-0018 interpretation |
|---|---|
| `Allow` | Continue all remaining FEATURE-0018 checks |
| `Deny` | Deny |
| `RequiresApproval` | Consume one exact system-selected ApprovalPolicy reference and evaluate it; never create a request automatically |
| `Indeterminate` | Fail closed |

Transport success is not authorization. The FEATURE-0017 fake stays digest-
keyed and contains no IAM, approval, exception, target, or governance logic.

Every grant-producing operation derives exactly one system-owned,
operation-local `ApprovalRequirement`: `NotRequired` or `Required` with one
exact system-selected published ApprovalPolicy version ref. The value is not
requester-authored, independently stored, or bearer authority. Missing or
ambiguous ApprovalRequirement evidence fails closed. `Required` with zero or
multiple applicable policies, or with an unpublished, mismatched, or unavailable
policy ref, fails closed and creates no ApprovalRequest. `NotRequired` carries
no ApprovalPolicy ref and an unexpected ref rejects.

The originating operation controller is the sole derivation authority before
intent acceptance: the grant controller for a RoleAssignment grant/proposal,
the privileged-access controller for a PrivilegedAccessRequest, and the
exception controller for an ExceptionProposal. The Membership controller is
the sole derivation authority for an AccessGroup assignment-effect expansion.
Derivation consumes mandatory
architecture rules plus exact trusted system-selected evidence; it does not
perform FEATURE-0020 resolution. The authorization and Approval-request
controllers validate and consume the result but cannot replace it.

JIT PrivilegedAccessRequest, effectively Privileged RoleAssignmentProposal,
and ExceptionProposal operations derive `Required`. A grant of an effectively
Ordinary role may derive `NotRequired`; the grant operation still requires the
intrinsically Privileged `roleassignment.grant` action under F18-RD-06.
An AccessGroup Membership assignment-effect expansion always derives
`NotRequired` in v1 because an AccessGroup cannot hold an effectively
Privileged role and Membership is not an ApprovalPolicy subject. `NotRequired`
never weakens F18-RD-09's complete membership-enabled grant ceiling and never
creates approver, reviewer, JIT, or break-glass eligibility. Adding a
membership-approval subject, group-based eligibility, or an entitlement-
manager exception requires separate architecture change control.
BreakGlass derives `Required`
unless the exact system-selected `privilegedAccessRule` in F18-RD-15 has
`breakGlassAllowed=true`; an accepted bypass derives `NotRequired` for
ApprovalRequest creation while retaining that exact rule as bypass evidence. A
missing, ambiguous, unpublished, or incompatible privileged-access rule fails
closed rather than falling back to approval. A FEATURE-0017
`RequiresApproval` outcome requires `Required` but never creates a request by
itself; an `Allow` outcome cannot weaken an independently mandatory `Required`.

A requester cannot select or replace the requirement or policy. FEATURE-0018
does not resolve GovernanceProfile assignment, hierarchy, inheritance,
applicability, or conflicts. Deterministic FEATURE-0018 fixtures provide trusted
preselected requirement and policy evidence. FEATURE-0020 later supplies
production selection through effective governance resolution without changing
FEATURE-0018 approval semantics.

### F18-RD-11 — Bounded FEATURE-0017 subject/target use

FEATURE-0017 v1 is not a universal principal/action/target authorization
engine. FEATURE-0018 may invoke FEATURE-0017 only for an exact UID- and
generation-pinned PrivilegedAccessRequest at Project, CloudPlatform, or
CloudProvider scope, using canonical action
`privilegedaccessrequest.submit`. It does not invoke FEATURE-0017 for ordinary
authorization, RoleAssignmentProposal, ExceptionProposal, AccessReview,
Membership, AccessGroup, or group evaluation. An unsupported subject or scope
never causes contract emulation. When a FEATURE-0017 evaluation is mandatory
but no compatible evaluation can be performed, the operation fails closed.

The FEATURE-0017 exclusion for Membership does not exempt an AccessGroup
assignment-effect expansion from authorization. Its
`MembershipEnabledGrantEnvelope`
and grantor ceiling are deterministic FEATURE-0018-local authorization algebra,
not policy emulation, approval evidence, or a new FEATURE-0017 subject.

Ordinary role/scope/target authorization and RBAC algebra remain FEATURE-0018-
owned and are not a custom policy engine. Preselected ApprovalPolicy and
governance evidence remain independently system supplied; absence of an
optional FEATURE-0017 evaluation never fabricates an outcome. Target-aware
dynamic policy for ordinary authorization requires a separately approved
FEATURE-0017 version; FEATURE-0018 must not overload `subjectRef` or emulate the
missing contract.

### F18-RD-12 — Bounded ApprovalPolicy

ApprovalPolicy is declarative and supports one to three linear ordered stages,
eligible approvers, required approval count per stage, decision expiry,
justification rules, self-approval restrictions, and separation of duties.
It is evaluated only for `ApprovalRequirement.Required`. A `NotRequired`
effectively Ordinary role grant or accepted break-glass bypass does not select
or evaluate an ApprovalPolicy. Break-glass bypass authority is defined solely
by F18-RD-15; an ApprovalPolicy always represents actual approval.
Its closed applicability key is `subjectKind + canonical action + ScopeKind +
optional targetKind + requestMode`. `subjectKind` is exactly
`PrivilegedAccessRequest`, `RoleAssignmentProposal`, or `ExceptionProposal`;
`requestMode`, when applicable, is exactly `JIT` or `BreakGlass` and is
permitted only for PrivilegedAccessRequest under F18-RD-14 and F18-RD-15.
Applicability never names an exact target instance,
evaluates an arbitrary attribute, or performs hierarchy/profile resolution.
The system-selected exact policy version must match every applicable key
component or evaluation fails closed under F18-RD-10.

The eligible-approver list is non-empty and every entry is one `EligibilityRef`
under F18-RD-02: exactly one Human PrincipalRef. It is satisfied only by exact
principal equality. RoleDefinition, RoleAssignment, AccessGroupRef, Membership,
external group claim, nested group, and dynamic group never satisfy or create
approver eligibility. The list uses OR semantics. The selected ApprovalPolicy
applicability and the independently required action determine exact operation
scope and target compatibility; EligibilityRef never broadens either.

Eligibility is necessary but never grants `approvalrequest.decide`; that
canonical action remains independently required. Empty, duplicate, malformed,
unsupported, inactive, stale, ambiguous, scope-incompatible, or target-
incompatible evidence fails closed.
Approver eligibility and separation of duties are re-evaluated when a decision
is accepted and immediately before downstream effect publication. One
principal may contribute at most one immutable effective decision to a stage.
Same-decision replay is idempotent; a different decision from that principal
conflicts. Decisions for an inactive or completed stage reject.

Separation of duties uses an exact current conflict set containing the
requester; every direct PrincipalRef beneficiary; the owner and every current
direct principal member of a beneficiary AccessGroup; the responsible Human for
a beneficiary Workload or System principal; and every beneficiary resolved
through an exact RoleAssignmentRef or PrivilegedAccessRequestRef. A conflicted
principal cannot approve the request. The set is re-evaluated both when each
approval decision is accepted and immediately before downstream effect
publication. If membership, ownership, responsibility, or referenced evidence
changes and creates a conflict, no effect is published from that approval; a
fresh approval request is required. A principal counted in one stage cannot
count in another stage of the same request. Submission authority, approver
eligibility, and effect-writing authority are independently evaluated; an
ApprovalRequest carries no delegated approval or effect authority.

Stages execute strictly in order. Any valid `Deny` terminates the whole request
as `Denied`; otherwise a stage succeeds when its configured approval quorum is
met by distinct eligible principals and only then activates the next stage.
Quorum is a positive integer. A runtime eligible set smaller than quorum
creates no approval; the stage remains pending until eligibility recovers or
the request's immutable `expiresAt` is reached, at which point the request
becomes `Expired`.
The policy's decision-expiry duration is greater than zero and no more than
twenty-four hours; a narrower platform-approved policy value may shorten but
never extend that ceiling.
ApprovalPolicy supports no abstain decision. Arbitrary branching, loops,
delegation chains, scripted escalation, expressions, and a general workflow
engine are prohibited. Curated single-stage, privileged two-stage, and
accelerated emergency-approval templates are the default; an emergency template
still performs approval and never implies bypass. Custom policy is an advanced
security-administrator action.

Effective lifecycle is `Draft -> Published -> Retired`; Draft is mutable and
published versions are immutable. F18-RD-07's common informational
supersession relationship applies: a superseded policy receives no new
requests, while exact version-pinned existing requests retain the policy until
their independent validity ends.

### F18-RD-13 — Immutable terminal approval evidence

ApprovalRequest lifecycle is:

```text
Pending -> Approved | Denied | Cancelled | Expired
```

Terminal outcomes are immutable. The subject is exactly one of a
`PrivilegedAccessRequestRef`, embedded immutable `RoleAssignmentProposal`, or
embedded immutable `ExceptionProposal`. A role proposal pins holder, role
version, scope, optional resourceRef, system-selected exact authorizing
Membership UID when required for a direct principal, exact Human
responsiblePartyRef for a Workload/System holder, exact accessReviewRuleRef for
Standing, and validity. A proposal for an effectively Privileged role must be
direct-PrincipalRef and TimeBound; AccessGroup or Standing privileged proposals
reject. An exception proposal pins control, subject, scope, interval,
justification, and compensating controls. Arbitrary or future-domain payloads
are rejected.

An ApprovalRequest is created only for `ApprovalRequirement.Required` and only
from an explicit, successfully accepted FEATURE-0018 business intent: a
submitted PrivilegedAccessRequest, an authorized RoleAssignmentProposal
submission, or an authorized ExceptionProposal submission. A `NotRequired`
ordinary grant proceeds only through its independently authorized owning effect
controller and creates no ApprovalRequest. The originating operation owns
acceptance of user intent; the Approval-request controller is the sole creator
of the system-owned immutable request and pins its exact subject, policy, scope,
requester, and correlation. The approval controller alone owns lifecycle status
and decisions. A FEATURE-0017 `RequiresApproval` outcome never creates a request
by itself. Missing or ambiguous requirement evidence, or missing, ambiguous,
unpublished, mismatched, or unavailable policy-selection evidence for
`Required`, fails closed without creating one.

The request pins the one exact system-selected published ApprovalPolicy version
accepted under F18-RD-10. At creation it also pins one immutable `expiresAt`
equal to the earliest of creation time plus the policy's decision-expiry
duration and the applicable exact subject bound: a PrivilegedAccessRequest's
`activationDeadline`, a TimeBound RoleAssignmentProposal's proposed
`expiresAt`, or an ExceptionProposal's requested `expiresAt`. A Standing
RoleAssignmentProposal contributes no temporal subject bound. All stages share
the resulting value and no replay or stage transition extends it. An Approved
request remains historically Approved,
but its evidence is current only while authoritative time is before
`expiresAt`. Evidence includes exact approver PrincipalRef,
decision, timestamp, subject, policy version, required justification, and
correlation. F18-RD-12's self-approval and separation-of-duties rules apply
without restatement. Material proposal change requires a new request. Approver
eligibility is rechecked at decision acceptance and downstream effect
publication under F18-RD-12. The first valid terminal transition
wins under inherited FEATURE-0012 optimistic-concurrency semantics. One
originating intent creates at most one ApprovalRequest. Idempotency binds the
originating subject identity, actor, scope, and request digest: same-key/same-
digest submission returns the existing request and same-key/different-content
conflicts. Initial and terminal publication follow the audit-before-publication
rule in F18-RD-20 and the replay rules in F18-RD-21.

Approval remains evidence and never performs the approved effect. After
approval, the owning controller re-evaluates current authorization and exact
evidence. The privileged-access controller may authorize and submit the exact
linked TimeBound RoleAssignment intent; the authorized grant controller may
authorize and submit the proposal's exact RoleAssignment intent; and the
exception controller creates the exact ExceptionGrant. Only the RoleAssignment
controller may publish either RoleAssignment intent under F18-RD-08. The
approval controller creates none of these effects. Denial, cancellation, or
expiry authorizes none of them.

### F18-RD-14 — JIT privileged access authorizes one temporary grant

The requester and resulting role holder are the same exact Human PrincipalRef.
Workload/System time-bounded administration uses an authorized
RoleAssignmentProposal with responsiblePartyRef, not interactive JIT or
break-glass.

Every PrivilegedAccessRequest pins system-owned immutable `submittedAt`,
`mode`, `activationDeadline`, and requested access duration. Mode is exactly
`JIT` or `BreakGlass`. The atomic publication instant of both the linked
TimeBound RoleAssignment and the request's `Active` transition must be strictly
before `activationDeadline`; beginning validation or controller work before the
deadline is insufficient. For JIT, the deadline is no later than `submittedAt +`
the exact selected ApprovalPolicy decision-expiry duration, bounded by
F18-RD-12's twenty-four-hour platform maximum. For BreakGlass, the deadline is
no later than fifteen minutes after `submittedAt`; its exact privileged-access
rule may shorten that bound. An applicable rule or policy may shorten only its
owned deadline or duration ceiling and never extend it. Requested access
duration begins only at successful activation; activation validity is not
access validity.

The exact system-selected privileged-access rule carries a non-empty
`requesterEligibilityRefs` list whose entries are only `EligibilityRef` values
under F18-RD-02: exact Human PrincipalRefs satisfied only by exact requester
equality. RoleDefinition, RoleAssignment, AccessGroupRef, Membership, external
group claim, nested group, and dynamic group never satisfy or create requester
eligibility. The list uses OR semantics. The selected privileged-access rule
applicability and the independently required action determine exact request
scope and target compatibility; EligibilityRef never broadens either.

Eligibility is re-evaluated at request admission and atomic activation
publication. It is necessary but never grants
`privilegedaccessrequest.submit` or permission to perform the requested
privileged actions; those checks remain independent. Empty, duplicate,
malformed, unsupported, inactive, stale, ambiguous, scope-incompatible, or
target-incompatible evidence fails
closed; the resulting assignment remains individual.

PrivilegedAccessRequest lifecycle is:

```text
Submitted -> PendingApproval | Active | Cancelled | Expired

PendingApproval -> Approved | Denied | Cancelled | Expired
Approved -> Active | Cancelled | Expired
Active -> Expired | Revoked
```

`Denied`, `Cancelled`, `Expired`, and `Revoked` are terminal. The privileged-
access controller maps exact terminal ApprovalRequest evidence to the corresponding
pre-activation request state. `Approved` is non-terminal evidence, not
authorization. If fresh activation checks fail, no RoleAssignment intent is
accepted or published and the request remains `Approved` with system-owned
activation readiness `Blocked` until the condition is corrected, the request is
cancelled, or its approval evidence or activation deadline expires. Activation
retries remain idempotent and must not extend validity.

Approval-path activation requires every Membership contextually required by
F18-RD-09, current `privilegedaccessrequest.submit` authority and exact pre-
authorized eligibility, one exact system-selected current ApprovalPolicy,
approved request, current AssuranceEvidence of at least AAL2 with
`phishingResistant=true`, bounded duration, and no conflicting or expired
evidence. The AssuranceEvidence must be no more than fifteen minutes old at
activation. This phishing-resistant AAL2 floor applies uniformly to normal JIT
and break-glass activation; the applicable privileged-access rule may require
AAL3 or a shorter assurance age but cannot weaken the floor. Request submission and
approval do not themselves activate access or require the activation assurance
unless their independently applicable policy says so. F18-RD-15 exclusively
defines the direct break-glass bypass evidence. Neither path requires the
requester to already possess the requested privileged actions. Activation
performs a fresh authorization evaluation and may submit exactly one individual
TimeBound RoleAssignment materialization intent linked immutably to the request
and to approval or exact bypass evidence as applicable. The RoleAssignment
controller performs F18-RD-08 validation and is the sole publisher. The request
becomes Active only when that publication succeeds; failure publishes neither
an Active request nor an assignment. Renewal is a new request. Denial,
cancellation, or pre-activation deadline expiry authorizes no assignment
intent. Once Active, the request becomes Expired when its linked RoleAssignment
expires; the activation deadline does not terminate already-activated access.
Expiry uses the authoritative UTC time source. Lifecycle transitions and
publication follow F18-RD-21 and F18-RD-20 respectively.

Normal JIT defaults to one hour and has an eight-hour platform maximum. Policy
may shorten but never extend it. Assurance evidence must be no more than
fifteen minutes old at activation. A published JIT RoleAssignment has
`notBefore` equal to the authoritative successful-activation instant and
`expiresAt` equal to the earliest of:

- `notBefore +` the approved duration;
- the pinned required Membership expiry; and
- applicable privileged-access rule and policy access-duration ceilings.

An absent non-applicable bound does not participate. If no effective interval
remains, activation fails closed and accepts or publishes no assignment intent.
Delayed activation never extends the activation deadline, approved duration,
Membership, or rule ceiling, and retry never changes a successfully established interval. Once
timely activation succeeds, the activation deadline and originating
ApprovalRequest expiry do not truncate the published assignment. Activation
attempted after the deadline requires a new request rather than grace or
automatic extension.

### F18-RD-15 — Constrained break-glass access

Break-glass is a PrivilegedAccessRequest mode, not a resource. It is individual
only and requires pre-authorized eligibility, fresh AssuranceEvidence of at
least AAL2 with `phishingResistant=true`, explicit emergency justification,
immediate protected audit, exact scope, automatic expiry, and retrospective
AccessReview. It requires one exact system-selected privileged-access rule. The
rule's `breakGlassAllowed` boolean is the sole authority that permits bypass of
normal pre-approval; ApprovalPolicy never permits bypass. Break-glass never
bypasses authentication, scope isolation, expiry, or audit.

When that exact rule has `breakGlassAllowed=true` and every authentication,
eligibility, scope, assurance, duration, authorization, review-rule, and audit
precondition passes, F18-RD-14's privileged-access activation trigger derives
`ApprovalRequirement.NotRequired`, may authorize submission of exactly one
linked TimeBound RoleAssignment materialization intent, and creates no
ApprovalRequest. The RoleAssignment controller alone may publish that intent
under F18-RD-08. The exceptional direct transition `Submitted -> Active` occurs
only when the assignment publication and linked-review acceptance succeed.
When the exact valid rule has `breakGlassAllowed=false`, an otherwise-valid
request derives `ApprovalRequirement.Required` and follows the normal
`Submitted -> PendingApproval` path with one exact applicable ApprovalPolicy
and all F18-RD-14 checks. Missing, multiple, unpublished, mismatched, or
unavailable privileged-access rule evidence fails closed; it never becomes an
approval fallback. Authentication, structural validation, eligibility, scope,
request authorization, review-rule, or required audit-obligation failure also
never falls back to approval and cannot be repaired by an approver. A failure
detected before acceptance is safely rejected and creates no
PrivilegedAccessRequest.

An accepted bypass-authorized request whose current assurance, authorization,
eligibility, duration, required audit-obligation acceptance, or other
activation evidence
becomes stale, unavailable, failed, or conflicting remains non-active in
`Submitted` with system-owned activation readiness `Blocked`; it creates
no ApprovalRequest and authorizes no publishable RoleAssignment intent.
Correction may trigger an idempotent retry before `activationDeadline`, but cannot
extend any validity bound.
Direct activation never fabricates approval evidence. F18-RD-14's interval
calculation applies with the break-glass platform maximum below; the activation
deadline never truncates access after timely activation.

Maximum duration is one hour, assurance age is at most fifteen minutes,
mutation renewal is prohibited, and review is due within twenty-four hours.
The exact privileged-access rule may only shorten these bounds. Audit failure
blocks activation.

Successful direct activation atomically materializes exactly one linked Pending
AccessReview with immutable reviewed-subject, activation/grant references,
system-selected `accessReviewRuleRef`, originating `privilegedAccessRule`
evidence, and the exact `dueAt` derived by F18-RD-16. The privileged-access
controller may create only this mandatory campaign intent; the review controller
alone owns review state and findings. Failure to accept the linked review or its
required audit obligation blocks activation. Retry is idempotent and cannot
create a second review. The review being overdue never extends access, which
expires independently; overdue-review consequences follow F18-RD-16.

### F18-RD-16 — Snapshot-based AccessReview

AccessReview campaign mode is exactly `StandingCertification`,
`BreakGlassRetrospective`, or `Manual`. StandingCertification and
BreakGlassRetrospective each review exactly one linked RoleAssignment and use
that assignment's exact `metadata.scopeRef`. A Manual campaign carries a
non-empty, deduplicated list of exact UID- and
resourceVersion-pinned Membership or RoleAssignment references; it has no
dynamic query, expression, filter, or future-item inclusion. An Organization-
scoped Manual campaign may include items at that Organization or its canonical
OrganizationUnit, Tenant, or Project descendants. An OrganizationUnit-scoped
campaign may include that exact OrganizationUnit and its same-tree nested
OrganizationUnit, Tenant, or Project descendants. A Tenant-scoped campaign may
include that exact Tenant and its Project descendants. Project, Platform,
CloudPlatform, and CloudProvider campaigns contain exact-scope items only.
Cross-tree and undefined containment reject.

Every campaign pins exactly one system-selected published
`accessReviewRuleRef`. A StandingCertification campaign uses the assignment's
already-pinned rule and has `dueAt = assignment.status.nextReviewDueAt`. A
Manual campaign has `dueAt = createdAt + reviewWindow` from its pinned rule. A
BreakGlassRetrospective campaign has `dueAt = activationAt + min(24 hours,
retrospectiveReviewDeadline)` where `retrospectiveReviewDeadline` comes from the
exact system-selected `privilegedAccessRule` that authorized activation; it
also pins the applicable access-review rule for reviewer eligibility and review
evidence. Missing, multiple, unpublished, scope-incompatible, or mode-
incompatible required rule evidence fails closed. The access-review rule's
non-empty `reviewerEligibilityRefs` list is authoritative for all three modes
and contains only `EligibilityRef` values under F18-RD-02: exact Human
PrincipalRefs satisfied only by exact reviewer equality. RoleDefinition,
RoleAssignment, AccessGroupRef, Membership, external group claim, nested group,
and dynamic group never satisfy or create reviewer eligibility. The list uses
OR semantics. The selected access-review rule applicability and the
independently required action determine exact AccessReview scope and reviewed-
target compatibility; existing canonical containment applies only where F18-
RD-09 defines it and undefined containment is exact-scope. EligibilityRef never
broadens scope or target reach. Eligibility is necessary but never grants
`accessreview.decide`; that canonical action remains independently required.
Empty, duplicate, malformed, unsupported, inactive, stale, ambiguous, scope-
incompatible, or target-incompatible
evidence fails closed. Reviewer eligibility is re-evaluated when a review
decision is accepted and immediately before atomic remediation publication.
The `ReviewBeneficiaryHumanSet` separation-of-duties rules apply after
eligibility resolution. A conflict blocks Retain, Replace, and due-time
advancement; independently authorized privilege-reducing RoleAssignment or
Membership revocation remains available through its existing canonical path.

AccessReview lifecycle is:

```text
Pending -> InProgress | Cancelled | ExpiredIncomplete
InProgress -> Completed | Cancelled | ExpiredIncomplete
```

`Completed`, `Cancelled`, and `ExpiredIncomplete` are immutable terminal campaign
states. The immutable snapshot is captured exactly once on entry to InProgress;
an empty, changed, or unresolvable required snapshot fails closed rather than
silently completing the campaign. Before snapshot acceptance the campaign
remains Pending with system-owned readiness `Blocked` until corrected,
cancelled, or `dueAt`, when it becomes ExpiredIncomplete.

AccessReview captures an immutable snapshot of exact membership/assignment
identities and versions, principals, roles, scopes, validity, reviewers, and
review-policy evidence. It may retain an embedded immutable, versioned
`UsageEvidenceSummary` written only by the review controller. The summary binds
its evidence window, `lastObservedUseAt`, observed canonical actions,
generation time, freshness, coverage, and protected source AuditEvent
references. Usage informs review but never grants access or becomes an
authorization input. Missing use is not non-use unless coverage is complete;
stale, absent, and incomplete evidence is explicit. RoleAssignment gains no
mutable `lastUsedAt`, and FEATURE-0018 adds no usage-evidence resource or fourth
DecisionProfile.

Per-item dispositions are closed by reviewed kind. A Membership item permits
`Retain | Revoke | Stale`; a RoleAssignment item permits
`Retain | Revoke | Replace | Stale`. Membership Revoke targets the exact
snapshotted Membership and never creates a replacement. RoleAssignment Revoke
authorizes the review controller to submit one exact revocation intent;
RoleAssignment Replace authorizes it to submit one exact replacement-and-
revocation materialization intent. The RoleAssignment controller alone publishes
the revocation or atomically publishes the new assignment and revokes the exact
old assignment under the same local audit/concurrency boundary; the review
controller cannot create, revoke, or replace either RoleAssignment directly.
The Replace decision pins the complete replacement proposal. It must preserve
the holder and may only reduce actions, scope, resource reach, delegation, or
validity; any expansion requires a separate authorized RoleAssignmentProposal
flow under F18-RD-13. The replacement independently satisfies Membership,
responsibility, validity, review-rule, privilege, and grantor/effect validation.
Revoke/Replace is exact-item-bound, idempotent, and effective-once under
F18-RD-21; required publication follows F18-RD-20. A changed item yields Stale
with no remediation.
`ExpiredIncomplete` is a campaign state, not an item decision, and incomplete
review is never certification. No principal may retain, replace, or advance
review eligibility for authority from which that principal currently benefits.
For every disposition that preserves or replaces effective authority, the
controller derives the exact `ReviewBeneficiaryHumanSet` for an assignment as
follows:

```text
Human role holder
  -> that Human

Workload or System role holder
  -> its current responsible Human

AccessGroup role holder
  -> the current group owner
  + every current direct Human member
  + the current responsible Human of every current direct Workload or System
    member
```

The expansion is bounded to current direct AccessGroup Memberships; nested and
dynamic groups remain excluded by F18-RD-04. A replacement uses the union of
the sets derived from both the reviewed assignment and the complete proposed
replacement. The acting reviewer must be a current Human PrincipalRef outside
that union.

For Retain and Replace, reviewer eligibility and the complete
`ReviewBeneficiaryHumanSet` are re-evaluated when the disposition is accepted
and at the same atomic decision/effect publication boundary that publishes
Retain certification and due advancement or Replace remediation. Missing,
stale, ambiguous, or unresolved group ownership, current direct Membership, or
responsible-Human evidence needed to derive the set fails closed with
`REVIEW_BENEFICIARY_UNRESOLVED`; it publishes no authority effect and does not
advance `nextReviewDueAt`. Revoke remains an independently authorized
privilege-reducing action and is not blocked merely because the reviewer is a
beneficiary or the beneficiary set cannot be completely resolved. A
Retain/Replace conflict returns `REVIEWER_CONFLICT`, publishes neither a final
authority effect nor `nextReviewDueAt`, and leaves the item available for a
different independently eligible reviewer. It never silently revokes or
retains access. A conflict discovered before publication cannot be overridden
by prior eligibility, a stale reviewer snapshot, approval, or an exception.
Membership-only review applies the equivalent Human expansion where retaining
the reviewed relationship would preserve access eligibility.

A `Retain` disposition for a Workload or System Membership or RoleAssignment
requires its responsiblePartyRef to resolve to a current Human PrincipalRef at
decision time. Failed responsibility validation cannot certify the item and
must produce a non-Retain disposition or leave the item undecided; it never
silently revokes an otherwise-effective production assignment.

Each item requires exactly one disposition from one currently eligible Human
PrincipalRef. The first valid disposition wins under F18-RD-21; replay of the
same decision is idempotent and a different decision conflicts. A campaign
becomes Completed only after every snapshotted item has one final disposition.
Stale completes that item but never certifies or remediates it. Any undecided
item at `dueAt` makes the campaign ExpiredIncomplete.

Privileged and guest access expires independently. An incomplete standing
production Workload/System review becomes overdue and escalated but is not
automatically revoked solely for missing telemetry or incomplete review unless
a prepublished policy requires suspension. Mode-specific due dates derive from
pinned rule evidence; the twenty-four-hour break-glass ceiling is the only v1
platform constant. Human and AccessGroup-held Standing
assignments become authorization-ineligible at their review due time until a
completed exact-version `StandingCertification + Retain` or an authorized
replacement exists. Manual and BreakGlassRetrospective Retain decisions never
restore Standing eligibility. A mandatory break-glass review that becomes
overdue emits protected escalation evidence and remains an unresolved security
finding; it cannot revive or extend the already expiring temporary assignment.

Access-review rule timing is mode-specific. `StandingCertification` requires
`reviewInterval`, `reviewWindow`, and `overdueEffect`, with
`0 < reviewWindow <= reviewInterval`. `Manual` requires `reviewWindow` only.
`BreakGlassRetrospective` carries no independent interval or window; its due
time comes from the exact privileged-access rule as specified above. Common
reviewer-eligibility and usage-evidence requirements remain present where the
mode requires them. On entering the window before a Standing assignment's
`nextReviewDueAt`, the review controller materializes at most one campaign bound
to that assignment, rule, exact assignment due time, and snapshot subject;
scheduler invocation carries no write authority. Same-boundary replay returns
the existing campaign. Failure to materialize never extends the assignment's
due time or effectiveness.

A completed `StandingCertification` campaign with `Retain` for the exact current
assignment version authorizes the review controller to submit the exact due-
advancement intent. The RoleAssignment controller alone sets the reviewed
Standing RoleAssignment's system-owned `status.nextReviewDueAt` to
`campaign.completedAt + reviewInterval`. A Manual or
BreakGlassRetrospective campaign never advances `nextReviewDueAt`. Replace
produces a replacement whose initial due time derives from its publication time
and pinned rule only after the RoleAssignment controller accepts and publishes
the exact intent. Review completion never extends an expired temporary
assignment. Certification and due-time advancement publish under the same local
audit/concurrency boundary.

### F18-RD-17 — Bounded immutable exception evidence

No ExceptionRequest resource exists. The flow is immutable ExceptionProposal
inside ApprovalRequest followed by approved ExceptionGrant. The grant pins
exact control, subject, scope, bounded override, interval, justification,
compensating controls, approval/decision evidence, and linkage.

Every accepted ExceptionProposal remains pending only while its exact
ApprovalRequest lifecycle and requested exception interval permit. A current
Approved request plus
successful final control, subject, scope, interval, overlap, compensating-
control, authorization, and audit validation produces `Grant`. ApprovalRequest
`Denied`, `Cancelled`, or `Expired` produces `Deny` with `ApprovalDenied`,
`ApprovalCancelled`, or `ApprovalExpired`, respectively. Final proposal
validation also produces terminal `Deny` with exactly one stable reason from
`InvalidProposal | UnsupportedControl | SubjectNotAllowed | ScopeNotAllowed |
DurationExceeded | OverlapConflict | Unauthorized | ApprovalDenied |
ApprovalCancelled | ApprovalExpired`. `ConcurrencyConflict | DependencyUnavailable |
AuditObligationUnavailable` are retryable operation failures: they publish no
terminal ExceptionProposal result and retry uses the same immutable proposal,
concurrency, and idempotency identity. If a retryable failure remains unresolved
at the immutable ApprovalRequest `expiresAt`, the proposal becomes terminal
`Deny` with `ApprovalExpired`. No unregistered reason changes terminality.
These are exception-controller decisions about the exact proposal and do not
fabricate or reinterpret an ApprovalRequest decision.

`subjectRef` is exactly one PrincipalRef, RoleAssignmentRef,
PrivilegedAccessRequestRef, or UID-pinned canonical ResourceRef whose kind is
registered by the control-owning definition. Adding another allowed canonical
subject kind requires a new control-definition version and conformance evidence,
not a core conditional. The versioned definition owning a control declares its
exact identity/version, exceptionability, allowed subject kinds, typed bounded-
override schema, required compensating-control kinds, and maximum duration.
`boundedOverride` must validate against that exact schema; an untyped map,
arbitrary expression, unknown field, or override of an undeclared control
dimension rejects. FEATURE-0018 consumes but does not redefine that control.
The platform ceiling is ninety days; policy may shorten it. Authentication
integrity, tenant/scope isolation, audit integrity, immutable evidence, and
unresolved conflict are non-exceptionable.

Final constraint composition is exact:

```text
maximumDuration = minimum(all applicable duration ceilings)
requiredCompensatingControls = union(all mandatory control sets)
allowedSubjects = intersection(all registered subject constraints)
allowedScopes = intersection(all registered scope constraints)
allowedOverrides = intersection(all registered override constraints)
```

An approved proposal must already satisfy the composed result exactly. The
exception controller rejects a nonconforming proposal and never silently
shortens its duration, removes an override, changes its subject or scope, or
adds compensating controls on the requester's behalf.

Subject/scope compatibility is exact. For a RoleAssignmentRef,
PrivilegedAccessRequestRef, or canonical ResourceRef subject,
`ExceptionGrant.metadata.scopeRef` must equal the referenced subject's exact
canonical governance scope. For a PrincipalRef subject, the declared
ExceptionGrant scope is authoritative and F18-RD-09's contextual Membership
requirements apply. An Organization-scoped exception does not implicitly cover
descendants; CloudProvider, CloudPlatform, and Platform exceptions are also
exact-scope only. The control-owning version registers allowed ScopeKinds in
addition to allowed subject kinds, typed override schema, maximum duration, and
required compensating-control kinds. Policy may narrow but not broaden them.
FEATURE-0018 performs no ancestor/descendant exception application; FEATURE-
0020 remains the sole effective resolver.

ExceptionGrant is append-only and uses the half-open validity interval
`[notBefore, expiresAt)`. Two grants overlap when they pin the same exact control
identity and version, UID-pinned subject, and canonical scope and their effective
intervals have a non-empty intersection. A new overlap is rejected regardless
of different justification, compensating controls, or bounded override; no
precedence or merge is inferred.

Early revocation is linked same-kind immutable evidence with `effect=Revoke`,
exact `revokedGrantRef`, and authoritative `revokedAt`. A valid revocation
before `expiresAt` shortens the original grant's effective interval without
mutating it. When `revokedAt <= notBefore`, the effective interval is empty;
otherwise it is `[notBefore, min(expiresAt, revokedAt))`. A revocation at or
after expiry is rejected as already expired. Overlap evaluation excludes time
after effective revocation, and a new grant may begin exactly at the prior
`expiresAt` or effective `revokedAt`. Concurrent issue/revoke follows F18-RD-21
and required publication follows F18-RD-20, so conflicting grants cannot both
become active. An exception never grants a role. FEATURE-0018 owns issuance,
exact revocation, and local non-overlap; FEATURE-0020 alone owns effective
application and cross-profile resolution.

ExceptionGrant has no mutable lifecycle or status. Its effectiveness is computed
from authoritative time, the immutable half-open interval, and any valid linked
revocation evidence. At effective expiry, the owning controller emits the
registered protected AuditEvent and may refresh a non-authoritative read model or
cache; it never mutates the grant, creates replacement evidence, or makes the
projection an authorization authority. Audit or projection failure cannot extend
the exception's effective interval.

### F18-RD-21 — Deterministic validation and safe denial

Validation precedence is exact:

```text
bounded transport
  -> authentication
  -> structural validation
  -> local semantic and cross-field validation
  -> safe scope/reference resolution
  -> operation-required assurance and contextually required current membership
  -> assignment/action/scope/resource/delegation and membership-enabled-grant-envelope evaluation
  -> operation-required EligibilityRef and separation-of-duties evaluation
  -> approval-requirement/policy/approval/exception evidence evaluation
  -> concurrency and idempotency
  -> required decision/audit-obligation evidence
  -> publication
```

Unknown, deferred, and future-owned fields reject. Safe denial occurs before
detailed inaccessible-reference disclosure. Tokens, assertions, provider
claims, confidential policy input, sensitive justification, credentials,
secrets, and evaluator diagnostics never enter unsafe errors, logs, or
projections.

Idempotency binds actor, operation, exact target, scope, and request digest.
Replay rechecks current authentication, authorization, and safe target
visibility. Same key/different content conflicts. Required audit failure
publishes no state change or completed replay result.

### F18-RD-24 — Progressive and normally hidden user friction

Routine authorization resolves silently. Users do not manage schema versions,
policy digests, issuer/subject, provenance, external object refs, exact role or
policy versions, authorization values, evidence refs, or correlation IDs.
These are system-owned.

The six user journeys are:

1. routine end-customer operation: business input and result or safe denial;
2. ordinary assignment: holder, named role, scope, explicit validity, and one
   status when system policy requires approval;
3. privileged request: role, scope, duration, justification, one request status,
   required approvers, readiness, and expiry;
4. approver decision: actionable requester/access/scope/duration/risk/quorum
   context and approve/deny action;
5. access review: holder/role/scope/validity/qualified usage summary and
   Retain/Revoke/Replace action; and
6. exception request: exact control, subject/scope, duration, justification,
   compensating controls, one request status, and effective period.

PrivilegedAccessRequest, ApprovalRequest, temporary RoleAssignment, and
evidence are one requester-facing journey. ExceptionProposal, ApprovalRequest,
and ExceptionGrant are also one requester-facing journey. Separate protected
resources remain internally auditable; no facade resource is added.
