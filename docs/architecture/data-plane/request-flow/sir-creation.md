# SIR Creation

Document:
  ID: sir-creation
  Title: SIR Creation
  Parent: request-flow
  Owner: SIR Runtime
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define SIR Creation stage of SDE Data Plane Request Flow
  - Keep request-flow.md concise
  - Define focused stage rules

Definition:
  SIR Creation defines how protocol-normalized intent becomes validated SIR.

Flow:
  - Protocol Runtime forwards protocol-normalized intent.
  - SIR Runtime creates SIR instance.
  - SIR Runtime validates SIR structure.
  - SIR Runtime validates SIR version.
  - SIR Runtime validates semantic references.
  - SIR Runtime emits validated SIR to Planning.

Rules:
  - SIR Runtime MUST preserve semantic intent.
  - SIR Runtime MUST NOT optimize by changing semantics.
  - Invalid SIR MUST NOT reach Planning.
  - Unsupported semantic construct MUST fail deterministically.

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
