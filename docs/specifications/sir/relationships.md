# Relationships

Document
- ID: sir-relationships
- Version: 1.0
- Status: Stable

Purpose
- Define semantic relationships
- Define resource associations
- Define dependency semantics

Rules

MUST
- Relate resources
- Preserve semantic meaning
- Be protocol independent
- Be engine independent

MUST NOT
- Represent physical implementation
- Represent relational joins
- Represent execution plans

Definition

A Relationship defines a semantic association between two or more Resources.

Relationships describe platform meaning rather than storage or execution.

Properties

Source

- Resource

Target

- Resource

Kind

- Semantic relationship type

Direction

- Directed
- Bidirectional

Cardinality

- One To One
- One To Many
- Many To One
- Many To Many

Metadata

- Optional
- Extensible

Core Relationship Kinds

Hierarchy

- Contains
- Parent Of
- Child Of

Dependency

- Depends On
- Required By

Ownership

- Owns
- Managed By

Reference

- References
- Referenced By

Composition

- Composed Of
- Part Of

Association

- Associated With

Capability

- Supports
- Supported By

Lifecycle

- Creates
- Consumes
- Produces

Characteristics

Every Relationship

MUST

- Connect valid resources
- Preserve semantic consistency

MAY

- Carry metadata
- Be versioned
- Be constrained

Relationship Invariants

- Source MUST exist
- Target MUST exist
- Kind MUST be defined
- Cardinality MUST be valid

Ownership

Owner
- Sovrunn

References

- concepts.md
- resources.md
- constraints.md
