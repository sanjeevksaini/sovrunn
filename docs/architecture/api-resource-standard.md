---
doc_type: architecture
title: API and Resource Standard
status: draft
phase: 2
ai_load_priority: always
ai_summary: Common API/resource conventions for Phase 2 and beyond, inspired by Kubernetes resource patterns.
---

# API and Resource Standard

## Purpose

Sovrunn resources must use consistent API and resource conventions so AI-generated features do not drift.

## Resource Shape

Preferred shape:

```yaml
apiVersion: sovrunn.io/v1alpha1
kind: ResourceKind
metadata:
  name: example
  labels: {}
  annotations: {}
spec: {}
status:
  phase: Pending
  conditions: []
```

## Required Concepts

- `apiVersion`
- `kind`
- `metadata`
- `spec`
- `status`
- `conditions`
- references using `*Ref` suffix
- validation error format
- status reason codes

## API Boundary Classification

Every resource/API must be classified as one of:

| Boundary | Meaning |
|---|---|
| Customer-facing | Safe for tenant/customer service consumption. |
| Provider/MSP-facing | Used by provider or MSP operators. |
| Internal engine-facing | Used by placement, policy, scoring, or operation internals. |
| Plugin-facing | Used by plugin contracts and execution boundaries. |

Customer-facing APIs must not expose low-level IaaS complexity by default.
