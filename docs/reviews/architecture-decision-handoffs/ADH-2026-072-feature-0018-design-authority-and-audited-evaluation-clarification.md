# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-072
- Date: 2026-09-02
- Source discussion: ChatGPT Project / Codex
- Related feature: FEATURE-0018 — Governance, IAM, Approval and Exception Foundation
- Related phase: Phase 2R
- Author: ChatGPT / Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

FEATURE-0018 design authority, local publication, and audited-evaluation clarification

## Summary

This clarification closes three interpretation gaps that caused repeated design
reviews to reopen already-approved FEATURE-0018 decisions. It distinguishes
FEATURE-0013 contract and validation ownership from FEATURE-0018-local publication
machinery, defines exactly when an internal authorization candidate becomes an
audited authorization evaluation, and confirms the implementation-mechanism choices
owned by FEATURE-0018 design. It adds no resource, action, route, state, writer,
error, conformance identifier, dependency, external effect, or runtime capability,
and it does not change F18-RD-01..24, REQ-F18-01..24, AC-F18-01..49, or their
ADH-2026-071 mappings.

## Classification

- Clarification

## Existing approved baseline

The active baseline is `ARCH-2026.08-PHASE2R-CANONICAL`. DEC-0060 and the approved
ADH-2026-070 three-file package remain the sole FEATURE-0018 product-semantic
authority. ADH-2026-071 remains the non-semantic conformance-executability
authority. In particular:

- F18-RD-19 adopts exactly three FEATURE-0013 DecisionProfiles and the existing
  DecisionRecord/AuditEvent envelope, validation, projection, validity, and linkage
  semantics;
- F18-RD-20 requires protected audit-obligation acceptance before an accepted
  mutation is published and before an AuthorizationResult is returned;
- F18-RD-21 orders concurrency and idempotency before required decision/audit
  evidence; and
- F18-RD-22 requires deterministic, side-effect-free, FEATURE-0018-local
  in-memory conformance and delegates concrete construction, validation, copying,
  and deterministic time-supply mechanics to design.

The approved architecture also requires all package authorities and canonical
registrations to agree. No ADH-2026-070 package file or registered artifact wins
over another when they conflict; a mismatch is an architecture blocker.

Relevant baseline references:

- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/decisions/DECISION_INDEX.md`
- `docs/decisions/DEC-0060-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-feature-0018-governance-iam-approval-exception-foundation.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-a-feature-0018-semantic-contract.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-070-appendix-b-feature-0018-registries-and-evidence.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-071-feature-0018-conformance-executability-reconciliation.md`
- `.kiro/specs/governance-iam-approval-exception-foundation/requirements.md`

## Decision or proposed decision

### 1. FEATURE-0013 and FEATURE-0018 ownership boundary

FEATURE-0013 owns the shared DecisionProfile, EvaluationResult, DecisionRecord,
and AuditEvent contracts; their carrier types; their common validation, projection,
validity, authority, reason, obligation, and audit-linkage rules; and the approved
extension mechanism used to register the three FEATURE-0018 profiles and the bounded
AuditEvent taxonomy.

FEATURE-0013 does not provide FEATURE-0018 with a transaction, acceptance service,
publisher, durable store, outbox, queue, allocator, lock manager, idempotency store,
or prepared-evidence commit API. The phrase “FEATURE-0013 local atomic boundary” in
F18-RD-20 means that FEATURE-0018 must preserve the inherited evidence contract and
its acceptance-before-publication invariant. It does not name an existing callable
FEATURE-0013 service and does not transfer mutation or transaction ownership to
FEATURE-0013.

FEATURE-0018 may implement a deterministic, feature-local, process-local,
non-durable in-memory publication mechanism that stages and publishes, as one
all-or-nothing aggregate change:

```text
FEATURE-0018 domain mutation(s), when applicable
+ required FEATURE-0013-compatible DecisionRecord/AuditEvent carrier values
+ completed caller idempotency result, when applicable
```

That mechanism is FEATURE-0018 implementation and conformance machinery. It must
consume FEATURE-0013 carrier construction and validation rules without creating a
second evidence envelope or writer, must perform zero external I/O, and must never
claim production durability or a FEATURE-0013 acceptance service.

### 2. Authorization candidate and audited evaluation boundary

Stages before concurrency/idempotency may compute one sealed internal authorization
candidate. A candidate is provisional implementation data only. It is not an
`AuthorizationResult`, DecisionRecord, AuditEvent, bearer value, completed
authorization evaluation, or downstream authority.

A candidate becomes a completed authorization evaluation only after every
operation-applicable check through F18-RD-21's concurrency/idempotency stage has
succeeded and the operation is ready to produce one exact current `Allow | Deny`
AuthorizationResult. At that point, the required `authorization.evaluated`
obligation and any mandatory DecisionRecord are prepared, validated, accepted, and
published before the result is released. Evidence failure returns the existing safe
internal error, publishes no result or completed replay, and grants no downstream
authority.

The following stage-10 terminal or coordination outcomes do not produce an
AuthorizationResult and therefore do not complete an authorization evaluation:

- request-binding `Conflict`;
- terminal optimistic-concurrency/CAS loss;
- `Wait` or retry/re-evaluation direction; and
- controller effective-once `NoOp`.

Their already-registered conflict, lifecycle, concurrency, or mutation-attempt
behavior remains unchanged; this handoff adds no new AuditEvent type or error.
Discarding a provisional candidate on one of these paths does not violate F18-RD-19,
because no completed `Allow | Deny` evaluation exists.

A completed-result replay is different: after current authentication,
authorization, and safe target visibility are re-evaluated and stage 10 selects the
matching stored result, that current evaluation must accept and publish exactly its
required authorization evidence before the stored successful result is returned.
Replay never extends the original effect or validity.

### 3. Design-owned implementation mechanics

Within the exact F18-RD, REQ, AC, route, writer, error, dependency, and Phase 2R
boundaries, FEATURE-0018 design owns the concrete implementation mechanics necessary
to make the approved behavior executable. This includes:

- package decomposition and acyclic import direction;
- sealed interfaces and owner-private constructors;
- immutable fixture construction, copying, validation, and deterministic clocks;
- process-local state roots, snapshots, locks, reservations, waiters, optimistic
  version checks, and idempotency tables;
- the FEATURE-0018-local non-durable unit-of-work protocol, transaction modes,
  context-bound change receipts, seal invariants, commit/abort, panic cleanup, and
  shutdown coordination;
- FEATURE-0013-compatible carrier preparation and validation wiring;
- literal route shapes for the already-closed caller operation surface; and
- static architecture-checker symbols, package/call rules, diagnostics, fixtures,
  and build integration.

These are design choices, not new product semantics. Design may not use them to add
or change a resource, action, field authority, scope, target binding, lifecycle,
eligibility rule, approval rule, audit obligation, error, route operation, writer,
dependency, persistence guarantee, external integration, or adjacent-feature
behavior. If a mechanism requires such a choice, design must stop with
`ARCHITECTURE_DECISION_REQUIRED`.

The requirements statement that no other concern is delegated means no other
observable product semantic or ownership decision is delegated. It does not prohibit
the design mechanics listed above.

### 4. Authority consistency and stale transcription

F18-RD-02, the approved requirements, the canonical contract catalog, the canonical
model, and the VS-000 registry must agree exactly for the fields they register. No
design may select one as the winner. Any mismatch stops with
`ARCHITECTURE_DECISION_REQUIRED`; it is never resolved by source precedence inside
FEATURE-0018 design.

The existing design sentence saying that the registry wins is a stale transcription
and must be removed. Correcting it does not reopen F18-RD-02.

### 5. FEATURE-0018 design-review classification rule

Every blocking or revision finding in a FEATURE-0018 design review must use exactly
one of these classifications:

| Classification | Exact use |
|---|---|
| `STALE_TRANSCRIPTION` | The design contradicts or duplicates a closed authority. Correct the design by exact citation/transcription; do not reopen the decision. |
| `DESIGN_EXECUTABILITY` | Approved behavior is closed, but the proposed package, API, state, concurrency, idempotency, evidence, checker, or test mechanism is incomplete or internally inconsistent. Correct design only. |
| `REQUIREMENT_GAP` | An approved architecture statement is absent or weakened in requirements. Stop design and reconcile requirements without inventing semantics. |
| `ARCHITECTURE_CLARIFICATION_REQUIRED` | Loaded approved authorities do not determine an observable semantic or ownership outcome. Stop and obtain bounded human-approved clarification. |
| `OUT_OF_SCOPE` | The proposed or requested concern belongs to an excluded/future feature, dependency owner, production runtime, or external system. Remove or defer it without designing the future feature. |

A reviewer may not classify a closed decision as a requirement or architecture gap
without citing the exact contradictory or missing authority and the exact target
text. A finding about an implementable mechanism must use
`DESIGN_EXECUTABILITY`, even when the mechanism protects a security invariant.
Repeated wording, preference, or an alternative implementation is not evidence that
a closed architecture decision is open.

## Rationale

The earlier review epochs repeatedly mixed product semantics, dependency ownership,
and implementation mechanics. This clarification supplies one stable interpretation:
FEATURE-0013 remains the evidence-contract owner; FEATURE-0018 owns its local
non-durable coordination; the audit obligation attaches to a completed Allow/Deny
evaluation after applicable stage-10 outcomes; and design has enough bounded latitude
to make the approved behavior executable without returning for new architecture on
every lock, receipt, interface, or checker detail.

The classification rule makes reviewers identify whether a finding is stale text,
an executable-design defect, a requirements omission, a genuine semantic gap, or an
excluded concern. Requiring an exact authority citation prevents already-closed IAM
decisions from being reopened merely because the design uses different wording.

## Reuse-before-build assessment

- Disposition: Reuse
- Summary of mature candidates / applicable standards:
  - Reuse ADH-2026-070 F18-RD-19 through F18-RD-22, ADH-2026-071 traceability,
    FEATURE-0013 carrier/validation contracts, FEATURE-0012 concurrency and
    idempotency semantics, and the existing generic Feature Factory review route.
- Sovrunn-owned responsibility summary:
  - FEATURE-0018 owns only its deterministic local publication coordination,
    design mechanics, exact evidence adoption, and local conformance.
- Non-goals summary:
  - No FEATURE-0013 service or contract extension; no production transaction,
    persistence, outbox, queue, allocator, route runtime, adapter, external I/O,
    new audit event, new error, or future-feature behavior.

## Phase impact

- Current phase allowed? Yes
- If not current phase, target phase: Not applicable
- Current phase boundary impact:
  - Clarification and deterministic Phase 2R design/conformance machinery only.
    Zero external I/O and non-durable in-memory boundaries remain unchanged.

## Conflict check

- Conflicts with accepted DEC/RFC? No
- Conflicting decisions, if any:
  - None. The handoff clarifies the interpretation of F18-RD-19 through
    F18-RD-22 without changing their observable outcomes.
- Resolution required:
  - Approved clarification handoff and design reconciliation only. No new DEC,
    RFC, ACR, baseline replacement, REQ, AC, route, or conformance registration.

## Required action

- Update architecture doc
- Update Kiro design.md
- Update feature control/review context
- Do not update Kiro requirements.md
- Do not generate Kiro tasks.md

## Impacted files

- `docs/reviews/architecture-decision-handoffs/ADH-2026-072-feature-0018-design-authority-and-audited-evaluation-clarification.md`
- `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `docs/features/FEATURE-0018-governance-iam-approval-exception-foundation.md`
- `.automation/features/FEATURE-0018.control.json`
- `.automation/context-projections/FEATURE-0018.design-review.spec.json`
- `.automation/context-projections/FEATURE-0018.design-review.md`
- `.automation/context-projections/FEATURE-0018.design-review.manifest.json`
- `.kiro/specs/governance-iam-approval-exception-foundation/design.md`
- `.automation/state/FEATURE-0018.json` only for approved hash reconciliation

The approved requirements remain byte-identical. Generated review outputs are
regenerated by the controlled review route and are not hand-edited.

## Impacted features

- FEATURE-0018: receives an executable design clarification only.
- FEATURE-0013: no change; its contract, validation, extension, and writer ownership
  remain unchanged.
- FEATURE-0012, FEATURE-0016, FEATURE-0017, and FEATURE-0019..0026: no change.

## Acceptance criteria for Kiro update

- [ ] Handoff is validated against the Architecture Operating System files.
- [ ] No F18-RD, REQ, AC, route operation, writer, state, error, conformance ID,
  dependency, or Phase 2R boundary changes.
- [ ] The design contains no registry-wins or other package-internal precedence rule.
- [ ] FEATURE-0013 is described only as carrier/contract/validation/extension owner;
  the local in-memory unit of work is explicitly FEATURE-0018-owned and non-durable.
- [ ] A provisional candidate is explicitly non-authoritative and becomes a completed
  audited evaluation only after applicable stage-10 success.
- [ ] Conflict, terminal CAS loss, Wait/retry, and controller NoOp produce no
  AuthorizationResult and are not represented as completed authorization evaluations.
- [ ] Replay performs fresh checks and publishes required current authorization
  evidence before returning the stored result.
- [ ] The transaction API contains enforceable context-bound applied-change receipts
  or an equivalent mechanism proving participant, owner, context, shadow/index,
  descriptor, and originating-result linkage at Seal.
- [ ] The review context contains the exact five-class finding rule and requires exact
  contradictory/missing authority citations before reopening closed decisions.
- [ ] Inherited semantics in design are replaced by exact citations wherever the text
  does not define a design-owned route, package, API, lock, publication mechanism,
  checker rule, or test.
- [ ] Requirements remain byte-identical and retain their approved digest.
- [ ] No tasks or implementation are generated.
- [ ] FEATURE-0018 requirements/design semantic checks, feature-control validation,
  strict documentation build, and applicable drift checks pass.

## Explicit instructions to Kiro

- Treat ADH-2026-070 and DEC-0060 as the sole product-semantic authority and
  ADH-2026-071 as the conformance-mapping authority.
- Use this handoff only to interpret the three boundaries above and classify design
  findings; do not add a 25th F18-RD group.
- Correct stale design text without requesting a new architecture decision.
- Resolve concrete API, receipt, package, lock, checker, and test completeness as
  `DESIGN_EXECUTABILITY` within the exact approved boundary.
- Do not modify requirements, routes, resources, actions, writers, errors,
  conformance mappings, dependencies, production behavior, or Phase 2R scope.
- Do not change FEATURE-0013 or claim that it provides a transaction or acceptance
  service.
- Do not generate tasks or code.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-09-02
- Notes: Approval is limited to the five corrections explicitly approved in the
  source discussion: the three clarification boundaries, the stale precedence fix,
  context-bound receipt completion, exact reviewer classification, and design
  deduplication. No product-semantic expansion is authorized.

