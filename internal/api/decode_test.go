package api

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func TestDecodePhaseOne_BoundedMalformedDuplicateAndExtraction(t *testing.T) {
	t.Parallel()
	lim := apivalid.DefaultLimits()

	oversized := []byte(strings.Repeat("a", lim.MaxObjectBytes+1))
	if _, prob := DecodePhaseOne(oversized, lim); prob == nil || prob.Code != apiproblem.CodeRequestTooLarge {
		t.Fatalf("oversized: %#v", prob)
	}

	if _, prob := DecodePhaseOne([]byte(`{"metadata":`), lim); prob == nil || prob.Code != apiproblem.CodeMalformedRequest {
		t.Fatalf("malformed: %#v", prob)
	}

	dupCases := []string{
		`{"metadata":{"name":"x"},"metadata":{"name":"y"}}`,
		`{"metadata":{"name":"x","name":"y"},"spec":{"operatingMarkets":["US"]}}`,
		`{"metadata":{"name":"x"},"spec":{"operatingMarkets":["US"],"operatingMarkets":["IN"]}}`,
	}
	for _, raw := range dupCases {
		if _, prob := DecodePhaseOne([]byte(raw), lim); prob == nil || prob.Code != apiproblem.CodeDuplicateField {
			t.Fatalf("duplicate %s: %#v", raw, prob)
		}
	}

	body := []byte(`{
		"metadata":{"name":"part"},
		"spec":{
			"cloudPlatformRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudPlatform","name":"plat","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
			"cloudProviderRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudProvider","name":"prov","uid":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
			"environment":"development"
		}
	}`)
	got, prob := DecodePhaseOne(body, lim)
	if prob != nil {
		t.Fatalf("phase one: %#v", prob)
	}
	if got.Name != "part" {
		t.Fatalf("name=%q", got.Name)
	}
	if ref, ok := got.Refs["spec.cloudPlatformRef"]; !ok || ref.UID != "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		t.Fatalf("platform ref: %#v", got.Refs)
	}
	if ref, ok := got.Refs["spec.cloudProviderRef"]; !ok || ref.UID != "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" {
		t.Fatalf("provider ref: %#v", got.Refs)
	}
}

func TestDecodePhaseTwoStrict_ClosedContractBeforeAccept(t *testing.T) {
	t.Parallel()
	var dst model.CloudPlatform
	prob := DecodePhaseTwoStrict(validate.KindCloudPlatform, validate.RequestClassCollectionCreate, []byte(`{
		"metadata":{"name":"plat","scopeRef":{"kind":"Platform","name":"platform","apiVersion":"v"}},
		"spec":{"ownerRegistration":{"legalName":"L","registrationIdentifier":"R","jurisdictionCode":"US"}}
	}`), &dst)
	if prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("scopeRef must fail classification: %#v", prob)
	}

	prob = DecodePhaseTwoStrict(validate.KindCloudPlatform, validate.RequestClassCollectionCreate, []byte(`{
		"metadata":{"name":"plat"},
		"spec":{"ownerRegistration":{"legalName":"L","registrationIdentifier":"R","jurisdictionCode":"US"}}
	}`), &dst)
	if prob != nil {
		t.Fatalf("valid create: %#v", prob)
	}
	if dst.Metadata.Name != "plat" {
		t.Fatalf("decoded name=%q", dst.Metadata.Name)
	}
}
