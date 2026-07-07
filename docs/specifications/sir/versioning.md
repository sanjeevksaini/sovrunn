# Versioning

Document
- ID: sir-versioning
- Version: 1.0
- Status: Stable

Purpose
- Define SIR versioning rules
- Define compatibility rules
- Define evolution rules

Rules

MUST
- Preserve semantic intent
- Prefer additive change
- Version public contracts
- Reject incompatible major versions

MUST NOT
- Break compatibility without RFC
- Rename published identifiers
- Reuse removed identifiers
- Silently change semantics

Version Model

SIR Version
- Version of SIR specification

Schema Version
- Version of serialized structure

Capability Version
- Version of capability contract

Extension Version
- Version of extension contract

Change Types

Patch
- Editorial fixes
- Clarifications
- No semantic change

Minor
- Additive fields
- Additive capabilities
- Additive metadata
- Backward compatible

Major
- Breaking semantic change
- Removed required field
- Changed semantic meaning
- Incompatible contract

Compatibility

Backward Compatible
- Older consumers can safely ignore new optional fields

Forward Compatible
- Newer consumers can read older versions

Incompatible
- Semantic meaning cannot be preserved

Deprecation

MUST
- Mark deprecated identifiers
- Provide replacement when available
- Preserve deprecated identifiers during compatibility window

MUST NOT
- Remove deprecated identifiers without major version

Identifier Stability

Published Identifiers

MUST
- Remain stable
- Remain globally unique
- Never be repurposed

Applies To
- Capability identifiers
- Resource kinds
- Operation categories
- Metadata keys
- Extension namespaces

Negotiation

Producer

MUST
- Declare SIR version
- Declare schema version
- Declare required capabilities

Consumer

MUST
- Validate supported version
- Reject unsupported major version
- Preserve semantic meaning

Ownership

Sovrunn owns
- SIR version policy
- Compatibility rules
- Deprecation rules

References
- sir.md
- serialization.md
- capabilities.md
- conformance.md
