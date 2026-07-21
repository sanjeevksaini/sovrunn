---
doc_type: standard
title: Phase 2 Reuse Assessment Standard
status: draft
phase: 2
ai_load_priority: always
ai_summary: Mandatory reuse assessment format for every Sovrunn feature.
---

# Phase 2 Reuse Assessment Standard

Every feature must include this section in its architecture contract, requirements, design, and review.

```markdown
## Reuse Assessment

### Existing mature solutions
- ...

### Decision
Reuse / Wrap / Extend / Build

### Reason
- ...

### Sovrunn-owned responsibility
- ...

### Non-goals
- ...

### Adapter boundary required?
Yes / No

### Future replacement risk
Low / Medium / High
```

## Decision Rules

| Decision | Use When |
|---|---|
| Reuse | Mature OSS or standard already solves a non-differentiating capability. |
| Wrap | Mature capability exists but must be governed, audited, or exposed through Sovrunn abstractions. |
| Extend | Existing tool is close but needs Sovrunn-specific behavior. |
| Build | Capability is core Sovrunn differentiation or no mature fit exists. |

## Sovrunn Differentiation

Build these:

```text
govern
decide
place
orchestrate
explain
audit
evidence
AI-assist
```

Reuse or wrap mature components for policy engines, IAM, secrets, workflow engines, Kubernetes controllers, PostgreSQL operators, observability, backup, traffic, AI runtimes, and persistence.
