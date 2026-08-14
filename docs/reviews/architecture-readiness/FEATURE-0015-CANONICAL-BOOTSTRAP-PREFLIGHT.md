# FEATURE-0015 Canonical Bootstrap Preflight

## Status

- Date: 2026-08-11
- Proposed controlling change: ADH-2026-045, canonical bootstrap with no alpha runtime migration
- Review purpose: identify every known architecture decision that requirements, design, or tasks would otherwise be forced to invent.
- Result: **DECISIONS CLOSED — controlling-authority application and admission validation required before Kiro requirements generation.**

This is a review artifact. It does not modify the active architecture, which continues to contain migration semantics until ADH-2026-045 is validated and applied atomically by Kiro.

## Intended FEATURE-0015 boundary after ADH-2026-045

FEATURE-0015 establishes canonical-only cloud identity resources through InfrastructureStack: CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, and InfrastructureStack. It reuses FEATURE-0012 resource/error conventions, FEATURE-0013 audit evidence, and compatible completed implementation infrastructure after FEATURE-0011 reuse assessment. ExecutionTarget moves in its entirety to FEATURE-0016. FEATURE-0015 does not implement runtime alpha migration, provider-native execution, target registration or qualification, governance/IAM, enrollment, catalog, placement, plugins, or provisioning.

### Approved single-authority provider onboarding model

The CloudPlatform control plane is the sole authoritative store for the CloudProvider, its topology, and its CloudProviderParticipation. A provider's own Sovrunn installation in a service location is not a federated control plane for FEATURE-0015 and does not create a duplicate canonical resource set. Its provider administrator uses the CloudPlatform-authoritative API to register the provider and topology, then requests participation. No remote deployment resource, cross-control-plane replication, remote-offer exchange, or federation trust protocol is in scope.

### Approved independent participation suspension holds

`CloudProviderParticipation.status.platformSuspended` and `status.providerSuspended` are durable controller-owned lifecycle evidence, not PATCHable spec fields. The CloudPlatform administrator controls only the former and the CloudProvider administrator controls only the latter through their respective explicit suspend/resume actions. Accepted participation is effectively `Active` only when both holds are false; it is effectively `Suspended` when either is true. A party's resume action clears only its own hold and cannot override the other party's suspension. Every hold-changing action must produce an AuditEvent with the actor, controlled hold, participation UID, result, request ID, and resulting effective state.

## Decision-closure record

The architecture owner approved the complete decision set in ADH-2026-045. The rows below remain as the audit trail of gaps that the active authority update and admission checker must close; they are not permission for requirements, design, or tasks to reinterpret the approved decisions.

| ID | Closure status |
|---|---|
| F15-CB-01 | Closed — ExecutionTarget entirely FEATURE-0016 |
| F15-CB-02 | Closed — participation action/state and evidence model in ADH-2026-045 |
| F15-CB-03 | Closed — provider-chain invariants and safe-denial ordering in ADH-2026-045 |
| F15-CB-04 | Closed — server-initialized Active status; no F0015 lifecycle/delete for identity/topology resources |
| F15-CB-05 | Closed — deterministic server-resolved scoped bootstrap grants |
| F15-CB-06 | Closed — canonical root/order, v1alpha1 API boundary, and no alpha routes |
| F15-CB-07 | Closed — AuditEvent action matrix and scheduler correlation |
| F15-CB-08 | Closed — merge-patch, resourceVersion, idempotency, and pair-race semantics |

## Architecture decision audit trail

| ID | Decision that must be made before requirements | Evidence of current gap | Recommended architecture outcome | Required proof before admission |
|---|---|---|---|---|
| F15-CB-01 | **ExecutionTarget feature ownership.** Does FEATURE-0015 create an incomplete target, or does FEATURE-0016 own the full contract? | VS0-SCHEMA-015 requires target status while F0015 prohibits persisting it. | **Resolved:** ExecutionTarget moves entirely to FEATURE-0016; F0015 ends at InfrastructureStack. | Registry, routes, writers, state, conformance, feature scopes, and the readiness checker assign all ExecutionTarget behavior only to FEATURE-0016. |
| F15-CB-02 | **CloudProviderParticipation acceptance protocol.** What exactly constitutes the two delegated acceptances needed for `Pending -> Active`, where is each recorded, and who writes it before FEATURE-0018 IAM/approval resources exist? | VS0-STATE-001 requires both delegated authorities; VS0-SCHEMA-010 has no acceptance field/ref; VS0-WRITER-004 names a contract authority but does not define the actor/action/proof. | Choose one bounded Phase 2R bootstrap protocol, such as two uid-pinned acceptance attestations held in controller-owned status and written only through two canonical actions by scoped principals. It must not introduce FEATURE-0018's generic ApprovalRequest or role model early. | Exact fields or action payloads, principal/scope rules, transition writer, duplicate/revoke behavior, errors, audit correlation, and positive/negative conformance are declared. |
| F15-CB-03 | **Topology and participation invariants.** Which references must share the same CloudProvider authority, and which failures are safe 404 denial versus authorized 422 invalidity? | Only CloudPlatform/participation scope-reference equality is currently explicit. Typed refs alone allow a Datacenter to reference another provider's location or an InfrastructureStack to reference another provider's FaultDomain. | Require UID equality across the topology chain: Datacenter↔HostingLocation, FaultDomain↔Datacenter, and InfrastructureStack↔FaultDomain. Inaccessible cross-provider targets return safe `RESOURCE_NOT_FOUND`; authorized same-provider structural failures return `VALIDATION_FAILED`. | Registry validations, exact violations, reference-resolution order, safe-projection rule, and local positive/negative conformance cases exist. |
| F15-CB-04 | **Lifecycle, status, and deletion model for F0015-owned managed resources.** Who initializes and changes CloudPlatform/CloudProvider/topology `status.phase`, and what is allowed to be deleted or decommissioned in this feature? | Schemas require phase status and use controller-only mutability, but no F0015 controller/lifecycle state machines define transitions. The catalog requires drain/decommission behavior before removal; F0015 currently claims CRUD. | Define a minimal bootstrap lifecycle: server creates the canonical identity resource in a declared initial phase; F0015 either supports only explicit approved phase transitions/deletion preconditions or intentionally exposes no deletion transition. Do not leave a controller-owned required field without an owner. | Initial state, transition authority, delete/decommission policy, dependency checks, terminal behavior, error mapping, audit events, and conformance cases are all declared. |
| F15-CB-05 | **Bootstrap authorization boundary.** How can F0015 enforce CloudPlatform, CloudProvider, and participation writers without importing FEATURE-0018 IAM/role/approval capability? | VS0-WRITER-002..004 name actors, while FEATURE-0018 owns Membership, RoleDefinition, and RoleAssignment. The canonical model forbids persona-name authorization. | Define a narrow bootstrap authorization contract based on authenticated typed principals plus canonical actions and scoped resource ownership. It is not a RoleDefinition/RoleAssignment implementation and must be replaced/consumed compatibly by FEATURE-0018. | Allowed actions, principal attributes, scope matching, denial/no-disclosure behavior, audit fields, FEATURE-0018 handoff, and test fixtures are explicit. |
| F15-CB-06 | **Canonical-only API and bootstrap input contract.** What exact API groups, versions, collection routes, create/update/delete behavior, and seed prerequisites are active in the first control plane? | The registry names API identities, while the older Phase 1 API contract exposes different `/v1` routes and alpha resources. The catalog requires exact routes and schemas before implementation. | Publish the F0015 canonical API surface and state that alpha routes are not registered. The CloudPlatform is the root resource and contains its immutable owner registration; a CloudProvider is registered in that same authoritative control plane before topology or participation creation. | Route/verb table, canonical schemas, request/response projections, reference-order rules, no-active-alpha-route proof, and API conformance cases exist. |
| F15-CB-07 | **F0015 audit and observability contract after migration removal.** Which remaining actions are material enough to produce AuditEvents, and which correlation fields are mandatory? | Existing audit section is dominated by removed plan/record/milestone events. FEATURE-0013 is reused but the canonical-bootstrap event set is not stated. | Require correlated, redacted FEATURE-0013 AuditEvents for participation creation/activation/suspension/termination, protected canonical-resource lifecycle changes, and authorization/safe-denial outcomes. Define when ordinary resource registration is logged only versus audited. | Event/action matrix, subject and actor references, request correlation, projection/redaction, no-secret rule, and observability/no-side-effect conformance exist. |
| F15-CB-08 | **Concurrency, idempotency, and update semantics.** What generic FEATURE-0012 behavior is adopted, and what F0015-specific uniqueness/update rules apply? | Pair uniqueness is stated for participation, but replay keys, optimistic concurrency, idempotent create/update behavior, and patch format are not bounded by the F0015 architecture. | Explicitly adopt FEATURE-0012 optimistic-concurrency and replay rules, add F0015 resource-name and participation-pair uniqueness, and expose PATCH only as `application/merge-patch+json`; PUT and DELETE are absent. | Conflict codes/violations, replay behavior, stale-version outcome, field-level patch matrix, and race/negative conformance are declared. |

## Approved participation-visibility boundary

The architecture owner has decided:

> CloudPlatform decides which enrolled customer Organizations may see or select an Active CloudProvider participation. Individual user visibility is derived from that Organization-level eligibility and scoped authorization; it is not stored directly on CloudProviderParticipation.

Consequences:

- An Active participation is a necessary but insufficient condition for customer visibility.
- FEATURE-0015 exposes no end-user CloudProvider projection; non-active participation states are never customer-visible.
- FEATURE-0021 owns Organization/enrollment eligibility and the CloudPlatform visibility decision.
- FEATURE-0018 owns individual-user authorization.
- FEATURE-0023 may select only participations that are Active and eligible under the resulting Organization-level visibility rule.

## Decisions already sufficiently bounded by the canonical architecture

The following are not new decision gaps once ADH-2026-045 is applied; they must be carried into the feature authority without reinterpretation:

- Seven canonical scope kinds; no generic Provider, ResourcePool, ProviderCapability, ServiceClass, EffectivePolicyContext, or six-scope active behavior.
- CloudPlatform ownership is distinct from CloudProvider operational supply.
- CloudProviderParticipation is the only relationship boundary between a platform and provider.
- Geography/topology identity remains provider-scoped and does not prove sovereignty or qualification.
- ExecutionTarget is a realization boundary, not a mandatory ResourcePool or a customer-visible topology resource.
- Phase 2R is deterministic, in-memory, side-effect-free, with no real provider, Kubernetes, OpenShift, PostgreSQL, credential, or native-object execution.
- FEATURE-0012 owns common metadata, typed references, error envelope, system-owned fields, and generic validation grammar.
- FEATURE-0013 owns AuditEvent; FEATURE-0015 may only consume it.

## Guardrails required in the architecture admission checker

The generic Phase 2R checker proposed by the completeness audit must reject FEATURE-0015 unless it can prove:

1. no active migration resource, writer, state machine, conformance case, route, or requirement remains;
2. no field is simultaneously required in FEATURE-0015 and owned for later introduction/activation;
3. every state transition names its stored evidence, writer, authorization action, error, audit event, and conformance case;
4. every cross-resource reference has an explicit same-authority invariant or an explicit permitted cross-boundary contract;
5. every protected writer is enforceable by an authorization contract available in this feature, without importing a future feature;
6. every managed-resource status field has an initial value, writer, transition/deletion rule, and projection rule;
7. every create/update/delete route has a canonical API version, scope, concurrency behavior, and stable failure mapping;
8. feature-local conformance covers each owned invariant and does not rely on a future-feature scenario; and
9. requirements, design, and tasks context contains one approved feature authority and does not copy stale migration ledgers.

## Admission decision

FEATURE-0015 is **not requirements-ready** until the controlling authorities are updated atomically and the architecture admission gate proves the approved decisions. At that point, the requirements stage may render the approved decisions; it must not choose among them.
