package validation

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateServiceClass_Valid(t *testing.T) {
	errs := ValidateServiceClass(validServiceClass())
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func TestValidateServiceClass_EmptyName(t *testing.T) {
	sc := validServiceClass()
	sc.Metadata.Name = ""
	errs := ValidateServiceClass(sc)
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateServiceClass_InvalidName(t *testing.T) {
	invalid := []string{"ABC", "a b", "-abc", "abc-", "a.b", "_abc"}
	for _, name := range invalid {
		sc := validServiceClass()
		sc.Metadata.Name = name
		errs := ValidateServiceClass(sc)
		if !hasFieldError(errs, "metadata.name") {
			t.Fatalf("name %q: got %v, want metadata.name error", name, errs)
		}
	}
}

func TestValidateServiceClass_NameTooLong(t *testing.T) {
	sc := validServiceClass()
	sc.Metadata.Name = strings.Repeat("a", 64)
	errs := ValidateServiceClass(sc)
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateServiceClass_EmptyCategory(t *testing.T) {
	sc := validServiceClass()
	sc.Spec.Category = ""
	errs := ValidateServiceClass(sc)
	if !hasFieldError(errs, "spec.category") {
		t.Fatalf("got %v, want spec.category error", errs)
	}
}

func TestValidateServiceClass_InvalidCategory(t *testing.T) {
	sc := validServiceClass()
	sc.Spec.Category = "NotACategory"
	errs := ValidateServiceClass(sc)
	if !hasFieldError(errs, "spec.category") {
		t.Fatalf("got %v, want spec.category error", errs)
	}
}

func TestValidateServiceClass_EmptyLifecycle(t *testing.T) {
	sc := validServiceClass()
	sc.Spec.Lifecycle = ""
	errs := ValidateServiceClass(sc)
	if !hasFieldError(errs, "spec.lifecycle") {
		t.Fatalf("got %v, want spec.lifecycle error", errs)
	}
}

func TestValidateServiceClass_InvalidLifecycle(t *testing.T) {
	sc := validServiceClass()
	sc.Spec.Lifecycle = "NotALifecycle"
	errs := ValidateServiceClass(sc)
	if !hasFieldError(errs, "spec.lifecycle") {
		t.Fatalf("got %v, want spec.lifecycle error", errs)
	}
}

func TestValidateServiceClass_EmptyDefaultPlanNameAccepted(t *testing.T) {
	sc := validServiceClass()
	sc.Spec.DefaultPlanName = ""
	errs := ValidateServiceClass(sc)
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors for empty defaultPlanName", errs)
	}
}

func TestValidateServiceClass_ValidDefaultPlanNameAccepted(t *testing.T) {
	sc := validServiceClass()
	sc.Spec.DefaultPlanName = "small"
	errs := ValidateServiceClass(sc)
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors for valid defaultPlanName", errs)
	}
}

func TestValidateServiceClass_InvalidDefaultPlanName(t *testing.T) {
	sc := validServiceClass()
	sc.Spec.DefaultPlanName = "Invalid Plan"
	errs := ValidateServiceClass(sc)
	if !hasFieldError(errs, "spec.defaultPlanName") {
		t.Fatalf("got %v, want spec.defaultPlanName error", errs)
	}
}

func TestValidateServiceClassPathSegment_Valid(t *testing.T) {
	errs := ValidateServiceClassPathSegment("postgres")
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func TestValidateServiceClassPathSegment_Invalid(t *testing.T) {
	errs := ValidateServiceClassPathSegment("Invalid Name")
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func validServiceClass() resources.ServiceClass {
	return resources.ServiceClass{
		Metadata: resources.Metadata{Name: "postgres"},
		Spec: resources.ServiceClassSpec{
			Category:  resources.CategoryDatabase,
			Lifecycle: resources.LifecycleActive,
		},
	}
}
