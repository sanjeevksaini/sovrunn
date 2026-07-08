# Protocol Response Mapping

Document:
  ID: protocol-response-mapping
  Title: Protocol Response Mapping
  Parent: protocol-execution
  Owner: Protocol Plugin
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Protocol Response Mapping stage of Protocol Execution
  - Keep protocol-execution.md concise
  - Define focused protocol-stage rules

Definition:
  Protocol Response Mapping defines how Result Model becomes protocol-compatible response.

Flow:
  - Protocol Runtime receives Result Model.
  - Protocol Plugin maps canonical fields to protocol response fields.
  - Protocol Plugin maps type and schema metadata where supported.
  - Protocol Plugin preserves cursor or stream references where applicable.
  - Protocol Runtime returns success or partial-success response.

Rules:
  - Raw downstream-native result MUST NOT be exposed.
  - Failed execution MUST NOT be converted into success.
  - Partial result state MUST remain explicit where supported.

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
