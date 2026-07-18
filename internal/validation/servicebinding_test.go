package validation

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateServiceBinding_Valid(t *testing.T) {
	valid := []string{"a", "a1", "a-b", strings.Repeat("a", 63)}
	for _, name := range valid {
		errs := ValidateServiceBinding(validServiceBindingWithName(name))
		if len(errs) != 0 {
			t.Errorf("name %q: got errors %v, want none", name, errs)
		}
	}
}

func TestValidateServiceBinding_EmptyName(t *testing.T) {
	errs := ValidateServiceBinding(validServiceBindingWithName(""))
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateServiceBinding_InvalidNameFormat(t *testing.T) {
	invalid := []string{"ABC", "a b", "-abc", "abc-", "a.b", "_abc"}
	for _, name := range invalid {
		errs := ValidateServiceBinding(validServiceBindingWithName(name))
		if !hasFieldError(errs, "metadata.name") {
			t.Fatalf("name %q: got %v, want metadata.name error", name, errs)
		}
	}
}

func TestValidateServiceBinding_NameTooLong(t *testing.T) {
	errs := ValidateServiceBinding(validServiceBindingWithName(strings.Repeat("a", 64)))
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateServiceBinding_EmptyServiceInstanceRef(t *testing.T) {
	sb := validServiceBinding()
	sb.Spec.ServiceInstanceRef = ""
	errs := ValidateServiceBinding(sb)
	if !hasFieldError(errs, "spec.serviceInstanceRef") {
		t.Fatalf("got %v, want spec.serviceInstanceRef error", errs)
	}
}

func TestValidateServiceBinding_InvalidServiceInstanceRef(t *testing.T) {
	invalid := []string{"ABC", "a b", "-abc", "abc-", strings.Repeat("a", 64)}
	for _, ref := range invalid {
		sb := validServiceBinding()
		sb.Spec.ServiceInstanceRef = ref
		errs := ValidateServiceBinding(sb)
		if !hasFieldError(errs, "spec.serviceInstanceRef") {
			t.Fatalf("serviceInstanceRef %q: got %v, want spec.serviceInstanceRef error", ref, errs)
		}
	}
}

func TestValidateServiceBinding_NilConsumerRef(t *testing.T) {
	sb := validServiceBinding()
	sb.Spec.ConsumerRef = nil
	errs := ValidateServiceBinding(sb)
	if !hasFieldError(errs, "spec.consumerRef") {
		t.Fatalf("got %v, want spec.consumerRef error", errs)
	}
}

func TestValidateServiceBinding_EmptyConsumerRefKind(t *testing.T) {
	sb := validServiceBinding()
	sb.Spec.ConsumerRef.Kind = ""
	errs := ValidateServiceBinding(sb)
	if !hasFieldError(errs, "spec.consumerRef.kind") {
		t.Fatalf("got %v, want spec.consumerRef.kind error", errs)
	}
}

func TestValidateServiceBinding_InvalidConsumerRefName(t *testing.T) {
	invalid := []string{"", "ABC", "a b", "-abc", "abc-", strings.Repeat("a", 64)}
	for _, name := range invalid {
		sb := validServiceBinding()
		sb.Spec.ConsumerRef.Name = name
		errs := ValidateServiceBinding(sb)
		if !hasFieldError(errs, "spec.consumerRef.name") {
			t.Fatalf("consumerRef.name %q: got %v, want spec.consumerRef.name error", name, errs)
		}
	}
}

func TestValidateServiceBinding_BindingTypeCredentialsAccepted(t *testing.T) {
	sb := validServiceBinding()
	sb.Spec.BindingType = resources.BindingTypeCredentials
	errs := ValidateServiceBinding(sb)
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors for bindingType credentials", errs)
	}
}

func TestValidateServiceBinding_InvalidBindingType(t *testing.T) {
	invalid := []string{"", "endpoint", "config", "Credentials", "CREDENTIALS"}
	for _, bt := range invalid {
		sb := validServiceBinding()
		sb.Spec.BindingType = bt
		errs := ValidateServiceBinding(sb)
		if !hasFieldError(errs, "spec.bindingType") {
			t.Fatalf("bindingType %q: got %v, want spec.bindingType error", bt, errs)
		}
	}
}

func TestValidateServiceBinding_AnyNonEmptyConsumerRefKindAccepted(t *testing.T) {
	kinds := []string{"Application", "Job", "CronJob", "custom-kind"}
	for _, kind := range kinds {
		sb := validServiceBinding()
		sb.Spec.ConsumerRef.Kind = kind
		errs := ValidateServiceBinding(sb)
		if len(errs) != 0 {
			t.Fatalf("consumerRef.kind %q: got %v, want no errors", kind, errs)
		}
	}
}

func TestValidateServiceBindingPathSegment_Valid(t *testing.T) {
	errs := ValidateServiceBindingPathSegment("app-binding")
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func TestValidateServiceBindingPathSegment_Invalid(t *testing.T) {
	errs := ValidateServiceBindingPathSegment("Invalid Name")
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func validServiceBinding() resources.ServiceBinding {
	return validServiceBindingWithName("app-binding")
}

func validServiceBindingWithName(name string) resources.ServiceBinding {
	return resources.ServiceBinding{
		Metadata: resources.Metadata{Name: name},
		Spec: resources.ServiceBindingSpec{
			ServiceInstanceRef: "pg-prod",
			ConsumerRef: &resources.ConsumerRef{
				Kind: "Application",
				Name: "payments-api",
			},
			BindingType: resources.BindingTypeCredentials,
		},
	}
}
