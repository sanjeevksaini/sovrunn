---
doc_type: architecture
title: Observability and Audit Baseline
status: draft
phase: 1
ai_load_priority: phase
ai_summary: Defines baseline logging, audit, metrics, trace, and operation correlation fields for Sovrunn Phase 1.
---

# Observability and Audit Baseline

## 1. Purpose

Sovrunn must be observable and auditable from the beginning.

Phase 1 does not need a full OpenTelemetry collector, metrics backend, log backend, or audit database. However, the field model and correlation model must be defined now.

The principle is:

```text
Every platform action must be explainable, traceable, and auditable.
```

## 2. Required Correlation Fields

Every mutating request should have:

```text
request_id
operation_id
trace_id
actor
organization
organizationUnit
tenant
project
resource_kind
resource_name
action
status
error_code
latency_ms
```

In Phase 1, some fields may be empty depending on resource scope.

## 3. Request ID

Clients may pass:

```text
X-Sovrunn-Request-ID
```

If absent, the API server generates one.

The response should include:

```text
X-Sovrunn-Request-ID
```

## 4. Operation ID

Mutating requests should eventually return:

```text
X-Sovrunn-Operation-ID
```

For Phase 1:

```text
Before FEATURE-0005:
  operation_id may be absent.

After FEATURE-0005:
  mutating requests must create and return an operation ID.
```

## 5. Actor

Phase 1 may use:

```text
system
anonymous-dev
```

Future phases should derive actor from identity provider/OIDC.

The API may accept a temporary development header:

```text
X-Sovrunn-Actor
```

This is not a security boundary.

## 6. Structured Logging

Minimum fields:

```text
level
timestamp
message
request_id
operation_id
resource_kind
resource_name
action
status
error_code
latency_ms
```

Do not log:

```text
passwords
tokens
secret values
private keys
credential payloads
authorization headers
```

## 7. Audit Events

An audit event records who did what, to which resource, under which organizational context, and with what result.

Audit event fields:

```text
audit_event_id
timestamp
actor
action
resource_kind
resource_name
organization
organizationUnit
tenant
project
request_id
operation_id
status
error_code
message
```

Phase 1 audit events may be in-memory or logged. Future audit events should be durable.

## 8. Metrics

Future metrics:

```text
sovrunn_api_requests_total
sovrunn_api_request_duration_ms
sovrunn_operations_total
sovrunn_operations_failed_total
sovrunn_registry_resources_total
sovrunn_validation_failures_total
```

Use low-cardinality labels:

```text
resource_kind
action
status
error_code
```

Avoid high-cardinality labels:

```text
resource_name
request_id
operation_id
actor
```

## 9. Tracing

Future OpenTelemetry spans should follow this structure:

```text
HTTP request span
  -> validation span
  -> registry span
  -> operation creation span
  -> audit span
  -> plugin span, future
```

Phase 1 may only generate trace IDs and structured logs.

## 10. Standard Error Codes

```text
VALIDATION_FAILED
RESOURCE_NOT_FOUND
RESOURCE_ALREADY_EXISTS
DELETE_BLOCKED
REFERENCE_INVALID
POLICY_DENIED
UNAUTHORIZED
FORBIDDEN
INTERNAL_ERROR
```

## 11. Operation and Audit Relationship

`Operation` is a lifecycle record.

`AuditEvent` is an accountability record.

Both are needed.

## 12. Phase 1 Non-Goals

Do not implement yet:

```text
full OpenTelemetry collector deployment
Prometheus endpoint if not needed
Loki/OpenSearch integration
distributed tracing backend
audit database
SIEM export
compliance report generation
```

## 13. Acceptance Criteria

Phase 1 observability baseline is satisfied when:

```text
request IDs are generated or propagated
mutating requests create operation records after FEATURE-0005
structured logs include action/resource/status/latency
standard error codes are used
secrets are not logged
resource context is available for audit
```

## 14. Final Principle

If an administrator asks what happened, who did it, where it happened, and why it failed, Sovrunn must be able to answer.
