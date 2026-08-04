---
doc_type: decision_record
title: Infrastructure Maintenance Authority
status: Accepted
phase: 2R
ai_load_priority: important
---

# Infrastructure Maintenance Authority

## Status

Accepted

## Context

External maintenance (provider-initiated changes, OS patches, hardware lifecycle) could bypass platform governance.

## Decision

External maintenance has explicit notice authority, target lifecycle ownership, epochs, fences, and mandatory requalification.

## Consequences

Provider maintenance cannot silently change execution targets. Fencing prevents execution on unqualified targets. Requalification is mandatory after maintenance.

## Migration

No existing maintenance authority exists; establishes the canonical pattern.

## Conformance Evidence

Maintenance epoch, fencing, and requalification contracts are modeled in Phase 2R.

## Source ADH Reference

ADH-2026-040
