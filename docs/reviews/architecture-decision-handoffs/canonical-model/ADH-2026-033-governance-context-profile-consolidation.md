# ADH-2026-033: EffectiveGovernanceContext and Profile Consolidation

- **Status:** Architecture-owner approved; repository integration pending
- **Handoff ID:** ADH-2026-033
- **Date:** 2 August 2026
- **Classification:** Planned Phase 2 model replacement before implementation
- **Canonical decision:** FCM-ADR-014
- **Affected baseline:** FEATURE-0013 references; planned FEATURE-0018–0021

## Context

The Phase 2 plan separates GovernanceProfile, SecurityProfile, DataPlacementPolicy, CostGuardrail and EffectivePolicyContext. The final model needs one explainable resolution boundary that includes governance profiles, sovereignty, entitlement restrictions, approvals and exceptions. Proliferating top-level profile kinds would increase customer/provider complexity and create overlapping precedence rules.

## Decision

1. `GovernanceProfile` is the compositional operating-control definition.
2. Security, backup, placement, approval, audit, evidence and cost controls are typed sections or versioned referenced components within GovernanceProfile when independent reuse is justified.
3. `SovereigntyProfile` remains separate because it has specialized dimensions, policy interpretation, evidence and continuous assessment.
4. `ProfileAssignment` applies exact profile versions to supported scopes/resources.
5. `EffectiveGovernanceContext` replaces the planned `EffectivePolicyContext` name and resolves profiles, entitlements, restrictions, approved exceptions and exact versions into one immutable decision input.
6. Decision domains consume the context and must not independently traverse governance hierarchy.

## Alternatives rejected

- **Keep every planned profile as a top-level peer:** overlapping ownership and precedence.
- **Keep EffectivePolicyContext name:** too narrow once commercial entitlement and approved exception context are included.
- **Embed resolved controls in each decision:** duplicates resolver semantics and destroys reuse.

## Consequences

- FEATURE-0018–0021 must be rebaselined before implementation.
- FEATURE-0013 fields and documentation referencing EffectivePolicyContext require a compatible alpha terminology amendment without changing its typed-reference/provenance semantics.
- Cost accounting and production compliance engines remain later capabilities; profile fields do not imply those engines exist.

## Migration

No runtime objects exist for FEATURE-0018–0021, so planned names are replaced before implementation. FEATURE-0013 schema/profile registry references are migrated atomically and compatibility fixtures updated.

## Security and sovereignty

Allowed sets intersect, prohibitions accumulate, strongest minimum and shortest evidence age win, and unresolved conflict denies. Only qualified publishers can publish governance/sovereignty definitions. Exceptions cannot override non-exceptionable legal controls.

## Conformance evidence

- Provider, portfolio, Organization, Tenant and Project assignments resolve once.
- Lower scope narrows but cannot weaken mandatory controls.
- Entitlement restriction and approved exception appear with provenance.
- Stale/revoked context is rejected by DecisionRecord profiles.
- Decision evaluator cannot perform hidden hierarchy traversal.

## Reassessment triggers

Three independently governed profile families require distinct ownership, lifecycle and resolution semantics that cannot safely compose under GovernanceProfile.
