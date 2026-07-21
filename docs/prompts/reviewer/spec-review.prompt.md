Review {{STAGE}} for {{FEATURE_ID}}.

Target file:
{{TARGET_PATH}}

Review goals:
- Check alignment with Sovrunn architecture and Phase 1 scope.
- Check consistency with completed features.
- Check non-goals and scope boundaries.
- Check implementation feasibility for a solo developer using the active Go version standard.
- Check security/privacy implications.
- Check whether it introduces hidden dependencies or scope creep.
- Check whether it is ready to move to the next stage.

Return exactly one verdict:
- APPROVED
- APPROVED_WITH_MINOR_FIXES
- NEEDS_REVISION
- BLOCKED

Then provide summary, required fixes, optional improvements, next-step prompt for Kiro if revision is needed, and approval note if approved.

Document content:

```markdown
{{DOCUMENT_CONTENT}}
```

## Phase 2 Reuse and Drift Gates

Every FEATURE-0011-and-later feature must include a reuse assessment that
conforms to the canonical standard:

`docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

Do not duplicate or redefine the assessment field schema in this document.
Populate the feature-level reuse summary and capability-level assessments
using the canonical fields and controlled vocabularies.

Architecture drift checks:

- no provider-specific hardcoding in core,
- no Kubernetes-only assumptions in core,
- no PostgreSQL lifecycle logic in core placement engine,
- no custom policy engine embedded in handlers,
- no raw secret storage,
- no customer-facing IaaS leakage,
- explainable decision object,
- defined audit behavior,
- preserved adapter boundaries.
