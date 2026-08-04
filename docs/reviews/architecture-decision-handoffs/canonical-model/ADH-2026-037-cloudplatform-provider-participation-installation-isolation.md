# ADH-2026-037: CloudPlatform Ownership, CloudProvider Participation and Installation Isolation

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-037
- **Date:** 4 August 2026
- **Classification:** New customer-cloud ownership and provider-supply boundary
- **Canonical decision:** FCM-ADR-018

## Context

A customer-facing cloud may be owned and branded by one organization while contracted CloudProviders operate different datacenters. NIC may own NIC Cloud, contract Yotta for Delhi/Pune and Jio for Bhubaneswar/Hyderabad, and present one cloud to government departments. Treating CloudProvider as both cloud owner and infrastructure operator cannot represent this arrangement cleanly.

## Decision

1. `CloudPlatform` is the customer-facing cloud product and governance boundary. It has exactly one owning `Organization`.
2. `CloudProvider` is an independent operational supply boundary and is never modeled as owned merely because it is contracted by the CloudPlatform owner.
3. `CloudProviderParticipation` is the governed relationship between exactly one CloudPlatform and exactly one CloudProvider. It records effective period, permitted locations, portfolios, responsibilities and protected agreement references.
4. `CloudEnrollment` joins a customer Organization to a CloudPlatform. It supersedes customer enrollment directly with a CloudProvider.
5. A `SovrunnInstallation` serves exactly one CloudProviderParticipation and therefore exactly one CloudPlatform and one CloudProvider.
6. Initially, at most one active production installation exists per CloudPlatform + CloudProvider + environment. HA replicas belong to one logical installation.
7. The CloudPlatform publishes the customer-facing portfolio, offering and plan. A participating CloudProvider supplies one or more governed realization mappings for eligible offerings.
8. Provider selection is `CustomerSelected`, `OwnerAssigned`, `PolicySelected` or `Fixed`. The authoritative PlacementDecision records the selected participation/provider; the customer projection discloses provider identity only when policy permits.
9. A CloudProvider may participate in another CloudPlatform only through a separate participation and isolated installation.

## Invariants

- One installation never executes for multiple CloudProviders.
- Provider credentials and protected handles remain isolated to that provider installation.
- CloudPlatform ownership does not grant unrestricted access to provider credentials.
- Cross-provider migration or failover requires a new placement, sovereignty and data-movement evaluation.
- Customer Organizations are related through CloudEnrollment and are never children of the owner Organization.

## Example

```text
Organization/NIC
  └── owns CloudPlatform/NIC-Cloud
        ├── participation → CloudProvider/Yotta
        │     └── SovrunnInstallation/nic-yotta-production
        └── participation → CloudProvider/Jio
              └── SovrunnInstallation/nic-jio-production

Organization/Department-of-Posts
  └── CloudEnrollment → CloudPlatform/NIC-Cloud
```

## Consequences

CloudProviderEnrollment becomes CloudEnrollment; customer-facing catalog ownership moves to CloudPlatform; provider realization eligibility moves to CloudProviderParticipation; placement and lifecycle records retain the selected participation. Existing Provider data requires classification rather than a blind rename.
