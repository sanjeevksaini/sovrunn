# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-018
- Date: 2026-07-29
- Source discussion: Sovrunn Architecture Governor ChatGPT Project — Feature 0014 Architecture
- Related feature: FEATURE-0014
- Related phase: Phase 2
- Author: Codex draft for architecture-owner review
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

FEATURE-0014 provider-neutral substrate topology boundary

## Summary

Define exactly five provider-neutral substrate resources—`Provider`, `ProviderLocation`, `ProviderDatacenter`, `DatacenterFailureDomain`, and `InfrastructureStack`—in the hierarchy shown by the approved sketches. Reuse the existing `Organization` resource as the optional supply owner above `Provider` and reuse FEATURE-0012 resource grammar. Classify FEATURE-0014 as `NOT_APPLICABLE` under the approved FEATURE-0013 downstream-adoption contract because it defines topology contracts only. Defer resource pools and capabilities to FEATURE-0015 and all adapter interfaces and provider integration to FEATURE-0016.

## Classification

New decision

## Existing approved baseline

The approved Phase 2 spine assigns provider-neutral hierarchy to FEATURE-0014, `ResourcePool` and `ProviderCapability` to FEATURE-0015, and adapter contracts to FEATURE-0016. It requires later features to reuse FEATURE-0012 resource grammar and FEATURE-0013 decision/accountability contracts without silent redefinition.

Relevant baseline references:

- `docs/phase2/PHASE2_EXECUTION_STRATEGY.md`
- `docs/phase2/PHASE2_ARCHITECTURE_SPINE.md`
- `docs/phase2/PHASE2_SCOPE.md`
- `docs/architecture/api-resource-standard.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-012-feature-0012-api-resource-standard.md`
- `docs/architecture/FEATURE-0013-decision-record-and-auditevent-standard.md`
- `docs/reviews/architecture-decision-handoffs/ADH-2026-017-feature-0013-consolidated-architecture.md`
- attached FEATURE-0014 topology sketch
- attached supply-owner sketch showing one owner organization governing multiple provider operators

The authoritative local repository records FEATURE-0013 as implemented and merged through PR #15 on 2026-07-29. `ADH-2026-017` is its sole controlling handoff; `ADH-2026-014`, `ADH-2026-015`, and `ADH-2026-016` are superseded downstream inputs retained only for provenance.

## Decision or proposed decision

1. FEATURE-0014 owns exactly `Provider`, `ProviderLocation`, `ProviderDatacenter`, `DatacenterFailureDomain`, and `InfrastructureStack`.
2. `ProviderLocation` is the only permitted kind and term for the location level. No alias resource, alias field, or alternative vocabulary is introduced.
3. The canonical contextual and containment hierarchy is:

   ```text
   Organization (existing; optional supply owner)
     -> Provider
       -> ProviderLocation
         -> ProviderDatacenter
           -> DatacenterFailureDomain
             -> InfrastructureStack
   ```

4. `Organization` is not owned or redefined by FEATURE-0014 and is not a sixth FEATURE-0014 resource.
5. Provider may be Platform-scoped or scoped to exactly one existing `Organization`, as already allowed by FEATURE-0012. The Provider's primary `scopeRef` identifies the supply owner or administrative authority.
6. A distinct owner/operator case is represented as `Organization: OwnerOrganization-A` with separately scoped `Provider: ProviderOperator-A` and `Provider: ProviderOperator-B` resources.
7. A same-party case is represented by role-specific resources such as `Organization: ProviderOperator-A` with `Provider: ProviderOperator-A` scoped to it; matching real-world identity does not merge the resource roles.
8. FEATURE-0012 `ownerRef` must not represent the supply-owner relationship. It remains reserved for lifecycle containment.
9. Each topology child has exactly one immediate parent. Each Provider descendant is immutably Provider-scoped, and its immediate parent must resolve within the same Provider scope.
10. `scopeRef` expresses FEATURE-0012 governance/security scope. Typed immediate-parent references express topology. Neither substitutes for the other.
11. All five resources use FEATURE-0012 canonical type identity, metadata, `ManagedResource` profile, operator-facing boundary, typed-reference grammar, ownership, status/conditions, validation, errors, concurrency, and conformance rules unchanged.
12. FEATURE-0014 records topology and minimal catalogue validity only. It does not model capacity, capabilities, placement eligibility, provider health discovery, or execution.
13. `InfrastructureStack` is a normalized substrate boundary within one failure domain. It is not an adapter, plugin, cluster, resource pool, or capability declaration.
14. FEATURE-0014 declares FEATURE-0013 adoption `NOT_APPLICABLE`: it produces or consumes none of the governed concepts listed by FEATURE-0013 section 28.2 and defines none of their semantics.
15. Provider-native IDs, SDK objects, endpoints, credentials, and secrets remain outside core resource contracts.
16. Physical/topological containment never implies network connectivity. Failure domains may or may not be connected within one datacenter, across datacenters, across ProviderLocations, or across Providers.
17. Absence of an explicit connectivity assertion means unknown, not connected and not disconnected. FEATURE-0014 introduces no connectivity boolean, link, route, peer graph, reachability status, latency, bandwidth, trust-zone, or health contract.
18. Any later placement, resilience, or data-movement behavior must use an explicitly owned and validated connectivity contract rather than infer connectivity from common ancestry or proximity. FEATURE-0053 is the current roadmap owner of `NetworkConnectivityProfile` foundation.
19. Registered parents may temporarily have zero children to support ordered onboarding, maintenance, decommissioning, and deletion. Zero children means topology-incomplete; it never permits a child to skip its required immediate parent.
20. A complete topology path contains all five FEATURE-0014 levels. Only an unbroken path ending in an InfrastructureStack may be reported as topology-complete; topology completeness is not capability, capacity, placement eligibility, or service readiness.
21. Every InfrastructureStack is a unique deployed stack instance identified by its immutable FEATURE-0012 `metadata.uid` and scoped resource identity. Multiple stacks may use the same technology without sharing identity.
22. Apache CloudStack, Cloud Foundry, Red Hat OpenShift, AWS IaaS, and Azure IaaS are illustrative associated technologies or environments only. They are not resource kinds, native core identities, or a closed compatibility enum.
23. Each InfrastructureStack belongs to exactly one `DatacenterFailureDomain` and must never span multiple failure domains. A deployment across failure domains is represented by distinct InfrastructureStack resources with distinct UIDs.

## Rationale

The existing `Organization` resource already provides the required owner/governance boundary. Reusing it avoids a duplicate owner model and directly supports both a multi-provider owner and a self-owned provider. The five FEATURE-0014 kinds preserve the supplied provider topology while creating the smallest provider-neutral substrate vocabulary needed by FEATURE-0015. Separating Organization scope, Provider scope, physical containment, and network connectivity preserves FEATURE-0012 authorization semantics and prevents placement from assuming reachability based on geography. Keeping capacity, capability, connectivity, adapters, and audit payloads outside the feature prevents dependency inversion and avoids pre-designing FEATURE-0015, FEATURE-0016, FEATURE-0053, or redefining FEATURE-0013.

## Reuse-before-build assessment

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - FEATURE-0012 API and resource standard
  - approved FEATURE-0013 consolidated architecture and `ADH-2026-017` as a preserved closed boundary
  - existing Phase 1 `Organization` resource
  - ISO country/jurisdiction codes where normalized geography requires them
  - common industry topology concepts: provider, provider location, datacenter, failure domain, infrastructure stack
- Sovrunn-owned responsibility summary:
  - stable provider-neutral resource semantics;
  - role separation between supply-owner Organization and infrastructure-operator Provider;
  - hierarchy and cardinality;
  - Provider scope and same-provider parent invariants;
  - typed topology references;
  - boundary classification and conformance.
- Non-goals summary:
  - ResourcePool and ProviderCapability;
  - a new Owner, ProviderOwner, or SupplyOwner resource;
  - network connectivity, routing, peering, latency, bandwidth, reachability, trust, or health contracts;
  - adapter interfaces and provider integration;
  - decision/audit/scope contract design;
  - capacity, placement, policy, runtime, plugins, operations, provisioning, persistence, billing, failover, or disaster recovery.

## Phase impact

- Current phase allowed? Yes
- If not current phase, target phase: N/A
- Current phase boundary impact:
  - Establishes the first half of the Wave B provider-neutral substrate model.
  - Supplies topology context for FEATURE-0015 without implementing FEATURE-0015 behavior.
  - Does not change the Phase 2 sequence or end-state.

## Conflict check

- Conflicts with accepted DEC/RFC? No.
- Conflicting decisions, if any:
  - None identified.
- Resolution required:
  - Architecture-owner approval.

## Required action

- Update architecture doc.
- Add the approved ADH under `docs/reviews/architecture-decision-handoffs/`.
- Replace the former stack-kind name with `InfrastructureStack` across current architecture, RFC, glossary, feature-index, phase-sequence, roadmap, and context references as one atomic terminology migration.
- Update current architecture/decision summaries and traceability only after approval.
- Initialize the Kiro spec using the canonical FEATURE-0014 slug only after approval.
- Generate Kiro `requirements.md` only; do not generate design, tasks, or implementation in the same step.

## Impacted files

- `docs/architecture/provider-neutral-resource-model.md` (new)
- `docs/reviews/architecture-decision-handoffs/ADH-2026-018-feature-0014-provider-neutral-resource-model.md` (new)
- `docs/context/CURRENT_ARCHITECTURE_BASELINE.md`
- `docs/context/CURRENT_DECISION_SUMMARY.md`
- `docs/context/CURRENT_PHASE_CONTEXT.md`
- `docs/glossary.md`
- `docs/phase2/PHASE2_FEATURE_SEQUENCE.md`
- `docs/roadmap/SOVRUNN_FEATURE_ROADMAP.md`
- `docs/architecture/development-phases.md`
- `docs/rfc/RFC-0024-provider-neutral-resource-model.md`
- `docs/features/FEATURE_INDEX.md`
- `docs/features/FEATURE-0014-provider-neutral-resource-model.md` (new)
- `docs/traceability/FEATURE_TRACEABILITY_MATRIX.md`
- `docs/reviews/reuse-assessments/FEATURE-0014-approval-evidence.md` (new)
- `docs/prompts/kiro/requirements.prompt.md`
- `docs/prompts/kiro/design.prompt.md`
- `docs/prompts/kiro/tasks.prompt.md`
- `docs/prompts/reviewer/spec-review.prompt.md`
- `docs/prompts/reviewer/approval-review.prompt.md`
- `scripts/feature-0014-boundary-check.py` (new)
- `scripts/kiro-stage.sh`
- `scripts/reviewer-stage.sh`
- `scripts/feature-gate.sh`
- `scripts/render-prompt.py`
- `Makefile`
- `.automation/features/FEATURE-0014.yaml` (new)
- `.automation/state/FEATURE-0014.json` (new)
- `.kiro/specs/provider-neutral-resource-model/.config.kiro` (new)
- `.kiro/specs/provider-neutral-resource-model/requirements.md` (later, only after approval)

## Impacted features

- FEATURE-0012: Consumed unchanged; no baseline modification.
- FEATURE-0013: Adoption is `NOT_APPLICABLE`; the approved closed boundary remains unchanged.
- FEATURE-0014: Establishes its complete architecture boundary and hierarchy.
- FEATURE-0015: Receives topology references later; no ResourcePool or ProviderCapability is defined here.
- FEATURE-0016: Receives integration context later; no adapter contract is defined here.
- FEATURE-0053: Remains the roadmap owner of `NetworkConnectivityProfile`; FEATURE-0014 supplies topology references only and makes no connectivity assertion.

## Acceptance criteria for Kiro update

- [ ] Handoff validated against current Architecture Operating System files.
- [ ] Approved FEATURE-0013 architecture and sole controlling `ADH-2026-017` are cited.
- [ ] Exactly five FEATURE-0014 resource kinds are present.
- [ ] No active architecture, RFC, glossary, feature-index, phase-sequence, roadmap, context, schema, or specification retains the superseded stack-kind name.
- [ ] Existing `Organization` is reused as the optional supply owner and is not redefined.
- [ ] Provider primary scope is exactly Platform or Organization as established by FEATURE-0012.
- [ ] Supply ownership uses Provider `scopeRef`; `ownerRef` is not used as governance/security scope.
- [ ] Both distinct-owner/operator and same-party owner/operator scenarios are covered by generic conformance examples.
- [ ] Same-datacenter, cross-datacenter, cross-location, and cross-provider ancestry imply neither connectivity nor isolation.
- [ ] Missing connectivity information is treated as unknown.
- [ ] No connectivity boolean, link, route, peer graph, reachability status, latency, bandwidth, trust-zone, or network-health field is introduced.
- [ ] `ProviderLocation` is the only permitted location-level kind and term.
- [ ] The hierarchy and one-parent cardinality match the attached sketch.
- [ ] Empty parents are accepted only as topology-incomplete registered resources.
- [ ] Every skipped-level relationship is rejected.
- [ ] A complete topology path contains all five FEATURE-0014 levels.
- [ ] Every InfrastructureStack has a distinct immutable UID even when multiple stacks use the same technology.
- [ ] Every InfrastructureStack has exactly one immutable failure-domain parent and cannot span failure domains.
- [ ] All descendants are Provider-scoped and immediate parents resolve in the same Provider scope.
- [ ] FEATURE-0012 grammar is reused without a new profile, scope kind, status grammar, or reference grammar.
- [ ] FEATURE-0013 adoption is declared `NOT_APPLICABLE` with the section 28.2 rationale.
- [ ] No FEATURE-0013 governed concept or closed vocabulary is produced, consumed, or redefined.
- [ ] No ResourcePool or ProviderCapability schema, semantics, requirement, or behavior is introduced.
- [ ] No adapter interface, provider client, SDK type, credential, discovery, plugin, or execution behavior is introduced.
- [ ] Provider-native identifiers do not become core identity or reference fields.
- [ ] Status contains current facts only; historical accountability remains with FEATURE-0013.
- [ ] Non-goals, architecture drift checks, security, sovereignty, observability, and reassessment triggers are explicit.
- [ ] Every normative requirement traces to at least one `F14-AD-*` decision or explicitly inherited FEATURE-0012 contract.
- [ ] Design introduces no uncited normative behavior and maps every requirement to a conforming representation.
- [ ] Tasks trace only to approved design and requirements and resolve no semantic choice.
- [ ] The closed decision register and stage-ownership matrix remain controlling architecture rather than being copied into a competing downstream source.
- [ ] The earlier-/adjacent-feature compatibility matrix is preserved and every inherited contract retains its canonical owner.
- [ ] Requirements generation loads only the approved task-specific context allowlist plus repository-mandated global context.
- [ ] FEATURE-0013 predecessor handoffs and superseded FEATURE-0014 text are excluded from semantic inputs.
- [ ] Source precedence is applied without allowing a lower-precedence summary, RFC, roadmap entry, example, or sketch to override normative architecture.
- [ ] Every normative requirement appears in the single-owner overlap ledger with exactly one owning feature and upstream decision/contract citation.
- [ ] A missing owner, duplicate owner, disallowed input, or semantic source conflict stops generation with `ARCHITECTURE_DECISION_REQUIRED`.
- [ ] Kiro follows the requirements → human approval → design → human approval → tasks stage state machine without combined or speculative generation.
- [ ] Each stage writes only its single allowlisted file and does not modify approved earlier-stage or architecture artifacts.
- [ ] Requirements trace all 21 closed decisions and all 30 risks without implementation choices or unresolved placeholders.
- [ ] Design introduces no new normative requirement and chooses only explicitly delegated representations.
- [ ] Tasks trace to approved design and requirements, include all required evidence work, and leave no semantic choice to implementation.
- [ ] Deterministic completeness checks fail the stage for missing traceability, forbidden concepts, aliases, unresolved markers, or cross-stage file changes.
- [ ] Examples and diagrams remain explanatory and non-normative.
- [ ] All architecture section 16.8 repository-readiness checks pass before requirements generation.
- [ ] Active architecture, RFC, glossary, feature-index, phase-sequence, roadmap, and context sources contain one consistent FEATURE-0014 vocabulary and boundary.
- [ ] Kiro and reviewer context manifests explicitly contain the approved FEATURE-0014 architecture and ADH and exclude conflicting/superseded inputs.
- [ ] Generic prompts and review gates are applicability-aware and do not require FEATURE-0014 to invent decision, audit, operation, or adjacent-feature behavior.
- [ ] A FEATURE-0014 automated boundary validator checks all decisions, risks, owners, normalization keys, forbidden semantics, context manifests, and stage outputs.
- [ ] The requirements normalization ledger has one normative home per semantic key and no unreviewed duplicate or conflicting requirement.
- [ ] Architecture status is not advanced to `READY_FOR_REQUIREMENTS` until repository alignment and tooling preflight pass.
- [ ] Any unresolved semantic choice produces `ARCHITECTURE_DECISION_REQUIRED` and stops the entire current generation stage.
- [ ] `F14-R01` through `F14-R30` each map to closed `F14-AD-*` decisions, mitigation, evidence, target residual, owner, and reassessment trigger.
- [ ] Requirements preserve every applicable risk mitigation and evidence obligation without accepting or downgrading residual risk.
- [ ] Design maps every risk control to a concrete component, schema, validation, fixture, or review mechanism.
- [ ] Tasks implement only approved risk controls and evidence; no task accepts, closes, or reclassifies risk.
- [ ] Residual-risk acceptance remains pending implementation evidence and human semantic review.
- [ ] Only `requirements.md` is generated after human approval.

## Explicit instructions to Kiro

- Treat the approved architecture document and this handoff as controlling inputs.
- Load only the approved consolidated FEATURE-0013 architecture and `ADH-2026-017` when validating the downstream boundary; do not load predecessor handoffs as active instructions.
- Generate `requirements.md` only after this handoff is approved.
- Do not introduce resource kinds beyond the locked five.
- Treat `InfrastructureStack` as a deliberate terminology replacement, not a sixth kind or compatibility alias; do not retain both names in active contracts.
- Reuse the existing `Organization` resource for the supply owner; do not introduce `Owner`, `ProviderOwner`, `SupplyOwner`, or another scope kind.
- Use Provider `scopeRef`, not `ownerRef`, for the owner/governance relationship.
- Use only `ProviderLocation` for the location level; introduce no alias kind, field, or term.
- Do not infer network connectivity or isolation from any topology relationship.
- Do not introduce network-connectivity fields or pre-design FEATURE-0053.
- Do not define ResourcePool, ProviderCapability, adapter interfaces, provider discovery, provider calls, or execution behavior.
- Do not modify FEATURE-0012 or FEATURE-0013 semantics.
- Treat `F14-AD-001` through `F14-AD-021` as a closed decision register. Translate them into testable requirements; do not reconsider them.
- Follow the stage-ownership matrix. Requirements specify testable outcomes, design selects only explicitly delegated representations, and tasks implement only approved design.
- Enforce the compatibility matrix, context allowlist, source-precedence rule, and single-owner overlap check before writing the first normative requirement.
- Enforce architecture section 16 as the closed Kiro generation contract for requirements, design, and tasks.
- Generate exactly one authorized stage and one allowlisted stage file at a time; never generate a later stage speculatively.
- Do not modify an approved earlier-stage file during later-stage generation.
- Run and report the deterministic completeness checks before presenting a stage for human review.
- Run the blocking repository-readiness gate before generating requirements. Report `REPOSITORY_CONTEXT_NOT_READY` for repository, prompt, reviewer, automation-state, or manifest failures; do not treat them as requirements revisions.
- Generate and validate the canonical requirement-normalization ledger; reference shared obligations instead of copying them.
- Do not satisfy generic prompt/reviewer checks by inventing decision, audit, operation, connectivity, capability, adapter, or runtime semantics.
- Produce the requirements overlap ledger as traceability within `requirements.md`; do not create a separate adoption or ownership-management system.
- Do not use superseded handoffs, generated artifacts, prior specs, automation logs, roadmap placeholders, examples, or sketches as independent semantic authority.
- Every normative requirement must cite an `F14-AD-*` decision or an explicitly inherited FEATURE-0012 contract.
- If a semantic choice is unresolved or outside the delegation matrix, emit `ARCHITECTURE_DECISION_REQUIRED` and stop the entire current generation stage without selecting a default.
- Preserve the complete `F14-R01`–`F14-R30` risk traceability. Do not remove, merge, downgrade, accept, or renumber risks in downstream stages.
- Treat target residual ratings as goals pending evidence, not as architecture-owner acceptance.
- Do not generate design, tasks, schemas, or Go code until separately authorized.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-07-29
- Notes:
  - Approval authorizes Kiro requirements generation only; normal later stage gates remain mandatory.

This Architecture Decision Handoff is not an approval by itself.
