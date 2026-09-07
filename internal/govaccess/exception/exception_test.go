package exception_test

import (
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/approvalreq"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/clock"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/evidence"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/exception"
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
	cd, _ := evidence.NewCandidateDigest([]byte("cand-exception"))
	depSet, _ := state.NewAuthorizationDependencyVersionSet(store.SnapshotVersion(), claim.Versions().Predicates())
	curr, _ := state.NewAuthorizationCurrentnessClaim(cd, depSet)
	adm, _ := state.NewCallerMutationAdmission([]state.ParticipantClaim{claim}, curr)
	actor := model.PrincipalRef{Issuer: "https://idp.example", Subject: "actor", PrincipalType: model.PrincipalTypeHuman}
	key, _ := idempotency.NewIdempotencyLookupKey(actor, "exceptiongrant.write", "ck-e")
	bind, _ := idempotency.NewRequestBinding(
		apimeta.TypedRef{APIVersion: model.APIVersionGov, Kind: model.KindExceptionGrant, Name: "e", UID: "e"},
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		[]byte("digest"),
	)
	lease, _ := store.InspectOrReserveCallerMutation(key, bind).Owner()
	tx, err := lease.Begin(adm)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Abort()
	desc, _ := evidence.NewAuthorizationCandidateDescriptor(cd, evidence.AuthorizationAllow, "exceptiongrant.write", "e")
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

func decisionFacts(t *testing.T, result model.ExceptionTerminalResult) model.ExceptionDecisionFacts {
	t.Helper()
	facts, ok := model.NewExceptionDecisionFacts(
		"v1",
		"subject-1",
		apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
		time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC),
		time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC),
		"override-type",
		[]string{"c1"},
		[]byte("approval"),
		result,
		"terminal",
		"proposal-1",
		func() string {
			if result == model.ExceptionTerminalGrant {
				return "grant-1"
			}
			return ""
		}(),
	)
	if !ok {
		t.Fatal("decision facts")
	}
	return facts
}

func TestExceptionDeriverAndTerminalDescriptorMatrix(t *testing.T) {
	t.Parallel()
	var _ approvalreq.ExceptionRequirementDeriverPort = exception.ApprovalRequirementDeriver{}
	req, ok := exception.ApprovalRequirementDeriver{}.DeriveForException(approvalreq.ExceptionInput{})
	if !ok || req.Mode() != model.ApprovalRequirementRequired {
		t.Fatal("exception deriver must return Required")
	}

	for _, result := range []model.ExceptionTerminalResult{model.ExceptionTerminalGrant, model.ExceptionTerminalDeny} {
		intent, fail := exception.NewPreparedExceptionIntent(
			"participant-e",
			"exception-1",
			"exception-one",
			apimeta.ScopeRef{TypedRef: apimeta.TypedRef{APIVersion: "core.sovrunn.io/v1alpha1", Kind: "Project", Name: "p1", UID: "proj-1"}},
			model.ExceptionGrantRecord{
				Effect:     model.ExceptionGrantEffectGrant,
				ControlRef: apimeta.TypedRef{APIVersion: model.APIVersionGov, Kind: "Control", Name: "control", UID: "ctrl-1"},
				SubjectUID: "subject-1",
				NotBefore:  time.Date(2026, 9, 7, 10, 0, 0, 0, time.UTC),
				ExpiresAt:  time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC),
			},
			decisionFacts(t, result),
			"0",
			"1",
		)
		if fail.Reason() != "" {
			t.Fatal(fail.Reason())
		}
		out := intent.FinalizeAt(issuePermit(t, intent.ParticipantClaim()))
		change, ok := out.AsFinalized()
		if !ok || !change.Sealed() {
			t.Fatalf("expected finalized for %s", result)
		}
		if result == model.ExceptionTerminalGrant {
			if change.PublicationKind() != state.ResourceMutation || len(change.EvidenceDescriptors()) != 2 {
				t.Fatal("grant requires two mutation descriptors")
			}
		} else {
			if change.PublicationKind() != state.DomainConclusion {
				t.Fatal("deny must publish domain conclusion")
			}
			if _, ok := change.ConclusionDescriptor(); !ok {
				t.Fatal("deny conclusion descriptor required")
			}
		}
	}
}
