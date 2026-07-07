# Capability Model

Document
- ID: sir-capability-model
- Version: 1.0
- Status: Stable

Purpose
- Define capability schema
- Define capability semantics
- Define capability lifecycle

Rules

MUST
- Be protocol independent
- Be engine independent
- Be deterministic
- Use stable canonical identifiers

MUST NOT
- Enumerate capabilities
- Reference specific engines
- Define implementation behavior

Definition

A Capability is a semantic contract describing functionality that MAY be provided by an implementation.

Ownership

Sovrunn owns
- Capability schema
- Capability semantics
- Capability compatibility

Capability Schema

Identifier

MUST
- Be globally unique
- Be stable
- Use lowercase
- Use dot notation

Examples

- data.read
- query.aggregate
- search.vector

Name

MUST
- Be human readable

Category

MUST
- Reference canonical capability category

Version

MUST
- Follow semantic versioning

Status

Allowed

- Stable
- Experimental
- Deprecated

Capability Relationships

MAY

- Require Capability
- Extend Capability
- Replace Capability
- Conflict With Capability

Compatibility

Compatible

- Semantic behavior preserved

Partial

- Semantic degradation declared

Incompatible

- Semantic behavior cannot be preserved

Lifecycle

Capability

↓

Published

↓

Adopted

↓

Deprecated

↓

Removed

References

- capabilities.md
- sir.md
