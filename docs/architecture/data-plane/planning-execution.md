# Planning Execution

Document:
  ID: planning-execution
  Title: Planning Execution
  Parent: sde-data-plane
  Owner: SDE Data Plane
  Layer: SDE Data Plane
  Type: Flow
  Version: 1.0
  Status: Stable

Purpose:
  - Define SDE Data Plane planning-stage execution behavior
  - Define how validated SIR becomes an Execution Plan
  - Define capability, policy, engine, and client-preference validation
  - Define Planning boundary with SDE Control Plane and SDE Runtime
  - Point detailed planning behavior to focused subflow files

Definition:
  Planning Execution is the SDE Data Plane flow that converts validated SIR into an immutable Execution Plan using approved runtime state, approved capability metadata, policy context, engine metadata, and client preferences when present.

  Planning Execution does not execute operations.

  Planning Execution does not invoke Engine Runtime.

  Planning Execution does not access Downstream Datastore.

Scope:
  In Scope:
    - Validated SIR intake
    - Approved runtime state resolution
    - Capability lookup
    - Capability validation
    - Policy validation
    - Engine candidate resolution
    - Client preference validation
    - Execution strategy selection
    - Execution Plan production
    - Planning failure handling

  Out of Scope:
    - Protocol parsing
    - SIR creation
    - Data Kernel orchestration
    - Engine Runtime delegation
    - Engine Plugin execution
    - Downstream datastore lifecycle management
    - Datastore Data Plane execution
    - SDE Control Plane authoritative state mutation

High-Level Flow:
  - Planning receives validated SIR from SIR Runtime.
  - Planning resolves approved runtime state.
  - Planning reads approved Capability Registry.
  - Planning resolves candidate engines and Engine Plugin bindings.
  - Planning validates required capabilities.
  - Planning validates policy constraints.
  - Planning validates client preferences when present.
  - Planning selects compatible execution strategy.
  - Planning produces immutable Execution Plan.
  - Planning returns Execution Plan to Data Kernel.

Flow Diagram:
  Validated SIR
    ↓
  Planning
    ↓
  Approved Runtime State
    ↓
  Capability Registry
    ↓
  Engine Metadata
    ↓
  Policy Context
    ↓
  Execution Strategy
    ↓
  Execution Plan
    ↓
  Data Kernel

Stage Map:
  SIR Intake:
    Document: planning-execution/sir-intake.md
    Owner: Planning

  Runtime State Resolution:
    Document: planning-execution/runtime-state-resolution.md
    Owner: Planning

  Capability Resolution:
    Document: planning-execution/capability-resolution.md
    Owner: Planning

  Policy Validation:
    Document: planning-execution/policy-validation.md
    Owner: Planning

  Engine Selection:
    Document: planning-execution/engine-selection.md
    Owner: Planning

  Execution Plan Production:
    Document: planning-execution/execution-plan-production.md
    Owner: Planning

Planning Inputs:
  - Validated SIR
  - Tenant context
  - Security context
  - Policy context
  - Runtime configuration
  - Capability Registry
  - Engine metadata
  - Plugin metadata
  - Client preferences when present

Planning Outputs:
  Success:
    - Execution Plan

  Failure:
    - Error Model

Planning Rules:
  - Planning MUST consume validated SIR only.
  - Planning MUST preserve SIR semantic intent.
  - Planning MUST consume approved SDE Control Plane state only.
  - Planning MUST use approved Capability Registry.
  - Planning MUST validate required capabilities.
  - Planning MUST validate policy constraints.
  - Planning MUST validate client preferences when present.
  - Planning MUST produce immutable Execution Plan.
  - Planning MUST NOT execute operations.
  - Planning MUST NOT invoke Data Kernel.
  - Planning MUST NOT invoke Engine Runtime.
  - Planning MUST NOT invoke Engine Plugin.
  - Planning MUST NOT access Downstream Datastore.
  - Planning MUST NOT modify SDE Control Plane authoritative state.

Capability Rules:
  - Required capabilities MUST be validated before Execution Plan emission.
  - Unsupported required capability MUST fail deterministically.
  - Capability downgrade MUST NOT be silent.
  - Capability metadata MUST come from approved Capability Registry.
  - Planning MUST NOT invent capabilities.
  - Planning MUST NOT consume unapproved Engine Plugin capability metadata.

Policy Rules:
  - Planning MUST apply policy constraints before Execution Plan emission.
  - Policy denial MUST prevent Execution Plan emission.
  - Policy decision context MUST be preserved where safe.
  - Planning MUST NOT bypass authorization or policy checks.
  - Planning MUST NOT encode unsafe policy internals into Execution Plan.

Engine Selection Rules:
  - Engine candidates MUST be selected from approved engine metadata.
  - Engine Plugin binding MUST be approved.
  - Client engine preference MUST be validated.
  - Required capability and policy constraints MUST be satisfied.
  - Engine selection MUST be deterministic for identical validated inputs and state view.

Execution Plan Rules:
  - Execution Plan MUST be immutable after production.
  - Execution Plan MUST reference SIR intent.
  - Execution Plan MUST define operation graph.
  - Execution Plan MUST define operation dependencies.
  - Execution Plan MUST define required capabilities.
  - Execution Plan MUST define execution constraints.
  - Execution Plan MUST NOT contain raw secrets.
  - Execution Plan MUST NOT contain SDE Control Plane mutation instructions.
  - Execution Plan MUST NOT expose downstream-native operations as platform contract.

Failure Rules:
  - Invalid planning input MUST produce Error Model.
  - Missing required runtime state MUST fail deterministically.
  - Unsupported capability MUST produce Error Model.
  - Policy denial MUST produce Error Model.
  - Unsupported client preference MUST produce Error Model.
  - No compatible engine candidate MUST produce Error Model.
  - Planning failure MUST NOT emit Execution Plan.
  - Planning failure MUST NOT mutate SDE Control Plane authoritative state.

Security Rules:
  - Planning MUST preserve tenant isolation.
  - Planning MUST preserve security context.
  - Planning MUST not expose secrets.
  - Planning MUST not leak policy internals.
  - Planning MUST not leak unsafe engine metadata.
  - Planning MUST use a consistent approved state view per request.

Invariants:
  - Planning starts with validated SIR.
  - Planning ends with Execution Plan or Error Model.
  - Planning is semantic-preserving.
  - Planning is capability-aware.
  - Planning is policy-aware.
  - Planning is engine-aware but does not execute engines.
  - Planning does not access Datastore Data Plane.
  - Planning does not manage datastore lifecycle.

Relationships:
  Parent:
    - data-plane.md
  Children:
    - planning-execution/sir-intake.md
    - planning-execution/runtime-state-resolution.md
    - planning-execution/capability-resolution.md
    - planning-execution/policy-validation.md
    - planning-execution/engine-selection.md
    - planning-execution/execution-plan-production.md
  Depends On:
    - data-plane-map.md
    - request-flow.md
    - ../runtime/planning.md
    - ../runtime/capability-registry.md
    - ../runtime/execution-plan.md
    - ../runtime/error-model.md
    - ../control-plane/core-control-plane/capability-governance.md
    - ../control-plane/core-control-plane/engine-registry.md
    - ../control-plane/core-control-plane/plugin-registry.md
  Used By:
    - request-flow.md
    - kernel-execution.md
    - engine-execution.md
    - Engine Plugin specifications

References:
  - data-plane.md
  - data-plane-map.md
  - request-flow.md
  - planning-execution/sir-intake.md
  - planning-execution/runtime-state-resolution.md
  - planning-execution/capability-resolution.md
  - planning-execution/policy-validation.md
  - planning-execution/engine-selection.md
  - planning-execution/execution-plan-production.md
  - ../runtime/planning.md
  - ../runtime/capability-registry.md
  - ../runtime/execution-plan.md
  - ../runtime/error-model.md
  - ../control-plane/core-control-plane/capability-governance.md
  - ../control-plane/core-control-plane/engine-registry.md
  - ../control-plane/core-control-plane/plugin-registry.md
