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

func seedHostingLocation(t *testing.T, rt *f15Runtime, name string) model.HostingLocation {
	t.Helper()
	seedPlatformAndProvider(t, rt, testProviderUID)
	mux := muxTopology(rt)
	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "seed-hl-" + name}
	rec := doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations",
		hostingLocationCreateBody(name, "US", "Austin", "US-TX", ""), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed HostingLocation status=%d body=%s", rec.Code, rec.Body.String())
	}
	var hl model.HostingLocation
	if err := json.Unmarshal(rec.Body.Bytes(), &hl); err != nil {
		t.Fatalf("seed decode: %v", err)
	}
	return hl
}

func TestHostingLocationItem_GetPatchWriterMethods(t *testing.T) {
	rt := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	mux := muxTopology(rt)
	hl := seedHostingLocation(t, rt, "loc-item")
	auditAfterCreate := rt.audit.Len()

	rec := doMux(t, mux, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+hl.Metadata.UID, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status=%d", rec.Code)
	}

	patchHeaders := map[string]string{
		"Content-Type": mediaMergePatchJSON,
		"If-Match":     hl.Metadata.ResourceVersion,
	}
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+hl.Metadata.UID,
		[]byte(`{"spec":{"description":"updated"}}`), patchHeaders)
	if rec.Code != http.StatusOK {
		t.Fatalf("patch description status=%d body=%s", rec.Code, rec.Body.String())
	}
	var updated model.HostingLocation
	_ = json.Unmarshal(rec.Body.Bytes(), &updated)
	if updated.Spec.Description != "updated" || rt.audit.Len() != auditAfterCreate+1 {
		t.Fatalf("patched=%#v audit=%d", updated, rt.audit.Len())
	}

	before := rt.audit.Len()
	patchHeaders["If-Match"] = "stale"
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+hl.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), patchHeaders)
	if rec.Code != http.StatusPreconditionFailed || rt.audit.Len() != before {
		t.Fatalf("stale: %d audit=%d", rec.Code, rt.audit.Len())
	}

	delete(patchHeaders, "If-Match")
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+hl.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), patchHeaders)
	if rec.Code != http.StatusPreconditionFailed {
		t.Fatalf("missing If-Match status=%d", rec.Code)
	}
	patchHeaders["If-Match"] = "bad\x01"
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+hl.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), patchHeaders)
	if rec.Code != http.StatusPreconditionFailed {
		t.Fatalf("malformed If-Match status=%d", rec.Code)
	}

	before = rt.audit.Len()
	patchHeaders["If-Match"] = updated.Metadata.ResourceVersion
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+hl.Metadata.UID,
		[]byte(`{"metadata":{"name":"renamed"}}`), patchHeaders)
	if rec.Code != http.StatusUnprocessableEntity || !hasViolation(decodeProblem(t, rec), validate.ViolationPatchImmutableField) {
		t.Fatalf("immutable name: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("immutable must not audit")
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+hl.Metadata.UID,
		[]byte(`{"status":{"phase":"Active"}}`), patchHeaders)
	if rec.Code != http.StatusForbidden || !hasViolation(decodeProblem(t, rec), validate.ViolationStatusFieldWrite) {
		t.Fatalf("status: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before+1 {
		t.Fatal("status must audit")
	}

	rtWrong := newF15Runtime(topologyReadOnlyGrants(testPrincipal, testProviderUID))
	seedPlatformAndProvider(t, rtWrong, testProviderUID)
	seeded, prob := rtWrong.store.CreateHostingLocation(model.HostingLocation{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation},
		Metadata: apimeta.ObjectMeta{
			Name: "loc-wrong", UID: "11111111111111111111111111111111", ResourceVersion: "1",
			ScopeRef: providerScopeRef(testProviderUID, "prov-a"),
		},
		Spec:   model.HostingLocationSpec{CountryCode: "US", Locality: "Austin"},
		Status: model.HostingLocationStatus{Phase: model.HostingLocationPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed: %#v", prob)
	}
	muxWrong := muxTopology(rtWrong)
	before = rtWrong.audit.Len()
	rec = doMux(t, muxWrong, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+seeded.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), map[string]string{
			"Content-Type": mediaMergePatchJSON, "If-Match": "1",
		})
	if rec.Code != http.StatusForbidden || rtWrong.audit.Len() != before {
		t.Fatalf("wrong-admin: %d audit=%d", rec.Code, rtWrong.audit.Len())
	}

	rtForge := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	hl2 := seedHostingLocation(t, rtForge, "loc-forge")
	muxForge := muxTopology(rtForge)
	before = rtForge.audit.Len()
	rec = doMux(t, muxForge, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+hl2.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), map[string]string{
			"Content-Type": "text/plain", "If-Match": "1",
			cloudmodel.HeaderBootstrapGrant: "forged",
		})
	if rec.Code != http.StatusForbidden || rtForge.audit.Len() != before+1 {
		t.Fatalf("forged before media: %d audit=%d", rec.Code, rtForge.audit.Len())
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPut, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+hl.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PUT %d", rec.Code)
	}
	rec = doMux(t, mux, http.MethodDelete, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/"+hl.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed || rt.audit.Len() != before {
		t.Fatalf("DELETE %d audit=%d", rec.Code, rt.audit.Len())
	}
	if rt.idemp.InFlightCountForTest() != 0 {
		t.Fatal("PATCH must not use idempotency InFlight state")
	}

	stages, ok := validate.StageSetFor(validate.KindHostingLocation, validate.RequestClassPatch)
	if !ok || stages.Defaulting == nil {
		t.Fatal("HostingLocation PATCH StageSet required")
	}
	_ = apiproblem.CodeStaleResourceVersion
}
