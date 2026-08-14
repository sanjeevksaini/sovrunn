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

func TestInfrastructureStackItem_GetPatchImmutableParent(t *testing.T) {
	rt := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	mux := muxTopology(rt)
	_, st := seedInfrastructureStack(t, rt, "stack-item")
	auditAfterCreate := rt.audit.Len()

	rec := doMux(t, mux, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/"+st.Metadata.UID, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status=%d", rec.Code)
	}

	patchHeaders := map[string]string{
		"Content-Type": mediaMergePatchJSON,
		"If-Match":     st.Metadata.ResourceVersion,
	}
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/"+st.Metadata.UID,
		[]byte(`{"spec":{"description":"updated"}}`), patchHeaders)
	if rec.Code != http.StatusOK {
		t.Fatalf("patch status=%d body=%s", rec.Code, rec.Body.String())
	}
	var updated model.InfrastructureStack
	_ = json.Unmarshal(rec.Body.Bytes(), &updated)
	if updated.Spec.Description != "updated" || rt.audit.Len() != auditAfterCreate+1 {
		t.Fatalf("patched=%#v audit=%d", updated, rt.audit.Len())
	}

	before := rt.audit.Len()
	patchHeaders["If-Match"] = updated.Metadata.ResourceVersion
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/"+st.Metadata.UID,
		[]byte(`{"metadata":{"name":"renamed"}}`), patchHeaders)
	if rec.Code != http.StatusUnprocessableEntity || !hasViolation(decodeProblem(t, rec), validate.ViolationPatchImmutableField) {
		t.Fatalf("immutable name: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("immutable name must not audit")
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/"+st.Metadata.UID,
		[]byte(`{"spec":{"faultDomainRef":{"apiVersion":"infrastructure.sovrunn.io/v1alpha1","kind":"FaultDomain","name":"x","uid":"11111111111111111111111111111111"}}}`), patchHeaders)
	if rec.Code != http.StatusUnprocessableEntity || !hasViolation(decodeProblem(t, rec), validate.ViolationPatchImmutableField) {
		t.Fatalf("immutable parent: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("immutable parent must not audit")
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/"+st.Metadata.UID,
		[]byte(`{"status":{"phase":"Active"}}`), patchHeaders)
	if rec.Code != http.StatusForbidden || rt.audit.Len() != before+1 {
		t.Fatalf("status: %d audit=%d", rec.Code, rt.audit.Len())
	}

	rtWrong := newF15Runtime(topologyReadOnlyGrants(testPrincipal, testProviderUID))
	seedPlatformAndProvider(t, rtWrong, testProviderUID)
	seeded, prob := rtWrong.store.CreateInfrastructureStack(model.InfrastructureStack{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionInfrastructureStack, Kind: model.KindInfrastructureStack},
		Metadata: apimeta.ObjectMeta{
			Name: "stack-wrong", UID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa01", ResourceVersion: "1",
			ScopeRef: providerScopeRef(testProviderUID, "prov-a"),
		},
		Spec: model.InfrastructureStackSpec{FaultDomainRef: apimeta.TypedRef{
			APIVersion: model.APIVersionFaultDomain, Kind: model.KindFaultDomain,
			Name: "fd", UID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa02",
		}},
		Status: model.InfrastructureStackStatus{Phase: model.InfrastructureStackPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed: %#v", prob)
	}
	muxWrong := muxTopology(rtWrong)
	before = rtWrong.audit.Len()
	rec = doMux(t, muxWrong, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/"+seeded.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), map[string]string{
			"Content-Type": mediaMergePatchJSON, "If-Match": "1",
		})
	if rec.Code != http.StatusForbidden || rtWrong.audit.Len() != before {
		t.Fatalf("wrong-admin: %d audit=%d", rec.Code, rtWrong.audit.Len())
	}

	rtForge := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	_, st2 := seedInfrastructureStack(t, rtForge, "stack-forge")
	muxForge := muxTopology(rtForge)
	before = rtForge.audit.Len()
	rec = doMux(t, muxForge, http.MethodPatch, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/"+st2.Metadata.UID,
		[]byte(`{"spec":{"description":"x"}}`), map[string]string{
			"Content-Type": "text/plain", "If-Match": "1",
			cloudmodel.HeaderBootstrapGrant: "forged",
		})
	if rec.Code != http.StatusForbidden || rtForge.audit.Len() != before+1 {
		t.Fatalf("forged before media: %d audit=%d", rec.Code, rtForge.audit.Len())
	}

	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPut, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/"+st.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PUT %d", rec.Code)
	}
	rec = doMux(t, mux, http.MethodDelete, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/"+st.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed || rt.audit.Len() != before {
		t.Fatalf("DELETE %d audit=%d", rec.Code, rt.audit.Len())
	}
	if rt.idemp.InFlightCountForTest() != 0 {
		t.Fatal("PATCH must not use idempotency InFlight state")
	}

	stages, ok := validate.StageSetFor(validate.KindInfrastructureStack, validate.RequestClassPatch)
	if !ok || stages.Defaulting == nil {
		t.Fatal("InfrastructureStack PATCH StageSet required")
	}
}
