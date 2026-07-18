package validation

import "github.com/sanjeevksaini/sovrunn/internal/resources"

// ValidateServiceInstance is a pure function. It validates all user-authored
// ServiceInstance identity and reference fields without performing I/O or
// registry lookups. Returns nil if the resource is valid.
// spec.parameters is not validated in Phase 1.
func ValidateServiceInstance(si resources.ServiceInstance) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateName(si.Metadata.Name)...)
	errs = append(errs, validateOrganizationRef(si.Spec.OrganizationRef)...)
	errs = append(errs, validateOrganizationUnitRef(si.Spec.OrganizationUnitRef)...)
	errs = append(errs, validateTenantRef(si.Spec.TenantRef)...)
	errs = append(errs, validateProjectRef(si.Spec.ProjectRef)...)
	errs = append(errs, validateServiceClassRef(si.Spec.ServiceClassRef)...)
	errs = append(errs, validateServicePlanRef(si.Spec.ServicePlanRef)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// ValidateServiceInstancePathSegment validates the single ServiceInstance URL
// path segment before a registry lookup. It only checks the name segment and
// maps it to metadata.name.
func ValidateServiceInstancePathSegment(name string) []resources.FieldError {
	return validateName(name)
}

// validateOrganizationRef validates spec.organizationRef using the same
// DNS-label rules as metadata.name. Existence is NOT checked here.
func validateOrganizationRef(organizationRef string) []resources.FieldError {
	if organizationRef == "" {
		return []resources.FieldError{{
			Field:   "spec.organizationRef",
			Message: "organizationRef is required",
		}}
	}
	if len(organizationRef) > 63 {
		return []resources.FieldError{{
			Field:   "spec.organizationRef",
			Message: "organizationRef must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(organizationRef) {
		return []resources.FieldError{{
			Field:   "spec.organizationRef",
			Message: "organizationRef must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}

// validateOrganizationUnitRef validates spec.organizationUnitRef only when
// present. An empty value is allowed (optional field).
func validateOrganizationUnitRef(organizationUnitRef string) []resources.FieldError {
	if organizationUnitRef == "" {
		return nil
	}
	if len(organizationUnitRef) > 63 {
		return []resources.FieldError{{
			Field:   "spec.organizationUnitRef",
			Message: "organizationUnitRef must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(organizationUnitRef) {
		return []resources.FieldError{{
			Field:   "spec.organizationUnitRef",
			Message: "organizationUnitRef must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}

// validateTenantRef validates spec.tenantRef using the same DNS-label rules as
// metadata.name. Existence is NOT checked here.
func validateTenantRef(tenantRef string) []resources.FieldError {
	if tenantRef == "" {
		return []resources.FieldError{{
			Field:   "spec.tenantRef",
			Message: "tenantRef is required",
		}}
	}
	if len(tenantRef) > 63 {
		return []resources.FieldError{{
			Field:   "spec.tenantRef",
			Message: "tenantRef must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(tenantRef) {
		return []resources.FieldError{{
			Field:   "spec.tenantRef",
			Message: "tenantRef must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}

// validateProjectRef validates spec.projectRef using the same DNS-label rules as
// metadata.name. Existence is NOT checked here.
func validateProjectRef(projectRef string) []resources.FieldError {
	if projectRef == "" {
		return []resources.FieldError{{
			Field:   "spec.projectRef",
			Message: "projectRef is required",
		}}
	}
	if len(projectRef) > 63 {
		return []resources.FieldError{{
			Field:   "spec.projectRef",
			Message: "projectRef must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(projectRef) {
		return []resources.FieldError{{
			Field:   "spec.projectRef",
			Message: "projectRef must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}

// validateServicePlanRef validates spec.servicePlanRef using the same DNS-label
// rules as metadata.name. Existence is NOT checked here.
func validateServicePlanRef(servicePlanRef string) []resources.FieldError {
	if servicePlanRef == "" {
		return []resources.FieldError{{
			Field:   "spec.servicePlanRef",
			Message: "servicePlanRef is required",
		}}
	}
	if len(servicePlanRef) > 63 {
		return []resources.FieldError{{
			Field:   "spec.servicePlanRef",
			Message: "servicePlanRef must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(servicePlanRef) {
		return []resources.FieldError{{
			Field:   "spec.servicePlanRef",
			Message: "servicePlanRef must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}
