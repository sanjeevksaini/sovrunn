---
doc_type: feature_architecture
feature: FEATURE-0012
title: API, Resource Naming, Status, and Validation Standard
status: approved_for_kiro
phase: 2
classification: architecture_baseline
reuse_disposition: Extend
ai_load_priority: critical
controlling_documents:
  - docs/phase2/PHASE2_ARCHITECTURE_SPINE.md
  - docs/phase2/PHASE2_EXECUTION_STRATEGY.md
  - docs/phase2/PHASE2_SCOPE.md
  - docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md
  - docs/context/CURRENT_ARCHITECTURE_BASELINE.md
required_handoff:
  - ADH-2026-012
  - ADH-2026-013
standard_maturity: draft
approval_date: 2026-07-22
approving_role: Sovrunn Architecture Owner
kiro_authorization: approved_for_requirements
---

# FEATURE-0012 Architecture

> **Approval record — 2026-07-22:** The FEATURE-0012 architecture has no identified unmitigated conflict with Sovrunn foundation principles or evaluated growth scenarios. Its intentional boundaries are explicit, owned, versioned, observable, auditable, replaceable, and testable. Remaining uncertainties have documented migration paths and reassessment triggers. Approved as the baseline for Kiro requirements, design, and tasks.
>
> This approval authorizes requirements generation only. Design, tasks, and Cursor implementation retain their separate approval gates. FEATURE-0011 applies as cross-phase governance; its current canonical repository path is retained for compatibility until a separately approved migration.


## 1. Purpose

FEATURE-0012 defines the provider-neutral API and resource grammar used by FEATURE-0013 and later Sovrunn features.

It standardizes how Sovrunn objects are:

- named and versioned;
- identified and scoped;
- referenced across boundaries;
- classified by lifecycle and audience;
- validated and updated;
- observed through status and conditions;
- reported through machine-readable errors;
- evolved without silent incompatibility.

FEATURE-0012 defines shared platform contracts, not the domain behavior of future resources.

## 2. Architecture intent

> Use a Kubernetes-inspired declarative resource model, extended into a Sovrunn-owned, HTTP-native, provider-neutral contract.

The architecture must be broad enough to support future resources, plugins, adapters, data sources, consumers, disconnected providers, and regulated workloads. The implementation must remain narrow: common types, schemas, validators, conformance tests, and Phase 1 compatibility analysis only.

```text
Broad architecture horizon.
Narrow FEATURE-0012 implementation scope.
Explicit boundaries.
Versioned evolution.
Executable conformance.
```

Normative terms `MUST`, `MUST NOT`, `SHOULD`, and `MAY` define requirement strength.

## 3. Reuse assessment

### 3.1 Disposition

```text
Feature-level disposition: Extend
```

### 3.2 Reused or extended foundations

Sovrunn will reuse or extend:

- Kubernetes API conventions for type metadata, desired/observed separation, and conditions;
- OpenAPI 3.1 and JSON Schema 2020-12 for machine-readable schemas;
- HTTP semantics for methods, status codes, conditional requests, and caching metadata;
- RFC 9457 Problem Details for error transport;
- RFC 6901 JSON Pointer for field paths;
- opaque ETag and `If-Match` semantics for optimistic concurrency.

### 3.3 Sovrunn-owned responsibilities

Sovrunn owns:

- the resource-profile taxonomy;
- API group, kind, field, collection, and resource-name conventions;
- common metadata, scope, identity, and reference semantics;
- customer, operator, internal, adapter, plugin, and governance boundaries;
- field ownership, mutability, status, and condition rules;
- validation stages, error codes, redaction rules, and extension governance;
- provider-neutrality constraints;
- compatibility, migration, conformance, and reassessment policy.

### 3.4 Responsibility boundary

External standards define generic syntax and semantics. Sovrunn defines how they are constrained and combined for a sovereign PaaS. Provider-, plugin-, and adapter-native types remain behind explicit boundaries and MUST NOT become Sovrunn core contracts.

## 4. Scope

### 4.1 In scope

- resource and object profiles;
- type metadata and naming;
- common metadata and identity;
- scope and reference grammar;
- API-boundary classification;
- ownership and mutability;
- `spec`, `status`, phase, and conditions;
- schema and semantic validation;
- problem/error responses;
- API maturity and compatibility rules;
- pagination and optimistic concurrency conventions;
- controlled extension points;
- shared implementation primitives and conformance fixtures;
- Phase 1 compatibility audit and explicit exceptions.

### 4.2 Non-goals

FEATURE-0012 MUST NOT implement:

- provider or substrate models owned by FEATURE-0014/0015;
- adapter protocols owned by FEATURE-0016;
- policy evaluation owned by FEATURE-0017+;
- DecisionObject or AuditEvent domain payloads owned by FEATURE-0013;
- placement behavior owned by FEATURE-0023;
- plugin taxonomy or execution owned by FEATURE-0024 and later phases;
- provisioning, persistence selection, workflow execution, billing, failover, or autonomous AI;
- a wholesale rewrite of Phase 1 resources or routes;
- vendor-specific core models;
- an unrestricted arbitrary extension mechanism.

## 5. Foundation-principle traceability

| Foundation principle | FEATURE-0012 architectural response |
|---|---|
| Reuse before build | Extend mature API, schema, HTTP, condition, and problem standards. |
| Provider-neutral core | Prohibit provider SDK types and provider-native fields in core/customer contracts. |
| Organization-first governance | Use explicit immutable scope and typed cross-scope references. |
| Tenant isolation | Deny cross-tenant references by default without existence disclosure. |
| Plugin-based delivery | Share common grammar while keeping plugin contracts semantically distinct. |
| Adapter boundaries | Translate external systems through versioned adapter-facing contracts. |
| Accuracy | Require strict validation, provenance, freshness, stable identities, and stale-reference detection. |
| Transparency | Require stable reason codes, conditions, boundaries, owners, and explainable failures. |
| Security | Apply least privilege, data classification, redaction, secret-reference-only rules, and bounded inputs. |
| Sovereignty | Keep provider choice replaceable; represent location separately from ownership scope; avoid mandatory remote control dependencies. |
| Auditability | Preserve immutable identity, actor/request/operation correlation, and references suitable for FEATURE-0013. |
| Flexibility | Use profiles, API versions, typed extensions, adapters, plugins, and explicit migration paths. |
| Governed AI | AI consumes authorized boundary-filtered contracts and stable reason codes; it receives no bypass boundary. |

## 6. Core architecture decisions

### 6.1 Canonical type identity

Externally exchanged typed objects MUST identify their contract with:

```yaml
apiVersion: <domain>.sovrunn.io/v1alpha1
kind: <SingularPascalCaseKind>
```

Rules:

- API group: lowercase DNS-style domain, grouped by Sovrunn domain.
- Version: `v1alpha1`, `v1beta1`, or `v1` maturity form.
- Kind: singular PascalCase.
- JSON/YAML field: lowerCamelCase.
- HTTP collection: plural kebab-case.
- Resource name: lowercase kebab-case, immutable, URL-safe, and unique within API group + kind + scope.
- Human label: mutable `displayName`, Unicode permitted.
- Phase 1 compatibility APIs MAY retain their existing group and routes until separately migrated.

The exact HTTP route migration is deferred. New Phase 2 APIs MUST use domain-grouped versioned routes and MUST NOT introduce unversioned public endpoints.

### 6.2 Canonical schema

Each external contract MUST have one canonical machine-readable schema. Generated code, documentation, examples, and SDKs are derivative artifacts and MUST be checked for consistency.

Every schema MUST declare machine-readable metadata equivalent to:

```yaml
x-sovrunn-profile: ManagedResource
x-sovrunn-boundary: customer-facing
x-sovrunn-allowed-scopes:
  - Project
x-sovrunn-stability: alpha
```

### 6.3 Resource profiles

Every externally exchanged object MUST select one approved profile.

#### A. Resource profile matrix

| Profile | Purpose | Persistent | Authoritative writer | Required shape | Mutability |
|---|---|---:|---|---|---|
| `ManagedResource` | Desired-state platform object | Yes | Client/operator for `spec`; Sovrunn for `status` | type metadata, metadata, spec, status | Spec mutable; status system-owned |
| `ObservedExternalResource` | Normalized observation from an external system | Usually | Adapter/discovery component | type metadata, metadata, status, provenance, freshness | Observation updates only |
| `VersionedDefinition` | Published capability, plugin, service, or schema definition | Yes | Authorized publisher | type metadata, metadata, spec, optional status | Draft mutable; published version immutable |
| `ImmutableRecord` | Decision, audit, evidence, or history record | Yes | Authorized producer | type metadata, metadata, record payload | Append-only |
| `LongRunningOperation` | Asynchronous lifecycle action | Yes | Requester for immutable request; executor for status | type metadata, metadata, spec, status | Request immutable after acceptance |
| `TransientRequestResult` | Evaluation, simulation, or command DTO | Optional | Caller/evaluator | typed input/result | Interaction-scoped |
| `EmbeddedValue` | Nested value without independent lifecycle | No | Parent writer | domain value only | Parent-owned |
| `ListEnvelope` | Paginated collection response | No | API server | type metadata, items, page metadata | Read-only |

Profile invariants:

- `ManagedResource`: `spec` is desired state; `status` is observed state.
- `ObservedExternalResource`: no customer-owned `spec`; source, observation time, and freshness are required.
- `VersionedDefinition`: published versions are immutable; changed contracts create new versions.
- `ImmutableRecord`: updates are prohibited; corrections create linked records.
- `LongRunningOperation`: target, action, requester, idempotency/correlation context, progress, retryability, and terminal result are representable.
- `TransientRequestResult`: MUST NOT be persisted merely to satisfy the resource grammar.
- `EmbeddedValue`: MUST NOT receive artificial identity or metadata.
- `ListEnvelope`: page tokens are opaque and ordering deterministic.

A new profile requires an approved architecture change.

### 6.4 Common metadata

Persistent resources use the applicable subset of:

```yaml
metadata:
  name: payments-production
  uid: <opaque-server-generated-id>
  displayName: Payments Production
  scopeRef: {}
  labels: {}
  annotations: {}
  generation: 4
  resourceVersion: <opaque-version>
  createdAt: <timestamp>
  updatedAt: <timestamp>
```

| Field | Owner | Mutable | Rule |
|---|---|---:|---|
| `name` | Authorized creator | No | Stable URL-safe identity within scope. |
| `uid` | Sovrunn | No | Globally unique, opaque, never reused. |
| `displayName` | Authorized resource owner | Yes | Human-readable; not identity. |
| `scopeRef` | Creator on create; validated by Sovrunn | Normally no | Security/governance ownership, not location. |
| `labels` | Authorized actors | Yes | Bounded, indexed classification; no secrets. |
| `annotations` | Namespaced owners | Yes | Bounded metadata; no secrets or ungoverned API contracts. |
| `generation` | Sovrunn | System-only | Changes when desired state changes. |
| `resourceVersion` | Sovrunn | System-only | Changes whenever stored representation changes. |
| timestamps | Sovrunn | System-only | UTC and normalized. |

The concrete UID algorithm is an implementation detail. Clients MUST treat UIDs and resource versions as opaque strings.

### 6.5 Scope model

Sovrunn models customer governance and provider supply as related but distinct scope structures.

```text
Customer governance:
Platform -> Organization -> OrganizationUnit -> Tenant -> Project

Provider supply:
Platform or Organization -> Provider -> provider-scoped resources
```

Provider is not a parent of Project. Customer use of provider resources occurs through authorized references, policy, and placement decisions.

#### B. Scope matrix

| Scope kind | Parent | Primary purpose | Typical resources |
|---|---|---|---|
| `Platform` | None | Global governance and catalog | PluginDefinition, global policy/definitions |
| `Organization` | Platform | Customer/MSP administrative boundary | Organization policy, shared catalog |
| `OrganizationUnit` | Organization or OrganizationUnit | Delegation and grouping | Delegated assignments and controls |
| `Tenant` | Organization or OrganizationUnit | Isolation and entitlement | Tenant policies, quotas |
| `Project` | Tenant | Workload/service ownership | ServiceInstance, PlacementRequest |
| `Provider` | Platform or Organization | Supply and infrastructure boundary | ResourcePool, AdapterConfiguration |
| Resource-local ownership | Parent resource | Lifecycle containment only | Child binding, component, operation |

Rules:

1. A persistent resource is either platform-scoped or has exactly one immutable primary `scopeRef`.
2. Each kind declares allowed scope kinds.
3. Identity uniqueness is `API group + kind + scope UID + name`.
4. Scope movement normally requires recreate-and-migrate.
5. Authorization resolves scope by UID, not display name.
6. `ownerRef` expresses lifecycle containment and MUST NOT replace security scope.
7. `scopeRef`, `ownerRef`, `location`, `sourceRef`, and `subjectRef` are distinct concepts.
8. Cross-tenant and cross-organization references are denied by default.
9. Denied cross-scope resolution MUST NOT disclose whether an inaccessible target exists.
10. Provider/resource-pool selection by customer-facing requests MUST occur through approved domain contracts rather than provider-native identifiers.

#### Operation allowed scopes (ADH-2026-013)

The canonical generic Operation contract declares exactly:

```yaml
x-sovrunn-allowed-scopes:
  - Platform
  - Organization
  - OrganizationUnit
  - Tenant
  - Project
  - Provider
```

Operation scope invariants:

1. Operation.scopeRef MUST equal the resolved canonical governance scope of Operation.targetRef.
2. For a platform-scoped target, Operation.scopeRef is canonically nil.
3. For a non-platform target, Operation.scopeRef identifies the target's governance scope by UID.
4. Operation.ownerRef MAY represent lifecycle containment but MUST NOT replace scopeRef or act as a governance or security scope.
5. A target/scope mismatch is rejected with a stable validation code and an RFC 6901 JSON Pointer path.
6. The six-value allowed-scope list does not grant authorization.
7. Target-kind constraints, caller authorization, and no-existence-disclosure rules remain mandatory.

### 6.6 References

References use a common typed base and domain-specific constrained aliases.

```yaml
resourcePoolRef:
  apiVersion: fabric.sovrunn.io/v1alpha1
  kind: ResourcePool
  name: sovereign-pool-a
  uid: <optional-immutable-id>
```

Rules:

- Singular fields end in `Ref`; collections end in `Refs`.
- `apiVersion`, `kind`, and name identify the requested target.
- UID MAY be omitted in human-authored input and returned after resolution.
- When name and UID are both present, they MUST identify the same object.
- Each reference schema constrains allowed kinds, scopes, and directions.
- Provider-native IDs MUST NOT act as core references.
- Secret values MUST be represented only through typed secret references.
- Generic references MAY provide a shared implementation base; public schemas SHOULD expose domain-specific reference types.

### 6.7 API boundaries

#### C1. Boundary matrix

| Boundary | Consumers | Allowed | Prohibited |
|---|---|---|---|
| `customer-facing` | Portal, CLI, SDK, GitOps, tenant automation | Product intent, safe status, actionable errors | Provider internals, secrets, inaccessible tenant data |
| `operator-facing` | Platform, MSP, provider operators | Administrative state, normalized infrastructure, diagnostics | Raw secrets and unrestricted customer data |
| `internal-engine-facing` | Policy, placement, entitlement, orchestration | Normalized internal contracts and decisions | Vendor SDK types as shared models |
| `adapter-facing` | External-system adapters | Translation contracts, provider handles, provenance | Leakage into customer/core schemas |
| `plugin-facing` | Plugin manager and implementations | Capabilities, operation contracts, validated results | Unrestricted control-plane access or policy bypass |
| `governance-only` | Architecture and review workflow | Decisions, assessments, approvals, traceability | Runtime credentials and customer secrets |

AI is a consumer of an authorized boundary, not a privileged boundary.

A single schema MUST NOT serve multiple audiences when their trust, sensitivity, or compatibility requirements differ. Separate boundary-specific views are preferred over hidden-field redaction.

### 6.8 Ownership and mutability

#### C2. Ownership matrix

| Area | Authoritative owner | Rule |
|---|---|---|
| Type metadata and schema | Sovrunn | Immutable for an object version. |
| Identity and scope | Sovrunn after validated creation | Immutable by default. |
| Desired `spec` | Declared customer/operator/system actor | Mutability defined per resource. |
| `status` | Declared controller, adapter, or plugin | Normal clients cannot write status. |
| Each condition type | One registered producer | No competing writers. |
| Provider-native handle | Adapter integration layer | Adapter/operator visibility only. |
| Plugin result | Plugin through validated contract | Cannot redefine core status grammar. |
| Stable error code/problem type | Sovrunn | Versioned machine contract. |
| Human message | Producing component | Informational; clients must not parse it. |
| Namespaced extension | Registered extension owner | Schema, boundary, version, and size controlled. |
| Secret value | External secret system | Never stored in common metadata, status, errors, or audit messages. |

Every mutable field and condition type MUST have one authoritative writer.

### 6.9 Desired state, status, phase, and conditions

Where applicable:

```text
spec   = desired/requested state
status = observed/evaluated state
```

`phase` is an optional coarse lifecycle summary. Conditions express independent current facts.

```yaml
status:
  observedGeneration: 4
  phase: Ready
  conditions:
    - type: Valid
      status: "True"
      reason: ValidationSucceeded
      message: Resource is valid.
      observedGeneration: 4
      lastTransitionTime: <timestamp>
```

Condition rules:

- `status` is `True`, `False`, or `Unknown`.
- `type` and `reason` are stable machine-readable PascalCase identifiers.
- `message` is human-readable and not a machine contract.
- `observedGeneration` identifies the desired-state generation evaluated.
- `lastTransitionTime` changes only when condition status changes.
- Conditions represent current observations, not event history.
- Historical decisions and actions belong to FEATURE-0013 records.
- Resources defining both phase and conditions MUST define consistency rules.

### 6.10 Validation

Validation occurs in ordered layers:

```text
1. HTTP/content and size checks
2. Safe JSON/YAML decoding
3. Duplicate- and unknown-field rejection
4. Structural schema validation
5. Deterministic defaulting
6. Semantic validation
7. Reference, kind, and scope validation
8. Authorization and policy validation
9. Later-feature capability/runtime validation
```

Rules:

- Unknown fields and duplicate keys are rejected by default.
- Defaulting MUST be deterministic, documented, and versioned.
- Structural, semantic, reference, authorization, and policy failures remain distinguishable.
- Validation MUST be safe for offline schema checks where external state is not required.
- Validation failures MUST identify stable codes and JSON Pointer field paths.
- Error messages and status MUST redact unauthorized or sensitive details.
- Validation limits MUST be finite for object size, nesting, metadata, conditions, references, and violations. Exact initial limits are defined in design and treated as reviewed platform configuration.

### 6.11 Error contract

Errors use Problem Details with Sovrunn extensions.

```json
{
  "type": "urn:sovrunn:problem:validation-failed",
  "title": "Validation failed",
  "status": 422,
  "detail": "One or more fields are invalid.",
  "instance": "/request-path",
  "code": "VALIDATION_FAILED",
  "requestId": "opaque-request-id",
  "violations": [
    {
      "field": "/spec/storage/sizeGiB",
      "code": "OUT_OF_RANGE",
      "message": "Value must be greater than zero."
    }
  ]
}
```

Baseline mappings:

| Failure | HTTP status |
|---|---:|
| Malformed request | 400 |
| Authentication required | 401 |
| Authorization denied | 403 |
| Not found | 404 |
| Lifecycle or uniqueness conflict | 409 |
| Structurally/semantically invalid resource | 422 |
| Stale resource version | 412 |
| Unsupported media type | 415 |
| Internal failure | 500 |
| Temporary dependency failure | 503 |

Problem types, codes, and violation codes are stable contracts. Messages MAY evolve. Responses MUST NOT expose credentials, stack traces, raw provider errors, sensitive policy inputs, or inaccessible resource details.

### 6.12 Updates and concurrency

Initial normative operations are create, get, list, full replace, and delete.

- Full replacement uses an opaque `resourceVersion` represented through HTTP ETag semantics.
- Protected updates MUST use `If-Match` or an equivalent checked resource version.
- Stale writes fail without overwriting current state.
- PATCH behavior is deferred until a specific patch contract is approved.
- Status updates use a separately authorized path or internal interface.
- Mutating APIs SHOULD support request correlation and idempotency where replay is possible.

### 6.13 Lists and pagination

List responses use `ListEnvelope`:

```yaml
apiVersion: core.sovrunn.io/v1alpha1
kind: ProjectList
items: []
page:
  nextPageToken: <opaque-token>
```

Rules:

- page tokens are opaque;
- ordering is deterministic;
- page size is bounded;
- total count is optional;
- tokens MUST NOT expose database offsets, provider details, or authorization context;
- list filtering is limited to explicitly indexed fields.

### 6.14 Extensions

Extensions are an explicit escape path, not a bypass.

An extension MUST have:

- a namespaced owner;
- a versioned schema;
- a declared boundary and data classification;
- finite size and nesting;
- validation and compatibility rules;
- no secret values;
- no provider-native leakage into customer/core contracts;
- no core decision dependency unless formally recognized.

When an extension becomes required by multiple independent consumers or core decisions, it MUST be reviewed for promotion into a typed Sovrunn contract.

### 6.15 API evolution

| Maturity | Compatibility expectation |
|---|---|
| Alpha | Breaking changes allowed only with explicit migration notes and review. |
| Beta | Breaking changes exceptional; migration and coexistence required. |
| Stable | Backward compatibility mandatory within the major version. |

Change classification:

| Change | Default classification |
|---|---|
| Add optional field | Compatible, subject to default/unknown-field rules |
| Add required field | Breaking |
| Remove or rename field | Breaking |
| Change field meaning, owner, or mutability | Breaking |
| Narrow enum or validation range | Breaking |
| Add enum value | Compatibility review required |
| Change reference target kind/scope | Compatibility review required |
| Expose internal data publicly | Security and boundary review required |
| Add registered extension | Additive within its extension contract |

Storage representation, served API version, and generated language types MUST remain distinguishable.

## 7. Future compatibility scenarios

#### D. Architecture conformance scenarios

FEATURE-0012 MUST demonstrate that its grammar can represent the following without implementing their domain behavior.

| Scenario | Profile | Scope/boundary | Required proof |
|---|---|---|---|
| Customer creates Project | ManagedResource | Tenant / customer | Stable identity, scope, strict validation, spec/status. |
| Operator registers ResourcePool | ManagedResource | Provider / operator | Provider-neutral capabilities; no vendor fields in core. |
| Adapter discovers external database | ObservedExternalResource | Provider / adapter | Provenance, freshness, stale/deleted semantics. |
| Publisher releases plugin contract | VersionedDefinition | Platform / plugin | Immutable published version and compatibility metadata. |
| Operator installs plugin | ManagedResource | Platform or Provider / operator | Definition, installation, and execution remain separate. |
| Operator configures adapter | ManagedResource | Provider / adapter | Secret references and native configuration isolation. |
| Placement engine evaluates request | TransientRequestResult | Project context / internal | Typed request/result without forced persistence. |
| Decision becomes auditable | ImmutableRecord | Project / governance | Immutable subject/actor/input references for FEATURE-0013. |
| Future provisioning executes | LongRunningOperation | Platform, Organization, OrganizationUnit, Tenant, Project, Provider / plugin | Idempotency, progress, retry, cancellation, terminal result; Operation.scopeRef equals the resolved canonical governance scope of Operation.targetRef. |
| Portal lists large collections | ListEnvelope | Customer/operator | Bounded opaque pagination and deterministic ordering. |
| Provider disconnects | ObservedExternalResource | Provider / adapter | Current, stale, unknown, and absent remain distinct. |
| External object is recreated under same name | ObservedExternalResource | Provider / adapter | UID prevents stale-reference rebinding. |
| Cross-tenant reference is attempted | Any reference | Customer | Denial without target-existence disclosure. |
| New cloud provider is added | Adapter-facing | Provider | No customer/core schema change. |
| New data-service plugin is added | Definition + operation | Multiple | No core grammar redesign. |
| AI explains denial | Existing authorized view | Customer/internal | Stable codes and safe context, no message scraping. |
| Phase 1 resource is migrated | ManagedResource | Compatibility API | Explicit version/migration; no silent reinterpretation. |

Required contract fixtures:

```text
Project
ResourcePool
DiscoveredDatabase
PluginDefinition
AdapterConfiguration
PlacementEvaluationRequest
Operation
AuditEvent
```

Fixtures prove schema fit, boundary classification, allowed scopes, ownership, strict parsing, reference behavior, and absence of later-phase execution.

## 8. Evolution safety mechanism

Architecture cannot prove that change will never be needed. FEATURE-0012 instead requires visible boundaries, executable constraints, and reversible evolution.

### 8.1 Boundary ledger

Every trust, API, scope, provider, plugin, adapter, lifecycle, data, and compatibility boundary MUST record:

```text
purpose; owner; producers; consumers; allowed/prohibited data;
authorization; audit; observability; failure behavior; versioning;
replacement path; migration path; reassessment trigger.
```

### 8.2 Architecture fitness functions

The feature gate and conformance suite MUST check at least:

1. every external schema declares profile, boundary, stability, and allowed scopes;
2. no core/customer schema imports or embeds provider SDK/native types;
3. every mutable field and condition has one owner;
4. unknown and duplicate fields fail;
5. references constrain kinds and scopes;
6. cross-tenant access fails without existence disclosure;
7. raw secret-like values are prohibited from metadata/status/errors;
8. externally sourced observations include provenance and freshness;
9. published definitions are immutable;
10. schema compatibility detects breaking changes;
11. object, metadata, condition, violation, reference, and page sizes are bounded;
12. errors use stable codes and JSON Pointer paths;
13. generated artifacts match the canonical schema;
14. later-feature runtime behavior is absent;
15. exceptions require an approved architecture handoff.

### 8.3 Reassessment triggers

FEATURE-0012 MUST be reassessed before or when introducing:

- the first non-Kubernetes provider integration;
- the first real provider adapter;
- the first remotely executed or data-path plugin;
- the first disconnected/federated control plane;
- the first external discovery source;
- cross-organization sharing;
- stable API promotion;
- high-frequency status updates;
- regulated/classified workloads;
- an object that cannot select an approved profile;
- a consumer requiring a new privileged view;
- an extension used by core decisions or multiple consumers;
- a backward-incompatible migration.

## 9. Risk register

#### E. FEATURE-0012 risks and controls

| ID | Risk | Primary control | Detection / response |
|---|---|---|---|
| F12-R01 | Overfit to Phase 1 or Kubernetes | Profile matrix and provider-neutral rules | Future-scenario fixtures; ADH for exceptions. |
| F12-R02 | Universal resource shape misrepresents lifecycle | Eight explicit profiles | Schema profile required; reclassify through review. |
| F12-R03 | Scope, location, source, and ownership become conflated | Separate typed concepts | Scope/reference negative tests. |
| F12-R04 | Provider-native data leaks into core/customer APIs | Adapter-only native contracts | Import/schema lint; move fields behind mapping. |
| F12-R05 | Plugin and adapter semantics are conflated | Separate boundaries and future resource families | Architecture classification tests. |
| F12-R06 | Generic or name-only references cause access/staleness errors | Typed refs plus optional immutable UID | Scope/kind and name/UID mismatch tests. |
| F12-R07 | Multiple writers corrupt status | One owner per field and condition | Ownership registry and concurrency tests. |
| F12-R08 | Conditions become history or phase explodes | Current-fact conditions; optional phase | Bounded status tests; move history to FEATURE-0013. |
| F12-R09 | Stale external observations produce inaccurate decisions | Required provenance, observed time, freshness | Freshness tests; mark stale/unknown and fail safely. |
| F12-R10 | Secrets or restricted data leak through metadata/errors | Data classification, redaction, SecretRef-only | Security scanning and negative response tests. |
| F12-R11 | Extensions become shadow APIs | Registered namespaced schemas | Reject unregistered use; promote mature contracts. |
| F12-R12 | API evolution breaks clients/plugins | Maturity and compatibility policy | Schema-diff gate; new version and migration plan. |
| F12-R13 | Unbounded objects/status/lists prevent scale | Explicit finite limits and opaque pagination | Limit tests and operational metrics. |
| F12-R14 | Documentation-only standard drifts | Shared primitives and executable conformance | Feature gate blocks non-conforming later features. |
| F12-R15 | Phase 1 migration expands into a rewrite | Compatibility audit and explicit exceptions | Scope gate; split migrations into approved features. |
| F12-R16 | AI or privileged consumers bypass normal boundaries | AI consumes authorized filtered views only | Access/redaction tests and audit. |

Residual risk is accepted only when documented with owner, corrective path, and reassessment trigger.

## 10. Security, sovereignty, accuracy, and audit requirements

Every schema field crossing an API boundary MUST identify or inherit:

```text
data classification;
authorized writer;
authorized readers;
mutability;
retention;
redaction behavior;
residency implications;
audit requirement.
```

Supported classifications:

```text
Public
Customer-visible
Tenant-confidential
Operator-confidential
Internal
Sensitive
Secret-reference-only
```

Additional invariants:

- external observations MUST expose provenance and freshness;
- inaccessible objects MUST not be revealed through error differences;
- provider replacement MUST not require customer-contract changes;
- disconnected external dependencies MUST produce explicit degraded/unknown status, not fabricated accuracy;
- auditing MUST link actors, requests, operations, subjects, versions, and later decisions without storing history in current resource status;
- AI-generated explanations MUST depend on stable structured context and authorized views.

## 11. FEATURE-0012 implementation contract

The approved implementation SHOULD contain only:

1. normative architecture and API/resource standard documentation;
2. shared type-metadata, metadata, reference, condition, problem, and validation primitives;
3. machine-readable schema conventions and schema metadata;
4. strict decoding and validation helpers;
5. compatibility/conformance tooling;
6. representative contract fixtures for the eight required objects;
7. Phase 1 compatibility report with conforming behavior, explicit exceptions, and migration candidates;
8. feature-gate checks for FEATURE-0013 and later adoption;
9. unit and integration tests for positive, negative, boundary, security, and compatibility behavior.

The implementation MUST NOT create functional provider, plugin, policy, placement, audit, or provisioning services.

## 12. Kiro generation contract

After `ADH-2026-012` and `ADH-2026-013` receive human approval, Kiro MUST treat this document as the architecture baseline.

### Requirements must

- convert every normative decision into testable requirements;
- include the A–E matrices and fitness functions as acceptance criteria;
- distinguish shared grammar from later-feature domain semantics;
- preserve provider neutrality, scope isolation, ownership, redaction, and compatibility;
- define Phase 1 compatibility deliverables without requiring a rewrite;
- include explicit non-goals and reassessment triggers.

### Design must

- identify the canonical schema source and derivative generation flow;
- define package/module boundaries for metadata, references, conditions, problems, validation, and conformance;
- define strict decoding, validation ordering, error mapping, and concurrency behavior;
- define machine-readable profile/boundary/scope annotations;
- define bounded configuration values and testing strategy;
- show how existing Phase 1 contracts coexist during migration;
- avoid dependencies on provider SDKs or later Phase 2 feature implementations.

### Tasks must

- implement the approved shared primitives and validators;
- add required fixtures and positive/negative conformance tests;
- add schema compatibility and boundary-lint checks;
- produce the Phase 1 compatibility report;
- wire FEATURE-0012 checks into the feature gate;
- update traceability and review artifacts;
- stop with `ARCHITECTURE_DECISION_REQUIRED` if a missing semantic decision is encountered.

## 13. Deferred but constrained decisions

The following are deliberately deferred and require separate approval when activated:

- exact HTTP route migration from Phase 1;
- concrete UID generation algorithm;
- PATCH format and field-ownership protocol;
- watch/change-stream protocol;
- production storage and indexing implementation;
- exact platform limits and tenancy scale targets;
- plugin execution and compatibility negotiation details;
- adapter transport and runtime protocol;
- policy inheritance/evaluation semantics;
- DecisionObject and AuditEvent payloads;
- provider/resource-pool domain fields;
- production identity, secrets, workflow, and observability integrations.

Deferred decisions MUST conform to the contracts and boundaries established here or use an approved architecture change.

## 14. Approval record

The architecture owner approved this baseline after confirming that:

1. the resource, scope, boundary, compatibility, and risk matrices are acceptable;
2. Sovrunn-owned responsibilities are explicit;
3. no known unmitigated conflict exists with platform foundation principles;
4. future scenarios fit without provider/customer coupling or arbitrary exceptions;
5. all intentional boundaries have owners, controls, migration paths, and reassessment triggers;
6. implementation scope remains limited to shared standards and conformance foundation;
7. residual risks are documented and accepted.

Approved statement:

> The FEATURE-0012 architecture has no identified unmitigated conflict with Sovrunn foundation principles or evaluated growth scenarios. Its intentional boundaries are explicit, owned, versioned, observable, auditable, replaceable, and testable. Remaining uncertainties have documented migration paths and reassessment triggers. Approved as the baseline for Kiro requirements, design, and tasks.

### Amendment record — ADH-2026-013

- Handoff: ADH-2026-013
- Classification: Clarification
- Approved by: Sanjeev Kumar
- Date: 2026-07-23
- Scope: Resolves the canonical Operation allowed-scope enumeration (all six Matrix B scopes) and the Operation-to-target scope equality invariant. No new architecture decision; no runtime behavior added.
