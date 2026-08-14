package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

func TestParticipationItem_GetOnlyNoPatchRegistration(t *testing.T) {
	rt := newF15Runtime(participationGrants(testPrincipal, testPlatformUID, testProviderUID))
	mux := muxParticipation(rt)
	seedPlatformAndProviderForParticipation(t, rt, testPlatformUID, testProviderUID)

	headers := map[string]string{"Content-Type": "application/json", "Idempotency-Key": "part-item-1"}
	rec := doMux(t, mux, http.MethodPost, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations",
		participationCreateBody("part-item", testPlatformUID, testProviderUID), headers)
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed create status=%d body=%s", rec.Code, rec.Body.String())
	}
	var created model.CloudProviderParticipation
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("seed decode: %v", err)
	}

	rec = doMux(t, mux, http.MethodGet, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/"+created.Metadata.UID, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("get status=%d body=%s", rec.Code, rec.Body.String())
	}
	var got model.CloudProviderParticipation
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("get decode: %v", err)
	}
	if got.Metadata.UID != created.Metadata.UID || got.Status.Phase != model.ParticipationPhasePending {
		t.Fatalf("got=%#v", got)
	}

	before := rt.audit.Len()
	// PATCH is not registered on the mux; ServeMux returns 405 for method mismatch.
	rec = doMux(t, mux, http.MethodPatch, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/"+created.Metadata.UID,
		[]byte(`{"spec":{"environment":"development"}}`), map[string]string{
			"Content-Type": "application/merge-patch+json",
			"If-Match":     created.Metadata.ResourceVersion,
		})
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PATCH must not be registered: status=%d body=%s", rec.Code, rec.Body.String())
	}
	if rt.audit.Len() != before {
		t.Fatal("PATCH must not audit")
	}

	rec = doMux(t, mux, http.MethodPut, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/"+created.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("PUT %d", rec.Code)
	}
	rec = doMux(t, mux, http.MethodDelete, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/"+created.Metadata.UID, nil, nil)
	if rec.Code != http.StatusMethodNotAllowed || rt.audit.Len() != before {
		t.Fatalf("DELETE %d audit=%d", rec.Code, rt.audit.Len())
	}

	// Inaccessible GET → audited safe 404.
	before = rt.audit.Len()
	rec = doMux(t, mux, http.MethodGet, "/apis/governance.sovrunn.io/v1alpha1/cloud-provider-participations/ffffffffffffffffffffffffffffffff", nil, nil)
	if rec.Code != http.StatusNotFound || !hasViolation(decodeProblem(t, rec), violationAuthorizationSafeDenial) {
		t.Fatalf("inaccessible get: %d %#v", rec.Code, decodeProblem(t, rec))
	}
	if rt.audit.Len() != before+1 {
		t.Fatalf("inaccessible get audit=%d", rt.audit.Len())
	}
}
