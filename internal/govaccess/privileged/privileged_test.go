package privileged_test

import (
	"context"
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
	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/approvalreq"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/policyseam"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/privileged"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/roleassign"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
	"github.com/sanjeevksaini/sovrunn/internal/policyeval"
)

type fakeBoundary struct {
	outcome policyeval.Outcome
}

func (f fakeBoundary) Evaluate(context.Context, policyeval.PolicyEvaluationRequest) (policyeval.PolicyEvaluationResult, decision.TimingEnvelope, policyeval.NonResult, error) {
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)
	return policyeval.PolicyEvaluationResult{
		Outcome: f.outcome, InputDigest: strings.Repeat("a", 64), EvaluatedAt: now,
	}, decision.TimingEnvelope{StartedAt: now, CompletedAt: now}, "", nil
}

func issuePermit(t *testing.T, claim state.ParticipantClaim) state.MutationFinalizationPermit {
	t.Helper()
	store := state.NewStore(clock.NewFixedClock(time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)))
	for _, p := range claim.Versions().Predicates() {
		store.SetResourceVersion(p.ResourceUID, p.ResourceVersion)
	}
	cd, _ := evidence.NewCandidateDigest([]byte("cand-privileged"))
	depSet, _ := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), claim.Versions().Predicates())
	curr, _ := state.NewAuthorizationCurrentnessClaim(cd, depSet)
	adm, _ := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "actor", PrincipalType: model.PrincipalTypeHuman}
	key, _ := idempotency.NewIdempotencyLookupKey(actor, "privilegedaccessrequest.write", "ck-p")
	bind, _ := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "PrivilegedAccessRequest", Name: "p", UID: "p"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		[]byte("digest"),
	)
	lease, _ := store.InspectOrReserveCallerMutation(key, bind).Owner()
	tx, err := lease.Begin(adm)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Abort()
	desc, _ := evidence.NewAuthorizationCandidateDescriptor(cd, evidence.AuthorizationAllow, "privilegedaccessrequest.write", "p")
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

func TestPrivilegedDeriverAndFinalizeAALFloor(t *testing.T) {
	t.Parallel()
	var _ approvalreq.PrivilegedRequirementDeriverPort = privileged.ApprovalRequirementDeriver{}
	req, ok := privileged.ApprovalRequirementDeriver{}.DeriveForPrivileged(approvalreq.PrivilegedInput{})
	if !ok || req.Mode() != model.ApprovalRequirementRequired {
		t.Fatal("privileged deriver must return Required")
	}

	p, fail := policyseam.NewPort(fakeBoundary{outcome: policyeval.OutcomeAllow})
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	now := time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)
	intent, fail := privileged.NewPreparedPrivilegedIntent(
		"participant-p",
		"par-1",
		"par-one",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		model.PrivilegedAccessRequestSpec{
			Mode: model.PrivilegedAccessModeJIT,
			RequesterRef: model.PrincipalRef{
				Issuer: "https://idp.example", Subject: "requester", PrincipalType: model.PrincipalTypeHuman,
			},
			RoleDefinitionRef:  apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "ops", UID: "rd-1"},
			ActivationDeadline: now.Add(2 * time.Hour),
			RequestedDuration:  "PT2H",
		},
		model.AssuranceEvidence{
			PrincipalRef:      model.PrincipalRef{Issuer: "https://idp.example", Subject: "requester", PrincipalType: model.PrincipalTypeHuman},
			AuthenticatedAt:   now,
			Level:             model.AssuranceLevelAAL2,
			PhishingResistant: true,
			SourceAuthority:   "authn",
		},
		p,
		"0",
		"1",
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if !intent.Sealed() || intent.HasPublicationDerivedFields() {
		t.Fatal("intent must be sealed and publication-time-free")
	}
	change, ok := intent.FinalizeAt(issuePermit(t, intent.ParticipantClaim())).AsFinalized()
	if !ok || !change.Sealed() {
		t.Fatal("expected finalized privileged change")
	}
	if len(change.EvidenceDescriptors()) != 1 || change.EvidenceDescriptors()[0].EventType() != "privilegedaccessrequest.published" {
		t.Fatal("descriptor mismatch")
	}
}

func TestActivationPortSubmission(t *testing.T) {
	t.Parallel()
	service := privileged.NewActivationService(roleassign.NewActivationPort())
	spec := model.RoleAssignmentSpec{
		RoleHolderRef: model.RoleHolderRef{
			Kind:      model.RoleHolderPrincipal,
			Principal: &model.PrincipalRef{Issuer: "https://idp.example", Subject: "u", PrincipalType: model.PrincipalTypeHuman},
		},
		RoleDefinitionRef:     apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "ops", UID: "rd-1"},
		RoleDefinitionVersion: "v1",
		Validity:              model.AssignmentValidityTimeBound,
		NotBefore:             ptrTime(time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)),
		ExpiresAt:             ptrTime(time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)),
	}
	intent, fail := service.SubmitActivate(
		"participant-a", "ra-1", "ra-1",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		spec,
		"0",
		"1",
	)
	if fail.Reason() != "" || intent.Operation() != roleassign.OpActivate || intent.TriggerPort() != "activation" {
		t.Fatalf("activate trigger failed: %s", fail.Reason())
	}
	rev, fail := service.SubmitPrivilegedRevoke(
		"participant-a", "ra-1", "ra-1",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		"1",
		"2",
	)
	if fail.Reason() != "" || rev.Operation() != roleassign.OpPrivilegedRevoke {
		t.Fatalf("revoke trigger failed: %s", fail.Reason())
	}
}

func ptrTime(t time.Time) *time.Time { return &t }

func TestPrivilegedOnlyPolicySeamCallerAndActivationPortCaller(t *testing.T) {
	t.Parallel()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	govaccessDir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), ".."))
	policySeamImports := map[string]int{}
	activationCallerFiles := map[string]int{}
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
		for _, imp := range file.Imports {
			if strings.Contains(strings.Trim(imp.Path.Value, `"`), "/govaccess/policyseam") {
				policySeamImports[rel]++
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
			if sel.Sel.Name == "NewActivationTriggerIntent" {
				activationCallerFiles[rel]++
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for file := range policySeamImports {
		if !strings.HasPrefix(file, "privileged/") {
			t.Fatalf("policyseam import outside privileged: %s", file)
		}
	}
	if policySeamImports["privileged/intent.go"] == 0 {
		t.Fatal("privileged must consume policyseam")
	}
	for file := range activationCallerFiles {
		if file == "roleassign/activationport.go" {
			continue
		}
		if file != "privileged/activation.go" {
			t.Fatalf("activation trigger caller outside privileged/activation.go: %s", file)
		}
	}
}

func TestPrivilegedConstructsNoRoleAssignmentTypes(t *testing.T) {
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
					t.Fatalf("privileged directly constructs roleassign type in %s", filepath.Base(fileName))
				}
				return true
			})
		}
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
