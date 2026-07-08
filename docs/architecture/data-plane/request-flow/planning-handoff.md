# Planning Handoff

Document:
  ID: planning-handoff
  Title: Planning Handoff
  Parent: request-flow
  Owner: Planning
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Planning Handoff stage of SDE Data Plane Request Flow
  - Keep request-flow.md concise
  - Define focused stage rules

Definition:
  Planning Handoff defines how validated SIR becomes an Execution Plan.

Flow:
  - Planning receives validated SIR.
  - Planning resolves approved runtime state.
  - Planning reads Capability Registry.
  - Planning validates capability requirements.
  - Planning validates client preferences and policy.
  - Planning produces immutable Execution Plan.

Rules:
  - Planning MUST consume approved capability metadata only.
  - Planning MUST NOT silently downgrade required capability.
  - Planning MUST NOT execute operations.
  - Execution Plan MUST preserve SIR intent.

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
