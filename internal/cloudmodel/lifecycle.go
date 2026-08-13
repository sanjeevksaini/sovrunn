package cloudmodel

import (
	"sync"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/cloudmodel/model"
)

// ViolationParticipationStateInvalid is the VS0-STATE-001 invalid-transition
// violation code (violations[].code only; design §5.2).
const ViolationParticipationStateInvalid apiproblem.ViolationCode = "VS0_PARTICIPATION_STATE_INVALID"

// ParticipationAction is a VS0-STATE-001 lifecycle action, including the
// internal scheduler expiry transition (not a public route).
type ParticipationAction string

const (
	ParticipationActionAccept         ParticipationAction = "accept"
	ParticipationActionReject         ParticipationAction = "reject"
	ParticipationActionWithdraw       ParticipationAction = "withdraw"
	ParticipationActionExpire         ParticipationAction = "expire"
	ParticipationActionSuspend        ParticipationAction = "suspend"
	ParticipationActionResume         ParticipationAction = "resume"
	ParticipationActionRequestRelease ParticipationAction = "request-release"
	ParticipationActionAcceptRelease  ParticipationAction = "accept-release"
	ParticipationActionDeclineRelease ParticipationAction = "decline-release"
)

// HoldParty identifies which independent suspension hold an actor may change.
type HoldParty string

const (
	HoldPartyPlatform HoldParty = "platform"
	HoldPartyProvider HoldParty = "provider"
)

// Exact server-resolved suspend/resume action names (design §6.1).
const (
	ActionSuspendPlatform = "participation.suspend.platform"
	ActionSuspendProvider = "participation.suspend.provider"
	ActionResumePlatform  = "participation.resume.platform"
	ActionResumeProvider  = "participation.resume.provider"
)

// ScopedActionGrant is a minimal server-resolved grant descriptor used to
// resolve exactly one suspend/resume hold actor. Task 9 owns full grant
// resolution; this type is the lifecycle matching surface.
type ScopedActionGrant struct {
	Action   string
	ScopeUID string
}

// ParticipationActionGuard serializes participation item-actions and the
// deterministic scheduler expiry path (design §4.6). Handlers and the
// scheduler must share one instance.
type ParticipationActionGuard struct {
	mu sync.Mutex
}

// Lock acquires the participation action concurrency guard.
func (g *ParticipationActionGuard) Lock() {
	if g == nil {
		return
	}
	g.mu.Lock()
}

// Unlock releases the participation action concurrency guard.
func (g *ParticipationActionGuard) Unlock() {
	if g == nil {
		return
	}
	g.mu.Unlock()
}

// InitialParticipationStatus returns the create-time status for a new
// participation: Pending, both holds false, requestExpiresAt = createdAt+7d.
func InitialParticipationStatus(createdAt time.Time) model.CloudProviderParticipationStatus {
	createdAt = createdAt.UTC()
	return model.CloudProviderParticipationStatus{
		Phase:             model.ParticipationPhasePending,
		PlatformSuspended: false,
		ProviderSuspended: false,
		RequestExpiresAt:  createdAt.Add(7 * 24 * time.Hour).Format(time.RFC3339),
	}
}

// DerivedEffectivePhase returns Active when both holds are false and
// Suspended when either hold is true (design §4.5).
func DerivedEffectivePhase(platformSuspended, providerSuspended bool) model.ParticipationPhase {
	if platformSuspended || providerSuspended {
		return model.ParticipationPhaseSuspended
	}
	return model.ParticipationPhaseActive
}

// ResolveSuspendResumeParty resolves exactly one matching current scoped
// suspend/resume grant. Zero or multiple matches fail closed with
// AUTHORIZATION_DENIED (design §6.1).
func ResolveSuspendResumeParty(
	action ParticipationAction,
	platformUID, providerUID string,
	grants []ScopedActionGrant,
) (HoldParty, *apiproblem.Problem) {
	var wantPlatform, wantProvider string
	switch action {
	case ParticipationActionSuspend:
		wantPlatform = ActionSuspendPlatform
		wantProvider = ActionSuspendProvider
	case ParticipationActionResume:
		wantPlatform = ActionResumePlatform
		wantProvider = ActionResumeProvider
	default:
		return "", participationStateInvalid("suspend/resume grant resolution requires suspend or resume")
	}

	var matches []HoldParty
	for _, g := range grants {
		switch {
		case g.Action == wantPlatform && g.ScopeUID == platformUID:
			matches = append(matches, HoldPartyPlatform)
		case g.Action == wantProvider && g.ScopeUID == providerUID:
			matches = append(matches, HoldPartyProvider)
		}
	}
	if len(matches) != 1 {
		return "", apiproblem.New(apiproblem.CodeAuthorizationDenied).WithDetail(
			"exactly one matching suspend/resume grant is required",
		)
	}
	return matches[0], nil
}

// ApplyParticipationAction validates the source state and returns the next
// status for a VS0-STATE-001 transition. On failure it returns
// CONFLICT/VS0_PARTICIPATION_STATE_INVALID and the caller must not write status.
//
// party is required for suspend/resume and ignored for other actions.
func ApplyParticipationAction(
	current model.CloudProviderParticipationStatus,
	action ParticipationAction,
	party HoldParty,
) (model.CloudProviderParticipationStatus, *apiproblem.Problem) {
	next := current
	switch action {
	case ParticipationActionAccept:
		if current.Phase != model.ParticipationPhasePending {
			return model.CloudProviderParticipationStatus{}, participationStateInvalid("accept requires Pending")
		}
		next.Phase = model.ParticipationPhaseActive
		return next, nil

	case ParticipationActionReject:
		if current.Phase != model.ParticipationPhasePending {
			return model.CloudProviderParticipationStatus{}, participationStateInvalid("reject requires Pending")
		}
		next.Phase = model.ParticipationPhaseRejected
		return next, nil

	case ParticipationActionWithdraw:
		if current.Phase != model.ParticipationPhasePending {
			return model.CloudProviderParticipationStatus{}, participationStateInvalid("withdraw requires Pending")
		}
		next.Phase = model.ParticipationPhaseWithdrawn
		return next, nil

	case ParticipationActionExpire:
		if current.Phase != model.ParticipationPhasePending {
			return model.CloudProviderParticipationStatus{}, participationStateInvalid("expire requires Pending")
		}
		next.Phase = model.ParticipationPhaseExpired
		return next, nil

	case ParticipationActionSuspend:
		if !holdsMutablePhase(current.Phase) {
			return model.CloudProviderParticipationStatus{}, participationStateInvalid("suspend denied outside Active/Suspended")
		}
		switch party {
		case HoldPartyPlatform:
			next.PlatformSuspended = true
		case HoldPartyProvider:
			next.ProviderSuspended = true
		default:
			return model.CloudProviderParticipationStatus{}, participationStateInvalid("suspend requires a hold party")
		}
		next.Phase = DerivedEffectivePhase(next.PlatformSuspended, next.ProviderSuspended)
		return next, nil

	case ParticipationActionResume:
		if !holdsMutablePhase(current.Phase) {
			return model.CloudProviderParticipationStatus{}, participationStateInvalid("resume denied outside Active/Suspended")
		}
		switch party {
		case HoldPartyPlatform:
			next.PlatformSuspended = false
		case HoldPartyProvider:
			next.ProviderSuspended = false
		default:
			return model.CloudProviderParticipationStatus{}, participationStateInvalid("resume requires a hold party")
		}
		next.Phase = DerivedEffectivePhase(next.PlatformSuspended, next.ProviderSuspended)
		return next, nil

	case ParticipationActionRequestRelease:
		if current.Phase != model.ParticipationPhaseActive && current.Phase != model.ParticipationPhaseSuspended {
			return model.CloudProviderParticipationStatus{}, participationStateInvalid("request-release requires Active or Suspended")
		}
		next.Phase = model.ParticipationPhaseTerminating
		// Holds are retained across Terminating.
		return next, nil

	case ParticipationActionAcceptRelease:
		if current.Phase != model.ParticipationPhaseTerminating {
			return model.CloudProviderParticipationStatus{}, participationStateInvalid("accept-release requires Terminating")
		}
		next.Phase = model.ParticipationPhaseTerminated
		return next, nil

	case ParticipationActionDeclineRelease:
		if current.Phase != model.ParticipationPhaseTerminating {
			return model.CloudProviderParticipationStatus{}, participationStateInvalid("decline-release requires Terminating")
		}
		next.Phase = DerivedEffectivePhase(next.PlatformSuspended, next.ProviderSuspended)
		return next, nil

	default:
		return model.CloudProviderParticipationStatus{}, participationStateInvalid("unknown participation action")
	}
}

func holdsMutablePhase(phase model.ParticipationPhase) bool {
	return phase == model.ParticipationPhaseActive || phase == model.ParticipationPhaseSuspended
}

func participationStateInvalid(detail string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeConflict).
		WithDetail(detail).
		WithViolations([]apiproblem.Violation{{
			Field:   "/status/phase",
			Code:    ViolationParticipationStateInvalid,
			Message: "participation transition is invalid",
		}})
}
