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



## FEATURE-0013 task guardrails

If `{{FEATURE_ID}}` is `FEATURE-0013`, these additional rules are mandatory:

- Implement only the requested `{{TASK_ID}}` from `{{TASKS_PATH}}`.
- Do not modify architecture.md, requirements.md, design.md, or tasks.md.
- Do not use any superseded or patched design draft as source input.
- Only `CONTRACT_NOW` task content may produce source changes; `INVARIANT_FOR_LATER` and `DEFERRED` content is documentation/traceability only.
- Preserve `metadata.scopeRef` as the sole scope authority; do not introduce `AuditScope`, a second scope field, a top-level scope alias, or owner/subject-derived scope.
- Preserve `DecisionRequest` as calling-domain-owned; do not create a shared/canonical/common DecisionRequest type, schema, registry entry, resource, or persistence model.
- Do not create PendingDecision, DeferredDecision, DecisionPending, or any FEATURE-0013-specific pending-response envelope.
- Reuse only the existing FEATURE-0012 baseline manifest and approvals mechanism; do not create a parallel schema-baseline manifest, schema-diff directory, or schema-approval directory.
- Treat `SecurityExceptionRef` as structural evidence only; do not create a security-exception approval service, approval workflow, authority model, issuance, revocation, or runtime enforcement mechanism.
- Matrix E work may stage architecture-owned values and blank human fields only; do not calculate, infer, downgrade, upgrade, or reinterpret residual risk.
- Do not implement erasure, key destruction, signing, cryptographic verification, AI runtime, production runtime service, adapter interface, or persistence mechanism unless the requested task explicitly authorizes a `CONTRACT_NOW` structural artifact.

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

Every FEATURE-0011-and-later feature must include a reuse assessment that
conforms to the canonical standard:

`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

Do not duplicate or redefine the assessment field schema in this document.
Populate the feature-level reuse summary and capability-level assessments
using the canonical fields and controlled vocabularies.

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
