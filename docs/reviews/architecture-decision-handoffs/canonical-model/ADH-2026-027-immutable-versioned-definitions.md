# ADH-2026-027: Immutable Published Definitions

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-027
- **Date:** 2 August 2026
- **Classification:** Shared versioning invariant
- **Canonical decision:** FCM-ADR-008

## Context

Service plans, runtime profiles, governance profiles, sovereignty profiles, policy bundles and placement profiles directly affect customer promises and decisions. Mutating a published definition would silently change the meaning of existing services and historical records.

## Decision

1. These contracts use the FEATURE-0012 VersionedDefinition profile.
2. Draft versions are mutable by one authorized publisher.
3. Published versions are immutable. A change creates a new version with explicit predecessor/supersession linkage.
4. ServiceInstances, contexts, decisions and plans reference exact versions.
5. Retirement prohibits new use but preserves existing references until migration or end of retention.
6. Deletion is restricted while any live or retained object references the version.

## Alternatives rejected

- **Mutable latest object:** destroys reproducibility.
- **Copy every definition into every instance:** creates duplication and inconsistent provenance.
- **Semantic version string without immutable identity:** insufficient enforcement.

## Consequences

- Publishers need draft, publish, deprecate, retire and supersede workflows.
- Storage retains multiple versions.
- Upgrade/migration becomes explicit and auditable.

## Security and sovereignty

Publication requires authorization, validation and audit. RegulatoryPolicyBundle versions additionally require qualified approval, source references, effective dates and integrity protection.

## Conformance evidence

- Published mutation rejected.
- New version and supersession accepted.
- Retired version denied for new instances but readable for existing records.
- Deletion denied while referenced.
- Historical decision resolves exact definition versions.

## Reassessment triggers

None for semantic immutability. Storage optimization may change physical retention only if canonical identity and reproducibility remain intact.
