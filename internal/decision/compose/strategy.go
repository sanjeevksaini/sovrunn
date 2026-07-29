// Package compose provides FEATURE-0013 pure composition strategies and the
// deterministic replay kernel (design §6.3, §8; F13-COMP-001/002; RID-06,
// RID-10; AD-005, AD-043; closures 2).
//
// Strategies are pure functions: they perform no I/O, hold no session state,
// and never become an independent fail-open authority. Strategy-level
// fail-open is effective only through a structurally valid governing-profile
// SecurityExceptionRef; otherwise the strategy fails closed. No approval
// workflow, issuance, or revocation is created here.
package compose

import (
	"encoding/json"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/apiproblem"
	"github.com/sanjeevksaini/sovrunn/internal/decision"
)

// Closed composition violation code used when a strategy fails closed on
// evaluator failure without an effective fail-open gate (architecture §15.2).
// Canonical constants live in internal/decision/validate (T-017); this
// reproduces the closed string so compose can emit Problem Details now.
const violationInputConflict apiproblem.ViolationCode = "DECISION_COMPOSITION_INPUT_CONFLICT"

// Strategy is a pure composition function (RID-06; design §8).
// It must not perform I/O, persistence, networking, or evaluator rerun.
type Strategy func(in StrategyInput) (StrategyOutput, *apiproblem.Problem)

// StrategyMetadata is the registered declaration for a versioned composition
// strategy (architecture §7.2; F13-COMP-001). It describes handling and
// budgets only. It never carries SecurityExceptionRef, approval, issuance,
// or revocation fields and is never an independent fail-open authority
// (RID-10; closure 2).
type StrategyMetadata struct {
	ID      string `json:"id"`
	Version string `json:"version"`

	AcceptedInputTypes []string `json:"acceptedInputTypes"`
	Ordering           string   `json:"ordering"`
	Precedence         string   `json:"precedence,omitempty"`
	ShortCircuit       bool     `json:"shortCircuit"`

	MissingBehavior  string `json:"missingBehavior"`
	TimeoutBehavior  string `json:"timeoutBehavior"`
	ConflictBehavior string `json:"conflictBehavior"`
	ErrorBehavior    string `json:"errorBehavior"`

	// RequestsFailOpen is a descriptive request only. It never authorizes
	// fail-open by itself. Effectiveness is gated by FailOpenEffective
	// against a valid governing-profile SecurityExceptionRef (RID-10).
	RequestsFailOpen bool `json:"requestsFailOpen,omitempty"`

	DeterministicOutput bool  `json:"deterministicOutput"`
	MaxFanOut           int   `json:"maxFanOut,omitempty"`
	MaxPayloadBytes     int64 `json:"maxPayloadBytes,omitempty"`
	MaxLatencyMs        int64 `json:"maxLatencyMs,omitempty"`

	ExplanationRules string `json:"explanationRules,omitempty"`
	EvidenceRules    string `json:"evidenceRules,omitempty"`
}

// StrategyInput is the pure in-memory input to a Strategy (design §8).
// Evaluations are captured results; strategies never invoke evaluators.
type StrategyInput struct {
	Meta            StrategyMetadata
	FailureBehavior decision.FailureBehavior
	Evaluations     []decision.EvaluationResult
	InputRefs       []string // ordered input identities for deterministic recording
}

// StrategyOutput is the deterministic composition result (architecture §7.2).
// Chosen strategy identity, ordered input identities, and output are recorded
// so the result can be reproduced without digest computation in FEATURE-0013.
type StrategyOutput struct {
	StrategyID      string   `json:"strategyId"`
	StrategyVersion string   `json:"strategyVersion"`
	InputRefs       []string `json:"inputRefs,omitempty"`

	TypedResult  json.RawMessage               `json:"typedResult,omitempty"`
	Rationale    decision.DecisionRationale    `json:"rationale"`
	Obligations  []decision.Obligation         `json:"obligations,omitempty"`
	Adjudication *decision.AdjudicationOutcome `json:"adjudication,omitempty"`

	// FailOpenApplied is true only when FailOpenEffective was true and the
	// strategy actually applied fail-open handling for a failure mode.
	FailOpenApplied bool `json:"failOpenApplied,omitempty"`
}

// FailOpenEffective reports whether strategy-requested fail-open is effective
// under RID-10 / closure 2 / architecture §27.9(2).
//
// Rules (fail-closed default):
//   - Strategy metadata alone never authorizes fail-open.
//   - RequestsFailOpen must be true on the strategy metadata.
//   - Governing profile FailPosture must be FAIL_OPEN.
//   - Governing profile must carry a structurally valid SecurityExceptionRef.
//
// Otherwise the result is false (fail closed). No approval workflow is
// created or consulted; SecurityExceptionRef is structural evidence only.
func FailOpenEffective(meta StrategyMetadata, fb decision.FailureBehavior) bool {
	if !meta.RequestsFailOpen {
		return false
	}
	if fb.FailPosture != decision.FailPostureFailOpen {
		return false
	}
	return exceptionEvidenceValid(fb.SecurityExceptionRef)
}

// ApplyFailure handles a composition failure mode under the RID-10 gate.
// When FailOpenEffective is false, it fails closed with a Problem.
// When effective, it returns a StrategyOutput that records fail-open
// application without inventing an alternate approval authority.
func ApplyFailure(in StrategyInput, mode string) (StrategyOutput, *apiproblem.Problem) {
	if !FailOpenEffective(in.Meta, in.FailureBehavior) {
		return StrategyOutput{}, composeProblem(
			violationInputConflict,
			"/composition/strategy",
			"composition strategy failed closed on "+mode+" without effective fail-open evidence",
		)
	}
	return StrategyOutput{
		StrategyID:      in.Meta.ID,
		StrategyVersion: in.Meta.Version,
		InputRefs:       cloneStrings(in.InputRefs),
		FailOpenApplied: true,
		Rationale: decision.DecisionRationale{
			ReasonCodes: []string{"COMPOSITION_FAIL_OPEN"},
			Reasons:     []string{"fail-open effective via governing-profile SecurityExceptionRef for mode " + mode},
			Warnings:    []string{"fail-open is profile-exception evidence only; strategy metadata is not an approval authority"},
		},
	}, nil
}

// exceptionEvidenceValid reports structural requiredness of SecurityExceptionRef
// for the strategy fail-open gate (design §7.7). Full profile validation
// (owner/scope/expiry-vs-decision-time/covered-mode registry) belongs to T-019;
// this gate only accepts structurally complete evidence and otherwise fails closed.
func exceptionEvidenceValid(ref *decision.SecurityExceptionRef) bool {
	if ref == nil {
		return false
	}
	if !typedRefUIDPinned(ref.ApprovalRef) {
		return false
	}
	if !typedRefUIDPinned(ref.OwnerRef) {
		return false
	}
	if !typedRefUIDPinned(ref.ApprovingAuthorityRef) {
		return false
	}
	if !nonEmptyStrings(ref.CompensatingControls) {
		return false
	}
	if !nonEmptyStrings(ref.CoveredFailureModes) {
		return false
	}
	if ref.EffectiveFrom == "" || ref.EffectiveUntil == "" {
		return false
	}
	// Non-positive interval fails closed (design §7.7 structural check).
	if ref.EffectiveUntil <= ref.EffectiveFrom {
		return false
	}
	if ref.AuditTreatment == "" || ref.ReassessmentTrigger == "" || ref.Purpose == "" {
		return false
	}
	return true
}

func typedRefUIDPinned(ref apimeta.TypedRef) bool {
	return ref.APIVersion != "" && ref.Kind != "" && ref.Name != "" && ref.UID != ""
}

func nonEmptyStrings(vals []string) bool {
	if len(vals) == 0 {
		return false
	}
	for _, v := range vals {
		if v == "" {
			return false
		}
	}
	return true
}

func cloneStrings(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func composeProblem(code apiproblem.ViolationCode, field, message string) *apiproblem.Problem {
	return apiproblem.New(apiproblem.CodeValidationFailed).
		WithDetail(message).
		WithViolations([]apiproblem.Violation{{
			Field:   field,
			Code:    code,
			Message: message,
		}})
}
