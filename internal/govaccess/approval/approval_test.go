package approval_test

import (
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/approval"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/exception"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/idempotency"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/membership"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/privileged"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

func issuePermit(t *testing.T, claim state.ParticipantClaim) state.MutationFinalizationPermit {
	t.Helper()
	store := state.NewStore(clock.NewFixedClock(time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC)))
	for _, p := range claim.Versions().Predicates() {
		store.SetResourceVersion(p.ResourceUID, p.ResourceVersion)
	}
	cd, _ := evidence.NewCandidateDigest([]byte("cand-approval"))
	depSet, _ := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), claim.Versions().Predicates())
	curr, _ := state.NewAuthorizationCurrentnessClaim(cd, depSet)
	adm, _ := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "actor", PrincipalType: model.PrincipalTypeHuman}
	key, _ := idempotency.NewIdempotencyLookupKey(actor, "approvalrequest.write", "ck-apr")
	bind, _ := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "ApprovalRequest", Name: "ar", UID: "ar"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		[]byte("digest"),
	)
	lease, _ := store.InspectOrReserveCallerMutation(key, bind).Owner()
	tx, err := lease.Begin(adm)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Abort()
	desc, _ := evidence.NewAuthorizationCandidateDescriptor(cd, evidence.AuthorizationAllow, "approvalrequest.write", "ar")
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

func TestApprovalConsumesOwnerDeriversAndSoD(t *testing.T) {
	t.Parallel()
	reqs, fail := approval.ConsumeDerivers(
		membership.ApprovalRequirementDeriver{},
		privileged.ApprovalRequirementDeriver{},
		exception.ApprovalRequirementDeriver{},
	)
	if fail.Reason() != "" || len(reqs) != 3 {
		t.Fatalf("consume derivers failed: %s", fail.Reason())
	}
	requester := model.PrincipalRef{Issuer: "https://idp.example", Subject: "r1", PrincipalType: model.PrincipalTypeHuman}
	decider := model.PrincipalRef{Issuer: "https://idp.example", Subject: "d1", PrincipalType: model.PrincipalTypeHuman}
	if !approval.ValidateStagesQuorumSoD(3, 2, requester, decider) {
		t.Fatal("expected valid quorum and SoD")
	}
	if approval.ValidateStagesQuorumSoD(1, 2, requester, decider) {
		t.Fatal("quorum cannot exceed stage population")
	}
	if approval.ValidateStagesQuorumSoD(2, 1, requester, requester) {
		t.Fatal("SoD should fail when requester equals decider")
	}
}

func TestApprovalFinalizeRespectsTerminalImmutability(t *testing.T) {
	t.Parallel()
	approver := model.PrincipalRef{Issuer: "https://idp.example", Subject: "approver-1", PrincipalType: model.PrincipalTypeHuman}
	elig := model.EligibilityRef{Principal: approver}
	intent, fail := approval.NewPreparedApprovalIntent(
		"participant-ap",
		"approvalreq-1",
		"approvalreq-one",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		model.ApprovalPolicy{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "ApprovalPolicy"},
			Metadata: apimeta.ObjectMeta{Name: "ap", UID: "ap-1"},
			Spec: model.ApprovalPolicySpec{
				Version:             "v1",
				PublicationState:    model.PublicationStatePublished,
				ApproverEligibility: []model.EligibilityRef{elig},
			},
		},
		model.ApprovalRequest{
			TypeMeta: apimeta.TypeMeta{APIVersion: model.APIVersionIAM, Kind: "ApprovalRequest"},
			Metadata: apimeta.ObjectMeta{Name: "ar", UID: "approvalreq-1"},
			Spec:     model.ApprovalRequestSpec{ApprovalPolicyRef: apimeta.TypedRef{APIVersion: model.APIVersionIAM, Kind: "ApprovalPolicy", Name: "ap", UID: "ap-1"}, SubjectKind: "PrivilegedAccessRequest", SubjectUID: "par-1"},
			Status:   model.ApprovalRequestStatus{Phase: "Decided", Decision: model.ApprovalDecisionApprove},
		},
		"0",
		"1",
	)
	if fail.Reason() != "" {
		t.Fatal(fail.Reason())
	}
	if _, ok := intent.FinalizeAt(issuePermit(t, intent.ParticipantClaim())).AsDomain(); !ok {
		t.Fatal("terminal approval decision must be immutable")
	}
}
