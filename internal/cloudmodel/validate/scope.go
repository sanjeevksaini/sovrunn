package validate

import (
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

// AllowedScopeKinds returns the closed scope-kind subset for kind
// (DEC-0037; design §4.1 scopes).
func AllowedScopeKinds(kind Kind) ([]apimeta.ScopeKind, bool) {
	switch kind {
	case KindCloudPlatform, KindCloudProvider:
		return []apimeta.ScopeKind{apimeta.ScopePlatform}, true
	case KindCloudProviderParticipation:
		return []apimeta.ScopeKind{apimeta.ScopeCloudPlatform}, true
	case KindHostingLocation, KindDatacenter, KindFaultDomain, KindInfrastructureStack:
		return []apimeta.ScopeKind{apimeta.ScopeCloudProvider}, true
	default:
		return nil, false
	}
}

// ValidateScopeKindSubset checks that scopeRef is within the per-kind
// seven-scope subset. Canonical Platform scope is represented as a nil
// ScopeRef (or explicit Kind==Platform before NormalizeScope).
func ValidateScopeKindSubset(kind Kind, scope *apimeta.ScopeRef) *apiproblem.Problem {
	allowed, ok := AllowedScopeKinds(kind)
	if !ok {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithDetail("unknown resource kind")
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
				return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
					Field:   "/metadata/scopeRef/uid",
					Code:    ViolationScopeKindInvalid,
					Message: "non-platform scopeRef.uid is required",
				}})
			}
			if !effective.Valid() {
				return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
					Field:   "/metadata/scopeRef/kind",
					Code:    ViolationScopeKindInvalid,
					Message: "scopeRef.kind is not a canonical ScopeKind",
				}})
			}
			return nil
		}
	}
	return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
		Field:   "/metadata/scopeRef/kind",
		Code:    ViolationScopeKindInvalid,
		Message: "scopeRef.kind is not permitted for this resource kind",
	}})
}

// ValidateParticipationScopeUIDInvariant enforces
// spec.cloudPlatformRef.uid == metadata.scopeRef.uid (REQ-F15-14; AC-F15-12).
func ValidateParticipationScopeUIDInvariant(p model.CloudProviderParticipation) *apiproblem.Problem {
	scope := apimeta.NormalizeScope(p.Metadata.ScopeRef)
	if scope == nil || apimeta.ScopeKind(scope.Kind) != apimeta.ScopeCloudPlatform {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/metadata/scopeRef",
			Code:    ViolationScopeKindInvalid,
			Message: "participation scopeRef.kind must be CloudPlatform",
		}})
	}
	if p.Spec.CloudPlatformRef.UID == "" || scope.UID == "" || p.Spec.CloudPlatformRef.UID != scope.UID {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/cloudPlatformRef/uid",
			Code:    ViolationScopeReferenceMismatch,
			Message: "cloudPlatformRef.uid must equal metadata.scopeRef.uid",
		}})
	}
	return nil
}
