---
doc_type: ai_factory
title: AI Feature Factory
status: draft
phase: 0
ai_load_priority: always
ai_summary: Defines the repeatable workflow used by AI and humans to develop Sovrunn features.
---

# AI Feature Factory

## 1. Purpose

The AI Feature Factory is the repeatable model for developing Sovrunn features.

Every feature should move through the same path:

```text
Feature Idea
  -> DEC
  -> ADR if tradeoff exists
  -> RFC
  -> Resource/API spec
  -> Implementation plan
  -> Tests
  -> Code
  -> Docs
  -> Demo
  -> Runbook
  -> Conformance checks
```

## 2. Factory Output Per Feature

Every major feature should produce or update:

1. DEC entry.
2. ADR if architectural tradeoff exists.
3. RFC.
4. Resource model.
5. API contract.
6. Controller/service design.
7. Plugin contract if relevant.
8. Validation rules.
9. Test plan.
10. Implementation tasks.
11. Code.
12. Unit tests.
13. Integration tests where applicable.
14. Documentation.
15. Demo script.
16. Runbook.

## 3. Feature Task Template

Use this structure:

```markdown
# AI Development Task: <Feature Name>

## Phase
Phase <number>: <phase name>

## Objective
<What this feature must achieve>

## In Scope
- ...

## Out of Scope
- ...

## Required Input Files
- ...

## Existing Code Paths
- ...

## Expected Output Files
- ...

## Acceptance Criteria
- ...

## Test Requirements
- ...

## Documentation Updates
- ...

## Constitutional Checks
- ...
```

## 4. Phase 1 Feature Sequence

Implement Phase 1 in this order:

1. FEATURE-0001: Organization Resource and Registry.
2. FEATURE-0002: OrganizationUnit Resource.
3. FEATURE-0003: Tenant Resource.
4. FEATURE-0004: Project Resource.
5. FEATURE-0005: Operation Resource.
6. FEATURE-0006: ServiceClass and ServicePlan.
7. FEATURE-0007: Plugin and Capability Registry.
8. FEATURE-0008: ServiceInstance and ServiceBinding.
9. FEATURE-0009: API server health/readiness.
10. FEATURE-0010: Basic CLI/API demo flow.

## 5. Feature Gate Rules

A feature is not complete until:

- code compiles,
- tests pass,
- failure modes are tested,
- docs are updated,
- acceptance criteria are met,
- relevant DEC/ADR/RFC references are updated,
- no out-of-scope work was added.

## 6. AI Implementation Rules

AI must:

- implement the smallest useful slice,
- avoid future-phase features,
- generate tests with code,
- keep resource names canonical,
- return explicit errors,
- preserve operation/audit model,
- use structured logs,
- avoid leaking secrets,
- update docs.

AI must not:

- build UI unless task says so,
- add billing,
- add marketplace,
- add multi-cluster federation,
- add advanced AI agents,
- add SDE transformation,
- integrate external systems before phase requires them.

## 7. Review Checklist

Before accepting AI output, review:

| Check | Question |
|---|---|
| Scope | Did AI stay inside the feature? |
| Terms | Did AI use glossary terms? |
| Decisions | Did AI follow DECISION_INDEX? |
| Constitution | Did AI comply with constitution.md? |
| Tests | Are validation and failure tests included? |
| Errors | Are failures explicit and actionable? |
| Docs | Are docs updated? |
| Simplicity | Did AI avoid unnecessary complexity? |

## 8. Human Approval

The founder approves:

- architecture,
- DEC changes,
- ADRs,
- RFCs,
- resource names,
- API contracts,
- security model,
- tenant model,
- plugin model,
- SDE hot-path rules.

AI may draft these but must not treat drafts as accepted without approval.

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
- preserved adapter boundaries.
