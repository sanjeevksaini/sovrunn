---
doc_type: decision_record
title: Extensibility by Data
status: accepted
phase: 2R
ai_load_priority: important
---

# Extensibility by Data

## Status

Accepted

## Context

Adding new countries, sectors, services, and backends could require core code changes or conditional logic.

## Decision

New countries, sectors, services, and backends enter through governed data, contracts, plugins, and adapters — not core conditionals.

## Consequences

Core remains stable. New capabilities are data-driven extensions. Plugin/adapter boundaries are the extensibility surface.

## Migration

Existing extension points are already adapter-based; this strengthens the commitment.

## Conformance Evidence

Synthetic service added without core schema change.

## Source ADH Reference

ADH-2026-031
