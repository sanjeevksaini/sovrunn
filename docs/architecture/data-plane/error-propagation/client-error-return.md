# Client Error Return

Document:
  ID: client-error-return
  Title: Client Error Return
  Parent: error-propagation
  Owner: Protocol Runtime
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Client Error Return stage of Error Propagation
  - Keep error-propagation.md concise
  - Define focused error-stage rules

Definition:
  Client Error Return defines how Protocol Runtime returns final protocol-compatible error response to the client.

Flow:
  - Protocol Plugin returns protocol-compatible error response.
  - Protocol Runtime validates response correlation.
  - Protocol Runtime writes error response to client connection.
  - Protocol Runtime records safe telemetry.
  - Protocol Runtime closes or preserves protocol state according to protocol rules.

Rules:
  - Error response MUST preserve request correlation.
  - Error response MUST accurately represent failure.
  - Response delivery failure MUST be recorded.
  - Response delivery MUST NOT mutate SDE Control Plane authoritative state.

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
