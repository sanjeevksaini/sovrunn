---
doc_type: rfc
title: RFC-0029 Governed PostgreSQL MVP
status: draft
phase: 3
ai_load_priority: high
ai_summary: RFC for first customer-testable Sovrunn MVP.
---

# RFC-0029: Governed PostgreSQL MVP

See `docs/mvp/MVP_001_GOVERNED_POSTGRESQL_PAAS.md`.

## Decision

The first MVP is governed PostgreSQL PaaS placement and provisioning on one substrate. Sovrunn must reuse a mature PostgreSQL operator or Helm chart rather than building PostgreSQL HA/lifecycle internals.

## Phase 2R Migration Note

**ADH-2026-042:** The MVP remains governed PostgreSQL PaaS placement and provisioning on one substrate. Phase 2R establishes the canonical service product, runtime, and requirement contracts (FEATURE-0022) and the PostgreSQL reference flow (`docs/architecture/reference-flows/postgresql-end-to-end.md`). Phase 3 executes the MVP using the canonical contracts rather than the old ServiceClass/ResourcePool/ProviderCapability topology. The reuse-first principle (mature PostgreSQL operator/Helm) is unchanged.
