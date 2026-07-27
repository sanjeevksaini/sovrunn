# Six-Scope AuditEvent Alpha Compatibility Decision

## Purpose

Records the approved compatibility decision for expanding the FEATURE-0012 alpha `AuditEvent` profile from Organization-only to six formal governance scopes per ADH-2026-014.

## Decision

| Attribute | Value |
|---|---|
| Decision type | Compatibility extension |
| Controlling handoff | ADH-2026-014 |
| Feature affected | FEATURE-0012 (source), FEATURE-0013 (owner) |
| Approval | Sanjeev Kumar, 2026-07-27 |
| Status | Approved |

## Scope expansion

| # | Scope | Prior status | New status |
|---|---|---|---|
| 1 | Platform | Not defined in alpha | Added |
| 2 | Organization | Defined in FEATURE-0012 alpha | Preserved |
| 3 | OrganizationUnit | Not defined in alpha | Added |
| 4 | Tenant | Not defined in alpha | Added |
| 5 | Project | Not defined in alpha | Added |
| 6 | ServiceInstance | Not defined in alpha | Added |

## Compatibility treatment

### Affected FEATURE-0012 alpha consumers

The FEATURE-0012 alpha `AuditEvent` profile defined Organization as the governance scope. The following contracts reference or consume this alpha profile:

1. FEATURE-0012 Operation resource — uses Organization scope for Operation-linked audit events.
2. FEATURE-0012 conformance fixtures — test AuditEvent creation with Organization scope.

### Migration path

- The Organization scope remains valid and unchanged.
- Five additional scopes are additive. Existing Organization-scoped events are expected to remain conformant.
- The design intent is that no FEATURE-0012 consumer is broken by the addition of new scope values (verification required by FEATURE-0013 fixtures and regression evidence).
- A scope field that previously accepted only `Organization` will accept the six-value enum. The prior value remains valid.

### Backward-compatible behavior preserved (design intent)

- All existing Organization-scoped `AuditEvent` instances are expected to remain valid.
- FEATURE-0012 conformance fixtures that create Organization-scoped events are expected to pass without modification (verification required by BC-AC-05).
- No field removal, rename, or semantic change applies to existing events.
- Scope does not grant authorization (per ADH-2026-014 and ADH-2026-013).

### Required fixture evidence (gates APPROVED_FOR_CURSOR)

- BC-AC-04: Six-scope compatibility fixtures must exist and pass for Platform, Organization, OrganizationUnit, Tenant, Project, and ServiceInstance.
- These fixtures are implementation artifacts produced during FEATURE-0013 design/task/implementation, not prerequisites for design authorization.

## FEATURE-0012 compatibility assessment

The existing FEATURE-0012 AuditEvent schema and bindings are Organization-only. The six-scope extension is intended to preserve Organization compatibility, but compatibility is not proven until:

1. FEATURE-0013 schemas defining the six-scope AuditEvent envelope are implemented;
2. FEATURE-0013 bindings for all six scopes exist;
3. Six-scope compatibility fixtures pass for Platform, Organization, OrganizationUnit, Tenant, Project, and ServiceInstance;
4. FEATURE-0012 regression tests confirm no breakage from the scope expansion;
5. Baseline checks verify no normative FEATURE-0012 contract is invalidated;
6. Feature gates pass for both FEATURE-0012 regression and FEATURE-0013 acceptance.

Until this evidence exists, the assessment is: the expansion is designed to be additive and backward-compatible, but conformance is unverified.

FEATURE-0012 human semantic review (`docs/reviews/feature-gates/FEATURE-0012-human-semantic-review-evidence.md`) confirmed the AuditEvent structure. The expansion to six scopes was explicitly anticipated by ADH-2026-013, which approved Operation declaration of all six governance scopes without granting authorization.

## Conclusion

The six-scope AuditEvent expansion is designed to be backward-compatible with FEATURE-0012. Compatibility is not yet proven. Executable compatibility fixtures will be produced during FEATURE-0013 implementation to provide runtime verification evidence. Until those fixtures, FEATURE-0012 regression tests, baseline checks, and feature gates pass, post-change conformance remains unverified.

## Human approval reference

ADH-2026-014 approval by Sanjeev Kumar on 2026-07-27 authorizes this compatibility extension.
