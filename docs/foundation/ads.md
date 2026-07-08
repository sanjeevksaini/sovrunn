# AI Documentation Specification

Document:
  ID: ads
  Version: 1.0
  Status: Stable

Purpose:
  - AI documentation contract
  - Deterministic document structure
  - Minimum tokens without semantic loss

Rules:
  - AI first
  - Human friendly
  - Information dense
  - One responsibility per document
  - One decision per RFC
  - Single source of truth
  - Reference never duplicate
  - Stable vocabulary
  - Canonical abbreviations
  - Markdown with structured fields
  - YAML for structured knowledge
  - No decorative prose
  - No unnecessary punctuation

DocumentSchema:

  Required:
    - Document
    - Purpose

  Document:
    Required:
      - ID
      - Version
      - Status

  Purpose:
    Required:
      - List

Structure:
  - Metadata first
  - Purpose second
  - Rules third when required
  - Content after rules
  - References last when required

Constraints:

  MUST:
    - preserve semantics
    - use canonical terms
    - use structured fields
    - use atomic list items
    - define once
    - reference canonical owner
    - minimize tokens

  MUST NOT:
    - duplicate definitions
    - use synonyms
    - mix responsibilities
    - add examples unless required
    - add rationale unless required
    - use decorative formatting
    - use filler text

Formats:

  Markdown:
    Use For:
      - prose-light specifications
      - architecture summaries
      - RFCs

  YAML:
    Use For:
      - ontology
      - ownership
      - runtime model
      - dependency graph
      - capability manifests
      - structured specifications

QualityGate:

  Check:
    - one responsibility
    - deterministic parsing
    - canonical vocabulary
    - no duplication
    - no unnecessary tokens
    - no undefined terms
    - no ownership conflict
