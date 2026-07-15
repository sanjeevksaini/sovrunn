package validation

import (
	"context"
	"regexp"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// dnsLabelRe is compiled once at package initialisation to avoid
// per-request allocation in a hot validation path.
var dnsLabelRe = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

// ValidateOrganization is a pure function. It validates all user-authored
// fields of org and returns all FieldErrors found in a single call
// (does not stop at the first error). Returns nil if the resource is valid.
func ValidateOrganization(ctx context.Context, org resources.Organization) []resources.FieldError {
	return validateName(org.Metadata.Name)
}

// ValidateNamePath validates a name extracted from a URL path segment.
// Used by Get, Update, and Delete handlers before the registry lookup.
func ValidateNamePath(ctx context.Context, name string) []resources.FieldError {
	return validateName(name)
}

func validateName(name string) []resources.FieldError {
	if name == "" {
		return []resources.FieldError{{
			Field:   "metadata.name",
			Message: "name is required",
		}}
	}
	if len(name) > 63 {
		return []resources.FieldError{{
			Field:   "metadata.name",
			Message: "name must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(name) {
		return []resources.FieldError{{
			Field:   "metadata.name",
			Message: "name must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}
