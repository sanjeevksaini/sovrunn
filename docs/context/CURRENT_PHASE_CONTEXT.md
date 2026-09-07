# Current Phase Context

Current active phase: Phase 2R.

Architecture baseline: `ARCH-2026.08-PHASE2R-CANONICAL`.

Controlling adoption: `ADH-2026-042`, `ACR-2026-001`, `ARCH-APPROVAL-2026-004`.

FEATURE-0018 status: implemented and merged; PR #20 merged 2026-09-07 as commit `2295e9e` into `phase2-reuse-first-paas-fabric-foundation`.

## Phase 2R Goal

Establish the canonical, implementation-neutral PaaS fabric foundation with seven-scope governance, qualified ExecutionTarget placement, sovereign evidence-backed decisions, and VS-000 cross-feature conformance — without real provider execution.

## Phase 2R Build Scope

- Canonical cloud model foundation via direct bootstrap (seven scopes, CloudPlatform, CloudProvider; ExecutionTarget owned by FEATURE-0016)
- Adapter boundary and ExecutionTarget qualification
- Policy evaluation abstraction with DecisionRecord linkage
- Governance, IAM, approval, and exception foundation
- Sovereignty facts, evidence, and policy foundation
- Assignment and effective governance resolution (EffectiveGovernanceContext)
- CloudEnrollment, personal onboarding, entitlement, and quota
- Service product, runtime, and requirement foundation (ServiceTypeDefinition, ServiceOffering, ServicePlan, ServiceRequirementSet)
- Sovereignty and placement decision v0 (DecisionRecord profiles)
- Plugin taxonomy and synthetic execution boundary
- AI-readable decision and operation context
- Slice 0 integration and conformance demo (VS-000)

## Phase 2R Completed Features

| Feature | Status | Approval |
|---|---|---|
| FEATURE-0011 Reuse Assessment Standard | Merged | Complete |
| FEATURE-0012 API, Resource Naming, Status, and Validation Standard | Merged through PR #14 | Final human approval 2026-07-24 |
| FEATURE-0013 Decision Record and AuditEvent Standard | Merged through PR #15 | Final review 2026-07-29 |
| FEATURE-0014 Provider-Neutral Resource Model | Merged through PR #16 | Final feature gate 2026-07-30 |
| FEATURE-0018 Governance, IAM, Approval and Exception Foundation | Implemented and merged; PR #20 merged 2026-09-07 as commit `2295e9e` into `phase2-reuse-first-paas-fabric-foundation` | Final feature gate passed 2026-09-07 |

## Phase 2R Next Planned Feature

| Feature | Status | Controlling decisions |
|---|---|---|
| FEATURE-0019 Sovereignty Facts, Evidence and Policy Foundation | Architecture not started | Pending architecture decision handoff |

## Approved FEATURE-0018 architecture application

ACR-2026-002 and DEC-0060 apply the approved ADH-2026-070 FEATURE-0018 package.
Final human and reuse approval are recorded, and independent security-review
renewal-05 passed with no blocking findings. The approved package limits approver,
privileged-requester and reviewer eligibility to exact Human PrincipalRefs;
roles, assignments, groups and Membership cannot manufacture that eligibility.
This record does not change the current feature execution position or
independently authorize requirements, design, tasks or implementation.

## Phase 2R Feature Sequence

See `docs/phase2/PHASE2_FEATURE_SEQUENCE.md` for the canonical order.

## FEATURE-0014 Completed Architecture Boundary

FEATURE-0014 remains completed history. Its alpha model resources (Provider, ProviderLocation, ProviderDatacenter, DatacenterFailureDomain, InfrastructureStack) are retained repository assets and reuse input under FEATURE-0011 (DEC-0059); they are not live control-plane data requiring runtime conversion. The active canonical model uses CloudPlatform, CloudProvider, HostingLocation, Datacenter, FaultDomain, InfrastructureStack, and ExecutionTarget, created directly by FEATURE-0015/0016 rather than migrated.

Later features consume FEATURE-0014 through its completed contracts:

- FEATURE-0011 controls reuse-assessment format and reuse-before-build gate.
- FEATURE-0012 controls API/resource grammar, seven-scope ScopeKind vocabulary (DEC-0037), typed references, status/validation shape, and Problem Details envelope.
- FEATURE-0013 controls DecisionRecord, DecisionProfile, EvaluationResult, AuditEvent, metadata.scopeRef as sole scope authority, and ServiceInstance as typed subject.
- FEATURE-0014 implementation history owns the alpha Provider/ProviderLocation/ProviderDatacenter/DatacenterFailureDomain/InfrastructureStack as retained repository assets assessed for reuse under FEATURE-0011; per DEC-0059 they are not runtime-migrated, since FEATURE-0015 creates canonical resources directly.

## Phase 2R Exit Criteria

- All Phase 2R features pass feature gate.
- VS-000 definition of done and conformance matrix pass.
- Sovereignty and placement simulation demonstrates selected, denied, requires-approval, and indeterminate outcomes.
- DecisionRecord and AuditEvent semantics remain owned by FEATURE-0013.
- Customer contracts contain no provider-native objects, raw secrets, or protected handles.
- No active desired-state API supports old and canonical kinds/scopes simultaneously.
- Phase 3 readiness review approves replacement of fake execution with one real PostgreSQL path.
