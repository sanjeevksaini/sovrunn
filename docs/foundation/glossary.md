# Glossary

Document:
  ID: glossary
  Title: Glossary
  Parent: foundation
  Owner: SDE Architecture Council
  Layer: Foundation
  Type: Reference
  Version: 1.1
  Status: Stable

Purpose:
  - Define canonical Sovrunn Data Engine terminology
  - Prevent ambiguous use of architecture, runtime, control plane, data plane, plugin, datastore, and AI terms
  - Provide stable vocabulary for architecture docs, specifications, RFCs, implementation docs, and AI-assisted reasoning

Rule:
  Use the canonical term exactly as defined here.

  Deprecated or ambiguous terms should be replaced with the canonical term listed in Deprecated Terms.

---

# Core Product Terms

Sovrunn:
  Definition:
    Parent entity and product family for sovereign, AI-native datastore platform capabilities.

  Notes:
    - Sovrunn is the broader product/company/product-family name.
    - Sovrunn Data Engine is a platform under Sovrunn.

Sovrunn Data Engine:
  Abbreviation:
    SDE

  Definition:
    AI-native sovereign datastore platform that provides governed data access, datastore management, runtime execution, specifications, plugin integration, and tenant-facing control capabilities.

  Notes:
    - Use Sovrunn Data Engine on first mention.
    - Use SDE after abbreviation is introduced.
    - Do not use Database Platform as a synonym.

SDE:
  Definition:
    Canonical abbreviation for Sovrunn Data Engine.

  Notes:
    - Use SDE for internal architecture documents after first expansion.
    - SDE includes Control Plane, Data Plane, Runtime, specifications, and extension surfaces.

SDE Platform:
  Definition:
    Broad term for Sovrunn Data Engine including Control Plane, Data Plane, Runtime, specifications, plugin frameworks, and operational capabilities.

SDE Tenant:
  Definition:
    Customer, organization, team, or logical consumer boundary that owns or operates tenant-scoped SDE resources.

Tenant Namespace:
  Definition:
    Tenant-scoped management boundary used by SDE Control Plane to isolate tenant-specific resources, configuration, credentials, policies, telemetry, workflows, and operational actions.

  Notes:
    - Tenant Namespace may map to Kubernetes namespace, cloud account, project, subscription, VPC, database namespace, secret namespace, monitoring namespace, or SDE logical namespace.
    - Do not assume Tenant Namespace means Kubernetes only.

---

# Plane Terms

SDE Control Plane:
  Definition:
    Control layer of Sovrunn Data Engine responsible for authoritative platform state, governance, registry metadata, tenant management, policy, workflows, lifecycle management, and pluggable management capabilities.

  Owns:
    - Authoritative platform state
    - Registries
    - Policy governance
    - Tenant management
    - Workflow coordination
    - Datastore Management Plane
    - Optional AI Control Plane extensions

  Must Not:
    - Execute tenant data requests directly
    - Replace SDE Data Plane request execution

SDE Data Plane:
  Definition:
    Runtime request execution plane of Sovrunn Data Engine that processes client requests using approved Control Plane state and delegates downstream execution through Engine Plugins.

  Owns:
    - Protocol execution
    - SIR handoff
    - Planning execution
    - Kernel execution
    - Engine execution
    - Result propagation
    - Error propagation

  Must Not:
    - Mutate SDE Control Plane authoritative state
    - Manage datastore lifecycle
    - Invoke Datastore Operator Plugins
    - Invoke Infrastructure Providers directly

SDE Runtime:
  Definition:
    Set of reusable runtime components used by SDE Data Plane to process requests, manage context, coordinate execution, and normalize results and errors.

  Includes:
    - Protocol Runtime
    - SIR Runtime
    - Planning
    - Data Kernel
    - Engine Runtime
    - Plugin Runtime
    - Session Runtime
    - Transaction Runtime
    - Execution Plan
    - Execution Context
    - Result Model
    - Error Model
    - Capability Registry

Datastore Data Plane:
  Definition:
    Native request execution plane of a Downstream Datastore.

  Notes:
    - Owned by the Downstream Datastore, not SDE.
    - SDE reaches Datastore Data Plane only through approved Engine Plugins.
    - Do not call this Engine Data Plane.

Datastore Management Plane:
  Abbreviation:
    DMP

  Definition:
    Pluggable management plane inside SDE Control Plane responsible for tenant-scoped Downstream Datastore lifecycle and operations.

  Notes:
    - DMP is the first pluggable management plane hosted through Management Plane Framework.
    - DMP manages datastores.
    - DMP is not the same as DMP Controller Runtime.
    - DMP does not execute tenant data requests.
    - DMP uses Datastore Operator Plugins and Infrastructure Providers.


Management Plane Framework:
  Definition:
    Shared framework inside SDE Control Plane that allows management domains to be added as governed, pluggable planes.

  Responsibilities:
    - Management plane registration
    - Management plane manifest validation
    - Management plane lifecycle governance
    - Controller runtime integration
    - Policy, workflow, audit, and observability integration
    - Management plane conformance

Pluggable Management Plane:
  Definition:
    Domain-specific control-plane capability hosted through the Management Plane Framework.

  Notes:
    - Datastore Management Plane is the first pluggable management plane.
    - A pluggable management plane must use SDE Control Plane Foundation Services.
    - A pluggable management plane must not execute tenant data-plane requests.

Management Plane Manifest:
  Definition:
    Versioned declaration describing a pluggable management plane, its domain, resources, APIs, dependencies, controller runtime requirements, compatibility, policy requirements, and conformance suite.

Management Plane Registry:
  Definition:
    Authoritative SDE Control Plane registry for pluggable management plane metadata, lifecycle state, compatibility, admission, and conformance status.

Management Plane Admission:
  Definition:
    Policy-governed process for validating, approving, rejecting, or enabling a pluggable management plane.

Management Plane Conformance:
  Definition:
    Validation suite and status proving that a pluggable management plane respects SDE Control Plane contracts, security boundaries, lifecycle rules, and non-data-plane execution constraints.

DMP Controller Runtime:
  Definition:
    Executable runtime that hosts and reconciles Datastore Management Plane resources, workflows, controllers, and plugin interactions.

  Notes:
    - DMP Controller Runtime is not the entire Datastore Management Plane.
    - `sde-dmp-controller` is the binary name for DMP Controller Runtime.

sde-dmp-controller:
  Definition:
    Binary executable for the DMP Controller Runtime.

  Notes:
    - It hosts and reconciles DMP.
    - It is not the whole DMP.

---

# Datastore Terms

Downstream Datastore:
  Definition:
    External or managed datastore integrated with SDE for storage, query, search, object, cache, graph, vector, or table functionality.

  Examples:
    - PostgreSQL
    - MySQL
    - MongoDB
    - Redis
    - Cassandra
    - OpenSearch
    - Neo4j
    - Milvus
    - S3
    - Iceberg
    - Delta Lake
    - Parquet

  Notes:
    - Preferred term over Downstream Database.
    - Use Downstream Datastore for generality.

Downstream Engine:
  Definition:
    Native execution engine of a Downstream Datastore when discussing execution semantics or Engine Plugin integration.

  Notes:
    - Use only when execution integration is central.
    - Prefer Downstream Datastore for general architecture.

Datastore Instance:
  Definition:
    Tenant-scoped or platform-managed concrete instance of a Downstream Datastore.

Datastore Profile:
  Definition:
    Approved template describing datastore class, size, topology, configuration, backup, scaling, monitoring, security, and operational defaults.

Datastore Request:
  Definition:
    Tenant or system request to create, configure, modify, observe, operate, or retire a Datastore Instance through the Datastore Management Plane.

Datastore Policy:
  Definition:
    Policy governing datastore configuration, access, backup, scaling, region, data residency, maintenance, security, and lifecycle behavior.

Datastore Workflow:
  Definition:
    Governed workflow executed by DMP to perform datastore lifecycle or operational activity.

Datastore Operation:
  Definition:
    Concrete lifecycle or operational action performed on a Datastore Instance.

dstoreOps:
  Definition:
    Sovrunn capability for managed Downstream Datastore operations powered by the Datastore Management Plane.

  Includes:
    - Provisioning
    - Configuration
    - Scaling
    - Backup
    - Restore
    - Patch
    - Upgrade
    - Monitoring
    - Retirement
    - Operational workflows

  Notes:
    - dstoreOps is an SDE capability.
    - dstoreOps is not SDE Data Plane request execution.
    - dstoreOps must operate through DMP governance.

---

# Plugin and Provider Terms

Plugin:
  Definition:
    Governed extension component that integrates an external protocol, engine, datastore lifecycle API, infrastructure substrate, or foundation implementation with SDE.

  Notes:
    - Always specify plugin type.

Protocol Plugin:
  Definition:
    SDE Data Plane plugin that integrates one client protocol into SDE.

  Responsibilities:
    - Parse protocol input
    - Preserve protocol-visible semantics
    - Produce protocol-normalized intent
    - Map Result Model to protocol-compatible response
    - Map Error Model to protocol-compatible error response

  Must Not:
    - Produce Execution Plan
    - Invoke Engine Runtime
    - Invoke Engine Plugin
    - Access Downstream Datastore
    - Manage datastore lifecycle
    - Mutate SDE Control Plane authoritative state

Engine Plugin:
  Definition:
    SDE Data Plane plugin that integrates SDE execution with a Downstream Engine.

  Responsibilities:
    - Receive execution fragments from Engine Runtime
    - Translate fragments into downstream-native operations
    - Invoke Downstream Datastore through approved interfaces
    - Map native result to Result Model
    - Map native error to Error Model

  Must Not:
    - Parse client protocol
    - Produce Execution Plan
    - Manage datastore lifecycle
    - Replace Datastore Operator Plugin
    - Invoke Infrastructure Provider
    - Mutate SDE Control Plane authoritative state

Datastore Operator Plugin:
  Definition:
    Datastore Management Plane plugin that integrates datastore lifecycle and operational actions with a specific Downstream Datastore.

  Responsibilities:
    - Provision datastore
    - Configure datastore
    - Scale datastore
    - Backup datastore
    - Restore datastore
    - Patch datastore
    - Upgrade datastore
    - Monitor datastore
    - Retire datastore

  Must Not:
    - Execute tenant data-plane requests
    - Replace Engine Plugin
    - Bypass DMP workflows, policy, or audit

Infrastructure Provider:
  Definition:
    Datastore Management Plane integration for provisioning or managing infrastructure substrate.

  Examples:
    - AWS
    - Azure
    - GCP
    - Kubernetes
    - VMware
    - Bare metal
    - Private cloud
    - Sovereign cloud
    - Hybrid cloud

  Notes:
    - Infrastructure Provider is not a Foundation Provider.
    - Infrastructure Provider supports DMP and datastore lifecycle.

Foundation Provider:
  Definition:
    Provider implementation for a Foundation Service.

  Examples:
    - Keycloak for Identity Service
    - OPA for Policy Service
    - Vault for Secrets Service
    - Temporal for Workflow Service
    - Kafka or Redpanda for Eventing Service
    - OpenTelemetry for Observability Service

  Notes:
    - Foundation Provider implements Foundation Service contracts.
    - Foundation Provider is not an Infrastructure Provider.

dstoreOps Extension:
  Definition:
    Extension package that adds dstoreOps operational workflows, detectors, remediations, or runbooks.

  Notes:
    - Future capability.
    - Must operate through Datastore Management Plane governance.

---

# Foundation Service Terms

Foundation Service:
  Definition:
    Stable SDE Control Plane service contract that provides shared platform capability used by Core Control Plane, Datastore Management Plane, AI Control Plane, and other components.

Identity Service:
  Definition:
    Foundation Service responsible for identity, subjects, principals, authentication context, and identity federation integration.

Authorization Service:
  Definition:
    Foundation Service responsible for authorization decisions and access control enforcement.

Tenant Management Service:
  Definition:
    Foundation Service responsible for tenant identity, tenant lifecycle, tenant metadata, tenant namespace association, and tenant-scoped boundaries.

Configuration Service:
  Definition:
    Foundation Service responsible for validated configuration storage and distribution.

Policy Service:
  Definition:
    Foundation Service responsible for policy evaluation, policy decisions, and policy-governed constraints.

Secrets Service:
  Definition:
    Foundation Service responsible for secrets, credential references, secret access control, and secret lifecycle integration.

Audit Service:
  Definition:
    Foundation Service responsible for durable audit records of security, control, workflow, and operational actions.

Workflow Service:
  Definition:
    Foundation Service responsible for orchestrating governed workflows.

Eventing Service:
  Definition:
    Foundation Service responsible for event publication, subscription, and event-driven integration.

Observability Service:
  Definition:
    Foundation Service responsible for logs, metrics, traces, telemetry, health signals, and observability integration.

Registry Framework:
  Definition:
    Foundation Service capability for creating governed registries and registry metadata contracts.

Plugin Framework:
  Definition:
    Foundation Service capability for plugin lifecycle, manifests, validation, admission, and runtime integration boundaries.

---

# Core Control Plane Terms

Core Control Plane:
  Definition:
    SDE Control Plane domain responsible for platform-level runtime, plugin, engine, capability, and deployment governance.

Runtime Registry:
  Definition:
    Registry containing approved runtime components, runtime versions, runtime compatibility metadata, and runtime lifecycle state.

Plugin Registry:
  Definition:
    Registry containing approved plugin metadata, plugin manifests, plugin lifecycle state, compatibility metadata, and admission decisions.

Engine Registry:
  Definition:
    Registry containing approved Downstream Engine metadata, supported versions, compatibility rules, and lifecycle state.

Capability Governance:
  Definition:
    Control Plane capability responsible for approving, versioning, validating, and governing SDE capabilities.

Deployment Governance:
  Definition:
    Control Plane capability responsible for governing deployment eligibility, environment promotion, runtime activation, and rollout constraints.

---

# Runtime Terms

Protocol Runtime:
  Definition:
    Runtime component that manages protocol request lifecycle, Protocol Plugin resolution, request context establishment, and protocol response return.

SIR Runtime:
  Definition:
    Runtime component that receives protocol-normalized intent and creates or validates Semantic Intermediate Representation.

Planning:
  Definition:
    Runtime component that converts validated SIR into immutable Execution Plan using approved runtime state, capability metadata, policy context, and engine metadata.

Data Kernel:
  Definition:
    Runtime orchestration component that executes immutable Execution Plan using immutable Execution Context and delegates downstream operations through Engine Runtime.

Engine Runtime:
  Definition:
    Runtime component that resolves approved Engine Plugin and delegates execution fragments to it.

Plugin Runtime:
  Definition:
    Runtime component that supports plugin loading, compatibility checks, invocation boundaries, and plugin execution isolation.

Session Runtime:
  Definition:
    Runtime component that manages session context and session lifecycle.

Transaction Runtime:
  Definition:
    Runtime component that manages transaction context, transaction references, transaction state, and transaction boundary semantics.

Execution Plan:
  Definition:
    Immutable planning output that defines operation graph, execution strategy, dependencies, capabilities, engine bindings, and constraints.

Execution Context:
  Definition:
    Immutable request-scoped runtime context containing request identity, trace identity, tenant context, security context, session reference, transaction reference, and execution metadata.

Result Model:
  Definition:
    Canonical SDE runtime representation of successful, partial, streaming, cursor, or continuation result.

Error Model:
  Definition:
    Canonical SDE runtime representation of failure, including code, category, message, severity, source, state, retry classification, timestamp, trace identifier, and safe details.

Capability Registry:
  Definition:
    Runtime-facing approved capability metadata used by Planning to validate required capabilities.

---

# Semantic and Specification Terms

Semantic Intermediate Representation:
  Abbreviation:
    SIR

  Definition:
    Canonical semantic representation used by SDE to express client request intent independently from client protocol and downstream engine.

SIR:
  Definition:
    Canonical abbreviation for Semantic Intermediate Representation.

Protocol-Normalized Intent:
  Definition:
    Protocol-layer intermediate output produced by Protocol Plugin and consumed by SIR Runtime.

  Notes:
    - It is not SIR.
    - It is not Execution Plan.
    - It is not downstream-native operation.

Execution Fragment:
  Definition:
    Engine-targeted execution unit derived from Execution Plan and sent by Engine Runtime to Engine Plugin.

Capability:
  Definition:
    Explicitly named behavior that may be required by SIR, offered by Engine Plugin, exposed by Protocol Plugin, governed by SDE Control Plane, and validated by Planning.

Capability Manifest:
  Definition:
    Manifest declaring capabilities supported by a plugin or engine integration.

Protocol Plugin Manifest:
  Definition:
    Manifest declaring Protocol Plugin identity, supported protocol versions, supported request forms, session behavior, transaction behavior, result mapping, error mapping, and capabilities.

Engine Plugin Manifest:
  Definition:
    Manifest declaring Engine Plugin identity, supported engine versions, execution forms, capabilities, translation behavior, result mapping, error mapping, and credential requirements.

Specification Layer:
  Definition:
    Set of SDE documents defining implementable contracts for versioning, serialization, capabilities, protocols, engines, manifests, and compatibility.

Versioning Specification:
  Definition:
    Specification defining version format, compatibility rules, lifecycle states, deprecation, breaking changes, and contract versioning behavior.

Serialization Specification:
  Definition:
    Specification defining canonical encoding, field rules, validation, null handling, unknown field handling, timestamp rules, identifier rules, and redaction.

Protocol Specification:
  Definition:
    Specification defining Protocol Plugin contracts and protocol integration behavior.

Engine Specification:
  Definition:
    Specification defining Engine Plugin contracts and Downstream Engine integration behavior.

Capability Specification:
  Definition:
    Specification defining SDE capability identity, support levels, manifests, validation, downgrade, emulation, and governance.

---

# Data Plane Flow Terms

Request Flow:
  Definition:
    End-to-end SDE Data Plane request lifecycle from client protocol request to protocol-compatible response.

Protocol Execution:
  Definition:
    Data Plane flow for protocol request entry, Protocol Plugin parsing, protocol-normalized intent creation, and protocol response or error mapping.

Planning Execution:
  Definition:
    Data Plane flow for converting validated SIR into immutable Execution Plan.

Kernel Execution:
  Definition:
    Data Plane flow for Data Kernel orchestration of Execution Plan.

Engine Execution:
  Definition:
    Data Plane flow for Engine Runtime and Engine Plugin downstream execution.

Result Propagation:
  Definition:
    Data Plane flow that carries successful or partial output from Downstream Datastore through Result Model to protocol-compatible client response.

Error Propagation:
  Definition:
    Data Plane flow that carries failure from runtime or Downstream Datastore through Error Model to protocol-compatible client error response.

---

# Datastore Management Plane Terms

Tenant Namespace Manager:
  Definition:
    DMP component responsible for resolving and managing tenant-specific namespaces.

Datastore Request Controller:
  Definition:
    DMP controller responsible for validating and reconciling DatastoreRequest resources.

Datastore Lifecycle Controller:
  Definition:
    DMP controller responsible for datastore lifecycle state and workflow coordination.

Datastore Configuration Controller:
  Definition:
    DMP controller responsible for datastore configuration changes.

Datastore Policy Controller:
  Definition:
    DMP controller responsible for datastore policy application and validation.

Datastore Credential Controller:
  Definition:
    DMP controller responsible for credential reference lifecycle and integration with Secrets Service.

Datastore Monitoring Controller:
  Definition:
    DMP controller responsible for datastore monitoring integration and health signal collection.

Backup Controller:
  Definition:
    DMP controller responsible for backup workflows.

Restore Controller:
  Definition:
    DMP controller responsible for restore workflows.

Scaling Controller:
  Definition:
    DMP controller responsible for datastore scaling workflows.

Patch Controller:
  Definition:
    DMP controller responsible for datastore patch workflows.

Upgrade Controller:
  Definition:
    DMP controller responsible for datastore upgrade workflows.

Retirement Controller:
  Definition:
    DMP controller responsible for datastore retirement workflows.

---

# AI Terms

AI Control Plane:
  Definition:
    Optional pluggable SDE Control Plane capability that exposes AI-assisted observation, recommendation, validation, tenant assistance, workflow initiation, and platform stabilization through governed Control Plane APIs.

  Notes:
    - AI Control Plane is part of SDE Control Plane.
    - AI Control Plane is not part of SDE Data Plane.
    - AI Control Plane must not bypass authorization, policy, workflow, audit, Datastore Management Plane, Datastore Operator Plugins, Infrastructure Providers, registry governance, or capability governance.
    - Detailed AI internals are intentionally deferred.

Tenant AI Agent:
  Definition:
    Tenant-facing AI interface that helps customers configure, integrate, monitor, and operate tenant-scoped SDE resources through approved Control Plane workflows.

  Notes:
    - Tenant AI Agent is tenant-scoped.
    - Tenant AI Agent may generate recommendations, configuration drafts, workflow requests, explanations, runbooks, and integration guidance.
    - Tenant AI Agent must not directly manage Downstream Datastores, infrastructure, plugins, secrets, or tenant data.

AI Recommendation:
  Definition:
    Non-authoritative AI-generated guidance that may explain, suggest, or propose an action but does not become an approved Control Plane change until validated and authorized.

AI Workflow Request:
  Definition:
    AI-generated or AI-assisted request submitted to Workflow Service for policy-governed execution.

AI Action Policy:
  Definition:
    Policy that classifies and constrains actions an AI capability may observe, recommend, validate, initiate, or execute.

AI-Generated Artifact:
  Definition:
    Artifact produced by AI, such as configuration draft, workflow request, runbook, remediation plan, integration guide, test draft, documentation draft, or code draft.

  Notes:
    - AI-generated artifacts are untrusted until validated.
    - AI-generated artifacts must not be treated as approved runtime artifacts without policy and validation.

AI Action Class:
  Definition:
    Risk classification for AI-assisted action.

  Classes:
    - Class 0: Read-only observation and explanation
    - Class 1: Generate draft recommendations or configuration artifacts
    - Class 2: Validate or simulate proposed action
    - Class 3: Apply reversible, non-destructive, policy-approved tuning
    - Class 4: Controlled operational change requiring workflow approval
    - Class 5: Destructive or irreversible action requiring explicit human approval and safety workflow

AI-Assisted Artifact Generation:
  Definition:
    Future AI capability for generating drafts of configuration, workflows, tests, documentation, runbooks, or code.

  Notes:
    - Not in initial AI Control Plane scope.
    - Generated artifacts are untrusted until validated.

AI Safe Autotuning:
  Definition:
    Future AI capability for recommending or applying reversible, non-destructive, policy-approved tuning.

  Notes:
    - Must include simulation, policy validation, rollback, audit, and safety limits.

AI Remediation:
  Definition:
    Future AI capability for recommending or initiating approved remediation workflows.

  Notes:
    - Destructive or irreversible remediation requires explicit human approval and safety workflow.

---

# RFC Terms

RFC:
  Definition:
    Reviewed design proposal or decision record that changes, extends, constrains, or clarifies SDE architecture, specifications, implementation, security, operations, or governance.

RFC Framework:
  Definition:
    Governance system for RFC authoring, numbering, review, status, acceptance, and source-of-truth updates.

RFC Status:
  Definition:
    Lifecycle state of an RFC.

  Values:
    - Draft
    - Review
    - Accepted
    - Rejected
    - Superseded
    - Deprecated
    - Stable

Source of Truth:
  Definition:
    Current authoritative architecture, specification, implementation, or operational documentation.

  Notes:
    - RFCs explain why.
    - Source-of-truth docs define what is currently true.

---

# Documentation Structure Terms

MAP:
  Definition:
    Documentation file that provides navigation, ownership, relationships, dependencies, and reading order.

ARCHITECTURE:
  Definition:
    Documentation file that defines design, structure, responsibilities, boundaries, invariants, and major flows.

FLOW:
  Definition:
    Documentation file that defines an end-to-end sequence of behavior.

SUBFLOW:
  Definition:
    Focused stage inside a larger flow.

CONTRACT:
  Definition:
    Documentation file that defines responsibilities, inputs, outputs, rules, constraints, compatibility, and failure behavior for a component, service, plugin, provider, or specification.

ADR:
  Definition:
    Architecture Decision Record.

  Notes:
    - RFCs are preferred for Sovrunn unless ADRs are explicitly introduced.

---

# Deprecated Terms

Control Plane:
  Replace With:
    SDE Control Plane

Data Plane:
  Replace With:
    SDE Data Plane

Database Platform:
  Replace With:
    Sovrunn Data Engine

Database Engine:
  Replace With:
    Downstream Engine

Downstream Database:
  Replace With:
    Downstream Datastore

Engine Data Plane:
  Replace With:
    Datastore Data Plane

Management Plane Plugin:
  Replace With:
    Datastore Operator Plugin, Infrastructure Provider, or Foundation Provider

Provider:
  Replace With:
    Foundation Provider or Infrastructure Provider

AI Plugin Generation:
  Replace With:
    AI-Assisted Artifact Generation, if discussing future draft generation

Non-Destructible Auto-Tuning:
  Replace With:
    Non-Destructive Auto-Tuning or AI Safe Autotuning

---

# Naming Rules

Rule 1:
  Use SDE Control Plane and SDE Data Plane when referring to Sovrunn Data Engine planes.

Rule 2:
  Use Datastore Data Plane for native downstream request execution plane.

Rule 3:
  Use Downstream Datastore unless specifically discussing native execution engine semantics.

Rule 4:
  Use Engine Plugin only for SDE Data Plane execution integration.

Rule 5:
  Use Datastore Operator Plugin only for DMP lifecycle integration.

Rule 6:
  Use Infrastructure Provider only for infrastructure substrate integration.

Rule 7:
  Use Foundation Provider only for implementation of Foundation Services.

Rule 8:
  Use AI Control Plane only as optional pluggable Control Plane extension until detailed AI scope is accepted through RFC.

Rule 9:
  Use Tenant AI Agent for tenant-facing AI interface.

Rule 10:
  Do not introduce new architecture terms without adding them to this glossary.
