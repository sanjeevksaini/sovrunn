---
doc_type: decision_record
title: ServiceBinding SecretRef Only
status: Accepted
phase: 2R
ai_load_priority: important
---

# ServiceBinding SecretRef Only

## Status

Accepted

## Context

ServiceBinding could expose raw credentials, full connection strings, or internal topology details to consumers.

## Decision

ServiceBinding is the per-consumer, separately revocable, SecretRef-only service-access boundary. No raw credentials appear in binding resources.

## Consequences

Each consumer gets an independently revocable binding. SecretRef is the only credential exposure mechanism. Credential rotation does not require binding recreation.

## Migration

Existing Phase 1 ServiceBinding resources remain historical. Phase 2R contracts enforce SecretRef-only.

## Conformance Evidence

No binding resource contains raw secrets or protected handles.

## Source ADH Reference

ADH-2026-034
