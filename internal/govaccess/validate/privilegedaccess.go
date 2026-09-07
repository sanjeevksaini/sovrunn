package validate

import (
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// ValidatePrivilegedAccessRequest validates PrivilegedAccessRequest rules
// (VS0-SCHEMA-063; REQ-F18-14/15).
func ValidatePrivilegedAccessRequest(r model.PrivilegedAccessRequest) *apiproblem.Problem {
	if prob := requireTypeMeta(r.TypeMeta, model.APIVersionIAM, "PrivilegedAccessRequest", ""); prob != nil {
		return prob
	}
	if prob := requireObjectName(r.Metadata); prob != nil {
		return prob
	}
	if prob := ValidateScopeApplicability(ContractPrivilegedAccessRequest, r.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if !r.Spec.Mode.Valid() {
		return invalidEnum("/spec/mode")
	}
	if prob := ValidatePrincipalRef(r.Spec.RequesterRef); prob != nil {
		return prob
	}
	if r.Spec.RequesterRef.PrincipalType != model.PrincipalTypeHuman {
		return validationFailed("/spec/requesterRef/principalType", "requester must be Human", apivalid.ViolationInvalidEnum)
	}
	if r.Spec.RoleDefinitionRef.Kind != RoleDefinitionKind {
		return validationFailed("/spec/roleDefinitionRef/kind", "roleDefinitionRef.kind must be RoleDefinition", apivalid.ViolationInvalidKind)
	}
	if prob := requireNonEmpty(r.Spec.RoleDefinitionRef.UID, "/spec/roleDefinitionRef/uid"); prob != nil {
		return prob
	}
	if r.Spec.ActivationDeadline.IsZero() {
		return missingField("/spec/activationDeadline")
	}
	if prob := requireNonEmpty(r.Spec.RequestedDuration, "/spec/requestedDuration"); prob != nil {
		return prob
	}
	return nil
}
