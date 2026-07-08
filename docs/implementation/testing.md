# Testing

Document:
  ID: implementation-testing
  Title: Testing
  Parent: implementation
  Owner: SDE Engineering
  Layer: Implementation
  Type: CONTRACT
  Version: 1.1
  Status: Draft

Purpose:
  - Define testing strategy for Sovrunn Data Engine
  - Ensure architecture boundaries are enforced by tests
  - Establish conformance testing for pluggable management planes, plugins, and providers
  - Support safe implementation of runtime, control plane, DMP, and future AI capabilities

Testing Principle:
  Tests must validate both behavior and architecture boundaries.

Test Categories:
  Unit Tests:
    Purpose:
      - Validate package-level behavior.

  Integration Tests:
    Purpose:
      - Validate interactions between modules.
    Examples:
      - Protocol Runtime with mock Protocol Plugin
      - Data Kernel with mock Engine Runtime
      - Management Plane Framework with mock management plane
      - DMP Controller Runtime with mock Datastore Operator Plugin
      - Workflow Service with mock policy

  Conformance Tests:
    Purpose:
      - Validate management plane, plugin, and provider compliance with SDE contracts.
    Required For:
      - Pluggable Management Planes
      - Protocol Plugins
      - Engine Plugins
      - Datastore Operator Plugins
      - Infrastructure Providers
      - Foundation Providers

  Compatibility Tests:
    Purpose:
      - Validate behavior against client libraries, downstream engines, and supported versions.

  End-to-End Tests:
    Purpose:
      - Validate full request or workflow.

  Security Tests:
    Purpose:
      - Validate authorization, policy, isolation, secrets, and audit.

  Performance Tests:
    Purpose:
      - Validate latency, throughput, resource usage, and scaling behavior.

  Failure Tests:
    Purpose:
      - Validate behavior under failure.

Architecture Boundary Tests:
  Purpose:
    - Prevent forbidden imports and illegal direct calls.

  Required Rules:
    - Protocol Plugin must not import Engine Plugin implementation.
    - Engine Plugin must not import Protocol Plugin implementation.
    - SDE Data Plane must not import DMP controllers.
    - SDE Data Plane must not import Management Plane controllers.
    - SDE Runtime must not import DMP controllers.
    - Datastore Operator Plugin must not import Data Plane runtime.
    - Infrastructure Provider must not import Data Plane runtime.
    - AI Control Plane must not be required by Data Plane.
    - DMP must use Management Plane Framework boundaries.

Management Plane Conformance:
  Must Test:
    - management plane manifest validation
    - management plane registration
    - lifecycle start/stop behavior
    - controller runtime integration
    - policy integration
    - workflow integration
    - audit integration
    - observability integration
    - failure classification
    - no tenant data-plane execution

DMP Conformance:
  Must Test:
    - DMP registers as pluggable management plane
    - DatastoreRequest reconciliation
    - DatastoreInstance status updates
    - Tenant Namespace handling
    - DatastoreProfile validation
    - DatastorePolicy validation
    - Datastore Operator Plugin invocation through DMP contracts
    - Infrastructure Provider invocation through DMP contracts
    - dstoreOps workflow integration
    - no direct Data Plane dependency

Protocol Plugin Conformance:
  Must Test:
    - manifest validation
    - supported protocol version declaration
    - request decode
    - protocol-normalized intent creation
    - result mapping
    - error mapping
    - session behavior
    - transaction behavior where applicable
    - unsupported feature handling

Engine Plugin Conformance:
  Must Test:
    - manifest validation
    - execution fragment handling
    - native operation translation
    - result mapping
    - error mapping
    - capability declaration
    - timeout behavior
    - retry classification

Datastore Operator Plugin Conformance:
  Must Test:
    - manifest validation
    - plan generation
    - dry-run behavior
    - apply behavior
    - idempotency
    - rollback or compensation where supported
    - policy integration
    - audit integration
    - failure classification

Infrastructure Provider Conformance:
  Must Test:
    - manifest validation
    - infrastructure plan
    - dry-run behavior
    - apply behavior
    - failure classification
    - no direct tenant data access

Foundation Provider Conformance:
  Must Test:
    - service contract compatibility
    - authentication/authorization behavior where applicable
    - policy decision behavior
    - audit event generation
    - secrets access behavior

AI Testing:
  Current Scope:
    - Placeholder only.

  Future Required Tests:
    - tenant isolation tests
    - action policy tests
    - recommendation safety tests
    - workflow request validation tests
    - audit tests
    - approval gate tests

Recommended Test Commands:
```bash
go test ./...
go test ./internal/managementplane/...
go test ./internal/dmp/...
go test ./internal/runtime/...
go test ./tests/conformance/...
go test ./tests/security/...
```

CI Requirements:
  Required:
    - unit tests
    - integration tests
    - lint
    - forbidden import checks
    - documentation build

  Before Management Plane Admission:
    - management plane manifest validation
    - management plane conformance tests
    - security tests
    - compatibility tests

  Before Plugin Admission:
    - plugin conformance tests
    - manifest validation
    - compatibility tests
    - security tests

  Before Release:
    - e2e tests
    - performance baseline
    - failure tests
    - upgrade tests
    - documentation validation

Invariants:
  - Tests must prove boundaries, not just behavior.
  - Management plane conformance is mandatory before management plane admission.
  - Plugin conformance is mandatory before registry admission.
  - DMP workflows must be tested for idempotency.
  - Security tests must include tenant isolation.
  - AI-generated code or artifacts are untrusted until tests pass.
