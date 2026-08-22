# FEATURE-0016 F16-R01 Paired Backing-Access Bridge — Design Addendum

## Status and boundary

This addendum is the sole design authority for the post-feature remediation
**F16-R01**. It consumes approved
[ADH-2026-066](../architecture-decision-handoffs/ADH-2026-066-feature-0016-coherent-backing-access-compatibility-bridge.md).
It does not reopen FEATURE-0016 requirements, the approved five-route surface,
the 128 local conformance meanings, or the completed twelve-task plan.

General FEATURE-0016 Cursor execution is frozen. Only F16-R01 may be planned,
reviewed, and executed after this addendum receives its own bounded approval.

**Founder approval:** Sanjeev Kumar approved F16-R01 on 2026-08-22. This
approval authorizes only the implementation scope and proof listed below.

## Closed implementation design

F16-R01 adds an F0016-owned `BackingAccessProvider` over one private,
read-only paired backing-read lease co-located in `internal/cloudmodel/store.go`.
The lease holds the existing FEATURE-0015 store lock once and exposes immutable
copies of the requested `CloudProviderParticipation` and
`InfrastructureStack` only to the provider callback. It is not an F0015 route,
schema, writer, lifecycle transition, or public API.

The provider accepts the authenticated principal, F0016 route intent, backing
UIDs, and inherited FEATURE-0012 grant evaluator. It returns only one internal
disposition:

- `SafeDenied`, mapped by F0016 to the existing audited safe 404;
- `AuthorizationDenied`, mapped by F0016 to the existing audited 403; or
- `Allowed`, containing immutable participation UID, derived CloudProvider UID,
  participation effective-Active value, InfrastructureStack UID/Active value/
  generation, and viability-fingerprint inputs.

No handler receives a lease, raw backing object, or raw existence boolean.
Ordinary create, GET, and LIST calls complete and release the lease within the
provider. Only final qualification invokes a lifecycle-service-scoped callback
that retains the fresh lease through this fixed order:

```text
FEATURE-0015 paired backing-read lease
  → FEATURE-0016 lifecycle-service mutex
  → inherited AuditEvent append
  → F0016 target publication
```

No path holding the lifecycle mutex may acquire the backing lease; no
FEATURE-0015 writer acquires the lifecycle mutex. Final qualification compares
only InfrastructureStack generation and viability fingerprint. Participation
generation is never captured or evaluated as a fence.

GET and LIST use a fresh `Allowed` result only for response-only
`effectiveAvailability`; they never mutate target state or ETag. Participation
suspension is unavailable for new admission only and never stops or mutates an
existing InfrastructureStack, ExecutionTarget, or workload.

## F16-R01 scope and proof

The exact writable/test paths are those already reserved by F16-R01:

- `internal/cloudmodel/store.go`, `internal/cloudmodel/store_test.go`;
- `internal/executiontarget/`;
- `internal/api/executiontarget_*.go` and matching tests; and
- `tests/conformance/feature_0016_test.go`.

Required proof is limited to ADH-2026-066: paired read non-interleaving,
principal-relative safe 404/authorized 403, final-fence non-publication,
absence of a participation-generation fence, response-only GET/LIST freshness,
the stated lock order, FEATURE-0015 regression, and FEATURE-0016 gate success.
No new route, field, state, error, violation, conformance case, or external
effect is authorized.

## Lean approval packet

Review only this addendum together with:

1. ADH-2026-066;
2. §4.15 of the FEATURE-0016 architecture authority;
3. the F16-R01 subsection of `tasks.md`; and
4. the existing FEATURE-0015 store lock contract in `internal/cloudmodel/store.go`.

The review question is limited to whether this addendum faithfully implements
ADH-2026-066 within F16-R01's reserved paths. It must not re-review or alter
the completed twelve-task plan, general FEATURE-0016 design, requirements, or
public contract.
