Feature: {{FEATURE_ID}} {{FEATURE_TITLE}}
Stage: Design

Generate only `{{DESIGN_PATH}}`. Do not modify requirements, architecture,
handoffs, automation configuration, source, schemas, tests, or tasks.

## Exact context boundary

Read only the files below. This list must exactly match the hashed design
context manifest; do not search for or load "similar implementations" or any
other repository file.

{{CONTEXT_FILES}}

Source precedence is: approved FEATURE-0014 architecture plus ADH-2026-018;
approved requirements; inherited FEATURE-0012 contracts; global repository
rules; supporting scope and reuse records. Lower-precedence context cannot add
or override semantics.

FEATURE-0013 verification is performed by repository preflight. Do not load or
summarize FEATURE-0013 architecture or handoffs. Its design-stage result here
is exactly `NOT_APPLICABLE`; add no decision, audit, operation, projection,
linkage, rationale, actor, subject, or composition machinery.

## Closed design boundary

- Preserve exactly Provider, ProviderLocation, ProviderDatacenter,
  DatacenterFailureDomain, and InfrastructureStack.
- Preserve Organization as the existing optional supply owner through Provider
  `scopeRef`; do not create an Owner kind or use `ownerRef` as scope authority.
- Preserve the exact hierarchy, one immutable immediate parent, same-Provider
  scope, registered-shell behavior, complete five-level path, immutable stack
  identity, one failure-domain parent, and leaf-first deletion.
- Physical containment, shared ancestry, shared owner, and geography imply
  neither connectivity nor isolation. Missing connectivity remains Unknown.
- Reuse FEATURE-0012 grammar unchanged. Select only representation mechanics
  delegated by architecture sections 14 and 15.
- Do not design ResourcePool, ProviderCapability, capacity, compatibility,
  placement, connectivity fields/graphs, adapters, discovery, credentials,
  endpoints, provider calls, provider SDK types, repositories, persistence,
  plugins, provisioning, runtime execution, or speculative extension points.
- Do not add a new normative requirement, public semantic, kind, relationship,
  status meaning, condition meaning, error family, operation, or interface.
- If any requirement cannot be represented without a semantic choice not
  closed or delegated by architecture, emit `ARCHITECTURE_DECISION_REQUIRED`
  and stop the entire stage without writing a partial design.

## Permitted representation decisions

Design may resolve only:

- exact FEATURE-0012-conforming API group and versioned routes;
- minimal field inventory required by approved requirements;
- exact condition names for registered validity and topology completeness;
- mature normalized geographic-code representation where applicable;
- bounded descriptive infrastructure-technology representation;
- finite initial limits and configuration mechanism;
- schema composition, indexes, package layout, and internal structure; and
- immediate-parent representation, including whether FEATURE-0012 `ownerRef`
  mirrors containment without becoming a second authority.

## Mandatory design contents

The design must include:

1. overview and design boundary;
2. resolved delegated design decisions with rationale;
3. component/package architecture with no adjacent-feature scaffolding;
4. resource and field model with an authoritative writer/mutability ledger;
5. API routes and FEATURE-0012 binding;
6. ordered structural, semantic, reference, authorization, concurrency, and
   deletion validation with inherited stable error behavior;
7. current-fact condition behavior without history or audit semantics;
8. security, redaction, no-existence disclosure, and multi-owner isolation;
9. conformance fixtures and verification mechanisms;
10. a requirement traceability table with exactly one row for every
    `F14-REQ-01` through `F14-REQ-31`, mapping each requirement to its
    representation, validator, error/evidence behavior, and no new norm;
11. an architecture traceability table enumerating every `F14-AD-001` through
    `F14-AD-021` without range shorthand;
12. a risk-control table enumerating every `F14-R01` through `F14-R30`, mapping
    each risk to a concrete component, schema constraint, validation, fixture,
    or review mechanism without changing rating, owner, or acceptance;
13. explicit absence ledger for FEATURE-0013/0015/0016/0053 semantics;
14. non-goals and unresolved semantic-gap result; and
15. design-to-requirements completeness and orphan report.

Do not add generic operation, audit, adapter, repository/storage, or provider
integration sections. Do not copy the closed decision/risk registers as a new
authority; record only design disposition and evidence mappings.

Use `SHALL` or `MUST` only when quoting an approved requirement. Design choices
must use descriptive design language and cite their originating requirement.

After writing, read the whole file back and verify that only `design.md` was
created or modified.

