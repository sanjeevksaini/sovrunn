# Protocol Runtime

Document:
  ID: protocol-runtime
  Title: Protocol Runtime
  Parent: runtime
  Owner: Protocol Runtime
  Layer: SDE Data Plane
  Type: Component Architecture
  Version: 1.0
  Status: Stable

Purpose:
  - Define Protocol Runtime
  - Define runtime position
  - Define execution responsibilities
  - Define boundaries and failure rules

Definition:
  Protocol Runtime manages protocol request lifecycle inside SDE Runtime.

Runtime Position:
  - Entry point for protocol-specific client request handling.
  - Uses Protocol Plugins.
  - Calls SIR Runtime.
  - Returns protocol-compatible response.

Responsibilities:
  MUST:
    - Accept client request
    - Resolve Protocol Plugin
    - Manage protocol session boundary
    - Forward semantic intent to SIR Runtime
    - Return protocol-compatible response

  MUST NOT:
    - Execute downstream datastore operations
    - Produce Execution Plan
    - Invoke Engine Plugin directly
    - Bypass SIR Runtime

Inputs:
  - Client protocol request
  - Protocol session context
  - Protocol Plugin metadata

Outputs:
  - Protocol response
  - Protocol error response
  - SIR input

State:
  - Protocol connection state
  - Protocol session reference
  - Protocol plugin binding

Execution Rules:
  - Select approved Protocol Plugin.
  - Preserve client-visible protocol semantics.
  - Forward only normalized protocol intent to SIR Runtime.

Failure Rules:
  - Convert runtime error to protocol-compatible error response.
  - Preserve Trace Identifier.
  - Fail closed when Protocol Plugin is unavailable.

Concurrency Rules:
  - Preserve execution isolation.
  - Avoid shared mutable execution state unless explicitly synchronized.
  - Preserve tenant isolation across concurrent executions.
  - Preserve session and transaction boundaries.

Security Rules:
  - Enforce authorized execution context.
  - Avoid exposing secrets.
  - Preserve safe error behavior.
  - Preserve trace and audit correlation where applicable.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - plugin-runtime.md
    - sir-runtime.md
    - result-model.md
    - error-model.md
  Used By:
    - execution-flow.md
    - session-flow.md
    - result-flow.md
    - error-flow.md

References:
  - runtime.md
  - runtime-map.md
  - plugin-runtime.md
  - sir-runtime.md
  - result-model.md
  - error-model.md
