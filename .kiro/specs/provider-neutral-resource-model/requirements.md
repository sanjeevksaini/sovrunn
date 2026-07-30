# FEATURE-0014 Provider-Neutral Resource Model — Requirements

- Feature: FEATURE-0014 — Provider-Neutral Resource Model
- Stage: Requirements
- Phase: Phase 2 — Reuse-First PaaS Fabric Foundation
- Kiro slug: provider-neutral-resource-model
- Controlling architecture: `docs/architecture/provider-neutral-resource-model.md`
- Controlling handoff: `ADH-2026-018` (Approved)
- Canonical reuse standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
- FEATURE-0013 downstream adoption: `NOT_APPLICABLE` (architecture section 9)

## 1. Introduction

FEATURE-0014 defines the minimum provider-neutral substrate topology for
Sovrunn using exactly five resource kinds: `Provider`, `ProviderLocation`,
`ProviderDatacenter`, `DatacenterFailureDomain`, and `InfrastructureStack`.
These resources describe where provider supply exists and how it is
structurally partitioned. They do not state available capacity or capability,
select a placement target, or communicate with provider systems.

This document translates the closed architecture decision register
(`F14-AD-001`–`F14-AD-021`) and the risk register (`F14-R01`–`F14-R30`) into
testable requirements. It does not reopen, weaken, extend, reorder, or replace
any closed decision. Every normative requirement cites at least one
`F14-AD-*` decision or an explicitly inherited FEATURE-0012 contract. Any
unresolved semantic choice outside the section 13 register and section 14–15
delegation matrix stops this stage with `ARCHITECTURE_DECISION_REQUIRED`; no
default is selected here.

Inherited contracts remain owned by their canonical feature. FEATURE-0014
reuses the existing Phase 1 `Organization` resource and the FEATURE-0012
resource grammar unchanged, and preserves the closed FEATURE-0013 boundary
without semantic adoption. Requirements that restate inherited behavior do so
by reference, not by redefinition.

Normative keywords `SHALL` and `MUST` are used only where traced to a closed
decision or an inherited contract. Requirements addressing several resources
uniformly are stated once and parameterized with per-resource acceptance
examples rather than duplicated per resource.

## 2. Glossary

New or feature-specific concepts introduced by FEATURE-0014. Terms already
defined by FEATURE-0012, FEATURE-0013, or the Phase 1 `Organization` model are
referenced, not redefined.

- **Provider**: A Sovrunn resource identifying an infrastructure supplier or
  provider-supply boundary. It is not the owner organization, adapter
  configuration, credential holder, capability aggregate, customer
  organization, or provider plugin (architecture section 6.1).
- **ProviderLocation**: The only permitted kind and term for a named
  geographic or operational provider grouping under a `Provider`
  (architecture section 6.2). No alias kind, alias field, or alternate
  location-level vocabulary exists.
- **ProviderDatacenter**: A provider-operated physical datacenter within
  exactly one `ProviderLocation`. It is not a cluster, resource pool, or
  placement target (architecture section 6.3).
- **DatacenterFailureDomain**: A provider-declared isolation boundary inside
  exactly one `ProviderDatacenter` whose members may share a common failure
  mode. FEATURE-0014 records identity and containment only (architecture
  section 6.4).
- **InfrastructureStack**: One uniquely identified deployed infrastructure
  stack contained by exactly one `DatacenterFailureDomain` (architecture
  section 6.5). It is the deliberate terminology replacement for the
  superseded stack-kind name and is not a sixth kind or compatibility alias.
- **Supply owner**: The existing `Organization` that governs which provider
  supply is registered and visible within its boundary, expressed only
  through a Provider's primary `scopeRef` (architecture section 5.1).
- **Provider operator**: The party under which provider topology is declared,
  represented by the `Provider` resource (architecture section 5.1).
- **Registered shell**: A parent resource that exists but currently has zero
  children; valid only as topology-incomplete, never as a license to skip a
  required level (architecture section 4, `F14-AD-007`).
- **Topology-complete path**: An unbroken `Provider` → `ProviderLocation` →
  `ProviderDatacenter` → `DatacenterFailureDomain` → `InfrastructureStack`
  path. Topology completeness is an administrative fact and is not capability,
  capacity, placement eligibility, or service readiness (`F14-AD-007`).
- **Connectivity Unknown**: The absence of an explicit, separately owned
  connectivity assertion; it means neither connected nor disconnected
  (`F14-AD-013`). Connectivity ownership is deferred to FEATURE-0053.

## 3. User stories

Personas: **Platform operator** (Sovrunn platform authority), **Provider
operator** (party declaring provider topology), **Governance owner**
(organization administrative authority), and **Downstream feature**
(FEATURE-0015 and later consumers of topology references).

Each user story is `REFERENCE_ONLY`: it names a need and points to the
canonical `F14-REQ-*` requirement(s) that own the behavior. A story introduces
no normative obligation and never restates requirement text.

- **US-01** — As a platform operator, I want to register a `Provider` under
  Platform or an existing Organization so that provider supply has a single
  governance scope. (F14-REQ-04, F14-REQ-05)
- **US-02** — As a governance owner, I want an existing `Organization` to own
  provider supply without a new owner resource so that governance stays in one
  model. (F14-REQ-04, F14-REQ-11)
- **US-03** — As a provider operator, I want to declare the location →
  datacenter → failure-domain → stack hierarchy so that substrate topology is
  described in a provider-neutral way. (F14-REQ-08, F14-REQ-09)
- **US-04** — As a provider operator, I want to register a parent before its
  children so that onboarding, maintenance, and decommissioning are ordered
  and safe. (F14-REQ-12, F14-REQ-13)
- **US-05** — As a downstream feature, I want to know when a topology path is
  complete versus merely registered so that I never treat an incomplete shell
  as usable. (F14-REQ-14)
- **US-06** — As a provider operator, I want each deployed stack to have a
  distinct immutable identity even when it uses the same technology as another
  stack so that stacks are never conflated. (F14-REQ-15, F14-REQ-16)
- **US-07** — As a security owner, I want detailed topology confined to
  operator-facing audiences with no-existence disclosure across scopes so that
  provider topology is not leaked. (F14-REQ-07, F14-REQ-22, F14-REQ-30,
  F14-REQ-31)
- **US-08** — As a platform architect, I want physical containment to never
  imply network reachability or isolation so that later placement never infers
  connectivity from ancestry. (F14-REQ-19, F14-REQ-20)
- **US-09** — As a lifecycle operator, I want parent deletion rejected while
  children exist so that topology is decommissioned leaf-first without orphans
  or silent reparenting. (F14-REQ-26)
- **US-10** — As an architecture owner, I want FEATURE-0015, FEATURE-0016,
  FEATURE-0053, and FEATURE-0013 semantics excluded so that FEATURE-0014 owns
  only topology. (F14-REQ-23, F14-REQ-24, F14-REQ-25)
- **US-11** — As an API architect, I want exactly one location-level term and
  one atomic terminology migration so that no alias or dual-name contract
  exists. (F14-REQ-02, F14-REQ-03)

## 4. Acceptance criteria (normative requirements)

Each requirement has a canonical semantic key `(owning feature, resource or
contract, actor, behavior, observable outcome)` recorded in the normalization
ledger (section 10.4). Acceptance criteria include positive, negative,
boundary, isolation, concurrency, deletion, terminology, and
architecture-drift cases as applicable. Inherited FEATURE-0012 behavior is
referenced through the inherited-contract family (F14-REQ-21) rather than
restated.

### 4.1 Resource kinds and terminology

**F14-REQ-01 — Exactly five owned resource kinds.**
Sovrunn SHALL define exactly `Provider`, `ProviderLocation`,
`ProviderDatacenter`, `DatacenterFailureDomain`, and `InfrastructureStack`,
and MUST NOT define any additional kind or compatibility alias for this
feature (`F14-AD-001`).
- Positive: the FEATURE-0014 schema/kind inventory resolves to exactly these
  five kinds.
- Negative: a sixth kind or an alias kind is rejected.
- Architecture-drift: a kind inventory containing fewer or more than five
  FEATURE-0014 kinds fails the boundary check.

**F14-REQ-02 — Single location-level term.**
`ProviderLocation` MUST be the only kind and term for the location level; no
alias kind, alias field, or alternate location-level vocabulary is permitted
(`F14-AD-006`).
- Positive: the location level is expressed only as `ProviderLocation`.
- Negative: any alternate location-level term or alias field is rejected.
- Terminology: a repository terminology scan finds no alternate location-level
  vocabulary in active FEATURE-0014 contracts.

**F14-REQ-03 — Atomic terminology migration for InfrastructureStack.**
`InfrastructureStack` MUST atomically replace the superseded stack-kind name in
active contracts, with no dual-name compatibility period or active alias
(`F14-AD-021`).
- Positive: active FEATURE-0014 contracts use only `InfrastructureStack`.
- Negative: the superseded stack-kind name appearing in an active contract
  fails the check.
- Terminology: a repository-wide scan finds the superseded name only in
  historical/migration/non-goal text, never as an active contract.

### 4.2 Owner, scope, and containment

**F14-REQ-04 — Reuse existing Organization as optional supply owner.**
The supply owner MUST be the existing Phase 1 `Organization` resource;
FEATURE-0014 MUST NOT define a new `Owner`, `ProviderOwner`, `SupplyOwner`, or
other owner resource or scope kind (`F14-AD-002`).
- Positive: a `Provider` references an existing Organization as supply owner.
- Negative: any new owner resource kind or new scope kind is rejected.

**F14-REQ-05 — Provider primary scope is Platform or Organization.**
A `Provider` primary `scopeRef` MUST be exactly `Platform` or one existing
`Organization`; `ownerRef` and unconstrained owner identifiers MUST NOT
express supply governance (`F14-AD-003`).
- Positive: a Platform-scoped and an Organization-scoped Provider are both
  accepted.
- Negative: expressing the supply owner through `ownerRef` or an unconstrained
  owner ID is rejected.
- Isolation: `metadata.scopeRef` remains the sole governance-scope authority.

**F14-REQ-06 — Descendants are immutably Provider-scoped.**
`ProviderLocation`, `ProviderDatacenter`, `DatacenterFailureDomain`, and
`InfrastructureStack` MUST each have the containing `Provider` as an immutable
primary `scopeRef` (`F14-AD-004`).
- Positive: each descendant resolves to exactly one Provider scope.
- Negative: attempting to mutate a descendant's Provider scope is rejected.

**F14-REQ-07 — Cross-provider references rejected without disclosure.**
A topology reference crossing Provider scope MUST be rejected using
FEATURE-0012 no-existence-disclosure behavior (`F14-AD-004`).
- Positive: a same-Provider immediate-parent reference resolves.
- Negative: a cross-provider parent reference fails without revealing target
  existence.
- Isolation: the failure response and error text disclose no cross-scope
  resource existence.

**F14-REQ-08 — Exact hierarchy with no skipped or reordered levels.**
The only permitted hierarchy MUST be `Provider` → `ProviderLocation` →
`ProviderDatacenter` → `DatacenterFailureDomain` → `InfrastructureStack`, with
no skipped, optional, reordered, or alternate parent level (`F14-AD-005`).
- Positive: a child attaches only to its required immediate parent kind.
- Negative: every skipped-level attachment (for instance a stack attaching
  directly to a datacenter) is rejected.

**F14-REQ-09 — Exactly one immutable immediate parent per child.**
Each child MUST have exactly one immutable immediate parent of the required
kind; reparenting MUST use recreate-and-migrate; ambiguous or multiple parents
MUST fail validation (`F14-AD-008`).
- Positive: a child has one typed immediate-parent reference.
- Negative: zero, multiple, or wrong-kind immediate parents are rejected.
- Negative: mutating an existing immediate-parent reference is rejected.

**F14-REQ-10 — Parent and child resolve to the same Provider scope UID.**
Every immediate parent and child MUST resolve to the same Provider scope UID
(`F14-AD-004`, `F14-AD-008`).
- Positive: parent and child share one Provider scope UID.
- Negative: a parent-child pair resolving to different Provider scopes is
  rejected.

**F14-REQ-11 — Supply-owner and provider-operator roles stay distinct.**
An `Organization` and a `Provider` MUST remain role-specific resources even
when they represent the same real-world party; equal names or legal identity
MUST NOT collapse their domain roles (`F14-AD-002`, `F14-AD-003`).
- Positive: distinct-party case (`Organization: OwnerOrganization-A` with
  `Provider: ProviderOperator-A` and `Provider: ProviderOperator-B`) is
  represented with separate scoped resources.
- Positive: same-party case (`Organization: ProviderOperator-A` with
  `Provider: ProviderOperator-A`) keeps two role-specific resources.
- Negative: a workflow assuming Organization UID equals Provider UID is
  rejected.

### 4.3 Registration and topology completeness

**F14-REQ-12 — Registered existence is distinct from completeness.**
Requirements MUST distinguish a registered resource from a topology-complete
path; a registered parent with zero children is valid only as a
topology-incomplete shell (`F14-AD-007`).
- Positive: an empty registered parent at any level is accepted as
  incomplete.
- Negative: an empty registered parent is not reported as topology-complete.

**F14-REQ-13 — Zero-child registered parents allowed, skips still rejected.**
A zero-child registered parent MUST be permitted to support ordered
onboarding, maintenance, decommissioning, and safe deletion, without
permitting any child to skip its required immediate parent (`F14-AD-007`).
- Positive: registering a parent before its children is accepted.
- Negative: a descendant attaching to a non-immediate ancestor is rejected.

**F14-REQ-14 — Topology completeness requires all five levels.**
Only an unbroken `Provider` → `ProviderLocation` → `ProviderDatacenter` →
`DatacenterFailureDomain` → `InfrastructureStack` path MAY be reported as
topology-complete; completeness MUST NOT be interpreted as capability,
capacity, placement eligibility, or service readiness (`F14-AD-007`).
- Positive: a full five-level path is reported topology-complete.
- Negative: a path missing any level is not reported topology-complete.
- Boundary: completeness carries no capability/capacity/placement meaning.

### 4.4 InfrastructureStack identity and technology

**F14-REQ-15 — Distinct immutable stack identity.**
Each `InfrastructureStack` MUST have a distinct immutable FEATURE-0012
`metadata.uid` and scoped identity; shared technology MUST NOT imply shared
stack identity (`F14-AD-009`).
- Positive: two stacks using the same technology retain distinct UIDs and
  parents.
- Negative: technology used as a lookup or reference key is rejected.

**F14-REQ-16 — One failure-domain parent per stack; no spanning.**
Each `InfrastructureStack` MUST belong to exactly one
`DatacenterFailureDomain` and MUST NOT span, belong to, or reference multiple
failure domains; cross-domain deployments MUST be represented as distinct
stacks (`F14-AD-010`).
- Positive: a stack has exactly one immutable failure-domain parent.
- Negative: zero, ambiguous, multiple, or cross-scope failure-domain parents
  are rejected.
- Negative: moving a stack across failure domains without recreate-and-migrate
  is rejected.

**F14-REQ-17 — Stack technology is descriptive only.**
Infrastructure-stack technology MUST be a bounded descriptive attribute and
MUST NOT become identity, a closed compatibility enum, or capability or
placement logic in FEATURE-0014 (`F14-AD-011`).
- Positive: technology is recorded as bounded descriptive metadata.
- Negative: matching or filtering that uses technology as capability evidence
  is rejected.
- Boundary: the architecture's illustrative technologies (Apache CloudStack,
  Cloud Foundry, Red Hat OpenShift, AWS IaaS, Azure IaaS) are explicitly
  non-normative and are not a closed core enum.

**F14-REQ-18 — No provider-native identity or references in core.**
Provider-native identifiers, SDK objects, endpoints, and native
availability-zone identifiers MUST NOT become Sovrunn core identity or
reference fields (`F14-AD-019`).
- Positive: core identity uses only FEATURE-0012 identity and scoped
  references.
- Negative: a native identifier used as core identity or a core reference is
  rejected.

### 4.5 Connectivity non-inference

**F14-REQ-19 — Containment never implies connectivity or isolation.**
No topology relationship, shared ancestor, sibling order, physical proximity,
or shared owner `Organization` MAY be interpreted as evidence of network
reachability, isolation, latency, bandwidth, or failure correlation; this
non-inference holds within a single `Provider` and across different `Provider`
scopes, including two Providers governed by the same owner `Organization`
(`F14-AD-012`).
- Positive: same-datacenter, cross-datacenter, cross-location, and
  cross-provider relationships carry no implied connectivity semantics,
  including two Providers governed by the same owner `Organization`.
- Negative: any inference of reachability from common ancestry — including a
  shared owner `Organization` across different Providers — is rejected.
- Negative: any inference of isolation from different ancestry — including
  different Providers under the same or a different owner `Organization` — is
  rejected.

**F14-REQ-20 — Missing connectivity means Unknown, with no connectivity fields.**
Absence of an explicit, separately owned connectivity assertion MUST mean
`Unknown`; FEATURE-0014 MUST NOT add a `connected` boolean, adjacency list,
route, endpoint, peer reference, reachability, quality, trust, or health field
to topology resources (`F14-AD-013`).
- Positive: connectivity between any two resources is `Unknown` by absence.
- Negative: introducing any connectivity field, graph, or status is rejected.
- Architecture-drift: a schema deny-list scan finds no connectivity field.

### 4.6 Inherited grammar and state

**F14-REQ-21 — Inherited FEATURE-0012 grammar family (reference only).**
All five kinds MUST reuse FEATURE-0012 grammar unchanged: canonical type
identity, ObjectMeta, `ManagedResource` profile, operator-facing boundary,
typed constrained references, ownership, status/conditions, ordered
validation, RFC 9457 Problem Details with stable codes and RFC 6901 paths,
concurrency (resourceVersion/ETag/If-Match), unknown/duplicate-field
rejection, and registered extensions; FEATURE-0014 MUST NOT add or reinterpret
this common grammar (`F14-AD-014`). This is an inherited-contract family;
domain requirements below specialize it only where FEATURE-0014 adds a
constraint.
- Positive: FEATURE-0012 conformance and regression suites remain applicable
  without modification.
- Negative: any new or modified profile, scope kind, reference, status,
  condition, error, concurrency, or extension grammar is rejected.

**F14-REQ-22 — Minimal system-owned state posture.**
Operator-declared topology MUST use the `ManagedResource` profile and
operator-facing boundary; `status`/conditions MUST represent current
system-owned facts only and MUST NOT become history, capability, capacity, or
audit (`F14-AD-015`).
- Positive: `status` reflects current validation and administrative facts.
- Negative: storing event history, capability, or capacity in `status` is
  rejected.
- Isolation: one registered producer owns each condition type; no external
  provider system is an authoritative writer in FEATURE-0014.

**F14-REQ-23 — FEATURE-0013 adoption is NOT_APPLICABLE.**
FEATURE-0014 MUST declare FEATURE-0013 adoption `NOT_APPLICABLE` and MUST NOT
produce, consume, specialize, project, audit, or redefine any governed
FEATURE-0013 concept (`DecisionRecord`, `DecisionProfile`, `EvaluationResult`,
`AuditEvent`, rationale, obligations, actor, subject, linkage, projection, or
composition) (`F14-AD-016`).
- Positive: the requirements declare `NOT_APPLICABLE` with the section 9
  rationale.
- Negative: introducing any governed FEATURE-0013 concept is rejected.
- Architecture-drift: a generic Phase 2 gate MUST NOT be satisfied by adding
  decision, audit, or operation semantics; such a demand is a tooling defect
  reported as `REPOSITORY_CONTEXT_NOT_READY`.

### 4.7 Adjacent-feature exclusions

**F14-REQ-24 — No ResourcePool or ProviderCapability semantics.**
FEATURE-0014 MUST NOT define or infer `ResourcePool`, `ProviderCapability`,
capacity, compatibility, placement candidacy, eligibility, or matching; these
are owned by FEATURE-0015 (`F14-AD-017`).
- Positive: schemas contain topology semantics only.
- Negative: any capacity, eligibility, compatibility, or placement field is
  rejected.
- Architecture-drift: a field deny-list scan finds no FEATURE-0015 semantics
  outside explicit non-goal/ownership text.

**F14-REQ-25 — No adapter or provider-integration semantics.**
FEATURE-0014 MUST NOT define adapter interfaces, provider clients, discovery
protocols, credentials, endpoints, provider calls, repositories, plugins,
operations, provisioning, or runtime execution; these are owned by
FEATURE-0016 or later work (`F14-AD-018`).
- Positive: no external-system boundary is introduced.
- Negative: any interface, SDK dependency, credential, endpoint, or provider
  call is rejected.
- Architecture-drift: an interface/dependency/SDK scan finds no integration
  artifact outside explicit non-goal/ownership text.

### 4.8 Lifecycle, sovereignty, bounds, concurrency, isolation, redaction

**F14-REQ-26 — Parent deletion rejected while children exist.**
Parent deletion MUST be rejected while children exist; FEATURE-0014 MUST
perform no cascade deletion and no silent reparenting; decommissioning MUST
proceed leaf-first (`F14-AD-020`).
- Positive: leaf-first deletion of a complete path succeeds.
- Negative: deleting a parent with existing children returns a
  child-existence conflict.
- Concurrency: concurrent delete/create and stale-version attempts are
  resolved by FEATURE-0012 conflict controls (see F14-REQ-29).
- Referential: an immutable UID is never reused after deletion.

**F14-REQ-27 — Geographic descriptors are declared topology only.**
Where a normalized geographic/sovereignty descriptor is carried, it MUST use a
bounded normalized shape and MUST be treated as declared topology, not as an
authoritative assignment, residency proof, compliance evidence, or policy
outcome (`F14-AD-011`, `ADH-2026-019`). FEATURE-0014 MUST validate syntax and
country-prefix agreement only and MUST NOT own or consult a geographic dataset,
library, membership lookup, or refresh lifecycle.
- Positive: a syntactically valid descriptor with a matching country prefix is
  accepted as declared topology metadata, even when assignment is unknown.
- Negative: malformed syntax or a subdivision/country prefix mismatch is
  rejected.
- Boundary: geography carries no authoritative-assignment, residency,
  compliance, or placement inference.

**F14-REQ-28 — Finite bounds and pagination inherited.**
All collections, strings, and payloads MUST be finite and MUST inherit
FEATURE-0012 bounds and pagination; unbounded recursive traversal or payloads
MUST NOT be introduced (`F14-AD-014`, `F14-AD-015`).
- Positive: bounded collections and paginated reads behave per FEATURE-0012.
- Negative: an unbounded field or traversal is rejected.

**F14-REQ-29 — Concurrency-safe onboarding and decommissioning.**
Concurrent onboarding, parent-mutation attempts, and decommissioning MUST use
immutable parents and FEATURE-0012 resourceVersion/ETag/If-Match with
deterministic recomputation and conflict errors (`F14-AD-014`, `F14-AD-020`).
- Positive: concurrent independent writes succeed with correct versioning.
- Negative: a stale-version write is rejected with a conflict error.
- Concurrency: completeness facts recompute deterministically after
  concurrent changes.

**F14-REQ-30 — Multi-owner isolation.**
The same external operator represented under multiple owner Organizations MUST
be treated as independent scoped identities; authorization MUST NOT use an
external operator name or native ID (`F14-AD-003`, `F14-AD-004`, `F14-AD-009`).
- Positive: same-operator/different-owner resources are isolated.
- Negative: cross-owner access by shared external identity is denied without
  existence disclosure.

**F14-REQ-31 — Native values and secrets are redacted.**
Provider-native identifiers, credentials, endpoints, and secrets MUST NOT
appear in common metadata, errors, conditions, or any audit-adjacent text and
MUST be redacted (`F14-AD-019`).
- Positive: responses, errors, and conditions contain no native secret values.
- Negative: a secret-bearing field or leaked native endpoint is rejected.

## 5. Non-goals (enforceable exclusions)

The following are locked non-goals restated for visibility. Each is an
enforceable exclusion whose enforcement lives in the cited owning requirement's
negative and architecture-drift acceptance criteria. Each entry is
`REFERENCE_ONLY`: it points to the canonical `F14-REQ-*` home and cited
`F14-AD-*` decision and introduces no independent normative obligation. Any
appearance in an active FEATURE-0014 contract fails the boundary check
(architecture sections 3.2 and 11).

- **NG-01** — No `ResourcePool` or any placement-candidate model (owned by
  FEATURE-0015; enforced by F14-REQ-24, `F14-AD-017`).
- **NG-02** — No `ProviderCapability`, capability vocabulary, discovery,
  validation, certification, or matching (owned by FEATURE-0015; enforced by
  F14-REQ-24, `F14-AD-017`).
- **NG-03** — No capacity, eligibility, or compatibility semantics (owned by
  FEATURE-0015; enforced by F14-REQ-24, `F14-AD-017`).
- **NG-04** — No adapter interfaces, provider clients, discovery protocols,
  repositories, credentials, or provider-native translation (owned by
  FEATURE-0016 or later work; enforced by F14-REQ-25, `F14-AD-018`).
- **NG-05** — No new or altered `DecisionRecord`, `DecisionProfile`,
  `EvaluationResult`, `AuditEvent`, reason-code, linkage, actor, subject,
  projection, or scope contract (owned by FEATURE-0013; enforced by
  F14-REQ-23, `F14-AD-016`).
- **NG-06** — No placement requests/decisions, policy evaluation, entitlement,
  runtime profiles, plugins, operations, provisioning, reconciliation,
  failover, disaster recovery, capacity, cost, billing, or production
  persistence (owned by FEATURE-0015, FEATURE-0016, and later work; enforced
  by F14-REQ-24 and F14-REQ-25, `F14-AD-017`, `F14-AD-018`).
- **NG-07** — No Kubernetes, OpenStack, VMware, AWS, Azure, or other vendor
  object as a core domain type (enforced by F14-REQ-18, `F14-AD-019`).
- **NG-08** — No new `Owner`, `ProviderOwner`, or `SupplyOwner` resource kind;
  ownership reuses the existing `Organization` (enforced by F14-REQ-04,
  `F14-AD-002`).
- **NG-09** — No network links, routes, peering, reachability, latency,
  bandwidth, trust zones, connectivity health, or `NetworkConnectivityProfile`
  (owned by FEATURE-0053; enforced by F14-REQ-19 and F14-REQ-20, `F14-AD-012`,
  `F14-AD-013`).
- **NG-10** — No new or altered FEATURE-0012 profile, scope kind, reference,
  status, error, concurrency, or extension grammar (enforced by F14-REQ-21,
  `F14-AD-014`).

## 6. Edge cases

Each edge case is `REFERENCE_ONLY`: it illustrates an observable situation and
cites the canonical `F14-REQ-*` requirement(s) that own the behavior. An edge
case introduces no new normative obligation.

- **EC-01** — Empty registered parent at each of the four parent levels:
  accepted as incomplete, never reported complete (F14-REQ-12, F14-REQ-14).
- **EC-02** — Skipped-level attachment at every boundary (location→stack,
  provider→datacenter, etc.): rejected (F14-REQ-08).
- **EC-03** — Two stacks with identical technology in two failure domains:
  distinct UIDs and parents retained (F14-REQ-15).
- **EC-04** — Stack with zero, ambiguous, or multiple failure-domain parents:
  rejected (F14-REQ-16).
- **EC-05** — Cross-provider parent reference: rejected without existence
  disclosure (F14-REQ-07).
- **EC-06** — Parent-child pair resolving to different Provider scopes:
  rejected (F14-REQ-10).
- **EC-07** — Same real-world party as both Organization and Provider: two
  role-specific resources, distinct UIDs (F14-REQ-11).
- **EC-08** — Same external operator under two owner Organizations: isolated
  scoped identities (F14-REQ-30).
- **EC-09** — Deletion of a parent that still has children: child-existence
  conflict (F14-REQ-26).
- **EC-10** — Concurrent create/delete and stale-version writes: FEATURE-0012
  conflict handling; deterministic completeness recomputation (F14-REQ-29).
- **EC-11** — Same-datacenter, cross-datacenter, cross-location, and
  cross-provider relationships — including two Providers governed by the same
  owner `Organization` — keep connectivity `Unknown` with no inference of
  reachability or isolation (F14-REQ-19, F14-REQ-20).
- **EC-12** — Malformed or country-prefix-inconsistent geographic descriptor:
  rejected; a syntactically valid unassigned descriptor is accepted without
  authoritative, residency, or compliance inference (F14-REQ-27).
- **EC-13** — Attempt to mutate an immediate-parent reference or Provider
  scope: rejected; requires recreate-and-migrate (F14-REQ-06, F14-REQ-09).
- **EC-14** — Superseded stack-kind name in an active contract: fails
  terminology check (F14-REQ-03).

## 7. Security and privacy requirements

Security and privacy obligations attach to their owning functional
requirements rather than forming duplicate requirement families. Each entry
below is `REFERENCE_ONLY`: it names the security or privacy concern and points
to the canonical `F14-REQ-*` home and cited `F14-AD-*` decision(s). No entry
states an independent `MUST`; the normative obligation lives in the cited
requirement.

- **SEC-01** — Sole governance-scope authority (`metadata.scopeRef` is the only
  scope authority; `ownerRef` and containment grant no authorization): owned by
  F14-REQ-05 and F14-REQ-06 (`F14-AD-003`, `F14-AD-004`).
- **SEC-02** — Operator-facing boundary with least-privilege projections for
  detailed datacenter and failure-domain topology; customer-facing or AI
  projection is a locked non-goal requiring separate review: owned by
  F14-REQ-22 (`F14-AD-015`).
- **SEC-03** — No-existence disclosure on cross-scope list, get, and reference:
  owned by F14-REQ-07 and F14-REQ-30 (`F14-AD-004`).
- **SEC-04** — Redaction of provider-native identifiers, credentials,
  endpoints, and secrets from metadata, errors, conditions, and audit-adjacent
  text: owned by F14-REQ-31 (`F14-AD-019`).
- **SEC-05** — Geographic descriptors are neither authoritative assignments nor
  residency/compliance evidence: owned by F14-REQ-27 (`F14-AD-011`,
  `ADH-2026-019`).
- **SEC-06** — Multi-owner isolation of one external operator, with
  authorization never keyed on external operator name or native ID: owned by
  F14-REQ-30 (`F14-AD-003`, `F14-AD-004`, `F14-AD-009`).
- **SEC-07** — Observability review is limited to contract validation, stable
  errors, and redaction; no decision, audit, or operation observability is
  required because FEATURE-0013 adoption is `NOT_APPLICABLE`: owned by
  F14-REQ-23 (`F14-AD-016`, architecture section 16.11).

## 8. Compatibility with completed Phase 1 features

Each entry below is `REFERENCE_ONLY`: it records a Phase 1 compatibility
concern and points to the canonical `F14-REQ-*` home and cited `F14-AD-*`
decision(s). No entry states an independent `MUST`; the normative obligation
lives in the cited requirement.

- **COMPAT-01** — Reuse of the existing Phase 1 `Organization` as the optional
  supply owner, without redefining or copying its identity, scope,
  authorization boundary, or lifecycle: owned by F14-REQ-04 and F14-REQ-11
  (`F14-AD-002`).
- **COMPAT-02** — Non-modification of the Phase 1 resource kinds
  (`Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Operation`,
  `ServiceClass`, `ServicePlan`, `Plugin`, `Capability`, `ServiceInstance`,
  `ServiceBinding`) and their registries is a consequence of reusing
  `Organization` unchanged and reusing the FEATURE-0012 grammar unchanged:
  owned by F14-REQ-04 and F14-REQ-21 (`F14-AD-002`, `F14-AD-014`). It adds no
  obligation beyond those homes.
- **COMPAT-03** — FEATURE-0014 introducing only new provider-topology kinds,
  and not changing the Phase 1 resource model or API-contract semantics those
  features rely on, is a consequence of the exactly-five-owned-kinds boundary
  and the unchanged inherited grammar: owned by F14-REQ-01, F14-REQ-04, and
  F14-REQ-21 (`F14-AD-001`, `F14-AD-002`, `F14-AD-014`). It introduces no
  `MUST` of its own.
- **COMPAT-04** — Backward compatibility of the reused Phase 2 FEATURE-0012
  grammar (which itself builds on the Phase 1 metadata/spec/status shape) with
  Phase 1 conventions: owned by F14-REQ-21 (`F14-AD-014`).
- **COMPAT-05** — The `Provider` scope kind is the existing FEATURE-0012/0013
  six-value `ScopeKind` value already present (`Platform`, `Organization`,
  `OrganizationUnit`, `Tenant`, `Project`, `Provider`); no seventh value is
  added: owned by F14-REQ-05 (`F14-AD-003`, `F14-AD-016`).

## 9. Design questions to resolve later (non-normative)

The following are bounded representation choices explicitly delegated to
design by architecture sections 14 and 15. They are non-normative here and
MUST NOT be resolved in requirements. No delegated choice may change a
resource's meaning, ownership, cardinality, scope, hierarchy, identity,
connectivity posture, deletion semantics, or feature boundary.

- **DQ-01** — Exact API group and versioned routes, within FEATURE-0012
  domain-group and maturity rules (one domain-grouped alpha API; no
  unversioned endpoint).
- **DQ-02** — The minimal domain field inventory needed to express the closed
  decisions (only fields necessary; no convenience fields).
- **DQ-03** — Exact condition names representing registered validity and
  topology completeness, using FEATURE-0012 condition grammar.
- **DQ-04** — Normalized geographic code representation and the specific
  mature standard/field, where applicable.
- **DQ-05** — Bounded descriptive infrastructure-technology representation
  that cannot become identity, capability, or compatibility logic.
- **DQ-06** — Initial finite limit values and the configuration mechanism.
- **DQ-07** — Schema composition, indexes, package layout, and internal
  implementation structure.
- **DQ-08** — Whether immediate-parent containment is represented by a
  domain-specific parent reference alone or additionally mirrored by
  FEATURE-0012 `ownerRef`, provided no second topology authority is created.

If any DQ item is discovered to require a semantic (not representational)
choice, design MUST stop with `ARCHITECTURE_DECISION_REQUIRED` rather than
select a default.

## 10. Architecture traceability

### 10.1 Decision traceability (`F14-AD-001`–`F14-AD-021`)

Every closed decision is individually translated into one or more testable
requirements. Requirements do not reinterpret, weaken, extend, or replace any
decision.

| Decision | Requirement(s) | Treatment |
|---|---|---|
| `F14-AD-001` | F14-REQ-01 | Translated now — exactly five kinds, no alias. |
| `F14-AD-002` | F14-REQ-04, F14-REQ-11 | Translated now — reuse Organization; no new owner. |
| `F14-AD-003` | F14-REQ-05, F14-REQ-30 | Translated now — Platform/Organization scope only. |
| `F14-AD-004` | F14-REQ-06, F14-REQ-07, F14-REQ-10, F14-REQ-30 | Translated now — immutable Provider scope. |
| `F14-AD-005` | F14-REQ-08 | Translated now — exact hierarchy, no skips. |
| `F14-AD-006` | F14-REQ-02 | Translated now — single location term. |
| `F14-AD-007` | F14-REQ-12, F14-REQ-13, F14-REQ-14 | Translated now — registration vs completeness. |
| `F14-AD-008` | F14-REQ-09, F14-REQ-10 | Translated now — one immutable parent. |
| `F14-AD-009` | F14-REQ-15, F14-REQ-30 | Translated now — distinct immutable stack UID. |
| `F14-AD-010` | F14-REQ-16 | Translated now — one failure-domain parent. |
| `F14-AD-011` | F14-REQ-17, F14-REQ-27 | Translated now — technology/geo descriptive only. |
| `F14-AD-012` | F14-REQ-19 | Translated now — no connectivity inference. |
| `F14-AD-013` | F14-REQ-20 | Translated now — Unknown; no connectivity fields. |
| `F14-AD-014` | F14-REQ-21, F14-REQ-28, F14-REQ-29 | Translated now — FEATURE-0012 grammar reuse. |
| `F14-AD-015` | F14-REQ-22 | Translated now — ManagedResource; status current-only. |
| `F14-AD-016` | F14-REQ-23 | Translated now — FEATURE-0013 `NOT_APPLICABLE`. |
| `F14-AD-017` | F14-REQ-24 | Translated now — no FEATURE-0015 semantics. |
| `F14-AD-018` | F14-REQ-25 | Translated now — no adapter/integration. |
| `F14-AD-019` | F14-REQ-18, F14-REQ-31 | Translated now — no native identity; redaction. |
| `F14-AD-020` | F14-REQ-26, F14-REQ-29 | Translated now — deletion rejected with children. |
| `F14-AD-021` | F14-REQ-03 | Translated now — atomic terminology migration. |

### 10.2 Risk traceability (`F14-R01`–`F14-R30`)

Every risk is individually preserved with its mapped requirement(s),
mitigation outcome, evidence obligation, target residual, owner, and
reassessment trigger. Target residual ratings are copied unchanged from
architecture section 12 and are goals pending implementation evidence and
human semantic review; this stage accepts no residual risk and downgrades,
merges, renumbers, or closes no risk.

| Risk | Requirement(s) | Mitigation outcome preserved | Evidence obligation | Target residual | Owner | Reassessment trigger |
|---|---|---|---|---:|---|---|
| `F14-R01` | F14-REQ-24, F14-REQ-01 | Topology only; no placement/capacity/capability | Field deny-list; changed-file review; negative fixtures | 1×5 Medium | FEATURE-0014 owner | Capacity/eligibility/compatibility/placement field proposed |
| `F14-R02` | F14-REQ-02, F14-REQ-01 | Only `ProviderLocation`; no alias/alternate term | Terminology scan; schema/route inventory | 1×3 Low | API architecture owner | Alias or alternate location term requested |
| `F14-R03` | F14-REQ-05, F14-REQ-06, F14-REQ-10 | `scopeRef` sole scope authority | Pos/neg scope fixtures; cross-scope authz tests | 1×5 Medium | Security/API owner | Second scope source or owner-based authz path |
| `F14-R04` | F14-REQ-04 | Reuse Organization; no new owner kind | Kind inventory; owner-scenario fixtures | 1×3 Low | Architecture owner | Organization cannot express a validated ownership case |
| `F14-R05` | F14-REQ-11 | Role-specific resources not collapsed | Same-party/distinct-party fixtures | 1×3 Low | Domain owner | Workflow assumes Organization UID equals Provider UID |
| `F14-R06` | F14-REQ-19, F14-REQ-20 | Connectivity non-inferable; missing = Unknown | Same/cross-parent non-inference fixtures; absence checks | 1×5 Medium | Network/placement owner | Placement consumes ancestry as connectivity |
| `F14-R07` | F14-REQ-20 | No adjacency/boolean connectivity field | Schema deny-list; semantic review | 1×4 Low | Network architecture owner | Early connectivity becomes a validated dependency |
| `F14-R08` | F14-REQ-18, F14-REQ-31 | Bounded descriptive only; no native leak/secrets | Provider-neutral scan; secret/redaction tests | 1×5 Medium | Security/integration owner | Native value required for core identity/output |
| `F14-R09` | F14-REQ-12, F14-REQ-13, F14-REQ-14 | Registered ≠ complete; unbroken five-level path | Empty-parent fixtures; completeness transition tests | 1×4 Low | Domain owner | Incomplete topology becomes a placement input |
| `F14-R10` | F14-REQ-15, F14-REQ-17 | Immutable UID per stack; shared tech ≠ identity | Duplicate-tech/distinct-UID fixtures | 1×3 Low | API owner | Technology used as lookup/reference identity |
| `F14-R11` | F14-REQ-16 | Exactly one immutable failure-domain parent | Zero/multiple/wrong-kind/cross-scope fixtures | 1×5 Medium | Domain owner | A deployment claims it cannot be distinct stacks |
| `F14-R12` | F14-REQ-08, F14-REQ-09 | One required immediate parent; same Provider scope | Skipped-level and wrong-parent fixtures | 1×5 Medium | Domain/API owner | A new topology form requires alternate containment |
| `F14-R13` | F14-REQ-26 | Reject deletion with children; leaf-first | Child-existence conflict; concurrent/stale/retry tests | 2×4 Medium | Lifecycle/API owner | Cascade/bulk decommissioning becomes a requirement |
| `F14-R14` | F14-REQ-22 | Status = Sovrunn-known facts only | Staleness/unknown-state fixtures; writer-ownership review | 2×3 Medium | Operator/domain owner | Discovery or external observation is introduced |
| `F14-R15` | F14-REQ-07, F14-REQ-22, F14-REQ-30 | Operator-facing; no-existence disclosure | Cross-scope list/get/reference tests; redaction review | 1×5 Medium | Security owner | Customer-facing or AI projection proposed |

| Risk | Requirement(s) | Mitigation outcome preserved | Evidence obligation | Target residual | Owner | Reassessment trigger |
|---|---|---|---|---:|---|---|
| `F14-R16` | F14-REQ-27 | Geography is syntax-normalized declared topology, not authoritative assignment or compliance proof | Malformed/prefix-mismatch fixtures; syntactically valid unassigned-code fixture; no-authority/no-residency-inference review | 2×4 Medium | Sovereignty/domain owner | Authoritative geography becomes policy/placement/attestation input |
| `F14-R17` | F14-REQ-28 | Finite bounds/pagination; no unbounded traversal | Boundary/property tests; pagination/max-size benchmarks | 2×3 Medium | API/performance owner | Provider topology exceeds reviewed limits |
| `F14-R18` | F14-REQ-29 | Immutable parents; resourceVersion/ETag/If-Match | Concurrency/stale-write/condition-consistency tests | 2×3 Medium | API/lifecycle owner | Multi-writer or external synchronization introduced |
| `F14-R19` | F14-REQ-03 | One atomic migration; no alias/dual-name | Old-name scan; schema/API diff; docs consistency | 1×4 Low | Architecture/API owner | Deployed external consumer needs compat migration |
| `F14-R20` | F14-REQ-17 | Technology descriptive, not compatibility enum | Enum/decision-use scan; negative compatibility tests | 1×4 Low | Domain/placement owner | Matching/filtering uses technology as capability |
| `F14-R21` | F14-REQ-22, F14-REQ-23 | Conditions current-only; F13 `NOT_APPLICABLE` | Status-history scan; FEATURE-0013 adoption gate | 1×3 Low | Architecture/audit owner | FEATURE-0014 produces/consumes governed F13 concept |
| `F14-R22` | F14-REQ-25 | No integration/runtime artifacts | Changed-file/dependency/interface scan; drift review | 1×5 Medium | Architecture owner | Implementation task adds an external-system boundary |
| `F14-R23` | F14-REQ-26, F14-REQ-15 | Immutable/non-reused UIDs; reject parent delete | UID non-reuse/stale-reference fixtures | 1×4 Low | API/later-feature owner | FEATURE-0015/0053 introduces durable references |
| `F14-R24` | F14-REQ-30, F14-REQ-07 | Independent scoped identity; no name/native authz | Same-operator/different-owner isolation and denial tests | 1×5 Medium | Security/domain owner | Cross-owner federation/shared-provider admin proposed |
| `F14-R25` | F14-REQ-01, F14-REQ-23, F14-REQ-24, F14-REQ-25 (whole-stage stop rule owned by architecture §14.1, §16.1) | Whole-stage `ARCHITECTURE_DECISION_REQUIRED` stop | Ambiguity fixtures; stage review of no unresolved choice | 1×5 Medium | Architecture owner | Downstream selects a value not closed or delegated |
| `F14-R26` | F14-REQ-01, F14-REQ-02, F14-REQ-03, F14-REQ-23, F14-REQ-24, F14-REQ-25 (context controls owned by architecture §2.2, §2.3, §16.8) | Task-specific allowlist; content-bound package | Terminology/ownership scan; context-manifest inspection | 1×5 Medium | Architecture/tooling owner | Allowlisted input conflicts or disallowed input enters |
| `F14-R27` | F14-REQ-01, F14-REQ-08, F14-REQ-12, F14-REQ-13, F14-REQ-14, F14-REQ-21 (normalization guard owned by architecture §16.9) | One normative home per semantic key | Normalization ledger; duplicate-key/similarity report | 1×4 Low | Requirements owner | Two requirements share a semantic key or conflict |
| `F14-R28` | F14-REQ-21, F14-REQ-23 (bidirectional exact-ID traceability owned by architecture §14.2, §16.10) | Bidirectional exact-ID traceability; no range-only | Automated enumeration of decisions/risks; orphan report | 1×5 Medium | Traceability/tooling owner | Missing ID, uncited statement, or orphan appears |
| `F14-R29` | F14-REQ-09, F14-REQ-15, F14-REQ-17, F14-REQ-21, F14-REQ-22, F14-REQ-26 (stage separation owned by architecture §16.1, §16.3) | Hard stage state machine; per-stage allowlists | Stage diff; design-language scan in requirements | 1×5 Medium | Kiro workflow owner | A stage modifies an earlier artifact or leaks a choice |
| `F14-R30` | F14-REQ-23, F14-REQ-24, F14-REQ-25 | Applicability-aware boundary; adjacent deny-list | Prompt/reviewer tests; forbidden-concept scan | 1×5 Medium | Architecture/tooling owner | Adjacent-feature semantics enter active artifacts |

Risks `F14-R25`–`F14-R29` are cross-cutting governance and process controls.
They trace to the canonical `F14-REQ-*` requirements listed above, while the
controlling workflow rules remain owned by the cited architecture sections
(2.2, 2.3, 14.1, 14.2, 16.1, 16.3, 16.8, 16.9, 16.10). Requirements reference
those controls and do not redefine or duplicate them as a competing source.

Residual-risk acceptance is human-owned and occurs only after implementation
evidence and semantic review (architecture section 12.1).

### 10.3 Single-owner overlap ledger

Every normative requirement has exactly one owning feature and one canonical
cited decision or inherited contract. Cross-cutting obligations (`US-*`,
`SEC-*`, `COMPAT-*`, `NG-*`, `EC-*`, `DQ-*`) reference these normative homes
and are recorded as `REFERENCE_ONLY` in the normalization ledger (10.4); they
introduce no second owner. No requirement is owned by two features.

| Requirement | Owning feature | Cited decision / inherited contract | Affected resource / field / behavior | Second owner? |
|---|---|---|---|---|
| F14-REQ-01 | FEATURE-0014 | `F14-AD-001` | Owned kind set | None |
| F14-REQ-02 | FEATURE-0014 | `F14-AD-006` | ProviderLocation term | None |
| F14-REQ-03 | FEATURE-0014 | `F14-AD-021` | InfrastructureStack naming | None |
| F14-REQ-04 | FEATURE-0014 | `F14-AD-002` (Organization owned by Phase 1) | Provider supply owner reference | None |
| F14-REQ-05 | FEATURE-0014 | `F14-AD-003` (scope grammar FEATURE-0012) | Provider `scopeRef` | None |
| F14-REQ-06 | FEATURE-0014 | `F14-AD-004` | Descendant Provider scope | None |
| F14-REQ-07 | FEATURE-0014 | `F14-AD-004` (disclosure FEATURE-0012) | Cross-provider reference | None |
| F14-REQ-08 | FEATURE-0014 | `F14-AD-005` | Hierarchy levels | None |
| F14-REQ-09 | FEATURE-0014 | `F14-AD-008` | Immediate parent reference | None |
| F14-REQ-10 | FEATURE-0014 | `F14-AD-004`, `F14-AD-008` | Parent/child scope UID | None |
| F14-REQ-11 | FEATURE-0014 | `F14-AD-002`, `F14-AD-003` | Owner/operator role separation | None |
| F14-REQ-12 | FEATURE-0014 | `F14-AD-007` | Registered vs complete | None |
| F14-REQ-13 | FEATURE-0014 | `F14-AD-007` | Zero-child registration | None |
| F14-REQ-14 | FEATURE-0014 | `F14-AD-007` | Topology completeness | None |
| F14-REQ-15 | FEATURE-0014 | `F14-AD-009` | InfrastructureStack UID identity | None |
| F14-REQ-16 | FEATURE-0014 | `F14-AD-010` | Stack failure-domain parent | None |

| Requirement | Owning feature | Cited decision / inherited contract | Affected resource / field / behavior | Second owner? |
|---|---|---|---|---|
| F14-REQ-17 | FEATURE-0014 | `F14-AD-011` | Stack technology descriptor | None |
| F14-REQ-18 | FEATURE-0014 | `F14-AD-019` | Native identity exclusion | None |
| F14-REQ-19 | FEATURE-0014 | `F14-AD-012` | Connectivity non-inference | None |
| F14-REQ-20 | FEATURE-0014 | `F14-AD-013` | Connectivity Unknown; no field | None |
| F14-REQ-21 | FEATURE-0012 (inherited) | `F14-AD-014` | Common resource grammar | None (referenced) |
| F14-REQ-22 | FEATURE-0014 | `F14-AD-015` (grammar FEATURE-0012) | Profile/boundary/status posture | None |
| F14-REQ-23 | FEATURE-0014 | `F14-AD-016` (concepts FEATURE-0013) | FEATURE-0013 `NOT_APPLICABLE` | None |
| F14-REQ-24 | FEATURE-0014 | `F14-AD-017` (semantics FEATURE-0015) | ResourcePool/capability exclusion | None |
| F14-REQ-25 | FEATURE-0014 | `F14-AD-018` (semantics FEATURE-0016) | Adapter/integration exclusion | None |
| F14-REQ-26 | FEATURE-0014 | `F14-AD-020` | Parent deletion behavior | None |
| F14-REQ-27 | FEATURE-0014 | `F14-AD-011` | Geographic descriptor | None |
| F14-REQ-28 | FEATURE-0014 | `F14-AD-014`, `F14-AD-015` (bounds FEATURE-0012) | Finite bounds/pagination | None |
| F14-REQ-29 | FEATURE-0014 | `F14-AD-014`, `F14-AD-020` (concurrency FEATURE-0012) | Concurrency controls | None |
| F14-REQ-30 | FEATURE-0014 | `F14-AD-003`, `F14-AD-004`, `F14-AD-009` | Multi-owner isolation | None |
| F14-REQ-31 | FEATURE-0014 | `F14-AD-019` | Native/secret redaction | None |

The overlap check passes: every requirement has exactly one owner, cites at
least one closed decision or inherited contract, and introduces no
adjacent-feature semantics. No requirement cites only a lower-precedence
summary.

### 10.4 Requirement normalization ledger

Each requirement is the single normative home for its semantic key
`(owning feature, resource/contract, actor, behavior, observable outcome)`.
No two requirements share a semantic key. Cross-cutting entries (`SEC-*`,
`COMPAT-*`, `NG-*`, `US-*`, `EC-*`, `DQ-*`) are `REFERENCE_ONLY` to these
homes.

| Requirement | Semantic key (abbrev.) | Upstream owner/decision | Referenced risks | Disposition |
|---|---|---|---|---|
| F14-REQ-01 | F14 · kind set · operator · define · exactly five | `F14-AD-001` | F14-R01, F14-R02, F14-R19 | NORMATIVE_HOME |
| F14-REQ-02 | F14 · location term · operator · name · single term | `F14-AD-006` | F14-R02 | NORMATIVE_HOME |
| F14-REQ-03 | F14 · stack naming · architect · migrate · no alias | `F14-AD-021` | F14-R19 | NORMATIVE_HOME |
| F14-REQ-04 | F14 · Provider owner · gov owner · reference · reuse Org | `F14-AD-002` | F14-R04 | NORMATIVE_HOME |
| F14-REQ-05 | F14 · Provider scopeRef · operator · scope · Platform/Org | `F14-AD-003` | F14-R03 | NORMATIVE_HOME |
| F14-REQ-06 | F14 · descendant scope · operator · scope · immutable | `F14-AD-004` | F14-R03, F14-R12 | NORMATIVE_HOME |
| F14-REQ-07 | F14 · cross-provider ref · operator · reject · no disclosure | `F14-AD-004` | F14-R15, F14-R24 | NORMATIVE_HOME |
| F14-REQ-08 | F14 · hierarchy · operator · attach · no skip/reorder | `F14-AD-005` | F14-R12 | NORMATIVE_HOME |
| F14-REQ-09 | F14 · parent ref · operator · attach · one immutable | `F14-AD-008` | F14-R11, F14-R12 | NORMATIVE_HOME |
| F14-REQ-10 | F14 · scope UID · system · resolve · same Provider | `F14-AD-004`, `F14-AD-008` | F14-R03, F14-R12 | NORMATIVE_HOME |
| F14-REQ-11 | F14 · roles · domain · separate · not collapsed | `F14-AD-002`, `F14-AD-003` | F14-R05, F14-R24 | NORMATIVE_HOME |
| F14-REQ-12 | F14 · registration · system · classify · registered≠complete | `F14-AD-007` | F14-R09 | NORMATIVE_HOME |
| F14-REQ-13 | F14 · empty parent · operator · register · no skip | `F14-AD-007` | F14-R09 | NORMATIVE_HOME |
| F14-REQ-14 | F14 · completeness · system · report · five-level path | `F14-AD-007` | F14-R09 | NORMATIVE_HOME |
| F14-REQ-15 | F14 · stack identity · system · identify · distinct UID | `F14-AD-009` | F14-R10, F14-R23 | NORMATIVE_HOME |
| F14-REQ-16 | F14 · stack parent · operator · attach · one domain | `F14-AD-010` | F14-R11 | NORMATIVE_HOME |

| Requirement | Semantic key (abbrev.) | Upstream owner/decision | Referenced risks | Disposition |
|---|---|---|---|---|
| F14-REQ-17 | F14 · stack technology · operator · describe · non-capability | `F14-AD-011` | F14-R20, F14-R16 | NORMATIVE_HOME |
| F14-REQ-18 | F14 · core identity · system · exclude · no native id | `F14-AD-019` | F14-R08 | NORMATIVE_HOME |
| F14-REQ-19 | F14 · connectivity · consumer · infer · none from ancestry | `F14-AD-012` | F14-R06 | NORMATIVE_HOME |
| F14-REQ-20 | F14 · connectivity · system · default · Unknown, no field | `F14-AD-013` | F14-R06, F14-R07 | NORMATIVE_HOME |
| F14-REQ-21 | FEATURE-0012 · grammar · system · reuse · unchanged | `F14-AD-014` | F14-R17, F14-R18, F14-R21 | INHERITED_REFERENCE |
| F14-REQ-22 | F14 · status · system · represent · current facts only | `F14-AD-015` | F14-R14, F14-R21 | NORMATIVE_HOME |
| F14-REQ-23 | F14 · F13 adoption · architect · declare · NOT_APPLICABLE | `F14-AD-016` | F14-R21, F14-R30 | NORMATIVE_HOME |
| F14-REQ-24 | F14 · FEATURE-0015 semantics · architect · exclude · none | `F14-AD-017` | F14-R01, F14-R30 | NORMATIVE_HOME |
| F14-REQ-25 | F14 · FEATURE-0016 semantics · architect · exclude · none | `F14-AD-018` | F14-R22, F14-R30 | NORMATIVE_HOME |
| F14-REQ-26 | F14 · deletion · operator · delete · reject with children | `F14-AD-020` | F14-R13, F14-R23 | NORMATIVE_HOME |
| F14-REQ-27 | F14 · geography · operator · declare · not compliance | `F14-AD-011` | F14-R16 | NORMATIVE_HOME |
| F14-REQ-28 | F14 · bounds · system · limit · finite/paginated | `F14-AD-014`, `F14-AD-015` | F14-R17 | NORMATIVE_HOME |
| F14-REQ-29 | F14 · concurrency · system · version · conflict-safe | `F14-AD-014`, `F14-AD-020` | F14-R18 | NORMATIVE_HOME |
| F14-REQ-30 | F14 · multi-owner · security · isolate · independent identity | `F14-AD-003`, `F14-AD-004`, `F14-AD-009` | F14-R24 | NORMATIVE_HOME |
| F14-REQ-31 | F14 · redaction · security · redact · no native/secret | `F14-AD-019` | F14-R08 | NORMATIVE_HOME |

No duplicate semantic key exists; no conflicting requirements share a key.
Cross-cutting `SEC-*`, `COMPAT-*`, `NG-*`, `US-*`, `EC-*`, and `DQ-*` entries
are `REFERENCE_ONLY` and add no normative home.

### 10.5 Feature-level reuse summary

Per the canonical `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` (format
`1.0.0`); the full capability-level assessments and Approved human-approval
evidence live in `docs/features/FEATURE-0014-provider-neutral-resource-model.md`
and `docs/reviews/reuse-assessments/FEATURE-0014-approval-evidence.md`. This
summary references those records and does not redefine the schema.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| Organization owner identity and FEATURE-0012 resource grammar | Extend | Reuse existing Organization and common grammar; add only constrained topology semantics. | Approved | ADH-2026-018; ADH-2026-012 |
| Five-resource provider-neutral topology semantics | Build | No external model supplies Sovrunn's exact owner/provider separation, immutable hierarchy, completeness, one-domain stack identity, and non-connectivity boundary. | Approved | ADH-2026-018; RFC-0024 |
| Provider integration and discovery | Wrap | Future provider-native translation belongs behind FEATURE-0016 adapters; FEATURE-0014 creates no adapter or runtime. | Deferred | ADH-2026-018; DEC-0036 |

## 11. Stage governance and completeness self-check

This artifact is the single authorized output of the Requirements stage
(architecture section 16.2). It creates or modifies no other file, and no
design or task content is present or implied.

Stage boundary (architecture section 16.1):

- Only `requirements.md` was generated under the approval of the architecture
  and `ADH-2026-018`.
- Design generation requires explicit human `APPROVED_FOR_DESIGN`; task
  generation requires explicit human `APPROVED_FOR_TASKS`.
- No source code, schema, generated client, migration, or test is authorized
  by this stage.

Deterministic completeness self-check (architecture section 16.6):

- Exactly five owned resource kinds; no active alias (F14-REQ-01, F14-REQ-03).
- All 21 decisions `F14-AD-001`–`F14-AD-021` individually traced (§10.1).
- All 30 risks `F14-R01`–`F14-R30` individually traced with preserved
  mitigation, evidence, target residual, owner, and reassessment trigger
  (§10.2).
- Every normative requirement has exactly one owner (§10.3) and one normative
  home per semantic key (§10.4).
- No ResourcePool, ProviderCapability, capacity, eligibility, compatibility,
  adapter, discovery, credential, endpoint, provider-call, plugin, operation,
  provisioning, runtime, or connectivity semantics appear outside marked
  non-goal/ownership text (F14-REQ-19, F14-REQ-20, F14-REQ-24, F14-REQ-25,
  §5).
- FEATURE-0013 adoption is exactly `NOT_APPLICABLE`; no decision/audit
  requirement was introduced by generic prompt text (F14-REQ-23, SEC-07).
- No placeholder (`TBD`, `TODO`, "implementation-defined", "as appropriate")
  and no unresolved semantic marker remains.
- No design-choice vocabulary appears as a normative requirement; delegated
  representation choices are isolated in §9 and marked non-normative.

Escalation rule (architecture sections 14.1, 16.8): an unresolved semantic
choice stops the stage with `ARCHITECTURE_DECISION_REQUIRED`; a repository,
prompt, reviewer, automation-state, or manifest failure stops with
`REPOSITORY_CONTEXT_NOT_READY`. Neither condition was encountered while
authoring these requirements against the approved package.
