Feature:
{{FEATURE_ID}} {{FEATURE_TITLE}}

Generate tasks.md only.

Do not implement code.
Do not modify source files.
Do not modify requirements.md.
Do not modify design.md unless you find a blocking contradiction.

Inputs:
- {{REQUIREMENTS_PATH}}
- {{DESIGN_PATH}}
- AGENTS.md
- docs/engineering/go-coding-guardrails.md
- docs/engineering/go-version-standard.md
- docs/engineering/go-observability-standard.md

Tool-output safety constraints:
- Use fs_write only for chunks of 50 lines or fewer.
- For files longer than 50 lines, create the file with fs_write using the first chunk, then use fs_append in chunks of 50 lines or fewer.
- Do not write the entire tasks.md in one fs_write call.
- Do not use one very large str_replace edit.
- Split content into logical sections.
- Write one section at a time.
- After writing, read the file back and verify it is complete.

Task generation rules:
- Create implementation tasks only.
- Each task must be small enough for Cursor to implement safely.
- Prefer one focused code area per task.
- Add tests close to the code they verify.
- Do not combine model, registry, handler, server wiring, and integration tests into one task.
- Each task must include objective, files, notes, tests, acceptance criteria, and commit message.
- Final task must include full Docker verification, guardrails, artifact cleanup, and clean git status.

Standard Docker verification command:

docker run --rm -v "$PWD":/src -w /src ${GO_DOCKER_IMAGE:-golang:1.22} sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'

Final Docker verification command:

docker run --rm -v "$PWD":/src -w /src ${GO_DOCKER_IMAGE:-golang:1.22} sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./... && go test -race ./... && go build ./cmd/sovrunn-api'

Final guardrails:
- rm -f sovrunn-api
- rm -rf bin
- no TODO({{FEATURE_ID}}) under internal or cmd
- no internal/api import of internal/server
- git status clean

Generate tasks.md only.

## Phase 2 Reuse and Drift Gates

Every generated feature must include:

```markdown
## Reuse Assessment

### Existing mature solutions
- ...

### Decision
Reuse / Wrap / Extend / Build

### Sovrunn-owned responsibility
- ...

### Adapter boundary required?
Yes / No

### Non-goals
- ...
```

Architecture drift checks:

- no provider-specific hardcoding in core,
- no Kubernetes-only assumptions in core,
- no PostgreSQL lifecycle logic in core placement engine,
- no custom policy engine embedded in handlers,
- no raw secret storage,
- no customer-facing IaaS leakage,
- explainable decision object,
- defined audit behavior,
- defined observability behavior,
- request/operation correlation where applicable,
- no secret or credential logging,
- preserved adapter boundaries.
