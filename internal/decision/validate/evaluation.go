package validate

import (
	"strings"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

// RFC 6901 pointers for the evaluation validation pass (design §9.1 step 4).
// Pointers are rooted at the EvaluationResult value; orchestration (T-026) may
// prefix them with /record/evaluationResults/{i} when validating an embedded
// capture.
const (
	ptrEvalType            = "/evaluator/type"
	ptrEvalVersion         = "/evaluator/version"
	ptrEvalResultStatus    = "/resultStatus"
	ptrEvalResult          = "/result"
	ptrEvalTimingStarted   = "/timing/startedAt"
	ptrEvalTimingCompleted = "/timing/completedAt"
	ptrEvalTimingDuration  = "/timing/durationMs"
	ptrEvalEvaluatedAt     = "/evaluatedAt"
	ptrEvalInputSnapshot   = "/inputSnapshotRef"
	ptrEvalInputIntegrity  = "/inputIntegrity"
	ptrEvalIntegrityState  = "/inputIntegrity/state"
	ptrEvalScopeEvidence   = "/scopeEvidence"
)

// governanceScopeDepth returns nesting depth on the Platform→Project chain.
// Lower values are wider. Provider is incomparable to that chain (-1).
func governanceScopeDepth(kind apimeta.ScopeKind) int {
	switch kind {
	case apimeta.ScopePlatform:
		return 0
	case apimeta.ScopeOrganization:
		return 1
	case apimeta.ScopeOrganizationUnit:
		return 2
	case apimeta.ScopeTenant:
		return 3
	case apimeta.ScopeProject:
		return 4
	default:
		return -1
	}
}

// ValidateEvaluation runs the FEATURE-0013 evaluation validation pass
// (design §9.1 step 4; F13-EVAL-001…005; AD-005).
//
// Checks (fail closed, short-circuit on first blocking violation):
//  1. evaluator type/version accepted by the governing profile
//  2. EvaluationResult structural requiredness and closed vocabularies
//  3. scopeEvidence equals or is a structural narrowing of recScope
//     (evidence only; never a second scope authority; no existence disclosure)
//  4. observed and registration limits against the inherited hierarchy
//
// Evaluator-supplied scope metadata is evidence only. Canonical scope remains
// the containing DecisionRecord's metadata.scopeRef (passed as recScope).
func ValidateEvaluation(e decision.EvaluationResult, recScope apimeta.ScopeIdentity, p decision.DecisionProfile) *apiproblem.Problem {
	// Structural validity before registry membership so empty identity maps to
	// DECISION_EVALUATION_RESULT_INVALID rather than TYPE_UNSUPPORTED.
	if prob := checkEvaluationResultStructure(e); prob != nil {
		return prob
	}

	accepted, ok := findAcceptedEvaluationType(e.Evaluator, p.Spec.AcceptedEvaluationTypes)
	if !ok {
		return evaluationProblem(CodeEvaluationTypeUnsupported, ptrEvalType,
			"evaluator type is not accepted by the governing decision profile")
	}

	if prob := checkEvaluationScopeEvidence(e.ScopeEvidence, recScope); prob != nil {
		return prob
	}
	if prob := checkEvaluationLimits(e, accepted, p); prob != nil {
		return prob
	}
	return nil
}

func findAcceptedEvaluationType(id decision.EvaluatorIdentity, accepted []decision.AcceptedEvaluationType) (decision.AcceptedEvaluationType, bool) {
	typeName := strings.TrimSpace(id.Type)
	if typeName == "" {
		return decision.AcceptedEvaluationType{}, false
	}
	version := strings.TrimSpace(id.Version)
	for _, a := range accepted {
		if a.Type != typeName {
			continue
		}
		if a.Version != "" && a.Version != version {
			continue
		}
		return a, true
	}
	return decision.AcceptedEvaluationType{}, false
}

func checkEvaluationResultStructure(e decision.EvaluationResult) *apiproblem.Problem {
	if strings.TrimSpace(e.Evaluator.Type) == "" {
		return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalType,
			"evaluator.type is required and must be non-empty")
	}
	if strings.TrimSpace(e.Evaluator.Version) == "" {
		return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalVersion,
			"evaluator.version is required and must be non-empty")
	}
	if !e.ResultStatus.Valid() {
		return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalResultStatus,
			"resultStatus must be a closed EvaluationResultStatus value")
	}

	if e.InputSnapshotRef == nil && e.InputIntegrity == nil {
		return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalInputSnapshot,
			"evaluation requires an inputSnapshotRef or inputIntegrity carrier")
	}
	if e.InputSnapshotRef != nil {
		if strings.TrimSpace(e.InputSnapshotRef.APIVersion) == "" ||
			strings.TrimSpace(e.InputSnapshotRef.Kind) == "" ||
			strings.TrimSpace(e.InputSnapshotRef.Name) == "" {
			return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalInputSnapshot,
				"inputSnapshotRef requires apiVersion, kind, and name")
		}
	}
	if e.InputIntegrity != nil {
		if integrityCarrierEmpty(e.InputIntegrity) {
			return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalInputIntegrity,
				"present inputIntegrity carrier must not be structurally empty")
		}
		if strings.TrimSpace(e.InputIntegrity.State) == "" {
			return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalIntegrityState,
				"inputIntegrity.state is required when inputIntegrity is present")
		}
	}

	if strings.TrimSpace(e.Timing.StartedAt) == "" {
		return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalTimingStarted,
			"timing.startedAt is required and must be UTC RFC3339")
	}
	if strings.TrimSpace(e.Timing.CompletedAt) == "" {
		return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalTimingCompleted,
			"timing.completedAt is required and must be UTC RFC3339")
	}
	if _, err := time.Parse(time.RFC3339, e.Timing.StartedAt); err != nil {
		// Also accept RFC3339Nano which encoding often emits.
		if _, err2 := time.Parse(time.RFC3339Nano, e.Timing.StartedAt); err2 != nil {
			return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalTimingStarted,
				"timing.startedAt must be UTC RFC3339")
		}
	}
	if _, err := time.Parse(time.RFC3339, e.Timing.CompletedAt); err != nil {
		if _, err2 := time.Parse(time.RFC3339Nano, e.Timing.CompletedAt); err2 != nil {
			return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalTimingCompleted,
				"timing.completedAt must be UTC RFC3339")
		}
	}
	if e.Timing.DurationMs < 0 {
		return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalTimingDuration,
			"timing.durationMs must not be negative")
	}
	if strings.TrimSpace(e.EvaluatedAt) == "" {
		return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalEvaluatedAt,
			"evaluatedAt is required and must be UTC RFC3339")
	}
	if _, err := time.Parse(time.RFC3339, e.EvaluatedAt); err != nil {
		if _, err2 := time.Parse(time.RFC3339Nano, e.EvaluatedAt); err2 != nil {
			return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalEvaluatedAt,
				"evaluatedAt must be UTC RFC3339")
		}
	}

	// SUCCESS requires a present normalized result payload (F13-EVAL-001).
	if e.ResultStatus == decision.EvaluationResultStatusSuccess && len(e.Result) == 0 {
		return evaluationProblem(CodeEvaluationResultInvalid, ptrEvalResult,
			"SUCCESS evaluation requires a present normalized result payload")
	}

	return nil
}

func integrityCarrierEmpty(c *decision.OpaqueIntegrityCarrier) bool {
	if c == nil {
		return true
	}
	return strings.TrimSpace(c.State) == "" &&
		strings.TrimSpace(c.CanonicalizationMethod) == "" &&
		strings.TrimSpace(c.DigestAlgorithm) == "" &&
		strings.TrimSpace(c.CoveredFields) == "" &&
		strings.TrimSpace(c.SignatureAlgorithm) == ""
}

// checkEvaluationScopeEvidence enforces F13-EVAL-003/004: evaluator scope
// metadata is evidence only. It must equal the record scope or be a structural
// narrowing on the Platform→Project chain. Unknown scope, unauthorized
// widening, mismatch, or inaccessible reference fails closed with
// DECISION_EVALUATION_SCOPE_MISMATCH and must not disclose target existence.
func checkEvaluationScopeEvidence(evidence *apimeta.ScopeRef, recScope apimeta.ScopeIdentity) *apiproblem.Problem {
	if evidence == nil {
		return nil
	}

	normalized := apimeta.NormalizeScope(evidence)
	if normalized != nil {
		kind := apimeta.ScopeKind(normalized.Kind)
		if !kind.Valid() {
			return evaluationProblem(CodeEvaluationScopeMismatch, ptrEvalScopeEvidence,
				"evaluation scope evidence is not permitted for the containing record scope")
		}
		if normalized.UID == "" {
			return evaluationProblem(CodeEvaluationScopeMismatch, ptrEvalScopeEvidence,
				"evaluation scope evidence is not permitted for the containing record scope")
		}
	}

	evID := apimeta.CanonicalScopeIdentity(normalized)
	if evID == recScope {
		return nil
	}

	recDepth := governanceScopeDepth(recScope.Kind)
	evDepth := governanceScopeDepth(evID.Kind)
	if recDepth < 0 || evDepth < 0 {
		// Provider or otherwise incomparable kinds: only exact identity is accepted.
		return evaluationProblem(CodeEvaluationScopeMismatch, ptrEvalScopeEvidence,
			"evaluation scope evidence is not permitted for the containing record scope")
	}
	if evDepth < recDepth {
		// Unauthorized widening (evidence broader than record).
		return evaluationProblem(CodeEvaluationScopeMismatch, ptrEvalScopeEvidence,
			"evaluation scope evidence is not permitted for the containing record scope")
	}
	if evDepth == recDepth {
		// Same nesting depth but different identity (kind or uid).
		return evaluationProblem(CodeEvaluationScopeMismatch, ptrEvalScopeEvidence,
			"evaluation scope evidence is not permitted for the containing record scope")
	}
	// Structural narrowing (deeper kind than record) is permitted as evidence
	// only; it never becomes a second scope authority.
	return nil
}

func checkEvaluationLimits(e decision.EvaluationResult, accepted decision.AcceptedEvaluationType, p decision.DecisionProfile) *apiproblem.Problem {
	profileLim := p.Spec.Limits
	plat := InheritedPlatformCeilings()

	// Edge 11: evaluator registration may only narrow profile ceilings.
	if accepted.Limits != nil {
		if prob := checkAcceptedEvaluatorLimits(*accepted.Limits, profileLim, plat); prob != nil {
			return prob
		}
	}

	var evalLim decision.ProfileLimits
	if accepted.Limits != nil {
		evalLim = *accepted.Limits
	}

	maxLatency := EffectiveInt64Limit(0, profileLim.MaxLatencyMs, evalLim.MaxLatencyMs)
	if maxLatency > 0 && e.Timing.DurationMs > maxLatency {
		return evaluationProblem(CodeEvaluationLimitExceeded, ptrEvalTimingDuration,
			"evaluation duration exceeds the effective MaxLatencyMs limit")
	}

	maxPayload := EffectiveInt64Limit(plat.MaxPayloadBytes, profileLim.MaxPayloadBytes, evalLim.MaxPayloadBytes)
	if maxPayload > 0 && int64(len(e.Result)) > maxPayload {
		return evaluationProblem(CodeEvaluationLimitExceeded, ptrEvalResult,
			"evaluation result payload exceeds the effective MaxPayloadBytes limit")
	}

	return nil
}

func checkAcceptedEvaluatorLimits(eval, profile decision.ProfileLimits, plat PlatformLimitCeilings) *apiproblem.Problem {
	checks := []struct {
		name     string
		eval     int64
		profile  int64
		platform int64
		field    string
	}{
		{"maxNodes", int64(eval.MaxNodes), int64(profile.MaxNodes), 0, "/evaluator"},
		{"maxEdges", int64(eval.MaxEdges), int64(profile.MaxEdges), 0, "/evaluator"},
		{"maxDepth", int64(eval.MaxDepth), int64(profile.MaxDepth), int64(plat.MaxDepth), "/evaluator"},
		{"maxWidth", int64(eval.MaxWidth), int64(profile.MaxWidth), 0, "/evaluator"},
		{"maxFanOut", int64(eval.MaxFanOut), int64(profile.MaxFanOut), 0, "/evaluator"},
		{"maxPayloadBytes", eval.MaxPayloadBytes, profile.MaxPayloadBytes, plat.MaxPayloadBytes, "/evaluator"},
		{"maxFieldCount", int64(eval.MaxFieldCount), int64(profile.MaxFieldCount), int64(plat.MaxFieldCount), "/evaluator"},
		{"maxCalls", int64(eval.MaxCalls), int64(profile.MaxCalls), 0, "/evaluator"},
		{"maxConcurrency", int64(eval.MaxConcurrency), int64(profile.MaxConcurrency), 0, "/evaluator"},
		{"maxRetries", int64(eval.MaxRetries), int64(profile.MaxRetries), 0, "/evaluator"},
		{"maxLatencyMs", eval.MaxLatencyMs, profile.MaxLatencyMs, 0, "/evaluator"},
		{"maxJurisdictions", int64(eval.MaxJurisdictions), int64(profile.MaxJurisdictions), 0, "/evaluator"},
	}
	for _, c := range checks {
		if c.eval <= 0 {
			continue
		}
		if c.platform > 0 && c.eval > c.platform {
			return evaluationProblem(CodeEvaluationLimitExceeded, c.field,
				"evaluator registration limit exceeds the FEATURE-0012 outer bound")
		}
		if c.profile <= 0 {
			return evaluationProblem(CodeEvaluationLimitExceeded, c.field,
				"evaluator registration limit is incomparable without a governing profile bound")
		}
		if c.eval > c.profile {
			return evaluationProblem(CodeEvaluationLimitExceeded, c.field,
				"evaluator registration may only narrow profile limits; widening is prohibited")
		}
	}
	if eval.MaxCostBudget != "" {
		if profile.MaxCostBudget == "" || eval.MaxCostBudget != profile.MaxCostBudget {
			return evaluationProblem(CodeEvaluationLimitExceeded, "/evaluator",
				"evaluator registration maxCostBudget must equal the governing profile opaque budget tag")
		}
	}
	return nil
}

func evaluationProblem(code apiproblem.ViolationCode, field, message string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail(message).
		WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    code,
			Message: message,
		}})
}
