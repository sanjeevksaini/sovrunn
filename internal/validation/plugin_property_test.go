package validation

import (
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Feature: plugin-capability-registry, Property 1: valid DNS-label names with valid enum values accepted
func TestProperty_ValidatePlugin_ValidInputs(t *testing.T) {
	f := func(name string) bool {
		if !isValidDNSLabel(name) {
			return true
		}
		p := validPlugin()
		p.Metadata.Name = name
		errs := ValidatePlugin(p)
		return len(errs) == 0
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: plugin-capability-registry, Property 2: arbitrary invalid strings rejected for Plugin name
func TestProperty_ValidatePlugin_InvalidNames(t *testing.T) {
	f := func(name string) bool {
		if isValidDNSLabel(name) {
			return true
		}
		p := validPlugin()
		p.Metadata.Name = name
		errs := ValidatePlugin(p)
		return hasFieldError(errs, "metadata.name")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: plugin-capability-registry, Property 3: valid pluginType values accepted; invalid values rejected
func TestProperty_ValidatePlugin_PluginType(t *testing.T) {
	validTypes := []string{
		resources.PluginTypeDStoreOps,
		resources.PluginTypeCacheOps,
		resources.PluginTypeStreamOps,
		resources.PluginTypeObjectOps,
		resources.PluginTypeGatewayOps,
		resources.PluginTypeFaasOps,
		resources.PluginTypeLBOps,
		resources.PluginTypeK8sOps,
		resources.PluginTypeBigDataOps,
		resources.PluginTypeSdeOps,
	}
	validSet := make(map[string]struct{}, len(validTypes))
	for _, pt := range validTypes {
		validSet[pt] = struct{}{}
	}

	f := func(idx uint8, pluginType string) bool {
		p := validPlugin()
		p.Spec.PluginType = validTypes[int(idx)%len(validTypes)]
		if hasFieldError(ValidatePlugin(p), "spec.pluginType") {
			return false
		}
		if _, ok := validSet[pluginType]; ok {
			return true
		}
		p.Spec.PluginType = pluginType
		return hasFieldError(ValidatePlugin(p), "spec.pluginType")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: plugin-capability-registry, Property 4: valid deploymentMode values accepted; invalid values rejected
func TestProperty_ValidatePlugin_DeploymentMode(t *testing.T) {
	validModes := []string{resources.DeploymentModeCompiledIn}
	validSet := map[string]struct{}{resources.DeploymentModeCompiledIn: {}}

	f := func(idx uint8, mode string) bool {
		p := validPlugin()
		p.Spec.DeploymentMode = validModes[int(idx)%len(validModes)]
		if hasFieldError(ValidatePlugin(p), "spec.deploymentMode") {
			return false
		}
		if _, ok := validSet[mode]; ok {
			return true
		}
		p.Spec.DeploymentMode = mode
		return hasFieldError(ValidatePlugin(p), "spec.deploymentMode")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
