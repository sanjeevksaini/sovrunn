package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
	"github.com/sanjeevksaini/sovrunn/internal/requestctx"
)

const testPrincipal = "bootstrap-principal"

type staticGrants struct {
	principal string
	grants    []cloudmodel.Grant
}

func (s *staticGrants) PrincipalID() string { return s.principal }

func (s *staticGrants) CoarseLookup(action string) []cloudmodel.Grant {
	var out []cloudmodel.Grant
	for _, g := range s.grants {
		if g.Action == action {
			out = append(out, g)
		}
	}
	return out
}

func (s *staticGrants) AuthorizeExact(action string, scope apimeta.ScopeIdentity, resourceUID, relatedPlatformUID string) bool {
	for _, g := range s.grants {
		if g.Action != action || g.Scope != scope {
			continue
		}
		if g.ResourceUID != "" && g.ResourceUID != resourceUID {
			continue
		}
		return true
	}
	return false
}

func fullPlatformGrants(principal string) *staticGrants {
	return &staticGrants{
		principal: principal,
		grants: []cloudmodel.Grant{
			{Action: cloudmodel.ActionCloudPlatformWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderRead, Scope: cloudmodel.PlatformRootScope},
		},
	}
}

func readOnlyPlatformGrants(principal string) *staticGrants {
	return &staticGrants{
		principal: principal,
		grants: []cloudmodel.Grant{
			{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderRead, Scope: cloudmodel.PlatformRootScope},
		},
	}
}

func noReadPlatformGrants(principal string) *staticGrants {
	return &staticGrants{
		principal: principal,
		grants: []cloudmodel.Grant{
			{Action: cloudmodel.ActionCloudPlatformWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderWrite, Scope: cloudmodel.PlatformRootScope},
		},
	}
}

type f15Runtime struct {
	store  *cloudmodel.Store
	audit  *cloudmodel.MemoryAuditAppender
	idemp  *cloudmodel.IdempotencyCoordinator
	pub    *cloudmodel.PublicationCoordinator
	grants GrantResolver
}

func newF15Runtime(grants GrantResolver) *f15Runtime {
	store := cloudmodel.NewStore()
	audit := &cloudmodel.MemoryAuditAppender{}
	idemp := cloudmodel.NewIdempotencyCoordinator(cloudmodel.IdempotencyCoordinatorConfig{})
	pub := cloudmodel.NewPublicationCoordinator(cloudmodel.PublicationCoordinatorConfig{
		Store: store, Audit: audit, Idempotency: idemp,
	})
	return &f15Runtime{store: store, audit: audit, idemp: idemp, pub: pub, grants: grants}
}

func (rt *f15Runtime) platformCollection() http.Handler {
	return NewCloudPlatformCollectionHandler(rt.store, rt.grants, rt.pub, rt.idemp)
}

func (rt *f15Runtime) platformItem() http.Handler {
	return NewCloudPlatformItemHandler(rt.store, rt.grants, rt.pub)
}

func (rt *f15Runtime) providerCollection() http.Handler {
	return NewCloudProviderCollectionHandler(rt.store, rt.grants, rt.pub, rt.idemp)
}

func (rt *f15Runtime) providerItem() http.Handler {
	return NewCloudProviderItemHandler(rt.store, rt.grants, rt.pub)
}

func withAuth(r *http.Request, principal string) *http.Request {
	r.Header.Set("Authorization", "Bearer "+principal)
	r.Header.Set("X-Sovrunn-Request-ID", "req-test-001")
	ctx := requestctx.WithRequestID(r.Context(), "req-test-001")
	return r.WithContext(ctx)
}

func platformCreateBody(name string) []byte {
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{"ownerRegistration":{"legalName":"Acme","registrationIdentifier":"R1","jurisdictionCode":"US"}}
	}`)
}

func muxPlatform(rt *f15Runtime) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /apis/core.sovrunn.io/v1alpha1/cloud-platforms", rt.platformCollection())
	mux.Handle("POST /apis/core.sovrunn.io/v1alpha1/cloud-platforms", rt.platformCollection())
	mux.Handle("GET /apis/core.sovrunn.io/v1alpha1/cloud-platforms/{uid}", rt.platformItem())
	mux.Handle("PATCH /apis/core.sovrunn.io/v1alpha1/cloud-platforms/{uid}", rt.platformItem())
	mux.Handle("PUT /apis/core.sovrunn.io/v1alpha1/cloud-platforms/{uid}", rt.platformItem())
	mux.Handle("DELETE /apis/core.sovrunn.io/v1alpha1/cloud-platforms/{uid}", rt.platformItem())
	return mux
}

func muxProvider(rt *f15Runtime) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /apis/core.sovrunn.io/v1alpha1/cloud-providers", rt.providerCollection())
	mux.Handle("POST /apis/core.sovrunn.io/v1alpha1/cloud-providers", rt.providerCollection())
	mux.Handle("GET /apis/core.sovrunn.io/v1alpha1/cloud-providers/{uid}", rt.providerItem())
	mux.Handle("PATCH /apis/core.sovrunn.io/v1alpha1/cloud-providers/{uid}", rt.providerItem())
	mux.Handle("PUT /apis/core.sovrunn.io/v1alpha1/cloud-providers/{uid}", rt.providerItem())
	mux.Handle("DELETE /apis/core.sovrunn.io/v1alpha1/cloud-providers/{uid}", rt.providerItem())
	return mux
}

func decodeProblem(t *testing.T, rec *httptest.ResponseRecorder) *apiproblem.Problem {
	t.Helper()
	var p apiproblem.Problem
	if err := json.Unmarshal(rec.Body.Bytes(), &p); err != nil {
		t.Fatalf("problem decode: %v body=%s", err, rec.Body.String())
	}
	return &p
}

func TestCloudPlatformCollection_ListCreateAuthAuditIdempotency(t *testing.T) {
	rt := newF15Runtime(fullPlatformGrants(testPrincipal))
	mux := muxPlatform(rt)

	// Missing auth → AUTH_REQUIRED, no audit.
	r := httptest.NewRequest(http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, r)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("missing auth status=%d", rec.Code)
	}
	if rt.audit.Len() != 0 {
		t.Fatalf("auth must not audit, got %d", rt.audit.Len())
	}

	// LIST empty authorized.
	rec = doMux(t, mux, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("list empty status=%d body=%s", rec.Code, rec.Body.String())
	}

	// Create closed contract → 201 + Active + audit.
	headers := map[string]string{
		"Content-Type":    "application/json",
		"Idempotency-Key": "create-1",
	}
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformCreateBody("plat-a"), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created model.CloudPlatform
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("create decode: %v", err)
	}
	if created.Status.Phase != model.CloudPlatformPhaseActive {
		t.Fatalf("initial status=%q", created.Status.Phase)
	}
	if created.Metadata.ScopeRef != nil {
		t.Fatalf("Platform-root scope must be nil, got %#v", created.Metadata.ScopeRef)
	}
	if created.Metadata.ResourceVersion != "1" || created.Metadata.UID == "" {
		t.Fatalf("server metadata incomplete: %#v", created.Metadata)
	}
	if rt.audit.Len() != 1 {
		t.Fatalf("create audit count=%d", rt.audit.Len())
	}

	// Same-key/same-digest replay → no additional audit.
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformCreateBody("plat-a"), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("replay status=%d", rec.Code)
	}
	if rt.audit.Len() != 1 {
		t.Fatalf("replay must not audit, count=%d", rt.audit.Len())
	}

	// Same-key/different-digest → CONFLICT, no audit.
	headers["Idempotency-Key"] = "create-1"
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformCreateBody("plat-b"), headers)
	if rec.Code != http.StatusConflict {
		t.Fatalf("digest mismatch status=%d", rec.Code)
	}
	p := decodeProblem(t, rec)
	if !hasViolation(p, cloudmodel.ViolationIdempotencyKeyReuseMismatch) {
		t.Fatalf("want VS0_IDEMPOTENCY_KEY_REUSE_MISMATCH, got %#v", p)
	}
	if rt.audit.Len() != 1 {
		t.Fatalf("mismatch must not audit")
	}

	// Duplicate name → ALREADY_EXISTS, no audit.
	headers["Idempotency-Key"] = "create-2"
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformCreateBody("plat-a"), headers)
	if rec.Code != http.StatusConflict {
		t.Fatalf("duplicate status=%d body=%s", rec.Code, rec.Body.String())
	}
	if decodeProblem(t, rec).Code != apiproblem.CodeAlreadyExists {
		t.Fatalf("duplicate code=%s", decodeProblem(t, rec).Code)
	}
	if rt.audit.Len() != 1 {
		t.Fatalf("duplicate must not audit")
	}

	// Malformed → no audit.
	headers["Idempotency-Key"] = "create-3"
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", []byte(`{`), headers)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("malformed status=%d", rec.Code)
	}
	if rt.audit.Len() != 1 {
		t.Fatalf("malformed must not audit")
	}

	// Client status → audited denial.
	headers["Idempotency-Key"] = "create-4"
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", []byte(`{
		"metadata":{"name":"plat-status"},
		"spec":{"ownerRegistration":{"legalName":"Acme","registrationIdentifier":"R1","jurisdictionCode":"US"}},
		"status":{"phase":"Active"}
	}`), headers)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status write status=%d", rec.Code)
	}
	if !hasViolation(decodeProblem(t, rec), validate.ViolationStatusFieldWrite) {
		t.Fatalf("status write violation missing")
	}
	if rt.audit.Len() != 2 {
		t.Fatalf("status write must audit exactly once more, count=%d", rt.audit.Len())
	}

	// Client system-owned metadata → audited denial.
	headers["Idempotency-Key"] = "create-4b"
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", []byte(`{
		"metadata":{"name":"plat-uid","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		"spec":{"ownerRegistration":{"legalName":"Acme","registrationIdentifier":"R1","jurisdictionCode":"US"}}
	}`), headers)
	if rec.Code != http.StatusForbidden || !hasViolation(decodeProblem(t, rec), validate.ViolationSystemOwnedFieldWrite) {
		t.Fatalf("system-owned: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != 3 {
		t.Fatalf("system-owned must audit, count=%d", rt.audit.Len())
	}

	// scopeRef → unaudited validation failure.
	before := rt.audit.Len()
	headers["Idempotency-Key"] = "create-5"
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", []byte(`{
		"metadata":{"name":"plat-scope","scopeRef":{"apiVersion":"v","kind":"Platform","name":"platform","uid":"platform"}},
		"spec":{"ownerRegistration":{"legalName":"Acme","registrationIdentifier":"R1","jurisdictionCode":"US"}}
	}`), headers)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("scopeRef status=%d", rec.Code)
	}
	if rt.audit.Len() != before {
		t.Fatalf("scopeRef must not audit")
	}

	// LIST ascending uid.
	headers["Idempotency-Key"] = "create-6"
	_ = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformCreateBody("plat-z"), headers)
	rec = doMux(t, mux, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", nil, nil)
	var list apimeta.ListEnvelope[model.CloudPlatform]
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("list decode: %v", err)
	}
	if len(list.Items) < 2 {
		t.Fatalf("list items=%d", len(list.Items))
	}
	for i := 1; i < len(list.Items); i++ {
		if list.Items[i-1].Metadata.UID >= list.Items[i].Metadata.UID {
			t.Fatalf("list not ascending by uid")
		}
	}

	// LIST without read grant → 403 + one audit.
	rt2 := newF15Runtime(noReadPlatformGrants(testPrincipal))
	mux2 := muxPlatform(rt2)
	rec = doMux(t, mux2, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", nil, nil)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("list no-read status=%d", rec.Code)
	}
	if rt2.audit.Len() != 1 {
		t.Fatalf("list no-read audit=%d", rt2.audit.Len())
	}

	// Audit-append failure → INTERNAL_ERROR, no publication.
	rt3 := newF15Runtime(fullPlatformGrants(testPrincipal))
	mux3 := muxPlatform(rt3)
	rt3.audit.SetFail(errors.New("audit down"))
	headers = map[string]string{"Content-Type": "application/json", "Idempotency-Key": "fail-audit"}
	rec = doMux(t, mux3, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", platformCreateBody("plat-fail"), headers)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("audit fail status=%d body=%s", rec.Code, rec.Body.String())
	}
	if rt3.store.HasCloudPlatform() {
		t.Fatal("audit failure must not publish")
	}
}

func TestCloudPlatformCollection_StageSetInvoked(t *testing.T) {
	rt := newF15Runtime(fullPlatformGrants(testPrincipal))
	mux := muxPlatform(rt)
	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "stageset-1"}
	// Invalid jurisdiction fails inside StageSet/semantic after closed-contract decode.
	body := []byte(`{
		"metadata":{"name":"plat-iso"},
		"spec":{"ownerRegistration":{"legalName":"Acme","registrationIdentifier":"R1","jurisdictionCode":"ZZ"}}
	}`)
	rec := doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", body, headers)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("StageSet ISO status=%d body=%s", rec.Code, rec.Body.String())
	}
	stages, ok := validate.StageSetFor(validate.KindCloudPlatform, validate.RequestClassCollectionCreate)
	if !ok || stages.Defaulting == nil || stages.Semantic == nil || stages.Reference == nil {
		t.Fatal("FEATURE-0015 must plug validators into inherited StageSet")
	}
}

func TestWriteRawJSON_AllowsOnlyJSONMediaAndReplayForcesJSON(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-platforms", nil)

	t.Run("unrecognized media type falls back to JSON", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writeRawJSON(recorder, request, http.StatusOK, "text/html", []byte(`{"ok":true}`))
		if got := recorder.Header().Get("Content-Type"); got != mediaJSON {
			t.Fatalf("Content-Type=%q, want %q", got, mediaJSON)
		}
		if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Fatalf("X-Content-Type-Options=%q, want nosniff", got)
		}
	})

	t.Run("stored replay media type cannot override JSON", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		writeReplay(recorder, request, &cloudmodel.ReplayResponse{
			StatusCode:  http.StatusCreated,
			ContentType: "text/html",
			Body:        []byte(`{"replayed":true}`),
		})
		if got := recorder.Header().Get("Content-Type"); got != mediaJSON {
			t.Fatalf("replay Content-Type=%q, want %q", got, mediaJSON)
		}
		if got := recorder.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Fatalf("replay X-Content-Type-Options=%q, want nosniff", got)
		}
	})
}

func doMux(t *testing.T, mux http.Handler, method, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == nil {
		r = httptest.NewRequest(method, path, nil)
	} else {
		r = httptest.NewRequest(method, path, bytes.NewReader(body))
	}
	r = withAuth(r, testPrincipal)
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, r)
	return rec
}
