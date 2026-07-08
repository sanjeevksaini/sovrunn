# Protocol Result Mapping

Document:
  ID: protocol-result-mapping
  Title: Protocol Result Mapping
  Parent: result-propagation
  Owner: Protocol Plugin
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Protocol Result Mapping stage of Result Propagation
  - Keep result-propagation.md concise
  - Define focused result-stage rules

Definition:
  Protocol Result Mapping defines how Protocol Plugin maps Result Model to protocol-compatible response.

Flow:
  - Protocol Runtime receives canonical Result Model.
  - Protocol Plugin applies protocol response rules.
  - Protocol Plugin maps values, schema, affected counts, and metadata.
  - Protocol Plugin maps cursor, stream, or continuation behavior.
  - Protocol Plugin emits protocol-compatible response.

Rules:
  - Protocol Plugin MUST preserve protocol-visible semantics.
  - Protocol Plugin MUST not expose unsafe internal metadata.
  - Raw downstream-native result MUST NOT be exposed.
  - Failed execution MUST NOT be converted into success.

Failure Rules:
  - Failure MUST produce or propagate Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST NOT expose raw downstream-native result.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.

Relationships:
  Parent:
    - ../result-propagation.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/result-model.md
    - ../../runtime/error-model.md
  Used By:
    - ../result-propagation.md

References:
  - ../result-propagation.md
  - ../data-plane.md
  - ../../runtime/result-model.md
  - ../../runtime/error-model.md
