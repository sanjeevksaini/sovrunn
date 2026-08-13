package cloudmodel

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

// Clock provides the current time for the deterministic expiry scheduler.
type Clock interface {
	Now() time.Time
}

// Timer is a resettable one-shot wait used by the scheduler loop.
type Timer interface {
	C() <-chan time.Time
	Stop() bool
	Reset(d time.Duration) bool
}

// TimerFactory constructs a Timer for a duration. Tests inject a controllable
// factory; production uses real timers.
type TimerFactory func(d time.Duration) Timer

// ExecutionIDGenerator produces a scheduler execution ID for each expiry attempt.
type ExecutionIDGenerator func() string

// ExpiryOperation is the api-server participation-expiry transition invoked by
// the scheduler. Task 8 supplies the audit-aware implementation. The scheduler
// never writes status or publishes expiry itself.
//
// The participation action concurrency guard is held across recheck and this
// call. The operation must not re-acquire that guard.
type ExpiryOperation func(ctx context.Context, uid, observedVersion, schedulerExecutionID string) error

// realClock is the production wall clock.
type realClock struct{}

func (realClock) Now() time.Time { return time.Now().UTC() }

type stdTimer struct {
	t *time.Timer
}

func (t *stdTimer) C() <-chan time.Time { return t.t.C }
func (t *stdTimer) Stop() bool          { return t.t.Stop() }
func (t *stdTimer) Reset(d time.Duration) bool {
	if !t.t.Stop() {
		select {
		case <-t.t.C:
		default:
		}
	}
	return t.t.Reset(d)
}

func stdTimerFactory(d time.Duration) Timer {
	return &stdTimer{t: time.NewTimer(d)}
}

// ParticipationExpiryScheduler wakes at the earliest Pending
// status.requestExpiresAt and invokes an injected expiry operation
// (design §4.6). It is not a public route or controller.
type ParticipationExpiryScheduler struct {
	store    *Store
	guard    *ParticipationActionGuard
	clock    Clock
	newTimer TimerFactory
	expire   ExpiryOperation
	genID    ExecutionIDGenerator

	signal chan struct{}
	stop   chan struct{}

	mu      sync.Mutex
	running bool
	stopped chan struct{}
}

// ParticipationExpirySchedulerConfig configures the scheduler core.
type ParticipationExpirySchedulerConfig struct {
	Store     *Store
	Guard     *ParticipationActionGuard
	Clock     Clock
	NewTimer  TimerFactory
	Expire    ExpiryOperation
	NewExecID ExecutionIDGenerator
}

// NewParticipationExpiryScheduler constructs a stopped scheduler core.
func NewParticipationExpiryScheduler(cfg ParticipationExpirySchedulerConfig) *ParticipationExpiryScheduler {
	clock := cfg.Clock
	if clock == nil {
		clock = realClock{}
	}
	newTimer := cfg.NewTimer
	if newTimer == nil {
		newTimer = stdTimerFactory
	}
	genID := cfg.NewExecID
	if genID == nil {
		genID = func() string { return "sched-" + time.Now().UTC().Format("20060102T150405.000000000") }
	}
	guard := cfg.Guard
	if guard == nil {
		guard = &ParticipationActionGuard{}
	}
	return &ParticipationExpiryScheduler{
		store:    cfg.Store,
		guard:    guard,
		clock:    clock,
		newTimer: newTimer,
		expire:   cfg.Expire,
		genID:    genID,
		signal:   make(chan struct{}, 1),
		stop:     make(chan struct{}),
	}
}

// Signal wakes the scheduler to recompute the next due time after a create
// or successful lifecycle change.
func (s *ParticipationExpiryScheduler) Signal() {
	if s == nil {
		return
	}
	select {
	case s.signal <- struct{}{}:
	default:
	}
}

// Start runs the scheduler loop until Stop. It is safe to call once.
func (s *ParticipationExpiryScheduler) Start(ctx context.Context) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.stopped = make(chan struct{})
	s.mu.Unlock()

	go s.loop(ctx)
}

// Stop requests shutdown and waits for any active tick to finish.
func (s *ParticipationExpiryScheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	select {
	case <-s.stop:
	default:
		close(s.stop)
	}
	s.mu.Lock()
	done := s.stopped
	s.mu.Unlock()
	if done != nil {
		<-done
	}
}

// EarliestPendingExpiresAt returns the earliest Pending requestExpiresAt, or
// false when no Pending participation exists.
func (s *ParticipationExpiryScheduler) EarliestPendingExpiresAt() (time.Time, bool) {
	return s.store.earliestPendingExpiresAt()
}

// DuePendingUIDs returns Pending participation UIDs whose requestExpiresAt is
// at or before now, in stable ascending UID order.
func (s *ParticipationExpiryScheduler) DuePendingUIDs(now time.Time) []string {
	return s.store.duePendingUIDs(now)
}

func (s *ParticipationExpiryScheduler) loop(ctx context.Context) {
	defer func() {
		s.mu.Lock()
		s.running = false
		if s.stopped != nil {
			close(s.stopped)
		}
		s.mu.Unlock()
	}()

	var timer Timer
	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

	for {
		dueAt, ok := s.EarliestPendingExpiresAt()
		now := s.clock.Now()

		if !ok {
			select {
			case <-ctx.Done():
				return
			case <-s.stop:
				return
			case <-s.signal:
				continue
			}
		}

		wait := dueAt.Sub(now)
		if wait < 0 {
			wait = 0
		}

		if timer == nil {
			timer = s.newTimer(wait)
		} else {
			timer.Reset(wait)
		}

		select {
		case <-ctx.Done():
			return
		case <-s.stop:
			return
		case <-s.signal:
			continue
		case <-timer.C():
			s.processDue(ctx, s.clock.Now())
		}
	}
}

func (s *ParticipationExpiryScheduler) processDue(ctx context.Context, now time.Time) {
	uids := s.DuePendingUIDs(now)
	for _, uid := range uids {
		if err := ctx.Err(); err != nil {
			return
		}
		select {
		case <-s.stop:
			return
		default:
		}
		s.attemptExpiry(ctx, uid, now)
	}
}

func (s *ParticipationExpiryScheduler) attemptExpiry(ctx context.Context, uid string, now time.Time) {
	s.guard.Lock()
	defer s.guard.Unlock()

	p, ok := s.store.GetParticipation(uid)
	if !ok {
		return
	}
	// Recheck Pending, still expired, and capture observed version.
	if p.Status.Phase != model.ParticipationPhasePending {
		return
	}
	exp, err := time.Parse(time.RFC3339, p.Status.RequestExpiresAt)
	if err != nil {
		return
	}
	if exp.After(now) {
		return
	}
	observedVersion := p.Metadata.ResourceVersion
	execID := s.genID()
	if s.expire == nil {
		return
	}
	// Injected operation owns audit-before-publication. Failures leave the
	// participation unpublished; the scheduler does not write status.
	_ = s.expire(ctx, uid, observedVersion, execID)
}

func (s *Store) earliestPendingExpiresAt() (time.Time, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var earliest time.Time
	found := false
	for _, p := range s.participations {
		if p.Status.Phase != model.ParticipationPhasePending {
			continue
		}
		exp, err := time.Parse(time.RFC3339, p.Status.RequestExpiresAt)
		if err != nil {
			continue
		}
		if !found || exp.Before(earliest) {
			earliest = exp
			found = true
		}
	}
	return earliest, found
}

func (s *Store) duePendingUIDs(now time.Time) []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var uids []string
	for uid, p := range s.participations {
		if p.Status.Phase != model.ParticipationPhasePending {
			continue
		}
		exp, err := time.Parse(time.RFC3339, p.Status.RequestExpiresAt)
		if err != nil {
			continue
		}
		if !exp.After(now) {
			uids = append(uids, uid)
		}
	}
	sort.Strings(uids)
	return uids
}
