# CI/CD

Document:
  ID: implementation-ci-cd
  Title: CI/CD
  Parent: implementation
  Owner: SDE Engineering
  Layer: Implementation
  Type: FLOW
  Version: 1.1
  Status: Draft

Purpose:
  - Define CI/CD expectations for Sovrunn Data Engine
  - Ensure code, docs, specifications, and RFCs are validated together
  - Establish quality gates for runtime, control plane, pluggable management planes, plugins, providers, and DMP workflows
  - Prepare for safe release management

CI Principle:
  CI must validate architecture boundaries, tests, documentation, management plane conformance, plugin conformance, and artifact integrity before merge.

Pipeline Stages:
  Stage 1:
    Name:
      - Checkout and setup
    Actions:
      - checkout code
      - set up Go
      - set up Python
      - install docs dependencies
      - restore cache

  Stage 2:
    Name:
      - Static validation
    Actions:
      - gofmt check
      - go vet
      - lint
      - forbidden import checks
      - management plane manifest validation
      - plugin manifest schema validation
      - API schema validation

  Stage 3:
    Name:
      - Unit tests
    Actions:
      - go test ./...
      - package-level tests
      - interface tests

  Stage 4:
    Name:
      - Architecture boundary tests
    Actions:
      - verify forbidden imports
      - verify management plane boundary rules
      - verify plugin boundary rules
      - verify Data Plane does not import DMP controllers
      - verify Data Plane does not import Management Plane controllers
      - verify AI Control Plane is not required by Data Plane

  Stage 5:
    Name:
      - Documentation validation
    Actions:
      - mkdocs build --strict
      - validate navigation
      - validate Mermaid rendering support
      - detect empty markdown files

  Stage 6:
    Name:
      - Integration tests
    Actions:
      - runtime integration tests
      - control plane integration tests
      - management plane framework integration tests
      - DMP integration tests
      - mock plugin tests

  Stage 7:
    Name:
      - Conformance tests
    Actions:
      - Management Plane conformance
      - Protocol Plugin conformance
      - Engine Plugin conformance
      - Datastore Operator Plugin conformance
      - Infrastructure Provider conformance
      - Foundation Provider conformance

  Stage 8:
    Name:
      - Security tests
    Actions:
      - tenant isolation tests
      - management plane authorization tests
      - authorization tests
      - policy denial tests
      - secret redaction tests
      - audit emission tests
      - dependency scan

  Stage 9:
    Name:
      - Build artifacts
    Actions:
      - build binaries
      - build containers
      - build management plane bundles
      - build plugin bundles
      - build docs site

  Stage 10:
    Name:
      - Publish artifacts
    Actions:
      - publish release candidates
      - publish container images
      - publish docs preview
      - publish management plane bundles to staging registry
      - publish plugin bundles to staging registry

Quality Gates:
  Pull Request Gate:
    Required:
      - lint
      - unit tests
      - docs build
      - forbidden import checks

  Main Branch Gate:
    Required:
      - integration tests
      - management plane conformance tests for changed management planes
      - conformance tests for changed plugins
      - security tests for security-sensitive changes
      - documentation validation

  Release Gate:
    Required:
      - full test suite
      - compatibility tests
      - upgrade tests
      - release notes
      - signed artifacts
      - approved RFC/source-of-truth consistency check

Management Plane Admission Gate:
  Required:
    - management plane manifest validation
    - management plane conformance tests
    - control plane integration tests
    - policy tests
    - audit tests
    - registry admission approval

Plugin Admission Gate:
  Required:
    - manifest validation
    - conformance tests
    - compatibility matrix
    - security review where applicable
    - registry admission approval

DMP Workflow Gate:
  Required:
    - DMP registered as pluggable management plane
    - idempotency tests
    - dry-run tests
    - policy tests
    - audit tests
    - rollback or compensation tests

AI Gate:
  Current:
    - AI Control Plane placeholder only.

  Future:
    - AI action policy tests
    - tenant isolation tests
    - recommendation validation tests
    - approval workflow tests
    - audit tests
    - rollback tests

Documentation Gate:
  Required:
    - docs build
    - RFC index updated when RFC added
    - glossary updated when terms added
    - architecture docs updated when RFC accepted
    - no broken nav

Release Artifacts:
  - binaries
  - container images
  - management plane bundles
  - Helm charts
  - CRD or resource schemas where applicable
  - plugin bundles
  - API schemas
  - documentation site
  - release notes
  - SBOM where applicable

Invariants:
  - No merge without docs build.
  - No management plane admission without conformance.
  - No plugin admission without conformance.
  - No DMP workflow release without idempotency tests.
  - No security-sensitive change without security review.
  - No AI action authority without policy and audit validation.
