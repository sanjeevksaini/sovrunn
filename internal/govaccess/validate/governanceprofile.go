package validate

import (
	"strconv"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// ValidateGovernanceProfile validates GovernanceProfile structural/local semantic
// rules (VS0-SCHEMA-024; REQ-F18-18).
func ValidateGovernanceProfile(p model.GovernanceProfile) *apiproblem.Problem {
	if prob := requireTypeMeta(p.TypeMeta, model.APIVersionGov, "GovernanceProfile", ""); prob != nil {
		return prob
	}
	if prob := requireObjectName(p.Metadata); prob != nil {
		return prob
	}
	if prob := ValidateScopeApplicability(ContractGovernanceProfile, p.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := requireNonEmpty(p.Spec.Version, "/spec/version"); prob != nil {
		return prob
	}
	if !p.Spec.PublicationState.Valid() {
		return invalidEnum("/spec/publicationState")
	}
	// GovernanceProfile v1 publicationState set excludes Suspended (VS0-SCHEMA-024).
	if p.Spec.PublicationState == model.PublicationStateSuspended {
		return invalidEnum("/spec/publicationState")
	}
	for i, ref := range p.Spec.ApprovalPolicyRefs {
		if ref.Kind != "ApprovalPolicy" {
			return validationFailed("/spec/approvalPolicyRefs/"+strconv.Itoa(i)+"/kind", "approvalPolicyRefs kind must be ApprovalPolicy", apivalid.ViolationInvalidKind)
		}
		if ref.UID == "" {
			return missingField("/spec/approvalPolicyRefs/" + strconv.Itoa(i) + "/uid")
		}
	}
	return nil
}
