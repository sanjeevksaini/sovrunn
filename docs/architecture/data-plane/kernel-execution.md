# Kernel Execution

Document:
  ID: kernel-execution
  Title: Kernel Execution
  Parent: sde-data-plane
  Owner: SDE Data Plane
  Layer: SDE Data Plane
  Type: Flow
  Version: 1.0
  Status: Stable

Purpose:
  - Define Data Kernel execution behavior inside SDE Data Plane
  - Define how Execution Plan and Execution Context are used
  - Define operation orchestration, dependency handling, and Engine Runtime delegation
  - Define partial result and failure handling boundaries
  - Point detailed kernel behavior to focused subflow files

Definition:
  Kernel Execution is the SDE Data Plane flow in which Data Kernel coordinates execution of an immutable Execution Plan using an immutable Execution Context.

  Data Kernel orchestrates runtime execution.

  Data Kernel does not access Downstream Datastore directly.

  Data Kernel does not invoke Engine Plugins directly.

  Data Kernel does not modify Execution Plan semantics.

Scope:
  In Scope:
    - Execution Plan intake
    - Execution Context intake
    - Execution readiness validation
    - Operation graph orchestration
    - Dependency handling
    - Session and transaction reference usage
    - Engine Runtime delegation
    - Result aggregation
    - Partial failure handling
    - Kernel execution completion

  Out of Scope:
    - Protocol parsing
    - SIR creation
    - Planning
    - Engine Plugin implementation
    - Downstream datastore native execution
    - Datastore lifecycle management
    - SDE Control Plane authoritative state mutation

High-Level Flow:
  - Data Kernel receives Execution Plan.
  - Data Kernel receives Execution Context.
  - Data Kernel validates execution readiness.
  - Data Kernel initializes in-flight execution state.
  - Data Kernel evaluates operation graph dependencies.
  - Data Kernel dispatches executable operations to Engine Runtime.
  - Engine Runtime returns Result Model or Error Model output.
  - Data Kernel aggregates operation output.
  - Data Kernel preserves partial result or failure state.
  - Data Kernel returns execution output to Protocol Runtime path.

Flow Diagram:
  Execution Plan
    ↓
  Execution Context
    ↓
  Data Kernel
    ↓
  Operation Graph
    ↓
  Dependency Evaluation
    ↓
  Engine Runtime Dispatch
    ↓
  Result Model or Error Model
    ↓
  Aggregation
    ↓
  Execution Output

Stage Map:
  Execution Intake:
    Document: kernel-execution/execution-intake.md
    Owner: Data Kernel

  Context Binding:
    Document: kernel-execution/context-binding.md
    Owner: Data Kernel

  Dependency Evaluation:
    Document: kernel-execution/dependency-evaluation.md
    Owner: Data Kernel

  Operation Dispatch:
    Document: kernel-execution/operation-dispatch.md
    Owner: Data Kernel

  Result Aggregation:
    Document: kernel-execution/result-aggregation.md
    Owner: Data Kernel

  Kernel Completion:
    Document: kernel-execution/kernel-completion.md
    Owner: Data Kernel

Kernel Inputs:
  - Execution Plan
  - Execution Context
  - Session reference
  - Transaction reference
  - Runtime configuration view
  - Engine Runtime availability
  - Result Model
  - Error Model

Kernel Outputs:
  Success:
    - Execution result
    - Aggregated Result Model
    - Stream or cursor reference where applicable

  Failure:
    - Error Model
    - Partial execution output where applicable
    - Unknown outcome marker where required

Kernel Rules:
  - Data Kernel MUST consume immutable Execution Plan.
  - Data Kernel MUST consume immutable Execution Context.
  - Data Kernel MUST preserve SIR semantic intent encoded in Execution Plan.
  - Data Kernel MUST preserve operation dependency ordering.
  - Data Kernel MUST coordinate session and transaction references.
  - Data Kernel MUST delegate downstream execution through Engine Runtime.
  - Data Kernel MUST NOT access Downstream Datastore directly.
  - Data Kernel MUST NOT invoke Engine Plugin directly.
  - Data Kernel MUST NOT manage datastore lifecycle.
  - Data Kernel MUST NOT modify Execution Plan semantics.
  - Data Kernel MUST NOT mutate SDE Control Plane authoritative state.

Dependency Rules:
  - Operation dependency graph MUST be respected.
  - Dependent operations MUST NOT run before required predecessors complete.
  - Independent operations MAY run concurrently only when Execution Plan allows it.
  - Failure of a dependency MUST stop dependent operations unless Execution Plan explicitly allows continuation.
  - Dependency evaluation MUST be deterministic for identical Execution Plan and Execution Context.

Session Rules:
  - Data Kernel MAY use session reference from Execution Context.
  - Data Kernel MUST NOT own session lifecycle.
  - Session lookup MUST preserve tenant isolation.
  - Invalid session reference MUST fail closed.
  - Session state mutation MUST follow Session Runtime rules.

Transaction Rules:
  - Data Kernel MAY use transaction reference from Execution Context.
  - Data Kernel MUST NOT replace Transaction Runtime.
  - Data Kernel MUST preserve transaction boundaries.
  - Data Kernel MUST NOT silently emulate unsupported transaction semantics.
  - Transaction uncertainty MUST be reported explicitly.

Delegation Rules:
  - Data Kernel MUST invoke Engine Runtime for downstream execution.
  - Data Kernel MUST pass Execution Context to Engine Runtime.
  - Data Kernel MUST pass only execution fragments allowed by Execution Plan.
  - Data Kernel MUST NOT bypass Engine Runtime.
  - Data Kernel MUST NOT invoke Datastore Operator Plugin.
  - Data Kernel MUST NOT invoke Infrastructure Provider.

Result Rules:
  - Data Kernel MUST accept Result Model from Engine Runtime.
  - Data Kernel MUST aggregate results according to Execution Plan.
  - Data Kernel MUST preserve result ordering where required.
  - Data Kernel MUST preserve partial result state.
  - Data Kernel MUST NOT expose raw downstream-native result.
  - Data Kernel MUST NOT convert failed execution into success.

Failure Rules:
  - Kernel failure MUST produce or propagate Error Model.
  - Kernel failure MUST preserve Trace Identifier.
  - Kernel failure MUST preserve Execution Identifier.
  - Kernel failure MUST preserve partial execution state.
  - Kernel failure MUST report unknown execution outcome explicitly.
  - Kernel failure MUST NOT corrupt SDE Control Plane authoritative state.
  - Kernel failure MUST NOT corrupt downstream datastore lifecycle state.

Concurrency Rules:
  - Data Kernel MUST preserve execution isolation.
  - Concurrent operations MUST respect Execution Plan dependency rules.
  - Shared mutable execution state MUST be avoided or explicitly synchronized.
  - Concurrent execution MUST preserve tenant isolation.
  - Concurrent result streams MUST preserve request and cursor boundaries.

Invariants:
  - Kernel Execution starts with Execution Plan and Execution Context.
  - Kernel Execution ends with Result Model, Error Model, or explicit partial output.
  - Execution Plan remains immutable.
  - Execution Context remains immutable.
  - Data Kernel is the orchestration boundary.
  - Engine Runtime is the downstream execution delegation boundary.
  - Data Kernel never accesses Datastore Data Plane directly.
  - Data Kernel never manages datastore lifecycle.

Relationships:
  Parent:
    - data-plane.md
  Children:
    - kernel-execution/execution-intake.md
    - kernel-execution/context-binding.md
    - kernel-execution/dependency-evaluation.md
    - kernel-execution/operation-dispatch.md
    - kernel-execution/result-aggregation.md
    - kernel-execution/kernel-completion.md
  Depends On:
    - data-plane-map.md
    - request-flow.md
    - planning-execution.md
    - ../runtime/data-kernel.md
    - ../runtime/execution-plan.md
    - ../runtime/execution-context.md
    - ../runtime/engine-runtime.md
    - ../runtime/result-model.md
    - ../runtime/error-model.md
  Used By:
    - request-flow.md
    - engine-execution.md
    - result-propagation.md
    - error-propagation.md

References:
  - data-plane.md
  - data-plane-map.md
  - request-flow.md
  - planning-execution.md
  - kernel-execution/execution-intake.md
  - kernel-execution/context-binding.md
  - kernel-execution/dependency-evaluation.md
  - kernel-execution/operation-dispatch.md
  - kernel-execution/result-aggregation.md
  - kernel-execution/kernel-completion.md
  - ../runtime/data-kernel.md
  - ../runtime/execution-plan.md
  - ../runtime/execution-context.md
  - ../runtime/engine-runtime.md
  - ../runtime/result-model.md
  - ../runtime/error-model.md
