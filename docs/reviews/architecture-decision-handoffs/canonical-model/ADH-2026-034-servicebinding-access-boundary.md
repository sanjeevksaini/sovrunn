# ADH-2026-034: ServiceBinding as the Managed-Service Access Boundary

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-034
- **Date:** 2 August 2026
- **Classification:** Phase 1 contract retained and strengthened
- **Canonical decision:** FCM-ADR-015
- **Affected baseline:** FEATURE-0008, planned FEATURE-0033, PostgreSQL MVP

## Context

A ready PostgreSQL ServiceInstance is not usable until an authorized consumer receives a governed endpoint/identity/credential relationship. The final canonical model omitted ServiceBinding even though Phase 1 already defines it and the PostgreSQL MVP requires SecretRef integration. Returning credentials in ServiceInstance status would violate least privilege and prevent per-consumer revocation.

## Decision

1. Retain `ServiceBinding` as a Project-scoped customer-facing ManagedResource.
2. It references exactly one ServiceInstance and one typed consumer identity/workload reference.
3. The referenced ServiceTypeDefinition supplies the versioned binding schema and supported access modes for that service family; the initial PostgreSQL form uses endpoint and secret/identity references.
4. Secret values are never stored in ServiceBinding; status exposes only authorized typed references and safe connection metadata.
5. Binding creation is an Operation and may require approval, rotation, expiry and network/policy validation.
6. A ServiceInstance cannot be deleted while active bindings exist unless an explicit authorized cascade/revocation flow is approved.

## Alternatives rejected

- **Credentials in ServiceInstance status:** broad exposure and shared blast radius.
- **Backend-native secret only:** loses Sovrunn governance, portability and lifecycle.
- **No first-class binding:** cannot model per-consumer access, rotation or revocation.

## Consequences

- The final canonical model and PostgreSQL reference flow add ServiceBinding after readiness.
- FEATURE-0033 evolves rather than replaces the Phase 1 concept.
- Workload identity may replace static credentials where supported without changing the binding relationship.

## Lifecycle and deletion

Requested → Provisioning → Ready → Rotating/Expiring → Revoked/Failed. ConsumerRef and ServiceInstanceRef are immutable. Rotation creates a governed operation and new secret/identity material behind references. Revocation is idempotent.

## Security and sovereignty

Authorization is evaluated for the requesting principal and consumer. Secret reads remain in the secret system’s policy boundary. Endpoint and network metadata are purpose-filtered. Every issue/rotate/revoke action is audited.

## Conformance evidence

- Authorized application receives a binding with SecretRef, not a value.
- Unauthorized consumer denied without secret/service disclosure.
- Independent bindings rotate/revoke separately.
- PostgreSQL and a second service family use different binding schemas without changing the core ServiceBinding envelope.
- Service deletion blocked by active binding.
- Retry produces no duplicate credential grant.
- Customer projection contains no backend credential or privileged endpoint detail.

## Reassessment triggers

All supported services use a standardized workload-identity protocol that fully replaces secret/endpoint bindings while preserving per-consumer authorization, rotation and revocation.
