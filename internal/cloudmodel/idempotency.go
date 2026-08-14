package cloudmodel

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
)

// ViolationIdempotencyKeyReuseMismatch is the VS0 changed-digest conflict code
// (violations[].code only; design DD-08, VS0-CF-F15-19).
const ViolationIdempotencyKeyReuseMismatch apiproblem.ViolationCode = "VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH"

// Completed idempotency retention bounds (ADH-2026-055).
const (
	CompletedRetention  = 24 * time.Hour
	MaxCompletedRecords = 10000
	idempotencyKeySep   = "\x00"
	contentTypeJSON     = "application/json"
)

// ReservationState is the idempotency namespace lifecycle (design DD-08).
type ReservationState string

const (
	ReservationInFlight  ReservationState = "InFlight"
	ReservationCompleted ReservationState = "Completed"
	ReservationAborted   ReservationState = "Aborted"
)

// IdempotencyNamespace is (principal, registered method/path, target UID when
// present, Idempotency-Key). Collection creates leave TargetUID empty.
type IdempotencyNamespace struct {
	PrincipalID string
	Pattern     string
	TargetUID   string
	Key         string
}

// key returns the internal map key. It must never be logged.
func (n IdempotencyNamespace) key() string {
	return n.PrincipalID + idempotencyKeySep + n.Pattern + idempotencyKeySep + n.TargetUID + idempotencyKeySep + n.Key
}

// Fingerprint returns a redacted correlation-safe namespace fingerprint for
// logs. It never exposes the raw Idempotency-Key or full namespace.
func (n IdempotencyNamespace) Fingerprint() string {
	sum := sha256.Sum256([]byte(n.key()))
	return hex.EncodeToString(sum[:])
}

// Digest is the SHA-256 of the canonical digest payload.
type Digest [32]byte

// String returns the lowercase hex encoding of the digest.
func (d Digest) String() string { return hex.EncodeToString(d[:]) }

// Equal reports whether digests are identical.
func (d Digest) Equal(other Digest) bool { return d == other }

// CanonicalDigestInput is the digest surface after phase-two classification
// (design DD-08). Body must already be classified canonical JSON. IfMatch is
// included only for existing-participation item actions.
type CanonicalDigestInput struct {
	Body    json.RawMessage
	Pattern string
	IfMatch string
}

// CanonicalDigest returns SHA-256 of deterministic JSON over body, pattern,
// and If-Match when present. It excludes Idempotency-Key, credentials,
// correlation, and server-assigned fields.
func CanonicalDigest(in CanonicalDigestInput) (Digest, error) {
	var bodyVal any
	if len(in.Body) == 0 {
		bodyVal = map[string]any{}
	} else if err := json.Unmarshal(in.Body, &bodyVal); err != nil {
		return Digest{}, err
	}
	payload := map[string]any{
		"body":    bodyVal,
		"pattern": in.Pattern,
	}
	if in.IfMatch != "" {
		payload["ifMatch"] = in.IfMatch
	}
	raw, err := marshalCanonical(payload)
	if err != nil {
		return Digest{}, err
	}
	return sha256.Sum256(raw), nil
}

// CompletedResult is the replayable successful response. Only status and body
// are stored; Content-Type is always application/json on replay. Entity and
// transport headers (ETag, Location, requestId, correlation) are never stored.
type CompletedResult struct {
	StatusCode int
	Body       []byte
}

// ReplayResponse is the caller-visible completed replay. Fresh request
// correlation is assigned by the HTTP layer; this value never carries
// requestId, ETag, Location, or other entity/transport headers.
type ReplayResponse struct {
	StatusCode  int
	Body        []byte
	ContentType string
}

// RecheckFunc repeats current authentication, authorization, and safe
// target/reference access before a completed replay is disclosed.
// A non-nil Problem denies without disclosing the stored result.
type RecheckFunc func(ctx context.Context) *apiproblem.Problem

// ReserveOutcome classifies the result of Reserve.
type ReserveOutcome int

const (
	// ReserveOutcomeReserved means this caller owns the InFlight reservation.
	ReserveOutcomeReserved ReserveOutcome = iota
	// ReserveOutcomeReplay means a completed same-digest record was returned.
	ReserveOutcomeReplay
	// ReserveOutcomeConflict means same-key/different-digest conflict.
	ReserveOutcomeConflict
	// ReserveOutcomeDenied means completed replay recheck denied disclosure.
	ReserveOutcomeDenied
)

// ReserveResult is the outcome of an idempotency reservation attempt.
type ReserveResult struct {
	Outcome ReserveOutcome
	Replay  *ReplayResponse
	Problem *apiproblem.Problem
}

// IdempotencyApplies reports whether method participates in FEATURE-0015
// idempotency state. Only POST collection creates and participation actions
// do; GET, LIST, and PATCH never require, look up, reserve, complete, or abort
// idempotency state (ADH-2026-056).
func IdempotencyApplies(method string) bool {
	return method == http.MethodPost
}

// IdempotencyCoordinator is the in-process idempotency reservation store
// (design DD-08). It provides no cross-process or durable replay guarantee.
type IdempotencyCoordinator struct {
	mu      sync.Mutex
	records map[string]*idempotencyRecord
	clock   Clock
	logger  *slog.Logger

	// waiters wake when a record reaches a terminal state.
	// Each waiter channel is closed exactly once.
}

type idempotencyRecord struct {
	ns          IdempotencyNamespace
	digest      Digest
	state       ReservationState
	result      *CompletedResult
	completedAt time.Time
	expiresAt   time.Time
	waiters     []chan struct{}
}

// IdempotencyCoordinatorConfig configures the coordinator.
type IdempotencyCoordinatorConfig struct {
	Clock  Clock
	Logger *slog.Logger
}

// NewIdempotencyCoordinator constructs an empty in-process coordinator.
func NewIdempotencyCoordinator(cfg IdempotencyCoordinatorConfig) *IdempotencyCoordinator {
	clock := cfg.Clock
	if clock == nil {
		clock = realClock{}
	}
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &IdempotencyCoordinator{
		records: make(map[string]*idempotencyRecord),
		clock:   clock,
		logger:  logger,
	}
}

// Reserve atomically reserves InFlight for a new namespace, waits for an
// InFlight same-digest peer, returns a completed replay after recheck, or
// returns CONFLICT on digest mismatch. On Aborted, the record is removed and
// reservation is retried from current request state.
func (c *IdempotencyCoordinator) Reserve(ctx context.Context, ns IdempotencyNamespace, digest Digest, recheck RecheckFunc) ReserveResult {
	for {
		c.mu.Lock()
		c.evictExpiredLocked(c.clock.Now())
		key := ns.key()
		rec, ok := c.records[key]
		if !ok {
			c.records[key] = &idempotencyRecord{
				ns:     ns,
				digest: digest,
				state:  ReservationInFlight,
			}
			c.logLocked("reserved", ns, "InFlight")
			c.mu.Unlock()
			return ReserveResult{Outcome: ReserveOutcomeReserved}
		}

		switch rec.state {
		case ReservationCompleted:
			if !rec.digest.Equal(digest) {
				c.mu.Unlock()
				return ReserveResult{
					Outcome: ReserveOutcomeConflict,
					Problem: idempotencyKeyReuseMismatch(),
				}
			}
			if !rec.expiresAt.After(c.clock.Now()) {
				delete(c.records, key)
				c.mu.Unlock()
				continue
			}
			stored := copyCompleted(rec.result)
			c.mu.Unlock()
			if recheck != nil {
				if prob := recheck(ctx); prob != nil {
					return ReserveResult{Outcome: ReserveOutcomeDenied, Problem: prob}
				}
			}
			return ReserveResult{
				Outcome: ReserveOutcomeReplay,
				Replay: &ReplayResponse{
					StatusCode:  stored.StatusCode,
					Body:        stored.Body,
					ContentType: contentTypeJSON,
				},
			}

		case ReservationInFlight:
			if !rec.digest.Equal(digest) {
				c.mu.Unlock()
				return ReserveResult{
					Outcome: ReserveOutcomeConflict,
					Problem: idempotencyKeyReuseMismatch(),
				}
			}
			ch := make(chan struct{})
			rec.waiters = append(rec.waiters, ch)
			c.mu.Unlock()

			select {
			case <-ctx.Done():
				c.detachWaiter(key, ch)
				return ReserveResult{
					Outcome: ReserveOutcomeDenied,
					Problem: apiproblem.New(apiproblem.CodeInternalError).WithDetail("request cancelled while waiting for idempotency reservation"),
				}
			case <-ch:
				// Terminal state reached; loop to observe Completed/Aborted.
			}

		case ReservationAborted:
			delete(c.records, key)
			c.mu.Unlock()
			continue

		default:
			c.mu.Unlock()
			return ReserveResult{
				Outcome: ReserveOutcomeDenied,
				Problem: apiproblem.New(apiproblem.CodeInternalError).WithDetail("unknown idempotency reservation state"),
			}
		}
	}
}

// Complete publishes a successful replayable result for an InFlight
// reservation. Only successful collection creates and successful participation
// actions may call Complete.
func (c *IdempotencyCoordinator) Complete(ns IdempotencyNamespace, digest Digest, result CompletedResult) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := ns.key()
	rec, ok := c.records[key]
	if !ok || rec.state != ReservationInFlight || !rec.digest.Equal(digest) {
		return
	}
	now := c.clock.Now()
	rec.state = ReservationCompleted
	rec.result = copyCompleted(&result)
	rec.completedAt = now
	rec.expiresAt = now.Add(CompletedRetention)
	c.wakeLocked(rec)
	c.enforceRetentionLocked()
	c.logLocked("completed", ns, "Completed")
}

// Abort transitions an InFlight reservation to Aborted, wakes waiters, and
// leaves no completed replay record. Used for validation, uniqueness,
// lifecycle-source-state, stale-version, audit-append failure, recovered
// panic, and graceful-shutdown outcomes.
func (c *IdempotencyCoordinator) Abort(ns IdempotencyNamespace) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.abortLocked(ns)
}

// AbortAll aborts every InFlight reservation and wakes all waiters. Invoked
// during graceful API-server shutdown after the scheduler stops.
func (c *IdempotencyCoordinator) AbortAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, rec := range c.records {
		if rec.state == ReservationInFlight {
			rec.state = ReservationAborted
			c.wakeLocked(rec)
			c.logLocked("aborted", rec.ns, "Aborted")
		}
	}
}

// StateForTest returns the reservation state for tests. The second result is
// false when the namespace has no record.
func (c *IdempotencyCoordinator) StateForTest(ns IdempotencyNamespace) (ReservationState, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evictExpiredLocked(c.clock.Now())
	rec, ok := c.records[ns.key()]
	if !ok {
		return "", false
	}
	return rec.state, true
}

// CompletedCountForTest returns the number of Completed records.
func (c *IdempotencyCoordinator) CompletedCountForTest() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.evictExpiredLocked(c.clock.Now())
	n := 0
	for _, rec := range c.records {
		if rec.state == ReservationCompleted {
			n++
		}
	}
	return n
}

// InFlightCountForTest returns the number of InFlight records.
func (c *IdempotencyCoordinator) InFlightCountForTest() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := 0
	for _, rec := range c.records {
		if rec.state == ReservationInFlight {
			n++
		}
	}
	return n
}

func (c *IdempotencyCoordinator) abortLocked(ns IdempotencyNamespace) {
	rec, ok := c.records[ns.key()]
	if !ok || rec.state != ReservationInFlight {
		return
	}
	rec.state = ReservationAborted
	c.wakeLocked(rec)
	c.logLocked("aborted", ns, "Aborted")
}

func (c *IdempotencyCoordinator) wakeLocked(rec *idempotencyRecord) {
	for _, ch := range rec.waiters {
		close(ch)
	}
	rec.waiters = nil
}

func (c *IdempotencyCoordinator) detachWaiter(key string, ch chan struct{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	rec, ok := c.records[key]
	if !ok {
		return
	}
	filtered := rec.waiters[:0]
	for _, w := range rec.waiters {
		if w != ch {
			filtered = append(filtered, w)
		}
	}
	rec.waiters = filtered
	// Cancellation detaches only this waiter; reservation remains InFlight.
}

func (c *IdempotencyCoordinator) evictExpiredLocked(now time.Time) {
	for key, rec := range c.records {
		if rec.state == ReservationCompleted && !rec.expiresAt.After(now) {
			delete(c.records, key)
		}
	}
}

func (c *IdempotencyCoordinator) enforceRetentionLocked() {
	type item struct {
		key       string
		expiresAt time.Time
		nsKey     string
	}
	var completed []item
	for key, rec := range c.records {
		if rec.state != ReservationCompleted {
			continue
		}
		completed = append(completed, item{key: key, expiresAt: rec.expiresAt, nsKey: key})
	}
	if len(completed) <= MaxCompletedRecords {
		return
	}
	sort.Slice(completed, func(i, j int) bool {
		if completed[i].expiresAt.Equal(completed[j].expiresAt) {
			return completed[i].nsKey < completed[j].nsKey
		}
		return completed[i].expiresAt.Before(completed[j].expiresAt)
	})
	overflow := len(completed) - MaxCompletedRecords
	for i := 0; i < overflow; i++ {
		delete(c.records, completed[i].key)
	}
}

func (c *IdempotencyCoordinator) logLocked(outcome string, ns IdempotencyNamespace, state string) {
	if c.logger == nil {
		return
	}
	c.logger.Info("cloudmodel.idempotency",
		"feature_id", "FEATURE-0015",
		"idempotency_namespace_fingerprint", ns.Fingerprint(),
		"outcome", outcome,
		"state", state,
	)
}

func idempotencyKeyReuseMismatch() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeConflict).
		WithDetail("Idempotency-Key was reused with a different request digest").
		WithViolations([]apiproblem.Violation{{
			Field:   "/headers/Idempotency-Key",
			Code:    ViolationIdempotencyKeyReuseMismatch,
			Message: "idempotency key reused with a different canonical digest",
		}})
}

func copyCompleted(in *CompletedResult) *CompletedResult {
	if in == nil {
		return &CompletedResult{}
	}
	out := &CompletedResult{StatusCode: in.StatusCode}
	if in.Body != nil {
		out.Body = append([]byte(nil), in.Body...)
	}
	return out
}

func marshalCanonical(v any) ([]byte, error) {
	switch t := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		buf := []byte{'{'}
		for i, k := range keys {
			if i > 0 {
				buf = append(buf, ',')
			}
			kb, err := json.Marshal(k)
			if err != nil {
				return nil, err
			}
			buf = append(buf, kb...)
			buf = append(buf, ':')
			vb, err := marshalCanonical(t[k])
			if err != nil {
				return nil, err
			}
			buf = append(buf, vb...)
		}
		buf = append(buf, '}')
		return buf, nil
	case []any:
		buf := []byte{'['}
		for i, elem := range t {
			if i > 0 {
				buf = append(buf, ',')
			}
			eb, err := marshalCanonical(elem)
			if err != nil {
				return nil, err
			}
			buf = append(buf, eb...)
		}
		buf = append(buf, ']')
		return buf, nil
	default:
		return json.Marshal(v)
	}
}
