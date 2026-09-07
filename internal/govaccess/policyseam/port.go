package policyseam

import (
	"context"

	"github.com/sanjeevksaini/sovrunn/internal/decision"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/policyeval"
)

// Port is the sole FEATURE-0018 package importing policyeval.
type Port struct {
	boundary policyBoundary
}

type policyBoundary interface {
	Evaluate(context.Context, policyeval.PolicyEvaluationRequest) (policyeval.PolicyEvaluationResult, decision.TimingEnvelope, policyeval.NonResult, error)
}

func NewPort(boundary policyBoundary) (Port, operation.MechanicalFailure) {
	if boundary == nil {
		return Port{}, operation.NewMechanicalFailure("policy_boundary_nil")
	}
	return Port{boundary: boundary}, operation.MechanicalFailure{}
}

type Evaluation struct {
	Outcome          string
	ReasonCodes      []string
	RequiresApproval bool
}

func (p Port) EvaluatePrivileged(
	ctx context.Context,
	req policyeval.PolicyEvaluationRequest,
) (Evaluation, operation.MechanicalFailure) {
	if p.boundary == nil {
		return Evaluation{}, operation.NewMechanicalFailure("policy_boundary_unset")
	}
	result, _, nonResult, err := p.boundary.Evaluate(ctx, req)
	if err != nil || nonResult != "" {
		return Evaluation{
			Outcome:     "Deny",
			ReasonCodes: []string{"POLICY_EVALUATION_UNAVAILABLE"},
		}, operation.MechanicalFailure{}
	}
	switch result.Outcome {
	case policyeval.OutcomeAllow:
		return Evaluation{Outcome: "Allow", ReasonCodes: append([]string(nil), result.ReasonCodes...)}, operation.MechanicalFailure{}
	case policyeval.OutcomeDeny:
		return Evaluation{Outcome: "Deny", ReasonCodes: append([]string(nil), result.ReasonCodes...)}, operation.MechanicalFailure{}
	case policyeval.OutcomeRequiresApproval:
		return Evaluation{
			Outcome:          "RequiresApproval",
			ReasonCodes:      append([]string(nil), result.ReasonCodes...),
			RequiresApproval: true,
		}, operation.MechanicalFailure{}
	case policyeval.OutcomeIndeterminate:
		return Evaluation{
			Outcome:     "Deny",
			ReasonCodes: []string{"POLICY_INDETERMINATE_MAPPED_TO_DENY"},
		}, operation.MechanicalFailure{}
	default:
		return Evaluation{
			Outcome:     "Deny",
			ReasonCodes: []string{"POLICY_OUTCOME_INVALID"},
		}, operation.MechanicalFailure{}
	}
}
