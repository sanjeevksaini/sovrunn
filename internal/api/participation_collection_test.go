package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func participationGrants(principal, platformUID, providerUID string) *cloudmodel.BootstrapGrantResolver {
	return cloudmodel.NewBootstrapGrantResolver(cloudmodel.BootstrapGrantConfig{
		PrincipalID:       principal,
		CloudPlatformUIDs: []string{platformUID},
		CloudProviderUIDs: []string{providerUID},
		ParticipationRelations: []cloudmodel.ParticipationRelation{{
			CloudPlatformUID: platformUID,
			CloudProviderUID: providerUID,
		}},
	})
}

func participationNoReadGrants(principal, platformUID, providerUID string) *staticGrants {
	return &staticGrants{
		principal: principal,
		grants: []cloudmodel.Grant{
			{Action: cloudmodel.ActionCloudPlatformRead, Scope: cloudmodel.PlatformRootScope},
			{Action: cloudmodel.ActionCloudProviderRead, Scope: cloudmodel.PlatformRootScope},
			{
				Action:                  cloudmodel.ActionParticipationRequestProvider,
				Scope:                   apimeta.ScopeIdentity{Kind: apimeta.ScopeCloudProvider, UID: providerUID},
				RelatedCloudPlatformUID: platformUID,
			},
		},
	}
}

func (rt *f15Runtime) participationCollection() http.Handler {
	return NewParticipationCollectionHandler(rt.store, rt.grants, rt.pub, rt.idemp)
}

func (rt *f15Runtime) participationItem() http.Handler {
	return NewParticipationItemHandler(rt.store, rt.grants, rt.pub)
}

func muxParticipation(rt *f15Runtime) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations", rt.participationCollection())
	mux.Handle("POST /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations", rt.participationCollection())
	mux.Handle("GET /apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/{uid}", rt.participationItem())
	return mux
}

func seedPlatformAndProviderForParticipation(t *testing.T, rt *f15Runtime, platformUID, providerUID string) {
	t.Helper()
	if !rt.store.HasCloudPlatform() {
		if _, prob := rt.store.CreateCloudPlatform(model.CloudPlatform{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform},
			Metadata: apimeta.ObjectMeta{Name: "root-plat", UID: platformUID},
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

func participationCreateBody(name, platformUID, providerUID string) []byte {
	return []byte(`{
		"metadata":{"name":"` + name + `"},
		"spec":{
			"cloudPlatformRef":{
				"apiVersion":"core.sovrunn.io/v1alpha1",
				"kind":"CloudPlatform",
				"name":"root-plat",
				"uid":"` + platformUID + `"
			},
			"cloudProviderRef":{
				"apiVersion":"core.sovrunn.io/v1alpha1",
				"kind":"CloudProvider",
				"name":"prov",
				"uid":"` + providerUID + `"
			},
			"environment":"development"
		}
	}`)
}

func TestParticipationCollection_ListCreateClosedContract(t *testing.T) {
	rt := newF15Runtime(participationGrants(testPrincipal, testPlatformUID, testProviderUID))
	mux := muxParticipation(rt)

	rec := doMuxUnauth(t, mux, http.MethodGet, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations", nil, nil)
	if rec.Code != http.StatusUnauthorized || rt.audit.Len() != 0 {
		t.Fatalf("auth required: %d audit=%d", rec.Code, rt.audit.Len())
	}

	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "part-root"}
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody("part-a", testPlatformUID, testProviderUID), headers)
	if rec.Code != http.StatusConflict || !hasViolation(decodeProblem(t, rec), violationCloudPlatformRootRequired) {
		t.Fatalf("root gate: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != 1 {
		t.Fatalf("root gate audit=%d", rt.audit.Len())
	}

	rtFail := newF15Runtime(participationGrants(testPrincipal, testPlatformUID, testProviderUID))
	muxFail := muxParticipation(rtFail)
	rtFail.audit.SetFail(errors.New("audit down"))
	rec = doMux(t, muxFail, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody("part-fail", testPlatformUID, testProviderUID), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "part-fail",
		})
	if rec.Code != http.StatusInternalServerError || rtFail.audit.Len() != 0 {
		t.Fatalf("root audit fail: %d audit=%d", rec.Code, rtFail.audit.Len())
	}

	seedPlatformAndProviderForParticipation(t, rt, testPlatformUID, testProviderUID)
	before := rt.audit.Len()
	headers["Idempotency-Key"] = "part-1"
	createdAtFloor := time.Now().UTC().Add(-time.Second)
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody("part-a", testPlatformUID, testProviderUID), headers)
	createdAtCeil := time.Now().UTC().Add(time.Second)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created model.CloudProviderParticipation
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("create decode: %v", err)
	}
	if created.Status.Phase != model.ParticipationPhasePending {
		t.Fatalf("phase=%q", created.Status.Phase)
	}
	if created.Status.PlatformSuspended || created.Status.ProviderSuspended {
		t.Fatalf("holds must be false: %#v", created.Status)
	}
	if created.Metadata.ScopeRef == nil || created.Metadata.ScopeRef.UID != testPlatformUID {
		t.Fatalf("server-derived CloudPlatform scope: %#v", created.Metadata.ScopeRef)
	}
	if created.Spec.CloudPlatformRef.UID != created.Metadata.ScopeRef.UID {
		t.Fatalf("cloudPlatformRef.uid/scopeRef.uid mismatch: %q vs %q",
			created.Spec.CloudPlatformRef.UID, created.Metadata.ScopeRef.UID)
	}
	expires, err := time.Parse(time.RFC3339, created.Status.RequestExpiresAt)
	if err != nil {
		t.Fatalf("requestExpiresAt parse: %v", err)
	}
	createdAt, err := time.Parse(time.RFC3339, created.Metadata.CreatedAt)
	if err != nil {
		t.Fatalf("createdAt parse: %v", err)
	}
	if createdAt.Before(createdAtFloor) || createdAt.After(createdAtCeil) {
		t.Fatalf("createdAt=%v outside [%v,%v]", createdAt, createdAtFloor, createdAtCeil)
	}
	wantExpiry := createdAt.Add(7 * 24 * time.Hour)
	if !expires.Equal(wantExpiry) {
		t.Fatalf("requestExpiresAt=%v want %v", expires, wantExpiry)
	}
	if rt.audit.Len() != before+1 {
		t.Fatalf("create audit=%d", rt.audit.Len())
	}

	// Inaccessible reference → audited safe 404.
	before = rt.audit.Len()
	headers["Idempotency-Key"] = "part-inacc"
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody("part-x", testPlatformUID, "dddddddddddddddddddddddddddddddd"), headers)
	if rec.Code != http.StatusNotFound || !hasViolation(decodeProblem(t, rec), violationAuthorizationSafeDenial) {
		t.Fatalf("inaccessible: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before+1 {
		t.Fatalf("inaccessible audit=%d", rt.audit.Len())
	}
	rt.audit.SetFail(errors.New("down"))
	beforeFail := rt.audit.Len()
	headers["Idempotency-Key"] = "part-inacc-fail"
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody("part-y", testPlatformUID, "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"), headers)
	if rec.Code != http.StatusInternalServerError || rt.audit.Len() != beforeFail {
		t.Fatalf("inaccessible audit fail: %d audit=%d", rec.Code, rt.audit.Len())
	}
	rt.audit.SetFail(nil)

	// If-Match on create → MALFORMED_REQUEST, no audit.
	before = rt.audit.Len()
	headers["Idempotency-Key"] = "part-ifmatch"
	headers["If-Match"] = "1"
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody("part-if", testPlatformUID, testProviderUID), headers)
	if rec.Code != http.StatusBadRequest || decodeProblem(t, rec).Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("If-Match on create: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before {
		t.Fatal("If-Match must not audit")
	}
	delete(headers, "If-Match")

	// Duplicate name → ALREADY_EXISTS, no audit.
	before = rt.audit.Len()
	headers["Idempotency-Key"] = "part-dup-name"
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody("part-a", testPlatformUID, testProviderUID), headers)
	if decodeProblem(t, rec).Code != apiproblem.CodeAlreadyExists || rt.audit.Len() != before {
		t.Fatalf("duplicate name: %d %#v audit=%d", rec.Code, decodeProblem(t, rec), rt.audit.Len())
	}

	// Duplicate non-terminal pair → ALREADY_EXISTS / VS0_PARTICIPATION_DUPLICATE, no audit.
	before = rt.audit.Len()
	headers["Idempotency-Key"] = "part-dup-pair"
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody("part-b", testPlatformUID, testProviderUID), headers)
	if decodeProblem(t, rec).Code != apiproblem.CodeAlreadyExists ||
		!hasViolation(decodeProblem(t, rec), validate.ViolationParticipationDuplicate) ||
		rt.audit.Len() != before {
		t.Fatalf("duplicate pair: %d %#v audit=%d", rec.Code, decodeProblem(t, rec), rt.audit.Len())
	}

	// Idempotency replay / mismatch — no extra audit.
	headers["Idempotency-Key"] = "part-1"
	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody("part-a", testPlatformUID, testProviderUID), headers)
	if rec.Code != http.StatusCreated || rt.audit.Len() != before {
		t.Fatalf("replay: %d audit=%d", rec.Code, rt.audit.Len())
	}
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody("part-c", testPlatformUID, testProviderUID), headers)
	if rec.Code != http.StatusConflict ||
		!hasViolation(decodeProblem(t, rec), cloudmodel.ViolationIdempotencyKeyReuseMismatch) ||
		rt.audit.Len() != before {
		t.Fatalf("idempotency mismatch: %d %#v audit=%d", rec.Code, decodeProblem(t, rec), rt.audit.Len())
	}

	// FEATURE-0021 deferred fields → UNKNOWN_FIELD, no audit.
	before = rt.audit.Len()
	headers["Idempotency-Key"] = "part-f21"
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations", []byte(`{
		"metadata":{"name":"part-f21"},
		"spec":{
			"cloudPlatformRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudPlatform","name":"root-plat","uid":"`+testPlatformUID+`"},
			"cloudProviderRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudProvider","name":"prov","uid":"`+testProviderUID+`"},
			"environment":"development",
			"providerSelectionModes":["CustomerSelected"]
		}
	}`), headers)
	if rec.Code != http.StatusBadRequest || decodeProblem(t, rec).Code != apiproblem.CodeUnknownField || rt.audit.Len() != before {
		t.Fatalf("FEATURE-0021 deferred: %d %#v audit=%d", rec.Code, decodeProblem(t, rec), rt.audit.Len())
	}
	headers["Idempotency-Key"] = "part-f21b"
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations", []byte(`{
		"metadata":{"name":"part-f21b"},
		"spec":{
			"cloudPlatformRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudPlatform","name":"root-plat","uid":"`+testPlatformUID+`"},
			"cloudProviderRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudProvider","name":"prov","uid":"`+testProviderUID+`"},
			"environment":"development",
			"permittedHostingLocationRefs":[]
		}
	}`), headers)
	if rec.Code != http.StatusBadRequest || decodeProblem(t, rec).Code != apiproblem.CodeUnknownField || rt.audit.Len() != before {
		t.Fatalf("permittedHostingLocationRefs: %d %#v audit=%d", rec.Code, decodeProblem(t, rec), rt.audit.Len())
	}

	// status audited; scopeRef/unknown/deferred unaudited.
	before = rt.audit.Len()
	headers["Idempotency-Key"] = "part-status"
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations", []byte(`{
		"metadata":{"name":"part-status"},
		"spec":{
			"cloudPlatformRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudPlatform","name":"root-plat","uid":"`+testPlatformUID+`"},
			"cloudProviderRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudProvider","name":"prov","uid":"`+testProviderUID+`"},
			"environment":"development"
		},
		"status":{"phase":"Active"}
	}`), headers)
	if rec.Code != http.StatusForbidden || rt.audit.Len() != before+1 {
		t.Fatalf("status: %d audit=%d", rec.Code, rt.audit.Len())
	}
	before = rt.audit.Len()
	headers["Idempotency-Key"] = "part-sys"
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations", []byte(`{
		"metadata":{"name":"part-sys","uid":"ffffffffffffffffffffffffffffffff"},
		"spec":{
			"cloudPlatformRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudPlatform","name":"root-plat","uid":"`+testPlatformUID+`"},
			"cloudProviderRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudProvider","name":"prov","uid":"`+testProviderUID+`"},
			"environment":"development"
		}
	}`), headers)
	if rec.Code != http.StatusForbidden || rt.audit.Len() != before+1 {
		t.Fatalf("system-owned metadata: %d audit=%d", rec.Code, rt.audit.Len())
	}
	before = rt.audit.Len()
	headers["Idempotency-Key"] = "part-scope"
	rec = doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations", []byte(`{
		"metadata":{"name":"part-scope","scopeRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudPlatform","name":"root-plat","uid":"`+testPlatformUID+`"}},
		"spec":{
			"cloudPlatformRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudPlatform","name":"root-plat","uid":"`+testPlatformUID+`"},
			"cloudProviderRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudProvider","name":"prov","uid":"`+testProviderUID+`"},
			"environment":"development"
		}
	}`), headers)
	if rec.Code != http.StatusUnprocessableEntity || rt.audit.Len() != before {
		t.Fatalf("scopeRef: %d audit=%d", rec.Code, rt.audit.Len())
	}

	// LIST without read grant.
	rtNoRead := newF15Runtime(participationNoReadGrants(testPrincipal, testPlatformUID, testProviderUID))
	seedPlatformAndProviderForParticipation(t, rtNoRead, testPlatformUID, testProviderUID)
	muxNoRead := muxParticipation(rtNoRead)
	rec = doMux(t, muxNoRead, http.MethodGet, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations", nil, nil)
	if rec.Code != http.StatusForbidden || rtNoRead.audit.Len() != 1 {
		t.Fatalf("list no-read: %d audit=%d", rec.Code, rtNoRead.audit.Len())
	}

	// Authorized LIST ascending metadata.uid.
	rec = doMux(t, mux, http.MethodGet, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status=%d", rec.Code)
	}
	var list apimeta.ListEnvelope[model.CloudProviderParticipation]
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatalf("list decode: %v", err)
	}
	if len(list.Items) != 1 || list.Items[0].Metadata.UID != created.Metadata.UID {
		t.Fatalf("list items=%#v", list.Items)
	}
	for i := 1; i < len(list.Items); i++ {
		if list.Items[i-1].Metadata.UID >= list.Items[i].Metadata.UID {
			t.Fatalf("list not ascending uid: %#v", list.Items)
		}
	}

	// Successful create audit-append failure → INTERNAL_ERROR, no publication.
	rt3 := newF15Runtime(participationGrants(testPrincipal, testPlatformUID, testProviderUID))
	seedPlatformAndProviderForParticipation(t, rt3, testPlatformUID, testProviderUID)
	mux3 := muxParticipation(rt3)
	rt3.audit.SetFail(errors.New("down"))
	rec = doMux(t, mux3, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody("part-x", testPlatformUID, testProviderUID), map[string]string{
			"Content-Type": "application/json", "Idempotency-Key": "part-x",
		})
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("create audit fail status=%d", rec.Code)
	}
	if len(rt3.store.ListParticipations()) != 0 {
		t.Fatal("create must not publish on audit failure")
	}

	stages, ok := validate.StageSetFor(validate.KindCloudProviderParticipation, validate.RequestClassCollectionCreate)
	if !ok || stages.Reference == nil || stages.Semantic == nil || stages.Defaulting == nil {
		t.Fatal("participation collection-create StageSet required")
	}
	if _, ok := validate.StageSetFor(validate.KindCloudProviderParticipation, validate.RequestClassPatch); ok {
		t.Fatal("participation must not expose a PATCH StageSet")
	}
}
