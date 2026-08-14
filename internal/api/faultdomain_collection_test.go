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

func faultDomainCreateBody(name, dcName, dcUID, desc string) []byte {
	descJSON := ""
	if desc != "" {
		descJSON = `,"description":"` + desc + `"`
	}
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{"datacenterRef":{
			"apiVersion":"infrastructure.sovrunn.io/v1alpha1",
			"kind":"Datacenter",
			"name":"` + dcName + `",
			"uid":"` + dcUID + `"
		}` + descJSON + `}
	}`)
}

func TestFaultDomainCollection_ListCreateParentScopeIsolation(t *testing.T) {
	rt := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	mux := muxTopology(rt)

	rec := doMuxUnauth(t, mux, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains", nil, nil)
	if rec.Code != http.StatusUnauthorized || rt.audit.Len() != 0 {
		t.Fatalf("auth required: %d audit=%d", rec.Code, rt.audit.Len())
	}

	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "fd-root"}
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains",
		faultDomainCreateBody("fd-a", "dc", "dddddddddddddddddddddddddddddddd", ""), headers)
	if rec.Code != http.StatusConflict || rt.audit.Len() != 1 {
		t.Fatalf("root gate: %d audit=%d", rec.Code, rt.audit.Len())
	}

	_, dc := seedDatacenter(t, rt, "dc-for-fd")
	before := rt.audit.Len()
	headers["Idempotency-Key"] = "fd-1"
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains",
		faultDomainCreateBody("fd-a", dc.Metadata.Name, dc.Metadata.UID, "zone-a"), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created model.FaultDomain
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	if created.Status.Phase != model.FaultDomainPhaseActive || created.Metadata.ScopeRef.UID != testProviderUID {
		t.Fatalf("created=%#v", created)
	}
	if rt.audit.Len() != before+1 {
		t.Fatalf("create audit=%d", rt.audit.Len())
	}

	rtCross := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	seedPlatformAndProvider(t, rtCross, testProviderUID)
	seedOtherProvider(t, rtCross, testOtherProviderUID)
	otherDC, prob := rtCross.store.CreateDatacenter(model.Datacenter{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionDatacenter, Kind: model.KindDatacenter},
		Metadata: apimeta.ObjectMeta{
			Name: "dc-other", UID: "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", ResourceVersion: "1",
			ScopeRef: providerScopeRef(testOtherProviderUID, "prov-b"),
		},
		Spec: model.DatacenterSpec{HostingLocationRef: apimeta.TypedRef{
			APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation,
			Name: "loc", UID: "ffffffffffffffffffffffffffffffff",
		}},
		Status: model.DatacenterStatus{Phase: model.DatacenterPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed other DC: %#v", prob)
	}
	muxCross := muxTopology(rtCross)
	before = rtCross.audit.Len()
	rec = doMux(t, muxCross, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains",
		faultDomainCreateBody("fd-x", otherDC.Metadata.Name, otherDC.Metadata.UID, ""), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "fd-inacc",
		})
	if rec.Code != http.StatusNotFound || rtCross.audit.Len() != before+1 {
		t.Fatalf("inaccessible: %d audit=%d", rec.Code, rtCross.audit.Len())
	}

	rtMismatch := newF15Runtime(topologyCrossProviderGrants(testPrincipal, testProviderUID, testOtherProviderUID))
	seedPlatformAndProvider(t, rtMismatch, testProviderUID)
	seedOtherProvider(t, rtMismatch, testOtherProviderUID)
	visDC, prob := rtMismatch.store.CreateDatacenter(model.Datacenter{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionDatacenter, Kind: model.KindDatacenter},
		Metadata: apimeta.ObjectMeta{
			Name: "dc-vis", UID: "99999999999999999999999999999999", ResourceVersion: "1",
			ScopeRef: providerScopeRef(testOtherProviderUID, "prov-b"),
		},
		Spec: model.DatacenterSpec{HostingLocationRef: apimeta.TypedRef{
			APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation,
			Name: "loc", UID: "88888888888888888888888888888888",
		}},
		Status: model.DatacenterStatus{Phase: model.DatacenterPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed visible DC: %#v", prob)
	}
	muxMismatch := muxTopology(rtMismatch)
	before = rtMismatch.audit.Len()
	rec = doMux(t, muxMismatch, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains",
		faultDomainCreateBody("fd-mm", visDC.Metadata.Name, visDC.Metadata.UID, ""), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "fd-mm",
		})
	if rec.Code != http.StatusUnprocessableEntity || !hasViolation(decodeProblem(t, rec), validate.ViolationTopologyProviderMismatch) {
		t.Fatalf("visible mismatch: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rtMismatch.audit.Len() != before {
		t.Fatal("visible mismatch must not audit")
	}

	before = rt.audit.Len()
	headers["Idempotency-Key"] = "fd-dup"
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains",
		faultDomainCreateBody("fd-a", dc.Metadata.Name, dc.Metadata.UID, ""), headers)
	if decodeProblem(t, rec).Code != apiproblem.CodeAlreadyExists || rt.audit.Len() != before {
		t.Fatalf("duplicate: %d audit=%d", rec.Code, rt.audit.Len())
	}

	rtNoRead := newF15Runtime(topologyWriteOnlyGrants(testPrincipal, testProviderUID))
	seedPlatformAndProvider(t, rtNoRead, testProviderUID)
	muxNoRead := muxTopology(rtNoRead)
	rec = doMux(t, muxNoRead, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains", nil, nil)
	if rec.Code != http.StatusForbidden || rtNoRead.audit.Len() != 1 {
		t.Fatalf("list no-read: %d audit=%d", rec.Code, rtNoRead.audit.Len())
	}

	rec = doMux(t, mux, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status=%d", rec.Code)
	}

	stages, ok := validate.StageSetFor(validate.KindFaultDomain, validate.RequestClassCollectionCreate)
	if !ok || stages.Reference == nil {
		t.Fatal("FaultDomain StageSet required")
	}
}

func seedFaultDomain(t *testing.T, rt *f15Runtime, name string) (model.Datacenter, model.FaultDomain) {
	t.Helper()
	_, dc := seedDatacenter(t, rt, "dc-for-"+name)
	mux := muxTopology(rt)
	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "seed-fd-" + name}
	rec := doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/fault-domains",
		faultDomainCreateBody(name, dc.Metadata.Name, dc.Metadata.UID, ""), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed FaultDomain status=%d body=%s", rec.Code, rec.Body.String())
	}
	var fd model.FaultDomain
	if err := json.Unmarshal(rec.Body.Bytes(), &fd); err != nil {
		t.Fatalf("seed decode: %v", err)
	}
	return dc, fd
}
