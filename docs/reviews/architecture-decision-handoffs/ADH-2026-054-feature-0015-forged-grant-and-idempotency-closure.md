# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-054
- Date: 2026-08-13
- Source discussion: FEATURE-0015 executable-design review recovery
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Close FEATURE-0015 forged-grant carrier and target-bound idempotency semantics

## Summary

This handoff closes the two remaining externally observable gaps found during
FEATURE-0015 design review. First, the existing audit contract requires an
authenticated header/body forged-grant attempt to return an audited
`AUTHORIZATION_DENIED` outcome, but it does not identify the exact header or
body member or the precedence relative to malformed and unknown-body handling.
Second, existing idempotency wording does not bind an existing-item action to
its concrete target UID or require current authorization before returning a
stored replay outcome. The correction defines both contracts without adding a
resource, lifecycle state, route, grant, Problem code, violation code, or
persistence dependency.

## Classification

Correction — missing externally observable input and idempotency identity
semantics in already-approved FEATURE-0015 behavior.

## Existing approved baseline

- `VS0-CF-F15-22` requires an authenticated header/body forged-grant attempt
  to return `AUTHORIZATION_DENIED`, ignore the claim, and emit exactly one
  redacted FEATURE-0013 AuditEvent.
- `REQ-F15-17` requires server-resolved bootstrap grants only; grants are never
  accepted from request headers or bodies.
- `REQ-F15-18`, `VS0-CF-F15-18`, and `VS0-CF-F15-19` require same-key replay
  and changed-digest conflict behavior for all F0015 collection creates and
  eight participation item actions.
- ADH-2026-048 requires validation failure, stale-version failure,
  audit-append failure, and recovered-panic handling to abort/release an
  in-flight idempotency reservation.
- ADH-2026-049 and ADH-2026-051 require append-before-publication for required
  AuditEvents, `INTERNAL_ERROR`/500 on required audit-append failure, no
  durable cross-store transaction claim, and allow an already-appended event
  to survive a later in-process publication failure.
- ADH-2026-053 corrected only the participation action URI shape and remains
  controlling.

## Correction

### 1. Exact forged-grant carriers and precedence

The only F0015 request carriers treated as an attempt to submit or substitute
a bootstrap grant are:

- HTTP header `X-Sovrunn-Bootstrap-Grant`, when present with any value; and
- top-level JSON member `bootstrapGrant`, when present in a syntactically valid
  and duplicate-free request body, with any JSON value.

Neither carrier is trusted, parsed as authorization, stored, defaulted, or
disclosed. Only the server-resolved current grant authorizes an operation.

The required precedence is:

1. Missing or invalid authentication returns existing `AUTH_REQUIRED`/401 with
   no durable AuditEvent.
2. After successful authentication, a present
   `X-Sovrunn-Bootstrap-Grant` header returns `AUTHORIZATION_DENIED`/403,
   produces the already-approved redacted AuditEvent, and stops before body
   decoding, current authorization, root gating, or reference resolution.
3. Without that header, phase-one decoding applies the existing bounded-body,
   JSON syntax, and duplicate-member checks. Their existing 400 outcomes take
   precedence over body-carrier detection and produce no durable AuditEvent.
4. After a syntactically valid, duplicate-free phase-one decode, a top-level
   `bootstrapGrant` member returns `AUTHORIZATION_DENIED`/403, produces the
   already-approved redacted AuditEvent, and stops before current
   authorization, root gating, reference resolution, idempotency, or normal
   phase-two unknown-field classification.
5. Any other unknown body member follows the existing closed-contract
   `UNKNOWN_FIELD`/400 behavior with no durable AuditEvent. No generic
   authorization-like-header heuristic is introduced.

This creates one explicit, narrow exception to the existing F15-30
non-empty-item-action-body rule: a syntactically valid, duplicate-free,
non-empty item-action body whose top-level member set includes
`bootstrapGrant` is governed by `VS0-CF-F15-22` and returns audited
`AUTHORIZATION_DENIED`/403. A non-empty item-action body without that member
remains governed by `VS0-CF-F15-30` and returns `MALFORMED_REQUEST`/400 with
no durable AuditEvent. The canonical REQ-F15-23 wording and exact
`VS0-CF-F15-30` ledger row must be updated to express that boundary.

### 2. Target-bound, currently-authorized idempotency

The idempotency namespace is:

```text
(authenticatedPrincipalUID, registered method/path pattern,
 concreteItemTargetUID when the route has {uid}, Idempotency-Key)
```

Collection-create routes have no item target UID; their successful classified
canonical request digest continues to distinguish distinct create requests.

For a same-key, same-digest completed replay, the API server must first repeat
current authentication, current server-resolved authorization, and safe
target/reference access for the incoming request. A revoked, narrowed, or
otherwise insufficient current grant returns its established denial outcome;
the stored result is not disclosed. A completed replay record for one concrete
item target must never satisfy a request for another target.

Only successful collection creates and successful participation item-action
transitions complete an idempotency record. Validation failure, uniqueness
failure, lifecycle-source-state failure, stale `If-Match`, required
AuditEvent-append failure, recovered panic, and graceful shutdown abort the
in-flight reservation, wake waiting callers, and leave no completed replay
record. A changed digest in the same namespace retains the existing
`CONFLICT`/409 `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH` outcome.

### 3. Design mechanics that must be made explicit

Without requesting another architecture decision, the design must:

- compare current `If-Match` after an authorized same-digest replay lookup but
  before lifecycle source-state validation;
- append every required audited-denial AuditEvent before publishing its denial
  response; required append failure replaces that response with
  `INTERNAL_ERROR`/500;
- state the already-approved limitation that an AuditEvent successfully
  appended before a later in-process publication failure may remain, without
  claiming durable cross-store atomicity; and
- enumerate the replayable and abort-only outcome classes exactly.

## Rationale

The correction makes the already-approved audit and idempotency behavior
executable without trusting client authorization claims, leaking a stored
result after access changes, or allowing a replay record for one participation
to satisfy another participation's action. It also fixes the required
precedence between malformed input and a body-carried forged claim.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse FEATURE-0012 `AUTH_REQUIRED`, `AUTHORIZATION_DENIED`,
    `UNKNOWN_FIELD`, `CONFLICT`, `STALE_RESOURCE_VERSION`, and
    `INTERNAL_ERROR` Problem Details outcomes.
  - Reuse FEATURE-0013 AuditEvent append and redaction contracts.
  - Reuse Go 1.22 request header access, JSON decoding, and standard-library
    SHA-256 facilities; no new router or persistence dependency.
- Sovrunn-owned responsibility summary:
  - Enforce the two exact forged-grant carriers, target-bound current-grant
    replay checks, and abort-only failure outcomes.
- Non-goals summary:
  - No new API route, resource field, action grant, lifecycle state,
    authorization role, Problem code, violation code, durable storage,
    cross-store transaction, router library, or future-feature behavior.

## Phase impact

- Current phase allowed: Yes.
- Current phase boundary impact: FEATURE-0015 contract-closure correction
  only; no future-feature capability is introduced.

## Conflict check

- Conflicts with accepted DEC/RFC: No.
- Resolution required: Founder approval of this correction; no DEC, RFC,
  baseline, roadmap, or feature-sequence update.

## Required action

- Update the F0015 registry, contract specification, architecture authority,
  feature authority, traceability, closure matrix, steering, control manifest,
  and deterministic checkers for the exact carrier, precedence, and
  target-bound idempotency contract.
- After authority validation, update only the affected requirements artifact
  text, canonical REQ-F15-17/18/23 rows, and its verbatim F15-18, F15-19,
  F15-22, and F15-30 conformance-ledger rows.
- Update only the affected design sections: two-phase decoding, pipeline order,
  idempotency/replay matrix, current-version ordering, audited-denial
  publication, test plan, and formal idempotency/audit model.
- Do not add or renumber REQ, AC, conformance, route, resource, lifecycle,
  Problem-code, or violation-code identifiers.
- Do not modify Go source or generate tasks.

## Strict Kiro read/write boundary

Kiro may read and write only:

- this handoff;
- `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`;
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`;
- `docs/traceability/VS-000_CONTRACT_TRACEABILITY_MATRIX.md`;
- `docs/reviews/architecture-readiness/FEATURE-0015-architecture-closure-matrix.md`;
- `.kiro/steering/slice0-contract.md`;
- `.automation/features/FEATURE-0015.control.json`;
- `scripts/feature-0015-architecture-readiness-check.py`;
- `scripts/vs000-contract-check.py`;
- `scripts/feature-contract-check.py`;
- `docs/formal/feature-0015/IdempotencyAuditPublication.tla`.
- `docs/formal/feature-0015/IdempotencyAuditPublication.cfg`.

Kiro must not search, glob, inspect, or follow references outside this list.
If another file is required, it must stop and report its exact path and reason
without opening it.

## Impacted features

- FEATURE-0015: forged-grant safety, replay isolation, and executable design
  closure.
- FEATURE-0016 through FEATURE-0026: no behavior change.

## Acceptance criteria for Kiro update

- [ ] `VS0-CF-F15-22` names exactly the two approved forged-grant carriers and
  their 401/403/400 precedence.
- [ ] `VS0-CF-F15-30` remains the exact no-AuditEvent malformed-action-body
  case only when the syntactically valid, duplicate-free body does not contain
  the reserved top-level `bootstrapGrant` member.
- [ ] `VS0-CF-F15-18` and `VS0-CF-F15-19` bind item-action replay/mismatch to
  the concrete `{uid}` target and require current authorization/safe access
  before any completed replay response.
- [ ] F0015 authorities state that only successful collection creates and
  successful participation item actions complete replay records; all listed
  failure classes abort and wake waiters.
- [ ] F0015 authorities state current `If-Match` precedes lifecycle source-state
  validation after authorized replay lookup.
- [ ] F0015 authorities state append-before-response for every audited denial,
  the `INTERNAL_ERROR`/500 substitution on required append failure, the allowed
  already-appended-event survivor, and no durable cross-store atomicity claim.
- [ ] The formal idempotency model proves no completed replay is shared across
  distinct concrete targets and that abort-only outcomes never complete a
  replay record.
- [ ] `FEATURE-0015.control.json` retains ADH-048 through ADH-053 and adds
  ADH-054 with no other field change.
- [ ] No Go source, task artifact, DEC, baseline, roadmap, feature sequence, or
  state file changes.
- [ ] `make feature-0015-architecture-readiness`,
  `make feature-contract-check FEATURE=FEATURE-0015`,
  `make feature-0015-formal-check`, `make vs000-contract-check`,
  `make phase2r-drift-check`, `git diff --check`, and `mkdocs build --strict`
  pass.

## Foreseeable-gap closure matrix

The following design-review gaps are closed by this handoff or explicitly
delegated as already-approved deterministic mechanics. Kiro must not create a
further architecture decision for an item in this table.

| Area | Closure in this handoff | Required proof |
|---|---|---|
| Forged-grant carrier ambiguity | Exactly one reserved header and one reserved top-level JSON member; no generic heuristic. | `VS0-CF-F15-22` and carrier/precedence readiness checks. |
| 401/403/400 precedence | Authentication first; header 403 next; malformed/duplicate body 400 next; valid body carrier 403 next; ordinary unknown field 400 last. | Two-phase decoder tests and `VS0-CF-F15-22`/`F15-30`. |
| Action replay across targets | Namespace contains the concrete `{uid}` for every item action. | Target-isolation conformance and TLA+ invariant. |
| Replay after access change | Current authentication, authorization, and safe target/reference access precede every completed replay response. | Revoked/narrowed-grant replay tests. |
| Replayable failure ambiguity | Only successful collection creates and successful participation actions complete; all listed failure classes abort and wake waiters. | Complete replay/abort matrix and TLA+ invariant. |
| Stale-version/lifecycle precedence | Current `If-Match` is checked before lifecycle source-state validation. | Stale-invalid-state handler test. |
| Audited-denial publication | Required audit append precedes every audited denial response; append failure substitutes `INTERNAL_ERROR`/500. | Audited-denial append-failure tests. |
| In-process failure after audit append | Already-appended AuditEvent may remain; no durable cross-store atomicity is claimed. | TLA+ model and failure-path test. |

No additional externally observable decision is left open by these closures.
Implementation details such as JSON parser APIs, lock/data-structure choice,
waiter primitive, and bounded in-process retention remain design-owned, so
long as they preserve the exact contracts above.

## Explicit instructions to Kiro

- Apply no decision beyond this handoff.
- Do not search outside the strict allowlist.
- Do not modify Go source, tasks, DEC, baseline, roadmap, feature sequence, or
  state files.
- Do not add a generic claim-carrier heuristic. Only the two exact carriers in
  this handoff receive the audited forged-grant outcome.
- Do not add a durable idempotency store, durable cross-store transaction, or
  new top-level Problem/violation code.
- Stop, do not infer, if this handoff cannot be represented with only the
  allowed files.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-13
- Notes: Approved bounded closure for forged-grant carriers, target-bound idempotency, and executable design ordering.
