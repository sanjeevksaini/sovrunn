package model

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func TestCloudProviderJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := CloudProvider{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: APIVersionCloudProvider,
			Kind:       KindCloudProvider,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "provider-a",
			UID:  "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		},
		Spec: CloudProviderSpec{
			DisplayName:      "Provider A",
			OperatingMarkets: []string{"US", "IN"},
		},
		Status: CloudProviderStatus{Phase: CloudProviderPhaseActive},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out CloudProvider
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Kind != KindCloudProvider {
		t.Fatalf("kind=%q, want CloudProvider", out.Kind)
	}
	if !reflect.DeepEqual(out.Spec.OperatingMarkets, []string{"US", "IN"}) {
		t.Fatalf("operatingMarkets=%v", out.Spec.OperatingMarkets)
	}
	if out.Status.Phase != CloudProviderPhaseActive {
		t.Fatalf("status.phase=%q", out.Status.Phase)
	}
}

func TestCloudProviderNoAlphaImport(t *testing.T) {
	t.Parallel()

	// Compile-time and runtime shape check: Spec is CloudProviderSpec, not
	// an alpha Provider type alias.
	var p CloudProvider
	p.Spec.OperatingMarkets = []string{"DE"}
	raw, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var top map[string]any
	if err := json.Unmarshal(raw, &top); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	spec, ok := top["spec"].(map[string]any)
	if !ok {
		t.Fatalf("spec missing: %s", string(raw))
	}
	if _, ok := spec["operatingMarkets"]; !ok {
		t.Fatalf("operatingMarkets missing from spec: %s", string(raw))
	}
}
