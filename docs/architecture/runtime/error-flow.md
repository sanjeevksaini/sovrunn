# Error Flow

Document:
  ID: error-flow
  Title: Error Flow
  Parent: runtime
  Owner: SDE Runtime
  Layer: SDE Data Plane
  Type: Flow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Error Flow
  - Define runtime sequence behavior
  - Separate flow behavior from component contracts

Definition:
  Error Flow defines how runtime and downstream errors propagate through SDE Runtime.

Flow:
  - Error is detected by runtime component or Engine Plugin.
  - Error is converted to Error Model.
  - Trace Identifier and Timestamp are preserved.
  - Error propagates to caller component.
  - Data Kernel preserves partial state where applicable.
  - Protocol Runtime receives Error Model.
  - Protocol Plugin maps Error Model to protocol error response.
  - Client receives protocol-compatible error.

Rules:
  - Never expose unsafe native error directly.
  - Never convert failure into success.
  - Preserve causal chain safely.
  - Preserve retry classification.

Outputs:
  - Error Model
  - Protocol error response

Failure Handling:
  - Unknown error still produces Error Model.
  - Unsafe details are redacted.
  - Timestamp is mandatory.

Boundaries:
  - This flow defines sequence behavior only.
  - Component ownership remains defined in component documents.
  - Contract structure remains defined in contract documents.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - error-model.md
    - engine-runtime.md
    - data-kernel.md
    - protocol-runtime.md
  Used By:
    - SDE Data Plane
    - SDE Runtime

References:
  - runtime.md
  - runtime-map.md
  - error-model.md
  - engine-runtime.md
  - data-kernel.md
  - protocol-runtime.md
