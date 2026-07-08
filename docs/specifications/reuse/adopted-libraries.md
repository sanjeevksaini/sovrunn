# Adopted Libraries

Document
- ID: adopted-libraries
- Version: 1.0
- Status: Stable

Purpose
- Define implementation libraries adopted by Sovrunn
- Prevent unnecessary implementation
- Promote mature open source software

Rules

MUST
- Adopt mature open source libraries
- Prefer actively maintained projects
- Isolate libraries behind adapters
- Replace libraries without affecting platform semantics

MUST NOT
- Fork adopted libraries
- Expose library APIs as platform APIs
- Couple architecture to library implementation

Libraries

### SQL Parsing

- Name: libpg_query
  Authority: PostgreSQL Global Development Group
  Scope:
  - PostgreSQL SQL Parsing
  - SQL AST

### PostgreSQL Protocol

- Name: pgproto3
  Authority: jackc
  Scope:
  - PostgreSQL Wire Protocol

### Columnar Data

- Name: Apache Arrow
  Authority: Apache Software Foundation
  Scope:
  - Memory Representation
  - IPC
  - Compute

### Serialization

- Name: Protocol Buffers
  Authority: Google
  Scope:
  - Internal Serialization

### Observability

- Name: OpenTelemetry SDK
  Authority: CNCF
  Scope:
  - Tracing
  - Metrics
  - Logging

### Object Storage

- Name: MinIO SDK
  Authority: MinIO
  Scope:
  - S3 Compatible Object Storage

### Container Runtime

- Name: containerd
  Authority: CNCF
  Scope:
  - Container Runtime

### Kubernetes

- Name: client-go
  Authority: Kubernetes
  Scope:
  - Kubernetes API

Compliance

Every Implementation

MUST
- Abstract adopted libraries
- Replace libraries without changing platform behavior
- Preserve platform contracts
