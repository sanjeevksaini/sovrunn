# ADH-2026-021: CloudEnrollment as the CloudPlatform–Customer Relationship

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-021
- **Date:** 2 August 2026
- **Classification:** New canonical relationship boundary
- **Canonical decision:** FCM-ADR-002

## Context

Making either a CloudPlatform or CloudProvider the parent of a customer Organization would imply ownership of the customer. A customer needs a durable commercial and governance relationship with the customer-facing CloudPlatform without duplicating its organization hierarchy. Provider assignment is a separate supply/placement concern.

## Decision

1. Owner Organization, CloudPlatform, CloudProvider and customer Organization remain independent boundaries.
2. `CloudEnrollment` is a CloudPlatform-scoped ManagedResource joining exactly one CloudPlatform and one customer Organization.
3. The relationship carries customer class, verification state, contract/support references and applied entitlement packages.
4. Actual grants are expressed through `ServiceEntitlement`; limits through `QuotaPolicy`.
5. Enrollment never contains customer tenants, projects, IAM roles or service instances.
6. One Organization may enroll with many CloudPlatforms; one CloudPlatform may enroll many Organizations.
7. Provider eligibility and assignment use `CloudProviderParticipation` and placement policy; they are not encoded as enrollment ownership.

## Alternatives rejected

- **CloudPlatform/CloudProvider → Organization hierarchy:** wrong ownership and weak portability.
- **Copy the Organization per provider:** fragments identity, audit and policy.
- **Enroll directly with CloudProvider:** couples customer contract identity to the selected execution supplier and obstructs provider substitution.
- **Put all grants directly on Organization:** mixes CloudPlatform authority and customer governance.

## Consequences

- CloudPlatform, provider and customer lifecycles remain independent.
- Cross-scope reference authorization and consent/verification are required.
- Terminating enrollment affects consumption but never deletes the Organization.

## Lifecycle and deletion

Enrollment progresses through Pending, Active, Suspended, Terminating and Terminated. CloudPlatformRef and OrganizationRef are immutable. Termination retains a governed summary and invokes an explicit policy for existing services; it never silently deletes them.

## Security and sovereignty

Enrollment lookups use no-existence-disclosure behavior. Contract documents and sensitive verification evidence are referenced, not embedded in customer-readable status.

## Conformance evidence

- Same Organization enrolled with two CloudPlatforms.
- One CloudPlatform enrolling multiple Organizations.
- Inactive enrollment denies new service creation.
- Termination leaves Organization hierarchy intact.
- Unauthorized cross-organization reference returns safe denial.

## Reassessment triggers

A proven marketplace/broker relationship requiring a distinct contracting principal beyond CloudPlatform and Organization, or regulation requiring a separate legally accountable intermediary resource.
