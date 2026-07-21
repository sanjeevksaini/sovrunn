---
doc_type: standard
title: Phase 2 Reuse Assessment Standard
status: draft
phase: 2
ai_load_priority: always
ai_summary: Canonical reuse assessment format, controlled vocabularies, validation rules, and risk-mitigation contract for FEATURE-0011 and later features.
---

# Phase 2 Reuse Assessment Standard

This document is the single canonical Sovrunn Reuse Assessment Standard.
It is the normative definition of the reuse assessment format, fields,
controlled vocabularies, validation rules, reassessment lifecycle, and
repository-level risk controls. Every other document must reference this
file rather than redefine the schema.

Format version: `1.0.0` (front matter `reuse_assessment_format_version`).

Document status remains **draft** until recorded human approval through the
FEATURE-0011 feature-review process. Approval of ADH-2026-011 authorizes
this work; it does not mark this canonical file Approved.

Controlling references:

- ADH-2026-011 (Approved)
- DEC-0026 — Reuse Before Build (Accepted)
- DEC-0036 — Adapter Boundaries Before External Integration (Accepted)
- RFC-0021 — Reuse-First Architecture
- PHASE2_ARCHITECTURE_SPINE Invariant A; contracts P2-C01, P2-C02, P2-C15

## 1. Purpose and applicability

Every FEATURE-0011-and-later Phase 2 feature must satisfy this contract
before source implementation begins. The standard makes reuse-before-build
decisions explicit, capability-level, auditable, and mitigable.

This standard is documentation, governance, and validation only. It does
not introduce a runtime `ReuseAssessment` resource, select a vendor, or
authorize production integration.

Applicability:

| Feature range | Mode |
|---|---|
| FEATURE-0011 and later | Strict — full structural and consistency validation |
| FEATURE-0001 through FEATURE-0010 | Legacy exemption from strict Phase 2 assessment checks |

## 2. Assessment model

### 2.1 Feature-level reuse summary (mandatory)

Every FEATURE-0011-and-later feature must include exactly one feature-level
reuse summary in its architecture contract, Kiro requirements, design,
tasks, and final review.

The summary enumerates each significant capability or decision unit with:

| Column | Meaning |
|---|---|
| Capability / decision unit | Identity of the assessed unit |
| Disposition | Exactly one of Reuse, Wrap, Extend, Build |
| Rationale | Concise justification |
| Decision status | Exactly one of Proposed, Approved, Deferred, Rejected, Superseded |
| Controlling reference | ADH, DEC, and/or RFC that controls the disposition |

Distinct dispositions must not be collapsed into a single feature-wide
label. Absence of the feature-level summary makes the feature non-compliant.

### 2.2 Capability-level assessment (zero or more)

A significant architectural capability is a decision unit whose disposition
affects ownership, coupling, or replacement cost. Each such unit may have
its own capability-level assessment. A single feature-level label must not
hide materially different component dispositions.

## 3. Controlled vocabularies

Any value outside a closed set is invalid.

| Field | Controlled values |
|---|---|
| Disposition | Reuse, Wrap, Extend, Build |
| Decision status | Proposed, Approved, Deferred, Rejected, Superseded |
| Adapter required | Yes, No |
| Allowed in current phase | Yes, No |
| Vendor-native types allowed | No, or an Approved exception reference |
| Replacement risk | Low, Medium, High |

### 3.1 Disposition semantics

| Disposition | Meaning |
|---|---|
| Reuse | Adopt a mature implementation, protocol, or standard substantially as provided; Sovrunn does not fork its core behavior; an adapter may still be required. |
| Wrap | Place a Sovrunn-owned contract around a mature capability without recreating the wrapped engine. |
| Extend | Add behavior through supported extension, composition, or compatible augmentation; a maintained fork requires separate approval and a maintenance assessment. |
| Build | Implement Sovrunn-owned differentiation or address the absence of an acceptable mature fit; requires rejection rationale for Reuse, Wrap, and Extend and defines long-term ownership. |

### 3.2 Decision status semantics

Decision status is the single canonical status field. There is no separate
approval-status field.

Only **Approved**, with recorded human approval, is authoritative and may
authorize progression to source implementation or implementation execution.
An assessment is non-authoritative until Approved status with recorded
human approval exists.

An Approved assessment status does not substitute for Kiro stage tokens
(`APPROVED_FOR_DESIGN`, `APPROVED_FOR_TASKS`) or the implementation
authorization gate (`APPROVED_FOR_CURSOR`).

A reuse assessment does not, by itself, approve a vendor or an architecture
change. Changes to disposition, responsibility boundary, or mitigation plan
after approval require a new Architecture Decision Handoff. Later features
must not silently override an earlier Approved assessment.

## 4. Capability-level mandatory fields

Normative field definitions occur only in this section. Supporting
documents must not redefine them.

### 4.1 Identity

| Field | Requirement |
|---|---|
| Feature identity | Must match `^FEATURE-[0-9]{4}$` and the active feature |
| Capability or decision-unit identity | Required |
| Assessment owner | Required |

### 4.2 Classification

| Field | Requirement |
|---|---|
| Disposition | Exactly one controlled value |
| Decision status | Exactly one controlled value |

### 4.3 Analysis

| Field | Requirement |
|---|---|
| Assessment scope | Required |
| Candidate category | Required |
| Mature candidates / applicable standards | Required |
| Relevant candidate strengths | Required |
| Material candidate constraints | Required |
| Rationale | Required |
| Selected foundation or approach | Required |

### 4.4 Boundary

| Field | Requirement |
|---|---|
| Sovrunn-owned responsibility | Required |
| Reused or external responsibility | Required |
| Data crossing the boundary | Required |
| Control crossing the boundary | Required |
| Adapter required | Exactly Yes or No |
| Adapter rationale | Required for both Yes and No |
| Adapter or contract identifier | Required; reserved literal `none` when Adapter required is No |
| Vendor-native types allowed | Exactly No, or an Approved exception reference |

Where disposition is Reuse, Wrap, or Extend against an external engine
expected to evolve or be replaced, the adapter decision must be justified
against DEC-0036.

### 4.5 Suitability

| Field | Requirement |
|---|---|
| Sovereignty and deployment fit | Required (include disconnected/air-gapped suitability where relevant) |
| Security and trust | Required |
| Operational and supportability | Required |
| Licensing and supply-chain | Required |
| Portability and provider-neutrality impact | Required |

Where a consideration does not apply, state non-applicability explicitly;
do not omit the field.

### 4.6 Phase and scope

| Field | Requirement |
|---|---|
| Allowed in current phase | Exactly Yes or No |
| Current-phase work | Required |
| Deferred work | Required |
| Explicit non-goals | Required |
| Exit or migration boundary | Required |
| Phase 2 non-goal acknowledgement | Required |

Future integration content may appear only under deferred work, explicit
non-goals, or an identified future-phase section.

### 4.7 Build justification (when disposition is Build)

| Field | Requirement |
|---|---|
| Why Reuse is insufficient | Required |
| Why Wrap is insufficient | Required |
| Why Extend is insufficient | Required |
| Protected Sovrunn differentiation and long-term ownership | Required |

### 4.8 Risk mitigation

| Field | Requirement |
|---|---|
| Applicable architecture risks | Required |
| Preventive controls | At least one per listed risk |
| Detection controls | At least one per listed risk |
| Corrective path | Required per listed risk |
| Residual risk | Required |
| Replacement risk | Exactly Low, Medium, or High |
| Reassessment triggers | Required |

A risk entry is incomplete unless it contains at least one preventive
control, at least one detection control, and a corrective path.

### 4.9 Traceability

| Field | Requirement |
|---|---|
| Related DEC / RFC / ADH references | Required; referenced records must exist |
| Linked acceptance criteria | Required |
| Validation and review evidence | Required |

FEATURE-0011 assessments must reference DEC-0026 and RFC-0021.
Adapter-related assessments must also reference DEC-0036.

### 4.10 Human-approval evidence (when decision status is Approved)

Recorded human-approval evidence must identify:

- the approving person or role;
- the approval date;
- the approved ADH or assessment-review reference;
- that the approval applies to the recorded disposition and responsibility
  boundary.

Evidence may resolve through an approved ADH that already contains the
approver and approval date.

## 5. Versioning and dependent artifacts

### 5.1 Format version

- Field name: `reuse_assessment_format_version`
- Initial value: `1.0.0`
- Syntax: semantic versioning `MAJOR.MINOR.PATCH`
- Increment rule: major for breaking field/vocabulary changes; minor for
  additive fields; patch for editorial clarification

Version markers:

- Markdown front matter: `reuse_assessment_format_version: 1.0.0`
- Shell / non-front-matter automation: `# reuse_assessment_format_version=1.0.0`

### 5.2 Version-bearing vs link-only artifacts

Required version-bearing artifacts (must declare a numeric version):

- this canonical standard;
- each actual reuse assessment;
- `scripts/reuse-assessment-check.sh`;
- complete fixtures/examples representing a versioned assessment.

Link-only supporting artifacts (reference the canonical path; do not
declare a numeric version unless they contain a complete version-bearing
example):

- Kiro, Cursor, and reviewer prompts;
- Feature Factory documents;
- templates;
- governance and policy documents.

### 5.3 Alignment rule

Dependent artifacts reference
`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` and must not duplicate or
redefine field definitions. When the canonical format changes: (1) update
this document and increment its version; (2) re-align supporting prompts
and Feature Factory documents; (3) update validation rules and fixtures;
(4) update traceability matrices.

## 6. Validation architecture

Validation has three layers. Automation validates form; humans validate
meaning. Automation must never approve architecture or select products.

```text
Assessment document
    |
    v
[Layer 1] Automated structural validation
    |
    v
[Layer 2] Automated consistency validation
    |
    v
[Layer 3] Human semantic review
    |
    v
Feature gate (Layers 1–2) + recorded human approval (Layer 3)
```

### 6.1 Validation contract

Entry point: `scripts/reuse-assessment-check.sh` (may invoke Python 3 for
deterministic Markdown parsing). No third-party validation framework.

Inputs:

- repository root;
- feature identifier;
- assessment artifact paths;
- validation mode (strict for FEATURE-0011+, legacy for FEATURE-0001–0010);
- optional gate-supplied normalized changed-file list (for RA-C13).

The validator resolves the canonical version from this file. Canonical file
missing, unreadable, or version missing/malformed is a configuration error
(exit 2), not RA-C09.

Diagnostics include: stable rule id, layer, feature id, file path,
section/field, message, severity (`error` | `warning`), corrective guidance.
Ordering: file path, then section/field, then rule id. Validation never
mutates inputs.

Exit codes:

| Code | Meaning |
|---|---|
| 0 | Validation passed |
| 1 | Validation failures found |
| 2 | Usage, configuration, or internal error |

Unknown or malformed feature identifiers fail safely (rejected, not
skipped). Human semantic approval is never inferred by automation.

### 6.2 Layer 1 — structural rules (RA-S01–RA-S10)

| Rule ID | Diagnostic |
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

### 6.3 Layer 2 — consistency rules (RA-C01–RA-C14)

| Rule ID | Rule |
|---|---|
| RA-C01 | Adapter rationale mandatory for Adapter required Yes and No |
| RA-C02 | Adapter-related assessments must reference DEC-0036 |
| RA-C03 | Future-integration content only under deferred work, non-goals, or future-phase headings; authoritative detection via `scripts/phase2-scope-check.sh` |
| RA-C04 | Conceptual examples must carry the exact label below |
| RA-C05 | Referenced DEC, RFC, and ADH records must exist |
| RA-C06 | Adapter or contract identifier mandatory; `none` when Adapter required is No |
| RA-C07 | Build requires Reuse/Wrap/Extend rejection rationale and long-term ownership |
| RA-C08 | Every risk requires preventive, detection, and corrective controls |
| RA-C09 | Required target version-bearing artifact version missing, malformed, or mismatched (exit 1); applies only after a valid canonical version is resolved |
| RA-C10 | Feature identifier must match `^FEATURE-[0-9]{4}$` |
| RA-C11 | Phase 2 scope acknowledgement present |
| RA-C12 | Operational artifacts must not redefine the canonical field schema |
| RA-C13 | Implementation-attempt paths for FEATURE-0011+ require Approved status and recorded human-approval evidence in the active feature's assessment |
| RA-C14 | Required operational artifacts must reference this canonical path |

RA-C12 and RA-C14 apply to operational prompts, reviewer prompts, Feature
Factory documents, templates, and governance/policy artifacts. Approved
Kiro requirements, design, and tasks specifications are excluded from
RA-C12 operational duplicated-schema enforcement.

### 6.4 Layer 3 — human semantic review

Human reviewers own:

- quality of the reuse analysis;
- correctness of responsibility boundaries;
- sovereignty suitability;
- adapter adequacy;
- Build justification;
- mitigation credibility;
- residual-risk acceptance;
- phase compliance.

### 6.5 Strict feature-gate enforcement

The strict feature gate runs Layers 1–2 for FEATURE-0011 and later.
FEATURE-0001 through FEATURE-0010 remain legacy-exempt. Any failed
structural or consistency check fails the gate. A final feature-gate review
must be recorded before merge. The gate is evidence for review; it is not
architecture approval.

Authoritative Git change-set discovery for RA-C13 lives in
`scripts/feature-gate.sh`. The validator consumes the gate-supplied list
and must not maintain a second discovery implementation.

Final merge-approval parsing requires the exact field:

```text
Final feature-review status: Approved
```

Other uses of the word Approved (assessment decision status, ADH
references, reuse-summary rows, Pending final review) must not satisfy
final merge approval.

## 7. Conceptual example labeling

Conceptual examples are permitted only to improve human and AI
understanding. Every conceptual example must be labeled exactly:

```text
Conceptual example — illustrative only and outside execution scope
```

Conceptual examples do not authorize implementation, product selection,
runtime execution, provider calls, plugin execution, provisioning, or
detailed design of later features.

> Conceptual example — illustrative only and outside execution scope

A later feature that introduces a policy evaluation abstraction might
record: Build the Sovrunn-owned evaluation request and result contracts;
Wrap a future external policy engine behind an adapter; Reuse a mature
policy language; and defer selecting the first production engine. This
illustrates capability-level dispositions only. It does not select a
product or authorize implementation.

## 8. Reassessment lifecycle

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

Every FEATURE-0011-and-later feature must document which architecture
risks apply, preventive controls, detection controls, a corrective path,
residual risk, and reassessment triggers.

## 9. Repository-level risk-control matrix

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
| Template divergence | Single canonical source; align-by-reference | RA-C12; RA-C14 | Re-align to canonical version | Documentation owner | Canonical format change |
| Traceability errors | Mandatory traceability fields | RA-C05 reference-existence check | Correct references | Architecture owner | Referenced record change |
| Silent later-feature overrides | Change-control via new ADH | Human review of the new ADH, prior assessment, and traceability records | Require ADH; restore Approved assessment | Architecture owner | Earlier contract change |

## 10. Non-goals

This standard shall not produce:

- a runtime `ReuseAssessment` API resource;
- selection, ranking, or approval of any production vendor or external
  engine;
- runtime provisioning, provider integration, or plugin execution;
- production persistence, billing, failover, or disaster recovery;
- autonomous AI operations;
- detailed design of FEATURE-0012 through FEATURE-0026;
- Go production code for FEATURE-0011;
- weakening of DEC-0026 or DEC-0036;
- changes to the Phase 2 feature sequence;
- a new architecture decision beyond ADH-2026-011.

If approved architecture sources are insufficient to supply a required
semantic decision, stop and report exactly:

```text
ARCHITECTURE_DECISION_REQUIRED
```

Do not invent missing values, approval evidence, owners, dates, boundaries,
or risk decisions.

## 11. Assessment outline (reference form)

Authors fill assessments using the fields defined in Section 4. The outline
below is structural guidance only; it does not redefine the schema.

```markdown
---
---

# FEATURE-NNNN — <title>

## Feature-level reuse summary

| Capability / decision unit | Disposition | Rationale | Decision status | Controlling reference |
|---|---|---|---|---|

## Capability assessment: <capability identity>

### Identity
### Classification
### Analysis
### Boundary
### Suitability
### Phase and scope
### Build justification  # only when disposition is Build
### Risk mitigation
### Traceability
### Human-approval evidence  # when decision status is Approved
```
