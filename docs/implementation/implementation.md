# Implementation

Document:
  ID: implementation
  Title: Implementation
  Parent: docs
  Owner: SDE Engineering
  Layer: Implementation
  Type: MAP
  Version: 1.1
  Status: Draft

Purpose:
  - Define the implementation documentation layer for Sovrunn Data Engine
  - Bridge architecture, specifications, RFCs, and executable code
  - Provide a stable map for repository structure, modules, naming, build, testing, and local development
  - Represent Datastore Management Plane as a pluggable management plane inside SDE Control Plane
  - Help human engineers and AI coding agents navigate implementation work without violating architecture boundaries

Definition:
  The Implementation layer describes how accepted SDE architecture and specifications are organized into repositories, modules, packages, build pipelines, tests, local development workflows, and controller runtimes.

  It does not redefine architecture.

  It translates architecture into code-ready structure.

Reading Order:
  1. repository.md
  2. modules.md
  3. skeleton.md
  4. naming.md
  5. coding.md
  6. build.md
  7. testing.md
  8. local-development.md
  9. ci-cd.md

Implementation Principles:
  - Architecture boundaries must be visible in module boundaries.
  - SDE Control Plane may host pluggable management planes.
  - Datastore Management Plane is the first pluggable management plane.
  - DMP Controller Runtime is the executable runtime that hosts and reconciles DMP.
  - SDE Control Plane code must not depend on SDE Data Plane runtime internals.
  - SDE Data Plane code must not mutate Control Plane authoritative state.
  - Protocol Plugin code must not invoke Engine Plugin code directly.
  - Engine Plugin code must not parse client protocol.
  - Datastore Operator Plugin code must not execute tenant data-plane requests.
  - Infrastructure Provider code must not be confused with Foundation Provider code.
  - AI Control Plane code is optional and pluggable; detailed implementation is deferred.
  - Tenant AI Agent code, when added later, must integrate through approved Control Plane APIs and workflows.

Implementation Domains:
  - cmd
  - internal/controlplane
  - internal/controlplane/managementplane
  - internal/managementplane
  - internal/dataplane
  - internal/runtime
  - internal/spec
  - internal/plugins
  - internal/dmp
  - internal/foundation
  - internal/security
  - internal/observability
  - internal/api
  - pkg/sdk
  - plugins
  - plugins/management-plane
  - deployments
  - tests

Architecture Mapping:
  SDE Control Plane:
    Maps To:
      - internal/controlplane
      - internal/controlplane/managementplane
      - internal/managementplane
      - internal/foundation
      - internal/security
      - internal/api/controlplane

  Pluggable Management Plane Framework:
    Maps To:
      - internal/managementplane
      - internal/controlplane/managementplane
      - plugins/management-plane

  Datastore Management Plane:
    Maps To:
      - internal/dmp
      - plugins/management-plane/datastore-management-plane
      - cmd/sde-dmp-controller

  DMP Controller Runtime:
    Maps To:
      - cmd/sde-dmp-controller
      - internal/managementplane/controller
      - internal/dmp/controllers

  SDE Data Plane:
    Maps To:
      - internal/dataplane
      - internal/runtime
      - internal/api/dataplane

  SDE Runtime:
    Maps To:
      - internal/runtime

  Protocol Plugins:
    Maps To:
      - plugins/protocol
      - internal/plugins/protocol

  Engine Plugins:
    Maps To:
      - plugins/engine
      - internal/plugins/engine

  Datastore Operator Plugins:
    Maps To:
      - plugins/datastore-operator
      - internal/dmp/operator

  Infrastructure Providers:
    Maps To:
      - plugins/infrastructure-provider
      - internal/dmp/infrastructure

  Foundation Providers:
    Maps To:
      - plugins/foundation-provider
      - internal/foundation/provider

  dstoreOps:
    Maps To:
      - internal/dmp/dstoreops
      - internal/dmp/workflow
      - internal/dmp/controllers

  AI Control Plane:
    Maps To:
      - internal/controlplane/ai
    Status:
      - Reserved
      - Optional
      - Deferred

Invariants:
  - Implementation docs must not contradict architecture docs.
  - Implementation modules must preserve plane boundaries.
  - DMP is a pluggable management plane, not merely a controller.
  - sde-dmp-controller is a DMP Controller Runtime, not the entire DMP.
  - Public contracts must be versioned.
  - Plugin contracts must be manifest-driven.
  - Management plane contracts must be manifest-driven.
  - Runtime behavior must be testable through conformance suites.
  - Control Plane state changes must go through approved services.
  - Tenant isolation must be enforced in APIs, services, workflows, management planes, and tests.
