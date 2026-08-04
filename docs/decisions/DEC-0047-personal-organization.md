---
doc_type: decision_record
title: Personal Organization
status: accepted
phase: 2R
ai_load_priority: important
---

# Personal Organization

## Status

Accepted

## Context

New users need a simple onboarding path without requiring enterprise organization setup.

## Decision

Personal onboarding creates a Personal Organization, default Tenant, and default Project.

## Consequences

Simple developer experience for individual users. Enterprise onboarding uses explicit Organization creation. Both paths share the same governance model.

## Migration

No existing personal organizations need migration; establishes the onboarding pattern.

## Conformance Evidence

Simple personal and enterprise onboarding tests both pass.

## Source ADH Reference

ADH-2026-030
