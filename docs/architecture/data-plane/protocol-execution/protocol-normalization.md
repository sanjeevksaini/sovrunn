# Protocol Normalization

Document:
  ID: protocol-normalization
  Title: Protocol Normalization
  Parent: protocol-execution
  Owner: Protocol Plugin
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Protocol Normalization stage of Protocol Execution
  - Keep protocol-execution.md concise
  - Define focused protocol-stage rules

Definition:
  Protocol Normalization defines how protocol-specific payload becomes protocol-normalized intent.

Flow:
  - Protocol Plugin parses protocol payload.
  - Protocol Plugin validates protocol version and feature flags.
  - Protocol Plugin preserves protocol-visible semantic modifiers.
  - Protocol Plugin creates protocol-normalized intent.
  - Protocol Runtime forwards normalized intent to SIR Runtime.

Rules:
  - Protocol-normalized intent is not SIR.
  - Protocol-normalized intent is not Execution Plan.
  - Protocol Plugin MUST NOT bind directly to Downstream Datastore.

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
