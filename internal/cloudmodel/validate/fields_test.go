package validate_test

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func TestValidateImmutablePatch_IdentityReferenceOwnerRegistration(t *testing.T) {
	t.Parallel()
	before := model.CloudPlatform{
		Metadata: apimeta.ObjectMeta{Name: "plat"},
		Spec: model.CloudPlatformSpec{
			OwnerRegistration: model.OwnerRegistration{
				LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
			},
			Description: "old",
		},
	}
	afterOK := before
	afterOK.Spec.Description = "new"
	if prob := validate.ValidateImmutablePatch(validate.KindCloudPlatform, before, afterOK); prob != nil {
		t.Fatalf("description patch: %#v", prob)
	}

	afterName := before
	afterName.Metadata.Name = "other"
	if prob := validate.ValidateImmutablePatch(validate.KindCloudPlatform, before, afterName); prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("name immutability: %#v", prob)
	}

	afterOwner := before
	afterOwner.Spec.OwnerRegistration.LegalName = "Other"
	if prob := validate.ValidateImmutablePatch(validate.KindCloudPlatform, before, afterOwner); prob == nil {
		t.Fatal("ownerRegistration must be immutable")
	}
}

func TestValidateImmutablePatch_TopologyNameAndParentRef(t *testing.T) {
	t.Parallel()
	before := model.Datacenter{
		Metadata: apimeta.ObjectMeta{
			Name: "dc1",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				Kind: string(apimeta.ScopeCloudProvider), UID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Name: "p", APIVersion: "v",
			}},
		},
		Spec: model.DatacenterSpec{
			HostingLocationRef: apimeta.TypedRef{
				APIVersion: model.APIVersionHostingLocation,
				Kind:       model.KindHostingLocation,
				Name:       "loc",
				UID:        "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			},
			Description: "d",
		},
	}
	afterDesc := before
	afterDesc.Spec.Description = "new"
	if prob := validate.ValidateImmutablePatch(validate.KindDatacenter, before, afterDesc); prob != nil {
		t.Fatalf("description: %#v", prob)
	}
	afterRef := before
	afterRef.Spec.HostingLocationRef.UID = "cccccccccccccccccccccccccccccccc"
	if prob := validate.ValidateImmutablePatch(validate.KindDatacenter, before, afterRef); prob == nil {
		t.Fatal("parent ref must be immutable")
	}
}
