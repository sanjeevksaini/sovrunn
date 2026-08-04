---
doc_type: decision_record
title: ServiceOffering Replaces ServiceClass
status: Accepted
phase: 2R
ai_load_priority: important
---

# ServiceOffering Replaces ServiceClass

## Status

Accepted

## Context

Global ServiceClass conflated reusable type definitions with platform-scoped product offerings. This prevented multi-platform catalog independence.

## Decision

ServiceOffering is the CloudPlatform product and references a reusable ServiceTypeDefinition. Global ServiceClass is migration input only and has no active canonical authority.

## Consequences

ServiceTypeDefinition is reusable across platforms. ServiceOffering is platform-scoped. ServiceClass references in completed features are historical.

## Migration

Old ServiceClass references map deterministically to ServiceTypeDefinition + CloudPlatform-scoped ServiceOffering.

## Conformance Evidence

ServiceClass has no active canonical authority. One mapping to ServiceTypeDefinition and ServiceOffering exists.

## Source ADH Reference

ADH-2026-032
