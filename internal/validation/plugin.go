package validation

import "github.com/sanjeevksaini/sovrunn/internal/resources"

// ValidatePlugin is a pure function. It validates all user-authored Plugin
// fields without performing I/O or registry lookups. It returns all
// FieldErrors found in a single call and returns nil if the resource is valid.
func ValidatePlugin(p resources.Plugin) []resources.FieldError {
	var errs []resources.FieldError
	errs = append(errs, validateName(p.Metadata.Name)...)
	errs = append(errs, validatePluginType(p.Spec.PluginType)...)
	errs = append(errs, validatePluginVersion(p.Spec.Version)...)
	errs = append(errs, validateServiceClassRefs(p.Spec.ServiceClassRefs)...)
	errs = append(errs, validateDeploymentMode(p.Spec.DeploymentMode)...)
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// ValidatePluginPathSegment validates the single Plugin URL path segment
// before a registry lookup. It only checks the name segment and maps it to
// metadata.name.
func ValidatePluginPathSegment(name string) []resources.FieldError {
	return validateName(name)
}

// validatePluginType checks that spec.pluginType is present and a known value.
func validatePluginType(pluginType string) []resources.FieldError {
	if pluginType == "" {
		return []resources.FieldError{{
			Field:   "spec.pluginType",
			Message: "pluginType is required",
		}}
	}
	switch pluginType {
	case resources.PluginTypeDStoreOps,
		resources.PluginTypeCacheOps,
		resources.PluginTypeStreamOps,
		resources.PluginTypeObjectOps,
		resources.PluginTypeGatewayOps,
		resources.PluginTypeFaasOps,
		resources.PluginTypeLBOps,
		resources.PluginTypeK8sOps,
		resources.PluginTypeBigDataOps,
		resources.PluginTypeSdeOps:
		return nil
	}
	return []resources.FieldError{{
		Field:   "spec.pluginType",
		Message: "pluginType must be one of: dStoreOps, cacheOps, streamOps, objectOps, gatewayOps, faasOps, lbOps, k8sOps, bigDataOps, sdeOps",
	}}
}

// validatePluginVersion checks that spec.version is present and non-empty.
func validatePluginVersion(version string) []resources.FieldError {
	if version == "" {
		return []resources.FieldError{{
			Field:   "spec.version",
			Message: "version is required",
		}}
	}
	return nil
}

// validateServiceClassRefs checks that spec.serviceClassRefs is non-nil with
// at least one entry, and that each entry is a valid DNS-label (1–63).
// Existence of referenced ServiceClasses is NOT checked here.
func validateServiceClassRefs(refs []string) []resources.FieldError {
	if len(refs) == 0 {
		return []resources.FieldError{{
			Field:   "spec.serviceClassRefs",
			Message: "serviceClassRefs is required and must contain at least one entry",
		}}
	}
	var errs []resources.FieldError
	for _, ref := range refs {
		errs = append(errs, validateServiceClassRefEntry(ref)...)
	}
	return errs
}

// validateServiceClassRefEntry validates one serviceClassRefs entry using the
// same DNS-label rules as metadata.name, with field mapped to
// spec.serviceClassRefs.
func validateServiceClassRefEntry(ref string) []resources.FieldError {
	if ref == "" {
		return []resources.FieldError{{
			Field:   "spec.serviceClassRefs",
			Message: "serviceClassRefs entries must be non-empty",
		}}
	}
	if len(ref) > 63 {
		return []resources.FieldError{{
			Field:   "spec.serviceClassRefs",
			Message: "serviceClassRefs entries must not exceed 63 characters",
		}}
	}
	if !dnsLabelRe.MatchString(ref) {
		return []resources.FieldError{{
			Field:   "spec.serviceClassRefs",
			Message: "serviceClassRefs entries must be a valid DNS label: lowercase alphanumeric and hyphens, no leading/trailing hyphens",
		}}
	}
	return nil
}

// validateDeploymentMode checks that spec.deploymentMode is present and a
// known Phase 1 value (compiled-in only).
func validateDeploymentMode(mode string) []resources.FieldError {
	if mode == "" {
		return []resources.FieldError{{
			Field:   "spec.deploymentMode",
			Message: "deploymentMode is required",
		}}
	}
	switch mode {
	case resources.DeploymentModeCompiledIn:
		return nil
	}
	return []resources.FieldError{{
		Field:   "spec.deploymentMode",
		Message: "deploymentMode must be one of: compiled-in",
	}}
}
