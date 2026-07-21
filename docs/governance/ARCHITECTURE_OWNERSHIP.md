# Architecture Ownership

This document defines ownership roles for long-term Sovrunn architecture and development.

## Architecture Owner

Responsible for:

- approving architecture baseline changes,
- accepting DEC/RFC changes,
- resolving architecture conflicts,
- protecting phase scope,
- approving replacements or new decisions.

## Feature Owner

Responsible for:

- feature requirements,
- design alignment,
- task sequencing,
- acceptance criteria,
- reuse assessment,
- architecture drift checks.

## Implementation Owner

Responsible for:

- code implementation,
- tests,
- error handling,
- observability,
- security hygiene,
- documentation updates.

## Reviewer

Responsible for:

- architecture drift review,
- reuse-before-build review,
- feature gate result review,
- test/lint/security validation.

## Security Reviewer

Responsible for:

- auth/authz implications,
- secret handling,
- tenant isolation,
- audit requirements,
- threat model implications.

## Observability Reviewer

Responsible for:

- logs,
- metrics,
- traces,
- request/operation correlation,
- audit-vs-log separation,
- no sensitive data in diagnostics.

## Product Owner

Responsible for:

- customer outcome clarity,
- MVP scenario acceptance,
- feature priority,
- roadmap revalidation.

## Multi-Developer Rule

No developer should implement a feature from roadmap placeholders alone. Implementation must start from approved feature specs that trace back to baseline, DEC/RFC, and architecture docs.
