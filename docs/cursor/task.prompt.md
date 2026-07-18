Implement {{FEATURE_ID}} Task {{TASK_ID}} only.

{{MODEL_RECOMMENDATIONS}}

Important spec paths:
- {{REQUIREMENTS_PATH}}
- {{DESIGN_PATH}}
- {{TASKS_PATH}}

Also use:
- AGENTS.md
- docs/engineering/go-coding-guardrails.md
- relevant existing implementation patterns in the repo

Task source:
Read {{TASKS_PATH}} and implement only Task {{TASK_ID}}.

Hard constraints:
- Implement Task {{TASK_ID}} only.
- Do not implement future tasks.
- Do not introduce external dependencies.
- Keep Go 1.21 compatible.
- Do not import internal/server from internal/api.
- Do not add unrelated future scope.
- Do not weaken existing tests to pass new code.
- Do not leave TODO({{FEATURE_ID}}) markers.
- Do not commit build artifacts.

Verification command:

docker run --rm -v "$PWD":/src -w /src golang:1.21 sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'

After implementation, report files changed, verification output, issues remaining, and scope confirmation.

Full tasks.md content for reference:

```markdown
{{TASKS_CONTENT}}
```
