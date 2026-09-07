package model_test

import (
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
)

func sampleAuthorizationInput(t *testing.T, at time.Time) model.AuthorizationInput {
	t.Helper()
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	holder := actor
	spec := model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{Kind: model.RoleHolderPrincipal, Principal: &holder},
		RoleDefinitionRef: apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "ops-admin", UID: "rd-1",
		},
		RoleDefinitionVersion: "v1",
		Validity:              model.AssignmentValidityStanding,
	}
	in, ok := model.NewAuthorizationInput(
		actor,
		"roleassignment.grant",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: string(apimeta.ScopeProject), Name: "p1", UID: "scope-1"}},
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "RoleAssignment", Name: "ra-1", UID: "ra-1"},
		at,
		"req-1",
		[]model.RoleAssignmentSpec{spec},
		nil,
		nil,
		map[string][]string{"rd-1@v1": {"roleassignment.grant"}},
		nil,
		nil,
		true,
		true,
		model.PolicyOutcomeAllow,
	)
	if !ok {
		t.Fatal("expected valid AuthorizationInput")
	}
	return in
}

func TestAuthorizationInputDeepCopyAndReadOnlyAccessors(t *testing.T) {
	t.Parallel()
	at := time.Date(2026, 9, 7, 9, 0, 0, 0, time.UTC)
	in := sampleAuthorizationInput(t, at)

	direct := in.DirectAssignments()
	direct[0].RoleDefinitionVersion = "mutated"
	if in.DirectAssignments()[0].RoleDefinitionVersion == "mutated" {
		t.Fatal("direct assignment accessor must return copy")
	}

	actions := in.RoleActions()
	actions["rd-1@v1"][0] = "mutated"
	if in.RoleActions()["rd-1@v1"][0] == "mutated" {
		t.Fatal("role action accessor must return deep copy")
	}

	typ := reflect.TypeOf(model.AuthorizationInput{})
	for i := 0; i < typ.NumMethod(); i++ {
		if len(typ.Method(i).Name) >= 3 && typ.Method(i).Name[:3] == "Set" {
			t.Fatalf("setter method forbidden: %s", typ.Method(i).Name)
		}
	}
}

func TestAuthorizationResultConstructionAndDeepCopy(t *testing.T) {
	t.Parallel()
	in := sampleAuthorizationInput(t, time.Date(2026, 9, 7, 9, 0, 0, 0, time.UTC))
	out, ok := model.NewAuthorizationResult(
		model.AuthorizationDecisionAllow,
		[]string{"MATCHED_ASSIGNMENT"},
		in.RequestID(),
		in.Action(),
		in.TargetRef(),
		in.EvaluatedAt(),
		in.EvaluatedAt().Add(time.Second),
		in.DirectAssignments(),
	)
	if !ok || !out.Sealed() {
		t.Fatal("expected sealed AuthorizationResult")
	}
	reasons := out.ReasonCodes()
	reasons[0] = "MUTATED"
	if out.ReasonCodes()[0] == "MUTATED" {
		t.Fatal("reason accessor must return copy")
	}
	cons := out.ContributingAssignments()
	cons[0].RoleDefinitionVersion = "mutated"
	if out.ContributingAssignments()[0].RoleDefinitionVersion == "mutated" {
		t.Fatal("contributing assignments accessor must return deep copy")
	}
}

func TestAuthorizationValuesValidateDeterministically(t *testing.T) {
	t.Parallel()
	f := func(sec uint8) bool {
		at := time.Date(2026, 9, 7, 9, 0, int(sec), 0, time.UTC)
		in1 := sampleAuthorizationInput(t, at)
		in2 := sampleAuthorizationInput(t, at)
		r1, ok1 := model.NewAuthorizationResult(
			model.AuthorizationDecisionDeny,
			[]string{"POLICY_DENY"},
			in1.RequestID(),
			in1.Action(),
			in1.TargetRef(),
			in1.EvaluatedAt(),
			in1.EvaluatedAt().Add(time.Millisecond),
			in1.DirectAssignments(),
		)
		r2, ok2 := model.NewAuthorizationResult(
			model.AuthorizationDecisionDeny,
			[]string{"POLICY_DENY"},
			in2.RequestID(),
			in2.Action(),
			in2.TargetRef(),
			in2.EvaluatedAt(),
			in2.EvaluatedAt().Add(time.Millisecond),
			in2.DirectAssignments(),
		)
		return in1.Sealed() == in2.Sealed() &&
			ok1 == ok2 &&
			r1.Sealed() == r2.Sealed() &&
			reflect.DeepEqual(r1.ReasonCodes(), r2.ReasonCodes())
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 32}); err != nil {
		t.Fatal(err)
	}
}
