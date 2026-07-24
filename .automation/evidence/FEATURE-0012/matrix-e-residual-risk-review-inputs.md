# FEATURE-0012 Matrix E — Residual-risk review inputs

Feature: FEATURE-0012
Requirement: F12-RISK-001
Human semantic review: F12-VERIFY-003

Status: **PENDING_HUMAN_REVIEW**

This worksheet stages Matrix E risks and implemented controls for human
residual-risk review. The coding agent does **not** accept residual risk,
assign an owner, invent a reviewer, invent a review date, or record approval.

Source register: `.kiro/specs/api-resource-naming-status-and-validation-standard/requirements.md` §4.17

## Human decision fields (must be filled by a human)

| Field | Value (human only) |
|---|---|
| Reviewer identity | _PENDING_HUMAN_REVIEW — not filled by agent_ |
| Review date | _PENDING_HUMAN_REVIEW — not filled by agent_ |
| Residual-risk decision | _PENDING_HUMAN_REVIEW — Accept / Reject / Defer_ |
| Decision notes | _PENDING_HUMAN_REVIEW_ |

## Matrix E review table

For each risk, a human must confirm control adequacy and record residual
acceptance (owner, corrective path, reassessment trigger) before residual
risk may be treated as accepted.

| ID | Risk | Primary control (implemented / documented) | Detection / response evidence | Human residual acceptance |
|---|---|---|---|---|
| F12-R01 | Overfit to Phase 1/Kubernetes | Profile matrix; provider-neutral rules | Future-scenario fixtures; ADH for exceptions; fitness check 2/14 | PENDING_HUMAN_REVIEW |
| F12-R02 | Universal shape misrepresents lifecycle | Eight explicit profiles | Schema profile required (`x-sovrunn-profile`); fitness check 1 | PENDING_HUMAN_REVIEW |
| F12-R03 | Scope/location/source/ownership conflated | Separate typed concepts (scopeRef/ownerRef/TypedRef) | Scope/reference negative fixtures; Property 4/11 | PENDING_HUMAN_REVIEW |
| F12-R04 | Provider-native data leaks into core/customer | Adapter-only native contracts; boundary ledger | Import/schema lint; Property 7; fitness check 2 | PENDING_HUMAN_REVIEW |
| F12-R05 | Plugin and adapter semantics conflated | Separate boundaries and resource families | Boundary ledger; Matrix D fixtures; fitness check 14 | PENDING_HUMAN_REVIEW |
| F12-R06 | Generic/name-only refs cause access/staleness errors | Typed refs + optional immutable UID | Kind/scope and name/UID mismatch tests; apiref constraints | PENDING_HUMAN_REVIEW |
| F12-R07 | Multiple writers corrupt status | One owner per field/condition; FieldPolicy | Ownership fitness check 3; operation-aware decode modes | PENDING_HUMAN_REVIEW |
| F12-R08 | Conditions become history or phase explodes | Current-fact conditions; optional phase | apicond transition Property 8; bounded MaxConditions | PENDING_HUMAN_REVIEW |
| F12-R09 | Stale observations produce inaccurate decisions | Provenance / observed time / freshness on observed resources | DiscoveredDatabase fixture/schema; fitness check 8 | PENDING_HUMAN_REVIEW |
| F12-R10 | Secrets/restricted data leak via metadata/errors | Classification, redaction, SecretRef-only | Fitness check 7; negative secret-like fixture scans | PENDING_HUMAN_REVIEW |
| F12-R11 | Extensions become shadow APIs | Registered namespaced schemas only | Unknown `x-sovrunn-*` fail-closed; fitness check 1a | PENDING_HUMAN_REVIEW |
| F12-R12 | API evolution breaks clients/plugins | Maturity + compatibility policy; baseline gate | Schema-diff; VerifyBaselineIntegrity/Approval; fitness check 10 | PENDING_HUMAN_REVIEW |
| F12-R13 | Unbounded objects/status/lists prevent scale | Finite Limits + opaque pagination | Limit defaults; fitness check 11; over-limit negatives | PENDING_HUMAN_REVIEW |
| F12-R14 | Documentation-only standard drifts | Shared primitives + executable conformance | Feature gate + api-conformance-check; fitness 1–15 | PENDING_HUMAN_REVIEW |
| F12-R15 | Phase 1 migration expands into a rewrite | Compatibility audit + explicit exceptions | `docs/api/PHASE1_COMPATIBILITY_REPORT.md`; fitness coverage | PENDING_HUMAN_REVIEW |
| F12-R16 | AI/privileged consumers bypass boundaries | AI consumes authorized filtered views only | Boundary ledger / field-policy; SafeDenial path tests | PENDING_HUMAN_REVIEW |

## Required human residual-acceptance record (template)

When a human accepts residual risk for any Matrix E item, record at least:

```text
risk_id: F12-Rxx
accepted: yes|no
owner: <human role or identity>
corrective_path: <action if risk materializes>
reassessment_trigger: <event that forces re-review>
reviewer: <human>
date: <YYYY-MM-DD>
```

Do not treat this evidence pack, Task 17.2 automation PASS, or architecture
ADH approvals as residual-risk acceptance.
