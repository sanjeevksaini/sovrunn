# Datastore Management Plane

Document:
  ID: datastore-management-plane
  Title: Datastore Management Plane
  Parent: management-plane-framework
  Owner: Datastore Management Plane
  Layer: SDE Control Plane
  Type: Architecture
  Version: 1.1
  Status: Stable

Purpose:
  - Define Datastore Management Plane architecture
  - Define DMP as the first pluggable management plane inside SDE Control Plane
  - Define dstoreOps relationship
  - Define downstream datastore lifecycle ownership
  - Separate datastore lifecycle from SDE Data Plane execution
  - Define registries, plugins, providers, controllers, and workflows

Definition:
  Datastore Management Plane is a pluggable Management Plane inside SDE Control Plane.

  It powers dstoreOps and manages tenant-scoped Downstream Datastore lifecycle and operations under SDE Control Plane authority.

  It manages lifecycle operations for Downstream Datastores.

  It does not execute client data requests and does not replace Engine Plugins.

Core Principle:
  DMP is the management plane for datastore lifecycle.

  SDE Data Plane is the runtime plane for tenant data request execution.

  These responsibilities must not be collapsed.

Architecture:
  Datastore Management Plane:
    Plugs Into:
      - Management Plane Framework

    Uses:
      - Foundation Services
      - Management Plane Controller Runtime
      - Datastore Operator Plugins
      - Infrastructure Providers

    Children:
      - dstoreOps
      - Tenant Namespace Manager
      - Datastore Registry
      - Datastore Request Controller
      - Datastore Instance Controller
      - Datastore Profile Controller
      - Datastore Policy Controller
      - Datastore Credential Controller
      - Datastore Operator Registry
      - Infrastructure Provider Registry
      - Lifecycle Controller
      - Provisioning Controller
      - Configuration Controller
      - Scaling Controller
      - Backup Controller
      - Restore Controller
      - Patch Controller
      - Upgrade Controller
      - Monitoring Controller
      - Retirement Controller

Component Model:
  dstoreOps:
    Role:
      - Product capability for managed Downstream Datastore operations.

  Tenant Namespace Manager:
    Role:
      - Resolves tenant-scoped management boundaries for datastore resources.

  Datastore Registry:
    Role:
      - Authoritative lifecycle metadata for managed and referenced Downstream Datastores.

  Datastore Request Controller:
    Role:
      - Reconciles DatastoreRequest resources into desired lifecycle workflows.

  Datastore Instance Controller:
    Role:
      - Maintains DatastoreInstance lifecycle and status.

  Datastore Operator Registry:
    Role:
      - Metadata and lifecycle state for Datastore Operator Plugins.

  Infrastructure Provider Registry:
    Role:
      - Metadata and lifecycle state for Infrastructure Providers.

  Lifecycle Controller:
    Role:
      - Coordinates lifecycle workflows across controllers, operator plugins, and infrastructure providers.

  Datastore Operator Plugin:
    Role:
      - Performs datastore-specific lifecycle operations under DMP control.

  Infrastructure Provider:
    Role:
      - Provisions or manages infrastructure substrate under DMP control where required.

DMP Controller Runtime:
  Definition:
    Executable runtime that hosts and reconciles DMP resources, workflows, and plugin interactions.

  Binary:
    - sde-dmp-controller

  Clarification:
    - DMP Controller Runtime runs DMP reconciliation.
    - It is not the entire Datastore Management Plane.
    - DMP also includes contracts, APIs, resources, workflows, registries, controllers, and plugin integrations.

Control Flow:
  Provisioning flow:
    - Tenant or system submits DatastoreRequest.
    - Tenant Namespace Manager resolves namespace.
    - Datastore Request Controller validates request.
    - Policy Service validates authorization, quota, residency, and allowed profiles.
    - Workflow Service creates lifecycle workflow.
    - Lifecycle Controller coordinates provisioning.
    - Infrastructure Provider provisions substrate if required.
    - Datastore Operator Plugin provisions datastore.
    - Datastore Instance Controller updates status.
    - Audit Service records all state-changing actions.

  Operation flow:
    - Tenant or operator requests backup, restore, scale, patch, upgrade, monitor, or retire.
    - DMP validates request through policy.
    - DMP creates DatastoreWorkflow.
    - DMP invokes relevant controllers.
    - Controllers invoke Datastore Operator Plugins and Infrastructure Providers through approved contracts.
    - DMP updates DatastoreInstance and operation status.

Boundaries:
  DMP Owns:
    - Tenant-scoped datastore lifecycle
    - DatastoreRequest reconciliation
    - DatastoreInstance lifecycle state
    - DatastoreProfile and DatastorePolicy application
    - dstoreOps workflows
    - Datastore Operator Plugin use
    - Infrastructure Provider use
    - Backup, restore, scale, patch, upgrade, monitoring, retirement workflows

  DMP Does Not Own:
    - Client protocol parsing
    - SIR generation
    - Execution Plan production
    - Data Kernel execution
    - Engine Plugin execution
    - Tenant data-plane request execution

  Datastore Operator Plugin Owns:
    - Datastore-specific lifecycle API integration.

  Datastore Operator Plugin Does Not Own:
    - DMP workflow authority
    - Policy decision authority
    - Tenant request execution

  Infrastructure Provider Owns:
    - Infrastructure substrate integration.

  Infrastructure Provider Does Not Own:
    - Datastore lifecycle semantics
    - Tenant request execution

AI Agent Boundary:
  Future Tenant AI Agent may submit tenant-scoped datastore workflow requests through approved Control Plane APIs.

  AI Agent Must Not:
    - Invoke DMP internals directly
    - Invoke Datastore Operator Plugins directly
    - Invoke Infrastructure Providers directly
    - Bypass policy, workflow, or audit

  Correct Flow:
    - Tenant AI Agent
    - Tenant-scoped Control Plane API
    - Authorization and Policy
    - Workflow Service
    - Datastore Management Plane
    - Datastore Operator Plugin
    - Infrastructure Provider where required
    - Tenant-specific Downstream Datastore

Invariants:
  - DMP is a pluggable management plane inside SDE Control Plane.
  - DMP plugs into Management Plane Framework.
  - DMP Controller Runtime is not the whole DMP.
  - DMP powers dstoreOps.
  - DMP does not execute tenant data-plane requests.
  - DMP does not replace Engine Plugins.
  - Engine Plugins do not manage datastore lifecycle.
  - Datastore Operator Plugins do not execute tenant data-plane requests.
  - Infrastructure Providers do not execute tenant data-plane requests.
  - All DMP state-changing actions must be authorized, policy-governed, workflow-driven, and audited.

Related Documents:
  - ../management-plane.md
  - ../management-plane-framework/management-plane-framework.md
  - dstoreops.md
  - datastore-operator-plugin.md
  - infrastructure-provider.md
