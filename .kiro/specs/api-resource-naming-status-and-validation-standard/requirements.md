---
reuse_assessment_format_version: 1.0.0
---

# Requirements Document

Feature: FEATURE-0012 — API, Resource Naming, Status, and Validation Standard
Phase: Phase 2 — Reuse-First PaaS Fabric Foundation (initial adoption phase; the standard applies across phases)
Stage: Requirements

## Metadata

| Field | Value |
|---|---|
| Feature ID | FEATURE-0012 |
| Feature title | API, Resource Naming, Status, and Validation Standard |
| Phase | Phase 2 — Reuse-First PaaS Fabric Foundation (initial adoption phase; standard applies across phases) |
| Spec stage | Requirements |
| Branch | feature-0012-api-resource-naming-status-and-validation-standard |
| Depends on | FEATURE-0011 (Reuse Assessment Standard) |
| Controlling handoff | ADH-2026-012 (Approved) |
| Canonical architecture | docs/architecture/api-resource-standard.md |
| Canonical reuse standard | docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md |
| Classification | Extend |
| Status | Draft — pending review |

## FEATURE-0012 reuse summary

Feature identity: FEATURE-0012 — API, Resource Naming, Status, and
Validation Standard.

This summary is populated only from the approved architecture baseline
(`docs/architecture/api-resource-standard.md`) and ADH-2026-012. Field
definitions are owned by `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
and are not redefined here.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| Sovrunn API/resource meta-model and conformance foundation | Extend | Extend mature HTTP, OpenAPI/JSON Schema, Problem Details, JSON Pointer, ETag concurrency, and selected Kubernetes API conventions with Sovrunn-owned sovereign scope, boundary, ownership, compatibility, and conformance rules. | Approved | ADH-2026-012; RFC-0022 (DEC-0026, DEC-0027, DEC-0036) |

Sovrunn owns the resource-profile taxonomy, API and naming conventions,
common metadata and identity, scope and reference semantics, API-boundary
classification, ownership and mutability rules, status and condition
grammar, validation and error contracts, provider-neutrality constraints,
compatibility policy, conformance rules, and reassessment triggers. HTTP
semantics, OpenAPI 3.1, JSON Schema 2020-12, RFC 9457 Problem Details,
RFC 6901 JSON Pointer, ETag/`If-Match` concurrency, and selected Kubernetes
API conventions are the reused or extended responsibility. The approved
feature-level disposition is Extend. Adapter required: No — FEATURE-0012
defines contract grammar and performs no external integration; later
integrations must use DEC-0036-compliant adapter boundaries.

The metadata Classification (Extend) does not replace this feature-level
reuse summary.

## Introduction

FEATURE-0012 defines the Sovrunn-owned, provider-neutral API and resource
grammar. This standard applies to ALL later applicable Sovrunn features
ACROSS PHASES that name, version, identify, scope, reference, classify,
validate, update, observe, or report externally exchanged objects.
Phase 2 is the initial adoption phase of the standard, not the limit of its
scope; FEATURE-0013 is the first adopting consumer, and later applicable
features in Phase 2 and subsequent phases adopt the same grammar. The
standard governs how Sovrunn objects are named, versioned, identified,
scoped, referenced, classified by audience and lifecycle, validated,
updated, observed through status and conditions, and reported through
machine-readable errors, and how they evolve without silent
incompatibility.

The architecture horizon is deliberately broad and cross-phase; the
FEATURE-0012 implementation scope is deliberately narrow. This feature
delivers shared platform contracts, shared implementation primitives,
machine-readable schema conventions, strict validation helpers, conformance
fixtures, a Phase 1 compatibility report, and feature-gate checks. The
normative boundary on what this feature may and may not contain lives in
F12-IMPL-001 and F12-IMPL-002; in particular, the domain behavior, runtime,
and persistence of any later feature stay outside this feature's
implementation.

Normative strength uses `MUST`, `MUST NOT`, `SHOULD`, and `MAY`. Every
normative requirement carries a stable `F12-<TOPIC>-<NNN>` identifier for
traceability; the stability and non-reuse of these identifiers is governed
by F12-TRACE-001. (The Matrix E risk identifiers `F12-R01`..`F12-R16` are a
separate, pre-existing numbering space and are retained unchanged.) These requirements convert every normative decision in
the approved architecture baseline into testable requirements, and include
the architecture A–E matrices and fitness functions as acceptance criteria.

Not every acceptance criterion is machine-verifiable. Automatable criteria
are verified by executable checks (`make fmt`, `make test`, `make vet`, and
the FEATURE-0012 feature gate); architecture approvals, exceptions, and
residual-risk acceptance require recorded human semantic review. Section
4.18 states these as two distinct verification requirements, and section 10
maps every architecture topic and Sovrunn foundation principle to the
requirement identifiers that satisfy it.

When an approved architecture source is insufficient to supply a required
semantic decision, the governing rule is F12-GOV-001: planning halts and
reports the F12-GOV-001 token rather than inventing
missing values, owners, dates, boundaries, limits, or risk decisions. The
normative statement of this discipline lives in F12-GOV-001.

## Glossary

Terms below are introduced or constrained by this feature. Canonical
platform terminology remains owned by `docs/glossary.md`.

| Term | Meaning in this feature |
|---|---|
| Resource profile | One of eight approved object shapes (`ManagedResource`, `ObservedExternalResource`, `VersionedDefinition`, `ImmutableRecord`, `LongRunningOperation`, `TransientRequestResult`, `EmbeddedValue`, `ListEnvelope`) selected by an externally exchanged object; profile selection is governed by F12-PROFILE-001. |
| API boundary | Classified audience for a schema: `customer-facing`, `operator-facing`, `internal-engine-facing`, `adapter-facing`, `plugin-facing`, or `governance-only`. |
| Scope kind (`scopeRef.kind`) | Governance/supply ownership context expressed by `scopeRef.kind`, limited to: `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`. Resource-local ownership is a distinct concept and is not among these scope kinds (see F12-OWNER-001). |
| `scopeRef` | Immutable primary security/governance ownership reference; not a location and not a lifecycle-containment reference. |
| Resource-local ownership | Lifecycle containment of a child object by a parent resource, expressed through `ownerRef`. It is a distinct concept from a formal security/governance scope; the rules that separate it from `scopeRef.kind` are governed by F12-OWNER-001. |
| `ownerRef` | Lifecycle-containment reference expressing resource-local ownership; its separation from security scope and from `scopeRef.kind` is governed by F12-OWNER-001. |
| `sourceRef` | Origin of an externally observed object. |
| `subjectRef` | Object a record or decision is about. |
| Typed reference | A `*Ref`/`*Refs` field carrying `apiVersion`, `kind`, `name`, and optional immutable `uid`, constrained by allowed kinds, scopes, and direction. |
| Condition | A stable, machine-readable current-fact observation (`type`, `status`, `reason`, `observedGeneration`, `lastTransitionTime`), not event history. |
| Problem | An RFC 9457 Problem Details response with Sovrunn extensions (`code`, `requestId`, `violations`). |
| Data classification | One of `Public`, `Customer-visible`, `Tenant-confidential`, `Operator-confidential`, `Internal`, `Sensitive`, `Secret-reference-only`. |
| Fitness function | An executable conformance/feature-gate check that enforces an architecture invariant. |
| Boundary ledger | The record of purpose, owner, producers, consumers, data, controls, and migration/reassessment paths for each declared boundary. |
| Provenance / freshness | Source identity, observation time, and staleness state required for externally observed objects. |
| Reassessment trigger | A defined condition that leads to architecture reassessment before implementation proceeds; the individual triggers and their effect are governed by F12-TRIGGER-001..F12-TRIGGER-013. |

## 3. User stories

- As a Sovrunn architecture owner, I want every externally exchanged object
  to declare a profile, boundary, allowed scopes, and maturity, so that
  audience, sensitivity, and compatibility are explicit and enforceable.
- As a feature author building FEATURE-0013 and later applicable features
  across phases, I want shared metadata, reference, condition, problem, and
  validation primitives, so that I inherit the common grammar instead of
  reinventing it.
- As a platform engineer, I want strict decoding that rejects unknown and
  duplicate fields with stable codes and JSON Pointer paths, so that
  malformed or ambiguous input fails deterministically.
- As a security reviewer, I want data classification, redaction, and
  secret-reference-only rules, so that credentials and restricted data
  never appear in metadata, status, errors, or audit messages.
- As a tenant administrator, I want cross-tenant and cross-organization
  references denied by default without existence disclosure, so that
  isolation holds and inaccessible objects are not revealed.
- As a platform operator, I want normalized, provider-neutral core and
  customer contracts, so that adding or replacing a provider does not
  change customer-facing schemas.
- As an adapter/plugin author, I want distinct adapter-facing and
  plugin-facing boundaries, so that native handles and provider data stay
  behind versioned contracts and cannot leak into core.
- As an API consumer, I want optimistic-concurrency updates via opaque
  `resourceVersion`/`If-Match`, so that stale writes fail without
  overwriting current state.
- As a portal developer, I want bounded, opaque, deterministically ordered
  pagination, so that large collections are safe to page without leaking
  offsets or provider details.
- As a compatibility reviewer, I want maturity levels and a schema-diff
  gate, so that breaking changes are detected and require a new version and
  migration plan.
- As a Phase 1 maintainer, I want a compatibility report with explicit
  exceptions and migration candidates, so that existing resources coexist
  without a rewrite or silent reinterpretation.
- As a governed AI consumer, I want stable reason codes and authorized
  boundary-filtered views, so that explanations rely on structured context
  and never bypass boundaries or scrape human messages.
- As an architecture reviewer, I want explicit reassessment triggers, so
  that new integration classes force a fresh architecture review before
  implementation instead of silently reusing assumptions.

## Requirements

Criteria are grouped by architecture decision. Each group is testable
through unit, integration, negative, boundary, security, or compatibility
tests. Groups 4.15–4.17 embed architecture matrices A–E; group 4.16 embeds
the fitness functions. Requirement strength uses `MUST`/`MUST NOT`/
`SHOULD`/`MAY`. Every normative criterion carries a stable
`F12-<TOPIC>-<NNN>` identifier.

### 4.0 Scope and applicability of the standard

1. **F12-SCOPESTD-001**: THE api-resource-standard SHALL apply to all later
   applicable Sovrunn features across phases; Phase 2 is its initial
   adoption phase and is NOT the limit of its scope.
2. **F12-SCOPESTD-002**: WHERE a later applicable Sovrunn feature, in any
   phase, defines objects that are named, versioned, identified, scoped,
   referenced, classified, validated, updated, observed, or reported across
   an API boundary, THE feature SHALL adopt this standard's grammar and
   contracts.
3. **F12-SCOPESTD-003**: FEATURE-0013 is the first adopting consumer of this
   standard; THE standard SHALL NOT be described as limited to FEATURE-0013
   or to Phase 2.
4. **F12-SCOPESTD-004**: IF a later applicable feature cannot adopt a
   provision of this standard, THEN THE feature SHALL record an approved
   architecture exception with owner, corrective path, and reassessment
   trigger rather than silently diverging.
5. **F12-TRACE-001**: Requirement identifiers MUST remain stable across
   revisions and MUST NOT be reused for a different semantic requirement.
6. **F12-GOV-001**: IF approved architecture sources do not provide a
   required semantic decision, THEN planning MUST stop and report exactly
   `ARCHITECTURE_DECISION_REQUIRED`; missing values, owners, dates,
   boundaries, limits, or risk decisions MUST NOT be invented.

### 4.1 Type identity, naming, and canonical schema

1. **F12-NAMING-001**: Externally exchanged typed objects MUST declare
   `apiVersion` (`<domain>.sovrunn.io/{v1alpha1|v1beta1|v1}`) and a singular
   PascalCase `kind`.
2. **F12-NAMING-002**: JSON/YAML fields MUST be lowerCamelCase; HTTP
   collections MUST be plural kebab-case; resource `name` MUST be lowercase
   kebab-case, immutable, URL-safe, and unique within API group + kind +
   scope.
3. **F12-NAMING-003**: `displayName` MUST be mutable, human-readable,
   Unicode-permitted, and MUST NOT be used as identity.
4. **F12-NAMING-004**: New APIs adopting this standard MUST use
   domain-grouped versioned routes and MUST NOT introduce unversioned
   public endpoints; Phase 1 compatibility APIs MAY retain existing
   group/routes until separately migrated.
5. **F12-NAMING-005**: Each external contract MUST have exactly one
   canonical machine-readable schema; generated code, docs, examples, and
   SDKs are derivative and MUST be checked for consistency against it.
6. **F12-NAMING-006**: Every schema MUST declare machine-readable metadata
   for profile, boundary, allowed scopes, and stability (e.g.
   `x-sovrunn-profile`, `x-sovrunn-boundary`, `x-sovrunn-allowed-scopes`,
   `x-sovrunn-stability`).

### 4.2 Resource profiles (Matrix A)

1. **F12-PROFILE-001**: Every externally exchanged object MUST select
   exactly one approved profile from Matrix A; introducing a new profile
   MUST require an approved architecture change.

Matrix A — Resource profile matrix:

| Profile | Purpose | Persistent | Authoritative writer | Required shape | Mutability |
|---|---|---|---|---|---|
| `ManagedResource` | Desired-state platform object | Yes | Client/operator for `spec`; Sovrunn for `status` | type metadata, metadata, spec, status | Spec mutable; status system-owned |
| `ObservedExternalResource` | Normalized external observation | Usually | Adapter/discovery component | type metadata, metadata, status, provenance, freshness | Observation updates only |
| `VersionedDefinition` | Published capability/plugin/service/schema | Yes | Authorized publisher | type metadata, metadata, spec, optional status | Draft mutable; published version immutable |
| `ImmutableRecord` | Decision/audit/evidence/history record | Yes | Authorized producer | type metadata, metadata, record payload | Append-only |
| `LongRunningOperation` | Asynchronous lifecycle action | Yes | Requester (request); executor (status) | type metadata, metadata, spec, status | Request immutable after acceptance |
| `TransientRequestResult` | Evaluation/simulation/command DTO | Optional | Caller/evaluator | typed input/result | Interaction-scoped |
| `EmbeddedValue` | Nested value, no independent lifecycle | No | Parent writer | domain value only | Parent-owned |
| `ListEnvelope` | Paginated collection response | No | API server | type metadata, items, page metadata | Read-only |

2. **F12-PROFILE-002**: Profile invariants MUST hold and be tested:
   `ManagedResource` separates desired `spec` from observed `status`;
   `ObservedExternalResource` has no customer-owned `spec` and requires
   source, observation time, and freshness; `VersionedDefinition` published
   versions are immutable; `ImmutableRecord` is append-only with linked
   corrections; `LongRunningOperation` represents target, action, requester,
   idempotency/correlation, progress, retryability, and terminal result;
   `TransientRequestResult` MUST NOT be persisted only to satisfy the
   grammar; `EmbeddedValue` MUST NOT receive artificial identity/metadata;
   `ListEnvelope` uses opaque tokens and deterministic ordering.

### 4.3 Common metadata

1. **F12-META-001**: Persistent resources MUST use the applicable subset of
   common metadata: `name`, `uid`, `displayName`, `scopeRef`, `labels`,
   `annotations`, `generation`, `resourceVersion`, `createdAt`, `updatedAt`.
2. **F12-META-002**: Field ownership and mutability MUST be enforced: `name`
   (creator, immutable), `uid` (Sovrunn, immutable, globally unique, opaque,
   never reused), `displayName` (owner, mutable), `scopeRef` (creator on
   create, validated, normally immutable), `labels`/`annotations`
   (authorized/namespaced owners, bounded, no secrets),
   `generation`/`resourceVersion`/timestamps (system-only).
3. **F12-META-003**: `generation` MUST change when desired state changes;
   `resourceVersion` MUST change whenever the stored representation changes;
   timestamps MUST be UTC-normalized.
4. **F12-META-004**: Clients MUST treat `uid` and `resourceVersion` as
   opaque strings; the UID algorithm is an implementation detail deferred to
   design.

### 4.4 Scope model (Matrix B)

1. **F12-SCOPE-001**: Customer governance and provider supply MUST be
   modeled as distinct scope structures; `Provider` MUST NOT be a parent of
   `Project`.

Matrix B — Scope matrix (valid `scopeRef.kind` values only):

| Scope kind (`scopeRef.kind`) | Parent | Primary purpose | Typical resources |
|---|---|---|---|
| `Platform` | None | Global governance and catalog | PluginDefinition, global policy |
| `Organization` | Platform | Customer/MSP administrative boundary | Organization policy, shared catalog |
| `OrganizationUnit` | Organization or OrganizationUnit | Delegation and grouping | Delegated assignments/controls |
| `Tenant` | Organization or OrganizationUnit | Isolation and entitlement | Tenant policies, quotas |
| `Project` | Tenant | Workload/service ownership | ServiceInstance, PlacementRequest |
| `Provider` | Platform or Organization | Supply/infrastructure boundary | ResourcePool, AdapterConfiguration |

Note — Resource-local ownership is NOT a `scopeRef.kind`. Resource-local
ownership is lifecycle containment of a child object by a parent resource,
expressed only through `ownerRef` (e.g. child binding, component, or
operation contained by a parent resource). It is not a formal
security/governance scope, does not appear in Matrix B, and MUST NOT be
supplied as a `scopeRef.kind`. See F12-OWNER-001.

2. **F12-SCOPE-002**: Scope rules MUST hold and be tested: a persistent
   resource is platform-scoped or has exactly one immutable primary
   `scopeRef`; each kind declares allowed scope kinds; identity uniqueness
   is `API group + kind + scope UID + name`; scope movement normally
   requires recreate-and-migrate; authorization resolves scope by UID, not
   display name; `scopeRef`, `ownerRef`, `location`, `sourceRef`, and
   `subjectRef` are distinct concepts; cross-tenant and cross-organization
   references are denied by default; denied cross-scope resolution MUST NOT
   disclose target existence; customer-facing provider/pool selection MUST
   use approved domain contracts, not provider-native identifiers.
3. **F12-OWNER-001**: `ownerRef` MUST express lifecycle containment
   (resource-local ownership) only; `ownerRef` MUST NOT be used as a
   `scopeRef.kind`, MUST NOT be used as a security/governance scope, and
   MUST NOT substitute for the primary `scopeRef`. A valid `scopeRef.kind`
   MUST be one of the Matrix B scope kinds.

### 4.5 References

1. **F12-REF-001**: References MUST use a common typed base with
   domain-specific constrained aliases: singular fields end in `Ref`,
   collections in `Refs`; each carries `apiVersion`, `kind`, `name`, and
   optional immutable `uid`.
2. **F12-REF-002**: UID MAY be omitted in human-authored input and returned
   after resolution; when name and UID are both present they MUST identify
   the same object.
3. **F12-REF-003**: Each reference schema MUST constrain allowed kinds,
   scopes, and direction; provider-native IDs MUST NOT act as core
   references; secret values MUST be represented only through typed secret
   references.
4. **F12-REF-004**: Generic references MAY provide a shared implementation
   base; public schemas SHOULD expose domain-specific reference types.

### 4.6 API boundaries (Matrix C1)

1. **F12-BOUNDARY-001**: Every schema MUST declare exactly one API boundary
   from Matrix C1 and MUST honor its allowed/prohibited data.

Matrix C1 — Boundary matrix:

| Boundary | Consumers | Allowed | Prohibited |
|---|---|---|---|
| `customer-facing` | Portal, CLI, SDK, GitOps, tenant automation | Product intent, safe status, actionable errors | Provider internals, secrets, inaccessible tenant data |
| `operator-facing` | Platform, MSP, provider operators | Administrative state, normalized infra, diagnostics | Raw secrets, unrestricted customer data |
| `internal-engine-facing` | Policy, placement, entitlement, orchestration | Normalized internal contracts and decisions | Vendor SDK types as shared models |
| `adapter-facing` | External-system adapters | Translation contracts, provider handles, provenance | Leakage into customer/core schemas |
| `plugin-facing` | Plugin manager and implementations | Capabilities, operation contracts, validated results | Unrestricted control-plane access or policy bypass |
| `governance-only` | Architecture and review workflow | Decisions, assessments, approvals, traceability | Runtime credentials, customer secrets |

2. **F12-BOUNDARY-002**: AI MUST be treated as a consumer of an authorized
   boundary, never a privileged boundary.
3. **F12-BOUNDARY-003**: A single schema MUST NOT serve multiple audiences
   whose trust, sensitivity, or compatibility requirements differ; separate
   boundary-specific views MUST be preferred over hidden-field redaction.

### 4.7 Ownership and mutability (Matrix C2)

1. **F12-OWNER-002**: Every mutable field and every condition type MUST have
   exactly one authoritative writer per Matrix C2; competing writers MUST be
   rejected.

Matrix C2 — Ownership matrix:

| Area | Authoritative owner | Rule |
|---|---|---|
| Type metadata and schema | Sovrunn | Immutable for an object version. |
| Identity and scope | Sovrunn after validated creation | Immutable by default. |
| Desired `spec` | Declared customer/operator/system actor | Mutability defined per resource. |
| `status` | Declared controller/adapter/plugin | Normal clients cannot write status. |
| Each condition type | One registered producer | No competing writers. |
| Provider-native handle | Adapter integration layer | Adapter/operator visibility only. |
| Plugin result | Plugin through validated contract | Cannot redefine core status grammar. |
| Stable error code/problem type | Sovrunn | Versioned machine contract. |
| Human message | Producing component | Informational; clients must not parse it. |
| Namespaced extension | Registered extension owner | Schema, boundary, version, size controlled. |
| Secret value | External secret system | Never in metadata, status, errors, or audit messages. |

### 4.8 Desired state, status, phase, and conditions

1. **F12-STATUS-001**: Where applicable, `spec` MUST express
   desired/requested state and `status` MUST express observed/evaluated
   state.
2. **F12-STATUS-002**: `phase` MAY provide a coarse lifecycle summary;
   conditions MUST express independent current facts, not event history.
3. **F12-STATUS-003**: Condition rules MUST hold and be tested: `status` is
   `True`/`False`/`Unknown`; `type` and `reason` are stable PascalCase
   machine identifiers; `message` is human-readable and not a machine
   contract; `observedGeneration` identifies the evaluated desired-state
   generation; `lastTransitionTime` changes only on status change.
4. **F12-STATUS-004**: Historical decisions and actions MUST belong to
   FEATURE-0013 records, not current resource status.
5. **F12-STATUS-005**: Resources defining both `phase` and conditions MUST
   define consistency rules between them.

### 4.9 Validation

1. **F12-VALIDATION-001**: Validation MUST occur in ordered, distinguishable
   layers: (1) HTTP/content/size checks, (2) safe JSON/YAML decoding, (3)
   duplicate- and unknown-field rejection, (4) structural schema validation,
   (5) deterministic defaulting, (6) semantic validation, (7)
   reference/kind/scope validation, (8) authorization/policy validation, (9)
   deferred later-feature capability/runtime validation.
2. **F12-VALIDATION-002**: Unknown fields and duplicate keys MUST be
   rejected by default.
3. **F12-VALIDATION-003**: Defaulting MUST be deterministic, documented, and
   versioned.
4. **F12-VALIDATION-004**: Structural, semantic, reference, authorization,
   and policy failures MUST remain distinguishable in results.
5. **F12-VALIDATION-005**: Validation MUST be safe for offline schema checks
   where external state is not required.
6. **F12-VALIDATION-006**: Validation failures MUST identify stable codes
   and JSON Pointer field paths; messages and status MUST redact
   unauthorized or sensitive detail.
7. **F12-VALIDATION-007**: Validation limits MUST be finite for object size,
   nesting, metadata, conditions, references, and violation counts; exact
   initial limits are deferred to design and treated as reviewed platform
   configuration.

### 4.10 Error contract

1. **F12-ERROR-001**: Errors MUST use RFC 9457 Problem Details with Sovrunn
   extensions (`type`, `title`, `status`, `detail`, `instance`, `code`,
   `requestId`, `violations[]` with `field` JSON Pointer, `code`,
   `message`).
2. **F12-ERROR-002**: Baseline HTTP mappings MUST hold: 400 malformed, 401
   auth required, 403 authorization denied, 404 not found, 409
   lifecycle/uniqueness conflict, 412 stale resource version, 415
   unsupported media type, 422 structurally/semantically invalid, 500
   internal failure, 503 temporary dependency failure.
3. **F12-ERROR-003**: Problem types, `code` values, and violation codes MUST
   be stable contracts; human messages MAY evolve and MUST NOT be parsed by
   clients.
4. **F12-ERROR-004**: Responses MUST NOT expose credentials, stack traces,
   raw provider errors, sensitive policy inputs, or inaccessible resource
   details.

### 4.11 Updates and concurrency

1. **F12-UPDATE-001**: Initial normative operations MUST be create, get,
   list, full replace, and delete.
2. **F12-UPDATE-002**: Full replacement MUST use an opaque `resourceVersion`
   via HTTP ETag semantics; protected updates MUST use `If-Match` or an
   equivalent checked resource version; stale writes MUST fail without
   overwriting current state.
3. **F12-UPDATE-003**: PATCH behavior MUST be deferred until a specific
   patch contract is approved; status updates MUST use a separately
   authorized path or internal interface.
4. **F12-UPDATE-004**: Mutating APIs SHOULD support request correlation and
   idempotency where replay is possible.

### 4.12 Lists and pagination

1. **F12-LIST-001**: List responses MUST use the `ListEnvelope` profile with
   `items` and `page.nextPageToken`.
2. **F12-LIST-002**: Page tokens MUST be opaque; ordering MUST be
   deterministic; page size MUST be bounded; total count MAY be omitted.
3. **F12-LIST-003**: Tokens MUST NOT expose database offsets, provider
   details, or authorization context; list filtering MUST be limited to
   explicitly indexed fields.

### 4.13 Extensions

1. **F12-EXT-001**: Extensions MUST be an explicit escape path, not a
   bypass; each extension MUST have a namespaced owner, a versioned schema, a
   declared boundary and data classification, finite size/nesting, and
   validation/compatibility rules.
2. **F12-EXT-002**: Extensions MUST NOT contain secret values, MUST NOT leak
   provider-native types into customer/core contracts, and MUST NOT create a
   core decision dependency unless formally recognized.
3. **F12-EXT-003**: When an extension becomes required by multiple
   independent consumers or core decisions, it MUST be reviewed for
   promotion into a typed Sovrunn contract.

### 4.14 API evolution and compatibility

1. **F12-EVOLVE-001**: Maturity levels MUST carry compatibility
   expectations: Alpha (breaking changes only with explicit migration notes
   and review), Beta (breaking changes exceptional; migration and
   coexistence required), Stable (backward compatibility mandatory within
   the major version).
2. **F12-EVOLVE-002**: Change classification MUST be applied and testable
   via a schema-diff gate: add optional field = compatible (subject to
   default/unknown-field rules); add required field, remove/rename field,
   change field meaning/owner/mutability, narrow enum or validation range =
   breaking; add enum value or change reference target kind/scope =
   compatibility review required; expose internal data publicly =
   security/boundary review required; add registered extension = additive
   within its contract.
3. **F12-EVOLVE-003**: Storage representation, served API version, and
   generated language types MUST remain distinguishable.

### 4.15 Future compatibility scenarios and fixtures (Matrix D)

1. **F12-FIXTURE-001**: FEATURE-0012 MUST demonstrate that its grammar can
   represent Matrix D scenarios without implementing their domain behavior.

Matrix D — Architecture conformance scenarios:

| Scenario | Profile | Scope/boundary | Required proof |
|---|---|---|---|
| Customer creates Project | ManagedResource | Tenant / customer | Stable identity, scope, strict validation, spec/status. |
| Operator registers ResourcePool | ManagedResource | Provider / operator | Provider-neutral capabilities; no vendor fields in core. |
| Adapter discovers external database | ObservedExternalResource | Provider / adapter | Provenance, freshness, stale/deleted semantics. |
| Publisher releases plugin contract | VersionedDefinition | Platform / plugin | Immutable published version and compatibility metadata. |
| Operator installs plugin | ManagedResource | Platform/Provider / operator | Definition, installation, execution remain separate. |
| Operator configures adapter | ManagedResource | Provider / adapter | Secret references and native-config isolation. |
| Placement engine evaluates request | TransientRequestResult | Project / internal | Typed request/result without forced persistence. |
| Decision becomes auditable | ImmutableRecord | Project / governance | Immutable subject/actor/input refs for FEATURE-0013. |
| Future provisioning executes | LongRunningOperation | Target scope / plugin | Idempotency, progress, retry, cancel, terminal result. |
| Portal lists large collections | ListEnvelope | Customer/operator | Bounded opaque pagination and deterministic ordering. |
| Provider disconnects | ObservedExternalResource | Provider / adapter | Current, stale, unknown, absent remain distinct. |
| External object recreated under same name | ObservedExternalResource | Provider / adapter | UID prevents stale-reference rebinding. |
| Cross-tenant reference attempted | Any reference | Customer | Denial without target-existence disclosure. |
| New cloud provider added | Adapter-facing | Provider | No customer/core schema change. |
| New data-service plugin added | Definition + operation | Multiple | No core grammar redesign. |
| AI explains denial | Existing authorized view | Customer/internal | Stable codes and safe context, no message scraping. |
| Phase 1 resource migrated | ManagedResource | Compatibility API | Explicit version/migration; no silent reinterpretation. |

2. **F12-FIXTURE-002**: Representative contract fixtures MUST exist for:
   `Project`, `ResourcePool`, `DiscoveredDatabase`, `PluginDefinition`,
   `AdapterConfiguration`, `PlacementEvaluationRequest`, `Operation`,
   `AuditEvent`. Fixtures MUST prove schema fit, boundary classification,
   allowed scopes, ownership, strict parsing, reference behavior, and the
   absence of later-phase execution.

### 4.16 Evolution safety: boundary ledger and fitness functions

1. **F12-LEDGER-001**: Every trust, API, scope, provider, plugin, adapter,
   lifecycle, data, and compatibility boundary MUST be recorded in a
   boundary ledger capturing: purpose; owner; producers; consumers;
   allowed/prohibited data; authorization; audit; observability; failure
   behavior; versioning; replacement path; migration path; reassessment
   trigger.
2. **F12-VERIFY-001**: The feature gate and conformance suite MUST implement
   executable fitness functions that check at least:

   1. every external schema declares profile, boundary, stability, and
      allowed scopes;
   2. no core/customer schema imports or embeds provider SDK/native types;
   3. every mutable field and condition has one owner;
   4. unknown and duplicate fields fail;
   5. references constrain kinds and scopes;
   6. cross-tenant access fails without existence disclosure;
   7. raw secret-like values are prohibited from metadata/status/errors;
   8. externally sourced observations include provenance and freshness;
   9. published definitions are immutable;
   10. schema compatibility detects breaking changes;
   11. object, metadata, condition, violation, reference, and page sizes
       are bounded;
   12. errors use stable codes and JSON Pointer paths;
   13. generated artifacts match the canonical schema;
   14. later-feature runtime behavior is absent;
   15. exceptions require an approved architecture handoff.

### 4.17 Risk register (Matrix E)

1. **F12-RISK-001**: FEATURE-0012 MUST carry the risk register below;
   residual risk is accepted only when documented with owner, corrective
   path, and reassessment trigger, and recorded through the human semantic
   review required by F12-VERIFY-002.

Matrix E — FEATURE-0012 risks and controls:

| ID | Risk | Primary control | Detection / response |
|---|---|---|---|
| F12-R01 | Overfit to Phase 1/Kubernetes | Profile matrix and provider-neutral rules | Future-scenario fixtures; ADH for exceptions. |
| F12-R02 | Universal shape misrepresents lifecycle | Eight explicit profiles | Schema profile required; reclassify via review. |
| F12-R03 | Scope/location/source/ownership conflated | Separate typed concepts | Scope/reference negative tests. |
| F12-R04 | Provider-native data leaks into core/customer | Adapter-only native contracts | Import/schema lint; move fields behind mapping. |
| F12-R05 | Plugin and adapter semantics conflated | Separate boundaries and resource families | Architecture classification tests. |
| F12-R06 | Generic/name-only refs cause access/staleness errors | Typed refs plus optional immutable UID | Scope/kind and name/UID mismatch tests. |
| F12-R07 | Multiple writers corrupt status | One owner per field/condition | Ownership registry and concurrency tests. |
| F12-R08 | Conditions become history or phase explodes | Current-fact conditions; optional phase | Bounded status tests; history to FEATURE-0013. |
| F12-R09 | Stale observations produce inaccurate decisions | Required provenance, observed time, freshness | Freshness tests; mark stale/unknown, fail safely. |
| F12-R10 | Secrets/restricted data leak via metadata/errors | Data classification, redaction, SecretRef-only | Security scanning and negative response tests. |
| F12-R11 | Extensions become shadow APIs | Registered namespaced schemas | Reject unregistered use; promote mature contracts. |
| F12-R12 | API evolution breaks clients/plugins | Maturity and compatibility policy | Schema-diff gate; new version and migration plan. |
| F12-R13 | Unbounded objects/status/lists prevent scale | Finite limits and opaque pagination | Limit tests and operational metrics. |
| F12-R14 | Documentation-only standard drifts | Shared primitives and executable conformance | Feature gate blocks non-conforming later features. |
| F12-R15 | Phase 1 migration expands into a rewrite | Compatibility audit and explicit exceptions | Scope gate; split migrations into approved features. |
| F12-R16 | AI/privileged consumers bypass boundaries | AI consumes authorized filtered views only | Access/redaction tests and audit. |

### 4.18 Implementation-contract boundary

1. **F12-IMPL-001**: The FEATURE-0012 implementation MUST contain only:
   normative architecture/standard documentation; shared type-metadata,
   metadata, reference, condition, problem, and validation primitives;
   machine-readable schema conventions and metadata; strict
   decoding/validation helpers; compatibility/conformance tooling; the eight
   required contract fixtures; a Phase 1 compatibility report; feature-gate
   checks; and unit and integration tests for positive, negative, boundary,
   security, and compatibility behavior.
2. **F12-IMPL-002**: The implementation MUST NOT create functional provider,
   plugin, policy, placement, audit, or provisioning services.
3. **F12-VERIFY-002**: Automatable acceptance criteria MUST be verified by
   executable checks through `make fmt`, `make test`, `make vet`, and the
   FEATURE-0012 feature gate (`make ff-feature-gate FEATURE=FEATURE-0012`).
   Not all acceptance criteria are machine-verifiable.
4. **F12-VERIFY-003**: Architecture approvals, granted exceptions, and
   residual-risk acceptance MUST be validated by recorded human semantic
   review; THE feature SHALL NOT claim that every acceptance criterion is
   verified by `make` commands alone. Semantic and governance approvals MUST
   record reviewer, decision, and date.

### 4.19 Reassessment triggers

1. **F12-TRIGGER-001**: IF a first non-Kubernetes provider integration is
   proposed, THEN THE Sovrunn architecture review process SHALL require
   architecture reassessment before implementation proceeds.
2. **F12-TRIGGER-002**: IF a first real provider adapter is proposed, THEN
   THE Sovrunn architecture review process SHALL require architecture
   reassessment before implementation proceeds.
3. **F12-TRIGGER-003**: IF a first remotely executed or data-path plugin is
   proposed, THEN THE Sovrunn architecture review process SHALL require
   architecture reassessment before implementation proceeds.
4. **F12-TRIGGER-004**: IF a first disconnected or federated control plane is
   proposed, THEN THE Sovrunn architecture review process SHALL require
   architecture reassessment before implementation proceeds.
5. **F12-TRIGGER-005**: IF a first external discovery source is proposed,
   THEN THE Sovrunn architecture review process SHALL require architecture
   reassessment before implementation proceeds.
6. **F12-TRIGGER-006**: IF cross-organization sharing is proposed, THEN THE
   Sovrunn architecture review process SHALL require architecture
   reassessment before implementation proceeds.
7. **F12-TRIGGER-007**: IF stable API promotion is proposed, THEN THE
   Sovrunn architecture review process SHALL require architecture
   reassessment before implementation proceeds.
8. **F12-TRIGGER-008**: IF high-frequency status updates are proposed, THEN
   THE Sovrunn architecture review process SHALL require architecture
   reassessment before implementation proceeds.
9. **F12-TRIGGER-009**: IF regulated or classified workloads are proposed,
   THEN THE Sovrunn architecture review process SHALL require architecture
   reassessment before implementation proceeds.
10. **F12-TRIGGER-010**: IF an object cannot select an approved profile, THEN
    THE Sovrunn architecture review process SHALL require architecture
    reassessment before implementation proceeds and THE object MUST NOT be
    force-fit into an existing profile.
11. **F12-TRIGGER-011**: IF a consumer requires a new privileged view, THEN
    THE Sovrunn architecture review process SHALL require architecture
    reassessment before implementation proceeds.
12. **F12-TRIGGER-012**: IF an extension becomes a core dependency or begins
    serving multiple consumers, THEN THE Sovrunn architecture review process
    SHALL require architecture reassessment before implementation proceeds.
13. **F12-TRIGGER-013**: IF any backward-incompatible migration is proposed,
    THEN THE Sovrunn architecture review process SHALL require architecture
    reassessment before implementation proceeds.

## 5. Non-goals

This section is a non-normative summary and introduces no new requirements.
It restates, for reader convenience, scope boundaries already governed by
the stable requirement identifiers cited in each item. The governing
requirements remain the sole normative source.

- Provider or substrate models (owned by FEATURE-0014/0015) are out of the
  FEATURE-0012 implementation, per the implementation-contract boundary
  (F12-IMPL-002; F12-IMPL-001).
- Adapter protocols (owned by FEATURE-0016) are out of the FEATURE-0012
  implementation, per F12-IMPL-002; adapter-native data stays behind the
  adapter-facing boundary (F12-BOUNDARY-001, F12-REF-003).
- Policy evaluation (owned by FEATURE-0017 and later) is out of the
  FEATURE-0012 implementation, per F12-IMPL-002.
- DecisionObject and AuditEvent domain payloads (owned by FEATURE-0013) are
  out of the FEATURE-0012 implementation, per F12-IMPL-002; historical
  decisions and actions belong to FEATURE-0013 records, not resource status
  (F12-STATUS-004).
- Placement behavior (owned by FEATURE-0023) is out of the FEATURE-0012
  implementation, per F12-IMPL-002.
- Plugin taxonomy and execution (owned by FEATURE-0024 and later phases) are
  out of the FEATURE-0012 implementation, per F12-IMPL-002.
- Provisioning, persistence selection, workflow execution, billing,
  failover, and autonomous AI operations are out of the FEATURE-0012
  implementation, per F12-IMPL-002.
- A wholesale rewrite of Phase 1 resources or routes is excluded; the
  compatibility report records exceptions and migration candidates without
  triggering a rewrite (F12-COMPAT-002, F12-COMPAT-003).
- Vendor-specific core models and provider-native core references are
  excluded from customer/core contracts (F12-REF-003, F12-EXT-002,
  F12-BOUNDARY-001).
- An unrestricted arbitrary extension mechanism is excluded; extensions are
  a governed escape path only (F12-EXT-001, F12-EXT-002).
- Exact HTTP route migration, PATCH format, watch/change-stream protocol,
  production storage/indexing, and stable API promotion are deferred and
  separately approved (F12-NAMING-004, F12-UPDATE-003, F12-TRIGGER-007,
  F12-TRIGGER-013).

Phase 2 remains a model, standard, decision, audit, adapter-boundary, and
simulation foundation. Consistent with F12-IMPL-001 and F12-IMPL-002,
FEATURE-0012 does not authorize real runtime integrations or later-phase
execution. Defining the standard as cross-phase (section 4.0,
F12-SCOPESTD-001) does not expand the FEATURE-0012 implementation scope:
FEATURE-0012 delivers the grammar and conformance tooling only, and later
applicable features adopt it under their own approvals (F12-SCOPESTD-002,
F12-SCOPESTD-004).

## 6. Edge cases

This section is a non-normative summary and introduces no new requirements.
Each item describes an edge case and points to the stable requirement
identifier(s) that already govern the expected behavior. The cited
requirements remain the sole normative source.

- An externally exchanged object fits no approved profile: the object is
  rejected and a reassessment trigger is raised, and it is not force-fit or
  invented into an existing profile (F12-TRIGGER-010, F12-PROFILE-001).
- Unknown or duplicate field present: rejected with a stable code and JSON
  Pointer path, with no silent drop (F12-VALIDATION-002, F12-VALIDATION-006).
- Reference supplies both `name` and `uid` that disagree: rejected as a
  mismatch (F12-REF-002).
- Reference target is in another tenant/organization: denied by default
  without disclosing whether the target exists (F12-SCOPE-002, F12-SEC-004).
- `ownerRef` supplied where a `scopeRef.kind` is expected: rejected, because
  resource-local ownership is not a scope kind (F12-OWNER-001).
- External object is deleted and recreated under the same name: UID mismatch
  prevents stale-reference rebinding (F12-REF-002, F12-SCOPE-002).
- External provider disconnects: the observation becomes explicitly
  stale/unknown/absent rather than fabricated as current (F12-SEC-005,
  F12-STATUS-003).
- Concurrent update with a stale `resourceVersion`: fails (412) without
  overwriting current state (F12-UPDATE-002, F12-ERROR-002).
- Two producers attempt to write the same condition type or status field:
  the non-authoritative writer is rejected (F12-OWNER-002).
- Secret-like value submitted in metadata, labels, annotations, status, or
  extension: rejected, because only typed secret references are allowed
  (F12-SEC-003, F12-REF-003, F12-EXT-002).
- Object, status, list page, violation list, or reference set exceeds a
  finite limit: rejected deterministically (F12-VALIDATION-007, F12-LIST-002).
- Client submits an unversioned or Phase-1-style route for a new adopting
  API: rejected, because versioned domain-grouped routes are required
  (F12-NAMING-004).
- Compatibility check detects a breaking schema change on an Alpha/Beta/
  Stable contract: blocked until a new version and migration plan exist
  (F12-EVOLVE-001, F12-EVOLVE-002).
- Generated artifact diverges from the canonical schema: reported as a
  conformance failure (F12-NAMING-005, F12-VERIFY-001).
- A required semantic decision is missing from approved sources: planning
  stops and reports `ARCHITECTURE_DECISION_REQUIRED` (F12-GOV-001).

## 7. Security and privacy requirements

1. **F12-SEC-001**: Every schema field crossing an API boundary MUST
   identify or inherit: data classification, authorized writer, authorized
   readers, mutability, retention, redaction behavior, residency
   implications, and audit requirement.
2. **F12-SEC-002**: Supported data classifications MUST be exactly:
   `Public`, `Customer-visible`, `Tenant-confidential`,
   `Operator-confidential`, `Internal`, `Sensitive`,
   `Secret-reference-only`.
3. **F12-SEC-003**: Least privilege MUST apply: strict decoding, typed
   scope/reference rules, bounded inputs, and secret-reference-only
   semantics are mandatory; raw secrets MUST NOT be stored in metadata,
   status, errors, or audit messages.
4. **F12-SEC-004**: Inaccessible objects MUST NOT be revealed through error
   differences, timing, or existence disclosure.
5. **F12-SEC-005**: External observations MUST expose provenance and
   freshness; disconnected dependencies MUST produce explicit
   degraded/unknown status rather than fabricated accuracy.
6. **F12-SEC-006**: Provider replacement MUST NOT require customer-contract
   changes; provider choice, location, and ownership scope MUST remain
   separable.
7. **F12-SEC-007**: Auditing MUST be able to link actors, requests,
   operations, subjects, versions, and later decisions without storing
   history in current resource status.
8. **F12-SEC-008**: AI-generated explanations MUST depend only on stable
   structured context and authorized boundary-filtered views; AI MUST
   receive no bypass boundary and MUST NOT scrape human messages.
9. **F12-SEC-009**: Open, locally processable schemas MUST support
   on-premise, disconnected, and air-gapped deployments; external
   availability MUST NOT be required for contract interpretation or
   structural validation.

## 8. Compatibility with completed Phase 1 features

1. **F12-COMPAT-001**: FEATURE-0012 MUST produce a Phase 1 compatibility
   report covering the completed Phase 1 resources and routes: Organization,
   OrganizationUnit, Tenant, Project, Operation, ServiceClass, ServicePlan,
   Plugin, Capability, ServiceInstance, ServiceBinding, and the
   health/readiness and demo-flow endpoints.
2. **F12-COMPAT-002**: The report MUST record, per contract, conforming
   behavior, explicit documented exceptions, and migration candidates; it
   MUST NOT trigger a wholesale rewrite.
3. **F12-COMPAT-003**: Phase 1 compatibility APIs MAY retain their existing
   API group and routes until a separately approved migration; their
   behavior MUST NOT be silently reinterpreted under the new grammar.
4. **F12-COMPAT-004**: New contracts adopting this standard MUST use the
   FEATURE-0012 grammar; existing Phase 1 contracts coexist during migration
   and any breaking change MUST follow the versioning, migration, and
   approval rules in 4.14.
5. **F12-COMPAT-005**: The metadata/spec/status resource shape, in-memory
   registry model, and stable error-code approach established in Phase 1
   MUST remain consistent with, and be generalized by, this standard.
6. **F12-COMPAT-006**: FEATURE-0013 and later applicable feature gates,
   across phases, MUST enforce adoption of this standard where relevant; a
   feature gate MUST fail an applicable feature that defines externally
   exchanged objects without adopting the grammar and contracts in this
   document, unless an approved architecture exception (F12-SCOPESTD-004) is
   recorded.

## 9. Design questions to resolve later in design.md

These questions are deliberately deferred. Design resolves them within the
approved contracts and boundaries established here; where a needed decision
is unavailable from approved architecture sources, design invokes F12-GOV-001
(architecture decision required), and where a provision cannot be adopted,
design invokes F12-SCOPESTD-004 (approved adoption exception) rather than
diverging silently. The governing normative force lives in those
requirements, not in this preamble:

1. What is the canonical schema source of truth and the derivative
   generation flow (code, docs, SDKs), and how is consistency enforced?
2. What are the Go package/module boundaries for type metadata, metadata,
   references, conditions, problems, validation, and conformance?
3. What is the concrete strict-decoding, validation-ordering, error-
   mapping, and concurrency implementation?
4. What are the exact finite platform limits (object size, nesting,
   metadata, labels/annotations, conditions, references, violations, page
   size) as reviewed configuration?
5. What is the concrete UID generation algorithm and the opaque
   `resourceVersion` representation?
6. How are machine-readable profile/boundary/scope/stability annotations
   represented in schemas and validated by fitness functions?
7. What is the exact HTTP route form for new adopting APIs, and how do
   Phase 1 compatibility routes coexist during migration?
8. How is the schema-compatibility (diff) gate implemented and wired into
   the feature gate?
9. How is the boundary ledger represented and kept in sync with schemas?
10. What is the testing strategy for positive, negative, boundary,
    security, and compatibility conformance, including the eight required
    fixtures?
11. How are PATCH, watch/change-stream, and status-update paths reserved
    without implementing them now?
12. How does the Phase 1 compatibility report capture exceptions and
    migration candidates in a reviewable, testable form?

Design, tasks, and implementation remain unauthorized until their separate
human approval tokens are issued. This requirements document is authorized
by ADH-2026-012 for the Requirements stage only.

## 10. Architecture Coverage Matrix

This final matrix maps each architecture topic and Sovrunn foundation
principle to the stable requirement identifiers that satisfy it. It exists
for traceability; it does not introduce new normative strength beyond the
referenced requirements.

| Architecture topic | Requirement IDs | Coverage status | Deferred owner |
|---|---|---|---|
| Accuracy | F12-STATUS-001, F12-STATUS-003, F12-SEC-005, F12-VALIDATION-004, F12-REF-002 | Covered | — |
| Transparency | F12-ERROR-001, F12-ERROR-003, F12-VALIDATION-006, F12-SEC-008, F12-VERIFY-003, F12-TRACE-001 | Covered | — |
| Security | F12-SEC-001, F12-SEC-002, F12-SEC-003, F12-SEC-004, F12-ERROR-004, F12-EXT-002, F12-BOUNDARY-003 | Covered | — |
| Sovereignty | F12-SEC-009, F12-SEC-006, F12-VALIDATION-005, F12-NAMING-001 | Covered | — |
| Auditability | F12-SEC-007, F12-STATUS-004, F12-PROFILE-002 (ImmutableRecord), F12-VERIFY-003, F12-RISK-001, F12-TRACE-001 | Covered | — |
| Flexibility | F12-EXT-001, F12-EXT-003, F12-EVOLVE-001, F12-EVOLVE-002, F12-PROFILE-001, F12-SCOPESTD-001, F12-SCOPESTD-002 | Covered | — |
| Provider neutrality | F12-SCOPE-001, F12-REF-003, F12-BOUNDARY-001, F12-SEC-006, F12-EXT-002, F12-VERIFY-001 (checks 2 and 8) | Covered | — |
| Tenant isolation | F12-SCOPE-002, F12-OWNER-001, F12-SEC-004, F12-BOUNDARY-001, F12-VERIFY-001 (check 6) | Covered | — |
| Governed AI | F12-BOUNDARY-002, F12-SEC-008, F12-ERROR-003 | Covered | — |
| Type identity and naming | F12-NAMING-001, F12-NAMING-002, F12-NAMING-003, F12-NAMING-004, F12-NAMING-005, F12-NAMING-006 | Covered | — |
| Resource profiles | F12-PROFILE-001, F12-PROFILE-002 | Covered | — |
| Common metadata and identity | F12-META-001, F12-META-002, F12-META-003, F12-META-004 | Covered | — |
| Scope model | F12-SCOPE-001, F12-SCOPE-002, F12-OWNER-001 | Covered | — |
| References | F12-REF-001, F12-REF-002, F12-REF-003, F12-REF-004 | Covered | — |
| API boundaries | F12-BOUNDARY-001, F12-BOUNDARY-002, F12-BOUNDARY-003 | Covered | — |
| Ownership and mutability | F12-OWNER-002, F12-OWNER-001 | Covered | — |
| Status and conditions | F12-STATUS-001, F12-STATUS-002, F12-STATUS-003, F12-STATUS-004, F12-STATUS-005 | Covered | — |
| Validation | F12-VALIDATION-001..F12-VALIDATION-007 | Covered | — |
| Error contract | F12-ERROR-001, F12-ERROR-002, F12-ERROR-003, F12-ERROR-004 | Covered | — |
| Updates and concurrency | F12-UPDATE-001, F12-UPDATE-002, F12-UPDATE-003, F12-UPDATE-004 | Covered | — |
| Lists and pagination | F12-LIST-001, F12-LIST-002, F12-LIST-003 | Covered | — |
| Extensions | F12-EXT-001, F12-EXT-002, F12-EXT-003 | Covered | — |
| API evolution and compatibility | F12-EVOLVE-001, F12-EVOLVE-002, F12-EVOLVE-003, F12-COMPAT-004 | Covered | — |
| Future compatibility scenarios and fixtures | F12-FIXTURE-001, F12-FIXTURE-002 | Covered | — |
| Boundary ledger and fitness functions | F12-LEDGER-001, F12-VERIFY-001 | Covered | — |
| Risk management | F12-RISK-001 (Matrix E: F12-R01..F12-R16) | Covered | — |
| Implementation-contract boundary | F12-IMPL-001, F12-IMPL-002 | Covered | — |
| Verification (automated + human) | F12-VERIFY-001, F12-VERIFY-002, F12-VERIFY-003 | Covered | — |
| Reassessment triggers | F12-TRIGGER-001..F12-TRIGGER-013 | Covered | — |
| Cross-phase applicability and adoption | F12-SCOPESTD-001, F12-SCOPESTD-002, F12-SCOPESTD-003, F12-SCOPESTD-004, F12-COMPAT-006, F12-GOV-001 | Covered | — |
| Traceability and governance | F12-TRACE-001, F12-GOV-001 | Covered | — |
| Phase 1 compatibility | F12-COMPAT-001, F12-COMPAT-002, F12-COMPAT-003, F12-COMPAT-004, F12-COMPAT-005 | Covered | — |
| DecisionObject/AuditEvent payloads | F12-IMPL-002, F12-STATUS-004 | Deferred | FEATURE-0013 |
| Provider models | F12-IMPL-002, F12-SCOPESTD-004 | Deferred | FEATURE-0014/0015 |
| Adapter protocols | F12-IMPL-002, F12-BOUNDARY-001 | Deferred | FEATURE-0016 |
| Policy evaluation | F12-IMPL-002 | Deferred | FEATURE-0017+ |
| Placement behavior | F12-IMPL-002 | Deferred | FEATURE-0023 |
| Plugin taxonomy/execution | F12-IMPL-002 | Deferred | FEATURE-0024 and later phases |
