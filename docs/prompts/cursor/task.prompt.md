Implement {{FEATURE_ID}} Task {{TASK_ID}} only.

{{MODEL_RECOMMENDATIONS}}

Important spec paths:
- {{REQUIREMENTS_PATH}}
- {{DESIGN_PATH}}
- {{TASKS_PATH}}

Also use:
- AGENTS.md
- docs/engineering/go-coding-guardrails.md
- docs/engineering/go-version-standard.md
- docs/engineering/go-observability-standard.md
- docs/architecture/observability-and-audit-baseline.md
- relevant existing implementation patterns in the repo

Task source:
Read {{TASKS_PATH}} and implement only Task {{TASK_ID}}.

Hard constraints:
- Implement Task {{TASK_ID}} only.
- Do not implement future tasks.
- Do not introduce external dependencies.
- Follow docs/engineering/go-version-standard.md.
- Do not import internal/server from internal/api.
- Do not add unrelated future scope.
- Do not weaken existing tests to pass new code.
- Do not leave TODO({{FEATURE_ID}}) markers.
- Do not commit build artifacts.

Verification command:

docker run --rm -v "$PWD":/src -w /src ${GO_DOCKER_IMAGE:-golang:1.22} sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'

After implementation, report:
- files changed
- verification output
- tests added/updated
- observability added or preserved
- audit behavior added or preserved
- request/operation correlation impact
- fields intentionally not logged
- security considerations
- issues remaining
- scope confirmation

Full tasks.md content for reference:

```markdown
{{TASKS_CONTENT}}
```

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
- request IDs and operation IDs are propagated where applicable,
- structured logs are used where applicable,
- no secrets, credentials, tokens, private keys, or connection strings are logged,
- preserved adapter boundaries.
