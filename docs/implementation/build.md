# Build

Document:
  ID: implementation-build
  Title: Build
  Parent: implementation
  Owner: SDE Engineering
  Layer: Implementation
  Type: CONTRACT
  Version: 1.1
  Status: Draft

Purpose:
  - Define build standards for Sovrunn Data Engine
  - Provide consistent local and CI build commands
  - Support reproducible binaries, containers, management-plane artifacts, and plugin artifacts
  - Ensure docs, specs, and code remain aligned

Build Principle:
  The build must be deterministic, repeatable, and validation-oriented.

Primary Build Tooling:
  - Go toolchain
  - Makefile
  - Docker
  - Helm
  - MkDocs
  - golangci-lint or equivalent
  - go test
  - schema validation tools
  - conformance test runners

Recommended Make Targets:

```makefile
.PHONY: all
all: lint test build

.PHONY: build
build:
	go build ./cmd/...

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: docs
docs:
	mkdocs build --strict

.PHONY: validate
validate:
	go test ./...
	mkdocs build --strict

.PHONY: dev
dev:
	go run ./cmd/sde-control-plane

.PHONY: dev-dmp
dev-dmp:
	go run ./cmd/sde-dmp-controller

.PHONY: clean
clean:
	rm -rf bin dist site
```

Build Outputs:
  Binaries:
    - bin/sde-control-plane
    - bin/sde-data-plane
    - bin/sde-management-plane-controller
    - bin/sde-dmp-controller
    - bin/sde-plugin-runner
    - bin/sde-cli

  Containers:
    - sovrunn/sde-control-plane
    - sovrunn/sde-data-plane
    - sovrunn/sde-management-plane-controller
    - sovrunn/sde-dmp-controller
    - sovrunn/sde-plugin-runner

  Management Plane Artifacts:
    - management plane manifests
    - datastore-management-plane bundle
    - management plane controller runtime configuration

  Plugin Artifacts:
    - protocol plugin bundles
    - engine plugin bundles
    - datastore operator plugin bundles
    - infrastructure provider bundles
    - foundation provider bundles

Build Profiles:
  local:
    - developer machine
    - mock providers
    - local configuration
    - debug logging

  dev:
    - shared development environment
    - non-production services

  test:
    - CI and integration tests
    - ephemeral dependencies
    - strict validation

  staging:
    - production-like environment
    - release candidates

  production:
    - signed images
    - locked configuration
    - strict policy

Validation During Build:
  - Go compilation
  - Unit tests
  - Lint
  - Forbidden import checks
  - Management plane manifest validation
  - Plugin manifest schema validation
  - API schema validation
  - MkDocs build
  - RFC link validation where available

Management Plane Build Rules:
  - Management plane manifests must be packaged with management plane artifact.
  - Management plane artifact must declare type, version, compatibility, required Foundation Services, and supported controller runtime.
  - DMP must be packaged as a management plane artifact when using pluggable deployment.
  - Management plane artifact must pass conformance tests before admission.

Plugin Build Rules:
  - Plugin manifests must be packaged with plugin artifact.
  - Plugin artifact must declare type.
  - Plugin artifact must declare contract version.
  - Plugin artifact must declare compatibility range.
  - Plugin artifact must pass conformance tests before registry admission.

Documentation Build:
  Required:
    - mkdocs build --strict

Recommended Local Build Flow:
```bash
make lint
make test
make docs
make build
```

Recommended CI Build Flow:
```bash
make lint
make test
make validate
make build
docker build .
```

Invariants:
  - Build must fail on compile errors.
  - Build must fail on forbidden imports.
  - Build must fail on broken docs navigation.
  - Management plane artifacts must not be admitted without manifest validation.
  - Plugin artifacts must not be admitted without manifest validation.
  - Production artifacts must be reproducible and versioned.
