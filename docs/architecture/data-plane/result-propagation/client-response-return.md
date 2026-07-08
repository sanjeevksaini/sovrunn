# Client Response Return

Document:
  ID: client-response-return
  Title: Client Response Return
  Parent: result-propagation
  Owner: Protocol Runtime
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Client Response Return stage of Result Propagation
  - Keep result-propagation.md concise
  - Define focused result-stage rules

Definition:
  Client Response Return defines how Protocol Runtime returns the final successful or partial response to the client.

Flow:
  - Protocol Plugin returns protocol-compatible response.
  - Protocol Runtime validates response correlation.
  - Protocol Runtime writes response to client connection.
  - Protocol Runtime records safe telemetry.
  - Protocol Runtime closes or preserves protocol state according to protocol rules.

Rules:
  - Response MUST preserve request correlation.
  - Response MUST accurately represent success or partial success.
  - Response delivery failure MUST be recorded.
  - Response delivery MUST NOT mutate SDE Control Plane authoritative state.

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
