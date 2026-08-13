# Architecture Decision Handoff

## Metadata

- Handoff ID: ADH-2026-053
- Date: 2026-08-13
- Source discussion: FEATURE-0015 design-review recovery
- Related feature: FEATURE-0015
- Related phase: Phase 2R
- Author: Codex
- Human approver: Sanjeev Kumar
- Approval status: Approved

## Decision title

Use ServeMux-compatible CloudProviderParticipation item-action URIs

## Summary

This is a bounded public-URI correction for the eight existing
CloudProviderParticipation item actions. The approved 22-logical-path and
35-explicit-registration contract remains unchanged, but the former
`/{uid}:<action>` notation cannot be registered as a Go 1.22
`http.ServeMux` pattern with `{uid}` available through `r.PathValue("uid")`.
This handoff replaces only that URI notation with a literal `actions` path
segment. It does not change action semantics, lifecycle, authorization,
headers, bodies, audit behavior, idempotency behavior, resource shape, or
route counts.

## Classification

Correction — implementation-incompatible public URI notation for an already
approved eight-action surface.

## Existing approved baseline

- FEATURE-0015 owns exactly 22 logical endpoint paths and exactly 35 explicit
  Go 1.22 `http.ServeMux` method/path registrations.
- The surface contains seven collection paths, seven item paths, and eight
  existing-participation item-action paths.
- The eight actions are `accept`, `reject`, `withdraw`, `suspend`, `resume`,
  `request-release`, `accept-release`, and `decline-release`.
- Each existing-participation action uses `POST`, has an empty JSON body, and
  requires `If-Match` plus `Idempotency-Key`.
- Scheduler expiry is internal to the API server and is not a public route.
- ADH-2026-045 through ADH-2026-052 remain controlling and are not
  superseded.

The current `/{uid}:<action>` notation is not executable using the mandated
Go 1.22 `http.ServeMux` wildcard semantics: `{uid}` is not a complete path
segment and therefore cannot be read reliably with `r.PathValue("uid")`.

## Correction

### 1. ServeMux-compatible action URI shape

Replace the retired action URI form:

```text
POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}:<action>
```

with the following eight public item-action paths:

```text
POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/accept
POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/reject
POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/withdraw
POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/suspend
POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/resume
POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/request-release
POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/accept-release
POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}/actions/decline-release
```

In each pattern, `{uid}` is a complete Go 1.22 `http.ServeMux` wildcard
segment and is read with:

```go
r.PathValue("uid")
```

The retired `/{uid}:<action>` URI form is not an allowed FEATURE-0015 route
form.

### 2. Preserved behavior

This correction preserves all existing approved behavior:

- exactly 22 logical endpoint paths;
- exactly 35 explicit Go 1.22 `http.ServeMux` method/path registrations;
- eight public existing-participation `POST` item-action routes;
- the existing action names and transition semantics;
- empty JSON action body;
- `If-Match` and `Idempotency-Key` requirements;
- current authorization, safe-denial, audit, idempotency, and status behavior;
- deterministic internal scheduler expiry with no public expiry route.

No route is added or removed. The eight corrected paths replace the same eight
logical action paths one-for-one.

## Rationale

The current URI notation cannot meet the approved implementation constraint:
explicit Go 1.22 `http.ServeMux` method/path registrations with a usable
participation UID wildcard. The corrected URI shape is conventional,
unambiguous, and makes the resource UID a complete wildcard segment while
retaining a literal action name. It preserves the approved route arithmetic
and does not introduce a router dependency or an internal method-dispatch
mechanism.

## Reuse-before-build assessment

Canonical standard: `docs/phase2/PHASE2_REUSE_ASSESSMENT_STANDARD.md`

- Disposition: Extend
- Summary of mature candidates / applicable standards:
  - Reuse Go 1.22 standard-library `http.ServeMux` explicit method/path
    registration and `Request.PathValue`.
  - Reuse the approved FEATURE-0015 participation lifecycle, authorization,
    idempotency, audit, and Problem Details contracts unchanged.
- Sovrunn-owned responsibility summary:
  - Register the existing eight participation actions using an executable,
    explicit standard-library route shape.
- Non-goals summary:
  - No new public action, resource, field, scope, state, transition, grant,
    header, body, audit category, idempotency semantic, Problem code,
    persistence dependency, router library, Go implementation, or generated
    Kiro artifact.

## Phase impact

- Current phase allowed: Yes.
- Current phase boundary impact: FEATURE-0015 route-shape correction only;
  no future-feature capability is introduced.

## Conflict check

- Conflicts with accepted DEC/RFC: No.
- Resolution required: Founder approval of this correction; no DEC, RFC,
  baseline, roadmap, or feature-sequence update.

## Required action

- Update architecture and feature authorities to name the eight corrected
  action paths.
- Update the contract specification and Slice 0 steering rule.
- Update the FEATURE-0015 architecture-readiness checker to fail closed if
  a retired action URI appears or a corrected action URI is absent.
- Add ADH-2026-053 to `feature.handoffs` without altering other manifest
  fields.
- Do not update the registry, requirements artifact, design artifact, tasks,
  Go source, DEC, baseline, roadmap, feature sequence, traceability matrix,
  closure matrix, or state files.

## Strict Kiro read/write boundary

Kiro may read only:

- this handoff;
- `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`;
- `.kiro/steering/slice0-contract.md`;
- `.automation/features/FEATURE-0015.control.json`;
- `scripts/feature-0015-architecture-readiness-check.py`.

Kiro may write only:

- `docs/architecture/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/features/FEATURE-0015-canonical-cloud-model-foundation.md`;
- `docs/architecture/vertical-slices/VS-000-contract-specification.md`;
- `.kiro/steering/slice0-contract.md`;
- `.automation/features/FEATURE-0015.control.json`;
- `scripts/feature-0015-architecture-readiness-check.py`;
- this handoff, only to add an application-status annotation.

Kiro must not search, glob, inspect, or follow references outside this list.
If another file is required, it must stop and report that exact file path and
reason without opening it.

## Impacted features

- FEATURE-0015: executable public action URI correction only.
- FEATURE-0016 through FEATURE-0026: no behavior change.

## Acceptance criteria for Kiro update

- [ ] All eight corrected action paths appear exactly in the architecture
  boundary, feature authority, and contract specification.
- [ ] The retired `/{uid}:<action>` form appears in none of the allowed active
  authorities.
- [ ] Each corrected path uses `{uid}` as a complete Go 1.22 `http.ServeMux`
  wildcard segment and names `r.PathValue("uid")` as its extraction method.
- [ ] The architecture continues to state exactly 22 logical endpoint paths
  and exactly 35 explicit method/path registrations.
- [ ] Existing item-action body, `If-Match`, `Idempotency-Key`, lifecycle,
  authorization, audit, and scheduler-expiry semantics are unchanged.
- [ ] The readiness checker fails closed for a missing corrected path or a
  reintroduced retired action URI.
- [ ] `FEATURE-0015.control.json` retains ADH-2026-048 through ADH-2026-052
  and adds ADH-2026-053 with no other field change.
- [ ] No generated Kiro artifact, Go source, registry, requirements, design,
  tasks, DEC, baseline, roadmap, feature sequence, traceability matrix,
  closure matrix, or state file changes.
- [ ] `make feature-0015-architecture-readiness`,
  `make feature-contract-check FEATURE=FEATURE-0015`,
  `make feature-0015-formal-check`, `git diff --check`, and
  `mkdocs build --strict` pass.

## Explicit instructions to Kiro

- Apply no decision beyond this handoff.
- Do not search outside the strict allowlist.
- Do not modify generated Kiro artifacts, Go source, requirements, design,
  tasks, registry, traceability, or state files.
- Do not add a routing library, custom router, reflection, wildcard
  catch-all, path-only registration, or internal HTTP-method dispatch.
- Do not alter the 22-logical-path or 35-registration contract.
- Stop, do not infer, if this handoff cannot be represented with only the
  allowed files.

## Human approval

- Approval status: Approved
- Approved by: Sanjeev Kumar
- Date: 2026-08-13
- Notes: Approved bounded correction for ServeMux-compatible participation action routes.
