package cloudmodel_test

import (
	"testing"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

func pendingStatus() model.CloudProviderParticipationStatus {
	return cloudmodel.InitialParticipationStatus(time.Date(2026, 8, 12, 0, 0, 0, 0, time.UTC))
}

func mustApply(t *testing.T, cur model.CloudProviderParticipationStatus, action cloudmodel.ParticipationAction, party cloudmodel.HoldParty) model.CloudProviderParticipationStatus {
	t.Helper()
	next, prob := cloudmodel.ApplyParticipationAction(cur, action, party)
	if prob != nil {
		t.Fatalf("action %s: %#v", action, prob)
	}
	return next
}

func assertStateInvalid(t *testing.T, cur model.CloudProviderParticipationStatus, action cloudmodel.ParticipationAction, party cloudmodel.HoldParty) {
	t.Helper()
	before := cur
	next, prob := cloudmodel.ApplyParticipationAction(cur, action, party)
	if prob == nil {
		t.Fatalf("action %s: want CONFLICT", action)
	}
	if prob.Code != apiproblem.CodeConflict {
		t.Fatalf("code=%q want CONFLICT", prob.Code)
	}
	if len(prob.Violations) == 0 || prob.Violations[0].Code != cloudmodel.ViolationParticipationStateInvalid {
		t.Fatalf("violations=%#v", prob.Violations)
	}
	if next.Phase != "" || next.PlatformSuspended || next.ProviderSuspended {
		t.Fatalf("invalid transition must not return a writable status: %#v", next)
	}
	// Caller must not write; source remains unchanged.
	if before.Phase != cur.Phase ||
		before.PlatformSuspended != cur.PlatformSuspended ||
		before.ProviderSuspended != cur.ProviderSuspended ||
		before.RequestExpiresAt != cur.RequestExpiresAt {
		t.Fatalf("input status mutated: before=%#v after=%#v", before, cur)
	}
}

func TestLifecycle_PendingToActiveRejectedWithdrawnExpired(t *testing.T) {
	t.Parallel()

	active := mustApply(t, pendingStatus(), cloudmodel.ParticipationActionAccept, "")
	if active.Phase != model.ParticipationPhaseActive {
		t.Fatalf("accept phase=%q", active.Phase)
	}

	rejected := mustApply(t, pendingStatus(), cloudmodel.ParticipationActionReject, "")
	if rejected.Phase != model.ParticipationPhaseRejected {
		t.Fatalf("reject phase=%q", rejected.Phase)
	}

	withdrawn := mustApply(t, pendingStatus(), cloudmodel.ParticipationActionWithdraw, "")
	if withdrawn.Phase != model.ParticipationPhaseWithdrawn {
		t.Fatalf("withdraw phase=%q", withdrawn.Phase)
	}

	expired := mustApply(t, pendingStatus(), cloudmodel.ParticipationActionExpire, "")
	if expired.Phase != model.ParticipationPhaseExpired {
		t.Fatalf("expire phase=%q", expired.Phase)
	}
}

func TestLifecycle_IndependentHoldsAndDerivedPhase(t *testing.T) {
	t.Parallel()

	st := mustApply(t, pendingStatus(), cloudmodel.ParticipationActionAccept, "")
	if st.Phase != model.ParticipationPhaseActive {
		t.Fatalf("phase=%q", st.Phase)
	}

	st = mustApply(t, st, cloudmodel.ParticipationActionSuspend, cloudmodel.HoldPartyPlatform)
	if st.Phase != model.ParticipationPhaseSuspended || !st.PlatformSuspended || st.ProviderSuspended {
		t.Fatalf("platform suspend: %#v", st)
	}

	st = mustApply(t, st, cloudmodel.ParticipationActionSuspend, cloudmodel.HoldPartyProvider)
	if st.Phase != model.ParticipationPhaseSuspended || !st.PlatformSuspended || !st.ProviderSuspended {
		t.Fatalf("both holds: %#v", st)
	}

	// Clearing one hold never reactivates while the other remains true.
	st = mustApply(t, st, cloudmodel.ParticipationActionResume, cloudmodel.HoldPartyPlatform)
	if st.Phase != model.ParticipationPhaseSuspended || st.PlatformSuspended || !st.ProviderSuspended {
		t.Fatalf("clear platform hold only: %#v", st)
	}

	st = mustApply(t, st, cloudmodel.ParticipationActionResume, cloudmodel.HoldPartyProvider)
	if st.Phase != model.ParticipationPhaseActive || st.PlatformSuspended || st.ProviderSuspended {
		t.Fatalf("both holds cleared: %#v", st)
	}

	if got := cloudmodel.DerivedEffectivePhase(false, false); got != model.ParticipationPhaseActive {
		t.Fatalf("derived=%q", got)
	}
	if got := cloudmodel.DerivedEffectivePhase(true, false); got != model.ParticipationPhaseSuspended {
		t.Fatalf("derived=%q", got)
	}
}

func TestLifecycle_ActiveSuspendedToTerminatingTerminated(t *testing.T) {
	t.Parallel()

	st := mustApply(t, pendingStatus(), cloudmodel.ParticipationActionAccept, "")
	st = mustApply(t, st, cloudmodel.ParticipationActionRequestRelease, "")
	if st.Phase != model.ParticipationPhaseTerminating {
		t.Fatalf("request-release phase=%q", st.Phase)
	}
	st = mustApply(t, st, cloudmodel.ParticipationActionAcceptRelease, "")
	if st.Phase != model.ParticipationPhaseTerminated {
		t.Fatalf("accept-release phase=%q", st.Phase)
	}

	// From Suspended → Terminating → Terminated with holds retained.
	st = mustApply(t, pendingStatus(), cloudmodel.ParticipationActionAccept, "")
	st = mustApply(t, st, cloudmodel.ParticipationActionSuspend, cloudmodel.HoldPartyProvider)
	st = mustApply(t, st, cloudmodel.ParticipationActionRequestRelease, "")
	if st.Phase != model.ParticipationPhaseTerminating || !st.ProviderSuspended {
		t.Fatalf("terminating retains holds: %#v", st)
	}
	st = mustApply(t, st, cloudmodel.ParticipationActionAcceptRelease, "")
	if st.Phase != model.ParticipationPhaseTerminated || !st.ProviderSuspended {
		t.Fatalf("terminated retains holds: %#v", st)
	}
}

func TestLifecycle_DeclineReleaseToPriorEffectivePhase(t *testing.T) {
	t.Parallel()

	// Both holds false → Active after decline-release.
	st := mustApply(t, pendingStatus(), cloudmodel.ParticipationActionAccept, "")
	st = mustApply(t, st, cloudmodel.ParticipationActionRequestRelease, "")
	st = mustApply(t, st, cloudmodel.ParticipationActionDeclineRelease, "")
	if st.Phase != model.ParticipationPhaseActive {
		t.Fatalf("decline → Active: %#v", st)
	}

	// Either hold true → Suspended after decline-release.
	st = mustApply(t, pendingStatus(), cloudmodel.ParticipationActionAccept, "")
	st = mustApply(t, st, cloudmodel.ParticipationActionSuspend, cloudmodel.HoldPartyPlatform)
	st = mustApply(t, st, cloudmodel.ParticipationActionRequestRelease, "")
	st = mustApply(t, st, cloudmodel.ParticipationActionDeclineRelease, "")
	if st.Phase != model.ParticipationPhaseSuspended || !st.PlatformSuspended {
		t.Fatalf("decline → Suspended: %#v", st)
	}
}

func TestLifecycle_SuspendResumeDeniedInPendingTerminatingTerminal(t *testing.T) {
	t.Parallel()

	phases := []model.CloudProviderParticipationStatus{
		{Phase: model.ParticipationPhasePending},
		{Phase: model.ParticipationPhaseTerminating},
		{Phase: model.ParticipationPhaseRejected},
		{Phase: model.ParticipationPhaseWithdrawn},
		{Phase: model.ParticipationPhaseExpired},
		{Phase: model.ParticipationPhaseTerminated},
	}
	for _, st := range phases {
		assertStateInvalid(t, st, cloudmodel.ParticipationActionSuspend, cloudmodel.HoldPartyPlatform)
		assertStateInvalid(t, st, cloudmodel.ParticipationActionResume, cloudmodel.HoldPartyProvider)
	}
}

func TestLifecycle_InvalidSourceStateTransitions(t *testing.T) {
	t.Parallel()

	active := mustApply(t, pendingStatus(), cloudmodel.ParticipationActionAccept, "")
	assertStateInvalid(t, active, cloudmodel.ParticipationActionAccept, "")
	assertStateInvalid(t, active, cloudmodel.ParticipationActionReject, "")
	assertStateInvalid(t, active, cloudmodel.ParticipationActionWithdraw, "")
	assertStateInvalid(t, active, cloudmodel.ParticipationActionExpire, "")
	assertStateInvalid(t, active, cloudmodel.ParticipationActionAcceptRelease, "")
	assertStateInvalid(t, active, cloudmodel.ParticipationActionDeclineRelease, "")

	pending := pendingStatus()
	assertStateInvalid(t, pending, cloudmodel.ParticipationActionRequestRelease, "")
	assertStateInvalid(t, pending, cloudmodel.ParticipationActionAcceptRelease, "")
	assertStateInvalid(t, pending, cloudmodel.ParticipationActionDeclineRelease, "")

	terminating := mustApply(t, active, cloudmodel.ParticipationActionRequestRelease, "")
	assertStateInvalid(t, terminating, cloudmodel.ParticipationActionAccept, "")
	assertStateInvalid(t, terminating, cloudmodel.ParticipationActionRequestRelease, "")
}

func TestLifecycle_ZeroOrMultipleSuspendResumeGrantRejection(t *testing.T) {
	t.Parallel()

	platformUID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	providerUID := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

	_, prob := cloudmodel.ResolveSuspendResumeParty(
		cloudmodel.ParticipationActionSuspend, platformUID, providerUID, nil,
	)
	if prob == nil || prob.Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("zero grants: %#v", prob)
	}

	_, prob = cloudmodel.ResolveSuspendResumeParty(
		cloudmodel.ParticipationActionResume, platformUID, providerUID,
		[]cloudmodel.ScopedActionGrant{
			{Action: cloudmodel.ActionResumePlatform, ScopeUID: platformUID},
			{Action: cloudmodel.ActionResumeProvider, ScopeUID: providerUID},
		},
	)
	if prob == nil || prob.Code != apiproblem.CodeAuthorizationDenied {
		t.Fatalf("multiple grants: %#v", prob)
	}

	party, prob := cloudmodel.ResolveSuspendResumeParty(
		cloudmodel.ParticipationActionSuspend, platformUID, providerUID,
		[]cloudmodel.ScopedActionGrant{
			{Action: cloudmodel.ActionSuspendPlatform, ScopeUID: platformUID},
			{Action: cloudmodel.ActionSuspendProvider, ScopeUID: "cccccccccccccccccccccccccccccccc"}, // wrong scope
		},
	)
	if prob != nil || party != cloudmodel.HoldPartyPlatform {
		t.Fatalf("exactly one match: party=%q prob=%#v", party, prob)
	}
}

func TestLifecycle_InitialStatus(t *testing.T) {
	t.Parallel()
	created := time.Date(2026, 8, 12, 15, 4, 5, 0, time.UTC)
	st := cloudmodel.InitialParticipationStatus(created)
	if st.Phase != model.ParticipationPhasePending {
		t.Fatalf("phase=%q", st.Phase)
	}
	if st.PlatformSuspended || st.ProviderSuspended {
		t.Fatalf("holds must be false: %#v", st)
	}
	want := created.Add(7 * 24 * time.Hour).Format(time.RFC3339)
	if st.RequestExpiresAt != want {
		t.Fatalf("requestExpiresAt=%q want %q", st.RequestExpiresAt, want)
	}
}
