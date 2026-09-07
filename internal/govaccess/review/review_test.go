package review_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/review"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/roleassign"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

func issuePermit(t *testing.T, claim state.ParticipantClaim) state.MutationFinalizationPermit {
	t.Helper()
	store := state.NewStore(clock.NewFixedClock(time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)))
	for _, p := range claim.Versions().Predicates() {
		store.SetResourceVersion(p.ResourceUID, p.ResourceVersion)
	}
	cd, _ := evidence.NewCandidateDigest([]byte("cand-review"))
	depSet, _ := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), claim.Versions().Predicates())
	curr, _ := state.NewAuthorizationCurrentnessClaim(cd, depSet)
	adm, _ := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "actor", PrincipalType: model.PrincipalTypeHuman}
	key, _ := idempotency.NewIdempotencyLookupKey(actor, "accessreview.write", "ck-r")
	bind, _ := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "AccessReview", Name: "r", UID: "r"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		[]byte("digest"),
	)
	lease, _ := store.InspectOrReserveCallerMutation(key, bind).Owner()
	tx, err := lease.Begin(adm)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Abort()
	desc, _ := evidence.NewAuthorizationCandidateDescriptor(cd, evidence.AuthorizationAllow, "accessreview.write", "r")
	fin, _ := evidence.NewFinalizedAuthorizationCandidate(desc, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	req := evidence.RequireCurrentAllow(fin, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	allow, _ := req.AsAllow()
	if fail := tx.AcceptCurrentAllow(allow); fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	permit, fail := tx.FinalizationPermit(claim)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return permit
}

func sampleAssignmentSpec() model.RoleAssignmentSpec {
	nb := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	exp := nb.Add(2 * time.Hour)
	return model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{
			Kind: model.RoleHolderPrincipal,
			Principal: &model.PrincipalRef{
				Issuer: "https://idp.example", Subject: "u1", PrincipalType: model.PrincipalTypeHuman,
			},
		},
		RoleDefinitionRef:     apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "ops", UID: "rd-1"},
		RoleDefinitionVersion: "v1",
		Validity:              model.AssignmentValidityTimeBound,
		NotBefore:             &nb,
		ExpiresAt:             &exp,
	}
}

func TestReviewFinalizeAndConflictDetection(t *testing.T) {
	t.Parallel()
	reviewer := model.PrincipalRef{Issuer: "https://idp.example", Subject: "reviewer-1", PrincipalType: model.PrincipalTypeHuman}
	elig := model.EligibilityRef{Principal: reviewer}
	intent, fail := review.NewPreparedReviewIntent(
		"participant-r", "review-1", "review-one",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		model.AccessReviewSpec{
			AccessReviewRuleRef: apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "AccessReviewRule", Name: "rule", UID: "rule-1"},
			ReviewerEligibility: []model.EligibilityRef{elig},
		},
		"0", "1",
		[]model.PrincipalRef{reviewer},
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if _, ok := intent.FinalizeAt(issuePermit(t, intent.ParticipantClaim())).AsDomain(); !ok {
		t.Fatal("beneficiary conflict must deny self-certification")
	}
}

func TestReviewRemediationPortSubmissions(t *testing.T) {
	t.Parallel()
	service := review.NewRemediationService(roleassign.NewReviewRemediationPort())
	scope := apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}}
	spec := sampleAssignmentSpec()
	if intent, fail := service.SubmitRevoke("p", "ra-1", "ra-1", scope, spec, "0", "1"); fail.Reason() != "" || intent.Operation() != roleassign.OpReviewRevoke {
		t.Fatalf("revoke trigger failed: %s", fail.Reason())
	}
	if intent, fail := service.SubmitReplace("p", "ra-1", "ra-1", scope, spec, spec, "1", "2"); fail.Reason() != "" || intent.Operation() != roleassign.OpReviewReplace {
		t.Fatalf("replace trigger failed: %s", fail.Reason())
	}
	if intent, fail := service.SubmitAdvanceReviewDue("p", "ra-1", "ra-1", scope, spec, "2", "3", time.Date(2026, 9, 8, 10, 0, 0, 0, time.UTC)); fail.Reason() != "" || intent.Operation() != roleassign.OpAdvanceReviewDue {
		t.Fatalf("advance trigger failed: %s", fail.Reason())
	}
}

func TestReviewOnlyReviewRemediationPortCaller(t *testing.T) {
	t.Parallel()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	govaccessDir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), ".."))
	callers := map[string]int{}
	err := filepath.WalkDir(govaccessDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.AllErrors)
		if err != nil {
			return err
		}
		rel := filepath.ToSlash(strings.TrimPrefix(path, govaccessDir+"/"))
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if sel.Sel.Name == "NewReviewRemediationTriggerIntent" {
				callers[rel]++
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for file := range callers {
		if file == "roleassign/reviewport.go" {
			continue
		}
		if file != "review/remediation.go" {
			t.Fatalf("review remediation trigger caller outside review/remediation.go: %s", file)
		}
	}
}
