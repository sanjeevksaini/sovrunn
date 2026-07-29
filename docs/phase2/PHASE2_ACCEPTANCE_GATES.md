---
doc_type: acceptance_gates
title: Phase 2 Acceptance Gates
status: draft
phase: 2
ai_load_priority: always
ai_summary: Mandatory gates for each Phase 2 feature under the Me + AI chain execution model.
---

# Phase 2 Acceptance Gates

Every Phase 2 feature must pass these gates before merge.

## 1. Architecture Contract Gate

- Feature purpose is clear.
- Dependencies are listed.
- Phase boundary is respected.
- Reuse assessment is complete.
- Customer, operator, internal-engine, adapter, plugin, or governance API boundary is identified.
- Non-goals are explicit.

FEATURE-0012-and-later resource/API contracts must also conform to:

`docs/architecture/api-resource-standard.md`

The architecture gate checks resource profile, allowed scope, boundary classification, ownership, provider neutrality, compatibility, and reassessment coverage.

## 2. Reuse Gate

Feature contract must answer:

```text
Can this be reused from mature OSS?
Should Sovrunn reuse, wrap, extend, or build?
What adapter boundary prevents recoding later?
```

FEATURE-0011-and-later assessments must conform to the canonical standard:

`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

The strict feature gate enforces structural and consistency validation
(RA-S* / RA-C*). Human semantic review remains required for architecture
approval. This gate section does not redefine the assessment field schema.

## 3. Architecture Drift Gate

Check:

- no provider-specific hardcoding in core,
- no Kubernetes-only assumption in core,
- no PostgreSQL lifecycle logic in core placement engine,
- no custom policy engine embedded in handlers,
- no raw secret storage,
- no customer-facing IaaS leakage,
- decision object is explainable,
- audit behavior is defined,
- adapter boundaries are preserved.

## 4. Quality Gate

- unit tests pass,
- integration tests pass where applicable,
- lint passes,
- gosec passes,
- race tests pass where applicable,
- generated docs are consistent with source-of-truth docs.

## 5. Human Acceptance Gate

The architecture owner approves with ChatGPT support before merge.

## 6. FEATURE-0013 Specific Gate Criteria

FEATURE-0013 (Decision Record and AuditEvent Standard) must additionally satisfy:

### 6.1 Terminology reconciliation

- All normative documents use `DecisionRecord` (not `DecisionObject`) for the common decision envelope.
- A terminology reconciliation traceability record exists.

### 6.2 Canonical scope and subject model

- `AuditEvent` supports FEATURE-0012's six governance scopes: Platform, Organization, OrganizationUnit, Tenant, Project, and Provider.
- The canonical `Platform` form is an absent/nil `metadata.scopeRef` (FEATURE-0012 `NormalizeScope`/`CanonicalScopeIdentity`): an absent `scopeRef` resolves deterministically to `Platform` where the contract permits `Platform`, and returns the stable required-scope error where it does not (absence is not automatically valid).
- No parallel scope enum, presence flag, `AuditScope` type, discriminator, alias, or independently mutable scope field is introduced.
- `ServiceInstance` is represented through a typed `subjectRef`, normally with Project as `metadata.scopeRef`; it is rejected as a `ScopeKind`.
- Compatibility fixtures cover all six FEATURE-0012 scopes, the canonical Platform form, invalid absence where Platform is disallowed, rejection of parallel scope sources, rejection of `ServiceInstance` as scope, and acceptance of a ServiceInstance subject within Project scope.
- FEATURE-0012 conformance is not regressed.

### 6.3 Contract-only scope

- No production runtime implementation exists.
- No persistence, workflow, or external service implementation exists.
- Schemas, validation, conformance fixtures, and documentation only.

### 6.4 Consolidated controlling reference

- `ADH-2026-017` is the approved single replacement handoff for `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`.
- ADH-2026-014, ADH-2026-015, and ADH-2026-016 are historical provenance and must not be loaded as separate downstream instructions.
- FEATURE-0013 requirements, design, tasks, implementation, conformance fixtures, Matrix E evidence, and final feature approval are complete through merged PR #15 (2026-07-29).
- AD-001 through AD-045 are preserved; AD-045 owns the public error binding and
  closed violation-code registry.
- Matrix E F13-R01 through F13-R31 are preserved with architecture-stage treatment.

### 6.5 Consolidated contract boundary criteria

- `DecisionRecord` uses `metadata.scopeRef` as its sole logical and serialized scope authority; no top-level `DecisionRecord` `scopeRef` or parallel scope source exists.
- Sensitivity classification uses the closed ordered provider-neutral vocabulary `PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED`; jurisdiction-specific labels are versioned profile mappings; profiles permitting captured output declare a sensitivity ceiling, content-category rules, and maximum fields/bytes/depth/value-types.
- Security validation is bounded to structural conformance; no secret/credential/PII/malware/DLP/content-scanning engine is introduced; comprehensive semantic content scanning is a later approved feature.
- Algorithm-agile carrier fields and structural trust metadata are contract-now; canonicalization algorithm/profile, digest-covered fields, signature algorithm, and cryptographic services remain DEFERRED under ADR-F13-002; RFC 8785 is illustrative only.
- Architecture section 17 scenarios map one-to-one to stable IDs `F13-CF-01` through `F13-CF-28`; coverage is counted by scenario ID with an explicit coverage-matrix entry per scenario.
- FEATURE-0013 final gates passed before PR #15 merge; ADH-2026-017 remains the controlling consolidated architecture reference.

### 6.6 Downstream adoption

- One normative downstream adoption contract exists.
- One lightweight gate check validates adoption section presence.
- No separate manifests, registries, indexes, or approval workflows are introduced.
