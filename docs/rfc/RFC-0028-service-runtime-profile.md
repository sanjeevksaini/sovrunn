---
doc_type: rfc
title: RFC-0028 ServiceRuntimeProfile
status: draft
phase: 2
ai_load_priority: high
ai_summary: RFC for ServiceRuntimeProfile as bridge between customer-facing plans and infrastructure capability requirements.
---

# RFC-0028: ServiceRuntimeProfile

## Summary

`ServicePlan` is customer-facing. `ServiceRuntimeProfile` is the internal bridge that defines runtime requirements, capabilities, and placement constraints.

## Decision

Placement must use ServiceRuntimeProfile rather than encoding infrastructure requirements directly into customer-facing ServicePlan APIs.
