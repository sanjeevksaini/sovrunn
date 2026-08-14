package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func seedCloudProvider(t *testing.T, rt *f15Runtime, name string) model.CloudProvider {
	t.Helper()
	seedPlatformForProvider(t, rt)
	mux := muxProvider(rt)
	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "seed-prov-" + name}
	rec := doMux(t, mux, http.MethodPost, "/apis/core.sovrunn.io/v1alpha1/cloud-providers",
		providerCreateBody(name, `["US"]`), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed provider status=%d body=%s", rec.Code, rec.Body.String())
	}
	var cp model.CloudProvider
	if err := json.Unmarshal(rec.Body.Bytes(), &cp); err != nil {
		t.Fatalf("seed decode: %v", err)
	}
	return cp
}

func TestCloudProviderItem_GetPatchWriterMethods(t *testing.T) {
	rt := newF15Runtime(fullPlatformGrants(testPrincipal))
	mux := muxProvider(rt)
	cp := seedCloudProvider(t, rt, "prov-item")
	auditAfterCreate := rt.audit.Len()

	rec := doMux(t, mux, http.MethodGet, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+cp.Metadata.UID, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status=%d", rec.Code)
	}

	patchHeaders := map[string]string{
		"Content-Type": mediaMergePatchJSON,
		"If-Match":     cp.Metadata.ResourceVersion,
	}
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+cp.Metadata.UID,
		[]byte(`{"spec":{"displayName":"Acme Cloud"}}`), patchHeaders)
	if rec.Code != http.StatusOK {
		t.Fatalf("patch displayName status=%d body=%s", rec.Code, rec.Body.String())
	}
	var updated model.CloudProvider
	_ = json.Unmarshal(rec.Body.Bytes(), &updated)
	if updated.Spec.DisplayName != "Acme Cloud" || rt.audit.Len() != auditAfterCreate+1 {
		t.Fatalf("patched=%#v audit=%d", updated, rt.audit.Len())
	}

	patchHeaders["If-Match"] = updated.Metadata.ResourceVersion
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+cp.Metadata.UID,
		[]byte(`{"spec":{"operatingMarkets":["US","IN"]}}`), patchHeaders)
	if rec.Code != http.StatusOK {
		t.Fatalf("patch markets status=%d body=%s", rec.Code, rec.Body.String())
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &updated)

	before := rt.audit.Len()
	patchHeaders["If-Match"] = "stale"
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+cp.Metadata.UID,
		[]byte(`{"spec":{"displayName":"x"}}`), patchHeaders)
	if rec.Code != http.StatusPreconditionFailed || rt.audit.Len() != before {
		t.Fatalf("stale: %d audit=%d", rec.Code, rt.audit.Len())
	}

	delete(patchHeaders, "If-Match")
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+cp.Metadata.UID,
		[]byte(`{"spec":{"displayName":"x"}}`), patchHeaders)
	if rec.Code != http.StatusPreconditionFailed {
		t.Fatalf("missing If-Match status=%d", rec.Code)
	}
	patchHeaders["If-Match"] = "bad\x01"
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+cp.Metadata.UID,
		[]byte(`{"spec":{"displayName":"x"}}`), patchHeaders)
	if rec.Code != http.StatusPreconditionFailed {
		t.Fatalf("malformed If-Match status=%d", rec.Code)
	}

	before = rt.audit.Len()
	patchHeaders["If-Match"] = updated.Metadata.ResourceVersion
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+cp.Metadata.UID,
		[]byte(`{"metadata":{"name":"renamed"}}`), patchHeaders)
	if rec.Code != http.StatusUnprocessableEntity || !hasViolation(decodeProblem(t, rec), validate.ViolationPatchImmutableField) {
		t.Fatalf("immutable name: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("immutable must not audit")
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+cp.Metadata.UID,
		[]byte(`{"status":{"phase":"Active"}}`), patchHeaders)
	if rec.Code != http.StatusForbidden || !hasViolation(decodeProblem(t, rec), validate.ViolationStatusFieldWrite) {
		t.Fatalf("status: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before+1 {
		t.Fatal("status must audit")
	}

	rtWrong := newF15Runtime(readOnlyPlatformGrants(testPrincipal))
	seedPlatformForProvider(t, rtWrong)
	seeded, prob := rtWrong.store.CreateCloudProvider(model.CloudProvider{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider},
		Metadata: apimeta.ObjectMeta{Name: "prov-wrong", UID: "11111111111111111111111111111111", ResourceVersion: "1"},
		Spec:     model.CloudProviderSpec{OperatingMarkets: []string{"US"}},
		Status:   model.CloudProviderStatus{Phase: model.CloudProviderPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed: %#v", prob)
	}
	muxWrong := muxProvider(rtWrong)
	before = rtWrong.audit.Len()
	rec = doMux(t, muxWrong, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+seeded.Metadata.UID,
		[]byte(`{"spec":{"displayName":"x"}}`), map[string]string{
			"Content-Type": mediaMergePatchJSON, "If-Match": "1",
		})
	if rec.Code != http.StatusForbidden || rtWrong.audit.Len() != before {
		t.Fatalf("wrong-admin: %d audit=%d", rec.Code, rtWrong.audit.Len())
	}

	rtForge := newF15Runtime(fullPlatformGrants(testPrincipal))
	cp2 := seedCloudProvider(t, rtForge, "prov-forge")
	muxForge := muxProvider(rtForge)
	before = rtForge.audit.Len()
	rec = doMux(t, muxForge, http.MethodPatch, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+cp2.Metadata.UID,
		[]byte(`{"spec":{"displayName":"x"}}`), map[string]string{
			"Content-Type": "text/plain", "If-Match": "1",
			cloudmodel.HeaderBootstrapGrant: "forged",
		})
	if rec.Code != http.StatusForbidden || rtForge.audit.Len() != before+1 {
		t.Fatalf("forged before media: %d audit=%d", rec.Code, rtForge.audit.Len())
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPut, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+cp.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PUT %d", rec.Code)
	}
	rec = doMux(t, mux, http.MethodDelete, "/apis/core.sovrunn.io/v1alpha1/cloud-providers/"+cp.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed || rt.audit.Len() != before {
		t.Fatalf("DELETE %d audit=%d", rec.Code, rt.audit.Len())
	}
	if rt.idemp.InFlightCountForTest() != 0 {
		t.Fatal("PATCH must not use idempotency InFlight state")
	}

	stages, ok := validate.StageSetFor(validate.KindCloudProvider, validate.RequestClassPatch)
	if !ok || stages.Defaulting == nil {
		t.Fatal("provider PATCH StageSet required")
	}
	_ = apiproblem.CodeStaleResourceVersion
}
