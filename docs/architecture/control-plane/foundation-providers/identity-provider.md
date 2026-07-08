# Identity Provider

Document:
  ID: identity-provider
  Title: Identity Provider
  Parent: control-plane-foundation
  Owner: Control Plane Foundation
  Layer: SDE Control Plane
  Type: Foundation Provider Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Identity Provider
  - Define responsibilities and boundaries
  - Support selective AI retrieval

Definition:
  Implements Identity Service using identity systems such as Keycloak, OIDC, LDAP, SAML, or cloud IAM.

Responsibilities:
  MUST:
    - Implement corresponding Foundation Service contract
    - Declare capability, configuration, lifecycle, and failure behavior
    - Preserve tenant, policy, security, and audit boundaries

  MUST NOT:
    - Expose provider-specific API as canonical consumer contract
    - Bypass Foundation Service
    - Own domain Management Plane semantics

Inputs:
  - Provider configuration
  - Provider credentials reference
  - Provider operation request

Outputs:
  - Provider operation result
  - Provider error

State:
  - Provider lifecycle metadata
  - Provider capability metadata

Failure Rules:
  - Normalize failure through Foundation Service
  - Avoid leaking unsafe provider internals

Relationships:
  Parent:
    - control-plane-foundation
  Depends On:
    - ../foundation-services/identity-service.md
  Used By:
    - Foundation Services

References:
  - foundation-providers.md
  - ../foundation-services/identity-service.md
