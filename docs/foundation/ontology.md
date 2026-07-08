# Ontology

Document:
  ID: ontology
  Version: 1.1
  Status: Stable

Purpose:
  - Canonical concepts
  - Canonical semantics
  - Architecture-wide concept map

Rules:
  - One concept one definition
  - One concept one owner
  - Reference never redefine
  - Use canonical glossary terms

Schema:
  Purpose:
  Owner:

Ontology:

  Foundation:

    Philosophy:
      Purpose: Engineering philosophy
      Owner: Foundation

    Vision:
      Purpose: Product vision
      Owner: Foundation

    Constitution:
      Purpose: Immutable platform laws
      Owner: Foundation

    Glossary:
      Purpose: Canonical vocabulary
      Owner: Foundation

    Ontology:
      Purpose: Canonical concept model
      Owner: Foundation

    Ownership:
      Purpose: Canonical ownership
      Owner: Foundation

    ADS:
      Purpose: AI documentation specification
      Owner: Foundation

    Style:
      Purpose: Documentation writing rules
      Owner: Foundation

  Platform:

    Sovrunn:
      Purpose: Parent entity and product family
      Owner: Foundation

    Sovrunn Data Engine:
      Purpose: Sovereign semantic execution platform across heterogeneous downstream datastores
      Owner: Architecture

    SDE:
      Purpose: Canonical abbreviation for Sovrunn Data Engine
      Owner: Glossary

  Specifications:

    Protocol:
      Purpose: Client interface specification
      Owner: Protocol Specification

    SIR:
      Purpose: Canonical semantic representation
      Owner: SIR Specification

    Capability:
      Purpose: Semantic capability specification
      Owner: Capability Specification

  Semantic:

    Resource:
      Purpose: Addressable semantic entity
      Owner: SIR Specification

    Operation:
      Purpose: Semantic intent applied to resources
      Owner: SIR Specification

    Expression:
      Purpose: Declarative semantic logic
      Owner: SIR Specification

    Relationship:
      Purpose: Semantic association between resources
      Owner: SIR Specification

    Constraint:
      Purpose: Semantic validation rule
      Owner: SIR Specification

    Metadata:
      Purpose: Non-behavior-changing semantic description
      Owner: SIR Specification

    Request:
      Purpose: Client or runtime request
      Owner: Protocol Specification

    Response:
      Purpose: Client or runtime response
      Owner: Protocol Specification

  Runtime Architecture:

    Runtime:
      Purpose: Execution environment that transforms semantic intent into downstream execution
      Owner: Architecture

    Protocol Runtime:
      Purpose: Runtime component that manages client protocol request lifecycle
      Owner: Runtime Architecture

    SIR Runtime:
      Purpose: Runtime component that manages live SIR instances
      Owner: Runtime Architecture

    Planning:
      Purpose: Runtime component that transforms SIR into Execution Plan
      Owner: Runtime Architecture

    Execution Plan:
      Purpose: Immutable runtime execution contract
      Owner: Runtime Architecture

    Execution Context:
      Purpose: Immutable execution-scoped runtime context
      Owner: Runtime Architecture

    Data Kernel:
      Purpose: Semantic execution orchestrator
      Owner: Runtime Architecture

    Engine Runtime:
      Purpose: Runtime component that delegates execution through Engine Plugins
      Owner: Runtime Architecture

    Capability Registry:
      Purpose: Runtime-facing capability lookup registry
      Owner: Runtime Architecture

    Result Model:
      Purpose: Canonical runtime result representation
      Owner: Runtime Architecture

    Error Model:
      Purpose: Canonical runtime failure representation
      Owner: Runtime Architecture

  SDE Control Plane:

    SDE Control Plane:
      Purpose: Management authority and management-plane host for SDE
      Owner: SDE Control Plane Architecture

    Control Plane Foundation:
      Purpose: Shared reusable management layer of SDE Control Plane
      Owner: SDE Control Plane Architecture

    Foundation Service:
      Purpose: Stable SDE Control Plane service contract
      Owner: Control Plane Foundation

    Foundation Provider:
      Purpose: Pluggable implementation of a Foundation Service
      Owner: Control Plane Foundation

    Management Plane:
      Purpose: Domain-specific control-plane component hosted under SDE Control Plane authority
      Owner: SDE Control Plane Architecture

    Management Plane Host:
      Purpose: Host for built-in and pluggable Management Planes
      Owner: SDE Control Plane Architecture

    Core Control Plane:
      Purpose: Built-in Management Plane for SDE runtime governance
      Owner: Core Control Plane

    Datastore Management Plane:
      Purpose: Optional pluggable Management Plane for downstream datastore lifecycle
      Owner: Datastore Management Plane

  SDE Data Plane:

    SDE Data Plane:
      Purpose: Runtime request execution plane of SDE
      Owner: SDE Data Plane Architecture

    Protocol Plugin:
      Purpose: Client protocol integration into SDE Data Plane
      Owner: Plugin Architecture

    Engine Plugin:
      Purpose: Downstream execution integration from SDE Data Plane to Downstream Engine
      Owner: Plugin Architecture

  Datastore Management:

    dstoreOps:
      Purpose: Product capability for managed downstream datastore operations
      Owner: Datastore Management Plane

    Datastore:
      Purpose: Downstream system that stores, serves, processes, or exposes data
      Owner: External

    Downstream Datastore:
      Purpose: External datastore used by SDE for native storage, processing, or data access
      Owner: External

    Downstream Engine:
      Purpose: Downstream datastore execution engine used through an Engine Plugin
      Owner: External

    Datastore SDE Data Plane:
      Purpose: Native execution and data access plane of a Downstream Datastore
      Owner: External

    Datastore Operator Plugin:
      Purpose: Downstream datastore lifecycle integration for Datastore Management Plane
      Owner: Datastore Management Plane

    Infrastructure Provider:
      Purpose: Infrastructure environment integration used by Datastore Management Plane
      Owner: Datastore Management Plane

  External:

    Application:
      Purpose: External client or system consuming SDE
      Owner: External

    Standard:
      Purpose: Adopted external specification
      Owner: Adopted Standards
