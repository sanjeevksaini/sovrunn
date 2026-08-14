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

func TestFaultDomainItem_GetPatchImmutableParent(t *testing.T) {
	rt := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	mux := muxTopology(rt)
	_, fd := seedFaultDomain(t, rt, "fd-item")
	auditAfterCreate := rt.audit.Len()

	rec := doMux(t, mux, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/"+fd.Metadata.UID, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status=%d", rec.Code)
	}

	patchHeaders := map[string]string{
		"Content-Type": mediaMergePatchJSON,
		"If-Match":     fd.Metadata.ResourceVersion,
	}
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/"+fd.Metadata.UID,
		[]byte(`{"spec":{"description":"updated"}}`), patchHeaders)
	if rec.Code != http.StatusOK {
		t.Fatalf("patch status=%d body=%s", rec.Code, rec.Body.String())
	}
	var updated model.FaultDomain
	_ = json.Unmarshal(rec.Body.Bytes(), &updated)
	if updated.Spec.Description != "updated" || rt.audit.Len() != auditAfterCreate+1 {
		t.Fatalf("patched=%#v audit=%d", updated, rt.audit.Len())
	}

	before := rt.audit.Len()
	patchHeaders["If-Match"] = updated.Metadata.ResourceVersion
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/"+fd.Metadata.UID,
		[]byte(`{"metadata":{"name":"renamed"}}`), patchHeaders)
	if rec.Code != http.StatusUnprocessableEntity || !hasViolation(decodeProblem(t, rec), validate.ViolationPatchImmutableField) {
		t.Fatalf("immutable name: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("immutable name must not audit")
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/"+fd.Metadata.UID,
		[]byte(`{"spec":{"datacenterRef":{"apiVersion":"infrastructure.sovrunn.io/v1alpha1","kind":"Datacenter","name":"x","uid":"11111111111111111111111111111111"}}}`), patchHeaders)
	if rec.Code != http.StatusUnprocessableEntity || !hasViolation(decodeProblem(t, rec), validate.ViolationPatchImmutableField) {
		t.Fatalf("immutable parent: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("immutable parent must not audit")
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/"+fd.Metadata.UID,
		[]byte(`{"status":{"phase":"Active"}}`), patchHeaders)
	if rec.Code != http.StatusForbidden || rt.audit.Len() != before+1 {
		t.Fatalf("status: %d audit=%d", rec.Code, rt.audit.Len())
	}

	rtWrong := newF15Runtime(topologyReadOnlyGrants(testPrincipal, testProviderUID))
	seedPlatformAndProvider(t, rtWrong, testProviderUID)
	seeded, prob := rtWrong.store.CreateFaultDomain(model.FaultDomain{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionFaultDomain, Kind: model.KindFaultDomain},
		Metadata: apimeta.ObjectMeta{
			Name: "fd-wrong", UID: "44444444444444444444444444444444", ResourceVersion: "1",
			ScopeRef: providerScopeRef(testProviderUID, "prov-a"),
		},
		Spec: model.FaultDomainSpec{DatacenterRef: apimeta.TypedRef{
			APIVersion: model.APIVersionDatacenter, Kind: model.KindDatacenter,
			Name: "dc", UID: "55555555555555555555555555555555",
		}},
		Status: model.FaultDomainStatus{Phase: model.FaultDomainPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed: %#v", prob)
	}
	muxWrong := muxTopology(rtWrong)
	before = rtWrong.audit.Len()
	rec = doMux(t, muxWrong, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/"+seeded.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), map[string]string{
			"Content-Type": mediaMergePatchJSON, "If-Match": "1",
		})
	if rec.Code != http.StatusForbidden || rtWrong.audit.Len() != before {
		t.Fatalf("wrong-admin: %d audit=%d", rec.Code, rtWrong.audit.Len())
	}

	rtForge := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	_, fd2 := seedFaultDomain(t, rtForge, "fd-forge")
	muxForge := muxTopology(rtForge)
	before = rtForge.audit.Len()
	rec = doMux(t, muxForge, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/"+fd2.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), map[string]string{
			"Content-Type": "text/plain", "If-Match": "1",
			cloudmodel.HeaderBootstrapGrant: "forged",
		})
	if rec.Code != http.StatusForbidden || rtForge.audit.Len() != before+1 {
		t.Fatalf("forged before media: %d audit=%d", rec.Code, rtForge.audit.Len())
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPut, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/"+fd.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PUT %d", rec.Code)
	}
	rec = doMux(t, mux, http.MethodDelete, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/"+fd.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed || rt.audit.Len() != before {
		t.Fatalf("DELETE %d audit=%d", rec.Code, rt.audit.Len())
	}
	if rt.idemp.InFlightCountForTest() != 0 {
		t.Fatal("PATCH must not use idempotency InFlight state")
	}

	stages, ok := validate.StageSetFor(validate.KindFaultDomain, validate.RequestClassPatch)
	if !ok || stages.Defaulting == nil {
		t.Fatal("FaultDomain PATCH StageSet required")
	}
}
