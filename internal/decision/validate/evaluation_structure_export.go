package validate

import "github.com/sanjeevksaini/sovrunn/internal/decision"

// IsEvaluationResultStructurallyValid exposes the existing FEATURE-0013
// structural pass without profile, record-scope, limit, Problem, or violation
// output. It delegates directly to checkEvaluationResultStructure and adds no
// rule (ADH-2026-068 / D-14 approved dependency extension).
func IsEvaluationResultStructurallyValid(result decision.EvaluationResult) bool {
	return checkEvaluationResultStructure(result) == nil
}
