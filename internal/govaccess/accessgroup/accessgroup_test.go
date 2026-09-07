package accessgroup_test

import (
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/accessgroup"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
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
	cd, fail := evidence.NewCandidateDigest([]byte("cand-accessgroup"))
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	depSet, err := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), claim.Versions().Predicates())
	if err != nil {
		t.Fatal(err)
	}
	curr, err := state.NewAuthorizationCurrentnessClaim(cd, depSet)
	if err != nil {
		t.Fatal(err)
	}
	adm, err := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	if err != nil {
		t.Fatal(err)
	}
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "actor", PrincipalType: model.PrincipalTypeHuman}
	key, fail := idempotency.NewIdempotencyLookupKey(actor, "accessgroup.write", "ck-ag")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	bind, fail := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: model.KindAccessGroup, Name: "ag-1", UID: "ag-1"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		[]byte("digest"),
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	out := store.InspectOrReserveCallerMutation(key, bind)
	lease, ok := out.Owner()
	if !ok {
		t.Fatal("owner lease expected")
	}
	tx, err := lease.Begin(adm)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Abort()
	desc, fail := evidence.NewAuthorizationCandidateDescriptor(cd, evidence.AuthorizationAllow, "accessgroup.write", "ag-1")
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	fin, fail := evidence.NewFinalizedAuthorizationCandidate(desc, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	req := evidence.RequireCurrentAllow(fin, tx.AuthorizationAdmissionProof(), tx.PublicationContext())
	allow, ok := req.AsAllow()
	if !ok {
		t.Fatal("allow proof expected")
	}
	if fail := tx.AcceptCurrentAllow(allow); fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	permit, fail := tx.FinalizationPermit(claim)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	return permit
}

func TestPreparedAndFinalizeAt(t *testing.T) {
	t.Parallel()
	scope := apimeta.ScopeRef{TypedRef: apimeta.TypedRef{
		APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "proj", UID: "proj-1",
	}}
	intent, fail := accessgroup.NewPreparedAccessGroupIntent(
		"participant-1", "ag-1", "ag-one", scope,
		model.AccessGroupSpec{
			OwnerRef: model.PrincipalRef{Issuer: "https://idp.example", Subject: "owner-1", PrincipalType: model.PrincipalTypeHuman},
		},
		"0", "1",
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if !intent.Sealed() || intent.HasPublicationDerivedFields() {
		t.Fatal("intent must be sealed and publication-time-free")
	}
	permit := issuePermit(t, intent.ParticipantClaim())
	out := intent.FinalizeAt(permit)
	change, ok := out.AsFinalized()
	if !ok {
		t.Fatalf("expected finalized, got %s", out.KindName())
	}
	if !change.Sealed() {
		t.Fatal("change must be sealed")
	}
	if !change.PermitBinding().Equal(permit.Binding()) {
		t.Fatal("permit binding mismatch")
	}
	want, fail := evidence.BindPublicationContext(permit.PublicationContext())
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if !change.ContextBinding().Equal(want) {
		t.Fatal("context binding mismatch")
	}
	if len(change.EvidenceDescriptors()) != 1 || change.EvidenceDescriptors()[0].EventType() != "accessgroup.published" {
		t.Fatal("expected one accessgroup descriptor")
	}
}
