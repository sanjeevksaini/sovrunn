package approval

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

// ConsumeDerivers proves the approval owner consumes typed deriver ports.
func ConsumeDerivers(
	membership approvalreq.MembershipRequirementDeriverPort,
	privileged approvalreq.PrivilegedRequirementDeriverPort,
	exception approvalreq.ExceptionRequirementDeriverPort,
) ([]model.ApprovalRequirement, operation.MechanicalFailure) {
	if membership == nil || privileged == nil || exception == nil {
		return nil, operation.NewMechanicalFailure("deriver_nil")
	}
	out := make([]model.ApprovalRequirement, 0, 3)
	req, ok := membership.DeriveForMembership(approvalreq.MembershipInput{})
	if !ok || !req.Valid() {
		return nil, operation.NewMechanicalFailure("membership_requirement_invalid")
	}
	out = append(out, req)
	req, ok = privileged.DeriveForPrivileged(approvalreq.PrivilegedInput{})
	if !ok || !req.Valid() {
		return nil, operation.NewMechanicalFailure("privileged_requirement_invalid")
	}
	out = append(out, req)
	req, ok = exception.DeriveForException(approvalreq.ExceptionInput{})
	if !ok || !req.Valid() {
		return nil, operation.NewMechanicalFailure("exception_requirement_invalid")
	}
	out = append(out, req)
	return out, operation.MechanicalFailure{}
}

// ValidateStagesQuorumSoD enforces bounded stage quorum and SoD.
func ValidateStagesQuorumSoD(eligibleCount int, quorum int, requester, decider model.PrincipalRef) bool {
	if eligibleCount <= 0 || quorum <= 0 || quorum > eligibleCount {
		return false
	}
	// SoD: requester cannot be decider.
	return !requester.Equal(decider)
}

type PreparedApprovalIntent struct {
	sealed          bool
	claim           state.ParticipantClaim
	resourceUID     string
	resourceName    string
	scopeRef        apimeta.ScopeRef
	policy          model.ApprovalPolicy
	request         model.ApprovalRequest
	expectedVersion string
	nextVersion     string
}

func (i PreparedApprovalIntent) Sealed() bool                             { return i.sealed }
func (i PreparedApprovalIntent) ParticipantClaim() state.ParticipantClaim { return i.claim }
func (i PreparedApprovalIntent) HasPublicationDerivedFields() bool        { return false }

func NewPreparedApprovalIntent(
	participantID state.ParticipantID,
	resourceUID string,
	resourceName string,
	scopeRef apimeta.ScopeRef,
	policy model.ApprovalPolicy,
	request model.ApprovalRequest,
	expectedVersion string,
	nextVersion string,
) (PreparedApprovalIntent, operation.MechanicalFailure) {
	if participantID == "" || resourceUID == "" || resourceName == "" {
		return PreparedApprovalIntent{}, operation.NewMechanicalFailure("approval_identity_incomplete")
	}
	if expectedVersion == "" || nextVersion == "" {
		return PreparedApprovalIntent{}, operation.NewMechanicalFailure("approval_versions_incomplete")
	}
	if policy.Spec.Version == "" || len(policy.Spec.ApproverEligibility) == 0 {
		return PreparedApprovalIntent{}, operation.NewMechanicalFailure("approval_policy_invalid")
	}
	if request.Spec.SubjectUID == "" {
		return PreparedApprovalIntent{}, operation.NewMechanicalFailure("approval_request_invalid")
	}
	raw, err := json.Marshal(struct {
		ResourceUID     string                `json:"resourceUid"`
		ResourceName    string                `json:"resourceName"`
		ScopeRef        apimeta.ScopeRef      `json:"scopeRef"`
		Policy          model.ApprovalPolicy  `json:"policy"`
		Request         model.ApprovalRequest `json:"request"`
		ExpectedVersion string                `json:"expectedVersion"`
		NextVersion     string                `json:"nextVersion"`
	}{
		ResourceUID:     resourceUID,
		ResourceName:    resourceName,
		ScopeRef:        scopeRef,
		Policy:          policy,
		Request:         request,
		ExpectedVersion: expectedVersion,
		NextVersion:     nextVersion,
	})
	if err != nil {
		return PreparedApprovalIntent{}, operation.NewMechanicalFailure("approval_intent_marshal")
	}
	sum := sha256.Sum256(raw)
	versionSet, err := state.NewExpectedVersionSet([]state.ExpectedVersionPredicate{{
		ResourceUID: resourceUID, ResourceVersion: expectedVersion,
	}})
	if err != nil {
		return PreparedApprovalIntent{}, operation.NewMechanicalFailure("approval_versionset_invalid")
	}
	claim, err := state.NewParticipantClaim(state.ApprovalOwner, participantID, sum[:], versionSet, state.OriginatingResult)
	if err != nil {
		return PreparedApprovalIntent{}, operation.NewMechanicalFailure("approval_claim_invalid")
	}
	return PreparedApprovalIntent{
		sealed:          true,
		claim:           claim,
		resourceUID:     resourceUID,
		resourceName:    resourceName,
		scopeRef:        scopeRef,
		policy:          policy,
		request:         request,
		expectedVersion: expectedVersion,
		nextVersion:     nextVersion,
	}, operation.MechanicalFailure{}
}
