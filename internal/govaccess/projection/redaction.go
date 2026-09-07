package projection

import "strings"

// SafeDenial is the customer-safe denial projection shape.
type SafeDenial struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// SafeDenialProjection returns a fixed disclosure-safe denial.
func SafeDenialProjection() SafeDenial {
	return SafeDenial{
		Code:    "RESOURCE_NOT_FOUND",
		Message: "requested resource is unavailable",
	}
}

var sensitiveKeyFragments = []string{
	"token",
	"secret",
	"credential",
	"password",
	"privatekey",
	"policydigest",
	"provenance",
	"correlationid",
	"issuer",
	"subject",
	"eligibility",
	"controller",
	"resourceversion",
	"audit",
}

// Redact removes sensitive and internal keys from a generic field map.
func Redact(input map[string]string) map[string]string {
	if len(input) == 0 {
		return map[string]string{}
	}
	out := make(map[string]string, len(input))
	for k, v := range input {
		if isSensitiveKey(k) {
			continue
		}
		out[k] = v
	}
	return out
}

func isSensitiveKey(key string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(key, "_", ""))
	normalized = strings.ReplaceAll(normalized, "-", "")
	if normalized == "kind" || normalized == "apiversion" || normalized == "status" {
		return false
	}
	for _, fragment := range sensitiveKeyFragments {
		if strings.Contains(normalized, fragment) {
			return true
		}
	}
	return normalized == "internalstate"
}
