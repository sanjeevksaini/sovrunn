# Management Plane

Document:
  ID: management-plane
  Title: Management Plane
  Parent: sde-control-plane
  Owner: SDE Control Plane
  Layer: SDE Control Plane
  Type: Architecture
  Version: 1.1
  Status: Stable

Purpose:
  - Define pluggable Management Plane architecture
  - Define Management Plane Framework responsibility
  - Define Management Plane Manifest and admission boundary
  - Define relationship with Core Control Plane
  - Define relationship with Datastore Management Plane
  - Clarify that Datastore Management Plane is the first pluggable management plane

Definition:
  A Management Plane is a domain-specific, pluggable control-plane domain hosted under SDE Control Plane authority.

  A Management Plane owns domain-specific management semantics while relying on Control Plane Foundation for shared services such as identity, authorization, policy, configuration, secrets, audit, workflow, eventing, observability, registry, and plugin lifecycle.

  Datastore Management Plane is the first Management Plane.

Core Principle:
  Management Planes are governed extensions of SDE Control Plane.

  A Management Plane must not bypass SDE Control Plane authority.

  A Management Plane must not execute SDE Data Plane tenant requests.

Architecture:
  SDE Control Plane:
    Contains:
      - Management Plane Framework

  Management Plane Framework:
    Hosts:
      - Datastore Management Plane
      - Future Management Planes

  Pluggable Management Plane:
    Provides:
      - Installable management domain
      - Domain metadata
      - Lifecycle hooks
      - Controller runtime integration
      - Required Foundation Service dependencies
      - Domain APIs
      - Domain workflows

  Datastore Management Plane:
    Role:
      - First pluggable management plane.
      - Manages tenant-scoped Downstream Datastore lifecycle and operations.
      - Powers dstoreOps.

Component Model:
  Management Plane Framework:
    Responsibilities:
      - Register
      - Validate
      - Authorize
      - Admit
      - Enable
      - Upgrade
      - Disable
      - Remove
      - Expose management-plane APIs
      - Integrate with Workflow Service
      - Integrate with Policy Service
      - Integrate with Audit Service
      - Integrate with Observability Service

  Management Plane Manifest:
    Responsibilities:
      - Declare management plane identity
      - Declare version
      - Declare domain
      - Declare required Foundation Services
      - Declare API surface
      - Declare controller runtime requirements
      - Declare supported resources
      - Declare compatibility
      - Declare security requirements

  Management Plane Controller Runtime:
    Responsibilities:
      - Host management-plane controllers
      - Reconcile management-plane resources
      - Execute approved workflows
      - Report health
      - Emit events, metrics, traces, and audit records

  Pluggable Management Plane:
    Responsibilities:
      - Own domain registries
      - Own domain lifecycle
      - Own domain workflows
      - Expose authorized management APIs
      - Emit domain events
      - Record domain audit

Control Flow:
  Admission flow:
    - Management Plane Manifest is submitted.
    - Management Plane Framework validates manifest.
    - Required Foundation Services are resolved.
    - Compatibility is checked.
    - Authorization Service validates enablement authority.
    - Policy Service evaluates admission policy.
    - Audit Service records admission decision.
    - Management Plane is admitted or rejected.

  Runtime flow:
    - Management Plane Controller Runtime starts.
    - Approved Management Plane controllers are loaded.
    - Management Plane resources are reconciled.
    - Domain workflows are executed through Workflow Service.
    - Policy and audit are applied to state-changing actions.

DMP Relationship:
  Datastore Management Plane:
    Is:
      - A pluggable management plane.

    Is Not:
      - The same thing as DMP Controller Runtime.
      - A SDE Data Plane component.
      - An Engine Plugin.
      - A Datastore Operator Plugin.

  DMP Controller Runtime:
    Is:
      - The executable runtime that hosts and reconciles DMP.

    Is Not:
      - The whole DMP.

Future Management Planes:
  Examples:
    - Cache Management Plane
    - Search Management Plane
    - Vector Management Plane
    - Data Pipeline Management Plane
    - Tenant Integration Management Plane

  Rule:
    - Future Management Planes must use Management Plane Framework admission, lifecycle, policy, workflow, audit, and observability boundaries.

Boundaries:
  Management Plane May:
    - Own domain-specific lifecycle
    - Own domain-specific resources
    - Own domain-specific workflows
    - Use Foundation Services
    - Use approved plugins/providers for its domain

  Management Plane Must Not:
    - Execute tenant data-plane requests
    - Bypass Authorization Service
    - Bypass Policy Service
    - Bypass Workflow Service
    - Bypass Audit Service
    - Bypass Secrets Service controls
    - Mutate SDE Data Plane runtime state directly
    - Invoke Engine Plugins directly unless explicitly defined by a future RFC

Invariants:
  - Management Plane is a pluggable SDE Control Plane domain.
  - Management Plane Framework owns admission and lifecycle of management planes.
  - Datastore Management Plane is the first pluggable management plane.
  - DMP Controller Runtime is not the whole DMP.
  - Management Planes must use Foundation Services.
  - Management Planes must be policy-governed and audited.
  - Management Planes must not execute tenant data-plane requests.

Related Documents:
  - control-plane.md
  - control-plane-map.md
  - management-plane-framework/management-plane-framework.md
  - datastore-management-plane/datastore-management-plane.md
  - control-plane-foundation.md
