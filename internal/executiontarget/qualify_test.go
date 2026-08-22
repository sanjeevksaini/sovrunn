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

func TestEvaluateObservation_ClosedProfile(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)
	base := Observation{
		ObserverID:        model.ObserverID,
		ObserverRevision:  "1",
		FactSchemaVersion: model.FactSchemaVersionV1,
		ObservedAt:        now,
		ExpiresAt:         now.Add(FactFreshnessWindow),
	}

	t.Run("qualified", func(t *testing.T) {
		obs := base
		obs.Facts = factsFromTruths(allSupportedTruths())
		eval, err := EvaluateObservation(obs, now)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if eval.Outcome != model.OutcomeQualified {
			t.Fatalf("outcome=%s", eval.Outcome)
		}
	})

	t.Run("rejected-any-unsupported", func(t *testing.T) {
		obs := base
		obs.Facts = factsFromTruths(unsupportedVMTruths())
		eval, err := EvaluateObservation(obs, now)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if eval.Outcome != model.OutcomeRejected {
			t.Fatalf("outcome=%s", eval.Outcome)
		}
	})

	t.Run("indeterminate-any-unknown", func(t *testing.T) {
		obs := base
		obs.Facts = factsFromTruths(unknownBlockTruths())
		eval, err := EvaluateObservation(obs, now)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if eval.Outcome != model.OutcomeIndeterminate {
			t.Fatalf("outcome=%s", eval.Outcome)
		}
	})

	t.Run("unsupported-outranks-unknown", func(t *testing.T) {
		truths := unsupportedVMTruths()
		truths.Storage.Block = model.FactUnknown
		obs := base
		obs.Facts = factsFromTruths(truths)
		eval, err := EvaluateObservation(obs, now)
		if err != nil {
			t.Fatalf("eval: %v", err)
		}
		if eval.Outcome != model.OutcomeRejected {
			t.Fatalf("outcome=%s", eval.Outcome)
		}
	})
}

func TestEvaluateObservation_ObserverFaults(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 21, 12, 0, 0, 0, time.UTC)
	base := Observation{
		ObserverID:        model.ObserverID,
		ObserverRevision:  "1",
		FactSchemaVersion: model.FactSchemaVersionV1,
		ObservedAt:        now,
		ExpiresAt:         now.Add(FactFreshnessWindow),
		Facts:             factsFromTruths(allSupportedTruths()),
	}

	cases := []struct {
		name string
		mut  func(*Observation)
	}{
		{"duplicate", func(o *Observation) {
			o.Facts = append(o.Facts, ProposedFact{Name: FactNameComputeVM, Truth: model.FactSupported})
		}},
		{"missing", func(o *Observation) {
			o.Facts = o.Facts[:3]
		}},
		{"malformed-provenance", func(o *Observation) {
			o.FactSchemaVersion = "bad"
		}},
		{"malformed-truth", func(o *Observation) {
			o.Facts[0].Truth = model.FactTruth("Maybe")
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			obs := base
			obs.Facts = append([]ProposedFact(nil), base.Facts...)
			tc.mut(&obs)
			_, err := EvaluateObservation(obs, now)
			if err == nil {
				t.Fatal("want observer fault")
			}
			if _, ok := err.(ObserverFault); !ok {
				t.Fatalf("want ObserverFault, got %T %v", err, err)
			}
		})
	}
}

func qualifyNS(targetUID, key string) IdempotencyNamespace {
	return IdempotencyNamespace{
		PrincipalUID:     "principal-1",
		RoutePattern:     "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify",
		CloudProviderUID: "cccccccccccccccccccccccccccccccc",
		ActionTargetUID:  targetUID,
		Key:              key,
	}
}

func activeBacking() BackingViability {
	return BackingViability{
		ParticipationUID:             "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		ParticipationEffectiveActive: true,
		ParticipationScopeUID:        "cccccccccccccccccccccccccccccccc",
		StackUID:                     "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		StackPhase:                   "Active",
		StackScopeUID:                "cccccccccccccccccccccccccccccccc",
		StackGeneration:              3,
	}
}

func newQualifyLifecycle(t *testing.T, src FixtureSource) (*ExecutionTargetLifecycleService, *Store, *IdempotencyTable, *memoryAuditAppender, *MapFixtureSource) {
	t.Helper()
	store := NewStore()
	idemp := NewIdempotencyTable()
	audit := &memoryAuditAppender{}
	now := fixedNow(time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC))
	mapSrc, _ := src.(*MapFixtureSource)
	if mapSrc == nil {
		mapSrc = NewMapFixtureSource()
		if src == nil {
			src = mapSrc
		}
	}
	obs := NewSyntheticObserver(src, now)
	svc := NewExecutionTargetLifecycleService(LifecycleConfig{
		Store:       store,
		Idempotency: idemp,
		Audit:       audit,
		Observer:    obs,
		Now:         now,
	})
	return svc, store, idemp, audit, mapSrc
}

func createActiveTarget(t *testing.T, svc *ExecutionTargetLifecycleService) model.ExecutionTarget {
	t.Helper()
	req := baseCreateReq("k-create-for-qualify")
	req.Digest = mustCreateDigest(t, `{"metadata":{"name":"target-1"}}`)
	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve create: %#v", res)
	}
	et, prob := svc.CommitCreate(context.Background(), req)
	if prob != nil {
		t.Fatalf("create: %#v", prob)
	}
	return et
}

func baseQualifyReq(targetUID string, key string) QualifyCommitRequest {
	return QualifyCommitRequest{
		TargetUID:  targetUID,
		Backing:    activeBacking(),
		FactSetUID: "fs-" + key,
		ResultUID:  "qr-" + key,
		Actor: apimeta.TypedRef{
			APIVersion: "identity.sovrunn.io/v1alpha1",
			Kind:       "Principal",
			Name:       "principal",
			UID:        "principal-1",
		},
		RequestID:   "req-qualify-" + key,
		AuditUID:    "audit-qualify-" + key,
		Idempotency: qualifyNS(targetUID, key),
		Digest:      DigestAction(),
		Completion:  CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{"qualified":true}`)},
	}
}

func TestQualify_HappyPath_Qualified(t *testing.T) {
	t.Parallel()

	svc, store, idemp, audit, fixtures := newQualifyLifecycle(t, nil)
	et := createActiveTarget(t, svc)
	fixtures.Set(et.Metadata.UID, presentFixture("fx-1", allSupportedTruths()))

	req := baseQualifyReq(et.Metadata.UID, "k-ok")
	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	beforeRV := et.Metadata.ResourceVersion

	got, prob := svc.Qualify(context.Background(), req)
	if prob != nil {
		t.Fatalf("qualify: %#v", prob)
	}
	if got.Status.Qualification != model.QualificationQualified {
		t.Fatalf("qualification=%s", got.Status.Qualification)
	}
	if got.Metadata.ResourceVersion == beforeRV {
		t.Fatal("resourceVersion must advance")
	}
	if got.Status.ObservedGeneration != 3 {
		t.Fatalf("observedGeneration=%d", got.Status.ObservedGeneration)
	}
	if got.Status.FactSetRef == nil || got.Status.QualificationResultRef == nil {
		t.Fatal("both record refs must be set")
	}
	fs, ok := store.GetFactSet(req.FactSetUID)
	if !ok {
		t.Fatal("fact set missing")
	}
	if fs.FactVersion != model.FactVersionV1 || fs.ObserverRevision != "fx-1" {
		t.Fatalf("fact set provenance: %#v", fs)
	}
	if fs.ViabilityFingerprint != req.Backing.Fingerprint() {
		t.Fatalf("viability fingerprint mismatch")
	}
	qr, ok := store.GetQualificationResult(req.ResultUID)
	if !ok || qr.Outcome != model.OutcomeQualified {
		t.Fatalf("result: ok=%v %#v", ok, qr)
	}
	st, ok := idemp.stateForTest(req.Idempotency, svc.now())
	if !ok || st != ReservationCompleted {
		t.Fatalf("idempotency state=%v ok=%v", st, ok)
	}
	if audit.Len() != 2 { // create + qualify
		t.Fatalf("audit len=%d", audit.Len())
	}
	if audit.Last().Record.Action != AuditActionQualify {
		t.Fatalf("last audit action=%s", audit.Last().Record.Action)
	}
	if svc.Observer().ExternalCallCount() != 0 {
		t.Fatalf("externalCallCount=%d", svc.Observer().ExternalCallCount())
	}
}

func TestQualify_MissingFixtureAndLogicalTimeout(t *testing.T) {
	t.Parallel()

	t.Run("missing", func(t *testing.T) {
		svc, store, _, _, _ := newQualifyLifecycle(t, nil)
		et := createActiveTarget(t, svc)
		req := baseQualifyReq(et.Metadata.UID, "k-miss")
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}
		got, prob := svc.Qualify(context.Background(), req)
		if prob != nil {
			t.Fatalf("qualify: %#v", prob)
		}
		if got.Status.Qualification != model.QualificationIndeterminate {
			t.Fatalf("qualification=%s", got.Status.Qualification)
		}
		fs, _ := store.GetFactSet(req.FactSetUID)
		if fs.Facts.Compute.VM != model.FactUnknown {
			t.Fatalf("want Unknown facts: %#v", fs.Facts)
		}
	})

	t.Run("logical-timeout", func(t *testing.T) {
		svc, store, _, _, fixtures := newQualifyLifecycle(t, nil)
		et := createActiveTarget(t, svc)
		fixtures.Set(et.Metadata.UID, logicalTimeoutFixture("to-1"))
		req := baseQualifyReq(et.Metadata.UID, "k-to")
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}
		got, prob := svc.Qualify(context.Background(), req)
		if prob != nil {
			t.Fatalf("qualify: %#v", prob)
		}
		if got.Status.Qualification != model.QualificationIndeterminate {
			t.Fatalf("qualification=%s", got.Status.Qualification)
		}
		fs, _ := store.GetFactSet(req.FactSetUID)
		if fs.ObserverRevision != "to-1" {
			t.Fatalf("observerRevision=%q", fs.ObserverRevision)
		}
	})
}

func TestQualify_ObserverFaultPublishesNothing(t *testing.T) {
	t.Parallel()

	svc, store, idemp, audit, fixtures := newQualifyLifecycle(t, nil)
	et := createActiveTarget(t, svc)
	beforeAudit := audit.Len()
	beforeRV := et.Metadata.ResourceVersion
	fixtures.Set(et.Metadata.UID, duplicateFaultFixture("bad"))

	req := baseQualifyReq(et.Metadata.UID, "k-fault")
	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	_, prob := svc.Qualify(context.Background(), req)
	if prob == nil || prob.Code != apiproblem.CodeInternalError {
		t.Fatalf("want INTERNAL_ERROR, got %#v", prob)
	}
	after, _ := store.GetExecutionTarget(et.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV {
		t.Fatal("fault must not mutate ETag")
	}
	if after.Status.FactSetRef != nil {
		t.Fatal("fault must not publish fact set ref")
	}
	if _, ok := store.GetFactSet(req.FactSetUID); ok {
		t.Fatal("fault must not publish fact set")
	}
	if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
		t.Fatal("fault must remove InFlight idempotency reservation")
	}
	if audit.Len() != beforeAudit {
		t.Fatal("fault must not append qualification AuditEvent")
	}
}

func TestQualify_BackingViabilityRecheck(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		mut    func(*BackingViability)
		code   apiproblem.ErrorCode
		viol   apiproblem.ViolationCode
		status int
	}{
		{
			name:   "participation-inactive",
			mut:    func(b *BackingViability) { b.ParticipationEffectiveActive = false },
			code:   apiproblem.CodeConflict,
			viol:   ViolationExecutionTargetParticipationUnavailable,
			status: http.StatusConflict,
		},
		{
			name:   "stack-inactive",
			mut:    func(b *BackingViability) { b.StackPhase = "Failed" },
			code:   apiproblem.CodeConflict,
			viol:   ViolationExecutionTargetStackUnavailable,
			status: http.StatusConflict,
		},
		{
			name:   "scope-mismatch",
			mut:    func(b *BackingViability) { b.StackScopeUID = "dddddddddddddddddddddddddddddddd" },
			code:   apiproblem.CodeValidationFailed,
			viol:   ViolationExecutionTargetScopeMismatch,
			status: http.StatusUnprocessableEntity,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc, store, idemp, audit, fixtures := newQualifyLifecycle(t, nil)
			et := createActiveTarget(t, svc)
			fixtures.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
			beforeAudit := audit.Len()
			beforeRV := et.Metadata.ResourceVersion

			req := baseQualifyReq(et.Metadata.UID, "k-"+tc.name)
			tc.mut(&req.Backing)

			if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
				t.Fatalf("reserve: %#v", res)
			}

			// Attach a waiter before the failing qualify to prove post-unlock wake.
			done, ok := svc.AttachWaiter(req.Idempotency)
			if !ok {
				t.Fatal("attach waiter")
			}
			woke := make(chan struct{})
			go func() {
				_ = svc.AwaitReservation(context.Background(), done)
				close(woke)
			}()

			_, prob := svc.Qualify(context.Background(), req)
			if prob == nil || prob.Code != tc.code {
				t.Fatalf("want %s, got %#v", tc.code, prob)
			}
			if len(prob.Violations) == 0 || prob.Violations[0].Code != tc.viol {
				t.Fatalf("want violation %s, got %#v", tc.viol, prob.Violations)
			}
			if apiproblem.StatusForCode(prob.Code) != tc.status {
				t.Fatalf("status=%d want %d", apiproblem.StatusForCode(prob.Code), tc.status)
			}

			select {
			case <-woke:
			case <-time.After(2 * time.Second):
				t.Fatal("waiter must wake after unlock")
			}

			after, _ := store.GetExecutionTarget(et.Metadata.UID)
			if after.Metadata.ResourceVersion != beforeRV || after.Status.FactSetRef != nil {
				t.Fatal("viability reject must not observe/publish")
			}
			if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
				t.Fatal("must detach InFlight reservation")
			}
			if audit.Len() != beforeAudit {
				t.Fatal("must not append AuditEvent")
			}
			if svc.Observer().ExternalCallCount() != 0 {
				t.Fatalf("externalCallCount=%d", svc.Observer().ExternalCallCount())
			}

			// Waiter re-entry: a new reservation must succeed.
			if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
				t.Fatalf("re-entry reserve: %#v", res)
			}
			_ = svc.AbortReservation(req.Idempotency)
		})
	}
}

func TestQualify_StaleFenceOrderedPredicate(t *testing.T) {
	t.Parallel()

	t.Run("active-maintenance-first", func(t *testing.T) {
		svc, store, idemp, audit, fixtures := newQualifyLifecycle(t, nil)
		et := createActiveTarget(t, svc)
		fixtures.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
		beforeAudit := audit.Len()
		beforeRV := et.Metadata.ResourceVersion

		// Plant an active Maintenance marker (TASK-F16-06 owns enter; unit seed here).
		store.beginPublication()
		staged, _ := store.stageSetMaintenanceMarker(model.CurrentMaintenanceMarker{
			TargetUID:        et.Metadata.UID,
			MaintenanceEpoch: et.Status.MaintenanceEpoch,
			Active:           true,
		})
		store.publish(staged)
		store.endPublication()

		req := baseQualifyReq(et.Metadata.UID, "k-maint")
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}
		_, prob := svc.Qualify(context.Background(), req)
		if prob == nil || prob.Code != apiproblem.CodeConflict {
			t.Fatalf("want CONFLICT, got %#v", prob)
		}
		if len(prob.Violations) == 0 || prob.Violations[0].Code != ViolationTargetMaintenance {
			t.Fatalf("want VS0_TARGET_MAINTENANCE, got %#v", prob.Violations)
		}
		after, _ := store.GetExecutionTarget(et.Metadata.UID)
		if after.Metadata.ResourceVersion != beforeRV {
			t.Fatal("maintenance abort must not mutate target")
		}
		if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
			t.Fatal("must remove InFlight")
		}
		if audit.Len() != beforeAudit {
			t.Fatal("must not qualify-audit")
		}
	})

	t.Run("epoch-stale", func(t *testing.T) {
		mapSrc := NewMapFixtureSource()
		started := make(chan struct{})
		release := make(chan struct{})
		blocking := &blockingFixtureSource{inner: mapSrc, started: started, release: release}
		svc, store, idemp, audit, _ := newQualifyLifecycle(t, blocking)
		et := createActiveTarget(t, svc)
		mapSrc.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
		beforeAudit := audit.Len()
		beforeRV := et.Metadata.ResourceVersion

		req := baseQualifyReq(et.Metadata.UID, "k-epoch")
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}

		errCh := make(chan *apiproblem.Problem, 1)
		go func() {
			_, prob := svc.Qualify(context.Background(), req)
			errCh <- prob
		}()
		<-started

		// Bump maintenance epoch after Qualifying fences were captured.
		store.beginPublication()
		cur, _ := store.lookupExecutionTarget(et.Metadata.UID)
		cur.Status.MaintenanceEpoch = 9
		staged, _ := store.stageUpdateExecutionTarget(cur)
		store.publish(staged)
		store.endPublication()
		close(release)

		prob := <-errCh
		if prob == nil || prob.Code != apiproblem.CodeStaleResourceVersion {
			t.Fatalf("want STALE_RESOURCE_VERSION, got %#v", prob)
		}
		if len(prob.Violations) == 0 || prob.Violations[0].Code != ViolationTargetEpochStale {
			t.Fatalf("want VS0_TARGET_EPOCH_STALE, got %#v", prob.Violations)
		}
		after, _ := store.GetExecutionTarget(et.Metadata.UID)
		// Epoch bump itself advanced RV; qualification must not publish refs.
		if after.Status.FactSetRef != nil {
			t.Fatal("epoch-stale must not publish qualification refs")
		}
		if after.Metadata.ResourceVersion == beforeRV {
			t.Fatal("epoch bump should have advanced RV independently")
		}
		if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
			t.Fatal("must remove InFlight")
		}
		if audit.Len() != beforeAudit {
			t.Fatal("must not qualify-audit")
		}
	})

	t.Run("generation-stale", func(t *testing.T) {
		svc, _, idemp, audit, fixtures := newQualifyLifecycle(t, nil)
		et := createActiveTarget(t, svc)
		fixtures.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
		beforeAudit := audit.Len()

		req := baseQualifyReq(et.Metadata.UID, "k-gen")
		req.RefreshBacking = func() BackingViability {
			b := req.Backing
			b.StackGeneration = 99
			return b
		}
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}
		_, prob := svc.Qualify(context.Background(), req)
		if prob == nil || prob.Code != apiproblem.CodeStaleResourceVersion {
			t.Fatalf("want STALE_RESOURCE_VERSION, got %#v", prob)
		}
		if len(prob.Violations) == 0 || prob.Violations[0].Code != ViolationTargetEpochStale {
			t.Fatalf("want VS0_TARGET_EPOCH_STALE, got %#v", prob.Violations)
		}
		if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
			t.Fatal("must remove InFlight")
		}
		if audit.Len() != beforeAudit {
			t.Fatal("must not qualify-audit")
		}
	})

	t.Run("viability-stale", func(t *testing.T) {
		svc, _, idemp, audit, fixtures := newQualifyLifecycle(t, nil)
		et := createActiveTarget(t, svc)
		fixtures.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
		beforeAudit := audit.Len()

		req := baseQualifyReq(et.Metadata.UID, "k-via")
		req.RefreshBacking = func() BackingViability {
			b := req.Backing
			b.ParticipationEffectiveActive = false
			return b
		}
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}
		_, prob := svc.Qualify(context.Background(), req)
		if prob == nil || prob.Code != apiproblem.CodeStaleResourceVersion {
			t.Fatalf("want STALE_RESOURCE_VERSION, got %#v", prob)
		}
		if len(prob.Violations) == 0 || prob.Violations[0].Code != ViolationExecutionTargetViabilityStale {
			t.Fatalf("want VS0_EXECUTION_TARGET_VIABILITY_STALE, got %#v", prob.Violations)
		}
		if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
			t.Fatal("must remove InFlight")
		}
		if audit.Len() != beforeAudit {
			t.Fatal("must not qualify-audit")
		}
	})
}

func TestQualify_AuditAppendFailure(t *testing.T) {
	t.Parallel()

	svc, store, idemp, audit, fixtures := newQualifyLifecycle(t, nil)
	et := createActiveTarget(t, svc)
	fixtures.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
	beforeRV := et.Metadata.ResourceVersion
	beforeAudit := audit.Len()
	audit.SetFail(errors.New("append failed"))

	req := baseQualifyReq(et.Metadata.UID, "k-audit-fail")
	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	done, ok := svc.AttachWaiter(req.Idempotency)
	if !ok {
		t.Fatal("attach")
	}
	woke := make(chan struct{})
	go func() {
		_ = svc.AwaitReservation(context.Background(), done)
		close(woke)
	}()

	_, prob := svc.Qualify(context.Background(), req)
	if prob == nil || prob.Code != apiproblem.CodeInternalError {
		t.Fatalf("want INTERNAL_ERROR, got %#v", prob)
	}
	select {
	case <-woke:
	case <-time.After(2 * time.Second):
		t.Fatal("waiter must wake after unlock")
	}
	after, _ := store.GetExecutionTarget(et.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV || after.Status.FactSetRef != nil {
		t.Fatal("failed append must leave target/ETag unmutated")
	}
	if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
		t.Fatal("must remove InFlight, not leave unchanged")
	}
	if audit.Len() != beforeAudit {
		t.Fatal("failed append must not retain AuditEvent")
	}
}

func TestQualify_CancellationAndShutdown(t *testing.T) {
	t.Parallel()

	t.Run("cancelled-before", func(t *testing.T) {
		svc, _, idemp, audit, fixtures := newQualifyLifecycle(t, nil)
		et := createActiveTarget(t, svc)
		fixtures.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
		beforeAudit := audit.Len()
		req := baseQualifyReq(et.Metadata.UID, "k-cancel")
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, prob := svc.Qualify(ctx, req)
		if prob == nil || prob.Code != apiproblem.CodeInternalError {
			t.Fatalf("want INTERNAL_ERROR, got %#v", prob)
		}
		if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
			t.Fatal("must remove InFlight")
		}
		if audit.Len() != beforeAudit {
			t.Fatal("must not audit")
		}
	})

	t.Run("shutdown-mid-flight", func(t *testing.T) {
		mapSrc := NewMapFixtureSource()
		started := make(chan struct{})
		release := make(chan struct{})
		blocking := &blockingFixtureSource{inner: mapSrc, started: started, release: release}
		svc, store, idemp, audit, _ := newQualifyLifecycle(t, blocking)
		et := createActiveTarget(t, svc)
		mapSrc.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
		beforeAudit := audit.Len()
		beforeRV := et.Metadata.ResourceVersion

		req := baseQualifyReq(et.Metadata.UID, "k-stop")
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}

		errCh := make(chan *apiproblem.Problem, 1)
		go func() {
			_, prob := svc.Qualify(context.Background(), req)
			errCh <- prob
		}()
		<-started
		svc.MarkQualificationStopped()
		close(release)
		prob := <-errCh
		if prob == nil || prob.Code != apiproblem.CodeInternalError {
			t.Fatalf("want INTERNAL_ERROR, got %#v", prob)
		}
		after, _ := store.GetExecutionTarget(et.Metadata.UID)
		if after.Metadata.ResourceVersion != beforeRV {
			t.Fatal("shutdown must not publish")
		}
		if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
			t.Fatal("must remove InFlight")
		}
		if audit.Len() != beforeAudit {
			t.Fatal("must not qualify-audit")
		}
	})
}

func TestQualify_RetireWins(t *testing.T) {
	t.Parallel()

	mapSrc := NewMapFixtureSource()
	started := make(chan struct{})
	release := make(chan struct{})
	blocking := &blockingFixtureSource{inner: mapSrc, started: started, release: release}
	svc, store, idemp, audit, _ := newQualifyLifecycle(t, blocking)
	et := createActiveTarget(t, svc)
	mapSrc.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))

	qreq := baseQualifyReq(et.Metadata.UID, "k-retire-wins-q")
	if res := svc.ReserveOrInspectReplay(qreq.Idempotency, qreq.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve qualify: %#v", res)
	}
	qDone, ok := svc.AttachWaiter(qreq.Idempotency)
	if !ok {
		t.Fatal("attach qualify waiter")
	}
	waiterWoke := make(chan struct{})
	go func() {
		_ = svc.AwaitReservation(context.Background(), qDone)
		close(waiterWoke)
	}()

	qualErr := make(chan *apiproblem.Problem, 1)
	go func() {
		_, prob := svc.Qualify(context.Background(), qreq)
		qualErr <- prob
	}()
	<-started

	retire := RetireCommitRequest{
		TargetUID:               et.Metadata.UID,
		ExpectedResourceVersion: et.Metadata.ResourceVersion,
		Actor:                   qreq.Actor,
		RequestID:               "req-retire-wins",
		AuditUID:                "audit-retire-wins",
		Idempotency:             retireNS(et.Metadata.UID, "k-retire-wins"),
		Digest:                  DigestAction(),
		Completion:              CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{}`)},
	}
	if res := svc.ReserveOrInspectReplay(retire.Idempotency, retire.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve retire: %#v", res)
	}
	retired, prob := svc.CommitRetire(context.Background(), retire)
	if prob != nil {
		t.Fatalf("retire: %#v", prob)
	}
	if retired.Status.Lifecycle != model.LifecycleRetired {
		t.Fatalf("lifecycle=%s", retired.Status.Lifecycle)
	}
	close(release)

	qprob := <-qualErr
	if qprob == nil || qprob.Code != apiproblem.CodeConflict {
		t.Fatalf("qualifying caller want CONFLICT, got %#v", qprob)
	}
	if len(qprob.Violations) == 0 || qprob.Violations[0].Code != ViolationTargetRetired {
		t.Fatalf("want VS0_TARGET_RETIRED, got %#v", qprob.Violations)
	}

	select {
	case <-waiterWoke:
	case <-time.After(2 * time.Second):
		t.Fatal("qualify waiter must wake after unlock")
	}

	if _, ok := idemp.stateForTest(qreq.Idempotency, svc.now()); ok {
		t.Fatal("qualify InFlight must be removed")
	}
	st, ok := idemp.stateForTest(retire.Idempotency, svc.now())
	if !ok || st != ReservationCompleted {
		t.Fatalf("retire must complete: ok=%v st=%v", ok, st)
	}
	if audit.Last().Record.Action != AuditActionRetire {
		t.Fatalf("retirement must be audited, last=%s", audit.Last().Record.Action)
	}
	if _, ok := store.GetFactSet(qreq.FactSetUID); ok {
		t.Fatal("no qualification result/factset on Retire-wins")
	}
}

func TestQualify_FailedRetireAuditLeavesQualifyingIntact(t *testing.T) {
	t.Parallel()

	mapSrc := NewMapFixtureSource()
	started := make(chan struct{})
	release := make(chan struct{})
	blocking := &blockingFixtureSource{inner: mapSrc, started: started, release: release}
	svc, store, idemp, audit, _ := newQualifyLifecycle(t, blocking)
	et := createActiveTarget(t, svc)
	mapSrc.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
	beforeRV := et.Metadata.ResourceVersion

	qreq := baseQualifyReq(et.Metadata.UID, "k-retire-fail-q")
	if res := svc.ReserveOrInspectReplay(qreq.Idempotency, qreq.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve qualify: %#v", res)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	qualErr := make(chan *apiproblem.Problem, 1)
	go func() {
		defer wg.Done()
		_, prob := svc.Qualify(context.Background(), qreq)
		qualErr <- prob
	}()
	<-started

	audit.SetFail(errors.New("retire append failed"))
	retire := RetireCommitRequest{
		TargetUID:               et.Metadata.UID,
		ExpectedResourceVersion: et.Metadata.ResourceVersion,
		Actor:                   qreq.Actor,
		RequestID:               "req-retire-fail",
		AuditUID:                "audit-retire-fail",
		Idempotency:             retireNS(et.Metadata.UID, "k-retire-fail"),
		Digest:                  DigestAction(),
		Completion:              CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{}`)},
	}
	if res := svc.ReserveOrInspectReplay(retire.Idempotency, retire.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve retire: %#v", res)
	}
	_, rprob := svc.CommitRetire(context.Background(), retire)
	if rprob == nil || rprob.Code != apiproblem.CodeInternalError {
		t.Fatalf("retire append fail: %#v", rprob)
	}

	// Qualifying must still be in flight.
	svc.lock()
	_, still := svc.qualifying[et.Metadata.UID]
	svc.unlock()
	if !still {
		t.Fatal("failed retire must leave Qualifying reservation intact")
	}
	st, ok := idemp.stateForTest(qreq.Idempotency, svc.now())
	if !ok || st != ReservationInFlight {
		t.Fatalf("qualify InFlight must remain: ok=%v st=%v", ok, st)
	}
	after, _ := store.GetExecutionTarget(et.Metadata.UID)
	if after.Status.Lifecycle != model.LifecycleActive || after.Metadata.ResourceVersion != beforeRV {
		t.Fatalf("target must remain Active/unmutated: %#v", after)
	}

	audit.SetFail(nil)
	close(release)
	wg.Wait()
	qprob := <-qualErr
	if qprob != nil {
		t.Fatalf("qualify should succeed after failed retire: %#v", qprob)
	}
}

func TestQualify_PanicCleansReservations(t *testing.T) {
	t.Parallel()

	svc, store, idemp, _, fixtures := newQualifyLifecycle(t, nil)
	et := createActiveTarget(t, svc)
	fixtures.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
	beforeRV := et.Metadata.ResourceVersion

	req := baseQualifyReq(et.Metadata.UID, "k-panic")
	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	svc.audit = &panicAuditAppender{}

	panicked := false
	func() {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		_, _ = svc.Qualify(context.Background(), req)
	}()
	if !panicked {
		t.Fatal("expected panic")
	}
	after, _ := store.GetExecutionTarget(et.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV {
		t.Fatal("panic must not publish")
	}
	if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
		t.Fatal("panic must remove InFlight")
	}
	svc.lock()
	_, still := svc.qualifying[et.Metadata.UID]
	svc.unlock()
	if still {
		t.Fatal("panic must clear Qualifying reservation")
	}
}

func TestQualify_UnderFinalBackingLease_DeniedAbortsWithoutAudit(t *testing.T) {
	t.Parallel()
	svc, store, idemp, audit, fixtures := newQualifyLifecycle(t, nil)
	et := createActiveTarget(t, svc)
	fixtures.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
	beforeAudit := audit.Len()
	beforeRV := et.Metadata.ResourceVersion

	req := baseQualifyReq(et.Metadata.UID, "k-lease-deny")
	req.UnderFinalBackingLease = func(commit FinalQualifyCommitFunc) (model.ExecutionTarget, *apiproblem.Problem, chan struct{}, bool) {
		_ = commit
		return model.ExecutionTarget{}, nil, nil, false
	}
	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	_, prob := svc.Qualify(context.Background(), req)
	if prob == nil || prob.Code != apiproblem.CodeStaleResourceVersion {
		t.Fatalf("want viability-stale, got %#v", prob)
	}
	if len(prob.Violations) == 0 || prob.Violations[0].Code != ViolationExecutionTargetViabilityStale {
		t.Fatalf("want VS0_EXECUTION_TARGET_VIABILITY_STALE, got %#v", prob.Violations)
	}
	after, _ := store.GetExecutionTarget(et.Metadata.UID)
	if after.Metadata.ResourceVersion != beforeRV || after.Status.FactSetRef != nil {
		t.Fatal("denied final lease must not publish")
	}
	if _, ok := idemp.stateForTest(req.Idempotency, svc.now()); ok {
		t.Fatal("must remove InFlight")
	}
	if audit.Len() != beforeAudit {
		t.Fatal("must not qualify-audit")
	}
}

func TestQualify_ParticipationGenerationNotFence(t *testing.T) {
	t.Parallel()
	svc, store, _, audit, fixtures := newQualifyLifecycle(t, nil)
	et := createActiveTarget(t, svc)
	fixtures.Set(et.Metadata.UID, presentFixture("fx", allSupportedTruths()))
	beforeAudit := audit.Len()

	req := baseQualifyReq(et.Metadata.UID, "k-part-gen")
	// Synthesize a final lease view that only "changes" participation generation
	// indirectly by keeping the same viability fingerprint (generation is not
	// part of BackingViability / fences).
	req.UnderFinalBackingLease = func(commit FinalQualifyCommitFunc) (model.ExecutionTarget, *apiproblem.Problem, chan struct{}, bool) {
		out, prob, wake := commit(req.Backing)
		return out, prob, wake, true
	}
	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	got, prob := svc.Qualify(context.Background(), req)
	if prob != nil {
		t.Fatalf("participation generation must not fence: %#v", prob)
	}
	if got.Status.FactSetRef == nil {
		t.Fatal("qualification must publish")
	}
	if audit.Len() != beforeAudit+1 {
		t.Fatalf("audit=%d", audit.Len()-beforeAudit)
	}
	after, _ := store.GetExecutionTarget(et.Metadata.UID)
	if after.Status.Qualification != model.QualificationQualified {
		t.Fatalf("qualification=%s", after.Status.Qualification)
	}
}
