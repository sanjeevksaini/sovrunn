# Protocol Error Mapping

Document:
  ID: protocol-error-mapping
  Title: Protocol Error Mapping
  Parent: error-propagation
  Owner: Protocol Plugin
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Protocol Error Mapping stage of Error Propagation
  - Keep error-propagation.md concise
  - Define focused error-stage rules

Definition:
  Protocol Error Mapping defines how Protocol Plugin maps Error Model to protocol-compatible error response.

Flow:
  - Protocol Runtime receives Error Model.
  - Protocol Plugin applies protocol error mapping rules.
  - Protocol Plugin maps safe code and category.
  - Protocol Plugin maps safe message.
  - Protocol Plugin preserves trace correlation where allowed.
  - Protocol Plugin redacts unsafe details.
  - Protocol Plugin emits protocol-compatible error response.

Rules:
  - Protocol response MUST preserve safe error semantics.
  - Raw downstream-native error MUST NOT be exposed.
  - Unsafe details MUST be redacted.
  - Failure MUST NOT be represented as success.

Failure Rules:
  - Failure during error handling MUST still produce or preserve Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST NOT expose unsafe details.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.
  - Failure MUST NOT be converted into success.

Relationships:
  Parent:
    - ../error-propagation.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/error-model.md
  Used By:
    - ../error-propagation.md

References:
  - ../error-propagation.md
  - ../data-plane.md
  - ../../runtime/error-model.md
