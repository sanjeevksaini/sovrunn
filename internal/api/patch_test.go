package api

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func TestApplyMergePatch_RFC7396StagedCloneAndNullRules(t *testing.T) {
	t.Parallel()
	current := model.CloudPlatform{
		TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionCloudPlatform, Kind: model.KindCloudPlatform},
		Metadata: apimeta.ObjectMeta{Name: "plat", UID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", ResourceVersion: "1"},
		Spec: model.CloudPlatformSpec{
			OwnerRegistration: model.OwnerRegistration{
				LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
			},
			Description: "old",
		},
		Status: model.CloudPlatformStatus{Phase: model.CloudPlatformPhaseActive},
	}

	merged, prob := ApplyMergePatch(validate.KindCloudPlatform, current, []byte(`{"spec":{"description":"new"}}`))
	if prob != nil {
		t.Fatalf("merge: %#v", prob)
	}
	got := merged.(model.CloudPlatform)
	if got.Spec.Description != "new" {
		t.Fatalf("description=%q", got.Spec.Description)
	}
	if current.Spec.Description != "old" {
		t.Fatal("current must remain an unmutated staged-clone source")
	}

	merged, prob = ApplyMergePatch(validate.KindCloudPlatform, current, []byte(`{"spec":{"description":null}}`))
	if prob != nil {
		t.Fatalf("null description: %#v", prob)
	}
	got = merged.(model.CloudPlatform)
	if got.Spec.Description != "" {
		t.Fatalf("null must remove optional description, got %q", got.Spec.Description)
	}
}

func TestApplyMergePatch_OperatingMarketsNullRejected(t *testing.T) {
	t.Parallel()
	current := model.CloudProvider{
		Metadata: apimeta.ObjectMeta{Name: "prov", UID: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", ResourceVersion: "3"},
		Spec:     model.CloudProviderSpec{DisplayName: "P", OperatingMarkets: []string{"US"}},
	}
	_, prob := ApplyMergePatch(validate.KindCloudProvider, current, []byte(`{"spec":{"operatingMarkets":null}}`))
	if prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("operatingMarkets null: %#v", prob)
	}
}

func TestApplyMergePatch_ImmutableAndUnknownRejectedBeforePublication(t *testing.T) {
	t.Parallel()
	current := model.CloudPlatform{
		Metadata: apimeta.ObjectMeta{Name: "plat", UID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", ResourceVersion: "1"},
		Spec: model.CloudPlatformSpec{
			OwnerRegistration: model.OwnerRegistration{
				LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
			},
		},
	}
	_, prob := ApplyMergePatch(validate.KindCloudPlatform, current, []byte(`{"metadata":{"name":"other"}}`))
	if prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("immutable name: %#v", prob)
	}
	_, prob = ApplyMergePatch(validate.KindCloudPlatform, current, []byte(`{"spec":{"ownerRegistration":{"legalName":"X","registrationIdentifier":"R1","jurisdictionCode":"US"}}}`))
	if prob == nil {
		t.Fatal("ownerRegistration patch must fail")
	}
	_, prob = ApplyMergePatch(validate.KindCloudPlatform, current, []byte(`{"spec":{"unknown":"x"}}`))
	if prob == nil || prob.Code != apiproblem.CodeUnknownField {
		t.Fatalf("unknown: %#v", prob)
	}
}

func TestApplyMergePatch_CloudProviderDisplayNameNull(t *testing.T) {
	t.Parallel()
	current := model.CloudProvider{
		Metadata: apimeta.ObjectMeta{Name: "prov", UID: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb", ResourceVersion: "1"},
		Spec:     model.CloudProviderSpec{DisplayName: "Nice", OperatingMarkets: []string{"US"}},
	}
	merged, prob := ApplyMergePatch(validate.KindCloudProvider, current, []byte(`{"spec":{"displayName":null}}`))
	if prob != nil {
		t.Fatalf("null displayName: %#v", prob)
	}
	got := merged.(model.CloudProvider)
	if got.Spec.DisplayName != "" {
		t.Fatalf("displayName=%q", got.Spec.DisplayName)
	}
	if len(got.Spec.OperatingMarkets) != 1 || got.Spec.OperatingMarkets[0] != "US" {
		t.Fatalf("markets=%v", got.Spec.OperatingMarkets)
	}
}
