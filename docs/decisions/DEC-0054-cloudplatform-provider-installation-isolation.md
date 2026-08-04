---
doc_type: decision_record
title: CloudPlatform Provider Installation Isolation
status: Accepted
phase: 2R
ai_load_priority: important
---

# CloudPlatform Provider Installation Isolation

## Status

Accepted

## Context

CloudPlatform ownership, CloudProvider participation, and installation isolation could be conflated as one boundary.

## Decision

CloudPlatform ownership, CloudProvider participation, and one-provider-per-installation isolation are distinct boundaries.

## Consequences

A CloudPlatform may be owned without being operated. A CloudProvider may participate without owning. Each installation has exactly one provider for credential isolation.

## Migration

Existing Provider records split into appropriate CloudPlatform ownership and CloudProvider participation.

## Conformance Evidence

Provider credential isolation per installation is enforced. Ownership and participation are independently testable.

## Source ADH Reference

ADH-2026-037
