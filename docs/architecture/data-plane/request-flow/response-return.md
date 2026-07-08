# Response Return

Document:
  ID: response-return
  Title: Response Return
  Parent: request-flow
  Owner: Protocol Runtime
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Response Return stage of SDE Data Plane Request Flow
  - Keep request-flow.md concise
  - Define focused stage rules

Definition:
  Response Return defines how canonical runtime output becomes a client-visible response.

Flow:
  - Data Kernel returns Result Model or Error Model.
  - Protocol Runtime receives canonical output.
  - Protocol Plugin maps Result Model to protocol success response.
  - Protocol Plugin maps Error Model to protocol error response.
  - Protocol Runtime returns protocol-compatible output.
  - Runtime telemetry records response completion.

Rules:
  - Failure MUST NOT be represented as success.
  - Protocol Plugin MUST redact unsafe details.
  - Raw downstream-native result or error MUST NOT be exposed.
  - Response MUST preserve request correlation.

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
