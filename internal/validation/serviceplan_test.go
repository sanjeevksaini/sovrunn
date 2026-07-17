package validation

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/resources"
)

func TestValidateServicePlan_Valid(t *testing.T) {
	errs := ValidateServicePlan(validServicePlan())
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func TestValidateServicePlan_EmptyName(t *testing.T) {
	sp := validServicePlan()
	sp.Metadata.Name = ""
	errs := ValidateServicePlan(sp)
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateServicePlan_InvalidName(t *testing.T) {
	invalid := []string{"ABC", "a b", "-abc", "abc-", "a.b", "_abc"}
	for _, name := range invalid {
		sp := validServicePlan()
		sp.Metadata.Name = name
		errs := ValidateServicePlan(sp)
		if !hasFieldError(errs, "metadata.name") {
			t.Fatalf("name %q: got %v, want metadata.name error", name, errs)
		}
	}
}

func TestValidateServicePlan_NameTooLong(t *testing.T) {
	sp := validServicePlan()
	sp.Metadata.Name = strings.Repeat("a", 64)
	errs := ValidateServicePlan(sp)
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateServicePlan_EmptyServiceClassName(t *testing.T) {
	sp := validServicePlan()
	sp.Spec.ServiceClassName = ""
	errs := ValidateServicePlan(sp)
	if !hasFieldError(errs, "spec.serviceClassName") {
		t.Fatalf("got %v, want spec.serviceClassName error", errs)
	}
}

func TestValidateServicePlan_InvalidServiceClassName(t *testing.T) {
	sp := validServicePlan()
	sp.Spec.ServiceClassName = "Invalid Class"
	errs := ValidateServicePlan(sp)
	if !hasFieldError(errs, "spec.serviceClassName") {
		t.Fatalf("got %v, want spec.serviceClassName error", errs)
	}
}

func TestValidateServicePlan_EmptyTier(t *testing.T) {
	sp := validServicePlan()
	sp.Spec.Tier = ""
	errs := ValidateServicePlan(sp)
	if !hasFieldError(errs, "spec.tier") {
		t.Fatalf("got %v, want spec.tier error", errs)
	}
}

func TestValidateServicePlan_InvalidTier(t *testing.T) {
	sp := validServicePlan()
	sp.Spec.Tier = "NotATier"
	errs := ValidateServicePlan(sp)
	if !hasFieldError(errs, "spec.tier") {
		t.Fatalf("got %v, want spec.tier error", errs)
	}
}

func TestValidateServicePlan_EmptyLifecycle(t *testing.T) {
	sp := validServicePlan()
	sp.Spec.Lifecycle = ""
	errs := ValidateServicePlan(sp)
	if !hasFieldError(errs, "spec.lifecycle") {
		t.Fatalf("got %v, want spec.lifecycle error", errs)
	}
}

func TestValidateServicePlan_InvalidLifecycle(t *testing.T) {
	sp := validServicePlan()
	sp.Spec.Lifecycle = "NotALifecycle"
	errs := ValidateServicePlan(sp)
	if !hasFieldError(errs, "spec.lifecycle") {
		t.Fatalf("got %v, want spec.lifecycle error", errs)
	}
}

func TestValidateServicePlan_ForbiddenParameterKeys(t *testing.T) {
	forbidden := []string{
		"apiKey",
		"ACCESSKEY",
		"secretKey",
		"privateKey",
		"password",
		"token",
		"auth",
		"credential",
		"secret",
	}
	for _, key := range forbidden {
		sp := validServicePlan()
		sp.Spec.Parameters = map[string]string{key: "value"}
		errs := ValidateServicePlan(sp)
		if !hasFieldError(errs, "spec.parameters") {
			t.Fatalf("key %q: got %v, want spec.parameters error", key, errs)
		}
	}
}

func TestValidateServicePlan_BenignParameterKeysAccepted(t *testing.T) {
	benign := []string{
		"key",
		"regionKey",
		"masterKeyCount",
		"partitionKey",
		"sortKey",
	}
	for _, key := range benign {
		sp := validServicePlan()
		sp.Spec.Parameters = map[string]string{key: "value"}
		errs := ValidateServicePlan(sp)
		if len(errs) != 0 {
			t.Fatalf("key %q: got %v, want no errors", key, errs)
		}
	}
}

func TestValidateServicePlanPathSegments_Valid(t *testing.T) {
	errs := ValidateServicePlanPathSegments("postgres", "small")
	if len(errs) != 0 {
		t.Fatalf("got %v, want no errors", errs)
	}
}

func TestValidateServicePlanPathSegments_InvalidServiceClassName(t *testing.T) {
	errs := ValidateServicePlanPathSegments("Invalid Class", "small")
	if !hasFieldError(errs, "spec.serviceClassName") {
		t.Fatalf("got %v, want spec.serviceClassName error", errs)
	}
}

func TestValidateServicePlanPathSegments_InvalidName(t *testing.T) {
	errs := ValidateServicePlanPathSegments("postgres", "Invalid Name")
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func TestValidateServicePlanPathSegments_BothInvalid(t *testing.T) {
	errs := ValidateServicePlanPathSegments("Invalid Class", "Invalid Name")
	if !hasFieldError(errs, "spec.serviceClassName") {
		t.Fatalf("got %v, want spec.serviceClassName error", errs)
	}
	if !hasFieldError(errs, "metadata.name") {
		t.Fatalf("got %v, want metadata.name error", errs)
	}
}

func validServicePlan() resources.ServicePlan {
	return resources.ServicePlan{
		Metadata: resources.Metadata{Name: "small"},
		Spec: resources.ServicePlanSpec{
			ServiceClassName: "postgres",
			Tier:             resources.TierSmall,
			Lifecycle:        resources.LifecycleActive,
		},
	}
}
