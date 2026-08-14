---
doc_type: decision_record
title: Canonical Bootstrap Replaces Alpha Runtime Migration
status: Accepted
supersedes: DEC-0058
phase: 2R
ai_load_priority: important
---

# Canonical Bootstrap Replaces Alpha Runtime Migration

## Status

Accepted

Supersedes: DEC-0058

## Context

DEC-0058 and ADH-2026-041 defined alpha migration as a signed cutover with no dual write/authority and immutable history preservation, on the premise that FEATURE-0001 through FEATURE-0014 constitute live control-plane data requiring conversion. That premise no longer holds: Sovrunn has no live control plane, customer data, production API estate, or persisted alpha state to convert. FEATURE-0001 through FEATURE-0014 are retained repository assets (implementation history), not a live deployment with state that must be migrated.

Building runtime migration machinery (`CanonicalMigrationPlan`, `CanonicalMigrationRecord`, a migration controller, or a cutover state machine) before a control plane exists would add a product surface, lifecycle, security boundary, controller authority, audit stream, and conformance burden for data that does not exist.

## Decision

1. The first executable Sovrunn control-plane release begins directly with canonical resource contracts. It does not implement alpha import, fixture conversion as a product capability, migration plans/records, a migration controller, migration milestones, write-freeze/cutover behavior, backup/restore migration evidence, or global migration completion.
2. FEATURE-0001 through FEATURE-0014 remain retained repository assets. Their code, documents, tests, standards, and compatible generic infrastructure are assessed for reuse or extension under FEATURE-0011; they do not constitute live control-plane data requiring conversion.
3. Alpha kinds, scope vocabulary, and desired-state endpoints are not active compatibility surfaces and must not be reintroduced as runtime APIs, writers, projections, or persistence contracts.
4. FEATURE-0015 owns direct creation of its canonical cloud-model resources (CloudPlatform, CloudProvider, CloudProviderParticipation, HostingLocation, Datacenter, FaultDomain, InfrastructureStack) with no migration authority.
5. `ExecutionTarget` moves in its entirety to FEATURE-0016. FEATURE-0015 ends at `InfrastructureStack`.
6. The one-authority principle (previously satisfied by a no-dual-write cutover) is now satisfied by beginning with canonical desired-state contracts only; there is no alpha authority to coexist with.
7. If historical fixture comparison remains useful during development, it is a bounded developer/test utility only — not a platform resource, runtime protocol, or migration authority.

## Consequences

`CanonicalMigrationPlan`, `CanonicalMigrationRecord`, the migration-controller and approved-migration-plan-publisher writer roles, the VS0-STATE-011 migration milestone state machine, and the migration failure mappings (VS0-MIG-F01..F03 and their VS0-CF-MIG*/VS0-CF-MIGF* conformance cases) are removed from active Phase 2R authorities. FEATURE-0015 requirements, design, and tasks must not describe alpha conversion, dual-write prevention, or cutover evidence as product behavior.

ADH-2026-041 (migration runtime semantics) is superseded in its entirety. ADH-2026-043 decisions 1–3 and 6 (the portions governing migration plan/record behavior) are superseded; ADH-2026-043's safe-denial, scope-reference integrity, audit reuse, and feature-local traceability decisions remain intact and unaffected. ADH-2026-044 (migration run cardinality) is superseded in its entirety, as it exists solely to govern migration-record cardinality.

## Migration

There is no migration to execute. This decision retires the migration model itself rather than replacing one migration mechanism with another. DEC-0058 is fully superseded.

## Conformance Evidence

Active Phase 2R authorities (canonical data model, contract catalog, VS-000 contract registry/specification/traceability, FEATURE-0015 feature and architecture boundary, and the FEATURE-0015 architecture-readiness checker) contain no `CanonicalMigrationPlan`, `CanonicalMigrationRecord`, migration-controller, migration milestone state machine, or migration failure mapping as active behavior. FEATURE-0015 ends at `InfrastructureStack` with no `ExecutionTarget` contract.

## Source ADH Reference

ADH-2026-045

