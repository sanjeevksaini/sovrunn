# Execution Flow

Document:
  ID: execution-flow
  Title: Execution Flow
  Parent: runtime
  Owner: SDE Runtime
  Layer: SDE Data Plane
  Type: Flow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Execution Flow
  - Define runtime sequence behavior
  - Separate flow behavior from component contracts

Definition:
  Execution Flow defines the end-to-end internal SDE Runtime sequence for executing a client request.

Flow:
  - Protocol Runtime receives client request.
  - Protocol Plugin parses protocol input.
  - SIR Runtime creates and validates SIR.
  - Planning validates capabilities and produces Execution Plan.
  - Execution Context is created.
  - Data Kernel coordinates Execution Plan execution.
  - Engine Runtime resolves Engine Plugin.
  - Engine Plugin invokes Downstream Datastore.
  - Result Model or Error Model is returned.
  - Protocol Runtime returns protocol-compatible response.

Rules:
  - Preserve SIR semantics.
  - Use approved capability metadata.
  - Use Engine Runtime and Engine Plugin boundary.
  - Never bypass Planning, Data Kernel, or Engine Runtime.

Outputs:
  - Protocol response
  - Protocol error response

Failure Handling:
  - Execution failure produces Error Model.
  - Partial execution state must be explicit.
  - Unknown outcome must not be hidden.

Boundaries:
  - This flow defines sequence behavior only.
  - Component ownership remains defined in component documents.
  - Contract structure remains defined in contract documents.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - protocol-runtime.md
    - sir-runtime.md
    - planning.md
    - execution-plan.md
    - execution-context.md
    - data-kernel.md
    - engine-runtime.md
    - result-model.md
    - error-model.md
  Used By:
    - SDE Data Plane
    - SDE Runtime

References:
  - runtime.md
  - runtime-map.md
  - protocol-runtime.md
  - sir-runtime.md
  - planning.md
  - execution-plan.md
  - execution-context.md
  - data-kernel.md
  - engine-runtime.md
  - result-model.md
  - error-model.md
