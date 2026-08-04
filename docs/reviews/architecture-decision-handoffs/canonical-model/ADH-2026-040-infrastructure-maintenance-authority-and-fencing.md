# ADH-2026-040: External Infrastructure Maintenance Authority, Drain State and Fencing

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-040
- **Date:** 4 August 2026
- **Classification:** ExecutionTarget lifecycle contract completion
- **Canonical decision:** FCM-ADR-021

## Context

Sovrunn does not own underlying infrastructure maintenance, but it must safely stop new placement, protect active operations and reassess affected services. A generic signal and a writable `Draining` status are insufficient because authority, state transitions and stale-plan fencing would be ambiguous.

## Decision

1. `InfrastructureMaintenanceNotice` is a CloudProvider-scoped ManagedResource for scheduled, emergency, updated or cancelled maintenance affecting one or more ExecutionTargets.
2. Only an authenticated `InfrastructureOperator` role delegated by the target's CloudProviderParticipation or a trusted adapter identity may create/update the notice. The CloudPlatform owner may impose a separately governed restriction through existing governance and decision contracts but cannot impersonate the infrastructure operator or clear provider maintenance.
3. The notice controller validates authority, target participation, expected target generation, maintenance window and idempotency key. It owns notice status.
4. Only the `ExecutionTargetLifecycleController` writes ExecutionTarget availability status. Providers, owners and customers never patch target status directly.
5. Qualification and availability are separate axes.

### Availability states

```text
Available
  → DrainPending
  → Draining
  → InMaintenance
  → Requalifying
  → Available
```

Allowed branches:

- `Available → Unavailable` for an unplanned outage;
- any active state → `Restricted` when policy/evidence disqualifies use;
- `DrainPending` or `Draining → Requalifying` when maintenance is cancelled after any placement impact;
- `InMaintenance → Degraded` when infrastructure reports completion but validation fails;
- `Requalifying → Restricted`, `Degraded` or `Unavailable` when gates fail.

Returning to Available always requires successful requalification; cancellation or provider completion alone is insufficient.

### Fencing

6. Notice acceptance increments a target `maintenanceEpoch` and issues a unique fence token.
7. Placement decisions and deployment plans pin target UID, expected generation, qualification/evidence snapshot and maintenance epoch.
8. DrainPending and later states deny new placement. Draining also prevents new disruptive lifecycle work unless it is an approved evacuation, recovery or maintenance-safe action.
9. PluginExecution validates the current epoch/fence before each externally consequential step. A stale token fails closed with a stable conflict code.
10. Existing operations reach a declared safe point, complete only when policy permits, or are superseded by explicit evacuation/recovery Operations.
11. Requalification refreshes health, capabilities/facts, connectivity, evidence and sovereignty assessments before a new epoch returns to Available.

## Responsibility boundary

The infrastructure operator owns the native maintenance action and its factual notice. Sovrunn owns placement suspension, service-impact assessment, governed customer-service actions, audit and requalification. Neither party silently writes the other's state.

## Conformance

- Unauthorized maintenance announcement and direct target-status mutation are denied.
- Duplicate notices with one idempotency key resolve to one maintenance epoch.
- Stale placement and PluginExecution fences cannot execute after drain acceptance.
- Cancellation does not bypass requalification.
- Provider maintenance completion cannot clear a separately governed owner restriction.
- New placement resumes only after all mandatory qualification and sovereignty gates pass.
