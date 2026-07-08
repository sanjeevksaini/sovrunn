# Runtime Map

Document:
  ID: runtime-map
  Title: Runtime Map
  Parent: architecture
  Owner: SDE Runtime
  Layer: SDE Data Plane
  Type: Map
  Version: 1.0
  Status: Stable

Purpose:
  - Define runtime document graph
  - Define runtime component relationships
  - Define runtime flow relationships
  - Define AI retrieval navigation
  - Prevent relationship inference from filenames alone

Runtime Hierarchy:
  SDE Runtime:
    Document: runtime.md
    Children:
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

Component Map:
  Protocol Runtime:
    Document: protocol-runtime.md
    Role:
      - Accepts client protocol requests
      - Selects Protocol Plugin
      - Produces protocol-compatible responses

  SIR Runtime:
    Document: sir-runtime.md
    Role:
      - Creates and validates SIR instances
      - Preserves semantic intent

  Planning:
    Document: planning.md
    Role:
      - Converts valid SIR into Execution Plan
      - Uses Capability Registry

  Data Kernel:
    Document: data-kernel.md
    Role:
      - Coordinates Execution Plan execution
      - Invokes Engine Runtime

  Engine Runtime:
    Document: engine-runtime.md
    Role:
      - Resolves and invokes Engine Plugins

  Plugin Runtime:
    Document: plugin-runtime.md
    Role:
      - Loads and governs runtime plugin lifecycle

  Session Runtime:
    Document: session-runtime.md
    Role:
      - Manages runtime session context

  Transaction Runtime:
    Document: transaction-runtime.md
    Role:
      - Manages SDE transaction context and lifecycle

Contract Map:
  Execution Plan:
    Document: execution-plan.md
    Role:
      - Immutable runtime execution contract

  Execution Context:
    Document: execution-context.md
    Role:
      - Immutable execution-scoped context

  Capability Registry:
    Document: capability-registry.md
    Role:
      - Runtime-facing approved capability lookup

  Result Model:
    Document: result-model.md
    Role:
      - Canonical SDE runtime result representation

  Error Model:
    Document: error-model.md
    Role:
      - Canonical SDE runtime error representation

Flow Map:
  Execution Flow:
    Document: execution-flow.md
    Path:
      - Protocol Runtime
      - SIR Runtime
      - Planning
      - Data Kernel
      - Engine Runtime
      - Engine Plugin
      - Downstream Datastore

  Session Flow:
    Document: session-flow.md
    Path:
      - Protocol Runtime
      - Session Runtime
      - Execution Context
      - Data Kernel

  Transaction Flow:
    Document: transaction-flow.md
    Path:
      - Protocol Runtime
      - Transaction Runtime
      - Execution Context
      - Data Kernel
      - Engine Runtime

  Result Flow:
    Document: result-flow.md
    Path:
      - Downstream Native Result
      - Engine Plugin
      - Result Model
      - Protocol Response

  Error Flow:
    Document: error-flow.md
    Path:
      - Runtime or Downstream Error
      - Error Model
      - Protocol Error Response

Execution Boundary Map:
  SDE Runtime:
    Owns:
      - Runtime execution components
      - Execution Context propagation
      - Execution Plan execution coordination
      - Runtime result and error normalization

  SDE Control Plane:
    Owns:
      - Authoritative management state
      - Runtime metadata governance
      - Capability governance
      - Plugin metadata governance
      - Engine metadata governance

  Datastore Data Plane:
    Owns:
      - Downstream native execution
      - Downstream native storage
      - Downstream native durability
      - Downstream native transaction implementation

Rules:
  - Runtime docs define internal components used by SDE Data Plane.
  - SDE Data Plane docs define outside-in execution-plane authority and request lifecycle.
  - Runtime component docs MUST NOT become full request-flow docs.
  - Flow docs MUST NOT redefine component ownership.
  - Contract docs MUST define one runtime contract only.
  - Runtime MUST consume approved SDE Control Plane state.
  - Runtime MUST NOT modify SDE Control Plane authoritative state.

Navigation:
  Use runtime.md when:
    - You need overall SDE Runtime architecture.

  Use runtime-map.md when:
    - You need relationships, ownership, and retrieval routing.

  Use component docs when:
    - You need one runtime component responsibility.

  Use contract docs when:
    - You need one immutable model or registry contract.

  Use flow docs when:
    - You need end-to-end runtime sequence behavior.

References:
  - runtime.md
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
  - ../data-plane/data-plane.md
  - ../control-plane/control-plane.md
