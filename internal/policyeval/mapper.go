package policyeval

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	decisionvalidate "github.com/sanjeevksaini/sovrunn/internal/decision/validate"
)

// errMappingInputInvalid is the sole sanitized mapper diagnostic (D-08 / §5.4).
var errMappingInputInvalid = errors.New("policy evaluation mapping input invalid")

// InputIdentity is the mapper-owned in-process FEATURE-0013 input-identity
// contract (D-14 / CDG-F17-05). Valid combinations are snapshot only,
// integrity only, or both. Both nil is invalid. It has no JSON tags or
// serialization path.
type InputIdentity struct {
	InputSnapshotRef *apimeta.TypedRef
	InputIntegrity   *decision.OpaqueIntegrityCarrier
}

// MapToEvaluationResult is the pure FEATURE-0013 evidence mapper (D-14).
// It validates structural mapper inputs, applies the closed status mapping,
// populates exactly the CDG-F17-05 field set, and delegates structural
// FEATURE-0013 validation to IsEvaluationResultStructurallyValid.
//
// It performs no time read, lookup, profile/scope validation, write,
// persistence, or external effect. On failure it returns a zero
// EvaluationResult and errMappingInputInvalid without mutating inputs.
func MapToEvaluationResult(
	result PolicyEvaluationResult,
	timing decision.TimingEnvelope,
	evaluator decision.EvaluatorIdentity,
	ident InputIdentity,
) (decision.EvaluationResult, error) {
	var zero decision.EvaluationResult

	if err := validateMapperPolicyResult(result); err != nil {
		return zero, errMappingInputInvalid
	}
	if !isCanonicalRFC3339NanoUTC(timing.StartedAt) || !isCanonicalRFC3339NanoUTC(timing.CompletedAt) {
		return zero, errMappingInputInvalid
	}
	if result.EvaluatedAt != timing.CompletedAt {
		return zero, errMappingInputInvalid
	}
	if strings.TrimSpace(evaluator.Type) == "" || strings.TrimSpace(evaluator.Version) == "" {
		return zero, errMappingInputInvalid
	}
	if err := validateInputIdentity(ident); err != nil {
		return zero, errMappingInputInvalid
	}

	status, ok := mapOutcomeToResultStatus(result.Outcome)
	if !ok {
		return zero, errMappingInputInvalid
	}

	ownedResult := copyPolicyEvaluationResult(result)
	raw, err := json.Marshal(ownedResult)
	if err != nil {
		return zero, errMappingInputInvalid
	}
	rawCopy := append(json.RawMessage(nil), raw...)

	out := decision.EvaluationResult{
		Evaluator: decision.EvaluatorIdentity{
			Type:    evaluator.Type,
			Version: evaluator.Version,
		},
		ResultStatus: status,
		Result:       rawCopy,
		Timing: decision.TimingEnvelope{
			StartedAt:   timing.StartedAt,
			CompletedAt: timing.CompletedAt,
			// DurationMs intentionally omitted (zero).
		},
		EvaluatedAt: result.EvaluatedAt,
		// Exact omitted fields remain zero: ExecutionConfig, TrustBoundary,
		// SafetyPolicyFilters, StructuralTrustState, ScopeEvidence.
	}
	if ident.InputSnapshotRef != nil {
		cp := *ident.InputSnapshotRef
		out.InputSnapshotRef = &cp
	}
	if ident.InputIntegrity != nil {
		cp := *ident.InputIntegrity
		out.InputIntegrity = &cp
	}

	if !decisionvalidate.IsEvaluationResultStructurallyValid(out) {
		return zero, errMappingInputInvalid
	}
	return out, nil
}

func mapOutcomeToResultStatus(o Outcome) (decision.EvaluationResultStatus, bool) {
	switch o {
	case OutcomeAllow, OutcomeDeny, OutcomeRequiresApproval:
		return decision.EvaluationResultStatusSuccess, true
	case OutcomeIndeterminate:
		return decision.EvaluationResultStatusIndeterminate, true
	default:
		return "", false
	}
}

func validateMapperPolicyResult(result PolicyEvaluationResult) error {
	if !result.Outcome.Valid() {
		return errMappingInputInvalid
	}
	n := len(result.ReasonCodes)
	if n < MinReasonCodes || n > MaxReasonCodes {
		return errMappingInputInvalid
	}
	seen := make(map[string]struct{}, n)
	for _, code := range result.ReasonCodes {
		if !ReasonCodeRegexp.MatchString(code) {
			return errMappingInputInvalid
		}
		if _, dup := seen[code]; dup {
			return errMappingInputInvalid
		}
		seen[code] = struct{}{}
	}
	if !validInputDigest(result.InputDigest) {
		return errMappingInputInvalid
	}
	// Phase 2R: every non-empty obligations value is rejected (D-14 §5.5).
	if len(result.Obligations) > 0 {
		return errMappingInputInvalid
	}
	if !isCanonicalRFC3339NanoUTC(result.EvaluatedAt) {
		return errMappingInputInvalid
	}
	return nil
}

func validateInputIdentity(ident InputIdentity) error {
	if ident.InputSnapshotRef == nil && ident.InputIntegrity == nil {
		return errMappingInputInvalid
	}
	if ident.InputSnapshotRef != nil {
		if issues := (apiref.Constraint{}).ValidateRef(*ident.InputSnapshotRef, "/inputSnapshotRef"); len(issues) > 0 {
			return errMappingInputInvalid
		}
	}
	if ident.InputIntegrity != nil {
		if integrityCarrierStructurallyInvalid(ident.InputIntegrity) {
			return errMappingInputInvalid
		}
	}
	return nil
}

func integrityCarrierStructurallyInvalid(c *decision.OpaqueIntegrityCarrier) bool {
	if c == nil {
		return true
	}
	// Empty carrier and blank/partial State are invalid (D-14 §4.6).
	if strings.TrimSpace(c.State) == "" {
		return true
	}
	return false
}

// isCanonicalRFC3339NanoUTC reports whether s is exactly the boundary's
// canonical UTC time.RFC3339Nano representation (Z form, round-trip stable).
func isCanonicalRFC3339NanoUTC(s string) bool {
	if s == "" {
		return false
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return false
	}
	return t.UTC().Format(time.RFC3339Nano) == s
}

func copyPolicyEvaluationResult(in PolicyEvaluationResult) PolicyEvaluationResult {
	out := PolicyEvaluationResult{
		Outcome:     in.Outcome,
		InputDigest: in.InputDigest,
		EvaluatedAt: in.EvaluatedAt,
	}
	if in.ReasonCodes != nil {
		out.ReasonCodes = append([]string(nil), in.ReasonCodes...)
	}
	// Non-empty obligations are rejected before copy; nil/empty stay absent.
	if len(in.Obligations) == 0 {
		out.Obligations = nil
	}
	return out
}
