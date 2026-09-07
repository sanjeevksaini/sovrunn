package roleassign

import (
	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// ActivationKind is Activate | PrivilegedRevoke for the RD-14 trigger.
type ActivationKind string

const (
	ActivationActivate         ActivationKind = "activate"
	ActivationPrivilegedRevoke ActivationKind = "privileged_revoke"
)

// Valid reports whether k is a closed ActivationKind.
func (k ActivationKind) Valid() bool {
	switch k {
	case ActivationActivate, ActivationPrivilegedRevoke:
		return true
	default:
		return false
	}
}

// ActivationTriggerIntent is the immutable RD-14 activation/privileged-revocation
// trigger submitted by privileged to ActivationPort.
type ActivationTriggerIntent struct {
	sealed          bool
	kind            ActivationKind
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

// NewActivationTriggerIntent seals an RD-14 activation or privileged-revocation trigger.
func NewActivationTriggerIntent(
	kind ActivationKind,
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
) (ActivationTriggerIntent, operation.MechanicalFailure) {
	if !kind.Valid() {
		return ActivationTriggerIntent{}, operation.NewMechanicalFailure("activation_kind_invalid")
	}
	if participantID == "" || resourceUID == "" || resourceName == "" {
		return ActivationTriggerIntent{}, operation.NewMechanicalFailure("activation_trigger_identity_incomplete")
	}
	if expectedVersion == "" || nextVersion == "" {
		return ActivationTriggerIntent{}, operation.NewMechanicalFailure("activation_trigger_versions_incomplete")
	}
	if kind == ActivationActivate && !hasSpec {
		return ActivationTriggerIntent{}, operation.NewMechanicalFailure("activation_spec_required")
	}
	if resultRole == "" {
		resultRole = state.OriginatingResult
	}
	out := ActivationTriggerIntent{
		sealed:          true,
		kind:            kind,
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

// Sealed reports whether the trigger was package-constructed.
func (t ActivationTriggerIntent) Sealed() bool { return t.sealed }

// Kind returns the sealed activation kind.
func (t ActivationTriggerIntent) Kind() ActivationKind { return t.kind }

// ActivationPort is the RD-14 workflow-trigger port (design §3.3). Only
// privileged may invoke Submit (F18-ARCH-006).
type ActivationPort struct{}

// NewActivationPort returns the RoleAssignment RD-14 activation trigger port.
func NewActivationPort() ActivationPort { return ActivationPort{} }

// Submit performs RoleAssignment-owned semantic preparation for an RD-14
// activation or privileged-revocation trigger.
func (ActivationPort) Submit(trigger ActivationTriggerIntent) (PreparedRoleAssignmentIntent, operation.MechanicalFailure) {
	if !trigger.sealed {
		return PreparedRoleAssignmentIntent{}, operation.NewMechanicalFailure("activation_trigger_unsealed")
	}
	op := OpActivate
	if trigger.kind == ActivationPrivilegedRevoke {
		op = OpPrivilegedRevoke
	}
	return newPreparedRoleAssignmentIntent(preparedInputs{
		operation:       op,
		participantID:   trigger.participantID,
		resourceUID:     trigger.resourceUID,
		resourceName:    trigger.resourceName,
		scopeRef:        trigger.scopeRef,
		spec:            trigger.spec,
		hasSpec:         trigger.hasSpec,
		expectedVersion: trigger.expectedVersion,
		nextVersion:     trigger.nextVersion,
		versions:        trigger.versions,
		resultRole:      trigger.resultRole,
		triggerPort:     "activation",
	})
}
