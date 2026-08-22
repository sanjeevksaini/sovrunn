package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget"
	etmodel "github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

const (
	testPartUID  = "11111111111111111111111111111111"
	testStackUID = "22222222222222222222222222222222"
	etCollection = "/apis/execution.sovrunn.io/v1alpha1/execution-targets"
)

type etRuntime struct {
	cloud   *cloudmodel.Store
	audit   *cloudmodel.MemoryAuditAppender
	life    *executiontarget.ExecutionTargetLifecycleService
	grants  GrantResolver
	handler *ExecutionTargetCollectionHandler
}

func newETRuntime(grants GrantResolver) *etRuntime {
	cloud := cloudmodel.NewStore()
	audit := &cloudmodel.MemoryAuditAppender{}
	life := executiontarget.NewExecutionTargetLifecycleService(executiontarget.LifecycleConfig{
		Audit: audit,
	})
	h := NewExecutionTargetCollectionHandler(life, cloud, grants, audit)
	return &etRuntime{cloud: cloud, audit: audit, life: life, grants: grants, handler: h}
}

func executionTargetGrants(principal, providerUID string, write, read bool) *staticGrants {
	g := &staticGrants{principal: principal}
	scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: providerUID}
	if write {
		g.grants = append(g.grants, cloudmodel.Grant{Action: ActionExecutionTargetWrite, Scope: scope})
	}
	if read {
		g.grants = append(g.grants, cloudmodel.Grant{Action: ActionExecutionTargetRead, Scope: scope})
	}
	return g
}

func seedActiveParticipationAndStack(t *testing.T, rt *etRuntime) (model.CloudProviderParticipation, model.InfrastructureStack) {
	t.Helper()
	seedPlatformAndProvider(t, &f15Runtime{store: rt.cloud}, testProviderUID)

	part, prob := rt.cloud.CreateParticipation(model.CloudProviderParticipation{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: model.APIVersionCloudProviderParticipation,
			Kind:       model.KindCloudProviderParticipation,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "part-active", UID: testPartUID, ResourceVersion: "1", Generation: 1,
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: string(apimeta.ScopeCloudPlatform),
				Name: "root-plat", UID: testPlatformUID,
			}},
		},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform,
				Name: "root-plat", UID: testPlatformUID,
			},
			CloudProviderRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider,
				Name: "prov-a", UID: testProviderUID,
			},
			Environment: model.ParticipationEnvironmentDevelopment,
		},
		Status: model.CloudProviderParticipationStatus{Phase: model.ParticipationPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed participation: %#v", prob)
	}

	stack, prob := rt.cloud.CreateInfrastructureStack(model.InfrastructureStack{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: model.APIVersionInfrastructureStack,
			Kind:       model.KindInfrastructureStack,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "stack-active", UID: testStackUID, ResourceVersion: "1", Generation: 3,
			ScopeRef: providerScopeRef(testProviderUID, "prov-a"),
		},
		Spec: model.InfrastructureStackSpec{
			FaultDomainRef: apimeta.TypedRef{
				APIVersion: model.APIVersionFaultDomain, Kind: model.KindFaultDomain,
				Name: "fd", UID: "33333333333333333333333333333333",
			},
		},
		Status: model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed stack: %#v", prob)
	}
	return part, stack
}

func executionTargetCreateBody(name, partUID, stackUID string) []byte {
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{
			"cloudProviderParticipationRef":{"uid":"` + partUID + `"},
			"infrastructureStackRef":{"uid":"` + stackUID + `"},
			"targetClass":"synthetic-iaas"
		}
	}`)
}

func createHeaders(key string) map[string]string {
	return map[string]string{
		"Content-Type":    "application/json",
		"Idempotency-Key": key,
		"Accept":          "application/xml",
	}
}

func assertJSONNosniff(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if got := rec.Header().Get("Content-Type"); got != mediaJSON {
		t.Fatalf("Content-Type=%q, want %q", got, mediaJSON)
	}
	if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options=%q, want nosniff", got)
	}
}

func assertProblemNosniff(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if got := rec.Header().Get("Content-Type"); got != apiproblem.MediaTypeProblemJSON {
		t.Fatalf("Content-Type=%q, want problem+json", got)
	}
	if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options=%q, want nosniff", got)
	}
}

func TestExecutionTargetCollection_CreateHappyPathAndReplay(t *testing.T) {
	rt := newETRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
	seedActiveParticipationAndStack(t, rt)

	headers := createHeaders("et-create-1")
	rec := doMux(t, rt.handler, http.MethodPost, etCollection, executionTargetCreateBody("target-a", testPartUID, testStackUID), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
	}
	assertJSONNosniff(t, rec)
	if rec.Header().Get("ETag") == "" {
		t.Fatal("create must return ETag")
	}
	if rec.Header().Get("Location") != "" {
		t.Fatal("create must not return Location")
	}
	if rec.Header().Get("Idempotency-Key") != "" {
		t.Fatal("response must not echo Idempotency-Key")
	}

	var created executiontarget.SafeExecutionTarget
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.Status.Lifecycle != etmodel.LifecycleActive ||
		created.Status.Qualification != etmodel.QualificationUnqualified ||
		created.EffectiveAvailability != etmodel.EffectiveUnavailable {
		t.Fatalf("created projection=%#v", created)
	}
	if created.Metadata.Generation != 1 || created.Spec.TargetClass != etmodel.TargetClassSyntheticIaaS {
		t.Fatalf("created metadata/spec=%#v", created)
	}
	if rt.audit.Len() != 1 {
		t.Fatalf("create audit=%d, want 1", rt.audit.Len())
	}
	stored, ok := rt.life.Store().GetExecutionTarget(created.Metadata.UID)
	if !ok {
		t.Fatal("target must be published")
	}
	if stored.Status.MaintenanceEpoch != 0 || stored.Status.ObservedGeneration != 3 {
		t.Fatalf("epoch/observedGeneration=%#v", stored.Status)
	}
	if stored.CloudProviderScopeUID != testProviderUID {
		t.Fatalf("derived scope=%q", stored.CloudProviderScopeUID)
	}

	// Completed replay: same key/digest → original body, fresh correlation, no second audit.
	before := rt.audit.Len()
	firstBody := append([]byte(nil), rec.Body.Bytes()...)
	rec2 := doMux(t, rt.handler, http.MethodPost, etCollection, executionTargetCreateBody("target-a", testPartUID, testStackUID), headers)
	if rec2.Code != http.StatusCreated {
		t.Fatalf("replay status=%d body=%s", rec2.Code, rec2.Body.String())
	}
	assertJSONNosniff(t, rec2)
	if !bytes.Equal(bytes.TrimSpace(rec2.Body.Bytes()), bytes.TrimSpace(firstBody)) {
		t.Fatalf("replay body mismatch\nfirst=%s\nreplay=%s", firstBody, rec2.Body.Bytes())
	}
	if rt.audit.Len() != before {
		t.Fatal("completed replay must not append a second AuditEvent")
	}
	if rec2.Header().Get("Idempotency-Key") != "" {
		t.Fatal("replay must not echo Idempotency-Key")
	}
}

func TestExecutionTargetCollection_ListAscendingUIDNoETagNoIdempotency(t *testing.T) {
	rt := newETRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
	seedActiveParticipationAndStack(t, rt)

	recCreateB := doMux(t, rt.handler, http.MethodPost, etCollection,
		executionTargetCreateBody("target-b", testPartUID, testStackUID), createHeaders("et-list-b"))
	if recCreateB.Code != http.StatusCreated {
		t.Fatalf("create b=%d %s", recCreateB.Code, recCreateB.Body.String())
	}
	var createdB executiontarget.SafeExecutionTarget
	if err := json.Unmarshal(recCreateB.Body.Bytes(), &createdB); err != nil {
		t.Fatalf("decode b: %v", err)
	}

	// Second target needs a different stack for live tuple uniqueness.
	stack2, prob := rt.cloud.CreateInfrastructureStack(model.InfrastructureStack{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionInfrastructureStack, Kind: model.KindInfrastructureStack},
		Metadata: apimeta.ObjectMeta{
			Name: "stack-2", UID: "44444444444444444444444444444444", ResourceVersion: "1", Generation: 1,
			ScopeRef: providerScopeRef(testProviderUID, "prov-a"),
		},
		Status: model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed stack2: %#v", prob)
	}
	_ = doMux(t, rt.handler, http.MethodPost, etCollection,
		executionTargetCreateBody("target-a", testPartUID, stack2.Metadata.UID), createHeaders("et-list-a"))

	// Retire one target so LIST includes both Active and Retired (F16-14).
	retireNS := executiontarget.IdempotencyNamespace{
		PrincipalUID:     testPrincipal,
		RoutePattern:     "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire",
		CloudProviderUID: testProviderUID,
		ActionTargetUID:  createdB.Metadata.UID,
		Key:              "retire-list-b",
	}
	if res := rt.life.ReserveOrInspectReplay(retireNS, executiontarget.DigestAction()); res.Outcome != executiontarget.ReservationOutcomeReserved {
		t.Fatalf("reserve retire: %#v", res)
	}
	if _, retireProb := rt.life.CommitRetire(context.Background(), executiontarget.RetireCommitRequest{
		TargetUID:               createdB.Metadata.UID,
		ExpectedResourceVersion: createdB.Metadata.ResourceVersion,
		Actor:                   actorRef(testPrincipal),
		RequestID:               "req-retire-list",
		AuditUID:                "audit-retire-list",
		Idempotency:             retireNS,
		Digest:                  executiontarget.DigestAction(),
		Completion:              executiontarget.CompletedResult{StatusCode: http.StatusOK, Body: []byte(`{}`)},
	}); retireProb != nil {
		t.Fatalf("retire: %#v", retireProb)
	}

	before := rt.audit.Len()
	headers := map[string]string{"Idempotency-Key": "ignored-on-list", "Accept": "text/plain"}
	rec := doMux(t, rt.handler, http.MethodGet, etCollection, nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status=%d body=%s", rec.Code, rec.Body.String())
	}
	assertJSONNosniff(t, rec)
	if rec.Header().Get("ETag") != "" {
		t.Fatal("LIST must not return ETag")
	}
	if rt.audit.Len() != before {
		t.Fatal("LIST must not touch audit for success")
	}

	var items []executiontarget.SafeExecutionTarget
	if err := json.Unmarshal(rec.Body.Bytes(), &items); err != nil {
		t.Fatalf("list decode: %v body=%s", err, rec.Body.String())
	}
	if len(items) != 2 {
		t.Fatalf("list len=%d, want 2", len(items))
	}
	if items[0].Metadata.UID > items[1].Metadata.UID {
		t.Fatalf("LIST must be ascending UID: %q then %q", items[0].Metadata.UID, items[1].Metadata.UID)
	}
	var sawActive, sawRetired bool
	for _, item := range items {
		switch item.Status.Lifecycle {
		case etmodel.LifecycleActive:
			sawActive = true
		case etmodel.LifecycleRetired:
			sawRetired = true
		}
	}
	if !sawActive || !sawRetired {
		t.Fatalf("LIST must include Active and Retired; items=%#v", items)
	}
}

func TestExecutionTargetCollection_AuthAndSafeDenial(t *testing.T) {
	t.Run("missing auth", func(t *testing.T) {
		rt := newETRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
		rec := doMuxUnauth(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("t", testPartUID, testStackUID), createHeaders("k"))
		if rec.Code != http.StatusUnauthorized || rt.audit.Len() != 0 {
			t.Fatalf("want 401 no audit, got %d audit=%d", rec.Code, rt.audit.Len())
		}
		assertProblemNosniff(t, rec)
	})

	t.Run("create without write F16-40", func(t *testing.T) {
		rt := newETRuntime(executionTargetGrants(testPrincipal, testProviderUID, false, true))
		seedActiveParticipationAndStack(t, rt)
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("t", testPartUID, testStackUID), createHeaders("no-write"))
		if rec.Code != http.StatusForbidden {
			t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
		}
		assertProblemNosniff(t, rec)
		if rt.audit.Len() != 1 {
			t.Fatalf("audited denial audit=%d", rt.audit.Len())
		}
		targets := rt.life.Store().ListExecutionTargets()
		if len(targets) != 0 {
			t.Fatal("must not publish on authz denial")
		}
	})

	t.Run("create denial append failure F16-40", func(t *testing.T) {
		rt := newETRuntime(executionTargetGrants(testPrincipal, testProviderUID, false, true))
		seedActiveParticipationAndStack(t, rt)
		rt.audit.SetFail(errors.New("append boom"))
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("t", testPartUID, testStackUID), createHeaders("no-write-fail"))
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
		}
		p := decodeProblem(t, rec)
		if p.Code != apiproblem.CodeInternalError {
			t.Fatalf("code=%s", p.Code)
		}
		if rt.audit.Len() != 0 {
			t.Fatal("failed append must leave no durable AuditEvent")
		}
		if len(rt.life.Store().ListExecutionTargets()) != 0 {
			t.Fatal("append failure must not publish")
		}
	})

	t.Run("list without read F16-112 member", func(t *testing.T) {
		rt := newETRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, false))
		rec := doMux(t, rt.handler, http.MethodGet, etCollection, nil, nil)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("status=%d", rec.Code)
		}
		assertProblemNosniff(t, rec)
		if rt.audit.Len() != 1 {
			t.Fatalf("list denial audit=%d", rt.audit.Len())
		}
	})

	t.Run("list denial append failure", func(t *testing.T) {
		rt := newETRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, false))
		rt.audit.SetFail(errors.New("append boom"))
		rec := doMux(t, rt.handler, http.MethodGet, etCollection, nil, nil)
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
		}
		p := decodeProblem(t, rec)
		if p.Code != apiproblem.CodeInternalError {
			t.Fatalf("code=%s", p.Code)
		}
		if p.Status == http.StatusForbidden {
			t.Fatal("must not disclose suppressed 403")
		}
		if rt.audit.Len() != 0 {
			t.Fatal("failed append must leave no durable AuditEvent")
		}
	})

	t.Run("inaccessible backing", func(t *testing.T) {
		rt := newETRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
		seedPlatformAndProvider(t, &f15Runtime{store: rt.cloud}, testProviderUID)
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("t", testPartUID, testStackUID), createHeaders("missing-backing"))
		if rec.Code != http.StatusNotFound || !hasViolation(decodeProblem(t, rec), violationAuthorizationSafeDenial) {
			t.Fatalf("want safe 404, got %d %#v", rec.Code, decodeProblem(t, rec))
		}
		if rt.audit.Len() != 1 {
			t.Fatalf("safe denial audit=%d", rt.audit.Len())
		}
	})
}

func TestExecutionTargetCollection_ClosedFieldBoundaryAndGraph(t *testing.T) {
	rt := newETRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
	part, stack := seedActiveParticipationAndStack(t, rt)

	t.Run("unknown adapterAuthorityRef", func(t *testing.T) {
		body := []byte(`{
			"metadata":{"name":"bad"},
			"spec":{
				"cloudProviderParticipationRef":{"uid":"` + part.Metadata.UID + `"},
				"infrastructureStackRef":{"uid":"` + stack.Metadata.UID + `"},
				"targetClass":"synthetic-iaas",
				"adapterAuthorityRef":{"uid":"x"}
			}
		}`)
		rec := doMux(t, rt.handler, http.MethodPost, etCollection, body, createHeaders("unk"))
		if rec.Code != http.StatusBadRequest || decodeProblem(t, rec).Code != apiproblem.CodeUnknownField {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
	})

	t.Run("status write", func(t *testing.T) {
		body := []byte(`{
			"metadata":{"name":"bad-status"},
			"spec":{
				"cloudProviderParticipationRef":{"uid":"` + part.Metadata.UID + `"},
				"infrastructureStackRef":{"uid":"` + stack.Metadata.UID + `"},
				"targetClass":"synthetic-iaas"
			},
			"status":{"lifecycle":"Retired"}
		}`)
		before := rt.audit.Len()
		rec := doMux(t, rt.handler, http.MethodPost, etCollection, body, createHeaders("status"))
		if rec.Code != http.StatusForbidden || !hasViolation(decodeProblem(t, rec), validate.ViolationStatusFieldWrite) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
		if rt.audit.Len() != before+1 {
			t.Fatal("status write must be audited")
		}
	})

	t.Run("ineffective participation", func(t *testing.T) {
		// Suspend the live Active participation → not effective-Active.
		cur, ok := rt.cloud.GetParticipation(part.Metadata.UID)
		if !ok {
			t.Fatal("missing participation")
		}
		cur.Status.PlatformSuspended = true
		cur.Status.Phase = model.ParticipationPhaseSuspended
		if _, prob := rt.cloud.UpdateParticipation(cur); prob != nil {
			t.Fatalf("suspend participation: %#v", prob)
		}
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("ineff", part.Metadata.UID, stack.Metadata.UID), createHeaders("ineff"))
		if rec.Code != http.StatusConflict ||
			!hasViolation(decodeProblem(t, rec), executiontarget.ViolationExecutionTargetParticipationUnavailable) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
		// Restore for later subtests.
		cur.Status.PlatformSuspended = false
		cur.Status.Phase = model.ParticipationPhaseActive
		if _, prob := rt.cloud.UpdateParticipation(cur); prob != nil {
			t.Fatalf("restore participation: %#v", prob)
		}
	})

	t.Run("non-Active stack", func(t *testing.T) {
		inactive, prob := rt.cloud.CreateInfrastructureStack(model.InfrastructureStack{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionInfrastructureStack, Kind: model.KindInfrastructureStack},
			Metadata: apimeta.ObjectMeta{
				Name: "stack-inactive", UID: "acacacacacacacacacacacacacacacac", ResourceVersion: "1", Generation: 1,
				ScopeRef: providerScopeRef(testProviderUID, "prov-a"),
			},
			Status: model.InfrastructureStackStatus{Phase: ""},
		})
		if prob != nil {
			t.Fatalf("seed inactive: %#v", prob)
		}
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("inactive-stack", part.Metadata.UID, inactive.Metadata.UID), createHeaders("inactive-stack"))
		if rec.Code != http.StatusConflict ||
			!hasViolation(decodeProblem(t, rec), executiontarget.ViolationExecutionTargetStackUnavailable) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
	})

	t.Run("cross-scope graph", func(t *testing.T) {
		seedOtherProvider(t, &f15Runtime{store: rt.cloud}, testOtherProviderUID)
		otherStack, prob := rt.cloud.CreateInfrastructureStack(model.InfrastructureStack{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionInfrastructureStack, Kind: model.KindInfrastructureStack},
			Metadata: apimeta.ObjectMeta{
				Name: "other-stack", UID: "66666666666666666666666666666666", ResourceVersion: "1", Generation: 1,
				ScopeRef: providerScopeRef(testOtherProviderUID, "prov-b"),
			},
			Status: model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
		})
		if prob != nil {
			t.Fatalf("seed other stack: %#v", prob)
		}
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("cross", part.Metadata.UID, otherStack.Metadata.UID), createHeaders("cross"))
		if rec.Code != http.StatusUnprocessableEntity ||
			!hasViolation(decodeProblem(t, rec), executiontarget.ViolationExecutionTargetScopeMismatch) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		headers := createHeaders("dup-1")
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("dup-name", part.Metadata.UID, stack.Metadata.UID), headers)
		if rec.Code != http.StatusCreated {
			t.Fatalf("first create=%d %s", rec.Code, rec.Body.String())
		}
		stackB, prob := rt.cloud.CreateInfrastructureStack(model.InfrastructureStack{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionInfrastructureStack, Kind: model.KindInfrastructureStack},
			Metadata: apimeta.ObjectMeta{
				Name: "stack-b", UID: "77777777777777777777777777777777", ResourceVersion: "1", Generation: 1,
				ScopeRef: providerScopeRef(testProviderUID, "prov-a"),
			},
			Status: model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
		})
		if prob != nil {
			t.Fatalf("seed stack-b: %#v", prob)
		}
		rec = doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("dup-name", part.Metadata.UID, stackB.Metadata.UID), createHeaders("dup-2"))
		if rec.Code != http.StatusConflict || decodeProblem(t, rec).Code != apiproblem.CodeAlreadyExists {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
	})

	t.Run("live tuple duplicate", func(t *testing.T) {
		stackT, prob := rt.cloud.CreateInfrastructureStack(model.InfrastructureStack{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionInfrastructureStack, Kind: model.KindInfrastructureStack},
			Metadata: apimeta.ObjectMeta{
				Name: "stack-tuple", UID: "b1b1b1b1b1b1b1b1b1b1b1b1b1b1b1b1", ResourceVersion: "1", Generation: 1,
				ScopeRef: providerScopeRef(testProviderUID, "prov-a"),
			},
			Status: model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
		})
		if prob != nil {
			t.Fatalf("seed stack-tuple: %#v", prob)
		}
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("tuple-a", part.Metadata.UID, stackT.Metadata.UID), createHeaders("tuple-a"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("first tuple create=%d %s", rec.Code, rec.Body.String())
		}
		rec = doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("tuple-b", part.Metadata.UID, stackT.Metadata.UID), createHeaders("tuple-b"))
		if rec.Code != http.StatusConflict || decodeProblem(t, rec).Code != apiproblem.CodeAlreadyExists {
			t.Fatalf("live tuple got %d %#v", rec.Code, decodeProblem(t, rec))
		}
	})

	t.Run("metadata.uid rejected", func(t *testing.T) {
		body := []byte(`{
			"metadata":{"name":"uid-reject","uid":"cccccccccccccccccccccccccccccccc"},
			"spec":{
				"cloudProviderParticipationRef":{"uid":"` + part.Metadata.UID + `"},
				"infrastructureStackRef":{"uid":"` + stack.Metadata.UID + `"},
				"targetClass":"synthetic-iaas"
			}
		}`)
		rec := doMux(t, rt.handler, http.MethodPost, etCollection, body, createHeaders("uid-reject"))
		if rec.Code != http.StatusForbidden || !hasViolation(decodeProblem(t, rec), validate.ViolationSystemOwnedFieldWrite) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
	})

	t.Run("metadata.scopeRef rejected", func(t *testing.T) {
		body := []byte(`{
			"metadata":{"name":"scope-reject","scopeRef":{"uid":"` + testProviderUID + `"}},
			"spec":{
				"cloudProviderParticipationRef":{"uid":"` + part.Metadata.UID + `"},
				"infrastructureStackRef":{"uid":"` + stack.Metadata.UID + `"},
				"targetClass":"synthetic-iaas"
			}
		}`)
		rec := doMux(t, rt.handler, http.MethodPost, etCollection, body, createHeaders("scope-reject"))
		if rec.Code != http.StatusForbidden || !hasViolation(decodeProblem(t, rec), validate.ViolationSystemOwnedFieldWrite) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
	})

	t.Run("SecretRef extension rejected", func(t *testing.T) {
		body := []byte(`{
			"metadata":{"name":"secret-reject"},
			"spec":{
				"cloudProviderParticipationRef":{"uid":"` + part.Metadata.UID + `"},
				"infrastructureStackRef":{"uid":"` + stack.Metadata.UID + `"},
				"targetClass":"synthetic-iaas",
				"secretRef":{"name":"x"}
			}
		}`)
		rec := doMux(t, rt.handler, http.MethodPost, etCollection, body, createHeaders("secret-reject"))
		if rec.Code != http.StatusBadRequest || decodeProblem(t, rec).Code != apiproblem.CodeUnknownField {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
	})

	t.Run("digest reuse mismatch", func(t *testing.T) {
		headers := createHeaders("reuse-key")
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("reuse-a", part.Metadata.UID, stack.Metadata.UID), headers)
		// May already conflict on name/tuple from prior subtests; use unique stack.
		_ = rec
		stackC, prob := rt.cloud.CreateInfrastructureStack(model.InfrastructureStack{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionInfrastructureStack, Kind: model.KindInfrastructureStack},
			Metadata: apimeta.ObjectMeta{
				Name: "stack-c", UID: "88888888888888888888888888888888", ResourceVersion: "1", Generation: 1,
				ScopeRef: providerScopeRef(testProviderUID, "prov-a"),
			},
			Status: model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
		})
		if prob != nil {
			t.Fatalf("seed stack-c: %#v", prob)
		}
		rec = doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("reuse-a", part.Metadata.UID, stackC.Metadata.UID), headers)
		if rec.Code != http.StatusCreated {
			t.Fatalf("first=%d %s", rec.Code, rec.Body.String())
		}
		rec = doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("reuse-b", part.Metadata.UID, stackC.Metadata.UID), headers)
		if rec.Code != http.StatusConflict ||
			!hasViolation(decodeProblem(t, rec), cloudmodel.ViolationIdempotencyKeyReuseMismatch) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
	})
}

func TestExecutionTargetCollection_HTTPContractF16_45_48_105_106(t *testing.T) {
	rt := newETRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
	seedActiveParticipationAndStack(t, rt)

	t.Run("non-json media 415", func(t *testing.T) {
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("t", testPartUID, testStackUID), map[string]string{
				"Content-Type": "text/plain", "Idempotency-Key": "media",
			})
		if rec.Code != http.StatusUnsupportedMediaType {
			t.Fatalf("status=%d", rec.Code)
		}
		assertProblemNosniff(t, rec)
	})

	t.Run("missing idempotency key 400", func(t *testing.T) {
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("t", testPartUID, testStackUID), map[string]string{
				"Content-Type": "application/json",
			})
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status=%d", rec.Code)
		}
		assertProblemNosniff(t, rec)
	})

	t.Run("If-Match on create F16-105", func(t *testing.T) {
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("t", testPartUID, testStackUID), map[string]string{
				"Content-Type": "application/json", "Idempotency-Key": "ifmatch", "If-Match": "1",
			})
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status=%d %#v", rec.Code, decodeProblem(t, rec))
		}
	})

	t.Run("query rejected F16-106 collection members", func(t *testing.T) {
		rec := doMux(t, rt.handler, http.MethodPost, etCollection+"?x=1",
			executionTargetCreateBody("t", testPartUID, testStackUID), createHeaders("q-post"))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("POST query status=%d", rec.Code)
		}
		rec = doMux(t, rt.handler, http.MethodGet, etCollection+"?x=1", nil, nil)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("GET query status=%d", rec.Code)
		}
	})

	t.Run("Accept ignored success and problem F16-45/48", func(t *testing.T) {
		stackD, prob := rt.cloud.CreateInfrastructureStack(model.InfrastructureStack{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionInfrastructureStack, Kind: model.KindInfrastructureStack},
			Metadata: apimeta.ObjectMeta{
				Name: "stack-d", UID: "99999999999999999999999999999999", ResourceVersion: "1", Generation: 1,
				ScopeRef: providerScopeRef(testProviderUID, "prov-a"),
			},
			Status: model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
		})
		if prob != nil {
			t.Fatalf("seed: %#v", prob)
		}
		rec := doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("accept-ok", testPartUID, stackD.Metadata.UID), createHeaders("accept-ok"))
		if rec.Code != http.StatusCreated {
			t.Fatalf("status=%d", rec.Code)
		}
		assertJSONNosniff(t, rec)

		rec = doMux(t, rt.handler, http.MethodPost, etCollection,
			executionTargetCreateBody("accept-bad", "ffffffffffffffffffffffffffffffff", testStackUID), createHeaders("accept-bad"))
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status=%d", rec.Code)
		}
		assertProblemNosniff(t, rec)
	})
}

func TestExecutionTargetCollection_CreateOwnerCancelAndPanicCleanup(t *testing.T) {
	rt := newETRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
	seedActiveParticipationAndStack(t, rt)

	t.Run("owner cancellation", func(t *testing.T) {
		ns := executiontarget.IdempotencyNamespace{
			PrincipalUID:     testPrincipal,
			RoutePattern:     PatternExecutionTargetCreate,
			CloudProviderUID: testProviderUID,
			Key:              "cancel-key",
		}
		body := executionTargetCreateBody("cancel-t", testPartUID, testStackUID)
		digest, err := executiontarget.DigestCreate(json.RawMessage(body))
		if err != nil {
			t.Fatal(err)
		}
		if res := rt.life.ReserveOrInspectReplay(ns, digest); res.Outcome != executiontarget.ReservationOutcomeReserved {
			t.Fatalf("pre-reserve: %#v", res)
		}
		req := executiontarget.CreateCommitRequest{
			Name: "cancel-t", UID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			CloudProviderScopeUID: testProviderUID,
			ParticipationRef:      apimeta.TypedRef{UID: testPartUID},
			StackRef:              apimeta.TypedRef{UID: testStackUID},
			StackGeneration:       3,
			TargetClass:           etmodel.TargetClassSyntheticIaaS,
			Actor:                 actorRef(testPrincipal),
			RequestID:             "req-cancel",
			AuditUID:              "audit-cancel",
			Idempotency:           ns,
			Digest:                digest,
			Completion:            executiontarget.CompletedResult{StatusCode: 201, Body: []byte(`{}`)},
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, prob := rt.life.CommitCreate(ctx, req)
		if prob == nil || prob.Code != apiproblem.CodeInternalError {
			t.Fatalf("cancelled create: %#v", prob)
		}
		if len(rt.life.Store().ListExecutionTargets()) != 0 {
			t.Fatal("cancelled create must not publish")
		}
	})

	t.Run("owner panic cleanup wakes waiter", func(t *testing.T) {
		ns := executiontarget.IdempotencyNamespace{
			PrincipalUID:     testPrincipal,
			RoutePattern:     PatternExecutionTargetCreate,
			CloudProviderUID: testProviderUID,
			Key:              "panic-key",
		}
		stackP, prob := rt.cloud.CreateInfrastructureStack(model.InfrastructureStack{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionInfrastructureStack, Kind: model.KindInfrastructureStack},
			Metadata: apimeta.ObjectMeta{
				Name: "stack-panic", UID: "abababababababababababababababab", ResourceVersion: "1", Generation: 1,
				ScopeRef: providerScopeRef(testProviderUID, "prov-a"),
			},
			Status: model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
		})
		if prob != nil {
			t.Fatalf("seed: %#v", prob)
		}
		body := executionTargetCreateBody("panic-t", testPartUID, stackP.Metadata.UID)
		digest, err := executiontarget.DigestCreate(json.RawMessage(body))
		if err != nil {
			t.Fatal(err)
		}
		if res := rt.life.ReserveOrInspectReplay(ns, digest); res.Outcome != executiontarget.ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}

		done, ok := rt.life.AttachWaiter(ns)
		if !ok {
			t.Fatal("attach waiter")
		}
		var woke atomic.Bool
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-done
			woke.Store(true)
		}()

		// Simulate owner panic cleanup: detach under lock via AbortReservation,
		// then one post-unlock waiter signal (TASK-F16-04 protocol).
		rt.life.AbortReservation(ns)
		wg.Wait()
		if !woke.Load() {
			t.Fatal("waiter must wake after abort/panic cleanup")
		}
		if len(rt.life.Store().ListExecutionTargets()) != 0 {
			t.Fatal("panic/abort path must not publish")
		}
	})
}

func TestExecutionTargetCollection_RouteBuilderTwoCollectionRegistrations(t *testing.T) {
	rt := newETRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
	// Mirror the unexposed F0016 mux builder registration surface from
	// server.NewExecutionTargetMux (TASK-F16-09): exactly the two collection
	// patterns, wired to the real handler. Importing the server package from
	// this package would create an import cycle.
	mux := http.NewServeMux()
	mux.Handle("GET "+etCollection, rt.handler)
	mux.Handle("POST "+etCollection, rt.handler)

	if PatternExecutionTargetCreate != "POST "+etCollection {
		t.Fatalf("create pattern=%q", PatternExecutionTargetCreate)
	}

	// Source-prove routes.go owns the unexposed builder with exactly the two
	// collection patterns and no item/action placeholder registrations yet.
	routesSrc, err := os.ReadFile("../server/routes.go")
	if err != nil {
		t.Fatalf("read routes.go: %v", err)
	}
	src := string(routesSrc)
	for _, want := range []string{
		"func NewExecutionTargetMux(",
		`routeExecutionTargetList   = "GET /apis/execution.sovrunn.io/v1alpha1/execution-targets"`,
		`routeExecutionTargetCreate = "POST /apis/execution.sovrunn.io/v1alpha1/execution-targets"`,
		"registerExecutionTargetPattern(mux, routeExecutionTargetList, h.Collection)",
		"registerExecutionTargetPattern(mux, routeExecutionTargetCreate, h.Collection)",
	} {
		if !strings.Contains(src, want) {
			t.Fatalf("routes.go missing %q", want)
		}
	}
	for _, forbidden := range []string{
		`"/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}"`,
		`"/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify"`,
		`"/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire"`,
		"StatusNotImplemented",
		"http.StatusNotImplemented",
	} {
		if strings.Contains(src, forbidden) {
			t.Fatalf("routes.go must not contain %q until TASK-F16-10/11", forbidden)
		}
	}

	seedActiveParticipationAndStack(t, rt)
	rec := doMux(t, mux, http.MethodPost, etCollection,
		executionTargetCreateBody("routed", testPartUID, testStackUID), createHeaders("route-post"))
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST via mux status=%d body=%s", rec.Code, rec.Body.String())
	}
	rec = doMux(t, mux, http.MethodGet, etCollection, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET via mux status=%d", rec.Code)
	}

	// No item/qualify/retire registrations yet (TASK-F16-10/11).
	rec = doMux(t, mux, http.MethodGet, etCollection+"/dddddddddddddddddddddddddddddddd", nil, nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("item must be unregistered: %d", rec.Code)
	}
}
