# Request Flow

Document:
  ID: request-flow
  Title: Request Flow
  Parent: sde-data-plane
  Owner: SDE Data Plane
  Layer: SDE Data Plane
  Type: Flow
  Version: 1.0
  Status: Stable

Purpose:
  - Define end-to-end SDE Data Plane request lifecycle
  - Define request-stage ownership
  - Define boundaries between SDE Data Plane, SDE Runtime, SDE Control Plane, and Datastore Data Plane
  - Point detailed stage behavior to focused subflow files

Definition:
  Request Flow is the end-to-end execution sequence used by SDE Data Plane to process a client data request.

  A request enters through a protocol boundary, is transformed into SIR, planned into an Execution Plan, executed through SDE Runtime, delegated to a Downstream Datastore through an Engine Plugin, and returned as a protocol-compatible response.

  Request Flow is an overview flow. Detailed stage behavior is defined in child subflow files.

Scope:
  In Scope:
    - Request entry
    - SIR creation and validation handoff
    - Planning handoff
    - Execution start
    - Downstream delegation
    - Response return
    - Failure propagation boundaries

  Out of Scope:
    - Runtime component contracts
    - Downstream datastore lifecycle management
    - Datastore Data Plane internals
    - SDE Control Plane authoritative state mutation
    - Detailed protocol parsing rules

High-Level Flow:
  - Client sends protocol-compatible request.
  - Protocol Runtime accepts request.
  - Protocol Plugin normalizes protocol intent.
  - SIR Runtime creates and validates SIR.
  - Planning validates capabilities and produces Execution Plan.
  - Execution Context is created.
  - Data Kernel coordinates execution.
  - Engine Runtime resolves Engine Plugin.
  - Engine Plugin delegates to Downstream Datastore.
  - Datastore Data Plane executes native operation.
  - Engine Plugin maps native output to Result Model or Error Model.
  - Protocol Plugin maps canonical output to protocol-compatible response.
  - Client receives response.

Flow Diagram:
  Client
    ↓
  Protocol Runtime
    ↓
  Protocol Plugin
    ↓
  SIR Runtime
    ↓
  Planning
    ↓
  Execution Plan + Execution Context
    ↓
  Data Kernel
    ↓
  Engine Runtime
    ↓
  Engine Plugin
    ↓
  Downstream Datastore / Datastore Data Plane
    ↓
  Result Model or Error Model
    ↓
  Protocol Plugin
    ↓
  Client

Stage Map:
  Request Entry:
    Document: request-flow/request-entry.md
    Owner: Protocol Runtime
    Output:
      - Request context
      - Protocol-normalized intent

  SIR Creation:
    Document: request-flow/sir-creation.md
    Owner: SIR Runtime
    Output:
      - Validated SIR

  Planning Handoff:
    Document: request-flow/planning-handoff.md
    Owner: Planning
    Output:
      - Execution Plan

  Execution Start:
    Document: request-flow/execution-start.md
    Owner: Data Kernel
    Output:
      - In-flight execution state

  Downstream Delegation:
    Document: request-flow/downstream-delegation.md
    Owner: Engine Runtime
    Output:
      - Result Model
      - Error Model

  Response Return:
    Document: request-flow/response-return.md
    Owner: Protocol Runtime
    Output:
      - Protocol-compatible response

Rules:
  - Request Flow MUST start at a protocol boundary.
  - Request Flow MUST preserve SIR semantics.
  - Request Flow MUST use approved SDE Control Plane state only.
  - Request Flow MUST use immutable Execution Plan.
  - Request Flow MUST use immutable Execution Context.
  - Request Flow MUST delegate downstream execution through Engine Plugin.
  - Request Flow MUST normalize success through Result Model.
  - Request Flow MUST normalize failure through Error Model.
  - Request Flow MUST NOT modify SDE Control Plane authoritative state.
  - Request Flow MUST NOT invoke Datastore Management Plane.
  - Request Flow MUST NOT invoke Datastore Operator Plugins.
  - Request Flow MUST NOT bypass Datastore Data Plane through unmanaged access.

Failure Rules:
  - Any failed request MUST produce Error Model.
  - Failure MUST preserve Trace Identifier.
  - Failure MUST preserve Timestamp.
  - Failure MUST preserve retry classification.
  - Partial result state MUST be explicit.
  - Unknown downstream outcome MUST be reported explicitly.
  - Failure MUST NOT be converted into success.

Security Rules:
  - Preserve tenant isolation.
  - Preserve execution isolation.
  - Protect Execution Context.
  - Protect session and transaction references.
  - Do not expose raw secrets.
  - Do not expose unsafe downstream-native errors.

Invariants:
  - Protocol Plugin is the protocol boundary.
  - SIR Runtime is the semantic representation boundary.
  - Planning is the Execution Plan boundary.
  - Data Kernel is the orchestration boundary.
  - Engine Runtime and Engine Plugin form the downstream execution boundary.
  - Datastore Data Plane remains owned by the Downstream Datastore.

Relationships:
  Parent:
    - data-plane.md
  Children:
    - request-flow/request-entry.md
    - request-flow/sir-creation.md
    - request-flow/planning-handoff.md
    - request-flow/execution-start.md
    - request-flow/downstream-delegation.md
    - request-flow/response-return.md
  Depends On:
    - data-plane-map.md
    - ../runtime/runtime.md
    - ../runtime/runtime-map.md
    - ../runtime/execution-flow.md
    - ../control-plane/control-plane.md
  Used By:
    - protocol-execution.md
    - planning-execution.md
    - kernel-execution.md
    - engine-execution.md
    - result-propagation.md
    - error-propagation.md

References:
  - data-plane.md
  - data-plane-map.md
  - request-flow/request-entry.md
  - request-flow/sir-creation.md
  - request-flow/planning-handoff.md
  - request-flow/execution-start.md
  - request-flow/downstream-delegation.md
  - request-flow/response-return.md
  - ../runtime/runtime.md
  - ../runtime/execution-flow.md
  - ../control-plane/control-plane.md
