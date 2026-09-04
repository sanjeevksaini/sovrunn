package validate

import (
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// ValidateApprovalRequest validates ApprovalRequest structural/local semantic rules
// (VS0-SCHEMA-066; REQ-F18-13).
func ValidateApprovalRequest(r model.ApprovalRequest) *apiproblem.Problem {
	if prob := requireTypeMeta(r.TypeMeta, model.APIVersionIAM, "ApprovalRequest", ""); prob != nil {
		return prob
	}
	if prob := requireObjectName(r.Metadata); prob != nil {
		return prob
	}
	if prob := ValidateScopeApplicability(ContractApprovalRequest, r.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if r.Spec.ApprovalPolicyRef.Kind != "ApprovalPolicy" {
		return validationFailed("/spec/approvalPolicyRef/kind", "approvalPolicyRef.kind must be ApprovalPolicy", apivalid.ViolationInvalidKind)
	}
	if prob := requireNonEmpty(r.Spec.ApprovalPolicyRef.UID, "/spec/approvalPolicyRef/uid"); prob != nil {
		return prob
	}
	if prob := requireNonEmpty(r.Spec.SubjectKind, "/spec/subjectKind"); prob != nil {
		return prob
	}
	if prob := requireNonEmpty(r.Spec.SubjectUID, "/spec/subjectUid"); prob != nil {
		return prob
	}
	if r.Status.Decision != "" && !r.Status.Decision.Valid() {
		return invalidEnum("/status/decision")
	}
	return nil
}
