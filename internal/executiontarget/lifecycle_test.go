package executiontarget

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

type memoryAuditAppender struct {
	mu     sync.Mutex
	events []apiconform.AuditEvent
	fail   error
	calls  int
}

func (a *memoryAuditAppender) SetFail(err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.fail = err
}

func (a *memoryAuditAppender) Append(_ context.Context, event apiconform.AuditEvent) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.calls++
	if a.fail != nil {
		return a.fail
	}
	a.events = append(a.events, event)
	return nil
}

func (a *memoryAuditAppender) Len() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.events)
}

func (a *memoryAuditAppender) Calls() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.calls
}

func (a *memoryAuditAppender) Last() apiconform.AuditEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.events) == 0 {
		return apiconform.AuditEvent{}
	}
	return a.events[len(a.events)-1]
}

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t.UTC() }
}

func newLifecycleForTest(t *testing.T, audit AuditAppender) (*ExecutionTargetLifecycleService, *Store, *IdempotencyTable, *memoryAuditAppender) {
	t.Helper()
	store := NewStore()
	idemp := NewIdempotencyTable()
	var mem *memoryAuditAppender
	if audit == nil {
		mem = &memoryAuditAppender{}
		audit = mem
	} else if m, ok := audit.(*memoryAuditAppender); ok {
		mem = m
	}
	svc := NewExecutionTargetLifecycleService(LifecycleConfig{
		Store:       store,
		Idempotency: idemp,
		Audit:       audit,
		Now:         fixedNow(time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)),
	})
	return svc, store, idemp, mem
}

func createNS(key string) IdempotencyNamespace {
	return IdempotencyNamespace{
		PrincipalUID:     "principal-1",
		RoutePattern:     "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets",
		CloudProviderUID: "cccccccccccccccccccccccccccccccc",
		Key:              key,
	}
}

func retireNS(targetUID, key string) IdempotencyNamespace {
	return IdempotencyNamespace{
		PrincipalUID:     "principal-1",
		RoutePattern:     "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire",
		CloudProviderUID: "cccccccccccccccccccccccccccccccc",
		ActionTargetUID:  targetUID,
		Key:              key,
	}
}

func baseCreateReq(key string) CreateCommitRequest {
	return CreateCommitRequest{
		Name:                  "target-1",
		UID:                   "dddddddddddddddddddddddddddddddd",
		CloudProviderScopeUID: "cccccccccccccccccccccccccccccccc",
		ParticipationRef:      participationRef("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		StackRef:              stackRef("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"),
		StackGeneration:       3,
		TargetClass:           model.TargetClassSyntheticIaaS,
		Actor: apimeta.TypedRef{
			APIVersion: "identity.sovrunn.io/v1alpha1",
			Kind:       "Principal",
			Name:       "principal",
			UID:        "principal-1",
		},
		RequestID:   "req-create-1",
		AuditUID:    "audit-create-1",
		Idempotency: createNS(key),
		Digest:      DigestAction(), // tests override with create digest when needed
		Completion:  CompletedResult{StatusCode: http.StatusCreated, Body: []byte(`{"ok":true}`)},
	}
}

func TestLifecycle_SoleCommitterMutexGuardsState(t *testing.T) {
	t.Parallel()

	svc := NewExecutionTargetLifecycleService(LifecycleConfig{})
	rt := reflect.TypeOf(svc).Elem()
	foundMu := false
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.Type.String() == "sync.Mutex" {
			foundMu = true
		}
	}
	if !foundMu {
		t.Fatal("ExecutionTargetLifecycleService must own the DD-02 sync.Mutex")
	}

	// IdempotencyTable must remain mutex-free; lifecycle mutex is the sole guard.
	var zero IdempotencyTable
	irt := reflect.TypeOf(zero)
	for i := 0; i < irt.NumField(); i++ {
		f := irt.Field(i)
		switch f.Type.String() {
		case "sync.Mutex", "sync.RWMutex":
			t.Fatalf("IdempotencyTable must not own an independent mutex; found %s", f.Type)
		}
	}
}

func TestLifecycle_ClientsNeverWriteStatus(t *testing.T) {
	t.Parallel()

	svc, store, idemp, audit := newLifecycleForTest(t, nil)
	req := baseCreateReq("k-status")
	req.ClientStatusAttempted = true
	req.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-1"}}`)

	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}

	_, prob := svc.CommitCreate(context.Background(), req)
	if prob == nil || prob.Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("status write must be denied: %#v", prob)
	}
	if len(prob.Violations) == 0 || prob.Violations[0].Code != "VS0_STATUS_FIELD_WRITE" {
		t.Fatalf("want VS0_STATUS_FIELD_WRITE, got %#v", prob.Violations)
	}
	if _, ok := store.GetExecutionTarget(req.UID); ok {
		t.Fatal("status-write attempt must not publish a target")
	}
	if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
		t.Fatal("denied create must remove InFlight reservation")
	}
	if audit.Len() != 0 {
		t.Fatal("status-write denial at commit must not append create AuditEvent")
	}
}

func TestLifecycle_CommitCreate_HappyPathAuditBeforePublication(t *testing.T) {
	t.Parallel()

	svc, store, idemp, audit := newLifecycleForTest(t, nil)
	req := baseCreateReq("k-create-ok")
	req.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-1"}}`)

	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}

	got, prob := svc.CommitCreate(context.Background(), req)
	if prob != nil {
		t.Fatalf("commit: %#v", prob)
	}
	if got.Status.Lifecycle != model.LifecycleActive || got.Status.Qualification != model.QualificationUnqualified {
		t.Fatalf("state=%#v", got.Status)
	}
	if got.Status.MaintenanceEpoch != 0 || got.Status.ObservedGeneration != 3 {
		t.Fatalf("epoch/observedGeneration=%#v", got.Status)
	}
	if got.CloudProviderScopeUID != req.CloudProviderScopeUID {
		t.Fatalf("derived scope=%q", got.CloudProviderScopeUID)
	}
	if got.Metadata.ScopeRef != nil {
		t.Fatal("scopeRef must remain unset")
	}
	if audit.Calls() != 1 || audit.Len() != 1 {
		t.Fatalf("exactly one append required; calls=%d len=%d", audit.Calls(), audit.Len())
	}
	if audit.Last().Record.Action != AuditActionCreate {
		t.Fatalf("action=%q", audit.Last().Record.Action)
	}
	st, ok := idemp.stateForTest(req.Idempotency, svc.now())
	if !ok || st != ReservationCompleted {
		t.Fatalf("completion state=%s ok=%v", st, ok)
	}
	if !store.nameReserved(req.CloudProviderScopeUID, req.Name) {
		t.Fatal("name must be reserved")
	}
	if !store.liveTupleOccupied(req.ParticipationRef.UID, req.StackRef.UID, model.TargetClassSyntheticIaaS) {
		t.Fatal("live tuple must be indexed")
	}
}

func TestLifecycle_CommitCreate_AppendFailureRemovesReservationNoPublication(t *testing.T) {
	t.Parallel()

	audit := &memoryAuditAppender{fail: errors.New("append boom")}
	svc, store, idemp, _ := newLifecycleForTest(t, audit)
	req := baseCreateReq("k-create-fail")
	req.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-1"}}`)

	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}

	// Attach a waiter before the failing commit to prove post-unlock wake.
	started := make(chan struct{})
	woken := make(chan error, 1)
	go func() {
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeInFlight {
			woken <- errors.New("expected InFlight peer")
			return
		}
		done, ok := svc.AttachWaiter(req.Idempotency)
		close(started)
		if !ok {
			woken <- errors.New("attachWaiter failed")
			return
		}
		woken <- svc.AwaitReservation(context.Background(), done)
	}()
	<-started
	time.Sleep(20 * time.Millisecond)

	_, prob := svc.CommitCreate(context.Background(), req)
	if prob == nil || prob.Code != apiproblem.CodeInternalError || prob.Status != http.StatusInternalServerError {
		t.Fatalf("want INTERNAL_ERROR/500, got %#v", prob)
	}
	if err := <-woken; err != nil {
		t.Fatalf("waiter must wake after unlock: %v", err)
	}
	if _, ok := store.GetExecutionTarget(req.UID); ok {
		t.Fatal("failed append must not publish target")
	}
	if store.nameReserved(req.CloudProviderScopeUID, req.Name) {
		t.Fatal("failed append must not reserve name")
	}
	if store.liveTupleOccupied(req.ParticipationRef.UID, req.StackRef.UID, model.TargetClassSyntheticIaaS) {
		t.Fatal("failed append must not index live tuple")
	}
	if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
		t.Fatal("failed append must remove InFlight reservation (not leave unchanged)")
	}
	if idemp.completedCountForTest(svc.now()) != 0 {
		t.Fatal("failed append must create no completion")
	}
	if audit.Len() != 0 {
		t.Fatal("failed append must leave no durable AuditEvent")
	}
}

func TestLifecycle_CommitRetire_ActiveIncludingMaintenance(t *testing.T) {
	t.Parallel()

	svc, store, idemp, audit := newLifecycleForTest(t, nil)
	create := baseCreateReq("k-retire-create")
	create.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-1"}}`)
	if res := svc.ReserveOrInspectReplay(create.Idempotency, create.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve create: %#v", res)
	}
	created, prob := svc.CommitCreate(context.Background(), create)
	if prob != nil {
		t.Fatalf("create: %#v", prob)
	}

	// Simulate Maintenance-active marker and existing FactSet/Result links.
	svc.lock()
	cur, _ := store.lookupExecutionTarget(created.Metadata.UID)
	fsRef := apimeta.TypedRef{UID: "fs-1", Kind: model.KindNormalizedTargetFactSet, APIVersion: model.APIVersionFactSet, Name: "fs-1"}
	qrRef := apimeta.TypedRef{UID: "qr-1", Kind: model.KindTargetQualificationResult, APIVersion: model.APIVersionQualificationResult, Name: "qr-1"}
	cur.Status.Qualification = model.QualificationQualified
	cur.Status.FactSetRef = &fsRef
	cur.Status.QualificationResultRef = &qrRef
	staged, _ := store.stageUpdateExecutionTarget(cur)
	store.publish(staged)
	mStaged, _ := store.stageSetMaintenanceMarker(model.CurrentMaintenanceMarker{
		TargetUID: created.Metadata.UID, MaintenanceEpoch: 1, Active: true,
	})
	store.publish(mStaged)
	fsPut, _ := store.stagePutFactSet(model.NormalizedTargetFactSet{UID: "fs-1", TargetRef: subjectRef(created)})
	store.publish(fsPut)
	qrPut, _ := store.stagePutQualificationResult(model.TargetQualificationResult{UID: "qr-1", TargetRef: subjectRef(created)})
	store.publish(qrPut)
	before, _ := store.lookupExecutionTarget(created.Metadata.UID)
	svc.unlock()

	retire := RetireCommitRequest{
		TargetUID:               created.Metadata.UID,
		ExpectedResourceVersion: before.Metadata.ResourceVersion,
		Actor:                   create.Actor,
		RequestID:               "req-retire-1",
		AuditUID:                "audit-retire-1",
		Idempotency:             retireNS(created.Metadata.UID, "k-retire-ok"),
		Digest:                  DigestAction(),
		Completion:              CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{"retired":true}`)},
	}
	if res := svc.ReserveOrInspectReplay(retire.Idempotency, retire.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve retire: %#v", res)
	}

	got, prob := svc.CommitRetire(context.Background(), retire)
	if prob != nil {
		t.Fatalf("retire: %#v", prob)
	}
	if got.Status.Lifecycle != model.LifecycleRetired || got.Status.Qualification != model.QualificationUnqualified {
		t.Fatalf("retired state=%#v", got.Status)
	}
	if got.Status.FactSetRef != nil || got.Status.QualificationResultRef != nil {
		t.Fatalf("links must clear on retire: %#v", got.Status)
	}
	if store.liveTupleOccupied(create.ParticipationRef.UID, create.StackRef.UID, model.TargetClassSyntheticIaaS) {
		t.Fatal("live tuple must release on retire")
	}
	if !store.nameReserved(create.CloudProviderScopeUID, create.Name) {
		t.Fatal("name reservation must be retained after retire")
	}
	// Maintenance marker is not cleared by this task's retire path.
	if _, ok := store.GetMaintenanceMarker(created.Metadata.UID); !ok {
		t.Fatal("TASK-F16-04 retire must not clear Maintenance marker (scheduler owns clear)")
	}
	if audit.Last().Record.Action != AuditActionRetire {
		t.Fatalf("retire audit action=%q", audit.Last().Record.Action)
	}
	st, ok := idemp.stateForTest(retire.Idempotency, svc.now())
	if !ok || st != ReservationCompleted {
		t.Fatalf("retire completion=%s ok=%v", st, ok)
	}
}

func TestLifecycle_CommitRetire_AppendFailurePreservesLinksAndTuple(t *testing.T) {
	t.Parallel()

	svc, store, idemp, okAudit := newLifecycleForTest(t, nil)
	create := baseCreateReq("k-retire-fail-create")
	create.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-1"}}`)
	if res := svc.ReserveOrInspectReplay(create.Idempotency, create.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve create: %#v", res)
	}
	created, prob := svc.CommitCreate(context.Background(), create)
	if prob != nil {
		t.Fatalf("create: %#v", prob)
	}
	_ = okAudit

	svc.lock()
	cur, _ := store.lookupExecutionTarget(created.Metadata.UID)
	fsRef := apimeta.TypedRef{UID: "fs-keep", Kind: model.KindNormalizedTargetFactSet, APIVersion: model.APIVersionFactSet, Name: "fs-keep"}
	qrRef := apimeta.TypedRef{UID: "qr-keep", Kind: model.KindTargetQualificationResult, APIVersion: model.APIVersionQualificationResult, Name: "qr-keep"}
	cur.Status.FactSetRef = &fsRef
	cur.Status.QualificationResultRef = &qrRef
	cur.Status.Qualification = model.QualificationQualified
	staged, _ := store.stageUpdateExecutionTarget(cur)
	store.publish(staged)
	before, _ := store.lookupExecutionTarget(created.Metadata.UID)
	svc.unlock()

	failAudit := &memoryAuditAppender{fail: errors.New("retire append boom")}
	svc.audit = failAudit

	retire := RetireCommitRequest{
		TargetUID:               created.Metadata.UID,
		ExpectedResourceVersion: before.Metadata.ResourceVersion,
		Actor:                   create.Actor,
		RequestID:               "req-retire-fail",
		AuditUID:                "audit-retire-fail",
		Idempotency:             retireNS(created.Metadata.UID, "k-retire-fail"),
		Digest:                  DigestAction(),
		Completion:              CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{}`)},
	}
	if res := svc.ReserveOrInspectReplay(retire.Idempotency, retire.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve retire: %#v", res)
	}

	started := make(chan struct{})
	woken := make(chan error, 1)
	go func() {
		if res := svc.ReserveOrInspectReplay(retire.Idempotency, retire.Digest); res.Outcome != ReservationOutcomeInFlight {
			woken <- errors.New("expected InFlight")
			return
		}
		done, ok := svc.AttachWaiter(retire.Idempotency)
		close(started)
		if !ok {
			woken <- errors.New("attach failed")
			return
		}
		woken <- svc.AwaitReservation(context.Background(), done)
	}()
	<-started
	time.Sleep(20 * time.Millisecond)

	_, prob = svc.CommitRetire(context.Background(), retire)
	if prob == nil || prob.Code != apiproblem.CodeInternalError {
		t.Fatalf("want INTERNAL_ERROR, got %#v", prob)
	}
	if err := <-woken; err != nil {
		t.Fatalf("waiter wake: %v", err)
	}

	after, ok := store.GetExecutionTarget(created.Metadata.UID)
	if !ok {
		t.Fatal("target must remain")
	}
	if after.Status.Lifecycle != model.LifecycleActive {
		t.Fatalf("lifecycle must remain Active, got %q", after.Status.Lifecycle)
	}
	if after.Status.FactSetRef == nil || after.Status.FactSetRef.UID != "fs-keep" {
		t.Fatalf("FactSet link must be preserved: %#v", after.Status.FactSetRef)
	}
	if after.Status.QualificationResultRef == nil || after.Status.QualificationResultRef.UID != "qr-keep" {
		t.Fatalf("Result link must be preserved: %#v", after.Status.QualificationResultRef)
	}
	if after.Metadata.ResourceVersion != before.Metadata.ResourceVersion {
		t.Fatalf("ETag must be unchanged: before=%s after=%s", before.Metadata.ResourceVersion, after.Metadata.ResourceVersion)
	}
	if !store.liveTupleOccupied(create.ParticipationRef.UID, create.StackRef.UID, model.TargetClassSyntheticIaaS) {
		t.Fatal("live tuple must remain on append failure")
	}
	if _, ok := idemp.stateForTest(retire.Idempotency, svc.now()); ok {
		t.Fatal("failed retire append must remove InFlight reservation")
	}
	if idemp.completedCountForTest(svc.now()) != 1 { // create completion only
		t.Fatalf("completed count=%d want 1 (create only)", idemp.completedCountForTest(svc.now()))
	}
}

func TestLifecycle_CommitRetire_SyntheticInFlightReservationNoQualifyAbort(t *testing.T) {
	t.Parallel()

	// Proves retire's own idempotency discipline with a synthetic InFlight
	// peer waiter. No Qualifying-reservation mechanism is invoked (TASK-F16-05).
	svc, _, _, _ := newLifecycleForTest(t, nil)
	create := baseCreateReq("k-synth-create")
	create.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-1"}}`)
	if res := svc.ReserveOrInspectReplay(create.Idempotency, create.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve create: %#v", res)
	}
	created, prob := svc.CommitCreate(context.Background(), create)
	if prob != nil {
		t.Fatalf("create: %#v", prob)
	}

	retire := RetireCommitRequest{
		TargetUID:               created.Metadata.UID,
		ExpectedResourceVersion: created.Metadata.ResourceVersion,
		Actor:                   create.Actor,
		RequestID:               "req-synth-retire",
		AuditUID:                "audit-synth-retire",
		Idempotency:             retireNS(created.Metadata.UID, "k-synth-retire"),
		Digest:                  DigestAction(),
		Completion:              CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{}`)},
	}
	if res := svc.ReserveOrInspectReplay(retire.Idempotency, retire.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve retire: %#v", res)
	}

	// Synthetic peer waiter on the retire idempotency key — never a qualify abort.
	started := make(chan struct{})
	woken := make(chan error, 1)
	go func() {
		if res := svc.ReserveOrInspectReplay(retire.Idempotency, retire.Digest); res.Outcome != ReservationOutcomeInFlight {
			woken <- errors.New("expected InFlight")
			return
		}
		done, ok := svc.AttachWaiter(retire.Idempotency)
		close(started)
		if !ok {
			woken <- errors.New("attach failed")
			return
		}
		woken <- svc.AwaitReservation(context.Background(), done)
	}()
	<-started
	time.Sleep(20 * time.Millisecond)

	got, prob := svc.CommitRetire(context.Background(), retire)
	if prob != nil {
		t.Fatalf("retire: %#v", prob)
	}
	if got.Status.Lifecycle != model.LifecycleRetired {
		t.Fatalf("lifecycle=%q", got.Status.Lifecycle)
	}
	if err := <-woken; err != nil {
		t.Fatalf("synthetic waiter wake: %v", err)
	}
}

func TestLifecycle_OwnerCancellation_CreateAndRetire(t *testing.T) {
	t.Parallel()

	svc, store, idemp, audit := newLifecycleForTest(t, nil)

	t.Run("create", func(t *testing.T) {
		req := baseCreateReq("k-cancel-create")
		req.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-1"}}`)
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, prob := svc.CommitCreate(ctx, req)
		if prob == nil || prob.Code != apiproblem.CodeInternalError {
			t.Fatalf("cancelled create: %#v", prob)
		}
		if _, ok := store.GetExecutionTarget(req.UID); ok {
			t.Fatal("cancelled create must not publish")
		}
		if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
			t.Fatal("cancelled create must remove InFlight reservation")
		}
		if audit.Len() != 0 {
			t.Fatal("cancelled create must not audit")
		}
	})

	t.Run("retire", func(t *testing.T) {
		create := baseCreateReq("k-cancel-retire-create")
		create.UID = "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		create.Name = "target-2"
		create.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-2"}}`)
		if res := svc.ReserveOrInspectReplay(create.Idempotency, create.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve create: %#v", res)
		}
		created, prob := svc.CommitCreate(context.Background(), create)
		if prob != nil {
			t.Fatalf("create: %#v", prob)
		}
		beforeAudit := audit.Len()

		retire := RetireCommitRequest{
			TargetUID:               created.Metadata.UID,
			ExpectedResourceVersion: created.Metadata.ResourceVersion,
			Actor:                   create.Actor,
			RequestID:               "req-cancel-retire",
			AuditUID:                "audit-cancel-retire",
			Idempotency:             retireNS(created.Metadata.UID, "k-cancel-retire"),
			Digest:                  DigestAction(),
			Completion:              CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{}`)},
		}
		if res := svc.ReserveOrInspectReplay(retire.Idempotency, retire.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve retire: %#v", res)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, prob = svc.CommitRetire(ctx, retire)
		if prob == nil || prob.Code != apiproblem.CodeInternalError {
			t.Fatalf("cancelled retire: %#v", prob)
		}
		after, _ := store.GetExecutionTarget(created.Metadata.UID)
		if after.Status.Lifecycle != model.LifecycleActive {
			t.Fatalf("cancelled retire must not publish retirement: %#v", after.Status)
		}
		if _, ok := idemp.stateForTest(retire.Idempotency, svc.now()); ok {
			t.Fatal("cancelled retire must remove InFlight reservation")
		}
		if audit.Len() != beforeAudit {
			t.Fatal("cancelled retire must not append retirement AuditEvent")
		}
	})
}

func TestLifecycle_OwnerPanic_CreateAndRetireCleanup(t *testing.T) {
	t.Parallel()

	svc, store, idemp, _ := newLifecycleForTest(t, nil)

	t.Run("create-panic", func(t *testing.T) {
		req := baseCreateReq("k-panic-create")
		req.UID = "ffffffffffffffffffffffffffffffff"
		req.Name = "target-panic-create"
		req.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-panic-create"}}`)
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}

		// Inject a panicking appender.
		svc.audit = &panicAuditAppender{}
		panicked := false
		func() {
			defer func() {
				if recover() != nil {
					panicked = true
				}
			}()
			_, _ = svc.CommitCreate(context.Background(), req)
		}()
		if !panicked {
			t.Fatal("expected panic from appender")
		}
		if _, ok := store.GetExecutionTarget(req.UID); ok {
			t.Fatal("panic create must not publish")
		}
		if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
			t.Fatal("panic create must remove InFlight reservation")
		}
	})

	t.Run("retire-panic", func(t *testing.T) {
		svc.audit = &memoryAuditAppender{}
		create := baseCreateReq("k-panic-retire-create")
		create.UID = "11111111111111111111111111111111"
		create.Name = "target-panic-retire"
		create.ParticipationRef = participationRef("22222222222222222222222222222222")
		create.StackRef = stackRef("33333333333333333333333333333333")
		create.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-panic-retire"}}`)
		if res := svc.ReserveOrInspectReplay(create.Idempotency, create.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve create: %#v", res)
		}
		created, prob := svc.CommitCreate(context.Background(), create)
		if prob != nil {
			t.Fatalf("create: %#v", prob)
		}

		retire := RetireCommitRequest{
			TargetUID:               created.Metadata.UID,
			ExpectedResourceVersion: created.Metadata.ResourceVersion,
			Actor:                   create.Actor,
			RequestID:               "req-panic-retire",
			AuditUID:                "audit-panic-retire",
			Idempotency:             retireNS(created.Metadata.UID, "k-panic-retire"),
			Digest:                  DigestAction(),
			Completion:              CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{}`)},
		}
		if res := svc.ReserveOrInspectReplay(retire.Idempotency, retire.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve retire: %#v", res)
		}
		svc.audit = &panicAuditAppender{}
		panicked := false
		func() {
			defer func() {
				if recover() != nil {
					panicked = true
				}
			}()
			_, _ = svc.CommitRetire(context.Background(), retire)
		}()
		if !panicked {
			t.Fatal("expected panic from appender")
		}
		after, _ := store.GetExecutionTarget(created.Metadata.UID)
		if after.Status.Lifecycle != model.LifecycleActive {
			t.Fatalf("panic retire must not publish: %#v", after.Status)
		}
		if _, ok := idemp.stateForTest(retire.Idempotency, svc.now()); ok {
			t.Fatal("panic retire must remove InFlight reservation")
		}
	})
}

type panicAuditAppender struct{}

func (p *panicAuditAppender) Append(context.Context, apiconform.AuditEvent) error {
	panic("injected audit panic")
}

func TestLifecycle_ReservationAPI_HandlersCannotComplete(t *testing.T) {
	t.Parallel()

	svc, _, idemp, _ := newLifecycleForTest(t, nil)
	ns := createNS("k-api")
	digest := mustCreateDigest(t, `{"metadata":{"name":"x"}}`)

	res := svc.ReserveOrInspectReplay(ns, digest)
	if res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	// No exported CompleteReservation — completion is commit-internal only.
	rt := reflect.TypeOf(svc)
	if _, ok := rt.MethodByName("CompleteReservation"); ok {
		t.Fatal("CompleteReservation must not be exported; completion is commit-internal")
	}

	if !svc.AbortReservation(ns) {
		t.Fatal("AbortReservation must remove InFlight")
	}
	if _, ok := idemp.stateForTest(ns, svc.now()); ok {
		t.Fatal("aborted reservation must be gone")
	}
}

func TestLifecycle_NameAndTupleUniquenessThroughSoleCommitter(t *testing.T) {
	t.Parallel()

	svc, store, _, _ := newLifecycleForTest(t, nil)
	first := baseCreateReq("k-uniq-1")
	first.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-1"}}`)
	if res := svc.ReserveOrInspectReplay(first.Idempotency, first.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	if _, prob := svc.CommitCreate(context.Background(), first); prob != nil {
		t.Fatalf("first create: %#v", prob)
	}

	dupName := baseCreateReq("k-uniq-2")
	dupName.UID = "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
	dupName.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-1-dup"}}`)
	if res := svc.ReserveOrInspectReplay(dupName.Idempotency, dupName.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve dup name: %#v", res)
	}
	_, prob := svc.CommitCreate(context.Background(), dupName)
	if prob == nil || prob.Code != apiproblem.CodeAlreadyExists {
		t.Fatalf("duplicate name: %#v", prob)
	}

	dupTuple := baseCreateReq("k-uniq-3")
	dupTuple.UID = "ffffffffffffffffffffffffffffffff"
	dupTuple.Name = "target-other"
	dupTuple.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-other"}}`)
	if res := svc.ReserveOrInspectReplay(dupTuple.Idempotency, dupTuple.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve dup tuple: %#v", res)
	}
	_, prob = svc.CommitCreate(context.Background(), dupTuple)
	if prob == nil || prob.Code != apiproblem.CodeAlreadyExists {
		t.Fatalf("duplicate live tuple: %#v", prob)
	}

	if got := store.ListExecutionTargets(); len(got) != 1 {
		t.Fatalf("want one published target, got %d", len(got))
	}
}

func TestLifecycle_ConcurrentCreatesSerialize(t *testing.T) {
	t.Parallel()

	svc, store, _, _ := newLifecycleForTest(t, nil)
	var success atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req := baseCreateReq("k-race")
			// Same name/tuple — only one may publish.
			req.UID = "dddddddddddddddddddddddddddddddd"
			req.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-1"}}`)
			_ = i
			res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest)
			switch res.Outcome {
			case ReservationOutcomeReserved:
				if _, prob := svc.CommitCreate(context.Background(), req); prob == nil {
					success.Add(1)
				}
			case ReservationOutcomeInFlight:
				done, ok := svc.AttachWaiter(req.Idempotency)
				if !ok {
					return
				}
				_ = svc.AwaitReservation(context.Background(), done)
			case ReservationOutcomeReplay, ReservationOutcomeConflict:
				return
			}
		}(i)
	}
	wg.Wait()
	if n := len(store.ListExecutionTargets()); n != 1 {
		t.Fatalf("sole committer must publish exactly one target, got %d (successes=%d)", n, success.Load())
	}
}
