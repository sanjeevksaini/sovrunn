# Resources

Document
- ID: sir-resources
- Version: 1.0
- Status: Stable

Purpose
- Define semantic resources
- Define resource identity
- Define resource hierarchy

Rules

MUST
- Represent semantic entities
- Be protocol independent
- Be engine independent
- Have stable identity

MUST NOT
- Represent physical storage
- Represent implementation details
- Depend on engine-specific concepts

Definition

A Resource is an addressable semantic entity that participates in one or more semantic operations.

A Resource represents business meaning rather than physical implementation.

Properties

Identity

- Globally unique
- Stable
- Immutable

Name

- Human readable
- Context scoped

Kind

- Defines semantic classification

Attributes

- Extensible
- Versioned

Metadata

- Optional
- Extensible

Hierarchy

Resources MAY contain child resources.

Example

Workspace
    ├── Catalog
    │     ├── Schema
    │     │      ├── Table
    │     │      ├── View
    │     │      └── Function
    │     └── Collection
    └── Object Store

Kinds

Core

- Workspace
- Catalog
- Schema
- Database
- Table
- View
- Collection
- Stream
- Topic
- Queue
- Function
- Procedure
- Model
- Index
- Object Store
- Bucket
- Object

Behavior

A Resource

MUST

- Support identity
- Support metadata
- Support capabilities
- Support authorization

MAY

- Support hierarchy
- Support relationships
- Support versioning
- Support lifecycle

Lifecycle

Create

↓

Discover

↓

Read

↓

Update

↓

Delete

Capabilities

Every Resource exposes one or more capabilities.

Examples

Table
- Read
- Write
- Scan
- Index

Object
- Read
- Write
- Delete

Stream
- Publish
- Subscribe

Ownership

Owner
- Sovrunn

References

- concepts.md
- capability-model.md
- metadata.md
