# Operations

Document
- ID: sir-operations
- Version: 1.0
- Status: Stable

Purpose
- Define semantic operations
- Define operation properties
- Define operation invariants

Rules

MUST
- Represent semantic intent
- Operate on Resources
- Use Expressions
- Remain protocol independent
- Remain engine independent

MUST NOT
- Represent protocol commands
- Represent engine commands
- Represent execution plans
- Represent runtime lifecycle
- Redefine Substrait operations

Definition

An Operation is semantic intent applied to one or more Resources using zero or more Expressions

Ownership

Sovrunn owns
- Semantic operation model
- Operation categories
- Operation invariants
- Capability requirements

Substrait owns
- Relational operators
- Relational transformations
- Relational execution semantics

Core Operation Categories

Access
- Read
- Discover
- Observe

Mutation
- Create
- Update
- Delete
- Replace
- Merge

Transformation
- Transform
- Derive
- Materialize

Invocation
- Invoke
- Execute

Movement
- Copy
- Move
- Replicate

Operation Properties

Target
- Resource

Input
- Resource
- Expression
- Constraint
- Metadata

Output
- Resource
- Metadata

Characteristics

MUST
- Declare target Resource
- Declare semantic intent
- Declare required Capability
- Preserve Resource semantics

MAY
- Read multiple Resources
- Produce multiple Resources
- Compose with other Operations
- Be idempotent
- Produce side effects

MUST NOT
- Declare execution location
- Declare execution order
- Declare engine selection
- Declare physical plan

Composition

Operations MAY compose when
- Resource semantics are compatible
- Expression semantics are compatible
- Capability requirements are satisfied

Invariants

Every Operation
MUST
- Be valid against target Resource
- Be valid against referenced Expressions
- Be valid against declared Constraints
- Preserve semantic meaning

References
- sir.md
- concepts.md
- resources.md
- relationships.md
- expressions.md
- constraints.md
- capability-model.md
- adopted-standards.md
