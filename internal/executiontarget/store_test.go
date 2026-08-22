package executiontarget

import (
	"sync"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

func participationRef(uid string) apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: "governance.sovrunn.io/v1alpha1",
		Kind:       "CloudProviderParticipation",
		Name:       "part",
		UID:        uid,
	}
}

func stackRef(uid string) apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: "infrastructure.sovrunn.io/v1alpha1",
		Kind:       "InfrastructureStack",
		Name:       "stack",
		UID:        uid,
	}
}

func TestNewCreateProposal_InitialState(t *testing.T) {
	t.Parallel()

	partUID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	stackUID := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	scopeUID := "cccccccccccccccccccccccccccccccc"
	et := model.NewCreateProposal(
		"target-1",
		"dddddddddddddddddddddddddddddddd",
		scopeUID,
		participationRef(partUID),
		stackRef(stackUID),
		7,
	)

	if et.CloudProviderScopeUID != scopeUID {
		t.Fatalf("derived scope=%q, want %q", et.CloudProviderScopeUID, scopeUID)
	}
	if et.Metadata.ScopeRef != nil {
		t.Fatalf("scopeRef must remain unset, got %#v", et.Metadata.ScopeRef)
	}
	if et.Status.Lifecycle != model.LifecycleActive {
		t.Fatalf("lifecycle=%q", et.Status.Lifecycle)
	}
	if et.Status.Qualification != model.QualificationUnqualified {
		t.Fatalf("qualification=%q", et.Status.Qualification)
	}
	if et.Status.MaintenanceEpoch != 0 {
		t.Fatalf("maintenanceEpoch=%d", et.Status.MaintenanceEpoch)
	}
	if et.Status.ObservedGeneration != 7 {
		t.Fatalf("observedGeneration=%d", et.Status.ObservedGeneration)
	}
	if et.Spec.TargetClass != model.TargetClassSyntheticIaaS {
		t.Fatalf("targetClass=%q", et.Spec.TargetClass)
	}
	if et.Status.FactSetRef != nil || et.Status.QualificationResultRef != nil {
		t.Fatalf("create proposal must have no record refs")
	}
	if et.APIVersion != model.APIVersionExecutionTarget || et.Kind != model.KindExecutionTarget {
		t.Fatalf("type meta=%s/%s", et.APIVersion, et.Kind)
	}
}

func TestNormalizedTargetFactSetAndResultFields(t *testing.T) {
	t.Parallel()

	fp := model.ComputeViabilityFingerprint("part-uid", true, "stack-uid", "Active")
	if fp == "" {
		t.Fatal("empty viability fingerprint")
	}
	fpInactive := model.ComputeViabilityFingerprint("part-uid", false, "stack-uid", "Active")
	if fp == fpInactive {
		t.Fatal("effective-active state must change fingerprint")
	}

	fs := model.NormalizedTargetFactSet{
		UID:                           "fs-1",
		TargetRef:                     apimeta.TypedRef{UID: "t-1", Kind: model.KindExecutionTarget, APIVersion: model.APIVersionExecutionTarget, Name: "t"},
		InfrastructureStackGeneration: 3,
		MaintenanceEpoch:              1,
		ViabilityFingerprint:          fp,
		ObserverID:                    model.ObserverID,
		ObserverRevision:              "rev-1",
		FactVersion:                   model.FactVersionV1,
		ObservedAt:                    "2026-08-21T00:00:00Z",
		ExpiresAt:                     "2026-08-21T01:00:00Z",
		Facts: model.FactSetTruths{
			Compute: model.FactComputeTruths{VM: model.FactSupported},
			Storage: model.FactStorageTruths{Block: model.FactSupported, Object: model.FactUnknown},
			Network: model.FactNetworkTruths{Private: model.FactUnsupported},
		},
	}
	if fs.ObserverID != model.ObserverID || fs.FactVersion != model.FactVersionV1 {
		t.Fatalf("provenance mismatch: %+v", fs)
	}

	result := model.TargetQualificationResult{
		UID:                           "qr-1",
		TargetRef:                     fs.TargetRef,
		FactSetRef:                    apimeta.TypedRef{UID: fs.UID, Kind: model.KindNormalizedTargetFactSet, APIVersion: model.APIVersionFactSet, Name: fs.UID},
		InfrastructureStackGeneration: 3,
		MaintenanceEpoch:              1,
		ProfileVersion:                model.ProfileVersion,
		Outcome:                       model.OutcomeQualified,
		ReasonCodes:                   []string{"ok"},
		EvaluatedAt:                   "2026-08-21T00:00:01Z",
	}
	if result.ProfileVersion != model.ProfileVersion {
		t.Fatalf("profileVersion=%q", result.ProfileVersion)
	}
	// Result must not carry a viabilityFingerprint field — verified by type shape.
	_ = result
}

func TestCurrentMaintenanceMarker(t *testing.T) {
	t.Parallel()

	m := model.CurrentMaintenanceMarker{
		TargetUID:        "dddddddddddddddddddddddddddddddd",
		MaintenanceEpoch: 2,
		Active:           true,
	}
	if !m.Active || m.MaintenanceEpoch != 2 {
		t.Fatalf("marker=%#v", m)
	}
}

func TestStore_StageDoesNotPublish(t *testing.T) {
	t.Parallel()

	s := NewStore()
	et := model.NewCreateProposal(
		"target-1",
		"dddddddddddddddddddddddddddddddd",
		"cccccccccccccccccccccccccccccccc",
		participationRef("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		stackRef("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"),
		1,
	)

	s.beginPublication()
	staged, prob := s.stageCreateExecutionTarget(et)
	if prob != nil {
		s.endPublication()
		t.Fatalf("stage: %#v", prob)
	}
	s.endPublication()

	if _, ok := s.GetExecutionTarget(et.Metadata.UID); ok {
		t.Fatal("staging alone must not publish the target")
	}

	s.beginPublication()
	s.abort(staged)
	s.endPublication()
	if _, ok := s.GetExecutionTarget(et.Metadata.UID); ok {
		t.Fatal("abort must leave target unpublished")
	}
}

func TestStore_PublishCreate_NameAndTupleIndexes(t *testing.T) {
	t.Parallel()

	s := NewStore()
	scope := "cccccccccccccccccccccccccccccccc"
	part := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	stack := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	et := model.NewCreateProposal("target-1", "dddddddddddddddddddddddddddddddd", scope, participationRef(part), stackRef(stack), 1)

	s.beginPublication()
	staged, prob := s.stageCreateExecutionTarget(et)
	if prob != nil {
		s.endPublication()
		t.Fatalf("stage: %#v", prob)
	}
	s.publish(staged)
	if !s.nameReserved(scope, "target-1") {
		s.endPublication()
		t.Fatal("name must be reserved after publish")
	}
	if !s.liveTupleOccupied(part, stack, model.TargetClassSyntheticIaaS) {
		s.endPublication()
		t.Fatal("live tuple must be indexed after publish")
	}
	s.endPublication()

	got, ok := s.GetExecutionTarget(et.Metadata.UID)
	if !ok {
		t.Fatal("published target missing")
	}
	if got.Status.Lifecycle != model.LifecycleActive || got.Status.Qualification != model.QualificationUnqualified {
		t.Fatalf("published state=%#v", got.Status)
	}
}

func TestStore_NameReservationRetainedAfterRetire(t *testing.T) {
	t.Parallel()

	s := NewStore()
	scope := "cccccccccccccccccccccccccccccccc"
	part := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	stack := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	et := model.NewCreateProposal("target-1", "dddddddddddddddddddddddddddddddd", scope, participationRef(part), stackRef(stack), 1)

	s.beginPublication()
	staged, prob := s.stageCreateExecutionTarget(et)
	if prob != nil {
		s.endPublication()
		t.Fatalf("stage create: %#v", prob)
	}
	s.publish(staged)

	cur, _ := s.lookupExecutionTarget(et.Metadata.UID)
	cur.Status.Lifecycle = model.LifecycleRetired
	cur.Status.Qualification = model.QualificationUnqualified
	upd, prob := s.stageUpdateExecutionTarget(cur)
	if prob != nil {
		s.endPublication()
		t.Fatalf("stage retire: %#v", prob)
	}
	s.publish(upd)

	if !s.nameReserved(scope, "target-1") {
		s.endPublication()
		t.Fatal("name reservation must be retained after retirement")
	}
	if s.liveTupleOccupied(part, stack, model.TargetClassSyntheticIaaS) {
		s.endPublication()
		t.Fatal("live tuple must be released after retirement")
	}

	again := model.NewCreateProposal("target-1", "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", scope, participationRef(part), stackRef(stack), 1)
	_, prob = s.stageCreateExecutionTarget(again)
	s.endPublication()
	if prob == nil || prob.Code != apiproblem.CodeAlreadyExists {
		t.Fatalf("retired name must block recreate: %#v", prob)
	}
}

func TestStore_LiveTupleAtMostOne(t *testing.T) {
	t.Parallel()

	s := NewStore()
	scope := "cccccccccccccccccccccccccccccccc"
	part := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	stack := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	first := model.NewCreateProposal("target-1", "dddddddddddddddddddddddddddddddd", scope, participationRef(part), stackRef(stack), 1)

	s.beginPublication()
	staged, prob := s.stageCreateExecutionTarget(first)
	if prob != nil {
		s.endPublication()
		t.Fatalf("stage: %#v", prob)
	}
	s.publish(staged)

	dup := model.NewCreateProposal("target-2", "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", scope, participationRef(part), stackRef(stack), 1)
	_, prob = s.stageCreateExecutionTarget(dup)
	s.endPublication()
	if prob == nil || prob.Code != apiproblem.CodeAlreadyExists {
		t.Fatalf("live duplicate tuple must be rejected: %#v", prob)
	}
}

func TestStore_LiveTupleFreedAllowsNewName(t *testing.T) {
	t.Parallel()

	s := NewStore()
	scope := "cccccccccccccccccccccccccccccccc"
	part := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	stack := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	first := model.NewCreateProposal("target-1", "dddddddddddddddddddddddddddddddd", scope, participationRef(part), stackRef(stack), 1)

	s.beginPublication()
	staged, prob := s.stageCreateExecutionTarget(first)
	if prob != nil {
		s.endPublication()
		t.Fatalf("stage: %#v", prob)
	}
	s.publish(staged)

	cur, _ := s.lookupExecutionTarget(first.Metadata.UID)
	cur.Status.Lifecycle = model.LifecycleRetired
	upd, prob := s.stageUpdateExecutionTarget(cur)
	if prob != nil {
		s.endPublication()
		t.Fatalf("retire: %#v", prob)
	}
	s.publish(upd)

	again := model.NewCreateProposal("target-2", "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", scope, participationRef(part), stackRef(stack), 1)
	staged2, prob := s.stageCreateExecutionTarget(again)
	if prob != nil {
		s.endPublication()
		t.Fatalf("tuple after retire should allow new name: %#v", prob)
	}
	s.publish(staged2)
	s.endPublication()

	if _, ok := s.GetExecutionTarget(again.Metadata.UID); !ok {
		t.Fatal("second target missing")
	}
}

func TestStore_FactSetResultAndMaintenanceMarkerStaging(t *testing.T) {
	t.Parallel()

	s := NewStore()
	fp := model.ComputeViabilityFingerprint("p", true, "s", "Active")
	fs := model.NormalizedTargetFactSet{
		UID:                  "fs-1",
		TargetRef:            apimeta.TypedRef{UID: "t-1"},
		ViabilityFingerprint: fp,
		ObserverID:           model.ObserverID,
		FactVersion:          model.FactVersionV1,
		Facts: model.FactSetTruths{
			Compute: model.FactComputeTruths{VM: model.FactSupported},
			Storage: model.FactStorageTruths{Block: model.FactSupported, Object: model.FactSupported},
			Network: model.FactNetworkTruths{Private: model.FactSupported},
		},
	}
	qr := model.TargetQualificationResult{
		UID:            "qr-1",
		TargetRef:      fs.TargetRef,
		FactSetRef:     apimeta.TypedRef{UID: fs.UID},
		ProfileVersion: model.ProfileVersion,
		Outcome:        model.OutcomeQualified,
		ReasonCodes:    []string{"ok"},
	}
	marker := model.CurrentMaintenanceMarker{TargetUID: "t-1", MaintenanceEpoch: 1, Active: true}

	s.beginPublication()
	fsStaged, prob := s.stagePutFactSet(fs)
	if prob != nil {
		s.endPublication()
		t.Fatalf("factset stage: %#v", prob)
	}
	qrStaged, prob := s.stagePutQualificationResult(qr)
	if prob != nil {
		s.endPublication()
		t.Fatalf("result stage: %#v", prob)
	}
	mStaged, prob := s.stageSetMaintenanceMarker(marker)
	if prob != nil {
		s.endPublication()
		t.Fatalf("marker stage: %#v", prob)
	}
	if _, ok := s.lookupFactSet(fs.UID); ok {
		s.endPublication()
		t.Fatal("fact set must not be visible before publish")
	}
	s.publish(fsStaged)
	s.publish(qrStaged)
	s.publish(mStaged)
	if _, ok := s.lookupFactSet(fs.UID); !ok {
		s.endPublication()
		t.Fatal("fact set missing after publish")
	}
	if _, ok := s.lookupQualificationResult(qr.UID); !ok {
		s.endPublication()
		t.Fatal("result missing after publish")
	}
	gotMarker, ok := s.lookupMaintenanceMarker("t-1")
	if !ok || !gotMarker.Active || gotMarker.MaintenanceEpoch != 1 {
		s.endPublication()
		t.Fatalf("marker=%#v ok=%v", gotMarker, ok)
	}

	s.publish(s.stageClearMaintenanceMarker("t-1"))
	if _, ok := s.lookupMaintenanceMarker("t-1"); ok {
		s.endPublication()
		t.Fatal("marker should be cleared")
	}
	s.publish(s.stageDeleteFactSet(fs.UID))
	s.publish(s.stageDeleteQualificationResult(qr.UID))
	s.endPublication()

	if _, ok := s.GetFactSet(fs.UID); ok {
		t.Fatal("fact set should be deleted")
	}
	if _, ok := s.GetQualificationResult(qr.UID); ok {
		t.Fatal("result should be deleted")
	}
}

func TestStore_ListAscendingUID(t *testing.T) {
	t.Parallel()

	s := NewStore()
	scope := "cccccccccccccccccccccccccccccccc"
	a := model.NewCreateProposal("a", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", scope, participationRef("11111111111111111111111111111111"), stackRef("22222222222222222222222222222222"), 1)
	b := model.NewCreateProposal("b", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", scope, participationRef("33333333333333333333333333333333"), stackRef("44444444444444444444444444444444"), 1)

	s.beginPublication()
	sa, pa := s.stageCreateExecutionTarget(a)
	sb, pb := s.stageCreateExecutionTarget(b)
	if pa != nil || pb != nil {
		s.endPublication()
		t.Fatalf("stage errors: %#v %#v", pa, pb)
	}
	s.publish(sa)
	s.publish(sb)
	s.endPublication()

	list := s.ListExecutionTargets()
	if len(list) != 2 {
		t.Fatalf("len=%d", len(list))
	}
	if list[0].Metadata.UID >= list[1].Metadata.UID {
		t.Fatalf("list not ascending by uid: %q then %q", list[0].Metadata.UID, list[1].Metadata.UID)
	}
}

func TestStore_ConcurrentCreateNameRace(t *testing.T) {
	s := NewStore()
	scope := "cccccccccccccccccccccccccccccccc"
	const n = 32
	var wg sync.WaitGroup
	wg.Add(n)
	errs := make([]*apiproblem.Problem, n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			uid, _ := apimeta.GenerateUID()
			part, _ := apimeta.GenerateUID()
			stack, _ := apimeta.GenerateUID()
			et := model.NewCreateProposal("same-name", uid, scope, participationRef(part), stackRef(stack), 1)
			s.beginPublication()
			staged, prob := s.stageCreateExecutionTarget(et)
			if prob != nil {
				s.endPublication()
				errs[i] = prob
				return
			}
			s.publish(staged)
			s.endPublication()
		}()
	}
	wg.Wait()
	success := 0
	for _, e := range errs {
		if e == nil {
			success++
			continue
		}
		if e.Code != apiproblem.CodeAlreadyExists {
			t.Fatalf("unexpected %#v", e)
		}
	}
	if success != 1 {
		t.Fatalf("want exactly one winner, got %d", success)
	}
}

func TestStore_EmptyListNonNil(t *testing.T) {
	t.Parallel()
	s := NewStore()
	list := s.ListExecutionTargets()
	if list == nil {
		t.Fatal("empty list must be non-nil")
	}
	if len(list) != 0 {
		t.Fatalf("len=%d", len(list))
	}
}
