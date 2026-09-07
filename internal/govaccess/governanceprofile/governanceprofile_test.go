package governanceprofile_test

import (
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/governanceprofile"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

func issuePermit(t *testing.T, claim state.ParticipantClaim) state.MutationFinalizationPermit {
	t.Helper()
	store := state.NewStore(clock.NewFixedClock(time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)))
	for _, p := range claim.Versions().Predicates() {
		store.SetResourceVersion(p.ResourceUID, p.ResourceVersion)
	}
	cd, _ := evidence.NewCandidateDigest([]byte("cand-gp"))
	depSet, _ := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), claim.Versions().Predicates())
	curr, _ := state.NewAuthorizationCurrentnessClaim(cd, depSet)
	adm, _ := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "actor", PrincipalType: model.PrincipalTypeHuman}
	key, _ := idempotency.NewIdempotencyLookupKey(actor, "governanceprofile.write", "ck-g")
	bind, _ := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: model.APIVersionGov, Kind: "GovernanceProfile", Name: "gp", UID: "gp"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Organization", Name: "org", UID: "org-1"}},
		[]byte("digest"),
	)
	lease, _ := store.InspectOrReserveCallerMutation(key, bind).Owner()
	tx, err := lease.Begin(adm)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Abort()
	desc, _ := evidence.NewAuthorizationCandidateDescriptor(cd, evidence.AuthorizationAllow, "governanceprofile.write", "gp")
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

func TestGovernanceProfileFinalize(t *testing.T) {
	t.Parallel()
	intent, fail := governanceprofile.NewPreparedGovernanceProfileIntent(
		"participant-g",
		"gp-1",
		"governanceprofile-one",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Organization", Name: "org", UID: "org-1"}},
		model.GovernanceProfileSpec{
			Version:          "v1",
			PublicationState: model.PublicationStatePublished,
			ApprovalPolicyRefs: []apimeta.TypedRef{
				{APIVersion: model.APIVersionIAM, Kind: "ApprovalPolicy", Name: "ap-1", UID: "ap-1"},
			},
		},
		"0",
		"1",
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if !intent.Sealed() {
		t.Fatal("intent must be sealed")
	}
	change, ok := intent.FinalizeAt(issuePermit(t, intent.ParticipantClaim())).AsFinalized()
	if !ok || !change.Sealed() || change.EvidenceDescriptors()[0].EventType() != "governanceprofile.published" {
		t.Fatal("governanceprofile finalization failed")
	}
}
