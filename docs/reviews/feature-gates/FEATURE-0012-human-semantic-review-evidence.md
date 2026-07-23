# FEATURE-0012 Human Semantic Review Evidence

Feature: FEATURE-0012 — API, Resource Naming, Status, and Validation Standard
Phase: Phase 2 — Reuse-First PaaS Fabric Foundation
Task: 17.3 Collect human-review evidence and mark PENDING_HUMAN_REVIEW
Requirements: F12-VERIFY-003, F12-RISK-001
Design: verification section (F12-VERIFY-003)

## Review status

Final feature-review status: PENDING_HUMAN_REVIEW

Assessment decision status: Approved (architecture / reuse assessment only;
see ADH-2026-012 / ADH-2026-013 and
`docs/reviews/reuse-assessments/FEATURE-0012-approval-evidence.md`)

This document stages evidence for human semantic review. It does **not**
authorize final merge. A human must later create or update
`docs/reviews/feature-gates/FEATURE-0012-approval-review.md` with an explicit
`Final feature-review status: Approved` line, reviewer identity, decision, and
date. The coding agent must not write that approval.

Evidence collected at (UTC): 2026-07-23T16:54:13Z
Repository HEAD at collection: `3cb1dd7f9d4884a674f1592f5cf99babee512357`
Branch: `feature-0012-api-resource-naming-status-and-validation-standard`

## Human semantic-review checklist (F12-VERIFY-003)

These items are **not** verified by `make` alone and remain for recorded human
review (reviewer, decision, and date required):

| Item | Evidence staged for review | Human decision |
|---|---|---|
| Architecture approvals and granted adoption exceptions (F12-SCOPESTD-004) | ADH-2026-012, ADH-2026-013; reuse-assessment approval evidence | PENDING_HUMAN_REVIEW |
| Residual-risk acceptance for Matrix E (F12-RISK-001) | `.automation/evidence/FEATURE-0012/matrix-e-residual-risk-review-inputs.md` | PENDING_HUMAN_REVIEW |
| Correctness of boundary classifications and responsibility boundaries | `docs/api/boundary-ledger.yaml` / `docs/api/BOUNDARY_LEDGER.md`; fitness ledger check | PENDING_HUMAN_REVIEW |
| Adequacy of Phase 1 compatibility exceptions and migration candidates | `docs/api/PHASE1_COMPATIBILITY_REPORT.md` | PENDING_HUMAN_REVIEW |
| Baseline protected-review external settings still pending | `docs/reviews/feature-gates/FEATURE-0012-baseline-protected-review-evidence.md` | PENDING_HUMAN_REVIEW |

## Automated verification evidence (F12-VERIFY-002 — supporting only)

Authoritative machine outputs are stored under
`.automation/evidence/FEATURE-0012/` (see `README.md` and `MANIFEST.txt`).

### Docker Go 1.22 verification

Command:

```bash
docker run --rm -v "$PWD":/src -w /src ${GO_DOCKER_IMAGE:-golang:1.22} sh -c 'test -z "$(gofmt -l .)" && go vet ./... && go test ./...'
```

Result: exit 0 (see `docker-verify.out` / `docker-verify-summary.txt`)

### Feature gate

Command:

```bash
make ff-feature-gate FEATURE=FEATURE-0012
```

Result: exit 0 / `SUCCESS: FEATURE-0012 passed Sovrunn feature gate`
(see `feature-gate.out` / `feature-gate-summary.txt`)

Notes from gate output (not human approval):

- WARN: `docs/reviews/feature-gates/FEATURE-0012-approval-review.md` not found
- WARN: final merge approval still requires `Final feature-review status: Approved`
- Task 17.2 automation flow log: PASS (`task-17.2-flow-log.json`)

### Changed-file inventory

| Inventory | Path |
|---|---|
| FEATURE-0012-scoped (primary) | `.automation/evidence/FEATURE-0012/changed-files.txt` |
| FEATURE-0012 commits | `.automation/evidence/FEATURE-0012/commits-feature-0012.txt` |
| Branch vs `main` (context) | `.automation/evidence/FEATURE-0012/changed-files-vs-main.txt` |

Primary inventory size at collection: 231 paths from first FEATURE-0012
architecture commit parent through HEAD.

## Compatibility-report evidence

| Field | Value |
|---|---|
| Path | `docs/api/PHASE1_COMPATIBILITY_REPORT.md` |
| SHA-256 | `73980ef7bbe634bc65f9bb725469681f7d0f0109820c1b27390f5854db570622` |

Human reviewers should confirm:

- every required Phase 1 contract is covered;
- exceptions are explicit and adequate;
- no wholesale Phase 1 rewrite was triggered;
- migration candidates are acceptable for later approved work.

## Boundary-ledger evidence

| Field | Value |
|---|---|
| YAML source of truth | `docs/api/boundary-ledger.yaml` |
| YAML SHA-256 | `4c1568e1cf99ca255dee2b06969575069f8e9be870e79ea1cc4a6eb8904c493f` |
| Generated Markdown view | `docs/api/BOUNDARY_LEDGER.md` |
| Markdown SHA-256 | `5c83c2e7a4381291078f309e691317b594b4cb970d2b35e371ad848c4c1a24dc` |

Human reviewers should confirm boundary purpose/owner/producers/consumers,
allowed/prohibited data, authorization, audit, observability, failure
behavior, versioning, replacement/migration paths, and reassessment triggers.

## Matrix E residual-risk review inputs

See `.automation/evidence/FEATURE-0012/matrix-e-residual-risk-review-inputs.md`
for F12-R01..F12-R16 with control/detection pointers and blank human
acceptance fields.

## Observability / audit / correlation notes for reviewers

Preserved by FEATURE-0012 grammar (no new runtime audit/operation service):

- Problem Details carry `requestId` for request correlation.
- Operation contract fields support target/action/requester/idempotency
  correlation for later FEATURE-0013 adoption.
- Structured, secret-free error posture is retained.
- Boundary ledger records per-boundary observability and audit expectations.

Intentionally not logged / not emitted by this feature:

- secrets, credentials, tokens, private keys, connection strings;
- raw secret-like values in metadata/labels/annotations/status/errors;
- existence-disclosing cross-scope denial details (SafeDenial 404 path).

## Explicit non-claims

- Automated PASS results are **not** human semantic approval.
- Architecture / reuse-assessment approval is **not** final feature merge
  approval.
- This artifact does **not** write an approval token.
- This artifact does **not** invent a reviewer or review date.
- This artifact does **not** accept Matrix E residual risk.
- `make` / feature-gate success does **not** satisfy F12-VERIFY-003 alone.

## Remaining human action

1. Review staged evidence under `.automation/evidence/FEATURE-0012/` and the
   documents listed above.
2. Complete Matrix E residual-risk acceptance (owner, corrective path,
   reassessment trigger) or reject/defer with recorded rationale.
3. After semantic review, record final merge approval separately in
   `docs/reviews/feature-gates/FEATURE-0012-approval-review.md` with:
   - `Final feature-review status: Approved`
   - reviewer identity
   - decision
   - date
4. Complete external baseline branch-protection PENDING items from Task 17.1
   evidence if still open.

Until those human actions occur, feature status remains
**PENDING_HUMAN_REVIEW**.
