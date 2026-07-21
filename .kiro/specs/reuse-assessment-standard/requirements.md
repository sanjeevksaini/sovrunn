# Requirements Document

Feature: FEATURE-0011 — Reuse Assessment Standard
Phase: Phase 2 — Reuse-First PaaS Fabric Foundation
Stage: Requirements

## Introduction

This document specifies the requirements for FEATURE-0011, the canonical
Sovrunn Reuse Assessment Standard. FEATURE-0011 is an Architecture
Operating System governance standard: documentation, governance, and
validation work only. It defines the mandatory reuse assessment contract
that every FEATURE-0011-and-later Phase 2 feature must satisfy before
implementation begins. It does not introduce a runtime resource, select a
vendor, or add Go production code, and it preserves all approved Phase 2
non-goals.

## Metadata

| Field | Value |
|---|---|
| Feature ID | FEATURE-0011 |
| Feature title | Reuse Assessment Standard |
| Phase | Phase 2 — Reuse-First PaaS Fabric Foundation |
| Spec stage | Requirements |
| Branch | feature-0011-reuse-assessment-standard |
| Architecture baseline | ARCH-2026.07-PHASE2-START |
| Controlling handoff | ADH-2026-011 (Approved) |
| Classification | Extend |
| Classification rationale | Extends the existing reuse-before-build baseline and reuse-first RFC without introducing a new architecture decision. |
| Artifact type | Architecture Operating System governance standard |
| Status | Draft — pending review |

## FEATURE-0011 reuse summary

Feature identity: FEATURE-0011 — Reuse Assessment Standard.

The following summary is populated only from the approved Architecture
Decision Handoff (ADH-2026-011) and the approved Phase 2 Architecture
Spine. It introduces no new capability classification.

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|
| Reuse Assessment Standard governance contract | Extend | Extends the existing reuse-before-build baseline, architecture-decision and RFC review practices, risk-control registers, and the existing draft Phase 2 reuse assessment format, rather than building a new governance mechanism. | Approved | ADH-2026-011 (DEC-0026, RFC-0021) |

Per ADH-2026-011, Sovrunn owns the four-disposition vocabulary,
capability-level assessment rules, sovereign deployment criteria,
provider-neutrality checks, adapter-boundary requirements, Phase 2 scope
controls, future-feature mitigation requirements, architecture
traceability, and feature-gate structure. General architecture-decision,
software-selection, and risk-management practices are the reused or
extended responsibility. The approved feature-level disposition is Extend.

The metadata Classification (Extend) does not replace this feature-level
reuse summary.

## Purpose

FEATURE-0011 defines the single canonical Sovrunn Reuse Assessment
Standard. The standard is the mandatory governance contract that every
FEATURE-0011-and-later Phase 2 feature must satisfy before implementation
begins.

The standard exists to make reuse-before-build decisions explicit,
capability-level, auditable, and mitigable, so that Sovrunn owns its
control-plane intelligence while reusing or wrapping mature external
foundations behind adapter boundaries.

FEATURE-0011 is documentation, governance, and validation work only. It
does not introduce a runtime resource, a production integration, vendor
selection, or Go production code.

## Background and controlling decisions

This standard extends the accepted reuse-before-build and adapter-boundary
decisions. It does not create new architecture decisions.

- DEC-0026 — Reuse Before Build (Accepted)
- DEC-0036 — Adapter Boundaries Before External Integration (Accepted)
- RFC-0021 — Reuse-First Architecture (Draft)
- ADH-2026-011 — Capability-Level Reuse Assessment and Risk-Mitigation
  Contract (Approved), treated as the controlling FEATURE-0011 decision
- Constitution Principles 4, 16, 17, 18, 19, 20
- PHASE2_ARCHITECTURE_SPINE Invariant A (reuse before build) and P2-C01,
  P2-C02, P2-C15 cross-feature contracts

## Scope

### In scope

- One canonical Reuse Assessment Standard document format.
- A required feature-level reuse summary for every FEATURE-0011-and-later
  feature.
- Zero or more capability-level reuse assessments per feature.
- A controlled disposition vocabulary: Reuse, Wrap, Extend, Build.
- Mandatory assessment fields, including risk-mitigation fields.
- Rules distinguishing automated validation from human semantic review.
- Rules for conceptual-example labeling.
- Repository-consistency rules that keep prompts, factory documents, and
  gates aligned with the canonical standard.
- Decision-status and change-control rules for assessments.
- Reassessment lifecycle and future-feature mitigation-record rules.

### Out of scope

- A runtime `ReuseAssessment` API resource.
- Selection or integration of any production vendor or external engine.
- Runtime provisioning, provider integration, or plugin execution.
- Production persistence, billing, failover, or disaster recovery.
- Autonomous AI operations.
- Detailed design of FEATURE-0012 through FEATURE-0026.
- Go production code for FEATURE-0011.
- Design-stage or task-stage content for this feature.

## Actors and stakeholders

| Actor | Interest in the standard |
|---|---|
| Architecture owner | Approves assessments; owns semantic quality of reuse analysis. |
| Kiro (planning) | Produces requirements, design, and tasks that embed the standard. |
| Cursor (implementation) | Implements approved specification tasks and must not independently change approved architecture semantics. |
| Reviewer | Confirms structure, controlled values, and semantic adequacy. |
| Feature gate automation | Performs structural and consistency validation only. |
| Future feature authors | Apply the standard to FEATURE-0012 and later features. |

## Assumptions

- The Architecture Decision Handoff ADH-2026-011 is Approved and is the
  controlling FEATURE-0011 decision.
- The reuse-before-build and adapter-boundary decisions remain Accepted
  and are not weakened by this standard.
- Phase 2 remains model-, standard-, decision-, and simulation-oriented,
  with all provisioning and integration deferred.
- A strict feature gate exists for FEATURE-0011 and later features and can
  perform structural and consistency checks.
- Human architecture review is available and remains the authority for
  semantic quality and architecture approval.

## Constraints

- The standard must remain a documentation and governance contract, not a
  runtime domain object.
- The controlled disposition vocabulary is exactly Reuse, Wrap, Extend,
  and Build; no other disposition values are permitted.
- The standard must not select, rank, or approve any production vendor,
  policy engine, identity provider, secret backend, workflow engine, or
  database operator.
- The standard must preserve provider-neutrality and Phase 2 non-goals.
- Automated validation must never constitute architecture approval or
  product selection.
- Conceptual examples must not authorize implementation or expand scope.
- No new architecture decision may be introduced; if the approved
  architecture is insufficient, the work must stop and report a conflict.

## Dependencies

- Governance dependency on ADH-2026-011, DEC-0026, DEC-0036, RFC-0021.
- Controlling scope from PHASE2_SCOPE, PHASE2_EXECUTION_STRATEGY,
  PHASE2_ACCEPTANCE_GATES, and PHASE2_ARCHITECTURE_SPINE.
- Supersedes and consolidates the draft PHASE2_REUSE_ASSESSMENT_STANDARD
  format into one canonical, versioned standard.
- FEATURE-0012 through FEATURE-0026 depend on this standard as a governance
  prerequisite but are not designed here.

## Glossary

| Term | Definition |
|---|---|
| Reuse Assessment | The governance record that classifies how a capability is delivered and records its risk mitigation. |
| Feature-level reuse summary | A per-feature summary enumerating every capability disposition. |
| Capability-level assessment | An assessment for one significant architectural capability or integration. |
| Disposition | One of the controlled values Reuse, Wrap, Extend, or Build. |
| Adapter boundary | A Sovrunn-owned contract that isolates core logic from an external engine (DEC-0036). |
| Decision status | The single canonical status of an assessment, drawn from the controlled vocabulary Proposed, Approved, Deferred, Rejected, and Superseded. Only Approved, with recorded human approval, is authoritative. |
| Reassessment trigger | A condition that requires an assessment to be revisited. |
| Feature gate | The strict automated check applied to FEATURE-0011 and later features. |

## Requirements

Requirements are grouped as functional, governance, validation,
traceability, and risk-mitigation. Acceptance criteria use EARS keywords
(WHEN, WHERE, WHILE, IF, THEN, THE, SHALL) and are independently testable.

**Functional requirements**

### Requirement 1: Canonical standard and feature-level reuse summary

**User Story:** As an architecture owner, I want one canonical reuse
assessment format with a mandatory feature-level summary, so that every
feature declares its reuse posture consistently.

#### Acceptance Criteria

1. THE standard SHALL define exactly one canonical Reuse Assessment format
   used by FEATURE-0011 and every later feature.
2. THE standard SHALL require every FEATURE-0011-and-later feature to
   include exactly one feature-level reuse summary.
3. WHERE a feature has multiple capabilities, THE feature-level summary
   SHALL enumerate each capability disposition without collapsing them
   into a single label.
4. THE standard SHALL require the feature-level summary to appear in the
   architecture contract, Kiro requirements, design, tasks, and final
   review.
5. IF a feature-level reuse summary is absent, THEN THE standard SHALL
   define the feature as non-compliant.

### Requirement 2: Capability-level reuse assessments

**User Story:** As a reviewer, I want reuse classification applied at the
capability level, so that reused, wrapped, and built components are not
hidden behind a single feature label.

#### Acceptance Criteria

1. THE standard SHALL allow zero or more capability-level assessments per
   feature, one per significant architectural capability or integration.
2. WHEN a feature contains distinct capabilities, THE standard SHALL
   require a separate assessment for each significant capability.
3. THE standard SHALL define "significant architectural capability" as a
   decision unit whose disposition affects ownership, coupling, or
   replacement cost.
4. THE standard SHALL prohibit a single feature-level label from hiding
   materially different component dispositions.

### Requirement 3: Controlled dispositions and decision status

**User Story:** As a feature gate maintainer, I want a fixed disposition
vocabulary and an explicit decision status, so that assessments are
machine-checkable and unambiguous.

#### Acceptance Criteria

1. THE standard SHALL define the disposition vocabulary as exactly Reuse,
   Wrap, Extend, and Build.
2. THE standard SHALL require each capability assessment to select exactly
   one disposition from that vocabulary.
3. THE standard SHALL define disposition semantics aligned with the
   approved Architecture Spine: Reuse adopts a mature implementation,
   protocol, or standard substantially as provided, where Sovrunn does not
   fork its core behavior and an adapter may still be required; Wrap places
   a Sovrunn-owned contract around a mature capability without recreating
   the wrapped engine; Extend adds behavior through supported extension,
   composition, or compatible augmentation, where a maintained fork
   requires separate approval and maintenance assessment; Build implements
   a Sovrunn-owned differentiating capability or addresses the absence of
   an acceptable mature fit, requires rejection rationale for Reuse, Wrap,
   and Extend, and defines long-term ownership.
4. IF a disposition value is not one of the four controlled values, THEN
   THE standard SHALL define the assessment as invalid.
5. THE standard SHALL define the decision-status vocabulary as exactly
   Proposed, Approved, Deferred, Rejected, and Superseded.
6. THE standard SHALL require exactly one decision status per capability
   assessment.
7. IF a decision status is not one of the five controlled values, THEN THE
   standard SHALL define the status as invalid.
8. THE standard SHALL state that only the Approved status, with recorded
   human approval, may authorize progression to source implementation or
   implementation execution.
9. THE standard SHALL state that progression between Kiro Requirements,
   Design, and Tasks stages remains governed by the separate
   stage-approval tokens, and that an Approved assessment status does not
   substitute for APPROVED_FOR_DESIGN, APPROVED_FOR_TASKS, or the
   implementation authorization gate.
10. THE standard SHALL state that an assessment is non-authoritative until
    an Approved status with recorded human approval exists.

### Requirement 4: Descriptive and responsibility fields

**User Story:** As an architecture owner, I want each assessment to state
scope, candidates, rationale, and responsibility boundaries, so that
ownership and coupling are explicit.

#### Acceptance Criteria

1. THE standard SHALL require an assessment scope statement for each
   capability assessment.
2. THE standard SHALL require the mature candidates or applicable
   standards considered to be listed.
3. THE standard SHALL require a rationale for the selected disposition.
4. THE standard SHALL require an explicit Sovrunn-owned responsibility
   field.
5. THE standard SHALL require an explicit reused or external responsibility
   field.
6. THE standard SHALL require the data that crosses the boundary to be
   described.
7. THE standard SHALL require the control that crosses the boundary to be
   described.
8. THE standard SHALL require the following mandatory canonical fields in
   each capability assessment: feature identity; capability or
   decision-unit identity; assessment owner; related DEC, RFC, and ADH
   references; candidate category; relevant candidate strengths; material
   candidate constraints; selected foundation or approach; adapter or
   contract identifier.
9. THE standard SHALL require a "vendor-native types allowed" field
   restricted to exactly No or an Approved exception reference.
10. THE standard SHALL require an "allowed in current phase" field
    restricted to exactly Yes or No.
11. THE standard SHALL require linked acceptance criteria to be recorded
    under the assessment traceability.

### Requirement 5: Adapter boundary fields

**User Story:** As an architect enforcing adapter boundaries, I want each
assessment to declare whether an adapter is required and why, so that
external engines remain replaceable.

#### Acceptance Criteria

1. THE standard SHALL require an explicit adapter requirement field with a
   value of exactly Yes or No.
2. THE standard SHALL require the automated gate to reject any
   adapter-requirement value other than Yes or No.
3. THE standard SHALL require a rationale for the adapter decision.
4. WHERE a disposition is Reuse, Wrap, or Extend against an external engine
   expected to evolve or be replaced, THE standard SHALL require the
   adapter requirement to be justified against DEC-0036.
5. IF an adapter rationale is missing for an assessment that touches a
   replaceable external engine, THEN THE standard SHALL define the
   assessment as invalid.

### Requirement 6: Suitability considerations

**User Story:** As a sovereignty-focused reviewer, I want each assessment
to evaluate sovereignty, security, operations, licensing, and
provider-neutrality, so that a technically capable component is also
suitable for sovereign deployment.

#### Acceptance Criteria

1. THE standard SHALL require a sovereignty and deployment-fit
   consideration, including disconnected or air-gapped suitability where
   relevant.
2. THE standard SHALL require a security and trust consideration.
3. THE standard SHALL require an operational and supportability
   consideration.
4. THE standard SHALL require a licensing and supply-chain consideration.
5. THE standard SHALL require a portability and provider-neutrality impact
   consideration.
6. WHERE a consideration does not apply, THE standard SHALL require an
   explicit statement of non-applicability rather than an omission.

### Requirement 7: Phase impact, non-goals, and exit boundary

**User Story:** As a Phase 2 scope owner, I want each assessment to
separate current-phase work from deferred work and to state non-goals and
an exit boundary, so that assessments cannot justify runtime integration
during Phase 2.

#### Acceptance Criteria

1. THE standard SHALL require an explicit current-phase work field.
2. THE standard SHALL require an explicit deferred work field.
3. THE standard SHALL require an explicit non-goals field.
4. THE standard SHALL require an exit or migration boundary field that
   describes how the capability could be replaced or migrated later.
5. WHEN an assessment describes a future integration, THE standard SHALL
   require that integration to appear only in deferred work, non-goals, or
   a future-phase section.
6. THE standard SHALL require each assessment to acknowledge the applicable
   Phase 2 non-goals.

### Requirement 8: Build rejection rationale

**User Story:** As a reviewer preventing unnecessary custom work, I want a
Build disposition to justify why reuse alternatives are insufficient, so
that Build is not an unreviewed default.

#### Acceptance Criteria

1. WHEN a disposition is Build, THE standard SHALL require an explicit
   explanation of why Reuse is insufficient.
2. WHEN a disposition is Build, THE standard SHALL require an explicit
   explanation of why Wrap is insufficient.
3. WHEN a disposition is Build, THE standard SHALL require an explicit
   explanation of why Extend is insufficient.
4. IF a Build disposition lacks rejection rationale for Reuse, Wrap, and
   Extend, THEN THE standard SHALL define the assessment as invalid.
5. THE standard SHALL require a Build disposition to reference the Sovrunn
   differentiation it protects.

### Requirement 9: Risk-mitigation fields

**User Story:** As an architecture owner, I want each assessment to carry
actionable risk-mitigation fields, so that documented risks always have
controls and an accepted residual position.

#### Acceptance Criteria

1. THE standard SHALL require a list of applicable architecture risks for
   each assessment.
2. THE standard SHALL require preventive controls for the identified risks.
3. THE standard SHALL require detection controls for the identified risks.
4. THE standard SHALL require a corrective path for the identified risks.
5. THE standard SHALL require an explicit residual risk statement.
6. THE standard SHALL require a replacement risk value drawn from an
   approved scale of Low, Medium, or High.
7. THE standard SHALL require reassessment triggers that define when the
   assessment must be revisited.
8. THE standard SHALL define a risk entry as incomplete unless it contains
   at least one preventive control, at least one detection control, and a
   corrective path.

### Requirement 10: Traceability and validation-evidence fields

**User Story:** As a reviewer, I want each assessment to trace to accepted
decisions and to carry validation evidence, so that assessments are
auditable and consistent with the architecture baseline.

#### Acceptance Criteria

1. THE standard SHALL require architecture traceability that references
   relevant DEC, RFC, architecture, and dependency sources.
2. THE standard SHALL require a validation and review evidence field.
3. WHEN an assessment references a DEC or RFC, THE standard SHALL require
   the reference to identify an existing decision or RFC record.
4. IF a referenced DEC or RFC record does not exist, THEN THE standard
   SHALL define the reference as invalid.
5. THE standard SHALL require FEATURE-0011 traceability to reference
   DEC-0026 and RFC-0021.
6. WHERE an assessment is adapter-related, THE standard SHALL require its
   traceability to reference DEC-0036 in addition to the FEATURE-0011
   references to DEC-0026 and RFC-0021.

**Governance requirements**

### Requirement 11: Decision status and change control

**User Story:** As a governance owner, I want assessments to be
non-authoritative until their decision status is Approved and to require a
handoff for changes, so that proposed candidates are never mistaken for
accepted decisions.

#### Acceptance Criteria

1. THE standard SHALL state that a reuse assessment does not approve a
   vendor or an architecture change by itself.
2. THE standard SHALL use the decision status as the single canonical
   status field and SHALL NOT define a separate approval status field.
3. THE standard SHALL require an Approved decision status with recorded
   human approval before an assessment is authoritative or before source
   implementation may begin.
4. WHEN a disposition, responsibility boundary, or mitigation plan changes
   after approval, THE standard SHALL require a new Architecture Decision
   Handoff.
5. THE standard SHALL prohibit later features from silently overriding an
   earlier Approved assessment.
6. IF source implementation is attempted without an Approved decision
   status, THEN THE standard SHALL define the action as non-compliant.

### Requirement 12: Conceptual example labeling

**User Story:** As an AI or human reader, I want conceptual examples to be
clearly non-authoritative, so that examples never trigger implementation
or scope expansion.

#### Acceptance Criteria

1. THE standard SHALL allow conceptual examples only to improve human and
   AI understanding.
2. WHEN a conceptual example is included, THE standard SHALL require it to
   be labeled exactly: "Conceptual example — illustrative only and outside
   execution scope".
3. THE standard SHALL state that conceptual examples do not authorize
   implementation, product selection, runtime execution, provider calls,
   plugin execution, provisioning, or detailed design of later features.
4. IF a conceptual example is present without the exact required label,
   THEN THE standard SHALL define the document as non-compliant.

**Validation requirements**

The standard SHALL define three distinct validation responsibilities.
Automated validation SHALL never approve architecture or select products.

### Requirement 13: Automated structural validation

**User Story:** As a feature gate maintainer, I want automated structural
checks, so that required sections and controlled values are always present.

#### Acceptance Criteria

1. THE standard SHALL define structural checks that verify all required
   sections exist.
2. THE standard SHALL define a check that each disposition uses the
   controlled vocabulary.
3. THE standard SHALL define a check that a decision status exists.
4. THE standard SHALL define a check that the adapter decision is explicit.
5. THE standard SHALL define a check that replacement risk uses an approved
   value.
6. THE standard SHALL define a check that non-goals exist.
7. THE standard SHALL define a check that the risk-mitigation fields exist.
8. THE standard SHALL define a check that architecture references exist.
9. THE standard SHALL define a check that decision status uses the exact
   controlled vocabulary (Proposed, Approved, Deferred, Rejected,
   Superseded).
10. THE standard SHALL define a check that the adapter requirement uses
    exactly Yes or No.
11. THE standard SHALL define a check that the assessment owner is present.
12. THE standard SHALL define a check that the capability or decision-unit
    identity is present.
13. THE standard SHALL define a check that the selected foundation or
    approach is present.
14. THE standard SHALL define a check that the adapter or contract
    identifier is present.
15. THE standard SHALL define a check that the current-phase field uses
    exactly Yes or No.
16. THE standard SHALL define a check that enforces the vendor-native-type
    rule (No or an Approved exception reference).
17. THE standard SHALL define a check that linked acceptance criteria are
    present.
18. THE standard SHALL define a check that each risk entry has both a
    preventive control and a detection control.

### Requirement 14: Automated consistency validation

**User Story:** As a reviewer, I want automated consistency checks, so that
invalid identifiers, missing rationale, and broken references are caught.

#### Acceptance Criteria

1. THE standard SHALL define a check that rejects invalid feature
   identifiers.
2. THE standard SHALL define a check that rejects unknown disposition
   values.
3. THE standard SHALL define a check that rejects a Build disposition with
   missing rejection rationale.
4. THE standard SHALL define a check that rejects a missing adapter
   rationale where an adapter is required.
5. THE standard SHALL define a check that rejects a missing Phase 2 scope
   acknowledgement.
6. THE standard SHALL define a check that detects inconsistent duplicated
   templates.
7. THE standard SHALL define a check that rejects nonexistent DEC or RFC
   references.
8. THE standard SHALL define a check that rejects implementation attempted
   without Approved status.
9. THE standard SHALL define a check that rejects a decision status outside
   the controlled vocabulary.
10. THE standard SHALL define a check that rejects an adapter requirement
    value other than Yes or No.
11. THE standard SHALL define a check that rejects a current-phase value
    other than Yes or No.
12. THE standard SHALL define a check that rejects a vendor-native-types
    value other than No or an Approved exception reference.
13. THE standard SHALL define a check that rejects a risk entry missing a
    preventive control, a detection control, or a corrective path.
14. THE standard SHALL define a check that rejects a missing adapter or
    contract identifier where an adapter is required.

### Requirement 15: Human semantic review

**User Story:** As an architecture owner, I want semantic quality to remain
a human responsibility, so that automation never substitutes for
architectural judgment.

#### Acceptance Criteria

1. THE standard SHALL assign human review responsibility for the quality of
   the reuse analysis.
2. THE standard SHALL assign human review responsibility for the
   correctness of responsibility boundaries.
3. THE standard SHALL assign human review responsibility for sovereignty
   suitability.
4. THE standard SHALL assign human review responsibility for adapter
   adequacy.
5. THE standard SHALL assign human review responsibility for Build
   justification.
6. THE standard SHALL assign human review responsibility for mitigation
   credibility.
7. THE standard SHALL assign human review responsibility for residual-risk
   acceptance.
8. THE standard SHALL assign human review responsibility for phase
   compliance.
9. THE standard SHALL state that automation must not approve architecture
   or select products.

**Traceability requirements**

### Requirement 16: Repository consistency and single canonical source

**User Story:** As a documentation owner, I want one versioned canonical
standard with all dependent documents aligned, so that the standard,
prompts, and gates cannot drift apart.

#### Acceptance Criteria

1. THE standard SHALL be defined in one canonical source document.
2. THE standard SHALL carry an explicit version for the canonical format.
3. WHEN the canonical format changes, THE standard SHALL require the
   version to change.
4. THE standard SHALL require Kiro requirements, design, and tasks prompts
   to align with the canonical format.
5. THE standard SHALL require Cursor prompts to align with the canonical
   format.
6. THE standard SHALL require reviewer prompts to align with the canonical
   format.
7. THE standard SHALL require Feature Factory documents to align with the
   canonical format.
8. THE standard SHALL require the controlled vocabulary to be consistent
   across all aligned documents.
9. THE standard SHALL require existing incorrect or inconsistent
   architecture references to be corrected without changing unrelated
   accepted decisions.
10. THE standard SHALL require FEATURE-0011 traceability entries in the
    feature and decision traceability matrices.

**Risk-mitigation requirements**

### Requirement 17: Reassessment lifecycle and future-feature records

**User Story:** As a governance owner, I want a reassessment lifecycle and
mandatory future-feature mitigation records, so that decisions stay current
and every later feature documents its risks and controls.

#### Acceptance Criteria

1. THE standard SHALL define reassessment triggers that include changes to
   requirements, licensing, maintenance status, security posture, or
   deployment context.
2. WHEN a reassessment trigger occurs, THE standard SHALL require the
   affected assessment to be revisited and its status re-evaluated.
3. THE standard SHALL require every FEATURE-0011-and-later feature to
   document which architecture risks apply.
4. THE standard SHALL require every such feature to document preventive
   controls, detection controls, and a corrective path.
5. THE standard SHALL require every such feature to document residual risk
   and reassessment triggers.
6. THE standard SHALL define the repository-level risks it mitigates,
   including incorrect classification, vendor-first architecture, direct
   external coupling, wrapper responsibility expansion, unjustified build,
   adapter omission, phase-scope leakage, missing replacement planning,
   status ambiguity, template divergence, traceability errors, and silent
   later-feature overrides.
7. THE standard SHALL provide preventive, detection, and corrective
   controls for each repository-level risk it defines.

### Requirement 18: Strict feature-gate enforcement

**User Story:** As a release owner, I want the strict feature gate to
enforce the standard for FEATURE-0011 and later, so that non-compliant
assessments cannot merge.

#### Acceptance Criteria

1. THE standard SHALL require the strict feature gate to validate assessment
   structure and controlled values for FEATURE-0011 and later features.
2. WHERE a feature is a Phase 1 legacy feature (FEATURE-0001 through
   FEATURE-0010), THE standard SHALL exempt it from strict Phase 2 section
   checks.
3. IF a required structural or consistency check fails, THEN THE standard
   SHALL require the feature gate to fail.
4. THE standard SHALL require a final feature-gate review to be recorded
   before merge.

## Non-goals

The following are explicitly out of scope for FEATURE-0011 and SHALL NOT be
produced by this standard:

- A runtime `ReuseAssessment` API resource of any kind.
- Selection, ranking, or approval of any production vendor or external
  engine, including policy engines, identity providers, secret backends,
  workflow engines, and database operators.
- Runtime provisioning, provider integration, or plugin execution design.
- Production persistence, billing, failover, or disaster recovery.
- Autonomous AI operations.
- Detailed design of FEATURE-0012 through FEATURE-0026.
- Go production code for FEATURE-0011.
- Any change to `CURRENT_ARCHITECTURE_BASELINE.md` absent a separately
  approved decision.
- Any new architecture decision beyond ADH-2026-011.
- Weakening of DEC-0026 or DEC-0036.
- Changes to the Phase 2 feature sequence.
- Design-stage or task-stage content within this requirements document.

## Conceptual example

> Conceptual example — illustrative only and outside execution scope

A later feature that introduces a policy evaluation abstraction might
record a capability assessment such as: Build the Sovrunn-owned evaluation
request and result contracts; Wrap a future external policy engine behind
an adapter; Reuse a mature policy language; and defer selecting the first
production engine. This illustrates capability-level dispositions only. It
does not select a product or authorize implementation.

## Non-normative requirements review checklist

This checklist is a non-normative reviewer aid. It does not replace or
modify the numbered acceptance criteria under Requirements 1–18, which
remain the normative source of truth.

FEATURE-0011 requirements are satisfied when:

1. One canonical Reuse Assessment format is established.
2. The format supports a mandatory feature-level summary and zero or more
   capability-level assessments.
3. The disposition vocabulary is exactly Reuse, Wrap, Extend, and Build,
   with unambiguous semantics.
4. A Build disposition requires explicit rejection rationale for Reuse,
   Wrap, and Extend.
5. Sovrunn-owned and external responsibilities, plus boundary data and
   control, are required fields.
6. Adapter requirement and rationale are required fields.
7. Sovereignty, security, operational, licensing, and provider-neutrality
   considerations are represented proportionately.
8. Current-phase work, deferred work, non-goals, and an exit or migration
   boundary are explicit.
9. Applicable risks, preventive controls, detection controls, corrective
   path, residual risk, replacement risk, and reassessment triggers are
   mandatory.
10. Architecture traceability and validation evidence are required, and
    FEATURE-0011 traceability references DEC-0026 and RFC-0021.
11. Automated structural and consistency validation is distinguished from
    human semantic review, and automation does not approve architecture or
    select products.
12. Conceptual examples are permitted and carry the exact required label.
13. One canonical, versioned source is defined, with Kiro, Cursor,
    reviewer, and Feature Factory documents aligned.
14. Decision-status enforcement, reassessment lifecycle, and future-feature
    mitigation records are required.
15. No runtime resource, vendor selection, Go production code, or Phase 2
    non-goal violation is introduced, and the feature sequence is
    unchanged.

## Unresolved issues

No open issues. The approved architecture in ADH-2026-011 and the Phase 2
controlling documents are sufficient to specify these requirements without
a new architecture decision.
