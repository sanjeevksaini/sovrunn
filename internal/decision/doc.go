// Package decision defines the FEATURE-0013 Decision Record and AuditEvent
// domain value types and closed vocabularies.
//
// Domain ownership (RID-01, design closure 7): canonical FEATURE-0013 domain
// types live in this package tree. internal/apiconform remains limited to
// TypeBinding registration, conformance adapters, and executable checks; it
// does not become the canonical domain owner and must not receive divergent
// copies of FEATURE-0012 base fields. The FEATURE-0013 AuditEvent
// payload/linkage extension value types are owned here.
//
// This package is stdlib-only at the domain root and may import internal/apimeta
// only. It does not import internal/server, internal/api, or runtime/provider
// packages.
//
// # Invariant and deferred notes (T-031; documentation only)
//
// The following obligations are recorded here as INVARIANT_FOR_LATER or
// DEFERRED. FEATURE-0013 carries carrier fields, structural validation, and
// conformance fixtures only. No persistence, delivery, cache, authZ runtime,
// cryptographic execution, observability runtime, vendor selection, or other
// production mechanism is implemented by this package.
//
// ## INVARIANT_FOR_LATER
//
// AD-007 — EffectivePolicyContext (architecture §7.3; F13-COMP-005):
// Runtime hierarchy resolution must produce one version-pinned
// EffectivePolicyContext with provenance, precedence, conflict resolution,
// effective time, and optional opaque integrity carrier. Decision domains must
// not independently traverse Platform/Organization/OrganizationUnit/Tenant/
// Project/Provider trees. FEATURE-0013 records only typed reference and
// provenance contract fields; no hierarchy resolver or cache runtime is
// created.
//
// AD-008 — Synchronous fast path (architecture §8.1, §9; F13-SCALE-005/006):
// Simple and bounded decisions should complete synchronously without a
// mandatory queue or durable workflow engine when inputs fit the declared
// latency budget. FEATURE-0013 provides the pure reference kernel
// (compose.Replay + strategies) only; no production worker, queue, or
// scale-out path is created.
//
// AD-010 — Orchestrate authority; choreograph reactions (architecture §8.3):
// Authoritative evaluation ordering, composition, retry policy, deadlines,
// human approval, and finalization remain orchestrated. After durable
// finalization, choreography may drive non-authoritative reactions
// (notifications, indexing, telemetry, cache invalidation). Event consumers
// must never silently alter or replace the authoritative decision. No
// orchestrator or choreography runtime is created here.
//
// AD-013 — Effective-once via idempotency (architecture §9; F13-SCALE-006):
// Exactly-once transport is not claimed. Business-level effective-once
// behavior comes from idempotency keys, immutable identity, and atomic
// persistence boundaries in a later owning feature. FEATURE-0013 implements
// only the retry-key and semantic-identity carriers (T-005) plus validation.
//
// AD-014 / AD-035 — Atomic durable acceptance and async audit materialization
// (design §11; architecture §10; F13-AUDIT-005/006/007):
// A final decision must not be reported durable until its DecisionRecord and
// required audit obligation are atomically accepted by the same local
// authoritative persistence boundary. The boundary preallocates the immutable
// AuditEvent identity and stores that identity in both the DecisionRecord
// reference and the audit obligation. Durable audit-obligation acceptance is
// sufficient for the synchronous response; full AuditEvent materialization,
// indexing, export, notification, SIEM delivery, cross-site replication, or
// regulator delivery may proceed asynchronously and must reuse the
// preallocated identity without mutating the DecisionRecord. Recovery,
// bounded delivery, poison-record quarantine, and reconciliation for partial
// downstream failure are later invariants. No transaction, outbox, event
// store, queue, or audit delivery service is created here.
//
// AD-039 / F13-AUTHZ-004/005/006 — Local authZ, every-hop duty, break-glass
// follow-up (architecture §12.1; design §12):
// Production authorization may use local verified artifacts, policy snapshots,
// and safe caches with profile-defined audience, scope, expiry, revocation,
// freshness, and fail-closed behavior; a remote central authorization call is
// not required on every request. Every human, workload, evaluator, and service
// hop in a decision or audit flow must be authenticated and authorized, bound
// to scope, purpose, action, and resource. Break-glass or emergency access
// requires follow-up review after the fact (the mandatory break-glass reason
// carrier itself is CONTRACT_NOW under F13-SOV-007). No authorization cache,
// session, enforcement, or follow-up-review workflow is created here.
//
// AD-041 — Tiered cryptographic assurance (architecture §12.3; F13-TRUST-004):
// Cryptographic assurance is tiered. Remote notary/HSM/trusted-time must not
// become a universal hot-path dependency. Local, batch, and asynchronous
// assurance profiles are later permitted. Selections and execution remain
// blocked by ADR-F13-002; FEATURE-0013 carries algorithm-agile structural
// TrustCarrier fields only.
//
// F13-SCALE-006 — Production scaling invariants (architecture §9):
// Stateless interchangeable workers, external durable state behind ports,
// compare-and-create semantics, effective-once behavior, bounded
// delivery/retry/backpressure/cache behavior, runtime resource budgets, no
// sticky sessions, and scale-to-zero where policy permits are later invariants
// carried and documented only.
//
// Architecture §16 — Observability and operability:
// Later owning features must expose latency, queueing, evaluation fan-out,
// cache result, composition strategy, failure class, retry, saturation, and
// audit-reconciliation health without exposing decision evidence or personal
// data as labels. Logs must be structured, redacted, locally routable,
// retention-controlled, and correlated by opaque identifiers. Health must
// distinguish evaluator, evidence, store, audit, identity, key, time, and
// synchronization dependencies. Operators must diagnose and recover without
// vendor remote access. FEATURE-0013 creates no metrics, tracing, logging, or
// health runtime.
//
// ## DEFERRED
//
// AD-011 — Runtime and vendor selection (architecture §8.4, §21; non-goals):
// Durable-execution systems, workflow engines, message brokers, queues,
// event-streaming platforms, databases, policy engines, identity products,
// cryptographic products, and AI/model runtimes are not selected. A later
// runtime decision must satisfy FEATURE-0013 contracts and demonstrate
// sovereign operability. FEATURE-0013 defines no wrapper, adapter interface,
// or runtime implementation.
//
// F13-SCALE-007 — Concrete scaling mechanisms:
// Every concrete store, transaction, lock, lease, queue, dead-letter
// processor, cache, worker topology, scaling policy, and production numeric
// default is deferred.
//
// ADR-F13-002 — Cryptographic selection and execution (architecture §12.3;
// F13-TRUST-004):
// Canonicalization algorithm/profile, digest-covered field set, signature
// algorithm, and cryptographic product/service selection remain deferred.
// FEATURE-0013 must not compute digests over content, canonicalize content,
// sign, verify, timestamp, notarize, enforce WORM retention, or claim that an
// opaque carrier proves content existed or was considered. RFC 8785 is
// illustrative only and is not selected. Structural trust-state validation
// (RID-08 empty-state rule; DECISION_TRUST_* codes) remains CONTRACT_NOW.
package decision
