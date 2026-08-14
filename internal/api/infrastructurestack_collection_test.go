package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func infrastructureStackCreateBody(name, fdName, fdUID, desc string) []byte {
	descJSON := ""
	if desc != "" {
		descJSON = `,"description":"` + desc + `"`
	}
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{"faultDomainRef":{
			"apiVersion":"infrastructure.sovrunn.io/v1alpha1",
			"kind":"FaultDomain",
			"name":"` + fdName + `",
			"uid":"` + fdUID + `"
		}` + descJSON + `}
	}`)
}

func TestInfrastructureStackCollection_ListCreateParentScopeIsolation(t *testing.T) {
	rt := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	mux := muxTopology(rt)

	rec := doMuxUnauth(t, mux, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks", nil, nil)
	if rec.Code != http.StatusUnauthorized || rt.audit.Len() != 0 {
		t.Fatalf("auth required: %d audit=%d", rec.Code, rt.audit.Len())
	}

	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "stack-root"}
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks",
		infrastructureStackCreateBody("stack-a", "fd", "dddddddddddddddddddddddddddddddd", ""), headers)
	if rec.Code != http.StatusConflict || rt.audit.Len() != 1 {
		t.Fatalf("root gate: %d audit=%d", rec.Code, rt.audit.Len())
	}

	_, fd := seedFaultDomain(t, rt, "fd-for-stack")
	before := rt.audit.Len()
	headers["Idempotency-Key"] = "stack-1"
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks",
		infrastructureStackCreateBody("stack-a", fd.Metadata.Name, fd.Metadata.UID, "k8s"), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created model.InfrastructureStack
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	if created.Status.Phase != model.InfrastructureStackPhaseActive || created.Metadata.ScopeRef.UID != testProviderUID {
		t.Fatalf("created=%#v", created)
	}
	if rt.audit.Len() != before+1 {
		t.Fatalf("create audit=%d", rt.audit.Len())
	}

	rtCross := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	seedPlatformAndProvider(t, rtCross, testProviderUID)
	seedOtherProvider(t, rtCross, testOtherProviderUID)
	otherFD, prob := rtCross.store.CreateFaultDomain(model.FaultDomain{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionFaultDomain, Kind: model.KindFaultDomain},
		Metadata: apimeta.ObjectMeta{
			Name: "fd-other", UID: "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", ResourceVersion: "1",
			ScopeRef: providerScopeRef(testOtherProviderUID, "prov-b"),
		},
		Spec: model.FaultDomainSpec{DatacenterRef: apimeta.TypedRef{
			APIVersion: model.APIVersionDatacenter, Kind: model.KindDatacenter,
			Name: "dc", UID: "ffffffffffffffffffffffffffffffff",
		}},
		Status: model.FaultDomainStatus{Phase: model.FaultDomainPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed other FD: %#v", prob)
	}
	muxCross := muxTopology(rtCross)
	before = rtCross.audit.Len()
	rec = doMux(t, muxCross, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks",
		infrastructureStackCreateBody("stack-x", otherFD.Metadata.Name, otherFD.Metadata.UID, ""), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "stack-inacc",
		})
	if rec.Code != http.StatusNotFound || rtCross.audit.Len() != before+1 {
		t.Fatalf("inaccessible: %d audit=%d", rec.Code, rtCross.audit.Len())
	}

	rtMismatch := newF15Runtime(topologyCrossProviderGrants(testPrincipal, testProviderUID, testOtherProviderUID))
	seedPlatformAndProvider(t, rtMismatch, testProviderUID)
	seedOtherProvider(t, rtMismatch, testOtherProviderUID)
	visFD, prob := rtMismatch.store.CreateFaultDomain(model.FaultDomain{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionFaultDomain, Kind: model.KindFaultDomain},
		Metadata: apimeta.ObjectMeta{
			Name: "fd-vis", UID: "77777777777777777777777777777777", ResourceVersion: "1",
			ScopeRef: providerScopeRef(testOtherProviderUID, "prov-b"),
		},
		Spec: model.FaultDomainSpec{DatacenterRef: apimeta.TypedRef{
			APIVersion: model.APIVersionDatacenter, Kind: model.KindDatacenter,
			Name: "dc", UID: "66666666666666666666666666666666",
		}},
		Status: model.FaultDomainStatus{Phase: model.FaultDomainPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed visible FD: %#v", prob)
	}
	muxMismatch := muxTopology(rtMismatch)
	before = rtMismatch.audit.Len()
	rec = doMux(t, muxMismatch, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks",
		infrastructureStackCreateBody("stack-mm", visFD.Metadata.Name, visFD.Metadata.UID, ""), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "stack-mm",
		})
	if rec.Code != http.StatusUnprocessableEntity || !hasViolation(decodeProblem(t, rec), validate.ViolationTopologyProviderMismatch) {
		t.Fatalf("visible mismatch: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rtMismatch.audit.Len() != before {
		t.Fatal("visible mismatch must not audit")
	}

	before = rt.audit.Len()
	headers["Idempotency-Key"] = "stack-dup"
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks",
		infrastructureStackCreateBody("stack-a", fd.Metadata.Name, fd.Metadata.UID, ""), headers)
	if decodeProblem(t, rec).Code != apiproblem.CodeAlreadyExists || rt.audit.Len() != before {
		t.Fatalf("duplicate: %d audit=%d", rec.Code, rt.audit.Len())
	}

	rtNoRead := newF15Runtime(topologyWriteOnlyGrants(testPrincipal, testProviderUID))
	seedPlatformAndProvider(t, rtNoRead, testProviderUID)
	muxNoRead := muxTopology(rtNoRead)
	rec = doMux(t, muxNoRead, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks", nil, nil)
	if rec.Code != http.StatusForbidden || rtNoRead.audit.Len() != 1 {
		t.Fatalf("list no-read: %d audit=%d", rec.Code, rtNoRead.audit.Len())
	}

	rec = doMux(t, mux, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status=%d", rec.Code)
	}

	stages, ok := validate.StageSetFor(validate.KindInfrastructureStack, validate.RequestClassCollectionCreate)
	if !ok || stages.Reference == nil {
		t.Fatal("InfrastructureStack StageSet required")
	}
}

func seedInfrastructureStack(t *testing.T, rt *f15Runtime, name string) (model.FaultDomain, model.InfrastructureStack) {
	t.Helper()
	_, fd := seedFaultDomain(t, rt, "fd-for-"+name)
	mux := muxTopology(rt)
	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "seed-stack-" + name}
	rec := doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks",
		infrastructureStackCreateBody(name, fd.Metadata.Name, fd.Metadata.UID, ""), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed InfrastructureStack status=%d body=%s", rec.Code, rec.Body.String())
	}
	var st model.InfrastructureStack
	if err := json.Unmarshal(rec.Body.Bytes(), &st); err != nil {
		t.Fatalf("seed decode: %v", err)
	}
	return fd, st
}
