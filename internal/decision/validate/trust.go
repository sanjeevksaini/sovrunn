package validate

import (
	"strings"
	"unicode/utf8"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

// RFC 6901 pointers for the trust validation pass (design §9.1 step 6).
// Pointers are rooted at the TrustCarrier value; orchestration (T-026) may
// prefix them with /record/trust or /spec/trust when validating an embedded
// carrier.
const (
	ptrTrust      = "/trust"
	ptrTrustState = "/trust/state"

	ptrTrustCanonicalization = "/trust/canonicalizationMethod"
	ptrTrustDigestAlgorithm  = "/trust/digestAlgorithm"
	ptrTrustCoveredFields    = "/trust/coveredFields"
	ptrTrustSignatureAlg     = "/trust/signatureAlgorithm"
)

// Opaque structural trust-state tokens that fail closed (design §9.2;
// F13-TRUST-003; AD-022). These are not a new closed vocabulary type; they are
// the opaque fail-closed transitions architecture requires before ADR-F13-002
// cryptographic selection. Values match the bundle offline trust gate.
const (
	TrustStateUnknown    = "unknown"
	TrustStateExpired    = "expired"
	TrustStateRevoked    = "revoked"
	TrustStateUnverified = "unverified"
)

// Schema-aligned structural identifier ceilings (api/schemas/_common/trust-carrier.json).
const (
	trustIDMaxLen      = 253
	trustCoveredMaxLen = 1024
)

// TrustInput is the input to ValidateTrust (design §8 / §9.1 step 6;
// F13-TRUST-001…005; RID-08; AD-022).
//
// Design §8 summary signature is (c *TrustCarrier, required bool).
// ExpectedState is the same pass's additional input required by F13-TRUST-002
// opaque expected-versus-observed comparison.
type TrustInput struct {
	Carrier  *decision.TrustCarrier
	Required bool

	// ExpectedState is the opaque expected trust-state token. When non-empty,
	// it must equal Carrier.State; mismatch fails closed with no digest,
	// signature, or cryptographic check.
	ExpectedState string
}

// ValidateTrust runs the FEATURE-0013 trust validation pass (design §9.1
// step 6; RID-08; F13-TRUST-001…005; AD-022).
//
// Checks (fail closed, short-circuit on first blocking violation):
//  1. required carrier absent → DECISION_TRUST_REQUIRED
//  2. optional absent with no expected state → accept
//  3. expected state set but carrier absent → DECISION_TRUST_MISMATCH
//  4. present-but-structurally-empty → DECISION_TRUST_UNKNOWN (RID-08)
//  5. state required when carrier present; identifier syntax ceilings
//  6. opaque fail-closed states (unknown/expired/revoked; unverified when required)
//  7. opaque expected≠observed → DECISION_TRUST_MISMATCH
//
// No canonicalization, digest computation, signing, or cryptographic
// verification is performed (F13-TRUST-004; ADR-F13-002).
func ValidateTrust(in TrustInput) *apiproblem.Problem {
	if in.Carrier == nil {
		if in.Required {
			return trustProblem(CodeTrustRequired, ptrTrust,
				"required trust carrier is absent")
		}
		if strings.TrimSpace(in.ExpectedState) != "" {
			return trustProblem(CodeTrustMismatch, ptrTrust,
				"expected trust state set but trust carrier absent")
		}
		return nil
	}

	if in.Carrier.IsStructurallyEmpty() {
		// RID-08: present-but-empty is malformed; never satisfies required trust.
		return trustProblem(CodeTrustUnknown, ptrTrust,
			"present-but-empty trust carrier is invalid")
	}

	if prob := checkTrustCarrierSyntax(in.Carrier); prob != nil {
		return prob
	}

	state := in.Carrier.State
	switch state {
	case TrustStateUnknown:
		return trustProblem(CodeTrustUnknown, ptrTrustState,
			"unknown trust state fails closed")
	case TrustStateExpired:
		return trustProblem(CodeTrustExpired, ptrTrustState,
			"expired trust state fails closed")
	case TrustStateRevoked:
		return trustProblem(CodeTrustRevoked, ptrTrustState,
			"revoked trust state fails closed")
	case TrustStateUnverified:
		// Architecture: when verified trust is required, unverified fails closed.
		if in.Required {
			return trustProblem(CodeTrustUnknown, ptrTrustState,
				"unverified trust state fails closed when verified trust is required")
		}
	}

	expected := strings.TrimSpace(in.ExpectedState)
	if expected != "" && state != expected {
		return trustProblem(CodeTrustMismatch, ptrTrustState,
			"expected versus observed trust state mismatch; no cryptographic check")
	}
	return nil
}

// checkTrustCarrierSyntax validates structural requiredness and identifier
// ceilings only (F13-TRUST-001/002). No algorithm is selected or executed.
func checkTrustCarrierSyntax(c *decision.TrustCarrier) *apiproblem.Problem {
	if strings.TrimSpace(c.State) == "" {
		return trustProblem(CodeTrustUnknown, ptrTrustState,
			"trust state is required when the trust carrier is present")
	}
	if utf8.RuneCountInString(c.State) > trustIDMaxLen {
		return trustProblem(CodeTrustUnknown, ptrTrustState,
			"trust state exceeds structural identifier length")
	}
	if c.CanonicalizationMethod != "" && utf8.RuneCountInString(c.CanonicalizationMethod) > trustIDMaxLen {
		return trustProblem(CodeTrustUnknown, ptrTrustCanonicalization,
			"canonicalizationMethod exceeds structural identifier length")
	}
	if c.DigestAlgorithm != "" && utf8.RuneCountInString(c.DigestAlgorithm) > trustIDMaxLen {
		return trustProblem(CodeTrustUnknown, ptrTrustDigestAlgorithm,
			"digestAlgorithm exceeds structural identifier length")
	}
	if c.CoveredFields != "" && utf8.RuneCountInString(c.CoveredFields) > trustCoveredMaxLen {
		return trustProblem(CodeTrustUnknown, ptrTrustCoveredFields,
			"coveredFields exceeds structural descriptor length")
	}
	if c.SignatureAlgorithm != "" && utf8.RuneCountInString(c.SignatureAlgorithm) > trustIDMaxLen {
		return trustProblem(CodeTrustUnknown, ptrTrustSignatureAlg,
			"signatureAlgorithm exceeds structural identifier length")
	}
	return nil
}

func trustProblem(code apiproblem.ViolationCode, field, message string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail(message).
		WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    code,
			Message: message,
		}})
}
