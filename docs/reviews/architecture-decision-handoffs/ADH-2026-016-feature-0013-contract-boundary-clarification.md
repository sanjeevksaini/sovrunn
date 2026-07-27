---
doc_type: architecture_decision_handoff
handoff_id: ADH-2026-016
feature: FEATURE-0013
title: Contract Boundary Clarification
status: Approved
classification: Clarification
phase: 2
approval_date: 2026-07-27
approving_role: Sovrunn Architecture Owner
---

# ADH-2026-016 — FEATURE-0013 Contract Boundary Clarification

## Metadata

- Handoff ID: ADH-2026-016
- Date: 2026-07-27
- Source discussion: Sovrunn Architecture Governor ChatGPT Project
- Related feature: FEATURE-0013
- Related phase: Phase 2
- Author: ChatGPT architecture synthesis with Sovrunn project owner
- Human approver: Sanjeev Kumar
- Reviewer: Sanjeev Kumar
- Approval status: Approved
- Approval date: 2026-07-27
- Predecessors: ADH-2026-014 and ADH-2026-015 remain controlling except where
  this handoff clarifies their unresolved contract boundaries.
- Requirements review context: the ADH-2026-015 review epoch exhausted its
  automated revision limit and stopped at a human gate.

## Decision title

Clarify the FEATURE-0013 contract boundaries for DecisionRecord scope
authority, provider-neutral sensitivity vocabulary, security validation scope,
canonicalization and trust boundary, exact conformance coverage, six-scope
evidence disposition, and DecisionObject migration.

## Summary

This handoff clarifies the unresolved contract boundaries of FEATURE-0013 that
remained open after ADH-2026-014 and ADH-2026-015. It confirms that
DecisionRecord uses FEATURE-0012 `metadata.scopeRef` as its sole logical and
serialized scope authority, defines a closed ordered provider-neutral
sensitivity vocabulary, fixes the security validation boundary at structural
and typed conformance rather than semantic content inspection, separates
CONTRACT_NOW carrier and structural trust metadata from DEFERRED cryptographic
selection, fixes the exact conformance coverage as scenarios F13-CF-01 through
F13-CF-28, sets the disposition of six-scope evidence to SUPERSEDED, and
completes the DecisionObject-to-DecisionRecord naming migration. This is a
clarification. It authorizes no runtime capability and no product, provider, or
algorithm selection. It returns requirements to architecture-remediation status
and does not resume or generate design, tasks, or implementation.

## Classification

- Clarification

This is a clarification of contract boundaries already implied by the approved
FEATURE-0013 baseline, ADH-2026-014, and ADH-2026-015. ADH-2026-014 and
ADH-2026-015 remain controlling except where this handoff clarifies their
unresolved contract boundaries. No new architecture is introduced.

## Existing approved baseline

The approved Phase 2 baseline establishes ScopeKind as the shared scope
vocabulary through FEATURE-0012, DecisionRecord and AuditEvent contract
standards through ADH-2026-014, and the seven-value canonical ScopeKind
vocabulary and `metadata.scopeRef` as sole canonical scope identity through
ADH-2026-015. FEATURE-0013 is contract-only work authorizing no runtime
capability.

Relevant baseline references:

- `docs/foundation/constitution.md`
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
- `docs/architecture/api-resource-standard.md`
- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- `docs/decisions/DEC-0026-reuse-before-build.md`
- `docs/decisions/DEC-0027-phase2-scope.md`
- `docs/traceability/ADH-2026-014-six-scope-auditevent-compatibility.md`
- `docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md`
- `ADH-2026-014`
- `ADH-2026-015`

## Decision or proposed decision

The Sovrunn Architecture Owner approves the following eight decision groups.

### 1. DecisionRecord scope authority

- DecisionRecord uses FEATURE-0012 `metadata.scopeRef` as its sole logical and
  serialized scope authority.
- A top-level DecisionRecord `scopeRef` or any second scope source is removed
  and prohibited.
- Platform uses absent/nil `metadata.scopeRef` where Platform is permitted.
- Every non-Platform scope uses non-nil `metadata.scopeRef` with canonical
  ScopeKind and UID semantics.
- Duplicate, conflicting, alternate, aliased, or parallel scope sources fail
  validation.
- The seven-value ScopeKind vocabulary approved by ADH-2026-015 is preserved.

### 2. Provider-neutral sensitivity vocabulary

- The closed ordered core values are:
  PUBLIC < INTERNAL < CONFIDENTIAL < RESTRICTED.
- Jurisdiction-specific classification labels and caveats are versioned profile
  mappings, not additional core enum values.
- Profiles permitting captured evaluator output must declare a sensitivity
  ceiling, allowed/prohibited content-category rules, maximum fields, maximum
  bytes, maximum depth, and allowed value types.
- Missing or unknown mandatory values invalidate the profile.
- Captured values above the declared ceiling fail validation.

### 3. Security validation boundary

- FEATURE-0013 owns structural schema validation, typed classification,
  declared bounds, projection restrictions, and deterministic conformance.
- FEATURE-0013 does not claim comprehensive semantic detection of secrets,
  credentials, personal data, malware, or policy violations in arbitrary
  content.
- Producers must redact prohibited content before submission.
- Runtime content inspection belongs to later approved security features.
- Kiro and Cursor must not invent secret-pattern, credential-pattern, PII,
  malware, DLP, or content-scanning engines in FEATURE-0013.

### 4. Canonicalization and trust boundary

- Algorithm-agile carrier fields and structural trust metadata are
  CONTRACT_NOW.
- Canonicalization algorithm/profile, digest-covered fields, signature
  algorithm, and cryptographic product/service selection are DEFERRED and
  remain blocked by ADR-F13-002.
- Actual cryptographic verification, signing, HSM, notary, trusted-time, key,
  revocation-distribution, and WORM services are not FEATURE-0013 artifacts.
- Structural fixtures may use opaque deterministic test assertions to validate
  state and failure semantics but must not claim cryptographic validity or
  select an algorithm.
- When a profile requires verified trust, an unknown, absent, expired, revoked,
  mismatched, or unverified trust state fails closed.
- RFC 8785 is illustrative only and cannot be selected by design, tasks,
  fixtures, or code without a later approved decision.

### 5. Exact conformance coverage

- Architecture section 17 scenarios map exactly and one-to-one to stable IDs
  F13-CF-01 through F13-CF-28 in their existing order.
- Additional AC-13.7, scope, security, and compatibility cases use separately
  named IDs and do not alter the canonical count.
- One fixture may cover multiple scenario IDs, but every scenario must have an
  explicit coverage-matrix entry; coverage is counted by scenario ID, not by
  fixture-file count.
- Contract-only fixtures must remain pure/in-memory and must not implement
  prohibited runtime systems.

### 6. Six-scope evidence disposition

- `docs/traceability/ADH-2026-014-six-scope-auditevent-compatibility.md`
  becomes SUPERSEDED historical evidence.
- ADH-2026-015 and
  `docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md` are
  authoritative for the seven-value contract.
- Historical content is preserved but must carry an explicit supersession
  notice and current disposition.

### 7. DecisionObject migration

- Active normative documents use DecisionRecord.
- Frozen FEATURE-0012 Kiro specifications may retain DecisionObject only with an
  explicit artifact-level note identifying it as the historical
  pre-ADH-2026-014 name for the same concept.
- No second schema, Go type, runtime alias, parallel contract, or migration
  implementation is created.
- Generated context files are regenerated from corrected sources and are not
  edited manually.
- Reconciliation evidence must distinguish corrected active normative documents
  from frozen historical artifacts.

### 8. Scope and non-goals

- No production decision/audit/security/cryptographic/profile-registry/
  synchronization service is authorized.
- No provider-specific component, parallel scope field, second decision type,
  adapter interface, queue, worker, persistence layer, or product selection is
  authorized.
- Matrix E identifiers and all unaffected ADH-2026-014/ADH-2026-015 decisions
  are preserved.

## Rationale

The ADH-2026-015 requirements review epoch exhausted its automated revision
limit and stopped at a human gate because several FEATURE-0013 contract
boundaries remained ambiguous: whether DecisionRecord carried a second scope
source, how provider-neutral sensitivity relates to jurisdiction-specific
labels, how far security validation extends into content, which trust
mechanics are contract-now versus deferred, how conformance scenarios are
counted, the disposition of six-scope evidence, and the residual
DecisionObject naming. Clarifying these boundaries without introducing new
architecture allows a fresh requirements-review epoch to proceed
deterministically. Fixing `metadata.scopeRef` as the sole scope authority
prevents divergent scope sources. A closed ordered sensitivity vocabulary with
versioned profiles keeps the contract provider-neutral while allowing
jurisdiction mappings without enum proliferation. Bounding security validation
to structural conformance prevents Kiro or Cursor from inventing unapproved
scanning engines. Separating CONTRACT_NOW carrier fields from DEFERRED
cryptographic selection preserves algorithm agility while keeping ADR-F13-002
blocking. Fixing the F13-CF-01 through F13-CF-28 count and counting coverage by
scenario ID makes conformance auditable. Marking six-scope evidence SUPERSEDED
while preserving it maintains an accurate historical record. Completing the
DecisionObject-to-DecisionRecord migration removes naming ambiguity without
creating a parallel schema.

## Reuse-before-build assessment

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - FEATURE-0012 ScopeKind vocabulary and `metadata.scopeRef` scope authority
  - FEATURE-0013 canonical architecture, ADH-2026-014 and ADH-2026-015
    decisions
  - RFC 8785 referenced only illustratively, not selected
- Sovrunn-owned responsibility summary:
  - contract-boundary clarification for scope authority, sensitivity
    vocabulary, security validation boundary, trust boundary, conformance
    coverage, evidence disposition, and naming migration
  - synchronization of architecture, requirements, active normative documents,
    superseded evidence, and historical compatibility notes
- Non-goals summary:
  - no runtime capability of any kind
  - no cryptographic product, service, or algorithm selection
  - no provider-specific component or parallel scope field
  - no second decision type, schema, Go type, adapter, queue, worker, or
    persistence layer

## Phase impact

- Current phase allowed: Yes
- Phase 2 output:
  - clarified DecisionRecord scope authority using `metadata.scopeRef`
  - closed ordered provider-neutral sensitivity vocabulary with versioned
    profiles
  - bounded security validation scope
  - CONTRACT_NOW versus DEFERRED trust boundary with ADR-F13-002 preserved
  - exact conformance coverage F13-CF-01 through F13-CF-28
  - superseded six-scope evidence disposition
  - completed DecisionObject-to-DecisionRecord migration in active normative
    documents
- Cross-phase effect:
  - later approved security features own runtime content inspection
  - later approved decisions own cryptographic algorithm and product selection
- Current phase boundary impact:
  - no infrastructure or lifecycle execution
  - no production runtime implementation
  - no change to the Phase 2 feature sequence

## Conflict check

- Conflicts with accepted DEC/RFC: No new decision conflict; this handoff
  clarifies unresolved contract boundaries within the existing baseline.
- Incomplete or ambiguous baseline statements clarified:
  - residual second scope source ambiguity for DecisionRecord
  - relationship between provider-neutral sensitivity and jurisdiction labels
  - extent of security validation into arbitrary content
  - contract-now versus deferred trust mechanics
  - conformance scenario counting method
  - disposition of six-scope AuditEvent compatibility evidence
  - residual DecisionObject naming in active documents
- Resolution required:
  - apply this approved ADH;
  - synchronize architecture, requirements, active normative documents,
    superseded evidence, and historical compatibility notes;
  - preserve Matrix E identifiers and all unaffected ADH-2026-014/ADH-2026-015
    decisions;
  - keep ADR-F13-002 blocking on deferred cryptographic selection.

## Required action

- Return requirements to architecture-remediation status.
- Synchronize architecture, requirements, active normative documents,
  superseded evidence, and historical compatibility notes.
- Start a fresh independently governed requirements-review epoch only after
  synchronization and validation.
- Do not generate or resume design, tasks, or implementation.
- The superseded design remains non-authoritative.
- Mark `docs/traceability/ADH-2026-014-six-scope-auditevent-compatibility.md`
  as SUPERSEDED historical evidence with an explicit supersession notice.
- Complete the DecisionObject-to-DecisionRecord migration in active normative
  documents; regenerate generated context files from corrected sources.
- Validate this handoff with `make arch-handoff-check`.

## Impacted files

- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-016-feature-0013-contract-boundary-clarification.md`
- `docs/traceability/ADH-2026-014-six-scope-auditevent-compatibility.md`
- `docs/traceability/ADH-2026-015-scope-vocabulary-compatibility.md`
- `.kiro/specs/decision-object-and-auditevent-standard/requirements.md`
- `.kiro/specs/decision-object-and-auditevent-standard/design.md`
- Canonical DecisionRecord and sensitivity schema files
- Go DecisionRecord bindings, sensitivity validators, and profile validators
- DecisionRecord, sensitivity, trust, and conformance test fixtures
- Generated context files regenerated from corrected sources

## Impacted features

- FEATURE-0012:
  - `metadata.scopeRef` remains the sole scope authority; frozen Kiro
    specifications may retain DecisionObject only with an explicit historical
    artifact note; no runtime behavior change.
- FEATURE-0013:
  - owns the clarified contract boundaries for scope authority, sensitivity
    vocabulary, security validation, trust, conformance coverage, evidence
    disposition, and naming migration; architecture, requirements, active
    normative documents, superseded evidence, and historical compatibility
    notes must be synchronized; requirements return to
    architecture-remediation status.

## Acceptance criteria for Kiro update

- [ ] DecisionRecord uses FEATURE-0012 `metadata.scopeRef` as its sole logical
      and serialized scope authority; no top-level `scopeRef` or second scope
      source exists.
- [ ] Platform uses absent/nil `metadata.scopeRef`; every non-Platform scope
      uses non-nil `metadata.scopeRef` with canonical ScopeKind and UID
      semantics.
- [ ] Duplicate, conflicting, alternate, aliased, or parallel scope sources
      fail validation; the seven-value ScopeKind vocabulary is preserved.
- [ ] The closed ordered sensitivity core is PUBLIC < INTERNAL < CONFIDENTIAL <
      RESTRICTED; jurisdiction labels and caveats are versioned profile
      mappings, not core enum values.
- [ ] Profiles permitting captured evaluator output declare a sensitivity
      ceiling, allowed/prohibited content-category rules, maximum fields,
      maximum bytes, maximum depth, and allowed value types.
- [ ] Missing or unknown mandatory profile values invalidate the profile;
      captured values above the declared ceiling fail validation.
- [ ] Security validation is bounded to structural schema validation, typed
      classification, declared bounds, projection restrictions, and
      deterministic conformance; no secret-pattern, credential-pattern, PII,
      malware, DLP, or content-scanning engine is introduced.
- [ ] Algorithm-agile carrier fields and structural trust metadata are
      CONTRACT_NOW; canonicalization algorithm/profile, digest-covered fields,
      signature algorithm, and cryptographic product/service selection remain
      DEFERRED and blocked by ADR-F13-002.
- [ ] Structural fixtures use opaque deterministic assertions and do not claim
      cryptographic validity or select an algorithm; RFC 8785 remains
      illustrative only.
- [ ] When a profile requires verified trust, an unknown, absent, expired,
      revoked, mismatched, or unverified trust state fails closed.
- [ ] Architecture section 17 scenarios map one-to-one to stable IDs
      F13-CF-01 through F13-CF-28 in existing order; additional AC-13.7, scope,
      security, and compatibility cases use separately named IDs.
- [ ] Every scenario has an explicit coverage-matrix entry; coverage is counted
      by scenario ID, not by fixture-file count; contract-only fixtures remain
      pure/in-memory.
- [ ] `docs/traceability/ADH-2026-014-six-scope-auditevent-compatibility.md`
      is marked SUPERSEDED historical evidence with an explicit supersession
      notice and current disposition; ADH-2026-015 and its compatibility note
      are authoritative for the seven-value contract.
- [ ] Active normative documents use DecisionRecord; frozen FEATURE-0012 Kiro
      specifications retain DecisionObject only with an explicit historical
      artifact note; no second schema, Go type, runtime alias, parallel
      contract, or migration implementation exists.
- [ ] Generated context files are regenerated from corrected sources and not
      edited manually; reconciliation evidence distinguishes corrected active
      normative documents from frozen historical artifacts.
- [ ] No production decision/audit/security/cryptographic/profile-registry/
      synchronization service, provider-specific component, parallel scope
      field, second decision type, adapter, queue, worker, persistence layer,
      or product selection is authorized.
- [ ] Matrix E identifiers and all unaffected ADH-2026-014/ADH-2026-015
      decisions are preserved.
- [ ] Requirements return to architecture-remediation status; design, tasks,
      and implementation are not generated or resumed; the superseded design
      remains non-authoritative.

## Explicit instructions to Kiro

- Apply this approved handoff as a clarification of FEATURE-0013 contract
  boundaries; ADH-2026-014 and ADH-2026-015 remain controlling except where
  this handoff clarifies their unresolved contract boundaries.
- Return requirements to architecture-remediation status.
- Synchronize architecture, requirements, active normative documents,
  superseded evidence, and historical compatibility notes.
- Start a fresh independently governed requirements-review epoch only after
  synchronization and validation.
- Do not generate or resume design, tasks, or implementation; the superseded
  design remains non-authoritative.
- Fix `metadata.scopeRef` as the sole DecisionRecord scope authority and remove
  any second scope source.
- Keep the sensitivity core closed and ordered; treat jurisdiction labels as
  versioned profile mappings only.
- Bound security validation to structural conformance; do not invent
  secret-pattern, credential-pattern, PII, malware, DLP, or content-scanning
  engines.
- Keep cryptographic canonicalization, digest coverage, signature algorithm,
  and product/service selection DEFERRED and blocked by ADR-F13-002; keep
  RFC 8785 illustrative only.
- Keep conformance scenarios F13-CF-01 through F13-CF-28 exact and counted by
  scenario ID.
- Mark six-scope AuditEvent compatibility evidence SUPERSEDED while preserving
  its content.
- Complete the DecisionObject-to-DecisionRecord migration in active documents
  without creating any parallel schema, Go type, alias, or migration
  implementation; regenerate generated context files from corrected sources.
- Preserve Matrix E identifiers and all unaffected ADH-2026-014/ADH-2026-015
  decisions.
- Do not introduce runtime capabilities of any kind and do not select any
  product, provider, or algorithm.
- Do not modify review JSON, automation state, approval tokens, or Go source
  files outside the explicit scope of this handoff.

## Human approval

- Approval status: Approved
- Status: Approved
- Approved by: Sanjeev Kumar
- Reviewer: Sanjeev Kumar
- Decision date: 2026-07-27
- Date: 2026-07-27
- Notes:
  - Approval covers all eight decision groups above: DecisionRecord scope
    authority; provider-neutral sensitivity vocabulary; security validation
    boundary; canonicalization and trust boundary; exact conformance coverage;
    six-scope evidence disposition; DecisionObject migration; and scope and
    non-goals.
  - ADH-2026-014 and ADH-2026-015 remain controlling except where this handoff
    clarifies their unresolved contract boundaries.
  - The ADH-2026-015 review epoch exhausted its automated revision limit and
    stopped at a human gate; this clarification enables a fresh independently
    governed requirements-review epoch.
  - Requirements return to architecture-remediation status; design, tasks, and
    implementation are not generated or resumed; the superseded design remains
    non-authoritative.
  - This approval authorizes contract-only clarification: no persistence,
    network, queue, worker, workflow, cryptographic service, external
    integration, or other runtime capability, and no product, provider, or
    algorithm selection.
