package validation

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateCapability_Valid(t *testing.T) {
	errs := ValidateCapability(validCapability())
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func TestValidateCapability_EmptyName(t *testing.T) {
	c := validCapability()
	c.Metadata.Name = ""
	errs := ValidateCapability(c)
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateCapability_InvalidName(t *testing.T) {
	invalid := []string{"ABC", "a b", "-abc", "abc-", "a.b", "_abc"}
	for _, name := range invalid {
		c := validCapability()
		c.Metadata.Name = name
		errs := ValidateCapability(c)
		if !hasFieldError(errs, "metadata.name") {
			t.Fatalf("name %q: got %v, want metadata.name error", name, errs)
		}
	}
}

func TestValidateCapability_NameTooLong(t *testing.T) {
	c := validCapability()
	c.Metadata.Name = strings.Repeat("a", 64)
	errs := ValidateCapability(c)
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateCapability_EmptyPluginRef(t *testing.T) {
	c := validCapability()
	c.Spec.PluginRef = ""
	errs := ValidateCapability(c)
	if !hasFieldError(errs, "spec.pluginRef") {
		t.Fatalf("got %v, want spec.pluginRef error", errs)
	}
}

func TestValidateCapability_InvalidPluginRef(t *testing.T) {
	invalid := []string{"ABC", "a b", "-abc", "abc-", strings.Repeat("a", 64)}
	for _, ref := range invalid {
		c := validCapability()
		c.Spec.PluginRef = ref
		errs := ValidateCapability(c)
		if !hasFieldError(errs, "spec.pluginRef") {
			t.Fatalf("pluginRef %q: got %v, want spec.pluginRef error", ref, errs)
		}
	}
}

func TestValidateCapability_EmptyServiceClassRef(t *testing.T) {
	c := validCapability()
	c.Spec.ServiceClassRef = ""
	errs := ValidateCapability(c)
	if !hasFieldError(errs, "spec.serviceClassRef") {
		t.Fatalf("got %v, want spec.serviceClassRef error", errs)
	}
}

func TestValidateCapability_InvalidServiceClassRef(t *testing.T) {
	invalid := []string{"ABC", "a b", "-abc", "abc-", strings.Repeat("a", 64)}
	for _, ref := range invalid {
		c := validCapability()
		c.Spec.ServiceClassRef = ref
		errs := ValidateCapability(c)
		if !hasFieldError(errs, "spec.serviceClassRef") {
			t.Fatalf("serviceClassRef %q: got %v, want spec.serviceClassRef error", ref, errs)
		}
	}
}

func TestValidateCapability_EmptyOperation(t *testing.T) {
	c := validCapability()
	c.Spec.Operation = ""
	errs := ValidateCapability(c)
	if !hasFieldError(errs, "spec.operation") {
		t.Fatalf("got %v, want spec.operation error", errs)
	}
}

func TestValidateCapability_InvalidOperation(t *testing.T) {
	c := validCapability()
	c.Spec.Operation = "NotAnOperation"
	errs := ValidateCapability(c)
	if !hasFieldError(errs, "spec.operation") {
		t.Fatalf("got %v, want spec.operation error", errs)
	}
}

func TestValidateCapabilityPathSegment_Valid(t *testing.T) {
	errs := ValidateCapabilityPathSegment("postgres-provision")
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func TestValidateCapabilityPathSegment_Invalid(t *testing.T) {
	errs := ValidateCapabilityPathSegment("Invalid Name")
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func validCapability() resources.Capability {
	return resources.Capability{
		Metadata: resources.Metadata{Name: "postgres-provision"},
		Spec: resources.CapabilitySpec{
			PluginRef:       "postgres-plugin",
			ServiceClassRef: "postgres",
			Operation:       resources.CapOpProvision,
		},
	}
}
