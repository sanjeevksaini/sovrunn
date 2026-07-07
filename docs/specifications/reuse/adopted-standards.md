# Adopted Standards

Document
- ID: adopted-standards
- Version: 1.0
- Status: Stable

Purpose
- Define normative specifications adopted by Sovrunn
- Define adoption scope
- Define specification authority

Rules

MUST
- Adopt authoritative specifications
- Preserve specification conformance
- Reference adopted specifications

MUST NOT
- Redefine adopted specifications
- Fork adopted specifications
- Duplicate transitive dependencies

Standards

### Semantic Model

- Name: Apache Arrow
  Authority: Apache Software Foundation
  Type: Open Specification
  Scope:
  - Primitive Types
  - Type System
  - Memory Model
  - Columnar Format
  - IPC Format

- Name: Substrait
  Authority: Linux Foundation
  Type: Open Specification
  Scope:
  - Relational Algebra
  - Expression Model
  - Function Model
  - Type System
  - Plan Interchange

### Client Protocols

- Name: PostgreSQL Wire Protocol
  Authority: PostgreSQL Global Development Group
  Type: Open Specification
  Scope:
  - PostgreSQL Client Compatibility

- Name: MySQL Client/Server Protocol
  Authority: Oracle
  Type: Open Specification
  Scope:
  - MySQL Client Compatibility

- Name: MongoDB Wire Protocol
  Authority: MongoDB Inc.
  Type: Open Specification
  Scope:
  - MongoDB Client Compatibility

- Name: Redis Serialization Protocol (RESP)
  Authority: Redis
  Type: Open Specification
  Scope:
  - Redis Client Compatibility

### Object Storage

- Name: Amazon S3 API
  Authority: Amazon Web Services
  Type: De Facto Standard
  Scope:
  - Object Storage
  - Blob Storage
  - Multipart Upload
  - Object Metadata

### Storage

- Name: Apache Parquet
  Authority: Apache Software Foundation
  Type: Open Specification
  Scope:
  - Columnar Storage

- Name: Apache Iceberg
  Authority: Apache Software Foundation
  Type: Open Specification
  Scope:
  - Table Format

### Observability

- Name: OpenTelemetry
  Authority: CNCF
  Type: Open Specification
  Scope:
  - Tracing
  - Metrics
  - Logging

### API

- Name: OpenAPI Specification
  Authority: OpenAPI Initiative
  Type: Open Specification
  Scope:
  - Management API

- Name: JSON Schema
  Authority: JSON Schema Organization
  Type: Open Specification
  Scope:
  - Configuration
  - Validation

### Serialization

- Name: Protocol Buffers
  Authority: Google
  Type: Open Specification
  Scope:
  - Internal Serialization

### Containers

- Name: OCI Image Specification
  Authority: Open Container Initiative
  Type: Open Specification
  Scope:
  - Container Images

- Name: OCI Runtime Specification
  Authority: Open Container Initiative
  Type: Open Specification
  Scope:
  - Container Runtime

- Name: Kubernetes API
  Authority: CNCF
  Type: Open Specification
  Scope:
  - Orchestration

### Time

- Name: RFC 3339
  Authority: IETF
  Type: RFC
  Scope:
  - Timestamp Representation
  - API Contracts
  - Configuration
  - Logging
  - Telemetry
  - Audit

- Name: ISO 8601
  Authority: ISO
  Type: International Standard
  Scope:
  - Date Representation
  - Time Representation
  - Duration
  - Interval

### Identity

- Name: RFC 9562
  Authority: IETF
  Type: RFC
  Scope:
  - Resource Identifier
  - Object Identifier
  - Request Identifier
  - Session Identifier
  - Transaction Identifier

Compliance

Every Specification
- MUST reference adopted standards

Every Architecture
- MUST conform to adopted standards

Every Implementation
- MUST implement adopted standards
