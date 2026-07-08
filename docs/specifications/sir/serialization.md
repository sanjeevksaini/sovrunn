# Serialization

Document
- ID: sir-serialization
- Version: 1.0
- Status: Stable

Purpose
- Define SIR serialization rules
- Define canonical representation formats
- Prevent custom serialization

Rules

MUST
- Preserve semantic intent
- Preserve deterministic structure
- Use adopted standards
- Support versioning
- Support validation

MUST NOT
- Define custom binary format
- Redefine adopted serialization formats
- Couple serialization to engine implementation

Canonical Formats

Primary

- Protocol Buffers

Secondary

- JSON

Data Representation

- Apache Arrow

Format Roles

Protocol Buffers
- Runtime transport
- Internal service communication
- Capability manifests
- SIR exchange

JSON
- Debugging
- Documentation
- Human readable inspection
- Test fixtures

Apache Arrow
- Value representation
- Columnar data
- Typed vectors
- IPC data exchange

Serialization Invariants

Every serialized SIR document

MUST
- Include SIR version
- Include schema version
- Use canonical capability identifiers
- Use adopted data types
- Preserve field names
- Preserve semantic meaning

Compatibility

MUST
- Support backward compatible decoding
- Reject unsupported major versions
- Ignore unknown optional fields when safe

MUST NOT
- Silently reinterpret semantic meaning
- Accept structurally invalid SIR

Ownership

Sovrunn owns
- SIR serialization profile
- Field semantics
- Compatibility rules

Adopted standards own
- Protocol Buffers encoding
- JSON encoding
- Apache Arrow representation

References
- sir.md
- versioning.md
- conformance.md
- adopted-standards.md
