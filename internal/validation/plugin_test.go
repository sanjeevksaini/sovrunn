package validation

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidatePlugin_Valid(t *testing.T) {
	errs := ValidatePlugin(validPlugin())
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func TestValidatePlugin_EmptyName(t *testing.T) {
	p := validPlugin()
	p.Metadata.Name = ""
	errs := ValidatePlugin(p)
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidatePlugin_InvalidName(t *testing.T) {
	invalid := []string{"ABC", "a b", "-abc", "abc-", "a.b", "_abc"}
	for _, name := range invalid {
		p := validPlugin()
		p.Metadata.Name = name
		errs := ValidatePlugin(p)
		if !hasFieldError(errs, "metadata.name") {
			t.Fatalf("name %q: got %v, want metadata.name error", name, errs)
		}
	}
}

func TestValidatePlugin_NameTooLong(t *testing.T) {
	p := validPlugin()
	p.Metadata.Name = strings.Repeat("a", 64)
	errs := ValidatePlugin(p)
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidatePlugin_EmptyPluginType(t *testing.T) {
	p := validPlugin()
	p.Spec.PluginType = ""
	errs := ValidatePlugin(p)
	if !hasFieldError(errs, "spec.pluginType") {
		t.Fatalf("got %v, want spec.pluginType error", errs)
	}
}

func TestValidatePlugin_InvalidPluginType(t *testing.T) {
	p := validPlugin()
	p.Spec.PluginType = "NotAPluginType"
	errs := ValidatePlugin(p)
	if !hasFieldError(errs, "spec.pluginType") {
		t.Fatalf("got %v, want spec.pluginType error", errs)
	}
}

func TestValidatePlugin_EmptyVersion(t *testing.T) {
	p := validPlugin()
	p.Spec.Version = ""
	errs := ValidatePlugin(p)
	if !hasFieldError(errs, "spec.version") {
		t.Fatalf("got %v, want spec.version error", errs)
	}
}

func TestValidatePlugin_NilServiceClassRefs(t *testing.T) {
	p := validPlugin()
	p.Spec.ServiceClassRefs = nil
	errs := ValidatePlugin(p)
	if !hasFieldError(errs, "spec.serviceClassRefs") {
		t.Fatalf("got %v, want spec.serviceClassRefs error", errs)
	}
}

func TestValidatePlugin_EmptyServiceClassRefs(t *testing.T) {
	p := validPlugin()
	p.Spec.ServiceClassRefs = []string{}
	errs := ValidatePlugin(p)
	if !hasFieldError(errs, "spec.serviceClassRefs") {
		t.Fatalf("got %v, want spec.serviceClassRefs error", errs)
	}
}

func TestValidatePlugin_InvalidServiceClassRefsEntries(t *testing.T) {
	invalid := []string{"", "ABC", "a b", "-abc", "abc-", strings.Repeat("a", 64)}
	for _, entry := range invalid {
		p := validPlugin()
		p.Spec.ServiceClassRefs = []string{entry}
		errs := ValidatePlugin(p)
		if !hasFieldError(errs, "spec.serviceClassRefs") {
			t.Fatalf("entry %q: got %v, want spec.serviceClassRefs error", entry, errs)
		}
	}
}

func TestValidatePlugin_EmptyDeploymentMode(t *testing.T) {
	p := validPlugin()
	p.Spec.DeploymentMode = ""
	errs := ValidatePlugin(p)
	if !hasFieldError(errs, "spec.deploymentMode") {
		t.Fatalf("got %v, want spec.deploymentMode error", errs)
	}
}

func TestValidatePlugin_InvalidDeploymentMode(t *testing.T) {
	p := validPlugin()
	p.Spec.DeploymentMode = "sidecar"
	errs := ValidatePlugin(p)
	if !hasFieldError(errs, "spec.deploymentMode") {
		t.Fatalf("got %v, want spec.deploymentMode error", errs)
	}
}

func TestValidatePluginPathSegment_Valid(t *testing.T) {
	errs := ValidatePluginPathSegment("postgres-plugin")
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func TestValidatePluginPathSegment_Invalid(t *testing.T) {
	errs := ValidatePluginPathSegment("Invalid Name")
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func validPlugin() resources.Plugin {
	return resources.Plugin{
		Metadata: resources.Metadata{Name: "postgres-plugin"},
		Spec: resources.PluginSpec{
			PluginType:       resources.PluginTypeDStoreOps,
			Version:          "1.0.0",
			ServiceClassRefs: []string{"postgres"},
			DeploymentMode:   resources.DeploymentModeCompiledIn,
		},
	}
}
