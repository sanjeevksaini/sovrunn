# ADH-2026-036: Independently Recoverable Sovrunn Platform Lifecycle Boundary

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-036
- **Date:** 4 August 2026
- **Classification:** New canonical operator-facing platform lifecycle boundary
- **Canonical decision:** FCM-ADR-017
- **Affected baseline:** Platform installation, release engineering, operations, backup/recovery and pre-production roadmap

## Context

The canonical model defines implementation-neutral lifecycle management for customer services but does not fully represent the lifecycle of the Sovrunn control plane itself. Production operation requires governed installation, release compatibility, upgrade, rollback, backup, restore, credential rotation, disaster-recovery testing and decommissioning. Depending exclusively on the affected Sovrunn control plane to repair or restore itself creates a bootstrap and recovery failure mode. At the same time, Sovrunn must not absorb ownership of the OpenShift, Kubernetes, OpenStack, AWS, OCI or datacenter lifecycle it reuses.

## Decision

1. `SovrunnInstallation` is a Platform-scoped ManagedResource representing one deployed Sovrunn control plane and its desired/observed release state.
2. `SovrunnRelease` is an immutable published VersionedDefinition containing signed component references and exact version identity. Directed transition, compatibility, recovery and typed-reference semantics are governed by `ADH-2026-039`.
3. `PlatformLifecyclePolicy` is an immutable published VersionedDefinition containing release-channel, maintenance, backup, restore-verification, approval, validation and rollback requirements.
4. `PlatformLifecyclePlan` is an ImmutableRecord derived from exact installation, current/target release, policy, backup, compatibility, evidence and health inputs. It is generated before acceptance, approved by exact reference and pinned atomically when the Operation is accepted.
5. Platform actions reuse the canonical `Operation` contract. Install, upgrade, rollback, backup, restore, rotate, disaster-recovery test and decommission do not create a parallel platform-only operation framework.
6. `PlatformLifecycleAgent` is a distinct system actor, not a customer or managed resource. It executes an approved immutable plan through an external operator, GitOps or bootstrap mechanism and may operate when the affected Sovrunn control plane is unavailable.
7. Signed release artifacts, recovery configuration, backup metadata and the minimum restoration authority path remain recoverable outside the installation they protect.
8. `PlatformHealth` is a platform-operator-safe projection of installation release, readiness, component, backup, recovery and assurance status; it is not independently managed desired state.
9. Underlying infrastructure lifecycle remains externally owned. Maintenance announcement authority, ExecutionTarget state ownership, epochs/fences, service-impact actions and requalification are governed by `ADH-2026-040`.
10. A Sovrunn control plane is not modeled as a customer `ServiceInstance`, and the PlatformLifecycleAgent cannot publish releases, approve its own changes, alter policy or grant itself authority.
11. Accepted platform lifecycle intent has one logically external, independently recoverable authoritative repository. The Sovrunn API is an authenticated request and operator-safe projection surface, not a second desired-state authority.
12. The initial implementation uses a signed GitOps lifecycle repository and one signed commit as the atomic activation boundary. Git is an implementation choice, not a canonical API dependency; an equivalent lifecycle service may replace it later without changing the model.
13. A platform release or other disruptive desired-state change is requested only through an idempotent `Operation`; direct mutation of `SovrunnInstallation.desiredReleaseRef` is prohibited.
14. Operation acceptance atomically binds the accepted Operation, desired release/action, active Operation, exact immutable plan, published policy version, approval, expected installation generation, unique fencing token and durable activation event.
15. External execution is asynchronous, idempotent, checkpointed and fenced. At most one disruptive lifecycle Operation may be active per installation.
16. Desired, current and last-stable releases remain distinct. Failure creates an explicit correlated rollback, restore or forward-recovery Operation rather than silently rewriting desired state.

## Responsibility boundary

> Sovrunn core owns the implementation-neutral lifecycle, governance and customer experience of managed services. An operator-facing platform lifecycle module manages Sovrunn installations and releases through an independently recoverable lifecycle agent. Underlying infrastructure remains externally operated; Sovrunn qualifies it and safely reacts to its lifecycle events.

## Alternatives rejected

- **Treat Sovrunn itself as a ServiceInstance:** mixes platform authority and customer-service governance, creates circular recovery dependencies and exposes the wrong lifecycle contract.
- **Let the running control plane exclusively upgrade and restore itself:** fails when the API, state store, certificates or controllers are unavailable.
- **Create a separate platform operation framework:** duplicates authorization, approval, idempotency, status, audit and retention semantics already provided by Operation.
- **Manage the underlying infrastructure lifecycle:** contradicts the reuse-first product boundary and duplicates infrastructure-provider ownership.
- **Leave platform lifecycle entirely to deployment documentation:** cannot provide versioned compatibility, governed approvals, evidence, audit or repeatable recovery assurance.
- **Use the Sovrunn database as the only lifecycle authority:** makes recovery depend on the installation being recovered and cannot provide the required external bootstrap boundary.
- **Keep independent authoritative copies in the Sovrunn database and Git:** creates split-brain and ambiguous conflict resolution.
- **Patch desired release and create an Operation separately:** permits accepted intent, audit identity and executable desired state to diverge during partial failure.

## Lifecycle and deletion

An installation progresses through Installing → Ready → PlanningChange → Changing → Verifying → Ready, with Degraded, RollingBack, Restoring, Failed and Decommissioning branches. A desired release change begins with an Operation request; the exact plan and approval are established before atomic acceptance activates desired state. A release, policy or plan remains retained while referenced. Decommission requires current backup and restore evidence, terminal platform Operations and approved impact handling.

```text
authorize request
  → validate generation/current release and target compatibility
  → generate immutable candidate plan
  → approve exact plan
  → revalidate freshness and health
  → atomically accept Operation and activate desired state
  → external agent reconciles with fencing and checkpoints
  → verify and publish status
```

## Security and sovereignty

Release publication, lifecycle request, approval, planning and execution are separated authorities. The agent has a distinct workload identity and short-lived scope-limited credentials. Plans contain SecretRefs only. Backup and recovery locations, keys, administrators and artifacts participate in governance and sovereignty evaluation where applicable. Every lifecycle action is correlated across approval, plan, Operation, evidence and audit records.

## External infrastructure maintenance flow

```text
Infrastructure maintenance signal
    → ExecutionTarget Draining
    → deny new placement
    → assess affected services and relationships
    → keep / fail over / migrate / stop under policy
    → infrastructure operator completes maintenance
    → refresh evidence and re-qualify target
```

## Conformance evidence

- Upgrade succeeds through the external agent using a signed release, immutable plan, approval and scoped authority.
- Failed upgrade rolls back without relying on the affected control-plane API.
- Restore succeeds from externally recoverable artifacts while the original installation is unavailable.
- Agent compromise cannot publish a release, approve a request, change policy or broaden its own role.
- Release or policy removal is restricted while retained installation, plan, operation, evidence or audit records reference it.
- Target maintenance blocks new placement and triggers governed service-impact actions without Sovrunn claiming to perform the underlying infrastructure upgrade.
- Platform health projections disclose no recovery secrets, internal credentials or unsafe topology.
- Failure during activation leaves neither a partially accepted Operation nor a changed desired release.
- Concurrent disruptive requests, stale expected generations and stale fencing tokens are rejected deterministically.
- Replaying the same idempotency key returns the original request outcome; conflicting reuse is rejected.
- An accepted change remains recoverable and executable from the external lifecycle authority while the Sovrunn API and database are unavailable.
- Failed execution produces an explicit correlated rollback, restore or forward-recovery Operation and preserves desired/current/last-stable history.

## Reassessment triggers

A proven external platform-management standard supplies equivalent identity, signed-release, compatibility, immutable-plan, approval, recovery, audit and authority-separation semantics; or production evidence demonstrates that a canonical concept can be safely reduced without weakening independent recovery or infrastructure ownership boundaries.
