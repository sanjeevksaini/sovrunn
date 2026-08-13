package validate

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/isocodes"
)

// alpha2SyntaxRe matches a two-letter uppercase ISO-3166-1 alpha-2 shape.
var alpha2SyntaxRe = regexp.MustCompile(`^[A-Z]{2}$`)

// adminAreaRe matches ISO-3166-2 shape: CC + hyphen + 1..3 alphanumerics
// (syntax only; no membership dataset — ADH-2026-048 decision 1).
var adminAreaRe = regexp.MustCompile(`^[A-Z]{2}-[A-Z0-9]{1,3}$`)

// ValidateAssignedAlpha2 rejects values that are not in the repository-owned
// assigned ISO-3166-1 alpha-2 dataset. Comparison is case-sensitive.
func ValidateAssignedAlpha2(field, code string) *apiproblem.Problem {
	if code == "" {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    apiproblem.ViolationOutOfRange,
			Message: "ISO-3166-1 alpha-2 code is required",
		}})
	}
	if !alpha2SyntaxRe.MatchString(code) {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    apiproblem.ViolationOutOfRange,
			Message: "ISO-3166-1 alpha-2 code must be two uppercase letters",
		}})
	}
	if !isocodes.IsAssignedAlpha2(code) {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    apiproblem.ViolationOutOfRange,
			Message: "ISO-3166-1 alpha-2 code is not an assigned value",
		}})
	}
	return nil
}

// ValidateOperatingMarkets validates a non-empty unique uppercase assigned
// ISO-3166-1 alpha-2 list (1..32) for CloudProvider.spec.operatingMarkets.
func ValidateOperatingMarkets(markets []string) *apiproblem.Problem {
	if len(markets) == 0 {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/operatingMarkets",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "operatingMarkets must be non-empty",
		}})
	}
	if len(markets) > MaxOperatingMarkets {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/operatingMarkets",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "operatingMarkets exceeds maximum length",
		}})
	}
	seen := make(map[string]struct{}, len(markets))
	for i, code := range markets {
		field := "/spec/operatingMarkets/" + itoa(i)
		if prob := ValidateAssignedAlpha2(field, code); prob != nil {
			return prob
		}
		if _, dup := seen[code]; dup {
			return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
				Field:   field,
				Code:    apiproblem.ViolationOutOfRange,
				Message: "operatingMarkets values must be unique",
			}})
		}
		seen[code] = struct{}{}
	}
	return nil
}

// ValidateAdministrativeAreaCode enforces syntax-plus-country-prefix only
// (ADH-2026-048 decision 1). Empty is permitted (optional field).
func ValidateAdministrativeAreaCode(countryCode, adminArea string) *apiproblem.Problem {
	if adminArea == "" {
		return nil
	}
	if len(adminArea) > MaxAdminAreaCodeLen {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/administrativeAreaCode",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "administrativeAreaCode exceeds maximum length",
		}})
	}
	if !adminAreaRe.MatchString(adminArea) {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/administrativeAreaCode",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "administrativeAreaCode must match ISO-3166-2 syntax",
		}})
	}
	if countryCode == "" || !strings.HasPrefix(adminArea, countryCode+"-") {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/administrativeAreaCode",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "administrativeAreaCode must use the countryCode prefix",
		}})
	}
	return nil
}

// IsUpperASCII reports whether s contains only ASCII uppercase letters A-Z.
func IsUpperASCII(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r > unicode.MaxASCII || !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b [12]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(b[pos:])
}
