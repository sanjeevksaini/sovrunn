# Go Observability Standard

This is the source-of-truth standard for observability in Sovrunn Go code.

All Go code must follow this standard unless a feature explicitly documents why observability is not applicable.

## Required Context Fields

Where applicable, logs, metrics, traces, decisions, and audit events must carry stable correlation fields:

- `request_id`
- `operation_id`
- `organization_id`
- `tenant_id`
- `project_id`
- `service_instance_id`
- `decision_id`
- `feature_id`

Not all fields apply to every package. Missing fields must be intentional, not accidental.

## Logging Rules

- Use structured logging.
- Avoid `fmt.Println`, `log.Println`, and ad hoc text logs in production paths.
- Do not log secrets, credentials, tokens, private keys, passwords, connection strings, or raw sensitive payloads.
- Use stable reason codes for expected errors.
- Log lifecycle transitions for operations and service instances.
- Log external adapter calls at boundary level without leaking sensitive inputs.

## Audit vs Logs

Application logs are diagnostic.

Audit events are evidence.

Security, governance, placement, provisioning, binding, credential, policy, and lifecycle decisions must be represented as `AuditEvent` records where required. They must not be treated as ordinary logs only.

## Metrics

Metrics should be added where they provide operational value, especially for:

- API request counts and latency,
- decision counts and outcomes,
- operation duration and failures,
- plugin execution duration and failures,
- adapter call duration and failures,
- resource reconciliation outcomes.

Metric names should be stable and low-cardinality.

## Tracing

Trace/span boundaries should be considered for:

- API request handling,
- policy evaluation,
- placement decision,
- operation execution,
- plugin execution,
- external provider/substrate calls.

OpenTelemetry is the preferred direction.

## Error Reporting

Errors should include:

- stable reason code,
- safe human-readable message,
- correlation IDs where available,
- wrapped cause where useful.

Do not return internal sensitive details to users.

## Review Gate

A feature touching Go code should describe:

- logs added or preserved,
- metrics impact,
- trace impact,
- audit behavior,
- sensitive fields intentionally not logged,
- correlation fields propagated.
