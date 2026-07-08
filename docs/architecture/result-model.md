# Result Model

Document

- ID: result-model
- Version: 1.0
- Status: Stable

Purpose

- Define canonical runtime result model
- Define result boundaries
- Define result ownership
- Define result propagation across Sovrunn Runtime

Definition

Result Model defines the canonical structure used by Sovrunn Runtime to represent execution outcomes.

A Result is produced by execution and returned through the runtime without changing SIR semantics.

Principles

MUST

- Preserve semantic intent
- Preserve execution outcome
- Be protocol independent
- Be engine independent
- Be deterministic
- Be serializable
- Separate successful results from errors

MUST NOT

- Expose downstream engine internal result formats directly
- Expose protocol-specific response formats directly
- Modify SIR
- Modify Execution Plan
- Hide partial execution state
- Hide execution failure

Result Categories

Result

MAY be one of

- Execution Result
- Operation Result
- Data Result
- Mutation Result
- Metadata Result
- Streaming Result
- Empty Result

Execution Result

Execution Result represents the outcome of executing an Execution Plan.

Execution Result

MUST contain

- Execution Identifier
- Execution State
- Operation Results
- Result Metadata

MAY contain

- Data Result
- Mutation Result
- Warning
- Partial Result Indicator

Execution State

Execution State

MUST be one of

- Succeeded
- Failed
- Partial
- Cancelled
- Timed Out

Operation Result

Operation Result represents the outcome of one operation inside an Execution Plan.

Operation Result

MUST contain

- Operation Identifier
- Operation State
- Required Capability
- Result Metadata

MAY contain

- Data Result
- Mutation Result
- Engine Reference
- Warning

Data Result

Data Result represents returned data.

Data Result

MUST contain

- Schema
- Rows or Batches
- Format
- Result Metadata

Data Result

MAY contain

- Cursor Reference
- Continuation Token
- Row Count
- Batch Count

Data Representation

Data Result

MUST use adopted data representation standards.

Primary representation

- Apache Arrow

Secondary representation

- JSON

Rules

MUST

- Use Arrow for typed tabular or columnar data
- Use JSON for debugging, testing, and human-readable output
- Preserve type information
- Preserve nullability
- Preserve ordering when semantically required

MUST NOT

- Invent custom runtime data encoding
- Lose type information
- Convert data silently without declared format

Mutation Result

Mutation Result represents state-changing operation outcome.

Mutation Result

MUST contain

- Affected Resource
- Mutation Type
- Affected Count
- Mutation State

MAY contain

- Generated Identifier
- Version
- Timestamp
- Engine Reference

Metadata Result

Metadata Result represents metadata returned by discovery, introspection, or runtime operations.

Metadata Result

MUST contain

- Metadata Type
- Metadata Values

MAY contain

- Resource Reference
- Capability Reference
- Version

Streaming Result

Streaming Result represents a result that is delivered incrementally.

Streaming Result

MUST contain

- Stream Identifier
- Stream State
- Schema or Event Type
- Continuation Reference

MAY contain

- Checkpoint
- Offset
- Partition
- Watermark

Empty Result

Empty Result represents successful execution without returned data.

Empty Result

MUST contain

- Execution Identifier
- Execution State
- Result Metadata

Result Metadata

Result Metadata

MAY contain

- Execution Duration
- Row Count
- Batch Count
- Engine Reference
- Capability Reference
- Trace Identifier
- Warnings
- Partial Result Indicator

Result Metadata

MUST NOT

- Change result semantics
- Contain business data
- Hide errors

Result Propagation

Result flows through

- Engine Plugin
- Engine Runtime
- Data Kernel
- Protocol Runtime
- Protocol Plugin
- Client

Propagation Rules

Engine Plugin

MUST

- Convert downstream engine native result into Sovrunn Result Model
- Preserve semantic equivalence
- Preserve data type information
- Preserve mutation outcome

MUST NOT

- Return raw native engine result directly to Data Kernel
- Hide downstream engine execution failure

Engine Runtime

MUST

- Receive Sovrunn Result from Engine Plugin
- Validate result structure
- Return result to Data Kernel

MUST NOT

- Convert result into protocol format
- Modify semantic meaning

Data Kernel

MUST

- Aggregate Operation Results
- Produce Execution Result
- Preserve operation ordering when required
- Preserve partial result state

MUST NOT

- Convert result into protocol format
- Hide failed operations

Protocol Runtime

MUST

- Receive Execution Result
- Delegate protocol response formatting to Protocol Plugin

MUST NOT

- Expose Sovrunn internal result format directly unless explicitly requested

Protocol Plugin

MUST

- Convert Sovrunn Result into protocol-compatible response
- Preserve client-visible semantics
- Preserve errors using Error Model

MUST NOT

- Invent successful result when execution failed
- Hide partial result state

Partial Results

Partial Result means some operations succeeded and some operations failed or were cancelled.

Partial Result

MUST

- Identify completed operations
- Identify failed operations
- Preserve deterministic state
- Include Error Model reference

MUST NOT

- Be reported as full success
- Be discarded silently

Cursor Results

Cursor Result represents deferred or incremental access to a Data Result.

Cursor Result

MUST contain

- Cursor Reference
- Schema
- Continuation State
- Expiration Policy

Cursor Result

MUST NOT

- Own Session lifecycle
- Own Transaction lifecycle

Session Runtime owns session lifecycle.

Transaction Runtime owns transaction lifecycle.

Result and Error Separation

Result Model represents successful, partial, or cancelled execution outcomes.

Error Model represents failure details.

Result Model

MUST reference Error Model when failure exists.

Result Model

MUST NOT redefine

- Error codes
- Error categories
- Error propagation
- Retry classification

Ownership

Sovrunn owns

- Result Model
- Execution Result
- Operation Result
- Data Result
- Mutation Result
- Metadata Result
- Streaming Result
- Empty Result
- Result propagation rules

Engine Plugin owns

- Native engine result conversion
- Native type mapping
- Native mutation outcome mapping

Downstream Engine owns

- Native result production
- Native data formats
- Native execution output

Protocol Plugin owns

- Protocol-specific response formatting

References

- architecture.md
- runtime.md
- execution-plan.md
- execution-context.md
- data-kernel.md
- engine-runtime.md
- plugin-runtime.md
- protocol-runtime.md
- specifications/reuse/adopted-standards.md
- specifications/sir/serialization.md
