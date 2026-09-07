package command_test

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
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/command"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/roleassign"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/uow"
)

type recordingCoordinator struct {
	called bool
	req    uow.CallerMutationRequest
}

func (r *recordingCoordinator) CoordinateCallerMutation(_ context.Context, req uow.CallerMutationRequest) uow.MutationOutcome {
	r.called = true
	r.req = req
	return uow.MutationOutcome{}
}

func (r *recordingCoordinator) CoordinateControllerMutation(context.Context, uow.ControllerMutationRequest) uow.MutationOutcome {
	return uow.MutationOutcome{}
}

func (r *recordingCoordinator) CoordinateEvidenceOnly(context.Context, uow.EvidenceOnlyRequest) uow.MutationOutcome {
	return uow.MutationOutcome{}
}

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
	return uow.NewFailureOutcome(operation.NewMechanicalFailure("finalize_outcome_invalid"))
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

func sampleSpec() model.RoleAssignmentSpec {
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	nb := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	exp := time.Date(2026, 12, 1, 0, 0, 0, 0, time.UTC)
	return model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{Kind: model.RoleHolderPrincipal, Principal: &actor},
		RoleDefinitionRef: apimeta.TypedRef{
			APIVersion: model.APIVersionIAM,
			Kind:       "RoleDefinition",
			Name:       "admin",
			UID:        "rd-1",
		},
		RoleDefinitionVersion: "v1",
		Validity:              model.AssignmentValidityTimeBound,
		NotBefore:             &nb,
		ExpiresAt:             &exp,
	}
}

func sampleAuthorizationInput(t *testing.T, at time.Time) model.AuthorizationInput {
	t.Helper()
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	in, ok := model.NewAuthorizationInput(
		actor,
		"roleassignment.revoke",
		sampleScope(),
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: roleassign.KindRoleAssignment, Name: "ra-1", UID: "ra-1"},
		at,
		"req-revoke-1",
		[]model.RoleAssignmentSpec{sampleSpec()},
		nil,
		nil,
		map[string][]string{"rd-1@v1": []string{"roleassignment.revoke"}},
		nil,
		nil,
		true,
		true,
		model.PolicyOutcomeAllow,
	)
	if !ok {
		t.Fatal("authorization input")
	}
	return in
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
	desc, fail := evidence.NewAuthorizationCandidateDescriptor(
		finalized.Candidate().CandidateDigest(),
		outcome,
		finalized.Candidate().Input().Action(),
		finalized.Candidate().Input().TargetRef().UID,
	)
	if fail.Reason() != "" {
		return evidence.FinalizedAuthorizationCandidate{}, fail
	}
	return evidence.NewFinalizedAuthorizationCandidate(desc, admission, ctx)
}

func TestRoleAssignmentRevokeServiceOrchestration(t *testing.T) {
	t.Parallel()

	coord := &recordingCoordinator{}
	svc := command.NewRoleAssignmentRevokeService(authzeval.NewEvaluator(), roleassign.NewDirectRevokePort(), coord)
	store := state.NewStore(clock.NewFixedClock(time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)))
	store.SetResourceVersion("ra-1", "1")
	authInput := sampleAuthorizationInput(t, time.Date(2026, 9, 7, 11, 0, 0, 0, time.UTC))
	directIntent, fail := roleassign.NewDirectRevokeIntent(
		"p-revoke", "ra-1", "ra-1", sampleScope(), sampleSpec(), true, "1", "2", nil, state.OriginatingResult,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}

	var preparedIntent roleassign.PreparedRoleAssignmentIntent
	var builtCandidate authzeval.OperationCandidate
	_, fail = svc.Execute(context.Background(), command.RoleAssignmentRevokeCommand{
		AuthorizationInput: authInput,
		DirectRevokeIntent: directIntent,
		BindCandidate: func(candidate authzeval.AuthorizationCandidate) (authzeval.OperationCandidate, operation.MechanicalFailure) {
			dep, err := state.NewAuthorizationDependencyVersionSet(
				store.SnapshotVersion(),
				[]state.ExpectedVersionPredicate{{ResourceUID: "ra-1", ResourceVersion: "1"}},
			)
			if err != nil {
				t.Fatal(err)
			}
			curr, err := state.NewAuthorizationCurrentnessClaim(candidate.CandidateDigest(), dep)
			if err != nil {
				t.Fatal(err)
			}
			opCandidate, f := authzeval.NewOperationCandidate(candidate, curr)
			if f.Reason() == "" {
				builtCandidate = opCandidate
			}
			return opCandidate, f
		},
		BuildCallerMutation: func(
			candidate authzeval.OperationCandidate,
			intent roleassign.PreparedRoleAssignmentIntent,
		) (uow.CallerMutationRequest, operation.MechanicalFailure) {
			preparedIntent = intent
			claim := intent.ParticipantClaim()
			adm, err := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, candidate.Currentness())
			if err != nil {
				t.Fatal(err)
			}
			actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
			key, f := idempotency.NewIdempotencyLookupKey(actor, "roleassignment.revoke", "idem-rd08")
			if f.Reason() != "" {
				t.Fatal(f.Reason())
			}
			binding, f := idempotency.NewRequestBinding(
				apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: roleassign.KindRoleAssignment, Name: "ra-1", UID: "ra-1"},
				sampleScope(),
				[]byte("digest-rd08"),
			)
			if f.Reason() != "" {
				t.Fatal(f.Reason())
			}
			return uow.CallerMutationRequest{
				LookupKey:     key,
				Binding:       binding,
				Admission:     adm,
				Candidate:     candidate,
				BuildEvidence: buildEvidenceFinalized,
				Participants:  []uow.FinalizableIntent{roleAssignIntentAdapter{intent: intent}},
			}, operation.MechanicalFailure{}
		},
	})
	if fail.Reason() != "" {
		t.Fatalf("execute fail=%s", fail.Reason())
	}
	if !coord.called {
		t.Fatal("expected coordinator to be called")
	}
	if !preparedIntent.Sealed() || preparedIntent.Operation() != roleassign.OpDirectRevoke {
		t.Fatalf("expected direct revoke prepared intent, sealed=%v op=%s", preparedIntent.Sealed(), preparedIntent.Operation())
	}
	if !builtCandidate.Sealed() || builtCandidate.Candidate().Decision() != model.AuthorizationDecisionAllow {
		t.Fatal("expected allow operation candidate")
	}
	if !coord.req.Candidate.Sealed() {
		t.Fatal("coordinator request candidate must be sealed")
	}
	if len(coord.req.Participants) != 1 {
		t.Fatalf("participants=%d", len(coord.req.Participants))
	}
	if !coord.req.Participants[0].ParticipantClaim().Equal(preparedIntent.ParticipantClaim()) {
		t.Fatal("participant claim must come from prepared direct-revoke intent")
	}
}

func TestCommandConstructsNoDomainValue(t *testing.T) {
	t.Parallel()
	files := parseCommandPackageFiles(t)
	forbiddenCalls := map[string]bool{
		"newPreparedRoleAssignmentIntent":   true,
		"NewGrantTriggerIntent":             true,
		"NewActivationTriggerIntent":        true,
		"NewReviewRemediationTriggerIntent": true,
		"NewPreparedRoleAssignmentIntent":   true,
		"NewFinalizedRoleAssignmentChange":  true,
		"NewAuthorizationResult":            true,
		"NewDomainDecisionMaterial":         true,
		"NewDomainConclusionDescriptor":     true,
		"NewMutationDescriptor":             true,
	}
	for name, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			switch fn := call.Fun.(type) {
			case *ast.Ident:
				if forbiddenCalls[fn.Name] {
					t.Fatalf("%s calls forbidden constructor %s", name, fn.Name)
				}
			case *ast.SelectorExpr:
				if forbiddenCalls[fn.Sel.Name] {
					t.Fatalf("%s calls forbidden constructor %s", name, fn.Sel.Name)
				}
			}
			return true
		})
	}
}

func TestCommandOwnsNoWriterSemantics(t *testing.T) {
	t.Parallel()
	files := parseCommandPackageFiles(t)
	forbiddenMethods := map[string]bool{
		"InspectOrReserveCallerMutation": true,
		"BeginControllerMutation":        true,
		"BeginEvidenceTransaction":       true,
		"SetResourceVersion":             true,
		"ApplyFinalized":                 true,
		"AcceptPreparedEvidence":         true,
		"StageCompletedResult":           true,
		"FinalizationPermit":             true,
		"Seal":                           true,
		"Commit":                         true,
		"Abort":                          true,
	}
	for name, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if forbiddenMethods[sel.Sel.Name] {
				t.Fatalf("%s calls writer-owned method %s", name, sel.Sel.Name)
			}
			return true
		})
	}
}

func TestPropertyCommandNeverImportsStateDirectly(t *testing.T) {
	t.Parallel()
	prop := func(_ uint8) bool {
		files := parseCommandPackageFiles(t)
		for _, f := range files {
			for _, imp := range f.Imports {
				path := strings.Trim(imp.Path.Value, `"`)
				if strings.Contains(path, "/govaccess/state") {
					return false
				}
			}
		}
		return true
	}
	if err := quick.Check(prop, &quick.Config{MaxCount: 8}); err != nil {
		t.Fatal(err)
	}
}

func parseCommandPackageFiles(t *testing.T) map[string]*ast.File {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller")
	}
	dir := filepath.Dir(thisFile)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	files := map[string]*ast.File{}
	for _, pkg := range pkgs {
		for name, f := range pkg.Files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			files[name] = f
		}
	}
	return files
}
