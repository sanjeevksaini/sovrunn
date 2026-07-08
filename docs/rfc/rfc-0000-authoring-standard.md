# RFC-0000: Authoring Standard

Document:
  ID: rfc-0000
  Title: Authoring Standard
  Parent: rfc-index
  Owner: SDE Architecture Council
  Layer: Governance
  Type: RFC
  Version: 2.0
  Status: Stable

RFC Number:
  0000

Status:
  Stable

Created:
  2026-07-08

Updated:
  2026-07-08

Authors:
  - SDE Architecture Council

Reviewers:
  - Architecture: Required
  - Security: Required when applicable
  - Operations: Required when applicable
  - Specification: Required when applicable
  - Implementation: Required when applicable
  - Tenant Experience: Required when applicable

Affected Areas:
  - RFC Framework
  - Documentation Governance
  - Source-of-Truth Management
  - AI-Assisted Authoring

Source-of-Truth Documents:
  - docs/rfc/README.md
  - docs/rfc/index.md
  - docs/rfc/rfc-template.md
  - docs/rfc/rfc-0000-authoring-standard.md

Supersedes:
  - None

Superseded By:
  - None

---

# Summary

This RFC defines the authoring standard for all Sovrunn Data Engine RFCs.

It defines purpose, lifecycle, required sections, numbering, review expectations, source-of-truth update rules, plugin-specific requirements, DMP-specific requirements, dstoreOps-specific requirements, and AI-specific requirements.

---

# Problem Statement

SDE is a scalable and extensible platform.

It includes:
  - SDE Control Plane
  - SDE Data Plane
  - SDE Runtime
  - Protocol Plugins
  - Engine Plugins
  - Datastore Management Plane
  - Datastore Operator Plugins
  - Infrastructure Providers
  - Foundation Services
  - Foundation Providers
  - dstoreOps
  - AI Control Plane reservation
  - Tenant AI Agent reservation
  - Security and policy governance
  - Implementation and operations

Without a formal RFC standard, architecture decisions may become fragmented, inconsistent, duplicated, difficult to review, and difficult for AI agents to retrieve accurately.

---

# Goals

This RFC establishes:
  - RFC purpose
  - RFC lifecycle
  - RFC required sections
  - RFC numbering ranges
  - Source-of-truth update rules
  - Review requirements
  - Plugin RFC requirements
  - DMP and dstoreOps RFC requirements
  - AI RFC requirements
  - AI-assisted authoring boundaries

---

# Non-Goals

This RFC does not:
  - Decide runtime architecture
  - Decide data plane architecture
  - Decide control plane architecture
  - Decide plugin framework internals
  - Decide AI Control Plane internals
  - Approve autonomous AI behavior
  - Replace architecture documents
  - Replace specification documents

---

# Proposal

Adopt the SDE RFC Framework as defined in:
  - README.md
  - index.md
  - rfc-template.md
  - rfc-0000-authoring-standard.md

All architecture-significant decisions must use the RFC process.

All accepted RFCs must update affected source-of-truth documents.

---

# Authoring Rules

Rule 1:
  One RFC should contain one major decision or tightly related decision set.

Rule 2:
  RFCs must state goals and non-goals.

Rule 3:
  RFCs must identify architecture impact.

Rule 4:
  RFCs must identify specification impact.

Rule 5:
  RFCs must identify security impact.

Rule 6:
  RFCs must identify operational impact.

Rule 7:
  RFCs must identify compatibility impact.

Rule 8:
  RFCs must list source-of-truth documents that must be updated.

Rule 9:
  RFCs must record alternatives considered.

Rule 10:
  RFCs must record risks and mitigation.

Rule 11:
  RFCs must be written for both human and AI retrieval.

Rule 12:
  RFCs must not redefine canonical terms without updating glossary.

Rule 13:
  RFCs must not bypass accepted architecture boundaries.

Rule 14:
  Accepted RFCs may become Stable only after source-of-truth updates are complete.

---

# AI-Reasoning Authoring Rules

RFCs must be easy for AI agents to retrieve and reason over.

Therefore:
  - Use explicit headings.
  - Use canonical terminology.
  - Avoid hidden assumptions.
  - State "No impact" explicitly where applicable.
  - Include relationships to architecture and specification files.
  - Keep each RFC focused.
  - Include non-goals.
  - Include rejected alternatives.
  - Avoid ambiguous terms such as "database platform", "engine data plane", or "management plugin".

---

# Canonical Terms

RFCs must use terms from the glossary.

Examples:
  - Sovrunn Data Engine
  - SDE
  - SDE Control Plane
  - SDE Data Plane
  - SDE Runtime
  - Datastore Data Plane
  - Downstream Datastore
  - Downstream Engine
  - Protocol Plugin
  - Engine Plugin
  - Datastore Operator Plugin
  - Infrastructure Provider
  - Foundation Provider
  - Datastore Management Plane
  - dstoreOps
  - AI Control Plane
  - Tenant AI Agent

Deprecated terms must not be introduced in new RFCs.

---

# Numbering Standard

RFC numbers must be assigned from the ranges in README.md and index.md.

The number must indicate the primary decision domain.

If an RFC affects multiple domains, choose the range based on the primary owner.

Examples:
  - Protocol Plugin framework decision: 1000-1199
  - PostgreSQL Protocol Plugin decision: 1200-1299
  - Engine Plugin framework decision: 2000-2199
  - DMP framework decision: 3000-3199
  - Datastore Operator Plugin framework decision: 3200-3399
  - dstoreOps framework decision: 6000-6199
  - AI Control Plane decision: 9000-9199
  - Tenant AI Agent decision: 9200-9399

---

# Required Sections

Every RFC must contain:
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

Plugin RFCs must also include:
  - Plugin Impact

AI RFCs must also include:
  - AI Impact

DMP or dstoreOps RFCs must also include:
  - DMP and dstoreOps Impact

---

# Source-of-Truth Rule

RFCs are not the source of current truth.

Architecture, specification, implementation, and operations documents are the source of current truth.

When an RFC is Accepted:
  - Required source-of-truth documents must be updated.
  - Index must be updated.
  - Related documents must link back to the RFC where appropriate.

An RFC becomes Stable only when:
  - Decision is accepted.
  - Required source-of-truth updates are complete.
  - MkDocs or documentation validation passes.
  - Architecture consistency has been checked.

---

# Review Standard

Architecture review is always required.

Security review is required when RFC affects:
  - tenant isolation
  - identity
  - authorization
  - policy
  - secrets
  - audit
  - plugin admission
  - AI action authority
  - destructive operations
  - runtime execution boundaries

Operations review is required when RFC affects:
  - deployment
  - lifecycle management
  - backup
  - restore
  - scaling
  - patch
  - upgrade
  - monitoring
  - incidents
  - dstoreOps
  - AI tuning or remediation

Specification review is required when RFC affects:
  - versioning
  - serialization
  - protocol contracts
  - engine contracts
  - capability contracts
  - manifests
  - compatibility rules

Tenant experience review is required when RFC affects:
  - tenant APIs
  - Tenant AI Agent
  - tenant onboarding
  - tenant workflows
  - tenant namespace
  - tenant-visible behavior

---

# Plugin RFC Standard

Plugin RFCs must clearly identify plugin type:
  - Protocol Plugin
  - Engine Plugin
  - Datastore Operator Plugin
  - Foundation Provider
  - Infrastructure Provider

Plugin RFCs must state:
  - Runtime boundary
  - Control Plane boundary
  - Manifest impact
  - Registry impact
  - Validation requirements
  - Conformance requirements
  - Compatibility matrix
  - Failure behavior

Plugin RFCs must not blur boundaries.

Protocol Plugin:
  - Does not produce Execution Plan.
  - Does not invoke Engine Plugin.
  - Does not access Downstream Datastore.

Engine Plugin:
  - Does not parse client protocol.
  - Does not manage datastore lifecycle.
  - Does not invoke Infrastructure Provider.

Datastore Operator Plugin:
  - Does not execute tenant data-plane requests.
  - Does not replace Engine Plugin.
  - Does not bypass DMP workflows, policy, or audit.

---

# DMP and dstoreOps RFC Standard

DMP and dstoreOps RFCs must state:
  - Tenant namespace impact
  - DatastoreRequest impact
  - DatastoreInstance impact
  - DatastoreProfile impact
  - DatastorePolicy impact
  - Datastore Operator Plugin impact
  - Infrastructure Provider impact
  - Workflow impact
  - Backup and restore impact
  - Monitoring impact
  - Tenant isolation impact

DMP RFCs must preserve:
  - DMP as lifecycle management plane
  - DMP as non-data-plane control surface
  - DMP as workflow-governed
  - DMP as policy-audited
  - DMP as plugin-integrated but not plugin-bypassed

---

# AI RFC Standard

AI RFCs must state:
  - AI scope
  - AI action class
  - Tenant isolation impact
  - Policy requirements
  - Approval requirements
  - Audit requirements
  - Rollback requirements
  - Safety boundaries
  - Whether AI acts only as recommender or may initiate workflow

AI RFCs must preserve:
  - AI Control Plane is optional and pluggable until accepted otherwise.
  - AI Control Plane is part of SDE Control Plane, not SDE Data Plane.
  - Tenant AI Agent is tenant-scoped.
  - AI must not bypass Control Plane services.
  - AI must not directly manage Downstream Datastores.
  - AI must not directly invoke Datastore Operator Plugins.
  - AI must not directly invoke Infrastructure Providers.
  - AI-generated artifacts are untrusted until validated.
  - AI must not autonomously perform destructive actions.

---

# AI-Assisted RFC Drafting

AI may help:
  - Draft RFC text
  - Summarize context
  - Compare alternatives
  - Identify impacted documents
  - Generate review checklists
  - Detect terminology inconsistencies

AI must not:
  - Approve RFCs
  - Replace review
  - Decide security posture
  - Decide tenant isolation policy
  - Approve plugin admission
  - Approve destructive operations

---

# Architecture Impact

This RFC establishes RFC governance.

Architecture impact:
  - Adds formal decision process
  - Adds scalable numbering
  - Adds domain-specific review expectations
  - Adds AI-specific governance rules

---

# Specification Impact

This RFC does not change SDE technical specifications.

It requires future specification changes to go through RFC review when architecture-significant.

---

# Security Impact

This RFC improves security governance by requiring security review for security-relevant RFCs and AI-relevant RFCs.

---

# Operational Impact

This RFC improves operational governance by requiring operations review for DMP, dstoreOps, lifecycle, incident, AI tuning, and remediation changes.

---

# Compatibility Impact

This RFC is backward-compatible with existing documentation.

Existing RFC drafts should be renumbered or updated to match this framework where needed.

---

# Source-of-Truth Updates

Updated:
  - docs/rfc/README.md
  - docs/rfc/index.md
  - docs/rfc/rfc-template.md
  - docs/rfc/rfc-0000-authoring-standard.md

---

# Alternatives Considered

Alternative:
  Use small 0000-0999 range for all RFCs.

Reason Rejected:
  SDE is extensible across Protocol Plugins, Engine Plugins, DMP, Datastore Operator Plugins, Infrastructure Providers, Foundation Providers, dstoreOps, and AI. Small ranges would become crowded.

Alternative:
  Avoid AI RFC ranges until AI is fully designed.

Reason Rejected:
  AI Control Plane and Tenant AI Agent are strategic future capabilities. Reserving ranges now avoids collision without forcing detailed design now.

Alternative:
  Allow AI-generated plugin deployment as an RFC category.

Reason Rejected:
  On-the-fly plugin generation and deployment is risky. AI-assisted artifact generation is reserved for future drafts, with validation required.

---

# Risks

Risk:
  RFC process becomes too heavy.

Mitigation:
  Apply RFCs only to architecture-significant decisions.

Risk:
  RFCs become stale.

Mitigation:
  RFCs explain why; source-of-truth docs define what is true.

Risk:
  AI-generated RFCs introduce hallucinated architecture.

Mitigation:
  Require review and source-of-truth validation.

---

# Decision

Decision:
  Accepted.

Rationale:
  SDE requires scalable, reviewable, AI-reasoning-friendly governance for long-term architecture and extensibility.
