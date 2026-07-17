package validation

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateProject_ValidNames(t *testing.T) {
	valid := []string{"a", "a1", "a-b", strings.Repeat("a", 63)}
	for _, name := range valid {
		errs := ValidateProject(resources.Project{
			Metadata: resources.Metadata{Name: name},
			Spec: resources.ProjectSpec{
				OrganizationName:     "nic",
				OrganizationUnitName: "ministry-health",
				TenantName:           "payments",
			},
		})
		if len(errs) != 0 {
			t.Errorf("name %q: got errors %v, want none", name, errs)
		}
	}
}

func TestValidateProject_EmptyName(t *testing.T) {
	errs := ValidateProject(validProjectWithName(""))
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateProject_InvalidNameFormat(t *testing.T) {
	invalid := []string{"ABC", "a b", "-abc", "abc-", "a.b", "_abc"}
	for _, name := range invalid {
		errs := ValidateProject(validProjectWithName(name))
		if !hasFieldError(errs, "metadata.name") {
			t.Fatalf("name %q: got %v, want metadata.name error", name, errs)
		}
	}
}

func TestValidateProject_NameTooLong(t *testing.T) {
	errs := ValidateProject(validProjectWithName(strings.Repeat("a", 64)))
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateProject_EmptyOrganizationName(t *testing.T) {
	project := validProject()
	project.Spec.OrganizationName = ""
	errs := ValidateProject(project)
	if !hasFieldError(errs, "spec.organizationName") {
		t.Fatalf("got %v, want spec.organizationName error", errs)
	}
}

func TestValidateProject_InvalidOrganizationNameFormat(t *testing.T) {
	project := validProject()
	project.Spec.OrganizationName = "NIC_Org"
	errs := ValidateProject(project)
	if !hasFieldError(errs, "spec.organizationName") {
		t.Fatalf("got %v, want spec.organizationName error", errs)
	}
}

func TestValidateProject_OrganizationNameTooLong(t *testing.T) {
	project := validProject()
	project.Spec.OrganizationName = strings.Repeat("a", 64)
	errs := ValidateProject(project)
	if !hasFieldError(errs, "spec.organizationName") {
		t.Fatalf("got %v, want spec.organizationName error", errs)
	}
}

func TestValidateProject_EmptyOrganizationUnitName(t *testing.T) {
	project := validProject()
	project.Spec.OrganizationUnitName = ""
	errs := ValidateProject(project)
	if !hasFieldError(errs, "spec.organizationUnitName") {
		t.Fatalf("got %v, want spec.organizationUnitName error", errs)
	}
}

func TestValidateProject_InvalidOrganizationUnitNameFormat(t *testing.T) {
	project := validProject()
	project.Spec.OrganizationUnitName = "Health Unit"
	errs := ValidateProject(project)
	if !hasFieldError(errs, "spec.organizationUnitName") {
		t.Fatalf("got %v, want spec.organizationUnitName error", errs)
	}
}

func TestValidateProject_OrganizationUnitNameTooLong(t *testing.T) {
	project := validProject()
	project.Spec.OrganizationUnitName = strings.Repeat("a", 64)
	errs := ValidateProject(project)
	if !hasFieldError(errs, "spec.organizationUnitName") {
		t.Fatalf("got %v, want spec.organizationUnitName error", errs)
	}
}

func TestValidateProject_EmptyTenantName(t *testing.T) {
	project := validProject()
	project.Spec.TenantName = ""
	errs := ValidateProject(project)
	if !hasFieldError(errs, "spec.tenantName") {
		t.Fatalf("got %v, want spec.tenantName error", errs)
	}
}

func TestValidateProject_InvalidTenantNameFormat(t *testing.T) {
	project := validProject()
	project.Spec.TenantName = "Payments Tenant"
	errs := ValidateProject(project)
	if !hasFieldError(errs, "spec.tenantName") {
		t.Fatalf("got %v, want spec.tenantName error", errs)
	}
}

func TestValidateProject_TenantNameTooLong(t *testing.T) {
	project := validProject()
	project.Spec.TenantName = strings.Repeat("a", 64)
	errs := ValidateProject(project)
	if !hasFieldError(errs, "spec.tenantName") {
		t.Fatalf("got %v, want spec.tenantName error", errs)
	}
}

func TestValidateProjectPathSegments_InvalidOrgNameMapsToSpecOrganizationName(t *testing.T) {
	errs := ValidateProjectPathSegments("Invalid Org", "ministry-health", "payments", "prod")
	if !hasFieldError(errs, "spec.organizationName") {
		t.Fatalf("got %v, want spec.organizationName error", errs)
	}
}

func TestValidateProjectPathSegments_InvalidOUNameMapsToSpecOrganizationUnitName(t *testing.T) {
	errs := ValidateProjectPathSegments("nic", "Invalid OU", "payments", "prod")
	if !hasFieldError(errs, "spec.organizationUnitName") {
		t.Fatalf("got %v, want spec.organizationUnitName error", errs)
	}
}

func TestValidateProjectPathSegments_InvalidTenantNameMapsToSpecTenantName(t *testing.T) {
	errs := ValidateProjectPathSegments("nic", "ministry-health", "Invalid Tenant", "prod")
	if !hasFieldError(errs, "spec.tenantName") {
		t.Fatalf("got %v, want spec.tenantName error", errs)
	}
}

func TestValidateProjectPathSegments_InvalidNameMapsToMetadataName(t *testing.T) {
	errs := ValidateProjectPathSegments("nic", "ministry-health", "payments", "Invalid Project")
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateProjectPathSegments_AllValid(t *testing.T) {
	errs := ValidateProjectPathSegments("nic", "ministry-health", "payments", "prod")
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func validProject() resources.Project {
	return validProjectWithName("prod")
}

func validProjectWithName(name string) resources.Project {
	return resources.Project{
		Metadata: resources.Metadata{Name: name},
		Spec: resources.ProjectSpec{
			OrganizationName:     "nic",
			OrganizationUnitName: "ministry-health",
			TenantName:           "payments",
		},
	}
}
