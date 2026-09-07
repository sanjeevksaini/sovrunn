package command

import (
	"context"

	"github.com/sanjeevksaini/sovrunn/internal/govaccess/authzeval"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/roleassign"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/uow"
)

// CandidateBinder binds an evaluated candidate to an OperationCandidate.
// The binder is supplied by the caller so command owns no state/currentness APIs.
type CandidateBinder func(authzeval.AuthorizationCandidate) (authzeval.OperationCandidate, operation.MechanicalFailure)

// RoleAssignmentRevokeCommand is the RD-08 caller-operation input for command orchestration.
type RoleAssignmentRevokeCommand struct {
	AuthorizationInput  model.AuthorizationInput
	BindCandidate       CandidateBinder
	DirectRevokeIntent  roleassign.DirectRevokeIntent
	BuildCallerMutation func(
		candidate authzeval.OperationCandidate,
		intent roleassign.PreparedRoleAssignmentIntent,
	) (uow.CallerMutationRequest, operation.MechanicalFailure)
}

type candidateEvaluator interface {
	Evaluate(input model.AuthorizationInput) (authzeval.AuthorizationCandidate, operation.MechanicalFailure)
}

type directRevokePreparer interface {
	Prepare(in roleassign.DirectRevokeIntent) (roleassign.PreparedRoleAssignmentIntent, operation.MechanicalFailure)
}

// RoleAssignmentRevokeService orchestrates the RD-08 caller operation:
// evaluate -> roleassign.DirectRevokePort -> uow.MutationCoordinatorPort.
type RoleAssignmentRevokeService struct {
	evaluator   candidateEvaluator
	preparer    directRevokePreparer
	coordinator uow.MutationCoordinatorPort
}

// NewRoleAssignmentRevokeService wires the command-only RD-08 orchestration service.
func NewRoleAssignmentRevokeService(
	evaluator candidateEvaluator,
	preparer directRevokePreparer,
	coordinator uow.MutationCoordinatorPort,
) RoleAssignmentRevokeService {
	return RoleAssignmentRevokeService{
		evaluator:   evaluator,
		preparer:    preparer,
		coordinator: coordinator,
	}
}

// Execute performs in-process caller orchestration only.
func (s RoleAssignmentRevokeService) Execute(
	ctx context.Context,
	cmd RoleAssignmentRevokeCommand,
) (uow.MutationOutcome, operation.MechanicalFailure) {
	var zero uow.MutationOutcome
	if s.evaluator == nil {
		return zero, operation.NewMechanicalFailure("command_evaluator_nil")
	}
	if s.preparer == nil {
		return zero, operation.NewMechanicalFailure("command_preparer_nil")
	}
	if s.coordinator == nil {
		return zero, operation.NewMechanicalFailure("command_coordinator_nil")
	}
	if cmd.BindCandidate == nil {
		return zero, operation.NewMechanicalFailure("command_bind_candidate_nil")
	}
	if cmd.BuildCallerMutation == nil {
		return zero, operation.NewMechanicalFailure("command_build_request_nil")
	}

	candidate, fail := s.evaluator.Evaluate(cmd.AuthorizationInput)
	if fail.Reason() != "" {
		return zero, fail
	}
	operationCandidate, fail := cmd.BindCandidate(candidate)
	if fail.Reason() != "" {
		return zero, fail
	}
	intent, fail := s.preparer.Prepare(cmd.DirectRevokeIntent)
	if fail.Reason() != "" {
		return zero, fail
	}

	request, fail := cmd.BuildCallerMutation(operationCandidate, intent)
	if fail.Reason() != "" {
		return zero, fail
	}

	return s.coordinator.CoordinateCallerMutation(ctx, request), operation.MechanicalFailure{}
}
