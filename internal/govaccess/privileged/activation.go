package privileged

import (
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/roleassign"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// ActivationService submits only RD-14 activation/privileged-revoke trigger
// intents to roleassign.ActivationPort. It never constructs RoleAssignment
// prepared/finalized types.
type ActivationService struct {
	port roleassign.ActivationPort
}

func NewActivationService(port roleassign.ActivationPort) ActivationService {
	return ActivationService{port: port}
}

func (s ActivationService) SubmitActivate(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec model.RoleAssignmentSpec,
	expectedVersion string,
	nextVersion string,
) (roleassign.PreparedRoleAssignmentIntent, operation.MechanicalFailure) {
	var zero roleassign.PreparedRoleAssignmentIntent
	trigger, fail := roleassign.NewActivationTriggerIntent(
		roleassign.ActivationActivate,
		participantID,
		resourceUID,
		resourceName,
		scopeRef,
		spec,
		true,
		expectedVersion,
		nextVersion,
		nil,
		state.OriginatingResult,
	)
	if fail.Reason() != "" {
		return zero, fail
	}
	return s.port.Submit(trigger)
}

func (s ActivationService) SubmitPrivilegedRevoke(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	expectedVersion string,
	nextVersion string,
) (roleassign.PreparedRoleAssignmentIntent, operation.MechanicalFailure) {
	var zero roleassign.PreparedRoleAssignmentIntent
	trigger, fail := roleassign.NewActivationTriggerIntent(
		roleassign.ActivationPrivilegedRevoke,
		participantID,
		resourceUID,
		resourceName,
		scopeRef,
		model.RoleAssignmentSpec{},
		false,
		expectedVersion,
		nextVersion,
		nil,
		state.OriginatingResult,
	)
	if fail.Reason() != "" {
		return zero, fail
	}
	return s.port.Submit(trigger)
}
