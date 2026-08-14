package model

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func TestCloudPlatformJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := CloudPlatform{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: APIVersionCloudPlatform,
			Kind:       KindCloudPlatform,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "acme-platform",
			UID:  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: "core.sovrunn.io/v1alpha1",
				Kind:       string(apimeta.ScopePlatform),
				Name:       "platform",
				UID:        apimeta.PlatformScopeUID,
			}},
			Generation:      1,
			ResourceVersion: "1",
			CreatedAt:       "2026-08-12T00:00:00Z",
			UpdatedAt:       "2026-08-12T00:00:00Z",
		},
		Spec: CloudPlatformSpec{
			OwnerRegistration: OwnerRegistration{
				LegalName:              "Acme Cloud Ltd",
				RegistrationIdentifier: "REG-001",
				JurisdictionCode:       "US",
			},
			Description: "primary platform",
		},
		Status: CloudPlatformStatus{Phase: CloudPlatformPhaseActive},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(raw, &top); err != nil {
		t.Fatalf("unmarshal top: %v", err)
	}
	for _, key := range []string{"apiVersion", "kind", "metadata", "spec", "status"} {
		if _, ok := top[key]; !ok {
			t.Fatalf("missing top-level key %q in %s", key, string(raw))
		}
	}
	if strings.Contains(string(raw), `"resources."`) || strings.Contains(string(raw), "ProviderSpec") {
		t.Fatalf("must not reuse alpha Provider types: %s", string(raw))
	}

	var out CloudPlatform
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Kind != KindCloudPlatform || out.APIVersion != APIVersionCloudPlatform {
		t.Fatalf("TypeMeta mismatch: %+v", out.TypeMeta)
	}
	if out.Spec.OwnerRegistration.JurisdictionCode != "US" {
		t.Fatalf("ownerRegistration round-trip failed: %+v", out.Spec.OwnerRegistration)
	}
	if out.Status.Phase != CloudPlatformPhaseActive {
		t.Fatalf("status.phase=%q, want Active", out.Status.Phase)
	}
	if out.Metadata.ScopeRef == nil || out.Metadata.ScopeRef.Kind != string(apimeta.ScopePlatform) {
		t.Fatalf("metadata.scopeRef not preserved: %+v", out.Metadata.ScopeRef)
	}
}

func TestCloudPlatformMetadataSpecStatusSeparation(t *testing.T) {
	t.Parallel()

	raw := []byte(`{
		"apiVersion":"core.sovrunn.io/v1alpha1",
		"kind":"CloudPlatform",
		"metadata":{"name":"p1"},
		"spec":{"ownerRegistration":{"legalName":"L","registrationIdentifier":"R","jurisdictionCode":"IN"}},
		"status":{"phase":"Active"}
	}`)
	var cp CloudPlatform
	if err := json.Unmarshal(raw, &cp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if cp.Metadata.Name != "p1" {
		t.Fatalf("metadata.name=%q", cp.Metadata.Name)
	}
	if cp.Spec.OwnerRegistration.LegalName != "L" {
		t.Fatalf("spec leaked into metadata or lost: %+v", cp.Spec)
	}
	if cp.Status.Phase != CloudPlatformPhaseActive {
		t.Fatalf("status.phase=%q", cp.Status.Phase)
	}
}
