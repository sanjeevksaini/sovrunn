package membership

import (
	"crypto/sha256"
	"encoding/json"

	"github.com/sanjeevksaini/sovrunn/internal/apimeta"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/approvalreq"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/model"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/operation"
	"github.com/sanjeevksaini/sovrunn/internal/govaccess/state"
)

type DomainOutcome struct {
	sealed bool
	code   string
}

func (o DomainOutcome) Sealed() bool { return o.sealed }
func (o DomainOutcome) Code() string { return o.code }

// ApprovalRequirementDeriver owns Membership approval derivation.
type ApprovalRequirementDeriver struct{}

func (ApprovalRequirementDeriver) DeriverKind() string { return "Membership" }

func (ApprovalRequirementDeriver) DeriveForMembership(approvalreq.MembershipInput) (model.ApprovalRequirement, bool) {
	req := model.NewApprovalRequirementNotRequired()
	return req, req.Valid()
}

type PreparedMembershipIntent struct {
	sealed               bool
	claim                state.ParticipantClaim
	resourceUID          string
	resourceName         string
	scopeRef             apimeta.ScopeRef
	spec                 model.MembershipSpec
	protected            model.MembershipProtected
	expectedVersion      string
	nextVersion          string
	assignmentEffectSize int
}

func (i PreparedMembershipIntent) Sealed() bool                             { return i.sealed }
func (i PreparedMembershipIntent) ParticipantClaim() state.ParticipantClaim { return i.claim }
func (i PreparedMembershipIntent) HasPublicationDerivedFields() bool        { return false }
func (i PreparedMembershipIntent) AssignmentEffectSize() int                { return i.assignmentEffectSize }

func NewPreparedMembershipIntent(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	spec model.MembershipSpec,
	protected model.MembershipProtected,
	expectedVersion string,
	nextVersion string,
	assignmentEffectSize int,
) (PreparedMembershipIntent, operation.MechanicalFailure) {
	if participantID == "" || resourceUID == "" || resourceName == "" {
		return PreparedMembershipIntent{}, operation.NewMechanicalFailure("membership_identity_incomplete")
	}
	if expectedVersion == "" || nextVersion == "" {
		return PreparedMembershipIntent{}, operation.NewMechanicalFailure("membership_versions_incomplete")
	}
	if assignmentEffectSize < 0 {
		return PreparedMembershipIntent{}, operation.NewMechanicalFailure("membership_assignment_effect_negative")
	}
	raw, err := json.Marshal(struct {
		ResourceUID          string                    `json:"resourceUid"`
		ResourceName         string                    `json:"resourceName"`
		ScopeRef             apimeta.ScopeRef          `json:"scopeRef"`
		Spec                 model.MembershipSpec      `json:"spec"`
		Protected            model.MembershipProtected `json:"protected"`
		ExpectedVersion      string                    `json:"expectedVersion"`
		NextVersion          string                    `json:"nextVersion"`
		AssignmentEffectSize int                       `json:"assignmentEffectSize"`
	}{
		ResourceUID:          resourceUID,
		ResourceName:         resourceName,
		ScopeRef:             scopeRef,
		Spec:                 spec,
		Protected:            protected,
		ExpectedVersion:      expectedVersion,
		NextVersion:          nextVersion,
		AssignmentEffectSize: assignmentEffectSize,
	})
	if err != nil {
		return PreparedMembershipIntent{}, operation.NewMechanicalFailure("membership_intent_marshal")
	}
	sum := sha256.Sum256(raw)
	versionSet, err := state.NewExpectedVersionSet([]state.ExpectedVersionPredicate{{
		ResourceUID: resourceUID, ResourceVersion: expectedVersion,
	}})
	if err != nil {
		return PreparedMembershipIntent{}, operation.NewMechanicalFailure("membership_versionset_invalid")
	}
	claim, err := state.NewParticipantClaim(state.MembershipOwner, participantID, sum[:], versionSet, state.OriginatingResult)
	if err != nil {
		return PreparedMembershipIntent{}, operation.NewMechanicalFailure("membership_claim_invalid")
	}
	return PreparedMembershipIntent{
		sealed:               true,
		claim:                claim,
		resourceUID:          resourceUID,
		resourceName:         resourceName,
		scopeRef:             scopeRef,
		spec:                 spec,
		protected:            protected,
		expectedVersion:      expectedVersion,
		nextVersion:          nextVersion,
		assignmentEffectSize: assignmentEffectSize,
	}, operation.MechanicalFailure{}
}
