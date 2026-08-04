# Sovrunn Final Canonical Architecture — Approval Record

**Record ID:** ARCH-APPROVAL-2026-004
**Date:** 4 August 2026
**Requested by:** Sanjeev Kumar
**Human approval basis:** Explicit instruction in the current Codex task to complete approval of `ADH-2026-020` through `ADH-2026-041` and proceed through the Slice 0 charter
**Approval status:** Approved for controlled repository adoption
**Repository authority status:** Pending merge into the Sovrunn Git repository

## 1. Approved package

The human instruction approves the semantic decisions in the contiguous package:

```text
ADH-2026-020 through ADH-2026-041
```

The approved content is bound by:

- `sovrunn-canonical-adoption-content-manifest.sha256`;
- `sovrunn-finalized-data-model.md` version 1.7;
- `sovrunn-final-canonical-contract-catalog.md` version 1.7;
- `sovrunn-postgresql-reference-flow.md` version 1.2;
- `sovrunn-integrated-vertical-slice-delivery-plan.md` version 1.2.

Changes to a content-bound file after this approval require a new digest, classification and architecture-owner review. Editorial path changes during repository import are allowed only when semantic content and traceability are preserved.

## 2. Approval scope

Approval covers:

1. adoption of the final canonical semantic and contract model;
2. classification of the package as a coordinated correction, replacement, extension and new-decision set;
3. supersession of mandatory `ResourcePool` and `ProviderCapability` placement assumptions;
4. separation of `CloudPlatform`, `CloudProvider`, `CloudProviderParticipation` and `CloudEnrollment`;
5. coordinated alpha migration with no dual desired-state authority;
6. Phase 2 rebaseline from FEATURE-0015 onward;
7. Slice 0 chartering as a synthetic, side-effect-free integration proof;
8. Kiro validation and generation of repository architecture/spec updates under the normal staged workflow.

Approval does not authorize:

- direct production implementation from the canonical model alone;
- skipping Kiro requirements, design and tasks gates;
- modifying completed FEATURE-0001 through FEATURE-0014 history;
- executing the migration without backup, dry run, conformance and explicit activation;
- weakening human approval, security, sovereignty, evidence, audit or feature gates.

## 3. Decision-index treatment

Repository adoption shall add `DEC-0037` through `DEC-0058` as the compact accepted-decision mapping for `ADH-2026-020` through `ADH-2026-041`. Existing `DEC-0032` and `DEC-0033` become `Superseded` by the ExecutionTarget decision; completed implementation history remains unchanged.

`ADH-2026-042` is the single content-bound repository adoption envelope. It consolidates the approved package for downstream Kiro use in the same manner that `ADH-2026-017` consolidated earlier FEATURE-0013 handoffs. The individual ADH files remain decision provenance; Kiro must not reinterpret them as independent mutable authorities.

## 4. Conditions before repository authority

The package becomes repository authority only when all of these occur in one reviewed change:

- the consolidated handoff passes `make arch-handoff-check`;
- the Architecture Change Request and accepted DEC records are added;
- the decision index, glossary, current baseline and architecture version are updated;
- Phase 2 scope, sequence, architecture spine, roadmap and feature index agree;
- Kiro steering/context manifests reference the new controlling package;
- obsolete active statements are removed or explicitly marked superseded;
- documentation and architecture checks pass;
- the Git diff is reviewed and merged by the repository architecture owner.

## 5. Approval statement

> The final canonical architecture package is approved for controlled integration into the Sovrunn repository. This record is human-directed approval evidence, but it does not itself replace the repository as the source of truth.
