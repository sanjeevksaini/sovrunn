---
doc_type: decision_record
title: Alpha Migration No Dual Authority
status: Accepted
phase: 2R
ai_load_priority: important
---

# Alpha Migration No Dual Authority

## Status

Accepted

## Context

Migrating from the alpha model to the canonical model could create a period of dual authority where both old and new terminology coexist as active.

## Decision

Alpha migration is one signed cutover with no dual write/authority and immutable history preservation.

## Consequences

Old and new writers never coexist. Migration is atomic. Completed history remains immutable. Rollback is full revert, not selective restoration.

## Migration

The canonical migration plan defines write freeze, verified backup/restore, deterministic mapping, no dual authority.

## Conformance Evidence

Old and new writers never coexist. Deterministic dry-run conversion passes. Zero ambiguous authorities after cutover.

## Source ADH Reference

ADH-2026-041
