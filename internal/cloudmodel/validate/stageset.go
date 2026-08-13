package validate

import (
	"context"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

// RequestClass selects which FEATURE-0015 validator composition to bind into
// the inherited FEATURE-0012 StageSet (design DD-02; §5.1).
type RequestClass string

const (
	RequestClassCollectionCreate RequestClass = "collection-create"
	RequestClassPatch            RequestClass = "patch"
	RequestClassItemAction       RequestClass = "item-action"
)

// StageSetFor returns a StageSet that reuses FEATURE-0012 DefaultingStage /
// ValidationStage / StageSet contracts with FEATURE-0015 validator adapters.
// It does not invoke HTTP handlers and does not fork or bypass apivalid.
func StageSetFor(kind Kind, class RequestClass) (apivalid.StageSet, bool) {
	if _, ok := CreateClassification(kind); !ok {
		return apivalid.StageSet{}, false
	}
	switch class {
	case RequestClassCollectionCreate:
		return apivalid.StageSet{
			Defaulting: apivalid.NewCommonDefaulting(),
			Semantic:   &semanticAdapter{kind: kind, class: class},
			Reference:  &referenceAdapter{kind: kind, class: class},
		}, true
	case RequestClassPatch:
		if kind == KindCloudProviderParticipation {
			return apivalid.StageSet{}, false
		}
		return apivalid.StageSet{
			Defaulting: apivalid.NewCommonDefaulting(),
			Semantic:   &semanticAdapter{kind: kind, class: class},
			Reference:  &referenceAdapter{kind: kind, class: class},
		}, true
	case RequestClassItemAction:
		if kind != KindCloudProviderParticipation {
			return apivalid.StageSet{}, false
		}
		return apivalid.StageSet{
			Defaulting: apivalid.NewCommonDefaulting(),
			Semantic:   &semanticAdapter{kind: kind, class: class},
			Reference:  &referenceAdapter{kind: kind, class: class},
		}, true
	default:
		return apivalid.StageSet{}, false
	}
}

// semanticAdapter plugs FEATURE-0015 schema/ISO validators into layer 6.
type semanticAdapter struct {
	kind  Kind
	class RequestClass
}

var _ apivalid.ValidationStage = (*semanticAdapter)(nil)

func (s *semanticAdapter) Validate(_ context.Context, object any) ([]apiproblem.Violation, error) {
	if s == nil {
		return nil, apivalid.ErrSemanticInternal
	}
	if s.class == RequestClassItemAction {
		// Empty-body item actions: explicit no-op still invoked by StageSet.
		return nil, nil
	}
	if prob := ValidateSchema(s.kind, object); prob != nil {
		if prob.Code == apiproblem.CodeInternalError {
			return nil, apivalid.ErrSemanticInternal
		}
		return append([]apiproblem.Violation(nil), prob.Violations...), nil
	}
	return nil, nil
}

// referenceAdapter plugs FEATURE-0015 scope/reference validators into layer 7.
type referenceAdapter struct {
	kind  Kind
	class RequestClass
}

var _ apivalid.ValidationStage = (*referenceAdapter)(nil)

func (r *referenceAdapter) Validate(_ context.Context, object any) ([]apiproblem.Violation, error) {
	if r == nil {
		return nil, apivalid.ErrSemanticInternal
	}
	if r.class == RequestClassItemAction {
		return nil, nil
	}
	prob := validateReferenceStage(r.kind, object)
	if prob == nil {
		return nil, nil
	}
	if prob.Code == apiproblem.CodeInternalError {
		return nil, apivalid.ErrSemanticInternal
	}
	return append([]apiproblem.Violation(nil), prob.Violations...), nil
}

func validateReferenceStage(kind Kind, object any) *apiproblem.Problem {
	switch kind {
	case KindCloudPlatform:
		cp, ok := object.(model.CloudPlatform)
		if !ok {
			return typeMismatchProblem()
		}
		return ValidateScopeKindSubset(kind, cp.Metadata.ScopeRef)
	case KindCloudProvider:
		cp, ok := object.(model.CloudProvider)
		if !ok {
			return typeMismatchProblem()
		}
		return ValidateScopeKindSubset(kind, cp.Metadata.ScopeRef)
	case KindCloudProviderParticipation:
		p, ok := object.(model.CloudProviderParticipation)
		if !ok {
			return typeMismatchProblem()
		}
		if prob := ValidateParticipationRefs(p); prob != nil {
			return prob
		}
		if prob := ValidateScopeKindSubset(kind, p.Metadata.ScopeRef); prob != nil {
			return prob
		}
		if apimeta.NormalizeScope(p.Metadata.ScopeRef) != nil {
			return ValidateParticipationScopeUIDInvariant(p)
		}
		return nil
	case KindHostingLocation:
		hl, ok := object.(model.HostingLocation)
		if !ok {
			return typeMismatchProblem()
		}
		return ValidateScopeKindSubset(kind, hl.Metadata.ScopeRef)
	case KindDatacenter:
		dc, ok := object.(model.Datacenter)
		if !ok {
			return typeMismatchProblem()
		}
		if prob := ValidateTypedRefStructural(dc.Spec.HostingLocationRef, "/spec/hostingLocationRef", model.KindHostingLocation, model.APIVersionHostingLocation); prob != nil {
			return prob
		}
		return ValidateScopeKindSubset(kind, dc.Metadata.ScopeRef)
	case KindFaultDomain:
		fd, ok := object.(model.FaultDomain)
		if !ok {
			return typeMismatchProblem()
		}
		if prob := ValidateTypedRefStructural(fd.Spec.DatacenterRef, "/spec/datacenterRef", model.KindDatacenter, model.APIVersionDatacenter); prob != nil {
			return prob
		}
		return ValidateScopeKindSubset(kind, fd.Metadata.ScopeRef)
	case KindInfrastructureStack:
		st, ok := object.(model.InfrastructureStack)
		if !ok {
			return typeMismatchProblem()
		}
		if prob := ValidateTypedRefStructural(st.Spec.FaultDomainRef, "/spec/faultDomainRef", model.KindFaultDomain, model.APIVersionFaultDomain); prob != nil {
			return prob
		}
		return ValidateScopeKindSubset(kind, st.Metadata.ScopeRef)
	default:
		return apiproblem.New(apiproblem.CodeValidationFailed).WithDetail("unknown resource kind")
	}
}
