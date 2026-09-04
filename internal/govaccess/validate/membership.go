package validate

import (
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

// ValidateMembership validates Membership structural/local semantic rules
// (VS0-SCHEMA-021; REQ-F18-05).
func ValidateMembership(m model.Membership) *apiproblem.Problem {
	if prob := requireTypeMeta(m.TypeMeta, model.APIVersionIAM, MembershipKind, ""); prob != nil {
		return prob
	}
	if prob := requireObjectName(m.Metadata); prob != nil {
		return prob
	}
	if prob := ValidateScopeApplicability(ContractMembership, m.Metadata.ScopeRef); prob != nil {
		return prob
	}
	if prob := ValidatePrincipalRef(m.Spec.PrincipalRef); prob != nil {
		return prob
	}
	if !m.Spec.ContextKind.Valid() {
		return invalidEnum("/spec/contextKind")
	}
	if !m.Spec.MembershipType.Valid() {
		return invalidEnum("/spec/membershipType")
	}
	switch m.Spec.ContextKind {
	case model.MembershipContextAccessGroup:
		if m.Spec.AccessGroupRef == nil {
			return missingField("/spec/accessGroupRef")
		}
		if prob := ValidateAccessGroupRef(*m.Spec.AccessGroupRef); prob != nil {
			return prob
		}
	default:
		if m.Spec.AccessGroupRef != nil {
			return validationFailed("/spec/accessGroupRef", "accessGroupRef is permitted only for AccessGroup context", apivalid.ViolationInvalidEnum)
		}
	}
	if m.Spec.MembershipType == model.MembershipTypeGuest {
		if m.Spec.PrincipalRef.PrincipalType != model.PrincipalTypeHuman {
			return validationFailed("/spec/membershipType", "Guest membership requires a Human principal", apivalid.ViolationInvalidEnum)
		}
		if m.Spec.SponsorRef == nil {
			return missingField("/spec/sponsorRef")
		}
		if prob := ValidatePrincipalRef(*m.Spec.SponsorRef); prob != nil {
			return prob
		}
		if m.Spec.SponsorRef.PrincipalType != model.PrincipalTypeHuman {
			return validationFailed("/spec/sponsorRef/principalType", "Guest sponsor must be Human", apivalid.ViolationInvalidEnum)
		}
		if m.Spec.ExpiresAt == nil || m.Spec.ExpiresAt.IsZero() {
			return missingField("/spec/expiresAt")
		}
	}
	switch m.Spec.PrincipalRef.PrincipalType {
	case model.PrincipalTypeWorkload, model.PrincipalTypeSystem:
		if m.Spec.MembershipType != model.MembershipTypeStandard {
			return validationFailed("/spec/membershipType", "Workload/System membership must be Standard", apivalid.ViolationInvalidEnum)
		}
		if m.Spec.ResponsiblePartyRef == nil {
			return missingField("/spec/responsiblePartyRef")
		}
		if prob := ValidatePrincipalRef(*m.Spec.ResponsiblePartyRef); prob != nil {
			return prob
		}
		if m.Spec.ResponsiblePartyRef.PrincipalType != model.PrincipalTypeHuman {
			return validationFailed("/spec/responsiblePartyRef/principalType", "responsiblePartyRef must be Human", apivalid.ViolationInvalidEnum)
		}
	}
	if prob := requireNonEmpty(m.Protected.SourceAuthority, "/protected/sourceAuthority"); prob != nil {
		return prob
	}
	if m.Protected.ProvisionedAt.IsZero() {
		return missingField("/protected/provisionedAt")
	}
	if m.Status.Phase != "" && !m.Status.Phase.Valid() {
		return invalidEnum("/status/phase")
	}
	return nil
}
