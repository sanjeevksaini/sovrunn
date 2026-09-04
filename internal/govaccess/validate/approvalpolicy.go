package validate

import (
	"strconv"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// ValidateApprovalPolicy validates ApprovalPolicy structural/local semantic rules
// (VS0-SCHEMA-065; REQ-F18-12).
func ValidateApprovalPolicy(p model.ApprovalPolicy) *apiproblem.Problem {
	if prob := requireTypeMeta(p.TypeMeta, model.APIVersionIAM, "ApprovalPolicy", ""); prob != nil {
		return prob
	}
	if prob := requireObjectName(p.Metadata); prob != nil {
		return prob
	}
	if prob := ValidateScopeApplicability(ContractApprovalPolicy, p.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := requireNonEmpty(p.Spec.Version, "/spec/version"); prob != nil {
		return prob
	}
	if !p.Spec.PublicationState.Valid() {
		return invalidEnum("/spec/publicationState")
	}
	if len(p.Spec.ApproverEligibility) == 0 {
		return missingField("/spec/approverEligibilityRefs")
	}
	for i, elig := range p.Spec.ApproverEligibility {
		if prob := ValidateEligibilityRef(elig); prob != nil {
			code := apivalid.ViolationInvalidEnum
			if len(prob.Violations) > 0 {
				code = prob.Violations[0].Code
			}
			return validationFailed("/spec/approverEligibilityRefs/"+strconv.Itoa(i), "invalid EligibilityRef", code)
		}
	}
	return nil
}
