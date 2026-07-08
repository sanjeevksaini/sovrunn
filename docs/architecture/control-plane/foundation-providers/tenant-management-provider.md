# Tenant Management Provider

Document:
  ID: tenant-management-provider
  Title: Tenant Management Provider
  Parent: control-plane-foundation
  Owner: Control Plane Foundation
  Layer: SDE Control Plane
  Type: Foundation Provider Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Tenant Management Provider
  - Define responsibilities and boundaries
  - Support selective AI retrieval

Definition:
  Implements Tenant Management Service using an approved tenant metadata backend.

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
    - ../foundation-services/tenant-management-service.md
  Used By:
    - Foundation Services

References:
  - foundation-providers.md
  - ../foundation-services/tenant-management-service.md
