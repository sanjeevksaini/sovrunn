package validation

import (
	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// ValidateOrganizationUnit is a pure function. It validates all
// user-authored fields of ou and returns all FieldErrors found in a
// single call (does not stop at the first error). Returns nil if the
// resource is valid.
//
// It does not perform cross-registry checks such as parent Organization
// existence — that is the responsibility of the OUHandler in later tasks.
// It accepts no context.Context because it performs no I/O or
// cancellation-aware work.
func ValidateOrganizationUnit(ou resources.OrganizationUnit) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateName(ou.Metadata.Name)...)
	errs = append(errs, validateOrganizationName(ou.Spec.OrganizationName)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// ValidateOUPathSegments validates the organizationName and name path
// segments extracted from a URL path. Used by Get, Update, and Delete
// handlers before the registry lookup. Context-free because it performs
// no I/O.
func ValidateOUPathSegments(orgName, name string) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateName(name)...)
	errs = append(errs, validateOrganizationName(orgName)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// validateOrganizationName validates spec.organizationName using the same
// DNS-label rules as metadata.name, reusing the shared dnsLabelRe from
// organization.go.
func validateOrganizationName(orgName string) []resources.FieldError {
	if orgName == "" {
		return []resources.FieldError{{
			Field:   "spec.organizationName",
			Message: "organizationName is required",
		}}
	}
	if len(orgName) > 63 {
		return []resources.FieldError{{
			Field:   "spec.organizationName",
			Message: "organizationName must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(orgName) {
		return []resources.FieldError{{
			Field:   "spec.organizationName",
			Message: "organizationName must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}
