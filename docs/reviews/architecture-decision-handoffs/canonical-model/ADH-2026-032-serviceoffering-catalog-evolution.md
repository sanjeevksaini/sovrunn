# ADH-2026-032: ServiceOffering and CloudPlatform-Scoped Catalog Evolution

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-032
- **Date:** 2 August 2026
- **Classification:** Phase 1 catalog replacement/migration
- **Canonical decision:** FCM-ADR-013
- **Affected baseline:** FEATURE-0006, FEATURE-0007, FEATURE-0008, FEATURE-0010

## Context

Phase 1 implements a global mutable `ServiceClass` and `ServicePlan`. The final product model requires each CloudPlatform to publish versioned offerings and plans with product-specific availability, profiles, commercial terms and runtime contracts. Participating CloudProviders publish realization eligibility separately. Retaining a global ServiceClass as the customer product authority would prevent independent CloudPlatform catalogs and conflict with enrollment and entitlement.

## Decision

1. `ServiceOffering` supersedes `ServiceClass` as the customer-consumable CloudPlatform product definition.
2. `ServiceTypeDefinition` is the separate reusable implementation-neutral service-family contract. It defines bounded parameter, lifecycle-action, safe-status, binding, usage-meter, runtime-compatibility and upgrade schemas.
3. `ServiceOffering` and `ServicePlan` are CloudPlatform-scoped VersionedDefinitions; ServiceTypeDefinition may be Platform-, CloudPlatform- or approved-authority-scoped.
4. Every ServiceOffering references exactly one ServiceTypeDefinition version. A ServiceInstance references an exact ServicePlan version; the plan references exactly one ServiceOffering version.
5. Existing ServiceClass data is migration input for both product offerings and, where semantics are reusable, ServiceTypeDefinitions; it is not a permanent alias.
6. The Phase 1 plugin `Capability` registry remains a plugin lifecycle-support declaration until separately reworked. It is not `ProviderCapability` and does not prove target compatibility.
7. Customer catalog projections show only offerings/plans permitted by active enrollment and entitlement.

## Alternatives rejected

- **Keep global ServiceClass:** cannot represent independent CloudPlatform catalogs.
- **Make ServiceOffering an alias:** leaves two public names for one concept.
- **Embed offering fields in ServicePlan:** duplicates product identity and lifecycle.
- **Use ServiceOffering as both product and technical service-type contract:** prevents several providers from packaging the same reusable service semantics independently.

## Consequences

- Phase 1 routes, types, fixtures, demo and plugin references require alpha migration.
- Existing ServicePlan identity changes from global/composite to CloudPlatform-scoped versioned identity.
- SDE and future service roadmap entries use ServiceOffering terminology.

## Migration

1. Extract reusable technical semantics from each existing ServiceClass into a ServiceTypeDefinition version.
2. Create one ServiceOffering version for each cloud product under the selected bootstrap CloudPlatform and reference the appropriate ServiceTypeDefinition.
3. Convert ServicePlans and ServiceInstances to typed exact-version references.
4. Convert plugin ServiceClassRefs to ServiceTypeDefinition compatibility references; product-facing references use ServiceOffering.
5. Provide a bounded migration report for unresolved or ambiguous global classes.
6. Remove old create/update routes after the approved alpha migration window.

## Security and sovereignty

Catalog reads are enrollment/entitlement filtered. Provider-private plans and unpublished definitions must not be discoverable across scopes.

## Conformance evidence

- Two CloudPlatforms publish distinct PostgreSQL offerings without identity collision.
- Both offerings may reuse one approved PostgreSQL ServiceTypeDefinition without sharing commercial identity or terms.
- Published versions are immutable.
- Legacy ServiceClass migrates deterministically.
- Customer cannot list or reference an unentitled offering.
- Plugin capability does not become placement authority.

## Reassessment triggers

A cross-CloudPlatform marketplace requires additional publication, trust or supply-chain semantics that cannot be represented by shared ServiceTypeDefinitions and product-specific offerings. Such a marketplace contract must not replace CloudPlatform product identity or customer enrollment.
