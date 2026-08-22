package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget"
	etmodel "github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

type etActionRuntime struct {
	*etRuntime
	fixtures *executiontarget.MapFixtureSource
	qualify  *ExecutionTargetActionHandler
	retire   *ExecutionTargetActionHandler
	item     *ExecutionTargetItemHandler
}

func newETActionRuntime(grants GrantResolver) *etActionRuntime {
	cloud := cloudmodel.NewStore()
	audit := &cloudmodel.MemoryAuditAppender{}
	fixtures := executiontarget.NewMapFixtureSource()
	life := executiontarget.NewExecutionTargetLifecycleService(executiontarget.LifecycleConfig{
		Audit:    audit,
		Observer: executiontarget.NewSyntheticObserver(fixtures, nil),
	})
	base := &etRuntime{cloud: cloud, audit: audit, life: life, grants: grants,
		handler: NewExecutionTargetCollectionHandler(life, cloud, grants, audit)}
	qualify := NewExecutionTargetQualifyHandler(life, cloud, grants, audit)
	qualify.Fixtures = fixtures
	return &etActionRuntime{
		etRuntime: base,
		fixtures:  fixtures,
		qualify:   qualify,
		retire:    NewExecutionTargetRetireHandler(life, cloud, grants, audit),
		item:      NewExecutionTargetItemHandler(life, cloud, grants, audit),
	}
}

// newETActionRuntimeWithObserver builds an action runtime whose lifecycle
// observer may block (Retire/Maintenance-wins HTTP re-tests). Fixtures remains
// the non-blocking MapFixtureSource used for qualify completion prediction.
func newETActionRuntimeWithObserver(grants GrantResolver, observer *executiontarget.SyntheticObserver, fixtures *executiontarget.MapFixtureSource) *etActionRuntime {
	cloud := cloudmodel.NewStore()
	audit := &cloudmodel.MemoryAuditAppender{}
	if fixtures == nil {
		fixtures = executiontarget.NewMapFixtureSource()
	}
	if observer == nil {
		observer = executiontarget.NewSyntheticObserver(fixtures, nil)
	}
	life := executiontarget.NewExecutionTargetLifecycleService(executiontarget.LifecycleConfig{
		Audit:    audit,
		Observer: observer,
	})
	base := &etRuntime{cloud: cloud, audit: audit, life: life, grants: grants,
		handler: NewExecutionTargetCollectionHandler(life, cloud, grants, audit)}
	qualify := NewExecutionTargetQualifyHandler(life, cloud, grants, audit)
	qualify.Fixtures = fixtures
	return &etActionRuntime{
		etRuntime: base,
		fixtures:  fixtures,
		qualify:   qualify,
		retire:    NewExecutionTargetRetireHandler(life, cloud, grants, audit),
		item:      NewExecutionTargetItemHandler(life, cloud, grants, audit),
	}
}

func executionTargetActionGrants(principal, providerUID string, write, read, qualify bool) *staticGrants {
	g := executionTargetGrants(principal, providerUID, write, read)
	if qualify {
		scope := apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: providerUID}
		g.grants = append(g.grants, cloudmodel.Grant{Action: ActionExecutionTargetQualify, Scope: scope})
	}
	return g
}

func muxActions(rt *etActionRuntime) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET "+etCollection, rt.handler)
	mux.Handle("POST "+etCollection, rt.handler)
	mux.Handle(PatternExecutionTargetItemGet, rt.item)
	mux.Handle(PatternExecutionTargetQualify, rt.qualify)
	mux.Handle(PatternExecutionTargetRetire, rt.retire)
	return mux
}

func etActionPath(uid, action string) string {
	return etItemBase + uid + "/actions/" + action
}

func etActionHeaders(key, ifMatch string) map[string]string {
	return map[string]string{
		"Content-Type":    "application/json",
		"Idempotency-Key": key,
		"If-Match":        ifMatch,
		"Accept":          "application/xml",
	}
}

func createTargetForAction(t *testing.T, rt *etActionRuntime, name, key string) executiontarget.SafeExecutionTarget {
	t.Helper()
	seedActiveParticipationAndStack(t, rt.etRuntime)
	rec := doMux(t, muxActions(rt), http.MethodPost, etCollection,
		executionTargetCreateBody(name, testPartUID, testStackUID), createHeaders(key))
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created executiontarget.SafeExecutionTarget
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	return created
}

func TestExecutionTargetActions_QualifyHappyPathReplayAndETag(t *testing.T) {
	rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
	created := createTargetForAction(t, rt, "qualify-ok", "create-q-ok")
	before := rt.audit.Len()

	headers := etActionHeaders("qualify-1", created.Metadata.ResourceVersion)
	rec := doMux(t, muxActions(rt), http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("qualify status=%d body=%s", rec.Code, rec.Body.String())
	}
	assertJSONNosniff(t, rec)
	if rec.Header().Get("ETag") == "" || rec.Header().Get("ETag") == created.Metadata.ResourceVersion {
		t.Fatalf("ETag must advance; got %q", rec.Header().Get("ETag"))
	}
	if rec.Header().Get("Idempotency-Key") != "" {
		t.Fatal("must not echo Idempotency-Key")
	}
	if rt.audit.Len() != before+1 {
		t.Fatalf("audit len=%d want %d", rt.audit.Len(), before+1)
	}

	var proj executiontarget.SafeExecutionTarget
	if err := json.Unmarshal(rec.Body.Bytes(), &proj); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if proj.Status.Lifecycle != etmodel.LifecycleActive {
		t.Fatalf("lifecycle=%s", proj.Status.Lifecycle)
	}
	if proj.Status.Qualification != etmodel.QualificationIndeterminate {
		t.Fatalf("qualification=%s want Indeterminate", proj.Status.Qualification)
	}
	raw := rec.Body.String()
	for _, forbidden := range []string{"factSetRef", "qualificationResultRef", "observerID", "conditions"} {
		if strings.Contains(raw, forbidden) {
			t.Fatalf("safe projection must omit %q", forbidden)
		}
	}

	rec2 := doMux(t, muxActions(rt), http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil, headers)
	if rec2.Code != http.StatusOK {
		t.Fatalf("replay status=%d", rec2.Code)
	}
	if rec2.Body.String() != rec.Body.String() {
		t.Fatal("replay body must match stored successful body")
	}
	if rt.audit.Len() != before+1 {
		t.Fatal("replay must not append a second AuditEvent")
	}
}

func TestExecutionTargetActions_RetireHappyPathClearsRefs(t *testing.T) {
	rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
	created := createTargetForAction(t, rt, "retire-ok", "create-r-ok")
	before := rt.audit.Len()

	headers := etActionHeaders("retire-1", created.Metadata.ResourceVersion)
	rec := doMux(t, muxActions(rt), http.MethodPost, etActionPath(created.Metadata.UID, "retire"), nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("retire status=%d body=%s", rec.Code, rec.Body.String())
	}
	assertJSONNosniff(t, rec)
	var proj executiontarget.SafeExecutionTarget
	if err := json.Unmarshal(rec.Body.Bytes(), &proj); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if proj.Status.Lifecycle != etmodel.LifecycleRetired ||
		proj.Status.Qualification != etmodel.QualificationUnqualified ||
		proj.EffectiveAvailability != etmodel.EffectiveUnavailable {
		t.Fatalf("projection=%#v", proj)
	}
	if rt.audit.Len() != before+1 {
		t.Fatal("retirement must append exactly one AuditEvent")
	}
	stored, ok := rt.life.Store().GetExecutionTarget(created.Metadata.UID)
	if !ok || stored.Status.FactSetRef != nil || stored.Status.QualificationResultRef != nil {
		t.Fatal("retire must clear record refs")
	}
}

func TestExecutionTargetActions_GrammarAndHeaderForm(t *testing.T) {
	rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
	created := createTargetForAction(t, rt, "grammar", "create-grammar")
	path := etActionPath(created.Metadata.UID, "qualify")
	mux := muxActions(rt)
	before := rt.audit.Len()

	t.Run("missing If-Match F16-24", func(t *testing.T) {
		h := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "no-if"}
		rec := doMux(t, mux, http.MethodPost, path, nil, h)
		if rec.Code != http.StatusPreconditionFailed || decodeProblem(t, rec).Code != apiproblem.CodeStaleResourceVersion {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
		assertProblemNosniff(t, rec)
	})
	t.Run("missing Idempotency-Key", func(t *testing.T) {
		h := map[string]string{"Content-Type": "application/json", "If-Match": created.Metadata.ResourceVersion}
		rec := doMux(t, mux, http.MethodPost, path, nil, h)
		if rec.Code != http.StatusBadRequest || decodeProblem(t, rec).Code != apiproblem.CodeMalformedRequest {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
	})
	t.Run("non-zero body F16-25", func(t *testing.T) {
		h := etActionHeaders("body", created.Metadata.ResourceVersion)
		rec := doMux(t, mux, http.MethodPost, path, []byte(`{}`), h)
		if rec.Code != http.StatusBadRequest || decodeProblem(t, rec).Code != apiproblem.CodeMalformedRequest {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
	})
	t.Run("weak If-Match F16-71", func(t *testing.T) {
		h := etActionHeaders("weak", `W/"1"`)
		rec := doMux(t, mux, http.MethodPost, path, nil, h)
		if rec.Code != http.StatusPreconditionFailed {
			t.Fatalf("got %d", rec.Code)
		}
	})
	t.Run("wildcard If-Match F16-72", func(t *testing.T) {
		h := etActionHeaders("wild", "*")
		rec := doMux(t, mux, http.MethodPost, path, nil, h)
		if rec.Code != http.StatusPreconditionFailed {
			t.Fatalf("got %d", rec.Code)
		}
	})
	t.Run("multiple If-Match F16-73", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, path, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", "multi")
		req.Header.Add("If-Match", "1")
		req.Header.Add("If-Match", "2")
		req = withAuth(req, testPrincipal)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusPreconditionFailed {
			t.Fatalf("got %d", rec.Code)
		}
	})
	t.Run("malformed If-Match before body F16-86", func(t *testing.T) {
		h := etActionHeaders("malformed", "bad\x01")
		rec := doMux(t, mux, http.MethodPost, path, []byte(`{"x":1}`), h)
		if rec.Code != http.StatusPreconditionFailed {
			t.Fatalf("header-form must precede body; got %d", rec.Code)
		}
	})
	t.Run("stale If-Match F16-74", func(t *testing.T) {
		h := etActionHeaders("stale", "999")
		rec := doMux(t, mux, http.MethodPost, path, nil, h)
		if rec.Code != http.StatusPreconditionFailed || decodeProblem(t, rec).Code != apiproblem.CodeStaleResourceVersion {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
	})
	t.Run("query rejected F16-109", func(t *testing.T) {
		h := etActionHeaders("query", created.Metadata.ResourceVersion)
		rec := doMux(t, mux, http.MethodPost, path+"?x=1", nil, h)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d", rec.Code)
		}
	})
	t.Run("malformed UID F16-110", func(t *testing.T) {
		h := etActionHeaders("baduid", "1")
		rec := doMux(t, mux, http.MethodPost, etActionPath("not-a-uid", "qualify"), nil, h)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("got %d", rec.Code)
		}
	})
	if rt.audit.Len() != before {
		t.Fatal("grammar failures must not audit")
	}
}

func TestExecutionTargetActions_AuthorizationAndSafeDenial(t *testing.T) {
	t.Run("qualify without grant", func(t *testing.T) {
		rtWrite := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, false))
		created := createTargetForAction(t, rtWrite, "no-qual", "create-no-qual")
		before := rtWrite.audit.Len()
		rec := doMux(t, muxActions(rtWrite), http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil,
			etActionHeaders("q-deny", created.Metadata.ResourceVersion))
		if rec.Code != http.StatusForbidden || decodeProblem(t, rec).Code != apiproblem.CodeAuthorizationDenied {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
		if rtWrite.audit.Len() != before+1 {
			t.Fatal("authorization denial must append exactly one AuditEvent")
		}
	})
	t.Run("inaccessible action safe 404", func(t *testing.T) {
		rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
		seedActiveParticipationAndStack(t, rt.etRuntime)
		before := rt.audit.Len()
		missing := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		rec := doMux(t, muxActions(rt), http.MethodPost, etActionPath(missing, "qualify"), nil,
			etActionHeaders("missing", "1"))
		if rec.Code != http.StatusNotFound || !hasViolation(decodeProblem(t, rec), violationAuthorizationSafeDenial) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
		if rt.audit.Len() != before+1 {
			t.Fatal("safe denial must append exactly one AuditEvent")
		}
	})
	t.Run("append failure F16-112 returns only 500", func(t *testing.T) {
		rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
		seedActiveParticipationAndStack(t, rt.etRuntime)
		rt.audit.SetFail(errors.New("append boom"))
		missing := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
		rec := doMux(t, muxActions(rt), http.MethodPost, etActionPath(missing, "retire"), nil,
			etActionHeaders("fail-append", "1"))
		p := decodeProblem(t, rec)
		if rec.Code != http.StatusInternalServerError || p.Code != apiproblem.CodeInternalError {
			t.Fatalf("got %d %#v", rec.Code, p)
		}
		if hasViolation(p, violationAuthorizationSafeDenial) {
			t.Fatal("must not disclose suppressed 404")
		}
		if rt.audit.Len() != 0 {
			t.Fatal("failed append must leave no durable AuditEvent")
		}
	})
}

func TestExecutionTargetActions_QualifyBackingAdmission(t *testing.T) {
	rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
	created := createTargetForAction(t, rt, "admit", "create-admit")
	mux := muxActions(rt)
	path := etActionPath(created.Metadata.UID, "qualify")

	t.Run("ineffective participation unaudited 409", func(t *testing.T) {
		part, _ := rt.cloud.GetParticipation(testPartUID)
		part.Status.Phase = model.ParticipationPhaseSuspended
		rt.cloud.BeginPublication()
		staged, _ := rt.cloud.StageUpdateParticipation(part)
		rt.cloud.Publish(staged)
		rt.cloud.EndPublication()
		before := rt.audit.Len()
		rec := doMux(t, mux, http.MethodPost, path, nil, etActionHeaders("part-unavail", created.Metadata.ResourceVersion))
		if rec.Code != http.StatusConflict ||
			!hasViolation(decodeProblem(t, rec), executiontarget.ViolationExecutionTargetParticipationUnavailable) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
		if rt.audit.Len() != before {
			t.Fatal("participation denial must not audit")
		}
		part.Status.Phase = model.ParticipationPhaseActive
		rt.cloud.BeginPublication()
		staged, _ = rt.cloud.StageUpdateParticipation(part)
		rt.cloud.Publish(staged)
		rt.cloud.EndPublication()
	})

	t.Run("inactive stack unaudited 409", func(t *testing.T) {
		stack, _ := rt.cloud.GetInfrastructureStack(testStackUID)
		stack.Status.Phase = model.InfrastructureStackPhase("Failed")
		rt.cloud.BeginPublication()
		staged, _ := rt.cloud.StageUpdateInfrastructureStack(stack)
		rt.cloud.Publish(staged)
		rt.cloud.EndPublication()
		before := rt.audit.Len()
		rec := doMux(t, mux, http.MethodPost, path, nil, etActionHeaders("stack-unavail", created.Metadata.ResourceVersion))
		if rec.Code != http.StatusConflict ||
			!hasViolation(decodeProblem(t, rec), executiontarget.ViolationExecutionTargetStackUnavailable) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
		if rt.audit.Len() != before {
			t.Fatal("stack denial must not audit")
		}
		stack.Status.Phase = model.InfrastructureStackPhaseActive
		rt.cloud.BeginPublication()
		staged, _ = rt.cloud.StageUpdateInfrastructureStack(stack)
		rt.cloud.Publish(staged)
		rt.cloud.EndPublication()
	})

	t.Run("scope mismatch unaudited 422", func(t *testing.T) {
		stack, _ := rt.cloud.GetInfrastructureStack(testStackUID)
		stack.Metadata.ScopeRef = &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
			APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider,
			Name: "other", UID: "dddddddddddddddddddddddddddddddd",
		}}
		rt.cloud.BeginPublication()
		staged, _ := rt.cloud.StageUpdateInfrastructureStack(stack)
		rt.cloud.Publish(staged)
		rt.cloud.EndPublication()
		before := rt.audit.Len()
		rec := doMux(t, mux, http.MethodPost, path, nil, etActionHeaders("scope-mismatch", created.Metadata.ResourceVersion))
		if rec.Code != http.StatusUnprocessableEntity ||
			!hasViolation(decodeProblem(t, rec), executiontarget.ViolationExecutionTargetScopeMismatch) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
		if rt.audit.Len() != before {
			t.Fatal("scope mismatch must not audit")
		}
		stack.Metadata.ScopeRef = providerScopeRef(testProviderUID, "prov-a")
		rt.cloud.BeginPublication()
		staged, _ = rt.cloud.StageUpdateInfrastructureStack(stack)
		rt.cloud.Publish(staged)
		rt.cloud.EndPublication()
	})

	t.Run("inaccessible backing audited 404", func(t *testing.T) {
		rt2 := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
		seedActiveParticipationAndStack(t, rt2.etRuntime)
		missingPart := "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
		missingStack := "ffffffffffffffffffffffffffffffff"
		uid := "cccccccccccccccccccccccccccccccc"
		ns := executiontarget.IdempotencyNamespace{
			PrincipalUID:     testPrincipal,
			RoutePattern:     PatternExecutionTargetCreate,
			CloudProviderUID: testProviderUID,
			Key:              "seed-orphan-backing",
		}
		if res := rt2.life.ReserveOrInspectReplay(ns, executiontarget.DigestAction()); res.Outcome != executiontarget.ReservationOutcomeReserved {
			t.Fatalf("reserve: %#v", res)
		}
		_, prob := rt2.life.CommitCreate(context.Background(), executiontarget.CreateCommitRequest{
			Name:                  "orphan-backing",
			UID:                   uid,
			CloudProviderScopeUID: testProviderUID,
			ParticipationRef:      apimeta.TypedRef{UID: missingPart},
			StackRef:              apimeta.TypedRef{UID: missingStack},
			StackGeneration:       1,
			Actor:                 actorRef(testPrincipal),
			RequestID:             "req-orphan",
			AuditUID:              "audit-orphan",
			Idempotency:           ns,
			Digest:                executiontarget.DigestAction(),
			Completion:            executiontarget.CompletedResult{StatusCode: http.StatusCreated, Body: []byte(`{}`)},
		})
		if prob != nil {
			t.Fatalf("seed create: %#v", prob)
		}
		before := rt2.audit.Len()
		rec := doMux(t, muxActions(rt2), http.MethodPost, etActionPath(uid, "qualify"), nil,
			etActionHeaders("backing-miss", "1"))
		if rec.Code != http.StatusNotFound || !hasViolation(decodeProblem(t, rec), violationAuthorizationSafeDenial) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
		if rt2.audit.Len() != before+1 {
			t.Fatal("inaccessible backing must append exactly one AuditEvent")
		}
	})
}

func TestExecutionTargetActions_RetiredQualifyAndRetireAgain(t *testing.T) {
	rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
	created := createTargetForAction(t, rt, "retired", "create-retired")
	mux := muxActions(rt)

	rec := doMux(t, mux, http.MethodPost, etActionPath(created.Metadata.UID, "retire"), nil,
		etActionHeaders("retire-first", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusOK {
		t.Fatalf("retire: %d", rec.Code)
	}
	var retired executiontarget.SafeExecutionTarget
	_ = json.Unmarshal(rec.Body.Bytes(), &retired)

	q := doMux(t, mux, http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil,
		etActionHeaders("q-retired", retired.Metadata.ResourceVersion))
	if q.Code != http.StatusConflict || !hasViolation(decodeProblem(t, q), executiontarget.ViolationTargetRetired) {
		t.Fatalf("qualify retired: %d %#v", q.Code, decodeProblem(t, q))
	}

	r2 := doMux(t, mux, http.MethodPost, etActionPath(created.Metadata.UID, "retire"), nil,
		etActionHeaders("retire-again", retired.Metadata.ResourceVersion))
	if r2.Code != http.StatusConflict || !hasViolation(decodeProblem(t, r2), executiontarget.ViolationTargetRetired) {
		t.Fatalf("retire again: %d %#v", r2.Code, decodeProblem(t, r2))
	}
}

func TestExecutionTargetActions_RetireWinsOverInFlightQualify(t *testing.T) {
	fixtures := executiontarget.NewMapFixtureSource()
	started := make(chan struct{})
	release := make(chan struct{})
	blocking := &blockingActionFixture{inner: fixtures, started: started, release: release}
	rt := newETActionRuntimeWithObserver(
		executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true),
		executiontarget.NewSyntheticObserver(blocking, nil),
		fixtures,
	)
	created := createTargetForAction(t, rt, "retire-wins", "create-rw")
	mux := muxActions(rt)
	path := etActionPath(created.Metadata.UID, "qualify")
	beforeAudit := rt.audit.Len()

	qualDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		qualDone <- doMux(t, mux, http.MethodPost, path, nil,
			etActionHeaders("q-inflight", created.Metadata.ResourceVersion))
	}()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("qualify did not reach observation")
	}

	retireRec := doMux(t, mux, http.MethodPost, etActionPath(created.Metadata.UID, "retire"), nil,
		etActionHeaders("retire-win", created.Metadata.ResourceVersion))
	if retireRec.Code != http.StatusOK {
		t.Fatalf("retire: %d %s", retireRec.Code, retireRec.Body.String())
	}
	close(release)

	var qRec *httptest.ResponseRecorder
	select {
	case qRec = <-qualDone:
	case <-time.After(3 * time.Second):
		t.Fatal("qualify did not finish after retire")
	}
	if qRec.Code != http.StatusConflict || !hasViolation(decodeProblem(t, qRec), executiontarget.ViolationTargetRetired) {
		t.Fatalf("qualifying caller: %d %#v", qRec.Code, decodeProblem(t, qRec))
	}
	// Winning retirement audited; no qualification AuditEvent on top of create+retire.
	if rt.audit.Len() != beforeAudit+1 {
		t.Fatalf("audit len=%d want create+retire only (%d)", rt.audit.Len(), beforeAudit+1)
	}
}

func TestExecutionTargetActions_MaintenanceWinsOverInFlightQualify(t *testing.T) {
	fixtures := executiontarget.NewMapFixtureSource()
	started := make(chan struct{})
	release := make(chan struct{})
	blocking := &blockingActionFixture{inner: fixtures, started: started, release: release}
	rt := newETActionRuntimeWithObserver(
		executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true),
		executiontarget.NewSyntheticObserver(blocking, nil),
		fixtures,
	)
	created := createTargetForAction(t, rt, "maint-wins", "create-mw")
	mux := muxActions(rt)
	path := etActionPath(created.Metadata.UID, "qualify")
	beforeAudit := rt.audit.Len()

	qualDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		qualDone <- doMux(t, mux, http.MethodPost, path, nil,
			etActionHeaders("q-maint-inflight", created.Metadata.ResourceVersion))
	}()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("qualify did not reach observation")
	}

	_, prob := rt.life.CommitMaintenanceEnter(context.Background(), executiontarget.MaintenanceCommitRequest{
		TargetUID:                      created.Metadata.UID,
		ExpectedInfrastructureStackGen: 3,
		ExpectedMaintenanceEpoch:       0,
		Actor:                          actorRef(testPrincipal),
		RequestID:                      "maint-wins",
		AuditUID:                       "audit-maint-wins",
	})
	if prob != nil {
		t.Fatalf("maintenance enter: %#v", prob)
	}
	close(release)

	var qRec *httptest.ResponseRecorder
	select {
	case qRec = <-qualDone:
	case <-time.After(3 * time.Second):
		t.Fatal("qualify did not finish after maintenance")
	}
	if qRec.Code != http.StatusConflict || !hasViolation(decodeProblem(t, qRec), executiontarget.ViolationTargetMaintenance) {
		t.Fatalf("qualifying caller: %d %#v", qRec.Code, decodeProblem(t, qRec))
	}
	if rt.audit.Len() != beforeAudit+1 {
		t.Fatalf("audit len=%d want create+maintenance only", rt.audit.Len())
	}
}

// blockingActionFixture blocks Lookup until release is closed so HTTP qualify
// can hold observation open for Retire/Maintenance-wins re-tests.
type blockingActionFixture struct {
	inner   executiontarget.FixtureSource
	started chan struct{}
	release chan struct{}
}

func (b *blockingActionFixture) Lookup(targetUID string) (executiontarget.TargetBoundFixture, bool) {
	var (
		fx executiontarget.TargetBoundFixture
		ok bool
	)
	if b.inner != nil {
		fx, ok = b.inner.Lookup(targetUID)
	}
	select {
	case <-b.started:
	default:
		close(b.started)
	}
	<-b.release
	return fx, ok
}

func TestExecutionTargetActions_MaintenanceDeniesQualifyOnly(t *testing.T) {
	rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
	created := createTargetForAction(t, rt, "maint", "create-maint")
	mux := muxActions(rt)

	_, prob := rt.life.CommitMaintenanceEnter(context.Background(), executiontarget.MaintenanceCommitRequest{
		TargetUID:                      created.Metadata.UID,
		ExpectedInfrastructureStackGen: 3,
		ExpectedMaintenanceEpoch:       0,
		Actor:                          actorRef(testPrincipal),
		RequestID:                      "maint-enter",
		AuditUID:                       "audit-maint",
	})
	if prob != nil {
		t.Fatalf("maintenance enter: %#v", prob)
	}
	et, _ := rt.life.Store().GetExecutionTarget(created.Metadata.UID)
	rv := et.Metadata.ResourceVersion

	q := doMux(t, mux, http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil,
		etActionHeaders("q-maint", rv))
	if q.Code != http.StatusConflict || !hasViolation(decodeProblem(t, q), executiontarget.ViolationTargetMaintenance) {
		t.Fatalf("qualify during maintenance: %d %#v", q.Code, decodeProblem(t, q))
	}

	r := doMux(t, mux, http.MethodPost, etActionPath(created.Metadata.UID, "retire"), nil,
		etActionHeaders("r-maint", rv))
	if r.Code != http.StatusOK {
		t.Fatalf("retire during maintenance: %d %s", r.Code, r.Body.String())
	}
}

func TestExecutionTargetActions_StaleIfMatchBeforeMaintenance(t *testing.T) {
	rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
	created := createTargetForAction(t, rt, "stale-maint", "create-sm")
	_, prob := rt.life.CommitMaintenanceEnter(context.Background(), executiontarget.MaintenanceCommitRequest{
		TargetUID:                      created.Metadata.UID,
		ExpectedInfrastructureStackGen: 3,
		ExpectedMaintenanceEpoch:       0,
		Actor:                          actorRef(testPrincipal),
		RequestID:                      "maint-stale",
		AuditUID:                       "audit-sm",
	})
	if prob != nil {
		t.Fatalf("enter: %#v", prob)
	}
	rec := doMux(t, muxActions(rt), http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil,
		etActionHeaders("stale-before-maint", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusPreconditionFailed || decodeProblem(t, rec).Code != apiproblem.CodeStaleResourceVersion {
		t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
	}
}

func TestExecutionTargetActions_PreCommitAbortWakesWaiter(t *testing.T) {
	rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
	created := createTargetForAction(t, rt, "waiter", "create-waiter")
	mux := muxActions(rt)
	path := etActionPath(created.Metadata.UID, "qualify")

	ns := executiontarget.IdempotencyNamespace{
		PrincipalUID:     testPrincipal,
		RoutePattern:     PatternExecutionTargetQualify,
		CloudProviderUID: testProviderUID,
		ActionTargetUID:  created.Metadata.UID,
		Key:              "shared-key",
	}
	if res := rt.life.ReserveOrInspectReplay(ns, executiontarget.DigestAction()); res.Outcome != executiontarget.ReservationOutcomeReserved {
		t.Fatalf("owner reserve: %#v", res)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var waiterCode int
	go func() {
		defer wg.Done()
		rec := doMux(t, mux, http.MethodPost, path, nil, etActionHeaders("shared-key", created.Metadata.ResourceVersion))
		waiterCode = rec.Code
	}()
	time.Sleep(30 * time.Millisecond)

	if !rt.life.AbortReservation(ns) {
		t.Fatal("abort owner reservation")
	}
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		if waiterCode == 0 {
			t.Fatal("empty status")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("waiter did not wake after abort")
	}
}

func TestExecutionTargetActions_ProblemIgnoresAccept(t *testing.T) {
	rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
	created := createTargetForAction(t, rt, "accept", "create-accept")
	h := etActionHeaders("accept-prob", "stale-rv")
	h["Accept"] = "text/plain"
	rec := doMux(t, muxActions(rt), http.MethodPost, etActionPath(created.Metadata.UID, "retire"), nil, h)
	assertProblemNosniff(t, rec)
}

func TestExecutionTargetActions_ReplayDenialAppendFailures_F16_103_104(t *testing.T) {
	t.Run("F16-103 same-key replay authz denial append fails", func(t *testing.T) {
		grants := executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true)
		rt := newETActionRuntime(grants)
		created := createTargetForAction(t, rt, "f103", "create-f103")
		headers := etActionHeaders("replay-authz", created.Metadata.ResourceVersion)
		rec := doMux(t, muxActions(rt), http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil, headers)
		if rec.Code != http.StatusOK {
			t.Fatalf("qualify: %d %s", rec.Code, rec.Body.String())
		}
		// Revoke qualify grant; keep principal authentication.
		*grants = *executionTargetActionGrants(testPrincipal, testProviderUID, true, true, false)
		rt.audit.SetFail(errors.New("append boom"))
		rec2 := doMux(t, muxActions(rt), http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil, headers)
		p := decodeProblem(t, rec2)
		if rec2.Code != http.StatusInternalServerError || p.Code != apiproblem.CodeInternalError {
			t.Fatalf("got %d %#v", rec2.Code, p)
		}
		if p.Status == http.StatusForbidden || hasViolation(p, apiproblem.ViolationCode("VS0_AUTHORIZATION_SAFE_DENIAL")) {
			t.Fatal("must not disclose suppressed 403")
		}
	})

	t.Run("F16-104 same-key replay safe-denial append fails", func(t *testing.T) {
		rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
		created := createTargetForAction(t, rt, "f104", "create-f104")
		headers := etActionHeaders("replay-safe", created.Metadata.ResourceVersion)
		rec := doMux(t, muxActions(rt), http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil, headers)
		if rec.Code != http.StatusOK {
			t.Fatalf("qualify: %d", rec.Code)
		}
		// Empty cloud store makes backing inaccessible while the completed
		// same-key record remains; admission safe-denies before replay.
		rt.qualify.Cloud = cloudmodel.NewStore()
		rt.audit.SetFail(errors.New("append boom"))
		rec2 := doMux(t, muxActions(rt), http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil, headers)
		p := decodeProblem(t, rec2)
		if rec2.Code != http.StatusInternalServerError || p.Code != apiproblem.CodeInternalError {
			t.Fatalf("got %d %#v", rec2.Code, p)
		}
		if p.Status == http.StatusNotFound || hasViolation(p, violationAuthorizationSafeDenial) {
			t.Fatal("must not disclose suppressed 404")
		}
	})
}

func TestExecutionTargetActions_DigestConflictAndQualifyInProgress(t *testing.T) {
	t.Run("changed digest conflict", func(t *testing.T) {
		rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
		created := createTargetForAction(t, rt, "digest", "create-digest")
		ns := executiontarget.IdempotencyNamespace{
			PrincipalUID:     testPrincipal,
			RoutePattern:     PatternExecutionTargetQualify,
			CloudProviderUID: testProviderUID,
			ActionTargetUID:  created.Metadata.UID,
			Key:              "digest-key",
		}
		var other executiontarget.Digest
		other[0] = 0xff
		if res := rt.life.ReserveOrInspectReplay(ns, other); res.Outcome != executiontarget.ReservationOutcomeReserved {
			t.Fatalf("seed other digest: %#v", res)
		}
		rec := doMux(t, muxActions(rt), http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil,
			etActionHeaders("digest-key", created.Metadata.ResourceVersion))
		if rec.Code != http.StatusConflict ||
			!hasViolation(decodeProblem(t, rec), cloudmodel.ViolationIdempotencyKeyReuseMismatch) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
		rt.life.AbortReservation(ns)
	})

	t.Run("different-key while Qualifying", func(t *testing.T) {
		fixtures := executiontarget.NewMapFixtureSource()
		started := make(chan struct{})
		release := make(chan struct{})
		blocking := &blockingActionFixture{inner: fixtures, started: started, release: release}
		rt := newETActionRuntimeWithObserver(
			executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true),
			executiontarget.NewSyntheticObserver(blocking, nil),
			fixtures,
		)
		created := createTargetForAction(t, rt, "inprog", "create-inprog")
		mux := muxActions(rt)

		go func() {
			_ = doMux(t, mux, http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil,
				etActionHeaders("owner-q", created.Metadata.ResourceVersion))
		}()
		select {
		case <-started:
		case <-time.After(3 * time.Second):
			t.Fatal("owner qualify did not start")
		}
		rec := doMux(t, mux, http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil,
			etActionHeaders("other-key", created.Metadata.ResourceVersion))
		if rec.Code != http.StatusConflict ||
			!hasViolation(decodeProblem(t, rec), executiontarget.ViolationTargetQualificationInProgress) {
			t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
		}
		close(release)
	})
}

func TestExecutionTargetActions_FailedRetireAppendLeavesQualifyIntact(t *testing.T) {
	fixtures := executiontarget.NewMapFixtureSource()
	started := make(chan struct{})
	release := make(chan struct{})
	blocking := &blockingActionFixture{inner: fixtures, started: started, release: release}
	rt := newETActionRuntimeWithObserver(
		executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true),
		executiontarget.NewSyntheticObserver(blocking, nil),
		fixtures,
	)
	created := createTargetForAction(t, rt, "fail-retire", "create-fr")
	mux := muxActions(rt)

	qualDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		qualDone <- doMux(t, mux, http.MethodPost, etActionPath(created.Metadata.UID, "qualify"), nil,
			etActionHeaders("q-survive", created.Metadata.ResourceVersion))
	}()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("qualify did not start")
	}

	rt.audit.SetFail(errors.New("retire append boom"))
	retireRec := doMux(t, mux, http.MethodPost, etActionPath(created.Metadata.UID, "retire"), nil,
		etActionHeaders("retire-fail", created.Metadata.ResourceVersion))
	if retireRec.Code != http.StatusInternalServerError {
		t.Fatalf("retire append failure: %d %s", retireRec.Code, retireRec.Body.String())
	}
	rt.audit.SetFail(nil)
	close(release)

	var qRec *httptest.ResponseRecorder
	select {
	case qRec = <-qualDone:
	case <-time.After(3 * time.Second):
		t.Fatal("qualify did not finish")
	}
	if qRec.Code != http.StatusOK {
		t.Fatalf("qualify must complete after failed retire: %d %s", qRec.Code, qRec.Body.String())
	}
	et, _ := rt.life.Store().GetExecutionTarget(created.Metadata.UID)
	if et.Status.Lifecycle != etmodel.LifecycleActive {
		t.Fatal("failed retire must not publish retirement")
	}
}

func TestExecutionTargetActions_RetireCancelledContext(t *testing.T) {
	rt := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
	created := createTargetForAction(t, rt, "cancel-retire", "create-cr")
	mux := muxActions(rt)

	req := httptest.NewRequest(http.MethodPost, etActionPath(created.Metadata.UID, "retire"), nil)
	for k, v := range etActionHeaders("cancel-retire", created.Metadata.ResourceVersion) {
		req.Header.Set(k, v)
	}
	req = withAuth(req, testPrincipal)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("cancelled retire: %d %s", rec.Code, rec.Body.String())
	}
	et, _ := rt.life.Store().GetExecutionTarget(created.Metadata.UID)
	if et.Status.Lifecycle != etmodel.LifecycleActive {
		t.Fatal("cancelled retire must not publish")
	}
}

func TestExecutionTargetActions_RetireWithoutWriteGrant(t *testing.T) {
	rtWrite := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, true, true, true))
	created := createTargetForAction(t, rtWrite, "no-write", "create-nw")
	rtDenied := newETActionRuntime(executionTargetActionGrants(testPrincipal, testProviderUID, false, true, true))
	// Share lifecycle/cloud/audit so the target exists under denied grants.
	rtDenied.life = rtWrite.life
	rtDenied.cloud = rtWrite.cloud
	rtDenied.audit = rtWrite.audit
	rtDenied.retire = NewExecutionTargetRetireHandler(rtWrite.life, rtWrite.cloud,
		executionTargetActionGrants(testPrincipal, testProviderUID, false, true, true), rtWrite.audit)
	before := rtWrite.audit.Len()
	rec := doMux(t, muxActions(rtDenied), http.MethodPost, etActionPath(created.Metadata.UID, "retire"), nil,
		etActionHeaders("retire-deny", created.Metadata.ResourceVersion))
	if rec.Code != http.StatusForbidden || decodeProblem(t, rec).Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("got %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rtWrite.audit.Len() != before+1 {
		t.Fatal("retire authz denial must append exactly one AuditEvent")
	}
}
