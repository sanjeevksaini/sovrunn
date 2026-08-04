# Current Phase Context

Current active phase: Phase 2R.

Architecture baseline: `ARCH-2026.08-PHASE2R-CANONICAL`.

Controlling adoption: `ADH-2026-042`, `ACR-2026-001`, `ARCH-APPROVAL-2026-004`.

## Phase 2R Goal

Establish the canonical, implementation-neutral PaaS fabric foundation with seven-scope governance, qualified ExecutionTarget placement, sovereign evidence-backed decisions, and VS-000 cross-feature conformance — without real provider execution.

## Phase 2R Build Scope

- Canonical cloud model and alpha migration foundation (seven scopes, CloudPlatform, CloudProvider, ExecutionTarget)
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

## Phase 2R Next Planned Feature

| Feature | Status | Controlling decisions |
|---|---|---|
| FEATURE-0015 Canonical Cloud Model and Alpha Migration Foundation | Architecture ready; requirements not yet generated | DEC-0037, DEC-0038, DEC-0041, DEC-0042, DEC-0054, DEC-0058; ADH-2026-020/024/025/037/041 |

## Phase 2R Feature Sequence

See `docs/phase2/PHASE2_FEATURE_SEQUENCE.md` for the canonical order.

## FEATURE-0014 Completed Architecture Boundary

FEATURE-0014 remains completed history. Its alpha model resources (Provider, ProviderLocation, ProviderDatacenter, DatacenterFailureDomain, InfrastructureStack) are implementation artifacts that migrate through the canonical cutover (DEC-0058). The active canonical model uses CloudPlatform, CloudProvider, HostingLocation, Datacenter, FaultDomain, InfrastructureStack, and ExecutionTarget.

Later features consume FEATURE-0014 through its completed contracts:

- FEATURE-0011 controls reuse-assessment format and reuse-before-build gate.
- FEATURE-0012 controls API/resource grammar, seven-scope ScopeKind vocabulary (DEC-0037), typed references, status/validation shape, and Problem Details envelope.
- FEATURE-0013 controls DecisionRecord, DecisionProfile, EvaluationResult, AuditEvent, metadata.scopeRef as sole scope authority, and ServiceInstance as typed subject.
- FEATURE-0014 implementation history owns the alpha Provider/ProviderLocation/ProviderDatacenter/DatacenterFailureDomain/InfrastructureStack. These migrate per DEC-0037/DEC-0041/DEC-0058.

## Phase 2R Exit Criteria

- All Phase 2R features pass feature gate.
- VS-000 definition of done and conformance matrix pass.
- Sovereignty and placement simulation demonstrates selected, denied, requires-approval, and indeterminate outcomes.
- DecisionRecord and AuditEvent semantics remain owned by FEATURE-0013.
- Customer contracts contain no provider-native objects, raw secrets, or protected handles.
- No active desired-state API supports old and canonical kinds/scopes simultaneously.
- Phase 3 readiness review approves replacement of fake execution with one real PostgreSQL path.
