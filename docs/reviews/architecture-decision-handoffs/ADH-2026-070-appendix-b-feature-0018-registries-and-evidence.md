# ADH-2026-070 Appendix B — FEATURE-0018 registries and evidence contract

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

Appendix B carries F18-RD-02, F18-RD-06, F18-RD-18 through
F18-RD-20, F18-RD-22, and F18-RD-23. Appendix A carries the remaining domain
semantic groups. All cross-group references resolve through the core's exact
decision index.

The decision text originated in the former single-file ADH-2026-070 payload and
includes the human-approved F18-SEC-001 through F18-SEC-005 corrections applied
after the first independent security review. It must be reviewed and approved
as part of this exact renewed three-file package.

## Normative decision groups

### F18-RD-02 — Closed contract inventory and profiles

FEATURE-0018 owns this closed inventory:

| Contract | Resource profile | Sole responsibility |
|---|---|---|
| `PrincipalRef` | EmbeddedValue | Stable reference to an already-authenticated Human, Workload, or System identity; identity lifecycle remains external |
| `AccessGroup` | ManagedResource | Non-authenticating, access-only authorization subject with explicit ownership and direct auditable principal membership |
| `Membership` | ManagedResource | Associates a principal with an Organization or CloudProvider belonging context, or records direct membership in a same-scope AccessGroup relationship; grants no permission by itself |
| `RoleDefinition` | VersionedDefinition | Defines an immutable published version of canonical Sovrunn actions |
| `RoleAssignment` | ManagedResource | Grants one published role version to one role holder at one canonical scope, optionally narrowed to one exact resource, with explicit validity |
| `PrivilegedAccessRequest` | LongRunningOperation | Requests justified, approved, time-bounded privileged access |
| `AccessReview` | LongRunningOperation | Reviews an immutable snapshot of access and retains its findings/remediation evidence |
| `ApprovalPolicy` | VersionedDefinition | Defines bounded approval stages, eligibility, quorum, expiry, and separation of duties |
| `ApprovalRequest` | LongRunningOperation | Retains approval of one bounded FEATURE-0018 request or immutable proposal |
| `ExceptionGrant` | ImmutableRecord | Retains an approved, scoped, time-bounded exception or its linked revocation evidence |
| `GovernanceProfile` | VersionedDefinition | Composes only FEATURE-0018-owned governance references and constraints in v1 |

The table above is the closed canonical output-contract inventory; only rows
with persistent resource profiles are independently addressable. FEATURE-0018
also owns this closed supporting-value inventory; none of these values has an
independent endpoint, lifecycle, storage authority, or user journey:

| Supporting contract | Resource profile | Sole responsibility |
|---|---|---|
| `RoleHolderRef` | EmbeddedValue | Closed union of exactly one `PrincipalRef` or kind-constrained `AccessGroupRef` |
| `AssuranceEvidence` | TransientRequestResult | Operation-local already-validated authentication-assurance carrier |
| `AuthorizationInput` | TransientRequestResult | Exact actor, action, exact-resource target, create-parent scope, or scope-only target, time, assurance, request, and correlation input |
| `AuthorizationResult` | TransientRequestResult | Non-bearer `Allow | Deny` outcome with exact contributing provenance |
| `ApprovalRequirement` | TransientRequestResult | System-derived operation-local closed union of `NotRequired` or `Required` with one exact published `ApprovalPolicy` version ref |
| `EligibilityRef` | EmbeddedValue | Kind-constrained reference containing exactly one Human PrincipalRef; no role, assignment, group, Membership, claim, or other relationship can satisfy it |
| `RoleAssignmentProposal` | EmbeddedValue | Immutable proposed assignment retained only when required by an accepted flow |
| `MembershipEnabledGrantEnvelope` | TransientRequestResult | Operation-local collection of exact group-held assignment action/reach/validity tuples newly enabled or extended by an AccessGroup Membership assignment-effect expansion |
| `ExceptionProposal` | EmbeddedValue | Immutable proposed exception retained only when required by an accepted flow |
| `UsageEvidenceSummary` | EmbeddedValue | Immutable review-only usage summary with protected source evidence |
| `GovernanceApplicability` | EmbeddedValue | Closed applicability tuple shared by the permitted FEATURE-0018 rule components |
| `ActionTargetBinding` | EmbeddedValue | Closed discriminated union binding one canonical action to exactly one ExactResource target kind, CreateParent parent ScopeKind, or ScopeOnly ScopeKind |

The exact initial supporting enums are:

```text
MembershipContextKind = Organization | CloudProvider | AccessGroup
MembershipType = Standard | Guest
PrivilegedAccessMode = JIT | BreakGlass
ApprovalDecision = Approve | Deny
ExceptionGrantEffect = Grant | Revoke
AssignmentValidityMode = Standing | TimeBound
ActionTargetMode = ExactResource | CreateParent | ScopeOnly
```

These names are canonical contract values, not personas, workflow states,
provider-native values, or extensible strings. Adding a value requires
versioned contract change and conformance evidence.

Scope applicability is closed and deterministic per versioned canonical
resource contract. Core authorization and validation consume governed scope-
applicability registrations generically and do not hard-code FEATURE-0018
resource kinds. Customer policy may narrow but never expand a registered
contract. Adding an existing canonical ScopeKind to a resource requires
contract-owner approval, compatibility assessment, a versioned contract update,
and updated conformance evidence, but should not require a core-engine change.
Adding a new ScopeKind or changing scope containment is a core architecture
change. Configuration, plugins, adapters, external claims, and generic extension
fields cannot add a ScopeKind or expand registered applicability.

The initial FEATURE-0018 contract version registers exactly:

| Contract | Allowed canonical `metadata.scopeRef` kinds |
|---|---|
| `PrincipalRef` | None; embedded values have no independent scope |
| `AccessGroup` | Organization, CloudProvider |
| `Membership` | Organization, CloudProvider |
| `RoleDefinition` | Platform, Organization, CloudProvider |
| `RoleAssignment` | Platform, CloudPlatform, CloudProvider, Organization, OrganizationUnit, Tenant, Project |
| `PrivilegedAccessRequest` | Platform, CloudPlatform, CloudProvider, Organization, OrganizationUnit, Tenant, Project |
| `AccessReview` | Platform, CloudPlatform, CloudProvider, Organization, OrganizationUnit, Tenant, Project |
| `ApprovalPolicy` | Platform, Organization, CloudProvider |
| `ApprovalRequest` | Platform, CloudPlatform, CloudProvider, Organization, OrganizationUnit, Tenant, Project |
| `ExceptionGrant` | Platform, CloudPlatform, CloudProvider, Organization, OrganizationUnit, Tenant, Project |
| `GovernanceProfile` | Platform, Organization, CloudProvider |

Target-bound contracts use the exact target governance scope; the registry
does not make every target kind valid at every listed scope. Existing
canonical target/scope and containment rules still apply, and `resourceRef`
may only narrow within the registered `metadata.scopeRef`.

Definition publication scope and reference compatibility are separate from
authorization containment and never grant an action:

| Published definition scope | Compatible consuming scope |
|---|---|
| Platform | Any scope registered for the consuming FEATURE-0018 contract |
| Organization | The same Organization or its OrganizationUnit, Tenant, or Project descendants |
| CloudProvider | The same exact CloudProvider only |

Accordingly, a Platform-scoped RoleDefinition may be assigned at any registered
RoleAssignment scope, an Organization-scoped RoleDefinition only within its
Organization tree, and a CloudProvider-scoped RoleDefinition only at that exact
provider. The same rule validates an exact selected ApprovalPolicy against an
ApprovalRequest. A GovernanceProfile may reference an ApprovalPolicy at its own
scope or at Platform; a Platform GovernanceProfile may reference only a
Platform policy. CloudPlatform-specific RoleDefinition, ApprovalPolicy, and
GovernanceProfile publication is not part of the initial contract; a
CloudPlatform-scoped consumer uses a compatible Platform definition. Policy may
narrow these relationships but never broaden them. Undefined or cross-tree
reference compatibility fails closed.

No `ExceptionRequest`, workflow, facade, generic condition, or external-group
resource is added.

System-owned field authority is separated as follows:

| Contract | Accepted-intent or immutable-spec writer | Status, result, or record writer |
|---|---|---|
| `AccessGroup` | Authorized access administrator | Access-group controller |
| `Membership` | Membership administrator or trusted provisioner submits an exact intent and no canonical state | Identity/membership controller; it alone validates any membership-enabled grant envelope and writes specification, protected provenance/freshness, lifecycle and eligibility projection |
| `RoleDefinition` | Security administrator | Publication controller |
| `RoleAssignment` | Administrators and grant, approval, privileged-access, or review controllers submit only exact immutable intents; RoleAssignment controller alone writes the immutable specification | RoleAssignment controller alone writes lifecycle, status, revocation, replacement, expiry, and review-due effects |
| `PrivilegedAccessRequest` | Requester through accepted-request boundary | Privileged-access controller |
| `ApprovalPolicy` | Security or governance administrator | Publication controller |
| `ApprovalRequest` | Originating accepted FEATURE-0018 intent; Approval-request controller materializes the immutable request | Approval controller |
| `AccessReview` | Authorized campaign owner submits Manual intent; schedulers or the privileged-access controller submit mandatory campaign intents; review controller alone materializes the canonical campaign | Review controller |
| `ExceptionGrant` | Exception controller after approval | Immutable record; no mutable status writer |
| `GovernanceProfile` | Governance administrator | Publication controller |

Automatic version supersession is written only by the publication controller
that successfully publishes the successor. It is not a caller action or a
second definition writer.

The RoleAssignment controller is the sole canonical RoleAssignment publisher
and sole writer of its specification, lifecycle, status, revocation,
replacement, expiry, and review-due effects. An administrator or originating
grant, approval, privileged-access, or review workflow may authorize and submit
only an exact immutable intent; it cannot directly create, modify, revoke,
replace, expire, or publish any part of a RoleAssignment.

The closed v1 operation and mutability surface is:

| Contract | Caller-visible or authorized-administrator operations | Owning-controller-only effects |
|---|---|---|
| `AccessGroup` | create, get, list, update permitted owner/display intent, suspend, resume, retire | status publication |
| `Membership` | create/provision, get, list, synchronize, suspend, reactivate, revoke | derive and validate the exact membership-enabled grant envelope for every AccessGroup assignment-effect expansion; publish protected freshness/expiry audit and non-authoritative assignment-effect projection |
| `RoleDefinition` | create Draft, get, list, update Draft, publish, suspend, restore, retire | automatic supersession linkage |
| `RoleAssignment` | submit proposal or direct grant, get, list, revoke | validate and publish the exact approved grant, revocation, expiry, and review-eligibility effects |
| `PrivilegedAccessRequest` | submit, get, list, cancel own, revoke own, administratively revoke under the distinct privileged action | map approval, activate, block, and expire |
| `ApprovalPolicy` | create Draft, get, list, update Draft, publish, retire | automatic supersession linkage |
| `ApprovalRequest` | get, list, decide, cancel originating pending intent, administratively cancel under the distinct privileged action | create from accepted intent, expire, and publish result |
| `AccessReview` | create authorized campaign, get, list, start, decide item, remediate, cancel | capture snapshot, publish readiness, complete, and expire incomplete campaign |
| `ExceptionGrant` | submit proposal, get, list, revoke | issue immutable grant/revocation evidence; emit expiry audit and refresh only a non-authoritative read projection without mutating canonical evidence |
| `GovernanceProfile` | create Draft, get, list, update Draft, publish, retire | automatic supersession linkage |

No FEATURE-0018 contract supports generic `PUT`, unrestricted `PATCH`, hard
deletion, or caller-authored status. FEATURE-0012 domain-grouped versioned route
rules govern design of literal HTTP paths, but design cannot add an operation,
change its authority, or convert a controller effect into caller permission.

Schedulers may invoke an owning controller but never become field writers.
System-owned provenance, exact versions, evidence references, correlation
identifiers, lifecycle results, and authorization values reject user
authorship.

### F18-RD-06 — Distributed action ownership and central role composition

The feature owning an operation owns its canonical action meaning.
FEATURE-0018 validates registered actions and composes them into roles without
renaming, splitting, or reinterpreting them. FEATURE-0016's
`executiontarget.read`, `executiontarget.qualify`, and
`executiontarget.write` are consumed unchanged. The accepted create/retire
granularity of `executiontarget.write` remains an explicit dependency
limitation. The action-owning feature also owns any intrinsic privileged-
handling classification attached to its registered action. FEATURE-0018
consumes that classification without weakening or reinterpreting it; a missing
or ambiguous required classification is fail-safe Privileged under F18-RD-07.

FEATURE-0018 authorizes canonical Sovrunn control-plane actions only. Provider-
native principals, groups, roles, policies, permissions, and credentials are
neither canonical FEATURE-0018 resources nor FEATURE-0018 outputs.
FEATURE-0018 does not create, publish, synchronize, translate, or expose them.

When a separately owned execution or integration feature carries an authorized
operation to an ExecutionTarget, FEATURE-0018 authorization is necessary but
not sufficient. The responsible controller or adapter must use its own least-
privileged native identity, and the ExecutionTarget's native IAM must
independently authorize the native effect. Sovrunn authorization and native IAM
authorization therefore compose by intersection: both must permit the
operation. An Allow from either layer cannot substitute for authorization by
the other, and neither layer may override a Deny from the other.

Direct, out-of-band administration of the IaaS environment remains
CloudProvider-owned, is outside FEATURE-0018 authority, and must not be
represented as an effect of a Sovrunn RoleAssignment. FEATURE-0018 defines only
this trust boundary; native identity acquisition, credential management,
permission translation, and effect execution remain responsibilities of their
separately approved owning features.

The table below is the exact initial assignable action registry consumed by
FEATURE-0018. FEATURE-0018 owns the meaning of its own contract rows;
FEATURE-0016 remains the sole semantic and registration owner for the three
ExecutionTarget rows. `Ordinary` means the action may appear in an ordinary
role but still requires an independently applicable assignment and every
contextual guardrail; `Privileged` forces the individual TimeBound path under
F18-RD-07. Lifecycle actions are distinct to preserve least privilege.

| Contract | Canonical action | Operation meaning | Intrinsic class |
|---|---|---|---|
| AccessGroup | `accessgroup.read` | Read an authorized safe projection | Ordinary |
| AccessGroup | `accessgroup.write` | Create or update allowed owner/display intent | Privileged |
| AccessGroup | `accessgroup.suspend`, `accessgroup.resume`, `accessgroup.retire` | Request the exact named lifecycle transition | Privileged |
| Membership | `membership.read` | Read an authorized safe projection | Ordinary |
| Membership | `membership.write` | Create trusted belonging/group intent or refresh synchronized provenance; an AccessGroup assignment-effect expansion additionally passes F18-RD-09's membership-enabled grant ceiling | Privileged |
| Membership | `membership.suspend`, `membership.reactivate`, `membership.revoke` | Request the exact named lifecycle transition; AccessGroup reactivation additionally passes F18-RD-09's membership-enabled grant ceiling | Privileged |
| RoleDefinition | `roledefinition.read` | Read an authorized safe definition projection | Ordinary |
| RoleDefinition | `roledefinition.write` | Create or update Draft intent | Privileged |
| RoleDefinition | `roledefinition.publish`, `roledefinition.suspend`, `roledefinition.restore`, `roledefinition.retire` | Request the exact named lifecycle transition | Privileged |
| RoleAssignment | `roleassignment.read` | Read an authorized safe grant projection | Ordinary |
| RoleAssignment | `roleassignment.grant` | Submit one direct non-delegable grant intent or one RoleAssignmentProposal within the grantor ceiling | Privileged |
| RoleAssignment | `roleassignment.delegate` | Additionally authorize granting a role containing `roleassignment.grant` or `roleassignment.delegate` | Privileged |
| RoleAssignment | `roleassignment.revoke` | Revoke one exact assignment | Privileged |
| PrivilegedAccessRequest | `privilegedaccessrequest.read` | Read an authorized requester/approver-safe projection | Ordinary |
| PrivilegedAccessRequest | `privilegedaccessrequest.submit`, `privilegedaccessrequest.cancel`, `privilegedaccessrequest.revoke` | Submit, cancel, or revoke only the actor's own bounded request/grant | Ordinary |
| PrivilegedAccessRequest | `privilegedaccessrequest.adminrevoke` | Revoke another principal's exact active temporary grant | Privileged |
| ApprovalPolicy | `approvalpolicy.read` | Read an authorized safe definition projection | Ordinary |
| ApprovalPolicy | `approvalpolicy.write` | Create or update Draft intent | Privileged |
| ApprovalPolicy | `approvalpolicy.publish`, `approvalpolicy.retire` | Request the exact named lifecycle transition | Privileged |
| ApprovalRequest | `approvalrequest.read`, `approvalrequest.decide` | Read an eligible projection or submit one eligible decision | Ordinary |
| ApprovalRequest | `approvalrequest.cancel` | Cancel only the actor's originating pending intent | Ordinary |
| ApprovalRequest | `approvalrequest.admincancel` | Administratively cancel another actor's exact pending intent | Privileged |
| AccessReview | `accessreview.read`, `accessreview.decide` | Read an eligible projection or submit one eligible item decision | Ordinary |
| AccessReview | `accessreview.create` | Create one bounded Manual campaign intent using F18-RD-16's explicit item-selection rules | Privileged |
| AccessReview | `accessreview.start`, `accessreview.remediate`, `accessreview.cancel` | Start, apply approved remediation, or cancel one campaign | Privileged |
| ExceptionGrant | `exceptiongrant.read`, `exceptiongrant.propose` | Read an authorized safe projection or submit one bounded proposal | Ordinary |
| ExceptionGrant | `exceptiongrant.revoke` | Issue linked revocation evidence for one exact grant | Privileged |
| GovernanceProfile | `governanceprofile.read` | Read an authorized safe definition projection | Ordinary |
| GovernanceProfile | `governanceprofile.write` | Create or update Draft intent | Privileged |
| GovernanceProfile | `governanceprofile.publish`, `governanceprofile.retire` | Request the exact named lifecycle transition | Privileged |
| ExecutionTarget (FEATURE-0016-owned) | `executiontarget.read` | Consume the existing FEATURE-0016 operation meaning unchanged | Ordinary |
| ExecutionTarget (FEATURE-0016-owned) | `executiontarget.qualify` | Consume the existing FEATURE-0016 operation meaning unchanged | Privileged |
| ExecutionTarget (FEATURE-0016-owned) | `executiontarget.write` | Consume the existing FEATURE-0016 operation meaning unchanged | Privileged |

The initial target-binding registrations below are exact. `OrgProvider` means
only `Organization | CloudProvider`; `DefinitionScopes` means only `Platform |
Organization | CloudProvider`; and `AllResourceScopes` means only `Platform |
CloudPlatform | CloudProvider | Organization | OrganizationUnit | Tenant |
Project`. These are table abbreviations, not additional canonical values.

| Canonical action or exact action set | `ExactResource` target kind | `CreateParent` parent ScopeKinds | `ScopeOnly` ScopeKinds |
|---|---|---|---|
| `accessgroup.read` | AccessGroup | — | OrgProvider |
| `accessgroup.write` | AccessGroup | OrgProvider | — |
| `accessgroup.suspend`, `accessgroup.resume`, `accessgroup.retire` | AccessGroup | — | — |
| `membership.read` | Membership | — | OrgProvider |
| `membership.write` | Membership | OrgProvider | — |
| `membership.suspend`, `membership.reactivate`, `membership.revoke` | Membership | — | — |
| `roledefinition.read` | RoleDefinition | — | DefinitionScopes |
| `roledefinition.write` | RoleDefinition | DefinitionScopes | — |
| `roledefinition.publish`, `roledefinition.suspend`, `roledefinition.restore`, `roledefinition.retire` | RoleDefinition | — | — |
| `roleassignment.read` | RoleAssignment | — | AllResourceScopes |
| `roleassignment.grant`, `roleassignment.delegate` | — | AllResourceScopes | — |
| `roleassignment.revoke` | RoleAssignment | — | — |
| `privilegedaccessrequest.read` | PrivilegedAccessRequest | — | AllResourceScopes |
| `privilegedaccessrequest.submit` | — | AllResourceScopes | — |
| `privilegedaccessrequest.cancel`, `privilegedaccessrequest.revoke`, `privilegedaccessrequest.adminrevoke` | PrivilegedAccessRequest | — | — |
| `approvalpolicy.read` | ApprovalPolicy | — | DefinitionScopes |
| `approvalpolicy.write` | ApprovalPolicy | DefinitionScopes | — |
| `approvalpolicy.publish`, `approvalpolicy.retire` | ApprovalPolicy | — | — |
| `approvalrequest.read` | ApprovalRequest | — | AllResourceScopes |
| `approvalrequest.decide`, `approvalrequest.cancel`, `approvalrequest.admincancel` | ApprovalRequest | — | — |
| `accessreview.read` | AccessReview | — | AllResourceScopes |
| `accessreview.create` | — | AllResourceScopes | — |
| `accessreview.decide`, `accessreview.start`, `accessreview.remediate`, `accessreview.cancel` | AccessReview | — | — |
| `exceptiongrant.read` | ExceptionGrant | — | AllResourceScopes |
| `exceptiongrant.propose` | — | AllResourceScopes | — |
| `exceptiongrant.revoke` | ExceptionGrant | — | — |
| `governanceprofile.read` | GovernanceProfile | — | DefinitionScopes |
| `governanceprofile.write` | GovernanceProfile | DefinitionScopes | — |
| `governanceprofile.publish`, `governanceprofile.retire` | GovernanceProfile | — | — |
| `executiontarget.read` | ExecutionTarget | — | CloudProvider |
| `executiontarget.qualify` | ExecutionTarget | — | — |
| `executiontarget.write` | ExecutionTarget | CloudProvider | — |

Every canonical action registers one or more exact discriminated target-binding
variants:

```text
ExactResource { action, targetKind }
CreateParent  { action, parentScopeKind }
ScopeOnly     { action, scopeKind }
```

Exactly one variant applies to each registration. Existing-resource operations
use `ExactResource`; collection operations use `ScopeOnly`; and creation,
submission, proposal, and grant operations use `CreateParent`. One action may
register multiple explicit variants where its accepted operation meaning
requires them. A mixed variant, missing discriminator or required kind, empty,
wildcard, unknown, or ambiguous registration fails closed. FEATURE-0018
consumes action-owner-supplied registrations generically and does not infer
target behavior or hard-code another feature's resource kinds. FEATURE-0016
must publish the three exact ExecutionTarget registrations and intrinsic
classes above through its owning extension/reconciliation process; until then,
an absent registration makes that action unassignable and unusable in
FEATURE-0018 rather than permitting FEATURE-0018 to invent or reinterpret it.

For ApprovalPolicy applicability, the canonical action is exactly
`privilegedaccessrequest.submit`, `roleassignment.grant` (and additionally
`roleassignment.delegate` when applicable), or `exceptiongrant.propose` for the
three allowed subject kinds. Expiry and controller reconciliation are automatic
effects, not assignable actions. Controllers in F18-RD-02 may materialize only
the exact effect authorized by accepted intent and current evidence; controller
identity never permits selection or enlargement of holder, role, scope,
duration, decision, review remediation, or exception.

### F18-RD-18 — FEATURE-0018-limited GovernanceProfile v1

GovernanceProfile remains the versioned compositional envelope established by
ADH-2026-033, not an assignment or effective context. FEATURE-0018 v1 activates
only:

- `approvalPolicyRefs`;
- `privilegedAccessRules`;
- `accessReviewRules`;
- `exceptionRules`; and
- `auditRequirements`.

ApprovalPolicy refs pin exact published versions. Applicability keys are unique
and ambiguous overlaps reject. No component is syntactically mandatory.
Generic maps and speculative sovereignty, placement, backup, cost,
entitlement, quota, or execution fields are prohibited in the FEATURE-0018 v1
contract. Future domain semantics require separate architecture change
control.

Every component entry uses the closed applicability tuple `subjectKind +
canonicalAction + ScopeKind + optional targetKind + optional requestMode` and a
stable `ruleId`; a field is present only where its component contract permits
it. Exact duplicate or overlapping tuples within one profile reject rather than
merge. The five component schemas are closed as follows:

| Component | Exact v1 semantic fields |
|---|---|
| `approvalPolicyRefs` | applicability tuple + exact published ApprovalPolicy version ref |
| `privilegedAccessRules` | applicability tuple + allowed exact RoleDefinition version selectors + non-empty `EligibilityRef` requester list + minimum assurance level + maximum assurance age + maximum JIT duration + break-glass allowed boolean + maximum break-glass duration + phishing-resistance requirement (constant true for JIT/BreakGlass) + retrospective-review deadline |
| `accessReviewRules` | applicability tuple + non-empty `EligibilityRef` reviewer list + usage-evidence freshness/coverage requirement + exactly one mode-specific variant: `StandingCertification { reviewInterval, reviewWindow, overdueEffect Ineligible | Escalate, allowedStandingHolderKinds, allowedStandingScopeKinds, allowedStandingTargetKinds }`; `Manual { reviewWindow }`; or `BreakGlassRetrospective {}` whose deadline is owned by the applicable privileged-access rule |
| `exceptionRules` | applicability tuple + exact exceptionable control version refs + rule maximum duration + required compensating-control kinds |
| `auditRequirements` | applicability tuple + registered required AuditEvent types + DecisionRecord requirement `Never | MaterialDeniedPrivileged | Always` + registered retention/projection class refs |

The initial component-specific applicability registrations are closed:

| Component | Allowed `subjectKind + canonicalAction + requestMode` combinations |
|---|---|
| `approvalPolicyRefs` | `PrivilegedAccessRequest + privilegedaccessrequest.submit` with `JIT` or `BreakGlass`; `RoleAssignmentProposal + roleassignment.grant + absent` and additionally `roleassignment.delegate + absent` when delegation is proposed; `ExceptionProposal + exceptiongrant.propose + absent` |
| `privilegedAccessRules` | `PrivilegedAccessRequest + privilegedaccessrequest.submit` with `JIT` or `BreakGlass`; `RoleAssignmentProposal + roleassignment.grant + absent` and additionally `roleassignment.delegate + absent` when delegation is proposed |
| `accessReviewRules` | `RoleAssignment + roleassignment.grant + absent` for StandingCertification; `PrivilegedAccessRequest + privilegedaccessrequest.submit + BreakGlass` for BreakGlassRetrospective; `AccessReview + accessreview.create + absent` for Manual |
| `exceptionRules` | `ExceptionProposal + exceptiongrant.propose + absent` only |
| `auditRequirements` | A `subjectKind + canonicalAction` pair already registered by F18-RD-06, including the three reused FEATURE-0016 ExecutionTarget actions; `subjectKind` is the action-owning canonical contract kind, and `requestMode` is absent except `JIT` or `BreakGlass` for `PrivilegedAccessRequest + privilegedaccessrequest.submit` |

For every row, ScopeKind must be registered for the governed operation and
subject. `targetKind` is present only when the exact operation contract permits
an exact resource target; it is absent for a heterogeneous Manual review.
Unknown, cross-owner, future-action, or otherwise unregistered combinations
reject. Extending these combinations requires a versioned FEATURE-0018 contract
change and conformance evidence; configuration, a profile instance, a plugin,
or external input cannot register a combination.

Rule values may only narrow the platform ceilings and mandatory controls in
F18-RD-12 through F18-RD-20. `Ineligible` is mandatory for overdue Human and
AccessGroup-held Standing access; production Workload/System rules may select
`Ineligible` or `Escalate`. A syntactically optional component becomes
operationally required when an operation depends on that rule: absence,
ambiguity, scope incompatibility, or an unsupported value fails closed. These
schemas define local validation only; FEATURE-0020 still exclusively selects
and resolves effective profile evidence.

`auditRequirements.DecisionRecordRequirement` may require DecisionRecords in
additional cases. Its `Never` value means no policy-added DecisionRecord and
cannot suppress any DecisionRecord made mandatory by F18-RD-19.

Effective lifecycle is `Draft -> Published -> Retired`; published versions are
immutable. F18-RD-07's common informational supersession relationship applies,
so a superseded profile receives no new selections while independently valid
existing exact-version pins remain usable. Retirement is terminal: a retired
version receives no new selection and is unavailable to an operation requiring
current rule evidence. Already-published effects continue under their own
validity unless their explicitly defined runtime dependencies require current
rule evidence; retained decisions, requests, grants, and historical pins remain
auditable. Neither supersession nor retirement silently migrates an artifact to
a successor version. ProfileAssignment, inheritance, conflict/exception
resolution, and EffectiveGovernanceContext remain FEATURE-0020-owned.
GovernanceProfile is an advanced administrator contract and is absent from
routine user journeys.

Any ApprovalPolicy or privileged-access, review, exception, or audit rule used
by FEATURE-0018 is supplied as exact system-selected evidence. FEATURE-0018
validates but never resolves profile applicability or precedence. Before
FEATURE-0020 exists, deterministic fixtures provide the preselected evidence
solely for FEATURE-0018-local conformance; no requester or routine user selects
it. Absence of an applicable `approvalPolicyRefs` entry is not by itself
authorization for `ApprovalRequirement.NotRequired`; that requirement must be
explicit trusted operation-local evidence under F18-RD-10.

### F18-RD-19 — FEATURE-0013 adoption

Register exactly these FEATURE-0013 DecisionProfiles through the approved
extension process:

- `authorization-decision/v1`;
- `approval-decision/v1`; and
- `exception-decision/v1`.

Their adopting-domain semantics are closed:

| DecisionProfile | Sole authority and subject | Exact profile input/result/validity |
|---|---|---|
| `authorization-decision/v1` | Authorization controller; exact-resource target, create-parent scope, or scope-only target | Input is the exact AuthorizationInput, contributing resource/version refs, scope relationship, and guardrail evidence; result is `Allow | Deny`; validity is the evaluation instant only and is never bearer authority |
| `approval-decision/v1` | Approval controller; exact ApprovalRequest | Input is exact policy version, immutable subject/proposal digest, active stage, eligible PrincipalRefs, and counted decisions; aggregate result is `Approved | Denied`; validity ends no later than ApprovalRequest `expiresAt` |
| `exception-decision/v1` | Exception controller; exact ExceptionProposal | Input is exact control version, subject, scope, interval, typed override, compensating controls, and approval evidence; result is `Grant | Deny`; `Deny` is decision-instant evidence, while `Grant` validity is bounded by the approved exception interval; neither is bearer authority, grants a role, or resolves effective governance |

All three reuse FEATURE-0013's existing envelope, authority, reason,
obligation, projection, validity, and audit-linkage mechanics. Authorization has
a protected full-evidence projection and a redacted safe-denial projection.
Approval cancellation and expiry remain lifecycle AuditEvents rather than
fabricated approval decisions. No adopting profile changes FEATURE-0013 or
performs a downstream effect.

FEATURE-0018 creates no competing envelope. AccessReview remains its own
retained LRO and does not add a fourth profile. Every authorization evaluation
emits one protected outcome AuditEvent. Material, denied, and privileged
authorization may additionally use `authorization-decision/v1` when policy
requires. RoleDefinition suspension, restoration, and retirement remain the
explicit mandatory `authorization-decision/v1` cases under F18-RD-07.

Every terminal `Approved` or `Denied` ApprovalRequest emits exactly one
mandatory `approval-decision/v1` DecisionRecord. `Cancelled` and `Expired`
remain lifecycle AuditEvents without a fabricated approval decision unless an
applicable auditRequirement independently requires an additional record. Every
accepted ExceptionProposal reaches F18-RD-17's exact terminal `Grant | Deny`
outcome and emits exactly one mandatory `exception-decision/v1`
DecisionRecord; `Deny` creates no ExceptionGrant. Required DecisionRecord
publication follows F18-RD-20 and is effective-once under F18-RD-21. An
applicable auditRequirement may require additional records but `Never` cannot
suppress these mandatory cases. Status contains only current state; lifecycle
history remains in AuditEvent and applicable DecisionRecord evidence.

### F18-RD-20 — Audit before authorization-changing publication

Required protected AuditEvent-obligation acceptance follows FEATURE-0013's
local atomic boundary and precedes every authorization-changing publication:
AccessGroup create/update/suspend/resume/retire; Membership
create/synchronize/freshness-expiry/guest-expiry/suspend/reactivate/revoke;
RoleDefinition publish/suspend/restore/retire; RoleAssignment
create/revoke/expiry/review-overdue ineligibility/review certification/due-time
advancement;
ApprovalPolicy and GovernanceProfile publication/retirement/supersession;
ApprovalRequest creation, terminal decision, and Approved-evidence expiry;
privileged activation,
blocked-readiness publication, expiry, or revocation; AccessReview lifecycle,
snapshot-blocked readiness, item decision, and remediation; ExceptionProposal
terminal decision; ExceptionGrant issuance, revocation, or expiry;
and completed idempotency result. Failure returns a safe internal failure and
publishes none of the state change or completed replay result. Application logs
never substitute for the protected obligation or AuditEvent.

For an AccessGroup Membership assignment-effect expansion, the protected evidence
pins the actor or canonical System provisioner, beneficiary, exact Membership
intent or current version, AccessGroup version, every group-held RoleAssignment
version in the derived `MembershipEnabledGrantEnvelope`, the applicable
Membership-operation and `roleassignment.grant` authorities, every per-action
delegable witness, computed effect interval, decision reason, request, and
correlation. This evidence belongs in protected DecisionRecord/AuditEvent
provenance; it is not copied into user-authored Membership fields and never
becomes bearer authority.

For approver, privileged-requester, and reviewer eligibility, protected
evidence pins the exact Human actor, selected EligibilityRef, exact policy or
rule version, evaluation time, scope and target, the separately evaluated
canonical action authority, separation-of-duties result where applicable, and stable
reason. The value and its evidence are non-bearer, are re-evaluated at the
authoritative decision/effect boundary, and are never copied into user-authored
Membership fields.

Approval publication and downstream effect publication are separate atomic
operations. The approval boundary atomically publishes the terminal
`Approved | Denied` ApprovalRequest state, mandatory `approval-decision/v1`
DecisionRecord, and required audit evidence; it publishes no RoleAssignment,
PrivilegedAccessRequest activation, AccessReview, or ExceptionGrant.

An ordinary assignment publication is a later freshly authorized atomic
boundary. Privileged activation is a later freshly authorized atomic boundary
that includes the RoleAssignment publication, the corresponding
PrivilegedAccessRequest `Active` state, any mandatory linked break-glass
AccessReview acceptance, and their required audit evidence. Failure publishes
none of that activation boundary.

The exception controller's final `Grant | Deny` is one separate atomic boundary
containing its mandatory `exception-decision/v1` DecisionRecord, required audit
evidence, and the exact resulting ExceptionGrant when the result is `Grant`.
That boundary cannot partially publish. No downstream effect is folded into the
earlier approval boundary.

Every other accepted mutation represented by the registered taxonomy—including
Draft create/update, proposal submission, request submission, synchronization,
and campaign creation—also requires protected audit-obligation acceptance
before publication even when it does not yet change effective authorization.

Every authorization evaluation must also accept its protected audit obligation
before returning either Allow or Deny. Acceptance failure returns a safe
internal error, no AuthorizationResult, and authorizes no downstream effect.
Later materialization may be asynchronous only as FEATURE-0013 permits and
never changes the pre-publication obligation boundary.

The initial registered FEATURE-0018 AuditEvent type taxonomy is exactly:

```text
authorization.evaluated
accessgroup.created | accessgroup.updated | accessgroup.suspended |
  accessgroup.resumed | accessgroup.retired
membership.created | membership.synchronized | membership.suspended |
  membership.freshness-expired | membership.expired |
  membership.reactivated | membership.revoked
roledefinition.draft-created | roledefinition.draft-updated |
  roledefinition.published | roledefinition.suspended |
  roledefinition.restored | roledefinition.retired |
  roledefinition.superseded
roleassignment.proposed | roleassignment.created | roleassignment.revoked |
  roleassignment.expired | roleassignment.review-overdue |
  roleassignment.review-certified | roleassignment.review-due-advanced
approvalpolicy.draft-created | approvalpolicy.draft-updated |
  approvalpolicy.published | approvalpolicy.retired |
  approvalpolicy.superseded
approvalrequest.created | approvalrequest.approved | approvalrequest.denied |
  approvalrequest.cancelled | approvalrequest.expired |
  approvalrequest.evidence-expired
privilegedaccessrequest.submitted | privilegedaccessrequest.approval-mapped |
  privilegedaccessrequest.activation-blocked |
  privilegedaccessrequest.activated | privilegedaccessrequest.cancelled |
  privilegedaccessrequest.expired | privilegedaccessrequest.revoked
accessreview.created | accessreview.started | accessreview.snapshot-blocked |
  accessreview.item-decided |
  accessreview.completed | accessreview.cancelled |
  accessreview.expired-incomplete | accessreview.remediated
exceptiongrant.proposed | exceptiongrant.issued | exceptiongrant.revoked |
  exceptiongrant.expired
governanceprofile.draft-created | governanceprofile.draft-updated |
  governanceprofile.published | governanceprofile.retired |
  governanceprofile.superseded
```

Each type uses FEATURE-0013 actor, action, subject, scope, outcome, correlation,
decision/evidence linkage, redaction, and projection semantics. The taxonomy
excludes structural noise, assertions, confidential policy input, secrets, and
inaccessible target details. A new event type requires FEATURE-0013 extension
registration and FEATURE-0018 conformance evidence.

Authoritative-time effectiveness is independent of successful status/event
materialization: reaching Membership freshness/guest expiry, RoleAssignment
expiry, a pre-activation PrivilegedAccessRequest `activationDeadline`, the
linked RoleAssignment expiry for an Active PrivilegedAccessRequest,
ApprovalRequest `expiresAt`, AccessReview `dueAt`, or ExceptionGrant effective
expiry fails closed at the boundary without grace. Audit or projection failure
never extends access or exception effect;
the owning controller retries only the permitted retained status, event, or
non-authoritative projection materialization. ExceptionGrant itself remains
immutable under F18-RD-17.

### F18-RD-22 — Deterministic in-memory foundation and local conformance

FEATURE-0018 uses synthetic fixtures that are immutable to consumers, rejects
invalid fixture sets before evaluation, obtains time from a deterministic UTC
source, and performs zero network or external I/O. Design owns the concrete
copying, construction, validation, and time-supply mechanisms. Fakes cover
approval-to-request-state mapping, denial, cancellation, blocked JIT activation,
successful JIT activation, break-glass direct and normal paths, expiry,
revocation, review, and exception evidence without simulating an IdP, policy
engine, workflow engine, or CloudProvider.

Conformance includes ordinary `NotRequired` direct grant, mandatory `Required`
JIT/privileged-role/exception paths, `breakGlassAllowed=true` direct bypass,
`false` normal approval, missing or ambiguous rule/requirement/policy evidence,
every allowed and rejected GovernanceApplicability component registration,
responsible-party inactivity at create/replace/retain boundaries and non-per-
request runtime behavior, all three exact AccessReview due-time calculations,
reviewer-rule mismatch, mandatory approval and exception DecisionRecord
emission, RoleDefinition supersession/suspension/retirement effects, every exact
action/mode/kind target-binding tuple and mismatch, admission approval current
before effect publication and provenance-only afterward, separate approval,
activation, and exception atomic-failure boundaries, activation immediately
before/at/after `activationDeadline` and active access continuing past that
deadline only until RoleAssignment expiry, Human-only `EligibilityRef` OR,
missing, stale, and unsupported-kind cases for approver, requester, and reviewer lists, validation-
precedence and safe-resolution ordering,
inactive-owner privilege-increase denial and privilege-reduction recovery,
`Standard | Guest` Membership validation with replacement/renewal rejection,
exact-version StandingCertification Retain advancement with Manual and
BreakGlassRetrospective non-advancement, Manual campaign containment at every
ScopeKind, every ExceptionProposal terminal mapping, immutable ExceptionGrant
expiry, and atomic publication failure.

Conformance also covers MembershipRelationshipKey uniqueness, terminal Guest
expiry versus recoverable synchronized staleness, the three exact FEATURE-0016-
owned ExecutionTarget target bindings and intrinsic classes, each malformed
ActionTargetBinding union variant, grantor-ceiling evaluation at publication
without later cascading revocation, per-action resource-A-to-B,
resource-to-scope, scope-to-resource and cross-assignment-synthesis cases,
deterministic Standing holder/scope/target registrations, indirect-beneficiary
separation of duties for both approval and AccessReview Retain/Replace/due
advancement at atomic decision/effect publication, including AccessGroup-held
authority whose current direct Workload/System member resolves to the acting
reviewer as its responsible Human; unresolved or stale group ownership,
Membership, or responsibility evidence; current and proposed replacement
beneficiary unions; privilege-reducing Revoke behavior when expansion is
unresolved; AccessGroup Membership addition to an empty group; matching bounded
membership-enabled grant envelopes; missing action, scope, resource, or temporal
witnesses; TimeBound administration attempting to enable Standing access;
self- or accomplice-add without a sufficient envelope; trusted-provisioner
overreach; stale-to-current and `freshUntil` extension; privilege-reducing
Membership transitions with unresolved dependencies; and concurrent Membership
and group-held RoleAssignment publication where the second publisher must
observe and authorize the effective access or conflict/retry; exact Human
PrincipalRef eligibility; structural rejection of RoleDefinition,
RoleAssignment, AccessGroupRef, Membership, external/nested/dynamic group, and
other relationship kinds; separate required-action authorization; inactive,
stale, scope-incompatible, target-incompatible, and decision-to-publication eligibility-change cases for
approver, requester, and reviewer lists,
missing/stale/non-phishing-resistant/below-AAL2 activation evidence for normal
JIT and break-glass, mode-specific access-review rule fields, exception
terminal/retryable reasons and narrowing algebra without silent rewrite,
GovernanceProfile supersession/retirement effects, and the Sovrunn/native-IAM
intersection boundary without provider-native objects, credentials, or external
effects.

Every material positive has applicable malformed, missing, unauthorized,
expired, revoked, stale, cross-scope, replay, audit-failure, and redaction
negatives. Race proof covers terminal approval and request-state mapping, JIT
and break-glass activation, revoke/use, review/remediation, and exception
revocation. Only FEATURE-0018-owned conformance identifiers and local proof
count as feature acceptance. No integration proof substitutes for local proof.

### F18-RD-23 — Standards-validation gate

Before final architecture approval, every applicable invariant maps to
FEATURE-0018 proof, a reused dependency, or a named exclusion under NIST SP
800-207, NIST SP 800-53 Rev. 5, NIST SP 800-63-4, NIST SP 800-162, ISO/IEC
27001/27002, CSA CCM, CIS Controls 5/6/8, and OWASP ASVS.

OIDC, OAuth, SCIM, WebAuthn, and SPIFFE are paper-only future-adapter
compatibility checks. They do not authorize adapter design or implementation.
The applicable ASVS profile and CIS/CSA shared-responsibility mappings must be
closed before final architecture approval. Formal ISO conformity claims remain
outside this handoff without licensed/current clause-level assessment.
