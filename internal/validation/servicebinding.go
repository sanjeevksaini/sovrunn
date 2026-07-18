package validation

import "github.com/sanjeevksaini/sovrunn/internal/resources"

// ValidateServiceBinding is a pure function. It validates all user-authored
// ServiceBinding identity and reference fields without performing I/O or
// registry lookups. Returns nil if the resource is valid.
func ValidateServiceBinding(sb resources.ServiceBinding) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateName(sb.Metadata.Name)...)
	errs = append(errs, validateServiceInstanceRef(sb.Spec.ServiceInstanceRef)...)
	errs = append(errs, validateConsumerRef(sb.Spec.ConsumerRef)...)
	errs = append(errs, validateBindingType(sb.Spec.BindingType)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// ValidateServiceBindingPathSegment validates the single ServiceBinding URL
// path segment before a registry lookup. It only checks the name segment and
// maps it to metadata.name.
func ValidateServiceBindingPathSegment(name string) []resources.FieldError {
	return validateName(name)
}

// validateServiceInstanceRef validates spec.serviceInstanceRef using the same
// DNS-label rules as metadata.name. Existence is NOT checked here.
func validateServiceInstanceRef(serviceInstanceRef string) []resources.FieldError {
	if serviceInstanceRef == "" {
		return []resources.FieldError{{
			Field:   "spec.serviceInstanceRef",
			Message: "serviceInstanceRef is required",
		}}
	}
	if len(serviceInstanceRef) > 63 {
		return []resources.FieldError{{
			Field:   "spec.serviceInstanceRef",
			Message: "serviceInstanceRef must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(serviceInstanceRef) {
		return []resources.FieldError{{
			Field:   "spec.serviceInstanceRef",
			Message: "serviceInstanceRef must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}

// validateConsumerRef validates that spec.consumerRef is non-nil and that its
// kind and name fields meet Phase 1 rules. Kind has no enum restriction.
func validateConsumerRef(consumerRef *resources.ConsumerRef) []resources.FieldError {
	if consumerRef == nil {
		return []resources.FieldError{{
			Field:   "spec.consumerRef",
			Message: "consumerRef is required",
		}}
	}
	var errs []resources.FieldError
	if consumerRef.Kind == "" {
		errs = append(errs, resources.FieldError{
			Field:   "spec.consumerRef.kind",
			Message: "consumerRef.kind is required",
		})
	}
	errs = append(errs, validateConsumerRefName(consumerRef.Name)...)
	return errs
}

// validateConsumerRefName validates spec.consumerRef.name using the same
// DNS-label rules as metadata.name.
func validateConsumerRefName(name string) []resources.FieldError {
	if name == "" {
		return []resources.FieldError{{
			Field:   "spec.consumerRef.name",
			Message: "consumerRef.name is required",
		}}
	}
	if len(name) > 63 {
		return []resources.FieldError{{
			Field:   "spec.consumerRef.name",
			Message: "consumerRef.name must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(name) {
		return []resources.FieldError{{
			Field:   "spec.consumerRef.name",
			Message: "consumerRef.name must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}

// validateBindingType checks that spec.bindingType is the Phase 1 allowed
// value "credentials".
func validateBindingType(bindingType string) []resources.FieldError {
	if bindingType == "" {
		return []resources.FieldError{{
			Field:   "spec.bindingType",
			Message: "bindingType is required",
		}}
	}
	if bindingType != resources.BindingTypeCredentials {
		return []resources.FieldError{{
			Field:   "spec.bindingType",
			Message: "bindingType must be credentials",
		}}
	}
	return nil
}
