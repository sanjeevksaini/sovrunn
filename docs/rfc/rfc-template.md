# RFC-NNNN: Title

Document:
  ID: rfc-NNNN
  Title: Title
  Parent: rfc-index
  Owner: TBD
  Layer: TBD
  Type: RFC
  Version: 0.1
  Status: Draft

RFC Number:
  NNNN

Status:
  Draft

Created:
  YYYY-MM-DD

Updated:
  YYYY-MM-DD

Authors:
  - TBD

Reviewers:
  - Architecture:
  - Security:
  - Operations:
  - Specification:
  - Implementation:
  - Tenant Experience:

Affected Areas:
  - TBD

Source-of-Truth Documents:
  - TBD

Supersedes:
  - None

Superseded By:
  - None

---

# Summary

Briefly describe the decision or proposal.

The summary must be short enough for quick retrieval but precise enough to identify the decision.

---

# Problem Statement

Describe the problem.

Include:
  - Current limitation
  - Why it matters
  - What breaks or becomes risky without this RFC
  - Who is affected

---

# Goals

This RFC intends to:
  - Goal 1
  - Goal 2
  - Goal 3

---

# Non-Goals

This RFC does not intend to:
  - Non-goal 1
  - Non-goal 2
  - Non-goal 3

---

# Context

Describe the relevant background.

Include links or references to:
  - Architecture documents
  - Specification documents
  - Previous RFCs
  - Implementation modules
  - Operational concerns
  - Security constraints

---

# Proposal

Describe the proposed decision.

Use clear language.

Avoid mixing multiple independent decisions into one RFC.

---

# Architecture Impact

Describe impact on:
  - SDE Control Plane
  - SDE Data Plane
  - SDE Runtime
  - Datastore Management Plane
  - Foundation Services
  - Foundation Providers
  - Protocol Plugins
  - Engine Plugins
  - Datastore Operator Plugins
  - Infrastructure Providers
  - dstoreOps
  - AI Control Plane
  - Tenant AI Agent

Use "No impact" explicitly when applicable.

---

# Specification Impact

Describe impact on:
  - Versioning Specification
  - Serialization Specification
  - Capability Specifications
  - Protocol Specifications
  - Engine Specifications
  - Manifest schemas
  - Compatibility rules

Use "No impact" explicitly when applicable.

---

# Security Impact

Describe impact on:
  - Tenant isolation
  - Identity
  - Authorization
  - Policy
  - Secrets
  - Audit
  - Encryption
  - Approval workflows
  - Plugin admission
  - AI action authority

Use "No impact" explicitly when applicable.

---

# Operational Impact

Describe impact on:
  - Deployment
  - Monitoring
  - Scaling
  - Backup
  - Restore
  - Patch
  - Upgrade
  - Incident response
  - DR
  - SLOs
  - Runbooks
  - dstoreOps workflows

Use "No impact" explicitly when applicable.

---

# Compatibility Impact

Describe compatibility impact on:
  - Existing tenants
  - Existing APIs
  - Existing plugins
  - Existing manifests
  - Existing runtime behavior
  - Existing datastore operations
  - Existing workflows

State whether this is:
  - Backward-compatible
  - Forward-compatible
  - Breaking
  - Migration-required

---

# Plugin Impact

Complete this section if the RFC affects plugins.

Plugin Type:
  - Protocol Plugin
  - Engine Plugin
  - Datastore Operator Plugin
  - Foundation Provider
  - Infrastructure Provider
  - None

Manifest Impact:
  - TBD

Registry Impact:
  - TBD

Validation Impact:
  - TBD

Conformance Impact:
  - TBD

Runtime Boundary:
  - TBD

Control Plane Boundary:
  - TBD

---

# AI Impact

Complete this section if the RFC affects AI Control Plane or Tenant AI Agent.

AI Scope:
  - No impact
  - AI Control Plane
  - Tenant AI Agent
  - AI Observation
  - AI Safe Autotuning
  - AI Remediation
  - AI Evaluation
  - AI-Assisted Artifact Generation

AI Action Classes:
  - Class 0: Read-only observation and explanation
  - Class 1: Generate draft recommendations or configuration artifacts
  - Class 2: Validate or simulate proposed action
  - Class 3: Apply reversible, non-destructive, policy-approved tuning
  - Class 4: Controlled operational change requiring workflow approval
  - Class 5: Destructive or irreversible action requiring explicit human approval and safety workflow

Policy Requirements:
  - TBD

Approval Requirements:
  - TBD

Audit Requirements:
  - TBD

Rollback Requirements:
  - TBD

Tenant Isolation Impact:
  - TBD

Safety Boundaries:
  - TBD

---

# DMP and dstoreOps Impact

Complete this section if the RFC affects Datastore Management Plane or dstoreOps.

Tenant Namespace Impact:
  - TBD

DatastoreRequest Impact:
  - TBD

DatastoreInstance Impact:
  - TBD

Datastore Operator Plugin Impact:
  - TBD

Infrastructure Provider Impact:
  - TBD

Workflow Impact:
  - TBD

Backup/Restore Impact:
  - TBD

Monitoring Impact:
  - TBD

---

# Source-of-Truth Updates

List documents that must be updated if this RFC is accepted.

Required Updates:
  - docs/...

Optional Updates:
  - docs/...

Stable Status Requirement:
  This RFC can become Stable only after required source-of-truth updates are completed and validated.

---

# Alternatives Considered

## Alternative 1

Description:
  TBD

Pros:
  - TBD

Cons:
  - TBD

Reason Rejected:
  TBD

## Alternative 2

Description:
  TBD

Pros:
  - TBD

Cons:
  - TBD

Reason Rejected:
  TBD

---

# Risks

Risk:
  - TBD

Mitigation:
  - TBD

Residual Risk:
  - TBD

---

# Rollout Plan

Phase 1:
  - TBD

Phase 2:
  - TBD

Phase 3:
  - TBD

Rollback:
  - TBD

---

# Migration Plan

Required:
  - Yes / No

Steps:
  - TBD

Tenant Impact:
  - TBD

Operator Impact:
  - TBD

---

# Validation Plan

Architecture Validation:
  - TBD

Specification Validation:
  - TBD

Security Validation:
  - TBD

Operational Validation:
  - TBD

Implementation Validation:
  - TBD

AI Validation:
  - TBD

---

# Open Questions

- Question 1
- Question 2
- Question 3

---

# Decision

Decision:
  Pending

Rationale:
  TBD

Approval:
  - Architecture:
  - Security:
  - Operations:
  - Product:
  - Implementation:
