---
doc_type: architecture_decision_handoff
handoff_id: ADH-2026-015
feature: FEATURE-0013
title: Scope Vocabulary Clarification
status: Approved
classification: Extension
phase: 2
approval_date: 2026-07-27
approving_role: Sovrunn Architecture Owner
---

# ADH-2026-015 — FEATURE-0013 Scope Vocabulary Clarification

## Metadata

- Handoff ID: ADH-2026-015
- Date: 2026-07-27
- Source discussion: Sovrunn Architecture Governor ChatGPT Project
- Related feature: FEATURE-0013
- Related phase: Phase 2
- Author: ChatGPT architecture synthesis with Sovrunn project owner
- Human approver: Sanjeev Kumar
- Approval status: Approved
- Predecessor: ADH-2026-014 remains controlling except that its six-scope
  AuditEvent statement is clarified and superseded by this handoff.

## Decision title

Clarify ScopeKind as the single canonical shared vocabulary and establish
ServiceInstance as the seventh canonical scope for AuditEvent.

## Summary

This handoff clarifies that ScopeKind is the single canonical shared vocabulary
for scope identity across Sovrunn, establishes the seven canonical values
(Platform, Organization, OrganizationUnit, Tenant, Project, Provider, and
ServiceInstance), and confirms that AuditEvent uses `metadata.scopeRef` as its
sole canonical scope identity drawn from this vocabulary. It supersedes the
six-scope AuditEvent statement in ADH-2026-014 by reconciling that AuditEvent
list with the FEATURE-0012 shared vocabulary: ServiceInstance is added to the
shared ScopeKind vocabulary, Provider is retained, and all seven values are
permitted for AuditEvent. No separate AuditEvent scope
enum or independently mutable scope field is permitted. This is contract-only
work authorizing no runtime capability.

## Classification

- Extension

This is an additive clarification of the approved ScopeKind vocabulary and
AuditEvent scope identity. ADH-2026-014 remains controlling for all other
decisions; only its six-scope AuditEvent statement is superseded.

## Existing approved baseline

The approved Phase 2 baseline establishes ScopeKind as the shared scope
vocabulary through FEATURE-0012, with values Platform, Organization,
OrganizationUnit, Tenant, Project, and Provider. ADH-2026-014 extended
AuditEvent from Organization-only to six formal governance scopes.

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
- `ADH-2026-014`

## Decision or proposed decision

The Sovrunn Architecture Owner approves the following:

1. ScopeKind is the single canonical shared vocabulary.
2. The canonical values are Platform, Organization, OrganizationUnit, Tenant,
   Project, Provider, and ServiceInstance.
3. ServiceInstance is added additively; Provider is retained unchanged.
4. AuditEvent uses `metadata.scopeRef` as its sole canonical scope identity.
5. AuditEvent initially permits all seven canonical scopes.
6. No separate AuditEvent scope enum or independently mutable scope field is
   permitted.
7. Existing FEATURE-0012 serialized ScopeKind values remain valid.
8. Existing Organization-scoped AuditEvents remain valid.
9. Canonical schemas, Go bindings, validators, fixtures, compatibility
   evidence, architecture, requirements, and design must be synchronized.
10. This is contract-only work and authorizes no persistence, network, queue,
    worker, workflow, cryptographic service, external integration, or other
    runtime capability.

## Rationale

ADH-2026-014 required ServiceInstance for AuditEvent even though
ServiceInstance was absent from FEATURE-0012's shared ScopeKind vocabulary,
while FEATURE-0012 included Provider and ADH-2026-014 omitted Provider from the
AuditEvent list. Adding ServiceInstance to the shared vocabulary and permitting
all seven canonical values for AuditEvent resolves both mismatches without
removing or renaming any existing FEATURE-0012 serialized value. Confirming
`metadata.scopeRef` as the sole canonical scope identity eliminates ambiguity
about where scope is expressed and prevents divergent enum proliferation.

## Reuse-before-build assessment

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - FEATURE-0012 ScopeKind vocabulary and conformance foundation
  - FEATURE-0013 canonical architecture and ADH-2026-014 decisions
- Sovrunn-owned responsibility summary:
  - additive ScopeKind value (ServiceInstance)
  - sole canonical scope identity confirmation (`metadata.scopeRef`)
  - schema, binding, validator, and fixture synchronization
- Non-goals summary:
  - no runtime capability of any kind
  - no new resource kinds or governance primitives
  - no product or provider selection

## Phase impact

- Current phase allowed: Yes
- Phase 2 output:
  - updated canonical ScopeKind schema with seven values
  - updated Go bindings, validators, and fixtures
  - compatibility evidence preserving FEATURE-0012 values
  - synchronized architecture, requirements, and design documentation
- Cross-phase effect:
  - later features consuming ScopeKind inherit the seven-value vocabulary
- Current phase boundary impact:
  - no infrastructure or lifecycle execution
  - no production runtime implementation
  - no change to the Phase 2 feature sequence

## Conflict check

- Conflicts with accepted DEC/RFC: Yes, one controlled baseline delta is
  approved by this handoff.
- Conflicting or incomplete baseline statements:
  - ADH-2026-014 states six formal governance scopes for AuditEvent; this
    handoff establishes seven canonical scopes.
- Resolution required:
  - apply this approved ADH;
  - update schemas, Go bindings, validators, fixtures, architecture,
    requirements, and design to reflect seven canonical ScopeKind values;
  - preserve all existing FEATURE-0012 serialized values as valid;
  - preserve compatibility evidence.

## Required action

- Supersede the ADH-2026-014 six-scope AuditEvent statement with the
  seven-scope vocabulary established here.
- The previous requirements `APPROVED_FOR_DESIGN` is superseded for continued
  design work by ADH-2026-015.
- Requirements must receive fresh independent approval before design resumes.
- The existing design and its review history must be preserved as evidence.
- Tasks and implementation remain unauthorized.
- Synchronize canonical schemas, Go bindings, validators, fixtures,
  compatibility evidence, architecture, requirements, and design.
- Validate this handoff with `make arch-handoff-check`.
- Kiro must not choose mappings, aliases, removals, additional scopes, or
  applicability exclusions.
- Any future scope change requires a separately approved architecture handoff.

## Impacted files

- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-015-feature-0013-scope-vocabulary-clarification.md`
- `.kiro/specs/decision-object-and-auditevent-standard/requirements.md`
- `.kiro/specs/decision-object-and-auditevent-standard/design.md`
- Canonical ScopeKind schema files
- Go ScopeKind bindings and validators
- ScopeKind and AuditEvent test fixtures

## Impacted features

- FEATURE-0012:
  - existing serialized ScopeKind values remain valid; receives additive
    ServiceInstance value in ScopeKind vocabulary and related compatibility
    evidence; no runtime behavior change.
- FEATURE-0013:
  - owns the clarified seven-scope AuditEvent contract via `metadata.scopeRef`
    as sole canonical scope identity; schemas, bindings, validators, fixtures,
    and documentation must be synchronized.

## Acceptance criteria for Kiro update

- [ ] Value-by-value compatibility: every existing FEATURE-0012 serialized
      ScopeKind value (Platform, Organization, OrganizationUnit, Tenant,
      Project, Provider) remains valid and unchanged.
- [ ] Sole canonical scope identity: AuditEvent uses `metadata.scopeRef` as the
      only field expressing scope; no parallel scope enum or independently
      mutable scope field exists.
- [ ] Seven-scope schema/binding/validator synchronization: canonical schemas,
      Go bindings, and validators all reflect exactly Platform, Organization,
      OrganizationUnit, Tenant, Project, Provider, and ServiceInstance.
- [ ] Regression coverage for all FEATURE-0012 values: test fixtures and
      validation tests cover each of the six previously approved values.
- [ ] ServiceInstance positive fixtures: at least one positive test fixture
      exercises ServiceInstance as a valid ScopeKind value in AuditEvent.
- [ ] Contradictory-scope rejection: validators reject any AuditEvent that
      carries a scope value not in the canonical seven-value vocabulary.
- [ ] No parallel enum: no separate AuditEvent-specific scope enum or
      independently mutable scope field is introduced anywhere in schemas,
      bindings, or validators.
- [ ] Existing design and review history preserved as evidence.
- [ ] Tasks and implementation remain unauthorized.
- [ ] Kiro does not choose mappings, aliases, removals, additional scopes, or
      applicability exclusions.

## Explicit instructions to Kiro

- Supersede the ADH-2026-014 six-scope AuditEvent statement with the
  seven-scope vocabulary from this handoff.
- The previous requirements `APPROVED_FOR_DESIGN` is superseded; requirements
  must receive fresh independent approval before design resumes.
- Preserve the existing design and its review history as evidence.
- Do not authorize tasks or implementation.
- Synchronize canonical schemas, Go bindings, validators, fixtures,
  compatibility evidence, architecture, requirements, and design to reflect
  seven canonical ScopeKind values.
- Do not choose mappings, aliases, removals, additional scopes, or
  applicability exclusions.
- Any future scope change requires a separately approved architecture handoff.
- Do not modify review JSON, automation state, approval tokens, or Go source
  files outside the explicit scope of this handoff.
- Do not introduce runtime capabilities of any kind.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-07-27
- Notes:
  - ADH-2026-014 remains controlling except that its six-scope AuditEvent
    statement is clarified and superseded by this handoff.
  - ServiceInstance is added additively to the canonical ScopeKind vocabulary;
    Provider is retained unchanged.
  - ScopeKind is confirmed as the single canonical shared vocabulary.
  - `metadata.scopeRef` is confirmed as the sole canonical scope identity for
    AuditEvent.
  - All seven canonical scopes are initially permitted for AuditEvent.
  - The previous requirements `APPROVED_FOR_DESIGN` is superseded; fresh
    independent approval is required before design resumes.
  - Tasks and implementation remain unauthorized.
  - This approval authorizes contract-only work: no persistence, network,
    queue, worker, workflow, cryptographic service, external integration, or
    other runtime capability.
