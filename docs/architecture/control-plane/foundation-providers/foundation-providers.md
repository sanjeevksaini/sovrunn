# Foundation Providers

Document:
  ID: foundation-providers
  Title: Foundation Providers
  Parent: control-plane-foundation
  Owner: Control Plane Foundation
  Layer: SDE Control Plane
  Type: Architecture
  Version: 1.0
  Status: Stable

Purpose:
  - Define Foundation Provider architecture
  - Define provider binding model
  - Define provider catalog as index
  - Keep provider category contracts in child files
  - Prevent confusion with Infrastructure Provider, Engine Plugin, and Datastore Operator Plugin

Definition:
  Foundation Providers are pluggable implementations of Foundation Services.

  A Foundation Provider is not the consumer-facing contract. Consumers use Foundation Services. Foundation Services bind to approved Foundation Providers.

Architecture:
  Foundation Provider:
    Role:
      - Implements a Foundation Service category
      - Integrates with external or internal provider technology
      - Declares capabilities, configuration, lifecycle, and failure behavior

  Foundation Service:
    Role:
      - Owns canonical consumer contract
      - Owns provider abstraction boundary
      - Normalizes provider behavior

Provider Catalog:
  Identity and Access:
    - Identity Provider
    - Authorization Provider

  Tenant and Configuration:
    - Tenant Management Provider
    - Configuration Provider

  Policy and Secrets:
    - Policy Provider
    - Secrets Provider

  Governance and Telemetry:
    - Audit Provider
    - Observability Provider

  Coordination:
    - Workflow Provider
    - Eventing Provider

  Frameworks:
    - Registry Provider
    - Plugin Provider

Binding Model:
  - Provider is registered through Plugin Framework.
  - Provider metadata is validated.
  - Provider capability is checked.
  - Provider configuration is approved.
  - Foundation Service binds to approved provider.
  - Consumers continue using Foundation Service contract.

Extension Model:
  Provider replacement MUST:
    - Preserve Foundation Service contract
    - Preserve security behavior
    - Preserve tenant isolation
    - Preserve audit requirements
    - Preserve error normalization

Security Model:
  Foundation Providers MUST:
    - Protect provider credentials
    - Avoid exposing unsafe provider internals
    - Preserve tenant and authorization context
    - Support audit requirements where applicable

Failure Model:
  Provider failure MUST:
    - Be translated through Foundation Service
    - Preserve safe diagnostic information
    - Avoid leaking unsafe provider-native errors
    - Preserve retry classification where applicable

Invariants:
  - Foundation Provider is not Infrastructure Provider.
  - Foundation Provider is not Engine Plugin.
  - Foundation Provider is not Datastore Operator Plugin.
  - Consumers do not depend directly on providers by default.
  - Provider-specific behavior must not become canonical SDE behavior.

Relationships:
  Parent:
    - control-plane-foundation.md
  Children:
    - identity-provider.md
    - authorization-provider.md
    - tenant-management-provider.md
    - configuration-provider.md
    - policy-provider.md
    - secrets-provider.md
    - audit-provider.md
    - workflow-provider.md
    - eventing-provider.md
    - observability-provider.md
    - registry-provider.md
    - plugin-provider.md
  Depends On:
    - foundation-services/foundation-services.md
  Used By:
    - Foundation Services

References:
  - ../control-plane-foundation.md
  - ../foundation-services/foundation-services.md
  - identity-provider.md
  - authorization-provider.md
  - tenant-management-provider.md
  - configuration-provider.md
  - policy-provider.md
  - secrets-provider.md
  - audit-provider.md
  - workflow-provider.md
  - eventing-provider.md
  - observability-provider.md
  - registry-provider.md
  - plugin-provider.md
