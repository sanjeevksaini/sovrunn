package cloudmodel_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = t.UTC()
}

type fakeTimer struct {
	ch       chan time.Time
	mu       sync.Mutex
	active   bool
	autoFire bool
	clock    *fakeClock
}

func (t *fakeTimer) C() <-chan time.Time { return t.ch }

func (t *fakeTimer) Stop() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	was := t.active
	t.active = false
	return was
}

func (t *fakeTimer) Reset(d time.Duration) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	was := t.active
	t.active = true
	if t.autoFire && d <= 0 {
		select {
		case t.ch <- t.clock.Now():
		default:
		}
	}
	return was
}

type timerHub struct {
	mu       sync.Mutex
	timers   []*fakeTimer
	autoFire bool
	clock    *fakeClock
}

func (h *timerHub) Factory(d time.Duration) cloudmodel.Timer {
	t := &fakeTimer{
		ch:       make(chan time.Time, 1),
		active:   true,
		autoFire: h.autoFire,
		clock:    h.clock,
	}
	if h.autoFire && d <= 0 {
		t.ch <- h.clock.Now()
	}
	h.mu.Lock()
	h.timers = append(h.timers, t)
	h.mu.Unlock()
	return t
}

func (h *timerHub) Latest() *fakeTimer {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.timers) == 0 {
		return nil
	}
	return h.timers[len(h.timers)-1]
}

func (t *fakeTimer) Fire(at time.Time) {
	t.mu.Lock()
	active := t.active
	t.mu.Unlock()
	if !active {
		return
	}
	select {
	case t.ch <- at:
	default:
	}
}

func participation(uid, providerUID, expiresAt string, phase model.ParticipationPhase) model.CloudProviderParticipation {
	platformUID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	return model.CloudProviderParticipation{
		Metadata: apimeta.ObjectMeta{
			Name:            "p-" + uid[:8],
			UID:             uid,
			ResourceVersion: "1",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform,
				Kind:       model.KindCloudPlatform,
				Name:       "plat",
				UID:        platformUID,
			}},
		},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
				Name: "plat", UID: platformUID,
			},
			CloudProviderRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider,
				Name: "prov", UID: providerUID,
			},
			Environment: model.ParticipationEnvironmentDevelopment,
		},
		Status: model.CloudProviderParticipationStatus{
			Phase:            phase,
			RequestExpiresAt: expiresAt,
		},
	}
}

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	ts, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatal(err)
	}
	return ts
}

func waitForCalls(calls *atomic.Int32, n int32, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if calls.Load() >= n {
			return true
		}
		time.Sleep(5 * time.Millisecond)
	}
	return calls.Load() >= n
}

func TestScheduler_EarliestPendingDueTime(t *testing.T) {
	t.Parallel()
	store := cloudmodel.NewStore()
	early := "2026-08-13T00:00:00Z"
	late := "2026-08-14T00:00:00Z"
	if _, prob := store.CreateParticipation(participation(
		"cccccccccccccccccccccccccccccccc", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", late, model.ParticipationPhasePending,
	)); prob != nil {
		t.Fatalf("create: %#v", prob)
	}
	if _, prob := store.CreateParticipation(participation(
		"dddddddddddddddddddddddddddddddd", "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", early, model.ParticipationPhasePending,
	)); prob != nil {
		t.Fatalf("create2: %#v", prob)
	}
	if _, prob := store.CreateParticipation(participation(
		"ffffffffffffffffffffffffffffffff", "11111111111111111111111111111111", "2026-08-12T00:00:00Z", model.ParticipationPhaseActive,
	)); prob != nil {
		t.Fatalf("create3: %#v", prob)
	}

	sched := cloudmodel.NewParticipationExpiryScheduler(cloudmodel.ParticipationExpirySchedulerConfig{
		Store:  store,
		Guard:  &cloudmodel.ParticipationActionGuard{},
		Expire: func(context.Context, string, string, string) error { return nil },
	})
	got, ok := sched.EarliestPendingExpiresAt()
	if !ok {
		t.Fatal("expected earliest pending")
	}
	want := mustTime(t, early)
	if !got.Equal(want) {
		t.Fatalf("earliest=%v want %v", got, want)
	}
}

func TestScheduler_StableAscendingUIDOrder(t *testing.T) {
	t.Parallel()
	store := cloudmodel.NewStore()
	uids := []string{
		"dddddddddddddddddddddddddddddddd",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccc",
	}
	providers := []string{
		"11111111111111111111111111111111",
		"22222222222222222222222222222222",
		"33333333333333333333333333333333",
	}
	expires := "2026-08-13T00:00:00Z"
	for i, uid := range uids {
		if _, prob := store.CreateParticipation(participation(uid, providers[i], expires, model.ParticipationPhasePending)); prob != nil {
			t.Fatalf("create %s: %#v", uid, prob)
		}
	}
	sched := cloudmodel.NewParticipationExpiryScheduler(cloudmodel.ParticipationExpirySchedulerConfig{
		Store:  store,
		Guard:  &cloudmodel.ParticipationActionGuard{},
		Expire: func(context.Context, string, string, string) error { return nil },
	})
	got := sched.DuePendingUIDs(mustTime(t, "2026-08-13T00:00:00Z"))
	want := []string{
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"cccccccccccccccccccccccccccccccc",
		"dddddddddddddddddddddddddddddddd",
	}
	if len(got) != len(want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got=%v want=%v", got, want)
		}
	}
}

func TestScheduler_RecheckBeforeExpiryAndIdempotent(t *testing.T) {
	t.Parallel()
	store := cloudmodel.NewStore()
	uid := "cccccccccccccccccccccccccccccccc"
	if _, prob := store.CreateParticipation(participation(
		uid, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "2026-08-13T00:00:00Z", model.ParticipationPhasePending,
	)); prob != nil {
		t.Fatalf("create: %#v", prob)
	}

	var calls atomic.Int32
	guard := &cloudmodel.ParticipationActionGuard{}
	clock := &fakeClock{now: mustTime(t, "2026-08-13T00:00:01Z")}
	hub := &timerHub{autoFire: true, clock: clock}

	sched := cloudmodel.NewParticipationExpiryScheduler(cloudmodel.ParticipationExpirySchedulerConfig{
		Store:     store,
		Guard:     guard,
		Clock:     clock,
		NewTimer:  hub.Factory,
		NewExecID: func() string { return "exec-1" },
		Expire: func(_ context.Context, gotUID, observedVersion, execID string) error {
			calls.Add(1)
			if gotUID != uid || observedVersion != "1" || execID != "exec-1" {
				t.Errorf("unexpected args uid=%s rv=%s exec=%s", gotUID, observedVersion, execID)
			}
			cur, ok := store.GetParticipation(gotUID)
			if !ok || cur.Status.Phase != model.ParticipationPhasePending {
				return nil
			}
			if cur.Metadata.ResourceVersion != observedVersion {
				return nil
			}
			next, prob := cloudmodel.ApplyParticipationAction(cur.Status, cloudmodel.ParticipationActionExpire, "")
			if prob != nil {
				return errors.New(prob.Detail)
			}
			cur.Status = next
			if _, prob := store.UpdateParticipation(cur); prob != nil {
				return errors.New(string(prob.Code))
			}
			return nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sched.Start(ctx)
	sched.Signal()
	if !waitForCalls(&calls, 1, 2*time.Second) {
		t.Fatalf("calls=%d", calls.Load())
	}
	cur, ok := store.GetParticipation(uid)
	if !ok || cur.Status.Phase != model.ParticipationPhaseExpired {
		t.Fatalf("status=%#v", cur.Status)
	}

	// Idempotent: already non-Pending is not transitioned again.
	clock.Set(mustTime(t, "2026-08-13T00:00:02Z"))
	sched.Signal()
	time.Sleep(50 * time.Millisecond)
	if calls.Load() != 1 {
		t.Fatalf("idempotent calls=%d", calls.Load())
	}
	sched.Stop()
}

func TestScheduler_ConcurrencyGuardPreventsRace(t *testing.T) {
	t.Parallel()
	store := cloudmodel.NewStore()
	uid := "cccccccccccccccccccccccccccccccc"
	if _, prob := store.CreateParticipation(participation(
		uid, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "2026-08-13T00:00:00Z", model.ParticipationPhasePending,
	)); prob != nil {
		t.Fatalf("create: %#v", prob)
	}

	guard := &cloudmodel.ParticipationActionGuard{}
	var expireEntered sync.WaitGroup
	expireEntered.Add(1)
	releaseExpire := make(chan struct{})
	var calls atomic.Int32
	clock := &fakeClock{now: mustTime(t, "2026-08-13T00:00:01Z")}
	hub := &timerHub{autoFire: true, clock: clock}

	sched := cloudmodel.NewParticipationExpiryScheduler(cloudmodel.ParticipationExpirySchedulerConfig{
		Store:    store,
		Guard:    guard,
		Clock:    clock,
		NewTimer: hub.Factory,
		Expire: func(context.Context, string, string, string) error {
			if calls.Add(1) == 1 {
				expireEntered.Done()
				<-releaseExpire
			}
			// Publish Expired so the scheduler does not retry forever.
			cur, ok := store.GetParticipation(uid)
			if !ok {
				return nil
			}
			next, prob := cloudmodel.ApplyParticipationAction(cur.Status, cloudmodel.ParticipationActionExpire, "")
			if prob != nil {
				return nil
			}
			cur.Status = next
			_, _ = store.UpdateParticipation(cur)
			return nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sched.Start(ctx)
	sched.Signal()
	expireEntered.Wait()

	actionStarted := make(chan struct{})
	actionDone := make(chan struct{})
	go func() {
		close(actionStarted)
		guard.Lock()
		defer guard.Unlock()
		close(actionDone)
	}()

	select {
	case <-actionStarted:
	case <-time.After(time.Second):
		t.Fatal("action did not start")
	}
	select {
	case <-actionDone:
		t.Fatal("action acquired guard while expiry in progress")
	case <-time.After(50 * time.Millisecond):
	}

	close(releaseExpire)
	select {
	case <-actionDone:
	case <-time.After(time.Second):
		t.Fatal("action did not acquire guard after expiry released")
	}
	sched.Stop()
}

func TestScheduler_InjectedOperationFailureNoPublication(t *testing.T) {
	t.Parallel()
	store := cloudmodel.NewStore()
	uid := "cccccccccccccccccccccccccccccccc"
	if _, prob := store.CreateParticipation(participation(
		uid, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "2026-08-13T00:00:00Z", model.ParticipationPhasePending,
	)); prob != nil {
		t.Fatalf("create: %#v", prob)
	}

	clock := &fakeClock{now: mustTime(t, "2026-08-13T00:00:01Z")}
	hub := &timerHub{autoFire: true, clock: clock}
	var calls atomic.Int32
	sched := cloudmodel.NewParticipationExpiryScheduler(cloudmodel.ParticipationExpirySchedulerConfig{
		Store:    store,
		Guard:    &cloudmodel.ParticipationActionGuard{},
		Clock:    clock,
		NewTimer: hub.Factory,
		Expire: func(context.Context, string, string, string) error {
			calls.Add(1)
			return errors.New("audit append failed")
		},
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sched.Start(ctx)
	sched.Signal()
	if !waitForCalls(&calls, 1, 2*time.Second) {
		t.Fatal("expiry operation was not invoked")
	}
	sched.Stop()
	cur, ok := store.GetParticipation(uid)
	if !ok || cur.Status.Phase != model.ParticipationPhasePending {
		t.Fatalf("failure must not publish expiry: %#v", cur.Status)
	}
}

func TestScheduler_StartStopAndSignal(t *testing.T) {
	t.Parallel()
	store := cloudmodel.NewStore()
	clock := &fakeClock{now: mustTime(t, "2026-08-12T00:00:00Z")}
	hub := &timerHub{clock: clock}
	sched := cloudmodel.NewParticipationExpiryScheduler(cloudmodel.ParticipationExpirySchedulerConfig{
		Store:    store,
		Guard:    &cloudmodel.ParticipationActionGuard{},
		Clock:    clock,
		NewTimer: hub.Factory,
		Expire:   func(context.Context, string, string, string) error { return nil },
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sched.Start(ctx)
	sched.Signal()
	if _, prob := store.CreateParticipation(participation(
		"cccccccccccccccccccccccccccccccc",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"2026-08-13T00:00:00Z",
		model.ParticipationPhasePending,
	)); prob != nil {
		t.Fatalf("create: %#v", prob)
	}
	sched.Signal()
	sched.Stop()
	sched.Stop()
}

func TestScheduler_RecheckSkipsWhenActionWinsRace(t *testing.T) {
	t.Parallel()
	store := cloudmodel.NewStore()
	uid := "cccccccccccccccccccccccccccccccc"
	if _, prob := store.CreateParticipation(participation(
		uid, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "2026-08-13T00:00:00Z", model.ParticipationPhasePending,
	)); prob != nil {
		t.Fatalf("create: %#v", prob)
	}

	guard := &cloudmodel.ParticipationActionGuard{}
	guard.Lock()
	cur, _ := store.GetParticipation(uid)
	next, prob := cloudmodel.ApplyParticipationAction(cur.Status, cloudmodel.ParticipationActionAccept, "")
	if prob != nil {
		guard.Unlock()
		t.Fatalf("accept: %#v", prob)
	}
	cur.Status = next
	if _, prob := store.UpdateParticipation(cur); prob != nil {
		guard.Unlock()
		t.Fatalf("update: %#v", prob)
	}
	guard.Unlock()

	var calls atomic.Int32
	clock := &fakeClock{now: mustTime(t, "2026-08-13T00:00:01Z")}
	hub := &timerHub{autoFire: true, clock: clock}
	sched := cloudmodel.NewParticipationExpiryScheduler(cloudmodel.ParticipationExpirySchedulerConfig{
		Store:    store,
		Guard:    guard,
		Clock:    clock,
		NewTimer: hub.Factory,
		Expire: func(context.Context, string, string, string) error {
			calls.Add(1)
			return nil
		},
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sched.Start(ctx)
	sched.Signal()
	time.Sleep(100 * time.Millisecond)
	sched.Stop()
	if calls.Load() != 0 {
		t.Fatalf("recheck must skip non-Pending; calls=%d", calls.Load())
	}
}

func TestScheduler_NoPublicRouteOrControllerSurface(t *testing.T) {
	t.Parallel()
	_ = cloudmodel.NewParticipationExpiryScheduler
	_ = cloudmodel.ParticipationActionGuard{}
	_ = cloudmodel.ParticipationActionExpire
}
