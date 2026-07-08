# Tenant Management Service

Document:
  ID: tenant-management-service
  Title: Tenant Management Service
  Parent: control-plane-foundation
  Owner: Control Plane Foundation
  Layer: SDE Control Plane
  Type: Foundation Service Contract
  Version: 1.0
  Status: Stable

Purpose:
  - Define Tenant Management Service
  - Define responsibilities and boundaries
  - Support selective AI retrieval

Definition:
  Tenant Management Service owns tenant identity, tenant metadata, and tenant-level governance context.

Responsibilities:
  MUST:
    - Create tenant metadata
    - Update tenant metadata
    - Resolve tenant context
    - Preserve tenant isolation

  MUST NOT:
    - Store tenant application data
    - Execute tenant data requests
    - Bypass authorization

Inputs:
  - Tenant request
  - Actor identity
  - Policy context

Outputs:
  - Tenant context
  - Tenant metadata

State:
  - Tenant registry metadata

Failure Rules:
  - Return deterministic error
  - Preserve traceability
  - Avoid leaking unsafe internal details

Relationships:
  Parent:
    - control-plane-foundation
  Depends On:
    - foundation-providers/tenant-management-provider.md
  Used By:
    - Management Planes
    - SDE Data Plane where authorized

References:
  - foundation-services.md
  - foundation-providers/tenant-management-provider.md
