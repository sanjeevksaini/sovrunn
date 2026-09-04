package validate

import (
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// ValidateAccessGroup validates AccessGroup structural/local semantic rules
// (VS0-SCHEMA-062; REQ-F18-04).
func ValidateAccessGroup(g model.AccessGroup) *apiproblem.Problem {
	if prob := requireTypeMeta(g.TypeMeta, model.APIVersionIAM, model.KindAccessGroup, ""); prob != nil {
		return prob
	}
	if prob := requireObjectName(g.Metadata); prob != nil {
		return prob
	}
	if prob := ValidateScopeApplicability(ContractAccessGroup, g.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := ValidatePrincipalRef(g.Spec.OwnerRef); prob != nil {
		return prob
	}
	if g.Spec.OwnerRef.PrincipalType != model.PrincipalTypeHuman {
		return validationFailed("/spec/ownerRef/principalType", "AccessGroup owner must be Human", apivalid.ViolationInvalidEnum)
	}
	if g.Status.Phase != "" && !g.Status.Phase.Valid() {
		return invalidEnum("/status/phase")
	}
	return nil
}
