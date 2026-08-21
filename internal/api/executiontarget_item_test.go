package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/executiontarget"
	etmodel "github.com/sanjeevksaini/sovrunn/internal/executiontarget/model"
)

const etItemBase = "/apis/execution.sovrunn.io/v1alpha1/execution-targets/"

type etItemRuntime struct {
	*etRuntime
	item *ExecutionTargetItemHandler
}

func newETItemRuntime(grants GrantResolver) *etItemRuntime {
	rt := newETRuntime(grants)
	item := NewExecutionTargetItemHandler(rt.life, rt.cloud, grants, rt.audit)
	return &etItemRuntime{etRuntime: rt, item: item}
}

func etItemPath(uid string) string {
	return etItemBase + uid
}

// muxItem registers collection + item GET so PathValue("uid") is populated.
func muxItem(rt *etItemRuntime) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET "+etCollection, rt.handler)
	mux.Handle("POST "+etCollection, rt.handler)
	mux.Handle(PatternExecutionTargetItemGet, rt.item)
	return mux
}

func createExecutionTargetForItem(t *testing.T, rt *etItemRuntime, name, key string) executiontarget.SafeExecutionTarget {
	t.Helper()
	seedActiveParticipationAndStack(t, rt.etRuntime)
	rec := doMux(t, muxItem(rt), http.MethodPost, etCollection,
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

func TestExecutionTargetItem_GetHappyPathETagSafeProjectionNoIdempotency(t *testing.T) {
	rt := newETItemRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
	created := createExecutionTargetForItem(t, rt, "item-get-ok", "item-create-ok")
	beforeAudit := rt.audit.Len()

	headers := map[string]string{
		"Idempotency-Key": "ignored-on-item-get",
		"Accept":          "application/xml",
	}
	rec := doMux(t, muxItem(rt), http.MethodGet, etItemPath(created.Metadata.UID), nil, headers)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	assertJSONNosniff(t, rec)
	if got := rec.Header().Get("ETag"); got == "" || got != created.Metadata.ResourceVersion {
		t.Fatalf("ETag=%q, want %q", got, created.Metadata.ResourceVersion)
	}
	if rec.Header().Get("Idempotency-Key") != "" {
		t.Fatal("response must not echo Idempotency-Key")
	}
	if rt.audit.Len() != beforeAudit {
		t.Fatal("successful item GET must not append AuditEvent")
	}

	var proj executiontarget.SafeExecutionTarget
	if err := json.Unmarshal(rec.Body.Bytes(), &proj); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if proj.Metadata.UID != created.Metadata.UID ||
		proj.Status.Lifecycle != etmodel.LifecycleActive ||
		proj.Status.Qualification != etmodel.QualificationUnqualified {
		t.Fatalf("projection=%#v", proj)
	}
	raw := rec.Body.String()
	for _, forbidden := range []string{
		"factSetRef", "qualificationResultRef", "observerID", "NormalizedTargetFactSet",
		"TargetQualificationResult", "handle", "facts", "conditions",
	} {
		if strings.Contains(raw, forbidden) {
			t.Fatalf("safe projection must omit %q; body=%s", forbidden, raw)
		}
	}

	// GET must not create, wait on, read, or replay an idempotency record:
	// the same key remains available for a later create reservation.
	ns := executiontarget.IdempotencyNamespace{
		PrincipalUID:     testPrincipal,
		RoutePattern:     PatternExecutionTargetCreate,
		CloudProviderUID: testProviderUID,
		Key:              "ignored-on-item-get",
	}
	inspect := rt.life.ReserveOrInspectReplay(ns, executiontarget.DigestAction())
	if inspect.Outcome != executiontarget.ReservationOutcomeReserved {
		t.Fatalf("GET must not consume idempotency key; outcome=%v", inspect.Outcome)
	}
	rt.life.AbortReservation(ns)
}

func TestExecutionTargetItem_AuthSafeDenialAndWithoutRead(t *testing.T) {
	t.Run("missing auth 401", func(t *testing.T) {
		rt := newETItemRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
		rec := doMuxUnauth(t, muxItem(rt), http.MethodGet, etItemPath("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), nil, nil)
		if rec.Code != http.StatusUnauthorized || rt.audit.Len() != 0 {
			t.Fatalf("want 401 no audit, got %d audit=%d", rec.Code, rt.audit.Len())
		}
		assertProblemNosniff(t, rec)
		p := decodeProblem(t, rec)
		if p.Code != apiproblem.CodeAuthRequired {
			t.Fatalf("code=%s", p.Code)
		}
	})

	t.Run("inaccessible direct GET safe 404 F16-17", func(t *testing.T) {
		rt := newETItemRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
		rec := doMux(t, muxItem(rt), http.MethodGet, etItemPath("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"), nil, nil)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
		}
		assertProblemNosniff(t, rec)
		p := decodeProblem(t, rec)
		if p.Code != apiproblem.CodeResourceNotFound || !hasViolation(p, violationAuthorizationSafeDenial) {
			t.Fatalf("want safe 404, got %#v", p)
		}
		if rt.audit.Len() != 1 {
			t.Fatalf("safe denial audit=%d", rt.audit.Len())
		}
	})

	t.Run("without read after scope derive 403", func(t *testing.T) {
		writeOnly := executionTargetGrants(testPrincipal, testProviderUID, true, false)
		rtWrite := newETItemRuntime(writeOnly)
		created := createExecutionTargetForItem(t, rtWrite, "item-no-read", "item-no-read-create")

		failAudit := &cloudmodel.MemoryAuditAppender{}
		item := NewExecutionTargetItemHandler(rtWrite.life, rtWrite.cloud, writeOnly, failAudit)
		rtDenied := &etItemRuntime{etRuntime: rtWrite.etRuntime, item: item}
		rtDenied.audit = failAudit

		rec := doMux(t, muxItem(rtDenied), http.MethodGet, etItemPath(created.Metadata.UID), nil, nil)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
		}
		assertProblemNosniff(t, rec)
		p := decodeProblem(t, rec)
		if p.Code != apiproblem.CodeAuthorizationDenied {
			t.Fatalf("code=%s", p.Code)
		}
		if failAudit.Len() != 1 {
			t.Fatalf("authz denial audit=%d", failAudit.Len())
		}
	})
}

func TestExecutionTargetItem_DenialAppendFailure_F16_102_and_F16_112(t *testing.T) {
	t.Run("inaccessible GET append failure F16-102", func(t *testing.T) {
		rt := newETItemRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
		rt.audit.SetFail(errors.New("append boom"))
		rec := doMux(t, muxItem(rt), http.MethodGet, etItemPath("cccccccccccccccccccccccccccccccc"), nil, map[string]string{
			"Accept": "text/plain",
		})
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
		}
		assertProblemNosniff(t, rec)
		p := decodeProblem(t, rec)
		if p.Code != apiproblem.CodeInternalError {
			t.Fatalf("code=%s", p.Code)
		}
		if p.Status == http.StatusNotFound || hasViolation(p, violationAuthorizationSafeDenial) {
			t.Fatal("must not disclose suppressed safe 404")
		}
		if rt.audit.Len() != 0 {
			t.Fatal("failed append must leave no durable AuditEvent")
		}
	})

	t.Run("direct GET without read append failure F16-112 member", func(t *testing.T) {
		writeOnly := executionTargetGrants(testPrincipal, testProviderUID, true, false)
		rtWrite := newETItemRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
		created := createExecutionTargetForItem(t, rtWrite, "item-112", "item-112-create")

		failAudit := &cloudmodel.MemoryAuditAppender{}
		failAudit.SetFail(errors.New("append boom"))
		item := NewExecutionTargetItemHandler(rtWrite.life, rtWrite.cloud, writeOnly, failAudit)
		rtDenied := &etItemRuntime{etRuntime: rtWrite.etRuntime, item: item}
		rtDenied.audit = failAudit

		rec := doMux(t, muxItem(rtDenied), http.MethodGet, etItemPath(created.Metadata.UID), nil, nil)
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
		if failAudit.Len() != 0 {
			t.Fatal("failed append must leave no durable AuditEvent")
		}
	})
}

func TestExecutionTargetItem_HTTPContract_F16_45_48_106_107_108(t *testing.T) {
	rt := newETItemRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
	created := createExecutionTargetForItem(t, rt, "item-contract", "item-contract-create")
	mux := muxItem(rt)

	t.Run("Accept ignored success F16-45", func(t *testing.T) {
		rec := doMux(t, mux, http.MethodGet, etItemPath(created.Metadata.UID), nil, map[string]string{
			"Accept": "application/xml",
		})
		if rec.Code != http.StatusOK {
			t.Fatalf("status=%d", rec.Code)
		}
		assertJSONNosniff(t, rec)
	})

	t.Run("Accept ignored problem F16-48", func(t *testing.T) {
		rec := doMux(t, mux, http.MethodGet, etItemPath("dddddddddddddddddddddddddddddddd"), nil, map[string]string{
			"Accept": "text/html",
		})
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status=%d", rec.Code)
		}
		assertProblemNosniff(t, rec)
	})

	t.Run("query rejected before safe resolution F16-106 item member", func(t *testing.T) {
		before := rt.audit.Len()
		rec := doMux(t, mux, http.MethodGet, etItemPath(created.Metadata.UID)+"?x=1", nil, nil)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
		}
		p := decodeProblem(t, rec)
		if p.Code != apiproblem.CodeMalformedRequest {
			t.Fatalf("code=%s", p.Code)
		}
		if rt.audit.Len() != before {
			t.Fatal("query rejection must not audit or resolve")
		}
	})

	t.Run("malformed UID before safe resolution F16-107", func(t *testing.T) {
		before := rt.audit.Len()
		for _, bad := range []string{
			"not-a-uid",
			"ABCDEFABCDEFABCDEFABCDEFABCDEFAB", // uppercase
			"short",
			"zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz", // 33 chars
		} {
			rec := doMux(t, mux, http.MethodGet, etItemPath(bad), nil, nil)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("uid=%q status=%d body=%s", bad, rec.Code, rec.Body.String())
			}
			p := decodeProblem(t, rec)
			if p.Code != apiproblem.CodeMalformedRequest {
				t.Fatalf("uid=%q code=%s", bad, p.Code)
			}
		}
		if rt.audit.Len() != before {
			t.Fatal("malformed UID must not audit or lookup")
		}
	})

	t.Run("trailing slash transport-only F16-108", func(t *testing.T) {
		// Item registration is exact `{uid}` (no trailing slash). ServeMux
		// therefore does not invoke the handler; transport-only 404/nosniff
		// for trailing-slash forms is owned by ExecutionTargetTransportGuard
		// (TASK-F16-08). Prove the item pattern is not registered with a
		// trailing slash and the handler is never reached.
		var called bool
		spy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			rt.item.ServeHTTP(w, r)
		})
		itemMux := http.NewServeMux()
		itemMux.Handle(PatternExecutionTargetItemGet, spy)
		rec := doMux(t, itemMux, http.MethodGet, etItemPath(created.Metadata.UID)+"/", nil, nil)
		if called {
			t.Fatal("trailing-slash path must not reach the item handler")
		}
		if rec.Code != http.StatusNotFound {
			t.Fatalf("mux trailing-slash status=%d", rec.Code)
		}
		if strings.Contains(rec.Body.String(), "problem") || strings.Contains(rec.Body.String(), "AUTH") {
			t.Fatalf("trailing slash must not return Problem body: %s", rec.Body.String())
		}

		guardSrc, err := os.ReadFile("../server/executiontarget_guard.go")
		if err != nil {
			t.Fatalf("read guard: %v", err)
		}
		src := string(guardSrc)
		for _, want := range []string{
			"Trailing-slash item path",
			"writeExecutionTargetTransport(w, http.StatusNotFound, \"\")",
			"contentTypeOptionsNosniff",
		} {
			if !strings.Contains(src, want) {
				t.Fatalf("guard.go missing transport trailing-slash contract %q", want)
			}
		}
	})
}

func TestExecutionTargetItem_RouteBuilderItemRegistration(t *testing.T) {
	rt := newETItemRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
	created := createExecutionTargetForItem(t, rt, "item-routed", "item-routed-create")
	mux := muxItem(rt)

	if PatternExecutionTargetItemGet != "GET "+etItemBase+"{uid}" {
		t.Fatalf("item pattern=%q", PatternExecutionTargetItemGet)
	}

	routesSrc, err := os.ReadFile("../server/routes.go")
	if err != nil {
		t.Fatalf("read routes.go: %v", err)
	}
	src := string(routesSrc)
	for _, want := range []string{
		"func NewExecutionTargetMux(",
		`routeExecutionTargetGet    = "GET /apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}"`,
		"registerExecutionTargetPattern(mux, routeExecutionTargetGet, h.Item)",
		"Item       http.Handler",
	} {
		if !strings.Contains(src, want) {
			t.Fatalf("routes.go missing %q", want)
		}
	}
	for _, forbidden := range []string{
		`"/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/qualify"`,
		`"/apis/execution.sovrunn.io/v1alpha1/execution-targets/{uid}/actions/retire"`,
		"StatusNotImplemented",
		"http.StatusNotImplemented",
	} {
		if strings.Contains(src, forbidden) {
			t.Fatalf("routes.go must not contain %q until TASK-F16-11", forbidden)
		}
	}

	rec := doMux(t, mux, http.MethodGet, etItemPath(created.Metadata.UID), nil, map[string]string{
		"Accept": "text/plain",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("item GET via mux status=%d body=%s", rec.Code, rec.Body.String())
	}
	assertJSONNosniff(t, rec)
	if rec.Header().Get("ETag") == "" {
		t.Fatal("item GET must return ETag")
	}

	// Qualify/retire remain unregistered (TASK-F16-11).
	rec = doMux(t, mux, http.MethodPost, etItemPath(created.Metadata.UID)+"/actions/qualify", nil, nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("qualify must be unregistered: %d", rec.Code)
	}
}

func TestExecutionTargetItem_IntegrationEndToEnd(t *testing.T) {
	rt := newETItemRuntime(executionTargetGrants(testPrincipal, testProviderUID, true, true))
	created := createExecutionTargetForItem(t, rt, "item-e2e", "item-e2e-create")
	mux := muxItem(rt)

	rec := doMux(t, mux, http.MethodGet, etItemPath(created.Metadata.UID), nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("e2e GET status=%d body=%s", rec.Code, rec.Body.String())
	}
	var proj executiontarget.SafeExecutionTarget
	if err := json.Unmarshal(rec.Body.Bytes(), &proj); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if proj.Metadata.UID != created.Metadata.UID || proj.EffectiveAvailability == "" {
		t.Fatalf("e2e projection=%#v", proj)
	}

	// Safe-denial path end-to-end through the same mux.
	before := rt.audit.Len()
	rec = doMux(t, mux, http.MethodGet, etItemPath("eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"), nil, nil)
	if rec.Code != http.StatusNotFound || !hasViolation(decodeProblem(t, rec), violationAuthorizationSafeDenial) {
		t.Fatalf("e2e safe denial=%d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before+1 {
		t.Fatalf("e2e safe denial audit=%d", rt.audit.Len())
	}
}
