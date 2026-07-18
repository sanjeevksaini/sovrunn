package validation

import "github.com/sanjeevksaini/sovrunn/internal/resources"

// ValidateCapability is a pure function. It validates all user-authored
// Capability fields without performing I/O or registry lookups. It returns all
// FieldErrors found in a single call and returns nil if the resource is valid.
// spec.supported defaults to false (Go zero-value) and is not validated.
// spec.description is optional and not format-validated.
func ValidateCapability(c resources.Capability) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateName(c.Metadata.Name)...)
	errs = append(errs, validatePluginRef(c.Spec.PluginRef)...)
	errs = append(errs, validateServiceClassRef(c.Spec.ServiceClassRef)...)
	errs = append(errs, validateCapabilityOperation(c.Spec.Operation)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// ValidateCapabilityPathSegment validates the single Capability URL path
// segment before a registry lookup. It only checks the name segment and maps
// it to metadata.name.
func ValidateCapabilityPathSegment(name string) []resources.FieldError {
	return validateName(name)
}

// validatePluginRef validates spec.pluginRef using the same DNS-label rules as
// metadata.name. Existence of the referenced Plugin is NOT checked here.
func validatePluginRef(pluginRef string) []resources.FieldError {
	if pluginRef == "" {
		return []resources.FieldError{{
			Field:   "spec.pluginRef",
			Message: "pluginRef is required",
		}}
	}
	if len(pluginRef) > 63 {
		return []resources.FieldError{{
			Field:   "spec.pluginRef",
			Message: "pluginRef must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(pluginRef) {
		return []resources.FieldError{{
			Field:   "spec.pluginRef",
			Message: "pluginRef must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}

// validateServiceClassRef validates spec.serviceClassRef using the same
// DNS-label rules as metadata.name. Existence of the referenced ServiceClass
// is NOT checked here.
func validateServiceClassRef(serviceClassRef string) []resources.FieldError {
	if serviceClassRef == "" {
		return []resources.FieldError{{
			Field:   "spec.serviceClassRef",
			Message: "serviceClassRef is required",
		}}
	}
	if len(serviceClassRef) > 63 {
		return []resources.FieldError{{
			Field:   "spec.serviceClassRef",
			Message: "serviceClassRef must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(serviceClassRef) {
		return []resources.FieldError{{
			Field:   "spec.serviceClassRef",
			Message: "serviceClassRef must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}

// validateCapabilityOperation checks that spec.operation is present and a
// known CapabilityOperation value.
func validateCapabilityOperation(operation string) []resources.FieldError {
	if operation == "" {
		return []resources.FieldError{{
			Field:   "spec.operation",
			Message: "operation is required",
		}}
	}
	switch operation {
	case resources.CapOpValidate,
		resources.CapOpPlan,
		resources.CapOpProvision,
		resources.CapOpConfigure,
		resources.CapOpBind,
		resources.CapOpObserve,
		resources.CapOpScale,
		resources.CapOpUpgrade,
		resources.CapOpBackup,
		resources.CapOpRestore,
		resources.CapOpRotateCredentials,
		resources.CapOpUnbind,
		resources.CapOpDelete:
		return nil
	}
	return []resources.FieldError{{
		Field:   "spec.operation",
		Message: "operation must be one of: Validate, Plan, Provision, Configure, Bind, Observe, Scale, Upgrade, Backup, Restore, RotateCredentials, Unbind, Delete",
	}}
}
