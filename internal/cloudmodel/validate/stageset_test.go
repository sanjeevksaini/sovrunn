package validate_test

import (
	"context"
	"testing"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/apivalid"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/validate"
)

func TestStageSetFor_RetainsInheritedOrderWithoutHTTP(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		kind  validate.Kind
		class validate.RequestClass
		obj   any
	}{
		{
			name:  "collection-create",
			kind:  validate.KindCloudPlatform,
			class: validate.RequestClassCollectionCreate,
			obj: model.CloudPlatform{
				Metadata: apimeta.ObjectMeta{Name: "plat"},
				Spec: model.CloudPlatformSpec{
					OwnerRegistration: model.OwnerRegistration{
						LegalName: "Acme", RegistrationIdentifier: "R1", JurisdictionCode: "US",
					},
				},
			},
		},
		{
			name:  "patch",
			kind:  validate.KindCloudProvider,
			class: validate.RequestClassPatch,
			obj: model.CloudProvider{
				Metadata: apimeta.ObjectMeta{Name: "prov"},
				Spec:     model.CloudProviderSpec{OperatingMarkets: []string{"US"}},
			},
		},
		{
			name:  "item-action",
			kind:  validate.KindCloudProviderParticipation,
			class: validate.RequestClassItemAction,
			obj:   model.CloudProviderParticipation{Metadata: apimeta.ObjectMeta{Name: "part"}},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			stages, ok := validate.StageSetFor(tc.kind, tc.class)
			if !ok {
				t.Fatal("StageSetFor failed")
			}
			if stages.Defaulting == nil || stages.Semantic == nil || stages.Reference == nil {
				t.Fatal("all StageSet slots required")
			}

			var order []string
			wrapped := apivalid.StageSet{
				Defaulting: orderDefaulting{inner: stages.Defaulting, calls: &order},
				Semantic:   orderValidation{inner: stages.Semantic, calls: &order, label: "semantic"},
				Reference:  orderValidation{inner: stages.Reference, calls: &order, label: "reference"},
			}

			ctx := context.Background()
			obj, err := wrapped.Defaulting.Apply(ctx, tc.obj)
			if err != nil {
				t.Fatalf("defaulting: %v", err)
			}
			if _, err := wrapped.Semantic.Validate(ctx, obj); err != nil {
				t.Fatalf("semantic: %v", err)
			}
			if _, err := wrapped.Reference.Validate(ctx, obj); err != nil {
				t.Fatalf("reference: %v", err)
			}
			want := []string{"defaulting", "semantic", "reference"}
			if len(order) != 3 || order[0] != want[0] || order[1] != want[1] || order[2] != want[2] {
				t.Fatalf("order=%v want %v", order, want)
			}
		})
	}
}

type orderDefaulting struct {
	inner apivalid.DefaultingStage
	calls *[]string
}

func (o orderDefaulting) Apply(ctx context.Context, object any) (any, error) {
	*o.calls = append(*o.calls, "defaulting")
	return o.inner.Apply(ctx, object)
}

type orderValidation struct {
	inner apivalid.ValidationStage
	calls *[]string
	label string
}

func (o orderValidation) Validate(ctx context.Context, object any) ([]apiproblem.Violation, error) {
	*o.calls = append(*o.calls, o.label)
	return o.inner.Validate(ctx, object)
}

func TestStageSetFor_RejectsInvalidCombinations(t *testing.T) {
	t.Parallel()
	if _, ok := validate.StageSetFor(validate.KindCloudPlatform, validate.RequestClassItemAction); ok {
		t.Fatal("platform item-action must be rejected")
	}
	if _, ok := validate.StageSetFor(validate.KindCloudProviderParticipation, validate.RequestClassPatch); ok {
		t.Fatal("participation PATCH must be rejected")
	}
}
