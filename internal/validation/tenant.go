package validation

import (
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// ValidateTenant is a pure function. It validates all user-authored Tenant
// identity and parent-reference fields without performing I/O or registry
// lookups. Returns nil if the resource is valid.
func ValidateTenant(t resources.Tenant) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateName(t.Metadata.Name)...)
	errs = append(errs, validateOrganizationName(t.Spec.OrganizationName)...)
	errs = append(errs, validateOrganizationUnitName(t.Spec.OrganizationUnitName)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// ValidateTenantPathSegments validates Tenant URL path segments before a
// registry lookup. Path validation maps public fields to the Tenant body
// fields that the segments represent.
func ValidateTenantPathSegments(orgName, ouName, name string) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateOrganizationName(orgName)...)
	errs = append(errs, validateOrganizationUnitName(ouName)...)
	errs = append(errs, validateName(name)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// validateOrganizationUnitName validates spec.organizationUnitName using the
// same DNS-label rules as metadata.name.
func validateOrganizationUnitName(ouName string) []resources.FieldError {
	if ouName == "" {
		return []resources.FieldError{{
			Field:   "spec.organizationUnitName",
			Message: "organizationUnitName is required",
		}}
	}
	if len(ouName) > 63 {
		return []resources.FieldError{{
			Field:   "spec.organizationUnitName",
			Message: "organizationUnitName must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(ouName) {
		return []resources.FieldError{{
			Field:   "spec.organizationUnitName",
			Message: "organizationUnitName must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}
