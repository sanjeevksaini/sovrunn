---
doc_type: execution_strategy
title: Phase 2 Execution Strategy
status: active
ai_load_priority: critical
---

# Phase 2 Execution Strategy

## Purpose

This document defines how Sovrunn Phase 2 must be designed and executed.

Phase 2 introduces the reusable governance, policy, decision, adapter, provider-neutral, plugin-taxonomy, and placement foundation for Sovrunn.

The goal is to avoid two risks:

1. Designing all Phase 2 features in excessive detail before implementation.
2. Implementing Phase 2 features independently without a shared architecture spine.

The approved approach is:

```text
Phase 2 architecture spine first.
Then strict sequential feature execution.
```

## Core Execution Principle

Phase 2 must use a hybrid architecture approach.

```text
High-level architecture for all Phase 2: yes.
Detailed architecture per feature: only when starting that feature.
Code and tests: only for the current feature.
```

Do not fully design all Phase 2 features upfront.

Do not implement multiple Phase 2 features at the same time.

Do not move to the next feature until the current feature passes its feature gate.

## Approved Workflow

Each Phase 2 feature must follow this sequence:

```text
1. Architecture discussion
2. Architecture Decision Handoff
3. Human approval
4. Kiro requirements, design, and tasks
5. Cursor implementation
6. Tests and validation
7. Feature gate
8. Commit / PR
9. Move to next feature
```

## Phase 2 Architecture Spine

Before detailed execution, Phase 2 requires a shared architecture spine.

The spine must define:

- Phase 2 dependency graph
- shared terminology
- shared object model
- common API and resource conventions
- decision and audit pattern
- provider-neutral core boundaries
- adapter, plugin, and core separation
- strict Phase 2 non-goals
- expected end-state of the FEATURE-0026 integration demo

The spine is not a full detailed design of every Phase 2 feature.

It is a controlling map that prevents feature-level architecture drift.

## Dependency Rule

Phase 2 features are dependency ordered.

Later features must not redefine concepts already established by earlier features.

If a later feature requires changing an earlier decision, it must use the Architecture Decision Handoff process.

## Phase 2 Mini-Waves

### Wave A: Standards and Guardrails

```text
FEATURE-0011: Reuse Assessment Standard
FEATURE-0012: API, Resource Naming, Status, and Validation Standard
FEATURE-0013: Decision Object and AuditEvent Standard
```

These features define the standards used by all later Phase 2 work.

### Wave B: Provider-Neutral Substrate Model

```text
FEATURE-0014: Provider-Neutral Resource Model
FEATURE-0015: ResourcePool and ProviderCapability Model
FEATURE-0016: Adapter Boundary Foundation
```

These features define how Sovrunn avoids hardcoding Kubernetes, IaaS, or provider-specific assumptions into core.

### Wave C: Policy and Governance Foundation

```text
FEATURE-0017: Policy Evaluation Abstraction
FEATURE-0018: GovernanceProfile and SecurityProfile Foundation
FEATURE-0019: DataPlacementPolicy and CostGuardrail Minimal Foundation
FEATURE-0020: ProfileAssignment and EffectivePolicyContext
FEATURE-0021: Minimal ServiceEntitlement and Quota Placeholder
```

These features define whether a request is allowed and what effective policy context applies.

### Wave D: Runtime Intent and Placement Decisions

```text
FEATURE-0022: ServiceRuntimeProfile Foundation
FEATURE-0023: PlacementRequest and PlacementDecision v0
FEATURE-0024: Plugin Taxonomy Foundation
FEATURE-0025: AI-Readable Decision Context
FEATURE-0026: Phase 2 Integration Demo
```

These features connect service intent, policy, placement, plugin taxonomy, audit, and explanation.

## Feature Execution Rules

For every Phase 2 feature:

- Start from the approved Phase 2 architecture spine.
- Keep feature scope narrow.
- Produce an Architecture Decision Handoff before Kiro updates.
- Use Kiro for requirements, design, and tasks.
- Use Cursor only after Kiro tasks are approved.
- Require tests before feature completion.
- Run the feature gate before merge.
- Do not add later-phase runtime behavior early.

## Phase 2 Non-Goals

Phase 2 must not implement:

- real PostgreSQL provisioning
- real plugin execution
- full OPA or Cedar integration
- full Keycloak, Vault, or Temporal integration
- Kubernetes-only assumptions inside Sovrunn core
- production persistence
- billing
- autonomous AI operations
- failover or disaster recovery execution
- production compliance engine

## Accuracy, Speed, and Correctness Guidance

The approved optimization model is:

```text
Accuracy:
  Use Phase 2 architecture spine first.

Speed:
  Keep each feature as a small vertical slice.

Correctness:
  Complete architecture, spec, implementation, tests, and gate for one feature before starting the next.
```

The combined strategy is:

```text
Spine first.
Then sequential feature execution.
```

## ChatGPT Project Usage

The ChatGPT Project should be used for architecture discussion and decision handoff preparation.

The Project must remember:

```text
Do not fully design all Phase 2 features upfront.
Do not code multiple Phase 2 features together.
First define the Phase 2 architecture spine.
Then execute FEATURE-0011 through FEATURE-0026 sequentially.
Each feature must complete architecture, Kiro spec, implementation, tests, and feature gate before the next feature begins.
```

## Kiro Usage

Kiro must use this document when generating or updating Phase 2 specs.

Kiro should not create requirements, design, or tasks that violate the Phase 2 execution strategy.

Kiro must preserve Phase 2 dependency order unless an approved Architecture Decision Handoff changes the sequence.

## Cursor Usage

Cursor must implement only the current approved feature.

Cursor must not pre-implement later Phase 2 features unless explicitly approved by Kiro tasks and human review.

Cursor must not introduce runtime provisioning, plugin execution, or provider-specific shortcuts during Phase 2 unless the current feature explicitly allows it.

## Feature Gate Expectation

The feature gate must enforce strict AOS validation for FEATURE-0011 and later.

Phase 1 features FEATURE-0001 through FEATURE-0010 remain legacy baseline features and are exempt from strict Phase 2 AOS section checks.

For FEATURE-0011 and later, the feature gate should validate:

- Reuse Assessment
- Acceptance Criteria
- Architecture Drift checks
- Observability considerations
- Security considerations
- Non-goals
- generated artifact hygiene
- tests
- lint
- security scan
- phase scope boundary

## Architecture Decision Handoff Expectation

Any cross-feature architecture change must produce an Architecture Decision Handoff.

The handoff must describe:

- decision title
- classification
- approved baseline
- proposed change
- rationale
- reuse-before-build assessment
- phase impact
- conflict check
- required action
- impacted files
- impacted features
- Kiro acceptance criteria
- explicit Kiro instructions
- human approval status

## Completion Rule

A Phase 2 feature is complete only when:

```text
Architecture decision is approved.
Kiro requirements/design/tasks are complete.
Implementation is complete.
Tests pass.
Feature gate passes.
Review is complete.
Commit or PR is merged.
```

Only then should the next feature begin.
