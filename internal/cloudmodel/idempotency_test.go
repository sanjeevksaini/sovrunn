package cloudmodel_test

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
)

func testNS(principal, pattern, target, key string) cloudmodel.IdempotencyNamespace {
	return cloudmodel.IdempotencyNamespace{
		PrincipalID: principal,
		Pattern:     pattern,
		TargetUID:   target,
		Key:         key,
	}
}

func mustDigest(t *testing.T, body string, pattern, ifMatch string) cloudmodel.Digest {
	t.Helper()
	d, err := cloudmodel.CanonicalDigest(cloudmodel.CanonicalDigestInput{
		Body:    json.RawMessage(body),
		Pattern: pattern,
		IfMatch: ifMatch,
	})
	if err != nil {
		t.Fatalf("CanonicalDigest: %v", err)
	}
	return d
}

func TestIdempotencyApplies_PATCHExcluded(t *testing.T) {
	t.Parallel()
	if cloudmodel.IdempotencyApplies(http.MethodPatch) {
		t.Fatal("PATCH must not participate in idempotency state")
	}
	if cloudmodel.IdempotencyApplies(http.MethodGet) {
		t.Fatal("GET must not participate in idempotency state")
	}
	if !cloudmodel.IdempotencyApplies(http.MethodPost) {
		t.Fatal("POST must participate in idempotency state")
	}
}

func TestCanonicalDigest_ExcludesKeyAndIncludesIfMatch(t *testing.T) {
	t.Parallel()
	pattern := "POST /v1/cloud-provider-participations/{uid}/accept"
	d1 := mustDigest(t, `{}`, pattern, `"1"`)
	d2 := mustDigest(t, `{}`, pattern, `"2"`)
	if d1.Equal(d2) {
		t.Fatal("changed If-Match must change digest")
	}
	d3 := mustDigest(t, `{"metadata":{"name":"a"}}`, "POST /v1/cloud-platforms", "")
	d4 := mustDigest(t, `{"metadata":{"name":"b"}}`, "POST /v1/cloud-platforms", "")
	if d3.Equal(d4) {
		t.Fatal("changed body must change digest")
	}
	// Same classified body/pattern/If-Match are stable regardless of map iteration.
	d5 := mustDigest(t, `{"a":1,"b":2}`, pattern, `"1"`)
	d6 := mustDigest(t, `{"b":2,"a":1}`, pattern, `"1"`)
	if !d5.Equal(d6) {
		t.Fatal("canonical digest must sort object keys")
	}
}

func TestIdempotencyNamespace_FingerprintRedactsKey(t *testing.T) {
	t.Parallel()
	ns := testNS("principal-1", "POST /v1/cloud-platforms", "", "secret-key-value")
	fp := ns.Fingerprint()
	if fp == "" || fp == ns.Key || len(fp) != 64 {
		t.Fatalf("fingerprint = %q", fp)
	}
	if fp == ns.PrincipalID+ns.Pattern+ns.Key {
		t.Fatal("fingerprint must not equal concatenated namespace")
	}
}

func TestIdempotency_ReserveCompleteReplay(t *testing.T) {
	t.Parallel()
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
	ns := testNS("p1", "POST /v1/cloud-platforms", "", "k1")
	digest := mustDigest(t, `{"metadata":{"name":"plat"}}`, ns.Pattern, "")

	res := c.Reserve(context.Background(), ns, digest, nil)
	if res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("first reserve: %#v", res)
	}
	body := []byte(`{"metadata":{"name":"plat","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}}`)
	c.Complete(ns, digest, cloudmodel.CompletedResult{StatusCode: http.StatusCreated, Body: body})

	recheckCalled := false
	res = c.Reserve(context.Background(), ns, digest, func(ctx context.Context) *apiproblem.Problem {
		recheckCalled = true
		return nil
	})
	if res.Outcome != cloudmodel.ReserveOutcomeReplay || res.Replay == nil {
		t.Fatalf("replay: %#v", res)
	}
	if !recheckCalled {
		t.Fatal("completed replay requires current auth/authz/safe-access recheck")
	}
	if res.Replay.StatusCode != http.StatusCreated || string(res.Replay.Body) != string(body) {
		t.Fatalf("replay payload: %#v", res.Replay)
	}
	if res.Replay.ContentType != "application/json" {
		t.Fatalf("content-type = %q", res.Replay.ContentType)
	}
}

func TestIdempotency_SameKeyDifferentDigestConflict(t *testing.T) {
	t.Parallel()
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
	ns := testNS("p1", "POST /v1/cloud-platforms", "", "k1")
	d1 := mustDigest(t, `{"metadata":{"name":"a"}}`, ns.Pattern, "")
	d2 := mustDigest(t, `{"metadata":{"name":"b"}}`, ns.Pattern, "")

	if res := c.Reserve(context.Background(), ns, d1, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	c.Complete(ns, d1, cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{"ok":true}`)})

	res := c.Reserve(context.Background(), ns, d2, nil)
	if res.Outcome != cloudmodel.ReserveOutcomeConflict || res.Problem == nil {
		t.Fatalf("conflict: %#v", res)
	}
	if res.Problem.Code != apiproblem.CodeConflict {
		t.Fatalf("code = %s", res.Problem.Code)
	}
	if len(res.Problem.Violations) == 0 || res.Problem.Violations[0].Code != cloudmodel.ViolationIdempotencyKeyReuseMismatch {
		t.Fatalf("violations: %#v", res.Problem.Violations)
	}
	if res.Replay != nil {
		t.Fatal("conflict must not disclose stored result")
	}
}

func TestIdempotency_ChangedIfMatchConflictOnAction(t *testing.T) {
	t.Parallel()
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
	ns := testNS("p1", "POST /v1/cloud-provider-participations/{uid}/accept", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "k-action")
	d1 := mustDigest(t, `{}`, ns.Pattern, `"1"`)
	d2 := mustDigest(t, `{}`, ns.Pattern, `"2"`)
	if res := c.Reserve(context.Background(), ns, d1, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	res := c.Reserve(context.Background(), ns, d2, nil)
	if res.Outcome != cloudmodel.ReserveOutcomeConflict {
		t.Fatalf("changed If-Match while InFlight: %#v", res)
	}
}

func TestIdempotency_TargetUIDIsolatesNamespaces(t *testing.T) {
	t.Parallel()
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
	pattern := "POST /v1/cloud-provider-participations/{uid}/accept"
	nsA := testNS("p1", pattern, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "same-key")
	nsB := testNS("p1", pattern, "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "same-key")
	d := mustDigest(t, `{}`, pattern, `"1"`)
	if res := c.Reserve(context.Background(), nsA, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("nsA: %#v", res)
	}
	if res := c.Reserve(context.Background(), nsB, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("nsB must reserve independently: %#v", res)
	}
}

func TestIdempotency_StatesInFlightCompletedAborted(t *testing.T) {
	t.Parallel()
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
	ns := testNS("p1", "POST /v1/cloud-platforms", "", "k-state")
	d := mustDigest(t, `{"metadata":{"name":"x"}}`, ns.Pattern, "")

	if res := c.Reserve(context.Background(), ns, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	if st, ok := c.StateForTest(ns); !ok || st != cloudmodel.ReservationInFlight {
		t.Fatalf("state = %s ok=%v", st, ok)
	}
	c.Complete(ns, d, cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{}`)})
	if st, ok := c.StateForTest(ns); !ok || st != cloudmodel.ReservationCompleted {
		t.Fatalf("completed state = %s ok=%v", st, ok)
	}

	ns2 := testNS("p1", "POST /v1/cloud-platforms", "", "k-abort")
	if res := c.Reserve(context.Background(), ns2, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("reserve2: %#v", res)
	}
	c.Abort(ns2)
	if st, ok := c.StateForTest(ns2); !ok || st != cloudmodel.ReservationAborted {
		t.Fatalf("aborted state = %s ok=%v", st, ok)
	}
	// Aborted record is removed before re-reservation.
	if res := c.Reserve(context.Background(), ns2, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("re-reserve after abort: %#v", res)
	}
}

func TestIdempotency_WaiterDetachesOnCancellation(t *testing.T) {
	t.Parallel()
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
	ns := testNS("p1", "POST /v1/cloud-platforms", "", "k-cancel")
	d := mustDigest(t, `{"metadata":{"name":"x"}}`, ns.Pattern, "")
	if res := c.Reserve(context.Background(), ns, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("owner: %#v", res)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan cloudmodel.ReserveResult, 1)
	go func() {
		done <- c.Reserve(ctx, ns, d, nil)
	}()
	time.Sleep(20 * time.Millisecond)
	cancel()
	res := <-done
	if res.Outcome != cloudmodel.ReserveOutcomeDenied {
		t.Fatalf("cancelled waiter: %#v", res)
	}
	if st, ok := c.StateForTest(ns); !ok || st != cloudmodel.ReservationInFlight {
		t.Fatalf("cancellation must not abort owner reservation: %s ok=%v", st, ok)
	}
	c.Complete(ns, d, cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{"ok":1}`)})
}

func TestIdempotency_AbortWakesWaiters(t *testing.T) {
	t.Parallel()
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
	ns := testNS("p1", "POST /v1/cloud-platforms", "", "k-wake")
	d := mustDigest(t, `{"metadata":{"name":"x"}}`, ns.Pattern, "")
	if res := c.Reserve(context.Background(), ns, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("owner: %#v", res)
	}

	started := make(chan struct{})
	done := make(chan cloudmodel.ReserveResult, 1)
	go func() {
		close(started)
		done <- c.Reserve(context.Background(), ns, d, nil)
	}()
	<-started
	time.Sleep(20 * time.Millisecond)
	c.Abort(ns)
	res := <-done
	// After abort, waiter removes Aborted and re-reserves.
	if res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("waiter after abort: %#v", res)
	}
}

func TestIdempotency_AbortReasons(t *testing.T) {
	t.Parallel()
	reasons := []string{
		"validation", "uniqueness", "lifecycle-source-state", "stale-version",
		"audit-failure", "recovered-panic", "graceful-shutdown",
	}
	for _, reason := range reasons {
		reason := reason
		t.Run(reason, func(t *testing.T) {
			t.Parallel()
			c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
			ns := testNS("p1", "POST /v1/cloud-platforms", "", "k-"+reason)
			d := mustDigest(t, `{"metadata":{"name":"x"}}`, ns.Pattern, "")
			if res := c.Reserve(context.Background(), ns, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
				t.Fatalf("reserve: %#v", res)
			}
			if reason == "graceful-shutdown" {
				c.AbortAll()
			} else {
				c.Abort(ns)
			}
			if st, ok := c.StateForTest(ns); !ok || st != cloudmodel.ReservationAborted {
				t.Fatalf("%s: state=%s ok=%v", reason, st, ok)
			}
			if c.CompletedCountForTest() != 0 {
				t.Fatalf("%s left a completed record", reason)
			}
		})
	}
}

func TestIdempotency_ReplayRequiresRecheck(t *testing.T) {
	t.Parallel()
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
	ns := testNS("p1", "POST /v1/cloud-platforms", "", "k-recheck")
	d := mustDigest(t, `{"metadata":{"name":"x"}}`, ns.Pattern, "")
	if res := c.Reserve(context.Background(), ns, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	secretBody := []byte(`{"secret":"must-not-disclose"}`)
	c.Complete(ns, d, cloudmodel.CompletedResult{StatusCode: 201, Body: secretBody})

	denied := apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("grant revoked")
	res := c.Reserve(context.Background(), ns, d, func(ctx context.Context) *apiproblem.Problem {
		return denied
	})
	if res.Outcome != cloudmodel.ReserveOutcomeDenied || res.Problem != denied {
		t.Fatalf("denied replay: %#v", res)
	}
	if res.Replay != nil {
		t.Fatal("revoked grant must not disclose stored result")
	}
}

func TestIdempotency_RetentionEviction(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 8, 13, 12, 0, 0, 0, time.UTC)}
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{Clock: clock})

	// Insert MaxCompletedRecords+2 completed records with staggered expiry.
	const extra = 2
	total := cloudmodel.MaxCompletedRecords + extra
	for i := 0; i < total; i++ {
		ns := testNS("p1", "POST /v1/cloud-platforms", "", "ret-"+padIndex(i))
		d := mustDigest(t, `{"metadata":{"name":"n"}}`, ns.Pattern, "")
		if res := c.Reserve(context.Background(), ns, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
			t.Fatalf("reserve %d: %#v", i, res)
		}
		// Advance clock slightly so completedAt/expiresAt order is deterministic.
		clock.Set(clock.Now().Add(time.Millisecond))
		c.Complete(ns, d, cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{}`)})
	}
	if got := c.CompletedCountForTest(); got != cloudmodel.MaxCompletedRecords {
		t.Fatalf("completed count = %d, want %d", got, cloudmodel.MaxCompletedRecords)
	}

	// Earliest-expiry (first inserted) must be gone; latest remain.
	first := testNS("p1", "POST /v1/cloud-platforms", "", "ret-"+padIndex(0))
	if _, ok := c.StateForTest(first); ok {
		t.Fatal("earliest completed record must be evicted")
	}
	last := testNS("p1", "POST /v1/cloud-platforms", "", "ret-"+padIndex(total-1))
	if st, ok := c.StateForTest(last); !ok || st != cloudmodel.ReservationCompleted {
		t.Fatalf("latest completed must remain: %s ok=%v", st, ok)
	}

	// InFlight is never evicted by the completed-record cap.
	inFlightNS := testNS("p1", "POST /v1/cloud-platforms", "", "inflight-keep")
	d := mustDigest(t, `{"metadata":{"name":"z"}}`, inFlightNS.Pattern, "")
	if res := c.Reserve(context.Background(), inFlightNS, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("inflight reserve: %#v", res)
	}
	// Force another completed insertion cycle.
	for i := 0; i < 3; i++ {
		ns := testNS("p1", "POST /v1/cloud-platforms", "", "more-"+padIndex(i))
		dig := mustDigest(t, `{"metadata":{"name":"m"}}`, ns.Pattern, "")
		_ = c.Reserve(context.Background(), ns, dig, nil)
		clock.Set(clock.Now().Add(time.Millisecond))
		c.Complete(ns, dig, cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{}`)})
	}
	if st, ok := c.StateForTest(inFlightNS); !ok || st != cloudmodel.ReservationInFlight {
		t.Fatalf("InFlight must never be evicted: %s ok=%v", st, ok)
	}
	if c.InFlightCountForTest() < 1 {
		t.Fatal("expected InFlight record")
	}
}

func TestIdempotency_TwentyFourHourExpiry(t *testing.T) {
	t.Parallel()
	clock := &fakeClock{now: time.Date(2026, 8, 13, 0, 0, 0, 0, time.UTC)}
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{Clock: clock})
	ns := testNS("p1", "POST /v1/cloud-platforms", "", "k-24h")
	d := mustDigest(t, `{"metadata":{"name":"x"}}`, ns.Pattern, "")
	_ = c.Reserve(context.Background(), ns, d, nil)
	c.Complete(ns, d, cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{}`)})

	clock.Set(clock.Now().Add(cloudmodel.CompletedRetention - time.Second))
	if st, ok := c.StateForTest(ns); !ok || st != cloudmodel.ReservationCompleted {
		t.Fatalf("before expiry: %s ok=%v", st, ok)
	}
	clock.Set(clock.Now().Add(2 * time.Second))
	if _, ok := c.StateForTest(ns); ok {
		t.Fatal("completed record must expire after 24 hours")
	}
	// Expired namespace is processed as new.
	if res := c.Reserve(context.Background(), ns, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("after expiry: %#v", res)
	}
}

func TestIdempotency_LexicalEvictionOnExpiryTie(t *testing.T) {
	clock := &fixedClock{now: time.Date(2026, 8, 13, 15, 0, 0, 0, time.UTC)}
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{Clock: clock})

	// Fill to capacity with identical expiry, then insert one more.
	for i := 0; i < cloudmodel.MaxCompletedRecords; i++ {
		ns := testNS("p1", "POST /v1/cloud-platforms", "", "tie-"+padIndex(i))
		d := mustDigest(t, `{"metadata":{"name":"n"}}`, ns.Pattern, "")
		_ = c.Reserve(context.Background(), ns, d, nil)
		c.Complete(ns, d, cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{}`)})
	}
	// Insert lexical-earliest key with same expiry timestamp — on next insert,
	// earliest-expiry then lexical eviction removes the smallest namespace key.
	overflow := testNS("p1", "POST /v1/cloud-platforms", "", "tie-overflow")
	d := mustDigest(t, `{"metadata":{"name":"n"}}`, overflow.Pattern, "")
	_ = c.Reserve(context.Background(), overflow, d, nil)
	c.Complete(overflow, d, cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{}`)})

	first := testNS("p1", "POST /v1/cloud-platforms", "", "tie-"+padIndex(0))
	if _, ok := c.StateForTest(first); ok {
		t.Fatal("lexical-earliest completed record must be evicted on expiry tie")
	}
	if got := c.CompletedCountForTest(); got != cloudmodel.MaxCompletedRecords {
		t.Fatalf("count = %d", got)
	}
}

func TestIdempotency_ConcurrentSameDigestWaiterReplay(t *testing.T) {
	c := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
	ns := testNS("p1", "POST /v1/cloud-platforms", "", "k-race")
	d := mustDigest(t, `{"metadata":{"name":"x"}}`, ns.Pattern, "")

	var wg sync.WaitGroup
	const waiters = 8
	results := make([]cloudmodel.ReserveResult, waiters)
	ownerReady := make(chan struct{})
	startWaiters := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		res := c.Reserve(context.Background(), ns, d, nil)
		if res.Outcome != cloudmodel.ReserveOutcomeReserved {
			t.Errorf("owner: %#v", res)
			close(ownerReady)
			return
		}
		close(ownerReady)
		<-startWaiters
		time.Sleep(10 * time.Millisecond)
		c.Complete(ns, d, cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{"uid":"a"}`)})
	}()

	<-ownerReady
	wg.Add(waiters)
	for i := 0; i < waiters; i++ {
		i := i
		go func() {
			defer wg.Done()
			<-startWaiters
			results[i] = c.Reserve(context.Background(), ns, d, func(ctx context.Context) *apiproblem.Problem {
				return nil
			})
		}()
	}
	close(startWaiters)
	wg.Wait()

	for i, res := range results {
		if res.Outcome != cloudmodel.ReserveOutcomeReplay || res.Replay == nil {
			t.Fatalf("waiter %d: %#v", i, res)
		}
		if string(res.Replay.Body) != `{"uid":"a"}` {
			t.Fatalf("waiter %d body = %s", i, res.Replay.Body)
		}
	}
}

func padIndex(i int) string {
	return string([]byte{
		'0' + byte((i/10000)%10),
		'0' + byte((i/1000)%10),
		'0' + byte((i/100)%10),
		'0' + byte((i/10)%10),
		'0' + byte(i%10),
	})
}

// fixedClock never moves unless Set is called; used for lexical-tie tests.
type fixedClock struct {
	now time.Time
}

func (c *fixedClock) Now() time.Time { return c.now.UTC() }
