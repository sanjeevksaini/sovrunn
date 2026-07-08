# Plugin Provider

Document:
  ID: plugin-provider
  Title: Plugin Provider
  Parent: control-plane-foundation
  Owner: Control Plane Foundation
  Layer: SDE Control Plane
  Type: Foundation Provider Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Plugin Provider
  - Define responsibilities and boundaries
  - Support selective AI retrieval

Definition:
  Implements Plugin Framework storage or artifact discovery using OCI artifacts, filesystem, Kubernetes CRDs, or Git.

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
    - ../foundation-services/plugin-framework.md
  Used By:
    - Foundation Services

References:
  - foundation-providers.md
  - ../foundation-services/plugin-framework.md
