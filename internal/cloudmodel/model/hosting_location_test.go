package model

import (
	"encoding/json"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func TestHostingLocationJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := HostingLocation{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: APIVersionHostingLocation,
			Kind:       KindHostingLocation,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "mumbai",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: APIVersionCloudProvider,
				Kind:       string(apimeta.ScopeCloudProvider),
				Name:       "provider-a",
				UID:        "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			}},
		},
		Spec: HostingLocationSpec{
			CountryCode:            "IN",
			Locality:               "Mumbai",
			AdministrativeAreaCode: "IN-MH",
			Description:            "west india",
		},
		Status: HostingLocationStatus{Phase: HostingLocationPhaseActive},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var top map[string]any
	if err := json.Unmarshal(raw, &top); err != nil {
		t.Fatalf("unmarshal top: %v", err)
	}
	spec := top["spec"].(map[string]any)
	for _, forbidden := range []string{"parentRef", "hostingLocationRef", "datacenterRef", "faultDomainRef", "providerRef"} {
		if _, ok := spec[forbidden]; ok {
			t.Fatalf("HostingLocation must have no parent reference field %q", forbidden)
		}
	}

	var out HostingLocation
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Spec.CountryCode != "IN" || out.Spec.Locality != "Mumbai" {
		t.Fatalf("spec round-trip failed: %+v", out.Spec)
	}
}
