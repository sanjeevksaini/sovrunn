package validation

import (
	"testing"
	"testing/quick"
)

// Feature: serviceclass-serviceplan, Property 3: ValidateServicePlan rejects forbidden parameter keys
func TestProperty_ValidateServicePlan_ForbiddenParameterKeys(t *testing.T) {
	forbidden := []string{
		"password", "secret", "token", "credential", "auth",
		"apikey", "accesskey", "secretkey", "privatekey",
	}
	f := func(prefix, suffix uint8) bool {
		p := string(rune('a' + (prefix % 26)))
		s := string(rune('a' + (suffix % 26)))
		for _, stem := range forbidden {
			key := p + stem + s
			sp := validServicePlan()
			sp.Spec.Parameters = map[string]string{key: "value"}
			errs := ValidateServicePlan(sp)
			if !hasFieldError(errs, "spec.parameters") {
				return false
			}
		}
		return true
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Fatal(err)
	}
}
