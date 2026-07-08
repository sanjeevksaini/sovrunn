# Planning

Document
- ID: planning
- Version: 1.0
- Status: Stable

Purpose
- Define Planning
- Define planning responsibilities
- Define planning boundaries
- Define transformation from SIR to Execution Plan

Definition

Planning transforms validated SIR instances into Execution Plans while preserving semantic intent.

Principles

MUST

- Preserve semantic intent
- Produce deterministic Execution Plans
- Be protocol independent
- Be engine independent
- Use Capability Registry
- Preserve SIR boundaries

MUST NOT

- Modify SIR semantics
- Execute operations
- Access downstream engines
- Perform physical execution
- Implement protocol logic

Responsibilities

- Consume validated SIR
- Resolve required capabilities
- Validate capability availability
- Transform semantic operations
- Produce Execution Plan
- Preserve semantic equivalence

Planning Lifecycle

Receive SIR

↓

Validate Planning Preconditions

↓

Resolve Capabilities

↓

Transform Semantic Operations

↓

Generate Execution Plan

↓

Validate Execution Plan

↓

Transfer Execution Plan

Planning Inputs

- SIR Instance
- Capability Registry

Planning Output

- Execution Plan

Capability Interaction

Planning

MUST

- Query Capability Registry
- Resolve required capabilities
- Detect unsupported capabilities
- Preserve capability contracts

MUST NOT

- Discover capabilities
- Modify capability definitions
- Invent capabilities

Transformation Rules

Planning

MUST

- Preserve Resources
- Preserve Relationships
- Preserve Expressions
- Preserve Operations
- Preserve Constraints

MUST NOT

- Rewrite semantic meaning
- Remove required semantics
- Introduce undefined behavior

Execution Plan

Planning

MUST

- Produce deterministic plan
- Preserve semantic equivalence
- Be version compatible

MUST NOT

- Execute plan
- Optimize plan
- Bind to engine implementation

Error Handling

Planning

MUST

- Detect planning failures
- Detect unsupported capabilities
- Produce deterministic errors

MUST NOT

- Produce partial Execution Plans
- Modify SIR to avoid failure

Runtime Characteristics

MUST

- Be stateless
- Be deterministic
- Support concurrent planning
- Scale horizontally

Ownership

Sovrunn owns

- Planning
- Planning lifecycle
- SIR transformation
- Execution Plan generation

Capability Registry owns

- Capability discovery
- Capability metadata

Execution Plan owns

- Runtime execution model

References

- architecture.md
- runtime.md
- sir-runtime.md
- capability-registry.md
- execution-plan.md
- specifications/sir/sir.md
- specifications/sir/capability-model.md
- specifications/sir/capabilities.md
