# ADH-2026-030: Personal Organization and Default Customer Governance

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-030
- **Date:** 2 August 2026
- **Classification:** Onboarding and governance simplification
- **Canonical decision:** FCM-ADR-011

## Context

An individual customer may use the cloud without a registered company. Bypassing Organization/Tenant/Project governance for individuals would create a parallel security model; forcing them to manually configure every level would import hyperscaler complexity.

## Decision

1. Individual registration creates a Personal Organization.
2. The onboarding workflow creates one default Tenant and one default Project.
3. The CloudProvider creates an enrollment in Active or Pending state according to identity, eligibility and payment policy.
4. An individual entitlement package grants only appropriate plans, regions, quotas and support.
5. The default hierarchy may remain visually hidden until team, environment or policy separation is needed.
6. Conversion to a verified legal-entity Organization uses an explicit governed migration; identity is not silently reclassified.

## Alternatives rejected

- **No Organization for individuals:** forks authorization, audit and ownership.
- **Synthetic provider-owned organization:** gives the provider incorrect ownership.
- **Mandatory manual hierarchy setup:** harms the simple customer experience.

## Consequences

- One governance model serves individuals, startups, enterprises and government.
- Default-resource naming, ownership transfer and account-recovery rules are required.
- Personal Organizations need stronger anti-abuse and payment controls appropriate to self-service.

## Security and sovereignty

Defaults do not grant broad access. Deny-by-default IAM, MFA/risk controls and least privilege apply identically. Provider activation policy decides when consumption begins.

## Conformance evidence

- Self-service active and approval-required pending flows.
- Personal Organization/default Tenant/default Project creation is idempotent.
- No implied permissions from membership.
- Legal-entity migration preserves audit and requires approval.
- Deleting the personal account safely handles active services and retained records.

## Reassessment triggers

Consumer regulation, age/identity requirements or reseller distribution requires a distinct legal contracting flow; governance resources remain reusable.
