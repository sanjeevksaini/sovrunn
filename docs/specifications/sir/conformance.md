# Conformance

Document
- ID: sir-conformance
- Version: 1.0
- Status: Stable

Purpose
- Define SIR conformance requirements
- Define implementation obligations
- Define interoperability requirements

Rules

MUST
- Preserve semantic intent
- Preserve interoperability
- Preserve deterministic behavior
- Conform to adopted specifications

MUST NOT
- Modify SIR semantics
- Introduce incompatible behavior
- Misrepresent implementation capabilities

Conformance Targets

- Protocol Runtime
- Planner
- Engine Runtime
- Engine Plugin
- Client SDK

Conformance Levels

Level 1

Name
- SIR Producer

MUST
- Produce valid SIR
- Produce valid serialization
- Produce supported SIR version

Examples
- PostgreSQL Protocol Runtime
- MySQL Protocol Runtime
- MongoDB Protocol Runtime
- REST API Runtime

Level 2

Name
- SIR Consumer

MUST
- Consume valid SIR
- Validate SIR
- Preserve semantic meaning

Examples
- Planner
- Optimizer

Level 3

Name
- SIR Executor

MUST
- Execute equivalent semantics
- Preserve operation semantics
- Preserve expression semantics
- Preserve constraint semantics

Examples
- PostgreSQL Engine Runtime
- MongoDB Engine Runtime
- Redis Engine Runtime

Level 4

Name
- SIR Provider

MUST
- Publish capability manifest
- Advertise supported capabilities
- Advertise supported versions

Examples
- PostgreSQL Plugin
- MongoDB Plugin
- Redis Plugin

Validation

Every implementation

MUST validate

- SIR Version
- Schema Version
- Resource Model
- Relationship Model
- Expression Model
- Operation Model
- Constraint Model
- Capability Requirements
- Serialization Rules

Interoperability

Conforming implementations

MUST
- Exchange valid SIR
- Preserve semantic meaning
- Preserve canonical identifiers
- Preserve capability identifiers

MUST NOT
- Depend on engine-specific semantics
- Depend on protocol-specific semantics

Compliance

A conforming implementation

MUST

- Implement required SIR constructs
- Preserve semantic compatibility
- Preserve version compatibility
- Preserve capability contracts

Ownership

Sovrunn owns
- Conformance specification
- Compliance requirements
- Interoperability requirements

References

- sir.md
- serialization.md
- versioning.md
- capability-model.md
- capabilities.md
- adopted-standards.md
