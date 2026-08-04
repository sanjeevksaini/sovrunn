# ADH-2026-026: Reuse DecisionRecord for Sovereignty and Placement

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-026
- **Date:** 2 August 2026
- **Classification:** FEATURE-0013 downstream adoption
- **Canonical decision:** FCM-ADR-007

## Context

Sovereignty assessment and placement are governed conclusions with authority, rationale, evidence, obligations and audit requirements. FEATURE-0013 already owns the canonical DecisionRecord envelope, DecisionProfile extension, EvaluationResult and projection rules. New standalone decision envelopes would fragment accountability.

## Decision

1. `SovereigntyAssessment` is a domain concept implemented through a registered `sovereignty-assessment` DecisionProfile and typed DecisionRecord result.
2. `PlacementDecision` is a domain concept implemented through a registered `placement-decision` DecisionProfile and typed DecisionRecord result.
3. Raw evaluators may produce bounded EvaluationResults; only an authorized decision service produces authoritative DecisionRecords.
4. Placement consumes the authoritative sovereignty decision when applicable; it never treats telemetry or AI output as authority.
5. Each adopting feature declares FEATURE-0013 roles, exact contract version, profiles, risks and conformance evidence.

## Alternatives rejected

- **Standalone assessment/placement record kinds with new envelopes:** duplicates scope, authority, rationale and projection semantics.
- **Store conclusion in ServiceInstance status:** mutable status is not an auditable decision.
- **Use Operation result as decision:** execution state cannot replace authorization or placement authority.

## Consequences

- DecisionProfile and typed-result registries must support both families.
- Customer and operator explanations are authorized projections, not unrestricted canonical records.
- Decision and audit storage/retention costs must be managed without weakening accountability.

## Security and sovereignty

Canonical records minimize embedded sensitive data and reference protected evidence. AI may receive only authorized DecisionContext projections and cannot mutate, approve or execute.

## FEATURE-0013 adoption

Applicability: `APPLICABLE`

Roles:

- `EVALUATION_PRODUCER`
- `DECISION_PRODUCER`
- `DECISION_CONSUMER`
- `AUDIT_PRODUCER`
- `PROJECTION_CONSUMER`
- `OPERATION_OWNER`
- `COMPOSITE_ADOPTER`

Profiles introduced:

- `sovereignty-assessment-v1`
- `placement-decision-v1`

Contracts reused:

- `EvaluationResult`
- `DecisionProfile`
- `DecisionRecord`
- `DecisionContext`
- `AuditEvent`
- decision-linked `Operation`

Inherited risks and exact FEATURE-0013 contract versions must be selected by the implementing feature architecture. No exception to the common envelope or authority rules is approved by this ADR.

## Conformance evidence

- Registered profile and typed-result validation.
- Authoritative/advisory confusion negative tests.
- Exact scope and subject linkage.
- Evidence/context freshness and fail-closed behavior.
- Decision-linked AuditEvent and Operation obligations.
- Customer projection redaction tests.

## Reassessment triggers

FEATURE-0013 changes its extension model, or performance evidence shows that the common envelope cannot meet an approved placement latency SLO without a reviewed compatible optimization.
