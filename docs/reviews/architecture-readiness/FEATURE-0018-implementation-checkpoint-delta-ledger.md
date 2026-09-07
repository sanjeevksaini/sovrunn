# FEATURE-0018 Implementation-Checkpoint Working-Draft Delta Ledger

Controlling handoff: ADH-2026-077 (FEATURE-0018 implementation checkpoint and
task-amendment control).

This ledger records, once and concisely, the classified difference between the
frozen approved FEATURE-0018 semantic design baseline and the later unapproved
working `design.md` draft that ADH-2026-077 discards. Per ADH-2026-077 §1 no
duplicate mechanics prose is preserved for historical recovery; only this
classified ledger is retained.

## Baseline and draft identity

- Frozen approved baseline (retained as `design.md`):
  - commit: `7d4447d7c1d07c37da4af738dfb399adaa9f4ebc`
  - path: `.kiro/specs/governance-iam-approval-exception-foundation/design.md`
  - sha256: `4eb84a0002d8e58de466806b2dcd18085608a0a62d97cc78684b7e98210d92a8`
  - line count: 2573
- Frozen approved requirements (unchanged):
  - sha256: `aed8f347ea154c30292cdd0d8f43841fad6d419b8cf789becacc8b57d9c2a9e4`
- Discarded working draft (not an additional product-semantic authority):
  - sha256: `bc516a50f2705b8ec168b8ed21476ddd88667b37de8e0654933329bd67915219`
  - line count: 3993
- Draft-to-baseline diffstat: +1812 / -392 across the single `design.md` file.

## Summary

The discarded working draft (produced under ADH-2026-075 design normalization
plus the Design Mechanics Contract, and ADH-2026-076 review-authority
realignment) expanded `design.md` by roughly 1420 net lines, almost entirely as
implementation-mechanics prose: exact Go 1.22 signatures, sealed types, named
closed result types, factory ownership, checker symbol lists, and per-mechanic
test intent. ADH-2026-075 made `design-mechanics-contract.json` the sole
authority for that mechanics surface, so the draft prose became a
non-authoritative mirror of the contract. The heading structure is identical to
the baseline except one new heading (DD-14) that elaborates ADH-2026-074 §2
publication-instant mechanics. No delta changes FEATURE-0018 IAM product
semantics.

## Classified delta table

| Delta ID | Location / Section | Change summary | Classification | Disposition |
|---|---|---|---|---|
| D-01 | §5.4 In-memory atomic publication protocol | ~470-line expansion: sealed API surface, named result types, notation; cited as contract mirror | MECHANICS_EXPANSION | DISCARD |
| D-02 | DD-13 (checker enforcement) | Expanded checker symbol ledger and compile-vs-static enforcement mechanics | MECHANICS_EXPANSION | DISCARD |
| D-03 | DD-14 (new heading) | New `operation.PublicationInstant` immutable value, factories, one clock read, overflow rules | NEW_HEADING | DISCARD |
| D-04 | §4.4 HTTP path contract ledger | Route ledger reformatted/expanded; same paths, still contract-only | SEMANTIC_NEUTRAL | DISCARD |
| D-05 | §7.2 Unit/property tests | Added focused mechanics test suites (permit registry, nonce/serial overflow, instant flow) | MECHANICS_EXPANSION | DISCARD |
| D-06 | §5.2 error/outcome mapping | Prose re-expressed; only inherited INTERNAL_ERROR/CONFLICT touched, outcomes unchanged | SEMANTIC_NEUTRAL | DISCARD |
| D-07 | §10.3 / §10.4 unresolved + self-verification | Restated resolution of DD-05..14 mechanics; no new decision | DUPLICATE_PROSE | DISCARD |
| D-08 | Global (multiple DD/§ bodies) | Contract-citation wording threaded throughout ("authority is the Design Mechanics Contract") | DUPLICATE_PROSE | DISCARD |

## Semantic-change statement

No delta changes FEATURE-0018 IAM product semantics. Verified: no new or altered
resource kind, action, HTTP route (the route ledger stays contract-only with
identical paths), writer, public/top-level error code (the diff touches only
inherited `INTERNAL_ERROR`/`CONFLICT`), event taxonomy,
authorization/approval/review/exception semantics, dependency ownership, or
Phase 2R boundary. The local conformance range is unchanged at
`VS0-CF-F18-01..54` in both versions. No `POTENTIAL_SEMANTIC` delta was found.

All deltas are implementation-mechanics elaboration whose authority
ADH-2026-075 relocated to
`.kiro/specs/governance-iam-approval-exception-foundation/design-mechanics-contract.json`.
They are safely discarded with the frozen baseline retained as the approved
design, and no duplicate mechanics prose is preserved for recovery. The Design
Mechanics Contract and its deterministic checker remain available as
implementation checks only; per ADH-2026-077 §4 they are not an LLM
design-review authority after the checkpoint.
