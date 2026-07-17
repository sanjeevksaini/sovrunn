package validation

import "github.com/sanjeevksaini/sovrunn/internal/resources"

// ValidateProject is a pure function. It validates all user-authored Project
// identity and parent-reference fields without performing I/O or registry
// lookups. Returns nil if the resource is valid.
func ValidateProject(p resources.Project) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateName(p.Metadata.Name)...)
	errs = append(errs, validateOrganizationName(p.Spec.OrganizationName)...)
	errs = append(errs, validateOrganizationUnitName(p.Spec.OrganizationUnitName)...)
	errs = append(errs, validateTenantName(p.Spec.TenantName)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// ValidateProjectPathSegments validates Project URL path segments before a
// registry lookup. Path validation maps public fields to the Project body
// fields that the segments represent.
func ValidateProjectPathSegments(orgName, ouName, tenantName, name string) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateOrganizationName(orgName)...)
	errs = append(errs, validateOrganizationUnitName(ouName)...)
	errs = append(errs, validateTenantName(tenantName)...)
	errs = append(errs, validateName(name)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// validateTenantName validates spec.tenantName using the same DNS-label rules
// as metadata.name.
func validateTenantName(tenantName string) []resources.FieldError {
	if tenantName == "" {
		return []resources.FieldError{{
			Field:   "spec.tenantName",
			Message: "tenantName is required",
		}}
	}
	if len(tenantName) > 63 {
		return []resources.FieldError{{
			Field:   "spec.tenantName",
			Message: "tenantName must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(tenantName) {
		return []resources.FieldError{{
			Field:   "spec.tenantName",
			Message: "tenantName must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}
