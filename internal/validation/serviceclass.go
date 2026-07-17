package validation

import "github.com/sanjeevksaini/sovrunn/internal/resources"

// ValidateServiceClass is a pure function. It validates all user-authored
// ServiceClass fields without performing I/O or registry lookups. It returns
// all FieldErrors found in a single call and returns nil if the resource is
// valid.
func ValidateServiceClass(sc resources.ServiceClass) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateName(sc.Metadata.Name)...)
	errs = append(errs, validateCategory(sc.Spec.Category)...)
	errs = append(errs, validateLifecycle(sc.Spec.Lifecycle)...)
	errs = append(errs, validateDefaultPlanName(sc.Spec.DefaultPlanName)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// ValidateServiceClassPathSegment validates the single ServiceClass URL path
// segment before a registry lookup. It only checks the name segment and maps
// it to metadata.name.
func ValidateServiceClassPathSegment(name string) []resources.FieldError {
	return validateName(name)
}

// validateCategory checks that spec.category is present and a known category.
func validateCategory(category string) []resources.FieldError {
	if category == "" {
		return []resources.FieldError{{
			Field:   "spec.category",
			Message: "category is required",
		}}
	}
	switch category {
	case resources.CategoryDatabase,
		resources.CategoryCache,
		resources.CategoryObjectStorage,
		resources.CategoryStream,
		resources.CategoryGateway,
		resources.CategoryFunction,
		resources.CategoryAnalytics,
		resources.CategoryOther:
		return nil
	}
	return []resources.FieldError{{
		Field:   "spec.category",
		Message: "category must be one of: Database, Cache, ObjectStorage, Stream, Gateway, Function, Analytics, Other",
	}}
}

// validateLifecycle checks that spec.lifecycle is present and a known lifecycle
// value. It is shared in spirit with ServicePlan but kept local to ServiceClass
// in this task.
func validateLifecycle(lifecycle string) []resources.FieldError {
	if lifecycle == "" {
		return []resources.FieldError{{
			Field:   "spec.lifecycle",
			Message: "lifecycle is required",
		}}
	}
	switch lifecycle {
	case resources.LifecyclePreview,
		resources.LifecycleActive,
		resources.LifecycleDeprecated,
		resources.LifecycleRetired:
		return nil
	}
	return []resources.FieldError{{
		Field:   "spec.lifecycle",
		Message: "lifecycle must be one of: Preview, Active, Deprecated, Retired",
	}}
}

// validateDefaultPlanName validates spec.defaultPlanName only when present. It
// does NOT verify that the referenced plan exists. An empty value is allowed.
func validateDefaultPlanName(defaultPlanName string) []resources.FieldError {
	if defaultPlanName == "" {
		return nil
	}
	if len(defaultPlanName) > 63 {
		return []resources.FieldError{{
			Field:   "spec.defaultPlanName",
			Message: "defaultPlanName must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(defaultPlanName) {
		return []resources.FieldError{{
			Field:   "spec.defaultPlanName",
			Message: "defaultPlanName must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}
