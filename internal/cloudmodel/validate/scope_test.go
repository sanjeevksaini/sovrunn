package validate_test

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func TestValidateScopeKindSubset_SevenScopePerKind(t *testing.T) {
	t.Parallel()
	cases := []struct {
		kind    validate.Kind
		allowed apimeta.ScopeKind
	}{
		{validate.KindCloudPlatform, apimeta.ScopePlatform},
		{validate.KindCloudProvider, apimeta.ScopePlatform},
		{validate.KindCloudProviderParticipation, apimeta.ScopeCloudPlatform},
		{validate.KindHostingLocation, apimeta.ScopeCloudProvider},
		{validate.KindDatacenter, apimeta.ScopeCloudProvider},
		{validate.KindFaultDomain, apimeta.ScopeCloudProvider},
		{validate.KindInfrastructureStack, apimeta.ScopeCloudProvider},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(string(tc.kind), func(t *testing.T) {
			t.Parallel()
			got, ok := validate.AllowedScopeKinds(tc.kind)
			if !ok || len(got) != 1 || got[0] != tc.allowed {
				t.Fatalf("got %#v", got)
			}
			var scope *apimeta.ScopeRef
			if tc.allowed == apimeta.ScopePlatform {
				scope = nil
			} else {
				scope = &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
					APIVersion: "core.sovrunn.io/v1alpha1",
					Kind:       string(tc.allowed),
					Name:       "scope",
					UID:        "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				}}
			}
			if prob := validate.ValidateScopeKindSubset(tc.kind, scope); prob != nil {
				t.Fatalf("allowed scope rejected: %#v", prob)
			}
			bad := &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: "v", Kind: string(apimeta.ScopeOrganization), Name: "o", UID: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			}}
			if prob := validate.ValidateScopeKindSubset(tc.kind, bad); prob == nil {
				t.Fatal("expected scope kind rejection")
			}
		})
	}
}

func TestValidateParticipationScopeUIDInvariant(t *testing.T) {
	t.Parallel()
	uid := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	p := model.CloudProviderParticipation{
		Metadata: apimeta.ObjectMeta{
			Name: "part",
			ScopeRef: &apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform,
				Kind:       model.KindCloudPlatform,
				Name:       "plat",
				UID:        uid,
			}},
		},
		Spec: model.CloudProviderParticipationSpec{
			CloudPlatformRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudPlatform,
				Kind:       model.KindCloudPlatform,
				Name:       "plat",
				UID:        uid,
			},
			CloudProviderRef: apimeta.TypedRef{
				APIVersion: model.APIVersionCloudProvider,
				Kind:       model.KindCloudProvider,
				Name:       "prov",
				UID:        "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			},
			Environment: model.ParticipationEnvironmentDevelopment,
		},
	}
	if prob := validate.ValidateParticipationScopeUIDInvariant(p); prob != nil {
		t.Fatalf("matching uids: %#v", prob)
	}
	p.Spec.CloudPlatformRef.UID = "cccccccccccccccccccccccccccccccc"
	prob := validate.ValidateParticipationScopeUIDInvariant(p)
	if prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("mismatch: %#v", prob)
	}
	if len(prob.Violations) == 0 || prob.Violations[0].Code != validate.ViolationScopeReferenceMismatch {
		t.Fatalf("want VS0_SCOPE_REFERENCE_MISMATCH, got %#v", prob.Violations)
	}
}
