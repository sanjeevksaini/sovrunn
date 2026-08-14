package model

import (
	"encoding/json"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func TestFaultDomainJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := FaultDomain{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: APIVersionFaultDomain,
			Kind:       KindFaultDomain,
		},
		Metadata: apimeta.ObjectMeta{Name: "fd-1"},
		Spec: FaultDomainSpec{
			DatacenterRef: apimeta.TypedRef{
				APIVersion: APIVersionDatacenter,
				Kind:       KindDatacenter,
				Name:       "dc-1",
				UID:        "eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
			},
		},
		Status: FaultDomainStatus{Phase: FaultDomainPhaseActive},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out FaultDomain
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Spec.DatacenterRef.Kind != KindDatacenter || out.Spec.DatacenterRef.UID == "" {
		t.Fatalf("immutable parent ref lost: %+v", out.Spec.DatacenterRef)
	}
}
