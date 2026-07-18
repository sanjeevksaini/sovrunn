package validation

import (
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

// Feature: serviceinstance-servicebinding, Property 3: ValidateServiceBinding accepts valid DNS-label names
func TestProperty_ValidateServiceBinding_ValidNames(t *testing.T) {
	f := func(name, instanceRef, consumerName string) bool {
		if !isValidDNSLabel(name) ||
			!isValidDNSLabel(instanceRef) ||
			!isValidDNSLabel(consumerName) {
			return true
		}
		errs := ValidateServiceBinding(resources.ServiceBinding{
			Metadata: resources.Metadata{Name: name},
			Spec: resources.ServiceBindingSpec{
				ServiceInstanceRef: instanceRef,
				ConsumerRef: &resources.ConsumerRef{
					Kind: "Application",
					Name: consumerName,
				},
				BindingType: resources.BindingTypeCredentials,
			},
		})
		return len(errs) == 0
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceinstance-servicebinding, Property 4: ValidateServiceBinding rejects arbitrary invalid strings
func TestProperty_ValidateServiceBinding_InvalidNames(t *testing.T) {
	f := func(s string) bool {
		if isValidDNSLabel(s) {
			return true
		}

		sb := validServiceBinding()
		sb.Metadata.Name = s
		if !hasFieldError(ValidateServiceBinding(sb), "metadata.name") {
			return false
		}

		sb = validServiceBinding()
		sb.Spec.ServiceInstanceRef = s
		if !hasFieldError(ValidateServiceBinding(sb), "spec.serviceInstanceRef") {
			return false
		}

		sb = validServiceBinding()
		sb.Spec.ConsumerRef.Name = s
		return hasFieldError(ValidateServiceBinding(sb), "spec.consumerRef.name")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}

// Feature: serviceinstance-servicebinding, Property 5: valid bindingType credentials accepted; invalid rejected
func TestProperty_ValidateServiceBinding_BindingType(t *testing.T) {
	validTypes := []string{resources.BindingTypeCredentials}
	validSet := map[string]struct{}{resources.BindingTypeCredentials: {}}

	f := func(idx uint8, bindingType string) bool {
		sb := validServiceBinding()
		sb.Spec.BindingType = validTypes[int(idx)%len(validTypes)]
		if hasFieldError(ValidateServiceBinding(sb), "spec.bindingType") {
			return false
		}
		if _, ok := validSet[bindingType]; ok {
			return true
		}
		sb.Spec.BindingType = bindingType
		return hasFieldError(ValidateServiceBinding(sb), "spec.bindingType")
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
