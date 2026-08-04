# ADH-2026-024: Separate Physical Geography from Sovereignty

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-024
- **Date:** 2 August 2026
- **Classification:** Topology and sovereignty boundary correction
- **Canonical decision:** FCM-ADR-005

## Context

Country is essential to data residency but does not describe a full hosting location and does not prove sovereignty. Sovereignty also depends on data processing, backups, key custody, administrative control, control planes, dependencies and evidence. Provider-prefixed geography incorrectly suggests that the provider owns geographic facts.

## Decision

1. Use implementation-independent `HostingLocation`, `Datacenter` and `FaultDomain` resources.
2. Country, subdivision, locality and optional physical classifications are attributes of HostingLocation.
3. Physical facts are separate from `SovereigntyProfile`, `RegulatoryPolicyBundle`, `SovereigntyFactSet`, `EvidenceRecord` and sovereignty decisions.
4. Topology is a graph of provider-registered facts and typed relationships, not proof of legal compliance or network connectivity.
5. New geography and sovereignty dimensions enter through governed registries and versioned schemas.

## Alternatives rejected

- **Country as the topology root:** too coarse for latency, datacenter choice and resilience.
- **ProviderLocation/ProviderDatacenter:** confuses fact registration with ownership.
- **Sovereignty boolean on location:** misleading and impossible to maintain as law/evidence changes.

## Consequences

- FEATURE-0014 topology must be superseded with migration mapping.
- Customer regions may map to several locations and targets.
- Legal/policy interpretation and evidence freshness become explicit dependencies.

## Security and sovereignty

Sensitive facility and fault-domain detail has operator/internal boundaries. Customer views expose only approved location and outcome summaries. Missing evidence yields Unknown/Indeterminate, never inferred compliance.

## Conformance evidence

- Same location produces different outcomes under two sovereignty profiles.
- Country match with failed key custody is rejected.
- Stale evidence produces indeterminate assessment.
- Common topology ancestry implies neither connectivity nor sovereignty.

## Reassessment triggers

A recognized external standard supplies a stable, extensible hosting-and-sovereignty model that can replace the Sovrunn contracts without losing evidence, versioning or decision linkage.
