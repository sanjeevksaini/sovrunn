# Ownership

Document:
  ID: ownership
  Version: 1.1
  Status: Stable

Purpose:
  - Ownership domains
  - Ownership boundaries
  - Authority hierarchy

Rules:
  - One owner per concept
  - Owner defines
  - Non owner references
  - Reference never redefine
  - Use glossary canonical terms

Domains:

  Foundation:
    Owns:
      - Philosophy
      - Vision
      - Constitution
      - Glossary
      - Ontology
      - Ownership
      - ADS
      - Style

  Protocol Specification:
    Owns:
      - Protocol
      - Request
      - Response

  SIR Specification:
    Owns:
      - SIR
      - Resource
      - Operation
      - Expression
      - Relationship
      - Constraint
      - Metadata

  Capability Specification:
    Owns:
      - Capability
      - Capability Model
      - Capability Identifiers

  Reuse Specification:
    Owns:
      - Adopted Standards
      - Adopted Algorithms
      - Adopted Architecture Patterns
      - Adopted Libraries
      - Reuse Authority
      - Standard Mapping

  Runtime Architecture:
    Owns:
      - Runtime
      - Protocol Runtime
      - SIR Runtime
      - Planning
      - Execution Plan
      - Execution Context
      - Data Kernel
      - Engine Runtime
      - Capability Registry
      - Session Runtime
      - Transaction Runtime
      - Result Model
      - Error Model

  SDE Control Plane Architecture:
    Owns:
      - SDE Control Plane
      - Control Plane Foundation
      - Foundation Service
      - Foundation Provider
      - Management Plane
      - Management Plane Host

  Core Control Plane:
    Owns:
      - Runtime Registry
      - Plugin Registry
      - Engine Registry
      - Capability Governance
      - Deployment Governance
      - Runtime Governance
      - Extension Governance

  Datastore Management Plane:
    Owns:
      - dstoreOps
      - Datastore Registry
      - Datastore Operator Registry
      - Infrastructure Provider Registry
      - Datastore Lifecycle Management
      - Datastore Operator Plugin Coordination
      - Infrastructure Provider Coordination

  SDE Data Plane Architecture:
    Owns:
      - SDE Data Plane
      - Request Execution Flow
      - Protocol Execution
      - Planning Execution
      - Kernel Execution
      - Engine Execution
      - Result Propagation
      - Error Propagation

  Plugin Architecture:
    Owns:
      - Protocol Plugin
      - Engine Plugin
      - Plugin Runtime
      - Plugin Boundaries

  Datastore Operator Plugin:
    Owns:
      - Datastore-specific lifecycle integration
      - Lifecycle operation mapping
      - Lifecycle error mapping
      - Lifecycle state reporting

  Infrastructure Provider:
    Owns:
      - Infrastructure environment integration
      - Infrastructure operation mapping
      - Infrastructure provider error mapping
      - Infrastructure provider state reporting

  Engine Plugin:
    Owns:
      - Downstream execution integration
      - Native result mapping
      - Native error mapping
      - Capability Manifest

  Implementation:
    Owns:
      - Source Code
      - Tests
      - Build
      - Deployment

  External:
    Owns:
      - Application
      - Datastore
      - Downstream Datastore
      - Downstream Engine
      - Datastore SDE Data Plane
      - Standard

Authority:

  Highest:
    - Constitution

  Concept:
    - Ontology

  Vocabulary:
    - Glossary

  Ownership:
    - Ownership

  Semantics:
    - SIR Specification

  Capability:
    - Capability Specification

  Runtime:
    - Runtime Architecture

  Management:
    - SDE Control Plane Architecture

  Execution:
    - SDE Data Plane Architecture

  Code:
    - Implementation
