# ADH-2026-039: Directed Release Compatibility, Recovery Semantics and Lifecycle Reference Contracts

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-039
- **Date:** 4 August 2026
- **Classification:** Platform lifecycle contract completion
- **Canonical decision:** FCM-ADR-020

## Context

Release versions alone cannot prove that an installation can upgrade or recover safely. Compatibility is directional and spans stored schemas, APIs, plugins, adapters, decision profiles, clients, lifecycle agents and backups. Existing lifecycle references also require typed semantic ownership.

## Decision

### Directed compatibility

1. Every permitted transition is an explicit directed `ReleaseCompatibilityContract` edge from one exact SovrunnRelease to another. Semantic-version proximity never implies compatibility.
2. Missing, expired, withdrawn or unverified compatibility means the transition is prohibited.
3. The contract covers API/resource read-write compatibility, stored-schema migration, component protocol compatibility, plugin/adapter ranges, DecisionProfile/result schemas, CLI/portal/SDK ranges, PlatformLifecycleAgent and plan-schema ranges, event/audit schemas, and backup/restore formats including key requirements.
4. An accepted PlatformLifecyclePlan pins the exact compatibility contract and every referenced definition/evidence version.

### Recovery semantics

`RecoveryMode` has three executable values:

| Mode | Meaning |
|---|---|
| `InPlaceSupported` | The exact compatibility edge permits a verified direct return to the previous release without restoring platform state |
| `RestoreRequired` | Direct downgrade is forbidden; recovery restores the verified pre-change backup using the contractually compatible release and repository |
| `ForwardRecoveryOnly` | Returning to the previous state is unsupported; the only permitted recovery is a validated forward repair or superseding release |

`rollbackProhibitedAfterMigration` is a separate irreversible-boundary constraint, not a fourth recovery action. It identifies the exact plan step/checkpoint after which in-place rollback is illegal and must be paired with `RestoreRequired` or `ForwardRecoveryOnly`.

The plan declares recovery behavior for at least pre-change, pre-irreversible-boundary and post-irreversible-boundary failure. It may not invent a recovery path absent from the compatibility contract.

### Typed lifecycle references

| Contract | Classification and authority | Required semantics |
|---|---|---|
| `ComponentManifest` | Published immutable VersionedDefinition; release publisher | Exact component identities, artifact digests, versions, dependencies, rollout ordering, required configuration/schema versions and health checks |
| `ArtifactProvenanceRecord` | Immutable EvidenceRecord profile; approved provenance collector/publisher | Source revision, builder identity, build attestation, signatures, SBOM, digest binding, vulnerability-policy result and validity |
| `MaintenanceWindow` | Published immutable VersionedDefinition; policy/schedule publisher | Time zone, recurrence or bounded interval, duration, notice, blackout periods, emergency behavior and approval requirements |
| `RecoveryRepository` | ManagedResource registering an external recovery system; recovery administrator/spec and qualification controller/status | Repository type, geography, sovereignty, encryption, integrity, retention, restore-format support, SecretRefs and current qualification; never raw credentials |
| `ValidationGateDefinition` | Published immutable VersionedDefinition; gate publisher | Evaluator contract, typed inputs, phase, success/failure semantics, timeout, criticality, retry, evidence output and fail-open prohibition where mandatory |
| `ReleaseCompatibilityContract` | Published immutable VersionedDefinition; release compatibility authority | Directed release edge and all compatibility/recovery rules described above |

## Acceptance gates

Before Operation acceptance, Sovrunn verifies exact current/target releases, compatibility edge, manifest/provenance, agent, plugins/adapters, clients subject to policy, backup format, RecoveryRepository qualification, required pre-change backup, validation gates and recovery path. Any unknown mandatory dimension denies activation.

## Conformance

- In-place rollback succeeds only on an edge explicitly permitting it.
- RestoreRequired proves restore from the pinned backup and repository while the affected API is unavailable.
- ForwardRecoveryOnly rejects downgrade and accepts only a compatible superseding plan.
- Crossing the irreversible checkpoint changes the permitted recovery path deterministically.
- An incompatible plugin, agent, decision profile or backup format denies activation before execution.
