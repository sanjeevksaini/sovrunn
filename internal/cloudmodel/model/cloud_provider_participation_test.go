package model

import (
	"encoding/json"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
)

func TestCloudProviderParticipationJSONRoundTrip(t *testing.T) {
	t.Parallel()

	in := CloudProviderParticipation{
		TypeMeta: apimeta.TypeMeta{
			APIVersion: APIVersionCloudProviderParticipation,
			Kind:       KindCloudProviderParticipation,
		},
		Metadata: apimeta.ObjectMeta{
			Name: "participation-1",
			UID:  "cccccccccccccccccccccccccccccccc",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: APIVersionCloudPlatform,
				Kind:       string(apimeta.ScopeCloudPlatform),
				Name:       "acme-platform",
				UID:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			}},
		},
		Spec: CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{
				APIVersion: APIVersionCloudPlatform,
				Kind:       KindCloudPlatform,
				Name:       "acme-platform",
				UID:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			},
			CloudProviderRef: apimeta.TypedRef{
				APIVersion: APIVersionCloudProvider,
				Kind:       KindCloudProvider,
				Name:       "provider-a",
				UID:        "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			},
			Environment: ParticipationEnvironmentDevelopment,
		},
		Status: CloudProviderParticipationStatus{
			Phase:             ParticipationPhasePending,
			PlatformSuspended: false,
			ProviderSuspended: false,
			RequestExpiresAt:  "2026-08-19T00:00:00Z",
		},
	}

	raw, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out CloudProviderParticipation
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Spec.Environment != ParticipationEnvironmentDevelopment {
		t.Fatalf("environment=%q", out.Spec.Environment)
	}
	if out.Spec.CloudPlatformRef.UID == "" || out.Spec.CloudProviderRef.UID == "" {
		t.Fatalf("UID-pinned refs lost: %+v", out.Spec)
	}
	if out.Status.Phase != ParticipationPhasePending {
		t.Fatalf("phase=%q", out.Status.Phase)
	}
	if out.Status.PlatformSuspended || out.Status.ProviderSuspended {
		t.Fatalf("holds must round-trip false: %+v", out.Status)
	}
	if out.Status.RequestExpiresAt != "2026-08-19T00:00:00Z" {
		t.Fatalf("requestExpiresAt=%q", out.Status.RequestExpiresAt)
	}
	if _, ok := jsonField(raw, "spec", "providerSelectionModes"); ok {
		t.Fatalf("FEATURE-0021 field must not be present: %s", string(raw))
	}
	if _, ok := jsonField(raw, "spec", "permittedHostingLocationRefs"); ok {
		t.Fatalf("FEATURE-0021 field must not be present: %s", string(raw))
	}
}

func jsonField(raw []byte, path ...string) (any, bool) {
	var cur any
	if err := json.Unmarshal(raw, &cur); err != nil {
		return nil, false
	}
	for _, key := range path {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		cur, ok = m[key]
		if !ok {
			return nil, false
		}
	}
	return cur, true
}
