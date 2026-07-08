# Local Development

Document:
  ID: implementation-local-development
  Title: Local Development
  Parent: implementation
  Owner: SDE Engineering
  Layer: Implementation
  Type: FLOW
  Version: 1.1
  Status: Draft

Purpose:
  - Define local development workflow for Sovrunn Data Engine
  - Support one-developer and small-team development
  - Make local testing repeatable
  - Support Go development, MkDocs validation, pluggable management plane testing, and future local Kubernetes deployment

Assumptions:
  - Initial development uses Go.
  - Local mode may use mock Foundation Providers, mock Management Planes, mock Datastore Operator Plugins, and mock Infrastructure Providers.
  - Full production dependencies are not required for skeleton development.
  - DMP is developed as the first pluggable management plane.

Required Tools:
  - Go
  - Git
  - Make
  - Docker or compatible runtime
  - Python
  - MkDocs
  - mkdocs-material
  - pymdown-extensions
  - kubectl optional
  - k3s or kind optional
  - Helm optional

Initial Setup:
```bash
git clone <repo>
cd sovrunn

go mod tidy

python3 -m venv .venv
source .venv/bin/activate
python -m pip install --upgrade pip
python -m pip install mkdocs mkdocs-material pymdown-extensions
```

Documentation Validation:
```bash
source .venv/bin/activate
mkdocs serve
```

Open:
```text
http://127.0.0.1:8000
```

Code Validation:
```bash
go test ./...
make lint
make build
```

Recommended Local Config Layout:
```text
configs/
  local/
    control-plane.yaml
    data-plane.yaml
    management-plane-controller.yaml
    dmp-controller.yaml
    management-plane-registry.yaml
    plugin-registry.yaml
    datastore-profiles.yaml
    policies.yaml
```

Local Services:
  Control Plane:
    Command:
      - go run ./cmd/sde-control-plane

  Data Plane:
    Command:
      - go run ./cmd/sde-data-plane

  Generic Management Plane Controller:
    Command:
      - go run ./cmd/sde-management-plane-controller

    Status:
      - Optional future runtime.

  DMP Controller Runtime:
    Command:
      - go run ./cmd/sde-dmp-controller

  CLI:
    Command:
      - go run ./cmd/sde-cli

Local Development Flow:
  1. Update docs or code.
  2. Run unit tests.
  3. Run boundary checks.
  4. Run docs build.
  5. Run local service.
  6. Validate health endpoint.
  7. Commit with related docs and tests.

Local Health Endpoints:
  Control Plane:
    - http://localhost:8080/healthz
    - http://localhost:8080/readyz

  Data Plane:
    - http://localhost:8081/healthz
    - http://localhost:8081/readyz

  DMP Controller Runtime:
    - http://localhost:8082/healthz
    - http://localhost:8082/readyz

Mock Mode:
  Purpose:
    - Allow development before real providers, management planes, and plugins exist.

  Mock Components:
    - mock identity provider
    - mock authorization provider
    - mock policy service
    - mock workflow service
    - mock audit service
    - mock management plane
    - mock DMP
    - mock protocol plugin
    - mock engine plugin
    - mock datastore operator plugin
    - mock infrastructure provider

  Rule:
    - Mock components must implement the same interfaces as real components.

Local Kubernetes Mode:
  Optional:
    - kind
    - k3s
    - minikube

  Use For:
    - management-plane controller testing
    - DMP controller reconciliation testing
    - namespace isolation testing
    - Helm chart validation
    - operator plugin integration testing

Local DMP Mode:
  Initial:
    - DMP registered as pluggable management plane
    - mock DatastoreRequest
    - mock workflow execution
    - mock Datastore Operator Plugin
    - mock Infrastructure Provider

  Later:
    - PostgreSQL Datastore Operator Plugin
    - local Kubernetes namespace
    - local PostgreSQL deployment

AI Local Mode:
  Status:
    - Deferred.

  Rule:
    - No local AI dependency required for core SDE development.
    - AI Control Plane placeholder must not block local startup.

Common Commands:
```bash
find . -name ".DS_Store" -delete
find docs -name "*.md" -type f -size 0
find docs -name "*.md" -type f -size +12k -exec ls -lh {} \;
mkdocs build --strict
go test ./...
```

Invariants:
  - Local mode must not require production secrets.
  - Local mode must not require cloud access for skeleton.
  - Local mode must preserve architecture boundaries.
  - Mock mode must be replaceable with real providers.
  - DMP must run as pluggable management plane.
  - DMP Controller Runtime must not be confused with the whole DMP.
  - AI Control Plane must remain optional.
