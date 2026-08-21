package executiontarget

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// AuditAppender appends a redacted FEATURE-0013 AuditEvent. FEATURE-0016 does
// not own a durable AuditEvent store; the appender is injected (VS0-WRITER-011
// writer authority remains FEATURE-0013).
type AuditAppender interface {
	Append(ctx context.Context, event apiconform.AuditEvent) error
}

// ReservationInspectOutcome classifies ReserveOrInspectReplay results for
// handlers. Completion remains commit-internal.
type ReservationInspectOutcome int

const (
	// ReservationOutcomeReserved means this caller owns a new InFlight reservation.
	ReservationOutcomeReserved ReservationInspectOutcome = iota
	// ReservationOutcomeReplay means a completed same-digest record is available.
	ReservationOutcomeReplay
	// ReservationOutcomeConflict means same-namespace/different-digest conflict.
	ReservationOutcomeConflict
	// ReservationOutcomeInFlight means a same-digest InFlight peer owns the reservation.
	ReservationOutcomeInFlight
)

// ReservationInspectResult is the handler-facing reserve/replay inspection
// outcome. Handlers never mutate or complete the table.
type ReservationInspectResult struct {
	Outcome ReservationInspectOutcome
	Replay  *CompletedResult
}

// CreateCommitRequest is the closed create commit input for
// ExecutionTargetLifecycleService.CommitCreate. Handlers supply
// participation-derived CloudProvider scope and stack generation after safe
// backing access; the service validates the closed create-field boundary and
// persists Active/Unqualified at maintenance epoch 0.
type CreateCommitRequest struct {
	Name                  string
	UID                   string
	CloudProviderScopeUID string
	ParticipationRef      apimeta.TypedRef
	StackRef              apimeta.TypedRef
	StackGeneration       int64
	TargetClass           model.TargetClass

	// ClientStatusAttempted is true when the create envelope attempted to
	// supply any status or server-owned identity field. Commit rejects these
	// without mutation (clients-never-write-status).
	ClientStatusAttempted bool

	Actor     apimeta.TypedRef
	RequestID string
	AuditUID  string

	Idempotency IdempotencyNamespace
	Digest      Digest
	Completion  CompletedResult
}

// RetireCommitRequest is the Active-state retire commit input. When a
// Qualifying reservation is in flight for the same target, a successful
// retirement aborts that reservation (Retire-wins; TASK-F16-05).
type RetireCommitRequest struct {
	TargetUID               string
	ExpectedResourceVersion string

	Actor     apimeta.TypedRef
	RequestID string
	AuditUID  string

	Idempotency IdempotencyNamespace
	Digest      Digest
	Completion  CompletedResult
}

// BackingViability is the already-resolved FEATURE-0015 backing snapshot used
// for the under-mutex qualify recheck and fence capture (REQ-F16-02).
type BackingViability struct {
	ParticipationUID             string
	ParticipationEffectiveActive bool
	ParticipationScopeUID        string
	StackUID                     string
	StackPhase                   string
	StackScopeUID                string
	StackGeneration              int64
}

// StackActive reports whether the InfrastructureStack phase is Active.
func (b BackingViability) StackActive() bool {
	return b.StackPhase == "Active"
}

// Fingerprint returns the closed viability fingerprint for fence capture.
func (b BackingViability) Fingerprint() model.ViabilityFingerprint {
	return model.ComputeViabilityFingerprint(
		b.ParticipationUID,
		b.ParticipationEffectiveActive,
		b.StackUID,
		b.StackPhase,
	)
}

// QualifyCommitRequest is the synchronous qualify input for
// ExecutionTargetLifecycleService.Qualify (DD-04). The action idempotency
// reservation must already exist; this path opens a distinct internal
// Qualifying reservation, observes unlocked, then commits under lock.
type QualifyCommitRequest struct {
	TargetUID string

	// Backing is the already-resolved participation/stack snapshot for the
	// under-mutex viability recheck and initial fence capture.
	Backing BackingViability
	// RefreshBacking, when non-nil, is invoked under the commit mutex to obtain
	// the current generation/viability fences. Nil means Backing remains current.
	RefreshBacking func() BackingViability

	FactSetUID string
	ResultUID  string

	Actor     apimeta.TypedRef
	RequestID string
	AuditUID  string

	Idempotency IdempotencyNamespace
	Digest      Digest
	Completion  CompletedResult
}

// capturedQualifyFences are the fences sealed when the Qualifying reservation
// opens (DD-04).
type capturedQualifyFences struct {
	stackGeneration      int64
	maintenanceEpoch     int64
	viabilityFingerprint model.ViabilityFingerprint
}

// qualifyingReservation is the lifecycle-service-only in-flight Qualifying
// reservation. It is never persisted or projected.
type qualifyingReservation struct {
	targetUID   string
	fences      capturedQualifyFences
	idempotency IdempotencyNamespace
	digest      Digest
	// aborted is set when Retire/Maintenance/shutdown wins the race.
	aborted *apiproblem.Problem
}

// ExecutionTargetLifecycleService is the sole committer of ExecutionTarget
// state and internal records (VS0-WRITER-006; DD-02). One service-owned mutex
// guards all F0016 targets, indexes, internal records, Maintenance markers,
// Qualifying reservations, and idempotency reservations. There are no
// per-target locks.
type ExecutionTargetLifecycleService struct {
	mu sync.Mutex

	store      *Store
	idemp      *IdempotencyTable
	audit      AuditAppender
	observer   *SyntheticObserver
	qualifying map[string]*qualifyingReservation // targetUID -> reservation
	stopped    bool
	now        func() time.Time
}

// LifecycleConfig configures the sole-committer lifecycle service.
type LifecycleConfig struct {
	Store       *Store
	Idempotency *IdempotencyTable
	Audit       AuditAppender
	Observer    *SyntheticObserver
	// Now injects the clock for idempotency retention and observation. Nil
	// defaults to time.Now().UTC. TASK-F16-06 owns the shared Clock interface.
	Now func() time.Time
}

// NewExecutionTargetLifecycleService constructs the VS0-WRITER-006 sole
// committer. Store and IdempotencyTable must be non-nil.
func NewExecutionTargetLifecycleService(cfg LifecycleConfig) *ExecutionTargetLifecycleService {
	now := cfg.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	store := cfg.Store
	if store == nil {
		store = NewStore()
	}
	idemp := cfg.Idempotency
	if idemp == nil {
		idemp = NewIdempotencyTable()
	}
	observer := cfg.Observer
	if observer == nil {
		observer = NewSyntheticObserver(NewMapFixtureSource(), now)
	}
	return &ExecutionTargetLifecycleService{
		store:      store,
		idemp:      idemp,
		audit:      cfg.Audit,
		observer:   observer,
		qualifying: make(map[string]*qualifyingReservation),
		now:        now,
	}
}

// Store returns the process-local registry. Mutation remains sole-committer
// owned; callers must not invoke package-private staging helpers.
func (s *ExecutionTargetLifecycleService) Store() *Store {
	return s.store
}

// lock acquires the DD-02 sole-committer mutex and the store publication lock
// required by staging helpers. Fixed order: service mu, then store publication.
// Idempotency table mutations always occur under this same critical section.
func (s *ExecutionTargetLifecycleService) lock() {
	s.mu.Lock()
	s.store.beginPublication()
}

// unlock releases the store publication lock then the DD-02 mutex.
func (s *ExecutionTargetLifecycleService) unlock() {
	s.store.endPublication()
	s.mu.Unlock()
}

// ReserveOrInspectReplay reserves a new InFlight entry, returns a completed
// replay payload, reports digest conflict, or reports a same-digest InFlight
// peer. Holds the DD-02 mutex only for the table operation.
func (s *ExecutionTargetLifecycleService) ReserveOrInspectReplay(ns IdempotencyNamespace, digest Digest) ReservationInspectResult {
	s.lock()
	defer s.unlock()
	res := s.idemp.reserveOrInspect(ns, digest, s.now())
	out := ReservationInspectResult{}
	switch res.Outcome {
	case reserveOutcomeReserved:
		out.Outcome = ReservationOutcomeReserved
	case reserveOutcomeReplay:
		out.Outcome = ReservationOutcomeReplay
		out.Replay = res.Replay
	case reserveOutcomeConflict:
		out.Outcome = ReservationOutcomeConflict
	case reserveOutcomeInFlight:
		out.Outcome = ReservationOutcomeInFlight
	}
	return out
}

// AttachWaiter returns the InFlight reservation's done channel for select.
// ok is false when the namespace is not InFlight. Holds the DD-02 mutex only
// for the table operation.
func (s *ExecutionTargetLifecycleService) AttachWaiter(ns IdempotencyNamespace) (done <-chan struct{}, ok bool) {
	s.lock()
	defer s.unlock()
	return s.idemp.attachWaiter(ns)
}

// DetachWaiter records that a waiter stopped waiting after request
// cancellation. Cancellation detaches only; the InFlight reservation is left
// unchanged. Holds the DD-02 mutex only for the table operation.
func (s *ExecutionTargetLifecycleService) DetachWaiter(ns IdempotencyNamespace) {
	s.lock()
	defer s.unlock()
	s.idemp.detachWaiter(ns)
}

// AwaitReservation selects the reservation done channel or the request
// context. A nil error means the waiter woke on done and must re-enter at
// authentication with no lifecycle mutex held. Context cancellation means the
// waiter must DetachWaiter and must not abort the owner.
func (s *ExecutionTargetLifecycleService) AwaitReservation(ctx context.Context, done <-chan struct{}) error {
	return awaitReservationDone(ctx, done)
}

// AbortReservation removes an InFlight reservation without creating a
// completion, unlocks, then closes the captured done channel exactly once so
// waiters become runnable only after protected state is stable.
func (s *ExecutionTargetLifecycleService) AbortReservation(ns IdempotencyNamespace) bool {
	s.lock()
	done, ok := s.idemp.finalizeAbort(ns)
	s.unlock()
	closeReservationDone(done)
	return ok
}

// CommitCreate validates the closed create-field boundary, derives/persists
// CloudProvider-scoped Active/Unqualified state at epoch 0, appends the
// required create AuditEvent exactly once before publication, publishes, then
// marks the linked idempotency reservation complete. A failed required append
// leaves store state unmutated, removes the linked InFlight reservation,
// creates no completion, wakes waiters only after unlock, and returns
// INTERNAL_ERROR/500.
func (s *ExecutionTargetLifecycleService) CommitCreate(ctx context.Context, req CreateCommitRequest) (model.ExecutionTarget, *apiproblem.Problem) {
	var (
		wake    chan struct{}
		staged  *stagedOutcome
		holding bool
	)
	defer func() {
		if r := recover(); r != nil {
			if !holding {
				s.lock()
				holding = true
			}
			if staged != nil {
				s.store.abort(staged)
			}
			wake, _ = s.idemp.finalizeAbort(req.Idempotency)
			if holding {
				s.unlock()
				holding = false
			}
			closeReservationDone(wake)
			panic(r)
		}
	}()

	if err := ctx.Err(); err != nil {
		s.abortLinkedReservation(req.Idempotency)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeInternalError).WithDetail("create cancelled before publication")
	}

	if prob := validateClosedCreate(req); prob != nil {
		s.abortLinkedReservation(req.Idempotency)
		return model.ExecutionTarget{}, prob
	}

	et := model.NewCreateProposal(
		req.Name,
		req.UID,
		req.CloudProviderScopeUID,
		req.ParticipationRef,
		req.StackRef,
		req.StackGeneration,
	)
	// NewCreateProposal fixes TargetClass; reject if caller attempted another.
	if req.TargetClass != "" && req.TargetClass != model.TargetClassSyntheticIaaS {
		s.abortLinkedReservation(req.Idempotency)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/targetClass",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "targetClass must be synthetic-iaas",
		}})
	}

	s.lock()
	holding = true
	var prob *apiproblem.Problem
	staged, prob = s.store.stageCreateExecutionTarget(et)
	if prob != nil {
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, prob
	}

	event := AssembleSuccessfulCreate(RedactedAuditInput{
		UID:                    req.AuditUID,
		RequestID:              req.RequestID,
		Actor:                  req.Actor,
		Subject:                subjectRef(et),
		SubjectResourceVersion: et.Metadata.ResourceVersion,
	})

	if s.audit == nil {
		s.store.abort(staged)
		staged = nil
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, internalAuditError("audit appender is not configured")
	}
	if err := s.audit.Append(ctx, event); err != nil {
		s.store.abort(staged)
		staged = nil
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, internalAuditError("required AuditEvent append failed")
	}

	s.store.publish(staged)
	staged = nil
	completion := req.Completion
	if completion.StatusCode == 0 {
		completion.StatusCode = http.StatusCreated
	}
	wake, _ = s.idemp.finalizeComplete(req.Idempotency, req.Digest, completion, s.now())
	out, _ := s.store.lookupExecutionTarget(et.Metadata.UID)
	s.unlock()
	holding = false
	closeReservationDone(wake)
	return out, nil
}

// CommitRetire retires an Active ExecutionTarget (including while Maintenance
// is active). It appends the required retirement AuditEvent exactly once
// before publication, clears current FactSet/Result links, releases the live
// tuple, retains the scope-name reservation, publishes Retired/Unqualified,
// then marks the linked idempotency reservation complete.
//
// When a Qualifying reservation is in flight for the same target, a successful
// retirement applies Retire-wins: it aborts that reservation (no qualification
// result, completion, or qualification AuditEvent) and the qualifying caller's
// outcome is CONFLICT/409 + VS0_TARGET_RETIRED. A failed required append
// preserves existing FactSet/Result links, the live tuple, and any in-flight
// Qualifying reservation, removes the linked InFlight retire idempotency
// reservation, creates no completion, wakes waiters only after unlock, and
// returns INTERNAL_ERROR/500.
func (s *ExecutionTargetLifecycleService) CommitRetire(ctx context.Context, req RetireCommitRequest) (model.ExecutionTarget, *apiproblem.Problem) {
	var (
		wake     chan struct{}
		qualWake chan struct{}
		staged   *stagedOutcome
		holding  bool
	)
	defer func() {
		if r := recover(); r != nil {
			if !holding {
				s.lock()
				holding = true
			}
			if staged != nil {
				s.store.abort(staged)
			}
			wake, _ = s.idemp.finalizeAbort(req.Idempotency)
			if holding {
				s.unlock()
				holding = false
			}
			closeReservationDone(wake)
			panic(r)
		}
	}()

	if err := ctx.Err(); err != nil {
		s.abortLinkedReservation(req.Idempotency)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeInternalError).WithDetail("retire cancelled before publication")
	}

	s.lock()
	holding = true
	cur, ok := s.store.lookupExecutionTarget(req.TargetUID)
	if !ok {
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeResourceNotFound)
	}
	if cur.Status.Lifecycle != model.LifecycleActive {
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
			Field:   "/status/lifecycle",
			Code:    ViolationTargetRetired,
			Message: "ExecutionTarget is not Active",
		}})
	}
	if req.ExpectedResourceVersion != "" && cur.Metadata.ResourceVersion != req.ExpectedResourceVersion {
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeStaleResourceVersion)
	}

	// Preserve FactSet/Result for append-failure proofs; clear only on publish.
	next := cur
	next.Status.Lifecycle = model.LifecycleRetired
	next.Status.Qualification = model.QualificationUnqualified
	next.Status.FactSetRef = nil
	next.Status.QualificationResultRef = nil

	var prob *apiproblem.Problem
	staged, prob = s.store.stageUpdateExecutionTarget(next)
	if prob != nil {
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, prob
	}

	// Subject resourceVersion is the pre-publication version (current ETag).
	event := AssembleRetirement(RedactedAuditInput{
		UID:                    req.AuditUID,
		RequestID:              req.RequestID,
		Actor:                  req.Actor,
		Subject:                subjectRef(cur),
		SubjectResourceVersion: cur.Metadata.ResourceVersion,
	})

	if s.audit == nil {
		s.store.abort(staged)
		staged = nil
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, internalAuditError("audit appender is not configured")
	}
	if err := s.audit.Append(ctx, event); err != nil {
		// Failed append must leave any in-flight Qualifying reservation intact.
		s.store.abort(staged)
		staged = nil
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, internalAuditError("required AuditEvent append failed")
	}

	// Retire-wins: abort in-flight Qualifying only after successful audit.
	qualWake = s.abortQualifyingLocked(req.TargetUID, retiredWhileQualifyingProblem())

	s.store.publish(staged)
	staged = nil
	completion := req.Completion
	if completion.StatusCode == 0 {
		completion.StatusCode = http.StatusOK
	}
	wake, _ = s.idemp.finalizeComplete(req.Idempotency, req.Digest, completion, s.now())
	out, _ := s.store.lookupExecutionTarget(req.TargetUID)
	s.unlock()
	holding = false
	closeReservationDone(wake)
	closeReservationDone(qualWake)
	return out, nil
}

// Qualify runs the synchronous qualification path (DD-04): under-mutex backing
// viability recheck, internal Qualifying reservation and fence capture, unlocked
// observation, then a separate relocked commit dispatch with the ADH-2026-060
// ordered predicate and required AuditEvent append before publication.
func (s *ExecutionTargetLifecycleService) Qualify(ctx context.Context, req QualifyCommitRequest) (model.ExecutionTarget, *apiproblem.Problem) {
	var (
		wake    chan struct{}
		res     *qualifyingReservation
		holding bool
	)
	defer func() {
		if r := recover(); r != nil {
			if !holding {
				s.lock()
				holding = true
			}
			wake = s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
			if holding {
				s.unlock()
				holding = false
			}
			closeReservationDone(wake)
			panic(r)
		}
	}()

	if err := ctx.Err(); err != nil {
		s.abortLinkedReservation(req.Idempotency)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeInternalError).WithDetail("qualify cancelled before observation")
	}

	s.lock()
	holding = true

	if s.stopped {
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeInternalError).WithDetail("qualification stopped")
	}

	if prob := recheckBackingViability(req.Backing); prob != nil {
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, prob
	}

	cur, ok := s.store.lookupExecutionTarget(req.TargetUID)
	if !ok {
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeResourceNotFound)
	}
	if cur.Status.Lifecycle != model.LifecycleActive {
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
			Field:   "/status/lifecycle",
			Code:    ViolationTargetRetired,
			Message: "ExecutionTarget is not Active",
		}})
	}
	if _, busy := s.qualifying[req.TargetUID]; busy {
		wake, _ = s.idemp.finalizeAbort(req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
			Field:   "/status/qualification",
			Code:    ViolationTargetQualificationInProgress,
			Message: "qualification already in progress",
		}})
	}

	res = &qualifyingReservation{
		targetUID: req.TargetUID,
		fences: capturedQualifyFences{
			stackGeneration:      req.Backing.StackGeneration,
			maintenanceEpoch:     cur.Status.MaintenanceEpoch,
			viabilityFingerprint: req.Backing.Fingerprint(),
		},
		idempotency: req.Idempotency,
		digest:      req.Digest,
	}
	s.qualifying[req.TargetUID] = res
	s.unlock()
	holding = false

	// Unlocked observation (external-effect-free).
	obs := s.observer.Observe(req.TargetUID)

	if err := ctx.Err(); err != nil {
		s.lock()
		holding = true
		wake = s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeInternalError).WithDetail("qualify cancelled after observation")
	}

	eval, evalErr := EvaluateObservation(obs, s.now())

	s.lock()
	holding = true

	if s.stopped {
		wake = s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeInternalError).WithDetail("qualification stopped")
	}
	if res.aborted != nil {
		// Winning transition already removed the reservation from the map and
		// will close the linked idempotency done channel after its unlock.
		s.unlock()
		holding = false
		return model.ExecutionTarget{}, res.aborted
	}
	if s.qualifying[req.TargetUID] != res {
		wake = s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeInternalError).WithDetail("qualifying reservation lost")
	}
	if evalErr != nil {
		wake = s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		s.unlock()
		holding = false
		closeReservationDone(wake)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeInternalError).WithDetail("observer fault")
	}

	out, prob, wake := s.commitQualifyLocked(ctx, req, res, obs, eval)
	s.unlock()
	holding = false
	closeReservationDone(wake)
	return out, prob
}

// commitQualifyLocked applies the ADH-2026-060 ordered commit predicate and,
// on success, appends the required AuditEvent before atomic publication.
// Caller must hold the sole-committer mutex. The returned done channel must be
// closed only after unlock.
func (s *ExecutionTargetLifecycleService) commitQualifyLocked(
	ctx context.Context,
	req QualifyCommitRequest,
	res *qualifyingReservation,
	obs Observation,
	eval ProfileEvaluation,
) (model.ExecutionTarget, *apiproblem.Problem, chan struct{}) {
	cur, ok := s.store.lookupExecutionTarget(req.TargetUID)
	if !ok || cur.Status.Lifecycle != model.LifecycleActive {
		wake := s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		return model.ExecutionTarget{}, retiredWhileQualifyingProblem(), wake
	}

	// ADH-2026-060 ordered commit predicate.
	if marker, hasMarker := s.store.lookupMaintenanceMarker(req.TargetUID); hasMarker && marker.Active {
		wake := s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		return model.ExecutionTarget{}, maintenanceWhileQualifyingProblem(), wake
	}
	if cur.Status.MaintenanceEpoch != res.fences.maintenanceEpoch {
		wake := s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		return model.ExecutionTarget{}, epochStaleProblem(), wake
	}

	currentBacking := req.Backing
	if req.RefreshBacking != nil {
		currentBacking = req.RefreshBacking()
	}
	if currentBacking.StackGeneration != res.fences.stackGeneration {
		wake := s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		return model.ExecutionTarget{}, epochStaleProblem(), wake
	}
	if currentBacking.Fingerprint() != res.fences.viabilityFingerprint {
		wake := s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		return model.ExecutionTarget{}, viabilityStaleProblem(), wake
	}

	if req.FactSetUID == "" || req.ResultUID == "" {
		wake := s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeInternalError).WithDetail("fact set and result uids are required"), wake
	}

	fs := BuildFactSet(
		req.FactSetUID,
		cur,
		obs,
		eval,
		res.fences.stackGeneration,
		res.fences.maintenanceEpoch,
		res.fences.viabilityFingerprint,
	)
	result := BuildQualificationResult(
		req.ResultUID,
		cur,
		req.FactSetUID,
		eval,
		res.fences.stackGeneration,
		res.fences.maintenanceEpoch,
	)

	next := cur
	next.Status.Qualification = qualificationStateForOutcome(eval.Outcome)
	next.Status.ObservedGeneration = res.fences.stackGeneration
	fsRef := apimeta.TypedRef{
		APIVersion: model.APIVersionFactSet,
		Kind:       model.KindNormalizedTargetFactSet,
		Name:       fs.UID,
		UID:        fs.UID,
	}
	resRef := apimeta.TypedRef{
		APIVersion: model.APIVersionQualificationResult,
		Kind:       model.KindTargetQualificationResult,
		Name:       result.UID,
		UID:        result.UID,
	}
	next.Status.FactSetRef = &fsRef
	next.Status.QualificationResultRef = &resRef

	stagedFS, prob := s.store.stagePutFactSet(fs)
	if prob != nil {
		wake := s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		return model.ExecutionTarget{}, prob, wake
	}
	stagedRes, prob := s.store.stagePutQualificationResult(result)
	if prob != nil {
		s.store.abort(stagedFS)
		wake := s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		return model.ExecutionTarget{}, prob, wake
	}
	stagedET, prob := s.store.stageUpdateExecutionTarget(next)
	if prob != nil {
		s.store.abort(stagedFS)
		s.store.abort(stagedRes)
		wake := s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		return model.ExecutionTarget{}, prob, wake
	}

	event := AssembleCompletedQualification(RedactedAuditInput{
		UID:                    req.AuditUID,
		RequestID:              req.RequestID,
		Actor:                  req.Actor,
		Subject:                subjectRef(cur),
		SubjectResourceVersion: cur.Metadata.ResourceVersion,
	})

	if s.audit == nil {
		s.store.abort(stagedFS)
		s.store.abort(stagedRes)
		s.store.abort(stagedET)
		wake := s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		return model.ExecutionTarget{}, internalAuditError("audit appender is not configured"), wake
	}
	if err := s.audit.Append(ctx, event); err != nil {
		s.store.abort(stagedFS)
		s.store.abort(stagedRes)
		s.store.abort(stagedET)
		wake := s.abortQualifyAndIdempotencyLocked(res, req.Idempotency)
		return model.ExecutionTarget{}, internalAuditError("required AuditEvent append failed"), wake
	}

	s.store.publish(stagedFS)
	s.store.publish(stagedRes)
	s.store.publish(stagedET)
	delete(s.qualifying, req.TargetUID)

	completion := req.Completion
	if completion.StatusCode == 0 {
		completion.StatusCode = http.StatusOK
	}
	wake, _ := s.idemp.finalizeComplete(req.Idempotency, req.Digest, completion, s.now())
	out, _ := s.store.lookupExecutionTarget(req.TargetUID)
	return out, nil, wake
}

// MarkQualificationStopped prevents new Qualifying work and aborts in-flight
// Qualifying reservations together with their linked InFlight idempotency
// reservations. TASK-F16-06/11 own ordered shutdown; this is the abort hook.
func (s *ExecutionTargetLifecycleService) MarkQualificationStopped() {
	s.lock()
	s.stopped = true
	var wakes []chan struct{}
	for uid, qres := range s.qualifying {
		qres.aborted = apiproblem.New(apiproblem.CodeInternalError).WithDetail("qualification stopped")
		done, _ := s.idemp.finalizeAbort(qres.idempotency)
		if done != nil {
			wakes = append(wakes, done)
		}
		delete(s.qualifying, uid)
	}
	s.unlock()
	for _, done := range wakes {
		closeReservationDone(done)
	}
}

// Observer returns the synthetic observer bound to this lifecycle service.
func (s *ExecutionTargetLifecycleService) Observer() *SyntheticObserver {
	return s.observer
}

// abortLinkedReservation finalizes/detaches an InFlight reservation without
// publication or completion, unlocks, then closes done exactly once.
func (s *ExecutionTargetLifecycleService) abortLinkedReservation(ns IdempotencyNamespace) {
	s.lock()
	done, _ := s.idemp.finalizeAbort(ns)
	s.unlock()
	closeReservationDone(done)
}

// abortQualifyingLocked marks and removes an in-flight Qualifying reservation
// for targetUID, aborts its linked idempotency reservation, and returns the
// captured done channel for post-unlock close. Caller must hold the mutex.
func (s *ExecutionTargetLifecycleService) abortQualifyingLocked(targetUID string, prob *apiproblem.Problem) chan struct{} {
	qres, ok := s.qualifying[targetUID]
	if !ok {
		return nil
	}
	qres.aborted = prob
	delete(s.qualifying, targetUID)
	done, _ := s.idemp.finalizeAbort(qres.idempotency)
	return done
}

// abortQualifyAndIdempotencyLocked clears the Qualifying reservation (when
// still owned) and removes the linked InFlight idempotency reservation.
// Caller must hold the mutex. Returns the done channel for post-unlock close.
func (s *ExecutionTargetLifecycleService) abortQualifyAndIdempotencyLocked(qres *qualifyingReservation, ns IdempotencyNamespace) chan struct{} {
	if qres != nil {
		if s.qualifying[qres.targetUID] == qres {
			delete(s.qualifying, qres.targetUID)
		}
	}
	done, _ := s.idemp.finalizeAbort(ns)
	return done
}

func recheckBackingViability(b BackingViability) *apiproblem.Problem {
	if !b.ParticipationEffectiveActive {
		return apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/cloudProviderParticipationRef",
			Code:    ViolationExecutionTargetParticipationUnavailable,
			Message: "CloudProviderParticipation is not effective Active",
		}})
	}
	if !b.StackActive() {
		return apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/infrastructureStackRef",
			Code:    ViolationExecutionTargetStackUnavailable,
			Message: "InfrastructureStack is not Active",
		}})
	}
	if b.ParticipationScopeUID == "" || b.StackScopeUID == "" || b.ParticipationScopeUID != b.StackScopeUID {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec",
			Code:    ViolationExecutionTargetScopeMismatch,
			Message: "participation and stack CloudProvider scopes do not match",
		}})
	}
	return nil
}

func retiredWhileQualifyingProblem() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
		Field:   "/status/lifecycle",
		Code:    ViolationTargetRetired,
		Message: "ExecutionTarget was retired while qualification was in flight",
	}})
}

func maintenanceWhileQualifyingProblem() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeConflict).WithViolations([]apiproblem.Violation{{
		Field:   "/status/lifecycle",
		Code:    ViolationTargetMaintenance,
		Message: "ExecutionTarget entered Maintenance while qualification was in flight",
	}})
}

func epochStaleProblem() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeStaleResourceVersion).WithViolations([]apiproblem.Violation{{
		Field:   "/status/maintenanceEpoch",
		Code:    ViolationTargetEpochStale,
		Message: "captured maintenance epoch or InfrastructureStack generation is stale",
	}})
}

func viabilityStaleProblem() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeStaleResourceVersion).WithViolations([]apiproblem.Violation{{
		Field:   "/status",
		Code:    ViolationExecutionTargetViabilityStale,
		Message: "captured viability fingerprint is stale",
	}})
}

func validateClosedCreate(req CreateCommitRequest) *apiproblem.Problem {
	if req.ClientStatusAttempted {
		return apiproblem.New(apiproblem.CodeAuthorizationDenied).WithViolations([]apiproblem.Violation{{
			Field:   "/status",
			Code:    "VS0_STATUS_FIELD_WRITE",
			Message: "clients never write status",
		}})
	}
	if req.Name == "" {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/metadata/name",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "metadata.name is required",
		}})
	}
	if req.UID == "" {
		return apiproblem.New(apiproblem.CodeInternalError).WithDetail("server-assigned uid is required before create commit")
	}
	if req.CloudProviderScopeUID == "" {
		return apiproblem.New(apiproblem.CodeInternalError).WithDetail("derived CloudProvider scope is required before create commit")
	}
	if req.ParticipationRef.UID == "" {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/cloudProviderParticipationRef/uid",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "participation uid is required",
		}})
	}
	if req.StackRef.UID == "" {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/infrastructureStackRef/uid",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "stack uid is required",
		}})
	}
	if req.TargetClass != "" && req.TargetClass != model.TargetClassSyntheticIaaS {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/targetClass",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "targetClass must be synthetic-iaas",
		}})
	}
	return nil
}

func subjectRef(et model.ExecutionTarget) apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: model.APIVersionExecutionTarget,
		Kind:       model.KindExecutionTarget,
		Name:       et.Metadata.Name,
		UID:        et.Metadata.UID,
	}
}

func internalAuditError(detail string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeInternalError).WithDetail(detail)
}
