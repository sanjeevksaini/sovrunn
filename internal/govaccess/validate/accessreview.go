package validate

import (
	"strconv"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// ValidateAccessReview validates AccessReview structural/local semantic rules
// (VS0-SCHEMA-064; REQ-F18-16).
func ValidateAccessReview(r model.AccessReview) *apiproblem.Problem {
	if prob := requireTypeMeta(r.TypeMeta, model.APIVersionIAM, "AccessReview", ""); prob != nil {
		return prob
	}
	if prob := requireObjectName(r.Metadata); prob != nil {
		return prob
	}
	if prob := ValidateScopeApplicability(ContractAccessReview, r.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := requireNonEmpty(r.Spec.AccessReviewRuleRef.UID, "/spec/accessReviewRuleRef/uid"); prob != nil {
		return prob
	}
	if len(r.Spec.ReviewerEligibility) == 0 {
		return missingField("/spec/reviewerEligibilityRefs")
	}
	for i, elig := range r.Spec.ReviewerEligibility {
		if prob := ValidateEligibilityRef(elig); prob != nil {
			code := apivalid.ViolationInvalidEnum
			if len(prob.Violations) > 0 {
				code = prob.Violations[0].Code
			}
			return validationFailed("/spec/reviewerEligibilityRefs/"+strconv.Itoa(i), "invalid EligibilityRef", code)
		}
	}
	return nil
}
