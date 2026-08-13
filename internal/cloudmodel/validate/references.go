package validate

import (
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apiref"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

// TopologyParents carries resolved parent identities for containment checks.
// Callers supply already-authorized, resolved parents; this package performs
// no store I/O (design §6.1.1 import graph).
type TopologyParents struct {
	HostingLocation  *model.HostingLocation
	Datacenter       *model.Datacenter
	FaultDomain      *model.FaultDomain
	CloudProviderUID string
}

// ValidateTypedRefStructural checks apiVersion/kind/name/uid-pinned shape and
// allowed kind without resolving existence.
func ValidateTypedRefStructural(ref apimeta.TypedRef, field string, allowedKind, allowedAPIVersion string) *apiproblem.Problem {
	c := apiref.Constraint{AllowedKinds: []string{allowedKind}}
	issues := c.ValidateRef(apiref.TypedRef(ref), field)
	if len(issues) > 0 {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   issues[0].Path,
			Code:    apiproblem.ViolationCode(issues[0].Code),
			Message: issues[0].Message,
		}})
	}
	if ref.UID == "" || !apimeta.IsGeneratedUIDFormat(ref.UID) {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   field + "/uid",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "reference uid must be present and UID-pinned",
		}})
	}
	if allowedAPIVersion != "" && ref.APIVersion != allowedAPIVersion {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   field + "/apiVersion",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "reference apiVersion is not permitted",
		}})
	}
	return nil
}

// ValidateDatacenterParentChain checks hostingLocationRef containment and
// same-provider UID preservation against the resolved HostingLocation.
func ValidateDatacenterParentChain(dc model.Datacenter, parent model.HostingLocation) *apiproblem.Problem {
	if prob := ValidateTypedRefStructural(dc.Spec.HostingLocationRef, "/spec/hostingLocationRef", model.KindHostingLocation, model.APIVersionHostingLocation); prob != nil {
		return prob
	}
	if dc.Spec.HostingLocationRef.UID != parent.Metadata.UID {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/hostingLocationRef/uid",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "hostingLocationRef does not resolve to the supplied parent",
		}})
	}
	return validateSameProviderScope(dc.Metadata.ScopeRef, parent.Metadata.ScopeRef, "/metadata/scopeRef")
}

// ValidateFaultDomainParentChain checks datacenterRef containment and
// same-provider UID preservation.
func ValidateFaultDomainParentChain(fd model.FaultDomain, parent model.Datacenter) *apiproblem.Problem {
	if prob := ValidateTypedRefStructural(fd.Spec.DatacenterRef, "/spec/datacenterRef", model.KindDatacenter, model.APIVersionDatacenter); prob != nil {
		return prob
	}
	if fd.Spec.DatacenterRef.UID != parent.Metadata.UID {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/datacenterRef/uid",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "datacenterRef does not resolve to the supplied parent",
		}})
	}
	return validateSameProviderScope(fd.Metadata.ScopeRef, parent.Metadata.ScopeRef, "/metadata/scopeRef")
}

// ValidateInfrastructureStackParentChain checks faultDomainRef containment and
// same-provider UID preservation.
func ValidateInfrastructureStackParentChain(stack model.InfrastructureStack, parent model.FaultDomain) *apiproblem.Problem {
	if prob := ValidateTypedRefStructural(stack.Spec.FaultDomainRef, "/spec/faultDomainRef", model.KindFaultDomain, model.APIVersionFaultDomain); prob != nil {
		return prob
	}
	if stack.Spec.FaultDomainRef.UID != parent.Metadata.UID {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   "/spec/faultDomainRef/uid",
			Code:    apiproblem.ViolationOutOfRange,
			Message: "faultDomainRef does not resolve to the supplied parent",
		}})
	}
	return validateSameProviderScope(stack.Metadata.ScopeRef, parent.Metadata.ScopeRef, "/metadata/scopeRef")
}

// ValidateParticipationRefs checks UID-pinned CloudPlatform and CloudProvider
// references on participation create.
func ValidateParticipationRefs(p model.CloudProviderParticipation) *apiproblem.Problem {
	if prob := ValidateTypedRefStructural(p.Spec.CloudPlatformRef, "/spec/cloudPlatformRef", model.KindCloudPlatform, model.APIVersionCloudPlatform); prob != nil {
		return prob
	}
	if prob := ValidateTypedRefStructural(p.Spec.CloudProviderRef, "/spec/cloudProviderRef", model.KindCloudProvider, model.APIVersionCloudProvider); prob != nil {
		return prob
	}
	return nil
}

func validateSameProviderScope(child, parent *apimeta.ScopeRef, field string) *apiproblem.Problem {
	cn := apimeta.NormalizeScope(child)
	pn := apimeta.NormalizeScope(parent)
	if cn == nil || pn == nil {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    ViolationTopologyProviderMismatch,
			Message: "topology resources require CloudProvider scope",
		}})
	}
	if apimeta.ScopeKind(cn.Kind) != apimeta.ScopeCloudProvider || apimeta.ScopeKind(pn.Kind) != apimeta.ScopeCloudProvider {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    ViolationTopologyProviderMismatch,
			Message: "topology resources require CloudProvider scope",
		}})
	}
	if cn.UID == "" || pn.UID == "" || cn.UID != pn.UID {
		return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    ViolationTopologyProviderMismatch,
			Message: "child and parent must share the same CloudProvider scope UID",
		}})
	}
	return nil
}
