package executiontarget

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"
)

func testNS(principal, pattern, cloudProvider, target, key string) IdempotencyNamespace {
	return IdempotencyNamespace{
		PrincipalUID:     principal,
		RoutePattern:     pattern,
		CloudProviderUID: cloudProvider,
		ActionTargetUID:  target,
		Key:              key,
	}
}

func mustCreateDigest(t *testing.T, body string) Digest {
	t.Helper()
	d, err := DigestCreate(json.RawMessage(body))
	if err != nil {
		t.Fatalf("DigestCreate: %v", err)
	}
	return d
}

func padIndex(i int) string {
	return fmt.Sprintf("%05d", i)
}

func TestIdempotencyApplies_GETAndLISTExcluded(t *testing.T) {
	t.Parallel()
	if IdempotencyApplies(http.MethodGet) {
		t.Fatal("GET must not participate in idempotency state")
	}
	if IdempotencyApplies(http.MethodHead) {
		t.Fatal("HEAD must not participate in idempotency state")
	}
	if !IdempotencyApplies(http.MethodPost) {
		t.Fatal("POST must participate in idempotency state")
	}
}

func TestIdempotencyTable_HasNoIndependentMutex(t *testing.T) {
	t.Parallel()
	var zero IdempotencyTable
	rt := reflect.TypeOf(zero)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		switch f.Type.String() {
		case "sync.Mutex", "sync.RWMutex":
			t.Fatalf("IdempotencyTable must not own an independent mutex; found %s", f.Type)
		}
	}
}

func TestDigestCreate_RecursivelySortedExcludesHeaders(t *testing.T) {
	t.Parallel()
	d1 := mustCreateDigest(t, `{"metadata":{"name":"a"},"spec":{"targetClass":"synthetic-iaas"}}`)
	d2 := mustCreateDigest(t, `{"spec":{"targetClass":"synthetic-iaas"},"metadata":{"name":"a"}}`)
	if !d1.Equal(d2) {
		t.Fatal("create digest must recursively sort object keys")
	}
	d3 := mustCreateDigest(t, `{"metadata":{"name":"b"},"spec":{"targetClass":"synthetic-iaas"}}`)
	if d1.Equal(d3) {
		t.Fatal("changed create JSON must change digest")
	}
	// DigestCreate has no If-Match / Accept / auth / correlation parameters by
	// construction; DigestAction is constant for all actions.
	a1 := DigestAction()
	a2 := DigestAction()
	if !a1.Equal(a2) {
		t.Fatal("action digest must be SHA-256 of zero bytes and stable")
	}
	want := Digest(sha256.Sum256(nil))
	if !a1.Equal(want) {
		t.Fatalf("action digest = %s, want %s", a1, want)
	}
	if d1.Equal(a1) {
		t.Fatal("create digest must not equal action zero-byte digest")
	}
}

func TestIdempotencyNamespace_IncludesCloudProviderAndTarget(t *testing.T) {
	t.Parallel()
	nsA := testNS("p1", "POST /v1/execution-targets/{uid}/qualify", "cp-a", "tgt-a", "same-key")
	nsB := testNS("p1", "POST /v1/execution-targets/{uid}/qualify", "cp-a", "tgt-b", "same-key")
	nsC := testNS("p1", "POST /v1/execution-targets", "cp-b", "", "same-key")
	if nsA.key() == nsB.key() {
		t.Fatal("different action-target UIDs must form independent namespaces")
	}
	if nsA.key() == nsC.key() {
		t.Fatal("different derived CloudProvider scopes must form independent namespaces")
	}
}

func TestIdempotency_ReserveCompleteReplay(t *testing.T) {
	t.Parallel()
	table := NewIdempotencyTable()
	now := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	ns := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "k1")
	digest := mustCreateDigest(t, `{"metadata":{"name":"tgt"}}`)

	res := table.reserveOrInspect(ns, digest, now)
	if res.Outcome != reserveOutcomeReserved {
		t.Fatalf("first reserve: %#v", res)
	}
	if st, ok := table.stateForTest(ns, now); !ok || st != ReservationInFlight {
		t.Fatalf("state = %s ok=%v", st, ok)
	}
	body := []byte(`{"metadata":{"uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","name":"tgt"}}`)
	done, ok := table.finalizeComplete(ns, digest, CompletedResult{StatusCode: http.StatusCreated, Body: body}, now)
	if !ok || done == nil {
		t.Fatal("finalizeComplete must capture done channel")
	}
	closeReservationDone(done)

	res = table.reserveOrInspect(ns, digest, now)
	if res.Outcome != reserveOutcomeReplay || res.Replay == nil {
		t.Fatalf("replay: %#v", res)
	}
	if res.Replay.StatusCode != http.StatusCreated || string(res.Replay.Body) != string(body) {
		t.Fatalf("replay payload: %#v", res.Replay)
	}
	if st, ok := table.stateForTest(ns, now); !ok || st != ReservationCompleted {
		t.Fatalf("completed state = %s ok=%v", st, ok)
	}
}

func TestIdempotency_SameKeyDifferentDigestConflict(t *testing.T) {
	t.Parallel()
	table := NewIdempotencyTable()
	now := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	ns := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "k1")
	d1 := mustCreateDigest(t, `{"metadata":{"name":"a"}}`)
	d2 := mustCreateDigest(t, `{"metadata":{"name":"b"}}`)

	if res := table.reserveOrInspect(ns, d1, now); res.Outcome != reserveOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	done, ok := table.finalizeComplete(ns, d1, CompletedResult{StatusCode: 201, Body: []byte(`{"ok":true}`)}, now)
	if !ok {
		t.Fatal("complete")
	}
	closeReservationDone(done)

	res := table.reserveOrInspect(ns, d2, now)
	if res.Outcome != reserveOutcomeConflict {
		t.Fatalf("conflict: %#v", res)
	}
	if st, ok := table.stateForTest(ns, now); !ok || st != ReservationCompleted {
		t.Fatalf("original completion must remain: %s ok=%v", st, ok)
	}
}

func TestIdempotency_InFlightDigestConflict(t *testing.T) {
	t.Parallel()
	table := NewIdempotencyTable()
	now := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	ns := testNS("p1", "POST /v1/execution-targets/{uid}/retire", "cp-1", "tgt-1", "k-action")
	d1 := DigestAction()
	// Simulate a different digest while InFlight (should not occur for actions
	// in production, but conflict must still be reported).
	var d2 Digest
	copy(d2[:], d1[:])
	d2[0] ^= 0xff

	if res := table.reserveOrInspect(ns, d1, now); res.Outcome != reserveOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	res := table.reserveOrInspect(ns, d2, now)
	if res.Outcome != reserveOutcomeConflict {
		t.Fatalf("InFlight conflict: %#v", res)
	}
	if st, ok := table.stateForTest(ns, now); !ok || st != ReservationInFlight {
		t.Fatalf("owner must remain InFlight: %s ok=%v", st, ok)
	}
}

func TestIdempotency_IndependentAcrossActionTargetsAndScopes(t *testing.T) {
	t.Parallel()
	table := NewIdempotencyTable()
	now := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	pattern := "POST /v1/execution-targets/{uid}/qualify"
	digest := DigestAction()

	nsA := testNS("p1", pattern, "cp-1", "tgt-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "same-key")
	nsB := testNS("p1", pattern, "cp-1", "tgt-bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", "same-key")
	if res := table.reserveOrInspect(nsA, digest, now); res.Outcome != reserveOutcomeReserved {
		t.Fatalf("nsA: %#v", res)
	}
	if res := table.reserveOrInspect(nsB, digest, now); res.Outcome != reserveOutcomeReserved {
		t.Fatalf("nsB must reserve independently: %#v", res)
	}

	createA := testNS("p1", "POST /v1/execution-targets", "cp-aaaa", "", "same-key")
	createB := testNS("p1", "POST /v1/execution-targets", "cp-bbbb", "", "same-key")
	cd := mustCreateDigest(t, `{"metadata":{"name":"x"}}`)
	if res := table.reserveOrInspect(createA, cd, now); res.Outcome != reserveOutcomeReserved {
		t.Fatalf("createA: %#v", res)
	}
	if res := table.reserveOrInspect(createB, cd, now); res.Outcome != reserveOutcomeReserved {
		t.Fatalf("createB must reserve independently across CloudProvider scope: %#v", res)
	}
}

func TestIdempotency_WaiterSelectsDoneOrContext_CancelDetachesOnly(t *testing.T) {
	t.Parallel()
	table := NewIdempotencyTable()
	var mu sync.Mutex
	now := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	ns := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "k-cancel")
	digest := mustCreateDigest(t, `{"metadata":{"name":"x"}}`)

	mu.Lock()
	if res := table.reserveOrInspect(ns, digest, now); res.Outcome != reserveOutcomeReserved {
		mu.Unlock()
		t.Fatalf("owner: %#v", res)
	}
	mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	waiterErr := make(chan error, 1)
	go func() {
		mu.Lock()
		if res := table.reserveOrInspect(ns, digest, now); res.Outcome != reserveOutcomeInFlight {
			mu.Unlock()
			waiterErr <- fmt.Errorf("expected InFlight, got %#v", res)
			return
		}
		done, ok := table.attachWaiter(ns)
		mu.Unlock()
		if !ok {
			waiterErr <- fmt.Errorf("attachWaiter failed")
			return
		}
		err := awaitReservationDone(ctx, done)
		if err != nil {
			mu.Lock()
			table.detachWaiter(ns)
			mu.Unlock()
		}
		waiterErr <- err
	}()

	time.Sleep(20 * time.Millisecond)
	cancel()
	err := <-waiterErr
	if err == nil {
		t.Fatal("cancelled waiter must return context error")
	}
	mu.Lock()
	st, ok := table.stateForTest(ns, now)
	mu.Unlock()
	if !ok || st != ReservationInFlight {
		t.Fatalf("cancellation must only detach waiter; owner remains InFlight: %s ok=%v", st, ok)
	}

	// Owner can still complete after waiter detach.
	mu.Lock()
	done, ok := table.finalizeComplete(ns, digest, CompletedResult{StatusCode: 201, Body: []byte(`{}`)}, now)
	mu.Unlock()
	if !ok {
		t.Fatal("owner complete after detach")
	}
	closeReservationDone(done)
}

func TestIdempotency_AbortWakesWaitersWithoutCompletion(t *testing.T) {
	t.Parallel()
	table := NewIdempotencyTable()
	var mu sync.Mutex
	now := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	ns := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "k-abort")
	digest := mustCreateDigest(t, `{"metadata":{"name":"x"}}`)

	mu.Lock()
	if res := table.reserveOrInspect(ns, digest, now); res.Outcome != reserveOutcomeReserved {
		mu.Unlock()
		t.Fatalf("owner: %#v", res)
	}
	mu.Unlock()

	started := make(chan struct{})
	woken := make(chan error, 1)
	go func() {
		mu.Lock()
		if res := table.reserveOrInspect(ns, digest, now); res.Outcome != reserveOutcomeInFlight {
			mu.Unlock()
			woken <- fmt.Errorf("expected InFlight: %#v", res)
			return
		}
		done, ok := table.attachWaiter(ns)
		mu.Unlock()
		close(started)
		if !ok {
			woken <- fmt.Errorf("attachWaiter failed")
			return
		}
		woken <- awaitReservationDone(context.Background(), done)
	}()
	<-started
	time.Sleep(20 * time.Millisecond)

	mu.Lock()
	done, ok := table.finalizeAbort(ns)
	mu.Unlock()
	if !ok {
		t.Fatal("finalizeAbort")
	}
	closeReservationDone(done)

	if err := <-woken; err != nil {
		t.Fatalf("waiter wake: %v", err)
	}
	mu.Lock()
	_, present := table.stateForTest(ns, now)
	completed := table.completedCountForTest(now)
	mu.Unlock()
	if present {
		t.Fatal("abort must remove reservation; no Aborted retention")
	}
	if completed != 0 {
		t.Fatal("abort must not create a completion")
	}

	// After abort, a new reservation can be established.
	mu.Lock()
	res := table.reserveOrInspect(ns, digest, now)
	mu.Unlock()
	if res.Outcome != reserveOutcomeReserved {
		t.Fatalf("re-reserve after abort: %#v", res)
	}
}

func TestIdempotency_CompleteClosesDoneExactlyOnceAfterUnlock(t *testing.T) {
	t.Parallel()
	table := NewIdempotencyTable()
	var mu sync.Mutex
	now := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	ns := testNS("p1", "POST /v1/execution-targets/{uid}/retire", "cp-1", "tgt-1", "k-done")
	digest := DigestAction()

	mu.Lock()
	_ = table.reserveOrInspect(ns, digest, now)
	waiterDone, ok := table.attachWaiter(ns)
	if !ok {
		mu.Unlock()
		t.Fatal("attach")
	}
	// Capture under mutex, then unlock before close (DD-05 two-phase protocol).
	done, ok := table.finalizeComplete(ns, digest, CompletedResult{StatusCode: 200, Body: []byte(`{}`)}, now)
	mu.Unlock()
	if !ok {
		t.Fatal("finalizeComplete")
	}

	select {
	case <-waiterDone:
		t.Fatal("done must not be closed while mutex protocol is still holding the channel pre-close")
	default:
	}
	closeReservationDone(done)
	select {
	case <-waiterDone:
	case <-time.After(time.Second):
		t.Fatal("done must wake waiters after post-unlock close")
	}
}

func TestIdempotency_TwentyFourHourExpiry(t *testing.T) {
	t.Parallel()
	table := NewIdempotencyTable()
	now := time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC)
	ns := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "k-24h")
	digest := mustCreateDigest(t, `{"metadata":{"name":"x"}}`)

	_ = table.reserveOrInspect(ns, digest, now)
	done, ok := table.finalizeComplete(ns, digest, CompletedResult{StatusCode: 201, Body: []byte(`{}`)}, now)
	if !ok {
		t.Fatal("complete")
	}
	closeReservationDone(done)

	before := now.Add(CompletedRetention - time.Second)
	if st, ok := table.stateForTest(ns, before); !ok || st != ReservationCompleted {
		t.Fatalf("before expiry: %s ok=%v", st, ok)
	}
	after := now.Add(CompletedRetention)
	if _, ok := table.stateForTest(ns, after); ok {
		t.Fatal("completed record must expire after 24 hours")
	}
	if res := table.reserveOrInspect(ns, digest, after); res.Outcome != reserveOutcomeReserved {
		t.Fatalf("after expiry processed as new: %#v", res)
	}
}

func TestIdempotency_RetentionEviction_InFlightNeverEvicts(t *testing.T) {
	table := NewIdempotencyTable()
	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)

	const extra = 2
	total := MaxCompletedRecords + extra
	for i := 0; i < total; i++ {
		ns := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "ret-"+padIndex(i))
		d := mustCreateDigest(t, `{"metadata":{"name":"n"}}`)
		if res := table.reserveOrInspect(ns, d, now); res.Outcome != reserveOutcomeReserved {
			t.Fatalf("reserve %d: %#v", i, res)
		}
		now = now.Add(time.Millisecond)
		done, ok := table.finalizeComplete(ns, d, CompletedResult{StatusCode: 201, Body: []byte(`{}`)}, now)
		if !ok {
			t.Fatalf("complete %d", i)
		}
		closeReservationDone(done)
	}
	if got := table.completedCountForTest(now); got != MaxCompletedRecords {
		t.Fatalf("completed count = %d, want %d", got, MaxCompletedRecords)
	}

	first := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "ret-"+padIndex(0))
	if _, ok := table.stateForTest(first, now); ok {
		t.Fatal("earliest completed record must be evicted")
	}
	last := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "ret-"+padIndex(total-1))
	if st, ok := table.stateForTest(last, now); !ok || st != ReservationCompleted {
		t.Fatalf("latest completed must remain: %s ok=%v", st, ok)
	}

	inFlightNS := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "inflight-keep")
	d := mustCreateDigest(t, `{"metadata":{"name":"z"}}`)
	if res := table.reserveOrInspect(inFlightNS, d, now); res.Outcome != reserveOutcomeReserved {
		t.Fatalf("inflight reserve: %#v", res)
	}
	for i := 0; i < 3; i++ {
		ns := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "more-"+padIndex(i))
		dig := mustCreateDigest(t, `{"metadata":{"name":"m"}}`)
		_ = table.reserveOrInspect(ns, dig, now)
		now = now.Add(time.Millisecond)
		done, ok := table.finalizeComplete(ns, dig, CompletedResult{StatusCode: 201, Body: []byte(`{}`)}, now)
		if !ok {
			t.Fatalf("complete more %d", i)
		}
		closeReservationDone(done)
	}
	if st, ok := table.stateForTest(inFlightNS, now); !ok || st != ReservationInFlight {
		t.Fatalf("InFlight must never be evicted: %s ok=%v", st, ok)
	}
	if table.inFlightCountForTest() < 1 {
		t.Fatal("expected InFlight record")
	}
}

func TestIdempotency_LexicalEvictionOnExpiryTie(t *testing.T) {
	table := NewIdempotencyTable()
	now := time.Date(2026, 8, 21, 15, 0, 0, 0, time.UTC)

	for i := 0; i < MaxCompletedRecords; i++ {
		ns := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "tie-"+padIndex(i))
		d := mustCreateDigest(t, `{"metadata":{"name":"n"}}`)
		_ = table.reserveOrInspect(ns, d, now)
		done, ok := table.finalizeComplete(ns, d, CompletedResult{StatusCode: 201, Body: []byte(`{}`)}, now)
		if !ok {
			t.Fatalf("complete %d", i)
		}
		closeReservationDone(done)
	}

	overflow := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "tie-overflow")
	d := mustCreateDigest(t, `{"metadata":{"name":"n"}}`)
	_ = table.reserveOrInspect(overflow, d, now)
	done, ok := table.finalizeComplete(overflow, d, CompletedResult{StatusCode: 201, Body: []byte(`{}`)}, now)
	if !ok {
		t.Fatal("overflow complete")
	}
	closeReservationDone(done)

	// At identical expiry, earliest lexical namespace key is evicted first.
	earliest := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "tie-"+padIndex(0))
	if _, ok := table.stateForTest(earliest, now); ok {
		t.Fatal("lexical-earliest completed record must be evicted on capacity overflow")
	}
	if got := table.completedCountForTest(now); got != MaxCompletedRecords {
		t.Fatalf("completed count = %d, want %d", got, MaxCompletedRecords)
	}
}

func TestIdempotency_FinalizeAbortAll(t *testing.T) {
	t.Parallel()
	table := NewIdempotencyTable()
	now := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	ns1 := testNS("p1", "POST /v1/execution-targets", "cp-1", "", "a1")
	ns2 := testNS("p1", "POST /v1/execution-targets/{uid}/qualify", "cp-1", "tgt-1", "a2")
	d1 := mustCreateDigest(t, `{"metadata":{"name":"x"}}`)
	d2 := DigestAction()

	_ = table.reserveOrInspect(ns1, d1, now)
	_ = table.reserveOrInspect(ns2, d2, now)
	doneComplete, ok := table.finalizeComplete(ns1, d1, CompletedResult{StatusCode: 201, Body: []byte(`{}`)}, now)
	if !ok {
		t.Fatal("complete ns1")
	}
	closeReservationDone(doneComplete)

	channels := table.finalizeAbortAll()
	if len(channels) != 1 {
		t.Fatalf("abortAll should capture one InFlight done channel, got %d", len(channels))
	}
	for _, ch := range channels {
		closeReservationDone(ch)
	}
	if table.inFlightCountForTest() != 0 {
		t.Fatal("no InFlight may remain after abortAll")
	}
	if st, ok := table.stateForTest(ns1, now); !ok || st != ReservationCompleted {
		t.Fatalf("completed entries must survive abortAll: %s ok=%v", st, ok)
	}
}

func TestIdempotency_ProtocolUnderExternalMutex_Race(t *testing.T) {
	table := NewIdempotencyTable()
	var mu sync.Mutex
	now := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
	ns := testNS("p1", "POST /v1/execution-targets", "cp-race", "", "k-race")
	digest := mustCreateDigest(t, `{"metadata":{"name":"race"}}`)

	mu.Lock()
	if res := table.reserveOrInspect(ns, digest, now); res.Outcome != reserveOutcomeReserved {
		mu.Unlock()
		t.Fatalf("owner: %#v", res)
	}
	mu.Unlock()

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			res := table.reserveOrInspect(ns, digest, now)
			var done <-chan struct{}
			if res.Outcome == reserveOutcomeInFlight {
				done, _ = table.attachWaiter(ns)
			}
			mu.Unlock()
			if res.Outcome == reserveOutcomeInFlight && done != nil {
				_ = awaitReservationDone(context.Background(), done)
			}
		}()
	}

	time.Sleep(30 * time.Millisecond)
	mu.Lock()
	done, ok := table.finalizeComplete(ns, digest, CompletedResult{StatusCode: 201, Body: []byte(`{}`)}, now)
	mu.Unlock()
	if !ok {
		t.Fatal("finalizeComplete")
	}
	closeReservationDone(done)
	wg.Wait()
}
