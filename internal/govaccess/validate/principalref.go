package validate

import (
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// ValidatePrincipalRef validates a PrincipalRef value (VS0-SCHEMA-020; REQ-F18-03).
func ValidatePrincipalRef(ref model.PrincipalRef) *apiproblem.Problem {
	if prob := requireNonEmpty(ref.Issuer, "/issuer"); prob != nil {
		return prob
	}
	if prob := requireNonEmpty(ref.Subject, "/subject"); prob != nil {
		return prob
	}
	if !ref.PrincipalType.Valid() {
		return invalidEnum("/principalType")
	}
	return nil
}

// ValidateAccessGroupRef validates a UID-pinned AccessGroupRef (DD-03; REQ-F18-04).
func ValidateAccessGroupRef(ref model.AccessGroupRef) *apiproblem.Problem {
	if ref.APIVersion != "" && ref.APIVersion != model.APIVersionIAM {
		return validationFailed("/apiVersion", "AccessGroupRef apiVersion must be iam.sovrunn.io/v1alpha1", apivalid.ViolationInvalidAPIVersion)
	}
	if ref.Kind != model.KindAccessGroup {
		return validationFailed("/kind", "AccessGroupRef kind must be AccessGroup", apivalid.ViolationInvalidKind)
	}
	if prob := requireNonEmpty(ref.Name, "/name"); prob != nil {
		return prob
	}
	if prob := requireNonEmpty(ref.UID, "/uid"); prob != nil {
		return prob
	}
	return nil
}

// ValidateRoleHolderRef validates the closed Principal|AccessGroup union (REQ-F18-04).
func ValidateRoleHolderRef(ref model.RoleHolderRef) *apiproblem.Problem {
	if !ref.Kind.Valid() {
		return invalidEnum("/kind")
	}
	switch ref.Kind {
	case model.RoleHolderPrincipal:
		if ref.Principal == nil || ref.AccessGroup != nil {
			return validationFailed("/principal", "RoleHolderRef Principal variant requires exactly one PrincipalRef", apivalid.ViolationInvalidEnum)
		}
		return ValidatePrincipalRef(*ref.Principal)
	case model.RoleHolderAccessGroup:
		if ref.AccessGroup == nil || ref.Principal != nil {
			return validationFailed("/accessGroup", "RoleHolderRef AccessGroup variant requires exactly one AccessGroupRef", apivalid.ViolationInvalidEnum)
		}
		return ValidateAccessGroupRef(*ref.AccessGroup)
	default:
		return invalidEnum("/kind")
	}
}

// ValidateEligibilityRef validates exact-Human eligibility (REQ-F18-12; SEC-06).
func ValidateEligibilityRef(ref model.EligibilityRef) *apiproblem.Problem {
	if prob := ValidatePrincipalRef(ref.Principal); prob != nil {
		return prob
	}
	if ref.Principal.PrincipalType != model.PrincipalTypeHuman {
		return validationFailed("/principal/principalType", "EligibilityRef requires exactly one Human PrincipalRef", apivalid.ViolationInvalidEnum)
	}
	return nil
}

// ValidateActionTargetBinding validates the closed ActionTargetBinding union (REQ-F18-06).
func ValidateActionTargetBinding(b model.ActionTargetBinding) *apiproblem.Problem {
	if !b.Mode.Valid() {
		return invalidEnum("/mode")
	}
	if prob := requireNonEmpty(b.Action, "/action"); prob != nil {
		return prob
	}
	switch b.Mode {
	case model.ActionTargetExactResource:
		if !b.IsExactResource() {
			return validationFailed("/targetKind", "ExactResource requires targetKind only", apivalid.ViolationInvalidEnum)
		}
	case model.ActionTargetCreateParent:
		if !b.IsCreateParent() || !b.ParentScopeKind.Valid() {
			return validationFailed("/parentScopeKind", "CreateParent requires a valid parentScopeKind only", apivalid.ViolationInvalidEnum)
		}
	case model.ActionTargetScopeOnly:
		if !b.IsScopeOnly() || !b.ScopeKind.Valid() {
			return validationFailed("/scopeKind", "ScopeOnly requires a valid scopeKind only", apivalid.ViolationInvalidEnum)
		}
	default:
		return invalidEnum("/mode")
	}
	return nil
}

// ValidateAssuranceEvidence validates operation-local assurance (REQ-F18-03).
func ValidateAssuranceEvidence(ev model.AssuranceEvidence, actor model.PrincipalRef) *apiproblem.Problem {
	if prob := ValidatePrincipalRef(ev.PrincipalRef); prob != nil {
		return prob
	}
	if !ev.PrincipalRef.Equal(actor) {
		return validationFailed("/principalRef", "AssuranceEvidence principalRef must equal authenticated actor", apivalid.ViolationInvalidEnum)
	}
	if ev.AuthenticatedAt.IsZero() {
		return missingField("/authenticatedAt")
	}
	if !ev.Level.Valid() {
		return invalidEnum("/level")
	}
	if prob := requireNonEmpty(ev.SourceAuthority, "/sourceAuthority"); prob != nil {
		return prob
	}
	if prob := requireNonEmpty(ev.IntegrityReference, "/integrityReference"); prob != nil {
		return prob
	}
	return nil
}
