package validate_test

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func TestValidateAssignedAlpha2_USINAcceptedZZRejected(t *testing.T) {
	t.Parallel()
	for _, code := range []string{"US", "IN"} {
		if prob := validate.ValidateAssignedAlpha2("/spec/countryCode", code); prob != nil {
			t.Fatalf("%s: %#v", code, prob)
		}
	}
	prob := validate.ValidateAssignedAlpha2("/spec/ownerRegistration/jurisdictionCode", "ZZ")
	if prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("ZZ: %#v", prob)
	}
}

func TestValidateAdministrativeAreaCode_Prefix(t *testing.T) {
	t.Parallel()
	if prob := validate.ValidateAdministrativeAreaCode("US", "US-CA"); prob != nil {
		t.Fatalf("valid: %#v", prob)
	}
	if prob := validate.ValidateAdministrativeAreaCode("US", "IN-MH"); prob == nil {
		t.Fatal("wrong prefix must fail")
	}
	if prob := validate.ValidateAdministrativeAreaCode("US", ""); prob != nil {
		t.Fatalf("optional empty: %#v", prob)
	}
}

func TestValidateOperatingMarkets(t *testing.T) {
	t.Parallel()
	if prob := validate.ValidateOperatingMarkets([]string{"US", "IN"}); prob != nil {
		t.Fatalf("valid: %#v", prob)
	}
	if prob := validate.ValidateOperatingMarkets([]string{"US", "US"}); prob == nil {
		t.Fatal("duplicates must fail")
	}
	if prob := validate.ValidateOperatingMarkets([]string{"ZZ"}); prob == nil {
		t.Fatal("unassigned must fail")
	}
}
