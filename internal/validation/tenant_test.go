package validation

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateTenant_ValidNames(t *testing.T) {
	valid := []string{"a", "a1", "a-b", strings.Repeat("a", 63)}
	for _, name := range valid {
		errs := ValidateTenant(resources.Tenant{
			Metadata: resources.Metadata{Name: name},
			Spec: resources.TenantSpec{
				OrganizationName:     "nic",
				OrganizationUnitName: "ministry-health",
			},
		})
		if len(errs) != 0 {
			t.Errorf("name %q: got errors %v, want none", name, errs)
		}
	}
}

func TestValidateTenant_EmptyName(t *testing.T) {
	errs := ValidateTenant(validTenantWithName(""))
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateTenant_InvalidNameFormat(t *testing.T) {
	invalid := []string{"ABC", "a b", "-abc", "abc-", "a.b", "_abc"}
	for _, name := range invalid {
		errs := ValidateTenant(validTenantWithName(name))
		if !hasFieldError(errs, "metadata.name") {
			t.Fatalf("name %q: got %v, want metadata.name error", name, errs)
		}
	}
}

func TestValidateTenant_NameTooLong(t *testing.T) {
	errs := ValidateTenant(validTenantWithName(strings.Repeat("a", 64)))
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateTenant_EmptyOrganizationName(t *testing.T) {
	tnt := validTenant()
	tnt.Spec.OrganizationName = ""
	errs := ValidateTenant(tnt)
	if !hasFieldError(errs, "spec.organizationName") {
		t.Fatalf("got %v, want spec.organizationName error", errs)
	}
}

func TestValidateTenant_InvalidOrganizationNameFormat(t *testing.T) {
	tnt := validTenant()
	tnt.Spec.OrganizationName = "NIC_Org"
	errs := ValidateTenant(tnt)
	if !hasFieldError(errs, "spec.organizationName") {
		t.Fatalf("got %v, want spec.organizationName error", errs)
	}
}

func TestValidateTenant_OrganizationNameTooLong(t *testing.T) {
	tnt := validTenant()
	tnt.Spec.OrganizationName = strings.Repeat("a", 64)
	errs := ValidateTenant(tnt)
	if !hasFieldError(errs, "spec.organizationName") {
		t.Fatalf("got %v, want spec.organizationName error", errs)
	}
}

func TestValidateTenant_EmptyOrganizationUnitName(t *testing.T) {
	tnt := validTenant()
	tnt.Spec.OrganizationUnitName = ""
	errs := ValidateTenant(tnt)
	if !hasFieldError(errs, "spec.organizationUnitName") {
		t.Fatalf("got %v, want spec.organizationUnitName error", errs)
	}
}

func TestValidateTenant_InvalidOrganizationUnitNameFormat(t *testing.T) {
	tnt := validTenant()
	tnt.Spec.OrganizationUnitName = "Health Unit"
	errs := ValidateTenant(tnt)
	if !hasFieldError(errs, "spec.organizationUnitName") {
		t.Fatalf("got %v, want spec.organizationUnitName error", errs)
	}
}

func TestValidateTenant_OrganizationUnitNameTooLong(t *testing.T) {
	tnt := validTenant()
	tnt.Spec.OrganizationUnitName = strings.Repeat("a", 64)
	errs := ValidateTenant(tnt)
	if !hasFieldError(errs, "spec.organizationUnitName") {
		t.Fatalf("got %v, want spec.organizationUnitName error", errs)
	}
}

func TestValidateTenantPathSegments_InvalidOrgNameMapsToSpecOrganizationName(t *testing.T) {
	errs := ValidateTenantPathSegments("Invalid Org", "ministry-health", "prod")
	if !hasFieldError(errs, "spec.organizationName") {
		t.Fatalf("got %v, want spec.organizationName error", errs)
	}
}

func TestValidateTenantPathSegments_InvalidOUNameMapsToSpecOrganizationUnitName(t *testing.T) {
	errs := ValidateTenantPathSegments("nic", "Invalid OU", "prod")
	if !hasFieldError(errs, "spec.organizationUnitName") {
		t.Fatalf("got %v, want spec.organizationUnitName error", errs)
	}
}

func TestValidateTenantPathSegments_InvalidNameMapsToMetadataName(t *testing.T) {
	errs := ValidateTenantPathSegments("nic", "ministry-health", "Invalid Tenant")
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateTenantPathSegments_AllValid(t *testing.T) {
	errs := ValidateTenantPathSegments("nic", "ministry-health", "prod")
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func validTenant() resources.Tenant {
	return validTenantWithName("prod")
}

func validTenantWithName(name string) resources.Tenant {
	return resources.Tenant{
		Metadata: resources.Metadata{Name: name},
		Spec: resources.TenantSpec{
			OrganizationName:     "nic",
			OrganizationUnitName: "ministry-health",
		},
	}
}
