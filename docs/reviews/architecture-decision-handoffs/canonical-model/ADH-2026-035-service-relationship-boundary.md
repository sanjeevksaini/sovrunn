# ADH-2026-035: Implementation-Neutral Service Relationship Boundary

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-035
- **Date:** 4 August 2026
- **Classification:** New canonical managed-resource and extension-contract boundary
- **Canonical decision:** FCM-ADR-016
- **Affected baseline:** FEATURE-0006–0008 service contracts; planned placement, plugin, adapter and connectivity work

## Context

Sovrunn must realize relationships between independently managed services without importing execution-infrastructure objects into its core. For example, a virtual-machine service may privately consume a PostgreSQL service while each is assigned to a provider-defined logical network placement. AWS may realize the intent with VPCs, subnets, routes and security groups; OCI, OpenStack and OpenShift use different native constructs. Making those objects core resources would break implementation neutrality. Modeling only a ServiceBinding would also be insufficient because credentials do not express dependency, lifecycle, latency, placement, connectivity or sovereignty obligations.

## Decision

1. `ServiceRelationshipDefinition` is an immutable VersionedDefinition that declares an allowed directed relationship between service families.
2. It owns the relationship parameter schema, compatible source/target ServiceTypeDefinition selectors, requirement resolver, applicable DecisionRecord profiles, supported binding forms and default lifecycle semantics.
3. `ServiceRelationship` is a ManagedResource scoped to the source ServiceInstance's Project, connecting exactly one source ServiceInstance to exactly one target ServiceInstance through an exact ServiceRelationshipDefinition version.
4. Relationship parameters are bounded and schema-validated; arbitrary property bags are prohibited.
5. A required relationship contributes to ServiceRequirementSet and EffectiveGovernanceContext and may trigger sovereignty, placement or relationship-eligibility DecisionRecords.
6. `ServiceBinding` separately realizes consumer identity, endpoint and SecretRef access. A binding neither proves nor replaces relationship eligibility or connectivity.
7. `ServiceDeploymentPlan` contains implementation-neutral relationship actions. Plugins interpret service semantics; ExecutionTarget adapters translate them into protected implementation-native operations.
8. Relationship change, failure or deletion follows explicit lifecycle policy and may trigger reassessment, re-placement, re-binding, degradation, compensation or deletion restriction.
9. Cross-Project relationships require authorization at both ends. Cross-Organization or cross-CloudProvider relationships require separately approved sharing, supply or federation contracts and are denied until those contract types exist and an applicable contract is active.
10. Implementation-native networking types, identifiers and fields are prohibited from customer and canonical core contracts.

## Core anti-drift rule

Sovrunn core models service outcomes, relationships, requirements, decisions and lifecycle. It must not contain API types, fields, conditionals or workflow branches for AWS VPCs, OCI VCNs, OpenStack networks, Kubernetes NetworkPolicies or equivalent implementation-native constructs.

Implementation-native networking objects and identifiers are restricted to provider configuration, service plugins, execution adapters and protected execution records. They must not become customer or canonical core contracts.

## Alternatives rejected

- **Put VPC/subnet/network-policy resources in core:** couples Sovrunn to current infrastructure implementations.
- **Use ServiceBinding as the relationship:** conflates credentials with dependency, placement and lifecycle.
- **Keep relationships only in ServiceCompositionDefinition:** cannot represent independently created or subsequently connected ServiceInstances.
- **Untyped relationship labels and parameters:** make governance, portability and compatibility unverifiable.
- **Infer connectivity from common network or topology ancestry:** membership does not prove routing, policy, reachability or latency.

## Consequences

- The canonical customer model gains ServiceRelationship; the provider/ecosystem extension model gains ServiceRelationshipDefinition.
- ServiceTypeDefinition advertises compatible relationship contracts without embedding implementation-native fields.
- Relationship evaluators adopt FEATURE-0013 DecisionRecord profiles rather than creating a competing decision envelope.
- The same relationship contract can be realized by different plugins/adapters on AWS, OCI, OpenStack, OpenShift or other targets.
- Provider-packaged profiles keep ordinary customer choices simple while allowing advanced required/preferred relationship intent.

## Lifecycle and deletion

Requested → Evaluating → Planning → Realizing → Ready → Degraded/Reassessing → Removing → Removed/Failed. SourceRef, targetRef and relationshipDefinitionRef are immutable. A required relationship may restrict target deletion. Removing a relationship revokes or updates associated bindings and implementation state only through explicit idempotent operations.

## Security and sovereignty

Both service references are authorized before existence is disclosed. Relationship parameters and status contain no credentials or protected native handles. Data-moving or administrative relationships participate in sovereignty evaluation. Plugins and adapters use operation- and target-scoped authority; all relationship decisions and lifecycle actions are audited.

## Conformance evidence

- VM-to-PostgreSQL private-consumption intent succeeds on two different execution implementations without changing core/customer fields.
- The relationship references provider-packaged logical placement profiles, not native VPC/subnet identifiers.
- Same logical network ancestry without explicit validated relationship eligibility remains Unknown or denied.
- ServiceBinding success does not mark ServiceRelationship Ready until required evaluation and realization gates pass.
- Required relationship blocks unsafe target deletion; optional relationship follows its declared failure policy.
- Cross-scope reference denial reveals neither service existence nor native topology.
- No core schema contains VPC, VCN, subnet, route, security-group, NetworkPolicy or equivalent implementation-specific fields.

## Reassessment triggers

An adopted cross-platform service-relationship standard provides equivalent identity, schema, lifecycle, decision and security semantics; or three service families demonstrate that a common relationship semantic should be promoted from an extension definition into a typed core field without introducing implementation coupling.
