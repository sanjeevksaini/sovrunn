# ADH-2026-023: Industry and Domain Clouds as ServicePortfolios

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-023
- **Date:** 2 August 2026
- **Classification:** Product-model boundary
- **Canonical decision:** FCM-ADR-004

## Context

Sovrunn may power government, financial-services, healthcare, AI and other domain clouds. Creating a provider or core resource kind for every industry would duplicate the platform and hard-code market-specific semantics into the control plane.

## Decision

1. An industry or domain cloud is a versioned `ServicePortfolio`.
2. A portfolio composes ServiceOfferings, mandatory GovernanceProfiles, SovereigntyProfiles, evidence requirements, placement profiles and support terms.
3. Service-specific behavior belongs in management plugins; sector interpretation belongs in versioned policy/profile data.
4. Multiple portfolios may reuse the same offering and runtime profile.

## Alternatives rejected

- **`BankingCloud`, `GovernmentCloud`, etc. kinds:** unbounded type growth and duplicated workflow.
- **Separate Sovrunn deployments for every industry:** prevents shared-core reuse.
- **Catalog folders only:** insufficient to package mandatory governance and assurance.

## Consequences

- One core supports many products and jurisdictions.
- Portfolio publication requires versioning and validation of all referenced profiles.
- Marketing labels cannot imply compliance without a valid assessment/evidence model.

## Security and sovereignty

Portfolio membership does not automatically satisfy its controls. Every ServiceInstance resolves effective governance and receives evidence-backed decisions.

## Conformance evidence

- Government and financial portfolios reuse Managed PostgreSQL.
- Mandatory profile omission fails publication.
- Published portfolio version is immutable.
- Customer sees only portfolios granted by enrollment entitlements.

## Reassessment triggers

A domain requires a fundamentally different ownership, tenancy or lifecycle model that cannot be represented through offerings, profiles, plugins and policy bundles.
