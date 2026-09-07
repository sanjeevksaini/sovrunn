package review

import (
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/roleassign"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// RemediationService submits only RD-16 remediation trigger intents to
// roleassign.ReviewRemediationPort.
type RemediationService struct {
	port roleassign.ReviewRemediationPort
}

func NewRemediationService(port roleassign.ReviewRemediationPort) RemediationService {
	return RemediationService{port: port}
}

func (s RemediationService) SubmitRevoke(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec model.RoleAssignmentSpec,
	expectedVersion string,
	nextVersion string,
) (roleassign.PreparedRoleAssignmentIntent, operation.MechanicalFailure) {
	var zero roleassign.PreparedRoleAssignmentIntent
	trigger, fail := roleassign.NewReviewRemediationTriggerIntent(
		roleassign.RemediationRevoke,
		participantID,
		resourceUID,
		resourceName,
		scopeRef,
		spec,
		true,
		model.RoleAssignmentSpec{},
		false,
		expectedVersion,
		nextVersion,
		nil,
		state.OriginatingResult,
		nil,
	)
	if fail.Reason() != "" {
		return zero, fail
	}
	return s.port.Submit(trigger)
}

func (s RemediationService) SubmitReplace(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	current model.RoleAssignmentSpec,
	replacement model.RoleAssignmentSpec,
	expectedVersion string,
	nextVersion string,
) (roleassign.PreparedRoleAssignmentIntent, operation.MechanicalFailure) {
	var zero roleassign.PreparedRoleAssignmentIntent
	trigger, fail := roleassign.NewReviewRemediationTriggerIntent(
		roleassign.RemediationReplace,
		participantID,
		resourceUID,
		resourceName,
		scopeRef,
		current,
		true,
		replacement,
		true,
		expectedVersion,
		nextVersion,
		nil,
		state.OriginatingResult,
		nil,
	)
	if fail.Reason() != "" {
		return zero, fail
	}
	return s.port.Submit(trigger)
}

func (s RemediationService) SubmitAdvanceReviewDue(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec model.RoleAssignmentSpec,
	expectedVersion string,
	nextVersion string,
	nextReviewDueAt time.Time,
) (roleassign.PreparedRoleAssignmentIntent, operation.MechanicalFailure) {
	var zero roleassign.PreparedRoleAssignmentIntent
	due := nextReviewDueAt.UTC()
	trigger, fail := roleassign.NewReviewRemediationTriggerIntent(
		roleassign.RemediationAdvanceReviewDue,
		participantID,
		resourceUID,
		resourceName,
		scopeRef,
		spec,
		true,
		model.RoleAssignmentSpec{},
		false,
		expectedVersion,
		nextVersion,
		nil,
		state.OriginatingResult,
		&due,
	)
	if fail.Reason() != "" {
		return zero, fail
	}
	return s.port.Submit(trigger)
}
