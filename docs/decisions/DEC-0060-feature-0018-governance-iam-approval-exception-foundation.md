---
doc_type: decision_record
title: FEATURE-0018 Governance, IAM, Approval and Exception Foundation
status: Accepted
phase: 2R
ai_load_priority: important
---

# DEC-0060: FEATURE-0018 Governance, IAM, Approval and Exception Foundation

## Status

Accepted

## Context

The Phase 2R canonical model names FEATURE-0018 but does not yet provide a
complete, deterministic enterprise authorization contract. Its preliminary
text omits access-only groups and retained approval/exception resources, binds
assignments directly to principals, and leaves scope, validity, lifecycle,
writers, decision provenance, access review, privileged access and provider-IAM
composition ambiguous. Those gaps could create multiple authorization
authorities and leak future or provider-native behavior into Sovrunn core.

## Decision

1. Adopt F18-RD-01 through F18-RD-24 in the content-bound ADH-2026-070 package
   as the sole detailed FEATURE-0018 semantic authority.
2. FEATURE-0018 owns `AccessGroup`, `Membership`, `RoleDefinition`,
   `RoleAssignment`, `PrivilegedAccessRequest`, `AccessReview`,
   `ApprovalPolicy`, `ApprovalRequest`, `ExceptionGrant`, and a v1-limited
   `GovernanceProfile`, plus the closed supporting values in F18-RD-02.
   `PrincipalRef` is an embedded reference to externally authenticated Human,
   Workload or System identity, not an identity lifecycle resource.
3. A `RoleAssignment` grants one immutable published RoleDefinition version to
   one `RoleHolderRef`—exactly one PrincipalRef or AccessGroupRef—at one
   registered canonical scope, optionally narrowed to one exact resource.
   Direct, auditable group membership is supported; external claims, nested
   groups and dynamic groups do not grant access.
4. Applicable grants combine by union. Scope/target applicability, policy
   guardrails, assurance, membership/responsibility, approval provenance,
   validity, review state, suspension/revocation and other dependencies combine
   by intersection. Any Deny, missing/ambiguous registration or failed required
   dependency denies authorization.
5. Every assignment has explicit `Standing` or `TimeBound` validity. Privileged,
   guest, JIT and break-glass assignments are TimeBound and expire without
   grace. Standing access is policy-limited, immediately revocable and
   periodically certified.
6. The RoleAssignment controller is the sole assignment publisher and sole
   specification, lifecycle, status and effect writer. Administrators and
   grant, approval, privileged-access and review controllers can submit exact
   immutable intents only. Delegated publication requires
   `roleassignment.grant` over the complete proposed reach and one independent
   per-action witness covering scope, target-binding mode, exact resource reach,
   validity and delegation; partial assignments cannot synthesize broader reach.
7. Approval is admission-time evidence. Approval publication and downstream
   effect publication are distinct atomic boundaries. Continuing authorization
   depends on the published effect's validity, dependencies, review state and
   runtime guardrails, not on continued freshness of the originating request.
8. FEATURE-0018 registers exactly `authorization-decision/v1`,
   `approval-decision/v1` and `exception-decision/v1` through FEATURE-0013 and
   registers its bounded AuditEvent taxonomy. FEATURE-0013 retains envelope,
   publication and audit-linkage ownership; AccessReview is not a fourth profile.
9. FEATURE-0018 consumes FEATURE-0017 evaluation through the bounded approved
   seam and consumes FEATURE-0016's existing ExecutionTarget actions with exact
   additive action-target and intrinsic-class metadata. It does not reinterpret
   either dependency.
10. FEATURE-0018 authorizes canonical Sovrunn control-plane actions only. Any
    native ExecutionTarget effect requires both a Sovrunn Allow and independent
    native-IAM authorization of a distinct least-privileged controller identity.
    Neither Allow substitutes for the other and neither layer overrides a Deny.
    FEATURE-0018 creates, translates, synchronizes and exposes no provider-native
    principal, group, role, policy, permission or credential.
11. `GovernanceProfile` v1 composes only FEATURE-0018-owned governance rules.
    FEATURE-0020 retains ProfileAssignment, inheritance, conflict/exception
    application and EffectiveGovernanceContext resolution.
12. Phase 2R provides deterministic, in-memory, zero-external-I/O fixtures and
    conformance. Real identity, policy, workflow, provider-IAM and execution
    integrations; persistence; UI; routes; and later-feature semantics remain
    excluded unless separately approved.
13. Normal JIT and break-glass activation share a non-weakenable floor of
    phishing-resistant AssuranceEvidence at AAL2 or higher and no more than
    fifteen minutes old. AccessReview excludes direct and indirect beneficiaries
    from retaining, replacing or advancing their own effective authority,
    including the responsible Human behind every direct Workload/System member
    of an AccessGroup that holds the reviewed authority.
14. Every AccessGroup Membership assignment-effect expansion is an atomic grant-
    producing boundary. The Membership controller derives an operation-local
    `MembershipEnabledGrantEnvelope`; a non-empty envelope requires the
    applicable Membership authority, `roleassignment.grant` over complete
    reach, and one non-synthesized per-action scope/target/resource/validity
    witness. Standing effects require Standing witnesses, trusted provisioners
    receive no bypass, and privilege-reducing Membership transitions remain
    available.
15. Approver, privileged-requester and reviewer eligibility use one closed
    non-resource `EligibilityRef` containing exactly one Human PrincipalRef.
    RoleDefinition, RoleAssignment, AccessGroupRef, Membership, claim and other
    relationships cannot satisfy or create eligibility in v1. Eligibility
    lists are non-empty, OR-composed, independently action-authorized, re-
    evaluated at their decision/effect boundaries, and fail closed on invalid
    or unresolved evidence. Exact policy/rule applicability and the separate
    canonical action alone determine scope and target reach.

## Consequences

Sovrunn gains a provider-neutral, enterprise-capable and auditable zero-trust
authorization foundation without making routine customer operations expose IAM
workflows. Administrators receive explicit role, group, approval, review and
exception resources; ordinary users normally see only the result or a bounded
request/approval journey when risk requires it.

The canonical model, catalog, glossary, dependency registries and traceability
must be reconciled as one change. Any new scope kind, containment rule, role
condition language, group model, action meaning, adapter, or future-feature
behavior requires separate change control.

## Alternatives Considered

- Direct-principal-only RBAC.
- Direct authorization from external IdP group claims.
- Provider-native IAM as the Sovrunn authorization model.
- A generic ABAC/condition language in v1.
- One combined workflow/governance resource.

All were rejected or deferred for the boundary, auditability, portability and
simplicity reasons recorded in ACR-2026-002 and ADH-2026-070.

## Supersedes

None. This decision narrows preliminary FEATURE-0018 wording and extends
DEC-0050 without changing FEATURE-0020 ownership of effective governance.

## Related

- ACR: ACR-2026-002
- ADH: ADH-2026-070 core plus Appendices A and B
- RFC: RFC-0012, RFC-0021, RFC-0022, RFC-0023 (preserved; no semantic replacement)
- Features: FEATURE-0011, FEATURE-0012, FEATURE-0013, FEATURE-0016,
  FEATURE-0017, FEATURE-0018, FEATURE-0020
- Architecture docs: `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md`

## Validation

Validate through the 36-item reconciliation ledger, canonical and dependency
drift checks, FEATURE-0011 reuse assessment, F18-RD-22 conformance inventory,
F18-RD-23 standards mapping, six journey-evidence records, strict MkDocs build,
applicable Structurizr check, independent security review, exact human approval
and the later FEATURE-0018 feature gate. Acceptance establishes architecture
authority only; requirements, design, tasks and implementation remain separate
human-gated stages.
