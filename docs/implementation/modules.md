# Modules

Document:
  ID: implementation-modules
  Title: Modules
  Parent: implementation
  Owner: SDE Engineering
  Layer: Implementation
  Type: ARCHITECTURE
  Version: 1.1
  Status: Draft

Purpose:
  - Define implementation modules for Sovrunn Data Engine
  - Map architecture components to code modules
  - Represent Datastore Management Plane as a pluggable management plane
  - Define module responsibilities, boundaries, and forbidden dependencies
  - Support safe implementation by human engineers and AI coding agents

Module Design Principle:
  A module should have one architectural responsibility.

  A module must not hide cross-plane dependencies.

  SDE Control Plane may host pluggable management planes.

  Datastore Management Plane is implemented as the first pluggable management plane.

Core Modules:
  platform:
    Path:
      - internal/platform
    Responsibility:
      - Shared platform primitives
      - Identifiers
      - Common errors
      - Configuration loading
      - Lifecycle helpers
      - Common validation utilities
    Must Not:
      - Implement Control Plane business logic
      - Implement Data Plane request execution
      - Access Downstream Datastores

  spec:
    Path:
      - internal/spec
    Responsibility:
      - Versioning rules
      - Serialization rules
      - Capability models
      - Protocol specification models
      - Engine specification models
      - Management Plane specification models
      - Manifest schemas
    Must Not:
      - Execute requests
      - Manage runtime state
      - Manage datastore lifecycle

  runtime:
    Path:
      - internal/runtime
    Responsibility:
      - Protocol Runtime
      - SIR Runtime
      - Planning
      - Data Kernel
      - Engine Runtime
      - Plugin Runtime
      - Session Runtime
      - Transaction Runtime
      - Result Model
      - Error Model
    Must Not:
      - Own tenant lifecycle
      - Own datastore lifecycle
      - Mutate Control Plane authoritative state
      - Call Datastore Operator Plugins
      - Depend on DMP controllers

  dataplane:
    Path:
      - internal/dataplane
    Responsibility:
      - Data Plane service process
      - Request ingress
      - Runtime composition
      - Protocol endpoint serving
      - Data Plane response lifecycle
    Must Not:
      - Manage datastore lifecycle
      - Invoke Infrastructure Providers
      - Modify Control Plane registries
      - Import DMP controllers

  controlplane:
    Path:
      - internal/controlplane
    Responsibility:
      - Core Control Plane
      - Runtime Registry
      - Plugin Registry
      - Engine Registry
      - Capability Governance
      - Deployment Governance
      - Management Plane Framework integration
      - Optional AI Control Plane reservation
    Must Not:
      - Execute tenant data requests
      - Replace Data Plane runtime
      - Hard-code datastore lifecycle behavior that belongs to DMP

  controlplane_managementplane:
    Path:
      - internal/controlplane/managementplane
    Responsibility:
      - Register pluggable management planes
      - Govern management plane lifecycle
      - Expose management plane APIs through Control Plane
      - Validate management plane manifests
      - Coordinate with Foundation Services
    Must Not:
      - Contain datastore-specific lifecycle logic
      - Execute tenant data requests

  managementplane:
    Path:
      - internal/managementplane
    Responsibility:
      - Shared framework for pluggable management planes
      - Management plane contracts
      - Management plane manifest model
      - Management plane controller runtime abstractions
      - Workflow, policy, audit, and observability integration
      - Management plane admission and lifecycle
    Must Not:
      - Contain datastore-specific DMP logic
      - Execute tenant data-plane requests
      - Replace Foundation Services

  dmp:
    Path:
      - internal/dmp
    Responsibility:
      - Datastore Management Plane implementation
      - Tenant Namespace Manager
      - Datastore Request Controller
      - Datastore Instance Controller
      - Datastore Profile and Policy handling
      - dstoreOps workflows
      - Datastore Operator Plugin integration
      - Infrastructure Provider integration
    Must Not:
      - Execute tenant data-plane requests
      - Parse client protocol
      - Replace Engine Runtime
      - Bypass Management Plane Framework governance

  foundation:
    Path:
      - internal/foundation
    Responsibility:
      - Foundation Service interfaces and implementations
      - Identity
      - Authorization
      - Tenant Management
      - Configuration
      - Policy
      - Secrets
      - Audit
      - Workflow
      - Eventing
      - Observability
      - Registry
      - Plugin Framework
    Must Not:
      - Own domain-specific DMP lifecycle logic
      - Own Data Plane execution logic

  plugins:
    Path:
      - internal/plugins
      - plugins
    Responsibility:
      - Plugin interfaces
      - Plugin adapters
      - Plugin implementations
      - Plugin manifests
      - Plugin conformance tests
    Types:
      - Protocol Plugin
      - Engine Plugin
      - Pluggable Management Plane
      - Datastore Operator Plugin
      - Infrastructure Provider
      - Foundation Provider

  security:
    Path:
      - internal/security
    Responsibility:
      - Authn helpers
      - Authz helpers
      - Tenant isolation utilities
      - Policy helpers
      - Redaction utilities
      - Security validation

  observability:
    Path:
      - internal/observability
    Responsibility:
      - Logging
      - Metrics
      - Tracing
      - Health checks
      - Telemetry emitters

  api:
    Path:
      - internal/api
      - api
    Responsibility:
      - API handlers
      - API schemas
      - API transport bindings
      - OpenAPI and Protobuf contracts

Module Boundary Matrix:
| From | May Depend On | Must Not Depend On |
|---|---|---|
| cmd | internal modules | plugin implementation internals directly |
| dataplane | runtime, spec, observability, security | dmp controllers, management plane controllers, infrastructure providers |
| runtime | spec, platform, plugin interfaces | controlplane state, dmp controllers, management plane controllers |
| controlplane | foundation, spec, platform, managementplane | dataplane request handlers |
| managementplane | foundation, spec, platform | dataplane runtime, datastore-specific DMP implementation |
| dmp | managementplane, foundation, datastore operator interfaces, infrastructure provider interfaces | dataplane runtime execution |
| protocol plugin | protocol runtime interfaces, spec | engine plugin, datastore operator plugin |
| engine plugin | engine runtime interfaces, spec | protocol plugin, DMP controllers |
| management plane plugin | management plane interfaces, foundation | dataplane runtime |
| datastore operator plugin | DMP operator interfaces | dataplane runtime |
| infrastructure provider | DMP infrastructure interfaces | dataplane runtime |
| foundation provider | foundation service interfaces | dataplane runtime |

AI Module:
  Path:
    - internal/controlplane/ai
  Status:
    - Reserved
    - Optional
    - Deferred
  Initial Rule:
    - May contain only placeholder interfaces or future extension points.
    - Must not become a dependency of SDE Data Plane.
    - Must integrate through Control Plane services.

Module Invariants:
  - Modules must align with architecture domains.
  - DMP is a pluggable management plane inside SDE Control Plane.
  - DMP Controller Runtime is not the whole DMP.
  - Shared types must not become hidden global state.
  - Runtime must remain deterministic and testable.
  - Plugins must be replaceable through manifests and registries.
  - Pluggable management planes must be admitted through Management Plane Framework.
  - DMP workflows must be policy-governed and audited.
  - AI integration must remain optional until accepted through RFC.
