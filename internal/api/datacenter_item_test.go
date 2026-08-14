package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func seedDatacenter(t *testing.T, rt *f15Runtime, name string) (model.HostingLocation, model.Datacenter) {
	t.Helper()
	hl := seedHostingLocation(t, rt, "loc-for-"+name)
	mux := muxTopology(rt)
	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "seed-dc-" + name}
	rec := doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterCreateBody(name, hl.Metadata.Name, hl.Metadata.UID, ""), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed Datacenter status=%d body=%s", rec.Code, rec.Body.String())
	}
	var dc model.Datacenter
	if err := json.Unmarshal(rec.Body.Bytes(), &dc); err != nil {
		t.Fatalf("seed decode: %v", err)
	}
	return hl, dc
}

func TestDatacenterItem_GetPatchImmutableParent(t *testing.T) {
	rt := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	mux := muxTopology(rt)
	_, dc := seedDatacenter(t, rt, "dc-item")
	auditAfterCreate := rt.audit.Len()

	rec := doMux(t, mux, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters/"+dc.Metadata.UID, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status=%d", rec.Code)
	}

	patchHeaders := map[string]string{
		"Content-Type": mediaMergePatchJSON,
		"If-Match":     dc.Metadata.ResourceVersion,
	}
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters/"+dc.Metadata.UID,
		[]byte(`{"spec":{"description":"updated"}}`), patchHeaders)
	if rec.Code != http.StatusOK {
		t.Fatalf("patch status=%d body=%s", rec.Code, rec.Body.String())
	}
	var updated model.Datacenter
	_ = json.Unmarshal(rec.Body.Bytes(), &updated)
	if updated.Spec.Description != "updated" || rt.audit.Len() != auditAfterCreate+1 {
		t.Fatalf("patched=%#v audit=%d", updated, rt.audit.Len())
	}

	before := rt.audit.Len()
	patchHeaders["If-Match"] = "stale"
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters/"+dc.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), patchHeaders)
	if rec.Code != http.StatusPreconditionFailed || rt.audit.Len() != before {
		t.Fatalf("stale: %d audit=%d", rec.Code, rt.audit.Len())
	}

	before = rt.audit.Len()
	patchHeaders["If-Match"] = updated.Metadata.ResourceVersion
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters/"+dc.Metadata.UID,
		[]byte(`{"metadata":{"name":"renamed"}}`), patchHeaders)
	if rec.Code != http.StatusUnprocessableEntity || !hasViolation(decodeProblem(t, rec), validate.ViolationPatchImmutableField) {
		t.Fatalf("immutable name: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("immutable name must not audit")
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters/"+dc.Metadata.UID,
		[]byte(`{"spec":{"hostingLocationRef":{"apiVersion":"infrastructure.sovrunn.io/v1alpha1","kind":"HostingLocation","name":"x","uid":"11111111111111111111111111111111"}}}`), patchHeaders)
	if rec.Code != http.StatusUnprocessableEntity || !hasViolation(decodeProblem(t, rec), validate.ViolationPatchImmutableField) {
		t.Fatalf("immutable parent: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("immutable parent must not audit")
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters/"+dc.Metadata.UID,
		[]byte(`{"status":{"phase":"Active"}}`), patchHeaders)
	if rec.Code != http.StatusForbidden || !hasViolation(decodeProblem(t, rec), validate.ViolationStatusFieldWrite) {
		t.Fatalf("status: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before+1 {
		t.Fatal("status must audit")
	}

	rtWrong := newF15Runtime(topologyReadOnlyGrants(testPrincipal, testProviderUID))
	seedPlatformAndProvider(t, rtWrong, testProviderUID)
	seeded, prob := rtWrong.store.CreateDatacenter(model.Datacenter{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionDatacenter, Kind: model.KindDatacenter},
		Metadata: apimeta.ObjectMeta{
			Name: "dc-wrong", UID: "22222222222222222222222222222222", ResourceVersion: "1",
			ScopeRef: providerScopeRef(testProviderUID, "prov-a"),
		},
		Spec: model.DatacenterSpec{HostingLocationRef: apimeta.TypedRef{
			APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation,
			Name: "loc", UID: "33333333333333333333333333333333",
		}},
		Status: model.DatacenterStatus{Phase: model.DatacenterPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed: %#v", prob)
	}
	muxWrong := muxTopology(rtWrong)
	before = rtWrong.audit.Len()
	rec = doMux(t, muxWrong, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters/"+seeded.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), map[string]string{
			"Content-Type": mediaMergePatchJSON, "If-Match": "1",
		})
	if rec.Code != http.StatusForbidden || rtWrong.audit.Len() != before {
		t.Fatalf("wrong-admin: %d audit=%d", rec.Code, rtWrong.audit.Len())
	}

	rtForge := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	_, dc2 := seedDatacenter(t, rtForge, "dc-forge")
	muxForge := muxTopology(rtForge)
	before = rtForge.audit.Len()
	rec = doMux(t, muxForge, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters/"+dc2.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), map[string]string{
			"Content-Type": "text/plain", "If-Match": "1",
			cloudmodel.HeaderBootstrapGrant: "forged",
		})
	if rec.Code != http.StatusForbidden || rtForge.audit.Len() != before+1 {
		t.Fatalf("forged before media: %d audit=%d", rec.Code, rtForge.audit.Len())
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPut, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters/"+dc.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PUT %d", rec.Code)
	}
	rec = doMux(t, mux, http.MethodDelete, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters/"+dc.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed || rt.audit.Len() != before {
		t.Fatalf("DELETE %d audit=%d", rec.Code, rt.audit.Len())
	}
	if rt.idemp.InFlightCountForTest() != 0 {
		t.Fatal("PATCH must not use idempotency InFlight state")
	}

	stages, ok := validate.StageSetFor(validate.KindDatacenter, validate.RequestClassPatch)
	if !ok || stages.Defaulting == nil {
		t.Fatal("Datacenter PATCH StageSet required")
	}
}
