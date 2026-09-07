package grantintent_test

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
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/approvalreq"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/grantintent"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/roleassign"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

func testProposal(t *testing.T, privileged bool) model.RoleAssignmentProposal {
	t.Helper()
	principal := model.PrincipalRef{Issuer: "https://idp.example", Subject: "human-1", PrincipalType: model.PrincipalTypeHuman}
	spec := model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{Kind: model.RoleHolderPrincipal, Principal: &principal},
		RoleDefinitionRef: apimeta.TypedRef{
			APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "admin", UID: "rd-1",
		},
		RoleDefinitionVersion: "v1",
		Validity:              model.AssignmentValidityStanding,
	}
	var policy *apimeta.TypedRef
	if privileged {
		nb := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
		exp := nb.Add(2 * time.Hour)
		spec.Validity = model.AssignmentValidityTimeBound
		spec.NotBefore = &nb
		spec.ExpiresAt = &exp
		ref := apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "ApprovalPolicy", Name: "jit", UID: "ap-1"}
		policy = &ref
	}
	proposal, ok := model.NewRoleAssignmentProposal(
		"proposal-1", "ra-1", "ra-one",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
			APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "proj", UID: "scope-1",
		}},
		spec, "1", "2", privileged, policy, nil,
	)
	if !ok {
		t.Fatal("proposal creation failed")
	}
	return proposal
}

func TestRoleGrantRequirementDeriver(t *testing.T) {
	t.Parallel()
	deriver := grantintent.RoleGrantRequirementDeriver{}

	ordinary := testProposal(t, false)
	ordinaryReq, ok := deriver.DeriveForRoleGrant(approvalreq.RoleGrantInput{Proposal: ordinary, Privileged: false})
	if !ok || ordinaryReq.Mode() != model.ApprovalRequirementNotRequired {
		t.Fatal("ordinary grant should derive NotRequired")
	}

	privileged := testProposal(t, true)
	policy, _ := privileged.PolicyRef()
	privReq, ok := deriver.DeriveForRoleGrant(approvalreq.RoleGrantInput{
		Proposal: privileged, Privileged: true, PolicyRef: &policy,
	})
	if !ok || privReq.Mode() != model.ApprovalRequirementRequired {
		t.Fatal("privileged grant should derive Required")
	}
}

func TestAdmissionDirectGrantAndProposal(t *testing.T) {
	t.Parallel()
	deriver := grantintent.RoleGrantRequirementDeriver{}
	proposal := testProposal(t, false)
	admitted, fail := grantintent.AdmitRoleAssignmentProposal(state.ParticipantID("participant-1"), proposal, deriver)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if !admitted.Sealed() || admitted.Proposal().ProposalUID() == "" {
		t.Fatal("expected sealed admitted proposal")
	}

	spec := proposal.Spec()
	direct, fail := grantintent.AdmitDirectGrant(
		state.ParticipantID("participant-2"),
		grantintent.DirectGrant{
			ResourceUID:     proposal.ResourceUID(),
			ResourceName:    proposal.ResourceName(),
			ScopeRef:        proposal.ScopeRef(),
			Spec:            spec,
			ExpectedVersion: proposal.ExpectedVersion(),
			NextVersion:     proposal.NextVersion(),
			Privileged:      false,
		},
		deriver,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if direct.Requirement().Mode() != model.ApprovalRequirementNotRequired {
		t.Fatal("direct ordinary grant must derive NotRequired")
	}
}

func TestSubmissionUsesGrantPortWithImmutableTrigger(t *testing.T) {
	t.Parallel()
	deriver := grantintent.RoleGrantRequirementDeriver{}
	admitted, fail := grantintent.AdmitRoleAssignmentProposal(state.ParticipantID("participant-3"), testProposal(t, false), deriver)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	intent, fail := grantintent.SubmitToGrantPort(roleassign.NewGrantPort(), admitted)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if !intent.Sealed() || intent.TriggerPort() != "grant" || intent.Operation() != roleassign.OpGrant {
		t.Fatalf("unexpected roleassign intent: %#v", intent)
	}
}

func TestGrantintentConstructsNoRoleassignPreparedOrFinalizedTypes(t *testing.T) {
	t.Parallel()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	dir := filepath.Dir(thisFile)
	pkgs, err := parser.ParseDir(token.NewFileSet(), dir, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, pkg := range pkgs {
		for fileName, f := range pkg.Files {
			if strings.HasSuffix(fileName, "_test.go") {
				continue
			}
			ast.Inspect(f, func(n ast.Node) bool {
				cl, ok := n.(*ast.CompositeLit)
				if !ok {
					return true
				}
				switch compositeTypeName(cl.Type) {
				case "PreparedRoleAssignmentIntent", "FinalizedRoleAssignmentChange":
					t.Fatalf("forbidden direct construction in %s", filepath.Base(fileName))
				}
				return true
			})
		}
	}
}

func TestOnlyGrantintentInvokesGrantPortAndRoleassignNeverDerivesApproval(t *testing.T) {
	t.Parallel()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	govaccessDir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), ".."))
	err := filepath.WalkDir(govaccessDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.AllErrors)
		if err != nil {
			return err
		}
		rel := filepath.ToSlash(strings.TrimPrefix(path, filepath.ToSlash(govaccessDir)+"/"))
		for _, imp := range file.Imports {
			ip := strings.Trim(imp.Path.Value, `"`)
			if strings.Contains(rel, "/roleassign/") && strings.Contains(ip, "/govaccess/approvalreq") {
				t.Fatalf("roleassign must not import approvalreq: %s", rel)
			}
		}
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			if sel.Sel.Name == "Submit" {
				if x, ok := sel.X.(*ast.Ident); ok && x.Name == "port" {
					if !strings.Contains(rel, "grantintent/submission.go") {
						t.Fatalf("GrantPort submit call outside grantintent: %s", rel)
					}
				}
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func compositeTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return t.Sel.Name
	default:
		return ""
	}
}
