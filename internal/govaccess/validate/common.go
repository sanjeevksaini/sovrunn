// Package validate provides deterministic FEATURE-0018 structural and local
// semantic validators (REQ-F18-21; design §5.1–§5.2; DD-03).
//
// Validators emit only inherited FEATURE-0012 Problem/violation codes and
// reuse apivalid.SafeDenial for inaccessible-resource protection.
package validate

import (
	"strings"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
)

// ContractKind identifies a FEATURE-0018 owned contract for scope registries.
type ContractKind string

const (
	ContractAccessGroup             ContractKind = "AccessGroup"
	ContractMembership              ContractKind = "Membership"
	ContractRoleDefinition          ContractKind = "RoleDefinition"
	ContractRoleAssignment          ContractKind = "RoleAssignment"
	ContractPrivilegedAccessRequest ContractKind = "PrivilegedAccessRequest"
	ContractAccessReview            ContractKind = "AccessReview"
	ContractApprovalPolicy          ContractKind = "ApprovalPolicy"
	ContractApprovalRequest         ContractKind = "ApprovalRequest"
	ContractExceptionGrant          ContractKind = "ExceptionGrant"
	ContractGovernanceProfile       ContractKind = "GovernanceProfile"
)

// AllowedScopeKinds returns the F18-RD-02 registered scope-kind set for kind.
func AllowedScopeKinds(kind ContractKind) ([]apimeta.ScopeKind, bool) {
	switch kind {
	case ContractAccessGroup, ContractMembership:
		return []apimeta.ScopeKind{apimeta.ScopeOrganization, apimeta.ScopeCloudProvider}, true
	case ContractRoleDefinition, ContractApprovalPolicy, ContractGovernanceProfile:
		return []apimeta.ScopeKind{apimeta.ScopePlatform, apimeta.ScopeOrganization, apimeta.ScopeCloudProvider}, true
	case ContractRoleAssignment, ContractPrivilegedAccessRequest, ContractAccessReview,
		ContractApprovalRequest, ContractExceptionGrant:
		return apimeta.AllScopeKinds(), true
	default:
		return nil, false
	}
}

func validationFailed(field, message string, code apiproblem.ViolationCode) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
		Field:   field,
		Code:    code,
		Message: message,
	}})
}

func missingField(field string) *apiproblem.Problem {
	return validationFailed(field, "required field is missing or empty", apivalid.ViolationInvalidEnum)
}

func invalidEnum(field string) *apiproblem.Problem {
	return validationFailed(field, "value is not a closed enum member", apivalid.ViolationInvalidEnum)
}

// ValidateScopeApplicability checks metadata.scopeRef against the registered
// per-contract ScopeKind set (REQ-F18-02; F18-RD-02).
func ValidateScopeApplicability(kind ContractKind, scope *apimeta.ScopeRef) *apiproblem.Problem {
	allowed, ok := AllowedScopeKinds(kind)
	if !ok {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithDetail("unknown contract kind")
	}
	normalized := apimeta.NormalizeScope(scope)
	effective := apimeta.ScopePlatform
	if normalized != nil {
		effective = apimeta.ScopeKind(normalized.Kind)
	}
	for _, a := range allowed {
		if a == effective {
			if effective == apimeta.ScopePlatform {
				return nil
			}
			if normalized == nil || normalized.UID == "" {
				return validationFailed("/metadata/scopeRef/uid", "non-platform scopeRef.uid is required", apivalid.ViolationScopeRefRequired)
			}
			if !effective.Valid() {
				return invalidEnum("/metadata/scopeRef/kind")
			}
			return nil
		}
	}
	return validationFailed("/metadata/scopeRef/kind", "scopeRef.kind is not permitted for this contract", apivalid.ViolationInvalidEnum)
}

// AccessGroupHolderScopeCompatibility validates RoleAssignment scope against an
// AccessGroup holder's owning scope (REQ-F18-04; F18-RD-02 reference compatibility).
//
// Organization-scoped groups may hold assignments at that Organization or its
// OrganizationUnit/Tenant/Project descendants. CloudProvider-scoped groups may
// hold only at that exact CloudProvider. Platform and CloudPlatform AccessGroup
// assignments are prohibited.
func AccessGroupHolderScopeCompatibility(
	groupScope *apimeta.ScopeRef,
	assignmentScope *apimeta.ScopeRef,
) *apiproblem.Problem {
	gNorm := apimeta.NormalizeScope(groupScope)
	aNorm := apimeta.NormalizeScope(assignmentScope)
	if gNorm == nil {
		return validationFailed("/metadata/scopeRef", "AccessGroup scopeRef is required", apivalid.ViolationScopeRefRequired)
	}
	gKind := apimeta.ScopeKind(gNorm.Kind)
	aKind := apimeta.ScopePlatform
	if aNorm != nil {
		aKind = apimeta.ScopeKind(aNorm.Kind)
	}
	switch gKind {
	case apimeta.ScopeOrganization:
		switch aKind {
		case apimeta.ScopeOrganization, apimeta.ScopeOrganizationUnit, apimeta.ScopeTenant, apimeta.ScopeProject:
			if aNorm == nil || aNorm.UID == "" {
				return validationFailed("/metadata/scopeRef/uid", "assignment scopeRef.uid is required", apivalid.ViolationScopeRefRequired)
			}
			if aKind == apimeta.ScopeOrganization && aNorm.UID != gNorm.UID {
				return validationFailed("/metadata/scopeRef/uid", "assignment Organization scope must match AccessGroup owning scope", apivalid.ViolationInvalidEnum)
			}
			return nil
		default:
			return validationFailed("/metadata/scopeRef/kind", "AccessGroup-held assignment scope is incompatible", apivalid.ViolationInvalidEnum)
		}
	case apimeta.ScopeCloudProvider:
		if aKind != apimeta.ScopeCloudProvider {
			return validationFailed("/metadata/scopeRef/kind", "CloudProvider AccessGroup may only hold CloudProvider assignments", apivalid.ViolationInvalidEnum)
		}
		if aNorm == nil || aNorm.UID == "" || aNorm.UID != gNorm.UID {
			return validationFailed("/metadata/scopeRef/uid", "assignment CloudProvider scope must equal AccessGroup owning scope", apivalid.ViolationScopeRefRequired)
		}
		return nil
	default:
		return validationFailed("/metadata/scopeRef/kind", "AccessGroup must be Organization or CloudProvider scoped", apivalid.ViolationInvalidEnum)
	}
}

// SafeDenyInaccessible returns the inherited safe-denial Problem for an
// inaccessible reference (REQ-F18-21; PRIV-01; design §5.2). Existence and
// sensitive policy detail are never disclosed.
func SafeDenyInaccessible() *apiproblem.Problem {
	return apivalid.SafeDenial(apivalid.DenyNotDisclosed)
}

// ResolveReferenceVisibility maps found+authorized into a safe outcome.
// Unauthorized or missing targets both yield RESOURCE_NOT_FOUND.
func ResolveReferenceVisibility(found, authorized bool) *apiproblem.Problem {
	if found && authorized {
		return nil
	}
	return SafeDenyInaccessible()
}

func requireNonEmpty(value, field string) *apiproblem.Problem {
	if strings.TrimSpace(value) == "" {
		return missingField(field)
	}
	return nil
}

func requireTypeMeta(tm apimeta.TypeMeta, apiVersion, kind, fieldPrefix string) *apiproblem.Problem {
	if tm.APIVersion != apiVersion {
		return validationFailed(fieldPrefix+"/apiVersion", "apiVersion mismatch", apivalid.ViolationInvalidAPIVersion)
	}
	if tm.Kind != kind {
		return validationFailed(fieldPrefix+"/kind", "kind mismatch", apivalid.ViolationInvalidKind)
	}
	return nil
}

func requireObjectName(meta apimeta.ObjectMeta) *apiproblem.Problem {
	return requireNonEmpty(meta.Name, "/metadata/name")
}

// RoleDefinitionKind is the RoleDefinition kind token.
const RoleDefinitionKind = "RoleDefinition"

// MembershipKind is the Membership kind token.
const MembershipKind = "Membership"
