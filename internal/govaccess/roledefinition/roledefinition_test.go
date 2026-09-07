package roledefinition_test

import (
	"strings"
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/roledefinition"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

func issuePermit(t *testing.T, claim state.ParticipantClaim) state.MutationFinalizationPermit {
	t.Helper()
	store := state.NewStore(clock.NewFixedClock(time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)))
	for _, p := range claim.Versions().Predicates() {
		store.SetResourceVersion(p.ResourceUID, p.ResourceVersion)
	}
	cd, _ := evidence.NewCandidateDigest([]byte("cand-roledefinition"))
	depSet, _ := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), claim.Versions().Predicates())
	curr, _ := state.NewAuthorizationCurrentnessClaim(cd, depSet)
	adm, _ := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "actor", PrincipalType: model.PrincipalTypeHuman}
	key, _ := idempotency.NewIdempotencyLookupKey(actor, "roledefinition.write", "ck-rd")
	bind, _ := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "rd", UID: "rd"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		[]byte("digest"),
	)
	lease, _ := store.InspectOrReserveCallerMutation(key, bind).Owner()
	tx, err := lease.Begin(adm)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Abort()
	desc, _ := evidence.NewAuthorizationCandidateDescriptor(cd, evidence.AuthorizationAllow, "roledefinition.write", "rd")
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

func TestRoleDefinitionFinalizeAndClassification(t *testing.T) {
	t.Parallel()
	scope := apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Organization", Name: "org", UID: "org-1"}}
	intent, fail := roledefinition.NewPreparedRoleDefinitionIntent(
		"participant-rd", "rd-1", "ops-admin", scope,
		model.RoleDefinitionSpec{
			Version:          "v1",
			Actions:          []string{"roleassignment.grant"},
			PublicationState: model.PublicationStatePublished,
			SupersedesRef: &apimeta.TypedRef{
				APIVersion: model.APIVersionIAM, Kind: "RoleDefinition", Name: "ops-admin-v0", UID: "rd-0",
			},
		},
		"0", "1",
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	permit := issuePermit(t, intent.ParticipantClaim())
	out := intent.FinalizeAt(permit)
	change, ok := out.AsFinalized()
	if !ok {
		t.Fatalf("expected finalized, got %s", out.KindName())
	}
	if len(change.EvidenceDescriptors()) != 1 || change.EvidenceDescriptors()[0].EventType() != "roledefinition.published" {
		t.Fatal("expected roledefinition event descriptor")
	}
	result, ok := change.CanonicalResult()
	if !ok {
		t.Fatal("result required")
	}
	payload := string(result.Bytes())
	if payload == "" {
		t.Fatal("result payload empty")
	}
	if want := `"baseClassification":"Privileged"`; !contains(payload, want) {
		t.Fatalf("expected effective classification %s in payload", want)
	}
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
