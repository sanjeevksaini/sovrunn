# Implementation Skeleton

Document:
  ID: implementation-skeleton
  Title: Implementation Skeleton
  Parent: implementation
  Owner: SDE Engineering
  Layer: Implementation
  Type: CONTRACT
  Version: 1.1
  Status: Draft

Purpose:
  - Define the first buildable Sovrunn Data Engine code skeleton
  - Provide a safe starting point for incremental implementation
  - Preserve architecture boundaries from day one
  - Include the Management Plane Framework and DMP-as-pluggable-plane model from day one
  - Enable AI coding agents to generate code in the correct modules

Skeleton Goal:
  The initial skeleton should compile, start processes, expose health endpoints, load configuration, initialize registries, initialize the Management Plane Framework, register DMP as the first pluggable management plane, and wire placeholder interfaces without implementing full database functionality.

MVP Skeleton Scope:
  - Go module initialization
  - Thin command entrypoints
  - Configuration loading
  - Structured logging
  - Health endpoints
  - Runtime package skeleton
  - Control Plane package skeleton
  - Management Plane Framework skeleton
  - Datastore Management Plane skeleton
  - DMP Controller Runtime skeleton
  - Plugin interface skeletons
  - Management plane manifest skeleton
  - Specification model skeletons
  - Test scaffolding
  - Local Docker Compose or Kubernetes manifests

Out of Scope:
  - Full PostgreSQL protocol compatibility
  - Full engine execution
  - Production DMP lifecycle management
  - Real cloud infrastructure provisioning
  - Generic multi-management-plane runtime beyond placeholder
  - AI Control Plane implementation
  - Production security hardening

Initial Commands:
  cmd/sde-control-plane:
    Responsibility:
      - Start Control Plane API server
      - Initialize Foundation Services
      - Initialize Core Control Plane registries
      - Initialize Management Plane Framework
      - Register DMP as available pluggable management plane
      - Expose health endpoint

  cmd/sde-data-plane:
    Responsibility:
      - Start Data Plane gateway
      - Initialize runtime components
      - Load Protocol Plugin registry metadata
      - Expose health endpoint
      - Accept no-op or mock requests initially

  cmd/sde-management-plane-controller:
    Responsibility:
      - Optional future generic controller runtime for pluggable management planes
      - Host approved management planes
      - Execute generic management-plane reconciliation loops

  cmd/sde-dmp-controller:
    Responsibility:
      - Start the DMP Controller Runtime
      - Host and reconcile Datastore Management Plane resources
      - Load mock DatastoreRequest resources
      - Execute no-op DMP workflows initially
      - Expose health endpoint

  cmd/sde-cli:
    Responsibility:
      - Provide local developer commands
      - Validate manifests
      - Validate management plane manifests
      - Inspect configuration
      - Run local conformance tests

Recommended First Package Creation Order:
  1. internal/platform/errors
  2. internal/platform/identifiers
  3. internal/platform/config
  4. internal/observability/logs
  5. internal/observability/health
  6. internal/spec/versioning
  7. internal/spec/serialization
  8. internal/spec/capability
  9. internal/spec/managementplane
  10. internal/managementplane/framework
  11. internal/managementplane/registry
  12. internal/managementplane/controller
  13. internal/runtime/result
  14. internal/runtime/error
  15. internal/runtime/session
  16. internal/runtime/transaction
  17. internal/runtime/protocol
  18. internal/runtime/sir
  19. internal/runtime/planning
  20. internal/runtime/kernel
  21. internal/runtime/engine
  22. internal/plugins/protocol
  23. internal/plugins/engine
  24. internal/foundation/policy
  25. internal/foundation/audit
  26. internal/foundation/workflow
  27. internal/controlplane/registry
  28. internal/controlplane/managementplane
  29. internal/dmp/request
  30. internal/dmp/controllers

Initial Interface Examples:

Management Plane Interface:
```go
package managementplane

type Plane interface {
    Name() string
    Version() string
    Kind() string
    Manifest() Manifest
    Start(ctx Context) error
    Stop(ctx Context) error
}
```

Management Plane Controller Interface:
```go
package managementplane

type Controller interface {
    Name() string
    Reconcile(ctx Context, request ReconcileRequest) (ReconcileResult, error)
}
```

Protocol Plugin Interface:
```go
package protocol

type Plugin interface {
    Name() string
    Version() string
    Capabilities() []string
    Decode(ctx Context, input []byte) (NormalizedIntent, error)
    EncodeResult(ctx Context, result Result) ([]byte, error)
    EncodeError(ctx Context, err Error) ([]byte, error)
}
```

Engine Plugin Interface:
```go
package engine

type Plugin interface {
    Name() string
    Version() string
    Capabilities() []string
    Execute(ctx Context, fragment ExecutionFragment) (Result, error)
}
```

Datastore Operator Plugin Interface:
```go
package datastoreoperator

type Plugin interface {
    Name() string
    Version() string
    SupportedEngines() []string
    Plan(ctx Context, request DatastoreRequest) (OperationPlan, error)
    Apply(ctx Context, plan OperationPlan) (OperationResult, error)
}
```

Infrastructure Provider Interface:
```go
package infrastructureprovider

type Provider interface {
    Name() string
    Version() string
    Plan(ctx Context, request InfrastructureRequest) (InfrastructurePlan, error)
    Apply(ctx Context, plan InfrastructurePlan) (InfrastructureResult, error)
}
```

Skeleton Acceptance Criteria:
  - `go test ./...` passes.
  - All command entrypoints compile.
  - Health endpoints respond locally.
  - Management Plane Framework interfaces compile.
  - DMP registers as first pluggable management plane.
  - Plugin interfaces compile.
  - No forbidden imports exist.
  - Configuration loads from local config.
  - Logs include trace/request identifiers.
  - MkDocs documentation builds.
  - No placeholder package violates architecture boundaries.

Initial Local Run:
```bash
make dev
make test
make lint
make docs
```

Invariants:
  - Skeleton must compile before feature implementation.
  - Skeleton must include tests from the beginning.
  - Skeleton must not introduce direct coupling between Protocol Plugin and Engine Plugin.
  - Skeleton must not allow Data Plane to invoke DMP lifecycle operations.
  - Skeleton must model DMP as a pluggable management plane.
  - Skeleton must not implement AI Control Plane beyond placeholder extension point.
