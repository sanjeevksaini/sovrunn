package executiontarget

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

func newSchedulerFixture(t *testing.T) (*ExecutionTargetLifecycleService, *Store, *IdempotencyTable, *memoryAuditAppender, *FakeClock, *Scheduler) {
	t.Helper()
	clock := NewFakeClock(time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC))
	store := NewStore()
	idemp := NewIdempotencyTable()
	audit := &memoryAuditAppender{}
	svc := NewExecutionTargetLifecycleService(LifecycleConfig{
		Store:       store,
		Idempotency: idemp,
		Audit:       audit,
		Clock:       clock,
	})
	sched := NewScheduler(SchedulerConfig{Lifecycle: svc, Clock: clock})
	return svc, store, idemp, audit, clock, sched
}

func createTargetForScheduler(t *testing.T, svc *ExecutionTargetLifecycleService, name, uid, key string) model.ExecutionTarget {
	t.Helper()
	req := baseCreateReq(key)
	req.Name = name
	req.UID = uid
	// Unique live tuple per target: participation UID and stack UID derived from
	// the target UID so concurrent fixtures in one store do not collide.
	partUID := uid[:16] + "aaaaaaaaaaaaaaaa"
	stackUID := uid[:16] + "bbbbbbbbbbbbbbbb"
	req.ParticipationRef = participationRef(partUID)
	req.StackRef = stackRef(stackUID)
	req.Digest = mustCreateDigest(t, `{"metadata":{"name":"`+name+`"}}`)
	req.Idempotency = createNS(key)
	req.Idempotency.CloudProviderUID = req.CloudProviderScopeUID
	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve create: %#v", res)
	}
	et, prob := svc.CommitCreate(context.Background(), req)
	if prob != nil {
		t.Fatalf("create: %#v", prob)
	}
	return et
}

func plantExpiredFactSet(t *testing.T, store *Store, et model.ExecutionTarget, fsUID string, expiresAt time.Time) model.ExecutionTarget {
	t.Helper()
	fs := model.NormalizedTargetFactSet{
		UID:                           fsUID,
		TargetRef:                     subjectRef(et),
		InfrastructureStackGeneration: et.Status.ObservedGeneration,
		MaintenanceEpoch:              et.Status.MaintenanceEpoch,
		ViabilityFingerprint:          "fp",
		ObserverID:                    model.ObserverID,
		ObserverRevision:              "1",
		FactVersion:                   model.FactVersionV1,
		ObservedAt:                    expiresAt.Add(-FactFreshnessWindow).UTC().Format(time.RFC3339),
		ExpiresAt:                     expiresAt.UTC().Format(time.RFC3339),
		Facts: model.FactSetTruths{
			Compute: model.FactComputeTruths{VM: model.FactSupported},
			Storage: model.FactStorageTruths{Block: model.FactSupported, Object: model.FactSupported},
			Network: model.FactNetworkTruths{Private: model.FactSupported},
		},
	}
	qrUID := "qr-" + fsUID
	qr := model.TargetQualificationResult{
		UID:       qrUID,
		TargetRef: subjectRef(et),
		FactSetRef: apimeta.TypedRef{
			APIVersion: model.APIVersionFactSet,
			Kind:       model.KindNormalizedTargetFactSet,
			Name:       fsUID,
			UID:        fsUID,
		},
		InfrastructureStackGeneration: et.Status.ObservedGeneration,
		MaintenanceEpoch:              et.Status.MaintenanceEpoch,
		ProfileVersion:                model.ProfileVersion,
		Outcome:                       model.OutcomeQualified,
		ReasonCodes:                   []string{"ok"},
		EvaluatedAt:                   expiresAt.Add(-FactFreshnessWindow).UTC().Format(time.RFC3339),
	}

	store.beginPublication()
	cur, _ := store.lookupExecutionTarget(et.Metadata.UID)
	fsRef := apimeta.TypedRef{
		APIVersion: model.APIVersionFactSet,
		Kind:       model.KindNormalizedTargetFactSet,
		Name:       fsUID,
		UID:        fsUID,
	}
	qrRef := apimeta.TypedRef{
		APIVersion: model.APIVersionQualificationResult,
		Kind:       model.KindTargetQualificationResult,
		Name:       qrUID,
		UID:        qrUID,
	}
	cur.Status.Qualification = model.QualificationQualified
	cur.Status.FactSetRef = &fsRef
	cur.Status.QualificationResultRef = &qrRef
	stagedFS, _ := store.stagePutFactSet(fs)
	stagedQR, _ := store.stagePutQualificationResult(qr)
	stagedET, _ := store.stageUpdateExecutionTarget(cur)
	store.publish(stagedFS)
	store.publish(stagedQR)
	store.publish(stagedET)
	out, _ := store.lookupExecutionTarget(et.Metadata.UID)
	store.endPublication()
	return out
}

func baseMaintenanceTrigger(et model.ExecutionTarget, op MaintenanceTriggerOperation, auditUID string) MaintenanceTrigger {
	return MaintenanceTrigger{
		TargetUID:                      et.Metadata.UID,
		ExpectedInfrastructureStackGen: et.Status.ObservedGeneration,
		ExpectedMaintenanceEpoch:       et.Status.MaintenanceEpoch,
		Operation:                      op,
		Actor:                          MaintenanceSystemActor(),
		RequestID:                      "corr-" + auditUID,
		AuditUID:                       auditUID,
	}
}

func TestScheduler_ExpirySuccess_UnchangedETagAndFreshness(t *testing.T) {
	t.Parallel()

	svc, store, _, audit, clock, sched := newSchedulerFixture(t)
	et := createTargetForScheduler(t, svc, "t-exp", "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", "k-exp")
	expires := clock.Now().Add(-time.Second)
	et = plantExpiredFactSet(t, store, et, "fs-exp-1", expires)
	beforeRV := et.Metadata.ResourceVersion
	beforeQual := et.Status.Qualification
	beforeAudit := audit.Len()
	fsUID := et.Status.FactSetRef.UID

	sched.ProcessExpiryTick(context.Background())

	after, _ := store.GetExecutionTarget(et.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV {
		t.Fatal("expiry must leave ETag unchanged")
	}
	if after.Status.Qualification != beforeQual {
		t.Fatal("expiry must leave qualification unchanged")
	}
	if after.Status.FactSetRef == nil || after.Status.FactSetRef.UID != fsUID {
		t.Fatal("expiry must leave FactSet link unchanged")
	}
	if !svc.FactSetExpired(fsUID) {
		t.Fatal("freshness input must mark fact set expired")
	}
	if !svc.FactSetExpiryAudited(fsUID) {
		t.Fatal("expiry audit must succeed once")
	}
	if audit.Len() != beforeAudit+1 {
		t.Fatalf("want exactly one expiry AuditEvent, audit=%d", audit.Len())
	}
	if audit.Last().Record.Action != AuditActionExpiry {
		t.Fatalf("action=%s", audit.Last().Record.Action)
	}
	if sched.ExternalCallCount() != 0 || svc.Observer().ExternalCallCount() != 0 {
		t.Fatal("externalCallCount must remain 0")
	}
}

func TestScheduler_ExpiryAuditRetryOncePerClockSecond(t *testing.T) {
	t.Parallel()

	svc, store, _, audit, clock, sched := newSchedulerFixture(t)
	et := createTargetForScheduler(t, svc, "t-retry", "ffffffffffffffffffffffffffffffff", "k-retry")
	et = plantExpiredFactSet(t, store, et, "fs-retry-1", clock.Now().Add(-time.Second))
	fsUID := et.Status.FactSetRef.UID
	beforeRV := et.Metadata.ResourceVersion
	audit.SetFail(errors.New("append boom"))

	sched.ProcessExpiryTick(context.Background())
	if svc.FactSetExpiryAudited(fsUID) {
		t.Fatal("failed append must not mark audited")
	}
	if !svc.FactSetExpired(fsUID) {
		t.Fatal("freshness must still flip on failed audit")
	}
	after, _ := store.GetExecutionTarget(et.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV {
		t.Fatal("failed expiry audit must leave ETag unchanged")
	}
	if audit.Len() != 1 { // create only
		t.Fatalf("failed append must not retain expiry event, len=%d", audit.Len())
	}
	if audit.Calls() < 2 { // create + failed expiry attempt
		t.Fatalf("want expiry append attempt, calls=%d", audit.Calls())
	}

	clock.Advance(time.Second)
	sched.ProcessExpiryTick(context.Background())
	if svc.FactSetExpiryAudited(fsUID) || audit.Len() != 1 {
		t.Fatal("still failing: must retry without publishing event")
	}

	audit.SetFail(nil)
	clock.Advance(time.Second)
	sched.ProcessExpiryTick(context.Background())
	if !svc.FactSetExpiryAudited(fsUID) {
		t.Fatal("successful retry must mark audited")
	}
	if audit.Len() != 2 || audit.Last().Record.Action != AuditActionExpiry {
		t.Fatalf("want one expiry event after success, len=%d action=%s", audit.Len(), audit.Last().Record.Action)
	}

	clock.Advance(time.Second)
	beforeCalls := audit.Calls()
	sched.ProcessExpiryTick(context.Background())
	if audit.Calls() != beforeCalls {
		t.Fatal("already-audited expiry must not append again")
	}
}

func TestScheduler_RetiredAndMaintenanceOutrankExpiry(t *testing.T) {
	t.Parallel()

	t.Run("retired", func(t *testing.T) {
		svc, store, _, audit, clock, sched := newSchedulerFixture(t)
		et := createTargetForScheduler(t, svc, "t-ret", "11111111111111111111111111111111", "k-ret-exp")
		et = plantExpiredFactSet(t, store, et, "fs-ret", clock.Now().Add(-time.Second))
		fsUID := et.Status.FactSetRef.UID

		rreq := RetireCommitRequest{
			TargetUID:   et.Metadata.UID,
			Actor:       MaintenanceSystemActor(),
			RequestID:   "req-retire",
			AuditUID:    "audit-retire",
			Idempotency: retireNS(et.Metadata.UID, "k-ret-exp"),
			Digest:      DigestAction(),
			Completion:  CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{}`)},
		}
		if res := svc.ReserveOrInspectReplay(rreq.Idempotency, rreq.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve retire: %#v", res)
		}
		if _, prob := svc.CommitRetire(context.Background(), rreq); prob != nil {
			t.Fatalf("retire: %#v", prob)
		}
		beforeAudit := audit.Len()
		sched.ProcessExpiryTick(context.Background())
		if svc.FactSetExpired(fsUID) || svc.FactSetExpiryAudited(fsUID) {
			t.Fatal("Retired must outrank expiry")
		}
		if audit.Len() != beforeAudit {
			t.Fatal("Retired expiry must not append")
		}
	})

	t.Run("maintenance", func(t *testing.T) {
		svc, store, _, audit, clock, sched := newSchedulerFixture(t)
		et := createTargetForScheduler(t, svc, "t-mnt", "22222222222222222222222222222222", "k-mnt-exp")
		et = plantExpiredFactSet(t, store, et, "fs-mnt", clock.Now().Add(-time.Second))
		beforeAudit := audit.Len()

		trig := baseMaintenanceTrigger(et, MaintenanceTriggerEnter, "audit-mnt-outrank")
		if _, prob := sched.DeliverMaintenanceTrigger(context.Background(), trig); prob != nil {
			t.Fatalf("enter: %#v", prob)
		}
		afterEnter, _ := store.GetExecutionTarget(et.Metadata.UID)
		afterEnter = plantExpiredFactSet(t, store, afterEnter, "fs-mnt-2", clock.Now().Add(-time.Second))
		fsUID2 := afterEnter.Status.FactSetRef.UID
		beforeTickAudit := audit.Len()
		sched.ProcessExpiryTick(context.Background())
		if svc.FactSetExpired(fsUID2) || svc.FactSetExpiryAudited(fsUID2) {
			t.Fatal("Maintenance must outrank expiry")
		}
		if audit.Len() != beforeTickAudit {
			t.Fatal("Maintenance expiry must not append")
		}
		if audit.Len() <= beforeAudit {
			t.Fatal("maintenance enter itself must still be audited")
		}
	})
}

func TestScheduler_MaintenanceEnterClear_HappyAndStale(t *testing.T) {
	t.Parallel()

	svc, store, _, audit, _, sched := newSchedulerFixture(t)
	et := createTargetForScheduler(t, svc, "t-mc", "33333333333333333333333333333333", "k-mc")
	et = plantExpiredFactSet(t, store, et, "fs-mc", time.Now().UTC().Add(time.Hour))
	beforeRV := et.Metadata.ResourceVersion
	beforeEpoch := et.Status.MaintenanceEpoch

	stale := baseMaintenanceTrigger(et, MaintenanceTriggerEnter, "audit-stale-enter")
	stale.ExpectedMaintenanceEpoch = beforeEpoch + 1
	beforeAudit := audit.Len()
	got, prob := sched.DeliverMaintenanceTrigger(context.Background(), stale)
	if prob != nil {
		t.Fatalf("stale enter must be silent no-op, got %#v", prob)
	}
	if got.Metadata.ResourceVersion != beforeRV {
		t.Fatal("stale enter must not advance RV")
	}
	if audit.Len() != beforeAudit {
		t.Fatal("stale enter must not audit")
	}
	if _, ok := store.GetMaintenanceMarker(et.Metadata.UID); ok {
		t.Fatal("stale enter must not set marker")
	}

	enter := baseMaintenanceTrigger(et, MaintenanceTriggerEnter, "audit-enter-ok")
	got, prob = sched.DeliverMaintenanceTrigger(context.Background(), enter)
	if prob != nil {
		t.Fatalf("enter: %#v", prob)
	}
	if got.Status.MaintenanceEpoch != beforeEpoch+1 {
		t.Fatalf("epoch=%d", got.Status.MaintenanceEpoch)
	}
	if got.Metadata.ResourceVersion == beforeRV {
		t.Fatal("enter must advance resourceVersion")
	}
	if got.Status.Qualification != model.QualificationUnqualified {
		t.Fatalf("qualification=%s", got.Status.Qualification)
	}
	if got.Status.FactSetRef != nil || got.Status.QualificationResultRef != nil {
		t.Fatal("enter must clear both links")
	}
	marker, ok := store.GetMaintenanceMarker(et.Metadata.UID)
	if !ok || !marker.Active || marker.MaintenanceEpoch != got.Status.MaintenanceEpoch {
		t.Fatalf("marker=%#v ok=%v", marker, ok)
	}
	if audit.Last().Record.Action != AuditActionMaintenanceEnter {
		t.Fatalf("action=%s", audit.Last().Record.Action)
	}

	staleClear := baseMaintenanceTrigger(got, MaintenanceTriggerClear, "audit-stale-clear")
	staleClear.ExpectedInfrastructureStackGen = got.Status.ObservedGeneration + 9
	beforeClearAudit := audit.Len()
	beforeClearRV := got.Metadata.ResourceVersion
	got2, prob := sched.DeliverMaintenanceTrigger(context.Background(), staleClear)
	if prob != nil {
		t.Fatalf("stale clear: %#v", prob)
	}
	if got2.Metadata.ResourceVersion != beforeClearRV || audit.Len() != beforeClearAudit {
		t.Fatal("stale clear must not mutate or audit")
	}

	clear := baseMaintenanceTrigger(got, MaintenanceTriggerClear, "audit-clear-ok")
	got3, prob := sched.DeliverMaintenanceTrigger(context.Background(), clear)
	if prob != nil {
		t.Fatalf("clear: %#v", prob)
	}
	if got3.Status.MaintenanceEpoch != got.Status.MaintenanceEpoch+1 {
		t.Fatalf("clear epoch=%d", got3.Status.MaintenanceEpoch)
	}
	if got3.Status.FactSetRef != nil || got3.Status.QualificationResultRef != nil {
		t.Fatal("clear must clear links")
	}
	if got3.Status.Qualification != model.QualificationUnqualified {
		t.Fatalf("qualification=%s", got3.Status.Qualification)
	}
	if _, ok := store.GetMaintenanceMarker(et.Metadata.UID); ok {
		t.Fatal("clear must remove marker")
	}
	if audit.Last().Record.Action != AuditActionMaintenanceClear {
		t.Fatalf("action=%s", audit.Last().Record.Action)
	}
	if sched.ExternalCallCount() != 0 {
		t.Fatal("externalCallCount must remain 0")
	}
}

func TestScheduler_MaintenanceEnter_AuditFailure_F16_100(t *testing.T) {
	t.Parallel()

	svc, store, idemp, audit, _, sched := newSchedulerFixture(t)
	et := createTargetForScheduler(t, svc, "t-100", "44444444444444444444444444444444", "k-100")
	et = plantExpiredFactSet(t, store, et, "fs-100", time.Now().UTC().Add(time.Hour))
	beforeRV := et.Metadata.ResourceVersion
	beforeEpoch := et.Status.MaintenanceEpoch
	beforeFS := et.Status.FactSetRef.UID
	beforeAudit := audit.Len()

	mapSrc := NewMapFixtureSource()
	started := make(chan struct{})
	release := make(chan struct{})
	blocking := &blockingFixtureSource{inner: mapSrc, started: started, release: release}
	svc.observer = NewSyntheticObserver(blocking, svc.now)
	mapSrc.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
	qreq := baseQualifyReq(et.Metadata.UID, "k-100-q")
	if res := svc.ReserveOrInspectReplay(qreq.Idempotency, qreq.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve qualify: %#v", res)
	}
	errCh := make(chan *apiproblem.Problem, 1)
	go func() {
		_, prob := svc.Qualify(context.Background(), qreq)
		errCh <- prob
	}()
	<-started

	audit.SetFail(errors.New("enter append boom"))
	trig := baseMaintenanceTrigger(et, MaintenanceTriggerEnter, "audit-f16-100")
	_, prob := sched.DeliverMaintenanceTrigger(context.Background(), trig)
	if prob == nil || prob.Code != apiproblem.CodeInternalError {
		t.Fatalf("want INTERNAL_ERROR F16-100, got %#v", prob)
	}

	after, _ := store.GetExecutionTarget(et.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV || after.Status.MaintenanceEpoch != beforeEpoch {
		t.Fatal("failed enter must not advance epoch/ETag")
	}
	if after.Status.FactSetRef == nil || after.Status.FactSetRef.UID != beforeFS {
		t.Fatal("failed enter must preserve FactSet link")
	}
	if _, ok := store.GetMaintenanceMarker(et.Metadata.UID); ok {
		t.Fatal("failed enter must not set marker")
	}
	if audit.Len() != beforeAudit {
		t.Fatal("failed enter must not retain AuditEvent")
	}
	if st, ok := idemp.stateForTest(qreq.Idempotency, svc.now()); !ok || st != ReservationInFlight {
		t.Fatalf("failed enter must not abort in-flight qualification, state=%v ok=%v", st, ok)
	}

	audit.SetFail(nil) // clear so the surviving qualification can append
	close(release)
	qprob := <-errCh
	if qprob != nil {
		t.Fatalf("qualification should complete after failed enter left it intact: %#v", qprob)
	}
}

func TestScheduler_MaintenanceClear_AuditFailure_F16_101(t *testing.T) {
	t.Parallel()

	svc, store, _, audit, _, sched := newSchedulerFixture(t)
	et := createTargetForScheduler(t, svc, "t-101", "55555555555555555555555555555555", "k-101")
	enter := baseMaintenanceTrigger(et, MaintenanceTriggerEnter, "audit-enter-101")
	et, prob := sched.DeliverMaintenanceTrigger(context.Background(), enter)
	if prob != nil {
		t.Fatalf("enter: %#v", prob)
	}
	beforeRV := et.Metadata.ResourceVersion
	beforeEpoch := et.Status.MaintenanceEpoch
	beforeAudit := audit.Len()

	audit.SetFail(errors.New("clear append boom"))
	clear := baseMaintenanceTrigger(et, MaintenanceTriggerClear, "audit-f16-101")
	_, prob = sched.DeliverMaintenanceTrigger(context.Background(), clear)
	if prob == nil || prob.Code != apiproblem.CodeInternalError {
		t.Fatalf("want INTERNAL_ERROR F16-101, got %#v", prob)
	}
	after, _ := store.GetExecutionTarget(et.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV || after.Status.MaintenanceEpoch != beforeEpoch {
		t.Fatal("failed clear must not advance epoch/ETag")
	}
	marker, ok := store.GetMaintenanceMarker(et.Metadata.UID)
	if !ok || !marker.Active {
		t.Fatal("failed clear must leave marker active")
	}
	if audit.Len() != beforeAudit {
		t.Fatal("failed clear must not retain clear AuditEvent")
	}
}

func TestScheduler_MaintenanceWins_QualifyingCallerOutcome(t *testing.T) {
	t.Parallel()

	svc, store, idemp, audit, _, sched := newSchedulerFixture(t)
	et := createTargetForScheduler(t, svc, "t-89", "66666666666666666666666666666666", "k-89")
	et = plantExpiredFactSet(t, store, et, "fs-89", time.Now().UTC().Add(time.Hour))

	mapSrc := NewMapFixtureSource()
	started := make(chan struct{})
	release := make(chan struct{})
	blocking := &blockingFixtureSource{inner: mapSrc, started: started, release: release}
	svc.observer = NewSyntheticObserver(blocking, svc.now)
	mapSrc.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))

	qreq := baseQualifyReq(et.Metadata.UID, "k-89-q")
	if res := svc.ReserveOrInspectReplay(qreq.Idempotency, qreq.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	errCh := make(chan *apiproblem.Problem, 1)
	go func() {
		_, prob := svc.Qualify(context.Background(), qreq)
		errCh <- prob
	}()
	<-started

	other := createTargetForScheduler(t, svc, "t-89b", "77777777777777777777777777777777", "k-89b")
	otherQ := baseQualifyReq(other.Metadata.UID, "k-89-other")
	if res := svc.ReserveOrInspectReplay(otherQ.Idempotency, otherQ.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve other: %#v", res)
	}

	beforeQualifyAudit := audit.Len()
	trig := baseMaintenanceTrigger(et, MaintenanceTriggerEnter, "audit-maint-wins")
	got, prob := sched.DeliverMaintenanceTrigger(context.Background(), trig)
	if prob != nil {
		t.Fatalf("enter: %#v", prob)
	}
	if got.Status.Qualification != model.QualificationUnqualified || got.Status.FactSetRef != nil {
		t.Fatal("enter must clear links and persist Unqualified")
	}
	close(release)

	qprob := <-errCh
	if qprob == nil || qprob.Code != apiproblem.CodeConflict {
		t.Fatalf("want CONFLICT, got %#v", qprob)
	}
	if len(qprob.Violations) == 0 || qprob.Violations[0].Code != ViolationTargetMaintenance {
		t.Fatalf("want VS0_TARGET_MAINTENANCE, got %#v", qprob.Violations)
	}
	if _, ok := idemp.stateForTest(qreq.Idempotency, svc.now()); ok {
		t.Fatal("maintenance-wins must abort linked qualify idempotency reservation")
	}
	if st, ok := idemp.stateForTest(otherQ.Idempotency, svc.now()); !ok || st != ReservationInFlight {
		t.Fatalf("other target reservation must be untouched, state=%v ok=%v", st, ok)
	}
	foundQualify := false
	foundMaint := false
	audit.mu.Lock()
	for _, ev := range audit.events[beforeQualifyAudit:] {
		if ev.Record.Action == AuditActionQualify {
			foundQualify = true
		}
		if ev.Record.Action == AuditActionMaintenanceEnter {
			foundMaint = true
		}
	}
	audit.mu.Unlock()
	if foundQualify {
		t.Fatal("maintenance-wins must not emit qualification AuditEvent")
	}
	if !foundMaint {
		t.Fatal("winning maintenance entry must be audited")
	}
	svc.AbortReservation(otherQ.Idempotency)
}

func TestLifecycle_Shutdown_ClosesAdmissionAndWakesWaiters(t *testing.T) {
	t.Parallel()

	svc, _, idemp, _, _, sched := newSchedulerFixture(t)
	et := createTargetForScheduler(t, svc, "t-sd", "88888888888888888888888888888888", "k-sd")

	nsCreate := createNS("k-sd-inflight")
	nsQualify := qualifyNS(et.Metadata.UID, "k-sd-q")
	nsRetire := retireNS(et.Metadata.UID, "k-sd-r")
	for _, ns := range []IdempotencyNamespace{nsCreate, nsQualify, nsRetire} {
		if res := svc.ReserveOrInspectReplay(ns, DigestAction()); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve %#v: %#v", ns, res)
		}
	}
	done, ok := svc.AttachWaiter(nsCreate)
	if !ok {
		t.Fatal("attach")
	}
	woke := make(chan struct{})
	go func() {
		_ = svc.AwaitReservation(context.Background(), done)
		close(woke)
	}()

	sched.StopAcceptance()
	if !sched.AcceptanceStopped() {
		t.Fatal("acceptance must stop first")
	}
	trig := baseMaintenanceTrigger(et, MaintenanceTriggerEnter, "audit-post-stop")
	_, prob := sched.DeliverMaintenanceTrigger(context.Background(), trig)
	if prob == nil {
		t.Fatal("post-stop trigger must not mutate")
	}

	svc.Shutdown()
	if !svc.AdmissionClosed() {
		t.Fatal("admission must be closed")
	}
	select {
	case <-woke:
	case <-time.After(2 * time.Second):
		t.Fatal("shutdown must wake waiters after unlock")
	}
	if idemp.inFlightCountForTest() != 0 {
		t.Fatalf("all InFlight must be aborted, count=%d", idemp.inFlightCountForTest())
	}

	res := svc.ReserveOrInspectReplay(createNS("k-sd-after"), DigestAction())
	if res.Outcome == ReservationOutcomeReserved {
		t.Fatal("no InFlight reservation may appear after shutdown sweep")
	}
	if idemp.inFlightCountForTest() != 0 {
		t.Fatal("post-stop reserve must create no InFlight")
	}
}

func TestLifecycle_Shutdown_ConcurrentNoNewInFlight(t *testing.T) {
	t.Parallel()

	svc, _, idemp, _, _, _ := newSchedulerFixture(t)
	var wg sync.WaitGroup
	start := make(chan struct{})

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			ns := createNS("k-conc-" + string(rune('a'+i)))
			_ = svc.ReserveOrInspectReplay(ns, DigestAction())
		}(i)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		svc.Shutdown()
	}()
	close(start)
	wg.Wait()

	if idemp.inFlightCountForTest() != 0 {
		t.Fatalf("InFlight count=%d after shutdown", idemp.inFlightCountForTest())
	}
	res := svc.ReserveOrInspectReplay(createNS("k-conc-after"), DigestAction())
	if res.Outcome == ReservationOutcomeReserved {
		t.Fatal("no InFlight reservation may appear after abort sweep begins")
	}
	if idemp.inFlightCountForTest() != 0 {
		t.Fatal("post-shutdown reserve must create no InFlight")
	}
}

func TestScheduler_MaintenanceEnter_NoQualifyDoesNotCancelIdempotency(t *testing.T) {
	t.Parallel()

	svc, _, idemp, _, _, sched := newSchedulerFixture(t)
	et := createTargetForScheduler(t, svc, "t-noq", "99999999999999999999999999999999", "k-noq")
	ns := createNS("k-noq-other")
	if res := svc.ReserveOrInspectReplay(ns, DigestAction()); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	trig := baseMaintenanceTrigger(et, MaintenanceTriggerEnter, "audit-noq")
	if _, prob := sched.DeliverMaintenanceTrigger(context.Background(), trig); prob != nil {
		t.Fatalf("enter: %#v", prob)
	}
	if st, ok := idemp.stateForTest(ns, svc.now()); !ok || st != ReservationInFlight {
		t.Fatalf("no in-flight qualify means no idempotency cancel, state=%v ok=%v", st, ok)
	}
	svc.AbortReservation(ns)
}
