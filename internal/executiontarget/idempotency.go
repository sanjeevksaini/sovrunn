package executiontarget

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sort"
	"time"
)

// Completed idempotency retention bounds (REQ-F16-06; design DD-05; ADH-2026-058
// clause 5).
const (
	CompletedRetention  = 24 * time.Hour
	MaxCompletedRecords = 10000
	idempotencyKeySep   = "\x00"
)

// ReservationState is the F0016 idempotency entry lifecycle. Only InFlight and
// Completed are retained; abort removes the reservation entirely.
type ReservationState string

const (
	ReservationInFlight  ReservationState = "InFlight"
	ReservationCompleted ReservationState = "Completed"
)

// IdempotencyNamespace is the exact F0016-owned reservation key
// (principal UID, route pattern, derived CloudProvider UID, concrete
// action-target UID when applicable, Idempotency-Key). Collection creates leave
// ActionTargetUID empty.
type IdempotencyNamespace struct {
	PrincipalUID     string
	RoutePattern     string
	CloudProviderUID string
	ActionTargetUID  string
	Key              string
}

// key returns the internal map key. It must never be logged.
func (n IdempotencyNamespace) key() string {
	return n.PrincipalUID + idempotencyKeySep +
		n.RoutePattern + idempotencyKeySep +
		n.CloudProviderUID + idempotencyKeySep +
		n.ActionTargetUID + idempotencyKeySep +
		n.Key
}

// Digest is the SHA-256 of the canonical digest payload.
type Digest [32]byte

// String returns the lowercase hex encoding of the digest.
func (d Digest) String() string { return hex.EncodeToString(d[:]) }

// Equal reports whether digests are identical.
func (d Digest) Equal(other Digest) bool { return d == other }

// DigestCreate returns SHA-256 of recursively sorted allowed create JSON.
// If-Match, Accept, authentication, and request/correlation IDs are not inputs
// and must not be folded into the digest (REQ-F16-06).
func DigestCreate(allowedCreateJSON json.RawMessage) (Digest, error) {
	var bodyVal any
	if len(allowedCreateJSON) == 0 {
		bodyVal = map[string]any{}
	} else if err := json.Unmarshal(allowedCreateJSON, &bodyVal); err != nil {
		return Digest{}, err
	}
	raw, err := marshalCanonical(bodyVal)
	if err != nil {
		return Digest{}, err
	}
	return sha256.Sum256(raw), nil
}

// DigestAction returns SHA-256 of zero bytes. Qualify and retire digests never
// include If-Match, Accept, authentication, or request/correlation IDs.
func DigestAction() Digest {
	return sha256.Sum256(nil)
}

// IdempotencyApplies reports whether the HTTP method participates in FEATURE-0016
// idempotency state. Only POST create/qualify/retire do; GET and LIST never
// require, inspect, reserve, wait on, read, or replay a record.
func IdempotencyApplies(method string) bool {
	return method == http.MethodPost
}

// CompletedResult is the replayable successful response. Only status and body
// are stored; entity and transport headers are never retained.
type CompletedResult struct {
	StatusCode int
	Body       []byte
}

// reserveOutcome classifies reserveOrInspect results.
type reserveOutcome int

const (
	// reserveOutcomeReserved means this caller owns the new InFlight reservation.
	reserveOutcomeReserved reserveOutcome = iota
	// reserveOutcomeReplay means a completed same-digest record is available.
	reserveOutcomeReplay
	// reserveOutcomeConflict means same-namespace/different-digest conflict.
	reserveOutcomeConflict
	// reserveOutcomeInFlight means a same-digest InFlight peer owns the reservation.
	reserveOutcomeInFlight
)

// reserveResult is the outcome of a package-private reserve-or-inspect call.
type reserveResult struct {
	Outcome reserveOutcome
	Replay  *CompletedResult
}

// IdempotencyTable is the F0016-owned in-process idempotency reservation store
// (design DD-05). It has no independent mutex; ExecutionTargetLifecycleService
// exclusively mutates it under the single DD-02 mutex. Handlers never mutate it
// directly.
type IdempotencyTable struct {
	records map[string]*idempotencyRecord
}

type idempotencyRecord struct {
	ns        IdempotencyNamespace
	digest    Digest
	state     ReservationState
	result    *CompletedResult
	expiresAt time.Time
	// done is the single unbuffered wake channel owned by an InFlight
	// reservation. It is nil after finalize captures it for post-unlock close.
	done chan struct{}
}

// NewIdempotencyTable constructs an empty F0016-owned idempotency table.
func NewIdempotencyTable() *IdempotencyTable {
	return &IdempotencyTable{records: make(map[string]*idempotencyRecord)}
}

// reserveOrInspect reserves a new InFlight entry, returns a completed replay
// payload, reports digest conflict, or reports that a same-digest InFlight peer
// owns the reservation. The caller must hold the lifecycle mutex. now drives
// at-expiry eviction before inspection.
func (t *IdempotencyTable) reserveOrInspect(ns IdempotencyNamespace, digest Digest, now time.Time) reserveResult {
	t.evictExpired(now)
	key := ns.key()
	rec, ok := t.records[key]
	if !ok {
		t.records[key] = &idempotencyRecord{
			ns:     ns,
			digest: digest,
			state:  ReservationInFlight,
			done:   make(chan struct{}), // unbuffered; owned by this reservation
		}
		return reserveResult{Outcome: reserveOutcomeReserved}
	}

	switch rec.state {
	case ReservationCompleted:
		if !rec.digest.Equal(digest) {
			return reserveResult{Outcome: reserveOutcomeConflict}
		}
		return reserveResult{
			Outcome: reserveOutcomeReplay,
			Replay:  copyCompleted(rec.result),
		}
	case ReservationInFlight:
		if !rec.digest.Equal(digest) {
			return reserveResult{Outcome: reserveOutcomeConflict}
		}
		return reserveResult{Outcome: reserveOutcomeInFlight}
	default:
		// Unknown state is treated as absent so a deterministic retry can reserve.
		delete(t.records, key)
		t.records[key] = &idempotencyRecord{
			ns:     ns,
			digest: digest,
			state:  ReservationInFlight,
			done:   make(chan struct{}),
		}
		return reserveResult{Outcome: reserveOutcomeReserved}
	}
}

// attachWaiter returns the InFlight reservation's done channel for select.
// The caller must hold the lifecycle mutex. ok is false when the namespace is
// not InFlight.
func (t *IdempotencyTable) attachWaiter(ns IdempotencyNamespace) (done <-chan struct{}, ok bool) {
	rec, found := t.records[ns.key()]
	if !found || rec.state != ReservationInFlight || rec.done == nil {
		return nil, false
	}
	return rec.done, true
}

// detachWaiter records that a waiter stopped waiting after request cancellation.
// Cancellation detaches only; the InFlight reservation is left unchanged. With
// one shared done channel there is no per-waiter channel to remove; this method
// exists so the lifecycle protocol can invoke detach under its mutex after
// cancellation without aborting the owner reservation.
func (t *IdempotencyTable) detachWaiter(ns IdempotencyNamespace) {
	_ = ns
	// Intentionally no mutation: shared done channel + InFlight retention.
}

// finalizeComplete transitions an InFlight reservation to Completed, captures
// its done channel for post-unlock close, and enforces retention. The caller
// must hold the lifecycle mutex, unlock, then call closeReservationDone exactly
// once on the returned channel. ok is false when there is no matching InFlight.
func (t *IdempotencyTable) finalizeComplete(ns IdempotencyNamespace, digest Digest, result CompletedResult, now time.Time) (done chan struct{}, ok bool) {
	key := ns.key()
	rec, found := t.records[key]
	if !found || rec.state != ReservationInFlight || !rec.digest.Equal(digest) {
		return nil, false
	}
	done = rec.done
	rec.done = nil
	rec.state = ReservationCompleted
	rec.result = copyCompleted(&result)
	rec.expiresAt = now.Add(CompletedRetention)
	t.enforceRetention()
	return done, true
}

// finalizeAbort removes an InFlight reservation without creating a completion
// and captures its done channel for post-unlock close. The caller must hold the
// lifecycle mutex, unlock, then call closeReservationDone exactly once.
func (t *IdempotencyTable) finalizeAbort(ns IdempotencyNamespace) (done chan struct{}, ok bool) {
	key := ns.key()
	rec, found := t.records[key]
	if !found || rec.state != ReservationInFlight {
		return nil, false
	}
	done = rec.done
	delete(t.records, key)
	return done, true
}

// finalizeAbortAll removes every InFlight reservation and returns their done
// channels for post-unlock close. Used by lifecycle shutdown. Completed entries
// are left untouched.
func (t *IdempotencyTable) finalizeAbortAll() []chan struct{} {
	var channels []chan struct{}
	for key, rec := range t.records {
		if rec.state != ReservationInFlight {
			continue
		}
		if rec.done != nil {
			channels = append(channels, rec.done)
		}
		delete(t.records, key)
	}
	return channels
}

// closeReservationDone closes a captured done channel exactly once after the
// lifecycle mutex has been released. It is package-private so handlers never
// receive an exported reservation-close capability.
func closeReservationDone(done chan struct{}) {
	if done == nil {
		return
	}
	close(done)
}

// awaitReservationDone selects the reservation done channel or the request
// context. A nil error means the waiter woke on done and must re-enter at
// authentication with no lifecycle mutex held. context.Canceled / DeadlineExceeded
// means the waiter must detach under the mutex and must not abort the owner.
func awaitReservationDone(ctx context.Context, done <-chan struct{}) error {
	if done == nil {
		return context.Canceled
	}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *IdempotencyTable) evictExpired(now time.Time) {
	for key, rec := range t.records {
		if rec.state == ReservationCompleted && !rec.expiresAt.After(now) {
			delete(t.records, key)
		}
	}
}

func (t *IdempotencyTable) enforceRetention() {
	type item struct {
		key       string
		expiresAt time.Time
	}
	var completed []item
	for key, rec := range t.records {
		if rec.state != ReservationCompleted {
			continue
		}
		completed = append(completed, item{key: key, expiresAt: rec.expiresAt})
	}
	if len(completed) <= MaxCompletedRecords {
		return
	}
	sort.Slice(completed, func(i, j int) bool {
		if completed[i].expiresAt.Equal(completed[j].expiresAt) {
			return completed[i].key < completed[j].key
		}
		return completed[i].expiresAt.Before(completed[j].expiresAt)
	})
	overflow := len(completed) - MaxCompletedRecords
	for i := 0; i < overflow; i++ {
		delete(t.records, completed[i].key)
	}
}

func (t *IdempotencyTable) stateForTest(ns IdempotencyNamespace, now time.Time) (ReservationState, bool) {
	t.evictExpired(now)
	rec, ok := t.records[ns.key()]
	if !ok {
		return "", false
	}
	return rec.state, true
}

func (t *IdempotencyTable) completedCountForTest(now time.Time) int {
	t.evictExpired(now)
	n := 0
	for _, rec := range t.records {
		if rec.state == ReservationCompleted {
			n++
		}
	}
	return n
}

func (t *IdempotencyTable) inFlightCountForTest() int {
	n := 0
	for _, rec := range t.records {
		if rec.state == ReservationInFlight {
			n++
		}
	}
	return n
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
