---
doc_type: feature
id: FEATURE-001A
title: Reuse Assessment Standard
status: draft
phase: 2
reuse_assessment_format_version: 1.0.0
depends_on: []
ai_load_priority: feature
ai_summary: Canonical Reuse Assessment Standard governance contract for FEATURE-0011 and later Phase 2 features.
controlling_handoff: ADH-2026-011
---

# FEATURE-0011 — Reuse Assessment Standard

## Purpose

FEATURE-0011 establishes the single canonical Sovrunn Reuse Assessment
Standard as an Architecture Operating System governance contract. It
defines the mandatory reuse assessment format, controlled vocabularies,
validation rules, and risk-mitigation fields that every FEATURE-0011-and-
later Phase 2 feature must satisfy before source implementation begins.

This feature is documentation, governance, and validation work only. It
does not introduce a runtime resource, select a vendor, or add Go
production code.

Canonical standard:
`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

Controlling handoff: ADH-2026-011 (Approved).

## Acceptance criteria

1. One canonical, versioned Reuse Assessment Standard exists at
   `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` with
   `reuse_assessment_format_version: 1.0.0`.
2. Every FEATURE-0011-and-later feature includes a feature-level reuse
   summary and zero or more capability-level assessments conforming to the
   canonical field set.
3. Controlled dispositions are exactly Reuse, Wrap, Extend, and Build.
4. Decision status uses exactly Proposed, Approved, Deferred, Rejected, and
   Superseded; only Approved with recorded human approval is authoritative.
5. Automated structural (RA-S01–RA-S10) and consistency (RA-C01–RA-C14)
   validation exists and is enforced by the strict feature gate for
   FEATURE-0011+.
6. Human semantic review remains separate from automation.
7. Prompts, Feature Factory documents, templates, and gates align by
   reference to the canonical standard without duplicating the schema.
8. No runtime `ReuseAssessment` resource, vendor selection, provider
   integration, plugin execution, persistence, billing, failover, or
   autonomous AI behavior is introduced.
9. Phase 2 non-goals and the Phase 2 feature sequence remain unchanged.
10. Final merge requires recorded feature-gate review with
    `Final feature-review status: Approved`.

## Feature-level reuse summary

Feature identity: FEATURE-0011 — Reuse Assessment Standard.

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

## Capability assessment: Reuse Assessment Standard governance contract

### Identity

| Field | Value |
|---|---|
| Feature identity | FEATURE-001A |
| Capability or decision-unit identity | Reuse Assessment Standard governance contract |
| Assessment owner | Sovrunn Architecture Owner |

### Classification

| Field | Value |
|---|---|
| Disposition | Extend |
| Decision status | Approved |

### Analysis

| Field | Value |
|---|---|
| Assessment scope | Define the mandatory governance contract for capability-level reuse assessment, risk mitigation, validation, and feature-gate enforcement for FEATURE-0011 and later Phase 2 features. |
| Candidate category | Architecture governance and software-selection practice |
| Mature candidates / applicable standards | Architecture decision records; RFC review practices; software adoption and build-versus-buy assessments; risk-control registers; the existing Sovrunn reuse-before-build baseline (DEC-0026); adapter-boundary practice (DEC-0036); the existing draft Phase 2 reuse assessment format; RFC-0021 reuse-first architecture. |
| Relevant candidate strengths | Established reuse-before-build baseline; existing disposition vocabulary seed; mature ADR/RFC/risk-register patterns; approved ADH-2026-011 controlling decision. |
| Material candidate constraints | Draft format lacked versioning, mandatory mitigation fields, automated structural/consistency rules, and repository-alignment enforcement; feature-wide single-label classification is insufficient for multi-disposition features. |
| Rationale | Extending the existing draft and reuse-first baseline preserves accepted decisions while adding the mandatory capability-level contract, mitigation fields, and gate validation required by ADH-2026-011, without inventing a new governance mechanism. |
| Selected foundation or approach | Consolidate and version `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md` as the single canonical source; align dependent prompts, templates, and gates by reference; implement deterministic structural and consistency validation. |

### Boundary

| Field | Value |
|---|---|
| Sovrunn-owned responsibility | Four-disposition vocabulary; capability-level assessment rules; sovereign deployment criteria; provider-neutrality checks; adapter-boundary requirements; Phase 2 scope controls; future-feature mitigation requirements; architecture traceability; feature-gate structure. |
| Reused or external responsibility | General architecture-decision practices; software-selection practices; risk-management practices; mature component documentation and conformance evidence. |
| Data crossing the boundary | Assessment metadata, disposition values, risk-control records, DEC/RFC/ADH references, and validation diagnostics exchanged between governance documents and automation. |
| Control crossing the boundary | Human architecture approval and decision-status authority remain human-owned; automation enforces structure and consistency only. |
| Adapter required | No |
| Adapter rationale | This capability is a documentation-and-validation governance contract with no replaceable external engine integration; no adapter boundary is required under DEC-0036 for this assessment unit. |
| Adapter or contract identifier | none |
| Vendor-native types allowed | No |

### Suitability

| Field | Value |
|---|---|
| Sovereignty and deployment fit | Suitable for sovereign and disconnected environments because the standard is documentation and local repository validation only; no external runtime dependency is introduced. |
| Security and trust | No secrets, credentials, or trust boundaries are introduced; validation must not log secrets and must not mutate inputs. |
| Operational and supportability | Operates through existing repository scripts and feature-gate workflow; no new operational runtime service. |
| Licensing and supply-chain | Uses repository Markdown and Bash/Python standard-library tooling only; no third-party validation framework. |
| Portability and provider-neutrality impact | Provider-neutral by design; explicitly forbids vendor selection and provider-specific runtime coupling in Phase 2. |

### Phase and scope

| Field | Value |
|---|---|
| Allowed in current phase | Yes |
| Current-phase work | Canonical standard consolidation and versioning; FEATURE-0011 assessment instance; prompt/template/factory/gate alignment; validator and fixtures; pending feature-review evidence. |
| Deferred work | Application of the standard to FEATURE-0012 through FEATURE-0026 detailed design; any future production vendor evaluation; runtime assessment resources. |
| Explicit non-goals | Runtime `ReuseAssessment` API resource; production vendor or engine selection; provider provisioning or plugin execution; persistence, billing, failover, or disaster recovery; autonomous AI operations; Go production code for FEATURE-0011; detailed design of later Phase 2 features; weakening DEC-0026 or DEC-0036; changing the Phase 2 feature sequence. |
| Exit or migration boundary | If a later approved ADH supersedes this contract, publish a new format version, re-align dependents, and mark prior assessments Superseded where required. |
| Phase 2 non-goal acknowledgement | This assessment acknowledges and preserves all Phase 2 non-goals: no real provider provisioning, no full external-engine integrations, no runtime assessment resource, no Go production services for this feature, and no autonomous remediation. |

### Risk mitigation

#### Applicable architecture risks

1. Incorrect capability classification
2. Vendor-first architecture
3. Direct external coupling
4. Wrapper responsibility expansion
5. Unjustified custom build
6. Adapter omission
7. Phase-scope leakage
8. Missing replacement planning
9. Architecture-status ambiguity
10. Template divergence
11. Traceability errors
12. Silent later-feature overrides

#### Preventive controls

- Capability-level disposition rules and closed vocabularies in the canonical standard
- Provider-neutrality and non-goals fields; no vendor selection in Phase 2
- Mandatory adapter fields and DEC-0036 traceability for adapter-related assessments
- Wrap semantics that prohibit recreating wrapped engines
- Build triple rejection-rationale requirement
- Current-phase Yes/No, deferred work, and non-goals fields
- Exit/migration boundary and replacement-risk fields
- Single canonical decision-status field
- Align-by-reference rule for prompts, templates, and factory documents
- Mandatory DEC/RFC/ADH traceability fields
- Change-control requiring a new ADH for post-approval disposition or boundary changes

#### Detection controls

- RA-S01–RA-S10 structural validation
- RA-C01–RA-C14 consistency validation
- `scripts/phase2-scope-check.sh` for Phase 2 scope-phrase placement (RA-C03)
- Feature-gate orchestration of assessment path resolution and RA-C13
- Human semantic review (Layer 3)

#### Corrective path

- Fail the feature gate on structural or consistency violations
- Require correction of assessment fields or references
- Require a new Architecture Decision Handoff when disposition, responsibility boundary, or mitigation plan changes after approval
- Restore or reaffirm Approved assessments rather than allowing silent overrides

#### Residual risk

Residual risk is accepted that human reviewers may still misjudge semantic
quality of an otherwise structurally valid assessment. Automation cannot
eliminate that residual; Layer 3 review owns it.

#### Replacement risk

Medium

#### Reassessment triggers

- Change to requirements, licensing, maintenance status, security posture, or deployment context affecting this governance contract
- Canonical format version change
- Approved ADH that alters disposition vocabulary, mandatory fields, or gate enforcement
- Discovery that repository alignment has drifted from the canonical standard

### Traceability

| Field | Value |
|---|---|
| Related DEC / RFC / ADH references | ADH-2026-011; DEC-0026; DEC-0036; RFC-0021 |
| Linked acceptance criteria | Acceptance criteria 1–10 in this document; Requirements 1–18 in `.kiro/specs/reuse-assessment-standard/requirements.md` |
| Validation and review evidence | Canonical standard `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`; validator `scripts/reuse-assessment-check.sh`; feature gate `scripts/feature-gate.sh`; fixtures under `tests/reuse-assessment/`; pending final review at `docs/reviews/feature-gates/FEATURE-0011-approval-review.md` |

### Human-approval evidence

| Field | Value |
|---|---|
| Approving person or role | Sovrunn project owner |
| Approval date | 2026-07-21 |
| Approved ADH or assessment-review reference | ADH-2026-011 (`docs/reviews/architecture-decision-handoffs/ADH-2026-011-feature-0011-reuse-assessment-standard.md`) |
| Scope of approval | Approves the Extend disposition and the Sovrunn-owned versus reused/extended responsibility boundary recorded in this assessment, as stated in ADH-2026-011 Human approval |

This assessment decision status is **Approved** for the governance-contract
capability through ADH-2026-011. That assessment approval is distinct from
final FEATURE-0011 merge review, which remains a separate human action and
must not be inferred from this document.
