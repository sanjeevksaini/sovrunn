package validation

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateOrganizationUnit_ValidNames(t *testing.T) {
	valid := []string{"a", "a1", "a-b", strings.Repeat("a", 63)}
	for _, name := range valid {
		errs := ValidateOrganizationUnit(resources.OrganizationUnit{
			Metadata: resources.Metadata{Name: name},
			Spec:     resources.OrganizationUnitSpec{OrganizationName: "nic"},
		})
		if len(errs) != 0 {
			t.Errorf("name %q: got errors %v, want none", name, errs)
		}
	}
}

func TestValidateOrganizationUnit_EmptyName(t *testing.T) {
	errs := ValidateOrganizationUnit(resources.OrganizationUnit{
		Spec: resources.OrganizationUnitSpec{OrganizationName: "nic"},
	})
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOrganizationUnit_UppercaseName(t *testing.T) {
	errs := ValidateOrganizationUnit(resources.OrganizationUnit{
		Metadata: resources.Metadata{Name: "ABC"},
		Spec:     resources.OrganizationUnitSpec{OrganizationName: "nic"},
	})
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOrganizationUnit_NameWithSpaces(t *testing.T) {
	errs := ValidateOrganizationUnit(resources.OrganizationUnit{
		Metadata: resources.Metadata{Name: "a b"},
		Spec:     resources.OrganizationUnitSpec{OrganizationName: "nic"},
	})
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOrganizationUnit_NameTooLong(t *testing.T) {
	errs := ValidateOrganizationUnit(resources.OrganizationUnit{
		Metadata: resources.Metadata{Name: strings.Repeat("a", 64)},
		Spec:     resources.OrganizationUnitSpec{OrganizationName: "nic"},
	})
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOrganizationUnit_NameLeadingHyphen(t *testing.T) {
	errs := ValidateOrganizationUnit(resources.OrganizationUnit{
		Metadata: resources.Metadata{Name: "-abc"},
		Spec:     resources.OrganizationUnitSpec{OrganizationName: "nic"},
	})
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOrganizationUnit_NameTrailingHyphen(t *testing.T) {
	errs := ValidateOrganizationUnit(resources.OrganizationUnit{
		Metadata: resources.Metadata{Name: "abc-"},
		Spec:     resources.OrganizationUnitSpec{OrganizationName: "nic"},
	})
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOrganizationUnit_EmptyOrganizationName(t *testing.T) {
	errs := ValidateOrganizationUnit(resources.OrganizationUnit{
		Metadata: resources.Metadata{Name: "ministry-health"},
	})
	if !hasFieldError(errs, "spec.organizationName") {
		t.Fatalf("got %v, want spec.organizationName error", errs)
	}
}

func TestValidateOrganizationUnit_InvalidOrganizationNameFormat(t *testing.T) {
	errs := ValidateOrganizationUnit(resources.OrganizationUnit{
		Metadata: resources.Metadata{Name: "ministry-health"},
		Spec:     resources.OrganizationUnitSpec{OrganizationName: "NIC_Org"},
	})
	if !hasFieldError(errs, "spec.organizationName") {
		t.Fatalf("got %v, want spec.organizationName error", errs)
	}
}

func TestValidateOrganizationUnit_ValidOrganizationName(t *testing.T) {
	errs := ValidateOrganizationUnit(resources.OrganizationUnit{
		Metadata: resources.Metadata{Name: "ministry-health"},
		Spec:     resources.OrganizationUnitSpec{OrganizationName: "nic"},
	})
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func TestValidateOUPathSegments_InvalidNameMapsToMetadataName(t *testing.T) {
	errs := ValidateOUPathSegments("nic", "Invalid Name")
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateOUPathSegments_InvalidOrgNameMapsToSpecOrganizationName(t *testing.T) {
	errs := ValidateOUPathSegments("Invalid Org", "ministry-health")
	if !hasFieldError(errs, "spec.organizationName") {
		t.Fatalf("got %v, want spec.organizationName error", errs)
	}
}

func TestValidateOUPathSegments_BothValid(t *testing.T) {
	errs := ValidateOUPathSegments("nic", "ministry-health")
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func hasFieldError(errs []resources.FieldError, field string) bool {
	for _, e := range errs {
		if e.Field == field {
			return true
		}
	}
	return false
}
