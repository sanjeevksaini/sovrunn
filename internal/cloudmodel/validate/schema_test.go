package validate_test

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func TestValidateSchema_CloudPlatformRequiredAndISO(t *testing.T) {
	t.Parallel()
	cp := model.CloudPlatform{
		Metadata: apimeta.ObjectMeta{Name: "acme-platform"},
		Spec: model.CloudPlatformSpec{
			OwnerRegistration: model.OwnerRegistration{
				LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
			},
		},
	}
	if prob := validate.ValidateSchema(validate.KindCloudPlatform, cp); prob != nil {
		t.Fatalf("valid: %#v", prob)
	}
	cp.Spec.OwnerRegistration.JurisdictionCode = "ZZ"
	if prob := validate.ValidateSchema(validate.KindCloudPlatform, cp); prob == nil {
		t.Fatal("unassigned jurisdiction must fail")
	}
	cp.Spec.OwnerRegistration.JurisdictionCode = "US"
	cp.Metadata.Name = ""
	if prob := validate.ValidateSchema(validate.KindCloudPlatform, cp); prob == nil {
		t.Fatal("missing name must fail")
	}
}

func TestValidateSchema_CloudProviderOperatingMarkets(t *testing.T) {
	t.Parallel()
	cp := model.CloudProvider{
		Metadata: apimeta.ObjectMeta{Name: "prov"},
		Spec:     model.CloudProviderSpec{OperatingMarkets: []string{"US", "IN"}},
	}
	if prob := validate.ValidateSchema(validate.KindCloudProvider, cp); prob != nil {
		t.Fatalf("valid: %#v", prob)
	}
	cp.Spec.OperatingMarkets = nil
	if prob := validate.ValidateSchema(validate.KindCloudProvider, cp); prob == nil {
		t.Fatal("empty markets must fail")
	}
}

func TestValidateSchema_ParticipationEnvironmentAndRefs(t *testing.T) {
	t.Parallel()
	platformUID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	p := model.CloudProviderParticipation{
		Metadata: apimeta.ObjectMeta{
			Name: "part",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform,
				Kind:       model.KindCloudPlatform,
				Name:       "plat",
				UID:        platformUID,
			}},
		},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform,
				Kind:       model.KindCloudPlatform,
				Name:       "plat",
				UID:        platformUID,
			},
			CloudProviderRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudProvider,
				Kind:       model.KindCloudProvider,
				Name:       "prov",
				UID:        "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			},
			Environment: model.ParticipationEnvironmentDevelopment,
		},
	}
	if prob := validate.ValidateSchema(validate.KindCloudProviderParticipation, p); prob != nil {
		t.Fatalf("valid: %#v", prob)
	}
	p.Spec.Environment = "production"
	if prob := validate.ValidateSchema(validate.KindCloudProviderParticipation, p); prob == nil {
		t.Fatal("non-development environment must fail")
	}
}

func TestValidateSchema_HostingLocationConstraints(t *testing.T) {
	t.Parallel()
	hl := model.HostingLocation{
		Metadata: apimeta.ObjectMeta{
			Name: "loc",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudProvider,
				Kind:       string(apimeta.ScopeCloudProvider),
				Name:       "prov",
				UID:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			}},
		},
		Spec: model.HostingLocationSpec{
			CountryCode:            "US",
			Locality:               "Austin",
			AdministrativeAreaCode: "US-TX",
		},
	}
	if prob := validate.ValidateSchema(validate.KindHostingLocation, hl); prob != nil {
		t.Fatalf("valid: %#v", prob)
	}
	hl.Spec.AdministrativeAreaCode = "IN-MH"
	prob := validate.ValidateSchema(validate.KindHostingLocation, hl)
	if prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("prefix mismatch: %#v", prob)
	}
}

func TestValidateSchema_TopologyRefsRequired(t *testing.T) {
	t.Parallel()
	scope := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudProvider,
		Kind:       string(apimeta.ScopeCloudProvider),
		Name:       "prov",
		UID:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}}
	dc := model.Datacenter{
		Metadata: apimeta.ObjectMeta{Name: "dc", ScopeRef: scope},
		Spec: model.DatacenterSpec{
			HostingLocationRef: apimeta.TypedRef{
				APIVersion: model.APIVersionHostingLocation,
				Kind:       model.KindHostingLocation,
				Name:       "loc",
				UID:        "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			},
		},
	}
	if prob := validate.ValidateSchema(validate.KindDatacenter, dc); prob != nil {
		t.Fatalf("datacenter: %#v", prob)
	}
	dc.Spec.HostingLocationRef.UID = ""
	if prob := validate.ValidateSchema(validate.KindDatacenter, dc); prob == nil {
		t.Fatal("uid-pinned ref required")
	}
}
