# SDE Runtime

Document:
  ID: runtime
  Title: SDE Runtime
  Parent: architecture
  Owner: SDE Runtime
  Layer: SDE Data Plane
  Type: Architecture
  Version: 1.0
  Status: Stable

Purpose:
  - Define SDE Runtime architecture
  - Define internal runtime components used by SDE Data Plane
  - Define runtime execution boundaries
  - Define runtime relationship with SDE Control Plane
  - Define runtime relationship with Datastore Data Plane

Definition:
  SDE Runtime is the internal runtime component architecture used by SDE Data Plane to execute client data requests.

  SDE Runtime transforms protocol-specific client intent into SIR, converts valid SIR into Execution Plan, coordinates execution through Data Kernel, invokes Engine Runtime, delegates downstream execution through Engine Plugins, and returns canonical Result Model or Error Model output to Protocol Runtime.

Architecture:
  SDE Runtime:
    Components:
      - Protocol Runtime
      - SIR Runtime
      - Planning
      - Data Kernel
      - Engine Runtime
      - Plugin Runtime
      - Session Runtime
      - Transaction Runtime
      - Capability Registry
      - Execution Plan
      - Execution Context
      - Result Model
      - Error Model

Runtime Position:
  SDE Data Plane:
    Uses:
      - SDE Runtime

  SDE Control Plane:
    Provides:
      - Approved runtime metadata
      - Approved plugin metadata
      - Approved engine metadata
      - Approved capability metadata
      - Runtime configuration
      - Policy context

  Datastore Data Plane:
    Provides:
      - Native execution
      - Native storage
      - Native transaction implementation
      - Native durability

Component Model:
  Protocol Runtime:
    - Handles protocol request lifecycle.

  SIR Runtime:
    - Creates and validates semantic representation.

  Planning:
    - Produces immutable Execution Plan.

  Data Kernel:
    - Coordinates execution of Execution Plan.

  Engine Runtime:
    - Resolves and invokes Engine Plugins.

  Plugin Runtime:
    - Manages runtime plugin loading and lifecycle.

  Session Runtime:
    - Manages SDE session context.

  Transaction Runtime:
    - Manages SDE transaction context.

Execution Model:
  Request path:
    - Protocol Runtime receives client request.
    - Protocol Plugin parses protocol input.
    - SIR Runtime creates valid SIR.
    - Planning produces Execution Plan.
    - Execution Context is attached.
    - Data Kernel coordinates execution.
    - Engine Runtime invokes Engine Plugin.
    - Engine Plugin invokes Downstream Datastore.
    - Result or Error is normalized.
    - Protocol Runtime returns protocol-compatible response.

State Model:
  Runtime may hold:
    - Session context references
    - Transaction context references
    - Execution-scoped context
    - Runtime component state
    - Plugin instance state
    - In-flight execution state
    - Result stream state
    - Error propagation state

  Runtime must not hold:
    - Authoritative SDE Control Plane state
    - Secret material beyond authorized execution scope
    - Downstream datastore administrative state
    - Cross-tenant execution state

Control Plane Interaction:
  Runtime consumes:
    - Configuration
    - Policy
    - Runtime metadata
    - Plugin metadata
    - Engine metadata
    - Capability metadata
    - Tenant metadata

  Runtime MUST:
    - Consume approved state only
    - Use consistent state view per execution
    - Fail deterministically when required state is missing
    - Preserve tenant isolation

  Runtime MUST NOT:
    - Modify SDE Control Plane authoritative state
    - Invent capability metadata
    - Bypass capability governance
    - Bypass plugin governance

Security Model:
  Runtime MUST:
    - Preserve tenant isolation
    - Preserve execution isolation
    - Protect Execution Context
    - Protect session and transaction references
    - Avoid leaking secrets
    - Avoid unsafe downstream-native error exposure

Failure Model:
  Runtime failure MUST:
    - Produce Error Model entry
    - Preserve Trace Identifier
    - Preserve Execution Identifier
    - Preserve Timestamp
    - Preserve retry classification
    - Preserve partial result state where applicable
    - Avoid converting failure into success

Invariants:
  - Runtime executes through SDE Data Plane authority.
  - Runtime does not own SDE Control Plane authoritative state.
  - Runtime does not manage downstream datastore lifecycle.
  - Runtime does not replace Datastore Data Plane.
  - Runtime delegates downstream execution only through Engine Runtime and Engine Plugins.
  - Runtime preserves SIR semantics.

Boundaries:
  Owns:
    - Runtime execution architecture
    - Runtime component interaction
    - Execution Context propagation
    - Execution Plan execution coordination
    - Result and Error normalization

  Does Not Own:
    - SDE Control Plane governance
    - Downstream datastore lifecycle
    - Downstream datastore native storage
    - Datastore Data Plane

Relationships:
  Parent:
    - architecture
  Children:
    - runtime-map.md
    - execution-flow.md
    - session-flow.md
    - transaction-flow.md
    - result-flow.md
    - error-flow.md
    - protocol-runtime.md
    - sir-runtime.md
    - planning.md
    - data-kernel.md
    - engine-runtime.md
    - plugin-runtime.md
    - session-runtime.md
    - transaction-runtime.md
    - capability-registry.md
    - execution-plan.md
    - execution-context.md
    - result-model.md
    - error-model.md
  Depends On:
    - ../control-plane/control-plane.md
    - ../data-plane/data-plane.md
  Used By:
    - SDE Data Plane

References:
  - runtime-map.md
  - execution-flow.md
  - protocol-runtime.md
  - sir-runtime.md
  - planning.md
  - data-kernel.md
  - engine-runtime.md
  - execution-context.md
  - result-model.md
  - error-model.md
  - ../data-plane/data-plane.md
  - ../control-plane/control-plane.md
