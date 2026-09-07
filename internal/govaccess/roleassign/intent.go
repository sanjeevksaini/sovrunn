// Package roleassign is the sole RoleAssignment writer boundary (VS0-WRITER-008).
//
// It alone constructs sealed PreparedRoleAssignmentIntent and
// FinalizedRoleAssignmentChange values. Approval derivation is owned elsewhere;
// this package never imports approvalreq (DD-01, DD-13 Rule 007).
package roleassign

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

// KindRoleAssignment is the RoleAssignment resource kind (VS0-SCHEMA-023).
const KindRoleAssignment = "RoleAssignment"

// OperationKind is the closed RoleAssignment mutation operation set admitted
// through the workflow-trigger and direct-revoke ports.
type OperationKind string

const (
	OpGrant            OperationKind = "grant"
	OpActivate         OperationKind = "activate"
	OpPrivilegedRevoke OperationKind = "privileged_revoke"
	OpReviewRevoke     OperationKind = "review_revoke"
	OpReviewReplace    OperationKind = "review_replace"
	OpAdvanceReviewDue OperationKind = "advance_review_due"
	OpDirectRevoke     OperationKind = "direct_revoke"
)

// Valid reports whether k is a closed OperationKind.
func (k OperationKind) Valid() bool {
	switch k {
	case OpGrant, OpActivate, OpPrivilegedRevoke, OpReviewRevoke, OpReviewReplace,
		OpAdvanceReviewDue, OpDirectRevoke:
		return true
	default:
		return false
	}
}

// DomainOutcome is a sealed non-publication RoleAssignment domain/lifecycle
// outcome returned by FinalizeAt (DD-01). It publishes nothing.
type DomainOutcome struct {
	sealed bool
	code   string
}

// Code returns the stable domain outcome code.
func (o DomainOutcome) Code() string { return o.code }

// Sealed reports whether o was package-constructed.
func (o DomainOutcome) Sealed() bool { return o.sealed }

func newDomainOutcome(code string) DomainOutcome {
	if code == "" {
		code = "domain_outcome"
	}
	return DomainOutcome{sealed: true, code: code}
}

// PreparedRoleAssignmentIntent is the publication-time-free sealed RoleAssignment
// intent (DD-01). It carries validated operation meaning and expected-version/CAS
// inputs only — never publication timestamp, final state, canonical result,
// after-state digest, MutationDescriptor, carrier ID, or completed idempotency.
type PreparedRoleAssignmentIntent struct {
	sealed          bool
	claim           state.ParticipantClaim
	operation       OperationKind
	resourceUID     string
	resourceName    string
	scopeRef        apimeta.ScopeRef
	spec            model.RoleAssignmentSpec
	hasSpec         bool
	expectedVersion string
	nextVersion     string
	nextReviewDueAt *time.Time // absolute due time supplied by review trigger; never publication time
	replacementSpec model.RoleAssignmentSpec
	hasReplacement  bool
	triggerPort     string // grant|activation|review|direct_revoke
}

// ParticipantClaim returns the sealed admitted participant claim.
func (i PreparedRoleAssignmentIntent) ParticipantClaim() state.ParticipantClaim {
	return i.claim
}

// Operation returns the sealed operation kind.
func (i PreparedRoleAssignmentIntent) Operation() OperationKind { return i.operation }

// ResourceUID returns the sealed RoleAssignment UID.
func (i PreparedRoleAssignmentIntent) ResourceUID() string { return i.resourceUID }

// TriggerPort returns the admitting workflow/direct port identity.
func (i PreparedRoleAssignmentIntent) TriggerPort() string { return i.triggerPort }

// Sealed reports whether i was constructed through a roleassign port.
func (i PreparedRoleAssignmentIntent) Sealed() bool { return i.sealed }

// HasPublicationDerivedFields reports whether the intent illegally carries
// publication-derived material. Always false for sealed intents (DD-01 Rule 009).
func (i PreparedRoleAssignmentIntent) HasPublicationDerivedFields() bool {
	return false
}

type preparedInputs struct {
	operation       OperationKind
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
	nextReviewDueAt *time.Time
	replacementSpec model.RoleAssignmentSpec
	hasReplacement  bool
	triggerPort     string
}

func newPreparedRoleAssignmentIntent(in preparedInputs) (PreparedRoleAssignmentIntent, operation.MechanicalFailure) {
	if !in.operation.Valid() {
		return PreparedRoleAssignmentIntent{}, operation.NewMechanicalFailure("operation_invalid")
	}
	if in.participantID == "" || in.resourceUID == "" || in.resourceName == "" {
		return PreparedRoleAssignmentIntent{}, operation.NewMechanicalFailure("intent_identity_incomplete")
	}
	if in.expectedVersion == "" || in.nextVersion == "" {
		return PreparedRoleAssignmentIntent{}, operation.NewMechanicalFailure("intent_versions_incomplete")
	}
	if in.triggerPort == "" {
		return PreparedRoleAssignmentIntent{}, operation.NewMechanicalFailure("intent_trigger_port_empty")
	}
	if in.resultRole == "" {
		in.resultRole = state.OriginatingResult
	}
	if !in.resultRole.Valid() {
		return PreparedRoleAssignmentIntent{}, operation.NewMechanicalFailure("intent_result_role_invalid")
	}
	if len(in.versions) == 0 {
		in.versions = []state.ExpectedVersionPredicate{{
			ResourceUID: in.resourceUID, ResourceVersion: in.expectedVersion,
		}}
	}
	versionSet, err := state.NewExpectedVersionSet(in.versions)
	if err != nil {
		return PreparedRoleAssignmentIntent{}, operation.NewMechanicalFailure("intent_version_set_invalid")
	}

	digest, fail := intentDigest(in)
	if fail.Reason() != "" {
		return PreparedRoleAssignmentIntent{}, fail
	}
	claim, err := state.NewParticipantClaim(
		state.RoleAssignmentOwner,
		in.participantID,
		digest,
		versionSet,
		in.resultRole,
	)
	if err != nil {
		return PreparedRoleAssignmentIntent{}, operation.NewMechanicalFailure("intent_claim_invalid")
	}

	out := PreparedRoleAssignmentIntent{
		sealed:          true,
		claim:           claim,
		operation:       in.operation,
		resourceUID:     in.resourceUID,
		resourceName:    in.resourceName,
		scopeRef:        in.scopeRef,
		expectedVersion: in.expectedVersion,
		nextVersion:     in.nextVersion,
		triggerPort:     in.triggerPort,
	}
	if in.hasSpec {
		out.spec = copyRoleAssignmentSpec(in.spec)
		out.hasSpec = true
	}
	if in.hasReplacement {
		out.replacementSpec = copyRoleAssignmentSpec(in.replacementSpec)
		out.hasReplacement = true
	}
	if in.nextReviewDueAt != nil {
		t := in.nextReviewDueAt.UTC()
		out.nextReviewDueAt = &t
	}
	return out, operation.MechanicalFailure{}
}

func intentDigest(in preparedInputs) ([]byte, operation.MechanicalFailure) {
	payload := struct {
		Operation       OperationKind            `json:"operation"`
		ResourceUID     string                   `json:"resourceUid"`
		ResourceName    string                   `json:"resourceName"`
		ScopeRef        apimeta.ScopeRef         `json:"scopeRef"`
		Spec            model.RoleAssignmentSpec `json:"spec,omitempty"`
		HasSpec         bool                     `json:"hasSpec"`
		ExpectedVersion string                   `json:"expectedVersion"`
		NextVersion     string                   `json:"nextVersion"`
		TriggerPort     string                   `json:"triggerPort"`
		NextReviewDueAt *time.Time               `json:"nextReviewDueAt,omitempty"`
		Replacement     model.RoleAssignmentSpec `json:"replacement,omitempty"`
		HasReplacement  bool                     `json:"hasReplacement"`
	}{
		Operation:       in.operation,
		ResourceUID:     in.resourceUID,
		ResourceName:    in.resourceName,
		ScopeRef:        in.scopeRef,
		HasSpec:         in.hasSpec,
		ExpectedVersion: in.expectedVersion,
		NextVersion:     in.nextVersion,
		TriggerPort:     in.triggerPort,
		HasReplacement:  in.hasReplacement,
	}
	if in.hasSpec {
		payload.Spec = copyRoleAssignmentSpec(in.spec)
	}
	if in.hasReplacement {
		payload.Replacement = copyRoleAssignmentSpec(in.replacementSpec)
	}
	if in.nextReviewDueAt != nil {
		t := in.nextReviewDueAt.UTC()
		payload.NextReviewDueAt = &t
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil, operation.NewMechanicalFailure("intent_digest_marshal")
	}
	sum := sha256.Sum256(raw)
	return sum[:], operation.MechanicalFailure{}
}

func copyRoleAssignmentSpec(in model.RoleAssignmentSpec) model.RoleAssignmentSpec {
	out := in
	if in.NotBefore != nil {
		t := in.NotBefore.UTC()
		out.NotBefore = &t
	}
	if in.ExpiresAt != nil {
		t := in.ExpiresAt.UTC()
		out.ExpiresAt = &t
	}
	if in.ResourceRef != nil {
		r := *in.ResourceRef
		out.ResourceRef = &r
	}
	if in.MembershipRef != nil {
		r := *in.MembershipRef
		out.MembershipRef = &r
	}
	if in.ResponsiblePartyRef != nil {
		r := *in.ResponsiblePartyRef
		out.ResponsiblePartyRef = &r
	}
	if in.AccessReviewRuleRef != nil {
		r := *in.AccessReviewRuleRef
		out.AccessReviewRuleRef = &r
	}
	return out
}
