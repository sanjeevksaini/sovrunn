package validation

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateServiceInstance_Valid(t *testing.T) {
	valid := []string{"a", "a1", "a-b", strings.Repeat("a", 63)}
	for _, name := range valid {
		errs := ValidateServiceInstance(validServiceInstanceWithName(name))
		if len(errs) != 0 {
			t.Errorf("name %q: got errors %v, want none", name, errs)
		}
	}
}

func TestValidateServiceInstance_EmptyName(t *testing.T) {
	errs := ValidateServiceInstance(validServiceInstanceWithName(""))
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateServiceInstance_InvalidNameFormat(t *testing.T) {
	invalid := []string{"ABC", "a b", "-abc", "abc-", "a.b", "_abc"}
	for _, name := range invalid {
		errs := ValidateServiceInstance(validServiceInstanceWithName(name))
		if !hasFieldError(errs, "metadata.name") {
			t.Fatalf("name %q: got %v, want metadata.name error", name, errs)
		}
	}
}

func TestValidateServiceInstance_NameTooLong(t *testing.T) {
	errs := ValidateServiceInstance(validServiceInstanceWithName(strings.Repeat("a", 64)))
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateServiceInstance_EmptyOrganizationRef(t *testing.T) {
	si := validServiceInstance()
	si.Spec.OrganizationRef = ""
	errs := ValidateServiceInstance(si)
	if !hasFieldError(errs, "spec.organizationRef") {
		t.Fatalf("got %v, want spec.organizationRef error", errs)
	}
}

func TestValidateServiceInstance_EmptyTenantRef(t *testing.T) {
	si := validServiceInstance()
	si.Spec.TenantRef = ""
	errs := ValidateServiceInstance(si)
	if !hasFieldError(errs, "spec.tenantRef") {
		t.Fatalf("got %v, want spec.tenantRef error", errs)
	}
}

func TestValidateServiceInstance_EmptyProjectRef(t *testing.T) {
	si := validServiceInstance()
	si.Spec.ProjectRef = ""
	errs := ValidateServiceInstance(si)
	if !hasFieldError(errs, "spec.projectRef") {
		t.Fatalf("got %v, want spec.projectRef error", errs)
	}
}

func TestValidateServiceInstance_EmptyServiceClassRef(t *testing.T) {
	si := validServiceInstance()
	si.Spec.ServiceClassRef = ""
	errs := ValidateServiceInstance(si)
	if !hasFieldError(errs, "spec.serviceClassRef") {
		t.Fatalf("got %v, want spec.serviceClassRef error", errs)
	}
}

func TestValidateServiceInstance_EmptyServicePlanRef(t *testing.T) {
	si := validServiceInstance()
	si.Spec.ServicePlanRef = ""
	errs := ValidateServiceInstance(si)
	if !hasFieldError(errs, "spec.servicePlanRef") {
		t.Fatalf("got %v, want spec.servicePlanRef error", errs)
	}
}

func TestValidateServiceInstance_EmptyOrganizationUnitRefAccepted(t *testing.T) {
	si := validServiceInstance()
	si.Spec.OrganizationUnitRef = ""
	errs := ValidateServiceInstance(si)
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors for empty organizationUnitRef", errs)
	}
}

func TestValidateServiceInstance_InvalidOrganizationUnitRef(t *testing.T) {
	invalid := []string{"ABC", "a b", "-abc", "abc-", strings.Repeat("a", 64)}
	for _, ref := range invalid {
		si := validServiceInstance()
		si.Spec.OrganizationUnitRef = ref
		errs := ValidateServiceInstance(si)
		if !hasFieldError(errs, "spec.organizationUnitRef") {
			t.Fatalf("organizationUnitRef %q: got %v, want spec.organizationUnitRef error", ref, errs)
		}
	}
}

func TestValidateServiceInstance_ValidOrganizationUnitRef(t *testing.T) {
	si := validServiceInstance()
	si.Spec.OrganizationUnitRef = "ministry-health"
	errs := ValidateServiceInstance(si)
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func TestValidateServiceInstance_ParametersNilAccepted(t *testing.T) {
	si := validServiceInstance()
	si.Spec.Parameters = nil
	errs := ValidateServiceInstance(si)
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors for nil parameters", errs)
	}
}

func TestValidateServiceInstance_ParametersEmptyMapAccepted(t *testing.T) {
	si := validServiceInstance()
	si.Spec.Parameters = map[string]string{}
	errs := ValidateServiceInstance(si)
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors for empty parameters map", errs)
	}
}

func TestValidateServiceInstance_InvalidRequiredRefs(t *testing.T) {
	cases := []struct {
		field string
		set   func(*resources.ServiceInstance, string)
	}{
		{"spec.organizationRef", func(si *resources.ServiceInstance, v string) { si.Spec.OrganizationRef = v }},
		{"spec.tenantRef", func(si *resources.ServiceInstance, v string) { si.Spec.TenantRef = v }},
		{"spec.projectRef", func(si *resources.ServiceInstance, v string) { si.Spec.ProjectRef = v }},
		{"spec.serviceClassRef", func(si *resources.ServiceInstance, v string) { si.Spec.ServiceClassRef = v }},
		{"spec.servicePlanRef", func(si *resources.ServiceInstance, v string) { si.Spec.ServicePlanRef = v }},
	}
	invalid := []string{"ABC", "a b", "-abc", "abc-", strings.Repeat("a", 64)}
	for _, tc := range cases {
		for _, v := range invalid {
			si := validServiceInstance()
			tc.set(&si, v)
			errs := ValidateServiceInstance(si)
			if !hasFieldError(errs, tc.field) {
				t.Fatalf("%s %q: got %v, want %s error", tc.field, v, errs, tc.field)
			}
		}
	}
}

func TestValidateServiceInstancePathSegment_Valid(t *testing.T) {
	errs := ValidateServiceInstancePathSegment("pg-prod")
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func TestValidateServiceInstancePathSegment_Invalid(t *testing.T) {
	errs := ValidateServiceInstancePathSegment("Invalid Name")
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func validServiceInstance() resources.ServiceInstance {
	return validServiceInstanceWithName("pg-prod")
}

func validServiceInstanceWithName(name string) resources.ServiceInstance {
	return resources.ServiceInstance{
		Metadata: resources.Metadata{Name: name},
		Spec: resources.ServiceInstanceSpec{
			OrganizationRef: "nic",
			TenantRef:       "payments",
			ProjectRef:      "prod",
			ServiceClassRef: "postgresql",
			ServicePlanRef:  "small",
		},
	}
}
