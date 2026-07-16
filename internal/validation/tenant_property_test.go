package validation

import (
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Feature: tenant-resource, Property 1: ValidateTenant accepts valid DNS-label names
func TestProperty_ValidateTenant_ValidNames(t *testing.T) {
	f := func(name, orgName, ouName string) bool {
		if !isValidDNSLabel(name) || !isValidDNSLabel(orgName) || !isValidDNSLabel(ouName) {
			return true
		}
		errs := ValidateTenant(resources.Tenant{
			Metadata: resources.Metadata{Name: name},
			Spec: resources.TenantSpec{
				OrganizationName:     orgName,
				OrganizationUnitName: ouName,
			},
		})
		return len(errs) == 0
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: tenant-resource, Property 2: ValidateTenant rejects invalid metadata.name values
func TestProperty_ValidateTenant_InvalidNames(t *testing.T) {
	f := func(name string) bool {
		if isValidDNSLabel(name) {
			errs := ValidateTenant(validTenantWithName(name))
			return len(errs) == 0
		}
		errs := ValidateTenant(validTenantWithName(name))
		return hasFieldError(errs, "metadata.name")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: tenant-resource, Property 2: ValidateTenant rejects invalid spec.organizationName values
func TestProperty_ValidateTenant_InvalidOrganizationNames(t *testing.T) {
	f := func(orgName string) bool {
		tnt := validTenant()
		tnt.Spec.OrganizationName = orgName
		errs := ValidateTenant(tnt)
		if isValidDNSLabel(orgName) {
			return len(errs) == 0
		}
		return hasFieldError(errs, "spec.organizationName")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: tenant-resource, Property 2: ValidateTenant rejects invalid spec.organizationUnitName values
func TestProperty_ValidateTenant_InvalidOrganizationUnitNames(t *testing.T) {
	f := func(ouName string) bool {
		tnt := validTenant()
		tnt.Spec.OrganizationUnitName = ouName
		errs := ValidateTenant(tnt)
		if isValidDNSLabel(ouName) {
			return len(errs) == 0
		}
		return hasFieldError(errs, "spec.organizationUnitName")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
