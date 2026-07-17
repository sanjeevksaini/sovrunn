package validation

import (
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Feature: project-resource, Property 1: ValidateProject accepts valid DNS-label names
func TestProperty_ValidateProject_ValidNames(t *testing.T) {
	f := func(name, orgName, ouName, tenantName string) bool {
		if !isValidDNSLabel(name) ||
			!isValidDNSLabel(orgName) ||
			!isValidDNSLabel(ouName) ||
			!isValidDNSLabel(tenantName) {
			return true
		}
		errs := ValidateProject(resources.Project{
			Metadata: resources.Metadata{Name: name},
			Spec: resources.ProjectSpec{
				OrganizationName:     orgName,
				OrganizationUnitName: ouName,
				TenantName:           tenantName,
			},
		})
		return len(errs) == 0
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: project-resource, Property 2: ValidateProject rejects invalid metadata.name values
func TestProperty_ValidateProject_InvalidNames(t *testing.T) {
	f := func(name string) bool {
		if isValidDNSLabel(name) {
			errs := ValidateProject(validProjectWithName(name))
			return len(errs) == 0
		}
		errs := ValidateProject(validProjectWithName(name))
		return hasFieldError(errs, "metadata.name")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: project-resource, Property 2: ValidateProject rejects invalid spec.organizationName values
func TestProperty_ValidateProject_InvalidOrganizationNames(t *testing.T) {
	f := func(orgName string) bool {
		project := validProject()
		project.Spec.OrganizationName = orgName
		errs := ValidateProject(project)
		if isValidDNSLabel(orgName) {
			return len(errs) == 0
		}
		return hasFieldError(errs, "spec.organizationName")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: project-resource, Property 2: ValidateProject rejects invalid spec.organizationUnitName values
func TestProperty_ValidateProject_InvalidOrganizationUnitNames(t *testing.T) {
	f := func(ouName string) bool {
		project := validProject()
		project.Spec.OrganizationUnitName = ouName
		errs := ValidateProject(project)
		if isValidDNSLabel(ouName) {
			return len(errs) == 0
		}
		return hasFieldError(errs, "spec.organizationUnitName")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: project-resource, Property 2: ValidateProject rejects invalid spec.tenantName values
func TestProperty_ValidateProject_InvalidTenantNames(t *testing.T) {
	f := func(tenantName string) bool {
		project := validProject()
		project.Spec.TenantName = tenantName
		errs := ValidateProject(project)
		if isValidDNSLabel(tenantName) {
			return len(errs) == 0
		}
		return hasFieldError(errs, "spec.tenantName")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
