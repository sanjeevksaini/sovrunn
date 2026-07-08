# SDE Data Plane Map

Document:
  ID: data-plane-map
  Title: SDE Data Plane Map
  Parent: architecture
  Owner: SDE Data Plane
  Layer: SDE Data Plane
  Type: Map
  Version: 1.0
  Status: Stable

Purpose:
  - Define SDE Data Plane document graph
  - Define execution-plane relationships
  - Define boundary relationships with SDE Control Plane and Datastore Data Plane
  - Define AI retrieval navigation
  - Prevent relationship inference from filenames alone

Data Plane Hierarchy:
  SDE Data Plane:
    Document: data-plane.md
    Children:
      - Request Flow
      - Protocol Execution
      - Planning Execution
      - Kernel Execution
      - Engine Execution
      - Result Propagation
      - Error Propagation

Architecture Documents:
  SDE Data Plane:
    Document: data-plane.md
    Type: Architecture
    Role:
      - Defines execution-plane authority
      - Defines boundaries
      - Defines component interaction at data-plane level

Flow Documents:
  Request Flow:
    Document: request-flow.md
    Type: Flow
    Role:
      - End-to-end client request lifecycle across SDE Data Plane

  Protocol Execution:
    Document: protocol-execution.md
    Type: Flow
    Role:
      - Protocol request handling and protocol response behavior

  Planning Execution:
    Document: planning-execution.md
    Type: Flow
    Role:
      - SIR to Execution Plan behavior inside SDE Data Plane

  Kernel Execution:
    Document: kernel-execution.md
    Type: Flow
    Role:
      - Data Kernel orchestration behavior

  Engine Execution:
    Document: engine-execution.md
    Type: Flow
    Role:
      - Engine Runtime and Engine Plugin execution behavior

  Result Propagation:
    Document: result-propagation.md
    Type: Flow
    Role:
      - Result flow from Datastore Data Plane to protocol response

  Error Propagation:
    Document: error-propagation.md
    Type: Flow
    Role:
      - Error flow from runtime or downstream failure to protocol error response

Runtime Dependency Map:
  SDE Data Plane:
    Uses:
      - SDE Runtime

  SDE Runtime:
    Provides:
      - Protocol Runtime
      - SIR Runtime
      - Planning
      - Data Kernel
      - Engine Runtime
      - Plugin Runtime
      - Session Runtime
      - Transaction Runtime
      - Execution Plan
      - Execution Context
      - Result Model
      - Error Model

Control Plane Dependency Map:
  SDE Data Plane:
    Consumes:
      - Approved runtime configuration
      - Approved policy context
      - Approved plugin metadata
      - Approved engine metadata
      - Approved capability metadata
      - Tenant metadata
      - Runtime topology metadata

  SDE Control Plane:
    Owns:
      - Authoritative management state
      - Capability governance
      - Engine registry
      - Plugin registry
      - Runtime registry
      - Policy and configuration authority

Datastore Boundary Map:
  SDE Data Plane:
    Owns:
      - Client request execution
      - Semantic execution orchestration
      - Runtime execution coordination
      - Result and error normalization
      - Engine Plugin delegation

  Datastore Data Plane:
    Owns:
      - Native datastore execution
      - Native storage
      - Native optimizer
      - Native transaction implementation
      - Native durability
      - Native access path

  Boundary:
    - SDE Data Plane reaches Datastore Data Plane only through Engine Plugin and downstream datastore interface.

Execution Boundary Map:
  SDE Data Plane:
    MUST:
      - Execute client data requests
      - Preserve SIR semantics
      - Use approved SDE Control Plane state
      - Use SDE Runtime components
      - Delegate downstream execution through Engine Plugins
      - Normalize results through Result Model
      - Normalize failures through Error Model

    MUST NOT:
      - Manage downstream datastore lifecycle
      - Invoke Datastore Management Plane
      - Invoke Datastore Operator Plugins
      - Invoke Infrastructure Providers
      - Modify SDE Control Plane authoritative state
      - Replace Datastore Data Plane
      - Bypass Planning, Data Kernel, Engine Runtime, or Engine Plugins

Relationship Map:
  Parent:
    - docs/architecture/architecture.md

  Depends On:
    - docs/architecture/runtime/runtime.md
    - docs/architecture/runtime/runtime-map.md
    - docs/architecture/control-plane/control-plane.md
    - docs/architecture/control-plane/control-plane-map.md
    - docs/foundation/glossary.md
    - docs/foundation/ontology.md
    - docs/foundation/ownership.md

  Used By:
    - SDE architecture documentation
    - Runtime execution documentation
    - Future protocol and engine plugin specifications

Navigation:
  Use data-plane.md when:
    - You need SDE Data Plane authority, architecture, and boundaries.

  Use data-plane-map.md when:
    - You need document relationships, boundary maps, and retrieval routing.

  Use request-flow.md when:
    - You need end-to-end request lifecycle.

  Use protocol-execution.md when:
    - You need protocol-side data-plane behavior.

  Use planning-execution.md when:
    - You need SIR-to-plan data-plane behavior.

  Use kernel-execution.md when:
    - You need Data Kernel orchestration sequence.

  Use engine-execution.md when:
    - You need Engine Runtime and Engine Plugin execution sequence.

  Use result-propagation.md when:
    - You need success and partial-result propagation.

  Use error-propagation.md when:
    - You need failure propagation and protocol error mapping.

Rules:
  - data-plane.md MUST remain an architecture document, not a full request-flow catalog.
  - Flow documents MUST define sequence behavior.
  - Runtime component contracts MUST remain in docs/architecture/runtime.
  - SDE Data Plane documents MUST NOT duplicate runtime component contracts.
  - SDE Data Plane documents MUST focus on execution-plane authority and cross-component flow.
  - Relationships MUST be explicit in metadata and relationship sections.

References:
  - data-plane.md
  - ../runtime/runtime.md
  - ../runtime/runtime-map.md
  - ../control-plane/control-plane.md
  - ../control-plane/control-plane-map.md
  - request-flow.md
  - protocol-execution.md
  - planning-execution.md
  - kernel-execution.md
  - engine-execution.md
  - result-propagation.md
  - error-propagation.md
