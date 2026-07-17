package validation

import (
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// forbiddenParamSubstrings are matched case-insensitively against each
// parameter KEY. The plain substring "key" is intentionally NOT listed; only
// the composite secret-bearing phrases below trigger rejection.
var forbiddenParamSubstrings = []string{
	"password", "secret", "token", "credential", "auth",
	"apikey", "accesskey", "secretkey", "privatekey",
}

// ValidateServicePlan is a pure function. It validates all user-authored
// ServicePlan fields without performing I/O or registry lookups. It returns all
// FieldErrors found in a single call and returns nil if the resource is valid.
func ValidateServicePlan(sp resources.ServicePlan) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateName(sp.Metadata.Name)...)
	errs = append(errs, validateServiceClassName(sp.Spec.ServiceClassName)...)
	errs = append(errs, validateTier(sp.Spec.Tier)...)
	errs = append(errs, validateLifecycle(sp.Spec.Lifecycle)...)
	errs = append(errs, validateParameters(sp.Spec.Parameters)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// ValidateServicePlanPathSegments validates the two ServicePlan URL path
// segments before a registry lookup. It maps the serviceClassName segment to
// spec.serviceClassName and the name segment to metadata.name.
func ValidateServicePlanPathSegments(serviceClassName, name string) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateServiceClassName(serviceClassName)...)
	errs = append(errs, validateName(name)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// validateServiceClassName validates spec.serviceClassName using the same
// DNS-label rules as metadata.name.
func validateServiceClassName(serviceClassName string) []resources.FieldError {
	if serviceClassName == "" {
		return []resources.FieldError{{
			Field:   "spec.serviceClassName",
			Message: "serviceClassName is required",
		}}
	}
	if len(serviceClassName) > 63 {
		return []resources.FieldError{{
			Field:   "spec.serviceClassName",
			Message: "serviceClassName must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(serviceClassName) {
		return []resources.FieldError{{
			Field:   "spec.serviceClassName",
			Message: "serviceClassName must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}

// validateTier checks that spec.tier is present and a known tier value.
func validateTier(tier string) []resources.FieldError {
	if tier == "" {
		return []resources.FieldError{{
			Field:   "spec.tier",
			Message: "tier is required",
		}}
	}
	switch tier {
	case resources.TierDev,
		resources.TierSmall,
		resources.TierMedium,
		resources.TierLarge,
		resources.TierProduction,
		resources.TierCustom:
		return nil
	}
	return []resources.FieldError{{
		Field:   "spec.tier",
		Message: "tier must be one of: Dev, Small, Medium, Large, Production, Custom",
	}}
}

// validateParameters rejects parameter keys that carry secret-bearing
// substrings (case-insensitive). It inspects keys only and never stores or
// echoes the offending value. A single FieldError is returned regardless of how
// many keys offend.
func validateParameters(parameters map[string]string) []resources.FieldError {
	for key := range parameters {
		if isForbiddenParamKey(key) {
			return []resources.FieldError{{
				Field:   "spec.parameters",
				Message: "parameters must not contain secret-bearing keys",
			}}
		}
	}
	return nil
}

// isForbiddenParamKey lowercases the key once, then checks containment against
// the forbidden substrings. The bare substring "key" is allowed.
func isForbiddenParamKey(key string) bool {
	lk := strings.ToLower(key)
	for _, s := range forbiddenParamSubstrings {
		if strings.Contains(lk, s) {
			return true
		}
	}
	return false
}
