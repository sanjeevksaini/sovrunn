# ADH-2026-022: Separate Entitlement from Quota

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-022
- **Date:** 2 August 2026
- **Classification:** Governance invariant
- **Canonical decision:** FCM-ADR-003

## Context

Catalog access and consumption limits answer different questions. Combining them creates ambiguous denials and makes commercial grants, organizational restrictions and capacity controls hard to evolve independently.

## Decision

1. `ServiceEntitlement` answers what an enrolled Organization may consume: portfolios, offerings, plans, regions and restrictions.
2. `QuotaPolicy` independently answers how much may be consumed over a scope, metric and period.
3. Both must allow a request; quota cannot create entitlement and entitlement cannot bypass quota.
4. Descendants may narrow inherited entitlement or quota but cannot enlarge a CloudProvider grant.
5. Denial reasons distinguish `NOT_ENTITLED` from `QUOTA_EXCEEDED` without leaking inaccessible resources.

## Alternatives rejected

- **One access-and-limit object:** conflates authority and quantity.
- **Quota embedded in ServicePlan:** cannot express customer, tenant or project allocation.
- **Customer-only quota:** omits provider commercial and operational limits.

## Consequences

- Two evaluator inputs are required in the authorization gate.
- Effective access and effective quota need separate customer projections.
- Revocation and limit reduction require impact analysis for running services.

## Lifecycle and deletion

Entitlements and quota policies are mutable desired state with versioned generations. Revocation or reduction is explicit and audited. Deletion restores inherited policy only after impact preview; it does not automatically terminate workloads.

## Security and sovereignty

Quota values may be commercially sensitive. Customer projections show effective limits and use but omit other customers and provider capacity.

## Conformance evidence

- Entitled and within quota succeeds.
- Entitled but over quota produces `QUOTA_EXCEEDED`.
- Within quota but not entitled produces `NOT_ENTITLED`.
- Descendant restriction succeeds; descendant expansion fails.
- Concurrent reservation cannot oversubscribe quota.

## Reassessment triggers

Introduction of reservations, prepaid commitments or shared organizational budgets that cannot be expressed as additive policy contracts without changing the invariant.
