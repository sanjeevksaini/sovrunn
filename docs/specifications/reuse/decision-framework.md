# Decision Framework

Document
- ID: reuse-decision-framework
- Version: 1.0
- Status: Stable

Purpose
- Define deterministic adoption decisions
- Ensure consistent architecture evolution
- Prevent unnecessary ownership

Decision Flow

Step 1

Question
- Does authoritative knowledge already exist

Action
- Adopt knowledge

Examples
- Relational Algebra
- Compiler Theory
- Consensus Theory

Step 2

Question
- Does an open specification already exist

Action
- Adopt specification

Examples
- Apache Arrow
- OpenTelemetry
- OpenAPI
- JSON Schema

Step 3

Question
- Does a proven architecture pattern exist

Action
- Adopt architecture pattern

Examples
- Microkernel
- Compiler Pipeline
- Plugin Architecture
- Shared Nothing

Step 4

Question
- Does a proven algorithm exist

Action
- Adopt algorithm

Examples
- Volcano Optimizer
- Cascades
- MVCC
- Raft

Step 5

Question
- Does a mature implementation exist

Action
- Adopt implementation

Examples
- libpg_query
- pgproto3
- Apache Arrow
- OpenTelemetry SDK

Step 6

Question
- Can the adopted implementation be extended

Action
- Extend

Step 7

Question
- Is new implementation required

Action
- Build only differentiated capability

Decision Rules

MUST
- Evaluate every step sequentially
- Stop when an acceptable solution exists
- Justify every build decision
- Prefer adoption over implementation
- Prefer extension over replacement

MUST NOT
- Skip decision steps
- Build before evaluation
- Duplicate adopted capability
- Fork mature implementations without justification

Exit Criteria

Adopt
- Existing solution satisfies requirements

Extend
- Existing solution satisfies most requirements

Build
- No acceptable solution exists
- Differentiated capability required

RFC Requirements

Every build decision MUST include
- Alternatives evaluated
- Rejected alternatives
- Justification
- Expected differentiation
- Long-term ownership impact
