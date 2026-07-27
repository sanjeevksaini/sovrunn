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

### 6.2 Canonical seven-value ScopeKind and AuditEvent

- `AuditEvent` supports Platform, Organization, OrganizationUnit, Tenant, Project, Provider, and ServiceInstance through `metadata.scopeRef` as its sole logical and serialized scope authority.
- The canonical `Platform` form is an absent/nil `metadata.scopeRef` (FEATURE-0012 `NormalizeScope`/`CanonicalScopeIdentity`): an absent `scopeRef` resolves deterministically to `Platform` where the contract permits `Platform`, and returns the stable required-scope error where it does not (absence is not automatically valid).
- No parallel scope enum, presence flag, `AuditScope` type, discriminator, alias, or independently mutable scope field is introduced.
- Compatibility fixtures exist for all seven AuditEvent scopes — a positive `Platform` fixture using canonical absent/nil `metadata.scopeRef`, positive non-Platform fixtures using `metadata.scopeRef`, and a negative absence fixture asserting rejection when `Platform` is disallowed — with regression coverage for the six pre-existing FEATURE-0012 ScopeKind values.
- FEATURE-0012 conformance is not regressed.

### 6.3 Contract-only scope

- No production runtime implementation exists.
- No persistence, workflow, or external service implementation exists.
- Schemas, validation, conformance fixtures, and documentation only.

### 6.4 Controlling references

- ADH-2026-014 and ADH-2026-015 are the joint controlling handoffs. ADH-2026-015 supersedes only the earlier six-scope `AuditEvent` statement; all other ADH-2026-014 decisions remain controlling.
- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md` is the canonical architecture.
- `docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md` records the value-by-value compatibility evidence.
- AD-001 through AD-044 are preserved.
- Matrix E F13-R01 through F13-R31 are preserved with architecture-stage treatment.

### 6.5 Downstream adoption

- One normative downstream adoption contract exists.
- One lightweight gate check validates adoption section presence.
- No separate manifests, registries, indexes, or approval workflows are introduced.
