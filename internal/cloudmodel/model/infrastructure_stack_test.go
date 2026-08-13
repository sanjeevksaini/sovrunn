package model

import (
	"encoding/json"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func TestInfrastructureStackJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := InfrastructureStack{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: APIVersionInfrastructureStack,
			Kind:       KindInfrastructureStack,
		},
		Metadata: apimeta.ObjectMeta{Name: "stack-1"},
		Spec: InfrastructureStackSpec{
			FaultDomainRef: apimeta.TypedRef{
				APIVersion: APIVersionFaultDomain,
				Kind:       KindFaultDomain,
				Name:       "fd-1",
				UID:        "ffffffffffffffffffffffffffffffff",
			},
			Description: "k8s stack",
		},
		Status: InfrastructureStackStatus{Phase: InfrastructureStackPhaseActive},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out InfrastructureStack
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Kind != KindInfrastructureStack {
		t.Fatalf("kind=%q", out.Kind)
	}
	if out.Spec.FaultDomainRef.Kind != KindFaultDomain || out.Spec.FaultDomainRef.UID == "" {
		t.Fatalf("immutable parent ref lost: %+v", out.Spec.FaultDomainRef)
	}
}
