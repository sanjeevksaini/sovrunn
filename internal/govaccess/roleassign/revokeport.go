package roleassign

import (
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// DirectRevokeIntent is the immutable RD-08 caller revoke intent prepared through
// DirectRevokePort. It is not a fourth workflow trigger (design §3.3).
type DirectRevokeIntent struct {
	sealed          bool
	participantID   state.ParticipantID
	resourceUID     string
	resourceName    string
	scopeRef        apimeta.ScopeRef
	spec            model.RoleAssignmentSpec
	hasSpec         bool
	expectedVersion string
	nextVersion     string
	versions        []state.ExpectedVersionPredicate
	resultRole      state.ResultRole
}

// NewDirectRevokeIntent seals an RD-08 direct revoke preparation request.
func NewDirectRevokeIntent(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec model.RoleAssignmentSpec,
	hasSpec bool,
	expectedVersion string,
	nextVersion string,
	versions []state.ExpectedVersionPredicate,
	resultRole state.ResultRole,
) (DirectRevokeIntent, operation.MechanicalFailure) {
	if participantID == "" || resourceUID == "" || resourceName == "" {
		return DirectRevokeIntent{}, operation.NewMechanicalFailure("direct_revoke_identity_incomplete")
	}
	if expectedVersion == "" || nextVersion == "" {
		return DirectRevokeIntent{}, operation.NewMechanicalFailure("direct_revoke_versions_incomplete")
	}
	if resultRole == "" {
		resultRole = state.OriginatingResult
	}
	out := DirectRevokeIntent{
		sealed:          true,
		participantID:   participantID,
		resourceUID:     resourceUID,
		resourceName:    resourceName,
		scopeRef:        scopeRef,
		hasSpec:         hasSpec,
		expectedVersion: expectedVersion,
		nextVersion:     nextVersion,
		versions:        append([]state.ExpectedVersionPredicate(nil), versions...),
		resultRole:      resultRole,
	}
	if hasSpec {
		out.spec = copyRoleAssignmentSpec(spec)
	}
	return out, operation.MechanicalFailure{}
}

// Sealed reports whether the intent was package-constructed.
func (i DirectRevokeIntent) Sealed() bool { return i.sealed }

// DirectRevokePort performs RoleAssignment-owned semantic preparation for the
// RD-08 caller revoke operation. Only command.RoleAssignmentRevokeService may
// invoke Prepare (F18-ARCH-006). It is not a workflow trigger, route, or writer.
type DirectRevokePort struct{}

// NewDirectRevokePort returns the RoleAssignment RD-08 direct revoke port.
func NewDirectRevokePort() DirectRevokePort { return DirectRevokePort{} }

// Prepare performs only RoleAssignment-owned semantic preparation and returns
// the same sealed PreparedRoleAssignmentIntent type used by workflow triggers.
func (DirectRevokePort) Prepare(in DirectRevokeIntent) (PreparedRoleAssignmentIntent, operation.MechanicalFailure) {
	if !in.sealed {
		return PreparedRoleAssignmentIntent{}, operation.NewMechanicalFailure("direct_revoke_unsealed")
	}
	return newPreparedRoleAssignmentIntent(preparedInputs{
		operation:       OpDirectRevoke,
		participantID:   in.participantID,
		resourceUID:     in.resourceUID,
		resourceName:    in.resourceName,
		scopeRef:        in.scopeRef,
		spec:            in.spec,
		hasSpec:         in.hasSpec,
		expectedVersion: in.expectedVersion,
		nextVersion:     in.nextVersion,
		versions:        in.versions,
		resultRole:      in.resultRole,
		triggerPort:     "direct_revoke",
	})
}
