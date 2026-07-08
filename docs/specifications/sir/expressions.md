# Expressions

Document
- ID: sir-expressions
- Version: 1.0
- Status: Stable

Purpose
- Define semantic expressions
- Define expression ownership
- Define expression boundaries

Rules

MUST
- Represent declarative semantic logic
- Remain protocol independent
- Remain engine independent
- Reference adopted expression models

MUST NOT
- Represent protocol syntax
- Represent engine syntax
- Represent execution plan
- Redefine Substrait expressions

Definition

An Expression is declarative semantic logic used by SIR constructs

Ownership

Substrait owns
- Relational expressions
- Scalar functions
- Aggregate functions
- Function signatures
- Expression typing

Sovrunn owns
- Semantic expression usage
- Cross resource expression context
- Non relational expression extension
- Capability linked expression requirements

Expression Roles

- Predicate
- Projection
- Computation
- Aggregation
- Transformation
- Reference
- Condition

Properties

Input
- Resource
- Attribute
- Value
- Metadata

Output
- Value
- Boolean
- Resource Reference
- Metadata

Characteristics

MUST
- Be deterministic when inputs are deterministic
- Preserve semantic meaning
- Declare required Capability

MAY
- Reference multiple Resources
- Compose with other Expressions
- Use adopted functions
- Use extension functions

Extension

Sovrunn MAY define expression extensions for
- Graph semantics
- Vector semantics
- Object semantics
- Stream semantics
- Policy semantics

Extension MUST
- Preserve SIR semantics
- Declare Capability
- Avoid redefining Substrait semantics

References
- sir.md
- concepts.md
- resources.md
- relationships.md
- capability-model.md
- adopted-standards.md
