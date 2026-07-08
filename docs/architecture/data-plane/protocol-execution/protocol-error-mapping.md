# Protocol Error Mapping

Document:
  ID: protocol-error-mapping
  Title: Protocol Error Mapping
  Parent: protocol-execution
  Owner: Protocol Plugin
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Protocol Error Mapping stage of Protocol Execution
  - Keep protocol-execution.md concise
  - Define focused protocol-stage rules

Definition:
  Protocol Error Mapping defines how Error Model becomes protocol-compatible error response.

Flow:
  - Protocol Runtime receives Error Model.
  - Protocol Plugin maps safe error code and category.
  - Protocol Plugin maps safe error message.
  - Protocol Plugin preserves trace correlation where allowed.
  - Protocol Plugin redacts unsafe details.
  - Protocol Runtime returns protocol-compatible error response.

Rules:
  - Raw downstream-native error MUST NOT be exposed.
  - Unsafe internal details MUST be redacted.
  - Failure MUST NOT be hidden as success.

Failure Rules:
  - Failure MUST produce or propagate Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST NOT expose unsafe internal details.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.

Relationships:
  Parent:
    - ../protocol-execution.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/protocol-runtime.md
  Used By:
    - ../protocol-execution.md

References:
  - ../protocol-execution.md
  - ../data-plane.md
  - ../../runtime/protocol-runtime.md
