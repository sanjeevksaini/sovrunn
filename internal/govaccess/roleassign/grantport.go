package roleassign

import (
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// GrantTriggerIntent is the immutable RD-13 grant/materialization trigger
// submitted by grantintent to GrantPort. Callers never construct a
// PreparedRoleAssignmentIntent directly.
type GrantTriggerIntent struct {
	sealed          bool
	participantID   state.ParticipantID
	resourceUID     string
	resourceName    string
	scopeRef        apimeta.ScopeRef
	spec            model.RoleAssignmentSpec
	expectedVersion string
	nextVersion     string
	versions        []state.ExpectedVersionPredicate
	resultRole      state.ResultRole
	nextReviewDueAt *time.Time
}

// NewGrantTriggerIntent seals an RD-13 grant/materialization trigger.
func NewGrantTriggerIntent(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec model.RoleAssignmentSpec,
	expectedVersion string,
	nextVersion string,
	versions []state.ExpectedVersionPredicate,
	resultRole state.ResultRole,
	nextReviewDueAt *time.Time,
) (GrantTriggerIntent, operation.MechanicalFailure) {
	if participantID == "" || resourceUID == "" || resourceName == "" {
		return GrantTriggerIntent{}, operation.NewMechanicalFailure("grant_trigger_identity_incomplete")
	}
	if expectedVersion == "" || nextVersion == "" {
		return GrantTriggerIntent{}, operation.NewMechanicalFailure("grant_trigger_versions_incomplete")
	}
	if !spec.Validity.Valid() {
		return GrantTriggerIntent{}, operation.NewMechanicalFailure("grant_trigger_validity_invalid")
	}
	if resultRole == "" {
		resultRole = state.OriginatingResult
	}
	out := GrantTriggerIntent{
		sealed:          true,
		participantID:   participantID,
		resourceUID:     resourceUID,
		resourceName:    resourceName,
		scopeRef:        scopeRef,
		spec:            copyRoleAssignmentSpec(spec),
		expectedVersion: expectedVersion,
		nextVersion:     nextVersion,
		versions:        append([]state.ExpectedVersionPredicate(nil), versions...),
		resultRole:      resultRole,
	}
	if nextReviewDueAt != nil {
		t := nextReviewDueAt.UTC()
		out.nextReviewDueAt = &t
	}
	return out, operation.MechanicalFailure{}
}

// Sealed reports whether the trigger was constructed through NewGrantTriggerIntent.
func (t GrantTriggerIntent) Sealed() bool { return t.sealed }

// GrantPort is the RD-13 workflow-trigger port (design §3.3). Only grantintent
// may invoke Submit (F18-ARCH-006).
type GrantPort struct{}

// NewGrantPort returns the RoleAssignment RD-13 grant trigger port.
func NewGrantPort() GrantPort { return GrantPort{} }

// Submit performs RoleAssignment-owned semantic preparation for an RD-13
// grant/materialization trigger and returns a sealed PreparedRoleAssignmentIntent.
func (GrantPort) Submit(trigger GrantTriggerIntent) (PreparedRoleAssignmentIntent, operation.MechanicalFailure) {
	if !trigger.sealed {
		return PreparedRoleAssignmentIntent{}, operation.NewMechanicalFailure("grant_trigger_unsealed")
	}
	return newPreparedRoleAssignmentIntent(preparedInputs{
		operation:       OpGrant,
		participantID:   trigger.participantID,
		resourceUID:     trigger.resourceUID,
		resourceName:    trigger.resourceName,
		scopeRef:        trigger.scopeRef,
		spec:            trigger.spec,
		hasSpec:         true,
		expectedVersion: trigger.expectedVersion,
		nextVersion:     trigger.nextVersion,
		versions:        trigger.versions,
		resultRole:      trigger.resultRole,
		nextReviewDueAt: trigger.nextReviewDueAt,
		triggerPort:     "grant",
	})
}
