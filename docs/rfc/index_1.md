# RFC Index

Document:
  ID: rfc-index
  Title: RFC Index
  Parent: rfc-framework
  Owner: SDE Architecture Council
  Layer: Governance
  Type: Index
  Version: 2.0
  Status: Stable

Purpose:
  - List SDE RFCs by number, title, status, owner, and affected source-of-truth areas
  - Provide a stable navigation surface for humans and AI agents
  - Prevent duplicate RFC numbers and overlapping decision ownership

Index Rules:
  - Every RFC must be listed here before review.
  - RFC number must not be reused.
  - Superseded RFCs remain listed.
  - Rejected RFCs remain listed.
  - Stable RFCs must reference updated source-of-truth documents.

---

# RFC Status Values

Draft:
  RFC is being written.

Review:
  RFC is under review.

Accepted:
  RFC has been accepted, but source-of-truth updates may still be pending.

Stable:
  RFC has been accepted and source-of-truth documents have been updated.

Rejected:
  RFC was not accepted.

Superseded:
  RFC was replaced by a newer RFC.

Deprecated:
  RFC was previously accepted but is no longer recommended.

---

# Active RFCs

| RFC | Title | Status | Owner | Source-of-Truth Impact |
|---|---|---:|---|---|
| RFC-0000 | Authoring Standard | Stable | SDE Architecture Council | docs/rfc |
| RFC-9000 | AI Control Plane Reservation | Reserved | SDE Architecture Council | docs/architecture/control-plane/ai-control-plane.md |

---

# Reserved RFC Ranges

## 0000-0999 Governance, Foundation, Core Architecture

| Range | Area |
|---:|---|
| 0000-0099 | RFC Governance and Authoring Standards |
| 0100-0199 | Foundation, Glossary, Ontology, Architecture Principles |
| 0200-0299 | Runtime Architecture |
| 0300-0399 | SDE Data Plane Architecture |
| 0400-0499 | SDE Control Plane Architecture |
| 0500-0599 | Specification Framework |
| 0600-0799 | Capability Specifications |
| 0800-0999 | Reserved Core Architecture |

## 1000-1999 Protocol Plugins

| Range | Area |
|---:|---|
| 1000-1199 | Protocol Plugin Framework |
| 1200-1299 | PostgreSQL Protocol Plugin |
| 1300-1399 | MySQL Protocol Plugin |
| 1400-1499 | MongoDB Protocol Plugin |
| 1500-1599 | Redis Protocol Plugin |
| 1600-1699 | REST, gRPC, and Native Protocol Plugins |
| 1700-1999 | Reserved Protocol Plugins |

## 2000-2999 Engine Plugins

| Range | Area |
|---:|---|
| 2000-2199 | Engine Plugin Framework |
| 2200-2299 | PostgreSQL Engine Plugin |
| 2300-2399 | MySQL Engine Plugin |
| 2400-2499 | MongoDB Engine Plugin |
| 2500-2599 | Redis Engine Plugin |
| 2600-2699 | Object and Table Engine Plugins |
| 2700-2799 | Search, Vector, and Graph Engine Plugins |
| 2800-2899 | Cassandra and Distributed Datastore Engine Plugins |
| 2900-2999 | Reserved Engine Plugins |

## 3000-3999 Datastore Management Plane and Datastore Operator Plugins

| Range | Area |
|---:|---|
| 3000-3199 | Datastore Management Plane Framework |
| 3200-3399 | Datastore Operator Plugin Framework |
| 3400-3499 | Datastore Operator Plugin Implementations |
| 3500-3699 | DMP Lifecycle Controllers |
| 3700-3899 | DMP Policy, Credentials, Namespace, and Tenant Operations |
| 3900-3999 | Reserved DMP |

## 4000-4999 Infrastructure Providers

| Range | Area |
|---:|---|
| 4000-4199 | Infrastructure Provider Framework |
| 4200-4299 | Kubernetes Infrastructure Provider |
| 4300-4399 | AWS Infrastructure Provider |
| 4400-4499 | Azure Infrastructure Provider |
| 4500-4599 | GCP Infrastructure Provider |
| 4600-4699 | VMware, Bare Metal, Private Cloud, Sovereign Cloud, and Hybrid Cloud Providers |
| 4700-4999 | Reserved Infrastructure Providers |

## 5000-5999 Foundation Services and Foundation Providers

| Range | Area |
|---:|---|
| 5000-5199 | Foundation Service Framework |
| 5200-5499 | Foundation Provider Framework and Implementations |
| 5500-5999 | Reserved Foundation |

## 6000-6999 dstoreOps

| Range | Area |
|---:|---|
| 6000-6199 | dstoreOps Framework |
| 6200-6299 | dstoreOps Policy Model |
| 6300-6399 | Provisioning Workflows |
| 6400-6499 | Scaling Workflows |
| 6500-6599 | Backup and Restore Workflows |
| 6600-6699 | Patch and Upgrade Workflows |
| 6700-6799 | Monitoring, Health, and Retirement Workflows |
| 6800-6999 | Reserved dstoreOps |

## 7000-7999 Security, Policy, Governance, and Compliance

| Range | Area |
|---:|---|
| 7000-7199 | Security Architecture |
| 7200-7399 | Policy and Governance |
| 7400-7599 | Audit and Compliance |
| 7600-7999 | Reserved Security and Governance |

## 8000-8999 Implementation

| Range | Area |
|---:|---|
| 8000-8199 | Repository and Module Structure |
| 8200-8399 | Coding Standards and SDKs |
| 8400-8599 | Build, Test, and Release |
| 8600-8799 | Deployment and Environment Structure |
| 8800-8999 | Reserved Implementation |

## 9000-9999 AI Control Plane and AI-Native Operations

| Range | Area |
|---:|---|
| 9000-9199 | AI Control Plane |
| 9200-9399 | Tenant AI Agent |
| 9400-9499 | AI Observation and Stability |
| 9500-9599 | AI Safe Autotuning |
| 9600-9699 | AI Remediation and Runbooks |
| 9700-9799 | AI Evaluation and Safety |
| 9800-9899 | AI-Assisted Artifact Generation |
| 9900-9999 | Reserved AI Future Use |

---

# Planned Initial RFCs

| RFC | Proposed Title | Area |
|---:|---|---|
| 0100 | SDE Architecture Documentation Backbone | Foundation |
| 0200 | SDE Runtime Core Model | Runtime |
| 0300 | SDE Data Plane Request Lifecycle | Data Plane |
| 0400 | SDE Control Plane Core Model | Control Plane |
| 0500 | SDE Specification Framework | Specifications |
| 0600 | SDE Capability Model | Capabilities |
| 1000 | Protocol Plugin Framework | Protocol Plugins |
| 2000 | Engine Plugin Framework | Engine Plugins |
| 3000 | Datastore Management Plane Framework | DMP |
| 3200 | Datastore Operator Plugin Framework | DMP Plugins |
| 4000 | Infrastructure Provider Framework | Infrastructure Providers |
| 5000 | Foundation Service Framework | Foundation |
| 5200 | Foundation Provider Framework | Foundation Providers |
| 6000 | dstoreOps Framework | dstoreOps |
| 7000 | Tenant Isolation and Security Model | Security |
| 8000 | Repository and Module Structure | Implementation |
| 9000 | AI Control Plane Reservation | AI Control Plane |

---

# Superseded RFCs

None.

---

# Rejected RFCs

None.

---

# Deprecated RFCs

None.
