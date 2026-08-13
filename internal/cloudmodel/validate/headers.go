package validate

import (
	"encoding/json"
	"strings"
	"unicode"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// ValidateIdempotencyKey checks a required Idempotency-Key header value
// (ADH-2026-048; VS-000 limits.idempotencyKey=128).
// Missing, empty, malformed, or over-length values return MALFORMED_REQUEST.
func ValidateIdempotencyKey(value string) *apiproblem.Problem {
	if value == "" {
		return apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("Idempotency-Key is required")
	}
	if strings.TrimSpace(value) == "" {
		return apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("Idempotency-Key must not be empty")
	}
	if len(value) > MaxIdempotencyKeyLen {
		return apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("Idempotency-Key exceeds maximum length")
	}
	for _, r := range value {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) || unicode.IsSpace(r) {
			return apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("Idempotency-Key is malformed")
		}
	}
	return nil
}

// ValidateParticipationCreateHeaders rejects If-Match on participation
// collection create and requires a valid Idempotency-Key.
func ValidateParticipationCreateHeaders(idempotencyKey, ifMatch string) *apiproblem.Problem {
	if ifMatch != "" {
		return apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("If-Match is not accepted on participation create")
	}
	return ValidateIdempotencyKey(idempotencyKey)
}

// ValidateRequiredIfMatch requires a syntactically non-empty If-Match token.
// Missing or malformed values map to STALE_RESOURCE_VERSION (ADH-2026-046/056).
// Currency against the stored resourceVersion is checked separately.
func ValidateRequiredIfMatch(ifMatch string) *apiproblem.Problem {
	if ifMatch == "" || strings.TrimSpace(ifMatch) == "" {
		return apiproblem.New(apiproblem.CodeStaleResourceVersion).WithDetail("If-Match is required")
	}
	for _, r := range ifMatch {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return apiproblem.New(apiproblem.CodeStaleResourceVersion).WithDetail("If-Match is malformed")
		}
	}
	return nil
}

// CheckIfMatchCurrent compares a validated If-Match token to the current
// opaque resourceVersion using the inherited FEATURE-0012 helper.
func CheckIfMatchCurrent(ifMatch, currentResourceVersion string) *apiproblem.Problem {
	if prob := ValidateRequiredIfMatch(ifMatch); prob != nil {
		return prob
	}
	return apivalid.CheckIfMatch(ifMatch, currentResourceVersion)
}

// ValidateEmptyActionBody enforces the participation item-action empty-body
// contract (ADH-2026-047 decision 4). EOF/zero bytes are accepted. Whitespace,
// `{}`, or any other non-empty body is MALFORMED_REQUEST, except a valid
// duplicate-free top-level bootstrapGrant object which is reported as
// AUTHORIZATION_DENIED for the forged-grant audited-denial path.
func ValidateEmptyActionBody(body []byte) *apiproblem.Problem {
	if len(body) == 0 {
		return nil
	}
	trimmed := strings.TrimSpace(string(body))
	if trimmed == "" {
		return apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("participation action body must be empty")
	}
	var root map[string]json.RawMessage
	if err := json.Unmarshal([]byte(trimmed), &root); err != nil {
		return apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("participation action body must be empty")
	}
	if len(root) == 1 {
		if raw, ok := root["bootstrapGrant"]; ok && len(raw) > 0 {
			return apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail("forged bootstrapGrant is denied")
		}
	}
	return apiproblem.New(apiproblem.CodeMalformedRequest).WithDetail("participation action body must be empty")
}
