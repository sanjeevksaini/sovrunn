package validate_test

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func providerScope(uid string) *apimeta.ScopeRef {
	return &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: model.APIVersionCloudProvider,
		Kind:       string(apimeta.ScopeCloudProvider),
		Name:       "prov",
		UID:        uid,
	}}
}

func TestValidateDatacenterParentChain_SameProvider(t *testing.T) {
	t.Parallel()
	providerUID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	locUID := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	loc := model.HostingLocation{
		Metadata: apimeta.ObjectMeta{Name: "loc", UID: locUID, ScopeRef: providerScope(providerUID)},
	}
	dc := model.Datacenter{
		Metadata: apimeta.ObjectMeta{Name: "dc", ScopeRef: providerScope(providerUID)},
		Spec: model.DatacenterSpec{
			HostingLocationRef: apimeta.TypedRef{
				APIVersion: model.APIVersionHostingLocation,
				Kind:       model.KindHostingLocation,
				Name:       "loc",
				UID:        locUID,
			},
		},
	}
	if prob := validate.ValidateDatacenterParentChain(dc, loc); prob != nil {
		t.Fatalf("same provider: %#v", prob)
	}

	dc.Metadata.ScopeRef = providerScope("cccccccccccccccccccccccccccccccc")
	prob := validate.ValidateDatacenterParentChain(dc, loc)
	if prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("cross-provider: %#v", prob)
	}
	if len(prob.Violations) == 0 || prob.Violations[0].Code != validate.ViolationTopologyProviderMismatch {
		t.Fatalf("want VS0_TOPOLOGY_PROVIDER_MISMATCH, got %#v", prob.Violations)
	}
}

func TestValidateFaultDomainAndStackParentChain(t *testing.T) {
	t.Parallel()
	providerUID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	dcUID := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	fdUID := "cccccccccccccccccccccccccccccccc"
	dc := model.Datacenter{
		Metadata: apimeta.ObjectMeta{Name: "dc", UID: dcUID, ScopeRef: providerScope(providerUID)},
	}
	fd := model.FaultDomain{
		Metadata: apimeta.ObjectMeta{Name: "fd", UID: fdUID, ScopeRef: providerScope(providerUID)},
		Spec: model.FaultDomainSpec{
			DatacenterRef: apimeta.TypedRef{
				APIVersion: model.APIVersionDatacenter,
				Kind:       model.KindDatacenter,
				Name:       "dc",
				UID:        dcUID,
			},
		},
	}
	if prob := validate.ValidateFaultDomainParentChain(fd, dc); prob != nil {
		t.Fatalf("fd chain: %#v", prob)
	}
	stack := model.InfrastructureStack{
		Metadata: apimeta.ObjectMeta{Name: "stack", ScopeRef: providerScope(providerUID)},
		Spec: model.InfrastructureStackSpec{
			FaultDomainRef: apimeta.TypedRef{
				APIVersion: model.APIVersionFaultDomain,
				Kind:       model.KindFaultDomain,
				Name:       "fd",
				UID:        fdUID,
			},
		},
	}
	if prob := validate.ValidateInfrastructureStackParentChain(stack, fd); prob != nil {
		t.Fatalf("stack chain: %#v", prob)
	}
}
