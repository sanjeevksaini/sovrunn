package validation

import (
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Feature: organizationunit-resource, Property 1: ValidateOrganizationUnit accepts all valid DNS-label names
func TestProperty_ValidateOrganizationUnit_ValidNames(t *testing.T) {
	f := func(name, orgName string) bool {
		if !isValidDNSLabel(name) || !isValidDNSLabel(orgName) {
			return true
		}
		errs := ValidateOrganizationUnit(resources.OrganizationUnit{
			Metadata: resources.Metadata{Name: name},
			Spec:     resources.OrganizationUnitSpec{OrganizationName: orgName},
		})
		return len(errs) == 0
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: organizationunit-resource, Property 2: ValidateOrganizationUnit rejects invalid or arbitrary strings outside the DNS-label domain
func TestProperty_ValidateOrganizationUnit_InvalidNames(t *testing.T) {
	f := func(name string) bool {
		validOrgName := "nic"
		if isValidDNSLabel(name) {
			errs := ValidateOrganizationUnit(resources.OrganizationUnit{
				Metadata: resources.Metadata{Name: name},
				Spec:     resources.OrganizationUnitSpec{OrganizationName: validOrgName},
			})
			return len(errs) == 0
		}
		errs := ValidateOrganizationUnit(resources.OrganizationUnit{
			Metadata: resources.Metadata{Name: name},
			Spec:     resources.OrganizationUnitSpec{OrganizationName: validOrgName},
		})
		for _, e := range errs {
			if e.Field == "metadata.name" {
				return true
			}
		}
		return false
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: organizationunit-resource, Property 2: ValidateOrganizationUnit rejects invalid spec.organizationName values outside the DNS-label domain
func TestProperty_ValidateOrganizationUnit_InvalidOrganizationNames(t *testing.T) {
	f := func(orgName string) bool {
		validName := "ministry-health"
		if isValidDNSLabel(orgName) {
			errs := ValidateOrganizationUnit(resources.OrganizationUnit{
				Metadata: resources.Metadata{Name: validName},
				Spec:     resources.OrganizationUnitSpec{OrganizationName: orgName},
			})
			return len(errs) == 0
		}
		errs := ValidateOrganizationUnit(resources.OrganizationUnit{
			Metadata: resources.Metadata{Name: validName},
			Spec:     resources.OrganizationUnitSpec{OrganizationName: orgName},
		})
		for _, e := range errs {
			if e.Field == "spec.organizationName" {
				return true
			}
		}
		return false
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
