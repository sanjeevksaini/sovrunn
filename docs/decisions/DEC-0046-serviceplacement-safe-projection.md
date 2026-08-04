---
doc_type: decision_record
title: ServicePlacement Safe Projection
status: accepted
phase: 2R
ai_load_priority: important
---

# ServicePlacement Safe Projection

## Status

Accepted

## Context

Customers should know where their service runs but not have access to canonical infrastructure topology or provider credentials.

## Decision

ServicePlacement is a safe immutable customer projection, not canonical topology access. It contains no provider-native objects, target credentials, raw secrets, or protected handles.

## Consequences

Customer APIs expose ServicePlacement with redacted/projected information. Internal topology remains internal.

## Migration

No existing ServicePlacement exists; this establishes the customer-safe boundary.

## Conformance Evidence

Customer APIs contain no provider-native object, target credential, raw secret, or protected handle.

## Source ADH Reference

ADH-2026-029
