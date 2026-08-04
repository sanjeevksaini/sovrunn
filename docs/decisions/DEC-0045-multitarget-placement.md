---
doc_type: decision_record
title: Multi-Target Placement
status: accepted
phase: 2R
ai_load_priority: important
---

# Multi-Target Placement

## Status

Accepted

## Context

A service may need to span multiple execution targets (e.g., primary + standby) with explicit cross-target constraints.

## Decision

One service may be realized across an atomic set of execution targets with explicit cross-target constraints.

## Consequences

ServiceDeploymentPlan may reference multiple ExecutionTargets. Cross-target constraints (latency, replication) are explicit. Phase 2R models the contract; Phase 3 executes it.

## Migration

No existing multi-target deployments need migration; this establishes the forward contract.

## Conformance Evidence

Phase 2R models multi-target placement in synthetic scenarios.

## Source ADH Reference

ADH-2026-028
