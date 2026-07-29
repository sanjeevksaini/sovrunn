---
doc_type: architecture_decision_handoff
handoff_id: ADH-2026-017
feature: FEATURE-0013
title: Consolidated Decision Record and AuditEvent Architecture
status: Approved
classification: Replacement
phase: 2
approval_date: 2026-07-28
approving_role: Sovrunn Architecture Owner
supersedes_on_approval:
  - ADH-2026-014
  - ADH-2026-015
  - ADH-2026-016
canonical_architecture: docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md
architecture_content_sha256: b52e83316657fdc8be31db0be954f7a31d72752c6c3584ba1e0f98f40ffdfb74
dependency_content_sha256:
  "docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md": f9bd5ad9cf0aa10b06175204ed388637e64fe740fadc121d30983e00f0786623
  "docs/features/FEATURE-0011-reuse-assessment-standard.md": ad1f551a9f6f9710dd1e2a9e0ea0a0189e6240716c3fe016a44aa4f32bc34454
  "docs/reviews/feature-gates/FEATURE-0011-approval-review.md": 934ac514cbe54dec22ae6d293fda3dddafb2e6af37eddf262e79b33fa06d7393
  "docs/reviews/architecture-decision-handoffs/ADH-2026-011-feature-0011-reuse-assessment-standard.md": 78c28407bba69116e1c16b2a33ca5ae04ab55db5625425deca5ec3f23ed8946e
  "docs/architecture/api-resource-standard.md": b39e353324ba41c20dda07bd1c5e77bda2bbd103f2c6392f1ec1e5743cb505df
  "docs/features/FEATURE-0012-api-resource-naming-status-and-validation-standard.md": 47c47462b81e75c2ab63fbcfba49001f5611d081b838ed881f666da20fc61b0d
  ".kiro/specs/api-resource-naming-status-and-validation-standard/requirements.md": b9e068f0d782f38f028ad3daa835d58c6ad1824bb6e70a33bb4d363b09a95d06
  ".kiro/specs/api-resource-naming-status-and-validation-standard/design.md": 2e0cead4e8a82735c10bc79ba6133377b0ad6c51528f80b52fe7e03cabc50b0e
  "docs/reviews/feature-gates/FEATURE-0012-approval-review.md": 6d649372550eed3c07fe5080d542073ea9a2e66b9f863d6cad64d7a81b8f2169
  "docs/reviews/architecture-decision-handoffs/ADH-2026-012-feature-0012-api-resource-standard.md": a1613acb79974efa5958fe5f5daef24a0c488c30fb9e6a39a42b8f2ae25c98d2
  "docs/reviews/architecture-decision-handoffs/ADH-2026-013-operation-allowed-scopes.md": 0eef8061b631681111bea5a81a4da36281892e23bd743a3d36a0d3a9c8a40f8c
  "api/schemas/_common/scope-ref.json": 53abad345e304970afb8bddb157106fe54c107849b31135fd53def50b1177f7f
  "api/schemas/audit-event.json": c60dd8e0f9b5e5388ba1982ae316d619ea3ea23fb560b0b933d6b998f94a5ef0
  "internal/apimeta/scope.go": f8b0eaaf622ff4299ad520b249f8785bb73df67396ee46dc354db07ae67ab63f
---

# ADH-2026-017 — FEATURE-0013 Consolidated Architecture

## Metadata

- Handoff ID: ADH-2026-017
- Date: 2026-07-28
- Source discussion: Sovrunn Architecture Governor ChatGPT Project
- Related feature: FEATURE-0013
- Related phase: Phase 2, with cross-phase contract applicability
- Author: ChatGPT architecture synthesis with Sovrunn project owner
- Human approver: Sanjeev Kumar
- Approval status: Approved
- Classification: Replacement
- Canonical architecture:
  `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- Historical predecessors: ADH-2026-014, ADH-2026-015, ADH-2026-016

## Decision title

Adopt one complete, internally consistent FEATURE-0013 architecture and one
controlling architecture-decision handoff.

## Summary

This approved handoff replaces the fragmented normative effect of
ADH-2026-014, ADH-2026-015, and ADH-2026-016 with one approval envelope for the
complete FEATURE-0013 architecture. On approval, every normative decision in
sections 1 through 28, every AD-001 through AD-045 decision as reconciled by
the main body, every F13-R01 through F13-R31 architecture-stage treatment, the
closed downstream-stage boundary, and the downstream-adoption contract are
approved together. The canonical architecture main body is the decision
content; this handoff records its approval and exact content identity rather
than duplicating a second copy that could drift.

ADH-2026-014, ADH-2026-015, and ADH-2026-016 remain immutable historical
provenance after replacement. They are not controlling inputs for Kiro,
reviewers, design, tasks, or implementation after this handoff is approved.

## Classification

- Replacement

This is a governance consolidation and semantic correction. It replaces three
partially superseding FEATURE-0013 handoffs with one authoritative handoff. It
does not erase their history and does not authorize any runtime capability.

## Existing approved baseline

FEATURE-0011 owns the mandatory reuse-assessment governance contract.
FEATURE-0012 owns the shared API/resource grammar, including resource profiles,
type and common metadata, the six governance `ScopeKind` values, typed
references, boundaries, ownership, validation layering, Problem Details, and
compatibility rules. Earlier FEATURE-0013 handoffs established the decision and
audit direction but distributed active semantics across three records.

Relevant baseline references:

- `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
- `docs/architecture/api-resource-standard.md`
- `.kiro/specs/api-resource-naming-status-and-validation-standard/requirements.md`
- `.kiro/specs/api-resource-naming-status-and-validation-standard/design.md`
- `docs/reviews/feature-gates/FEATURE-0012-approval-review.md`
- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- ADH-2026-014, ADH-2026-015, and ADH-2026-016 as historical provenance

## Decision or proposed decision

The Sovrunn Architecture Owner is asked to approve the following as one
indivisible FEATURE-0013 architecture decision:

1. The exact approved content of
   `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
   is the sole normative FEATURE-0013 architecture.
2. The architecture content approved by this handoff includes sections 1
   through 28, the complete decision register, Matrix E v2 treatments,
   acceptance criteria, deferred decisions, reassessment triggers,
   anti-overengineering guardrails, and downstream-adoption contract.
3. AD-001 through AD-045 are an inventory and scorecard of decisions in the
   main body; if abbreviated scorecard wording differs from the main body, the
   main body controls.
4. `DecisionRecord`, `EvaluationResult`, `AuditEvent`, `Operation`, and
   `DecisionContext` remain distinct and retain the meanings and ownership
   assigned by the canonical architecture.
5. FEATURE-0013 inherits FEATURE-0012 resource profiles, metadata, references,
   boundaries, validation layers, Problem Details, schema conventions, and
   compatibility policy. It does not create parallel shared grammar.
6. The shared `ScopeKind` vocabulary remains the six FEATURE-0012 governance
   scopes: Platform, Organization, OrganizationUnit, Tenant, Project, and
   Provider.
7. `ServiceInstance` is not a `ScopeKind`. A decision or audit event concerning
   a ServiceInstance uses its Project governance scope in
   `metadata.scopeRef` and identifies the ServiceInstance through a typed
   `subjectRef` or `subjectRefs`. `ownerRef` is used only for genuine lifecycle
   containment and does not replace scope or subject identity.
8. `metadata.scopeRef` is the sole governance-scope authority. The serialized
   Platform form is absence of `scopeRef`; explicit `null` is not a second
   canonical serialized representation.
9. FEATURE-0013 owns decision/audit domain semantics, declarative profile and
   semantic-registry contracts, deterministic bounded composition,
   conformance, and compatibility. It owns no adapter interface or production
   runtime service.
10. FEATURE-0011 reuse governance applies directly. The canonical architecture
    must contain exactly one conforming feature-level reuse summary before
    approval.
11. FEATURE-0012 platform/schema validation ceilings are absolute outer bounds;
    a DecisionProfile may only narrow them, and an evaluator registration may
    only narrow the profile. The effective limit is the most restrictive
    applicable value.
12. Public failures use the FEATURE-0012 RFC 9457 Problem Details envelope.
    Architecture section 15 binds FEATURE-0013 validation to the inherited
    validation-failed type, HTTP 422, top-level `VALIDATION_FAILED`, RFC 6901
    pointers, and an architecture-owned closed `violations[].code` registry.
    Downstream stages may map cases and validators but may not invent public
    codes or error semantics.
13. `DecisionRecord` and `AuditEvent` conform to the FEATURE-0012
    `ImmutableRecord` shape: type metadata, valid common metadata, and one
    `record` payload. Examples are normative for shape and must validate.
14. Immutable records are never rewritten for correction, revocation, or
    erasure. Linked records, separately controlled payload custody, or later
    approved cryptographic erasure may preserve a tombstone without mutating
    the canonical record.
15. FEATURE-0013 sensitivity is a separate handling dimension mapped from, and
    never weakening or replacing, FEATURE-0012 `DataClassification`.
16. Structural trust carriers are CONTRACT_NOW; cryptographic algorithm,
    covered-field, signing, verification, product, and service selections are
    deferred under ADR-F13-002.
17. Requirements translate architecture into testable obligations; design may
    decide only delegated representation and implementation mechanics; tasks
    implement only approved CONTRACT_NOW work. Missing semantics require
    `ARCHITECTURE_DECISION_REQUIRED`.
18. Phase 2 remains contract-only. No persistence, queue, worker, workflow,
    production registry, provider integration, cryptographic service, policy
    engine, AI runtime, or other production decision/audit service is approved.
19. This approval supersedes ADH-2026-014, ADH-2026-015, and
    ADH-2026-016 in full as normative inputs. Those files remain historical
    evidence only.
20. Any later FEATURE-0013 semantic change requires a new approved handoff; it
    must not reactivate or reinterpret an earlier superseded handoff.
21. Final approval changes each accepted reuse-summary decision from `Proposed`
    to `Approved` and records the human decision; deferred rows remain
    `Deferred`. The approved architecture must contain no `Proposed` reuse row.
22. `DecisionRequest` is a calling-domain-owned conceptual input role. It
    inherits FEATURE-0012 `TransientRequestResult`, reference, scope,
    validation, limit, and error boundaries, but FEATURE-0013 defines no shared
    request schema, Go type, registry entry, resource, or persistence contract.
23. Completed handling returns `DecisionRecord`; accepted asynchronous handling
    returns the existing FEATURE-0012 `Operation`; rejected or unaccepted
    handling returns inherited Problem Details. FEATURE-0013 defines no
    pending-decision response envelope or status vocabulary.
24. FEATURE-0013 reuses FEATURE-0012 risk-governance principles but not its risk
    identifiers. `F12-R01` through `F12-R16` remain owned exclusively by
    FEATURE-0012; FEATURE-0013 owns the separate `F13-R01` through `F13-R31`
    namespace. The namespaces are never merged, renumbered, aliased, inherited,
    or mapped by ordinal position.
25. The section 19 implementation-class assignment for every AD-001 through
    AD-045 is architecture-owned and immutable downstream. Requirements,
    design, and tasks must reproduce each exact `CONTRACT_NOW`,
    `INVARIANT_FOR_LATER`, or `DEFERRED` value. Only `CONTRACT_NOW` may produce
    FEATURE-0013 implementation artifacts; persistent graph storage,
    synchronization runtime, cryptographic execution, and any unspecified
    non-final response remain unauthorized.

The approved `architecture_content_sha256` is the SHA-256 of the exact
canonical architecture bytes approved on 2026-07-28. A hash mismatch
invalidates downstream authorization until a new human decision is recorded.
The `dependency_content_sha256` map binds the exact approved FEATURE-0011 and
FEATURE-0012 standards, approvals, architecture/specification artifacts, scope
schema, AuditEvent schema, and ScopeKind binding used by the consolidation. A
dependency mismatch also invalidates downstream authorization until the delta
is reconciled and a new human architecture decision updates the lock.

## Rationale

One canonical architecture plus one content-bound approval handoff eliminates
the ambiguity created by layered corrections. It prevents Kiro or a reviewer
from choosing among historical statements, prevents partial supersession from
becoming a hidden precedence system, and makes architectural change visible.
Keeping historical handoffs preserves governance evidence without treating
history as active instruction.

Separating Project governance scope from ServiceInstance subject identity
preserves FEATURE-0012 semantics and avoids promoting every resource type into
the scope hierarchy. Binding the approval to exact architecture content avoids
duplicating the full architecture in two files while still approving every
decision as one package.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - FEATURE-0011 Reuse Assessment Standard
  - FEATURE-0012 API/resource grammar and conformance foundation
  - ADR/RFC architecture approval and immutable provenance practices
- Sovrunn-owned responsibility summary:
  - consolidated decision/audit domain architecture
  - content-bound single-handoff approval
  - FEATURE-0013 semantic registries, conformance, and downstream adoption
- Non-goals summary:
  - no deletion of historical approvals
  - no parallel architecture copy in the handoff
  - no runtime implementation, vendor selection, or stage authorization

## Phase impact

- Current phase allowed: Yes
- Current work: architecture consolidation, dependency reconciliation,
  approval controls, and later requirements regeneration
- Current phase boundary impact: contract-only; no runtime capability
- Cross-phase effect: later adopting features reference the single canonical
  architecture version and ADH-2026-017 only

## Conflict check

- Conflicts with accepted DEC/RFC: No unresolved conflict is permitted at
  approval time.
- Historical conflicts replaced by this handoff:
  - fragmented ADH-2026-014/015/016 precedence;
  - ServiceInstance-as-ScopeKind wording;
  - any example or rule inconsistent with FEATURE-0011/0012 inheritance.
- Resolution required:
  - finish canonical architecture corrections;
  - validate FEATURE-0011 and FEATURE-0012 reconciliation;
  - record exact architecture SHA-256;
  - obtain one fresh human approval.

## Required action

- Update and validate the canonical FEATURE-0013 architecture.
- Preserve this Approved handoff and its exact architecture and dependency
  digest bindings; any content change requires a new human architecture decision.
- Exclude ADH-2026-014/015/016 from downstream Kiro and reviewer context.
- Treat ADH-2026-014/015/016 as Superseded historical records.
- Update baseline, phase context, decision summary, RFC, feature index,
  traceability, and generated context to reference ADH-2026-017 only.
- Regenerate requirements.md from the approved architecture; do not reuse the
  pre-consolidation requirements or design.
- Keep design, tasks, and implementation unauthorized until their normal gates.

## Impacted files

- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-017-feature-0013-consolidated-architecture.md`
- ADH-2026-014, ADH-2026-015, ADH-2026-016 metadata/history notices
- `.kiro/agents/sovrunn-spec.md`
- `docs/prompts/kiro/requirements.prompt.md`
- `docs/prompts/kiro/design.prompt.md`
- `docs/prompts/kiro/tasks.prompt.md`
- `scripts/feature-0013-architecture-boundary-check.py`
- `scripts/reviewer-stage.sh`
- active context, phase, RFC, feature-index, and traceability documents

## Impacted features

- FEATURE-0011: reuse-assessment format remains controlling and must be
  satisfied by the FEATURE-0013 architecture.
- FEATURE-0012: shared grammar remains controlling; ServiceInstance remains a
  Project-scoped subject rather than a new governance scope.
- FEATURE-0013: receives one architecture and one controlling handoff.
- Later adopters: reference ADH-2026-017 and the exact approved architecture
  version only.

## Acceptance criteria for Kiro update

- [ ] ADH-2026-017 is the only active controlling FEATURE-0013 handoff.
- [ ] The architecture content SHA-256 matches the approved file exactly.
- [ ] Every `dependency_content_sha256` entry matches its exact repository file,
      and the locked path set contains no missing or additional dependency.
- [ ] ADH-2026-014/015/016 are excluded from Kiro and reviewer context.
- [ ] The canonical architecture contains one conforming FEATURE-0011 reuse
      summary.
- [ ] Every accepted reuse-summary row is `Approved` and every intentionally
      deferred row is `Deferred`; no `Proposed` reuse decision remains.
- [ ] FEATURE-0012 profile, metadata, scope, reference, validation, error,
      limit, schema, and compatibility semantics are inherited explicitly.
- [ ] The public error binding and closed violation-code registry are owned by
      architecture; downstream stages cannot create or reinterpret codes.
- [ ] ServiceInstance is represented through typed subject identity, not
      ScopeKind.
- [ ] DecisionRecord and AuditEvent examples conform to ImmutableRecord and
      ObjectMeta.
- [ ] No common FEATURE-0013 `DecisionRequest` schema/type is generated; any
      calling-domain request remains a conforming `TransientRequestResult`.
- [ ] The outcome boundary uses only `DecisionRecord`, `Operation`, and inherited
      Problem Details; no pending-decision response envelope is introduced.
- [ ] The generic `Operation` contract and payload remain owned by FEATURE-0012;
      FEATURE-0013 reuses but does not redefine them.
- [ ] FEATURE-0012 `F12-Rxx` and FEATURE-0013 `F13-Rxx` identifiers remain
      disjoint; only the risk-governance approach is reused.
- [ ] Every AD-001–AD-045 traceability entry reproduces exactly the
      architecture-owned implementation class, and only `CONTRACT_NOW` entries
      can produce FEATURE-0013 implementation artifacts.
- [ ] No downstream artifact introduces a typed non-final response, persistent
      graph store, synchronization runtime, signing/verification execution, or
      other production mechanism prohibited by section 27.
- [ ] Every downstream artifact contains architecture traceability and no
      unresolved semantic choice.
- [ ] Downstream traceability enumerates every architecture section 1–28,
      AD-001–AD-045, F13-R01–F13-R31, and every stable conformance ID
      individually; range-only assertions do not pass.
- [ ] Superseded requirements and design are not used as inputs.
- [ ] Design, tasks, and implementation remain gated separately.

## Explicit instructions to Kiro

- Load only the exact approved canonical architecture and ADH-2026-017 as the
  FEATURE-0013 architecture decision.
- Do not load ADH-2026-014, ADH-2026-015, or ADH-2026-016 as instructions.
- Treat older handoffs, requirements, design, reviews, and compatibility files
  only as historical evidence when a human explicitly requests provenance.
- Do not invent a missing semantic, mapping, alias, scope, error convention,
  limit, algorithm, runtime, or product.
- Stop with `ARCHITECTURE_DECISION_REQUIRED` on any conflict or missing
  architecture decision.
- Generate only the requested stage artifact.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-07-28
- Notes:
  - Approval covers the exact canonical architecture and SHA-256 recorded above.
  - Approval supersedes ADH-2026-014, ADH-2026-015, and ADH-2026-016 as
    normative inputs without deleting their historical record.
  - This handoff authorizes requirements regeneration only. Design, tasks, and
    implementation remain subject to their separate approval gates.
