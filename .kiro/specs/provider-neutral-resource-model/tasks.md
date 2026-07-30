# FEATURE-0014 Provider-Neutral Resource Model — Tasks

- Feature: FEATURE-0014 — Provider-Neutral Resource Model
- Stage: Tasks
- Phase: Phase 2 — Reuse-First PaaS Fabric Foundation
- Kiro slug: provider-neutral-resource-model
- Controlling architecture: `docs/architecture/provider-neutral-resource-model.md`
- Controlling handoff: `ADH-2026-018` (Approved)
- Controlling requirements: `.kiro/specs/provider-neutral-resource-model/requirements.md`
- Controlling design: `.kiro/specs/provider-neutral-resource-model/design.md`
- FEATURE-0013 downstream adoption: `NOT_APPLICABLE` (architecture section 9)

## 0. Task-stage boundary

These tasks implement only the design section 3.1 `IMPLEMENT` set (I-1…I-7) and
its supporting FEATURE-0012 grammar reuse. No task creates a repository,
persistence engine, storage component, lookup/resolver interface, state-provider
abstraction, live HTTP handler, CRUD executor, controller, reconciler, adapter,
provider integration, capability, capacity, placement, connectivity, decision,
audit, or operation artifact. Behaviors classified `CONTRACT_ONLY / NO_TASK`
(design section 3.1, items C-1…C-8) are proven by offline validators, pure
deterministic helpers over explicitly supplied in-memory resource values, and
conformance fixtures — never by production state code. `EXCLUDED` mechanisms
(design section 3.1, items X-1…X-7) are enforced absent by deny-list and
dependency scans and produce no implementation task (see the No-task ledger).

Package layout, file names, and internal implementation structure below are the
representation choices delegated to design by architecture section 15 and
resolved in design DD-06/DD-07 ("follow the live repository conventions"): flat
per-kind files in `internal/resources` and `internal/validation`, flat schema
files in `api/schemas` composing `api/schemas/_common/` fragments, and
conformance fixtures under `tests/conformance/fixtures/` with checks in
`internal/apiconform`. No task resolves a semantic choice.

All five kinds share `apiVersion: fabric.sovrunn.io/v1alpha1`, singular
PascalCase kinds, `ManagedResource` profile, `operator-facing` boundary, and
stability `alpha` (design DD-01, §4). Reviewed finite limits are exactly those
in design §6.4.

---

## Task 1 — Provider resource representation and FEATURE-0014 type identity

Objective: Define the `Provider` Go type and the shared FEATURE-0014 API
group/version and kind identity constants, with no domain `spec` fields.

Requirements: F14-REQ-01, F14-REQ-04, F14-REQ-05, F14-REQ-21, F14-REQ-22
Design: §4.1, §4.2 (`Provider` row), §4.3, DD-01, DD-02, DD-03
Decisions: F14-AD-001, F14-AD-002, F14-AD-003, F14-AD-014, F14-AD-015
Risks: F14-R01, F14-R02, F14-R04
Implementation class: IMPLEMENT
Files: `internal/resources/provider.go`, `internal/resources/provider_test.go`
Notes:
- Declare the shared constants `apiVersion = "fabric.sovrunn.io/v1alpha1"` and
  the five kind name constants (`Provider`, `ProviderLocation`,
  `ProviderDatacenter`, `DatacenterFailureDomain`, `InfrastructureStack`) in
  this root file; do not duplicate them elsewhere.
- Embed the inherited FEATURE-0012 `internal/apimeta` object metadata and
  status types and the `internal/apicond` condition type; reuse them unchanged
  (F14-REQ-21). Do not redefine metadata, scope, status, or condition grammar.
- `Provider.spec` carries no domain field; governance scope is expressed only
  through inherited `metadata.scopeRef` (Platform or one `Organization`).
- No `ownerRef`-as-scope, no owner/`ProviderOwner`/`SupplyOwner` field, no
  provider-native identifier, no connectivity/capacity/capability field.
- `status` holds `observedGeneration` and current-fact conditions `Valid` and
  `TopologyComplete` only (DD-03); it is system-owned and is not history.
Tests:
- Positive: JSON round-trip preserves `apiVersion`, `kind`, metadata, empty
  spec, and status/conditions; kind and apiVersion constants equal the design
  values.
- Negative: unmarshalling a payload containing any domain `spec` field, an
  owner field, or a connectivity/capacity field is rejectable (field absent
  from the type; struct has no such tag).
- Boundary: `status.conditions` type accommodates the two condition types and
  no more than the §6.4 condition bound is implied by the type contract.
Acceptance criteria: `Provider` type compiles, marshals/unmarshals with the
canonical shape, exposes the shared identity constants, and contains no domain
`spec`, owner, native, connectivity, or capacity field; tests pass under `go
test ./internal/resources/...`.
Commit message: `feat(FEATURE-0014): add Provider resource type and fabric API identity constants`

---

## Task 2 — ProviderLocation resource representation with optional geo descriptor

Objective: Define the `ProviderLocation` Go type with an optional descriptive
`spec.geo` value and no separate parent field.

Requirements: F14-REQ-02, F14-REQ-06, F14-REQ-21, F14-REQ-22, F14-REQ-27
Design: §4.1, §4.2 (`ProviderLocation` row), §4.3, DD-02, DD-04, §6.2
Decisions: F14-AD-006, F14-AD-004, F14-AD-011, F14-AD-014, F14-AD-015
Risks: F14-R02, F14-R16
Implementation class: IMPLEMENT
Files: `internal/resources/providerlocation.go`, `internal/resources/providerlocation_test.go`
Notes:
- Kind is exactly `ProviderLocation`; no alias kind, alias field, or alternate
  location-level vocabulary (F14-REQ-02).
- No parent field: the immediate parent `Provider` is identified only by the
  inherited immutable `metadata.scopeRef` (design §4.2, DD-08).
- `spec.geo` is optional and shaped `{ countryCode, subdivisionCode? }` as
  descriptive declared topology only (DD-04); it carries no residency,
  compliance, authoritative-assignment, or placement meaning. Syntax and
  country-prefix validation belongs to Task 12, not to the type; no dataset or
  membership lookup is authorized.
- Reuse `internal/apimeta`/`internal/apicond` unchanged; add no connectivity,
  capacity, or native field.
Tests:
- Positive: round-trip with and without `spec.geo`; kind constant is
  `ProviderLocation`.
- Negative: no field named for an alternate location term exists; no parent or
  connectivity field exists on the type.
- Boundary: absence of `spec.geo` is representable and carries no default.
Acceptance criteria: `ProviderLocation` type compiles and round-trips, exposes
only optional `spec.geo`, has no parent/alias/connectivity field, and tests
pass.
Commit message: `feat(FEATURE-0014): add ProviderLocation resource type with optional geo descriptor`

---

## Task 3 — ProviderDatacenter resource representation and providerLocationRef

Objective: Define the `ProviderDatacenter` Go type with one immutable typed
`spec.providerLocationRef` constrained to kind `ProviderLocation`.

Requirements: F14-REQ-06, F14-REQ-08, F14-REQ-09, F14-REQ-21, F14-REQ-22
Design: §4.1, §4.2 (`ProviderDatacenter` row), §4.3, DD-02, DD-08
Decisions: F14-AD-004, F14-AD-005, F14-AD-008, F14-AD-014, F14-AD-015
Risks: F14-R12
Implementation class: IMPLEMENT
Files: `internal/resources/providerdatacenter.go`, `internal/resources/providerdatacenter_test.go`
Notes:
- Declare the typed constrained parent-reference alias over `internal/apiref`
  for `providerLocationRef` (kind `ProviderLocation`, optional resolved `uid`),
  reusing FEATURE-0012 reference grammar unchanged (DD-08, F14-REQ-21).
- Exactly one immediate parent reference; the field is the only topology
  authority. `metadata.scopeRef` remains the containing `Provider` and the sole
  governance authority. Immutability of the reference and scope is a contract
  enforced by stateful validation (design §7.1 stages 3–5,
  `CONTRACT_ONLY / NO_TASK`); the type only shapes the single typed reference.
Tests:
- Positive: round-trip with a valid `providerLocationRef`; kind is
  `ProviderDatacenter`.
- Negative: the type admits no second parent field, no cross-kind parent, and
  no connectivity/capacity field.
- Boundary: reference shape matches the `_common/typed-ref` fragment fields.
Acceptance criteria: `ProviderDatacenter` type compiles and round-trips with
exactly one typed `providerLocationRef`; tests pass.
Commit message: `feat(FEATURE-0014): add ProviderDatacenter resource type with providerLocationRef`

---

## Task 4 — DatacenterFailureDomain resource representation and providerDatacenterRef

Objective: Define the `DatacenterFailureDomain` Go type with one immutable typed
`spec.providerDatacenterRef` constrained to kind `ProviderDatacenter`.

Requirements: F14-REQ-06, F14-REQ-08, F14-REQ-09, F14-REQ-21, F14-REQ-22
Design: §4.1, §4.2 (`DatacenterFailureDomain` row), §4.3, DD-02, DD-08
Decisions: F14-AD-004, F14-AD-005, F14-AD-008, F14-AD-014, F14-AD-015
Risks: F14-R12
Implementation class: IMPLEMENT
Files: `internal/resources/datacenterfailuredomain.go`, `internal/resources/datacenterfailuredomain_test.go`
Notes:
- Declare the typed constrained parent-reference alias over `internal/apiref`
  for `providerDatacenterRef` (kind `ProviderDatacenter`), reusing FEATURE-0012
  reference grammar unchanged.
- Exactly one immediate parent reference; `metadata.scopeRef` remains the
  containing `Provider`. This kind records identity and containment only — no
  correlated-risk, availability, quorum, resilience, or placement field
  (architecture §6.4).
Tests:
- Positive: round-trip with a valid `providerDatacenterRef`; kind is
  `DatacenterFailureDomain`.
- Negative: no second/wrong-kind parent field; no connectivity/capacity/
  resilience field on the type.
- Boundary: reference shape matches the `_common/typed-ref` fragment fields.
Acceptance criteria: `DatacenterFailureDomain` type compiles and round-trips
with exactly one typed `providerDatacenterRef`; tests pass.
Commit message: `feat(FEATURE-0014): add DatacenterFailureDomain resource type with providerDatacenterRef`

---

## Task 5 — InfrastructureStack resource representation, parent ref, and technology

Objective: Define the `InfrastructureStack` Go type with one immutable typed
`spec.datacenterFailureDomainRef` (kind `DatacenterFailureDomain`) and an
optional bounded descriptive `spec.technology`.

Requirements: F14-REQ-03, F14-REQ-06, F14-REQ-08, F14-REQ-09, F14-REQ-15, F14-REQ-16, F14-REQ-17, F14-REQ-18, F14-REQ-21, F14-REQ-22
Design: §4.1, §4.2 (`InfrastructureStack` row), §4.3, DD-02, DD-05, DD-08, §6.3
Decisions: F14-AD-003, F14-AD-005, F14-AD-008, F14-AD-009, F14-AD-010, F14-AD-011, F14-AD-014, F14-AD-015, F14-AD-019, F14-AD-021
Risks: F14-R08, F14-R10, F14-R11, F14-R19, F14-R20
Implementation class: IMPLEMENT
Files: `internal/resources/infrastructurestack.go`, `internal/resources/infrastructurestack_test.go`
Notes:
- The only active kind name is `InfrastructureStack`; the superseded stack-kind
  name must not appear as an active type, constant, tag, or alias (F14-REQ-03,
  F14-AD-021).
- Declare the typed constrained parent-reference alias over `internal/apiref`
  for `datacenterFailureDomainRef` (kind `DatacenterFailureDomain`); exactly one
  immutable immediate parent; a stack must not carry more than one failure-
  domain reference (F14-REQ-16).
- Distinct immutable identity is the inherited `metadata.uid` plus scoped
  identity (F14-REQ-15); shared technology never implies shared identity.
- `spec.technology` is an optional single bounded descriptive string; it is not
  identity, a closed enum, or a lookup/reference key. Length and character-set
  validation and normalization belong to Task 16, not to the type.
- No provider-native identifier, SDK object, endpoint, credential, native
  availability-zone ID, capacity, capability, or connectivity field.
Tests:
- Positive: round-trip with a valid `datacenterFailureDomainRef` and with/
  without `spec.technology`; kind is `InfrastructureStack`.
- Negative: the superseded stack-kind name is absent from the type, tags, and
  constants; no second failure-domain reference field; no native/capacity/
  connectivity field.
- Boundary: two stack values with identical `technology` retain distinct `uid`
  and distinct parent references (identity is not technology).
Acceptance criteria: `InfrastructureStack` type compiles and round-trips with
exactly one failure-domain reference and an optional descriptive technology
string, uses only the active kind name, and tests pass.
Commit message: `feat(FEATURE-0014): add InfrastructureStack resource type with failure-domain ref and technology`

---

## Task 6 — Provider canonical JSON Schema

Objective: Author the canonical `Provider` JSON Schema composing the shared
FEATURE-0012 fragments and adding no domain `spec` field.

Requirements: F14-REQ-01, F14-REQ-05, F14-REQ-18, F14-REQ-20, F14-REQ-21, F14-REQ-24, F14-REQ-28
Design: DD-01, DD-02, DD-07, §4.2, §5, §6.4
Decisions: F14-AD-001, F14-AD-003, F14-AD-014, F14-AD-017, F14-AD-019
Risks: F14-R01, F14-R03, F14-R08
Implementation class: IMPLEMENT
Files: `api/schemas/provider.json`
Notes:
- JSON Schema 2020-12 composing `api/schemas/_common/` fragments (`type-meta`,
  `object-meta`, `scope-ref`, `condition`), adding no new common grammar
  (F14-REQ-21).
- Constrain `metadata.scopeRef.kind` to exactly `Platform` or `Organization`;
  no seventh scope kind, no owner field as scope.
- `spec` has no domain property. Reject unknown and duplicate fields
  (`additionalProperties: false` per fragment convention).
- Encode the §6.4 finite bounds already carried by the shared fragments; add no
  connectivity, capacity, capability, adapter, or native field.
Tests:
- Executed by the Task 19 conformance harness against this schema; a JSON
  schema hosts no Go test itself.
- Positive: a minimal `Provider` with `scopeRef.kind: Platform` and with
  `scopeRef.kind: Organization` validates.
- Negative: `scopeRef.kind: Tenant` (or any non-Platform/Organization value),
  an owner field used as scope, an unknown/duplicate field, and any domain
  `spec` property each fail validation.
- Boundary: `metadata.name` at 253 chars validates and 254 fails; the schema
  resolves all `_common` `$ref`s.
Acceptance criteria: `provider.json` validates as 2020-12, resolves all
`_common` `$ref`s, constrains `scopeRef.kind` to Platform/Organization, exposes
no domain `spec` field, and is accepted by the Task 19 registry/type-binding.
Commit message: `feat(FEATURE-0014): add Provider canonical JSON Schema`

---

## Task 7 — ProviderLocation canonical JSON Schema

Objective: Author the canonical `ProviderLocation` JSON Schema with the optional
descriptive `spec.geo` and no parent field.

Requirements: F14-REQ-02, F14-REQ-06, F14-REQ-18, F14-REQ-20, F14-REQ-21, F14-REQ-27, F14-REQ-28
Design: DD-01, DD-02, DD-04, DD-07, §4.2, §6.2, §6.4
Decisions: F14-AD-004, F14-AD-006, F14-AD-011, F14-AD-014, F14-AD-019
Risks: F14-R02, F14-R16
Implementation class: IMPLEMENT
Files: `api/schemas/provider-location.json`
Notes:
- Kind fixed to `ProviderLocation`; collection `provider-locations` (design §5).
- Constrain `metadata.scopeRef.kind` to exactly `Provider`.
- `spec.geo` optional object: `countryCode` `^[A-Z]{2}$` (exactly 2 chars),
  optional `subdivisionCode` `^[A-Z]{2}-[A-Z0-9]{1,3}$` (≤ 6 chars) whose
  country prefix must equal `countryCode` (design §6.2). Format only; no
  authoritative membership is asserted or checked.
- No alias location term, no parent field, no connectivity/capacity/native
  field; reject unknown/duplicate fields; §6.4 bounds.
Tests:
- Executed by the Task 19 conformance harness against this schema.
- Positive: a `ProviderLocation` with `scopeRef.kind: Provider`, with and
  without a well-formed `spec.geo`, validates.
- Negative: a non-`Provider` scope kind, an alternate location-term field, a
  parent field, `countryCode` not `^[A-Z]{2}$`, a `subdivisionCode` whose
  country prefix differs from `countryCode`, or any unknown field fails.
- Boundary: `countryCode` exactly 2 chars validates; `subdivisionCode` at 6
  chars validates and 7 fails.
Acceptance criteria: `provider-location.json` validates as 2020-12, constrains
`scopeRef.kind` to `Provider`, encodes the geo format constraints, has no
alias/parent/connectivity field, and is accepted by the Task 19 registry.
Commit message: `feat(FEATURE-0014): add ProviderLocation canonical JSON Schema`

---

## Task 8 — ProviderDatacenter canonical JSON Schema

Objective: Author the canonical `ProviderDatacenter` JSON Schema with the
required immutable `spec.providerLocationRef`.

Requirements: F14-REQ-06, F14-REQ-08, F14-REQ-09, F14-REQ-18, F14-REQ-20, F14-REQ-21, F14-REQ-28
Design: DD-01, DD-02, DD-07, DD-08, §4.2, §5, §6.4
Decisions: F14-AD-004, F14-AD-005, F14-AD-008, F14-AD-014, F14-AD-019
Risks: F14-R12
Implementation class: IMPLEMENT
Files: `api/schemas/provider-datacenter.json`
Notes:
- Kind `ProviderDatacenter`; collection `provider-datacenters`.
- Constrain `metadata.scopeRef.kind` to `Provider`.
- `spec.providerLocationRef` required, composing `_common/typed-ref`, with
  `kind` fixed to `ProviderLocation`; exactly one reference (no array).
- No connectivity/capacity/native field; reject unknown/duplicate fields; §6.4
  bounds.
Tests:
- Executed by the Task 19 conformance harness against this schema.
- Positive: a `ProviderDatacenter` with `scopeRef.kind: Provider` and one
  `providerLocationRef` of `kind: ProviderLocation` validates.
- Negative: a missing reference, two references (array), a wrong-kind
  reference, a non-`Provider` scope, or any unknown field fails.
- Boundary: exactly one reference is accepted; zero and two are rejected.
Acceptance criteria: `provider-datacenter.json` validates as 2020-12, requires a
single `providerLocationRef` constrained to `ProviderLocation`, and is accepted
by the Task 19 registry.
Commit message: `feat(FEATURE-0014): add ProviderDatacenter canonical JSON Schema`

---

## Task 9 — DatacenterFailureDomain canonical JSON Schema

Objective: Author the canonical `DatacenterFailureDomain` JSON Schema with the
required immutable `spec.providerDatacenterRef`.

Requirements: F14-REQ-06, F14-REQ-08, F14-REQ-09, F14-REQ-18, F14-REQ-19, F14-REQ-20, F14-REQ-21, F14-REQ-28
Design: DD-01, DD-02, DD-07, DD-08, §4.2, §5, §6.4
Decisions: F14-AD-004, F14-AD-005, F14-AD-008, F14-AD-012, F14-AD-013, F14-AD-014, F14-AD-019
Risks: F14-R06, F14-R07, F14-R12
Implementation class: IMPLEMENT
Files: `api/schemas/datacenter-failure-domain.json`
Notes:
- Kind `DatacenterFailureDomain`; collection `datacenter-failure-domains`.
- Constrain `metadata.scopeRef.kind` to `Provider`.
- `spec.providerDatacenterRef` required, composing `_common/typed-ref`, `kind`
  fixed to `ProviderDatacenter`; exactly one reference.
- Deny-list: no `connected` boolean, adjacency, route, peer, reachability,
  latency, bandwidth, trust, health, quorum, correlated-risk, or resilience
  field (F14-REQ-19, F14-REQ-20). Reject unknown/duplicate fields; §6.4 bounds.
Tests:
- Executed by the Task 19 conformance harness against this schema.
- Positive: a `DatacenterFailureDomain` with `scopeRef.kind: Provider` and one
  `providerDatacenterRef` of `kind: ProviderDatacenter` validates.
- Negative: a missing/duplicate/wrong-kind reference, a non-`Provider` scope,
  or any connectivity/isolation/resilience field fails.
- Boundary: exactly one reference accepted; the connectivity deny-list scan
  finds no matching property anywhere in the schema.
Acceptance criteria: `datacenter-failure-domain.json` validates as 2020-12,
requires a single `providerDatacenterRef`, contains no connectivity/isolation
field, and is accepted by the Task 19 registry.
Commit message: `feat(FEATURE-0014): add DatacenterFailureDomain canonical JSON Schema`

---

## Task 10 — InfrastructureStack canonical JSON Schema

Objective: Author the canonical `InfrastructureStack` JSON Schema with the
required immutable `spec.datacenterFailureDomainRef` and optional bounded
`spec.technology`.

Requirements: F14-REQ-03, F14-REQ-08, F14-REQ-09, F14-REQ-15, F14-REQ-16, F14-REQ-17, F14-REQ-18, F14-REQ-20, F14-REQ-21, F14-REQ-24, F14-REQ-28
Design: DD-01, DD-02, DD-05, DD-07, DD-08, §4.2, §5, §6.3, §6.4
Decisions: F14-AD-005, F14-AD-008, F14-AD-009, F14-AD-010, F14-AD-011, F14-AD-014, F14-AD-017, F14-AD-019, F14-AD-021
Risks: F14-R10, F14-R11, F14-R19, F14-R20
Implementation class: IMPLEMENT
Files: `api/schemas/infrastructure-stack.json`
Notes:
- Active kind name only: `InfrastructureStack`; collection
  `infrastructure-stacks`. The superseded stack-kind name must not appear as an
  active `kind`, `$id`, `title`, property, or enum value (F14-REQ-03).
- Constrain `metadata.scopeRef.kind` to `Provider`.
- `spec.datacenterFailureDomainRef` required, composing `_common/typed-ref`,
  `kind` fixed to `DatacenterFailureDomain`; exactly one reference (no array,
  no spanning).
- `spec.technology` optional string `^[A-Za-z0-9][A-Za-z0-9 ._+-]{0,99}$`
  (≤ 100 chars); not an enum, not identity, not a reference key (F14-REQ-17).
- No capability/capacity/placement/adapter/native/connectivity field; reject
  unknown/duplicate fields; §6.4 bounds.
Tests:
- Executed by the Task 19 conformance harness against this schema.
- Positive: an `InfrastructureStack` with `scopeRef.kind: Provider`, one
  `datacenterFailureDomainRef` of `kind: DatacenterFailureDomain`, and with/
  without a valid `technology` validates.
- Negative: a missing/duplicate/wrong-kind reference, a non-`Provider` scope,
  the superseded stack-kind name as `kind`/`$id`/property/enum, a `technology`
  over 100 chars or with a disallowed character, or any capability/capacity/
  connectivity field fails.
- Boundary: `technology` at exactly 100 chars validates and 101 fails; exactly
  one failure-domain reference is accepted.
Acceptance criteria: `infrastructure-stack.json` validates as 2020-12, uses only
the active kind name, requires a single `datacenterFailureDomainRef`, encodes
`technology` as a bounded non-enum string, and is accepted by the Task 19
registry.
Commit message: `feat(FEATURE-0014): add InfrastructureStack canonical JSON Schema`

---

## Task 12 — ProviderLocation offline validator (structural + semantic, incl. geo)

Objective: Implement the offline structural and semantic validator for
`ProviderLocation`, including `scopeRef.kind`, bounds, status-not-user-authored,
and deterministic geo syntax plus country-prefix agreement.

Requirements: F14-REQ-02, F14-REQ-06, F14-REQ-21, F14-REQ-22, F14-REQ-27, F14-REQ-28, F14-REQ-31
Design: §7.1 (stages 1–2, offline IMPLEMENT), DD-04, §6.2, §6.4, §4.3
Decisions: F14-AD-004, F14-AD-006, F14-AD-011, F14-AD-014, F14-AD-015, F14-AD-019
Risks: F14-R02, F14-R08, F14-R16, F14-R17
Implementation class: IMPLEMENT
Files: `internal/validation/providerlocation.go`, `internal/validation/providerlocation_test.go`
Notes:
- Reuse `internal/apivalid` structural primitives and `DefaultLimits` for common
  bounds (§6.4); add no new grammar.
- Structural: required fields, DNS-style `metadata.name`, unknown/duplicate-
  field rejection, bounded sizes, and rejection of user-authored `status`,
  `resourceVersion`, `generation`, and timestamps (F14-REQ-22, guardrails §10).
- Semantic: `metadata.scopeRef.kind` must be `Provider`; `spec.geo` (when
  present) must pass the §6.2 format and country-prefix agreement checks.
  Malformed or prefix-inconsistent values are rejected; syntactically valid
  unassigned values are accepted without authoritative meaning (F14-REQ-27).
- Emit inherited RFC 9457 Problem Details with stable codes and RFC 6901 paths;
  responses/errors carry no native identifiers or secrets (F14-REQ-31).
- This is a pure offline validator over a supplied resource value; it performs
  no lookup of other resources and owns no state.
Tests:
- Positive: a valid `ProviderLocation` with and without a valid `geo` passes.
- Negative: wrong `scopeRef.kind`; user-authored `status`; unknown field;
  malformed or prefix-inconsistent geo descriptor — each rejected with the
  correct stable code and pointer.
- Boundary: a syntactically valid unassigned descriptor is accepted and carries
  no authoritative, residency, compliance, or placement inference.
- Boundary: geo `countryCode` exactly 2 chars; `subdivisionCode` ≤ 6 chars with
  matching country prefix; label/annotation/name bounds at the §6.4 edges.
Acceptance criteria: the validator accepts conforming and syntactically valid
unassigned values, rejects malformed and prefix-inconsistent values with the
inherited stable error contract, performs no dataset/library lookup, and tests
pass.
Commit message: `feat(FEATURE-0014): add ProviderLocation offline validator with geo syntax checks`

---

## Task 13 — Provider offline validator (structural + semantic)

Objective: Implement the offline structural and semantic validator for
`Provider`, enforcing Platform/Organization scope and the empty domain spec.

Requirements: F14-REQ-04, F14-REQ-05, F14-REQ-21, F14-REQ-22, F14-REQ-28, F14-REQ-31
Design: section 7.1 (stages 1–2), §4.2, §4.3, §6.4
Decisions: F14-AD-002, F14-AD-003, F14-AD-014, F14-AD-015, F14-AD-019
Risks: F14-R03, F14-R04
Implementation class: IMPLEMENT
Files: `internal/validation/provider.go`, `internal/validation/provider_test.go`
Notes:
- Structural: required fields, DNS-style name, unknown/duplicate-field
  rejection, bounds, status-not-user-authored (§6.4, F14-REQ-22).
- Semantic: `metadata.scopeRef.kind` must be exactly `Platform` or
  `Organization`; reject any owner field or `ownerRef` used as governance scope
  (F14-REQ-05); `spec` must contain no domain field.
- Reuse `internal/apivalid`; emit inherited Problem Details; redact native
  values/secrets (F14-REQ-31). Pure offline validator; no cross-resource lookup.
Tests:
- Positive: Platform-scoped and Organization-scoped `Provider` both pass.
- Negative: `scopeRef.kind` other than Platform/Organization; owner field or
  `ownerRef`-as-scope; any domain `spec` field; user-authored `status` — each
  rejected.
- Boundary: name/label/annotation bounds at §6.4 edges.
Acceptance criteria: the validator enforces Platform/Organization sole scope and
the empty domain spec with the inherited error contract; tests pass.
Commit message: `feat(FEATURE-0014): add Provider offline validator with scope enforcement`

---

## Task 14 — ProviderDatacenter offline validator (structural + semantic)

Objective: Implement the offline validator for `ProviderDatacenter`, enforcing
`Provider` scope and a single required `providerLocationRef` of the correct
kind.

Requirements: F14-REQ-06, F14-REQ-08, F14-REQ-09, F14-REQ-21, F14-REQ-22, F14-REQ-28, F14-REQ-31
Design: §7.1 (stages 1–2), DD-08, §4.2, §4.3, §6.4
Decisions: F14-AD-004, F14-AD-005, F14-AD-008, F14-AD-014, F14-AD-015, F14-AD-019
Risks: F14-R12
Implementation class: IMPLEMENT
Files: `internal/validation/providerdatacenter.go`, `internal/validation/providerdatacenter_test.go`
Notes:
- Structural + status-not-user-authored as in Task 13.
- Semantic: `metadata.scopeRef.kind` must be `Provider`; exactly one
  `spec.providerLocationRef` with `kind == ProviderLocation`; zero, multiple, or
  wrong-kind immediate parents are rejected (F14-REQ-08, F14-REQ-09). Skipped-
  level attachment (for example a datacenter referencing a `Provider` or a
  failure domain) is rejected by kind constraint.
- Live parent existence and same-`Provider` resolution are stateful
  (`CONTRACT_ONLY / NO_TASK`, design §7.1 stage 3) and are handled by Tasks
  16–17 helpers and Task 20 fixtures, not here.
Tests:
- Positive: a datacenter with one valid `providerLocationRef` passes.
- Negative: missing reference; two references; wrong-kind reference; wrong
  `scopeRef.kind`; user-authored status — each rejected.
- Boundary: reference field cardinality is exactly one.
Acceptance criteria: the validator enforces `Provider` scope and a single
correctly typed parent reference with the inherited error contract; tests pass.
Commit message: `feat(FEATURE-0014): add ProviderDatacenter offline validator`

---

## Task 15 — DatacenterFailureDomain offline validator (structural + semantic)

Objective: Implement the offline validator for `DatacenterFailureDomain`,
enforcing `Provider` scope and a single required `providerDatacenterRef`.

Requirements: F14-REQ-06, F14-REQ-08, F14-REQ-09, F14-REQ-19, F14-REQ-20, F14-REQ-21, F14-REQ-22, F14-REQ-28, F14-REQ-31
Design: §7.1 (stages 1–2), DD-08, §4.2, §4.3, §6.4, §8
Decisions: F14-AD-004, F14-AD-005, F14-AD-008, F14-AD-012, F14-AD-013, F14-AD-014, F14-AD-015, F14-AD-019
Risks: F14-R06, F14-R07, F14-R12
Implementation class: IMPLEMENT
Files: `internal/validation/datacenterfailuredomain.go`, `internal/validation/datacenterfailuredomain_test.go`
Notes:
- Structural + status-not-user-authored as in Task 13.
- Semantic: `scopeRef.kind` must be `Provider`; exactly one
  `spec.providerDatacenterRef` with `kind == ProviderDatacenter`; zero/multiple/
  wrong-kind rejected. Confirm the type/schema carries no connectivity field
  (defense in depth for F14-REQ-19/F14-REQ-20; primary evidence is the Task 9
  schema and Task 20 deny-list).
Tests:
- Positive: a failure domain with one valid `providerDatacenterRef` passes.
- Negative: missing/multiple/wrong-kind reference; wrong scope; user-authored
  status — each rejected.
- Boundary: exactly one reference; no connectivity field accepted.
Acceptance criteria: the validator enforces `Provider` scope and a single
correctly typed parent reference and introduces no connectivity semantics; tests
pass.
Commit message: `feat(FEATURE-0014): add DatacenterFailureDomain offline validator`

---

## Task 16 — InfrastructureStack offline validator (structural + semantic, incl. technology)

Objective: Implement the offline validator for `InfrastructureStack`, enforcing
`Provider` scope, a single required `datacenterFailureDomainRef`, and bounded
descriptive `technology` with deterministic normalization.

Requirements: F14-REQ-06, F14-REQ-08, F14-REQ-09, F14-REQ-15, F14-REQ-16, F14-REQ-17, F14-REQ-18, F14-REQ-21, F14-REQ-22, F14-REQ-28, F14-REQ-31
Design: §7.1 (stages 1–2), DD-05, DD-08, §4.2, §4.3, §6.3, §6.4
Decisions: F14-AD-005, F14-AD-008, F14-AD-009, F14-AD-010, F14-AD-011, F14-AD-014, F14-AD-015, F14-AD-019, F14-AD-021
Risks: F14-R10, F14-R11, F14-R17, F14-R20
Implementation class: IMPLEMENT
Files: `internal/validation/infrastructurestack.go`, `internal/validation/infrastructurestack_test.go`
Notes:
- Structural + status-not-user-authored as in Task 13.
- Semantic: `scopeRef.kind` must be `Provider`; exactly one
  `spec.datacenterFailureDomainRef` with `kind == DatacenterFailureDomain`;
  zero/multiple/cross-kind rejected; no spanning (F14-REQ-16).
- `spec.technology` (when present): apply the design §6.3 deterministic
  normalization in the fixed order (reject non-ASCII outside the allowed set;
  trim leading/trailing ASCII spaces; collapse internal space runs; preserve
  case) then validate `^[A-Za-z0-9][A-Za-z0-9 ._+-]{0,99}$`; store verbatim; do
  not match against an enum or use as identity/lookup key (F14-REQ-17,
  F14-REQ-15).
- The superseded stack-kind name must not be introduced by this validator.
Tests:
- Positive: a stack with one valid parent reference and with/without a valid
  normalized `technology` passes; two stacks with identical technology validate
  independently (identity is `uid`, not technology).
- Negative: missing/multiple/cross-kind parent; wrong scope; `technology` with
  a control/non-ASCII character, over 100 chars, or leading/trailing space after
  normalization; user-authored status — each rejected.
- Boundary: `technology` exactly 100 chars accepted; 101 rejected; single
  internal space preserved, double collapsed.
Acceptance criteria: the validator enforces a single failure-domain parent and
bounded descriptive technology with deterministic normalization and the
inherited error contract; tests pass.
Commit message: `feat(FEATURE-0014): add InfrastructureStack offline validator with bounded technology`

---

## Task 17 — Deterministic scope-agreement and hierarchy-ordering helper

Objective: Implement a pure deterministic helper that evaluates same-`Provider`
scope agreement and hierarchy-level ordering over explicitly supplied in-memory
resource values, with no lookup, state, or persistence.

Requirements: F14-REQ-07, F14-REQ-08, F14-REQ-10, F14-REQ-13, F14-REQ-30
Design: section 1.3, §3 (component responsibilities), §3.1 (I-5; C-2/C-8 CONTRACT_ONLY), §4.4, §7.1 (stage 3, CONTRACT_ONLY)
Decisions: F14-AD-004, F14-AD-005, F14-AD-008
Risks: F14-R03, F14-R12, F14-R24
Implementation class: IMPLEMENT
Files: `internal/validation/topology.go`, `internal/validation/topology_test.go`
Notes:
- The helper takes a parent value and a child value already supplied by the
  caller (or fixtures) and evaluates: (a) both resolve to the same `Provider`
  scope UID; (b) the child's immediate parent kind is the required kind for the
  child kind (hierarchy ordering, no skips/reorders). It performs no lookup,
  reads no store, and owns no state (design §1.3, I-5).
- Cross-`Provider` parent references and parent/child scope-UID mismatches
  evaluate to a rejection outcome; the live rejection response and no-existence
  disclosure against authoritative state are `CONTRACT_ONLY / NO_TASK` (design
  §3.1, C-2/C-8) proven by Task 20 fixtures.
- Multi-owner isolation (F14-REQ-30) is expressed as independent scope UIDs; the
  helper never keys on external operator name or native ID.
Tests:
- Positive: same-`Provider` parent/child with correct hierarchy ordering
  evaluates to agreement.
- Negative: different `Provider` scope UIDs; correct scope but wrong hierarchy
  level (skip/reorder) — each evaluates to rejection.
- Boundary/isolation: two providers under the same owner `Organization` are
  distinct scope UIDs and do not agree across providers; same external operator
  under two owners yields independent scope UIDs.
- Concurrency: the helper is a pure function and is safe under concurrent
  invocation (verified with a `-race` parallel test over supplied values).
Acceptance criteria: the helper deterministically decides scope agreement and
hierarchy ordering over supplied values, performs no lookup/state access, and
tests pass including under `-race`.
Commit message: `feat(FEATURE-0014): add deterministic scope-agreement and hierarchy helper`

---

## Task 18 — Deterministic topology-completeness helper

Objective: Implement a pure deterministic helper that evaluates registered-vs-
complete topology and the `TopologyComplete` condition over an explicitly
supplied in-memory set of resource values.

Requirements: F14-REQ-12, F14-REQ-13, F14-REQ-14, F14-REQ-26, F14-REQ-29
Design: §1.3, §3.1 (I-5; C-3/C-4/C-5/C-7 CONTRACT_ONLY), §4.4, §7.1 (stage 6, CONTRACT_ONLY), §7.3, §8, DD-03
Decisions: F14-AD-007, F14-AD-020
Risks: F14-R09, F14-R13, F14-R18, F14-R23
Implementation class: IMPLEMENT
Files: `internal/validation/completeness.go`, `internal/validation/completeness_test.go`
Notes:
- Given a supplied set of resources, deterministically compute whether an
  unbroken `Provider` → `ProviderLocation` → `ProviderDatacenter` →
  `DatacenterFailureDomain` → `InfrastructureStack` path exists beneath a
  resource; a leaf `InfrastructureStack` is complete when it is itself `Valid`
  (DD-03). Completeness carries no capability/capacity/placement/readiness
  meaning (F14-REQ-14).
- A zero-child registered parent evaluates as incomplete but valid (F14-REQ-12,
  F14-REQ-13); skips are never complete.
- Provide a child-existence evaluation over the supplied set as the contract
  basis for deletion rejection while children exist (F14-REQ-26); the live
  deletion executor and authoritative child query are `CONTRACT_ONLY / NO_TASK`
  (design §3.1, C-3/C-5), proven by Task 20 fixtures.
- Recomputation is deterministic and order-independent so completeness converges
  after concurrent create/delete (F14-REQ-29, design §7.3); live persistence-
  backed recomputation is `CONTRACT_ONLY / NO_TASK` (C-4/C-7).
Tests:
- Positive: a full five-level supplied path is complete; leaf stack `Valid` is
  complete.
- Negative: any missing level is incomplete; an empty registered parent is
  incomplete, not complete; a supplied parent with an existing child yields a
  child-existence (delete-blocked) outcome.
- Boundary: completeness result is independent of input ordering.
- Concurrency: pure function safe under `-race` parallel evaluation; repeated
  evaluation over the same set is stable.
Acceptance criteria: the helper deterministically distinguishes registered from
complete, evaluates child existence for deletion protection, is order- and
concurrency-independent, and tests pass including under `-race`.
Commit message: `feat(FEATURE-0014): add deterministic topology-completeness helper`

---

## Task 19 — Schema registration and type-binding conformance

Objective: Register the five canonical schemas and bind them to the five Go
types in the inherited conformance harness, proving exactly-five-kinds, grammar
reuse, and provider neutrality.

Requirements: F14-REQ-01, F14-REQ-02, F14-REQ-03, F14-REQ-18, F14-REQ-21, F14-REQ-22
Design: §3 (inherited primitives), §3.1 (I-6), DD-01, DD-07, §5, §9.2
Decisions: F14-AD-001, F14-AD-006, F14-AD-014, F14-AD-019, F14-AD-021
Risks: F14-R02, F14-R08, F14-R19, F14-R28
Implementation class: IMPLEMENT
Files: `internal/apiconform/feature0014_bindings.go`, `internal/apiconform/feature0014_bindings_test.go`
Notes:
- Register the five schemas with the existing `internal/apiconform`
  schema/type-binding mechanism (following the existing binding conventions in
  `bindings.go`/`schemaregistry.go`); reuse the shared primitives (`apimeta`,
  `apiref`, `apicond`, `apiproblem`, `apischema`, `apivalid`) unchanged
  (F14-REQ-21).
- Assert the FEATURE-0014 kind inventory resolves to exactly the five kinds and
  their collections (design §5), with no alias kind and no active superseded
  stack-kind name (F14-REQ-01, F14-REQ-02, F14-REQ-03).
- Reuse the existing provider-neutrality conformance to assert no native
  identifier/SDK/endpoint appears in the five contracts (F14-REQ-18).
Tests:
- Positive: each schema binds to its Go type; canonical shape and status/
  condition grammar validate.
- Negative: a sixth/alias kind, an alternate location term, or the superseded
  stack-kind name fails the inventory assertion.
- Boundary: the inventory count is exactly five.
Acceptance criteria: all five schema/type bindings register and pass, the kind
inventory is exactly five with no alias, and provider-neutrality holds; tests
pass.
Commit message: `feat(FEATURE-0014): register schemas and type bindings in conformance harness`

---

## Task 20 — Positive conformance fixtures for owner/operator, hierarchy, and completeness

Objective: Add valid conformance fixtures and tests covering owner/operator
scenarios, the exact hierarchy, empty registered parents, completeness, and
distinct stack identity, using explicitly supplied in-memory state.

Requirements: F14-REQ-04, F14-REQ-05, F14-REQ-08, F14-REQ-11, F14-REQ-12, F14-REQ-13, F14-REQ-14, F14-REQ-15, F14-REQ-16, F14-REQ-26, F14-REQ-27, F14-REQ-30
Design: section 1.3, §3.1 (I-6/I-7; C-1…C-8 proven by fixtures), §9.2, §4.4
Decisions: F14-AD-002, F14-AD-003, F14-AD-005, F14-AD-007, F14-AD-009, F14-AD-010, F14-AD-020
Risks: F14-R04, F14-R05, F14-R09, F14-R11, F14-R13, F14-R24
Implementation class: IMPLEMENT
Files: `tests/conformance/fixtures/provider.json`, `tests/conformance/fixtures/provider-location.json`, `tests/conformance/fixtures/provider-datacenter.json`, `tests/conformance/fixtures/datacenter-failure-domain.json`, `tests/conformance/fixtures/infrastructure-stack.json`, `internal/apiconform/feature0014_positive_test.go`
Notes:
- Include a distinct-owner/operator case (`Organization: OwnerOrganization-A` →
  `Provider: ProviderOperator-A` and `Provider: ProviderOperator-B`) and a
  same-party case (`Organization: ProviderOperator-A` → `Provider:
  ProviderOperator-A`), proving two role-specific resources with distinct UIDs
  (F14-REQ-11).
- Cover empty registered parents at each level accepted as incomplete, a full
  five-level path reported complete, two same-technology stacks with distinct
  UIDs/parents, and leaf-first deletion of a complete path succeeding
  (F14-REQ-12–16, F14-REQ-26) — all evaluated with the Task 17/18 helpers over
  supplied fixture state (the live handlers/store are `CONTRACT_ONLY / NO_TASK`).
- Include a valid `geo` code fixture carrying no residency/compliance inference
  (F14-REQ-27) and a multi-owner isolation fixture (F14-REQ-30).
Tests:
- Positive: each fixture validates against its schema and validator; owner/
  operator scenarios resolve to independent scoped identities; the full path is
  complete; leaf-first deletion sequence is accepted.
- Boundary: empty parent at each level is valid-but-incomplete; two identical-
  technology stacks retain distinct UIDs.
- Isolation: same-operator/different-owner fixtures are independent scope UIDs.
Acceptance criteria: all positive fixtures validate and the scenario assertions
hold over supplied in-memory state without any production store; tests pass.
Commit message: `test(FEATURE-0014): add positive conformance fixtures for topology scenarios`

---

## Task 21 — Negative, boundary, isolation, and deny-list conformance fixtures

Objective: Add negative and deny-list conformance fixtures and tests covering
skipped levels, multiple/cross-scope parents, cross-provider no-existence
disclosure, connectivity absence, excluded adjacent-feature semantics, geo/
technology rejection, deletion-with-children, and stale-version conflict.

Requirements: F14-REQ-06, F14-REQ-07, F14-REQ-08, F14-REQ-09, F14-REQ-10, F14-REQ-16, F14-REQ-17, F14-REQ-19, F14-REQ-20, F14-REQ-22, F14-REQ-23, F14-REQ-24, F14-REQ-25, F14-REQ-26, F14-REQ-27, F14-REQ-29, F14-REQ-30, F14-REQ-31
Design: section 3.1 (C-1…C-8, X-1…X-7), §7.1 (stages 3–6, CONTRACT_ONLY), §7.2, §9.2, §13
Decisions: F14-AD-004, F14-AD-005, F14-AD-008, F14-AD-010, F14-AD-011, F14-AD-012, F14-AD-013, F14-AD-015, F14-AD-016, F14-AD-017, F14-AD-018, F14-AD-019, F14-AD-020
Risks: F14-R01, F14-R06, F14-R07, F14-R08, F14-R11, F14-R12, F14-R14, F14-R15, F14-R16, F14-R20, F14-R21, F14-R22, F14-R24, F14-R30
Implementation class: IMPLEMENT
Files: `tests/conformance/fixtures/negative/feature0014/` (negative fixture set), `internal/apiconform/feature0014_negative_test.go`
Notes:
- Negative topology: skipped-level attachment at every boundary; zero/multiple/
  wrong-kind immediate parents; a stack referencing multiple failure domains;
  parent/child scope-UID mismatch (F14-REQ-06–10, F14-REQ-16).
- Isolation: cross-provider parent reference and cross-scope get/list/reference
  return the same response as a genuinely absent target, disclosing no existence
  (F14-REQ-07, F14-REQ-30, F14-REQ-31).
- Connectivity deny-list: schema/field scan finds no `connected` boolean,
  adjacency, route, peer, reachability, latency, bandwidth, trust, or health
  field on any of the five kinds; same/cross-datacenter, cross-location, and
  cross-provider (including two providers under one owner `Organization`)
  relationships carry no connectivity inference (F14-REQ-19, F14-REQ-20).
- Adjacent-feature deny-list: field/interface scan finds no ResourcePool,
  ProviderCapability, capacity, eligibility, compatibility, placement,
  adapter, discovery, credential, endpoint, provider-call, plugin, operation,
  provisioning, runtime, repository, or persistence artifact, and no governed
  FEATURE-0013 concept (`DecisionRecord`, `DecisionProfile`, `EvaluationResult`,
  `AuditEvent`, rationale, obligation, actor, subject, linkage, projection,
  composition) outside non-goal/ownership text (F14-REQ-23, F14-REQ-24,
  F14-REQ-25).
- Field handling: malformed or country-prefix-inconsistent geo descriptors are
  rejected while syntactically valid unassigned descriptors are accepted
  without authoritative meaning; malformed/over-length/non-ASCII technology and
  status-history in `status` are rejected (F14-REQ-22, F14-REQ-27, F14-REQ-17).
- Deletion/concurrency contract (proven over supplied state; live executor is
  `CONTRACT_ONLY / NO_TASK`): deleting a parent with existing children yields a
  child-existence (delete-blocked) conflict with no cascade or silent
  reparenting; a stale-version write is rejected with a conflict; a deleted UID
  is not reused (F14-REQ-26, F14-REQ-29).
Tests:
- Negative: each fixture is rejected with the correct inherited stable code and
  RFC 6901 pointer, or (for cross-scope) with the no-existence-disclosure
  response.
- Boundary: geo/technology edge violations; exactly-one-parent violations.
- Isolation: cross-provider/cross-owner denials reveal no target existence and
  no native value.
- Concurrency: stale-version and concurrent delete/create over supplied state
  resolve deterministically (verified under `-race`).
Acceptance criteria: every negative, deny-list, isolation, and concurrency
fixture behaves as specified, no excluded semantic or connectivity field is
present, and tests pass including under `-race`.
Commit message: `test(FEATURE-0014): add negative and deny-list conformance fixtures`

---

## Task 22 — Repository-approved verification, guardrail, boundary, and artifact checks

Objective: Run the approved Go/Docker format, vet, test, race, build,
guardrail, boundary, artifact-cleanup, and changed-file checks for FEATURE-0014
and confirm they pass for this feature's committed changes.

Requirements: F14-REQ-01, F14-REQ-03, F14-REQ-18, F14-REQ-21, F14-REQ-23, F14-REQ-24, F14-REQ-25, F14-REQ-28, F14-REQ-29, F14-REQ-31
Design: section 3, §3.1 (X-1…X-7 enforced absent), §9.2, §15
Decisions: F14-AD-001, F14-AD-014, F14-AD-016, F14-AD-017, F14-AD-018, F14-AD-019, F14-AD-021
Risks: F14-R19, F14-R22, F14-R25, F14-R26, F14-R28, F14-R29, F14-R30
Implementation class: IMPLEMENT
Files: none (runs approved verification commands; creates or modifies no source file)
Notes:
- Run, from the engineering guardrails and `Makefile`, in this order: `make fmt`;
  `make vet`; `make test`; `make test-race` (equivalently `go test -race
  ./...`); `make build`; `make ff-guardrails FEATURE=FEATURE-0014`; and
  `make feature-0014-architecture-boundary-check STAGE=tasks`. Also run the
  Docker verification `./scripts/verify.sh` (image `golang:1.22`, which runs
  `gofmt -l .`, `go vet ./...`, `go test ./...`).
- Artifact cleanup before the changed-file check: `find . -name ".DS_Store"
  -delete` and remove `bin/` and any `sovrunn-api` binary so the guardrail
  build-artifact check passes; do not commit `site/`, generated prompts, or zip
  archives.
- Changed-file check: confirm the working tree changes are limited to the
  FEATURE-0014 approved paths from Tasks 1–10 and Tasks 12–21 (`internal/resources/*`,
  `internal/validation/*`, `api/schemas/*.json`, `internal/apiconform/*`,
  `tests/conformance/fixtures/*`) and that `gofmt -l .` reports no unformatted
  file. `internal/api` must not import `internal/server` (guardrail) and no
  unfinished FEATURE-0014 work marker may remain.
- This task must not promise a clean working tree before this feature's own
  approved Task 1–21 changes are committed; the changed-file and artifact checks
  apply to the committed FEATURE-0014 changes.
Tests:
- Positive: all Go packages build; `make test` and `make test-race` pass; the
  boundary check passes for `STAGE=tasks`; guardrails pass.
- Negative: an injected out-of-boundary change (for example an added
  connectivity field, an alias kind, or the superseded stack-kind name) causes
  the boundary check or a conformance deny-list test to fail (confirming the
  gate detects drift).
- Boundary: the changed-file set is exactly the FEATURE-0014 approved paths.
- Concurrency: `go test -race ./...` reports no data race.
Acceptance criteria: all listed commands pass for the committed FEATURE-0014
changes, the working tree contains only approved FEATURE-0014 paths plus this
feature's committed changes, no build artifact or generated file is staged, and
the FEATURE-0014 boundary check passes for the tasks stage.
Commit message: `chore(FEATURE-0014): run approved verification, guardrail, and boundary checks`

---

## Requirement-to-task ledger

Every `F14-REQ-01`–`F14-REQ-31` maps to one or more task IDs or an explicit
`NO_TASK` disposition.

| Requirement | Task(s) |
|---|---|
| F14-REQ-01 | Task 1, Task 6, Task 7, Task 8, Task 9, Task 10, Task 19, Task 22 |
| F14-REQ-02 | Task 2, Task 7, Task 19 |
| F14-REQ-03 | Task 5, Task 10, Task 19, Task 22 |
| F14-REQ-04 | Task 1, Task 13, Task 20 |
| F14-REQ-05 | Task 1, Task 6, Task 13, Task 20 |
| F14-REQ-06 | Task 3, Task 4, Task 5, Task 12, Task 13, Task 14, Task 15, Task 16, Task 21 |
| F14-REQ-07 | Task 17, Task 21 |
| F14-REQ-08 | Task 3, Task 4, Task 5, Task 8, Task 10, Task 14, Task 15, Task 16, Task 17, Task 21 |
| F14-REQ-09 | Task 3, Task 4, Task 5, Task 8, Task 9, Task 10, Task 14, Task 15, Task 16, Task 21 |
| F14-REQ-10 | Task 17, Task 21 |
| F14-REQ-11 | Task 20 |
| F14-REQ-12 | Task 18, Task 20 |
| F14-REQ-13 | Task 17, Task 18, Task 20 |
| F14-REQ-14 | Task 18, Task 20 |
| F14-REQ-15 | Task 5, Task 16, Task 20 |
| F14-REQ-16 | Task 5, Task 10, Task 16, Task 20, Task 21 |
| F14-REQ-17 | Task 5, Task 10, Task 16, Task 21 |
| F14-REQ-18 | Task 5, Task 6, Task 7, Task 8, Task 9, Task 10, Task 19, Task 22 |
| F14-REQ-19 | Task 9, Task 15, Task 21 |
| F14-REQ-20 | Task 6, Task 7, Task 8, Task 9, Task 10, Task 15, Task 21 |
| F14-REQ-21 | Task 1, Task 2, Task 3, Task 4, Task 5, Task 6, Task 7, Task 8, Task 9, Task 10, Task 12, Task 13, Task 14, Task 15, Task 16, Task 19, Task 22 |
| F14-REQ-22 | Task 1, Task 2, Task 12, Task 13, Task 14, Task 15, Task 16, Task 19, Task 21 |
| F14-REQ-23 | Task 21, Task 22 |
| F14-REQ-24 | Task 6, Task 10, Task 21, Task 22 |
| F14-REQ-25 | Task 21, Task 22 |
| F14-REQ-26 | Task 18, Task 20, Task 21 |
| F14-REQ-27 | Task 2, Task 7, Task 12, Task 20, Task 21 |
| F14-REQ-28 | Task 6, Task 7, Task 8, Task 9, Task 10, Task 12, Task 13, Task 14, Task 15, Task 16, Task 22 |
| F14-REQ-29 | Task 18, Task 21, Task 22 |
| F14-REQ-30 | Task 17, Task 20, Task 21 |
| F14-REQ-31 | Task 12, Task 13, Task 14, Task 15, Task 16, Task 21 |

No `F14-REQ-*` is unmapped; no `NO_TASK` disposition is required for any
requirement. The `CONTRACT_ONLY / NO_TASK` production mechanisms that underlie
F14-REQ-07, F14-REQ-10, F14-REQ-26, and F14-REQ-29 (live handlers, authoritative
lookup, persistence-backed concurrency, live deletion) are recorded in the
No-task ledger; their required contract outcomes are proven by the mapped helper
and fixture tasks above.

## Decision-to-task ledger

Every `F14-AD-001`–`F14-AD-021`, enumerated individually.

| Decision | Task(s) |
|---|---|
| F14-AD-001 | Task 1, Task 6, Task 7, Task 8, Task 9, Task 10, Task 19, Task 22 |
| F14-AD-002 | Task 1, Task 13, Task 20 |
| F14-AD-003 | Task 1, Task 5, Task 6, Task 13, Task 20 |
| F14-AD-004 | Task 3, Task 4, Task 5, Task 7, Task 8, Task 9, Task 12, Task 13, Task 14, Task 15, Task 16, Task 17, Task 21 |
| F14-AD-005 | Task 3, Task 4, Task 5, Task 8, Task 10, Task 14, Task 15, Task 16, Task 17, Task 21 |
| F14-AD-006 | Task 2, Task 7, Task 19 |
| F14-AD-007 | Task 18, Task 20 |
| F14-AD-008 | Task 3, Task 4, Task 5, Task 8, Task 9, Task 10, Task 14, Task 15, Task 16, Task 21 |
| F14-AD-009 | Task 5, Task 16, Task 20 |
| F14-AD-010 | Task 5, Task 10, Task 16, Task 20 |
| F14-AD-011 | Task 2, Task 5, Task 7, Task 10, Task 12, Task 16, Task 21 |
| F14-AD-012 | Task 9, Task 15, Task 21 |
| F14-AD-013 | Task 6, Task 7, Task 8, Task 9, Task 10, Task 15, Task 21 |
| F14-AD-014 | Task 1, Task 2, Task 3, Task 4, Task 5, Task 6, Task 7, Task 8, Task 9, Task 10, Task 12, Task 13, Task 14, Task 15, Task 16, Task 19, Task 22 |
| F14-AD-015 | Task 1, Task 2, Task 12, Task 13, Task 14, Task 15, Task 16, Task 21 |
| F14-AD-016 | Task 21, Task 22 |
| F14-AD-017 | Task 6, Task 10, Task 21, Task 22 |
| F14-AD-018 | Task 21, Task 22 |
| F14-AD-019 | Task 1, Task 5, Task 6, Task 7, Task 8, Task 9, Task 10, Task 12, Task 13, Task 14, Task 15, Task 16, Task 19, Task 21, Task 22 |
| F14-AD-020 | Task 18, Task 20, Task 21 |
| F14-AD-021 | Task 5, Task 10, Task 19, Task 22 |

## Risk-to-evidence ledger

Every `F14-R01`–`F14-R30`, enumerated individually, mapped to a task test or
review evidence. Target residual ratings are carried unchanged from architecture
section 12 / requirements section 10.2 and are not accepted, downgraded, merged,
renumbered, or closed here.

| Risk | Mitigation evidence (task test / review) |
|---|---|
| F14-R01 | Field deny-list and negative fixtures (Task 21); kind/inventory conformance (Task 19); boundary check (Task 22). |
| F14-R02 | Terminology/inventory conformance (Task 19); single location kind/term schema+type (Task 2, Task 7). |
| F14-R03 | Provider scope validator pos/neg (Task 13); scope-agreement helper (Task 17); cross-scope negatives (Task 21). |
| F14-R04 | Provider type/validator with reused `Organization` owner, no owner kind (Task 1, Task 13); owner-scenario fixtures (Task 20). |
| F14-R05 | Same-party/distinct-party role-separation fixtures (Task 20). |
| F14-R06 | Connectivity non-inference fixtures and schema deny-list (Task 21); failure-domain schema/validator absence (Task 9, Task 15). |
| F14-R07 | Connectivity-field schema deny-list scan (Task 9, Task 21). |
| F14-R08 | Provider-neutrality conformance (Task 19); redaction negatives and native-value deny-list (Task 21); descriptive-only technology (Task 5, Task 16). |
| F14-R09 | Empty-parent-incomplete and completeness-transition fixtures (Task 18, Task 20). |
| F14-R10 | Duplicate-technology/distinct-UID fixtures (Task 5, Task 16, Task 20). |
| F14-R11 | Zero/multiple/wrong-kind/cross-scope failure-domain-parent negatives (Task 16, Task 21); positive single-parent fixture (Task 20). |
| F14-R12 | Skipped-level and wrong-parent negatives (Task 14, Task 15, Task 16, Task 17, Task 21). |
| F14-R13 | Child-existence delete-blocked and leaf-first fixtures; concurrent/stale/retry tests (Task 18, Task 20, Task 21). |
| F14-R14 | Status current-fact-only validator checks; status-history negative scan (Task 12–16, Task 21). |
| F14-R15 | Cross-scope no-existence-disclosure and redaction fixtures (Task 21). |
| F14-R16 | Malformed/prefix-mismatch negatives, valid-unassigned acceptance, and no-authority/no-residency-inference review (Task 12, Task 20, Task 21). |
| F14-R17 | Finite-bounds and technology-length boundary tests (Task 12, Task 16); §6.4 bounds in schemas (Task 6–10). |
| F14-R18 | Stale-version and deterministic-recompute concurrency tests (Task 18, Task 21). |
| F14-R19 | Old-name absence in types/schemas (Task 5, Task 10); inventory conformance and boundary old-name scan (Task 19, Task 22). |
| F14-R20 | Technology-as-descriptive negatives (no enum/compatibility use) (Task 16, Task 21). |
| F14-R21 | Status-history scan and FEATURE-0013 forbidden-concept fixtures (Task 21); boundary applicability check (Task 22). |
| F14-R22 | Interface/dependency/changed-file scan and adjacent-feature deny-list (Task 21, Task 22). |
| F14-R23 | UID non-reuse and deletion-rejection fixtures (Task 18, Task 20, Task 21). |
| F14-R24 | Same-operator/different-owner isolation and cross-owner denial tests (Task 17, Task 20, Task 21). |
| F14-R25 | This tasks.md resolves no semantic choice; delegated-only representation, boundary check (Task 22) and stage review evidence. |
| F14-R26 | Boundary terminology/context scans (Task 22); tasks authored from the approved package only (this stage's context boundary). |
| F14-R27 | One normative home per requirement preserved via the Requirement-to-task and Orphan ledgers; requirements-review evidence. |
| F14-R28 | Exact-ID enumeration in these ledgers and the Orphan report; boundary enumeration check (Task 22). |
| F14-R29 | This stage creates only tasks.md and resolves no earlier-stage artifact; stage-diff review and boundary check (Task 22). |
| F14-R30 | Adjacent-feature and connectivity deny-list fixtures (Task 21); forbidden-concept boundary scan (Task 22). |

## No-task ledger

Invariants, conformance-only controls, inherited behavior, and non-goals that
intentionally produce no implementation task.

- **Inherited FEATURE-0012 grammar (F14-REQ-21).** Type identity, ObjectMeta,
  `ManagedResource` profile, operator-facing boundary, typed references,
  conditions, RFC 9457 Problem Details, RFC 6901 paths, concurrency
  (resourceVersion/ETag/If-Match), unknown/duplicate-field rejection, and
  extensions are reused unchanged from `internal/apimeta`, `internal/apiref`,
  `internal/apicond`, `internal/apiproblem`, `internal/apischema`,
  `internal/apivalid`, and `api/schemas/_common/`. No task modifies them; the
  FEATURE-0012 regression suite runs via Task 22.
- **CONTRACT_ONLY / NO_TASK production mechanisms (design §3.1, C-1…C-8).** Live
  create/get/list/replace/delete handlers; authoritative target/parent
  resolution; authoritative child lookup; production topology-completeness
  recomputation; production deletion execution; stateful/live pagination;
  persistence-backed concurrency and UID history; and no-existence disclosure
  based on authoritative state. FEATURE-0014 owns these as contract and
  conformance expectations only; their required outcomes are proven by the
  offline validators, the Task 17/18 deterministic helpers, and the Task 20/21
  fixtures over explicitly supplied in-memory state. No production state code is
  authored.
- **EXCLUDED mechanisms (design §3.1, X-1…X-7).** Repositories/registries;
  persistence/storage engines; lookup/resolver interfaces and state-provider
  abstractions; adapters and provider connectivity (FEATURE-0016);
  capability/capacity/placement (FEATURE-0015); decision/audit/operation
  (FEATURE-0013, `NOT_APPLICABLE`); provisioning/controllers/reconcilers. None is
  implemented or owned as a contract surface; absence is enforced by the Task 21
  deny-list fixtures and the Task 22 interface/dependency and boundary scans.
- **Non-goals NG-01…NG-10 (requirements §5).** Each is an enforceable exclusion
  realized by the negative and deny-list acceptance criteria of its owning
  requirement (Task 21) and the boundary scan (Task 22); none produces an
  independent implementation task.
- **Governance/process risks F14-R25…F14-R29.** Stage separation, context
  binding, normalization, and traceability are controlled by this stage's
  authoring discipline, these ledgers, and the boundary check (Task 22); they
  produce no feature code task.
- **Terminology migration across docs/RFC/glossary/roadmap/context.** Performed
  as repository-update work before Kiro stage generation (architecture §16.2);
  not a FEATURE-0014 implementation task. Tasks 5, 10, 19, and 22 only ensure the
  active Go/schema contracts use the single `InfrastructureStack` name.

## Orphan report

- Every task above maps to at least one `F14-REQ-*` requirement and at least one
  design section / `DD-*` element and closed `F14-AD-*` decision, as shown in
  each task's labels and in the Requirement-to-task and Decision-to-task ledgers.
- Every design section 3.1 `IMPLEMENT` element maps to a task: I-1 (five Go
  types) → Tasks 1–5; I-2 (five JSON Schemas) → Tasks 6–10; I-3 (offline
  structural validators) → Tasks 12–16; I-4 (offline semantic validators) →
  Tasks 12–16; I-5 (pure deterministic helpers) → Tasks 17–18; I-6 (schema/type-
  binding + conformance fixtures) → Tasks 19–21; I-7 (contract and excluded-
  semantics tests) → tests within Tasks 1–10 and Tasks 12–22. Geographic
  descriptors require syntax and prefix validation only; no dataset task exists.
- Every design `CONTRACT_ONLY / NO_TASK` (C-1…C-8) and `EXCLUDED` (X-1…X-7)
  element maps to an explicit No-task ledger entry, with its contract outcome (
  where applicable) proven by a mapped helper/fixture task.
- No task implements a repository, persistence, storage, resolver, adapter,
  provider integration, capability, capacity, placement, connectivity, decision,
  audit, operation, provisioning, reconciliation, runtime, or FEATURE-0013/0015/
  0016/0053 semantic; no task exists solely to implement a non-goal, governance
  statement, risk prose, or adjacent feature.

NO_ORPHANS
