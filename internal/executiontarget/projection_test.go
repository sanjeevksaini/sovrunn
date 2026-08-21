package executiontarget

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

func TestProject_ClosedSafeJSONShape(t *testing.T) {
	t.Parallel()

	svc, _, _, _, fixtures := newQualifyLifecycle(t, nil)
	et := createActiveTarget(t, svc)
	fixtures.Set(et.Metadata.UID, presentFixture("rev-1", allSupportedTruths()))

	req := baseQualifyReq(et.Metadata.UID, "proj-shape")
	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	qualified, prob := svc.Qualify(context.Background(), req)
	if prob != nil {
		t.Fatalf("qualify: %#v", prob)
	}

	snap, ok := svc.CaptureProjectionSnapshot(qualified.Metadata.UID, activeBacking())
	if !ok {
		t.Fatal("expected snapshot")
	}
	proj := Project(snap)

	raw, err := json.Marshal(proj)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		t.Fatalf("unmarshal top: %v", err)
	}

	wantTop := map[string]bool{
		"metadata": true, "spec": true, "status": true, "effectiveAvailability": true,
	}
	if len(top) != len(wantTop) {
		t.Fatalf("top-level keys=%v want exactly %v", keysOf(top), keysOfBool(wantTop))
	}
	for k := range wantTop {
		if _, ok := top[k]; !ok {
			t.Fatalf("missing top-level key %q in %s", k, string(raw))
		}
	}

	var meta map[string]any
	if err := json.Unmarshal(top["metadata"], &meta); err != nil {
		t.Fatalf("metadata: %v", err)
	}
	assertExactKeys(t, meta, "uid", "name", "generation", "resourceVersion")
	if meta["uid"] != qualified.Metadata.UID || meta["name"] != qualified.Metadata.Name {
		t.Fatalf("metadata identity mismatch: %#v", meta)
	}
	if int64(meta["generation"].(float64)) != qualified.Metadata.Generation {
		t.Fatalf("generation=%v", meta["generation"])
	}
	if meta["resourceVersion"] != qualified.Metadata.ResourceVersion {
		t.Fatalf("resourceVersion=%v", meta["resourceVersion"])
	}

	var spec map[string]any
	if err := json.Unmarshal(top["spec"], &spec); err != nil {
		t.Fatalf("spec: %v", err)
	}
	assertExactKeys(t, spec, "cloudProviderParticipationRef", "infrastructureStackRef", "targetClass")
	if spec["targetClass"] != string(model.TargetClassSyntheticIaaS) {
		t.Fatalf("targetClass=%v", spec["targetClass"])
	}

	var status map[string]any
	if err := json.Unmarshal(top["status"], &status); err != nil {
		t.Fatalf("status: %v", err)
	}
	assertExactKeys(t, status, "lifecycle", "qualification")
	if status["lifecycle"] != string(model.LifecycleActive) {
		t.Fatalf("lifecycle=%v", status["lifecycle"])
	}
	if status["qualification"] != string(model.QualificationQualified) {
		t.Fatalf("qualification=%v", status["qualification"])
	}

	var ea string
	if err := json.Unmarshal(top["effectiveAvailability"], &ea); err != nil {
		t.Fatalf("effectiveAvailability: %v", err)
	}
	if ea != string(model.EffectiveAvailable) {
		t.Fatalf("effectiveAvailability=%q", ea)
	}

	// Risk 4 (task-scope): never disclose facts/results/observer/handles/links
	// or other non-approved members in the projected output. Nested TypedRef
	// apiVersion/kind/name/uid on the two immutable spec refs are part of those
	// approved spec fields.
	forbiddenTopOrStatus := []string{
		"conditions", "factSetRef", "qualificationResultRef", "maintenanceEpoch",
		"observedGeneration", "observerID", "observerRevision", "facts", "result",
		"handle", "availability", "labels", "annotations", "displayName",
		"createdAt", "updatedAt", "scopeRef", "CloudProviderScopeUID", "TypeMeta",
	}
	for _, f := range forbiddenTopOrStatus {
		if _, ok := top[f]; ok {
			t.Fatalf("forbidden top-level key %q present in projection: %s", f, string(raw))
		}
		if _, ok := status[f]; ok {
			t.Fatalf("forbidden status key %q present in projection: %s", f, string(raw))
		}
	}
	if _, ok := top["apiVersion"]; ok {
		t.Fatal("top-level apiVersion must not be projected")
	}
	if _, ok := top["kind"]; ok {
		t.Fatal("top-level kind must not be projected")
	}
	// Internal links exist on the committed target but must not project.
	if qualified.Status.FactSetRef == nil || qualified.Status.QualificationResultRef == nil {
		t.Fatal("qualified target must have internal links for this disclosure test")
	}
	if containsSubstring(string(raw), qualified.Status.FactSetRef.UID) {
		t.Fatalf("factSetRef uid leaked: %s", string(raw))
	}
	if containsSubstring(string(raw), qualified.Status.QualificationResultRef.UID) {
		t.Fatalf("qualificationResultRef uid leaked: %s", string(raw))
	}
}

func TestComputeEffectiveAvailability_OrderedTruthTable(t *testing.T) {
	t.Parallel()

	base := model.ExecutionTarget{
		Metadata: apimeta.ObjectMeta{
			UID:             "dddddddddddddddddddddddddddddddd",
			Name:            "target-1",
			Generation:      1,
			ResourceVersion: "2",
		},
		Spec: model.ExecutionTargetSpec{
			CloudProviderParticipationRef: participationRef("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			InfrastructureStackRef:        stackRef("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"),
			TargetClass:                   model.TargetClassSyntheticIaaS,
		},
		Status: model.ExecutionTargetStatus{
			Lifecycle:     model.LifecycleActive,
			Qualification: model.QualificationQualified,
		},
	}

	cases := []struct {
		name string
		snap ProjectionSnapshot
		want model.EffectiveAvailability
	}{
		{
			name: "retired-outranks-everything",
			snap: ProjectionSnapshot{
				Target: func() model.ExecutionTarget {
					et := base
					et.Status.Lifecycle = model.LifecycleRetired
					et.Status.Qualification = model.QualificationUnqualified
					return et
				}(),
				MaintenanceActive: true,
				FactsFresh:        true,
				BackingViable:     true,
			},
			want: model.EffectiveUnavailable,
		},
		{
			name: "active-maintenance-marker",
			snap: ProjectionSnapshot{
				Target:            base,
				MaintenanceActive: true,
				FactsFresh:        true,
				BackingViable:     true,
			},
			want: model.EffectiveMaintenance,
		},
		{
			name: "active-qualified-fresh-viable",
			snap: ProjectionSnapshot{
				Target:            base,
				MaintenanceActive: false,
				FactsFresh:        true,
				BackingViable:     true,
			},
			want: model.EffectiveAvailable,
		},
		{
			name: "active-qualified-stale-facts",
			snap: ProjectionSnapshot{
				Target:            base,
				MaintenanceActive: false,
				FactsFresh:        false,
				BackingViable:     true,
			},
			want: model.EffectiveUnavailable,
		},
		{
			name: "active-qualified-nonviable-backing",
			snap: ProjectionSnapshot{
				Target:            base,
				MaintenanceActive: false,
				FactsFresh:        true,
				BackingViable:     false,
			},
			want: model.EffectiveUnavailable,
		},
		{
			name: "active-unqualified",
			snap: ProjectionSnapshot{
				Target: func() model.ExecutionTarget {
					et := base
					et.Status.Qualification = model.QualificationUnqualified
					return et
				}(),
				FactsFresh:    false,
				BackingViable: true,
			},
			want: model.EffectiveUnavailable,
		},
		{
			name: "active-rejected",
			snap: ProjectionSnapshot{
				Target: func() model.ExecutionTarget {
					et := base
					et.Status.Qualification = model.QualificationRejected
					return et
				}(),
				FactsFresh:    true,
				BackingViable: true,
			},
			want: model.EffectiveUnavailable,
		},
		{
			name: "active-indeterminate",
			snap: ProjectionSnapshot{
				Target: func() model.ExecutionTarget {
					et := base
					et.Status.Qualification = model.QualificationIndeterminate
					return et
				}(),
				FactsFresh:    true,
				BackingViable: true,
			},
			want: model.EffectiveUnavailable,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ComputeEffectiveAvailability(tc.snap)
			if got != tc.want {
				t.Fatalf("got %s want %s", got, tc.want)
			}
			proj := Project(tc.snap)
			if proj.EffectiveAvailability != tc.want {
				t.Fatalf("Project effectiveAvailability=%s want %s", proj.EffectiveAvailability, tc.want)
			}
			if proj.Status.Lifecycle != tc.snap.Target.Status.Lifecycle {
				t.Fatalf("projected lifecycle=%s", proj.Status.Lifecycle)
			}
			if proj.Status.Qualification != tc.snap.Target.Status.Qualification {
				t.Fatalf("projected qualification=%s", proj.Status.Qualification)
			}
		})
	}
}

func TestCaptureProjectionSnapshot_UnderMutex_TruthTable(t *testing.T) {
	t.Parallel()

	t.Run("create-projects-unavailable", func(t *testing.T) {
		t.Parallel()
		svc, _, _, _, _ := newQualifyLifecycle(t, nil)
		et := createActiveTarget(t, svc)

		snap, ok := svc.CaptureProjectionSnapshot(et.Metadata.UID, activeBacking())
		if !ok {
			t.Fatal("expected snapshot")
		}
		if snap.FactsFresh {
			t.Fatal("create has no facts; FactsFresh must be false")
		}
		if snap.MaintenanceActive {
			t.Fatal("create must not have Maintenance marker")
		}
		proj := Project(snap)
		if proj.EffectiveAvailability != model.EffectiveUnavailable {
			t.Fatalf("got %s", proj.EffectiveAvailability)
		}
		if proj.Status.Lifecycle != model.LifecycleActive || proj.Status.Qualification != model.QualificationUnqualified {
			t.Fatalf("status=%#v", proj.Status)
		}
		if proj.Metadata.ResourceVersion != et.Metadata.ResourceVersion {
			t.Fatalf("etag changed: %s vs %s", proj.Metadata.ResourceVersion, et.Metadata.ResourceVersion)
		}
	})

	t.Run("qualified-fresh-viable-available", func(t *testing.T) {
		t.Parallel()
		svc, _, _, _, fixtures := newQualifyLifecycle(t, nil)
		et := createActiveTarget(t, svc)
		fixtures.Set(et.Metadata.UID, presentFixture("rev-1", allSupportedTruths()))
		req := baseQualifyReq(et.Metadata.UID, "proj-available")
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}
		qualified, prob := svc.Qualify(context.Background(), req)
		if prob != nil {
			t.Fatalf("qualify: %#v", prob)
		}

		snap, ok := svc.CaptureProjectionSnapshot(qualified.Metadata.UID, activeBacking())
		if !ok {
			t.Fatal("expected snapshot")
		}
		if !snap.FactsFresh || !snap.BackingViable || snap.MaintenanceActive {
			t.Fatalf("snap freshness/viability/marker=%#v", snap)
		}
		proj := Project(snap)
		if proj.EffectiveAvailability != model.EffectiveAvailable {
			t.Fatalf("got %s", proj.EffectiveAvailability)
		}
	})

	t.Run("maintenance-marker-projects-maintenance", func(t *testing.T) {
		t.Parallel()
		svc, store, _, _, fixtures := newQualifyLifecycle(t, nil)
		et := createActiveTarget(t, svc)
		fixtures.Set(et.Metadata.UID, presentFixture("rev-1", allSupportedTruths()))
		req := baseQualifyReq(et.Metadata.UID, "proj-maint")
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}
		qualified, prob := svc.Qualify(context.Background(), req)
		if prob != nil {
			t.Fatalf("qualify: %#v", prob)
		}

		enter := MaintenanceCommitRequest{
			TargetUID:                      qualified.Metadata.UID,
			ExpectedInfrastructureStackGen: qualified.Status.ObservedGeneration,
			ExpectedMaintenanceEpoch:       qualified.Status.MaintenanceEpoch,
			Actor:                          ExpirySystemActor(),
			RequestID:                      "req-maint-enter",
			AuditUID:                       "audit-maint-enter",
		}
		after, prob := svc.CommitMaintenanceEnter(context.Background(), enter)
		if prob != nil {
			t.Fatalf("maintenance enter: %#v", prob)
		}
		if marker, ok := store.GetMaintenanceMarker(after.Metadata.UID); !ok || !marker.Active {
			t.Fatalf("marker=%#v ok=%v", marker, ok)
		}

		snap, ok := svc.CaptureProjectionSnapshot(after.Metadata.UID, activeBacking())
		if !ok {
			t.Fatal("expected snapshot")
		}
		if !snap.MaintenanceActive {
			t.Fatal("MaintenanceActive must be true")
		}
		proj := Project(snap)
		if proj.EffectiveAvailability != model.EffectiveMaintenance {
			t.Fatalf("got %s", proj.EffectiveAvailability)
		}
		if proj.Status.Qualification != model.QualificationUnqualified {
			t.Fatalf("qualification=%s", proj.Status.Qualification)
		}
	})

	t.Run("retired-projects-unavailable", func(t *testing.T) {
		t.Parallel()
		svc, _, _, _, _ := newQualifyLifecycle(t, nil)
		et := createActiveTarget(t, svc)
		retire := RetireCommitRequest{
			TargetUID:               et.Metadata.UID,
			ExpectedResourceVersion: et.Metadata.ResourceVersion,
			Actor: apimeta.TypedRef{
				APIVersion: "identity.sovrunn.io/v1alpha1",
				Kind:       "Principal",
				Name:       "principal",
				UID:        "principal-1",
			},
			RequestID:   "req-retire-proj",
			AuditUID:    "audit-retire-proj",
			Idempotency: retireNS(et.Metadata.UID, "k-retire-proj"),
			Digest:      DigestAction(),
			Completion:  CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{}`)},
		}
		if res := svc.ReserveOrInspectReplay(retire.Idempotency, retire.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}
		retired, prob := svc.CommitRetire(context.Background(), retire)
		if prob != nil {
			t.Fatalf("retire: %#v", prob)
		}

		snap, ok := svc.CaptureProjectionSnapshot(retired.Metadata.UID, activeBacking())
		if !ok {
			t.Fatal("expected snapshot")
		}
		proj := Project(snap)
		if proj.EffectiveAvailability != model.EffectiveUnavailable {
			t.Fatalf("got %s", proj.EffectiveAvailability)
		}
		if proj.Status.Lifecycle != model.LifecycleRetired {
			t.Fatalf("lifecycle=%s", proj.Status.Lifecycle)
		}
	})

	t.Run("fact-expiry-projects-unavailable-unchanged-etag", func(t *testing.T) {
		t.Parallel()
		now := time.Date(2026, 8, 21, 10, 0, 0, 0, time.UTC)
		svc, _, _, _, fixtures := newQualifyLifecycle(t, nil)
		svc.now = fixedNow(now)
		et := createActiveTarget(t, svc)
		fixtures.Set(et.Metadata.UID, presentFixture("rev-1", allSupportedTruths()))
		req := baseQualifyReq(et.Metadata.UID, "proj-expiry")
		if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}
		qualified, prob := svc.Qualify(context.Background(), req)
		if prob != nil {
			t.Fatalf("qualify: %#v", prob)
		}
		beforeRV := qualified.Metadata.ResourceVersion

		snapFresh, ok := svc.CaptureProjectionSnapshot(qualified.Metadata.UID, activeBacking())
		if !ok || Project(snapFresh).EffectiveAvailability != model.EffectiveAvailable {
			t.Fatalf("pre-expiry want Available, snap=%#v", snapFresh)
		}

		// Advance clock past expiresAt and run expiry tick (freshness input only).
		expiredAt := now.Add(FactFreshnessWindow + time.Second)
		svc.now = fixedNow(expiredAt)
		svc.ProcessExpiryTick(context.Background(), expiredAt)

		snapStale, ok := svc.CaptureProjectionSnapshot(qualified.Metadata.UID, activeBacking())
		if !ok {
			t.Fatal("expected snapshot")
		}
		if snapStale.FactsFresh {
			t.Fatal("facts must not be fresh after expiry")
		}
		proj := Project(snapStale)
		if proj.EffectiveAvailability != model.EffectiveUnavailable {
			t.Fatalf("got %s", proj.EffectiveAvailability)
		}
		if proj.Metadata.ResourceVersion != beforeRV {
			t.Fatalf("ETag changed on expiry: %s -> %s", beforeRV, proj.Metadata.ResourceVersion)
		}
		cur, _ := svc.Store().GetExecutionTarget(qualified.Metadata.UID)
		if cur.Metadata.ResourceVersion != beforeRV {
			t.Fatalf("target mutated on expiry: %s -> %s", beforeRV, cur.Metadata.ResourceVersion)
		}
	})

	t.Run("missing-target", func(t *testing.T) {
		t.Parallel()
		svc, _, _, _, _ := newQualifyLifecycle(t, nil)
		_, ok := svc.CaptureProjectionSnapshot("missing-uid-000000000000000000000000", activeBacking())
		if ok {
			t.Fatal("expected ok=false for missing target")
		}
	})
}

func TestProject_ViabilityOnlyChange_NoMutationOrETagChange(t *testing.T) {
	t.Parallel()

	svc, _, _, _, fixtures := newQualifyLifecycle(t, nil)
	et := createActiveTarget(t, svc)
	fixtures.Set(et.Metadata.UID, presentFixture("rev-1", allSupportedTruths()))
	req := baseQualifyReq(et.Metadata.UID, "proj-viability")
	if res := svc.ReserveOrInspectReplay(req.Idempotency, req.Digest); res.Outcome != ReservationOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}
	qualified, prob := svc.Qualify(context.Background(), req)
	if prob != nil {
		t.Fatalf("qualify: %#v", prob)
	}
	beforeRV := qualified.Metadata.ResourceVersion
	beforeGen := qualified.Metadata.Generation

	viable := activeBacking()
	snapViable, ok := svc.CaptureProjectionSnapshot(qualified.Metadata.UID, viable)
	if !ok {
		t.Fatal("expected snapshot")
	}
	projViable := Project(snapViable)
	if projViable.EffectiveAvailability != model.EffectiveAvailable {
		t.Fatalf("viable backing got %s", projViable.EffectiveAvailability)
	}

	nonViable := viable
	nonViable.ParticipationEffectiveActive = false
	snapNonViable, ok := svc.CaptureProjectionSnapshot(qualified.Metadata.UID, nonViable)
	if !ok {
		t.Fatal("expected snapshot")
	}
	projNonViable := Project(snapNonViable)
	if projNonViable.EffectiveAvailability != model.EffectiveUnavailable {
		t.Fatalf("non-viable backing got %s", projNonViable.EffectiveAvailability)
	}

	// Stack phase non-viable likewise.
	stackDown := viable
	stackDown.StackPhase = "Failed"
	snapStack, ok := svc.CaptureProjectionSnapshot(qualified.Metadata.UID, stackDown)
	if !ok {
		t.Fatal("expected snapshot")
	}
	if Project(snapStack).EffectiveAvailability != model.EffectiveUnavailable {
		t.Fatal("failed stack must project Unavailable")
	}

	// Restore viable → Available again, still no mutation.
	snapRestored, ok := svc.CaptureProjectionSnapshot(qualified.Metadata.UID, viable)
	if !ok {
		t.Fatal("expected snapshot")
	}
	if Project(snapRestored).EffectiveAvailability != model.EffectiveAvailable {
		t.Fatal("restored viable backing must project Available")
	}

	cur, ok := svc.Store().GetExecutionTarget(qualified.Metadata.UID)
	if !ok {
		t.Fatal("target missing")
	}
	if cur.Metadata.ResourceVersion != beforeRV || cur.Metadata.Generation != beforeGen {
		t.Fatalf("viability-only change mutated target: before rv=%s gen=%d after rv=%s gen=%d",
			beforeRV, beforeGen, cur.Metadata.ResourceVersion, cur.Metadata.Generation)
	}
	if projViable.Metadata.ResourceVersion != beforeRV ||
		projNonViable.Metadata.ResourceVersion != beforeRV ||
		Project(snapStack).Metadata.ResourceVersion != beforeRV {
		t.Fatal("projected ETag must remain unchanged across viability-only changes")
	}
}

func TestProject_OmitsInternalLinksEvenWhenPresentOnTarget(t *testing.T) {
	t.Parallel()

	et := model.ExecutionTarget{
		Metadata: apimeta.ObjectMeta{
			UID:             "dddddddddddddddddddddddddddddddd",
			Name:            "target-1",
			Generation:      1,
			ResourceVersion: "5",
		},
		Spec: model.ExecutionTargetSpec{
			CloudProviderParticipationRef: participationRef("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			InfrastructureStackRef:        stackRef("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"),
			TargetClass:                   model.TargetClassSyntheticIaaS,
		},
		Status: model.ExecutionTargetStatus{
			Lifecycle:          model.LifecycleActive,
			Qualification:      model.QualificationQualified,
			MaintenanceEpoch:   2,
			ObservedGeneration: 9,
			FactSetRef: &apimeta.TypedRef{
				APIVersion: model.APIVersionFactSet,
				Kind:       model.KindNormalizedTargetFactSet,
				Name:       "fs-1",
				UID:        "fs-1",
			},
			QualificationResultRef: &apimeta.TypedRef{
				APIVersion: model.APIVersionQualificationResult,
				Kind:       model.KindTargetQualificationResult,
				Name:       "qr-1",
				UID:        "qr-1",
			},
		},
		CloudProviderScopeUID: "cccccccccccccccccccccccccccccccc",
	}

	proj := Project(ProjectionSnapshot{
		Target:            et,
		MaintenanceActive: false,
		FactsFresh:        true,
		BackingViable:     true,
	})
	raw, err := json.Marshal(proj)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	blob := string(raw)
	for _, key := range []string{
		"factSetRef", "qualificationResultRef", "maintenanceEpoch",
		"observedGeneration", "CloudProviderScopeUID", "conditions", "observerID",
	} {
		if containsJSONKey(blob, key) {
			t.Fatalf("forbidden key %q present: %s", key, blob)
		}
	}
	for _, leak := range []string{"fs-1", "qr-1"} {
		if containsSubstring(blob, leak) {
			t.Fatalf("internal link value %q leaked: %s", leak, blob)
		}
	}
	if proj.EffectiveAvailability != model.EffectiveAvailable {
		t.Fatalf("got %s", proj.EffectiveAvailability)
	}
}

func assertExactKeys(t *testing.T, m map[string]any, want ...string) {
	t.Helper()
	if len(m) != len(want) {
		t.Fatalf("keys=%v want exactly %v", keysOfAny(m), want)
	}
	for _, k := range want {
		if _, ok := m[k]; !ok {
			t.Fatalf("missing key %q in %v", k, keysOfAny(m))
		}
	}
}

func keysOf(m map[string]json.RawMessage) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func keysOfBool(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func keysOfAny(m map[string]any) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// containsJSONKey reports whether raw JSON text includes `"key"` as a quoted token.
func containsJSONKey(blob, key string) bool {
	return strings.Contains(blob, `"`+key+`"`)
}

func containsSubstring(s, sub string) bool {
	return strings.Contains(s, sub)
}
