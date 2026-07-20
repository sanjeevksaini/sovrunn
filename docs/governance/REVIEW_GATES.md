# Review Gates

This document defines required review gates for Sovrunn architecture and development.

## Feature Gate

Before a feature is complete, run:

```bash
make ff-feature-gate FEATURE=<FEATURE-ID>
```

A feature is not complete until this gate passes.

## Architecture Drift Gate

Reject the feature if it:

- bypasses reuse-before-build,
- violates phase scope,
- hardcodes provider-specific behavior in core,
- hardcodes Kubernetes-only assumptions in core,
- puts PostgreSQL lifecycle logic into core placement,
- introduces custom policy engine logic,
- stores raw secrets,
- misses required audit events,
- lacks observability behavior,
- exposes low-level IaaS complexity to customer APIs.

## Phase Gate

At the end of each phase:

- verify feature traceability matrix,
- review accepted and pending decisions,
- close or carry open questions,
- update current architecture baseline,
- revalidate roadmap placeholders,
- create next phase context.

## Monthly Architecture Review

Monthly review should assess:

- architecture drift,
- stale decisions,
- reuse-before-build adherence,
- documentation accuracy,
- feature delivery alignment,
- open questions,
- next-month architecture focus.
