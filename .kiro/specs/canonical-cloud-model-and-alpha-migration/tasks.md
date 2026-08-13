# FEATURE-0015 — Canonical Cloud Model Foundation: Tasks

## 1. Identity, stage, inputs, and execution rules

| Field | Value |
|-------|-------|
| Feature | FEATURE-0015 — Canonical Cloud Model Foundation |
| Spec | `.kiro/specs/canonical-cloud-model-and-alpha-migration/` |
| Stage | Tasks |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase | 2R |
| Architecture authority | `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`; `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md` |
| Approved requirements | `.kiro/specs/canonical-cloud-model-and-alpha-migration/requirements.md` |
| Approved design | `.kiro/specs/canonical-cloud-model-and-alpha-migration/design.md` |
| Controlling Decisions | DEC-0037, DEC-0041, DEC-0042, DEC-0054, DEC-0059 |
| Controlling Handoffs | ADH-2026-020, ADH-2026-024, ADH-2026-025, ADH-2026-037, ADH-2026-042, ADH-2026-043 (non-migration), ADH-2026-045, ADH-2026-046, ADH-2026-047, ADH-2026-048, ADH-2026-049, ADH-2026-050, ADH-2026-051, ADH-2026-052, ADH-2026-053, ADH-2026-054, ADH-2026-055, ADH-2026-056, ADH-2026-057 |
| Owned resources | CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, InfrastructureStack |
| Module path | `github.com/sanjeevksaini/sovrunn` (Go 1.22) |

### 1.1 Purpose and execution rule

This stage collectively decomposes the approved design into independently testable, vertically sliced implementation tasks and inseparable work packages. Each path enumerated under either `Writable paths` or `Tests` is an exact writable path for that task; no unlisted implementation or test path may be changed without updating the task plan. No task implements a CONTRACT_ONLY or EXCLUDED element. Dependency order is explicit; independent tasks are marked.

### 1.2 Execution rules

- Implement tasks in dependency order.
- Each task is independently testable.
- Do not implement future-feature behavior or EXCLUDED elements.
- Every task cites exact requirements, design sections/decisions, architecture decisions, and applicable risks.
- The final task verifies repository cleanliness and boundary validation.
- Before Cursor begins any implementation task, load and follow `AGENTS.md`, `docs/engineering/go-coding-guardrails.md`, `docs/engineering/go-observability-standard.md`, and `docs/architecture/observability-and-audit-baseline.md`. These inherited guardrails are consumed by reference and are not redefined by this task plan.

---

## 2. Dependency-ordered implementation tasks

There are exactly eleven implementation commit tasks: 1, 4, 6, 8, 9, 11, 13, 14, 15, 16, and 17; Task 18 is the final non-commit verification checkpoint. Work packages 2, 3, 5, 7, 10, and 12 are inseparable internal sections of their stated parent task: they have no separate task or commit boundary, and their paths, tests, and acceptance criteria are included in the parent task's one commit.

Tasks are presented in dependency order. Each task includes: exact writable paths, tests, verification commands, acceptance criteria, security/observability impact, exclusions, and commit message.

### Task 1: Canonical resource foundation (types, ScopeKind, and ISO dataset)

**Purpose:** Implement the seven owned resource types (CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, InfrastructureStack) with metadata/spec/status shapes per VS0-SCHEMA-008..014.

**Dependencies:** None (foundation).

**Requirements traceability:** REQ-F15-01 (CloudPlatform), REQ-F15-02 (CloudProvider), REQ-F15-03 (CloudProviderParticipation), REQ-F15-04 (topology chain).

**Design traceability:** Design §3.1 (canonical resource types), §4.1 (owned resource shapes), §8.1 (IMPLEMENT classification), DD-01 (dedicated canonical package).

**Architecture decisions:** DEC-0037 (CloudPlatform/CloudProvider split), DEC-0041 (topology), DEC-0054 (participation), DEC-0059 (canonical bootstrap); ADH-2026-045 (immutable ownerRegistration), ADH-2026-047 (closed create contract).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this task.

**Writable paths:**
- `internal/cloudmodel/model/cloud_platform.go`
- `internal/cloudmodel/model/cloud_provider.go`
- `internal/cloudmodel/model/cloud_provider_participation.go`
- `internal/cloudmodel/model/hosting_location.go`
- `internal/cloudmodel/model/datacenter.go`
- `internal/cloudmodel/model/fault_domain.go`
- `internal/cloudmodel/model/infrastructure_stack.go`
- `internal/cloudmodel/model/doc.go`

**Tests:**
- `internal/cloudmodel/model/cloud_platform_test.go`
- `internal/cloudmodel/model/cloud_provider_test.go`
- `internal/cloudmodel/model/cloud_provider_participation_test.go`
- `internal/cloudmodel/model/hosting_location_test.go`
- `internal/cloudmodel/model/datacenter_test.go`
- `internal/cloudmodel/model/fault_domain_test.go`
- `internal/cloudmodel/model/infrastructure_stack_test.go`

Test JSON marshaling/unmarshaling; verify metadata/spec/status separation; verify no alpha-type reuse.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/cloudmodel/model/...
```

**Acceptance criteria:**
- All seven owned kinds have Go types with JSON tags.
- Each type uses canonical `metadata` (ObjectMeta with ScopeRef), `spec` (kind-specific fields), `status` (api-server-owned fields).
- CloudPlatform has immutable `spec.ownerRegistration{legalName,registrationIdentifier,jurisdictionCode}`.
- CloudProvider has `spec.operatingMarkets[]`.
- CloudProviderParticipation has `spec.cloudPlatformRef`, `spec.cloudProviderRef`, `spec.environment`; `status.phase`, `status.platformSuspended`, `status.providerSuspended`, `status.requestExpiresAt`.
- Datacenter, FaultDomain, and InfrastructureStack have immutable parent references; HostingLocation has no parent reference.
- No alpha `Provider`-like types are created or imported.
- AC-F15-01, AC-F15-02, AC-F15-03 structural foundation.

**Security/observability impact:**
- Security: Types carry no credential values; references are UID-pinned identifiers only.
- Observability: Types are structured for JSON logging.

**Exclusions:**
- No ExecutionTarget types (FEATURE-0016).
- No ServiceRegion, ServicePlan, or catalog types (FEATURE-0022).
- No `spec.providerSelectionModes` or `spec.permittedHostingLocationRefs` (FEATURE-0021).
- No alpha type reuse from `internal/resources`.

**Commit message:**
```
feat(cloudmodel): add canonical cloud model resource types

Introduce CloudPlatform, CloudProvider, CloudProviderParticipation,
HostingLocation, Datacenter, FaultDomain, and InfrastructureStack
types with metadata/spec/status shapes per VS0-SCHEMA-008..014.

- CloudPlatform: immutable ownerRegistration (ADH-2026-045)
- CloudProvider: operatingMarkets[]
- CloudProviderParticipation: lifecycle with independent holds
- HostingLocation has no parent reference; Datacenter/FaultDomain/InfrastructureStack have immutable parent references
- Includes work package 2: canonical ScopeKind and retained-fixture compatibility.
- Includes work package 3: version-pinned assigned ISO-3166-1 alpha-2 dataset and strict validation.

Refs: FEATURE-0015 Task 1, REQ-F15-01..04, DEC-0037/0041/0054/0059
```

---

#### Work package 2 (within Task 1): Shared ScopeKind canonicalization

**Purpose:** Update `internal/apimeta.ScopeKind` to the canonical seven-value vocabulary (Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider); update retained alpha fixtures/tests to compile against canonical CloudProvider; reject retired Provider scope kind on canonical paths.

**External dependencies:** None beyond the parent Task 1 foundation.

**Internal sequencing:** Implement after Task 1’s canonical resource-type definitions; this is parent-task sequencing, not a dependency of Task 1 on itself.

**Requirements traceability:** REQ-F15-10 (seven-scope vocabulary).

**Design traceability:** Design DD-03 (canonicalize shared ScopeKind), §3.1 (compatibility implementation), §8.1 (IMPLEMENT/CONTRACT_ONLY hybrid).

**Architecture decisions:** DEC-0037 (seven canonical scope kinds); ADH-2026-020 (retire ambiguous Provider).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this work package.

**Included writable paths:**
- `internal/apimeta/scope.go` (add `CloudPlatform`, `CloudProvider` constants; update validation)
- `internal/apimeta/scope_test.go`
- `internal/apiconform/canonical_schemas_test.go`
- `internal/apiconform/feature0014_bindings.go`
- `internal/apiconform/feature0014_negative_test.go`
- `internal/apiconform/fitness_ref.go`
- `internal/apiconform/fixtures.go`
- `internal/apiconform/fixtures_valid_test.go`
- `internal/apivalid/authz_test.go`
- `internal/apivalid/operation_target_scope_property_test.go`
- `internal/apivalid/platform_scope_property_test.go`
- `internal/apivalid/safe_denial_property_test.go`
- `internal/decision/validate/scope_test.go`
- `internal/decision/validate/validate.go`
- `internal/resources/provider.go`
- `internal/resources/provider_test.go`
- `internal/resources/providerlocation.go`
- `internal/resources/providerlocation_test.go`
- `internal/resources/providerdatacenter.go`
- `internal/resources/providerdatacenter_test.go`
- `internal/resources/datacenterfailuredomain.go`
- `internal/resources/datacenterfailuredomain_test.go`
- `internal/resources/infrastructurestack.go`
- `internal/resources/infrastructurestack_test.go`
- `internal/validation/completeness_test.go`
- `internal/validation/datacenterfailuredomain.go`
- `internal/validation/datacenterfailuredomain_test.go`
- `internal/validation/infrastructurestack.go`
- `internal/validation/infrastructurestack_test.go`
- `internal/validation/provider_test.go`
- `internal/validation/providerdatacenter.go`
- `internal/validation/providerdatacenter_test.go`
- `internal/validation/providerlocation.go`
- `internal/validation/providerlocation_test.go`
- `internal/validation/topology.go`
- `internal/validation/topology_test.go`
- `internal/apischema/schema_support_property_test.go`
- `api/schemas/_common/scope-ref.json`
- `api/schemas/adapter-configuration.json`
- `api/schemas/audit-event.json`
- `api/schemas/baseline/_common/scope-ref.json`
- `api/schemas/baseline/adapter-configuration.json`
- `api/schemas/baseline/audit-event.json`
- `api/schemas/baseline/decision-profile.json`
- `api/schemas/baseline/decision-record.json`
- `api/schemas/baseline/discovered-database.json`
- `api/schemas/baseline/evaluation-result.json`
- `api/schemas/baseline/operation.json`
- `api/schemas/baseline/resource-pool.json`
- `api/schemas/datacenter-failure-domain.json`
- `api/schemas/decision-profile.json`
- `api/schemas/decision-record.json`
- `api/schemas/discovered-database.json`
- `api/schemas/evaluation-result.json`
- `api/schemas/infrastructure-stack.json`
- `api/schemas/operation.json`
- `api/schemas/provider-datacenter.json`
- `api/schemas/provider-location.json`
- `api/schemas/provider.json`
- `api/schemas/resource-pool.json`
- `tests/conformance/fixtures/adapter-configuration.json`
- `tests/conformance/fixtures/adapter-configuration.yaml`
- `tests/conformance/fixtures/datacenter-failure-domain.json`
- `tests/conformance/fixtures/decision/positive/F13-COMPAT-06.json`
- `tests/conformance/fixtures/decision/positive/F13-COMPAT-08.json`
- `tests/conformance/fixtures/decision/positive/F13-SCOPE-06.json`
- `tests/conformance/fixtures/decision/positive/_shared-bundle.json`
- `tests/conformance/fixtures/discovered-database.json`
- `tests/conformance/fixtures/discovered-database.yaml`
- `tests/conformance/fixtures/infrastructure-stack.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-08.a.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-08.b.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-08.c.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-08.d.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-08.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-11.a.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-11.b.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-11.c.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-11.d.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-11.e.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-11.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-12.a.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-12.b.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-12.c.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-12.d.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-12.e.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-12.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-16.a.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-16.b.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-16.c.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-16.d.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-16.e.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-16.f.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-16.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-23.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-24.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-25.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-26.a.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-26.b.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-26.c.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-26.d.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-26.e.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-26.f.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-26.g.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-26.h.json`
- `tests/conformance/fixtures/negative/decision/F13-CF-26.json`
- `tests/conformance/fixtures/negative/decision/F13-EVAL-03.json`
- `tests/conformance/fixtures/negative/decision/F13-EVAL-04.json`
- `tests/conformance/fixtures/negative/decision/F13-EVAL-05.json`
- `tests/conformance/fixtures/negative/decision/F13-EVAL-06.json`
- `tests/conformance/fixtures/negative/decision/F13-EVAL-07.json`
- `tests/conformance/fixtures/negative/decision/F13-SCOPE-07.json`
- `tests/conformance/fixtures/negative/decision/F13-SCOPE-09.json`
- `tests/conformance/fixtures/negative/decision/F13-SCOPE-10.json`
- `tests/conformance/fixtures/negative/decision/F13-SEC-01.json`
- `tests/conformance/fixtures/negative/decision/F13-SEC-02.json`
- `tests/conformance/fixtures/negative/decision/F13-SEC-03.json`
- `tests/conformance/fixtures/negative/decision/F13-SEC-04.json`
- `tests/conformance/fixtures/negative/decision/F13-SEC-05.json`
- `tests/conformance/fixtures/negative/decision/F13-SEC-06.json`
- `tests/conformance/fixtures/negative/decision/F13-SEC-07.json`
- `tests/conformance/fixtures/negative/decision/F13-SEC-08.json`
- `tests/conformance/fixtures/negative/decision/F13-TRUST-03.json`
- `tests/conformance/fixtures/negative/feature0014/geo-malformed-country.json`
- `tests/conformance/fixtures/negative/feature0014/geo-prefix-mismatch.json`
- `tests/conformance/fixtures/negative/feature0014/geo-unassigned-valid.json`
- `tests/conformance/fixtures/negative/feature0014/multiple-parent-datacenter.json`
- `tests/conformance/fixtures/negative/feature0014/multiple-parent-stack.json`
- `tests/conformance/fixtures/negative/feature0014/skipped-level-datacenter-parent-provider.json`
- `tests/conformance/fixtures/negative/feature0014/skipped-level-failure-domain-parent-location.json`
- `tests/conformance/fixtures/negative/feature0014/skipped-level-stack-parent-datacenter.json`
- `tests/conformance/fixtures/negative/feature0014/skipped-level-stack-parent-location.json`
- `tests/conformance/fixtures/negative/feature0014/status-history.json`
- `tests/conformance/fixtures/negative/feature0014/technology-non-ascii.json`
- `tests/conformance/fixtures/negative/feature0014/technology-over-length.json`
- `tests/conformance/fixtures/negative/feature0014/wrong-kind-parent-failure-domain.json`
- `tests/conformance/fixtures/negative/feature0014/zero-parent-datacenter.json`
- `tests/conformance/fixtures/negative/feature0014/zero-parent-stack.json`
- `tests/conformance/fixtures/operation-provider.json`
- `tests/conformance/fixtures/operation-provider.yaml`
- `tests/conformance/fixtures/provider-datacenter.json`
- `tests/conformance/fixtures/provider-location.json`
- `tests/conformance/fixtures/provider.json`
- `tests/conformance/fixtures/resource-pool.json`
- `tests/conformance/fixtures/resource-pool.yaml`
- `api/schemas/baseline/BASELINE_MANIFEST.json`
- `api/schemas/baseline/BASELINE_APPROVALS.json`

**Included tests:**
- `internal/apimeta/scope_test.go`: verify seven canonical values; verify rejected legacy values; verify CloudPlatform/CloudProvider acceptance.
- Updated retained alpha tests, including `datacenterfailuredomain_test.go`: verify compilation and fixture alignment; if the named file needs no edit after the ScopeKind change, its existing focused test must prove that fact.

**Included verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/apimeta/...
go test -v ./internal/resources/... # retained alpha compatibility
```

**Included acceptance criteria:**
- `internal/apimeta.ScopeKind` enum has exactly seven canonical values: Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider.
- Canonical validation accepts CloudPlatform and CloudProvider.
- Canonical validation rejects retired `Provider` scope kind.
- Retained alpha resource fixtures/tests compile and pass against canonical CloudProvider scope.
- No legacy `Provider` ScopeKind compatibility API is exposed.
- AC-F15-01 (scope subsets), AC-F15-07 (stale-concept anti-drift).

**Security/observability impact:**
- Security: Scope vocabulary alignment prevents ambiguous authorization boundaries.
- Observability: Logs/metrics use canonical scope names.

**Exclusions:**
- No new ScopeKind values beyond the canonical seven.
- No six-scope vocabulary compatibility surface.

**Parent commit inclusion:** The Task 1 commit includes this work package's ScopeKind and retained-fixture changes.

---

#### Work package 3 (within Task 1): Assigned ISO-3166-1 alpha-2 dataset and lookup

**Purpose:** Implement compiled-in, version-pinned assigned ISO-3166-1 alpha-2 dataset (as of 2026-08-12) for strict case-sensitive uppercase validation of jurisdictionCode, operatingMarkets[], and countryCode; no ISO-3166-2 dataset or input normalization.

**External dependencies:** None beyond the parent Task 1 foundation.

**Internal sequencing:** Implement after Task 1’s canonical resource-type definitions; this is parent-task sequencing, not a dependency of Task 1 on itself.

**Requirements traceability:** REQ-F15-22 (assigned-code semantics).

**Design traceability:** Design DD-10 (repository-owned dataset), §3.1 (isocodes), §5.1 (validation pipeline), §8.1 (IMPLEMENT).

**Architecture decisions:** ADH-2026-048 decision 1 (assigned-code semantics; no ISO-3166-2).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this work package.

**Included writable paths:**
- `internal/cloudmodel/isocodes/dataset.go` (assigned codes as of 2026-08-12)
- `internal/cloudmodel/isocodes/lookup.go` (membership function)
- `internal/cloudmodel/isocodes/doc.go`

**Included tests:**
- `internal/cloudmodel/isocodes/dataset_test.go`: verify known assigned uppercase codes (US, IN, etc.); verify unassigned uppercase codes (ZZ) and lowercase input (`us`) are rejected; verify no input normalization.

**Included verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/cloudmodel/isocodes/...
```

**Included acceptance criteria:**
- Fixed dataset contains assigned ISO-3166-1 alpha-2 codes as of 2026-08-12.
- Lookup function returns true for assigned codes, false for unassigned codes.
- Validation is case-sensitive: lowercase input is rejected rather than normalized.
- No network request or third-party runtime dependency.
- No ISO-3166-2 membership dataset introduced.
- AC-F15-18 (assigned-code rejection).

**Security/observability impact:**
- Security: Deterministic validation; no external dependency.
- Observability: Validation failures are structured logs.

**Exclusions:**
- No ISO-3166-2 membership dataset.
- No network-based code verification.

**Parent commit inclusion:** The Task 1 commit includes this work package's ISO dataset and lookup changes.

---

### Task 4: Deterministic validation and in-memory store

**Purpose:** Implement pure, deterministic validation functions for closed create-contract classification, scope-kind subset, immutable-field, reference/containment, assigned-code, and header/body-input checks.

**Dependencies:** Task 1 (types, canonical ScopeKind, and ISO dataset).

**Requirements traceability:** REQ-F15-01 (CloudPlatform validation), REQ-F15-02 (CloudProvider validation), REQ-F15-03 (participation validation), REQ-F15-04 (topology validation), REQ-F15-14 (scope-reference UID invariant), REQ-F15-19 (closed create contract), REQ-F15-20 (scope derivation), REQ-F15-21 (participation body), REQ-F15-22 (assigned ISO), REQ-F15-23 (malformed inputs).

**Design traceability:** Design §5.1 (per-route pipeline), §5.3 (determinism), §7.1 (validate tests), §8.1 (IMPLEMENT).

**Architecture decisions:** DEC-0037 (scope kinds), DEC-0041 (topology), DEC-0054 (participation); ADH-2026-047 (closed contract, scope derivation, immutability), ADH-2026-048 (ISO assigned codes, header/body rules).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this task.

**Writable paths:**
- `internal/cloudmodel/validate/classify.go` (closed create-contract classification)
- `internal/cloudmodel/validate/scope.go` (scope-kind subset; scope-reference UID invariant)
- `internal/cloudmodel/validate/fields.go` (immutable/system-owned/unknown field rejection)
- `internal/cloudmodel/validate/references.go` (reference/containment; topology graph)
- `internal/cloudmodel/validate/iso.go` (assigned-code validation; administrativeAreaCode syntax)
- `internal/cloudmodel/validate/headers.go` (Idempotency-Key/If-Match/action-body preconditions)
- `internal/cloudmodel/validate/schema.go` (complete VS0-SCHEMA-008..014 constraint validation)
- `internal/api/decode.go` (bounded two-phase JSON decode, duplicate-member detection, route-safe reference extraction)
- `internal/api/patch.go` (RFC 7396 staged merge-patch application)
- `internal/cloudmodel/validate/doc.go`

**Tests:**
- `internal/cloudmodel/validate/classify_test.go`: per-kind client-required/client-optional allowlists; server-owned/status/scope/unknown/deferred field rejection.
- `internal/cloudmodel/validate/scope_test.go`: seven-scope subset per kind; cloudPlatformRef.uid = scopeRef.uid.
- `internal/cloudmodel/validate/fields_test.go`: immutable identity/reference/ownerRegistration; PATCHable description; name immutability.
- `internal/cloudmodel/validate/references_test.go`: topology parent chain; same-provider UID preservation.
- `internal/cloudmodel/validate/iso_test.go`: assigned US/IN accepted; unassigned ZZ rejected; administrativeAreaCode prefix.
- `internal/cloudmodel/validate/headers_test.go`: missing/empty/malformed/over-length Idempotency-Key; If-Match on participation create; non-empty action body.
- `internal/cloudmodel/validate/schema_test.go`: every registered required field, field type/range/format/reference constraint, and per-kind request-contract constraint from VS0-SCHEMA-008..014.
- `internal/api/decode_test.go`: inherited bounded-body limit; malformed JSON; duplicate object member at every relevant nesting level; phase-one route-safe extraction; phase-two strict closed-contract decode.
- `internal/api/patch_test.go`: RFC 7396 staged clone; `null` removal only for optional mutable `spec.description`/`spec.displayName`; required `spec.operatingMarkets: null` rejected; post-merge schema/immutable/system-owned/deferred/unknown-field validation before publication.

Property tests where applicable.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/cloudmodel/validate/...
```

**Acceptance criteria:**
- All validation functions are pure (no I/O, no wall-clock reads except injected timestamps, no map-iteration-order dependence).
- Closed create-contract classification: only client-required and client-optional fields accepted per kind.
- Scope-kind subset: CloudPlatform/CloudProvider accept only Platform; participation only CloudPlatform; topology only CloudProvider.
- Immutable fields: identity, references, ownerRegistration, topology name.
- Assigned ISO: jurisdictionCode, operatingMarkets[], countryCode reject unassigned codes.
- Header/body preconditions: malformed Idempotency-Key, If-Match on participation create, non-empty action body detected.
- Complete schema validation: all VS0-SCHEMA-008..014 field, format, range, reference, mutability, and request-contract constraints are checked before publication.
- Two-phase decode: bounded malformed/duplicate JSON is rejected before strict classification; phase two performs closed-contract classification before digesting, storing, or replaying.
- PATCH: RFC 7396 merges into a staged clone; `null` and post-merge validation follow Design §6.1.1 exactly.
- AC-F15-01 (validation), AC-F15-12 (UID invariant), AC-F15-15 (closed contract), AC-F15-16 (scope derivation), AC-F15-17 (participation body), AC-F15-18 (ISO), AC-F15-19 (malformed inputs).

**Security/observability impact:**
- Security: Deterministic validation prevents injection; no external I/O.
- Observability: Validation failures are structured Problem Details.

**Exclusions:**
- No FEATURE-0021 `spec.providerSelectionModes` or `spec.permittedHostingLocationRefs` validation (rejected as unknown fields).
- No ExecutionTarget validation (FEATURE-0016).
- No ServiceRegion or catalog validation (FEATURE-0022).

**Commit message:**
```
feat(cloudmodel): add deterministic validation functions

Implement pure validation for closed create-contract classification,
scope-kind subset, immutable-field, reference/containment,
assigned-code, and header/body-input checks per ADH-2026-047/048.

- Closed contract: client-required/client-optional field allowlists
- Scope subsets: Platform/CloudPlatform/CloudProvider per kind
- Immutable: identity, references, ownerRegistration, topology name
- Assigned ISO: reject unassigned codes
- Header/body: Idempotency-Key/If-Match/action-body preconditions
- Includes work package 5: in-memory store, uniqueness, resourceVersion, lock, and staging primitives.

Refs: FEATURE-0015 Task 4, REQ-F15-01/02/03/04/14/19/20/21/22/23,
ADH-2026-047, ADH-2026-048
```

---

#### Work package 5 (within Task 4): In-memory store with name/pair uniqueness and resourceVersion

**Purpose:** Implement in-memory registry with name uniqueness in server-derived scope, non-terminal participation pair uniqueness, resourceVersion, and publication coordinator.

**External dependencies:** Task 1 (types).

**Internal sequencing:** Implement after the validation section earlier in parent Task 4; this is parent-task sequencing, not a dependency of Task 4 on itself.

**Requirements traceability:** REQ-F15-01 (CloudPlatform uniqueness), REQ-F15-02 (CloudProvider uniqueness), REQ-F15-03 (participation pair uniqueness), REQ-F15-15 (audit atomicity).

**Design traceability:** Design DD-05 (in-memory store), §3.1 (store), §5.1 step 14 (publication coordinator), §8.1 (IMPLEMENT).

**Architecture decisions:** DEC-0037 (scope), DEC-0054 (participation); ADH-2026-045 (canonical bootstrap), ADH-2026-047 (scope derivation), ADH-2026-052 (uniqueness proofs).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this work package.

**Included writable paths:**
- `internal/cloudmodel/store.go` (per-kind registries, name/pair uniqueness, resourceVersion, publication lock and staging primitives)
- `internal/cloudmodel/doc.go`

**Included tests:**
- `internal/cloudmodel/store_test.go`: name uniqueness within Platform-root scope for CloudPlatform/CloudProvider; name uniqueness within CloudProvider scope for topology; CloudProviderParticipation name uniqueness within its server-derived CloudPlatform scope; separate non-terminal participation pair uniqueness; resourceVersion increment on mutation; concurrent create-name and participation-name/pair races (`go test -race`).

**Included verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/cloudmodel/...
go test -race ./internal/cloudmodel/...
```

**Included acceptance criteria:**
- Per-kind in-memory registries for the seven owned kinds.
- CloudPlatform/CloudProvider: name uniqueness within server-derived Platform-root scope.
- Topology: name uniqueness within server-derived CloudProvider scope.
- CloudProviderParticipation: name uniqueness within its server-derived CloudPlatform scope.
- CloudProviderParticipation: exactly one non-terminal participation per (cloudPlatformUID, cloudProviderUID) pair.
- resourceVersion increments on every mutation.
- Store lock serializes under-lock recheck and exposes staged-outcome primitives; Task 8's audit coordinator owns AuditEvent append-before-publication orchestration.
- Concurrent same-name creates: exactly one winner; loser receives ALREADY_EXISTS.
- AC-F15-01 (uniqueness), AC-F15-02 (participation pair).

**Security/observability impact:**
- Security: In-memory only; no durable/external persistence dependency.
- Observability: Store operations are logged with request correlation.

**Exclusions:**
- No durable storage, Kubernetes CRDs, or database persistence.
- No cross-process replay guarantee or persistence.

**Parent commit inclusion:** The Task 4 commit includes this work package's store, uniqueness, resourceVersion, lock, and staging changes.

---

### Task 6: Participation lifecycle and deterministic expiry

**Purpose:** Implement VS0-STATE-001 transitions, independent suspension holds (platformSuspended, providerSuspended), derived effective phase, and lifecycle validation.

**Dependencies:** Task 1 (types) and Task 4 (validation/store lock and staging primitives).

**Requirements traceability:** REQ-F15-03 (lifecycle), REQ-F15-06 (independent holds), REQ-F15-07 (no mutable spec; action preconditions).

**Design traceability:** Design §4.5 (lifecycle representation), §5.1 step 13 (source-state validation), §7.1 (lifecycle tests), §8.1 (IMPLEMENT).

**Architecture decisions:** DEC-0054 (participation), ADH-2026-037 (participation boundary), ADH-2026-045 (lifecycle actions), ADH-2026-046 (independent holds).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this task.

**Writable paths:**
- `internal/cloudmodel/lifecycle.go` (VS0-STATE-001 transitions, hold logic, derived phase, source-state validation)

**Tests:**
- `internal/cloudmodel/lifecycle_test.go`: Pending→Active/Rejected/Withdrawn/Expired; Active↔Suspended via independent holds; Active/Suspended→Terminating→Terminated; decline-release→prior effective phase; suspend/resume denied in Pending/Terminating/terminal; clearing one hold never reactivates while the other remains true; zero/multiple matching suspend/resume grant rejection; invalid source-state transitions.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/cloudmodel/...
```

**Acceptance criteria:**
- All VS0-STATE-001 transitions implemented.
- Independent holds: each administrator may set/clear only its own hold.
- Effective Suspended when either hold true; Active only when both false.
- Clearing one hold does not reactivate while the other remains true.
- Suspend/resume denied with CONFLICT/VS0_PARTICIPATION_STATE_INVALID in Pending/Terminating/terminal phases.
- Invalid source-state action returns CONFLICT with no status write.
- AC-F15-02 (lifecycle), AC-F15-10 (independent holds), AC-F15-11 (lifecycle changes).

**Security/observability impact:**
- Security: Hold-changing actions are scoped to the acting party's grant.
- Observability: Lifecycle transitions are logged and audited.

**Exclusions:**
- No end-user visibility eligibility (FEATURE-0021).
- No customer/Organization enrollment logic (FEATURE-0021).

**Commit message:**
```
feat(cloudmodel): add participation lifecycle and independent holds

Implement VS0-STATE-001 transitions with independent suspension holds
(platformSuspended, providerSuspended) and derived effective phase
per DEC-0054, ADH-2026-037, ADH-2026-046.

- Pending→Active/Rejected/Withdrawn/Expired
- Active↔Suspended via independent holds
- Active/Suspended→Terminating→Terminated
- Clearing one hold never reactivates while other remains true
- Suspend/resume denied in Pending/Terminating/terminal
- Includes work package 7: deterministic scheduler and API-server lifecycle wiring.

Refs: FEATURE-0015 Task 6, REQ-F15-03/06/07, DEC-0054, ADH-2026-046
```

---

#### Work package 7 (within Task 6): Deterministic participation-expiry scheduler

**Purpose:** Implement internal scheduler component that invokes api-server expiry transition for Pending participations at requestExpiresAt; not a public route or controller.

**External dependencies:** Task 1 (types) and Task 4 (store lock and staging primitives).

**Internal sequencing:** Implement after the lifecycle core earlier in parent Task 6; this is not a dependency of Task 6 on itself.

**Requirements traceability:** REQ-F15-06 (expiry as internal transition).

**Design traceability:** Design §4.6 (deterministic scheduler), §5.1 (scheduler expiry), §8.1 (IMPLEMENT).

**Architecture decisions:** DEC-0054 (participation), ADH-2026-045 (expiry is api-server transition), ADH-2026-046 decision 1 (scheduler is system actor, not separate writer).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this work package.

**Included writable paths:**
- `internal/cloudmodel/scheduler.go` (expiry scheduler with injected clock/ticker, stable-order processing, concurrency guard, recheck, api-server expiry invocation)
- `internal/server/server.go` (start/stop scheduler with server lifecycle)
- `cmd/sovrunn-api/main.go` (wire graceful shutdown into server stop)

**Included tests:**
- `internal/cloudmodel/scheduler_test.go`: due-time computation from earliest Pending requestExpiresAt; signal on create/action; stable ascending UID order; concurrency guard prevents race; recheck before expiry; idempotent expiry; scheduler start/stop; no public route or controller.

**Included verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/cloudmodel/...
```

**Included acceptance criteria:**
- Scheduler wakes at earliest Pending requestExpiresAt; waits if none.
- Processes due candidates in stable ascending UID order.
- Acquires same concurrency guard as participation actions.
- Rechecks Pending, expired, and version before invoking api-server expiry.
- Invokes api-server expiry transition rather than writing status directly.
- Idempotent: already non-Pending participation is not transitioned again.
- Scheduler started during api-server startup; stopped on graceful shutdown.
- No public route, controller, or external persistence dependency.
- AC-F15-02 (participation lifecycle and expiry).

**Security/observability impact:**
- Security: Scheduler is an internal system actor with no external surface.
- Observability: Expiry transitions are logged and audited with scheduler execution ID.

**Exclusions:**
- No public expiry route or action.
- No separate participation controller authority.

**Parent commit inclusion:** The Task 6 commit includes this work package's scheduler and server start/stop wiring.

---

### Task 8: Idempotency and audit-publication coordination

**Purpose:** Implement idempotency namespace, digest, reservation states (InFlight/Completed/Aborted), waiter/abort/release, 24-hour/10,000-entry retention, and completed-replay semantics.

**Dependencies:** Task 1 (types), Task 4 (validation/store lock and staging primitives), and Task 6 (scheduler/server lifecycle composition).

**Requirements traceability:** REQ-F15-18 (idempotency semantics).

**Design traceability:** Design DD-08 (idempotency namespace and digest), §5.1 step 9 (conditional idempotency branch), §7.1 (idempotency tests), §8.1 (IMPLEMENT).

**Architecture decisions:** ADH-2026-045 (idempotency key), ADH-2026-050 (all-route coverage), ADH-2026-054 (replay isolation and recheck), ADH-2026-055 (retention), ADH-2026-056 (replay response).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this task.

**Writable paths:**
- `internal/cloudmodel/idempotency.go` (namespace, digest, reservation, waiter/abort/release, retention/eviction)
- `internal/server/server.go` (final lifecycle composition: stop scheduler, abort in-flight idempotency reservations, wake waiters, then complete server shutdown)
- `cmd/sovrunn-api/main.go` (process signal/shutdown wiring for the composed server lifecycle)
- `internal/cloudmodel/idempotency_test.go`
- `internal/server/server_test.go` (combined graceful shutdown: scheduler stops, in-flight reservation aborts, waiters wake, no further scheduler transition occurs, process/server completion follows)

**Tests:**
- `internal/cloudmodel/idempotency_test.go`: namespace = (principal, registered pattern, target UID when present, key); same-key/same-digest replay returns stored success; same-key/different-digest returns CONFLICT/VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH; changed If-Match on action = different digest; InFlight/Completed/Aborted states; waiter detachment on cancellation; abort wakes waiters; completed-only 24-hour/10,000-entry retention; earliest-expiry then lexical eviction; InFlight never evicted; fresh request correlation on replay; validation, uniqueness, lifecycle-source-state, stale-version, audit-failure, recovered-panic, and graceful-shutdown aborts; completed replay requires current auth/authz/safe-access; race tests (`go test -race`).

**Verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/cloudmodel/...
go test -race ./internal/cloudmodel/...
```

**Acceptance criteria:**
- Idempotency namespace: (principal, registered method/path, target UID when present, key).
- Canonical digest: SHA-256 of deterministic JSON (body, pattern, If-Match when action); excludes key, credentials, correlation, server-assigned fields.
- Reservation: first caller atomically reserves InFlight; concurrent same-key/same-digest waits; same-key/different-digest returns CONFLICT.
- Completed replay: returns stored success before current-version comparison; requires current auth/authz/safe-access recheck; revoked grant denies without stored-result disclosure.
- Only successful collection creates and successful participation actions complete replay records.
- Validation, uniqueness, lifecycle-source-state, stale-version, audit-append failure, panic, shutdown abort reservation and wake waiters.
- Graceful API-server shutdown invokes the coordinator's in-flight reservation abort, wakes waiting callers, and completes before process exit.
- Final server composition stops the scheduler before completing idempotency abort/waiter wake-up and process exit; tests prove this combined ordering.
- The 24-hour and 10,000-entry cap applies only to Completed records; earliest-expiry then lexical eviction; InFlight is never evicted.
- Fresh request correlation on every request; requestId, correlation, ETag, Location, entity/transport headers not replayed.
- PATCH never requires, looks up, reserves, completes, or aborts idempotency state.
- AC-F15-11 (idempotency on creates/actions), AC-F15-19 (changed-digest conflict).

**Security/observability impact:**
- Security: Replay requires current auth/authz/safe-access recheck before disclosure.
- Observability: Idempotency events log only a redacted correlation-safe namespace fingerprint and outcome; never log raw Idempotency-Key, credentials, request body, or full namespace.

**Exclusions:**
- No cross-process replay guarantee or durable idempotency store.
- No GET, LIST, or PATCH idempotency state.

**Commit message:**
```
feat(cloudmodel): add idempotency coordinator with replay semantics

Implement idempotency namespace, digest, reservation states
(InFlight/Completed/Aborted), waiter/abort/release, 24-hour/10,000-
entry retention, and completed-replay per ADH-2026-050/054/055/056.

- Namespace: (principal, pattern, target UID, key)
- Same-key/same-digest: replay stored success
- Same-key/different-digest: CONFLICT/VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH
- Completed replay requires current auth/authz/safe-access recheck
- Validation/stale-version/audit-failure aborts; panic/shutdown aborts
- 24-hour/10,000-entry retention; InFlight never evicted
- Fresh request correlation on every request
- PATCH has no idempotency state
- Includes work package 10: AuditEvent append-before-publication orchestration.

Refs: FEATURE-0015 Task 8, REQ-F15-18, ADH-2026-050, ADH-2026-054,
ADH-2026-055, ADH-2026-056
```

---

#### Work package 10 (within Task 8): Audit/mutation atomic coordinator and append-failure mapping

**Purpose:** Implement the audit-publication orchestration that uses Task 4 store lock/staging primitives to append the required AuditEvent before staged resource/lifecycle/idempotency publication; append-failure → unpublished + INTERNAL_ERROR.

**External dependencies:** Task 4 (store lock and staging primitives).

**Internal sequencing:** Implement after the idempotency coordinator section earlier in parent Task 8; this is parent-task sequencing, not a dependency of Task 8 on itself.

**Requirements traceability:** REQ-F15-15 (audit atomicity and append-failure outcome).

**Design traceability:** Design DD-09 (audit/mutation coordinator), §5.1 step 14 (publication coordinator), §6.2 (privacy/redaction), §8.1 (IMPLEMENT).

**Architecture decisions:** ADH-2026-043 (preserved audit contract), ADH-2026-049 (INTERNAL_ERROR, not DEPENDENCY_UNAVAILABLE), ADH-2026-051 (audited-denial boundary).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this work package.

**Included writable paths:**
- `internal/cloudmodel/audit.go` (audit-publication orchestration using Task 4 under-lock recheck/staging primitives, AuditEvent append, staged-outcome publication)
- `internal/cloudmodel/audit_test.go`

**Included tests:**
- `internal/cloudmodel/audit_test.go`: audited success (collection create/PATCH, participation create/action/expiry) appends AuditEvent before publication; audited denial (status write, system-owned-field write, forged-grant, root denial, LIST-without-read-grant, safe-denial) appends AuditEvent before response; audit-append failure suppresses staged success/denial and substitutes INTERNAL_ERROR/500; no DEPENDENCY_UNAVAILABLE for this outcome; audit-append failure leaves neither mutation, lifecycle/status transition, nor idempotency completion published; an already-appended AuditEvent may remain after a later in-process publication failure; redacted AuditEvent has request correlation and no secret or inaccessible-resource disclosure.

**Included verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/cloudmodel/...
```

**Included acceptance criteria:**
- Task 4 store primitives serialize under-lock recheck (create name/pair uniqueness) and staging; this work package owns the AuditEvent append-before-publication orchestration.
- Required AuditEvent appended before staged resource/lifecycle/idempotency publication.
- Audited success: collection create/PATCH, participation create/action/expiry.
- Audited denial: status write, system-owned-field write, forged-grant, root denial, LIST-without-read-grant, safe-denial.
- Audit-append failure: staged success/denial suppressed; INTERNAL_ERROR/500 substituted; no mutation/lifecycle/status/idempotency published.
- No DEPENDENCY_UNAVAILABLE (503) for audit-append failure (no durable/external dependency).
- Redacted AuditEvent: request correlation; no secret, credential, inaccessible-resource disclosure, raw Idempotency-Key, request body, or complete idempotency namespace.
- AC-F15-13 (audit evidence). The append-failure outcome is REQ-F15-15 and VS0-CF-F15-24; AC-F15-15 is the separate closed-field-boundary criterion.

**Security/observability impact:**
- Security: AuditEvent redaction prevents secret/credential disclosure.
- Observability: One AuditEvent per audited outcome; audit-append failure is logged as INTERNAL_ERROR.

**Exclusions:**
- No durable or external AuditEvent store (FEATURE-0013 writer authority).
- No cross-process or persistent audit guarantee beyond FEATURE-0013 scope.

**Parent commit inclusion:** The Task 8 commit includes this work package's AuditEvent append-before-publication orchestration.

---

### Task 9: Bootstrap-grant resolver (server-resolved only)

**Purpose:** Implement server-resolved, deterministic bootstrap grant resolver; header/body grant claims ignored; no persisted roles/memberships/assignments.

**Dependencies:** None (independent).

**Requirements traceability:** REQ-F15-17 (bootstrap authorization).

**Design traceability:** Design DD-07 (bootstrap-grant resolver), §5.1 steps 3–5 (coarse action-grant lookup and forged-grant denial), §6.1 (authorization), §8.1 (IMPLEMENT).

**Architecture decisions:** ADH-2026-045 (bootstrap grants), ADH-2026-054 (forged-grant denial precedence).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this task.

**Writable paths:**
- `internal/cloudmodel/grants.go` (server-resolved grant lookup; header/body claim rejection)

**Tests:**
- `internal/cloudmodel/grants_test.go`: each exact server-resolved action and scope listed below is returned only for its allowed resource; X-Sovrunn-Bootstrap-Grant header claim ignored and denied; valid duplicate-free top-level bootstrapGrant body claim ignored and denied; zero/multiple suspend/resume matches fail closed; no persisted roles/memberships/assignments.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/cloudmodel/...
```

**Acceptance criteria:**
- Server-resolved, deterministic bootstrap grants only.
- Grant claims in X-Sovrunn-Bootstrap-Grant header or valid duplicate-free top-level bootstrapGrant body member are ignored and denied with audited AUTHORIZATION_DENIED.
- Bootstrap principal receives only these server-resolved grants: `cloudplatform.write` and `cloudprovider.write` at Platform-root scope; `topology.write` at CloudProvider UID scope; `participation.request.provider` at referenced CloudProvider UID plus CloudPlatform relation; `participation.accept.platform`, `participation.reject.platform`, and `participation.request-release.platform` at CloudPlatform UID scope; `participation.withdraw.provider`, `participation.accept-release.provider`, and `participation.decline-release.provider` at CloudProvider UID scope; and exactly one matching platform/provider suspend or resume action grant at its respective CloudPlatform/CloudProvider UID scope. Read grants remain `cloudplatform.read`, `cloudprovider.read`, `topology.read`, and `participation.read` at their exact design §6.1 scopes. Zero or multiple suspend/resume grant matches fail closed.
- No persisted roles, memberships, or assignments (FEATURE-0018 authority).
- REQ-F15-17 bootstrap boundary is proven by VS0-CF-F15-22; AC-F15-19 covers the valid duplicate-free forged-grant carrier's audited denial.

**Security/observability impact:**
- Security: Grants never accepted from header or body; deterministic server-resolved only.
- Observability: Forged-grant attempts are audited with AUTHORIZATION_DENIED.

**Exclusions:**
- No FEATURE-0018 roles, memberships, role assignments, or approval policies.
- No Organization/Tenant membership or individual authorization logic.

**Commit message:**
```
feat(cloudmodel): add bootstrap-grant resolver (server-resolved only)

Implement server-resolved, deterministic bootstrap grant resolver;
header/body grant claims ignored and denied per ADH-2026-045/054.

- Server-resolved action/scope/UID-restricted grants only
- X-Sovrunn-Bootstrap-Grant header claim: ignored, audited denial
- Valid duplicate-free top-level bootstrapGrant claim: ignored, audited denial
- No persisted roles/memberships/assignments
- Exact bootstrap grants and scopes: cloudplatform.write/cloudprovider.write (Platform-root), topology.write (CloudProvider UID), request/accept/reject/release actions at their declared participation CloudProvider or CloudPlatform UID, exact-one suspend/resume grant resolution, and read grants at their design §6.1 scopes

Refs: FEATURE-0015 Task 9, REQ-F15-17, ADH-2026-045, ADH-2026-054
```

---

### Task 11: CloudPlatform and CloudProvider collection/item handlers

**Purpose:** Implement CloudPlatform collection (LIST/create) and item (GET/PATCH) handlers with per-route pipeline order, closed create contract, PATCH-only update, scope derivation, validation, writer enforcement, idempotency, and audit.

**Dependencies:** Task 1 (types), Task 4 (validation/store primitives), Task 8 (idempotency/audit coordination), Task 9 (grants).

**Requirements traceability:** REQ-F15-01 (CloudPlatform validation), REQ-F15-08 (PATCH-only), REQ-F15-11 (writer enforcement), REQ-F15-15 (audit atomicity), REQ-F15-17 (server-resolved grants), REQ-F15-18 (create idempotency), REQ-F15-19 (closed create), REQ-F15-20 (scope derivation), REQ-F15-23 (malformed/prohibited inputs), REQ-F15-24 (AUTH_REQUIRED).

**Design traceability:** Design §4.2 (route model), §4.3 (closed create contract), §5.1 (per-route pipeline), §6.1 (authorization), §7.1 (handler tests), §8.1 (IMPLEMENT).

**Architecture decisions:** DEC-0037 (CloudPlatform), ADH-2026-045 (ownerRegistration), ADH-2026-047 (closed contract, scope derivation), ADH-2026-056 (PATCH conditional update, reads).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this task.

**Writable paths:**
- `internal/api/cloudplatform_collection.go` (LIST, create handlers)
- `internal/api/cloudplatform_item.go` (GET, PATCH handlers)

**Tests:**
- `internal/api/cloudplatform_collection_test.go`: authenticated LIST with cloudplatform.read returns ascending metadata.uid; LIST without read grant returns 403 with exactly one AuditEvent; inaccessible GET returns safe 404 with exactly one redacted AuditEvent; authorized create with closed contract returns 201 + exact initial status; server-derived Platform-root scope; Idempotency-Key required; same-key/same-digest replay and same-key/different-digest CONFLICT create no additional AuditEvent; duplicate name ALREADY_EXISTS and malformed input create no AuditEvent; missing auth AUTH_REQUIRED creates no AuditEvent; client-supplied status/system-owned metadata is audited, while client-supplied scopeRef/unknown/deferred field is unaudited; successful create appends exactly one AuditEvent; audit-append failure returns INTERNAL_ERROR with no publication.
- `internal/api/cloudplatform_item_test.go`: authorized GET returns resource; authorized PATCH spec.description returns 200 + updated resource and exactly one AuditEvent; stale/missing/malformed If-Match returns `STALE_RESOURCE_VERSION` / 412 with no AuditEvent; immutable ownerRegistration/name PATCH returns `VALIDATION_FAILED` / 422 with `VS0_PATCH_IMMUTABLE_FIELD` and no AuditEvent; client status write returns audited `AUTHORIZATION_DENIED` / 403 with `VS0_STATUS_FIELD_WRITE`; wrong-administrator PATCH returns `AUTHORIZATION_DENIED` / 403 before semantic processing with no AuditEvent; forged bootstrap header is denied/audited `AUTHORIZATION_DENIED` / 403 before an otherwise unsupported PATCH media type could return 415; PUT/DELETE return 405 with no AuditEvent; PATCH never requires idempotency state.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/api/...
```

**Acceptance criteria:**
- CloudPlatform collection: GET /apis/core.sovrunn.io/v1alpha1/cloud-platforms (LIST); POST (create).
- CloudPlatform item: GET /apis/core.sovrunn.io/v1alpha1/cloud-platforms/{uid}; PATCH (spec.description only).
- Per-route pipeline order: auth → route/method gate → headers → media/decode → coarse action-grant → root gate → safe resolution/scope derivation/exact authz → strict decode/classify → idempotency (create only) → scope derivation (create) → validation → version comparison (PATCH) → publication coordinator.
- Closed create contract: only metadata.name, spec.ownerRegistration.{legalName, registrationIdentifier, jurisdictionCode}, optional metadata.displayName, optional spec.description accepted.
- Server-derived Platform-root scope; no client-supplied scope.
- PATCH-only; PUT/DELETE return 405; PATCH requires application/merge-patch+json and If-Match.
- Writer enforcement: only cloud-platform-admin may PATCH spec; wrong-administrator denied AUTHORIZATION_DENIED before semantic processing.
- Idempotency: create requires key; PATCH has no idempotency state.
- Audit: successful create/PATCH and audited denials append AuditEvent; append-failure returns INTERNAL_ERROR.
- AC-F15-01 (CloudPlatform collection create, GET/LIST, permitted PATCH), AC-F15-06 (PATCH-only), AC-F15-08 (writer enforcement), AC-F15-13 (audit evidence), AC-F15-15 (closed contract), AC-F15-16 (scope derivation), AC-F15-20 (AUTH_REQUIRED).

**Security/observability impact:**
- Security: Authentication required; writer-boundary enforcement; safe denial; no secret disclosure.
- Observability: Request correlation, structured logs, AuditEvent correlation.

**Exclusions:**
- No PUT or DELETE (405).
- No mutable fields beyond spec.description.
- No client-supplied scope or status.

**Commit message:**
```
feat(api): add CloudPlatform collection and item handlers

Implement CloudPlatform LIST/create/GET/PATCH handlers with per-route
pipeline, closed create contract, PATCH-only update, scope derivation,
writer enforcement, idempotency (create only), and audit atomicity
per ADH-2026-045/047/056.

- Collection: LIST, create with Platform-root scope derivation
- Item: GET, PATCH (spec.description only)
- PUT/DELETE return 405; PATCH requires If-Match
- Writer enforcement: cloud-platform-admin only
- Idempotency: create requires key; PATCH has no state
- Audit: successful create/PATCH, audited denials; append-failure → INTERNAL_ERROR
- Includes work package 12: CloudProvider collection/item handler changes and tests.

Refs: FEATURE-0015 Task 11, REQ-F15-01/08/11/19/20/24, DEC-0037,
ADH-2026-047, ADH-2026-056
```

---

#### Work package 12 (within Task 11): CloudProvider collection and item handlers

**Purpose:** Implement CloudProvider collection (LIST/create) and item (GET/PATCH) handlers with per-route pipeline, closed create contract, PATCH-only update (displayName, operatingMarkets), assigned-ISO validation, scope derivation, writer enforcement, idempotency, and audit.

**Dependencies:** Task 1 (types and ISO dataset), Task 4 (validation/store primitives), Task 8 (idempotency/audit coordination), Task 9 (grants).

**Requirements traceability:** REQ-F15-02 (CloudProvider validation), REQ-F15-08 (PATCH-only), REQ-F15-11 (writer enforcement), REQ-F15-15 (audit atomicity), REQ-F15-17 (server-resolved grants), REQ-F15-18 (create idempotency), REQ-F15-19 (closed create), REQ-F15-20 (scope derivation), REQ-F15-22 (assigned ISO), REQ-F15-23 (malformed/prohibited inputs), REQ-F15-24 (AUTH_REQUIRED).

**Design traceability:** Design §4.2 (route model), §4.3 (closed create contract), §5.1 (per-route pipeline), §6.1 (authorization), §7.1 (handler tests), §8.1 (IMPLEMENT).

**Architecture decisions:** DEC-0037 (CloudProvider), ADH-2026-047 (closed contract, scope derivation), ADH-2026-048 (assigned ISO), ADH-2026-056 (PATCH conditional update, reads).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this work package.

**Included writable paths:**
- `internal/api/cloudprovider_collection.go` (LIST, create handlers)
- `internal/api/cloudprovider_item.go` (GET, PATCH handlers)

**Included tests:**
- `internal/api/cloudprovider_collection_test.go`: authenticated LIST with cloudprovider.read; LIST without read grant 403 with exactly one AuditEvent; authorized create with closed contract returns 201 + exact initial status; server-derived Platform-root scope; root gate absent CloudPlatform returns audited CONFLICT/VS0_CLOUDPLATFORM_ROOT_REQUIRED and audit-append failure substitutes INTERNAL_ERROR with no publication; non-empty operatingMarkets with assigned codes (US, IN) accepted; unassigned code (ZZ), duplicate name ALREADY_EXISTS, malformed input, replay/mismatch, and missing auth have no AuditEvent; client-supplied status/system-owned metadata is audited, while client-supplied scopeRef/unknown/deferred fields are unaudited; successful create appends exactly one AuditEvent; audit-append failure INTERNAL_ERROR with no publication.
- `internal/api/cloudprovider_item_test.go`: authorized GET returns resource; authorized PATCH spec.displayName or spec.operatingMarkets returns 200 + updated with exactly one AuditEvent; stale/missing/malformed If-Match returns `STALE_RESOURCE_VERSION` / 412 with no AuditEvent; immutable name PATCH returns `VALIDATION_FAILED` / 422 with `VS0_PATCH_IMMUTABLE_FIELD` and no AuditEvent; client status write returns audited `AUTHORIZATION_DENIED` / 403 with `VS0_STATUS_FIELD_WRITE`; wrong-administrator PATCH returns `AUTHORIZATION_DENIED` / 403 before semantic processing with no AuditEvent; forged bootstrap header is denied/audited `AUTHORIZATION_DENIED` / 403 before unsupported PATCH media is evaluated; PUT/DELETE 405 and PATCH idempotency absence create no AuditEvent.

**Included verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/api/...
```

**Included acceptance criteria:**
- CloudProvider collection: GET /apis/core.sovrunn.io/v1alpha1/cloud-providers (LIST); POST (create).
- CloudProvider item: GET /apis/core.sovrunn.io/v1alpha1/cloud-providers/{uid}; PATCH (spec.displayName, spec.operatingMarkets).
- Closed create contract: metadata.name, non-empty spec.operatingMarkets[], optional metadata.displayName, optional spec.displayName.
- Server-derived Platform-root scope.
- Root gate: CloudProvider create is denied until CloudPlatform exists with audited `CONFLICT` / `VS0_CLOUDPLATFORM_ROOT_REQUIRED`; AuditEvent append failure returns `INTERNAL_ERROR` and publishes neither denial nor mutation.
- Assigned ISO: operatingMarkets[] accepts assigned codes; rejects unassigned codes (ZZ) with VALIDATION_FAILED.
- PATCH-only; PUT/DELETE 405; PATCH requires application/merge-patch+json and If-Match.
- Writer enforcement: cloud-provider-admin only; wrong-administrator denied before semantic processing.
- Idempotency: create requires key; PATCH has no state.
- Audit: successful create/PATCH and only the registered audited denials append exactly one AuditEvent; all other validation/version/uniqueness/replay denials append none; append-failure INTERNAL_ERROR with no publication.
- AC-F15-01 (CloudProvider collection create, GET/LIST, permitted PATCH), AC-F15-06 (PATCH-only), AC-F15-08 (writer enforcement), AC-F15-13 (audit evidence), AC-F15-15 (closed contract), AC-F15-16 (scope derivation), AC-F15-18 (assigned ISO), AC-F15-20 (AUTH_REQUIRED).

**Security/observability impact:**
- Security: Authentication required; writer-boundary enforcement; assigned-ISO validation.
- Observability: Request correlation, structured logs, AuditEvent correlation.

**Exclusions:**
- No PUT or DELETE.
- No client-supplied scope or status.

**Parent commit inclusion:** The Task 11 commit includes this work package's CloudProvider collection/item handler changes and tests.

---

### Task 13: HTTP handlers for topology resources (HostingLocation, Datacenter, FaultDomain, InfrastructureStack)

**Purpose:** Implement collection (LIST/create) and item (GET/PATCH) handlers for the four topology resources with per-route pipeline, closed create contract, HostingLocation scope derived from topology.write with no parent reference, immutable Datacenter/FaultDomain/InfrastructureStack parent references, assigned-ISO validation (HostingLocation), PATCH-only update (description only), writer enforcement, idempotency, and audit.

**Dependencies:** Task 1 (types and ISO dataset), Task 4 (validation/store primitives), Task 8 (idempotency/audit coordination), Task 9 (grants), Task 11 (CloudPlatform and CloudProvider handlers ensure root exists).

**Requirements traceability:** REQ-F15-04 (topology validation), REQ-F15-08 (PATCH-only), REQ-F15-09 (safe cross-provider denial), REQ-F15-11 (writer enforcement), REQ-F15-15 (audit atomicity), REQ-F15-16 (root gate), REQ-F15-17 (server-resolved grants), REQ-F15-18 (create idempotency), REQ-F15-19 (closed create), REQ-F15-20 (scope derivation), REQ-F15-22 (assigned ISO for HostingLocation), REQ-F15-23 (malformed/prohibited inputs), REQ-F15-24 (AUTH_REQUIRED).

**Design traceability:** Design §4.2 (route model), §4.3 (closed create contract), §4.4 (scope derivation), §5.1 (per-route pipeline step 6 root gate), §6.1 (authorization), §7.1 (handler tests), §8.1 (IMPLEMENT).

**Architecture decisions:** DEC-0041 (topology), ADH-2026-045 (root gate), ADH-2026-047 (closed contract, scope derivation, immutability), ADH-2026-048 (assigned ISO), ADH-2026-056 (PATCH conditional update, reads).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this task.

**Writable paths:**
- `internal/api/hostinglocation_collection.go`, `internal/api/hostinglocation_item.go`
- `internal/api/datacenter_collection.go`, `internal/api/datacenter_item.go`
- `internal/api/faultdomain_collection.go`, `internal/api/faultdomain_item.go`
- `internal/api/infrastructurestack_collection.go`, `internal/api/infrastructurestack_item.go`

**Tests:**
- `internal/api/hostinglocation_collection_test.go`, `internal/api/datacenter_collection_test.go`, `internal/api/faultdomain_collection_test.go`, `internal/api/infrastructurestack_collection_test.go`: authenticated LIST with topology.read returns ascending metadata.uid; LIST without read grant returns 403 with exactly one AuditEvent; authorized create with closed contract returns 201 + exact initial status and one AuditEvent; HostingLocation: server-derived CloudProvider scope from topology.write grant, assigned countryCode (US accepted, ZZ rejected), optional administrativeAreaCode (syntax-plus-prefix only); Datacenter/FaultDomain/InfrastructureStack: server-derived CloudProvider scope from resolved immutable parent reference; root gate and inaccessible reference produce exactly one required redacted AuditEvent, with append failure substituting INTERNAL_ERROR/no publication; duplicate name, malformed/unknown/scope input, unassigned ISO, visible mismatch, idempotency replay/mismatch, and missing authentication produce no AuditEvent; client-supplied status/system-owned metadata is audited; successful create appends AuditEvent.
- `internal/api/hostinglocation_item_test.go`, `internal/api/datacenter_item_test.go`, `internal/api/faultdomain_item_test.go`, `internal/api/infrastructurestack_item_test.go`: authorized GET returns resource; authorized PATCH spec.description returns 200 + updated with exactly one AuditEvent; stale/missing/malformed If-Match returns `STALE_RESOURCE_VERSION` / 412 with no AuditEvent; immutable name and, for Datacenter/FaultDomain/InfrastructureStack only, immutable parent-reference PATCH return `VALIDATION_FAILED` / 422 with `VS0_PATCH_IMMUTABLE_FIELD` and no AuditEvent; client status write returns audited `AUTHORIZATION_DENIED` / 403 with `VS0_STATUS_FIELD_WRITE`; wrong-administrator PATCH returns `AUTHORIZATION_DENIED` / 403 before semantic processing with no AuditEvent; forged bootstrap header is denied/audited `AUTHORIZATION_DENIED` / 403 before unsupported PATCH media is evaluated; PUT/DELETE 405 and PATCH idempotency absence create no AuditEvent.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/api/...
```

**Acceptance criteria:**
- Four topology kinds: HostingLocation, Datacenter, FaultDomain, InfrastructureStack.
- Collection paths: GET /apis/infrastructure.sovrunn.io/v1alpha1/{kind} (LIST); POST (create).
- Item paths: GET /apis/infrastructure.sovrunn.io/v1alpha1/{kind}/{uid}; PATCH (spec.description only).
- HostingLocation closed create: metadata.name, spec.countryCode (assigned), spec.locality, optional spec.administrativeAreaCode (syntax-plus-prefix), optional spec.description.
- Datacenter/FaultDomain/InfrastructureStack closed create: metadata.name, immutable parent ref, optional spec.description.
- Root gate: create denied until CloudPlatform exists (CONFLICT/VS0_CLOUDPLATFORM_ROOT_REQUIRED).
- Scope derivation: HostingLocation from topology.write grant CloudProvider UID; descendants from resolved parent.
- Assigned ISO: countryCode accepts assigned, rejects unassigned (ZZ) with VALIDATION_FAILED.
- administrativeAreaCode: syntax-plus-country-prefix validation only; no ISO-3166-2 dataset.
- Immutable: name (identity) for all topology kinds; parent references only for Datacenter, FaultDomain, and InfrastructureStack. HostingLocation has no parent reference.
- PATCHable: spec.description only.
- PUT/DELETE 405; PATCH requires application/merge-patch+json and If-Match.
- Writer enforcement: cloud-provider-admin only; wrong-administrator denied before semantic processing.
- Cross-provider isolation: inaccessible reference safe 404; visible mismatch VALIDATION_FAILED/VS0_TOPOLOGY_PROVIDER_MISMATCH.
- Idempotency: create requires key; PATCH has no state.
- Audit: successful create/PATCH and only registered audited denials append exactly one AuditEvent; validation/version/uniqueness/replay denials append none; append-failure INTERNAL_ERROR with no publication.
- AC-F15-01 (topology collection create, GET/LIST, permitted PATCH), AC-F15-04 (cross-provider denial), AC-F15-06 (PATCH-only), AC-F15-08 (writer enforcement), AC-F15-13 (audit evidence), AC-F15-14 (root requirement), AC-F15-15 (closed contract), AC-F15-16 (scope derivation), AC-F15-18 (assigned ISO), AC-F15-20 (AUTH_REQUIRED).

**Security/observability impact:**
- Security: Authentication required; writer-boundary enforcement; safe cross-provider denial; assigned-ISO validation.
- Observability: Request correlation, structured logs, AuditEvent correlation.

**Exclusions:**
- No PUT or DELETE.
- No ISO-3166-2 membership dataset (administrativeAreaCode syntax-plus-prefix only).
- No client-supplied scope or status.

**Commit message:**
```
feat(api): add topology resource handlers (HostingLocation, Datacenter, FaultDomain, InfrastructureStack)

Implement LIST/create/GET/PATCH handlers for four topology resources
with per-route pipeline, closed create contract, immutable parent
references, single-source scope derivation, assigned-ISO validation
(HostingLocation), root gate, PATCH-only update (description),
writer enforcement, idempotency (create only), and audit per
ADH-2026-045/047/048/056.

- HostingLocation: scope from topology.write grant; assigned countryCode; administrativeAreaCode syntax-plus-prefix
- Datacenter/FaultDomain/InfrastructureStack: scope from resolved parent
- Root gate: create denied until CloudPlatform exists (CONFLICT/VS0_CLOUDPLATFORM_ROOT_REQUIRED)
- Immutable: name for all topology kinds; parent references only for Datacenter/FaultDomain/InfrastructureStack (HostingLocation has none)
- PATCHable: spec.description only
- PUT/DELETE 405; PATCH requires If-Match
- Writer enforcement: cloud-provider-admin only
- Cross-provider: inaccessible safe 404; visible mismatch VALIDATION_FAILED
- Idempotency: create requires key; PATCH has no state
- Audit: successful create/PATCH, audited denials; append-failure → INTERNAL_ERROR

Refs: FEATURE-0015 Task 13, REQ-F15-04/08/11/16/19/20/22/24,
DEC-0041, ADH-2026-045, ADH-2026-047, ADH-2026-048, ADH-2026-056
```

---

### Task 14: HTTP handlers for CloudProviderParticipation collection and item (GET-only)

**Purpose:** Implement CloudProviderParticipation collection (LIST/create) and item (GET-only; no PATCH) handlers with per-route pipeline, closed create contract (exactly four decision-4 fields), pair uniqueness, server-derived CloudPlatform scope, initial Pending status with false holds and seven-day expiry, root gate, Idempotency-Key-only precondition (no If-Match on create), idempotency, and audit.

**Dependencies:** Task 1 (types), Task 4 (validation/store primitives), Task 6 (lifecycle), Task 8 (idempotency/audit coordination), Task 9 (grants), Task 11 (CloudPlatform exists for scope derivation).

**Requirements traceability:** REQ-F15-03 (participation registration), REQ-F15-07 (create preconditions), REQ-F15-14 (scope-reference UID invariant), REQ-F15-15 (audit atomicity), REQ-F15-16 (root gate), REQ-F15-17 (server-resolved grants), REQ-F15-18 (create idempotency), REQ-F15-19 (closed create), REQ-F15-20 (scope derivation), REQ-F15-21 (participation create body), REQ-F15-23 (If-Match on create rejected), REQ-F15-24 (AUTH_REQUIRED).

**Design traceability:** Design §4.2 (route model, no PATCH registration for participation), §4.3 (closed create contract), §4.4 (scope derivation), §4.5 (lifecycle), §5.1 (per-route pipeline), §6.1 (authorization), §7.1 (handler tests), §8.1 (IMPLEMENT).

**Architecture decisions:** DEC-0054 (participation), ADH-2026-045 (root gate), ADH-2026-046 (create Idempotency-Key only; actions require both preconditions), ADH-2026-047 (closed contract decision 4), ADH-2026-048 (If-Match on create MALFORMED_REQUEST), ADH-2026-056 (reads).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this task.

**Writable paths:**
- `internal/api/participation_collection.go` (LIST, create handlers)
- `internal/api/participation_item.go` (GET-only handler; no PATCH registration)

**Tests:**
- `internal/api/participation_collection_test.go`: authenticated LIST with participation.read returns ascending metadata.uid; LIST without read grant 403 with exactly one AuditEvent; authorized create with exactly four fields (metadata.name, spec.cloudPlatformRef, spec.cloudProviderRef, spec.environment=development) returns 201 + Pending/false holds/seven-day expiry and exactly one AuditEvent; server-derived CloudPlatform scope from cloudPlatformRef; cloudPlatformRef.uid = scopeRef.uid invariant enforced; root gate and inaccessible reference produce exactly one required redacted AuditEvent, with append failure substituting INTERNAL_ERROR/no publication; mismatch, duplicate participation name or non-terminal pair, malformed body/headers, idempotency replay/mismatch, and missing authentication produce no AuditEvent; client-supplied status/system-owned metadata is audited, while client-supplied scopeRef/unknown/deferred fields are unaudited; FEATURE-0021 fields (providerSelectionModes, permittedHostingLocationRefs) rejected as UNKNOWN_FIELD.
- `internal/api/participation_item_test.go`: authorized GET returns resource; no PATCH registration (item lifecycle changes only through explicit action routes).

**Verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/api/...
```

**Acceptance criteria:**
- CloudProviderParticipation collection: GET /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations (LIST); POST (create).
- CloudProviderParticipation item: GET /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid} (no PATCH registration).
- Closed create contract: exactly metadata.name, spec.cloudPlatformRef, spec.cloudProviderRef, spec.environment (development only); no optional fields.
- FEATURE-0021 fields (providerSelectionModes, permittedHostingLocationRefs) rejected as UNKNOWN_FIELD.
- Server-derived CloudPlatform scope from resolved cloudPlatformRef.
- cloudPlatformRef.uid = scopeRef.uid invariant; mismatch VALIDATION_FAILED/VS0_SCOPE_REFERENCE_MISMATCH.
- Initial status: phase=Pending, platformSuspended=false, providerSuspended=false, requestExpiresAt=createdAt+7 days.
- Root gate: create denied until CloudPlatform exists (CONFLICT/VS0_CLOUDPLATFORM_ROOT_REQUIRED).
- Duplicate non-terminal (platform, provider) pair: ALREADY_EXISTS/VS0_PARTICIPATION_DUPLICATE.
- Create preconditions: Idempotency-Key required; If-Match on create MALFORMED_REQUEST.
- Pair uniqueness is keyed by the non-terminal `(cloudPlatformUID, cloudProviderUID)` pair. Collection-create idempotency is separately keyed by `(principal, registered method/path pattern, Idempotency-Key)` with the classified body in its digest; same-key/same-digest replay; same-key/different-digest CONFLICT.
- Audit: successful create and only registered audited denials append exactly one AuditEvent; validation/uniqueness/replay denials append none; append-failure INTERNAL_ERROR with no publication.
- AC-F15-02 (participation lifecycle), AC-F15-11 (no mutable spec; create Idempotency-Key only), AC-F15-12 (UID invariant), AC-F15-13 (audit evidence), AC-F15-14 (root requirement), AC-F15-15 (closed contract), AC-F15-16 (scope derivation), AC-F15-17 (participation create body), AC-F15-19 (If-Match on create MALFORMED_REQUEST), AC-F15-20 (AUTH_REQUIRED).

**Security/observability impact:**
- Security: Authentication required; server-derived scope; no client scope/status; pair uniqueness.
- Observability: Request correlation, structured logs, AuditEvent correlation.

**Exclusions:**
- No PATCH registration (lifecycle changes through explicit action routes only).
- No spec.providerSelectionModes or spec.permittedHostingLocationRefs (FEATURE-0021).
- No PUT or DELETE.

**Commit message:**
```
feat(api): add CloudProviderParticipation collection and item (GET-only)

Implement CloudProviderParticipation LIST/create/GET handlers with
per-route pipeline, closed create contract (exactly four decision-4
fields), pair uniqueness, server-derived CloudPlatform scope, initial
Pending status with false holds and seven-day expiry, root gate,
Idempotency-Key-only precondition (no If-Match on create), idempotency,
and audit per ADH-2026-045/046/047/048/056.

- Closed create: metadata.name, cloudPlatformRef, cloudProviderRef, environment (development)
- FEATURE-0021 fields (providerSelectionModes, permittedHostingLocationRefs) rejected
- Server-derived CloudPlatform scope from cloudPlatformRef
- cloudPlatformRef.uid = scopeRef.uid invariant enforced
- Initial: Pending, false holds, requestExpiresAt=createdAt+7d
- Root gate: denied until CloudPlatform exists (CONFLICT/VS0_CLOUDPLATFORM_ROOT_REQUIRED)
- Pair uniqueness: one non-terminal (platform, provider)
- Create: Idempotency-Key required; If-Match MALFORMED_REQUEST
- No PATCH registration (lifecycle through actions only)
- Pair uniqueness and collection-create idempotency remain separate mechanisms; idempotency uses principal, registered pattern, key, and classified-body digest
- Audit: successful create, audited denials; append-failure → INTERNAL_ERROR

Refs: FEATURE-0015 Task 14, REQ-F15-03/07/14/16/19/20/21/23/24,
DEC-0054, ADH-2026-046, ADH-2026-047, ADH-2026-048, ADH-2026-056
```

---

### Task 15: HTTP handlers for eight CloudProviderParticipation item-action routes

**Purpose:** Implement eight participation item-action routes (accept, reject, withdraw, suspend, resume, request-release, accept-release, decline-release) with per-route pipeline, empty JSON body contract, If-Match + Idempotency-Key preconditions, lifecycle source-state validation, independent hold logic, action-scoped authorization, lifecycle transition, idempotency (item actions bind concrete target UID), and audit.

**Dependencies:** Task 1 (types), Task 4 (validation/store primitives), Task 6 (lifecycle and scheduler), Task 8 (idempotency/audit coordination), Task 9 (grants), Task 14 (participation collection/item exists).

**Requirements traceability:** REQ-F15-03 (lifecycle), REQ-F15-06 (independent holds), REQ-F15-07 (action preconditions), REQ-F15-15 (audit atomicity), REQ-F15-17 (server-resolved grants), REQ-F15-18 (idempotency), REQ-F15-21 (empty action body), REQ-F15-23 (non-empty body MALFORMED_REQUEST), REQ-F15-24 (AUTH_REQUIRED).

**Design traceability:** Design §4.2 (eight action routes), §4.5 (lifecycle), §5.1 (per-route pipeline step 4 body validation, step 13 source-state validation), §6.1 (action-scoped authorization), §7.1 (handler tests), §8.1 (IMPLEMENT).

**Architecture decisions:** DEC-0054 (participation), ADH-2026-046 (action preconditions, independent holds), ADH-2026-050 (all-route idempotency coverage), ADH-2026-053 (action routes use /actions/<action>, not /{uid}:<action>), ADH-2026-054 (item-action idempotency binds concrete target UID; replay rechecks current auth/safe-access), ADH-2026-048 (non-empty body without valid duplicate-free top-level bootstrapGrant MALFORMED_REQUEST).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this task.

**Writable paths:**
- `internal/api/participation_actions.go` (eight action handlers)

**Tests:**
- `internal/api/participation_actions_test.go`: each of eight actions (accept, reject, withdraw, suspend, resume, request-release, accept-release, decline-release) with valid source state returns 200 + updated participation and exactly one AuditEvent; empty JSON body (EOF/zero bytes) accepted; whitespace/`{}`/non-empty body without valid duplicate-free top-level bootstrapGrant MALFORMED_REQUEST with no AuditEvent; If-Match + Idempotency-Key required; missing/malformed/stale If-Match returns `STALE_RESOURCE_VERSION` / 412 with no AuditEvent; when both a stale If-Match and invalid lifecycle source state are supplied, stale-version 412 wins and lifecycle validation is not reached; a current If-Match plus invalid source state returns `CONFLICT` / 409 with `VS0_PARTICIPATION_STATE_INVALID` and no AuditEvent; same-key/different-digest returns `CONFLICT` / 409 with `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH` and no AuditEvent; same-key/same-digest replays without a second AuditEvent; valid duplicate-free bootstrapGrant, forged headers, safe inaccessible target/reference, and any other registered audited denial append exactly one redacted AuditEvent; audit-append failure always returns INTERNAL_ERROR with no response/mutation/lifecycle/idempotency publication; action-scoped authorization (accept/reject/request-release: CloudPlatform grant; withdraw/accept-release/decline-release: CloudProvider grant; suspend/resume: resolve exactly one matching current scoped grant; zero/multiple matches fail closed); suspend sets only acting party hold; resume clears only acting party hold; clearing one hold does not reactivate while other remains true; suspend/resume denied in Pending/Terminating/terminal; decline-release resolves to correct hold-derived prior effective Active/Suspended; missing auth AUTH_REQUIRED; idempotency binds concrete target UID (no replay for another target); replay rechecks current auth/authz/safe-access (revoked grant denies without stored-result disclosure); action registration uses /actions/<action>, not /{uid}:<action>.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/api/...
```

**Acceptance criteria:**
- Eight action routes: POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/<action> where <action> is accept, reject, withdraw, suspend, resume, request-release, accept-release, decline-release.
- {uid} is a complete Go 1.22 http.ServeMux wildcard segment, read with r.PathValue("uid").
- Empty JSON body contract: EOF/zero bytes accepted; whitespace/`{}`/non-empty body without valid duplicate-free top-level bootstrapGrant MALFORMED_REQUEST.
- Action preconditions: If-Match + Idempotency-Key required; missing/malformed If-Match STALE_RESOURCE_VERSION.
- Lifecycle source-state validation: invalid source state CONFLICT/VS0_PARTICIPATION_STATE_INVALID.
- Independent hold logic: suspend sets only acting party hold; resume clears only acting party hold; clearing one does not reactivate while other true.
- Suspend/resume denied in Pending/Terminating/terminal phases.
- Action-scoped authorization: accept/reject/request-release (CloudPlatform grant); withdraw/accept-release/decline-release (CloudProvider grant); suspend/resume (resolve exactly one matching current scoped grant; zero/multiple fail closed).
- Lifecycle transitions per VS0-STATE-001.
- Idempotency: item actions bind concrete target UID; same-key/same-digest (including same normalized If-Match) replays stored success; same-key/different-digest (including changed normalized If-Match) CONFLICT/VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH; no replay for another target.
- Replay rechecks: current auth/authz/safe-access; revoked grant denies without stored-result disclosure.
- Audit: successful action and audited denials append AuditEvent; append-failure INTERNAL_ERROR.
- AC-F15-02 (lifecycle), AC-F15-10 (independent holds), AC-F15-11 (action preconditions), AC-F15-17 (empty action body), AC-F15-19 (non-empty body MALFORMED_REQUEST), AC-F15-20 (AUTH_REQUIRED).

**Security/observability impact:**
- Security: Action-scoped authorization; If-Match + Idempotency-Key required; replay rechecks current auth/authz/safe-access.
- Observability: Action transitions are logged and audited with request correlation.

**Exclusions:**
- No retired /{uid}:<action> form.
- No public expiry route (scheduler invokes api-server expiry transition).
- No PATCH registration (lifecycle through actions only).

**Commit message:**
```
feat(api): add eight CloudProviderParticipation item-action handlers

Implement eight participation item-action routes (accept, reject,
withdraw, suspend, resume, request-release, accept-release,
decline-release) with per-route pipeline, empty JSON body, If-Match +
Idempotency-Key preconditions, lifecycle source-state validation,
independent hold logic, action-scoped authorization, lifecycle
transition, idempotency (binds concrete target UID), and audit per
ADH-2026-046/050/053/054.

- Eight routes: /actions/<action> (not /{uid}:<action>)
- Empty body: EOF/zero bytes; whitespace/`{}`/non-empty without valid bootstrapGrant MALFORMED_REQUEST
- Preconditions: If-Match + Idempotency-Key; missing/malformed/stale If-Match STALE_RESOURCE_VERSION
- Source-state: invalid state CONFLICT/VS0_PARTICIPATION_STATE_INVALID
- Independent holds: suspend/resume set/clear only acting party hold; clearing one never reactivates while other true
- Suspend/resume denied in Pending/Terminating/terminal
- Action-scoped authorization: CloudPlatform (accept/reject/request-release); CloudProvider (withdraw/accept-release/decline-release); suspend/resume (resolve exactly one; zero/multiple fail closed)
- Idempotency: binds target UID; same-key/same-digest (same If-Match) replays success; same-key/different-digest (changed If-Match) CONFLICT
- Replay rechecks: current auth/authz/safe-access; revoked grant denies without disclosure
- Audit: successful action, audited denials; append-failure → INTERNAL_ERROR

Refs: FEATURE-0015 Task 15, REQ-F15-03/06/07/18/21/23/24, DEC-0054,
ADH-2026-046, ADH-2026-050, ADH-2026-053, ADH-2026-054
```

---

### Task 16: Route registration (35 explicit Go 1.22 http.ServeMux patterns)

**Purpose:** Register exactly 35 explicit Go 1.22 http.ServeMux method/path patterns for 22 logical endpoint paths in internal/server; no path-only or catch-all registration, reflection, or internal HTTP-method dispatch; HEAD matched by GET returns 405.

**Dependencies:** Task 11 (CloudPlatform and CloudProvider handlers), Task 13 (topology handlers), Task 14 (participation collection/item), Task 15 (participation actions).

**Requirements traceability:** REQ-F15-08 (PATCH-only/unsupported-method surface) and REQ-F15-24 (owned method/path patterns require authentication); route registration is not the collective implementation of REQ-F15-01 through REQ-F15-24.

**Design traceability:** Design DD-06 (exact route model), §4.2 (22 logical paths, 35 registrations), §5.1 (HEAD returns 405), §7.1 (route tests), §8.1 (IMPLEMENT).

**Architecture decisions:** ADH-2026-051 (22 logical paths, 35 registrations; no path-only/catch-all/reflection/dispatch), ADH-2026-053 (participation actions use /actions/<action>), ADH-2026-055 (HEAD returns 405).

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this task.

**Writable paths:**
- `internal/server/routes.go` (35 explicit method/path registrations; no path-only or catch-all registration)

**Tests:**
- `internal/server/routes_test.go`: verify exactly 35 explicit registrations for 22 logical paths (14 collection LIST/create + 12 PATCHable-item GET/PATCH + 1 participation-item GET-only + 8 participation actions); verify HEAD matched by GET returns 405 without adding registration; verify unsupported methods return 405; verify {uid} is a complete wildcard segment extracted with r.PathValue("uid"); verify no path-only or catch-all handler; verify no reflection or internal HTTP-method dispatch; verify /actions/<action> form (not /{uid}:<action>) for participation actions.

**Verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/server/...
```

**Acceptance criteria:**
- Exactly 35 explicit Go 1.22 http.ServeMux method/path patterns registered for 22 logical endpoint paths.
- 14 collection: GET and POST for CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, InfrastructureStack.
- 12 PATCHable-item: GET and PATCH for CloudPlatform, CloudProvider, HostingLocation, Datacenter, FaultDomain, InfrastructureStack.
- 1 participation-item: GET-only for CloudProviderParticipation (no PATCH registration).
- 8 participation actions: POST for accept, reject, withdraw, suspend, resume, request-release, accept-release, decline-release using /actions/<action>.
- No path-only or catch-all registration, reflection, or internal HTTP-method dispatch.
- {uid} is a complete Go 1.22 http.ServeMux wildcard segment, not a catch-all.
- HEAD matched by GET returns 405 without adding registration.
- PUT/DELETE on CloudPlatform, CloudProvider, topology return 405.
- AC-F15-01 (route surface), AC-F15-03 (no ExecutionTarget route), AC-F15-06 (PUT/DELETE 405).

**Security/observability impact:**
- Security: Explicit method/path registration prevents unintended route exposure.
- Observability: Route registration is logged at server startup.

**Exclusions:**
- No path-only or catch-all handler.
- No ExecutionTarget or ServiceRegion routes (FEATURE-0016/0022).
- No migration routes (DEC-0059).

**Commit message:**
```
feat(server): register 35 explicit Go 1.22 http.ServeMux patterns

Register exactly 35 explicit method/path patterns for 22 logical
endpoint paths per ADH-2026-051/053/055.

- 14 collection: GET+POST for 7 kinds
- 12 PATCHable-item: GET+PATCH for 6 kinds
- 1 participation-item: GET-only (no PATCH)
- 8 participation actions: POST /actions/<action>
- No path-only/catch-all/reflection/dispatch
- {uid} is complete wildcard, not catch-all
- HEAD matched by GET returns 405
- PUT/DELETE on CloudPlatform/CloudProvider/topology return 405

Refs: FEATURE-0015 Task 16, REQ-F15-08/24, ADH-2026-051,
ADH-2026-053, ADH-2026-055
```

---

### Task 17: Conformance tests for FEATURE-0015-local scenarios (VS0-CF-X03, VS0-CF-F15-01..41)

**Purpose:** Implement conformance tests for every FEATURE-0015-local scenario: VS0-CF-X03 (cross-provider safe denial) and VS0-CF-F15-01 through VS0-CF-F15-41 (FEATURE-0015 local proof cases).

**Dependencies:** Task 11 (CloudPlatform and CloudProvider handlers), Task 13 (topology handlers), Task 14 (participation collection/item), Task 15 (participation actions), and Task 16 (route registration).

**Requirements traceability:** All requirements (REQ-F15-01 through REQ-F15-24); all acceptance criteria (AC-F15-01 through AC-F15-20).

**Design traceability:** Design §7.1 (conformance tests), §7.2 (acceptance-scenario coverage), requirements §4.26 (acceptance-scenario coverage), requirements §4.28 (exact conformance semantics ledger), §8.1 (IMPLEMENT).

**Architecture decisions:** All controlling decisions and handoffs.

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this task.

**Writable paths:**
- `internal/apiconform/feature_0015_test.go` (one scenario per local conformance ID)

**Tests:**
- VS0-CF-X03: cross-provider target reference denied with safe RESOURCE_NOT_FOUND/VS0_AUTHORIZATION_SAFE_DENIAL; exactly one redacted AuditEvent with request correlation and no target-existence disclosure; audit-append failure INTERNAL_ERROR/500 with no publication.
- VS0-CF-F15-01 through VS0-CF-F15-41: implement exact inputs, expectedState, expectedError, expectedViolation (where registered), expectedSideEffects, and gate from registry (requirements §4.28).
- Anti-drift/non-runtime gates (AC-F15-05, AC-F15-07) remain verified by named repository checks, not invented runtime conformance.

Each test asserts exact Problem code, HTTP status, violation code (when registered), state outcome, and side effects (mutation/lifecycle/status/idempotency/AuditEvent presence or absence).

**Verification commands:**
```bash
make fmt
make test
make vet
go test -v ./internal/apiconform/...
go test -race ./internal/apiconform/... # concurrency-sensitive cases
```

**Acceptance criteria:**
- One conformance test per FEATURE-0015-local case: VS0-CF-X03, VS0-CF-F15-01..41.
- Each test asserts exact registry semantics (requirements §4.28): inputs, expectedState, expectedError, expectedViolation, expectedSideEffects, gate.
- Covers all approved acceptance criteria (AC-F15-01 through AC-F15-20).
- Covers all approved requirements (REQ-F15-01 through REQ-F15-24).
- Anti-drift/non-runtime gates (AC-F15-05, AC-F15-07) remain verified by repository checks, not runtime tests.
- All acceptance scenarios pass.

**Security/observability impact:**
- Security: Conformance tests verify safe denial, audit atomicity, replay recheck, forged-grant denial, writer enforcement.
- Observability: Conformance tests verify AuditEvent correlation and redaction.

**Exclusions:**
- No downstream-owned cases (VS0-CF-HP01, VS0-CF-F01, VS0-CF-F09, VS0-CF-T01) as local proof.
- No migration conformance (retired: VS0-CF-MIG01/02, VS0-CF-MIGF01/02/03).

**Commit message:**
```
test(conformance): add FEATURE-0015 local scenario tests

Implement conformance tests for VS0-CF-X03 and VS0-CF-F15-01..41
per requirements §4.26, §4.28, covering all approved acceptance
criteria (AC-F15-01..20) and requirements (REQ-F15-01..24).

- VS0-CF-X03: cross-provider safe denial with audited non-disclosure
- VS0-CF-F15-01..41: exact registry inputs/expectedState/expectedError/expectedViolation/expectedSideEffects
- All approved acceptance scenarios pass
- Anti-drift/non-runtime gates verified by repository checks

Refs: FEATURE-0015 Task 17, all REQs, all ACs, all controlling
DEC/ADH
```

---

### Task 18: Repository verification and boundary validation (non-commit final checkpoint)

**Purpose:** Perform full repository verification, boundary validation, controlled API-server smoke verification, and clean-tree confirmation before marking FEATURE-0015 complete. This is a non-commit checkpoint: all generated output is confined to known current-run temporary paths outside the working tree.

**Dependencies:** All eleven implementation commit tasks: 1, 4, 6, 8, 9, 11, 13, 14, 15, 16, and 17.

**Requirements traceability:** All requirements (REQ-F15-01 through REQ-F15-24).

**Design traceability:** Design completeness and unresolved-decision report (§10); §7.3 (verification commands).

**Architecture decisions:** All controlling decisions and handoffs.

**Risks addressed:** Approved architecture boundary, anti-drift, and execution controls applicable to this checkpoint.

**Writable paths:**
- None (verification only).

**Tests:**
- None (runs existing test suites and verification commands).

**Verification commands:**
```bash
# FEATURE-0015 gates
make feature-0015-architecture-readiness
make feature-contract-check FEATURE=FEATURE-0015
make feature-0015-formal-check

# Slice 0 and Phase 2R gates
make vs000-contract-check
make phase2r-drift-check

# Repository integrity; documentation output is confined to a known current-run
# temporary directory, never the configured working-tree site directory.
docs_dir="$(mktemp -d)"
mkdocs build --strict --site-dir "$docs_dir"
rm -rf "$docs_dir"
make structurizr-check

# Go quality gates (formatting has already been performed and committed by
# the preceding implementation tasks; this non-mutating checkpoint only
# verifies that no formatter-induced change remains)
make test
make vet
go test -race ./...

# Feature gate
make ff-feature-gate FEATURE=FEATURE-0015

# Controlled API runtime smoke check. This preserves start/readiness/probe
# failure status and always terminates and waits for the actual child.
set -eu
api_log="$(mktemp)"
server_dir="$(mktemp -d)"
server_bin="$server_dir/sovrunn-api"
server_pid=''
cleanup() {
  smoke_status=$?
  if [ -n "$server_pid" ] && kill -0 "$server_pid" 2>/dev/null; then
    kill "$server_pid" 2>/dev/null || true
    wait "$server_pid" 2>/dev/null || true
  fi
  rm -f "$api_log"
  rm -rf "$server_dir"
  trap - EXIT INT TERM
  exit "$smoke_status"
}
trap cleanup EXIT INT TERM
go build -o "$server_bin" ./cmd/sovrunn-api
"$server_bin" --config configs/sovrunn-api.local.yaml >"$api_log" 2>&1 & server_pid=$!
ready=false
for attempt in $(seq 1 30); do
  if curl --connect-timeout 1 --max-time 3 -fsS http://127.0.0.1:8080/readyz >/dev/null; then ready=true; break; fi
  kill -0 "$server_pid" 2>/dev/null || { cat "$api_log"; exit 1; }
  sleep 1
done
[ "$ready" = true ] || { cat "$api_log"; exit 1; }
curl --connect-timeout 1 --max-time 3 -fsS http://127.0.0.1:8080/healthz
curl --connect-timeout 1 --max-time 3 -fsS http://127.0.0.1:8080/readyz
kill "$server_pid"
set +e
wait "$server_pid"
wait_status=$?
set -e
server_pid=''
case "$wait_status" in
  0|143) ;; # graceful exit or expected SIGTERM termination
  *) trap - EXIT INT TERM; rm -f "$api_log"; rm -rf "$server_dir"; exit "$wait_status" ;;
esac
trap - EXIT INT TERM
rm -f "$api_log"
rm -rf "$server_dir"

# Source-of-truth cleanliness checks
git diff --check
git diff
git status --short
```

**Acceptance criteria:**
- `make feature-0015-architecture-readiness` passes (field classification, route-contract completeness, single-scope-derivation-source, stale-wording/absent-create-field/client-supplied-scope/unowned-participation-field/missing-create-status rejection per ADH-2026-047 decision 5).
- `make feature-contract-check FEATURE=FEATURE-0015` passes (complete pipeline: body/headers → authn/authz → root/reference/safe-denial ordering → structural/semantic validation → exact Problem outcome → state/status result → audit effect → idempotency/race outcome → local conformance test per ADH-2026-048 decision 4).
- `make feature-0015-formal-check` runs TLC successfully (and fails if TLC is unavailable) for `ParticipationLifecycle.tla`/`.cfg` and `IdempotencyAuditPublication.tla`/`.cfg`, proving the modeled participation/hold and idempotency/audit-publication invariants.
- `make vs000-contract-check` passes (registry YAML parses; schemas/states/Problem mappings are structurally valid; every required local conformance case is registered; no field in implementation is absent from the registry).
- `make phase2r-drift-check` passes and active-authority exclusions remain enforced.
- `mkdocs build --strict --site-dir <current-run-mktemp-dir>` passes and emits no documentation output into the working tree.
- `make structurizr-check` passes.
- The preceding implementation tasks have run `make fmt`; this non-mutating checkpoint runs `make test`, `make vet`, and `go test -race ./...` and fails if `git diff` reveals a formatter-induced change.
- `make ff-feature-gate FEATURE=FEATURE-0015` passes (no drift, no missing acceptance, no staged generated artifacts, Phase 2R scope boundaries satisfied).
- A temporary-path `go build -o <current-run-mktemp-dir>/sovrunn-api ./cmd/sovrunn-api` succeeds, then the checkpoint starts that actual binary; `/healthz` and `/readyz` both return successful responses; the owned process is then terminated and successfully waited. No `bin/` artifact is created or updated.
- The bounded smoke procedure starts the direct server binary (which must not fork descendants), uses unique `mktemp` artifacts, waits at most 30 seconds for readiness, uses bounded curl connect/overall timeouts, preserves start/readiness/probe/termination failure status, explicitly accepts only graceful exit or expected SIGTERM status, and cannot leave the owned API-server process running.
- `git diff --check` passes; `git diff` and `git status --short` are reviewed.
- Any unexpected or generated artifact is reported as a failed cleanliness check. It is not deleted or otherwise cleaned up by this task; only an exact, current-run artifact explicitly identified by its generating command may be handled under separate user direction.
- Any formatter-induced change detected by `git diff` or `git status --short` fails this non-mutating checkpoint; it is not silently accepted or cleaned up.
- All acceptance criteria (AC-F15-01 through AC-F15-20) satisfied.
- No orphan IDs (every referenced REQ, AC, DEC, ADH, VS0-SCHEMA, VS0-WRITER, VS0-STATE, VS0-CF, violation/Problem code exists in loaded authorities).
- No ARCHITECTURE_DECISION_REQUIRED, REQUIREMENT_CLARIFICATION_REQUIRED, BOUNDARY_CHANGE_REQUIRED, DEPENDENCY_APPROVAL_REQUIRED, or SECURITY_REVIEW_REQUIRED condition triggered.

**Security/observability impact:**
- Security: Verification confirms no secret disclosure, safe denial, writer enforcement, audit atomicity.
- Observability: Verification confirms AuditEvent correlation, redaction, structured logging.

**Exclusions:**
- No implementation beyond approved design.
- No future-feature activation.
- No prohibited concept reintroduction.

**Commit message:** None — Task 18 is a verification checkpoint. It records results only and does not produce a commit or alter the working tree.

---

## 3. No-task ledger (CONTRACT_ONLY, EXCLUDED, or anti-drift verification gates)

The following approved design elements are not decomposed into implementation tasks because they are consumed by reference, excluded, or verified by anti-drift/non-runtime gates per ADH-2026-046 decision 3 and design §8.

| Element | Disposition | Verification method |
|---------|-------------|---------------------|
| REQ-F15-12 (existing FEATURE-0012 Problem codes only) | CONTRACT_ONLY; anti-drift/non-runtime | Proven by construction (§5.2 closed set); verified by `make vs000-contract-check` |
| REQ-F15-13 (no migration plan/record/controller/cutover) | CONTRACT_ONLY; anti-drift/non-runtime | Proven by absence (§6.4 retired tombstones); verified by `make phase2r-drift-check` |
| AC-F15-05 (no migration mechanism exists) | CONTRACT_ONLY; anti-drift/non-runtime | Proven by absence; verified by `make phase2r-drift-check` |
| AC-F15-07 (stale-concept anti-drift gate) | CONTRACT_ONLY; anti-drift/non-runtime | Proven by prohibition (§10.1); verified by `make vs000-contract-check` and `make phase2r-drift-check` |
| AC-F15-09 (FEATURE-0012 codes only) | CONTRACT_ONLY; contract-reuse claim | Proven by closed set (requirements §4.28); no standalone runtime scenario |
| FEATURE-0011 reuse governance contract | CONTRACT_ONLY | Consumed by reference (design §1.2); no FEATURE-0015 task changes it |
| FEATURE-0012 Problem Details, ObjectMeta, TypedRef, Condition, reference constraints, validation pipeline | CONTRACT_ONLY (ScopeKind is explicit DD-03 compatibility exception) | Consumed by reference (design §3.2); ScopeKind canonicalization is Task 1 work package 2 |
| FEATURE-0013 DecisionRecord/AuditEvent envelope, linkage, writer | CONTRACT_ONLY | Consumed by reference (design §3.2); no FEATURE-0015 task changes FEATURE-0013 authority |
| Writer boundaries VS0-WRITER-002/003/004/005 | CONTRACT_ONLY | Consumed by reference (design §6.1); enforced by handler authorization |
| Existing observability baseline | CONTRACT_ONLY | Consumed as-is (design §6.3); redacted logs/metrics use existing baseline |
| ExecutionTarget and its lifecycle/schema/routes/status/writer/conformance (FEATURE-0016) | EXCLUDED | No task implements; proven by absence (design §1.3, §10); verified by `make vs000-contract-check` |
| PolicyEvaluationRequest/Result, PolicyEngineAdapter (FEATURE-0017) | EXCLUDED | No task implements; proven by absence (design §1.3, §10) |
| GovernanceProfile, Membership, RoleDefinition, RoleAssignment (FEATURE-0018) | EXCLUDED | No task implements; F0015 persists no roles/memberships (design §9 grants) |
| SovereigntyProfile, SovereigntyFactSet, EvidenceRecord (FEATURE-0019) | EXCLUDED | No task implements; proven by absence (design §1.3, §10) |
| ProfileAssignment, EffectiveGovernanceContext (FEATURE-0020) | EXCLUDED | No task implements; proven by absence (design §1.3, §10) |
| CloudEnrollment, EntitlementPackage, ServiceEntitlement, QuotaPolicy, provider-selection intent, spec.providerSelectionModes, spec.permittedHostingLocationRefs, Organization visibility eligibility (FEATURE-0021) | EXCLUDED | No task implements; F0021 fields rejected (design §10); proven by absence (design §1.3, §10) |
| ServiceTypeDefinition, ServiceOffering, ServicePlan, ServiceRuntimeProfile, ServiceRequirementSet, ServiceRegion (FEATURE-0022) | EXCLUDED | No task implements; proven by absence (design §1.3, §10) |
| DecisionProfiles, ServicePlacement, PlacementDecision (FEATURE-0023) | EXCLUDED | No task implements; proven by absence (design §1.3, §10) |
| PluginExecution, ServiceDeploymentPlan, fake harness (FEATURE-0024) | EXCLUDED | No task implements; proven by absence (design §1.3, §10) |
| DecisionOperationContext, AI-readable projection (FEATURE-0025) | EXCLUDED | No task implements; proven by absence (design §1.3, §10) |
| Slice 0 integration demo, cross-feature orchestration; happy-path proof VS0-CF-HP01, trace proof VS0-CF-T01 (FEATURE-0026) | EXCLUDED | No task implements; downstream cases cited as integration references only (design §10) |
| Federation, remote deployment resource, signed remote offer, cross-control-plane replication; platform lifecycle (No Phase 2R owner) | EXCLUDED | No task implements; proven by absence (design §10) |
| Real provisioning (Phase 3) | EXCLUDED | No task implements; proven by absence per PHASE2R non-goals |
| Prohibited stale concepts (ResourcePool, ProviderCapability, generic Provider, ServiceClass, EffectivePolicyContext, six-scope vocabulary, SovrunnInstallation, CanonicalMigrationPlan/Record/controller/cutover) | EXCLUDED; anti-drift verification | Prohibited (design §10.1); verified by `make vs000-contract-check` and `make phase2r-drift-check` |

---

## 4. Canonical coverage ledger

Every approved REQ and AC appears exactly once below with its task disposition. No ID is created, renumbered, merged, split, reinterpreted, or omitted.

### 4.1 Requirements coverage

| REQ ID | Task(s) | Coverage summary |
|--------|---------|------------------|
| REQ-F15-01 | Task 1, 4, 11 | CloudPlatform type, validation/uniqueness, handlers |
| REQ-F15-02 | Task 1, 4, 11 | CloudProvider type, ISO dataset, validation/uniqueness, handlers |
| REQ-F15-03 | Task 1, 4, 6, 14, 15 | CloudProviderParticipation type, lifecycle, pair uniqueness, handlers, actions |
| REQ-F15-04 | Task 1, 4, 13 | Topology chain types, ISO dataset (HostingLocation), validation/uniqueness, handlers |
| REQ-F15-05 | No-task ledger (EXCLUDED) | Feature ends at InfrastructureStack; no ExecutionTarget introduced |
| REQ-F15-06 | Task 6, 15 | Lifecycle actions, independent holds, scheduler expiry, action handlers |
| REQ-F15-07 | Task 6, 14, 15 | No mutable participation spec; create Idempotency-Key only; action preconditions |
| REQ-F15-08 | Task 11, 13, 14 | PATCH-only merge-patch surface; PUT/DELETE 405 |
| REQ-F15-09 | Task 4, 13 | Safe cross-provider denial validation and handler enforcement |
| REQ-F15-10 | Task 1 | Seven-scope vocabulary canonicalization |
| REQ-F15-11 | Task 4, 9, 11, 13, 15 | Writer enforcement validation, grants, handler denial |
| REQ-F15-12 | No-task ledger (CONTRACT_ONLY) | Existing FEATURE-0012 codes only; proven by construction |
| REQ-F15-13 | No-task ledger (CONTRACT_ONLY) | No migration mechanism; proven by absence |
| REQ-F15-14 | Task 4, 14 | Scope-reference UID invariant validation and handler enforcement |
| REQ-F15-15 | Task 8, 11, 13, 14, 15 | Audit atomicity, append-failure mapping, handler coordination |
| REQ-F15-16 | Task 11, 13, 14 | CloudPlatform-root requirement validation and handler enforcement, including CloudProvider create in Task 11 work package 12 |
| REQ-F15-17 | Task 9, 11, 13, 14, 15 | Bootstrap-grant resolver, forged-grant denial, handler enforcement |
| REQ-F15-18 | Task 8, 11, 13, 14, 15 | Idempotency coordinator, replay semantics, handler integration |
| REQ-F15-19 | Task 4, 11, 13, 14 | Closed create-contract validation and handler enforcement |
| REQ-F15-20 | Task 4, 11, 13, 14 | Single-source scope derivation validation and handler enforcement |
| REQ-F15-21 | Task 4, 14, 15 | Participation create body and empty item-action body validation |
| REQ-F15-22 | Task 1, 4, 11, 13 | Assigned ISO-3166-1 alpha-2 dataset, validation, handler enforcement |
| REQ-F15-23 | Task 4, 11, 13, 14, 15 | Malformed/prohibited-input validation (headers, body) |
| REQ-F15-24 | Task 11, 13, 14, 15 | AUTH_REQUIRED on every owned pattern with no side effect |

### 4.2 Acceptance criteria coverage

| AC ID | Task(s) | Coverage summary |
|-------|---------|------------------|
| AC-F15-01 | Task 1, 4, 11, 13, 16, 17 | CloudPlatform/CloudProvider/topology collection create, GET/LIST, permitted PATCH with scope subsets; ends at InfrastructureStack |
| AC-F15-02 | Task 1, 4, 6, 14, 15, 17 | Participation lifecycle, acceptance guard, pair uniqueness |
| AC-F15-03 | No-task ledger (EXCLUDED), Task 16, 17 | No ExecutionTarget behavior introduced; verified by route surface |
| AC-F15-04 | Task 4, 13, 17 | Cross-provider target reference denied safely |
| AC-F15-05 | No-task ledger (CONTRACT_ONLY) | No migration mechanism exists; anti-drift/non-runtime gate |
| AC-F15-06 | Task 11, 13, 16, 17 | PUT/DELETE 405; PATCH-only merge-patch |
| AC-F15-07 | No-task ledger (CONTRACT_ONLY) | Stale-concept anti-drift gate; verified by repository checks |
| AC-F15-08 | Task 4, 9, 11, 13, 15, 17 | Writer enforcement; status/system-owned field denial; safe cross-provider denial |
| AC-F15-09 | No-task ledger (CONTRACT_ONLY) | FEATURE-0012 codes only; contract-reuse claim; closed-set proof |
| AC-F15-10 | Task 6, 15, 17 | Independent holds; suspend/resume denied in terminal/Pending/Terminating; clearing one never reactivates |
| AC-F15-11 | Task 6, 14, 15, 17 | No mutable participation spec; create Idempotency-Key only; action If-Match + Idempotency-Key |
| AC-F15-12 | Task 4, 14, 17 | Scope-reference UID invariant enforced; mismatch rejected |
| AC-F15-13 | Task 8, 11, 13, 14, 15, 17 | Lifecycle/hold-changing actions, create/PATCH, denials produce correlated redacted AuditEvent |
| AC-F15-14 | Task 13, 14, 17 | CloudPlatform-root requirement enforced until CloudPlatform exists |
| AC-F15-15 | Task 4, 11, 13, 14, 17 | Closed field boundary; server-owned/status/scope/unknown/deferred rejected; 201 + exact initial status |
| AC-F15-16 | Task 4, 11, 13, 14, 17 | Single-source server-side scope derivation; client scope never honored |
| AC-F15-17 | Task 4, 14, 15, 17 | Participation create accepts exactly four fields; initializes Pending/false holds/seven-day expiry; rejects F0021 fields; item actions empty body |
| AC-F15-18 | Task 1, 4, 11, 13, 17 | Syntactically valid but unassigned ISO rejected; administrativeAreaCode syntax-plus-prefix only |
| AC-F15-19 | Task 4, 9, 11, 13, 14, 15, 17 | Malformed Idempotency-Key, If-Match on create, non-empty action body MALFORMED_REQUEST; forged-grant AUTHORIZATION_DENIED |
| AC-F15-20 | Task 11, 13, 14, 15, 17 | Missing/invalid authentication AUTH_REQUIRED with no side effect |

---

## 5. Orphan and unresolved report

### 5.1 Orphan ID check

All referenced IDs (REQ-F15-01 through REQ-F15-24, AC-F15-01 through AC-F15-20, DEC-0037/0041/0042/0054/0059, ADH-2026-020/024/025/037/042/043/045/046/047/048/049/050/051/052/053/054/055/056/057, VS0-SCHEMA-001..014, VS0-WRITER-002/003/004/005/011, VS0-STATE-001, VS0-CF-X03, VS0-CF-F15-01..41, Problem codes MALFORMED_REQUEST/UNKNOWN_FIELD/DUPLICATE_FIELD/REQUEST_TOO_LARGE/AUTH_REQUIRED/AUTHORIZATION_DENIED/RESOURCE_NOT_FOUND/ALREADY_EXISTS/CONFLICT/STALE_RESOURCE_VERSION/UNSUPPORTED_MEDIA_TYPE/VALIDATION_FAILED/INTERNAL_ERROR, violation codes VS0_PARTICIPATION_STATE_INVALID/VS0_PARTICIPATION_DUPLICATE/VS0_PARTICIPATION_ACCEPTANCE_MISSING/VS0_STATUS_FIELD_WRITE/VS0_SYSTEM_OWNED_FIELD_WRITE/VS0_SCOPE_KIND_INVALID/VS0_PATCH_IMMUTABLE_FIELD/VS0_CLOUDPLATFORM_ROOT_REQUIRED/VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH/VS0_SCOPE_REFERENCE_MISMATCH/VS0_TOPOLOGY_PROVIDER_MISMATCH/VS0_AUTHORIZATION_SAFE_DENIAL/VS0_AUTH_REQUIRED) exist in the loaded authorities.

Task 9 has no `AC-F15-22` reference: its bootstrap-grant boundary is `REQ-F15-17` with exact local proof `VS0-CF-F15-22`; the associated acceptance coverage is the forged-grant exception within `AC-F15-19`.

Retired tombstones (VS0-SCHEMA-057/060/061, VS0-WRITER-020/021, VS0-STATE-011, VS0-CF-MIG01/02, VS0-CF-MIGF01/02/03) are explicitly named as retired and never reused; they are not orphan IDs.

### 5.2 Unresolved decisions

None. No `ARCHITECTURE_DECISION_REQUIRED`, `REQUIREMENT_CLARIFICATION_REQUIRED`, `BOUNDARY_CHANGE_REQUIRED`, `DEPENDENCY_APPROVAL_REQUIRED`, or `SECURITY_REVIEW_REQUIRED` condition was triggered. All requirements map exactly to approved design elements; all design elements are either IMPLEMENT, CONTRACT_ONLY, or EXCLUDED with clear disposition.

### 5.3 Implementation completeness

- The eleven implementation commit tasks—1, 4, 6, 8, 9, 11, 13, 14, 15, 16, and 17—with inseparable internal work packages 2, 3, 5, 7, 10, and 12 collectively cover all IMPLEMENT design elements; Task 18 is the non-commit verification checkpoint.
- All 24 approved requirements (REQ-F15-01 through REQ-F15-24) are covered exactly once in §4.1.
- All 20 approved acceptance criteria (AC-F15-01 through AC-F15-20) are covered exactly once in §4.2.
- All CONTRACT_ONLY and EXCLUDED design elements are enumerated in the no-task ledger (§3) with their exact verification method.
- Anti-drift/non-runtime gates (AC-F15-05, AC-F15-07, REQ-F15-12, REQ-F15-13) remain verified by named repository checks per ADH-2026-046 decision 3.
- No task implements a non-goal, future-feature behavior, or prohibited concept.
- Every task cites exact requirements, design sections/decisions, architecture decisions, and applicable risks.
- Dependency order is explicit; independent tasks are marked.
- Task 18 performs full repository verification, controlled runtime smoke checks, boundary validation, and fail-closed clean-tree reporting.

---

## Model Execution Report:
- Tool: kiro
- Stage or task: FEATURE-0015 Tasks (`.kiro/specs/canonical-cloud-model-and-alpha-migration/tasks.md`)
- Recommended priority list: 1) Task decomposition (`claude-sonnet-4.5`), effort: medium; 2) Complex fallback (`claude-opus-4.8`), effort: high
- Selected model: claude-sonnet-4.5
- Effort/reasoning setting: medium (structured implementation sequencing and acceptance criteria)
- Fallback used: no
- Fallback reason: none

STAGE_STATUS: COMPLETE
