# Request Entry

Document:
  ID: request-entry
  Title: Request Entry
  Parent: request-flow
  Owner: Protocol Runtime
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Request Entry stage of SDE Data Plane Request Flow
  - Keep request-flow.md concise
  - Define focused stage rules

Definition:
  Request Entry defines how SDE Data Plane accepts a client protocol request and creates request-scoped context.

Flow:
  - Client sends protocol-compatible request.
  - Protocol Runtime accepts request through approved protocol listener.
  - Protocol Runtime resolves approved Protocol Plugin.
  - Protocol Plugin parses protocol input.
  - Protocol Plugin produces protocol-normalized intent.
  - Protocol Runtime creates request context.

Rules:
  - Request MUST enter through approved protocol boundary.
  - Protocol Runtime MUST preserve Request Identifier and Trace Identifier.
  - Protocol Plugin MUST NOT select Downstream Datastore.
  - Protocol Plugin MUST NOT produce Execution Plan.
  - Malformed input MUST fail deterministically.

Failure Rules:
  - Failure MUST produce or propagate Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.
  - Failure MUST NOT hide partial or uncertain execution state.

Relationships:
  Parent:
    - ../request-flow.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/runtime.md
  Used By:
    - ../request-flow.md

References:
  - ../request-flow.md
  - ../data-plane.md
  - ../../runtime/runtime.md
