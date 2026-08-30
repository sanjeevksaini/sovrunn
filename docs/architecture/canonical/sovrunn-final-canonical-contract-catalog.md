# Sovrunn Final Canonical Contract Catalog

**Status:** Architecture-owner approved target architecture; FEATURE-0018 repository amendments complete; feature-level schemas pending
**Version:** 1.8
**Date:** 30 August 2026
**Companion to:** `sovrunn-finalized-data-model.md`

**Formal decisions:** [Final canonical ADR package](sovrunn-final-canonical-adrs/README.md)
**Reference conformance scenario:** [Yotta–Department of Posts PostgreSQL flow](sovrunn-postgresql-reference-flow.md)

## 1. Purpose

This catalog defines the architecture contract for every concept in the Sovrunn final canonical model. It fixes:

- canonical and API names;
- purpose, ownership and scope;
- resource profile and API boundary;
- identity and reference rules;
- authoritative writers of desired and observed state;
- mutability, lifecycle and deletion behavior;
- security and visibility boundaries;
- relationships, cardinalities and invariants;
- skeletal YAML proving that the contracts compose.

It does not replace feature requirements, OpenAPI schemas, controller design or adapter-specific payloads.

## 2. Required compatibility amendments

The final target model intentionally corrects parts of the current Phase 2 repository baseline. These changes must be approved and migrated together; later feature authors must not reconcile them independently.

| Amendment | Current repository baseline | Final canonical contract | Required action |
|---|---|---|---|
| Cloud ownership and provider identity | `Provider` conflates customer cloud ownership, product publication and infrastructure operation | Owner `Organization` + `CloudPlatform` for the customer-facing cloud; independent `CloudProvider` joined by `CloudProviderParticipation` | Amend FEATURE-0012 scope vocabulary and dependent schemas; classify rather than blindly rename alpha records; bind every installation to one participation |
| Physical topology | `ProviderLocation → ProviderDatacenter → DatacenterFailureDomain → InfrastructureStack` strict containment | `HostingLocation`, `Datacenter`, `FaultDomain`, `InfrastructureStack` (FEATURE-0015) and `ExecutionTarget` (FEATURE-0016) as provider-registered facts and relationships | Supersede FEATURE-0014 architecture through an approved amendment; FEATURE-0015 creates the canonical topology resources directly (DEC-0059; no migration mapping) |
| Placement and sovereignty conclusions | FEATURE-0013 owns the only common decision envelope | `PlacementDecision` and `SovereigntyAssessment` remain domain concepts | Implement each as a registered `DecisionProfile` and typed `DecisionRecord`; do not create competing decision envelopes |
| Plural conceptual names | `SovereigntyFacts` and `ServiceRequirements` are conceptual labels | API kinds must be singular PascalCase | Use `SovereigntyFactSet` and `ServiceRequirementSet` as API kinds while retaining the conceptual labels in explanatory text |
| Infrastructure allocation | FEATURE-0015 planned `ResourcePool` and provider capability resources | No mandatory resource pool or provider-wide capability in the core | Reframe post-0014 roadmap around execution targets, normalized facts, requirements and adapters |
| CloudPlatform product catalog | FEATURE-0006 implements a global mutable `ServiceClass` and `ServicePlan` | Reusable versioned `ServiceTypeDefinition` plus CloudPlatform-scoped, versioned `ServiceOffering` and `ServicePlan` definitions | Split reusable technical semantics from cloud-product identity; FEATURE-0022 creates the canonical catalog resources directly; FEATURE-0006 remains retained repository history assessed for reuse (DEC-0059) |
| Customer enrollment | Customer Organization enrolls directly with the old Provider/CloudProvider product boundary | `CloudEnrollment` joins customer Organization to CloudPlatform; provider realization eligibility derives from participation and placement | FEATURE-0021 creates `CloudEnrollment` directly; there is no old/new dual authority to prohibit because there is no runtime migration (DEC-0059) |
| Governance resolution | Planned FEATURE-0018–0021 split overlapping policy/profile kinds and use `EffectivePolicyContext` | Compositional `GovernanceProfile`, separate `SovereigntyProfile` and immutable `EffectiveGovernanceContext` | Rebaseline the planned features and amend FEATURE-0013 terminology without weakening provenance or typed-reference rules |
| Governance, IAM, approval and exception | Preliminary identity inventory binds assignments directly to principals and does not close group, validity, workflow, writer, evidence or native-IAM boundaries | FEATURE-0018 closed contract inventory, direct-member `AccessGroup`, `RoleHolderRef`, exact scopes/targets, Standing/TimeBound assignments, retained approvals/reviews/exceptions, audit-before-publication and two-plane IAM intersection | Apply ACR-2026-002, DEC-0060 and ADH-2026-070 without adding provider-native IAM or future-feature resolution |
| Managed-service access | FEATURE-0008 implements `ServiceBinding` and planned FEATURE-0033 extends it | `ServiceBinding` remains the per-consumer, Project-scoped access boundary | Retain the resource; require SecretRef-only delivery, independent rotation/revocation and lifecycle protection |
| Sovrunn platform lifecycle | No complete canonical model for installing, upgrading, backing up, restoring or recovering the Sovrunn control plane | `SovrunnInstallation`, immutable `SovrunnRelease`, `PlatformLifecyclePolicy`, immutable `PlatformLifecyclePlan`, canonical `Operation`, external `PlatformLifecycleAgent` and safe `PlatformHealth` projection | Add an operator-facing platform lifecycle module before production MVP; keep recovery execution independent from the affected installation and do not absorb underlying infrastructure lifecycle |
| Release and recovery compatibility | Version strings and skeletal compatibility fields do not define all executable transitions | Explicit directed `ReleaseCompatibilityContract`, typed lifecycle references, three recovery actions and an irreversible-boundary constraint | Deny any transition lacking a complete published edge and verified recovery path |
| Infrastructure maintenance | Normalized signal and Draining status lack exact authority and concurrency semantics | `InfrastructureMaintenanceNotice`, lifecycle-controller-owned availability, maintenance epoch/fence and mandatory requalification | Separate provider notice authority from Sovrunn target-state authority; reject stale placement/execution |
| Canonical bootstrap (no alpha runtime migration) | FEATURE-0001–0014 are retained repository assets, not live control-plane state | No `CanonicalMigrationPlan`/`CanonicalMigrationRecord`, migration controller, or cutover state machine; canonical resources are created directly | Retire the migration model per DEC-0059 (supersedes DEC-0058); reject reintroduction of migration-plan/record contracts as active behavior |

Until those amendments are approved in the repository, this document describes the target architecture rather than claiming that the current schemas already implement it.

## 3. Common API contract

### 3.1 Type and metadata

Externally exchanged objects use:

```yaml
apiVersion: <domain>.sovrunn.io/v1alpha1
kind: <SingularPascalCaseKind>
metadata:
  name: <immutable-kebab-case-name>
  uid: <server-generated-opaque-uid>
  displayName: <mutable-human-name>
  scopeRef:
    apiVersion: core.sovrunn.io/v1alpha1
    kind: <ScopeKind>
    name: <scope-name>
    uid: <resolved-scope-uid>
  generation: 1
  resourceVersion: <opaque-concurrency-token>
  createdAt: <rfc3339>
  updatedAt: <rfc3339>
```

Universal rules:

1. Identity uniqueness is API group + kind + scope UID + name.
2. `name`, `uid` and `scopeRef` are immutable.
3. Authorization resolves scope by UID, never by display name.
4. `ownerRef` expresses lifecycle containment only and never replaces `scopeRef`.
5. References use `apiVersion`, `kind`, `name` and optional resolved `uid`.
6. Cross-organization and cross-provider references deny by default and do not disclose inaccessible targets.
7. Provider-native identifiers and secrets never become Sovrunn identities.
8. Secret material is represented only by typed secret references.

For readability, scalar values shown in skeletal fields ending in `Ref` or `Refs` abbreviate the complete typed-reference object. Normative schemas must use constrained typed references; the shorthand is not an alternative API representation.

### 3.2 Final scope vocabulary

```text
Platform
Organization
OrganizationUnit
Tenant
Project
CloudPlatform
CloudProvider
```

`CloudPlatform` is the customer product/governance scope and `CloudProvider` is the operational supply scope. Together they replace the ambiguous current alpha `Provider` scope after classification. Resource-local containment is not an eighth security scope.

### 3.3 Resource profiles

| Code | Profile | Contract |
|---|---|---|
| MR | `ManagedResource` | Client/operator owns mutable desired `spec`; controller owns `status` |
| VD | `VersionedDefinition` | Draft is mutable; a published version is immutable and later retired or superseded |
| OER | `ObservedExternalResource` | Authorized collector owns normalized observations, provenance and freshness |
| IR | `ImmutableRecord` | System-produced append-only record; corrections create linked records |
| LRO | `LongRunningOperation` | Request becomes immutable after acceptance; executor owns progress and terminal status |
| EV | `EmbeddedValue` | Parent-owned value with no independent identity or lifecycle |

### 3.4 API boundaries

| Code | Boundary | Rule |
|---|---|---|
| CF | customer-facing | Product intent and safe status only; no provider internals or secrets |
| OF | operator-facing | Provider administration and normalized infrastructure; no raw secrets |
| POF | platform-operator-facing | Sovrunn installation, release, lifecycle, recovery and health operations; never exposed to managed-service customers |
| IE | internal-engine-facing | Normalized policy, placement and orchestration contracts |
| AF | adapter-facing | Backend translation, handles and provenance; never leaked into core/customer contracts |
| PF | plugin-facing | Validated service-lifecycle requests and results; no policy bypass |
| GO | governance-only | Canonical decisions, evidence, approvals and audit; audience-specific projections required |

Where multiple audiences need different data, Sovrunn publishes distinct projections rather than one schema with conditionally hidden fields.

### 3.5 Status and conditions

MR, VD, OER and LRO resources may expose system-owned status:

```yaml
status:
  observedGeneration: 1
  phase: Ready
  conditions:
    - type: Valid
      status: "True"
      reason: ValidationSucceeded
      observedGeneration: 1
      lastTransitionTime: <rfc3339>
```

Conditions represent current facts, not history. Decisions, evidence and activity history belong in immutable records.

### 3.6 Lifecycle and deletion codes

| Code | Behavior |
|---|---|
| RESTRICT | Reject deletion while dependants exist |
| RETIRE | Stop new use but retain the definition for existing references |
| DECOMMISSION | Drain or terminate use through an explicit lifecycle before deletion |
| CASCADE-EXPLICIT | Cascade only when explicitly requested, authorized and impact-previewed |
| RETAIN | Preserve according to audit/evidence retention and legal-hold policy |

## 4. Architecture decisions

| ID | Final decision |
|---|---|
| `ADH-2026-020` | Generic `Provider` is retired; its records are classified into CloudPlatform ownership and/or CloudProvider operation under ADH-2026-037/041 |
| `ADH-2026-021` | Customer governance is independent and consumption uses `CloudEnrollment` with CloudPlatform; ADH-2026-037 supersedes direct provider enrollment |
| `ADH-2026-022` | Entitlement answers what may be consumed; quota independently answers how much |
| `ADH-2026-023` | Industry/domain clouds are CloudPlatform-scoped `ServicePortfolio` configurations, not new core provider or resource kinds |
| `ADH-2026-024` | Physical geography is modeled independently from sovereignty and independently from governance scope |
| `ADH-2026-025` | Existing infrastructure may be registered through `InfrastructureStack`; all lifecycle action is directed to an approved `ExecutionTarget`, including external or federated realization endpoints; no mandatory `ResourcePool` exists |
| `ADH-2026-026` | `SovereigntyAssessment` and `PlacementDecision` adopt FEATURE-0013 `DecisionRecord` profiles and typed results |
| `ADH-2026-027` | Published product, runtime, governance, sovereignty and policy definitions are immutable by version |
| `ADH-2026-028` | A service may be composed across multiple execution targets; cross-target relationships are evaluated explicitly |
| `ADH-2026-029` | Customer placement visibility is a safe immutable projection, not access to provider topology or credentials |
| `ADH-2026-030` | A personal customer receives a Personal Organization, default Tenant and default Project; governance is simplified, not bypassed |
| `ADH-2026-031` | Countries, sectors, service families and backends enter through governed versioned contracts—not core workflow conditionals |
| `ADH-2026-032` | `ServiceOffering` supersedes the global `ServiceClass` as the CloudPlatform product; it references a reusable versioned `ServiceTypeDefinition` |
| `ADH-2026-033` | Governance controls compose into one `EffectiveGovernanceContext`; sovereignty remains a separate explicit profile and assessment |
| `ADH-2026-034` | `ServiceBinding` is the separately revocable, SecretRef-only managed-service access boundary |
| `ADH-2026-035` | `ServiceRelationshipDefinition` and `ServiceRelationship` model cross-service semantics and lifecycle while plugins/adapters alone own implementation-native translation |
| `ADH-2026-036` | Sovrunn installations and releases use an operator-facing, independently recoverable platform lifecycle boundary with one external accepted-intent authority and atomic Operation-first activation; external infrastructure remains externally operated |
| `ADH-2026-037` | CloudPlatform ownership, CloudProvider participation and one-provider-per-installation are separate governed boundaries |
| `ADH-2026-038` | Platform sovereignty evaluates the SovrunnInstallation and every dependency able to control, observe, change, decrypt or recover it |
| `ADH-2026-039` | Release transitions are explicit directed compatibility contracts with typed references, executable recovery modes and irreversible checkpoints |
| `ADH-2026-040` | Infrastructure maintenance separates provider notice authority from Sovrunn target state and uses epochs, fences and mandatory requalification |
| `ADH-2026-041` | Superseded by ADH-2026-045/DEC-0059: alpha migration is not implemented as a coordinated cutover because there is no live alpha state to convert |

## 5. CloudPlatform product and provider-supply contracts

### 5.1 Contract matrix

| Concept / API kind | Purpose and ownership | Profile, scope and writers | Mutability and lifecycle | Boundary, relationships and invariants |
|---|---|---|---|---|
| `CloudPlatform` | Named customer-facing cloud product and governance boundary | MR; Platform scope (no Organization scope/reference); bootstrap/cloud-owner administrator owns spec; platform controller owns status | Identity, scope, and immutable `spec.ownerRegistration` (legalName, registrationIdentifier, jurisdictionCode) immutable; `spec.description` and approved operating configuration mutable via PATCH only; no PUT/DELETE; no lifecycle action or deletion exposed | CF brand/product projection plus operator canonical view; owns catalog and enrollments, never owns participating providers or customer Organizations; `ADH-2026-037`/`045` |
| `CloudProvider` | Independent infrastructure/execution supply operator | MR; Platform scope; supply administrator owns spec; provider controller owns status | Identity/scope immutable; `spec.displayName` and non-empty `spec.operatingMarkets[]` mutable via PATCH only; no PUT/DELETE; no lifecycle action or deletion exposed | OF canonical resource plus optional CF provider-summary projection; registers topology/targets and maintenance, never owns customer catalog or Organizations; no deployment reference or credential; `ADH-2026-020`/`037`/`045` |
| `CloudProviderParticipation` | Governed supply relationship between one CloudPlatform and one CloudProvider | MR; CloudPlatform scope; cloud-owner contracting authority and delegated provider authority act through explicit lifecycle actions; participation controller owns status | PlatformRef/providerRef immutable; no mutable spec fields; lifecycle changes (accept/reject/withdraw/expire/suspend/resume/request-release/accept-release/decline-release) occur only through explicit actions, never PATCH; terminate rather than erase; RETAIN summary | OF/GO; one non-terminal participation per exact CloudPlatform+CloudProvider pair; effective `Active` only when accepted and both `platformSuspended`/`providerSuspended` holds are false; protected agreement data never enters customer projection; `ADH-2026-037`/`045` |
| `ServiceRegion` | Customer-selectable service geography | MR; CloudPlatform scope; product administrator owns spec; availability controller owns status | Mappings may change through generation; existing instances preserve placement record; RETIRE then RESTRICT | CF projection and owner/provider views; maps through eligible participations to one or more HostingLocations/Datacenters/ExecutionTargets; must not claim sovereignty |
| `ServicePortfolio` | Curated general, technology, sovereign or industry product bundle | VD; CloudPlatform scope; authorized product publisher owns spec | Draft mutable; published version immutable; RETIRE | CF published projection; contains offerings and required profiles; `ADH-2026-023`/`027` |
| `ServiceTypeDefinition` | Implementation-neutral extension contract for one managed-service family | VD; Platform, CloudPlatform or approved publisher scope; approved service-contract publisher owns spec | Draft mutable; published version immutable; RETIRE; referenced versions retained | IE/PF with CF schema/documentation projection; defines parameter, action, status, binding, meter, runtime-compatibility and upgrade schemas; no provider-native payloads; `ADH-2026-031`/`032` |
| `ServiceRelationshipDefinition` | Implementation-neutral extension contract for a directed relationship between service families | VD; Platform, CloudPlatform or approved publisher scope; approved service-contract publisher owns spec | Draft mutable; published version immutable; RETIRE; referenced versions retained | IE/PF with CF selectable summary; defines source/target type selectors, parameter schema, requirement resolver, DecisionRecord profiles, binding compatibility and lifecycle defaults; `ADH-2026-035` |
| `ServiceOffering` | Managed application, data, AI or governed-IaaS outcome | VD; CloudPlatform scope; service product owner publishes | Draft mutable; published immutable; RETIRE | CF; supersedes global `ServiceClass` as product identity, references exactly one ServiceTypeDefinition version, is referenced by portfolios and owns a plan set; `ADH-2026-032`/`037` |
| `ServicePlan` | Selectable service configuration and commercial/operational tier | VD; CloudPlatform scope; service product owner publishes | Draft mutable; published immutable; RETIRE; active instances pin version | CF; references one offering, runtime profile and supported regions/profiles; parameter schema owns version choices and limits |
| `EntitlementPackage` | Reusable template for catalog, region, support and initial quota grants | VD; CloudPlatform scope; commercial/governance administrator publishes | Published version immutable; RETIRE | Cloud-owner view with customer summary; applied through enrollment and materialized as entitlements; package is not itself authorization |

### 5.2 Skeletal YAML

```yaml
apiVersion: core.sovrunn.io/v1alpha1
kind: CloudPlatform
metadata:
  name: nic-cloud
  displayName: NIC Cloud
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Platform, name: sovrunn}
spec:
  description: National Informatics Centre sovereign cloud
  ownerRegistration:
    legalName: National Informatics Centre
    registrationIdentifier: NIC-GOV-IN-0001
    jurisdictionCode: IN
status:
  phase: Active
---
apiVersion: core.sovrunn.io/v1alpha1
kind: CloudProvider
metadata:
  name: yotta
  displayName: Yotta
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Platform, name: sovrunn}
spec:
  operatingMarkets: [IN]
status:
  phase: Active
---
apiVersion: governance.sovrunn.io/v1alpha1
kind: CloudProviderParticipation
metadata:
  name: nic-cloud-yotta-production
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
spec:
  cloudPlatformRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
  cloudProviderRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
  environment: production
  permittedHostingLocationRefs: [india-delhi-ncr, india-pune]
  providerSelectionModes: [CustomerSelected, PlatformAssigned]
status:
  phase: Active
---
apiVersion: catalog.sovrunn.io/v1alpha1
kind: ServiceRegion
metadata:
  name: nic-cloud-india-north
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
spec:
  displayName: India North
  hostingLocationRefs: [india-delhi-ncr]
  executionTargetRefs: [yotta-openshift-gov-primary]
status:
  availability: Available
---
apiVersion: catalog.sovrunn.io/v1alpha1
kind: ServicePortfolio
metadata:
  name: government-cloud-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
spec:
  version: v1
  portfolioClass: industry
  offeringRefs: [managed-postgresql-v1]
  publicationState: Published
---
apiVersion: catalog.sovrunn.io/v1alpha1
kind: ServiceTypeDefinition
metadata:
  name: managed-postgresql-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Platform, name: sovrunn}
spec:
  version: v1
  parameterSchemaRef: managed-postgresql-parameters-v1
  actionSchemaRef: managed-postgresql-actions-v1
  statusProjectionSchemaRef: managed-postgresql-status-v1
  bindingSchemaRef: managed-postgresql-binding-v1
  meterSchemaRef: managed-postgresql-meters-v1
  upgradeContractRef: managed-postgresql-upgrades-v1
  compatibleRuntimeProfileSelector: database.postgresql/v1
  publicationState: Published
---
apiVersion: catalog.sovrunn.io/v1alpha1
kind: ServiceRelationshipDefinition
metadata:
  name: private-service-consumption-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Platform, name: sovrunn}
spec:
  version: v1
  allowedSourceServiceTypeRefs: [compute.virtual-machine-v1, application.runtime-v1]
  allowedTargetServiceTypeRefs: [managed-postgresql-v1]
  parameterSchemaRef: private-service-consumption-parameters-v1
  requirementResolverRef: private-service-consumption-resolver-v1
  decisionProfileRefs: [relationship-eligibility-v1]
  supportedBindingTypeRefs: [database-client-identity-v1]
  lifecycleDefaults:
    sourceReadiness: RequireRelationshipReady
    targetDeletion: RestrictWhileRequired
  publicationState: Published
---
apiVersion: catalog.sovrunn.io/v1alpha1
kind: ServiceOffering
metadata:
  name: managed-postgresql-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
spec:
  version: v1
  serviceTypeDefinitionRef: managed-postgresql-v1
  planRefs: [postgres-government-ha-v1]
  publicationState: Published
---
apiVersion: catalog.sovrunn.io/v1alpha1
kind: ServicePlan
metadata:
  name: postgres-government-ha-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
spec:
  offeringRef: managed-postgresql-v1
  runtimeProfileRef: postgres-ha-runtime-v1
  parameterSchemaRef: postgres-government-ha-parameters-v1
  supportedServiceRegionRefs: [nic-cloud-india-north]
  publicationState: Published
---
apiVersion: governance.sovrunn.io/v1alpha1
kind: EntitlementPackage
metadata:
  name: government-production-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
spec:
  version: v1
  portfolioRefs: [government-cloud-v1]
  planRefs: [postgres-government-ha-v1]
  serviceRegionRefs: [nic-cloud-india-north]
  publicationState: Published
```

## 6. CloudPlatform–customer relationship contracts

### 6.1 Contract matrix

| Concept / API kind | Purpose and ownership | Profile, scope and writers | Mutability and lifecycle | Boundary, relationships and invariants |
|---|---|---|---|---|
| `CloudEnrollment` | Durable commercial/governance relationship between one CloudPlatform and one customer Organization | MR; CloudPlatform scope; enrollment service/cloud-owner administrator owns spec; enrollment controller owns status | CloudPlatformRef and OrganizationRef immutable; commercial references/versioned defaults mutable by governed operation; terminate rather than erase; RETAIN summary | Cloud-owner canonical plus CF projection; M:N bridge; never owns Organization hierarchy; provider assignment remains separate; no cross-scope existence disclosure; `ADH-2026-021`/`037` |
| `ServiceEntitlement` | Actual grant of portfolios, offerings, plans and regions under an enrollment | MR; CloudPlatform scope; entitlement administrator/controller owns spec; entitlement evaluator owns status | Grants/validity mutable with impact analysis; revocation is explicit; RETAIN | Cloud-owner view with CF effective projection; requires active enrollment; descendants may narrow but never enlarge grant; `ADH-2026-022` |
| `QuotaPolicy` | Independent consumption/rate/cost limit | MR; CloudPlatform, Organization, Tenant or Project scope; scope policy owner writes spec; quota controller owns status | Limits and validity mutable; deletion restores inherited effective quota only after impact preview | CF for effective quota, operator view for full policy; target must be same scope or descendant; must not imply entitlement; `ADH-2026-022` |

### 6.2 Skeletal YAML

```yaml
apiVersion: governance.sovrunn.io/v1alpha1
kind: CloudEnrollment
metadata:
  name: department-of-posts-nic-cloud
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
spec:
  cloudPlatformRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
  organizationRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Organization, name: department-of-posts}
  entitlementPackageRefs: [government-production-v1]
status:
  phase: Active
---
apiVersion: governance.sovrunn.io/v1alpha1
kind: ServiceEntitlement
metadata:
  name: department-of-posts-postgresql
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
spec:
  enrollmentRef: department-of-posts-nic-cloud
  planRefs: [postgres-government-ha-v1]
  serviceRegionRefs: [nic-cloud-india-north]
status:
  effective: "True"
---
apiVersion: governance.sovrunn.io/v1alpha1
kind: QuotaPolicy
metadata:
  name: department-of-posts-postgresql-quota
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
spec:
  enrollmentRef: department-of-posts-nic-cloud
  limits:
    maximumInstances: 20
    maximumTotalStorageGiB: 20480
```

## 7. Customer governance contracts

### 7.1 Contract matrix

| Concept / API kind | Purpose and ownership | Profile, scope and writers | Mutability and lifecycle | Boundary, relationships and invariants |
|---|---|---|---|---|
| `Organization` | Durable customer ownership and governance boundary for a person or legal entity | MR; Platform scope; registry/customer administrator owns spec; organization controller owns status | Type and verified identity changes governed; DECOMMISSION; RESTRICT descendants/enrollments | CF; owns tenants and optional units; personal signup creates Personal Organization; `ADH-2026-030` |
| `OrganizationUnit` | Optional delegation/grouping inside Organization | MR; Organization or OrganizationUnit scope; organization administrator owns spec; governance controller owns status | Scope movement requires recreate/migrate; RESTRICT children | CF; never mandatory and does not create independent tenancy or enrollment |
| `Tenant` | Primary customer security, service and data-isolation boundary | MR; Organization or OrganizationUnit scope; organization administrator owns spec; tenant controller owns status | Scope immutable; DECOMMISSION; RESTRICT projects/services | CF; default tenant created for simple onboarding |
| `Project` | Everyday workload/application/environment workspace | MR; Tenant scope; authorized customer owns spec; project controller owns status | Scope immutable; RESTRICT by default; CASCADE-EXPLICIT only with impact preview | CF; owns ServiceInstances; no provider-scope parentage |
| `ServiceInstance` | Customer desired state for one managed service outcome | MR; Project scope; customer owns allowed spec fields; service controller owns status | Plan version pinned; fields mutable only per plan; major lifecycle changes use Operation; delete creates Delete Operation and finalizes after backend cleanup | CF; references entitled plan and embedded placement intent; backend handles prohibited |
| `ServiceRelationship` | Directed governed relationship between two independently meaningful ServiceInstances | MR; source Project scope; authorized customer/provider owns allowed spec; relationship controller owns status | Source, target and definition version immutable; parameter changes governed; deletion/remediation through Operation; required relationships may restrict target deletion | CF safe intent/status plus IE/PF realization; exact ServiceRelationshipDefinition version, bounded parameters and authorization at both endpoints required; cross-Organization and cross-CloudProvider use additionally approved sharing, supply or federation contracts and are denied until those contracts exist; no native handles; `ADH-2026-035` |
| `ServiceBinding` | Per-consumer access relationship to one managed ServiceInstance | MR; Project scope; authorized customer/workload owns allowed intent; binding controller owns status | ServiceInstance and consumer references immutable; rotate/revoke through Operation; delete/revoke removes access before finalization | CF to authorized principals; status contains safe connection metadata and typed SecretRef only; raw credentials prohibited; `ADH-2026-034` |

### 7.2 Skeletal YAML

```yaml
apiVersion: core.sovrunn.io/v1alpha1
kind: Organization
metadata:
  name: department-of-posts
spec:
  organizationType: GovernmentDepartment
status:
  phase: Active
---
apiVersion: core.sovrunn.io/v1alpha1
kind: OrganizationUnit
metadata:
  name: postal-technology
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Organization, name: department-of-posts}
spec:
  description: Optional delegated technology unit
---
apiVersion: core.sovrunn.io/v1alpha1
kind: Tenant
metadata:
  name: department-of-posts-production
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Organization, name: department-of-posts}
spec:
  isolationClass: production
---
apiVersion: core.sovrunn.io/v1alpha1
kind: Project
metadata:
  name: postal-applications
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Tenant, name: department-of-posts-production}
spec:
  description: Postal application managed services
---
apiVersion: services.sovrunn.io/v1alpha1
kind: ServiceInstance
metadata:
  name: postal-postgresql
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Project, name: postal-applications}
spec:
  servicePlanRef: postgres-government-ha-v1
  parameters:
    engineVersion: "17"
    storageGiB: 1000
  placementIntent:
    requiredServiceRegionRefs: [nic-cloud-india-north]
    placementProfileRef: single-region-ha-v1
    sovereigntyProfileRef: india-government-data-v1
status:
  phase: Pending
---
apiVersion: services.sovrunn.io/v1alpha1
kind: ServiceRelationship
metadata:
  name: application-consumes-postgresql
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Project, name: postal-applications}
spec:
  relationshipDefinitionRef: private-service-consumption-v1
  sourceServiceInstanceRef: application-vm
  targetServiceInstanceRef: postal-postgresql
  required: true
  parameters:
    sourcePlacementProfileRef: application-zone
    targetPlacementProfileRef: data-zone
    relationshipProfileRef: private-low-latency
    maximumLatencyMs: 5
status:
  phase: Evaluating
---
apiVersion: services.sovrunn.io/v1alpha1
kind: ServiceBinding
metadata:
  name: postal-app-postgresql
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Project, name: postal-applications}
spec:
  serviceInstanceRef: postal-postgresql
  consumerRef:
    kind: WorkloadIdentity
    name: postal-application
  accessMode: ReadWrite
status:
  phase: Ready
  endpoint: postal-postgresql.services.sovrunn.example
  secretRef:
    kind: Secret
    name: postal-app-postgresql-credentials
```

### 7.3 Persona realization and authorization contracts

`Persona` is a non-normative description of job intent, responsibilities and user journey. It is not an API resource, identity kind, role, scope or authorization input. Persona documentation may map a journey to one or more default role templates; actual authority exists only through the contracts below.

| Concept / API kind | Purpose and ownership | Profile, scope and writers | Mutability and lifecycle | Boundary, relationships and invariants |
|---|---|---|---|---|
| `PrincipalRef` | Embedded stable reference to an already-authenticated Human, Workload or System identity | EV only; identity lifecycle remains external; resolver validates issuer/subject and current identity state | Immutable value; identity change produces a new reference | Email/display name are not authority; unresolved identity fails closed |
| `AccessGroup` | Non-authenticating, access-only authorization subject | MR; Organization or CloudProvider scope; authorized administrator writes intent; AccessGroup controller writes status | Owner/display intent may update; suspend/resume/retire explicit | Direct auditable principal members only; not a ScopeKind; external, nested and dynamic groups grant nothing |
| `Membership` | Records Organization/CloudProvider belonging or direct same-scope AccessGroup membership without granting permission or workflow eligibility | MR; Organization or CloudProvider scope; administrator/trusted provisioner submits intent; membership controller alone validates grant envelope and writes canonical state | Relationship immutable; synchronize only system-owned provenance/freshness; suspend/reactivate/revoke; Guest expiry terminal | Exact relationship-key uniqueness; Standard or Guest; an AccessGroup assignment-effect expansion derives `MembershipEnabledGrantEnvelope` and passes complete non-synthesized action/scope/target/resource/validity dominance; privilege reduction remains available |
| `RoleDefinition` | Versioned permission template containing only registered canonical Sovrunn actions | VD; Platform, Organization or CloudProvider scope; security administrator drafts; publication controller owns lifecycle status | Draft mutable; published version immutable; supersede, emergency suspend/restore, retire; retained while referenced | Privilege classification is monotonic; no backend IAM action, persona, arbitrary condition or wildcard authority |
| `RoleAssignment` | Grants one exact published RoleDefinition version to one `RoleHolderRef` at one registered scope, optionally narrowed to one exact resource | MR; seven registered resource scopes; accepted workflows submit exact intents; RoleAssignment controller alone writes spec/status/lifecycle and every revoke/replace/expiry/review-due effect | Immutable published spec; explicit revoke/expiry/review effects; replacement creates a new assignment | `Standing` or `TimeBound`; exact Membership/responsibility/review pins; grantor ceiling includes per-action scope/target/resource-reach/validity/delegation dominance; no cross-assignment synthesis or workflow write |
| `PrivilegedAccessRequest` | Requests one justified individual JIT or break-glass temporary grant | LRO; seven registered scopes; requester writes accepted intent; privileged controller owns state | Immutable request, activation deadline and duration; cancel/block/activate/expire/revoke retained | JIT and break-glass activation require fresh phishing-resistant AAL2-or-higher assurance, non-empty exact-Human PrincipalRef-only `EligibilityRef` requester list, separate submit authority, separation of duties, no-grace expiry and linked review where required |
| `ApprovalPolicy` | Versioned bounded stages, eligibility, quorum, expiry and separation-of-duties policy | VD; Platform, Organization or CloudProvider scope; security/governance administrator drafts; publication controller owns lifecycle | Draft mutable; published immutable; supersede/retire; retained while referenced | Non-empty exact-Human PrincipalRef-only `EligibilityRef` approver list with separate decide authority and decision/effect recheck; never authorizes break-glass bypass; operation-local `ApprovalRequirement` selects one exact version |
| `ApprovalRequest` | Retains approval evidence for one bounded FEATURE-0018 intent or immutable proposal | LRO; seven registered scopes; originating controller materializes; approval controller owns result | Immutable subject/proposal; exactly one terminal decision; cancel/expire retained | Approval is current at downstream publication, then immutable provenance rather than continuing bearer authority |
| `AccessReview` | Reviews an immutable access snapshot with usage and remediation evidence | LRO; seven registered scopes; campaign owner or mandated controller originates; review controller owns state | Snapshot immutable after start; item decisions/completion/remediation retained | StandingCertification, Manual or BreakGlassRetrospective variants; non-empty exact-Human PrincipalRef-only `EligibilityRef` reviewer list with independent action and decision/effect recheck; only exact-version StandingCertification Retain advances `nextReviewDueAt`; beneficiary conflict set includes accountable Humans behind direct Workload/System AccessGroup members and prevents direct or indirect self-certification |
| `ExceptionGrant` | Immutable approved exception grant or linked revocation evidence | IR; seven registered scopes; exception controller is sole writer after approval | Append-only Grant/Revoke evidence; expiry is time-effective without mutation | Exact control/subject/scope/interval; deterministic overlap and narrowing; FEATURE-0020 alone applies it to effective governance |
| `AuthorizationInput` / `AuthorizationResult` | Transient exact authorization request and non-bearer `Allow | Deny` result | TRR; authorization controller writes result; exact provenance retained in required evidence | Evaluation-instant result only; never reusable authority | Grant union plus guardrail/dependency intersection; exact target binding; safe denial and audit-obligation acceptance |

The remaining closed FEATURE-0018 supporting-value inventory is
`RoleHolderRef`, `AssuranceEvidence`, `ApprovalRequirement`, `EligibilityRef`,
`RoleAssignmentProposal`, `MembershipEnabledGrantEnvelope`, `ExceptionProposal`, `UsageEvidenceSummary`,
`GovernanceApplicability`, and `ActionTargetBinding`. Together with
`AuthorizationInput` and `AuthorizationResult`, these are EmbeddedValue or
TransientRequestResult contracts only: no endpoint, independent lifecycle,
storage authority or user journey exists.

Scope applicability is closed per version and consumed generically. AccessGroup
and Membership allow Organization/CloudProvider; RoleDefinition,
ApprovalPolicy and GovernanceProfile allow Platform/Organization/CloudProvider;
RoleAssignment, PrivilegedAccessRequest, AccessReview, ApprovalRequest and
ExceptionGrant allow all seven canonical ScopeKinds. Customer policy may narrow
but not expand this registry. A new ScopeKind or containment rule is a core
architecture change. Exact publication-scope/reference compatibility and the
complete operation/action/target-binding registries remain F18-RD-02/F18-RD-06
authority and must be transcribed without inference into executable schemas.

Canonical processing path:

```text
Persona documentation
    → maps journey to default RoleDefinition templates

Authenticated PrincipalRef
  + canonical Sovrunn action
  + target resource and resolved scope
  + pinned RoleAssignment and RoleDefinition version
  + Membership/responsibility, policy guardrails, assurance and current time
      → AuthorizationResult
          → allow or deny
          → material DecisionRecord profile and/or AuditEvent
```

Authorization must never branch on persona name, UI label, organizational job title, email address, unprovisioned external group claim or backend IAM role. Default roles are packaged conveniences, not permanent enums; providers and customers may create narrower roles but may never grant actions outside their own delegated authority. Applicable grants combine by union and every guardrail/dependency combines by intersection; missing or ambiguous evidence denies.

An AccessGroup Membership assignment-effect expansion is grant-producing even
though Membership itself contributes no action. At its atomic publication boundary,
the Membership controller derives the exact current group-held assignment
action/reach/validity envelope and applies F18-RD-09's Membership-operation,
`roleassignment.grant`, per-action witness, temporal-dominance, no-synthesis,
trusted-provisioner and coherent-snapshot rules. No user-authored Membership
field stores or supplies that authority. Membership cannot establish approver,
reviewer, JIT, or break-glass eligibility. All three eligibility lists use the
same non-empty `EligibilityRef` OR contract containing exact Human
PrincipalRefs only. RoleDefinition, RoleAssignment, AccessGroupRef, Membership,
claim, and other relationships never qualify in v1. Eligibility remains
separate from the required canonical action; exact policy/rule applicability
and that action alone determine scope and target reach.

```yaml
apiVersion: iam.sovrunn.io/v1alpha1
kind: RoleDefinition
metadata:
  name: application-developer-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
spec:
  version: v1
  actions:
    - services.instances.create
    - services.instances.read
    - services.operations.read
    - services.relationships.manage
    - services.bindings.request
  publicationState: Published
---
apiVersion: iam.sovrunn.io/v1alpha1
kind: RoleAssignment
metadata:
  name: developer-001-postal-applications
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Project, name: postal-applications}
spec:
  roleHolderRef:
    principalRef:
      issuer: https://identity.department-of-posts.example
      subject: 2c48198f-74d7-4c41-8d5e-21dd9001d3b7
  roleDefinitionRef: application-developer-v1
  validity:
    mode: TimeBound
    notBefore: 2027-08-03T16:00:00Z
    expiresAt: 2027-08-04T00:00:00Z
```

## 8. Physical topology and execution contracts

### 8.1 Contract matrix

All five resources are CloudProvider-scoped. Their relationships describe registered facts, not security scope, network connectivity, placement eligibility or Sovrunn ownership.

| Concept / API kind | Purpose and ownership | Profile and writers | Mutability and lifecycle | Boundary, relationships and invariants |
|---|---|---|---|---|
| `HostingLocation` | Normalized physical geography including country, subdivision and locality | MR; provider topology administrator owns spec; verification controller owns status | Scope immutable; material geographic correction creates reviewed generation or replacement; RESTRICT references | OF with safe CF projection; country is an attribute, not the entire model; `ADH-2026-024` |
| `Datacenter` | Physical facility or equivalent hosting site | MR; provider topology administrator/spec; topology controller/status | HostingLocationRef immutable; certifications/display facts mutable with provenance; DECOMMISSION then RESTRICT | OF; may be provider- or third-party-operated; registration never implies ownership |
| `FaultDomain` | Meaningful failure-isolation boundary | MR; provider/adapter declares spec; topology controller verifies status | Parent relationships immutable; DECOMMISSION then RESTRICT | OF; must identify type and provenance; sibling relationship never implies connectivity |
| `InfrastructureStack` | Existing execution platform such as OpenShift, Kubernetes, OpenStack, AWS or OCI | MR; provider integration administrator/spec; integration controller/status | Implementation type and topology identity immutable; integration config changes controlled; DECOMMISSION | OF/AF projections; does not own or replace native scheduler; optional when the ExecutionTarget is not infrastructure-backed; `ADH-2026-025` |
| `ExecutionTarget` | Approved adapter-addressable realization boundary for service lifecycle actions | MR; provider integration administrator/spec; qualification/adapter controllers own distinct status conditions | Realization pattern and backing-environment identity immutable; adapter version and credentialRef rotatable; Drain → DECOMMISSION → RESTRICT active placements | OF canonical, AF execution view, CF only through ServicePlacement; backing environment may be infrastructure, external managed API, edge control plane or future federation endpoint; `ADH-2026-025` |
| `InfrastructureMaintenanceNotice` | Scheduled, emergency, updated or cancelled external maintenance notice for one or more targets | MR; CloudProvider scope; delegated InfrastructureOperator/trusted adapter owns spec; notice controller owns status | Target refs immutable after acceptance; schedule/impact updates use generation and authorization; terminal notice retained with epoch/fence history | OF/IE with safe CF impact projection; only target lifecycle controller changes availability; `ADH-2026-040` |

### 8.2 Skeletal YAML

```yaml
apiVersion: topology.sovrunn.io/v1alpha1
kind: HostingLocation
metadata:
  name: india-delhi-ncr
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
spec:
  countryCode: IN
  subdivision: Delhi NCR
---
apiVersion: topology.sovrunn.io/v1alpha1
kind: Datacenter
metadata:
  name: yotta-qualified-dc-a
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
spec:
  hostingLocationRef: india-delhi-ncr
  operatorRef: yotta
---
apiVersion: topology.sovrunn.io/v1alpha1
kind: FaultDomain
metadata:
  name: yotta-dc-a-zone-a
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
spec:
  datacenterRef: yotta-qualified-dc-a
  domainType: InfrastructureZone
  nativeRef: {protectedHandle: zone-a}
---
apiVersion: topology.sovrunn.io/v1alpha1
kind: InfrastructureStack
metadata:
  name: yotta-openshift-government-production
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
spec:
  implementationType: OpenShift
  datacenterRefs: [yotta-qualified-dc-a]
status:
  integrationState: Registered
---
apiVersion: execution.sovrunn.io/v1alpha1
kind: ExecutionTarget
metadata:
  name: yotta-openshift-gov-primary
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
spec:
  realizationPattern: InfrastructurePlatform
  backingEnvironmentRef:
    apiVersion: topology.sovrunn.io/v1alpha1
    kind: InfrastructureStack
    name: yotta-openshift-government-production
  hostingLocationRef: india-delhi-ncr
  datacenterRef: yotta-qualified-dc-a
  faultDomainRefs: [yotta-dc-a-zone-a, yotta-dc-a-zone-b]
  adapterRef: openshift-adapter-v1
  credentialRef: yotta-openshift-gov-primary-credential
status:
  qualification: Qualified
  health: Available
```

### 8.3 External infrastructure lifecycle coordination

Underlying infrastructure lifecycle remains externally owned. `InfrastructureMaintenanceNotice` is a CloudProvider-scoped ManagedResource authored only by a delegated InfrastructureOperator or trusted adapter identity. It references affected targets, a typed MaintenanceWindow, expected target generations, impact, urgency and idempotency key. The notice controller owns notice status; only the ExecutionTargetLifecycleController owns target availability status.

Qualification and availability are independent. Normal availability progresses `Available → DrainPending → Draining → InMaintenance → Requalifying → Available`, with Unavailable, Degraded and Restricted branches. Cancellation or provider completion cannot bypass Requalifying.

Notice acceptance increments `maintenanceEpoch` and creates a fence. Placement decisions/plans pin target UID, generation, qualification/evidence snapshot and epoch. DrainPending and later states deny new placement. Each externally consequential PluginExecution step validates the current fence and fails closed when stale. Requalification refreshes health, facts, connectivity, evidence and sovereignty before availability is restored.

The CloudPlatform owner may impose a separate restriction through existing governance and decision contracts but cannot impersonate the provider's infrastructure operator or clear provider maintenance. Provider maintenance completion cannot clear that owner restriction.

Backend-native maintenance objects remain adapter-facing. Customer projections disclose affected services, timing and safe planned outcomes only when authorized; they do not expose sensitive topology or credentials. Sovrunn coordinates managed-service impact but never reports that it performed the underlying stack upgrade or repair.

### 8.4 InfrastructureMaintenanceNotice skeleton

```yaml
apiVersion: execution.sovrunn.io/v1alpha1
kind: InfrastructureMaintenanceNotice
metadata:
  name: yotta-delhi-maintenance-2026-08-20
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
spec:
  targetRefs: [yotta-openshift-gov-primary]
  maintenanceWindowRef: yotta-delhi-window-2026-08-20
  expectedTargetGenerations: {yotta-openshift-gov-primary: 12}
  impact: Disruptive
  urgency: Scheduled
  idempotencyKey: yotta-delhi-2026-08-20
status:
  phase: Accepted
  maintenanceEpoch: 44
  fenceToken: 01J4YOTTAMAINT0044
```

## 9. Governance contracts

### 9.1 Contract matrix

| Concept / API kind | Purpose and ownership | Profile, scope and writers | Mutability and lifecycle | Boundary, relationships and invariants |
|---|---|---|---|---|
| `GovernanceProfile` | v1 composition of only FEATURE-0018-owned privileged-access, access-review, exception and audit rule references/constraints | VD; Platform, CloudProvider or Organization scope; governance administrator drafts; publication controller owns lifecycle | Draft mutable; published immutable; supersession blocks new selection while preserving valid pins; retirement blocks new/current-rule use; retained for audit | GO canonical with safe summaries; no security, sovereignty, placement, backup, cost or other future-domain semantics; FEATURE-0020 owns assignment and effective resolution; DEC-0060/F18-RD-18 |
| `ProfileAssignment` | Applies a specific profile version to a target scope/resource | MR; same scope as target; authorized scope policy administrator owns spec; policy resolver owns status | ProfileRef and targetRef immutable; validity mutable by governed operation; RESTRICT active decision dependencies | GO/OF; allowed sets intersect, prohibitions accumulate, strongest minimum and shortest evidence age win |
| `EffectiveGovernanceContext` | Immutable resolved governance, entitlement, restriction and approved-exception snapshot for one decision or operation | IR; same canonical scope as subject; policy resolver is sole writer | Append-only; supersession creates new record; RETAIN | IE/GO; supersedes the planned `EffectivePolicyContext` term; contains references/versions, not copied secrets; exact contextRef is mandatory for authoritative decisions; `ADH-2026-033` |

### 9.2 Skeletal YAML

```yaml
apiVersion: policy.sovrunn.io/v1alpha1
kind: GovernanceProfile
metadata:
  name: government-production-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
spec:
  version: v1
  approvalPolicyRefs: [government-production-approval-v1]
  privilegedAccessRules: [government-production-privileged-v1]
  accessReviewRules: [government-production-standing-review-v1]
  exceptionRules: [government-production-exception-v1]
  auditRequirements: [government-production-iam-audit-v1]
  publicationState: Published
---
apiVersion: policy.sovrunn.io/v1alpha1
kind: ProfileAssignment
metadata:
  name: posts-production-government-profile
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Tenant, name: department-of-posts-production}
spec:
  profileRef: government-production-v1
  targetRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Tenant, name: department-of-posts-production}
status:
  effective: "True"
---
apiVersion: policy.sovrunn.io/v1alpha1
kind: EffectiveGovernanceContext
metadata:
  name: postal-postgresql-context-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Project, name: postal-applications}
record:
  subjectRef: {apiVersion: services.sovrunn.io/v1alpha1, kind: ServiceInstance, name: postal-postgresql}
  resolvedProfileRefs: [government-production-v1, india-government-data-v1]
  effectiveControlsRef: controls-snapshot-8f42
```

## 10. Sovereignty contracts

### 10.1 Contract matrix

| Concept / API kind | Purpose and ownership | Profile, scope and writers | Mutability and lifecycle | Boundary, relationships and invariants |
|---|---|---|---|---|
| `SovereigntyProfile` | Desired multidimensional sovereignty outcome | VD; Platform, CloudPlatform, CloudProvider or Organization scope; qualified publisher owns spec | Draft mutable; published immutable; RETIRE | GO canonical with CF selectable summary; dimensions come from governed registry; `ADH-2026-027`/`031`/`033` |
| `RegulatoryPolicyBundle` | Signed, effective-dated executable interpretation of applicable sources | VD; Platform, CloudPlatform or qualified-authority scope; qualified policy authority publishes | Approved version immutable; supersession links new version; RETAIN while referenced/legal hold | GO only, with safe summary; explicitly not legislation or legal advice; source and approval references mandatory |
| Concept `SovereigntyFacts`; API kind `SovereigntyFactSet` | Current normalized facts about one subject | OER; same scope as subject; authorized adapter/collector writes observations; fact controller owns freshness status | Observation updates only; every update carries observedAt/provenance; stale becomes Unknown, never silently valid | IE/AF/GO; no customer-written facts; no facts inferred from geography alone |
| `EvidenceRecord` | Immutable proof supporting facts or decisions | IR; same scope as subject; authorized evidence collector sole writer | Append-only; correction/supersession linked; RETAIN under evidence/legal policy | GO; evidence payload may be external immutable reference; integrity, collection time, validity and custody required |
| `PlatformDependencySnapshot` | Complete dependency closure for one SovrunnInstallation revision and assessment | IR; Platform scope; platform dependency resolver sole writer | Immutable; any dependency or relevant version change creates a new snapshot; RETAIN with decisions and lifecycle history | GO/POF safe summary; includes every system or actor able to control, observe, change, decrypt, execute or recover the installation; SecretRefs/identifiers only; `ADH-2026-038` |
| Concept `SovereigntyAssessment`; API implementation `DecisionRecord` with `sovereignty-assessment` profile | Authoritative evidence-backed conclusion for one subject and context | IR through FEATURE-0013; subject scope; registered evaluator/decision service writes | Final and immutable; changed evidence/policy creates new decision; RETAIN | GO canonical plus CF/OF decision projections; must reference EffectiveGovernanceContext, fact/evidence snapshot and typed result; `ADH-2026-026` |

### 10.2 Skeletal YAML

```yaml
apiVersion: sovereignty.sovrunn.io/v1alpha1
kind: SovereigntyProfile
metadata:
  name: india-government-data-v1
spec:
  version: v1
  dimensions:
    workloadExecutionCountries: [IN]
    primaryDataCountries: [IN]
    backupDataCountries: [IN]
    maximumEvidenceAge: 24h
  publicationState: Published
---
apiVersion: sovereignty.sovrunn.io/v1alpha1
kind: RegulatoryPolicyBundle
metadata:
  name: india-government-cloud-baseline-2026-08
spec:
  version: 2026-08
  effectiveFrom: 2026-08-01T00:00:00Z
  sourceRefs: [approved-government-cloud-requirements]
  approvalRef: qualified-policy-approval-2026-08
  publicationState: Published
---
apiVersion: sovereignty.sovrunn.io/v1alpha1
kind: SovereigntyFactSet
metadata:
  name: yotta-gov-primary-facts
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
provenance:
  collectorRef: yotta-target-fact-collector-v1
  observedAt: 2026-08-02T08:00:00Z
status:
  subjectRef: {apiVersion: execution.sovrunn.io/v1alpha1, kind: ExecutionTarget, name: yotta-openshift-gov-primary}
  facts:
    workloadExecutionCountry: IN
    encryptionKeyCustodyCountry: IN
  evidenceRefs: [yotta-gov-key-custody-evidence-001]
---
apiVersion: evidence.sovrunn.io/v1alpha1
kind: EvidenceRecord
metadata:
  name: yotta-gov-key-custody-evidence-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
record:
  subjectRef: {apiVersion: execution.sovrunn.io/v1alpha1, kind: ExecutionTarget, name: yotta-openshift-gov-primary}
  evidenceType: KeyCustodyAttestation
  sourceRef: approved-key-management-collector
  collectedAt: 2026-08-02T08:00:00Z
  validUntil: 2026-08-03T08:00:00Z
  integrityRef: signed-evidence-envelope-8f42
---
apiVersion: decisions.sovrunn.io/v1alpha1
kind: DecisionRecord
metadata:
  name: postal-postgresql-sovereignty-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Project, name: postal-applications}
record:
  profileRef: sovereignty-assessment-v1
  subjectRef: {apiVersion: services.sovrunn.io/v1alpha1, kind: ServiceInstance, name: postal-postgresql}
  effectiveContextRef: postal-postgresql-context-001
  evidenceSnapshotRef: postal-postgresql-evidence-001
  authority: Authoritative
  typedResult:
    outcome: Satisfied
    validUntil: 2026-08-03T08:00:00Z
```

### 10.3 Platform-sovereignty composition

```yaml
apiVersion: sovereignty.sovrunn.io/v1alpha1
kind: PlatformDependencySnapshot
metadata:
  name: nic-cloud-yotta-platform-dependencies-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Platform, name: sovrunn}
record:
  installationRef: nic-cloud-yotta-production
  installationGeneration: 7
  dependencies:
    - {role: LifecycleAuthority, subjectRef: nic-cloud-yotta-lifecycle-repository}
    - {role: ArtifactRegistry, subjectRef: sovrunn-production-registry}
    - {role: KeyManagement, subjectRef: yotta-sovrunn-kms}
    - {role: RecoveryRepository, subjectRef: yotta-sovrunn-recovery}
  factSetRefs: [nic-cloud-yotta-platform-facts-001]
  evidenceSnapshotRef: nic-cloud-yotta-platform-evidence-001
---
apiVersion: decisions.sovrunn.io/v1alpha1
kind: DecisionRecord
metadata:
  name: nic-cloud-yotta-platform-sovereignty-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Platform, name: sovrunn}
record:
  profileRef: platform-sovereignty-assessment-v1
  subjectRef: {apiVersion: platform.sovrunn.io/v1alpha1, kind: SovrunnInstallation, name: nic-cloud-yotta-production}
  dependencySnapshotRef: nic-cloud-yotta-platform-dependencies-001
  regulatoryPolicyBundleRefs: [india-government-cloud-baseline-2026-08]
  authority: Authoritative
  typedResult:
    outcome: Satisfied
    validUntil: 2026-08-03T08:00:00Z
```

The snapshot includes all dependencies able to control, observe, change, decrypt, execute or recover the installation. Missing or stale mandatory evidence denies a satisfactory assessment. When a service profile requires platform sovereignty, the current satisfactory platform assessment is a placement and execution prerequisite; any dependency, administrator path, key-custody, repository, processor, location or evidence change triggers reassessment.

## 11. Customer placement contracts

### 11.1 Contract matrix

| Concept / API kind | Purpose and ownership | Profile, scope and writers | Mutability and lifecycle | Boundary, relationships and invariants |
|---|---|---|---|---|
| `ServicePlacementProfile` | CloudPlatform-packaged HA, latency, co-location and DR outcome | VD; CloudPlatform scope; reliability/product owner publishes | Published version immutable; RETIRE | CF summary and IE definition; references outcome constraints, never raw scheduler objects |
| `ServicePlacementIntent` | Customer required/preferred/prohibited/provider-managed placement choices | EV within ServiceInstance spec; customer owns allowed values | Mutable only while plan/lifecycle permits; change may trigger new evaluation and Operation | CF; references entitled regions, optional exposed datacenters, placement and sovereignty profiles; no ExecutionTarget IDs |
| `ServicePlacement` | Customer-safe immutable projection of actual placement, resilience and sovereignty result | IR; Project scope; placement projection controller sole writer | New placement/reassessment creates a new record; ServiceInstance status points to current record; RETAIN | CF; may show region, optional datacenter and outcome summaries; must exclude credentials, protected handles and sensitive topology; `ADH-2026-029` |

### 11.2 Skeletal YAML

```yaml
apiVersion: placement.sovrunn.io/v1alpha1
kind: ServicePlacementProfile
metadata:
  name: single-region-ha-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
spec:
  version: v1
  minimumFaultDomains: 2
  runtimeDataMaximumLatencyMs: 5
  backupRelationship: IndependentApprovedSite
  publicationState: Published
---
# ServicePlacementIntent is embedded in ServiceInstance; it is not a resource.
placementIntent:
  requiredServiceRegionRefs: [nic-cloud-india-north]
  preferredCloudProviderRefs: [yotta]
  preferredDatacenterRefs: [yotta-qualified-dc-a]
  placementProfileRef: single-region-ha-v1
  sovereigntyProfileRef: india-government-data-v1
---
apiVersion: placement.sovrunn.io/v1alpha1
kind: ServicePlacement
metadata:
  name: postal-postgresql-placement-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Project, name: postal-applications}
record:
  serviceInstanceRef: postal-postgresql
  serviceRegionRef: nic-cloud-india-north
  selectedCloudProviderRef: yotta
  selectedParticipationRef: nic-cloud-yotta-production
  displayedDatacenterRefs: [yotta-qualified-dc-a]
  placementProfileRef: single-region-ha-v1
  sovereigntyOutcome: Satisfied
  placementDecisionRef: postal-postgresql-placement-decision-001
```

## 12. Internal service-realization contracts

### 12.1 Contract matrix

| Concept / API kind | Purpose and ownership | Profile, scope and writers | Mutability and lifecycle | Boundary, relationships and invariants |
|---|---|---|---|---|
| `ServiceRuntimeProfile` | Versioned implementation-neutral realization of a plan and ServiceTypeDefinition | VD; Platform or CloudPlatform scope; service plugin/product owner publishes | Published immutable; RETIRE; active plans pin version | IE/PF; declares realization pattern, compatible service type/plugin versions and optional composition reference; provider-specific eligibility is a separate participation realization mapping; no backend objects |
| `ServiceCompositionDefinition` | Optional reusable component/dependency graph for a composite managed outcome | VD; Platform or CloudPlatform scope; approved service architect publishes | Draft mutable; published immutable; RETIRE; active runtime profiles pin version | IE/PF; declares component roles, pinned ServiceRelationshipDefinition dependencies, lifecycle order, failure and compensation policy; independently governed components become child ServiceInstances and their meaningful dependency edges become ServiceRelationships; ServiceBinding remains a separate access contract |
| Concept `ServiceRequirements`; API kind `ServiceRequirementSet` | Resolved requirements for one service revision and its applicable relationships | IR; Project scope; requirements resolver sole writer | Immutable; changed request/context/relationship creates new record; RETAIN with service history | IE; derived from exact plan/runtime/request/relationship/governance versions; finite typed requirements, no arbitrary property bag |
| Concept `PlacementDecision`; API implementation `DecisionRecord` with `placement-decision` profile | Authoritative selection of one or more compatible targets | IR through FEATURE-0013; Project scope; placement decision service writes | Final immutable; reassessment creates new decision; RETAIN | GO canonical plus IE/CF projections; references candidates, requirements, sovereignty decision, selected targets, reasons and obligations; `ADH-2026-026`/`028` |
| `ServiceDeploymentPlan` | Immutable executable DAG realizing one selected placement | IR; Project scope; deployment planner sole writer | Immutable; change creates new revision; RETAIN while operation/audit references exist | IE/PF; references exact placement decision and plugin/adapter versions; contains secret refs only |
| `Operation` | Customer/operator-visible long-running lifecycle action | LRO; scope must equal target governance scope; requester owns accepted request; operation controller owns status | Request immutable after acceptance; status progresses to terminal; cancellation explicit; RETAIN per operations policy | CF/PF projections; existing FEATURE-0012 contract reused; idempotency and target/scope match required |
| `PluginExecution` | Internal execution of one plan step through a versioned plugin/adapter/target | LRO; same scope as Operation target; orchestrator writes request; plugin executor owns status | Request immutable; retry attempt creates linked execution when required; RETAIN operational history | PF/AF; least-privileged target credentials, no policy bypass, stable result codes |

### 12.2 Skeletal YAML

```yaml
apiVersion: runtime.sovrunn.io/v1alpha1
kind: ServiceRuntimeProfile
metadata:
  name: postgres-ha-runtime-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudPlatform, name: nic-cloud}
spec:
  version: v1
  serviceTypeDefinitionRef: managed-postgresql-v1
  realizationPattern: ProvisionedPlatformService
  servicePluginRef: postgresql-management-plugin-v1
  componentRoles: [PrimaryDatabase, StandbyDatabase, Backup, Monitoring]
  publicationState: Published
---
# Optional and used only by composite offerings.
apiVersion: runtime.sovrunn.io/v1alpha1
kind: ServiceCompositionDefinition
metadata:
  name: governed-digital-application-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: CloudProvider, name: yotta}
spec:
  version: v1
  components:
    - {role: application, serviceTypeDefinitionRef: application-runtime-v1}
    - {role: database, serviceTypeDefinitionRef: managed-postgresql-v1, independentlyManaged: true}
    - {role: messaging, serviceTypeDefinitionRef: managed-kafka-v1, independentlyManaged: true}
  dependencies:
    - {from: application, to: database, relationshipDefinitionRef: private-service-consumption-v1, required: true}
    - {from: application, to: messaging, relationshipDefinitionRef: event-publication-v1, required: true}
  failurePolicy: CompensateCreatedComponents
  publicationState: Published
---
apiVersion: runtime.sovrunn.io/v1alpha1
kind: ServiceRequirementSet
metadata:
  name: postal-postgresql-requirements-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Project, name: postal-applications}
record:
  serviceInstanceRef: postal-postgresql
  servicePlanRef: postgres-government-ha-v1
  runtimeProfileRef: postgres-ha-runtime-v1
  governanceContextRef: postal-postgresql-context-001
  requirements:
    engineMajorVersion: "17"
    minimumFaultDomains: 2
    backupRelationship: IndependentApprovedSite
---
apiVersion: decisions.sovrunn.io/v1alpha1
kind: DecisionRecord
metadata:
  name: postal-postgresql-placement-decision-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Project, name: postal-applications}
record:
  profileRef: placement-decision-v1
  subjectRef: {apiVersion: services.sovrunn.io/v1alpha1, kind: ServiceInstance, name: postal-postgresql}
  requirementsRef: postal-postgresql-requirements-001
  sovereigntyDecisionRef: postal-postgresql-sovereignty-001
  authority: Authoritative
  typedResult:
    selectedExecutionTargetRefs: [yotta-openshift-gov-primary, yotta-openshift-gov-backup]
    outcome: Selected
---
apiVersion: orchestration.sovrunn.io/v1alpha1
kind: ServiceDeploymentPlan
metadata:
  name: postal-postgresql-deployment-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Project, name: postal-applications}
record:
  serviceInstanceRef: postal-postgresql
  placementDecisionRef: postal-postgresql-placement-decision-001
  steps:
    - {id: provision-storage, action: StorageProvision}
    - {id: deploy-primary, action: PostgreSQLDeployPrimary, dependsOn: [provision-storage]}
    - {id: verify-service, action: PostgreSQLVerify, dependsOn: [deploy-primary]}
---
apiVersion: ops.sovrunn.io/v1alpha1
kind: Operation
metadata:
  name: provision-postal-postgresql-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Project, name: postal-applications}
spec:
  targetRef: {apiVersion: services.sovrunn.io/v1alpha1, kind: ServiceInstance, name: postal-postgresql}
  action: Provision
  idempotencyKey: request-7f42
status:
  phase: Running
  progressPercent: 60
---
apiVersion: execution.sovrunn.io/v1alpha1
kind: PluginExecution
metadata:
  name: deploy-postal-primary-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Project, name: postal-applications}
spec:
  operationRef: provision-postal-postgresql-001
  deploymentPlanRef: postal-postgresql-deployment-001
  planStepId: deploy-primary
  pluginRef: postgresql-management-plugin-v1
  adapterRef: openshift-adapter-v1
  executionTargetRef: yotta-openshift-gov-primary
status:
  phase: Succeeded
  terminalCode: EXECUTION_SUCCEEDED
```

## 13. Sovrunn platform lifecycle contracts

### 13.1 Contract matrix

| Concept / API kind | Purpose and ownership | Profile, scope and writers | Mutability and lifecycle | Boundary, relationships and invariants |
|---|---|---|---|---|
| `SovrunnInstallation` | Desired and observed state of one deployed Sovrunn control plane serving one provider participation | MR; Platform scope; authorized platform lifecycle administrator requests change; authoritative lifecycle repository owns accepted desired state; installation and health controllers own distinct observed fields | Identity, scope and participationRef immutable; no direct desired-release mutation; accepted Operation atomically activates desired release; backup/impact gates precede DECOMMISSION | POF/GO; keeps desired, current, last-stable and active Operation distinct; serves exactly one CloudPlatform + CloudProvider participation; Sovrunn API is request/projection, not a second authority; must not be modeled as a customer ServiceInstance; `ADH-2026-036`/`037` |
| `SovrunnRelease` | Immutable release, component manifest, schema/migration metadata and compatibility contract | VD; Platform scope; approved release publisher owns definition | Draft mutable; published version immutable; retire/supersede; referenced releases retained | POF/IE; signed artifact digests, provenance, compatibility, upgrade and rollback edges required; no runtime credentials |
| `PlatformLifecyclePolicy` | Versioned upgrade channel, maintenance, backup, approval, validation and rollback requirements | VD; Platform scope; authorized platform policy publisher owns definition | Draft mutable; published version immutable; retire/supersede; referenced versions retained | POF/GO; installation may select an allowed policy; policy cannot weaken mandatory platform controls |
| `PlatformLifecyclePlan` | Exact immutable sequence generated for one proposed installation action or release transition and pinned upon acceptance | IR; Platform scope with one SovrunnInstallation subject; platform lifecycle planner sole writer | Immutable after generation; changed release, policy, health, backup or compatibility inputs create a new candidate plan; unaccepted plans expire or are retained under audit policy | IE/GO/POF summary; references exact current/target release, policy, backup, validation, component steps and recovery strategy; accepted Operation binds exact plan and approval; SecretRefs only |
| `Operation` targeted at `SovrunnInstallation` | Sole command boundary for install, upgrade, rollback, backup, restore, rotation, DR test or decommission | Existing LRO; Platform scope; authorized requester proposes request; acceptance authority and Operation controller own accepted state/status | Idempotent request; immutable after acceptance; at most one active disruptive Operation per installation; terminal result retained; retry/cancel/superseding recovery explicit | POF/GO; acceptance atomically binds desired state, exact plan/policy/approval, expected generation, fencing token and activation event; backend execution remains asynchronous |
| `PlatformHealth` | Safe current installation release, readiness, component, backup, recovery and assurance projection | Projection of SovrunnInstallation status and retained evidence; health/projection controllers write | Recomputed from observations; history retained through evidence, operations and audit rather than mutable projection | POF; not an independently managed resource; excludes secrets and unsafe internal topology |
| `PlatformLifecycleAgent` | Executes approved lifecycle work independently of the affected Sovrunn control plane | System actor, not a customer or managed resource; distinct workload identity; external operator/GitOps/bootstrap mechanism owns execution | Agent release and credentials independently rotatable; execution attempt correlated to Operation; authority expires or is revoked | IE/GO; may execute only an approved plan with scoped authority; cannot publish releases, approve its own request, change policy or grant authorization |
| `ComponentManifest` | Exact signed component set for one release | VD; Platform scope; approved release publisher | Draft mutable; published immutable; retained while referenced | POF/IE; pins artifact digests, versions, dependencies, rollout order, configuration/schema requirements and health checks |
| `ArtifactProvenanceRecord` | Integrity and supply-chain evidence for release artifacts | EvidenceRecord profile; approved provenance publisher/collector | Append-only; expires/supersedes; retained under evidence policy | GO/POF summary; source revision, builder, attestation, signatures, SBOM, digest binding, vulnerability-policy result and validity |
| `MaintenanceWindow` | Governed time and approval boundary for lifecycle or infrastructure work | VD; applicable owner/provider scope; schedule publisher | Draft mutable; published immutable; supersede/retire | GO/OF/POF projection; time zone, recurrence or bounded interval, duration, notice, blackout, emergency and approval rules |
| `RecoveryRepository` | Registered external backup/recovery system | MR; Platform or provider-participation scope; recovery administrator/spec and qualification controller/status | Identity/type immutable; location, retention, SecretRefs and qualification controlled; decommission only after retained recovery dependencies end | POF/GO; geography, sovereignty, encryption, integrity, formats and restore evidence; no raw credentials |
| `ValidationGateDefinition` | Reusable exact validation contract for lifecycle plan phases | VD; Platform scope; approved gate publisher | Published immutable; retire; retained while plans reference it | IE/POF; evaluator, typed inputs, phase, success/failure, timeout, criticality, retry and evidence output; mandatory gates fail closed |
| `ReleaseCompatibilityContract` | Explicit directed compatibility and recovery edge between two exact releases | VD; Platform scope; release compatibility authority | Published immutable; withdraw prevents new use but retained history remains | POF/IE/GO; covers schemas, components, plugins/adapters, decision profiles, clients, agent, audit/events and backup formats; defines phase-specific recovery and irreversible checkpoint |

`POF` means platform-operator-facing and is a specialized safe operator projection. It does not expose platform lifecycle resources to managed-service customers.

### 13.2 Skeletal YAML

```yaml
apiVersion: platform.sovrunn.io/v1alpha1
kind: SovrunnRelease
metadata:
  name: sovrunn-1-0-0
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Platform, name: sovrunn}
spec:
  version: 1.0.0
  componentManifestRef: sovrunn-1-0-0-components
  artifactProvenanceRef: sovrunn-1-0-0-provenance
  schemaCompatibility:
    readableFrom: [0.9.0]
    upgradeFrom: [0.9.0]
    rollbackTo: [0.9.0]
  publicationState: Published
---
apiVersion: platform.sovrunn.io/v1alpha1
kind: PlatformLifecyclePolicy
metadata:
  name: production-stable-v1
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Platform, name: sovrunn}
spec:
  version: v1
  releaseChannel: Stable
  maintenanceWindowRef: production-platform-window
  preChangeBackupRequired: true
  restoreVerificationRequired: true
  approvalPolicyRef: production-platform-change-v1
  rollbackRequired: true
  publicationState: Published
---
apiVersion: platform.sovrunn.io/v1alpha1
kind: SovrunnInstallation
metadata:
  name: yotta-sovrunn-production
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Platform, name: sovrunn}
spec:
  participationRef: nic-cloud-yotta-production
  desiredReleaseRef: sovrunn-1-0-0
  lifecyclePolicyRef: production-stable-v1
  recoveryRepositoryRef: yotta-sovrunn-recovery
status:
  phase: Changing
  currentReleaseRef: sovrunn-0-9-0
  lastStableReleaseRef: sovrunn-0-9-0
  activeOperationRef: yotta-sovrunn-upgrade-operation-001
  lastVerifiedBackupRef: yotta-sovrunn-backup-2026-08-04
  conditions:
    - {type: Available, status: "True", reason: ComponentsReady}
---
apiVersion: platform.sovrunn.io/v1alpha1
kind: PlatformLifecyclePlan
metadata:
  name: yotta-sovrunn-upgrade-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Platform, name: sovrunn}
record:
  installationRef: yotta-sovrunn-production
  fromReleaseRef: sovrunn-0-9-0
  toReleaseRef: sovrunn-1-0-0
  releaseCompatibilityContractRef: sovrunn-0-9-0-to-1-0-0
  componentManifestRef: sovrunn-1-0-0-components
  artifactProvenanceRef: sovrunn-1-0-0-provenance
  lifecyclePolicyRef: production-stable-v1
  maintenanceWindowRef: production-platform-window
  recoveryRepositoryRef: yotta-sovrunn-recovery
  preChangeBackupRef: yotta-sovrunn-backup-2026-08-04
  validationGateRefs: [control-plane-ready, schema-compatible, restore-verified]
  steps: [verify-backup, migrate-schema, roll-components, verify-health]
  recovery:
    beforeIrreversibleBoundary: InPlaceSupported
    afterIrreversibleBoundary: RestoreRequired
    rollbackProhibitedAfterStep: migrate-schema
---
apiVersion: ops.sovrunn.io/v1alpha1
kind: Operation
metadata:
  name: yotta-sovrunn-upgrade-operation-001
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Platform, name: sovrunn}
request:
  action: Upgrade
  targetRef: {apiVersion: platform.sovrunn.io/v1alpha1, kind: SovrunnInstallation, name: yotta-sovrunn-production}
  targetReleaseRef: sovrunn-1-0-0
  expectedInstallationGeneration: 7
  idempotencyKey: yotta-prod-upgrade-1-0-0
status:
  phase: Accepted
  platformLifecyclePlanRef: yotta-sovrunn-upgrade-001
  platformLifecyclePolicyRef: production-stable-v1
  approvalRequestRef: yotta-prod-upgrade-approval-001
  fencingToken: 01J4YOTTAUPGRADE001
```

### 13.3 Independently recoverable execution

```text
External bootstrap/operator plane
        ↓
Sovrunn control plane
        ↓
Managed customer services
```

Authoritative signed release artifacts, recovery configuration, backup metadata and the minimum restoration authority path must remain recoverable outside the SovrunnInstallation they protect. The affected control plane must not be required to complete its own restore or rollback. PlatformLifecycleAgent execution uses short-lived scoped credentials, plan and Operation correlation, idempotency and complete audit.

### 13.4 Lifecycle authority and atomic activation

The canonical contract requires one logically external authoritative lifecycle repository. For the first vertical slice, use a signed GitOps repository and one signed commit as the activation boundary. A future dedicated lifecycle service may replace this physical implementation without changing the canonical resource semantics.

| Store or interface | Authoritative responsibility |
|---|---|
| External lifecycle repository | Accepted desired installation state, accepted Operation identity, pinned plan/policy/approval references, generation, fencing token and activation event |
| Signed artifact registry | Immutable Sovrunn release artifacts, manifests and provenance |
| External secret manager | Lifecycle-agent, artifact-access and recovery credentials |
| Immutable backup repository | Control-plane backups, integrity evidence and restore metadata |
| Sovrunn API | Authenticated request/proposal workflow, acceptance orchestration and operator-safe projections; never a competing desired-state authority |

The activation workflow is: authorize → validate expected generation/current release → validate target compatibility → create immutable candidate plan → approve exact plan → revalidate freshness and health → atomically accept and activate → reconcile externally → verify and publish status.

The atomic activation contains the accepted `Operation`, new `desiredReleaseRef`, `activeOperationRef`, exact plan/policy/approval references, validated generation, unique fencing token and durable activation event/outbox record. If any element cannot be committed, none is activated.

Execution after acceptance is asynchronous and must be idempotent, checkpointed and fenced. Direct PATCH of `desiredReleaseRef` is rejected. At most one disruptive Operation may be active for an installation. A stale generation, reused conflicting idempotency key or stale fencing token is rejected. Failure produces an explicit correlated rollback, restore or forward-recovery Operation; it does not silently revert desired state.

### 13.5 Directed compatibility and recovery contract

Compatibility is a directed edge; `fromReleaseRef → toReleaseRef` must be published explicitly. Unknown compatibility is prohibited, not assumed from semantic versions.

```yaml
apiVersion: platform.sovrunn.io/v1alpha1
kind: ReleaseCompatibilityContract
metadata:
  name: sovrunn-0-9-0-to-1-0-0
  scopeRef: {apiVersion: core.sovrunn.io/v1alpha1, kind: Platform, name: sovrunn}
spec:
  fromReleaseRef: sovrunn-0-9-0
  toReleaseRef: sovrunn-1-0-0
  componentManifestRef: sovrunn-1-0-0-components
  artifactProvenanceRef: sovrunn-1-0-0-provenance
  compatibility:
    resourceSchemas: CompatibleWithMigration
    pluginAndAdapterContractRange: ">=0.9 <2.0"
    decisionProfileSchemaRange: ">=1.0 <2.0"
    lifecycleAgentRange: ">=1.0 <2.0"
    backupFormatReaders: [backup-format-v3]
  recovery:
    beforeIrreversibleBoundary: InPlaceSupported
    afterIrreversibleBoundary: RestoreRequired
    rollbackProhibitedAfterStep: migrate-platform-schema-v4
  publicationState: Published
```

Executable recovery modes are `InPlaceSupported`, `RestoreRequired` and `ForwardRecoveryOnly`. `rollbackProhibitedAfterMigration` or `rollbackProhibitedAfterStep` is a boundary constraint, not an executable recovery mode; after that boundary it must resolve to RestoreRequired or ForwardRecoveryOnly.

Before acceptance, the planner verifies the exact edge, manifest, provenance, installed plugins/adapters, decision profiles, applicable client policy, lifecycle agent, backup format, qualified RecoveryRepository, pre-change backup and every mandatory ValidationGateDefinition. Any unknown mandatory dimension denies activation.

### 13.6 Canonical bootstrap (no alpha runtime migration)

Per DEC-0059 (superseding DEC-0058), there is no `CanonicalMigrationPlan`, `CanonicalMigrationRecord`, migration controller, or cutover state machine. FEATURE-0001 through FEATURE-0014 are retained repository assets and reuse input under FEATURE-0011, not live control-plane state requiring conversion. The first executable control-plane release creates canonical resources directly. Historical DecisionRecords and AuditEvents from FEATURE-0001–0014 remain in their original immutable representation as repository history and are never rewritten. Obsolete alpha write endpoints, kinds, and scope values are not carried forward as active compatibility surfaces.

The platform lifecycle boundary manages Sovrunn software and state. It does not install, upgrade or repair OpenShift, Kubernetes, OpenStack, AWS, OCI, datacenter facilities or equivalent external substrates. Those systems expose normalized maintenance and health signals through the ExecutionTarget boundary described in Section 8.3.

## 14. Future-service extension boundaries

These boundaries are part of the canonical extensibility model but do not make every concept a mandatory core or MVP resource.

| Boundary | Canonical rule | Promotion gate |
|---|---|---|
| New service family | Add a published `ServiceTypeDefinition`, compatible runtime profile and versioned plugin/adapter implementation | Must pass parameter, action, projection, binding, meter and upgrade-schema compatibility tests without core schema changes |
| Service-to-service relationship | Add a published `ServiceRelationshipDefinition` and Project-scoped `ServiceRelationship`; resolve requirements and registered DecisionRecord profiles before realization | Must prove two-sided authorization, directed semantics, bounded parameters, lifecycle/deletion behavior and adapter-independent conformance |
| Composite service | Add an optional `ServiceCompositionDefinition`; use child ServiceInstances only for independently governed components and instantiate their meaningful dependency edges as ServiceRelationships | Must define pinned relationship contracts, lifecycle, partial-failure, compensation and typed output behavior; ServiceBinding remains separately authorized access |
| Governed data | Reserve `DataAsset`, `DataFlow` and customer-safe `DataPlacement` | Promote when data is shared, independently moved, survives a service or needs a distinct sovereignty lifecycle |
| Metering and charging | Reserve `MeterDefinition`, immutable `UsageRecord`, `RatingPolicy` and immutable `ChargeRecord` | Promote when independently auditable cross-service consumption and charging are required; raw telemetry and quota remain separate |
| External supply/federation | Reserve `ServiceSupplyRelationship` and federation agreement/target contracts | Promote when an independently accountable supplier or CloudProvider realizes part of an offering |

Reserved contracts require a separate approved handoff and end-to-end reference flow before becoming API kinds. No extension may create a parallel authorization, sovereignty, placement, operation or secret-delivery path.

### 14.1 Core anti-drift boundary

Sovrunn core models service outcomes, relationships, requirements, decisions and lifecycle. It must not contain API types, fields, conditionals or workflow branches for AWS VPCs, OCI VCNs, OpenStack networks, Kubernetes NetworkPolicies or equivalent implementation-native constructs.

Implementation-native networking objects and identifiers are restricted to provider configuration, service plugins, execution adapters and protected execution records. They must not become customer or canonical core contracts. Shared topology or provider-defined logical-network ancestry never proves connectivity; an applicable registered relationship evaluator and verified realization outcome are required.

## 15. Canonical creation and mutation authority

| Actor | May author | Must not author |
|---|---|---|
| Platform architecture/operator | Platform-scoped definitions, CloudProvider registrations, global policy registries | Customer service intent or provider-native facts without authority |
| Approved platform release publisher | SovrunnRelease drafts and signed published versions within delegated scope | Installation desired state, lifecycle approval, plans, execution or health status |
| Release compatibility authority | ReleaseCompatibilityContract, ComponentManifest and approved ValidationGateDefinition drafts/publication within delegated scope | Installation desired state, lifecycle approval or execution status |
| Platform lifecycle administrator | SovrunnInstallation desired release/policy intent and authorized platform Operation requests | Published release contents, immutable lifecycle plans, approval decisions, agent execution status or underlying infrastructure lifecycle |
| Platform lifecycle planner | Immutable PlatformLifecyclePlan derived from exact release, policy, approval, backup, compatibility and health inputs | Desired installation state, approval authority or execution status |
| Platform lifecycle agent | Execution status and protected handles for an approved PlatformLifecyclePlan and correlated Operation | Release publishing, desired state, lifecycle policy, approval, authorization decisions or unplanned underlying infrastructure changes |
| Infrastructure operator or trusted adapter | InfrastructureMaintenanceNotice desired facts within delegated CloudProviderParticipation/target authority | ExecutionTarget availability/qualification status, owner restrictions or service placement decisions |
| ExecutionTarget lifecycle controller | Target availability state, maintenance epoch/fence and requalification status | Native maintenance facts, customer intent or provider agreement state |
| (retired) Canonical migration controller | None — role retired by DEC-0059; no runtime migration authority exists | Transformation outputs, CanonicalMigrationRecord, or any migration-plan/record write — all removed from active Phase 2R authorities |
| Approved service-contract publisher | ServiceTypeDefinition, ServiceRelationshipDefinition and, when applicable, ServiceCompositionDefinition drafts and published versions within delegated scope | Provider commercial terms, customer instances, decisions or execution state |
| CloudPlatform administrator | Platform catalog, enrollments, entitlements, quota envelope, approved product/governance/placement profiles and participation requests within delegated scope | Provider credentials or native operations, customer Organization ownership or customer Project contents |
| CloudProvider administrator | Participation response, topology declarations, realization mappings, maintenance notices, approved operating profiles and targets within delegated scope | CloudPlatform catalog/enrollment grants, customer Organization ownership or customer Project contents |
| Customer organization administrator | Organization hierarchy, membership/policy assignments and customer restrictions | Provider grants, target facts or wider entitlement |
| Authorized identity/security administrator | Membership intent, RoleDefinition drafts/publication requests, exact scoped RoleAssignment grant intents and authorized AccessReview campaign intents within delegated authority | Direct RoleAssignment publication/status, persona-based grants, permissions outside delegated scope, authorization outcomes or silent mutation of published role versions |
| RoleAssignment controller | Publish/revoke/expire and write status for the exact accepted immutable RoleAssignment intent | Selecting or enlarging holder, role, scope, resource, validity or delegation; originating workflow semantics |
| Project user | ServiceInstance intent, authorized ServiceRelationship intent and ServiceBinding requests within entitlement, quota and role | Status, decisions, evidence, deployment plans, raw credentials, native infrastructure objects or execution credentials |
| Policy resolver | EffectiveGovernanceContext and policy status | Source profiles or customer intent |
| Evidence collector/adapter | SovereigntyFactSet observations and EvidenceRecords within authorized subjects | Desired customer state or decision authority |
| Platform dependency resolver | Immutable PlatformDependencySnapshot from exact installation and dependency facts | Dependency desired state, evidence contents, policy or assessment outcome |
| Sovereignty/placement decision service | FEATURE-0013 DecisionRecords under registered profiles | Policy inputs, facts, evidence or execution state |
| Deployment planner | ServiceDeploymentPlan | Placement authority or backend execution |
| Operation controller | Operation status and orchestration | Customer request fields after acceptance |
| Service plugin/adapter | PluginExecution status and protected backend handles | Entitlement, governance or placement decisions |
| Projection controller | ServicePlacement customer-safe record | Canonical decisions or sensitive provider topology |
| Binding controller | ServiceBinding status, safe endpoint metadata and typed SecretRef | Secret values, provider target credentials or unrelated consumer bindings |
| Relationship controller | ServiceRelationship status, safe rationale and decision/operation references | Source definitions, customer intent, canonical decisions, bindings, secret values or implementation-native handles |

## 16. Cross-resource deletion rules

1. A CloudProvider cannot be deleted while products, enrollments, targets or retained records reference it.
2. Published definitions are retired, not mutated or immediately deleted.
3. A ServiceTypeDefinition or ServiceCompositionDefinition version cannot be removed while an offering, runtime profile, instance history, decision or retained record references it.
4. A ServiceRelationshipDefinition version cannot be removed while a ServiceRelationship, requirement set, decision, plan, operation or retained history references it.
5. Terminating an enrollment does not delete the customer Organization or its hierarchy.
6. Revoking entitlement does not silently delete running services; a governed transition policy decides deny-new, grace, migrate or terminate.
7. Deleting a Project or ServiceInstance requires an impact preview and explicit lifecycle operations.
8. A ServiceInstance cannot finalize deletion while an active required ServiceRelationship or ServiceBinding depends on it; relationships must follow declared lifecycle policy and bindings must be revoked first.
9. Datacenters, stacks and targets must be drained/decommissioned before removal and cannot be removed while active placement records depend on them.
10. Effective contexts, decisions, evidence, deployment plans, operations and placement projections obey retention/legal-hold rules and are never cascade-deleted with a mutable parent.
11. Correction of an immutable record creates a linked correction, supersession or revocation record.
12. FEATURE-0018 resources have no hard-delete surface. AccessGroup and Membership use explicit suspend/reactivate/revoke/retire behavior; RoleAssignment uses explicit revoke/expiry/review effects; published definitions retire and remain resolvable while referenced; ApprovalRequest, AccessReview and ExceptionGrant evidence follows retention/legal hold. No parent deletion or workflow outcome cascade-deletes authorization or audit history.
13. A SovrunnInstallation cannot finalize decommissioning until required backups and restore evidence are current, active platform Operations are terminal and control-plane impact has been approved.
14. A published SovrunnRelease or PlatformLifecyclePolicy version cannot be removed while an installation, plan, operation, evidence or retained audit record references it.
15. PlatformLifecyclePlan, platform-targeted Operation, backup/restore evidence and lifecycle audit history follow retention and legal-hold rules and are never cascade-deleted with a SovrunnInstallation.
16. A ReleaseCompatibilityContract, ComponentManifest, provenance record, MaintenanceWindow, RecoveryRepository or ValidationGateDefinition cannot be removed while a retained release, plan, Operation, backup or migration record references it.
17. An InfrastructureMaintenanceNotice and its epoch/fence history are retained while any decision, plan, execution, service-impact Operation or audit record references them.
18. CanonicalMigrationPlan, CanonicalMigrationRecord and legacy-to-canonical identity mappings follow retention and legal-hold rules and are never deleted merely because cutover completed.
19. PlatformDependencySnapshot and referenced evidence remain retained while any platform assessment, placement, execution, lifecycle plan or audit record depends on them.

## 17. Minimum feature specification required before implementation

Each feature implementing one of these concepts must define:

- exact API group, version, kind and collection route;
- canonical JSON Schema/OpenAPI schema and generated-type bindings;
- allowed `scopeRef` kinds and constrained reference schemas;
- field-level classification, readers, writer, mutability, retention, redaction, residency and audit requirements;
- defaulting, structural, semantic, scope/reference and authorization validation;
- status phase/condition ownership and consistency rules;
- lifecycle/state transitions and deletion/finalizer behavior;
- idempotency, optimistic concurrency and replay behavior;
- customer, CloudProvider-operator, platform-operator, internal, plugin and adapter projections;
- positive, negative, cross-scope and no-existence-disclosure fixtures;
- compatibility and migration rules;
- for a new service family, the complete ServiceTypeDefinition extension surface and proof that core schemas/workflow remain unchanged;
- for a composite service, component ownership, child-instance criteria, dependency, compensation and typed binding/output rules;
- for a ServiceRelationshipDefinition, directed source/target compatibility, parameter schema, resolver, decision profiles, bindings, lifecycle, scope and no-native-object conformance;
- for platform lifecycle, signed release provenance, complete directed compatibility dimensions, phase-specific recovery modes, irreversible boundary, typed reference contracts, upgrade/rollback/restore state machine, external-agent recovery, authority separation and proof that underlying infrastructure lifecycle remains external;
- for infrastructure maintenance, notice authority, separate qualification/availability states, expected generation, idempotency, epoch/fence, safe-point behavior and mandatory requalification;
- for alpha migration, inventory classification, deterministic mapping, write freeze, backup/restore, historical integrity, cutover, rejection of obsolete writers and no-dual-authority conformance;
- FEATURE-0013 adoption section when producing or consuming decisions, projections, audit or decision-linked operations;
- an approved ADR or handoff for any deviation from this catalog.

## 18. Recommended implementation sequence

Define complete schemas through one production-shaped PostgreSQL vertical slice rather than finalizing all fields in isolation:

1. Approve `ADH-2026-020` through `ADH-2026-041` and the coordinated ownership/scope/catalog/topology/governance/sovereignty/relationship/platform-lifecycle/migration amendments.
2. Execute the approved CanonicalMigrationPlan for common scope, ownership, catalog and topology contracts while the APIs remain alpha; never run legacy and canonical desired-state writers together.
3. Implement CloudProvider, Organization, Enrollment, ServiceRegion, ServiceTypeDefinition, ServiceRelationshipDefinition, ServiceOffering, ServicePlan, Entitlement and Project contracts.
4. Implement ServiceInstance, ServiceRelationship, ServiceBinding, RuntimeProfile and RequirementSet.
5. Implement ExecutionTarget, adapter qualification, facts and evidence.
6. Implement governance resolution and sovereignty DecisionRecord profile.
7. Implement placement DecisionRecord profile and multi-target decision rules.
8. Implement immutable deployment plan, Operation and PluginExecution.
9. Publish ServicePlacement projection and run the end-to-end Yotta/Department of Posts/PostgreSQL conformance scenario.
10. Prove reuse on a non-infrastructure or second infrastructure target, then add a second service family through ServiceTypeDefinition without modifying core schemas.
11. Before production MVP, implement SovrunnInstallation, SovrunnRelease, ComponentManifest, ReleaseCompatibilityContract, PlatformLifecyclePolicy, RecoveryRepository, ValidationGateDefinition, immutable PlatformLifecyclePlan and Operation-first change through an independently recoverable PlatformLifecycleAgent; use one external lifecycle authority and atomic activation, then prove in-place, restore-required and forward-only paths as declared with the control plane unavailable.
12. Introduce governed-data, metering or federation modules only after their promotion gates are demonstrated.

## 19. Architecture completion test

The canonical contract is ready for implementation only when:

- every kind has exactly one semantic owner and one authoritative writer per mutable field/condition;
- all current repository conflicts have approved migrations;
- no customer contract leaks provider-native or secret data;
- no authorization path depends on a persona name, UI label, job title, email address or backend IAM role;
- default roles resolve only to canonical Sovrunn actions, assignments pin exact role versions and scope-boundary negative tests deny wider authority;
- enrollment, entitlement, quota and scope resolution pass cross-boundary negative tests;
- sovereignty and placement use registered FEATURE-0013 profiles rather than parallel decision envelopes;
- immutable records preserve exact input versions and provenance;
- a Sovrunn installation can be upgraded, rolled back and restored through an external lifecycle agent using signed releases, immutable plans, scoped authority and externally recoverable artifacts when the affected control plane is unavailable;
- platform desired-state activation and Operation acceptance cannot partially commit, direct desired-release mutation is rejected, concurrent disruptive requests are fenced and the API cannot become a competing lifecycle authority;
- every permitted release edge proves all compatibility dimensions, exact typed references, a valid recovery path and irreversible-boundary behavior before activation;
- an authorized InfrastructureMaintenanceNotice advances one maintenance epoch, drains the target, blocks stale placement/execution, triggers governed service-impact actions and requires requalification without making Sovrunn the infrastructure lifecycle owner;
- a migration dry run and cutover preserve historical integrity, produce zero unresolved mappings, run one desired-state authority and reject obsolete writers after activation;
- multi-target placement and failure behavior are explicit;
- PostgreSQL and one materially different service family use distinct ServiceTypeDefinitions while sharing the same ServiceInstance/decision/operation core;
- an ExecutionTarget can model a qualified non-infrastructure realization endpoint without changing the customer contract;
- one VM-to-PostgreSQL relationship is realized on two infrastructure implementations using the same ServiceRelationship contract and no implementation-native core fields;
- ServiceBinding, relationship eligibility and connectivity/realization outcomes remain independently testable;
- deletion cannot orphan services, targets, evidence or decisions;
- skeletal YAML validates against the proposed canonical schemas;
- the PostgreSQL vertical slice can be traced from customer request to customer-safe placement without semantic gaps.

---

**Architecture note:** The YAML in this catalog is skeletal and semantic. Exact field names, enums and payload shapes become normative only after approved feature architecture, requirements, schemas and conformance fixtures. Regulatory policy content requires qualified legal, security and compliance approval.
