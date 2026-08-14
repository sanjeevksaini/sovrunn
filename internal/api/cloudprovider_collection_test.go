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
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func providerCreateBody(name string, markets string) []byte {
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{"operatingMarkets":` + markets + `}
	}`)
}

func seedPlatformForProvider(t *testing.T, rt *f15Runtime) {
	t.Helper()
	if _, prob := rt.store.CreateCloudPlatform(model.CloudPlatform{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform},
		Metadata: apimeta.ObjectMeta{Name: "root-plat", UID: "ffffffffffffffffffffffffffffffff"},
		Spec: model.CloudPlatformSpec{OwnerRegistration: model.OwnerRegistration{
			LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
		}},
		Status: model.CloudPlatformStatus{Phase: model.CloudPlatformPhaseActive},
	}); prob != nil {
		t.Fatalf("seed platform: %#v", prob)
	}
}

func TestCloudProviderCollection_ListCreateRootGateISO(t *testing.T) {
	rt := newF15Runtime(fullPlatformGrants(testPrincipal))
	mux := muxProvider(rt)

	// Missing auth.
	r := doMuxUnauth(t, mux, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-providers", nil, nil)
	if r.Code != http.StatusUnauthorized || rt.audit.Len() != 0 {
		t.Fatalf("auth required: %d audit=%d", r.Code, rt.audit.Len())
	}

	// Root gate before any CloudPlatform.
	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "prov-1"}
	rec := doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers",
		providerCreateBody("prov-a", `["US","IN"]`), headers)
	if rec.Code != http.StatusConflict {
		t.Fatalf("root gate status=%d body=%s", rec.Code, rec.Body.String())
	}
	if !hasViolation(decodeProblem(t, rec), violationCloudPlatformRootRequired) {
		t.Fatalf("root gate violation missing: %#v", decodeProblem(t, rec))
	}
	if rt.audit.Len() != 1 {
		t.Fatalf("root gate audit=%d", rt.audit.Len())
	}

	// Audit-append failure on root gate → INTERNAL_ERROR, no publication of denial.
	rtFail := newF15Runtime(fullPlatformGrants(testPrincipal))
	muxFail := muxProvider(rtFail)
	rtFail.audit.SetFail(errors.New("audit down"))
	rec = doMux(t, muxFail, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers",
		providerCreateBody("prov-fail", `["US"]`), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "prov-fail",
		})
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("root audit fail status=%d", rec.Code)
	}
	if rtFail.audit.Len() != 0 {
		t.Fatal("failed append must publish neither denial nor resource")
	}

	seedPlatformForProvider(t, rt)
	// Authorized create with assigned ISO.
	headers["Idempotency-Key"] = "prov-2"
	before := rt.audit.Len()
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers",
		providerCreateBody("prov-a", `["US","IN"]`), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created model.CloudProvider
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	if created.Status.Phase != model.CloudProviderPhaseActive || created.Metadata.ScopeRef != nil {
		t.Fatalf("created=%#v", created)
	}
	if rt.audit.Len() != before+1 {
		t.Fatalf("create audit=%d", rt.audit.Len())
	}

	// Unassigned ZZ → VALIDATION_FAILED, no audit.
	before = rt.audit.Len()
	headers["Idempotency-Key"] = "prov-zz"
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers",
		providerCreateBody("prov-zz", `["ZZ"]`), headers)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("ZZ status=%d body=%s", rec.Code, rec.Body.String())
	}
	if rt.audit.Len() != before {
		t.Fatal("ZZ must not audit")
	}

	// Duplicate name.
	before = rt.audit.Len()
	headers["Idempotency-Key"] = "prov-dup"
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers",
		providerCreateBody("prov-a", `["US"]`), headers)
	if decodeProblem(t, rec).Code != apiproblem.CodeAlreadyExists || rt.audit.Len() != before {
		t.Fatalf("duplicate: %d %#v audit=%d", rec.Code, decodeProblem(t, rec), rt.audit.Len())
	}

	// Replay / mismatch no extra audit.
	headers["Idempotency-Key"] = "prov-2"
	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers",
		providerCreateBody("prov-a", `["US","IN"]`), headers)
	if rec.Code != http.StatusCreated || rt.audit.Len() != before {
		t.Fatalf("replay: %d audit=%d", rec.Code, rt.audit.Len())
	}
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers",
		providerCreateBody("prov-b", `["US"]`), headers)
	if rec.Code != http.StatusConflict || rt.audit.Len() != before {
		t.Fatalf("mismatch: %d audit=%d", rec.Code, rt.audit.Len())
	}

	// status audited; scopeRef unaudited.
	headers["Idempotency-Key"] = "prov-status"
	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers", []byte(`{
		"metadata":{"name":"prov-status"},
		"spec":{"operatingMarkets":["US"]},
		"status":{"phase":"Active"}
	}`), headers)
	if rec.Code != http.StatusForbidden || rt.audit.Len() != before+1 {
		t.Fatalf("status: %d audit=%d", rec.Code, rt.audit.Len())
	}
	before = rt.audit.Len()
	headers["Idempotency-Key"] = "prov-scope"
	rec = doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers", []byte(`{
		"metadata":{"name":"prov-scope","scopeRef":{"apiVersion":"v","kind":"Platform","name":"p","uid":"platform"}},
		"spec":{"operatingMarkets":["US"]}
	}`), headers)
	if rec.Code != http.StatusUnprocessableEntity || rt.audit.Len() != before {
		t.Fatalf("scopeRef: %d audit=%d", rec.Code, rt.audit.Len())
	}

	// LIST without read grant.
	rtNoRead := newF15Runtime(noReadPlatformGrants(testPrincipal))
	muxNoRead := muxProvider(rtNoRead)
	rec = doMux(t, muxNoRead, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-providers", nil, nil)
	if rec.Code != http.StatusForbidden || rtNoRead.audit.Len() != 1 {
		t.Fatalf("list no-read: %d audit=%d", rec.Code, rtNoRead.audit.Len())
	}

	// Successful create audit-append failure.
	rt3 := newF15Runtime(fullPlatformGrants(testPrincipal))
	seedPlatformForProvider(t, rt3)
	mux3 := muxProvider(rt3)
	rt3.audit.SetFail(errors.New("down"))
	rec = doMux(t, mux3, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers",
		providerCreateBody("prov-x", `["US"]`), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "prov-x",
		})
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("create audit fail status=%d", rec.Code)
	}
	if len(rt3.store.ListCloudProviders()) != 0 {
		t.Fatal("create must not publish on audit failure")
	}

	stages, ok := validate.StageSetFor(validate.KindCloudProvider, validate.RequestClassCollectionCreate)
	if !ok || stages.Reference == nil {
		t.Fatal("provider StageSet required")
	}
}

func doMuxUnauth(t *testing.T, mux http.Handler, method, path string, body []byte, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == nil {
		r = httptest.NewRequest(method, path, nil)
	} else {
		r = httptest.NewRequest(method, path, bytes.NewReader(body))
	}
	for k, v := range headers {
		r.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, r)
	return rec
}
