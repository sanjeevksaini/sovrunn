package cloudmodel_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiconform"
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

func actorRef(uid string) apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: "identity.sovrunn.io/v1alpha1",
		Kind:       "Principal",
		Name:       "principal",
		UID:        uid,
	}
}

func subjectPlatform(uid string) apimeta.TypedRef {
	return apimeta.TypedRef{
		APIVersion: model.APIVersionCloudPlatform,
		Kind:       model.KindCloudPlatform,
		Name:       "plat",
		UID:        uid,
	}
}

func newPub(t *testing.T) (*cloudmodel.Store, *cloudmodel.MemoryAuditAppender, *cloudmodel.IdempotencyCoordinator, *cloudmodel.PublicationCoordinator) {
	t.Helper()
	store := cloudmodel.NewStore()
	audit := &cloudmodel.MemoryAuditAppender{}
	idemp := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
	pub := cloudmodel.NewPublicationCoordinator(cloudmodel.PublicationCoordinatorConfig{
		Store:       store,
		Audit:       audit,
		Idempotency: idemp,
	})
	return store, audit, idemp, pub
}

func TestPublication_AuditedSuccessCollectionCreate(t *testing.T) {
	t.Parallel()
	store, audit, idemp, pub := newPub(t)
	ns := testNS("p1", "POST /v1/cloud-platforms", "", "k1")
	d := mustDigest(t, `{"metadata":{"name":"plat"}}`, ns.Pattern, "")
	if res := idemp.Reserve(context.Background(), ns, d, nil); res.Outcome != cloudmodel.ReserveOutcomeReserved {
		t.Fatalf("reserve: %#v", res)
	}

	cp := platformScoped("plat", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	store.BeginPublication()
	staged, prob := store.StageCreateCloudPlatform(cp)
	if prob != nil {
		t.Fatalf("stage: %#v", prob)
	}
	event := cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
		UID:       "audit-1",
		RequestID: "req-1",
		Action:    cloudmodel.AuditActionCollectionCreate,
		Outcome:   apiconform.AuditOutcomeSucceeded,
		Reason:    "created",
		Actor:     actorRef("principal-1"),
		Subject:   subjectPlatform(cp.Metadata.UID),
	})
	body, _ := json.Marshal(cp)
	pubProb := pub.PublishSuccess(context.Background(), cloudmodel.AuditedSuccess{
		Staged: staged,
		Event:  event,
		Idempotency: &cloudmodel.IdempotencyCompletion{
			Namespace: ns,
			Digest:    d,
			Result:    cloudmodel.CompletedResult{StatusCode: http.StatusCreated, Body: body},
		},
	})
	store.EndPublication()
	if pubProb != nil {
		t.Fatalf("publish: %#v", pubProb)
	}
	if _, ok := store.GetCloudPlatform(cp.Metadata.UID); !ok {
		t.Fatal("mutation must be published after audit append")
	}
	if audit.Len() != 1 || audit.Events()[0].Record.RequestID != "req-1" {
		t.Fatalf("audit events: %#v", audit.Events())
	}
	if st, ok := idemp.StateForTest(ns); !ok || st != cloudmodel.ReservationCompleted {
		t.Fatalf("idempotency: %s ok=%v", st, ok)
	}
}

func TestPublication_AuditedSuccessPatchAndParticipation(t *testing.T) {
	t.Parallel()
	store, audit, _, pub := newPub(t)
	cp := platformScoped("plat", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if _, prob := store.CreateCloudPlatform(cp); prob != nil {
		t.Fatalf("create: %#v", prob)
	}
	cp.Spec.Description = "updated"
	store.BeginPublication()
	staged, prob := store.StageUpdateCloudPlatform(cp)
	if prob != nil {
		t.Fatalf("stage patch: %#v", prob)
	}
	if pubProb := pub.PublishSuccess(context.Background(), cloudmodel.AuditedSuccess{
		Staged: staged,
		Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
			UID: "audit-patch", RequestID: "req-patch",
			Action: cloudmodel.AuditActionResourcePatch, Outcome: apiconform.AuditOutcomeSucceeded,
			Reason: "patched", Actor: actorRef("p1"), Subject: subjectPlatform(cp.Metadata.UID),
		}),
	}); pubProb != nil {
		t.Fatalf("patch publish: %#v", pubProb)
	}
	store.EndPublication()

	prov := providerScoped("prov", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	if _, prob := store.CreateCloudProvider(prov); prob != nil {
		t.Fatalf("provider: %#v", prob)
	}
	part := model.CloudProviderParticipation{
		Metadata: apimeta.ObjectMeta{
			Name: "part", UID: "cccccccccccccccccccccccccccccccc",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
				Name: "plat", UID: cp.Metadata.UID,
			}},
		},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform, Name: "plat", UID: cp.Metadata.UID},
			CloudProviderRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider, Name: "prov", UID: prov.Metadata.UID},
			Environment:      model.ParticipationEnvironmentDevelopment,
		},
		Status: cloudmodel.InitialParticipationStatus(time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)),
	}
	store.BeginPublication()
	stagedPart, prob := store.StageCreateParticipation(part)
	if prob != nil {
		t.Fatalf("stage part: %#v", prob)
	}
	if pubProb := pub.PublishSuccess(context.Background(), cloudmodel.AuditedSuccess{
		Staged: stagedPart,
		Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
			UID: "audit-part-create", RequestID: "req-part",
			Action: cloudmodel.AuditActionParticipationCreate, Outcome: apiconform.AuditOutcomeSucceeded,
			Reason: "requested", Actor: actorRef("p1"),
			Subject: apimeta.TypedRef{APIVersion: model.APIVersionCloudProviderParticipation, Kind: model.KindCloudProviderParticipation, Name: "part", UID: part.Metadata.UID},
		}),
	}); pubProb != nil {
		t.Fatalf("part create: %#v", pubProb)
	}
	store.EndPublication()

	got, ok := store.GetParticipation(part.Metadata.UID)
	if !ok {
		t.Fatal("participation missing")
	}
	next, prob := cloudmodel.ApplyParticipationAction(got.Status, cloudmodel.ParticipationActionAccept, "")
	if prob != nil {
		t.Fatalf("accept: %#v", prob)
	}
	got.Status = next
	store.BeginPublication()
	stagedAction, prob := store.StageUpdateParticipation(got)
	if prob != nil {
		t.Fatalf("stage action: %#v", prob)
	}
	if pubProb := pub.PublishSuccess(context.Background(), cloudmodel.AuditedSuccess{
		Staged: stagedAction,
		Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
			UID: "audit-part-action", RequestID: "req-accept",
			Action: cloudmodel.AuditActionParticipationAction, Outcome: apiconform.AuditOutcomeSucceeded,
			Reason: "accepted", Actor: actorRef("p1"),
			Subject: apimeta.TypedRef{APIVersion: model.APIVersionCloudProviderParticipation, Kind: model.KindCloudProviderParticipation, Name: "part", UID: part.Metadata.UID},
		}),
	}); pubProb != nil {
		t.Fatalf("action publish: %#v", pubProb)
	}
	store.EndPublication()

	if audit.Len() != 3 {
		t.Fatalf("audit len = %d", audit.Len())
	}
}

func TestPublication_AuditedDenials(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name   string
		action string
	}{
		{"status-write", cloudmodel.AuditActionStatusWrite},
		{"system-owned-field-write", cloudmodel.AuditActionSystemOwnedFieldWrite},
		{"forged-grant", cloudmodel.AuditActionForgedGrant},
		{"root-denial", cloudmodel.AuditActionRootDenial},
		{"list-without-read-grant", cloudmodel.AuditActionListWithoutReadGrant},
		{"safe-denial", cloudmodel.AuditActionSafeDenial},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, audit, _, pub := newPub(t)
			denial := apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail(tc.name)
			got := pub.PublishDenial(context.Background(), cloudmodel.AuditedDenial{
				Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
					UID: "audit-" + tc.name, RequestID: "req-" + tc.name,
					Action: tc.action, Outcome: apiconform.AuditOutcomeDenied,
					Reason: tc.name, Actor: actorRef("p1"),
					Subject: subjectPlatform("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
				}),
				Problem: denial,
			})
			if got != denial {
				t.Fatalf("got %#v", got)
			}
			if audit.Len() != 1 {
				t.Fatalf("audit len = %d", audit.Len())
			}
			if audit.Events()[0].Record.RequestID != "req-"+tc.name {
				t.Fatalf("correlation missing: %#v", audit.Events()[0])
			}
		})
	}
}

func TestPublication_AppendFailureSuppressesSuccess(t *testing.T) {
	t.Parallel()
	store, audit, idemp, pub := newPub(t)
	audit.SetFail(errors.New("append boom"))
	ns := testNS("p1", "POST /v1/cloud-platforms", "", "k-fail")
	d := mustDigest(t, `{"metadata":{"name":"plat"}}`, ns.Pattern, "")
	_ = idemp.Reserve(context.Background(), ns, d, nil)

	cp := platformScoped("plat", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	store.BeginPublication()
	staged, prob := store.StageCreateCloudPlatform(cp)
	if prob != nil {
		t.Fatalf("stage: %#v", prob)
	}
	pubProb := pub.PublishSuccess(context.Background(), cloudmodel.AuditedSuccess{
		Staged: staged,
		Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
			UID: "audit-fail", RequestID: "req-fail",
			Action: cloudmodel.AuditActionCollectionCreate, Outcome: apiconform.AuditOutcomeSucceeded,
			Reason: "created", Actor: actorRef("p1"), Subject: subjectPlatform(cp.Metadata.UID),
		}),
		Idempotency: &cloudmodel.IdempotencyCompletion{
			Namespace: ns, Digest: d,
			Result: cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{}`)},
		},
	})
	store.EndPublication()

	if pubProb == nil || pubProb.Code != apiproblem.CodeInternalError {
		t.Fatalf("want INTERNAL_ERROR, got %#v", pubProb)
	}
	if pubProb.Code == apiproblem.CodeDependencyUnavailable {
		t.Fatal("DEPENDENCY_UNAVAILABLE must not be used for audit-append failure")
	}
	if _, ok := store.GetCloudPlatform(cp.Metadata.UID); ok {
		t.Fatal("mutation must not be published on append failure")
	}
	if audit.Len() != 0 {
		t.Fatal("failed append must not retain event")
	}
	if st, ok := idemp.StateForTest(ns); !ok || st != cloudmodel.ReservationAborted {
		t.Fatalf("idempotency must abort: %s ok=%v", st, ok)
	}
}

func TestPublication_AppendFailureSuppressesDenial(t *testing.T) {
	t.Parallel()
	_, audit, _, pub := newPub(t)
	audit.SetFail(errors.New("append boom"))
	denial := apiproblem.New(apiproblem.CodeAuthorizationDenied)
	got := pub.PublishDenial(context.Background(), cloudmodel.AuditedDenial{
		Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
			UID: "audit-denial-fail", RequestID: "req-d",
			Action: cloudmodel.AuditActionStatusWrite, Outcome: apiconform.AuditOutcomeDenied,
			Reason: "status", Actor: actorRef("p1"), Subject: subjectPlatform("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		}),
		Problem: denial,
	})
	if got == nil || got.Code != apiproblem.CodeInternalError {
		t.Fatalf("want INTERNAL_ERROR, got %#v", got)
	}
	if got == denial {
		t.Fatal("original denial must be suppressed")
	}
	if strings.Contains(string(got.Code), "DEPENDENCY") {
		t.Fatal("no DEPENDENCY_UNAVAILABLE")
	}
}

func TestPublication_AppendThenPublishFailureLeavesAudit(t *testing.T) {
	t.Parallel()
	store, audit, idemp, pub := newPub(t)
	ns := testNS("p1", "POST /v1/cloud-platforms", "", "k-postfail")
	d := mustDigest(t, `{"metadata":{"name":"plat"}}`, ns.Pattern, "")
	_ = idemp.Reserve(context.Background(), ns, d, nil)

	cp := platformScoped("plat", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	store.BeginPublication()
	staged, _ := store.StageCreateCloudPlatform(cp)
	pubProb := pub.PublishSuccessAllowingPostAppendFailure(
		context.Background(),
		cloudmodel.AuditedSuccess{
			Staged: staged,
			Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
				UID: "audit-remain", RequestID: "req-remain",
				Action: cloudmodel.AuditActionCollectionCreate, Outcome: apiconform.AuditOutcomeSucceeded,
				Reason: "created", Actor: actorRef("p1"), Subject: subjectPlatform(cp.Metadata.UID),
			}),
			Idempotency: &cloudmodel.IdempotencyCompletion{
				Namespace: ns, Digest: d,
				Result: cloudmodel.CompletedResult{StatusCode: 201, Body: []byte(`{}`)},
			},
		},
		func() error { return errors.New("publish boom") },
	)
	store.EndPublication()

	if pubProb == nil || pubProb.Code != apiproblem.CodeInternalError {
		t.Fatalf("want INTERNAL_ERROR, got %#v", pubProb)
	}
	if audit.Len() != 1 {
		t.Fatalf("already-appended AuditEvent may remain: len=%d", audit.Len())
	}
	if _, ok := store.GetCloudPlatform(cp.Metadata.UID); ok {
		t.Fatal("mutation must remain unpublished")
	}
	if st, ok := idemp.StateForTest(ns); !ok || st != cloudmodel.ReservationAborted {
		t.Fatalf("idempotency must not complete: %s ok=%v", st, ok)
	}
}

func TestPublication_RedactedAuditHasCorrelationNoSecrets(t *testing.T) {
	t.Parallel()
	_, audit, _, pub := newPub(t)
	secretish := "Idempotency-Key: raw-secret-key\nAuthorization: Bearer tok\n{\"password\":\"x\"}"
	_ = pub.PublishDenial(context.Background(), cloudmodel.AuditedDenial{
		Event: cloudmodel.NewRedactedAuditEvent(cloudmodel.RedactedAuditInput{
			UID: "audit-redact", RequestID: "corr-123",
			Action: cloudmodel.AuditActionSafeDenial, Outcome: apiconform.AuditOutcomeDenied,
			Reason: "safe", Actor: actorRef("p1"),
			Subject: subjectPlatform("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		}),
		Problem: apiproblem.New(apiproblem.CodeResourceNotFound),
	})
	raw, err := json.Marshal(audit.Events()[0])
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)
	if !strings.Contains(s, "corr-123") {
		t.Fatal("request correlation required")
	}
	if strings.Contains(s, secretish) || strings.Contains(s, "raw-secret-key") || strings.Contains(s, "Bearer tok") || strings.Contains(s, "password") {
		t.Fatalf("secrets leaked in AuditEvent: %s", s)
	}
}

func TestPublication_ExpiryAppendBeforePublication(t *testing.T) {
	t.Parallel()
	store, audit, _, pub := newPub(t)
	cp := platformScoped("plat", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	prov := providerScoped("prov", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	_, _ = store.CreateCloudPlatform(cp)
	_, _ = store.CreateCloudProvider(prov)
	part := model.CloudProviderParticipation{
		Metadata: apimeta.ObjectMeta{
			Name: "part", UID: "cccccccccccccccccccccccccccccccc", ResourceVersion: "1",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
				Name: "plat", UID: cp.Metadata.UID,
			}},
		},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform, Name: "plat", UID: cp.Metadata.UID},
			CloudProviderRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider, Name: "prov", UID: prov.Metadata.UID},
			Environment:      model.ParticipationEnvironmentDevelopment,
		},
		Status: cloudmodel.InitialParticipationStatus(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
	}
	if _, prob := store.CreateParticipation(part); prob != nil {
		t.Fatalf("create: %#v", prob)
	}

	if err := pub.ExpireParticipation(context.Background(), part.Metadata.UID, "1", "sched-exec-1"); err != nil {
		t.Fatalf("expire: %v", err)
	}
	got, ok := store.GetParticipation(part.Metadata.UID)
	if !ok || got.Status.Phase != model.ParticipationPhaseExpired {
		t.Fatalf("phase = %#v ok=%v", got.Status.Phase, ok)
	}
	if audit.Len() != 1 {
		t.Fatalf("audit len = %d", audit.Len())
	}
	ev := audit.Events()[0]
	if ev.Record.Action != cloudmodel.AuditActionParticipationExpiry {
		t.Fatalf("action = %s", ev.Record.Action)
	}
	if ev.Record.ActorRef.UID != cloudmodel.SchedulerSystemActorUID {
		t.Fatalf("system actor = %#v", ev.Record.ActorRef)
	}
	if ev.Record.RequestID != "sched-exec-1" {
		t.Fatalf("scheduler execution correlation = %s", ev.Record.RequestID)
	}
}

func TestPublication_ExpiryAppendFailureNoPublication(t *testing.T) {
	t.Parallel()
	store, audit, _, pub := newPub(t)
	audit.SetFail(errors.New("append boom"))
	cp := platformScoped("plat", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	prov := providerScoped("prov", "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	_, _ = store.CreateCloudPlatform(cp)
	_, _ = store.CreateCloudProvider(prov)
	part := model.CloudProviderParticipation{
		Metadata: apimeta.ObjectMeta{
			Name: "part", UID: "cccccccccccccccccccccccccccccccc", ResourceVersion: "1",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
				Name: "plat", UID: cp.Metadata.UID,
			}},
		},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform, Name: "plat", UID: cp.Metadata.UID},
			CloudProviderRef: apimeta.TypedRef{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider, Name: "prov", UID: prov.Metadata.UID},
			Environment:      model.ParticipationEnvironmentDevelopment,
		},
		Status: cloudmodel.InitialParticipationStatus(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
	}
	if _, prob := store.CreateParticipation(part); prob != nil {
		t.Fatalf("create: %#v", prob)
	}
	err := pub.ExpireParticipation(context.Background(), part.Metadata.UID, "1", "sched-exec-fail")
	if err == nil || err.Error() != string(apiproblem.CodeInternalError) {
		t.Fatalf("want INTERNAL_ERROR, got %v", err)
	}
	got, _ := store.GetParticipation(part.Metadata.UID)
	if got.Status.Phase != model.ParticipationPhasePending {
		t.Fatalf("phase must remain Pending, got %s", got.Status.Phase)
	}
	if audit.Len() != 0 {
		t.Fatal("no audit on append failure")
	}
}
