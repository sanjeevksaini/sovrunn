# Protocol Execution

Document:
  ID: protocol-execution
  Title: Protocol Execution
  Parent: sde-data-plane
  Owner: SDE Data Plane
  Layer: SDE Data Plane
  Type: Flow
  Version: 1.0
  Status: Stable

Purpose:
  - Define protocol-side execution behavior inside SDE Data Plane
  - Define Protocol Runtime and Protocol Plugin responsibilities
  - Define protocol request normalization
  - Define protocol response and error mapping
  - Point detailed protocol behavior to focused subflow files

Definition:
  Protocol Execution is the SDE Data Plane flow that accepts a protocol-compatible client request, resolves the appropriate Protocol Plugin, parses protocol input, normalizes protocol intent for SIR Runtime, and maps canonical SDE runtime output back into protocol-compatible response.

  Protocol Execution is an overview flow. Detailed protocol behavior is defined in child subflow files.

Scope:
  In Scope:
    - Protocol request entry
    - Protocol Plugin resolution
    - Protocol normalization
    - Protocol session and transaction context
    - Protocol response mapping
    - Protocol error mapping

  Out of Scope:
    - SIR semantic validation
    - Planning
    - Execution Plan creation
    - Data Kernel orchestration
    - Engine Plugin execution
    - Downstream datastore lifecycle
    - Datastore Data Plane execution

High-Level Flow:
  - Client sends protocol request.
  - Protocol Runtime accepts request.
  - Protocol Runtime resolves approved Protocol Plugin.
  - Protocol Plugin parses protocol input.
  - Protocol Plugin produces protocol-normalized intent.
  - Protocol Runtime resolves session and transaction references where applicable.
  - Protocol Runtime forwards normalized intent to SIR Runtime.
  - SDE Runtime executes request outside protocol flow.
  - Protocol Runtime receives Result Model or Error Model.
  - Protocol Plugin maps canonical output to protocol-compatible response.
  - Client receives response.

Flow Diagram:
  Client
    ↓
  Protocol Runtime
    ↓
  Protocol Plugin
    ↓
  Protocol-Normalized Intent
    ↓
  SIR Runtime
    ↓
  SDE Runtime Execution
    ↓
  Result Model or Error Model
    ↓
  Protocol Plugin
    ↓
  Client

Stage Map:
  Protocol Request Entry:
    Document: protocol-execution/protocol-request-entry.md
    Owner: Protocol Runtime

  Protocol Plugin Resolution:
    Document: protocol-execution/protocol-plugin-resolution.md
    Owner: Protocol Runtime

  Protocol Normalization:
    Document: protocol-execution/protocol-normalization.md
    Owner: Protocol Plugin

  Protocol Session Transaction:
    Document: protocol-execution/protocol-session-transaction.md
    Owner: Protocol Runtime

  Protocol Response Mapping:
    Document: protocol-execution/protocol-response-mapping.md
    Owner: Protocol Plugin

  Protocol Error Mapping:
    Document: protocol-execution/protocol-error-mapping.md
    Owner: Protocol Plugin

Protocol-Normalized Intent:
  Meaning:
    - Intermediate protocol-layer output consumed by SIR Runtime.
    - It is not SIR.
    - It is not Execution Plan.
    - It is not downstream-native operation.

  MUST Include:
    - Protocol identity
    - Protocol operation kind
    - Parsed request structure
    - Protocol semantic modifiers
    - Session reference when applicable
    - Transaction intent when applicable
    - Request metadata
    - Tenant context when available
    - Security context when available

  MUST NOT Include:
    - Raw secrets
    - Engine Plugin invocation data
    - Downstream datastore credentials
    - Downstream-native execution commands
    - SDE Control Plane mutation commands

Protocol Runtime Rules:
  - MUST accept protocol request through approved listener.
  - MUST resolve approved Protocol Plugin.
  - MUST establish request context.
  - MUST manage protocol session boundary.
  - MUST preserve transaction references where applicable.
  - MUST forward normalized intent to SIR Runtime.
  - MUST return protocol-compatible response.
  - MUST NOT execute downstream datastore operations.
  - MUST NOT bypass SIR Runtime.
  - MUST NOT invoke Engine Plugin directly.

Protocol Plugin Rules:
  - MUST parse protocol input.
  - MUST preserve protocol-visible semantics.
  - MUST produce protocol-normalized intent.
  - MUST map Result Model to protocol response.
  - MUST map Error Model to protocol error response.
  - MUST redact unsafe internal details.
  - MUST NOT select Downstream Datastore directly.
  - MUST NOT produce Execution Plan.
  - MUST NOT invoke Engine Runtime or Engine Plugin.
  - MUST NOT manage datastore lifecycle.

Security Rules:
  - Preserve tenant isolation.
  - Preserve protocol session isolation.
  - Protect request context.
  - Protect trace context.
  - Do not expose raw secrets.
  - Do not expose unsafe downstream-native results or errors.

Failure Rules:
  - Protocol parse failure MUST produce Error Model.
  - Protocol Plugin resolution failure MUST produce Error Model.
  - Session or transaction resolution failure MUST fail closed.
  - Result mapping failure MUST produce protocol-compatible internal error.
  - Error mapping failure MUST produce protocol-compatible internal error.
  - Failure MUST NOT be converted into success.

Invariants:
  - Protocol Execution starts at protocol boundary.
  - Protocol Plugin is mandatory for protocol-specific behavior.
  - Protocol Runtime does not execute downstream operations.
  - Protocol Plugin does not produce Execution Plan.
  - Protocol Execution does not bypass SIR Runtime.
  - Raw downstream-native output is never exposed directly.

Relationships:
  Parent:
    - data-plane.md
  Children:
    - protocol-execution/protocol-request-entry.md
    - protocol-execution/protocol-plugin-resolution.md
    - protocol-execution/protocol-normalization.md
    - protocol-execution/protocol-session-transaction.md
    - protocol-execution/protocol-response-mapping.md
    - protocol-execution/protocol-error-mapping.md
  Depends On:
    - data-plane-map.md
    - request-flow.md
    - ../runtime/protocol-runtime.md
    - ../runtime/plugin-runtime.md
    - ../runtime/session-runtime.md
    - ../runtime/transaction-runtime.md
    - ../runtime/result-model.md
    - ../runtime/error-model.md
  Used By:
    - request-flow.md
    - result-propagation.md
    - error-propagation.md
    - Protocol Plugin specifications

References:
  - data-plane.md
  - data-plane-map.md
  - request-flow.md
  - protocol-execution/protocol-request-entry.md
  - protocol-execution/protocol-plugin-resolution.md
  - protocol-execution/protocol-normalization.md
  - protocol-execution/protocol-session-transaction.md
  - protocol-execution/protocol-response-mapping.md
  - protocol-execution/protocol-error-mapping.md
  - ../runtime/protocol-runtime.md
  - ../runtime/result-model.md
  - ../runtime/error-model.md
