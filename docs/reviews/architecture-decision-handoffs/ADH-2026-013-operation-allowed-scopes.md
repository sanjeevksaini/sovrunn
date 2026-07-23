# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-013
- Date: 2026-07-23
- Source discussion: Sovrunn Architecture Governor ChatGPT Project
- Related feature: FEATURE-0012
- Related phase: Phase 2
- Author: Sanjeev Kumar
- Human approver: Sanjeev
- Approval status: Approved

## Decision title

Canonical Operation allowed scopes

## Summary

The canonical generic Operation contract supports all six formal Matrix B
governance scopes. Each individual Operation must use the same canonical
governance scope as its resolved target resource.

The allowed-scope declaration does not grant access. Target-kind constraints,
caller authorization, and no-existence-disclosure rules remain mandatory.

## Classification

Clarification

## Existing approved baseline

The approved FEATURE-0012 architecture classifies Operation as a
LongRunningOperation at the plugin-facing boundary and states that it uses
the target scope.

Matrix B defines these formal governance scopes:

- Platform
- Organization
- OrganizationUnit
- Tenant
- Project
- Provider

The approved baseline does not enumerate the exact machine-readable
x-sovrunn-allowed-scopes list for Operation.

Relevant references:

- docs/architecture/api-resource-standard.md
- ADH-2026-012
- F12-GOV-001
- F12-SCOPE-002
- F12-REF-001
- F12-FIXTURE-001
- F12-FIXTURE-002

## Decision or proposed decision

The canonical generic Operation schema declares exactly:

    x-sovrunn-allowed-scopes:
      - Platform
      - Organization
      - OrganizationUnit
      - Tenant
      - Project
      - Provider

Every individual Operation MUST use the same canonical primary governance
scope as its resolved target resource.

For a platform-scoped target:

- Operation.scopeRef is canonically absent or nil.
- PlatformScopeUID is used for scoped identity.

For every non-platform target:

- Operation.scopeRef identifies the target resource's governance scope.
- Scope identity is resolved by UID.
- An Operation whose scope differs from its target's resolved scope is
  rejected with a stable validation code and RFC 6901 JSON Pointer path.

Operation.ownerRef MAY identify the target as lifecycle containment, but it
MUST NOT replace scopeRef or act as a security or governance scope.

The six-value allowed-scope list does not authorize an Operation.
Target-kind constraints, caller authorization, and no-existence-disclosure
rules remain mandatory.

## Rationale

Operation is a generic LongRunningOperation grammar used to represent
actions against platform, organization, tenant, project, and provider-domain
resources.

Restricting the generic contract to only one domain would require separate
Operation grammars or a later compatibility-impacting expansion.

Allowing all six formal scopes preserves the generic contract. The mandatory
target-scope equality invariant prevents callers from placing an Operation
in an unrelated governance scope.

## Reuse-before-build assessment

- Decision: Extend
- Mature reusable options considered:
  - Kubernetes-style operation and condition conventions
  - JSON Schema 2020-12 annotations
  - typed resource references
- Sovrunn-owned responsibility:
  - governance-scope vocabulary
  - Operation-to-target scope equality
  - canonical platform-scope representation
  - target-kind constraints
  - authorization and no-existence-disclosure rules
- Non-goals:
  - operation execution
  - workflow orchestration
  - persistence
  - provider execution
  - plugin execution
  - policy-engine implementation

## Phase impact

- Current phase allowed? Yes
- If not current phase, target phase: N/A
- Current phase boundary impact:
  - Clarifies a machine-readable contract value.
  - Does not add runtime behavior.
  - Does not implement operation execution.
  - Preserves F12-IMPL-001 and F12-IMPL-002.

## Conflict check

- Conflicts with accepted DEC/RFC? No
- Conflicting decisions, if any:
  - None identified
- Resolution required:
  - Update architecture doc
  - Update Kiro design.md
  - Update Kiro tasks.md
  - Update relevant traceability
- Architecture Change Request, DEC, or RFC required:
  - No. This is a bounded clarification of an already approved contract.

## Required action

- Update architecture doc
- Update Kiro design.md
- Update Kiro tasks.md
- Update Operation fixture and traceability references
- Do not modify requirements.md
- Do not implement source code while applying this handoff

## Impacted files

- docs/architecture/api-resource-standard.md
- .kiro/specs/api-resource-naming-status-and-validation-standard/design.md
- .kiro/specs/api-resource-naming-status-and-validation-standard/tasks.md

## Impacted features

- FEATURE-0012:
  - Resolves the canonical Operation allowed-scope annotation.
- FEATURE-0013:
  - Later adopts the clarified Operation grammar without adding runtime
    behavior to FEATURE-0012.

## Acceptance criteria for Kiro update

- [ ] Handoff validated against Architecture Operating System files
- [ ] Operation declares all six exact Matrix B scopes
- [ ] Operation.scopeRef equals its resolved target scope
- [ ] Platform-scoped Operation uses canonical nil scopeRef
- [ ] PlatformScopeUID is used for platform-scoped identity
- [ ] ownerRef is not treated as a governance scope
- [ ] Target/scope mismatch receives a stable code and JSON Pointer
- [ ] The six-value list is not treated as an authorization grant
- [ ] No runtime or later-feature behavior is introduced
- [ ] requirements.md remains unchanged
- [ ] design.md records this handoff as a controlling input
- [ ] tasks.md contains no Operation-scope ARCHITECTURE_DECISION_REQUIRED marker
- [ ] Traceability remains complete

## Explicit instructions to Kiro

- Apply only the clarification recorded in this handoff.
- Do not narrow or broaden the six-value scope list.
- Do not treat the allowed-scope list as an authorization grant.
- Preserve target-kind constraints and caller authorization.
- Preserve canonical nil scopeRef for platform-scoped resources.
- Do not treat ownerRef as a security or governance scope.
- Do not introduce operation execution, persistence, workflow, provider,
  plugin, or policy-engine behavior.
- Do not modify requirements.md.
- Do not modify Go source while applying the architecture and spec updates.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-07-23
- Notes:
  - Approval authorizes only the exact six-scope declaration and the
    Operation-to-target scope equality invariant.

This Architecture Decision Handoff is ready for Kiro validation and repo update only after human approval.
