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
accepted_decision: DEC-0060
change_request: ACR-2026-002
sole_architecture_authority: docs/architecture/FEATURE-0018-governance-iam-approval-exception-foundation.md
requirements_authority: .kiro/specs/governance-iam-approval-exception-foundation/requirements.md
ai_load_priority: feature
ai_summary: FEATURE-0018 implementable design preserving approved resources, routes, ownership, mappings and exclusions. Stages 7-9 produce only a sealed non-authoritative candidate with a complete dependency-version claim; state proves admission currentness and authzeval owns publication-time candidate revalidation. Caller and authorization-bearing controller mutations produce audited evaluations, while non-assignable automatic effects publish mutation evidence only. Mutation owners return typed domain-versus-mechanical finalization outcomes and expose neutral result material. Context-bound proofs and runtime Seal checks enforce atomic publication; the static checker covers only imports, declarations and typed calls. Evidence alone constructs required carrier sets, atomic commit precedes result release, and Phase 2R wires deterministic clocks only.
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
`CanonicalResult() (operation.ResultMaterial, bool)` and `ApplyTo(state.StateEditor) error`. A
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

The approved workflow-trigger, RoleAssignment-writer, Membership-writer, approval-derivation and
FEATURE-0017 seam ownership remains exactly as registered by F18-RD-08/10/11/13/14/16 and the
§3.3 writer ledger; this design represents those authorities with narrow typed ports rather than
restating their behavioral contracts. The compiler enforces sealed-type visibility; the named
DD-13 checker enforces imports and authorized call sites. No external dependency is introduced.

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
`state.Store` captures the authoritative `publicationInstant` under `stateMu`; §5.4 defines its
revalidation and reuse. (Delegating: REQ-F18-22; governing: OPS-01.)

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

The exact API, lock order, mode-specific seal/termination rules, panic/cancellation/shutdown
lifecycle and commit visibility order are defined once in §5.4. The mechanism is deterministic,
process-local and non-durable; `## FEATURE-0013 Adoption` is authoritative for the unchanged
dependency boundary. (Governing: REQ-F18-20/21, F18-RD-20/21.)

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

`InFlightTable` is keyed by `IdempotencyLookupKey`; the committed aggregate root stores
`CompletedRecord` values under the same key. Every reservation and completed record contains the
full `RequestBinding`. Target, scope and digest are therefore compared **after** lookup and are
never part of the lookup key. This makes reuse of one caller key for another request detectable.

This table applies only to caller-keyed mutations. Controller-owned operations have no
caller key, reservation or completed replay result; `BeginControllerMutation` reuses inherited
resource-version/transition CAS so a repeated reconciliation is the registered effective-once,
same-decision or no-op outcome. No synthetic caller key or second idempotency mechanism is created.

Every replay and woken waiter restarts the ordered pipeline and reaches stage 10 with a fresh sealed
`OperationCandidate` and safe visibility finding before any stored result can be selected. The
stage-11 evidence set is then accepted before the stored result is returned. Completed-record
retention honours the 24h horizon computed from the injected clock and never extends JIT validity.
The typed `IdempotencyLookupKey`,
`RequestBinding`, `InFlightReservation` and `CompletedRecord` shapes, pure lifecycle helpers, and
replay-recheck ports live in the neutral child package
`internal/govaccess/idempotency`, which imports only `model` and `apivalid`; `state` stores those
values, and `uow` is the sole lifecycle caller. §5.4 is authoritative for inspection, waiting,
reservation ownership and publication ordering. (Governing: REQ-F18-13, REQ-F18-21; OPS-03;
FEATURE-0012 §6.12.)

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
commit. Root only wires these components. The exact denied, read, replay and mutation flows are
defined once in §5.5. FEATURE-0013 supplies only carrier and validation contracts
(ADH-2026-072 §§1–2; REQ-F18-20/21).

### DD-13 — Exactly what the compiler enforces vs. what the supplemental static package check enforces

The isolation claims of DD-01/§3.4 are separated into two enforcement mechanisms, and no claim
is stronger than the mechanism that backs it. Neither mechanism is `VS0-CF-F18-50`: that registered
case is the current-feature/prohibited-concept architecture fixture (REQ-F18-01) and is **not**
redefined by this design (see the supplemental-check note below and §7.1).

- **What the Go compiler enforces.** Owner-private implementations behind sealed exported
  interfaces prevent another package from constructing a domain `PreparedIntent` or
  `FinalizedPreparedChange`, evidence `PreparedEvidenceChange`, `CallerReservationLease`,
  transaction, `MutationFinalizationPermit`, `FinalizationPermitToken`, `PermitBinding`,
  `AppliedChangeReceipt` or `StateEditor`. Only owner intent interfaces expose their owner-package
  `FinalizeAt(state.MutationFinalizationPermit)` operation, and a finalized change can only embed a
  `state.PermitBinding` that a transaction minted and handed to that exact `FinalizeAt` call, so the
  binding is unforgeable by construction; only finalized changes expose `PublicationKind`,
  `PermitBinding`, the canonical ordered `EvidenceDescriptors()` set (for a `ResourceMutation`), `ConclusionDescriptor` (for a
  `DomainConclusion`), `DomainDecisionMaterials`, `CanonicalResult` and `ApplyTo(StateEditor)`. Evidence exposes read-only
  validation/linkage fields and no editor port. Mutation transactions expose `ApplyFinalized`, not
  an editor; only a state-owned transaction can create a receipt, and only `evidence` can create the
  receipt's sealed applied proof. `uow` rejects nil/zero sealed values. Package
  privacy protects the live root, both mutexes, the permit-issuance registry and reservation table. Neither the compiler nor the
  supplemental checker proves dynamic call-path order, context identity, version equality,
  permit-issuance provenance or Seal cardinality; the sealed runtime protocol and focused tests do.
- **What the named supplemental checker enforces (import graph + typed call graph).** The exact
  standard-library-only artifact is:

  ```text
  source:            scripts/feature0018-architecture-check/main.go
  tests:             scripts/feature0018-architecture-check/main_test.go
  valid fixture:     scripts/feature0018-architecture-check/testdata/valid/
  invalid fixtures:  scripts/feature0018-architecture-check/testdata/invalid/
  command:           make feature-0018-architecture-check
  feature gate:      make ff-feature-gate FEATURE=FEATURE-0018
  ```

  It uses `go/build`, `go/parser`, `go/ast`, `go/token`, `go/types` and `go/importer`; package and
  file traversal is sorted, network access is forbidden, and any unresolved package or type makes
  the check fail closed. It protects these exact design symbols:

  ```text
  state.InspectOrReserveCallerMutation
  state.ExpectedVersionPredicate
  state.ExpectedVersionSet
  state.NewExpectedVersionSet
  state.AuthorizationDependencyVersionSet
  state.NewAuthorizationDependencyVersionSet
  state.AuthorizationCurrentnessClaim
  state.NewAuthorizationCurrentnessClaim
  evidence.AuthorizationAdmissionProof
  evidence.BindAuthorizationAdmission
  evidence.CurrentAllowProof
  evidence.RequireCurrentAllow
  evidence.RouteDenyToEvidenceOnly
  evidence.DomainDecisionMaterial
  evidence.NewApprovalDecisionMaterial
  evidence.NewExceptionDecisionMaterial
  model.MutationEventFacts
  operation.AppliedChangeLink
  operation.NewAppliedChangeLink
  operation.AppliedConclusionLink
  operation.NewAppliedConclusionLink
  operation.ParticipantBinding
  operation.PublicationBinding
  operation.ShadowDeltaDigest
  operation.IndexDeltaDigest
  operation.MutationDigest
  operation.CompletionBinding
  operation.NewCompletionBinding
  state.ParticipantOwner
  state.RoleAssignmentOwner
  state.MembershipOwner
  state.RoleDefinitionOwner
  state.AccessGroupOwner
  state.ApprovalOwner
  state.PrivilegedOwner
  state.ReviewOwner
  state.ExceptionOwner
  state.GovernanceProfileOwner
  state.ParticipantID
  state.IntentDigest
  state.ResultRole
  state.OriginatingResult
  state.SupportingResult
  state.ParticipantClaim
  state.NewParticipantClaim
  state.MutationMode
  state.CallerMutationMode
  state.ControllerMutationMode
  state.MutationAuthorityMode
  state.AuthorizationBearing
  state.AutomaticMutationOnly
  state.MutationAdmission
  state.NewCallerMutationAdmission
  state.NewAuthorizedControllerMutationAdmission
  state.NewAutomaticControllerMutationAdmission
  state.ContextBoundChange
  state.ContextBoundChange.PublicationKind
  state.ContextBoundChange.PermitBinding
  state.ContextBoundChange.EvidenceDescriptors
  state.ContextBoundChange.ConclusionDescriptor
  state.PublicationKind
  state.ResourceMutation
  state.DomainConclusion
  state.FinalizationPermitToken
  state.PermitBinding
  state.MutationFinalizationPermit
  state.MutationFinalizationPermit.Binding
  state.AppliedChangeReceipt
  state.AppliedChangeReceipt.EvidenceProof
  state.AppliedChangeReceipt.ConclusionProof
  state.AppliedChangeReceipt.PermitBinding
  state.AppliedChangeReceipt.PublicationKind
  state.CommitReceipt
  state.CallerReservationLease
  state.CallerReservationLease.Begin
  state.CallerReservationLease.Abort
  state.CallerMutationTransaction
  state.CallerMutationTransaction.PublicationContext
  state.CallerMutationTransaction.AuthorizationAdmissionProof
  state.CallerMutationTransaction.AcceptCurrentAllow
  state.CallerMutationTransaction.FinalizationPermit
  state.CallerMutationTransaction.ApplyFinalized
  state.CallerMutationTransaction.StageCompletedResult
  state.CallerMutationTransaction.AcceptPreparedEvidence
  state.CallerMutationTransaction.Seal
  state.CallerMutationTransaction.Commit
  state.CallerMutationTransaction.Abort
  state.BeginControllerMutation
  state.ControllerMutationTransaction
  state.ControllerMutationTransaction.PublicationContext
  state.ControllerMutationTransaction.AuthorizationAdmissionProof
  state.ControllerMutationTransaction.AcceptCurrentAllow
  state.ControllerMutationTransaction.FinalizationPermit
  state.ControllerMutationTransaction.ApplyFinalized
  state.ControllerMutationTransaction.AcceptPreparedEvidence
  state.ControllerMutationTransaction.Seal
  state.ControllerMutationTransaction.Commit
  state.ControllerMutationTransaction.Abort
  state.BeginEvidenceTransaction
  state.EvidenceTransaction
  state.EvidenceTransaction.PublicationContext
  state.EvidenceTransaction.AuthorizationAdmissionProof
  state.EvidenceTransaction.AcceptPreparedEvidence
  state.EvidenceTransaction.Seal
  state.EvidenceTransaction.Commit
  state.EvidenceTransaction.Abort
  state.StateEditor
  roleassign.PreparedRoleAssignmentIntent
  roleassign.newPreparedRoleAssignmentIntent
  roleassign.FinalizedRoleAssignmentChange
  roleassign.PreparedRoleAssignmentIntent.FinalizeAt
  roleassign.FinalizedRoleAssignmentChange.EvidenceDescriptors
  roleassign.FinalizedRoleAssignmentChange.DomainDecisionMaterials
  roleassign.FinalizedRoleAssignmentChange.CanonicalResult
  roleassign.FinalizedRoleAssignmentChange.ApplyTo
  operation.ResultMaterial
  operation.NewResultMaterial
  operation.FinalizationOutcome[T, D]
  operation.Finalized
  operation.Domain
  operation.MechanicalFailure
  evidence.AuthorizationCandidateDescriptor
  evidence.NewAuthorizationCandidateDescriptor
  evidence.FinalizedAuthorizationCandidate
  evidence.NewFinalizedAuthorizationCandidate
  evidence.CompletedAuthorizationEvaluation
  evidence.CompleteAuthorizationCandidate
  evidence.PublicationContextBinding
  evidence.PreparedEvidenceChange
  evidence.NewMutationDescriptor
  evidence.ValidateExceptionTerminalDescriptorSet
  evidence.DomainConclusionDescriptor
  evidence.NewDomainConclusionDescriptor
  evidence.AppliedMutationProof
  evidence.AppliedConclusionProof
  evidence.BindAppliedChange
  evidence.BindAppliedConclusion
  evidence.NewAuthorizedMutationPlan
  evidence.NewAutomaticMutationPlan
  evidence.NewDomainConclusionPlan
  evidence.AuthorizedMutationPlan
  evidence.AutomaticMutationPlan
  evidence.DomainConclusionPlan
  evidence.PrepareMutationCarrierSet
  evidence.PrepareAuthorizationCarrierSet
  evidence.PrepareDomainConclusionCarrierSet
  idempotency.PreparedCompletedResult
  idempotency.PrepareCompletedResult
  authzeval.AuthorizationCandidate
  authzeval.OperationCandidate
  evidence.FinalizedAuthorizationCandidate
  authzeval.FinalizeCandidateAt
  authzeval.ReleaseAfterPublication
  roleassign.GrantPort
  roleassign.ActivationPort
  roleassign.ReviewRemediationPort
  roleassign.DirectRevokePort
  command.RoleAssignmentRevokeService
  uow.MutationCoordinatorPort
  ```

  Rule 009 additionally protects this closed owner-finalization table; for every row it protects
  the intent's `FinalizeAt(state.MutationFinalizationPermit)` and the finalized change's
  `PublicationKind`, `PermitBinding`, `EvidenceDescriptors`, `ConclusionDescriptor`, `CanonicalResult`
  and `ApplyTo` methods:

  | Owner package | Exact claim-owner constant | Sealed intent | Sealed finalized change |
  |---|---|---|---|
  | `roleassign` | `state.RoleAssignmentOwner` | `PreparedRoleAssignmentIntent` | `FinalizedRoleAssignmentChange` |
  | `membership` | `state.MembershipOwner` | `PreparedMembershipIntent` | `FinalizedMembershipChange` |
  | `roledefinition` | `state.RoleDefinitionOwner` | `PreparedRoleDefinitionIntent` | `FinalizedRoleDefinitionChange` |
  | `accessgroup` | `state.AccessGroupOwner` | `PreparedAccessGroupIntent` | `FinalizedAccessGroupChange` |
  | `approval` | `state.ApprovalOwner` | `PreparedApprovalIntent` | `FinalizedApprovalChange` |
  | `privileged` | `state.PrivilegedOwner` | `PreparedPrivilegedIntent` | `FinalizedPrivilegedChange` |
  | `review` | `state.ReviewOwner` | `PreparedReviewIntent` | `FinalizedReviewChange` |
  | `exception` | `state.ExceptionOwner` | `PreparedExceptionIntent` | `FinalizedExceptionChange` |
  | `governanceprofile` | `state.GovernanceProfileOwner` | `PreparedGovernanceProfileIntent` | `FinalizedGovernanceProfileChange` |

  The checker emits one stable, checker-internal diagnostic for each rule:

  | Rule | Failure condition |
  |---|---|
  | `F18-ARCH-001_UNAUTHORIZED_COMMIT_CALLER` | A package other than `uow` calls caller/controller/evidence transaction lifecycle methods; a lease/transaction lacks immediate deferred `Abort`; a capability escapes; code waits while holding one; a `Begin` accepts anything other than a sealed `MutationAdmission`; or a path lacks terminal transfer, commit or abort. |
  | `F18-ARCH-002_ROLEASSIGNMENT_WRITER` | A package other than `roleassign` declares, constructs or directly writes a RoleAssignment prepared intent, finalized change, spec or lifecycle state. |
| `F18-ARCH-003_EVIDENCE_IMPORT` | `evidence` imports `state`, `authzeval`, `uow` or root; another package declares/implements candidate descriptors, finalized candidate values, current-Allow proofs, completed evaluations, domain-decision materials, mutation descriptors, domain-conclusion descriptors, applied proofs (mutation or conclusion), digests, carrier plans (mutation or domain-conclusion), IDs, carrier sets or `PreparedEvidenceChange`; the authorization-candidate constructors are invoked outside `authzeval`; `RequireCurrentAllow` is invoked outside `uow`; the approval/exception material constructors, `NewDomainConclusionDescriptor` or `NewMutationDescriptor` are invoked outside their exact semantic-owner finalizer; `BindAppliedChange`/`BindAppliedConclusion` is invoked outside a state transaction; the domain-conclusion plan or carrier-set constructors are invoked outside `uow`; or `uow` constructs a FEATURE-0013 carrier instead of calling the evidence preparation API. |
  | `F18-ARCH-004_EDITOR_ESCAPE` | A domain package returns `StateEditor`; places it in a package variable, struct, composite, map, slice or channel; passes it to a goroutine/defer/closure; assigns it to a longer-lived alias; or calls anything except its registered direct mutation methods inside `ApplyTo`. |
| `F18-ARCH-005_PACKAGE_DIRECTION` | A domain owner imports `uow`, `command`, `idempotency` or root; `state` imports `authzeval`, `uow`, `command` or a semantic owner; `evidence`/`idempotency` imports `state`; `idempotency` imports `evidence`; `operation` imports a non-leaf package; any package other than `policyseam` imports `policyeval`; a package edge is absent from §3.4; a Phase-2R file calls `time.Now` or defines/wires `UTCClock`; or another edge violates §3.4. |
| `F18-ARCH-006_TRIGGER_CALLER` | Any caller other than `grantintent`, `privileged` or `review` invokes `GrantPort`, `ActivationPort` or `ReviewRemediationPort`, respectively; membership submits a RoleAssignment intent; a caller other than `command.RoleAssignmentRevokeService` invokes `DirectRevokePort`; or direct revoke bypasses `uow.MutationCoordinatorPort`. |
  | `F18-ARCH-007_POLICY_SEAM` | A package other than `privileged` calls `policyseam.Port`, or `roleassign` derives approval/imports `approvalreq`. |
| `F18-ARCH-008_EVIDENCE_PUBLISHER` | A package other than `uow` invokes `RequireCurrentAllow`, `CompleteAuthorizationCandidate`, evidence carrier preparation, completed-result staging, transaction acceptance or `authzeval.ReleaseAfterPublication`; a `ReleaseAfterPublication` call lacks the sealed `state.CommitReceipt` returned by the applicable transaction's `Commit`; a pre-stage-10 publisher symbol is declared; or root contains anything beyond constructor/bundle wiring. This static rule does not claim that runtime evidence/commit ordering or current-Allow validation succeeded. |
| `F18-ARCH-009_PUBLICATION_FINALIZATION` | An owner intent declares publication-derived state/digest/result/descriptor/decision-material/permit-binding fields; a finalizer call site is outside `uow` or is not passed a transaction-issued `MutationFinalizationPermit`; a finalized change omits `PublicationKind`, `PermitBinding`, or `EvidenceDescriptors()` when it declares `ResourceMutation`; `uow`/`state` declares or constructs an owner finalized-change concrete type; an owner finalizer calls clock/state/lock/I/O APIs; `NewResultMaterial` is called outside the originating owner's `FinalizeAt`; the approval/exception decision-material constructor, `NewMutationDescriptor`, `ValidateExceptionTerminalDescriptorSet`, or `NewDomainConclusionDescriptor` is called outside its corresponding owner `FinalizeAt`; `ApplyTo` is called directly rather than by state transaction implementation; or an owner imports/exposes an idempotency-owned result instead of neutral `operation.ResultMaterial`. This static rule does not prove runtime context, current-Allow, permit-provenance, descriptor cardinality, decision-material or outcome equality. |
| `F18-ARCH-010_APPLIED_CHANGE_RECEIPT` | A non-owner or an owner using another owner's constant calls `NewExpectedVersionSet`/`NewParticipantClaim`; a package other than `state` calls `operation.NewAppliedChangeLink` or `operation.NewAppliedConclusionLink`, or declares/constructs `AppliedChangeReceipt`; a package other than `evidence` calls `operation.NewCompletionBinding`; `uow` or an owner obtains/uses `StateEditor` directly; or a package other than `uow` calls mutation `AcceptCurrentAllow`/`FinalizationPermit`/`ApplyFinalized`/`Seal`. Runtime claim/context/current-Allow/permit-provenance/decision-material/digest/delta/cardinality equality is excluded from this checker rule. |
| `F18-ARCH-011_AUTHORIZATION_CURRENTNESS` | A package other than `authzeval` calls `NewAuthorizationDependencyVersionSet` or `NewAuthorizationCurrentnessClaim`; a package other than `uow` calls `authzeval.FinalizeCandidateAt` or one of the three mode-specific admission constructors; a package other than `state` calls `BindAuthorizationAdmission`; an automatic constructor call supplies a currentness value; an authorized constructor call lacks one; an authorization-bearing owner finalizer is reachable before `AcceptCurrentAllow` and `FinalizationPermit`; or `uow` calls a semantic time-validation function directly. This rule checks declarations and typed calls, not runtime version, Allow result or time equality. |

  Invalid fixtures each violate exactly one rule. Rule 001 fixtures include foreign inspection,
  lease/transaction acquisition without immediate deferred abort, capability escape, waiting with
  a live capability, descriptor-bearing `Begin`, and a missing terminal path. Rule 003 includes a
  `uow` carrier literal, foreign evidence implementation, approval/exception material or descriptor
  construction outside the exact owner finalizer, plus `BindAppliedChange`/
  `BindAuthorizationAdmission` outside `state`. Rule 008 has isolated pre-stage
  publisher/completed-evaluation, non-`uow` completion/preparation/acceptance,
  foreign current-Allow construction, early-result-release call site and root-invocation fixtures.
  Rule 009 has isolated
  publication-derived intent field, foreign finalizer caller, owner clock/state read,
  finalization without a transaction permit, a finalized change omitting `PublicationKind`/
  `PermitBinding`, a foreign `NewDomainConclusionDescriptor` caller, `uow`/`state`
  semantic-construction, direct `ApplyTo`, and owner-to-idempotency-result fixtures. Rule 010 has
  isolated foreign neutral-link (`NewAppliedChangeLink`/`NewAppliedConclusionLink`)/receipt/
  completion-binding construction, editor extraction, foreign `AcceptCurrentAllow`,
  `FinalizationPermit`, `ApplyFinalized` and `Seal` fixtures. Rule 011 has isolated foreign currentness
  construction/finalization, wrong mode-specific admission-constructor use and direct-uow semantic-recheck
  fixtures, plus an authorization-bearing finalizer ordered before the sealed current-Allow gate.
  Valid fixtures prove the permitted imports and typed call sites for waiter return,
  caller lease transfer, owner finalization (both `ResourceMutation` and `DomainConclusion`
  publication kinds), transaction-only application, evidence preparation (mutation and
  domain-conclusion carrier plans), neutral linkage, all three transaction modes, both controller
  authority modes, sealed current-Allow admission, post-commit release and root-only wiring. Remaining fixtures
  cover the writer, editor, package-direction, trigger and policy rules without changing their
  registered meaning.

  **Runtime Seal enforcement is separate.** The checker proves only declarations, imports,
  authorized typed calls and the explicitly named local AST patterns above. It does not prove
  runtime version/currentness, publication-context identity, permit-issuance provenance,
  result/digest/delta linkage, receipt cardinality, carrier cardinality, all-path ordering or
  successful Seal. The sealed transaction
  APIs perform those validations at runtime; the focused tests in §7.2 prove every positive and
  negative dynamic invariant. A dynamic rejection is never claimed as a checker fixture.

  These diagnostics are implementation-gate evidence, not public API Problem codes and not new
  VS0 conformance identifiers. `VS0-CF-F18-50` retains only its registered REQ-F18-01
  current-feature/prohibited-concept semantics. (Governing: SEC-05, F18-RD-10/11/13/14/16.)

---

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
| `internal/govaccess/` (root, composition/wiring only) | Controller wiring and composition **only**: constructs child packages; injects the single deterministic clock into `state.Store` and time-dependent owners; wires caller command services, the pure evaluator, `uow` coordinator and post-commit result releaser without invoking them; and loads the exact local FEATURE-0013 adoption bundle through existing `decision/bundle.Load`. Constructs no record, holds no live state, opens no transaction, obtains no editor and exposes no pre-stage-10 publisher | child owner packages, `command`, `uow`, `authzeval`, `clock`, `decision/bundle` |
| `internal/govaccess/model/` | Go structs and closed enums for the ten contracts and the supporting values; the four distinct FEATURE-0018-owned reference values `PrincipalRef` (closed-category identity value, `VS0-SCHEMA-020`), `AccessGroupRef` (UID-pinned `AccessGroup` reference embedding `apimeta.TypedRef` pinned to `Kind==AccessGroup`), `RoleHolderRef` (closed `PrincipalRef`-or-`AccessGroupRef` discriminated union), and `EligibilityRef` (exactly one Human `PrincipalRef`); plus `ActionTargetBinding`; the closed enums of requirements §3.3; and immutable deep-copied `ApprovalDecisionFacts`, `ExceptionDecisionFacts`, and owner-finalized `MutationEventFacts` projections. The first two contain exactly the corresponding F18-RD-19 profile inputs and terminal result; `MutationEventFacts` contains only the already-registered event taxonomy, resource/version linkage, action and after-state digest. These projections are design-level carrier inputs, not resources, decisions or public contracts. Only `AccessGroupRef` reuses `apimeta.TypedRef`; the other three do not (DD-03) | `apimeta` (via `AccessGroupRef` only) |
| `internal/govaccess/operation/` (**package**) | Neutral leaf protocol values shared by semantic owners, state, evidence and idempotency: immutable deep-copied `ResultMaterial`; sealed mechanical `AppliedChangeLink`, `AppliedConclusionLink` and `CompletionBinding`; plus generic `FinalizationOutcome[T,D]` with exactly one typed `Finalized(T)`, `Domain(D)` or `MechanicalFailure` variant, where the `Finalized(T)` change declares a closed `PublicationKind` (`ResourceMutation` or `DomainConclusion`). `state` alone invokes `NewAppliedChangeLink` from its private admitted claim, exact publication binding and actual shadow/index deltas; `evidence` alone invokes `NewCompletionBinding` from its validated plan's publication and mutation digests. Neither neutral value contains authorization, approval, exception or lifecycle meaning. Each owner supplies its own already-registered domain/lifecycle outcome type as `D`; `MechanicalFailure` can identify only local invariant, evidence-construction or publication-mechanism failure. It defines no public error, lifecycle or domain semantics | `model` |
| `internal/govaccess/clock/` (**package**) | True leaf: `Clock` (`NowUTC() time.Time`) plus deterministic `FixedClock` and manually advanced fixture clock only. No `UTCClock`, `time.Now` or production wall-clock wiring is implemented in Phase 2R. The root injects one instance into `state.Store` and owner-local clock ports (DD-05) | — |
| `internal/govaccess/validate/` | Deterministic per-contract structural + local-semantic + cross-field + scope/reference validators, ordered per the F18-RD-21 precedence (§5.1); emits inherited Problem/violation codes only | `apivalid`, `apiref`, `apiproblem`, `apicond` |
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
| `internal/govaccess/state/` (**package**) | Neutral non-durable storage. Owns root/shadows, `reservationMu`, `stateMu`, publication sequence, in-flight/completed tables, accepted carrier storage, sealed dependency-version sets/currentness claims, `MutationAdmission`, transaction-issued `MutationFinalizationPermit` with its transaction-private permit-issuance registry and unforgeable `PermitBinding`, the closed `PublicationKind`, the transaction-private `AppliedChangeReceipt` registry, and the sealed post-install `CommitReceipt`. Exposes caller inspection/lease, controller begin, evidence-only begin, and the three sealed transaction interfaces in §5.4. Every successful authorization-bearing begin revalidates the complete admitted dependency set, captures the deterministic context and asks `evidence.BindAuthorizationAdmission` for one proof binding both while still under `stateMu`. An authorization-bearing transaction issues no finalization permit until `AcceptCurrentAllow` validates the exact evidence-owned proof; automatic mode issues only an automatic permit. `ApplyFinalized` first proves the change's `PermitBinding` against its permit-issuance registry (matching unconsumed token, equal participant claim, exact context) and consumes that entry, then invokes the admitted finalized change against the private editor; for a `ResourceMutation` it constructs a neutral `operation.AppliedChangeLink` from the admitted claim/context and non-empty actual shadow/index deltas and asks `evidence.BindAppliedChange` to create the proof and receipt, and for a `DomainConclusion` it proves an empty delta, constructs `operation.AppliedConclusionLink`, and asks `evidence.BindAppliedConclusion`. `Commit` installs exactly one sealed root/evidence state and then returns its opaque `CommitReceipt`; state never constructs/finalizes domain meaning, performs semantic time validation, or constructs descriptors, canonical results, domain-decision materials or carriers | `model`, `operation`, `idempotency`, `decision`, `evidence`, `clock` |
| `internal/govaccess/uow/` (**package**) | Sole stage-10/11/12 coordinator and implementation of the narrow `MutationCoordinatorPort` consumed by `command`. It carries owner-sealed intents and, for authorization-bearing paths, the candidate/currentness claim into stage 10; manages caller leases and caller/controller/evidence transactions; passes the transaction's exact admission proof/context to `authzeval.FinalizeCandidateAt`; asks `evidence.RequireCurrentAllow` for the unforgeable mutation gate; submits that proof to `AcceptCurrentAllow`; obtains an owner-bound transaction permit; and only then calls each semantic owner and submits finalized changes through `ApplyFinalized`. A Deny aborts this transaction and is completed/audited only through a fresh evidence-only transaction. It then requests evidence-owned completed-evaluation proofs and the applicable carrier plan, plus idempotency-owned completed results through the plan's neutral `operation.CompletionBinding` where applicable. It accepts sealed values, commits, receives (but cannot construct) the state transaction's `CommitReceipt`, and only then passes that receipt to `authzeval.ReleaseAfterPublication` when a completed evaluation exists. It constructs or mutates no semantic value, editor, permit, receipt, descriptor, decision material, result or carrier and performs no semantic revalidation. `Coordinator` owns cancellation/admission/waiting; capabilities remain stack-local | `state`, child owner packages, `authzeval`, `operation`, `idempotency`, `evidence`, `decision` |
| `internal/govaccess/evidence/` (**package**) | Sole constructor/validator of `AuthorizationCandidateDescriptor`, `FinalizedAuthorizationCandidate`, `AuthorizationAdmissionProof`, sealed `CurrentAllowProof`, context-bound `CompletedAuthorizationEvaluation`, profile-specific sealed `DomainDecisionMaterial`, canonical ordered `MutationDescriptor` sets, `DomainConclusionDescriptor`, `AppliedMutationProof`, `AppliedConclusionProof`, canonical `MutationDigest`, the distinct sealed `AuthorizedMutationPlan`/`AutomaticMutationPlan`/`DomainConclusionPlan` values, deterministic carrier IDs, exact `CarrierSet`, and sealed `PreparedEvidenceChange` under consumed writer `VS0-WRITER-011`. `state` alone may request an admission proof after version revalidation; `authzeval.FinalizeCandidateAt` alone may request a finalized-candidate value after semantic revalidation; only `uow` may request `RequireCurrentAllow`; and only the `approval`/`exception` owner finalizers may request their corresponding decision material from exact immutable `model.ApprovalDecisionFacts`/`model.ExceptionDecisionFacts`. Only the exception owner finalizer may request terminal-exception descriptors: its Grant set is exactly `exceptionproposal.decided`, then `exceptiongrant.issued`; its Deny conclusion descriptor is exactly `exceptionproposal.decided`. `ValidateExceptionTerminalDescriptorSet` binds each to the exact immutable proposal, terminal reason, mandatory `exception-decision/v1` record material and containing operation/correlation evidence using the existing FEATURE-0013 envelope/linkage semantics. `BindAppliedChange`/`BindAppliedConclusion` consume only neutral `operation.AppliedChangeLink`/`operation.AppliedConclusionLink` plus evidence-owned descriptors/materials—never a state type. Authorized and automatic plans independently enforce all applicable F18-RD-19 mandatory domain DecisionRecords; authorization evidence is additive and never their condition. Plans expose only a neutral `operation.CompletionBinding` to idempotency. It has no editor and imports neither `state`, `authzeval`, `idempotency` nor `uow` | `decision`, `model`, `operation` |
| `internal/govaccess/policyseam/` (**package**) | **Sole** FEATURE-0017 consumer for the pinned PrivilegedAccessRequest subject; exposes `policyseam.Port`; the only package importing `internal/policyeval`; `Indeterminate`→`Deny` mapping (DD-01) | `policyeval` |
| `internal/govaccess/fixtures/` (**package**) | Immutable synthetic fixture builders + graph validator + deterministic fake executor (DD-06) | `model`, `clock` |
| `internal/govaccess/idempotency/` (**package**) | Neutral owner of `IdempotencyLookupKey`, `RequestBinding`, `InFlightReservation`, sealed `PreparedCompletedResult`, `CompletedRecord`, and pure lifecycle helpers. `PrepareCompletedResult` wraps the owner-supplied neutral `operation.ResultMaterial` and binds it to the exact lookup/request binding plus evidence-validated neutral `operation.CompletionBinding`; it never imports or names an evidence type. These types apply only to caller-keyed mutations; controller mutations use existing transition/version CAS and no synthetic key | `apivalid`, `model`, `operation` |
| `scripts/feature0018-architecture-check/` | Named standard-library static checker, unit tests and positive/negative fixtures for DD-13; invoked by `make feature-0018-architecture-check` and the FEATURE-0018 feature gate | Go parser/type checker |
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

The authoritative package rows are §3.2; this section records only design-owned edges and
enforcement:

```text
evidence -> model, decision, operation
operation -> model
idempotency -> model, apivalid, operation
state -> model, operation, idempotency, decision, evidence, clock
authzeval -> model, decision, operation, evidence, state, clock
domain owners -> leaf ports, operation, evidence (descriptor constructor), state (sealed claim/finalization capability values and StateEditor type only)
uow -> domain owners, operation, state, idempotency, evidence, authzeval
command -> authzeval, roleassign, uow
root -> domain owners, command, uow, clock, decision/bundle
```

There is no reverse edge: `evidence`/`idempotency` do not import `state`; owners do not import
`uow`, `command` or root; `state` does not import `authzeval`, `uow`, `command` or an owner; and
nothing domain-local imports root. The `authzeval -> state` edge uses only sealed neutral
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
`FinalizeAt`. It may not import or name `Store`, any transaction, lease, reservation, mutex, root,
editor-acquisition API, receipt constructor, admission constructor, `ApplyFinalized`, `Seal`,
`Commit`, or `Abort`. Thus the sealed owner-finalization API is usable without granting a domain
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
  the committed state transaction; no pre-stage-10 publisher exists.

These are representations of F18-RD-08/10/11/13/14/16 and the §3.3 writer ledger, not duplicate
semantic definitions. DD-13 checks the exact imports and call sites; the compiler enforces sealed
implementations and private state. The resulting graph is acyclic.

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
action form `/{uid}/actions/<action>` read via `r.PathValue("uid")`). No generic `PUT`,
unrestricted `PATCH`, or hard `DELETE` (REQ-F18-02):

```text
Collection : POST | GET   /apis/<group>/v1alpha1/<plural-kebab>
Item       : GET          /apis/<group>/v1alpha1/<plural-kebab>/{uid}
Draft PATCH: PATCH        /apis/<group>/v1alpha1/<plural-kebab>/{uid}   (Draft only, VersionedDefinition kinds)
Action     : POST         /apis/<group>/v1alpha1/<plural-kebab>/{uid}/actions/<action>
```

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

| Kind (group) | Collection POST / GET | Item GET | Draft PATCH (update) | Item actions (POST `/{uid}/actions/<action>`) | Appendix B caller operations covered |
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
| ExceptionGrant (`governance.sovrunn.io`) | GET only `/apis/governance.sovrunn.io/v1alpha1/exception-grants` (append-only; controller issues evidence) | `…/{uid}` | — | `propose` (collection action — `exceptiongrant.propose`; submit immutable `ExceptionProposal`), `revoke` (submit exact revoke intent; controller issues linked immutable evidence) | propose *(collection action — see below)*, get, list, revoke |

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

### 5.4 In-memory atomic publication protocol (DD-02, DD-10, DD-11, DD-12)

This section is authoritative for the FEATURE-0018-local, deterministic and non-durable
publication mechanics. It does not define a FEATURE-0013 service.

**Exact sealed API.** Only `uow` may call this lifecycle surface. Domain owners provide their
owner-sealed two-phase values; `evidence` and `idempotency` provide the sealed inputs named below:

```text
operation.NewResultMaterial(canonicalSuccessfulResultBytes)
  -> operation.ResultMaterial | operation.MechanicalFailure

state.NewExpectedVersionSet([]state.ExpectedVersionPredicate)
  -> state.ExpectedVersionSet | error

state.NewAuthorizationDependencyVersionSet(
  expectedSnapshotVersion,
  []state.ExpectedVersionPredicate
) -> state.AuthorizationDependencyVersionSet | error

state.NewAuthorizationCurrentnessClaim(
  evidence.CandidateDigest,
  state.AuthorizationDependencyVersionSet
) -> state.AuthorizationCurrentnessClaim | error

state.NewParticipantClaim(
  state.ParticipantOwner,
  state.ParticipantID,
  state.IntentDigest,
  state.ExpectedVersionSet,
  state.OriginatingResult | state.SupportingResult
) -> state.ParticipantClaim | error

state.NewCallerMutationAdmission(
  []state.ParticipantClaim,
  state.AuthorizationCurrentnessClaim
) -> state.MutationAdmission | error

state.NewAuthorizedControllerMutationAdmission(
  []state.ParticipantClaim,
  state.AuthorizationCurrentnessClaim
) -> state.MutationAdmission | error

state.NewAutomaticControllerMutationAdmission(
  []state.ParticipantClaim
) -> state.MutationAdmission | error

InspectOrReserveCallerMutation(IdempotencyLookupKey, RequestBinding)
  -> Conflict | Replay(CompletedRecord) | Wait(WaitHandle)
   | Owner(CallerReservationLease)

CallerReservationLease:
  Begin(state.MutationAdmission)
    -> CallerMutationTransaction | Conflict
  Abort()

CallerMutationTransaction:
  PublicationContext() evidence.PublicationContext
  AuthorizationAdmissionProof() evidence.AuthorizationAdmissionProof
  AcceptCurrentAllow(evidence.CurrentAllowProof)
    -> nil | operation.MechanicalFailure
  FinalizationPermit(state.ParticipantClaim)
    -> state.MutationFinalizationPermit | operation.MechanicalFailure
  ApplyFinalized(state.ContextBoundChange)
    -> state.AppliedChangeReceipt | operation.MechanicalFailure
  StageCompletedResult(idempotency.PreparedCompletedResult)
    -> nil | operation.MechanicalFailure
  AcceptPreparedEvidence(evidence.PreparedEvidenceChange)
    -> nil | operation.MechanicalFailure
  Seal() -> nil | operation.MechanicalFailure
  Commit() -> state.CommitReceipt
  Abort()

BeginControllerMutation(state.MutationAdmission)
  -> ControllerMutationTransaction | Conflict | NoOp

ControllerMutationTransaction:
  PublicationContext() evidence.PublicationContext
  AuthorizationAdmissionProof() (evidence.AuthorizationAdmissionProof, bool)
  AcceptCurrentAllow(evidence.CurrentAllowProof)
    -> nil | operation.MechanicalFailure
  FinalizationPermit(state.ParticipantClaim)
    -> state.MutationFinalizationPermit | operation.MechanicalFailure
  ApplyFinalized(state.ContextBoundChange)
    -> state.AppliedChangeReceipt | operation.MechanicalFailure
  AcceptPreparedEvidence(evidence.PreparedEvidenceChange)
    -> nil | operation.MechanicalFailure
  Seal() -> nil | operation.MechanicalFailure
  Commit() -> state.CommitReceipt
  Abort()

BeginEvidenceTransaction(state.AuthorizationCurrentnessClaim)
  -> EvidenceTransaction | RetryEvaluation

EvidenceTransaction:
  PublicationContext() evidence.PublicationContext
  AuthorizationAdmissionProof() evidence.AuthorizationAdmissionProof
  AcceptPreparedEvidence(evidence.PreparedEvidenceChange)
    -> nil | operation.MechanicalFailure
  Seal() -> nil | operation.MechanicalFailure
  Commit() -> state.CommitReceipt
  Abort()

Per semantic-owner package (owner-private implementations):
  PreparedIntent:
    ParticipantClaim() state.ParticipantClaim
    FinalizeAt(state.MutationFinalizationPermit)
      -> operation.FinalizationOutcome[FinalizedPreparedChange, OwnerDomainOutcome]

  FinalizedPreparedChange:
    implements state.ContextBoundChange
    exposes no additional publication method

state.ContextBoundChange:
  ParticipantClaim() state.ParticipantClaim
  ContextBinding() evidence.PublicationContextBinding
  PermitBinding() state.PermitBinding
  PublicationKind() state.PublicationKind
  EvidenceDescriptors() []evidence.MutationDescriptor               // non-empty iff ResourceMutation;
                                                                    // canonical event/resource/version order
  ConclusionDescriptor() (evidence.DomainConclusionDescriptor, bool) // present iff DomainConclusion
  DomainDecisionMaterials() []evidence.DomainDecisionMaterial
  CanonicalResult() (operation.ResultMaterial, bool)
  ApplyTo(state.StateEditor) error   // ResourceMutation applies a non-empty delta;
                                     // DomainConclusion applies an empty delta and creates no resource

state.MutationFinalizationPermit:
  ParticipantClaim() state.ParticipantClaim
  PublicationContext() evidence.PublicationContext
  AuthorityMode() state.MutationAuthorityMode
  Binding() state.PermitBinding

state.PermitBinding:
  (opaque sealed value derived from a transaction-scoped state.FinalizationPermitToken;
   equality-comparable; mintable only by the issuing transaction, embeddable by the owner
   FinalizeAt that received the permit, and verifiable only against that transaction's
   private permit-issuance registry)

state.PublicationKind = ResourceMutation | DomainConclusion  // closed enum, Valid()

state.AppliedChangeReceipt:
  ParticipantClaim() state.ParticipantClaim
  PublicationContext() evidence.PublicationContext
  PermitBinding() state.PermitBinding
  PublicationKind() state.PublicationKind
  EvidenceProof() (evidence.AppliedMutationProof, bool)      // present iff ResourceMutation
  ConclusionProof() (evidence.AppliedConclusionProof, bool)  // present iff DomainConclusion

state.CommitReceipt:
  PublicationContext() evidence.PublicationContext
  AcceptedEvidence() evidence.PreparedEvidenceChange
  (opaque sealed value constructed only by the committing state transaction after its one
   shadow-root/evidence installation; it proves no new semantic result or carrier)

evidence.AuthorizationAdmissionProof:
  CandidateDigest() evidence.CandidateDigest
  DependencyVersionDigest() evidence.DependencyVersionDigest
  PublicationContext() evidence.PublicationContext

evidence.BindAuthorizationAdmission(
  evidence.CandidateDigest,
  evidence.DependencyVersionDigest,
  evidence.PublicationContext
) -> evidence.AuthorizationAdmissionProof

evidence.RequireCurrentAllow(
  evidence.FinalizedAuthorizationCandidate,
  evidence.AuthorizationAdmissionProof,
  evidence.PublicationContext
) -> evidence.CurrentAllowProof | evidence.RouteDenyToEvidenceOnly
   | operation.MechanicalFailure

evidence.CurrentAllowProof:
  CandidateDigest() evidence.CandidateDigest
  AdmissionDigest() evidence.AdmissionDigest
  PublicationContext() evidence.PublicationContext

evidence.NewApprovalDecisionMaterial(
  model.ApprovalDecisionFacts,
  evidence.PublicationContext
) -> evidence.DomainDecisionMaterial | operation.MechanicalFailure

evidence.NewExceptionDecisionMaterial(
  model.ExceptionDecisionFacts,
  evidence.PublicationContext
) -> evidence.DomainDecisionMaterial | operation.MechanicalFailure

evidence.NewMutationDescriptor(
  model.MutationEventFacts,
  evidence.PublicationContext
) -> evidence.MutationDescriptor | operation.MechanicalFailure
  // called only by the exact semantic-owner FinalizeAt; one call per element of
  // its canonical ordered EvidenceDescriptors() set

evidence.NewDomainConclusionDescriptor(
  model.ExceptionDecisionFacts,
  evidence.PublicationContext
) -> evidence.DomainConclusionDescriptor | operation.MechanicalFailure
  // `exceptionproposal.decided` AuditEvent descriptor for a publishable domain conclusion
  // (the terminal ExceptionProposal Deny); called only by the exception owner FinalizeAt

evidence.ValidateExceptionTerminalDescriptorSet(
  model.ExceptionDecisionFacts,
  []evidence.MutationDescriptor | evidence.DomainConclusionDescriptor,
  []evidence.DomainDecisionMaterial,
  evidence.PublicationContext
) -> ok | operation.MechanicalFailure
  // validates the ADH-2026-073 exact Grant/Deny matrix and existing FEATURE-0013 linkage:
  // each descriptor/material binds the exact immutable ExceptionProposal, terminal reason,
  // exception-decision/v1 DecisionRecord and containing operation/correlation evidence.

operation.NewAppliedChangeLink(
  operation.ParticipantBinding,
  operation.PublicationBinding,
  operation.ShadowDeltaDigest,
  operation.IndexDeltaDigest,
  operation.ResultMaterial | none
) -> operation.AppliedChangeLink | operation.MechanicalFailure

evidence.BindAppliedChange(
  operation.AppliedChangeLink,
  []evidence.MutationDescriptor,
  []evidence.DomainDecisionMaterial
) -> evidence.AppliedMutationProof | operation.MechanicalFailure

operation.NewAppliedConclusionLink(
  operation.ParticipantBinding,
  operation.PublicationBinding
) -> operation.AppliedConclusionLink | operation.MechanicalFailure
  // no shadow/index delta digest: a publishable domain conclusion mutates no resource

evidence.BindAppliedConclusion(
  operation.AppliedConclusionLink,
  evidence.DomainConclusionDescriptor,
  []evidence.DomainDecisionMaterial
) -> evidence.AppliedConclusionProof | operation.MechanicalFailure

operation.NewCompletionBinding(
  operation.PublicationBinding,
  operation.MutationDigest
) -> operation.CompletionBinding | operation.MechanicalFailure

idempotency.PrepareCompletedResult(
  idempotency.IdempotencyLookupKey,
  idempotency.RequestBinding,
  operation.ResultMaterial,
  operation.CompletionBinding
) -> idempotency.PreparedCompletedResult | operation.MechanicalFailure

authzeval.FinalizeCandidateAt(
  authzeval.OperationCandidate,
  evidence.AuthorizationAdmissionProof,
  evidence.PublicationContext
) -> operation.FinalizationOutcome[evidence.FinalizedAuthorizationCandidate, authzeval.ReevaluateCandidate]

evidence.NewFinalizedAuthorizationCandidate(
  evidence.AuthorizationCandidateDescriptor,
  evidence.AuthorizationAdmissionProof,
  evidence.PublicationContext
) -> evidence.FinalizedAuthorizationCandidate | operation.MechanicalFailure

evidence.CompleteAuthorizationCandidate(
  evidence.FinalizedAuthorizationCandidate,
  evidence.AuthorizationAdmissionProof,
  evidence.PublicationContext
) -> evidence.CompletedAuthorizationEvaluation | operation.MechanicalFailure

evidence.AuthorizedMutationPlan / evidence.AutomaticMutationPlan / evidence.DomainConclusionPlan:
  (three distinct opaque sealed evidence values. Each is constructed only by its identically named
   New*Plan constructor, carries the exact PublicationContext and proof set, and is consumable only
   by the corresponding Prepare*CarrierSet signature below.)

evidence.NewAuthorizedMutationPlan(
  evidence.PublicationContext,
  evidence.CompletedAuthorizationEvaluation,
  []evidence.AppliedMutationProof
) -> evidence.AuthorizedMutationPlan | operation.MechanicalFailure
  // only uow; authorization-bearing ResourceMutation transaction

evidence.NewAutomaticMutationPlan(
  evidence.PublicationContext,
  []evidence.AppliedMutationProof
) -> evidence.AutomaticMutationPlan | operation.MechanicalFailure
  // only uow; AutomaticMutationOnly ResourceMutation transaction

evidence.NewDomainConclusionPlan(
  evidence.PublicationContext,
  evidence.CompletedAuthorizationEvaluation,
  []evidence.AppliedConclusionProof
) -> evidence.DomainConclusionPlan | operation.MechanicalFailure
  // only uow; publishable DomainConclusion transaction

evidence.PrepareMutationCarrierSet(
  evidence.AuthorizedMutationPlan | evidence.AutomaticMutationPlan
) -> evidence.PreparedEvidenceChange | operation.MechanicalFailure
  // only uow; accepts an AuthorizedMutationPlan or AutomaticMutationPlan only

evidence.PrepareAuthorizationCarrierSet(
  evidence.CompletedAuthorizationEvaluation
) -> evidence.PreparedEvidenceChange | operation.MechanicalFailure
  // only uow; evidence-only denied/read/replay transaction only

evidence.PrepareDomainConclusionCarrierSet(
  evidence.DomainConclusionPlan
) -> evidence.PreparedEvidenceChange | operation.MechanicalFailure
  // only uow; accepts a DomainConclusionPlan only

authzeval.ReleaseAfterPublication(
  evidence.CompletedAuthorizationEvaluation,
  state.CommitReceipt,
  operation.ResultMaterial | idempotency.PreparedCompletedResult
) -> model.AuthorizationResult | operation.MechanicalFailure
  // only uow, only after the transaction's successful Commit; the receipt proves exact accepted
  // carrier publication for the same completed evaluation and returns no new carrier or authority
```

`PreparedIntent` and `FinalizedPreparedChange` above are a required interface pattern independently
sealed in each semantic-owner package, not shared semantic implementations. Each intent constructs
one sealed `state.ParticipantClaim` only through `state.NewParticipantClaim`, containing its stable
participant identity, exact closed semantic-owner identity, immutable pre-publication intent digest,
an `ExpectedVersionSet` produced from the owner-known FEATURE-0012 version predicates through
`state.NewExpectedVersionSet`, and result role `OriginatingResult | SupportingResult`; it contains no
publication-time value. DD-13 permits each registered owner package to call that constructor only
with its own exact owner constant. `uow` passes the complete claim
set to exactly one of the three mode-specific admission constructors. They validate uniqueness,
exactly one originating participant for caller mode, no originating participant for controller
mode, a mandatory currentness claim for caller/authorization-bearing controller mode, and no claim
for automatic-mutation-only mode without inventing domain semantics. No invalid mode/claim
combination is representable through the public constructor surface.
Only the originating semantic owner's `FinalizeAt` may call `operation.NewResultMaterial`, using the
canonical successful-result encoding already fixed for that operation. `CanonicalResult` returns
the neutral sealed `operation.ResultMaterial` with `true` only for that
admitted originating owner; supporting changes return `false`. `state.ParticipantOwner`,
`state.ParticipantID`, `state.IntentDigest`,
`state.ExpectedVersionSet`, `state.ResultRole`, `state.ParticipantClaim`, `state.MutationMode`,
`state.MutationAuthorityMode`, `state.AuthorizationDependencyVersionSet`,
`state.AuthorizationCurrentnessClaim`, `evidence.AuthorizationAdmissionProof`,
`state.MutationAdmission`, `state.MutationFinalizationPermit`, `state.AppliedChangeReceipt`,
`operation.ResultMaterial`, `operation.AppliedChangeLink`, `operation.CompletionBinding`,
`evidence.AuthorizationCandidateDescriptor`, `evidence.CompletedAuthorizationEvaluation`,
`evidence.CurrentAllowProof`, `evidence.DomainDecisionMaterial`, `evidence.MutationDescriptor`,
`evidence.AppliedMutationProof`,
`evidence.PublicationContext`, `evidence.PublicationContextBinding`,
`idempotency.PreparedCompletedResult` and
`evidence.PreparedEvidenceChange` are sealed values. `uow` cannot implement or mutate any of them
or construct their underlying FEATURE-0013 carriers. The concrete lease, transactions, editor,
wait handle and roots remain private or sealed and may not escape the owning `uow` call.

`authzeval` is the only caller of `state.NewAuthorizationDependencyVersionSet` and
`state.NewAuthorizationCurrentnessClaim`. The dependency set contains the coherent snapshot version
and every exact resource/version fact consumed by stages 5–9, including applicable principal,
membership, group, role-definition version, role-assignment, scope/target, policy, approval,
exception and review evidence. It contains no permission semantics; `state` can therefore compare
the complete closed set to its current root without interpreting it. Each `Begin` rechecks that set
under `stateMu` before capturing the context. Success asks
`evidence.BindAuthorizationAdmission` to create a sealed `AuthorizationAdmissionProof` binding the
exact claim digest, revalidated versions and captured context while `stateMu` remains held.
Mutation-only admission forbids this claim and proof. Currentness alone is not mutation authority.
For every caller and authorization-bearing controller transaction, `uow` must next finalize the
candidate through `authzeval`, call `evidence.RequireCurrentAllow`, and submit the resulting exact
sealed `CurrentAllowProof` to `AcceptCurrentAllow`. The constructor validates that the finalized
result is exactly `Allow` and binds candidate digest, admission proof and publication context. It
cannot construct a proof for `Deny`; that result is the internal `RouteDenyToEvidenceOnly`
direction. Until the exact proof is accepted, `FinalizationPermit`, `ApplyFinalized`, carrier-plan
construction, `Seal` and `Commit` all reject. Automatic-mutation-only mode rejects a
`CurrentAllowProof` and issues only an automatic permit. Candidate completion and transaction
`Seal` require the exact admission/current-Allow chain whenever the admission is
authorization-bearing.

`FinalizationPermit(claim)` is the only way to obtain the owner-bound permit required by an owner
`FinalizeAt`. It verifies the admitted participant and, in authorization-bearing mode, the accepted
current-`Allow` proof; on success it mints a fresh unique `state.FinalizationPermitToken`, records
`{token -> admitted claim, unconsumed}` in the transaction's private permit-issuance registry, and
returns a `MutationFinalizationPermit` carrying the exact context, authority mode and the derived
`state.PermitBinding`. The permit exposes no editor or authorization result. `ApplyFinalized` is the
only application entry point exposed by either mutation transaction and accepts both a
`ResourceMutation` and a `DomainConclusion` finalized change. The transaction first proves permit
provenance at runtime: it reads `change.PermitBinding()`, requires a matching unconsumed entry in
its own permit-issuance registry whose recorded claim equals `change.ParticipantClaim()`, requires
the change's exact publication-context binding, and atomically marks that registry entry consumed so
one permit yields at most one applied change. A missing, foreign, already-consumed, claim-mismatched
or context-mismatched binding is a `MechanicalFailure` and the transaction aborts with no
application. On success it invokes `ApplyTo` against the private detached editor and computes the
actual shadow and index delta digests. For a `ResourceMutation` the transaction alone calls
`operation.NewAppliedChangeLink` from neutral copies of the claim identity, publication binding,
non-empty actual digests and optional result material, then calls
`evidence.BindAppliedChange(link, mutationDescriptor, domainDecisionMaterials)`, receiving one sealed
`AppliedMutationProof`. For a `DomainConclusion` the transaction requires the computed shadow and
index deltas to be empty (proving no domain resource was created), calls
`operation.NewAppliedConclusionLink` from the claim identity and publication binding, then calls
`evidence.BindAppliedConclusion(link, conclusionDescriptor, domainDecisionMaterials)`, receiving one
sealed `AppliedConclusionProof`. No `state` type crosses either evidence API. The transaction records
the proof and the application in its private registry and returns one immutable transaction-created
`AppliedChangeReceipt` exposing only the matching claim, context, publication kind, permit binding
and the applicable proof. No caller can construct a receipt or proof or directly obtain the editor. The receipt and proof attest local application linkage; neither is a domain
resource, FEATURE-0013 carrier, writer, or externally observable result.

**Synchronization and caller inspection.** `state.Store` owns `reservationMu` and `stateMu`. A duplicate
may acquire `reservationMu` while an owner holds `stateMu`; this is what makes the waiter path
reachable. `InspectOrReserveCallerMutation` executes this bounded algorithm:

1. under `reservationMu`, return `Wait` for the same in-flight binding or `Conflict` for a
   different one;
2. if none exists, release `reservationMu`, acquire `stateMu` then `reservationMu`, and recheck;
3. return `Replay`/`Conflict` if the committed or in-flight binding determines it; otherwise insert
   a reservation with an opaque deterministic process-local owner token and return
   `Owner(CallerReservationLease)`;
4. release both locks before returning.

No code holds `reservationMu` while acquiring `stateMu`, and no waiter holds either lock.

**Caller lease and stage 10.** `CallerReservationLease` holds only the owner token. `uow` immediately
installs `defer lease.Abort()`. Before `Begin`, `uow` collects the owner-produced participant claims
and asks `state.NewCallerMutationAdmission(claims, currentnessClaim)` to bind their exact mutation
`ExpectedVersionSet`, candidate dependency versions,
participant cardinality, owner identity, intent digest, and single caller-result role. It adds no
semantic value.
`lease.Begin(admission)` acquires `stateMu`, then briefly `reservationMu` to verify ownership, and
revalidates all admitted mutation predicates and authorization dependency versions while no
finalized domain change, descriptor or publication carrier exists. On failure it removes only its
reservation, signals `Retry`, releases
both locks, invalidates the lease and returns the registered conflict. On success it transfers
reservation ownership to a sealed `CallerMutationTransaction`, releases `reservationMu` while
retaining `stateMu`, creates the exact `AuthorizationAdmissionProof`, invalidates the lease and
returns. `uow` immediately installs `defer tx.Abort()`.

**Controller stage 10.** A controller-owned operation has no client idempotency key, reservation,
waiter or completed replay result. `uow` uses
`state.NewAuthorizedControllerMutationAdmission(claims, currentnessClaim)` only for an
authorization-bearing controller trigger and `state.NewAutomaticControllerMutationAdmission(claims)`
only for expiry, reconciliation or another automatic, non-assignable effect. The constructor split
makes a missing or fabricated currentness claim unrepresentable.
`BeginControllerMutation(admission)` acquires
`stateMu`, checks the admission's exact expected resource versions, participant set, registered
transition and, when authorization-bearing, the complete candidate dependency-version set, and
returns `Conflict` or
the registered effective-once `NoOp` before any carrier preparation. On success it returns a
sealed `ControllerMutationTransaction` with an admission proof only in authorization-bearing
mode; `uow` immediately defers `Abort`. A synthetic caller key or controller-specific replay table
is forbidden.

**Authoritative publication context and owner finalization.** Every successful `Begin` captures one
`evidence.PublicationContext { publicationInstant, publicationSequence }` while holding `stateMu`.
The instant comes only from the injected deterministic `clock.Clock`; Phase 2R has no `UTCClock`,
`time.Now` call or production wall-clock wiring. Only after `Begin` succeeds does `uow` read that
context once. On an authorization-bearing path it first passes the context and exact admission
proof to `authzeval.FinalizeCandidateAt`; `ReevaluateCandidate` aborts and restarts before any
mutation-owner finalization. Success yields the sealed `evidence.FinalizedAuthorizationCandidate`
held unchanged by `uow`. Once the candidate is current—or immediately for automatic
mutation-only mode—`uow` passes the context to each participating mutation semantic owner's
`FinalizeAt`. Each mutation-owner finalizer embeds the exact `permit.Binding()` into the change it
returns; for a `ResourceMutation` it then:

1. revalidates its applicable time-dependent predicates against `publicationInstant`, including
   the strict `publicationInstant < activationDeadline` rule and applicable assignment,
   membership, approval, review or exception validity boundaries;
2. derives its immutable final domain state and canonical after-state digest using that instant;
3. supplies the canonical operation result only when it is the registered originating owner; and
4. derives a non-empty, canonical event/resource/version-ordered
   `EvidenceDescriptors()` set through `evidence.NewMutationDescriptor` for each exact
   registered taxonomy event, resource/version linkage, action and after-state digest.

For a terminal ExceptionProposal `Grant`, the exception owner derives exactly two descriptors in
that canonical set: first `exceptionproposal.decided`, then `exceptiongrant.issued`. Both are bound
through the existing FEATURE-0013 envelope/linkage semantics to the exact immutable
ExceptionProposal, its terminal `Grant` reason, the mandatory `exception-decision/v1`
DecisionRecord material, and the containing operation/correlation evidence; the second descriptor
is additionally bound to the issued immutable ExceptionGrant. The exception owner supplies the same
facts to `evidence.ValidateExceptionTerminalDescriptorSet`. No third descriptor or event is
permitted. `exceptiongrant.proposed` remains submission-only evidence and is not a terminal
descriptor.

An owner that instead reaches a publishable terminal domain conclusion creating no resource — the
exact vehicle for a terminal ExceptionProposal `Deny` under F18-RD-17 — returns a `Finalized` change
whose `PublicationKind()` is `DomainConclusion`. It revalidates its publication-time predicates
against `publicationInstant`, applies an empty state delta (creating no ExceptionGrant or other
domain resource), requests its mandatory `exception-decision/v1` domain-decision material through
`evidence.NewExceptionDecisionMaterial`, invokes `evidence.NewDomainConclusionDescriptor` for its
exact `exceptionproposal.decided` terminal-decision taxonomy event, and supplies the terminal Deny
result when it is the originating owner. The descriptor and decision material are bound through
the existing FEATURE-0013 linkage semantics to the exact immutable ExceptionProposal, its terminal
`Deny` reason, the mandatory `exception-decision/v1` DecisionRecord, and the containing
operation/correlation evidence. The exception owner supplies those sealed values to
`evidence.ValidateExceptionTerminalDescriptorSet`, never invokes `evidence.NewMutationDescriptor`,
creates no ExceptionGrant, and emits neither `exceptiongrant.issued` nor a second decision event.

The mutation-owner finalizer returns a sealed
`operation.FinalizationOutcome[OwnerFinalizedPreparedChange, OwnerDomainOutcome]`. `Finalized`
contains one immutable `FinalizedPreparedChange` of exactly one `PublicationKind`
(`ResourceMutation` or `DomainConclusion`), both of which are published; `Domain` contains only the
owner's existing registered non-publication denial, deadline, expiry, validity, lifecycle or conflict
outcome, which returns its existing mapping and publishes nothing; `MechanicalFailure` identifies only
a local invariant, descriptor-construction or publication-mechanism failure. A publishable terminal
Deny is therefore a `DomainConclusion` `Finalized` outcome, never a `Domain(D)` non-publication
outcome. Neither failure variant contains a new public code. A `Domain(D)` outcome is returned using
its existing mapping and never becomes `INTERNAL_ERROR`; a mechanical failure uses `INTERNAL_ERROR`
only where F18-RD-20/§5.2 already assigns that inherited evidence/publication failure. The finalizer
never rereads time or state and cannot retain the transaction context. For `Finalized`, `uow` passes
the change to the transaction's `ApplyFinalized`; only the transaction invokes `ApplyTo` on its
private detached editor and returns an `AppliedChangeReceipt`. `uow` collects the immutable receipts,
sealed descriptors and the one originating `operation.ResultMaterial`, but cannot construct or alter
them. A `Domain(D)`, `MechanicalFailure` or application failure terminates stage 10: `uow` aborts
before mutation-plan, conclusion-plan, completed-result or carrier construction and publishes nothing. All participants must receive the same context;
placeholders, mutable prepared values, a second clock read and cross-context mixing are forbidden.
The context's sequence becomes committed only with the root; abort does not consume it.

**Completed mutation data flow and stage 11.** After every admitted participant has exactly one
successfully registered `AppliedChangeReceipt`, the finalized state, actual shadow/index delta
digests, descriptors, applied proofs and originating result material are closed. On an
authorization-bearing path, `uow` now invokes
`evidence.CompleteAuthorizationCandidate(finalizedCandidate,
authorizationAdmissionProof, publicationContext)`. That constructor validates exact claim,
proof and context linkage and returns the transaction's sole sealed
`CompletedAuthorizationEvaluation`; this is the exact completed-evaluation point defined by
ADH-2026-072 §2. Automatic mutation-only paths perform neither call and have no completed
evaluation.

`uow` may then collect only the receipts' read-only `AppliedMutationProof` values. An
authorization-bearing path invokes
`evidence.NewAuthorizedMutationPlan(publicationContext, completedEvaluation, appliedProofs)`; an
automatic mutation-only path invokes
`evidence.NewAutomaticMutationPlan(publicationContext, appliedProofs)`. The separate constructors
make a fabricated or omitted authorization evaluation unrepresentable. (A transaction whose applied
participants are publishable terminal domain conclusions instead follows the resource-less
domain-conclusion path defined immediately below.)
`evidence` validates every proof/descriptor/context linkage, sorts them canonically, derives the aggregate
mutation digest, preallocates deterministic carrier
IDs from the context, operation/evaluation digest and stable ordinal, and returns a sealed plan.
For caller mode only,
`idempotency.PrepareCompletedResult` binds the exact originating `operation.ResultMaterial` to the lookup key, request binding,
mutation digest and publication context. `uow` passes the sealed completed result to
`StageCompletedResult`.

`evidence.PrepareMutationCarrierSet` is then the sole constructor and validator of the required
carrier set. `uow` passes the sealed result to `AcceptPreparedEvidence`; it never constructs a
`DecisionRecord`, `AuditEvent`, identifier, digest or carrier literal. The exact carrier set is:

- for an authorization-bearing caller or controller mutation, exactly one protected
  `authorization.evaluated` `AuditEvent` for the supplied completed evaluation; for an
  `AutomaticMutationOnly` controller effect, no authorization event and no authorization
  DecisionRecord;
- one `AuditEvent` for every staged `MutationDescriptor` in the approved atomic effect group;
  duplicate descriptors for the same event/resource/version linkage are rejected, while the same
  taxonomy event applied to different resources remains separate (there may therefore be more than
  one mutation event, and none replaces the authorization event);
- exactly the mandatory `DecisionRecord` values required by F18-RD-19 for that operation, when an
  authorization evaluation exists, plus
  exactly any records selected by the applicable registered `auditRequirement`; and
- no duplicate, unrequired or unlinked carrier.

For a terminal ExceptionProposal `Grant`, this validation additionally requires its one exception
participant's descriptor set to be exactly, in canonical order,
`exceptionproposal.decided` then `exceptiongrant.issued`; it requires both descriptors and the
mandatory `exception-decision/v1` material to bind the same exact immutable proposal, terminal
reason, containing operation/correlation evidence and PublicationContext. The issuance descriptor
must also bind the one immutable issued ExceptionGrant. Missing, extra, reordered, duplicate or
cross-proposal/correlation-linked values reject the sealed plan. This is the ADH-2026-073 Grant
matrix, not a new event or writer.

For a transaction whose applied participants are publishable terminal domain conclusions that create
no resource (the ExceptionProposal `Deny` case), `uow` instead collects the receipts' read-only
`AppliedConclusionProof` values and invokes
`evidence.NewDomainConclusionPlan(publicationContext, completedEvaluation, appliedConclusionProofs)`
followed by `evidence.PrepareDomainConclusionCarrierSet`. That plan validates every conclusion
proof/descriptor/context linkage, forbids any `MutationDescriptor`, `AppliedMutationProof`, resource
delta, ExceptionGrant, or completed caller result that would imply a resource mutation, and produces
exactly: one protected `authorization.evaluated` `AuditEvent` for the completed evaluation; one
`AuditEvent` for each staged `DomainConclusionDescriptor`; the mandatory F18-RD-19 domain
`DecisionRecord` for each terminal conclusion (exactly one `exception-decision/v1` record and one
`exceptionproposal.decided` AuditEvent for an ExceptionProposal `Deny`, each linked through the
existing FEATURE-0013 envelope/linkage semantics to the exact immutable proposal, terminal reason,
and containing operation/correlation evidence); any records selected by the applicable registered `auditRequirement`; and
no duplicate, unrequired or unlinked carrier. Its `PreparedEvidenceChange` is accepted through the
same `AcceptPreparedEvidence` path and committed atomically, so a terminal Deny publishes its
mandatory decision and audit evidence while creating no ExceptionGrant.

Carrier identifiers, mutation linkage and the protected audit ID references are therefore supplied
before `Seal` without creating a FEATURE-0013 allocator or service. `PreparedEvidenceChange`
contains the full carrier set, not a single assumed event.

**Mode-specific seal invariants.** The transaction owns an internal receipt registry, so every
claim below is runtime-checkable rather than inferred from a static call site. Seal is exact and
closed:

- `CallerMutationTransaction` requires `AuthorizationBearing`, the exact transaction-created
  authorization admission proof and linked completed evaluation, the admitted participant count
  to equal the registered receipt count, and exactly one transaction-created receipt per claim with matching owner, intent
  digest and publication context; each receipt's evidence-owned proof bound to the actual
  shadow/index delta digests; one and only one originating
  result material from the admitted originating participant; the exact receipt-linked accepted
  carrier set; and exactly one
  completed result matching its owned reservation, request binding, mutation digest and context;
- `ControllerMutationTransaction` requires exact admitted-participant/receipt cardinality and the
  same owner/context/shadow/index/descriptor/proof linkage, the exact receipt-linked accepted carrier
  set, no originating caller result, and rejects every reservation, waiter and completed idempotency
  value. In `AuthorizationBearing` mode it requires the exact admission proof, completed
  evaluation and authorization evidence. In `AutomaticMutationOnly` mode it rejects every
  currentness claim, admission proof, completed evaluation, authorization event and authorization
  DecisionRecord, while requiring the exact mutation evidence; and
- `EvidenceTransaction` requires no domain/index mutation and no idempotency state, and requires
  the exact admission proof, completed evaluation and accepted authorization carrier set for its
  completed candidate.

When a transaction's applied participants are publishable terminal domain conclusions
(`PublicationKind == DomainConclusion`), Seal additionally requires each such receipt to carry an
`AppliedConclusionProof` bound to an empty shadow/index delta, requires the accepted carrier set to be
the exact `evidence.PrepareDomainConclusionCarrierSet` output (the mandatory F18-RD-19 domain
`DecisionRecord`, the one `exceptionproposal.decided` event linked to the exact immutable proposal,
terminal reason and containing operation/correlation evidence, and — when authorization-bearing — the
`authorization.evaluated` event), and rejects any `MutationDescriptor`, `AppliedMutationProof`,
non-empty delta, ExceptionGrant, or caller completed result. A conclusion and a resource mutation are
never mixed in one transaction, and every applied receipt's `PermitBinding` must match a consumed
entry in the transaction's permit-issuance registry.

For a terminal ExceptionProposal `Grant` represented by a `ResourceMutation`, Seal instead requires
the exact validated two-descriptor Grant set and exactly one mandatory
`exception-decision/v1` DecisionRecord linked to the same proposal, terminal reason and containing
operation/correlation evidence; it requires one immutable ExceptionGrant and rejects a missing,
extra, reordered, duplicate or cross-linked descriptor. For a terminal `Deny`, Seal requires the
DomainConclusion matrix above and rejects every ExceptionGrant, `exceptiongrant.issued` event and
second decision event.

Every mode rejects a missing, extra, duplicate, differently linked or differently timestamped
value. Mutation modes also reject an unadmitted participant, caller-constructed receipt,
pre-context descriptor, unfinalized or mutable change, change finalized with another context,
unmatched or reused permit binding, and
descriptor/digest/result linkage inconsistent with the actual applied shadow and indexes. Failure
aborts without publication. Because the applicable stage-10 outcome and
owner finalization precede carrier construction, evidence failure cannot mask a registered conflict
or extend an expired boundary.

**Commit and abort.** After `Seal`, the non-fallible terminal section installs one shadow root under
`stateMu`. Caller commit installs domain/index changes, carrier set, completed result and committed
publication sequence, then acquires `reservationMu`, removes only its reservation, signals
`Committed`, releases both locks, creates the sealed `CommitReceipt` bound to that accepted carrier
set/context, and invalidates every capability. Controller commit installs only its domain/index
changes, carrier set and sequence before returning the corresponding receipt. A publishable
domain-conclusion transaction installs only its accepted carrier set and committed publication
sequence with no domain/index change and no ExceptionGrant before returning its receipt.
Evidence-only commit installs only its carrier set and sequence before returning its receipt. `Abort` discards the shadow, releases `stateMu`, and—only for caller mode—removes
the owned reservation under the required `stateMu -> reservationMu` order and signals `Retry`.
No abort publishes or consumes a sequence.

**Evidence-only stage 10/11.** After a denied/read/replay `OperationCandidate`, `uow` calls
`BeginEvidenceTransaction(currentnessClaim)`. The transaction checks the complete evaluated
snapshot/dependency-version set; mismatch returns `RetryEvaluation` before completed-evaluation or
carrier preparation. Success captures the authoritative publication context and creates the exact
authorization admission proof under `stateMu`. `uow` passes the candidate, proof and context to
`authzeval.FinalizeCandidateAt`; `ReevaluateCandidate` publishes nothing and restarts the pipeline.
Only after authzeval returns a sealed current candidate does `uow` invoke
`evidence.CompleteAuthorizationCandidate(finalizedCandidate, admissionProof,
publicationContext)`. Only that successful invocation completes the current evaluation. `uow` then asks
`evidence.PrepareAuthorizationCarrierSet(completedEvaluation)` for the exact set: exactly one protected
`authorization.evaluated` `AuditEvent`, plus exactly the authorization `DecisionRecord` values
mandated by F18-RD-19 and the applicable registered `auditRequirement` (possibly none). Acceptance,
seal and commit then precede release of the result or replay. This transaction exposes no editor,
reservation or completed-result method.

**Panic, cancellation and shutdown.** Every lease and transaction is stack-local to one `uow`
operation and has an immediate deferred idempotent abort. `uow.Coordinator` owns only an admission
gate, cancellation context and active-operation wait group—never transaction handles. Shutdown
closes admission, cancels owners/waiters and waits for active operations. Waiters select cancellation
without locks; owners check cancellation before `Begin`, `Seal` and `Commit`, then unwind through
their deferred abort. All work while `stateMu` is held is bounded and local: only the statically
known semantic-owner `FinalizeAt`/`ApplyTo` and sealed evidence/idempotency constructors may run;
arbitrary injected callbacks, state/clock reentry, goroutines and external I/O are forbidden. Thus
shutdown never needs cross-goroutine transaction access. A commit that has already
entered its non-fallible terminal section finishes; shutdown waits for it rather than attempting a
cross-goroutine abort.

**Replay and atomic boundaries.** `Committed` waiters and completed-record replays restart the
ordered pipeline and perform fresh authentication, authorization and safe-visibility checks before
stage-10 replay selection. `Retry` and `RetryEvaluation` restart from the required snapshot stage.
The approved atomic-effect groupings and separation of approval publication from downstream-effect
publication are consumed from F18-RD-20; failure publishes none of a transaction's domain, evidence
or caller idempotency state.

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
   the sealed current candidate; stage 11 prepares, accepts and commits the exact
   authorization carrier set. Only after commit may `uow` invoke
   `authzeval.ReleaseAfterPublication` and return the redacted Deny or protected Allow projection.
3. **Caller mutation.** Stage 10 returns Conflict, Wait, Replay or Owner. Conflict returns before
   completed-evaluation/evidence preparation and discards the candidate. Wait/retry also discards it
   and restarts after the signal. Replay reaches an evidence-only transaction only after the fresh
   replay checks; its current candidate is finalized by `authzeval`, completed and audited before the stored successful result
   is returned. Owner begins the caller transaction without descriptors, obtains its context, asks
   `authzeval` to finalize the admitted candidate, then finalizes and applies every mutation-owner
   intent using that exact context. It completes the candidate only after all applicable stage-10
   checks succeed, derives and stages the exact completed result and carrier set, commits atomically,
   and only then releases the result.
4. **Authorization-bearing controller mutation.** Stage 10 returns Conflict, effective-once NoOp or
   a controller transaction using exact transition/version and authorization-dependency checks. No
   synthetic caller key or completed replay record is created. Conflict/NoOp discard the candidate
   and create no completed evaluation or `AuthorizationResult`. A successful transaction captures
   its context, asks `authzeval` to finalize the candidate, and only then permits the mutation owner
   to finalize/apply its change. The completed evaluation plus mutation plan determines the exact
   carrier set; commit precedes any operation-local result release. When the owner finalizer instead
   reaches a publishable terminal domain conclusion that creates no resource — the ExceptionProposal
   `Deny` case (F18-RD-17/19) — the finalized change is a `DomainConclusion`: `ApplyFinalized`
   consumes the permit and records an empty-delta conclusion receipt, `uow` builds the exact
   `evidence.NewDomainConclusionPlan`/`PrepareDomainConclusionCarrierSet` carrier set (the mandatory
   `exception-decision/v1` DecisionRecord, the terminal-decision `AuditEvent`, and the
   `authorization.evaluated` event), and commit atomically publishes that evidence while creating no
   ExceptionGrant, before any operation-local result release. A non-publication `Domain(D)` outcome
   instead returns its existing mapping and publishes nothing.
5. **Automatic mutation-only controller effect.** Expiry, reconciliation and another closed
   non-assignable automatic effect enters stage 10 without an `OperationCandidate` or currentness
   claim. Conflict/effective-once NoOp returns its existing outcome. Success finalizes and applies
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
downstream authority. `Commit` is the declared non-fallible terminal state transition and is not an
error source. A conflict already selected at stage 10 never enters the evidence-failure path.

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
(DD-05). **`VS0-CF-F18-54` is explicitly excluded from runtime Go tests**: it is the standards
matrix architecture-gate obligation (REQ-F18-23), realized as the §7.3/§8.2 mapping-matrix
check, not as a Go conformance test. **`VS0-CF-F18-50` is likewise an architecture-gate case
(REQ-F18-01), not a runtime Go case**: it is the registered current-feature/prohibited-concept
fixture asserting that a fixture containing a future-owned/prohibited resource, field, adapter,
credential flow, or provider-native IAM object is rejected before evaluation/publication with zero
external I/O (§8.2). The DD-13 supplemental static package check is **separate** design-owned gate
evidence and does not redefine or stand in for `VS0-CF-F18-50`. Runtime cases (the `{01..37,
39..49,51,52,53}` set) are the only ones realized as Go conformance tests; `VS0-CF-F18-50` and
`VS0-CF-F18-54` are architecture/standards gates. `VS0-CF-F18-38` is a permanent tombstone (no
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

### 7.2 Unit and property tests (Go guardrails §30)

The canonical REQ/AC and conformance ledgers define domain test semantics; unit tests do not repeat
them. Design-owned mechanics require these additional focused suites:

- each owner proves private `PreparedIntent` and `FinalizedPreparedChange` construction, registered
  port use and no second writer; prepared intents are immutable and contain no publication time,
  final state/result, after-state digest, descriptor, carrier ID or completed-result value;
  `evidence` proves that only its descriptor/plan/carrier-set constructors create mutation digests,
  carrier IDs, `DecisionRecord`/`AuditEvent` values or `PreparedEvidenceChange`; it proves the sealed
  value has no editor port and imports neither `state` nor `uow`;
- package-DAG/checker fixtures prove every owner may use only the enumerated sealed `state`
  capability types and claim constructors with its own constant, while any owner reference to a
  store, transaction, lease, reservation, lock, root, editor-acquisition or commit/abort method
  fails `F18-ARCH-005`; a `ResourceMutation` omitting `EvidenceDescriptors()` fails
  `F18-ARCH-009`, while a canonical ordered multi-descriptor implementation is a valid fixture;
- ordering tests prove stages 7–9 yield only a non-releasable candidate, stage-10
  Conflict is returned even when stage-11 evidence failure is injected, and no audit publication or
  `AuthorizationResult` occurs before the applicable stage-10 outcome;
- authorization-path tests separately prove denied, standalone-read, replay, caller-mutation and
  authorization-bearing controller-mutation flows; every released result follows successful
  exact-carrier acceptance and commit, while evidence failure returns the safe internal failure with
  no result or authority; automatic mutation-only effects prove exact mutation evidence with no
  candidate, completed evaluation, `authorization.evaluated` event, authorization DecisionRecord or
  released result;
- two-lock race tests hold an owner `CallerMutationTransaction`/`stateMu` open while a duplicate
  obtains a `WaitHandle` under `reservationMu`; prove the `stateMu -> reservationMu` order; prove no
  wait with either lock; and prove install-completed → remove-reservation → signal-`Committed`
  visibility;
- reservation tests cover same/different bindings, owner-lease immediate abort, lease-to-transaction
  transfer, descriptor-free caller/controller `Begin`, CAS failure before owner finalization,
  owner-only cleanup, `Retry`, replay rechecks, cancellation and the inherited retention/JIT
  behavior;
- admission tests prove only registered owner packages construct expected-version sets and claims,
  each claim uses that package's exact owner constant, caller mode admits exactly one originating
  result participant, controller mode admits none, and duplicate participant/intent identities or
  a mode mismatch reject before transaction begin; only `authzeval` constructs the complete
  snapshot/dependency version set and currentness claim; caller/authorization-bearing controller
  modes require it, automatic-mutation-only mode forbids it, and any changed dependency causes
  stage-10 retry/conflict before context capture or evidence preparation;
- two-phase finalization tests prove the transaction context is unavailable before successful
  `Begin`, authorization-bearing paths run `authzeval.FinalizeCandidateAt` first and abort/restart
  on `ReevaluateCandidate` before any mutation finalizer, every participating owner receives that exact context and its owner-bound `FinalizeAt(state.MutationFinalizationPermit)`, the finalized change
  embeds the transaction-issued `PermitBinding` and declares exactly one `PublicationKind`, final state/result/after-state
  digest/descriptor derive only inside owner `FinalizeAt`, no owner rereads clock/state, supporting
  changes expose no canonical result, only the transaction invokes `ApplyTo`, and any finalization
  returns a typed `Finalized | Domain(owner type) | MechanicalFailure`; non-publication `Domain(D)` outcomes retain their
  existing mappings, mechanical failures use only the existing applicable publication-failure
  mapping, and any non-finalized outcome/application failure aborts before completed-evaluation,
  mutation-digest, conclusion-digest, completed-result or carrier construction;
- permit-provenance tests prove `FinalizationPermit` mints one fresh token per issuance and records it unconsumed, `FinalizeAt` embeds exactly that permit's `PermitBinding`, and `ApplyFinalized` accepts a change only when its `PermitBinding` matches a still-unconsumed registry entry whose recorded claim and publication context equal the change's, consumes the entry so a permit yields at most one applied change, and returns `MechanicalFailure` (no application) for a missing, foreign, already-consumed, claim-mismatched or context-mismatched binding;
- terminal-exception evidence tests prove the closed ADH-2026-073 matrix. A terminal
  ExceptionProposal `Grant` is a `ResourceMutation` with a canonical ordered descriptor set of
  exactly `exceptionproposal.decided` then `exceptiongrant.issued`, one immutable ExceptionGrant
  and one mandatory `exception-decision/v1` DecisionRecord. A terminal `Deny` is a
  `DomainConclusion` `Finalized` outcome (never `Domain(D)`), applies an empty shadow/index delta,
  creates no ExceptionGrant, and its conclusion plan atomically accepts exactly one
  `exceptionproposal.decided` AuditEvent and the mandatory `exception-decision/v1` DecisionRecord.
  In both cases, the event(s) and record bind the same exact immutable proposal, terminal reason,
  and containing operation/correlation evidence through existing FEATURE-0013 linkage semantics.
  The tests reject missing, extra, reordered, duplicate, cross-proposal, cross-reason or
  cross-correlation descriptors/material; Deny also rejects every `exceptiongrant.issued` event,
  mutation proof, non-empty delta, ExceptionGrant and caller completed result. A non-publication
  `Domain(D)` outcome publishes nothing;
- receipt tests prove `ApplyFinalized` accepts only an admitted claim bound to the transaction's
  exact context and a matching unconsumed permit binding; only `state` constructs the receipt; only
  `evidence.BindAppliedChange` (resource mutation) or `evidence.BindAppliedConclusion` (domain
  conclusion) constructs its proof; a mutation proof binds actual shadow/index deltas, descriptor and
  optional originating result while a conclusion proof binds an empty delta and its conclusion
  descriptor; and Seal rejects missing, duplicate, foreign, wrong-owner, wrong-intent, wrong-context,
  wrong-permit, wrong-delta, wrong-descriptor and wrong-result receipts/proofs;
- coordinator API tests prove only `uow` calls each exact evidence plan/carrier constructor with
  its permitted plan variant and sealed proof set; only a successful `Commit` produces a
  `state.CommitReceipt`; and `authzeval.ReleaseAfterPublication` rejects a missing, foreign or
  context/evidence-mismatched receipt and is unreachable before commit;
- currentness/runtime-Seal tests prove transaction begin compares the complete dependency set,
  creates the admission proof only after successful comparison, and Seal dynamically rejects a
  wrong/missing claim, receipt, publication context, result, digest, delta, participant/cardinality,
  carrier or authority mode; none of these dynamic equality assertions is attributed to the static
  checker;
- caller-mutation tests require domain change + exact carrier set + a sealed completed result bound
  to the lookup key, request binding, mutation digest and publication context; controller-mutation
  tests require domain change + exact carrier set while forbidding caller keys, reservations,
  waiters and completed records; both reject missing/extra values, prove all-or-nothing commit and
  reject every method after terminal invalidation;
- carrier-cardinality tests cover one-event and grouped multi-event mutations, mandatory and
  policy-selected DecisionRecords, stable deterministic ordering/IDs, and rejection of every
  missing, duplicate, unrequired, unlinked or differently timestamped carrier;
- evidence-only-mode tests forbid editor/domain/reservation/completed state, require exactly one
  `authorization.evaluated` AuditEvent plus the exact required DecisionRecord set, and prove commit
  before result release;
- candidate-boundary tests prove `AuthorizationCandidateDescriptor` is non-authoritative;
  only `authzeval.FinalizeCandidateAt` performs candidate-owned publication-time revalidation;
  `CompleteAuthorizationCandidate` cannot run before applicable stage-10 success or with another
  admission proof/context; Conflict/CAS loss/Wait/retry/NoOp and automatic-mutation-only effects
  never complete an evaluation; and replay,
  standalone and successful mutation paths each complete exactly once before carrier preparation;
- deterministic-time tests prove one publication instant is captured under `stateMu` before owner
  finalization, strict activation succeeds only when `publicationInstant < activationDeadline`,
  equality/expiry fails with no publication, all finalized state/descriptors/results and staged
  completed/evidence values share that instant, mixing contexts is rejected, abort consumes no
  publication sequence, and implemented Phase-2R paths contain no `time.Now`/`UTCClock` wiring;
- controller automatic-operation tests prove exact transition/version CAS, registered Conflict or
  effective-once NoOp, absence of any synthetic client idempotency key or replay table, and the
  exact separation of authorization-bearing versus automatic-mutation-only carrier plans;
- direct-revoke topology tests prove the only path is
  `command.RoleAssignmentRevokeService -> roleassign.DirectRevokePort +
  uow.MutationCoordinatorPort`, with no `roleassign -> uow/command`, direct state/evidence access,
  fourth workflow trigger, route or second writer;
- panic/cancellation/shutdown tests prove stack-local deferred abort, closed admission, waiter
  cancellation and `Coordinator` wait-group completion without storing transaction handles; and
- `make feature-0018-architecture-check` runs every DD-13 valid/invalid fixture, including the
  two-lock lifecycle, all three transaction modes, evidence-constructor ownership (including the
  domain-conclusion descriptor/proof/plan constructors), rule-008
  no-early-publication/post-commit-release split, rule-009 begin/context/permit-bound finalization
  order, rule-010 transaction-only application/receipt/conclusion-link ownership, rule-011 currentness
  call ownership, deterministic-clock ban and exact package DAG. Dynamic Seal, permit-provenance and
  cardinality/equality cases run only in the runtime suites above, not as checker fixtures.

The FEATURE-0013 adoption tests remain those in `## FEATURE-0013 Adoption`. Runtime domain vectors
remain exactly §7.1. Concurrency suites run under `-race`; all tests are deterministic and perform
zero external I/O.

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
complete unless it passes.

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
- REQ-F18-01 scope/exclusion authority (proved by the registered `VS0-CF-F18-50`
  current-feature/prohibited-concept architecture fixture) and REQ-F18-23 standards-validation gate
  (proved by the `VS0-CF-F18-54` standards
  mapping-matrix check) — architecture-gate/boundary obligations, **not runtime Go tests**
  (§7.1/§7.3). The DD-13 `feature-0018-architecture-check` target is separate design-owned gate evidence
  for the DD-01 package topology and does **not** redefine, replace, or stand in for
  `VS0-CF-F18-50`.
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
  design-mechanic ownership and reviewer classification only).
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
| Two-lock, F18-local non-durable publication with caller/controller/evidence transactions and authorization-bearing/automatic-mutation-only authority modes; stage-10 begin checks complete candidate dependencies and captures one context; transaction-only `ApplyFinalized` creates exact context-bound receipts/proofs before completed-evaluation and carrier preparation (DD-02/DD-10, §5.4) | ADH-2026-072 §§2–3; REQ-F18-20/21/22 | owner packages, `authzeval/`, `state/`, `uow/`, `evidence/`, `idempotency/`, `clock/` | VS0-CF-F18-17 (`CONFLICT` precedence/replay); VS0-CF-F18-20 (accepted evidence before all-or-nothing publication); VS0-CF-F18-53 (deterministic zero-I/O); focused currentness/receipt/Seal tests |
| Single permit-gated `FinalizeAt(state.MutationFinalizationPermit)` with unforgeable `PermitBinding` minted per issuance, transaction-private permit-issuance registry, and runtime `ApplyFinalized` permit-provenance check consuming one entry per permit (DD-01/DD-13, §5.4) | REQ-F18-20/21; ADH-2026-072 §3 | `state/`, owner packages, `uow/` | focused permit-provenance/finalization tests; `make feature-0018-architecture-check` (rules 009/010) |
| Publishable terminal `DomainConclusion` evidence path for a resource-less terminal decision — ExceptionProposal `Deny` — that creates no ExceptionGrant yet atomically publishes the mandatory `exception-decision/v1` DecisionRecord and required AuditEvents after stage 10 and before result release (DD-01, §5.4/§5.5) | F18-RD-17/19/20 | `exception/`, `evidence/`, `state/`, `uow/` | VS0-CF-F18-22 (non-exceptionable terminal Deny), VS0-CF-F18-43 (approved terminal Deny); focused domain-conclusion publication tests |
| ADH-2026-073 terminal-ExceptionProposal evidence matrix: canonical ordered Grant descriptor set (`exceptionproposal.decided`, then `exceptiongrant.issued`) and one-event Deny conclusion, each linked to the exact immutable proposal, terminal reason, mandatory `exception-decision/v1` record and containing operation/correlation evidence through FEATURE-0013 | REQ-F18-19/20; ADH-2026-073 | `exception/`, `evidence/`, `state/`, `uow/` | focused Grant/Deny descriptor-set, linkage and Seal rejection tests |
| Caller idempotency lookup key separated from exact target/scope/digest binding; independently observable reservation table and direct waiter handle; controller-owned operations use only exact transition/version CAS; replay rechecks after wake (DD-11, §5.4) | REQ-F18-13/21 | `idempotency/`, `state/`, `uow/` | VS0-CF-F18-17 plus payload/target/scope conflict and controller effective-once vectors |
| Stages 7–9 produce a non-authoritative candidate plus complete dependency-version claim; state proves currentness, `authzeval` owns publication-time revalidation, and only applicable success completes an audited Allow/Deny evaluation; automatic mutation-only effects never fabricate one (DD-12, §5.1/§5.5) | ADH-2026-072 §2; REQ-F18-19/20/21 | `authzservice.go`, `authzeval/`, `evidence/`, `uow/`, `state/` | VS0-CF-F18-20 (audit-before-result); VS0-CF-F18-17 (conflict precedence and replay recheck); focused automatic-effect tests |
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
| Standards mapping matrix as architecture gate, not runtime test (§7.1/§7.3) | REQ-F18-23 | mapping-matrix check | VS0-CF-F18-54 (non-runtime standards gate) |

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
mandatory F18-RD-17/19/20 decision and audit evidence while creating no ExceptionGrant. FEATURE-0013
remains contract/validation-only (ADH-2026-072 §1); the candidate-to-completed-evaluation point is
closed by ADH-2026-072 §2. No stop condition is active.

**Review classification.** Every subsequent FEATURE-0018 design finding must be exactly one of
`STALE_TRANSCRIPTION`, `DESIGN_EXECUTABILITY`, `REQUIREMENT_GAP`,
`ARCHITECTURE_CLARIFICATION_REQUIRED`, or `OUT_OF_SCOPE`, using the definitions in
ADH-2026-072 §5. A reviewer must cite the exact contradictory/missing authority and target text
before reopening a closed decision. An incomplete API, lock, receipt, checker or test mechanism is
`DESIGN_EXECUTABILITY`, not an architecture gap.

### 10.4 Final self-verification

- Approved REQ/AC semantics, ADH-2026-070/071 mappings, routes, writers, exclusions and the
  `VS0-CF-F18-01..54` identifier ledger are unchanged.
- FEATURE-0013 supplies decision/audit carrier contracts, strict local bundle loading and bounded
  validation only. FEATURE-0018 claims no FEATURE-0013 transaction, acceptance service, allocator,
  durable store, outbox, queue or persistence behavior.
- Stages 7–9 produce only a sealed internal candidate. No authorization audit or
  `AuthorizationResult` is published before the applicable stage-10 coherence/conflict/replay/CAS
  outcome.
- The exact F18-RD-20 acceptance point is successful `AcceptPreparedEvidence` on the applicable
  sealed transaction. `evidence` alone constructs and validates candidate/completed-evaluation
  descriptors, admission/finalized-candidate/applied proofs, mutation descriptors/digests, carrier
  IDs, exact carrier sets and `PreparedEvidenceChange`; `uow` only coordinates sealed values.
- Stage 10 completes before any stage-11 carrier construction, validation or acceptance failure can
  surface. A registered conflict therefore cannot be masked by evidence failure.
- Caller/controller `Begin` accepts a sealed mode-specific `MutationAdmission`; authorization-bearing
  admission includes the candidate's complete dependency-version claim and automatic-mutation-only
  admission cannot include one. State revalidates versions, captures the context and binds the
  admission proof before any final state, result, after-state digest or descriptor exists.
  `authzeval` finalizes only candidate-owned time semantics; each mutation owner finalizes only its
  own sealed intent. Transaction-only `ApplyFinalized` applies changes and creates the sole
  applied-receipt/proof linkage.
- `reservationMu` makes in-flight waiters observable while an owner holds `stateMu`; lock order is
  `stateMu -> reservationMu`, and waiting holds neither. Owner leases and transactions immediately
  defer idempotent abort and remain stack-local.
- Caller/controller mutation Seal verifies exact admitted-participant/transaction-created-receipt
  cardinality and owner/intent/context/shadow/index/descriptor/proof linkage. Caller mode also
  requires its owned reservation and one matching completed result; controller mode forbids caller
  idempotency. Authorization-bearing modes require the currentness/admission/completed-evaluation
  chain; automatic mutation-only mode rejects that chain and requires mutation carriers only.
  Evidence-only Seal forbids domain/idempotency state and requires the exact completed-evaluation
  carrier set.
- Every `FinalizeAt` uses the single `FinalizeAt(state.MutationFinalizationPermit)` contract across
  DD-01, DD-13, the component descriptions, the checker rules, the exact API and the tests. The
  finalized change embeds an unforgeable `state.PermitBinding`, and `ApplyFinalized` proves provenance
  at runtime by consuming a matching unconsumed entry in the transaction's private permit-issuance
  registry (equal claim and context), rejecting a missing, foreign, consumed, claim- or
  context-mismatched binding.
- A publishable terminal domain conclusion that creates no resource — the ExceptionProposal `Deny`
  case — is a `DomainConclusion` `Finalized` outcome, applies an empty delta, creates no
  ExceptionGrant, and atomically publishes exactly one mandatory `exception-decision/v1`
  DecisionRecord plus the terminal-decision and `authorization.evaluated` AuditEvents via
  `evidence.PrepareDomainConclusionCarrierSet` after stage 10 and before result release. A
  non-publication `Domain(D)` outcome publishes nothing and returns its existing mapping.
- Every successful transaction start captures one authoritative publication instant under
  `stateMu` before finalization; `authzeval` revalidates candidate time predicates and each mutation
  owner revalidates its own domain predicates against it. The same context is used in finalized
  state/descriptors/results, completed results and evidence. Phase 2R implements no `UTCClock`,
  `time.Now` or production wall-clock wiring.
- Shutdown cancels through `uow.Coordinator` and waits for owning operations to unwind; it stores no
  transaction handle and performs no cross-goroutine abort.
- Only `uow` invokes evidence preparation, transaction acceptance and
  `authzeval.ReleaseAfterPublication`, and release occurs only after successful commit. Rules
  001/003/005/008/009/010/011 cover only statically checkable lifecycle, constructor, import,
  deterministic-clock, call-site and receipt/application ownership patterns. Runtime tests—not the
  checker—prove version/context/digest/delta/cardinality equality, Seal and all-path ordering.
- `VS0-CF-F18-12/47/48/49` require `success == false` and outcome class `SafeError` while
  remaining agnostic about Problem code, violation channel, HTTP status and precedence. Runtime Go
  cases are exactly `01..37,39..49,51..53`; case 50 is architecture-only, 54 standards-only, and
  38 remains the tombstone.
- The exact `## FEATURE-0013 Adoption` section assigns profile definitions to
  `evidence/profiles.go`, strict `decision/bundle.Load` composition to the root, and tests for
  exactly three profiles, duplicate/invalid rejection and absence of a fourth profile.
- Section 9 retains exactly 24 compact REQ disposition rows and 49 exact AC mapping rows; inherited
  schema/writer/state/error ownership is cited to ADH-2026-071 and requirements rather than copied.
- No live route, external I/O, real IdP, provider-native IAM mutation, durable persistence,
  future-feature mechanism, new resource, action, mapping, writer or public error code is
  introduced. No tasks are generated.
---

## Model Execution Report

```text
Model Execution Report:
- Tool: Kiro initial design; Codex applied human-approved ADH-2026-072 reconciliation; Kiro applied
  reviewer-directed permit-provenance and domain-conclusion-publication reconciliation
- Stage or task: FEATURE-0018 design clarification reconciliation
- Scope: design.md; automatic-effect/evaluation separation, typed finalization outcomes,
  candidate-currentness admission, semantic-owner publication-time validation, neutral results,
  acyclic direct revoke, static/runtime proof separation and inherited-semantics deduplication;
  plus (this revision) a single permit-gated `FinalizeAt(state.MutationFinalizationPermit)` contract
  with an enforceable transaction-owned permit-issuance registry and unforgeable `PermitBinding`, and
  an exact publishable terminal `DomainConclusion` evidence path (ExceptionProposal `Deny`) that
  creates no ExceptionGrant yet atomically publishes the mandatory `exception-decision/v1` record and
  required AuditEvents
- Architecture changed: no
- Requirements or mappings changed: no
- Routes, writers, ownership, phase boundary, context budget changed: no
- Tasks generated: no
- Dependency expansion: no; FEATURE-0013 remains carrier/bundle/validation-only
- Reconciled mechanics: descriptor-free stage-10 begin; single-context capture; owner-sealed
  publication-time-free intent followed by exact-context immutable finalization; transaction-only
  application; state-owned receipts with evidence-owned applied proofs; explicit candidate-to-
  completed-evaluation transition; unchanged atomic commit/post-commit release ordering;
  single permit-gated `FinalizeAt(state.MutationFinalizationPermit)` with transaction-owned
  permit-issuance registry and runtime `ApplyFinalized` permit-provenance check; publishable terminal
  `DomainConclusion` evidence path (ExceptionProposal `Deny`) with no ExceptionGrant and atomic
  mandatory decision/audit acceptance
- Receipt: retained exactly once
```
STAGE_STATUS: COMPLETE
