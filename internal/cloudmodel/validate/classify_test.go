package validate_test

import (
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func TestClassifyCreateContract_CloudPlatformAllowlist(t *testing.T) {
	t.Parallel()
	body := []byte(`{
		"metadata":{"name":"acme-platform","displayName":"Acme"},
		"spec":{"ownerRegistration":{"legalName":"Acme Inc","registrationIdentifier":"REG-1","jurisdictionCode":"US"},"description":"d"}
	}`)
	if prob := validate.ClassifyCreateContract(validate.KindCloudPlatform, body); prob != nil {
		t.Fatalf("valid create: %#v", prob)
	}
}

func TestClassifyCreateContract_RejectsServerOwnedStatusScopeUnknownDeferred(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		kind validate.Kind
		body string
		code apiproblem.ErrorCode
	}{
		{
			name: "status",
			kind: validate.KindCloudPlatform,
			body: `{"metadata":{"name":"p"},"spec":{"ownerRegistration":{"legalName":"L","registrationIdentifier":"R","jurisdictionCode":"US"}},"status":{"phase":"Active"}}`,
			code: apiproblem.CodeAuthorizationDenied,
		},
		{
			name: "system-owned uid",
			kind: validate.KindCloudPlatform,
			body: `{"metadata":{"name":"p","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},"spec":{"ownerRegistration":{"legalName":"L","registrationIdentifier":"R","jurisdictionCode":"US"}}}`,
			code: apiproblem.CodeAuthorizationDenied,
		},
		{
			name: "scopeRef",
			kind: validate.KindCloudProvider,
			body: `{"metadata":{"name":"p","scopeRef":{"apiVersion":"v","kind":"Platform","name":"platform"}},"spec":{"operatingMarkets":["US"]}}`,
			code: apiproblem.CodeValidationFailed,
		},
		{
			name: "unknown field",
			kind: validate.KindCloudProvider,
			body: `{"metadata":{"name":"p"},"spec":{"operatingMarkets":["US"],"extra":"x"}}`,
			code: apiproblem.CodeUnknownField,
		},
		{
			name: "deferred providerSelectionModes",
			kind: validate.KindCloudProviderParticipation,
			body: `{
				"metadata":{"name":"part"},
				"spec":{
					"cloudPlatformRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudPlatform","name":"p","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
					"cloudProviderRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudProvider","name":"q","uid":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
					"environment":"development",
					"providerSelectionModes":["CustomerSelected"]
				}
			}`,
			code: apiproblem.CodeUnknownField,
		},
		{
			name: "deferred permittedHostingLocationRefs",
			kind: validate.KindCloudProviderParticipation,
			body: `{
				"metadata":{"name":"part"},
				"spec":{
					"cloudPlatformRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudPlatform","name":"p","uid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
					"cloudProviderRef":{"apiVersion":"core.sovrunn.io/v1alpha1","kind":"CloudProvider","name":"q","uid":"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
					"environment":"development",
					"permittedHostingLocationRefs":[]
				}
			}`,
			code: apiproblem.CodeUnknownField,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			prob := validate.ClassifyCreateContract(tc.kind, []byte(tc.body))
			if prob == nil || prob.Code != tc.code {
				t.Fatalf("got %#v, want code %s", prob, tc.code)
			}
		})
	}
}

func TestClassifyCreateContract_PerKindRequiredOptional(t *testing.T) {
	t.Parallel()
	kinds := []validate.Kind{
		validate.KindCloudPlatform,
		validate.KindCloudProvider,
		validate.KindCloudProviderParticipation,
		validate.KindHostingLocation,
		validate.KindDatacenter,
		validate.KindFaultDomain,
		validate.KindInfrastructureStack,
	}
	for _, kind := range kinds {
		kind := kind
		t.Run(string(kind), func(t *testing.T) {
			t.Parallel()
			class, ok := validate.CreateClassification(kind)
			if !ok {
				t.Fatal("missing classification")
			}
			if len(class.ClientRequired) == 0 {
				t.Fatal("client-required must be non-empty")
			}
			if kind == validate.KindCloudProviderParticipation && len(class.ClientOptional) != 0 {
				t.Fatal("participation has no client-optional fields")
			}
			if kind == validate.KindCloudProviderParticipation && len(class.Deferred) != 2 {
				t.Fatalf("deferred want 2, got %d", len(class.Deferred))
			}
		})
	}
}

func TestClassifyPatchContract_RejectsImmutableAndUnknown(t *testing.T) {
	t.Parallel()
	prob := validate.ClassifyPatchContract(validate.KindCloudPlatform, []byte(`{"spec":{"description":"x"}}`))
	if prob != nil {
		t.Fatalf("valid patch: %#v", prob)
	}
	prob = validate.ClassifyPatchContract(validate.KindCloudPlatform, []byte(`{"metadata":{"name":"other"}}`))
	if prob == nil || prob.Code != apiproblem.CodeValidationFailed {
		t.Fatalf("name patch: %#v", prob)
	}
	prob = validate.ClassifyPatchContract(validate.KindCloudProvider, []byte(`{"spec":{"displayName":"n","operatingMarkets":["US"]}}`))
	if prob != nil {
		t.Fatalf("provider patch: %#v", prob)
	}
}
