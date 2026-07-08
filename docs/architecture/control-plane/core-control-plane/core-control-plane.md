# Core Control Plane

Document:
  ID: core-control-plane
  Title: Core Control Plane
  Parent: management-plane
  Owner: Core Control Plane
  Layer: SDE Control Plane
  Type: Architecture
  Version: 1.0
  Status: Stable

Purpose:
  - Define Core Control Plane architecture
  - Define built-in runtime governance Management Plane
  - Define registries and governance components
  - Separate runtime governance from SDE Data Plane execution

Definition:
  Core Control Plane is the built-in Management Plane responsible for SDE runtime governance.

  It owns authoritative metadata for runtime instances, plugins, downstream engine registrations, capability approval, and deployment governance.

Architecture:
  Core Control Plane:
    Children:
      - Runtime Registry
      - Plugin Registry
      - Engine Registry
      - Capability Governance
      - Deployment Governance

Component Model:
  Runtime Registry:
    Role:
      - Tracks SDE runtime instances, topology, health, version, and availability metadata.

  Plugin Registry:
    Role:
      - Tracks plugin metadata and lifecycle state.

  Engine Registry:
    Role:
      - Tracks downstream engine registration metadata used by SDE Data Plane.

  Capability Governance:
    Role:
      - Validates, approves, versions, and publishes Engine Plugin capability metadata.

  Deployment Governance:
    Role:
      - Governs rollout, compatibility, upgrade, downgrade, and rollback.

Control Flow:
  Plugin registration flow:
    - Plugin metadata is submitted.
    - Plugin Registry validates identity and compatibility.
    - Policy Service validates governance rules.
    - Plugin lifecycle state is recorded.
    - Approved plugin metadata becomes discoverable.

  Capability approval flow:
    - Engine Plugin publishes Capability Manifest.
    - Capability Governance validates manifest.
    - Approved capability metadata is published to runtime Capability Registry.
    - Planning consumes approved capability metadata during SDE Data Plane execution.

  Engine registration flow:
    - Engine metadata is registered.
    - Engine Registry validates engine identity, plugin binding, endpoint metadata, and ownership mode.
    - Approved metadata is made available to SDE Data Plane through controlled interface.

State Model:
  Core Control Plane owns:
    - Runtime instance metadata
    - Plugin metadata
    - Engine metadata
    - Capability approval state
    - Deployment state
    - Compatibility state

Security Model:
  Core Control Plane MUST:
    - Enforce authorization for registry mutation
    - Validate plugin identity
    - Validate runtime identity
    - Prevent unapproved plugins from runtime use
    - Prevent unapproved capabilities from Planning use
    - Prevent raw secrets from being stored in registries

Failure Model:
  Core Control Plane failures MUST:
    - Preserve registry consistency
    - Prevent partial capability approval
    - Prevent unsafe deployment state
    - Preserve rollback metadata
    - Avoid publishing unapproved metadata

Invariants:
  - Core Control Plane governs runtime metadata but does not execute runtime requests.
  - Capability Governance approves capabilities but does not invent capabilities.
  - Engine Registry stores execution metadata but does not execute engines.
  - Plugin Registry stores lifecycle metadata but does not execute plugins.
  - SDE Data Plane consumes approved metadata only.

Boundaries:
  Owns:
    - Runtime Registry
    - Plugin Registry
    - Engine Registry
    - Capability Governance
    - Deployment Governance

  Does Not Own:
    - SDE Data Plane request execution
    - Datastore lifecycle management
    - Datastore Data Plane
    - Foundation Provider implementation

Relationships:
  Parent:
    - management-plane.md
  Children:
    - runtime-registry.md
    - plugin-registry.md
    - engine-registry.md
    - capability-governance.md
    - deployment-governance.md
  Depends On:
    - foundation-services/foundation-services.md
    - registry-framework.md
    - plugin-framework.md
  Used By:
    - SDE Data Plane
    - Planning
    - Engine Runtime
    - Protocol Runtime
    - Datastore Management Plane

References:
  - ../management-plane.md
  - runtime-registry.md
  - plugin-registry.md
  - engine-registry.md
  - capability-governance.md
  - deployment-governance.md
  - ../../runtime/runtime.md
  - ../../data-plane/data-plane.md
