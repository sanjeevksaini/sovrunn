---
doc_type: rfc
title: RFC-0022 API Resource Standard
status: draft
phase: 2
ai_load_priority: high
ai_summary: RFC for Sovrunn resource shape, status, conditions, references, and API boundary classification.
---

# RFC-0022: API Resource Standard

See `docs/architecture/api-resource-standard.md` for the current standard.

## Decision

Phase 2 resources must use consistent `apiVersion`, `kind`, `metadata`, `spec`, `status`, and `conditions` patterns.

Each resource must be classified as customer-facing, provider/MSP-facing, internal engine-facing, or plugin-facing.
