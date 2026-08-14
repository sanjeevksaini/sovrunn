package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

const (
	testProviderUID      = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	testOtherProviderUID = "cccccccccccccccccccccccccccccccc"
	testPlatformUID      = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
)

func providerScopeID(uid string) apimeta.ScopeIdentity {
	return apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: uid}
}

func topologyGrants(principal, providerUID string) *staticGrants {
	return &staticGrants{
		principal: principal,
		grants: []cloudmodel.Grant{
			{Action: cloudmodel.ActionCloudPlatformWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderRead, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionTopologyWrite, Scope: providerScopeID(providerUID)},
			{Action: cloudmodel.ActionTopologyRead, Scope: providerScopeID(providerUID)},
		},
	}
}

func topologyReadOnlyGrants(principal, providerUID string) *staticGrants {
	return &staticGrants{
		principal: principal,
		grants: []cloudmodel.Grant{
			{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderRead, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionTopologyRead, Scope: providerScopeID(providerUID)},
		},
	}
}

func topologyWriteOnlyGrants(principal, providerUID string) *staticGrants {
	return &staticGrants{
		principal: principal,
		grants: []cloudmodel.Grant{
			{Action: cloudmodel.ActionCloudPlatformWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionTopologyWrite, Scope: providerScopeID(providerUID)},
		},
	}
}

func topologyCrossProviderGrants(principal, writeUID, readExtraUID string) *staticGrants {
	return &staticGrants{
		principal: principal,
		grants: []cloudmodel.Grant{
			{Action: cloudmodel.ActionCloudPlatformWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderWrite, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderRead, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionTopologyWrite, Scope: providerScopeID(writeUID)},
			{Action: cloudmodel.ActionTopologyRead, Scope: providerScopeID(writeUID)},
			{Action: cloudmodel.ActionTopologyRead, Scope: providerScopeID(readExtraUID)},
		},
	}
}

func (rt *f15Runtime) hostingLocationCollection() http.Handler {
	return NewHostingLocationCollectionHandler(rt.store, rt.grants, rt.pub, rt.idemp)
}

func (rt *f15Runtime) hostingLocationItem() http.Handler {
	return NewHostingLocationItemHandler(rt.store, rt.grants, rt.pub)
}

func (rt *f15Runtime) datacenterCollection() http.Handler {
	return NewDatacenterCollectionHandler(rt.store, rt.grants, rt.pub, rt.idemp)
}

func (rt *f15Runtime) datacenterItem() http.Handler {
	return NewDatacenterItemHandler(rt.store, rt.grants, rt.pub)
}

func (rt *f15Runtime) faultDomainCollection() http.Handler {
	return NewFaultDomainCollectionHandler(rt.store, rt.grants, rt.pub, rt.idemp)
}

func (rt *f15Runtime) faultDomainItem() http.Handler {
	return NewFaultDomainItemHandler(rt.store, rt.grants, rt.pub)
}

func (rt *f15Runtime) infrastructureStackCollection() http.Handler {
	return NewInfrastructureStackCollectionHandler(rt.store, rt.grants, rt.pub, rt.idemp)
}

func (rt *f15Runtime) infrastructureStackItem() http.Handler {
	return NewInfrastructureStackItemHandler(rt.store, rt.grants, rt.pub)
}

func muxTopology(rt *f15Runtime) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations", rt.hostingLocationCollection())
	mux.Handle("POST /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations", rt.hostingLocationCollection())
	mux.Handle("GET /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/{uid}", rt.hostingLocationItem())
	mux.Handle("PATCH /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/{uid}", rt.hostingLocationItem())
	mux.Handle("PUT /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/{uid}", rt.hostingLocationItem())
	mux.Handle("DELETE /apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations/{uid}", rt.hostingLocationItem())

	mux.Handle("GET /apis/infrastructure.sovrunn.io/v1alpha1/datacenters", rt.datacenterCollection())
	mux.Handle("POST /apis/infrastructure.sovrunn.io/v1alpha1/datacenters", rt.datacenterCollection())
	mux.Handle("GET /apis/infrastructure.sovrunn.io/v1alpha1/datacenters/{uid}", rt.datacenterItem())
	mux.Handle("PATCH /apis/infrastructure.sovrunn.io/v1alpha1/datacenters/{uid}", rt.datacenterItem())
	mux.Handle("PUT /apis/infrastructure.sovrunn.io/v1alpha1/datacenters/{uid}", rt.datacenterItem())
	mux.Handle("DELETE /apis/infrastructure.sovrunn.io/v1alpha1/datacenters/{uid}", rt.datacenterItem())

	mux.Handle("GET /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains", rt.faultDomainCollection())
	mux.Handle("POST /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains", rt.faultDomainCollection())
	mux.Handle("GET /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/{uid}", rt.faultDomainItem())
	mux.Handle("PATCH /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/{uid}", rt.faultDomainItem())
	mux.Handle("PUT /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/{uid}", rt.faultDomainItem())
	mux.Handle("DELETE /apis/infrastructure.sovrunn.io/v1alpha1/fault-domains/{uid}", rt.faultDomainItem())

	mux.Handle("GET /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks", rt.infrastructureStackCollection())
	mux.Handle("POST /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks", rt.infrastructureStackCollection())
	mux.Handle("GET /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/{uid}", rt.infrastructureStackItem())
	mux.Handle("PATCH /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/{uid}", rt.infrastructureStackItem())
	mux.Handle("PUT /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/{uid}", rt.infrastructureStackItem())
	mux.Handle("DELETE /apis/infrastructure.sovrunn.io/v1alpha1/infrastructure-stacks/{uid}", rt.infrastructureStackItem())
	return mux
}

func seedPlatformAndProvider(t *testing.T, rt *f15Runtime, providerUID string) {
	t.Helper()
	if !rt.store.HasCloudPlatform() {
		if _, prob := rt.store.CreateCloudPlatform(model.CloudPlatform{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform},
			Metadata: apimeta.ObjectMeta{Name: "root-plat", UID: testPlatformUID},
			Spec: model.CloudPlatformSpec{OwnerRegistration: model.OwnerRegistration{
				LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
			}},
			Status: model.CloudPlatformStatus{Phase: model.CloudPlatformPhaseActive},
		}); prob != nil {
			t.Fatalf("seed platform: %#v", prob)
		}
	}
	if _, ok := rt.store.GetCloudProvider(providerUID); ok {
		return
	}
	if _, prob := rt.store.CreateCloudProvider(model.CloudProvider{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider},
		Metadata: apimeta.ObjectMeta{Name: "prov-" + providerUID[:4], UID: providerUID, ResourceVersion: "1"},
		Spec:     model.CloudProviderSpec{OperatingMarkets: []string{"US"}},
		Status:   model.CloudProviderStatus{Phase: model.CloudProviderPhaseActive},
	}); prob != nil {
		t.Fatalf("seed provider: %#v", prob)
	}
}

func seedOtherProvider(t *testing.T, rt *f15Runtime, providerUID string) {
	t.Helper()
	if _, prob := rt.store.CreateCloudProvider(model.CloudProvider{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudProvider, Kind: model.KindCloudProvider},
		Metadata: apimeta.ObjectMeta{Name: "prov-b", UID: providerUID, ResourceVersion: "1"},
		Spec:     model.CloudProviderSpec{OperatingMarkets: []string{"IN"}},
		Status:   model.CloudProviderStatus{Phase: model.CloudProviderPhaseActive},
	}); prob != nil {
		t.Fatalf("seed other provider: %#v", prob)
	}
}

func providerScopeRef(uid, name string) *apimeta.ScopeRef {
	return &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudProvider,
		Kind:       string(apimeta.ScopeCloudProvider),
		Name:       name,
		UID:        uid,
	}}
}

func hostingLocationCreateBody(name, country, locality, admin, desc string) []byte {
	adminJSON := ""
	if admin != "" {
		adminJSON = `,"administrativeAreaCode":"` + admin + `"`
	}
	descJSON := ""
	if desc != "" {
		descJSON = `,"description":"` + desc + `"`
	}
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{"countryCode":"` + country + `","locality":"` + locality + `"` + adminJSON + descJSON + `}
	}`)
}

func TestHostingLocationCollection_ListCreateISORootGate(t *testing.T) {
	rt := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	mux := muxTopology(rt)

	rec := doMuxUnauth(t, mux, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations", nil, nil)
	if rec.Code != http.StatusUnauthorized || rt.audit.Len() != 0 {
		t.Fatalf("auth required: %d audit=%d", rec.Code, rt.audit.Len())
	}

	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "hl-1"}
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations",
		hostingLocationCreateBody("loc-a", "US", "Austin", "US-TX", "desc"), headers)
	if rec.Code != http.StatusConflict || !hasViolation(decodeProblem(t, rec), violationCloudPlatformRootRequired) {
		t.Fatalf("root gate: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != 1 {
		t.Fatalf("root gate audit=%d", rt.audit.Len())
	}

	rtFail := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	muxFail := muxTopology(rtFail)
	rtFail.audit.SetFail(errors.New("audit down"))
	rec = doMux(t, muxFail, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations",
		hostingLocationCreateBody("loc-fail", "US", "Austin", "", ""), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "hl-fail",
		})
	if rec.Code != http.StatusInternalServerError || rtFail.audit.Len() != 0 {
		t.Fatalf("root audit fail: %d audit=%d", rec.Code, rtFail.audit.Len())
	}

	seedPlatformAndProvider(t, rt, testProviderUID)
	before := rt.audit.Len()
	headers["Idempotency-Key"] = "hl-2"
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations",
		hostingLocationCreateBody("loc-a", "US", "Austin", "US-TX", "west"), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created model.HostingLocation
	_ = json.Unmarshal(rec.Body.Bytes(), &created)
	if created.Status.Phase != model.HostingLocationPhaseActive {
		t.Fatalf("initial status=%#v", created.Status)
	}
	if created.Metadata.ScopeRef == nil || created.Metadata.ScopeRef.UID != testProviderUID {
		t.Fatalf("scope from topology.write grant: %#v", created.Metadata.ScopeRef)
	}
	if rt.audit.Len() != before+1 {
		t.Fatalf("create audit=%d", rt.audit.Len())
	}

	before = rt.audit.Len()
	headers["Idempotency-Key"] = "hl-zz"
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations",
		hostingLocationCreateBody("loc-zz", "ZZ", "Nowhere", "", ""), headers)
	if rec.Code != http.StatusUnprocessableEntity || rt.audit.Len() != before {
		t.Fatalf("ZZ: %d audit=%d", rec.Code, rt.audit.Len())
	}

	before = rt.audit.Len()
	headers["Idempotency-Key"] = "hl-dup"
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations",
		hostingLocationCreateBody("loc-a", "US", "Dallas", "", ""), headers)
	if decodeProblem(t, rec).Code != apiproblem.CodeAlreadyExists || rt.audit.Len() != before {
		t.Fatalf("duplicate: %d audit=%d", rec.Code, rt.audit.Len())
	}

	headers["Idempotency-Key"] = "hl-2"
	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations",
		hostingLocationCreateBody("loc-a", "US", "Austin", "US-TX", "west"), headers)
	if rec.Code != http.StatusCreated || rt.audit.Len() != before {
		t.Fatalf("replay: %d audit=%d", rec.Code, rt.audit.Len())
	}
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations",
		hostingLocationCreateBody("loc-b", "US", "Houston", "", ""), headers)
	if rec.Code != http.StatusConflict || rt.audit.Len() != before {
		t.Fatalf("mismatch: %d audit=%d", rec.Code, rt.audit.Len())
	}

	before = rt.audit.Len()
	headers["Idempotency-Key"] = "hl-status"
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations", []byte(`{
		"metadata":{"name":"loc-status"},
		"spec":{"countryCode":"US","locality":"Austin"},
		"status":{"phase":"Active"}
	}`), headers)
	if rec.Code != http.StatusForbidden || rt.audit.Len() != before+1 {
		t.Fatalf("status: %d audit=%d", rec.Code, rt.audit.Len())
	}
	before = rt.audit.Len()
	headers["Idempotency-Key"] = "hl-scope"
	rec = doMux(t, mux, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations", []byte(`{
		"metadata":{"name":"loc-scope","scopeRef":{"apiVersion":"v","kind":"CloudProvider","name":"p","uid":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}},
		"spec":{"countryCode":"US","locality":"Austin"}
	}`), headers)
	if rec.Code != http.StatusUnprocessableEntity || rt.audit.Len() != before {
		t.Fatalf("scopeRef: %d audit=%d", rec.Code, rt.audit.Len())
	}

	rtNoRead := newF15Runtime(topologyWriteOnlyGrants(testPrincipal, testProviderUID))
	seedPlatformAndProvider(t, rtNoRead, testProviderUID)
	muxNoRead := muxTopology(rtNoRead)
	rec = doMux(t, muxNoRead, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations", nil, nil)
	if rec.Code != http.StatusForbidden || rtNoRead.audit.Len() != 1 {
		t.Fatalf("list no-read: %d audit=%d", rec.Code, rtNoRead.audit.Len())
	}

	rec = doMux(t, mux, http.MethodGet, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status=%d", rec.Code)
	}
	var list apimeta.ListEnvelope[model.HostingLocation]
	_ = json.Unmarshal(rec.Body.Bytes(), &list)
	if len(list.Items) != 1 || list.Items[0].Metadata.UID != created.Metadata.UID {
		t.Fatalf("list items=%#v", list.Items)
	}

	rt3 := newF15Runtime(topologyGrants(testPrincipal, testProviderUID))
	seedPlatformAndProvider(t, rt3, testProviderUID)
	mux3 := muxTopology(rt3)
	rt3.audit.SetFail(errors.New("down"))
	rec = doMux(t, mux3, http.MethodPost, "/apis/infrastructure.sovrunn.io/v1alpha1/hosting-locations",
		hostingLocationCreateBody("loc-x", "US", "Austin", "", ""), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "hl-x",
		})
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("create audit fail status=%d", rec.Code)
	}
	if len(rt3.store.ListHostingLocations()) != 0 {
		t.Fatal("create must not publish on audit failure")
	}

	stages, ok := validate.StageSetFor(validate.KindHostingLocation, validate.RequestClassCollectionCreate)
	if !ok || stages.Reference == nil || stages.Semantic == nil || stages.Defaulting == nil {
		t.Fatal("HostingLocation StageSet required")
	}
}
