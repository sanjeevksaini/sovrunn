# Result Flow

Document:
  ID: result-flow
  Title: Result Flow
  Parent: runtime
  Owner: SDE Runtime
  Layer: SDE Data Plane
  Type: Flow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Result Flow
  - Define runtime sequence behavior
  - Separate flow behavior from component contracts

Definition:
  Result Flow defines how successful and partial results propagate through SDE Runtime.

Flow:
  - Downstream Datastore returns native result.
  - Engine Plugin maps native result to Result Model.
  - Engine Runtime returns Result Model.
  - Data Kernel aggregates or streams result.
  - Protocol Runtime receives canonical result.
  - Protocol Plugin maps Result Model to protocol response.
  - Client receives protocol-compatible response.

Rules:
  - Do not expose raw native result directly.
  - Preserve type and schema metadata.
  - Preserve partial or streaming state.

Outputs:
  - Result Model
  - Protocol response

Failure Handling:
  - Invalid result mapping produces Error Model.
  - Partial result must remain explicit.

Boundaries:
  - This flow defines sequence behavior only.
  - Component ownership remains defined in component documents.
  - Contract structure remains defined in contract documents.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - result-model.md
    - engine-runtime.md
    - data-kernel.md
    - protocol-runtime.md
  Used By:
    - SDE Data Plane
    - SDE Runtime

References:
  - runtime.md
  - runtime-map.md
  - result-model.md
  - engine-runtime.md
  - data-kernel.md
  - protocol-runtime.md
