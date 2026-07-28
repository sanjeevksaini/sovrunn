---
doc_type: traceability
title: ADH-2026-015 Scope Vocabulary Compatibility Evidence
status: historical-pending-consolidation
phase: 2
feature_id: FEATURE-0013
proposed_successor: ADH-2026-017
depends_on:
  - FEATURE-0012
ai_load_priority: historical
ai_summary: Historical compatibility analysis for ADH-2026-015. It is not a downstream architecture input. Proposed ADH-2026-017 replaces the fragmented FEATURE-0013 handoffs with the consolidated architecture, preserves FEATURE-0012's six ScopeKind values, and models ServiceInstance as a typed subject.
---

# ADH-2026-015 — Scope Vocabulary Compatibility Evidence

## 1. Purpose and status

This document records the historical compatibility analysis performed for
ADH-2026-015. It is retained for provenance only and must not be loaded as a
requirements, design, task, implementation, or reviewer instruction.

**This document is not authoritative for the consolidated FEATURE-0013
architecture.** Proposed replacement ADH-2026-017 preserves FEATURE-0012's six
governance `ScopeKind` values and represents `ServiceInstance` as a typed
subject under its governing scope. On approval, ADH-2026-017 and the exact
digest-bound consolidated architecture become the sole normative input.

**ADH-2026-016 relationship.** ADH-2026-016 (Approved, 2026-07-27) does **not**
change this seven-value contract. It is referenced here only for the **newly
clarified boundaries** it added on top of the seven-value contract:
`DecisionRecord` `metadata.scopeRef` as the sole scope authority (no top-level
or parallel scope source), the provider-neutral sensitivity vocabulary, the
structural-only security-validation boundary, the structural-trust versus
deferred-cryptographic boundary (ADR-F13-002), the exact conformance coverage
`F13-CF-01` through `F13-CF-28`, and the SUPERSEDED disposition of the six-scope
evidence. See `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
section 30.

- Controlling handoffs: **ADH-2026-014** (Approved, 2026-07-27),
  **ADH-2026-015** (Approved, 2026-07-27), and **ADH-2026-016** (Approved,
  2026-07-27). ADH-2026-015 supersedes only the earlier six-scope `AuditEvent`
  statement in ADH-2026-014; ADH-2026-016 clarifies previously unresolved
  contract boundaries without changing the seven-value contract; all other
  ADH-2026-014 decisions remain controlling.
- Evidence status: **PENDING_IMPLEMENTATION_EVIDENCE**.
- Requirements status dependency: continued design, tasks, and implementation
  remain **unauthorized** until requirements receive fresh independent approval.

This file does **not** assert that any schema, Go binding, validator, fixture,
or test has been updated or has passed. It records the current inspected state
and the obligations that must be met later.

## 2. Approved decision being evidenced

ADH-2026-015 approves exactly:

1. `ScopeKind` is the single canonical shared vocabulary for scope identity.
2. The canonical values are exactly: `Platform`, `Organization`,
   `OrganizationUnit`, `Tenant`, `Project`, `Provider`, and `ServiceInstance`.
3. `ServiceInstance` is an additive extension; `Provider` remains unchanged.
4. `AuditEvent` initially permits all seven canonical values.
5. `AuditEvent` uses `metadata.scopeRef` as its sole canonical scope identity.
6. No separate `AuditEvent` scope enum or independently mutable scope field is
   permitted.
7. Existing FEATURE-0012 serialized `ScopeKind` values and existing
   Organization-scoped `AuditEvent`s remain valid.

## 3. FEATURE-0012 artifacts inspected (read-only)

The following canonical FEATURE-0012 files were inspected without modification.

| Artifact class | Exact file path | Relevant content observed |
|---|---|---|
| Canonical JSON schema | `api/schemas/_common/scope-ref.json` | `kind` enum lists exactly six values: `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `Provider`. `ServiceInstance` is absent. |
| Baseline JSON schema copy | `api/schemas/baseline/_common/scope-ref.json` | Baseline-pinned copy of the six-value `scope-ref.json` enum. |
| Common metadata schema | `api/schemas/_common/object-meta.json` | `metadata.scopeRef` is defined via `$ref` to `scope-ref.json`; scope is expressed only through `metadata.scopeRef`. |
| AuditEvent schema | `api/schemas/audit-event.json` | `x-sovrunn-allowed-scopes: ["Organization"]`; `x-sovrunn-stability: alpha`. Scope is carried only through `metadata` (`object-meta.json` → `scope-ref.json`). There is no separate `AuditEvent` scope property or scope enum. |
| Go binding | `internal/apimeta/scope.go` | `type ScopeKind string`; constants `ScopePlatform`, `ScopeOrganization`, `ScopeOrganizationUnit`, `ScopeTenant`, `ScopeProject`, `ScopeProvider`; `AllScopeKinds()` returns exactly those six; `Valid()` accepts only those six; `NormalizeScope` and `CanonicalScopeIdentity` handle the canonical nil-platform form. |
| Go binding test | `internal/apimeta/scope_test.go` | `TestMatrixBScopeKinds` asserts the closed set is exactly the six values and that an unknown kind (`Cluster`) is invalid. |
| Scope reference validator | `internal/apivalid/stage_reference.go` | Enforces `scopeRef.kind` against the Matrix B set and validates against schema `x-sovrunn-allowed-scopes`. |
| Allowed-scopes annotation reader | `internal/apischema/annotations.go` | Reads the `x-sovrunn-allowed-scopes` schema annotation used to constrain per-kind scope acceptance. |
| Conformance harness | `internal/apiconform/fitness_ref.go`, `internal/apiconform/fixtures.go`, `internal/apiconform/fixtures_test.go`, `internal/apiconform/canonical_schemas_test.go` | Drive the canonical schema and fixture conformance, including scope validation. |
| Positive AuditEvent fixture | `tests/conformance/fixtures/audit-event.json` (+ `.yaml`) | Uses `metadata.scopeRef` with `kind: Organization`. Confirms `AuditEvent` scope identity is expressed only through `metadata.scopeRef`. |
| Positive scope fixtures | `tests/conformance/fixtures/operation-platform.json`, `operation-organization.json`, `operation-organizationunit.json`, `operation-tenant.json`, `operation-project.json`, `operation-provider.json` (+ `.yaml`) | Positive coverage for each of the six existing `ScopeKind` values. |
| Negative scope fixture | `tests/conformance/fixtures/negative/invalid-scope-kind.json` | A non-Matrix-B kind (`Cluster`) is rejected. |

## 4. Value-by-value compatibility table

Every inherited controlled-vocabulary value is classified using the required
change taxonomy: **added, removed, renamed, mapped, restricted, expanded, or
semantically reinterpreted**.

| Value | In FEATURE-0012 `ScopeKind` | In FEATURE-0013 (ADH-2026-015) | Classification | Serialization impact | Notes |
|---|---|---|---|---|---|
| `Platform` | Present | Present | Unchanged (no delta) | None | Canonical logical and serialized form is an absent/nil `metadata.scopeRef` per FEATURE-0012 D-16 (`NormalizeScope`, `internal/apimeta/scope.go:84-92`; `CanonicalScopeIdentity`, `internal/apimeta/scope.go:107-112`); unchanged. Absent `scopeRef` resolves to `Platform` only where the contract permits `Platform`, else the stable required-scope error applies. |
| `Organization` | Present | Present | Unchanged (no delta) | None | Existing Organization-scoped `AuditEvent`s remain valid. |
| `OrganizationUnit` | Present | Present | Unchanged (no delta) | None | — |
| `Tenant` | Present | Present | Unchanged (no delta) | None | — |
| `Project` | Present | Present | Unchanged (no delta) | None | — |
| `Provider` | Present | Present | Unchanged (no delta) — retained | None | ADH-2026-014's six-scope `AuditEvent` list had omitted `Provider`; ADH-2026-015 retains it. Reconciled by permitting all seven values for `AuditEvent`. |
| `ServiceInstance` | Absent | Present | **Added** (additive extension) | Additive enum value only; no change to existing serialized values | ADH-2026-014 required `ServiceInstance` for `AuditEvent` although it was absent from the FEATURE-0012 shared vocabulary. ADH-2026-015 adds it to the shared `ScopeKind` vocabulary. |

Additional controlled difference at the `AuditEvent` profile level:

| Contract element | FEATURE-0012 current | FEATURE-0013 (ADH-2026-015) | Classification |
|---|---|---|---|
| `AuditEvent` `x-sovrunn-allowed-scopes` | `["Organization"]` (alpha) | All seven canonical `ScopeKind` values | **Expanded** (from one permitted scope to seven) |
| `AuditEvent` scope identity field | `metadata.scopeRef` only (no separate scope enum) | `metadata.scopeRef` only (confirmed sole authority; no parallel enum) | Unchanged / confirmed authoritative |

## 5. Serialization and backward-compatibility assessment

- All six existing FEATURE-0012 serialized `ScopeKind` values remain valid and
  unchanged. No value is removed, renamed, mapped, restricted, or semantically
  reinterpreted.
- `ServiceInstance` is a purely additive enum value. Adding it does not alter
  the meaning or serialized form of any existing value.
- The canonical platform-scope form (absent/nil `scopeRef`, `PlatformScopeUID`
  sentinel) defined by FEATURE-0012 is not affected.
- Existing Organization-scoped `AuditEvent`s remain valid because
  `Organization` remains permitted and `metadata.scopeRef` remains the scope
  identity field.
- The `AuditEvent` scope-acceptance change is an **expansion** of
  `x-sovrunn-allowed-scopes` from `["Organization"]` to all seven canonical
  values. Expansion of permitted values is backward compatible for previously
  accepted records.

## 6. Migration classification

- Overall migration classification: **additive, backward-compatible extension**.
- `ScopeKind` vocabulary: expanded by one additive value (`ServiceInstance`).
- `AuditEvent` permitted scopes: expanded from one to seven.
- No breaking change to any existing serialized value or previously accepted
  `AuditEvent`.
- No renamed, removed, mapped, restricted, or semantically reinterpreted value.

## 7. Affected consumers

- `api/schemas/_common/scope-ref.json` (and its baseline copy) — the `kind`
  enum would need the additive `ServiceInstance` value.
- `api/schemas/audit-event.json` — `x-sovrunn-allowed-scopes` would need
  expansion to the seven canonical values.
- `internal/apimeta/scope.go` — `ScopeKind` constants, `AllScopeKinds()`, and
  `Valid()` would need the additive `ServiceInstance` value.
- `internal/apimeta/scope_test.go` — the closed-set assertion would need to
  reflect seven values.
- `internal/apivalid/stage_reference.go` and `internal/apischema/annotations.go`
  — scope and allowed-scopes validation would need to accept the seven-value
  vocabulary and the expanded `AuditEvent` acceptance.
- `internal/apiconform/*` and `tests/conformance/fixtures/*` — positive
  `ServiceInstance` coverage and expanded `AuditEvent` scope coverage would be
  required.
- Downstream Phase 2 features that inherit `ScopeKind` would inherit the
  seven-value vocabulary.

## 8. Regression obligations (to be satisfied later, not yet met)

The following obligations must be satisfied during a later, separately approved
implementation stage. They are recorded here as obligations only.

- Every one of the six existing FEATURE-0012 values must continue to validate.
- Existing Organization-scoped `AuditEvent` fixtures must continue to pass.
- At least one positive fixture must exercise `ServiceInstance` as a valid
  `ScopeKind` value for `AuditEvent`.
- The positive `Platform` fixture must use the canonical absent/nil
  `metadata.scopeRef` form; an absent `metadata.scopeRef` must resolve
  deterministically to `Platform` only where the contract permits `Platform`.
- A negative fixture must assert that an absent `metadata.scopeRef` is rejected
  with the stable required-scope error under a contract that does not permit
  `Platform` (absence is not automatically valid).
- Validators must reject any scope value outside the canonical seven.
- No separate `AuditEvent` scope enum or independently mutable scope field may
  be introduced anywhere in schemas, bindings, or validators.
- Contradictory or duplicate scope identity (any scope source other than
  `metadata.scopeRef`, or conflicting scope declarations) must be rejected. A
  canonical absent/nil `metadata.scopeRef` resolving to `Platform` is not a
  contradictory or duplicate scope.

## 9. Remaining implementation evidence (not produced)

The following have **not** been done and are explicitly out of scope for this
documentation-only correction:

- No JSON schema enum was changed.
- No Go binding was changed.
- No validator was changed.
- No fixture was added or changed.
- No test was added, changed, run, or passed.

These artifacts must be synchronized only after requirements receive fresh
independent approval and a separately authorized implementation stage begins.

## 10. Dependency Contract Reconciliation summary

> **Correction note (2026-07-27).** The "Required fields" row previously stated
> that `metadata.scopeRef` was "required via `object-meta.json`". That statement
> was internally inconsistent with the canonical nil-`Platform` representation
> defined in ADH-2026-015 and clarified by ADH-2026-016 (and with SV-AC-03,
> AC-24, and CB-AC-02 of the requirements). It has been corrected to describe
> `metadata.scopeRef` as the sole scope source that is conditionally present —
> absent/nil for canonical `Platform` where permitted, non-nil and required for
> every non-Platform scope, and rejected when absent under a contract that
> disallows `Platform`. No other contract change is introduced.

| Reconciliation dimension | Evidence |
|---|---|
| Exact dependency version/files | FEATURE-0012 `sovrunn.io` `v1alpha1` grammar; files enumerated in section 3. |
| Enums | `ScopeKind` expanded additively by `ServiceInstance` (six → seven). |
| Required fields | `AuditEvent` scope identity is `metadata.scopeRef` as the **sole** scope source. It is **conditionally present, not universally required**: the canonical `Platform` form is an absent/nil `metadata.scopeRef` where the contract permits `Platform`; a non-nil `metadata.scopeRef` is required for every non-Platform scope and whenever the applicable contract does not permit `Platform` (in which case an absent `metadata.scopeRef` is rejected with the stable required-scope error). `object-meta.json` defines `metadata.scopeRef` via `$ref` to `scope-ref.json`; it does not make the field unconditionally required. |
| Cardinality | Single scope identity per record; no additional scope field. |
| Ownership | `ScopeKind` vocabulary owned as shared canonical vocabulary; `AuditEvent` scope contract owned by FEATURE-0013. |
| Serialization | Additive value only; existing serialized values unchanged. Canonical `Platform` scope remains an absent/nil `metadata.scopeRef` per FEATURE-0012 `NormalizeScope`/`CanonicalScopeIdentity` (unchanged). |
| Validation | Expanded acceptance to seven values; reject values outside the seven; reject non-`metadata.scopeRef` scope sources. Absent `metadata.scopeRef` resolves to `Platform` where permitted, else returns the stable required-scope error (absence is not automatically valid); canonical nil-`Platform` is not a contradictory/duplicate scope. |
| Security semantics | Scope does not grant authorization; `scopeRef.uid` resolves authorization (unchanged). |
| Go bindings | `internal/apimeta/scope.go` requires additive `ServiceInstance` (pending). |
| Fixtures | `ServiceInstance` positive fixtures and expanded `AuditEvent` coverage required (pending). |
| Migration classification | Additive, backward-compatible extension. |
| Approval reference | ADH-2026-015 (Approved 2026-07-27, Sanjeev Kumar); joint with ADH-2026-014. |

## 11. Escalation rule

Kiro and Cursor must not independently resolve any unapproved dependency delta
between the FEATURE-0012 inherited vocabulary and the FEATURE-0013 contract.
Any difference not explicitly approved by ADH-2026-014 or ADH-2026-015 must be
recorded as `ARCHITECTURE_DECISION_REQUIRED` and escalated to the Architecture
Owner. No mapping, alias, removal, additional scope, or applicability exclusion
may be chosen by an agent.

## 12. Traceability

- Controlling handoffs: `docs/reviews/architecture-decision-handoffs/ADH-2026-014-feature-0013-decision-record-and-auditevent-standard.md`, `docs/reviews/architecture-decision-handoffs/ADH-2026-015-feature-0013-scope-vocabulary-clarification.md`
- Architecture: `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md` (section 29)
- Requirements: `.kiro/specs/decision-object-and-auditevent-standard/requirements.md` (section 17)
- Related earlier terminology record: `docs/traceability/ADH-2026-014-terminology-reconciliation.md`
