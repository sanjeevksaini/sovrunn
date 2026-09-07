package uow_test

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/authzeval"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/roleassign"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/uow"
)

type roleAssignIntentAdapter struct {
	intent roleassign.PreparedRoleAssignmentIntent
}

func (a roleAssignIntentAdapter) ParticipantClaim() state.ParticipantClaim {
	return a.intent.ParticipantClaim()
}

func (a roleAssignIntentAdapter) FinalizeAt(permit state.MutationFinalizationPermit) uow.FinalizedIntentOutcome {
	out := a.intent.FinalizeAt(permit)
	if change, ok := out.AsFinalized(); ok {
		return uow.NewFinalizedChangeOutcome(change)
	}
	if domain, ok := out.AsDomain(); ok {
		return uow.NewDomainOutcome(domain.Code())
	}
	if fail, ok := out.AsMechanicalFailure(); ok {
		return uow.NewFailureOutcome(fail)
	}
	return uow.NewFailureOutcome(operation.NewMechanicalFailure("roleassign_finalize_invalid"))
}

func buildEvidenceFinalized(
	finalized authzeval.FinalizedCandidate,
	admission evidence.AuthorizationAdmissionProof,
	ctx evidence.PublicationContext,
) (evidence.FinalizedAuthorizationCandidate, operation.MechanicalFailure) {
	outcome := evidence.AuthorizationDeny
	if finalized.Candidate().Decision() == model.AuthorizationDecisionAllow {
		outcome = evidence.AuthorizationAllow
	}
	target := finalized.Candidate().Input().TargetRef().UID
	if target == "" {
		target = finalized.Candidate().Input().TargetRef().Name
	}
	desc, fail := evidence.NewAuthorizationCandidateDescriptor(
		finalized.Candidate().CandidateDigest(),
		outcome,
		finalized.Candidate().Input().Action(),
		target,
	)
	if fail.Reason() != "" {
		return evidence.FinalizedAuthorizationCandidate{}, fail
	}
	return evidence.NewFinalizedAuthorizationCandidate(desc, admission, ctx)
}

func sampleScope() apimeta.ScopeRef {
	return apimeta.ScopeRef{
		TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1",
			Kind:       string(apimeta.ScopeProject),
			Name:       "proj-1",
			UID:        "proj-1",
		},
	}
}

func sampleSpec(validity model.AssignmentValidityMode) model.RoleAssignmentSpec {
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	spec := model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{Kind: model.RoleHolderPrincipal, Principal: &actor},
		RoleDefinitionRef: apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "admin", UID: "rd-1",
		},
		RoleDefinitionVersion: "v1",
		Validity:              validity,
	}
	if validity == model.AssignmentValidityTimeBound {
		nb := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
		exp := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
		spec.NotBefore = &nb
		spec.ExpiresAt = &exp
	} else {
		rule := apimeta.TypedRef{APIVersion: model.APIVersionGov, Kind: "AccessReviewRule", Name: "standing", UID: "arr-1"}
		spec.AccessReviewRuleRef = &rule
	}
	return spec
}

func allowCandidate(t *testing.T, at time.Time) authzeval.AuthorizationCandidate {
	t.Helper()
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	spec := sampleSpec(model.AssignmentValidityStanding)
	in, ok := model.NewAuthorizationInput(
		actor,
		"roleassignment.grant",
		sampleScope(),
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: roleassign.KindRoleAssignment, Name: "ra-1", UID: "ra-1"},
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
		t.Fatal("authorization input")
	}
	candidate, fail := authzeval.NewEvaluator().Evaluate(in)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return candidate
}

func denyCandidate(t *testing.T, at time.Time) authzeval.AuthorizationCandidate {
	t.Helper()
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	spec := sampleSpec(model.AssignmentValidityStanding)
	in, ok := model.NewAuthorizationInput(
		actor,
		"roleassignment.grant",
		sampleScope(),
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: roleassign.KindRoleAssignment, Name: "ra-1", UID: "ra-1"},
		at,
		"req-2",
		[]model.RoleAssignmentSpec{spec},
		nil,
		nil,
		map[string][]string{"rd-1@v1": []string{"roleassignment.grant"}},
		nil,
		nil,
		true,
		false,
		model.PolicyOutcomeAllow,
	)
	if !ok {
		t.Fatal("authorization input")
	}
	candidate, fail := authzeval.NewEvaluator().Evaluate(in)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return candidate
}

func setupCallerRequest(t *testing.T, candidate authzeval.AuthorizationCandidate, adapter uow.FinalizableIntent) (*state.Store, uow.CallerMutationRequest) {
	t.Helper()
	store := state.NewStore(clock.NewFixedClock(time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)))
	store.SetResourceVersion("ra-1", "1")

	claim := adapter.ParticipantClaim()
	dep, err := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), claim.Versions().Predicates())
	if err != nil {
		t.Fatal(err)
	}
	curr, err := state.NewAuthorizationCurrentnessClaim(candidate.CandidateDigest(), dep)
	if err != nil {
		t.Fatal(err)
	}
	adm, err := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	if err != nil {
		t.Fatal(err)
	}
	opCandidate, fail := authzeval.NewOperationCandidate(candidate, curr)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	key, fail := idempotency.NewIdempotencyLookupKey(actor, "roleassignment.grant", "idem-1")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	binding, fail := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: roleassign.KindRoleAssignment, Name: "ra-1", UID: "ra-1"},
		sampleScope(),
		[]byte("digest-1"),
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return store, uow.CallerMutationRequest{
		LookupKey:     key,
		Binding:       binding,
		Admission:     adm,
		Candidate:     opCandidate,
		BuildEvidence: buildEvidenceFinalized,
		Participants:  []uow.FinalizableIntent{adapter},
	}
}

func mustGrantIntent(t *testing.T) roleassign.PreparedRoleAssignmentIntent {
	return mustGrantIntentWithRole(t, state.OriginatingResult)
}

func mustGrantIntentWithRole(t *testing.T, role state.ResultRole) roleassign.PreparedRoleAssignmentIntent {
	t.Helper()
	tr, fail := roleassign.NewGrantTriggerIntent(
		"p-grant", "ra-1", "ra-1", sampleScope(), sampleSpec(model.AssignmentValidityTimeBound), "1", "2", nil, role, nil,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	intent, fail := roleassign.NewGrantPort().Submit(tr)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return intent
}

func TestCallerMutationAllowFinalizeAndRelease(t *testing.T) {
	t.Parallel()
	intent := roleAssignIntentAdapter{intent: mustGrantIntent(t)}
	store, req := setupCallerRequest(t, allowCandidate(t, time.Date(2026, 9, 7, 11, 0, 0, 0, time.UTC)), intent)
	coordinator := uow.NewCoordinator(store)
	out := coordinator.CoordinateCallerMutation(context.Background(), req)
	if out.Kind() != uow.OutcomeReleased {
		t.Fatalf("kind=%s", out.Kind())
	}
	if _, ok := out.AuthorizationResult(); !ok {
		t.Fatal("authorization result missing")
	}
	if _, ok := out.CommitReceipt(); !ok {
		t.Fatal("commit receipt missing")
	}
}

func TestCallerMutationDenyRoutesToEvidenceOnly(t *testing.T) {
	t.Parallel()
	intent := roleAssignIntentAdapter{intent: mustGrantIntent(t)}
	store, req := setupCallerRequest(t, denyCandidate(t, time.Date(2026, 9, 7, 11, 0, 0, 0, time.UTC)), intent)
	coordinator := uow.NewCoordinator(store)
	out := coordinator.CoordinateCallerMutation(context.Background(), req)
	if out.Kind() != uow.OutcomeReleased {
		t.Fatalf("kind=%s", out.Kind())
	}
	result, ok := out.AuthorizationResult()
	if !ok {
		t.Fatal("authorization result missing")
	}
	if result.Decision() != model.AuthorizationDecisionDeny {
		t.Fatalf("decision=%s", result.Decision())
	}
}

func TestCallerMutationDependencyRevalidationConflict(t *testing.T) {
	t.Parallel()
	intent := roleAssignIntentAdapter{intent: mustGrantIntent(t)}
	store := state.NewStore(clock.NewFixedClock(time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)))
	store.SetResourceVersion("ra-1", "1")
	candidate := allowCandidate(t, time.Date(2026, 9, 7, 11, 0, 0, 0, time.UTC))
	claim := intent.ParticipantClaim()
	dep, err := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), []state.ExpectedVersionPredicate{{ResourceUID: "ra-1", ResourceVersion: "9"}})
	if err != nil {
		t.Fatal(err)
	}
	curr, err := state.NewAuthorizationCurrentnessClaim(candidate.CandidateDigest(), dep)
	if err != nil {
		t.Fatal(err)
	}
	adm, err := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	if err != nil {
		t.Fatal(err)
	}
	opCandidate, fail := authzeval.NewOperationCandidate(candidate, curr)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	key, fail := idempotency.NewIdempotencyLookupKey(actor, "roleassignment.grant", "idem-conflict")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	binding, fail := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: roleassign.KindRoleAssignment, Name: "ra-1", UID: "ra-1"},
		sampleScope(),
		[]byte("digest-x"),
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	out := uow.NewCoordinator(store).CoordinateCallerMutation(context.Background(), uow.CallerMutationRequest{
		LookupKey: key, Binding: binding, Admission: adm, Candidate: opCandidate, BuildEvidence: buildEvidenceFinalized, Participants: []uow.FinalizableIntent{intent},
	})
	if out.Kind() != uow.OutcomeConflict {
		t.Fatalf("kind=%s", out.Kind())
	}
}

func TestCallerMutationReplay(t *testing.T) {
	t.Parallel()
	intent := roleAssignIntentAdapter{intent: mustGrantIntent(t)}
	store, req := setupCallerRequest(t, allowCandidate(t, time.Date(2026, 9, 7, 11, 0, 0, 0, time.UTC)), intent)
	coord := uow.NewCoordinator(store)
	first := coord.CoordinateCallerMutation(context.Background(), req)
	if first.Kind() != uow.OutcomeReleased {
		t.Fatalf("first kind=%s", first.Kind())
	}
	second := coord.CoordinateCallerMutation(context.Background(), req)
	if second.Kind() != uow.OutcomeReplay {
		t.Fatalf("second kind=%s", second.Kind())
	}
	if rec, ok := second.Replay(); !ok || !rec.Sealed() {
		t.Fatal("replay completed record missing")
	}
}

func TestControllerMutationAutomatic(t *testing.T) {
	t.Parallel()
	store := state.NewStore(clock.NewFixedClock(time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)))
	store.SetResourceVersion("ra-1", "1")
	intent := roleAssignIntentAdapter{intent: mustGrantIntentWithRole(t, state.SupportingResult)}
	adm, err := state.NewAutomaticControllerMutationAdmission([]state.ParticipantClaim{intent.ParticipantClaim()})
	if err != nil {
		t.Fatal(err)
	}
	out := uow.NewCoordinator(store).CoordinateControllerMutation(context.Background(), uow.ControllerMutationRequest{
		Admission: adm, BuildEvidence: buildEvidenceFinalized, Participants: []uow.FinalizableIntent{intent},
	})
	if out.Kind() != uow.OutcomeReleased {
		t.Fatalf("kind=%s", out.Kind())
	}
}

func TestEvidenceOnlyTransaction(t *testing.T) {
	t.Parallel()
	store := state.NewStore(clock.NewFixedClock(time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)))
	store.SetResourceVersion("ra-1", "1")
	candidate := denyCandidate(t, time.Date(2026, 9, 7, 11, 0, 0, 0, time.UTC))
	dep, err := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), []state.ExpectedVersionPredicate{{ResourceUID: "ra-1", ResourceVersion: "1"}})
	if err != nil {
		t.Fatal(err)
	}
	curr, err := state.NewAuthorizationCurrentnessClaim(candidate.CandidateDigest(), dep)
	if err != nil {
		t.Fatal(err)
	}
	opCandidate, fail := authzeval.NewOperationCandidate(candidate, curr)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	out := uow.NewCoordinator(store).CoordinateEvidenceOnly(context.Background(), uow.EvidenceOnlyRequest{
		Candidate: opCandidate, BuildEvidence: buildEvidenceFinalized,
	})
	if out.Kind() != uow.OutcomeReleased {
		t.Fatalf("kind=%s", out.Kind())
	}
}

func TestRoleAssignmentPublicationBoundaryThroughUOW(t *testing.T) {
	t.Parallel()
	makeIntents := func(t *testing.T) []roleassign.PreparedRoleAssignmentIntent {
		t.Helper()
		spec := sampleSpec(model.AssignmentValidityTimeBound)
		scope := sampleScope()

		grantTrigger, fail := roleassign.NewGrantTriggerIntent("p-g", "ra-1", "ra-1", scope, spec, "1", "2", nil, state.OriginatingResult, nil)
		if fail.Reason() != "" {
			t.Fatal(fail.Reason())
		}
		grant, fail := roleassign.NewGrantPort().Submit(grantTrigger)
		if fail.Reason() != "" {
			t.Fatal(fail.Reason())
		}

		actTrigger, fail := roleassign.NewActivationTriggerIntent(roleassign.ActivationActivate, "p-a", "ra-1", "ra-1", scope, spec, true, "1", "2", nil, state.OriginatingResult)
		if fail.Reason() != "" {
			t.Fatal(fail.Reason())
		}
		activation, fail := roleassign.NewActivationPort().Submit(actTrigger)
		if fail.Reason() != "" {
			t.Fatal(fail.Reason())
		}

		reviewTrigger, fail := roleassign.NewReviewRemediationTriggerIntent(roleassign.RemediationRevoke, "p-r", "ra-1", "ra-1", scope, spec, true, model.RoleAssignmentSpec{}, false, "1", "2", nil, state.OriginatingResult, nil)
		if fail.Reason() != "" {
			t.Fatal(fail.Reason())
		}
		review, fail := roleassign.NewReviewRemediationPort().Submit(reviewTrigger)
		if fail.Reason() != "" {
			t.Fatal(fail.Reason())
		}

		directIn, fail := roleassign.NewDirectRevokeIntent("p-d", "ra-1", "ra-1", scope, spec, true, "1", "2", nil, state.OriginatingResult)
		if fail.Reason() != "" {
			t.Fatal(fail.Reason())
		}
		direct, fail := roleassign.NewDirectRevokePort().Prepare(directIn)
		if fail.Reason() != "" {
			t.Fatal(fail.Reason())
		}
		return []roleassign.PreparedRoleAssignmentIntent{grant, activation, review, direct}
	}

	for _, raw := range makeIntents(t) {
		intent := roleAssignIntentAdapter{intent: raw}
		store, req := setupCallerRequest(t, allowCandidate(t, time.Date(2026, 9, 7, 11, 0, 0, 0, time.UTC)), intent)
		out := uow.NewCoordinator(store).CoordinateCallerMutation(context.Background(), req)
		if out.Kind() != uow.OutcomeReleased {
			t.Fatalf("port=%s op=%s kind=%s", raw.TriggerPort(), raw.Operation(), out.Kind())
		}
	}
}

func TestPropertyUOWConstructsNoSemanticValue(t *testing.T) {
	t.Parallel()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	dir := filepath.Dir(thisFile)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	forbiddenCallees := map[string]bool{
		"NewAuthorizationInput":         true,
		"NewAuthorizationResult":        true,
		"NewMutationEventFacts":         true,
		"NewMutationDescriptor":         true,
		"NewDomainConclusionDescriptor": true,
		"NewDomainDecisionMaterial":     true,
	}
	for _, pkg := range pkgs {
		for name, file := range pkg.Files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			ast.Inspect(file, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				switch fn := call.Fun.(type) {
				case *ast.SelectorExpr:
					if forbiddenCallees[fn.Sel.Name] {
						t.Fatalf("uow must not construct semantic value via %s in %s", fn.Sel.Name, name)
					}
				case *ast.Ident:
					if forbiddenCallees[fn.Name] {
						t.Fatalf("uow must not call %s in %s", fn.Name, name)
					}
				}
				return true
			})
		}
	}

	prop := func(v uint8) bool {
		store := state.NewStore(clock.NewFixedClock(time.Date(2026, 9, 7, 12, 0, int(v), 0, time.UTC)))
		return uow.NewCoordinator(store) != nil
	}
	if err := quick.Check(prop, &quick.Config{MaxCount: 8}); err != nil {
		t.Fatal(err)
	}
}
