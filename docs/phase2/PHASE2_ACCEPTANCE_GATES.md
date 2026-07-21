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
- Customer/provider/internal/plugin API boundary is identified.
- Non-goals are explicit.

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
