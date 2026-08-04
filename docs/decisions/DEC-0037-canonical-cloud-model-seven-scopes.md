---
doc_type: decision_record
title: Canonical Cloud Model Seven Scopes
status: Accepted
phase: 2R
ai_load_priority: important
---

# Canonical Cloud Model Seven Scopes

## Status

Accepted

## Context

The alpha model used ambiguous generic Provider as combined product-owner and infrastructure-operator. This merged two distinct responsibilities and prevented correct classification. The existing six governance scopes did not model CloudPlatform ownership.

## Decision

Retire ambiguous Provider by classifying records into CloudPlatform ownership and/or CloudProvider operation. Introduce seven canonical scope kinds: Platform, Organization, OrganizationUnit, Tenant, Project, CloudPlatform, CloudProvider.

## Consequences

All scope vocabulary must use seven kinds. Generic Provider is migration input only. Active authorities must not use the six-scope vocabulary after cutover.

## Migration

Existing Provider fixtures convert to CloudPlatform or CloudProvider based on their actual role. FEATURE-0014 Provider registry remains completed history.

## Conformance Evidence

Drift checks verify no active authority uses six-scope vocabulary or generic Provider as combined concept.

## Source ADH Reference

ADH-2026-020
