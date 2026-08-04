# ADH-2026-020: CloudProvider Identity and Scope

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-020
- **Date:** 2 August 2026
- **Classification:** Breaking alpha terminology and scope correction
- **Affected baseline:** FEATURE-0012, FEATURE-0013, FEATURE-0014
- **Canonical decision:** FCM-ADR-001

## Amendment

`ADH-2026-037` and `ADH-2026-041` refine this decision by separating the customer-facing `CloudPlatform` from the operational `CloudProvider`. The decision below is the amended canonical form; the earlier direct product/enrollment interpretation is superseded.

## Context

The repository currently uses `Provider` for both cloud-product ownership and infrastructure operation. The generic word also conflicts with identity providers, service providers, plugin providers and ordinary prose. The canonical model must classify those roles rather than rename both into one new kind.

## Decision

1. `CloudProvider` is the canonical resource kind for the cloud-supply/operator boundary; `CloudPlatform` is the separate customer product/governance boundary.
2. The scope vocabulary becomes `Platform`, `Organization`, `OrganizationUnit`, `Tenant`, `Project`, `CloudPlatform`, `CloudProvider`.
3. A `CloudProvider` is Platform-scoped or Organization-scoped. The scope represents administrative governance, not customer ownership.
4. Customer Organizations remain independent and consume a CloudPlatform through `CloudEnrollment`; provider supply is joined through `CloudProviderParticipation`.
5. Generic `Provider` is not retained as a simultaneous alias in active contracts.

## Alternatives rejected

- **Keep `Provider`:** too ambiguous for a long-lived public model.
- **Name it `InfrastructureProvider`:** incorrectly narrows a CloudProvider that offers application, data, AI and IaaS services.
- **Rename every Provider record to CloudProvider:** preserves the original conflation of product ownership and infrastructure operation.

## Consequences

- Positive: precise language, stable customer/provider separation and room for other provider concepts.
- Negative: coordinated breaking change to schemas, fixtures, references, routes, documentation and conformance checks.
- No runtime compatibility alias is introduced unless migration evidence proves one is necessary and time-bounded.

## Migration

1. Approve a FEATURE-0012 amendment replacing the scope enum atomically.
2. Classify `Provider` records into CloudPlatform ownership and/or CloudProvider operation while APIs are alpha.
3. Update FEATURE-0013 scope fixtures without changing DecisionRecord semantics.
4. Supersede FEATURE-0014 `Provider` identity and migrate stored objects preserving UID lineage where technically safe.
5. Reject mixed old/new scope graphs after the migration window.

## Security and sovereignty

Authorization continues to resolve immutable `scopeRef.uid`. The rename grants no authority. Cross-provider references remain denied by default and must not disclose target existence.

## Conformance evidence

- Exactly seven final scope kinds.
- Old `Provider` rejected after migration.
- Platform- and Organization-scoped CloudProvider positive fixtures.
- Mixed-name and cross-scope negative fixtures.
- DecisionRecord, AuditEvent and Operation scope compatibility tests.

## Reassessment triggers

A demonstrable need for a generic supply-party abstraction shared by three independent resource families, or a stable API compatibility obligation that makes immediate removal unsafe.
