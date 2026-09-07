package idempotency

import (
	"crypto/sha256"
	"sync"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
)

// ReservationTable owns reservation-table operations over Task-2 value types.
// It stores in-flight reservations and completed records by lookup key digest.
type ReservationTable struct {
	mu        sync.Mutex
	inflight  map[string]inflightEntry
	completed map[string]completedEntry
}

type inflightEntry struct {
	reservation InFlightReservation
	wait        chan struct{}
}

type completedEntry struct {
	record      CompletedRecord
	retainUntil time.Time
}

// NewReservationTable constructs an empty reservation table.
func NewReservationTable() *ReservationTable {
	return &ReservationTable{
		inflight:  map[string]inflightEntry{},
		completed: map[string]completedEntry{},
	}
}

// LookupResult is a closed lookup result for an idempotency key.
type LookupResult struct {
	hasInFlight bool
	inFlight    InFlightReservation
	wait        <-chan struct{}

	hasCompleted bool
	completed    CompletedRecord
	retainUntil  time.Time
}

// InFlight returns the reservation and waiter when an in-flight entry exists.
func (r LookupResult) InFlight() (InFlightReservation, <-chan struct{}, bool) {
	if !r.hasInFlight {
		return InFlightReservation{}, nil, false
	}
	return r.inFlight, r.wait, true
}

// Completed returns the completed record and retention horizon if present.
func (r LookupResult) Completed() (CompletedRecord, time.Time, bool) {
	if !r.hasCompleted {
		return CompletedRecord{}, time.Time{}, false
	}
	return r.completed, r.retainUntil, true
}

// Lookup returns the current reservation/completed entry for key.
func (t *ReservationTable) Lookup(key IdempotencyLookupKey) (LookupResult, operation.MechanicalFailure) {
	if !key.Sealed() {
		return LookupResult{}, operation.NewMechanicalFailure("idempotency_lookup_unsealed")
	}
	d := lookupDigest(key)
	t.mu.Lock()
	defer t.mu.Unlock()
	if c, ok := t.completed[d]; ok {
		return LookupResult{
			hasCompleted: true,
			completed:    c.record,
			retainUntil:  c.retainUntil,
		}, operation.MechanicalFailure{}
	}
	if i, ok := t.inflight[d]; ok {
		return LookupResult{
			hasInFlight: true,
			inFlight:    i.reservation,
			wait:        i.wait,
		}, operation.MechanicalFailure{}
	}
	return LookupResult{}, operation.MechanicalFailure{}
}

// Reserve stores reservation for key if not already present.
func (t *ReservationTable) Reserve(reservation InFlightReservation) (bool, <-chan struct{}, operation.MechanicalFailure) {
	if !reservation.Sealed() {
		return false, nil, operation.NewMechanicalFailure("idempotency_reservation_unsealed")
	}
	d := lookupDigest(reservation.LookupKey())
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.completed[d]; ok {
		return false, nil, operation.NewMechanicalFailure("idempotency_reserve_completed_exists")
	}
	if existing, ok := t.inflight[d]; ok {
		return false, existing.wait, operation.MechanicalFailure{}
	}
	wait := make(chan struct{})
	t.inflight[d] = inflightEntry{reservation: reservation, wait: wait}
	return true, wait, operation.MechanicalFailure{}
}

// Complete transitions an in-flight reservation into completed state.
func (t *ReservationTable) Complete(prepared PreparedCompletedResult, retainUntil time.Time) operation.MechanicalFailure {
	if !prepared.Sealed() {
		return operation.NewMechanicalFailure("idempotency_complete_prepared_unsealed")
	}
	if retainUntil.IsZero() {
		return operation.NewMechanicalFailure("idempotency_complete_retain_until_zero")
	}
	d := lookupDigest(prepared.LookupKey())
	record, fail := NewCompletedRecord(prepared.LookupKey(), prepared.Binding(), prepared.Result(), prepared.Completion())
	if fail.Reason() != "" {
		return fail
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	entry, ok := t.inflight[d]
	if !ok {
		return operation.NewMechanicalFailure("idempotency_complete_reservation_missing")
	}
	delete(t.inflight, d)
	t.completed[d] = completedEntry{record: record, retainUntil: retainUntil.UTC()}
	close(entry.wait)
	return operation.MechanicalFailure{}
}

// ReleaseInFlight aborts an in-flight reservation without publishing a completed record.
func (t *ReservationTable) ReleaseInFlight(key IdempotencyLookupKey) operation.MechanicalFailure {
	if !key.Sealed() {
		return operation.NewMechanicalFailure("idempotency_release_lookup_unsealed")
	}
	d := lookupDigest(key)
	t.mu.Lock()
	defer t.mu.Unlock()
	if entry, ok := t.inflight[d]; ok {
		delete(t.inflight, d)
		close(entry.wait)
	}
	return operation.MechanicalFailure{}
}

// PruneExpired removes completed records that have passed retainUntil.
func (t *ReservationTable) PruneExpired(now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for k, rec := range t.completed {
		if !now.Before(rec.retainUntil) {
			delete(t.completed, k)
		}
	}
}

func lookupDigest(key IdempotencyLookupKey) string {
	sum := sha256.Sum256([]byte(key.ActorRef().Issuer + "|" + key.ActorRef().Subject + "|" + key.Operation() + "|" + key.ClientIdempotencyKey()))
	return string(sum[:])
}
