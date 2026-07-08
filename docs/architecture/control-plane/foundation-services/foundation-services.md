# Foundation Services

Document:
  ID: foundation-services
  Title: Foundation Services
  Parent: control-plane-foundation
  Owner: Control Plane Foundation
  Layer: SDE Control Plane
  Type: Architecture
  Version: 1.0
  Status: Stable

Purpose:
  - Define Foundation Service architecture
  - Define common service rules
  - Define service catalog as index
  - Keep service contracts in child files
  - Prevent provider-specific coupling

Definition:
  Foundation Services are stable SDE Control Plane contracts exposed by Control Plane Foundation.

  They provide shared management capabilities used by Management Planes and authorized SDE components.

Architecture:
  Foundation Services:
    Categories:
      - Identity and access
      - Tenant and configuration
      - Policy and secrets
      - Audit and observability
      - Workflow and eventing
      - Registry and plugin framework

Service Catalog:
  Identity and Access:
    - Identity Service
    - Authorization Service

  Tenant and Configuration:
    - Tenant Management Service
    - Configuration Service

  Policy and Secrets:
    - Policy Service
    - Secrets Service

  Governance and Telemetry:
    - Audit Service
    - Observability Service

  Coordination:
    - Workflow Service
    - Eventing Service

  Frameworks:
    - Registry Framework
    - Plugin Framework

Service Interaction Model:
  Consumers:
    - Management Planes
    - Core Control Plane
    - Datastore Management Plane
    - Authorized SDE runtime components

  Provider Boundary:
    - Consumers call Foundation Services.
    - Foundation Services call Foundation Providers.
    - Consumers do not directly depend on providers unless explicitly authorized.

Control Flow:
  - Consumer calls Foundation Service.
  - Service validates request shape and authorization context.
  - Service applies service-level policy.
  - Service delegates to bound provider where required.
  - Provider-specific response is normalized.
  - Service returns stable result or service-level error.

State Model:
  Foundation Service state includes:
    - Service contract version
    - Binding metadata
    - Consumer permissions
    - Service-level configuration
    - Service-level policy

Security Model:
  Foundation Services MUST:
    - Preserve tenant isolation
    - Preserve service authorization
    - Avoid leaking provider credentials
    - Avoid exposing unsafe provider details
    - Use Secrets Service for secrets

Failure Model:
  Foundation Service failures MUST:
    - Return canonical service errors
    - Preserve traceability
    - Preserve auditability when security-sensitive
    - Avoid exposing unsafe provider internals

Invariants:
  - Foundation Services are consumer-facing contracts.
  - Foundation Providers are implementation details.
  - Service contracts are stable across provider replacement.
  - Configuration Service does not store secrets.
  - Secrets Service owns secret access.
  - Audit Service is not replaced by Observability Service.
  - Eventing Service is not authoritative registry state.
  - Workflow Service coordinates long-running management operations but does not own domain semantics.

Relationships:
  Parent:
    - control-plane-foundation.md
  Children:
    - identity-service.md
    - authorization-service.md
    - tenant-management-service.md
    - configuration-service.md
    - policy-service.md
    - secrets-service.md
    - audit-service.md
    - workflow-service.md
    - eventing-service.md
    - observability-service.md
    - registry-framework.md
    - plugin-framework.md
  Depends On:
    - foundation-providers/foundation-providers.md
  Used By:
    - management-plane.md
    - core-control-plane/core-control-plane.md
    - datastore-management-plane/datastore-management-plane.md

References:
  - ../control-plane-foundation.md
  - ../foundation-providers/foundation-providers.md
  - identity-service.md
  - authorization-service.md
  - tenant-management-service.md
  - configuration-service.md
  - policy-service.md
  - secrets-service.md
  - audit-service.md
  - workflow-service.md
  - eventing-service.md
  - observability-service.md
  - registry-framework.md
  - plugin-framework.md
