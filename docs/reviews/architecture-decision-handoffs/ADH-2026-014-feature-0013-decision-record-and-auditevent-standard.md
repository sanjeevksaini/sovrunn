---
doc_type: architecture_decision_handoff
handoff_id: ADH-2026-014
feature: FEATURE-0013
title: Decision Record and AuditEvent Standard
status: Proposed
classification: Extension
phase: 2
approval_date: null
approving_role: null
---

# ADH-2026-014 — FEATURE-0013 Decision Record and AuditEvent Standard

## Metadata

- Handoff ID: ADH-2026-014
- Date: 2026-07-27
- Source discussion: Sovrunn Architecture Governor ChatGPT Project
- Related feature: FEATURE-0013
- Related phase: Phase 2, with cross-phase contract applicability
- Author: ChatGPT architecture synthesis with Sovrunn project owner
- Human approver: _PENDING_HUMAN_REVIEW_
- Approval status: Proposed

## Decision title

Establish the provider-neutral DecisionRecord, DecisionProfile, evaluation,
composition, rationale, obligation, AuditEvent-linkage, and downstream-adoption
standard for Sovrunn.

## Summary

This proposed handoff establishes the common, immutable, provider- and
runtime-neutral contract for governed decisions across authorization,
governance, security, placement, sovereignty, lifecycle, availability,
compatibility, economics, and AI-assisted use cases. It separates atomic
evaluation, governed decision, audit accountability, asynchronous operation,
and authorized explanation; supports typed decision profiles without changing
the common envelope; requires deterministic bounded composition and locally
available digest-pinned synchronous inputs; and preserves Phase 2 as
contract-only work through explicit anti-overengineering guardrails.

Canonical architecture:

`docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`

Controlling reuse standard:

`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

## Classification

- Extension

The extension includes two controlled corrections that require explicit human
approval:

1. replace the ambiguous common term `DecisionObject` with `DecisionRecord`;
2. extend the FEATURE-0012 alpha `AuditEvent` profile from Organization-only to
   the six formal governance scopes without treating scope as authorization.

## Existing approved baseline

The approved Phase 2 spine requires a provider-neutral, policy-aware,
auditable, explainable decision fabric without infrastructure execution. It
keeps `DecisionObject`, `AuditEvent`, and `Operation` distinct; treats structured
explainability as mandatory; keeps AI a consumer rather than authority; resolves
hierarchical governance into `EffectivePolicyContext`; and places external
engines and providers behind adapter boundaries.

Relevant baseline references:

- `docs/foundation/constitution.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md`
- `docs/phase2/PHASE2_EXECUTION_STRATEGY.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
- `docs/architecture/api-resource-standard.md`
- `docs/decisions/DEC-0026-reuse-before-build.md`
- `docs/decisions/DEC-0027-phase2-scope.md`
- `docs/decisions/DEC-0036-adapter-boundaries.md`
- `docs/rfc/RFC-0023-decision-and-audit-standard.md`
- `ADH-2026-012`
- `ADH-2026-013`

## Decision or proposed decision

Subject to recorded human approval of AD-001 through AD-044 and Matrix E risks
F13-R01 through F13-R31:

1. Use `DecisionRecord` as the immutable common governed-conclusion envelope;
   retain `EvaluationResult`, `AuditEvent`, `Operation`, and `DecisionContext`
   as distinct concepts.
2. Extend the envelope through governed, versioned `DecisionProfile` contracts
   and typed result schemas rather than modifying or forking the envelope.
3. Use the closed primary-form vocabulary: adjudication, selection, ranking,
   classification, resolution, allocation, plan, assessment, and
   recommendation.
4. Keep authority orthogonal to form: authoritative, advisory,
   recommendation, and simulation.
5. Limit `ALLOWED`, `DENIED`, and `REQUIRES_APPROVAL` to adjudication facets;
   keep evaluator failure and typed domain results separate.
6. Make every persisted decision final and immutable. Corrections,
   supersessions, and revocations are new linked final records; effective state
   is calculated by projection.
7. Require deterministic composition from captured normalized evaluation
   results. Replays do not rerun external, probabilistic, or AI evaluators.
8. Resolve hierarchy once into a versioned, time-bounded, digest-pinned
   `EffectivePolicyContext`; do not create hidden per-domain inheritance engines.
9. Permit bounded decision-requirements DAGs for dependency and composition,
   not as general workflow languages or production execution engines.
10. Keep the simple path synchronous and local. `Operation` owns long-running,
    waiting, retrying, cancellation, and human-approval lifecycle.
11. Orchestrate authoritative composition; permit choreography only for
    non-authoritative post-decision reactions.
12. Require interchangeable stateless workers, separate retry idempotency from
    semantic decision identity, and bound calls, bytes, concurrency, retries,
    time, jurisdiction, and cost.
13. Atomically accept the immutable decision and preallocated audit obligation
    at the same local persistence boundary; permit downstream audit
    materialization and delivery asynchronously.
14. Prohibit implicit universal synchronous dependencies on remote registries,
    policy resolvers, audit indexes, evidence exporters, AI models, identity
    control planes, cryptographic notaries, provider control planes, or vendor
    services.
15. Treat sovereignty as data, operational, technical, legal, cryptographic,
    supply-chain, and AI/model control, with connected through air-gapped
    conformance profiles.
16. Require least-privilege authorized projections, structured rationale,
    enforceable obligations, provenance, redaction, bounded evidence, and
    algorithm-agile integrity.
17. Treat authorization as a decision profile while deferring authentication,
    identity resolution, policy evaluation, enforcement, caching, and revocation
    implementation to their owning features.
18. Keep AI optional and non-authoritative; AI consumes only authorized
    `DecisionContext` projections and cannot mutate, approve, or execute.
19. Carry Matrix E v2 risk governance with stable IDs, explicit treatment,
    evidence, ownership, reassessment triggers, acceptance expiry, and
    human-only residual disposition.
20. Constrain FEATURE-0013 implementation to `CONTRACT_NOW` schemas, static
    formats, validation, compatibility, pure reference composition, fixtures,
    documentation, traceability, and gates.
21. Require one normative downstream adoption contract, one short adoption
    section per dependent feature, and one lightweight feature-gate check; do
    not create manifests, registries, indexes, or separate approval workflows
    until observed triggers justify them and a human approves the added cost.

The detailed AD-001 through AD-044 scorecard, Matrix E v2 register, conformance
model, Kiro guardrails, and downstream adoption contract in the canonical
architecture are normative parts of this proposed decision.

## Rationale

Sovrunn needs one reusable decision grammar that supports a simple allow/deny
and a bounded complex composite outcome without conflating evidence, authority,
accountability, or execution state. A stable envelope plus typed profiles
avoids repeated architecture changes while retaining domain correctness.
Locally available bounded inputs protect latency and horizontal scale. Explicit
authority and projections make AI useful without allowing it to bypass policy,
approval, or audit. Contract-only guardrails prevent Kiro from turning the
standard into a premature production platform.

## Reuse-before-build assessment

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - FEATURE-0012 API/resource grammar and conformance foundation
  - bounded decision-requirements concepts
  - standard event-envelope concepts for optional transport mapping
  - standard telemetry correlation concepts
  - canonical JSON and content-digest concepts
  - established software-supply-chain evidence concepts
  - policy-evaluator and durable-execution systems as later adapter/runtime
    candidates only
- Sovrunn-owned responsibility summary:
  - common decision/evaluation/audit/operation distinctions
  - `DecisionRecord`, `DecisionProfile`, forms, authority, typed-result,
    rationale, obligation, provenance, identity, correction, and linkage
    semantics
  - deterministic bounded composition and conformance
  - sovereign projection, portability, compatibility, and risk contracts
  - anti-overengineering and downstream-adoption governance
- Non-goals summary:
  - no production decision, authorization, policy, placement, audit, evidence,
    identity, registry, synchronization, cryptographic, provider, or AI service
  - no database, queue, cache, workflow engine, controller, worker, outbox,
    distributed coordination, or external integration
  - no runtime or provider product selection

## Phase impact

- Current phase allowed: Yes
- Phase 2 output:
  - normative schemas and static contract formats
  - validation and stable errors
  - compatibility and conformance fixtures/checks
  - pure reference composition only where required to prove semantics
  - architecture, risk, traceability, and feature-gate evidence
- Cross-phase effect:
  - later decision-producing and decision-consuming features adopt the standard
    through the lightweight downstream section unless an approved architecture
    change supersedes it
- Current phase boundary impact:
  - no infrastructure or lifecycle execution
  - no production runtime implementation
  - no change to the Phase 2 feature sequence

## Conflict check

- Conflicts with accepted DEC/RFC: Yes, two controlled baseline deltas require
  approval.
- Conflicting or incomplete baseline statements:
  - the Phase 2 spine uses `DecisionObject`; this proposal uses
    `DecisionRecord`.
  - the FEATURE-0012 alpha `AuditEvent` profile is Organization-only; this
    proposal requires six formal governance scopes.
- Resolution required:
  - approve this ADH;
  - update the architecture spine, current baseline, current decision summary,
    RFC, feature index, traceability, and relevant gates;
  - preserve compatibility and migration evidence.

## Required action

- Approve, reject, or defer AD-001 through AD-044.
- Complete per-risk Matrix E disposition for F13-R01 through F13-R31.
- Approve or reject the `DecisionRecord` terminology correction.
- Approve or reject the six-scope `AuditEvent` correction.
- Validate this handoff with `make arch-handoff-check`.
- After approval, update the canonical architecture status to
  `approved-for-kiro-requirements` and record this handoff as controlling.
- Update baseline, RFC, feature-index, traceability, agent-loading, and gate
  references.
- Initialize the Kiro spec path resolved by `docs/features/FEATURE_INDEX.md`.
- Generate `requirements.md` only and wait for `APPROVED_FOR_DESIGN`.

## Impacted files

- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-014-feature-0013-decision-record-and-auditevent-standard.md`
- `docs/rfc/RFC-0023-decision-and-audit-standard.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_DECISION_SUMMARY.md`
- `docs/features/FEATURE_INDEX.md`
- `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md`
- `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md`
- `docs/phase2/PHASE2_ACCEPTANCE_GATES.md`
- `.kiro/agents/sovrunn-spec.md`
- `.kiro/specs/decision-object-and-auditevent-standard/.config.kiro`
- `.kiro/specs/decision-object-and-auditevent-standard/requirements.md`

## Impacted features

- FEATURE-0012:
  - receives only the explicitly approved `AuditEvent` scope compatibility
    correction and related evidence; no runtime behavior.
- FEATURE-0013:
  - owns the common decision and audit-linkage contracts and conformance.
- Later evaluation, effective-context, placement, authorization, operation,
  projection, AI, and audit adopters:
  - reference the lightweight downstream adoption contract;
  - implement their own domain/runtime behavior only through separately
    approved features.

## Acceptance criteria for Kiro update

- [ ] Handoff validated against Architecture Operating System files.
- [ ] Approval status is Approved with human approver and date before Kiro uses
      it as controlling input.
- [ ] Canonical architecture status is `approved-for-kiro-requirements`.
- [ ] AD-001 through AD-044 are preserved by reference or testable requirement
      coverage without being silently changed.
- [ ] Matrix E F13-R01 through F13-R31 IDs, controls, verification, ownership,
      and human dispositions are preserved.
- [ ] `DecisionRecord`, `DecisionProfile`, `EvaluationResult`, `AuditEvent`,
      `Operation`, and `DecisionContext` remain distinct.
- [ ] Simple and bounded complex use cases are covered without authorizing a
      production graph/workflow engine.
- [ ] Provider/runtime neutrality, statelessness, local data-plane inputs,
      sovereignty, security, projection, AI, audit, identity, compatibility,
      and boundedness invariants are testable.
- [ ] Every requirement is classified `CONTRACT_NOW`,
      `INVARIANT_FOR_LATER`, or `DEFERRED`.
- [ ] `INVARIANT_FOR_LATER` and `DEFERRED` items generate no production runtime
      implementation tasks.
- [ ] Lightweight downstream adoption remains documentation plus one gate; no
      manifest/registry/index/approval framework is introduced.
- [ ] Requirements include explicit non-goals, reassessment triggers, and
      `ARCHITECTURE_DECISION_REQUIRED` behavior.
- [ ] Requirements do not choose implementation products or generate source
      code.

## Explicit instructions to Kiro

- Generate `requirements.md` only.
- Use the FEATURE-0013 slug resolved from `docs/features/FEATURE_INDEX.md`.
- Treat the approved canonical architecture and this approved ADH as
  controlling inputs.
- Do not introduce decisions beyond AD-001 through AD-044.
- Do not renumber, merge, delete, downgrade, transfer, or accept Matrix E risks.
- Do not invent owners, limits, SLOs, retention periods, legal values,
  reviewers, dates, approvals, or product selections.
- Record unresolved semantic gaps as `ARCHITECTURE_DECISION_REQUIRED`.
- Do not build prohibited section 27 runtime behavior.
- Do not create structured downstream manifests or a central adoption registry.
- Do not generate design, tasks, or implementation until separately approved.
- Do not modify Go source while applying this architecture handoff.

## Human approval

- Approval status: Proposed
- Approved by: _PENDING_HUMAN_REVIEW_
- Date: _PENDING_HUMAN_REVIEW_
- Notes:
  - This file stages the proposed handoff only.
  - It does not accept Matrix E residual risk.
  - It does not authorize Kiro requirements until separately human-approved.
