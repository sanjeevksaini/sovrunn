# ADH-2026-031: Extend Countries, Sectors, Services and Backends Through Data and Contracts

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-031
- **Date:** 2 August 2026
- **Classification:** Extensibility and anti-drift invariant
- **Canonical decision:** FCM-ADR-012

## Context

Sovrunn must support changing sovereignty rules, new industries, additional services and heterogeneous infrastructure without adding country conditionals, industry-specific core kinds or vendor SDK types to the control-plane model.

## Decision

1. Countries and physical classifications enter through HostingLocation facts and governed registries.
2. Sovereignty changes enter through versioned SovereigntyProfiles, RegulatoryPolicyBundles, fact/evidence schemas and registered decision profiles.
3. Industries enter through ServicePortfolios, GovernanceProfiles and service plugins.
4. New service families enter through immutable ServiceTypeDefinitions, compatible runtime profiles and versioned plugins/adapters; service-specific fields never enter core ServiceInstance, decision or operation schemas.
5. Backends enter through versioned adapters, ExecutionTargets and normalized facts; ExecutionTargets may represent infrastructure, external APIs, edge or future federation endpoints.
6. Extensions require a namespaced owner, versioned schema, declared boundary/classification, bounded size, validation and compatibility rules.
7. Governed data, metering/charging and external supply/federation are reserved domain extension boundaries, not mandatory MVP resources. Promotion requires a separate approved handoff and end-to-end reference flow.
8. An extension used by core decisions or three independent consumers is reviewed for promotion to a typed core contract.
9. Arbitrary property bags and country/vendor/service conditionals in orchestration are prohibited.

## Alternatives rejected

- **Hard-code each country/industry/vendor:** produces combinatorial drift.
- **Unrestricted extensions:** creates opaque semantics and security bypass.
- **Lowest-common-denominator schema:** cannot express sovereign evidence or service lifecycle accurately.
- **Add a core resource/workflow for every service family:** creates combinatorial product drift and couples the platform to today's catalog.

## Consequences

- Strong schema registries, version governance and compatibility tests are required.
- Legal/policy content has a separate qualified publication lifecycle.
- Plugin and adapter ecosystems can grow without changing customer APIs.
- ServiceTypeDefinition provides a governed extension surface for parameters, actions, projections, bindings, meters and upgrades.

## Security and sovereignty

Extensions inherit authorization, redaction, residency, provenance and audit requirements. No extension can weaken core policy, create authority or expose secrets.

## Conformance evidence

- Add a new country policy bundle without core-code change.
- Add a second industry portfolio reusing existing offerings.
- Add a second backend adapter without customer-schema change.
- Add Kafka or DNS through a new ServiceTypeDefinition without modifying core ServiceInstance, decision or operation schemas.
- Unknown/unregistered extension rejection.
- Core-decision dependency on unpromoted extension rejected.

## Reassessment triggers

Repeated promotion pressure indicates missing core vocabulary; extension-governance overhead exceeds demonstrated benefit; or an external standard becomes suitable for direct adoption.
