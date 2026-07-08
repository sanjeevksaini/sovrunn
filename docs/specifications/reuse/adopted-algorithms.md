# Adopted Algorithms

Document
- ID: adopted-algorithms
- Version: 1.0
- Status: Stable

Purpose
- Define normative algorithms adopted by Sovrunn
- Define computational foundations
- Prevent algorithm reinvention

Rules

MUST
- Adopt proven algorithms
- Preserve algorithm semantics
- Reference authoritative publications

MUST NOT
- Reinvent adopted algorithms
- Modify adopted semantics
- Treat implementation optimizations as platform algorithms

Algorithms

### Relational Processing

- Name: Relational Algebra
  Authority: A Relational Model of Data for Large Shared Data Banks (1970)
  Scope:
  - Relational Semantics
  - Logical Transformation

### Query Optimization

- Name: Volcano Optimizer
  Authority: The Volcano Optimizer Generator (1993)
  Scope:
  - Cost Based Optimization
  - Physical Planning

- Name: Cascades Framework
  Authority: The Cascades Framework for Query Optimization (1995)
  Scope:
  - Rule Based Optimization
  - Memoization
  - Search Space Exploration

### Transactions

- Name: Multi-Version Concurrency Control (MVCC)
  Authority: Database Research Community
  Scope:
  - Concurrent Transactions
  - Snapshot Isolation

### Distributed Consensus

- Name: Raft
  Authority: In Search of an Understandable Consensus Algorithm (2014)
  Scope:
  - Consensus
  - Leader Election
  - Replicated State Machine

### Data Distribution

- Name: Consistent Hashing
  Authority: Consistent Hashing and Random Trees (1997)
  Scope:
  - Data Partitioning
  - Elastic Scaling

- Name: Rendezvous Hashing
  Authority: A Name-Based Mapping Scheme (1996)
  Scope:
  - Data Placement
  - Request Routing

Compliance

Every Architecture
- MUST preserve adopted algorithm semantics

Every Implementation
- MUST implement equivalent behavior
