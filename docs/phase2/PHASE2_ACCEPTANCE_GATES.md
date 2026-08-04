---
doc_type: acceptance_gates
title: Phase 2R Acceptance Gates
status: approved
phase: 2R
ai_load_priority: always
ai_summary: Mandatory gates for each Phase 2R feature under canonical model architecture.
---

# Phase 2R Acceptance Gates

Every Phase 2R feature must pass these gates before merge.

## 1. Architecture Contract Gate

- Feature purpose is clear.
- Dependencies are listed.
- Phase boundary is respected.
- Reuse assessment is complete.
- Customer, operator, internal-engine, adapter, plugin, or governance API boundary is identified.
- Non-goals are explicit.
- Canonical model terminology is used (no generic Provider, no ResourcePool/ProviderCapability as active concepts).
- Seven-scope vocabulary is respected.

FEATURE-0012-and-later resource/API contracts must also conform to:

`docs/architecture/api-resource-standard.md`

The architecture gate checks resource profile, allowed scope, boundary classification, ownership, implementation neutrality, compatibility, and reassessment coverage.

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
- adapter boundaries are preserved,
- no mandatory ResourcePool or ProviderCapability in active contracts,
- no generic Provider as combined owner/operator concept,
- no six-scope vocabulary in active authorities,
- ServiceClass has no active canonical authority,
- ExecutionTarget is the placement boundary,
- EffectiveGovernanceContext is the resolved governance term,
- published definitions are immutable by version,
- customer APIs contain no provider-native objects, raw secrets, or protected handles.

## 4. Quality Gate

- Tests cover happy paths and failure paths.
- Structured errors follow API contract.
- Observability fields are present.
- Race conditions are tested where relevant.

## 5. Human Approval Gate

- Feature gate script passes.
- Human or approved reviewer confirms acceptance criteria.
- Feature may proceed to next in sequence.

## 6. VS-000 Conformance Gate (FEATURE-0026 only)

- Cross-feature orchestration passes positive, denial, stale-input, authorization, idempotency, fencing, redaction, and no-side-effect conformance.
- VS-000 definition of done is met.
