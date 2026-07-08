# RFC Framework

Document:
  ID: rfc-framework
  Title: RFC Framework
  Parent: docs
  Owner: SDE Architecture Council
  Layer: Governance
  Type: Framework
  Version: 2.0
  Status: Stable

Purpose:
  - Define how Sovrunn Data Engine RFCs are authored, reviewed, accepted, rejected, superseded, and maintained
  - Provide scalable RFC numbering for architecture, runtime, data plane, control plane, plugins, providers, dstoreOps, security, AI, implementation, and operations
  - Ensure RFCs preserve source-of-truth architecture and specification integrity
  - Make design decisions easy for humans and AI agents to retrieve, reason over, and validate

Definition:
  RFC means Request for Comments.

  In SDE, an RFC is a reviewed design proposal or decision record that changes, extends, constrains, or clarifies SDE architecture, specifications, implementation, security, operations, extensibility, or governance.

Core Rule:
  RFCs explain why a decision was made.

  Source-of-truth documents define what is currently true.

  When an RFC is accepted, affected architecture, specification, implementation, or operational documents must be updated.

Non-Goal:
  RFCs are not the primary source of current architecture truth.

  RFCs must not become stale shadow architecture.

---

# RFC Document Roles

RFC documents are used for:
  - Major architecture decisions
  - Public or internal contract changes
  - Runtime behavior changes
  - Control Plane behavior changes
  - Data Plane behavior changes
  - Plugin framework decisions
  - Protocol Plugin decisions
  - Engine Plugin decisions
  - Datastore Management Plane decisions
  - Datastore Operator Plugin decisions
  - Infrastructure Provider decisions
  - Foundation Provider decisions
  - Capability model decisions
  - Security and policy decisions
  - AI Control Plane decisions
  - Tenant AI Agent decisions
  - Operations and autonomous operations decisions
  - Implementation structure decisions

RFC documents are not used for:
  - Small typo fixes
  - Formatting-only changes
  - Local refactors without architectural effect
  - Temporary notes
  - Unreviewed brainstorming
  - Replacing architecture documents

---

# RFC Lifecycle

Draft:
  Meaning:
    RFC is being written and is not ready for review.

Review:
  Meaning:
    RFC is ready for architecture, implementation, security, operations, and product review.

Accepted:
  Meaning:
    RFC decision is approved.

  Requirement:
    Source-of-truth documents must be updated.

Rejected:
  Meaning:
    RFC was considered and explicitly not accepted.

Superseded:
  Meaning:
    RFC was replaced by a newer RFC.

Deprecated:
  Meaning:
    RFC was previously accepted but is no longer recommended.

Stable:
  Meaning:
    RFC is accepted and the corresponding source-of-truth documents have been updated and validated.

---

# RFC Numbering

SDE uses domain-separated RFC ranges.

Numbering is intentionally sparse to allow long-term extensibility.

## 0000-0999: Governance, Foundation, Core Architecture

0000-0099:
  Name: RFC Governance and Authoring Standards
  Scope:
    - RFC framework
    - RFC template
    - Authoring standards
    - Review process
    - Numbering policy

0100-0199:
  Name: Foundation, Glossary, Ontology, Architecture Principles
  Scope:
    - Canonical terminology
    - Ontology
    - Ownership
    - Documentation structure
    - Architecture principles

0200-0299:
  Name: Runtime Architecture
  Scope:
    - SDE Runtime
    - Protocol Runtime
    - SIR Runtime
    - Planning
    - Data Kernel
    - Engine Runtime
    - Plugin Runtime
    - Session Runtime
    - Transaction Runtime
    - Execution Context
    - Execution Plan
    - Result Model
    - Error Model

0300-0399:
  Name: SDE Data Plane Architecture
  Scope:
    - Request Flow
    - Protocol Execution
    - Planning Execution
    - Kernel Execution
    - Engine Execution
    - Result Propagation
    - Error Propagation

0400-0499:
  Name: SDE Control Plane Architecture
  Scope:
    - Core Control Plane
    - Foundation Services
    - Foundation Providers
    - Datastore Management Plane
    - Control Plane maps and boundaries

0500-0599:
  Name: Specification Framework
  Scope:
    - Versioning
    - Serialization
    - Protocol specification framework
    - Engine specification framework
    - Manifest specification framework

0600-0799:
  Name: Capability Specifications
  Scope:
    - Capability model
    - Transaction capability
    - Security capability
    - Object capability
    - Cache capability
    - Search capability
    - Indexing capability
    - Streaming capability
    - Federation capability
    - Vector capability
    - Graph capability

0800-0999:
  Name: Reserved Core Architecture
  Scope:
    - Reserved for future core architecture expansion

---

# 1000-1999: Protocol Plugins

1000-1199:
  Name: Protocol Plugin Framework
  Scope:
    - Protocol Plugin contract
    - Protocol Plugin lifecycle
    - Protocol Plugin manifest
    - Protocol conformance harness
    - Protocol Runtime integration
    - SIR handoff rules

1200-1299:
  Name: PostgreSQL Protocol Plugin
  Scope:
    - PostgreSQL wire compatibility
    - Simple Query
    - Extended Query
    - Prepared statements
    - Portals
    - Runtime parameters
    - PostgreSQL error and result mapping

1300-1399:
  Name: MySQL Protocol Plugin
  Scope:
    - MySQL protocol compatibility
    - Session behavior
    - Authentication integration
    - Request normalization
    - Result and error mapping

1400-1499:
  Name: MongoDB Protocol Plugin
  Scope:
    - MongoDB wire protocol compatibility
    - Command mapping
    - Document request normalization
    - Result and error mapping

1500-1599:
  Name: Redis Protocol Plugin
  Scope:
    - RESP compatibility
    - Command normalization
    - Streaming and pub/sub semantics
    - Result and error mapping

1600-1699:
  Name: REST, gRPC, and Native Protocol Plugins
  Scope:
    - REST protocol plugin
    - gRPC protocol plugin
    - Native SDE protocol
    - API compatibility

1700-1999:
  Name: Reserved Protocol Plugins
  Scope:
    - Future protocol plugin families

---

# 2000-2999: Engine Plugins

2000-2199:
  Name: Engine Plugin Framework
  Scope:
    - Engine Plugin contract
    - Engine Plugin lifecycle
    - Engine Plugin manifest
    - Engine conformance harness
    - Engine Runtime integration
    - Execution Fragment rules

2200-2299:
  Name: PostgreSQL Engine Plugin

2300-2399:
  Name: MySQL Engine Plugin

2400-2499:
  Name: MongoDB Engine Plugin

2500-2599:
  Name: Redis Engine Plugin

2600-2699:
  Name: Object and Table Engine Plugins
  Scope:
    - S3
    - Iceberg
    - Delta Lake
    - Parquet

2700-2799:
  Name: Search, Vector, and Graph Engine Plugins
  Scope:
    - OpenSearch
    - Milvus
    - Neo4j

2800-2899:
  Name: Cassandra and Distributed Datastore Engine Plugins

2900-2999:
  Name: Reserved Engine Plugins

---

# 3000-3999: Datastore Management Plane and Datastore Operator Plugins

3000-3199:
  Name: Datastore Management Plane Framework
  Scope:
    - DMP architecture
    - DMP API model
    - Tenant Namespace model
    - DatastoreRequest model
    - DatastoreInstance model
    - DatastoreProfile model
    - DMP workflow model

3200-3399:
  Name: Datastore Operator Plugin Framework
  Scope:
    - Datastore Operator Plugin contract
    - Datastore Operator Plugin manifest
    - Lifecycle operation contract
    - DMP integration
    - Plugin admission

3400-3499:
  Name: Datastore Operator Plugin Implementations
  Scope:
    - PostgreSQL Operator Plugin
    - MySQL Operator Plugin
    - MongoDB Operator Plugin
    - Redis Operator Plugin
    - Cassandra Operator Plugin
    - OpenSearch Operator Plugin
    - Other datastore operator plugins

3500-3699:
  Name: DMP Lifecycle Controllers
  Scope:
    - Provisioning Controller
    - Configuration Controller
    - Scaling Controller
    - Backup Controller
    - Restore Controller
    - Patch Controller
    - Upgrade Controller
    - Monitoring Controller
    - Retirement Controller

3700-3899:
  Name: DMP Policy, Credentials, Namespace, and Tenant Operations
  Scope:
    - Tenant Namespace Manager
    - Datastore Policy Controller
    - Datastore Credential Controller
    - Datastore access policies
    - Tenant-scoped lifecycle governance

3900-3999:
  Name: Reserved DMP

---

# 4000-4999: Infrastructure Providers

4000-4199:
  Name: Infrastructure Provider Framework
  Scope:
    - Infrastructure Provider contract
    - Infrastructure Provider manifest
    - Infrastructure Provider lifecycle
    - DMP integration
    - Infrastructure substrate abstraction

4200-4299:
  Name: Kubernetes Infrastructure Provider

4300-4399:
  Name: AWS Infrastructure Provider

4400-4499:
  Name: Azure Infrastructure Provider

4500-4599:
  Name: GCP Infrastructure Provider

4600-4699:
  Name: VMware, Bare Metal, Private Cloud, Sovereign Cloud, and Hybrid Cloud Providers

4700-4999:
  Name: Reserved Infrastructure Providers

---

# 5000-5999: Foundation Services and Foundation Providers

5000-5199:
  Name: Foundation Service Framework
  Scope:
    - Foundation Service contract
    - Service boundaries
    - Shared Control Plane services
    - Service compatibility

5200-5499:
  Name: Foundation Provider Framework and Implementations
  Scope:
    - Identity Provider
    - Authorization Provider
    - Policy Provider
    - Secrets Provider
    - Workflow Provider
    - Audit Provider
    - Eventing Provider
    - Observability Provider
    - Registry Provider
    - Plugin Provider

5500-5999:
  Name: Reserved Foundation

---

# 6000-6999: dstoreOps

6000-6199:
  Name: dstoreOps Framework
  Scope:
    - dstoreOps architecture
    - dstoreOps workflow model
    - Operational capability model
    - Runbook model
    - Lifecycle integration

6200-6299:
  Name: dstoreOps Policy Model

6300-6399:
  Name: Provisioning Workflows

6400-6499:
  Name: Scaling Workflows

6500-6599:
  Name: Backup and Restore Workflows

6600-6699:
  Name: Patch and Upgrade Workflows

6700-6799:
  Name: Monitoring, Health, and Retirement Workflows

6800-6999:
  Name: Reserved dstoreOps

---

# 7000-7999: Security, Policy, Governance, and Compliance

7000-7199:
  Name: Security Architecture
  Scope:
    - Trust model
    - Tenant isolation
    - Authentication
    - Authorization
    - Secrets
    - Encryption
    - Zero trust

7200-7399:
  Name: Policy and Governance
  Scope:
    - Policy model
    - Approval model
    - Admission model
    - Guardrails
    - Compliance constraints

7400-7599:
  Name: Audit and Compliance
  Scope:
    - Audit model
    - Evidence model
    - Compliance export
    - Regulated environment support

7600-7999:
  Name: Reserved Security and Governance

---

# 8000-8999: Implementation

8000-8199:
  Name: Repository and Module Structure

8200-8399:
  Name: Coding Standards and SDKs

8400-8599:
  Name: Build, Test, and Release

8600-8799:
  Name: Deployment and Environment Structure

8800-8999:
  Name: Reserved Implementation

---

# 9000-9999: AI Control Plane and AI-Native Operations

9000-9199:
  Name: AI Control Plane
  Scope:
    - AI Control Plane architecture
    - AI governance boundaries
    - AI integration with Control Plane services
    - AI observation, recommendation, validation, and workflow initiation

9200-9399:
  Name: Tenant AI Agent
  Scope:
    - Tenant AI Agent Interface
    - Tenant intent handling
    - Tenant context resolution
    - Tenant datastore planning
    - Tenant configuration assistance
    - Tenant workflow request generation
    - Tenant guardrails

9400-9499:
  Name: AI Observation and Stability
  Scope:
    - Runtime observation
    - Issue detection
    - Root-cause hypothesis
    - Stability analysis
    - Runtime health summarization

9500-9599:
  Name: AI Safe Autotuning
  Scope:
    - Non-destructive tuning
    - Policy-guarded tuning
    - Simulation and rollback
    - Tuning safety limits

9600-9699:
  Name: AI Remediation and Runbooks
  Scope:
    - Remediation recommendation
    - Runbook generation
    - Workflow-assisted remediation
    - Incident explanation

9700-9799:
  Name: AI Evaluation and Safety
  Scope:
    - AI evaluation harness
    - AI action policy
    - AI safety governance
    - Prompt/model safety
    - Validation and simulation

9800-9899:
  Name: AI-Assisted Artifact Generation
  Scope:
    - Configuration drafts
    - Workflow request drafts
    - Tenant integration artifacts
    - Documentation drafts
    - Future developer-assist plugin draft generation

9900-9999:
  Name: Reserved AI Future Use

---

# RFC Required Sections

Every RFC must include:
  - Title
  - Metadata
  - Summary
  - Problem Statement
  - Goals
  - Non-Goals
  - Context
  - Proposal
  - Architecture Impact
  - Specification Impact
  - Security Impact
  - Operational Impact
  - Compatibility Impact
  - Source-of-Truth Updates
  - Alternatives Considered
  - Risks
  - Rollout Plan
  - Open Questions
  - Decision

Plugin RFCs must additionally include:
  - Plugin type
  - Runtime or Control Plane boundary
  - Manifest impact
  - Registry impact
  - Validation and conformance impact
  - Compatibility matrix
  - Failure behavior

AI RFCs must additionally include:
  - AI action classes
  - Tenant isolation impact
  - Policy requirements
  - Approval requirements
  - Audit requirements
  - Rollback requirements
  - Safety boundaries

DMP and dstoreOps RFCs must additionally include:
  - Tenant namespace impact
  - Datastore lifecycle impact
  - Workflow impact
  - Datastore Operator Plugin impact
  - Infrastructure Provider impact
  - Backup/restore impact where applicable

---

# Source-of-Truth Update Rule

Accepted RFCs must update source-of-truth documents.

Examples:
  - Architecture RFC updates docs/architecture
  - Specification RFC updates docs/specifications
  - Runtime RFC updates docs/architecture/runtime and implementation docs
  - Data Plane RFC updates docs/architecture/data-plane
  - Control Plane RFC updates docs/architecture/control-plane
  - Protocol Plugin RFC updates docs/specifications/protocol and plugin docs
  - Engine Plugin RFC updates docs/specifications/engine and plugin docs
  - DMP RFC updates docs/architecture/control-plane/datastore-management-plane
  - AI RFC updates docs/architecture/control-plane/ai-control-plane.md when applicable

Rule:
  An RFC may be Accepted before source-of-truth docs are updated.

  An RFC may become Stable only after source-of-truth docs are updated and validated.

---

# Review Requirements

Architecture Review:
  Required for all RFCs.

Security Review:
  Required when RFC affects:
    - identity
    - authorization
    - policy
    - secrets
    - audit
    - tenant isolation
    - AI action authority
    - plugin admission
    - runtime execution

Operations Review:
  Required when RFC affects:
    - deployment
    - monitoring
    - backup
    - restore
    - upgrade
    - scaling
    - incident response
    - dstoreOps
    - AI remediation or tuning

Specification Review:
  Required when RFC affects:
    - protocol contracts
    - engine contracts
    - capability contracts
    - manifests
    - versioning
    - serialization

Implementation Review:
  Required when RFC affects:
    - repository structure
    - modules
    - SDKs
    - build
    - tests
    - runtime code
    - plugin code

Tenant Experience Review:
  Required when RFC affects:
    - tenant APIs
    - Tenant AI Agent
    - tenant onboarding
    - tenant workflows
    - tenant namespaces
    - tenant-visible behavior

---

# AI and Generated Content Rule

AI may assist in drafting RFCs.

AI-generated RFC content must be reviewed.

AI must not be treated as final authority for:
  - architecture acceptance
  - security acceptance
  - production runtime behavior
  - plugin admission
  - tenant isolation policy
  - destructive operation policy

Generated artifacts are untrusted until validated.

---

# File Naming

RFC files must use this format:

  rfc-NNNN-short-title.md

Examples:
  - rfc-0000-authoring-standard.md
  - rfc-0200-runtime-core-model.md
  - rfc-0300-data-plane-request-lifecycle.md
  - rfc-1000-protocol-plugin-framework.md
  - rfc-2000-engine-plugin-framework.md
  - rfc-3000-datastore-management-plane-framework.md
  - rfc-6000-dstoreops-framework.md
  - rfc-9000-ai-control-plane-reservation.md

---

# Relationship to ADRs

SDE uses RFCs as the primary decision record.

ADRs may be introduced later for smaller local implementation decisions, but architecture-level decisions should use RFCs.
