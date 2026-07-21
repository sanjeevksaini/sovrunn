---
doc_type: architecture
title: Decision and Audit Standard
status: draft
phase: 2
ai_load_priority: always
ai_summary: Common decision and audit event schema for governance, security, placement, scaling, traffic, and future decisions.
---

# Decision and Audit Standard

## Purpose

Sovrunn must make important platform decisions explainable, auditable, and AI-readable.

## Decision Base Shape

```yaml
decision: ALLOWED # ALLOWED | DENIED | REQUIRES_APPROVAL
reasonCodes: []
humanReadableReasons: []
selectedTarget: null
rejectedAlternatives: []
policyReferences: []
riskLevel: low
suggestedActions: []
auditEventRef: null
```

## Decision Types

- GovernanceDecision
- SecurityDecision
- PlacementDecision
- DataMovementDecision
- ScalingDecision
- TrafficDecision
- FailoverDecision
- ComplianceDecision

## AuditEvent Base Shape

```yaml
eventType: decision.created
actorRef: null
subjectRef: null
resourceRef: null
decisionRef: null
operationRef: null
outcome: success
timestamp: null
reason: null
metadata: {}
```

## Rules

- Every meaningful platform change must create or link to an Operation.
- Every meaningful decision must create or link to an AuditEvent.
- Decision reasons must be structured enough for AI explanation.
- Denials must include actionable reasons and alternatives when possible.
