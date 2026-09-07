package roleassign

import (
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// RemediationKind is Revoke | Replace | AdvanceReviewDue for the RD-16 trigger.
type RemediationKind string

const (
	RemediationRevoke           RemediationKind = "revoke"
	RemediationReplace          RemediationKind = "replace"
	RemediationAdvanceReviewDue RemediationKind = "advance_review_due"
)

// Valid reports whether k is a closed RemediationKind.
func (k RemediationKind) Valid() bool {
	switch k {
	case RemediationRevoke, RemediationReplace, RemediationAdvanceReviewDue:
		return true
	default:
		return false
	}
}

// ReviewRemediationTriggerIntent is the immutable RD-16
// revocation/replacement/due-advancement trigger submitted by review to
// ReviewRemediationPort.
type ReviewRemediationTriggerIntent struct {
	sealed          bool
	kind            RemediationKind
	participantID   state.ParticipantID
	resourceUID     string
	resourceName    string
	scopeRef        apimeta.ScopeRef
	spec            model.RoleAssignmentSpec
	hasSpec         bool
	replacementSpec model.RoleAssignmentSpec
	hasReplacement  bool
	expectedVersion string
	nextVersion     string
	versions        []state.ExpectedVersionPredicate
	resultRole      state.ResultRole
	nextReviewDueAt *time.Time
}

// NewReviewRemediationTriggerIntent seals an RD-16 review remediation trigger.
func NewReviewRemediationTriggerIntent(
	kind RemediationKind,
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec model.RoleAssignmentSpec,
	hasSpec bool,
	replacementSpec model.RoleAssignmentSpec,
	hasReplacement bool,
	expectedVersion string,
	nextVersion string,
	versions []state.ExpectedVersionPredicate,
	resultRole state.ResultRole,
	nextReviewDueAt *time.Time,
) (ReviewRemediationTriggerIntent, operation.MechanicalFailure) {
	if !kind.Valid() {
		return ReviewRemediationTriggerIntent{}, operation.NewMechanicalFailure("remediation_kind_invalid")
	}
	if participantID == "" || resourceUID == "" || resourceName == "" {
		return ReviewRemediationTriggerIntent{}, operation.NewMechanicalFailure("review_trigger_identity_incomplete")
	}
	if expectedVersion == "" || nextVersion == "" {
		return ReviewRemediationTriggerIntent{}, operation.NewMechanicalFailure("review_trigger_versions_incomplete")
	}
	if kind == RemediationReplace && !hasReplacement {
		return ReviewRemediationTriggerIntent{}, operation.NewMechanicalFailure("replace_spec_required")
	}
	if kind == RemediationAdvanceReviewDue && nextReviewDueAt == nil {
		return ReviewRemediationTriggerIntent{}, operation.NewMechanicalFailure("review_due_required")
	}
	if resultRole == "" {
		resultRole = state.OriginatingResult
	}
	out := ReviewRemediationTriggerIntent{
		sealed:          true,
		kind:            kind,
		participantID:   participantID,
		resourceUID:     resourceUID,
		resourceName:    resourceName,
		scopeRef:        scopeRef,
		hasSpec:         hasSpec,
		hasReplacement:  hasReplacement,
		expectedVersion: expectedVersion,
		nextVersion:     nextVersion,
		versions:        append([]state.ExpectedVersionPredicate(nil), versions...),
		resultRole:      resultRole,
	}
	if hasSpec {
		out.spec = copyRoleAssignmentSpec(spec)
	}
	if hasReplacement {
		out.replacementSpec = copyRoleAssignmentSpec(replacementSpec)
	}
	if nextReviewDueAt != nil {
		t := nextReviewDueAt.UTC()
		out.nextReviewDueAt = &t
	}
	return out, operation.MechanicalFailure{}
}

// Sealed reports whether the trigger was package-constructed.
func (t ReviewRemediationTriggerIntent) Sealed() bool { return t.sealed }

// Kind returns the sealed remediation kind.
func (t ReviewRemediationTriggerIntent) Kind() RemediationKind { return t.kind }

// ReviewRemediationPort is the RD-16 workflow-trigger port (design §3.3). Only
// review may invoke Submit (F18-ARCH-006).
type ReviewRemediationPort struct{}

// NewReviewRemediationPort returns the RoleAssignment RD-16 remediation trigger port.
func NewReviewRemediationPort() ReviewRemediationPort { return ReviewRemediationPort{} }

// Submit performs RoleAssignment-owned semantic preparation for an RD-16
// revocation, replacement, or review-due advancement trigger.
func (ReviewRemediationPort) Submit(trigger ReviewRemediationTriggerIntent) (PreparedRoleAssignmentIntent, operation.MechanicalFailure) {
	if !trigger.sealed {
		return PreparedRoleAssignmentIntent{}, operation.NewMechanicalFailure("review_trigger_unsealed")
	}
	var op OperationKind
	switch trigger.kind {
	case RemediationRevoke:
		op = OpReviewRevoke
	case RemediationReplace:
		op = OpReviewReplace
	case RemediationAdvanceReviewDue:
		op = OpAdvanceReviewDue
	default:
		return PreparedRoleAssignmentIntent{}, operation.NewMechanicalFailure("remediation_kind_invalid")
	}
	return newPreparedRoleAssignmentIntent(preparedInputs{
		operation:       op,
		participantID:   trigger.participantID,
		resourceUID:     trigger.resourceUID,
		resourceName:    trigger.resourceName,
		scopeRef:        trigger.scopeRef,
		spec:            trigger.spec,
		hasSpec:         trigger.hasSpec,
		replacementSpec: trigger.replacementSpec,
		hasReplacement:  trigger.hasReplacement,
		expectedVersion: trigger.expectedVersion,
		nextVersion:     trigger.nextVersion,
		versions:        trigger.versions,
		resultRole:      trigger.resultRole,
		nextReviewDueAt: trigger.nextReviewDueAt,
		triggerPort:     "review",
	})
}
