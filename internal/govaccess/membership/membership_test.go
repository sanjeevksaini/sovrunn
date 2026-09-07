package membership_test

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
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/membership"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

func issuePermit(t *testing.T, claim state.ParticipantClaim) state.MutationFinalizationPermit {
	t.Helper()
	store := state.NewStore(clock.NewFixedClock(time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)))
	for _, p := range claim.Versions().Predicates() {
		store.SetResourceVersion(p.ResourceUID, p.ResourceVersion)
	}
	cd, _ := evidence.NewCandidateDigest([]byte("cand-membership"))
	depSet, _ := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), claim.Versions().Predicates())
	curr, _ := state.NewAuthorizationCurrentnessClaim(cd, depSet)
	adm, _ := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "actor", PrincipalType: model.PrincipalTypeHuman}
	key, _ := idempotency.NewIdempotencyLookupKey(actor, "membership.write", "ck-m")
	bind, _ := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "Membership", Name: "m", UID: "m"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		[]byte("digest"),
	)
	lease, _ := store.InspectOrReserveCallerMutation(key, bind).Owner()
	tx, err := lease.Begin(adm)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Abort()
	desc, _ := evidence.NewAuthorizationCandidateDescriptor(cd, evidence.AuthorizationAllow, "membership.write", "m")
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

func TestMembershipDeriverAndFinalizeExpansion(t *testing.T) {
	t.Parallel()
	var _ approvalreq.MembershipRequirementDeriverPort = membership.ApprovalRequirementDeriver{}
	req, ok := membership.ApprovalRequirementDeriver{}.DeriveForMembership(approvalreq.MembershipInput{})
	if !ok || req.Mode() != model.ApprovalRequirementNotRequired {
		t.Fatal("membership deriver must be deterministic NotRequired")
	}

	scope := apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}}
	intent, fail := membership.NewPreparedMembershipIntent(
		"participant-m",
		"m-1",
		"membership-1",
		scope,
		model.MembershipSpec{
			PrincipalRef:   model.PrincipalRef{Issuer: "https://idp.example", Subject: "u1", PrincipalType: model.PrincipalTypeHuman},
			ContextKind:    model.MembershipContextAccessGroup,
			MembershipType: model.MembershipTypeStandard,
			AccessGroupRef: &model.AccessGroupRef{TypedRef: apimeta.TypedRef{
				APIVersion: model.APIVersionIAM, Kind: model.KindAccessGroup, Name: "ag-1", UID: "ag-1",
			}},
		},
		model.MembershipProtected{
			SourceAuthority: "idp-sync",
			ProvisionedAt:   time.Date(2026, 9, 7, 9, 0, 0, 0, time.UTC),
		},
		"0",
		"1",
		2,
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if intent.AssignmentEffectSize() != 2 {
		t.Fatal("assignment effect expansion should remain owner-local")
	}
	change, ok := intent.FinalizeAt(issuePermit(t, intent.ParticipantClaim())).AsFinalized()
	if !ok || !change.Sealed() {
		t.Fatal("expected finalized membership change")
	}
	if change.EvidenceDescriptors()[0].EventType() != "membership.published" {
		t.Fatal("membership event type")
	}
}

func TestMembershipNeverSubmitsRoleAssignmentTriggerIntent(t *testing.T) {
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
		if d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if !strings.Contains(filepath.ToSlash(path), "/membership/") {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.AllErrors)
		if err != nil {
			return err
		}
		for _, imp := range file.Imports {
			if strings.Contains(strings.Trim(imp.Path.Value, `"`), "/govaccess/roleassign") {
				t.Fatalf("membership must not import roleassign: %s", filepath.Base(path))
			}
		}
		ast.Inspect(file, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if !ok {
				return true
			}
			switch id.Name {
			case "ActivationTriggerIntent", "ReviewRemediationTriggerIntent", "PreparedRoleAssignmentIntent", "FinalizedRoleAssignmentChange":
				t.Fatalf("membership references roleassign trigger/type: %s in %s", id.Name, filepath.Base(path))
			}
			return true
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
