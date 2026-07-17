package validation

import (
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Feature: serviceclass-serviceplan, Property 1: ValidateServiceClass accepts valid inputs
func TestProperty_ValidateServiceClass_ValidInputs(t *testing.T) {
	f := func(name string) bool {
		if !isValidDNSLabel(name) {
			return true
		}
		errs := ValidateServiceClass(resources.ServiceClass{
			Metadata: resources.Metadata{Name: name},
			Spec: resources.ServiceClassSpec{
				Category:  resources.CategoryDatabase,
				Lifecycle: resources.LifecycleActive,
			},
		})
		return len(errs) == 0
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceclass-serviceplan, Property 2: ValidateServiceClass rejects invalid names
func TestProperty_ValidateServiceClass_InvalidNames(t *testing.T) {
	f := func(name string) bool {
		if isValidDNSLabel(name) {
			return true
		}
		errs := ValidateServiceClass(resources.ServiceClass{
			Metadata: resources.Metadata{Name: name},
			Spec: resources.ServiceClassSpec{
				Category:  resources.CategoryDatabase,
				Lifecycle: resources.LifecycleActive,
			},
		})
		return hasFieldError(errs, "metadata.name")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
