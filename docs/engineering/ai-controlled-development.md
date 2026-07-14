---
doc_type: engineering_standard
title: AI Controlled Development Standard
status: draft
phase: 1
ai_load_priority: always
ai_summary: Defines how AI tools may generate Sovrunn code while preserving architectural control, test discipline, and feature boundaries.
---

# AI Controlled Development Standard

## 1. Purpose

Sovrunn is built with heavy AI assistance, but AI must remain under architectural control.

The principle is:

```text
AI may generate implementation.
AI must not own architecture.
```

## 2. Required Inputs Before AI Coding

An AI coding task must include:

```text
feature ID
feature spec
resource spec
API contract
acceptance criteria
non-goals
relevant engineering standards
```

For Phase 1, AI must load:

```text
docs/foundation/constitution.md
docs/decisions/DECISION_INDEX.md
docs/glossary.md
docs/ai/AI_CONTEXT_GUIDE.md
docs/features/FEATURE_SEQUENCE.md
docs/resource-specs/RESOURCE_MODEL_PHASE1.md
docs/api/API_CONTRACT_PHASE1.md
docs/engineering/go-style.md
docs/engineering/package-layout.md
docs/engineering/testing-standard.md
docs/engineering/ai-controlled-development.md
```

Then it must load only the current feature file.

## 3. What AI May Do

AI may generate:

```text
Go structs
validation functions
in-memory registry methods
REST handlers
tests
curl examples
README updates
demo scripts
small refactors
test-gap analysis
```

## 4. What AI Must Not Do

AI must not:

```text
invent new architecture
create new resource kinds without approval
rename accepted canonical terms
skip feature specs
skip tests
change feature sequence
implement future features early
introduce Kubernetes CRDs before approval
introduce a persistent database before approval
introduce plugin execution before approval
introduce AI agent execution before approval
introduce UI or portal code before approval
introduce billing implementation before approval
bypass validation rules
remove audit/operation hooks
change security-sensitive behavior silently
```

## 5. Feature Boundary Rule

AI must implement one feature at a time.

While implementing `FEATURE-0001 Organization Resource and Registry`, AI may implement:

```text
Organization resource
Organization registry
Organization validation
Organization API
minimal health/readiness/version endpoints
Organization tests
```

It must not implement:

```text
Tenant
Project
ServiceInstance
Plugin execution
database storage
Kubernetes CRDs
GitOps controller
```

## 6. Architecture Drift Control

Every AI-generated change must be checked against:

```text
constitution.md
DECISION_INDEX.md
glossary.md
FEATURE_SEQUENCE.md
current feature file
RESOURCE_MODEL_PHASE1.md
API_CONTRACT_PHASE1.md
```

If AI proposes a change that conflicts with those files, reject the change unless a founder-approved architecture update is made.

## 7. Test-Gated Rule

No AI-generated feature is accepted without tests.

Minimum tests:

```text
resource validation tests
registry create/get/list/update/delete tests
duplicate name tests
missing reference tests where applicable
delete-blocked tests where applicable
API handler happy-path tests
API handler error-path tests
```

## 8. Dependency Control

AI must not add new dependencies unless there is a clear reason.

Before adding a dependency, AI must explain:

```text
why it is needed
what alternative exists in the standard library
whether it affects security or supply-chain risk
whether it is required for Phase 1
```

## 9. Review Checklist

Before accepting AI-generated code, verify:

```text
feature ID is clear
scope matches the feature file
canonical terms are used
resource shape follows metadata/spec/status
validation rules are implemented
errors are deterministic
tests pass
non-goals are respected
API contract is followed
logs do not leak secrets
new dependencies are justified
```

## 10. Commit Rule

Each feature should be committed separately.

Recommended commit format:

```text
feat(FEATURE-0001): implement Organization resource and registry
test(FEATURE-0001): add Organization registry and API tests
docs(FEATURE-0001): add Organization demo commands
```

## 11. Final Principle

AI accelerates Sovrunn development. Architecture remains founder-controlled, spec-first, test-gated, and review-driven.
