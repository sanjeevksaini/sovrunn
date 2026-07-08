# Versioning Specification

Document:
  ID: versioning-specification
  Title: Versioning Specification
  Parent: specifications
  Owner: SDE Specification Layer
  Layer: Specification
  Type: Contract
  Version: 1.0
  Status: Draft

Purpose:
  - Define versioning rules for Sovrunn Data Engine specifications
  - Define compatibility rules across SIR, capabilities, protocols, engines, plugins, and runtime contracts
  - Define how SDE components evolve without breaking approved integrations

Definition:
  Versioning Specification defines the canonical version model used by SDE architecture, specifications, runtime contracts, Control Plane metadata, Data Plane plugins, and downstream datastore integrations.

Scope:
  In Scope:
    - Specification versioning
    - Contract versioning
    - Capability versioning
    - Protocol versioning
    - Engine Plugin versioning
    - SIR versioning
    - Compatibility rules
    - Deprecation rules
    - Breaking-change rules

  Out of Scope:
    - Product release planning
    - Git branching strategy
    - Package manager implementation
    - Customer release notes
    - Runtime upgrade orchestration

Version Format:
  SDE specifications MUST use semantic versioning style:

    MAJOR.MINOR.PATCH

  MAJOR:
    Meaning: Breaking contract change.

  MINOR:
    Meaning: Backward-compatible capability or behavior addition.

  PATCH:
    Meaning: Backward-compatible clarification, correction, or non-behavioral update.

Versioned Objects:
  - SIR
  - Capability Specification
  - Protocol Specification
  - Engine Specification
  - Serialization Specification
  - Protocol Plugin Manifest
  - Engine Plugin Manifest
  - Capability Manifest
  - Execution Plan schema
  - Execution Context schema
  - Result Model schema
  - Error Model schema
  - Control Plane registry metadata

Compatibility Rules:
  - Consumers MUST declare supported versions.
  - Producers MUST declare emitted versions.
  - Runtime MUST reject unsupported MAJOR versions unless explicit compatibility adapter exists.
  - Runtime MAY accept older MINOR versions when backward compatibility is guaranteed.
  - Runtime MAY accept PATCH differences when behavior is unchanged.
  - Compatibility decisions MUST be deterministic.
  - Compatibility decisions MUST be observable through safe telemetry.

Breaking Change Rules:
  A change is breaking when it:
    - Removes a required field
    - Renames a required field
    - Changes semantics of an existing field
    - Changes default behavior visible to clients or plugins
    - Removes a capability
    - Weakens required validation
    - Changes error classification incompatibly
    - Changes wire-visible protocol behavior incompatibly

Non-Breaking Change Rules:
  A change is non-breaking when it:
    - Adds optional field with safe default
    - Adds new capability without changing existing capability behavior
    - Adds new protocol feature behind explicit negotiation
    - Clarifies documentation without changing behavior
    - Adds additional safe error detail
    - Adds new enum value only when consumers are required to ignore unknown values safely

Deprecation Rules:
  - Deprecated versions MUST remain documented.
  - Deprecated versions MUST identify replacement version where available.
  - Deprecated versions MUST identify removal target where known.
  - Runtime MUST emit safe deprecation telemetry when deprecated contract is used.
  - Deprecation MUST NOT silently change behavior.

Negotiation Rules:
  - Protocol Plugins MUST negotiate protocol versions where protocol supports negotiation.
  - Engine Plugins MUST declare supported engine contract versions.
  - Capability Manifests MUST declare supported capability versions.
  - Planning MUST validate version compatibility before Execution Plan emission.
  - Engine Runtime MUST validate Engine Plugin compatibility before invocation.

Registry Rules:
  - Registry entries MUST include version.
  - Registry entries MUST include compatibility metadata where applicable.
  - Registry entries MUST identify lifecycle state.
  - Registry entries MUST identify deprecation state.
  - Registry entries MUST not overwrite incompatible metadata without explicit version change.

Lifecycle States:
  - Draft
  - Experimental
  - Stable
  - Deprecated
  - Removed

Invariants:
  - Every executable contract is versioned.
  - Compatibility is explicit.
  - Breaking changes require MAJOR version change.
  - Deprecated does not mean removed.
  - Runtime must fail closed on unsupported contract version.
  - Versioning must preserve tenant safety and execution determinism.

Relationships:
  Parent:
    - ../specifications
  Depends On:
    - ../sir/versioning.md
    - ../capability/capability.md
    - ../protocol/protocol.md
    - ../engine/engine.md
    - ../serialization/serialization.md
  Used By:
    - SIR Specification
    - Capability Specification
    - Protocol Specification
    - Engine Specification
    - Plugin Manifests
    - Runtime registries
    - SDE Control Plane registries

References:
  - ../sir/versioning.md
  - ../capability/capability.md
  - ../protocol/protocol.md
  - ../engine/engine.md
  - ../serialization/serialization.md
