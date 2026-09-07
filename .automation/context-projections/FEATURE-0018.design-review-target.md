---
doc_type: ai_design_review_target_projection
feature: FEATURE-0018
stage: design-review-target
authority: non-authoritative-semantic-projection-of-raw-design
controlling_handoff: ADH-2026-076
generated: true
---

# FEATURE-0018 Design-Review Target Projection (semantic, hash-bound)

This is a non-authoritative, mechanically generated semantic projection of
the raw `design.md` submitted for review, produced under ADH-2026-076 clause 3 as the
controlled FEATURE-0018 design-review **target**. It retains the semantic
architecture, ownership, routes, validation, security, traceability, non-goals
and classification content and omits only the contract-governed executable
mechanics (for which the hash-bound Design Mechanics Contract is the sole
authority under ADH-2026-075) and the non-normative revision history.

The raw `design.md` and the Design Mechanics Contract remain authoritative.
This projection is bound to their exact bytes and to a successful deterministic
contract-check receipt; a source-hash mismatch or a failing contract check
blocks regeneration, so the raw design mechanics can never override the contract
at review time.

## Hash-bound review-target binding

- Raw design source: `.kiro/specs/governance-iam-approval-exception-foundation/design.md`
- Raw design SHA-256 (also hashed by `semantic-delta.py`): `bc516a50f2705b8ec168b8ed21476ddd88667b37de8e0654933329bd67915219`
- Design Mechanics Contract source: `.kiro/specs/governance-iam-approval-exception-foundation/design-mechanics-contract.json`
- Design Mechanics Contract SHA-256: `028d345dc5aa340d50784155d5f7ce9cd1b4640c452ffedd25efebbfb6a159e6`
- Deterministic contract-check receipt (`make feature-0018-design-contract-check`): `PASS: FEATURE-0018 design-mechanics contract is internally consistent, hash-bound to design.md, and cited by design.md`

## Omitted from this semantic target (authoritative elsewhere)

- Contract-governed executable mechanics — DD-13, DD-14, §5.4 in-memory atomic
  publication protocol, and §7.2 focused-mechanics unit/property tests — are the
  Design Mechanics Contract's authority (ADH-2026-075) and are supplied to the
  reviewer as the hash-bound contract in the controlled design-review context,
  not restated here.
- Non-normative history — §10.4 Final self-verification and the Model Execution
  Report — is omitted; it asserts no product semantics.

---

## Retained semantic design content (exact excerpts from raw `design.md`)

<!-- BEGIN RETAINED SEMANTIC EXCERPT: raw design.md lines 1-575 (Front matter, section index, §1 identity/boundary, §2 decisions DD-01..DD-12) -->
---
doc_type: design
feature: FEATURE-0018
title: Governance, IAM, Approval and Exception Foundation
stage: design
status: draft
baseline: ARCH-2026.08-PHASE2R-CANONICAL
controlling_handoff: ADH-2026-070
conformance_handoff: ADH-2026-071
design_clarification_handoff: ADH-2026-072
audit_taxonomy_correction_handoff: ADH-2026-073
local_publication_checkpoint_handoff: ADH-2026-074
design_normalization_handoff: ADH-2026-075
accepted_decision: DEC-0060
change_request: ACR-2026-002
sole_architecture_authority: docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md
requirements_authority: .kiro/specs/governance-iam-approval-exception-foundation/requirements.md
executable_mechanics_authority: .kiro/specs/governance-iam-approval-exception-foundation/design-mechanics-contract.json
ai_load_priority: feature
ai_summary: FEATURE-0018 implementable design preserving approved resources, routes, ownership, mappings and exclusions. Stages 7-9 produce only a sealed non-authoritative candidate with a complete dependency-version claim; state proves admission currentness and authzeval owns publication-time candidate revalidation. Caller and authorization-bearing controller mutations produce audited evaluations, while non-assignable automatic effects publish mutation evidence only. Mutation owners return typed domain-versus-mechanical finalization outcomes and expose neutral result material. Context-bound proofs and runtime Seal checks enforce atomic publication; the static checker covers imports, declarations, typed calls, and the internally satisfiable checker-enforced closed state.ContextBoundChange/ApplyTo/StateEditor/permit-binding boundary that Go 1.22 cannot compiler-seal (state alone declares the interface/editor/ApplyTo; owners implement only their own change and name StateEditor only as their own ApplyTo parameter), and the checker is the required boundary for source-level copied-binding wrappers while runtime permit/context validation guarantees only rejection of absent, mismatched, foreign-transaction, reused, or context-mismatched bindings. Evidence alone constructs required carrier sets, atomic commit precedes result release, and Phase 2R wires deterministic clocks only.
---

# FEATURE-0018 Design: Governance, IAM, Approval and Exception Foundation

> This document translates the approved FEATURE-0018 requirements
> (`.kiro/specs/governance-iam-approval-exception-foundation/requirements.md`,
> REQ-F18-01..24 / AC-F18-01..49) into an implementable representation. It owns
> only implementation mechanics that the requirements or architecture explicitly
> delegate (requirements §9) plus repository-convention structure (Go packages,
> file layout, validation ordering, error mapping, test wiring). It does **not**
> own, add, reinterpret, or transfer any observable behavior, contract, enum,
> scope, action, writer, lifecycle, state, error, or conformance meaning. Those
> remain fixed by F18-RD-01 through F18-RD-24 as carried by ADH-2026-070 (core,
> Appendix A, Appendix B), ADH-2026-071 (conformance IDs only), ADH-2026-072
> (design-authority and audited-evaluation clarification), ADH-2026-073
> (terminal-ExceptionProposal AuditEvent registration), DEC-0060, and the
> sole architecture authority. Inherited FEATURE-0012/0013/0016/0017 contracts are
> reused by exact reference from their canonical owners and are never forked.

---

> **Design Mechanics Contract (SOLE authority for executable mechanics; ADH-2026-075).**
> The machine-checkable, hash-bound
> `.kiro/specs/governance-iam-approval-exception-foundation/design-mechanics-contract.json`
> is the **single, sole, and authoritative source** for this design's non-product-semantic
> executable mechanics:
> the exhaustive repository-owned package-edge ledger; the named Go 1.22 API signatures, result
> variants, factories and permitted callers; the `StateEditor`/`ContextBoundChange` method surface
> and owner-to-method allowlist; the transaction modes and their Begin/admission/finalization/Seal/
> Commit/Abort/reservation-cleanup/waiter-notification/result-release state machines; the
> deterministic identity, publication-instant, sequence, overflow and failure rules; the
> static-checker scope, rules, protected symbols, test-only seams and fixtures; and the exact
> non-runtime proof artifact,
> command and ownership mappings for the active conformance gates. Per ADH-2026-075 this design is
> normalized: the numbered sections below retain only feature **scope, component intent, routes,
> ownership, test intent and traceability** and **cite the contract** for every machine-checkable
> executable-mechanics form rather than restate it. Any executable-mechanics detail (an API
> signature, result variant, factory, permitted caller, checker rule, owner-finalization row,
> transaction-state transition, lock order, identity/sequence/overflow rule, test-seam or fixture)
> appearing in a numbered section is present only as an **exact citation to the contract** and never
> as an independent restatement.
>
> **Authority model (ADH-2026-075).** For **executable mechanics**, the contract is the sole
> authority: a numbered design section **never overrides** the contract, and on any divergence the
> contract's machine-checkable form prevails and the numbered section is corrected to cite it. The
> contract narrows and mirrors the design; it never broadens, adds, or reinterprets any observable
> behavior, contract, enum, scope, action, writer, lifecycle, state, error, DecisionProfile,
> AuditEvent, REQ, AC, conformance identifier, dependency boundary or Phase 2R exclusion. For
> **product semantics** (observable behavior, contracts, enums, scopes, actions, writers,
> lifecycle, state, errors, DecisionProfiles, AuditEvents, REQ/AC and conformance meaning), the
> cited semantic authorities (ADH-2026-070..074, DEC-0060, requirements, the sole architecture
> authority) remain controlling and the contract must conform to them; the contract asserts no
> product semantics and defers every semantic question to those authorities. The contract may never
> be overridden on an executable-mechanics point by a numbered section, and a numbered section may
> never be overridden on a product-semantics point by the contract.
>
> The deterministic pre-review checker `make feature-0018-design-contract-check`
> (`scripts/feature0018-design-contract-check.py`) validates the contract's internal consistency,
> its Go-1.22 signature closedness, and its SHA-256 hash binding to this `design.md`, and confirms
> this `design.md` cites the contract, before the LLM design review. The contract is included as a
> hash-bound source in the controlled design-review projection
> (`.automation/context-projections/FEATURE-0018.design-review.spec.json` /
> `.manifest.json` / `.md`) so independent semantic review covers its contents.

---

> **Canonical section index (format aliases only).** The headings in this block add no
> design semantics, IDs, or normative content. Each is a thin pointer to the domain-specific
> numbered section that already owns that concern; the numbered sections (§1–§10, plus the
> mandatory Architecture Traceability and Conformance Mapping sections) remain authoritative
> and unchanged.

## Overview

See §1 (Identity, stage, inputs, and closed boundary) and §2 (Resolved design decisions) for
the feature scope, closed boundary, and the design decisions that translate the approved
requirements into an implementable representation.

## Architecture

See §3 (Components and repository paths), specifically §3.3 (component boundaries and
single-writer map) and §3.4 (dependency direction, typed ports, and the acyclic prepare/commit
topology), together with §5.4 (in-memory atomic publication protocol) for the runtime structure.

## Components and Interfaces

See §3.2 (proposed FEATURE-0018 component layout), §3.4 (typed ports and dependency edges), and
the exact sealed API surface in §5.4 (in-memory atomic publication protocol).

## Data Models

See §4 (Data / API representation) — §4.1 contract→profile correspondence, §4.3 enums and typed
discriminators, §4.4 literal HTTP path ledger, and §4.5 field ownership — plus the reference and
value types described in DD-03 (§2).

## Error Handling

See §5 (Validation and deterministic error behavior) — §5.1 validation precedence, §5.2 error
and outcome mapping (inherited FEATURE-0012 codes only), §5.3 fail-closed/determinism, and the
publication/abort behavior in §5.4–§5.5.

## Testing Strategy

See §7 (Test and conformance strategy) — §7.1 local conformance, §7.2 unit and property/focused
runtime tests, and §7.3 executable conformance gating.

---

## 1. Identity, stage, inputs, and closed boundary

### 1.1 Identity

| Field | Value |
|---|---|
| Feature | FEATURE-0018 — Governance, IAM, Approval and Exception Foundation |
| Stage | Design (translation of approved requirements into implementable representation) |
| Phase / order | Phase 2R / order 8 |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Sole architecture authority | `docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md` + ADH-2026-070 core/Appendix A/Appendix B |
| Conformance ID authority | ADH-2026-071 (registers stable proof IDs and REQ/AC mappings only) |
| Design clarification authority | ADH-2026-072 (FEATURE-0013/F18 ownership, completed audited-evaluation boundary, design-owned mechanics, and review classification only) |
| Accepted change records | ACR-2026-002; DEC-0060 (Accepted 2026-08-30) |
| Direct requirements input | `.kiro/specs/governance-iam-approval-exception-foundation/requirements.md` |
| Module / toolchain | `github.com/sanjeevksaini/sovrunn`, Go 1.22, stdlib-first (`gopkg.in/yaml.v3` is the only current external dependency) |
| Public routes / persistence / real adapters | None authorized (requirements §1). Design records literal path **contract** shapes only; it wires no live public route, durable store, or real adapter |
| Stage boundary | Modifies only this `design.md`; generates no requirements, tasks, schema, route, automation, source, or architecture change |

### 1.2 Inputs consumed (exact context manifest)

Behavioral authority: the approved requirements plus F18-RD-01..24 via ADH-2026-070,
ADH-2026-071, ADH-2026-072, and DEC-0060. Mechanics/standards consumed by reference:

- FEATURE-0012 API/resource standard (`docs/architecture/api-resource-standard.md`):
  resource profiles, seven ScopeKinds, `metadata.scopeRef` as sole scope authority,
  typed references, layered validation order (§6.10), Problem/error contract (§6.11),
  ETag/`If-Match` optimistic concurrency (§6.12), pagination (§6.13), boundaries (§6.7).
- FEATURE-0013 decision/audit standard (design projection
  `.automation/context-projections/FEATURE-0018.design-dependencies.md`): the immutable
  `DecisionRecord`/`AuditEvent` envelopes, `DecisionProfile` extension mechanism (§6.9.1),
  evidence validation and audit-obligation invariants (§10), sensitivity vocabulary (§12.4),
  and lightweight adoption contract (§28). It supplies no transaction, acceptance service,
  allocator, publisher, durable store, outbox, queue, lock manager, or idempotency store
  (ADH-2026-072 §1).
- FEATURE-0016 ExecutionTarget action registrations
  (`docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`):
  `executiontarget.read`, `executiontarget.qualify`, `executiontarget.write` consumed
  unchanged.
- FEATURE-0017 policy-evaluation seam (`docs/architecture/policy-evaluation-abstraction.md`
  and live `internal/policyeval`): `PolicyEngineAdapter`, `PolicyEvaluationRequest`/
  `PolicyEvaluationResult`, closed `Outcome` vocabulary, digest-keyed deterministic fake,
  pure mapper, injected `TimeSource`.
- Engineering standards: Go guardrails (`docs/engineering/go-coding-guardrails.md`),
  controller/reconciliation model (`docs/architecture/controller-reconciliation-model.md`),
  observability and audit baseline (`docs/architecture/observability-and-audit-baseline.md`).
- Anti-drift steering (`.kiro/steering/slice0-contract.md`) and Slice-0 registry identity.

Live repository conventions verified before stating any package/file (see §3.1).

### 1.3 Closed boundary (design may / may not)

Design **may** choose only mechanisms within the exact approved observable boundary:

1. Concrete deterministic fixture copying/construction/validation and deterministic UTC
   time-supply mechanisms — delegated by REQ-F18-22 (F18-RD-22).
2. The literal HTTP path shape for the closed operation surface under FEATURE-0012
   domain-grouped versioned route rules, without adding an operation, changing authority,
   or converting a controller effect into caller permission — delegated by REQ-F18-02
   (F18-RD-02).
3. Repository-convention mechanics that the FEATURE-0012 design contract expects: Go
   package/module boundaries, struct/JSON representation reusing inherited base types,
   validation-stage organization, error-code mapping to inherited validators, canonical
   schema file placement and profile/boundary/scope annotations, and test/conformance
   wiring.
4. The package, sealed-interface, process-local state, snapshot, locking, reservation,
   waiter, idempotency, FEATURE-0018-local unit-of-work, context-bound receipt, commit/abort,
   cleanup/shutdown, evidence-wiring, static-checker and focused-test mechanics enumerated by
   ADH-2026-072 §3. These mechanics may make closed behavior executable but may not create or
   change observable product semantics or ownership.

Design **may not**: redefine the contract inventory, enums, scopes, actions, writers,
lifecycles, states, validation precedence, eligibility rules, or exclusions; select real
IdP/policy/workflow engines, adapters, credentials, transport runtime, or persistence;
convert examples/journeys/standards mappings into behavior; reinterpret FEATURE-0012/0013/
0016/0017/0020 authority; infer an unregistered action, binding, scope, transition, or
evidence value; or weaken the future-feature exclusions or negative-leakage sentinels.
Any unresolved semantic choice halts with one stop condition (§10.3). The previously ambiguous
candidate/evaluation boundary and FEATURE-0013/local-publication ownership are closed by
ADH-2026-072; no other semantic choice remained (§10.3).

---

## 2. Resolved design decisions

Each decision below is a mechanic, not a semantic. Each cites its delegating or governing
authority. No decision changes any REQ/AC meaning.

### DD-01 — Domain package with Go-enforced subpackage boundaries and an acyclic prepare/commit topology

FEATURE-0018 follows the verified `internal/<domain>/` convention under
`internal/govaccess`. Package boundaries exist only where Go visibility or import direction must
enforce an approved owner boundary; they do not redefine the inherited domain semantics.

The exact component paths and import edges are defined once in §3.2–§3.4 and checked by DD-13.
In summary, no child imports the composition root, no semantic owner imports `uow`, and the neutral
producer/storage/coordinator edges remain acyclic.

Each domain owner has two sealed, owner-private representations. Stages 1–9 produce an immutable
`PreparedIntent` containing only already-validated operation meaning and expected-version/CAS
inputs; it contains no publication timestamp, final state, canonical result, after-state digest,
`MutationDescriptor`, carrier ID or completed idempotency value. After stage 10 captures the
transaction's single `evidence.PublicationContext`, only that same semantic owner implements the
`PreparedIntent.FinalizeAt(state.MutationFinalizationPermit)` constructor that `uow` invokes. The
permit is transaction-created and carries the transaction's exact `PublicationContext`, its
authority mode, and one unforgeable, transaction-scoped `state.PermitBinding` minted from a fresh
`state.FinalizationPermitToken` that the transaction records — still unconsumed — in its private
permit-issuance registry; an authorization-bearing transaction cannot issue a permit until it has
accepted an evidence-owned sealed current-`Allow` proof, while an automatic-mutation-only transaction
can issue only its distinct automatic permit. The constructor returns the owner's sealed
`operation.FinalizationOutcome[FinalizedPreparedChange, OwnerDomainOutcome]`: one immutable published
`FinalizedPreparedChange`, one owner-typed already-registered non-publication `Domain(D)` outcome
(which publishes nothing and returns its existing mapping), or one local mechanical failure. A
`FinalizedPreparedChange` embeds the exact `state.PermitBinding` supplied by the permit and declares
its `PublicationKind()` (`ResourceMutation` or `DomainConclusion`), exposing
`PublicationKind() state.PublicationKind`, `PermitBinding() state.PermitBinding`,
`DomainDecisionMaterials() []evidence.DomainDecisionMaterial`,
`CanonicalResult() (operation.ResultMaterial, bool)` and
`ApplyTo(state.StateEditor) operation.MechanicalFailure` (the closed failure representation fixed by
the contract `stateEditorSurface`; never a bare `error`). A
`ResourceMutation` additionally exposes its non-empty canonical ordered
`EvidenceDescriptors() []evidence.MutationDescriptor` set and applies a non-empty state delta; a publishable `DomainConclusion` — the exact vehicle for a terminal
ExceptionProposal `Deny` — instead exposes
`ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool)`, carries its mandatory
domain-decision material, applies an empty state delta, and creates no domain resource. A
`DomainConclusion` is a `Finalized`, published outcome and is never a `Domain(D)` non-publication
failure.
Finalization performs the owner's publication-time semantic validation, derives all time-dependent
state from the supplied context, and requests the canonical mutation-descriptor set (for a `ResourceMutation`) or
the domain-conclusion descriptor and mandatory decision material (for a `DomainConclusion`) from
`evidence`; it performs no state read, lock operation, clock read, evidence publication or external
I/O. Non-publication `Domain(D)` outcomes retain their canonical owner and existing mapping. Only an invariant, evidence-construction or local
publication-mechanism failure is mechanical; the design does not remap domain denial, deadline,
expiry, validity, lifecycle or conflict outcomes to `INTERNAL_ERROR`.

Only `uow` coordinates the sequence: it carries the sealed intent to stage 10, obtains the context
from the transaction, completes the authorization-bearing current-`Allow` gate when applicable,
obtains the exact owner-bound finalization permit, calls the owning package's finalizer with that
permit, and submits the
sealed finalized change through the transaction's `ApplyFinalized` method. The transaction alone
invokes the change's `ApplyTo` method against its private detached editor and returns an immutable
context-bound application receipt. `uow` and `state` cannot construct, alter or reinterpret either
owner value. `evidence` is deliberately different: it alone declares and constructs
authorization-candidate descriptors, authorization-admission/finalized-candidate proofs,
current-`Allow` proofs, domain-decision materials,
completed-evaluation proofs, mutation descriptors,
FEATURE-0013 carriers and sealed `PreparedEvidenceChange`; it exposes no `ApplyTo(StateEditor)`
method and imports neither `state` nor `uow`. Physical application/publication by `state`/`uow` is
not semantic writer authority.

`state.ContextBoundChange` is deliberately an **exported interface**: `state` alone declares it, and
each owner's `FinalizedPreparedChange` must satisfy it so the state transaction's
`ApplyFinalized`/`ApplyTo` path can apply an opaque, already-finalized change without importing any
owner concrete type. Go 1.22 provides no mechanism to seal an exported interface, so the compiler
alone cannot prevent a foreign package from declaring, implementing, embedding, adapting, wrapping,
or exposing `state.ContextBoundChange` or `ApplyTo`, nor from extracting
`MutationFinalizationPermit.Binding`. This design therefore treats `state.ContextBoundChange` and
`ApplyTo` as a **checker-enforced closed boundary**, not a compiler-sealed one, under five exact,
internally satisfiable rules:

1. `state` alone may declare the `state.ContextBoundChange` interface, declare and construct
   `state.StateEditor`, and invoke `ApplyTo` from its transaction implementation. Each of the nine
   registered owner packages in the DD-13 owner-finalization table may **directly implement only its
   own single registered concrete `ContextBoundChange` finalized-change type**, satisfying the
   interface's method set on exactly that type. No package — the registered owner package itself
   included — may embed, adapt, wrap, alias, forward, delegate, or otherwise re-expose
   `state.ContextBoundChange` or `ApplyTo`, its own or another owner's finalized change, or a copied
   binding through any other type. Each such owner may name `state.StateEditor` **only** as the
   synchronous parameter of its own registered `ApplyTo` implementation, using only its registered
   direct editor methods there. No other package (including `uow`, `command`, root, or a foreign
   helper) may declare, implement, embed, wrap, adapt, alias, forward, delegate, or re-expose either
   symbol.
2. Each registered owner may receive and synchronously use `state.StateEditor` **only** as the
   parameter of its own registered direct `ApplyTo` implementation, using only its registered direct
   editor methods inside that exact body. Owners may not construct, acquire, retain, return, capture,
   escape, forward, embed, wrap, adapt, or otherwise expose `state.StateEditor`, and may not use it
   anywhere outside that exact `ApplyTo` body. Only the `state` transaction implementation may invoke
   `ApplyTo` (against its private detached editor inside `ApplyFinalized`); `uow` and every other
   package may neither obtain a `StateEditor` nor call `ApplyTo` directly.
3. Only each matching owner `FinalizeAt` method may extract `MutationFinalizationPermit.Binding`
   and embed it into the finalized change it returns; `uow` and foreign packages may not extract,
   adapt, or forward a permit binding.
4. DD-13 rule `F18-ARCH-012_CONTEXT_BOUND_CHANGE` enforces rules 1–3 with the exact protected
   symbols, one-rule positive fixtures for all nine registered owners (each declaring its own single
   concrete finalized-change type, naming `StateEditor` only as its own `ApplyTo` parameter, and
   using only its registered direct editor methods), and one-rule negative fixtures for foreign
   implementation, wrapper/embedding, delegated `ApplyTo`, `uow` `ApplyTo`, editor escape, an
   isolated owner-wrapper (a registered owner re-exposing `ContextBoundChange`/`ApplyTo` through
   anything other than its single registered concrete type), and foreign permit-binding extraction.
   This static checker is the required enforcement boundary for source-level copied-binding wrappers.
5. Runtime permit, claim, context, receipt, and consumed-token validation guarantees rejection only
   of a change carrying an absent, mismatched, foreign-transaction, reused, or context-mismatched
   binding: such a change carries no matching unconsumed permit binding and is rejected at
   `ApplyFinalized`/`Seal` (§5.4). This runtime validation does not independently detect a
   source-level wrapper that copies an otherwise-valid binding; that wrapper is prevented statically
   by the DD-13 checker (rule 4), not at runtime.

No design claim asserts that the compiler prevents foreign implementation, embedding, or `ApplyTo`
invocation of `state.ContextBoundChange`, and no design claim asserts that runtime validation
independently detects a copied-binding wrapper.

The approved workflow-trigger, RoleAssignment-writer, Membership-writer, approval-derivation and
FEATURE-0017 seam ownership remains exactly as registered by F18-RD-08/10/11/13/14/16 and the
§3.3 writer ledger; this design represents those authorities with narrow typed ports rather than
restating their behavioral contracts. The compiler enforces sealed-type visibility only for the
owner-private concrete implementations and the sealed value types (not for the exported
`state.ContextBoundChange` interface or `ApplyTo`); the named DD-13 checker enforces imports,
authorized call sites, and the `state.ContextBoundChange`/`ApplyTo` closed boundary above. No
external dependency is introduced.

### DD-02 — Deterministic, in-process, side-effect-free foundation (no live server wiring)

Per requirements §1 (no public routes/persistence/real adapters) and architecture §3, the
implementable artifact is a deterministic in-memory contract library plus local conformance,
not a live HTTP service. Controllers are synchronous in-process functions; no `internal/server`
registration, durable store, background reconciler, or external client is created. The neutral
`state/` aggregate owns a package-private root, a short-held `reservationMu` for in-flight
coordination, and a `stateMu` for coherent snapshots and publication. The locks are not
per-resource maps and follow the exact order and lifetime in DD-10/§5.4. All reads return deep
copies. Owners prepare opaque changes but never acquire, retain, commit or swap state; only `uow`
coordinates the sealed state capabilities. (Governing: REQ-F18-22; architecture §3;
controller-reconciliation model.)

### DD-03 — Reuse inherited base types by import, never fork

Canonical grammar types are imported from their owners and never redefined:
`apimeta.TypedRef`/`ObjectMeta`/`TypeMeta`/`ScopeRef` (FEATURE-0012 `internal/apimeta`);
the layered validation stages and Problem mapping (`internal/apivalid`, `internal/apiproblem`,
`internal/apiref`, `internal/apicond`); `DecisionProfile`/`DecisionRecord`/`AuditEvent`/
`EvaluationResult` (`internal/decision`); the policy seam
(`internal/policyeval.PolicyEngineAdapter`, `PolicyEvaluationRequest`, `PolicyEvaluationResult`,
`Outcome`, `TimeSource`); and the FEATURE-0016 action identifiers as consumed constants.

The four FEATURE-0018-owned reference values are **distinct** value types with their own
canonical semantics (requirements §3.2); they are **not** aliases, wrappers, or constraints
over `apimeta.TypedRef`, and collapsing them would lose required identity, discriminator, and
union semantics. Each is defined in `internal/govaccess/model` and reuses only the appropriate
inherited FEATURE-0012 primitive by composition — never by aliasing the whole `TypedRef`:

- `PrincipalRef` — the FEATURE-0018-owned already-authenticated issuer/subject identity value
  (REQ-F18-03). It carries the closed principal category (`Human | Workload | System`) plus the
  stable issuer/subject identity fields fixed by `VS0-SCHEMA-020`; it is not a `TypedRef` and
  carries no `apiVersion`/`kind`/`name`/`uid` resource-reference grammar. Its category
  discriminator is a closed Go string enum with `Valid()`.
- `AccessGroupRef` — a UID-pinned canonical reference to exactly one `AccessGroup` resource
  (REQ-F18-04). This is the only one of the four that is a resource reference, so it reuses
  `apimeta.TypedRef` by an embedded field pinned to `Kind == AccessGroup` (validated, not
  redefined); it is not interchangeable with `PrincipalRef`.
- `RoleHolderRef` — a **closed discriminated union** of exactly one `PrincipalRef` **or** one
  `AccessGroupRef` (REQ-F18-04/08), with exactly one active variant and a required
  discriminator; a mixed, empty, or unknown holder fails closed. It is modeled as a struct with
  one discriminator field and the two mutually exclusive optional variant fields, never as a
  bare `TypedRef`.
- `EligibilityRef` — an exact-Human eligibility entry containing **exactly one Human**
  `PrincipalRef` (REQ-F18-12/14/16), satisfied only by exact principal equality. It wraps one
  `PrincipalRef` constrained to category `Human`; roles, assignments, groups, `Membership`, and
  claims can never populate or satisfy it (SEC-06). It is not a `TypedRef` and is not a
  `RoleHolderRef`.

Each type has a deterministic structural validator that rejects any cross-type substitution.
(Governing: COMPAT-01..04; requirements §3.2; SEC-06; FEATURE-0013 adoption §28.5.)

### DD-04 — FEATURE-0012 resource-profile correspondence is consumed from the registry, not selected

Each contract's FEATURE-0012 resource profile is fixed by F18-RD-02 and its `VS0-SCHEMA-*`
registry entry; design consumes it generically. The correspondence below is the
representation mapping design **follows**, grounded in the explicit requirement lifecycle
text. F18-RD-02, approved requirements, the canonical model/catalog, and registered artifacts
must agree exactly. Design selects no winner; any mismatch stops with
`ARCHITECTURE_DECISION_REQUIRED` (ADH-2026-072 §4). See §4.1.

### DD-05 — Deterministic UTC time supply (delegated mechanic 1a)

Time is obtained only from an injected deterministic UTC source, reusing the established
`policyeval.TimeSource`/`FixedTimeSource` pattern and never calling `time.Now`. The true leaf
`internal/govaccess/clock` package defines `Clock` (`NowUTC() time.Time`) and deterministic
`FixedClock`/manually advanced fixture implementations only. Phase 2R does not implement or wire
`UTCClock`, system wall time or a production clock adapter; that runtime responsibility remains
outside this feature. The root injects one clock into `state.Store` and every time-dependent owner.
`state.Store` captures the authoritative `publicationInstant` under `stateMu` as the immutable
neutral `operation.PublicationInstant` defined by DD-14, reading the injected clock exactly once;
§5.4 and DD-14 define its construction, binding, extraction, access and reuse. (Delegating:
REQ-F18-22; governing: OPS-01.)

### DD-06 — Deterministic synthetic fixtures: immutable, validated-before-use, zero I/O (delegated mechanic 1b)

Fixtures are constructed as in-package Go value builders returning deep copies (no shared
mutable slices/maps exposed to consumers, per Go guardrails §13/§27). A `fixtureset`
constructor validates the whole graph (referential closure, enum membership, scope/reference
compatibility, no future-owned/prohibited field) and rejects an invalid set **before** any
evaluation, returning a deterministic error. Fixtures perform zero network/file/external I/O
and increment no external-call counter (an in-memory `externalCallCount` sentinel stays `0`).
The fake privileged/approval/JIT/break-glass/expiry/revocation/review/exception executor is a
pure deterministic function of fixture state and the injected clock; it simulates no IdP,
policy engine, workflow engine, or CloudProvider. (Delegating: REQ-F18-22; governing:
OPS-02, VS0-CF-F18-53.)

### DD-07 — Literal HTTP path shape as a contract-only ledger (delegated mechanic 2)

The literal path shape for the closed operation surface follows FEATURE-0012 §6.1
domain-grouped versioned route rules and the established Go 1.22 `http.ServeMux` collection/
item/action forms already used by FEATURE-0015/0016 (ADH-2026-053 action form
`/{uid}/actions/<action>` read via `r.PathValue("uid")`). Because no public route is
authorized this phase, this is a **CONTRACT_ONLY/NO_TASK** design artifact (§4.4): it records
the path templates that a future transport feature would register; it wires no live route and
converts no controller effect into caller permission. The caller-operation surface each path
covers is transcribed one-to-one from the ADH-2026-070 Appendix B closed v1 operation-and-
mutability table (F18-RD-02) and its F18-RD-06 `ActionTargetBinding` registrations, so no caller
operation is dropped, renamed, or reclassified and no controller-only effect is promoted. The
exact API-group token per kind is owned by the canonical contract catalog / finalized data model;
§4.4 binds each path to that registration and adopts only the FEATURE-0012 domain-grouping rule.
(Delegating: REQ-F18-02; governing: F18-RD-06, ADH-2026-070 Appendix B.)

### DD-08 — Authorization/approval/review/exception logic is deterministic FEATURE-0018-local algebra

The exact algebra and FEATURE-0017 invocation semantics are consumed from REQ-F18-09..11 and
F18-RD-09/11 without restatement. Design represents them as pure deterministic functions over an
immutable snapshot and injected clock. The pure evaluator is never a bearer boundary; only the
post-stage-10 finalization path in DD-12/§5.5 may release `AuthorizationResult`. No policy engine
is introduced.

### DD-09 — Three FEATURE-0013 profiles registered via the existing extension mechanism; no competing envelope

`authorization-decision/v1`, `approval-decision/v1`, and `exception-decision/v1` are registered
through the FEATURE-0013 `DecisionProfile` extension process (projection §6.9.1) and reuse the
existing envelope, authority, reason, obligation, projection, validity, and audit-linkage
mechanics unchanged. The exact `## FEATURE-0013 Adoption` section assigns profile-definition,
bundle-loading, composition and test ownership. FEATURE-0018 creates no fourth profile and no parallel envelope; AccessReview stays a
LongRunningOperation. (Governing: REQ-F18-19; FEATURE-0013 §6.9.1, §28.)

### DD-10 — FEATURE-0018-local, non-durable all-or-nothing in-memory publication

FEATURE-0018 uses independent reservation and aggregate-state locks so duplicates can observe an
in-flight owner without entering its publication critical section. It defines distinct sealed
`CallerMutationTransaction`, `ControllerMutationTransaction`, and `EvidenceTransaction` modes.
Caller mode owns a reservation and completed result; controller mode uses the inherited exact
transition/version CAS and forbids caller-idempotency state; evidence-only mode publishes no domain
state. Mutation admission also declares `AuthorizationBearing` or `AutomaticMutationOnly`.
Caller mutations are always authorization-bearing. A controller mutation is authorization-bearing
only when its closed trigger actually carries a current `Allow | Deny` candidate; expiry,
reconciliation and other non-assignable automatic effects are mutation-only. `evidence` alone
constructs and validates carrier sets, while `uow` coordinates acceptance.
Mutation publication is explicitly two-phase: the semantic owner prepares a publication-time-free
intent before stage 10; reservation/CAS then captures the authoritative context; the owner finalizes
its domain change, canonical result and descriptor from that exact context before evidence or
idempotency derivation. Neither `Begin` API accepts a finalized change or descriptor.

The captured context's `publicationSequence` is **provisional**: `state` derives it as
`current + 1` while holding `stateMu`, binds that provisional value into the sealed, state-issued
`evidence.PublicationContext`, and persists it in `state.Store` only as part of a successful
aggregate commit — after revalidating, still under `stateMu`, that the provisional value is still
the next sequence for the committed root. An abort, failed CAS, cancellation, evidence failure or
panic-cleanup path publishes nothing and consumes no sequence, so no aborted work creates an
observable sequence gap. This is local, deterministic in-memory state management; it is neither
durable sequencing nor a FEATURE-0013 responsibility (ADH-2026-074 §1).

The exact API, lock order, mode-specific seal/termination rules, panic/cancellation/shutdown
lifecycle and commit visibility order are defined once in §5.4. The mechanism is deterministic,
process-local and non-durable; `## FEATURE-0013 Adoption` is authoritative for the unchanged
dependency boundary. (Governing: REQ-F18-20/21, F18-RD-20/21; ADH-2026-074 §1.)

### DD-11 — Idempotency/concurrency reuse FEATURE-0012 mechanics with FEATURE-0018 replay rechecks

Optimistic concurrency uses the inherited opaque `resourceVersion`/`If-Match` contract
(FEATURE-0012 §6.12); the first valid terminal transition wins under CAS (§5.4).

The caller-supplied lookup identity is deliberately separate from the immutable request binding:

```text
IdempotencyLookupKey {
  actorRef
  operation
  clientIdempotencyKey
}

RequestBinding {
  exactTarget
  scopeRef
  canonicalRequestDigest
}
```

`state.Store` alone owns the private in-flight entry and table keyed by `IdempotencyLookupKey`,
`InspectOrReserveCallerMutation`, `CallerReservationLease`, reservation cleanup and waiter
signaling; the committed aggregate root stores `CompletedRecord` values under the same key.
`idempotency` owns no in-flight table, reservation type or lifecycle ownership. Every in-flight
entry and completed record contains the full `RequestBinding`. Target, scope and digest are
therefore compared **after** lookup and are never part of the lookup key. This makes reuse of one
caller key for another request detectable.

This table applies only to caller-keyed mutations. Controller-owned operations have no
caller key, reservation or completed replay result; `BeginControllerMutation` reuses inherited
resource-version/transition CAS so a repeated reconciliation is the registered effective-once,
same-decision or no-op outcome. No synthetic caller key or second idempotency mechanism is created.

Every replay and woken waiter restarts the ordered pipeline and reaches stage 10 with a fresh sealed
`OperationCandidate` and safe visibility finding before any stored result can be selected. The
stage-11 evidence set is then accepted before the stored result is returned. The neutral
`operation.CompletionBinding` carries the sealed state-issued publication instant, and
`idempotency.PrepareCompletedResult` derives the fixed 24-hour completed-record deadline from that
one instant — performing no second clock read and importing neither `state`, `evidence` nor `clock`;
completed-record retention therefore honours the 24h horizon computed from that sealed instant and
never extends JIT validity (ADH-2026-074 §2).
The typed `IdempotencyLookupKey`,
`RequestBinding` and `CompletedRecord` shapes, the neutral `PreparedCompletedResult`, and the pure
completed-record lifecycle helpers live in the neutral child package
`internal/govaccess/idempotency`, which imports only `model`, `apivalid` and `operation`.
`idempotency` owns no in-flight table, reservation type, caller-reservation lease or reservation/
waiter lifecycle. `state.Store` alone owns the private in-flight entry/table,
`InspectOrReserveCallerMutation`, `CallerReservationLease`, reservation cleanup and waiter
signaling; it stores the neutral idempotency values, and `uow` is the sole caller of the
completed-record lifecycle helpers. §5.4 is authoritative for inspection,
waiting, reservation ownership and publication ordering. (Governing: REQ-F18-13, REQ-F18-21; OPS-03;
FEATURE-0012 §6.12; ADH-2026-074 §§2–3.)

### DD-12 — Audited authorization-evaluation service wrapping the pure algebra (F18-RD-20)

`authzeval.Evaluator` is pure and returns only a sealed internal
`AuthorizationCandidate`. Stages 8–9 add their sealed internal findings to one
`OperationCandidate` carrying an evidence-owned `AuthorizationCandidateDescriptor` and an
authzeval-owned, state-verifiable `AuthorizationCurrentnessClaim`; none of those
values is an authorization evaluation, `AuthorizationResult`, audit publication, bearer value, or
downstream authority. After the operation's applicable stage-10 coherence/CAS checks and
dependency-version validation succeed, `uow` passes the exact transaction context and its sealed
admission proof to `authzeval.FinalizeCandidateAt`. That semantic-owner method revalidates the
candidate's publication-time predicates and returns either a sealed current candidate or its
already-registered non-publication outcome. `uow` then asks
`evidence.CompleteAuthorizationCandidate` to bind the sealed current candidate to the same proof
and context. For a mutation, `evidence.RequireCurrentAllow` first converts only a finalized
current `Allow` candidate into a sealed `CurrentAllowProof`; a finalized `Deny` yields only the
internal `RouteDenyToEvidenceOnly` direction, aborts the mutation transaction, and can be completed
and audited only through a fresh evidence-only transaction. The mutation transaction must accept
that exact `CurrentAllowProof` before it can issue any semantic-owner finalization permit.
Successful construction of the sealed `CompletedAuthorizationEvaluation` after all applicable
stage-10 mutation finalization succeeds is the
single point at which the candidate becomes a completed current `Allow | Deny` evaluation. Its
required carrier set must then be accepted and committed before release. Conflict, terminal CAS
loss, Wait/retry, controller effective-once NoOp, and automatic mutation-only paths discard or never
create a candidate, return no AuthorizationResult, and are not completed authorization evaluations.
Only `uow` may invoke `authzeval.ReleaseAfterPublication`, and only after successful evidence
commit; that call constructs and releases the `AuthorizationResult` only, takes no
`operation.ResultMaterial`, and a caller mutation's owner-finalized `operation.ResultMaterial` is
instead consumed solely by `idempotency.PrepareCompletedResult` for its completed idempotency record.
Root only wires these components. The exact denied, read, replay and mutation flows are
defined once in §5.5. FEATURE-0013 supplies only carrier and validation contracts
(ADH-2026-072 §§1–2; REQ-F18-20/21).

<!-- END RETAINED SEMANTIC EXCERPT: raw design.md lines 1-575 -->

<!-- BEGIN RETAINED SEMANTIC EXCERPT: raw design.md lines 904-1442 (§3 components/ownership/writers/topology, §4 data/API/route ledger, §5.1-§5.3 validation and error mapping) -->
## 3. Components and repository paths

### 3.1 Verified live conventions (basis for all stated paths)

- Domain feature package: `internal/<domain>/` with `model/` (structs/enums) and `validate/`
  (deterministic validators) subpackages — confirmed for `internal/cloudmodel/{model,validate}`
  and `internal/executiontarget/model`.
- Shared FEATURE-0012 primitives: `internal/apimeta`, `internal/apivalid`, `internal/apiref`,
  `internal/apiproblem`, `internal/apicond`, `internal/apischema`, `internal/apiconform`.
- FEATURE-0013 domain: `internal/decision` (+ `compose/`, `graph/`, `bundle/`).
- FEATURE-0017 seam: `internal/policyeval`.
- FEATURE-0016 domain: `internal/executiontarget`.
- Conformance: `tests/conformance/feature_0016_*.go` / `feature_0015_test.go` precedent.
- Canonical machine-readable contracts: the manifest-controlled VS-000 contract registry
  (`docs/architecture/vertical-slices/VS-000-contract-registry.yaml`) carries the closed,
  machine-readable schema entry for each kind inline (`required`/`optional`/`refs`/
  `classification`/`redaction`/`mutability`/`retention`/`slice0Constraint`); the FEATURE-0018
  schema IDs `VS0-SCHEMA-020..024,062..068` are supplied there and **carry no `externalRef`**
  to any `api/schemas/*.json` file (verified: unlike FEATURE-0012/0013 entries 001–005, no
  FEATURE-0018 entry declares an `externalRef`). No per-kind `api/schemas/<kind>.json` file
  exists or is registered for FEATURE-0018.
- Module `github.com/sanjeevksaini/sovrunn`, Go 1.22, stdlib-first.

### 3.2 Proposed FEATURE-0018 component layout (IMPLEMENT unless marked)

| Path | Responsibility | Reuses |
|---|---|---|
| `internal/govaccess/` (root, composition/wiring only) | Controller wiring and composition **only**: constructs child packages; injects the single deterministic clock into `state.Store` and time-dependent owners; wires caller command services, the pure evaluator, `uow` coordinator and post-commit result releaser without invoking them; and obtains the exact local FEATURE-0013 adoption bundle bytes from `evidence/profiles.go` and loads them through existing `decision/bundle.Load`. Constructs no record, holds no live state, opens no transaction, obtains no editor and exposes no pre-stage-10 publisher | `roleassign`, `membership`, `roledefinition`, `accessgroup`, `approval`, `privileged`, `review`, `exception`, `governanceprofile`, `grantintent`, `approvalreq`, `policyseam`, `validate`, `projection`, `fixtures`, `command`, `uow`, `authzeval`, `state`, `evidence`, `clock`, `decision/bundle` |
| `internal/govaccess/model/` | Go structs and closed enums for the ten contracts and the supporting values; the four distinct FEATURE-0018-owned reference values `PrincipalRef` (closed-category identity value, `VS0-SCHEMA-020`), `AccessGroupRef` (UID-pinned `AccessGroup` reference embedding `apimeta.TypedRef` pinned to `Kind==AccessGroup`), `RoleHolderRef` (closed `PrincipalRef`-or-`AccessGroupRef` discriminated union), and `EligibilityRef` (exactly one Human `PrincipalRef`); plus `ActionTargetBinding`; the closed enums of requirements §3.3; and immutable deep-copied `ApprovalDecisionFacts`, `ExceptionDecisionFacts`, and owner-finalized `MutationEventFacts` projections. The first two contain exactly the corresponding F18-RD-19 profile inputs and terminal result; `MutationEventFacts` contains only the already-registered event taxonomy, resource/version linkage, action and after-state digest. These projections are design-level carrier inputs, not resources, decisions or public contracts. Only `AccessGroupRef` reuses `apimeta.TypedRef`; the other three do not (DD-03) | `apimeta` (via `AccessGroupRef` only) |
| `internal/govaccess/operation/` (**package**) | Neutral leaf protocol values shared by semantic owners, state, evidence and idempotency: immutable deep-copied `ResultMaterial`; the immutable sealed `PublicationInstant` (DD-14) wrapping one unexported UTC `time.Time` with read-only `UTC()`/`Equal` accessors, constructed through the exported `operation.NewPublicationInstant(time.Time)` factory that validates a non-zero canonical UTC input and that only `state` is permitted to call (DD-13 checker, not the compiler); sealed mechanical `AppliedChangeLink`, `AppliedConclusionLink` and `CompletionBinding` (the last carrying the sealed state-issued `PublicationInstant`); plus generic `FinalizationOutcome[T,D]` with exactly one typed `Finalized(T)`, `Domain(D)` or `MechanicalFailure` variant, where the `Finalized(T)` change declares a closed `PublicationKind` (`ResourceMutation` or `DomainConclusion`). `state` alone invokes `NewAppliedChangeLink` from its private admitted claim, exact publication binding and actual shadow/index deltas; `evidence` alone invokes `NewCompletionBinding` from its validated plan's `PublicationInstant`/binding and mutation digests. Neither neutral value contains authorization, approval, exception or lifecycle meaning. Each owner supplies its own already-registered domain/lifecycle outcome type as `D`; `MechanicalFailure` can identify only local invariant, evidence-construction or publication-mechanism failure. It defines no public error, lifecycle or domain semantics | `model` |
| `internal/govaccess/clock/` (**package**) | True leaf: `Clock` (`NowUTC() time.Time`) plus deterministic `FixedClock` and manually advanced fixture clock only. No `UTCClock`, `time.Now` or production wall-clock wiring is implemented in Phase 2R. The root injects one instance into `state.Store` and owner-local clock ports (DD-05) | — |
| `internal/govaccess/validate/` | Deterministic per-contract structural + local-semantic + cross-field + scope/reference validators, ordered per the F18-RD-21 precedence (§5.1); emits inherited Problem/violation codes only; validates the FEATURE-0018 `model` contract types by reference | `model`, `apivalid`, `apiref`, `apiproblem`, `apicond` |
| `internal/govaccess/authzeval/` (**package**) — files `authz.go`, `authzservice.go`, `eligibility.go` | Pure `Evaluator` produces sealed `AuthorizationCandidate`; stage coordinators build sealed `OperationCandidate` carrying one evidence-owned descriptor and one `AuthorizationCurrentnessClaim` whose complete evaluated snapshot/dependency-version set is constructed through the neutral `state` API. Neither candidate is a completed evaluation or externally releasable. `FinalizeCandidateAt` alone revalidates candidate-owned publication-time predicates against the exact transaction context and evidence-owned admission proof without state/clock reads; it returns a typed finalized candidate or the internal `ReevaluateCandidate` direction when the candidate is no longer current. Re-evaluation is not a public error or completed evaluation. For a mutation, only `evidence.RequireCurrentAllow` can convert that finalized candidate into the sealed current-`Allow` proof accepted by the transaction; Deny is routed to a fresh evidence-only transaction. `ReleaseAfterPublication` is the sole `AuthorizationResult` constructor and DD-13 permits only `uow` to call it after applicable stage-10 success, completed-evaluation construction, successful evidence publication and receipt of the committing state transaction's sealed `CommitReceipt`. It never imports `uow`, root or a domain mutation owner | `apimeta`, `model`, `decision`, `operation`, `evidence`, `state`, `clock` |
| `internal/govaccess/command/` (**package**) | In-process caller-command orchestration only. Its `RoleAssignmentRevokeService` evaluates the existing RD-08 caller operation, asks `roleassign.DirectRevokePort` for the owner-sealed intent, and passes candidate + intent to the abstract `uow.MutationCoordinatorPort`. It owns no route, transaction, state, result, evidence or RoleAssignment writer semantics | `authzeval`, `roleassign`, `uow` |
| `internal/govaccess/grantintent/` (**package**) | **Limited to** the direct non-delegable RoleAssignment grant intent, one `RoleAssignmentProposal` admission, and `RoleGrantRequirementDeriver` derivation (`RoleGrantRequirementDeriver.DeriveForRoleGrant`); derives the RoleGrant `ApprovalRequirement`; when a grant is to be published it submits an exact immutable trigger intent to `roleassign.GrantPort` (RD-13); it never constructs a `PreparedRoleAssignmentIntent` or `FinalizedRoleAssignmentChange` and owns no activation/revocation/replacement/due-advancement trigger (DD-01) | `model`, `approvalreq`, `roleassign` |
| `internal/govaccess/roleassign/` (**package**) | Sole RoleAssignment publisher/writer boundary — **VS0-WRITER-008** (spec, lifecycle, status, revocation, replacement, expiry, review-due). It is the sole constructor of sealed `PreparedRoleAssignmentIntent` and `FinalizedRoleAssignmentChange`: the intent is publication-time-free; its owner-implemented `FinalizeAt(state.MutationFinalizationPermit)`, invoked only by `uow`, uses the permit's `PublicationContext` to validate time, derive final state/result and request the exact descriptor, and embeds the permit's `PermitBinding` into the finalized change; only the finalized change exposes publication kind/permit binding/descriptor/neutral result/`ApplyTo`. Exposes **three distinct narrow typed workflow-trigger ports** keeping RD-13/14/16 separate: `GrantPort` (grant/materialization intent, called only by `grantintent`), `ActivationPort` (activation + privileged revocation, called only by `privileged`), `ReviewRemediationPort` (revocation/replacement/due-advancement, called only by `review`). Separately exposes `DirectRevokePort.Prepare(DirectRevokeIntent)` for the existing RD-08 caller operation; it performs only RoleAssignment-owned semantic preparation and returns the same sealed intent type. `command.RoleAssignmentRevokeService`, not `roleassign`, owns caller-flow coordination through `uow.MutationCoordinatorPort`. The direct port is **not** a fourth workflow trigger, route or writer. `roleassign` does not derive approval and imports neither `approvalreq`, `uow` nor `command` (DD-01/DD-13) | `decision`, `model`, `operation`, `state`, `evidence` |
| `internal/govaccess/membership/` (**package**) | Sole identity-membership writer boundary — **VS0-WRITER-022**; `MembershipRelationshipKey` uniqueness; freshness/provenance; derives its own `ApprovalRequirement`; **publishes Membership state and any AccessGroup assignment-effect expansion entirely within this owner** through its owner-sealed publication-time-free intent and context-finalized change; it **submits no RoleAssignment grant intent** to `grantintent` or `roleassign` and never constructs a RoleAssignment intent/change (DD-01) | `decision`, `model`, `operation`, `state`, `evidence`, `approvalreq` |
| `internal/govaccess/roledefinition/` (**package**) | RoleDefinition version lifecycle, classification (base/effective), supersession; sole constructor of its publication-time-free intent and context-finalized change (DD-01) | `decision`, `model`, `operation`, `state`, `evidence` |
| `internal/govaccess/accessgroup/` (**package**) | AccessGroup lifecycle, owner/display-intent update, ownership resolution, membership-eligibility state; sole constructor of its publication-time-free intent and context-finalized change (DD-01) | `model`, `operation`, `state`, `evidence` |
| `internal/govaccess/approval/` (**package**) | ApprovalPolicy version lifecycle/validation; ApprovalRequest lifecycle/evidence/terminal-immutability; **consumes** each originating controller's derived `ApprovalRequirement` (does **not** derive it, does **not** import `roleassign`, does **not** import `policyseam`); sole constructor of its publication-time-free intent and context-finalized change (DD-01) | `decision`, `model`, `operation`, `state`, `evidence`, `approvalreq` |
| `internal/govaccess/approvalreq/` (**package**) | Typed `ApprovalRequirementDeriver` **ports** (§3.4), one per intent kind, implemented independently by each originating operation owner (`RoleGrantRequirementDeriver` by `grantintent`; membership/privileged/exception derivers by their owners); imports neither `roleassign` nor `policyseam`; carries no shared derivation body | `model` |
| `internal/govaccess/privileged/` (**package**) | PrivilegedAccessRequest JIT + break-glass flow, activation floor, deadlines; derives its own `ApprovalRequirement`; **sole caller** wired to `policyseam.Port` for the pinned subject; submits activation and privileged revocation intents to `roleassign.ActivationPort` (RD-14) only; sole constructor of its publication-time-free intent and context-finalized change (DD-01) | `decision`, `model`, `operation`, `state`, `evidence`, `policyseam`, `roleassign`, `approvalreq` |
| `internal/govaccess/review/` (**package**) | AccessReview snapshot, dispositions, beneficiary conflict; submits revocation/replacement/due-advancement intents to `roleassign.ReviewRemediationPort` (RD-16) only; sole constructor of its publication-time-free intent and context-finalized change (DD-01) | `decision`, `model`, `operation`, `state`, `evidence`, `roleassign` |
| `internal/govaccess/exception/` (**package**) | ExceptionProposal → ExceptionGrant append-only evidence, overlap rejection, linked revocation; derives its own `ApprovalRequirement`; sole constructor of its publication-time-free intent and context-finalized change (DD-01) | `decision`, `model`, `operation`, `state`, `evidence`, `approvalreq` |
| `internal/govaccess/governanceprofile/` (**package**) | GovernanceProfile v1 component/applicability validation (local only); sole constructor of its publication-time-free intent and context-finalized change (DD-01) | `decision`, `model`, `operation`, `state`, `evidence` |
| `internal/govaccess/state/` (**package**) | Neutral non-durable storage. Owns root/shadows, `reservationMu`, `stateMu`, publication sequence, in-flight/completed tables, accepted carrier storage, sealed dependency-version sets/currentness claims, `MutationAdmission`, transaction-issued `MutationFinalizationPermit` with its transaction-private permit-issuance registry and unforgeable `PermitBinding`, the closed `PublicationKind`, the transaction-private `AppliedChangeReceipt` registry, and the sealed post-install `CommitReceipt`. Its sole constructor `state.NewStore(clock.Clock)` obtains one `evidence.StoreNonce` via `evidence.NewStoreNonce()` and fails closed (returning `operation.MechanicalFailure`, no partial Store) if that nonce counter is exhausted. Every identity counter it owns — `transactionSerial`, per-root `publicationSequence`, `FinalizationPermitToken` — is strictly monotonic, never reused, and fail-closed at `math.MaxUint64` (no wraparound to 0 or a reused value). Exposes caller inspection/lease, controller begin, evidence-only begin, and the three sealed transaction interfaces in §5.4. Every successful authorization-bearing begin revalidates the complete admitted dependency set, captures the deterministic context and asks `evidence.BindAuthorizationAdmission` for one proof binding both while still under `stateMu`. An authorization-bearing transaction issues no finalization permit until `AcceptCurrentAllow` validates the exact evidence-owned proof; automatic mode issues only an automatic permit. `ApplyFinalized` first proves the change's `PermitBinding` against its permit-issuance registry (matching unconsumed token, equal participant claim, exact context) and consumes that entry, then invokes the admitted finalized change against the private editor; for a `ResourceMutation` it constructs a neutral `operation.AppliedChangeLink` from the admitted claim/context and non-empty actual shadow/index deltas and asks `evidence.BindAppliedChange` to create the proof and receipt, and for a `DomainConclusion` it proves an empty delta, constructs `operation.AppliedConclusionLink`, and asks `evidence.BindAppliedConclusion`. `Commit` rejects (via the `state.CommitOutcome` `MechanicalFailure` variant) any call before a successful `Seal` or with an absent required carrier/completed-result linkage, installing nothing; otherwise it installs exactly one sealed root/evidence state and returns its opaque `CommitReceipt` in the `Committed` variant; the provisional `current + 1` publication sequence bound into the captured context is persisted only at that successful aggregate commit, after `state` revalidates under `stateMu` that it is still the next sequence for the committed root, and every abort, failed CAS, cancellation, evidence failure or panic-cleanup path consumes no sequence; state never constructs/finalizes domain meaning, performs semantic time validation, or constructs descriptors, canonical results, domain-decision materials or carriers | `model`, `operation`, `idempotency`, `decision`, `evidence`, `clock` |
| `internal/govaccess/uow/` (**package**) | Sole stage-10/11/12 coordinator and implementation of the narrow `MutationCoordinatorPort` consumed by `command`. It carries owner-sealed intents and, for authorization-bearing paths, the candidate/currentness claim into stage 10; manages caller leases and caller/controller/evidence transactions; passes the transaction's exact admission proof/context to `authzeval.FinalizeCandidateAt`; asks `evidence.RequireCurrentAllow` for the unforgeable mutation gate; submits that proof to `AcceptCurrentAllow`; obtains an owner-bound transaction permit; and only then calls each semantic owner and submits finalized changes through `ApplyFinalized`. A Deny aborts this transaction and is completed/audited only through a fresh evidence-only transaction. It then requests evidence-owned completed-evaluation proofs and the applicable carrier plan, plus idempotency-owned completed results through the plan's neutral `operation.CompletionBinding` where applicable. It accepts sealed values, commits, receives (but cannot construct) the state transaction's `CommitReceipt`, and only then passes that receipt to `authzeval.ReleaseAfterPublication` when a completed evaluation exists. It constructs or mutates no semantic value, editor, permit, receipt, descriptor, decision material, result or carrier and performs no semantic revalidation. `Coordinator` owns cancellation/admission/waiting; capabilities remain stack-local | `roleassign`, `membership`, `roledefinition`, `accessgroup`, `approval`, `privileged`, `review`, `exception`, `governanceprofile`, `grantintent`, `operation`, `state`, `idempotency`, `evidence`, `authzeval`, `decision` |
| `internal/govaccess/evidence/` (**package**) | Sole constructor/validator of `AuthorizationCandidateDescriptor`, `FinalizedAuthorizationCandidate`, `AuthorizationAdmissionProof`, sealed `CurrentAllowProof`, context-bound `CompletedAuthorizationEvaluation`, profile-specific sealed `DomainDecisionMaterial`, canonical ordered `MutationDescriptor` sets, `DomainConclusionDescriptor`, `AppliedMutationProof`, `AppliedConclusionProof`, canonical `MutationDigest`, the distinct sealed `AuthorizedMutationPlan`/`AutomaticMutationPlan`/`DomainConclusionPlan` values, deterministic carrier IDs, exact `CarrierSet`, and sealed `PreparedEvidenceChange` under consumed writer `VS0-WRITER-011`. `state` alone may request an admission proof after version revalidation; `authzeval.FinalizeCandidateAt` alone may request a finalized-candidate value after semantic revalidation; only `uow` may request `RequireCurrentAllow`; and only the `approval`/`exception` owner finalizers may request their corresponding decision material from exact immutable `model.ApprovalDecisionFacts`/`model.ExceptionDecisionFacts`. Only the exception owner finalizer may request terminal-exception descriptors: its Grant set is exactly `exceptionproposal.decided`, then `exceptiongrant.issued`; its Deny conclusion descriptor is exactly `exceptionproposal.decided`. `ValidateExceptionTerminalGrantDescriptorSet` (Grant path) and `ValidateExceptionTerminalDenyDescriptor` (Deny path) bind each to the exact immutable proposal, terminal reason, mandatory `exception-decision/v1` record material and containing operation/correlation evidence using the existing FEATURE-0013 envelope/linkage semantics. `BindAppliedChange`/`BindAppliedConclusion` consume only neutral `operation.AppliedChangeLink`/`operation.AppliedConclusionLink` plus evidence-owned descriptors/materials—never a state type. Authorized and automatic plans independently enforce all applicable F18-RD-19 mandatory domain DecisionRecords; authorization evidence is additive and never their condition. Plans expose only a neutral `operation.CompletionBinding` to idempotency. It has no editor and imports neither `state`, `authzeval`, `idempotency` nor `uow` | `decision`, `model`, `operation` |
| `internal/govaccess/policyseam/` (**package**) | **Sole** FEATURE-0017 consumer for the pinned PrivilegedAccessRequest subject; exposes `policyseam.Port`; the only package importing `internal/policyeval`; `Indeterminate`→`Deny` mapping (DD-01) | `policyeval` |
| `internal/govaccess/fixtures/` (**package**) | Immutable synthetic fixture builders + graph validator + deterministic fake executor (DD-06) | `model`, `clock` |
| `internal/govaccess/idempotency/` (**package**) | Neutral owner of `IdempotencyLookupKey`, `RequestBinding`, sealed `PreparedCompletedResult`, `CompletedRecord`, and pure completed-record lifecycle helpers only. It owns no in-flight entry/table, reservation type, `CallerReservationLease`, reservation cleanup or waiter signaling — those are owned solely by `state.Store` (DD-11; ADH-2026-074 §3). `PrepareCompletedResult` wraps the owner-supplied neutral `operation.ResultMaterial` and binds it to the exact lookup/request binding plus the evidence-validated neutral `operation.CompletionBinding`, and derives the fixed 24-hour completed-record deadline from that binding's sealed publication instant with **no second clock read**; it never imports or names a `state`, `evidence` or `clock` type. These types apply only to caller-keyed mutations; controller mutations use existing transition/version CAS and no synthetic key | `apivalid`, `model`, `operation` |
| `scripts/feature0018-architecture-check/` | Named standard-library static checker, unit tests and positive/negative fixtures for DD-13, including the `F18-ARCH-012_CONTEXT_BOUND_CHANGE` closed-boundary rule (state-only declaration of `state.ContextBoundChange`/`StateEditor`/`ApplyTo`, owner-only implementation/embedding/exposure of `state.ContextBoundChange`, each owner naming `StateEditor` only as its own `ApplyTo` parameter with no editor escape, owner-`FinalizeAt`-only permit-binding extraction, and the source-level copied-binding-wrapper rejection for which this checker is the required enforcement boundary); invoked by `make feature-0018-architecture-check` and the FEATURE-0018 feature gate | Go parser/type checker |
| `internal/govaccess/projection/` (**package**) | Progressive-disclosure safe projections (REQ-F18-24) and redacted safe-denial projection | `model`, `decision` |
| `tests/conformance/feature_0018_*.go` | Runtime cases `{VS0-CF-F18-01..37,39..49,51..53}` and inherited `VS0-CF-F01/F02/X01/X02`; 38/50/54 retain their non-runtime classifications in §7.1 | conformance harness |
| Canonical machine-readable schema entries for `VS0-SCHEMA-020..024,062..068` (**CONTRACT reference, not an owned artifact**) | The closed schema for each FEATURE-0018 kind is already supplied inline in the manifest-controlled `docs/architecture/vertical-slices/VS-000-contract-registry.yaml` (with the finalized data model / canonical contract catalog as the field-authority owner). These entries carry **no `externalRef`**, so **no `api/schemas/<kind>.json` file exists, is registered, or is created by this feature**; the `validate/` structs consume the registry entries by reference | VS-000 registry |
| Literal HTTP path ledger (**CONTRACT_ONLY**, §4.4) | Route-shape contract for the closed operation surface; no live registration | — |

### 3.3 Component boundaries and single-writer map

Each mutable path resolves to exactly one semantic writer (SEC-05, REQ-F18-02). Owned writers: the
`roleassign` package implements **VS0-WRITER-008** (`RoleAssignment.spec`/`.status`, sole
publisher of every downstream RoleAssignment effect and sole constructor of
`PreparedRoleAssignmentIntent` and `FinalizedRoleAssignmentChange`) reached through its three distinct RD-13/14/16 workflow-trigger ports
(`GrantPort`/`ActivationPort`/`ReviewRemediationPort`) and, for the RD-08 `roleassignment.revoke`
caller operation, through `roleassign.DirectRevokePort` called only by
`command.RoleAssignmentRevokeService`. The command service supplies the sealed candidate and
owner-prepared intent to `uow.MutationCoordinatorPort`; it neither publishes nor imports state.
This is not a fourth workflow trigger, route or writer, so no other package can publish; the
`membership` package implements **VS0-WRITER-022**
(`Membership.spec`/`.protected`/`.status`) and publishes Membership state — including any AccessGroup
assignment-effect expansion — entirely within itself, submitting no RoleAssignment grant intent.
Consumed writers, unchanged:
**VS0-WRITER-007** (authorized-governance-publisher) on its FEATURE-0018 registered paths
`GovernanceProfile.spec` and `RoleDefinition.spec` only; **VS0-WRITER-011**
(authorized-decision-or-audit-producer) for `DecisionRecord.record`/`AuditEvent.record`
(constructed and validated as the sole authorized opaque `PreparedEvidenceChange` by `evidence`,
then accepted by the applicable `uow`-coordinated transaction before any result release).
**VS0-WRITER-023**
(delegated-governance-admin, `ProfileAssignment.spec`) is FEATURE-0020-owned and is neither
implemented nor consumed. AccessGroup, ApprovalPolicy, ApprovalRequest, PrivilegedAccessRequest,
AccessReview, and ExceptionGrant have no registered Slice-0 writer ID; their sole-writer authority
is fixed by their F18-RD and is realized by the corresponding component above without introducing a
competing writer. The neutral `uow.Commit` physical swap (DD-01) is **not** semantic writer
authority: it only applies opaque finalized changes that their semantic owners authorized and
constructed from the transaction's exact publication context, so single-writer authority per path
is preserved even though the physical mutation is
centralized. Approval publication and downstream effect publication remain distinct atomic
operations; only the `roleassign` package may construct or publish a RoleAssignment effect.

### 3.4 Dependency direction, typed ports, acyclic prepare/commit topology, and the compile-enforced FEATURE-0017 seam

The authoritative package rows are §3.2; this section is **the one exhaustive ledger of allowed
design-owned package edges**. Every allowed import edge appears exactly once below and is reconciled
one-for-one with the concrete package dependency declared in the §3.2 "Reuses" column for that
package; there are no generic domain-owner placeholders. This ledger **governs repository-owned
package edges only** — imports of `github.com/sanjeevksaini/sovrunn/...` packages. Approved Go 1.22
standard-library imports (for example `time`, `sync`, `context`, `errors`, `fmt`, `sort`) remain
permitted everywhere and are neither listed here nor treated as ledger edges;
`F18-ARCH-005_PACKAGE_DIRECTION` classifies each import by path prefix and applies the closed-set
check only to repository-owned import paths, ignoring standard-library and `apimeta`/`apivalid`/
`apiref`/`apiproblem`/`apicond`/`apischema`/`apiconform`/`decision`/`policyeval` inherited-primitive
edges already enumerated below. Any repository-owned `internal/govaccess/*` edge not listed here is
forbidden, and `F18-ARCH-005_PACKAGE_DIRECTION` (DD-13) enforces the closed set by reference to this
ledger alone:

```text
model -> apimeta
operation -> model
clock -> (none)
idempotency -> model, apivalid, operation
evidence -> model, decision, operation
validate -> model, apivalid, apiref, apiproblem, apicond
projection -> model, decision
policyseam -> policyeval
fixtures -> model, clock
approvalreq -> model
state -> model, operation, idempotency, decision, evidence, clock
authzeval -> apimeta, model, decision, operation, evidence, state, clock
grantintent -> model, approvalreq, roleassign
roleassign -> decision, model, operation, state, evidence
membership -> decision, model, operation, state, evidence, approvalreq
roledefinition -> decision, model, operation, state, evidence
accessgroup -> model, operation, state, evidence
approval -> decision, model, operation, state, evidence, approvalreq
privileged -> decision, model, operation, state, evidence, policyseam, roleassign, approvalreq
review -> decision, model, operation, state, evidence, roleassign
exception -> decision, model, operation, state, evidence, approvalreq
governanceprofile -> decision, model, operation, state, evidence
command -> authzeval, roleassign, uow
uow -> roleassign, membership, roledefinition, accessgroup, approval, privileged, review, exception, governanceprofile, grantintent, operation, state, idempotency, evidence, authzeval, decision
root -> roleassign, membership, roledefinition, accessgroup, approval, privileged, review, exception, governanceprofile, grantintent, approvalreq, policyseam, validate, projection, fixtures, command, uow, authzeval, state, evidence, clock, decision/bundle
```

Every listed edge is exactly a concrete package edge from a §3.2 row; the nine domain-owner rows
(`roleassign`, `membership`, `roledefinition`, `accessgroup`, `approval`, `privileged`, `review`,
`exception`, `governanceprofile`) and the `grantintent`, `command`, `uow` and `root` rows are each
enumerated here in full rather than by a generic "domain owners ->" placeholder. Each domain owner's
`state`/`evidence` edge is capability-only in the exact sense defined immediately below; its
`operation`/`model`/`decision` edges are the neutral value/base-type imports of §3.2. There is no reverse edge: `evidence`/`idempotency` do not import `state`; `idempotency` does not
import `evidence`; `authzeval` does not import `idempotency`; owners do not import
`uow`, `command` or root; `state` does not import `authzeval`, `uow`, `command` or an owner; and
nothing domain-local imports root. The `root -> authzeval`, `root -> state` and `root -> evidence`
edges are composition/wiring only: root constructs the pure evaluator and post-commit result
releaser, injects the single deterministic clock into `state.Store`, and obtains the local
FEATURE-0013 adoption bundle bytes from `evidence/profiles.go` for `decision/bundle.Load`, without
invoking any transaction, evidence-preparation, publication or result method. The `authzeval -> state` edge uses only sealed neutral
dependency/currentness constructors, evidence admission-proof views and the post-commit
`CommitReceipt`; state never calls back into `authzeval`.
`state -> evidence` exists only so the sealed transaction can bind authorization admission, accept
the evidence-owned current-`Allow` proof, bind a neutral `operation.AppliedChangeLink` through
`evidence.BindAppliedChange`, consume `PreparedEvidenceChange`, and compare sealed proof linkage at
`Seal`; it does not transfer descriptor, proof, decision-material, carrier-construction or
validation authority. `BindAppliedChange` accepts no `state` type. The completion edge is
`evidence -> operation <- idempotency`: evidence creates a neutral `operation.CompletionBinding`
from its validated plan, and idempotency consumes that neutral value without importing evidence.

The domain-owner → `state` edge is deliberately capability-only. An owner may import exactly
`ParticipantClaim`, `ExpectedVersionSet`, `ParticipantOwner`, `ResultRole`,
`MutationFinalizationPermit`, `PermitBinding`, `PublicationKind`, `ContextBoundChange`, and the
`StateEditor` parameter type; it may call only `NewParticipantClaim` and `NewExpectedVersionSet`
with its own registered owner constant, and receive/use a permit only as the argument to its own
`FinalizeAt`. It names `state.StateEditor` only as the synchronous parameter of its own registered
`ApplyTo` implementation, uses only its registered direct editor methods there, and never acquires,
retains, returns, captures, embeds, wraps, adapts, or otherwise exposes the editor. It may not import
or name `Store`, any transaction, lease, reservation, mutex, root,
editor-acquisition API, receipt constructor, admission constructor, `ApplyFinalized`, `Seal`,
`Commit`, or `Abort`. Because `state.ContextBoundChange` is an exported interface that Go 1.22
cannot seal, the DD-13 `F18-ARCH-012_CONTEXT_BOUND_CHANGE` rule — not the compiler — enforces the
closed boundary: `state` alone declares the `ContextBoundChange` interface, declares and constructs
`StateEditor`, and invokes `ApplyTo`; the right to implement, embed, adapt, wrap, or expose
`ContextBoundChange`/`ApplyTo` is confined to exactly
the nine registered owner packages, each through its own concrete finalized-change type; `ApplyTo`
invocation and `StateEditor` acquisition are confined to the `state` transaction implementation; each
owner may name `StateEditor` only as its own `ApplyTo` parameter and use only its registered direct
editor methods there; and extraction of `MutationFinalizationPermit.Binding` is confined to each
matching owner `FinalizeAt`
(DD-01 rules 1–4). An owner names `ContextBoundChange` only as the interface its own concrete change
satisfies and never implements another owner's change, wraps a foreign change, or invokes `ApplyTo`.
Thus the sealed owner-finalization API is usable without granting a domain
owner state lifecycle or physical-publication authority.

The registered ownership boundaries are represented by narrow typed ports:

- originating operation owners implement their own typed `ApprovalRequirementDeriver` ports;
- `policyseam` is the only FEATURE-0017 importer and `privileged` its only wired caller;
- the three approved workflow triggers call the corresponding `roleassign` ports; direct revoke
  enters `command.RoleAssignmentRevokeService`, which calls `roleassign.DirectRevokePort` for
  semantic preparation and the acyclic `uow.MutationCoordinatorPort` for stage-10/11/12 mechanics;
- each domain owner alone constructs its publication-time-free intent and, only after receiving a
  transaction-issued owner-bound `MutationFinalizationPermit` through its `FinalizeAt`, constructs
  the immutable finalized change, canonical result and exact descriptor; the approval and exception
  owners also request their exact sealed F18-RD-19 domain-decision material from `evidence`;
- `evidence` alone constructs/validates candidate descriptors, completed-evaluation and
  current-`Allow`/applied-change proofs, domain-decision materials, mutation descriptors, digests,
  carrier plans, carrier IDs, carrier sets and `PreparedEvidenceChange`;
- `authzeval` alone constructs authorization-currentness claims and revalidates candidate-owned
  publication-time predicates; domain mutation owners revalidate only their own mutation
  predicates; `uow` merely passes the transaction's exact admission proof and context; and
- `uow` alone calls the current-`Allow` gate, state reservation/transaction APIs and
  `authzeval.ReleaseAfterPublication`, supplying the unforgeable `CommitReceipt` returned only by
  the committed state transaction; no pre-stage-10 publisher exists; and
- `state` alone declares the `state.ContextBoundChange` interface, declares and constructs
  `state.StateEditor`, and invokes `ApplyTo`; each of the nine registered owner packages implements
  `state.ContextBoundChange` (through its
  own concrete finalized-change type) and extracts its permit binding, naming `StateEditor` only as
  its own `ApplyTo` parameter and using only its registered direct editor methods there without
  acquiring, retaining, returning, capturing, embedding, wrapping, adapting, or exposing it; only the
  `state` transaction
  invokes `ApplyTo`; and no `uow`, `command`, root, or foreign package implements, embeds, wraps,
  adapts, or invokes `ContextBoundChange`/`ApplyTo`, nor copies an otherwise-valid `PermitBinding`
  into a shape-satisfying change (DD-01 rules 1–5; DD-13
  `F18-ARCH-012_CONTEXT_BOUND_CHANGE`).

These are representations of F18-RD-08/10/11/13/14/16 and the §3.3 writer ledger, not duplicate
semantic definitions. DD-13 checks the exact imports and call sites; the compiler enforces sealed
implementations and private state, except that the exported `state.ContextBoundChange` interface and
its `ApplyTo` method are a checker-enforced closed boundary (not compiler-sealed). The DD-13 checker
is the required enforcement boundary for source-level copied-binding wrappers; the runtime
permit/context/receipt validation guarantees rejection only of absent, mismatched,
foreign-transaction, reused, or context-mismatched bindings. The resulting graph is acyclic.

---

## 4. Data / API representation

### 4.1 Contract → FEATURE-0012 profile correspondence (consumed, not selected — DD-04)

F18-RD-02, requirements §3.1–§3.3, and the registered `VS0-SCHEMA-020..024,062..068`
entries are the sole profile, field, lifecycle and supporting-value authority. Design consumes
those registrations through `model/` and `validate/`; it does not restate them or select a source
on disagreement. Startup/fixture consistency validation compares the loaded registrations and
fails with `ARCHITECTURE_DECISION_REQUIRED` on any mismatch (ADH-2026-072 §4).

### 4.2 Allowed scopes and reference compatibility (consumed from F18-RD-02)

F18-RD-02/04 and REQ-F18-02/04 are the sole allowed-scope and reference-compatibility authority.
`validate/` consumes their versioned registrations generically, never expands a set, and treats
compatibility as validation rather than permission. The exact test vectors come from the canonical
REQ/AC and conformance ledgers; this design adds no scope or containment rule.

### 4.3 Enums and typed discriminators

The closed enums of requirements §3.3 (`MembershipContextKind`, `MembershipType`,
`PrivilegedAccessMode`, `ApprovalDecision`, `ExceptionGrantEffect`, `AssignmentValidityMode`,
`ActionTargetMode`) are declared as Go string types with a `Valid()` method and an exhaustive
`All…()` accessor, mirroring the `policyeval.Outcome`/`NonResult` and `decision.FailPosture`
convention. `ActionTargetBinding` is a discriminated union (`ExactResource`/`CreateParent`/
`ScopeOnly`) with exactly one active variant; a mixed/missing/empty/wildcard/unknown/ambiguous
binding fails closed (REQ-F18-06). No enum value is added, widened, or made extensible.

### 4.4 Literal HTTP path contract ledger (CONTRACT_ONLY/NO_TASK — DD-07)

Route grammar (FEATURE-0012 §6.1; Go 1.22 `http.ServeMux` method/path patterns; ADH-2026-053
item-action form `/{uid}/actions/<action>` read via `r.PathValue("uid")`, plus the collection-action
form `/actions/<action>` that carries no `{uid}` segment). No generic `PUT`, unrestricted `PATCH`,
or hard `DELETE` (REQ-F18-02):

```text
Collection       : POST | GET   /apis/<group>/v1alpha1/<plural-kebab>
Item             : GET          /apis/<group>/v1alpha1/<plural-kebab>/{uid}
Draft PATCH      : PATCH        /apis/<group>/v1alpha1/<plural-kebab>/{uid}   (Draft only, VersionedDefinition kinds)
Item action      : POST         /apis/<group>/v1alpha1/<plural-kebab>/{uid}/actions/<action>
Collection action: POST         /apis/<group>/v1alpha1/<plural-kebab>/actions/<action>
```

The **Item action** form carries a complete `{uid}` `http.ServeMux` wildcard segment (read via
`r.PathValue("uid")`) and targets one existing resource. The **Collection action** form carries **no
`{uid}` segment**: it targets the collection itself and is authorized against the operation's
registered `CreateParent.parentScopeKind` (F18-RD-06), not against an existing item. The only
FEATURE-0018 collection action is `exceptiongrant.propose` at
`/apis/governance.sovrunn.io/v1alpha1/exception-grants/actions/propose`; it is a distinct grammar
category from every UID-bearing item action and is never placed in the item-action column. No other
collection action, and no generic `PUT`, unrestricted `PATCH`, or hard `DELETE`, exists
(REQ-F18-02).

The `<group>` token is the literal FEATURE-0012 domain-group **owned by the VS-000 contract
registry** (`iam.sovrunn.io` or `governance.sovrunn.io`) for each kind, resolved below from the
registry `identity` field of the routed kinds `VS0-SCHEMA-021..024,062..067` (`VS0-SCHEMA-020`
PrincipalRef and `VS0-SCHEMA-068` supporting values are `EmbeddedValue` with no route). This
ledger adopts only the domain-grouping rule and binds to that exact registration; it invents no
group token. **Every row is transcribed one-to-one from the ADH-2026-070 Appendix B closed v1
operation-and-mutability surface (F18-RD-02) and the F18-RD-06 discriminated `ActionTargetBinding`
registrations** — no caller operation in that surface is dropped, renamed, or reclassified, and
no controller-only effect is promoted to a caller operation. Each caller action verb below maps
to the exact Appendix B caller operation named in the mapping column; the literal path shape is
the only design-owned mechanic. **Owning-controller-only effects are excluded** (Appendix B right
column) and are never caller-visible: RoleAssignment specification/lifecycle/status/`replace`/
`expiry`/`review-due` publication; privileged `activate`/`block`/`expire`; JIT/break-glass
RoleAssignment materialization; RoleDefinition/ApprovalPolicy/GovernanceProfile automatic
supersession linkage; ApprovalRequest creation, `expire`, and result publication; AccessReview
snapshot capture, readiness publication, `complete`, and expiry; ExceptionGrant immutable
grant/revocation issuance and expiry-audit/projection refresh; and every automatic
authoritative-time transition. The caller may *authorize and submit an exact immutable intent*
for a controller-published effect (Appendix B: RoleAssignment `revoke`, review `remediate`,
privileged `revoke`/`admin-revoke`), but the sole publisher remains the owning controller.
**`CreateParent.parentScopeKind` is the authorization target scope of the submitting operation
(F18-RD-06); it does *not* mean the caller's route creates an ApprovalRequest.** An
`ActionTargetBinding` of `CreateParent` authorizes the submission against the exact registered
`parentScopeKind`; whether any ApprovalRequest is later materialized is a controller decision
taken only when the originating controller's derived `ApprovalRequirement` is `Required` — and
even then the controller-only ApprovalRequest creator materializes it, never the caller route.
`PrincipalRef` (`VS0-SCHEMA-020`, `EmbeddedValue`) has no route; it is an operation-local value.
No route is registered by this feature and no controller effect becomes a caller permission.

| Kind (group) | Collection POST / GET | Item GET | Draft PATCH (update) | Item actions (UID-bearing, POST `/{uid}/actions/<action>`) | Appendix B caller operations covered |
|---|---|---|---|---|---|
| AccessGroup (`iam.sovrunn.io`) | POST/GET `/apis/iam.sovrunn.io/v1alpha1/access-groups` | `…/{uid}` | — | `update`, `suspend`, `resume`, `retire` | create, get, list, **update permitted owner/display intent**, suspend, resume, retire |
| Membership (`iam.sovrunn.io`) | POST/GET `/apis/iam.sovrunn.io/v1alpha1/memberships` | `…/{uid}` | — | `synchronize`, `suspend`, `reactivate`, `revoke` | create/provision, get, list, synchronize, suspend, **reactivate**, revoke |
| RoleDefinition (`iam.sovrunn.io`) | POST/GET `/apis/iam.sovrunn.io/v1alpha1/role-definitions` | `…/{uid}` | Draft only (update Draft) | `publish`, `suspend`, `restore`, `retire` | create Draft, get, list, update Draft, publish, suspend, restore, retire |
| RoleAssignment (`iam.sovrunn.io`) | POST/GET `/apis/iam.sovrunn.io/v1alpha1/role-assignments` | `…/{uid}` | — | `revoke` (submit exact revoke intent; `roleassign` publishes) | **collection POST is the one `roleassignment.grant` entry point** — accepts either an authorized direct non-delegable grant intent or one `RoleAssignmentProposal` (`grantintent` accepts; `roleassign` publishes); get, list, revoke |
| PrivilegedAccessRequest (`iam.sovrunn.io`) | POST/GET `/apis/iam.sovrunn.io/v1alpha1/privileged-access-requests` | `…/{uid}` | — | `cancel`, `revoke`, `admin-revoke` | submit, get, list, cancel own, **revoke own**, **administratively revoke (distinct privileged action)** |
| AccessReview (`iam.sovrunn.io`) | POST/GET `/apis/iam.sovrunn.io/v1alpha1/access-reviews` | `…/{uid}` | — | `start`, `decide`, `remediate`, `cancel` | create authorized campaign, get, list, **start**, decide item, remediate, cancel |
| ApprovalPolicy (`iam.sovrunn.io`) | POST/GET `/apis/iam.sovrunn.io/v1alpha1/approval-policies` | `…/{uid}` | Draft only (update Draft) | `publish`, `retire` | create Draft, get, list, update Draft, publish, retire |
| ApprovalRequest (`iam.sovrunn.io`) | GET only `/apis/iam.sovrunn.io/v1alpha1/approval-requests` (controller is sole creator) | `…/{uid}` | — | `decide`, `cancel`, `admin-cancel` | get, list, decide, cancel originating pending intent, administratively cancel (distinct privileged action) |
| GovernanceProfile (`governance.sovrunn.io`) | POST/GET `/apis/governance.sovrunn.io/v1alpha1/governance-profiles` | `…/{uid}` | Draft only (update Draft) | `publish`, `retire` | create Draft, get, list, update Draft, publish, retire |
| ExceptionGrant (`governance.sovrunn.io`) | GET only `/apis/governance.sovrunn.io/v1alpha1/exception-grants` (append-only; controller issues evidence) | `…/{uid}` | — | `revoke` (submit exact revoke intent; controller issues linked immutable evidence) — *`propose` is a **collection action**, not a UID-bearing item action; see the collection-action ledger below* | propose *(collection action — see below)*, get, list, revoke |

**Collection-action ledger (distinct grammar category; no `{uid}` segment).** Exactly one
FEATURE-0018 caller operation uses the collection-action form; it is represented here separately
from every UID-bearing item action and is authorized against its registered
`CreateParent.parentScopeKind` (F18-RD-06), never against an existing item:

| Kind (group) | Collection action (POST `/actions/<action>`) | Appendix B caller operation | Authorization target |
|---|---|---|---|
| ExceptionGrant (`governance.sovrunn.io`) | `propose` → POST `/apis/governance.sovrunn.io/v1alpha1/exception-grants/actions/propose` (`exceptiongrant.propose`; submit the immutable `ExceptionProposal`) | propose | registered `CreateParent.parentScopeKind` |

No other kind registers a collection action; `RoleAssignment` proposal submission uses its own
collection `POST` (the `roleassignment.grant` entry point), not a collection action. This
collection-action row does not change the approved operation, path, ownership, or CONTRACT_ONLY
classification of `exceptiongrant.propose`; it only classifies it in its correct grammar category.

**Proposal submission uses the two owning-collection routes; there is no
`/approval-requests/actions/submit-*` route.** The prior draft's
`/approval-requests/actions/submit-exception-proposal` and
`/approval-requests/actions/submit-role-assignment-proposal` routes are **removed**, together with
all prose claiming that a `CreateParent` binding means the caller route creates an ApprovalRequest.
ApprovalRequest remains controller-created only (its collection is `GET` only). The two proposal
submissions are instead caller operations on their **own** owning collections, authorized against
their exact registered `parentScopeKind` (F18-RD-06):

- **RoleAssignment grant / proposal** — the single `RoleAssignment` collection POST is the one
  `roleassignment.grant` entry point (§4.4 table). It accepts **either** an authorized direct
  non-delegable grant intent **or** exactly one `RoleAssignmentProposal`. `grantintent` is the sole
  acceptor and derives the `ApprovalRequirement`; if the derived requirement is `Required`, the
  originating controller may later ask the **controller-only ApprovalRequest creator** to
  materialize an ApprovalRequest — that materialization is not caller route ownership and is not a
  separate caller route.

  ```text
  POST /apis/iam.sovrunn.io/v1alpha1/role-assignments
    (roleassignment.grant: accepts an authorized direct non-delegable grant intent OR one
     RoleAssignmentProposal; grantintent accepts, roleassign publishes; ApprovalRequest, if any,
     is materialized only by the controller-only ApprovalRequest creator when derived
     ApprovalRequirement == Required — never by this route)
  ```

- **ExceptionGrant proposal** — `exceptiongrant.propose` remains owned by ExceptionGrant and uses
  **one collection-action path under the ExceptionGrant collection**, authorized against the exact
  registered `parentScopeKind`:

  ```text
  POST /apis/governance.sovrunn.io/v1alpha1/exception-grants/actions/propose
    (exceptiongrant.propose: submit the immutable ExceptionProposal; authorized against the
     registered CreateParent parentScopeKind. The ExceptionGrant is issued only by the exception
     controller on terminal Grant; if derived ApprovalRequirement == Required the controller-only
     ApprovalRequest creator may later materialize a request — never this route)
  ```

Notes preserving the closed Appendix B / F18-RD-02 surface:

- **AccessGroup `update`** is the Appendix B "update permitted owner/display intent" caller
  operation — a constrained item action that may set only the responsible-owner reference and
  display metadata (owner must resolve to a current Human at every privilege-increasing
  operation, REQ-F18-04); it is not a generic `PUT`/`PATCH` and cannot alter scope, membership,
  or lifecycle.
- **Membership `reactivate`** (Suspended → Active, Appendix B) restores the exact suspended
  relationship and, when it makes a group-held assignment newly/longer effective, must pass the
  membership-enabled grant ceiling (REQ-F18-04/09). It is named `reactivate`, matching Appendix B
  exactly (not `resume`); `synchronize` updates only system-owned freshness/provenance.
- **PrivilegedAccessRequest** has three caller actions: `cancel` (caller-initiated early stop of
  own request before/at approval), `revoke` (caller revokes own active grant), and `admin-revoke`
  (the distinct privileged administrative revoke, Appendix B "administratively revoke under the
  distinct privileged action"). Each revoke submits an exact revoke intent; the `privileged` package
  submits it to `roleassign.ActivationPort` (the RD-14 trigger), the `roleassign`
  controller publishes the RoleAssignment revocation, and the privileged controller records the
  request transition. `activate`/`block`/`expire` and the `Active`/`Expired`/`Revoked-by-time`
  transitions remain controller-owned and are not caller actions.
- **AccessReview `start`** is the Appendix B caller operation that moves a `Pending` campaign into
  `InProgress`; the immutable snapshot is captured by the controller on that entry (REQ-F18-16).
  `decide` records an item disposition, `remediate` submits Revoke/Replace/due-advancement intents
  (through the `review` package to `roleassign.ReviewRemediationPort`, the RD-16 trigger — never
  through `grantintent`) for the `roleassign` controller to publish, and `complete`/`ExpiredIncomplete`
  remain controller-derived.
- **RoleAssignment** exposes only `revoke` as a caller item action (submit intent); the collection
  POST is the `roleassignment.grant` entry point; `replace`, `expiry`, and `review-due` are
  controller-owned review-remediation effects, never caller actions. The direct caller `revoke`
  (`roleassignment.revoke`, RD-08) is handled by the in-process caller ingress
  `command.RoleAssignmentRevokeService`, which evaluates the caller operation, asks
  `roleassign.DirectRevokePort` for the owner-sealed intent, and delegates stage-10/11/12 mechanics
  through `uow.MutationCoordinatorPort`. It is **not** one of the three RD-13/RD-14/RD-16 workflow
  triggers. `roleassign` remains the sole semantic preparer/publisher; the command ingress adds no
  route and no second writer.
- **ApprovalRequest** collection is `GET` only and is never caller-created by any route (no
  `submit-*` collection action exists). `admin-cancel` (ApprovalRequest) is a privileged
  administrative caller item action matching the Appendix B lifecycle. **ExceptionGrant** collection
  is `GET` only for reads; its `propose` collection action (`exceptiongrant.propose`) and item
  `revoke` are the only caller mutating operations, and the immutable evidence is issued solely by
  the exception controller.

This ledger is CONTRACT_ONLY: it produces no task and no source in this feature, registers no
live route, and adds no action, authority, or effect beyond the Appendix B / F18-RD-02/06 caller
operation surface.

### 4.5 Field ownership and system-owned fields

Clients never author `status`, identity/version/timestamp fields, `metadata.scopeRef`
(server-derived), evidence/provenance, correlation IDs, or any field owned by another feature
(ServiceRegion, `providerSelectionModes`, ExecutionTarget `status.*`, ProfileAssignment) —
these are rejected in structural/field-policy validation (FEATURE-0012 §6.8; Go guardrails
§10). Progressive disclosure (REQ-F18-24) is realized by `projection.go` returning
risk-proportional views over the same protected contracts; it adds no facade resource and
weakens no underlying contract.

---

## 5. Validation and deterministic error behavior

### 5.1 Validation precedence (exact F18-RD-21 order, realized over inherited stages)

`validate/` and the controllers together realize the F18-RD-21 precedence, layered on the
FEATURE-0012 ordered validation stages (§6.10) and the `internal/apivalid` pipeline:

1. bounded transport (inherited body/size/content-type limits)
2. authentication (already-authenticated `PrincipalRef`; absent auth → inherited VS0-F01)
3. structural validation (strict decode, unknown/duplicate/future-owned field rejection)
4. local semantic / cross-field validation
5. safe scope/reference resolution (safe denial precedes disclosure)
6. operation-required assurance + contextually required current membership
7. assignment/action/scope/resource/delegation + membership-enabled-grant-envelope evaluation
8. operation-required `EligibilityRef` + separation-of-duties evaluation
9. approval-requirement / policy / approval / exception evidence evaluation
10. concurrency + idempotency
11. required decision/audit-obligation evidence
12. publication

Stages 1–7's structural/reference/scope work, stage-8 eligibility/separation-of-duties, and
stage-9 approval/policy/exception **evaluation** are pure functions of the immutable snapshot,
request and injected deterministic clock; unknown, deferred and future-owned fields reject at
stage 3. Stages 7–9 produce only one sealed internal `OperationCandidate`. It is not a completed
authorization evaluation, `AuthorizationResult`, audit evidence, bearer value, or downstream
authority (ADH-2026-072 §2).

Stage 10 is completed before any stage-11 carrier construction, validation, acceptance or
publication can fail:

- a caller-keyed mutation resolves binding conflict, replay, wait or reservation ownership and
  revalidates both its mutation predicates and the authorization candidate's complete evaluated
  snapshot/dependency-version set;
- a controller-owned mutation performs its exact transition/version CAS and resolves a repeated
  effective-once operation without inventing a caller key; an authorization-bearing controller
  trigger also revalidates the candidate's complete dependency-version set, while expiry,
  reconciliation and another non-assignable automatic effect carries no candidate;
- a standalone read or a terminal `Allow`/`Deny` conclusion performs a coherent-snapshot check;
  the check covers every dependency version recorded by `authzeval`, and a change restarts
  evaluation rather than finalizing stale evidence; and
- a replay is selected only after fresh authentication, authorization and safe-visibility checks.

For a successful authorization-bearing begin, stage 10 returns a sealed admission proof linked to
the revalidated dependency set and captures the one authoritative publication context. `uow`
passes both values to `authzeval.FinalizeCandidateAt`; `authzeval`, not `uow`, revalidates the
candidate-owned publication-time predicates. For every successful mutation begin, `uow` invokes
each mutation semantic owner's exact-context finalizer and transaction application before any
mutation digest, completed result, carrier ID or carrier set is derived. A finalizer/application
returns either a finalized change, an already-registered domain outcome, or a mechanical failure;
none publishes state. After all path-applicable stage-10 checks and semantic-owner publication-time
revalidations succeed, `evidence.CompleteAuthorizationCandidate` binds the authzeval-finalized
candidate to the exact admission proof and context. This successful constructor call is the
candidate's single transition into a sealed `CompletedAuthorizationEvaluation`; it is still not
externally releasable authority. Automatic mutation-only effects skip candidate completion and
carry only their required mutation evidence.

A registered stage-10 conflict is returned immediately; the candidate is discarded,
so a later evidence failure cannot mask it. Only an applicable stage-10 success enters stage 11.
There, `evidence` constructs and validates the exact required carrier set from that completed
evaluation and the applicable sealed transaction accepts it. Stage 12 atomically commits the accepted carrier set together with any
domain mutation and caller completed record, then—and only then—`uow` may call
`authzeval.ReleaseAfterPublication` to release the operation-local `AuthorizationResult`. If stage
11 fails, the operation returns the inherited safe internal failure with no result, mutation,
completed replay or downstream authority. No pre-stage-10 authorization audit is published.

The path-specific mechanics, including denied/read/replay flows, are defined once in §5.4–§5.5.
Safe denial (stage 5) still precedes detailed inaccessible-reference disclosure (PRIV-01), and no
path reorders the authoritative twelve stages.

### 5.2 Error and outcome mapping (inherited codes only; per operation/case)

Design maps only the outcomes the requirements' Exact conformance semantics ledger already
fixes to an inherited FEATURE-0012 code, per operation/case; it introduces no new top-level
Problem code, violation code, HTTP status, or precedence rule, and it selects **no** code for
any case the requirements leave intentionally unspecified (ADH-2026-071). Rows below cite the
governing case; where the ledger records a case as an intentionally unspecified existing safe
error, the design records that fact and chooses nothing:

| Operation / case | Representation (as fixed by the requirements ledger) |
|---|---|
| Absent/invalid authentication (any owned route) | inherited `AUTH_REQUIRED`/401 (VS0-F01 → VS0-CF-F01); no mutation/lifecycle/idempotency/audit |
| Safe denial of an inaccessible reference | inherited `RESOURCE_NOT_FOUND`/404 (VS0-F02 → VS0-CF-F02); no existence disclosure |
| Guest RoleAssignment omitting a required TimeBound timestamp — the guest TimeBound violation (`VS0-CF-F18-32`; REQ-F18-05/08; a Guest assignment must be TimeBound with both timestamps and may not exceed the pinned guest Membership `expiresAt`; not a Membership failure) | `VALIDATION_FAILED`/422 with RFC 6901 pointer (the one structural case the ledger fixes) |
| Qualifier-only ExecutionTarget retire (`VS0-CF-F18-11`) | `AUTHORIZATION_DENIED`/403 |
| Authorization/self-approval/eligibility denial where the ledger registers `AUTHORIZATION_DENIED` (e.g. `VS0-CF-F18-14`, `VS0-CF-F18-15`; `VS0-CF-F18-41` as `none or AUTHORIZATION_DENIED`) | `AUTHORIZATION_DENIED`/403 exactly as that case registers |
| Safe replay / terminal-CAS loss / overlap where the ledger registers `CONFLICT` (e.g. `VS0-CF-F18-17`; same caller key with a different target, scope or digest binding) | `CONFLICT`/409 |
| Required audit-obligation failure (F18-RD-20) | `INTERNAL_ERROR`/500; publishes no state change or completed replay result |
| Domain adjudication (not an error) | terminal domain `Deny` / `RequiresApproval` / exception `Grant`\|`Deny` (never a Problem code) |
| Review beneficiary/reviewer conflict (`VS0-CF-F18-45`) | `violations[].code` `REVIEWER_CONFLICT` / `REVIEW_BENEFICIARY_UNRESOLVED` (never a top-level code) |
| `VS0-CF-F18-12` (ExecutionTarget create-only granularity limitation) | **code-agnostic safe-error oracle** — `success == false` and outcome class is `SafeError`; no unauthorized disclosure, mutation/publication/idempotency completion; audit follows registered authority; no assertion of Problem code, `violations[].code`, HTTP status, or precedence |
| `VS0-CF-F18-47` (membership/relationship/role evidence to manufacture eligibility) | **code-agnostic safe-error oracle** — `success == false`, outcome class `SafeError`, no unauthorized disclosure/effect, audit per registered authority; no code/channel/status/precedence assertion |
| `VS0-CF-F18-48` (empty/indirect/stale/incompatible reviewer list) | **code-agnostic safe-error oracle** — `success == false`, outcome class `SafeError`, no unauthorized disclosure/effect, audit per registered authority; no code/channel/status/precedence assertion |
| `VS0-CF-F18-49` (grantor/rule publisher makes a role/group/relationship qualify) | **code-agnostic safe-error oracle** — `success == false`, outcome class `SafeError`, no unauthorized disclosure/effect, audit per registered authority; no code/channel/status/precedence assertion |

The four cases above are realized with a **code-agnostic safe-error oracle**: each test first asserts
`success == false` and the generic outcome class `SafeError`, then asserts that the operation
terminates safely, discloses nothing to an unauthorized caller, performs no
mutation, publication, or idempotency completion, and produces the audit behavior its registered
authority requires — and the test **must not assert or choose** a Problem `code`, a `violations[].code`
channel, an HTTP status, or any precedence. A successful no-op is a test failure. They are therefore **not** classified as
`VALIDATION_FAILED`, `AUTHORIZATION_DENIED`, a `violations[].code` value, an HTTP status, or a
precedence — design leaves the exact error selection exactly as the requirements ledger records it
(ADH-2026-071 authorizes no such choice). All other structural/prohibited/eligibility failures that
the ledger does not fix to a specific code likewise inherit the FEATURE-0012 layered-validation
outcome without a design-chosen code. HTTP status is transport metadata, never a `code` value. No
`429`/quota code exists (quota is out of FEATURE-0018 scope).

### 5.3 Fail-closed and determinism

The observable rules are owned by requirements SEC-01/02, OPS-04 and REQ-F18-09/10/20. Design
enforces them through pure closed-enum validators, immutable snapshots, the injected clock,
code-agnostic safe-error oracles and the sealed publication protocol. No additional denial,
precedence or time semantic is selected here.

<!-- END RETAINED SEMANTIC EXCERPT: raw design.md lines 904-1442 -->

<!-- BEGIN RETAINED SEMANTIC EXCERPT: raw design.md lines 2591-2830 (§5.5 audited release, §6 security/privacy/observability/compat, FEATURE-0013 Adoption, §7.1 local conformance) -->
### 5.5 Provisional authorization and post-stage-10 audited release (DD-12; F18-RD-20 clause 3)

`authzservice.go` is an orchestration boundary, not an early publisher. Its pure evaluator runs
stages 7–9 over the immutable snapshot and returns one sealed internal `OperationCandidate`. The
candidate carries an evidence-owned `AuthorizationCandidateDescriptor`, but it is not a completed
authorization evaluation, is not returned to a caller, cannot authorize downstream work, and
cannot be converted to `AuthorizationResult` outside `uow`.

`uow` then follows exactly one path:

1. **Failure before an authorization candidate exists.** Bounded transport, authentication,
   structural, safe-resolution or other registered pre-candidate failure returns only its existing
   safe outcome. It creates no candidate, completed evaluation or `AuthorizationResult`; any
   separately registered mutation-attempt audit obligation remains governed by F18-RD-20.
2. **Standalone read or candidate Deny.** `BeginEvidenceTransaction(currentnessClaim)`
   performs the stage-10 complete dependency-version check. A mismatch discards the candidate and
   restarts evaluation. On success, `uow` passes the proof/context to
   `authzeval.FinalizeCandidateAt` and then invokes `evidence.CompleteAuthorizationCandidate` with
   the sealed current candidate to complete the evaluation; stage 11 prepares and the transaction
   accepts (`AcceptPreparedEvidence`) the exact authorization carrier set, then commits. Only after
   the transaction's successful `Commit` returns its sealed `state.CommitReceipt` may `uow` invoke
   `authzeval.ReleaseAfterPublication(completedEvaluation, commitReceipt)` — which constructs and
   releases the `AuthorizationResult` only, taking no `operation.ResultMaterial` — and return the
   redacted Deny or protected Allow projection. This standalone read/Deny path stages and commits no
   caller mutation, domain state, `operation.ResultMaterial` or completed idempotency record.
3. **Caller mutation.** Stage 10 returns Conflict, Wait, Replay or Owner. Conflict returns before
   completed-evaluation/evidence preparation and discards the candidate. Wait/retry also discards it
   and restarts after the signal. **Replay** reaches an evidence-only transaction only after fresh
   reauthentication, reauthorization and safe-visibility rechecks; `authzeval` finalizes the current
   candidate, `evidence.CompleteAuthorizationCandidate` completes the evaluation over the current
   request, the transaction accepts and commits exactly the authorization carrier set, and only after
   the commit receipt is returned does `uow` release the `AuthorizationResult` through
   `authzeval.ReleaseAfterPublication(completedEvaluation, commitReceipt)`. `uow` then returns the
   existing successful result recorded on the state-provided `CompletedRecord` (selected through the
   state reservation/replay API) via the unit of work; because the stored result is read from the
   `state`-owned `CompletedRecord` and `AuthorizationResult` is constructed by `authzeval` from the
   completed evaluation and commit receipt alone, replay introduces no `authzeval -> idempotency`
   import and passes no `idempotency`-owned value to the release API. **Owner** begins the caller
   transaction without descriptors, obtains its context, asks `authzeval` to finalize the admitted
   candidate, then finalizes and applies every mutation-owner intent using that exact context. The
   originating owner's `FinalizeAt` produces the single owner-finalized neutral
   `operation.ResultMaterial`; that result material is consumed **only** by
   `idempotency.PrepareCompletedResult` for the caller's completed idempotency record (staged through
   `StageCompletedResult`) and is never passed to `authzeval.ReleaseAfterPublication`. `uow` completes
   the candidate only after all applicable stage-10 checks succeed, derives and stages the exact
   completed result and carrier set, commits atomically, and only after the commit receipt is returned
   separately releases the `AuthorizationResult` through
   `authzeval.ReleaseAfterPublication(completedEvaluation, commitReceipt)`.
4. **Authorization-bearing controller mutation.** Stage 10 returns Conflict, effective-once NoOp or
   a controller transaction using exact transition/version and authorization-dependency checks. No
   synthetic caller key or completed replay record is created. Conflict/NoOp discard the candidate
   and create no completed evaluation or `AuthorizationResult`. A successful transaction captures
   its context, asks `authzeval` to finalize the candidate, and only then permits the mutation owner
   to finalize/apply its change. The completed evaluation plus mutation plan determines the exact
   carrier set; the transaction accepts and commits it, and only after the commit receipt is returned
   does `uow` separately release the `AuthorizationResult` through
   `authzeval.ReleaseAfterPublication(completedEvaluation, commitReceipt)` (taking no
   `operation.ResultMaterial`). When the owner finalizer instead
   reaches a publishable terminal domain conclusion that creates no resource — the ExceptionProposal
   `Deny` case (F18-RD-17/19) — the finalized change is a `DomainConclusion`: `ApplyFinalized`
   consumes the permit and records an empty-delta conclusion receipt, `uow` builds the exact
   `evidence.NewDomainConclusionPlan`/`PrepareDomainConclusionCarrierSet` carrier set (the mandatory
   `exception-decision/v1` DecisionRecord, the terminal-decision `AuditEvent`, and the
   `authorization.evaluated` event), and commit atomically publishes that evidence while creating no
   ExceptionGrant, before the same post-commit `AuthorizationResult` release. A non-publication
   `Domain(D)` outcome
   instead returns its existing mapping and publishes nothing.
5. **Automatic mutation-only controller effect.** Expiry, reconciliation and another closed
   non-assignable automatic effect enters stage 10 without an `OperationCandidate` or currentness
   claim, so it has **no completed authorization evaluation and no `AuthorizationResult`** unless
   existing approved semantics independently require one. Conflict/effective-once NoOp returns its
   existing outcome. Success finalizes and applies
   only the mutation-owner change, prepares exactly its F18-RD-20-required mutation carriers, and
   commits. It never invokes `authzeval.FinalizeCandidateAt`,
   `evidence.CompleteAuthorizationCandidate` or `authzeval.ReleaseAfterPublication`, and emits no
   `authorization.evaluated` event or authorization DecisionRecord.

All successful stage-10 transaction starts capture the authoritative deterministic publication
instant described in §5.4. Each mutation semantic owner performs only its own publication-time
recheck; `authzeval` performs candidate-owned publication-time revalidation. State verifies only
version/currentness linkage and permit provenance, and `uow` only coordinates. Non-publication
`Domain(D)` finalization outcomes retain their registered denial, deadline, expiry, validity,
lifecycle or conflict mapping and publish nothing; a publishable terminal `DomainConclusion` instead
publishes exactly its mandatory domain DecisionRecord and required AuditEvents with no resource
mutation and no ExceptionGrant.
The successful `evidence.CompleteAuthorizationCandidate` call—after those stage-10 checks and only
on an authorization-bearing path—is the candidate-to-completed-evaluation transition. An applicable
invariant, evidence-construction, evidence-validation or evidence-acceptance failure returns the
inherited `INTERNAL_ERROR`/500 with no `AuthorizationResult`, mutation, completed result or
downstream authority. `Commit` has two distinct phases (contract
`transactionModes.stateMachine.*.commit` and `goApiSurface.namedClosedResultTypes`
`state.CommitOutcome`): a **fallible pre-install validation phase** and a **non-fallible terminal
install phase**. In the fallible pre-install phase `Commit` returns the `state.CommitOutcome`
`MechanicalFailure` variant — installing no shadow root, carrier set, completed result or
publication sequence — when invoked before a successful `Seal` (no `Sealed` marker recorded), when a
still-required accepted-carrier or completed-result linkage is absent, or when it detects an internal
invariant breach; the deferred idempotent `Abort` then unwinds the transaction. Only after that
validation passes does the **terminal install** run, and that terminal section alone is non-fallible:
it installs exactly one shadow root, the carrier set and the revalidated sequence under `stateMu` and
returns the `Committed` variant with the opaque `state.CommitReceipt`. `Commit` is therefore an
error source **only** in its pre-install validation phase and never during terminal install. A
conflict already selected at stage 10 never enters the evidence-failure path.

FEATURE-0013 supplies only carrier types, profile/bundle rules and deterministic validation;
FEATURE-0018's sealed in-memory transaction is its own non-durable conformance machinery. The root
may wire constructors but cannot invoke transaction, evidence or release methods. `authz.go` stays
pure, `evidence` stays the sole carrier-set constructor, and `uow` stays the sole caller of evidence
preparation, state acceptance and `authzeval.ReleaseAfterPublication`.

---

## 6. Security, privacy, observability, compatibility, and versioning

### 6.1 Security (realizes SEC-01..09)

Requirements §6 and architecture F18-RD-08..21 remain the security authority. The design-owned
enforcement mechanisms are: closed value types and deterministic validators (§4–§5), immutable
snapshot evaluation (DD-08), owner-private prepared-intent/context-finalized-change pairs and typed ports (§3.2–§3.4), the
two-lock sealed publication protocol (§5.4), audited authorization service (§5.5), and the DD-13
checker. The conformance ledger, not this section, defines each observable security outcome.

### 6.2 Privacy and redaction (PRIV-01..03)

PRIV-01..03 own the observable rules. `projection/`, §5.2 safe-error oracles and FEATURE-0013's
canonical projection/handling vocabulary implement them without defining a second redaction
policy. Tests assert absence of forbidden material rather than copying the forbidden-value list.

### 6.3 Observability (Go guardrails §19–21; observability-and-audit baseline)

The Go guardrails and observability/audit baseline remain authoritative. Components expose only
their required redacted structured fields for deterministic tests; FEATURE-0018 creates no
logging, metrics or tracing runtime. Decision/audit history uses the adopted FEATURE-0013 carriers.

### 6.4 Compatibility and FEATURE-0013 adoption (COMPAT-01..05)

COMPAT-01..05 and §1.2 identify the exact inherited contracts. Design imports those contracts and
adds no compatibility shim or dual authority. FEATURE-0013 wiring is specified once in
`## FEATURE-0013 Adoption`; the ExecutionTarget limitation remains a cited canonical constraint,
not a design decision.

### 6.5 Versioning

FEATURE-0012 §6.15 and REQ-F18-07/12/18 own version semantics. `model` and `validate` consume the
registered version and immutability fields exactly; design introduces no migration or compatibility
rule.

### 6.6 Migration impact

None. No durable store, live route, or runtime object is created; the deliverable is a
deterministic in-memory contract library plus local conformance. Any future durable/transport
adoption is a separately approved feature that must conform to these contracts (FEATURE-0012
§13; controller-reconciliation Phase 2R note).

---

## FEATURE-0013 Adoption

FEATURE-0018 adopts FEATURE-0013 in the roles `DECISION_PRODUCER`, `AUDIT_PRODUCER`,
`DECISION_CONSUMER`, and `OPERATION_OWNER`. It registers exactly these profiles and no fourth:

```text
authorization-decision/v1
approval-decision/v1
exception-decision/v1
```

`internal/govaccess/evidence/profiles.go` is the sole FEATURE-0018 owner of deterministic strict
local `DecisionProfileBundle` bytes containing those three existing FEATURE-0013 carrier values
with the exact F18-RD-19 semantics. Root composition is the sole registration/composition point: it
passes those local bytes to the existing `internal/decision/bundle.Load` path and distributes the
resulting read-only `BundleView`. Root does not define, modify or directly publish a profile;
`evidence` does not create a second registry; and no global mutable registry, network lookup,
fallback version or competing envelope is introduced.

The registration tests must prove:

1. strict local load succeeds with exactly the three `(id, version)` entries;
2. each exact lookup returns the F18-RD-19 authority, input/result, validity, obligation,
   projection and audit-linkage semantics;
3. duplicate `(id, version)`, unknown version, malformed profile or missing required registry entry
   fails through existing FEATURE-0013 validation;
4. a fourth profile is absent and AccessReview remains a LongRunningOperation; and
5. authorization, approval and exception evidence construction pins the corresponding loaded
   profile version and never synthesizes or falls back to another profile.

This adoption consumes FEATURE-0013 types, bundle loading and validation unchanged. Transaction-
local acceptance remains FEATURE-0018 conformance machinery under DD-10/§5.4 and is not registered
as a FEATURE-0013 service.

---

## 7. Test and conformance strategy

### 7.1 Local conformance (REQ-F18-22; VS0-CF-F18-01..54)

Each active runtime `VS0-CF-F18-*` case in the ranges `{01..37, 39..49, 51..53}` is realized as a
deterministic table-driven test in `tests/conformance/feature_0018_*.go` following the
`feature_0016_test.go` precedent, using the immutable fixtures (DD-06) and injected `FixedClock`
(DD-05). **`VS0-CF-F18-54` is explicitly excluded from runtime Go tests** (proof classification:
non-runtime standards architecture gate, REQ-F18-23): its exact proof artifact is the completed
FEATURE-0018 standards-mapping matrix at
`docs/reviews/architecture-readiness/FEATURE-0018-standards-mapping.md` (the `VS0-CF-F18-54` row and its
per-invariant mappings), owned by the design-owned gate step **TASK-F18-19** (repository integration
and cleanup, which consolidates the standards-gate reconciliation), and proved by the exact command
`make ff-feature-gate FEATURE=FEATURE-0018` (which invokes the existing `scripts/feature-gate.sh`),
supported by `make vs000-contract-check`; a complete matrix passes and one unmapped applicable
invariant fails the final architecture-approval gate. It is realized as this mapping-matrix gate, not
as a Go conformance test or a DD-13 checker diagnostic.
**`VS0-CF-F18-50` is likewise an architecture-gate case (REQ-F18-01), not a runtime Go case**
(proof classification: non-runtime architecture fixture gate): its exact proof artifacts are the
existing negative fixture `internal/govaccess/fixtures/prohibited_concept.go` and its fixture-gate
assertion `internal/govaccess/fixtures/prohibited_concept_test.go`, owned by the design-owned gate
step **TASK-F18-15** (immutable fixtures + graph validator), proved by the exact command
`make test` and enforced through `make ff-feature-gate FEATURE=FEATURE-0018`. The fixture asserts
that a fixture containing a future-owned/prohibited resource, field, adapter, credential flow, or
provider-native IAM object is rejected by the graph validator before evaluation/publication with
zero external I/O (`externalCallCount == 0`). The DD-13 supplemental static package check
(`make feature-0018-architecture-check`) is **separate** design-owned gate evidence for the DD-01
package topology and does **not** redefine, replace, or stand in for `VS0-CF-F18-50`. Runtime cases
(the `{01..37,
39..49,51,52,53}` set) are the only ones realized as Go conformance tests; `VS0-CF-F18-50` and
`VS0-CF-F18-54` are architecture/standards gates with the exact artifacts, commands, owning tasks
and proof classifications named above. `VS0-CF-F18-38` is a permanent tombstone (no
active case)
and is asserted only by the negative fixture (§8.3). Each runtime test asserts exactly the
registered `expectedState`, `expectedError`, and `expectedSideEffects` for its case, and asserts
`externalCallCount == 0`. **`VS0-CF-F18-12/47/48/49` use a code-agnostic safe-error oracle** — the
test asserts `success == false` and generic outcome class `SafeError`, then asserts safe termination,
discloses nothing to an unauthorized caller, performs no mutation/publication/idempotency
completion, and produces the audit behavior its registered authority requires, and it **must not
assert or choose** a Problem `code`, a `violations[].code` channel, an HTTP status, or any
precedence (§5.2). Successful no-op completion fails these four cases. **`VS0-CF-F18-32` keeps its corrected label**: the guest RoleAssignment
TimeBound violation (`VALIDATION_FAILED`/422 with an RFC 6901 pointer; not a Membership failure).
The requirements’ Exact conformance semantics ledger fields are the
authority; tests neither add nor narrow a case. Positive, negative, race (CAS/terminal),
replay/idempotency, audit-failure, and redaction vectors are all covered per the case inventory.
Inherited `VS0-CF-F01`, `VS0-CF-F02`, `VS0-CF-X01`, `VS0-CF-X02` are reused as supplementary
security/safe-denial evidence and are never conflated with local cases; FEATURE-0015-owned
`VS0-CF-X03` is never FEATURE-0018-local proof.

<!-- END RETAINED SEMANTIC EXCERPT: raw design.md lines 2591-2830 -->

<!-- BEGIN RETAINED SEMANTIC EXCERPT: raw design.md lines 3019-3471 (§7.3 executable conformance gating, §8 classification, §9 traceability, §10.1-§10.3 non-goals/absence/unresolved) -->
### 7.3 Executable conformance gating

Executable conformance is a downstream feature-gate obligation (REQ-F18-22; architecture §10);
design specifies the vectors and wiring only. Standards-gate mapping (REQ-F18-23 /
`VS0-CF-F18-54`) is an architecture-gate obligation realized as a mapping-matrix check, **not**
a runtime Go test (§7.1). Mandatory verification commands, run before the feature is marked
complete and before advancing:

```bash
make fmt
make test
make vet
go test -race ./...
make feature-0018-architecture-check
make vs000-contract-check
make phase2r-drift-check
make ff-feature-gate FEATURE=FEATURE-0018
```

`make ff-feature-gate FEATURE=FEATURE-0018` is the mandatory final gate; the feature is not
complete unless it passes. The sole existing FEATURE-0018 feature-gate integration path is the
existing `scripts/feature-gate.sh`; any wiring of the FEATURE-0018 gate (including the
`feature-0018-architecture-check` target above) is integrated only through that existing script,
with `Makefile` changes permitted only where needed to wire that existing target. No alternative or
replacement feature-gate script is named or created (ADH-2026-074 §3).

---

## 8. Implementation classification ledger

### 8.1 IMPLEMENT (deterministic in-memory contracts, validators, controllers, fakes, conformance)

`internal/govaccess/model/`, `internal/govaccess/validate/`, the semantic-owner and neutral
packages `internal/govaccess/grantintent/`, `internal/govaccess/roleassign/`,
`internal/govaccess/membership/`, `internal/govaccess/roledefinition/`,
`internal/govaccess/accessgroup/`, `internal/govaccess/approval/`,
`internal/govaccess/privileged/`, `internal/govaccess/review/`, `internal/govaccess/exception/`,
`internal/govaccess/governanceprofile/`, `internal/govaccess/approvalreq/`,
`internal/govaccess/policyseam/`, `internal/govaccess/evidence/`,
`internal/govaccess/authzeval/` (files `authz.go`, `authzservice.go`, `eligibility.go`),
`internal/govaccess/idempotency/`, `internal/govaccess/state/`, `internal/govaccess/uow/`, the
leaf child packages `internal/govaccess/clock/`, `internal/govaccess/fixtures/`, and
`internal/govaccess/projection/`, the
root composition/wiring package `internal/govaccess/` (wiring files only), the design-owned
checker `scripts/feature0018-architecture-check/` and Make target
`feature-0018-architecture-check` (DD-13),
and `tests/conformance/feature_0018_*.go`.
These realize REQ-F18-02..22 and REQ-F18-24 observable behavior with zero external I/O.

### 8.2 CONTRACT_ONLY / NO_TASK

- Literal HTTP path ledger (§4.4) — delegated mechanic 2; documented route-shape contract, no
  live route (no public routes authorized).
- Canonical machine-readable schemas for `VS0-SCHEMA-020..024,062..068` — **already supplied by
  an existing manifest-controlled artifact; no new file is created or missing.** The closed,
  machine-readable schema for each FEATURE-0018 kind is carried inline in
  `docs/architecture/vertical-slices/VS-000-contract-registry.yaml`
  (`required`/`optional`/`refs`/`classification`/`redaction`/`mutability`/`retention`/
  `slice0Constraint`), with the finalized data model / canonical contract catalog as the
  field-authority owner. **Verified:** unlike FEATURE-0012/0013 entries 001–005, none of the
  FEATURE-0018 registry entries declares an `externalRef`, so **no `api/schemas/<kind>.json` file
  is registered, exists, or is owed by this feature.** The prior draft's twelve invented
  `api/schemas/*.json` paths are removed as non-existent and unregistered; `validate/` consumes
  the registry entries by reference. This is a CONTRACT reference with no implementation task
  because the artifact already exists and is manifest-controlled — it is not an absent artifact
  marked `NO_TASK`.
- REQ-F18-01 scope/exclusion authority and REQ-F18-23 standards-validation gate — architecture-gate/
  boundary obligations, **not runtime Go tests** (§7.1/§7.3), each with an exact proof artifact,
  command, owning task, and proof classification:
  - `VS0-CF-F18-50` (proof classification: non-runtime architecture fixture gate; REQ-F18-01) —
    exact artifacts `internal/govaccess/fixtures/prohibited_concept.go` and
    `internal/govaccess/fixtures/prohibited_concept_test.go`; owning task **TASK-F18-15**; exact
    commands `make test` and `make ff-feature-gate FEATURE=FEATURE-0018`. The fixture asserts a
    future-owned/prohibited resource, field, adapter, credential flow, or provider-native IAM object
    is rejected by the graph validator before evaluation/publication with `externalCallCount == 0`.
  - `VS0-CF-F18-54` (proof classification: non-runtime standards architecture gate; REQ-F18-23) —
    exact artifact `docs/reviews/architecture-readiness/FEATURE-0018-standards-mapping.md`; owning
    task **TASK-F18-19**; exact command `make ff-feature-gate FEATURE=FEATURE-0018` (invoking the
    existing `scripts/feature-gate.sh`), supported by `make vs000-contract-check`. A complete matrix
    passes; one unmapped applicable invariant fails the final architecture-approval gate.
  The DD-13 `feature-0018-architecture-check` target (`make feature-0018-architecture-check`) is
  separate design-owned gate evidence for the DD-01 package topology and does **not** redefine,
  replace, or stand in for `VS0-CF-F18-50` or `VS0-CF-F18-54`.
- Future adapter boundaries (real IdP/policy/workflow engine, provider-native IAM, transport,
  persistence) — named implementation-neutral boundaries only; no adapter designed or selected.

### 8.3 EXCLUDED (owned by another feature; neither implemented, referenced-as-behavior, nor pulled forward)

F18-RD-01 and requirements §7.2 are the sole excluded-feature and prohibited-concept inventory.
Every listed FEATURE-0019..0026 concern is classified `EXCLUDED`, receives no package, route,
adapter or task, and remains covered by the registered `VS0-CF-F18-50` architecture fixture. This
design neither reproduces that inventory nor turns it into behavior.

---

## 9. Requirement / decision / risk traceability

### Architecture Traceability

- **Decision/change records:** DEC-0060 (Accepted 2026-08-30); ACR-2026-002.
- **Controlling handoff package:** ADH-2026-070 core + Appendix A (semantic contract) +
  Appendix B (registries and evidence); ADH-2026-071 (conformance IDs and REQ/AC mappings only);
  ADH-2026-072 (FEATURE-0013/local-publication boundary, audited-evaluation completion point,
  design-mechanic ownership and reviewer classification only); ADH-2026-073
  (terminal-ExceptionProposal AuditEvent taxonomy registration — the exact Grant/Deny
  descriptor matrix — only); ADH-2026-074 (FEATURE-0018-local
  abort-safe provisional publication sequence, neutral completed-record publication-instant flow,
  `state.Store` as sole owner of caller-reservation lifecycle, and `scripts/feature-gate.sh` as the
  sole existing feature-gate integration path only).
- **Controlling design-normalization handoff:** ADH-2026-075 (FEATURE-0018 executable-design
  normalization and deterministic contract check). ADH-2026-075 makes the hash-bound
  `design-mechanics-contract.json` the sole authority for this design's non-product-semantic
  executable mechanics (package edges, Go 1.22 API signatures/result variants/factories/permitted
  callers, `StateEditor`/`ContextBoundChange` surface, transaction-mode state machines, identity/
  publication/sequence/overflow/failure rules, static-checker scope/rules/protected symbols/fixtures
  and non-runtime conformance-gate mappings); requires the numbered sections to cite the contract
  rather than restate those mechanics; adds the contract to the controlled design-review projection
  with verifiable hash binding; and authorizes the deterministic pre-review checker
  `scripts/feature0018-design-contract-check.py` (`make feature-0018-design-contract-check`).
  ADH-2026-075 changes no product semantics, REQ/AC, route, writer, event, error, dependency
  boundary or Phase 2R exclusion.
- **Reused decision authorities (by reference, not redefined):** DEC-0026, ADH-2026-012,
  DEC-0037 (FEATURE-0012 grammar/seven scopes); DEC-0043, ADH-2026-017, DEC-0044 (FEATURE-0013
  envelopes/immutable published definitions); ADH-2026-058..066 (FEATURE-0016 ExecutionTarget
  actions); DEC-0028, DEC-0036, ADH-2026-067..069 (FEATURE-0017 seam); DEC-0050, ADH-2026-033
  (GovernanceProfile envelope); DEC-0018/0019/0020/0035/0048 (governance/simplicity/
  extensibility); DEC-0022/0036, RFC-0012/0021/0022/0023 (non-authoritative future-adapter
  deferral). ADH-2026-049 informs the FEATURE-0018 `INTERNAL_ERROR`-on-audit-failure pattern
  via F18-RD-20.
- **Owned schema IDs:** VS0-SCHEMA-020, 021, 022, 023, 024, 062, 063, 064, 065, 066, 067, 068.
- **Writer IDs:** owned VS0-WRITER-008, VS0-WRITER-022; consumed unchanged VS0-WRITER-007
  (`GovernanceProfile.spec`, `RoleDefinition.spec`), VS0-WRITER-011 (`DecisionRecord.record`,
  `AuditEvent.record`); out-of-scope VS0-WRITER-023 (FEATURE-0020).
- **Decision profiles:** `authorization-decision/v1`, `approval-decision/v1`,
  `exception-decision/v1` (FEATURE-0013 registered extensions).
- **State authorities:** no FEATURE-0018-owned `VS0-STATE-*` registration; every resource
  lifecycle and effective projection is fixed by its F18-RD (REQ-F18-04/05/07/08/12/13/14/16/18)
  and realized by the corresponding §3.2 component. No conformance case substitutes for a state
  authority.
- **Error codes:** inherited FEATURE-0012 top-level `AUTH_REQUIRED`/401, `RESOURCE_NOT_FOUND`/404,
  `AUTHORIZATION_DENIED`/403, `CONFLICT`/409, `VALIDATION_FAILED`/422, `INTERNAL_ERROR`/500;
  inherited failure mappings VS0-F01→VS0-CF-F01, VS0-F02→VS0-CF-F02; violation codes (in
  `violations[].code` only) `REVIEWER_CONFLICT`, `REVIEW_BENEFICIARY_UNRESOLVED`. No new
  top-level code, violation code, HTTP status, or precedence is introduced.

### Canonical coverage ledger

Every approved REQ and every approved AC appears exactly once with its design disposition. No
ID is created, renumbered, merged, split, reinterpreted, or omitted. Disposition legend:
**IMPLEMENT** = realized by named §3.2 components + §7 conformance; **CONTRACT_ONLY** = §8.2
contract/gate artifact with no source task; **EXCLUDED** = §8.3 tombstone/negative case.

#### Requirements

| REQ | Title | Disposition | Primary component(s) / artifact |
|---|---|---|---|
| REQ-F18-01 | Current-feature-only authority | CONTRACT_ONLY (+ negative fixture) | §8.2 boundary; `fixtures/` `VS0-CF-F18-50` |
| REQ-F18-02 | Closed contract inventory and profiles | IMPLEMENT + CONTRACT_ONLY (§4.4 paths) | `model/`, `validate/`, §4.1 profiles, §4.4 path ledger |
| REQ-F18-03 | Stable principal identity | IMPLEMENT | `model/` `PrincipalRef`/`AssuranceEvidence`; `authz.go` |
| REQ-F18-04 | Direct principal and scoped AccessGroup assignments | IMPLEMENT | `accessgroup/`, `roleassign/`, `validate/` |
| REQ-F18-05 | Membership is non-authorizing | IMPLEMENT | `membership/` (VS0-WRITER-022), `validate/` |
| REQ-F18-06 | Distributed action ownership and central role composition | IMPLEMENT | `model/` `ActionTargetBinding`; `authz.go`; FEATURE-0016 action reuse |
| REQ-F18-07 | Versioned RoleDefinition | IMPLEMENT | `roledefinition/` (VS0-WRITER-007 consumed) |
| REQ-F18-08 | Scoped RoleAssignment | IMPLEMENT | `roleassign/` (VS0-WRITER-008) |
| REQ-F18-09 | Deterministic scoped authorization composition | IMPLEMENT | `authz.go` + `authzservice.go` (audited, DD-12), `membership/`, `eligibility.go` |
| REQ-F18-10 | FEATURE-0017 adoption without reinterpretation | IMPLEMENT | `policyseam/`, `approval/` |
| REQ-F18-11 | Bounded FEATURE-0017 subject/target use | IMPLEMENT | `policyseam/`, `privileged/` |
| REQ-F18-12 | Bounded ApprovalPolicy | IMPLEMENT | `approval/`, `eligibility.go`, `validate/` |
| REQ-F18-13 | Immutable terminal approval evidence | IMPLEMENT | `approval/`, `roleassign/`, `idempotency/` |
| REQ-F18-14 | JIT privileged access authorizes one temporary grant | IMPLEMENT | `privileged/` (→ `roleassign.ActivationPort`), `roleassign/`, `clock/` |
| REQ-F18-15 | Constrained break-glass access | IMPLEMENT | `privileged/`, `review/`, `evidence/` |
| REQ-F18-16 | Snapshot-based AccessReview | IMPLEMENT | `review/` (→ `roleassign.ReviewRemediationPort`), `eligibility.go`, `roleassign/`, `membership/` |
| REQ-F18-17 | Bounded immutable exception evidence | IMPLEMENT | `exception/`, `approvalreq/`, `approval/` |
| REQ-F18-18 | FEATURE-0018-limited GovernanceProfile v1 | IMPLEMENT | `governanceprofile/` (VS0-WRITER-007 consumed) |
| REQ-F18-19 | FEATURE-0013 adoption | IMPLEMENT | `evidence/`, `authzservice.go`, profile registration (DD-09) |
| REQ-F18-20 | Audit before authorization-changing publication | IMPLEMENT | `evidence/` (§5.4), `authzservice.go` (§5.5/DD-12); all publishing owner packages |
| REQ-F18-21 | Deterministic validation and safe denial | IMPLEMENT | `validate/` precedence (§5.1), `idempotency/`, `projection/` |
| REQ-F18-22 | Deterministic in-memory foundation and local conformance | IMPLEMENT | `clock/`, `fixtures/`, `tests/conformance/feature_0018_*.go` |
| REQ-F18-23 | Standards-validation gate | CONTRACT_ONLY | §8.2 architecture-gate matrix; `VS0-CF-F18-54` |
| REQ-F18-24 | Progressive and normally hidden user friction | IMPLEMENT | `projection/` |

#### Acceptance criteria

| AC | Disposition | Owning conformance case / note |
|---|---|---|
| AC-F18-01 | IMPLEMENT | VS0-CF-F18-01 |
| AC-F18-02 | IMPLEMENT | VS0-CF-F18-02 |
| AC-F18-03 | IMPLEMENT | VS0-CF-F18-03 |
| AC-F18-04 | IMPLEMENT | VS0-CF-F18-04 |
| AC-F18-05 | IMPLEMENT | VS0-CF-F18-05 |
| AC-F18-06 | IMPLEMENT | VS0-CF-F18-06 |
| AC-F18-07 | IMPLEMENT | VS0-CF-F18-07 |
| AC-F18-08 | IMPLEMENT | VS0-CF-F18-08 |
| AC-F18-09 | IMPLEMENT | VS0-CF-F18-09 |
| AC-F18-10 | IMPLEMENT | VS0-CF-F18-10 |
| AC-F18-11 | IMPLEMENT | VS0-CF-F18-11 |
| AC-F18-12 | IMPLEMENT | VS0-CF-F18-12 (code-agnostic non-success safe-error oracle — §5.2/§7.1) |
| AC-F18-13 | IMPLEMENT | VS0-CF-F18-13 |
| AC-F18-14 | IMPLEMENT | VS0-CF-F18-14 |
| AC-F18-15 | IMPLEMENT | VS0-CF-F18-15 |
| AC-F18-16 | IMPLEMENT | VS0-CF-F18-16 |
| AC-F18-17 | IMPLEMENT | VS0-CF-F18-17 |
| AC-F18-18 | IMPLEMENT | VS0-CF-F18-18 |
| AC-F18-19 | IMPLEMENT | VS0-CF-F18-19 |
| AC-F18-20 | IMPLEMENT | VS0-CF-F18-20 |
| AC-F18-21 | IMPLEMENT | VS0-CF-F18-21 |
| AC-F18-22 | IMPLEMENT | VS0-CF-F18-22 |
| AC-F18-23 | IMPLEMENT | VS0-CF-F18-23 |
| AC-F18-24 | IMPLEMENT | VS0-CF-F18-24 |
| AC-F18-25 | IMPLEMENT | VS0-CF-F18-25 |
| AC-F18-26 | IMPLEMENT | VS0-CF-F18-26 |
| AC-F18-27 | IMPLEMENT | VS0-CF-F18-27 |
| AC-F18-28 | IMPLEMENT | VS0-CF-F18-28 |
| AC-F18-29 | IMPLEMENT | VS0-CF-F18-29 |
| AC-F18-30 | IMPLEMENT | VS0-CF-F18-30 |
| AC-F18-31 | IMPLEMENT | VS0-CF-F18-31 |
| AC-F18-32 | IMPLEMENT | VS0-CF-F18-32 |
| AC-F18-33 | IMPLEMENT | VS0-CF-F18-33 |
| AC-F18-34 | IMPLEMENT | VS0-CF-F18-34 |
| AC-F18-35 | IMPLEMENT | VS0-CF-F18-35 |
| AC-F18-36 | IMPLEMENT | VS0-CF-F18-36 |
| AC-F18-37 | IMPLEMENT | VS0-CF-F18-37 |
| AC-F18-38 | EXCLUDED | Retired duplicate of F18-SCN-29; permanent tombstone `VS0-CF-F18-38`; enforced as negative case (§8.3, requirements §7.3); never an active obligation |
| AC-F18-39 | IMPLEMENT | VS0-CF-F18-39 |
| AC-F18-40 | IMPLEMENT | VS0-CF-F18-40 |
| AC-F18-41 | IMPLEMENT | VS0-CF-F18-41 |
| AC-F18-42 | IMPLEMENT | VS0-CF-F18-42 |
| AC-F18-43 | IMPLEMENT | VS0-CF-F18-43 |
| AC-F18-44 | IMPLEMENT | VS0-CF-F18-44 |
| AC-F18-45 | IMPLEMENT | VS0-CF-F18-45 |
| AC-F18-46 | IMPLEMENT | VS0-CF-F18-46 |
| AC-F18-47 | IMPLEMENT | VS0-CF-F18-47 (code-agnostic non-success safe-error oracle — §5.2/§7.1) |
| AC-F18-48 | IMPLEMENT | VS0-CF-F18-48 (code-agnostic non-success safe-error oracle — §5.2/§7.1) |
| AC-F18-49 | IMPLEMENT | VS0-CF-F18-49 (code-agnostic non-success safe-error oracle — §5.2/§7.1) |

### Conformance Mapping

Every active AC maps one-to-one to the local conformance ID with the same numeric suffix, per
ADH-2026-071 and the requirements Acceptance-to-conformance mapping: each active acceptance
criterion maps to the FEATURE-0018 conformance case bearing the same numeric suffix, for the
suffix values in {01..37, 39..49}, as recorded row-by-row in the Acceptance criteria table above.
`AC-F18-38` is excluded (tombstone `VS0-CF-F18-38`; no active row).
REQ-only proof cases `VS0-CF-F18-50..54` back REQ-F18-01/18/19/(22,20)/22/23 as recorded in the
requirements Requirement-to-conformance mapping. Design adds, renames, and reinterprets no case;
each case retains its exact registry owner (FEATURE-0018) and semantics. Inherited
`VS0-CF-F01/F02/X01/X02` remain supplementary; FEATURE-0015-owned `VS0-CF-X03` is never
FEATURE-0018-local proof.

#### Canonical ownership detail (citation only)

Per-requirement schema, writer, lifecycle/state, error and conformance ownership remains
authoritative in ADH-2026-071 and requirements §§8.3, 8.6, 10.3 and 10.4. This design does not
copy or reinterpret that semantic matrix. The compact REQ disposition table and exact per-AC
mapping above are the complete design-stage coverage ledgers.

### 9.4 Risk traceability

Independent-review findings closed by the approved package (F18-SEC-001..005; F18-REN-001/002;
F18-REN2-001/002; F18-REN3-001/002; F18-REN4-001; renewal-05 clean) are preserved by design as
follows: sole-writer boundaries (§3.3) and no second RoleAssignment/decision/audit writer;
exact-Human-only eligibility with structural rejection of relationship-derived eligibility
(`eligibility.go`, VS0-CF-F18-47/49); the membership-enabled grant ceiling with no
trusted-provisioner bypass (`membership/`, VS0-CF-F18-46); the indirect-beneficiary conflict
set (`eligibility.go`, VS0-CF-F18-45). Applicable FEATURE-0012/0013 inherited risks are
mitigated by reuse-not-fork (DD-03), single-writer per field (F12-R07), bounded inputs and
redaction (F12-R10/R13), and immutable published definitions (F12-R12 / DEC-0044). Residual
risk remains Medium until independent security review and executable proof; expected
Low-to-Medium after feature-gate evidence (requirements §8.4). No new implementation risk is
introduced beyond deterministic in-memory scope.

### 9.5 Design-delta → component → conformance traceability

Compact map of each design-owned mechanic (the only deltas this stage introduces over the
approved requirements) to its realizing component(s) and its exact conformance evidence. No new
behavior is added; each row is a representation/wiring choice grounded in the cited authority.

| Design delta (DD/§) | Requirement authority | Realizing component(s) | Conformance evidence |
|---|---|---|---|
| Distinct `PrincipalRef`/`AccessGroupRef`/`RoleHolderRef`/`EligibilityRef` value types (DD-03, §3.2) | REQ-F18-03/04; SEC-06 | `model/`, `eligibility.go`, `validate/` | VS0-CF-F18-01/02/05; VS0-CF-F18-47/49 (code-agnostic non-success safe-error oracle) |
| Acyclic topology + layered enforcement: sole FEATURE-0017 caller, distinct RD-13/14/16 ports, acyclic RD-08 direct revoke, neutral result/finalization values, non-escaping editor and sole transaction-lifecycle owner (DD-01/DD-13, §3.4) | F18-RD-08/10/11/13/14/16; SEC-05 | owner packages, `operation/`, `command/`, `evidence/`, `state/`, `uow/`, `scripts/feature0018-architecture-check/` | named static checker plus focused runtime writer/port/lifecycle/finalization tests; `VS0-CF-F18-50` is not topology evidence |
| Two-lock, F18-local non-durable publication with caller/controller/evidence transactions and authorization-bearing/automatic-mutation-only authority modes; stage-10 begin checks complete candidate dependencies and captures one context; transaction-only `ApplyFinalized` creates exact context-bound receipts/proofs before completed-evaluation and carrier preparation (DD-02/DD-10, §5.4) | ADH-2026-072 §§2–3; REQ-F18-20/21/22 | owner packages, `authzeval/`, `state/`, `uow/`, `evidence/`, `idempotency/`, `clock/` | VS0-CF-F18-17 (`CONFLICT` precedence/replay); VS0-CF-F18-53 (deterministic zero-I/O); focused ordering tests proving accepted evidence precedes all-or-nothing publication (no registered case directly proves this design mechanic — §7.2); focused currentness/receipt/Seal tests |
| Single permit-gated `FinalizeAt(state.MutationFinalizationPermit)` with unforgeable `PermitBinding` minted per issuance, transaction-private permit-issuance registry, and runtime `ApplyFinalized` permit-provenance check consuming one entry per permit (DD-01/DD-13, §5.4) | REQ-F18-20/21; ADH-2026-072 §3 | `state/`, owner packages, `uow/` | focused permit-provenance/finalization tests; `make feature-0018-architecture-check` (rules 009/010) |
| Checker-enforced closed `state.ContextBoundChange`/`ApplyTo` boundary (exported interface Go 1.22 cannot seal): `state` alone declares the interface, declares/constructs `StateEditor`, and invokes `ApplyTo`; only the nine registered owner packages may implement/embed/adapt/wrap/expose `ContextBoundChange` or `ApplyTo`, each naming `StateEditor` only as its own `ApplyTo` parameter with no editor escape; only each matching owner `FinalizeAt` may extract `MutationFinalizationPermit.Binding`; the DD-13 checker is the required enforcement boundary for source-level copied-binding wrappers, and runtime permit/context/receipt validation guarantees rejection only of absent/mismatched/foreign-transaction/reused/context-mismatched bindings (DD-01 rules 1–5, DD-13 `F18-ARCH-012`, §3.4/§5.4) | SEC-05; F18-RD-10/11/13/14/16; REQ-F18-20/21 | owner packages, `state/`, `uow/`, `scripts/feature0018-architecture-check/` | `make feature-0018-architecture-check` (rule 012, `state`+nine-owner positive + foreign-implementation/wrapper/delegated-`ApplyTo`/`uow`-`ApplyTo`/editor-escape/copied-binding-wrapper/foreign-permit-binding negatives); focused context-bound-change boundary runtime tests |
| Publishable terminal `DomainConclusion` evidence path for a resource-less terminal decision — ExceptionProposal `Deny` — that creates no ExceptionGrant yet atomically publishes the mandatory `exception-decision/v1` DecisionRecord and required AuditEvents after stage 10 and before result release (DD-01, §5.4/§5.5) | F18-RD-17/19/20 | `exception/`, `evidence/`, `state/`, `uow/` | VS0-CF-F18-22 (non-exceptionable terminal Deny), VS0-CF-F18-43 (approved terminal Deny); focused domain-conclusion publication tests |
| ADH-2026-073 terminal-ExceptionProposal evidence matrix: canonical ordered Grant descriptor set (`exceptionproposal.decided`, then `exceptiongrant.issued`) and one-event Deny conclusion, each linked to the exact immutable proposal, terminal reason, mandatory `exception-decision/v1` record and containing operation/correlation evidence through FEATURE-0013 | REQ-F18-19/20; ADH-2026-073 | `exception/`, `evidence/`, `state/`, `uow/` | focused Grant/Deny descriptor-set, linkage and Seal rejection tests |
| Abort-safe provisional publication sequence (`state` derives `current + 1` under `stateMu`, persists only at successful aggregate commit after next-sequence revalidation; abort, failed CAS, cancellation, evidence failure and panic-cleanup consume no sequence); immutable neutral `operation.PublicationInstant` publication-instant flow (`state` alone invokes, under `stateMu` from one clock read, the exported `operation.NewPublicationInstant(time.Time)` and `evidence.NewPublicationContext(publicationSequence, operation.PublicationInstant, evidence.TransactionIdentity)` factories (with `evidence.TransactionIdentity` minted via `evidence.NewTransactionIdentity`) — DD-13 confines all three call sites to `state` — → extracted by `evidence` for `operation.NewCompletionBinding` → read by `idempotency.PrepareCompletedResult` for the fixed 24-hour deadline with no second clock read and no `state`/`evidence`/`clock` import; runtime validation and `Seal` reject zero/foreign/forged/context-mismatched values); a sealed transaction-unique `evidence.PublicationContextBinding` derived by `evidence.NewPublicationContext` from an `evidence.TransactionIdentity` (evidence-owned `evidence.StoreNonce` + strictly-monotonic never-reused `transactionSerial`) independent of sequence and instant, so `ApplyFinalized`/`Seal` reject an aborted or foreign transaction's binding even when sequence and fixed instant repeat; and `state.Store` as sole owner of `InspectOrReserveCallerMutation`, `CallerReservationLease`, the in-flight table, reservation cleanup and the single authoritative waiter-notification protocol (record outcome + remove reservation under lock, release both locks, broadcast only after unlock, waiters gate on the recorded outcome for lost-wakeup safety), with `idempotency` retaining only neutral value types and pure completed-record helpers (DD-10/DD-11/DD-14, §5.4) | ADH-2026-074 §§1–3; REQ-F18-20/21 | `state/`, `operation/`, `evidence/`, `idempotency/`, `uow/` | VS0-CF-F18-17 (replay/conflict precedence); VS0-CF-F18-53 (deterministic zero-I/O); focused deterministic-time/sequence/completed-record/publication-instant, transaction-identity/context-binding and two-lock waiter-notification tests |
| Caller idempotency lookup key separated from exact target/scope/digest binding; independently observable reservation table and direct waiter handle; controller-owned operations use only exact transition/version CAS; replay rechecks after wake via the authoritative waiter-notification protocol (recorded terminal outcome + re-inspection) (DD-11, §5.4) | REQ-F18-13/21 | `idempotency/`, `state/`, `uow/` | VS0-CF-F18-17 plus payload/target/scope conflict and controller effective-once vectors |
| Stages 7–9 produce a non-authoritative candidate plus complete dependency-version claim; state proves currentness, `authzeval` owns publication-time revalidation, and only applicable success completes an audited Allow/Deny evaluation; automatic mutation-only effects never fabricate one (DD-12, §5.1/§5.5) | ADH-2026-072 §2; REQ-F18-19/20/21 | `authzservice.go`, `authzeval/`, `evidence/`, `uow/`, `state/` | focused audit-before-result ordering tests (no registered case directly proves this design mechanic — §7.2); VS0-CF-F18-17 (conflict precedence and replay recheck); focused automatic-effect tests |
| Named standard-library topology checker covers only imports/declarations/typed calls/local AST patterns; sealed runtime validators and tests separately prove equality, cardinality and ordering (DD-13) | SEC-05; F18-RD-10/11/13/14/16 | `scripts/feature0018-architecture-check/`, transaction tests | `make feature-0018-architecture-check` plus focused runtime Seal suites; supplemental gate evidence only |
| Pure vs stateful stage split (§5.1) | REQ-F18-21 | `validate/`, `authz.go`, `eligibility.go` (pure) vs `authzservice.go`/`uow/` (stateful) | VS0-CF-F18-21 (validation precedence + safe denial) |
| Per-owner typed `ApprovalRequirementDeriver` ports (RoleGrant deriver owned by `grantintent`, not `roleassign`); approval consumes, never derives (§3.4) | F18-RD-10 | `approvalreq/`, `grantintent/`, `approval/`, owners | VS0-CF-F18-13/40/41 |
| Bounded FEATURE-0017 seam invoked only by `privileged/` (§3.4, DD-08) | F18-RD-11 | `policyseam/`, `privileged/` | VS0-CF-F18-13/17 (`Indeterminate→Deny`) |
| Distinct RD-13/RD-14/RD-16 RoleAssignment workflow triggers via `roleassign.GrantPort`/`ActivationPort`/`ReviewRemediationPort`; separate acyclic `command.RoleAssignmentRevokeService -> roleassign.DirectRevokePort + uow.MutationCoordinatorPort` RD-08 caller ingress (not a fourth trigger, no route, no second writer); Membership publication owned in-owner (§3.4, DD-01) | F18-RD-08/13/14/16 | `command/`, `grantintent/`, `privileged/`, `review/`, `membership/`, `roleassign/`, `uow/` | VS0-CF-F18-13/17/19/20 (per RD trigger); §7.2 distinct-port/direct-revoke tests |
| Literal route ledger rebuilt one-to-one from Appendix B caller ops; proposals on owning collections (RoleAssignment collection POST = `roleassignment.grant`; `exceptiongrant.propose` collection action); no `/approval-requests/actions/submit-*` route; `CreateParent.parentScopeKind` = authorization target scope only (§4.4, DD-07) | F18-RD-02/06 | CONTRACT_ONLY ledger; no source | no runtime case (route-shape contract) |
| Code-agnostic safe-error oracle for `VS0-CF-F18-12/47/48/49` (`success=false`, `SafeError`, no disclosure/effect; no code/channel/status/precedence selected) (§5.2/§7.1) | REQ-F18-16/21; ADH-2026-071 | `tests/conformance/feature_0018_*.go` | VS0-CF-F18-12/47/48/49 |
| Exact FEATURE-0013 adoption bundle and root composition registration | REQ-F18-19; F18-RD-19 | `evidence/profiles.go`, root wiring, existing `decision/bundle.Load` | exact-three, duplicate, unknown-version, malformed and no-fourth-profile tests |
| Canonical schemas consumed from the manifest-controlled VS-000 registry inline entries; no `api/schemas/*.json` created (§8.2) | REQ-F18-02 | CONTRACT reference to VS-000 registry; `validate/` consumes | no runtime case (registry-owned content) |
| Deterministic `Clock` leaf package + one state-captured authoritative publication instant + immutable fixtures; no Phase-2R production wall clock (DD-05/DD-06, §5.4) | REQ-F18-14/20/22 | `clock/`, `state/`, `uow/`, `fixtures/` | VS0-CF-F18-17/20 plus VS0-CF-F18-53 (`externalCallCount==0`) |
| Standards mapping matrix as architecture gate, not runtime test (§7.1/§7.3) | REQ-F18-23 | mapping-matrix check | VS0-CF-F18-54 (non-runtime standards gate); exact artifact `docs/reviews/architecture-readiness/FEATURE-0018-standards-mapping.md`, owning task TASK-F18-19, command `make ff-feature-gate FEATURE=FEATURE-0018` |
| Exact Go-1.22 §5.4 API: named closed result types (`state.ReserveOutcome`/`CallerBeginOutcome`/`ControllerBeginOutcome`/`EvidenceBeginOutcome`/`CommitOutcome`, `evidence.CurrentAllowOutcome`/`TerminalDescriptorValidation`, `operation.FinalizationOutcome[T,D]`) replacing every anonymous multi-arm `->`; the sealed `operation.MechanicalFailure` type replacing every open `error` return; complete Begin/Store-construction/Commit failure propagation; exact Commit-after-Seal rejection via the `CommitOutcome` `MechanicalFailure` variant (§5.4 notation legend + API) | REQ-F18-20/21; ADH-2026-072 §3; ADH-2026-074 §§1–2 | `operation/`, `state/`, `evidence/`, `authzeval/`, `uow/`, `scripts/feature0018-architecture-check/` | focused §7.2 result-type/mechanical-failure, Begin-failure, Commit-after-Seal and coordinator-API tests; `make feature-0018-architecture-check` (symbol ledger + rule-010 call-site confinement); no registered runtime case directly proves this design mechanic |
| Deterministic process-local unique identity issuance and fail-closed `uint64` exhaustion for `evidence.StoreNonce` (evidence-package-global mutexed counter, `state.NewStore` failure propagation, no partial Store), `transactionSerial`, per-root `publicationSequence` and `state.FinalizationPermitToken` (each strictly monotonic, never reused, fail-closed at `math.MaxUint64`, no wraparound); concurrent Store construction (DD-14, §5.4) | REQ-F18-20/21; ADH-2026-074 §§1–2 | `evidence/`, `state/`, `operation/`, `uow/`, `scripts/feature0018-architecture-check/` | focused §7.2 identity-issuance/exhaustion tests (distinct/concurrent nonces, exhaustion fail-closed, monotonic never-reused serial across committed/aborted begins, sequence-overflow and zero-sequence rejection) plus the aborted-binding-replay negative; `make feature-0018-architecture-check` (protected-symbol confinement of `evidence.NewStoreNonce`/`NewTransactionIdentity`/`NewPublicationContext` and `state.NewStore`) |

No delta introduces a semantic, ownership transfer, error code, route, writer, state, or
conformance case beyond the approved requirements.

---

## 10. Non-goals, absence ledger, and unresolved report

### 10.1 Non-goals (design-level)

F18-RD-01, REQ-F18-01 and requirements §7 are the sole product and future-feature exclusion
inventory; design does not repeat or extend it. Design-specific absences are recorded once in
§10.2. The §4.4 paths remain contract shapes only, and §8.2 owns artifact classification.

### 10.2 Absence ledger (deliberately not present, with reason)

| Absent element | Reason | Authority |
|---|---|---|
| `internal/server` route registration / `net/http` handlers | No public routes authorized; deterministic in-process foundation | requirements §1; architecture §3 |
| Durable persistence / database / queue | Phase 2R in-memory only; side-effect-free | REQ-F18-22; OPS-02 |
| Real IdP/OIDC/SCIM, policy engine, workflow engine, provider adapters | Deferred, non-authoritative future hypotheses; fresh FEATURE-0011 assessment required | architecture §8 (Wrap/Deferred rows) |
| New top-level Problem / violation code, HTTP status, precedence | Only inherited codes permitted | ADH-2026-071; error anti-drift |
| `VS0-STATE-*` FEATURE-0018 registration | Lifecycles are F18-RD-owned; no state-machine ID assigned | requirements §8.6 |
| VS0-WRITER-023 (`ProfileAssignment.spec`) | FEATURE-0020-owned; out of scope | requirements §8.3 |
| `api/schemas/<kind>.json` files for FEATURE-0018 kinds | Not registered and not owed: the FEATURE-0018 VS0-SCHEMA entries carry no `externalRef`; the machine-readable schema is supplied inline in the manifest-controlled VS-000 contract registry (with the finalized data model/catalog as field authority) | §3.1, §8.2; VS-000 registry (no `externalRef` on VS0-SCHEMA-020..024,062..068) |
| Fourth DecisionProfile / competing envelope | Exactly three profiles; AccessReview stays an LRO | REQ-F18-19 |
| Excluded-feature mechanisms (FEATURE-0019..0026) | Out of scope; enforced by leakage sentinels + VS0-CF-F18-50 | requirements §7.2; architecture §9.2 |

### 10.3 Unresolved report

No unresolved product-semantic or ownership decision remains. Requirements §9 plus
ADH-2026-072 §3 authorize the bounded mechanics resolved by DD-05..07 and DD-10..13: deterministic
fixtures/time, literal route shapes, local state/locking/idempotency, owner-controlled
prepare/finalize, transaction-created context-bound receipts, evidence-owned proofs/carriers,
commit/abort/shutdown, and the named checker. Within those same §3 mechanics, this revision also
resolves the single permit-gated `FinalizeAt(state.MutationFinalizationPermit)` contract with an
enforceable transaction-owned permit-issuance registry and unforgeable `PermitBinding`, and the
publishable terminal `DomainConclusion` evidence path (ExceptionProposal `Deny`) that publishes the
mandatory F18-RD-17/19/20 decision and audit evidence while creating no ExceptionGrant. This
revision further makes exact, within those same §3 mechanics, the ADH-2026-074 §§1–3 corrections:
the abort-safe provisional `current + 1` publication sequence persisted only at successful aggregate
commit (with abort/CAS-loss/cancellation/evidence-failure/panic-cleanup consuming no sequence); the
immutable neutral `operation.PublicationInstant` (DD-14) built by `state` alone, under `stateMu`
from one deterministic clock read, by invoking the exported `operation.NewPublicationInstant(time.Time)`
factory (which validates a non-zero canonical UTC input) and the exported
`evidence.NewPublicationContext(publicationSequence, operation.PublicationInstant, evidence.TransactionIdentity)`
factory (which validates the
instant, the non-zero provisional sequence and the transaction identity and derives the sealed
`evidence.PublicationContextBinding`) — with the DD-13 checker, not the compiler, confining the
`operation.NewPublicationInstant`, `evidence.NewTransactionIdentity` and `evidence.NewPublicationContext`
call sites
to `state`, and runtime validation plus `Seal` rejecting zero/foreign/forged/context-mismatched
values — extracted by
`evidence` for `operation.NewCompletionBinding`, and read by `idempotency.PrepareCompletedResult` to
derive the fixed 24-hour completed-record deadline without a second clock read or a
`state`/`evidence`/`clock` import; and `state.Store` as the sole owner of the private in-flight
entry/table, `InspectOrReserveCallerMutation`, `CallerReservationLease`, reservation cleanup and
waiter signaling, with `idempotency` retaining only the neutral `IdempotencyLookupKey`,
`RequestBinding`, `PreparedCompletedResult`, `CompletedRecord` values and pure completed-record
lifecycle helpers. Within those same §3 mechanics, this revision (design-reconciliation-10)
additionally closes the `state.ContextBoundChange` cross-package enforcement gap without any
compiler-sealing claim Go 1.22 cannot back: `state.ContextBoundChange` and its `ApplyTo` are an
exported interface treated as a checker-enforced closed boundary (DD-01 rules 1–5, DD-13
`F18-ARCH-012`), whereby only the nine registered owner packages may
implement/embed/adapt/wrap/expose the interface or `ApplyTo`, only the `state` transaction
invokes `ApplyTo`, only each matching owner `FinalizeAt` extracts
`MutationFinalizationPermit.Binding`, and runtime permit/claim/context/receipt/consumed-token
validation rejects a forged, foreign, reused, or context-mismatched change carrying no matching
binding;
the compiler enforcement claim over that interface is removed and replaced by the checker rule plus
runtime backstop. Within those same §3 mechanics, this revision (design-reconciliation-11) resolves
the two remaining `ContextBoundChange` enforcement inconsistencies without any semantic change.
First, it makes `F18-ARCH-012` internally satisfiable: `state` alone may declare the
`state.ContextBoundChange` interface, declare and construct `StateEditor`, and invoke `ApplyTo` from
its transaction; each of the nine registered owner finalized-change types may name `StateEditor` only
as the synchronous parameter of its own registered `ApplyTo` implementation and use only its
registered direct editor methods there; owners may not acquire, retain, return, capture, embed, wrap,
adapt, or otherwise expose `StateEditor`; all other packages remain forbidden from declaring,
implementing, embedding, adapting, wrapping, exposing, invoking, or delegating
`ContextBoundChange`/`ApplyTo`; and only the matching owner `FinalizeAt` may extract
`MutationFinalizationPermit.Binding`. Second, it corrects the runtime guarantee: the DD-13 static
checker is the required enforcement boundary for source-level copied-binding wrappers, while runtime
permit/context/receipt validation guarantees only rejection of absent, mismatched,
foreign-transaction, reused, or context-mismatched bindings and no longer claims independent
detection of wrappers that copy an otherwise-valid binding. This is a `DESIGN_EXECUTABILITY`
refinement of the DD-13 static/runtime split, not
a semantic, ownership, route, writer, error, state or conformance change. Within those same §3
mechanics, this revision (design-reconciliation-12) resolves four remaining design-executability
inconsistencies without any semantic change: (1) it makes `F18-ARCH-012` the sole canonical
`ContextBoundChange`/`StateEditor` restriction and reconciles `F18-ARCH-004`, `F18-ARCH-005` and
`F18-ARCH-010` so
that editor construction, acquisition, parameter use, retention, return, capture, extraction, escape,
forwarding, embedding, wrapping, adaptation, and all use — and every editor negative fixture — are
governed exclusively by `F18-ARCH-012`, which permits each
registered owner to receive and synchronously use `state.StateEditor` only as the parameter of its
own registered direct `ApplyTo` implementation, using only its registered direct editor methods, and
rejects every other editor use; `F18-ARCH-004`, `F18-ARCH-005` and `F18-ARCH-010` assert no
editor-scope condition, own no editor fixture, and cite `F18-ARCH-012` rather than restate divergent
editor wording (the stale Rule 005 and Rule 010 editor-fixture assignments are removed); (2) it requires each registered owner to directly implement only its
own single registered concrete `ContextBoundChange` type and rejects embedding, adapting, wrapping,
aliasing, forwarding, delegation, or copied-binding re-exposure by every package including the
registered owner package itself, adding an isolated owner-wrapper negative fixture; (3) it replaces
the invalid Go pseudo-union release API with the exact
`authzeval.ReleaseAfterPublication(evidence.CompletedAuthorizationEvaluation, state.CommitReceipt)
-> AuthorizationResult | operation.MechanicalFailure`, which constructs/releases an
`AuthorizationResult` only and takes no `operation.ResultMaterial`, keeping
`idempotency.PreparedCompletedResult` internal to `idempotency` and adding no `authzeval -> idempotency`
import; and (4) it makes §3.4 the one exhaustive allowed package-edge ledger, adds the `root -> authzeval`
composition/wiring edge, and has `F18-ARCH-005_PACKAGE_DIRECTION` reference only that ledger. Within
those same §3 mechanics, this revision (design-reconciliation-13) closes five remaining
design-executability/stale-transcription inconsistencies without any semantic change: (1) it
reconciles §3.2 and §3.4 one-for-one by adding the actually required `root -> state`,
`root -> evidence` (composition/wiring: clock injection into `state.Store` and FEATURE-0013 bundle
bytes from `evidence/profiles.go`) and `validate -> model` edges to the exhaustive §3.4 ledger and to
the matching §3.2 "Reuses" columns, introducing no cycle or new owner; (2) it replaces the two
remaining Go-1.22-invalid pseudo-union input APIs with exact closed typed APIs —
`PrepareMutationCarrierSet(AuthorizedMutationPlan | AutomaticMutationPlan)` becomes the distinct
`evidence.PrepareAuthorizedMutationCarrierSet(evidence.AuthorizedMutationPlan)` and
`evidence.PrepareAutomaticMutationCarrierSet(evidence.AutomaticMutationPlan)`, and
`ValidateExceptionTerminalDescriptorSet([]MutationDescriptor | DomainConclusionDescriptor, …)` becomes
the distinct `evidence.ValidateExceptionTerminalGrantDescriptorSet(…, []evidence.MutationDescriptor,
…)` and `evidence.ValidateExceptionTerminalDenyDescriptor(…, evidence.DomainConclusionDescriptor, …)`,
preserving the distinct authorized, automatic, Grant-descriptor and Deny-conclusion paths with no
`any`, open interface, or implementation choice left to tasks; (3) it extends
`F18-ARCH-009`/`F18-ARCH-012`, their protected-symbol inventory and fixtures/tests so that only the
matching owner `FinalizeAt` may construct a nonzero registered concrete finalized-change value —
rejecting owner-local helper construction, semantic cloning, or copied still-unconsumed
`PermitBinding` embedding into that same concrete type outside that exact `FinalizeAt` — while stating
that ordinary immutable transfers (returning the value from `FinalizeAt`, passing it by value
into/out of `uow`, submitting it unchanged to `ApplyFinalized`, and reading its accessors) remain
permitted; (4) it makes `CallerReservationLease.Abort` and `CallerMutationTransaction.Abort` exact
and internally consistent, specifying the `stateMu -> reservationMu` acquisition order, terminal
outcome recording and reservation removal under the locks, capability invalidation, unlock order, and
notification strictly after both locks are released per the single authoritative
waiter-notification protocol, with no wait, `Signal`, `Broadcast`, or waiter wake while either lock
is held, waiters gating on the recorded terminal outcome for lost-wakeup safety, and no publication
or consumed sequence on abort/CAS-loss/cancellation/evidence-failure/panic
cleanup; and (5) it adds ADH-2026-073 to the §9 controlling-handoff list and corrects §9.5 so
`VS0-CF-F18-20` retains its registered AC-F18-20 (Assignment-changes-during-review) meaning,
replacing its prior audit-before-result/accepted-evidence repurposing with focused design-mechanic
tests where no registered case directly proves the mechanic. FEATURE-0013
remains contract/validation-only (ADH-2026-072 §1; ADH-2026-074 summary); the
candidate-to-completed-evaluation point is
closed by ADH-2026-072 §2. No stop condition is active.

**Review classification.** Every subsequent FEATURE-0018 design finding must be exactly one of
`STALE_TRANSCRIPTION`, `DESIGN_EXECUTABILITY`, `REQUIREMENT_GAP`,
`ARCHITECTURE_CLARIFICATION_REQUIRED`, or `OUT_OF_SCOPE`, using the definitions in
ADH-2026-072 §5. A reviewer must cite the exact contradictory/missing authority and target text
before reopening a closed decision. An incomplete API, lock, receipt, checker or test mechanism is
`DESIGN_EXECUTABILITY`, not an architecture gap.

<!-- END RETAINED SEMANTIC EXCERPT: raw design.md lines 3019-3471 -->
