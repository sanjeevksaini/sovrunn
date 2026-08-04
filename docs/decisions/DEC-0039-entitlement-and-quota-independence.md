---
doc_type: decision_record
title: Entitlement and Quota Independence
status: Accepted
phase: 2R
ai_load_priority: important
---

# Entitlement and Quota Independence

## Status

Accepted

## Context

Entitlement (what may be consumed) and quota (how much) were conflated.

## Decision

Entitlement answers what may be consumed; quota independently answers how much. Both may deny; denial by either is sufficient.

## Consequences

EntitlementPackage and QuotaPolicy are separate resources with independent evaluation.

## Migration

Existing entitlement placeholders split into distinct entitlement and quota authorities.

## Conformance Evidence

Tests prove entitlement and quota independently deny.

## Source ADH Reference

ADH-2026-022
