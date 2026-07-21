---
doc_type: rfc
title: RFC-0021 Reuse-First Architecture
status: draft
phase: 2
ai_load_priority: high
ai_summary: Establishes reuse-before-build as mandatory feature design policy.
---

# RFC-0021: Reuse-First Architecture

## Summary

Sovrunn must prefer mature open-source and open-standard foundations before custom implementation.

## Decision

Every feature must classify work as:

```text
Reuse / Wrap / Extend / Build
```

The mandatory assessment format, controlled vocabularies, mitigation
fields, and validation rules are defined by the canonical standard:

`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

This RFC does not redefine that schema. FEATURE-0011 consolidates and
versions the standard; later features consume it by reference.

Sovrunn builds governance, decisioning, placement, orchestration, audit, evidence, AI-readable context, and customer/provider experience.

Sovrunn reuses or wraps mature systems for policy, IAM, secrets, workflow, Kubernetes control loops, PostgreSQL operations, observability, backup, networking, traffic, AI runtime, and persistence.
