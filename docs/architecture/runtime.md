# Runtime

Document:
  ID: runtime
  Version: 1.0
  Status: Stable

Purpose:
  - Canonical runtime model
  - Runtime component responsibilities
  - Runtime dependency graph

Rules:
  - Runtime is protocol independent
  - Runtime is engine independent
  - Dependencies are directional
  - Cyclic dependencies prohibited
  - Specifications never depend on runtime
  - Applications never depend on engines
  - Engines never depend on applications

Components:

  Application:
    Purpose: Consume platform capabilities
    DependsOn:
      - Protocol
    Produces:
      - Request

  Protocol:
    Purpose: Define client interface
    DependsOn: []
    Consumes:
      - Request
    Produces:
      - Protocol Request

  Protocol Runtime:
    Purpose: Translate protocol request into SIR
    DependsOn:
      - Protocol
      - SIR
    Consumes:
      - Protocol Request
    Produces:
      - SIR

  SIR:
    Purpose: Represent semantic intent
    DependsOn: []
    Consumes:
      - Protocol Runtime
    Produces:
      - Semantic Representation

  Capability:
    Purpose: Describe supported semantic behavior
    DependsOn:
      - SIR
    Produces:
      - Capability Metadata

  Planning:
    Purpose: Produce execution plan
    DependsOn:
      - SIR
      - Capability
    Consumes:
      - SIR
      - Capability Metadata
    Produces:
      - Execution Plan

  Data Kernel:
    Purpose: Coordinate execution
    DependsOn:
      - Planning
      - Engine Runtime
    Consumes:
      - Execution Plan
    Produces:
      - Engine Request

  Engine Runtime:
    Purpose: Translate engine request into engine operation
    DependsOn:
      - Engine
    Consumes:
      - Engine Request
    Produces:
      - Engine Operation

  Engine:
    Purpose: Execute engine operation
    DependsOn: []
    Consumes:
      - Engine Operation
    Produces:
      - Result

Constraints:

  Forbidden:
    - Application -> Engine
    - Engine -> Application
    - Engine -> Protocol
    - Protocol -> Engine
    - Specification -> Runtime

Flow:
  - Application
  - Protocol
  - Protocol Runtime
  - SIR
  - Planning
  - Data Kernel
  - Engine Runtime
  - Engine
