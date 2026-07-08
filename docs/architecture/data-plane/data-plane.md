# SDE Data Plane

Document:
  ID: sde-data-plane
  Title: SDE Data Plane
  Parent: architecture
  Owner: SDE Data Plane
  Layer: SDE Data Plane
  Type: Architecture
  Version: 1.0
  Status: Stable

Purpose:
  - Define SDE Data Plane authority
  - Define runtime request execution boundary
  - Define relationship with SDE Runtime
  - Define relationship with SDE Control Plane
  - Define relationship with Datastore Data Plane
  - Define decomposition for detailed flow documents

Definition:
  SDE Data Plane is the runtime request execution plane of Sovrunn Data Engine.

  SDE Data Plane executes client data requests by using SDE Runtime components to transform protocol-specific intent into SIR, produce Execution Plans, coordinate execution, invoke Engine Plugins, normalize results, normalize errors, and return protocol-compatible responses.

  SDE Data Plane does not manage downstream datastore lifecycle.

  SDE Data Plane does not replace Datastore Data Plane.

Architecture:
  SDE Data Plane:
    Uses:
      - SDE Runtime
      - Protocol Plugins
      - Engine Plugins
      - Approved SDE Control Plane state

    Delegates To:
      - Downstream Datastore through Engine Plugin

    Returns:
      - Protocol-compatible response
      - Protocol-compatible error response

Execution Architecture:
  Client Request:
    - Enters through protocol boundary.

  Protocol Layer:
    - Protocol Runtime and Protocol Plugin handle protocol semantics.

  Semantic Layer:
    - SIR Runtime creates and validates SIR.

  Planning Layer:
    - Planning converts valid SIR into Execution Plan using approved capabilities.

  Orchestration Layer:
    - Data Kernel coordinates Execution Plan execution using Execution Context.

  Engine Layer:
    - Engine Runtime resolves Engine Plugin.
    - Engine Plugin translates execution fragment into downstream-native operation.

  Datastore Boundary:
    - Downstream Datastore executes through Datastore Data Plane.

  Response Layer:
    - Engine Plugin maps native result or error into SDE Result Model or Error Model.
    - Protocol Plugin maps canonical output into protocol-compatible response.

Component Relationship:
  SDE Data Plane:
    Internal Runtime:
      - Protocol Runtime
      - SIR Runtime
      - Planning
      - Data Kernel
      - Engine Runtime
      - Plugin Runtime
      - Session Runtime
      - Transaction Runtime

    Runtime Contracts:
      - Execution Plan
      - Execution Context
      - Capability Registry
      - Result Model
      - Error Model

    Plugin Boundaries:
      - Protocol Plugin
      - Engine Plugin

    External Boundaries:
      - SDE Control Plane
      - Downstream Datastore
      - Datastore Data Plane

Control Plane Relationship:
  SDE Control Plane provides approved management state.

  SDE Data Plane consumes:
    - Configuration
    - Policy
    - Tenant metadata
    - Runtime metadata
    - Plugin metadata
    - Engine metadata
    - Capability metadata
    - Datastore endpoint metadata when required

  SDE Data Plane MUST:
    - Consume approved state only
    - Use consistent state view per execution
    - Fail deterministically when required state is unavailable
    - Preserve tenant isolation
    - Preserve policy enforcement

  SDE Data Plane MUST NOT:
    - Modify SDE Control Plane authoritative state
    - Invent management configuration
    - Invent engine metadata
    - Consume unapproved capability metadata
    - Invoke Datastore Management Plane
    - Invoke Datastore Operator Plugins
    - Depend directly on Foundation Providers unless explicitly authorized

Datastore Relationship:
  Downstream Datastore owns:
    - Native execution
    - Native storage
    - Native optimizer
    - Native transaction implementation
    - Native durability
    - Datastore Data Plane

  SDE Data Plane owns:
    - Semantic execution orchestration
    - Runtime execution coordination
    - Engine Plugin delegation
    - Result normalization
    - Error normalization

  Boundary Rule:
    - SDE Data Plane accesses Downstream Datastore only through Engine Plugin and approved downstream interface.

Extension Model:
  Protocol Plugin:
    Role:
      - Integrates client protocol semantics into SDE Data Plane.

  Engine Plugin:
    Role:
      - Integrates SDE Data Plane execution with Downstream Engine.

  Datastore Operator Plugin:
    Role:
      - Integrates Datastore Management Plane with datastore lifecycle operations.

  Infrastructure Provider:
    Role:
      - Integrates Datastore Management Plane with infrastructure environments.

  Boundary Rules:
    - Protocol Plugin participates in SDE Data Plane.
    - Engine Plugin participates in SDE Data Plane.
    - Datastore Operator Plugin does not participate in SDE Data Plane request execution.
    - Infrastructure Provider does not participate in SDE Data Plane request execution.
    - Foundation Provider is not a runtime execution plugin.

Execution Model:
  Request execution follows:
    - Client Request
    - Protocol Runtime
    - SIR Runtime
    - Planning
    - Execution Plan
    - Execution Context
    - Data Kernel
    - Engine Runtime
    - Engine Plugin
    - Downstream Datastore
    - Datastore Data Plane
    - Result Model or Error Model
    - Protocol Response

Security Model:
  SDE Data Plane MUST:
    - Preserve tenant isolation
    - Preserve execution isolation
    - Enforce runtime authorization decisions
    - Protect Execution Context
    - Protect session and transaction references
    - Avoid exposing secrets
    - Avoid exposing unsafe downstream-native errors

Failure Model:
  SDE Data Plane failure MUST:
    - Produce Error Model entry
    - Preserve Trace Identifier
    - Preserve Execution Identifier
    - Preserve Timestamp
    - Preserve retry classification
    - Preserve partial result state where applicable
    - Avoid corrupting SDE Control Plane authoritative state
    - Avoid hiding uncertain execution state

Invariants:
  - SDE Data Plane executes client data requests.
  - SDE Data Plane preserves SIR semantics.
  - SDE Data Plane uses SDE Runtime components.
  - SDE Data Plane consumes approved SDE Control Plane state.
  - SDE Data Plane delegates downstream execution through Engine Plugins.
  - SDE Data Plane does not manage downstream datastore lifecycle.
  - SDE Data Plane does not replace Datastore Data Plane.
  - SDE Data Plane does not bypass Planning, Data Kernel, Engine Runtime, or Engine Plugin boundary.

Boundaries:
  Owns:
    - Runtime request execution authority
    - Cross-component execution flow
    - Protocol-to-SIR execution path
    - SIR-to-Execution Plan execution path
    - Execution Plan orchestration path
    - Engine Plugin delegation path
    - Result propagation
    - Error propagation

  Does Not Own:
    - SDE Control Plane governance
    - SDE Control Plane authoritative state
    - Downstream datastore lifecycle
    - Datastore Operator Plugin lifecycle operations
    - Infrastructure Provider lifecycle operations
    - Datastore Data Plane
    - Downstream native execution semantics

Relationships:
  Parent:
    - architecture
  Children:
    - data-plane-map.md
    - request-flow.md
    - protocol-execution.md
    - planning-execution.md
    - kernel-execution.md
    - engine-execution.md
    - result-propagation.md
    - error-propagation.md
  Depends On:
    - ../runtime/runtime.md
    - ../runtime/runtime-map.md
    - ../control-plane/control-plane.md
    - ../control-plane/control-plane-map.md
  Used By:
    - SDE Runtime documentation
    - Protocol Plugin specifications
    - Engine Plugin specifications
    - Future request-flow documents

References:
  - data-plane-map.md
  - ../runtime/runtime.md
  - ../runtime/runtime-map.md
  - ../runtime/execution-flow.md
  - ../runtime/result-flow.md
  - ../runtime/error-flow.md
  - ../control-plane/control-plane.md
  - ../control-plane/control-plane-map.md
  - request-flow.md
  - protocol-execution.md
  - planning-execution.md
  - kernel-execution.md
  - engine-execution.md
  - result-propagation.md
  - error-propagation.md
