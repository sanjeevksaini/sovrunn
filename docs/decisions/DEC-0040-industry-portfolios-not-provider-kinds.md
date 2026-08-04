---
doc_type: decision_record
title: Industry Portfolios Not Provider Kinds
status: Accepted
phase: 2R
ai_load_priority: important
---

# Industry Portfolios Not Provider Kinds

## Status

Accepted

## Context

Industry and domain clouds could be modeled as new core provider kinds, causing type explosion.

## Decision

Industry and domain clouds are versioned ServicePortfolios, not new core provider kinds.

## Consequences

Adding a new industry does not require core schema changes. ServicePortfolio is a data-driven extension.

## Migration

No existing industry provider kinds exist; this is a preventive boundary.

## Conformance Evidence

Synthetic service added without core schema change in Phase 2R tests.

## Source ADH Reference

ADH-2026-023
