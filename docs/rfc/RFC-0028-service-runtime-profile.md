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

## Phase 2R Migration Note

**ADH-2026-042 / DEC-0049:** The customer-facing catalog evolves from ServiceClass/ServicePlan to ServiceTypeDefinition + CloudPlatform-scoped ServiceOffering + versioned ServicePlan. ServiceRuntimeProfile bridges to immutable ServiceRequirementSet (DEC-0044). The principle that placement uses internal runtime requirements rather than customer-facing plan APIs remains valid and is strengthened by the canonical model's separation of product ownership (ServiceOffering) from runtime realization (ServiceRuntimeProfile/ServiceRequirementSet).
