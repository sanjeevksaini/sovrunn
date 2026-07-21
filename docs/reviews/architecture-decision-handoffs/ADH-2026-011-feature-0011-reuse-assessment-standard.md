# Architecture Decision Handoff

## Metadata

- **Handoff ID:** ADH-2026-011
- **Date:** 2026-07-21
- **Source discussion:** Sovrunn Phase 2 Architecture Spine and FEATURE-0011 Decision Preparation
- **Related feature:** FEATURE-0011
- **Related phase:** Phase 2
- **Author:** System Architect and Architecture Reviewer
- **Human approver:** Sovrunn project owner
- **Approval evidence:** Approval recorded in architecture review conversation
- **Approval status:** Approved

## Decision title

Capability-Level Reuse Assessment and Risk-Mitigation Contract for Phase 2 and Later Features

## Summary

FEATURE-0011 will establish the mandatory Sovrunn Reuse Assessment contract used by FEATURE-0011 and every later feature.

Assessments will:

- classify significant architectural capability units as Reuse, Wrap, Extend, or Build;
- identify Sovrunn and external responsibilities;
- determine whether an adapter boundary is required;
- evaluate sovereignty, security, operations, licensing, and replacement considerations;
- identify applicable architecture risks;
- define preventive, detection, and corrective controls;
- preserve current-phase boundaries;
- provide architecture traceability;
- define reassessment triggers.

The standard is an Architecture Operating System governance contract.

It will not:

- create a runtime API resource;
- select future integration products;
- implement a production integration;
- introduce Go production code.

## Classification

**Extend**

## Existing approved baseline

The approved baseline establishes that:

- reuse before build is mandatory;
- every feature requires a Reuse Assessment;
- Sovrunn owns governance, policy context, placement, decisions, audit, explanations, and customer/provider experience;
- mature systems should be reused or wrapped for policy, identity, secrets, workflow, observability, service operators, and persistence;
- adapter boundaries are required before integrating replaceable external engines;
- Phase 2 is limited to models, standards, decisions, audit, adapter boundaries, policy context, plugin taxonomy, and placement simulation;
- real provider provisioning and full external-engine integrations are outside Phase 2.

Relevant references:

- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/phase2/PHASE2_EXECUTION_STRATEGY.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md`
- `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
- `docs/phase2/PHASE2_ACCEPTANCE_GATES.md`
- `docs/decisions/DEC-0026-reuse-before-build.md`
- `docs/decisions/DEC-0036-adapter-boundaries.md`
- `docs/rfc/RFC-0021-reuse-first-architecture.md`
- `docs/foundation/constitution.md`

## Approved decision

FEATURE-0011 will define one canonical Reuse Assessment contract with the following rules:

1. Every FEATURE-0011-and-later feature must contain a feature-level reuse summary.
2. Each significant architectural capability or integration may have its own assessment.
3. Each assessment must select exactly one disposition for each assessed capability:
   - Reuse;
   - Wrap;
   - Extend;
   - Build.
4. Each assessment must identify:
   - assessment scope;
   - mature candidates or applicable standards;
   - rationale;
   - Sovrunn-owned responsibility;
   - external or reused responsibility;
   - adapter requirement and rationale;
   - sovereignty and deployment fit;
   - security and trust considerations;
   - operational and support implications;
   - licensing and supply-chain considerations;
   - provider-neutrality impact;
   - current-phase impact;
   - explicit non-goals;
   - applicable architecture risks;
   - preventive controls;
   - detection controls;
   - corrective path;
   - residual risk;
   - replacement risk;
   - reassessment triggers;
   - architecture traceability.
5. A feature may have multiple dispositions when it contains distinct capabilities.
6. A Build decision must explicitly explain why Reuse, Wrap, and Extend are insufficient.
7. Conceptual examples may be included to improve human and AI understanding.
8. Every conceptual example must be explicitly marked as illustrative and outside execution scope.
9. A reuse assessment does not approve a vendor or architecture change by itself.
10. Approved changes to a disposition, responsibility boundary, or mitigation plan require a new Architecture Decision Handoff.
11. The strict feature gate will validate assessment structure and controlled values.
12. Human architecture review remains responsible for semantic quality.
13. FEATURE-0011 is documentation, governance, and validation work only.
14. FEATURE-0011 will not introduce a runtime resource, production integration, or Go service.

## Rationale

A feature-wide single-label classification is insufficient because most Sovrunn platform features combine Sovrunn-owned control-plane intelligence with reused infrastructure.

Capability-level assessment makes responsibility boundaries explicit and prevents:

- external APIs leaking into core;
- wrappers becoming duplicate engines;
- Build becoming an unreviewed default;
- future tool candidates becoming prematurely approved;
- Phase 2 expanding into runtime integration;
- risks being documented without controls;
- repository inconsistencies propagating into later features.

Structural feature-gate validation creates repeatability.

Human review preserves architectural judgment.

Risk-mitigation fields make the assessment actionable through Kiro, Cursor, testing, the feature gate, and pull-request review.

## Reuse-before-build assessment

- **Decision:** Extend
- **Mature reusable options considered:**
  - architecture decision records;
  - RFC review practices;
  - software adoption and build-versus-buy assessments;
  - risk-control registers;
  - the existing Sovrunn reuse-before-build baseline;
  - the existing draft Phase 2 reuse assessment format.
- **Sovrunn-owned responsibility:**
  - four-disposition vocabulary;
  - capability-level assessment rules;
  - sovereign deployment criteria;
  - provider-neutrality checks;
  - adapter-boundary requirements;
  - Phase 2 scope controls;
  - future-feature mitigation requirements;
  - architecture traceability;
  - feature-gate structure.
- **Reused or extended responsibility:**
  - general architecture-decision practices;
  - software-selection practices;
  - risk-management practices;
  - mature component documentation and conformance evidence.
- **Non-goals:**
  - product procurement automation;
  - production vendor selection;
  - automatic architectural approval;
  - full assessment of later Phase 2 features;
  - runtime `ReuseAssessment` resources;
  - Go production code.

## Risk and mitigation impact

FEATURE-0011 will define reusable controls for:

- incorrect capability classification;
- vendor-first architecture;
- direct external coupling;
- wrapper responsibility expansion;
- unjustified custom build;
- adapter omission;
- phase-scope leakage;
- missing replacement planning;
- architecture-status ambiguity;
- template divergence;
- traceability errors;
- silent later-feature overrides.

Every later feature must document:

- which risks apply;
- preventive controls;
- detection controls;
- corrective path;
- residual risk;
- reassessment triggers.

## Phase impact

- **Current phase allowed?** Yes
- **Target phase:** Phase 2
- **Current phase boundary impact:**
  - establishes the governance gate required before all later Phase 2 features;
  - adds conceptual examples for human and AI understanding;
  - adds cross-feature risk-mitigation requirements;
  - does not add runtime behavior;
  - does not alter the approved Phase 2 sequence;
  - does not select or integrate deferred external engines.

## Conflict check

- **Conflicts with accepted DEC/RFC?** No
- **Conflicting decisions:** None identified
- **Resolution required:**
  - extend the existing standard and reuse-first RFC;
  - correct FEATURE-0011 traceability to reference DEC-0026;
  - verify canonical constitution references;
  - align templates and gates;
  - add mitigation fields without changing the Phase 2 objective;
  - use an Architecture Change Request only if repository validation discovers a conflicting accepted rule.

## Required action

- Update architecture documentation.
- Update the Phase 2 reuse assessment standard.
- Update Phase 2 acceptance gates.
- Update the reuse-first RFC where required.
- Create the FEATURE-0011 feature document.
- Generate and approve Kiro `requirements.md`.
- Generate and approve Kiro `design.md`.
- Generate and approve Kiro `tasks.md`.
- Update traceability matrices.
- Update feature-gate checks.
- Align Feature Factory and reviewer prompts.
- Add structural validation fixtures or equivalent evidence.
- Add risk-mitigation guidance for future features.
- Add explicit labeling rules for conceptual examples.

## Impacted files

Primary files:

- `docs/features/FEATURE-0011-reuse-assessment-standard.md`
- `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`
- `docs/phase2/PHASE2_ACCEPTANCE_GATES.md`
- `docs/rfc/RFC-0021-reuse-first-architecture.md`
- `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md`
- `docs/traceability/DECISION_TRACEABILITY_MATRIX.md`
- `scripts/feature-gate.sh`

Templates and automation to inspect and align:

- `docs/templates/ARCHITECTURE_DECISION_HANDOFF.md`
- `docs/templates/ARCHITECTURE_CHANGE_REQUEST.md`
- `docs/templates/RFC_TEMPLATE.md`
- `docs/templates/FEATURE_REVIEW_TEMPLATE.md`
- `docs/automation/FEATURE_FACTORY.md`
- `docs/ai/AI_FEATURE_FACTORY.md`
- `docs/prompts/kiro/requirements.prompt.md`
- `docs/prompts/kiro/design.prompt.md`
- `docs/prompts/kiro/tasks.prompt.md`
- `docs/prompts/cursor/task.prompt.md`
- `docs/prompts/reviewer/spec-review.prompt.md`
- `docs/prompts/reviewer/approval-review.prompt.md`

Kiro specification files to create:

- `.kiro/specs/reuse-assessment-standard/requirements.md`
- `.kiro/specs/reuse-assessment-standard/design.md`
- `.kiro/specs/reuse-assessment-standard/tasks.md`

## Impacted features

### FEATURE-0011

Creates and validates the canonical standard.

### FEATURE-0012 through FEATURE-0026

Must consume:

- the reuse-assessment contract;
- the mitigation structure;
- the controlled vocabulary;
- the conceptual-example labeling rule;
- the change-control rule.

Later features must not be designed in detail as part of FEATURE-0011.

### Future features

Inherit the assessment and mitigation contract unless a later approved decision changes it.

## Acceptance criteria for Kiro update

- [ ] Handoff validated against Architecture Operating System files.
- [ ] Human approval status checked before applying the handoff.
- [ ] FEATURE-0011 feature document created with explicit acceptance criteria.
- [ ] One canonical Reuse Assessment format established.
- [ ] Assessment supports capability-level classification and a feature-level summary.
- [ ] Controlled dispositions are exactly Reuse, Wrap, Extend, and Build.
- [ ] Disposition semantics are unambiguous.
- [ ] Build requires explicit rejection rationale for Reuse, Wrap, and Extend.
- [ ] Sovrunn and external responsibilities are required fields.
- [ ] Adapter requirement and rationale are required fields.
- [ ] Sovereignty, security, operational, licensing, and provider-neutrality considerations are represented proportionately.
- [ ] Current-phase impact and deferred work are explicit.
- [ ] Non-goals are mandatory.
- [ ] Applicable architecture risks are mandatory.
- [ ] Preventive controls are mandatory.
- [ ] Detection controls are mandatory.
- [ ] Corrective path is mandatory.
- [ ] Residual risk is mandatory.
- [ ] Replacement risk is mandatory.
- [ ] Reassessment triggers are mandatory.
- [ ] Conceptual examples are allowed for understanding.
- [ ] Conceptual examples are explicitly labeled as outside execution scope.
- [ ] Reuse assessment is explicitly non-authoritative until architecture approval.
- [ ] Strict feature-gate validation checks required structure and controlled values.
- [ ] Semantic quality remains a human review responsibility.
- [ ] Duplicated formats in prompts and factory documents are aligned with the canonical standard.
- [ ] FEATURE-0011 traceability points to DEC-0026 and RFC-0021.
- [ ] Existing incorrect or inconsistent architecture references are corrected.
- [ ] Repository risks have preventive, detection, and corrective controls.
- [ ] Future-feature mitigation requirements are documented.
- [ ] No provider, policy engine, identity provider, secret backend, workflow engine, or PostgreSQL operator is selected.
- [ ] No runtime `ReuseAssessment` resource is introduced.
- [ ] No production Go code is required.
- [ ] Phase 2 non-goals remain intact.
- [ ] Feature sequence remains unchanged.
- [ ] Appropriate automated validation evidence is included.
- [ ] Final feature-gate review is recorded before merge.

## Explicit instructions to Kiro

- Treat this handoff as Approved.
- Generate only the artifact authorized by the current Kiro stage.
- Do not introduce new architecture decisions beyond this handoff.
- Do not design FEATURE-0012 or later features in detail.
- Do not evaluate or select production vendors.
- Do not turn `ReuseAssessment` into a customer-facing, provider-facing, or internal runtime resource.
- Do not introduce Go production code for FEATURE-0011.
- Do not change the Phase 2 feature sequence.
- Do not weaken DEC-0026 or DEC-0036.
- Do not permit a single feature-level label to hide materially different component dispositions.
- Do not treat automated validation as architecture approval.
- Preserve the distinction between proposed candidates and accepted decisions.
- Include mitigation fields in the canonical assessment.
- Ensure later features identify applicable risks and controls.
- Mark conceptual examples as illustrative and outside execution scope.
- Do not interpret conceptual examples as implementation requirements.
- Do not update `CURRENT_ARCHITECTURE_BASELINE.md` unless a separately approved decision requires a baseline change.
- Correct traceability metadata without changing unrelated accepted decisions.
- Keep future integrations in Non-goals, Deferred, or Future Phase sections only.
- Stop and report a conflict if the approved architecture cannot be represented without a new decision.

## Human approval

- **Approval status:** Approved
- **Approved by:** Sovrunn project owner
- **Date:** 2026-07-21
- **Notes:**
  - Approval authorizes Kiro to create FEATURE-0011 requirements, design, and tasks through separate approval stages.
  - Approval does not authorize Cursor implementation until Kiro artifacts are reviewed and approved.
  - Approval of conceptual examples does not authorize execution or detailed design of later features.
