# Metadata

Document
- ID: sir-metadata
- Version: 1.0
- Status: Stable

Purpose
- Define semantic metadata
- Define metadata ownership
- Define metadata boundaries

Rules

MUST
- Describe SIR constructs
- Be optional
- Be extensible
- Be version compatible

MUST NOT
- Change semantic meaning
- Represent business data
- Represent runtime state
- Represent engine implementation

Definition

Metadata describes SIR constructs without modifying their semantic behavior.

Ownership

Sovrunn owns
- Metadata model
- Metadata semantics
- Reserved metadata

Metadata Scope

- Resource
- Relationship
- Expression
- Operation
- Constraint

Metadata Categories

Identity
- Identifier
- Name
- Description

Lifecycle
- Created
- Updated
- Version

Classification
- Labels
- Tags
- Category

Documentation
- Summary
- Documentation
- Reference

Extension
- Vendor Extension
- User Extension

Characteristics

MUST
- Be immutable within a SIR instance
- Be serializable
- Be independently interpretable
- Preserve semantic compatibility

MAY
- Be inherited
- Be extended
- Be ignored by implementations

Reserved Metadata

- id
- name
- description
- version
- labels
- annotations

Extension Rules

Extension

MUST
- Use implementation namespace
- Preserve reserved metadata
- Avoid naming conflicts

MUST NOT
- Override reserved metadata
- Redefine semantic constructs

Metadata Invariants

Every Metadata entry

MUST
- Have unique key
- Have deterministic value
- Be compatible with serialization

References

- sir.md
- resources.md
- relationships.md
- expressions.md
- operations.md
- constraints.md
