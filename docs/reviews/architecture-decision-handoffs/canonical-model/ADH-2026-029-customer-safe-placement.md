# ADH-2026-029: Customer-Safe ServicePlacement Projection

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-029
- **Date:** 2 August 2026
- **Classification:** Customer/security boundary
- **Canonical decision:** FCM-ADR-010

## Context

Customers need evidence of where services and data are hosted and whether location, HA and sovereignty outcomes are met. Canonical placement decisions and provider topology contain protected details that should not cross the customer boundary.

## Decision

1. `ServicePlacement` is a customer-facing immutable projection produced by a dedicated projection controller.
2. It may expose ServiceRegion, approved datacenter labels, component roles, resilience results, sovereignty outcome, assessment time/validity and current lifecycle summary.
3. It excludes credentials, provider-native handles, sensitive fault-domain topology, rejected confidential candidates and unrestricted evidence.
4. ServiceInstance status points to the current ServicePlacement record; changes create a new projection record.
5. Canonical DecisionRecords remain authoritative. ServicePlacement never becomes a second decision authority.

## Alternatives rejected

- **Expose PlacementDecision directly:** violates boundary and projection rules.
- **Only show region label:** insufficient for regulated customers and assurance.
- **Mutable placement fields in ServiceInstance status:** loses historical traceability.

## Consequences

- Projection schemas and redaction policy require explicit versioning and tests.
- Different customer/auditor views may require distinct projections.
- Safe explanation may be less detailed than provider diagnostics by design.

## Security and sovereignty

Field-level classification, purpose authorization and no-existence disclosure apply. Evidence is summarized or referenced through an authorized view.

## Conformance evidence

- Customer sees allowed location and outcome fields.
- Secret/provider-native/sensitive candidate leakage tests.
- Projection links exact decision and assessment.
- Reassessment creates a new projection and preserves prior record.
- Cross-project access denied safely.

## Reassessment triggers

A regulator or contracted customer requires a new purpose-specific projection with additional evidence; the canonical record remains unchanged.
