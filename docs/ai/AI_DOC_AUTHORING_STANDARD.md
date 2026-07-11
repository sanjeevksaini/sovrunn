---
doc_type: ai_standard
title: AI Document Authoring Standard
status: draft
phase: 0
ai_load_priority: always
ai_summary: Defines how Sovrunn docs must be written for human review, AI reasoning, retrieval, and token optimization.
---

# AI Document Authoring Standard

## 1. Purpose

Sovrunn documentation is platform source code.

All foundation docs, RFCs, ADRs, architecture docs, resource specs, plugin specs, and task files must be optimized for:

- human review,
- AI reasoning,
- retrieval,
- token efficiency,
- stable implementation,
- architecture consistency.

## 2. Core Rule

Every design file should be:

```text
precise
structured
canonical
non-repetitive
decision-linked
retrieval-friendly
token-efficient
implementation-oriented
```

Avoid:

```text
essay-like prose
marketing-heavy wording
ambiguous language
repeated vision text
many synonyms
large paragraphs
unanchored opinions
```

## 3. Metadata Header

Every file must start with front matter.

Example:

```yaml
---
doc_type: rfc
id: RFC-0012
title: Organization, Tenant, and Governance Model
status: draft
phase: 1
depends_on:
  - DEC-0018
  - DEC-0019
  - constitution.md
ai_load_priority: feature
ai_summary: Defines Organization, OrganizationUnit, Tenant, Project, policy inheritance, and audit behavior.
---
```

## 4. AI Load Priority

Use one of:

| Priority | Meaning |
|---|---|
| always | Load in most AI sessions. |
| phase | Load for work in a specific phase. |
| feature | Load when building a specific feature. |
| reference | Load when deeper context is needed. |
| optional | Load only on demand. |

## 5. Canonical Terms

Use canonical terms from `docs/glossary.md`.

Rule:

```text
One concept = one canonical term.
```

Use:

- Organization,
- OrganizationUnit,
- Tenant,
- Project,
- ServiceInstance,
- ServiceBinding,
- Plugin,
- Capability,
- Operation,
- ServiceOps.

Avoid unmanaged synonyms.

## 6. Structure

Use numbered sections.

Prefer:

```markdown
## 1. Purpose
## 2. Scope
## 3. Goals
## 4. Non-Goals
## 5. Definitions
## 6. Design
## 7. Validation Rules
## 8. Failure Modes
## 9. Acceptance Criteria
## 10. AI Implementation Guidance
```

## 7. Normative Language

Use crisp requirement statements.

Good:

```text
Sovrunn must create an Operation for every asynchronous platform change.
```

Bad:

```text
It may be useful to create some kind of operation record in many cases.
```

## 8. Goals, Non-Goals, Future Work

Every architecture doc and RFC should include:

- Goals,
- Non-Goals,
- Future Work where useful.

This prevents AI from implementing future-phase features early.

## 9. Decision References

Reference DEC IDs when relevant.

Example:

```text
This follows DEC-0018.
```

Use ADR references when a tradeoff exists.

## 10. Tables for Contracts

Use tables for:

- resource fields,
- status fields,
- lifecycle operations,
- validation rules,
- failure modes,
- API mappings.

Tables are token-efficient and implementation-friendly.

## 11. Structured Bullets

Prefer bullets over long prose.

Good:

```text
Validation rules:
- organizationRef must exist.
- name must be unique within parent scope.
- failed resources must include status.reason.
```

## 12. Examples

Use one canonical example per concept.

Avoid many variations unless needed.

Examples should be valid YAML or code whenever possible.

## 13. Avoid Repetition

Do not repeat the full Sovrunn vision in every file.

Use references:

```text
This RFC follows constitution.md and DEC-0018.
```

## 14. AI Implementation Guidance

Every RFC must include an `AI Implementation Guidance` section.

It should say:

- what to implement,
- what not to implement,
- which files to update,
- which tests to generate,
- which architecture rules apply.

## 15. Failure Modes

Every platform feature must define failure modes.

Example:

| Failure | Expected Behavior |
|---|---|
| Parent Organization missing | reject with `ParentNotFound` |
| Quota exceeded | fail Operation with `QuotaExceeded` |
| Plugin capability missing | fail fast with `CapabilityUnsupported` |

## 16. Acceptance Criteria

Acceptance criteria must be testable.

Good:

```text
Given an Organization exists,
when an OrganizationUnit is created with that organizationRef,
then the OrganizationUnit is stored and status.phase becomes Ready.
```

Bad:

```text
Governance should work properly.
```

## 17. Token Hierarchy

Use three tiers:

| Tier | File Type | Token Rule |
|---|---|---|
| 1 | Always-loaded canonical files | short, dense, low repetition |
| 2 | Feature specs/RFCs | detailed enough for implementation |
| 3 | Deep references/examples | long only when needed |

## 18. Always-Loaded File Target

Always-loaded docs should normally stay between 1,000 and 2,500 words.

Examples:

- vision.md,
- philosophy.md,
- constitution.md,
- DECISION_INDEX.md,
- glossary.md,
- AI_CONTEXT_GUIDE.md,
- AI_FEATURE_FACTORY.md,
- AI_DOC_AUTHORING_STANDARD.md.

## 19. Feature RFC Target

Feature RFCs may be longer, normally 2,000 to 5,000 words.

## 20. AI Authoring Checklist

Before accepting a doc, verify:

- canonical terms used,
- metadata present,
- numbered sections used,
- non-goals included,
- failure modes included,
- acceptance criteria included,
- related decisions listed,
- AI implementation guidance included,
- no unnecessary repetition,
- no future-phase implementation leakage.
