package model_test

import (
	"reflect"
	"testing"
	"testing/quick"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

func TestApprovalRequirementRequiredAndNotRequired(t *testing.T) {
	t.Parallel()
	notRequired := model.NewApprovalRequirementNotRequired()
	if !notRequired.Valid() || notRequired.Mode() != model.ApprovalRequirementNotRequired {
		t.Fatal("expected valid NotRequired")
	}
	if _, ok := notRequired.PolicyRef(); ok {
		t.Fatal("NotRequired must not carry policy ref")
	}

	required, ok := model.NewApprovalRequirementRequired(apimeta.TypedRef{
		APIVersion: model.APIVersionIAM, Kind: "ApprovalPolicy", Name: "pol-1", UID: "ap-1",
	})
	if !ok || !required.Valid() || !required.IsRequired() {
		t.Fatal("expected valid Required requirement")
	}
	ref, hasRef := required.PolicyRef()
	if !hasRef || ref.UID != "ap-1" {
		t.Fatal("expected policy ref")
	}
	ref.UID = "mutated"
	ref2, _ := required.PolicyRef()
	if ref2.UID == "mutated" {
		t.Fatal("policy accessor must return copy")
	}
}

func TestApprovalRequirementDeterministicAndNoSetter(t *testing.T) {
	t.Parallel()
	fn := func(uid string) bool {
		if uid == "" {
			uid = "ap-fixed"
		}
		a, okA := model.NewApprovalRequirementRequired(apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: "ApprovalPolicy", Name: "policy", UID: uid,
		})
		b, okB := model.NewApprovalRequirementRequired(apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: "ApprovalPolicy", Name: "policy", UID: uid,
		})
		refA, hasA := a.PolicyRef()
		refB, hasB := b.PolicyRef()
		return okA == okB && a.Valid() == b.Valid() && hasA == hasB && refA == refB
	}
	if err := quick.Check(fn, &quick.Config{MaxCount: 32}); err != nil {
		t.Fatal(err)
	}

	typ := reflect.TypeOf(model.ApprovalRequirement{})
	for i := 0; i < typ.NumMethod(); i++ {
		if len(typ.Method(i).Name) >= 3 && typ.Method(i).Name[:3] == "Set" {
			t.Fatalf("setter method forbidden: %s", typ.Method(i).Name)
		}
	}
}
