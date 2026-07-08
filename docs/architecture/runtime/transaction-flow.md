# Transaction Flow

Document:
  ID: transaction-flow
  Title: Transaction Flow
  Parent: runtime
  Owner: SDE Runtime
  Layer: SDE Data Plane
  Type: Flow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Transaction Flow
  - Define runtime sequence behavior
  - Separate flow behavior from component contracts

Definition:
  Transaction Flow defines how SDE Runtime handles transaction context and lifecycle.

Flow:
  - Protocol Runtime receives transaction operation.
  - Transaction Runtime creates or resolves transaction context.
  - Execution Context references transaction.
  - Data Kernel coordinates operations under transaction context.
  - Engine Runtime delegates transaction-relevant execution to Engine Plugin.
  - Commit or rollback intent is coordinated.
  - Transaction Runtime records final transaction state.

Rules:
  - Do not silently emulate unsupported transaction semantics.
  - Preserve explicit transaction boundaries.
  - Report uncertain transaction outcome explicitly.

Outputs:
  - Transaction state
  - Transaction error

Failure Handling:
  - Unsupported transaction capability fails deterministically.
  - Commit or rollback uncertainty must be reported through Error Model.

Boundaries:
  - This flow defines sequence behavior only.
  - Component ownership remains defined in component documents.
  - Contract structure remains defined in contract documents.

Relationships:
  Parent:
    - runtime.md
  Depends On:
    - transaction-runtime.md
    - execution-context.md
    - engine-runtime.md
    - error-model.md
  Used By:
    - SDE Data Plane
    - SDE Runtime

References:
  - runtime.md
  - runtime-map.md
  - transaction-runtime.md
  - execution-context.md
  - engine-runtime.md
  - error-model.md
