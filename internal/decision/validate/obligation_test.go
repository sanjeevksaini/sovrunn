package validate

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

func assertObligationViolation(t *testing.T, p *apiproblem.Problem, wantCode apiproblem.ViolationCode, wantField string) {
	t.Helper()
	if p == nil {
		t.Fatalf("expected Problem with %s at %s, got nil", wantCode, wantField)
	}
	if p.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("top-level code = %q, want %q", p.Code, apiproblem.CodeValidationFailed)
	}
	if p.Status != 422 {
		t.Fatalf("status = %d, want 422", p.Status)
	}
	if p.Type != "urn:sovrunn:problem:validation-failed" {
		t.Fatalf("type = %q, want urn:sovrunn:problem:validation-failed", p.Type)
	}
	if len(p.Violations) != 1 {
		t.Fatalf("violations len=%d, want 1: %#v", len(p.Violations), p.Violations)
	}
	v := p.Violations[0]
	if v.Code != wantCode {
		t.Fatalf("violations[0].code = %q, want %q", v.Code, wantCode)
	}
	if v.Field != wantField {
		t.Fatalf("violations[0].field = %q, want %q", v.Field, wantField)
	}
	if v.Message == "" {
		t.Fatal("violations[0].message must be non-empty (redactable)")
	}
}

func supportedObligations(ids ...string) map[string]bool {
	m := make(map[string]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m
}

func TestValidateObligations_HappyPathEmpty(t *testing.T) {
	t.Parallel()

	p := ValidateObligations(nil, nil)
	if p != nil {
		t.Fatalf("empty obligations must pass, got %#v", p)
	}
	p = ValidateObligations([]decision.Obligation{}, supportedObligations())
	if p != nil {
		t.Fatalf("empty obligations slice must pass, got %#v", p)
	}
}

func TestValidateObligations_HappyPathSupported(t *testing.T) {
	t.Parallel()

	obs := []decision.Obligation{
		{ID: "audit-retain", Mandatory: true, Payload: json.RawMessage(`{"retainDays":365}`)},
		{ID: "notify-owner", Mandatory: false},
	}
	p := ValidateObligations(obs, supportedObligations("audit-retain", "notify-owner"))
	if p != nil {
		t.Fatalf("supported obligations must pass, got %#v", p)
	}
}

func TestValidateObligations_HappyPathOptionalUnsupported(t *testing.T) {
	t.Parallel()

	// Optional obligations that are known but unsupported may be skipped.
	obs := []decision.Obligation{
		{ID: "notify-owner", Mandatory: false},
	}
	supported := map[string]bool{"notify-owner": false}
	p := ValidateObligations(obs, supported)
	if p != nil {
		t.Fatalf("optional unsupported obligation must pass, got %#v", p)
	}
}

func TestValidateObligations_Unknown(t *testing.T) {
	t.Parallel()

	obs := []decision.Obligation{
		{ID: "unknown-obligation", Mandatory: false},
	}
	p := ValidateObligations(obs, supportedObligations("audit-retain"))
	assertObligationViolation(t, p, CodeObligationUnknown, "/obligations/0/id")
}

func TestValidateObligations_UnknownMandatory(t *testing.T) {
	t.Parallel()

	// Unknown maps to DECISION_OBLIGATION_UNKNOWN even when mandatory
	// (F13-OBLIG-003 vocabulary vs mandatory-support distinction).
	obs := []decision.Obligation{
		{ID: "unknown-mandatory", Mandatory: true},
	}
	p := ValidateObligations(obs, supportedObligations("audit-retain"))
	assertObligationViolation(t, p, CodeObligationUnknown, "/obligations/0/id")
}

func TestValidateObligations_InvalidEmptyID(t *testing.T) {
	t.Parallel()

	obs := []decision.Obligation{
		{ID: "", Mandatory: true},
	}
	p := ValidateObligations(obs, supportedObligations("audit-retain"))
	assertObligationViolation(t, p, CodeObligationInvalid, "/obligations/0/id")
}

func TestValidateObligations_InvalidWhitespaceID(t *testing.T) {
	t.Parallel()

	obs := []decision.Obligation{
		{ID: "  audit-retain  ", Mandatory: true},
	}
	p := ValidateObligations(obs, supportedObligations("audit-retain"))
	assertObligationViolation(t, p, CodeObligationInvalid, "/obligations/0/id")
}

func TestValidateObligations_InvalidIDLength(t *testing.T) {
	t.Parallel()

	tooLong := strings.Repeat("a", obligationIDMaxLen+1)
	obs := []decision.Obligation{
		{ID: tooLong, Mandatory: false},
	}
	p := ValidateObligations(obs, map[string]bool{tooLong: true})
	assertObligationViolation(t, p, CodeObligationInvalid, "/obligations/0/id")
}

func TestValidateObligations_InvalidPayload(t *testing.T) {
	t.Parallel()

	obs := []decision.Obligation{
		{ID: "audit-retain", Mandatory: true, Payload: json.RawMessage(`{not-json}`)},
	}
	p := ValidateObligations(obs, supportedObligations("audit-retain"))
	assertObligationViolation(t, p, CodeObligationInvalid, "/obligations/0/payload")
}

func TestValidateObligations_UnsupportedMandatory(t *testing.T) {
	t.Parallel()

	// Edge 14 / F13-OBLIG-002: mandatory type known but not supported by the
	// enforcement point → DECISION_OBLIGATION_UNSUPPORTED_MANDATORY.
	obs := []decision.Obligation{
		{ID: "audit-retain", Mandatory: true},
	}
	supported := map[string]bool{"audit-retain": false}
	p := ValidateObligations(obs, supported)
	assertObligationViolation(t, p, CodeObligationUnsupportedMandatory, "/obligations/0/id")
}

func TestValidateObligations_ShortCircuitFirstIndex(t *testing.T) {
	t.Parallel()

	obs := []decision.Obligation{
		{ID: "audit-retain", Mandatory: true},
		{ID: "unknown-later", Mandatory: true},
	}
	p := ValidateObligations(obs, supportedObligations("audit-retain"))
	assertObligationViolation(t, p, CodeObligationUnknown, "/obligations/1/id")
}

func TestValidateObligations_InvalidBeforeUnknown(t *testing.T) {
	t.Parallel()

	// Structural invalidity is checked before vocabulary membership.
	obs := []decision.Obligation{
		{ID: "", Mandatory: true},
	}
	p := ValidateObligations(obs, nil)
	assertObligationViolation(t, p, CodeObligationInvalid, "/obligations/0/id")
}

func TestValidateObligations_NilSupportedWithPresentObligation(t *testing.T) {
	t.Parallel()

	obs := []decision.Obligation{
		{ID: "audit-retain", Mandatory: false},
	}
	p := ValidateObligations(obs, nil)
	assertObligationViolation(t, p, CodeObligationUnknown, "/obligations/0/id")
}
