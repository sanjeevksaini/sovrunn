# Capability Registry

Document
- ID: capability-registry
- Version: 1.0
- Status: Stable

Purpose
- Define Capability Registry
- Define capability discovery
- Define capability resolution
- Define runtime capability services

Definition

Capability Registry is the runtime service that maintains the authoritative registry of capabilities published by Engine Plugins.

The Capability Registry is the single source of truth for runtime capability discovery.

Principles

MUST

- Preserve capability integrity
- Preserve capability consistency
- Preserve canonical identifiers
- Support deterministic capability resolution
- Support concurrent access

MUST NOT

- Invent capabilities
- Modify published capabilities
- Access downstream engines directly
- Depend on engine implementation

Responsibilities

Capability Registry

- Register capabilities
- Index capabilities
- Resolve capabilities
- Publish capability metadata
- Maintain capability versions
- Maintain engine capability mappings

Capability Sources

Engine Plugin

MUST

- Publish Capability Manifest

Capability Registry

MUST

- Validate Capability Manifest
- Register published capabilities
- Reject invalid manifests

Registry Model

Contains

- Registered Engines
- Registered Engine Plugins
- Capability Manifests
- Capability Catalog References
- Capability Versions

Queries

Planning

MAY query

- Engine supports capability
- Capability version
- Capability compatibility
- Engines supporting capability

Engine Runtime

MAY query

- Engine metadata
- Registered plugin

Runtime Lifecycle

Engine Registration

↓

Capability Manifest Validation

↓

Capability Registration

↓

Capability Index Update

↓

Capability Available

Validation

Capability Registry

MUST

- Validate canonical capability identifiers
- Validate capability versions
- Validate manifest structure
- Validate duplicate registrations

MUST NOT

- Accept unknown required identifiers
- Accept invalid capability versions
- Accept malformed manifests

Characteristics

MUST

- Be deterministic
- Be thread safe
- Be horizontally scalable
- Support concurrent queries

Ownership

Capability Registry owns

- Runtime capability state
- Capability indexing
- Capability resolution

Engine Plugin owns

- Capability Manifest
- Capability publication

Capability Specification owns

- Capability schema
- Capability identifiers
- Capability semantics

References

- runtime.md
- planning.md
- plugin-runtime.md
- engine-runtime.md
- specifications/sir/capability-model.md
- specifications/sir/capabilities.md
