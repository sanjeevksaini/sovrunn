# Design Document

Feature: FEATURE-0011 — Reuse Assessment Standard
Phase: Phase 2 — Reuse-First PaaS Fabric Foundation
Stage: Design
Controlling handoff: ADH-2026-011 (Approved)

## Overview

FEATURE-0011 delivers the single canonical Sovrunn Reuse Assessment
Standard as an Architecture Operating System governance artifact. This
design describes a documentation-and-validation architecture, not a
runtime software component. There is no runtime `ReuseAssessment`
resource, no vendor selection, and no Go production code.

The design defines four coordinated elements:

1. A canonical, versioned standard document that specifies the required
   reuse assessment format, controlled vocabularies, and field semantics.
2. A reuse assessment document model: a mandatory feature-level summary
   plus zero or more capability-level assessments.
3. A three-layer validation architecture: automated structural checks,
   automated consistency checks, and human semantic review.
4. A repository-consistency model that keeps prompts, Feature Factory
   documents, and the feature gate aligned with the canonical standard.

The design maps directly to the approved requirements (Requirements 1–18)
and preserves the approved FEATURE-0011 feature-level reuse summary
(disposition: Extend) and all controlled vocabularies.

## FEATURE-0011 reuse summary

This is the approved feature-level reuse summary copied from the approved
requirements without changing its meaning.

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

## Controlling inputs

This design is derived only from the approved controlling inputs and
introduces no new architecture decision.

| Input | Role in this design |
|---|---|
| requirements.md (Approved) | Normative acceptance criteria (Requirements 1–18). |
| ADH-2026-011 (Approved) | Controlling decision; approves the Extend disposition and the mandatory field and mitigation contract. |
| PHASE2_ARCHITECTURE_SPINE | Disposition semantics, Invariant A, and cross-feature contracts P2-C01, P2-C02, P2-C15. |
| DEC-0026, DEC-0036, RFC-0021 | Reuse-before-build and adapter-boundary basis. |
| PHASE2_ACCEPTANCE_GATES, PHASE2_SCOPE | Feature-gate expectations and Phase 2 boundaries. |

## Design goals and non-goals

### Design goals

- Specify one canonical, versioned standard as the single source of truth.
- Make every mandatory assessment field and controlled vocabulary
  unambiguous and machine-checkable.
- Separate automated validation from human architecture judgment.
- Keep the standard documentation-only and provider-neutral.
- Preserve the approved feature-level reuse summary and Phase 2 non-goals.

### Design non-goals

- No runtime `ReuseAssessment` API resource or data store.
- No vendor, engine, or product selection.
- No runtime provisioning, provider integration, or plugin execution.
- No detailed design of FEATURE-0012 through FEATURE-0026.
- No new architecture decision beyond ADH-2026-011.
- No Go production code and no task-stage content.

## Design principles

- Single canonical source: exactly one document defines the standard;
  every other document references it rather than redefining it.
- Controlled vocabularies are closed sets; any value outside a set is
  invalid by definition.
- Capability-level classification is primary; the feature-level summary
  aggregates but never hides capability dispositions.
- Automation checks structure and controlled values; it never approves
  architecture or selects products.
- Human review owns semantic quality and architecture approval.
- Decision status is the single canonical status; only Approved with
  recorded human approval is authoritative.

## Architecture

FEATURE-0011 is a governance artifact, so its architecture is a
documentation-and-validation architecture rather than a runtime software
architecture. It has no services, endpoints, or persisted data.

```text
Canonical Reuse Assessment Standard (versioned document)
        |
        | defines format, fields, controlled vocabularies
        v
Reuse assessment content (feature-level summary + capability assessments)
        |
        | validated by
        v
Three-layer validation
  Layer 1 structural  ->  Layer 2 consistency  ->  Layer 3 human review
        |
        v
Strict feature gate (Layers 1–2) + recorded human approval (Layer 3)
        |
        v
Aligned prompts, Feature Factory docs, and traceability matrices
```

The remaining architecture subsections specify the canonical artifact, the
logical components and interfaces, the documentation data models, and the
validation behavior.

## Canonical artifact architecture

The canonical Reuse Assessment Standard is the existing canonical draft
designated by the approved Architecture Spine, not a new competing source:

- Exact canonical file path: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`.
- Normative document: this file is the single normative definition of the
  reuse assessment format, fields, and controlled vocabularies. Every other
  document referencing the standard is supporting, not normative.
- Format-version field: the canonical file records an explicit
  `reuse_assessment_format_version` field in its front matter.
- Initial format version: `1.0.0`.
- Versioning rule: the format version uses semantic versioning; any change
  to required fields, controlled vocabularies, or validation rules
  increments it (major for breaking field/vocabulary changes, minor for
  additive fields, patch for editorial clarification).
- Dependent-artifact referencing rule: every dependent artifact references
  the canonical file path and must not duplicate or redefine field
  definitions. Only required version-bearing artifacts declare a numeric
  version. Link-only supporting artifacts do not declare or duplicate the
  numeric version unless they contain a complete version-bearing example.

The following statuses are distinct and must not be conflated:

- Approval of ADH-2026-011: already granted; it authorizes this work.
- Canonical designation of the repository path: the Architecture Spine
  designates `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` as canonical.
- Current document status: the canonical file remains a draft to be
  consolidated by FEATURE-0011; it is not yet an approved standard.
- Final human approval of the implemented standard: granted only through
  the feature review process after implementation.

Cursor and automated tasks must not independently mark the governance
artifact Approved. Any final status transition requires recorded human
approval through the feature review process.

Requirement 16 is satisfied because FEATURE-0011 consolidates the existing
draft into this one canonical, versioned source. No new canonical location
is introduced, so no `ARCHITECTURE_DECISION_REQUIRED` condition arises.

```text
Canonical Reuse Assessment Standard (versioned)
    |
    +--> Kiro requirements / design / tasks prompts (aligned)
    +--> Cursor prompts (aligned)
    +--> Reviewer prompts (aligned)
    +--> Feature Factory documents (aligned)
    +--> Feature gate checks (structural + consistency)
    +--> Feature and decision traceability matrices
```

When the canonical format changes, its version changes and dependent
documents are re-aligned. This satisfies Requirement 16 (single canonical
source, versioning, and alignment).

## Repository artifact map

Paths below were verified against the repository. Artifacts marked [NEW]
are required by the approved requirements (Requirements 13–14, 16) and do
not yet exist; all others already exist and are aligned by reference.

| Artifact role | Repository path | Planned FEATURE-0011 change | Status | Requirements |
|---|---|---|---|---|
| Canonical reuse-standard document | `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` | Consolidate into the canonical format; add `reuse_assessment_format_version`; add mandatory fields, controlled vocabularies, and mitigation fields | Normative | 1,3,4,5,6,7,8,9,10,16 |
| Kiro prompt artifacts | `docs/prompts/kiro/requirements.prompt.md`, `docs/prompts/kiro/design.prompt.md`, `docs/prompts/kiro/tasks.prompt.md` | Align to canonical format by reference | Supporting | 16 |
| Kiro decision policy | `docs/automation/KIRO_DECISION_POLICY.md` | Align stage/approval wording to decision-status vocabulary | Supporting | 11,16 |
| Cursor prompt artifacts | `docs/prompts/cursor/task.prompt.md` | Align to canonical format by reference | Supporting | 16 |
| Reviewer prompt artifacts | `docs/prompts/reviewer/spec-review.prompt.md`, `docs/prompts/reviewer/approval-review.prompt.md` | Align semantic-review checklist to Layer 3 responsibilities | Supporting | 15,16 |
| Feature Factory documents | `docs/automation/FEATURE_FACTORY.md`, `docs/ai/AI_FEATURE_FACTORY.md` | Align to canonical format by reference | Supporting | 16 |
| Strict feature-gate entry point | `scripts/feature-gate.sh`, `scripts/phase2-scope-check.sh` | Orchestrate structural/consistency validation for FEATURE-0011+; `phase2-scope-check.sh` remains authoritative for Phase 2 scope-phrase detection (RA-C03) | Supporting | 13,14,18 |
| Reuse-assessment validator / validation rules | `scripts/reuse-assessment-check.sh` [NEW] | Bash entry point that may invoke Python 3 for deterministic Markdown parsing; implements Layer 1/Layer 2 rules; consumes or is orchestrated alongside `scripts/phase2-scope-check.sh` for RA-C03 rather than duplicating its blocked-phrase list | Supporting [NEW] | 13,14 |
| Feature traceability matrix | `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md` | Add FEATURE-0011 entry | Supporting | 16 |
| Decision traceability matrix | `docs/traceability/DECISION_TRACEABILITY_MATRIX.md` | Add FEATURE-0011 → DEC-0026/DEC-0036/RFC-0021 mapping | Supporting | 10,16 |
| Validation fixtures / tests | `tests/reuse-assessment/` [NEW] | Valid/invalid assessment fixtures per correctness property | Supporting [NEW] | 13,14 |
| FEATURE-0011 feature document | `docs/features/FEATURE-0011-reuse-assessment-standard.md` [NEW] | Create feature document with acceptance criteria and the approved reuse summary | Supporting [NEW] | 1,16 |
| Phase 2 acceptance gates | `docs/phase2/PHASE2_ACCEPTANCE_GATES.md` | Inspect and update only if inconsistent with strict validation | Supporting | 18 |
| Reuse-first RFC | `docs/rfc/RFC-0021-reuse-first-architecture.md` | Inspect and update only if inconsistent; add canonical-version alignment note | Supporting | 10,16 |
| Architecture Decision Handoff template | `docs/templates/ARCHITECTURE_DECISION_HANDOFF.md` | Inspect and update only if inconsistent with the canonical format | Supporting | 11,16 |
| Architecture Change Request template | `docs/templates/ARCHITECTURE_CHANGE_REQUEST.md` | Inspect and update only if inconsistent | Supporting | 11,16 |
| RFC template | `docs/templates/RFC_TEMPLATE.md` | Inspect and update only if inconsistent | Supporting | 16 |
| Feature review template | `docs/templates/FEATURE_REVIEW_TEMPLATE.md` | Align semantic-review checklist to Layer 3 responsibilities | Supporting | 15,16 |
| Human review evidence | `docs/reviews/feature-gates/FEATURE-0011-approval-review.md` [NEW] | Record final human approval of the implemented standard | Supporting [NEW] | 15,18 |

## Components and Interfaces

The architecture is composed of logical (documentation) components. None is
a runtime service; each is a document or a governance control.

| Component | Responsibility | Interface (contract) |
|---|---|---|
| Canonical Standard document | Defines the format, fields, and controlled vocabularies | The versioned standard consumed by all other components |
| Feature-level reuse summary | Aggregates capability dispositions per feature | Required section in feature contract, requirements, design, tasks, review |
| Capability-level assessment | Records disposition and mitigation for one capability | The mandatory field set defined in Data Models |
| Structural validator (Layer 1) | Checks presence and controlled values | Pass/fail signals to the feature gate |
| Consistency validator (Layer 2) | Checks cross-field rules and references | Pass/fail signals to the feature gate |
| Human semantic review (Layer 3) | Judges meaning and approves | Recorded human approval / decision status |
| Strict feature gate | Runs Layers 1–2 and records the review | Gate decision before merge |
| Aligned prompts and Factory docs | Reference the canonical standard | Alignment-by-reference, single vocabulary |

The primary interface authors interact with is the assessment document
format itself: a filled feature-level summary plus capability-level
assessments conforming to the Data Models below.

## Data Models

These are documentation schemas for assessment content. They are not
runtime data structures and are never persisted as a runtime resource. The
standard defines two nested structures.

### Feature-level reuse summary

Required in every FEATURE-0011-and-later feature and present in the
architecture contract, Kiro requirements, design, tasks, and final review
(Requirement 1). It aggregates each significant capability or decision
unit with its disposition, a concise rationale, decision status, and
controlling reference. It must not collapse distinct dispositions into a
single label (Requirement 2).

### Capability-level assessment

One assessment per significant architectural capability or decision unit
(Requirement 2). Each assessment carries the mandatory fields below.

```text
Identity
  - feature identity
  - capability or decision-unit identity
  - assessment owner
Classification
  - disposition (Reuse | Wrap | Extend | Build)
  - decision status (Proposed | Approved | Deferred | Rejected | Superseded)
Analysis
  - assessment scope
  - candidate category
  - mature candidates / applicable standards
  - relevant candidate strengths
  - material candidate constraints
  - rationale
  - selected foundation or approach
Boundary
  - Sovrunn-owned responsibility
  - reused or external responsibility
  - data crossing the boundary
  - control crossing the boundary
  - adapter required (Yes | No) + rationale
  - adapter or contract identifier
  - vendor-native types allowed (No | Approved exception reference)
```

The capability-level assessment continues with suitability, phase, risk,
and traceability fields.

```text
Suitability
  - sovereignty and deployment fit
  - security and trust
  - operational and supportability
  - licensing and supply-chain
  - portability and provider-neutrality impact
Phase and scope
  - allowed in current phase (Yes | No)
  - current-phase work
  - deferred work
  - explicit non-goals
  - exit or migration boundary
  - Phase 2 non-goal acknowledgement
Build justification (only when disposition is Build)
  - why Reuse is insufficient
  - why Wrap is insufficient
  - why Extend is insufficient
  - protected Sovrunn differentiation and long-term ownership
Risk mitigation
  - applicable architecture risks
  - preventive controls
  - detection controls
  - corrective path
  - residual risk
  - replacement risk (Low | Medium | High)
  - reassessment triggers
Traceability
  - related DEC / RFC / ADH references
  - linked acceptance criteria
  - validation and review evidence
```

Adapter-related assessments additionally trace to DEC-0036 alongside the
FEATURE-0011 references to DEC-0026 and RFC-0021 (Requirement 10).

## Controlled vocabularies

The standard defines these closed vocabularies. Any value outside a set is
invalid.

| Field | Controlled values | Requirement |
|---|---|---|
| Disposition | Reuse, Wrap, Extend, Build | 3 |
| Decision status | Proposed, Approved, Deferred, Rejected, Superseded | 3 |
| Adapter required | Yes, No | 5 |
| Allowed in current phase | Yes, No | 4 |
| Vendor-native types allowed | No, Approved exception reference | 4 |
| Replacement risk | Low, Medium, High | 9 |

### Disposition semantics (aligned with the Architecture Spine)

- Reuse: adopt a mature implementation, protocol, or standard substantially
  as provided; Sovrunn does not fork its core behavior; an adapter may
  still be required.
- Wrap: place a Sovrunn-owned contract around a mature capability without
  recreating the wrapped engine.
- Extend: add behavior through supported extension, composition, or
  compatible augmentation; a maintained fork requires separate approval and
  a maintenance assessment.
- Build: implement Sovrunn-owned differentiation or address the absence of
  an acceptable mature fit; requires rejection rationale for Reuse, Wrap,
  and Extend and defines long-term ownership.

### Decision status semantics

Only Approved, with recorded human approval, is authoritative and may
authorize progression to source implementation or implementation
execution. Progression between Kiro Requirements, Design, and Tasks stages
remains governed by the separate stage-approval tokens; an Approved
assessment status does not substitute for APPROVED_FOR_DESIGN,
APPROVED_FOR_TASKS, or the implementation authorization gate.

## Validation architecture

Validation is layered into three responsibilities with a strict boundary:
automation validates form; humans validate meaning.

```text
Assessment document
    |
    v
[Layer 1] Automated structural validation  -> required sections/fields,
                                               controlled values present
    |
    v
[Layer 2] Automated consistency validation  -> cross-field rules, invalid
                                                values, broken references
    |
    v
[Layer 3] Human semantic review  -> analysis quality, boundary correctness,
                                     sovereignty, Build justification,
                                     mitigation credibility, phase compliance
    |
    v
Feature gate decision (pass/fail) + recorded human approval
```

### Validation contract

The automated layers implement one deterministic, non-mutating validation
interface. The implementation approach is fixed here.
The repository-standard Bash entry point is
`scripts/reuse-assessment-check.sh`.
It may invoke Python 3 for deterministic Markdown parsing and validation,
following the existing repository automation pattern. No third-party
validation framework or runtime service is introduced. Tasks implement this
approved approach and do not choose another engine.

Inputs:

- repository root;
- feature identifier;
- assessment artifact paths;
- validation mode (strict for FEATURE-0011+, legacy for FEATURE-0001–0010).

The caller must not supply the authoritative canonical version. The
validator resolves the version itself.

Version-resolution flow:

- locate the fixed canonical file
  `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`;
- read `reuse_assessment_format_version` from that file;
- treat the following as validator configuration errors (exit code 2), not
  RA-C09 failures: canonical file missing; canonical version field missing;
  canonical version malformed; canonical file unreadable;
- only after a valid authoritative canonical version is resolved, read the
  declared aligned version from each required version-bearing artifact;
- apply RA-C09 to those target artifacts.

The canonical document itself never fails RA-C09; problems with the
canonical file are configuration errors (exit 2).

Two distinct mechanisms exist and must not be conflated.

Canonical-source reference (link only): the following supporting artifacts
link to `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` and are not
required to declare a numeric version unless they contain a version-bearing
example: Kiro prompt artifacts; Cursor prompt artifacts; reviewer prompt
artifacts; Feature Factory documents; templates; governance and policy
documents.

Numeric version declaration: the following required version-bearing
artifacts must declare a numeric version, using
`reuse_assessment_format_version: 1.0.0` in Markdown front matter or
`# reuse_assessment_format_version=1.0.0` in a shell or non-front-matter
automation artifact: the canonical Reuse Assessment Standard; each actual
reuse assessment; `scripts/reuse-assessment-check.sh`; and any fixture or
example intended to represent a complete versioned assessment.

RA-C09 is a validation failure (exit code 1) that applies only to required
target version-bearing artifacts: a required target assessment version
missing; the validator version marker missing; a required complete
fixture/example version missing; a target version malformed; or a target
version different from the canonical version. Link-only supporting
artifacts are not required to declare or duplicate the numeric version.

Diagnostic output: each issue includes at least a stable rule identifier;
validation layer; feature identifier; file path; section or field;
human-readable message; severity; and corrective guidance where
appropriate.

Behavior:

- diagnostics are emitted in a deterministic order;
- validation never mutates any document;
- a non-zero result is returned when structural or consistency validation
  fails;
- a successful result is returned when all applicable automated rules pass;
- FEATURE-0001 through FEATURE-0010 are explicitly exempt (legacy mode);
- FEATURE-0011 and later use strict mode;
- unknown or malformed feature identifiers fail safely (rejected, not
  skipped);
- human semantic approval is recorded separately and is never inferred by
  automation.

### Layer 1 — structural validation (Requirement 13)

Verifies that all required sections and fields exist and that controlled
values are drawn from their vocabularies: sections present; disposition in
vocabulary; decision status exists and uses the exact five values; adapter
decision explicit as Yes/No; replacement risk uses an approved value;
non-goals present; risk-mitigation fields present; architecture references
present; assessment owner, capability/decision-unit identity, selected
foundation or approach, and adapter/contract identifier present;
current-phase Yes/No; vendor-native-type rule enforced; linked acceptance
criteria present; each risk entry has both a preventive and a detection
control.

### Layer 2 — consistency validation (Requirement 14)

Applies cross-field and reference rules, each mapped to a stable RA-C*
identifier: adapter rationale mandatory for any Adapter required value
(RA-C01); adapter-related DEC-0036 reference (RA-C02); future-integration
placement (RA-C03); exact conceptual-example label (RA-C04); DEC/RFC/ADH
reference existence (RA-C05); mandatory adapter/contract identifier, which
must be the reserved literal `none` when Adapter required is No (RA-C06);
Build triple rejection rationale and ownership (RA-C07); risk triple
controls (RA-C08); required version-bearing artifact version present, well
formed, and matching (RA-C09); invalid feature identifier (RA-C10); missing
Phase 2 scope acknowledgement (RA-C11); duplicated canonical schema
definition (RA-C12); implementation attempted without Approved decision
status (RA-C13); and missing canonical-source reference (RA-C14). Every
Layer 2 failure maps to an RA-C* identifier in the range RA-C01 through
RA-C14.

### Consistency rule design (detailed)

Each rule has a stable identifier and is deterministic.

| Rule ID | Rule | Failure condition |
|---|---|---|
| RA-C01 | Adapter rationale is mandatory whether Adapter required is Yes or No | Adapter rationale absent for any value of Adapter required |
| RA-C02 | Adapter-related assessments must reference DEC-0036 | Adapter-related assessment omits a DEC-0036 reference |
| RA-C03 | Future integration content appears only under deferred work, explicit non-goals, or an identified future-phase section | Future integration content appears elsewhere |
| RA-C04 | A conceptual example carries the exact label | A conceptual example lacks "Conceptual example — illustrative only and outside execution scope" |
| RA-C05 | Referenced DEC, RFC, and ADH records must exist | Any referenced DEC/RFC/ADH record is not found |
| RA-C06 | Adapter or contract identifier is a mandatory field | Field absent |
| RA-C07 | Build requires all three rejection rationales and long-term ownership | Any of Reuse/Wrap/Extend rejection rationale or ownership missing |
| RA-C08 | Every risk requires preventive, detection, and corrective controls | Any risk missing one of the three |
| RA-C09 | Required target version-bearing artifact has a missing, malformed, or mismatched version (exit 1); applies only after a valid canonical version is resolved | A required target assessment, the validator version marker, or a complete fixture/example lacks, malforms, or mismatches `reuse_assessment_format_version` against the canonical version |
| RA-C10 | Invalid feature identifier | The identifier is missing or does not match `^FEATURE-[0-9]{4}$` |
| RA-C11 | Missing Phase 2 scope acknowledgement | A FEATURE-0011-and-later assessment does not acknowledge the applicable Phase 2 non-goals |
| RA-C12 | Duplicated canonical schema definition | An operational prompt, template, Feature Factory document, or policy artifact redefines the canonical field schema instead of referencing the canonical document |
| RA-C13 | Implementation without Approved decision status | Source implementation or implementation execution is attempted without an Approved decision status and recorded human approval |
| RA-C14 | Missing canonical-source reference | An operational prompt, reviewer prompt, Feature Factory document, template, governance document, or policy artifact required to align with FEATURE-0011 does not reference `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` |

Adapter or contract identifier when Adapter required is No: the field
remains mandatory and must carry the reserved literal value `none` (meaning
"no adapter or contract applies"). An empty or omitted value fails RA-C06;
`none` satisfies it without weakening the mandatory-field requirement.

### Structural rule identifiers and diagnostic behavior

Structural diagnostics (Layer 1) each have a stable identifier:

| Rule ID | Structural diagnostic |
|---|---|
| RA-S01 | Required section missing |
| RA-S02 | Required field missing |
| RA-S03 | Invalid disposition |
| RA-S04 | Invalid decision status |
| RA-S05 | Invalid adapter value |
| RA-S06 | Invalid phase value |
| RA-S07 | Invalid vendor-native-types value |
| RA-S08 | Invalid replacement-risk value |
| RA-S09 | Missing risk-control field |
| RA-S10 | Missing traceability field |

Every automated diagnostic (RA-S* structural and RA-C* consistency) carries
exactly one stable identifier.

Severity vocabulary:

- error: gate-blocking structural or consistency failure;
- warning: explicitly non-blocking information only.

Exit codes:

- 0: validation passed;
- 1: validation failures found;
- 2: usage, configuration, or internal validator error.

Diagnostic ordering: sort by file path, then section or field, then rule
identifier. Validation is non-mutating.

Deterministic detection boundaries:

- future-integration placement (RA-C03): RA-C03 and the strict feature gate
  reuse or invoke the existing scope-text validation logic in
  `scripts/phase2-scope-check.sh`. Its existing blocked-pattern and
  allowed-heading logic remains the authoritative source for Phase 2 scope
  phrase detection. `scripts/reuse-assessment-check.sh` must not create a
  second independent blocked-phrase list; it may consume the scope
  checker's result or be orchestrated alongside it by
  `scripts/feature-gate.sh`.
- duplicated canonical schema (RA-C12) and missing canonical-source
  reference (RA-C14): the canonical standard remains the only normative
  operational schema definition. RA-C12 detects prohibited duplicated
  schema definitions; RA-C14 detects a missing canonical-source link. Both
  apply to operational prompts, reviewer prompts, Feature Factory documents,
  templates, and governance and policy artifacts. RA-C12 scans for
  prohibited duplicated field-definition blocks while allowing canonical
  links and version markers; link-only artifacts do not need a numeric
  version declaration. Approved Kiro requirements, design, and tasks
  specifications remain excluded from operational duplicated-schema
  enforcement; they are historical specification artifacts, not competing
  canonical sources, but should reference the canonical artifact where
  required by their stage template.
- implementation without Approved decision status (RA-C13): orchestrated by
  `scripts/feature-gate.sh`, and it must consume the same authoritative
  changed-file set the gate uses rather than a second independent
  implementation. The implementation change set is the union of: committed
  feature-branch changes relative to the phase integration branch merge
  base; staged and unstaged tracked changes relative to HEAD; and untracked,
  non-ignored files. The phase integration branch defaults to
  `phase2-reuse-first-paas-fabric-foundation`; the existing `PHASE_BRANCH`
  environment variable or repository gate configuration may override the
  default. The conceptual collection flow is: resolve the merge base
  between HEAD and `PHASE_BRANCH`; collect committed feature changes from
  the merge base through HEAD; collect staged and unstaged tracked changes;
  collect untracked non-ignored files using Git; then normalize,
  de-duplicate, and sort paths deterministically. Implementation-attempt
  paths are files in the resulting set under `cmd/`, `internal/`, `pkg/`,
  `api/`, `scripts/`, or other executable source directories already
  recognized by the repository gate. Changes to documentation or
  `.kiro/specs/` alone do not constitute source implementation. When
  implementation-attempt paths exist for FEATURE-0011 or later, RA-C13
  verifies that the relevant reuse assessment has decision status Approved;
  that recorded human-approval evidence exists; that the evidence references
  an approved Architecture Decision Handoff or an approved
  assessment-review record; and that the evidence identifies the approver
  and approval date, directly or through the referenced approved artifact.
  The approved ADH may provide pre-implementation architecture approval;
  `docs/reviews/feature-gates/FEATURE-0011-approval-review.md` records the
  later final merge review; final merge-review evidence is not required
  before implementation begins; and assessment approval does not substitute
  for Kiro stage tokens or implementation authorization. If the phase
  branch, merge base, or required Git metadata cannot be resolved, RA-C13
  reports a configuration diagnostic and exits with code 2; it must not
  silently reduce the scan to working-tree changes or skip the check.

Automation does not claim semantic understanding; these are deterministic
textual and structural checks only.

### Layer 3 — human semantic review (Requirement 15)

Human reviewers own reuse-analysis quality, responsibility-boundary
correctness, sovereignty suitability, adapter adequacy, Build
justification, mitigation credibility, residual-risk acceptance, and phase
compliance. Automation must not approve architecture or select products.

## Feature gate integration (Requirement 18)

The strict feature gate runs Layer 1 and Layer 2 for FEATURE-0011 and later
features. Phase 1 legacy features (FEATURE-0001 through FEATURE-0010) are
exempt from strict Phase 2 section checks. Any failed structural or
consistency check fails the gate, and a final feature-gate review is
recorded before merge. The gate is evidence for review; it is not
architecture approval.

## Repository consistency and alignment (Requirement 16)

The canonical standard is the only definition of the format. Kiro
requirements/design/tasks prompts, Cursor prompts, reviewer prompts, and
Feature Factory documents align to it by reference and share one controlled
vocabulary. Existing incorrect or inconsistent architecture references are
corrected without changing unrelated accepted decisions. FEATURE-0011
traceability entries are recorded in the feature and decision traceability
matrices. The RA-C12 duplicated-schema check detects inconsistent
duplicated field-definition blocks, and the RA-C14 canonical-source-
reference check detects a required artifact that omits the canonical link,
so that drift is caught early.

## Versioning and alignment

- Format-version syntax: semantic versioning `MAJOR.MINOR.PATCH`.
- Where recorded: the `reuse_assessment_format_version` field in the front
  matter of `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`.
- Canonical-source linking (link only): operational prompts, policies,
  Feature Factory documents, reviewer artifacts, and templates reference the
  canonical file path and do not declare a numeric version unless they
  contain a complete version-bearing example.
- Numeric version declaration: only the canonical standard, each actual
  reuse assessment, `scripts/reuse-assessment-check.sh`, and complete
  versioned assessment fixtures or examples declare a numeric version.
- Drift detection: RA-C09 compares the numeric version only for required
  version-bearing artifacts. A link-only supporting artifact is checked for
  the canonical path reference and duplicated-schema avoidance (RA-C12), not
  for a numeric version declaration.
- Drift severity: a version mismatch on a required version-bearing artifact
  fails consistency validation.
- Update ordering when the canonical format changes: (1) update the
  canonical document and increment its version; (2) re-align supporting
  prompts and Feature Factory documents; (3) update validation rules and
  fixtures; (4) update traceability matrices.
- Duplicated schema definitions are prohibited: only the canonical document
  defines fields; the RA-C12 duplicated-schema check flags any operational
  artifact that redefines the schema instead of referencing it.

## Reassessment lifecycle (Requirement 17)

```text
Approved assessment
    |
    v
Reassessment trigger occurs
  (requirements | licensing | maintenance status |
   security posture | deployment context change)
    |
    v
Assessment revisited; decision status re-evaluated
    |
    +--> remains Approved  (re-affirmed)
    +--> becomes Superseded (new assessment via ADH)
    +--> becomes Deferred / Rejected (with rationale)
```

Every FEATURE-0011-and-later feature documents which architecture risks
apply, its preventive and detection controls, its corrective path, its
residual risk, and its reassessment triggers.

## Repository-level risk-control matrix (Requirement 17)

The standard defines the following repository-level risks and their
non-product-specific controls.

| Risk | Preventive control | Detection control | Corrective path | Residual-risk owner | Reassessment trigger |
|---|---|---|---|---|---|
| Incorrect classification | Capability-level disposition rules and semantics | Structural check on disposition vocabulary; human review | Reclassify via new assessment | Architecture owner | Capability or dependency change |
| Vendor-first architecture | Provider-neutrality field; no vendor selection | Human review of candidates | Remove vendor coupling; re-assess | Architecture owner | New candidate considered |
| Direct external coupling | Adapter-boundary requirement | RA-C01/RA-C02 checks; human review | Introduce adapter boundary | Architecture owner | Integration scope change |
| Wrapper responsibility expansion | Wrap semantics (no engine recreation) | Human review of boundary fields | Narrow wrapper scope | Architecture owner | Wrapper scope growth |
| Unjustified Build | Build rejection-rationale requirement | RA-C07 check; human review | Provide rationale or re-disposition | Architecture owner | Mature fit emerges |
| Adapter omission | Mandatory adapter fields | RA-C01/RA-C06 checks | Add adapter decision/identifier | Architecture owner | New external dependency |
| Phase 2 scope leakage | Current-phase Yes/No and non-goals | RA-C03 check; phase2-scope-check.sh | Move content to deferred/non-goals | Phase 2 scope owner | Scope boundary change |
| Missing replacement planning | Exit/migration boundary; replacement risk | Structural checks on those fields | Add exit boundary and replacement risk | Architecture owner | Component maturity change |
| Decision-status ambiguity | Single canonical decision-status field | RA-S04 decision-status check | Set exact status | Architecture owner | Status change |
| Template divergence | Single canonical source; align-by-reference | RA-C12 duplicated-schema check; RA-C14 canonical-source-reference check | Re-align to canonical version | Documentation owner | Canonical format change |
| Traceability errors | Mandatory traceability fields | RA-C05 reference-existence check | Correct references | Architecture owner | Referenced record change |
| Silent later-feature overrides | Change-control via new ADH | Human review of the new ADH, prior assessment, and traceability records | Require ADH; restore Approved assessment | Architecture owner | Earlier contract change |

## Change control and approval flow (Requirement 11)

Decision status is the single canonical status field; there is no separate
approval-status field. An assessment is authoritative, and source
implementation may begin, only when its decision status is Approved with
recorded human approval; rule RA-C13 enforces this implementation-
authorization control. A change to a disposition, responsibility
boundary, or mitigation plan after approval requires a new Architecture
Decision Handoff. Later features must not silently override an earlier
Approved assessment. A reuse assessment does not, by itself, approve a
vendor or an architecture change.

## Conceptual example handling (Requirement 12)

The standard permits conceptual examples only to aid human and AI
understanding. Every conceptual example must carry the exact label
"Conceptual example — illustrative only and outside execution scope" and
must not authorize implementation, product selection, runtime execution,
provider calls, plugin execution, provisioning, or detailed design of later
features.

> Conceptual example — illustrative only and outside execution scope

A future feature introducing a policy evaluation abstraction might record:
Build the Sovrunn-owned evaluation request/result contracts; Wrap a future
external policy engine behind an adapter; Reuse a mature policy language;
and defer selecting the first production engine. This illustrates
capability-level dispositions only; it selects no product and authorizes no
implementation.

## Correctness Properties

These properties must hold for any conforming reuse assessment and are the
basis for validation evidence.

Property 1: Every disposition value is a member of {Reuse, Wrap, Extend,
Build}.
**Validates: Requirements 3.1, 3.2, 3.4**

Property 2: Every decision status is a member of {Proposed, Approved,
Deferred, Rejected, Superseded}, with exactly one per capability
assessment.
**Validates: Requirements 3.5, 3.6, 3.7**

Property 3: The feature-level summary enumerates every significant
capability, so no distinct disposition is hidden by aggregation.
**Validates: Requirements 1.3, 2.4**

Property 4: If disposition is Build, rejection rationale exists for Reuse,
Wrap, and Extend, and long-term ownership is stated.
**Validates: Requirements 8.1, 8.2, 8.3, 8.4, 8.5**

Property 5: Adapter required is exactly Yes or No (an adapter rationale and
an adapter or contract identifier are always mandatory; see Properties 11
and 12).
**Validates: Requirements 5.1, 5.2**

Property 6: Allowed-in-current-phase is Yes or No, and vendor-native types
allowed is No or an Approved exception reference.
**Validates: Requirements 4.9, 4.10**

Property 7: Replacement risk is Low, Medium, or High.
**Validates: Requirements 9.6**

Property 8: Every listed risk has at least one preventive control, at least
one detection control, and a corrective path.
**Validates: Requirements 9.8**

Property 9: Decision status is the only status field, and only Approved
with recorded human approval is authoritative.
**Validates: Requirements 11.2, 11.3, 11.6**

Property 10: Every referenced DEC, RFC, and ADH record exists, and
adapter-related assessments reference DEC-0036 alongside DEC-0026 and
RFC-0021.
**Validates: Requirements 10.3, 10.4, 10.5, 10.6**

Property 11: An adapter rationale is present regardless of whether Adapter
required is Yes or No (rule RA-C01).
**Validates: Requirements 5.2, 5.3**

Property 12: The adapter or contract identifier is always present; when
Adapter required is No it is the reserved literal `none` (rule RA-C06).
**Validates: Requirements 4.8**

Property 13: Future-integration content appears only under deferred work,
explicit non-goals, or an identified future-phase section (rule RA-C03).
**Validates: Requirements 7.5**

Property 14: Every conceptual example carries the exact required label
(rule RA-C04).
**Validates: Requirements 12.2, 12.4**

Property 15: Every required version-bearing artifact declares a present,
well-formed `reuse_assessment_format_version` that matches the canonical
version (rule RA-C09).
**Validates: Requirements 16.2, 16.3**

## Error Handling

- Structural or consistency violations cause the feature gate to fail and
  the assessment to be classified non-compliant; the failing rule is
  identified for correction.
- A value outside a controlled vocabulary is reported as invalid rather
  than coerced.
- A missing or nonexistent DEC/RFC/ADH reference is reported as an invalid
  reference.
- Source implementation attempted without an Approved decision status is
  rejected as non-compliant (RA-C13, exit code 1).
- Canonical-configuration errors are distinct from validation failures: a
  missing, unreadable, or version-missing/malformed canonical file is a
  configuration error (exit code 2), not an RA-C09 failure. RA-C09 target
  failures are validation failures (exit code 1) and apply only after a
  valid canonical version is resolved.
- If the approved architecture is insufficient to complete an assessment
  without a new decision, work stops and reports the exact token
  `ARCHITECTURE_DECISION_REQUIRED` rather than inventing a classification.

## Testing Strategy

- Maintain fixtures of valid and invalid assessment documents covering each
  correctness property (Property 1–15), including boundary cases for every
  controlled vocabulary.
- Verify structural checks (Layer 1, rules RA-S01 through RA-S10) detect
  missing sections/fields and out-of-vocabulary values.
- Verify consistency checks (Layer 2) detect cross-field violations,
  broken references, duplicated-template drift, and unapproved
  implementation attempts.
- Verify each detailed consistency rule (RA-C01 through RA-C14), including:
  adapter rationale required for both Yes and No; adapter-related DEC-0036
  reference; future-integration placement; exact conceptual-example label;
  DEC/RFC/ADH existence; adapter/contract identifier including the reserved
  `none`; Build triple rejection rationale and ownership; risk triple
  controls; version-bearing-artifact version presence and match; invalid
  feature identifier; missing Phase 2 scope acknowledgement; duplicated
  canonical schema definition (RA-C12); implementation attempted without an
  Approved decision status; and missing canonical-source reference (RA-C14).
- Verify canonical-configuration errors exit 2 (canonical file missing,
  unreadable, or with a missing/malformed version) while RA-C09 target
  failures exit 1, and that the canonical document itself never fails
  RA-C09.
- Verify RA-C10 uses `^FEATURE-[0-9]{4}$` and that these values fail:
  `XFEATURE-0011`, `FEATURE-0011-extra`, `FEATURE-011`, `FEATURE-001A`.
- Verify RA-C13 change-set discovery detects an implementation file that
  is: already committed on the feature branch; an unstaged tracked change;
  a staged change; and an untracked non-ignored file; and that
  documentation-only feature changes are not treated as implementation.
- Verify RA-C13 approval logic: an implementation change with Approved
  status and valid human evidence passes; with Proposed status fails; with
  Approved status but missing human evidence fails.
- Verify RA-C13 fails safely: an unavailable phase branch, merge base, or
  required Git metadata yields a configuration diagnostic with exit code 2
  (not a silent skip and not reduced to working-tree changes).
- Verify the feature gate fails on any Layer 1/Layer 2 violation and exempts
  FEATURE-0001 through FEATURE-0010 from strict Phase 2 section checks.
- Human semantic review (Layer 3) is validated by checklist and recorded
  approval; it is intentionally not automated.

## Requirements traceability

| Requirement | Design coverage |
|---|---|
| 1 Canonical standard and feature-level summary | Canonical artifact architecture; Data Models (feature-level summary). |
| 2 Capability-level assessments | Data Models (capability-level assessment). |
| 3 Controlled dispositions and decision status | Controlled vocabularies; disposition and decision-status semantics. |
| 4 Descriptive and responsibility fields | Capability-level assessment identity, analysis, and boundary fields. |
| 5 Adapter boundary fields | Boundary fields; controlled vocabularies (Adapter required). |
| 6 Suitability considerations | Suitability fields. |
| 7 Phase impact, non-goals, exit boundary | Phase and scope fields. |
| 8 Build rejection rationale | Build justification fields. |
| 9 Risk-mitigation fields | Risk-mitigation fields; controlled vocabularies (Replacement risk). |
| 10 Traceability and validation evidence | Traceability fields; adapter DEC-0036 tracing. |
| 11 Decision status and change control | Change control and approval flow. |
| 12 Conceptual example labeling | Conceptual example handling. |
| 13 Automated structural validation | Validation architecture Layer 1; Validation contract. |
| 14 Automated consistency validation | Validation architecture Layer 2; Consistency rule design (RA-C01–RA-C14); Validation contract. |
| 15 Human semantic review | Validation architecture Layer 3. |
| 16 Repository consistency and canonical source | Canonical artifact architecture; Repository artifact map; Repository consistency and alignment; Versioning and alignment. |
| 17 Reassessment lifecycle and future-feature records | Reassessment lifecycle; Repository-level risk-control matrix. |
| 18 Strict feature-gate enforcement | Feature gate integration. |

## Feature identity and guardrails

- Feature: FEATURE-0011 — Reuse Assessment Standard.
- Phase: Phase 2 — Reuse-First PaaS Fabric Foundation.
- Approved feature-level disposition: Extend (preserved from requirements
  and ADH-2026-011).
- This design is documentation and governance only: no runtime resource,
  no vendor selection, no provider integration, no plugin execution, no
  persistence, no billing, no failover, no autonomous AI operations, and
  no Go production code.
- No FEATURE-0012 through FEATURE-0026 detailed design is included.
- No new architecture decision is introduced; all content derives from the
  approved requirements, ADH-2026-011, and the Phase 2 Architecture Spine.
- All Phase 2 non-goals are preserved.

## Verification

- Traceability: every requirement (1–18) maps to a design section above.
- Controlled vocabularies match the approved requirements exactly.
- The approved feature-level reuse summary (Extend) is preserved and not
  contradicted.
- Design remains within Phase 2 scope and introduces no runtime artifact.
