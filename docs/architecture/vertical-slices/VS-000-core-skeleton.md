# VS-000: Executable Canonical Core Skeleton

**Status:** Approved target charter; repository merge pending
**Phase:** Phase 2R
**Contributing features:** FEATURE-0015 through FEATURE-0026
**Controlling architecture:** ADH-2026-020–041 through consolidated ADH-2026-042
**Provider fixture:** NIC CloudPlatform with Yotta CloudProvider participation
**Customer fixture:** Department of Posts
**Service fixture:** Synthetic Echo Service
**Backend fixture:** Deterministic in-memory fake target, plugin and adapter
**External side effects:** Prohibited

**Exact Slice 0 integration profile:**
`docs/architecture/vertical-slices/VS-000-contract-specification.md` with
machine-readable authority in `VS-000-contract-registry.yaml`. This charter
owns cross-feature intent; the detailed profile owns Slice 0 schema, writer,
state, error and conformance precision without replacing feature requirements.

## 1. Purpose

Prove that Sovrunn's canonical service-first contracts compose end to end before integrating PostgreSQL, OpenShift or another real backend.

The slice validates the framework, ownership boundaries, writer rules, authorization, decisions, audit, idempotency and information-disclosure controls. It does not validate infrastructure or database technology.

## 2. Provider and customer outcomes

### Cloud owner outcome

NIC can publish one simple synthetic service through NIC Cloud, enroll the Department of Posts, and allow either explicit Yotta selection or governed provider assignment without exposing Yotta credentials or backend objects.

### CloudProvider outcome

Yotta can register one fake execution target, expose normalized facts and operate through an isolated participation without becoming the owner of NIC Cloud or the customer Organization.

### Customer outcome

An authorized Department of Posts user requests one service using a plan, service region and optional provider preference. The user receives explainable status, safe placement and a revocable binding without learning internal target or credential details.

### Platform outcome

Sovrunn traces the request through one authorization path, one governance context, FEATURE-0013 DecisionRecords, one immutable plan, one Operation, one fake PluginExecution and correlated AuditEvents.

## 3. Goals

1. Exercise every canonical boundary required by the first real PostgreSQL slice.
2. Prove the CloudPlatform/CloudProvider/customer separation and one-provider-per-installation invariant.
3. Prove implementation-neutral product, requirement, decision and execution contracts.
4. Prove that governance and sovereignty inputs are versioned and fail closed when missing or stale.
5. Prove Operation-first, idempotent, auditable simulated execution.
6. Prove safe customer projection and SecretRef-only binding.
7. Prove that AI-readable context is a projection with no decision or execution authority.
8. Produce conformance evidence that supports Alpha promotion only after a later real slice.

## 4. Non-goals

- Real infrastructure or workload provisioning.
- Real Kubernetes/OpenShift, cloud SDK, operator, Helm or GitOps calls.
- PostgreSQL-specific fields or lifecycle behavior.
- Production authentication, policy, secret, workflow, evidence or observability integrations.
- Multi-target HA, failover, migration, capacity reservation or data movement.
- Native network, subnet, route, security-group or NetworkPolicy modeling.
- Platform installation, upgrade or restore execution.
- Customer billing, marketplace or UI.
- Autonomous AI action.
- Stable API commitment.

## 5. Fixture graph

```text
Organization/NIC
  └── CloudPlatform/NIC-Cloud
        ├── ServiceRegion/India-North
        ├── ServicePortfolio/Government-Development
        ├── ServiceOffering/Synthetic-Echo
        │     └── ServicePlan/Echo-Small
        ├── CloudProviderParticipation/NIC-Yotta-Development
        │     └── CloudProvider/Yotta
        └── CloudEnrollment/Posts-NIC-Cloud
              └── Organization/Department-of-Posts
                    └── Tenant/Development
                          └── Project/Slice-0
                                ├── ServiceInstance/Postal-Echo
                                └── ServiceBinding/Postal-Echo-Client

CloudProvider/Yotta
  └── HostingLocation/India-Delhi-NCR
        └── Datacenter/Yotta-Fixture-DC-A
              └── FaultDomain/FD-A
                    └── InfrastructureStack/Fake-Stack-A
                          └── ExecutionTarget/Fake-Yotta-Target-A
```

Fixture names are synthetic. They do not assert a real NIC/Yotta contract, facility, certification or customer deployment.

## 6. Participating canonical contracts

| Domain | Contracts used in Slice 0 | Owning feature |
|---|---|---|
| Common grammar | metadata, scopeRef, ownerRef, typed references, generation/resourceVersion, status/conditions, stable Problem Details | FEATURE-0012 retained |
| Decision/audit | DecisionProfile, DecisionRecord, AuditEvent | FEATURE-0013 retained |
| Cloud model | CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, InfrastructureStack | FEATURE-0015 |
| Integration | ExecutionTarget, normalized target facts, qualification, and synthetic observer boundary; consumes CloudProviderParticipation and InfrastructureStack read-only | FEATURE-0016 |
| Policy | PolicyEvaluationRequest, PolicyEvaluationResult, fake PolicyEngineAdapter | FEATURE-0017 |
| IAM/governance | PrincipalRef, AccessGroup, Membership, RoleDefinition, RoleAssignment, Human PrincipalRef-only EligibilityRef, PrivilegedAccessRequest, AccessReview, ApprovalPolicy, ApprovalRequest, ExceptionGrant and FEATURE-0018-limited GovernanceProfile | FEATURE-0018 |
| Sovereignty inputs | SovereigntyProfile, RegulatoryPolicyBundle, SovereigntyFactSet, EvidenceRecord | FEATURE-0019 |
| Resolution | ProfileAssignment, EffectiveGovernanceContext | FEATURE-0020 |
| Customer access to catalog | CloudEnrollment, EntitlementPackage, ServiceEntitlement, QuotaPolicy/reservation | FEATURE-0021 |
| Service product | ServicePortfolio, ServiceTypeDefinition, ServiceOffering, ServicePlan, ServicePlacementProfile, ServiceRuntimeProfile, ServiceRequirementSet, ServiceRegion | FEATURE-0022 |
| Decisions/projection | sovereignty and placement DecisionRecord profiles, ServicePlacement | FEATURE-0023 |
| Simulated execution | plugin/adapter manifests, ServiceDeploymentPlan, Operation reuse, minimal PluginExecution and fake harness | FEATURE-0024 |
| Explanation | bounded AI-readable decision/operation context | FEATURE-0025 |
| Integration proof | fixtures, orchestrated simulation, conformance report and demo | FEATURE-0026 |
| Binding | existing ServiceBinding strengthened by ADH-2026-034; fake SecretRef only | FEATURE-0008 retained plus FEATURE-0021/0024 adoption |

No slice feature may redefine a contract owned by another feature.

## 7. Actors and authorization boundaries

| Actor | Allowed actions | Explicitly prohibited |
|---|---|---|
| NIC CloudPlatform administrator | Publish fixture catalog/profiles, create participation and enrollment within delegated scope | Provider credentials, customer Project content, authoritative decisions |
| Yotta infrastructure operator | Register fixture topology/target and normalized provider facts within participation | NIC catalog grants, customer intent, target status owned by controllers |
| Department organization administrator | Manage customer hierarchy, membership and narrower policy | CloudPlatform entitlement grants, provider topology or evidence |
| Project service consumer | Create/get/delete synthetic ServiceInstance and request/revoke own binding | Status, decisions, plans, target IDs, raw secret values |
| Policy adapter | Return deterministic evaluation result | Mutate request, profiles, decisions or execution |
| Evidence collector | Publish fixture facts/evidence for authorized subject | Desired state or assessment outcome |
| Decision service | Create immutable authoritative DecisionRecord | Alter evidence, policy or customer intent |
| Deployment planner | Create immutable ServiceDeploymentPlan | Select targets or execute steps |
| Operation controller | Accept idempotent Operation and own progress | Mutate accepted request or bypass final decisions |
| Fake plugin/adapter executor | Execute only approved fake steps and publish protected result | External API calls, policy bypass, decision/plan mutation |
| Projection controller | Publish safe ServicePlacement and explanation | Canonical decision mutation or protected topology disclosure |
| Binding controller | Publish safe endpoint metadata and fake SecretRef | Embed or log secret value |

Authorization uses authenticated PrincipalRef, an exact registered action/target,
applicable scoped RoleAssignments and explicit dependencies/guardrails. Grants
combine by union and guardrails by intersection; ambiguity fails closed.
Persona labels, job titles, email, provider-native roles and unprovisioned
external group claims grant no authority. DEC-0060 remains Proposed, so this is
an architecture reconciliation record rather than executable Slice 0 authority.

## 8. Happy path

1. Authenticate the fixture project user and resolve stable PrincipalRef.
2. Authorize `services.instances.create` at Project/Slice-0.
3. Verify active CloudEnrollment with NIC Cloud.
4. Verify entitlement to Synthetic Echo, Echo Small and India North.
5. Reserve one instance against quota using request idempotency key.
6. Validate the exact published ServicePlan and bounded parameters.
7. Resolve exact GovernanceProfile/ProfileAssignment versions into immutable EffectiveGovernanceContext.
8. Resolve ServiceRuntimeProfile and request into immutable ServiceRequirementSet.
9. Discover active CloudProviderParticipation candidates allowed by provider-selection intent.
10. Load only qualified/available fake ExecutionTargets with fresh normalized facts and evidence.
11. Create authoritative sovereignty DecisionRecord using the registered profile.
12. Create authoritative placement DecisionRecord selecting NIC-Yotta-Development/Fake-Yotta-Target-A.
13. Create immutable ServiceDeploymentPlan containing abstract `Create`, `Verify` and `PublishBindingEligibility` steps; no native payload.
14. Accept one idempotent Operation linked to the ServiceInstance and exact plan.
15. Execute steps through the fake plugin/adapter, producing protected PluginExecution records with no external effects.
16. Verify deterministic synthetic readiness from the fake result contract.
17. Publish customer-safe ServicePlacement showing NIC Cloud, India North, selected provider only if disclosure policy permits, and summarized sovereignty/decision outcome.
18. Authorize and issue ServiceBinding for one fake workload consumer, publishing safe endpoint metadata and a SecretRef identifier only.
19. Publish correlated AuditEvents and bounded AI-readable explanation.
20. On delete, revoke binding, execute fake deletion idempotently, release quota and finalize the ServiceInstance without erasing retained decisions/audit.

## 9. Cross-feature orchestration milestones

These are slice milestones, not new resource status enums. Each feature design owns its exact status vocabulary.

```text
Request accepted
  → Governance resolved
  → Requirements resolved
  → Sovereignty decided
  → Placement decided
  → Plan created
  → Operation accepted
  → Fake execution completed
  → Readiness verified
  → Safe placement published
  → Binding issued
  → Deletion and access revocation verified
```

No milestone may be inferred from logs. Authoritative records and controller-owned status establish progression.

## 10. Critical failure paths

| ID | Failure | Required result and stable semantic category |
|---|---|---|
| VS0-F01 | Unauthenticated request | Deny before scope or resource disclosure; no quota or operation side effect |
| VS0-F02 | Unauthorized Project action | Safe denial with no-existence disclosure behavior |
| VS0-F03 | Inactive CloudEnrollment | Deny consumption; customer hierarchy remains intact |
| VS0-F04 | Missing entitlement | Deny even when quota exists |
| VS0-F05 | Quota exhausted | Deny even when entitlement exists; no leaked capacity details |
| VS0-F06 | Invalid/unpublished plan version | Reject deterministically before decision or operation |
| VS0-F07 | Conflicting governance profiles | Produce conflict/approval-required outcome; never choose silently |
| VS0-F08 | Missing/stale sovereignty evidence | Indeterminate/deny according to mandatory profile; never infer success from country |
| VS0-F09 | Requested provider participation unavailable | Respect required/preferred semantics; deny or choose allowed alternative through PlacementDecision |
| VS0-F10 | Target unqualified, draining or maintenance epoch changed | Exclude candidate or fail stale work closed |
| VS0-F11 | No compatible candidate | Immutable denied PlacementDecision with safe reasons and alternatives |
| VS0-F12 | Repeated create with same idempotency key and same request | Return same accepted resource/Operation identity; do not duplicate quota or execution |
| VS0-F13 | Reused idempotency key with different request | Conflict; no mutation of original request |
| VS0-F14 | Plan/decision generation changed before operation acceptance | Reject stale generation and require new plan |
| VS0-F15 | Fake plugin failure | Operation fails with normalized code; compensation/release behavior follows plan; no readiness/binding |
| VS0-F16 | Fake verification failure after create | Service not Ready; no binding; explicit recovery/delete path |
| VS0-F17 | Binding requested before readiness or by wrong consumer | Deny; no SecretRef issued |
| VS0-F18 | Secret value or protected handle appears in response/log/audit/explanation | Test failure and release blocker |
| VS0-F19 | AI explanation attempts decision/approval/execution | Deny by contract and authorization |
| VS0-F20 | Delete with active binding | Revoke access first; finalization waits for fake cleanup or records explicit retained failure |

Exact public error codes and HTTP bindings must be owned by feature requirements/design and inherited FEATURE-0012 Problem Details rules. This charter does not invent parallel error envelopes.

## 11. Idempotency, concurrency and compensation

- Every create/delete/binding request carries a stable idempotency key.
- Mutation validates resourceVersion/generation and applicable decision, evidence and target epochs.
- Accepted Operation request and ServiceDeploymentPlan are immutable.
- Duplicate execution attempts use stable operation/plan/step identity and cannot repeat externally consequential work; the fake harness must prove this behavior.
- Failed pre-execution evaluation releases provisional quota without creating an Operation.
- Failure after Operation acceptance follows explicit plan compensation; no silent desired-state rollback occurs.
- Deletion retains DecisionRecord, Operation and AuditEvent according to policy.

## 12. Security and information-disclosure tests

Required tests verify:

- cross-Organization, cross-Project and cross-provider references deny by default;
- a CloudPlatform administrator cannot read provider credentials;
- a CloudProvider administrator cannot read customer Project content;
- the customer cannot enumerate ExecutionTargets or protected evidence;
- customer status, errors and explanation contain no target-native ID, credential, token, SecretRef value or protected topology;
- fake credentials are stored only in the approved fake secret provider and represented by typed SecretRef;
- structured logs, traces, metrics labels, AuditEvents and test snapshots contain no secret value;
- privileged roles are time-bound where used;
- decision and audit projections are audience-specific;
- all externally supplied strings and payloads are bounded and strictly validated.

## 13. Observability and audit

Every request carries one request ID and trace context. Every accepted lifecycle action has one Operation ID. Every fake plan step has one PluginExecution ID. Logs, metrics and traces correlate these identifiers but remain operational telemetry.

Minimum metrics:

- request count/latency/error by canonical action and safe result category;
- decision count/latency/outcome by registered profile;
- operation and fake execution duration/result;
- stale evidence/target/epoch denials;
- quota reservation/release outcomes;
- projection and binding issuance failures.

Minimum AuditEvents:

- enrollment/entitlement/quota decision where policy requires;
- material authorization outcome;
- sovereignty and placement decision linkage;
- Operation acceptance and terminal outcome;
- binding issue/revoke;
- privileged access or exception use;
- migration/conformance fixture activation where applicable.

AuditEvent is not a log line. Log absence never removes required audit evidence.

## 14. Demo script

The FEATURE-0026 demo must show:

1. fixture bootstrap;
2. successful synthetic service request;
3. trace from request through exact decisions, plan, Operation and fake executions;
4. safe placement and binding output;
5. repeat request proving idempotency;
6. one entitlement denial;
7. one stale-evidence sovereignty denial;
8. one target-draining/stale-epoch denial;
9. one fake execution failure with no binding;
10. deletion, binding revocation and quota release;
11. audit correlation and safe AI-readable explanation;
12. proof that the fake adapter made zero external calls.

## 15. Definition of done

Slice 0 is complete only when:

- all contributing features have approved requirements, design and tasks and pass their individual gates;
- canonical migration fixtures pass with no dual authority;
- the full happy path passes deterministically from a clean state;
- VS0-F01 through VS0-F20 have explicit automated coverage or approved non-applicability rationale;
- every authoritative field has one documented writer;
- every reference has scope/kind/identity validation;
- every immutable record pins exact input versions;
- no raw secret, protected handle or provider-native object crosses the customer/core boundary;
- fake execution is idempotent and proves zero external side effects;
- logs, metrics, traces and AuditEvents are correlated and separated correctly;
- the demo and conformance report are repeatable;
- `make fmt`, `make test`, `make vet`, applicable race/security checks and `mkdocs build --strict` pass;
- the feature and Phase 2R human acceptance gates approve the evidence;
- Phase 3 is not started until the Slice 0 exit review is accepted.

## 16. API maturity evidence

Slice 0 proves only Experimental contract integration. It does not promote contracts to Alpha by itself. Alpha requires the real PostgreSQL Slice 1 to pass positive, negative, lifecycle and security conformance. Beta and Stable require later production and portability evidence defined by the canonical delivery plan.

## 17. Instructions to Kiro

1. Treat this charter as cross-feature acceptance context, not a substitute for feature requirements.
2. Generate one feature and one stage at a time in the approved Phase 2R order.
3. Do not resolve any semantic gap in requirements, design or tasks; stop with `ARCHITECTURE_DECISION_REQUIRED`.
4. Preserve FEATURE-0012 and FEATURE-0013 ownership; do not redefine their envelopes.
5. Implement only deterministic fake execution in Phase 2R and explicitly prove zero external effects.
6. Keep PostgreSQL, OpenShift, native network and provider SDK details out of Slice 0.
7. Map every requirement to ADH/DEC, owning feature, slice gate and conformance ID.
8. Do not generate design without `APPROVED_FOR_DESIGN` or tasks without `APPROVED_FOR_TASKS`.
