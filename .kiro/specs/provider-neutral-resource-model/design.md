# FEATURE-0014 Provider-Neutral Resource Model — Design

- Feature: FEATURE-0014 — Provider-Neutral Resource Model
- Stage: Design
- Phase: Phase 2 — Reuse-First PaaS Fabric Foundation
- Kiro slug: provider-neutral-resource-model
- Controlling architecture: `docs/architecture/provider-neutral-resource-model.md`
- Controlling handoff: `ADH-2026-018` (Approved)
- Inherited grammar: `docs/architecture/api-resource-standard.md` (FEATURE-0012), `ADH-2026-012`, `ADH-2026-013`
- Controlling requirements: `.kiro/specs/provider-neutral-resource-model/requirements.md`
- FEATURE-0013 downstream adoption: `NOT_APPLICABLE` (architecture section 9)

## 1. Overview and design boundary

### 1.1 Overview

This design maps the approved FEATURE-0014 requirements (`F14-REQ-01`–`F14-REQ-31`)
onto a concrete representation using the FEATURE-0012 resource grammar unchanged.
It resolves only the representation choices explicitly delegated by architecture
sections 14 and 15. It introduces no new normative requirement, public semantic,
kind, relationship, status meaning, condition meaning, error family, operation, or
interface.

The feature owns exactly five resource kinds — `Provider`, `ProviderLocation`,
`ProviderDatacenter`, `DatacenterFailureDomain`, and `InfrastructureStack` — and
reuses the existing Phase 1 `Organization` resource as the optional supply owner
through the `Provider` primary `scopeRef` (`F14-AD-001`, `F14-AD-002`, F14-REQ-01,
F14-REQ-04). The design binds these resources to the FEATURE-0012
`ManagedResource` profile, `operator-facing` boundary, typed constrained
references, ordered validation, RFC 9457 Problem Details, and optimistic
concurrency (`F14-AD-014`, `F14-AD-015`, F14-REQ-21, F14-REQ-22).

### 1.2 Design boundary

This design is the single authorized output of the Design stage
(architecture section 16.2). It creates or modifies only
`.kiro/specs/provider-neutral-resource-model/design.md`. It does not modify
requirements, architecture, handoffs, automation configuration, source code,
schemas, tests, or tasks. It produces no `tasks.md` content and no
implementation artifact.

Design language is descriptive. The keywords `SHALL` and `MUST` appear only
inside a direct quotation of an approved requirement; every design choice below
cites its originating requirement and closed decision.

Within the closed boundary this design:

- preserves exactly the five owned kinds, `Organization` as the existing optional
  supply owner through `scopeRef`, the exact hierarchy, one immutable immediate
  parent, same-`Provider` scope, registered-shell behavior, the complete
  five-level path, immutable stack identity, one failure-domain parent, and
  leaf-first deletion (`F14-AD-001`–`F14-AD-010`, `F14-AD-020`);
- treats physical containment, shared ancestry, shared owner, and geography as
  implying neither connectivity nor isolation, and keeps missing connectivity
  `Unknown` (`F14-AD-012`, `F14-AD-013`, F14-REQ-19, F14-REQ-20);
- reuses FEATURE-0012 grammar unchanged and selects only representation mechanics
  delegated by architecture sections 14 and 15 (`F14-AD-014`);
- designs no `ResourcePool`, `ProviderCapability`, capacity, compatibility,
  placement, connectivity field/graph, adapter, discovery, credential, endpoint,
  provider call, provider SDK type, repository, persistence engine, plugin,
  provisioning, runtime execution, or speculative extension point (`F14-AD-016`–
  `F14-AD-019`, F14-REQ-23, F14-REQ-24, F14-REQ-25).

## 2. Resolved delegated design decisions

The decisions below resolve only the bounded representation choices delegated by
architecture sections 14–15 and enumerated as `DQ-01`–`DQ-08` in requirements
section 9. No decision changes a resource's meaning, ownership, cardinality,
scope, hierarchy, identity, connectivity posture, deletion semantics, or feature
boundary (architecture section 15).

### DD-01 — API group and versioned routes (resolves DQ-01; F14-REQ-21)

- **Decision.** The five kinds share one domain-grouped alpha API group
  `fabric.sovrunn.io`, version `v1alpha1`, so `apiVersion` is
  `fabric.sovrunn.io/v1alpha1`. HTTP collections use the FEATURE-0012
  domain-grouped versioned route form `/apis/fabric.sovrunn.io/v1alpha1/<collection>`
  and `/apis/fabric.sovrunn.io/v1alpha1/<collection>/{name}`, with a separately
  authorized status path `/{name}/status`. No unversioned public endpoint is
  introduced.
- **Rationale.** FEATURE-0012 section 6.1 requires new Phase 2 APIs to use
  domain-grouped versioned routes and singular PascalCase kinds, lowerCamelCase
  fields, and plural kebab-case collections. `fabric.sovrunn.io` is the
  provider-supply/infrastructure domain group already used by the FEATURE-0012
  reference example (`resourcePoolRef: apiVersion: fabric.sovrunn.io/v1alpha1`),
  which keeps FEATURE-0014 topology references stable for later features without
  implying any FEATURE-0015 semantics. Alpha maturity matches the stage-ownership
  matrix ("one domain-grouped alpha API").
- **Boundary.** This is representation only; it adds no public compatibility
  semantic beyond FEATURE-0012 (architecture section 14, "Exact API group and
  routes").

### DD-02 — Minimal field inventory (resolves DQ-02; F14-REQ-01, F14-REQ-05–F14-REQ-17, F14-REQ-27)

- **Decision.** Each kind carries only the fields necessary to express the closed
  decisions. Identity, human label, and governance scope live in inherited
  FEATURE-0012 `metadata`; desired topology lives in a minimal `spec`; observed
  facts live in system-owned `status`. The full inventory is section 4.
- **Rationale.** Architecture section 14 ("Initial domain field inventory")
  restricts fields to those necessary to express locked semantics and forbids
  adding domain meaning through a convenient field. `Provider.spec` therefore
  carries no domain field (identity and scope are metadata); descendants carry
  only their required immutable typed parent reference plus, where delegated, a
  bounded descriptive attribute.
- **Boundary.** No capacity, capability, eligibility, connectivity, adapter, or
  provider-native field is added (`F14-AD-017`–`F14-AD-019`, F14-REQ-18,
  F14-REQ-24, F14-REQ-25).

### DD-03 — Condition names for registered validity and topology completeness (resolves DQ-03; F14-REQ-12, F14-REQ-14, F14-REQ-22)

- **Decision.** Two FEATURE-0012 condition types express current facts:
  - `Valid` — the resource passed all applicable structural, semantic, reference,
    and scope validation (registered validity).
  - `TopologyComplete` — an unbroken required child path exists beneath this
    resource down to at least one `InfrastructureStack` (topology completeness).
  Both use tri-state `status` (`True`/`False`/`Unknown`), stable PascalCase
  `reason` codes, `observedGeneration`, and `lastTransitionTime`, exactly as
  FEATURE-0012 section 6.9 defines.
- **Rationale.** Architecture delegates "exact condition names representing
  registered validity and topology completeness, using FEATURE-0012 condition
  grammar." `Valid` reuses the FEATURE-0012 example condition; `TopologyComplete`
  names the registered-vs-complete distinction of `F14-AD-007` without adding a
  new grammar. On an `InfrastructureStack` (the leaf level), `TopologyComplete`
  is `True` when the stack itself is `Valid`, because the required path
  terminates at the stack.
- **Boundary.** `TopologyComplete` is an administrative topology fact only; it
  carries no capability, capacity, placement-eligibility, or service-readiness
  meaning (F14-REQ-14). Conditions are system-owned current facts and never
  become history, capability, capacity, or audit (F14-REQ-22).

### DD-04 — Normalized geographic-code representation (resolves DQ-04; F14-REQ-27)

- **Decision.** Where a `ProviderLocation` carries a geographic/sovereignty
  descriptor, the design represents it as an optional `spec.geo` value with a
  mature normalized `countryCode` (ISO 3166-1 alpha-2, exactly two uppercase
  letters) and an optional `subdivisionCode` (ISO 3166-2, bounded). The value is
  descriptive declared topology only.
- **Rationale.** Architecture delegates "mature normalized geographic-code
  representation where applicable." ISO 3166-1/3166-2 are the mature standards
  named in the reuse assessment (architecture section 10).
- **Boundary.** Geography is not residency proof, compliance evidence, or a policy
  outcome, and drives no placement (F14-REQ-27, `F14-AD-011`). The field is
  optional; its absence carries no inference.

### DD-05 — Bounded descriptive infrastructure-technology representation (resolves DQ-05; F14-REQ-17)

- **Decision.** `InfrastructureStack.spec.technology` is a single bounded
  descriptive string (bounded length and character set, section 6). It is not
  validated against a closed core enum, is not identity, and is not a lookup or
  reference key.
- **Rationale.** Architecture delegates "bounded descriptive
  infrastructure-technology representation that cannot become identity,
  capability, or compatibility logic." The architecture's illustrative
  technologies (Apache CloudStack, Cloud Foundry, Red Hat OpenShift, AWS IaaS,
  Azure IaaS) are explicitly non-normative examples, not enum values.
- **Boundary.** Technology never becomes identity, a closed compatibility enum,
  or capability/placement logic (F14-REQ-17, `F14-AD-011`); stack identity is the
  immutable `metadata.uid` plus scoped identity (F14-REQ-15).

### DD-06 — Finite initial limits and configuration mechanism (resolves DQ-06; F14-REQ-28)

- **Decision.** All collections, strings, and payloads are finite, and this
  design fixes an exact initial value for every one of them in section 6.4
  (fully resolving DQ-06). This includes whole-object size, JSON nesting depth,
  `metadata.name`/`displayName` length, label entry count, label key and value
  lengths, annotation total size, `status.conditions` count, typed references
  per field, Problem Details violation count, geographic-code lengths,
  `technology` length, and list page size. No FEATURE-0014 bound is deferred to
  another feature to supply. The values are loaded through a validated
  configuration struct at startup and are treated as reviewed platform
  configuration, not user-authored input.
- **Rationale.** Architecture section 15 delegates "initial finite limits and
  configuration values" to design, so the exact numbers are chosen here rather
  than inherited implicitly. The chosen values equal the reviewed platform limit
  set already applied by the shared FEATURE-0012 validation primitives
  (`internal/apivalid` `DefaultLimits`), keeping FEATURE-0014 consistent with the
  rest of the platform without redefining the grammar. The Go configuration
  package pattern (`docs/engineering/go-coding-guardrails.md` section 25)
  supplies the loading mechanism.
- **Boundary.** No value is left unbounded, no bound is stated only as
  "inherited", and no limit silently changes a semantic (F14-REQ-28,
  `F14-AD-014`, `F14-AD-015`).

### DD-07 — Schema composition, indexes, and package layout (resolves DQ-07; F14-REQ-21)

- **Decision.** One canonical JSON Schema (2020-12) per kind composes the shared
  FEATURE-0012 type-metadata, `ObjectMeta`, reference, and status/condition
  fragments (under `api/schemas/_common/`) and adds only the FEATURE-0014 `spec`
  fields. Schema files and Go packages follow the live repository conventions in
  section 3. Logical query behaviors (section 4.4) support unique identity, child
  lookup for deletion and completeness, and paginated reads. This design
  specifies those required query behaviors but names no storage component and
  defines no registry, repository, persistence engine, or adapter interface.
- **Rationale.** Architecture delegates "schema composition, indexes, package
  layout, and internal implementation structure." Composition over the shared
  `api/schemas/_common/` fragments enforces grammar reuse without forking it.
  Schema filenames (`api/schemas/*.json`) and Go packages (`internal/resources`,
  `internal/validation`, `internal/api`) follow the live conventions used by the
  existing Phase 1 resources and the Go guardrails (short meaningful names).
- **Boundary.** No storage component, registry, repository, persistence engine,
  or adapter interface is designed or assumed to pre-exist; only the required
  query behaviors are stated (closed design boundary).

### DD-08 — Immediate-parent representation and `ownerRef` mirroring (resolves DQ-08; F14-REQ-06, F14-REQ-08, F14-REQ-09, F14-REQ-10)

- **Decision.** Each descendant's immediate topology parent is represented solely
  by one typed, constrained, immutable domain reference in `spec`
  (`providerLocationRef`, `providerDatacenterRef`, `datacenterFailureDomainRef`);
  `ProviderLocation` has no separate parent field because its parent `Provider` is
  already identified by `metadata.scopeRef`. FEATURE-0012 `ownerRef` is **not**
  populated to encode the topology parent, so no second topology authority exists.
  Governance scope remains exclusively `metadata.scopeRef`.
- **Rationale.** Architecture section 5 delegates "whether immediate-parent
  containment is represented by a domain-specific parent reference alone or
  additionally mirrored by FEATURE-0012 `ownerRef`, provided no second topology
  authority is created." Choosing the domain reference alone gives one
  unambiguous topology authority and one governance authority (`scopeRef`),
  avoiding the dual-authority failure condition of the single-owner overlap check
  (architecture section 2.4). Deletion protection and completeness are computed
  from this single domain reference.
- **Boundary.** A parent reference never replaces `scopeRef`, grants
  authorization, or changes isolation (F14-REQ-05, F14-REQ-06); `ownerRef`
  retains its inherited FEATURE-0012 lifecycle-containment meaning and is not
  repurposed (`F14-AD-003`).

### 2.1 Unresolved semantic-gap check

Every approved requirement is representable using only the DD-01–DD-08 delegated
choices. No requirement required a semantic choice outside the section 13 closed
register or the section 14–15 delegation matrix. Result:
`NO_UNRESOLVED_SEMANTIC_GAP`. No `ARCHITECTURE_DECISION_REQUIRED` condition was
encountered while authoring this design.

## 3. Component and package architecture

The design adds FEATURE-0014-only components. It creates no adjacent-feature
scaffolding (no pool, capability, adapter, discovery, connectivity, decision, or
audit package).

```text
cmd/sovrunn-api/          # entrypoint (inherited; no FEATURE-0014 logic)
internal/resources/       # five resource structs, one file per kind: provider.go, providerlocation.go,
                          #   providerdatacenter.go, datacenterfailuredomain.go, infrastructurestack.go
internal/validation/      # per-kind offline structural/semantic validators + stateful reference/scope/
                          #   deletion validators, one file per kind (provider.go, providerlocation.go, …)
internal/api/             # per-kind HTTP handlers + decoders binding the DD-01 routes, one pair per kind
                          #   (provider_handler.go + provider_decode.go, …)
api/schemas/              # one canonical JSON Schema per kind: provider.json, provider-location.json,
                          #   provider-datacenter.json, datacenter-failure-domain.json, infrastructure-stack.json
                          #   (composing api/schemas/_common/ fragments; described here, authored in tasks/impl)
tests/conformance/        # conformance fixtures under tests/conformance/fixtures/, checks in internal/apiconform
                          #   (described in section 9; authored in tasks/impl)
```

The layout follows the live flat per-kind file convention used by the existing
Phase 1 resources (for example `internal/resources/organization.go`,
`internal/validation/organization.go`, `internal/api/org_handler.go`) and the
existing flat schema directory (`api/schemas/*.json` with shared fragments in
`api/schemas/_common/`). No per-feature subpackage is introduced, so no path
deviates from the live repository conventions.

Inherited, not designed here:

- FEATURE-0012 shared primitives are reused unchanged (F14-REQ-21): type
  metadata (`internal/apimeta`), typed constrained references (`internal/apiref`),
  conditions (`internal/apicond`), Problem Details (`internal/apiproblem`),
  schema/route support (`internal/apischema`), strict decoding and ordered
  validation (`internal/apivalid`), and conformance checks
  (`internal/apiconform`), plus the shared `api/schemas/_common/` fragments.
- Logical lookup behavior: the schemas, validators, handlers, and conformance
  fixtures express the identity, parent-child, scope, and pagination behaviors
  listed in section 4.4. FEATURE-0014 defines no repository, persistence
  mechanism, storage component, or production-persistence task; those remain
  outside this feature boundary.

Component responsibilities:

- `internal/resources` — declares the five typed structs (one file per kind)
  with explicit JSON tags, embedding the shared `apimeta` metadata/status types
  and the `apicond` condition type; declares the three typed constrained
  parent-reference aliases over `apiref` (DD-08).
- `internal/validation` — implements the ordered validation pipeline of
  section 7.1 (one validator file per kind), split into offline
  (structural/semantic) validators callable without external state and stateful
  (reference/scope/authorization/concurrency/deletion) validators that read
  existing resource state (architecture section 16.4).
- `internal/api` — thin, context-aware handlers and decoders (one pair per kind)
  that decode safely, invoke validators, perform the section 4.4 logical lookups,
  recompute `Valid`/`TopologyComplete`, and emit inherited Problem Details;
  handlers hold no business logic beyond wiring (go-coding-guardrails sections 8,
  34).

## 4. Resource and field model

All five kinds use profile `ManagedResource`, boundary `operator-facing`,
stability `alpha`, and one canonical schema each (`F14-AD-014`, `F14-AD-015`,
F14-REQ-21, F14-REQ-22). Inherited `metadata`, `status`, and condition fields are
reused unchanged and are not redefined here.

### 4.1 Common inherited shape (reused, not redefined)

- `apiVersion`: `fabric.sovrunn.io/v1alpha1`; `kind`: the singular PascalCase kind.
- `metadata`: `name` (immutable, DNS-style), `uid` (opaque, server-generated,
  never reused), `displayName` (mutable human label), `scopeRef`, `labels`,
  `annotations`, `generation`, `resourceVersion`, `createdAt`, `updatedAt`.
- `status`: `observedGeneration`, `conditions[]` with `Valid` and
  `TopologyComplete` (DD-03). No `phase` or other lifecycle-vocabulary field is
  introduced, and no lifecycle vocabulary or phase-to-condition consistency rule
  is authorized (F14-REQ-22).

### 4.2 Per-kind spec inventory (minimal; DD-02)

| Kind | `metadata.scopeRef` | Required `spec` domain fields | Optional `spec` domain fields |
|---|---|---|---|
| `Provider` | `Platform` or one `Organization` | none | none |
| `ProviderLocation` | the containing `Provider` | none (parent is `scopeRef`) | `geo` = `{ countryCode, subdivisionCode? }` (DD-04) |
| `ProviderDatacenter` | the containing `Provider` | `providerLocationRef` (immutable, kind `ProviderLocation`) | none |
| `DatacenterFailureDomain` | the containing `Provider` | `providerDatacenterRef` (immutable, kind `ProviderDatacenter`) | none |
| `InfrastructureStack` | the containing `Provider` | `datacenterFailureDomainRef` (immutable, kind `DatacenterFailureDomain`) | `technology` (bounded descriptive string, DD-05) |

Typed reference shape (inherited FEATURE-0012 section 6.6): `{ apiVersion, kind,
name, uid? }`, constrained to the required parent kind and to `Provider` scope.

### 4.3 Authoritative writer and mutability ledger

Every field and condition has exactly one authoritative writer (architecture
section 16.4; FEATURE-0012 section 6.8). No field lacks a writer.

| Field / condition | Authoritative writer | Mutable | Rule / citation |
|---|---|---|---|
| `apiVersion`, `kind` | Sovrunn (schema) | No | Immutable per object version (F14-REQ-21). |
| `metadata.name` | Authorized creator | No | Stable URL-safe identity within scope (F14-REQ-21). |
| `metadata.uid` | Sovrunn | No | Opaque, globally unique, never reused (F14-REQ-15, F14-REQ-26). |
| `metadata.displayName` | Authorized operator | Yes | Human label, not identity (F14-REQ-21). |
| `metadata.scopeRef` | Creator on create; validated by Sovrunn | No | Sole governance authority; immutable (F14-REQ-05, F14-REQ-06). |
| `metadata.labels` / `metadata.annotations` | Authorized operator / namespaced owner | Yes | Bounded; no secrets or native values (F14-REQ-28, F14-REQ-31). |
| `metadata.generation` / `resourceVersion` / timestamps | Sovrunn | System-only | System-owned (F14-REQ-21, F14-REQ-29). |
| `spec.geo` (`ProviderLocation`) | Authorized operator | Yes | Descriptive normalized code; not compliance (F14-REQ-27). |
| `spec.providerLocationRef` / `providerDatacenterRef` / `datacenterFailureDomainRef` | Authorized operator on create | No | One immutable immediate parent (F14-REQ-08, F14-REQ-09). |
| `spec.technology` (`InfrastructureStack`) | Authorized operator | Yes | Bounded descriptive only; not identity/capability (F14-REQ-17). |
| `status.observedGeneration` | registered system-owned status producer | System-only | Current facts (F14-REQ-22). |
| condition `Valid` | one registered producer | System-only | Registered validity (F14-REQ-12, F14-REQ-22). |
| condition `TopologyComplete` | one registered producer | System-only | Completeness; recomputed deterministically (F14-REQ-14, F14-REQ-29). |

No connectivity, capacity, capability, adapter, credential, endpoint, or
provider-native field exists in any inventory above (F14-REQ-18, F14-REQ-20,
F14-REQ-24, F14-REQ-25, F14-REQ-31).

### 4.4 Logical indexes (DD-07)

- Unique identity index on `(apiVersion group, kind, scope UID, name)` — enforces
  FEATURE-0012 identity uniqueness (F14-REQ-21).
- Parent index on the immediate-parent reference UID — supports child lookup for
  deletion rejection (F14-REQ-26) and `TopologyComplete` recomputation
  (F14-REQ-14, F14-REQ-29).
- Scope index on `Provider` scope UID — supports same-scope resolution
  (F14-REQ-10) and scoped, paginated list reads (F14-REQ-28).

These are logical lookup and conformance behaviors required by the resource
contract; FEATURE-0014 defines no storage or persistence implementation.

## 5. API routes and FEATURE-0012 binding

Routes follow DD-01. Each kind exposes the FEATURE-0012 initial normative
operations (create, get, list, full replace, delete) plus a separately authorized
status path (FEATURE-0012 sections 6.12, 6.13).

| Collection | Kind | List envelope |
|---|---|---|
| `/apis/fabric.sovrunn.io/v1alpha1/providers` | `Provider` | `ProviderList` |
| `/apis/fabric.sovrunn.io/v1alpha1/provider-locations` | `ProviderLocation` | `ProviderLocationList` |
| `/apis/fabric.sovrunn.io/v1alpha1/provider-datacenters` | `ProviderDatacenter` | `ProviderDatacenterList` |
| `/apis/fabric.sovrunn.io/v1alpha1/datacenter-failure-domains` | `DatacenterFailureDomain` | `DatacenterFailureDomainList` |
| `/apis/fabric.sovrunn.io/v1alpha1/infrastructure-stacks` | `InfrastructureStack` | `InfrastructureStackList` |

Binding rules (all inherited from FEATURE-0012, F14-REQ-21):

- `POST` create returns 201; `GET` get/list returns 200; `PUT` full replace
  returns 200 and requires `If-Match` with the opaque `resourceVersion`;
  `DELETE` returns 204 on success.
- `status` updates use the separately authorized status path; normal clients
  cannot write `status` (F14-REQ-22).
- Lists use `ListEnvelope` with opaque page tokens, deterministic ordering,
  bounded page size, and filtering limited to indexed fields (F14-REQ-28).
- User-authored `status`, `resourceVersion`, `generation`, and timestamps are
  rejected on create/replace (F14-REQ-22; go-coding-guardrails section 10).

## 6. Field representation, limits, and configuration

### 6.1 Names and identity

`metadata.name` follows the inherited DNS-style pattern
`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`; `metadata.uid` is opaque and never reused
(F14-REQ-15, F14-REQ-26, F14-REQ-21).

### 6.2 Geographic code (`ProviderLocation.spec.geo`; DD-04)

- `countryCode`: ISO 3166-1 alpha-2, exactly two uppercase ASCII letters
  (`^[A-Z]{2}$`).
- `subdivisionCode`: optional ISO 3166-2 code of the form
  `<countryCode>-<subdivision>`, where `<subdivision>` is 1–3 uppercase ASCII
  letters or digits (`^[A-Z]{2}-[A-Z0-9]{1,3}$`; total length ≤ 6). When present,
  its two-letter country prefix MUST equal `countryCode`.
- **Authoritative dataset.** Validation accepts only codes present in a pinned,
  embedded snapshot of the ISO 3166-1 alpha-2 country list and the ISO 3166-2
  subdivision list. The snapshot is a build-time data asset committed to the
  repository; its version is recorded as the ISO 3166 publication date the
  snapshot was taken from, stored in a manifest beside the data file. Validation
  performs no network lookup and no runtime fetch — it is a pure membership test
  against the embedded snapshot.
- **Update mechanism.** Refreshing the dataset is an explicit, reviewed change
  that replaces the embedded snapshot and bumps its recorded version identifier.
  "Unknown" codes are exactly those absent from the currently embedded snapshot.
  A dataset refresh changes no resource meaning and is never user-authored input.
- Descriptive only; a value that fails the format check or is absent from the
  embedded snapshot is rejected as malformed/unknown (F14-REQ-27).

### 6.3 Technology descriptor (`InfrastructureStack.spec.technology`; DD-05)

- **Deterministic validation.** The accepted value matches
  `^[A-Za-z0-9][A-Za-z0-9 ._+-]{0,99}$` — a single line of 1–100 characters
  drawn only from ASCII letters, digits, space, and the punctuation set
  `. _ + -`, and starting with a letter or digit. No control characters, no
  non-ASCII/Unicode characters, and no leading or trailing space (after
  normalization) are permitted.
- **Deterministic normalization**, applied before validation and storage in this
  fixed order: (1) reject the value if it contains any character outside the
  ASCII set above; (2) trim leading and trailing ASCII spaces; (3) collapse each
  internal run of spaces to a single space. Case is preserved (not folded),
  because the descriptor is not identity and case-folding would imply an equality
  contract.
- The normalized value is stored verbatim. It is not matched against a closed
  enum, not case-folded for comparison, and not used as identity or a
  reference/lookup key (F14-REQ-17, F14-REQ-18, F14-REQ-15).

### 6.4 Reviewed initial finite limits (DD-06; F14-REQ-28)

| Limit | Reviewed initial value | Citation |
|---|---|---|
| whole-object size | ≤ 1,048,576 bytes (1 MiB) | F14-REQ-28 |
| JSON nesting depth | ≤ 32 levels | F14-REQ-28 |
| `metadata.name` length | ≤ 253 characters | F14-REQ-21, F14-REQ-28 |
| `metadata.displayName` length | ≤ 253 characters | F14-REQ-28 |
| `metadata.labels` entries | ≤ 64 entries | F14-REQ-28 |
| label key length | ≤ 63 characters | F14-REQ-28 |
| label value length | ≤ 253 characters | F14-REQ-28 |
| `metadata.annotations` total size | ≤ 262,144 bytes (256 KiB) | F14-REQ-28, F14-REQ-31 |
| `status.conditions` entries | ≤ 32 conditions | F14-REQ-28 |
| typed references per field | ≤ 64 references | F14-REQ-28 |
| Problem Details violations per response | ≤ 100 violations | F14-REQ-21, F14-REQ-28 |
| `geo.countryCode` length | exactly 2 characters | F14-REQ-27 |
| `geo.subdivisionCode` length | ≤ 6 characters | F14-REQ-27 |
| `technology` length | ≤ 100 characters | F14-REQ-17 |
| list page size | default 50, maximum 200 | F14-REQ-28 |

Every bound above is an exact finite value fixed by this design (resolving
DQ-06); no FEATURE-0014 bound is left to another feature to supply. The values
are the reviewed platform limit set and are loaded through a validated
configuration struct at startup, treated as reviewed platform configuration and
not user-authored input (go-coding-guardrails section 25).

## 7. Validation, error behavior, and concurrency

### 7.1 Ordered validation pipeline

The pipeline mirrors FEATURE-0012 section 6.10, split into offline and stateful
stages (architecture section 16.4). Structural and semantic checks are safe for
offline schema validation; reference, authorization, concurrency, and deletion
checks require access to existing resource state.

| Order | Stage | Offline / stateful | FEATURE-0014 checks | Inherited stable error behavior |
|---|---|---|---|---|
| 1 | Structural | Offline | required fields, DNS name, unknown/duplicate-field rejection, bounded sizes (section 6.4), status not user-authored | 400 malformed / 422 `VALIDATION_FAILED` with RFC 6901 path (F14-REQ-21, F14-REQ-22) |
| 2 | Semantic | Offline | `scopeRef.kind` allowed per kind; required parent-reference kind correct; `geo` normalized; `technology` bounded/descriptive; hierarchy has no skipped/reordered level | 422 `VALIDATION_FAILED` (F14-REQ-05, F14-REQ-08, F14-REQ-17, F14-REQ-27) |
| 3 | Reference | Stateful | parent reference resolves to an existing resource of the required kind; parent and child resolve to the same `Provider` scope UID; cross-provider rejected | 422/404 with no-existence disclosure (F14-REQ-07, F14-REQ-09, F14-REQ-10) |
| 4 | Authorization | Stateful | caller authorized for the resolved `scopeRef`; multi-owner isolation; cross-scope denied without disclosure | 403 / no-existence disclosure (F14-REQ-05, F14-REQ-30) |
| 5 | Concurrency | Stateful | `If-Match` / `resourceVersion` checked on replace and delete; immutable `scopeRef` and parent reference enforced against stored version | 412 stale version; 422 on immutable-field mutation (F14-REQ-06, F14-REQ-09, F14-REQ-29) |
| 6 | Deletion | Stateful | reject deletion while children exist (via parent index); no cascade; no silent reparenting; leaf-first | 409 delete-blocked conflict (F14-REQ-26) |

### 7.2 Inherited error behavior (no new error family)

All failures use the inherited FEATURE-0012 RFC 9457 Problem Details contract with
stable codes and RFC 6901 JSON Pointer paths (FEATURE-0012 section 6.11). No new
problem type, code, violation code, or error family is introduced (F14-REQ-21).
Cross-scope and cross-provider denials return the same response as a genuinely
absent target, disclosing no existence (F14-REQ-07, F14-REQ-30). Responses,
errors, and conditions carry no provider-native identifiers, credentials,
endpoints, or secrets (F14-REQ-31).

### 7.3 Concurrency and deterministic recomputation

Immutable parents plus `resourceVersion`/ETag/`If-Match` provide conflict-safe
onboarding and decommissioning; `TopologyComplete` is recomputed deterministically
from the parent index after concurrent create/delete so that completeness facts
converge (F14-REQ-14, F14-REQ-29).

## 8. Current-fact condition behavior

`status` holds current system-owned facts only and never becomes history,
capability, capacity, or audit (F14-REQ-22, `F14-AD-015`).

- `Valid` reflects the outcome of the latest validation for the observed
  `generation`; `reason` is a stable PascalCase code; `lastTransitionTime` changes
  only on status change (F14-REQ-12).
- `TopologyComplete` reflects, at the observed `generation`, whether an unbroken
  required child path to at least one `InfrastructureStack` currently exists
  beneath the resource; on an `InfrastructureStack` it is `True` when the stack is
  `Valid` (DD-03, F14-REQ-14). It carries no capability, capacity, placement, or
  readiness meaning.
- Absence of a connectivity assertion is represented as no field and no condition;
  connectivity is `Unknown` by absence and is never inferred from ancestry,
  proximity, sibling order, or shared owner `Organization` (F14-REQ-19,
  F14-REQ-20). No decision, audit, or operation semantics are recorded
  (F14-REQ-23).

## 9. Security, redaction, no-existence disclosure, isolation, and conformance

### 9.1 Security and privacy realization

- Sole governance authority is `metadata.scopeRef`; `ownerRef` and containment
  grant no authorization (F14-REQ-05, F14-REQ-06; DD-08).
- `operator-facing` boundary with least-privilege reads; customer-facing or AI
  projection is a locked non-goal requiring separate review (F14-REQ-22).
- No-existence disclosure on cross-scope list/get/reference (F14-REQ-07,
  F14-REQ-30).
- Redaction of provider-native identifiers, credentials, endpoints, and secrets
  from metadata, errors, and conditions (F14-REQ-31).
- Multi-owner isolation: the same external operator under two owner Organizations
  yields independent scoped identities; authorization never keys on external
  operator name or native ID (F14-REQ-30).

### 9.2 Conformance fixtures and verification mechanisms

Fixtures (described here; authored during implementation) realize the architecture
section 11 fitness checks and the risk evidence obligations:

- exactly-five-kinds and zero-alias schema/kind inventory (F14-REQ-01, F14-REQ-03).
- distinct-owner/operator case (`Organization: OwnerOrganization-A` →
  `Provider: ProviderOperator-A`/`ProviderOperator-B`) and same-party case
  (`Organization: ProviderOperator-A` → `Provider: ProviderOperator-A`)
  (F14-REQ-11).
- empty registered parents at each level accepted as incomplete; every
  skipped-level relationship rejected (F14-REQ-08, F14-REQ-12, F14-REQ-13,
  F14-REQ-14).
- two same-technology stacks retain distinct UIDs and parents; stack with zero,
  ambiguous, multiple, or cross-scope failure-domain parents rejected
  (F14-REQ-15, F14-REQ-16).
- cross-provider / cross-scope reference rejected without existence disclosure;
  parent/child scope-UID mismatch rejected (F14-REQ-07, F14-REQ-10).
- same-datacenter, cross-datacenter, cross-location, and cross-provider
  relationships (including two Providers under one owner Organization) carry no
  connectivity semantics; schema deny-list finds no connectivity field
  (F14-REQ-19, F14-REQ-20).
- invalid/unknown geographic code rejected; valid code carries no residency
  inference (F14-REQ-27).
- parent deletion with children returns a child-existence conflict; leaf-first
  deletion succeeds; UID non-reuse (F14-REQ-26).
- stale-version write rejected; deterministic completeness recomputation after
  concurrent change (F14-REQ-29).
- field/interface deny-list finds no FEATURE-0015/0016 semantics and no
  FEATURE-0013 governed concept outside non-goal/ownership text (F14-REQ-23,
  F14-REQ-24, F14-REQ-25).

Verification mechanisms: inherited FEATURE-0012 conformance/regression suites
(unchanged), the FEATURE-0014 boundary validator (`scripts/feature-0014-boundary-check.py`),
schema/route inventory and terminology scans, and human semantic review. This
design authors none of these; it specifies what they verify.

## 10. Requirement traceability (`F14-REQ-01`–`F14-REQ-31`)

Exactly one row per requirement. Each row maps the requirement to its
representation, validator, error/evidence behavior, and confirms no new normative
behavior is introduced (architecture section 16.4).

| Requirement | Representation | Validator (section 7.1 stage) | Error / evidence behavior | New norm? |
|---|---|---|---|---|
| F14-REQ-01 | DD-01/DD-02 five kinds, one schema each | structural (schema/kind inventory) | boundary check; inventory fixture | No |
| F14-REQ-02 | `ProviderLocation` sole location kind/term | structural + semantic | terminology scan; route/schema inventory | No |
| F14-REQ-03 | `InfrastructureStack` sole active name (DD-01) | structural | old-name scan; schema/API diff | No |
| F14-REQ-04 | `Provider.scopeRef` → existing `Organization`; no owner kind | semantic (scope kind) | kind inventory; owner-scenario fixture | No |
| F14-REQ-05 | `Provider.scopeRef` = Platform/Organization; `ownerRef` not scope | semantic + authorization | pos/neg scope fixtures; cross-scope authz tests | No |
| F14-REQ-06 | Immutable `scopeRef`=Provider on descendants (4.2/4.3) | semantic + concurrency | immutable-field rejection test | No |
| F14-REQ-07 | Typed constrained ref; no-existence disclosure | reference | cross-provider denial fixture (no disclosure) | No |
| F14-REQ-08 | Exact hierarchy via required parent-ref kind (4.2) | semantic | skipped/reordered-level negative fixtures | No |
| F14-REQ-09 | One immutable typed parent ref (DD-08) | semantic + concurrency | zero/multiple/wrong-kind + mutation tests | No |
| F14-REQ-10 | Parent/child same Provider scope UID (4.4 scope index) | reference | scope-UID mismatch fixture | No |
| F14-REQ-11 | Role-specific `Organization` and `Provider` resources | semantic | same-party/distinct-party fixtures | No |
| F14-REQ-12 | `Valid` condition; registered ≠ complete (DD-03) | status recomputation | empty-parent accepted-incomplete fixture | No |
| F14-REQ-13 | Zero-child registration allowed; skips rejected | semantic | register-before-children + skip-reject fixtures | No |
| F14-REQ-14 | `TopologyComplete` condition over five-level path (DD-03) | status recomputation | full-path complete / missing-level incomplete fixtures | No |
| F14-REQ-15 | Immutable `metadata.uid` + scoped identity (4.3) | structural + reference | duplicate-technology/distinct-UID fixture | No |
| F14-REQ-16 | `datacenterFailureDomainRef` one immutable parent (4.2) | semantic + reference | zero/multiple/cross-scope parent fixtures | No |
| F14-REQ-17 | `spec.technology` bounded descriptive (DD-05, 6.3) | structural + semantic | enum/decision-use scan; negative compatibility test | No |
| F14-REQ-18 | No provider-native identity/reference fields (4.3) | structural | provider-neutral schema scan | No |
| F14-REQ-19 | No connectivity field; no inference (4.3, 8) | structural | same/cross-parent non-inference fixtures | No |
| F14-REQ-20 | Connectivity `Unknown` by absence; deny-list (4.3, 8) | structural | schema deny-list scan | No |
| F14-REQ-21 | Inherited FEATURE-0012 grammar unchanged (all sections) | all stages | FEATURE-0012 conformance/regression suite | No |
| F14-REQ-22 | `ManagedResource`/operator-facing; status current-only (DD-03, 8) | structural + status | status-history scan; writer-ownership review | No |
| F14-REQ-23 | No FEATURE-0013 machinery; `NOT_APPLICABLE` (1.2, 8) | boundary scan | forbidden-concept scan; applicability check | No |
| F14-REQ-24 | No ResourcePool/capability/capacity fields (4.3) | structural | field deny-list scan | No |
| F14-REQ-25 | No adapter/integration interfaces or deps (3) | boundary scan | interface/dependency/SDK scan | No |
| F14-REQ-26 | Parent index deletion rejection; no cascade (7.1 stage 6) | deletion | child-existence conflict; UID non-reuse fixtures | No |
| F14-REQ-27 | `geo` normalized code, descriptive only (DD-04, 6.2) | structural + semantic | invalid/unknown-code + no-residency-inference fixtures | No |
| F14-REQ-28 | Finite limits + pagination (DD-06, 6.4) | structural | boundary/pagination tests | No |
| F14-REQ-29 | ETag/If-Match; deterministic recompute (7.3) | concurrency | stale-write + recompute tests | No |
| F14-REQ-30 | Independent scoped identity; no name/native authz (9.1) | authorization | same-operator/different-owner isolation tests | No |
| F14-REQ-31 | Redaction of native values/secrets (7.2, 9.1) | all stages (output) | secret/redaction response tests | No |

## 11. Architecture-decision traceability (`F14-AD-001`–`F14-AD-021`)

Each closed decision is enumerated individually (no range shorthand) with its
design disposition.

| Decision | Design disposition and evidence |
|---|---|
| `F14-AD-001` | DD-01/DD-02 define exactly five kinds; no alias (F14-REQ-01). |
| `F14-AD-002` | `Provider.scopeRef` reuses existing `Organization`; no owner kind (F14-REQ-04, F14-REQ-11). |
| `F14-AD-003` | `scopeRef` = Platform/Organization; `ownerRef` not scope (DD-08, F14-REQ-05). |
| `F14-AD-004` | Descendants immutably Provider-scoped; cross-provider rejected without disclosure (F14-REQ-06, F14-REQ-07, F14-REQ-10). |
| `F14-AD-005` | Required parent-ref kind enforces the only hierarchy; no skips (F14-REQ-08). |
| `F14-AD-006` | `ProviderLocation` is the only location kind/term (F14-REQ-02). |
| `F14-AD-007` | `Valid` vs `TopologyComplete` separate registration from completeness (DD-03, F14-REQ-12–F14-REQ-14). |
| `F14-AD-008` | One immutable typed immediate parent (DD-08, F14-REQ-09). |
| `F14-AD-009` | Immutable `metadata.uid` + scoped identity per stack (F14-REQ-15). |
| `F14-AD-010` | One `datacenterFailureDomainRef`; no spanning (F14-REQ-16). |
| `F14-AD-011` | `spec.technology` and `geo` descriptive only (DD-04, DD-05, F14-REQ-17, F14-REQ-27). |
| `F14-AD-012` | No inference from containment/ancestry/proximity (F14-REQ-19). |
| `F14-AD-013` | Connectivity `Unknown` by absence; no connectivity field (F14-REQ-20). |
| `F14-AD-014` | FEATURE-0012 grammar reused unchanged; no new grammar (F14-REQ-21, F14-REQ-28, F14-REQ-29). |
| `F14-AD-015` | `ManagedResource`/operator-facing; status system-owned current facts (DD-03, F14-REQ-22). |
| `F14-AD-016` | No FEATURE-0013 adoption machinery; `NOT_APPLICABLE` (F14-REQ-23). |
| `F14-AD-017` | No ResourcePool/ProviderCapability/placement fields (F14-REQ-24). |
| `F14-AD-018` | No adapter/discovery/provider-integration artifacts (F14-REQ-25). |
| `F14-AD-019` | No provider-native identity/reference; redaction (F14-REQ-18, F14-REQ-31). |
| `F14-AD-020` | Deletion rejected while children exist; no cascade/reparent (F14-REQ-26). |
| `F14-AD-021` | Single active `InfrastructureStack` name; no alias (DD-01, F14-REQ-03). |

## 12. Risk-control realization (`F14-R01`–`F14-R30`)

Each risk maps to a concrete design component, schema constraint, validation,
fixture, or review mechanism. Ratings, owners, and residual-risk acceptance are
copied unchanged from architecture section 12 and are not modified here
(architecture section 12.1).

| Risk | Design control realization | Evidence mechanism | Target residual (unchanged) | Owner (unchanged) |
|---|---|---|---|---|
| `F14-R01` | No placement/capacity/capability fields (4.3); topology-only schema | field deny-list; negative fixtures | 1×5 Medium | FEATURE-0014 owner |
| `F14-R02` | Single location kind/term (DD-01) | terminology scan; route/schema inventory | 1×3 Low | API architecture owner |
| `F14-R03` | `scopeRef` sole authority; `ownerRef` not scope (DD-08, 7.1 stage 4) | pos/neg scope + cross-scope authz tests | 1×5 Medium | Security/API owner |
| `F14-R04` | Reuse `Organization`; no owner kind (4.2) | kind inventory; owner-scenario fixtures | 1×3 Low | Architecture owner |
| `F14-R05` | Role-specific resources (F14-REQ-11) | same-/distinct-party fixtures | 1×3 Low | Domain owner |
| `F14-R06` | No inference logic; no connectivity field (4.3, 8) | same/cross-parent non-inference fixtures | 1×5 Medium | Network/placement owner |
| `F14-R07` | Schema deny-list excludes adjacency/boolean (4.3) | schema deny-list; semantic review | 1×4 Low | Network architecture owner |
| `F14-R08` | Bounded descriptive tech; redaction (DD-05, 7.2) | provider-neutral scan; secret/redaction tests | 1×5 Medium | Security/integration owner |
| `F14-R09` | `Valid`/`TopologyComplete` separation (DD-03) | empty-parent + completeness-transition fixtures | 1×4 Low | Domain owner |
| `F14-R10` | Immutable UID identity (4.3); technology not a key | duplicate-tech/distinct-UID fixtures | 1×3 Low | API owner |
| `F14-R11` | One immutable failure-domain parent (4.2) | zero/multiple/wrong-kind/cross-scope fixtures | 1×5 Medium | Domain owner |
| `F14-R12` | Required parent-ref kind + same scope (7.1 stages 2–3) | skipped-level/wrong-parent fixtures | 1×5 Medium | Domain/API owner |
| `F14-R13` | Deletion stage 6 rejects with children; leaf-first | child-existence conflict; concurrent/stale/retry tests | 2×4 Medium | Lifecycle/API owner |
| `F14-R14` | Status = current Sovrunn facts only (8) | staleness/unknown-state fixtures; writer review | 2×3 Medium | Operator/domain owner |
| `F14-R15` | Operator-facing; no-existence disclosure (9.1) | cross-scope list/get/reference + redaction tests | 1×5 Medium | Security owner |
| `F14-R16` | `geo` declared topology, not compliance (DD-04) | invalid/unknown-code + no-residency-inference review | 2×4 Medium | Sovereignty/domain owner |
| `F14-R17` | Finite bounds/pagination (DD-06, 6.4) | boundary/property + pagination benchmarks | 2×3 Medium | API/performance owner |
| `F14-R18` | Immutable parents + ETag/If-Match (7.3) | concurrency/stale-write/condition-consistency tests | 2×3 Medium | API/lifecycle owner |
| `F14-R19` | Single active name (DD-01) | old-name scan; schema/API diff | 1×4 Low | Architecture/API owner |
| `F14-R20` | Technology descriptive, not enum (DD-05) | enum/decision-use scan; negative compatibility tests | 1×4 Low | Domain/placement owner |
| `F14-R21` | Current-fact conditions; F13 `NOT_APPLICABLE` (8, 1.2) | status-history scan; FEATURE-0013 adoption gate | 1×3 Low | Architecture/audit owner |
| `F14-R22` | No integration/runtime artifacts (3) | changed-file/dependency/interface scan; drift review | 1×5 Medium | Architecture owner |
| `F14-R23` | Immutable non-reused UID; deletion rejection (4.3, 7.1) | UID non-reuse/stale-reference fixtures | 1×4 Low | API/later-feature owner |
| `F14-R24` | Independent scoped identity; no native authz (9.1) | same-operator/different-owner isolation + denial tests | 1×5 Medium | Security/domain owner |
| `F14-R25` | Delegated-only decisions (section 2); NO_UNRESOLVED_SEMANTIC_GAP (2.1) | ambiguity review; stage diff | 1×5 Medium | Architecture owner |
| `F14-R26` | Context bound to approved package; single active vocabulary (1) | terminology/ownership scan; context-manifest inspection | 1×5 Medium | Architecture/tooling owner |
| `F14-R27` | One design element per requirement; no duplicate norm (10, 13) | duplicate/similarity review | 1×4 Low | Requirements owner |
| `F14-R28` | Exact-ID enumeration in sections 10–12; orphan report (13) | automated decision/risk enumeration; orphan report | 1×5 Medium | Traceability/tooling owner |
| `F14-R29` | Design adds no new norm; delegated-only choices (2, 16.4 compliance) | design-language/semantic-delta review; stage diff | 1×5 Medium | Kiro workflow owner |
| `F14-R30` | Adjacent-feature absence ledger (section 13); deny-list scans | forbidden-concept scan; interface/dependency review | 1×5 Medium | Architecture/tooling owner |

Residual-risk acceptance remains human-owned and pending implementation evidence
and semantic review (architecture section 12.1); this design accepts, downgrades,
merges, renumbers, or closes no risk.

## 13. Absence ledger for adjacent-feature semantics

Explicit record that FEATURE-0013/0015/0016/0053 machinery is absent and where its
absence is enforced. No such section is designed; this ledger records absence and
the enforcing mechanism only.

| Owning feature | Absent concepts | Enforcement in this design |
|---|---|---|
| FEATURE-0013 | `DecisionRecord`, `DecisionProfile`, `EvaluationResult`, `AuditEvent`, rationale, obligation, actor, subject, linkage, projection, composition | `NOT_APPLICABLE` (1.2); status is current-fact conditions only (8); forbidden-concept scan (F14-REQ-23) |
| FEATURE-0015 | `ResourcePool`, `ProviderCapability`, capacity, compatibility, eligibility, placement, matching | No such field in any inventory (4.3); field deny-list (F14-REQ-24) |
| FEATURE-0016 | adapter interfaces, provider clients, discovery, credentials, endpoints, provider calls, repositories, plugins, provisioning, runtime execution | No such package/interface/dependency (3); interface/SDK scan (F14-REQ-25) |
| FEATURE-0053 | connectivity boolean, adjacency, route, peer, reachability, latency, bandwidth, trust, health, `NetworkConnectivityProfile` | No connectivity field; `Unknown` by absence (4.3, 8); schema deny-list (F14-REQ-19, F14-REQ-20) |

## 14. Non-goals and unresolved semantic-gap result

- Non-goals `NG-01`–`NG-10` from requirements section 5 are preserved as
  enforceable exclusions realized by the negative and architecture-drift
  acceptance criteria of their owning requirements (sections 10–12) and by the
  absence ledger (section 13). This design adds no independent normative
  obligation.
- No generic operation, audit, adapter, repository/storage, or provider
  integration section is designed (architecture section 16.4).
- Unresolved semantic-gap result: `NO_UNRESOLVED_SEMANTIC_GAP` (section 2.1). No
  `ARCHITECTURE_DECISION_REQUIRED` and no `REPOSITORY_CONTEXT_NOT_READY` condition
  was encountered while authoring this design against the approved package.

## 15. Design-to-requirements completeness and orphan report

Deterministic completeness self-check (architecture section 16.6):

- Exactly five owned resource kinds and zero active alias (DD-01, DD-02;
  F14-REQ-01, F14-REQ-03).
- All 21 decisions `F14-AD-001`–`F14-AD-021` individually mapped to a design
  disposition (section 11; no range shorthand).
- All 30 risks `F14-R01`–`F14-R30` individually mapped to a design control and
  evidence mechanism with unchanged rating and owner (section 12).
- Every requirement `F14-REQ-01`–`F14-REQ-31` maps to a representation, validator,
  and error/evidence behavior, and every design element in sections 3–9 traces
  back to at least one requirement (section 10). No orphan design element exists:
  each package, field, condition, route, limit, and fixture cites an originating
  requirement.
- No ResourcePool, ProviderCapability, capacity, eligibility, compatibility,
  adapter, discovery, credential, endpoint, provider-call, plugin, operation,
  provisioning, runtime, or connectivity semantic appears outside marked
  non-goal/ownership/absence text (sections 13–14).
- FEATURE-0013 adoption is exactly `NOT_APPLICABLE`; no decision/audit/operation
  semantic was introduced (F14-REQ-23).
- No new normative requirement, kind, relationship, status meaning, condition
  meaning, error family, operation, or interface was introduced; every design
  choice is a delegated representation (sections 2, 10 "New norm?" column all
  "No").
- No placeholder (`TBD`, `TODO`, "implementation-defined", "as appropriate") and
  no unresolved semantic marker remains.
- Stage boundary: only `design.md` is produced; no requirements, architecture,
  handoff, automation, source, schema, test, or `tasks.md` artifact was created or
  modified. Task generation requires a separate human `APPROVED_FOR_TASKS`
  (architecture section 16.1).

Orphan report: 0 orphan requirements, 0 orphan decisions, 0 orphan risks, 0 orphan
design elements.
