package validate

import (
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

// ValidateImmutablePatch compares before/after resource values and rejects
// immutable identity, reference, and owner-registration changes. For create,
// use ClassifyCreateContract (ADH-2026-047 decision 3).
func ValidateImmutablePatch(kind Kind, before, after any) *apiproblem.Problem {
	switch kind {
	case KindCloudPlatform:
		b, bok := before.(model.CloudPlatform)
		a, aok := after.(model.CloudPlatform)
		if !bok || !aok {
			return internalTypeMismatch()
		}
		if b.Metadata.Name != a.Metadata.Name {
			return immutableField("/metadata/name")
		}
		if !sameScope(b.Metadata.ScopeRef, a.Metadata.ScopeRef) {
			return immutableField("/metadata/scopeRef")
		}
		if b.Spec.OwnerRegistration != a.Spec.OwnerRegistration {
			return immutableField("/spec/ownerRegistration")
		}
	case KindCloudProvider:
		b, bok := before.(model.CloudProvider)
		a, aok := after.(model.CloudProvider)
		if !bok || !aok {
			return internalTypeMismatch()
		}
		if b.Metadata.Name != a.Metadata.Name {
			return immutableField("/metadata/name")
		}
		if !sameScope(b.Metadata.ScopeRef, a.Metadata.ScopeRef) {
			return immutableField("/metadata/scopeRef")
		}
	case KindHostingLocation:
		b, bok := before.(model.HostingLocation)
		a, aok := after.(model.HostingLocation)
		if !bok || !aok {
			return internalTypeMismatch()
		}
		if b.Metadata.Name != a.Metadata.Name {
			return immutableField("/metadata/name")
		}
		if !sameScope(b.Metadata.ScopeRef, a.Metadata.ScopeRef) {
			return immutableField("/metadata/scopeRef")
		}
		if b.Spec.CountryCode != a.Spec.CountryCode || b.Spec.Locality != a.Spec.Locality ||
			b.Spec.AdministrativeAreaCode != a.Spec.AdministrativeAreaCode {
			return immutableField("/spec")
		}
	case KindDatacenter:
		b, bok := before.(model.Datacenter)
		a, aok := after.(model.Datacenter)
		if !bok || !aok {
			return internalTypeMismatch()
		}
		if b.Metadata.Name != a.Metadata.Name {
			return immutableField("/metadata/name")
		}
		if !sameScope(b.Metadata.ScopeRef, a.Metadata.ScopeRef) {
			return immutableField("/metadata/scopeRef")
		}
		if b.Spec.HostingLocationRef != a.Spec.HostingLocationRef {
			return immutableField("/spec/hostingLocationRef")
		}
	case KindFaultDomain:
		b, bok := before.(model.FaultDomain)
		a, aok := after.(model.FaultDomain)
		if !bok || !aok {
			return internalTypeMismatch()
		}
		if b.Metadata.Name != a.Metadata.Name {
			return immutableField("/metadata/name")
		}
		if !sameScope(b.Metadata.ScopeRef, a.Metadata.ScopeRef) {
			return immutableField("/metadata/scopeRef")
		}
		if b.Spec.DatacenterRef != a.Spec.DatacenterRef {
			return immutableField("/spec/datacenterRef")
		}
	case KindInfrastructureStack:
		b, bok := before.(model.InfrastructureStack)
		a, aok := after.(model.InfrastructureStack)
		if !bok || !aok {
			return internalTypeMismatch()
		}
		if b.Metadata.Name != a.Metadata.Name {
			return immutableField("/metadata/name")
		}
		if !sameScope(b.Metadata.ScopeRef, a.Metadata.ScopeRef) {
			return immutableField("/metadata/scopeRef")
		}
		if b.Spec.FaultDomainRef != a.Spec.FaultDomainRef {
			return immutableField("/spec/faultDomainRef")
		}
	default:
		return apiproblem.New(apiproblem.CodeValidationFailed).WithDetail("unknown resource kind")
	}
	return nil
}

func immutableField(field string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).WithViolations([]apiproblem.Violation{{
		Field:   field,
		Code:    ViolationPatchImmutableField,
		Message: "immutable field cannot be changed",
	}})
}

func internalTypeMismatch() *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeInternalError).WithDetail("immutable-field type mismatch")
}

func sameScope(a, b *apimeta.ScopeRef) bool {
	na := apimeta.NormalizeScope(a)
	nb := apimeta.NormalizeScope(b)
	if na == nil && nb == nil {
		return true
	}
	if na == nil || nb == nil {
		return false
	}
	return na.APIVersion == nb.APIVersion &&
		na.Kind == nb.Kind &&
		na.Name == nb.Name &&
		na.UID == nb.UID
}
