# ADH-2026-038: Evidence-Backed Platform Sovereignty and Dependency Composition

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-038
- **Date:** 4 August 2026
- **Classification:** Extension of sovereignty and platform-lifecycle boundaries
- **Canonical decision:** FCM-ADR-019

## Context

An in-country control-plane runtime is insufficient to establish platform sovereignty when lifecycle state, artifacts, keys, backups, telemetry or administrators remain under unassessed external control.

## Decision

1. Platform sovereignty reuses `SovereigntyProfile`, `RegulatoryPolicyBundle`, `SovereigntyFactSet`, `EvidenceRecord`, `ProfileAssignment` and the sovereignty `DecisionRecord` profile; it does not create a second assessment engine.
2. The assessment subject is `SovrunnInstallation`.
3. `PlatformDependencySnapshot` is an immutable assessment input enumerating every dependency capable of controlling, observing, changing, decrypting or recovering the installation.
4. Required dependency classes include control-plane runtime/state, lifecycle authority and agent, artifact registry, secret/KMS/HSM, backup/recovery, identity/privileged access, observability/audit, network/DNS, support path and external AI/automation.
5. The applicable profile evaluates hosting, metadata, legal/organizational control, cryptographic custody, administration, software supply chain, recovery, telemetry and external processing.
6. Missing required dependencies or stale/missing mandatory evidence produce `Indeterminate`, never success.
7. A service placement requiring sovereign operation is ineligible when the selected installation's platform sovereignty assessment is expired, Indeterminate or NotSatisfied.
8. A material dependency, legal, location, administrator, key, release, telemetry or evidence change triggers reassessment.
9. Existing services are not automatically deleted after adverse reassessment; policy governs restrict, drain, remediate, migrate, fail over or stop actions.

## Responsibility separation

The CloudPlatform owner or approved authority publishes minimum profiles; qualified legal/compliance authority publishes policy bundles; providers declare configuration; approved collectors produce facts/evidence; an authorized decision service evaluates; separate approvers govern exceptions; provider lifecycle operators remediate.

## Core statement

> A SovrunnInstallation is sovereign only when the installation and every dependency capable of controlling, observing, changing, decrypting or recovering it satisfy the applicable evidence-backed SovereigntyProfile.
