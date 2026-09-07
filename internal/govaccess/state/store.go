package state

import (
	"crypto/sha256"
	"sync"
	"sync/atomic"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
)

// Store is the neutral in-memory aggregate root (DD-02, DD-10).
// Lock order: never hold reservationMu while acquiring stateMu.
type Store struct {
	reservationMu sync.Mutex
	stateMu       sync.Mutex

	clk clock.Clock

	root     map[string]resourceEntry
	versions map[string]string // uid -> version
	snapshot string
	sequence uint64

	inflight  map[string]idempotency.InFlightReservation // lookup key digest -> reservation
	completed map[string]idempotency.CompletedRecord

	stateMuHeld atomic.Bool
}

// NewStore constructs an empty in-memory store with the injected clock.
func NewStore(clk clock.Clock) *Store {
	return &Store{
		clk:       clk,
		root:      map[string]resourceEntry{},
		versions:  map[string]string{},
		snapshot:  "0",
		inflight:  map[string]idempotency.InFlightReservation{},
		completed: map[string]idempotency.CompletedRecord{},
	}
}

// SetResourceVersion seeds a resource version for dependency revalidation tests.
func (s *Store) SetResourceVersion(uid, version string) {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	s.versions[uid] = version
	s.root[uid] = resourceEntry{version: version, payload: []byte(uid)}
}

// SnapshotVersion returns the current aggregate snapshot version.
func (s *Store) SnapshotVersion() string {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	return s.snapshot
}

// InspectOutcome is the closed InspectOrReserveCallerMutation result.
type InspectOutcome struct {
	kind      inspectKind
	conflict  bool
	replay    idempotency.CompletedRecord
	hasReplay bool
	waitKey   string
	hasWait   bool
	lease     *CallerReservationLease
}

type inspectKind uint8

const (
	inspectNone inspectKind = iota
	inspectConflict
	inspectReplay
	inspectWait
	inspectOwner
)

// IsConflict reports Conflict.
func (o InspectOutcome) IsConflict() bool { return o.kind == inspectConflict }

// Replay returns a CompletedRecord when Replay.
func (o InspectOutcome) Replay() (idempotency.CompletedRecord, bool) {
	return o.replay, o.hasReplay
}

// WaitKey returns a waiter key when Wait.
func (o InspectOutcome) WaitKey() (string, bool) { return o.waitKey, o.hasWait }

// Owner returns the reservation lease when Owner.
func (o InspectOutcome) Owner() (*CallerReservationLease, bool) {
	if o.kind != inspectOwner || o.lease == nil {
		return nil, false
	}
	return o.lease, true
}

// InspectOrReserveCallerMutation implements the bounded reservation algorithm
// (design §5.4). Lock order: never hold reservationMu while acquiring stateMu.
func (s *Store) InspectOrReserveCallerMutation(
	key idempotency.IdempotencyLookupKey,
	binding idempotency.RequestBinding,
) InspectOutcome {
	lookup := lookupDigest(key)

	s.reservationMu.Lock()
	if res, ok := s.inflight[lookup]; ok {
		if !res.Binding().Equal(binding) {
			s.reservationMu.Unlock()
			return InspectOutcome{kind: inspectConflict, conflict: true}
		}
		s.reservationMu.Unlock()
		return InspectOutcome{kind: inspectWait, hasWait: true, waitKey: lookup}
	}
	s.reservationMu.Unlock()

	s.stateMu.Lock()
	s.reservationMu.Lock()
	if rec, ok := s.completed[lookup]; ok {
		s.reservationMu.Unlock()
		s.stateMu.Unlock()
		if !rec.Binding().Equal(binding) {
			return InspectOutcome{kind: inspectConflict, conflict: true}
		}
		return InspectOutcome{kind: inspectReplay, replay: rec, hasReplay: true}
	}
	if res, ok := s.inflight[lookup]; ok {
		same := res.Binding().Equal(binding)
		s.reservationMu.Unlock()
		s.stateMu.Unlock()
		if !same {
			return InspectOutcome{kind: inspectConflict, conflict: true}
		}
		return InspectOutcome{kind: inspectWait, hasWait: true, waitKey: lookup}
	}
	token := sha256.Sum256([]byte(lookup + ":owner"))
	res, fail := idempotency.NewInFlightReservation(key, binding, token[:])
	if fail.Reason() != "" {
		s.reservationMu.Unlock()
		s.stateMu.Unlock()
		return InspectOutcome{kind: inspectConflict, conflict: true}
	}
	s.inflight[lookup] = res
	lease := &CallerReservationLease{store: s, lookup: lookup, token: copyBytes(token[:]), binding: binding, key: key}
	s.reservationMu.Unlock()
	s.stateMu.Unlock()
	return InspectOutcome{kind: inspectOwner, lease: lease}
}

// CallerReservationLease holds reservation ownership before Begin.
type CallerReservationLease struct {
	store   *Store
	lookup  string
	token   []byte
	binding idempotency.RequestBinding
	key     idempotency.IdempotencyLookupKey
	invalid bool
}

// Abort releases an unused reservation.
func (l *CallerReservationLease) Abort() {
	if l == nil || l.invalid {
		return
	}
	l.store.reservationMu.Lock()
	defer l.store.reservationMu.Unlock()
	delete(l.store.inflight, l.lookup)
	l.invalid = true
}

// Begin starts a CallerMutationTransaction after dependency revalidation.
func (l *CallerReservationLease) Begin(admission MutationAdmission) (CallerMutationTransaction, error) {
	if l == nil || l.invalid {
		return CallerMutationTransaction{}, ErrConflict
	}
	if !admission.sealed || admission.mode != CallerMutationMode {
		return CallerMutationTransaction{}, ErrConflict
	}

	l.store.stateMu.Lock()
	l.store.stateMuHeld.Store(true)
	l.store.reservationMu.Lock()
	res, ok := l.store.inflight[l.lookup]
	if !ok || !digestEqual(res.OwnerToken(), l.token) {
		l.store.reservationMu.Unlock()
		l.store.releaseStateMu()
		return CallerMutationTransaction{}, ErrConflict
	}
	l.store.reservationMu.Unlock()

	if err := l.store.revalidateAdmission(admission); err != nil {
		l.store.reservationMu.Lock()
		delete(l.store.inflight, l.lookup)
		l.store.reservationMu.Unlock()
		l.store.releaseStateMu()
		l.invalid = true
		return CallerMutationTransaction{}, err
	}

	ctx, proof, err := l.store.captureContextAndProof(admission)
	if err != nil {
		l.store.reservationMu.Lock()
		delete(l.store.inflight, l.lookup)
		l.store.reservationMu.Unlock()
		l.store.releaseStateMu()
		l.invalid = true
		return CallerMutationTransaction{}, err
	}

	base := &baseTransaction{
		store:             l.store,
		admission:         admission,
		context:           ctx,
		admissionProof:    proof,
		hasAdmissionProof: true,
		editor:            newMemoryEditor(l.store.root),
		permits:           map[string]*permitRegistryEntry{},
	}
	l.invalid = true
	return CallerMutationTransaction{base: base}, nil
}

// BeginControllerMutation starts a controller mutation transaction.
func (s *Store) BeginControllerMutation(admission MutationAdmission) (ControllerMutationTransaction, error) {
	if !admission.sealed || admission.mode != ControllerMutationMode {
		return ControllerMutationTransaction{}, ErrConflict
	}
	s.stateMu.Lock()
	s.stateMuHeld.Store(true)
	if err := s.revalidateAdmission(admission); err != nil {
		s.releaseStateMu()
		return ControllerMutationTransaction{}, err
	}
	ctx, proof, err := s.captureContextAndProof(admission)
	if err != nil {
		s.releaseStateMu()
		return ControllerMutationTransaction{}, err
	}
	base := &baseTransaction{
		store:     s,
		admission: admission,
		context:   ctx,
		editor:    newMemoryEditor(s.root),
		permits:   map[string]*permitRegistryEntry{},
	}
	if admission.authority == AuthorizationBearing {
		base.admissionProof = proof
		base.hasAdmissionProof = true
	}
	return ControllerMutationTransaction{base: base}, nil
}

// BeginEvidenceTransaction starts an evidence-only transaction.
func (s *Store) BeginEvidenceTransaction(claim AuthorizationCurrentnessClaim) (EvidenceTransaction, error) {
	if !claim.sealed {
		return EvidenceTransaction{}, ErrRetryEvaluation
	}
	s.stateMu.Lock()
	s.stateMuHeld.Store(true)
	if claim.dependency.expectedSnapshotVersion != s.snapshot {
		s.releaseStateMu()
		return EvidenceTransaction{}, ErrRetryEvaluation
	}
	if !predicatesMatch(s.versions, claim.dependency.predicates) {
		s.releaseStateMu()
		return EvidenceTransaction{}, ErrRetryEvaluation
	}
	s.sequence++
	ctx, fail := evidence.NewPublicationContext(s.clk.NowUTC(), s.sequence)
	if fail.Reason() != "" {
		s.releaseStateMu()
		return EvidenceTransaction{}, ErrRetryEvaluation
	}
	dep, fail := claim.dependency.AsDependencyVersionDigest()
	if fail.Reason() != "" {
		s.releaseStateMu()
		return EvidenceTransaction{}, ErrRetryEvaluation
	}
	proof, fail := evidence.BindAuthorizationAdmission(claim.candidate, dep, ctx)
	if fail.Reason() != "" {
		s.releaseStateMu()
		return EvidenceTransaction{}, ErrRetryEvaluation
	}
	base := &baseTransaction{
		store:             s,
		context:           ctx,
		admissionProof:    proof,
		hasAdmissionProof: true,
		editor:            newMemoryEditor(s.root),
		admission: MutationAdmission{
			sealed:    true,
			mode:      ControllerMutationMode,
			authority: AuthorizationBearing,
		},
	}
	return EvidenceTransaction{base: base}, nil
}

func (s *Store) revalidateAdmission(admission MutationAdmission) error {
	for _, c := range admission.participants {
		if !predicatesMatch(s.versions, c.versions.predicates) {
			return ErrConflict
		}
	}
	if admission.hasCurrent {
		if admission.currentness.dependency.expectedSnapshotVersion != s.snapshot {
			return ErrConflict
		}
		if !predicatesMatch(s.versions, admission.currentness.dependency.predicates) {
			return ErrConflict
		}
	}
	return nil
}

func (s *Store) captureContextAndProof(admission MutationAdmission) (evidence.PublicationContext, evidence.AuthorizationAdmissionProof, error) {
	s.sequence++
	ctx, fail := evidence.NewPublicationContext(s.clk.NowUTC(), s.sequence)
	if fail.Reason() != "" {
		return evidence.PublicationContext{}, evidence.AuthorizationAdmissionProof{}, ErrConflict
	}
	if admission.authority != AuthorizationBearing {
		return ctx, evidence.AuthorizationAdmissionProof{}, nil
	}
	dep, fail := admission.currentness.dependency.AsDependencyVersionDigest()
	if fail.Reason() != "" {
		return evidence.PublicationContext{}, evidence.AuthorizationAdmissionProof{}, ErrConflict
	}
	proof, fail := evidence.BindAuthorizationAdmission(admission.currentness.candidate, dep, ctx)
	if fail.Reason() != "" {
		return evidence.PublicationContext{}, evidence.AuthorizationAdmissionProof{}, ErrConflict
	}
	return ctx, proof, nil
}

func (s *Store) installRoot(root map[string]resourceEntry) {
	s.root = root
	s.versions = map[string]string{}
	for uid, e := range root {
		s.versions[uid] = e.version
	}
	s.snapshot = formatUint64(s.sequence)
}

func (s *Store) storeCompleted(prep idempotency.PreparedCompletedResult) {
	rec, fail := idempotency.NewCompletedRecord(prep.LookupKey(), prep.Binding(), prep.Result(), prep.Completion())
	if fail.Reason() != "" {
		return
	}
	s.reservationMu.Lock()
	defer s.reservationMu.Unlock()
	lookup := lookupDigest(prep.LookupKey())
	s.completed[lookup] = rec
	delete(s.inflight, lookup)
}

func (s *Store) releaseStateMu() {
	if s.stateMuHeld.Swap(false) {
		s.stateMu.Unlock()
	}
}

// LockOrderObservation is used by tests to assert reservationMu → stateMu order.
func (s *Store) LockOrderObservation(fn func()) {
	s.reservationMu.Lock()
	s.reservationMu.Unlock()
	s.stateMu.Lock()
	s.reservationMu.Lock()
	fn()
	s.reservationMu.Unlock()
	s.stateMu.Unlock()
}

func lookupDigest(key idempotency.IdempotencyLookupKey) string {
	h := sha256.Sum256([]byte(key.ActorRef().Issuer + "|" + key.ActorRef().Subject + "|" + key.Operation() + "|" + key.ClientIdempotencyKey()))
	return string(h[:])
}

func formatUint64(v uint64) string {
	const digits = "0123456789"
	if v == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = digits[v%10]
		v /= 10
	}
	return string(buf[i:])
}
