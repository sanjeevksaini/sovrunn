# SIR Runtime

Document
- ID: sir-runtime
- Version: 1.0
- Status: Stable

Purpose
- Define SIR instance runtime
- Define SIR lifecycle
- Define SIR validation
- Define SIR exchange boundaries

Definition

SIR Runtime manages live SIR instances produced by Protocol Runtime and consumed by Planning

Rules

MUST
- Create SIR instances
- Validate SIR instances
- Preserve SIR semantics
- Preserve SIR version compatibility
- Transfer SIR through runtime contracts
- Dispose SIR instances safely

MUST NOT
- Modify semantic intent
- Perform planning
- Execute operations
- Access Engine
- Depend on Protocol

Lifecycle

Create

↓

Validate

↓

Serialize

↓

Transfer

↓

Consume

↓

Dispose

Responsibilities

Create
- Build SIR instance from Protocol Runtime input

Validate
- Validate structure
- Validate version
- Validate resources
- Validate operations
- Validate expressions
- Validate constraints
- Validate capability references

Serialize
- Use SIR serialization rules
- Preserve semantic meaning

Transfer
- Pass SIR to Planning
- Preserve ownership boundary

Dispose
- Release runtime resources
- Preserve audit visibility

Validation

SIR Runtime

MUST
- Reject invalid SIR
- Reject unsupported major version
- Reject unknown required fields
- Produce deterministic validation errors

MUST NOT
- Validate Engine capability
- Validate execution feasibility
- Rewrite SIR for optimization

Ownership

Sovrunn owns
- SIR Runtime
- SIR instance lifecycle
- SIR validation
- SIR transfer rules

SIR Specification owns
- SIR semantics
- SIR structure
- SIR serialization
- SIR versioning

Planning owns
- SIR consumption
- Execution Plan production

References
- runtime.md
- architecture.md
- planning.md
- specifications/sir/sir.md
- specifications/sir/serialization.md
- specifications/sir/versioning.md
- specifications/sir/conformance.md
