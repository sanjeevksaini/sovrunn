# Design Document

## Overview

See Introduction and §1 below for stage identity, authorities, and closed
boundary.

---

## Introduction

This document is the FEATURE-0016 Design stage for Adapter Boundary and
ExecutionTarget Qualification. It translates the approved FEATURE-0016
requirements (`.kiro/specs/adapter-boundary-and-executiontarget-qualification/requirements.md`),
the architecture boundary
(`docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`),
and the manifest-controlled FEATURE-0016 handoff set — ADH-2026-058 as
foundational contract closure together with ADH-2026-060, ADH-2026-061,
ADH-2026-063, ADH-2026-064, ADH-2026-065, and ADH-2026-066 — into an implementable
representation. It introduces no product semantics, transfers no ownership, and
designs no excluded adjacent-feature mechanism. It chooses only the mechanics
explicitly delegated to design by requirements §9 and the architecture boundary
§7: in-memory registry data structures, handler wiring behind the fixed
pre-ServeMux guard, idempotency reservation/waiter internal data structures,
injected-clock scheduling internals, package/file layout, locks, and internal
helper composition.

Every observable behavior below is fixed by the authorities above. Where this
document names a struct, package, or field name, that name is a design-local
mechanic; it never adds, renames, or removes an observable field, route, state,
error code, projection member, audit effect, or violation code.

---

## 1. Identity, stage, inputs, and closed boundary

### 1.1 Identity

| Field | Value |
|-------|-------|
| Feature | FEATURE-0016 — Adapter Boundary and ExecutionTarget Qualification |
| Stage | Design |
| Kiro slug | `adapter-boundary-and-executiontarget-qualification` |
| Baseline | ARCH-2026.08-PHASE2R-CANONICAL |
| Phase | 2R |
| Order | 6 |
| Owned resources | `ExecutionTarget` (VS0-SCHEMA-015), `NormalizedTargetFactSet` (VS0-SCHEMA-016), `TargetQualificationResult` (VS0-SCHEMA-017) |
| Owned state machine | `VS0-STATE-004` |
| Owned writer | `VS0-WRITER-006` (`ExecutionTargetLifecycleService`) |
| Controlling decisions | DEC-0036, DEC-0042, DEC-0057 |
| Controlling handoffs | ADH-2026-025, ADH-2026-040, ADH-2026-042, ADH-2026-045, ADH-2026-058, ADH-2026-060, ADH-2026-061, ADH-2026-063, ADH-2026-064, ADH-2026-065, ADH-2026-066 |
| Local conformance | VS0-CF-F16-01..128 |

### 1.2 Stage inputs (consumed by reference, not redefined)

- Approved requirements (this spec's `requirements.md`): the sole source of REQ
  and AC meanings and the 128-case conformance ledger (VS0-CF-F16-01..128).
- `docs/architecture/FEATURE-0016-adapter-boundary-and-executiontarget-qualification.md`
  together with ADH-2026-058 (foundational contract closure) and the completing
  handoffs ADH-2026-060 (Maintenance-race outcome), ADH-2026-061
  (Maintenance-entry link-clearing), ADH-2026-063 (collection-create
  phase-one/strict-classification precedence), ADH-2026-064 (VS-000
  ExecutionTarget ownership correction), ADH-2026-065 (create phase-one
  reference extraction and exact classification cases), and ADH-2026-066
  (private coherent backing-access compatibility bridge): the current semantic
  authority for VS0-SCHEMA-015..017 and VS0-STATE-004. ADH-2026-058 is
  foundational but not the sole current semantic authority.
- FEATURE-0012 (`internal/apimeta`, `internal/apiproblem`, `internal/apivalid`,
  `internal/apiref`, `internal/apicond`, `internal/apischema`): ObjectMeta,
  TypedRef, Problem Details, Condition, media semantics, optimistic-concurrency
  conventions — reused, not forked.
- FEATURE-0013 (`internal/decision`): AuditEvent atomic append-before-publication
  semantics — reused, not redefined.
- FEATURE-0015 (`internal/cloudmodel`): `CloudProviderParticipation` and
  `InfrastructureStack` read-only prior-feature authority, plus the
  server-resolved authorization, safe-denial, and idempotency conventions to be
  reused — never modified.

### 1.3 Closed boundary for this stage

This design modifies only
`.kiro/specs/adapter-boundary-and-executiontarget-qualification/design.md`. It
generates no source code, schema, task, automation, architecture, prompt, or
validator. It decides only mechanics delegated by requirements §9 and the
architecture boundary §7. It emits exactly one approved stop condition and halts
if a semantic choice were found missing; §10.3 records that none was
encountered.

Fail-closed stop conditions available to this stage:
`ARCHITECTURE_DECISION_REQUIRED`, `REQUIREMENT_CLARIFICATION_REQUIRED`,
`BOUNDARY_CHANGE_REQUIRED`, `DEPENDENCY_APPROVAL_REQUIRED`,
`SECURITY_REVIEW_REQUIRED`.

---

## Architecture

Architecture is realized by DD-01 through DD-17 (§2) and component layout (§3).

---

## 2. Resolved design decisions

Each decision below is a mechanic that realizes an already-fixed observable
behavior. No decision resolves a product semantic; all product semantics are
supplied by the authorities in §1.2.

### DD-01 — Domain service package placement

FEATURE-0015 realized its canonical cloud-model domain service in a single
cohesive package `internal/cloudmodel/` (`store.go`, `lifecycle.go`,
`idempotency.go`, `scheduler.go`, `grants.go`, `audit.go`, `model/`,
`validate/`). Following that verified live convention, FEATURE-0016 places its
domain service in a sibling package `internal/executiontarget/`
(`store.go`, `lifecycle.go`, `qualify.go`, `observer.go`, `idempotency.go`,
`scheduler.go`, `clock.go`, `projection.go`, `audit.go`, `backing_access.go`, `violations.go`,
`model/`), as verified live in the repository. This keeps the ExecutionTarget
lifecycle, observer, qualification, expiry/maintenance scheduling, idempotency
reservation table, projection, and audit assembly in one small, testable
package. ADH-2026-066 adds only a private F0016 compatibility extension in
`internal/cloudmodel/store.go` and `store_test.go`; FEATURE-0015's product
authority, routes, schemas, lifecycle, and public behavior remain read-only.
Rationale: minimal dependencies, explicit
structs, and clear ownership boundaries per the engineering steering; it also
physically separates the `execution.sovrunn.io` group from the FEATURE-0015
`governance.sovrunn.io` group.

### DD-02 — Sole committer is a single in-process service type

`ExecutionTargetLifecycleService` (VS0-WRITER-006) is the only type that mutates
persisted ExecutionTarget state and the internal `NormalizedTargetFactSet` /
`TargetQualificationResult` records, owns the target-bound in-process
current-Maintenance marker, and computes the response-only `effectiveAvailability`
projection. All create/qualify/retire/maintenance/expiry mutations funnel through
this type under one service-owned `sync.Mutex` guarding all F0016 targets,
indexes, internal records, Maintenance markers, and reservations (including
their `abortCause` and `done` channel). There are no nested F0016 locks or
per-target locks. Observation and profile evaluation run with this mutex
unlocked. Initial backing access runs before this mutex. Final qualification
commit takes the ADH-2026-066 backing-read lease first and then this mutex,
retaining that order through audit append and publication; no code holding this
mutex can acquire the backing lease. Clients and HTTP handlers never write
status.

### DD-03 — Synthetic observer is a target-bound fixture reader

`sovrunn.synthetic-iaas-observer/v1` is implemented as an external-effect-free,
clock-driven component that reads a target-bound in-memory fixture and returns
either four named facts (each Supported/Unsupported/Unknown) or a classified
deterministic outcome (missing fixture; fixture-declared logical timeout). It
performs no network I/O and never blocks on a real call; the "logical timeout"
is a fixture-declared outcome resolved against the injected clock, not an elapsed
wait. It proposes facts only; it never commits.

### DD-04 — Qualification is synchronous with a captured-fence reservation

Qualify is executed synchronously on the request goroutine. The fixed pipeline
ordering for qualify is (per §5.1 and reviewer correction):

1. Authentication, header-form/grammar validation, UID-only safe target
   resolution, derived-scope authorization.
2. **Replay/admission** (step 9): idempotency replay or new InFlight owner
   reservation establishment only; it does not acquire target Qualifying state.
3. **Current If-Match** (step 10): stale ETag check against persisted
   resourceVersion; a post-reservation rejection detaches the owner reservation
   through the lifecycle service before its waiter signal is closed.
4. **Fresh backing validation and initial fence capture** (step 8a): a fresh
   `BackingAccessProvider` result (DD-17) is obtained using the target's persisted
   participation and stack UIDs. This validates current backing viability
   (ineffective participation or non-Active stack → abort with registered
   outcome) and supplies the initial captured fences (InfrastructureStack
   generation and viability fingerprint).
5. **Lifecycle state evaluation and Qualifying acquisition** (step 11): under
   the lifecycle mutex, reject Retired, active-Maintenance, or different-key
   Qualifying; otherwise acquire the target-local Qualifying reservation and
   capture maintenance epoch.
6. **Observer invocation and profile evaluation** — without the lifecycle-service
   mutex.
7. **Pre-commit recheck and fenced commit** — obtain a fresh
   `BackingAccessProvider` lease, compare the registered backing fences, then
   acquire the lifecycle mutex while retaining the lease; audit append and
   publication occur only while both are held in the ADH-2026-066 lock order.

The `Qualifying` in-flight reservation (never persisted, never projected)
records:

- The captured maintenance epoch (from target state when Qualifying is acquired).
- The initial backing fences (InfrastructureStack generation and viability
  fingerprint) captured at step 4. Participation generation is not a fence.
- An `abortCause` field (initially nil) that, once set under the mutex, is
  immutable for the lifetime of the reservation.

The linked idempotency reservation is established before Qualifying state and
is linked to it only once step 5 succeeds; it carries:

- The `done` channel (unbuffered broadcast signal), allocated at InFlight
  establishment. Closed exactly once under the lifecycle-service mutex.
- The `abortOutcome` field, set by `releaseReservation` from the effective
  abort reason (which respects the qualification reservation's first-winner
  `abortCause`). Read by the owner after mutex release.

At commit, the ordered predicate sequence is:

1. **Active-Maintenance marker** (checked first per ADH-2026-060): if an active
   current-Maintenance marker exists, abort qualification with `CONFLICT`/409
   + `VS0_TARGET_MAINTENANCE` (VS0-CF-F16-89).
2. **Retired predicate**: if the target is now Retired, abort qualification
   with `CONFLICT`/409 + `VS0_TARGET_RETIRED` (VS0-CF-F16-88).
3. **Maintenance epoch fence**: if the captured maintenance epoch differs from
   the current maintenance epoch (and no active marker exists), abort with
   `STALE_RESOURCE_VERSION`/412 + `VS0_TARGET_EPOCH_STALE` (VS0-CF-F16-75).
4. **InfrastructureStack generation fence**: if the recheck stack generation
   differs from the captured value, abort with `STALE_RESOURCE_VERSION`/412
   + `VS0_TARGET_EPOCH_STALE` (VS0-CF-F16-29).
5. **Viability fingerprint fence**: if the recheck viability fingerprint
   differs from the captured value, abort with `STALE_RESOURCE_VERSION`/412
   + `VS0_EXECUTION_TARGET_VIABILITY_STALE` (VS0-CF-F16-30).

These outcomes are mutually exclusive: at most one applies to a given
commit attempt. The active-Maintenance marker is always checked first
(ADH-2026-060 "an active current-Maintenance marker always selects
CONFLICT/409 + VS0_TARGET_MAINTENANCE first"); the epoch-stale case never
applies when an active current-Maintenance marker is present. In every one of
these aborts, no qualification conclusion, no qualification AuditEvent, and no
idempotency completion is published; if a Maintenance transition independently
won the race, that transition's own AuditEvent remains audited regardless of
the qualification abort.

Only after this ordered decision and the required AuditEvent append succeeds
does the service atomically replace both current record refs, refresh
`observedGeneration`, advance resourceVersion/ETag, mark idempotency complete,
and publish. It then releases the mutex and wakes same-key waiters via the
`done` channel.

#### Abort-cause retention (first-winner rule)

When a Retire or Maintenance transition wins the race and aborts the in-flight
qualification (DD-04a, DD-12), the `releaseReservation` primitive (DD-05) sets
the qualification reservation's `abortCause` field exactly once under the
lifecycle-service mutex. The `abortCause` is immutable once set: a later
transition cannot overwrite a previously recorded cause. This ensures the
qualifying owner always receives the response corresponding to the first
winning transition. Concretely:

- If Maintenance entry sets `abortCause = VS0_TARGET_MAINTENANCE` first, a
  subsequent Retire cannot change it to `VS0_TARGET_RETIRED`.
- If Retire sets `abortCause = VS0_TARGET_RETIRED` first, a subsequent
  Maintenance entry cannot change it to `VS0_TARGET_MAINTENANCE`.
- The cleanup that fixes the first abort cause atomically detaches the channel;
  its caller closes that captured channel after unlocking, so the owner observes
  a complete state transition rather than a partial in-lock wake-up.

The qualifying owner, upon reacquiring the mutex for its own commit (or upon
observing its `done` channel closed before mutex acquisition), reads
`abortOutcome` from its held idempotency reservation pointer and returns the
corresponding response deterministically.

#### Done-channel close ownership

DD-05 is the sole close protocol. Every completion or abort detaches the
channel and invalidates its field while the lifecycle mutex is held, then
unlocks and closes the captured channel once. This covers success,
Retire-wins, Maintenance-entry, stale fences, audit failure, observer fault,
panic, cancellation, and shutdown. No close occurs under the mutex and no
handler has an independent cleanup path.

### DD-04a — Retire-wins commit path (qualification abort on retirement)

When a retire request passes lifecycle-state evaluation (step 11) and commits
while a Qualifying reservation is in flight for the same target, the retirement
wins. The exact commit path is:

1. The retire handler acquires the lifecycle-service mutex at commit time.
2. The Retired-at-commit predicate is checked: if the target is already Retired
   (another retire committed between reservation and this commit), abort with
   `CONFLICT`/409 + `VS0_TARGET_RETIRED`.
3. The retire commits in audit-before-publication order:
   a. Construct the retirement AuditEvent.
   b. Call `AuditAppender.Append`. If the append **fails**, the retire aborts:
      no state mutation, no ETag advance, no idempotency completion, and
      `INTERNAL_ERROR`/500 is returned. The qualification reservation is
      unaffected (it remains in flight; the retire simply failed).
   c. Only after successful append: atomically persist `Retired/Unqualified`,
      clear both current FactSet/Result links, release the backing tuple,
      advance resourceVersion/ETag, and mark the retire idempotency entry
      complete.
4. As a post-commit side-effect inside the same mutex hold, the service
   locates the in-flight Qualifying reservation for this target (at most one
   exists). If found, it calls
   `releaseReservation(qualifyIdempotencyKey, VS0_TARGET_RETIRED)` (DD-05)
   which applies the first-winner rule: if `abortCause` is already set (a
   prior Maintenance entry won), it is not overwritten; if nil, it is set to
   `VS0_TARGET_RETIRED`. The cleanup primitive removes both the qualification
   and its linked idempotency reservations, detaches the `done` channel under
   the mutex, then closes it after unlock to wake waiters.
5. The mutex is released; the retirement response (200, Retired/Unqualified,
   effectiveAvailability Unavailable, new ETag) is returned.
6. The qualifying owner, upon reacquiring the mutex for its own commit (or
   upon observing its `done` channel closed), reads `abortOutcome` from its
   held idempotency reservation pointer and deterministically returns the
   corresponding response: `CONFLICT`/409 + `VS0_TARGET_RETIRED` (if retire
   was the first winner) or `CONFLICT`/409 + `VS0_TARGET_MAINTENANCE` (if
   Maintenance was the first winner), with no qualification result,
   completion, or qualification AuditEvent
   (VS0-CF-F16-88 / VS0-CF-F16-89 respectively).

This path ensures no stale qualification conclusion can publish after
retirement. The retire's own idempotency completion is retained and replayable.
The winning retirement is always audited; the aborted qualification produces
no AuditEvent or idempotency completion.

### DD-05 — Idempotency reservation table: unified signaling, cleanup, and waiter model

FEATURE-0016 reuses the FEATURE-0012/0015 idempotency conventions (SHA-256 of
recursively sorted allowed create JSON, or SHA-256 of zero bytes for actions;
`If-Match`, `Accept`, authentication, and request/correlation IDs excluded) but
owns its own in-process reservation table keyed by the exact namespace
`(principal UID, route pattern, derived CloudProvider UID, concrete
action-target UID when applicable, Idempotency-Key)`. The table holds `InFlight`
and `Completed` entries. This is F0016-owned state cleared on restart; it does
not touch FEATURE-0015's table. GET and LIST never require, inspect, reserve,
wait on, read, or replay a record.

#### Unified reservation object with broadcast signal

Every InFlight idempotency reservation — whether for create, qualify, or
retire — carries a **broadcast completion signal** (`done` channel). This
ensures that same-key waiters for any action type have a defined wait/wake
mechanism. The reservation object contains:

- `state`: `InFlight` or `Completed`.
- `done`: an unbuffered channel (`make(chan struct{})`) allocated at the
  instant the reservation transitions to `InFlight`. The channel is created
  atomically with the InFlight state establishment under the lifecycle-service
  mutex; no waiter can observe an InFlight reservation without its channel.
- `abortOutcome`: initially nil; once set under the mutex, immutable. Records
  the owner-only abort reason so the owner can return the correct response
  after map removal.
- `storedResponse`: set only on successful completion; contains the replayable
  response (status, body, ETag, Content-Type).
- `digest`: the SHA-256 body digest (for replay-mismatch detection).
- `qualificationRef`: for qualify only, a reference to the linked per-target
  qualification reservation (nil for create/retire).

For qualify, an additional **qualification reservation** exists per target
(at most one). It contains `abortCause`, the captured fences, and the target
UID. The qualification reservation does **not** own a separate channel; the
idempotency reservation's `done` channel is the sole broadcast signal for both
the qualify owner and its waiters. The idempotency reservation is created first
with `qualificationRef=nil`. Only after If-Match, backing admission, and
lifecycle checks succeed does one lifecycle-mutex hold insert Qualifying state
and link `qualificationRef`. A waiter can wait on the InFlight channel before
that point, but cannot infer that Qualifying state exists.

#### Channel ownership and close rules (all action types)

Every InFlight idempotency reservation owns one `done` channel. The lifecycle
service performs the only state transition: under its mutex it records a
Completed or abort outcome, removes/updates map entries, detaches the channel,
and sets its field to nil. It then unlocks and the one caller holding that
captured channel closes it. Completion, owner abort, Retire-wins,
Maintenance-entry, and shutdown all use this same two-phase protocol. A nil
field cannot be detached twice, so exactly one close is possible; wake-up is
therefore always after unlock and never while another waiter could observe
partial reservation state.

#### Unified cleanup primitive (`releaseReservation`)

One cleanup primitive is used on **every** unsuccessful path after reservation
establishment. It is always called under the lifecycle-service mutex and
returns the one detached channel for a post-unlock close. Its behavior is:

```text
releaseReservation(idempotencyKey, abortReason) -> capturedDone:
  1. Look up the InFlight idempotency reservation by key.
     If not found or already Completed → no-op (idempotent).
  2. If the reservation has a qualificationRef:
     a. Check the qualification reservation's abortCause:
        - If already set (first-winner rule): use the existing abortCause as
          the effective abort reason (ignore the incoming abortReason).
        - If nil: set abortCause = abortReason.
     b. Remove the qualification reservation from the per-target map.
  3. Record abortOutcome = effective abort reason on the idempotency
     reservation (owner-only; used if the owner re-enters the mutex).
  4. Remove the idempotency reservation from the namespace table.
  5. Detach any non-nil done channel, set its field to nil, and return it.
```

This primitive guarantees:
- All linked reservations are atomically removed in one mutex hold.
- The `abortOutcome` is retained for the owner's return path even after map
  removal: the owner holds a direct pointer to the reservation object (Go GC
  keeps it alive) and reads `abortOutcome` from that pointer after the mutex
  is released.
- The caller closes the returned channel exactly once after unlocking, waking all waiters.
- Subsequent calls for the same key are no-ops.

Every unsuccessful path after reservation calls this primitive:
- Stale current If-Match (step 10 in §5.1).
- Unavailable backing (qualify step 8a: ineffective participation or non-Active
  stack).
- Retired/Maintenance lifecycle state (step 11).
- Different-key qualification conflict (step 11).
- Fenced-commit abort (steps 3-7 in DD-13.5).
- Observer fault (step 8 in DD-13.5).
- Audit append failure (step 11 in DD-13.5).
- Panic recovery (via defer).
- Shutdown abort (DD-06 step 3).

For create and retire, the same primitive is used but `qualificationRef` is
nil, so step 2 is skipped and only the idempotency reservation is removed.

#### Waiter model (applies to create, qualify, and retire)

A **waiter** is a subsequent request with the same idempotency key that arrives
while the owner's reservation is InFlight. Waiters do **not** receive the
owner's unpublished outcome. Their behavior:

1. The waiter finds the InFlight idempotency reservation under the mutex. It
   reads the `done` channel. Because the channel is allocated atomically with
   the InFlight state, the channel is always non-nil when the reservation is
   InFlight and visible to the waiter.
2. The waiter releases the mutex and blocks on
   `select { case <-done: case <-ctx.Done(): }`.
3. **If `ctx.Done()` fires first:** the waiter detaches and returns a
   context-cancellation error to its caller. The reservation is unaffected;
   other waiters remain blocked.
4. **If `done` fires first (channel closed):** the waiter wakes up. It does
   **not** read the owner's result directly. Instead, it re-enters the
   pipeline at authentication (§5.1 step 2) and repeats the full sequence:
   authorization, safe access, validation, and reservation lookup. At
   reservation lookup:
   - If the idempotency reservation is now `Completed` (owner succeeded):
     the waiter receives the stored replay response.
   - If the idempotency reservation was removed (owner aborted): the waiter
     may establish a new InFlight reservation and become the new owner.

#### External abort paths (Retire-wins / Maintenance-entry)

When Retire-wins (DD-04a) or Maintenance-entry (DD-12) aborts a qualification:

1. Under the lifecycle-service mutex, the aborter calls
   `releaseReservation(qualifyIdempotencyKey, abortReason)` where
   `abortReason` is `VS0_TARGET_RETIRED` or `VS0_TARGET_MAINTENANCE`.
2. The primitive applies the first-winner rule, removes both reservations,
   and detaches the channel for exactly-once close after unlock.
3. Waiters wake, re-enter pipeline, find no reservation, may establish new one.

This order ensures: (a) no waiter can observe a partially-removed state (mutex
protects the entire sequence); (b) the qualifying owner, upon reacquiring the
mutex for its own commit or upon observing its `done` channel closed (before
mutex acquisition), reads `abortOutcome` from its held reservation pointer and
returns the corresponding error without attempting commit.

#### Owner outcome retention after map removal

- **Owner succeeds:** The idempotency reservation transitions to `Completed`
  with the stored response. For qualify, the qualification reservation is
  removed. The `done` channel is closed (waking waiters). Waiters that wake
  and re-enter will find a `Completed` reservation and receive the replay.
- **Owner aborts:** `releaseReservation` removes reservation(s) from the maps.
  The owner retains a direct Go pointer to the reservation object containing
  `abortOutcome`. After mutex release, the owner reads `abortOutcome` to
  determine its response. The object is GC-eligible only after the owner
  returns. No completion is stored; waiters that wake find no reservation and
  may become new owners.

#### Retention and eviction

Retention is 24 hours and 10,000 completed entries; eviction is at-expiry then
earliest-expiry/lexical-namespace at capacity. `InFlight` entries never evict.
Invalid, unknown, and duplicate bodies are never digested or reserved; for
collection create, phase one never digests or reserves a body (DD-11).

### DD-06 — Injected clock drives expiry scheduling, audit-debt lifecycle, and shutdown

A single injected clock (a `Clock` interface per DD-13.4 with a real monotonic
implementation for production and a deterministic fake for tests, realized in
`internal/executiontarget/clock.go`) drives fact `expiresAt` freshness checks
and the expiry-audit scheduling loop. Expiry runs as the sole background worker
with no request caller; it never returns 500 and writes exactly one expiry
AuditEvent per expiry event once its append succeeds.

#### Expiry entry model (keyed by target UID + FactSet identity)

The expiry map is keyed by `(target UID, factSetRef)` — not by target UID
alone. This allows at most one current **pending** timer per target while
retaining any number of historical **auditDebt** obligations for superseded
FactSets. Each entry exists in exactly one of three states:

1. **`pending`** — a timer is scheduled for `expiresAt`; the timer has not
   fired. A superseding target state change (new qualification, Maintenance
   entry, retirement) may cancel the timer and remove this entry. No audit
   obligation exists yet. At most one `pending` entry exists per target at
   any time (enforced by the scheduler: scheduling a new pending timer for a
   target first cancels and removes any existing pending entry for that
   target).

2. **`auditDebt`** — the timer has fired and the expiry event is now a durable
   in-process audit obligation. The AuditEvent has not yet been successfully
   appended. A superseding target state change (new qualification, Maintenance,
   retirement) **must not** remove or cancel this entry; the audit obligation
   persists until append succeeds. Only a successful `AuditAppender.Append`
   transitions the entry to `audited` and removes it. Multiple `auditDebt`
   entries may coexist for the same target (one per distinct superseded
   factSetRef), each with its own independent retry timer.

3. **`audited`** — the AuditEvent was successfully appended. The entry is
   immediately removed from the map. This is a terminal state.

The critical invariant is: once a timer fires and the entry transitions to
`auditDebt`, no target state change may erase that obligation before the append
succeeds (VS0-CF-F16-79).

#### Expiry scheduling and timer mechanics

When a qualification concludes and a new `NormalizedTargetFactSet` is linked,
the scheduler (under the lifecycle-service mutex) performs:

1. If a `pending` entry exists for this target (with any factSetRef): stop its
   timer and remove it. (Only one pending entry per target; a superseded
   pending timer never fires.)
2. Register a new `pending` entry keyed by `(target UID, newFactSetRef)`.
3. Schedule a clock timer (via `Clock.AfterFunc`) set to fire at `expiresAt`,
   storing the `Timer` handle on the entry.

**Superseding while `pending`:** If a new qualification replaces the current
FactSet, or Maintenance/Retirement commits, and a pending entry exists for this
target: the timer is stopped, and the pending entry is removed. No audit
obligation exists.

**Superseding while `auditDebt`:** If a new qualification replaces the current
FactSet, or Maintenance/Retirement commits, and an `auditDebt` entry exists for
this target (with the old factSetRef): the entry is **not** removed. Its retry
loop continues independently until append succeeds. The entry's `factSetRef`
remains the original (now-superseded) one, and the AuditEvent records the
historical expiry. After successful append, the entry is removed. Any new
qualification's pending timer operates independently.

#### Timer ownership and generation tracking

Each `pending` entry stores:
- `timer`: the `Timer` handle from `Clock.AfterFunc`, used for `Stop()`.
- `generation`: a monotonically increasing per-target generation counter,
  incremented each time a new pending entry is created for that target. The
  timer callback carries the generation it was scheduled with.

Each `auditDebt` entry stores:
- `retryTimer`: the `Timer` handle for the current retry, used for `Stop()` at
  shutdown.
- `retryGeneration`: incremented at each retry schedule. The retry callback
  carries this value.

This explicit timer ownership ensures deduplication: a callback whose carried
generation does not match the entry's current generation is a stale no-op.

#### Timer callback and retry

When a pending timer fires:

1. The callback calls `admissionGate.Enter()` (see shutdown protocol below).
   If rejected (gate closed), the callback returns immediately (no-op).
2. `defer admissionGate.Exit()` ensures exit on all paths.
3. Acquire the lifecycle-service mutex.
4. Look up the entry by `(target UID, factSetRef)`. Check:
   - If not found → release mutex; return (entry was removed while pending).
   - If found but not `pending` → release mutex; return (already auditDebt;
     stale callback from a stopped-but-already-fired timer).
   - If found and `pending` but entry generation ≠ callback's carried
     generation → release mutex; return (stale callback).
5. Transition the entry to `auditDebt`.
6. Construct the expiry AuditEvent; call `AuditAppender.Append`.
7. On success: remove the entry (terminal `audited`); release mutex; return.
8. On failure: schedule a retry timer for one clock-second later via
   `Clock.AfterFunc`, storing the handle as `retryTimer` and incrementing
   `retryGeneration`; release mutex; return.

Retry callbacks follow the same protocol:

1. `admissionGate.Enter()` → if rejected, return.
2. `defer admissionGate.Exit()`.
3. Acquire lifecycle-service mutex.
4. Look up entry by `(target UID, factSetRef)`. Check:
   - If not found → release mutex; return (impossible in correct operation
     since `auditDebt` entries cannot be externally removed, but defensive).
   - If `retryGeneration` ≠ callback's carried generation → release mutex;
     return (stale).
5. Call `AuditAppender.Append`.
6. On success: remove entry; release mutex; return.
7. On failure: schedule next retry (same protocol as step 8 above); release
   mutex; return.

#### Deduplication proof

At most one in-flight callback is effective per entry at any time:
- For `pending` entries: at most one pending timer exists per target. The
  generation check rejects stale callbacks from a previously-stopped timer
  that fired between `Stop()` and `AfterFunc` completion.
- For `auditDebt` entries: at most one retry timer is active (stored as
  `retryTimer`). The `retryGeneration` check rejects stale retry callbacks.
- Result: exactly one successful AuditEvent per expiry event.

#### Fake-clock test behavior

With the deterministic fake clock:
- `Advance(d)` fires all timers whose deadline ≤ new current time,
  synchronously, in deadline order.
- No wall-clock time elapses; tests can advance the clock by 60 seconds and
  observe the expiry callback fire immediately.
- Retry: advancing the clock by one second fires the retry timer
  synchronously; tests can supply a failing then succeeding AuditAppender to
  prove exactly one successful AuditEvent with zero wall-clock wait.

#### Shutdown protocol — mutex-serialized admission gate

On shutdown signal, the scheduler must guarantee that no AuditEvent is appended
after the shutdown boundary (VS0-CF-F16-77), and that no Maintenance trigger
can commit after the boundary. The design uses a **mutex-serialized admission
gate** that eliminates the `sync.WaitGroup` Add/Wait race entirely:

**Admission gate structure.** The `admissionGate` uses a single `sync.Mutex`
(the **gate mutex**, distinct from the lifecycle-service mutex) plus a counter
and closed flag, all protected by the gate mutex:

```text
type admissionGate struct {
    mu       sync.Mutex
    closed   bool       // initially false
    inflight int        // initially 0
    drained  chan struct{} // allocated; closed when inflight reaches 0 after close
}
```

- `Enter() → bool`: lock gate mutex; if `closed`, unlock and return false;
  otherwise increment `inflight`, unlock, return true.
- `Exit()`: lock gate mutex; decrement `inflight`; if `closed` and
  `inflight == 0`, close `drained` channel; unlock.
- `Close()`: lock gate mutex; set `closed = true`; if `inflight == 0`, close
  `drained` immediately and unlock; otherwise unlock and block on `<-drained`.

**Race-safety proof.** All state (`closed`, `inflight`) is accessed only under
the gate mutex. No atomic/WaitGroup interaction exists. An `Enter` that
observes `closed=false` and increments `inflight` does so in a single mutex
hold; `Close` sets `closed=true` in a single mutex hold. These two operations
are totally ordered by the gate mutex. If `Enter` acquired the gate mutex
first, `inflight` is already incremented before `Close` inspects it; `Close`
will wait on `drained`. If `Close` acquired the gate mutex first, `closed` is
already true when `Enter` inspects it; `Enter` returns false. No interleaving
between flag-check and counter-increment is possible. QED.

**Admission scope.** The admission gate governs entry for:
- Expiry timer callbacks (pending and retry).
- Maintenance-trigger intake (DD-07 fixture triggers).

Both must call `admissionGate.Enter()` before they can acquire the
lifecycle-service mutex or call `AuditAppender.Append`. This ensures the
shutdown boundary applies uniformly to all background/trigger-initiated
mutations and audit appends.

Request-triggered create, qualify, and retire use a separate
`acceptingMutations` flag guarded by the lifecycle-service mutex. Reservation
creation and every audited publication recheck this flag while holding that
mutex. It is not a public response or violation: once false, an owner follows
the existing internal abort cleanup and no new InFlight reservation or
publication is permitted.

**Shutdown sequence:**

1. **Close the admission gate.** `admissionGate.Close()` prevents new timer
   callbacks and Maintenance triggers from entering, then blocks until all
   currently-admitted callbacks/triggers exit. After `Close` returns, no
   callback or trigger can call `AuditAppender.Append` or commit a Maintenance
   mutation.

2. **One lifecycle-mutex shutdown section.** Acquire the lifecycle mutex once;
   set `acceptingMutations=false`; stop all pending/retry timers; detach every
   InFlight create, qualify, and retire reservation; and mark every remaining
   pending/auditDebt entry invalidated. Collect detached channels, then release
   this mutex. Timers admitted before step 1 already drained in `Close()`.

3. **Wake and stop.** Close each collected channel exactly once after unlock,
   then stop HTTP serving.

**Post-close guarantee.** After step 1 returns, for background work only:
- Every callback/trigger that was admitted completed before `Close` returned.
- No future callback/trigger can pass `Enter`.
- Therefore no expiry or Maintenance-trigger `AuditAppender.Append` call and
  no Maintenance-trigger mutation/audit occurs after the close boundary.

This protocol does not depend on clearing committed target state. The last
committed target projection, ETag, and record links remain unchanged through
shutdown (VS0-CF-F16-77).

#### Timer callback entry protocol

Every timer callback (pending or retry) follows this exact sequence:

```text
1. if !admissionGate.Enter() → return (gate closed; no-op)
2. defer admissionGate.Exit()
3. Acquire lifecycle-service mutex
4. Check entry state and generation (see timer callback section above)
   - If stale/missing/invalidated → release mutex; return
5. [If pending: transition to auditDebt]
6. Construct AuditEvent; call AuditAppender.Append
7. On success: remove entry; release mutex; return
8. On failure: schedule retry timer; release mutex; return
```

The admission gate Enter is the first operation; Exit is guaranteed on every
path via defer (including panic).

### DD-07 — Maintenance enters/clears only through a target-bound fixture trigger

Maintenance enter/clear is delivered exclusively through a deterministic,
target-bound synthetic-observer fixture trigger carrying target UID, expected
InfrastructureStack generation, expected maintenance epoch, operation, system
identity, and correlation. It is not an HTTP route and not a generic event bus.
The lifecycle service validates the fence (matching current epoch/generation)
before applying; stale-fenced triggers produce no state change, AuditEvent, or
record publication (VS0-CF-F16-80,83).

### DD-08 — Pre-ServeMux transport guard wraps the F0016 mux registrations

Following the FEATURE-0015 live convention of centralized route registration in
`internal/server/routes.go`, FEATURE-0016 registers exactly five Go 1.22
`http.ServeMux` method/path patterns and wraps them with a fixed pre-ServeMux
method/path guard (an `http.Handler` decorator realized in
`internal/server/executiontarget_guard.go`, installed ahead of the mux). The
guard supplies transport-only 405/404 outcomes for HEAD, PUT, other unmatched
methods, and every trailing-slash variant, each with
`X-Content-Type-Options: nosniff` and no Problem body, authentication,
authorization, audit, idempotency, or lifecycle effect. Valid registered pairs
bypass the guard and enter the handler at authentication.

### DD-09 — Request pipeline is an explicit ordered sequence

The handler realizes the fixed precedence from REQ-F16-05 as an explicit,
short-circuiting sequence of stages, reusing FEATURE-0012 pipeline stages
(`internal/apivalid`) for decode/classification/structural/reference validation
and FEATURE-0015 conventions (`internal/cloudmodel` grants + safe resolver) for
authentication, authorization, and safe access. The order is fixed and testable;
see §5. The collection-create route additionally realizes the bounded-read →
phase-one-extraction → authorization/safe-access → single-strict-classification
precedence fixed by DD-11.

### DD-10 — Reuse inherited envelopes, errors, and helpers by reference

Problem Details, ObjectMeta, TypedRef, Condition, media handling, ETag/If-Match
handling, and the AuditEvent append boundary are reused from their canonical
owners (`internal/apiproblem`, `internal/apimeta`, `internal/apiref`,
`internal/apicond`, `internal/apivalid`, `internal/decision`). FEATURE-0016
defines no parallel envelope and adds only the eight new violation codes
(`internal/executiontarget/violations.go`) to the violation surface (top-level
Problem codes remain FEATURE-0012-owned).

### DD-11 — Collection-create phase-one and single-strict-classification precedence (ADH-2026-063/065)

For the create-only route
`POST /apis/execution.sovrunn.io/v1alpha1/execution-targets`, the collection
handler realizes the exact observable precedence fixed by ADH-2026-063 and made
exact by ADH-2026-065. This is a request-ordering mechanic; it introduces no new
route, field, resource, state, top-level Problem code, dependency, or external
effect. The ordered mechanic is:

1. Transport guard (DD-08), authentication, and method/media/header-form
   validation run first, unchanged. A collection create carrying an `If-Match`
   header is `MALFORMED_REQUEST`/400 (VS0-CF-F16-105).
2. A bounded single body read runs next, reusing the inherited FEATURE-0012
   single-body-read limit. If the body exceeds that limit, the request returns
   the existing `REQUEST_TOO_LARGE`/400 **before** phase-one extraction,
   authorization, or safe access (VS0-CF-F16-126, anchored on VS0-CF-F16-53).
   Otherwise the exact bytes are retained.
3. Phase one extracts **only** `spec.cloudProviderParticipationRef.uid` and
   `spec.infrastructureStackRef.uid` from the retained bytes. It emits no
   body-classification Problem and never classifies, canonicalizes, digests, or
   reserves a body. If it cannot extract exactly one syntactically usable UID at
   **both** required locations, the create fails closed with the existing
   audited `AUTHORIZATION_DENIED`/403 **before** derived-scope authorization and
   safe backing access, disclosing no body-classification detail and no
   backing-resource existence, and performing no strict classification,
   canonicalization, digest, idempotency reservation, observer invocation,
   mutation, publication, or completion (VS0-CF-F16-128, anchored on
   VS0-CF-F16-08).
4. Derived-scope `executiontarget.write` authorization and safe backing access
   are evaluated using the phase-one-extracted references. An authenticated
   caller lacking the grant over safe-accessible backing whose body is malformed
   JSON or a duplicate top-level member retains denial-family precedence with the
   existing audited `AUTHORIZATION_DENIED`/403 and no body classification
   (VS0-CF-F16-123, anchored on VS0-CF-F16-08). An inaccessible backing reference
   retains safe-denial precedence with the existing `RESOURCE_NOT_FOUND`/404 +
   `VS0_AUTHORIZATION_SAFE_DENIAL` and no body classification (VS0-CF-F16-124,
   anchored on VS0-CF-F16-09).
5. Only for an authorized, safe-accessible request are the retained same bytes
   strictly classified **exactly once**. Malformed JSON only returns the existing
   `MALFORMED_REQUEST`/400 (VS0-CF-F16-125, as VS0-CF-F16-51); a duplicate
   top-level member only returns the existing `DUPLICATE_FIELD`/400
   (VS0-CF-F16-127, as VS0-CF-F16-52). These are distinct exact non-family
   classification cases, not a mixed equivalence family. Graph/reference
   validation, idempotency, audit, and publication follow, unchanged.

Action `If-Match` header-form precedence (§5.1 step 10; VS0-CF-F16-70..74,86) is
preserved and unaffected by this create-only mechanic. VS0-CF-F16-123 and 124 use
named finite equivalence families (malformed JSON or duplicate top-level member),
each member producing the identical declared observable outcome.

### DD-12 — Maintenance-entry link-clearing and target-local abort (ADH-2026-061)

Every successful Maintenance entry committed by the lifecycle service
unconditionally clears both the current `NormalizedTargetFactSet` and
`TargetQualificationResult` links and persists `Active/Unqualified`, advancing
maintenance epoch and resourceVersion/ETag, whether or not a qualification is in
flight; it never preserves the last completed conclusion, and the successful
entry is always audited (VS0-CF-F16-42). When a qualification is in flight at
Maintenance entry, the entry additionally aborts only that target's in-flight
qualification: the lifecycle service (under the same mutex hold as the
Maintenance commit) calls
`releaseReservation(qualifyIdempotencyKey, VS0_TARGET_MAINTENANCE)` (DD-05).
The `releaseReservation` primitive applies the first-winner rule: if the
qualification reservation's `abortCause` is already set (a prior Retire won),
it is not overwritten; otherwise it is set to `VS0_TARGET_MAINTENANCE`. The
primitive removes both the qualification reservation and its linked idempotency
reservation, detaches its `done` channel under the mutex, and wakes only that qualification's
waiters — never all reservations for the target. The qualifying caller keeps
the response determined by `abortOutcome` (which is `VS0_TARGET_MAINTENANCE`
if Maintenance was the first winner, or `VS0_TARGET_RETIRED` if a prior Retire
was) with no qualification result, completion, or qualification AuditEvent
(VS0-CF-F16-89). When no qualification is in flight, no idempotency reservation
is cancelled, and the Maintenance entry is still audited. A successful
Maintenance clear likewise persists `Active/Unqualified` with both current links
cleared and is audited (VS0-CF-F16-43). VS0-CF-F16-42, VS0-CF-F16-43, and
VS0-CF-F16-89 retain their exact observable semantics unchanged. This is a
commit-composition mechanic within DD-02/DD-07; it adds no state, field, or
violation.

### DD-13 — Internal interfaces, immutable-value contracts, and lock ordering

This design decision defines the internal Go interface contracts and lock
ordering for the five cross-component boundaries consumed by the lifecycle
service. These are F0016-internal contracts; they add no route, field, state,
violation, or external effect.

#### DD-13.1 Observer interface

```text
Interface: Observer
Method:    Observe(ctx, target, clock) → (ObservationResult, error)
```

Immutable-value contract:
- `ObservationResult` carries exactly four named truth values, `observedAt`,
  `expiresAt`, observer ID, observer revision, and `factSchemaVersion=v1`
  provenance. Each truth value is exactly one of Supported/Unsupported/Unknown.
- The observer never mutates target state, never acquires the lifecycle-service
  mutex, never performs network I/O, and never blocks on a real clock.
- A missing fixture produces four Unknown facts at `now`/`now+60s`. A
  fixture-declared logical timeout produces four Unknown facts at
  `now`/`now+60s`. Both are resolved against the injected clock value at
  invocation time.
- An error return indicates an observer fault (duplicate/missing/malformed
  provenance); the lifecycle service treats this as a commit abort.

Lock ownership: Observer holds **no lock**. It is invoked by the qualify path
with the lifecycle-service mutex **unlocked**.

#### DD-13.2 BackingAccessProvider interface (ADH-2026-066)

`BackingAccessProvider` is an F0016-internal, principal-aware boundary.  It is
the only F0016 component permitted to acquire the private FEATURE-0015 paired
backing-read lease.  The lease returns immutable copies of the requested
CloudProviderParticipation and InfrastructureStack under one FEATURE-0015
store lock; it is not a new F0015 public API, writer, resource, route, or
schema.

```text
Access(ctx, principal, routeIntent, participationUID, stackUID, grantEvaluator)
  -> SafeDenied | AuthorizationDenied | Allowed(BackingView)

WithFinalQualificationLease(ctx, principal, participationUID, stackUID,
  grantEvaluator, func(AllowedBackingView) CommitDisposition) -> CommitDisposition
```

While the lease is held, the provider evaluates inherited FEATURE-0012 grants
and returns only the caller-relative safe disposition and the immutable values
needed by the route: derived CloudProvider scope,
effective participation state, stack phase/generation, and the deterministic
viability fingerprint.  It must not expose a backing resource's existence to an
unauthorized caller.  An inaccessible backing reference therefore follows the
registered audited safe-404 behavior, rather than an existence-derived result.

The provider has three uses:

1. **Collection create:** after bounded read and successful phase-one UID
   extraction, it supplies authorization scope, safe backing access, and graph
   validation.  The lease is released before strict classification and no
   backing generation is persisted as a create fence.
2. **GET/LIST:** after route authorization, it supplies a fresh response-only
   viability view for `effectiveAvailability`; this cannot mutate target state
   or ETag.
3. **Qualify:** it supplies initial backing admission and the initial registered
   fences.  At final commit it supplies a fresh lease that remains held through
   the lifecycle mutex, inherited audit append, and publication.

Only three qualification fences exist: InfrastructureStack generation,
target-local maintenance epoch, and viability fingerprint. Participation
generation is not a fence. A participation suspension means unavailable for
new admission only; it never stops or mutates existing stacks, targets, or
workloads.

The required cross-package order is:

```text
FEATURE-0015 paired backing-read lease
  -> FEATURE-0016 lifecycle-service mutex
  -> inherited AuditAppender lock
  -> F0016 publication
```

No code holding the lifecycle mutex may acquire the lease. FEATURE-0015
writers never acquire the lifecycle mutex. The final qualification recheck is
performed while the final lease remains held, so no FEATURE-0015 backing write
can interleave between that recheck and publication.

#### DD-13.3 Audit appender interface

```text
Interface: AuditAppender
Method:    Append(ctx, AuditEvent) → error
```

Immutable-value contract:
- The `AuditEvent` value passed is a complete, redacted FEATURE-0013 record.
  Once constructed it is never mutated by the caller.
- An `error` return means the append failed; the caller must not publish the
  associated mutation.
- The appender never acquires the lifecycle-service mutex.

Lock ownership: AuditAppender holds **no F0016 lock**. It is called with the
lifecycle-service mutex **held** at commit time (the audit append is part of
the atomic commit sequence). The appender internally may acquire its own
FEATURE-0013 lock; there is no cross-package circular dependency because
FEATURE-0013 never calls back into FEATURE-0016.

Cross-package lock ordering (global): `lifecycle-service mutex` →
`AuditAppender internal lock`. Never the reverse.

#### DD-13.4 Clock interface and scheduling primitives

```text
Interface: Clock
Method:    Now() → time.Time
Method:    NewTimer(d) → Timer
Method:    AfterFunc(d, fn) → Timer

Interface: Timer
Method:    Stop() → bool
Method:    C() → <-chan time.Time   (only for NewTimer timers)
```

Immutable-value contract:
- `Now()` returns a monotonic instant; successive calls are non-decreasing.
- `NewTimer(d)` returns a Timer that fires exactly once after duration `d`
  relative to the clock's current time.
- `AfterFunc(d, fn)` schedules `fn` to run once after `d`; `Stop()` returns
  true if the firing was prevented.
- In the fake clock, advancing time by `Advance(d)` fires all timers whose
  deadline is at or before the new current time, synchronously and
  deterministically, in deadline order. No wall-clock time elapses.
- In the real clock, `NewTimer`/`AfterFunc` delegate to the standard library.

Lock ownership: The Clock itself holds **no F0016 lock**. Timer callbacks
(used by the expiry scheduler) acquire the lifecycle-service mutex only when
they need to commit an expiry AuditEvent or update pending-expiry state.

#### DD-13.5 Lifecycle coordinator commit contract (ADH-2026-066)

The lifecycle service is the sole F0016 publication authority. Idempotency
admission and acquisition of the target-local Qualifying reservation are
separate: replay/owner admission occurs first; current If-Match, safe backing
admission, and lifecycle preconditions occur next; only then may the service
capture maintenance epoch and acquire Qualifying state. This preserves the
registered action precedence and prevents a failed earlier check from creating
Qualifying work.

Observer work runs without a backing lease or lifecycle mutex. Its final
qualification commit is one critical section:

```text
acquire F0015 paired backing-read lease
→ acquire F0016 lifecycle mutex
→ reject closed admission or a finalized owner reservation
→ resolve active Maintenance, Retired, maintenance epoch,
  Stack generation, and viability-fingerprint predicates
→ construct and append AuditEvent
→ only on append success publish links, ETag, completed replay, and expiry state
→ detach under mutex; unlock; close captured done channel once; release lease
```

The only backing fences are InfrastructureStack generation and viability
fingerprint; the only target-local fence is maintenance epoch. Every failure
detaches its linked reservation under the mutex, retains the owner outcome,
unlocks, and then closes the captured channel exactly once. It cannot publish a
result, advance ETag, or create replayable completion.

Active Maintenance wins before Retired, followed by maintenance epoch, Stack
generation, and viability fingerprint. This is ADH-2026-060's mutually
exclusive rule. ADH-2026-061 remains intact: a successful Maintenance entry
clears both links and aborts only the same target's in-flight qualification and
linked idempotency reservation.

Shutdown closes lifecycle/idempotency admission under this same mutex before
detaching all InFlight create, qualify, and retire reservations. Each owner
rechecks admission before publication. `abortReasonShutdown` is an internal
cleanup reason, never a public violation code.

### DD-14 — Phase-one UID extraction algorithm (collection create only)

Phase-one extraction is a non-classifying, security-critical extraction that
runs after the bounded body read succeeds and before derived-scope authorization
or safe backing access. Its sole purpose is to extract exactly two UID values
from the retained raw bytes so the pipeline can derive scope for authorization.
It introduces no route, field, state, violation, or external effect.

#### Algorithm specification

Input: the retained raw bytes from the successful bounded body read (guaranteed
≤ the bounded-read limit).

**Scanner model.** The scanner is a forward-only, single-pass, byte-level
JSON tokenizer. It maintains a logical nesting stack tracking the current
object-key path from root. It processes tokens in strict left-to-right order
with no backtracking. The scanner's purpose is solely to identify string values
at the two target paths; it is not a JSON decoder and does not construct a
parsed tree.

**Token semantics:**

- **Object key recognition.** An object key is a JSON string token
  (RFC 8259 §7) immediately preceding a `:` separator within an object
  context. Key equality uses the **decoded** key value: JSON escape sequences
  (`\"`, `\\`, `\/`, `\b`, `\f`, `\n`, `\r`, `\t`, `\uXXXX`, and
  `\uXXXX\uXXXX` surrogate pairs) are decoded to their UTF-8 byte
  representation before comparison. Two keys are equal if and only if their
  decoded UTF-8 byte sequences are identical. This means `"spec"` and
  `"\u0073pec"` are the same key.

- **String value extraction.** A value at a target path is recognized as a
  JSON string if the token begins with `"` and the scanner can identify the
  closing unescaped `"`. The extracted candidate UID is the decoded content
  between the quotes (with escape sequences resolved to UTF-8). If decoding
  encounters an invalid escape sequence within a target-path string value, the
  value is treated as unextractable at that occurrence (not a scanner abort).

- **Invalid UTF-8 handling.** The scanner operates on raw bytes. JSON strings
  containing bytes that are not valid UTF-8 (outside of `\uXXXX` escapes)
  are not rejected by the scanner; they are passed through as raw bytes. If
  such bytes appear in an extracted candidate UID, the canonical UID grammar
  check (step 3) will reject them since the grammar requires ASCII-only
  characters.

- **Malformed JSON tokens.** Tokenization follows RFC 8259 exactly. On the
  first malformed key, escape, surrogate pair, delimiter, literal, number, or
  unterminated string, the scanner makes one deterministic choice: if both
  complete target paths have already occurred exactly once with valid canonical
  UID strings, it returns those two extracted values and leaves the malformed
  body for the single later strict classification; otherwise it returns
  unextractable for every target path not already completed. It never attempts
  speculative recovery, balanced-pair guessing, or discretionary skip mode.
  It advances monotonically to that decision or EOF and never recurses.

- **Structural nesting.** The scanner maintains a bounded nesting depth. For
  the two target paths (`spec.cloudProviderParticipationRef.uid` and
  `spec.infrastructureStackRef.uid`), the maximum meaningful depth is 4 levels
  (root object → `spec` → ref object → `uid`). The scanner tracks depth up to
  a fixed maximum (≥ the body size limit / 2, to handle pathological nesting).
  Objects/arrays at depths beyond the target paths are skipped without
  descending: the scanner tracks balanced braces/brackets but does not push
  them onto the path stack.

- **Duplicate path counting.** Every occurrence of either complete target path
  is counted when its decoded key sequence is reached at the correct nesting
  depth, irrespective of whether its value is a string, number, boolean, null,
  object, or array. Duplicate intermediate keys (including duplicate `spec`)
  contribute their own sub-tree traversals. A duplicate is therefore always
  ambiguous; a string/non-string pair cannot be accepted.

**Steps:**

1. **Scan the retained bytes.** The scanner processes the bytes left-to-right,
   tracking the nesting path. For each target path it counts every occurrence
   and retains a candidate only when its sole occurrence is a decoded string.

2. **Occurrence counting rules:**
   - Duplicate top-level `spec` keys: each occurrence is visited.
   - Nested duplicate keys within a `spec` object: each occurrence is visited.
   - A non-string value counts as an occurrence and is unextractable.
   - A candidate UID exists only when exactly one occurrence exists and that
     occurrence is a decoded string.
   - Zero or more-than-one occurrence makes the path unextractable.

3. **Canonical UID grammar validation.** Each extracted candidate UID is
   validated against the existing canonical UID grammar (the same syntactic
   validation reused by FEATURE-0012 `TypedRef` UID fields). This is a pure
   syntactic check; it does not perform any backing resolution or existence
   test. A candidate that does not match the canonical UID grammar is
   unextractable.

4. **Success predicate.** Phase one succeeds if and only if each path occurs
   exactly once and each sole value is a syntactically valid canonical UID.
   The two extracted UIDs are the phase-one result.

5. **Fail-closed result.** If either path is unextractable (zero occurrences,
   more-than-one occurrences, non-string value, invalid escape in target value,
   or invalid UID grammar), phase one fails closed. The pipeline returns the
   existing audited `AUTHORIZATION_DENIED`/403 (VS0-CF-F16-128) before
   derived-scope authorization and safe backing access, disclosing no
   body-classification detail and no backing-resource existence.

#### Bounded linear-time guarantee

The scanner processes each byte of the input at most a constant number of times
(once for primary scan; balanced-pair tracking in skip mode revisits bytes but
advances monotonically). Total time is O(N) where N is the byte length of the
retained body (which is bounded by the body-read limit). No recursion, no
backtracking, no quadratic re-scanning.

#### Guarantees

- Phase one **never** emits a body-classification Problem (no
  `MALFORMED_REQUEST`, `DUPLICATE_FIELD`, `REQUEST_TOO_LARGE`, or
  `UNKNOWN_FIELD`). All body-classification outcomes are deferred to strict
  classification (step 7 in §5.1), which runs exactly once, only for an
  authorized, safe-accessible request.
- Phase one **never** classifies, canonicalizes, digests, or reserves a body.
  No idempotency digest is computed from the retained bytes at this stage.
- Phase one **never** performs strict JSON decoding. Malformed JSON (missing
  commas, trailing garbage, unclosed strings) may still yield extractable UIDs
  at the two paths; the malformed body is then correctly rejected later by
  strict classification (VS0-CF-F16-125) only if authorization and safe access
  pass. If the malformed body prevents extraction of both UIDs, phase one fails
  closed (VS0-CF-F16-128).
- The retained raw bytes are passed unchanged to strict classification; phase
  one does not modify, copy with normalization, or re-encode them.

#### Duplicate behavior

A duplicate `spec` top-level key in the retained JSON means the scanner visits
both `spec` objects. If both contain `cloudProviderParticipationRef.uid` as a
string, the occurrence count for that path is ≥ 2 and the path is
unextractable. This is the correct behavior: the body will later be rejected by
strict classification as `DUPLICATE_FIELD`/400 if it reaches that stage, but
phase one itself emits no classification Problem and fails closed with 403.

#### Security-critical outcome determinism

For a given byte sequence, the scanner produces a deterministic result. Two
implementations conforming to this specification must agree on the extraction
outcome for any input. The critical security boundary is: if the scanner
reports "extractable" for both paths, the two UIDs are unambiguously identified
from the raw bytes. If there is any ambiguity (duplicate paths, malformation
that obscures boundaries, escape-decoding failure), the scanner reports
"unextractable" and the pipeline fails closed with 403 — never proceeding to
authorization with a potentially-incorrect scope derivation.

### DD-15 — Fresh backing reads for GET/LIST effectiveAvailability projection

`effectiveAvailability` must reflect viability-only changes (backing
participation becoming ineffective, or stack leaving Active phase) without
mutating the target's ETag (VS0-CF-F16-119,120). This means GET and LIST must
read current backing viability state at response time, not rely solely on the
last-committed target state.

#### Item GET projection

After authorization succeeds, the item GET handler calls the F0016-owned
`BackingAccessProvider` using the target's persisted backing UIDs. This call
occurs **outside** the
lifecycle-service mutex (it is a read-only backing lookup). The returned
view's `ParticipationEffectiveActive` and `StackPhase` feed the
`effectiveAvailability` truth table (§4.6) together with the target's committed
lifecycle/qualification state and current-Maintenance marker.

The ETag returned in the response is the target's committed `resourceVersion`
(unchanged by the fresh backing read). A viability-only change therefore does
not alter the ETag — the response reflects current availability without
implying a target mutation.

Safe-denial ordering is preserved: if the target is inaccessible at the
UID-only safe resolution step (§5.1 step 4), the handler returns safe 404
before reaching backing-access acquisition.

#### Collection LIST projection

After authorization succeeds, the LIST handler iterates accessible targets in
ascending-UID order. For each target, it calls the F0016-owned
`BackingAccessProvider` to compute
`effectiveAvailability` at projection time. Each call is outside the
lifecycle-service mutex. LIST has no ETag.

Deterministic ascending-UID ordering is maintained because the iteration key
is the target's `metadata.uid` (stable), not a field derived from the fresh
view.

#### Performance note (design-local, not observable)

These backing reads use the private paired immutable backing-read lease defined
by DD-17. For LIST,
the total time is O(N) where N is the number of accessible targets. No network
I/O or blocking occurs.

### DD-16 — Derived-scope index: internal non-authoritative cache

At ExecutionTarget create time, the lifecycle service records the mapping
`(target UID → CloudProvider UID)` in a process-local index, derived from the
participation's `Spec.CloudProviderRef.UID` as read from the initial
`BackingAccessProvider` lease. This index is:

- **Not an ExecutionTarget field.** It does not appear in `metadata`, `spec`,
  or `status` of the ExecutionTarget schema (VS0-SCHEMA-015). The
  ExecutionTarget's persisted schema remains exactly the closed set in §4.1.
- **Not a persisted `scopeRef`.** No `scopeRef` field is persisted or
  projected. Scope derives solely from the immutable
  `spec.cloudProviderParticipationRef.uid` reference at runtime.
- **Not an authoritative scope record.** If the index were unavailable (e.g.
  after restart before re-population), the canonical derivation path is to
  resolve the target's persisted `spec.cloudProviderParticipationRef.uid` via
  `BackingAccessProvider` to obtain `ParticipationCloudProviderUID`.
  The index is a performance optimization that avoids a backing read on every
  item-route request.
- **Populated at create time only.** The `spec.cloudProviderParticipationRef`
  reference is immutable for the lifetime of the target, so the derived
  CloudProvider UID never changes. The index entry is removed when the target
  is removed from the store (restart clears all in-process state).
- **Process-local and restart-cleared.** On restart, the index is empty; it is
  repopulated as targets are recreated. No durable storage or cross-process
  sharing exists.

Item GET/qualify/retire handlers use this index to derive the CloudProvider
scope for authorization without issuing a backing-access call. The index
lookup is a map read under the lifecycle-service mutex. It discloses no backing-resource
existence to the caller because the target must first pass UID-only safe
resolution (step 4 in §5.1) before the index is consulted.

### DD-17 — Coherent, principal-aware backing access (ADH-2026-066)

This decision is implemented by the active DD-13.2/DD-13.5 contracts. No
sequential backing read, participation-generation fence, or lifecycle-first
final commit is an implementation option.

`BackingAccessProvider` is F0016-internal. It accepts the authenticated
principal, F0016 route intent, both backing UIDs, and the inherited F0012 grant
evaluator. It invokes a private `internal/cloudmodel` backing-read lease. The
lease holds one existing store lock while it supplies immutable copies of both
backing resources and while the provider evaluates caller-relative safe access.
The provider returns only safe denial, authorization denial, or an allowed
snapshot containing derived CloudProvider UID, effective-Active participation,
InfrastructureStack phase/generation, and a deterministic viability
fingerprint.

Its result is a closed union, never a raw store result:

```text
SafeDenied              // no raw backing existence or fields escape
AuthorizationDenied     // scope is known only inside the provider
Allowed(BackingView)    // immutable paired values and derived scope
```

`BackingView` contains participation UID, InfrastructureStack UID,
participation effective-Active state, stack phase and generation, derived
CloudProvider scope, and viability fingerprint. Its
copies cannot be mutated by either store after construction. The provider owns
lease release on every safe-denial, authorization-denial, error, cancellation,
or panic path. Ordinary create, GET, and LIST calls release internally before
returning `Allowed`. Only `WithFinalQualificationLease` retains a lease, solely
inside the lifecycle-service callback; it releases with defer/finally after the
prescribed commit section. The provider never returns raw existence booleans,
raw backing objects, or a lease to an HTTP handler.

Initial create/qualify access uses a lease and releases it before observation.
Final qualify commit obtains a fresh lease, compares only the registered
InfrastructureStack-generation and viability-fingerprint fences, and then
acquires the lifecycle mutex while retaining the lease. The fixed order is:

```text
FEATURE-0015 backing lease → F0016 lifecycle mutex → inherited AuditEvent append → F0016 publication
```

A changed registered fence aborts before append, publication, or replayable
completion. Participation generation is not a fence. GET/LIST obtain a fresh
provider result outside the lifecycle mutex solely for response-only
`effectiveAvailability`.

For create, an initial provider `SafeDenied` or `AuthorizationDenied` returns
the already registered audited safe-404 or authorization result and skips
classification, reservation, and publication as required by the route
precedence. For qualification, an initial non-Allowed result or a final-lease
recheck failure detaches the linked owner reservation under the lifecycle mutex,
unlocks, closes its waiter signal once, and produces the registered result. GET
and LIST use only an Allowed view for projection; a denial follows their
existing non-disclosing route result. The derived-scope cache is an optimization
only; a cache miss uses this provider and never falls back to sequential reads.

No F0016 path holding the lifecycle mutex may acquire the backing lease, and no
FEATURE-0015 writer acquires the F0016 mutex. Participation suspension means
unavailable for new admission only: it does not stop, suspend, or mutate an
existing InfrastructureStack, ExecutionTarget, or already-realized workload.

---

## Components and Interfaces

See §3 (components/paths) and §4 (data/API representation) below.

---

## 3. Components and repository paths

All paths below follow verified live conventions (FEATURE-0016 domain package
`internal/executiontarget/`; FEATURE-0015 domain package
`internal/cloudmodel/`; handlers `internal/api/*_collection.go`, `*_item.go`,
`*_actions.go`; mux `internal/server/routes.go` plus
`internal/server/executiontarget_guard.go`; conformance
`tests/conformance/`). Package placement and file decomposition are fixed by
DD-01 and this table; task sequencing and implementation ordering belong to the
tasks stage.

### 3.1 F0016-owned components (IMPLEMENT)

| Component | Path | Responsibility |
|-----------|------|----------------|
| ExecutionTarget domain types | `internal/executiontarget/model/` (`types.go`, `enums.go`) | In-memory structs for ExecutionTarget, NormalizedTargetFactSet, TargetQualificationResult, the current-Maintenance marker, and the viability fingerprint; enum consts for lifecycle/qualification/effectiveAvailability/fact truth values |
| In-memory store | `internal/executiontarget/store.go` | Process-local registry; scope-unique name reservation (retained after retirement); at-most-one live `(participation UID, stack UID, targetClass)` tuple index; ascending-`metadata.uid` LIST ordering; derived-scope index (DD-16) |
| Lifecycle service (VS0-WRITER-006) | `internal/executiontarget/lifecycle.go` | Sole committer of target state + internal records; Qualifying reservation; atomic conclusion publication; maintenance enter/clear with unconditional link-clearing (DD-12); retirement; current-Maintenance marker; effectiveAvailability computation |
| Synthetic observer | `internal/executiontarget/observer.go` | `sovrunn.synthetic-iaas-observer/v1`; target-bound, clock-driven, external-effect-free fact proposal and classified missing-fixture/logical-timeout outcomes; fixed provenance |
| Qualification engine | `internal/executiontarget/qualify.go` | Closed profile evaluation: any Unsupported → Rejected; else any Unknown → Indeterminate; else Qualified; observer-fault detection (duplicate/missing/malformed) |
| Idempotency reservation table | `internal/executiontarget/idempotency.go` | F0016-owned namespace/digest/reservation/waiter table with unified broadcast signal and releaseReservation cleanup primitive (DD-05); retention/eviction; waiter recheck of current authorization + safe access |
| Expiry/maintenance scheduler + clock | `internal/executiontarget/scheduler.go`, `internal/executiontarget/clock.go` | Injected-clock expiry worker with (target UID, factSetRef)-keyed entry lifecycle, timer generation tracking, once-per-clock-second audit retry; target-bound fixture maintenance-trigger intake; mutex-serialized admission gate shutdown covering both timer callbacks and Maintenance triggers (DD-06) |
| Safe projection | `internal/executiontarget/projection.go` | Closed safe ExecutionTarget JSON projection incl. response-only effectiveAvailability; omits facts/results/observer/handles/links |
| Audit assembly | `internal/executiontarget/audit.go` | Builds redacted FEATURE-0013 AuditEvents for the F0016 audit matrix; delegates append to the inherited FEATURE-0013 append boundary |
| Principal-aware backing access | `internal/executiontarget/backing_access.go`, `backing_access_test.go` | F0016-owned `BackingAccessProvider`: converts the private paired lease into safe-denial, authorization-denial, or immutable allowed backing view; never returns raw backing existence |
| Private paired backing-read lease | `internal/cloudmodel/store.go`, `store_test.go` | ADH-2026-066 F0016 compatibility extension co-located with the inherited store; returns paired immutable copies under one store lock and changes no FEATURE-0015 product contract |
| Eight new violation constants | `internal/executiontarget/violations.go` | F0016-owned `violations[]` constants only; the `internal/apiproblem` top-level Problem enum/map is reused unchanged |
| HTTP handlers | `internal/api/executiontarget_collection.go`, `internal/api/executiontarget_item.go`, `internal/api/executiontarget_actions.go` | Collection POST/GET, item GET, qualify POST, retire POST; pipeline wiring per §5, including the DD-11 collection-create precedence |
| Route registration + transport guard | `internal/server/routes.go` (five registrations) + `internal/server/executiontarget_guard.go` | Five exact ServeMux patterns and the fixed pre-ServeMux method/path guard |

### 3.2 Reused prior-feature components (CONTRACT_ONLY/NO_TASK — consumed by reference)

| Component | Path | Reuse |
|-----------|------|-------|
| Problem Details + HTTP/URN map | `internal/apiproblem` | Top-level Problem codes and mappings; F0016 adds no top-level code |
| ObjectMeta / TypedRef | `internal/apimeta`, `internal/apiref` | Identity/version fields and typed-reference shape validation |
| Validation pipeline | `internal/apivalid` | Decode, strict classification, structural/semantic/reference stages, bounded body read, concurrency (`If-Match`) helpers |
| Condition | `internal/apicond` | Inherited Condition type (not projected by F0016) |
| Authorization / safe resolver / idempotency conventions | `internal/cloudmodel` (grants, safe denial, idempotency digest) | Server-resolved grants, safe-denial semantics, digest algorithm conventions |
| AuditEvent append boundary | `internal/decision` | Atomic append-before-publication; F0016 does not redefine |

### 3.3 Excluded components (EXCLUDED — never designed here)

Real adapters, provider-native types, credentials/SecretRef, external calls,
external persistence, controllers beyond the lifecycle service, placement,
provisioning, plugin execution, customer target projection, maintenance-notice
resources, IAM resources, the concepts named in the non-goals list in §10.1, and
all FEATURE-0017/0019/0022/0023/0024 concepts (§10.1).

---

## Data Models

See §4 (Data/API representation) below.

---

## 4. Data/API representation

### 4.1 ExecutionTarget (VS0-SCHEMA-015)

Create body accepts exactly (closed set): `metadata.name`,
`spec.cloudProviderParticipationRef.uid`, `spec.infrastructureStackRef.uid`, and
immutable `spec.targetClass="synthetic-iaas"`. Server-assigned only:
`metadata.uid`, `metadata.resourceVersion`, `metadata.generation`,
`status.lifecycle`, `status.qualification`, `status.maintenanceEpoch`,
`status.observedGeneration`, `status.factSetRef`, `status.qualificationResultRef`.
Scope derives solely from `spec.cloudProviderParticipationRef`; no `scopeRef` is
persisted or projected. Create persists `Active/Unqualified` at maintenance epoch
0, sets `observedGeneration` to the referenced InfrastructureStack generation,
holds no record refs, and projects `effectiveAvailability=Unavailable`.

Safe projection (create 201 / item GET 200 / qualify 200 / retire 200; LIST
ascending-UID array, no ETag) is exactly: `metadata.uid`, `metadata.name`,
`metadata.generation`, `metadata.resourceVersion`, the three immutable spec
fields, `status.lifecycle`, `status.qualification`, and response-only
`effectiveAvailability`. No public `conditions` member; facts, results, observer
data, handles, and internal record links are omitted. No `Location` header.

### 4.2 NormalizedTargetFactSet (VS0-SCHEMA-016) — internal-only, never projected

Fields: target ref, referenced InfrastructureStack generation, maintenance
epoch, viability fingerprint, observer ID/revision, `factVersion=v1` (the sole
persisted FactSet version field), `observedAt`/`expiresAt`, and the four named
truth values for `compute.vm`, `storage.block`, `storage.object`,
`network.private`, each exactly one of Supported/Unsupported/Unknown.
`factSchemaVersion=v1` is fixed observer provenance describing the emitted fact
schema and is not a second persisted FactSet field.

### 4.3 TargetQualificationResult (VS0-SCHEMA-017) — internal-only, never projected

Fields: target ref and fact-set ref, only the registered
`infrastructureStackGeneration` and `maintenanceEpoch` fences (the FactSet is
referenced by ref; `viabilityFingerprint` is FactSet-only and is not duplicated
on Result), profile version, outcome (Qualified/Rejected/Indeterminate), reason
codes, and evaluation time. Each record is immutable once created, but a target
retains links only to its current FactSet/Result pair: successful qualification
atomically replaces both links, while retirement and every Maintenance entry
clear both links (DD-12).

### 4.4 Fixed provenance and observer identity

`observerID=sovrunn.synthetic-iaas-observer/v1`,
`profileVersion=synthetic-iaas/v1`, `factSchemaVersion=v1` (observer provenance
only), and the immutable selected fixture revision.

### 4.5 HTTP surface (five routes; ADH-2026-058 clause 3)

| Method | Path | Handler | ETag | Idempotency |
|--------|------|---------|------|-------------|
| POST | `/apis/execution.sovrunn.io/v1alpha1/execution-targets` | collection create | yes (201) | yes |
| GET | `/apis/execution.sovrunn.io/v1alpha1/execution-targets` | collection LIST | no | never |
| GET | `/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}` | item GET | yes (200) | never |
| POST | `/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify` | qualify | yes (200) | yes |
| POST | `/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire` | retire | yes (200) | yes |

`{uid}` is a complete Go 1.22 `http.ServeMux` wildcard read with
`r.PathValue("uid")`. Success media `application/json`; error media
`application/problem+json`; every response (including transport-only 404/405 and
replay) sends `X-Content-Type-Options: nosniff`. `Accept` is ignored; no route
accepts a query parameter.

### 4.6 effectiveAvailability projection (response-only; no persisted writer; no ETag effect)

Ordered truth table (from committed state, freshness, maintenance, viability):

1. Retired → `Unavailable`.
2. else Active with current Maintenance marker active → `Maintenance`.
3. else only Active + Qualified + fresh current facts + viable backing → `Available`.
4. every other Active combination → `Unavailable`.

A viability-only change does not mutate the target or change its ETag
(VS0-CF-F16-119,120).

---

## Correctness Properties

Invariants are stated in §9.2; formal models in §7.6. Each property below is
one row of the §9.2 table.

### Property 1: Sole-committer status writes

For all requests and all ExecutionTarget/NormalizedTargetFactSet/
TargetQualificationResult records, only `ExecutionTargetLifecycleService`
mutates status/internal records; no client or handler write occurs (DD-02;
VS0-CF-F16-03,04,91,92).

**Validates: REQ-F16-08**

### Property 2: No external effect

For all observer invocations, no network I/O, credential use, or external
persistence occurs (DD-03; §8.3; VS0-CF-F16-46).

**Validates: REQ-F16-11**

### Property 3: Audit-before-publication

For all committing mutations, the required AuditEvent append succeeds before
the mutation is published; a failed append leaves target, ETag, and record
links unchanged (DD-04, §5.3; VS0-CF-F16-39,40,82,99..104).

**Validates: REQ-F16-09**

### Property 4: Safe-denial non-disclosure

For all requests against an inaccessible target or backing reference, the
response never confirms existence (§5.1 step 4 item UID-only safe resolution,
§5.1 step 6 collection-create safe backing access, §5.4, §6.2;
VS0-CF-F16-09,17,19,66,124,128).

**Validates: REQ-F16-02**

### Property 5: Closed route surface with side-effect-free guard

For all requests, exactly five registered routes exist and the pre-ServeMux
transport guard produces no authentication, authorization, audit, idempotency,
or lifecycle effect (DD-08; VS0-CF-F16-44,113..118).

**Validates: REQ-F16-03**

### Property 6: Stale-fence non-publication

For all qualification commits, a stale captured fence or a
Retired/Maintenance-at-commit predicate at commit time aborts with no target
mutation, AuditEvent, or completion. The pre-commit recheck lease (DD-13.2
linearization rule) obtained outside the mutex detects any FEATURE-0015
mutation that changed a generation or viability fingerprint between initial
capture and recheck. The lifecycle-service mutex protects target-local state
(maintenance epoch, Retired flag, active-Maintenance marker) and ensures those
predicates are authoritative at commit time. Together they ensure no
stale-fenced conclusion publishes. For collection create, no recheck is needed
because no pre-existing target state exists to fence; each individual backing
lease is authoritative through its defined use (DD-13.2). (DD-04,
DD-04a, DD-12, DD-13.5, §5.1 step 12; VS0-CF-F16-29,30,75,81,82,88).

**Validates: REQ-F16-08**

### Property 7: Collection-create phase-one and single-strict-classification precedence

For all collection-create requests, the bounded body read rejects oversized
bodies first; phase one extracts only the two required UIDs and fails closed
with audited 403 when either is unextractable; the BackingAccessProvider view
determines accessibility before authorization; authorization and safe access
precede any body classification; and the retained same bytes are strictly
classified exactly once only for an authorized, safe-accessible request (DD-11,
DD-14, §5.4; VS0-CF-F16-123..128).

**Validates: REQ-F16-05**

### Property 8: No post-shutdown AuditEvent or Maintenance mutation

For all shutdown sequences, the mutex-serialized admission gate's `Close()`
prevents new timer callbacks and Maintenance triggers from entering, then
drains all currently-admitted callbacks/triggers. Both timer callbacks and
Maintenance triggers must pass `admissionGate.Enter()` before they can acquire
the lifecycle-service mutex, call `AuditAppender.Append`, or commit a
Maintenance mutation; after `Close()` returns, no `Enter()` can succeed and all
previously-admitted callbacks/triggers have exited. Therefore no AuditEvent is
appended and no Maintenance mutation commits after the shutdown boundary (DD-06
shutdown protocol; VS0-CF-F16-77).

**Validates: REQ-F16-09**

### Property 9: First-winner abort-cause immutability

For all qualification reservations, the `abortCause` field on the qualification
reservation is set at most once under the lifecycle-service mutex. The
`releaseReservation` primitive (DD-05) checks `abortCause` before setting it; a
non-nil value is never overwritten. The qualifying owner reads `abortOutcome`
from its held reservation pointer to determine its response deterministically,
regardless of subsequent target transitions or map removal (DD-04 abort-cause
retention; VS0-CF-F16-88, VS0-CF-F16-89).

**Validates: REQ-F16-08**

---

## Error Handling

See §5 (validation/error behavior) and §6.3 (audit matrix/append-failure).

---

## 5. Validation and deterministic error behavior

### 5.1 Request precedence (REQ-F16-05; ADH clause 5) — fixed order

The numbered stages below are the abstract pipeline; §5.1.1 maps them to each
route with the correct per-route ordering (item operations resolve the target
by UID before authorization; collection create extracts references before
authorization; LIST authorizes before enumeration).

1. Closed transport method/path guard (pre-ServeMux; DD-08).
2. Authentication → `AUTH_REQUIRED`/401 on missing/invalid, no audit, no record.
3. Method/media/header-form validation (application/json for POST; exactly one
   Idempotency-Key; qualify/retire: exactly one strong opaque If-Match + zero-byte
   body; query parameter or If-Match on collection create → `MALFORMED_REQUEST`/400;
   non-JSON media → `UNSUPPORTED_MEDIA_TYPE`/415; malformed UID segment →
   `MALFORMED_REQUEST`/400 before safe resolution). Header-form failure precedes
   body classification. Empty/malformed/overlong Idempotency-Key →
   `MALFORMED_REQUEST`/400.
4. Phase-one safe resolution (route-specific ordering):
   - **Item GET/qualify/retire**: validate the path `{uid}` canonical grammar
     (malformed → `MALFORMED_REQUEST`/400 at step 3); resolve that UID to a
     target existence/accessibility predicate; an inaccessible or non-existent
     target returns safe `RESOURCE_NOT_FOUND`/404 +
     `VS0_AUTHORIZATION_SAFE_DENIAL` (audited) **before** derived-scope
     authorization (VS0-CF-F16-17,19,66). Only after the target is safe-resolved
     does the pipeline derive the CloudProvider scope for authorization: the
     target's persisted `spec.cloudProviderParticipationRef.uid` is used to
     look up the `ParticipationCloudProviderUID` via an internal,
     non-authoritative derived-scope index (see DD-16 below). This index is
     not an ExecutionTarget field, not a persisted `scopeRef`, and not an
     authoritative scope record; it is a process-local lookup cache derived
     solely from the immutable `spec.cloudProviderParticipationRef.uid`
     reference and the participation's own `Spec.CloudProviderRef.UID` at
     create time. The lookup discloses no backing-resource existence to the
     caller.
   - **Collection create**: run the bounded body read; if the body exceeds the
     existing limit → `REQUEST_TOO_LARGE`/400 before phase-one extraction,
     authorization, or safe access (VS0-CF-F16-126). Otherwise extract only
     `spec.cloudProviderParticipationRef.uid` and
     `spec.infrastructureStackRef.uid` (see §5.4, DD-11, DD-14). If phase-one
     extraction succeeds, call `BackingAccessProvider.Access` (DD-13.2). Its
     closed disposition governs the next step: participation `SafeDenied`
     returns the registered audited safe 404 before authorization; an
     `AuthorizationDenied` returns the registered authorization outcome; an
     `Allowed` view supplies derived CloudProvider scope for authorization.
     After authorization, stack `SafeDenied` returns the registered audited
     safe 404. No route consumes raw backing existence flags.
5. Authorization (uses the scope derived in step 4): create/retire
   `executiontarget.write`; reads `executiontarget.read`; qualify
   `executiontarget.qualify`. For item routes this occurs only after UID-only
   safe resolution succeeds; for collection create it occurs after phase-one
   extraction and provider safe-access evaluation succeeds
   (VS0-CF-F16-111 proves authorization follows safe resolution).
6. Safe backing access (collection create only): a provider stack `SafeDenied` → safe
   `RESOURCE_NOT_FOUND`/404 + `VS0_AUTHORIZATION_SAFE_DENIAL` (audited);
   never discloses existence (VS0-CF-F16-09,124). For qualify/retire, backing
   viability is validated at step 8a (below) after replay/reservation, using
   the target's persisted backing references and a fresh provider view.
7. Strict classification (single, on retained bytes for create per §5.4): unknown
   field / `adapterAuthorityRef` / SecretRef → `UNKNOWN_FIELD`/400; client status
   write → `AUTHORIZATION_DENIED`/403 + `VS0_STATUS_FIELD_WRITE`; client
   identity/metadata write → `AUTHORIZATION_DENIED`/403 +
   `VS0_SYSTEM_OWNED_FIELD_WRITE`; malformed body → `MALFORMED_REQUEST`/400;
   duplicate top-level member → `DUPLICATE_FIELD`/400.
8. Graph/reference validation (collection create only): missing required field /
   invalid typed-ref shape / non-`synthetic-iaas` targetClass →
   `VALIDATION_FAILED`/422; from the step-4 `Allowed` view:
   ineffective participation (`ParticipationEffectiveActive=false`) →
   `CONFLICT`/409 + `VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE`;
   non-Active stack (`StackPhase≠Active`) →
   `CONFLICT`/409 + `VS0_EXECUTION_TARGET_STACK_UNAVAILABLE`;
   cross-CloudProvider graph (`ParticipationCloudProviderUID ≠
   StackCloudProviderUID`) → `VALIDATION_FAILED`/422 +
   `VS0_EXECUTION_TARGET_SCOPE_MISMATCH`; duplicate name or live duplicate
   tuple → `ALREADY_EXISTS`/409.
8a. Qualify-time backing validation (qualify only): after replay/reservation
   (step 9) and current If-Match comparison (step 10) succeed, and before
   lifecycle state evaluation (step 11), the pipeline obtains a fresh
   `BackingAccessProvider` view using the target's persisted participation and stack
   UIDs. This validates current backing viability: ineffective participation
   → `CONFLICT`/409 + `VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE`
   (VS0-CF-F16-10); non-Active stack → `CONFLICT`/409 +
   `VS0_EXECUTION_TARGET_STACK_UNAVAILABLE` (VS0-CF-F16-11). The `Allowed` view
   also supplies the initial captured fences for the qualification reservation
   (DD-04). This step does not apply to retire (retire releases the tuple
   regardless of backing state).
9. Replay/reservation (create/qualify/retire only): completed replay →
   original successful status/body + fresh correlation, no second audit;
   changed digest → `CONFLICT`/409 + `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`;
   independent reservations across action targets and derived scopes.
10. Current If-Match version comparison (actions): syntactically valid stale →
    `STALE_RESOURCE_VERSION`/412 (after replay/reservation, before backing
    validation for qualify and Maintenance state evaluation for retire).
11. Lifecycle state: Retired → `CONFLICT`/409 + `VS0_TARGET_RETIRED`;
    Maintenance → `CONFLICT`/409 + `VS0_TARGET_MAINTENANCE`; different-key qualify
    while Qualifying → `CONFLICT`/409 + `VS0_TARGET_QUALIFICATION_IN_PROGRESS`.
12. Atomic commit: fence recheck at commit — active current-Maintenance marker
    (checked first per ADH-2026-060) → `CONFLICT`/409 +
    `VS0_TARGET_MAINTENANCE`; otherwise target Retired at commit →
    `CONFLICT`/409 + `VS0_TARGET_RETIRED` (Retire-wins; see DD-04a);
    otherwise a changed captured maintenance epoch →
    `STALE_RESOURCE_VERSION`/412 + `VS0_TARGET_EPOCH_STALE`; a changed
    InfrastructureStack generation remains the independent
    `STALE_RESOURCE_VERSION`/412 + `VS0_TARGET_EPOCH_STALE` case; these
    outcomes are mutually exclusive per commit attempt (ADH-2026-060); stale
    viability → `STALE_RESOURCE_VERSION`/412 +
    `VS0_EXECUTION_TARGET_VIABILITY_STALE`; observer fault
    (duplicate/missing/malformed) → `INTERNAL_ERROR`/500 with reservation
    abort and no publication; required AuditEvent append failure →
    `INTERNAL_ERROR`/500 with no AuditEvent, mutation, or completion and no
    suppressed-outcome disclosure.

### 5.1.1 Five-route stage matrix

The fixed sequence above applies per route with the per-route safe-resolution
ordering stated below. For item operations, UID-only safe target resolution
precedes derived-scope authorization. For collection create, phase-one
extraction and provider safe-access evaluation precede authorization and safe
backing access. For qualify, backing validation (step 8a) occurs after
replay/reservation and current If-Match, using a fresh provider view. For
LIST, authorization precedes accessible-target enumeration. This matrix does not
add a route or outcome; it makes the five existing route contracts executable. A
waiter that wakes after an owner abort re-enters at authentication and repeats
every applicable stage before it can establish a new reservation.

Conformance lists below are non-exhaustive representative subsets; the
canonical ledger (§ Canonical coverage ledger) is the authoritative mapping.

| Route | Applicable stages in order | Skipped stages |
|-------|----------------------------|----------------|
| Collection create | transport guard → authentication → POST media/key grammar (If-Match rejected) → bounded body read (oversized → REQUEST_TOO_LARGE) → phase-one extraction of only `spec.cloudProviderParticipationRef.uid` and `spec.infrastructureStackRef.uid` (unextractable → fail-closed audited 403) → principal-aware BackingAccessProvider lease (inaccessible participation → safe 404 + VS0_AUTHORIZATION_SAFE_DENIAL) → derived-scope `executiontarget.write` authorization → safe backing access (inaccessible stack → safe 404) → single strict classification of retained bytes → graph/reference and uniqueness validation → reservation/replay → audit-before-publication commit | path-UID resolution; action If-Match; lifecycle action state; qualify-time backing validation |
| Collection LIST | transport guard → authentication → `executiontarget.read` authorization → accessible-target enumeration → fresh BackingAccessProvider view per target (outside lifecycle mutex) → effectiveAvailability computation → ascending-UID projection | body decode/classification; backing graph validation; safe item resolution; idempotency; If-Match; lifecycle/fence commit |
| Item GET | transport guard → authentication → malformed-UID validation → UID-only safe target resolution (inaccessible → audited safe 404) → derive CloudProvider scope from derived-scope index (DD-16) → `executiontarget.read` authorization → fresh BackingAccessProvider view (outside lifecycle mutex) → effectiveAvailability computation → safe projection | body decode/classification; idempotency; If-Match; lifecycle/fence commit |
| Qualify action | transport guard → authentication → action media/key/If-Match/zero-byte grammar → malformed-UID validation → UID-only safe target resolution (inaccessible → audited safe 404) → derive CloudProvider scope from derived-scope index (DD-16) → `executiontarget.qualify` authorization → replay/reservation → current If-Match → qualify-time backing validation via BackingAccessProvider (ineffective participation → F16-10; non-Active stack → F16-11; view supplies initial captured fences) → lifecycle state → fenced observer invocation → final paired lease recheck → qualification commit and audit-before-publication | client graph body fields; create uniqueness; safe backing access (step 6) |
| Retire action | transport guard → authentication → action media/key/If-Match/zero-byte grammar → malformed-UID validation → UID-only safe target resolution (inaccessible → audited safe 404) → derive CloudProvider scope from derived-scope index (DD-16) → `executiontarget.write` authorization → replay/reservation → current If-Match → lifecycle state → audit-before-publication retirement commit | client graph body fields; observer/profile evaluation; qualification fences; safe backing access; qualify-time backing validation |

### 5.2 Closed local violation set (REQ-F16-10)

Exactly eight new codes, appearing only in `violations[]`:
`VS0_TARGET_RETIRED`, `VS0_TARGET_MAINTENANCE`,
`VS0_TARGET_QUALIFICATION_IN_PROGRESS`, `VS0_TARGET_EPOCH_STALE`,
`VS0_EXECUTION_TARGET_PARTICIPATION_UNAVAILABLE`,
`VS0_EXECUTION_TARGET_STACK_UNAVAILABLE`,
`VS0_EXECUTION_TARGET_SCOPE_MISMATCH`,
`VS0_EXECUTION_TARGET_VIABILITY_STALE`. Top-level Problem codes remain
FEATURE-0012-owned; HTTP status is transport metadata, never a Problem `code`;
no `429`/quota code exists. `VS0_STATUS_FIELD_WRITE`,
`VS0_SYSTEM_OWNED_FIELD_WRITE`, `VS0_AUTHORIZATION_SAFE_DENIAL`, and
`VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH` are inherited violation codes reused, not
introduced.

### 5.3 Deterministic single-writer commit

All commits serialize through the lifecycle service (DD-02). The registered
writer conflict for VS0-WRITER-006 is `STALE_RESOURCE_VERSION`. The commit step
is all-or-nothing: audit append precedes publication; a failed append leaves the
committed target, ETag, and record links unchanged.

### 5.4 Collection-create phase-one and single-strict-classification precedence (ADH-2026-063/065)

This subsection fixes, without altering §5.1, the exact create-only precedence
realized by DD-11. The six exact cases are:

| Case | Condition | Outcome |
|------|-----------|---------|
| VS0-CF-F16-126 | Oversized collection-create body exceeding the bounded single-body-read limit | `REQUEST_TOO_LARGE`/400 (as VS0-CF-F16-53) before phase-one extraction, authorization, or safe access; no phase-one extraction, classification, digest, reservation, mutation, AuditEvent, or completion |
| VS0-CF-F16-128 | Bounded read succeeds but phase one cannot extract exactly one syntactically usable UID at both required locations | Fail-closed existing audited `AUTHORIZATION_DENIED`/403 (as VS0-CF-F16-08) before derived-scope authorization and safe backing access; no body/backing disclosure; no strict classification, canonicalization, digest, reservation, observer invocation, mutation, publication, or completion; no added violation or top-level Problem code |
| VS0-CF-F16-123 | Authenticated caller lacking `executiontarget.write` over safe-accessible backing; body malformed JSON or duplicate top-level member | Existing audited `AUTHORIZATION_DENIED`/403 (as VS0-CF-F16-08); phase one extracts only the two allowed UIDs and never classifies, canonicalizes, digests, or reserves the body |
| VS0-CF-F16-124 | Inaccessible backing reference; body malformed JSON or duplicate top-level member | Safe `RESOURCE_NOT_FOUND`/404 + `VS0_AUTHORIZATION_SAFE_DENIAL` (as VS0-CF-F16-09); no body classification, digest, reservation, mutation, or completion |
| VS0-CF-F16-125 | Authorized, safe-accessible request; retained body malformed JSON only | Single strict classification returns `MALFORMED_REQUEST`/400 (as VS0-CF-F16-51); exact non-family case, distinct from VS0-CF-F16-127; no mutation, AuditEvent, digest, reservation, or completion |
| VS0-CF-F16-127 | Authorized, safe-accessible request; retained body duplicate top-level member only | Single strict classification returns `DUPLICATE_FIELD`/400 (as VS0-CF-F16-52); exact non-family case, distinct from VS0-CF-F16-125; no mutation, AuditEvent, digest, reservation, or completion |

Action `If-Match` header-form precedence (§5.1 step 10) is preserved. VS0-CF-F16-125,
VS0-CF-F16-127, and VS0-CF-F16-128 map to REQ-F16-05; VS0-CF-F16-123, VS0-CF-F16-124,
and VS0-CF-F16-126 are unchanged precedence rows.

---

## 6. Security, privacy, observability, compatibility, and versioning

### 6.1 Security

- Authentication precedes all F0016 effects; missing/invalid → 401 with no audit
  or record (VS0-CF-F16-07,61..64,93..97).
- Grants: `executiontarget.write` (create/retire), `executiontarget.read`
  (LIST/GET), `executiontarget.qualify` (qualify), each bound to the derived
  CloudProvider scope of the target UID for item operations.
- For collection create, phase one fails closed with audited
  `AUTHORIZATION_DENIED`/403 when either required UID reference is unextractable,
  before derived-scope authorization and safe backing access, disclosing no body
  or backing detail (DD-11, §5.4; VS0-CF-F16-128).
- Client status/identity writes are audited 403 with the inherited
  `VS0_STATUS_FIELD_WRITE` / `VS0_SYSTEM_OWNED_FIELD_WRITE` violations; clients
  never write status (VS0-WRITER-005/006).
- Network-exposed routes are authenticated and authorized by design; there is no
  unauthenticated mutating path. The pre-ServeMux guard is transport-only and
  performs no authentication decision, so it never exposes resource data.

### 6.2 Privacy / redaction

- ExecutionTarget is Provider-confidential; NormalizedTargetFactSet and
  TargetQualificationResult are internal-only and never projected.
- Safe denial (404 + `VS0_AUTHORIZATION_SAFE_DENIAL`) never confirms existence of
  an inaccessible target or backing reference.
- No secret value on any path; no SecretRef extension is active; no credentials
  exist. Replay discloses only the original successful status/body with fresh
  correlation; denial replays disclose no stored result.

### 6.3 Observability

- Reuse the inherited FEATURE-0013 AuditEvent as the durable audit signal. The
  audit matrix audits: successful create, completed qualification, retirement,
  maintenance enter/clear, authenticated authorization denial, and safe denial.
  A required request-triggered append failure → `INTERNAL_ERROR`/500 with no
  AuditEvent, mutation, or completion (VS0-CF-F16-39,40,82,99..104,112).
- Expiry is the sole background exception: no request caller, no 500, one expiry
  AuditEvent once its append succeeds, retrying once per injected-clock second
  (VS0-CF-F16-41,79).
- Structured request/correlation fields follow the FEATURE-0015 handler
  conventions and the observability-and-audit baseline; no secret is logged.

### 6.4 Compatibility and versioning

- Reuses FEATURE-0012 Problem Details, ObjectMeta, TypedRef, Condition, media,
  and optimistic-concurrency conventions by reference; defines no parallel
  envelope. API group/version is `execution.sovrunn.io/v1alpha1`.
- Consumes FEATURE-0015 `CloudProviderParticipation` and `InfrastructureStack`
  read-only; modifies no FEATURE-0015 document, schema, writer, state, or
  conformance case.
- Per ADH-2026-064, FEATURE-0016 remains the sole owner and activator of the
  ExecutionTarget contract; its being a declared prerequisite for a dependent
  feature (e.g. FEATURE-0022 ServiceRegion) does not transfer activation to that
  feature and does not require FEATURE-0016 to introduce or activate any
  dependent-feature field.
- ADH-2026-058 atomically replaced the placeholder VS0-SCHEMA-015..017 and
  VS0-STATE-004; there is no alias, migration path, or legacy compatibility for
  `status.availability`, `Draining`, or publicly projected `Qualifying`. No
  migration impact: all state is new and in-process.

### 6.5 Operational

- All state is in-process. Restart clears F0016-owned targets, records,
  idempotency reservations, and pending work; it never deletes inherited
  FEATURE-0013 AuditEvents and makes no external observer call (VS0-CF-F16-46).
- Shutdown ordering: (1) close the mutex-serialized admission gate (prevents
  new timer callbacks and Maintenance triggers from entering, then drains all
  currently-admitted — see DD-06 shutdown protocol); (2) under the
  lifecycle-service mutex: cancel all outstanding timers; abort all in-flight
  qualification and create/retire reservations via `releaseReservation`
  (DD-05), closing their `done` channels and waking waiters; invalidate all
  remaining expiry entries; (3) release the mutex; (4) stop HTTP serving. The
  admission gate close guarantees no `AuditAppender.Append` call and no
  Maintenance mutation occurs after step 1 completes (VS0-CF-F16-77). The last
  committed target projection, ETag, and record links remain unchanged through
  shutdown.

---

## Testing Strategy

See §7 (test and conformance strategy) below.

---

## 7. Test and conformance strategy

### 7.1 Placement and reuse

Feature-local conformance cases follow the live convention in
`tests/conformance/` (`feature_0016_test.go`, `feature_0016_fixtures.go`),
reusing the FEATURE-0015 test harness/fixtures patterns. Unit and property tests
live beside the implementation in `internal/executiontarget/` and
`internal/api/`, and route/guard tests in `internal/server/`
(`executiontarget_guard_test.go`, `server_test.go`), matching FEATURE-0015's
`*_test.go` co-location.

### 7.2 Conformance coverage obligation

Every one of the 128 FEATURE-0016-local cases (VS0-CF-F16-01..128) is realized as
a test asserting its exact registered `expectedState`, `expectedError`,
`expectedViolation` (where registered), and `expectedSideEffects`. Equivalence
families (VS0-CF-F16-106,112,114,115,118) must execute every named member; the
named finite equivalence families for the create-precedence rows
(VS0-CF-F16-123,124: malformed JSON or duplicate top-level member) must execute
every named member. Determinism is guaranteed by the injected clock and
target-bound fixtures; no test relies on wall-clock elapse or network. No
shared/downstream case (`VS0-CF-F10`, `VS0-CF-X03`, `VS0-CF-HP01`) is counted as
F0016-local proof.

### 7.3 Test categories

- Create-field boundary and scope derivation (01,02,04,54..57,84,91,92).
- Safe access before graph validation (09..13,17,19,65,66,98).
- Collection-create phase-one/strict-classification precedence (123..128).
- Route/guard transport outcomes (44,45,48,105..110,113..118) and precedence
  (05,06,07,08,24,25,58..60,70..74,85,86).
- Idempotency replay/reuse/isolation/eviction/panic/waiter (14,16,31..35,46,47,
  49,50,76..78,81,82,87..89,99).
- Observer/facts/conclusions (20..23,67..69,90) and lifecycle/availability
  (26..30,36,38,39,41..43,75,79..83,88,89,100,101,119..122).
- Audit matrix and append-failure (39,40,82,99..104,112) and future boundary (46).

### 7.4 Race and concurrency test obligations

The following focused race tests are required in addition to the conformance
suite. Each exercises a specific interleaving that the design's synchronization
protocol must survive under `go test -race`:

1. **Expiry timer replacement race.** A qualification commits and schedules an
   expiry timer at T₁; a second qualification commits before T₁ fires,
   replacing the pending entry and stopping the first timer. Prove: exactly
   one expiry AuditEvent is produced (for the second FactSet), and the first
   timer callback is a deduplication no-op.

2. **Expiry audit-debt survives superseding state.** A timer fires
   (entry → `auditDebt`) and the first append fails. Before the retry fires,
   a Maintenance entry commits. Prove: the `auditDebt` entry is not removed
   by Maintenance, the retry fires, append succeeds, and exactly one expiry
   AuditEvent is produced for the original (now-superseded) FactSet.

3. **Shutdown versus callback entry/append.** A timer fires and its callback
   passes `admissionGate.Enter()`. Concurrently, shutdown calls
   `admissionGate.Close()`. Prove: either the callback completes its append
   before `Close` returns (valid), or the callback finds the entry invalidated
   and produces no append (valid). No append occurs after `Close` returns.

4. **Shutdown versus callback that has not yet entered.** Shutdown calls
   `admissionGate.Close()`. A timer fires after `closed=true`. Prove: the
   callback is rejected at `Enter()` and produces no append.

5. **Retire versus in-flight qualification (DD-04a).** A qualify request
   establishes a reservation and begins observer invocation. Concurrently, a
   retire request commits. Prove: the qualification owner observes
   `abortCause = VS0_TARGET_RETIRED`, returns `CONFLICT`/409, produces no
   AuditEvent, no completion, and the retire's own AuditEvent is produced.

6. **Maintenance entry versus in-flight qualification (DD-12).** Same as above
   but with Maintenance entry. Prove: `abortCause = VS0_TARGET_MAINTENANCE`,
   exactly the qualification's linked idempotency reservation is removed
   (not others), and only that qualification's waiters are woken.

7. **Same-key waiter wake-up and re-entry.** Owner qualifies and commits.
   Two waiters are blocked on the `done` channel. Prove: both wake, both
   re-enter the pipeline, both receive the stored replay response from the
   now-Completed idempotency reservation.

8. **Same-key waiter wake-up on owner abort.** Owner's audit append fails;
   reservation is removed. Two waiters wake. Prove: both re-enter, find no
   reservation, one becomes the new owner, the other waits again.

9. **Paired backing lease under concurrent FEATURE-0015 writes.** Hold an
   Allowed provider lease while a concurrent backing writer attempts mutation.
   Prove: the provider view is paired and immutable, the writer cannot interleave
   through final qualification publication, and no lock inversion/deadlock
   occurs.

10. **Principal-aware result containment.** Exercise SafeDenied,
    AuthorizationDenied, and Allowed provider outcomes. Prove: denied callers
    receive only registered non-disclosing results; raw backing existence never
    reaches a handler or projection.

11. **Fence and projection boundaries.** Change participation generation only
    during qualification and prove it is not a stale fence; change Stack
    generation or viability and prove the registered stale outcome. Verify GET
    and LIST use a fresh Allowed view for `effectiveAvailability` without an
    ETag mutation.

### 7.5 Fuzz and property tests for phase-one scanner

The phase-one UID extraction scanner (DD-14) is security-critical: it
determines whether the pipeline proceeds to authorization or fails closed. The
following fuzz/property tests are required:

1. **Fuzz: arbitrary bytes.** `testing/F` fuzz test feeds arbitrary `[]byte`
   to the scanner. Prove: the scanner never panics, never produces more than
   two UIDs, and terminates in bounded time proportional to input length.

2. **Fuzz: valid JSON with injected paths.** Fuzz test generates syntactically
   valid JSON with random nesting, duplicate keys, and the two target paths at
   random positions. Prove: the scanner correctly counts occurrences and
   extracts the UID when exactly one valid occurrence exists at each path.

3. **Property: duplicate key non-extraction.** For any JSON with ≥ 2
   occurrences of either complete target path—regardless of value types—the
   scanner returns unextractable for that path. This includes a string/non-string
   pair, duplicate intermediate keys, and escaped-equivalent key spellings.

4. **Property: non-string value non-extraction.** For any JSON where the
   target path resolves to a non-string (number, boolean, null, object, array),
   the scanner returns unextractable.

5. **Property: malformed JSON handling.** For any byte sequence that is not
   valid JSON, the scanner either extracts valid UIDs (if present at the
   target paths despite surrounding malformation) or returns unextractable.
   It never panics or hangs.

6. **Property: canonical UID grammar enforcement.** For any extracted string,
   the grammar check correctly accepts/rejects per the FEATURE-0012 canonical
   UID grammar.

7. **Property: malformed boundaries and termination.** Invalid escapes,
   malformed surrogate pairs, invalid UTF-8, unbalanced nesting, and every skip
   path terminate linearly and either extract the same two unambiguous UIDs or
   fail closed.

### 7.6 Formal proof references (owned by architecture; referenced only)

`docs/formal/feature-0016/ExecutionTargetLifecycle.tla` (Retired terminality,
availability truth table, Maintenance/Retirement priority, non-public Qualifying)
and `docs/formal/feature-0016/ExecutionTargetFence.tla` (audit-before-conclusion,
stale/cancelled/failed non-mutation) are consumed as design constraints, not
re-authored.

### 7.7 Verification commands (tasks-stage execution)

`make fmt`, `make test`, `make vet`, `go test -race ./...`,
`go test ./internal/cloudmodel/...`,
`make feature-contract-check FEATURE=FEATURE-0016`,
`make vs000-contract-check`, and `make phase2r-drift-check` are component
checks. The authoritative final gate is
`make ff-feature-gate FEATURE=FEATURE-0016`; it must pass after implementation.
Design authors no runtime; these commands run when implementation exists.

---

## 8. Implementation classification ledger

### 8.1 IMPLEMENT (produces tasks / code later)

| Element | Component (§3) | Realizes |
|---------|----------------|----------|
| ExecutionTarget in-memory type + store + tuple/name index + derived-scope index | model/, store.go | Create-time persistence and scope-derived, name/tuple-unique storage of ExecutionTarget; non-authoritative derived-scope index (DD-16) |
| ExecutionTargetLifecycleService (sole committer) | lifecycle.go | Sole-committer lifecycle state management and expiry/maintenance scheduling |
| Qualifying reservation + captured-fence commit | lifecycle.go, qualify.go | Sole-committer atomic, fence-checked qualification commit with first-winner abort-cause retention, final paired-lease recheck, and Maintenance-first predicate ordering (DD-04, DD-04a, DD-13.5) |
| Maintenance-entry link-clearing + target-local abort | lifecycle.go, scheduler.go | Sole-committer Maintenance-entry semantics with first-winner abort-cause check and releaseReservation cleanup per DD-05/DD-12 |
| Synthetic observer + fixtures | observer.go | Synthetic fact observation feeding qualification (DD-13.1 Observer interface) |
| Qualification profile evaluation | qualify.go | Synthetic fact observation feeding qualification |
| NormalizedTargetFactSet + TargetQualificationResult records | model/ | Synthetic fact observation and the sole-committer qualification commit it feeds |
| effectiveAvailability projection | projection.go, lifecycle.go | Sole-committer lifecycle state management and its response-only availability projection |
| Idempotency reservation/waiter table with unified broadcast signal | idempotency.go | F0016-owned idempotency handling for create/qualify/retire; every InFlight reservation carries a done channel; unified releaseReservation cleanup primitive (DD-05) |
| Injected-clock expiry worker + per-expiry bookkeeping + retry | scheduler.go, clock.go | Expiry/maintenance scheduling with (target UID, factSetRef)-keyed entries and timer generation tracking per DD-06 and DD-13.4; pending/auditDebt/audited entry lifecycle; mutex-serialized admission gate shutdown covering timer callbacks and Maintenance triggers |
| Target-bound maintenance-trigger intake + current-Maintenance marker | scheduler.go, lifecycle.go | Expiry/maintenance scheduling and sole-committer lifecycle state management; admission gate governs trigger entry |
| Shutdown ordering | scheduler.go, lifecycle.go | Background gate close → one lifecycle-mutex request-admission closure, timer cancellation, reservation detachment, and entry invalidation → post-unlock waiter signals → HTTP stop (DD-06) |
| Phase-one UID extractor (collection create) | `internal/api/executiontarget_collection.go` | Non-classifying bounded extraction per DD-14; deterministic for F16-123..128 |
| Five handlers + request pipeline (incl. collection-create phase-one precedence) | `internal/api/executiontarget_*.go` | The closed HTTP surface, its fixed ordered request pipeline, and DD-11 precedence |
| Five ServeMux registrations + pre-ServeMux guard | `internal/server/routes.go`, `internal/server/executiontarget_guard.go` | The closed HTTP surface |
| Safe projection + closed representation | projection.go | The closed HTTP surface's response representation |
| Audit assembly (redacted, matrix) | audit.go | Sole-committer lifecycle state management and expiry/maintenance scheduling audit effects |
| Coherent principal-aware backing bridge | `internal/executiontarget/backing_access.go`, `backing_access_test.go`; `internal/cloudmodel/store.go`, `store_test.go` | ADH-2026-066 paired lease, caller-relative safe results, registered-fence final commit, and FEATURE-0015 compatibility regression proof |
| Eight new violation constants | `internal/executiontarget/violations.go` | F0016-owned `violations[]` constants; `internal/apiproblem` Problem enum/map is reused unchanged |
| Fresh backing reads for GET/LIST effectiveAvailability | `internal/api/executiontarget_item.go`, `internal/api/executiontarget_collection.go` | BackingAccessProvider view per target outside lifecycle mutex; viability-only changes reflected without ETag mutation (DD-15; VS0-CF-F16-119,120) |
| F0016-local conformance suite (128 cases) | `tests/conformance/feature_0016_test.go` | Full coverage of the feature's design elements above |
| Race tests (§7.4) | `internal/executiontarget/*_test.go` | Focused interleaving tests: expiry replacement, audit-debt survival, shutdown vs. callback, Retire/Maintenance vs. qualification, same-key waiters |
| Fuzz/property tests for phase-one scanner (§7.5) | `internal/api/executiontarget_collection_test.go` | Fuzz: arbitrary bytes and valid-JSON injection; Properties: duplicate-key, non-string, malformed JSON, UID grammar |

### 8.2 CONTRACT_ONLY/NO_TASK (reused by reference; no new task)

| Element | Owner | Disposition |
|---------|-------|-------------|
| Problem Details + HTTP/URN mapping | FEATURE-0012 / `internal/apiproblem` | Reused; F0016 adds no top-level code |
| ObjectMeta / TypedRef / Condition | FEATURE-0012 / `internal/apimeta`,`apiref`,`apicond` | Reused as-is |
| Validation pipeline stages, bounded body read, media, If-Match/ETag | FEATURE-0012 / `internal/apivalid` | Reused; F0016 wires order, adds no stage semantics |
| AuditEvent atomic append-before-publication | FEATURE-0013 / `internal/decision` | Reused; not redefined |
| Server-resolved authorization, safe denial, digest conventions | FEATURE-0015 / `internal/cloudmodel` | Reused read-only. The ADH-2026-066 private lease extension is F0016 IMPLEMENT work listed above; FEATURE-0015 remains unchanged as a product feature. |
| Inherited violation codes (`VS0_STATUS_FIELD_WRITE`, `VS0_SYSTEM_OWNED_FIELD_WRITE`, `VS0_AUTHORIZATION_SAFE_DENIAL`, `VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH`) | prior features | Reused; not introduced |
| Formal models (lifecycle, fence) | architecture / `docs/formal/feature-0016` | Referenced as constraints |

### 8.3 EXCLUDED (must never be designed or built by F0016)

Real adapter/provider-native type/credential/SecretRef/external call/external
persistence; placement, provisioning, plugin execution, realization, plugin
taxonomy, Crossplane dependency; customer target projection; maintenance-notice
resource; IAM resource; the explicit non-goals listed in §10.1; persisted
`status.availability`, `Draining`, public `Qualifying`, public `conditions`,
`scopeRef`, `spec.participationRef`, `adapterAuthorityRef`; FEATURE-0017
PolicyEvaluationRequest/Result/PolicyEngineAdapter; FEATURE-0019
SovereigntyProfile/SovereigntyFactSet/EvidenceRecord; FEATURE-0022
ServiceRegion / `spec.executionTargetRefs` / `ServiceRegion.status.availability`;
FEATURE-0023 DecisionProfiles/ServicePlacement/PlacementDecision/customer target
projection; FEATURE-0024 PluginExecution/ServiceDeploymentPlan/real realization.

---

## 9. Requirement/decision/risk traceability

### 9.1 Architecture Traceability

- Consumed decisions/handoffs: DEC-0036, DEC-0042, DEC-0057; ADH-2026-025,
  ADH-2026-040, ADH-2026-042, ADH-2026-045, ADH-2026-058 (foundational contract
  closure), ADH-2026-060, ADH-2026-061, ADH-2026-063, ADH-2026-064,
  ADH-2026-065, ADH-2026-066 (coherent principal-aware backing bridge).
- Schema IDs: VS0-SCHEMA-015 (ExecutionTarget), VS0-SCHEMA-016
  (NormalizedTargetFactSet), VS0-SCHEMA-017 (TargetQualificationResult) —
  referenced, not redefined.
- Writer ID: VS0-WRITER-006 (`ExecutionTargetLifecycleService`; conflict
  `STALE_RESOURCE_VERSION`).
- State ID: VS0-STATE-004 (Active/Retired × Unqualified/Qualified/Rejected/
  Indeterminate; Qualifying is an in-flight reservation only).
- Error codes: top-level Problem codes owned by FEATURE-0012; eight new
  `violations[]` codes per §5.2.
- Conformance: VS0-CF-F16-01..128 (all F0016-local).

### 9.2 Risk-relevant invariants and enforcement

| Invariant | Design enforcement |
|-----------|--------------------|
| Clients never write status; single sole committer | DD-02; handlers never mutate status; VS0-CF-F16-03,04,91,92 |
| No external effect/credential/persistence | DD-03; observer external-effect-free; §8.3; VS0-CF-F16-46 |
| Audit-before-publication; no partial mutation on append failure | DD-04, DD-13.5, §5.3; VS0-CF-F16-39,40,82,99..104 |
| Safe denial never discloses inaccessible resources | §5.1 step 4 (item UID-only safe resolution; collection-create principal-aware backing lease), §5.4, §6.2; VS0-CF-F16-09,17,19,66,124,128 |
| Closed five-route surface; transport guard has no side effect | DD-08; VS0-CF-F16-44,113..118 |
| Stale fences cannot publish a conclusion or replayable completion; final paired lease keeps backing values stable through qualification publication | DD-04, DD-04a, DD-12, DD-13.2, DD-13.5, §5.1 step 12; VS0-CF-F16-29,30,75,81,82,88 |
| Collection-create phase-one/single-strict-classification precedence | DD-11, DD-14, §5.4; VS0-CF-F16-123..128 |
| No post-shutdown AuditEvent or Maintenance mutation; mutex-serialized admission gate close drains all admitted timer callbacks and Maintenance triggers before shutdown proceeds | DD-06 shutdown protocol (mutex-serialized admissionGate Enter/Close covering timer callbacks and Maintenance triggers + entry invalidation + reservation abort); VS0-CF-F16-77 |
| First-winner abort-cause is immutable once set; qualifying owner response is deterministic | DD-04 abort-cause retention; DD-04a, DD-12; VS0-CF-F16-88,89 |

### 9.3 Conformance Mapping (behavior → design disposition → conformance IDs)

The canonical coverage ledger (below) is the single source of REQ-F16-*/AC-F16-*
identifier mappings; this table restates the same coverage in plain-language
behavior terms, pointing only at conformance IDs.

| Behavior | Design disposition | Conformance IDs |
|----|--------------------|-----------------|
| Create-time persistence and field acceptance | model/,store.go,lifecycle.go create path (§4.1,§5.1) | VS0-CF-F16-01,02,57,84 |
| Safe access ordering before graph validation | safe resolver and BackingAccessProvider before graph validation (§5.1 steps 4/6/8/8a; DD-13.2) | VS0-CF-F16-09..13,65,98 |
| Collection-create phase-one/strict-classification precedence | bounded read → phase-one extraction → principal-aware backing lease → authorization/safe access → single strict classification (DD-11,§5.4) | VS0-CF-F16-123..128 |
| Qualify-time backing validation | Fresh BackingAccessProvider view after replay/reservation and If-Match; ineffective/non-Active returns registered outcomes (§5.1 step 8a; DD-13.2) | VS0-CF-F16-10,11,29,30 |
| Closed five-route HTTP surface and transport guard | five registrations + pre-ServeMux guard (DD-08,§4.5) | VS0-CF-F16-44,61..64,66,93..96,105..118 |
| Header-form request grammar validation | header-form validation before body classification (§5.1 step 3) | VS0-CF-F16-05,06,24,25,58..60,70..74,85,86 |
| Idempotency replay and reservation behavior | idempotency table (DD-05,§5.1 step 9) | VS0-CF-F16-14,16,31..35,46,47,49,50,76..78,81,82,86..89,99 |
| Synthetic observer fact behavior | observer + qualification engine (DD-03,DD-04,§4.2) | VS0-CF-F16-20..23,67..69,90 |
| Sole-committer lifecycle and effectiveAvailability projection (incl. fresh backing reads for GET/LIST) | lifecycle service + effectiveAvailability (DD-02,DD-15,§4.6) | VS0-CF-F16-01,20..22,26..30,36,39,41..43,75,77,79..83,87,88,99..101,119..122 |
| Retire/Maintenance precedence over in-flight qualification | active-Maintenance marker checked first (ADH-2026-060); first-winner abort-cause retention (DD-04); Retire-wins commit path (DD-04a); Maintenance-entry target-local abort (DD-12); F16-75 is the inactive-marker changed-maintenance-epoch `STALE_RESOURCE_VERSION`/412 case, F16-89 is the active-Maintenance-entry-wins `CONFLICT`/409 case; F16-88 is the Retire-wins `CONFLICT`/409 case | VS0-CF-F16-26..28,75,80,83,88,89 |
| Expiry audit-debt lifecycle and shutdown ordering | (target UID, factSetRef)-keyed pending/auditDebt/audited entry states + timer generation tracking + mutex-serialized admission gate shutdown covering timer callbacks and Maintenance triggers (DD-06,§6.5) | VS0-CF-F16-41,77,79 |
| Audit matrix and append-failure behavior | audit matrix + append-failure (§6.3) | VS0-CF-F16-39,40,82,99..104,112 |
| Closed local violation set | eight closed violations (§5.2) | VS0-CF-F16-10..12,26..30,38,75 |
| Absence of external effect; restart behavior | no external effect; restart clears in-memory (§6.5,§8.3) | VS0-CF-F16-46 |

---

## Canonical coverage ledger

Every approved requirement and every approved acceptance criterion from
`requirements.md` appears below exactly once with its design disposition. No ID
is created, renumbered, merged, split, reinterpreted, or omitted. REQ and AC
meanings are preserved exactly from `requirements.md`.

### Requirements

| ID | Design disposition | Classification |
|----|--------------------|----------------|
| REQ-F16-01 | ExecutionTarget in-memory type, store, closed create-field decode, participation-derived scope, server-assigned identity/status; create persists Active/Unqualified at epoch 0 (DD-01,DD-02; §3.1,§4.1,§5.1) | IMPLEMENT |
| REQ-F16-02 | Safe resolver before graph validation; effective-Active participation + Active stack same-scope check; scope-unique name retained after retirement; at-most-one live tuple index (§3.1 store.go,§5.1 steps 4/6/8) | IMPLEMENT |
| REQ-F16-03 | Five ServeMux registrations + fixed pre-ServeMux transport guard (DD-08; §3.1,§4.5) | IMPLEMENT |
| REQ-F16-04 | Request grammar (media/Idempotency-Key/If-Match/zero-byte), ETag placement, ignored Accept, closed safe projection, status codes (§4.1,§4.5,§5.1 step 3) | IMPLEMENT |
| REQ-F16-05 | Explicit ordered request pipeline with short-circuit and fixed precedence, including the collection-create bounded-read → phase-one extraction (fail-closed audited 403 when unextractable) → authorization/safe access → single strict classification precedence (DD-09,DD-11,DD-14; §5.1,§5.4) | IMPLEMENT |
| REQ-F16-06 | F0016-owned idempotency namespace/digest/reservation/waiter/retention/eviction; phase one never digests or reserves; GET/LIST never touch idempotency (DD-05; §3.1,§5.1 step 9) | IMPLEMENT |
| REQ-F16-07 | Synthetic observer + closed qualification profile + four named facts + missing-fixture/timeout outcomes + fixed provenance + observer-fault faulting (DD-03,DD-04; §3.1,§4.2,§4.4) | IMPLEMENT |
| REQ-F16-08 | Sole-committer lifecycle service; lifecycle/qualification axes; non-projected Qualifying; response-only effectiveAvailability truth table; atomic audited conclusion; Maintenance-entry link-clearing (DD-12); Retire-wins abort (DD-04a); stale-fence abort (DD-02,DD-04; §4.6,§5.3) | IMPLEMENT |
| REQ-F16-09 | Injected-clock expiry with retry and per-expiry bookkeeping; target-bound maintenance trigger with unconditional link-clearing and target-local abort (DD-12); retire sole exit; audit matrix; shutdown ordering (DD-06,DD-07; §3.1,§6.3,§6.5) | IMPLEMENT |
| REQ-F16-10 | Eight new `violations[]` codes only; top-level Problems remain FEATURE-0012 (§5.2) | IMPLEMENT (violation surface); top-level codes CONTRACT_ONLY (reused) |
| REQ-F16-11 | No external effect/adapter/credential/Crossplane; restart clears F0016-owned in-memory state, preserves inherited AuditEvents; ownership not transferred to any dependent feature (ADH-2026-064; §6.4,§6.5,§8.3) | IMPLEMENT (restart/in-memory boundary); external-effect absence CONTRACT_ONLY |

### Acceptance criteria (immutable canonical reference and design-disposition mapping)

The canonical AC ledger is the single immutable source of acceptance-criterion
definitions; it lives in `requirements.md` §4 and is referenced, never
redefined, by this design. The table below maps each canonical AC to its design
disposition and conformance IDs without restating or paraphrasing the criterion
text. The "AC" column cites the immutable identifier; the "Design disposition"
column states where in this design the criterion is realized; the
"Conformance ID" column copies the conformance IDs exactly from the
canonical ledger. For the verbatim criterion text, refer to `requirements.md`
§4 Canonical acceptance ledger.

| AC | Design disposition (where realized) | Conformance ID | Classification |
|----|--------------------------------------|----------------|----------------|
| AC-F16-01 | Create path: model/, store.go, lifecycle.go (§4.1, §5.1 steps 4/7/8, §5.4) | VS0-CF-F16-01,02,57,84 | IMPLEMENT |
| AC-F16-02 | Safe resolution and graph validation: §5.1 steps 4/6/8/8a; DD-13.2 backing lease; §5.4 | VS0-CF-F16-09..13,65,98 | IMPLEMENT |
| AC-F16-03 | Five registrations + pre-ServeMux guard: DD-08, §4.5, §3.1 | VS0-CF-F16-44,61..64,66,93..96,105..118 | IMPLEMENT |
| AC-F16-04 | Header-form validation: §5.1 step 3; action grammar; collection-create precedence §5.4 | VS0-CF-F16-05,06,24,25,58..60,70..74,85..86 | IMPLEMENT |
| AC-F16-05 | Idempotency table: DD-05, §5.1 step 9 | VS0-CF-F16-14,16,31..35,46..47,49..50,76..78,81..82,86..89,99 | IMPLEMENT |
| AC-F16-06 | Synthetic observer: DD-03, DD-13.1, qualify.go, observer.go | VS0-CF-F16-20..23,67..69,90 | IMPLEMENT |
| AC-F16-07 | Lifecycle service + effectiveAvailability: DD-02, DD-04, §4.6, lifecycle.go, projection.go | VS0-CF-F16-01,20..22,26..30,36,39,41..43,75,77,79..83,87..88,99..101,119..122 | IMPLEMENT |
| AC-F16-08 | Retire-wins (DD-04a) + Maintenance-entry abort (DD-12); first-winner rule (DD-04); §5.1 steps 11/12 | VS0-CF-F16-26..28,75,80,83,88..89 | IMPLEMENT |
| AC-F16-09 | Expiry scheduler: DD-06 (timer, retry, shutdown protocol); §6.5 | VS0-CF-F16-41,77,79 | IMPLEMENT |
| AC-F16-10 | Audit assembly: audit.go, §6.3; AuditAppender interface DD-13.3 | VS0-CF-F16-39,40,82,99..104,112 | IMPLEMENT |
| AC-F16-11 | Violation constants: §5.2, violations.go | VS0-CF-F16-10..12,26..30,38,75 | IMPLEMENT |
| AC-F16-12 | No external effect: §6.5, §8.3; restart clears F0016 state only | VS0-CF-F16-46 | IMPLEMENT (restart proof); external-effect absence CONTRACT_ONLY |

### Supplemental conformance evidence (post-handoff ADH-2026-063/065 coverage)

The following table maps the six collection-create precedence cases
(VS0-CF-F16-123..128) to the acceptance criteria they additionally prove. These
are supplemental evidence rows introduced by ADH-2026-063 and ADH-2026-065;
they do not alter the canonical AC conformance IDs above.

| Canonical criterion | Supplemental conformance IDs | Rationale |
|----|------------------------------|-----------|
| Canonical safe-access acceptance proof | VS0-CF-F16-124,128 | Create-precedence safe-denial (F16-124) and fail-closed denial (F16-128) additionally prove safe reference access precedes graph validation for collection create; the canonical acceptance-ID mapping remains in the ledger below. |
| Canonical request-precedence acceptance proof | VS0-CF-F16-123,124,125,126,127,128 | The full six create-precedence cases additionally prove header-form/body-classification ordering for the collection-create route; the canonical acceptance-ID mapping remains in the ledger below. |

---

## 10. Non-goals, absence ledger, and unresolved report

### 10.1 Non-goals and adjacent-feature exclusions

FEATURE-0016 does not design or build any element listed in §8.3. It introduces
no real adapter, provider-native type, raw credential, external call, external
persistence, placement, execution, plugin, customer target API,
maintenance-notice resource, or IAM resource. FEATURE-0017/0019/0022/0023/0024
concepts appear only as enforceable exclusions, never as design mechanisms.

The following are the explicit non-goals for this feature:

- Non-goal: FEATURE-0016 must not introduce, activate, or design `ResourcePool`.
- Non-goal: FEATURE-0016 must not introduce, activate, or design `ProviderCapability`.
- Non-goal: FEATURE-0016 must not introduce, activate, or design a generic `Provider`.
- Non-goal: FEATURE-0016 must not introduce or activate any ExecutionTarget subtype as an active platform concept.

### 10.2 Absence ledger (deliberately absent by contract)

- Persisted `status.availability`; `Draining`; publicly projected `Qualifying`.
- Public `conditions` member on ExecutionTarget.
- `scopeRef` (scope derived solely from participation); `spec.participationRef`
  (renamed to `spec.cloudProviderParticipationRef`).
- `adapterAuthorityRef`; any `SecretRef` extension.
- `viabilityFingerprint` on TargetQualificationResult (FactSet-only); a second
  persisted FactSet version field (only `factVersion=v1` persists).
- `Location` header; conditional GET; `Accept` negotiation; any query parameter.
- PATCH/PUT/DELETE/HEAD public routes; watch/filter/pagination/fact/result/
  adapter-selection routes.
- A second body classification pass, or any phase-one body classification,
  canonicalization, digest, or reservation on collection create (single strict
  classification only, after authorization and safe access; DD-11).
- `429`/quota HTTP code; any top-level Problem code introduced by F0016.
- Any alias, migration path, or legacy compatibility for the above.

### 10.3 Unresolved report

None. ADH-2026-058 supplies the foundational executable contract closure, and
the manifest-controlled completing set — ADH-2026-060 (Maintenance-race
outcome), ADH-2026-061 (Maintenance-entry link-clearing), ADH-2026-063
(collection-create phase-one/strict-classification precedence), ADH-2026-064
(ExecutionTarget ownership correction), ADH-2026-065 (create phase-one
reference extraction and exact classification cases), and ADH-2026-066 (private
coherent backing-access bridge) — supplies the remaining create-precedence,
maintenance-entry, ownership, exact-classification, and coherent-access
mechanics. Together they form one complete, closed, executable contract, and
every requirement and acceptance criterion maps to existing registry entries and
a design disposition above. No `ARCHITECTURE_DECISION_REQUIRED`,
`REQUIREMENT_CLARIFICATION_REQUIRED`, `BOUNDARY_CHANGE_REQUIRED`,
`DEPENDENCY_APPROVAL_REQUIRED`, or `SECURITY_REVIEW_REQUIRED` condition was
encountered; no semantic choice remains open for the tasks stage.

---

Model Execution Report:
- Tool: kiro
- Stage or task: FEATURE-0016 Design revision 2 (adapter-boundary-and-executiontarget-qualification/design.md)
- Revision scope: ADH-2026-066 coherent, principal-aware backing-access rewrite; DD-13.2 and DD-13.5 now define the single paired-lease contract used through qualification publication.
- Selected model: claude-sonnet-4-20250514
- Fallback used: no

STAGE_STATUS: COMPLETE
