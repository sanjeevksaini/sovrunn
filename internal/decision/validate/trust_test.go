package validate

import (
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

func assertTrustViolation(t *testing.T, p *apiproblem.Problem, wantCode apiproblem.ViolationCode, wantField string) {
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

func validTrustCarrier() *decision.TrustCarrier {
	return &decision.TrustCarrier{
		State:                  "trusted-root-a",
		CanonicalizationMethod: "canonicalization:profile-declared",
		DigestAlgorithm:        "digest:profile-declared",
		CoveredFields:          "covered:decision-body",
		SignatureAlgorithm:     "signature:profile-declared",
	}
}

func TestValidateTrust_HappyPathOptionalAbsent(t *testing.T) {
	t.Parallel()

	p := ValidateTrust(TrustInput{Carrier: nil, Required: false})
	if p != nil {
		t.Fatalf("optional absent trust must pass, got %#v", p)
	}
}

func TestValidateTrust_HappyPathPresent(t *testing.T) {
	t.Parallel()

	p := ValidateTrust(TrustInput{
		Carrier:       validTrustCarrier(),
		Required:      true,
		ExpectedState: "trusted-root-a",
	})
	if p != nil {
		t.Fatalf("valid required trust must pass, got %#v", p)
	}
}

func TestValidateTrust_RequiredAbsent(t *testing.T) {
	t.Parallel()

	p := ValidateTrust(TrustInput{Carrier: nil, Required: true})
	assertTrustViolation(t, p, CodeTrustRequired, ptrTrust)
}

func TestValidateTrust_PresentEmpty(t *testing.T) {
	t.Parallel()

	// RID-08: present-but-structurally-empty fails closed for optional and required.
	for _, required := range []bool{false, true} {
		required := required
		t.Run(funcName(required), func(t *testing.T) {
			t.Parallel()
			p := ValidateTrust(TrustInput{
				Carrier:  &decision.TrustCarrier{},
				Required: required,
			})
			assertTrustViolation(t, p, CodeTrustUnknown, ptrTrust)
		})
	}
}

func funcName(required bool) string {
	if required {
		return "required"
	}
	return "optional"
}

func TestValidateTrust_OpaqueExpectedObservedMismatch(t *testing.T) {
	t.Parallel()

	p := ValidateTrust(TrustInput{
		Carrier:       &decision.TrustCarrier{State: "trusted-root-a"},
		Required:      false,
		ExpectedState: "trusted-root-b",
	})
	assertTrustViolation(t, p, CodeTrustMismatch, ptrTrustState)
}

func TestValidateTrust_ExpectedSetButCarrierAbsent(t *testing.T) {
	t.Parallel()

	p := ValidateTrust(TrustInput{
		Carrier:       nil,
		Required:      false,
		ExpectedState: "trusted-root-a",
	})
	assertTrustViolation(t, p, CodeTrustMismatch, ptrTrust)
}

func TestValidateTrust_FailClosedStates(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		state string
		code  apiproblem.ViolationCode
	}{
		{name: "unknown", state: TrustStateUnknown, code: CodeTrustUnknown},
		{name: "expired", state: TrustStateExpired, code: CodeTrustExpired},
		{name: "revoked", state: TrustStateRevoked, code: CodeTrustRevoked},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p := ValidateTrust(TrustInput{
				Carrier:  &decision.TrustCarrier{State: tc.state},
				Required: false,
			})
			assertTrustViolation(t, p, tc.code, ptrTrustState)
		})
	}
}

func TestValidateTrust_UnverifiedFailsWhenRequired(t *testing.T) {
	t.Parallel()

	p := ValidateTrust(TrustInput{
		Carrier:  &decision.TrustCarrier{State: TrustStateUnverified},
		Required: true,
	})
	assertTrustViolation(t, p, CodeTrustUnknown, ptrTrustState)
}

func TestValidateTrust_UnverifiedAllowedWhenOptional(t *testing.T) {
	t.Parallel()

	// Optional carrier may carry opaque "unverified" without claiming crypto validity.
	p := ValidateTrust(TrustInput{
		Carrier:  &decision.TrustCarrier{State: TrustStateUnverified},
		Required: false,
	})
	if p != nil {
		t.Fatalf("optional unverified opaque state must pass structurally, got %#v", p)
	}
}

func TestValidateTrust_MissingStateWhenCarrierPresent(t *testing.T) {
	t.Parallel()

	p := ValidateTrust(TrustInput{
		Carrier: &decision.TrustCarrier{
			DigestAlgorithm: "digest:profile-declared",
		},
		Required: false,
	})
	assertTrustViolation(t, p, CodeTrustUnknown, ptrTrustState)
}

func TestValidateTrust_IdentifierSyntaxCeiling(t *testing.T) {
	t.Parallel()

	tooLong := strings.Repeat("a", trustIDMaxLen+1)
	p := ValidateTrust(TrustInput{
		Carrier: &decision.TrustCarrier{
			State:           "trusted-root-a",
			DigestAlgorithm: tooLong,
		},
	})
	assertTrustViolation(t, p, CodeTrustUnknown, ptrTrustDigestAlgorithm)
}

func TestValidateTrust_NoCryptographicBehavior(t *testing.T) {
	t.Parallel()

	// Structural acceptance of opaque algorithm identifiers must not imply
	// canonicalization, digest, signing, or verification (F13-TRUST-004).
	p := ValidateTrust(TrustInput{
		Carrier: &decision.TrustCarrier{
			State:                  "opaque-token",
			CanonicalizationMethod: "illustrative-only",
			DigestAlgorithm:        "illustrative-only",
			CoveredFields:          "/record",
			SignatureAlgorithm:     "illustrative-only",
		},
		Required:      true,
		ExpectedState: "opaque-token",
	})
	if p != nil {
		t.Fatalf("opaque structural carrier must pass without crypto execution, got %#v", p)
	}
}

func TestValidateTrust_RequiredAbsentTakesPrecedenceOverExpected(t *testing.T) {
	t.Parallel()

	p := ValidateTrust(TrustInput{
		Carrier:       nil,
		Required:      true,
		ExpectedState: "trusted-root-a",
	})
	assertTrustViolation(t, p, CodeTrustRequired, ptrTrust)
}
