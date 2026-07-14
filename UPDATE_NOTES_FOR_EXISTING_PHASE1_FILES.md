---
doc_type: update_notes
title: Update Notes for Existing Phase 1 Files
status: draft
phase: 1
ai_load_priority: reference
ai_summary: Shows the small edits recommended for existing Phase 1 implementation-pack files after adding best-practice documents.
---

# Update Notes for Existing Phase 1 Files

## 1. Update `docs/prompts/AI_IMPLEMENTATION_PROMPT_PHASE1.md`

In `Load First`, add:

```text
docs/engineering/ai-controlled-development.md
docs/engineering/context-engineering-standard.md
docs/architecture/controller-reconciliation-model.md
docs/architecture/observability-and-audit-baseline.md
docs/architecture/gitops-desired-state-model.md
```

Add this rule:

```text
Do not implement a feature unless the feature file, resource spec, API contract, and acceptance criteria are loaded.
```

Add this output requirement:

```text
Confirm which documents were used as context.
Confirm which non-goals were intentionally not implemented.
```

## 2. Update `docs/features/FEATURE_SEQUENCE.md`

Add this section:

```markdown
## Engineering Practices Applied

Phase 1 implementation must follow:

- controlled AI-assisted development
- context engineering
- spec-first implementation
- controller/reconciliation resource model
- GitOps-friendly resource shape
- observability and audit baseline
- test-gated code generation
- human review for architecture changes
```

## 3. Update `docs/resource-specs/RESOURCE_MODEL_PHASE1.md`

Add this section:

```markdown
## Desired-State Rule

All resources must be designed as desired-state resources.

- `metadata` represents identity and classification.
- `spec` represents desired state.
- `status` represents observed state.
- `operation` records lifecycle activity.

User input must not directly set `status`.
```

Add this section:

```markdown
## GitOps Compatibility

Every Phase 1 resource must be serializable as YAML/JSON and safe to store in Git.

Rules:

- names must be stable
- references must be explicit
- spec must be declarative
- status must be system-owned
- secrets must be referenced, not embedded
```

## 4. Update `docs/api/API_CONTRACT_PHASE1.md`

Add this section:

```markdown
## Standard Request Headers

Optional request headers:

| Header | Purpose |
|---|---|
| `X-Sovrunn-Request-ID` | Client-provided request correlation ID |
| `X-Sovrunn-Actor` | Temporary development actor identity |
| `X-Sovrunn-Organization` | Optional organization context |
| `X-Sovrunn-Tenant` | Optional tenant context |

The server must generate a request ID if one is not provided.
```

Add this section:

```markdown
## Standard Response Headers

| Header | Purpose |
|---|---|
| `X-Sovrunn-Request-ID` | Request correlation ID |
| `X-Sovrunn-Operation-ID` | Operation ID for mutating requests after FEATURE-0005 |
```

Add this standard error-code list:

```text
VALIDATION_FAILED
RESOURCE_NOT_FOUND
RESOURCE_ALREADY_EXISTS
DELETE_BLOCKED
REFERENCE_INVALID
POLICY_DENIED
UNAUTHORIZED
FORBIDDEN
INTERNAL_ERROR
```

## 5. Optional: Update `README.md`

Add:

```markdown
## Best-Practice Additions

Phase 1 also follows:

- AI Controlled Development Standard
- Context Engineering Standard
- Controller and Reconciliation Model
- Observability and Audit Baseline
- GitOps Desired-State Model
```

## 6. No Feature Sequence Change

Do not change the feature order.

The additions strengthen implementation discipline but do not expand Phase 1 scope.
