# Repository Structure

Document:
  ID: implementation-repository
  Title: Repository Structure
  Parent: implementation
  Owner: SDE Engineering
  Layer: Implementation
  Type: ARCHITECTURE
  Version: 1.1
  Status: Draft

Purpose:
  - Define the recommended repository layout for Sovrunn Data Engine
  - Make architecture boundaries visible in code organization
  - Represent Datastore Management Plane as a pluggable management plane inside SDE Control Plane
  - Support modular implementation, plugin extensibility, testing, deployment, and AI-assisted development
  - Avoid coupling between SDE Control Plane, SDE Data Plane, SDE Runtime, Management Planes, plugins, and providers

Repository Principle:
  Repository layout must reflect SDE architecture.

  Code organization should make illegal dependencies obvious.

  SDE Control Plane may host pluggable management planes.

  Datastore Management Plane is the first pluggable management plane.

  DMP Controller Runtime is the executable controller runtime that hosts and reconciles Datastore Management Plane resources and workflows.

---

# Recommended Repository Layout

```text
sovrunn/
  README.md
  go.mod
  go.sum
  Makefile
  mkdocs.yml

  docs/
    architecture/
    specifications/
    foundation/
    implementation/
    rfc/

  cmd/
    sde-control-plane/
    sde-data-plane/
    sde-management-plane-controller/
    sde-dmp-controller/
    sde-plugin-runner/
    sde-cli/

  internal/
    api/
      controlplane/
      dataplane/
      managementplane/
      dmp/
      admin/

    controlplane/
      core/
      registry/
      capability/
      deployment/
      managementplane/
      ai/

    dataplane/
      gateway/
      request/
      response/
      errors/

    runtime/
      protocol/
      sir/
      planning/
      kernel/
      engine/
      plugin/
      session/
      transaction/
      result/
      error/

    managementplane/
      framework/
      registry/
      lifecycle/
      workflow/
      policy/
      admission/
      controller/
      api/

    dmp/
      namespace/
      request/
      instance/
      profile/
      policy/
      credential/
      workflow/
      controllers/
      dstoreops/
      operator/
      infrastructure/

    foundation/
      identity/
      authorization/
      tenant/
      configuration/
      policy/
      secrets/
      audit/
      workflow/
      eventing/
      observability/
      registry/
      plugin/
      provider/

    plugins/
      protocol/
      engine/
      managementplane/
      datastoreoperator/
      infrastructureprovider/
      foundationprovider/

    spec/
      versioning/
      serialization/
      capability/
      protocol/
      engine/
      managementplane/
      manifest/

    security/
      authn/
      authz/
      policy/
      isolation/
      redaction/

    observability/
      logs/
      metrics/
      traces/
      health/

    platform/
      config/
      lifecycle/
      errors/
      identifiers/

  pkg/
    sdk/
      go/
    client/
    types/

  plugins/
    protocol/
      postgresql/
      mysql/
      mongodb/
      redis/
      rest/
      grpc/
      native/

    engine/
      postgresql/
      mysql/
      mongodb/
      redis/
      cassandra/
      opensearch/
      neo4j/
      milvus/
      s3/
      iceberg/
      delta-lake/
      parquet/

    management-plane/
      datastore-management-plane/

    datastore-operator/
      postgresql/
      mysql/
      mongodb/
      redis/
      cassandra/
      opensearch/

    infrastructure-provider/
      kubernetes/
      aws/
      azure/
      gcp/
      vmware/
      baremetal/

    foundation-provider/
      identity/
      authorization/
      policy/
      secrets/
      workflow/
      audit/
      eventing/
      observability/

  api/
    openapi/
    protobuf/
    schemas/

  deployments/
    docker/
    kubernetes/
    helm/
    local/
    terraform/

  configs/
    local/
    dev/
    test/
    examples/

  tests/
    unit/
    integration/
    conformance/
    e2e/
    compatibility/
    security/
    performance/
    fixtures/

  tools/
    codegen/
    lint/
    docs/
    validation/

  scripts/
    dev/
    test/
    build/
    release/
```

---

# Top-Level Directory Rules

docs:
  Purpose:
    - Human and AI-readable architecture, specifications, implementation guidance, and RFCs.

  Rule:
    - Docs are the architecture and specification source of truth.
    - Code must not silently contradict docs.
    - Accepted RFCs must update affected source-of-truth docs.

cmd:
  Purpose:
    - Executable entrypoints only.

  Rule:
    - cmd packages should be thin.
    - Business logic belongs in internal modules.
    - cmd may initialize configuration, logging, dependencies, servers, controller runtimes, and graceful shutdown.

internal:
  Purpose:
    - Main private implementation code.

  Rule:
    - Cross-domain imports must respect architecture boundaries.

pkg:
  Purpose:
    - Public SDKs, client libraries, and stable external Go packages.

  Rule:
    - Do not place unstable internal logic here.

plugins:
  Purpose:
    - Plugin implementations and plugin-specific packaging.

  Rule:
    - Plugin code must conform to plugin contracts and manifests.
    - Plugin implementation folders are not the source of plugin authority.
    - Plugin Registry and manifests define plugin admission.

api:
  Purpose:
    - External API contracts, OpenAPI definitions, Protobuf definitions, JSON schemas.

  Rule:
    - API definitions must be versioned.

deployments:
  Purpose:
    - Deployment artifacts for local, Kubernetes, Helm, Docker, and cloud environments.

configs:
  Purpose:
    - Environment-specific configuration examples and local development configuration.

tests:
  Purpose:
    - Cross-module tests, conformance suites, compatibility tests, and fixtures.

tools:
  Purpose:
    - Developer and CI tooling.

scripts:
  Purpose:
    - Repeatable shell or automation scripts.

---

# Executable Entrypoints

cmd/sde-control-plane:
  Meaning:
    - Starts the SDE Control Plane API and core governance services.

  Responsibilities:
    - Start Control Plane APIs
    - Initialize Foundation Services
    - Initialize Core Control Plane registries
    - Initialize Management Plane Framework
    - Register available pluggable management planes
    - Expose health and readiness endpoints

  Must Not:
    - Execute tenant data requests
    - Run long-running datastore lifecycle operations directly
    - Bypass Workflow Service

cmd/sde-data-plane:
  Meaning:
    - Starts the SDE Data Plane request execution service.

  Responsibilities:
    - Start protocol listeners
    - Load runtime configuration
    - Initialize SDE Runtime
    - Resolve approved Protocol Plugins and Engine Plugins
    - Execute tenant requests through runtime flow

  Must Not:
    - Mutate Control Plane authoritative state
    - Manage datastore lifecycle
    - Invoke Datastore Operator Plugins
    - Invoke Infrastructure Providers

cmd/sde-management-plane-controller:
  Meaning:
    - Generic executable controller runtime for pluggable management planes.

  Responsibilities:
    - Host one or more approved pluggable management planes
    - Run management-plane reconciliation loops
    - Coordinate management-plane workflows
    - Enforce policy and audit boundaries through Control Plane services

  Status:
    - Optional future generic runtime.

cmd/sde-dmp-controller:
  Meaning:
    - Specialized DMP Controller Runtime.

  Responsibilities:
    - Host and reconcile Datastore Management Plane resources
    - Watch DatastoreRequest resources
    - Manage DatastoreInstance state
    - Coordinate DMP lifecycle workflows
    - Invoke Datastore Operator Plugins through DMP contracts
    - Invoke Infrastructure Providers through DMP contracts where required
    - Emit audit events, metrics, and workflow status

  Clarification:
    - DMP is the pluggable management plane.
    - sde-dmp-controller is only the executable controller runtime for DMP.
    - sde-dmp-controller is not the entire Datastore Management Plane.

  Must Not:
    - Execute tenant data-plane requests
    - Parse client protocol
    - Invoke Engine Plugins
    - Bypass Policy Service, Workflow Service, or Audit Service

cmd/sde-plugin-runner:
  Meaning:
    - Starts isolated plugin runtime process when plugins run out-of-process.

  Responsibilities:
    - Load approved plugin bundles
    - Enforce plugin execution boundaries
    - Expose plugin invocation interface
    - Report plugin health and metrics

cmd/sde-cli:
  Meaning:
    - Command-line tool for developers, operators, and automation.

  Responsibilities:
    - Validate manifests
    - Inspect registry metadata
    - Submit local test requests
    - Run conformance tests
    - Inspect health
    - Support local development

---

# Internal Domain Rules

internal/controlplane:
  Purpose:
    - Implements SDE Control Plane core.

  Includes:
    - Core Control Plane
    - Runtime Registry
    - Plugin Registry
    - Engine Registry
    - Capability Governance
    - Deployment Governance
    - Management Plane Framework integration
    - Optional AI Control Plane extension point

  Must Not:
    - Execute tenant data requests
    - Replace SDE Data Plane runtime
    - Hard-code datastore lifecycle logic that belongs to DMP

internal/controlplane/managementplane:
  Purpose:
    - Control Plane integration point for pluggable management planes.

  Responsibilities:
    - Register management plane types
    - Validate management plane manifests
    - Govern management plane lifecycle
    - Expose management plane APIs through Control Plane
    - Coordinate with Foundation Services

  Examples:
    - Datastore Management Plane
    - Future Cache Management Plane
    - Future Search Management Plane
    - Future Tenant Integration Management Plane

internal/managementplane:
  Purpose:
    - Shared framework for pluggable management planes.

  Responsibilities:
    - Management plane contract
    - Management plane manifest model
    - Controller runtime abstractions
    - Workflow integration
    - Policy integration
    - Admission and lifecycle handling

  Must Not:
    - Contain datastore-specific logic
    - Execute tenant data-plane requests
    - Replace Foundation Services

internal/dmp:
  Purpose:
    - Implements Datastore Management Plane as a pluggable management plane.

  Responsibilities:
    - Tenant Namespace handling
    - DatastoreRequest handling
    - DatastoreInstance reconciliation
    - DatastoreProfile handling
    - DatastorePolicy handling
    - Datastore credential references
    - dstoreOps workflows
    - Datastore lifecycle controllers
    - Datastore Operator Plugin integration
    - Infrastructure Provider integration

  Must Not:
    - Execute tenant data-plane requests
    - Parse client protocol
    - Replace Engine Plugin
    - Bypass management-plane framework governance
    - Bypass Policy Service, Workflow Service, or Audit Service

internal/dataplane:
  Purpose:
    - Implements SDE Data Plane service process and request lifecycle.

  Must Not:
    - Import DMP controllers
    - Invoke Datastore Operator Plugins
    - Invoke Infrastructure Providers
    - Mutate Control Plane authoritative state

internal/runtime:
  Purpose:
    - Implements reusable SDE Runtime components.

  Includes:
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
    - Own datastore lifecycle
    - Depend on DMP controllers
    - Depend on Control Plane mutable state

internal/foundation:
  Purpose:
    - Implements Foundation Service contracts and provider integration.

  Must Not:
    - Contain domain-specific DMP lifecycle logic
    - Execute data-plane requests

internal/plugins:
  Purpose:
    - Contains internal plugin contracts, adapters, and shared plugin runtime integration.

  Notes:
    - Actual plugin implementations may live under top-level plugins/.
    - Internal interfaces must remain stable and versioned.

internal/spec:
  Purpose:
    - Contains implementable representation of SDE specifications.

  Includes:
    - Versioning
    - Serialization
    - Capability
    - Protocol
    - Engine
    - Management Plane
    - Manifest models

internal/security:
  Purpose:
    - Shared security primitives.

  Includes:
    - Authn helpers
    - Authz helpers
    - Policy helpers
    - Tenant isolation utilities
    - Redaction utilities

internal/observability:
  Purpose:
    - Shared observability implementation.

  Includes:
    - Logs
    - Metrics
    - Traces
    - Health checks

---

# Plugin Repository Rules

plugins/protocol:
  Purpose:
    - Protocol Plugin implementations.

  Must Not:
    - Invoke Engine Plugins directly
    - Access Downstream Datastores
    - Manage datastore lifecycle

plugins/engine:
  Purpose:
    - Engine Plugin implementations.

  Must Not:
    - Parse client protocol
    - Manage datastore lifecycle
    - Invoke Infrastructure Providers

plugins/management-plane:
  Purpose:
    - Pluggable management plane implementations.

  Current:
    - datastore-management-plane

  Future:
    - cache-management-plane
    - search-management-plane
    - vector-management-plane
    - tenant-integration-management-plane

  Notes:
    - A management plane plugin provides a management domain.
    - DMP is the first management plane implementation.

plugins/datastore-operator:
  Purpose:
    - Datastore Operator Plugin implementations used by DMP.

  Must Not:
    - Execute tenant data-plane requests
    - Replace Engine Plugins
    - Bypass DMP workflows

plugins/infrastructure-provider:
  Purpose:
    - Infrastructure Provider implementations used by pluggable management planes, especially DMP.

  Must Not:
    - Execute tenant data-plane requests
    - Access tenant data
    - Replace Datastore Operator Plugins

plugins/foundation-provider:
  Purpose:
    - Foundation Provider implementations for Foundation Services.

  Must Not:
    - Be confused with Infrastructure Providers
    - Execute tenant data-plane requests

---

# Dependency Direction

Allowed:
  - cmd → internal
  - internal/controlplane → internal/foundation
  - internal/controlplane → internal/managementplane
  - internal/managementplane → internal/foundation
  - internal/managementplane → internal/spec
  - internal/dmp → internal/managementplane
  - internal/dmp → internal/foundation
  - internal/dmp → internal/plugins/datastoreoperator interfaces
  - internal/dmp → internal/plugins/infrastructureprovider interfaces
  - internal/dataplane → internal/runtime
  - internal/dataplane → internal/spec
  - internal/runtime → internal/spec
  - internal/runtime → internal/plugins interfaces

Not Allowed:
  - internal/runtime → internal/controlplane mutable state
  - internal/runtime → internal/dmp controllers
  - internal/dataplane → internal/dmp controllers
  - internal/dataplane → internal/managementplane controllers
  - protocol plugin → engine plugin direct import
  - engine plugin → protocol plugin direct import
  - datastore operator plugin → dataplane runtime
  - infrastructure provider → dataplane runtime
  - foundation provider → dataplane runtime
  - AI Control Plane → direct datastore or infrastructure APIs

---

# Management Plane Model

Management Plane Framework:
  Definition:
    Shared framework inside SDE Control Plane that allows management domains to be added as governed, pluggable planes.

  Responsibilities:
    - Management plane registration
    - Management plane manifest validation
    - Management plane lifecycle
    - Management plane controller runtime integration
    - Management plane API exposure
    - Policy, workflow, audit, and observability integration

Datastore Management Plane:
  Definition:
    Pluggable management plane inside SDE Control Plane responsible for tenant-scoped Downstream Datastore lifecycle and operations.

  Relationship:
    - DMP plugs into Management Plane Framework.
    - DMP uses Foundation Services.
    - DMP uses Datastore Operator Plugins.
    - DMP uses Infrastructure Providers where required.
    - DMP powers dstoreOps.

DMP Controller Runtime:
  Definition:
    Executable runtime that hosts and reconciles DMP resources, workflows, and plugin interactions.

  Relationship:
    - DMP Controller Runtime runs DMP.
    - DMP Controller Runtime is not the whole DMP.

dstoreOps:
  Definition:
    Managed datastore operations capability powered by DMP.

  Relationship:
    - dstoreOps workflows run through DMP.
    - dstoreOps does not execute data-plane requests.

---

# Boundary Enforcement

Required:
  - Use Go internal package boundaries.
  - Use interface packages for plugin contracts.
  - Use management plane interfaces for DMP integration.
  - Use linters for forbidden imports.
  - Use conformance tests for plugin boundaries.
  - Use conformance tests for management plane boundaries.
  - Use review gates for cross-domain dependencies.

Recommended Checks:
  - Data Plane must not import internal/dmp.
  - Runtime must not import internal/dmp.
  - Runtime must not import internal/controlplane mutable packages.
  - DMP must not import protocol plugin implementation packages.
  - DMP must not import engine plugin implementation packages.
  - Datastore Operator Plugins must not import Data Plane runtime.
  - Infrastructure Providers must not import Data Plane runtime.
  - AI Control Plane must not be required by Data Plane binaries.

---

# Naming Clarifications

DMP:
  Meaning:
    - Datastore Management Plane.

  Use For:
    - The pluggable management plane.

DMP Controller Runtime:
  Meaning:
    - The executable runtime that hosts and reconciles DMP.

  Use For:
    - cmd/sde-dmp-controller behavior.

sde-dmp-controller:
  Meaning:
    - Binary name for the DMP Controller Runtime.

Management Plane Plugin:
  Avoid:
    - This term is too generic.

  Prefer:
    - Pluggable Management Plane
    - Datastore Management Plane
    - Datastore Operator Plugin
    - Infrastructure Provider
    - Foundation Provider

Provider:
  Avoid:
    - Provider without qualifier.

  Prefer:
    - Infrastructure Provider
    - Foundation Provider

---

# Migration From Previous Repository Model

Previous Wording:
  - DMP appeared as a fixed subsystem directly under SDE Control Plane.
  - sde-dmp-controller could be misread as the whole DMP.

Corrected Wording:
  - DMP is a pluggable management plane inside SDE Control Plane.
  - sde-dmp-controller is the DMP Controller Runtime.
  - Management Plane Framework is the host framework for pluggable management planes.
  - Datastore Operator Plugins are used by DMP.
  - Infrastructure Providers are used by DMP.

Repository Change:
  Add:
    - internal/managementplane/
    - internal/controlplane/managementplane/
    - plugins/management-plane/datastore-management-plane/

  Keep:
    - internal/dmp/
    - cmd/sde-dmp-controller/
    - plugins/datastore-operator/
    - plugins/infrastructure-provider/

---

# Invariants

- SDE Control Plane may host pluggable management planes.
- DMP is a pluggable management plane inside SDE Control Plane.
- DMP is not the same thing as the DMP Controller Runtime.
- sde-dmp-controller is an executable controller runtime for DMP.
- DMP must operate through Management Plane Framework governance.
- DMP must use Foundation Services for policy, workflow, audit, identity, authorization, secrets, and observability.
- DMP must invoke Datastore Operator Plugins only through approved DMP contracts.
- DMP must invoke Infrastructure Providers only through approved DMP contracts.
- SDE Data Plane must not depend on DMP controllers.
- SDE Runtime must not depend on DMP controllers.
- Protocol Plugins must not invoke Engine Plugins directly.
- Engine Plugins must not manage datastore lifecycle.
- Datastore Operator Plugins must not execute tenant data-plane requests.
- Infrastructure Providers must not execute tenant data-plane requests.
- AI Control Plane remains optional and pluggable until later RFCs define scope.
