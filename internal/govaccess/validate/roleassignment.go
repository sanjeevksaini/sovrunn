package validate

import (
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// ValidateRoleAssignment validates RoleAssignment structural/local semantic rules
// (VS0-SCHEMA-023; REQ-F18-08). When the holder is an AccessGroup, groupScope
// must be supplied for F18-RD-02 reference-compatibility checks.
func ValidateRoleAssignment(a model.RoleAssignment, groupScope *apimeta.ScopeRef) *apiproblem.Problem {
	if prob := requireTypeMeta(a.TypeMeta, model.APIVersionIAM, "RoleAssignment", ""); prob != nil {
		return prob
	}
	if prob := requireObjectName(a.Metadata); prob != nil {
		return prob
	}
	if prob := ValidateScopeApplicability(ContractRoleAssignment, a.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := ValidateRoleHolderRef(a.Spec.RoleHolderRef); prob != nil {
		return prob
	}
	if a.Spec.RoleDefinitionRef.Kind != RoleDefinitionKind {
		return validationFailed("/spec/roleDefinitionRef/kind", "roleDefinitionRef.kind must be RoleDefinition", apivalid.ViolationInvalidKind)
	}
	if prob := requireNonEmpty(a.Spec.RoleDefinitionRef.UID, "/spec/roleDefinitionRef/uid"); prob != nil {
		return prob
	}
	if prob := requireNonEmpty(a.Spec.RoleDefinitionVersion, "/spec/roleDefinitionVersion"); prob != nil {
		return prob
	}
	if !a.Spec.Validity.Valid() {
		return invalidEnum("/spec/validity")
	}
	if a.Spec.Validity == model.AssignmentValidityTimeBound {
		if a.Spec.NotBefore == nil || a.Spec.ExpiresAt == nil {
			return missingField("/spec/notBefore")
		}
		if !a.Spec.ExpiresAt.After(*a.Spec.NotBefore) {
			return validationFailed("/spec/expiresAt", "expiresAt must be after notBefore", apiproblem.ViolationOutOfRange)
		}
	}
	if a.Spec.RoleHolderRef.IsAccessGroup() {
		if a.Spec.MembershipRef != nil {
			return validationFailed("/spec/membershipRef", "membershipRef is prohibited for AccessGroup holders", apivalid.ViolationInvalidEnum)
		}
		if prob := AccessGroupHolderScopeCompatibility(groupScope, a.Metadata.ScopeRef); prob != nil {
			return prob
		}
	}
	if a.Spec.RoleHolderRef.IsPrincipal() {
		pt := a.Spec.RoleHolderRef.Principal.PrincipalType
		if pt == model.PrincipalTypeWorkload || pt == model.PrincipalTypeSystem {
			if a.Spec.ResponsiblePartyRef == nil {
				return missingField("/spec/responsiblePartyRef")
			}
			if prob := ValidatePrincipalRef(*a.Spec.ResponsiblePartyRef); prob != nil {
				return prob
			}
			if a.Spec.ResponsiblePartyRef.PrincipalType != model.PrincipalTypeHuman {
				return validationFailed("/spec/responsiblePartyRef/principalType", "responsiblePartyRef must be Human", apivalid.ViolationInvalidEnum)
			}
		}
	}
	if a.Spec.Validity == model.AssignmentValidityStanding && a.Spec.AccessReviewRuleRef == nil {
		return missingField("/spec/accessReviewRuleRef")
	}
	if a.Status.Phase == "" {
		return missingField("/status/phase")
	}
	return nil
}
