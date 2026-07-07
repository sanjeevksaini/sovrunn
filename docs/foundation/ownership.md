# Ownership

Document:
  ID: ownership
  Version: 1.0
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

Domains:

  Foundation:
    Owns:
      - Philosophy
      - Vision
      - Constitution
      - Glossary
      - Ontology
      - Ownership
      - Architectural Dependencies
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

  Architecture:
    Owns:
      - Protocol Runtime
      - Planning
      - Data Kernel
      - Engine Runtime
      - Execution Plan
      - Runtime Flow

  Implementation:
    Owns:
      - Source Code
      - Tests
      - Build
      - Deployment

  Adopted Standards:
    Owns:
      - External Standards
      - Reuse Authority
      - Standard Mapping

  External:
    Owns:
      - Application
      - Engine
      - Standard

Authority:

  Highest:
    - Constitution

  Concept:
    - Ontology

  Ownership:
    - Ownership

  Semantics:
    - SIR Specification

  Runtime:
    - Architecture

  Code:
    - Implementation
