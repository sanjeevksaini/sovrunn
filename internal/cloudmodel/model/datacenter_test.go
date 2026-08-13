package model

import (
	"encoding/json"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func TestDatacenterJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := Datacenter{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: APIVersionDatacenter,
			Kind:       KindDatacenter,
		},
		Metadata: apimeta.ObjectMeta{Name: "dc-1"},
		Spec: DatacenterSpec{
			HostingLocationRef: apimeta.TypedRef{
				APIVersion: APIVersionHostingLocation,
				Kind:       KindHostingLocation,
				Name:       "mumbai",
				UID:        "dddddddddddddddddddddddddddddddd",
			},
			Description: "primary dc",
		},
		Status: DatacenterStatus{Phase: DatacenterPhaseActive},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out Datacenter
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Spec.HostingLocationRef.Kind != KindHostingLocation || out.Spec.HostingLocationRef.UID == "" {
		t.Fatalf("immutable parent ref lost: %+v", out.Spec.HostingLocationRef)
	}
}
