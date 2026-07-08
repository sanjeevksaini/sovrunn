# SDE Control Plane

Document:
  ID: sde-control-plane
  Title: SDE Control Plane
  Parent: architecture
  Owner: SDE Control Plane
  Layer: SDE Control Plane
  Type: Architecture
  Version: 1.1
  Status: Stable

Purpose:
  - Define SDE Control Plane authority
  - Define SDE Control Plane architecture
  - Define management and governance boundaries
  - Define relationship with SDE Data Plane
  - Define relationship with the Management Plane Framework
  - Define Datastore Management Plane as the first pluggable management plane
  - Define relationship with Control Plane Foundation
  - Reserve optional AI Control Plane extension boundary

Definition:
  SDE Control Plane is the management authority, governance plane, and pluggable management-plane host for Sovrunn Data Engine.

  SDE Control Plane governs SDE configuration, policy, tenant context, runtime metadata, plugin metadata, engine metadata, capability approval, management-plane lifecycle, and approved state consumed by SDE Data Plane.

Core Principle:
  SDE Control Plane hosts governed control domains.

  Datastore Management Plane is not a fixed hard-coded subsystem.

  Datastore Management Plane is the first pluggable management plane hosted through the Management Plane Framework.

Architecture:
  SDE Control Plane:
    Components:
      - Control Plane Foundation
      - Core Control Plane
      - Management Plane Framework
      - Datastore Management Plane
      - Optional AI Control Plane
      - SDE Data Plane Interface

  Control Plane Foundation:
    Provides:
      - Foundation Services
      - Foundation Provider binding

  Core Control Plane:
    Provides:
      - Runtime governance
      - Plugin governance
      - Engine metadata governance
      - Capability governance
      - Deployment governance
      - Management plane governance

  Management Plane Framework:
    Provides:
      - Pluggable management plane registration
      - Management plane manifest validation
      - Management plane lifecycle governance
      - Management plane controller runtime integration
      - Management plane API exposure
      - Policy, workflow, audit, and observability integration

  Datastore Management Plane:
    Role:
      - First pluggable management plane.

    Provides:
      - dstoreOps
      - Tenant-scoped Downstream Datastore lifecycle management
      - DatastoreRequest reconciliation
      - DatastoreInstance state management
      - Datastore Operator Plugin integration
      - Infrastructure Provider integration

  Optional AI Control Plane:
    Role:
      - Reserved pluggable Control Plane extension.

    Provides:
      - Future AI-assisted observation, recommendation, validation, tenant assistance, and workflow initiation.

  SDE Data Plane Interface:
    Provides:
      - Approved state publication boundary
      - Runtime configuration access
      - Runtime metadata access
      - Capability metadata access
      - Engine metadata access

Component Model:
  Control Plane Foundation:
    Role:
      - Shared foundation layer for reusable control-plane services.

    Includes:
      - Identity Service
      - Authorization Service
      - Tenant Management Service
      - Configuration Service
      - Policy Service
      - Secrets Service
      - Audit Service
      - Workflow Service
      - Eventing Service
      - Observability Service
      - Registry Framework
      - Plugin Framework

  Core Control Plane:
    Role:
      - Owns platform-level governance metadata.

    Includes:
      - Runtime Registry
      - Plugin Registry
      - Engine Registry
      - Capability Governance
      - Deployment Governance

  Management Plane Framework:
    Role:
      - Provides the host framework for pluggable management planes.

    Includes:
      - Management Plane Registry
      - Management Plane Manifest
      - Management Plane Admission
      - Management Plane Lifecycle
      - Management Plane Controller Runtime integration
      - Management Plane Conformance

  Datastore Management Plane:
    Role:
      - Pluggable management plane for Downstream Datastore lifecycle and operations.

    Includes:
      - dstoreOps
      - Tenant Namespace Manager
      - Datastore Registry
      - Datastore Request Controller
      - Datastore Instance Controller
      - Datastore Operator Registry
      - Infrastructure Provider Registry
      - Lifecycle Controllers

  AI Control Plane:
    Role:
      - Optional pluggable extension, reserved for later scope.

    Boundary:
      - Must use Control Plane APIs, policy, workflow, audit, DMP, and registries.
      - Must not bypass Control Plane governance.

Control Flow:
  Runtime state publication:
    - Core Control Plane maintains approved runtime and plugin metadata.
    - SDE Data Plane consumes approved state through controlled interfaces.
    - SDE Data Plane does not mutate Control Plane authoritative state.

  Management plane registration:
    - Management Plane implementation is registered.
    - Management Plane Manifest is validated.
    - Required Foundation Services are resolved.
    - Management Plane is admitted by policy.
    - Management Plane Controller Runtime is authorized to reconcile resources.

  DMP operation:
    - Tenant or system submits DatastoreRequest.
    - DMP validates request using policy and tenant context.
    - DMP creates or updates DatastoreWorkflow.
    - DMP invokes Datastore Operator Plugin through approved DMP contract.
    - DMP invokes Infrastructure Provider through approved DMP contract when required.
    - DMP records audit, events, metrics, and status.

  Tenant request execution:
    - SDE Data Plane executes client requests.
    - SDE Data Plane uses approved runtime metadata.
    - Engine Plugins access Downstream Datastores for request execution.
    - DMP does not participate in request execution.

Boundaries:
  SDE Control Plane Owns:
    - Authoritative state
    - Registries
    - Governance
    - Policy decisions
    - Workflow coordination
    - Management plane lifecycle
    - Datastore lifecycle through DMP

  SDE Control Plane Does Not Own:
    - Per-request tenant data execution
    - Protocol wire handling
    - Engine execution fragments
    - Downstream native data-plane execution

  SDE Data Plane Owns:
    - Tenant request execution
    - Protocol execution
    - Planning execution
    - Kernel execution
    - Engine execution
    - Result and error propagation

  SDE Data Plane Must Not:
    - Invoke Datastore Operator Plugins
    - Invoke Infrastructure Providers
    - Reconcile DatastoreRequests
    - Mutate Control Plane authoritative state

Invariants:
  - SDE Control Plane is authoritative for platform state.
  - SDE Control Plane hosts pluggable management planes through Management Plane Framework.
  - Datastore Management Plane is the first pluggable management plane.
  - DMP Controller Runtime is not the entire DMP.
  - Datastore lifecycle belongs to DMP, not SDE Data Plane.
  - Tenant data request execution belongs to SDE Data Plane, not DMP.
  - Protocol Plugins do not access Downstream Datastores.
  - Engine Plugins do not manage datastore lifecycle.
  - Datastore Operator Plugins do not execute tenant data-plane requests.
  - Infrastructure Providers do not execute tenant data-plane requests.
  - AI Control Plane remains optional and pluggable until later RFCs define scope.

Related Documents:
  - control-plane-map.md
  - control-plane-foundation.md
  - management-plane.md
  - management-plane-framework/management-plane-framework.md
  - datastore-management-plane/datastore-management-plane.md
  - core-control-plane/core-control-plane.md
  - ai-control-plane.md
