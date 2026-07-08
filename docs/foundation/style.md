# Style

Document:
  ID: style
  Version: 1.0
  Status: Stable

Purpose:
  - Documentation writing rules

Principles:
  - AI First
  - Human Friendly
  - Information Dense
  - Minimum Token Consumption
  - Single Source Of Truth
  - One Responsibility Per Document
  - One Decision Per RFC
  - Stable Vocabulary
  - Canonical Abbreviations
  - Reference Instead Of Duplicate
  - Architecture Before Implementation
  - Reuse Before Build
  - Vertical Slice First

Writing:

  MUST:
    - Use Markdown
    - Use YAML for structured knowledge
    - Use structured fields
    - Use canonical vocabulary
    - Use canonical abbreviations
    - Use atomic statements
    - Use consistent ordering
    - Keep documents self contained
    - Reference canonical documents

  MUST NOT:
    - Duplicate knowledge
    - Introduce synonyms
    - Mix responsibilities
    - Repeat definitions
    - Repeat examples
    - Add decorative prose
    - Add unnecessary punctuation
    - Add unnecessary whitespace

Vocabulary:

  Case:
    Documents: Pascal Case
    Concepts: Pascal Case
    Abbreviations: Upper Case
    Files: kebab-case
    YAML Fields: Pascal Case

  Canonical:
    Semantic Intermediate Representation: SIR
    Request For Comments: RFC
    AI Documentation Specification: ADS
    Minimum Viable Product: MVP

Lists:

  Rules:
    - Short
    - Atomic
    - Deterministic
    - Ordered
    - Parallel

Headings:

  Order:
    - Document
    - Purpose
    - Rules
    - Content

References:

  MUST:
    - Reference canonical owner
    - Reference adopted standards
    - Reference instead of duplicate

Evolution:

  MUST:
    - Preserve compatibility
    - Preserve vocabulary
    - Preserve document responsibility
