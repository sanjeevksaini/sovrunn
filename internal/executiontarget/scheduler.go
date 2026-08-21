package executiontarget

import (
	"context"
	"sync"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

// MaintenanceTriggerOperation is the closed maintenance fixture-trigger op set.
type MaintenanceTriggerOperation string

const (
	// MaintenanceTriggerEnter enters Maintenance for a fenced target.
	MaintenanceTriggerEnter MaintenanceTriggerOperation = "enter"
	// MaintenanceTriggerClear clears Maintenance for a fenced target.
	MaintenanceTriggerClear MaintenanceTriggerOperation = "clear"
)

// MaintenanceTrigger is the deterministic, target-bound synthetic-observer
// fixture trigger carrying fence, operation, system identity, and correlation
// (DD-07; ADH-2026-058 clause 8). It is never an HTTP route or event bus.
type MaintenanceTrigger struct {
	TargetUID                      string
	ExpectedInfrastructureStackGen int64
	ExpectedMaintenanceEpoch       int64
	Operation                      MaintenanceTriggerOperation
	Actor                          apimeta.TypedRef
	RequestID                      string
	AuditUID                       string
}

// Scheduler owns injected-clock expiry ticks and target-bound maintenance
// trigger intake (design §3.1 scheduler.go; DD-06; DD-07).
type Scheduler struct {
	lifecycle *ExecutionTargetLifecycleService
	clock     Clock

	mu                sync.Mutex
	acceptanceStopped bool
	running           bool
	stopCh            chan struct{}
	stoppedCh         chan struct{}
	tickEvery         time.Duration
	externalCallCount int64 // always remains 0 (Risk 2)
}

// SchedulerConfig configures the expiry/maintenance scheduler.
type SchedulerConfig struct {
	Lifecycle *ExecutionTargetLifecycleService
	Clock     Clock
	// TickEvery defaults to one second for the production loop. Tests typically
	// call ProcessExpiryTick directly against a FakeClock.
	TickEvery time.Duration
}

// NewScheduler constructs a stopped scheduler bound to the lifecycle service.
func NewScheduler(cfg SchedulerConfig) *Scheduler {
	clock := cfg.Clock
	if clock == nil {
		clock = MonotonicClock{}
	}
	tick := cfg.TickEvery
	if tick <= 0 {
		tick = time.Second
	}
	return &Scheduler{
		lifecycle: cfg.Lifecycle,
		clock:     clock,
		tickEvery: tick,
		stopCh:    make(chan struct{}),
	}
}

// Clock returns the injected clock.
func (s *Scheduler) Clock() Clock {
	return s.clock
}

// ExternalCallCount remains 0: maintenance triggers and expiry never perform
// network I/O (Risk 2).
func (s *Scheduler) ExternalCallCount() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.externalCallCount
}

// AcceptanceStopped reports whether expiry and Maintenance trigger acceptance
// has been closed (shutdown step 1).
func (s *Scheduler) AcceptanceStopped() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.acceptanceStopped
}

// DeliverMaintenanceTrigger dispatches a target-bound fixture trigger to the
// lifecycle maintenance commit methods. After StopAcceptance/Stop, triggers
// produce no reservation, mutation, or AuditEvent.
func (s *Scheduler) DeliverMaintenanceTrigger(ctx context.Context, trigger MaintenanceTrigger) (model.ExecutionTarget, *apiproblem.Problem) {
	s.mu.Lock()
	stopped := s.acceptanceStopped
	s.mu.Unlock()
	if stopped || s.lifecycle == nil {
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeInternalError).WithDetail("maintenance trigger acceptance stopped")
	}

	req := MaintenanceCommitRequest{
		TargetUID:                      trigger.TargetUID,
		ExpectedInfrastructureStackGen: trigger.ExpectedInfrastructureStackGen,
		ExpectedMaintenanceEpoch:       trigger.ExpectedMaintenanceEpoch,
		Actor:                          trigger.Actor,
		RequestID:                      trigger.RequestID,
		AuditUID:                       trigger.AuditUID,
	}
	switch trigger.Operation {
	case MaintenanceTriggerEnter:
		return s.lifecycle.CommitMaintenanceEnter(ctx, req)
	case MaintenanceTriggerClear:
		return s.lifecycle.CommitMaintenanceClear(ctx, req)
	default:
		return model.ExecutionTarget{}, apiproblem.New(apiproblem.CodeValidationFailed).WithDetail("unknown maintenance trigger operation")
	}
}

// ProcessExpiryTick runs one expiry pass at the injected clock's current time.
// Deterministic tests Advance the FakeClock then call this once per
// injected-clock second.
func (s *Scheduler) ProcessExpiryTick(ctx context.Context) {
	s.mu.Lock()
	stopped := s.acceptanceStopped
	s.mu.Unlock()
	if stopped || s.lifecycle == nil {
		return
	}
	s.lifecycle.ProcessExpiryTick(ctx, s.clock.Now().UTC())
}

// StopAcceptance closes expiry and Maintenance trigger acceptance without
// aborting reservations. TASK-F16-11 orders: StopAcceptance → lifecycle
// Shutdown → HTTP shutdown.
func (s *Scheduler) StopAcceptance() {
	s.mu.Lock()
	s.acceptanceStopped = true
	s.mu.Unlock()
}

// Start runs the once-per-tickEvery expiry loop until Stop. Safe to call once.
func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	if s.running || s.acceptanceStopped {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.stoppedCh = make(chan struct{})
	s.mu.Unlock()
	go s.loop(ctx)
}

// Stop closes trigger acceptance, stops the expiry loop, and waits for the
// active tick to finish.
func (s *Scheduler) Stop() {
	s.StopAcceptance()
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()

	select {
	case <-s.stopCh:
	default:
		close(s.stopCh)
	}
	s.mu.Lock()
	done := s.stoppedCh
	s.mu.Unlock()
	if done != nil {
		<-done
	}
}

func (s *Scheduler) loop(ctx context.Context) {
	defer func() {
		s.mu.Lock()
		s.running = false
		if s.stoppedCh != nil {
			close(s.stoppedCh)
		}
		s.mu.Unlock()
	}()

	ticker := time.NewTicker(s.tickEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.ProcessExpiryTick(ctx)
		}
	}
}

// MaintenanceSystemActor returns the deterministic system actor TypedRef used
// for maintenance fixture triggers.
func MaintenanceSystemActor() apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: model.APIVersionExecutionTarget,
		Kind:       "SystemActor",
		Name:       "executiontarget-maintenance-trigger",
		UID:        "system:executiontarget-maintenance-trigger",
	}
}
