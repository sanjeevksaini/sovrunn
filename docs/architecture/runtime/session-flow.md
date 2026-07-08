# Session Flow

Document:
  ID: session-flow
  Title: Session Flow
  Parent: runtime
  Owner: SDE Runtime
  Layer: SDE Data Plane
  Type: Flow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Session Flow
  - Define runtime sequence behavior
  - Separate flow behavior from component contracts

Definition:
  Session Flow defines how SDE Runtime creates, resolves, uses, and expires session context.

Flow:
  - Protocol Runtime receives session-affecting request.
  - Session Runtime creates or resolves session context.
  - Runtime parameters are applied where authorized.
  - Execution Context references session.
  - Data Kernel executes using session reference.
  - Session Runtime updates session-scoped state where applicable.
  - Session expires or closes according to policy.

Rules:
  - Preserve tenant isolation.
  - Do not leak session state across clients.
  - Externalize session context where stateless runtime requires it.

Outputs:
  - Session reference
  - Session error

Failure Handling:
  - Unauthorized session lookup fails closed.
  - Expired session produces deterministic error.

Boundaries:
  - This flow defines sequence behavior only.
  - Component ownership remains defined in component documents.
  - Contract structure remains defined in contract documents.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - session-runtime.md
    - execution-context.md
    - protocol-runtime.md
  Used By:
    - SDE Data Plane
    - SDE Runtime

References:
  - runtime.md
  - runtime-map.md
  - session-runtime.md
  - execution-context.md
  - protocol-runtime.md
