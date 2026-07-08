# Constraints

Document
- ID: sir-constraints
- Version: 1.0
- Status: Stable

Purpose
- Define semantic constraints
- Define constraint invariants
- Define semantic validation boundaries

Rules

MUST
- Represent semantic validity
- Remain protocol independent
- Remain engine independent
- Be deterministic

MUST NOT
- Represent authorization
- Represent engine limitations
- Represent storage limitations
- Represent runtime validation

Definition

A Constraint defines a semantic rule that MUST be satisfied for a SIR construct to be considered valid.

Ownership

Sovrunn owns
- Semantic constraints
- Constraint model
- Constraint invariants

Planner owns
- Runtime validation

Engine owns
- Physical constraints

Constraint Categories

Structural

- Resource validity
- Relationship validity
- Expression validity
- Operation validity

Behavioral

- Capability requirement
- Composition validity
- Semantic consistency

Compatibility

- Version compatibility
- Extension compatibility

Constraint Scope

Resource

Relationship

Expression

Operation

Metadata

Capability

Constraint Properties

Identity
- Stable
- Versioned

Target
- SIR Construct

Severity
- Error
- Warning

Characteristics

MUST
- Be declarative
- Be deterministic
- Be independently evaluable
- Produce deterministic outcome

MAY
- Reference Resources
- Reference Expressions
- Reference Capabilities
- Compose with other Constraints

Constraint Invariants

Every Constraint

MUST
- Target one or more SIR constructs
- Preserve semantic correctness
- Be version compatible

MUST NOT
- Depend on protocol
- Depend on engine
- Depend on runtime implementation

Violation

Violation of a Constraint

MUST
- Produce deterministic failure
- Identify violated Constraint
- Preserve validation consistency

References

- sir.md
- concepts.md
- resources.md
- relationships.md
- expressions.md
- operations.md
- capability-model.md
