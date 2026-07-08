# Runtime State Resolution

Document:
  ID: runtime-state-resolution
  Title: Runtime State Resolution
  Parent: planning-execution
  Owner: Planning
  Layer: SDE Data Plane
  Type: Subflow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Runtime State Resolution stage of Planning Execution
  - Keep planning-execution.md concise
  - Define focused planning-stage rules

Definition:
  Runtime State Resolution defines how Planning obtains an approved and consistent state view for planning.

Flow:
  - Planning resolves tenant context.
  - Planning resolves runtime configuration.
  - Planning resolves policy context.
  - Planning resolves approved engine metadata.
  - Planning resolves approved plugin metadata.
  - Planning creates a consistent planning state view.

Rules:
  - Planning MUST consume approved SDE Control Plane state only.
  - Planning MUST use one consistent state view per request.
  - Planning MUST fail when required state is unavailable.
  - Planning MUST NOT modify SDE Control Plane authoritative state.
  - Planning MUST NOT invent engine or plugin metadata.

Failure Rules:
  - Failure MUST produce or propagate Error Model.
  - Failure MUST preserve Trace Identifier where available.
  - Failure MUST NOT emit Execution Plan when planning validation fails.
  - Failure MUST NOT mutate SDE Control Plane authoritative state.

Relationships:
  Parent:
    - ../planning-execution.md
  Depends On:
    - ../data-plane.md
    - ../../runtime/planning.md
    - ../../runtime/error-model.md
  Used By:
    - ../planning-execution.md

References:
  - ../planning-execution.md
  - ../data-plane.md
  - ../../runtime/planning.md
  - ../../runtime/error-model.md
