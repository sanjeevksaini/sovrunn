package authzeval_test

import (
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/authzeval"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

func mkInput(t *testing.T, at time.Time, opts ...func(*model.RoleAssignmentSpec, *model.AuthorizationInput)) model.AuthorizationInput {
	t.Helper()
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	spec := model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{Kind: model.RoleHolderPrincipal, Principal: &actor},
		RoleDefinitionRef: apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "ops", UID: "rd-1",
		},
		RoleDefinitionVersion: "v1",
		Validity:              model.AssignmentValidityStanding,
	}
	target := apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "RoleAssignment", Name: "ra-1", UID: "ra-1"}
	input, ok := model.NewAuthorizationInput(
		actor,
		"roleassignment.grant",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: string(apimeta.ScopeProject), Name: "proj", UID: "scope-1"}},
		target,
		at,
		"req-1",
		[]model.RoleAssignmentSpec{spec},
		nil,
		nil,
		map[string][]string{"rd-1@v1": []string{"roleassignment.grant"}},
		nil,
		nil,
		true,
		true,
		model.PolicyOutcomeAllow,
	)
	if !ok {
		t.Fatal("input invalid")
	}
	return input
}

func makeSealedReceipt(t *testing.T, now time.Time) (state.CommitReceipt, evidence.AuthorizationAdmissionProof) {
	t.Helper()
	store := state.NewStore(clock.NewFixedClock(now))
	store.SetResourceVersion("ra-1", "1")
	candidateDigest, fail := evidence.NewCandidateDigest([]byte("candidate"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	dep, err := state.NewAuthorizationDependencyVersionSet(
		store.SnapshotVersion(),
		[]state.ExpectedVersionPredicate{{ResourceUID: "ra-1", ResourceVersion: "1"}},
	)
	if err != nil {
		t.Fatal(err)
	}
	claim, err := state.NewAuthorizationCurrentnessClaim(candidateDigest, dep)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := store.BeginEvidenceTransaction(claim)
	if err != nil {
		t.Fatal(err)
	}
	desc, fail := evidence.NewAuthorizationCandidateDescriptor(candidateDigest, evidence.AuthorizationAllow, "roleassignment.grant", "ra-1")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	fin, fail := evidence.NewFinalizedAuthorizationCandidate(desc, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	completed, fail := evidence.CompleteAuthorizationCandidate(fin, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	prepared, fail := evidence.PrepareAuthorizationCarrierSet(completed)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if fail := tx.AcceptPreparedEvidence(prepared); fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if fail := tx.Seal(); fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	receipt, fail := tx.Commit()
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return receipt, tx.AuthorizationAdmissionProof()
}

func TestEvaluatorGrantUnionAndGuardrailIntersection(t *testing.T) {
	t.Parallel()
	ev := authzeval.NewEvaluator()
	at := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	input := mkInput(t, at)
	candidate, fail := ev.Evaluate(input)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if candidate.Decision() != model.AuthorizationDecisionAllow {
		t.Fatalf("decision=%s", candidate.Decision())
	}

	denyInput, ok := model.NewAuthorizationInput(
		input.Actor(),
		input.Action(),
		input.ScopeRef(),
		input.TargetRef(),
		at,
		input.RequestID(),
		input.DirectAssignments(),
		nil,
		nil,
		input.RoleActions(),
		nil,
		nil,
		true,
		false,
		model.PolicyOutcomeAllow,
	)
	if !ok {
		t.Fatal("expected deny input construction")
	}
	denyCandidate, fail := ev.Evaluate(denyInput)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if denyCandidate.Decision() != model.AuthorizationDecisionDeny {
		t.Fatal("guardrail deny must deny")
	}
}

func TestEvaluatorIndependentApplicabilityMembershipValidityLifecycle(t *testing.T) {
	t.Parallel()
	ev := authzeval.NewEvaluator()
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	groupRef := model.AccessGroupRef{TypedRef: apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: model.KindAccessGroup, Name: "ops", UID: "grp-1"}}
	nb := time.Date(2026, 9, 7, 9, 0, 0, 0, time.UTC)
	exp := nb.Add(time.Hour)
	groupSpec := model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{Kind: model.RoleHolderAccessGroup, AccessGroup: &groupRef},
		RoleDefinitionRef: apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "ops", UID: "rd-2",
		},
		RoleDefinitionVersion: "v1",
		Validity:              model.AssignmentValidityTimeBound,
		NotBefore:             &nb,
		ExpiresAt:             &exp,
	}
	input, ok := model.NewAuthorizationInput(
		actor, "roleassignment.grant",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: string(apimeta.ScopeProject), Name: "proj", UID: "scope-1"}},
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "RoleAssignment", Name: "ra", UID: "ra"},
		nb.Add(30*time.Minute),
		"req-2",
		nil,
		[]model.RoleAssignmentSpec{groupSpec},
		[]string{"grp-1"},
		map[string][]string{"rd-2@v1": []string{"roleassignment.grant"}},
		nil,
		nil,
		true,
		true,
		model.PolicyOutcomeAllow,
	)
	if !ok {
		t.Fatal("input invalid")
	}
	candidate, fail := ev.Evaluate(input)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if candidate.Decision() != model.AuthorizationDecisionAllow {
		t.Fatal("group assignment with active membership should allow")
	}

	expiredInput, ok := model.NewAuthorizationInput(
		actor, input.Action(), input.ScopeRef(), input.TargetRef(),
		exp.Add(time.Second), input.RequestID(),
		nil, []model.RoleAssignmentSpec{groupSpec}, []string{"grp-1"},
		input.RoleActions(), nil, nil, true, true, model.PolicyOutcomeAllow,
	)
	if !ok {
		t.Fatal("expired input invalid")
	}
	expiredCandidate, fail := ev.Evaluate(expiredInput)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if expiredCandidate.Decision() != model.AuthorizationDecisionDeny {
		t.Fatal("expired grant must deny")
	}
}

func TestPrivilegedConstraintNoStandingOrAccessGroup(t *testing.T) {
	t.Parallel()
	spec := model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{
			Kind: model.RoleHolderAccessGroup,
			AccessGroup: &model.AccessGroupRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionIAM, Kind: model.KindAccessGroup, Name: "g", UID: "g-1",
			}},
		},
		RoleDefinitionRef: apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "priv", UID: "rd-priv",
		},
		RoleDefinitionVersion: "v1",
		Validity:              model.AssignmentValidityStanding,
	}
	if authzeval.PrivilegedRoleEligibleGrant(spec) {
		t.Fatal("standing/group privileged grant must be rejected")
	}
}

func TestFinalizeCandidateAtRevalidatesPredicates(t *testing.T) {
	t.Parallel()
	ev := authzeval.NewEvaluator()
	nb := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	exp := nb.Add(2 * time.Minute)
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	spec := model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{Kind: model.RoleHolderPrincipal, Principal: &actor},
		RoleDefinitionRef: apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "ops", UID: "rd-1",
		},
		RoleDefinitionVersion: "v1",
		Validity:              model.AssignmentValidityTimeBound,
		NotBefore:             &nb,
		ExpiresAt:             &exp,
	}
	in, ok := model.NewAuthorizationInput(
		actor, "roleassignment.grant",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: string(apimeta.ScopeProject), Name: "proj", UID: "scope-1"}},
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "RoleAssignment", Name: "ra", UID: "ra"},
		nb.Add(time.Minute), "req-finalize",
		[]model.RoleAssignmentSpec{spec}, nil, nil,
		map[string][]string{"rd-1@v1": []string{"roleassignment.grant"}},
		nil, nil, true, true, model.PolicyOutcomeAllow,
	)
	if !ok {
		t.Fatal("input")
	}
	candidate, fail := ev.Evaluate(in)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	receipt, admission := makeSealedReceipt(t, exp.Add(time.Minute))
	out := authzeval.FinalizeCandidateAt(candidate, admission, receipt.PublicationContext().PublicationInstant())
	if _, ok := out.AsReevaluate(); !ok {
		t.Fatal("expected reevaluate when publication is after expiry")
	}
}

func TestReleaseAfterPublicationRequiresCommitReceipt(t *testing.T) {
	t.Parallel()
	ev := authzeval.NewEvaluator()
	in := mkInput(t, time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC))
	candidate, fail := ev.Evaluate(in)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	receipt, admission := makeSealedReceipt(t, time.Date(2026, 9, 7, 10, 5, 0, 0, time.UTC))
	finalized := authzeval.FinalizeCandidateAt(candidate, admission, receipt.PublicationContext().PublicationInstant())
	ready, ok := finalized.AsFinalized()
	if !ok {
		t.Fatal("expected finalized candidate")
	}
	if _, fail := authzeval.ReleaseAfterPublication(ready, state.CommitReceipt{}); fail.Reason() == "" {
		t.Fatal("expected unsealed receipt failure")
	}
	result, fail := authzeval.ReleaseAfterPublication(ready, receipt)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if result.Decision() != model.AuthorizationDecisionAllow {
		t.Fatalf("decision=%s", result.Decision())
	}
}

func TestCandidateIsNotAuthorizationResultProperty(t *testing.T) {
	t.Parallel()
	fn := func(s uint8) bool {
		at := time.Date(2026, 9, 7, 10, 0, int(s), 0, time.UTC)
		in := mkInput(t, at)
		cand, fail := authzeval.NewEvaluator().Evaluate(in)
		if fail.Reason() != "" {
			return false
		}
		return reflect.TypeOf(cand).Name() != reflect.TypeOf(model.AuthorizationResult{}).Name()
	}
	if err := quick.Check(fn, &quick.Config{MaxCount: 32}); err != nil {
		t.Fatal(err)
	}
}
