# Coding Standards

Document:
  ID: implementation-coding
  Title: Coding Standards
  Parent: implementation
  Owner: SDE Engineering
  Layer: Implementation
  Type: CONTRACT
  Version: 1.1
  Status: Draft

Purpose:
  - Define coding standards for Sovrunn Data Engine
  - Preserve architecture boundaries in code
  - Reflect Datastore Management Plane as a pluggable management plane inside SDE Control Plane
  - Support reliable implementation by humans and AI coding agents
  - Establish rules for errors, context, logging, validation, and testing

Language:
  Primary:
    - Go

General Rules:
  - Keep packages small and cohesive.
  - Prefer explicit interfaces at architecture boundaries.
  - Avoid global mutable state.
  - Pass context.Context through request lifecycles.
  - Return typed errors or wrapped errors with stable categories.
  - Validate inputs at API boundaries.
  - Use structured logging.
  - Emit metrics for important runtime, workflow, and management-plane events.
  - Keep command entrypoints thin.
  - Do not mix transport, business logic, and persistence in one package.

Architecture Boundary Rules:
  SDE Control Plane:
    - May host pluggable management planes.
    - Must use Management Plane Framework to register and govern management planes.
    - Must not hard-code datastore lifecycle logic that belongs to DMP.

  Management Plane Framework:
    - Must define generic management plane contracts.
    - Must not contain datastore-specific lifecycle logic.
    - Must integrate with Foundation Services for policy, workflow, audit, and observability.
    - Must not execute tenant data-plane requests.

  Datastore Management Plane:
    - Must be implemented as a pluggable management plane.
    - Must use Management Plane Framework governance.
    - Must use Foundation Services for policy, workflow, audit, identity, authorization, secrets, and observability.
    - Must invoke Datastore Operator Plugins through approved DMP contracts.
    - Must invoke Infrastructure Providers through approved DMP contracts.

  DMP Controller Runtime:
    - Hosts and reconciles DMP resources and workflows.
    - Is not the whole DMP.
    - Must not execute tenant data-plane requests.

  SDE Data Plane:
    - Must not mutate Control Plane authoritative state.
    - Must not invoke DMP lifecycle controllers.
    - Must not invoke Management Plane controllers.
    - Must not invoke Infrastructure Providers.

  Protocol Plugin:
    - Must not invoke Engine Plugin.
    - Must not produce Execution Plan.
    - Must not access Downstream Datastore.

  Engine Plugin:
    - Must not parse client protocol.
    - Must not manage datastore lifecycle.
    - Must not invoke Infrastructure Provider.

  Datastore Operator Plugin:
    - Must not execute tenant data-plane requests.
    - Must not replace Engine Plugin.
    - Must not bypass DMP workflows, policy, or audit.

  AI Control Plane:
    - Must remain optional and pluggable.
    - Must not be required by Data Plane execution.
    - Must not bypass Control Plane services.

Context Rules:
  Required Context Values:
    - tenant id
    - request id
    - trace id
    - actor id where available
    - namespace where applicable
    - management plane id where applicable
    - authorization context where applicable

Error Rules:
  Error Model:
    - Runtime and API errors must map to SDE Error Model.
    - Native downstream errors must be mapped by Engine Plugin.
    - Protocol-visible errors must be mapped by Protocol Plugin.
    - Management plane errors must be mapped to stable control-plane error categories.

  Error Categories:
    - validation
    - authorization
    - policy
    - capability
    - plugin
    - managementplane
    - dmp
    - runtime
    - datastore
    - infrastructure
    - timeout
    - conflict
    - internal

Logging Rules:
  - Use structured logs.
  - Include trace id and request id.
  - Include management plane id for management-plane workflows.
  - Include tenant id only where allowed by policy.
  - Redact secrets.
  - Do not log raw credentials.
  - Do not log tenant data payloads by default.
  - Log workflow decisions and policy denials.

Metrics Rules:
  Emit metrics for:
    - request count
    - request latency
    - error count
    - plugin invocation count
    - plugin latency
    - downstream latency
    - management plane registration count
    - management plane reconciliation count
    - DMP workflow count
    - workflow duration
    - policy denial count
    - datastore operation count

Tracing Rules:
  - Use OpenTelemetry-compatible tracing.
  - Span boundaries should follow architecture boundaries.
  - Trace Protocol Execution, Planning Execution, Kernel Execution, Engine Execution, and Result/Error Propagation.
  - Trace Management Plane Framework operations.
  - Trace DMP workflows and plugin calls.

Validation Rules:
  - Validate external input at API boundaries.
  - Validate manifests before registry admission.
  - Validate Management Plane manifests before management-plane admission.
  - Validate Execution Plan before Kernel Execution.
  - Validate DatastoreRequest before DMP workflow.
  - Validate policy before operational action.
  - Validate AI-generated artifacts before use.

Security Coding Rules:
  - Enforce tenant isolation at boundaries.
  - Enforce management plane authorization.
  - Use authorization checks before state changes.
  - Use policy checks before workflows.
  - Do not expose internal errors to tenants.
  - Redact sensitive details.
  - Do not bypass Secrets Service.

AI-Assisted Coding Rules:
  AI coding agents may:
    - Generate package skeletons.
    - Generate tests.
    - Generate interface implementations.
    - Generate docs.
    - Suggest refactors.

  AI coding agents must not:
    - Change architecture boundaries without RFC.
    - Introduce forbidden imports.
    - Treat DMP as a fixed non-pluggable subsystem.
    - Add direct cloud/datastore access where DMP or plugins are required.
    - Add destructive behavior without explicit approval.
    - Treat generated code as accepted without tests.

Code Review Checklist:
  - Does package match architecture owner?
  - Are imports allowed?
  - Is DMP modeled as a pluggable management plane?
  - Is sde-dmp-controller treated as DMP Controller Runtime only?
  - Are tenant boundaries preserved?
  - Are errors mapped correctly?
  - Are logs safe?
  - Are metrics/traces emitted?
  - Are tests included?
  - Are public contracts versioned?
  - Are plugin boundaries respected?
  - Are management plane boundaries respected?
  - Are docs updated if behavior changed?
