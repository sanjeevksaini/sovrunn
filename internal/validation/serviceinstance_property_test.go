package validation

import (
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Feature: serviceinstance-servicebinding, Property 1: ValidateServiceInstance accepts valid DNS-label names
func TestProperty_ValidateServiceInstance_ValidNames(t *testing.T) {
	f := func(name, orgRef, ouRef, tenantRef, projectRef, classRef, planRef string) bool {
		if !isValidDNSLabel(name) ||
			!isValidDNSLabel(orgRef) ||
			!isValidDNSLabel(tenantRef) ||
			!isValidDNSLabel(projectRef) ||
			!isValidDNSLabel(classRef) ||
			!isValidDNSLabel(planRef) {
			return true
		}
		// organizationUnitRef is optional: empty is accepted; non-empty must be a DNS label.
		if ouRef != "" && !isValidDNSLabel(ouRef) {
			return true
		}
		errs := ValidateServiceInstance(resources.ServiceInstance{
			Metadata: resources.Metadata{Name: name},
			Spec: resources.ServiceInstanceSpec{
				OrganizationRef:     orgRef,
				OrganizationUnitRef: ouRef,
				TenantRef:           tenantRef,
				ProjectRef:          projectRef,
				ServiceClassRef:     classRef,
				ServicePlanRef:      planRef,
			},
		})
		return len(errs) == 0
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceinstance-servicebinding, Property 2: ValidateServiceInstance rejects arbitrary invalid strings
func TestProperty_ValidateServiceInstance_InvalidNames(t *testing.T) {
	f := func(s string) bool {
		if isValidDNSLabel(s) {
			return true
		}

		si := validServiceInstance()
		si.Metadata.Name = s
		if !hasFieldError(ValidateServiceInstance(si), "metadata.name") {
			return false
		}

		si = validServiceInstance()
		si.Spec.OrganizationRef = s
		if !hasFieldError(ValidateServiceInstance(si), "spec.organizationRef") {
			return false
		}

		si = validServiceInstance()
		si.Spec.TenantRef = s
		if !hasFieldError(ValidateServiceInstance(si), "spec.tenantRef") {
			return false
		}

		si = validServiceInstance()
		si.Spec.ProjectRef = s
		if !hasFieldError(ValidateServiceInstance(si), "spec.projectRef") {
			return false
		}

		si = validServiceInstance()
		si.Spec.ServiceClassRef = s
		if !hasFieldError(ValidateServiceInstance(si), "spec.serviceClassRef") {
			return false
		}

		si = validServiceInstance()
		si.Spec.ServicePlanRef = s
		if !hasFieldError(ValidateServiceInstance(si), "spec.servicePlanRef") {
			return false
		}

		// organizationUnitRef is optional: empty is accepted; non-empty invalid values are rejected.
		si = validServiceInstance()
		si.Spec.OrganizationUnitRef = s
		if s == "" {
			return !hasFieldError(ValidateServiceInstance(si), "spec.organizationUnitRef")
		}
		return hasFieldError(ValidateServiceInstance(si), "spec.organizationUnitRef")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
