package validation

import (
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Feature: plugin-capability-registry, Property 5: valid DNS-label names with valid enum operation accepted
func TestProperty_ValidateCapability_ValidInputs(t *testing.T) {
	f := func(name string) bool {
		if !isValidDNSLabel(name) {
			return true
		}
		c := validCapability()
		c.Metadata.Name = name
		errs := ValidateCapability(c)
		return len(errs) == 0
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: plugin-capability-registry, Property 6: arbitrary invalid strings rejected for Capability name/pluginRef/serviceClassRef
func TestProperty_ValidateCapability_InvalidRefs(t *testing.T) {
	f := func(s string) bool {
		if isValidDNSLabel(s) {
			return true
		}
		c := validCapability()
		c.Metadata.Name = s
		if !hasFieldError(ValidateCapability(c), "metadata.name") {
			return false
		}
		c = validCapability()
		c.Spec.PluginRef = s
		if !hasFieldError(ValidateCapability(c), "spec.pluginRef") {
			return false
		}
		c = validCapability()
		c.Spec.ServiceClassRef = s
		return hasFieldError(ValidateCapability(c), "spec.serviceClassRef")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: plugin-capability-registry, Property 7: valid operation values accepted; invalid values rejected
func TestProperty_ValidateCapability_Operation(t *testing.T) {
	validOps := []string{
		resources.CapOpValidate,
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
		resources.CapOpDelete,
	}
	validSet := make(map[string]struct{}, len(validOps))
	for _, op := range validOps {
		validSet[op] = struct{}{}
	}

	f := func(idx uint8, operation string) bool {
		c := validCapability()
		c.Spec.Operation = validOps[int(idx)%len(validOps)]
		if hasFieldError(ValidateCapability(c), "spec.operation") {
			return false
		}
		if _, ok := validSet[operation]; ok {
			return true
		}
		c.Spec.Operation = operation
		return hasFieldError(ValidateCapability(c), "spec.operation")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
