package idempotency

import "github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"

// ReplayDecisionKind is the closed replay decision set after key lookup.
type ReplayDecisionKind string

const (
	ReplayDecisionProceed  ReplayDecisionKind = "Proceed"
	ReplayDecisionReplay   ReplayDecisionKind = "Replay"
	ReplayDecisionWait     ReplayDecisionKind = "Wait"
	ReplayDecisionConflict ReplayDecisionKind = "Conflict"
)

// ReplayDecision captures replay-recheck outcome over lookup + RequestBinding.
type ReplayDecision struct {
	kind       ReplayDecisionKind
	replay     CompletedRecord
	hasReplay  bool
	wait       <-chan struct{}
	hasWait    bool
	reserveHit bool
}

// Kind returns the closed replay decision.
func (d ReplayDecision) Kind() ReplayDecisionKind { return d.kind }

// ReplayRecord returns the completed record when Kind is Replay.
func (d ReplayDecision) ReplayRecord() (CompletedRecord, bool) {
	if !d.hasReplay {
		return CompletedRecord{}, false
	}
	return d.replay, true
}

// WaitChannel returns the waiter channel when Kind is Wait.
func (d ReplayDecision) WaitChannel() (<-chan struct{}, bool) {
	if !d.hasWait {
		return nil, false
	}
	return d.wait, true
}

// EvaluateReplayRecheck compares RequestBinding after lookup-key match.
func EvaluateReplayRecheck(lookup LookupResult, binding RequestBinding) (ReplayDecision, operation.MechanicalFailure) {
	if !binding.Sealed() {
		return ReplayDecision{}, operation.NewMechanicalFailure("replay_binding_unsealed")
	}
	if rec, _, ok := lookup.Completed(); ok {
		if !rec.Binding().Equal(binding) {
			return ReplayDecision{kind: ReplayDecisionConflict}, operation.MechanicalFailure{}
		}
		return ReplayDecision{kind: ReplayDecisionReplay, replay: rec, hasReplay: true}, operation.MechanicalFailure{}
	}
	if inFlight, wait, ok := lookup.InFlight(); ok {
		if !inFlight.Binding().Equal(binding) {
			return ReplayDecision{kind: ReplayDecisionConflict}, operation.MechanicalFailure{}
		}
		return ReplayDecision{kind: ReplayDecisionWait, wait: wait, hasWait: true}, operation.MechanicalFailure{}
	}
	return ReplayDecision{kind: ReplayDecisionProceed}, operation.MechanicalFailure{}
}
