package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func datacenterCreateBody(name, hlName, hlUID, desc string) []byte {
	descJSON := ""
	if desc != "" {
		descJSON = `,"description":"` + desc + `"`
	}
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{"hostingLocationRef":{
			"apiVersion":"infrastructure.sovrunn.io/v1alpha1",
			"kind":"HostingLocation",
			"name":"` + hlName + `",
			"uid":"` + hlUID + `"
		}` + descJSON + `}
	}`)
}

func TestDatacenterCollection_ListCreateParentScopeIsolation(t *testing.T) {
	rt := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	mux := muxTopology(rt)

	rec := doMuxUnauth(t, mux, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters", nil, nil)
	if rec.Code != http.StatusUnauthorized || rt.audit.Len() != 0 {
		t.Fatalf("auth required: %d audit=%d", rec.Code, rt.audit.Len())
	}

	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "dc-root"}
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterCreateBody("dc-a", "loc", "dddddddddddddddddddddddddddddddd", ""), headers)
	if rec.Code != http.StatusConflict || !hasViolation(decodeProblem(t, rec), violationCloudPlatformRootRequired) {
		t.Fatalf("root gate: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != 1 {
		t.Fatalf("root gate audit=%d", rt.audit.Len())
	}

	hl := seedHostingLocation(t, rt, "loc-dc")
	before := rt.audit.Len()
	headers["Idempotency-Key"] = "dc-1"
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterCreateBody("dc-a", hl.Metadata.Name, hl.Metadata.UID, "primary"), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created model.Datacenter
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	if created.Status.Phase != model.DatacenterPhaseActive {
		t.Fatalf("status=%#v", created.Status)
	}
	if created.Metadata.ScopeRef == nil || created.Metadata.ScopeRef.UID != testProviderUID {
		t.Fatalf("scope from parent: %#v", created.Metadata.ScopeRef)
	}
	if rt.audit.Len() != before+1 {
		t.Fatalf("create audit=%d", rt.audit.Len())
	}

	// Inaccessible cross-provider parent → audited safe 404.
	rtCross := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	seedPlatformAndProvider(t, rtCross, testProviderUID)
	seedOtherProvider(t, rtCross, testOtherProviderUID)
	otherHL, prob := rtCross.store.CreateHostingLocation(model.HostingLocation{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation},
		Metadata: apimeta.ObjectMeta{
			Name: "loc-other", UID: "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", ResourceVersion: "1",
			ScopeRef: providerScopeRef(testOtherProviderUID, "prov-b"),
		},
		Spec:   model.HostingLocationSpec{CountryCode: "IN", Locality: "Mumbai"},
		Status: model.HostingLocationStatus{Phase: model.HostingLocationPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed other HL: %#v", prob)
	}
	muxCross := muxTopology(rtCross)
	before = rtCross.audit.Len()
	rec = doMux(t, muxCross, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterCreateBody("dc-x", otherHL.Metadata.Name, otherHL.Metadata.UID, ""), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "dc-inacc",
		})
	if rec.Code != http.StatusNotFound || !hasViolation(decodeProblem(t, rec), violationAuthorizationSafeDenial) {
		t.Fatalf("inaccessible: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rtCross.audit.Len() != before+1 {
		t.Fatalf("inaccessible audit=%d", rtCross.audit.Len())
	}

	// Visible parent without matching topology.write → VS0_TOPOLOGY_PROVIDER_MISMATCH, no audit.
	rtMismatch := newF15Runtime(topologyCrossProviderGrants(testPrincipal, testProviderUID, testOtherProviderUID))
	seedPlatformAndProvider(t, rtMismatch, testProviderUID)
	seedOtherProvider(t, rtMismatch, testOtherProviderUID)
	otherHL2, prob := rtMismatch.store.CreateHostingLocation(model.HostingLocation{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionHostingLocation, Kind: model.KindHostingLocation},
		Metadata: apimeta.ObjectMeta{
			Name: "loc-vis", UID: "ffffffffffffffffffffffffffffffff", ResourceVersion: "1",
			ScopeRef: providerScopeRef(testOtherProviderUID, "prov-b"),
		},
		Spec:   model.HostingLocationSpec{CountryCode: "IN", Locality: "Pune"},
		Status: model.HostingLocationStatus{Phase: model.HostingLocationPhaseActive},
	})
	if prob != nil {
		t.Fatalf("seed visible HL: %#v", prob)
	}
	muxMismatch := muxTopology(rtMismatch)
	before = rtMismatch.audit.Len()
	rec = doMux(t, muxMismatch, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterCreateBody("dc-mm", otherHL2.Metadata.Name, otherHL2.Metadata.UID, ""), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "dc-mm",
		})
	if rec.Code != http.StatusUnprocessableEntity || !hasViolation(decodeProblem(t, rec), validate.ViolationTopologyProviderMismatch) {
		t.Fatalf("visible mismatch: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rtMismatch.audit.Len() != before {
		t.Fatal("visible mismatch must not audit")
	}

	before = rt.audit.Len()
	headers["Idempotency-Key"] = "dc-dup"
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterCreateBody("dc-a", hl.Metadata.Name, hl.Metadata.UID, ""), headers)
	if decodeProblem(t, rec).Code != apiproblem.CodeAlreadyExists || rt.audit.Len() != before {
		t.Fatalf("duplicate: %d audit=%d", rec.Code, rt.audit.Len())
	}

	headers["Idempotency-Key"] = "dc-1"
	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterCreateBody("dc-a", hl.Metadata.Name, hl.Metadata.UID, "primary"), headers)
	if rec.Code != http.StatusCreated || rt.audit.Len() != before {
		t.Fatalf("replay: %d audit=%d", rec.Code, rt.audit.Len())
	}

	before = rt.audit.Len()
	headers["Idempotency-Key"] = "dc-status"
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters", []byte(`{
		"metadata":{"name":"dc-status"},
		"spec":{"hostingLocationRef":{"apiVersion":"infrastructure.sovrunn.io/v1alpha1","kind":"HostingLocation","name":"`+hl.Metadata.Name+`","uid":"`+hl.Metadata.UID+`"}},
		"status":{"phase":"Active"}
	}`), headers)
	if rec.Code != http.StatusForbidden || rt.audit.Len() != before+1 {
		t.Fatalf("status: %d audit=%d", rec.Code, rt.audit.Len())
	}

	rtNoRead := newF15Runtime(topologyWriteOnlyGrants(testPrincipal, testProviderUID))
	seedPlatformAndProvider(t, rtNoRead, testProviderUID)
	muxNoRead := muxTopology(rtNoRead)
	rec = doMux(t, muxNoRead, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters", nil, nil)
	if rec.Code != http.StatusForbidden || rtNoRead.audit.Len() != 1 {
		t.Fatalf("list no-read: %d audit=%d", rec.Code, rtNoRead.audit.Len())
	}

	rec = doMux(t, mux, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status=%d", rec.Code)
	}
	var list apimeta.ListEnvelope[model.Datacenter]
	_ = json.Unmarshal(rec.Body.Bytes(), &list)
	if len(list.Items) != 1 {
		t.Fatalf("list items=%d", len(list.Items))
	}

	rtFail := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	hlFail := seedHostingLocation(t, rtFail, "loc-fail-dc")
	muxFail := muxTopology(rtFail)
	rtFail.audit.SetFail(errors.New("down"))
	rec = doMux(t, muxFail, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/datacenters",
		datacenterCreateBody("dc-fail", hlFail.Metadata.Name, hlFail.Metadata.UID, ""), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "dc-fail",
		})
	if rec.Code != http.StatusInternalServerError || len(rtFail.store.ListDatacenters()) != 0 {
		t.Fatalf("audit fail: %d published=%d", rec.Code, len(rtFail.store.ListDatacenters()))
	}

	stages, ok := validate.StageSetFor(validate.KindDatacenter, validate.RequestClassCollectionCreate)
	if !ok || stages.Reference == nil {
		t.Fatal("Datacenter StageSet required")
	}
}
