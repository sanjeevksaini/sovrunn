# Concepts

Document
- ID: sir-concepts
- Version: 1.0
- Status: Stable

Purpose
- Define the conceptual model of SIR
- Define ownership boundaries
- Define adopted semantic concepts

Rules

MUST
- Adopt concepts before defining new concepts
- Reference canonical specifications
- Preserve semantic intent

MUST NOT
- Duplicate adopted concepts
- Redefine adopted semantics

Concept Model

### Semantic Intent

Definition

Semantic intent represents **what** the application wants to accomplish.

Semantic intent MUST remain independent of:

- Protocol
- Query Language
- Engine
- Physical Execution

Owner
- Sovrunn

---

### Value

Definition

A value represents an immutable semantic value.

Owner
- Apache Arrow

Reference
- Apache Arrow Type System

---

### Type

Definition

A type defines the representation and semantics of a value.

Owner
- Apache Arrow

Reference
- Apache Arrow Type System

---

### Expression

Definition

An expression represents declarative computation over values.

Owner
- Substrait

Reference
- Substrait Expression Model

---

### Relation

Definition

A relation represents a logical dataset.

Owner
- Substrait

Reference
- Substrait Relational Model

---

### Plan

Definition

A plan represents a logical description of computation.

Owner
- Substrait

Reference
- Substrait Plan Specification

---

### Capability

Definition

A capability represents a semantic feature supported by an engine.

Owner
- Sovrunn

---

### Resource

Definition

A resource represents an addressable semantic object.

Owner
- Sovrunn

---

### Operation

Definition

An operation represents an action performed on a resource.

Owner
- Sovrunn

---

### Constraint

Definition

A constraint defines semantic validity.

Owner
- Sovrunn

---

### Metadata

Definition

Metadata describes semantic objects.

Owner
- Sovrunn

---

### Relationship

Definition

A relationship defines semantic association between resources.

Owner
- Sovrunn

Concept Ownership

Apache Arrow
- Value
- Type

Substrait
- Relation
- Expression
- Plan

Sovrunn
- Semantic Intent
- Capability
- Resource
- Operation
- Relationship
- Constraint
- Metadata
