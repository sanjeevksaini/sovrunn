# Semantic Intermediate Representation

Document
- ID: sir
- Version: 1.0
- Status: Stable

Purpose
- Define Semantic Intermediate Representation
- Define semantic execution contract
- Define platform semantic boundary

Definition

Semantic Intermediate Representation (SIR) is the canonical semantic representation of user intent within Sovrunn.

SIR separates semantic intent from:

- Client protocols
- Query languages
- Data engines
- Execution strategies
- Physical implementations

Goals

SIR MUST

- Preserve semantic intent
- Be protocol independent
- Be engine independent
- Be deterministic
- Be extensible
- Be versioned

SIR MUST NOT

- Represent protocol syntax
- Represent engine syntax
- Represent physical execution
- Represent engine implementation
- Depend on a specific database

Responsibilities

SIR defines

- Resources
- Operations
- Expressions
- Relationships
- Constraints
- Metadata
- Capability Model

SIR does not define

- Wire protocols
- Runtime architecture
- Execution plans
- Storage formats
- Transport protocols
- Engine implementations

Lifecycle

Application

↓

Protocol

↓

Protocol Runtime

↓

SIR

↓

Planning

↓

Execution Plan

↓

Data Kernel

↓

Engine Runtime

↓

Engine

Ownership

Sovrunn owns

- Semantic model
- Semantic contracts
- Semantic validation
- Capability model

Sovrunn adopts

- Apache Arrow
- Substrait
- Adopted standards
- Adopted algorithms

Design Principles

Semantic First

- Preserve meaning

Declarative

- Describe intent
- Never describe execution

Canonical

- One semantic representation
- Multiple protocol representations
- Multiple engine implementations

Extensible

- Additive evolution
- Versioned specifications
- Backward compatibility

Compliance

Every Protocol Runtime

MUST

- Produce valid SIR

Every Planner

MUST

- Consume SIR

Every Engine Runtime

MUST

- Preserve SIR semantics

Every Engine

MUST

- Execute equivalent semantics

References

- ontology.md
- adopted-standards.md
- adopted-algorithms.md
- adopted-architecture-patterns.md
